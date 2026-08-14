package main

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestV1603AuthorizationFreshFundamentalsWithStaleCheckIsPartialNotStaleData(t *testing.T) {
	now := time.Date(2026, 8, 11, 14, 0, 0, 0, time.UTC)
	st, quotes, bars, fundamentals, last, health := completeResearchFixture(now)
	f := fundamentals["AAPL"]
	f.UpdatedAt = now.Add(-2 * time.Hour).UnixMilli()
	fundamentals["AAPL"] = f
	last["research-fundamentals:AAPL"] = now.Add(-3 * 24 * time.Hour).UnixMilli()
	pkg := buildResearchPackageTruth(st, quotes, bars, fundamentals, nil, nil, nil, last, health, agreedResearchRec(now), now)
	c, _ := researchComponent(pkg, "Fundamentals")
	if c.State != "PARTIAL" {
		t.Fatalf("fresh data + stale target check should be PARTIAL, got %+v", c)
	}
}

func TestV1603AuthorizationFreshReceiptCannotRescueOldProviderQuote(t *testing.T) {
	now := time.Date(2026, 8, 11, 15, 30, 0, 0, time.UTC)
	q := Quote{Symbol: "AAPL", Price: 200, Source: "alpaca", ProviderTimestamp: now.Add(-10 * time.Minute).UnixMilli(), UpdatedAt: now.Add(-time.Second).UnixMilli()}
	pa, ra, current := quoteEvidenceAges(q, now.UnixMilli())
	if current || pa < int64(9*time.Minute/time.Millisecond) || ra > int64(2*time.Second/time.Millisecond) {
		t.Fatalf("old market data rescued by receipt time: pa=%d ra=%d current=%v", pa, ra, current)
	}
}

func TestV1603AuthorizationPreviouslyBackfilledSymbolKeepsStampWhenAnotherSymbolPartial(t *testing.T) {
	oldBase := alpacaDataBaseURL
	defer func() { alpacaDataBaseURL = oldBase }()
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v1/corporate-actions" {
			http.NotFound(w, r)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"forward_splits":[{"id":"ca1","symbol":"AAPL","process_date":"2024-06-01","old_rate":1,"new_rate":4},{"id":"ca2","symbol":"MSFT","process_date":"2024-06-01","old_rate":1,"new_rate":2}],"next_page_token":"loop"}`))
	}))
	defer srv.Close()
	alpacaDataBaseURL = srv.URL
	root := t.TempDir()
	t.Setenv("XDG_CONFIG_HOME", root)
	t.Setenv("HOME", root)
	t.Setenv("APPDATA", root)
	app, err := NewApplication()
	if err != nil {
		t.Fatal(err)
	}
	registerApplicationCleanup(t, app)
	app.engine = NewEngine(app)
	old := time.Now().Add(-24 * time.Hour).UnixMilli()
	app.engine.lastUpdated["corporate-actions-backfill:AAPL"] = old
	app.mu.Lock()
	app.state.Watchlists = []Watchlist{{ID: "day", Symbols: []string{"AAPL", "MSFT"}}, {ID: "swing"}, {ID: "long"}, {ID: "discovery"}}
	app.mu.Unlock()
	app.engine.refreshAlpacaCorporateActions(context.Background(), "k", "s")
	if app.engine.lastUpdated["corporate-actions-backfill:AAPL"] != old {
		t.Fatal("existing proven backfill stamp was corrupted")
	}
	if app.engine.lastUpdated["corporate-actions-backfill:MSFT"] != 0 {
		t.Fatal("partial new-symbol backfill incorrectly stamped complete")
	}
	if !strings.Contains(strings.ToLower(app.engine.health["corporate-actions"]), "partial") {
		t.Fatalf("partial batch not visible: %q", app.engine.health["corporate-actions"])
	}
}

func TestV1603AuthorizationThirtySecondFutureBoundaryIsNotSilentlyInvalid(t *testing.T) {
	now := time.Date(2026, 8, 11, 15, 30, 0, 0, time.UTC).UnixMilli()
	q := Quote{Symbol: "AAPL", Price: 200, Source: "alpaca", ProviderTimestamp: now + 30_000, UpdatedAt: now}
	_, _, current := quoteEvidenceAges(q, now)
	if !current {
		t.Fatal("exact allowed skew boundary unexpectedly invalid")
	}
	q.ProviderTimestamp = now + 30_001
	_, _, current = quoteEvidenceAges(q, now)
	if current {
		t.Fatal("timestamp beyond allowed skew boundary treated current")
	}
}

func TestV1603AuthorizationFutureTargetRefreshCannotMakeAIReady(t *testing.T) {
	app := newTestApplication(t)
	sym := "AAA"
	now := time.Date(2026, 8, 11, 11, 0, 0, 0, easternLocation())
	fresh := now.UnixMilli()
	app.engine.mu.Lock()
	app.engine.quotes[sym] = Quote{Symbol: sym, Price: 100, ProviderTimestamp: fresh, UpdatedAt: fresh, Source: "Alpaca"}
	app.engine.bars[sym] = map[string][]Bar{"intraday": {{T: now.Unix(), C: 100}}, "daily": {{T: now.Unix(), C: 100}}, "weekly": {{T: now.Unix(), C: 100}}}
	app.engine.fundamentals[sym] = FundamentalSnapshot{Symbol: sym, RevenueGrowth: 10, UpdatedAt: fresh, Source: "Finnhub"}
	for _, k := range []string{"research-news:" + sym, "research-earnings:" + sym, "research-fundamentals:" + sym, "research-sec:" + sym} {
		app.engine.lastUpdated[k] = fresh
	}
	app.engine.lastUpdated["research-fundamentals:"+sym] = now.Add(5 * time.Minute).UnixMilli()
	app.engine.health["research-sec:"+sym] = "healthy"
	app.engine.mu.Unlock()
	ready, issues := app.engine.researchPackageReadinessAt(sym, now)
	if ready || !strings.Contains(strings.ToLower(strings.Join(issues, " ")), "fundamentals") {
		t.Fatalf("future target refresh incorrectly made AI ready: ready=%v issues=%v", ready, issues)
	}
}

func TestV1603AuthorizationUniverseHealthCannotReplaceSelectedTickerCheck(t *testing.T) {
	app := newTestApplication(t)
	sym := "AAA"
	now := time.Date(2026, 8, 11, 11, 0, 0, 0, easternLocation())
	fresh := now.UnixMilli()
	app.engine.mu.Lock()
	app.engine.quotes[sym] = Quote{Symbol: sym, Price: 100, ProviderTimestamp: fresh, UpdatedAt: fresh, Source: "Alpaca"}
	app.engine.bars[sym] = map[string][]Bar{"intraday": {{T: now.Unix(), C: 100}}, "daily": {{T: now.Unix(), C: 100}}, "weekly": {{T: now.Unix(), C: 100}}}
	app.engine.fundamentals[sym] = FundamentalSnapshot{Symbol: sym, RevenueGrowth: 10, UpdatedAt: fresh, Source: "Finnhub"}
	app.engine.health["news"] = "healthy"
	app.engine.health["earnings"] = "healthy"
	app.engine.health["fundamentals"] = "healthy"
	app.engine.health["filings"] = "healthy"
	app.engine.health["research-sec:"+sym] = "healthy"
	app.engine.lastUpdated["research-sec:"+sym] = fresh
	app.engine.mu.Unlock()
	ready, issues := app.engine.researchPackageReadinessAt(sym, now)
	joined := strings.ToLower(strings.Join(issues, " "))
	if ready || !strings.Contains(joined, "news") || !strings.Contains(joined, "earnings") || !strings.Contains(joined, "fundamentals") {
		t.Fatalf("universe health substituted for selected-ticker checks: ready=%v issues=%v", ready, issues)
	}
}
