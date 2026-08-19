package main

import "testing"

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
