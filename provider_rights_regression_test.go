package main

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
	"time"
)

func clearProviderRightsBundleForTest(t *testing.T) {
	t.Helper()
	t.Setenv(providerRightsBundlePathEnv, "")
	t.Setenv(providerRightsBundleSHA256Env, "")
}

func approvedHostedRightsFixtureFor(provider string) ProviderDataRightsMetadata {
	rights := ProviderDataRightsMetadata{
		PolicyVersion:       providerDataRightsPolicyVersion,
		Provider:            provider,
		ReviewState:         providerRightsApproved,
		CommercialUse:       providerRightsApproved,
		MultiUserUse:        providerRightsApproved,
		Proxying:            providerRightsApproved,
		CachingRetention:    providerRightsApproved,
		Redistribution:      providerRightsApproved,
		Display:             providerRightsApproved,
		AIUse:               providerRightsApproved,
		UsageLimits:         "governed provider contract limits recorded",
		Attribution:         "provider attribution requirements recorded",
		AllowedEnvironments: []string{"stage", "prod"},
		EffectiveAt:         "2026-08-01T00:00:00Z",
		ExpiresAt:           "2026-12-01T00:00:00Z",
		RenewalState:        providerRightsRenewalCurrent,
		EvidenceBound:       true,
		EvidenceRef:         "rights/provider/example/2026-08",
		EvidenceDigest:      "sha256:" + strings.Repeat("a", 64),
		ReviewedAt:          "2026-08-20T00:00:00Z",
		Detail:              "fixture-only reviewed rights record",
	}
	rights.CommercialReadiness = evaluateProviderCommercialReadiness(rights)
	return rights
}

func approvedHostedRightsFixture() ProviderDataRightsMetadata {
	return approvedHostedRightsFixtureFor("Finnhub")
}

func bindHostedRightsBundleForTest(t *testing.T, records ...ProviderDataRightsMetadata) string {
	t.Helper()
	t.Setenv(runtimeModeEnv, "hosted")
	t.Setenv(hostedEnvironmentEnv, "prod")
	bundle := ProviderDataRightsBundle{PolicyVersion: providerRightsBundlePolicyVersion, Records: records}
	raw, err := json.Marshal(bundle)
	if err != nil {
		t.Fatal(err)
	}
	path := filepath.Join(t.TempDir(), "provider-rights.json")
	if err := os.WriteFile(path, raw, 0600); err != nil {
		t.Fatal(err)
	}
	sum := sha256.Sum256(raw)
	t.Setenv(providerRightsBundlePathEnv, path)
	t.Setenv(providerRightsBundleSHA256Env, "sha256:"+hex.EncodeToString(sum[:]))
	return path
}

func TestV184ProviderDataRightsDefaultFailsClosed(t *testing.T) {
	clearProviderRightsBundleForTest(t)
	providers := []string{"Alpaca", "Finnhub", "Twelve Data", "Marketaux", "FRED", "SEC EDGAR", "yfinance", "CBOE"}
	for _, provider := range providers {
		rights := providerDataRightsMetadata(provider)
		if rights.PolicyVersion != providerDataRightsPolicyVersion || rights.Provider != provider {
			t.Fatalf("%s provider/policy binding = %+v", provider, rights)
		}
		if rights.ReviewState != providerRightsUnreviewed {
			t.Fatalf("%s review state = %q; want conservative UNREVIEWED", provider, rights.ReviewState)
		}
		if rights.CommercialUse != providerRightsNotAsserted || rights.Redistribution != providerRightsNotAsserted || rights.AIUse != providerRightsNotAsserted {
			t.Fatalf("%s implicitly asserted rights: %+v", provider, rights)
		}
		if rights.EvidenceBound {
			t.Fatalf("%s claims rights evidence without a bound evidence record", provider)
		}
	}
}

func TestV184ProviderRightsAreSeparateFromOperationalEntitlement(t *testing.T) {
	clearProviderRightsBundleForTest(t)
	hop := ProviderRouteHop{Provider: "Finnhub", Entitlement: providerCapabilitySupported, DataRights: providerDataRightsMetadata("Finnhub")}
	raw, err := json.Marshal(hop)
	if err != nil {
		t.Fatal(err)
	}
	var got map[string]any
	if err := json.Unmarshal(raw, &got); err != nil {
		t.Fatal(err)
	}
	if got["entitlement"] != providerCapabilitySupported {
		t.Fatalf("operational entitlement changed: %#v", got["entitlement"])
	}
	rights, ok := got["dataRights"].(map[string]any)
	if !ok || rights["reviewState"] != providerRightsUnreviewed || rights["commercialUse"] != providerRightsNotAsserted {
		t.Fatalf("rights metadata not conservative/separate: %#v", got["dataRights"])
	}
}

func TestV184ProviderRightsDoNotChangeSmartRouterScore(t *testing.T) {
	cap := ProviderCapabilityStateRecord{Provider: "Finnhub", Dataset: "News", State: providerCapabilitySupported}
	circuit := providerCircuit{}
	telemetry := ProviderRequestDiagnostics{Provider: "Finnhub", Successes: 10}
	before := smartRouteScore("Finnhub", "News", 1, WorkTierUserActionable, cap, circuit, telemetry, "regular")
	_ = evaluateProviderHostedRightsDecision(approvedHostedRightsFixture(), providerHostedUseProductionServing, "prod", time.Date(2026, 8, 25, 12, 0, 0, 0, time.UTC))
	after := smartRouteScore("Finnhub", "News", 1, WorkTierUserActionable, cap, circuit, telemetry, "regular")
	if !reflect.DeepEqual(before, after) {
		t.Fatalf("rights governance mutated Smart Router score: before=%+v after=%+v", before, after)
	}
}

func TestHOST001ProviderRightsMetadataCoversHostedLegalDimensions(t *testing.T) {
	clearProviderRightsBundleForTest(t)
	rights := providerDataRightsMetadata("Finnhub")
	for name, value := range map[string]string{
		"provider":    rights.Provider,
		"commercial":  rights.CommercialUse,
		"multi-user":  rights.MultiUserUse,
		"proxy":       rights.Proxying,
		"cache":       rights.CachingRetention,
		"redisplay":   rights.Redistribution,
		"display":     rights.Display,
		"AI/derived":  rights.AIUse,
		"renewal":     rights.RenewalState,
	} {
		if value == "" {
			t.Fatalf("%s rights dimension is missing", name)
		}
	}
	if rights.EvidenceBound || len(rights.AllowedEnvironments) != 0 || rights.EffectiveAt != "" || rights.ExpiresAt != "" {
		t.Fatalf("unreviewed provider fabricated hosted rights evidence: %+v", rights)
	}
}

func TestHOST002WorkingProviderKeyNeverGrantsHostedRights(t *testing.T) {
	t.Setenv(runtimeModeEnv, "hosted")
	t.Setenv(hostedEnvironmentEnv, "prod")
	clearProviderRightsBundleForTest(t)
	rights := providerDataRightsMetadata("Finnhub")
	now := time.Date(2026, 8, 25, 12, 0, 0, 0, time.UTC)
	decision := evaluateProviderHostedRightsDecision(rights, providerHostedUseCommercialMultiUser, "prod", now)
	if decision.Allowed || decision.State != providerHostedRightsBlocked {
		t.Fatalf("unbound provider rights unexpectedly allowed hosted use: %+v", decision)
	}
	if decision.EvidenceRef != "" || decision.EvidenceDigest != "" {
		t.Fatalf("default rights fabricated provenance: %+v", decision)
	}
}

func TestHOST002HostedRightsRequireBoundReviewableProvenance(t *testing.T) {
	rights := approvedHostedRightsFixture()
	now := time.Date(2026, 8, 25, 12, 0, 0, 0, time.UTC)
	for _, mutate := range []func(*ProviderDataRightsMetadata){
		func(r *ProviderDataRightsMetadata) { r.Provider = "" },
		func(r *ProviderDataRightsMetadata) { r.EvidenceBound = false },
		func(r *ProviderDataRightsMetadata) { r.EvidenceRef = "" },
		func(r *ProviderDataRightsMetadata) { r.EvidenceDigest = "not-a-sha" },
		func(r *ProviderDataRightsMetadata) { r.ReviewState = providerRightsUnreviewed },
	} {
		candidate := rights
		mutate(&candidate)
		decision := evaluateProviderHostedRightsDecision(candidate, providerHostedUseCommercialMultiUser, "prod", now)
		if decision.Allowed {
			t.Fatalf("missing provider/review/provenance unexpectedly allowed: %+v", candidate)
		}
	}
}

func TestHOST002ProviderRightsBundleIsSHA256Pinned(t *testing.T) {
	path := bindHostedRightsBundleForTest(t, approvedHostedRightsFixture())
	if got := providerDataRightsMetadata("Finnhub"); !got.EvidenceBound || got.Provider != "Finnhub" {
		t.Fatalf("valid pinned bundle did not load reviewed provider record: %+v", got)
	}
	if err := os.WriteFile(path, []byte(`{"tampered":true}`), 0600); err != nil {
		t.Fatal(err)
	}
	if got := providerDataRightsMetadata("Finnhub"); got.EvidenceBound || got.ReviewState != providerRightsUnreviewed {
		t.Fatalf("tampered bundle remained trusted: %+v", got)
	}
}

func TestHOST003HostedRightsPurposeEligibilityIsFailClosed(t *testing.T) {
	rights := approvedHostedRightsFixture()
	now := time.Date(2026, 8, 25, 12, 0, 0, 0, time.UTC)
	for _, purpose := range []string{
		providerHostedUseCommercialMultiUser,
		providerHostedUseProxy,
		providerHostedUseCacheRetention,
		providerHostedUseRedisplay,
		providerHostedUseAI,
		providerHostedUseLiveFanout,
		providerHostedUseProductionServing,
	} {
		decision := evaluateProviderHostedRightsDecision(rights, purpose, "prod", now)
		if !decision.Allowed || decision.State != providerHostedRightsAllowed {
			t.Fatalf("approved %s purpose blocked: %+v", purpose, decision)
		}
	}
	if decision := evaluateProviderHostedRightsDecision(rights, "UNKNOWN", "prod", now); decision.Allowed {
		t.Fatalf("unknown hosted purpose must fail closed: %+v", decision)
	}
}

func TestHOST003RightsExpiryDowngradeAndEnvironmentBlockDeterministically(t *testing.T) {
	rights := approvedHostedRightsFixture()
	beforeExpiry := time.Date(2026, 11, 30, 23, 59, 59, 0, time.UTC)
	atExpiry := time.Date(2026, 12, 1, 0, 0, 0, 0, time.UTC)
	if decision := evaluateProviderHostedRightsDecision(rights, providerHostedUseProductionServing, "prod", beforeExpiry); !decision.Allowed {
		t.Fatalf("valid rights blocked before expiry: %+v", decision)
	}
	if decision := evaluateProviderHostedRightsDecision(rights, providerHostedUseProductionServing, "prod", atExpiry); decision.Allowed {
		t.Fatalf("expired rights remained eligible: %+v", decision)
	}

	downgraded := rights
	downgraded.CachingRetention = providerRightsDenied
	if decision := evaluateProviderHostedRightsDecision(downgraded, providerHostedUseProductionServing, "prod", beforeExpiry); decision.Allowed {
		t.Fatalf("rights downgrade did not block production cache/persistence eligibility: %+v", decision)
	}
	if decision := evaluateProviderHostedRightsDecision(rights, providerHostedUseProductionServing, "dev", beforeExpiry); decision.Allowed {
		t.Fatalf("unapproved environment unexpectedly eligible: %+v", decision)
	}
}

func TestHOST003ExecutableRouterFailsClosedThenAdmitsReviewedRights(t *testing.T) {
	e := newV1801Engine(t)
	e.app.mu.Lock()
	e.app.secrets.Finnhub = "working-test-key"
	e.app.mu.Unlock()
	t.Setenv(runtimeModeEnv, "hosted")
	t.Setenv(hostedEnvironmentEnv, "prod")
	clearProviderRightsBundleForTest(t)

	called := 0
	attempts := map[string]providerRouteAttempt{
		"Finnhub": func(context.Context) bool { called++; return true },
	}
	if provider, ok := e.executeProviderRoute(context.Background(), "News", attempts); ok || provider != "" || called != 0 {
		t.Fatalf("unreviewed hosted provider was attempted: provider=%q ok=%v called=%d", provider, ok, called)
	}

	bindHostedRightsBundleForTest(t, approvedHostedRightsFixture())
	if provider, ok := e.executeProviderRoute(context.Background(), "News", attempts); !ok || provider != "Finnhub" || called != 1 {
		t.Fatalf("reviewed hosted provider was not reachable through canonical router: provider=%q ok=%v called=%d", provider, ok, called)
	}
}

func TestHOST003RouterObservabilityShowsRightsBlockWithoutChangingEntitlement(t *testing.T) {
	e := newV1801Engine(t)
	e.app.mu.Lock()
	e.app.secrets.Finnhub = "working-test-key"
	settings := clone(e.app.state.Settings)
	secrets := clone(e.app.secrets)
	e.app.mu.Unlock()
	t.Setenv(runtimeModeEnv, "hosted")
	t.Setenv(hostedEnvironmentEnv, "prod")
	clearProviderRightsBundleForTest(t)

	snap := e.buildProviderRouterSnapshot(settings, secrets, map[string]Quote{}, map[string]int64{})
	found := false
	for _, route := range snap.Routes {
		if route.Dataset != "News" {
			continue
		}
		for _, hop := range route.Route {
			if hop.Provider == "Finnhub" {
				found = true
				if hop.Health != "RIGHTS BLOCKED" || hop.Recovery != "SUPPRESSED" {
					t.Fatalf("rights denial not diagnosable in router: %+v", hop)
				}
				if hop.Entitlement == "RIGHTS BLOCKED" {
					t.Fatalf("legal rights were incorrectly folded into operational entitlement: %+v", hop)
				}
			}
		}
	}
	if !found {
		t.Fatal("Finnhub News hop missing from router diagnostics")
	}
}
