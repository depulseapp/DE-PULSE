package main

import (
	"sort"
	"strings"

	"depulse/internal/providerlifecycle"
)

const providerRegistrationContractVersion = "provider-registration-v1.1.0"

// ProviderDatasetContract is the bounded adaptive work/evidence contract for one
// provider capability on one canonical dataset. It is deliberately descriptive:
// canonical data ownership, lifecycle/readiness, routing, health, freshness,
// persistence and degradation remain owned by the existing DE.PULSE systems.
type ProviderDatasetContract struct {
	Dataset           string
	Capability        string
	Priority          int
	Lifecycle         string
	ContractVersion   string
	CanonicalOwner    string
	Consumer          string
	AdapterContract   string
	SchemaContract    string
	TimestampContract string
	FreshnessContract string
	FailureContract   string
	RightsContract    string
	EvidenceRef       string
	ApprovalRef       string
	InvalidationRule  string
	ExpectedDelay     string
	Uses              []string
}

// ProviderRegistration is the single provider-onboarding descriptor consumed
// by the existing router/capability owners. Adding a provider still requires
// code for its adapter/normalizer and governed evidence, but provider metadata
// and route adoption no longer need parallel switches/maps.
type ProviderRegistration struct {
	Name       string
	QuotaLabel string
	CostClass  string
	Configured func(Settings, Secrets) bool
	Routes     []ProviderDatasetContract
}

func inheritedProductionRoute(provider, dataset, capability string, priority int, expectedDelay string, uses ...string) ProviderDatasetContract {
	return ProviderDatasetContract{
		Dataset: dataset, Capability: capability, Priority: priority, Lifecycle: providerlifecycle.Production,
		ContractVersion: providerRegistrationContractVersion,
		CanonicalOwner:  dataset,
		Consumer:        strings.Join(uses, ", "),
		AdapterContract: "existing provider adapter/normalizer on the exact-head production baseline",
		SchemaContract:  "canonical adapter/schema contract inherited from exact-head production baseline",
		TimestampContract: "canonical provider timestamp semantics inherited from exact-head production baseline",
		FreshnessContract: "canonical Data Freshness/SLO owner evaluates age/cadence; provider registration cannot redefine freshness",
		FailureContract: "Smart Router v2 capability classification, circuit, rate-limit, fallback and truthful degradation contract",
		RightsContract: "provider-data-rights governance remains separate from operational entitlement and fails closed for unbound commercial/redistribution/AI rights",
		EvidenceRef: "#79 provider/data-health production closure plus exact-head regression evidence for " + provider + " / " + dataset,
		ApprovalRef: "existing production route inherited from the #79/#84 governed provider program; any semantic change requires fresh exact-head qualification",
		InvalidationRule: "adapter/schema/timestamp/freshness/failure/rights/priority/config semantics or relevant source fingerprint change invalidates inherited evidence",
		ExpectedDelay: expectedDelay,
		Uses: append([]string(nil), uses...),
	}
}

func tradeInsightProductionHistoryRoute(priority int) ProviderDatasetContract {
	return ProviderDatasetContract{
		Dataset: canonicalHistoricalBarsDataset, Capability: "Adjusted daily OHLCV / corporate-action corroboration", Priority: priority,
		Lifecycle: providerlifecycle.Production, ContractVersion: providerRegistrationContractVersion,
		CanonicalOwner: canonicalHistoricalBarsDataset,
		Consumer: "Historical Bars, Research, Corporate Actions, Adaptive Intelligence",
		AdapterContract: "native Go TradeInsight daily-history adapter through the canonical Historical Bars loader",
		SchemaContract: "verified /trading-data/v1/ohlc daily OHLCV contract; raw/adjusted and corporate-action semantics remain explicit",
		TimestampContract: "completed daily bars normalized by canonical Historical Bars owner",
		FreshnessContract: "completed-bar cadence evaluated by canonical Historical Bars freshness; never presented as live/intraday truth",
		FailureContract: "TradeInsight adapter failures remain capability-scoped and fall through Smart Router v2 without poisoning unrelated capabilities",
		RightsContract: "operational entitlement is separate from provider data-rights/commercial readiness; unbound rights fail closed",
		EvidenceRef: "#78 qualification evidence and #84 production adoption evidence",
		ApprovalRef: "#84 exact-head governed production adoption",
		InvalidationRule: "endpoint/schema/adjustment/timestamp/rights/authority or adapter behavior change requires new qualification and approval evidence",
		ExpectedDelay: "Completed-bar cadence",
		Uses: []string{"Historical Bars", "Research", "Corporate Actions", "Adaptive Intelligence"},
	}
}

func providerRegistrations() []ProviderRegistration {
	has := func(v string) bool { return strings.TrimSpace(v) != "" }
	return []ProviderRegistration{
		{Name: "Alpaca", QuotaLabel: "Entitlement / feed dependent", CostClass: "Broker/data entitlement",
			Configured: func(_ Settings, s Secrets) bool { return has(s.AlpacaKey) && has(s.AlpacaSecret) },
			Routes: []ProviderDatasetContract{
				inheritedProductionRoute("Alpaca", "US Live Equities", "IEX quotes / snapshots / liquidity", 1, "Near-live when entitled", "Day", "Market Open Prep", "Trade Readiness", "Decision Queue"),
				inheritedProductionRoute("Alpaca", canonicalUSAssetUniverseDataset, "U.S. asset universe / instrument identity", 1, "Provider asset-directory cadence", "Symbol Validation", "Instrument Identity"),
				inheritedProductionRoute("Alpaca", canonicalUSMarketCalendarDataset, "U.S. market calendar", 1, "Official market-calendar cadence", "Session Awareness", "Preparation Jobs"),
				inheritedProductionRoute("Alpaca", canonicalUSCorporateActionsDataset, "Corporate actions", 1, "Corporate-action dissemination cadence", "Corporate Actions", "Historical Normalization", "Research"),
				inheritedProductionRoute("Alpaca", canonicalHistoricalBarsDataset, "Historical OHLCV", 1, "Completed-bar cadence", "Historical Bars", "Research", "Adaptive Intelligence"),
			}},
		{Name: "Finnhub", QuotaLabel: "API plan / endpoint dependent", CostClass: "Free tier / optional paid upgrade",
			Configured: func(_ Settings, s Secrets) bool { return has(s.Finnhub) },
			Routes: []ProviderDatasetContract{
				inheritedProductionRoute("Finnhub", "US Live Equities", "Primary U.S. equity", 2, "Near-live when entitled", "Day", "Swing", "Long", "Discovery", "Trade Readiness"),
				inheritedProductionRoute("Finnhub", "News", "Company news", 1, "Near-live when entitled", "News", "Research", "Catalyst Watch"),
				inheritedProductionRoute("Finnhub", "Earnings", "Earnings calendar / surprise", 1, "Provider event/filing cadence", "Earnings", "Catalyst Watch", "Research"),
				inheritedProductionRoute("Finnhub", "Fundamentals", "Company fundamentals", 1, "Provider fundamental-update cadence", "Research", "Swing", "Long"),
			}},
		{Name: tradeInsightProviderName, QuotaLabel: "Runtime tier / response headers", CostClass: "Beta / free tier",
			Configured: func(_ Settings, s Secrets) bool { return tradeInsightConfigured(s.TradeInsight) },
			Routes: []ProviderDatasetContract{tradeInsightProductionHistoryRoute(2)}},
		{Name: "Twelve Data", QuotaLabel: "Credit based", CostClass: "Free tier / optional paid upgrade",
			Configured: func(_ Settings, s Secrets) bool { return has(s.TwelveData) },
			Routes: []ProviderDatasetContract{
				inheritedProductionRoute("Twelve Data", "US Live Equities", "U.S. equity fallback", 3, "Near-live when entitled", "Data Freshness", "Recovery"),
				inheritedProductionRoute("Twelve Data", canonicalGlobalMarketContextDataset, "FX / direct global context", 1, "Near-live when entitled", "Market Regime", "Dashboard", "Research"),
				inheritedProductionRoute("Twelve Data", "VIX / Indices", "VIX / indices", 1, "Near-live when entitled", "VIX", "Market Regime", "Day", "Swing", "Long"),
				inheritedProductionRoute("Twelve Data", canonicalHistoricalBarsDataset, "Historical OHLCV fallback", 3, "Completed-bar cadence", "Historical Bars", "Recovery"),
			}},
		{Name: "Marketaux", QuotaLabel: "Request quota", CostClass: "Free tier / optional paid upgrade",
			Configured: func(_ Settings, s Secrets) bool { return has(s.Marketaux) },
			Routes: []ProviderDatasetContract{inheritedProductionRoute("Marketaux", "News", "Stock news fallback", 2, "Near-live when entitled", "News", "Dashboard", "Research", "Catalyst Watch")}},
		{Name: "FRED", QuotaLabel: "Free API key", CostClass: "Free tier / optional paid upgrade",
			Configured: func(_ Settings, s Secrets) bool { return has(s.FRED) },
			Routes: []ProviderDatasetContract{inheritedProductionRoute("FRED", "Macro", "Rates / credit / conditions / USD", 1, "Series release cadence", "Market Regime", "Swing", "Long", "Research", "Trade Readiness")}},
		{Name: "SEC", QuotaLabel: "Fair-access policy", CostClass: "Public / no API fee",
			Configured: func(s Settings, _ Secrets) bool { return has(s.SECEmail) },
			Routes: []ProviderDatasetContract{inheritedProductionRoute("SEC", "Fundamentals", "SEC fundamental authority/fallback", 2, "Filing dissemination cadence", "Research", "Fundamentals")}},
		{Name: "SEC EDGAR", QuotaLabel: "Fair-access policy", CostClass: "Public / no API fee",
			Configured: func(s Settings, _ Secrets) bool { return has(s.SECEmail) },
			Routes: []ProviderDatasetContract{inheritedProductionRoute("SEC EDGAR", "SEC", "Direct SEC/EDGAR filings and Form 4 authority", 1, "Filing dissemination cadence", "SEC", "Research", "Catalyst Watch")}},
		{Name: "yfinance", QuotaLabel: "Public recovery · best effort", CostClass: "Public / no API fee",
			Configured: func(_ Settings, _ Secrets) bool { return true },
			Routes: []ProviderDatasetContract{
				inheritedProductionRoute("yfinance", "VIX / Indices", "Recovery-only VIX/index context", 2, "Recovery-only; may be delayed", "Data Freshness", "Recovery"),
				inheritedProductionRoute("yfinance", canonicalHistoricalBarsDataset, "Recovery-only historical bars", 4, "Recovery-only; may be delayed", "Historical Bars", "Recovery"),
				inheritedProductionRoute("yfinance", "Earnings", "Recovery-only earnings", 2, "Recovery-only; may be delayed", "Earnings", "Recovery"),
				inheritedProductionRoute("yfinance", "Fundamentals", "Recovery-only fundamentals", 3, "Recovery-only; may be delayed", "Fundamentals", "Recovery"),
			}},
		{Name: "CBOE", QuotaLabel: "Public official/delayed", CostClass: "Public / no API fee",
			Configured: func(_ Settings, _ Secrets) bool { return true },
			Routes: []ProviderDatasetContract{inheritedProductionRoute("CBOE", "VIX / Indices", "Official VIX validation / delayed close", 3, "Official delayed/close validation", "VIX", "Market Regime", "Data Freshness")}},
		{Name: "BLS", QuotaLabel: "Official API", CostClass: "Public / free API", Configured: func(_ Settings, _ Secrets) bool { return true }},
		{Name: "EIA", QuotaLabel: "Free API key", CostClass: "Public / free API", Configured: func(_ Settings, s Secrets) bool { return has(s.EIA) }},
	}
}

func providerDatasetContractValidationErrors(c ProviderDatasetContract) []string {
	var out []string
	require := func(name, value string) {
		if strings.TrimSpace(value) == "" { out = append(out, name+" missing") }
	}
	require("dataset", c.Dataset)
	require("capability", c.Capability)
	require("contractVersion", c.ContractVersion)
	require("canonicalOwner", c.CanonicalOwner)
	require("consumer", c.Consumer)
	require("adapterContract", c.AdapterContract)
	require("schemaContract", c.SchemaContract)
	require("timestampContract", c.TimestampContract)
	require("freshnessContract", c.FreshnessContract)
	require("failureContract", c.FailureContract)
	require("rightsContract", c.RightsContract)
	require("evidenceRef", c.EvidenceRef)
	require("approvalRef", c.ApprovalRef)
	require("invalidationRule", c.InvalidationRule)
	if c.Priority <= 0 { out = append(out, "priority missing") }
	if len(c.Uses) == 0 { out = append(out, "downstream uses missing") }
	return out
}

func providerDatasetContractProductionReady(c ProviderDatasetContract) bool {
	return providerlifecycle.CanonicalLifecycle(c.Lifecycle) == providerlifecycle.Production && len(providerDatasetContractValidationErrors(c)) == 0
}

func routeChainsFromProviderRegistrations(regs []ProviderRegistration) map[string][]string {
	type ranked struct { name string; priority int }
	byDataset := map[string][]ranked{}
	for _, reg := range regs {
		for _, route := range reg.Routes {
			if providerDatasetContractProductionReady(route) {
				byDataset[route.Dataset] = append(byDataset[route.Dataset], ranked{reg.Name, route.Priority})
			}
		}
	}
	out := make(map[string][]string, len(byDataset))
	for dataset, rows := range byDataset {
		sort.SliceStable(rows, func(i, j int) bool {
			if rows[i].priority != rows[j].priority { return rows[i].priority < rows[j].priority }
			return rows[i].name < rows[j].name
		})
		for _, row := range rows { out[dataset] = append(out[dataset], row.name) }
	}
	return out
}

func providerRegistrationLookup(provider string) (ProviderRegistration, bool) {
	want := providerKey(provider)
	for _, reg := range providerRegistrations() {
		if providerKey(reg.Name) == want { return reg, true }
	}
	return ProviderRegistration{}, false
}

func providerConfiguredFromRegistration(provider string, settings Settings, secrets Secrets) bool {
	reg, ok := providerRegistrationLookup(provider)
	return ok && reg.Configured != nil && reg.Configured(settings, secrets)
}

func providerQuotaFromRegistration(provider string) string {
	if reg, ok := providerRegistrationLookup(provider); ok && strings.TrimSpace(reg.QuotaLabel) != "" { return reg.QuotaLabel }
	return "Provider dependent"
}

func providerCostFromRegistration(provider string) string {
	if reg, ok := providerRegistrationLookup(provider); ok && strings.TrimSpace(reg.CostClass) != "" { return reg.CostClass }
	return "Provider dependent"
}

func providerExpectedDelayFromRegistration(dataset, provider string) string {
	if reg, ok := providerRegistrationLookup(provider); ok {
		for _, route := range reg.Routes {
			if route.Dataset == dataset && strings.TrimSpace(route.ExpectedDelay) != "" { return route.ExpectedDelay }
		}
	}
	return "Near-live when entitled"
}
