package main

import (
	"strings"
	"testing"
	"time"
)

func TestV1611LongStructureRequiresWeeklyEvidence(t *testing.T) {
	now := time.Date(2026, 8, 11, 15, 0, 0, 0, time.UTC)
	bars := map[string]map[string][]Bar{"SPY": {"weekly": v161BarsEnding("SPY", 10, 400, 2, now.Add(-7*24*time.Hour), 7*24*time.Hour), "daily": v161BarsEnding("SPY", 80, 400, 2, now.Add(-24*time.Hour), 24*time.Hour)}}
	got := marketStructureFor(bars, "long", now)
	if got.State != "UNAVAILABLE" || !strings.Contains(strings.ToLower(got.Detail), "weekly") || got.Coverage != 10 {
		t.Fatalf("daily substitution escaped: %+v", got)
	}
}
func TestV1611StructureRejectsStaleHistoricalBars(t *testing.T) {
	now := time.Now()
	end := now.Add(-120 * 24 * time.Hour)
	bars := map[string]map[string][]Bar{"SPY": {"daily": v161BarsEnding("SPY", 30, 400, 1, end, 24*time.Hour)}}
	got := marketStructureFor(bars, "swing", now)
	if got.State != "UNAVAILABLE" || !strings.Contains(strings.ToLower(got.Detail), "stale") {
		t.Fatalf("stale structure accepted: %+v", got)
	}
	if got.UpdatedAt != end.Unix()*1000 {
		t.Fatalf("timestamp not evidence-derived: got %d want %d", got.UpdatedAt, end.Unix()*1000)
	}
}
func TestV1611RelativeStrengthRejectsStaleBarsAndUsesEvidenceTimestamp(t *testing.T) {
	now := time.Now()
	stale := now.Add(-120 * 24 * time.Hour)
	bars := map[string]map[string][]Bar{"NVDA": {"daily": v161BarsEnding("NVDA", 30, 100, 1, stale, 24*time.Hour)}, "SPY": {"daily": v161BarsEnding("SPY", 30, 400, .5, stale, 24*time.Hour)}}
	got := relativeStrengthFor(bars, "NVDA", "SPY", "swing", now)
	if got.State != "UNAVAILABLE" || got.UpdatedAt != stale.Unix()*1000 {
		t.Fatalf("stale rs accepted/current-stamped: %+v", got)
	}
	nEnd, sEnd := now.Add(-24*time.Hour), now.Add(-48*time.Hour)
	bars["NVDA"]["daily"] = v161BarsEnding("NVDA", 30, 100, 1, nEnd, 24*time.Hour)
	bars["SPY"]["daily"] = v161BarsEnding("SPY", 30, 400, .5, sEnd, 24*time.Hour)
	fresh := relativeStrengthFor(bars, "NVDA", "SPY", "swing", now)
	if fresh.State == "UNAVAILABLE" {
		t.Fatalf("fresh rs unavailable: %+v", fresh)
	}
	if fresh.UpdatedAt != sEnd.Unix()*1000 {
		t.Fatalf("rs stamped calculation time: %+v want %d", fresh, sEnd.Unix()*1000)
	}
}
func TestV1611LiquidityWithoutValidBidAskIsUnknown(t *testing.T) {
	now := time.Now()
	q := v161Quote("AAPL", 100, 0, now)
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
