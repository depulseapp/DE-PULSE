package main

import (
	"fmt"
	"testing"
	"time"
)

func v1610RadarPerfRows(n int) []ScannerResult {
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

func BenchmarkV1610OpportunityRadar500Candidates(b *testing.B) {
	rows := v1610RadarPerfRows(500)
	now := time.Date(2026, 8, 12, 13, 0, 0, 0, easternLocation())
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = selectOpportunityPromotions(rows, nil, now)
	}
}

func BenchmarkV1610OpportunityMetricEnrichment500Candidates(b *testing.B) {
	snapshots := make([]alpacaLiveSnapshot, 500)
	bases := make([]ScannerResult, 500)
	now := time.Date(2026, 8, 12, 13, 0, 0, 0, easternLocation())
	for i := range snapshots {
		s := radarTestSnapshot()
		s.LatestTrade.Price += float64(i % 20)
		snapshots[i] = s
		bases[i] = scannerScoreFromSnapshot(fmt.Sprintf("E%04d", i), "day", s)
	}
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		for i := range snapshots {
			_ = enrichOpportunityMetrics(bases[i], snapshots[i], "regular", now, i%13 == 0, i%17 == 0)
		}
	}
}
