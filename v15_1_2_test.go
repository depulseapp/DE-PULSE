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

func TestV1512StrictUserTickerParsing(t *testing.T) {
	good := map[string]string{
		"aapl":    "AAPL",
		" BRK.B ": "BRK.B",
		"BF-B":    "BF-B",
		"ABC1":    "ABC1",
	}
	for raw, want := range good {
		got, ok := parseUserTicker(raw)
		if !ok || got != want {
			t.Fatalf("parseUserTicker(%q)=%q,%v want %q,true", raw, got, ok, want)
		}
	}
	for _, raw := range []string{"", "VIX", "^VIX", "AA/PL", "AAPL<script>", "1AAPL", "ABCDEFGHI", "AAPL$"} {
		if got, ok := parseUserTicker(raw); ok {
			t.Fatalf("malformed/dedicated symbol %q unexpectedly accepted as %q", raw, got)
		}
	}
}

func TestV1512DeskMembershipRejectsMalformedAndVIX(t *testing.T) {
	app := newTestApplication(t)
	for _, sym := range []string{"VIX", "AA/PL", "AAPL<script>"} {
		rr := httptest.NewRecorder()
		req := httptest.NewRequest("POST", "/api/desk/membership", strings.NewReader(`{"symbol":"`+sym+`","desk":"day","active":true}`))
		req.Header.Set("Content-Type", "application/json")
		app.handleDeskMembership(rr, req)
		if rr.Code != 400 {
			t.Fatalf("symbol %q status=%d body=%s", sym, rr.Code, rr.Body.String())
		}
	}
}

func TestV1512MasterRestoreRejectsDedicatedVIX(t *testing.T) {
	app := newTestApplication(t)
	rr := httptest.NewRecorder()
	req := httptest.NewRequest("POST", "/api/master-symbol/restore", strings.NewReader(`{"symbol":"VIX","membership":{"day":true}}`))
	req.Header.Set("Content-Type", "application/json")
	app.handleMasterSymbolRestore(rr, req)
	if rr.Code != 400 {
		t.Fatalf("VIX restore status=%d body=%s", rr.Code, rr.Body.String())
	}
}

func TestV1512PermanentDeskSanitizationRemovesInvalidAndVIX(t *testing.T) {
	st := defaultState()
	st.Watchlists = []Watchlist{{ID: "day", Name: "Day", Symbols: []string{"aapl", "VIX", "AA/PL", "BRK.B"}}}
	ensureDedicatedDeskWatchlists(&st, defaultState())
	wl, ok := watchlistValueByID(st.Watchlists, "day")
	if !ok {
		t.Fatal("missing day watchlist")
	}
	got := strings.Join(wl.Symbols, "|")
	if got != "AAPL|BRK.B" {
		t.Fatalf("sanitized day symbols=%s", got)
	}
}

func TestV1512LiveModeProviderEligibility(t *testing.T) {
	cases := []struct {
		name string
		s    Secrets
		want bool
	}{
		{"none", Secrets{}, false},
		{"finnhub", Secrets{Finnhub: "x"}, true},
		{"alpaca complete", Secrets{AlpacaKey: "k", AlpacaSecret: "s"}, true},
		{"alpaca key only", Secrets{AlpacaKey: "k"}, false},
		{"twelve", Secrets{TwelveData: "x"}, true},
	}
	for _, tc := range cases {
		if got := liveEquityProviderConfigured(tc.s); got != tc.want {
			t.Fatalf("%s got=%v want=%v", tc.name, got, tc.want)
		}
	}
}

func TestV1512YahooPercentageNormalization(t *testing.T) {
	cases := map[float64]float64{
		0:     0,
		0.235: 23.5,
		-0.12: -12,
		1.5:   150,
		23.5:  23.5,
		-18:   -18,
	}
	for in, want := range cases {
		if got := yahooPercentValue(in); got != want {
			t.Fatalf("yahooPercentValue(%v)=%v want=%v", in, got, want)
		}
	}
}

func TestV1512CriticalDelayedQuoteDegradesReadiness(t *testing.T) {
	rows := []FreshnessDiagnostic{
		{Dataset: "Quotes", State: "DELAYED", Provider: "fallback", Reason: "15m delay"},
		{Dataset: "VIX", State: "FRESH", Provider: "Twelve Data"},
	}
	usable, degraded, ex := readinessFreshnessGate(rows, []string{"Quotes", "VIX"}, time.Now())
	if usable != 1 || !degraded || len(ex) != 1 || ex[0].Severity != "HIGH" || checkpointAttention(ex, degraded) != "DATA DEGRADED" {
		t.Fatalf("unexpected readiness result usable=%d degraded=%v ex=%+v", usable, degraded, ex)
	}
}

func TestV1512ProfileImportSanitizesTradingDesks(t *testing.T) {
	app := newTestApplication(t)
	profile := map[string]any{
		"settings": defaultState().Settings,
		"watchlists": []Watchlist{
			{ID: "day", Name: "Day", Symbols: []string{"AAPL", "VIX", "AA/PL"}},
			{ID: "swing", Name: "Swing", Symbols: []string{"BRK.B"}},
			{ID: "long", Name: "Long", Symbols: []string{"META"}},
			{ID: "discovery", Name: "Discovery", Symbols: []string{"TSLA"}},
		},
		"ui": defaultState().UI,
	}
	b, _ := json.Marshal(profile)
	body, _ := json.Marshal(map[string]json.RawMessage{"profile": b})
	rr := httptest.NewRecorder()
	req := httptest.NewRequest("POST", "/api/profile/import", strings.NewReader(string(body)))
	req.Header.Set("Content-Type", "application/json")
	app.handleImport(rr, req)
	if rr.Code != 200 {
		t.Fatalf("import status=%d body=%s", rr.Code, rr.Body.String())
	}
	app.mu.RLock()
	day, _ := watchlistValueByID(app.state.Watchlists, "day")
	app.mu.RUnlock()
	if strings.Join(day.Symbols, "|") != "AAPL" {
		t.Fatalf("imported day symbols not sanitized: %v", day.Symbols)
	}
}

func TestV1512FailedLiveStartCleansRuntimeContext(t *testing.T) {
	app := newTestApplication(t)
	app.state.Settings.DataMode = "live"
	if err := app.engine.Start(); err == nil {
		t.Fatal("expected live start without provider credentials to fail")
	}
	app.engine.mu.RLock()
	status := app.engine.status
	cancel := app.engine.cancel
	app.engine.mu.RUnlock()
	if status != "stopped" || cancel != nil {
		t.Fatalf("failed start left dirty runtime state status=%q cancelNil=%v", status, cancel == nil)
	}
}

func TestV1512StartRejectsStoppingRuntime(t *testing.T) {
	app := newTestApplication(t)
	app.engine.mu.Lock()
	app.engine.status = "stopping"
	app.engine.mu.Unlock()
	if err := app.engine.Start(); err == nil || !strings.Contains(strings.ToLower(err.Error()), "stopping") {
		t.Fatalf("expected explicit stopping error, got %v", err)
	}
	app.engine.mu.RLock()
	status := app.engine.status
	app.engine.mu.RUnlock()
	if status != "stopping" {
		t.Fatalf("Start changed stopping status to %q", status)
	}
}

func TestV1512MasterRemoveRejectsMalformedTicker(t *testing.T) {
	app := newTestApplication(t)
	app.mu.Lock()
	_, _, _ = applyDeskMembershipLocked(&app.state, "day", "AAPL", true)
	app.mu.Unlock()
	rr := httptest.NewRecorder()
	req := httptest.NewRequest("POST", "/api/master-symbol/remove", strings.NewReader(`{"symbol":"AA/PL"}`))
	req.Header.Set("Content-Type", "application/json")
	app.handleMasterSymbolRemove(rr, req)
	if rr.Code != 400 {
		t.Fatalf("malformed remove status=%d body=%s", rr.Code, rr.Body.String())
	}
	app.mu.RLock()
	m := deskMembershipsLocked(&app.state, "AAPL")
	app.mu.RUnlock()
	if !m["day"] {
		t.Fatal("malformed remove unexpectedly removed AAPL")
	}
}

func TestV1512ResearchRefreshRejectsMalformedAndVIX(t *testing.T) {
	app := newTestApplication(t)
	for _, sym := range []string{"VIX", "AA/PL", "AAPL<script>"} {
		rr := httptest.NewRecorder()
		req := httptest.NewRequest("POST", "/api/research/refresh", strings.NewReader(`{"symbol":"`+sym+`"}`))
		req.Header.Set("Content-Type", "application/json")
		app.handleResearchRefresh(rr, req)
		if rr.Code != 400 {
			t.Fatalf("symbol %q research status=%d body=%s", sym, rr.Code, rr.Body.String())
		}
	}
}

func TestV1512SettingsNormalizesResearchAIMode(t *testing.T) {
	app := newTestApplication(t)
	settings := defaultState().Settings
	settings.ResearchAIMode = "nonsense"
	body, _ := json.Marshal(map[string]any{"settings": settings})
	rr := httptest.NewRecorder()
	req := httptest.NewRequest("POST", "/api/settings/save", strings.NewReader(string(body)))
	req.Header.Set("Content-Type", "application/json")
	app.handleSettingsSave(rr, req)
	if rr.Code != 200 {
		t.Fatalf("settings status=%d body=%s", rr.Code, rr.Body.String())
	}
	app.mu.RLock()
	got := app.state.Settings.ResearchAIMode
	app.mu.RUnlock()
	if got != "manual" {
		t.Fatalf("ResearchAIMode=%q want manual", got)
	}
}

func TestV1512SessionAwareWaitCapsAtMarketBoundaries(t *testing.T) {
	loc := easternLocation()
	cases := []struct {
		name    string
		now     time.Time
		desired time.Duration
		wantMax time.Duration
		wantMin time.Duration
	}{
		{"overnight-to-premarket", time.Date(2026, 8, 10, 3, 50, 0, 0, loc), 30 * time.Minute, 11 * time.Minute, 9 * time.Minute},
		{"regular-to-afterhours", time.Date(2026, 8, 10, 15, 58, 0, 0, loc), 30 * time.Minute, 3 * time.Minute, time.Minute},
		{"afterhours-to-overnight", time.Date(2026, 8, 10, 19, 58, 0, 0, loc), 30 * time.Minute, 3 * time.Minute, time.Minute},
		{"weekend-keeps-normal-cadence", time.Date(2026, 8, 9, 12, 0, 0, 0, loc), 30 * time.Minute, 31 * time.Minute, 29 * time.Minute},
	}
	for _, tc := range cases {
		got := capWaitToSessionBoundary(tc.now, tc.desired)
		if got < tc.wantMin || got > tc.wantMax {
			t.Fatalf("%s got=%s expected [%s,%s]", tc.name, got, tc.wantMin, tc.wantMax)
		}
	}
}

func TestV1512PreparationRetryCapAndNewDayResetSemantics(t *testing.T) {
	loc := easternLocation()
	now := time.Date(2026, 8, 10, 9, 30, 0, 0, loc)
	day := now.Format("2006-01-02")
	p := PreparationJobStatus{TradingDay: day, State: "DATA DEGRADED", AttemptCount: 2, LastAttempt: now.Add(-6 * time.Minute).UnixMilli()}
	if preparationRanForTradingDay(p, now) {
		t.Fatal("two failed/degraded attempts should not suppress a permitted retry")
	}
	if !preparationRetryDue(p, now, 5*time.Minute) {
		t.Fatal("retry should be due after minimum gap")
	}
	p.AttemptCount = 3
	if !preparationRanForTradingDay(p, now) {
		t.Fatal("three same-day attempts should trip retry cap")
	}
	tomorrow := now.AddDate(0, 0, 1)
	if preparationRanForTradingDay(p, tomorrow) {
		t.Fatal("prior-day attempt cap must not suppress the next trading day")
	}
}

func TestV1512SparseFreshnessUsesCheckAgeNotLatestItemAge(t *testing.T) {
	app := newTestApplication(t)
	now := time.Now()
	app.engine.mu.Lock()
	app.engine.news = []NewsItem{{Headline: "Old but latest", Datetime: now.Add(-2 * time.Hour).Unix()}}
	app.engine.filings = []FilingItem{{Symbol: "NVDA", Form: "8-K", FiledAt: now.Add(-48 * time.Hour).Format("2006-01-02")}}
	app.engine.macroMetrics = map[string]MacroMetric{"cpi": {Key: "cpi", UpdatedAt: now.Add(-20 * 24 * time.Hour).UnixMilli(), Status: "current"}}
	app.engine.lastUpdated["news"] = now.Add(-2 * time.Minute).UnixMilli()
	app.engine.lastUpdated["filings"] = now.Add(-2 * time.Minute).UnixMilli()
	app.engine.lastUpdated["macro"] = now.Add(-2 * time.Minute).UnixMilli()
	app.engine.health["news"] = "healthy · Finnhub"
	app.engine.health["filings"] = "healthy · SEC"
	app.engine.health["macro"] = "healthy · FRED"
	last := clone(app.engine.lastUpdated)
	health := clone(app.engine.health)
	quotes := clone(app.engine.quotes)
	app.engine.mu.Unlock()
	rows, _ := app.engine.buildFreshnessDiagnostics(quotes, last, health)
	by := map[string]FreshnessDiagnostic{}
	for _, r := range rows {
		by[r.Dataset] = r
	}
	for _, dataset := range []string{"News", "SEC Filings", "Macro"} {
		r := by[dataset]
		if r.State == "STALE" || r.State == "ERROR" || r.State == "UNAVAILABLE" {
			t.Fatalf("%s incorrectly stale despite recent successful check: %+v", dataset, r)
		}
		if r.CheckAgeMs <= 0 || r.DataAgeMs <= r.CheckAgeMs {
			t.Fatalf("%s did not preserve separate check/data ages: %+v", dataset, r)
		}
	}
}

func TestV1512EarningsRefreshCadenceIntensifiesAroundRelease(t *testing.T) {
	loc := easternLocation()
	bmo := time.Date(2026, 8, 10, 7, 0, 0, 0, loc)
	events := []EarningsItem{{Symbol: "NVDA", Date: "2026-08-10", Hour: "bmo"}}
	if got := earningsRefreshIntervalFrom(events, bmo); got > 10*time.Minute {
		t.Fatalf("BMO premarket cadence too slow: %s", got)
	}
	prior := time.Date(2026, 8, 9, 12, 0, 0, 0, loc)
	if got := earningsRefreshIntervalFrom(events, prior); got > 30*time.Minute {
		t.Fatalf("within-24h cadence too slow: %s", got)
	}
	far := time.Date(2026, 8, 1, 12, 0, 0, 0, loc)
	if got := earningsRefreshIntervalFrom(events, far); got != 2*time.Hour {
		t.Fatalf("normal cadence=%s want 2h", got)
	}
}

func TestV1512TickerSelectionRejectsMalformedButAllowsVIX(t *testing.T) {
	app := newTestApplication(t)
	for _, sym := range []string{"AA/PL", "AAPL<script>", "^VIX"} {
		rr := httptest.NewRecorder()
		req := httptest.NewRequest("POST", "/api/ui/ticker", strings.NewReader(`{"symbol":"`+sym+`"}`))
		req.Header.Set("Content-Type", "application/json")
		app.handleTicker(rr, req)
		if rr.Code != 400 {
			t.Fatalf("symbol %q selection status=%d body=%s", sym, rr.Code, rr.Body.String())
		}
	}
	rr := httptest.NewRecorder()
	req := httptest.NewRequest("POST", "/api/ui/ticker", strings.NewReader(`{"symbol":"VIX"}`))
	req.Header.Set("Content-Type", "application/json")
	app.handleTicker(rr, req)
	if rr.Code != 200 {
		t.Fatalf("VIX selection status=%d body=%s", rr.Code, rr.Body.String())
	}
}

func TestV1512MergeStateRepairsInvalidSelectedTicker(t *testing.T) {
	st := defaultState()
	st.UI.SelectedTicker = "AA/PL"
	got := mergeState(st)
	if got.UI.SelectedTicker == "" || got.UI.SelectedTicker == "AAPL" || strings.Contains(got.UI.SelectedTicker, "/") {
		t.Fatalf("invalid selected ticker was not safely repaired: %q", got.UI.SelectedTicker)
	}
}

func TestV1512Form4KeepsLegitimateSameTableDuplicatesButSuppressesMirroredDerivative(t *testing.T) {
	tx := `<transactionDate><value>2026-08-08</value></transactionDate><transactionCoding><transactionCode>P</transactionCode></transactionCoding><transactionAmounts><transactionShares><value>100</value></transactionShares><transactionPricePerShare><value>10</value></transactionPricePerShare><transactionAcquiredDisposedCode><value>A</value></transactionAcquiredDisposedCode></transactionAmounts><postTransactionAmounts><sharesOwnedFollowingTransaction><value>1000</value></sharesOwnedFollowingTransaction></postTransactionAmounts>`
	xmlBody := `<?xml version="1.0"?><ownershipDocument><reportingOwner><reportingOwnerId><rptOwnerName>Jane Trader</rptOwnerName></reportingOwnerId><reportingOwnerRelationship><isOfficer>1</isOfficer><officerTitle>CFO</officerTitle></reportingOwnerRelationship></reportingOwner><nonDerivativeTable><nonDerivativeTransaction>` + tx + `</nonDerivativeTransaction><nonDerivativeTransaction>` + tx + `</nonDerivativeTransaction></nonDerivativeTable><derivativeTable><derivativeTransaction>` + tx + `</derivativeTransaction></derivativeTable></ownershipDocument>`
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { _, _ = w.Write([]byte(xmlBody)) }))
	defer srv.Close()
	item := FilingItem{Symbol: "AAPL", Form: "4", FiledAt: "2026-08-08", URL: srv.URL}
	enrichForm4(context.Background(), srv.Client(), nil, srv.URL, &item)
	if len(item.Transactions) != 2 {
		t.Fatalf("transactions=%d want 2 legitimate same-table purchases; %+v", len(item.Transactions), item.Transactions)
	}
	if item.Signal != "Buy" || item.Shares != 200 || item.Value != 2000 {
		t.Fatalf("wrong aggregate after dedupe: signal=%s shares=%v value=%v", item.Signal, item.Shares, item.Value)
	}
}

func TestV1512ResearchReadinessDoesNotMaskTargetSECFailureWithGlobalFreshness(t *testing.T) {
	app := newTestApplication(t)
	app.engine.seedDemo()
	app.engine.mu.Lock()
	app.engine.health["research-sec:NVDA"] = "setup required · SEC contact email"
	app.engine.lastUpdated["filings"] = time.Now().UnixMilli()
	app.engine.mu.Unlock()
	ready, issues := app.engine.researchPackageReadiness("NVDA")
	if ready {
		t.Fatalf("research incorrectly ready despite target SEC failure; issues=%v", issues)
	}
	found := false
	for _, issue := range issues {
		if strings.Contains(issue, "SEC Filings") {
			found = true
		}
	}
	if !found {
		t.Fatalf("expected SEC readiness issue, got %v", issues)
	}
}

func TestV1512StrictJSONRejectsTrailingPayload(t *testing.T) {
	var in struct {
		Symbol string `json:"symbol"`
	}
	req := httptest.NewRequest("POST", "/", strings.NewReader(`{"symbol":"AAPL"} {"extra":1}`))
	if err := decodeJSON(req, &in); err == nil {
		t.Fatal("decodeJSON accepted trailing JSON value")
	}
}

func TestV1512PostOnlyRejectsMutationGET(t *testing.T) {
	called := false
	h := postOnly(func(w http.ResponseWriter, r *http.Request) { called = true; w.WriteHeader(204) })
	rr := httptest.NewRecorder()
	h(rr, httptest.NewRequest(http.MethodGet, "/api/runtime/start", nil))
	if rr.Code != http.StatusMethodNotAllowed || called {
		t.Fatalf("GET mutation guard status=%d called=%v", rr.Code, called)
	}
}

func TestV1512LegacyRemoveRejectsMalformedTicker(t *testing.T) {
	app := newTestApplication(t)
	app.mu.Lock()
	if wl := findWatchlistInState(&app.state, "discovery"); wl != nil {
		wl.Symbols = []string{"AAPL"}
	}
	app.mu.Unlock()
	rr := httptest.NewRecorder()
	req := httptest.NewRequest("POST", "/api/watchlists/remove-symbol", strings.NewReader(`{"watchlistId":"discovery","symbol":"AA/PL"}`))
	req.Header.Set("Content-Type", "application/json")
	app.handleRemoveSymbol(rr, req)
	if rr.Code != 400 {
		t.Fatalf("malformed legacy remove status=%d body=%s", rr.Code, rr.Body.String())
	}
	app.mu.RLock()
	wl, _ := watchlistValueByID(app.state.Watchlists, "discovery")
	app.mu.RUnlock()
	if !contains(wl.Symbols, "AAPL") {
		t.Fatal("malformed ticker removed normalized AAPL")
	}
}

func TestV1512TargetedRefreshOutcomeSparseCheckVsDataAge(t *testing.T) {
	before := FreshnessDiagnostic{Dataset: "News", ProviderTimestamp: 1000, ReceivedAt: 2000, FreshnessBasis: "last successful news check", State: "FRESH"}
	after := before
	after.ReceivedAt = 3000
	changed, msg := targetedRefreshOutcome(before, after, true, true)
	if changed || !strings.Contains(msg, "Check Age reset") || !strings.Contains(msg, "Data Age") {
		t.Fatalf("sparse unchanged refresh outcome changed=%v msg=%q", changed, msg)
	}
	after.ProviderTimestamp = 4000
	changed, msg = targetedRefreshOutcome(before, after, true, true)
	if !changed || !strings.Contains(msg, "newer provider observation") {
		t.Fatalf("new provider refresh outcome changed=%v msg=%q", changed, msg)
	}
	changed, msg = targetedRefreshOutcome(before, after, true, false)
	if changed || !strings.Contains(strings.ToLower(msg), "failed") {
		t.Fatalf("failed refresh outcome changed=%v msg=%q", changed, msg)
	}
}

func TestV1512ScopedFreshnessUsesOldestCriticalObservation(t *testing.T) {
	app := newTestApplication(t)
	now := time.Now()
	fresh := now.Add(-20 * time.Second).UnixMilli()
	old := now.Add(-3 * time.Hour).UnixMilli()
	app.engine.mu.Lock()
	app.engine.quotes["AAPL"] = Quote{Symbol: "AAPL", Price: 100, ProviderTimestamp: fresh, UpdatedAt: fresh}
	app.engine.quotes["MSFT"] = Quote{Symbol: "MSFT", Price: 200, ProviderTimestamp: old, UpdatedAt: old}
	app.engine.bars["AAPL"] = map[string][]Bar{"intraday": {{T: fresh / 1000, C: 100}}}
	app.engine.bars["MSFT"] = map[string][]Bar{"intraday": {{T: old / 1000, C: 200}}}
	last := clone(app.engine.lastUpdated)
	health := clone(app.engine.health)
	quotes := clone(app.engine.quotes)
	app.engine.mu.Unlock()
	rows, _ := app.engine.buildFreshnessDiagnostics(quotes, last, health, []string{"AAPL", "MSFT"}, []string{"AAPL", "MSFT"})
	by := map[string]FreshnessDiagnostic{}
	for _, r := range rows {
		by[r.Dataset] = r
	}
	if got := by["Quotes"].ProviderTimestamp; got != old {
		t.Fatalf("Quotes aggregate used freshest timestamp %d; want oldest critical %d", got, old)
	}
	wantBarOld := (old / 1000) * 1000
	if got := by["Intraday Bars"].ProviderTimestamp; got != wantBarOld {
		t.Fatalf("Intraday aggregate used freshest timestamp %d; want oldest critical %d", got, wantBarOld)
	}
}

func TestV1512ScopedFreshnessMissingCriticalSymbolCannotLookCurrent(t *testing.T) {
	app := newTestApplication(t)
	now := time.Now().UnixMilli()
	app.engine.mu.Lock()
	app.engine.quotes["AAPL"] = Quote{Symbol: "AAPL", Price: 100, ProviderTimestamp: now, UpdatedAt: now}
	app.engine.bars["AAPL"] = map[string][]Bar{"intraday": {{T: now / 1000, C: 100}}}
	last := clone(app.engine.lastUpdated)
	health := clone(app.engine.health)
	quotes := clone(app.engine.quotes)
	app.engine.mu.Unlock()
	rows, _ := app.engine.buildFreshnessDiagnostics(quotes, last, health, []string{"AAPL", "MSFT"}, []string{"AAPL", "MSFT"})
	by := map[string]FreshnessDiagnostic{}
	for _, r := range rows {
		by[r.Dataset] = r
	}
	for _, name := range []string{"Quotes", "Intraday Bars"} {
		if by[name].State != "STALE" || !strings.Contains(strings.ToLower(by[name].Reason), "missing") {
			t.Fatalf("%s missing coverage incorrectly current: %+v", name, by[name])
		}
	}
}

func TestV1512FutureProviderTimestampFallsBackToReceipt(t *testing.T) {
	now := time.Now().UnixMilli()
	receipt := now - int64(time.Minute/time.Millisecond)
	future := now + int64(24*time.Hour/time.Millisecond)
	got, anomaly := safeFreshnessTimestamp(future, receipt, now)
	if !anomaly || got != receipt {
		t.Fatalf("future provider timestamp got=%d anomaly=%v want receipt=%d", got, anomaly, receipt)
	}
	got, anomaly = safeFreshnessTimestamp(future, 0, now)
	if !anomaly || got != 0 {
		t.Fatalf("future timestamp without receipt got=%d anomaly=%v", got, anomaly)
	}
}

func TestV1512ResearchTargetReadyDespiteUnrelatedUniverseStaleness(t *testing.T) {
	app := newTestApplication(t)
	sym := "AAA"
	now := time.Now()
	ms := now.UnixMilli()
	app.engine.mu.Lock()
	app.engine.quotes[sym] = Quote{Symbol: sym, Price: 100, ProviderTimestamp: ms, UpdatedAt: ms}
	app.engine.quotes["OLDX"] = Quote{Symbol: "OLDX", Price: 10, ProviderTimestamp: now.Add(-24 * time.Hour).UnixMilli(), UpdatedAt: now.Add(-24 * time.Hour).UnixMilli()}
	app.engine.bars[sym] = map[string][]Bar{"intraday": {{T: now.Unix(), C: 100}}, "daily": {{T: now.Add(-24 * time.Hour).Unix(), C: 99}}}
	app.engine.fundamentals[sym] = FundamentalSnapshot{Symbol: sym, RevenueGrowth: 5, UpdatedAt: ms, Source: "Finnhub"}
	for _, k := range []string{"research-news:" + sym, "research-earnings:" + sym, "research-fundamentals:" + sym, "research-sec:" + sym} {
		app.engine.lastUpdated[k] = ms
	}
	app.engine.health["research-sec:"+sym] = "healthy · SEC EDGAR"
	app.engine.mu.Unlock()
	ready, issues := app.engine.researchPackageReadiness(sym)
	if !ready || len(issues) != 0 {
		t.Fatalf("unrelated stale universe masked target-scoped research: ready=%v issues=%v", ready, issues)
	}
}

func TestV1512StopTimeoutRemainsStoppingUntilWorkersExit(t *testing.T) {
	app := newTestApplication(t)
	oldTimeout := runtimeStopTimeout
	runtimeStopTimeout = 20 * time.Millisecond
	defer func() { runtimeStopTimeout = oldTimeout }()
	app.engine.mu.Lock()
	app.engine.status = "running"
	app.engine.cancel = func() {}
	app.engine.mu.Unlock()
	app.engine.wg.Add(1)
	app.engine.Stop()
	app.engine.mu.RLock()
	status := app.engine.status
	app.engine.mu.RUnlock()
	if status != "stopping" {
		t.Fatalf("timed-out stop falsely reported %q instead of stopping", status)
	}
	if err := app.engine.Start(); err == nil {
		t.Fatal("Start succeeded while prior workers were still stopping")
	}
	app.engine.wg.Done()
	deadline := time.Now().Add(time.Second)
	for time.Now().Before(deadline) {
		app.engine.mu.RLock()
		status = app.engine.status
		app.engine.mu.RUnlock()
		if status == "stopped" {
			return
		}
		time.Sleep(5 * time.Millisecond)
	}
	t.Fatalf("runtime did not finalize stop after workers exited; status=%q", status)
}

func TestV1512AIRejectsMalformedTickerBeforeProviderCall(t *testing.T) {
	app := newTestApplication(t)
	_, err := app.GenerateAIForUser(context.Background(), bootstrapOwnerID, AIRequest{Kind: "ticker", Ticker: "AA/PL"})
	if err == nil || !strings.Contains(strings.ToLower(err.Error()), "invalid ticker") {
		t.Fatalf("malformed AI ticker error=%v", err)
	}
}

func TestV1512Form4PreservesAllReportingOwners(t *testing.T) {
	tx := `<transactionDate><value>2026-08-08</value></transactionDate><transactionCoding><transactionCode>P</transactionCode></transactionCoding><transactionAmounts><transactionShares><value>50</value></transactionShares><transactionPricePerShare><value>20</value></transactionPricePerShare><transactionAcquiredDisposedCode><value>A</value></transactionAcquiredDisposedCode></transactionAmounts>`
	xmlBody := `<?xml version="1.0"?><ownershipDocument>` +
		`<reportingOwner><reportingOwnerId><rptOwnerName>Jane Trader</rptOwnerName></reportingOwnerId><reportingOwnerRelationship><isOfficer>1</isOfficer><officerTitle>CFO</officerTitle></reportingOwnerRelationship></reportingOwner>` +
		`<reportingOwner><reportingOwnerId><rptOwnerName>John Director</rptOwnerName></reportingOwnerId><reportingOwnerRelationship><isDirector>1</isDirector></reportingOwnerRelationship></reportingOwner>` +
		`<nonDerivativeTable><nonDerivativeTransaction>` + tx + `</nonDerivativeTransaction></nonDerivativeTable></ownershipDocument>`
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { _, _ = w.Write([]byte(xmlBody)) }))
	defer srv.Close()
	item := FilingItem{Symbol: "AAPL", Form: "4", FiledAt: "2026-08-08", URL: srv.URL}
	enrichForm4(context.Background(), srv.Client(), nil, srv.URL, &item)
	if !strings.Contains(item.Actor, "Jane Trader") || !strings.Contains(item.Actor, "John Director") {
		t.Fatalf("reporting owners lost: actor=%q", item.Actor)
	}
	if !strings.Contains(item.Role, "CFO") || !strings.Contains(item.Role, "Director") {
		t.Fatalf("reporting owner roles lost: role=%q", item.Role)
	}
	if len(item.Transactions) != 1 || item.Transactions[0].Actor != item.Actor {
		t.Fatalf("transaction did not retain all reporting owners: %+v", item.Transactions)
	}
}

func TestV1512SignalValidationRejectsMalformedTickerInsteadOfNormalizing(t *testing.T) {
	app := newTestApplication(t)
	for _, body := range []string{
		`{"symbol":"AA/PL","horizon":"day"}`,
		`{"symbol":"VIX","horizon":"day"}`,
		`{"symbol":"AAPL<script>","horizon":"day"}`,
	} {
		rr := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPost, "/api/signal-validation/record", strings.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		app.handleSignalValidationRecord(rr, req)
		if rr.Code != http.StatusBadRequest {
			t.Fatalf("malformed signal accepted body=%s status=%d response=%s", body, rr.Code, rr.Body.String())
		}
	}

	rr := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/api/signal-validation/record", strings.NewReader(`{"symbol":" aapl ","horizon":" DAY "}`))
	req.Header.Set("Content-Type", "application/json")
	app.handleSignalValidationRecord(rr, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("valid normalized signal rejected: %d %s", rr.Code, rr.Body.String())
	}
	var out struct {
		Snapshot SignalSnapshot `json:"snapshot"`
	}
	if err := json.Unmarshal(rr.Body.Bytes(), &out); err != nil {
		t.Fatal(err)
	}
	if out.Snapshot.Symbol != "AAPL" || out.Snapshot.Horizon != "day" {
		t.Fatalf("signal not canonically normalized: %+v", out.Snapshot)
	}
}

func TestV1512GlobalRemoveSelectedTickerActuallyLeavesMasterTracking(t *testing.T) {
	app := newTestApplication(t)
	const sym = "ZZZZ"
	app.mu.Lock()
	for _, id := range deskIDs() {
		setMembershipLocked(&app.state, id, sym, true)
	}
	app.state.UI.SelectedTicker = sym
	app.mu.Unlock()

	rr := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/api/master-symbol/remove", strings.NewReader(`{"symbol":"ZZZZ"}`))
	req.Header.Set("Content-Type", "application/json")
	app.handleMasterSymbolRemove(rr, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("remove status=%d body=%s", rr.Code, rr.Body.String())
	}
	app.mu.RLock()
	st := clone(app.state)
	app.mu.RUnlock()
	if normalizeSymbol(st.UI.SelectedTicker) == sym {
		t.Fatalf("removed ticker remained selected: %q", st.UI.SelectedTicker)
	}
	for _, got := range masterSymbolsFromState(st) {
		if got == sym {
			t.Fatalf("removed ticker remained in master tracking: %v", masterSymbolsFromState(st))
		}
	}
	if activeDeskCount(deskMembershipsLocked(&st, sym)) != 0 {
		t.Fatalf("removed ticker retained desk membership: %+v", deskMembershipsLocked(&st, sym))
	}
}

func TestV1512ManualPrepRunsDoNotConsumeAutomaticRetryBudget(t *testing.T) {
	app := newTestApplication(t)
	e := app.engine
	for i := 0; i < 5; i++ {
		e.setPreparationRich("pre-market-prep", "RUNNING", "manual", false, false, "", nil, nil, nil)
		e.setPreparationRich("pre-market-prep", "DATA DEGRADED", "manual failed", false, false, "DATA DEGRADED", nil, nil, nil)
	}
	e.mu.RLock()
	p := e.preparations["pre-market-prep"]
	e.mu.RUnlock()
	if p.AttemptCount != 0 {
		t.Fatalf("manual runs consumed automatic retry budget: %d", p.AttemptCount)
	}
	if preparationRanForTradingDay(p, time.Now()) {
		t.Fatal("manual failures incorrectly suppressed future automatic checkpoint")
	}

	for i := 0; i < 3; i++ {
		e.setPreparationRich("pre-market-prep", "RUNNING", "scheduled", false, false, "", nil, nil, nil)
		e.setPreparationRich("pre-market-prep", "DATA DEGRADED", "scheduled failed", false, false, "DATA DEGRADED", nil, nil, nil)
	}
	e.mu.RLock()
	p = e.preparations["pre-market-prep"]
	e.mu.RUnlock()
	if p.AttemptCount != 3 || !preparationRanForTradingDay(p, time.Now()) {
		t.Fatalf("automatic retry cap not enforced: attempts=%d status=%v", p.AttemptCount, preparationRanForTradingDay(p, time.Now()))
	}
}

func TestV1512LateRegularSessionCatalystGetsFullNextSessionLifecycle(t *testing.T) {
	loc := easternLocation()
	// Monday 3:30 PM ET leaves only 30 minutes to the regular close, so a promised
	// 60-minute reaction checkpoint cannot finish that day.
	trigger := time.Date(2026, 8, 10, 15, 30, 0, 0, loc)
	until := catalystCompletionAt(trigger.UnixMilli())
	wantDay := "2026-08-11"
	if got := until.In(loc).Format("2006-01-02"); got != wantDay {
		t.Fatalf("late catalyst completed too early: got=%s until=%s", got, until.In(loc))
	}
	if until.In(loc).Hour() != 16 || until.In(loc).Minute() != 15 {
		t.Fatalf("unexpected next-session completion=%s", until.In(loc))
	}

	// A trigger with a full hour remaining may complete in the same session.
	early := time.Date(2026, 8, 10, 14, 59, 0, 0, loc)
	until = catalystCompletionAt(early.UnixMilli())
	if got := until.In(loc).Format("2006-01-02"); got != "2026-08-10" {
		t.Fatalf("sufficiently early catalyst unexpectedly rolled forward: %s", until.In(loc))
	}
}

func TestV1512CBOEVIXOfficialReferenceRemainsDelayedOutsideRegularSession(t *testing.T) {
	loc := easternLocation()
	cases := []struct {
		name string
		now  time.Time
		age  time.Duration
	}{
		{"premarket", time.Date(2026, 8, 10, 8, 0, 0, 0, loc), 18 * time.Hour},
		{"overnight", time.Date(2026, 8, 10, 2, 0, 0, 0, loc), 12 * time.Hour},
		{"weekend", time.Date(2026, 8, 9, 12, 0, 0, 0, loc), 44 * time.Hour},
	}
	for _, tc := range cases {
		nowMs := tc.now.UnixMilli()
		ts := tc.now.Add(-tc.age).UnixMilli()
		state, reason, _, _, _ := freshnessState("VIX", "CBOE", marketSessionET(tc.now), ts, "healthy", nowMs)
		if state != "DELAYED" || !strings.Contains(reason, "not live") {
			t.Fatalf("%s CBOE reference state=%s reason=%q", tc.name, state, reason)
		}
	}
}

func TestV1512FreshnessUnknownClocksDoNotPretendToBeZeroSecondsOld(t *testing.T) {
	app := newTestApplication(t)
	app.engine.mu.Lock()
	app.engine.news = nil
	app.engine.earnings = nil
	app.engine.filings = nil
	app.engine.fundamentals = map[string]FundamentalSnapshot{}
	app.engine.options = map[string]OptionsContext{}
	for _, k := range []string{"news", "earnings", "filings", "fundamentals", "options"} {
		delete(app.engine.lastUpdated, k)
		app.engine.health[k] = "unavailable"
	}
	last := clone(app.engine.lastUpdated)
	health := clone(app.engine.health)
	quotes := clone(app.engine.quotes)
	app.engine.mu.Unlock()
	rows, _ := app.engine.buildFreshnessDiagnostics(quotes, last, health)
	by := map[string]FreshnessDiagnostic{}
	for _, row := range rows {
		by[row.Dataset] = row
	}
	for _, name := range []string{"News", "Earnings", "SEC Filings", "Fundamentals", "Options"} {
		row := by[name]
		if row.CheckAgeMs != -1 {
			t.Fatalf("%s missing check clock displayed as age=%d instead of unknown", name, row.CheckAgeMs)
		}
		if row.DataAgeMs != -1 {
			t.Fatalf("%s missing data clock displayed as age=%d instead of unknown", name, row.DataAgeMs)
		}
	}
}

func TestV1512MalformedProviderNewsTimestampStaysUnknown(t *testing.T) {
	if got := parseRFC3339Unix("not-a-timestamp"); got != 0 {
		t.Fatalf("malformed provider time became %d instead of unknown", got)
	}
	if got := parseRFC3339Unix("2026-08-10T16:30:00Z"); got <= 0 {
		t.Fatalf("valid RFC3339 time rejected: %d", got)
	}
}

func TestV1512DeskMembershipExhaustiveStateTransitions(t *testing.T) {
	desks := []string{"day", "swing", "long"}
	for mask := 1; mask < 8; mask++ {
		for _, clicked := range desks {
			app := newTestApplication(t)
			const sym = "EDGE"
			app.mu.Lock()
			for i, desk := range desks {
				setMembershipLocked(&app.state, desk, sym, mask&(1<<i) != 0)
			}
			before := deskMembershipsLocked(&app.state, sym)
			app.mu.Unlock()
			wasActive := before[clicked]
			activeBefore := activeDeskCount(before)

			rr := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodPost, "/api/desk/membership", strings.NewReader(`{"symbol":"EDGE","desk":"`+clicked+`"}`))
			req.Header.Set("Content-Type", "application/json")
			app.handleDeskMembership(rr, req)
			if rr.Code != http.StatusOK {
				t.Fatalf("mask=%03b click=%s status=%d body=%s", mask, clicked, rr.Code, rr.Body.String())
			}
			app.mu.RLock()
			after := deskMembershipsLocked(&app.state, sym)
			app.mu.RUnlock()
			if !wasActive && !after[clicked] {
				t.Fatalf("mask=%03b inactive %s did not add: %+v", mask, clicked, after)
			}
			if wasActive && activeBefore > 1 && after[clicked] {
				t.Fatalf("mask=%03b active %s did not remove with other desk remaining: %+v", mask, clicked, after)
			}
			if wasActive && activeBefore == 1 && !after[clicked] {
				t.Fatalf("mask=%03b final %s membership was removed: %+v", mask, clicked, after)
			}
			if activeDeskCount(after) < 1 {
				t.Fatalf("mask=%03b click=%s left zero desks: %+v", mask, clicked, after)
			}
		}
	}
}
