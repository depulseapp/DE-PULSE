package main

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"
	"testing"
	"time"
)

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
	client := adaptProviderUniverseClient(func(r *http.Request) (*http.Response, error) {
		t.Fatal("canceled provider request should not reach transport")
		return nil, context.Canceled
	})
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	if rows, ok := e.loadUSSymbolUniverseWithClient(ctx, "test-key", "test-secret", client); ok || rows != nil {
		t.Fatalf("canceled production load unexpectedly succeeded: ok=%v rows=%v", ok, rows)
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
