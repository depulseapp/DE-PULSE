package main

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
	"time"
)

func TestCleanCredential(t *testing.T) {
	got := cleanCredential("  \"abc\n123\"  ")
	if got != "abc123" {
		t.Fatalf("got %q", got)
	}
}
func TestMergeStateMigration(t *testing.T) {
	st := mergeState(AppState{})
	if len(st.Watchlists) < 3 {
		t.Fatalf("expected default lists")
	}
	if st.Settings.SwingWatchlistID == "" || st.Settings.DayWatchlistID == "" || st.Settings.LongWatchlistID == "" {
		t.Fatal("desk assignments missing")
	}
	if st.Settings.MarketContext <= 0 {
		t.Fatal("market context default missing")
	}
}
func TestUniqueSymbols(t *testing.T) {
	out := uniqueSymbols([]string{" aapl ", "AAPL", "msft", ""})
	if len(out) != 2 || out[0] != "AAPL" || out[1] != "MSFT" {
		t.Fatalf("%v", out)
	}
}
func TestDemoRuntimeAndCache(t *testing.T) {
	dir := t.TempDir()
	app := &Application{configDir: dir, hub: NewHub(), sessionKey: "x", state: defaultState()}
	app.engine = NewEngine(app)
	if err := app.engine.Start(); err != nil {
		t.Fatal(err)
	}
	time.Sleep(50 * time.Millisecond)
	snap := app.engine.Snapshot()
	if snap.Status != "running" || len(snap.Quotes) == 0 || len(snap.Bars) == 0 {
		t.Fatalf("bad snapshot: %s q=%d b=%d", snap.Status, len(snap.Quotes), len(snap.Bars))
	}
	app.engine.Stop()
	if _, err := os.Stat(filepath.Join(dir, "market-cache.json")); err != nil {
		t.Fatal("cache not saved")
	}
	e2 := NewEngine(app)
	if len(e2.quotes) == 0 || len(e2.bars) == 0 {
		t.Fatal("cache not restored")
	}
}
func TestProviderMissingCredentials(t *testing.T) {
	r := testFinnhub(context.Background(), "")
	if r.OK || r.Status != "missing" {
		t.Fatalf("%+v", r)
	}
	a := testAlpaca(context.Background(), "", "")
	if a.OK || a.Status != "missing" {
		t.Fatalf("%+v", a)
	}
}
func TestAtomicWrite(t *testing.T) {
	p := filepath.Join(t.TempDir(), "x.json")
	if err := atomicWrite(p, []byte("ok"), 0600); err != nil {
		t.Fatal(err)
	}
	b, _ := os.ReadFile(p)
	if string(b) != "ok" {
		t.Fatal(string(b))
	}
}

func newTestApplication(t *testing.T) *Application {
	t.Helper()
	app := &Application{configDir: t.TempDir(), hub: NewHub(), sessionKey: "test-session", state: defaultState()}
	app.engine = NewEngine(app)
	return app
}

func TestPriceCacheRoundTrip(t *testing.T) {
	app := newTestApplication(t)
	app.engine.mu.Lock()
	app.engine.quotes["SPY"] = Quote{Symbol: "SPY", Price: 600, PreviousClose: 598, UpdatedAt: time.Now().UnixMilli()}
	app.engine.bars["SPY"] = map[string][]Bar{"daily": {{T: time.Now().Unix(), O: 590, H: 605, L: 588, C: 600, V: 1000}}}
	app.engine.mu.Unlock()
	if err := app.engine.saveCache(); err != nil {
		t.Fatal(err)
	}
	other := NewEngine(app)
	if other.quotes["SPY"].Price != 600 || len(other.bars["SPY"]["daily"]) != 1 {
		t.Fatalf("cache mismatch: %+v %+v", other.quotes["SPY"], other.bars["SPY"])
	}
	if other.quotes["SPY"].DataState != "cache" || other.quotes["SPY"].FeedType != "cache" {
		t.Fatalf("restored quote must be explicitly cached: %+v", other.quotes["SPY"])
	}
}

func TestTrackedSymbolsCapAndDeduplication(t *testing.T) {
	app := newTestApplication(t)
	var syms []string
	for i := 0; i < 140; i++ {
		syms = append(syms, fmt.Sprintf("T%03d", i))
	}
	for i := range app.state.Watchlists {
		if app.state.Watchlists[i].ID == "swing" {
			app.state.Watchlists[i].Symbols = syms
		}
	}
	got := app.engine.trackedSymbols()
	if len(got) != 100 {
		t.Fatalf("expected 100 symbols, got %d", len(got))
	}
	seen := map[string]bool{}
	for _, s := range got {
		if seen[s] {
			t.Fatalf("duplicate %s", s)
		}
		seen[s] = true
	}
}

func TestHubDropsStaleClientsWithoutBlocking(t *testing.T) {
	h := NewHub()
	ch := h.Subscribe()
	for i := 0; i < 1000; i++ {
		h.Broadcast(map[string]int{"n": i})
	}
	select {
	case <-ch:
	case <-time.After(time.Second):
		t.Fatal("expected at least one event")
	}
	h.Unsubscribe(ch)
}

func TestExtractAITextVariants(t *testing.T) {
	if got := extractAIText(map[string]any{"output_text": "  hello  "}); got != "hello" {
		t.Fatalf("got %q", got)
	}
	data := map[string]any{"output": []any{map[string]any{"content": []any{map[string]any{"text": "one"}, map[string]any{"output_text": "two"}}}}}
	if got := extractAIText(data); got != "one\ntwo" {
		t.Fatalf("got %q", got)
	}
}

func TestNewsMatches(t *testing.T) {
	if !newsMatchesResearchSymbol(NewsItem{Symbols: []string{"NVDA"}}, "NVDA") {
		t.Fatal("direct symbol should match")
	}
	if !newsMatchesResearchSymbol(NewsItem{Related: "AAPL,NVDA"}, "NVDA") {
		t.Fatal("related symbol should match")
	}
	if newsMatchesResearchSymbol(NewsItem{Symbols: []string{"TSLA"}}, "NVDA") {
		t.Fatal("unexpected match")
	}
}

func TestVIXIncludedInCommonMarketAndDemo(t *testing.T) {
	found := false
	for _, symbol := range generalSymbols {
		if symbol == "VIX" {
			found = true
			break
		}
	}
	if !found {
		t.Fatal("VIX missing from common market symbols")
	}
	app := newTestApplication(t)
	app.engine.seedDemo()
	snap := app.engine.Snapshot()
	if snap.Quotes["VIX"].Price <= 0 {
		t.Fatal("demo VIX quote missing")
	}
	if len(snap.Bars["VIX"]["daily"]) < 20 {
		t.Fatal("demo VIX history missing")
	}
}

func TestDedicatedDeskWatchlistMigrationFromSharedList(t *testing.T) {
	shared := Watchlist{ID: "shared", Name: "My Watchlist", Symbols: []string{"AAPL", "MSFT"}}
	st := AppState{
		Watchlists: []Watchlist{shared},
		Settings: Settings{
			SwingWatchlistID: "shared",
			DayWatchlistID:   "shared",
			LongWatchlistID:  "shared",
		},
	}
	merged := mergeState(st)
	for _, id := range []string{"swing", "day", "long"} {
		wl, ok := watchlistValueByID(merged.Watchlists, id)
		if !ok {
			t.Fatalf("missing dedicated %s list", id)
		}
		if len(wl.Symbols) != 2 || wl.Symbols[0] != "AAPL" || wl.Symbols[1] != "MSFT" {
			t.Fatalf("unexpected %s symbols: %v", id, wl.Symbols)
		}
	}
	if merged.Settings.SwingWatchlistID != "swing" || merged.Settings.DayWatchlistID != "day" || merged.Settings.LongWatchlistID != "long" {
		t.Fatalf("desk assignments not canonical: %+v", merged.Settings)
	}
	merged.Watchlists[0].Symbols[0] = "NVDA"
	day, _ := watchlistValueByID(merged.Watchlists, "day")
	if day.Symbols[0] != "AAPL" {
		t.Fatal("dedicated desk lists share mutable symbol storage")
	}
}

func TestPublicStateRepairsMissingDiscoveryWatchlist(t *testing.T) {
	app := newTestApplication(t)
	app.mu.Lock()
	var kept []Watchlist
	for _, wl := range app.state.Watchlists {
		if wl.ID != "discovery" {
			kept = append(kept, wl)
		}
	}
	app.state.Watchlists = kept
	app.state.Settings.DiscoveryWatchlistID = ""
	pub := app.publicStateLockedForUser(bootstrapOwnerID)
	app.mu.Unlock()
	if pub.Settings.DiscoveryWatchlistID != "discovery" {
		t.Fatalf("discovery watchlist id was not repaired: %q", pub.Settings.DiscoveryWatchlistID)
	}
	wl, ok := watchlistValueByID(pub.Watchlists, "discovery")
	if !ok {
		t.Fatal("public state did not synthesize Discovery watchlist")
	}
	if wl.Symbols == nil {
		t.Fatal("Discovery symbols must be a safe empty array, not nil")
	}
}

func TestTrackedSymbolsUseDedicatedDesksOnly(t *testing.T) {
	app := newTestApplication(t)
	app.state.Watchlists = append(app.state.Watchlists, Watchlist{ID: "archived", Name: "Archived", Symbols: []string{"ZZZZ"}})
	got := app.engine.trackedSymbols()
	for _, symbol := range got {
		if symbol == "ZZZZ" {
			t.Fatal("archived generic watchlist should not be subscribed")
		}
	}
}

func TestMacNativeShellUsesRegisteredSubclassByName(t *testing.T) {
	s := productionGoSourceForTest(t)
	if strings.Contains(s, "const Delegate=ObjC.registerSubclass") || strings.Contains(s, "Delegate.alloc") {
		t.Fatal("macOS shell must not use ObjC.registerSubclass return value; JXA returns undefined")
	}
	if !strings.Contains(s, "$.DePulseDelegateV121.alloc.init") {
		t.Fatal("macOS shell must instantiate the registered subclass through the Objective-C bridge")
	}
}

func TestMacNativeShellProvidesStandardEditShortcuts(t *testing.T) {
	s := productionGoSourceForTest(t)
	for _, marker := range []string{
		"initWithTitle('Edit')",
		"initWithTitleActionKeyEquivalent('Cut','cut:','x')",
		"initWithTitleActionKeyEquivalent('Copy','copy:','c')",
		"initWithTitleActionKeyEquivalent('Paste','paste:','v')",
		"initWithTitleActionKeyEquivalent('Select All','selectAll:','a')",
	} {
		if !strings.Contains(s, marker) {
			t.Fatalf("macOS WKWebView shell is missing standard edit shortcut marker %q", marker)
		}
	}
}

func TestV15FinnhubFailoverRespectsFreeLimit(t *testing.T) {
	app := newTestApplication(t)
	var syms []string
	for i := 0; i < 80; i++ {
		syms = append(syms, fmt.Sprintf("L%03d", i))
	}
	for i := range app.state.Watchlists {
		if app.state.Watchlists[i].ID == "swing" {
			app.state.Watchlists[i].Symbols = syms
		}
	}
	got := app.engine.liveSymbols()
	if len(got) == 0 || len(got) > finnhubPlanMaxSymbols {
		t.Fatalf("Finnhub failover must remain within %d-symbol ceiling, got %d", finnhubPlanMaxSymbols, len(got))
	}
	for _, symbol := range got {
		if symbol == "VIX" {
			t.Fatal("VIX should be resolved through the dedicated snapshot path")
		}
	}
}

func TestMarketSessionET(t *testing.T) {
	loc, err := time.LoadLocation("America/New_York")
	if err != nil {
		t.Fatal(err)
	}
	cases := []struct {
		name string
		at   time.Time
		want string
	}{
		{"overnight-thursday", time.Date(2026, 8, 6, 21, 0, 0, 0, loc), "overnight"},
		{"overnight-early", time.Date(2026, 8, 7, 2, 0, 0, 0, loc), "overnight"},
		{"pre", time.Date(2026, 8, 6, 8, 0, 0, 0, loc), "pre-market"},
		{"regular", time.Date(2026, 8, 6, 10, 0, 0, 0, loc), "regular"},
		{"after", time.Date(2026, 8, 6, 17, 0, 0, 0, loc), "after-hours"},
		{"friday-night", time.Date(2026, 8, 7, 21, 0, 0, 0, loc), "weekend"},
		{"saturday", time.Date(2026, 8, 8, 12, 0, 0, 0, loc), "weekend"},
		{"sunday-day", time.Date(2026, 8, 9, 12, 0, 0, 0, loc), "weekend"},
		{"sunday-overnight-open", time.Date(2026, 8, 9, 20, 30, 0, 0, loc), "overnight"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := marketSessionET(tc.at); got != tc.want {
				t.Fatalf("got %q, want %q", got, tc.want)
			}
		})
	}
}

func TestRESTSnapshotDoesNotOverwriteRecentWebSocketPrice(t *testing.T) {
	app := newTestApplication(t)
	now := time.Now().UnixMilli()
	app.engine.mu.Lock()
	app.engine.quotes["AAPL"] = Quote{Symbol: "AAPL", Price: 211.25, Source: "finnhub-websocket", ProviderTimestamp: now, UpdatedAt: now, PreviousClose: 205}
	app.engine.mu.Unlock()
	app.engine.mergeFinnhubSnapshot("AAPL", finnhubQuoteResponse{Current: 209, Previous: 206, Open: 207, High: 212, Low: 204, Timestamp: time.Now().Unix()})
	got := app.engine.Snapshot().Quotes["AAPL"]
	if got.Price != 211.25 || got.Source != "finnhub-websocket" {
		t.Fatalf("REST overwrote recent stream price: %+v", got)
	}
	if got.PreviousClose != 206 {
		t.Fatalf("REST close context was not merged: %+v", got)
	}
}

func TestWSClientReadsFinnhubTradeFrame(t *testing.T) {
	clientConn, serverConn := net.Pipe()
	defer clientConn.Close()
	defer serverConn.Close()
	ws := &WSClient{conn: clientConn, r: bufio.NewReader(clientConn), w: bufio.NewWriter(clientConn)}
	payload := []byte(`{"type":"trade","data":[{"s":"AAPL","p":211.25,"t":1770000000000,"v":10}]}`)
	go func() {
		frame := append([]byte{0x81, byte(len(payload))}, payload...)
		_, _ = serverConn.Write(frame)
	}()
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	got, err := ws.ReadText(ctx)
	if err != nil {
		t.Fatal(err)
	}
	if string(got) != string(payload) {
		t.Fatalf("got %q, want %q", got, payload)
	}
}

func TestAlpacaOvernightPrefersRealTimeIndicativeQuote(t *testing.T) {
	var snap alpacaLiveSnapshot
	snap.LatestTrade.Price = 100
	snap.LatestTrade.Time = "2026-08-07T01:00:00Z"
	snap.LatestQuote.Bid = 101
	snap.LatestQuote.Ask = 103
	snap.LatestQuote.Time = "2026-08-07T01:14:59Z"
	snap.MinuteBar.Close = 102.25
	snap.MinuteBar.Time = "2026-08-07T01:14:00Z"
	price, stamp, kind := alpacaSnapshotPrice(snap, "overnight", "overnight")
	if price != 102 || kind != "indicative-quote" {
		t.Fatalf("overnight should prefer real-time indicative quote midpoint, got price=%v kind=%q", price, kind)
	}
	if stamp != providerTimeMillis(snap.LatestQuote.Time) {
		t.Fatalf("overnight timestamp should come from latest quote")
	}
}

func TestAlpacaExtendedHoursPrefersFreshQuote(t *testing.T) {
	var snap alpacaLiveSnapshot
	snap.LatestTrade.Price = 100
	snap.LatestTrade.Time = "2026-08-07T20:00:00Z"
	snap.LatestQuote.Bid = 101
	snap.LatestQuote.Ask = 103
	snap.LatestQuote.Time = "2026-08-07T21:14:59Z"
	price, stamp, kind := alpacaSnapshotPrice(snap, "iex", "after-hours")
	if price != 102 || kind != "quote" {
		t.Fatalf("after-hours should prefer the fresh quote midpoint, got price=%v kind=%q", price, kind)
	}
	if stamp != providerTimeMillis(snap.LatestQuote.Time) {
		t.Fatal("after-hours timestamp should come from the latest quote")
	}
}

func TestV15AlpacaRegularSnapshotRestoresPreferredSource(t *testing.T) {
	app := newTestApplication(t)
	now := time.Now().UnixMilli()
	app.engine.mu.Lock()
	app.engine.quotes["AAPL"] = Quote{Symbol: "AAPL", Price: 211.25, Source: "finnhub-websocket", ProviderTimestamp: now, UpdatedAt: now, PreviousClose: 205}
	app.engine.mu.Unlock()
	var snap alpacaLiveSnapshot
	snap.LatestTrade.Price = 209
	snap.LatestTrade.Time = time.Now().UTC().Format(time.RFC3339Nano)
	applied := app.engine.mergeAlpacaLiveSnapshot("AAPL", 209, now, snap, "iex", "trade")
	got := app.engine.Snapshot().Quotes["AAPL"]
	if !applied || got.Price != 209 || got.Source != "alpaca-iex-snapshot" {
		t.Fatalf("Alpaca snapshot must restore preferred v15 source: applied=%v quote=%+v", applied, got)
	}
}

func TestAlpacaOvernightCanReplaceOlderSessionPrice(t *testing.T) {
	app := newTestApplication(t)
	old := time.Now().Add(-10 * time.Minute).UnixMilli()
	app.engine.mu.Lock()
	app.engine.quotes["AAPL"] = Quote{Symbol: "AAPL", Price: 210, Source: "finnhub-websocket", ProviderTimestamp: old, UpdatedAt: old, PreviousClose: 205}
	app.engine.mu.Unlock()
	var snap alpacaLiveSnapshot
	snap.LatestQuote.Bid = 211
	snap.LatestQuote.Ask = 213
	snap.LatestQuote.Time = time.Now().UTC().Format(time.RFC3339Nano)
	price, stamp, kind := alpacaSnapshotPrice(snap, "overnight", "overnight")
	applied := app.engine.mergeAlpacaLiveSnapshot("AAPL", price, stamp, snap, "overnight", kind)
	got := app.engine.Snapshot().Quotes["AAPL"]
	if !applied || got.Price != 212 || got.Source != "alpaca-overnight-indicative-snapshot" {
		t.Fatalf("overnight indicative quote should become current price: applied=%v quote=%+v", applied, got)
	}
}

func TestV12MigrationAddsDiscoveryAndEngineDefaults(t *testing.T) {
	st := AppState{Version: 10, Settings: Settings{DataMode: "demo", SwingWatchlistID: "swing", DayWatchlistID: "day", LongWatchlistID: "long"}, Watchlists: []Watchlist{{ID: "swing", Name: "Swing", Symbols: []string{"AAA"}}, {ID: "day", Name: "Day", Symbols: []string{"BBB"}}, {ID: "long", Name: "Long", Symbols: []string{"CCC"}}}}
	got := mergeState(st)
	if got.Version != 17 {
		t.Fatalf("version=%d", got.Version)
	}
	if !got.Settings.DayEnabled || !got.Settings.SwingEnabled || !got.Settings.LongEnabled {
		t.Fatalf("engines not defaulted on: %+v", got.Settings)
	}
	if got.Settings.OvernightDataMode != "auto" {
		t.Fatalf("overnight mode=%q", got.Settings.OvernightDataMode)
	}
	if got.Settings.DiscoveryWatchlistID != "discovery" {
		t.Fatalf("discovery id=%q", got.Settings.DiscoveryWatchlistID)
	}
	if wl, ok := watchlistValueByID(got.Watchlists, "discovery"); !ok || wl.Name != "Discovery Watchlist" {
		t.Fatalf("discovery watchlist missing: %+v", got.Watchlists)
	}
}

func TestEngineTogglesReduceRequiredSymbolsButKeepSharedSymbols(t *testing.T) {
	app := newTestApplication(t)
	for i := range app.state.Watchlists {
		switch app.state.Watchlists[i].ID {
		case "day":
			app.state.Watchlists[i].Symbols = []string{"DAYONLY", "NVDA"}
		case "swing":
			app.state.Watchlists[i].Symbols = []string{"SWINGONLY", "NVDA"}
		case "long":
			app.state.Watchlists[i].Symbols = []string{"LONGONLY"}
		case "discovery":
			app.state.Watchlists[i].Symbols = []string{"DISC"}
		}
	}
	app.state.Settings.DayEnabled = false
	app.state.Settings.SwingEnabled = true
	app.state.Settings.LongEnabled = false
	got := app.engine.trackedSymbols()
	for _, want := range []string{"NVDA", "SWINGONLY", "DISC", "SPY"} {
		if !contains(got, want) {
			t.Fatalf("missing %s in %v", want, got)
		}
	}
	for _, unwanted := range []string{"DAYONLY", "LONGONLY"} {
		if contains(got, unwanted) {
			t.Fatalf("disabled-engine symbol %s still tracked: %v", unwanted, got)
		}
	}
}

func TestEngineToggleHandlerPersistsState(t *testing.T) {
	app := newTestApplication(t)
	body := strings.NewReader(`{"engine":"day","enabled":false}`)
	req := httptest.NewRequest("POST", "/api/engine/toggle", body)
	rr := httptest.NewRecorder()
	app.handleEngineToggle(rr, req)
	if rr.Code != 200 {
		t.Fatalf("code=%d body=%s", rr.Code, rr.Body.String())
	}
	if app.state.Settings.DayEnabled {
		t.Fatal("day engine should be off")
	}
	var pub PublicState
	if err := json.Unmarshal(rr.Body.Bytes(), &pub); err != nil {
		t.Fatal(err)
	}
	if pub.Settings.DayEnabled {
		t.Fatal("response still reports day enabled")
	}
}

func TestClearMarketCacheClearsMemoryAndDiskButPreservesProfile(t *testing.T) {
	app := newTestApplication(t)
	app.state.Settings.DayEnabled = false
	app.secrets.Finnhub = "secret-key"
	app.engine.mu.Lock()
	app.engine.quotes["AAPL"] = Quote{Symbol: "AAPL", Price: 200}
	app.engine.bars["AAPL"] = map[string][]Bar{"daily": {{T: 1, C: 200}}}
	app.engine.fundamentals["AAPL"] = FundamentalSnapshot{Symbol: "AAPL", PERatio: 20}
	app.engine.mu.Unlock()
	if err := app.engine.saveCache(); err != nil {
		t.Fatal(err)
	}
	req := httptest.NewRequest("POST", "/api/cache/clear", strings.NewReader(`{}`))
	rr := httptest.NewRecorder()
	app.handleCacheClear(rr, req)
	if rr.Code != 200 {
		t.Fatalf("code=%d body=%s", rr.Code, rr.Body.String())
	}
	snap := app.engine.Snapshot()
	if len(snap.Quotes) != 0 || len(snap.Bars) != 0 || len(snap.Fundamentals) != 0 {
		t.Fatalf("cache memory not cleared: q=%d b=%d f=%d", len(snap.Quotes), len(snap.Bars), len(snap.Fundamentals))
	}
	if _, err := os.Stat(app.cachePath()); !os.IsNotExist(err) {
		t.Fatalf("cache file still exists: %v", err)
	}
	if app.secrets.Finnhub != "secret-key" || app.state.Settings.DayEnabled {
		t.Fatal("profile or secrets were changed by cache clear")
	}
	if app.state.LastCacheCleared <= 0 {
		t.Fatal("cache-clear timestamp was not persisted")
	}
	app.mu.RLock()
	pub := app.publicStateLockedForUser(bootstrapOwnerID)
	app.mu.RUnlock()
	if pub.LastCacheCleared != app.state.LastCacheCleared {
		t.Fatalf("public cache-clear timestamp mismatch: public=%d state=%d", pub.LastCacheCleared, app.state.LastCacheCleared)
	}
	data, err := os.ReadFile(app.statePath())
	if err != nil {
		t.Fatal(err)
	}
	var persisted AppState
	if err := json.Unmarshal(data, &persisted); err != nil {
		t.Fatal(err)
	}
	if persisted.LastCacheCleared != app.state.LastCacheCleared {
		t.Fatalf("cache-clear timestamp not stored: file=%d state=%d", persisted.LastCacheCleared, app.state.LastCacheCleared)
	}
}

func TestDemoDiscoveryScannerReturnsExternalCandidates(t *testing.T) {
	app := newTestApplication(t)
	app.state.Settings.DataMode = "demo"
	req := httptest.NewRequest("POST", "/api/discovery/scan", strings.NewReader(`{"mode":"day"}`))
	rr := httptest.NewRecorder()
	app.handleDiscoveryScan(rr, req)
	if rr.Code != 200 {
		t.Fatalf("code=%d body=%s", rr.Code, rr.Body.String())
	}
	var scan ScannerState
	if err := json.Unmarshal(rr.Body.Bytes(), &scan); err != nil {
		t.Fatal(err)
	}
	if scan.Status != "complete" || len(scan.Results) < 5 || scan.Scanned < len(scan.Results) {
		t.Fatalf("bad scanner state: %+v", scan)
	}
	tracked := []string{}
	for _, id := range []string{app.state.Settings.SwingWatchlistID, app.state.Settings.DayWatchlistID, app.state.Settings.LongWatchlistID} {
		if wl, ok := watchlistValueByID(app.state.Watchlists, id); ok {
			tracked = append(tracked, wl.Symbols...)
		}
	}
	tracked = uniqueSymbols(tracked)
	foundExternal := false
	for _, x := range scan.Results {
		if !contains(tracked, x.Symbol) {
			foundExternal = true
			break
		}
	}
	if !foundExternal {
		t.Fatal("scanner did not discover any symbol outside desk watchlists")
	}
}

func TestAlpacaBoatsPrefersLiveTrade(t *testing.T) {
	var snap alpacaLiveSnapshot
	snap.LatestTrade.Price = 105
	snap.LatestTrade.Time = "2026-08-07T02:00:01Z"
	snap.LatestQuote.Bid = 104
	snap.LatestQuote.Ask = 106
	snap.LatestQuote.Time = "2026-08-07T02:00:02Z"
	price, _, kind := alpacaSnapshotPrice(snap, "boats", "overnight")
	if price != 105 || kind != "trade" {
		t.Fatalf("boats should prefer true latest trade, price=%v kind=%s", price, kind)
	}
}

func TestMarketSessionHolidayAndEarlyClose(t *testing.T) {
	loc, err := time.LoadLocation("America/New_York")
	if err != nil {
		t.Fatal(err)
	}
	if got := marketSessionET(time.Date(2026, 7, 3, 10, 0, 0, 0, loc)); got != "closed" { // July 4 2026 is Saturday; observed Friday July 3.
		t.Fatalf("Independence Day observed session=%q", got)
	}
	if got := marketSessionET(time.Date(2026, 11, 27, 12, 30, 0, 0, loc)); got != "regular" {
		t.Fatalf("Black Friday before early close=%q", got)
	}
	if got := marketSessionET(time.Date(2026, 11, 27, 13, 30, 0, 0, loc)); got != "after-hours" {
		t.Fatalf("Black Friday after early close=%q", got)
	}
}

func TestDemoEngineScopingReducesUnneededBars(t *testing.T) {
	app := newTestApplication(t)
	app.mu.Lock()
	app.state.Settings.DayEnabled = false
	app.state.Settings.SwingEnabled = true
	app.state.Settings.LongEnabled = false
	app.state.Watchlists = []Watchlist{
		{ID: "swing", Name: "Swing Watchlist", Symbols: []string{"SWGX"}},
		{ID: "day", Name: "Day Trade Watchlist", Symbols: []string{"DAYX"}},
		{ID: "long", Name: "Long-Term Watchlist", Symbols: []string{"LNGX"}},
		{ID: "discovery", Name: "Discovery Watchlist", Symbols: []string{}},
	}
	app.mu.Unlock()
	app.engine.seedDemo()
	snap := app.engine.Snapshot()
	if _, ok := snap.Quotes["DAYX"]; ok {
		t.Fatal("disabled Day-only symbol should not be tracked")
	}
	if len(snap.Bars["SWGX"]["daily"]) < 200 {
		t.Fatal("enabled Swing symbol should retain daily history")
	}
	if len(snap.Bars["SWGX"]["intraday"]) != 0 || len(snap.Bars["SWGX"]["weekly"]) != 0 {
		t.Fatal("Swing-only symbol should not allocate Day or Long-Term bars")
	}
}

func TestAddSymbolPersistsAndHydratesRunningDemo(t *testing.T) {
	app := newTestApplication(t)
	app.state.Settings.DataMode = "demo"
	if err := app.engine.Start(); err != nil {
		t.Fatal(err)
	}
	defer app.engine.Stop()
	req := httptest.NewRequest("POST", "/api/watchlists/add-symbol", strings.NewReader(`{"watchlistId":"swing","symbol":"SHOP"}`))
	rr := httptest.NewRecorder()
	app.handleAddSymbol(rr, req)
	if rr.Code != 200 {
		t.Fatalf("code=%d body=%s", rr.Code, rr.Body.String())
	}
	app.mu.RLock()
	wl := findWatchlistInState(&app.state, "swing")
	persisted := wl != nil && contains(wl.Symbols, "SHOP")
	app.mu.RUnlock()
	if !persisted {
		t.Fatal("SHOP was not persisted to Swing watchlist")
	}
	snap := app.engine.Snapshot()
	if snap.Quotes["SHOP"].Price <= 0 {
		t.Fatalf("new Demo symbol was not hydrated immediately: %+v", snap.Quotes["SHOP"])
	}
}

func TestV121PublicStateAndHealthExposeBuildIdentity(t *testing.T) {
	app := newTestApplication(t)
	app.mu.RLock()
	pub := app.publicStateLockedForUser(bootstrapOwnerID)
	app.mu.RUnlock()
	if pub.Version != appVersion || pub.BuildID != buildID {
		t.Fatalf("unexpected build identity: version=%q build=%q", pub.Version, pub.BuildID)
	}
	req := httptest.NewRequest("GET", "/api/health", nil)
	rr := httptest.NewRecorder()
	app.routes().ServeHTTP(rr, req)
	if rr.Code != 200 {
		t.Fatalf("health code=%d body=%s", rr.Code, rr.Body.String())
	}
	if !strings.Contains(rr.Body.String(), `"version":"`+appVersion+`"`) || !strings.Contains(rr.Body.String(), buildID) {
		t.Fatalf("health did not expose version/build: %s", rr.Body.String())
	}
}

func TestV121TerminateWindowRemovesInstanceFileWithoutWindowPID(t *testing.T) {
	app := newTestApplication(t)
	data, err := json.Marshal(instanceInfo{URL: "http://127.0.0.1:9999", BackendPID: os.Getpid(), WindowPID: 0})
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(instancePath(app.configDir), data, 0600); err != nil {
		t.Fatal(err)
	}
	app.terminateAppWindow()
	if _, err := os.Stat(instancePath(app.configDir)); !os.IsNotExist(err) {
		t.Fatalf("instance file should be removed, stat err=%v", err)
	}
}

func TestV121ScannerScoreUsesGapSpreadAndRelativeVolume(t *testing.T) {
	var snap alpacaLiveSnapshot
	snap.LatestTrade.Price = 110
	snap.LatestQuote.Bid = 109.9
	snap.LatestQuote.Ask = 110.1
	snap.DailyBar.Open = 108
	snap.DailyBar.Close = 110
	snap.DailyBar.Volume = 2_000_000
	snap.PrevDailyBar.Close = 100
	snap.PrevDailyBar.Volume = 1_000_000
	got := scannerScoreFromSnapshot("TEST", "day", snap)
	if got.GapPercent < 7.5 || got.GapPercent > 8.5 {
		t.Fatalf("expected ~8%% gap from daily open vs previous close, got %.2f", got.GapPercent)
	}
	if got.SpreadPercent <= 0 || got.RelativeVolume < 1.9 {
		t.Fatalf("expected spread and RVOL to be populated: %+v", got)
	}
	if got.Score <= 0 {
		t.Fatalf("expected positive scanner score: %+v", got)
	}
}

func TestV15DynamicAllocatorUsesAlpacaThenFinnhubOverflow(t *testing.T) {
	app := newTestApplication(t)
	var day, swing []string
	for i := 0; i < 42; i++ {
		day = append(day, fmt.Sprintf("D%03d", i))
	}
	for i := 0; i < 20; i++ {
		swing = append(swing, fmt.Sprintf("S%03d", i))
	}
	for i := range app.state.Watchlists {
		switch app.state.Watchlists[i].ID {
		case "day":
			app.state.Watchlists[i].Symbols = day
		case "swing":
			app.state.Watchlists[i].Symbols = swing
		}
	}
	alloc := app.engine.multiFeedAllocation()
	if len(alloc.Alpaca) != alpacaActiveTarget {
		t.Fatalf("Alpaca primary active=%d want=%d", len(alloc.Alpaca), alpacaActiveTarget)
	}
	if len(alloc.Finnhub) == 0 || len(alloc.Finnhub) > finnhubActiveTarget {
		t.Fatalf("Finnhub overflow=%d", len(alloc.Finnhub))
	}
	seen := map[string]bool{}
	for _, x := range alloc.Finnhub {
		seen[x] = true
	}
	for _, x := range alloc.Alpaca {
		if seen[x] {
			t.Fatalf("duplicate live allocation %s", x)
		}
	}
	for _, pinned := range []string{"GLD", "SLV", "USO"} {
		if !containsLiveSymbol(alloc.Alpaca, pinned) {
			t.Fatalf("pinned tradable %s missing from preferred Alpaca live pool", pinned)
		}
	}
}

func TestV1433PassiveMarketContextUsesSnapshotPool(t *testing.T) {
	app := newTestApplication(t)
	alloc := app.engine.multiFeedAllocation()
	containsSym := func(xs []string, s string) bool {
		for _, x := range xs {
			if x == s {
				return true
			}
		}
		return false
	}
	if containsSym(alloc.Finnhub, "EWJ") || containsSym(alloc.Alpaca, "EWJ") {
		t.Fatalf("passive EWJ context should not consume live slots: %+v", alloc)
	}
	if !containsSym(alloc.Snapshot, "EWJ") {
		t.Fatalf("EWJ should remain in snapshot pool: %+v", alloc.Snapshot)
	}
	for _, pinned := range []string{"GLD", "SLV", "USO"} {
		if !containsSym(alloc.Alpaca, pinned) {
			t.Fatalf("%s should stay pinned in the preferred Alpaca live pool", pinned)
		}
	}
}

func TestV1433AllocatorPrioritizesDayBeforeSwing(t *testing.T) {
	app := newTestApplication(t)
	var day, swing []string
	for i := 0; i < 35; i++ {
		day = append(day, fmt.Sprintf("DAY%02d", i))
	}
	for i := 0; i < 25; i++ {
		swing = append(swing, fmt.Sprintf("SW%02d", i))
	}
	for i := range app.state.Watchlists {
		switch app.state.Watchlists[i].ID {
		case "day":
			app.state.Watchlists[i].Symbols = day
		case "swing":
			app.state.Watchlists[i].Symbols = swing
		case "long":
			app.state.Watchlists[i].Symbols = nil
		case "discovery":
			app.state.Watchlists[i].Symbols = nil
		}
	}
	alloc := app.engine.multiFeedAllocation()
	live := map[string]bool{}
	for _, x := range alloc.Alpaca {
		live[x] = true
	}
	for _, x := range alloc.Finnhub {
		live[x] = true
	}
	for _, x := range day {
		if !live[x] {
			t.Fatalf("day symbol %s should be promoted before lower-priority overflow", x)
		}
	}
	if len(alloc.Finnhub) == 0 {
		t.Fatal("expected lower-priority overflow into Finnhub")
	}
}

func TestV123SnapshotDoesNotOverwriteRecentAlpacaIEXStream(t *testing.T) {
	app := newTestApplication(t)
	now := time.Now().UnixMilli()
	app.engine.mu.Lock()
	app.engine.quotes["SPY"] = Quote{Symbol: "SPY", Price: 640.25, Source: "alpaca-iex-websocket-trade", ProviderTimestamp: now, UpdatedAt: now, PreviousClose: 638}
	app.engine.lastAlpacaStreamAt = now
	app.engine.lastAlpacaStreamSymbol = "SPY"
	app.engine.mu.Unlock()
	var snap alpacaLiveSnapshot
	snap.LatestTrade.Price = 639.50
	snap.LatestTrade.Time = time.Now().UTC().Format(time.RFC3339Nano)
	applied := app.engine.mergeAlpacaLiveSnapshot("SPY", 639.50, now, snap, "iex", "trade")
	got := app.engine.Snapshot().Quotes["SPY"]
	if applied || got.Price != 640.25 || got.Source != "alpaca-iex-websocket-trade" {
		t.Fatalf("REST snapshot overwrote recent Alpaca IEX stream: applied=%v quote=%+v", applied, got)
	}
}

func TestV123FeedDiagnosticsExposeAlpacaIEXStream(t *testing.T) {
	app := newTestApplication(t)
	now := time.Now().UnixMilli()
	app.engine.mu.Lock()
	app.engine.status = "running"
	app.engine.mode = "live"
	app.engine.alpacaWebSocketConnected = true
	app.engine.alpacaSubscribedSymbols["SPY"] = true
	app.engine.lastAlpacaStreamAt = now
	app.engine.lastAlpacaStreamSymbol = "SPY"
	app.engine.mu.Unlock()
	// FeedState is session-aware; on an active extended/regular session the IEX
	// stream must be exposed as streaming. On a closed session the diagnostics
	// still must expose the connected provider and latest stream timestamp.
	snap := app.engine.Snapshot()
	if !snap.Feed.AlpacaWebSocketConnected || len(snap.Feed.AlpacaSubscribedSymbols) != 1 || snap.Feed.LastAlpacaStreamSymbol != "SPY" {
		t.Fatalf("missing Alpaca IEX diagnostics: %+v", snap.Feed)
	}
	session := marketSessionET(time.Now())
	if session != "closed" && session != "weekend" && session != "overnight" && snap.Feed.FeedState != "streaming" {
		t.Fatalf("expected active-session IEX feedState, got %q for %q", snap.Feed.FeedState, session)
	}
}

func TestV1233SettingsSavePersistsTimestamp(t *testing.T) {
	app := newTestApplication(t)
	if app.state.SettingsSavedAt != 0 {
		t.Fatalf("expected fresh test state to have no settings timestamp, got %d", app.state.SettingsSavedAt)
	}
	payload, err := json.Marshal(map[string]any{"settings": app.state.Settings})
	if err != nil {
		t.Fatal(err)
	}
	req := httptest.NewRequest("POST", "/api/settings/save", strings.NewReader(string(payload)))
	rr := httptest.NewRecorder()
	app.handleSettingsSave(rr, req)
	if rr.Code != 200 {
		t.Fatalf("code=%d body=%s", rr.Code, rr.Body.String())
	}
	var pub PublicState
	if err := json.Unmarshal(rr.Body.Bytes(), &pub); err != nil {
		t.Fatal(err)
	}
	if pub.SettingsSavedAt <= 0 || app.state.SettingsSavedAt != pub.SettingsSavedAt {
		t.Fatalf("settings save timestamp missing: public=%d state=%d", pub.SettingsSavedAt, app.state.SettingsSavedAt)
	}
	data, err := os.ReadFile(app.statePath())
	if err != nil {
		t.Fatal(err)
	}
	var persisted AppState
	if err := json.Unmarshal(data, &persisted); err != nil {
		t.Fatal(err)
	}
	if persisted.SettingsSavedAt != pub.SettingsSavedAt {
		t.Fatalf("timestamp was not persisted: file=%d public=%d", persisted.SettingsSavedAt, pub.SettingsSavedAt)
	}
}

func TestV1235DefaultAIProviderIsGroq(t *testing.T) {
	st := defaultState()
	if st.Settings.AIProvider != "groq" {
		t.Fatalf("expected groq default, got %q", st.Settings.AIProvider)
	}
	if st.Settings.GroqModel != "openai/gpt-oss-120b" || st.Settings.GeminiModel != "gemini-3.1-flash-lite" {
		t.Fatalf("unexpected default AI models: %+v", st.Settings)
	}
}

func TestV1235SettingsSavePersistsAIProviderAndKeys(t *testing.T) {
	app := newTestApplication(t)
	settings := app.state.Settings
	settings.AIProvider = "gemini"
	payload, err := json.Marshal(map[string]any{
		"settings":      settings,
		"groqKey":       "gsk_test_123456789",
		"geminiKey":     "gemini_test_123456789",
		"openRouterKey": "sk-or-test-123456789",
	})
	if err != nil {
		t.Fatal(err)
	}
	req := httptest.NewRequest("POST", "/api/settings/save", strings.NewReader(string(payload)))
	rr := httptest.NewRecorder()
	app.handleSettingsSave(rr, req)
	if rr.Code != 200 {
		t.Fatalf("code=%d body=%s", rr.Code, rr.Body.String())
	}
	app.mu.RLock()
	defer app.mu.RUnlock()
	if app.state.Settings.AIProvider != "gemini" {
		t.Fatalf("AI provider not saved: %q", app.state.Settings.AIProvider)
	}
	if app.secrets.Groq == "" || app.secrets.Gemini == "" || app.secrets.OpenRouter == "" {
		t.Fatalf("AI keys were not persisted: %+v", app.secrets)
	}
	pub := app.publicStateLockedForUser(bootstrapOwnerID)
	if !pub.HasGroqKey || !pub.HasGeminiKey || !pub.HasOpenRouterKey {
		t.Fatalf("public key flags missing: %+v", pub)
	}
}

func TestV1235AIDoesNotSilentlyFallback(t *testing.T) {
	app := newTestApplication(t)
	app.state.Settings.AIProvider = "gemini"
	_, err := app.GenerateAIForUser(context.Background(), bootstrapOwnerID, AIRequest{Question: "test"})
	if err == nil || !strings.Contains(strings.ToLower(err.Error()), "gemini") {
		t.Fatalf("expected Gemini missing-key error without provider fallback, got %v", err)
	}
	app.state.Settings.AIProvider = "groq"
	_, err = app.GenerateAIForUser(context.Background(), bootstrapOwnerID, AIRequest{Question: "test"})
	if err == nil || !strings.Contains(strings.ToLower(err.Error()), "groq") {
		t.Fatalf("expected Groq missing-key error without provider fallback, got %v", err)
	}
}

func TestV1237SECFormMeaningsAreCompactAndUseful(t *testing.T) {
	cases := map[string]string{
		"8-K":      "Material Company Update",
		"10-Q":     "Quarterly Report",
		"10-K":     "Annual Report",
		"4":        "Insider Transaction",
		"13F-HR":   "Institutional Holdings Report",
		"SC 13D/A": "Major Shareholder Disclosure",
		"S-3":      "Shelf Registration / Offering",
		"424B5":    "Offering Terms",
		"144":      "Planned Insider Sale",
		"DEF 14A":  "Proxy Statement",
	}
	for form, want := range cases {
		got, _ := secFormMeaning(form)
		if got != want {
			t.Fatalf("%s meaning=%q want=%q", form, got, want)
		}
	}
}

func TestV1237OpenRouterModesAndFallbacks(t *testing.T) {
	fast := openRouterConfig("fast", "")
	if fast.Primary != "openai/gpt-5.6-luna" || len(fast.Fallback) != 2 || fast.Fallback[0] != "x-ai/grok-4.20" {
		t.Fatalf("unexpected fast route: %+v", fast)
	}
	balanced := openRouterConfig("balanced", "")
	if balanced.Primary != "x-ai/grok-4.20" || len(balanced.Fallback) != 2 {
		t.Fatalf("unexpected balanced route: %+v", balanced)
	}
	powerful := openRouterConfig("powerful", "")
	if powerful.Primary != "anthropic/claude-sonnet-5" || len(powerful.Fallback) != 2 {
		t.Fatalf("unexpected powerful route: %+v", powerful)
	}
	specific := openRouterConfig("specific", "x-ai/grok-4.5")
	if specific.Primary != "x-ai/grok-4.5" || len(specific.Fallback) != 0 {
		t.Fatalf("unexpected specific route: %+v", specific)
	}
	if allowedOpenRouterModel("google/gemini-3.1-pro") {
		t.Fatal("OpenRouter specific mode must stay restricted to GPT, Grok, or Claude families")
	}
}

func TestV1436Form4ClassifiesBuySellAndOtherWithoutFalseBullishLabel(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/xml")
		_, _ = io.WriteString(w, `<?xml version="1.0"?><ownershipDocument>
<reportingOwner><reportingOwnerId><rptOwnerName>Jane Example</rptOwnerName></reportingOwnerId><reportingOwnerRelationship><isOfficer>1</isOfficer><officerTitle>CFO</officerTitle></reportingOwnerRelationship></reportingOwner>
<nonDerivativeTable>
<nonDerivativeTransaction><transactionDate><value>2026-08-08</value></transactionDate><transactionCoding><transactionCode>P</transactionCode></transactionCoding><transactionAmounts><transactionShares><value>1000</value></transactionShares><transactionPricePerShare><value>25.50</value></transactionPricePerShare><transactionAcquiredDisposedCode><value>A</value></transactionAcquiredDisposedCode></transactionAmounts><postTransactionAmounts><sharesOwnedFollowingTransaction><value>12000</value></sharesOwnedFollowingTransaction></postTransactionAmounts></nonDerivativeTransaction>
</nonDerivativeTable></ownershipDocument>`)
	}))
	defer server.Close()
	item := FilingItem{Form: "4", Meaning: "Insider Transaction", Category: "insider"}
	enrichForm4(context.Background(), server.Client(), nil, server.URL, &item)
	if item.Signal != "Buy" || item.Meaning != "Insider Buy" || item.TransactionCode != "P" || item.TransactionType != "Open-market/private purchase" {
		t.Fatalf("unexpected classified buy: %+v", item)
	}
	if item.Actor != "Jane Example" || item.Role != "CFO" || item.Shares != 1000 || item.Price != 25.5 || item.Value != 25500 || item.OwnershipAfter != 12000 {
		t.Fatalf("missing Form 4 details: %+v", item)
	}

	serverOther := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/xml")
		_, _ = io.WriteString(w, `<?xml version="1.0"?><ownershipDocument><nonDerivativeTable><nonDerivativeTransaction><transactionDate><value>2026-08-08</value></transactionDate><transactionCoding><transactionCode>A</transactionCode></transactionCoding><transactionAmounts><transactionShares><value>500</value></transactionShares><transactionPricePerShare><value>0</value></transactionPricePerShare><transactionAcquiredDisposedCode><value>A</value></transactionAcquiredDisposedCode></transactionAmounts><postTransactionAmounts><sharesOwnedFollowingTransaction><value>2500</value></sharesOwnedFollowingTransaction></postTransactionAmounts></nonDerivativeTransaction></nonDerivativeTable></ownershipDocument>`)
	}))
	defer serverOther.Close()
	other := FilingItem{Form: "4", Meaning: "Insider Transaction", Category: "insider"}
	enrichForm4(context.Background(), serverOther.Client(), nil, serverOther.URL, &other)
	if other.Signal != "Other" || other.Meaning != "Insider Other" || other.TransactionCode != "A" || other.TransactionType != "Award / grant" {
		t.Fatalf("award/grant must remain OTHER, got %+v", other)
	}
	if strings.EqualFold(other.Signal, "buy") || strings.Contains(strings.ToUpper(other.Meaning), "BUY") {
		t.Fatalf("award/grant was mislabeled as insider buy: %+v", other)
	}
}

func TestV1237SECIntelligenceSummaryIsSignalFirst(t *testing.T) {
	now := time.Now()
	items := []FilingItem{
		{Symbol: "META", Form: "4", FiledAt: now.AddDate(0, 0, -2).Format("2006-01-02"), Meaning: "Insider trade", Category: "insider", Signal: "buy", Role: "CEO", Value: 420000},
		{Symbol: "META", Form: "8-K", FiledAt: now.AddDate(0, 0, -5).Format("2006-01-02"), Meaning: "Material event", Category: "material"},
		{Symbol: "META", Form: "SC 13G", FiledAt: now.AddDate(0, 0, -7).Format("2006-01-02"), Meaning: "Major ownership", Category: "ownership"},
		{Symbol: "META", Form: "S-3", FiledAt: now.AddDate(0, 0, -12).Format("2006-01-02"), Meaning: "Shelf registration", Category: "offering"},
	}
	got := buildSECIntelligence("META", items)
	if got.InsiderBuys != 1 || got.InsiderSells != 0 || got.InsiderBuyValue != 420000 {
		t.Fatalf("unexpected insider summary: %+v", got)
	}
	if got.MaterialEvents != 1 || got.OwnershipChanges != 1 || got.OfferingRisk != "High" {
		t.Fatalf("unexpected SEC signals: %+v", got)
	}
	if got.Institutional != "13F · Quarterly" {
		t.Fatalf("13F must be explicitly quarterly, got %q", got.Institutional)
	}
	if len(got.Signals) == 0 || got.Signals[0].Label != "BUY" {
		t.Fatalf("expected compact insider signal first: %+v", got.Signals)
	}
}

func TestV1301MarketSessionBoundaryContext(t *testing.T) {
	loc, err := time.LoadLocation("America/New_York")
	if err != nil {
		t.Fatal(err)
	}
	cases := []struct {
		name                    string
		at                      time.Time
		wantSession, wantAction string
		wantWeekday             time.Weekday
		wantHour, wantMinute    int
	}{
		{"premarket", time.Date(2026, 8, 7, 8, 0, 0, 0, loc), "pre-market", "Opens", time.Friday, 9, 30},
		{"regular", time.Date(2026, 8, 7, 11, 0, 0, 0, loc), "regular", "Closes", time.Friday, 16, 0},
		{"afterhours", time.Date(2026, 8, 7, 18, 0, 0, 0, loc), "after-hours", "Ends", time.Friday, 20, 0},
		{"early-close", time.Date(2026, 11, 27, 12, 0, 0, 0, loc), "regular", "Closes", time.Friday, 13, 0},
		{"friday-weekend", time.Date(2026, 8, 7, 21, 0, 0, 0, loc), "weekend", "Next overnight", time.Sunday, 20, 0},
		{"saturday-weekend", time.Date(2026, 8, 8, 13, 0, 0, 0, loc), "weekend", "Next overnight", time.Sunday, 20, 0},
		{"sunday-before-overnight", time.Date(2026, 8, 9, 19, 0, 0, 0, loc), "weekend", "Next overnight", time.Sunday, 20, 0},
		{"sunday-overnight", time.Date(2026, 8, 9, 20, 30, 0, 0, loc), "overnight", "Pre-market", time.Monday, 4, 0},
		{"holiday-weekend-roll", time.Date(2026, 9, 6, 18, 0, 0, 0, loc), "weekend", "Next overnight", time.Monday, 20, 0},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := marketSessionET(tc.at); got != tc.wantSession {
				t.Fatalf("session=%q want %q", got, tc.wantSession)
			}
			ms, action := marketSessionBoundaryET(tc.at)
			if action != tc.wantAction {
				t.Fatalf("action=%q want %q", action, tc.wantAction)
			}
			b := time.UnixMilli(ms).In(loc)
			if b.Weekday() != tc.wantWeekday || b.Hour() != tc.wantHour || b.Minute() != tc.wantMinute {
				t.Fatalf("boundary=%s want %s %02d:%02d ET", b, tc.wantWeekday, tc.wantHour, tc.wantMinute)
			}
		})
	}
}

func TestV1327LegacyOptionalWorkspaceFieldsAreIgnoredAndPurged(t *testing.T) {
	raw := []byte(`{"version":14,"settings":{"dataMode":"demo","showAlerts":true,"showPortfolio":true,"showJournal":true},"watchlists":[{"id":"swing","name":"Swing","symbols":["NVDA"]},{"id":"day","name":"Day","symbols":["NVDA"]},{"id":"long","name":"Long","symbols":["SPY"]},{"id":"discovery","name":"Discovery","symbols":[]}],"ui":{"scopeType":"watchlist","watchlistId":"swing","selectedTicker":"NVDA"},"alerts":[{"id":"old-alert"}],"positions":[{"id":"old-position"}],"journal":[{"id":"old-journal"}]}`)
	var st AppState
	if err := json.Unmarshal(raw, &st); err != nil {
		t.Fatal(err)
	}
	st = mergeState(st)
	if st.Version != 17 {
		t.Fatalf("expected canonical state version 17, got %d", st.Version)
	}
	out, err := json.Marshal(st)
	if err != nil {
		t.Fatal(err)
	}
	text := string(out)
	for _, forbidden := range []string{`"showAlerts"`, `"showPortfolio"`, `"showJournal"`, `"alerts"`, `"positions"`, `"journal"`} {
		if strings.Contains(text, forbidden) {
			t.Fatalf("legacy field survived canonical state: %s", forbidden)
		}
	}
}

func TestV1327MaintenanceRunIsDiagnosticAndPreservesMarketCacheAndSubscriptions(t *testing.T) {
	app := newTestApplication(t)
	app.mu.Lock()
	if err := app.saveLocked(); err != nil {
		app.mu.Unlock()
		t.Fatal(err)
	}
	app.mu.Unlock()

	app.engine.mu.Lock()
	app.engine.quotes["SPY"] = Quote{Symbol: "SPY", Price: 600, UpdatedAt: time.Now().UnixMilli()}
	app.engine.subscribedSymbols["SPY"] = true
	app.engine.alpacaSubscribedSymbols["QQQ"] = true
	app.engine.mu.Unlock()
	if err := app.engine.saveCache(); err != nil {
		t.Fatal(err)
	}
	before, err := os.Stat(app.cachePath())
	if err != nil {
		t.Fatal(err)
	}
	beforeFinnhub := len(app.engine.subscribedSymbols)
	beforeAlpaca := len(app.engine.alpacaSubscribedSymbols)

	req := httptest.NewRequest("POST", "/api/maintenance/run", strings.NewReader(`{}`))
	rr := httptest.NewRecorder()
	app.handleMaintenanceRun(rr, req)
	if rr.Code != 200 {
		t.Fatalf("status=%d body=%s", rr.Code, rr.Body.String())
	}
	var report MaintenanceReport
	if err := json.Unmarshal(rr.Body.Bytes(), &report); err != nil {
		t.Fatal(err)
	}
	if report.Version != appVersion || report.BuildID != buildID || len(report.Checks) < 10 {
		t.Fatalf("unexpected report: %+v", report)
	}
	if app.state.MaintenanceLastRun == 0 {
		t.Fatal("maintenance last-run timestamp was not persisted")
	}
	after, err := os.Stat(app.cachePath())
	if err != nil {
		t.Fatal(err)
	}
	if after.Size() != before.Size() {
		t.Fatalf("weekly maintenance modified cache size: before=%d after=%d", before.Size(), after.Size())
	}
	app.engine.mu.RLock()
	afterFinnhub := len(app.engine.subscribedSymbols)
	afterAlpaca := len(app.engine.alpacaSubscribedSymbols)
	app.engine.mu.RUnlock()
	if beforeFinnhub != afterFinnhub || beforeAlpaca != afterAlpaca {
		t.Fatalf("weekly maintenance changed subscriptions: Finnhub %d->%d Alpaca %d->%d", beforeFinnhub, afterFinnhub, beforeAlpaca, afterAlpaca)
	}
}

func TestV1327ProfileExportDoesNotContainRemovedWorkspaceData(t *testing.T) {
	app := newTestApplication(t)
	req := httptest.NewRequest("GET", "/api/profile/export", nil)
	rr := httptest.NewRecorder()
	app.handleExport(rr, req)
	if rr.Code != 200 {
		t.Fatalf("status=%d", rr.Code)
	}
	body := rr.Body.String()
	for _, forbidden := range []string{`"alerts"`, `"positions"`, `"journal"`, `"showAlerts"`, `"showPortfolio"`, `"showJournal"`} {
		if strings.Contains(body, forbidden) {
			t.Fatalf("removed workspace field present in export: %s", forbidden)
		}
	}
}

func TestV1324UpgradePreservesV1323TradingProfile(t *testing.T) {
	dir := t.TempDir()
	legacy := defaultState()
	legacy.Settings.DataMode = "live"
	legacy.Settings.OvernightDataMode = "auto"
	legacy.Settings.SECEmail = "qa-upgrade@example.com"
	legacy.Settings.SwingEnabled = true
	legacy.Settings.DayEnabled = false
	legacy.Settings.LongEnabled = true
	legacy.UI.SelectedTicker = "AMD"
	for i := range legacy.Watchlists {
		switch legacy.Watchlists[i].ID {
		case "swing":
			legacy.Watchlists[i].Symbols = []string{"AMD", "NVDA", "META"}
		case "day":
			legacy.Watchlists[i].Symbols = []string{"TSLA", "PLTR"}
		case "long":
			legacy.Watchlists[i].Symbols = []string{"SPY", "QQQ", "AAPL"}
		case "discovery":
			legacy.Watchlists[i].Symbols = []string{"AVGO"}
		}
	}
	legacy.MaintenanceLastRun = 1760000000000
	data, err := json.MarshalIndent(legacy, "", "  ")
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "state.json"), data, 0600); err != nil {
		t.Fatal(err)
	}
	secrets := Secrets{Finnhub: "fh-secret", AlpacaKey: "ak", AlpacaSecret: "as", Groq: "gq"}
	secretData, _ := json.Marshal(secrets)
	if err := os.WriteFile(filepath.Join(dir, "secrets.json"), secretData, 0600); err != nil {
		t.Fatal(err)
	}
	app := &Application{configDir: dir, hub: NewHub(), sessionKey: "upgrade-test"}
	app.load()
	app.engine = NewEngine(app)
	if app.state.Settings.DataMode != "live" || app.state.Settings.OvernightDataMode != "auto" || app.state.Settings.SECEmail != "qa-upgrade@example.com" {
		t.Fatalf("settings changed on upgrade: %+v", app.state.Settings)
	}
	if app.state.Settings.DayEnabled || !app.state.Settings.SwingEnabled || !app.state.Settings.LongEnabled {
		t.Fatalf("engine toggles changed: %+v", app.state.Settings)
	}
	if app.state.UI.SelectedTicker != "AMD" {
		t.Fatalf("selected ticker changed: %+v", app.state.UI)
	}
	checks := map[string][]string{"swing": {"AMD", "NVDA", "META"}, "day": {"TSLA", "PLTR"}, "long": {"SPY", "QQQ", "AAPL"}, "discovery": {"AVGO"}}
	for id, want := range checks {
		wl, ok := watchlistValueByID(app.state.Watchlists, id)
		if !ok || !reflect.DeepEqual(wl.Symbols, want) {
			t.Fatalf("%s watchlist changed: got=%v want=%v", id, wl.Symbols, want)
		}
	}
	if app.secrets.Finnhub != "fh-secret" || app.secrets.AlpacaKey != "ak" || app.secrets.AlpacaSecret != "as" || app.secrets.Groq != "gq" {
		t.Fatal("provider credentials were not preserved")
	}
	if app.state.MaintenanceLastRun != legacy.MaintenanceLastRun {
		t.Fatal("maintenance timestamp was not preserved")
	}
}

func TestV1327MasterMarketUniverseIncludesMarketInstrumentsAndSeparatesVIX(t *testing.T) {
	st := defaultState()
	master := masterSymbolsFromState(st)
	for _, symbol := range masterMarketSymbols {
		if !contains(master, symbol) {
			t.Fatalf("master store missing Market Instrument %s: %v", symbol, master)
		}
	}
	if contains(master, "VIX") {
		t.Fatalf("VIX must not be in the equity/ETF Master Symbol Store: %v", master)
	}
	special := specialIndexSymbolsFromState(st)
	if len(special) != 1 || special[0] != "VIX" {
		t.Fatalf("unexpected special-index universe: %v", special)
	}
	required := requiredSymbolsFromState(st)
	if !contains(required, "VIX") || !contains(required, "SPY") || !contains(required, "QQQ") || !contains(required, "TLT") {
		t.Fatalf("required universe did not combine master + special stores: %v", required)
	}
	countVIX := 0
	for _, symbol := range required {
		if symbol == "VIX" {
			countVIX++
		}
	}
	if countVIX != 1 {
		t.Fatalf("VIX should appear exactly once in required universe, got %d", countVIX)
	}
}

func TestV1327HistorySpecsHydrateDeskSymbolAndCommonRegimeInputs(t *testing.T) {
	st := defaultState()
	for i := range st.Watchlists {
		switch st.Watchlists[i].ID {
		case "day":
			st.Watchlists[i].Symbols = []string{"DAYX"}
		case "swing":
			st.Watchlists[i].Symbols = []string{"SWGX"}
		case "long":
			st.Watchlists[i].Symbols = []string{"LNGX"}
		}
	}
	specs := historySpecsForState(st, nil)
	byName := map[string][]string{}
	for _, spec := range specs {
		byName[spec.Name] = spec.Symbols
		if contains(spec.Symbols, "VIX") {
			t.Fatalf("VIX must remain outside Alpaca equity history specs: %+v", spec)
		}
	}
	for _, symbol := range []string{"DAYX", "SPY", "QQQ", "IWM", "XLK"} {
		if !contains(byName["intraday"], symbol) {
			t.Fatalf("intraday hydration missing %s: %v", symbol, byName["intraday"])
		}
	}
	for _, symbol := range []string{"SWGX", "LNGX", "SPY", "QQQ", "IWM", "XLK", "TLT"} {
		if !contains(byName["daily"], symbol) {
			t.Fatalf("daily hydration missing %s: %v", symbol, byName["daily"])
		}
	}
	if !contains(byName["weekly"], "LNGX") {
		t.Fatalf("weekly hydration missing Long-Term symbol: %v", byName["weekly"])
	}

	one := historySpecsForState(st, []string{"SWGX"})
	if len(one) != 1 || one[0].Name != "daily" || len(one[0].Symbols) != 1 || one[0].Symbols[0] != "SWGX" {
		t.Fatalf("new Swing symbol should request only its required daily bars: %+v", one)
	}
	dayOnly := historySpecsForState(st, []string{"DAYX"})
	if len(dayOnly) != 1 || dayOnly[0].Name != "intraday" || len(dayOnly[0].Symbols) != 1 || dayOnly[0].Symbols[0] != "DAYX" {
		t.Fatalf("new Day symbol should request only intraday bars: %+v", dayOnly)
	}
	longOnly := historySpecsForState(st, []string{"LNGX"})
	if len(longOnly) != 2 || longOnly[0].Name != "daily" || longOnly[1].Name != "weekly" {
		t.Fatalf("new Long-Term symbol should request daily + weekly bars: %+v", longOnly)
	}
}

func TestV1327HistoryHydrationTriggersAreWired(t *testing.T) {
	src := productionGoSourceForTest(t)
	for _, marker := range []string{
		"e.requestHistoryHydration(symbol)",
		"a.engine.requestHistoryHydration()",
		"if running && in.Enabled",
		"e.sessionAwareHistoryLoop(ctx)",
		"func (e *Engine) sessionAwareHistoryLoop",
	} {
		if !strings.Contains(src, marker) {
			t.Fatalf("missing history-hydration wiring marker %q", marker)
		}
	}
}

func TestV1327LiveCacheClearImmediatelyRehydratesHistoricalBars(t *testing.T) {
	app := newTestApplication(t)
	app.mu.Lock()
	app.state.Settings.DataMode = "live"
	app.secrets.AlpacaKey = "test-key"
	app.secrets.AlpacaSecret = "test-secret"
	app.mu.Unlock()
	app.engine.mu.Lock()
	app.engine.status = "running"
	app.engine.mode = "live"
	app.engine.bars["SPY"] = map[string][]Bar{"daily": {{T: 1, C: 500}}}
	app.engine.mu.Unlock()
	if err := app.engine.saveCache(); err != nil {
		t.Fatal(err)
	}

	oldBase := alpacaDataBaseURL
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v2/stocks/bars" {
			http.NotFound(w, r)
			return
		}
		syms := strings.Split(r.URL.Query().Get("symbols"), ",")
		payload := map[string]any{"bars": map[string]any{}}
		out := payload["bars"].(map[string]any)
		for _, sym := range syms {
			if sym == "" || sym == "VIX" {
				continue
			}
			out[sym] = []map[string]any{{"c": 101.0, "h": 102.0, "l": 99.0, "o": 100.0, "v": 1000.0, "t": "2026-08-07T20:00:00Z"}}
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(payload)
	}))
	defer server.Close()
	alpacaDataBaseURL = server.URL
	defer func() { alpacaDataBaseURL = oldBase }()

	req := httptest.NewRequest("POST", "/api/cache/clear", strings.NewReader(`{}`))
	rr := httptest.NewRecorder()
	app.handleCacheClear(rr, req)
	if rr.Code != 200 {
		t.Fatalf("code=%d body=%s", rr.Code, rr.Body.String())
	}
	deadline := time.Now().Add(3 * time.Second)
	for time.Now().Before(deadline) {
		snap := app.engine.Snapshot()
		if len(snap.Bars["SPY"]["intraday"]) > 0 && len(snap.Bars["SPY"]["daily"]) > 0 {
			if app.state.LastCacheCleared <= 0 {
				t.Fatal("cache clear timestamp missing after live hydration")
			}
			return
		}
		time.Sleep(20 * time.Millisecond)
	}
	t.Fatalf("required bars did not rehydrate immediately after cache clear: %+v", app.engine.Snapshot().Bars["SPY"])
}

func TestV1329StableResearchDatasetsPersistAcrossCacheRestart(t *testing.T) {
	app := newTestApplication(t)
	now := time.Now().UnixMilli()
	app.engine.mu.Lock()
	app.engine.news = []NewsItem{{Headline: "Material product update", Source: "QA", Symbols: []string{"AAA"}, Datetime: now / 1000}}
	app.engine.earnings = []EarningsItem{{Symbol: "AAA", Date: "2026-08-20", Hour: "amc"}}
	app.engine.filings = []FilingItem{{ID: "aaa-10q", Symbol: "AAA", Form: "10-Q", FiledAt: "2026-08-07", URL: "https://example.com/10q"}}
	app.engine.secIntelligence = map[string]SECIntelligenceSummary{"AAA": {Symbol: "AAA", MaterialEvents: 1, LatestForm: "10-Q"}}
	app.engine.scanner = ScannerState{Mode: "swing", Status: "ready", Results: []ScannerResult{{Symbol: "AAA", Score: 77}}, UpdatedAt: now}
	app.engine.lastUpdated["news"] = now - 1000
	app.engine.lastUpdated["earnings"] = now - 2000
	app.engine.lastUpdated["filings"] = now - 3000
	app.engine.mu.Unlock()
	if err := app.engine.saveCache(); err != nil {
		t.Fatal(err)
	}

	reloaded := NewEngine(app)
	reloaded.mu.RLock()
	defer reloaded.mu.RUnlock()
	if len(reloaded.news) != 1 || reloaded.news[0].Headline != "Material product update" {
		t.Fatalf("news cache not restored: %+v", reloaded.news)
	}
	if len(reloaded.earnings) != 1 || reloaded.earnings[0].Symbol != "AAA" {
		t.Fatalf("earnings cache not restored: %+v", reloaded.earnings)
	}
	if len(reloaded.filings) != 1 || reloaded.filings[0].Form != "10-Q" {
		t.Fatalf("filings cache not restored: %+v", reloaded.filings)
	}
	if reloaded.secIntelligence["AAA"].LatestForm != "10-Q" {
		t.Fatalf("SEC intelligence cache not restored: %+v", reloaded.secIntelligence)
	}
	if reloaded.scanner.Mode != "swing" || len(reloaded.scanner.Results) != 1 {
		t.Fatalf("scanner cache not restored: %+v", reloaded.scanner)
	}
	if reloaded.lastUpdated["news"] != now-1000 {
		t.Fatalf("last-updated metadata not restored: %+v", reloaded.lastUpdated)
	}
}

func TestV1329StructuredAIParserNormalizesAndBoundsOutput(t *testing.T) {
	raw := `{"verdict":"favorable","confidence":140,"reasons":["one","two","three","four"],"risks":["a","b","c","d"],"catalyst":"earnings","bestFitHorizon":"Swing","nextAction":"Open Swing Desk","summary":"Evidence is constructive.","details":"bounded evidence"}`
	out := parseAIStructuredPayload(raw)
	if out.Verdict != "FAVORABLE" || out.Confidence != 100 || out.BestFitHorizon != "swing" {
		t.Fatalf("normalization failed: %+v", out)
	}
	if len(out.Reasons) != 3 || len(out.Risks) != 3 {
		t.Fatalf("AI reasons/risks were not bounded: %+v", out)
	}

	fenced := "```json\n" + raw + "\n```"
	out = parseAIStructuredPayload(fenced)
	if out.Verdict != "FAVORABLE" || out.BestFitHorizon != "swing" {
		t.Fatalf("fenced JSON was not parsed: %+v", out)
	}

	fallback := parseAIStructuredPayload("This provider ignored the JSON schema.")
	if fallback.Verdict != "INFORMATIONAL" || fallback.Confidence != 0 || fallback.Details == "" {
		t.Fatalf("unstructured fallback unsafe: %+v", fallback)
	}
}

func TestV1329CacheRefreshEndpointIsNonDestructiveInDemo(t *testing.T) {
	app := newTestApplication(t)
	app.engine.mu.Lock()
	app.engine.status = "running"
	app.engine.mode = "demo"
	app.engine.quotes["AAA"] = Quote{Symbol: "AAA", Price: 123}
	app.engine.mu.Unlock()

	req := httptest.NewRequest("POST", "/api/cache/refresh", strings.NewReader(`{}`))
	rr := httptest.NewRecorder()
	app.handleCacheRefresh(rr, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("status=%d body=%s", rr.Code, rr.Body.String())
	}
	if !strings.Contains(rr.Body.String(), "no provider refresh is required") {
		t.Fatalf("unexpected demo refresh response: %s", rr.Body.String())
	}
	if app.engine.Snapshot().Quotes["AAA"].Price != 123 {
		t.Fatal("manual stale-data refresh modified/cleared canonical quote state")
	}
}

func TestV1329CacheClearAlsoClearsPersistedResearchDatasets(t *testing.T) {
	app := newTestApplication(t)
	app.engine.mu.Lock()
	app.engine.news = []NewsItem{{Headline: "x"}}
	app.engine.earnings = []EarningsItem{{Symbol: "AAA"}}
	app.engine.filings = []FilingItem{{ID: "f", Symbol: "AAA"}}
	app.engine.secIntelligence = map[string]SECIntelligenceSummary{"AAA": {Symbol: "AAA"}}
	app.engine.scanner = ScannerState{Mode: "day", Status: "ready", Results: []ScannerResult{{Symbol: "AAA", Score: 80}}}
	app.engine.mu.Unlock()
	if err := app.engine.saveCache(); err != nil {
		t.Fatal(err)
	}

	req := httptest.NewRequest("POST", "/api/cache/clear", strings.NewReader(`{}`))
	rr := httptest.NewRecorder()
	app.handleCacheClear(rr, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("status=%d body=%s", rr.Code, rr.Body.String())
	}
	snap := app.engine.Snapshot()
	if len(snap.News) != 0 || len(snap.Earnings) != 0 || len(snap.Filings) != 0 || len(snap.SECIntelligence) != 0 || len(snap.Scanner.Results) != 0 {
		t.Fatalf("research/scanner cache not cleared: news=%d earnings=%d filings=%d sec=%d scanner=%d", len(snap.News), len(snap.Earnings), len(snap.Filings), len(snap.SECIntelligence), len(snap.Scanner.Results))
	}
}

func TestV1402MigratesV1330Schema16ProfileWithoutLoss(t *testing.T) {
	dir := t.TempDir()
	before := AppState{
		Version: 16,
		Settings: Settings{
			DataMode: "live", AIProvider: "openrouter", GroqModel: "openai/gpt-oss-120b",
			OpenRouterMode: "specific", OpenRouterSpecificModel: "openai/gpt-5.6-sol", GeminiModel: "gemini-3.1-flash-lite",
			SECEmail: "qa@example.com", AutoStart: true, SignalProfile: "balanced", MarketContext: 17, EarningsPenalty: 8,
			SwingEnabled: true, DayEnabled: false, LongEnabled: true, OvernightDataMode: "indicative",
			SwingWatchlistID: "swing", DayWatchlistID: "day", LongWatchlistID: "long", DiscoveryWatchlistID: "discovery",
		},
		Watchlists: []Watchlist{
			{ID: "swing", Name: "Swing Watchlist", Symbols: []string{"NVDA", "AMD"}},
			{ID: "day", Name: "Day Trade Watchlist", Symbols: []string{"TSLA"}},
			{ID: "long", Name: "Long-Term Watchlist", Symbols: []string{"MSFT", "AAPL"}},
			{ID: "discovery", Name: "Discovery Watchlist", Symbols: []string{"META"}},
		},
		UI:              UIState{ScopeType: "watchlist", WatchlistID: "swing", SelectedTicker: "NVDA"},
		SettingsSavedAt: 1777777000000, MaintenanceLastRun: 1777777100000, LastCacheCleared: 1777777200000,
	}
	secrets := Secrets{Finnhub: "fh-secret", AlpacaKey: "ap-key", AlpacaSecret: "ap-secret", Groq: "g-secret", OpenRouter: "or-secret", Gemini: "gm-secret"}
	stateBytes, _ := json.MarshalIndent(before, "", "  ")
	secretBytes, _ := json.MarshalIndent(secrets, "", "  ")
	if err := os.WriteFile(filepath.Join(dir, "state.json"), stateBytes, 0600); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "secrets.json"), secretBytes, 0600); err != nil {
		t.Fatal(err)
	}

	a := &Application{configDir: dir, hub: NewHub(), sessionKey: "migration-qa"}
	a.load()
	checks := []struct {
		name string
		ok   bool
	}{
		{"schema", a.state.Version == 17},
		{"data mode", a.state.Settings.DataMode == "live"},
		{"overnight mode", a.state.Settings.OvernightDataMode == "indicative"},
		{"v14 global mode default", a.state.Settings.GlobalProviderMode == "auto"},
		{"v14 options mode default", a.state.Settings.OptionsDataMode == "auto"},
		{"v14 macro event mode default", a.state.Settings.MacroEventModeEnabled},
		{"AI provider", a.state.Settings.AIProvider == "openrouter"},
		{"AI mode", a.state.Settings.OpenRouterMode == "specific"},
		{"AI model", a.state.Settings.OpenRouterSpecificModel == "openai/gpt-5.6-sol"},
		{"SEC email", a.state.Settings.SECEmail == "qa@example.com"},
		{"engine states", a.state.Settings.SwingEnabled && !a.state.Settings.DayEnabled && a.state.Settings.LongEnabled},
		{"signal settings", a.state.Settings.SignalProfile == "balanced" && a.state.Settings.MarketContext == 17 && a.state.Settings.EarningsPenalty == 8},
		{"swing watchlist", reflect.DeepEqual(a.state.Watchlists[0].Symbols, []string{"NVDA", "AMD"})},
		{"day watchlist", reflect.DeepEqual(a.state.Watchlists[1].Symbols, []string{"TSLA"})},
		{"long watchlist", reflect.DeepEqual(a.state.Watchlists[2].Symbols, []string{"MSFT", "AAPL"})},
		{"discovery watchlist", reflect.DeepEqual(a.state.Watchlists[3].Symbols, []string{"META"})},
		{"selected ticker", a.state.UI.SelectedTicker == "NVDA" && a.state.UI.WatchlistID == "swing"},
		{"saved timestamps", a.state.SettingsSavedAt == before.SettingsSavedAt && a.state.MaintenanceLastRun == before.MaintenanceLastRun && a.state.LastCacheCleared == before.LastCacheCleared},
		{"Finnhub secret", a.secrets.Finnhub == secrets.Finnhub},
		{"Alpaca secrets", a.secrets.AlpacaKey == secrets.AlpacaKey && a.secrets.AlpacaSecret == secrets.AlpacaSecret},
		{"AI secrets", a.secrets.Groq == secrets.Groq && a.secrets.OpenRouter == secrets.OpenRouter && a.secrets.Gemini == secrets.Gemini},
	}
	for _, c := range checks {
		if !c.ok {
			t.Fatalf("migration check failed: %s", c.name)
		}
	}
}

func TestV1330PreMarketAndWeeklyMaintenanceWindows(t *testing.T) {
	loc, err := time.LoadLocation("America/New_York")
	if err != nil {
		t.Fatal(err)
	}
	cases := []struct {
		name         string
		when         time.Time
		prep, weekly bool
	}{
		{"monday prep", time.Date(2026, time.August, 10, 3, 30, 0, 0, loc), true, false},
		{"too early", time.Date(2026, time.August, 10, 3, 5, 0, 0, loc), false, false},
		{"after prep", time.Date(2026, time.August, 10, 3, 55, 0, 0, loc), false, false},
		{"saturday integrity", time.Date(2026, time.August, 8, 10, 30, 0, 0, loc), false, true},
		{"saturday before window", time.Date(2026, time.August, 8, 9, 30, 0, 0, loc), false, false},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := preMarketPrepWindow(tc.when); got != tc.prep {
				t.Fatalf("prep=%v want %v", got, tc.prep)
			}
			if got := weeklyIntegrityWindow(tc.when); got != tc.weekly {
				t.Fatalf("weekly=%v want %v", got, tc.weekly)
			}
		})
	}
}

func TestV1330RefreshDueUsesDatasetTTL(t *testing.T) {
	app := newTestApplication(t)
	now := time.Now()
	app.engine.mu.Lock()
	app.engine.lastUpdated["history"] = now.Add(-5 * time.Minute).UnixMilli()
	app.engine.lastUpdated["fundamentals"] = now.Add(-25 * time.Hour).UnixMilli()
	app.engine.lastUpdated["news"] = 0
	app.engine.mu.Unlock()
	if app.engine.refreshDue("history", 15*time.Minute, now) {
		t.Fatal("fresh history should not refresh")
	}
	if !app.engine.refreshDue("fundamentals", 24*time.Hour, now) {
		t.Fatal("25h fundamentals should refresh")
	}
	if !app.engine.refreshDue("news", 10*time.Minute, now) {
		t.Fatal("missing timestamp must refresh")
	}
}

func TestV1330CacheFingerprintSkipsUnchangedPhysicalRewrite(t *testing.T) {
	app := newTestApplication(t)
	app.engine.mu.Lock()
	app.engine.quotes["AAA"] = Quote{Symbol: "AAA", Price: 123, UpdatedAt: time.Now().UnixMilli()}
	app.engine.mu.Unlock()
	if err := app.engine.saveCache(); err != nil {
		t.Fatal(err)
	}
	st1, err := os.Stat(app.cachePath())
	if err != nil {
		t.Fatal(err)
	}
	b1, err := os.ReadFile(app.cachePath())
	if err != nil {
		t.Fatal(err)
	}
	time.Sleep(20 * time.Millisecond)
	if err := app.engine.saveCache(); err != nil {
		t.Fatal(err)
	}
	st2, err := os.Stat(app.cachePath())
	if err != nil {
		t.Fatal(err)
	}
	b2, err := os.ReadFile(app.cachePath())
	if err != nil {
		t.Fatal(err)
	}
	if !st1.ModTime().Equal(st2.ModTime()) {
		t.Fatalf("unchanged cache was physically rewritten: %v -> %v", st1.ModTime(), st2.ModTime())
	}
	if !reflect.DeepEqual(b1, b2) {
		t.Fatal("unchanged cache bytes changed")
	}
	app.engine.mu.Lock()
	q := app.engine.quotes["AAA"]
	q.Price = 124
	app.engine.quotes["AAA"] = q
	app.engine.mu.Unlock()
	time.Sleep(20 * time.Millisecond)
	if err := app.engine.saveCache(); err != nil {
		t.Fatal(err)
	}
	st3, _ := os.Stat(app.cachePath())
	if st3.ModTime().Equal(st2.ModTime()) {
		t.Fatal("changed canonical data did not persist")
	}
}

func TestV1330EarningsAndFinancialFilingFingerprintsAreMaterialOnly(t *testing.T) {
	eps := 1.25
	rev := 1000.0
	a := []EarningsItem{{Symbol: "AAA", Date: "2026-08-07", Quarter: 2, Year: 2026, EPSActual: &eps, RevenueActual: &rev}, {Symbol: "AAA", Date: "2026-11-01", Quarter: 3, Year: 2026}}
	b := append([]EarningsItem(nil), a...)
	if reportedEarningsFingerprint(a) != reportedEarningsFingerprint(b) {
		t.Fatal("same reported earnings fingerprint differs")
	}
	eps2 := 1.30
	b[0].EPSActual = &eps2
	if reportedEarningsFingerprint(a) == reportedEarningsFingerprint(b) {
		t.Fatal("changed reported EPS not detected")
	}
	f1 := []FilingItem{{Symbol: "AAA", Form: "10-Q", FiledAt: "2026-08-07", ReportDate: "2026-06-30"}, {Symbol: "AAA", Form: "8-K", FiledAt: "2026-08-08"}}
	f2 := append([]FilingItem(nil), f1...)
	f2[1].FiledAt = "2026-08-09"
	if financialFilingFingerprint(f1) != financialFilingFingerprint(f2) {
		t.Fatal("non-financial 8-K should not invalidate fundamentals fingerprint")
	}
	f2[0].FiledAt = "2026-08-08"
	if financialFilingFingerprint(f1) == financialFilingFingerprint(f2) {
		t.Fatal("new 10-Q should invalidate fundamentals fingerprint")
	}
}

func TestV1330IntegrityAndPreMarketDemoActionsAreNonDestructive(t *testing.T) {
	app := newTestApplication(t)
	for i := range app.state.Watchlists {
		if app.state.Watchlists[i].ID == "swing" {
			app.state.Watchlists[i].Symbols = []string{"AAA", "BBB"}
		}
	}
	beforeLists := clone(app.state.Watchlists)
	beforeSettings := clone(app.state.Settings)
	app.engine.mu.Lock()
	app.engine.status = "running"
	app.engine.mode = "demo"
	app.engine.subscribedSymbols = map[string]bool{"AAA": true}
	app.engine.quotes["AAA"] = Quote{Symbol: "AAA", Price: 123}
	app.engine.mu.Unlock()
	beforeSubs := clone(app.engine.subscribedSymbols)
	rr := httptest.NewRecorder()
	app.handlePreMarketPrep(rr, httptest.NewRequest(http.MethodPost, "/api/cache/pre-market-prep", strings.NewReader(`{}`)))
	if rr.Code != http.StatusOK {
		t.Fatalf("prep status=%d body=%s", rr.Code, rr.Body.String())
	}
	rr = httptest.NewRecorder()
	app.handleIntegrityCheck(rr, httptest.NewRequest(http.MethodPost, "/api/cache/integrity", strings.NewReader(`{}`)))
	if rr.Code != http.StatusOK {
		t.Fatalf("integrity status=%d body=%s", rr.Code, rr.Body.String())
	}
	if !reflect.DeepEqual(beforeLists, app.state.Watchlists) || !reflect.DeepEqual(beforeSettings, app.state.Settings) {
		t.Fatal("maintenance changed user profile")
	}
	app.engine.mu.RLock()
	defer app.engine.mu.RUnlock()
	if !reflect.DeepEqual(beforeSubs, app.engine.subscribedSymbols) {
		t.Fatal("maintenance changed subscriptions")
	}
	if app.engine.quotes["AAA"].Price != 123 {
		t.Fatal("maintenance changed canonical quote")
	}
	if app.engine.lastUpdated["pre-market-prep"] == 0 || app.engine.lastUpdated["weekly-integrity"] == 0 {
		t.Fatal("maintenance timestamps missing")
	}
}

func TestV1330FriendlySECNames(t *testing.T) {
	checks := map[string]string{"10-Q": "Quarterly Report", "10-K": "Annual Report", "8-K": "Material Company Update", "4": "Insider Transaction", "13F-HR": "Institutional Holdings Report", "13D": "Major Shareholder Disclosure", "S-1": "Registration / Share Offering", "DEF 14A": "Proxy Statement"}
	for form, want := range checks {
		got, _ := secFormMeaning(form)
		if got != want {
			t.Fatalf("%s=%q want %q", form, got, want)
		}
	}
}

func TestV1330DemoCacheClearImmediatelyRehydratesHistoricalBars(t *testing.T) {
	app := newTestApplication(t)
	app.engine.mu.Lock()
	app.engine.status = "running"
	app.engine.mode = "demo"
	app.engine.bars["SPY"] = map[string][]Bar{"daily": {{T: 1, C: 500}}}
	app.engine.mu.Unlock()

	req := httptest.NewRequest("POST", "/api/cache/clear", strings.NewReader(`{}`))
	rr := httptest.NewRecorder()
	app.handleCacheClear(rr, req)
	if rr.Code != 200 {
		t.Fatalf("code=%d body=%s", rr.Code, rr.Body.String())
	}
	snap := app.engine.Snapshot()
	if len(snap.Bars["SPY"]["daily"]) == 0 {
		t.Fatalf("demo required daily bars did not rehydrate immediately after cache clear: %+v", snap.Bars["SPY"])
	}
	dayHydrated := false
	for _, symbol := range []string{"NVDA", "AAPL", "TSLA"} {
		if len(snap.Bars[symbol]["intraday"]) > 0 {
			dayHydrated = true
			break
		}
	}
	if !dayHydrated {
		t.Fatal("demo Day intraday bars did not rehydrate immediately after cache clear")
	}
	if app.state.LastCacheCleared <= 0 {
		t.Fatal("cache clear timestamp missing after demo hydration")
	}
}

func TestV1330ProfileExportPreservesAllNonSecretTradingSettings(t *testing.T) {
	app := newTestApplication(t)
	app.mu.Lock()
	app.state.Settings.SignalProfile = "aggressive"
	app.state.Settings.MarketContext = 23
	app.state.Settings.EarningsPenalty = 14
	app.state.Settings.DayEnabled = true
	app.state.Settings.SwingEnabled = false
	app.state.Settings.LongEnabled = true
	app.state.Settings.OvernightDataMode = "indicative"
	app.state.Settings.AIProvider = "gemini"
	app.state.Settings.AutoStart = true
	app.state.Settings.DataMode = "live"
	app.mu.Unlock()

	rr := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/profile/export", nil)
	app.handleExport(rr, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("export code=%d body=%s", rr.Code, rr.Body.String())
	}
	var profile struct {
		Settings Settings `json:"settings"`
	}
	if err := json.Unmarshal(rr.Body.Bytes(), &profile); err != nil {
		t.Fatal(err)
	}
	got := profile.Settings
	if got.SignalProfile != "aggressive" || got.MarketContext != 23 || got.EarningsPenalty != 14 || !got.DayEnabled || got.SwingEnabled || !got.LongEnabled || got.OvernightDataMode != "indicative" || got.AIProvider != "gemini" {
		t.Fatalf("non-secret settings did not round-trip in export: %+v", got)
	}
	if got.DataMode != "demo" || got.AutoStart {
		t.Fatalf("export safety settings incorrect: dataMode=%s autoStart=%v", got.DataMode, got.AutoStart)
	}
}

func TestV15ReservedAlpacaSlotsAreUsedByUrgentPromotion(t *testing.T) {
	app := newTestApplication(t)
	var day, swing []string
	for i := 0; i < 42; i++ {
		day = append(day, fmt.Sprintf("D%03d", i))
	}
	for i := 0; i < 20; i++ {
		swing = append(swing, fmt.Sprintf("S%03d", i))
	}
	for i := range app.state.Watchlists {
		switch app.state.Watchlists[i].ID {
		case "day":
			app.state.Watchlists[i].Symbols = day
		case "swing":
			app.state.Watchlists[i].Symbols = swing
		}
	}
	app.state.UI.SelectedTicker = "URGENT"
	alloc := app.engine.multiFeedAllocation()
	if len(alloc.Alpaca) != alpacaActiveTarget {
		t.Fatalf("selected symbol is now Tier 0 normal capacity and should not consume reserve: got %d want %d", len(alloc.Alpaca), alpacaActiveTarget)
	}
	if !containsLiveSymbol(alloc.Alpaca, "URGENT") {
		t.Fatalf("selected Tier 0 symbol did not receive preferred Alpaca capacity: %+v", alloc.Alpaca)
	}
}

func TestV1433DecisionQueueHintUsesReserveCapacity(t *testing.T) {
	app := newTestApplication(t)
	var day, swing []string
	for i := 0; i < 42; i++ {
		day = append(day, fmt.Sprintf("D%03d", i))
	}
	for i := 0; i < 20; i++ {
		swing = append(swing, fmt.Sprintf("S%03d", i))
	}
	for i := range app.state.Watchlists {
		switch app.state.Watchlists[i].ID {
		case "day":
			app.state.Watchlists[i].Symbols = day
		case "swing":
			app.state.Watchlists[i].Symbols = swing
		}
	}
	app.state.UI.SelectedTicker = ""
	app.engine.mu.Lock()
	app.engine.livePriorityHints["QUEUEX"] = time.Now().UnixMilli()
	app.engine.mu.Unlock()
	alloc := app.engine.multiFeedAllocation()
	if !containsLiveSymbol(alloc.Alpaca, "QUEUEX") {
		t.Fatalf("Decision Queue hint did not promote into Alpaca reserve: %+v", alloc.Alpaca)
	}
	if !alloc.Urgent["QUEUEX"] {
		t.Fatalf("Decision Queue hint not marked urgent: %+v", alloc.Urgent)
	}
}

func TestV1433AlpacaIEXBecomesLiveFailoverWhenFinnhubDisconnected(t *testing.T) {
	app := newTestApplication(t)
	var day []string
	for i := 0; i < 40; i++ {
		day = append(day, fmt.Sprintf("D%03d", i))
	}
	for i := range app.state.Watchlists {
		if app.state.Watchlists[i].ID == "day" {
			app.state.Watchlists[i].Symbols = day
		}
	}
	app.engine.mu.Lock()
	app.engine.webSocketConnected = false
	app.engine.mu.Unlock()
	got := app.engine.alpacaIEXSymbols()
	if len(got) != alpacaPlanMaxSymbols {
		t.Fatalf("Finnhub failover should use Alpaca up to plan ceiling: got %d want %d", len(got), alpacaPlanMaxSymbols)
	}
	for _, pinned := range []string{"GLD", "SLV", "USO"} {
		if !containsLiveSymbol(got, pinned) {
			t.Fatalf("Alpaca failover should protect pinned tradable %s: %+v", pinned, got)
		}
	}
}

func TestV1433CoverageStateRequiresConfirmedSubscription(t *testing.T) {
	alloc := liveAllocation{Finnhub: []string{"AAA"}, Priority: map[string]int{"AAA": 1}}
	now := time.Now().UnixMilli()
	quotes := map[string]Quote{"AAA": {Symbol: "AAA", Price: 100, Source: "finnhub-websocket", ProviderTimestamp: now, UpdatedAt: now}}
	states := buildLiveCoverageStatesFrom(alloc, quotes, false, false, map[string]bool{"AAA": true}, nil)
	if states["AAA"].State == "FINNHUB LIVE" {
		t.Fatalf("coverage must not claim FINNHUB LIVE while socket is disconnected: %+v", states["AAA"])
	}
	quotes["AAA"] = Quote{Symbol: "AAA", Price: 100, Source: "alpaca-iex-websocket-trade", ProviderTimestamp: now, UpdatedAt: now}
	states = buildLiveCoverageStatesFrom(alloc, quotes, false, true, nil, map[string]bool{"AAA": true})
	if states["AAA"].State != "ALPACA IEX LIVE" {
		t.Fatalf("confirmed Alpaca observation should be canonical preferred live state: %+v", states["AAA"])
	}
}

func TestV1433LiveCoverageCountDeduplicatesFailoverSubscriptions(t *testing.T) {
	got := uniqueLiveSubscriptionCount(map[string]bool{"AAA": true, "BBB": true}, map[string]bool{"AAA": true, "CCC": true})
	if got != 3 {
		t.Fatalf("unique live subscription count=%d want=3", got)
	}
}
