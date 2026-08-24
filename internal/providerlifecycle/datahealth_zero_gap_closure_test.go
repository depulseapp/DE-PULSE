package providerlifecycle

import (
	"encoding/json"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

type zeroGapClosure struct {
	Schema                        string `json:"schema"`
	ScopeID                       string `json:"scopeId"`
	Issue                         int    `json:"issue"`
	GeneralRoutingAuthority       string `json:"generalRoutingAuthority"`
	NoExecutionBoundary           string `json:"noExecutionBoundary"`
	DirectSECEDGARForm4Authority  string `json:"directSECEDGARForm4Authority"`
	TradeInsightProductionPromotion struct {
		ApprovedByIssue              int      `json:"approvedByIssue"`
		ProductionGateIssue          int      `json:"productionGateIssue"`
		PromotedCapabilities         []string `json:"promotedCapabilities"`
		HardGated                    []string `json:"hardGated"`
		PromotionMode                string   `json:"promotionMode"`
		ProductionActivationCondition string   `json:"productionActivationCondition"`
	} `json:"tradeInsightProductionPromotion"`
	EvidenceBundles map[string]struct {
		RegressionIDs             []string `json:"regressionIds"`
		RuntimeEvidence           string   `json:"runtimeEvidence"`
		SupportedPlatformEvidence []string `json:"supportedPlatformEvidence"`
	} `json:"evidenceBundles"`
	RowDispositions []struct {
		Provider         string `json:"provider"`
		Capability       string `json:"capability"`
		FinalDisposition string `json:"finalDisposition"`
		EvidenceBundle   string `json:"evidenceBundle"`
	} `json:"rowDispositions"`
	FaultRecoveryMatrix []struct {
		ID                  string   `json:"id"`
		RegressionIDs       []string `json:"regressionIds"`
		ExpectedDisposition string   `json:"expectedDisposition"`
	} `json:"faultRecoveryMatrix"`
	RealisticWorkloadEvidence []string `json:"realisticWorkloadEvidence"`
	NativeEvidenceContracts   []string `json:"nativeEvidenceContracts"`
	ExactHeadEvidencePolicy   string   `json:"exactHeadEvidencePolicy"`
	ResidualsPolicy           string   `json:"residualsPolicy"`
}

type zeroGapBaseline struct {
	Providers []struct {
		Provider     string `json:"provider"`
		Capabilities []struct {
			Capability     string   `json:"capability"`
			CanonicalOwner string   `json:"canonicalOwner"`
			Consumers      []string `json:"consumers"`
		} `json:"capabilities"`
	} `json:"providers"`
}

type nativeProofContract struct {
	Schema        string   `json:"schema"`
	ScopeID       string   `json:"scopeId"`
	Issue         int      `json:"issue"`
	Platform      string   `json:"platform"`
	QualifiedJob  string   `json:"qualifiedJob"`
	EvidenceOwner string   `json:"evidenceOwner"`
	RequiredChecks []string `json:"requiredChecks"`
	Binding       string   `json:"binding"`
}

func loadZeroGapClosure(t *testing.T) zeroGapClosure {
	t.Helper()
	return loadJSON[zeroGapClosure](t, repoFile("governance", "data-health", "datahealth-zero-gap-closure.json"))
}

func currentRegressionSource(t *testing.T) string {
	t.Helper()
	root := repoFile()
	var b strings.Builder
	err := filepath.WalkDir(root, func(path string, entry fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if entry.IsDir() {
			rel, relErr := filepath.Rel(root, path)
			if relErr == nil {
				prefix := filepath.ToSlash(rel)
				if strings.HasPrefix(prefix, "release/") || strings.HasPrefix(prefix, "internal/vendorcrypto/") {
					return filepath.SkipDir
				}
			}
			return nil
		}
		if !strings.HasSuffix(entry.Name(), "_test.go") {
			return nil
		}
		raw, readErr := os.ReadFile(path)
		if readErr != nil {
			return readErr
		}
		b.Write(raw)
		b.WriteByte('\n')
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
	return b.String()
}

func requireRegressionIDsExist(t *testing.T, source string, ids []string) {
	t.Helper()
	if len(ids) == 0 {
		t.Fatal("regression evidence must not be empty")
	}
	for _, id := range ids {
		if strings.TrimSpace(id) == "" || !strings.Contains(source, "func "+id+"(") {
			t.Fatalf("zero-gap evidence references missing current regression %q", id)
		}
	}
}

func TestDataHealth84ClosureExhaustsAllBaselineCapabilityRows(t *testing.T) {
	closure := loadZeroGapClosure(t)
	if closure.Schema != "DE.PULSE-DATAHEALTH-ZERO-GAP-CLOSURE-1" || closure.ScopeID != "ADAPT-DATAHEALTH-CLOSURE-001" || closure.Issue != 84 {
		t.Fatalf("zero-gap closure identity drift: %+v", closure)
	}
	if closure.GeneralRoutingAuthority != "Smart Provider Router v2" || closure.NoExecutionBoundary != "PRESERVED" || closure.DirectSECEDGARForm4Authority != "PRESERVED" {
		t.Fatalf("permanent architecture boundary drift: %+v", closure)
	}

	baseline := loadJSON[zeroGapBaseline](t, repoFile("governance", "data-health", "provider-capability-matrix.json"))
	want := map[string]bool{}
	for _, provider := range baseline.Providers {
		for _, capability := range provider.Capabilities {
			key := provider.Provider + "::" + capability.Capability
			if capability.CanonicalOwner == "" || len(capability.Consumers) == 0 {
				t.Fatalf("baseline row lacks exact source owner/consumers: %s", key)
			}
			want[key] = true
		}
	}
	if len(want) != 26 {
		t.Fatalf("baseline capability count=%d want=26", len(want))
	}

	got := map[string]bool{}
	for _, row := range closure.RowDispositions {
		key := row.Provider + "::" + row.Capability
		if got[key] {
			t.Fatalf("duplicate zero-gap row disposition: %s", key)
		}
		got[key] = true
		if !want[key] {
			t.Fatalf("zero-gap row is outside canonical baseline: %s", key)
		}
		if row.FinalDisposition != "PRODUCTION_SUPPORTED" {
			t.Fatalf("row %s final disposition=%q", key, row.FinalDisposition)
		}
		bundle, ok := closure.EvidenceBundles[row.EvidenceBundle]
		if !ok || strings.TrimSpace(bundle.RuntimeEvidence) == "" || len(bundle.RegressionIDs) == 0 {
			t.Fatalf("row %s has incomplete evidence bundle %q", key, row.EvidenceBundle)
		}
		platforms := strings.Join(bundle.SupportedPlatformEvidence, " | ")
		for _, required := range []string{"Qualified macOS native lifecycle rehearsal", "Qualified Windows native runtime rehearsal"} {
			if !strings.Contains(platforms, required) {
				t.Fatalf("row %s missing supported-platform evidence %q", key, required)
			}
		}
	}
	if len(got) != len(want) {
		t.Fatalf("zero-gap row coverage=%d want=%d", len(got), len(want))
	}
	for key := range want {
		if !got[key] {
			t.Fatalf("baseline capability has no final zero-gap disposition: %s", key)
		}
	}
}

func TestDataHealth84FaultRecoveryMatrixIsExecutableAndComplete(t *testing.T) {
	closure := loadZeroGapClosure(t)
	regressionSource := currentRegressionSource(t)
	required := map[string]bool{
		"missing_invalid_credential":                 true,
		"auth_401_403":                              true,
		"rate_limit_429_quota":                      true,
		"server_5xx":                                true,
		"timeout_unreachable_offline":               true,
		"malformed_partial_schema_drift":            true,
		"stale_temporally_invalid":                  true,
		"preferred_failure_healthy_fallback":        true,
		"fallback_exhaustion":                       true,
		"contradictory_independent_providers":       true,
		"cache_miss_hit_expired_warm_recovery":      true,
		"local_pressure":                            true,
		"restart_recovery":                          true,
		"optional_provider_failure_isolation":       true,
		"provider_recovery_hysteresis_no_flapping":  true,
	}
	if len(closure.FaultRecoveryMatrix) != len(required) {
		t.Fatalf("fault/recovery matrix=%d want=%d", len(closure.FaultRecoveryMatrix), len(required))
	}
	seen := map[string]bool{}
	for _, row := range closure.FaultRecoveryMatrix {
		if !required[row.ID] || seen[row.ID] {
			t.Fatalf("unexpected/duplicate fault matrix row %q", row.ID)
		}
		seen[row.ID] = true
		if row.ExpectedDisposition != "TRUTHFUL_SCOPED_RECOVERY_OR_FAIL_CLOSED" {
			t.Fatalf("fault %s has unsafe disposition %q", row.ID, row.ExpectedDisposition)
		}
		requireRegressionIDsExist(t, regressionSource, row.RegressionIDs)
	}
	for id := range required {
		if !seen[id] {
			t.Fatalf("required fault/recovery scenario missing: %s", id)
		}
	}
	for name, bundle := range closure.EvidenceBundles {
		if strings.TrimSpace(bundle.RuntimeEvidence) == "" {
			t.Fatalf("evidence bundle %s missing runtime evidence", name)
		}
		requireRegressionIDsExist(t, regressionSource, bundle.RegressionIDs)
	}
	requireRegressionIDsExist(t, regressionSource, closure.RealisticWorkloadEvidence)
}

func TestDataHealth84NativeProofContractsRequireExactHeadQualifiedJobs(t *testing.T) {
	closure := loadZeroGapClosure(t)
	if len(closure.NativeEvidenceContracts) != 2 {
		t.Fatalf("native proof contracts=%d want=2", len(closure.NativeEvidenceContracts))
	}
	workflowRaw, err := os.ReadFile(repoFile(".github", "workflows", "ci-qualified.yml"))
	if err != nil {
		t.Fatal(err)
	}
	workflow := string(workflowRaw)
	platforms := map[string]string{
		"macOS Apple Silicon": "Qualified macOS native lifecycle rehearsal",
		"Windows x64":         "Qualified Windows native runtime rehearsal",
	}
	seen := map[string]bool{}
	for _, rel := range closure.NativeEvidenceContracts {
		contract := loadJSON[nativeProofContract](t, repoFile(strings.Split(filepath.ToSlash(rel), "/")...))
		if contract.Schema != "DE.PULSE-DATAHEALTH-NATIVE-PROOF-1" || contract.ScopeID != closure.ScopeID || contract.Issue != 84 {
			t.Fatalf("native proof identity drift in %s: %+v", rel, contract)
		}
		wantJob, ok := platforms[contract.Platform]
		if !ok || contract.QualifiedJob != wantJob || seen[contract.Platform] {
			t.Fatalf("native platform/job drift in %s: %+v", rel, contract)
		}
		seen[contract.Platform] = true
		if contract.Binding != "EXACT_HEAD_QUALIFIED_PASS_REQUIRED" || len(contract.RequiredChecks) == 0 {
			t.Fatalf("native proof must bind exact-head Qualified and concrete checks: %+v", contract)
		}
		if _, err := os.Stat(repoFile(strings.Split(filepath.ToSlash(contract.EvidenceOwner), "/")...)); err != nil {
			t.Fatalf("native evidence owner missing for %s: %v", contract.Platform, err)
		}
		if !strings.Contains(workflow, contract.QualifiedJob) {
			t.Fatalf("Qualified workflow does not own required native job %q", contract.QualifiedJob)
		}
	}
	if len(seen) != len(platforms) {
		t.Fatalf("supported native platform coverage=%d want=%d", len(seen), len(platforms))
	}
	if !strings.Contains(closure.ExactHeadEvidencePolicy, "GitHub Actions") || !strings.Contains(closure.ResidualsPolicy, "Zero unexplained") {
		t.Fatalf("post-commit evidence/residual policy drift: %+v", closure)
	}
}
