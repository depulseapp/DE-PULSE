package main

import (
	"testing"
	"time"
)

func TestLatestProviderEvidenceMillisUsesNewestObservation(t *testing.T) {
	older := "2026-08-20T14:00:00Z"
	newer := "2026-08-20T14:05:00Z"
	got := latestProviderEvidenceMillis(older, newer)
	want := providerTimeMillis(newer)
	if got != want {
		t.Fatalf("latest provider evidence = %d, want %d", got, want)
	}
}

func TestScannerUnknownEvidenceDoesNotBecomeFresh(t *testing.T) {
	got := scannerScoreFromSnapshot("AAPL", "day", alpacaLiveSnapshot{})
	if got.UpdatedAt != 0 {
		t.Fatalf("unknown provider evidence became fresh: UpdatedAt=%d", got.UpdatedAt)
	}
}

func TestScannerUsesProviderEvidenceNotRetrievalTime(t *testing.T) {
	var snap alpacaLiveSnapshot
	snap.LatestTrade.Time = "2026-08-20T13:00:00Z"
	retrievedAt := time.Date(2026, 8, 20, 15, 0, 0, 0, time.UTC)
	got := scannerScoreFromSnapshot("AAPL", "day", snap)
	want := providerTimeMillis(snap.LatestTrade.Time)
	if got.UpdatedAt != want {
		t.Fatalf("scanner evidence time = %d, want provider observation %d", got.UpdatedAt, want)
	}
	if got.UpdatedAt == retrievedAt.UnixMilli() {
		t.Fatal("scanner freshness incorrectly used local retrieval time")
	}
}

func TestScannerChoosesNewestTradeOrQuoteEvidence(t *testing.T) {
	var snap alpacaLiveSnapshot
	snap.LatestTrade.Time = "2026-08-20T14:01:00Z"
	snap.LatestQuote.Time = "2026-08-20T14:03:00Z"
	got := scannerScoreFromSnapshot("AAPL", "day", snap)
	want := providerTimeMillis(snap.LatestQuote.Time)
	if got.UpdatedAt != want {
		t.Fatalf("scanner evidence time = %d, want newest provider observation %d", got.UpdatedAt, want)
	}
}

func TestOpportunityRadarPreservesScannerEvidenceTime(t *testing.T) {
	var snap alpacaLiveSnapshot
	snap.LatestTrade.Time = "2026-08-20T14:04:00Z"
	base := scannerScoreFromSnapshot("AAPL", "day", snap)
	got := enrichOpportunityMetrics(base, snap, "regular", time.Date(2026, 8, 20, 15, 0, 0, 0, time.UTC), false, false)
	if got.UpdatedAt != base.UpdatedAt {
		t.Fatalf("radar changed shared evidence time: got %d, scanner %d", got.UpdatedAt, base.UpdatedAt)
	}
}
