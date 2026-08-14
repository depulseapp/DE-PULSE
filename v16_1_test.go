package main

import (
	"strings"
	"testing"
	"time"
)

func v161BarsEnding(symbol string, count int, start, step float64, end time.Time, interval time.Duration) []Bar {
	out := make([]Bar, 0, count)
	first := end.Add(-time.Duration(count-1) * interval)
	for i := 0; i < count; i++ {
		c := start + float64(i)*step
		out = append(out, Bar{T: first.Add(time.Duration(i) * interval).Unix(), O: c - .2, H: c + .5, L: c - .5, C: c, V: 1_000_000})
	}
	return out
}
func v161Bars(symbol string, count int, start, step float64) []Bar {
	return v161BarsEnding(symbol, count, start, step, time.Now().Add(-24*time.Hour), 24*time.Hour)
}
func v161Quote(symbol string, price, change float64, now time.Time) Quote {
	return Quote{Symbol: symbol, Price: price, ChangePercent: change, Bid: price - .02, Ask: price + .02, ProviderTimestamp: now.UnixMilli(), UpdatedAt: now.UnixMilli(), Source: "test"}
}

func TestV161MarketStructureUsesExactHorizons(t *testing.T) {
	now := time.Now()
	bars := map[string]map[string][]Bar{"SPY": {
		"intraday": v161BarsEnding("SPY", 20, 500, .2, now.Add(-5*time.Minute), 5*time.Minute),
		"daily":    v161BarsEnding("SPY", 40, 480, .7, now.Add(-24*time.Hour), 24*time.Hour),
		"weekly":   v161BarsEnding("SPY", 40, 400, 2.0, now.Add(-7*24*time.Hour), 7*24*time.Hour),
	}}
	for _, h := range []string{"day", "swing", "long"} {
		if got := marketStructureFor(bars, h, now); got.State == "UNAVAILABLE" {
			t.Fatalf("%s unavailable: %+v", h, got)
		}
	}
}
func TestV161MarketStructureDoesNotUseFutureQuote(t *testing.T) {
	// Structure is derived from canonical bars and has no quote argument; this
	// regression protects the design against future quote substitution.
	now := time.Now()
	bars := map[string]map[string][]Bar{"SPY": {"daily": v161BarsEnding("SPY", 30, 400, 1, now.Add(-24*time.Hour), 24*time.Hour)}}
	got := marketStructureFor(bars, "swing", now)
	if got.State == "UNAVAILABLE" {
		t.Fatalf("expected canonical-bar structure: %+v", got)
	}
	if got.UpdatedAt > now.UnixMilli() {
		t.Fatalf("future update: %+v", got)
	}
}
func TestV161BreadthUsesExplicitDenominatorAndStaleCannotVote(t *testing.T) {
	now := time.Now()
	qs := map[string]Quote{}
	for i, s := range marketIntelligenceBreadthUniverse {
		q := v161Quote(s, 100, float64(i%3)-1, now)
		qs[s] = q
	}
	stale := qs[marketIntelligenceBreadthUniverse[0]]
	stale.ProviderTimestamp = now.Add(-10 * 24 * time.Hour).UnixMilli()
	stale.UpdatedAt = stale.ProviderTimestamp
	stale.ChangePercent = 99
	qs[marketIntelligenceBreadthUniverse[0]] = stale
	b := marketBreadthStateWithBars(qs, nil, now)
	if b.Expected != len(marketIntelligenceBreadthUniverse) || b.Fresh != len(marketIntelligenceBreadthUniverse)-1 {
		t.Fatalf("breadth truth: %+v", b)
	}
	if b.Advancers >= b.Fresh {
		t.Fatalf("stale member appears to vote: %+v", b)
	}
}
func TestV161BreadthDegradedWithholdsParticipation(t *testing.T) {
	now := time.Now()
	qs := map[string]Quote{"SPY": v161Quote("SPY", 500, 2, now), "QQQ": v161Quote("QQQ", 400, 2, now)}
	b := marketBreadthStateWithBars(qs, nil, now)
	if b.State != "DEGRADED" || b.ParticipationPct != nil {
		t.Fatalf("degraded breadth overstated: %+v", b)
	}
}
func TestV161TradeabilityRequiresSPYQQQVIX(t *testing.T) {
	now := time.Now()
	q := map[string]Quote{"SPY": v161Quote("SPY", 500, 0, now), "QQQ": v161Quote("QQQ", 400, 0, now)}
	got := marketTradeability(q, MarketBreadthState{State: "DEGRADED"}, GlobalMarketContext{}, now)
	if got.State != "DATA DEGRADED" || !strings.Contains(strings.Join(got.Blockers, " "), "VIX") {
		t.Fatalf("tradeability=%+v", got)
	}
}
func TestV161RelativeStrengthMissingBenchmarkIsUnavailable(t *testing.T) {
	now := time.Now()
	bars := map[string]map[string][]Bar{"NVDA": {"daily": v161BarsEnding("NVDA", 30, 100, 1, now.Add(-24*time.Hour), 24*time.Hour)}}
	got := relativeStrengthFor(bars, "NVDA", "SPY", "swing", now)
	if got.State != "UNAVAILABLE" {
		t.Fatalf("rs=%+v", got)
	}
}
func TestV161RelativeStrengthExactHorizonNoShortening(t *testing.T) {
	now := time.Now()
	bars := map[string]map[string][]Bar{"NVDA": {"daily": v161BarsEnding("NVDA", 15, 100, 1, now.Add(-24*time.Hour), 24*time.Hour)}, "SPY": {"daily": v161BarsEnding("SPY", 15, 400, .5, now.Add(-24*time.Hour), 24*time.Hour)}}
	got := relativeStrengthFor(bars, "NVDA", "SPY", "swing", now)
	if got.State != "UNAVAILABLE" {
		t.Fatalf("short horizon accepted: %+v", got)
	}
}
func TestV161CanonicalClassificationAndIndustryRegime(t *testing.T) {
	c := classificationForSymbol("nvda")
	if c.SectorBenchmark != "XLK" || c.IndustryBenchmark != "SMH" || c.Industry != "Semiconductors" {
		t.Fatalf("classification=%+v", c)
	}
	now := time.Now()
	bars := map[string]map[string][]Bar{"SMH": {"daily": v161BarsEnding("SMH", 30, 200, 1, now.Add(-24*time.Hour), 24*time.Hour)}}
	r := benchmarkRegime("Industry", c.Industry, c.IndustryBenchmark, bars, now)
	if r.State == "UNAVAILABLE" {
		t.Fatalf("industry regime=%+v", r)
	}
}
func TestV161LiquidityMappingUnknownIsNeverSafe(t *testing.T) {
	got := liquiditySlippageState("AAPL", map[string]LiquidityState{"AAPL": {Symbol: "AAPL", State: "UNKNOWN", Detail: "no spread"}})
	if got.State != "UNKNOWN" {
		t.Fatalf("liquidity=%+v", got)
	}
}
func TestV161SnapshotWiresAllCommittedContext(t *testing.T) {
	now := time.Now()
	bars := map[string]map[string][]Bar{}
	for _, s := range []string{"SPY", "QQQ", "NVDA", "XLK", "SMH"} {
		bars[s] = map[string][]Bar{"daily": v161BarsEnding(s, 40, 100, 1, now.Add(-24*time.Hour), 24*time.Hour), "weekly": v161BarsEnding(s, 40, 100, 1, now.Add(-7*24*time.Hour), 7*24*time.Hour), "intraday": v161BarsEnding(s, 20, 100, .1, now.Add(-5*time.Minute), 5*time.Minute)}
	}
	qs := map[string]Quote{}
	for _, s := range marketIntelligenceBreadthUniverse {
		qs[s] = v161Quote(s, 100, 1, now)
	}
	qs["VIX"] = v161Quote("VIX", 18, -1, now)
	qs["NVDA"] = v161Quote("NVDA", 130, 1, now)
	raw := deriveLiquidityStatesWithContext(qs, nil, nil, now)
	st := AppState{UI: UIState{SelectedTicker: "NVDA"}}
	got := buildMarketIntelligenceSnapshotWithContext(st, qs, bars, raw, GlobalMarketContext{Tone: "CONSTRUCTIVE"}, MarketTradeabilityContext{}, now)
	if got.SelectedSymbol != "NVDA" || got.Classification.Industry != "Semiconductors" || got.Tradeability.State == "DATA DEGRADED" || len(got.RelativeStrength) < 3 {
		t.Fatalf("snapshot=%+v", got)
	}
}
