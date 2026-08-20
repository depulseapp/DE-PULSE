package main

import (
	"context"
	"reflect"
	"strings"
	"sync/atomic"
	"testing"
	"time"
)

func TestV1881DiscoveryUniverseEligibilityIsExplicit(t *testing.T) {
	tests := []struct {
		name  string
		asset alpacaAsset
		want  bool
	}{
		{"active tradable Nasdaq equity", alpacaAsset{Symbol: "AAPL", Status: "active", Tradable: true, Exchange: "NASDAQ"}, true},
		{"active tradable NYSE equity", alpacaAsset{Symbol: "IBM", Status: "ACTIVE", Tradable: true, Exchange: "nyse"}, true},
		{"inactive excluded", alpacaAsset{Symbol: "AAPL", Status: "inactive", Tradable: true, Exchange: "NASDAQ"}, false},
		{"nontradable excluded", alpacaAsset{Symbol: "AAPL", Status: "active", Tradable: false, Exchange: "NASDAQ"}, false},
		{"unsupported exchange excluded", alpacaAsset{Symbol: "AAPL", Status: "active", Tradable: true, Exchange: "OTC"}, false},
		{"class-like symbol excluded", alpacaAsset{Symbol: "BRK.B", Status: "active", Tradable: true, Exchange: "NYSE"}, false},
		{"long symbol excluded", alpacaAsset{Symbol: "ABCDEF", Status: "active", Tradable: true, Exchange: "NASDAQ"}, false},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if got := canonicalUSUniverseAssetEligible(tc.asset); got != tc.want {
				t.Fatalf("eligibility=%v want=%v asset=%+v", got, tc.want, tc.asset)
			}
		})
	}
}

func TestV1881DiscoveryUniverseDoesNotSilentlyRequireOptions(t *testing.T) {
	urls := canonicalUSSymbolUniverseAssetURLs()
	if len(urls) != 2 {
		t.Fatalf("expected paper/live Alpaca asset endpoints, got %v", urls)
	}
	for _, raw := range urls {
		if strings.Contains(strings.ToLower(raw), "has_options") || strings.Contains(strings.ToLower(raw), "attributes=") {
			t.Fatalf("broad discovery universe must not silently require option capability: %s", raw)
		}
		if !strings.Contains(raw, canonicalUSSymbolUniverseProviderQuery) {
			t.Fatalf("missing explicit active U.S.-equity provider query: %s", raw)
		}
	}
}

func TestV1881UniverseRetrievalTimeIsNotEvidenceTime(t *testing.T) {
	e := newV188UniverseTestEngine(t)
	e.canonicalUSUniverseLoader = func(context.Context, string, string) ([]string, bool) {
		return []string{"AAPL", "MSFT"}, true
	}
	now := time.Date(2026, 8, 20, 14, 30, 0, 0, time.UTC)
	rows := e.canonicalUSSymbolUniverse(context.Background(), "key", "secret", now)
	if !reflect.DeepEqual(rows, []string{"AAPL", "MSFT"}) {
		t.Fatalf("unexpected universe: %v", rows)
	}
	timing := e.canonicalUSUniverseTimingSnapshot()
	if timing.RetrievedAtMS != now.UnixMilli() {
		t.Fatalf("retrieval time mismatch: got=%d want=%d", timing.RetrievedAtMS, now.UnixMilli())
	}
	if timing.EvidenceAtMS != 0 || timing.EvidenceTimeState != canonicalUSSymbolUniverseEvidenceTimeUnknown {
		t.Fatalf("provider evidence time must remain UNKNOWN, got %+v", timing)
	}
	e.mu.RLock()
	health := e.health[canonicalUSSymbolUniverseHealthKey]
	e.mu.RUnlock()
	if !strings.Contains(health, canonicalUSSymbolUniverseDiagnosticLabel) || !strings.Contains(health, "evidence time UNKNOWN") {
		t.Fatalf("neutral/unknown evidence-time diagnostic missing: %q", health)
	}
}

func TestV1881CanceledRefreshDoesNotPoisonRetryState(t *testing.T) {
	e := newV188UniverseTestEngine(t)
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	var canceledCalls int32
	e.canonicalUSUniverseLoader = func(ctx context.Context, _, _ string) ([]string, bool) {
		atomic.AddInt32(&canceledCalls, 1)
		<-ctx.Done()
		return nil, false
	}
	now := time.Date(2026, 8, 20, 14, 30, 0, 0, time.UTC)
	rows := e.canonicalUSSymbolUniverse(ctx, "key", "secret", now)
	if !reflect.DeepEqual(rows, uniqueSymbols(discoverySeedUniverse)) {
		t.Fatalf("canceled owner without cache should return deterministic seed, got %v", rows)
	}
	e.mu.RLock()
	retryAt := e.canonicalUSUniverseRetryAt
	refresh := e.canonicalUSUniverseRefresh
	e.mu.RUnlock()
	if retryAt != 0 {
		t.Fatalf("caller cancellation must not create provider retry suppression, retryAt=%d", retryAt)
	}
	if refresh != nil {
		t.Fatal("single-flight owner must be cleared after cancellation")
	}
	if atomic.LoadInt32(&canceledCalls) != 1 {
		t.Fatalf("expected one canceled loader call, got %d", canceledCalls)
	}

	var recoveryCalls int32
	e.canonicalUSUniverseLoader = func(context.Context, string, string) ([]string, bool) {
		atomic.AddInt32(&recoveryCalls, 1)
		return []string{"AAPL", "NVDA"}, true
	}
	recovered := e.canonicalUSSymbolUniverse(context.Background(), "key", "secret", now.Add(time.Second))
	if !reflect.DeepEqual(recovered, []string{"AAPL", "NVDA"}) || atomic.LoadInt32(&recoveryCalls) != 1 {
		t.Fatalf("live caller should immediately recover after canceled owner: rows=%v calls=%d", recovered, recoveryCalls)
	}
}

func TestV1881LoaderPanicReleasesSingleFlightOwner(t *testing.T) {
	e := newV188UniverseTestEngine(t)
	e.canonicalUSUniverseLoader = func(context.Context, string, string) ([]string, bool) {
		panic("synthetic universe loader panic")
	}
	now := time.Date(2026, 8, 20, 14, 30, 0, 0, time.UTC)
	func() {
		defer func() {
			if recover() == nil {
				t.Fatal("expected synthetic loader panic")
			}
		}()
		_ = e.canonicalUSSymbolUniverse(context.Background(), "key", "secret", now)
	}()

	e.mu.RLock()
	refresh := e.canonicalUSUniverseRefresh
	e.mu.RUnlock()
	if refresh != nil {
		t.Fatal("panic must not strand canonical universe single-flight owner")
	}

	var recoveryCalls int32
	e.canonicalUSUniverseLoader = func(context.Context, string, string) ([]string, bool) {
		atomic.AddInt32(&recoveryCalls, 1)
		return []string{"AAPL", "MSFT", "NVDA"}, true
	}
	rows := e.canonicalUSSymbolUniverse(context.Background(), "key", "secret", now.Add(time.Second))
	if !reflect.DeepEqual(rows, []string{"AAPL", "MSFT", "NVDA"}) || atomic.LoadInt32(&recoveryCalls) != 1 {
		t.Fatalf("universe must recover immediately after panic cleanup: rows=%v calls=%d", rows, recoveryCalls)
	}
}
