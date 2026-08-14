package main

import (
	"context"
	"crypto/sha256"
	"encoding/json"
	"os"
	"sort"
	"strings"
	"time"
)

func NewEngine(app *Application) *Engine {
	bootstrapStarted := time.Now()
	e := &Engine{app: app, status: "stopped", mode: "demo", message: "Runtime is stopped", health: map[string]string{"quotes": "stopped", "history": "cached", "news": "stopped", "earnings": "stopped", "filings": "stopped", "fundamentals": "cached", "scanner": "idle", "cache-refresh": "ready", "pre-market-prep": "READY · next window pending", "market-open-prep": "READY · next window pending", "catalyst-watch": "READY · event driven", "provider-capabilities": "checking", "market-calendar": "checking", "market-activity": "checking", "corporate-actions": "checking", "symbol-intelligence": "checking", "ai": "ready"}, lastUpdated: map[string]int64{}, quotes: map[string]Quote{}, history: map[string][]HistoryPoint{}, bars: map[string]map[string][]Bar{}, fundamentals: map[string]FundamentalSnapshot{}, news: []NewsItem{}, earnings: []EarningsItem{}, filings: []FilingItem{}, secIntelligence: map[string]SECIntelligenceSummary{}, scanner: ScannerState{Mode: "day", Status: "idle", Message: "Run a scan to discover candidates.", Results: []ScannerResult{}}, globalDirect: map[string]GlobalDriver{}, macroMetrics: map[string]MacroMetric{}, macroEvents: []MacroEvent{}, eventReactions: []EventReaction{}, options: map[string]OptionsContext{}, capabilities: []CapabilityStatus{}, signalValidation: SignalValidationState{Snapshots: []SignalSnapshot{}, Message: "Collecting validated real-outcome snapshots."}, preparations: initialPreparationJobs(time.Now()), symbolIntelligence: map[string]SymbolIntelligence{}, catalystReactions: map[string]CatalystReactionState{}, marketOpenFlags: map[string][]string{}, marketOpenCheckpoint: MarketOpenCheckpoint{State: "NOT RUN YET"}, alpacaCalendar: map[string]AlpacaCalendarDay{}, marketActivity: MarketActivityState{}, corporateActions: []CorporateAction{}, manualActions: initialManualActions(), livePriorityHints: map[string]int64{}, lastBroadcast: map[string]time.Time{}, subscribedSymbols: map[string]bool{}, alpacaSubscribedSymbols: map[string]bool{}, providerCircuits: map[string]providerCircuit{}, providerCapabilityCircuits: map[string]providerCircuit{}, providerCapabilityStates: map[string]ProviderCapabilityStateRecord{}, smartRouterScorecard: SmartRouterScorecard{PolicyVersion: smartRouterPolicyVersion}, rapidMoveEvents: map[string]RapidMoveEvent{}, rapidMoveRecent: []RapidMoveEvent{}, rapidMoveScorecard: RapidMoveScorecard{PolicyVersion: rapidMovePolicyVersion}, providerQuotes: map[string]map[string]Quote{}, rawHistoryCoverage: map[string]RawHistoryCoverage{}, liquidityBaselines: map[string]LiquidityBaseline{}, workload: NewWorkloadController(), providerTelemetry: NewProviderTelemetry(), runtimeSLOTracker: NewRuntimeSLOTracker()}
	e.loadCache()
	cacheQuotes := len(e.quotes)
	persistedApplied := e.loadPersistedCanonicalQuotes()
	targets := 0
	if app != nil {
		app.mu.RLock()
		targets = len(masterSymbolsFromState(app.processingStateLocked()))
		app.mu.RUnlock()
	}
	e.startupProfile = RuntimeStartupDiagnostics{BootstrapDurationMs: time.Since(bootstrapStarted).Milliseconds(), CacheQuotesLoaded: cacheQuotes, PersistedQuotesApplied: persistedApplied, WarmStartQuotes: len(e.quotes), WarmStartTargetQuotes: targets}
	if targets > 0 {
		e.startupProfile.WarmStartCoveragePct = len(e.quotes) * 100 / targets
		if e.startupProfile.WarmStartCoveragePct > 100 {
			e.startupProfile.WarmStartCoveragePct = 100
		}
	}
	if app != nil && app.persistence != nil {
		e.initialPersistenceBytes = app.persistence.Diagnostics().Store.StorageBytes
	}
	e.sampleRuntimeLoad()
	return e
}
func (e *Engine) Snapshot() RuntimeSnapshot {
	e.app.mu.RLock()
	appState := e.app.processingStateLocked()
	settings := clone(appState.Settings)
	secrets := clone(e.app.secrets)
	overnightMode := settings.OvernightDataMode
	e.app.mu.RUnlock()
	e.mu.RLock()
	defer e.mu.RUnlock()
	subscribed := make([]string, 0, len(e.subscribedSymbols))
	for symbol := range e.subscribedSymbols {
		subscribed = append(subscribed, symbol)
	}
	sort.Strings(subscribed)
	alpacaSubscribed := make([]string, 0, len(e.alpacaSubscribedSymbols))
	for symbol := range e.alpacaSubscribedSymbols {
		alpacaSubscribed = append(alpacaSubscribed, symbol)
	}
	sort.Strings(alpacaSubscribed)
	session := marketSessionET(time.Now())
	now := time.Now().UnixMilli()
	recentFinnhub := e.lastTradeAt > 0 && now-e.lastTradeAt <= 30_000
	recentAlpacaStream := e.lastAlpacaStreamAt > 0 && now-e.lastAlpacaStreamAt <= 30_000
	recentAlpaca := e.lastAlpacaAt > 0 && now-e.lastAlpacaAt <= 30_000
	feedState := "stopped"
	if e.mode == "demo" && e.status == "running" {
		feedState = "demo"
	} else if session == "closed" || session == "weekend" {
		if e.status == "starting" || e.status == "running" || e.status == "degraded" {
			feedState = "connected-market-closed"
		}
	} else if session == "overnight" {
		if recentAlpaca && (e.lastAlpacaFeed == "overnight" || e.lastAlpacaFeed == "boats") {
			feedState = "overnight-streaming"
		} else if e.status == "starting" || e.status == "running" || e.status == "degraded" {
			feedState = "overnight-idle"
		}
	} else if recentAlpacaStream {
		feedState = "streaming"
	} else if recentAlpaca {
		feedState = "snapshot-live"
	} else if recentFinnhub {
		feedState = "finnhub-fallback"
	} else if e.alpacaWebSocketConnected || e.webSocketConnected {
		feedState = "connected-idle"
	} else if e.status == "starting" || e.status == "running" || e.status == "degraded" {
		feedState = "reconnecting"
	}
	boundaryAt, boundaryAction := marketSessionBoundaryET(time.Now())
	alloc := multiFeedAllocationWithHints(appState, clone(e.catalystReactions), clone(e.marketOpenFlags), clone(e.livePriorityHints), time.Now())
	feed := FeedDiagnostics{WebSocketConnected: e.webSocketConnected, SubscribedSymbols: subscribed, LastMessageAt: e.lastMessageAt, LastTradeAt: e.lastTradeAt, LastTradeSymbol: e.lastTradeSymbol, AlpacaWebSocketConnected: e.alpacaWebSocketConnected, AlpacaSubscribedSymbols: alpacaSubscribed, LastAlpacaStreamAt: e.lastAlpacaStreamAt, LastAlpacaStreamSymbol: e.lastAlpacaStreamSymbol, LastAlpacaAt: e.lastAlpacaAt, LastAlpacaSymbol: e.lastAlpacaSymbol, AlpacaLiveFeed: e.lastAlpacaFeed, OvernightDataMode: overnightMode, OvernightLiveAvailable: e.lastAlpacaFeed == "boats", MarketSession: session, SessionBoundaryAt: boundaryAt, SessionBoundaryAction: boundaryAction, FeedState: feedState, FinnhubMaxSymbols: finnhubPlanMaxSymbols, FinnhubReserveSlots: finnhubPlanMaxSymbols - finnhubActiveTarget, AlpacaMaxSymbols: alpacaPlanMaxSymbols, AlpacaReserveSlots: alpacaPlanMaxSymbols - alpacaActiveTarget, SnapshotSymbols: clone(alloc.Snapshot), TrackedSymbols: len(alloc.Finnhub) + len(alloc.Alpaca) + len(alloc.Snapshot), LiveSymbols: uniqueLiveSubscriptionCount(e.subscribedSymbols, e.alpacaSubscribedSymbols)}
	quotes := clone(e.quotes)
	metrics := clone(e.macroMetrics)
	direct := clone(e.globalDirect)
	global := deriveGlobalMarketContext(quotes, direct, metrics, settings.GlobalProviderMode)
	options := clone(e.options)
	events := updateEventLifecycles(clone(e.macroEvents), time.Now())
	eventMode := eventModeFor(events, time.Now(), settings.MacroEventModeEnabled)
	if eventMode.Active {
		for _, ev := range events {
			if ev.ID == eventMode.EventID {
				eventMode.AffectedSymbols, eventMode.AffectedSectors = affectedContext(ev, e.trackedSymbols())
				break
			}
		}
		eventMode.Prepared = true
		eventMode.QueuePrepared = true
	}
	validation := evaluateSignalSnapshotsWithActions(clone(e.signalValidation), clone(e.bars), clone(e.corporateActions), e.mode)
	validationLearning := buildValidationLearningSnapshot(appState, validation, clone(e.scanner), clone(e.bars), clone(e.macroEvents), clone(e.earnings), time.Now())
	capabilities := buildCapabilities(settings, secrets, e.mode, global, options, metrics, events)
	providerRouter := e.buildProviderRouterSnapshot(settings, secrets, quotes, clone(e.lastUpdated))

	quoteFreshnessSymbols := activeDeskSymbolsFromState(appState)
	if sym, ok := parseUserTicker(appState.UI.SelectedTicker); ok {
		quoteFreshnessSymbols = uniqueSymbols(append(quoteFreshnessSymbols, sym))
	}
	intradayFreshnessSymbols := []string{}
	if settings.DayEnabled {
		if wl, ok := watchlistValueByID(appState.Watchlists, settings.DayWatchlistID); ok {
			intradayFreshnessSymbols = userTradingSymbols(wl.Symbols)
		}
	}
	freshness, freshnessSummary := e.buildFreshnessDiagnostics(quotes, clone(e.lastUpdated), clone(e.health), quoteFreshnessSymbols, intradayFreshnessSymbols)
	providerReconciliation := buildProviderReconciliation(providerRouter, clone(e.providerQuotes), quotes, now)
	providerRouter.Scorecard.SourceDisagreements = providerReconciliationConflictCount(providerReconciliation)
	corporateTruth := buildCorporateActionTruth(clone(e.corporateActions), clone(e.bars), now, clone(e.rawHistoryCoverage))
	researchTruth := buildResearchPackageTruth(appState, quotes, clone(e.bars), clone(e.fundamentals), clone(e.news), clone(e.earnings), clone(e.filings), clone(e.lastUpdated), clone(e.health), providerReconciliation, time.Now(), ResearchPackageContext{CatalystReactions: clone(e.catalystReactions), Global: global})
	evidenceSnapshot := buildEvidenceSnapshot(researchTruth, providerReconciliation, corporateTruth)
	researchTruth.EvidenceSnapshotID = evidenceSnapshot.ID
	rawLiquidity := deriveLiquidityStatesWithContext(quotes, clone(e.bars), clone(e.liquidityBaselines), time.Now())
	marketIntelligence := buildMarketIntelligenceSnapshotWithContext(appState, quotes, clone(e.bars), rawLiquidity, global, MarketTradeabilityContext{EventMode: eventMode, Freshness: freshness, Scanner: clone(e.scanner), Options: options}, time.Now())
	rapidMove := e.buildRapidMoveStateLocked(time.Now())
	eventIntelligence := buildEventIntelligenceSnapshot(clone(e.news), events, eventMode, clone(e.eventReactions), clone(e.catalystReactions), clone(e.earnings), clone(e.lastUpdated), clone(e.health), time.Now(), SmartNotificationContext{Validation: validation, SEC: clone(e.secIntelligence), Freshness: freshness, ProviderRouter: providerRouter, Scanner: clone(e.scanner), RapidMove: rapidMove})
	alternativeIntelligence := buildContextAlternativeIntelligenceSnapshot(appState, quotes, options, metrics, marketIntelligence, clone(e.news), clone(e.filings), time.Now())
	adaptiveDataPolicy := buildAdaptiveDataPolicyState(clone(e.scanner), clone(e.health), time.Now())
	shadowControl := buildShadowControlState(clone(e.scanner), time.Now())
	preparations := preparationsWithMarketTradeability(clone(e.preparations), marketIntelligence.Tradeability, time.Now())
	runtimeLoad := clone(e.runtimeLoad)
	degradation := deriveRuntimeDegradation(e.status, e.mode, feed, freshness, providerRouter, runtimeLoad)
	if e.runtimeSLOTracker != nil {
		runtimeLoad.Recovery = e.runtimeSLOTracker.Observe(freshness, degradation, time.Now())
	}
	runtimeSLO := buildRuntimeSLOAssessmentWithContext(e.status, e.mode, feed, freshness, runtimeLoad, clone(e.scanner), appState, quotes, time.Now())
	return RuntimeSnapshot{Status: e.status, Mode: e.mode, Message: e.message, StartedAt: e.startedAt, LastError: e.lastError, Health: clone(e.health), LastUpdated: clone(e.lastUpdated), Quotes: quotes, History: clone(e.history), Bars: clone(e.bars), Fundamentals: clone(e.fundamentals), News: clone(e.news), Earnings: clone(e.earnings), Filings: clone(e.filings), SECIntelligence: clone(e.secIntelligence), Scanner: clone(e.scanner), Feed: feed, Global: global, MacroMetrics: metrics, MacroEvents: events, EventMode: eventMode, EventReactions: clone(e.eventReactions), Options: options, Capabilities: capabilities, SignalValidation: validation, ValidationLearning: validationLearning, Preparations: preparations, Liquidity: rawLiquidity, Intelligence: deriveIntelligenceStates(metrics), ProviderRegistry: enrichProviderCapabilityRegistry(buildProviderCapabilityRegistry(settings, secrets, e.health, e.symbolIntelligence, e.globalDirect), providerRouter), SymbolIntelligence: clone(e.symbolIntelligence), CatalystReactions: clone(e.catalystReactions), MarketOpenFlags: clone(e.marketOpenFlags), MarketOpenCheckpoint: clone(e.marketOpenCheckpoint), MarketActivity: clone(e.marketActivity), CorporateActions: clone(e.corporateActions), LiveCoverage: buildLiveCoverageStatesFrom(alloc, quotes, e.webSocketConnected, e.alpacaWebSocketConnected, clone(e.subscribedSymbols), clone(e.alpacaSubscribedSymbols)), ManualActions: clone(e.manualActions), ProviderRouter: providerRouter, RapidMove: rapidMove, ProviderReconciliation: providerReconciliation, ResearchPackage: researchTruth, EvidenceSnapshot: evidenceSnapshot, CorporateActionTruth: corporateTruth, Freshness: freshness, FreshnessSummary: freshnessSummary, MarketIntelligence: marketIntelligence, EventIntelligence: eventIntelligence, AlternativeIntelligence: alternativeIntelligence, AdaptiveDataPolicy: adaptiveDataPolicy, ShadowControl: shadowControl, RuntimeLoad: runtimeLoad, Degradation: degradation, RuntimeSLO: runtimeSLO}
}

func preparationsWithMarketTradeability(preparations map[string]PreparationJobStatus, tradeability MarketTradeabilityState, now time.Time) map[string]PreparationJobStatus {
	state := strings.TrimSpace(strings.ToUpper(tradeability.State))
	if state == "" {
		return preparations
	}
	updatedAt := tradeability.UpdatedAt
	if updatedAt == 0 {
		updatedAt = now.UnixMilli()
	}
	line := "Market Tradeability · " + state
	if detail := strings.TrimSpace(tradeability.Detail); detail != "" {
		line += " · " + detail
	}
	for _, key := range []string{"pre-market-prep", "market-open-prep"} {
		p, ok := preparations[key]
		if !ok {
			continue
		}
		p.Summary = appendUniqueString(p.Summary, line)
		if state == "WAIT" || state == "REDUCE SIZE" || state == "DATA DEGRADED" {
			severity := "medium"
			if state == "WAIT" || state == "DATA DEGRADED" {
				severity = "high"
			}
			p.Exceptions = appendUniqueCheckpointException(p.Exceptions, CheckpointException{Symbol: "Market", Reason: "Market Tradeability " + state, Severity: severity, Target: "market-intelligence", Source: "Market Intelligence", UpdatedAt: updatedAt})
		}
		preparations[key] = p
	}
	return preparations
}

func appendUniqueCheckpointException(values []CheckpointException, value CheckpointException) []CheckpointException {
	for _, existing := range values {
		if existing.Symbol == value.Symbol && existing.Reason == value.Reason && existing.Source == value.Source {
			return values
		}
	}
	return append(values, value)
}

func (e *Engine) loadCache() {
	data, err := os.ReadFile(e.app.cachePath())
	if err != nil {
		return
	}
	var c MarketCache
	if json.Unmarshal(data, &c) != nil {
		return
	}
	e.quotes = c.Quotes

	for symbol, q := range e.quotes {
		q.Symbol = normalizeSymbol(symbol)
		q.DataState = "cache"
		q.FeedType = "cache"
		e.quotes[symbol] = q
	}
	e.history = c.History
	e.bars = c.Bars
	e.fundamentals = c.Fundamentals
	e.news = c.News
	e.earnings = c.Earnings
	e.filings = c.Filings
	e.secIntelligence = c.SECIntelligence
	if len(c.Scanner.Results) > 0 || c.Scanner.UpdatedAt > 0 {
		e.scanner = c.Scanner
	}
	if c.GlobalDirect != nil {
		e.globalDirect = c.GlobalDirect
	}
	if c.MacroMetrics != nil {
		e.macroMetrics = c.MacroMetrics
	}
	if c.MacroEvents != nil {
		e.macroEvents = c.MacroEvents
	}
	if c.EventReactions != nil {
		e.eventReactions = c.EventReactions
	}
	if c.Options != nil {
		e.options = c.Options
	}
	if c.SignalValidation.Snapshots != nil {
		e.signalValidation = c.SignalValidation
	}
	if c.Preparations != nil {
		for k, v := range c.Preparations {
			e.preparations[k] = v
		}
	}
	if c.CatalystReactions != nil {
		e.catalystReactions = c.CatalystReactions
	}
	if c.MarketOpenFlags != nil {
		e.marketOpenFlags = c.MarketOpenFlags
	}
	if c.MarketOpenCheckpoint.RunAt > 0 || c.MarketOpenCheckpoint.State != "" {
		e.marketOpenCheckpoint = c.MarketOpenCheckpoint
	}
	if c.CorporateActions != nil {
		e.corporateActions = mergeCorporateActionLedger(nil, c.CorporateActions, c.SavedAt)
	}
	if c.RawHistoryCoverage != nil {
		e.rawHistoryCoverage = c.RawHistoryCoverage
	}
	if c.LiquidityBaselines != nil {
		e.liquidityBaselines = c.LiquidityBaselines
	}
	if c.ProviderCapabilityStates != nil {
		e.providerCapabilityStates = c.ProviderCapabilityStates
	}
	if c.RapidMoveEvents != nil {
		e.rapidMoveEvents = c.RapidMoveEvents
	}
	if c.RapidMoveRecent != nil {
		e.rapidMoveRecent = c.RapidMoveRecent
	}
	if c.RapidMoveScorecard.PolicyVersion != "" || c.RapidMoveScorecard.Observations > 0 {
		e.rapidMoveScorecard = c.RapidMoveScorecard
	}
	for k, v := range c.LastUpdated {
		if v > e.lastUpdated[k] {
			e.lastUpdated[k] = v
		}
	}
	if e.quotes == nil {
		e.quotes = map[string]Quote{}
	}
	if e.history == nil {
		e.history = map[string][]HistoryPoint{}
	}
	if e.bars == nil {
		e.bars = map[string]map[string][]Bar{}
	}
	if e.fundamentals == nil {
		e.fundamentals = map[string]FundamentalSnapshot{}
	}
	if e.secIntelligence == nil {
		e.secIntelligence = map[string]SECIntelligenceSummary{}
	}
	if e.liquidityBaselines == nil {
		e.liquidityBaselines = map[string]LiquidityBaseline{}
	}
	e.lastUpdated["cache"] = c.SavedAt
}
func (e *Engine) loadPersistedCanonicalQuotes() int {
	if e.app == nil || e.app.persistence == nil {
		return 0
	}
	persisted := e.app.persistence.LoadQuotes()
	if len(persisted) == 0 {
		return 0
	}
	applied := 0
	for symbol, candidate := range persisted {
		symbol = normalizeSymbol(symbol)
		current, ok := e.quotes[symbol]
		if ok && maxInt64(current.UpdatedAt, current.ProviderTimestamp) >= maxInt64(candidate.UpdatedAt, candidate.ProviderTimestamp) {
			continue
		}
		candidate.Symbol = symbol
		// Persisted labels are never trusted as timeless freshness truth. The
		// normal freshness pipeline re-derives CURRENT/STALE/DELAYED/DEGRADED.
		candidate.DataState = "persisted"
		candidate.FeedType = "persisted"
		e.quotes[symbol] = candidate
		applied++
	}
	return applied
}

func (e *Engine) saveCache() error {
	e.mu.RLock()
	lastUpdated := clone(e.lastUpdated)
	delete(lastUpdated, "cache")
	base := MarketCache{Quotes: clone(e.quotes), History: clone(e.history), Bars: clone(e.bars), Fundamentals: clone(e.fundamentals), News: clone(e.news), Earnings: clone(e.earnings), Filings: clone(e.filings), SECIntelligence: clone(e.secIntelligence), Scanner: clone(e.scanner), GlobalDirect: clone(e.globalDirect), MacroMetrics: clone(e.macroMetrics), MacroEvents: clone(e.macroEvents), EventReactions: clone(e.eventReactions), Options: clone(e.options), SignalValidation: clone(e.signalValidation), Preparations: clone(e.preparations), CatalystReactions: clone(e.catalystReactions), MarketOpenFlags: clone(e.marketOpenFlags), MarketOpenCheckpoint: clone(e.marketOpenCheckpoint), CorporateActions: clone(e.corporateActions), RawHistoryCoverage: clone(e.rawHistoryCoverage), LiquidityBaselines: clone(e.liquidityBaselines), ProviderCapabilityStates: clone(e.providerCapabilityStates), RapidMoveEvents: clone(e.rapidMoveEvents), RapidMoveRecent: clone(e.rapidMoveRecent), RapidMoveScorecard: clone(e.rapidMoveScorecard), LastUpdated: lastUpdated, SavedAt: 0}
	e.mu.RUnlock()
	fingerprintData, err := json.Marshal(base)
	if err != nil {
		return err
	}
	hash := sha256.Sum256(fingerprintData)
	e.mu.RLock()
	unchanged := e.lastCacheHashSet && e.lastCacheHash == hash
	e.mu.RUnlock()
	if unchanged {
		return nil
	}
	base.SavedAt = time.Now().UnixMilli()
	data, err := json.Marshal(base)
	if err != nil {
		return err
	}
	if err := atomicWrite(e.app.cachePath(), data, 0600); err != nil {
		return err
	}
	if e.app.persistence != nil {
		e.app.persistence.EnqueueQuotes(base.Quotes)
	}
	e.mu.Lock()
	e.lastCacheHash = hash
	e.lastCacheHashSet = true
	e.lastUpdated["cache"] = base.SavedAt
	e.mu.Unlock()
	return nil
}

func (e *Engine) setStatus(status, message string) {
	e.mu.Lock()
	e.status = status
	e.message = message
	e.mu.Unlock()
	e.app.broadcastRuntime()
}

func (e *Engine) setHealth(key, value string) {
	e.mu.Lock()
	e.health[key] = value

	e.lastUpdated["health:"+key] = time.Now().UnixMilli()
	e.mu.Unlock()
	e.app.broadcastRuntime()
}

func (e *Engine) setError(prefix string, err error) {
	e.mu.Lock()
	e.lastError = prefix + ": " + err.Error()
	e.mu.Unlock()
	e.app.broadcastRuntime()
}

func (e *Engine) trackedSymbols() []string {
	e.app.mu.RLock()
	defer e.app.mu.RUnlock()
	unique := requiredSymbolsFromState(e.app.processingStateLocked())
	return unique[:minInt(100, len(unique))]
}

// onSymbolSetChanged makes watchlist edits immediately useful instead of
// waiting for the next periodic provider cycle. Demo symbols are hydrated
// synchronously so the UI can render a complete row as soon as Add returns.
// Live mode updates the WebSocket subscription immediately and performs a
// lightweight provider refresh plus required historical-bar hydration in the background.
func (e *Engine) hydrateFocusQuote(symbol string) {
	symbol = normalizeSymbol(symbol)
	if symbol == "" {
		return
	}
	e.mu.RLock()
	status, mode, ws, alpacaWS := e.status, e.mode, e.ws, e.alpacaWS
	e.mu.RUnlock()
	if status != "running" && status != "degraded" {
		return
	}
	if mode == "demo" {
		e.ensureDemoSymbol(symbol)
		return
	}
	if ws != nil {
		e.syncLiveSubscriptions(ws)
	}
	if alpacaWS != nil {
		e.syncAlpacaSubscriptions(alpacaWS)
	}
	e.app.mu.RLock()
	finnhubKey := strings.TrimSpace(e.app.secrets.Finnhub)
	alpacaKey := strings.TrimSpace(e.app.secrets.AlpacaKey)
	alpacaSecret := strings.TrimSpace(e.app.secrets.AlpacaSecret)
	e.app.mu.RUnlock()
	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
		defer cancel()
		if finnhubKey != "" {
			e.refreshSingleFinnhubSnapshot(ctx, finnhubKey, symbol)
		}
		if alpacaKey != "" && alpacaSecret != "" {
			e.refreshSingleAlpacaSnapshot(ctx, alpacaKey, alpacaSecret, symbol)
		}
	}()
}

func (e *Engine) onSymbolSetChanged(symbol string) {
	symbol = normalizeSymbol(symbol)
	if symbol == "" {
		return
	}
	e.hydrateFocusQuote(symbol)

	e.requestHistoryHydration(symbol)
}
