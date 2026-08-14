package main

import (
	"fmt"
	"testing"
	"time"
)

func v1611ClosureCandidates(n int) []ScannerResult {
	rows := make([]ScannerResult, 0, n)
	for i := 0; i < n; i++ {
		rows = append(rows, ScannerResult{
			Symbol: fmt.Sprintf("T%04d", i), Price: 25 + float64(i%50),
			DollarVolume:          10_000_000 + float64(i)*100_000,
			SpreadPercent:         .05 + float64(i%5)*.01,
			SessionRelativeVolume: 1.4 + float64(i%20)/10,
			RangeExpansion:        1.3 + float64(i%10)/10,
			ChangePercent:         3 + float64(i%12)/3,
			OpportunityScore:      80 + float64(i%20), PriceConfirmation: "CONFIRMED",
		})
	}
	return rows
}

func TestV1611MajorClosureRadarPromotionStateBoundedAcrossLongRun(t *testing.T) {
	rows := v1611ClosureCandidates(500)
	now := time.Date(2026, 8, 12, 15, 0, 0, 0, time.UTC)
	var prev []OpportunityPromotion
	cycles := 10_000
	if raceBuild {
		// Preserve the full 10,000-cycle longevity acceptance in the normal gate.
		// Under -race, instrumentation makes the O(candidates*cycles) stress loop
		// materially slower; 1,000 cycles still exercises bounded state repeatedly
		// while the dedicated non-race performance/closure gate retains full depth.
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

func TestV1611MajorClosureRadarUniverseRotationRemainsBounded(t *testing.T) {
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

func TestV1611MajorClosureAdaptiveCadencesRemainBounded(t *testing.T) {
	cases := []struct {
		session         string
		promo, degraded bool
	}{
		{"regular", false, false}, {"regular", true, false}, {"pre-market", false, false},
		{"after-hours", false, true}, {"overnight", false, true}, {"closed", false, true}, {"weekend", false, true},
	}
	for _, c := range cases {
		d := opportunityRadarCadence(c.session, c.promo, c.degraded)
		if d < 30*time.Second || d > 30*time.Minute {
			t.Fatalf("unbounded cadence %s: %s", c.session, d)
		}
	}
}

func BenchmarkV1611MajorClosurePromotionSelection500(b *testing.B) {
	rows := v1611ClosureCandidates(500)
	now := time.Date(2026, 8, 12, 15, 0, 0, 0, time.UTC)
	var prev []OpportunityPromotion
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		prev = selectOpportunityPromotions(rows, prev, now.Add(time.Duration(i)*time.Second))
	}
	_ = prev
}
