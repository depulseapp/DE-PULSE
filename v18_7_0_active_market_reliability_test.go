package main

import (
	"context"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

func waitV1870ProviderQueue(t *testing.T, w *WorkloadController, want int) WorkClassDiagnostics {
	t.Helper()
	deadline := time.Now().Add(time.Second)
	for time.Now().Before(deadline) {
		for _, row := range w.Diagnostics() {
			if row.Class == "provider-rest" && row.Queued >= want {
				return row
			}
		}
		time.Sleep(5 * time.Millisecond)
	}
	for _, row := range w.Diagnostics() {
		if row.Class == "provider-rest" {
			t.Fatalf("provider queue never reached %d: %+v", want, row)
		}
	}
	t.Fatalf("provider-rest diagnostics missing")
	return WorkClassDiagnostics{}
}

func TestV1870ActiveMarketDuplicateSnapshotBurstCoalesces(t *testing.T) {
	broker := NewBroadSnapshotBroker(64)
	now := time.Date(2026, 8, 19, 14, 0, 0, 0, easternLocation())
	started := make(chan struct{})
	releaseFetch := make(chan struct{})
	var calls atomic.Int64
	fetch := func(_ context.Context, symbols []string) (map[string]alpacaLiveSnapshot, error) {
		if calls.Add(1) == 1 {
			close(started)
		}
		<-releaseFetch
		return broadSnapshotFixture(symbols), nil
	}

	const requesters = 16
	type result struct {
		rows int
		err  error
	}
	results := make(chan result, requesters)
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		rows, _, err := broker.Acquire(context.Background(), "iex", []string{"AAPL", "MSFT", "NVDA", "META", "SPY"}, 30*time.Second, now, fetch)
		results <- result{rows: len(rows), err: err}
	}()
	<-started

	for i := 1; i < requesters; i++ {
		wg.Add(1)
		go func(offset int) {
			defer wg.Done()
			rows, _, err := broker.Acquire(context.Background(), "iex", []string{"SPY", "META", "NVDA", "MSFT", "AAPL"}, 30*time.Second, now.Add(time.Duration(offset)*time.Millisecond), fetch)
			results <- result{rows: len(rows), err: err}
		}(i)
	}
	time.Sleep(30 * time.Millisecond)
	close(releaseFetch)
	wg.Wait()
	close(results)

	for got := range results {
		if got.err != nil || got.rows != 5 {
			t.Fatalf("snapshot burst result invalid: rows=%d err=%v", got.rows, got.err)
		}
	}
	if calls.Load() != 1 {
		t.Fatalf("duplicate active-market snapshot burst must collapse to one provider call, got %d", calls.Load())
	}
	diag := broker.Diagnostics()
	if diag.ProviderFetches != 1 || diag.CoalescedWaiters < 1 {
		t.Fatalf("coalescing must be observable under burst load: %+v", diag)
	}
}

func TestV1870ActiveMarketProviderPressureIsBoundedAndTruthful(t *testing.T) {
	w := NewWorkloadController()
	var provider WorkClassDiagnostics
	for _, row := range w.Diagnostics() {
		if row.Class == "provider-rest" {
			provider = row
			break
		}
	}
	if provider.Capacity <= 0 || provider.MaxQueue <= 0 || provider.ReservedCritical <= 0 {
		t.Fatalf("provider pressure contract unavailable: %+v", provider)
	}

	holds := make([]func(), 0, provider.Capacity)
	for i := 0; i < provider.Capacity; i++ {
		release, ok := w.TryAcquireTier("provider-rest", WorkTierMarketCritical)
		if !ok {
			t.Fatalf("failed to establish provider pressure at slot %d/%d", i+1, provider.Capacity)
		}
		holds = append(holds, release)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	done := make(chan struct{}, provider.MaxQueue)
	for i := 0; i < provider.MaxQueue; i++ {
		go func() {
			release, ok := w.AcquireTier(ctx, "provider-rest", WorkTierBackground)
			if ok {
				release()
			}
			done <- struct{}{}
		}()
	}
	pressure := waitV1870ProviderQueue(t, w, provider.MaxQueue)
	if pressure.Queued != pressure.MaxQueue {
		t.Fatalf("provider queue must stop exactly at its hard bound: %+v", pressure)
	}

	fresh := []FreshnessDiagnostic{{Dataset: "Quotes", State: "LIVE"}, {Dataset: "VIX", State: "FRESH"}}
	load := RuntimeLoadDiagnostics{Workload: w.Diagnostics()}
	degraded := deriveRuntimeDegradation("running", "live", FeedDiagnostics{MarketSession: "regular", FeedState: "streaming"}, fresh, ProviderRouterSnapshot{}, load)
	if degraded.ReasonCode != "QUEUE_SATURATED" || degraded.PressureState != "PROTECTED" || degraded.Abstain || !degraded.CriticalUsable {
		t.Fatalf("bounded local overload must be explicit without falsely invalidating current critical evidence: %+v", degraded)
	}

	slo := buildRuntimeSLOAssessmentWithContext("running", "live", FeedDiagnostics{MarketSession: "regular", FeedState: "streaming"}, fresh, load, ScannerState{}, AppState{}, nil, time.Now())
	if slo.Status != "BLOCK" {
		t.Fatalf("full provider queue must fail active-market capacity SLO even while critical evidence remains usable: %+v", slo)
	}

	if release, ok := w.TryAcquireTier("provider-rest", WorkTierBackground); ok {
		release()
		t.Fatal("work beyond saturated provider capacity/queue must not be admitted immediately")
	}

	cancel()
	for _, release := range holds {
		release()
	}
	for i := 0; i < provider.MaxQueue; i++ {
		select {
		case <-done:
		case <-time.After(time.Second):
			t.Fatal("queued active-market worker did not cancel cleanly")
		}
	}
}
