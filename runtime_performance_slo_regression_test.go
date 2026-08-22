package main

import (
	"context"
	"testing"
	"time"
)

func sloCheckByName(a RuntimeSLOAssessment, name string) (RuntimeSLOCheck, bool) {
	for _, c := range a.Checks {
		if c.Name == name {
			return c, true
		}
	}
	return RuntimeSLOCheck{}, false
}

func currentQuotesForSymbols(symbols []string, at time.Time) map[string]Quote {
	out := map[string]Quote{}
	for _, s := range uniqueSymbols(symbols) {
		out[s] = Quote{Symbol: s, Price: 100, Bid: 99.99, Ask: 100.01, Source: "test-live", FeedType: "websocket", DataState: "live", ProviderTimestamp: at.UnixMilli(), UpdatedAt: at.UnixMilli()}
	}
	return out
}

func TestV173SelectedAndActionableFreshnessSLOsAreExplicit(t *testing.T) {
	now := time.Date(2026, time.August, 12, 14, 0, 0, 0, easternLocation())
	st := defaultState()
	quotes := currentQuotesForSymbols(actionableSLOSymbols(st), now)
	fresh := []FreshnessDiagnostic{{Dataset: "Quotes", State: "LIVE"}, {Dataset: "VIX", State: "FRESH"}}
	load := RuntimeLoadDiagnostics{GOMAXPROCS: 4, Startup: RuntimeStartupDiagnostics{WarmStartQuotes: 10, WarmStartTargetQuotes: 10, WarmStartCoveragePct: 100}}
	a := buildRuntimeSLOAssessmentWithContext("running", "live", FeedDiagnostics{MarketSession: "regular"}, fresh, load, ScannerState{}, st, quotes, now)
	if c, _ := sloCheckByName(a, "Selected-symbol freshness"); c.Status != "PASS" {
		t.Fatalf("selected SLO=%+v", c)
	}
	if c, _ := sloCheckByName(a, "Actionable-watchlist freshness"); c.Status != "PASS" {
		t.Fatalf("watchlist SLO=%+v", c)
	}
	q := quotes[st.UI.SelectedTicker]
	q.ProviderTimestamp = now.Add(-10 * time.Minute).UnixMilli()
	q.UpdatedAt = q.ProviderTimestamp
	quotes[st.UI.SelectedTicker] = q
	a = buildRuntimeSLOAssessmentWithContext("running", "live", FeedDiagnostics{MarketSession: "regular"}, fresh, load, ScannerState{}, st, quotes, now)
	if c, _ := sloCheckByName(a, "Selected-symbol freshness"); c.Status != "BLOCK" {
		t.Fatalf("stale selected symbol must block: %+v", c)
	}
}

func TestV173RecoveryTrackerMeasuresStaleAndDegradationRecovery(t *testing.T) {
	tr := NewRuntimeSLOTracker()
	base := time.Unix(1000, 0)
	tr.Observe([]FreshnessDiagnostic{{Dataset: "Quotes", State: "STALE"}}, RuntimeDegradationState{Code: "NETWORK"}, base)
	d := tr.Observe([]FreshnessDiagnostic{{Dataset: "Quotes", State: "LIVE"}}, RuntimeDegradationState{}, base.Add(45*time.Second))
	if d.StaleToCurrentEvents != 1 || d.LastStaleToCurrentMs != 45000 {
		t.Fatalf("stale recovery=%+v", d)
	}
	if d.DegradationRecoveryEvents != 1 || d.LastDegradationRecoveryMs != 45000 || d.DegradedSince != 0 {
		t.Fatalf("degradation recovery=%+v", d)
	}
}

func TestV173PersistenceWriteRateIsObservable(t *testing.T) {
	p := NewPersistenceManager(t.TempDir())
	defer p.Close()
	if !p.Diagnostics().Ready {
		t.Fatal("persistence not ready")
	}
	now := time.Now().UnixMilli()
	p.EnqueueQuotes(map[string]Quote{"AAPL": {Symbol: "AAPL", Price: 200, ProviderTimestamp: now, UpdatedAt: now, Source: "test"}})
	deadline := time.Now().Add(2 * time.Second)
	for time.Now().Before(deadline) {
		d := p.Diagnostics()
		if d.WriteBatchesLastMin > 0 && d.RowsWrittenLastMin > 0 {
			return
		}
		time.Sleep(10 * time.Millisecond)
	}
	t.Fatalf("recent DB write telemetry missing: %+v", p.Diagnostics())
}

func TestV173RuntimeLoadSamplesCPUStartupReuseAndStorage(t *testing.T) {
	dir := t.TempDir()
	p := NewPersistenceManager(dir)
	defer p.Close()
	app := &Application{configDir: dir, hub: NewHub(), persistence: p, state: defaultState(), aiCache: map[string]aiCacheEntry{}, httpTelemetry: NewRequestTelemetry()}
	e := NewEngine(app)
	e.mu.Lock()
	e.providerCallsAvoided = 20
	e.mu.Unlock()
	for i := 0; i < 100000; i++ {
		_ = i * i
	}
	time.Sleep(5 * time.Millisecond)
	e.sampleRuntimeLoad()
	e.mu.RLock()
	d := e.runtimeLoad
	e.mu.RUnlock()
	if d.GOMAXPROCS < 1 || d.CPUUtilizationPct < 0 {
		t.Fatalf("CPU telemetry=%+v", d)
	}
	if d.Startup.BootstrapDurationMs < 0 || d.Startup.WarmStartTargetQuotes == 0 {
		t.Fatalf("startup telemetry=%+v", d.Startup)
	}
	if d.ProviderCallsAvoided != 20 || d.CanonicalReuseHitRatePct < 0 {
		t.Fatalf("reuse telemetry=%+v", d)
	}
}

func TestV173ActiveMarketPressureProtectsCriticalWork(t *testing.T) {
	w := NewWorkloadController()
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	releases := []func(){}
	// Fill all provider slots with Tier 0/1 work; lower tiers must shed/reject before critical work is reclassified.
	for i := 0; i < 4; i++ {
		rel, ok := w.TryAcquireTier("provider-rest", WorkTierUserActionable)
		if !ok {
			t.Fatalf("critical acquire %d rejected", i)
		}
		releases = append(releases, rel)
	}
	if _, ok := w.TryAcquireTier("provider-rest", WorkTierBackground); ok {
		t.Fatal("background work entered saturated critical provider path")
	}
	for _, rel := range releases {
		rel()
	}
	rel, ok := w.AcquireTier(ctx, "provider-rest", WorkTierMarketCritical)
	if !ok {
		t.Fatal("Tier 0 failed after pressure released")
	}
	rel()
}

func TestV173SLOBudgetsBlockOnlyMeasuredSeverePressure(t *testing.T) {
	load := RuntimeLoadDiagnostics{GOMAXPROCS: 4, CPUUtilizationPct: 95, HeapAllocBytes: 100 * 1024 * 1024, Goroutines: 100, HTTP: HTTPRuntimeDiagnostics{InteractiveP95Ms: 1200}, Persistence: PersistenceDiagnostics{RowsWrittenLastMin: 3500}, Startup: RuntimeStartupDiagnostics{BootstrapDurationMs: 100, WarmStartTargetQuotes: 1, WarmStartQuotes: 1, WarmStartCoveragePct: 100}}
	fresh := []FreshnessDiagnostic{{Dataset: "Quotes", State: "LIVE"}, {Dataset: "VIX", State: "FRESH"}}
	a := buildRuntimeSLOAssessment("running", "live", FeedDiagnostics{MarketSession: "regular"}, fresh, load, ScannerState{})
	if a.Status != "BLOCK" {
		t.Fatalf("severe measured pressure must block: %+v", a)
	}
	for _, name := range []string{"Interactive API p95", "CPU utilization", "DB write rate"} {
		if c, ok := sloCheckByName(a, name); !ok || c.Status != "BLOCK" {
			t.Fatalf("%s=%+v ok=%v", name, c, ok)
		}
	}
}
