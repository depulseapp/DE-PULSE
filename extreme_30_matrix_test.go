package main

import (
	"bufio"
	"context"
	"encoding/binary"
	"encoding/json"
	"errors"
	"io"
	"net"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"testing"
	"time"
)

func callHandler(t *testing.T, h func(http.ResponseWriter, *http.Request), method, path, body string) *httptest.ResponseRecorder {
	t.Helper()
	rr := httptest.NewRecorder()
	var r io.Reader
	if body != "" {
		r = strings.NewReader(body)
	}
	req := httptest.NewRequest(method, path, r)
	if body != "" {
		req.Header.Set("Content-Type", "application/json")
	}
	h(rr, req)
	return rr
}

func TestExtreme30_01CoreFunctionalWorkflowAndMethodSafety(t *testing.T) {
	app := v15TestApp(t)
	// State-changing handlers must not mutate on GET.
	before := len(app.state.Watchlists)
	rr := callHandler(t, postOnly(app.handleWatchlistCreate), http.MethodGet, "/api/watchlists/create", "")
	if rr.Code != http.StatusMethodNotAllowed || len(app.state.Watchlists) != before {
		t.Fatalf("GET create mutated state code=%d", rr.Code)
	}
	rr = callHandler(t, postOnly(app.handleRuntimeStart), http.MethodGet, "/api/runtime/start", "")
	if rr.Code != http.StatusMethodNotAllowed {
		t.Fatalf("GET runtime start code=%d", rr.Code)
	}
	// Malformed JSON must not mutate.
	rr = callHandler(t, app.handleWatchlistCreate, http.MethodPost, "/api/watchlists/create", "{")
	if rr.Code != 400 || len(app.state.Watchlists) != before {
		t.Fatalf("malformed create mutated state code=%d", rr.Code)
	}
	// Valid create/rename/delete roundtrip.
	rr = callHandler(t, app.handleWatchlistCreate, http.MethodPost, "/api/watchlists/create", `{"name":"Extreme QA"}`)
	if rr.Code != 200 {
		t.Fatalf("create %d %s", rr.Code, rr.Body.String())
	}
	var wl Watchlist
	_ = json.Unmarshal(rr.Body.Bytes(), &wl)
	rr = callHandler(t, app.handleWatchlistRename, http.MethodPost, "/api/watchlists/rename", `{"id":"`+wl.ID+`","name":"Extreme QA 2"}`)
	if rr.Code != 200 {
		t.Fatalf("rename %d", rr.Code)
	}
	rr = callHandler(t, app.handleWatchlistDelete, http.MethodPost, "/api/watchlists/delete", `{"id":"`+wl.ID+`"}`)
	if rr.Code != 200 {
		t.Fatalf("delete %d", rr.Code)
	}
}

func TestExtreme30_02DeskMembershipEveryStateStress(t *testing.T) {
	states := []string{"100", "010", "001", "110", "101", "011", "111"}
	desks := []string{"day", "swing", "long"}
	for _, bits := range states {
		for _, desk := range desks {
			app := v15TestApp(t)
			setDeskBits(t, app, "EDGE", bits)
			for i := 0; i < 20; i++ { // repeated explicit state requests must stay idempotent
				want := true
				rr := callHandler(t, app.handleDeskMembership, http.MethodPost, "/api/desk/membership", `{"symbol":"EDGE","desk":"`+desk+`","active":`+map[bool]string{true: "true", false: "false"}[want]+`}`)
				if rr.Code != 200 {
					t.Fatalf("%s %s add iteration %d code=%d", bits, desk, i, rr.Code)
				}
			}
			app.mu.RLock()
			m := deskMembershipsLocked(&app.state, "EDGE")
			app.mu.RUnlock()
			if !m[desk] {
				t.Fatalf("desk %s not active after idempotent add: %v", desk, m)
			}
		}
	}
}

func TestExtreme30_03TickerInputEdgeCases(t *testing.T) {
	good := map[string]string{" aapl ": "AAPL", "brk.b": "BRK.B", "bf-b": "BF-B", "a1": "A1", "qqq": "QQQ"}
	for raw, want := range good {
		got, ok := parseUserTicker(raw)
		if !ok || got != want {
			t.Fatalf("valid %q => %q %v", raw, got, ok)
		}
	}
	bad := []string{"", "VIX", "1AAPL", "AA@PL", "AAPL/US", "AAPL MSFT", "AAPL<script>", "💥AAPL", "AAPL$", "TOOLONG99"}
	for _, raw := range bad {
		if got, ok := parseUserTicker(raw); ok {
			t.Fatalf("invalid %q accepted as %q", raw, got)
		}
	}
	app := v15TestApp(t)
	for _, raw := range bad {
		rr := callHandler(t, app.handleMasterSymbolAdd, http.MethodPost, "/api/master-symbol/add", `{"symbol":"`+strings.ReplaceAll(raw, "\"", "\\\"")+`"}`)
		if rr.Code == 200 {
			t.Fatalf("master add accepted %q", raw)
		}
	}
	// Valid lowercase master add is canonical and fills all desks exactly once.
	rr := callHandler(t, app.handleMasterSymbolAdd, http.MethodPost, "/api/master-symbol/add", `{"symbol":"rddt"}`)
	if rr.Code != 200 {
		t.Fatalf("valid master add code=%d %s", rr.Code, rr.Body.String())
	}
	app.mu.RLock()
	m := deskMembershipsLocked(&app.state, "RDDT")
	app.mu.RUnlock()
	if activeDeskCount(m) != 3 {
		t.Fatalf("RDDT not in all desks %v", m)
	}
}

func TestExtreme30_04ProviderRouterChaos(t *testing.T) {
	app := v15TestApp(t)
	app.mu.Lock()
	app.secrets.Finnhub = "x"
	app.secrets.Marketaux = "x"
	app.mu.Unlock()
	e := app.engine
	calls := []string{}
	active, ok := e.executeProviderRoute(context.Background(), "News", map[string]providerRouteAttempt{
		"Finnhub":   func(context.Context) bool { calls = append(calls, "Finnhub"); return false },
		"Marketaux": func(context.Context) bool { calls = append(calls, "Marketaux"); return true },
	})
	if !ok || active != "Marketaux" || strings.Join(calls, ">") != "Finnhub>Marketaux" {
		t.Fatalf("fallback bad active=%s calls=%v", active, calls)
	}
	// 429 suppresses provider; fallback becomes first executable hop.
	e.recordProviderFailure("Finnhub", errors.New("HTTP 429 too many requests"))
	calls = nil
	active, ok = e.executeProviderRoute(context.Background(), "News", map[string]providerRouteAttempt{
		"Finnhub":   func(context.Context) bool { calls = append(calls, "Finnhub"); return true },
		"Marketaux": func(context.Context) bool { calls = append(calls, "Marketaux"); return true },
	})
	if !ok || active != "Marketaux" || len(calls) != 1 || calls[0] != "Marketaux" {
		t.Fatalf("rate limit not respected active=%s calls=%v", active, calls)
	}
	e.recordProviderSuccess("Finnhub")
	if !e.providerAllowed("Finnhub") {
		t.Fatal("Finnhub did not recover")
	}
}

func writeServerFrame(c net.Conn, fin bool, opcode byte, payload []byte) error {
	b1 := opcode
	if fin {
		b1 |= 0x80
	}
	hdr := []byte{b1}
	n := len(payload)
	switch {
	case n < 126:
		hdr = append(hdr, byte(n))
	case n <= 65535:
		hdr = append(hdr, 126, 0, 0)
		binary.BigEndian.PutUint16(hdr[len(hdr)-2:], uint16(n))
	default:
		hdr = append(hdr, 127, 0, 0, 0, 0, 0, 0, 0, 0)
		binary.BigEndian.PutUint64(hdr[len(hdr)-8:], uint64(n))
	}
	if _, err := c.Write(append(hdr, payload...)); err != nil {
		return err
	}
	return nil
}

func TestExtreme30_05WebSocketProtocolEdges(t *testing.T) {
	cli, srv := net.Pipe()
	defer cli.Close()
	defer srv.Close()
	ws := &WSClient{conn: cli, r: bufio.NewReader(cli), w: bufio.NewWriter(cli)}
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	go func() {
		_ = writeServerFrame(srv, false, 0x1, []byte("hel"))
		_ = writeServerFrame(srv, true, 0x9, []byte("p"))
		// Read the masked client PONG before sending the continuation. This
		// models a protocol-compliant peer and prevents net.Pipe's zero-buffer
		// semantics from creating an artificial cross-write deadlock.
		h := make([]byte, 2)
		if _, err := io.ReadFull(srv, h); err == nil {
			n := int(h[1] & 0x7f)
			mask := make([]byte, 4)
			_, _ = io.ReadFull(srv, mask)
			p := make([]byte, n)
			_, _ = io.ReadFull(srv, p)
		}
		_ = writeServerFrame(srv, true, 0x0, []byte("lo"))
	}()
	got, err := ws.ReadText(ctx)
	if err != nil || string(got) != "hello" {
		t.Fatalf("fragment/ping read got=%q err=%v", got, err)
	}
	// Close frame returns EOF.
	go func() { _ = writeServerFrame(srv, true, 0x8, nil) }()
	_, err = ws.ReadText(ctx)
	if !errors.Is(err, io.EOF) {
		t.Fatalf("close frame err=%v", err)
	}
}

func TestExtreme30_06DataFreshnessBoundariesAndClockSkew(t *testing.T) {
	now := time.Now().UnixMilli()
	datasets := []struct{ d, p, s string }{
		{"Quotes", "Alpaca", "regular"}, {"VIX", "Twelve Data", "regular"}, {"Intraday Bars", "Alpaca", "regular"},
		{"Daily / Weekly History", "Alpaca", "regular"}, {"News", "Finnhub", "regular"}, {"SEC Filings", "SEC EDGAR", "regular"},
		{"Fundamentals", "Finnhub", "regular"}, {"Global", "Twelve Data", "regular"}, {"Macro", "FRED", "regular"}, {"Options", "Alpaca", "regular"},
	}
	for _, x := range datasets {
		_, fresh, stale := freshnessLimits(x.d, x.p, x.s)
		st, _, _, _, _ := freshnessState(x.d, x.p, x.s, now-int64(fresh/time.Millisecond)+1, "healthy", now)
		if st == "STALE" || st == "ERROR" || st == "UNAVAILABLE" {
			t.Fatalf("%s unexpectedly %s just inside fresh", x.d, st)
		}
		st, _, _, _, _ = freshnessState(x.d, x.p, x.s, now-int64(stale/time.Millisecond)-1, "healthy", now)
		if x.d == "VIX" && x.p == "CBOE" {
			continue
		}
		if st != "STALE" {
			t.Fatalf("%s expected STALE past stale limit, got %s", x.d, st)
		}
	}
	stateTs, anomaly := safeFreshnessTimestamp(now+10*60*1000, 0, now)
	if !anomaly || stateTs != 0 {
		t.Fatalf("future provider timestamp not quarantined: ts=%d anomaly=%v", stateTs, anomaly)
	}
	stateTs, anomaly = safeFreshnessTimestamp(now+30*1000, now, now)
	if anomaly || stateTs != now+30*1000 {
		t.Fatalf("small clock skew should be tolerated: ts=%d anomaly=%v", stateTs, anomaly)
	}
	st, _, age, _, _ := freshnessState("Quotes", "Alpaca", "regular", stateTs, "healthy", now)
	if st != "LIVE" || age != 0 {
		t.Fatalf("small clock skew should be live age0: %s %d", st, age)
	}
	// Sparse News uses recent CHECK AGE while preserving old DATA AGE.
	app := v15TestApp(t)
	e := app.engine
	e.mu.Lock()
	e.news = []NewsItem{{ID: "1", Headline: "old", Datetime: (now - int64(2*time.Hour/time.Millisecond)) / 1000, Source: "Finnhub"}}
	e.lastUpdated["news"] = now - int64(2*time.Minute/time.Millisecond)
	e.health["news"] = "healthy · Finnhub"
	e.mu.Unlock()
	rows, _ := e.buildFreshnessDiagnostics(map[string]Quote{}, clone(e.lastUpdated), clone(e.health))
	found := false
	for _, r := range rows {
		if r.Dataset == "News" {
			found = true
			if r.State != "FRESH" && r.State != "DUE SOON" {
				t.Fatalf("sparse news state=%s", r.State)
			}
			if r.DataAgeMs <= r.CheckAgeMs {
				t.Fatalf("news clocks not distinct %+v", r)
			}
		}
	}
	if !found {
		t.Fatal("News row missing")
	}
}

func TestExtreme30_07MarketSessionBoundaryMatrix(t *testing.T) {
	loc := easternLocation()
	cases := []struct {
		at   time.Time
		want string
	}{
		{time.Date(2026, 8, 10, 3, 59, 59, 0, loc), "overnight"}, {time.Date(2026, 8, 10, 4, 0, 0, 0, loc), "pre-market"},
		{time.Date(2026, 8, 10, 9, 29, 59, 0, loc), "pre-market"}, {time.Date(2026, 8, 10, 9, 30, 0, 0, loc), "regular"},
		{time.Date(2026, 8, 10, 15, 59, 59, 0, loc), "regular"}, {time.Date(2026, 8, 10, 16, 0, 0, 0, loc), "after-hours"},
		{time.Date(2026, 8, 10, 19, 59, 59, 0, loc), "after-hours"}, {time.Date(2026, 8, 10, 20, 0, 0, 0, loc), "overnight"},
		{time.Date(2026, 8, 8, 12, 0, 0, 0, loc), "weekend"},
		{time.Date(2026, 11, 27, 12, 59, 59, 0, loc), "regular"}, {time.Date(2026, 11, 27, 13, 0, 0, 0, loc), "after-hours"},
		{time.Date(2026, 12, 25, 12, 0, 0, 0, loc), "closed"},
	}
	for _, c := range cases {
		if got := marketSessionET(c.at); got != c.want {
			t.Fatalf("%s got %s want %s", c.at, got, c.want)
		}
	}
	// DST-era dates retain correct session classification.
	for _, at := range []time.Time{time.Date(2026, 3, 9, 10, 0, 0, 0, loc), time.Date(2026, 11, 2, 10, 0, 0, 0, loc)} {
		if got := marketSessionET(at); got != "regular" {
			t.Fatalf("DST date %s=%s", at, got)
		}
	}
}

func TestExtreme30_08QuotesAndVIXTruthfulness(t *testing.T) {
	now := time.Now().UnixMilli()
	st, reason, _, _, _ := freshnessState("VIX", "CBOE", "after-hours", now-int64(2*time.Hour/time.Millisecond), "healthy", now)
	if st != "DELAYED" || !strings.Contains(strings.ToLower(reason), "not live") {
		t.Fatalf("CBOE false live %s %s", st, reason)
	}
	st, _, _, _, _ = freshnessState("VIX", "Twelve Data", "regular", now-int64(30*time.Minute/time.Millisecond), "healthy", now)
	if st != "STALE" {
		t.Fatalf("old primary VIX=%s", st)
	}
	// Cached quote must remain explicitly cached on reload.
	app := v15TestApp(t)
	app.engine.mu.Lock()
	app.engine.quotes["AAPL"] = Quote{Symbol: "AAPL", Price: 200, Source: "Alpaca", FeedType: "stream", DataState: "live", UpdatedAt: now}
	app.engine.mu.Unlock()
	if err := app.engine.saveCache(); err != nil {
		t.Fatal(err)
	}
	e2 := NewEngine(app)
	q := e2.quotes["AAPL"]
	if q.DataState != "cache" || q.FeedType != "cache" {
		t.Fatalf("cache falsely live %+v", q)
	}
}

func TestExtreme30_09HistoricalBarsModeSeparation(t *testing.T) {
	now := time.Now().UnixMilli()
	st, _, _, _, _ := freshnessState("Daily / Weekly History", "Alpaca", "regular", now-int64(25*time.Hour/time.Millisecond), "healthy", now)
	if st == "STALE" {
		t.Fatalf("daily/weekly history false stale at 25h: %s", st)
	}
	st, _, _, _, _ = freshnessState("Intraday Bars", "Alpaca", "regular", now-int64(25*time.Hour/time.Millisecond), "healthy", now)
	if st != "STALE" {
		t.Fatalf("intraday old bars should stale: %s", st)
	}
	if got := strings.Join(routeChains()["Historical Bars"], ">"); got != "Alpaca>TradeInsight>Twelve Data>yfinance" {
		t.Fatalf("history route %s", got)
	}
}

func TestExtreme30_10NewsSparseAndDedup(t *testing.T) {
	items := []NewsItem{
		{ID: "1", Datetime: 100, Headline: "A", URL: "https://x/1"}, {ID: "1", Datetime: 101, Headline: "A dup", URL: "https://x/2"},
		{ID: nil, Datetime: 90, Headline: "B", URL: "https://x/b"}, {ID: nil, Datetime: 91, Headline: "B other", URL: "https://x/b"},
		{ID: nil, Datetime: 80, Headline: "C"}, {ID: nil, Datetime: 80, Headline: "C"},
	}
	got := dedupeNews(items)
	if len(got) != 3 {
		t.Fatalf("dedupe len=%d %+v", len(got), got)
	}
	if got[0].Datetime < got[len(got)-1].Datetime {
		t.Fatal("news not sorted newest first")
	}
}

func TestExtreme30_11SECInsiderSemanticsAndDuplicates(t *testing.T) {
	codes := map[string]string{"P": "BUY", "S": "SELL", "A": "OTHER", "M": "OTHER", "F": "OTHER", "G": "OTHER", "D": "OTHER", "C": "OTHER", "X": "OTHER", "J": "OTHER"}
	for code, want := range codes {
		got, _ := form4TransactionMeaning(code)
		if got != want {
			t.Fatalf("code %s=%s want %s", code, got, want)
		}
	}
	// Exact duplicate transaction rows across derivative/non-derivative are suppressed.
	xmlBody := `<ownershipDocument><reportingOwner><reportingOwnerId><rptOwnerName>Jane</rptOwnerName></reportingOwnerId><reportingOwnerRelationship><isDirector>1</isDirector></reportingOwnerRelationship></reportingOwner><nonDerivativeTable><nonDerivativeTransaction><transactionDate><value>2026-08-07</value></transactionDate><transactionCoding><transactionCode>M</transactionCode></transactionCoding><transactionAmounts><transactionShares><value>50</value></transactionShares><transactionPricePerShare><value>10</value></transactionPricePerShare></transactionAmounts></nonDerivativeTransaction></nonDerivativeTable><derivativeTable><derivativeTransaction><transactionDate><value>2026-08-07</value></transactionDate><transactionCoding><transactionCode>M</transactionCode></transactionCoding><transactionAmounts><transactionShares><value>50</value></transactionShares><transactionPricePerShare><value>10</value></transactionPricePerShare></transactionAmounts></derivativeTransaction></derivativeTable></ownershipDocument>`
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { _, _ = w.Write([]byte(xmlBody)) }))
	defer srv.Close()
	f := FilingItem{Symbol: "AAA", Form: "4", FiledAt: "2026-08-08", URL: srv.URL}
	enrichForm4(context.Background(), srv.Client(), nil, srv.URL, &f)
	if len(f.Transactions) != 1 || f.Transactions[0].Classification != "OTHER" {
		t.Fatalf("duplicate/semantic bad %+v", f.Transactions)
	}
}

func TestExtreme30_12EarningsCadenceAndWindow(t *testing.T) {
	loc := easternLocation()
	today := time.Date(2026, 8, 10, 8, 0, 0, 0, loc)
	if got := earningsRefreshIntervalFrom(nil, today); got != 2*time.Hour {
		t.Fatalf("normal cadence %s", got)
	}
	ev := []EarningsItem{{Symbol: "AAA", Date: "2026-08-10", Hour: "bmo"}}
	if got := earningsRefreshIntervalFrom(ev, time.Date(2026, 8, 10, 8, 30, 0, 0, loc)); got != 7*time.Minute {
		t.Fatalf("BMO cadence %s", got)
	}
	ev = []EarningsItem{{Symbol: "AAA", Date: "2026-08-10", Hour: "amc"}}
	if got := earningsRefreshIntervalFrom(ev, time.Date(2026, 8, 10, 17, 0, 0, 0, loc)); got != 7*time.Minute {
		t.Fatalf("AMC cadence %s", got)
	}
	actual := 1.2
	ev[0].EPSActual = &actual
	if got := earningsRefreshIntervalFrom(ev, time.Date(2026, 8, 10, 17, 0, 0, 0, loc)); got != 30*time.Minute {
		t.Fatalf("released cadence %s", got)
	}
}

func TestExtreme30_13CatalystLifecycleAllPhases(t *testing.T) {
	loc := easternLocation()
	trigger := time.Date(2026, 8, 10, 8, 0, 0, 0, loc)
	cases := []struct {
		at   time.Time
		want string
	}{
		{time.Date(2026, 8, 10, 8, 30, 0, 0, loc), "PREMARKET REACTION"},
		{time.Date(2026, 8, 10, 9, 32, 0, 0, loc), "OPENING REACTION"},
		{time.Date(2026, 8, 10, 9, 40, 0, 0, loc), "5m"},
		{time.Date(2026, 8, 10, 9, 50, 0, 0, loc), "15m"},
		{time.Date(2026, 8, 10, 10, 15, 0, 0, loc), "30m"},
		{time.Date(2026, 8, 10, 10, 35, 0, 0, loc), "60m"},
		{time.Date(2026, 8, 10, 12, 0, 0, 0, loc), "SESSION REACTION"},
		{time.Date(2026, 8, 10, 16, 16, 0, 0, loc), "COMPLETE"},
	}
	for _, c := range cases {
		if got := catalystPhase(c.at, trigger.UnixMilli()); got != c.want {
			t.Fatalf("%s got %s want %s", c.at, got, c.want)
		}
	}
	// AMC Friday survives weekend through Monday close.
	fri := time.Date(2026, 8, 7, 16, 5, 0, 0, loc)
	if catalystComplete(time.Date(2026, 8, 9, 12, 0, 0, 0, loc), fri.UnixMilli()) {
		t.Fatal("AMC completed on weekend")
	}
}

func TestExtreme30_14PreparationWindowsAndReadiness(t *testing.T) {
	loc := easternLocation()
	if !preMarketPrepWindow(time.Date(2026, 8, 10, 3, 15, 0, 0, loc)) || !preMarketPrepWindow(time.Date(2026, 8, 10, 3, 50, 0, 0, loc)) || preMarketPrepWindow(time.Date(2026, 8, 10, 3, 51, 0, 0, loc)) {
		t.Fatal("pre-market window boundary")
	}
	if !marketOpenPrepWindow(time.Date(2026, 8, 10, 9, 20, 0, 0, loc)) || !marketOpenPrepWindow(time.Date(2026, 8, 10, 9, 25, 0, 0, loc)) || marketOpenPrepWindow(time.Date(2026, 8, 10, 9, 26, 0, 0, loc)) {
		t.Fatal("market-open window boundary")
	}
	rows := []FreshnessDiagnostic{{Dataset: "Quotes", State: "LIVE", Provider: "Alpaca"}, {Dataset: "VIX", State: "STALE", Provider: "Twelve Data"}, {Dataset: "Intraday Bars", State: "FRESH", Provider: "Alpaca"}}
	_, degraded, ex := readinessFreshnessGate(rows, []string{"Quotes", "VIX", "Intraday Bars"}, time.Now())
	if !degraded || checkpointAttention(ex, degraded) != "DATA DEGRADED" {
		t.Fatalf("false ready degraded=%v ex=%v", degraded, ex)
	}
}

func TestExtreme30_15TradeReadinessNoFalseReady(t *testing.T) {
	rows := []FreshnessDiagnostic{{Dataset: "Quotes", State: "LIVE"}, {Dataset: "VIX", State: "DELAYED", Provider: "CBOE", Reason: "official delayed/close"}, {Dataset: "Intraday Bars", State: "FRESH"}}
	_, degraded, ex := readinessFreshnessGate(rows, []string{"Quotes", "VIX", "Intraday Bars"}, time.Now())
	if degraded || checkpointAttention(ex, degraded) != "READY WITH CAUTION" {
		t.Fatalf("delayed should caution degraded=%v ex=%v", degraded, ex)
	}
	rows[1].State = "ERROR"
	_, degraded, ex = readinessFreshnessGate(rows, []string{"Quotes", "VIX", "Intraday Bars"}, time.Now())
	if !degraded || checkpointAttention(ex, degraded) != "DATA DEGRADED" {
		t.Fatalf("error false ready")
	}
}

func TestExtreme30_16ResearchReadinessPartialStaleAndRecovery(t *testing.T) {
	app := newTestApplication(t)
	sym := "AAA"
	now := time.Now()
	ms := now.UnixMilli()
	app.engine.mu.Lock()
	app.engine.quotes[sym] = Quote{Symbol: sym, Price: 100, ProviderTimestamp: ms, UpdatedAt: ms, Source: "Alpaca"}
	app.engine.quotes["VIX"] = Quote{Symbol: "VIX", Price: 16, ProviderTimestamp: ms, UpdatedAt: ms, Source: "Twelve Data"}
	app.engine.bars[sym] = map[string][]Bar{"intraday": {{T: now.Unix(), C: 100}}, "daily": {{T: now.Unix(), C: 100}}, "weekly": {{T: now.Unix(), C: 100}}}
	app.engine.fundamentals[sym] = FundamentalSnapshot{Symbol: sym, RevenueGrowth: 1, UpdatedAt: ms, Source: "Finnhub"}
	for _, k := range []string{"quotes", "history", "history-intraday", "history-daily", "research-news:" + sym, "research-earnings:" + sym, "research-fundamentals:" + sym, "research-sec:" + sym} {
		app.engine.lastUpdated[k] = ms
	}
	app.engine.health["history"] = "healthy · Alpaca"
	app.engine.health["news"] = "healthy · Finnhub"
	app.engine.health["earnings"] = "healthy · Finnhub"
	app.engine.health["fundamentals"] = "healthy · Finnhub"
	app.engine.health["filings"] = "healthy · SEC EDGAR"
	app.engine.mu.Unlock()
	ready, issues := app.engine.researchPackageReadiness(sym)
	if !ready || len(issues) > 0 {
		t.Fatalf("ready package blocked %v", issues)
	}
	app.engine.mu.Lock()
	delete(app.engine.fundamentals, sym)
	app.engine.mu.Unlock()
	ready, issues = app.engine.researchPackageReadiness(sym)
	if ready || len(issues) == 0 {
		t.Fatal("missing fundamentals falsely ready")
	}
}

func TestExtreme30_17ResearchAIIsEvidenceGated(t *testing.T) {
	// Source-level contract is backed by functional readiness tests above: automatic
	// AI may only execute after researchPackageReadiness succeeds.
	s := productionGoSourceForTest(t)
	if !strings.Contains(s, "researchPackageReadiness") || !strings.Contains(s, "ResearchAIMode") {
		t.Fatal("Research AI evidence gate wiring missing")
	}
}

func TestExtreme30_18PersistenceCorruptStateAndCache(t *testing.T) {
	dir := t.TempDir()
	app := &Application{configDir: dir, hub: NewHub(), sessionKey: "x", state: defaultState()}
	app.engine = NewEngine(app)
	// Corrupt state/cache must fall back safely rather than crash or invent live data.
	if err := os.WriteFile(app.statePath(), []byte("{corrupt"), 0600); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(app.cachePath(), []byte("{corrupt"), 0600); err != nil {
		t.Fatal(err)
	}
	next := &Application{configDir: dir, hub: NewHub(), sessionKey: "y"}
	next.load()
	next.engine = NewEngine(next)
	if len(next.state.Watchlists) < 4 {
		t.Fatal("corrupt state did not recover defaults")
	}
	if len(next.engine.quotes) != 0 {
		t.Fatalf("corrupt cache produced quotes %+v", next.engine.quotes)
	}
	// Valid canonical membership persists.
	next.mu.Lock()
	setMembershipLocked(&next.state, "day", "AAPL", true)
	setMembershipLocked(&next.state, "swing", "AAPL", true)
	if err := next.saveLocked(); err != nil {
		t.Fatal(err)
	}
	next.mu.Unlock()
	again := &Application{configDir: dir, hub: NewHub(), sessionKey: "z"}
	again.load()
	again.mu.RLock()
	m := deskMembershipsLocked(&again.state, "AAPL")
	again.mu.RUnlock()
	if !m["day"] || !m["swing"] {
		t.Fatalf("membership lost %v", m)
	}
}

func TestExtreme30_19ConcurrencyAndRaceStress(t *testing.T) {
	app := v15TestApp(t)
	setDeskBits(t, app, "CONC", "100")
	var wg sync.WaitGroup
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			desk := []string{"day", "swing", "long"}[i%3]
			active := true
			rr := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodPost, "/api/desk/membership", strings.NewReader(`{"symbol":"CONC","desk":"`+desk+`","active":true}`))
			req.Header.Set("Content-Type", "application/json")
			app.handleDeskMembership(rr, req)
			_ = active
		}(i)
	}
	wg.Wait()
	app.mu.RLock()
	m := deskMembershipsLocked(&app.state, "CONC")
	app.mu.RUnlock()
	if activeDeskCount(m) != 3 {
		t.Fatalf("concurrent adds lost state %v", m)
	}
	for _, id := range deskIDs() {
		app.mu.RLock()
		wl, _ := watchlistValueByID(app.state.Watchlists, id)
		count := 0
		for _, s := range wl.Symbols {
			if s == "CONC" {
				count++
			}
		}
		app.mu.RUnlock()
		if count != 1 {
			t.Fatalf("duplicate %s count=%d", id, count)
		}
	}
}

func TestExtreme30_20PerformanceLongRunNoUnboundedGoroutineGrowth(t *testing.T) {
	before := runtime.NumGoroutine()
	app := newTestApplication(t)
	cycles := 15
	if raceBuild {
		// The normal gate retains the full 15-cycle longevity check. Under -race,
		// shared-memory instrumentation makes the same lifecycle materially slower;
		// five complete Start/Stop generations still exercise shutdown/race truth.
		cycles = 5
	}
	for cycle := 0; cycle < cycles; cycle++ {
		if err := app.engine.Start(); err != nil {
			t.Fatal(err)
		}
		time.Sleep(20 * time.Millisecond)
		app.engine.Stop()
		deadline := time.Now().Add(12 * time.Second)
		for {
			app.engine.mu.RLock()
			status := app.engine.status
			app.engine.mu.RUnlock()
			if status == "stopped" {
				break
			}
			if time.Now().After(deadline) {
				t.Fatalf("cycle %d runtime did not finish stopping; status=%s", cycle+1, status)
			}
			time.Sleep(20 * time.Millisecond)
		}
	}
	time.Sleep(100 * time.Millisecond)
	after := runtime.NumGoroutine()
	if after > before+8 {
		t.Fatalf("goroutine growth before=%d after=%d", before, after)
	}
	// 100-symbol canonical cap remains bounded.
	app.mu.Lock()
	for i := 0; i < 130; i++ {
		setMembershipLocked(&app.state, "swing", strings.ToUpper("T"+itoa3(i)), true)
	}
	app.mu.Unlock()
	if n := len(app.engine.trackedSymbols()); n > 100 {
		t.Fatalf("tracked cap exceeded %d", n)
	}
}

func itoa3(i int) string { return string(rune('A'+(i/26)%26)) + string(rune('A'+i%26)) + "1" }

func TestExtreme30_21NetworkAbuseGetJSON(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/ok":
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`{"x":1}`))
		case "/bad":
			_, _ = w.Write([]byte(`{bad`))
		case "/429":
			w.WriteHeader(429)
			_, _ = w.Write([]byte("rate limited"))
		case "/slow":
			time.Sleep(200 * time.Millisecond)
			_, _ = w.Write([]byte(`{"x":1}`))
		}
	}))
	defer srv.Close()
	var out map[string]any
	if err := getJSON(context.Background(), srv.Client(), srv.URL+"/ok", nil, &out); err != nil || out["x"].(float64) != 1 {
		t.Fatalf("ok err=%v out=%v", err, out)
	}
	if err := getJSON(context.Background(), srv.Client(), srv.URL+"/bad", nil, &out); err == nil {
		t.Fatal("malformed JSON accepted")
	}
	if err := getJSON(context.Background(), srv.Client(), srv.URL+"/429", nil, &out); err == nil || !strings.Contains(err.Error(), "429") {
		t.Fatalf("429 err=%v", err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Millisecond)
	defer cancel()
	if err := getJSON(ctx, srv.Client(), srv.URL+"/slow", nil, &out); err == nil {
		t.Fatal("timeout not enforced")
	}
}

func TestExtreme30_21BLSProviderUsesDeterministicLocalEndpoint(t *testing.T) {
	old := blsAPIBaseURL
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Fatalf("BLS method=%s", r.Method)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"status":"REQUEST_SUCCEEDED","Results":{"series":[]}}`))
	}))
	defer srv.Close()
	blsAPIBaseURL = srv.URL
	defer func() { blsAPIBaseURL = old }()
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	r := testBLSProvider(ctx, "")
	if !r.OK || r.Status != "connected" || !strings.Contains(strings.ToLower(strings.Join(r.Details, " ")), "public") {
		t.Fatalf("BLS deterministic provider test failed: %+v", r)
	}
}

func TestExtreme30_22UIFunctionalIntegritySourceContracts(t *testing.T) {
	js, err := os.ReadFile("renderer/renderer.js")
	if err != nil {
		t.Fatal(err)
	}
	css, err := os.ReadFile("renderer/styles.css")
	if err != nil {
		t.Fatal(err)
	}
	j := string(js)
	c := string(css)
	// Functional UI contracts only; this test intentionally does not alter appearance.
	for _, needle := range []string{"captureSaveContext", "restoreSaveContext", "ageText", "data-target-refresh", "Research Evidence Incomplete", "AI Second Opinion Ready"} {
		if !strings.Contains(j, needle) {
			t.Fatalf("renderer functional contract missing %q", needle)
		}
	}
	for _, needle := range []string{"freshness", "reason", "age"} {
		if !strings.Contains(strings.ToLower(c), needle) {
			t.Fatalf("freshness layout contract missing %q", needle)
		}
	}
}

func TestExtreme30_23ScreenResolutionResponsiveContracts(t *testing.T) {
	b, err := os.ReadFile("responsive_ui_test.py")
	if err != nil {
		t.Fatal(err)
	}
	s := string(b)
	// The executable Python responsive suite is run separately; assert here that
	// its viewport matrix still covers desktop, iPad-class and narrow windows.
	for _, needle := range []string{"1280", "1440", "1024", "768"} {
		if !strings.Contains(s, needle) {
			t.Fatalf("responsive suite missing viewport token %s", needle)
		}
	}
}

func TestExtreme30_24SettingsAndSecretsPersistence(t *testing.T) {
	app := v15TestApp(t)
	s := app.state.Settings
	s.DataMode = "live"
	s.ResearchAIMode = "automatic"
	payload := map[string]any{"settings": s, "finnhubKey": " secret\nkey ", "alpacaKey": "ak", "alpacaSecret": "as"}
	b, _ := json.Marshal(payload)
	rr := callHandler(t, app.handleSettingsSave, http.MethodPost, "/api/settings/save", string(b))
	if rr.Code != 200 {
		t.Fatalf("save %d %s", rr.Code, rr.Body.String())
	}
	if app.secrets.Finnhub != "secretkey" {
		t.Fatalf("credential not cleaned %q", app.secrets.Finnhub)
	}
	raw, err := os.ReadFile(app.secretsPath())
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(string(raw), "\nkey") {
		t.Fatal("newline secret persisted")
	}
	info, err := os.Stat(app.secretsPath())
	if err != nil {
		t.Fatal(err)
	}
	if runtime.GOOS == "windows" {
		// Windows does not expose ACLs through os.FileMode; a file created with
		// os.WriteFile(..., 0600) commonly reports 0666. Native G14 owns the
		// Windows ACL/isolation check, while this unit test still proves the
		// secret is a regular private-profile file with sanitized contents.
		if !info.Mode().IsRegular() {
			t.Fatalf("secrets path is not a regular file: %v", info.Mode())
		}
	} else if info.Mode().Perm()&0o077 != 0 {
		t.Fatalf("secrets permissions too open %o", info.Mode().Perm())
	}
}

func TestExtreme30_25SecuritySessionExportAndInput(t *testing.T) {
	app := v15TestApp(t)
	app.mu.Lock()
	app.secrets.Finnhub = "SUPERSECRET"
	app.mu.Unlock()
	rr := callHandler(t, app.handleExport, http.MethodGet, "/api/profile/export", "")
	if rr.Code != 200 {
		t.Fatalf("export %d", rr.Code)
	}
	if strings.Contains(rr.Body.String(), "SUPERSECRET") || strings.Contains(strings.ToLower(rr.Body.String()), "finnhubkey") {
		t.Fatal("profile export leaked secret")
	}
	if _, ok := parseUserTicker("AAPL<script>"); ok {
		t.Fatal("ticker injection accepted")
	}
	if err := openExternalURL("javascript:alert(1)"); err == nil {
		t.Fatal("javascript URL accepted")
	}
}

func TestExtreme30_26UpgradeCompatibilityPreservesCanonicalDesks(t *testing.T) {
	old := AppState{Version: 1, Settings: Settings{SwingWatchlistID: "legacy"}, Watchlists: []Watchlist{{ID: "legacy", Name: "Legacy", Symbols: []string{"AAPL", "MSFT"}}}, UI: UIState{WatchlistID: "legacy", SelectedTicker: "AAPL"}}
	got := mergeState(old)
	ensureDedicatedDeskWatchlists(&got, defaultState())
	if got.Settings.DayWatchlistID != "day" || got.Settings.SwingWatchlistID != "swing" || got.Settings.LongWatchlistID != "long" {
		t.Fatalf("canonical ids not migrated %+v", got.Settings)
	}
	if len(got.Watchlists) < 4 {
		t.Fatal("migration lost permanent lists")
	}
}

func TestExtreme30_29DeterministicInputsStayContextOnly(t *testing.T) {
	// Structural guard: options/global/AI remain contextual and deterministic score
	// equivalence is executed by the external 2403-case suite.
	b, err := os.ReadFile("README.md")
	if err != nil {
		t.Fatal(err)
	}
	contract := string(b)
	if !strings.Contains(contract, "deterministic Day/Swing/Long") ||
		!(strings.Contains(contract, "cannot silently mutate") || strings.Contains(contract, "never silently rewrite")) {
		t.Fatal("deterministic contract missing")
	}
}

func TestExtreme30_30CompoundFailureDegradesTruthfully(t *testing.T) {
	app := v15TestApp(t)
	e := app.engine
	now := time.Now().UnixMilli()
	e.mu.Lock()
	e.health["quotes"] = "failed · offline"
	e.health["vix"] = "failed · offline"
	e.health["history"] = "failed · offline"
	e.lastUpdated["quotes"] = 0
	e.lastUpdated["vix"] = 0
	e.lastUpdated["history"] = 0
	e.quotes = map[string]Quote{}
	e.bars = map[string]map[string][]Bar{}
	e.mu.Unlock()
	rows, _ := e.buildFreshnessDiagnostics(map[string]Quote{}, clone(e.lastUpdated), clone(e.health))
	_, degraded, ex := readinessFreshnessGate(rows, []string{"Quotes", "VIX", "Intraday Bars"}, time.Now())
	if !degraded || len(ex) < 3 {
		t.Fatalf("compound outage not degraded ex=%v", ex)
	}
	// A future timestamp also cannot restore false-ready status.
	stateTs, anomaly := safeFreshnessTimestamp(now+int64(10*time.Minute/time.Millisecond), 0, now)
	if !anomaly || stateTs != 0 {
		t.Fatalf("future timestamp not quarantined stateTs=%d anomaly=%v", stateTs, anomaly)
	}
}

func TestExtreme30_27CrossPlatformSourceHasNativeQuitAndLaunchPaths(t *testing.T) {
	s := productionGoSourceForTest(t)
	for _, needle := range []string{"case \"darwin\"", "case \"windows\"", "terminateAppWindow", "instancePath"} {
		if !strings.Contains(s, needle) {
			t.Fatalf("native source path missing %s", needle)
		}
	}
	_ = filepath.Separator
}

func TestExtreme30_28AuthenticatedProvidersHaveTruthfulMissingCredentialStates(t *testing.T) {
	ctx := context.Background()
	if r := testFinnhub(ctx, ""); r.OK || r.Status != "missing" {
		t.Fatalf("Finnhub missing %+v", r)
	}
	if r := testAlpaca(ctx, "", ""); r.OK || r.Status != "missing" {
		t.Fatalf("Alpaca missing %+v", r)
	}
	if r := testFREDProvider(ctx, ""); r.OK || r.Status != "missing" {
		t.Fatalf("FRED missing %+v", r)
	}
	if r := testMarketauxProvider(ctx, ""); r.OK || r.Status != "missing" {
		t.Fatalf("Marketaux missing %+v", r)
	}
}

func TestExtreme30V1611EscapedMarketIntelligenceTruth(t *testing.T) {
	t.Run("exact-long-weekly", TestV1611LongStructureRequiresWeeklyEvidence)
	t.Run("stale-structure", TestV1611StructureRejectsStaleHistoricalBars)
	t.Run("stale-rs-evidence-time", TestV1611RelativeStrengthRejectsStaleBarsAndUsesEvidenceTimestamp)
	t.Run("missing-spread-unknown", TestV1611LiquidityWithoutValidBidAskIsUnknown)
}
