package main

import (
	"encoding/json"
	"strings"
	"testing"
)

func readyTradeInsightEvidence() ProviderCapabilityLifecycleEvidence {
	return ProviderCapabilityLifecycleEvidence{
		Observations: 40, Successes: 40, P95LatencyMs: 900, FreshnessState: "FRESH",
		SchemaIntegrity: true, FallbackVerified: true, CorroborationSamples: 12, DisagreementPct: 1,
		IndependentProvenance: true, ConsumerUtilityProven: true, TruthBoundaryVerified: true,
	}
}

func TestADAPTProviderLifecycleReadyNeverAutoPromotes(t *testing.T) {
	policy := tradeInsightProductionReadinessPolicy("daily-history", canonicalHistoricalBarsDataset)
	if policy.Lifecycle != providerLifecycleShadow {
		t.Fatalf("fixture must begin SHADOW: %+v", policy)
	}
	got := EvaluateProviderCapabilityReadiness(policy, readyTradeInsightEvidence())
	if got.Readiness != providerReadinessReadyForPromotion {
		t.Fatalf("complete evidence should be ready for explicit promotion: %+v", got)
	}
	if got.Lifecycle != providerLifecycleShadow || policy.Lifecycle != providerLifecycleShadow {
		t.Fatalf("readiness evaluation must never mutate lifecycle: policy=%+v diagnostic=%+v", policy, got)
	}
	if got.PromotionMode != "EXPLICIT_GOVERNED_ONLY" {
		t.Fatalf("promotion must remain governed: %+v", got)
	}
}

func TestADAPTProviderLifecycleInsufficientEvidenceIsNotReady(t *testing.T) {
	policy := tradeInsightProductionReadinessPolicy("daily-history", canonicalHistoricalBarsDataset)
	got := EvaluateProviderCapabilityReadiness(policy, ProviderCapabilityLifecycleEvidence{Observations: 2, Successes: 2})
	if got.Readiness != providerReadinessInsufficientEvidence {
		t.Fatalf("shallow observations must stay insufficient: %+v", got)
	}
	if len(got.Reasons) == 0 {
		t.Fatal("insufficient evidence must explain non-admission")
	}
}

func TestADAPTProviderLifecycleAuthAndQualityFailuresBlockReadiness(t *testing.T) {
	policy := tradeInsightProductionReadinessPolicy("adjusted-history", canonicalHistoricalBarsDataset)
	ev := readyTradeInsightEvidence()
	ev.AuthFailures = 1
	ev.RateLimited = 8
	ev.DisagreementPct = 25
	got := EvaluateProviderCapabilityReadiness(policy, ev)
	if got.Readiness == providerReadinessReadyForPromotion {
		t.Fatalf("auth/rate-limit/disagreement failures must block readiness: %+v", got)
	}
	joined := strings.Join(got.Reasons, " | ")
	for _, want := range []string{"auth failures", "rate limited", "disagreement"} {
		if !strings.Contains(joined, want) {
			t.Fatalf("missing deterministic reason %q: %+v", want, got)
		}
	}
}

func TestADAPTProviderLifecycleDirectAuthorityIsNotRankPromoted(t *testing.T) {
	policy := providerLifecyclePolicy("SEC EDGAR", "Form 4", "SEC")
	got := EvaluateProviderCapabilityReadiness(policy, ProviderCapabilityLifecycleEvidence{})
	if got.Readiness != providerReadinessNotApplicable || !policy.DirectAuthority {
		t.Fatalf("direct SEC authority must be lifecycle N/A: policy=%+v diagnostic=%+v", policy, got)
	}
	if got.Lifecycle != providerLifecycleProduction {
		t.Fatalf("direct authority production status changed: %+v", got)
	}
}

func TestADAPTProviderLifecycleRuntimeSuppressionDoesNotChangeGovernedState(t *testing.T) {
	policy := tradeInsightProductionReadinessPolicy("daily-history", canonicalHistoricalBarsDataset)
	cap := ProviderCapabilityStateRecord{Provider: tradeInsightProviderName, Dataset: canonicalHistoricalBarsDataset, State: providerCapabilityRateLimited, FailureCount: 4, SuccessCount: 30}
	telemetry := ProviderRequestDiagnostics{Provider: tradeInsightProviderName, Requests: 34, Successes: 30, Errors: 4, RateLimited: 4, P95LatencyMs: 1200}
	ev := lifecycleEvidenceFromRuntime(cap, telemetry, "FRESH")
	got := EvaluateProviderCapabilityReadiness(policy, ev)
	if got.Lifecycle != providerLifecycleShadow {
		t.Fatalf("runtime suppression/cooldown must not promote or roll back governed lifecycle: %+v", got)
	}
	if got.Readiness == providerReadinessReadyForPromotion {
		t.Fatalf("runtime-only evidence without schema/fallback/provenance/utility proof must not be ready: %+v", got)
	}
}

func TestADAPTProviderLifecycleDiagnosticCannotExposeSecrets(t *testing.T) {
	got := EvaluateProviderCapabilityReadiness(
		tradeInsightProductionReadinessPolicy("daily-history", canonicalHistoricalBarsDataset),
		readyTradeInsightEvidence(),
	)
	raw, err := json.Marshal(got)
	if err != nil {
		t.Fatal(err)
	}
	text := strings.ToLower(string(raw))
	for _, forbidden := range []string{"api_key", "apikey", "authorization", "secret", "token"} {
		if strings.Contains(text, forbidden) {
			t.Fatalf("lifecycle diagnostic exposed credential-shaped field %q: %s", forbidden, raw)
		}
	}
}

func TestADAPTTradeInsightAdmissionUsesSharedLifecycleVocabulary(t *testing.T) {
	for _, id := range []string{"daily-history", "adjusted-history", "corporate-actions", "bulk-history", "congressional-trades"} {
		row, ok := tradeInsightCapabilityAdmissionLookup(id)
		if !ok {
			t.Fatalf("TradeInsight capability %s missing", id)
		}
		if canonicalProviderLifecycle(row.Lifecycle) != providerLifecycleShadow || row.Lifecycle != providerLifecycleShadow {
			t.Fatalf("TradeInsight %s must enter common lifecycle at SHADOW: %+v", id, row)
		}
	}
}
