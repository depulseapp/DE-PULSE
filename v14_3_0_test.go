package main

import (
	"strings"
	"testing"
	"time"
)

func TestV1430PreparationJobsAndMarketOpenWindow(t *testing.T) {
	loc := easternLocation()
	now := time.Date(2026, 8, 10, 8, 0, 0, 0, loc) // Monday
	jobs := initialPreparationJobs(now)
	for _, key := range []string{"pre-market-prep", "market-open-prep", "catalyst-watch"} {
		if _, ok := jobs[key]; !ok {
			t.Fatalf("missing preparation job %q", key)
		}
	}
	if !marketOpenPrepWindow(time.Date(2026, 8, 10, 9, 22, 0, 0, loc)) {
		t.Fatal("9:22 ET should be inside market-open prep window")
	}
	if marketOpenPrepWindow(time.Date(2026, 8, 10, 9, 26, 0, 0, loc)) {
		t.Fatal("9:26 ET should be outside market-open prep window")
	}
}

func TestV1430LiquidityHealth(t *testing.T) {
	now := time.Now()
	fresh := now.UnixMilli()
	states := deriveLiquidityStatesWithContext(map[string]Quote{
		"GOOD":  {Symbol: "GOOD", Price: 100, Bid: 99.99, Ask: 100.01, BidSize: 200, AskSize: 180, ProviderTimestamp: fresh},
		"WIDE":  {Symbol: "WIDE", Price: 10, Bid: 9.90, Ask: 10.10, ProviderTimestamp: fresh},
		"STALE": {Symbol: "STALE", Price: 50, Bid: 49.99, Ask: 50.01, ProviderTimestamp: now.Add(-3 * time.Minute).UnixMilli()},
	}, nil, nil, now)
	if states["GOOD"].State != "HEALTHY" {
		t.Fatalf("GOOD=%s", states["GOOD"].State)
	}
	if states["WIDE"].State != "RISK" {
		t.Fatalf("WIDE=%s", states["WIDE"].State)
	}
	if states["STALE"].State != "RISK" {
		t.Fatalf("STALE=%s", states["STALE"].State)
	}
	if states["GOOD"].BidSize != 200 || states["GOOD"].AskSize != 180 {
		t.Fatal("quote sizes not preserved")
	}
}

func TestV1430DerivedIntelligencePreservesSourceFreshness(t *testing.T) {
	stamp := time.Now().Add(-time.Hour).UnixMilli()
	states := deriveIntelligenceStates(map[string]MacroMetric{
		"DGS10":        {Key: "DGS10", Value: 4.25, Change5D: .30, UpdatedAt: stamp},
		"DGS2":         {Key: "DGS2", Value: 3.95, UpdatedAt: stamp},
		"BAMLH0A0HYM2": {Key: "BAMLH0A0HYM2", Value: 3.2, Change20D: .2, UpdatedAt: stamp},
		"NFCI":         {Key: "NFCI", Value: -.35, UpdatedAt: stamp},
		"DTWEXBGS":     {Key: "DTWEXBGS", Value: 122, Change20D: 2, UpdatedAt: stamp},
		"CPI_INDEX":    {Key: "CPI_INDEX", Value: 2.8, Change5D: -.2, UpdatedAt: stamp},
		"UNEMP":        {Key: "UNEMP", Value: 4.2, Change5D: .1, UpdatedAt: stamp},
		"CRUDE_STOCKS": {Key: "CRUDE_STOCKS", Value: 410, Change5D: -2.5, UpdatedAt: stamp},
	})
	if states["rates"].State != "SHOCK HIGHER" {
		t.Fatalf("rates=%+v", states["rates"])
	}
	if states["rates"].UpdatedAt != stamp {
		t.Fatalf("derived state freshness was rewritten: %d != %d", states["rates"].UpdatedAt, stamp)
	}
	for _, key := range []string{"rates", "credit", "financial-conditions", "dollar", "inflation", "labor", "energy"} {
		if _, ok := states[key]; !ok {
			t.Fatalf("missing derived intelligence %q", key)
		}
	}
}

func TestV1430ProviderCapabilityVocabulary(t *testing.T) {
	rows := buildProviderCapabilityRegistry(Settings{}, Secrets{}, map[string]string{}, map[string]SymbolIntelligence{}, map[string]GlobalDriver{})
	allowed := map[string]bool{"AVAILABLE": true, "PLAN LIMITED": true, "NOT ENTITLED": true, "TEMPORARILY UNAVAILABLE": true}
	if len(rows) < 6 {
		t.Fatalf("unexpected registry size %d", len(rows))
	}
	for _, row := range rows {
		if !allowed[row.Status] {
			t.Fatalf("invalid status %q for %s/%s", row.Status, row.Provider, row.Capability)
		}
	}
}

func TestV1430MaterialCatalystThresholds(t *testing.T) {
	for _, s := range []string{"Company raises guidance", "8-K material update", "secondary offering announced", "FDA regulatory decision"} {
		if !materialText(s) {
			t.Fatalf("expected material: %q", s)
		}
	}
	for _, s := range []string{"routine conference attendance", "weekly marketing update", "website redesign"} {
		if materialText(s) {
			t.Fatalf("unexpected material trigger: %q", s)
		}
	}
}

func TestV1436RatesIntelligenceFallsBackToOfficialTreasury(t *testing.T) {
	stamp := time.Now().Add(-time.Hour).UnixMilli()
	states := deriveIntelligenceStates(map[string]MacroMetric{
		"UST10Y": {Key: "UST10Y", Value: 4.15, Change5D: .11, UpdatedAt: stamp, Source: "U.S. Treasury", Status: "OFFICIAL"},
		"UST2Y":  {Key: "UST2Y", Value: 3.85, UpdatedAt: stamp, Source: "U.S. Treasury", Status: "OFFICIAL"},
	})
	rates, ok := states["rates"]
	if !ok {
		t.Fatal("Treasury core rates must keep Rates State available when FRED is unavailable")
	}
	if rates.State != "RISING" || rates.Source != "U.S. Treasury" || !strings.Contains(rates.Detail, "10Y−2Y 0.30pp") {
		t.Fatalf("unexpected Treasury fallback rates state: %+v", rates)
	}
}

func TestV1436MacroRatesAggregateDoesNotFalseDegradeWhenTreasuryIsHealthy(t *testing.T) {
	now := time.Now().UnixMilli()
	e := &Engine{
		macroMetrics: map[string]MacroMetric{
			"UST10Y": {Key: "UST10Y", Value: 4.15, UpdatedAt: now, Source: "U.S. Treasury", Status: "OFFICIAL"},
			"UST2Y":  {Key: "UST2Y", Value: 3.85, UpdatedAt: now, Source: "U.S. Treasury", Status: "OFFICIAL"},
		},
		lastUpdated: map[string]int64{"treasury": now},
		health: map[string]string{
			"treasury":   "healthy · official Treasury yields",
			"fred-rates": "temporarily unavailable · keeping cached FRED enrichment",
		},
	}
	e.reconcileMacroRatesHealth()
	if got := e.health["macro-rates"]; !strings.HasPrefix(got, "healthy · official Treasury core rates") {
		t.Fatalf("usable Treasury rates must not be falsely marked DEGRADED: %q", got)
	}
	if e.lastUpdated["macro-rates"] != now {
		t.Fatalf("aggregate freshness should use source freshness: %d != %d", e.lastUpdated["macro-rates"], now)
	}

	rows := buildProviderCapabilityRegistry(Settings{}, Secrets{FRED: "configured"}, e.health, map[string]SymbolIntelligence{}, map[string]GlobalDriver{})
	for _, row := range rows {
		if row.Provider == "FRED" {
			if row.Status == "AVAILABLE" {
				t.Fatalf("aggregate Treasury health must not falsely prove FRED availability: %+v", row)
			}
			return
		}
	}
	t.Fatal("FRED provider capability row missing")
}
