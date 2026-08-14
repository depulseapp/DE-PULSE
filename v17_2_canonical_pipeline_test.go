package main

import (
	"encoding/json"
	"testing"
	"time"
)

func newV172PipelineTestEngine(t *testing.T) (*Engine, *PersistenceManager) {
	t.Helper()
	dir := t.TempDir()
	p := NewPersistenceManager(dir)
	if d := p.Diagnostics(); !d.Ready {
		t.Fatalf("persistence unavailable: %+v", d)
	}
	app := &Application{configDir: dir, hub: NewHub(), persistence: p, state: defaultState(), aiCache: map[string]aiCacheEntry{}}
	e := &Engine{app: app, quotes: map[string]Quote{}, history: map[string][]HistoryPoint{}, bars: map[string]map[string][]Bar{}, lastUpdated: map[string]int64{}, lastBroadcast: map[string]time.Time{}, liquidityBaselines: map[string]LiquidityBaseline{}, providerQuotes: map[string]map[string]Quote{}, catalystReactions: map[string]CatalystReactionState{}, preparations: initialPreparationJobs(time.Now()), health: map[string]string{}}
	return e, p
}

func TestV172ImmaterialTicksUpdateMemoryButSuppressHeavyDownstream(t *testing.T) {
	e, p := newV172PipelineTestEngine(t)
	defer p.Close()
	base := time.Now().UnixMilli()
	e.quotes["NVDA"] = Quote{Symbol: "NVDA", Price: 100, Bid: 99.99, Ask: 100.01, Source: "test-live", FeedType: "websocket", DataState: "live", ProviderTimestamp: base, UpdatedAt: base}
	ch := e.app.hub.Subscribe()
	defer e.app.hub.Unsubscribe(ch)
	for i := 1; i <= 100; i++ {
		price := 100 + float64(i)*.001
		e.updateQuote("NVDA", Quote{Price: price, Bid: price - .01, Ask: price + .01, ProviderTimestamp: base + int64(i*100), FeedType: "websocket", DataState: "live"}, "test-live")
	}
	e.mu.RLock()
	q := e.quotes["NVDA"]
	hl := len(e.history["NVDA"])
	d := e.canonicalPipeline
	e.mu.RUnlock()
	if q.Price <= 100 || hl == 0 {
		t.Fatalf("memory truth did not update: q=%+v history=%d", q, hl)
	}
	if d.ReceivedQuoteChanges != 100 || d.MaterialQuoteChanges != 0 || d.SuppressedDownstream != 100 || d.PersistenceEnqueues != 0 {
		t.Fatalf("immaterial work not suppressed: %+v", d)
	}
	select {
	case <-ch:
	default:
		t.Fatal("lightweight UI quote update was suppressed")
	}
}

func TestV172MaterialPriceAndTruthStateChangesPropagate(t *testing.T) {
	e, p := newV172PipelineTestEngine(t)
	defer p.Close()
	base := time.Now().UnixMilli()
	e.quotes["AAPL"] = Quote{Symbol: "AAPL", Price: 100, Source: "alpaca-iex-websocket-trade", FeedType: "iex", DataState: "live", ProviderTimestamp: base, UpdatedAt: base}
	e.updateQuote("AAPL", Quote{Price: 100.10, ProviderTimestamp: base + 1000, FeedType: "iex", DataState: "live"}, "alpaca-iex-websocket-trade")
	e.updateQuote("AAPL", Quote{Price: 100.10, ProviderTimestamp: base + 2000, FeedType: "rest", DataState: "snapshot"}, "alpaca-iex-snapshot")
	e.mu.RLock()
	d := e.canonicalPipeline
	e.mu.RUnlock()
	if d.MaterialQuoteChanges != 2 || d.PersistenceEnqueues != 2 || d.SuppressedDownstream != 0 {
		t.Fatalf("material changes did not propagate: %+v", d)
	}
}

func TestV172CatalystQuotePropagationIsSymbolScoped(t *testing.T) {
	e, p := newV172PipelineTestEngine(t)
	defer p.Close()
	now := time.Now()
	base := now.UnixMilli()
	e.catalystReactions["AAPL"] = CatalystReactionState{Symbol: "AAPL", TriggerType: "NEWS", Trigger: "material event", TriggerAt: now.Add(-time.Minute).UnixMilli(), TriggerPrice: 100}
	e.quotes["AAPL"] = Quote{Symbol: "AAPL", Price: 100, PreviousClose: 99, Source: "test-live", FeedType: "websocket", DataState: "live", ProviderTimestamp: base, UpdatedAt: base}
	e.quotes["MSFT"] = Quote{Symbol: "MSFT", Price: 100, PreviousClose: 99, Source: "test-live", FeedType: "websocket", DataState: "live", ProviderTimestamp: base, UpdatedAt: base}
	e.updateQuote("MSFT", Quote{Price: 101, ProviderTimestamp: base + 1000, FeedType: "websocket", DataState: "live"}, "test-live")
	e.mu.RLock()
	unrelated := e.canonicalPipeline.CatalystEvaluations
	e.mu.RUnlock()
	if unrelated != 0 {
		t.Fatalf("unrelated symbol triggered catalyst: %d", unrelated)
	}
	e.updateQuote("AAPL", Quote{Price: 101, ProviderTimestamp: base + 2000, FeedType: "websocket", DataState: "live"}, "test-live")
	e.mu.RLock()
	related := e.canonicalPipeline.CatalystEvaluations
	e.mu.RUnlock()
	if related != 1 {
		t.Fatalf("active catalyst symbol evaluations=%d", related)
	}
}

func TestV172SnapshotProvidersUseCanonicalPropagationOwner(t *testing.T) {
	e, p := newV172PipelineTestEngine(t)
	defer p.Close()
	e.mergeFinnhubSnapshot("NVDA", finnhubQuoteResponse{Current: 200, Open: 198, High: 201, Low: 197, Previous: 199, Timestamp: time.Now().Unix()})
	e.mergeAlpacaLiveSnapshot("AAPL", 150, time.Now().UnixMilli(), alpacaLiveSnapshot{}, "iex", "trade")
	e.mu.RLock()
	d := e.canonicalPipeline
	e.mu.RUnlock()
	if d.ReceivedQuoteChanges != 2 || d.MaterialQuoteChanges != 2 || d.PersistenceEnqueues != 2 {
		t.Fatalf("snapshot providers bypassed pipeline: %+v", d)
	}
}

func TestV172DecisionLineagePayloadIsFrozenFromLaterOutcomes(t *testing.T) {
	s := SignalSnapshot{ID: "AAPL-day-evidence-1-v1", Symbol: "AAPL", Horizon: "day", Timestamp: 1000, Price: 100, Score: 82, Action: "WATCH", Readiness: "READY", EvidenceSnapshotID: "evidence-1", FormulaVersion: "protected-v1", FamilyScores: map[string]float64{"momentum": 80}, OutcomeState: "TARGET HIT", OutcomeUpdatedAt: 5000, MFE: 3.2, MAE: -.7}
	b := signalSnapshotPersistenceBatch(s)
	if len(b.Decisions) != 1 || len(b.Outcomes) != 1 {
		t.Fatalf("decision/outcome separation missing: %+v", b)
	}
	var d map[string]any
	if err := json.Unmarshal(b.Decisions[0].Payload, &d); err != nil {
		t.Fatal(err)
	}
	for _, k := range []string{"outcomeState", "outcomeUpdatedAt", "mfe", "mae", "outcomes", "entryTouchedAt", "targetTouchedAt", "invalidationAt"} {
		if _, ok := d[k]; ok {
			t.Fatalf("frozen decision leaked %q: %s", k, string(b.Decisions[0].Payload))
		}
	}
	if d["action"] != "WATCH" || d["formulaVersion"] != "protected-v1" {
		t.Fatalf("decision truth missing: %s", string(b.Decisions[0].Payload))
	}
}

func TestV172DerivedFeatureHashIgnoresOutcomeOnlyChanges(t *testing.T) {
	s := SignalSnapshot{ID: "NVDA-swing-1", Symbol: "NVDA", Horizon: "swing", Timestamp: 1000, FormulaVersion: "protected-v1", FamilyScores: map[string]float64{"trend": 71, "quality": 66}, Readiness: "READY", MarketRegime: "NEUTRAL", MarketTradeability: "NORMAL"}
	a := signalSnapshotPersistenceBatch(s)
	s.Timestamp = 9000
	s.OutcomeState = "TARGET HIT"
	s.OutcomeUpdatedAt = 9000
	s.MFE = 4.1
	b := signalSnapshotPersistenceBatch(s)
	if len(a.Features) != 1 || len(b.Features) != 1 || a.Features[0].SourceHash != b.Features[0].SourceHash {
		t.Fatalf("outcome-only change altered feature hash")
	}
}

func TestV172AlpacaIEXStreamDoesNotSelfDeadlockOnAllocationRead(t *testing.T) {
	e, p := newV172PipelineTestEngine(t)
	defer p.Close()
	observed := time.Date(2026, time.August, 12, 14, 0, 0, 0, easternLocation())
	done := make(chan struct{})
	go func() {
		e.mergeAlpacaIEXStreamAt("SPY", 600.25, 600.20, 600.30, observed.UnixMilli(), "trade", observed)
		close(done)
	}()
	select {
	case <-done:
	case <-time.After(500 * time.Millisecond):
		t.Fatal("Alpaca IEX stream self-deadlocked while resolving allocation")
	}
}
