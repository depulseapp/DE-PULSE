package main

import (
	"strings"
	"testing"
	"time"
)

// extreme30TestApp is the capability-owned in-memory Application fixture for
// the permanent Extreme-30 regression matrix. It preserves only the generic
// setup formerly supplied by the retired pre-v17 test stack.
func extreme30TestApp(t *testing.T) *Application {
	t.Helper()
	st := defaultState()
	ensureDedicatedDeskWatchlists(&st, defaultState())
	app := &Application{state: st, configDir: t.TempDir(), hub: NewHub(), sessionKey: "test"}
	app.engine = NewEngine(app)
	return app
}

// v15TestApp is a temporary call-site compatibility name inside the retained
// Extreme-30 matrix. Wave 4 removes this alias when the large matrix is
// capability-renamed/refactored without changing its assertions.
func v15TestApp(t *testing.T) *Application {
	return extreme30TestApp(t)
}

// setDeskBits provides deterministic desk-membership setup for the Extreme-30
// state-transition matrix.
func setDeskBits(t *testing.T, app *Application, sym, bits string) {
	t.Helper()
	if len(bits) != 3 {
		t.Fatalf("desk bit fixture %q must contain exactly three bits", bits)
	}
	app.mu.Lock()
	defer app.mu.Unlock()
	ensureDedicatedDeskWatchlists(&app.state, defaultState())
	for i, desk := range []string{"day", "swing", "long"} {
		setMembershipLocked(&app.state, desk, sym, bits[i] == '1')
	}
}

// extreme30BarsEnding is the capability-neutral historical-bar fixture kept by
// the permanent Extreme-30 market-intelligence truth assertions.
func extreme30BarsEnding(count int, start, step float64, end time.Time, interval time.Duration) []Bar {
	out := make([]Bar, 0, count)
	first := end.Add(-time.Duration(count-1) * interval)
	for i := 0; i < count; i++ {
		c := start + float64(i)*step
		out = append(out, Bar{T: first.Add(time.Duration(i) * interval).Unix(), O: c - .2, H: c + .5, L: c - .5, C: c, V: 1_000_000})
	}
	return out
}

func extreme30Quote(symbol string, price, change float64, now time.Time) Quote {
	return Quote{Symbol: symbol, Price: price, ChangePercent: change, Bid: price - .02, Ask: price + .02, ProviderTimestamp: now.UnixMilli(), UpdatedAt: now.UnixMilli(), Source: "test"}
}

func extreme30LongStructureRequiresWeeklyEvidence(t *testing.T) {
	now := time.Date(2026, 8, 11, 15, 0, 0, 0, time.UTC)
	bars := map[string]map[string][]Bar{
		"SPY": {
			"weekly": extreme30BarsEnding(10, 400, 2, now.Add(-7*24*time.Hour), 7*24*time.Hour),
			"daily":  extreme30BarsEnding(80, 400, 2, now.Add(-24*time.Hour), 24*time.Hour),
		},
	}
	got := marketStructureFor(bars, "long", now)
	if got.State != "UNAVAILABLE" || !strings.Contains(strings.ToLower(got.Detail), "weekly") || got.Coverage != 10 {
		t.Fatalf("daily substitution escaped: %+v", got)
	}
}

func extreme30StructureRejectsStaleHistoricalBars(t *testing.T) {
	now := time.Now()
	end := now.Add(-120 * 24 * time.Hour)
	bars := map[string]map[string][]Bar{
		"SPY": {"daily": extreme30BarsEnding(30, 400, 1, end, 24*time.Hour)},
	}
	got := marketStructureFor(bars, "swing", now)
	if got.State != "UNAVAILABLE" || !strings.Contains(strings.ToLower(got.Detail), "stale") {
		t.Fatalf("stale structure accepted: %+v", got)
	}
	if got.UpdatedAt != end.Unix()*1000 {
		t.Fatalf("timestamp not evidence-derived: got %d want %d", got.UpdatedAt, end.Unix()*1000)
	}
}

func extreme30RelativeStrengthRejectsStaleBarsAndUsesEvidenceTimestamp(t *testing.T) {
	now := time.Now()
	stale := now.Add(-120 * 24 * time.Hour)
	bars := map[string]map[string][]Bar{
		"NVDA": {"daily": extreme30BarsEnding(30, 100, 1, stale, 24*time.Hour)},
		"SPY":  {"daily": extreme30BarsEnding(30, 400, .5, stale, 24*time.Hour)},
	}
	got := relativeStrengthFor(bars, "NVDA", "SPY", "swing", now)
	if got.State != "UNAVAILABLE" || got.UpdatedAt != stale.Unix()*1000 {
		t.Fatalf("stale rs accepted/current-stamped: %+v", got)
	}
	nEnd, sEnd := now.Add(-24*time.Hour), now.Add(-48*time.Hour)
	bars["NVDA"]["daily"] = extreme30BarsEnding(30, 100, 1, nEnd, 24*time.Hour)
	bars["SPY"]["daily"] = extreme30BarsEnding(30, 400, .5, sEnd, 24*time.Hour)
	fresh := relativeStrengthFor(bars, "NVDA", "SPY", "swing", now)
	if fresh.State == "UNAVAILABLE" {
		t.Fatalf("fresh rs unavailable: %+v", fresh)
	}
	if fresh.UpdatedAt != sEnd.Unix()*1000 {
		t.Fatalf("rs stamped calculation time: %+v want %d", fresh, sEnd.Unix()*1000)
	}
}

func extreme30LiquidityWithoutValidBidAskIsUnknown(t *testing.T) {
	now := time.Now()
	q := extreme30Quote("AAPL", 100, 0, now)
	q.Bid, q.Ask = 0, 0
	raw := deriveLiquidityStatesWithContext(map[string]Quote{"AAPL": q}, nil, nil, now)
	if raw["AAPL"].State != "UNKNOWN" {
		t.Fatalf("raw missing spread safe: %+v", raw["AAPL"])
	}
	mapped := liquiditySlippageState("AAPL", raw)
	if mapped.State != "UNKNOWN" {
		t.Fatalf("mapped missing spread safe: %+v", mapped)
	}
	q.Bid, q.Ask = 101, 100
	raw = deriveLiquidityStatesWithContext(map[string]Quote{"AAPL": q}, nil, nil, now)
	if raw["AAPL"].State != "UNKNOWN" {
		t.Fatalf("crossed market safe: %+v", raw["AAPL"])
	}
}

// Temporary function-value aliases keep the retained Extreme-30 matrix source
// unchanged while the pre-v17 standalone test files remain deleted. Because
// these are variables rather than Test functions, Go does not rediscover the
// retired version-scoped tests as standalone test identities. Wave 4 removes
// the aliases when the matrix call sites are capability-renamed.
var (
	TestV1611LongStructureRequiresWeeklyEvidence                      = extreme30LongStructureRequiresWeeklyEvidence
	TestV1611StructureRejectsStaleHistoricalBars                      = extreme30StructureRejectsStaleHistoricalBars
	TestV1611RelativeStrengthRejectsStaleBarsAndUsesEvidenceTimestamp = extreme30RelativeStrengthRejectsStaleBarsAndUsesEvidenceTimestamp
	TestV1611LiquidityWithoutValidBidAskIsUnknown                     = extreme30LiquidityWithoutValidBidAskIsUnknown
)
