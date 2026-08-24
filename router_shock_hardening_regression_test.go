package main

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestV1806ProviderRouterScorecardSurfacesCanonicalSourceDisagreement(t *testing.T) {
	e := newV1801Engine(t)
	now := time.Now().UnixMilli()
	e.mu.Lock()
	e.quotes["NVDA"] = Quote{Symbol: "NVDA", Price: 100, ProviderTimestamp: now, UpdatedAt: now, Source: "finnhub-websocket", FeedType: "websocket", DataState: "live"}
	e.providerQuotes["NVDA"] = map[string]Quote{
		"Finnhub": {Symbol: "NVDA", Price: 100, ProviderTimestamp: now, UpdatedAt: now, Source: "finnhub-websocket", FeedType: "websocket", DataState: "live"},
		"Alpaca":  {Symbol: "NVDA", Price: 101, ProviderTimestamp: now, UpdatedAt: now, Source: "alpaca-iex", FeedType: "websocket", DataState: "live"},
	}
	e.mu.Unlock()
	snap := e.Snapshot()
	if snap.ProviderRouter.Scorecard.SourceDisagreements != 1 {
		t.Fatalf("router scorecard must report canonical reconciliation conflict truth, got %+v", snap.ProviderRouter.Scorecard)
	}
	found := false
	for _, row := range snap.ProviderReconciliation {
		if row.Symbol == "NVDA" && row.State == "CONFLICT" {
			found = true
		}
	}
	if !found {
		t.Fatalf("expected NVDA provider conflict in canonical reconciliation: %+v", snap.ProviderReconciliation)
	}
}

func TestV1806RapidMoveClassifiesMarketWideEventAsMarketShock(t *testing.T) {
	e := newV1801Engine(t)
	at := time.Date(2026, 8, 13, 14, 0, 0, 0, easternLocation())
	q := seedRapidMoveAt(t, e, "NVDA", 100, 106, "AGREED", false, at)
	now := at.UnixMilli()
	e.mu.Lock()
	e.history["SPY"] = []HistoryPoint{{T: now - 60_000, P: 100}, {T: now, P: 101.2}}
	e.mu.Unlock()
	e.evaluateRapidMoveObservation("NVDA", q)

	e.mu.RLock()
	ev := e.rapidMoveEvents["NVDA"]
	promotions := append([]OpportunityPromotion(nil), e.scanner.Radar.Promotions...)
	shocks := e.rapidMoveScorecard.MarketShockAlerts
	e.mu.RUnlock()
	if !ev.Alerted || ev.Classification != "MARKET_SHOCK" || !strings.HasPrefix(ev.MarketContext, "MARKET_WIDE") {
		t.Fatalf("market-wide rapid move must be explicitly classified as MARKET_SHOCK: %+v", ev)
	}
	if shocks != 1 {
		t.Fatalf("market shock scorecard must count first production alert exactly once: %+v", e.rapidMoveScorecard)
	}
	found := false
	for _, p := range promotions {
		if p.Symbol == "NVDA" && p.State == "MARKET SHOCK" {
			found = true
		}
	}
	if !found {
		t.Fatalf("market shock must retain canonical Radar promotion with market-wide semantics: %+v", promotions)
	}
}

func TestV1806RapidMoveHysteresisPreventsAlertStateThrash(t *testing.T) {
	e := newV1801Engine(t)
	at := time.Date(2026, 8, 13, 14, 0, 0, 0, easternLocation())
	q := seedRapidMoveAt(t, e, "AMD", 100, 106, "AGREED", true, at)
	e.evaluateRapidMoveObservation("AMD", q)
	e.mu.RLock()
	first := e.rapidMoveEvents["AMD"]
	e.mu.RUnlock()
	if !first.Alerted || (first.State != "CONFIRMED" && first.State != "EXTENDED") {
		t.Fatalf("expected initial production event: %+v", first)
	}

	// Small retracement remains above the FADING retention threshold but below
	// the production trigger. It must not regress the already-alerted event to EARLY.
	q.Price = 104.5
	q.Bid = 104.4
	q.Ask = 104.6
	q.ProviderTimestamp += 2_000
	q.UpdatedAt = q.ProviderTimestamp
	e.mu.Lock()
	e.providerQuotes["AMD"] = map[string]Quote{
		"Finnhub": q,
		"Alpaca":  {Symbol: "AMD", Price: 104.52, ProviderTimestamp: q.ProviderTimestamp, UpdatedAt: q.UpdatedAt, Source: "alpaca-iex", FeedType: "websocket", DataState: "live"},
	}
	e.mu.Unlock()
	e.evaluateRapidMoveObservation("AMD", q)

	e.mu.RLock()
	second := e.rapidMoveEvents["AMD"]
	retained := e.rapidMoveScorecard.HysteresisRetained
	e.mu.RUnlock()
	if second.State != first.State || !second.Alerted || retained < 1 {
		t.Fatalf("hysteresis must preserve material alerted state until FADING/RESOLVED logic owns de-escalation: before=%+v after=%+v score=%+v", first, second, e.rapidMoveScorecard)
	}
}

func TestV1806RapidMoveDedupPreservesOriginalProviderOutcomeAnchor(t *testing.T) {
	e := newV1801Engine(t)
	at := time.Date(2026, 8, 13, 14, 0, 0, 0, easternLocation())
	q := seedRapidMoveAt(t, e, "CRM", 100, 106, "AGREED", true, at)
	e.evaluateRapidMoveObservation("CRM", q)
	e.mu.RLock()
	first := e.rapidMoveEvents["CRM"]
	e.mu.RUnlock()
	if first.EventProviderAt <= 0 {
		t.Fatalf("expected provider-time event anchor: %+v", first)
	}

	// A deduped update must not move the original event anchor.
	q.ProviderTimestamp += 2 * 60_000
	q.UpdatedAt = q.ProviderTimestamp
	q.Price = 107
	q.Bid, q.Ask = 106.9, 107.1
	e.mu.Lock()
	e.history["CRM"] = []HistoryPoint{{T: q.ProviderTimestamp - 120_000, P: 100}, {T: q.ProviderTimestamp - 60_000, P: 100}}
	e.providerQuotes["CRM"] = map[string]Quote{"Finnhub": q, "Alpaca": {Symbol: "CRM", Price: 107.02, ProviderTimestamp: q.ProviderTimestamp, UpdatedAt: q.UpdatedAt, Source: "alpaca-iex"}}
	e.mu.Unlock()
	e.evaluateRapidMoveObservation("CRM", q)
	e.mu.RLock()
	mid := e.rapidMoveEvents["CRM"]
	e.mu.RUnlock()
	if mid.EventProviderAt != first.EventProviderAt {
		t.Fatalf("dedupe moved outcome anchor: first=%d mid=%d", first.EventProviderAt, mid.EventProviderAt)
	}

	// Twenty-one minutes from the original anchor must resolve even after the
	// intermediate live update.
	q.ProviderTimestamp = first.EventProviderAt + 21*60_000
	q.UpdatedAt = q.ProviderTimestamp
	q.Price = 108
	q.Bid, q.Ask = 107.9, 108.1
	e.mu.Lock()
	e.history["CRM"] = []HistoryPoint{{T: q.ProviderTimestamp - 120_000, P: 108}, {T: q.ProviderTimestamp - 60_000, P: 108}}
	e.providerQuotes["CRM"] = map[string]Quote{"Finnhub": q}
	e.mu.Unlock()
	e.evaluateRapidMoveObservation("CRM", q)
	e.mu.RLock()
	final := e.rapidMoveEvents["CRM"]
	e.mu.RUnlock()
	if final.State != "RESOLVED" || final.Outcome20mPct == nil {
		t.Fatalf("original provider anchor must drive durable 20m outcome resolution: %+v", final)
	}
}

func TestV1806RapidMoveGovernanceIsExplicitAndCannotAutoPromote(t *testing.T) {
	g := rapidMovePolicyGovernance()
	if g.DetectionStage != "PRODUCTION" || g.LearningStage != "SHADOW" || g.AutoPromotion {
		t.Fatalf("adaptive governance boundary drift: %+v", g)
	}
	want := []string{"SHADOW", "VALIDATED", "APPROVED", "PRODUCTION"}
	if len(g.PromotionPath) != len(want) {
		t.Fatalf("promotion path drift: %+v", g)
	}
	for i := range want {
		if g.PromotionPath[i] != want[i] {
			t.Fatalf("promotion path drift: %+v", g)
		}
	}
	if !strings.Contains(g.ProtectedFormulaImpact, "NONE") {
		t.Fatalf("protected deterministic formula boundary must be explicit: %+v", g)
	}
}

func configureTwelveDataGlobalRouteTest(t *testing.T, e *Engine) {
	t.Helper()
	oldGlobal := twelveGlobalSearch
	oldFutures := twelveFutureSearch
	t.Cleanup(func() {
		twelveGlobalSearch = oldGlobal
		twelveFutureSearch = oldFutures
	})
	twelveGlobalSearch = map[string]struct{ Label, Query, Group string }{
		"test_index": {Label: "Test Index", Query: "TEST INDEX", Group: "test"},
	}
	twelveFutureSearch = map[string]struct{ Label, Query string }{
		"test_future": {Label: "Test Future", Query: "TEST FUTURE"},
	}
	e.app.mu.Lock()
	e.app.secrets.TwelveData = "test-key"
	e.app.mu.Unlock()
}

func withTwelveDataBaseURL(t *testing.T, raw string) {
	t.Helper()
	old := twelveDataBaseURL
	twelveDataBaseURL = raw
	t.Cleanup(func() { twelveDataBaseURL = old })
}

func TestProviderRouterTwelveDataGlobalContextUsesSharedTransport(t *testing.T) {
	e := newV1801Engine(t)
	configureTwelveDataGlobalRouteTest(t, e)
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/symbol_search":
			typ := "Index"
			name := "Test Index"
			if strings.Contains(strings.ToLower(r.URL.Query().Get("symbol")), "future") {
				typ = "Future"
				name = "Test Future"
			}
			_, _ = fmt.Fprintf(w, `{"data":[{"symbol":"TEST","instrument_name":%q,"exchange":"TEST","country":"US","instrument_type":%q}],"status":"ok"}`, name, typ)
		case "/quote":
			_, _ = fmt.Fprintf(w, `{"symbol":"TEST","name":"Test","exchange":"TEST","currency":"USD","close":"100","open":"99","high":"101","low":"98","previous_close":"99","percent_change":"1.01","timestamp":%d,"status":"ok"}`, time.Now().Unix())
		default:
			http.NotFound(w, r)
		}
	}))
	defer server.Close()
	withTwelveDataBaseURL(t, server.URL)

	e.refreshDirectGlobal(context.Background(), "test-key")
	telemetry := telemetryForProvider(e.providerTelemetry.Diagnostics(), "Twelve Data")
	if telemetry.Successes < 4 {
		t.Fatalf("global context must use shared bounded provider transport; telemetry=%+v", telemetry)
	}
	e.mu.RLock()
	cap := e.providerCapabilityCircuits[providerCapabilityCircuitKey("Twelve Data", canonicalGlobalMarketContextDataset)]
	index := e.globalDirect["test_index"]
	future := e.globalDirect["test_future"]
	last := e.lastUpdated["global-direct"]
	e.mu.RUnlock()
	if cap.LastSuccess <= 0 || cap.Failures != 0 {
		t.Fatalf("successful global route must populate capability success state: %+v", cap)
	}
	if index.Source != "Twelve Data" || future.Source != "Twelve Data" || last <= 0 {
		t.Fatalf("successful routed refresh must publish direct global/futures context: index=%+v future=%+v last=%d", index, future, last)
	}
	chain := routeChains()[canonicalGlobalMarketContextDataset]
	if len(chain) != 1 || chain[0] != "Twelve Data" || providerInstrumentClass(canonicalGlobalMarketContextDataset) != "GLOBAL_MARKET" {
		t.Fatalf("global context must own an explicit Twelve Data capability route: chain=%+v class=%s", chain, providerInstrumentClass(canonicalGlobalMarketContextDataset))
	}
}

func TestProviderRouterTwelveDataGlobalFailureDoesNotSuppressLiveEquities(t *testing.T) {
	e := newV1801Engine(t)
	configureTwelveDataGlobalRouteTest(t, e)
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		http.Error(w, "provider failure", http.StatusBadGateway)
	}))
	defer server.Close()
	withTwelveDataBaseURL(t, server.URL)

	const priorStamp int64 = 1234567890
	e.mu.Lock()
	if e.globalDirect == nil {
		e.globalDirect = map[string]GlobalDriver{}
	}
	e.globalDirect["test_index"] = GlobalDriver{Key: "test_index", Label: "Prior Index", Value: 99, Source: "Twelve Data"}
	e.lastUpdated["global-direct"] = priorStamp
	e.health["global-direct"] = "healthy · prior direct context"
	beforeGlobal := e.providerCircuits[providerKey("Twelve Data")]
	e.mu.Unlock()

	e.refreshDirectGlobal(context.Background(), "test-key")
	e.mu.RLock()
	cap := e.providerCapabilityCircuits[providerCapabilityCircuitKey("Twelve Data", canonicalGlobalMarketContextDataset)]
	afterGlobal := e.providerCircuits[providerKey("Twelve Data")]
	stamp := e.lastUpdated["global-direct"]
	prior := e.globalDirect["test_index"]
	health := e.health["global-direct"]
	e.mu.RUnlock()
	if cap.Failures < 1 || cap.LastFailure <= 0 {
		t.Fatalf("terminal global failure must populate capability failure state: %+v", cap)
	}
	if afterGlobal.Failures != beforeGlobal.Failures || afterGlobal.OpenUntil != beforeGlobal.OpenUntil || afterGlobal.RateLimitedUntil != beforeGlobal.RateLimitedUntil || afterGlobal.LastFailure != beforeGlobal.LastFailure || afterGlobal.LastError != beforeGlobal.LastError {
		t.Fatalf("global context failure must not mutate Twelve Data global failure circuit: before=%+v after=%+v", beforeGlobal, afterGlobal)
	}
	if stamp != priorStamp || prior.Value != 99 {
		t.Fatalf("failed refresh must preserve prior canonical global context and freshness: stamp=%d prior=%+v", stamp, prior)
	}
	if !strings.Contains(strings.ToLower(health), "degraded") {
		t.Fatalf("terminal provider failure must make direct-global health explicit: %q", health)
	}
	if !e.providerAllowedFor("US Live Equities", "Twelve Data") {
		t.Fatal("global-context capability failure must not suppress Twelve Data US Live Equities eligibility")
	}
	if ok, state := e.providerCapabilityAllowed(canonicalGlobalMarketContextDataset, "Twelve Data", time.Now()); ok || state == "" {
		t.Fatalf("failed global capability must be actively suppressed until revalidation, ok=%v state=%q", ok, state)
	}
}

func TestProviderRouterTwelveDataGlobalCancellationIsNeutral(t *testing.T) {
	e := newV1801Engine(t)
	configureTwelveDataGlobalRouteTest(t, e)
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		http.Error(w, "should not determine provider health", http.StatusServiceUnavailable)
	}))
	defer server.Close()
	withTwelveDataBaseURL(t, server.URL)

	const priorStamp int64 = 2234567890
	e.mu.Lock()
	if e.globalDirect == nil {
		e.globalDirect = map[string]GlobalDriver{}
	}
	e.globalDirect["test_index"] = GlobalDriver{Key: "test_index", Label: "Prior Index", Value: 88, Source: "Twelve Data"}
	e.lastUpdated["global-direct"] = priorStamp
	e.health["global-direct"] = "healthy · prior direct context"
	beforeGlobal := e.providerCircuits[providerKey("Twelve Data")]
	beforeCap := e.providerCapabilityCircuits[providerCapabilityCircuitKey("Twelve Data", canonicalGlobalMarketContextDataset)]
	e.mu.Unlock()

	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	e.refreshDirectGlobal(ctx, "test-key")
	e.mu.RLock()
	afterGlobal := e.providerCircuits[providerKey("Twelve Data")]
	afterCap := e.providerCapabilityCircuits[providerCapabilityCircuitKey("Twelve Data", canonicalGlobalMarketContextDataset)]
	stamp := e.lastUpdated["global-direct"]
	prior := e.globalDirect["test_index"]
	health := e.health["global-direct"]
	e.mu.RUnlock()
	if afterCap.Failures != beforeCap.Failures || afterCap.LastFailure != beforeCap.LastFailure || afterCap.LastError != beforeCap.LastError {
		t.Fatalf("caller cancellation must remain neutral to global capability health: before=%+v after=%+v", beforeCap, afterCap)
	}
	if afterGlobal.Failures != beforeGlobal.Failures || afterGlobal.OpenUntil != beforeGlobal.OpenUntil || afterGlobal.RateLimitedUntil != beforeGlobal.RateLimitedUntil || afterGlobal.LastFailure != beforeGlobal.LastFailure || afterGlobal.LastError != beforeGlobal.LastError {
		t.Fatalf("caller cancellation must remain neutral to Twelve Data global circuit: before=%+v after=%+v", beforeGlobal, afterGlobal)
	}
	if stamp != priorStamp || prior.Value != 88 || health != "healthy · prior direct context" {
		t.Fatalf("canceled refresh must preserve prior global context, freshness and health: stamp=%d prior=%+v health=%q", stamp, prior, health)
	}
}
