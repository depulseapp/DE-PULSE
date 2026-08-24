package providerlifecycle

import (
	"fmt"
	"strings"
)

const (
	Shadow     = "SHADOW"
	Validated  = "VALIDATED"
	Approved   = "APPROVED"
	Production = "PRODUCTION"
)

const (
	ReadinessNotApplicable        = "N/A"
	ReadinessInsufficientEvidence = "INSUFFICIENT_EVIDENCE"
	ReadinessNotReady             = "NOT_READY"
	ReadinessReadyForPromotion    = "READY_FOR_PROMOTION"
	ReadinessProductionApproved   = "PRODUCTION_APPROVED"
)

// Policy is governance metadata only. Smart Provider Router v2 remains the
// runtime routing/admission authority; Evaluate never mutates Lifecycle.
type Policy struct {
	Provider                     string  `json:"provider"`
	Capability                   string  `json:"capability"`
	Dataset                      string  `json:"dataset,omitempty"`
	Lifecycle                    string  `json:"lifecycle"`
	AuthorityClass               string  `json:"authorityClass"`
	PromotionMode                string  `json:"promotionMode"`
	MinObservations              int64   `json:"minObservations,omitempty"`
	MinSuccessPct                float64 `json:"minSuccessPct,omitempty"`
	MaxP95LatencyMs              int64   `json:"maxP95LatencyMs,omitempty"`
	MaxAuthFailures              int64   `json:"maxAuthFailures,omitempty"`
	MaxRateLimitedPct            float64 `json:"maxRateLimitedPct,omitempty"`
	MinCorroborationSamples      int64   `json:"minCorroborationSamples,omitempty"`
	MaxDisagreementPct           float64 `json:"maxDisagreementPct,omitempty"`
	RequireFreshEvidence         bool    `json:"requireFreshEvidence,omitempty"`
	RequireSchemaIntegrity       bool    `json:"requireSchemaIntegrity,omitempty"`
	RequireFallbackProof         bool    `json:"requireFallbackProof,omitempty"`
	RequireIndependentProvenance bool    `json:"requireIndependentProvenance,omitempty"`
	RequireConsumerUtility       bool    `json:"requireConsumerUtility,omitempty"`
	RequireTruthBoundary         bool    `json:"requireTruthBoundary,omitempty"`
	DirectAuthority              bool    `json:"directAuthority,omitempty"`
}

// Evidence intentionally contains only non-secret measurements. Credentials,
// request headers, cookies and tokens have no representation in this contract.
type Evidence struct {
	Observations          int64   `json:"observations"`
	Successes             int64   `json:"successes"`
	Errors                int64   `json:"errors"`
	AuthFailures          int64   `json:"authFailures"`
	RateLimited           int64   `json:"rateLimited"`
	ServerFailures        int64   `json:"serverFailures"`
	P95LatencyMs          int64   `json:"p95LatencyMs,omitempty"`
	FreshnessState        string  `json:"freshnessState,omitempty"`
	SchemaIntegrity       bool    `json:"schemaIntegrity"`
	FallbackVerified      bool    `json:"fallbackVerified"`
	CorroborationSamples  int64   `json:"corroborationSamples"`
	DisagreementPct       float64 `json:"disagreementPct,omitempty"`
	IndependentProvenance bool    `json:"independentProvenance"`
	ConsumerUtilityProven bool    `json:"consumerUtilityProven"`
	TruthBoundaryVerified bool    `json:"truthBoundaryVerified"`
}

type Diagnostic struct {
	Provider       string   `json:"provider"`
	Capability     string   `json:"capability"`
	Dataset        string   `json:"dataset,omitempty"`
	Lifecycle      string   `json:"lifecycle"`
	Readiness      string   `json:"readiness"`
	PromotionMode  string   `json:"promotionMode"`
	AuthorityClass string   `json:"authorityClass"`
	Reasons        []string `json:"reasons,omitempty"`
}

func CanonicalLifecycle(value string) string {
	switch strings.ToUpper(strings.TrimSpace(value)) {
	case Shadow:
		return Shadow
	case Validated:
		return Validated
	case Approved:
		return Approved
	case Production:
		return Production
	default:
		return Shadow
	}
}

func freshnessUsable(state string) bool {
	switch strings.ToUpper(strings.TrimSpace(state)) {
	case "LIVE", "FRESH", "DUE SOON", "DELAYED":
		return true
	default:
		return false
	}
}

func pct(n, d int64) float64 {
	if d <= 0 {
		return 0
	}
	return float64(n) / float64(d) * 100
}

// Evaluate computes evidence sufficiency only. READY_FOR_PROMOTION is advisory
// evidence for an explicit governance action; this function cannot promote.
func Evaluate(policy Policy, ev Evidence) Diagnostic {
	policy.Lifecycle = CanonicalLifecycle(policy.Lifecycle)
	promotion := strings.TrimSpace(policy.PromotionMode)
	if promotion == "" {
		promotion = "EXPLICIT_GOVERNED_ONLY"
	}
	out := Diagnostic{
		Provider: policy.Provider, Capability: policy.Capability, Dataset: policy.Dataset,
		Lifecycle: policy.Lifecycle, PromotionMode: promotion, AuthorityClass: policy.AuthorityClass,
	}
	if policy.DirectAuthority {
		out.Readiness = ReadinessNotApplicable
		out.Reasons = []string{"direct authority is not rank-promoted through the provider lifecycle"}
		return out
	}
	if policy.Lifecycle == Production {
		out.Readiness = ReadinessProductionApproved
		out.Reasons = []string{"capability is explicitly governed for production; runtime suppression/recovery remains router owned"}
		return out
	}

	observed := ev.Observations
	if observed <= 0 {
		observed = ev.Successes + ev.Errors
	}
	if observed < policy.MinObservations {
		out.Readiness = ReadinessInsufficientEvidence
		out.Reasons = append(out.Reasons, fmt.Sprintf("observation depth %d/%d", observed, policy.MinObservations))
	}
	if observed > 0 && pct(ev.Successes, observed) < policy.MinSuccessPct {
		out.Reasons = append(out.Reasons, fmt.Sprintf("success %.2f%% < %.2f%%", pct(ev.Successes, observed), policy.MinSuccessPct))
	}
	if ev.AuthFailures > policy.MaxAuthFailures {
		out.Reasons = append(out.Reasons, fmt.Sprintf("auth failures %d > %d", ev.AuthFailures, policy.MaxAuthFailures))
	}
	if observed > 0 && policy.MaxRateLimitedPct >= 0 && pct(ev.RateLimited, observed) > policy.MaxRateLimitedPct {
		out.Reasons = append(out.Reasons, fmt.Sprintf("rate limited %.2f%% > %.2f%%", pct(ev.RateLimited, observed), policy.MaxRateLimitedPct))
	}
	if policy.MaxP95LatencyMs > 0 && ev.P95LatencyMs > policy.MaxP95LatencyMs {
		out.Reasons = append(out.Reasons, fmt.Sprintf("p95 latency %dms > %dms", ev.P95LatencyMs, policy.MaxP95LatencyMs))
	}
	if policy.RequireFreshEvidence && !freshnessUsable(ev.FreshnessState) {
		out.Reasons = append(out.Reasons, "freshness evidence is not usable")
	}
	if policy.RequireSchemaIntegrity && !ev.SchemaIntegrity {
		out.Reasons = append(out.Reasons, "schema/semantic integrity proof missing")
	}
	if policy.RequireFallbackProof && !ev.FallbackVerified {
		out.Reasons = append(out.Reasons, "fallback correctness proof missing")
	}
	if ev.CorroborationSamples < policy.MinCorroborationSamples {
		out.Reasons = append(out.Reasons, fmt.Sprintf("corroboration depth %d/%d", ev.CorroborationSamples, policy.MinCorroborationSamples))
	}
	if policy.MaxDisagreementPct >= 0 && ev.CorroborationSamples > 0 && ev.DisagreementPct > policy.MaxDisagreementPct {
		out.Reasons = append(out.Reasons, fmt.Sprintf("independent disagreement %.2f%% > %.2f%%", ev.DisagreementPct, policy.MaxDisagreementPct))
	}
	if policy.RequireIndependentProvenance && !ev.IndependentProvenance {
		out.Reasons = append(out.Reasons, "independent provenance proof missing")
	}
	if policy.RequireConsumerUtility && !ev.ConsumerUtilityProven {
		out.Reasons = append(out.Reasons, "consumer utility/outcome proof missing")
	}
	if policy.RequireTruthBoundary && !ev.TruthBoundaryVerified {
		out.Reasons = append(out.Reasons, "truth-boundary proof missing")
	}

	if len(out.Reasons) == 0 {
		out.Readiness = ReadinessReadyForPromotion
		out.Reasons = []string{"all deterministic readiness thresholds are satisfied; explicit governed promotion is still required"}
		return out
	}
	if out.Readiness == "" {
		out.Readiness = ReadinessNotReady
	}
	return out
}

func TradeInsightPolicy(capability, dataset string) Policy {
	return Policy{
		Provider: "TradeInsight", Capability: capability, Dataset: dataset, Lifecycle: Shadow,
		AuthorityClass: "ROUTED_SHADOW_FIRST", PromotionMode: "EXPLICIT_GOVERNED_ONLY",
		MinObservations: 20, MinSuccessPct: 95, MaxP95LatencyMs: 5000, MaxAuthFailures: 0, MaxRateLimitedPct: 10,
		MinCorroborationSamples: 5, MaxDisagreementPct: 5,
		RequireFreshEvidence: true, RequireSchemaIntegrity: true, RequireFallbackProof: true,
		RequireIndependentProvenance: true, RequireConsumerUtility: true, RequireTruthBoundary: true,
	}
}

func DirectAuthorityPolicy(provider, capability, dataset string) Policy {
	return Policy{
		Provider: provider, Capability: capability, Dataset: dataset, Lifecycle: Production,
		AuthorityClass: "DIRECT_AUTHORITY", PromotionMode: "EXPLICIT_GOVERNED_ONLY", DirectAuthority: true,
	}
}
