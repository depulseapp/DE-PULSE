package providerlifecycle

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func readyEvidence() Evidence {
	return Evidence{
		Observations: 40, Successes: 40, P95LatencyMs: 900, FreshnessState: "FRESH",
		SchemaIntegrity: true, FallbackVerified: true, CorroborationSamples: 12, DisagreementPct: 1,
		IndependentProvenance: true, ConsumerUtilityProven: true, TruthBoundaryVerified: true,
	}
}

func TestReadyNeverAutoPromotes(t *testing.T) {
	policy := TradeInsightPolicy("daily-history", "Historical Bars")
	got := Evaluate(policy, readyEvidence())
	if got.Readiness != ReadinessReadyForPromotion {
		t.Fatalf("complete evidence should be ready for explicit promotion: %+v", got)
	}
	if got.Lifecycle != Shadow || policy.Lifecycle != Shadow || got.PromotionMode != "EXPLICIT_GOVERNED_ONLY" {
		t.Fatalf("readiness must never mutate lifecycle: policy=%+v diagnostic=%+v", policy, got)
	}
}

func TestInsufficientEvidenceIsNotReady(t *testing.T) {
	got := Evaluate(TradeInsightPolicy("daily-history", "Historical Bars"), Evidence{Observations: 2, Successes: 2})
	if got.Readiness != ReadinessInsufficientEvidence || len(got.Reasons) == 0 {
		t.Fatalf("shallow observations must stay insufficient with reasons: %+v", got)
	}
}

func TestAuthRateLimitAndDisagreementBlockReadiness(t *testing.T) {
	policy := TradeInsightPolicy("adjusted-history", "Historical Bars")
	ev := readyEvidence()
	ev.AuthFailures = 1
	ev.RateLimited = 8
	ev.DisagreementPct = 25
	got := Evaluate(policy, ev)
	if got.Readiness == ReadinessReadyForPromotion {
		t.Fatalf("quality failures must block readiness: %+v", got)
	}
	joined := strings.Join(got.Reasons, " | ")
	for _, want := range []string{"auth failures", "rate limited", "disagreement"} {
		if !strings.Contains(joined, want) {
			t.Fatalf("missing deterministic reason %q: %+v", want, got)
		}
	}
}

func TestDirectAuthorityIsNotRankPromoted(t *testing.T) {
	policy := Policy{
		Provider: "SEC EDGAR", Capability: "Form 4", Dataset: "SEC", Lifecycle: Production,
		AuthorityClass: "DIRECT_AUTHORITY", PromotionMode: "EXPLICIT_GOVERNED_ONLY", DirectAuthority: true,
	}
	got := Evaluate(policy, Evidence{})
	if got.Readiness != ReadinessNotApplicable || got.Lifecycle != Production || !policy.DirectAuthority {
		t.Fatalf("direct authority must be production but rank-promotion N/A: policy=%+v diagnostic=%+v", policy, got)
	}
}

func TestRuntimePressureEvidenceCannotChangeGovernedLifecycle(t *testing.T) {
	policy := TradeInsightPolicy("daily-history", "Historical Bars")
	ev := Evidence{Observations: 34, Successes: 30, Errors: 4, RateLimited: 4, P95LatencyMs: 1200, FreshnessState: "FRESH"}
	got := Evaluate(policy, ev)
	if got.Lifecycle != Shadow {
		t.Fatalf("runtime suppression/cooldown evidence must not alter lifecycle: %+v", got)
	}
	if got.Readiness == ReadinessReadyForPromotion {
		t.Fatalf("runtime-only evidence without schema/fallback/provenance/utility proof must not be ready: %+v", got)
	}
}

func TestDiagnosticCannotExposeCredentialFields(t *testing.T) {
	got := Evaluate(TradeInsightPolicy("daily-history", "Historical Bars"), readyEvidence())
	raw, err := json.Marshal(got)
	if err != nil {
		t.Fatal(err)
	}
	text := strings.ToLower(string(raw))
	for _, forbidden := range []string{"api_key", "apikey", "authorization", "secret", "token"} {
		if strings.Contains(text, forbidden) {
			t.Fatalf("diagnostic exposed credential-shaped field %q: %s", forbidden, raw)
		}
	}
}

type lifecycleRegistry struct {
	Schema       string `json:"schema"`
	ScopeID      string `json:"scopeId"`
	Issue        int    `json:"issue"`
	Capabilities []struct {
		Provider       string `json:"provider"`
		Capability     string `json:"capability"`
		Lifecycle      string `json:"lifecycle"`
		Authority      string `json:"authority"`
		EvidenceStatus string `json:"evidenceStatus"`
		Readiness      string `json:"readiness"`
	} `json:"capabilities"`
}

type baselineMatrix struct {
	Providers []struct {
		Provider     string `json:"provider"`
		Capabilities []struct {
			Capability string `json:"capability"`
		} `json:"capabilities"`
	} `json:"providers"`
}

func repoFile(parts ...string) string {
	all := append([]string{"..", ".."}, parts...)
	return filepath.Join(all...)
}

func loadJSON[T any](t *testing.T, path string) T {
	t.Helper()
	raw, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	var out T
	if err := json.Unmarshal(raw, &out); err != nil {
		t.Fatalf("invalid %s: %v", path, err)
	}
	return out
}

func TestRegistryExhaustsBaselineCapabilityMatrix(t *testing.T) {
	matrix := loadJSON[baselineMatrix](t, repoFile("governance", "data-health", "provider-capability-matrix.json"))
	registry := loadJSON[lifecycleRegistry](t, repoFile("governance", "data-health", "provider-lifecycle-readiness.json"))
	if registry.Schema != "DE.PULSE-PROVIDER-CAPABILITY-LIFECYCLE-READINESS-1" || registry.ScopeID != "ADAPT-PROVIDER-LIFECYCLE-001" || registry.Issue != 83 {
		t.Fatalf("lifecycle registry identity drift: %+v", registry)
	}
	want := map[string]bool{}
	for _, provider := range matrix.Providers {
		for _, cap := range provider.Capabilities {
			want[provider.Provider+"::"+cap.Capability] = true
		}
	}
	got := map[string]bool{}
	tradeInsightRows := 0
	for _, row := range registry.Capabilities {
		key := row.Provider + "::" + row.Capability
		if got[key] {
			t.Fatalf("duplicate lifecycle capability: %s", key)
		}
		got[key] = true
		if row.Authority == "" || row.EvidenceStatus == "" || row.Readiness == "" {
			t.Fatalf("lifecycle row missing authority/evidence/readiness: %+v", row)
		}
		switch row.Lifecycle {
		case Shadow, Validated, Approved, Production:
		default:
			t.Fatalf("unsupported governed lifecycle %q for %s", row.Lifecycle, key)
		}
		if row.Provider == "TradeInsight" {
			tradeInsightRows++
			if row.Lifecycle != Shadow || row.EvidenceStatus != "COMMON_READINESS_REQUIRED" {
				t.Fatalf("TradeInsight must remain common SHADOW-first before #78: %+v", row)
			}
		}
	}
	if len(want) != 26 || len(got) != len(want) || tradeInsightRows != 3 {
		t.Fatalf("baseline/lifecycle coverage drift: baseline=%d lifecycle=%d tradeInsight=%d", len(want), len(got), tradeInsightRows)
	}
	for key := range want {
		if !got[key] {
			t.Fatalf("baseline capability missing lifecycle/readiness classification: %s", key)
		}
	}
	for key := range got {
		if !want[key] {
			t.Fatalf("lifecycle registry contains unowned capability outside #80 baseline: %s", key)
		}
	}
}
