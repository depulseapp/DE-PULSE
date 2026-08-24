package main

import (
	"context"
	"net/http"
	"reflect"
	"testing"
	"time"
)

func TestADAPTInstrumentIdentityCaptureUsesExistingRoutedUniverseFetch(t *testing.T) {
	e := newV1801Engine(t)
	configureAdaptProviderUniverseAlpaca(e)
	calls := 0
	client := adaptProviderUniverseClient(func(r *http.Request) (*http.Response, error) {
		calls++
		return adaptProviderUniverseResponse(http.StatusOK, `[
{"id":"asset-aapl","class":"us_equity","exchange":"NASDAQ","symbol":"AAPL","name":"Apple Inc.","status":"active","tradable":true},
{"id":"asset-otc","class":"us_equity","exchange":"OTC","symbol":"OTCX","name":"OTC Example","status":"active","tradable":true},
{"id":"asset-class","class":"us_equity","exchange":"NYSE","symbol":"BRK.B","name":"Berkshire Hathaway Class B","status":"active","tradable":true}
]`), nil
	})

	rows, ok := e.loadUSSymbolUniverseWithClient(context.Background(), "test-key", "test-secret", client)
	if !ok || !adaptProviderUniverseHas(rows, "AAPL") {
		t.Fatalf("routed universe acquisition failed: ok=%v rows=%v", ok, rows)
	}
	if calls != 1 {
		t.Fatalf("identity capture must not add provider requests; got %d HTTP calls", calls)
	}
	identity, found := e.instrumentIdentityForSymbol("aapl")
	if !found {
		t.Fatal("expected canonical AAPL identity from the same /v2/assets response")
	}
	if identity.Name != "Apple Inc." || identity.Exchange != "NASDAQ" || identity.AssetClass != "us_equity" || identity.ProviderAssetID != "asset-aapl" || identity.Source != "alpaca-assets" {
		t.Fatalf("unexpected canonical identity: %+v", identity)
	}
	if _, found := e.instrumentIdentityForSymbol("OTCX"); found {
		t.Fatal("unsupported-exchange asset must not enter canonical identity")
	}
	if _, found := e.instrumentIdentityForSymbol("BRK.B"); found {
		t.Fatal("symbol outside canonical U.S. universe boundary must not enter identity")
	}
}

func TestADAPTInstrumentIdentityPersistenceDoesNotMutateTradingRegistry(t *testing.T) {
	dir := t.TempDir()
	p := NewPersistenceManager(dir)
	now := time.Now().UnixMilli()
	registry := []SymbolRegistryRecord{
		{Symbol: "AAPL", FirstSeenAt: now - 1000, LastSeenAt: now, Active: true, Selected: true, ProcessingTier: 1, ProviderEligible: true},
		{Symbol: "MSFT", FirstSeenAt: now - 1000, LastSeenAt: now, Active: true, Selected: false, ProcessingTier: 2, ProviderEligible: true},
	}
	p.EnqueueSymbols(registry)
	p.flushPending()
	before := p.LoadSymbols()

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	written, err := p.SaveInstrumentIdentities(ctx, []InstrumentIdentityRecord{{
		Symbol: "AAPL", Name: "Apple Inc.", Exchange: "NASDAQ", AssetClass: "us_equity", ProviderAssetID: "asset-aapl", Source: "alpaca-assets", ObservedAt: now,
	}})
	cancel()
	if err != nil || written != 1 {
		_ = p.Close()
		t.Fatalf("identity persistence failed: written=%d err=%v", written, err)
	}
	after := p.LoadSymbols()
	if !reflect.DeepEqual(before, after) {
		_ = p.Close()
		t.Fatalf("partial identity persistence mutated trading registry:\nbefore=%+v\nafter=%+v", before, after)
	}
	if err := p.Close(); err != nil {
		t.Fatalf("close persistence: %v", err)
	}
}

func TestADAPTInstrumentIdentityWarmReuseAcrossRestart(t *testing.T) {
	dir := t.TempDir()
	p := NewPersistenceManager(dir)
	now := time.Now().UnixMilli()
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	_, err := p.SaveInstrumentIdentities(ctx, []InstrumentIdentityRecord{{
		Symbol: "NVDA", Name: "NVIDIA Corporation", Exchange: "NASDAQ", AssetClass: "us_equity", ProviderAssetID: "asset-nvda", Source: "alpaca-assets", ObservedAt: now,
	}})
	cancel()
	if err != nil {
		_ = p.Close()
		t.Fatalf("seed identity persistence: %v", err)
	}
	if err := p.Close(); err != nil {
		t.Fatalf("close seed persistence: %v", err)
	}

	p2 := NewPersistenceManager(dir)
	defer p2.Close()
	e := &Engine{app: &Application{persistence: p2}, health: map[string]string{}}
	identity, found := e.instrumentIdentityForSymbol("nvda")
	if !found {
		t.Fatal("persisted identity was not warm-reused after restart")
	}
	if identity.Name != "NVIDIA Corporation" || identity.ProviderAssetID != "asset-nvda" || identity.ObservedAt != now {
		t.Fatalf("warm identity changed across restart: %+v", identity)
	}
}

func TestADAPTInstrumentIdentityRejectsOlderOverwrite(t *testing.T) {
	dir := t.TempDir()
	p := NewPersistenceManager(dir)
	defer p.Close()
	newerAt := time.Now().UnixMilli()
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	_, err := p.SaveInstrumentIdentities(ctx, []InstrumentIdentityRecord{{
		Symbol: "AAPL", Name: "Apple Inc.", Exchange: "NASDAQ", AssetClass: "us_equity", Source: "alpaca-assets", ObservedAt: newerAt,
	}})
	cancel()
	if err != nil {
		t.Fatalf("save newer identity: %v", err)
	}
	ctx, cancel = context.WithTimeout(context.Background(), 3*time.Second)
	_, err = p.SaveInstrumentIdentities(ctx, []InstrumentIdentityRecord{{
		Symbol: "AAPL", Name: "Stale Name", Exchange: "NYSE", AssetClass: "us_equity", Source: "alpaca-assets", ObservedAt: newerAt - 10_000,
	}})
	cancel()
	if err != nil {
		t.Fatalf("save older identity: %v", err)
	}
	ctx, cancel = context.WithTimeout(context.Background(), 3*time.Second)
	rows, err := p.LoadInstrumentIdentities(ctx)
	cancel()
	if err != nil {
		t.Fatalf("load identity: %v", err)
	}
	if len(rows) != 1 || rows[0].Name != "Apple Inc." || rows[0].Exchange != "NASDAQ" || rows[0].ObservedAt != newerAt {
		t.Fatalf("older observation overwrote newer canonical identity: %+v", rows)
	}
}
