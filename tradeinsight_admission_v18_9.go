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
			EndpointEvidence: "Official tidata SDK performs client-side bounded fan-out over per-ticker daily history; no server-side bulk endpoint is assumed.", SchemaVerified: true, RuntimeEnabled: false, Lifecycle: "NOT_IMPLEMENTED",
			GateReason: "Bounded job/cache/backpressure consumer is not yet wired through the canonical history owner.",
		},
		{
			ID: "congressional-trades", Capability: "Congressional trades/disclosures", Disposition: "USE",
			Consumer: "Congressional Trading Intelligence; Research; catalyst/context correlation", Authority: "Alternative evidence only; disclosure lag remains explicit and it never becomes deterministic trade truth.",
			EndpointEvidence: "/trading-data/v1/congress/v1/trades", SchemaVerified: false, RuntimeEnabled: false, Lifecycle: "NOT_IMPLEMENTED",
			GateReason: "Public endpoint is evidenced, but the configured-key response schema has not been verified.",
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
