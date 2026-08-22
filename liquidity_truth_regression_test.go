package main

import (
	"testing"
	"time"
)

func TestV174LiquidityRiskRequiresCurrentMarketEvidence(t *testing.T) {
	now := time.Date(2026, time.August, 12, 10, 0, 0, 0, easternLocation())
	current := LiquidityState{State: "RISK", QuoteAgeMs: 1000, UpdatedAt: now.Add(-time.Second).UnixMilli(), Detail: "wide spread"}
	if !currentLiquidityMarketRisk(current, now) {
		t.Fatal("current measured liquidity risk must remain a market-open risk")
	}
	stale := current
	stale.QuoteAgeMs = int64((3 * time.Minute) / time.Millisecond)
	stale.UpdatedAt = now.Add(-3 * time.Minute).UnixMilli()
	if currentLiquidityMarketRisk(stale, now) {
		t.Fatal("stale quote evidence must be classified as freshness risk, not liquidity risk")
	}
	missing := current
	missing.UpdatedAt = 0
	if currentLiquidityMarketRisk(missing, now) {
		t.Fatal("missing quote timestamp must not become per-symbol liquidity risk")
	}
	future := current
	future.UpdatedAt = now.Add(time.Minute).UnixMilli()
	if currentLiquidityMarketRisk(future, now) {
		t.Fatal("future-skewed quote evidence must not become per-symbol liquidity risk")
	}
}

func TestV174DerivedStaleLiquidityDoesNotMasqueradeAsMarketRisk(t *testing.T) {
	now := time.Date(2026, time.August, 12, 10, 0, 0, 0, easternLocation())
	staleAt := now.Add(-5 * time.Minute).UnixMilli()
	stale := deriveLiquidityStatesWithContext(map[string]Quote{"AAA": {Symbol: "AAA", Price: 100, Bid: 99, Ask: 101, ProviderTimestamp: staleAt, UpdatedAt: staleAt}}, nil, nil, now)["AAA"]
	if stale.State != "RISK" {
		t.Fatalf("expected conservative raw liquidity state for stale quote, got %+v", stale)
	}
	if currentLiquidityMarketRisk(stale, now) {
		t.Fatal("raw conservative stale state must be filtered from Market Open LIQUIDITY RISK flags")
	}

	freshAt := now.Add(-time.Second).UnixMilli()
	freshWide := deriveLiquidityStatesWithContext(map[string]Quote{"AAA": {Symbol: "AAA", Price: 100, Bid: 99, Ask: 101, ProviderTimestamp: freshAt, UpdatedAt: freshAt}}, nil, nil, now)["AAA"]
	if freshWide.State != "RISK" || !currentLiquidityMarketRisk(freshWide, now) {
		t.Fatalf("genuine current wide-spread risk must remain visible: %+v", freshWide)
	}
}
