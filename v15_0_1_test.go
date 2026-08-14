package main

import (
	"context"
	"encoding/json"
	"net/http/httptest"
	"os"
	"strings"
	"sync"
	"testing"
	"time"
)

func membershipBits(m map[string]bool) string {
	var b strings.Builder
	for _, d := range []string{"day", "swing", "long"} {
		if m[d] {
			b.WriteByte('1')
		} else {
			b.WriteByte('0')
		}
	}
	return b.String()
}

func setDeskBits(t *testing.T, app *Application, sym, bits string) {
	t.Helper()
	app.mu.Lock()
	defer app.mu.Unlock()
	ensureDedicatedDeskWatchlists(&app.state, defaultState())
	for i, d := range []string{"day", "swing", "long"} {
		setMembershipLocked(&app.state, d, sym, bits[i] == '1')
	}
}

func callDeskActive(t *testing.T, app *Application, sym, desk string, active bool) map[string]any {
	t.Helper()
	body := `{"symbol":"` + sym + `","desk":"` + desk + `","active":` + map[bool]string{true: "true", false: "false"}[active] + `}`
	rr := httptest.NewRecorder()
	req := httptest.NewRequest("POST", "/api/desk/membership", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	app.handleDeskMembership(rr, req)
	if rr.Code != 200 {
		t.Fatalf("desk %s %v status=%d body=%s", desk, active, rr.Code, rr.Body.String())
	}
	var out map[string]any
	if err := json.Unmarshal(rr.Body.Bytes(), &out); err != nil {
		t.Fatal(err)
	}
	return out
}

func TestV1501DataFreshnessClocksAreIndependent(t *testing.T) {
	app := v15TestApp(t)
	e := app.engine
	e.mu.Lock()
	e.lastUpdated["news"] = 111
	e.mu.Unlock()
	e.setHealth("news", "degraded · provider unavailable")
	e.mu.RLock()
	dataStamp := e.lastUpdated["news"]
	healthStamp := e.lastUpdated["health:news"]
	e.mu.RUnlock()
	if dataStamp != 111 {
		t.Fatalf("health update rewrote data receipt timestamp: %d", dataStamp)
	}
	if healthStamp <= 0 {
		t.Fatal("health timestamp not recorded separately")
	}

	e.updateQuote("XYZ", Quote{Price: 10, Source: "test"}, "test-no-provider-clock")
	e.mu.RLock()
	q := e.quotes["XYZ"]
	recv := e.lastUpdated["quotes"]
	e.mu.RUnlock()
	if q.ProviderTimestamp != 0 {
		t.Fatalf("provider timestamp was synthesized from local clock: %d", q.ProviderTimestamp)
	}
	if q.UpdatedAt <= 0 || recv != q.UpdatedAt {
		t.Fatalf("receipt timestamp not independently recorded: q=%d recv=%d", q.UpdatedAt, recv)
	}

	now := time.Now().UnixMilli()
	e.mu.Lock()
	e.lastUpdated["cache"] = now - 5000
	e.lastUpdated["news"] = now - 1000
	e.news = []NewsItem{{Headline: "x", Datetime: (now - 3000) / 1000, Source: "Finnhub"}}
	e.mu.Unlock()
	rows, _ := e.buildFreshnessDiagnostics(map[string]Quote{}, clone(e.lastUpdated), map[string]string{"news": "healthy · Finnhub"})
	for _, r := range rows {
		if r.Dataset == "News" {
			if r.ProviderTimestamp == 0 || r.ReceivedAt == 0 || r.CacheAt == 0 {
				t.Fatalf("missing independent clocks: %+v", r)
			}
			if r.ProviderTimestamp == r.ReceivedAt || r.ReceivedAt == r.CacheAt {
				t.Fatalf("clocks unexpectedly aliased: %+v", r)
			}
			return
		}
	}
	t.Fatal("News freshness row not found")
}

func TestV1501RateLimitCircuitMetadataAndRecovery(t *testing.T) {
	app := v15TestApp(t)
	e := app.engine
	e.recordProviderFailure("Finnhub", context.DeadlineExceeded)
	e.recordProviderFailure("Finnhub", &providerTestError{"HTTP 429 too many requests"})
	if e.providerAllowed("Finnhub") {
		t.Fatal("rate-limited provider should be suppressed")
	}
	e.mu.RLock()
	c := e.providerCircuits[providerKey("Finnhub")]
	e.mu.RUnlock()
	if c.RateLimitedUntil <= time.Now().UnixMilli() {
		t.Fatal("rate-limit cooldown not recorded")
	}
	e.recordProviderSuccess("Finnhub")
	if !e.providerAllowed("Finnhub") {
		t.Fatal("provider did not recover after success")
	}
	// Latency and attempt metadata must be retained independently of circuit state.
	e.recordProviderLatency("Finnhub", time.Now().Add(-7*time.Millisecond))
	e.mu.RLock()
	c = e.providerCircuits[providerKey("Finnhub")]
	e.mu.RUnlock()
	if c.Attempts < 1 || c.LatencyMs < 1 {
		t.Fatalf("latency/attempt metadata missing after recovery: %+v", c)
	}
}

type providerTestError struct{ s string }

func (e *providerTestError) Error() string { return e.s }

func TestV1501ExecutableRouterOrderAndNoDuplicateAttempt(t *testing.T) {
	app := v15TestApp(t)
	app.mu.Lock()
	app.secrets.Finnhub = "x"
	app.state.Settings.SECEmail = "qa@example.com"
	app.mu.Unlock()
	e := app.engine
	order := []string{}
	_, ok := e.executeProviderRoute(context.Background(), "Fundamentals", map[string]providerRouteAttempt{
		"Finnhub":  func(context.Context) bool { order = append(order, "Finnhub"); return false },
		"SEC":      func(context.Context) bool { order = append(order, "SEC"); return false },
		"yfinance": func(context.Context) bool { order = append(order, "yfinance"); return true },
	})
	if !ok || strings.Join(order, ">") != "Finnhub>SEC>yfinance" {
		t.Fatalf("route order=%v ok=%v", order, ok)
	}

	// A successful hop is terminal; later providers must remain dormant.
	calls := map[string]int{}
	active, ok := e.executeProviderRoute(context.Background(), "VIX / Indices", map[string]providerRouteAttempt{
		"yfinance": func(context.Context) bool { calls["yfinance"]++; return true },
		"CBOE":     func(context.Context) bool { calls["CBOE"]++; return true },
	})
	if !ok || active != "yfinance" || calls["yfinance"] != 1 || calls["CBOE"] != 0 {
		t.Fatalf("duplicate/later provider call: active=%s calls=%v", active, calls)
	}
}

func TestV1501ProviderRegistryContainsOperationalMetadata(t *testing.T) {
	app := v15TestApp(t)
	e := app.engine
	app.mu.Lock()
	app.secrets.Finnhub = "x"
	app.secrets.TwelveData = "x"
	app.secrets.FRED = "x"
	app.state.Settings.SECEmail = "qa@example.com"
	app.mu.Unlock()
	e.recordProviderFailure("Finnhub", context.DeadlineExceeded)
	e.recordProviderSuccess("Finnhub")
	e.recordProviderLatency("Finnhub", time.Now().Add(-5*time.Millisecond))
	app.mu.RLock()
	settings := clone(app.state.Settings)
	secrets := clone(app.secrets)
	app.mu.RUnlock()
	e.mu.RLock()
	snap := e.buildProviderRouterSnapshot(settings, secrets, clone(e.quotes), clone(e.lastUpdated))
	e.mu.RUnlock()
	found := false
	for _, r := range snap.Routes {
		for _, h := range r.Route {
			if h.Provider == "Finnhub" {
				found = true
				if h.Priority <= 0 || h.Quota == "" || h.RateLimit == "" || h.CostClass == "" || h.ExpectedDelay == "" || h.Attempts <= 0 {
					t.Fatalf("incomplete hop metadata: %+v", h)
				}
			}
		}
	}
	if !found {
		t.Fatal("Finnhub hop missing")
	}
}

func TestV1501DeskMembershipAllSevenStatesAndTransitions(t *testing.T) {
	states := []string{"100", "010", "001", "110", "101", "011", "111"}
	desks := []string{"day", "swing", "long"}
	for _, bits := range states {
		for i, desk := range desks {
			app := v15TestApp(t)
			sym := "TST"
			setDeskBits(t, app, sym, bits)
			active := bits[i] == '1'
			desired := !active
			out := callDeskActive(t, app, sym, desk, desired)
			app.mu.RLock()
			got := membershipBits(deskMembershipsLocked(&app.state, sym))
			app.mu.RUnlock()
			n := strings.Count(bits, "1")
			want := []byte(bits)
			protected := false
			if active && n == 1 {
				want[i] = '1'
				protected = true
			} else if active {
				want[i] = '0'
			} else {
				want[i] = '1'
			}
			if got != string(want) {
				t.Fatalf("%s click %s active=%v got=%s want=%s out=%v", bits, desk, active, got, string(want), out)
			}
			if protected && out["protected"] != true {
				t.Fatalf("last desk not protected: %s %s %v", bits, desk, out)
			}
		}
	}
}

func TestV1501DeskMembershipExplicitRequestsAreIdempotent(t *testing.T) {
	app := v15TestApp(t)
	sym := "IDEM"
	setDeskBits(t, app, sym, "100")
	callDeskActive(t, app, sym, "swing", true)
	callDeskActive(t, app, sym, "swing", true) // duplicate rapid request must not toggle back
	app.mu.RLock()
	got := membershipBits(deskMembershipsLocked(&app.state, sym))
	app.mu.RUnlock()
	if got != "110" {
		t.Fatalf("duplicate add was not idempotent: %s", got)
	}
	callDeskActive(t, app, sym, "day", false)
	callDeskActive(t, app, sym, "day", false) // duplicate remove must remain removed
	app.mu.RLock()
	got = membershipBits(deskMembershipsLocked(&app.state, sym))
	app.mu.RUnlock()
	if got != "010" {
		t.Fatalf("duplicate remove was not idempotent: %s", got)
	}
}

func TestV1501DeskMembershipPersistsAcrossReload(t *testing.T) {
	dir := t.TempDir()
	st := defaultState()
	ensureDedicatedDeskWatchlists(&st, defaultState())
	app := &Application{state: st, configDir: dir, hub: NewHub(), sessionKey: "test"}
	app.engine = NewEngine(app)
	setDeskBits(t, app, "PERS", "101")
	app.mu.Lock()
	if err := app.saveLocked(); err != nil {
		app.mu.Unlock()
		t.Fatal(err)
	}
	app.mu.Unlock()
	next := &Application{configDir: dir, hub: NewHub(), sessionKey: "next"}
	next.load()
	next.engine = NewEngine(next)
	next.mu.RLock()
	got := membershipBits(deskMembershipsLocked(&next.state, "PERS"))
	next.mu.RUnlock()
	if got != "101" {
		t.Fatalf("membership did not persist: %s", got)
	}
	if _, err := os.Stat(next.statePath()); err != nil {
		t.Fatal(err)
	}
}

func TestV1501LegacyDeskAddRemoveUsesCanonicalRules(t *testing.T) {
	app := v15TestApp(t)
	sym := "LEG"
	setDeskBits(t, app, sym, "100")
	rr := httptest.NewRecorder()
	req := httptest.NewRequest("POST", "/api/watchlist/remove", strings.NewReader(`{"watchlistId":"day","symbol":"LEG"}`))
	req.Header.Set("Content-Type", "application/json")
	app.handleRemoveSymbol(rr, req)
	if rr.Code != 200 {
		t.Fatalf("legacy remove status=%d", rr.Code)
	}
	app.mu.RLock()
	got := membershipBits(deskMembershipsLocked(&app.state, sym))
	app.mu.RUnlock()
	if got != "100" {
		t.Fatalf("legacy endpoint bypassed last-desk protection: %s", got)
	}
	rr = httptest.NewRecorder()
	req = httptest.NewRequest("POST", "/api/watchlist/add", strings.NewReader(`{"watchlistId":"swing","symbol":"LEG"}`))
	req.Header.Set("Content-Type", "application/json")
	app.handleAddSymbol(rr, req)
	app.mu.RLock()
	got = membershipBits(deskMembershipsLocked(&app.state, sym))
	app.mu.RUnlock()
	if got != "110" {
		t.Fatalf("legacy add did not update canonical state: %s", got)
	}
}

func TestV1501DeskMembershipConcurrentAddNoDuplicatesAndPersists(t *testing.T) {
	app := v15TestApp(t)
	sym := "CONC"
	setDeskBits(t, app, sym, "100")
	var wg sync.WaitGroup
	errs := make(chan string, 24)
	for i := 0; i < 24; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			rr := httptest.NewRecorder()
			req := httptest.NewRequest("POST", "/api/desk/membership", strings.NewReader(`{"symbol":"CONC","desk":"swing","active":true}`))
			req.Header.Set("Content-Type", "application/json")
			app.handleDeskMembership(rr, req)
			if rr.Code != 200 {
				errs <- rr.Body.String()
			}
		}()
	}
	wg.Wait()
	close(errs)
	for msg := range errs {
		t.Fatalf("concurrent membership request failed: %s", msg)
	}
	app.mu.RLock()
	got := membershipBits(deskMembershipsLocked(&app.state, sym))
	var swingCount int
	for _, wl := range app.state.Watchlists {
		if wl.ID == "swing" {
			for _, s := range wl.Symbols {
				if normalizeSymbol(s) == sym {
					swingCount++
				}
			}
		}
	}
	app.mu.RUnlock()
	if got != "110" || swingCount != 1 {
		t.Fatalf("concurrent adds produced bad canonical state/duplicates: membership=%s swingCount=%d", got, swingCount)
	}
	// Canonical state must survive persistence after concurrent updates.
	app.mu.Lock()
	if err := app.saveLocked(); err != nil {
		app.mu.Unlock()
		t.Fatal(err)
	}
	app.mu.Unlock()
	next := &Application{configDir: app.configDir, hub: NewHub(), sessionKey: "reload"}
	next.load()
	next.engine = NewEngine(next)
	next.mu.RLock()
	reloaded := membershipBits(deskMembershipsLocked(&next.state, sym))
	next.mu.RUnlock()
	if reloaded != "110" {
		t.Fatalf("concurrent membership state did not persist: %s", reloaded)
	}
}

func TestV1501NoActiveLegacyTwelveFundamentalsRoute(t *testing.T) {
	if got := strings.Join(routeChains()["Fundamentals"], ">"); got != "Finnhub>SEC>yfinance" {
		t.Fatalf("approved Fundamentals route changed: %s", got)
	}
	app := v15TestApp(t)
	app.mu.Lock()
	app.secrets.TwelveData = "key"
	app.state.Settings.SECEmail = "qa@example.com"
	app.mu.Unlock()
	app.mu.RLock()
	settings := clone(app.state.Settings)
	secrets := clone(app.secrets)
	app.mu.RUnlock()
	rows := buildProviderCapabilityRegistry(settings, secrets, map[string]string{}, nil, nil)
	for _, row := range rows {
		if row.Provider == "Twelve Data" && strings.Contains(strings.ToLower(row.Capability), "fundamental") {
			t.Fatalf("legacy Twelve Data fundamentals still exposed as an active capability: %+v", row)
		}
	}
}
