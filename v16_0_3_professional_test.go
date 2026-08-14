package main

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func researchComponent(pkg ResearchPackageTruth, name string) (ResearchEvidenceComponent, bool) {
	for _, c := range pkg.Components {
		if c.Dataset == name {
			return c, true
		}
	}
	return ResearchEvidenceComponent{}, false
}

func agreedResearchRec(now time.Time) []ProviderReconciliationDecision {
	return []ProviderReconciliationDecision{{Dataset: "US Live Equities", Symbol: "AAPL", State: "AGREED", UpdatedAt: now.UnixMilli(), Observations: []ProviderQuoteObservation{
		{Provider: "Alpaca", Price: 200, ProviderTimestamp: now.Add(-time.Second).UnixMilli(), ReceivedAt: now.Add(-500 * time.Millisecond).UnixMilli()},
		{Provider: "Finnhub", Price: 200.01, ProviderTimestamp: now.Add(-time.Second).UnixMilli(), ReceivedAt: now.Add(-500 * time.Millisecond).UnixMilli()},
	}}}
}

func TestV1603ProfessionalRecentCheckCannotFreshenOldFundamentals(t *testing.T) {
	now := time.Date(2026, 8, 11, 14, 0, 0, 0, time.UTC)
	st, quotes, bars, fundamentals, last, health := completeResearchFixture(now)
	f := fundamentals["AAPL"]
	f.UpdatedAt = now.Add(-10 * 24 * time.Hour).UnixMilli()
	fundamentals["AAPL"] = f
	last["research-fundamentals:AAPL"] = now.Add(-time.Minute).UnixMilli()
	pkg := buildResearchPackageTruth(st, quotes, bars, fundamentals, nil, nil, nil, last, health, agreedResearchRec(now), now)
	c, ok := researchComponent(pkg, "Fundamentals")
	if !ok {
		t.Fatal("fundamentals component missing")
	}
	if c.State == "FRESH" || c.DataAgeMs < int64((9*24*time.Hour)/time.Millisecond) {
		t.Fatalf("old fundamental data was freshened by recent check: %+v", c)
	}
}

func TestV1603EdgeFutureFundamentalTimestampCannotBeFresh(t *testing.T) {
	now := time.Date(2026, 8, 11, 14, 0, 0, 0, time.UTC)
	st, quotes, bars, fundamentals, last, health := completeResearchFixture(now)
	f := fundamentals["AAPL"]
	f.UpdatedAt = now.Add(5 * time.Minute).UnixMilli()
	fundamentals["AAPL"] = f
	pkg := buildResearchPackageTruth(st, quotes, bars, fundamentals, nil, nil, nil, last, health, agreedResearchRec(now), now)
	c, _ := researchComponent(pkg, "Fundamentals")
	if c.State == "FRESH" || !strings.Contains(strings.ToLower(c.Detail), "future") {
		t.Fatalf("future-skewed fundamentals timestamp was treated as fresh: %+v", c)
	}
}

func TestV1603ProfessionalFutureQuoteUsesCanonicalSkewPolicy(t *testing.T) {
	now := time.Date(2026, 8, 11, 14, 0, 0, 0, time.UTC)
	st, quotes, bars, fundamentals, last, health := completeResearchFixture(now)
	q := quotes["AAPL"]
	q.ProviderTimestamp = now.Add(5 * time.Minute).UnixMilli()
	q.UpdatedAt = now.UnixMilli()
	quotes["AAPL"] = q
	pkg := buildResearchPackageTruth(st, quotes, bars, fundamentals, nil, nil, nil, last, health, agreedResearchRec(now), now)
	c, _ := researchComponent(pkg, "Quote")
	if c.State == "FRESH" || !strings.Contains(strings.ToLower(c.Detail), "future") {
		t.Fatalf("future quote bypassed canonical timestamp truth: %+v", c)
	}
}

func TestV1603EdgeSmallClockSkewIsToleratedConsistently(t *testing.T) {
	now := time.Date(2026, 8, 11, 15, 30, 0, 0, time.UTC)
	q := Quote{Symbol: "AAPL", Price: 200, Source: "alpaca", ProviderTimestamp: now.Add(20 * time.Second).UnixMilli(), UpdatedAt: now.UnixMilli()}
	pa, ra, current := quoteEvidenceAges(q, now.UnixMilli())
	if !current || pa != 0 || ra != 0 {
		t.Fatalf("small allowed skew rejected: providerAge=%d receiptAge=%d current=%v", pa, ra, current)
	}
}

func TestV1603ProfessionalRepeatedCorporateTokenNeverStampsBackfillComplete(t *testing.T) {
	oldBase := alpacaDataBaseURL
	defer func() { alpacaDataBaseURL = oldBase }()
	calls := 0
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v1/corporate-actions" {
			http.NotFound(w, r)
			return
		}
		calls++
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"forward_splits":[{"id":"ca1","symbol":"AAPL","process_date":"2024-06-01","old_rate":1,"new_rate":4}],"next_page_token":"same"}`))
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
	app.mu.Lock()
	app.state.Watchlists = []Watchlist{{ID: "day", Symbols: []string{"AAPL"}}, {ID: "swing"}, {ID: "long"}, {ID: "discovery"}}
	app.mu.Unlock()
	app.engine.refreshAlpacaCorporateActions(context.Background(), "k", "s")
	app.engine.mu.RLock()
	stamp := app.engine.lastUpdated["corporate-actions-backfill:AAPL"]
	h := app.engine.health["corporate-actions"]
	app.engine.mu.RUnlock()
	if calls < 2 || stamp != 0 || !strings.Contains(strings.ToLower(h), "partial") {
		t.Fatalf("truncated backfill falsely closed calls=%d stamp=%d health=%q", calls, stamp, h)
	}
}

func TestV1603ProfessionalNaturalCorporatePaginationStampsBackfill(t *testing.T) {
	oldBase := alpacaDataBaseURL
	defer func() { alpacaDataBaseURL = oldBase }()
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v1/corporate-actions" {
			http.NotFound(w, r)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		if r.URL.Query().Get("page_token") == "p2" {
			_, _ = w.Write([]byte(`{"cash_dividends":[{"id":"ca2","symbol":"AAPL","process_date":"2025-01-01","cash_rate":0.25}]}`))
			return
		}
		_, _ = w.Write([]byte(`{"forward_splits":[{"id":"ca1","symbol":"AAPL","process_date":"2024-06-01","old_rate":1,"new_rate":4}],"next_page_token":"p2"}`))
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
	app.mu.Lock()
	app.state.Watchlists = []Watchlist{{ID: "day", Symbols: []string{"AAPL"}}, {ID: "swing"}, {ID: "long"}, {ID: "discovery"}}
	app.mu.Unlock()
	app.engine.refreshAlpacaCorporateActions(context.Background(), "k", "s")
	app.engine.mu.RLock()
	stamp := app.engine.lastUpdated["corporate-actions-backfill:AAPL"]
	h := app.engine.health["corporate-actions"]
	app.engine.mu.RUnlock()
	if stamp == 0 || !strings.Contains(strings.ToLower(h), "healthy") {
		t.Fatalf("natural backfill did not close stamp=%d health=%q", stamp, h)
	}
}

func TestV1603EdgePartialCorporateBackfillRemainsRetryableAndRecovers(t *testing.T) {
	oldBase := alpacaDataBaseURL
	defer func() { alpacaDataBaseURL = oldBase }()
	fail := true
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v1/corporate-actions" {
			http.NotFound(w, r)
			return
		}
		if r.URL.Query().Get("page_token") == "p2" {
			if fail {
				http.Error(w, "temporary", 500)
				return
			}
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`{"cash_dividends":[{"id":"ca2","symbol":"AAPL","process_date":"2025-01-01","cash_rate":0.25}]}`))
			return
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"forward_splits":[{"id":"ca1","symbol":"AAPL","process_date":"2024-06-01","old_rate":1,"new_rate":4}],"next_page_token":"p2"}`))
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
	app.mu.Lock()
	app.state.Watchlists = []Watchlist{{ID: "day", Symbols: []string{"AAPL"}}, {ID: "swing"}, {ID: "long"}, {ID: "discovery"}}
	app.mu.Unlock()
	app.engine.refreshAlpacaCorporateActions(context.Background(), "k", "s")
	if app.engine.lastUpdated["corporate-actions-backfill:AAPL"] != 0 {
		t.Fatal("partial backfill stamped complete")
	}
	fail = false
	app.engine.refreshAlpacaCorporateActions(context.Background(), "k", "s")
	if app.engine.lastUpdated["corporate-actions-backfill:AAPL"] == 0 {
		t.Fatal("successful retry did not stamp complete")
	}
}

func TestV1603ProfessionalAIReadinessRejectsOldFundamentalsDespiteRecentCheck(t *testing.T) {
	app := newTestApplication(t)
	sym := "AAA"
	now := time.Date(2026, 8, 11, 11, 0, 0, 0, easternLocation())
	fresh := now.UnixMilli()
	app.engine.mu.Lock()
	app.engine.quotes[sym] = Quote{Symbol: sym, Price: 100, ProviderTimestamp: fresh, UpdatedAt: fresh, Source: "Alpaca"}
	app.engine.bars[sym] = map[string][]Bar{
		"intraday": {{T: now.Unix(), O: 99, H: 101, L: 98, C: 100, V: 1000}},
		"daily":    {{T: now.Unix(), O: 95, H: 101, L: 94, C: 100, V: 10000}},
		"weekly":   {{T: now.Unix(), O: 90, H: 102, L: 89, C: 100, V: 50000}},
	}
	app.engine.fundamentals[sym] = FundamentalSnapshot{Symbol: sym, RevenueGrowth: 10, UpdatedAt: now.Add(-10 * 24 * time.Hour).UnixMilli(), Source: "Finnhub"}
	for _, k := range []string{"quotes", "history", "history-intraday", "history-daily", "research-news:" + sym, "research-earnings:" + sym, "research-fundamentals:" + sym, "research-sec:" + sym} {
		app.engine.lastUpdated[k] = fresh
	}
	app.engine.health["news"] = "healthy · Finnhub"
	app.engine.health["earnings"] = "healthy · Finnhub"
	app.engine.health["fundamentals"] = "healthy · Finnhub"
	app.engine.health["filings"] = "healthy · SEC EDGAR"
	app.engine.health["history"] = "healthy · Alpaca"
	app.engine.mu.Unlock()
	ready, issues := app.engine.researchPackageReadinessAt(sym, now)
	if ready || !strings.Contains(strings.ToLower(strings.Join(issues, " ")), "fundamentals") {
		t.Fatalf("AI/readiness accepted stale fundamentals: ready=%v issues=%v", ready, issues)
	}
}

func TestV1603ProfessionalAIReadinessRejectsFutureQuote(t *testing.T) {
	app := newTestApplication(t)
	sym := "AAA"
	now := time.Date(2026, 8, 11, 11, 0, 0, 0, easternLocation())
	fresh := now.UnixMilli()
	app.engine.mu.Lock()
	app.engine.quotes[sym] = Quote{Symbol: sym, Price: 100, ProviderTimestamp: now.Add(5 * time.Minute).UnixMilli(), UpdatedAt: fresh, Source: "Alpaca"}
	app.engine.bars[sym] = map[string][]Bar{"intraday": {{T: now.Unix(), C: 100}}, "daily": {{T: now.Unix(), C: 100}}, "weekly": {{T: now.Unix(), C: 100}}}
	app.engine.fundamentals[sym] = FundamentalSnapshot{Symbol: sym, RevenueGrowth: 10, UpdatedAt: fresh, Source: "Finnhub"}
	for _, k := range []string{"quotes", "history", "history-intraday", "history-daily", "research-news:" + sym, "research-earnings:" + sym, "research-fundamentals:" + sym, "research-sec:" + sym} {
		app.engine.lastUpdated[k] = fresh
	}
	app.engine.health["news"] = "healthy"
	app.engine.health["earnings"] = "healthy"
	app.engine.health["fundamentals"] = "healthy"
	app.engine.health["filings"] = "healthy"
	app.engine.health["history"] = "healthy"
	app.engine.mu.Unlock()
	ready, issues := app.engine.researchPackageReadinessAt(sym, now)
	if ready || !strings.Contains(strings.ToLower(strings.Join(issues, " ")), "quote") {
		t.Fatalf("AI/readiness accepted future-skewed quote: ready=%v issues=%v", ready, issues)
	}
}
