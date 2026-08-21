package main

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"os"
	"strings"
	"testing"
)

func TestV189DeveloperCommandRequiresExplicitOptIn(t *testing.T) {
	called := false
	probe := func(context.Context) (tradeInsightSchemaFingerprint, error) {
		called = true
		return tradeInsightSchemaFingerprint{}, nil
	}

	for _, args := range [][]string{
		nil,
		{},
		{"--unrelated-flag"},
	} {
		var stdout, stderr bytes.Buffer
		handled, code := runDeveloperCommandWithProbe(args, &stdout, &stderr, probe)
		if handled {
			t.Fatalf("args %v unexpectedly handled", args)
		}
		if code != 0 {
			t.Fatalf("args %v returned code %d, want 0", args, code)
		}
		if stdout.Len() != 0 || stderr.Len() != 0 {
			t.Fatalf("args %v produced output without explicit opt-in", args)
		}
	}
	if called {
		t.Fatal("Congress schema probe ran without explicit developer flag")
	}
}

func TestV189DeveloperCommandRejectsAdditionalArgumentsWithoutProbing(t *testing.T) {
	called := false
	probe := func(context.Context) (tradeInsightSchemaFingerprint, error) {
		called = true
		return tradeInsightSchemaFingerprint{}, nil
	}
	var stdout, stderr bytes.Buffer
	handled, code := runDeveloperCommandWithProbe(
		[]string{tradeInsightCongressSchemaProbeFlag, "unexpected"},
		&stdout,
		&stderr,
		probe,
	)
	if !handled || code != 2 {
		t.Fatalf("handled=%v code=%d, want handled=true code=2", handled, code)
	}
	if called {
		t.Fatal("probe ran despite invalid developer command arguments")
	}
	if stdout.Len() != 0 {
		t.Fatalf("unexpected stdout: %q", stdout.String())
	}
	if !strings.Contains(stderr.String(), "does not accept additional arguments") {
		t.Fatalf("unexpected stderr: %q", stderr.String())
	}
}

func TestV189DeveloperCommandPrintsOnlyRedactedFingerprintAndStaysGated(t *testing.T) {
	fingerprint, err := tradeInsightSchemaFingerprintFromPayload([]byte(`[
		{"field_alpha":"secret-person-value","field_beta":123},
		{"field_alpha":"secret-symbol-value","field_beta":null}
	]`))
	if err != nil {
		t.Fatal(err)
	}
	before := tradeInsightCapabilityLifecycleTruth("congressional-trades")
	if before != "GATED" {
		t.Fatalf("precondition lifecycle=%q, want GATED", before)
	}

	called := 0
	probe := func(ctx context.Context) (tradeInsightSchemaFingerprint, error) {
		called++
		if ctx == nil {
			t.Fatal("probe context is nil")
		}
		return fingerprint, nil
	}
	var stdout, stderr bytes.Buffer
	handled, code := runDeveloperCommandWithProbe(
		[]string{tradeInsightCongressSchemaProbeFlag},
		&stdout,
		&stderr,
		probe,
	)
	if !handled || code != 0 {
		t.Fatalf("handled=%v code=%d stderr=%q", handled, code, stderr.String())
	}
	if called != 1 {
		t.Fatalf("probe calls=%d, want 1", called)
	}
	if stderr.Len() != 0 {
		t.Fatalf("unexpected stderr: %q", stderr.String())
	}

	output := stdout.String()
	for _, forbidden := range []string{"secret-person-value", "secret-symbol-value"} {
		if strings.Contains(output, forbidden) {
			t.Fatalf("developer command leaked provider value %q in %s", forbidden, output)
		}
	}
	for _, required := range []string{"field_alpha", "field_beta", "congressional-trades", "GATED"} {
		if !strings.Contains(output, required) {
			t.Fatalf("developer command output missing %q: %s", required, output)
		}
	}

	var report tradeInsightCongressSchemaProbeReport
	if err := json.Unmarshal(stdout.Bytes(), &report); err != nil {
		t.Fatalf("decode report: %v; output=%s", err, output)
	}
	if report.Capability != "congressional-trades" {
		t.Fatalf("capability=%q", report.Capability)
	}
	if report.Lifecycle != "GATED" {
		t.Fatalf("lifecycle=%q, want GATED", report.Lifecycle)
	}
	if report.Fingerprint.RowsObserved != 2 {
		t.Fatalf("rows observed=%d, want 2", report.Fingerprint.RowsObserved)
	}
	if after := tradeInsightCapabilityLifecycleTruth("congressional-trades"); after != before {
		t.Fatalf("probe command changed lifecycle from %q to %q", before, after)
	}
}

func TestV189DeveloperCommandFailsClosedOnProbeError(t *testing.T) {
	probe := func(context.Context) (tradeInsightSchemaFingerprint, error) {
		return tradeInsightSchemaFingerprint{}, errors.New("synthetic admission failure")
	}
	var stdout, stderr bytes.Buffer
	handled, code := runDeveloperCommandWithProbe(
		[]string{tradeInsightCongressSchemaProbeFlag},
		&stdout,
		&stderr,
		probe,
	)
	if !handled || code != 1 {
		t.Fatalf("handled=%v code=%d, want handled=true code=1", handled, code)
	}
	if stdout.Len() != 0 {
		t.Fatalf("unexpected stdout: %q", stdout.String())
	}
	if !strings.Contains(stderr.String(), "synthetic admission failure") {
		t.Fatalf("unexpected stderr: %q", stderr.String())
	}
	if lifecycle := tradeInsightCapabilityLifecycleTruth("congressional-trades"); lifecycle != "GATED" {
		t.Fatalf("failed probe changed lifecycle to %q", lifecycle)
	}
}

func TestV189DeveloperCommandFailsClosedWhenProbeUnavailable(t *testing.T) {
	var stdout, stderr bytes.Buffer
	handled, code := runDeveloperCommandWithProbe(
		[]string{tradeInsightCongressSchemaProbeFlag},
		&stdout,
		&stderr,
		nil,
	)
	if !handled || code != 1 {
		t.Fatalf("handled=%v code=%d, want handled=true code=1", handled, code)
	}
	if stdout.Len() != 0 {
		t.Fatalf("unexpected stdout: %q", stdout.String())
	}
	if !strings.Contains(stderr.String(), "probe unavailable") {
		t.Fatalf("unexpected stderr: %q", stderr.String())
	}
}

func TestV189DeveloperCommandGuardRunsBeforeApplicationInitialization(t *testing.T) {
	body, err := os.ReadFile("main.go")
	if err != nil {
		t.Fatalf("read main.go: %v", err)
	}
	probeGuard := bytes.Index(body, []byte("runDeveloperCommand(os.Args[1:]"))
	applicationInit := bytes.Index(body, []byte("NewApplication()"))
	if probeGuard < 0 {
		t.Fatal("main.go no longer contains the developer command guard")
	}
	if applicationInit < 0 {
		t.Fatal("main.go no longer contains NewApplication initialization")
	}
	if probeGuard > applicationInit {
		t.Fatal("developer command guard moved after NewApplication; admission probe could initialize runtime state")
	}
}
