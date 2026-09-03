package main

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
	"time"
)

func clearProviderRightsBundleForTest(t *testing.T) {
	t.Helper()
	// Existing fail-closed regressions exercise the future public-production
	// boundary explicitly. Development-mode regressions override this back to
	// audit-only after clearing the bundle.
	t.Setenv(providerRightsEnforcementModeEnv, providerRightsEnforcementPublicMode)
	t.Setenv(providerRightsBundlePathEnv, "")
	t.Setenv(providerRightsBundleSHA256Env, "")
}

func approvedHostedRightsFixtureFor(provider string) ProviderDataRightsMetadata {
	rights := ProviderDataRightsMetadata{
		PolicyVersion:       providerDataRightsPolicyVersion,
		Provider:            provider,
		ReviewState:         providerRightsApproved,
		CommercialUse:       providerRightsApproved,
		MultiUserUse:        providerRightsApproved,
		Proxying:            providerRightsApproved,
		CachingRetention:    providerRightsApproved,
		Redistribution:      providerRightsApproved,
		Display:             providerRightsApproved,
		AIUse:               providerRightsApproved,
		UsageLimits:         "governed provider contract limits recorded",
		Attribution:         "provider attribution requirements recorded",
		AllowedEnvironments: []string{"stage", "prod"},
		EffectiveAt:         "2026-08-01T00:00:00Z",
		ExpiresAt:           "2026-12-01T00:00:00Z",
		RenewalState:        providerRightsRenewalCurrent,
		EvidenceBound:       true,
		EvidenceRef:         "rights/provider/example/2026-08",
		EvidenceDigest:      "sha256:" + strings.Repeat("a", 64),
		ReviewedAt:          "2026-08-20T00:00:00Z",
		Detail:              "fixture-only reviewed rights record",
	}
	rights.CommercialReadiness = evaluateProviderCommercialReadiness(rights)
	return rights
}

func approvedHostedRightsFixture() ProviderDataRightsMetadata {
	return approvedHostedRightsFixtureFor("Finnhub")
}

func bindHostedRightsBundleForTest(t *testing.T, records ...ProviderDataRightsMetadata) string {
	t.Helper()
	t.Setenv(runtimeModeEnv, "hosted")
	t.Setenv(hostedEnvironmentEnv, "prod")
	t.Setenv(providerRightsEnforcementModeEnv, providerRightsEnforcementPublicMode)
	bundle := ProviderDataRightsBundle{PolicyVersion: providerRightsBundlePolicyVersion, Records: records}
	raw, err := json.Marshal(bundle)
	if err != nil {
		t.Fatal(err)
	}
	path := filepath.Join(t.TempDir(), "provider-rights.json")
	if err := os.WriteFile(path, raw, 0600); err != nil {
		t.Fatal(err)
	}
	sum := sha256.Sum256(raw)
	t.Setenv(providerRightsBundlePathEnv, path)
	t.Setenv(providerRightsBundleSHA256Env, "sha256:"+hex.EncodeToString(sum[:]))
	return path
}

func TestHOST003PublicProviderRightsEnforcementRequiresExplicitActivation(t *testing.T) {
	t.Setenv(runtimeModeEnv, "hosted")
	t.Setenv(hostedEnvironmentEnv, "prod")
	t.Setenv(providerRightsEnforcementModeEnv, "")
	if providerRightsEnforcementActive() {
		t.Fatal("hosted development activated provider-rights enforcement without explicit public-production mode")
	}
	t.Setenv(providerRightsEnforcementModeEnv, providerRightsEnforcementPublicMode)
	if !providerRightsEnforcementActive() {
		t.Fatal("explicit PUBLIC_PRODUCTION mode did not activate provider-rights enforcement")
	}
	t.Setenv(runtimeModeEnv, "desktop")
	if providerRightsEnforcementActive() {
		t.Fatal("desktop runtime activated public provider-rights enforcement")
	}
}

func TestHOST003HostedDevelopmentAuditsRightsWithoutBlockingRouter(t *testing.T) {
	e := newV1801Engine(t)
	e.app.mu.Lock()
	e.app.secrets.Finnhub = "working-test-key"
	e.app.mu.Unlock()
	t.Setenv(runtimeModeEnv, "hosted")
	t.Setenv(hostedEnvironmentEnv, "prod")
	clearProviderRightsBundleForTest(t)
	t.Setenv(providerRightsEnforcementModeEnv, "")

	now := time.Date(2026, 8, 26, 12, 0, 0, 0, time.UTC)
	strict := evaluateProviderHostedRightsDecision(providerDataRightsMetadata("Finnhub"), providerHostedUseProductionServing, "prod", now)
	if strict.Allowed || strict.State != providerHostedRightsBlocked {
		t.Fatalf("strict rights evaluation stopped reporting the unresolved production blocker: %+v", strict)
	}
	runtimeDecision := hostedProviderRightsDecision("Finnhub", providerHostedUseProductionServing, now)
	if !runtimeDecision.Allowed || runtimeDecision.State != providerHostedRightsAuditOnly || len(runtimeDecision.BlockingReasons) == 0 {
		t.Fatalf("hosted development did not preserve audit truth while allowing provider capacity: %+v", runtimeDecision)
	}

	called := 0
	attempts := map[string]providerRouteAttempt{
		"Finnhub": func(context.Context) bool { called++; return true },
	}
	if provider, ok := e.executeProviderRoute(context.Background(), "News", attempts); !ok || provider != "Finnhub" || called != 1 {
		t.Fatalf("hosted development rights audit suppressed configured router capacity: provider=%q ok=%v called=%d", provider, ok, called)
	}
}

func TestHOST003HostedDevelopmentDoesNotSuppressCachePersistenceStreamsOrServing(t *testing.T) {
	e := newV1801Engine(t)
	t.Setenv(runtimeModeEnv, "hosted")
	t.Setenv(hostedEnvironmentEnv, "prod")
	clearProviderRightsBundleForTest(t)
	t.Setenv(providerRightsEnforcementModeEnv, "")

	cache := MarketCache{
		Quotes: map[string]Quote{
			"SPY": {Symbol: "SPY", Price: 600, Source: "finnhub-rest", UpdatedAt: time.Now().UnixMilli()},
			"QQQ": {Symbol: "QQQ", Price: 500, Source: "alpaca-iex", UpdatedAt: time.Now().UnixMilli()},
		},
		History: map[string][]HistoryPoint{"SPY": {{T: time.Now().UnixMilli(), P: 600}}},
		News:    []NewsItem{{ID: "dev-news", Headline: "development provider capacity"}},
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
	if len(replay.Quotes) != 2 || len(replay.History) != 1 || len(replay.News) != 1 {
		t.Fatalf("hosted development cache was rights-filtered: quotes=%d history=%d news=%d", len(replay.Quotes), len(replay.History), len(replay.News))
	}

	inner := &providerRightsMemoryBackend{}
	backend := wrapHostedRightsPersistenceBackend(inner)
	if n, err := backend.SaveQuotes(context.Background(), cache.Quotes); err != nil || n != 2 || len(inner.quotes) != 2 {
		t.Fatalf("hosted development persistence lost provider quotes: n=%d err=%v quotes=%+v", n, err, inner.quotes)
	}
	batch := PersistenceIntelligenceBatch{Evidence: []EvidenceRecord{{ID: "external-dev", Kind: "provider", Source: "Finnhub", ObservedAt: time.Now().UnixMilli()}}}
	if n, err := backend.SaveIntelligence(context.Background(), batch); err != nil || n != 1 || len(inner.intelligence.Evidence) != 1 {
		t.Fatalf("hosted development persistence lost provider evidence: n=%d err=%v evidence=%+v", n, err, inner.intelligence.Evidence)
	}

	if len(e.effectiveFinnhubSymbols()) == 0 || len(e.effectiveAlpacaIEXSymbols()) == 0 {
		t.Fatal("hosted development rights audit collapsed live provider subscription capacity")
	}

	snap := RuntimeSnapshot{
		Status:   "running",
		Quotes:   cache.Quotes,
		History:  cache.History,
		News:     cache.News,
		Earnings: []EarningsItem{{Symbol: "SPY", Date: "2026-08-26"}},
		Health:   map[string]string{"existing": "healthy"},
	}
	served := enforceHostedRuntimeRightsSnapshot(snap)
	if len(served.Quotes) != 2 || len(served.History) != 1 || len(served.News) != 1 || len(served.Earnings) != 1 || served.Status != "running" {
		t.Fatalf("hosted development serving was rights-filtered: status=%s quotes=%d history=%d news=%d earnings=%d", served.Status, len(served.Quotes), len(served.History), len(served.News), len(served.Earnings))
	}
}

func TestV184ProviderDataRightsDefaultFailsClosed(t *testing.T) {
	clearProviderRightsBundleForTest(t)
	providers := []string{"Alpaca", "Finnhub", "Twelve Data", "Marketaux", "FRED", "SEC EDGAR", "yfinance", "CBOE"}
	for _, provider := range providers {
		rights := providerDataRightsMetadata(provider)
		if rights.PolicyVersion != providerDataRightsPolicyVersion || rights.Provider != provider {
			t.Fatalf("%s provider/policy binding = %+v", provider, rights)
		}
		if rights.ReviewState != providerRightsUnreviewed {
			t.Fatalf("%s review state = %q; want conservative UNREVIEWED", provider, rights.ReviewState)
		}
		if rights.CommercialUse != providerRightsNotAsserted || rights.Redistribution != providerRightsNotAsserted || rights.AIUse != providerRightsNotAsserted {
			t.Fatalf("%s implicitly asserted rights: %+v", provider, rights)
		}
		if rights.EvidenceBound {
			t.Fatalf("%s claims rights evidence without a bound evidence record", provider)
		}
	}
}

func TestV184ProviderRightsAreSeparateFromOperationalEntitlement(t *testing.T) {
	clearProviderRightsBundleForTest(t)
	hop := ProviderRouteHop{Provider: "Finnhub", Entitlement: providerCapabilitySupported, DataRights: providerDataRightsMetadata("Finnhub")}
	raw, err := json.Marshal(hop)
	if err != nil {
		t.Fatal(err)
	}
	var got map[string]any
	if err := json.Unmarshal(raw, &got); err != nil {
		t.Fatal(err)
	}
	if got["entitlement"] != providerCapabilitySupported {
		t.Fatalf("operational entitlement changed: %#v", got["entitlement"])
	}
	rights, ok := got["dataRights"].(map[string]any)
	if !ok || rights["reviewState"] != providerRightsUnreviewed || rights["commercialUse"] != providerRightsNotAsserted {
		t.Fatalf("rights metadata not conservative/separate: %#v", got["dataRights"])
	}
}

func TestV184ProviderRightsDoNotChangeSmartRouterScore(t *testing.T) {
	cap := ProviderCapabilityStateRecord{Provider: "Finnhub", Dataset: "News", State: providerCapabilitySupported}
	circuit := providerCircuit{}
	telemetry := ProviderRequestDiagnostics{Provider: "Finnhub", Successes: 10}
	before := smartRouteScore("Finnhub", "News", 1, WorkTierUserActionable, cap, circuit, telemetry, "regular")
	_ = evaluateProviderHostedRightsDecision(approvedHostedRightsFixture(), providerHostedUseProductionServing, "prod", time.Date(2026, 8, 25, 12, 0, 0, 0, time.UTC))
	after := smartRouteScore("Finnhub", "News", 1, WorkTierUserActionable, cap, circuit, telemetry, "regular")
	if !reflect.DeepEqual(before, after) {
		t.Fatalf("rights governance mutated Smart Router score: before=%+v after=%+v", before, after)
	}
}

func TestHOST001ProviderRightsMetadataCoversHostedLegalDimensions(t *testing.T) {
	clearProviderRightsBundleForTest(t)
	rights := providerDataRightsMetadata("Finnhub")
	for name, value := range map[string]string{
		"provider":   rights.Provider,
		"commercial": rights.CommercialUse,
		"multi-user": rights.MultiUserUse,
		"proxy":      rights.Proxying,
		"cache":      rights.CachingRetention,
		"redisplay":  rights.Redistribution,
		"display":    rights.Display,
		"AI/derived": rights.AIUse,
		"renewal":    rights.RenewalState,
	} {
		if value == "" {
			t.Fatalf("%s rights dimension is missing", name)
		}
	}
	if rights.EvidenceBound || len(rights.AllowedEnvironments) != 0 || rights.EffectiveAt != "" || rights.ExpiresAt != "" {
		t.Fatalf("unreviewed provider fabricated hosted rights evidence: %+v", rights)
	}
}

func TestHOST002WorkingProviderKeyNeverGrantsHostedRights(t *testing.T) {
	t.Setenv(runtimeModeEnv, "hosted")
	t.Setenv(hostedEnvironmentEnv, "prod")
	clearProviderRightsBundleForTest(t)
	rights := providerDataRightsMetadata("Finnhub")
	now := time.Date(2026, 8, 25, 12, 0, 0, 0, time.UTC)
	decision := evaluateProviderHostedRightsDecision(rights, providerHostedUseCommercialMultiUser, "prod", now)
	if decision.Allowed || decision.State != providerHostedRightsBlocked {
		t.Fatalf("unbound provider rights unexpectedly allowed hosted use: %+v", decision)
	}
	if decision.EvidenceRef != "" || decision.EvidenceDigest != "" {
		t.Fatalf("default rights fabricated provenance: %+v", decision)
	}
}

func TestHOST002HostedRightsRequireBoundReviewableProvenance(t *testing.T) {
	rights := approvedHostedRightsFixture()
	now := time.Date(2026, 8, 25, 12, 0, 0, 0, time.UTC)
	for _, mutate := range []func(*ProviderDataRightsMetadata){
		func(r *ProviderDataRightsMetadata) { r.Provider = "" },
		func(r *ProviderDataRightsMetadata) { r.EvidenceBound = false },
		func(r *ProviderDataRightsMetadata) { r.EvidenceRef = "" },
		func(r *ProviderDataRightsMetadata) { r.EvidenceDigest = "not-a-sha" },
		func(r *ProviderDataRightsMetadata) { r.ReviewState = providerRightsUnreviewed },
	} {
		candidate := rights
		mutate(&candidate)
		decision := evaluateProviderHostedRightsDecision(candidate, providerHostedUseCommercialMultiUser, "prod", now)
		if decision.Allowed {
			t.Fatalf("missing provider/review/provenance unexpectedly allowed: %+v", candidate)
		}
	}
}

func TestHOST002ProviderRightsBundleIsSHA256Pinned(t *testing.T) {
	path := bindHostedRightsBundleForTest(t, approvedHostedRightsFixture())
	if got := providerDataRightsMetadata("Finnhub"); !got.EvidenceBound || got.Provider != "Finnhub" {
		t.Fatalf("valid pinned bundle did not load reviewed provider record: %+v", got)
	}
	if err := os.WriteFile(path, []byte(`{"tampered":true}`), 0600); err != nil {
		t.Fatal(err)
	}
	if got := providerDataRightsMetadata("Finnhub"); got.EvidenceBound || got.ReviewState != providerRightsUnreviewed {
		t.Fatalf("tampered bundle remained trusted: %+v", got)
	}
}

func TestHOST003HostedRightsPurposeEligibilityIsFailClosed(t *testing.T) {
	rights := approvedHostedRightsFixture()
	now := time.Date(2026, 8, 25, 12, 0, 0, 0, time.UTC)
	for _, purpose := range []string{
		providerHostedUseCommercialMultiUser,
		providerHostedUseProxy,
		providerHostedUseCacheRetention,
		providerHostedUseRedisplay,
		providerHostedUseAI,
		providerHostedUseLiveFanout,
		providerHostedUseProductionServing,
	} {
		decision := evaluateProviderHostedRightsDecision(rights, purpose, "prod", now)
		if !decision.Allowed || decision.State != providerHostedRightsAllowed {
			t.Fatalf("approved %s purpose blocked: %+v", purpose, decision)
		}
	}
	if decision := evaluateProviderHostedRightsDecision(rights, "UNKNOWN", "prod", now); decision.Allowed {
		t.Fatalf("unknown hosted purpose must fail closed: %+v", decision)
	}
}

func TestHOST003RightsExpiryDowngradeAndEnvironmentBlockDeterministically(t *testing.T) {
	rights := approvedHostedRightsFixture()
	beforeExpiry := time.Date(2026, 11, 30, 23, 59, 59, 0, time.UTC)
	atExpiry := time.Date(2026, 12, 1, 0, 0, 0, 0, time.UTC)
	if decision := evaluateProviderHostedRightsDecision(rights, providerHostedUseProductionServing, "prod", beforeExpiry); !decision.Allowed {
		t.Fatalf("valid rights blocked before expiry: %+v", decision)
	}
	if decision := evaluateProviderHostedRightsDecision(rights, providerHostedUseProductionServing, "prod", atExpiry); decision.Allowed {
		t.Fatalf("expired rights remained eligible: %+v", decision)
	}

	downgraded := rights
	downgraded.CachingRetention = providerRightsDenied
	if decision := evaluateProviderHostedRightsDecision(downgraded, providerHostedUseProductionServing, "prod", beforeExpiry); decision.Allowed {
		t.Fatalf("rights downgrade did not block production cache/persistence eligibility: %+v", decision)
	}
	if decision := evaluateProviderHostedRightsDecision(rights, providerHostedUseProductionServing, "dev", beforeExpiry); decision.Allowed {
		t.Fatalf("unapproved environment unexpectedly eligible: %+v", decision)
	}
}

func TestHOST003ExecutableRouterFailsClosedThenAdmitsReviewedRights(t *testing.T) {
	e := newV1801Engine(t)
	e.app.mu.Lock()
	e.app.secrets.Finnhub = "working-test-key"
	e.app.mu.Unlock()
	t.Setenv(runtimeModeEnv, "hosted")
	t.Setenv(hostedEnvironmentEnv, "prod")
	clearProviderRightsBundleForTest(t)

	called := 0
	attempts := map[string]providerRouteAttempt{
		"Finnhub": func(context.Context) bool { called++; return true },
	}
	if provider, ok := e.executeProviderRoute(context.Background(), "News", attempts); ok || provider != "" || called != 0 {
		t.Fatalf("unreviewed hosted provider was attempted: provider=%q ok=%v called=%d", provider, ok, called)
	}

	bindHostedRightsBundleForTest(t, approvedHostedRightsFixture())
	if provider, ok := e.executeProviderRoute(context.Background(), "News", attempts); !ok || provider != "Finnhub" || called != 1 {
		t.Fatalf("reviewed hosted provider was not reachable through canonical router: provider=%q ok=%v called=%d", provider, ok, called)
	}
}

func TestHOST003RouterObservabilityShowsRightsBlockWithoutChangingEntitlement(t *testing.T) {
	e := newV1801Engine(t)
	e.app.mu.Lock()
	e.app.secrets.Finnhub = "working-test-key"
	settings := clone(e.app.state.Settings)
	secrets := clone(e.app.secrets)
	e.app.mu.Unlock()
	t.Setenv(runtimeModeEnv, "hosted")
	t.Setenv(hostedEnvironmentEnv, "prod")
	clearProviderRightsBundleForTest(t)

	snap := e.buildProviderRouterSnapshot(settings, secrets, map[string]Quote{}, map[string]int64{})
	found := false
	for _, route := range snap.Routes {
		if route.Dataset != "News" {
			continue
		}
		for _, hop := range route.Route {
			if hop.Provider == "Finnhub" {
				found = true
				if hop.Health != "RIGHTS BLOCKED" || hop.Recovery != "SUPPRESSED" {
					t.Fatalf("rights denial not diagnosable in router: %+v", hop)
				}
				if hop.Entitlement == "RIGHTS BLOCKED" {
					t.Fatalf("legal rights were incorrectly folded into operational entitlement: %+v", hop)
				}
			}
		}
	}
	if !found {
		t.Fatal("Finnhub News hop missing from router diagnostics")
	}
}

type providerRightsMemoryBackend struct {
	quotes       map[string]Quote
	intelligence PersistenceIntelligenceBatch
}

func (b *providerRightsMemoryBackend) Name() string               { return "memory" }
func (b *providerRightsMemoryBackend) Capabilities() []string     { return []string{"test"} }
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
			"SPY": {Symbol: "SPY", Price: 600, Source: "finnhub-rest", UpdatedAt: time.Now().UnixMilli()},
			"QQQ": {Symbol: "QQQ", Price: 500, Source: "alpaca-iex", UpdatedAt: time.Now().UnixMilli()},
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
