package main

import (
	"fmt"
	"math"
	"strings"
)

const (
	providerLifecycleShadow     = "SHADOW"
	providerLifecycleValidated  = "VALIDATED"
	providerLifecycleApproved   = "APPROVED"
	providerLifecycleProduction = "PRODUCTION"
)

const (
	providerReadinessNotApplicable       = "N/A"
	providerReadinessInsufficientEvidence = "INSUFFICIENT_EVIDENCE"
	providerReadinessNotReady            = "NOT_READY"
	providerReadinessReadyForPromotion   = "READY_FOR_PROMOTION"
	providerReadinessProductionApproved  = "PRODUCTION_APPROVED"
)

// ProviderCapabilityLifecyclePolicy is governance metadata only. Runtime health,
// circuit breaking, freshness, fallback and admission remain owned by Smart
// Provider Router v2 and the canonical Data Health owners. In particular,
// EvaluateProviderCapabilityReadiness never mutates Lifecycle.
type ProviderCapabilityLifecyclePolicy struct {
	Provider                     string `json:"provider"`
	Capability                   string `json:"capability"`
	Dataset                      string `json:"dataset,omitempty"`
	Lifecycle                    string `json:"lifecycle"`
	AuthorityClass               string `json:"authorityClass"`
	PromotionMode                string `json:"promotionMode"`
	MinObservations              int64  `json:"minObservations,omitempty"`
	MinSuccessPct                float64 `json:"minSuccessPct,omitempty"`
	MaxP95LatencyMs              int64  `json:"maxP95LatencyMs,omitempty"`
	MaxAuthFailures              int64  `json:"maxAuthFailures,omitempty"`
	MaxRateLimitedPct            float64 `json:"maxRateLimitedPct,omitempty"`
	MinCorroborationSamples      int64  `json:"minCorroborationSamples,omitempty"`
	MaxDisagreementPct           float64 `json:"maxDisagreementPct,omitempty"`
	RequireFreshEvidence         bool   `json:"requireFreshEvidence,omitempty"`
	RequireSchemaIntegrity       bool   `json:"requireSchemaIntegrity,omitempty"`
	RequireFallbackProof         bool   `json:"requireFallbackProof,omitempty"`
	RequireIndependentProvenance bool   `json:"requireIndependentProvenance,omitempty"`
	RequireConsumerUtility       bool   `json:"requireConsumerUtility,omitempty"`
	RequireTruthBoundary         bool   `json:"requireTruthBoundary,omitempty"`
	DirectAuthority              bool   `json:"directAuthority,omitempty"`
}

// ProviderCapabilityLifecycleEvidence contains only non-secret measurements.
// Credential values, request headers and tokens have no representation here.
type ProviderCapabilityLifecycleEvidence struct {
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

type ProviderCapabilityReadinessDiagnostic struct {
	Provider       string   `json:"provider"`
	Capability     string   `json:"capability"`
	Dataset        string   `json:"dataset,omitempty"`
	Lifecycle      string   `json:"lifecycle"`
	Readiness      string   `json:"readiness"`
	PromotionMode  string   `json:"promotionMode"`
	AuthorityClass string   `json:"authorityClass"`
	Reasons        []string `json:"reasons,omitempty"`
}

func canonicalProviderLifecycle(value string) string {
	switch strings.ToUpper(strings.TrimSpace(value)) {
	case providerLifecycleShadow:
		return providerLifecycleShadow
	case providerLifecycleValidated:
		return providerLifecycleValidated
	case providerLifecycleApproved:
		return providerLifecycleApproved
	case providerLifecycleProduction:
		return providerLifecycleProduction
	default:
		return providerLifecycleShadow
	}
}

func lifecycleFreshnessUsable(state string) bool {
	switch strings.ToUpper(strings.TrimSpace(state)) {
	case "LIVE", "FRESH", "DUE SOON", "DELAYED":
		return true
	default:
		return false
	}
}

func readinessPercent(n, d int64) float64 {
	if d <= 0 {
		return 0
	}
	return float64(n) / float64(d) * 100
}

// EvaluateProviderCapabilityReadiness computes evidence sufficiency only. A
// READY_FOR_PROMOTION result is advisory governance evidence and cannot alter
// SHADOW/VALIDATED/APPROVED/PRODUCTION lifecycle state automatically.
func EvaluateProviderCapabilityReadiness(policy ProviderCapabilityLifecyclePolicy, ev ProviderCapabilityLifecycleEvidence) ProviderCapabilityReadinessDiagnostic {
	policy.Lifecycle = canonicalProviderLifecycle(policy.Lifecycle)
	out := ProviderCapabilityReadinessDiagnostic{
		Provider: policy.Provider, Capability: policy.Capability, Dataset: policy.Dataset,
		Lifecycle: policy.Lifecycle, PromotionMode: defaultString(policy.PromotionMode, "EXPLICIT_GOVERNED_ONLY"),
		AuthorityClass: policy.AuthorityClass,
	}
	if policy.DirectAuthority {
		out.Readiness = providerReadinessNotApplicable
		out.Reasons = []string{"direct authority is not rank-promoted through the provider lifecycle"}
		return out
	}
	if policy.Lifecycle == providerLifecycleProduction {
		out.Readiness = providerReadinessProductionApproved
		out.Reasons = []string{"capability is already explicitly governed for production; runtime suppression/recovery remains Router v2 owned"}
		return out
	}

	observed := ev.Observations
	if observed <= 0 {
		observed = ev.Successes + ev.Errors
	}
	if observed < policy.MinObservations {
		out.Readiness = providerReadinessInsufficientEvidence
		out.Reasons = append(out.Reasons, fmt.Sprintf("observation depth %d/%d", observed, policy.MinObservations))
	}
	if observed > 0 && readinessPercent(ev.Successes, observed) < policy.MinSuccessPct {
		out.Reasons = append(out.Reasons, fmt.Sprintf("success %.2f%% < %.2f%%", readinessPercent(ev.Successes, observed), policy.MinSuccessPct))
	}
	if ev.AuthFailures > policy.MaxAuthFailures {
		out.Reasons = append(out.Reasons, fmt.Sprintf("auth failures %d > %d", ev.AuthFailures, policy.MaxAuthFailures))
	}
	if observed > 0 && policy.MaxRateLimitedPct >= 0 && readinessPercent(ev.RateLimited, observed) > policy.MaxRateLimitedPct {
		out.Reasons = append(out.Reasons, fmt.Sprintf("rate limited %.2f%% > %.2f%%", readinessPercent(ev.RateLimited, observed), policy.MaxRateLimitedPct))
	}
	if policy.MaxP95LatencyMs > 0 && ev.P95LatencyMs > policy.MaxP95LatencyMs {
		out.Reasons = append(out.Reasons, fmt.Sprintf("p95 latency %dms > %dms", ev.P95LatencyMs, policy.MaxP95LatencyMs))
	}
	if policy.RequireFreshEvidence && !lifecycleFreshnessUsable(ev.FreshnessState) {
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
		out.Readiness = providerReadinessReadyForPromotion
		out.Reasons = []string{"all deterministic readiness thresholds are satisfied; explicit governed promotion is still required"}
		return out
	}
	if out.Readiness == "" {
		out.Readiness = providerReadinessNotReady
	}
	return out
}

func defaultProviderLifecyclePolicy(provider, capability, dataset string) ProviderCapabilityLifecyclePolicy {
	return ProviderCapabilityLifecyclePolicy{
		Provider: provider, Capability: capability, Dataset: dataset, Lifecycle: providerLifecycleProduction,
		AuthorityClass: "ROUTED_PROVIDER", PromotionMode: "EXPLICIT_GOVERNED_ONLY",
	}
}

func tradeInsightProductionReadinessPolicy(capability, dataset string) ProviderCapabilityLifecyclePolicy {
	return ProviderCapabilityLifecyclePolicy{
		Provider: tradeInsightProviderName, Capability: capability, Dataset: dataset, Lifecycle: providerLifecycleShadow,
		AuthorityClass: "ROUTED_SHADOW_FIRST", PromotionMode: "EXPLICIT_GOVERNED_ONLY",
		MinObservations: 20, MinSuccessPct: 95, MaxP95LatencyMs: 5000, MaxAuthFailures: 0, MaxRateLimitedPct: 10,
		MinCorroborationSamples: 5, MaxDisagreementPct: 5,
		RequireFreshEvidence: true, RequireSchemaIntegrity: true, RequireFallbackProof: true,
		RequireIndependentProvenance: true, RequireConsumerUtility: true, RequireTruthBoundary: true,
	}
}

func providerLifecyclePolicy(provider, capability, dataset string) ProviderCapabilityLifecyclePolicy {
	if strings.EqualFold(provider, tradeInsightProviderName) {
		return tradeInsightProductionReadinessPolicy(capability, dataset)
	}
	if strings.EqualFold(provider, "SEC EDGAR") || strings.EqualFold(provider, "SEC") ||
		strings.EqualFold(provider, "CBOE") || strings.EqualFold(provider, "BLS") ||
		strings.EqualFold(provider, "EIA") || strings.EqualFold(provider, "U.S. Treasury") ||
		strings.EqualFold(provider, "BEA") || strings.EqualFold(provider, "TWSE") ||
		strings.EqualFold(provider, "Official Macro Calendars") {
		p := defaultProviderLifecyclePolicy(provider, capability, dataset)
		p.AuthorityClass = "DIRECT_AUTHORITY"
		p.DirectAuthority = true
		return p
	}
	return defaultProviderLifecyclePolicy(provider, capability, dataset)
}

// lifecycleEvidenceFromRuntime reuses existing Router v2 capability state and
// ProviderTelemetry instead of creating a provider-specific telemetry engine.
// Evidence that cannot be justified from these owners remains false/unknown and
// therefore cannot make a SHADOW capability READY_FOR_PROMOTION by accident.
func lifecycleEvidenceFromRuntime(cap ProviderCapabilityStateRecord, telemetry ProviderRequestDiagnostics, freshnessState string) ProviderCapabilityLifecycleEvidence {
	observations := cap.SuccessCount + cap.FailureCount
	if telemetry.Requests > observations {
		observations = telemetry.Requests
	}
	p95 := telemetry.P95LatencyMs
	if p95 <= 0 {
		p95 = telemetry.AverageLatencyMs
	}
	return ProviderCapabilityLifecycleEvidence{
		Observations: observations,
		Successes: int64(math.Max(float64(cap.SuccessCount), float64(telemetry.Successes))),
		Errors: int64(math.Max(float64(cap.FailureCount), float64(telemetry.Errors))),
		RateLimited: telemetry.RateLimited,
		P95LatencyMs: p95,
		FreshnessState: freshnessState,
	}
}
