package main

import (
	"fmt"
	"strings"
	"testing"
	"time"
)

func newV1801Engine(t *testing.T) *Engine {
	t.Helper()
	app := newTestApplication(t)
	return app.engine
}

func TestV1801SmartRouterClassifiesAndSuppressesNotEntitled(t *testing.T) {
	state, ttl, reason := classifyProviderCapabilityFailure(fmt.Errorf("Twelve Data HTTP 404: available starting with Pro plan"))
	if state != providerCapabilityNotEntitled || ttl < 6*time.Hour || !strings.Contains(strings.ToLower(reason), "pro") {
		t.Fatalf("unexpected classification: state=%s ttl=%s reason=%q", state, ttl, reason)
	}
	e := newV1801Engine(t)
	e.recordProviderCapabilityFailure("VIX / Indices", "Twelve Data", fmt.Errorf("HTTP 404: available starting with Pro plan"))
	if ok, got := e.providerCapabilityAllowed("VIX / Indices", "Twelve Data", time.Now()); ok || got != providerCapabilityNotEntitled {
		t.Fatalf("expected NOT_ENTITLED suppression, ok=%v state=%q", ok, got)
	}
	e.mu.RLock()
	avoided := e.providerCallsAvoided
	scoreAvoided := e.smartRouterScorecard.EntitlementAvoided
	e.mu.RUnlock()
	if avoided < 1 || scoreAvoided < 1 {
		t.Fatalf("suppression must record calls avoided: total=%d scorecard=%d", avoided, scoreAvoided)
	}
}

func TestV1801SmartRouterCapabilityCircuitIsolation(t *testing.T) {
	e := newV1801Engine(t)
	for i := 0; i < 3; i++ {
		e.recordProviderCapabilityCircuitFailure("News", "Finnhub", fmt.Errorf("temporary news endpoint failure"))
	}
	e.mu.RLock()
	news := e.capabilityCircuitStatusLocked("Finnhub", "News", time.Now().UnixMilli())
	live := e.capabilityCircuitStatusLocked("Finnhub", "US Live Equities", time.Now().UnixMilli())
	e.mu.RUnlock()
	if news != "OPEN" {
		t.Fatalf("expected News circuit OPEN, got %s", news)
	}
	if live != "CLOSED" {
		t.Fatalf("News failure must not open live-equities capability: %s", live)
	}
}

func TestV1801SmartRouterDynamicPreferredLiveProvider(t *testing.T) {
	e := newV1801Engine(t)
	et := easternLocation()
	regular := time.Date(2026, 8, 13, 14, 0, 0, 0, et)
	ranked := e.rankedProviderRoute("US Live Equities", WorkTierUserActionable, Settings{}, Secrets{
		Finnhub: "f", AlpacaKey: "a", AlpacaSecret: "s", TwelveData: "t",
	}, regular)
	if len(ranked) < 2 {
		t.Fatalf("expected ranked providers, got %+v", ranked)
	}
	if ranked[0].Provider != "Finnhub" || !ranked[0].Eligible {
		t.Fatalf("expected Finnhub preferred for live stream when suitable, got %+v", ranked)
	}
}

func TestV1801ProviderLatencyPercentiles(t *testing.T) {
	p50, p95 := latencyPercentiles([]int64{10, 20, 30, 40, 50, 60, 70, 80, 90, 100})
	if p50 != 50 || p95 != 90 {
		t.Fatalf("unexpected percentiles p50=%d p95=%d", p50, p95)
	}
}

func seedRapidMoveAt(t *testing.T, e *Engine, symbol string, startPrice, currentPrice float64, sourceAgreement string, catalyst bool, at time.Time) Quote {
	t.Helper()
	now := at.UnixMilli()
	e.mu.Lock()
	e.history[symbol] = []HistoryPoint{{T: now - 120_000, P: startPrice}, {T: now - 60_000, P: startPrice}}
	q := Quote{Symbol: symbol, Price: currentPrice, Bid: currentPrice * .999, Ask: currentPrice * 1.001, ProviderTimestamp: now, UpdatedAt: now, Source: "finnhub-websocket", FeedType: "websocket", DataState: "live"}
	switch sourceAgreement {
	case "AGREED":
		e.providerQuotes[symbol] = map[string]Quote{
			"Finnhub": q,
			"Alpaca":  {Symbol: symbol, Price: currentPrice * 1.0005, ProviderTimestamp: now, UpdatedAt: now, Source: "alpaca-iex", FeedType: "websocket", DataState: "live"},
		}
	case "CONFLICT":
		e.providerQuotes[symbol] = map[string]Quote{
			"Finnhub": q,
			"Alpaca":  {Symbol: symbol, Price: currentPrice * .95, ProviderTimestamp: now, UpdatedAt: now, Source: "alpaca-iex", FeedType: "websocket", DataState: "live"},
		}
	default:
		e.providerQuotes[symbol] = map[string]Quote{"Finnhub": q}
	}
	if catalyst {
		e.catalystReactions[symbol] = CatalystReactionState{Symbol: symbol, TriggerType: "GUIDANCE", Trigger: "Guidance raised", State: "TRIGGERED", TriggerAt: now - 20_000, UpdatedAt: now}
	}
	e.mu.Unlock()
	return q
}

func seedRapidMove(t *testing.T, e *Engine, symbol string, startPrice, currentPrice float64, sourceAgreement string, catalyst bool) Quote {
	t.Helper()
	return seedRapidMoveAt(t, e, symbol, startPrice, currentPrice, sourceAgreement, catalyst, time.Now())
}

func TestV1801RapidMoveFivePercentMinuteAlertsAndPromotes(t *testing.T) {
	e := newV1801Engine(t)
	q := seedRapidMove(t, e, "NVDA", 100, 106, "AGREED", true)
	e.evaluateRapidMoveObservation("NVDA", q)

	e.mu.RLock()
	ev, ok := e.rapidMoveEvents["NVDA"]
	promotions := append([]OpportunityPromotion(nil), e.scanner.Radar.Promotions...)
	_, liveHint := e.livePriorityHints["NVDA"]
	e.mu.RUnlock()
	if !ok || !ev.Alerted {
		t.Fatalf("expected material rapid-move alert, event=%+v", ev)
	}
	if ev.State != "CONFIRMED" && ev.State != "EXTENDED" {
		t.Fatalf("expected catalyst-confirmed/extended state, got %s", ev.State)
	}
	if ev.CatalystState != "CONFIRMED" || ev.SourceAgreement != "AGREED" {
		t.Fatalf("expected correlated evidence, event=%+v", ev)
	}
	found := false
	for _, p := range promotions {
		if p.Symbol == "NVDA" && p.State == "RAPID MOVE" {
			found = true
		}
	}
	if !found || !liveHint {
		t.Fatalf("expected Radar promotion and live priority, promotions=%+v live=%v", promotions, liveHint)
	}
	state := e.buildRapidMoveStateLocked(time.Now())
	notifications := rapidMoveNotifications(state, time.Now())
	if len(notifications) == 0 || notifications[0].Symbol != "NVDA" {
		t.Fatalf("expected prominent smart notification, got %+v", notifications)
	}
}

func TestV1801RapidMoveCanAlertBeforeCatalystArrives(t *testing.T) {
	e := newV1801Engine(t)
	q := seedRapidMove(t, e, "AMD", 100, 107, "AGREED", false)
	e.evaluateRapidMoveObservation("AMD", q)
	e.mu.RLock()
	ev := e.rapidMoveEvents["AMD"]
	e.mu.RUnlock()
	if !ev.Alerted {
		t.Fatalf("strong validated price event should be able to alert before news arrives: %+v", ev)
	}
	if ev.CatalystState != "UNEXPLAINED" || ev.State != "VALIDATING" && ev.State != "EXTENDED" {
		t.Fatalf("expected unexplained/validating event, got %+v", ev)
	}
}

func TestV1801RapidMoveSuppressesMechanicalAndSourceConflictNoise(t *testing.T) {
	t.Run("split", func(t *testing.T) {
		e := newV1801Engine(t)
		q := seedRapidMove(t, e, "AAPL", 100, 107, "AGREED", true)
		day := time.Now().In(easternLocation()).Format("2006-01-02")
		e.mu.Lock()
		e.corporateActions = []CorporateAction{{Symbol: "AAPL", Type: "stock_split", ExDate: day, Status: "EFFECTIVE", UpdatedAt: time.Now().UnixMilli()}}
		e.mu.Unlock()
		e.evaluateRapidMoveObservation("AAPL", q)
		e.mu.RLock()
		ev := e.rapidMoveEvents["AAPL"]
		e.mu.RUnlock()
		if ev.Alerted || ev.MechanicalRisk == "" {
			t.Fatalf("split must suppress alert: %+v", ev)
		}
	})
	t.Run("source-conflict", func(t *testing.T) {
		e := newV1801Engine(t)
		q := seedRapidMove(t, e, "MSFT", 100, 107, "CONFLICT", true)
		e.evaluateRapidMoveObservation("MSFT", q)
		e.mu.RLock()
		ev := e.rapidMoveEvents["MSFT"]
		e.mu.RUnlock()
		if ev.Alerted || ev.SourceAgreement != "CONFLICT" {
			t.Fatalf("source conflict must suppress production alert: %+v", ev)
		}
	})
}

func TestV1801RapidMoveSuppressesLowPriceNoise(t *testing.T) {
	e := newV1801Engine(t)
	q := seedRapidMove(t, e, "ABCD", 1.00, 1.10, "AGREED", true)
	e.evaluateRapidMoveObservation("ABCD", q)
	e.mu.RLock()
	ev := e.rapidMoveEvents["ABCD"]
	suppressed := e.rapidMoveScorecard.LiquiditySuppressed
	e.mu.RUnlock()
	if ev.Alerted || suppressed == 0 {
		t.Fatalf("low-price burst must remain suppressed: %+v scorecard=%+v", ev, e.rapidMoveScorecard)
	}
}

func TestV1801RapidMoveDedupesOneMarketEvent(t *testing.T) {
	e := newV1801Engine(t)
	q := seedRapidMove(t, e, "META", 100, 106, "AGREED", true)
	e.evaluateRapidMoveObservation("META", q)
	e.evaluateRapidMoveObservation("META", q)
	e.mu.RLock()
	alerts := e.rapidMoveScorecard.ProductionAlerts
	dup := e.rapidMoveScorecard.DuplicateUpdates
	e.mu.RUnlock()
	if alerts != 1 || dup < 1 {
		t.Fatalf("expected one alert with duplicate update, alerts=%d dup=%d", alerts, dup)
	}
}

func TestV1801RapidMoveOutcomeLearningResolves(t *testing.T) {
	e := newV1801Engine(t)
	now := time.Now().UnixMilli()
	e.mu.Lock()
	e.rapidMoveEvents["TSLA"] = RapidMoveEvent{ID: "rapid-tsla-test", TraceID: "rapid-tsla-test", Symbol: "TSLA", Direction: "UP", State: "CONFIRMED", Alerted: true, PolicyVersion: rapidMovePolicyVersion, DetectedAt: now - 21*60_000, UpdatedAt: now - 21*60_000, Price: 100, StartPrice: 95}
	e.history["TSLA"] = []HistoryPoint{{T: now - 120_000, P: 103}, {T: now - 60_000, P: 103}}
	e.mu.Unlock()
	q := Quote{Symbol: "TSLA", Price: 103, ProviderTimestamp: now, UpdatedAt: now, Source: "finnhub-websocket"}
	e.evaluateRapidMoveObservation("TSLA", q)
	e.mu.RLock()
	ev := e.rapidMoveEvents["TSLA"]
	resolved := e.rapidMoveScorecard.OutcomesResolved
	recent := len(e.rapidMoveRecent)
	e.mu.RUnlock()
	if ev.State != "RESOLVED" || ev.Outcome20mPct == nil || resolved != 1 || recent != 1 {
		t.Fatalf("expected durable outcome resolution, event=%+v resolved=%d recent=%d", ev, resolved, recent)
	}
}

func TestV1801RapidMoveCoverageTruthIsExplicitlyPartial(t *testing.T) {
	e := newV1801Engine(t)
	e.mu.Lock()
	e.subscribedSymbols["SPY"] = true
	e.subscribedSymbols["NVDA"] = true
	e.mu.Unlock()
	e.mu.Lock()
	c := e.rapidMoveCoverageLocked(time.Now())
	e.mu.Unlock()
	if c.State != "TIERED_PARTIAL" || c.LiveSymbols != 2 || !strings.Contains(c.Detail, "not claimed") {
		t.Fatalf("coverage truth must not imply full-market surveillance: %+v", c)
	}
}

func TestV1801SmartRouterNotConfiguredIsNotEntitlementFailure(t *testing.T) {
	e := newV1801Engine(t)
	ranked := e.rankedProviderRoute("US Live Equities", WorkTierUserActionable, Settings{}, Secrets{}, time.Now())
	if len(ranked) == 0 {
		t.Fatal("expected route candidates")
	}
	for _, row := range ranked {
		if row.Eligible {
			t.Fatalf("unconfigured provider unexpectedly eligible: %+v", row)
		}
		if row.State != providerCapabilityNotConfigured {
			t.Fatalf("unconfigured provider must remain NOT_CONFIGURED, not entitlement truth: %+v", row)
		}
	}
}

func TestV1801SmartRouterAvoidedTelemetrySeparatesReasons(t *testing.T) {
	t.Run("circuit", func(t *testing.T) {
		e := newV1801Engine(t)
		e.recordProviderCapabilityFailure("News", "Finnhub", fmt.Errorf("HTTP 429 rate limit"))
		if ok, state := e.providerCapabilityAllowed("News", "Finnhub", time.Now()); ok || state != providerCapabilityRateLimited {
			t.Fatalf("expected rate-limit suppression, ok=%v state=%s", ok, state)
		}
		e.mu.RLock()
		s := e.smartRouterScorecard
		e.mu.RUnlock()
		if s.CircuitAvoided != 1 || s.EntitlementAvoided != 0 || s.CapacityAvoided != 0 {
			t.Fatalf("wrong avoided-reason attribution: %+v", s)
		}
	})
	t.Run("capacity", func(t *testing.T) {
		e := newV1801Engine(t)
		e.recordProviderCapabilityFailure("US Live Equities", "Finnhub", fmt.Errorf("subscription capacity saturated"))
		if ok, state := e.providerCapabilityAllowed("US Live Equities", "Finnhub", time.Now()); ok || state != providerCapabilitySaturated {
			t.Fatalf("expected capacity suppression, ok=%v state=%s", ok, state)
		}
		e.mu.RLock()
		s := e.smartRouterScorecard
		e.mu.RUnlock()
		if s.CapacityAvoided != 1 || s.EntitlementAvoided != 0 || s.CircuitAvoided != 0 {
			t.Fatalf("wrong avoided-reason attribution: %+v", s)
		}
	})
}

func TestV1801RapidMoveExcludesLaggingProviderFromConflict(t *testing.T) {
	e := newV1801Engine(t)
	q := seedRapidMove(t, e, "AVGO", 100, 107, "AGREED", true)
	e.mu.Lock()
	lagging := e.providerQuotes["AVGO"]["Alpaca"]
	lagging.Price = 100
	lagging.ProviderTimestamp = q.ProviderTimestamp - 30_000
	lagging.UpdatedAt = q.UpdatedAt - 30_000
	e.providerQuotes["AVGO"]["Alpaca"] = lagging
	e.mu.Unlock()
	e.evaluateRapidMoveObservation("AVGO", q)
	e.mu.RLock()
	ev := e.rapidMoveEvents["AVGO"]
	e.mu.RUnlock()
	if ev.SourceAgreement != "SINGLE_SOURCE" || !ev.Alerted {
		t.Fatalf("lagging secondary source must be excluded rather than manufacture a conflict: %+v", ev)
	}
	if !strings.Contains(strings.ToLower(ev.SourceAgreementDetail), "lagging") {
		t.Fatalf("expected truthful lagging-source explanation: %q", ev.SourceAgreementDetail)
	}
}

func TestV1801RapidMoveOutcomeUsesProviderEventTimeAnchor(t *testing.T) {
	e := newV1801Engine(t)
	now := time.Now().UnixMilli()
	e.mu.Lock()
	e.rapidMoveEvents["CRM"] = RapidMoveEvent{
		ID: "rapid-crm-anchor", TraceID: "rapid-crm-anchor", Symbol: "CRM", Direction: "UP", State: "CONFIRMED", Alerted: true,
		PolicyVersion: rapidMovePolicyVersion, EventProviderAt: now - 21*60_000, DetectedAt: now - 60_000, UpdatedAt: now - 21*60_000,
		Price: 100, StartPrice: 95,
	}
	e.history["CRM"] = []HistoryPoint{{T: now - 120_000, P: 103}, {T: now - 60_000, P: 103}}
	e.mu.Unlock()
	q := Quote{Symbol: "CRM", Price: 103, ProviderTimestamp: now, UpdatedAt: now, Source: "finnhub-websocket"}
	e.evaluateRapidMoveObservation("CRM", q)
	e.mu.RLock()
	ev := e.rapidMoveEvents["CRM"]
	e.mu.RUnlock()
	if ev.Outcome20mPct == nil || ev.State != "RESOLVED" {
		t.Fatalf("outcome horizon must use provider/event time, not delayed detection wall clock: %+v", ev)
	}
}

func TestV1801RapidMoveCatalystRejectsFutureNewsInPointInTimeReplay(t *testing.T) {
	now := time.Now()
	news := []NewsItem{{Headline: "Guidance raised materially", Source: "TestWire", Datetime: now.Add(60 * time.Second).Unix(), Symbols: []string{"NVDA"}}}
	state, _, _, _ := rapidMoveCatalyst("NVDA", news, nil, nil, nil, now)
	if state != "UNEXPLAINED" {
		t.Fatalf("future headline leaked into point-in-time catalyst evidence: %s", state)
	}
}

func TestV1801RapidMoveCoverageSeparatesSubscriptionFromWindowReadiness(t *testing.T) {
	e := newV1801Engine(t)
	now := time.Now().UnixMilli()
	e.mu.Lock()
	e.subscribedSymbols["SPY"] = true
	e.subscribedSymbols["NVDA"] = true
	e.quotes["NVDA"] = Quote{Symbol: "NVDA", Price: 100, ProviderTimestamp: now, UpdatedAt: now, Source: "finnhub-websocket"}
	e.history["NVDA"] = []HistoryPoint{{T: now - 60_000, P: 98}, {T: now, P: 100}}
	c := e.rapidMoveCoverageLocked(time.UnixMilli(now))
	e.mu.Unlock()
	if c.SubscribedSymbols != 2 || c.LiveSymbols != 2 || c.WindowReadySymbols != 1 {
		t.Fatalf("coverage must distinguish subscribed from actual window-ready symbols: %+v", c)
	}
}

func TestV1801RapidMoveEscalationFromShadowToProductionCountsFirstAlert(t *testing.T) {
	e := newV1801Engine(t)
	regular := time.Date(2026, 8, 13, 14, 0, 0, 0, easternLocation())
	q := seedRapidMoveAt(t, e, "ORCL", 100, 104.7, "AGREED", false, regular)
	e.evaluateRapidMoveObservation("ORCL", q)
	e.mu.RLock()
	before := e.rapidMoveEvents["ORCL"]
	beforeAlerts := e.rapidMoveScorecard.ProductionAlerts
	e.mu.RUnlock()
	if before.Alerted || !before.ShadowWouldAlert || beforeAlerts != 0 {
		t.Fatalf("expected shadow-only precursor before production escalation: event=%+v alerts=%d", before, beforeAlerts)
	}

	q.Price = 106
	q.Bid = 105.9
	q.Ask = 106.1
	q.ProviderTimestamp += 2_000
	q.UpdatedAt = q.ProviderTimestamp
	e.mu.Lock()
	e.providerQuotes["ORCL"] = map[string]Quote{
		"Finnhub": q,
		"Alpaca":  {Symbol: "ORCL", Price: 106.02, ProviderTimestamp: q.ProviderTimestamp, UpdatedAt: q.UpdatedAt, Source: "alpaca-iex", FeedType: "websocket", DataState: "live"},
	}
	e.mu.Unlock()
	e.evaluateRapidMoveObservation("ORCL", q)

	e.mu.RLock()
	after := e.rapidMoveEvents["ORCL"]
	afterAlerts := e.rapidMoveScorecard.ProductionAlerts
	e.mu.RUnlock()
	if !after.Alerted || afterAlerts != 1 {
		t.Fatalf("shadow-to-production transition must be the first production alert: event=%+v alerts=%d", after, afterAlerts)
	}
	if after.TraceID != before.TraceID {
		t.Fatalf("material escalation must preserve one intelligence trace: before=%s after=%s", before.TraceID, after.TraceID)
	}
}

func TestV1801RapidMoveTransitionsToFadingBeforeResolution(t *testing.T) {
	e := newV1801Engine(t)
	now := time.Now().UnixMilli()
	ev := RapidMoveEvent{
		ID: "rapid-nvda-fade", TraceID: "rapid-nvda-fade", Symbol: "NVDA", Direction: "UP",
		State: "CONFIRMED", Alerted: true, PolicyVersion: rapidMovePolicyVersion,
		EventProviderAt: now - 90_000, DetectedAt: now - 90_000, UpdatedAt: now - 90_000,
		StartPrice: 100, Price: 106,
	}
	q := Quote{Symbol: "NVDA", Price: 102, ProviderTimestamp: now, UpdatedAt: now, Source: "finnhub-websocket"}
	e.mu.Lock()
	updated, changed := e.updateRapidMoveOutcomeLocked(ev, q, now)
	e.mu.Unlock()
	if !changed || updated.State != "FADING" {
		t.Fatalf("expected material giveback to transition to FADING before 20m resolution: %+v", updated)
	}
	if updated.Outcome20mPct != nil {
		t.Fatalf("FADING is an intermediate state and must not fabricate resolved outcome: %+v", updated)
	}
}
