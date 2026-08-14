package main

import (
	"math"
	"testing"
	"time"
)

func radarTestSnapshot() alpacaLiveSnapshot {
	var s alpacaLiveSnapshot
	s.LatestTrade.Price = 110
	s.LatestQuote.Bid = 109.9
	s.LatestQuote.Ask = 110.1
	s.DailyBar.Open = 102
	s.DailyBar.High = 112
	s.DailyBar.Low = 100
	s.DailyBar.Close = 110
	s.DailyBar.Volume = 700_000
	s.PrevDailyBar.Open = 100
	s.PrevDailyBar.High = 104
	s.PrevDailyBar.Low = 98
	s.PrevDailyBar.Close = 100
	s.PrevDailyBar.Volume = 1_000_000
	return s
}

func TestV1610RegularSessionRelativeVolumeIsTimeNormalized(t *testing.T) {
	loc := easternLocation()
	now := time.Date(2026, 8, 12, 12, 45, 0, 0, loc) // exactly half of 9:30-16:00
	s := radarTestSnapshot()
	rv := sessionRelativeVolumeFromSnapshot(s, "regular", now)
	if math.Abs(rv-1.4) > .02 {
		t.Fatalf("session RVOL = %.3f, want about 1.40", rv)
	}
	legacy := sessionRelativeVolumeFromSnapshot(s, "after-hours", now)
	if math.Abs(legacy-.7) > .001 {
		t.Fatalf("non-regular RVOL = %.3f, want .70 prior-session comparison", legacy)
	}
}

func TestV1610OpportunityMetricsUseVolumeVolatilityLiquidityPersistenceAndCatalyst(t *testing.T) {
	loc := easternLocation()
	now := time.Date(2026, 8, 12, 12, 0, 0, 0, loc)
	s := radarTestSnapshot()
	base := scannerScoreFromSnapshot("NVDA", "day", s)
	x := enrichOpportunityMetrics(base, s, "regular", now, true, true)
	if x.SessionRelativeVolume <= 1 || x.RangeExpansion <= 1 {
		t.Fatalf("expected normalized RVOL/range expansion: %+v", x)
	}
	if x.UnusualVolumeScore <= 0 || x.VolatilityScore <= 0 || x.OpportunityScore <= 0 {
		t.Fatalf("opportunity metrics missing: %+v", x)
	}
	if x.PriceConfirmation != "CONFIRMED" {
		t.Fatalf("price confirmation = %q", x.PriceConfirmation)
	}
	joined := ""
	for _, r := range x.Reasons {
		joined += r + "|"
	}
	for _, want := range []string{"Persistent across radar cycles", "Material catalyst/news context", "session-normalized volume", "range expansion", "price dislocation"} {
		if !containsFold(joined, want) {
			t.Fatalf("missing reason %q in %q", want, joined)
		}
	}
}

func containsFold(s, sub string) bool {
	return len(s) >= len(sub) && (stringContainsFold(s, sub))
}

func stringContainsFold(s, sub string) bool {
	// tiny helper avoids importing strings only for tests already covered elsewhere.
	for i := 0; i+len(sub) <= len(s); i++ {
		match := true
		for j := range sub {
			a, b := s[i+j], sub[j]
			if a >= 'A' && a <= 'Z' {
				a += 'a' - 'A'
			}
			if b >= 'A' && b <= 'Z' {
				b += 'a' - 'A'
			}
			if a != b {
				match = false
				break
			}
		}
		if match {
			return true
		}
	}
	return false
}

func TestV1610PromotionIsBoundedAndShadowCannotMutateProduction(t *testing.T) {
	now := time.Date(2026, 8, 12, 14, 0, 0, 0, easternLocation())
	rows := []ScannerResult{}
	for i := 0; i < 8; i++ {
		rows = append(rows, ScannerResult{Symbol: []string{"AAA", "BBB", "CCC", "DDD", "EEE", "FFF", "GGG", "HHH"}[i], Price: 50, DollarVolume: 50_000_000, SpreadPercent: .1, SessionRelativeVolume: 2, OpportunityScore: 95 - float64(i)})
	}
	promos := selectOpportunityPromotions(rows, nil, now)
	if len(promos) != opportunityMaxPromotions {
		t.Fatalf("promotions=%d want=%d", len(promos), opportunityMaxPromotions)
	}
	shadowOnly := ScannerResult{Symbol: "SHDW", Price: 50, DollarVolume: 50_000_000, SpreadPercent: .1, SessionRelativeVolume: 2, OpportunityScore: 74}
	if got := selectOpportunityPromotions([]ScannerResult{shadowOnly}, nil, now); len(got) != 0 {
		t.Fatalf("shadow threshold mutated production promotion: %+v", got)
	}
	if !opportunityPromotionEligible(shadowOnly, opportunityShadowFloor) {
		t.Fatal("expected candidate to qualify only for shadow observation")
	}
}

func TestV1610PromotionHysteresisAvoidsSubscriptionChurn(t *testing.T) {
	now := time.Date(2026, 8, 12, 14, 0, 0, 0, easternLocation())
	old := []OpportunityPromotion{{Symbol: "NVDA", Score: 84, State: "PROMOTED", PromotedAt: now.Add(-2 * time.Minute).UnixMilli()}}
	row := ScannerResult{Symbol: "NVDA", Price: 110, DollarVolume: 80_000_000, SpreadPercent: .08, SessionRelativeVolume: 1.6, OpportunityScore: 72}
	promos := selectOpportunityPromotions([]ScannerResult{row}, old, now)
	if len(promos) != 1 || promos[0].PromotedAt != old[0].PromotedAt {
		t.Fatalf("promotion hysteresis lost continuity: %+v", promos)
	}
	row.OpportunityScore = 65
	if got := selectOpportunityPromotions([]ScannerResult{row}, old, now); len(got) != 0 {
		t.Fatalf("weak candidate should demote/expire: %+v", got)
	}
}

func TestV1610RadarCadenceIsSessionAndProviderAware(t *testing.T) {
	if d := opportunityRadarCadence("regular", false, false); d != 2*time.Minute {
		t.Fatalf("regular=%s", d)
	}
	if d := opportunityRadarCadence("regular", true, false); d != time.Minute {
		t.Fatalf("hot=%s", d)
	}
	if d := opportunityRadarCadence("pre-market", false, false); d != 3*time.Minute {
		t.Fatalf("premarket=%s", d)
	}
	if d := opportunityRadarCadence("regular", true, true); d != 2*time.Minute {
		t.Fatalf("degraded=%s", d)
	}
	if radarSessionActive("weekend") {
		t.Fatal("weekend radar should pause")
	}
	if !radarSessionActive("overnight") {
		t.Fatal("provider-aware overnight radar should be eligible")
	}
}

func TestV1610AdaptiveDataPolicyTightensOnlyForHotSymbols(t *testing.T) {
	now := time.Date(2026, 8, 12, 13, 0, 0, 0, easternLocation())
	cold := buildAdaptiveDataPolicyState(ScannerState{}, map[string]string{"alpaca-live": "healthy"}, now)
	hot := buildAdaptiveDataPolicyState(ScannerState{Radar: OpportunityRadarState{Promotions: []OpportunityPromotion{{Symbol: "NVDA"}}}}, map[string]string{"alpaca-live": "healthy"}, now)
	if hot.IntradayHistoryCadence >= cold.IntradayHistoryCadence || hot.CachePersistCadence >= cold.CachePersistCadence {
		t.Fatalf("hot policy did not tighten: cold=%+v hot=%+v", cold, hot)
	}
	if len(hot.HotSymbols) != 1 || hot.HotSymbols[0] != "NVDA" {
		t.Fatalf("hot symbols=%v", hot.HotSymbols)
	}
}

func TestV1610ShadowControlIsReadOnly(t *testing.T) {
	s := ScannerState{Radar: OpportunityRadarState{Candidates: []ScannerResult{{Symbol: "AAA", Price: 20, DollarVolume: 30_000_000, SpreadPercent: .1, SessionRelativeVolume: 2, OpportunityScore: 74}}}}
	shadow := buildShadowControlState(s, time.Now())
	if shadow.PromotionPath != "SHADOW → VALIDATED → APPROVED → PRODUCTION" {
		t.Fatalf("path=%q", shadow.PromotionPath)
	}
	if len(shadow.Experiments) < 2 {
		t.Fatalf("experiments=%+v", shadow.Experiments)
	}
	for _, x := range shadow.Experiments {
		if x.CanMutateProduction {
			t.Fatalf("shadow experiment can mutate production: %+v", x)
		}
	}
}

func TestV1610OpportunityNotificationIsStableMaterialTransition(t *testing.T) {
	now := time.Now()
	p := OpportunityPromotion{Symbol: "NVDA", Score: 88, State: "PROMOTED", PromotedAt: now.Add(-time.Minute).UnixMilli(), ExpiresAt: now.Add(3 * time.Minute).UnixMilli(), Reasons: []string{"2.0x session-normalized volume"}}
	a := opportunityRadarNotifications(ScannerState{Radar: OpportunityRadarState{Promotions: []OpportunityPromotion{p}}}, now)
	b := opportunityRadarNotifications(ScannerState{Radar: OpportunityRadarState{Promotions: []OpportunityPromotion{p}}}, now.Add(30*time.Second))
	if len(a) != 1 || len(b) != 1 || a[0].ID != b[0].ID {
		t.Fatalf("notification was not stable: %+v %+v", a, b)
	}
}

func TestV1610MarketCriticalAndPinnedPriorityRemainProtected(t *testing.T) {
	st := defaultState()
	rows, priority := baselineLiveCandidatesFrom(st)
	if len(rows) < 5 || rows[0] != "SPY" || rows[1] != "QQQ" {
		t.Fatalf("market critical lead order=%v", rows[:minInt(5, len(rows))])
	}
	for _, sym := range []string{"SPY", "QQQ"} {
		if priority[sym] != 0 {
			t.Fatalf("%s priority=%d", sym, priority[sym])
		}
	}
	for _, sym := range []string{"GLD", "SLV", "USO"} {
		if priority[sym] != 1 {
			t.Fatalf("%s priority=%d", sym, priority[sym])
		}
	}
}
