package main

import (
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
