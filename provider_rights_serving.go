package main

import (
	"sort"
	"strings"
	"time"
)

func hostedSECServingAllowed(now time.Time) bool {
	return hostedProviderRightsAllowed("SEC EDGAR", providerHostedUseProductionServing, now)
}

func hostedSECDerivedServingAllowed(now time.Time) bool {
	return hostedSECServingAllowed(now) && hostedProviderRightsAllowed("SEC EDGAR", providerHostedUseAI, now)
}

func hostedRightsFilteredOptions(options map[string]OptionsContext, now time.Time) map[string]OptionsContext {
	if !isHostedRuntime() {
		return options
	}
	out := make(map[string]OptionsContext, len(options))
	for symbol, item := range options {
		provider := strings.TrimSpace(item.Provider)
		if provider != "" && hostedProviderRightsAllowed(provider, providerHostedUseProductionServing, now) {
			out[symbol] = item
		}
	}
	return out
}

func hostedRightsBlockedProviders(router ProviderRouterSnapshot) []string {
	if !isHostedRuntime() {
		return nil
	}
	seen := map[string]bool{}
	for _, route := range router.Routes {
		for _, hop := range route.Route {
			if hop.Configured && strings.EqualFold(hop.Health, "RIGHTS BLOCKED") {
				seen[hop.Provider] = true
			}
		}
	}
	out := make([]string, 0, len(seen))
	for provider := range seen {
		out = append(out, provider)
	}
	sort.Strings(out)
	return out
}

// enforceHostedRuntimeRightsSnapshot is the final shared-data serving guard.
// It extends the existing RuntimeSnapshot owner rather than introducing a new
// serving model. Provider-attributed quotes/options/SEC data are admitted only
// while their current rights evidence remains valid. Legacy mixed datasets that
// do not carry provider identity per record fail closed in hosted mode until the
// canonical provenance work can bind them (HOST-022); desktop remains unchanged.
func enforceHostedRuntimeRightsSnapshot(snap RuntimeSnapshot) RuntimeSnapshot {
	if !isHostedRuntime() {
		return snap
	}
	now := time.Now()
	out := snap
	out.Health = clone(snap.Health)
	if out.Health == nil {
		out.Health = map[string]string{}
	}

	out.Quotes = hostedRightsFilterQuotes(snap.Quotes, now)
	// These legacy cache/runtime collections do not retain provider identity on
	// every row. Serving them in a multi-user hosted process would make expiry or
	// downgrade non-enforceable, so they remain unavailable until provenance is
	// explicit rather than inferred from a successful connection.
	out.History = nil
	out.Bars = nil
	out.Fundamentals = nil
	out.News = nil
	out.Earnings = nil
	out.Global = GlobalMarketContext{}
	out.MacroMetrics = nil
	out.MacroEvents = nil
	out.EventMode = EventModeState{}
	out.EventReactions = nil
	out.Options = hostedRightsFilteredOptions(snap.Options, now)

	if hostedSECServingAllowed(now) {
		out.Filings = snap.Filings
		out.CorporateActions = snap.CorporateActions
	} else {
		out.Filings = nil
		out.CorporateActions = nil
	}
	if hostedSECDerivedServingAllowed(now) {
		out.SECIntelligence = snap.SECIntelligence
	} else {
		out.SECIntelligence = nil
	}

	// Derived/adaptive outputs currently aggregate provider data without a
	// provider-rights lineage on every output. Keep operational diagnostics, the
	// sole Smart Router, capabilities, workload/SLO and manual controls visible;
	// fail closed on market-derived payloads until HOST-022 binds that lineage.
	out.Scanner = ScannerState{}
	out.SignalValidation = SignalValidationState{}
	out.ValidationLearning = ValidationLearningSnapshot{}
	out.Preparations = nil
	out.Liquidity = nil
	out.Intelligence = nil
	out.SymbolIntelligence = nil
	out.CatalystReactions = nil
	out.MarketOpenFlags = nil
	out.MarketOpenCheckpoint = MarketOpenCheckpoint{}
	out.MarketActivity = MarketActivityState{}
	out.LiveCoverage = nil
	out.RapidMove = RapidMoveState{}
	out.ProviderReconciliation = nil
	out.ResearchPackage = ResearchPackageTruth{}
	out.EvidenceSnapshot = EvidenceSnapshot{}
	out.CorporateActionTruth = CorporateActionTruth{}
	out.MarketIntelligence = MarketIntelligenceSnapshot{}
	out.EventIntelligence = EventIntelligenceSnapshot{}
	out.AlternativeIntelligence = ContextAlternativeIntelligenceSnapshot{}
	out.AdaptiveDataPolicy = AdaptiveDataPolicyState{}
	out.ShadowControl = ShadowControlState{}

	if !hostedProviderRightsAllowed("Finnhub", providerHostedUseProductionServing, now) {
		out.Feed.WebSocketConnected = false
		out.Feed.SubscribedSymbols = nil
		out.Feed.LastTradeAt = 0
		out.Feed.LastTradeSymbol = ""
	}
	if !hostedProviderRightsAllowed("Alpaca", providerHostedUseProductionServing, now) {
		out.Feed.AlpacaWebSocketConnected = false
		out.Feed.AlpacaSubscribedSymbols = nil
		out.Feed.LastAlpacaStreamAt = 0
		out.Feed.LastAlpacaStreamSymbol = ""
		out.Feed.LastAlpacaAt = 0
		out.Feed.LastAlpacaSymbol = ""
	}
	out.Feed.LiveSymbols = uniqueLiveSubscriptionCount(boolSliceToMap(out.Feed.SubscribedSymbols), boolSliceToMap(out.Feed.AlpacaSubscribedSymbols))

	blocked := hostedRightsBlockedProviders(out.ProviderRouter)
	removedQuotes := len(out.Quotes) < len(snap.Quotes)
	if len(blocked) > 0 || removedQuotes {
		out.Health["provider-rights"] = "BLOCKED · hosted serving denied for unapproved/expired provider data"
		if out.Status == "running" {
			out.Status = "degraded"
		}
		out.Message = "Hosted provider rights are blocking unapproved market-data serving."
	} else {
		out.Health["provider-rights"] = "READY · hosted provider-rights serving guard active"
	}
	return out
}

func boolSliceToMap(values []string) map[string]bool {
	out := make(map[string]bool, len(values))
	for _, value := range values {
		out[value] = true
	}
	return out
}

func hostedQuoteServingAllowedForSymbol(a *Application, symbol string) bool {
	if !isHostedRuntime() {
		return true
	}
	if a == nil || a.engine == nil {
		return false
	}
	symbol = normalizeSymbol(symbol)
	a.engine.mu.RLock()
	q, ok := a.engine.quotes[symbol]
	a.engine.mu.RUnlock()
	return ok && providerQuoteHostedRightsAllowed(q, providerHostedUseProductionServing, time.Now())
}

func hostedNewsItemsForServing(items []NewsItem) []NewsItem {
	if !isHostedRuntime() {
		return items
	}
	// NewsItem.Source is publisher identity, not the provider that supplied the
	// record. Provider rights therefore cannot be proven per item yet.
	return nil
}

func hostedEarningsItemsForServing(items []EarningsItem) []EarningsItem {
	if !isHostedRuntime() {
		return items
	}
	// EarningsItem currently has no provider provenance field.
	return nil
}

func hostedFilingItemsForServing(items []FilingItem) []FilingItem {
	if !isHostedRuntime() || hostedSECServingAllowed(time.Now()) {
		return items
	}
	return nil
}

func hostedSECIntelligenceForServing(items map[string]SECIntelligenceSummary) map[string]SECIntelligenceSummary {
	if !isHostedRuntime() || hostedSECDerivedServingAllowed(time.Now()) {
		return items
	}
	return nil
}
