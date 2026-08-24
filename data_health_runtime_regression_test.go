package main

import (
	"strings"
	"testing"
	"time"
)

func adaptFreshCritical() []FreshnessDiagnostic {
	return []FreshnessDiagnostic{
		{Dataset: "Quotes", State: "LIVE"},
		{Dataset: "VIX", State: "FRESH"},
	}
}

func adaptContainsFold(values []string, needle string) bool {
	for _, value := range values {
		if strings.EqualFold(strings.TrimSpace(value), strings.TrimSpace(needle)) {
			return true
		}
	}
	return false
}

func TestADAPTDataHealthValidWarmEvidenceSurvivesTemporaryCriticalRouteOutage(t *testing.T) {
	router := ProviderRouterSnapshot{Routes: []ProviderRouteState{{Dataset: "US Live Equities", State: "UNAVAILABLE"}}}
	got := deriveRuntimeDegradation("degraded", "live", FeedDiagnostics{MarketSession: "regular", FeedState: "reconnecting"}, adaptFreshCritical(), router, RuntimeLoadDiagnostics{})
	if got.Code != "" || got.ReasonCode != "" || got.PressureState != "HEALTHY" || got.Abstain || !got.CriticalUsable {
		t.Fatalf("valid warm evidence must remain usable until canonical freshness expires: %+v", got)
	}
	if !got.WarmStateActive {
		t.Fatalf("temporary route outage must expose warm-state reuse diagnostics without extending evidence time: %+v", got)
	}
	if len(got.TransportIssues) == 0 {
		t.Fatalf("temporary route outage must remain visible as non-secret transport context: %+v", got)
	}
}

func TestADAPTDataHealthStaleWarmEvidenceCannotMaskCriticalRouteOutage(t *testing.T) {
	freshness := []FreshnessDiagnostic{{Dataset: "Quotes", State: "STALE"}, {Dataset: "VIX", State: "FRESH"}}
	router := ProviderRouterSnapshot{Routes: []ProviderRouteState{{Dataset: "US Live Equities", State: "UNAVAILABLE"}}}
	got := deriveRuntimeDegradation("degraded", "live", FeedDiagnostics{MarketSession: "regular", FeedState: "connected-idle"}, freshness, router, RuntimeLoadDiagnostics{})
	if got.Code != "PARTIAL COVERAGE" || got.ReasonCode != "LOW_COVERAGE" || got.PressureState != "DEGRADED" || !got.Abstain || got.CriticalUsable {
		t.Fatalf("stale required quote evidence must remain truthfully degraded: %+v", got)
	}
	if !adaptContainsFold(got.Affected, "Quotes") {
		t.Fatalf("critical degradation must identify the stale dataset: %+v", got)
	}
}

func TestADAPTDataHealthOptionalRouteFailureDoesNotCreateFalseGlobalDegradation(t *testing.T) {
	router := ProviderRouterSnapshot{Routes: []ProviderRouteState{{Dataset: "News", State: "UNAVAILABLE"}}}
	got := deriveRuntimeDegradation("degraded", "live", FeedDiagnostics{MarketSession: "regular", FeedState: "connected-idle"}, adaptFreshCritical(), router, RuntimeLoadDiagnostics{})
	if got.Code != "" || got.ReasonCode != "" || got.PressureState != "HEALTHY" || got.Abstain {
		t.Fatalf("optional route failure must remain isolated while required canonical evidence is usable: %+v", got)
	}
	if !got.WarmStateActive || len(got.TransportIssues) == 0 {
		t.Fatalf("isolated optional route failure must remain diagnosable without creating a global degraded state: %+v", got)
	}
}

func TestADAPTDataHealthFallbackRateLimitDoesNotInvalidateFreshEvidence(t *testing.T) {
	router := ProviderRouterSnapshot{Routes: []ProviderRouteState{{
		Dataset: "US Live Equities",
		State:   "DEGRADED",
		Route:   []ProviderRouteHop{{Provider: "Alpaca", RateLimit: "RATE LIMITED"}, {Provider: "Finnhub", Circuit: "CLOSED"}},
	}}}
	got := deriveRuntimeDegradation("degraded", "live", FeedDiagnostics{MarketSession: "regular", FeedState: "finnhub-fallback"}, adaptFreshCritical(), router, RuntimeLoadDiagnostics{})
	if got.Code != "" || got.PressureState != "HEALTHY" || got.Abstain || !got.CriticalUsable {
		t.Fatalf("rate-limit pressure with healthy serving fallback must not manufacture degradation: %+v", got)
	}
	if !got.FallbackActive || !got.WarmStateActive || len(got.TransportIssues) == 0 {
		t.Fatalf("fallback/rate-limit provenance must remain visible for diagnosis: %+v", got)
	}
}

func TestADAPTDataHealthRecoveredCanonicalStateCanUnlatchDegradedRuntimeStatus(t *testing.T) {
	tracker := NewRuntimeSLOTracker()
	t0 := time.Unix(1_800_000_000, 0)
	badFreshness := []FreshnessDiagnostic{{Dataset: "Quotes", State: "STALE"}, {Dataset: "VIX", State: "FRESH"}}
	degraded := deriveRuntimeDegradation("degraded", "live", FeedDiagnostics{MarketSession: "regular", FeedState: "connected-idle"}, badFreshness, ProviderRouterSnapshot{}, RuntimeLoadDiagnostics{})
	if degraded.Code == "" {
		t.Fatalf("fixture must begin degraded: %+v", degraded)
	}
	if got := tracker.StabilizeDegradation(degraded, t0); got.Code == "" {
		t.Fatalf("real degradation must become visible immediately: %+v", got)
	}

	healthyRaw := deriveRuntimeDegradation("degraded", "live", FeedDiagnostics{MarketSession: "regular", FeedState: "connected-idle"}, adaptFreshCritical(), ProviderRouterSnapshot{}, RuntimeLoadDiagnostics{})
	if healthyRaw.Code != "" {
		t.Fatalf("mutable runtime status alone must not re-latch a degradation after canonical evidence recovers: %+v", healthyRaw)
	}

	got := tracker.StabilizeDegradation(healthyRaw, t0.Add(time.Second))
	if got.Code == "" || got.PressureState != "RECOVERING" {
		t.Fatalf("first healthy sample must enter hysteresis recovery instead of unlatching immediately: %+v", got)
	}
	got = tracker.StabilizeDegradation(healthyRaw, t0.Add(3*time.Second))
	if got.Code == "" || got.PressureState != "RECOVERING" {
		t.Fatalf("second healthy sample must remain held during the stability window: %+v", got)
	}
	got = tracker.StabilizeDegradation(healthyRaw, t0.Add(6*time.Second))
	if got.Code != "" || got.PressureState != "HEALTHY" {
		t.Fatalf("sustained recovery must unlatch the stale degraded state: %+v", got)
	}
}
