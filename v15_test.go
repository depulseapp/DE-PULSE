package main

import (
	"context"
	"encoding/json"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func v15TestApp(t *testing.T) *Application {
	t.Helper()
	st := defaultState()
	ensureDedicatedDeskWatchlists(&st, defaultState())
	app := &Application{state: st, configDir: t.TempDir(), hub: NewHub(), sessionKey: "test"}
	app.engine = NewEngine(app)
	return app
}

func TestV15ReleaseIdentity(t *testing.T) {
	if appVersion == "" {
		t.Fatal("runtime version identity must not be empty")
	}
	if buildID == "" {
		t.Fatal("runtime build identity must not be empty")
	}
}

func TestV15ProviderRouteChains(t *testing.T) {
	r := routeChains()
	cases := map[string][]string{
		"US Live Equities": {"Alpaca", "Finnhub", "Twelve Data"},
		"VIX / Indices":    {"Twelve Data", "yfinance", "CBOE"},
		"Historical Bars":  {"Alpaca", "TradeInsight", "Twelve Data", "yfinance"},
		"News":             {"Finnhub", "Marketaux"},
		"Earnings":         {"Finnhub", "yfinance"},
		"Fundamentals":     {"Finnhub", "SEC", "yfinance"},
		"SEC":              {"SEC EDGAR"}, "Macro": {"FRED"},
	}
	for k, want := range cases {
		got := r[k]
		if strings.Join(got, "|") != strings.Join(want, "|") {
			t.Fatalf("%s route=%v want=%v", k, got, want)
		}
	}
}

func TestV15CircuitBreakerAndRecovery(t *testing.T) {
	app := v15TestApp(t)
	e := app.engine
	for i := 0; i < 3; i++ {
		e.recordProviderFailure("Finnhub", context.DeadlineExceeded)
	}
	if e.providerAllowed("Finnhub") {
		t.Fatal("Finnhub should be circuit-open after 3 failures")
	}
	e.recordProviderSuccess("Finnhub")
	if !e.providerAllowed("Finnhub") {
		t.Fatal("Finnhub should recover after success")
	}
}

func TestV15VIXDelayedFreshnessTruth(t *testing.T) {
	now := time.Now().UnixMilli()
	st, _, _, _, _ := freshnessState("VIX", "yfinance", "regular", now-8*60*1000, "healthy", now)
	if st != "DELAYED" {
		t.Fatalf("yfinance VIX state=%s", st)
	}
	st, _, _, _, _ = freshnessState("VIX", "CBOE", "closed", now-25*60*1000, "healthy", now)
	if st != "DELAYED" {
		t.Fatalf("CBOE VIX state=%s", st)
	}
	st, _, _, _, _ = freshnessState("VIX", "Twelve Data", "regular", now-40*60*1000, "healthy", now)
	if st != "STALE" {
		t.Fatalf("old direct VIX state=%s", st)
	}
}

func TestV15DeskMembershipToggleAndLastDeskProtection(t *testing.T) {
	app := v15TestApp(t)
	app.mu.Lock()
	for _, id := range deskIDs() {
		setMembershipLocked(&app.state, id, "AAPL", false)
	}
	setMembershipLocked(&app.state, "day", "AAPL", true)
	app.mu.Unlock()
	call := func(desk string) map[string]any {
		rr := httptest.NewRecorder()
		req := httptest.NewRequest("POST", "/api/desk/membership", strings.NewReader(`{"symbol":"AAPL","desk":"`+desk+`"}`))
		req.Header.Set("Content-Type", "application/json")
		app.handleDeskMembership(rr, req)
		if rr.Code != 200 {
			t.Fatalf("desk %s status=%d body=%s", desk, rr.Code, rr.Body.String())
		}
		var out map[string]any
		_ = json.Unmarshal(rr.Body.Bytes(), &out)
		return out
	}
	out := call("day")
	if out["protected"] != true {
		t.Fatalf("last desk must be protected: %v", out)
	}
	call("swing")
	app.mu.RLock()
	m := deskMembershipsLocked(&app.state, "AAPL")
	app.mu.RUnlock()
	if !m["day"] || !m["swing"] {
		t.Fatalf("add swing failed: %v", m)
	}
	call("day")
	app.mu.RLock()
	m = deskMembershipsLocked(&app.state, "AAPL")
	app.mu.RUnlock()
	if m["day"] || !m["swing"] {
		t.Fatalf("remove day from multi membership failed: %v", m)
	}
	call("long")
	call("long")
	app.mu.RLock()
	m = deskMembershipsLocked(&app.state, "AAPL")
	app.mu.RUnlock()
	if m["long"] {
		t.Fatalf("second long click should remove when swing remains: %v", m)
	}
}

func TestV15MasterRemoveAndUndo(t *testing.T) {
	app := v15TestApp(t)
	app.mu.Lock()
	for _, id := range deskIDs() {
		setMembershipLocked(&app.state, id, "NVDA", true)
	}
	app.mu.Unlock()
	rr := httptest.NewRecorder()
	req := httptest.NewRequest("POST", "/api/master-symbol/remove", strings.NewReader(`{"symbol":"NVDA"}`))
	req.Header.Set("Content-Type", "application/json")
	app.handleMasterSymbolRemove(rr, req)
	if rr.Code != 200 {
		t.Fatalf("remove=%d %s", rr.Code, rr.Body.String())
	}
	app.mu.RLock()
	m := deskMembershipsLocked(&app.state, "NVDA")
	app.mu.RUnlock()
	if activeDeskCount(m) != 0 {
		t.Fatalf("global remove left membership %v", m)
	}
	rr = httptest.NewRecorder()
	req = httptest.NewRequest("POST", "/api/master-symbol/restore", strings.NewReader(`{"symbol":"NVDA","membership":{"day":true,"swing":true,"long":true}}`))
	req.Header.Set("Content-Type", "application/json")
	app.handleMasterSymbolRestore(rr, req)
	if rr.Code != 200 {
		t.Fatalf("restore=%d %s", rr.Code, rr.Body.String())
	}
	app.mu.RLock()
	m = deskMembershipsLocked(&app.state, "NVDA")
	app.mu.RUnlock()
	if activeDeskCount(m) != 3 {
		t.Fatalf("undo failed %v", m)
	}
}

func TestV15SECForm4OtherIncludedWithoutFalseBuy(t *testing.T) {
	rows := []FilingItem{
		{Symbol: "AAPL", Form: "4", Category: "insider", Signal: "BUY", TransactionType: "BUY", TransactionCode: "P", Actor: "Director A", Shares: 100, Price: 200, Value: 20000, FiledAt: "2026-08-09"},
		{Symbol: "AAPL", Form: "4", Category: "insider", Signal: "OTHER", TransactionType: "OTHER", TransactionCode: "A", Actor: "Officer B", Shares: 500, FiledAt: "2026-08-08"},
		{Symbol: "AAPL", Form: "4", Category: "insider", Signal: "SELL", TransactionType: "SELL", TransactionCode: "S", Actor: "Officer C", Shares: 50, Price: 205, Value: 10250, FiledAt: "2026-08-07"},
	}
	x := buildSECIntelligence("AAPL", rows)
	if x.InsiderBuys != 1 || x.InsiderSells != 1 || x.InsiderOthers != 1 {
		t.Fatalf("classification=%+v", x)
	}
	if len(x.RecentTransactions) != 3 {
		t.Fatalf("transactions=%d", len(x.RecentTransactions))
	}
}

func TestV15FreshnessDiagnosticsExposeReasonFallbackImpact(t *testing.T) {
	app := v15TestApp(t)
	e := app.engine
	now := time.Now().UnixMilli()
	q := map[string]Quote{"VIX": {Symbol: "VIX", Price: 18, ProviderTimestamp: now - 5*60*1000, UpdatedAt: now, Source: "yfinance-vix:^VIX"}, "SPY": {Symbol: "SPY", Price: 600, UpdatedAt: now, Source: "alpaca-iex"}}
	last := map[string]int64{"quotes": now, "history": now, "news": now, "earnings": now, "filings": now, "fundamentals": now, "global": now, "macro": now, "options": now, "cache": now}
	health := map[string]string{"quotes": "healthy", "history": "healthy · Alpaca", "news": "healthy · Finnhub", "earnings": "healthy · Finnhub", "filings": "healthy · SEC EDGAR", "fundamentals": "healthy · Finnhub"}
	rows, sum := e.buildFreshnessDiagnostics(q, last, health)
	if len(rows) < 9 {
		t.Fatalf("freshness rows=%d", len(rows))
	}
	found := false
	for _, r := range rows {
		if r.Dataset == "VIX" {
			found = true
			if r.State != "DELAYED" || r.Fallback == "" || len(r.Affected) == 0 || r.Action != "vix" {
				t.Fatalf("VIX diagnostic=%+v", r)
			}
		}
	}
	if !found {
		t.Fatal("missing VIX diagnostic")
	}
	if sum.Delayed < 1 {
		t.Fatalf("summary=%+v", sum)
	}
}
