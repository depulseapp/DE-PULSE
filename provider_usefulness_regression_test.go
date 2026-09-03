package main

import (
	"context"
	"os"
	"strings"
	"sync"
	"testing"
	"time"
)

func usefulnessObservation(provider string, price float64, at int64) ProviderQuoteObservation {
	return ProviderQuoteObservation{Symbol: "SPY", Provider: provider, Price: price, ProviderTimestamp: at, ReceivedAt: at, Source: provider}
}

func usefulnessDecision(symbol, state, canonical string, at int64, observations ...ProviderQuoteObservation) ProviderReconciliationDecision {
	for i := range observations {
		observations[i].Symbol = symbol
	}
	return ProviderReconciliationDecision{Dataset: "US Live Equities", Symbol: symbol, State: state, CanonicalProvider: canonical, CanonicalValue: 100, Observations: observations, UpdatedAt: at}
}

func usefulnessDiagnosticsByProvider(rows []ProviderUsefulnessDiagnostic) map[string]ProviderUsefulnessDiagnostic {
	out := map[string]ProviderUsefulnessDiagnostic{}
	for _, row := range rows {
		out[row.Provider] = row
	}
	return out
}

func TestProviderUsefulnessSeparatesSemanticEvidenceFromTransportReliability(t *testing.T) {
	now := time.Now().UnixMilli()
	tracker := newProviderUsefulnessTracker()
	decisions := []ProviderReconciliationDecision{}
	for i, symbol := range []string{"SPY", "QQQ", "DIA", "IWM"} {
		at := now + int64(i*1000)
		decisions = append(decisions, usefulnessDecision(symbol, "AGREED", "Alpaca", at,
			usefulnessObservation("Alpaca", 100, at), usefulnessObservation("Finnhub", 100.02, at)))
	}
	conflictAt := now + 5000
	decisions = append(decisions, usefulnessDecision("XLK", "CONFLICT", "Alpaca", conflictAt,
		usefulnessObservation("Alpaca", 100, conflictAt), usefulnessObservation("Finnhub", 101, conflictAt)))
	singleAt := now + 6000
	decisions = append(decisions, usefulnessDecision("GLD", "SINGLE SOURCE", "Twelve Data", singleAt,
		usefulnessObservation("Twelve Data", 300, singleAt)))
	decisions = append(decisions, usefulnessDecision("VIX", "STALE", "yfinance", now+7000))

	if !tracker.observe(decisions, now+8000) {
		t.Fatal("expected first semantic evidence batch to change usefulness aggregates")
	}
	rows := usefulnessDiagnosticsByProvider(tracker.diagnostics())
	alpaca := rows["Alpaca"]
	if alpaca.EligibleSamples != 5 || alpaca.CrossSourceSamples != 5 || alpaca.AgreementSamples != 4 || alpaca.ConflictSamples != 1 || alpaca.CanonicalSelections != 5 {
		t.Fatalf("unexpected Alpaca semantic aggregate: %+v", alpaca)
	}
	if alpaca.State != "OBSERVING" || alpaca.AgreementPct == nil || *alpaca.AgreementPct != 80 {
		t.Fatalf("expected evidence-qualified 80%% agreement, got %+v", alpaca)
	}
	if alpaca.RoutingImpact != "ADVISORY_ONLY" {
		t.Fatalf("semantic usefulness must remain advisory, got %q", alpaca.RoutingImpact)
	}
	finnhub := rows["Finnhub"]
	if finnhub.CanonicalSelections != 0 || finnhub.CrossSourceSamples != 5 {
		t.Fatalf("unexpected Finnhub aggregate: %+v", finnhub)
	}
	twelve := rows["Twelve Data"]
	if twelve.State != "INSUFFICIENT" || twelve.SingleSourceSamples != 1 || twelve.AgreementPct != nil {
		t.Fatalf("single-source evidence must stay insufficient: %+v", twelve)
	}
	yfinance := rows["yfinance"]
	if yfinance.ExcludedSamples != 1 || yfinance.EligibleSamples != 0 {
		t.Fatalf("stale/invalid evidence must not become an eligible sample: %+v", yfinance)
	}
	if tracker.observe(decisions, now+9000) {
		t.Fatal("identical reconciliation decisions must be deduplicated")
	}
	if got := usefulnessDiagnosticsByProvider(tracker.diagnostics())["Alpaca"].EligibleSamples; got != 5 {
		t.Fatalf("duplicate snapshot inflated evidence count to %d", got)
	}
}

func TestProviderUsefulnessExcludesNonContemporaneousAndFutureEvidence(t *testing.T) {
	now := time.Now().UnixMilli()
	observations := []ProviderQuoteObservation{
		usefulnessObservation("Alpaca", 100, now),
		usefulnessObservation("Finnhub", 100, now-int64(3*time.Minute/time.Millisecond)),
	}
	kept := keepContemporaneousProviderObservations(observations)
	if len(kept) != 1 || kept[0].Provider != "Alpaca" {
		t.Fatalf("expected only freshest contemporaneous observation, got %+v", kept)
	}
	tracker := newProviderUsefulnessTracker()
	tracker.observe([]ProviderReconciliationDecision{usefulnessDecision("SPY", "SINGLE SOURCE", "Alpaca", now, kept...)}, now)
	rows := usefulnessDiagnosticsByProvider(tracker.diagnostics())
	if rows["Alpaca"].EligibleSamples != 1 {
		t.Fatalf("fresh observation was not counted: %+v", rows["Alpaca"])
	}
	if _, ok := rows["Finnhub"]; ok {
		t.Fatal("non-contemporaneous observation must not be counted as semantic evidence")
	}

	future := Quote{Price: 100, ProviderTimestamp: now + quoteMaxFutureSkewMs + 1, UpdatedAt: now}
	_, _, valid, detail := quoteEvidenceTimestampTruth(future, now)
	if valid || !strings.Contains(strings.ToLower(detail), "future") {
		t.Fatalf("material future timestamp must fail evidence validity: valid=%v detail=%q", valid, detail)
	}
	tracker.observe([]ProviderReconciliationDecision{usefulnessDecision("QQQ", "STALE", "Finnhub", now+1)}, now+1)
	row := usefulnessDiagnosticsByProvider(tracker.diagnostics())["Finnhub"]
	if row.ExcludedSamples != 1 || row.EligibleSamples != 0 {
		t.Fatalf("future/invalid exclusion must not be promoted to eligible evidence: %+v", row)
	}
}

func TestProviderUsefulnessSeenLedgerIsBounded(t *testing.T) {
	tracker := newProviderUsefulnessTracker()
	now := time.Now().UnixMilli()
	for i := 0; i < providerUsefulnessSeenLimit+40; i++ {
		symbol := "SYM" + jsonInt64(int64(i))
		at := now + int64(i)
		tracker.observe([]ProviderReconciliationDecision{usefulnessDecision(symbol, "SINGLE SOURCE", "Alpaca", at,
			usefulnessObservation("Alpaca", 100+float64(i), at))}, at)
	}
	if len(tracker.seenOrder) != providerUsefulnessSeenLimit || len(tracker.seen) != providerUsefulnessSeenLimit {
		t.Fatalf("seen ledger must remain bounded at %d, got order=%d map=%d", providerUsefulnessSeenLimit, len(tracker.seenOrder), len(tracker.seen))
	}
}

type providerUsefulnessFakeBackend struct {
	mu       sync.Mutex
	features map[string]DerivedFeatureRecord
}

func newProviderUsefulnessFakeBackend() *providerUsefulnessFakeBackend {
	return &providerUsefulnessFakeBackend{features: map[string]DerivedFeatureRecord{}}
}
func (b *providerUsefulnessFakeBackend) Name() string           { return "provider-usefulness-test" }
func (b *providerUsefulnessFakeBackend) Capabilities() []string { return []string{"intelligence"} }
func (b *providerUsefulnessFakeBackend) Init(context.Context) error {
	return nil
}
func (b *providerUsefulnessFakeBackend) UpsertSymbols(context.Context, []SymbolRegistryRecord) (int, error) {
	return 0, nil
}
func (b *providerUsefulnessFakeBackend) LoadSymbols(context.Context) ([]SymbolRegistryRecord, error) {
	return nil, nil
}
func (b *providerUsefulnessFakeBackend) SaveQuotes(context.Context, map[string]Quote) (int, error) {
	return 0, nil
}
func (b *providerUsefulnessFakeBackend) LoadQuotes(context.Context) (map[string]Quote, error) {
	return nil, nil
}
func (b *providerUsefulnessFakeBackend) SaveIntelligence(_ context.Context, batch PersistenceIntelligenceBatch) (int, error) {
	b.mu.Lock()
	defer b.mu.Unlock()
	for _, feature := range batch.Features {
		key := feature.Symbol + "|" + feature.FeatureKey + "|" + feature.FeatureVersion
		b.features[key] = feature
	}
	return batch.Len(), nil
}
func (b *providerUsefulnessFakeBackend) LoadIdentityState(context.Context) (IdentityPersistentState, error) {
	return IdentityPersistentState{}, nil
}
func (b *providerUsefulnessFakeBackend) SaveIdentityState(context.Context, IdentityPersistentState) error {
	return nil
}
func (b *providerUsefulnessFakeBackend) LoadUserWorkspaces(context.Context) ([]UserWorkspace, error) {
	return nil, nil
}
func (b *providerUsefulnessFakeBackend) SaveUserWorkspace(context.Context, UserWorkspace) error {
	return nil
}
func (b *providerUsefulnessFakeBackend) Stats(context.Context) (PersistenceStoreStats, error) {
	return PersistenceStoreStats{}, nil
}
func (b *providerUsefulnessFakeBackend) Close() error { return nil }
func (b *providerUsefulnessFakeBackend) ExportPersistenceArchive(context.Context) (PersistenceArchive, error) {
	b.mu.Lock()
	defer b.mu.Unlock()
	archive := PersistenceArchive{SchemaVersion: persistenceArchiveSchemaVersion, Features: make([]DerivedFeatureRecord, 0, len(b.features))}
	for _, feature := range b.features {
		archive.Features = append(archive.Features, feature)
	}
	return archive, nil
}
func (b *providerUsefulnessFakeBackend) RestorePersistenceArchive(context.Context, PersistenceArchive, string) error {
	return nil
}

func TestProviderUsefulnessPersistsBoundedAggregateAndDedupState(t *testing.T) {
	backend := newProviderUsefulnessFakeBackend()
	persistence := newPersistenceManagerWithBackend(backend)
	defer persistence.Close()
	tracker := newProviderUsefulnessTracker()
	tracker.bindPersistence(persistence)
	now := time.Now().UnixMilli()
	rows := []ProviderReconciliationDecision{usefulnessDecision("SPY", "AGREED", "Alpaca", now,
		usefulnessObservation("Alpaca", 100, now), usefulnessObservation("Finnhub", 100.01, now))}
	if !tracker.observe(rows, now) {
		t.Fatal("expected usefulness observation")
	}
	tracker.persist(persistence, now)
	persistence.flushPending()

	restarted := newProviderUsefulnessTracker()
	restarted.bindPersistence(persistence)
	got := usefulnessDiagnosticsByProvider(restarted.diagnostics())
	if got["Alpaca"].EligibleSamples != 1 || got["Finnhub"].EligibleSamples != 1 {
		t.Fatalf("restart did not restore bounded usefulness aggregate: %+v", got)
	}
	if restarted.observe(rows, now+1000) {
		t.Fatal("persisted decision signature must prevent restart double-counting")
	}
}

func TestProviderTelemetryUIIsPrivilegedAndRouterIndependent(t *testing.T) {
	ui, err := os.ReadFile("renderer/provider-telemetry-ui.js")
	if err != nil {
		t.Fatal(err)
	}
	text := string(ui)
	for _, token := range []string{"SUPER_OWNER", "OWNER", "ADMIN", "Transport reliability", "Operational Scorecard", "transport UNKNOWN", "Semantic Evidence", "successPct", "p50LatencyMs", "p95LatencyMs", "OBSERVABILITY ONLY", "ADVISORY ONLY", "no routing effect"} {
		if !strings.Contains(text, token) {
			t.Fatalf("provider telemetry UI missing %q", token)
		}
	}
	if strings.Contains(text, "privilegedRoles=new Set(['USER'") || strings.Contains(text, "privilegedRoles=new Set(['DEMO'") {
		t.Fatal("normal USER/DEMO roles must not receive the new internal provider diagnostics")
	}
	index, err := os.ReadFile("renderer/index.html")
	if err != nil {
		t.Fatal(err)
	}
	indexText := string(index)
	rendererAt := strings.Index(indexText, "renderer.js")
	telemetryAt := strings.Index(indexText, "provider-telemetry-ui.js")
	if rendererAt < 0 || telemetryAt <= rendererAt {
		t.Fatal("provider telemetry UI extension must load after the canonical renderer")
	}
	router, err := os.ReadFile("smart_router_v2.go")
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(string(router), "ProviderUsefulness") || strings.Contains(string(router), "providerUsefulness") {
		t.Fatal("semantic usefulness telemetry must not feed Smart Provider Router v2 ordering")
	}
	if strings.Contains(string(router), "ProviderOperationalScorecard") || strings.Contains(string(router), "providerScorecard") {
		t.Fatal("operational scorecard projection must not feed Smart Provider Router v2 ordering")
	}
	runtimeOwner, err := os.ReadFile("provider_usefulness_runtime.go")
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(runtimeOwner), "buildProviderReconciliation") {
		t.Fatal("semantic usefulness must derive from canonical reconciliation truth")
	}
}

func TestProviderOperationalScorecardsPreserveMeasuredAndUnknownTruth(t *testing.T) {
	agreement := 80.0
	registrations := []ProviderRegistration{
		{Name: "Measured", CostClass: "Contract tier", Configured: func(Settings, Secrets) bool { return true }, Routes: []ProviderDatasetContract{{Dataset: "Quotes", Capability: "Quotes"}}},
		{Name: "Unknown", CostClass: "Provider dependent", Configured: func(Settings, Secrets) bool { return true }, Routes: []ProviderDatasetContract{{Dataset: "Fundamentals", Capability: "Fundamentals"}}},
	}
	router := ProviderRouterSnapshot{Routes: []ProviderRouteState{
		{Dataset: "Quotes", Serving: "Measured", Route: []ProviderRouteHop{{Provider: "Measured", Configured: true, Health: "AVAILABLE", Circuit: "CLOSED", Attempts: 10, LastSuccess: 100}}},
		{Dataset: "Fundamentals", Route: []ProviderRouteHop{{Provider: "Unknown", Configured: true, Health: "AVAILABLE", Circuit: "CLOSED"}}},
	}}
	rows := buildProviderOperationalScorecards(registrations, Settings{}, Secrets{}, router,
		[]FreshnessDiagnostic{{Dataset: "Quotes", Provider: "Measured", State: "LIVE"}},
		[]ProviderRequestDiagnostics{{Provider: "Measured", Successes: 9, Errors: 1, SuccessPct: 90, P50LatencyMs: 25, P95LatencyMs: 75, BudgetPerMinute: 30, BudgetRemaining: 20}},
		[]ProviderUsefulnessDiagnostic{{Provider: "Measured", State: "OBSERVING", CrossSourceSamples: 10, AgreementPct: &agreement, RoutingImpact: "ADVISORY_ONLY"}},
		[]LiveSubscriptionBudgetDiagnostics{{Provider: "Measured stream", Capacity: 20, Available: 7}}, 123)
	byProvider := map[string]ProviderOperationalScorecard{}
	for _, row := range rows {
		byProvider[row.Provider] = row
	}
	measured := byProvider["Measured"]
	if measured.TransportMeasurementState != "MEASURED" || measured.SuccessPct == nil || *measured.SuccessPct != 90 || measured.P50LatencyMs == nil || *measured.P50LatencyMs != 25 || measured.P95LatencyMs == nil || *measured.P95LatencyMs != 75 {
		t.Fatalf("measured transport truth missing: %+v", measured)
	}
	if measured.FreshnessMeasurementState != "MEASURED" || measured.HeadroomMeasurementState != "MEASURED" || measured.RequestBudgetRemaining == nil || *measured.RequestBudgetRemaining != 20 || measured.LiveSubscriptionAvailable == nil || *measured.LiveSubscriptionAvailable != 7 {
		t.Fatalf("measured freshness/headroom truth missing: %+v", measured)
	}
	if measured.AgreementPct == nil || *measured.AgreementPct != 80 || measured.RoutingImpact != "OBSERVABILITY_ONLY" || measured.ObservedCostUSD != nil || measured.CostMeasurementState != "DECLARED_CLASS_ONLY" {
		t.Fatalf("semantic/cost truth must remain advisory and non-invented: %+v", measured)
	}
	unknown := byProvider["Unknown"]
	if unknown.TransportMeasurementState != "UNKNOWN" || unknown.FreshnessMeasurementState != "UNKNOWN" || unknown.HeadroomMeasurementState != "UNKNOWN" || unknown.SuccessPct != nil || unknown.P50LatencyMs != nil || unknown.AgreementPct != nil || unknown.ObservedCostUSD != nil {
		t.Fatalf("unmeasured metrics must remain explicitly unknown: %+v", unknown)
	}
}
