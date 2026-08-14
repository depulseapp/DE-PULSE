package main

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

// These tests intentionally inject states that DE.PULSE must reject. The gate is
// valuable only if these bad states cannot be upgraded to professional truth.
func TestV1603FaultInjectionStaleFundamentalsCannotBecomeFresh(t *testing.T) {
	now := time.Date(2026, 8, 11, 14, 0, 0, 0, time.UTC)
	st, quotes, bars, fundamentals, last, health := completeResearchFixture(now)
	f := fundamentals["AAPL"]
	f.UpdatedAt = now.Add(-10 * 24 * time.Hour).UnixMilli()
	fundamentals["AAPL"] = f
	last["research-fundamentals:AAPL"] = now.Add(-time.Minute).UnixMilli()
	pkg := buildResearchPackageTruth(st, quotes, bars, fundamentals, nil, nil, nil, last, health, agreedResearchRec(now), now)
	c, _ := researchComponent(pkg, "Fundamentals")
	if c.State == "FRESH" {
		t.Fatalf("fault injection escaped: stale fundamentals became FRESH: %+v", c)
	}
}

func TestV1603FaultInjectionFutureQuoteCannotBecomeFresh(t *testing.T) {
	now := time.Date(2026, 8, 11, 14, 0, 0, 0, time.UTC)
	st, quotes, bars, fundamentals, last, health := completeResearchFixture(now)
	q := quotes["AAPL"]
	q.ProviderTimestamp = now.Add(5 * time.Minute).UnixMilli()
	q.UpdatedAt = now.UnixMilli()
	quotes["AAPL"] = q
	pkg := buildResearchPackageTruth(st, quotes, bars, fundamentals, nil, nil, nil, last, health, agreedResearchRec(now), now)
	c, _ := researchComponent(pkg, "Quote")
	if c.State == "FRESH" {
		t.Fatalf("fault injection escaped: future quote became FRESH: %+v", c)
	}
}

func TestV1603FaultInjectionPartialCorporateBackfillCannotBecomeComplete(t *testing.T) {
	oldBase := alpacaDataBaseURL
	defer func() { alpacaDataBaseURL = oldBase }()
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v1/corporate-actions" {
			http.NotFound(w, r)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"forward_splits":[{"id":"x","symbol":"AAPL","process_date":"2025-01-01","old_rate":1,"new_rate":2}],"next_page_token":"loop"}`))
	}))
	defer srv.Close()
	alpacaDataBaseURL = srv.URL
	root := t.TempDir()
	t.Setenv("HOME", root)
	t.Setenv("XDG_CONFIG_HOME", root)
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
		t.Fatal("fault injection escaped: partial corporate backfill stamped complete")
	}
	if !strings.Contains(strings.ToLower(app.engine.health["corporate-actions"]), "partial") {
		t.Fatalf("fault injection escaped: partial corporate backfill not degraded: %q", app.engine.health["corporate-actions"])
	}
}

func TestV1603FaultInjectionCanonicalProviderReasonCannotContradictState(t *testing.T) {
	now := time.Date(2026, 8, 11, 15, 30, 0, 0, time.UTC).UnixMilli()
	router := ProviderRouterSnapshot{Routes: []ProviderRouteState{{Dataset: "US Live Equities", Active: "Alpaca", Route: []ProviderRouteHop{{Provider: "Alpaca", Priority: 1}, {Provider: "Finnhub", Priority: 2}}}}}
	obs := map[string]map[string]Quote{"AAPL": {"Alpaca": {Symbol: "AAPL", Price: 200, Source: "alpaca", ProviderTimestamp: now - 1000, UpdatedAt: now - 500}, "Finnhub": {Symbol: "AAPL", Price: 202, Source: "finnhub", ProviderTimestamp: now - 1000, UpdatedAt: now - 500}}}
	canonical := map[string]Quote{"AAPL": {Symbol: "AAPL", Price: 202, Source: "finnhub", ProviderTimestamp: now - 1000, UpdatedAt: now - 500}}
	rows := buildProviderReconciliation(router, obs, canonical, now)
	if len(rows) != 1 {
		t.Fatalf("unexpected reconciliation rows: %d", len(rows))
	}
	r := rows[0]
	if r.CanonicalProvider != "Finnhub" || !strings.Contains(strings.ToLower(r.Reason), "finnhub") {
		t.Fatalf("fault injection escaped: canonical state/reason contradiction: %+v", r)
	}
}
