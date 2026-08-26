package main

import (
	"sort"
	"strings"
	"time"
)

// runtimeAllowedSymbolsLocked returns the symbols a user may see in raw,
// symbol-keyed runtime payloads. Always-on market context is shared globally;
// personal actionable symbols come only from the authenticated workspace.
func (a *Application) runtimeAllowedSymbolsLocked(userID string) map[string]bool {
	allowed := map[string]bool{}
	for _, sym := range append(append([]string{}, masterMarketSymbols...), specialIndexSymbols...) {
		if s := normalizeSymbol(sym); s != "" {
			allowed[s] = true
		}
	}
	workspace := a.workspaceStateLocked(userID)
	for _, sym := range analysisSymbolsFromState(workspace) {
		if s := normalizeSymbol(sym); s != "" {
			allowed[s] = true
		}
	}
	return allowed
}

func runtimeSymbolAllowed(allowed map[string]bool, symbol string) bool {
	symbol = normalizeSymbol(symbol)
	return symbol == "" || allowed[symbol]
}

func filterRuntimeSymbolMap[T any](in map[string]T, allowed map[string]bool) map[string]T {
	if in == nil {
		return nil
	}
	out := make(map[string]T, len(in))
	for symbol, value := range in {
		if runtimeSymbolAllowed(allowed, symbol) {
			out[symbol] = value
		}
	}
	return out
}

func filterRuntimeDiagnosticMap[T any](in map[string]T, allowed map[string]bool) map[string]T {
	if in == nil {
		return nil
	}
	out := make(map[string]T, len(in))
	for key, value := range in {
		parts := strings.Split(key, ":")
		last := strings.ToUpper(strings.TrimSpace(parts[len(parts)-1]))
		if len(parts) > 1 {
			if _, valid := parseSelectableTicker(last); valid && !runtimeSymbolAllowed(allowed, last) {
				continue
			}
		}
		out[key] = value
	}
	return out
}

func filterRuntimeSymbols(items []string, allowed map[string]bool) []string {
	out := make([]string, 0, len(items))
	for _, symbol := range items {
		if runtimeSymbolAllowed(allowed, symbol) {
			out = append(out, symbol)
		}
	}
	return uniqueSymbols(out)
}

func (a *Application) runtimeSnapshotForUserFrom(userID string, snap RuntimeSnapshot) RuntimeSnapshot {
	// Legal/data-rights admission precedes the existing per-user privacy filter.
	// This preserves one canonical shared engine while preventing an authenticated
	// consumer from reaching provider data that the hosted process may not serve.
	snap = enforceHostedRuntimeRightsSnapshot(snap)
	a.mu.RLock()
	allowed := a.runtimeAllowedSymbolsLocked(userID)
	a.mu.RUnlock()
	out := snap
	out.Quotes = filterRuntimeSymbolMap(snap.Quotes, allowed)
	out.History = filterRuntimeSymbolMap(snap.History, allowed)
	out.Bars = filterRuntimeSymbolMap(snap.Bars, allowed)
	out.Fundamentals = filterRuntimeSymbolMap(snap.Fundamentals, allowed)
	out.SECIntelligence = filterRuntimeSymbolMap(snap.SECIntelligence, allowed)
	out.Options = filterRuntimeSymbolMap(snap.Options, allowed)
	out.Liquidity = filterRuntimeSymbolMap(snap.Liquidity, allowed)
	out.SymbolIntelligence = filterRuntimeSymbolMap(snap.SymbolIntelligence, allowed)
	out.CatalystReactions = filterRuntimeSymbolMap(snap.CatalystReactions, allowed)
	out.MarketOpenFlags = filterRuntimeSymbolMap(snap.MarketOpenFlags, allowed)
	out.LiveCoverage = filterRuntimeSymbolMap(snap.LiveCoverage, allowed)
	out.Health = filterRuntimeDiagnosticMap(snap.Health, allowed)
	out.LastUpdated = filterRuntimeDiagnosticMap(snap.LastUpdated, allowed)

	out.News = make([]NewsItem, 0, len(snap.News))
	for _, item := range snap.News {
		if len(item.Symbols) == 0 {
			out.News = append(out.News, item)
			continue
		}
		keep := false
		for _, symbol := range item.Symbols {
			if runtimeSymbolAllowed(allowed, symbol) {
				keep = true
				break
			}
		}
		if keep {
			out.News = append(out.News, item)
		}
	}
	out.Earnings = make([]EarningsItem, 0, len(snap.Earnings))
	for _, item := range snap.Earnings {
		if runtimeSymbolAllowed(allowed, item.Symbol) {
			out.Earnings = append(out.Earnings, item)
		}
	}
	out.Filings = make([]FilingItem, 0, len(snap.Filings))
	for _, item := range snap.Filings {
		if runtimeSymbolAllowed(allowed, item.Symbol) {
			out.Filings = append(out.Filings, item)
		}
	}
	out.CorporateActions = make([]CorporateAction, 0, len(snap.CorporateActions))
	for _, item := range snap.CorporateActions {
		if runtimeSymbolAllowed(allowed, item.Symbol) || runtimeSymbolAllowed(allowed, item.OldSymbol) || runtimeSymbolAllowed(allowed, item.NewSymbol) {
			out.CorporateActions = append(out.CorporateActions, item)
		}
	}
	out.ProviderReconciliation = make([]ProviderReconciliationDecision, 0, len(snap.ProviderReconciliation))
	for _, item := range snap.ProviderReconciliation {
		if runtimeSymbolAllowed(allowed, item.Symbol) {
			out.ProviderReconciliation = append(out.ProviderReconciliation, item)
		}
	}
	out.SignalValidation = snap.SignalValidation
	out.SignalValidation.Snapshots = make([]SignalSnapshot, 0, len(snap.SignalValidation.Snapshots))
	for _, item := range snap.SignalValidation.Snapshots {
		if runtimeSymbolAllowed(allowed, item.Symbol) {
			out.SignalValidation.Snapshots = append(out.SignalValidation.Snapshots, item)
		}
	}
	out.Feed = snap.Feed
	out.Feed.SubscribedSymbols = filterRuntimeSymbols(snap.Feed.SubscribedSymbols, allowed)
	out.Feed.AlpacaSubscribedSymbols = filterRuntimeSymbols(snap.Feed.AlpacaSubscribedSymbols, allowed)
	out.Feed.SnapshotSymbols = filterRuntimeSymbols(snap.Feed.SnapshotSymbols, allowed)
	out.Feed.TrackedSymbols = len(filterRuntimeSymbols(append(append([]string{}, out.Feed.SubscribedSymbols...), out.Feed.SnapshotSymbols...), allowed))
	out.Feed.LiveSymbols = len(uniqueSymbols(append(append([]string{}, out.Feed.SubscribedSymbols...), out.Feed.AlpacaSubscribedSymbols...)))
	if !runtimeSymbolAllowed(allowed, out.Feed.LastTradeSymbol) {
		out.Feed.LastTradeSymbol = ""
	}
	if !runtimeSymbolAllowed(allowed, out.Feed.LastAlpacaStreamSymbol) {
		out.Feed.LastAlpacaStreamSymbol = ""
	}
	if !runtimeSymbolAllowed(allowed, out.Feed.LastAlpacaSymbol) {
		out.Feed.LastAlpacaSymbol = ""
	}

	out.EventMode = snap.EventMode
	out.EventMode.AffectedSymbols = filterRuntimeSymbols(snap.EventMode.AffectedSymbols, allowed)
	out.EventReactions = make([]EventReaction, 0, len(snap.EventReactions))
	for _, reaction := range snap.EventReactions {
		copyReaction := reaction
		copyReaction.Moves = filterRuntimeSymbolMap(reaction.Moves, allowed)
		copyReaction.Baseline = filterRuntimeSymbolMap(reaction.Baseline, allowed)
		out.EventReactions = append(out.EventReactions, copyReaction)
	}
	out.AdaptiveDataPolicy = snap.AdaptiveDataPolicy
	out.AdaptiveDataPolicy.HotSymbols = filterRuntimeSymbols(snap.AdaptiveDataPolicy.HotSymbols, allowed)

	out.Freshness = make([]FreshnessDiagnostic, 0, len(snap.Freshness))
	for _, diag := range snap.Freshness {
		copyDiag := diag
		if len(diag.Affected) > 0 {
			copyDiag.Affected = filterRuntimeSymbols(diag.Affected, allowed)
			if len(copyDiag.Affected) == 0 {
				continue
			}
		}
		out.Freshness = append(out.Freshness, copyDiag)
	}
	return out
}

func (e *Engine) SnapshotForUser(userID string) RuntimeSnapshot {
	if e == nil {
		return RuntimeSnapshot{}
	}
	snap := enforceHostedRuntimeRightsSnapshot(e.Snapshot())
	if e.app == nil || e.app.workspaces == nil {
		return snap
	}
	return e.app.runtimeSnapshotForUserFrom(userID, snap)
}

// broadcastRuntime keeps one canonical engine snapshot, then performs only a
// presentation/privacy filter per user. It never creates per-user provider or
// scoring pipelines.
func (a *Application) broadcastRuntime() {
	if a == nil || a.hub == nil || a.engine == nil {
		return
	}
	snap := enforceHostedRuntimeRightsSnapshot(a.engine.Snapshot())
	a.mu.RLock()
	ids := make([]string, 0, len(a.workspaces))
	for userID := range a.workspaces {
		ids = append(ids, userID)
	}
	workspaceMode := a.workspaces != nil
	a.mu.RUnlock()
	if !workspaceMode {
		a.hub.Broadcast(map[string]any{"type": "runtime", "runtime": snap})
		return
	}
	for _, userID := range ids {
		a.hub.BroadcastUser(userID, map[string]any{"type": "runtime", "runtime": a.runtimeSnapshotForUserFrom(userID, snap)})
	}
}

func (a *Application) workspaceBroadcastUsers() ([]string, bool) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	if a.workspaces == nil {
		return nil, false
	}
	ids := make([]string, 0, len(a.workspaces))
	for userID := range a.workspaces {
		ids = append(ids, userID)
	}
	return ids, true
}

func (a *Application) symbolAllowedForUser(userID, symbol string) bool {
	a.mu.RLock()
	allowed := a.runtimeAllowedSymbolsLocked(userID)
	a.mu.RUnlock()
	return runtimeSymbolAllowed(allowed, symbol)
}

func (a *Application) broadcastSymbolEvent(symbol string, payload any) {
	if !hostedQuoteServingAllowedForSymbol(a, symbol) {
		return
	}
	ids, workspaceMode := a.workspaceBroadcastUsers()
	if !workspaceMode {
		a.hub.Broadcast(payload)
		return
	}
	for _, userID := range ids {
		if a.symbolAllowedForUser(userID, symbol) {
			a.hub.BroadcastUser(userID, payload)
		}
	}
}

func filterNewsItemsForAllowed(items []NewsItem, allowed map[string]bool) []NewsItem {
	out := make([]NewsItem, 0, len(items))
	for _, item := range items {
		if len(item.Symbols) == 0 {
			out = append(out, item)
			continue
		}
		for _, symbol := range item.Symbols {
			if runtimeSymbolAllowed(allowed, symbol) {
				out = append(out, item)
				break
			}
		}
	}
	return out
}

func filterEarningsItemsForAllowed(items []EarningsItem, allowed map[string]bool) []EarningsItem {
	out := make([]EarningsItem, 0, len(items))
	for _, item := range items {
		if runtimeSymbolAllowed(allowed, item.Symbol) {
			out = append(out, item)
		}
	}
	return out
}

func filterFilingItemsForAllowed(items []FilingItem, allowed map[string]bool) []FilingItem {
	out := make([]FilingItem, 0, len(items))
	for _, item := range items {
		if runtimeSymbolAllowed(allowed, item.Symbol) {
			out = append(out, item)
		}
	}
	return out
}

func (a *Application) broadcastNews(items []NewsItem) {
	items = hostedNewsItemsForServing(items)
	ids, workspaceMode := a.workspaceBroadcastUsers()
	if !workspaceMode {
		a.hub.Broadcast(map[string]any{"type": "news", "news": items})
		return
	}
	for _, userID := range ids {
		a.mu.RLock()
		allowed := a.runtimeAllowedSymbolsLocked(userID)
		a.mu.RUnlock()
		a.hub.BroadcastUser(userID, map[string]any{"type": "news", "news": filterNewsItemsForAllowed(items, allowed)})
	}
}

func (a *Application) broadcastEarnings(items []EarningsItem) {
	items = hostedEarningsItemsForServing(items)
	ids, workspaceMode := a.workspaceBroadcastUsers()
	if !workspaceMode {
		a.hub.Broadcast(map[string]any{"type": "earnings", "earnings": items})
		return
	}
	for _, userID := range ids {
		a.mu.RLock()
		allowed := a.runtimeAllowedSymbolsLocked(userID)
		a.mu.RUnlock()
		a.hub.BroadcastUser(userID, map[string]any{"type": "earnings", "earnings": filterEarningsItemsForAllowed(items, allowed)})
	}
}

func (a *Application) broadcastFilings(items []FilingItem, intelligence map[string]SECIntelligenceSummary) {
	items = hostedFilingItemsForServing(items)
	intelligence = hostedSECIntelligenceForServing(intelligence)
	ids, workspaceMode := a.workspaceBroadcastUsers()
	if !workspaceMode {
		payload := map[string]any{"type": "filings", "filings": items}
		if intelligence != nil {
			payload["secIntelligence"] = intelligence
		}
		a.hub.Broadcast(payload)
		return
	}
	for _, userID := range ids {
		a.mu.RLock()
		allowed := a.runtimeAllowedSymbolsLocked(userID)
		a.mu.RUnlock()
		payload := map[string]any{"type": "filings", "filings": filterFilingItemsForAllowed(items, allowed)}
		if intelligence != nil {
			payload["secIntelligence"] = filterRuntimeSymbolMap(intelligence, allowed)
		}
		a.hub.BroadcastUser(userID, payload)
	}
}

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
	return nil
}

func hostedEarningsItemsForServing(items []EarningsItem) []EarningsItem {
	if !isHostedRuntime() {
		return items
	}
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
