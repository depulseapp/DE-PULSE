package main

import (
	"testing"
	"time"
)

func TestV1681USActionableTickerBoundary(t *testing.T) {
	for _, sym := range []string{"SHOP.TO", "RY.TO", "BHP.AX", "7203.T", "SHEL.L"} {
		if _, ok := parseUserTicker(sym); ok {
			t.Fatalf("foreign listing %s crossed US actionable boundary", sym)
		}
	}
	for _, sym := range []string{"AAPL", "SHOP", "TSM", "ASML", "BRK.B", "GLD", "USO"} {
		if got, ok := parseUserTicker(sym); !ok || got != sym {
			t.Fatalf("US-listed/actionable syntax rejected: %s got=%s ok=%v", sym, got, ok)
		}
	}
}

func TestV1681MarketCriticalLivePriority(t *testing.T) {
	st := defaultState()
	st.UI.SelectedTicker = "NVDA"
	rows, pri := baselineLiveCandidatesFrom(st)
	if len(rows) < 3 || rows[0] != "SPY" || rows[1] != "QQQ" || rows[2] != "NVDA" {
		t.Fatalf("Tier 0 order must be SPY/QQQ/selected, got %v", rows[:minInt(len(rows), 6)])
	}
	if pri["SPY"] != 0 || pri["QQQ"] != 0 || pri["NVDA"] != 0 {
		t.Fatalf("market-critical priority drift: %+v", pri)
	}
	for _, sym := range []string{"GLD", "SLV", "USO"} {
		if pri[sym] > 1 {
			t.Fatalf("approved tradable %s lost Tier 1 priority: %+v", sym, pri)
		}
	}
}

func TestV1681GlobalMacroIsContextNotMarketCriticalEventMode(t *testing.T) {
	now := time.Now()
	cn := MacroEvent{ID: "cn-gdp", Region: "CN", Name: "Gross Domestic Product", Impact: "HIGH", StartsAt: now.Add(5 * time.Minute).UnixMilli(), TimeKnown: true}
	if got := macroEventProcessingClass(cn); got != "GLOBAL_CONTEXT" {
		t.Fatalf("global event classification=%s", got)
	}
	if eventModeFor([]MacroEvent{cn}, now, true).Active {
		t.Fatal("global context event incorrectly activated full US market-critical Event Mode")
	}
	us := MacroEvent{ID: "us-cpi", Region: "US", Name: "Consumer Price Index", Impact: "HIGH", StartsAt: now.Add(5 * time.Minute).UnixMilli(), TimeKnown: true}
	if !eventModeFor([]MacroEvent{us}, now, true).Active {
		t.Fatal("US high-impact event failed to activate market-critical Event Mode")
	}
}

func TestV1681EconomicCalendarUSVsGlobalContextScope(t *testing.T) {
	now := time.Now()
	rows := buildEconomicCalendar([]MacroEvent{
		{ID: "us", Region: "US", Name: "CPI", Impact: "HIGH", StartsAt: now.Add(time.Hour).UnixMilli(), TimeKnown: true, Source: "BLS"},
		{ID: "jp", Region: "JP", Name: "Monetary Policy Meeting", Impact: "HIGH", StartsAt: now.Add(2 * time.Hour).UnixMilli(), TimeKnown: true, Source: "BOJ"},
	}, nil, now)
	if len(rows) != 2 || rows[0].Scope != "US" || rows[1].Scope != "GLOBAL CONTEXT" {
		t.Fatalf("calendar scopes incorrect: %+v", rows)
	}
}
