package main

import (
	"context"
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"sync"
	"testing"
	"time"
)

func TestV170PersistentCanonicalQuoteWarmStartIsNeverTimelessLiveTruth(t *testing.T) {
	dir := t.TempDir()
	p := NewPersistenceManager(dir)
	if d := p.Diagnostics(); !d.Ready {
		t.Fatalf("persistence backend not ready: %+v", d)
	}
	original := Quote{Symbol: "NVDA", Price: 200.5, Bid: 200.4, Ask: 200.6, Source: "alpaca-iex", FeedType: "websocket", DataState: "live", ProviderTimestamp: 1_900, UpdatedAt: 2_000}
	p.EnqueueQuotes(map[string]Quote{"NVDA": original})
	if err := p.Close(); err != nil {
		t.Fatal(err)
	}

	reopened := NewPersistenceManager(dir)
	defer reopened.Close()
	loaded := reopened.LoadQuotes()
	q, ok := loaded["NVDA"]
	if !ok {
		t.Fatalf("persisted NVDA quote missing: %#v", loaded)
	}
	if q.Price != original.Price || q.Source != original.Source {
		t.Fatalf("persisted quote changed: got=%+v want=%+v", q, original)
	}
	if q.DataState != "persisted" || q.FeedType != "persisted" {
		t.Fatalf("warm-start quote retained timeless live labels: %+v", q)
	}
}

func TestV170EnginePrefersNewerPersistentCanonicalQuoteOverOlderJSONCache(t *testing.T) {
	dir := t.TempDir()
	p := NewPersistenceManager(dir)
	p.EnqueueQuotes(map[string]Quote{"NVDA": {Symbol: "NVDA", Price: 222, Source: "alpaca-iex", DataState: "live", ProviderTimestamp: 20_000, UpdatedAt: 20_100}})
	if err := p.Close(); err != nil {
		t.Fatal(err)
	}
	cache := MarketCache{Quotes: map[string]Quote{"NVDA": {Symbol: "NVDA", Price: 111, Source: "old-cache", ProviderTimestamp: 10_000, UpdatedAt: 10_100}}, SavedAt: 10_200}
	raw, _ := json.Marshal(cache)
	if err := os.WriteFile(filepath.Join(dir, "market-cache.json"), raw, 0600); err != nil {
		t.Fatal(err)
	}

	app := &Application{configDir: dir, hub: NewHub(), state: defaultState(), sessionKey: "v17-test", aiCache: map[string]aiCacheEntry{}}
	app.persistence = NewPersistenceManager(dir)
	defer app.persistence.Close()
	app.engine = NewEngine(app)
	app.engine.mu.RLock()
	q := app.engine.quotes["NVDA"]
	app.engine.mu.RUnlock()
	if q.Price != 222 || q.Source != "alpaca-iex" || q.DataState != "persisted" {
		t.Fatalf("engine did not use newer persisted canonical state: %+v", q)
	}
}

func TestV170SymbolRegistryPriorityContract(t *testing.T) {
	st := defaultState()
	st.UI.SelectedTicker = "NVDA"
	records := symbolRegistryRecords(st, time.Unix(100, 0))
	bySymbol := map[string]SymbolRegistryRecord{}
	for _, r := range records {
		bySymbol[r.Symbol] = r
	}
	for _, sym := range []string{"SPY", "QQQ", "NVDA"} {
		if r, ok := bySymbol[sym]; !ok || r.ProcessingTier != 0 {
			t.Fatalf("Tier 0 contract missing for %s: %+v", sym, r)
		}
	}
	for _, sym := range []string{"GLD", "SLV", "USO", "TSLA", "PLTR"} {
		if r, ok := bySymbol[sym]; !ok || r.ProcessingTier > 1 {
			t.Fatalf("Tier 1 actionable contract missing for %s: %+v", sym, r)
		}
	}
}

func TestV170SharedProviderBudgetIsBoundedAndReportsQueue(t *testing.T) {
	w := NewWorkloadController()
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	var initial WorkClassDiagnostics
	for _, d := range w.Diagnostics() {
		if d.Class == "provider-rest" {
			initial = d
		}
	}
	if initial.Capacity < 1 {
		t.Fatalf("provider budget missing: %+v", initial)
	}
	releases := make([]func(), 0, initial.Capacity)
	for i := 0; i < initial.Capacity; i++ {
		release, ok := w.Acquire(ctx, "provider-rest")
		if !ok {
			t.Fatal("could not acquire expected provider slot")
		}
		releases = append(releases, release)
	}
	acquired := make(chan func(), 1)
	go func() {
		release, ok := w.Acquire(ctx, "provider-rest")
		if ok {
			acquired <- release
		}
	}()
	time.Sleep(20 * time.Millisecond)
	var provider WorkClassDiagnostics
	for _, d := range w.Diagnostics() {
		if d.Class == "provider-rest" {
			provider = d
		}
	}
	if provider.InFlight != provider.Capacity || provider.Queued < 1 {
		t.Fatalf("provider budget not bounded/observable: %+v", provider)
	}
	releases[0]()
	select {
	case next := <-acquired:
		next()
	case <-ctx.Done():
		t.Fatal("queued provider work did not proceed after capacity freed")
	}
	for _, release := range releases[1:] {
		release()
	}
}

func TestV170RuntimeLoadDiagnosticsExposePersistenceAndWorkload(t *testing.T) {
	dir := t.TempDir()
	app := &Application{configDir: dir, hub: NewHub(), state: defaultState(), sessionKey: "v17-load", aiCache: map[string]aiCacheEntry{}}
	app.persistence = NewPersistenceManager(dir)
	defer app.persistence.Close()
	app.engine = NewEngine(app)
	snap := app.engine.Snapshot()
	if snap.RuntimeLoad.SampledAt == 0 || snap.RuntimeLoad.Goroutines <= 0 {
		t.Fatalf("runtime load profile missing: %+v", snap.RuntimeLoad)
	}
	if snap.RuntimeLoad.Persistence.Backend == "" {
		t.Fatalf("persistence diagnostics missing: %+v", snap.RuntimeLoad.Persistence)
	}
	if len(snap.RuntimeLoad.Workload) == 0 {
		t.Fatalf("workload diagnostics missing: %+v", snap.RuntimeLoad.Workload)
	}
}

func TestV170MaterialChangePersistenceSuppressesTickNoise(t *testing.T) {
	dir := t.TempDir()
	p := NewPersistenceManager(dir)
	base := Quote{Symbol: "NVDA", Price: 200, Bid: 199.99, Ask: 200.01, Source: "alpaca-iex", FeedType: "websocket", DataState: "live", ProviderTimestamp: 100_000, UpdatedAt: 100_100}
	p.EnqueueQuotes(map[string]Quote{"NVDA": base})
	noise := base
	noise.Price = 200.001
	noise.ProviderTimestamp += 1_000
	noise.UpdatedAt += 1_000
	p.EnqueueQuotes(map[string]Quote{"NVDA": noise})
	if got := p.Diagnostics().MaterialWritesSuppressed; got < 1 {
		t.Fatalf("expected immaterial tick write to be suppressed, diagnostics=%+v", p.Diagnostics())
	}
	material := noise
	material.Price = 201
	material.ProviderTimestamp += 2_000
	material.UpdatedAt += 2_000
	p.EnqueueQuotes(map[string]Quote{"NVDA": material})
	if err := p.Close(); err != nil {
		t.Fatal(err)
	}
	if d := p.Diagnostics(); d.Errors != 0 || d.RowsWritten < 1 {
		t.Fatalf("material persistence did not flush cleanly: %+v", d)
	}
}

func TestV170PersistentIntelligenceRepositoryFoundation(t *testing.T) {
	dir := t.TempDir()
	p := NewPersistenceManager(dir)
	caps := map[string]bool{}
	for _, c := range p.Diagnostics().Capabilities {
		caps[c] = true
	}
	for _, want := range []string{"global-symbol-registry", "canonical-quotes", "evidence-records", "decision-lineage", "outcome-history", "derived-feature-store"} {
		if !caps[want] {
			t.Fatalf("missing persistence capability %q: %+v", want, p.Diagnostics().Capabilities)
		}
	}
	now := time.Now().UnixMilli()
	p.EnqueueIntelligence(PersistenceIntelligenceBatch{
		Evidence:  []EvidenceRecord{{ID: "ev-1", Symbol: "NVDA", Kind: "research", ObservedAt: now, Source: "canonical", Provenance: "test", FreshnessState: "CURRENT", Payload: json.RawMessage(`{"price":200}`)}},
		Decisions: []DecisionLineageRecord{{ID: "dec-1", Symbol: "NVDA", Horizon: "day", EvidenceID: "ev-1", DecisionKind: "readiness", DecisionValue: "REVIEW", FormulaVersion: "protected", CreatedAt: now, Payload: json.RawMessage(`{"reason":"test"}`)}},
		Outcomes:  []OutcomeHistoryRecord{{ID: "out-1", DecisionID: "dec-1", Symbol: "NVDA", Horizon: "day", ObservedAt: now, OutcomeLabel: "PENDING", Payload: json.RawMessage(`{"window":"open"}`)}},
		Features:  []DerivedFeatureRecord{{Symbol: "NVDA", FeatureKey: "example", FeatureVersion: "v1", AsOf: now, SourceHash: "abc", Payload: json.RawMessage(`{"value":1}`)}},
	})
	if err := p.Close(); err != nil {
		t.Fatal(err)
	}
	if d := p.Diagnostics(); d.Errors != 0 || d.RowsWritten < 4 {
		t.Fatalf("structured persistence batch failed: %+v", d)
	}
}

func TestV170ProviderTelemetryAndAdaptiveDegradationReasonCodes(t *testing.T) {
	pt := NewProviderTelemetry()
	done := pt.begin("Finnhub")
	done(context.DeadlineExceeded)
	done429 := pt.begin("Alpaca")
	done429(&testHTTP429Error{})
	rows := pt.Diagnostics()
	if len(rows) != 2 {
		t.Fatalf("provider telemetry missing rows: %+v", rows)
	}
	var alpaca ProviderRequestDiagnostics
	for _, row := range rows {
		if row.Provider == "Alpaca" {
			alpaca = row
		}
	}
	if alpaca.Requests != 1 || alpaca.Errors != 1 || alpaca.RateLimited != 1 {
		t.Fatalf("rate-limit telemetry incorrect: %+v", alpaca)
	}
	state := deriveRuntimeDegradation("running", "live", FeedDiagnostics{MarketSession: "closed"}, []FreshnessDiagnostic{{Dataset: "Quotes", State: "FRESH"}, {Dataset: "VIX", State: "FRESH"}}, ProviderRouterSnapshot{}, RuntimeLoadDiagnostics{ProviderRequests: rows})
	if state.Code != "RATE LIMITED" || !state.CriticalUsable {
		t.Fatalf("adaptive degradation reason incorrect: %+v", state)
	}
}

type testHTTP429Error struct{}

func (*testHTTP429Error) Error() string { return "HTTP 429: too many requests" }

func TestV170WorkloadPressureTriggersLoadSheddingSignal(t *testing.T) {
	w := NewWorkloadController()
	releases := []func(){}
	for i := 0; i < 4; i++ {
		release, ok := w.TryAcquireTier("provider-rest", WorkTierUserActionable)
		if !ok {
			t.Fatal("expected provider slot")
		}
		releases = append(releases, release)
	}
	if !w.Pressured() {
		t.Fatal("full provider budget should report pressure for load shedding")
	}
	for _, release := range releases {
		release()
	}
	if w.Pressured() {
		t.Fatal("pressure should clear after provider capacity is released")
	}
}

func TestV170GlobalSymbolRegistryRetainsInactiveHistory(t *testing.T) {
	dir := t.TempDir()
	p := NewPersistenceManager(dir)
	if d := p.Diagnostics(); !d.Ready {
		t.Fatalf("persistence backend not ready: %+v", d)
	}
	first := int64(1_000)
	p.EnqueueSymbols([]SymbolRegistryRecord{{Symbol: "NVDA", FirstSeenAt: first, LastSeenAt: first, Active: true, ProcessingTier: 1, ProviderEligible: true}})
	p.EnqueueSymbols([]SymbolRegistryRecord{{Symbol: "SPY", FirstSeenAt: 2_000, LastSeenAt: 2_000, Active: true, ProcessingTier: 0, ProviderEligible: true}})
	if err := p.Close(); err != nil {
		t.Fatal(err)
	}

	reopened := NewPersistenceManager(dir)
	defer reopened.Close()
	records := reopened.LoadSymbols()
	bySymbol := map[string]SymbolRegistryRecord{}
	for _, r := range records {
		bySymbol[r.Symbol] = r
	}
	nvda, ok := bySymbol["NVDA"]
	if !ok {
		t.Fatalf("removed desk symbol was deleted from canonical registry: %+v", records)
	}
	if nvda.Active {
		t.Fatalf("removed desk symbol should be retained inactive: %+v", nvda)
	}
	if nvda.FirstSeenAt != first {
		t.Fatalf("first-seen history changed: got %d want %d", nvda.FirstSeenAt, first)
	}
	if spy, ok := bySymbol["SPY"]; !ok || !spy.Active || spy.ProcessingTier != 0 {
		t.Fatalf("active registry symbol incorrect: %+v", spy)
	}
}

func TestV170PersistenceSchemaVersionAndStoreStats(t *testing.T) {
	dir := t.TempDir()
	p := NewPersistenceManager(dir)
	p.EnqueueSymbols([]SymbolRegistryRecord{{Symbol: "SPY", FirstSeenAt: 1_000, LastSeenAt: 1_000, Active: true, ProcessingTier: 0, ProviderEligible: true}})
	p.EnqueueQuotes(map[string]Quote{"SPY": {Symbol: "SPY", Price: 600, Source: "test", ProviderTimestamp: 2_000, UpdatedAt: 2_010}})
	if err := p.Close(); err != nil {
		t.Fatal(err)
	}

	reopened := NewPersistenceManager(dir)
	defer reopened.Close()
	d := reopened.Diagnostics()
	if d.Store.SchemaVersion < 2 {
		t.Fatalf("schema migrations were not versioned/applied: %+v", d.Store)
	}
	if d.Store.SymbolCount < 1 || d.Store.ActiveSymbolCount < 1 || d.Store.CanonicalQuotes < 1 {
		t.Fatalf("store statistics incomplete: %+v", d.Store)
	}
	if d.Store.StorageBytes <= 0 {
		t.Fatalf("storage growth telemetry missing: %+v", d.Store)
	}
}

func TestV170SignalValidationFeedsEvidenceDecisionOutcomeAndFeatureStore(t *testing.T) {
	dir := t.TempDir()
	app := &Application{configDir: dir, hub: NewHub(), state: defaultState(), sessionKey: "v17-lineage", aiCache: map[string]aiCacheEntry{}}
	app.persistence = NewPersistenceManager(dir)
	app.engine = NewEngine(app)
	now := time.Now().UnixMilli()
	app.engine.recordSignalSnapshot(SignalSnapshot{
		ID: "sig-v17-1", Symbol: "NVDA", Horizon: "day", Timestamp: now, Price: 200, Score: 81,
		Action: "WATCH", Readiness: "REVIEW", EvidenceSnapshotID: "evidence-v17-1", FormulaVersion: "protected-v14.3.7",
		FamilyScores: map[string]float64{"trend": 82, "momentum": 77}, MarketRegime: "Neutral", MarketTradeability: "NORMAL",
		OutcomeState: "PENDING", OutcomeUpdatedAt: now,
	})
	if err := app.persistence.Close(); err != nil {
		t.Fatal(err)
	}
	d := app.persistence.Diagnostics()
	if d.Errors != 0 {
		t.Fatalf("lineage persistence errors: %+v", d)
	}
	if d.Store.EvidenceRows < 1 || d.Store.DecisionRows < 1 || d.Store.OutcomeRows < 1 || d.Store.FeatureRows < 1 {
		t.Fatalf("signal validation did not feed the complete persistent intelligence path: %+v", d.Store)
	}
}

func TestV170RegistryDoesNotFabricateSubscriptionOrProcessingTimestamps(t *testing.T) {
	st := defaultState()
	records := symbolRegistryRecords(st, time.Unix(123, 0))
	if len(records) == 0 {
		t.Fatal("expected canonical symbol registry records")
	}
	for _, r := range records {
		if r.LastSubscribedAt != 0 || r.LastProcessedAt != 0 {
			t.Fatalf("membership snapshot fabricated runtime timestamps for %s: %+v", r.Symbol, r)
		}
	}
}

type flakyPersistenceBackend struct {
	mu        sync.Mutex
	failQuote int
	quotes    map[string]Quote
}

func (b *flakyPersistenceBackend) Name() string               { return "flaky-test" }
func (b *flakyPersistenceBackend) Capabilities() []string     { return []string{"test"} }
func (b *flakyPersistenceBackend) Init(context.Context) error { return nil }
func (b *flakyPersistenceBackend) UpsertSymbols(context.Context, []SymbolRegistryRecord) (int, error) {
	return 0, nil
}
func (b *flakyPersistenceBackend) LoadSymbols(context.Context) ([]SymbolRegistryRecord, error) {
	return nil, nil
}
func (b *flakyPersistenceBackend) LoadQuotes(context.Context) (map[string]Quote, error) {
	b.mu.Lock()
	defer b.mu.Unlock()
	out := map[string]Quote{}
	for k, v := range b.quotes {
		out[k] = v
	}
	return out, nil
}
func (b *flakyPersistenceBackend) SaveQuotes(_ context.Context, quotes map[string]Quote) (int, error) {
	b.mu.Lock()
	defer b.mu.Unlock()
	if b.failQuote > 0 {
		b.failQuote--
		return 0, errors.New("transient persistence failure")
	}
	if b.quotes == nil {
		b.quotes = map[string]Quote{}
	}
	for k, v := range quotes {
		b.quotes[k] = v
	}
	return len(quotes), nil
}
func (b *flakyPersistenceBackend) SaveIntelligence(context.Context, PersistenceIntelligenceBatch) (int, error) {
	return 0, nil
}
func (b *flakyPersistenceBackend) LoadIdentityState(context.Context) (IdentityPersistentState, error) {
	return IdentityPersistentState{}, nil
}
func (b *flakyPersistenceBackend) SaveIdentityState(context.Context, IdentityPersistentState) error {
	return nil
}
func (b *flakyPersistenceBackend) LoadUserWorkspaces(context.Context) ([]UserWorkspace, error) {
	return nil, nil
}
func (b *flakyPersistenceBackend) SaveUserWorkspace(context.Context, UserWorkspace) error {
	return nil
}
func (b *flakyPersistenceBackend) Stats(context.Context) (PersistenceStoreStats, error) {
	return PersistenceStoreStats{}, nil
}
func (b *flakyPersistenceBackend) Close() error { return nil }

func newFlakyPersistenceManagerForTest(b PersistenceBackend) *PersistenceManager {
	p := &PersistenceManager{
		backend: b, wake: make(chan struct{}, 1), stop: make(chan struct{}), done: make(chan struct{}),
		pendingQuotes: map[string]Quote{}, lastAcceptedQuotes: map[string]Quote{},
	}
	p.diag.Backend = b.Name()
	p.diag.Capabilities = b.Capabilities()
	if err := b.Init(context.Background()); err != nil {
		p.diag.Errors++
		p.diag.LastError = err.Error()
		close(p.done)
		return p
	}
	p.diag.Ready = true
	go p.worker()
	return p
}

func TestV170AsyncPersistenceRetriesTransientFailureWithoutDroppingQuote(t *testing.T) {
	b := &flakyPersistenceBackend{failQuote: 1, quotes: map[string]Quote{}}
	p := newFlakyPersistenceManagerForTest(b)
	p.EnqueueQuotes(map[string]Quote{"NVDA": {Symbol: "NVDA", Price: 201, Source: "test", UpdatedAt: time.Now().UnixMilli()}})
	deadline := time.Now().Add(2 * time.Second)
	for time.Now().Before(deadline) {
		d := p.Diagnostics()
		b.mu.Lock()
		_, persisted := b.quotes["NVDA"]
		b.mu.Unlock()
		if persisted && d.RetryBatches >= 1 && d.WriteBatches >= 1 {
			if d.DroppedBatches != 0 {
				t.Fatalf("recovered transient batch was reported dropped: %+v", d)
			}
			if err := p.Close(); err != nil {
				t.Fatal(err)
			}
			return
		}
		time.Sleep(20 * time.Millisecond)
	}
	d := p.Diagnostics()
	_ = p.Close()
	t.Fatalf("transient quote write was not retried successfully: %+v", d)
}

func TestV171CanonicalSymbolWorkTierContract(t *testing.T) {
	a := newTestApplication(t)
	e := a.engine
	a.mu.Lock()
	a.state.UI.SelectedTicker = "NVDA"
	for i := range a.state.Watchlists {
		if a.state.Watchlists[i].ID == a.state.Settings.SwingWatchlistID {
			a.state.Watchlists[i].Symbols = append(a.state.Watchlists[i].Symbols, "MSFT")
		}
		if a.state.Watchlists[i].ID == a.state.Settings.DiscoveryWatchlistID {
			a.state.Watchlists[i].Symbols = append(a.state.Watchlists[i].Symbols, "ZZQX")
		}
	}
	a.mu.Unlock()
	if got := e.workTierForSymbol("SPY"); got != WorkTierMarketCritical {
		t.Fatalf("SPY tier=%d", got)
	}
	if got := e.workTierForSymbol("NVDA"); got != WorkTierMarketCritical {
		t.Fatalf("selected tier=%d", got)
	}
	if got := e.workTierForSymbol("GLD"); got != WorkTierUserActionable {
		t.Fatalf("GLD tier=%d", got)
	}
	if got := e.workTierForSymbol("MSFT"); got != WorkTierUserActionable {
		t.Fatalf("swing tier=%d", got)
	}
	if got := e.workTierForSymbol("ZZQX"); got != WorkTierRadarPromoted {
		t.Fatalf("discovery tier=%d", got)
	}
}
