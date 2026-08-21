package main

// tradeInsightCapabilityAdmission is the durable v18.9 admission contract for
// TradeInsight capabilities. It is metadata only: it does not create routes,
// fetch data, or bypass Smart Provider Router v2.
type tradeInsightCapabilityAdmission struct {
	ID               string
	Capability       string
	Disposition      string
	Consumer         string
	Authority        string
	EndpointEvidence string
	SchemaVerified   bool
	RuntimeEnabled   bool
	Lifecycle        string
	GateReason       string
}

func (a tradeInsightCapabilityAdmission) runtimeAdmitted() bool {
	if !a.RuntimeEnabled || !a.SchemaVerified {
		return false
	}
	switch a.Lifecycle {
	case "SHADOW", "VALIDATED", "APPROVED", "PRODUCTION":
		return true
	default:
		return false
	}
}

func tradeInsightCapabilityAdmissionLookup(id string) (tradeInsightCapabilityAdmission, bool) {
	for _, row := range tradeInsightCapabilityAdmissionRegistry() {
		if row.ID == id {
			return row, true
		}
	}
	return tradeInsightCapabilityAdmission{}, false
}

// tradeInsightCapabilityLifecycleTruth is the shared observability view of
// admission state. A capability that is not fully runtime-admitted is GATED,
// regardless of provider connectivity or whether an endpoint is advertised.
func tradeInsightCapabilityLifecycleTruth(id string) string {
	row, ok := tradeInsightCapabilityAdmissionLookup(id)
	if !ok || !row.runtimeAdmitted() {
		return "GATED"
	}
	return row.Lifecycle
}

// tradeInsightCanonicalOwnerDependency captures the G3 canonical-owner map
// from durable scope issue #61. It is intentionally provider-neutral at the
// ownership boundary so TradeInsight cannot create parallel state, scanners,
// telemetry, persistence, SEC truth, or routing subsystems.
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

// tradeInsightCapabilityAdmissionRegistry mirrors issue #61's admitted
// capability matrix. Capabilities without verified response schemas stay
// non-executable even when the vendor publicly advertises them.
func tradeInsightCapabilityAdmissionRegistry() []tradeInsightCapabilityAdmission {
	return []tradeInsightCapabilityAdmission{
		{
			ID: "daily-history", Capability: "Daily historical OHLCV", Disposition: "FALLBACK + STORE_FOR_HISTORY",
			Consumer: "Historical Bars; history backfill; outcome studies", Authority: "Smart Provider Router v2 is the sole routing authority; daily-only semantics and raw/adjusted provenance are retained.",
			EndpointEvidence: "/trading-data/v1/ohlc", SchemaVerified: true, RuntimeEnabled: true, Lifecycle: "SHADOW",
		},
		{
			ID: "adjusted-history", Capability: "Adjusted daily OHLCV", Disposition: "USE",
			Consumer: "Historical analytics; outcome evaluation", Authority: "Uses the canonical Historical Bars owner; adjusted and raw series must never be mixed silently.",
			EndpointEvidence: "/trading-data/v1/ohlc adjusted semantics", SchemaVerified: true, RuntimeEnabled: true, Lifecycle: "SHADOW",
		},
		{
			ID: "corporate-actions", Capability: "Dividends and stock splits", Disposition: "USE + CORROBORATE",
			Consumer: "Canonical corporate-action ledger; historical normalization; Research", Authority: "Supplemental evidence only; canonical corporate-action truth and provenance rules remain authoritative.",
			EndpointEvidence: "/trading-data/v1/ohlc corporate-action fields", SchemaVerified: true, RuntimeEnabled: true, Lifecycle: "SHADOW",
		},
		{
			ID: "bulk-history", Capability: "Bounded multi-ticker history", Disposition: "USE SELECTIVELY",
			Consumer: "Bounded history backfill and outcome-analysis jobs", Authority: "Reuse canonical history ownership; cache-first, budgeted and backpressure-aware; never continuous broad refetch.",
			EndpointEvidence: "Official tidata SDK performs client-side bounded fan-out over per-ticker daily history; no server-side bulk endpoint is assumed.", SchemaVerified: true, RuntimeEnabled: true, Lifecycle: "SHADOW",
			GateReason: "Canonical Historical Bars route now provides deduplicated, sequential client-side fan-out capped at 50 symbols over the verified per-ticker daily endpoint; SHADOW-only pending validation evidence and explicit promotion approval.",
		},
		{
			ID: "congressional-trades", Capability: "Congressional trades/disclosures", Disposition: "USE",
			Consumer: "Congressional Trading Intelligence; Research; catalyst/context correlation", Authority: "Alternative evidence only; disclosure lag remains explicit and it never becomes deterministic trade truth.",
			EndpointEvidence: "Official TradeInsight data-api/docs/insight-data.mdx: GET /trading-data/v1/congress/v1/trades with ticker/limit/offset, data/total/limit/offset envelope, and typed disclosure fields.", SchemaVerified: true, RuntimeEnabled: true, Lifecycle: "SHADOW",
			GateReason: "Official upstream REST schema verified; admitted SHADOW-only pending validation evidence and explicit promotion approval.",
		},
		{
			ID: "sec-form4", Capability: "SEC Form 4 normalized insider trades", Disposition: "CORROBORATE + ENRICH",
			Consumer: "SEC/Ownership Research; insider normalization; catalyst correlation", Authority: "Direct SEC/EDGAR remains primary authority; source-family evidence must not be double-counted.",
			EndpointEvidence: "", SchemaVerified: false, RuntimeEnabled: false, Lifecycle: "NOT_IMPLEMENTED",
			GateReason: "Exact REST endpoint and configured-key output schema are not yet verified.",
		},
		{
			ID: "top-movers", Capability: "Top movers", Disposition: "SHADOW + CORROBORATE",
			Consumer: "Opportunity Radar candidate discovery and provider-usefulness evaluation", Authority: "Candidate evidence only; existing DE.PULSE scanner and ranking remain canonical.",
			EndpointEvidence: "Official MCP get_top_movers capability only; production REST path/output schema not verified.", SchemaVerified: false, RuntimeEnabled: false, Lifecycle: "NOT_IMPLEMENTED",
			GateReason: "Do not infer a REST endpoint or parser from the MCP tool name.",
		},
		{
			ID: "symbol-search", Capability: "Ticker/company search", Disposition: "FALLBACK + CORROBORATE",
			Consumer: "Canonical symbol resolution and Add Symbol resilience", Authority: "Existing U.S.-equity boundary and canonical symbol validation remain final.",
			EndpointEvidence: "Official MCP search_ticker capability only; production REST path/output schema not verified.", SchemaVerified: false, RuntimeEnabled: false, Lifecycle: "NOT_IMPLEMENTED",
			GateReason: "Do not infer a REST endpoint or parser from the MCP tool name.",
		},
		{
			ID: "generic-market-price", Capability: "Generic market-price capability beyond documented daily OHLCV", Disposition: "FUTURE",
			Consumer: "Potential resilience evidence only if entitlement proves a useful non-duplicate surface", Authority: "No live-feed path and no intraday assumption.",
			EndpointEvidence: "", SchemaVerified: false, RuntimeEnabled: false, Lifecycle: "NOT_IMPLEMENTED",
			GateReason: "No independently verified and entitled production contract beyond daily OHLCV.",
		},
		{
			ID: "mcp-interface", Capability: "TradeInsight MCP server", Disposition: "FUTURE",
			Consumer: "Optional developer/AI research interface", Authority: "Never the production provider-routing owner.",
			EndpointEvidence: "Official MCP capability documentation", SchemaVerified: true, RuntimeEnabled: false, Lifecycle: "NOT_IMPLEMENTED",
			GateReason: "Explicitly excluded from the production data path in v18.9.",
		},
		{
			ID: "python-sdk", Capability: "tidata Python SDK", Disposition: "REFERENCE",
			Consumer: "Request, adjustment, error and bounded fan-out semantics for tests/design", Authority: "Production Go runtime remains native; no Python runtime dependency is introduced.",
			EndpointEvidence: "Official TradeInsight-Info/tidata SDK", SchemaVerified: true, RuntimeEnabled: false, Lifecycle: "NOT_IMPLEMENTED",
			GateReason: "Reference semantics only; not an executable production dependency.",
		},
		{
			ID: "vendor-derived-scores", Capability: "Vendor-derived scores or recommendations", Disposition: "FUTURE",
			Consumer: "Experimental corroboration only if a future entitlement exposes them", Authority: "Never deterministic Day/Swing/Long truth without independent validation and promotion evidence.",
			EndpointEvidence: "", SchemaVerified: false, RuntimeEnabled: false, Lifecycle: "NOT_IMPLEMENTED",
			GateReason: "No verified current capability contract.",
		},
	}
}
