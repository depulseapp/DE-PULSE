package main

import (
	"context"
	"reflect"
	"sync/atomic"
	"testing"
	"time"
)

func broadSnapshotFixture(symbols []string) map[string]alpacaLiveSnapshot {
	out := make(map[string]alpacaLiveSnapshot, len(symbols))
	for _, symbol := range symbols {
		out[normalizeSymbol(symbol)] = alpacaLiveSnapshot{}
	}
	return out
}

func TestV186BroadSnapshotBrokerReusesFreshObservations(t *testing.T) {
	broker := NewBroadSnapshotBroker(16)
	now := time.Date(2026, 8, 18, 14, 0, 0, 0, easternLocation())
	var calls atomic.Int64
	fetch := func(_ context.Context, symbols []string) (map[string]alpacaLiveSnapshot, error) {
		calls.Add(1)
		return broadSnapshotFixture(symbols), nil
	}

	first, firstResult, err := broker.Acquire(context.Background(), "iex", []string{"AAPL", "MSFT"}, time.Minute, now, fetch)
	if err != nil {
		t.Fatal(err)
	}
	if firstResult.ProviderCallAvoided || calls.Load() != 1 || len(first) != 2 {
		t.Fatalf("unexpected first acquisition: result=%+v calls=%d rows=%d", firstResult, calls.Load(), len(first))
	}

	second, secondResult, err := broker.Acquire(context.Background(), "iex", []string{"MSFT", "AAPL"}, time.Minute, now.Add(10*time.Second), fetch)
	if err != nil {
		t.Fatal(err)
	}
	if !secondResult.ProviderCallAvoided || calls.Load() != 1 || secondResult.SymbolsReused != 2 || len(second) != 2 {
		t.Fatalf("fresh reuse failed: result=%+v calls=%d rows=%d", secondResult, calls.Load(), len(second))
	}
	diag := broker.Diagnostics()
	if diag.ProviderFetches != 1 || diag.CacheOnlyRequests != 1 || diag.SymbolsReused < 2 {
		t.Fatalf("unexpected diagnostics after reuse: %+v", diag)
	}
}

func TestV186BroadSnapshotBrokerFetchesOnlyMissingSymbols(t *testing.T) {
	broker := NewBroadSnapshotBroker(16)
	now := time.Date(2026, 8, 18, 14, 0, 0, 0, easternLocation())
	var batches [][]string
	fetch := func(_ context.Context, symbols []string) (map[string]alpacaLiveSnapshot, error) {
		batches = append(batches, append([]string{}, symbols...))
		return broadSnapshotFixture(symbols), nil
	}

	if _, _, err := broker.Acquire(context.Background(), "iex", []string{"AAPL", "MSFT"}, time.Minute, now, fetch); err != nil {
		t.Fatal(err)
	}
	rows, result, err := broker.Acquire(context.Background(), "iex", []string{"MSFT", "NVDA"}, time.Minute, now.Add(5*time.Second), fetch)
	if err != nil {
		t.Fatal(err)
	}
	if len(batches) != 2 || !reflect.DeepEqual(batches[1], []string{"NVDA"}) {
		t.Fatalf("expected only missing NVDA to fetch, batches=%v", batches)
	}
	if result.SymbolsReused != 1 || result.SymbolsFetched != 1 || len(rows) != 2 {
		t.Fatalf("unexpected partial-reuse result: %+v rows=%d", result, len(rows))
	}
}

func TestV186BroadSnapshotBrokerCoalescesInFlightRequests(t *testing.T) {
	broker := NewBroadSnapshotBroker(16)
	now := time.Date(2026, 8, 18, 14, 0, 0, 0, easternLocation())
	started := make(chan struct{})
	release := make(chan struct{})
	var calls atomic.Int64
	fetch := func(_ context.Context, symbols []string) (map[string]alpacaLiveSnapshot, error) {
		if calls.Add(1) == 1 {
			close(started)
		}
		<-release
		return broadSnapshotFixture(symbols), nil
	}

	type outcome struct {
		result BroadSnapshotAcquireResult
		err    error
	}
	first := make(chan outcome, 1)
	second := make(chan outcome, 1)
	go func() {
		_, result, err := broker.Acquire(context.Background(), "iex", []string{"AAPL", "MSFT"}, time.Minute, now, fetch)
		first <- outcome{result: result, err: err}
	}()
	<-started
	go func() {
		_, result, err := broker.Acquire(context.Background(), "iex", []string{"AAPL", "MSFT"}, time.Minute, now.Add(time.Second), fetch)
		second <- outcome{result: result, err: err}
	}()
	time.Sleep(20 * time.Millisecond)
	close(release)

	one, two := <-first, <-second
	if one.err != nil || two.err != nil {
		t.Fatalf("coalesced acquisition errors: first=%v second=%v", one.err, two.err)
	}
	if calls.Load() != 1 {
		t.Fatalf("expected one provider fetch for overlapping in-flight requests, got %d", calls.Load())
	}
	if !one.result.Coalesced && !two.result.Coalesced {
		t.Fatalf("expected one waiter to report coalescing: first=%+v second=%+v", one.result, two.result)
	}
	diag := broker.Diagnostics()
	if diag.CoalescedWaiters < 1 || diag.ProviderFetches != 1 {
		t.Fatalf("unexpected coalescing diagnostics: %+v", diag)
	}
}

func TestV186BroadSnapshotBrokerExpiresAndBoundsCache(t *testing.T) {
	broker := NewBroadSnapshotBroker(2)
	now := time.Date(2026, 8, 18, 14, 0, 0, 0, easternLocation())
	var calls atomic.Int64
	fetch := func(_ context.Context, symbols []string) (map[string]alpacaLiveSnapshot, error) {
		calls.Add(1)
		return broadSnapshotFixture(symbols), nil
	}

	if _, _, err := broker.Acquire(context.Background(), "iex", []string{"AAPL"}, 10*time.Second, now, fetch); err != nil {
		t.Fatal(err)
	}
	if _, _, err := broker.Acquire(context.Background(), "iex", []string{"AAPL"}, 10*time.Second, now.Add(11*time.Second), fetch); err != nil {
		t.Fatal(err)
	}
	if calls.Load() != 2 {
		t.Fatalf("expired snapshot should refetch, calls=%d", calls.Load())
	}
	if _, _, err := broker.Acquire(context.Background(), "iex", []string{"MSFT"}, time.Minute, now.Add(12*time.Second), fetch); err != nil {
		t.Fatal(err)
	}
	if _, _, err := broker.Acquire(context.Background(), "iex", []string{"NVDA"}, time.Minute, now.Add(13*time.Second), fetch); err != nil {
		t.Fatal(err)
	}
	diag := broker.Diagnostics()
	if diag.Entries > 2 || diag.Evictions < 1 {
		t.Fatalf("bounded cache contract failed: %+v", diag)
	}
}

func TestV186BroadSnapshotRequestCanonicalization(t *testing.T) {
	raw := "https://data.alpaca.markets/v2/stocks/snapshots?symbols=MSFT%2CAAPL%2CMSFT&feed=iex"
	feed, symbols, ok := parseBroadSnapshotRequest(raw)
	if !ok || feed != "iex" || !reflect.DeepEqual(symbols, []string{"AAPL", "MSFT"}) {
		t.Fatalf("unexpected parsed request feed=%q symbols=%v ok=%v", feed, symbols, ok)
	}
	rewritten, err := broadSnapshotURLForSymbols(raw, []string{"NVDA", "AAPL"})
	if err != nil {
		t.Fatal(err)
	}
	feed, symbols, ok = parseBroadSnapshotRequest(rewritten)
	if !ok || feed != "iex" || !reflect.DeepEqual(symbols, []string{"AAPL", "NVDA"}) {
		t.Fatalf("unexpected rewritten request feed=%q symbols=%v ok=%v url=%s", feed, symbols, ok, rewritten)
	}
}
