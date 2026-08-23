package main

import (
	"os"
	"strings"
	"testing"
	"time"
)

func v1882CurrentQuote(symbol string, price float64, now time.Time) Quote {
	return Quote{
		Symbol:            symbol,
		Price:             price,
		UpdatedAt:         now.UnixMilli(),
		ProviderTimestamp: now.UnixMilli(),
		DataState:         "real-time",
		FeedType:          "live",
		Source:            "v18.8.2-test",
	}
}

func v1882FreshnessRow(t *testing.T, rows []FreshnessDiagnostic, dataset string) FreshnessDiagnostic {
	t.Helper()
	for _, row := range rows {
		if row.Dataset == dataset {
			return row
		}
	}
	t.Fatalf("freshness row %q not found", dataset)
	return FreshnessDiagnostic{}
}

func TestV1882MarketIntelligenceBreadthIsCanonicalAllocatorDemand(t *testing.T) {
	if got, want := len(broadBreadthUniverse), 15; got != want {
		t.Fatalf("canonical breadth universe size = %d, want %d", got, want)
	}
	if len(marketIntelligenceBreadthUniverse) != len(broadBreadthUniverse) {
		t.Fatalf("Market Intelligence breadth size drifted: got %d, canonical %d", len(marketIntelligenceBreadthUniverse), len(broadBreadthUniverse))
	}
	for i := range broadBreadthUniverse {
		if marketIntelligenceBreadthUniverse[i] != broadBreadthUniverse[i] {
			t.Fatalf("Market Intelligence breadth drift at %d: got %s, canonical %s", i, marketIntelligenceBreadthUniverse[i], broadBreadthUniverse[i])
		}
	}

	alloc := multiFeedAllocationWithHints(AppState{}, nil, nil, nil, time.Now())
	admitted := map[string]bool{}
	for _, symbol := range alloc.Alpaca {
		admitted[symbol] = true
	}
	for _, symbol := range alloc.Finnhub {
		admitted[symbol] = true
	}
	for _, symbol := range alloc.Snapshot {
		admitted[symbol] = true
	}
	for _, symbol := range broadBreadthUniverse {
		if symbol == "VIX" {
			t.Fatal("VIX must remain outside the equity breadth allocator")
		}
		if !admitted[symbol] {
			t.Fatalf("canonical Market Intelligence breadth symbol %s is not allocator-owned live or snapshot demand", symbol)
		}
	}
	for _, symbol := range []string{"SPY", "QQQ"} {
		if !containsLiveSymbol(alloc.Alpaca, symbol) && !containsLiveSymbol(alloc.Finnhub, symbol) {
			t.Fatalf("protected Market Intelligence benchmark %s is not in a live allocation pool", symbol)
		}
	}
}

func TestV1882SnapshotWiresCanonicalBreadthIntoQuoteFreshness(t *testing.T) {
	source, err := os.ReadFile("engine_core.go")
	if err != nil {
		t.Fatalf("read engine_core.go: %v", err)
	}
	contract := "quoteFreshnessSymbols = uniqueSymbols(append(quoteFreshnessSymbols, broadBreadthUniverse...))"
	if !strings.Contains(string(source), contract) {
		t.Fatalf("Engine.Snapshot must include canonical Market Intelligence breadth in the sole quote freshness scope")
	}
}

func TestV1882CanonicalFreshnessDetectsAndClearsMissingBreadthQuote(t *testing.T) {
	now := time.Now()
	quotes := map[string]Quote{}
	for i, symbol := range broadBreadthUniverse {
		quotes[symbol] = v1882CurrentQuote(symbol, 100+float64(i), now)
	}
	e := &Engine{bars: map[string]map[string][]Bar{}}

	rows, _ := e.buildFreshnessDiagnostics(quotes, map[string]int64{}, map[string]string{}, broadBreadthUniverse, nil)
	current := v1882FreshnessRow(t, rows, "Quotes")
	if current.State == "STALE" || current.State == "ERROR" || current.State == "UNAVAILABLE" {
		t.Fatalf("complete current breadth unexpectedly unhealthy: state=%s reason=%s", current.State, current.Reason)
	}

	delete(quotes, "XLU")
	rows, _ = e.buildFreshnessDiagnostics(quotes, map[string]int64{}, map[string]string{}, broadBreadthUniverse, nil)
	missing := v1882FreshnessRow(t, rows, "Quotes")
	if missing.State != "STALE" {
		t.Fatalf("missing admitted breadth quote state = %s, want STALE", missing.State)
	}
	if !strings.Contains(missing.Reason, "1 of 15") || !strings.Contains(strings.ToLower(missing.Reason), "targeted recovery") {
		t.Fatalf("missing breadth quote reason does not expose bounded recovery truth: %q", missing.Reason)
	}

	quotes["XLU"] = v1882CurrentQuote("XLU", 115, now)
	rows, _ = e.buildFreshnessDiagnostics(quotes, map[string]int64{}, map[string]string{}, broadBreadthUniverse, nil)
	recovered := v1882FreshnessRow(t, rows, "Quotes")
	if recovered.State == "STALE" || recovered.State == "ERROR" || recovered.State == "UNAVAILABLE" {
		t.Fatalf("restored breadth quote did not recover freshness: state=%s reason=%s", recovered.State, recovered.Reason)
	}
}

func TestV1882TradeabilityRequiresCurrentSPYQQQVIX(t *testing.T) {
	now := time.Now()
	base := map[string]Quote{
		"SPY": v1882CurrentQuote("SPY", 650, now),
		"QQQ": v1882CurrentQuote("QQQ", 590, now),
		"VIX": v1882CurrentQuote("VIX", 17, now),
	}

	normal := marketTradeabilityWithContext(base, MarketBreadthState{State: "CURRENT"}, GlobalMarketContext{}, MarketTradeabilityContext{}, now)
	if normal.State == "DATA DEGRADED" || normal.Score <= 0 {
		t.Fatalf("current SPY/QQQ/VIX should permit evaluated tradeability, got state=%s score=%d blockers=%v", normal.State, normal.Score, normal.Blockers)
	}

	for _, missing := range []string{"SPY", "QQQ", "VIX"} {
		quotes := map[string]Quote{}
		for symbol, quote := range base {
			quotes[symbol] = quote
		}
		delete(quotes, missing)
		got := marketTradeabilityWithContext(quotes, MarketBreadthState{State: "CURRENT"}, GlobalMarketContext{}, MarketTradeabilityContext{}, now)
		if got.State != "DATA DEGRADED" || got.Score != 0 {
			t.Fatalf("missing %s: state=%s score=%d, want DATA DEGRADED/0", missing, got.State, got.Score)
		}
		if len(got.Blockers) == 0 || !strings.Contains(got.Blockers[0], missing) {
			t.Fatalf("missing %s is not named in tradeability blocker: %v", missing, got.Blockers)
		}
	}
}

func TestV1882TradeabilityRejectsStaleCoreEvidence(t *testing.T) {
	now := time.Now()
	staleAt := now.Add(-120 * time.Hour)
	quotes := map[string]Quote{
		"SPY": v1882CurrentQuote("SPY", 650, staleAt),
		"QQQ": v1882CurrentQuote("QQQ", 590, now),
		"VIX": v1882CurrentQuote("VIX", 17, now),
	}
	got := marketTradeabilityWithContext(quotes, MarketBreadthState{State: "CURRENT"}, GlobalMarketContext{}, MarketTradeabilityContext{}, now)
	if got.State != "DATA DEGRADED" || got.Score != 0 {
		t.Fatalf("stale SPY evidence must be rejected, got state=%s score=%d", got.State, got.Score)
	}
	if len(got.Blockers) == 0 || !strings.Contains(got.Blockers[0], "SPY") {
		t.Fatalf("stale SPY blocker missing: %v", got.Blockers)
	}
}

func TestV1882ZeroOfFifteenBreadthRemainsUnavailableNotDirectionalZero(t *testing.T) {
	got := marketBreadthStateWithBars(map[string]Quote{}, map[string]map[string][]Bar{}, time.Now())
	if got.Expected != 15 || got.Fresh != 0 {
		t.Fatalf("empty breadth coverage = %d/%d, want 0/15", got.Fresh, got.Expected)
	}
	if got.State != "UNAVAILABLE" {
		t.Fatalf("0/15 breadth state = %s, want UNAVAILABLE", got.State)
	}
	if got.ParticipationPct != nil {
		t.Fatalf("0/15 breadth must not fabricate a directional participation percentage: %v", *got.ParticipationPct)
	}
}
