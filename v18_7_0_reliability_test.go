package main

import (
	"testing"
	"time"
)

func containsString(values []string, want string) bool {
	for _, value := range values {
		if value == want {
			return true
		}
	}
	return false
}

func TestV1870CanonicalReasonTaxonomyKeepsCompatibilityLabel(t *testing.T) {
	fresh := []FreshnessDiagnostic{{Dataset: "Quotes", State: "LIVE"}, {Dataset: "VIX", State: "FRESH"}}
	load := RuntimeLoadDiagnostics{Workload: []WorkClassDiagnostics{{Class: "provider-rest", Queued: 8, MaxQueue: 8, OldestQueueAgeMs: 2500}}}
	got := deriveRuntimeDegradation("running", "live", FeedDiagnostics{MarketSession: "regular"}, fresh, ProviderRouterSnapshot{}, load)
	if got.Code != "LOCAL LOAD" {
		t.Fatalf("compatibility display label changed unexpectedly: %+v", got)
	}
	if got.ReasonCode != "QUEUE_SATURATED" {
		t.Fatalf("expected canonical QUEUE_SATURATED reason, got %+v", got)
	}
	if got.PressureState != "PROTECTED" || got.Abstain || !got.CriticalUsable {
		t.Fatalf("local queue pressure must isolate optional work while critical truth remains usable: %+v", got)
	}
	if !containsString(got.AffectedConsumers, "Opportunity Radar") {
		t.Fatalf("expected broad-refresh blast radius to name Opportunity Radar: %+v", got)
	}
}

func TestV1870CriticalNetworkFailureRequiresAbstain(t *testing.T) {
	fresh := []FreshnessDiagnostic{{Dataset: "Quotes", State: "STALE"}, {Dataset: "VIX", State: "FRESH"}}
	got := deriveRuntimeDegradation("degraded", "live", FeedDiagnostics{MarketSession: "regular", FeedState: "reconnecting"}, fresh, ProviderRouterSnapshot{}, RuntimeLoadDiagnostics{})
	if got.ReasonCode != "NETWORK_FAILURE" || got.PressureState != "DEGRADED" || !got.Abstain || got.CriticalUsable {
		t.Fatalf("critical live-feed failure must be explicit and fail closed: %+v", got)
	}
	if !containsString(got.AffectedConsumers, "Day") || !containsString(got.AffectedConsumers, "Decision Queue") {
		t.Fatalf("critical quote blast radius is incomplete: %+v", got)
	}
}

func TestV1870OptionalStalenessIsolatedFromCriticalDecisionTruth(t *testing.T) {
	fresh := []FreshnessDiagnostic{
		{Dataset: "Quotes", State: "LIVE"},
		{Dataset: "VIX", State: "FRESH"},
		{Dataset: "Fundamentals", State: "STALE"},
	}
	got := deriveRuntimeDegradation("degraded", "live", FeedDiagnostics{MarketSession: "regular"}, fresh, ProviderRouterSnapshot{}, RuntimeLoadDiagnostics{})
	if got.Code != "PARTIAL COVERAGE" || got.ReasonCode != "LOW_COVERAGE" {
		t.Fatalf("expected capability-scoped partial coverage: %+v", got)
	}
	if !got.CriticalUsable || got.Abstain || got.PressureState != "PROTECTED" {
		t.Fatalf("optional fundamentals staleness must not collapse critical live truth: %+v", got)
	}
	if !containsString(got.AffectedConsumers, "Long") || !containsString(got.AffectedConsumers, "Research") {
		t.Fatalf("fundamentals blast radius must remain scoped to its consumers: %+v", got)
	}
	if containsString(got.AffectedConsumers, "Day") {
		t.Fatalf("fundamentals-only degradation must not contaminate Day: %+v", got)
	}
}

func TestV1870HealthyStateHasNoFalseAbstain(t *testing.T) {
	fresh := []FreshnessDiagnostic{{Dataset: "Quotes", State: "LIVE"}, {Dataset: "VIX", State: "FRESH"}}
	got := deriveRuntimeDegradation("running", "live", FeedDiagnostics{MarketSession: "regular", FeedState: "streaming"}, fresh, ProviderRouterSnapshot{}, RuntimeLoadDiagnostics{})
	if got.Code != "" || got.ReasonCode != "" || got.PressureState != "HEALTHY" || got.Abstain {
		t.Fatalf("healthy runtime should remain clean and non-blocking: %+v", got)
	}
}

func TestV1870RecoveryHysteresisRejectsTransientHealthySample(t *testing.T) {
	tracker := NewRuntimeSLOTracker()
	t0 := time.Unix(1_800_000_000, 0)
	degraded := RuntimeDegradationState{Code: "NETWORK", ReasonCode: "NETWORK_FAILURE", PressureState: "DEGRADED", CriticalUsable: false, Abstain: true, Detail: "feeds reconnecting"}
	healthy := RuntimeDegradationState{PressureState: "HEALTHY", CriticalUsable: true}

	if got := tracker.StabilizeDegradation(degraded, t0); got.Code != "NETWORK" {
		t.Fatalf("active degradation must surface immediately: %+v", got)
	}
	tracker.Observe(nil, degraded, t0)

	first := tracker.StabilizeDegradation(healthy, t0.Add(time.Second))
	if first.Code != "NETWORK" || first.PressureState != "RECOVERING" || !first.Abstain {
		t.Fatalf("one healthy observation must not clear a critical degradation: %+v", first)
	}
	d1 := tracker.Observe(nil, first, t0.Add(time.Second))
	if !d1.DegradationRecoveryPending || d1.DegradationHealthyObservations != 1 || d1.DegradationHealthyRequired != 3 {
		t.Fatalf("recovery confirmation diagnostics are incomplete: %+v", d1)
	}

	second := tracker.StabilizeDegradation(healthy, t0.Add(3*time.Second))
	if second.Code == "" {
		t.Fatalf("recovery cleared before the required confirmation count/stability window")
	}
	tracker.Observe(nil, second, t0.Add(3*time.Second))

	third := tracker.StabilizeDegradation(healthy, t0.Add(6*time.Second))
	if third.Code != "" || third.PressureState != "HEALTHY" {
		t.Fatalf("confirmed stable recovery should clear the held degradation: %+v", third)
	}
	d3 := tracker.Observe(nil, third, t0.Add(6*time.Second))
	if d3.DegradationRecoveryPending || d3.DegradationRecoveryEvents != 1 || d3.LastDegradationRecoveryMs < 6000 {
		t.Fatalf("confirmed recovery should be recorded exactly once: %+v", d3)
	}
}

func TestV1870RecoveryHysteresisResetsAfterRelapse(t *testing.T) {
	tracker := NewRuntimeSLOTracker()
	t0 := time.Unix(1_800_000_000, 0)
	degraded := RuntimeDegradationState{Code: "PROVIDER DEGRADED", ReasonCode: "PROVIDER_DOWN", PressureState: "DEGRADED", CriticalUsable: false, Abstain: true}
	healthy := RuntimeDegradationState{PressureState: "HEALTHY", CriticalUsable: true}

	tracker.StabilizeDegradation(degraded, t0)
	tracker.StabilizeDegradation(healthy, t0.Add(time.Second))
	tracker.StabilizeDegradation(degraded, t0.Add(2*time.Second))
	afterRelapse := tracker.StabilizeDegradation(healthy, t0.Add(8*time.Second))
	if afterRelapse.Code == "" || afterRelapse.PressureState != "RECOVERING" {
		t.Fatalf("a relapse must reset the healthy streak even when wall-clock time has elapsed: %+v", afterRelapse)
	}
}
