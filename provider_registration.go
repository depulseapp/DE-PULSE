package main

import (
	"sort"
	"strings"

	"depulse/internal/providerlifecycle"
)

const providerRegistrationContractVersion = "provider-registration-v1.0.0"

type ProviderDatasetContract struct {
	Dataset           string
	Priority          int
	Lifecycle         string
	ContractVersion   string
	SchemaContract    string
	TimestampContract string
	FailureContract   string
	EvidenceRef       string
	ExpectedDelay     string
}

type ProviderRegistration struct {
	Name       string
	QuotaLabel string
	CostClass  string
	Configured func(Settings, Secrets) bool
	Routes     []ProviderDatasetContract
}

func inheritedProductionRoute(dataset string, priority int, expectedDelay string) ProviderDatasetContract {
	return ProviderDatasetContract{
		Dataset: dataset, Priority: priority, Lifecycle: providerlifecycle.Production,
		ContractVersion: providerRegistrationContractVersion,
		SchemaContract: "canonical adapter/schema contract inherited from exact-head production baseline",
		TimestampContract: "canonical timestamp/freshness semantics inherited from exact-head production baseline",
		FailureContract: "Smart Router v2 capability, circuit, rate-limit, fallback and degradation contract",
		EvidenceRef: "#79 provider/data-health production closure and current executable regression baseline",
		ExpectedDelay: expectedDelay,
	}
}

func tradeInsightProductionHistoryRoute(priority int) ProviderDatasetContract {
	return ProviderDatasetContract{
		Dataset: canonicalHistoricalBarsDataset, Priority: priority, Lifecycle: providerlifecycle.Production,
		ContractVersion: providerRegistrationContractVersion,
		SchemaContract: "verified /trading-data/v1/ohlc daily OHLCV; raw/adjusted semantics remain explicit",
		TimestampContract: "completed daily bars normalized by canonical Historical Bars owner",
		FailureContract: "adapter failures remain capability-scoped and fall through Smart Router v2",
		EvidenceRef: "#78 qualification and #84 exact-head production adoption",
		ExpectedDelay: "Completed-bar cadence",
	}
}

func providerRegistrations() []ProviderRegistration {
	has := func(v string) bool { return strings.TrimSpace(v) != "" }
	return []ProviderRegistration{
		{Name: "Alpaca", QuotaLabel: "Entitlement / feed dependent", CostClass: "Broker/data entitlement",
			Configured: func(_ Settings, s Secrets) bool { return has(s.AlpacaKey) && has(s.AlpacaSecret) },
			Routes: []ProviderDatasetContract{
				inheritedProductionRoute("US Live Equities", 1, "Near-live when entitled"),
				inheritedProductionRoute(canonicalUSAssetUniverseDataset, 1, "Provider asset-directory cadence"),
				inheritedProductionRoute(canonicalUSMarketCalendarDataset, 1, "Official market-calendar cadence"),
				inheritedProductionRoute(canonicalUSCorporateActionsDataset, 1, "Corporate-action dissemination cadence"),
				inheritedProductionRoute(canonicalHistoricalBarsDataset, 1, "Completed-bar cadence"),
			}},
		{Name: "Finnhub", QuotaLabel: "API plan / endpoint dependent", CostClass: "Free tier / optional paid upgrade",
			Configured: func(_ Settings, s Secrets) bool { return has(s.Finnhub) },
			Routes: []ProviderDatasetContract{
				inheritedProductionRoute("US Live Equities", 2, "Near-live when entitled"),
				inheritedProductionRoute("News", 1, "Near-live when entitled"),
				inheritedProductionRoute("Earnings", 1, "Provider event/filing cadence"),
				inheritedProductionRoute("Fundamentals", 1, "Provider fundamental-update cadence"),
			}},
		{Name: tradeInsightProviderName, QuotaLabel: "Runtime tier / response headers", CostClass: "Beta / free tier",
			Configured: func(_ Settings, s Secrets) bool { return tradeInsightConfigured(s.TradeInsight) },
			Routes: []ProviderDatasetContract{tradeInsightProductionHistoryRoute(2)}},
		{Name: "Twelve Data", QuotaLabel: "Credit based", CostClass: "Free tier / optional paid upgrade",
			Configured: func(_ Settings, s Secrets) bool { return has(s.TwelveData) },
			Routes: []ProviderDatasetContract{
				inheritedProductionRoute("US Live Equities", 3, "Near-live when entitled"),
				inheritedProductionRoute(canonicalGlobalMarketContextDataset, 1, "Near-live when entitled"),
				inheritedProductionRoute("VIX / Indices", 1, "Near-live when entitled"),
				inheritedProductionRoute(canonicalHistoricalBarsDataset, 3, "Completed-bar cadence"),
			}},
		{Name: "Marketaux", QuotaLabel: "Request quota", CostClass: "Free tier / optional paid upgrade",
			Configured: func(_ Settings, s Secrets) bool { return has(s.Marketaux) },
			Routes: []ProviderDatasetContract{inheritedProductionRoute("News", 2, "Near-live when entitled")}},
		{Name: "FRED", QuotaLabel: "Free API key", CostClass: "Free tier / optional paid upgrade",
			Configured: func(_ Settings, s Secrets) bool { return has(s.FRED) },
			Routes: []ProviderDatasetContract{inheritedProductionRoute("Macro", 1, "Series release cadence")}},
		{Name: "SEC", QuotaLabel: "Fair-access policy", CostClass: "Public / no API fee",
			Configured: func(s Settings, _ Secrets) bool { return has(s.SECEmail) },
			Routes: []ProviderDatasetContract{inheritedProductionRoute("Fundamentals", 2, "Filing dissemination cadence")}},
		{Name: "SEC EDGAR", QuotaLabel: "Fair-access policy", CostClass: "Public / no API fee",
			Configured: func(s Settings, _ Secrets) bool { return has(s.SECEmail) },
			Routes: []ProviderDatasetContract{inheritedProductionRoute("SEC", 1, "Filing dissemination cadence")}},
		{Name: "yfinance", QuotaLabel: "Public recovery · best effort", CostClass: "Public / no API fee",
			Configured: func(_ Settings, _ Secrets) bool { return true },
			Routes: []ProviderDatasetContract{
				inheritedProductionRoute("VIX / Indices", 2, "Recovery-only; may be delayed"),
				inheritedProductionRoute(canonicalHistoricalBarsDataset, 4, "Recovery-only; may be delayed"),
				inheritedProductionRoute("Earnings", 2, "Recovery-only; may be delayed"),
				inheritedProductionRoute("Fundamentals", 3, "Recovery-only; may be delayed"),
			}},
		{Name: "CBOE", QuotaLabel: "Public official/delayed", CostClass: "Public / no API fee",
			Configured: func(_ Settings, _ Secrets) bool { return true },
			Routes: []ProviderDatasetContract{inheritedProductionRoute("VIX / Indices", 3, "Official delayed/close validation")}},
		{Name: "BLS", QuotaLabel: "Official API", CostClass: "Public / free API", Configured: func(_ Settings, _ Secrets) bool { return true }},
		{Name: "EIA", QuotaLabel: "Free API key", CostClass: "Public / free API", Configured: func(_ Settings, s Secrets) bool { return has(s.EIA) }},
	}
}

func providerDatasetContractProductionReady(c ProviderDatasetContract) bool {
	return providerlifecycle.CanonicalLifecycle(c.Lifecycle) == providerlifecycle.Production &&
		strings.TrimSpace(c.Dataset) != "" && c.Priority > 0 &&
		strings.TrimSpace(c.ContractVersion) != "" && strings.TrimSpace(c.SchemaContract) != "" &&
		strings.TrimSpace(c.TimestampContract) != "" && strings.TrimSpace(c.FailureContract) != "" &&
		strings.TrimSpace(c.EvidenceRef) != ""
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
