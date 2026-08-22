package main

import (
	"fmt"
	"testing"
	"time"
)

// performanceDailyBars is a test-only fixture for bounded product-performance
// evidence. It deliberately avoids restoring any pre-v17 version-named helper.
func performanceDailyBars(count int, start, step float64, now time.Time) []Bar {
	rows := make([]Bar, 0, count)
	base := now.AddDate(0, 0, -count+1).Truncate(24 * time.Hour)
	for i := 0; i < count; i++ {
		c := start + float64(i)*step
		rows = append(rows, Bar{
			T: base.AddDate(0, 0, i).Unix(),
			O: c - step*.25,
			H: c + 1,
			L: c - 1,
			C: c,
			V: 1_000_000,
		})
	}
	return rows
}

func performanceMarketFixture(n int) (AppState, map[string]Quote, map[string]map[string][]Bar) {
	st := defaultState()
	symbols := make([]string, 0, n)
	quotes := make(map[string]Quote, n)
	bars := make(map[string]map[string][]Bar, n)
	now := time.Now()
	for i := 0; i < n; i++ {
		sym := fmt.Sprintf("X%03d", i)
		symbols = append(symbols, sym)
		px := 50.0 + float64(i%200)
		quotes[sym] = Quote{
			Symbol:            sym,
			Price:             px,
			Bid:               px - .02,
			Ask:               px + .02,
			Volume:            1_000_000,
			ChangePercent:     float64((i%11)-5) / 10,
			ProviderTimestamp: now.UnixMilli(),
			UpdatedAt:         now.UnixMilli(),
		}
		bars[sym] = map[string][]Bar{
			"daily": performanceDailyBars(25, px-3, .12, now),
		}
	}
	st.Watchlists = []Watchlist{{ID: "swing", Name: "Swing Watchlist", Symbols: symbols}}
	st.Settings.SwingWatchlistID = "swing"
	st.Settings.DayWatchlistID = ""
	st.Settings.LongWatchlistID = ""
	st.Settings.DiscoveryWatchlistID = ""
	return st, quotes, bars
}

func performanceRadarSnapshot() alpacaLiveSnapshot {
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

func performanceRadarRows(n int) []ScannerResult {
	now := time.Date(2026, 8, 12, 13, 0, 0, 0, easternLocation())
	rows := make([]ScannerResult, 0, n)
	for i := 0; i < n; i++ {
		var s alpacaLiveSnapshot
		px := 20.0 + float64(i%180)
		s.LatestTrade.Price = px
		s.LatestQuote.Bid = px - .03
		s.LatestQuote.Ask = px + .03
		s.DailyBar.Open = px * .97
		s.DailyBar.High = px * 1.04
		s.DailyBar.Low = px * .95
		s.DailyBar.Close = px
		s.DailyBar.Volume = float64(350000 + (i%12)*80000)
		s.PrevDailyBar.Open = px * .96
		s.PrevDailyBar.High = px * .99
		s.PrevDailyBar.Low = px * .93
		s.PrevDailyBar.Close = px * .96
		s.PrevDailyBar.Volume = float64(700000 + (i%10)*50000)
		base := scannerScoreFromSnapshot(fmt.Sprintf("R%04d", i), "day", s)
		rows = append(rows, enrichOpportunityMetrics(base, s, "regular", now, i%13 == 0, i%17 == 0))
	}
	return rows
}

func performancePromotionCandidates(n int) []ScannerResult {
	rows := make([]ScannerResult, 0, n)
	for i := 0; i < n; i++ {
		rows = append(rows, ScannerResult{
			Symbol:                fmt.Sprintf("T%04d", i),
			Price:                 25 + float64(i%50),
			DollarVolume:          10_000_000 + float64(i)*100_000,
			SpreadPercent:         .05 + float64(i%5)*.01,
			SessionRelativeVolume: 1.4 + float64(i%20)/10,
			RangeExpansion:        1.3 + float64(i%10)/10,
			ChangePercent:         3 + float64(i%12)/3,
			OpportunityScore:      80 + float64(i%20),
			PriceConfirmation:     "CONFIRMED",
		})
	}
	return rows
}

func requirePerformanceBound(t *testing.T, name string, result testing.BenchmarkResult, maxNs, maxBytes int64) {
	t.Helper()
	ns := result.NsPerOp()
	bytes := result.AllocedBytesPerOp()
	if ns > maxNs || bytes > maxBytes {
		t.Fatalf("%s regression: ns/op=%d limit=%d B/op=%d limit=%d", name, ns, maxNs, bytes, maxBytes)
	}
}

func TestPerformanceHeatMap500SymbolsBounded(t *testing.T) {
	st, quotes, _ := performanceMarketFixture(500)
	now := time.Now()
	result := testing.Benchmark(func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			_ = v165HeatMapForState(st, quotes, now)
		}
	})
	requirePerformanceBound(t, "HeatMap500Symbols", result, 150_000_000, 12_000_000)
}

func TestPerformanceLiquidity500SymbolsBounded(t *testing.T) {
	_, quotes, bars := performanceMarketFixture(500)
	now := time.Now()
	base := map[string]LiquidityBaseline{}
	result := testing.Benchmark(func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			_ = deriveLiquidityStatesWithContext(quotes, bars, base, now)
		}
	})
	requirePerformanceBound(t, "Liquidity500Symbols", result, 250_000_000, 30_000_000)
}

func TestPerformanceOpportunityRadar500CandidatesBounded(t *testing.T) {
	rows := performanceRadarRows(500)
	now := time.Date(2026, 8, 12, 13, 0, 0, 0, easternLocation())
	result := testing.Benchmark(func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			_ = selectOpportunityPromotions(rows, nil, now)
		}
	})
	requirePerformanceBound(t, "OpportunityRadar500Candidates", result, 50_000_000, 20_000_000)
}

func TestPerformanceOpportunityMetricEnrichment500CandidatesBounded(t *testing.T) {
	snapshots := make([]alpacaLiveSnapshot, 500)
	bases := make([]ScannerResult, 500)
	now := time.Date(2026, 8, 12, 13, 0, 0, 0, easternLocation())
	for i := range snapshots {
		s := performanceRadarSnapshot()
		s.LatestTrade.Price += float64(i % 20)
		snapshots[i] = s
		bases[i] = scannerScoreFromSnapshot(fmt.Sprintf("E%04d", i), "day", s)
	}
	result := testing.Benchmark(func(b *testing.B) {
		b.ReportAllocs()
		for n := 0; n < b.N; n++ {
			for i := range snapshots {
				_ = enrichOpportunityMetrics(bases[i], snapshots[i], "regular", now, i%13 == 0, i%17 == 0)
			}
		}
	})
	requirePerformanceBound(t, "OpportunityMetricEnrichment500Candidates", result, 150_000_000, 40_000_000)
}

func TestPerformanceRadarPromotionStateBoundedAcrossLongRun(t *testing.T) {
	rows := performancePromotionCandidates(500)
	now := time.Date(2026, 8, 12, 15, 0, 0, 0, time.UTC)
	var prev []OpportunityPromotion
	cycles := 10_000
	if raceBuild {
		cycles = 1_000
	}
	for i := 0; i < cycles; i++ {
		prev = selectOpportunityPromotions(rows, prev, now.Add(time.Duration(i)*time.Second))
		if len(prev) > opportunityMaxPromotions {
			t.Fatalf("promotion state grew beyond bound at cycle %d: %d", i, len(prev))
		}
		for _, p := range prev {
			if p.ExpiresAt <= p.LastConfirmedAt {
				t.Fatalf("invalid promotion expiry at cycle %d: %+v", i, p)
			}
		}
	}
	if len(prev) != opportunityMaxPromotions {
		t.Fatalf("expected bounded full promotion set, got %d", len(prev))
	}
}

func TestPerformanceRadarUniverseRotationRemainsBounded(t *testing.T) {
	universe := make([]string, 2000)
	for i := range universe {
		universe[i] = fmt.Sprintf("U%04d", i)
	}
	cursor := 0
	for i := 0; i < 1000; i++ {
		sample, next := radarSampleUniverse(universe, cursor)
		if len(sample) > opportunityRotateCount+len(uniqueSymbols(discoverySeedUniverse)) {
			t.Fatalf("sample grew without bound: %d", len(sample))
		}
		if next < 0 || next >= len(universe) {
			t.Fatalf("cursor out of range: %d", next)
		}
		cursor = next
	}
}

func TestPerformanceAdaptiveCadencesRemainBounded(t *testing.T) {
	cases := []struct {
		session         string
		promo, degraded bool
	}{
		{"regular", false, false},
		{"regular", true, false},
		{"pre-market", false, false},
		{"after-hours", false, true},
		{"overnight", false, true},
		{"closed", false, true},
		{"weekend", false, true},
	}
	for _, c := range cases {
		d := opportunityRadarCadence(c.session, c.promo, c.degraded)
		if d < 30*time.Second || d > 30*time.Minute {
			t.Fatalf("unbounded cadence %s: %s", c.session, d)
		}
	}
}

func TestPerformancePromotionSelection500Bounded(t *testing.T) {
	rows := performancePromotionCandidates(500)
	now := time.Date(2026, 8, 12, 15, 0, 0, 0, time.UTC)
	var prev []OpportunityPromotion
	result := testing.Benchmark(func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			prev = selectOpportunityPromotions(rows, prev, now.Add(time.Duration(i)*time.Second))
		}
	})
	if result.NsPerOp() > 10_000_000 || result.AllocedBytesPerOp() > 5_000_000 || result.AllocsPerOp() > 1000 {
		t.Fatalf("PromotionSelection500 regression: ns/op=%d B/op=%d allocs/op=%d", result.NsPerOp(), result.AllocedBytesPerOp(), result.AllocsPerOp())
	}
}
