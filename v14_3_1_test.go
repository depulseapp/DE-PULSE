package main

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func withSingleTrackedSymbol(t *testing.T, sym string) func() {
	t.Helper()
	old := append([]string(nil), masterMarketSymbols...)
	masterMarketSymbols = []string{sym}
	return func() { masterMarketSymbols = old }
}

func TestV1431CapabilityRequiresVerifiedHealth(t *testing.T) {
	if got := capabilityStatusFromHealth(true, ""); got != "TEMPORARILY UNAVAILABLE" {
		t.Fatalf("configured-but-unverified must not be AVAILABLE: %s", got)
	}
	if got := capabilityStatusFromHealth(true, "checking"); got != "TEMPORARILY UNAVAILABLE" {
		t.Fatalf("checking must not be AVAILABLE: %s", got)
	}
	if got := capabilityStatusFromHealth(true, "healthy · verified provider response"); got != "AVAILABLE" {
		t.Fatalf("verified healthy capability should be AVAILABLE: %s", got)
	}
	if got := capabilityStatusFromHealth(true, "403 plan limited"); got != "PLAN LIMITED" {
		t.Fatalf("entitlement failure should be PLAN LIMITED: %s", got)
	}
}

func TestV1431FinnhubIntelligenceAdapterFixtures(t *testing.T) {
	defer withSingleTrackedSymbol(t, "AAPL")()
	oldBase, oldInterval := finnhubAPIBaseURL, finnhubMinRequestInterval
	defer func() { finnhubAPIBaseURL, finnhubMinRequestInterval = oldBase, oldInterval }()
	finnhubMinRequestInterval = 0
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch r.URL.Path {
		case "/stock/earnings":
			json.NewEncoder(w).Encode([]map[string]any{{"period": "2026-Q2", "actual": 2.2, "estimate": 2.0, "surprise": .2, "surprisePercent": 10.0}})
		case "/stock/peers":
			json.NewEncoder(w).Encode([]string{"MSFT", "GOOG"})
		case "/stock/recommendation":
			json.NewEncoder(w).Encode([]map[string]any{{"buy": 12, "strongBuy": 8, "hold": 4, "sell": 1, "strongSell": 0}})
		case "/stock/price-target":
			json.NewEncoder(w).Encode(map[string]any{"targetMean": 250.0})
		case "/stock/insider-transactions":
			json.NewEncoder(w).Encode(map[string]any{"data": []map[string]any{{"change": 1000.0}}})
		case "/stock/fund-ownership":
			json.NewEncoder(w).Encode(map[string]any{"ownership": []map[string]any{{"change": 500.0}, {"change": -100.0}}})
		default:
			http.NotFound(w, r)
		}
	}))
	defer srv.Close()
	finnhubAPIBaseURL = srv.URL
	app := newTestApplication(t)
	app.state.Watchlists = []Watchlist{{ID: "day", Name: "Day", Symbols: []string{"AAPL"}}}
	app.state.Settings.DayEnabled, app.state.Settings.SwingEnabled, app.state.Settings.LongEnabled = true, false, false
	app.state.Settings.DayWatchlistID = "day"
	app.state.UI.SelectedTicker = "AAPL"
	app.engine.refreshFinnhubIntelligence(context.Background(), "test")
	si := app.engine.symbolIntelligence["AAPL"]
	if len(si.EarningsSurprises) != 1 || len(si.Peers) != 2 || si.RecommendationTrend != "BULLISH" || si.PriceTarget != 250 || si.InstitutionalOwners != 2 {
		t.Fatalf("unexpected Finnhub fixture parse: %+v", si)
	}
}

func TestV1431AlpacaAdaptersFixtures(t *testing.T) {
	defer withSingleTrackedSymbol(t, "AAPL")()
	oldTrade, oldData := alpacaTradingBaseURL, alpacaDataBaseURL
	defer func() { alpacaTradingBaseURL, alpacaDataBaseURL = oldTrade, oldData }()
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch {
		case r.URL.Path == "/v2/calendar":
			json.NewEncoder(w).Encode([]AlpacaCalendarDay{{Date: "2026-08-10", Open: "09:30", Close: "16:00"}})
		case strings.Contains(r.URL.Path, "most-actives"):
			json.NewEncoder(w).Encode(map[string]any{"most_actives": []map[string]any{{"symbol": "AAPL", "volume": 12345, "percent_change": 2.5}}})
		case strings.Contains(r.URL.Path, "movers"):
			json.NewEncoder(w).Encode(map[string]any{"gainers": []map[string]any{{"symbol": "AAPL", "percent_change": 3.2}}, "losers": []map[string]any{{"symbol": "MSFT", "percent_change": -2.1}}})
		case r.URL.Path == "/v1/corporate-actions":
			if r.URL.Query().Get("symbols") != "AAPL" {
				t.Fatalf("corporate-actions request did not filter tracked symbol: %s", r.URL.RawQuery)
			}
			json.NewEncoder(w).Encode(map[string]any{
				"forward_splits":  []map[string]any{{"symbol": "AAPL", "new_rate": 2.0, "old_rate": 1.0, "process_date": "2026-08-10", "ex_date": "2026-08-11"}},
				"cash_dividends":  []map[string]any{{"symbol": "MSFT", "rate": 0.5, "process_date": "2026-08-10", "ex_date": "2026-08-11"}},
				"next_page_token": nil,
			})
		default:
			http.NotFound(w, r)
		}
	}))
	defer srv.Close()
	alpacaTradingBaseURL, alpacaDataBaseURL = srv.URL, srv.URL
	app := newTestApplication(t)
	app.state.Watchlists = []Watchlist{{ID: "day", Name: "Day", Symbols: []string{"AAPL"}}}
	app.state.Settings.DayEnabled, app.state.Settings.SwingEnabled, app.state.Settings.LongEnabled = true, false, false
	app.state.Settings.DayWatchlistID = "day"
	app.state.UI.SelectedTicker = "AAPL"
	ctx := context.Background()
	app.engine.refreshAlpacaMarketCalendar(ctx, "k", "s")
	app.engine.refreshAlpacaMarketActivity(ctx, "k", "s")
	app.engine.refreshAlpacaCorporateActions(ctx, "k", "s")
	if len(app.engine.alpacaCalendar) == 0 || len(app.engine.marketActivity.MostActive) != 1 || len(app.engine.corporateActions) != 1 || app.engine.corporateActions[0].Type != "forward_split" {
		t.Fatalf("Alpaca adapter fixtures not applied: cal=%d activity=%+v actions=%+v", len(app.engine.alpacaCalendar), app.engine.marketActivity, app.engine.corporateActions)
	}
}

func TestV1431TwelveFXFallbackFixtureAndRetiredFundamentalsRoute(t *testing.T) {
	old := twelveDataBaseURL
	defer func() { twelveDataBaseURL = old }()
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if r.URL.Path != "/quote" {
			http.NotFound(w, r)
			return
		}
		json.NewEncoder(w).Encode(map[string]any{"close": 1.25, "percent_change": .4, "status": "ok"})
	}))
	defer srv.Close()
	twelveDataBaseURL = srv.URL
	app := newTestApplication(t)
	app.engine.refreshTwelveFX(context.Background(), "key")
	if len(app.engine.globalDirect) == 0 {
		t.Fatal("Twelve FX fixture did not populate direct FX")
	}
	// Twelve Data fundamentals were intentionally retired from the active route in v15.
	// Preserve the still-supported FX behavior and guard against accidentally reviving
	// a second Fundamentals owner beside the canonical Finnhub -> SEC -> yfinance route.
	if got := strings.Join(routeChains()["Fundamentals"], ">"); got != "Finnhub>SEC>yfinance" {
		t.Fatalf("retired Twelve fundamentals route was reintroduced: %s", got)
	}
}

func TestV1431CatalystArmedThenTriggered(t *testing.T) {
	app := newTestApplication(t)
	loc := easternLocation()
	now := time.Date(2026, 8, 10, 7, 0, 0, 0, loc)
	app.engine.mu.Lock()
	app.engine.earnings = []EarningsItem{{Symbol: "AAPL", Date: "2026-08-10", Hour: "bmo"}}
	app.engine.quotes["AAPL"] = Quote{Symbol: "AAPL", Price: 200, Bid: 199.9, Ask: 200.1, UpdatedAt: now.UnixMilli(), ProviderTimestamp: now.UnixMilli()}
	app.engine.mu.Unlock()
	app.engine.evaluateCatalystWatch(now)
	if got := app.engine.preparations["catalyst-watch"].State; got != "ARMED" {
		t.Fatalf("expected ARMED before actual release, got %s", got)
	}
	if _, ok := app.engine.catalystReactions["AAPL"]; ok {
		t.Fatal("armed earnings must not create reaction before release confirmation")
	}
	actual := 1.5
	app.engine.mu.Lock()
	app.engine.earnings = []EarningsItem{{Symbol: "AAPL", Date: "2026-08-10", Hour: "bmo", EPSActual: &actual}}
	app.engine.mu.Unlock()
	app.engine.evaluateCatalystWatch(now.Add(time.Minute))
	r, ok := app.engine.catalystReactions["AAPL"]
	if !ok || (r.Phase != "TRIGGERED" && r.Phase != "REACTION" && r.Phase != "PREMARKET REACTION" && r.Phase != "OPENING REACTION") {
		t.Fatalf("expected confirmed earnings reaction, got %+v", r)
	}
}

func TestV1431PremarketSnapshotAndMarketOpenCheckpoint(t *testing.T) {
	loc := easternLocation()
	now := time.Date(2026, 8, 10, 9, 22, 0, 0, loc)
	bars := []Bar{
		{T: time.Date(2026, 8, 10, 8, 0, 0, 0, loc).Unix(), O: 100, H: 102, L: 99, C: 101, V: 1000},
		{T: time.Date(2026, 8, 10, 9, 15, 0, 0, loc).Unix(), O: 101, H: 104, L: 100, C: 103, V: 2000},
	}
	pm, ok := premarketSnapshotFromBars("AAPL", bars, Quote{PreviousClose: 98}, now)
	if !ok || pm.High != 104 || pm.Low != 99 || pm.Volume != 3000 || pm.GapPercent <= 0 {
		t.Fatalf("bad premarket snapshot: %+v", pm)
	}
}

func TestV1436CatalystWatchDoesNotPersistAsRunning(t *testing.T) {
	app := newTestApplication(t)
	loc := easternLocation()
	now := time.Date(2026, 8, 10, 7, 0, 0, 0, loc)
	actual := 1.5
	app.engine.mu.Lock()
	app.engine.earnings = []EarningsItem{{Symbol: "AAPL", Date: "2026-08-10", Hour: "bmo", EPSActual: &actual}}
	app.engine.quotes["AAPL"] = Quote{Symbol: "AAPL", Price: 200, Bid: 199.9, Ask: 200.1, UpdatedAt: now.UnixMilli(), ProviderTimestamp: now.UnixMilli()}
	app.engine.mu.Unlock()

	app.engine.evaluateCatalystWatch(now)
	state := app.engine.preparations["catalyst-watch"].State
	if state != "TRIGGERED" && state != "REACTION" {
		t.Fatalf("confirmed catalyst must expose event phase, not persistent RUNNING; got %s", state)
	}
	if state == "RUNNING" {
		t.Fatal("event-driven catalyst watcher must not remain RUNNING between evaluations")
	}

	// A subsequent quote evaluation may advance to REACTION, but must still not
	// masquerade as a background job that is permanently running.
	app.engine.mu.Lock()
	q := app.engine.quotes["AAPL"]
	q.Price = 204
	q.ProviderTimestamp = now.Add(time.Minute).UnixMilli()
	q.UpdatedAt = q.ProviderTimestamp
	app.engine.quotes["AAPL"] = q
	app.engine.mu.Unlock()
	app.engine.evaluateCatalystWatch(now.Add(time.Minute))
	if got := app.engine.preparations["catalyst-watch"].State; got != "REACTION" {
		t.Fatalf("expected REACTION after measurable post-trigger move, got %s", got)
	}
}
