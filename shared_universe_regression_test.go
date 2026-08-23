package main

import (
	"context"
	"reflect"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

func newV188UniverseTestEngine(t *testing.T) *Engine {
	t.Helper()
	return newTestApplication(t).engine
}

func TestV188RadarThenScannerReusesCanonicalUniverse(t *testing.T) {
	e := newV188UniverseTestEngine(t)
	var calls int32
	e.canonicalUSUniverseLoader = func(context.Context, string, string) ([]string, bool) {
		atomic.AddInt32(&calls, 1)
		return []string{"AAPL", "MSFT", "NVDA"}, true
	}
	now := time.Now()
	radar := e.opportunityUniverse(context.Background(), "key", "secret", now)
	scanner := e.scannerUniverse(context.Background(), "key", "secret")
	if got := atomic.LoadInt32(&calls); got != 1 {
		t.Fatalf("expected one shared universe load, got %d", got)
	}
	if !reflect.DeepEqual(radar, scanner) {
		t.Fatalf("radar/scanner universe mismatch: radar=%v scanner=%v", radar, scanner)
	}
}

func TestV188ScannerThenRadarReusesCanonicalUniverse(t *testing.T) {
	e := newV188UniverseTestEngine(t)
	var calls int32
	e.canonicalUSUniverseLoader = func(context.Context, string, string) ([]string, bool) {
		atomic.AddInt32(&calls, 1)
		return []string{"AAPL", "MSFT", "NVDA"}, true
	}
	scanner := e.scannerUniverse(context.Background(), "key", "secret")
	radar := e.opportunityUniverse(context.Background(), "key", "secret", time.Now())
	if got := atomic.LoadInt32(&calls); got != 1 {
		t.Fatalf("expected one shared universe load, got %d", got)
	}
	if !reflect.DeepEqual(scanner, radar) {
		t.Fatalf("scanner/radar universe mismatch: scanner=%v radar=%v", scanner, radar)
	}
}

func TestV188CanonicalUniverseTTLAndExpiry(t *testing.T) {
	e := newV188UniverseTestEngine(t)
	var calls int32
	e.canonicalUSUniverseLoader = func(context.Context, string, string) ([]string, bool) {
		n := atomic.AddInt32(&calls, 1)
		if n == 1 {
			return []string{"AAPL", "MSFT"}, true
		}
		return []string{"AAPL", "MSFT", "NVDA"}, true
	}
	now := time.Date(2026, 8, 19, 20, 0, 0, 0, time.UTC)
	first := e.canonicalUSSymbolUniverse(context.Background(), "key", "secret", now)
	beforeExpiry := e.canonicalUSSymbolUniverse(context.Background(), "key", "secret", now.Add(opportunityUniverseTTL-time.Millisecond))
	if got := atomic.LoadInt32(&calls); got != 1 {
		t.Fatalf("fresh 12h TTL should reuse cache, loads=%d", got)
	}
	if !reflect.DeepEqual(first, beforeExpiry) {
		t.Fatalf("universe changed inside TTL: first=%v cached=%v", first, beforeExpiry)
	}
	afterExpiry := e.canonicalUSSymbolUniverse(context.Background(), "key", "secret", now.Add(opportunityUniverseTTL))
	if got := atomic.LoadInt32(&calls); got != 2 {
		t.Fatalf("TTL expiry should cause exactly one refresh, loads=%d", got)
	}
	if len(afterExpiry) != 3 || afterExpiry[2] != "NVDA" {
		t.Fatalf("expected refreshed universe at TTL boundary, got %v", afterExpiry)
	}
}

func TestV188ConcurrentScannerRadarMissCoalesces(t *testing.T) {
	e := newV188UniverseTestEngine(t)
	var calls int32
	loaderStarted := make(chan struct{})
	releaseLoader := make(chan struct{})
	e.canonicalUSUniverseLoader = func(context.Context, string, string) ([]string, bool) {
		if atomic.AddInt32(&calls, 1) == 1 {
			close(loaderStarted)
		}
		<-releaseLoader
		return []string{"AAPL", "MSFT", "NVDA"}, true
	}

	ctx := context.Background()
	now := time.Now()
	start := make(chan struct{})
	entered := make(chan struct{}, 2)
	results := make(chan []string, 2)
	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		defer wg.Done()
		entered <- struct{}{}
		<-start
		results <- e.scannerUniverse(ctx, "key", "secret")
	}()
	go func() {
		defer wg.Done()
		entered <- struct{}{}
		<-start
		results <- e.opportunityUniverse(ctx, "key", "secret", now)
	}()
	<-entered
	<-entered
	close(start)
	<-loaderStarted
	close(releaseLoader)
	wg.Wait()
	close(results)

	if got := atomic.LoadInt32(&calls); got != 1 {
		t.Fatalf("concurrent Scanner/Radar misses should coalesce, loads=%d", got)
	}
	for rows := range results {
		if !reflect.DeepEqual(rows, []string{"AAPL", "MSFT", "NVDA"}) {
			t.Fatalf("unexpected coalesced result: %v", rows)
		}
	}
}

func TestV188FailedRefreshPreservesStaleTimestamp(t *testing.T) {
	e := newV188UniverseTestEngine(t)
	var calls int32
	e.canonicalUSUniverseLoader = func(context.Context, string, string) ([]string, bool) {
		if atomic.AddInt32(&calls, 1) == 1 {
			return []string{"AAPL", "MSFT"}, true
		}
		return nil, false
	}
	now := time.Date(2026, 8, 19, 20, 0, 0, 0, time.UTC)
	first := e.canonicalUSSymbolUniverse(context.Background(), "key", "secret", now)
	e.mu.RLock()
	originalAt := e.canonicalUSUniverseAt
	e.mu.RUnlock()

	failureAt := now.Add(opportunityUniverseTTL)
	stale := e.canonicalUSSymbolUniverse(context.Background(), "key", "secret", failureAt)
	e.mu.RLock()
	cachedAt := e.canonicalUSUniverseAt
	source := e.canonicalUSUniverseSource
	retryAt := e.canonicalUSUniverseRetryAt
	health := e.health["scanner-universe"]
	e.mu.RUnlock()

	if !reflect.DeepEqual(first, stale) {
		t.Fatalf("failed refresh should preserve prior universe: first=%v stale=%v", first, stale)
	}
	if cachedAt != originalAt {
		t.Fatalf("failed refresh must not make stale cache fresh: original=%d got=%d", originalAt, cachedAt)
	}
	if source != "stale-cache" || retryAt <= failureAt.UnixMilli() || !strings.Contains(health, "stale cache") {
		t.Fatalf("stale refresh diagnostics incorrect: source=%q retryAt=%d health=%q", source, retryAt, health)
	}
	if got := atomic.LoadInt32(&calls); got != 2 {
		t.Fatalf("expected initial load plus one failed refresh, loads=%d", got)
	}
}

func TestV188SeedFallbackIsIdentifiable(t *testing.T) {
	e := newV188UniverseTestEngine(t)
	var calls int32
	e.canonicalUSUniverseLoader = func(context.Context, string, string) ([]string, bool) {
		atomic.AddInt32(&calls, 1)
		return nil, false
	}
	now := time.Date(2026, 8, 19, 20, 0, 0, 0, time.UTC)
	rows := e.canonicalUSSymbolUniverse(context.Background(), "key", "secret", now)
	e.mu.RLock()
	cachedAt := e.canonicalUSUniverseAt
	source := e.canonicalUSUniverseSource
	retryAt := e.canonicalUSUniverseRetryAt
	health := e.health["scanner-universe"]
	e.mu.RUnlock()

	if !reflect.DeepEqual(rows, uniqueSymbols(discoverySeedUniverse)) {
		t.Fatalf("expected deterministic seed fallback, got %v", rows)
	}
	if cachedAt != 0 || source != "seed-fallback" || retryAt <= now.UnixMilli() || !strings.Contains(health, "seed fallback") {
		t.Fatalf("fallback diagnostics incorrect: cachedAt=%d source=%q retryAt=%d health=%q", cachedAt, source, retryAt, health)
	}
	_ = e.canonicalUSSymbolUniverse(context.Background(), "key", "secret", now.Add(time.Minute))
	if got := atomic.LoadInt32(&calls); got != 1 {
		t.Fatalf("retry backoff should prevent immediate repeated fallback loads, loads=%d", got)
	}
}

func TestV188ScannerRankingAndRadarPromotionRemainUnchanged(t *testing.T) {
	before := demoScannerResults("day")
	e := newV188UniverseTestEngine(t)
	e.canonicalUSUniverseLoader = func(context.Context, string, string) ([]string, bool) {
		return []string{"AAPL", "MSFT", "NVDA"}, true
	}
	_ = e.canonicalUSSymbolUniverse(context.Background(), "key", "secret", time.Now())
	after := demoScannerResults("day")
	if len(before) != len(after) {
		t.Fatalf("scanner ranking length changed: before=%d after=%d", len(before), len(after))
	}
	for i := range before {
		if before[i].Symbol != after[i].Symbol || before[i].Score != after[i].Score {
			t.Fatalf("scanner ranking changed at %d: before=%s %.4f after=%s %.4f", i, before[i].Symbol, before[i].Score, after[i].Symbol, after[i].Score)
		}
	}

	universe := []string{"AAPL", "MSFT", "NVDA", "AMD"}
	sampleBefore, cursorBefore := radarSampleUniverse(universe, 1)
	sampleAfter, cursorAfter := radarSampleUniverse(universe, 1)
	if cursorBefore != cursorAfter || !reflect.DeepEqual(sampleBefore, sampleAfter) {
		t.Fatalf("radar sampling changed: before=%v/%d after=%v/%d", sampleBefore, cursorBefore, sampleAfter, cursorAfter)
	}
	now := time.Now()
	candidate := ScannerResult{Symbol: "AAPL", OpportunityScore: 90, Price: 100, DollarVolume: 20_000_000, SpreadPercent: .1, SessionRelativeVolume: 2, Reasons: []string{"test"}}
	promotions := selectOpportunityPromotions([]ScannerResult{candidate}, nil, now)
	if len(promotions) != 1 || promotions[0].Symbol != "AAPL" || promotions[0].State != "PROMOTED" {
		t.Fatalf("radar promotion behavior changed: %+v", promotions)
	}
}
