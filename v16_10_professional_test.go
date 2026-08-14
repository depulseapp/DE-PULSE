package main

import (
	"os"
	"strings"
	"testing"
	"time"
)

func TestV1610ProfessionalIlliquidVolumeSpikeCannotPromote(t *testing.T) {
	rows := []ScannerResult{
		{Symbol: "PENNY", Price: 1.25, DollarVolume: 30_000_000, SpreadPercent: .2, SessionRelativeVolume: 9, OpportunityScore: 99},
		{Symbol: "WIDE", Price: 40, DollarVolume: 80_000_000, SpreadPercent: 1.25, SessionRelativeVolume: 4, OpportunityScore: 99},
		{Symbol: "THIN", Price: 40, DollarVolume: 1_000_000, SpreadPercent: .1, SessionRelativeVolume: 8, OpportunityScore: 99},
	}
	if got := selectOpportunityPromotions(rows, nil, time.Now()); len(got) != 0 {
		t.Fatalf("illiquid candidates promoted: %+v", got)
	}
}

func TestV1610ProfessionalRadarPromotionsCannotStarveActiveWatchlist(t *testing.T) {
	st := defaultState()
	// A lower-popularity active user symbol must remain in the normal pool even when
	// five scanner symbols consume every reserve slot.
	st.Settings.DayEnabled = true
	st.Settings.DayWatchlistID = "day"
	for i := range st.Watchlists {
		if st.Watchlists[i].ID == "day" {
			st.Watchlists[i].Symbols = []string{"OBSCURE"}
		}
	}
	hints := map[string]int64{}
	now := time.Now()
	for _, s := range []string{"HOT1", "HOT2", "HOT3", "HOT4", "HOT5"} {
		hints[s] = now.UnixMilli()
	}
	alloc := multiFeedAllocationWithHints(st, nil, nil, hints, now)
	found := false
	for _, s := range append(append([]string{}, alloc.Alpaca...), alloc.Finnhub...) {
		if s == "OBSCURE" {
			found = true
		}
	}
	if !found {
		t.Fatalf("active watchlist symbol starved by radar promotions: %+v", alloc)
	}
	for _, s := range []string{"HOT1", "HOT2", "HOT3", "HOT4", "HOT5"} {
		if !alloc.Urgent[s] {
			t.Fatalf("promotion %s not treated as urgent", s)
		}
	}
}

func TestV1610ProfessionalClosedMarketRadarDoesNotClaimLiveOpportunityState(t *testing.T) {
	for _, session := range []string{"closed", "weekend"} {
		if radarSessionActive(session) {
			t.Fatalf("%s unexpectedly active", session)
		}
		if d := opportunityRadarCadence(session, false, false); d < 10*time.Minute {
			t.Fatalf("%s cadence too aggressive: %s", session, d)
		}
	}
}

func TestV1610ProfessionalShadowThresholdNeverChangesPromotionSelection(t *testing.T) {
	now := time.Now()
	row := ScannerResult{Symbol: "AAA", Price: 50, DollarVolume: 50_000_000, SpreadPercent: .1, SessionRelativeVolume: 2, OpportunityScore: 75}
	if !opportunityPromotionEligible(row, opportunityShadowFloor) {
		t.Fatal("fixture should be shadow-qualified")
	}
	if opportunityPromotionEligible(row, opportunityProductionFloor) {
		t.Fatal("fixture should not be production-qualified")
	}
	if got := selectOpportunityPromotions([]ScannerResult{row}, nil, now); len(got) != 0 {
		t.Fatalf("shadow leaked into production: %+v", got)
	}
}

func TestV1610ProfessionalRadarDoesNotChangeProtectedDeskScoreFormula(t *testing.T) {
	b, err := os.ReadFile("opportunity_radar.go")
	if err != nil {
		t.Fatal(err)
	}
	s := string(b)
	for _, bad := range []string{"SetupScore =", "Action = \"BUY\"", "paper trading", "order execution"} {
		if strings.Contains(s, bad) {
			t.Fatalf("radar source crossed protected decision/execution boundary: %s", bad)
		}
	}
}
