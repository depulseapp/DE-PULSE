package main

import "testing"

func TestV170CriticalDecisionUsabilityRequiresActualCriticalRows(t *testing.T) {
	if criticalDecisionDataUsable(nil, "closed") {
		t.Fatal("missing critical freshness rows must not be reported usable")
	}
	rows := []FreshnessDiagnostic{{Dataset: "Quotes", State: "IDLE"}, {Dataset: "VIX", State: "IDLE"}}
	if !criticalDecisionDataUsable(rows, "closed") {
		t.Fatal("closed-session IDLE critical datasets should be usable when explicitly present")
	}
	rows[0].State = "STALE"
	if criticalDecisionDataUsable(rows, "regular") {
		t.Fatal("regular-session stale quotes must make critical decision data unusable")
	}
}

func TestV170DegradationPrefersLocalLoadAndReportsCriticalUsability(t *testing.T) {
	rows := []FreshnessDiagnostic{{Dataset: "Quotes", State: "LIVE"}, {Dataset: "VIX", State: "FRESH"}}
	load := RuntimeLoadDiagnostics{Workload: []WorkClassDiagnostics{{Class: "provider-rest", Queued: 2, OldestQueueAgeMs: 2500}}}
	got := deriveRuntimeDegradation("running", "live", FeedDiagnostics{MarketSession: "regular"}, rows, ProviderRouterSnapshot{}, load)
	if got.Code != "LOCAL LOAD" || !got.CriticalUsable {
		t.Fatalf("expected LOCAL LOAD with usable critical data, got %+v", got)
	}
}

func TestV170DegradationDistinguishesRateLimitAndPartialCoverage(t *testing.T) {
	rows := []FreshnessDiagnostic{{Dataset: "Quotes", State: "LIVE"}, {Dataset: "VIX", State: "FRESH"}, {Dataset: "News", State: "STALE"}}
	load := RuntimeLoadDiagnostics{ProviderRequests: []ProviderRequestDiagnostics{{Provider: "Finnhub", RequestsLastMin: 3, RateLimited: 1}}}
	got := deriveRuntimeDegradation("degraded", "live", FeedDiagnostics{MarketSession: "regular"}, rows, ProviderRouterSnapshot{}, load)
	if got.Code != "RATE LIMITED" {
		t.Fatalf("expected RATE LIMITED, got %+v", got)
	}
	load = RuntimeLoadDiagnostics{}
	got = deriveRuntimeDegradation("degraded", "live", FeedDiagnostics{MarketSession: "regular"}, rows, ProviderRouterSnapshot{}, load)
	if got.Code != "PARTIAL COVERAGE" || !got.CriticalUsable {
		t.Fatalf("expected PARTIAL COVERAGE with critical data usable, got %+v", got)
	}
}

func TestV170RuntimeSLOBlocksCriticalTruthButOnlyWarnsModerateLatency(t *testing.T) {
	fresh := []FreshnessDiagnostic{{Dataset: "Quotes", State: "LIVE"}, {Dataset: "VIX", State: "FRESH"}}
	load := RuntimeLoadDiagnostics{HTTP: HTTPRuntimeDiagnostics{InteractiveP95Ms: 400}, Goroutines: 100, HeapAllocBytes: 64 * 1024 * 1024}
	got := buildRuntimeSLOAssessment("running", "live", FeedDiagnostics{MarketSession: "regular"}, fresh, load, ScannerState{})
	if got.Status != "WARN" {
		t.Fatalf("expected WARN for moderate API latency, got %+v", got)
	}
	stale := []FreshnessDiagnostic{{Dataset: "Quotes", State: "STALE"}, {Dataset: "VIX", State: "FRESH"}}
	got = buildRuntimeSLOAssessment("running", "live", FeedDiagnostics{MarketSession: "regular"}, stale, load, ScannerState{})
	if got.Status != "BLOCK" {
		t.Fatalf("regular-session stale critical quote truth must block, got %+v", got)
	}
}
