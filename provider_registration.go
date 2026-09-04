package main

import (
	"crypto/sha256"
	"sort"
	"strings"

	"depulse/internal/providerlifecycle"
)

const providerRegistrationContractVersion = "provider-registration-v1.3.0"

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

type ProviderCapabilityDiagnosticRegistration struct {
	Capability string
	Detail     string
	Uses       []string
	Status     func(Settings, Secrets, map[string]string, map[string]SymbolIntelligence, map[string]GlobalDriver) string
}

// ProviderRegistration is the single provider-onboarding descriptor consumed
// by the existing router/capability owners. Adding a provider still requires
// code for its adapter/normalizer and governed evidence, but provider metadata,
// configuration invalidation, capability diagnostics and route adoption no
// longer need parallel provider maps.
type ProviderRegistration struct {
	Name                     string
	QuotaLabel               string
	CostClass                string
	Configured               func(Settings, Secrets) bool
	ConfigurationFingerprint func(Settings, Secrets) [32]byte
	Routes                   []ProviderDatasetContract
	Diagnostics              []ProviderCapabilityDiagnosticRegistration
}

func providerConfigurationFingerprint(parts ...string) [32]byte {
	return sha256.Sum256([]byte(strings.Join(parts, "\x00")))
}

func inheritedProductionRoute(provider, dataset, capability string, priority int, expectedDelay string, uses ...string) ProviderDatasetContract {
	return ProviderDatasetContract{
		Dataset: dataset, Capability: capability, Priority: priority, Lifecycle: providerlifecycle.Production,
		ContractVersion:   providerRegistrationContractVersion,
		CanonicalOwner:    dataset,
		Consumer:          strings.Join(uses, ", "),
		AdapterContract:   "existing provider adapter/normalizer on the exact-head production baseline",
		SchemaContract:    "canonical adapter/schema contract inherited from exact-head production baseline",
		TimestampContract: "canonical provider timestamp semantics inherited from exact-head production baseline",
		FreshnessContract: "canonical Data Freshness/SLO owner evaluates age/cadence; provider registration cannot redefine freshness",
		FailureContract:   "Smart Router v2 capability classification, circuit, rate-limit, fallback and truthful degradation contract",
		RightsContract:    "provider-data-rights governance remains separate from operational entitlement and fails closed for unbound commercial/redistribution/AI rights",
		EvidenceRef:       "#79 provider/data-health production closure plus exact-head regression evidence for " + provider + " / " + dataset,
		ApprovalRef:       "existing production route inherited from the #79/#84 governed provider program; any semantic change requires fresh exact-head qualification",
		InvalidationRule:  "adapter/schema/timestamp/freshness/failure/rights/priority/config semantics or relevant source fingerprint change invalidates inherited evidence",
		ExpectedDelay:     expectedDelay,
		Uses:              append([]string(nil), uses...),
	}
}

func tradeInsightProductionHistoryRoute(priority int) ProviderDatasetContract {
	return ProviderDatasetContract{
		Dataset: canonicalHistoricalBarsDataset, Capability: "Adjusted daily OHLCV / corporate-action corroboration", Priority: priority,
		Lifecycle: providerlifecycle.Production, ContractVersion: providerRegistrationContractVersion,
		CanonicalOwner:    canonicalHistoricalBarsDataset,
		Consumer:          "Historical Bars, Research, Corporate Actions, Adaptive Intelligence",
		AdapterContract:   "native Go TradeInsight daily-history adapter through the canonical Historical Bars loader",
		SchemaContract:    "verified /trading-data/v1/ohlc daily OHLCV contract; raw/adjusted and corporate-action semantics remain explicit",
		TimestampContract: "completed daily bars normalized by canonical Historical Bars owner",
		FreshnessContract: "completed-bar cadence evaluated by canonical Historical Bars freshness; never presented as live/intraday truth",
		FailureContract:   "TradeInsight adapter failures remain capability-scoped and fall through Smart Router v2 without poisoning unrelated capabilities",
		RightsContract:    "operational entitlement is separate from provider data-rights/commercial readiness; unbound rights fail closed",
		EvidenceRef:       "#78 qualification evidence and #84 production adoption evidence",
		ApprovalRef:       "#84 exact-head governed production adoption",
		InvalidationRule:  "endpoint/schema/adjustment/timestamp/rights/authority or adapter behavior change requires new qualification and approval evidence",
		ExpectedDelay:     "Completed-bar cadence",
		Uses:              []string{"Historical Bars", "Research", "Corporate Actions", "Adaptive Intelligence"},
	}
}

func providerRegistrations() []ProviderRegistration {
	has := func(v string) bool { return strings.TrimSpace(v) != "" }
	coreUses := []string{"Market Regime", "Day", "Swing", "Long", "Discovery", "Research", "Decision Queue", "Trade Readiness"}
	return []ProviderRegistration{
		{Name: "Alpaca", QuotaLabel: "Entitlement / feed dependent", CostClass: "Broker/data entitlement",
			Configured: func(_ Settings, s Secrets) bool { return has(s.AlpacaKey) && has(s.AlpacaSecret) },
			ConfigurationFingerprint: func(_ Settings, s Secrets) [32]byte {
				return providerConfigurationFingerprint("Alpaca", strings.TrimSpace(s.AlpacaKey), strings.TrimSpace(s.AlpacaSecret))
			},
			Routes: []ProviderDatasetContract{
				inheritedProductionRoute("Alpaca", "US Live Equities", "IEX quotes / snapshots / liquidity", 1, "Near-live when entitled", "Day", "Market Open Prep", "Trade Readiness", "Decision Queue"),
				inheritedProductionRoute("Alpaca", canonicalUSAssetUniverseDataset, "U.S. asset universe / instrument identity", 1, "Provider asset-directory cadence", "Symbol Validation", "Instrument Identity"),
				inheritedProductionRoute("Alpaca", canonicalUSMarketCalendarDataset, "U.S. market calendar", 1, "Official market-calendar cadence", "Session Awareness", "Preparation Jobs"),
				inheritedProductionRoute("Alpaca", canonicalUSCorporateActionsDataset, "Corporate actions", 1, "Corporate-action dissemination cadence", "Corporate Actions", "Historical Normalization", "Research"),
				inheritedProductionRoute("Alpaca", canonicalHistoricalBarsDataset, "Historical OHLCV", 1, "Completed-bar cadence", "Historical Bars", "Research", "Adaptive Intelligence"),
			},
			Diagnostics: []ProviderCapabilityDiagnosticRegistration{
				{Capability: "IEX quotes / snapshots / liquidity", Detail: "Bid/ask, size, spread, quote age and snapshot hydration.", Uses: []string{"Day", "Market Open Prep", "Trade Readiness", "Decision Queue"}, Status: func(_ Settings, s Secrets, health map[string]string, _ map[string]SymbolIntelligence, _ map[string]GlobalDriver) string {
					return capabilityStatusFromHealth(has(s.AlpacaKey) && has(s.AlpacaSecret), health["alpaca-live"])
				}},
				{Capability: "SIP movers / most active", Detail: "Discovery seed only when account entitlement permits.", Uses: []string{"Discovery", "Dashboard"}, Status: func(_ Settings, s Secrets, health map[string]string, _ map[string]SymbolIntelligence, _ map[string]GlobalDriver) string {
					if !has(s.AlpacaKey) || !has(s.AlpacaSecret) {
						return "NOT ENTITLED"
					}
					if strings.Contains(strings.ToLower(health["market-activity"]), "available") {
						return "AVAILABLE"
					}
					return "PLAN LIMITED"
				}},
			}},
		{Name: "Finnhub", QuotaLabel: "API plan / endpoint dependent", CostClass: "Free tier / optional paid upgrade",
			Configured: func(_ Settings, s Secrets) bool { return has(s.Finnhub) },
			ConfigurationFingerprint: func(_ Settings, s Secrets) [32]byte {
				return providerConfigurationFingerprint("Finnhub", strings.TrimSpace(s.Finnhub))
			},
			Routes: []ProviderDatasetContract{
				inheritedProductionRoute("Finnhub", "US Live Equities", "Primary U.S. equity", 2, "Near-live when entitled", "Day", "Swing", "Long", "Discovery", "Trade Readiness"),
				inheritedProductionRoute("Finnhub", "News", "Company news", 1, "Near-live when entitled", "News", "Research", "Catalyst Watch"),
				inheritedProductionRoute("Finnhub", "Earnings", "Earnings calendar / surprise", 1, "Provider event/filing cadence", "Earnings", "Catalyst Watch", "Research"),
				inheritedProductionRoute("Finnhub", "Fundamentals", "Company fundamentals", 1, "Provider fundamental-update cadence", "Research", "Swing", "Long"),
			},
			Diagnostics: []ProviderCapabilityDiagnosticRegistration{
				{Capability: "Primary U.S. equity + earnings/peers", Detail: "Live/REST plus earnings surprise and peer context.", Uses: append([]string(nil), coreUses...), Status: func(_ Settings, s Secrets, health map[string]string, _ map[string]SymbolIntelligence, _ map[string]GlobalDriver) string {
					return capabilityStatusFromHealth(has(s.Finnhub), health["quotes-rest"])
				}},
				{Capability: "Analyst / insider premium context", Detail: "Used only when endpoint entitlement returns real data; never required for deterministic scoring.", Uses: []string{"Research", "Swing", "Long", "Decision Queue"}, Status: func(_ Settings, s Secrets, _ map[string]string, intel map[string]SymbolIntelligence, _ map[string]GlobalDriver) string {
					if !has(s.Finnhub) {
						return "NOT ENTITLED"
					}
					for _, v := range intel {
						if v.RecommendationTrend != "" || v.PriceTarget > 0 || v.InsiderNetShares != 0 {
							return "AVAILABLE"
						}
					}
					return "PLAN LIMITED"
				}},
			}},
		{Name: tradeInsightProviderName, QuotaLabel: "Runtime tier / response headers", CostClass: "Beta / free tier",
			Configured: func(_ Settings, s Secrets) bool { return tradeInsightConfigured(s.TradeInsight) },
			ConfigurationFingerprint: func(_ Settings, s Secrets) [32]byte {
				return providerConfigurationFingerprint(tradeInsightProviderName, strings.TrimSpace(s.TradeInsight))
			},
			Routes: []ProviderDatasetContract{tradeInsightProductionHistoryRoute(2)},
			Diagnostics: []ProviderCapabilityDiagnosticRegistration{
				{Capability: "Adjusted daily OHLCV / corporate-action corroboration", Detail: "Smart Router v2 member of canonical Historical Bars. Daily adjusted OHLCV only; dividends/splits are supplemental evidence merged into the canonical corporate-action ledger.", Uses: []string{"Historical Bars", "Research", "Corporate Actions", "Adaptive Intelligence"}, Status: func(_ Settings, s Secrets, health map[string]string, _ map[string]SymbolIntelligence, _ map[string]GlobalDriver) string {
					if !tradeInsightConfigured(s.TradeInsight) {
						return "NOT ENTITLED"
					}
					history := strings.ToLower(strings.TrimSpace(health["history"]))
					actions := strings.ToLower(strings.TrimSpace(health["tradeinsight-corporate-actions"]))
					if strings.Contains(history, "tradeinsight") || strings.Contains(actions, "healthy") {
						return "AVAILABLE"
					}
					return "SHADOW"
				}},
			}},
		{Name: "Twelve Data", QuotaLabel: "Credit based", CostClass: "Free tier / optional paid upgrade",
			Configured: func(_ Settings, s Secrets) bool { return has(s.TwelveData) },
			ConfigurationFingerprint: func(s Settings, sec Secrets) [32]byte {
				return providerConfigurationFingerprint("Twelve Data", strings.TrimSpace(sec.TwelveData), strings.ToLower(strings.TrimSpace(s.GlobalProviderMode)))
			},
			Routes: []ProviderDatasetContract{
				inheritedProductionRoute("Twelve Data", "US Live Equities", "U.S. equity fallback", 3, "Near-live when entitled", "Data Freshness", "Recovery"),
				inheritedProductionRoute("Twelve Data", canonicalGlobalMarketContextDataset, "FX / direct global context", 1, "Near-live when entitled", "Market Regime", "Dashboard", "Research"),
				inheritedProductionRoute("Twelve Data", "VIX / Indices", "VIX / indices", 1, "Near-live when entitled", "VIX", "Market Regime", "Day", "Swing", "Long"),
				inheritedProductionRoute("Twelve Data", canonicalHistoricalBarsDataset, "Historical OHLCV fallback", 3, "Completed-bar cadence", "Historical Bars", "Recovery"),
			},
			Diagnostics: []ProviderCapabilityDiagnosticRegistration{
				{Capability: "FX / direct global context", Detail: "Direct FX/global where entitled; official/proxy/cache fallback remains truthful.", Uses: []string{"Market Regime", "Dashboard", "Research"}, Status: func(_ Settings, s Secrets, health map[string]string, _ map[string]SymbolIntelligence, direct map[string]GlobalDriver) string {
					if !has(s.TwelveData) {
						return "NOT ENTITLED"
					}
					for k := range direct {
						if strings.HasPrefix(k, "fx_") {
							return "AVAILABLE"
						}
					}
					return capabilityStatusFromHealth(true, health["global-direct"])
				}},
				{Capability: "VIX / indices / historical recovery", Detail: "v15 primary VIX/index route and first historical fallback after Alpaca.", Uses: []string{"Market Regime", "Dashboard", "Day", "Swing", "Long"}, Status: func(_ Settings, s Secrets, health map[string]string, _ map[string]SymbolIntelligence, _ map[string]GlobalDriver) string {
					return capabilityStatusFromHealth(has(s.TwelveData), health["vix"])
				}},
			}},
		{Name: "Marketaux", QuotaLabel: "Request quota", CostClass: "Free tier / optional paid upgrade",
			Configured: func(_ Settings, s Secrets) bool { return has(s.Marketaux) },
			ConfigurationFingerprint: func(_ Settings, s Secrets) [32]byte {
				return providerConfigurationFingerprint("Marketaux", strings.TrimSpace(s.Marketaux))
			},
			Routes: []ProviderDatasetContract{inheritedProductionRoute("Marketaux", "News", "Stock news fallback", 2, "Near-live when entitled", "News", "Dashboard", "Research", "Catalyst Watch")},
			Diagnostics: []ProviderCapabilityDiagnosticRegistration{
				{Capability: "Stock news fallback", Detail: "Supplemental/fallback company news when Finnhub is unavailable.", Uses: []string{"News", "Dashboard", "Research", "Catalyst Watch"}, Status: func(_ Settings, s Secrets, _ map[string]string, _ map[string]SymbolIntelligence, _ map[string]GlobalDriver) string {
					if has(s.Marketaux) {
						return "AVAILABLE"
					}
					return "NOT ENTITLED"
				}},
			}},
		{Name: "FRED", QuotaLabel: "Free API key", CostClass: "Free tier / optional paid upgrade",
			Configured: func(_ Settings, s Secrets) bool { return has(s.FRED) },
			ConfigurationFingerprint: func(_ Settings, s Secrets) [32]byte {
				return providerConfigurationFingerprint("FRED", strings.TrimSpace(s.FRED))
			},
			Routes: []ProviderDatasetContract{inheritedProductionRoute("FRED", "Macro", "Rates / credit / conditions / USD", 1, "Series release cadence", "Market Regime", "Swing", "Long", "Research", "Trade Readiness")},
			Diagnostics: []ProviderCapabilityDiagnosticRegistration{
				{Capability: "Rates / credit / conditions / USD", Detail: "Slow macro state; cadence-aware cache.", Uses: []string{"Market Regime", "Swing", "Long", "Research", "Trade Readiness"}, Status: func(_ Settings, s Secrets, health map[string]string, _ map[string]SymbolIntelligence, _ map[string]GlobalDriver) string {
					return capabilityStatusFromHealth(has(s.FRED), health["fred-rates"])
				}},
			}},
		{Name: "SEC", QuotaLabel: "Fair-access policy", CostClass: "Public / no API fee",
			Configured: func(s Settings, _ Secrets) bool { return has(s.SECEmail) },
			ConfigurationFingerprint: func(s Settings, _ Secrets) [32]byte {
				return providerConfigurationFingerprint("SEC", strings.TrimSpace(s.SECEmail))
			},
			Routes: []ProviderDatasetContract{inheritedProductionRoute("SEC", "Fundamentals", "SEC fundamental authority/fallback", 2, "Filing dissemination cadence", "Research", "Fundamentals")}},
		{Name: "SEC EDGAR", QuotaLabel: "Fair-access policy", CostClass: "Public / no API fee",
			Configured: func(s Settings, _ Secrets) bool { return has(s.SECEmail) },
			ConfigurationFingerprint: func(s Settings, _ Secrets) [32]byte {
				return providerConfigurationFingerprint("SEC EDGAR", strings.TrimSpace(s.SECEmail))
			},
			Routes: []ProviderDatasetContract{inheritedProductionRoute("SEC EDGAR", "SEC", "Direct SEC/EDGAR filings and Form 4 authority", 1, "Filing dissemination cadence", "SEC", "Research", "Catalyst Watch")}},
		{Name: "yfinance", QuotaLabel: "Public recovery · best effort", CostClass: "Public / no API fee",
			Configured: func(_ Settings, _ Secrets) bool { return true },
			ConfigurationFingerprint: func(_ Settings, _ Secrets) [32]byte {
				return providerConfigurationFingerprint("yfinance", "public-recovery")
			},
			Routes: []ProviderDatasetContract{
				inheritedProductionRoute("yfinance", "VIX / Indices", "Recovery-only VIX/index context", 2, "Recovery-only; may be delayed", "Data Freshness", "Recovery"),
				inheritedProductionRoute("yfinance", canonicalHistoricalBarsDataset, "Recovery-only historical bars", 4, "Recovery-only; may be delayed", "Historical Bars", "Recovery"),
				inheritedProductionRoute("yfinance", "Earnings", "Recovery-only earnings", 2, "Recovery-only; may be delayed", "Earnings", "Recovery"),
				inheritedProductionRoute("yfinance", "Fundamentals", "Recovery-only fundamentals", 3, "Recovery-only; may be delayed", "Fundamentals", "Recovery"),
			},
			Diagnostics: []ProviderCapabilityDiagnosticRegistration{
				{Capability: "Recovery-only public market context", Detail: "Fallback only for VIX, historical bars, earnings and fundamentals; never the primary live production feed.", Uses: []string{"Data Freshness", "Research", "Recovery"}, Status: func(_ Settings, _ Secrets, _ map[string]string, _ map[string]SymbolIntelligence, _ map[string]GlobalDriver) string {
					return "AVAILABLE"
				}},
			}},
		{Name: "CBOE", QuotaLabel: "Public official/delayed", CostClass: "Public / no API fee",
			Configured: func(_ Settings, _ Secrets) bool { return true },
			ConfigurationFingerprint: func(_ Settings, _ Secrets) [32]byte {
				return providerConfigurationFingerprint("CBOE", "public-official")
			},
			Routes: []ProviderDatasetContract{inheritedProductionRoute("CBOE", "VIX / Indices", "Official VIX validation / delayed close", 3, "Official delayed/close validation", "VIX", "Market Regime", "Data Freshness")},
			Diagnostics: []ProviderCapabilityDiagnosticRegistration{
				{Capability: "Official VIX validation / delayed close", Detail: "Authoritative VIX validation and delayed/official fallback only; not a general stock provider.", Uses: []string{"VIX", "Market Regime", "Data Freshness"}, Status: func(_ Settings, _ Secrets, _ map[string]string, _ map[string]SymbolIntelligence, _ map[string]GlobalDriver) string {
					return "AVAILABLE"
				}},
			}},
		{Name: "BLS", QuotaLabel: "Official API", CostClass: "Public / free API",
			Configured: func(_ Settings, _ Secrets) bool { return true },
			ConfigurationFingerprint: func(_ Settings, s Secrets) [32]byte {
				return providerConfigurationFingerprint("BLS", strings.TrimSpace(s.BLS))
			},
			Diagnostics: []ProviderCapabilityDiagnosticRegistration{
				{Capability: "Inflation / labor / wages / PPI", Detail: "Official release-triggered actuals; no invented consensus.", Uses: []string{"Market Regime", "Swing", "Long", "Research"}, Status: func(_ Settings, _ Secrets, health map[string]string, _ map[string]SymbolIntelligence, _ map[string]GlobalDriver) string {
					return capabilityStatusFromHealth(true, health["bls-actuals"])
				}},
			}},
		{Name: "EIA", QuotaLabel: "Free API key", CostClass: "Public / free API",
			Configured: func(_ Settings, s Secrets) bool { return has(s.EIA) },
			ConfigurationFingerprint: func(_ Settings, s Secrets) [32]byte {
				return providerConfigurationFingerprint("EIA", strings.TrimSpace(s.EIA))
			},
			Diagnostics: []ProviderCapabilityDiagnosticRegistration{
				{Capability: "Petroleum / natural gas / energy state", Detail: "Official energy release context.", Uses: []string{"Market Regime", "Research", "Swing", "Long"}, Status: func(_ Settings, s Secrets, health map[string]string, _ map[string]SymbolIntelligence, _ map[string]GlobalDriver) string {
					return capabilityStatusFromHealth(has(s.EIA), health["eia-actuals"])
				}},
			}},
	}
}

func providerDatasetContractValidationErrors(c ProviderDatasetContract) []string {
	var out []string
	require := func(name, value string) {
		if strings.TrimSpace(value) == "" {
			out = append(out, name+" missing")
		}
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
	if c.Priority <= 0 {
		out = append(out, "priority missing")
	}
	if len(c.Uses) == 0 {
		out = append(out, "downstream uses missing")
	}
	return out
}

func providerDatasetContractProductionReady(c ProviderDatasetContract) bool {
	return providerlifecycle.CanonicalLifecycle(c.Lifecycle) == providerlifecycle.Production && len(providerDatasetContractValidationErrors(c)) == 0
}

func routeChainsFromProviderRegistrations(regs []ProviderRegistration) map[string][]string {
	type ranked struct {
		name     string
		priority int
	}
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
			if rows[i].priority != rows[j].priority {
				return rows[i].priority < rows[j].priority
			}
			return rows[i].name < rows[j].name
		})
		for _, row := range rows {
			out[dataset] = append(out[dataset], row.name)
		}
	}
	return out
}

func providerRegistrationLookup(provider string) (ProviderRegistration, bool) {
	want := providerKey(provider)
	for _, reg := range providerRegistrations() {
		if providerKey(reg.Name) == want {
			return reg, true
		}
	}
	return ProviderRegistration{}, false
}

func providerConfiguredFromRegistration(provider string, settings Settings, secrets Secrets) bool {
	reg, ok := providerRegistrationLookup(provider)
	return ok && reg.Configured != nil && reg.Configured(settings, secrets)
}

func providerQuotaFromRegistration(provider string) string {
	if reg, ok := providerRegistrationLookup(provider); ok && strings.TrimSpace(reg.QuotaLabel) != "" {
		return reg.QuotaLabel
	}
	return "Provider dependent"
}

func providerCostFromRegistration(provider string) string {
	if reg, ok := providerRegistrationLookup(provider); ok && strings.TrimSpace(reg.CostClass) != "" {
		return reg.CostClass
	}
	return "Provider dependent"
}

func providerExpectedDelayFromRegistration(dataset, provider string) string {
	if reg, ok := providerRegistrationLookup(provider); ok {
		for _, route := range reg.Routes {
			if route.Dataset == dataset && strings.TrimSpace(route.ExpectedDelay) != "" {
				return route.ExpectedDelay
			}
		}
	}
	return "Near-live when entitled"
}

const adaptiveProviderRegistryContractVersion = "adaptive-provider-registry-v19.1.0"

const (
	providerAuthorityRoutable        = "ROUTABLE"
	providerAuthorityFallback        = "FALLBACK"
	providerAuthorityCorroborative   = "CORROBORATIVE"
	providerAuthorityDirectAuthority = "DIRECT_AUTHORITY"
	providerAuthorityNonRoutable     = "NON_ROUTABLE"

	providerRegistryCapabilityRegistered              = "REGISTERED"
	providerRegistryCapabilityRegisteredNotProduction = "REGISTERED_NOT_PRODUCTION"
)

// AdaptiveProviderRouteManifest is a read-only projection of the existing
// ProviderRegistration route contract. It does not own routing, health,
// freshness, lifecycle, rights, cache, persistence, subscriptions or canonical
// state. Smart Provider Router v2 remains the sole general selection authority.
type AdaptiveProviderRouteManifest struct {
	Dataset           string   `json:"dataset"`
	Capability        string   `json:"capability"`
	InstrumentClass   string   `json:"instrumentClass"`
	AuthorityClass    string   `json:"authorityClass"`
	Lifecycle         string   `json:"lifecycle"`
	ContractVersion   string   `json:"contractVersion"`
	CanonicalOwner    string   `json:"canonicalOwner"`
	SchemaContract    string   `json:"schemaContract"`
	TimestampContract string   `json:"timestampContract"`
	FreshnessContract string   `json:"freshnessContract"`
	RightsContract    string   `json:"rightsContract"`
	ExpectedDelay     string   `json:"expectedDelay,omitempty"`
	Uses              []string `json:"uses,omitempty"`
	ProductionReady   bool     `json:"productionReady"`
}

// AdaptiveProviderDiagnosticManifest exposes registered diagnostic capability
// metadata without exposing the diagnostic function or any credential material.
type AdaptiveProviderDiagnosticManifest struct {
	Capability string   `json:"capability"`
	Detail     string   `json:"detail,omitempty"`
	Uses       []string `json:"uses,omitempty"`
}

// AdaptiveProviderManifest is the generic manifest projected from the existing
// canonical ProviderRegistration descriptor. Configuration is represented only
// as presence/absence; secrets are never copied into this projection.
type AdaptiveProviderManifest struct {
	ProviderID      string                               `json:"providerId"`
	Name            string                               `json:"name"`
	ContractVersion string                               `json:"contractVersion"`
	Configured      bool                                 `json:"configured"`
	Configuration   string                               `json:"configuration"`
	QuotaLabel      string                               `json:"quotaLabel,omitempty"`
	CostClass       string                               `json:"costClass,omitempty"`
	Routable        bool                                 `json:"routable"`
	Routes          []AdaptiveProviderRouteManifest      `json:"routes,omitempty"`
	Diagnostics     []AdaptiveProviderDiagnosticManifest `json:"diagnostics,omitempty"`
}

// AdaptiveProviderRegistrySnapshot is intentionally a computed snapshot rather
// than a second runtime state owner. All entries come from providerRegistrations
// (or an injected registration slice in tests); Router eligibility remains in
// Smart Provider Router v2 and Data Health/freshness remains separately owned.
type AdaptiveProviderRegistrySnapshot struct {
	ContractVersion string                     `json:"contractVersion"`
	Providers       []AdaptiveProviderManifest `json:"providers"`
}

func adaptiveProviderAuthorityClass(provider string, route ProviderDatasetContract) string {
	switch {
	case strings.EqualFold(provider, "SEC EDGAR") && strings.EqualFold(route.Dataset, "SEC"):
		return providerAuthorityDirectAuthority
	case strings.EqualFold(provider, "SEC"):
		return providerAuthorityCorroborative
	case strings.EqualFold(provider, "CBOE"):
		return providerAuthorityCorroborative
	case strings.EqualFold(provider, "yfinance"):
		return providerAuthorityFallback
	default:
		return providerAuthorityRoutable
	}
}

func adaptiveProviderRouteManifest(provider string, route ProviderDatasetContract) AdaptiveProviderRouteManifest {
	return AdaptiveProviderRouteManifest{
		Dataset:           route.Dataset,
		Capability:        route.Capability,
		InstrumentClass:   providerInstrumentClass(route.Dataset),
		AuthorityClass:    adaptiveProviderAuthorityClass(provider, route),
		Lifecycle:         route.Lifecycle,
		ContractVersion:   route.ContractVersion,
		CanonicalOwner:    route.CanonicalOwner,
		SchemaContract:    route.SchemaContract,
		TimestampContract: route.TimestampContract,
		FreshnessContract: route.FreshnessContract,
		RightsContract:    route.RightsContract,
		ExpectedDelay:     route.ExpectedDelay,
		Uses:              append([]string(nil), route.Uses...),
		ProductionReady:   providerDatasetContractProductionReady(route),
	}
}

func adaptiveProviderManifest(reg ProviderRegistration, settings Settings, secrets Secrets) AdaptiveProviderManifest {
	configured := reg.Configured != nil && reg.Configured(settings, secrets)
	manifest := AdaptiveProviderManifest{
		ProviderID:      providerKey(reg.Name),
		Name:            reg.Name,
		ContractVersion: adaptiveProviderRegistryContractVersion,
		Configured:      configured,
		Configuration:   providerCapabilityNotConfigured,
		QuotaLabel:      reg.QuotaLabel,
		CostClass:       reg.CostClass,
		Routable:        len(reg.Routes) > 0,
	}
	if configured {
		manifest.Configuration = "CONFIGURED"
	}
	for _, route := range reg.Routes {
		manifest.Routes = append(manifest.Routes, adaptiveProviderRouteManifest(reg.Name, route))
	}
	for _, diagnostic := range reg.Diagnostics {
		manifest.Diagnostics = append(manifest.Diagnostics, AdaptiveProviderDiagnosticManifest{
			Capability: diagnostic.Capability,
			Detail:     diagnostic.Detail,
			Uses:       append([]string(nil), diagnostic.Uses...),
		})
	}
	sort.SliceStable(manifest.Routes, func(i, j int) bool {
		if manifest.Routes[i].Dataset != manifest.Routes[j].Dataset {
			return manifest.Routes[i].Dataset < manifest.Routes[j].Dataset
		}
		if manifest.Routes[i].Capability != manifest.Routes[j].Capability {
			return manifest.Routes[i].Capability < manifest.Routes[j].Capability
		}
		return manifest.Routes[i].Lifecycle < manifest.Routes[j].Lifecycle
	})
	sort.SliceStable(manifest.Diagnostics, func(i, j int) bool {
		return strings.ToLower(manifest.Diagnostics[i].Capability) < strings.ToLower(manifest.Diagnostics[j].Capability)
	})
	return manifest
}

func adaptiveProviderRegistrySnapshotFromRegistrations(regs []ProviderRegistration, settings Settings, secrets Secrets) AdaptiveProviderRegistrySnapshot {
	out := AdaptiveProviderRegistrySnapshot{ContractVersion: adaptiveProviderRegistryContractVersion}
	for _, reg := range regs {
		if strings.TrimSpace(reg.Name) == "" {
			continue
		}
		out.Providers = append(out.Providers, adaptiveProviderManifest(reg, settings, secrets))
	}
	sort.SliceStable(out.Providers, func(i, j int) bool {
		return out.Providers[i].ProviderID < out.Providers[j].ProviderID
	})
	return out
}

func adaptiveProviderRegistrySnapshot(settings Settings, secrets Secrets) AdaptiveProviderRegistrySnapshot {
	return adaptiveProviderRegistrySnapshotFromRegistrations(providerRegistrations(), settings, secrets)
}

// adaptiveProviderRegistryCapabilityState answers only registration/configuration
// truth. It deliberately does not answer Router eligibility or provider health.
// Those remain Smart Provider Router v2/Data Health responsibilities.
func adaptiveProviderRegistryCapabilityStateFromRegistrations(regs []ProviderRegistration, provider, dataset string, settings Settings, secrets Secrets) string {
	wantProvider := providerKey(provider)
	wantDataset := strings.TrimSpace(dataset)
	for _, reg := range regs {
		if providerKey(reg.Name) != wantProvider {
			continue
		}
		for _, route := range reg.Routes {
			if !strings.EqualFold(strings.TrimSpace(route.Dataset), wantDataset) {
				continue
			}
			if reg.Configured == nil || !reg.Configured(settings, secrets) {
				return providerCapabilityNotConfigured
			}
			if !providerDatasetContractProductionReady(route) {
				return providerRegistryCapabilityRegisteredNotProduction
			}
			return providerRegistryCapabilityRegistered
		}
		return providerCapabilityNotSupported
	}
	return providerCapabilityNotSupported
}

func adaptiveProviderRegistryCapabilityState(provider, dataset string, settings Settings, secrets Secrets) string {
	return adaptiveProviderRegistryCapabilityStateFromRegistrations(providerRegistrations(), provider, dataset, settings, secrets)
}
