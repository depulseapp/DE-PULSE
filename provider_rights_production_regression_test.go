package main

import (
	"context"
	"encoding/json"
	"testing"
	"time"
)

type providerRightsMemoryBackend struct {
	quotes       map[string]Quote
	intelligence PersistenceIntelligenceBatch
}

func (b *providerRightsMemoryBackend) Name() string           { return "memory" }
func (b *providerRightsMemoryBackend) Capabilities() []string { return []string{"test"} }
func (b *providerRightsMemoryBackend) Init(context.Context) error { return nil }
func (b *providerRightsMemoryBackend) UpsertSymbols(context.Context, []SymbolRegistryRecord) (int, error) {
	return 0, nil
}
func (b *providerRightsMemoryBackend) LoadSymbols(context.Context) ([]SymbolRegistryRecord, error) {
	return nil, nil
}
func (b *providerRightsMemoryBackend) SaveQuotes(_ context.Context, quotes map[string]Quote) (int, error) {
	b.quotes = clone(quotes)
	return len(quotes), nil
}
func (b *providerRightsMemoryBackend) LoadQuotes(context.Context) (map[string]Quote, error) {
	return clone(b.quotes), nil
}
func (b *providerRightsMemoryBackend) SaveIntelligence(_ context.Context, batch PersistenceIntelligenceBatch) (int, error) {
	b.intelligence = batch
	return batch.Len(), nil
}
func (b *providerRightsMemoryBackend) LoadIdentityState(context.Context) (IdentityPersistentState, error) {
	return IdentityPersistentState{}, nil
}
func (b *providerRightsMemoryBackend) SaveIdentityState(context.Context, IdentityPersistentState) error {
	return nil
}
func (b *providerRightsMemoryBackend) LoadUserWorkspaces(context.Context) ([]UserWorkspace, error) {
	return nil, nil
}
func (b *providerRightsMemoryBackend) SaveUserWorkspace(context.Context, UserWorkspace) error {
	return nil
}
func (b *providerRightsMemoryBackend) Stats(context.Context) (PersistenceStoreStats, error) {
	return PersistenceStoreStats{CanonicalQuotes: len(b.quotes)}, nil
}
func (b *providerRightsMemoryBackend) Close() error { return nil }

func TestHOST003HostedCachePersistsAndReplaysOnlyCurrentlyApprovedProviderData(t *testing.T) {
	bindHostedRightsBundleForTest(t, approvedHostedRightsFixtureFor("Finnhub"))
	cache := MarketCache{
		Quotes: map[string]Quote{
			"SPY":  {Symbol: "SPY", Price: 600, Source: "finnhub-rest", UpdatedAt: time.Now().UnixMilli()},
			"QQQ":  {Symbol: "QQQ", Price: 500, Source: "alpaca-iex", UpdatedAt: time.Now().UnixMilli()},
		},
		History: map[string][]HistoryPoint{"SPY": {{T: time.Now().UnixMilli(), P: 600}}},
		News:    []NewsItem{{ID: "legacy", Headline: "unproven provider provenance"}},
		SavedAt: time.Now().UnixMilli(),
	}
	raw, err := json.Marshal(cache)
	if err != nil {
		t.Fatal(err)
	}
	var replay MarketCache
	if err := json.Unmarshal(raw, &replay); err != nil {
		t.Fatal(err)
	}
	if len(replay.Quotes) != 1 || replay.Quotes["SPY"].Source != "finnhub-rest" {
		t.Fatalf("hosted cache did not retain only approved provider quote: %+v", replay.Quotes)
	}
	if len(replay.History) != 0 || len(replay.News) != 0 {
		t.Fatalf("mixed legacy cache data without provider provenance survived hosted persistence: history=%d news=%d", len(replay.History), len(replay.News))
	}

	downgraded := approvedHostedRightsFixtureFor("Finnhub")
	downgraded.CachingRetention = providerRightsDenied
	bindHostedRightsBundleForTest(t, downgraded)
	var afterDowngrade MarketCache
	if err := json.Unmarshal(raw, &afterDowngrade); err != nil {
		t.Fatal(err)
	}
	if len(afterDowngrade.Quotes) != 0 {
		t.Fatalf("rights downgrade resurrected cached quote on restart: %+v", afterDowngrade.Quotes)
	}
}

func TestHOST003CanonicalPersistenceWriteAndReplayRecheckRights(t *testing.T) {
	bindHostedRightsBundleForTest(t, approvedHostedRightsFixtureFor("Finnhub"))
	inner := &providerRightsMemoryBackend{}
	backend := wrapHostedRightsPersistenceBackend(inner)
	quotes := map[string]Quote{
		"SPY": {Symbol: "SPY", Price: 600, Source: "finnhub-websocket", UpdatedAt: time.Now().UnixMilli()},
		"QQQ": {Symbol: "QQQ", Price: 500, Source: "alpaca-iex-websocket", UpdatedAt: time.Now().UnixMilli()},
	}
	if n, err := backend.SaveQuotes(context.Background(), quotes); err != nil || n != 1 {
		t.Fatalf("persistence write filter = n=%d err=%v", n, err)
	}
	if len(inner.quotes) != 1 || inner.quotes["SPY"].Source != "finnhub-websocket" {
		t.Fatalf("unapproved provider crossed persistence write boundary: %+v", inner.quotes)
	}

	downgraded := approvedHostedRightsFixtureFor("Finnhub")
	downgraded.Redistribution = providerRightsDenied
	bindHostedRightsBundleForTest(t, downgraded)
	loaded, err := backend.LoadQuotes(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if len(loaded) != 0 {
		t.Fatalf("persisted provider quote remained replayable after rights downgrade: %+v", loaded)
	}
}

func TestHOST003RecognizedExternalEvidenceCannotPersistWithoutRights(t *testing.T) {
	t.Setenv(runtimeModeEnv, "hosted")
	t.Setenv(hostedEnvironmentEnv, "prod")
	clearProviderRightsBundleForTest(t)
	inner := &providerRightsMemoryBackend{}
	backend := wrapHostedRightsPersistenceBackend(inner)
	batch := PersistenceIntelligenceBatch{
		Evidence: []EvidenceRecord{
			{ID: "external", Kind: "provider", Source: "Finnhub", ObservedAt: time.Now().UnixMilli()},
			{ID: "internal", Kind: "internal", Source: "canonical-signal-validation", ObservedAt: time.Now().UnixMilli()},
		},
	}
	if n, err := backend.SaveIntelligence(context.Background(), batch); err != nil || n != 1 {
		t.Fatalf("intelligence rights filter = n=%d err=%v", n, err)
	}
	if len(inner.intelligence.Evidence) != 1 || inner.intelligence.Evidence[0].ID != "internal" {
		t.Fatalf("unapproved external evidence persisted: %+v", inner.intelligence.Evidence)
	}
}

func TestHOST003LiveSubscriptionDesiredSetRevokesOnRightsDowngrade(t *testing.T) {
	e := newV1801Engine(t)
	bindHostedRightsBundleForTest(t, approvedHostedRightsFixtureFor("Finnhub"))
	if got := e.effectiveFinnhubSymbols(); len(got) == 0 {
		t.Fatal("approved Finnhub hosted rights did not produce a desired live subscription set")
	}
	if got := e.effectiveAlpacaIEXSymbols(); len(got) != 0 {
		t.Fatalf("Alpaca had live subscriptions without a reviewed rights record: %+v", got)
	}

	downgraded := approvedHostedRightsFixtureFor("Finnhub")
	downgraded.MultiUserUse = providerRightsDenied
	bindHostedRightsBundleForTest(t, downgraded)
	if got := e.effectiveFinnhubSymbols(); len(got) != 0 {
		t.Fatalf("Finnhub live subscription set survived rights downgrade: %+v", got)
	}
}

func TestHOST003UserServingPathFiltersRightsBeforePrivacyProjection(t *testing.T) {
	e := newV1801Engine(t)
	bindHostedRightsBundleForTest(t, approvedHostedRightsFixtureFor("Finnhub"))
	snap := RuntimeSnapshot{
		Status: "running",
		Quotes: map[string]Quote{
			"SPY": {Symbol: "SPY", Price: 600, Source: "finnhub-websocket", UpdatedAt: time.Now().UnixMilli()},
			"QQQ": {Symbol: "QQQ", Price: 500, Source: "alpaca-iex-websocket", UpdatedAt: time.Now().UnixMilli()},
		},
		History:  map[string][]HistoryPoint{"SPY": {{T: time.Now().UnixMilli(), P: 600}}},
		News:     []NewsItem{{ID: "n1", Headline: "publisher identity is not provider provenance", Symbols: []string{"SPY"}}},
		Earnings: []EarningsItem{{Symbol: "SPY", Date: "2026-08-26"}},
		Health:   map[string]string{},
	}
	out := e.app.runtimeSnapshotForUserFrom("", snap)
	if len(out.Quotes) != 1 || out.Quotes["SPY"].Price != 600 {
		t.Fatalf("approved provider quote was not served correctly: %+v", out.Quotes)
	}
	if len(out.History) != 0 || len(out.News) != 0 || len(out.Earnings) != 0 {
		t.Fatalf("source-unbound market data crossed hosted serving boundary: history=%d news=%d earnings=%d", len(out.History), len(out.News), len(out.Earnings))
	}

	downgraded := approvedHostedRightsFixtureFor("Finnhub")
	downgraded.Display = providerRightsDenied
	bindHostedRightsBundleForTest(t, downgraded)
	out = e.app.runtimeSnapshotForUserFrom("", snap)
	if len(out.Quotes) != 0 || out.Status != "degraded" || out.Health["provider-rights"] == "" {
		t.Fatalf("serving did not fail closed/diagnose rights downgrade: status=%s health=%q quotes=%+v", out.Status, out.Health["provider-rights"], out.Quotes)
	}
}
