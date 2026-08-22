package main

// tradeInsightCanonicalOwnerDependency is certification-only metadata for the
// durable G3 canonical-owner map from scope issue #61. Keeping it in test scope
// prevents provider-specific ownership metadata from becoming runtime state.
type tradeInsightCanonicalOwnerDependency struct {
	Concern        string
	CanonicalOwner string
	Dependencies   []string
	Contract       string
}

func tradeInsightCanonicalOwnerDependencies() []tradeInsightCanonicalOwnerDependency {
	return []tradeInsightCanonicalOwnerDependency{
		{
			Concern:        "congressional-trading",
			CanonicalOwner: "Research/Event alternative-evidence boundary",
			Dependencies:   []string{"research_truth.go", "event_intelligence.go", "evidence_alternative.go", "EvidenceRecord"},
			Contract:       "Disclosure lag and provenance remain explicit; Event Intelligence stays fetch-free; Congressional evidence is SHADOW-only until promotion and never deterministic trade truth.",
		},
		{
			Concern:        "sec-form4-enrichment",
			CanonicalOwner: "sec_intelligence.go",
			Dependencies:   []string{"research_truth.go", "evidence_alternative.go", "edgar_enrichment.go", "EvidenceRecord"},
			Contract:       "Direct SEC/EDGAR remains authoritative; TradeInsight can only corroborate or enrich after endpoint/schema admission and source-family evidence is not double-counted.",
		},
		{
			Concern:        "opportunity-radar-movers",
			CanonicalOwner: "market_activity_corporate.go -> opportunity_radar.go / Discovery",
			Dependencies:   []string{"existing market-event candidate flow", "shared source metadata", "shared freshness metadata"},
			Contract:       "Provider movers enter the canonical Discovery/Opportunity Radar path as candidate evidence; no TradeInsight-specific scanner or ranking owner is allowed.",
		},
		{
			Concern:        "symbol-validation-fallback",
			CanonicalOwner: "PersistenceManager / Global Symbol Registry + existing symbol validation",
			Dependencies:   []string{"data.go", "router_extreme_v18_8.go"},
			Contract:       "Provider lookup is fallback-only after canonical/local misses; final validation remains U.S.-equity-only while preserving GLD/SLV/USO as explicit tradable exceptions.",
		},
		{
			Concern:        "provider-telemetry-usefulness",
			CanonicalOwner: "ProviderTelemetry / RuntimeLoad diagnostics + Smart Router capability records",
			Dependencies:   []string{"shared provider health", "loading", "staleness", "usefulness"},
			Contract:       "TradeInsight capability metrics attach to shared telemetry and usefulness owners; no provider-specific telemetry store is allowed.",
		},
		{
			Concern:        "freshness-cache-persistence",
			CanonicalOwner: "existing freshness/cache owners + generic EvidenceRecord / DerivedFeatureRecord persistence",
			Dependencies:   []string{"freshness_adaptive.go", "shared caches", "persistence_repository.go", "persistence_intelligence.go"},
			Contract:       "Reuse canonical freshness, cache and persistence ownership; no TradeInsight-specific shadow tables or duplicate canonical state are allowed.",
		},
		{
			Concern:        "quota-rate-limit-backpressure",
			CanonicalOwner: "provider capability registry + router health/backoff path",
			Dependencies:   []string{"shared provider health", "shared backoff", "Smart Provider Router v2 capability state"},
			Contract:       "HTTP 429/backpressure may temporarily de-rank the affected capability but must never reorder the fixed historical provider route.",
		},
		{
			Concern:        "shadow-lifecycle-promotion",
			CanonicalOwner: "Adaptive Intelligence / Smart Provider Router capability state + evidence logs",
			Dependencies:   []string{"shared capability state", "validation evidence", "review/promotion gate"},
			Contract:       "Lifecycle is shadow_registered -> shadow_running -> proposed_promotion -> approved_live; promotion requires evidence and explicit approval, never automatic promotion.",
		},
	}
}
