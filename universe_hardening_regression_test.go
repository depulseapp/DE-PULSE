package main

import (
	"context"
	"fmt"
	"io"
	"net/http"
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

type adaptProviderUniverseRoundTripFunc func(*http.Request) (*http.Response, error)

func (f adaptProviderUniverseRoundTripFunc) RoundTrip(r *http.Request) (*http.Response, error) {
	return f(r)
}

func adaptProviderUniverseClient(fn adaptProviderUniverseRoundTripFunc) *http.Client {
	return &http.Client{Transport: fn, Timeout: time.Second}
}

func adaptProviderUniverseResponse(status int, body string) *http.Response {
	return &http.Response{
		StatusCode: status,
		Status:     fmt.Sprintf("HTTP %d", status),
		Header:     make(http.Header),
		Body:       io.NopCloser(strings.NewReader(body)),
	}
}

func configureAdaptProviderUniverseAlpaca(e *Engine) {
	e.app.mu.Lock()
	e.app.secrets.AlpacaKey = "test-key"
	e.app.secrets.AlpacaSecret = "test-secret"
	e.app.mu.Unlock()
}

func adaptProviderUniverseHas(rows []string, want string) bool {
	for _, row := range rows {
		if row == want {
			return true
		}
	}
	return false
}

func TestADAPTProviderUniverseProductionUsesRouterV2(t *testing.T) {
	e := newV1801Engine(t)
	configureAdaptProviderUniverseAlpaca(e)
	calls := 0
	client := adaptProviderUniverseClient(func(r *http.Request) (*http.Response, error) {
		calls++
		return adaptProviderUniverseResponse(http.StatusOK, `[{"symbol":"ZZZZ","status":"active","tradable":true,"exchange":"NASDAQ"}]`), nil
	})

	e.mu.RLock()
	beforeDecisions := e.smartRouterScorecard.RouteDecisions
	e.mu.RUnlock()
	rows, ok := e.loadUSSymbolUniverseWithClient(context.Background(), "test-key", "test-secret", client)
	if !ok || !adaptProviderUniverseHas(rows, "ZZZZ") {
		t.Fatalf("production universe route failed: ok=%v rows=%v", ok, rows)
	}
	if calls != 1 {
		t.Fatalf("successful primary Alpaca asset endpoint should stop fallback, calls=%d", calls)
	}

	e.mu.RLock()
	global := e.providerCircuits[providerKey("Alpaca")]
	capability := e.providerCapabilityCircuits[providerCapabilityCircuitKey("Alpaca", canonicalUSAssetUniverseDataset)]
	afterDecisions := e.smartRouterScorecard.RouteDecisions
	e.mu.RUnlock()
	if global.LastSuccess == 0 || capability.LastSuccess == 0 {
		t.Fatalf("Router/provider success evidence missing: global=%+v capability=%+v", global, capability)
	}
	if afterDecisions != beforeDecisions+1 {
		t.Fatalf("production universe acquisition bypassed Router decision accounting: before=%d after=%d", beforeDecisions, afterDecisions)
	}
}

func TestADAPTProviderUniverseSyntheticLoaderRemainsRouterIsolated(t *testing.T) {
	e := newV1801Engine(t)
	e.canonicalUSUniverseLoader = func(context.Context, string, string) ([]string, bool) {
		return []string{"ZZZZ"}, true
	}
	e.mu.RLock()
	beforeDecisions := e.smartRouterScorecard.RouteDecisions
	e.mu.RUnlock()
	rows := e.canonicalUSSymbolUniverse(context.Background(), "", "", time.Now())
	if !adaptProviderUniverseHas(rows, "ZZZZ") {
		t.Fatalf("synthetic canonical loader did not remain isolated: %v", rows)
	}
	e.mu.RLock()
	afterDecisions := e.smartRouterScorecard.RouteDecisions
	e.mu.RUnlock()
	if afterDecisions != beforeDecisions {
		t.Fatalf("synthetic canonical loader unexpectedly entered Router: before=%d after=%d", beforeDecisions, afterDecisions)
	}
}

func TestADAPTProviderUniverseCancellationDoesNotPoisonRouterCircuit(t *testing.T) {
	e := newV1801Engine(t)
	configureAdaptProviderUniverseAlpaca(e)
	calls := 0
	client := adaptProviderUniverseClient(func(r *http.Request) (*http.Response, error) {
		calls++
		return nil, context.Canceled
	})
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	if rows, ok := e.loadUSSymbolUniverseWithClient(ctx, "test-key", "test-secret", client); ok || rows != nil {
		t.Fatalf("canceled production load unexpectedly succeeded: ok=%v rows=%v", ok, rows)
	}
	if calls > 1 {
		t.Fatalf("caller cancellation must stop same-provider fallback after first observed cancellation, calls=%d", calls)
	}

	e.mu.RLock()
	global := e.providerCircuits[providerKey("Alpaca")]
	capability := e.providerCapabilityCircuits[providerCapabilityCircuitKey("Alpaca", canonicalUSAssetUniverseDataset)]
	failureCount := int64(0)
	for _, state := range e.providerCapabilityStates {
		if state.Provider == "Alpaca" && state.Dataset == canonicalUSAssetUniverseDataset {
			failureCount += state.FailureCount
		}
	}
	e.mu.RUnlock()
	if global.Failures != 0 || global.LastFailure != 0 || capability.Failures != 0 || capability.LastFailure != 0 || failureCount != 0 {
		t.Fatalf("caller cancellation poisoned Router/provider health: global=%+v capability=%+v stateFailures=%d", global, capability, failureCount)
	}
}

func TestADAPTProviderUniverseSameProviderEndpointFallback(t *testing.T) {
	e := newV1801Engine(t)
	configureAdaptProviderUniverseAlpaca(e)
	calls := []string{}
	client := adaptProviderUniverseClient(func(r *http.Request) (*http.Response, error) {
		calls = append(calls, r.URL.Host)
		if len(calls) == 1 {
			return adaptProviderUniverseResponse(http.StatusInternalServerError, `{"message":"paper endpoint unavailable"}`), nil
		}
		return adaptProviderUniverseResponse(http.StatusOK, `[{"symbol":"ZZZZ","status":"active","tradable":true,"exchange":"NASDAQ"}]`), nil
	})
	rows, ok := e.loadUSSymbolUniverseWithClient(context.Background(), "test-key", "test-secret", client)
	if !ok || !adaptProviderUniverseHas(rows, "ZZZZ") {
		t.Fatalf("same-provider endpoint fallback failed: ok=%v rows=%v calls=%v", ok, rows, calls)
	}
	if len(calls) != 2 || calls[0] != "paper-api.alpaca.markets" || calls[1] != "api.alpaca.markets" {
		t.Fatalf("Alpaca endpoint fallback order changed: %v", calls)
	}
	e.mu.RLock()
	global := e.providerCircuits[providerKey("Alpaca")]
	capability := e.providerCapabilityCircuits[providerCapabilityCircuitKey("Alpaca", canonicalUSAssetUniverseDataset)]
	e.mu.RUnlock()
	if global.Failures != 0 || global.LastFailure != 0 || capability.Failures != 0 || capability.LastFailure != 0 {
		t.Fatalf("successful same-provider fallback must not count first endpoint as provider failure: global=%+v capability=%+v", global, capability)
	}
}

func TestADAPTProviderUniverseFailedProductionRefreshPreservesStaleTimestamp(t *testing.T) {
	e := newV1801Engine(t)
	configureAdaptProviderUniverseAlpaca(e)
	fail := false
	calls := 0
	client := adaptProviderUniverseClient(func(r *http.Request) (*http.Response, error) {
		calls++
		if fail {
			return adaptProviderUniverseResponse(http.StatusInternalServerError, `{"message":"assets unavailable"}`), nil
		}
		return adaptProviderUniverseResponse(http.StatusOK, `[{"symbol":"ZZZZ","status":"active","tradable":true,"exchange":"NASDAQ"}]`), nil
	})
	e.canonicalUSUniverseLoader = func(ctx context.Context, key, secret string) ([]string, bool) {
		return e.loadUSSymbolUniverseWithClient(ctx, key, secret, client)
	}

	firstNow := time.Now()
	first := e.canonicalUSSymbolUniverse(context.Background(), "test-key", "test-secret", firstNow)
	if !adaptProviderUniverseHas(first, "ZZZZ") {
		t.Fatalf("initial production-backed canonical load failed: %v", first)
	}
	e.mu.RLock()
	firstAt := e.canonicalUSUniverseAt
	e.mu.RUnlock()
	if firstAt != firstNow.UnixMilli() {
		t.Fatalf("unexpected initial canonical retrieval timestamp: got=%d want=%d", firstAt, firstNow.UnixMilli())
	}

	fail = true
	secondNow := firstNow.Add(opportunityUniverseTTL + time.Minute)
	second := e.canonicalUSSymbolUniverse(context.Background(), "test-key", "test-secret", secondNow)
	if !adaptProviderUniverseHas(second, "ZZZZ") {
		t.Fatalf("failed provider refresh did not preserve stale rows: %v", second)
	}
	e.mu.RLock()
	secondAt := e.canonicalUSUniverseAt
	source := e.canonicalUSUniverseSource
	retryAt := e.canonicalUSUniverseRetryAt
	e.mu.RUnlock()
	if secondAt != firstAt || source != "stale-cache" || retryAt <= secondNow.UnixMilli() {
		t.Fatalf("failed provider refresh corrupted stale-cache evidence time: firstAt=%d secondAt=%d source=%q retryAt=%d", firstAt, secondAt, source, retryAt)
	}
	if calls != 3 {
		t.Fatalf("expected one successful endpoint call then two terminal fallback calls, got %d", calls)
	}
}

func TestADAPTProviderUniverseFailureIsCapabilityScoped(t *testing.T) {
	e := newV1801Engine(t)
	configureAdaptProviderUniverseAlpaca(e)
	calls := 0
	client := adaptProviderUniverseClient(func(r *http.Request) (*http.Response, error) {
		calls++
		return adaptProviderUniverseResponse(http.StatusInternalServerError, `{"message":"assets unavailable"}`), nil
	})
	if rows, ok := e.loadUSSymbolUniverseWithClient(context.Background(), "test-key", "test-secret", client); ok || rows != nil {
		t.Fatalf("terminal provider failure unexpectedly succeeded: ok=%v rows=%v", ok, rows)
	}
	if calls != 2 {
		t.Fatalf("terminal provider failure must exhaust same-provider endpoints, calls=%d", calls)
	}

	e.mu.RLock()
	global := e.providerCircuits[providerKey("Alpaca")]
	universeCircuit := e.providerCapabilityCircuits[providerCapabilityCircuitKey("Alpaca", canonicalUSAssetUniverseDataset)]
	liveCircuit := e.providerCapabilityCircuits[providerCapabilityCircuitKey("Alpaca", "US Live Equities")]
	failureCount := int64(0)
	for _, state := range e.providerCapabilityStates {
		if state.Provider == "Alpaca" && state.Dataset == canonicalUSAssetUniverseDataset {
			failureCount += state.FailureCount
		}
	}
	e.mu.RUnlock()
	if global.Failures != 0 || global.LastFailure != 0 {
		t.Fatalf("universe endpoint failure leaked into global Alpaca circuit: %+v", global)
	}
	if universeCircuit.Failures != 1 || universeCircuit.LastFailure == 0 || failureCount != 1 {
		t.Fatalf("universe capability failure not recorded canonically: circuit=%+v stateFailures=%d", universeCircuit, failureCount)
	}
	if liveCircuit.Failures != 0 || liveCircuit.LastFailure != 0 || !e.providerAllowedFor("US Live Equities", "Alpaca") {
		t.Fatalf("universe failure suppressed unrelated Alpaca live-equities capability: %+v", liveCircuit)
	}
}

func TestADAPTProviderRequestLocalNeutralClassification(t *testing.T) {
	if !providerRequestFailureIsLocalNeutral(context.Background(), fmt.Errorf("BROAD_DISCOVERY deferred: provider request budget exhausted")) {
		t.Fatal("telemetry budget deferral must be local-neutral")
	}
	if !providerRequestFailureIsLocalNeutral(context.Background(), fmt.Errorf("provider request rejected by bounded workload budget")) {
		t.Fatal("workload budget rejection must be local-neutral")
	}
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	if !providerRequestFailureIsLocalNeutral(ctx, context.Canceled) {
		t.Fatal("caller cancellation must be local-neutral")
	}
	if providerRequestFailureIsLocalNeutral(context.Background(), fmt.Errorf("Alpaca HTTP 500")) {
		t.Fatal("genuine provider HTTP failure must not be local-neutral")
	}
}
