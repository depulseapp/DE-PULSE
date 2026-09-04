package main

import (
	"reflect"
	"strings"
	"testing"
	"time"

	"depulse/internal/providerlifecycle"
)

func resetProviderConfigurationObservationForTest(e *Engine) {
	providerConfigurationObservations.Delete(e)
}

func TestProviderRegistrationPreservesProductionRouteBaseline(t *testing.T) {
	expected := map[string][]string{
		"US Live Equities":                  {"Alpaca", "Finnhub", "Twelve Data"},
		canonicalUSAssetUniverseDataset:     {"Alpaca"},
		canonicalUSMarketCalendarDataset:    {"Alpaca"},
		canonicalUSCorporateActionsDataset:  {"Alpaca"},
		canonicalGlobalMarketContextDataset: {"Twelve Data"},
		"VIX / Indices":                     {"Twelve Data", "yfinance", "CBOE"},
		canonicalHistoricalBarsDataset:      {"Alpaca", tradeInsightProviderName, "Twelve Data", "yfinance"},
		"News":                              {"Finnhub", "Marketaux"},
		"Earnings":                          {"Finnhub", "yfinance"},
		"Fundamentals":                      {"Finnhub", "SEC", "yfinance"},
		"SEC":                               {"SEC EDGAR"},
		"Macro":                             {"FRED"},
	}
	got := routeChains()
	if !reflect.DeepEqual(got, expected) {
		t.Fatalf("provider registration changed production route baseline:\n got=%#v\nwant=%#v", got, expected)
	}
}

func TestProviderRegistrationProductionAdmissionFailsClosed(t *testing.T) {
	incomplete := ProviderDatasetContract{
		Dataset: "Synthetic Dataset", Capability: "Synthetic Capability", Priority: 1,
		Lifecycle: providerlifecycle.Production,
	}
	if providerDatasetContractProductionReady(incomplete) {
		t.Fatal("PRODUCTION label without adaptive evidence contract must fail closed")
	}
	if len(providerDatasetContractValidationErrors(incomplete)) == 0 {
		t.Fatal("incomplete production contract must expose validation errors")
	}
	regs := []ProviderRegistration{{Name: "Synthetic", Routes: []ProviderDatasetContract{incomplete}}}
	if chain := routeChainsFromProviderRegistrations(regs)["Synthetic Dataset"]; len(chain) != 0 {
		t.Fatalf("incomplete production contract entered canonical route: %v", chain)
	}
}

func TestProviderRegistrationRequiresExplicitProductionLifecycle(t *testing.T) {
	for _, lifecycle := range []string{providerlifecycle.Shadow, providerlifecycle.Validated, providerlifecycle.Approved} {
		route := inheritedProductionRoute("Synthetic", "Synthetic Dataset", "Synthetic Capability", 1, "Test cadence", "Test Consumer")
		route.Lifecycle = lifecycle
		regs := []ProviderRegistration{{Name: "Synthetic", Routes: []ProviderDatasetContract{route}}}
		if chain := routeChainsFromProviderRegistrations(regs)["Synthetic Dataset"]; len(chain) != 0 {
			t.Fatalf("%s capability entered canonical production route: %v", lifecycle, chain)
		}
	}

	route := inheritedProductionRoute("Synthetic", "Synthetic Dataset", "Synthetic Capability", 1, "Test cadence", "Test Consumer")
	regs := []ProviderRegistration{{Name: "Synthetic", Routes: []ProviderDatasetContract{route}}}
	if chain := routeChainsFromProviderRegistrations(regs)["Synthetic Dataset"]; !reflect.DeepEqual(chain, []string{"Synthetic"}) {
		t.Fatalf("complete governed production contract did not enter constructed route: %v", chain)
	}
}

func TestProviderRegistrationConfigurationOwnershipMatchesBaseline(t *testing.T) {
	settings := defaultState().Settings
	secrets := Secrets{
		Finnhub: "fh", TradeInsight: "ti", AlpacaKey: "ak", AlpacaSecret: "as",
		FRED: "fred", EIA: "eia", TwelveData: "td", Marketaux: "ma",
	}
	checks := map[string]bool{
		"Alpaca": true, "Finnhub": true, tradeInsightProviderName: true, "Twelve Data": true,
		"Marketaux": true, "FRED": true, "SEC": false, "SEC EDGAR": false,
		"yfinance": true, "CBOE": true, "BLS": true, "EIA": true,
	}
	for provider, want := range checks {
		if got := providerConfiguredFromRegistration(provider, settings, secrets); got != want {
			t.Fatalf("configured state mismatch for %s: got=%v want=%v", provider, got, want)
		}
	}
	settings.SECEmail = "depulse@example.com"
	if !providerConfiguredFromRegistration("SEC", settings, secrets) || !providerConfiguredFromRegistration("SEC EDGAR", settings, secrets) {
		t.Fatal("SEC registration must remain configured by SEC contact email")
	}
}

func TestProviderConfigurationChangeReopensOnlyEntitlementState(t *testing.T) {
	e := newV1801Engine(t)
	resetProviderConfigurationObservationForTest(e)
	settings := defaultState().Settings
	oldSecrets := Secrets{Finnhub: "old-finnhub", AlpacaKey: "alpaca-key", AlpacaSecret: "alpaca-secret"}
	if changed := e.refreshProviderConfigurationEntitlements(settings, oldSecrets); len(changed) == 0 {
		t.Fatal("first process observation must establish a bounded revalidation boundary")
	}

	now := time.Now()
	finnhubKey := providerCapabilityKey("Finnhub", "News", marketSessionET(now))
	alpacaKey := providerCapabilityKey("Alpaca", "US Live Equities", marketSessionET(now))
	e.mu.Lock()
	e.providerCapabilityStates[finnhubKey] = ProviderCapabilityStateRecord{
		Key: finnhubKey, Provider: "Finnhub", Dataset: "News", InstrumentClass: providerInstrumentClass("News"), Session: marketSessionET(now),
		State: providerCapabilityNotEntitled, Reason: "HTTP 403 plan limited", LastObservedAt: now.UnixMilli(), RevalidateAt: now.Add(12 * time.Hour).UnixMilli(), PolicyVersion: smartRouterPolicyVersion,
	}
	e.providerCapabilityStates[alpacaKey] = ProviderCapabilityStateRecord{
		Key: alpacaKey, Provider: "Alpaca", Dataset: "US Live Equities", InstrumentClass: providerInstrumentClass("US Live Equities"), Session: marketSessionET(now),
		State: providerCapabilitySupported, Reason: "capability served canonical route successfully", LastObservedAt: now.UnixMilli(), PolicyVersion: smartRouterPolicyVersion,
	}
	e.providerCapabilityCircuits[providerCapabilityCircuitKey("Finnhub", "News")] = providerCircuit{Failures: 3, OpenUntil: now.Add(time.Hour).UnixMilli(), LastError: "HTTP 403 plan limited"}
	e.providerCircuits[providerKey("Finnhub")] = providerCircuit{Failures: 3, OpenUntil: now.Add(time.Hour).UnixMilli(), LastError: "HTTP 403 plan limited"}
	e.mu.Unlock()

	if changed := e.refreshProviderConfigurationEntitlements(settings, oldSecrets); len(changed) != 0 {
		t.Fatalf("unchanged provider configuration must not churn entitlement evidence: %v", changed)
	}
	e.mu.RLock()
	unchanged := e.providerCapabilityStates[finnhubKey]
	e.mu.RUnlock()
	if unchanged.State != providerCapabilityNotEntitled {
		t.Fatalf("unchanged configuration unexpectedly reopened entitlement: %+v", unchanged)
	}

	newSecrets := oldSecrets
	newSecrets.Finnhub = "rotated-finnhub"
	changed := e.refreshProviderConfigurationEntitlements(settings, newSecrets)
	if !reflect.DeepEqual(changed, []string{"Finnhub"}) {
		t.Fatalf("only changed provider should be revalidated: %v", changed)
	}
	e.mu.RLock()
	finnhub := e.providerCapabilityStates[finnhubKey]
	alpaca := e.providerCapabilityStates[alpacaKey]
	capCircuit := e.providerCapabilityCircuits[providerCapabilityCircuitKey("Finnhub", "News")]
	globalCircuit := e.providerCircuits[providerKey("Finnhub")]
	e.mu.RUnlock()
	if finnhub.State != providerCapabilityUnknown || finnhub.RevalidateAt != 0 || !strings.Contains(finnhub.Reason, "fresh evidence") {
		t.Fatalf("changed provider credential did not reopen entitlement state: %+v", finnhub)
	}
	if capCircuit.Failures != 0 || capCircuit.OpenUntil != 0 || capCircuit.LastError != "" {
		t.Fatalf("entitlement-caused capability suppression was not cleared: %+v", capCircuit)
	}
	if globalCircuit.Failures != 0 || globalCircuit.OpenUntil != 0 || globalCircuit.LastError != "" {
		t.Fatalf("entitlement-caused provider suppression was not cleared: %+v", globalCircuit)
	}
	if alpaca.State != providerCapabilitySupported {
		t.Fatalf("healthy evidence for unchanged provider was disturbed: %+v", alpaca)
	}
	if strings.Contains(finnhub.Reason, "rotated-finnhub") || strings.Contains(finnhub.Reason, "old-finnhub") {
		t.Fatal("credential material leaked into entitlement diagnostics")
	}
}

func TestRankedProviderRouteObservesProviderConfigChangeBeforeScoring(t *testing.T) {
	e := newV1801Engine(t)
	resetProviderConfigurationObservationForTest(e)
	settings := defaultState().Settings
	oldSecrets := Secrets{Finnhub: "old-finnhub"}
	_ = e.rankedProviderRoute("News", WorkTierUserActionable, settings, oldSecrets, time.Now())

	now := time.Now()
	key := providerCapabilityKey("Finnhub", "News", marketSessionET(now))
	e.mu.Lock()
	e.providerCapabilityStates[key] = ProviderCapabilityStateRecord{
		Key: key, Provider: "Finnhub", Dataset: "News", InstrumentClass: providerInstrumentClass("News"), Session: marketSessionET(now),
		State: providerCapabilityNotEntitled, Reason: "HTTP 403 plan limited", LastObservedAt: now.UnixMilli(), RevalidateAt: now.Add(12 * time.Hour).UnixMilli(), PolicyVersion: smartRouterPolicyVersion,
	}
	e.mu.Unlock()

	newSecrets := oldSecrets
	newSecrets.Finnhub = "rotated-finnhub"
	ranked := e.rankedProviderRoute("News", WorkTierUserActionable, settings, newSecrets, now)
	var finnhub ProviderRouteScore
	found := false
	for _, row := range ranked {
		if row.Provider == "Finnhub" {
			finnhub, found = row, true
			break
		}
	}
	if !found {
		t.Fatalf("Finnhub missing from News route: %+v", ranked)
	}
	if !finnhub.Eligible || finnhub.State == providerCapabilityNotEntitled {
		t.Fatalf("changed provider configuration remained suppressed during canonical ranking: %+v", finnhub)
	}
}

func findAdaptiveProviderManifest(snapshot AdaptiveProviderRegistrySnapshot, provider string) (AdaptiveProviderManifest, bool) {
	want := providerKey(provider)
	for _, row := range snapshot.Providers {
		if row.ProviderID == want {
			return row, true
		}
	}
	return AdaptiveProviderManifest{}, false
}

func findAdaptiveProviderRoute(manifest AdaptiveProviderManifest, dataset, capability string) (AdaptiveProviderRouteManifest, bool) {
	for _, row := range manifest.Routes {
		if row.Dataset == dataset && row.Capability == capability {
			return row, true
		}
	}
	return AdaptiveProviderRouteManifest{}, false
}

func TestAdaptiveProviderRegistryProjectsEveryCanonicalRegistration(t *testing.T) {
	settings := defaultState().Settings
	settings.SECEmail = "depulse@example.com"
	secrets := Secrets{
		Finnhub: "fh", TradeInsight: "ti", AlpacaKey: "ak", AlpacaSecret: "as",
		FRED: "fred", EIA: "eia", TwelveData: "td", Marketaux: "ma",
	}
	regs := providerRegistrations()
	snapshot := adaptiveProviderRegistrySnapshotFromRegistrations(regs, settings, secrets)
	if snapshot.ContractVersion != adaptiveProviderRegistryContractVersion {
		t.Fatalf("registry contract version mismatch: %q", snapshot.ContractVersion)
	}
	if len(snapshot.Providers) != len(regs) {
		t.Fatalf("registry projection lost registrations: got=%d want=%d", len(snapshot.Providers), len(regs))
	}
	seen := map[string]bool{}
	for _, reg := range regs {
		manifest, ok := findAdaptiveProviderManifest(snapshot, reg.Name)
		if !ok {
			t.Fatalf("registration missing from adaptive registry projection: %s", reg.Name)
		}
		if seen[manifest.ProviderID] {
			t.Fatalf("duplicate provider id in registry: %s", manifest.ProviderID)
		}
		seen[manifest.ProviderID] = true
		if manifest.Routable != (len(reg.Routes) > 0) {
			t.Fatalf("routable projection mismatch for %s: %+v", reg.Name, manifest)
		}
		for _, route := range reg.Routes {
			projected, ok := findAdaptiveProviderRoute(manifest, route.Dataset, route.Capability)
			if !ok {
				t.Fatalf("route missing from registry projection: %s / %s / %s", reg.Name, route.Dataset, route.Capability)
			}
			if projected.Lifecycle != route.Lifecycle || projected.CanonicalOwner != route.CanonicalOwner || projected.ContractVersion != route.ContractVersion {
				t.Fatalf("registry rewrote governed route contract for %s / %s: got=%+v source=%+v", reg.Name, route.Dataset, projected, route)
			}
			if projected.ProductionReady != providerDatasetContractProductionReady(route) {
				t.Fatalf("production readiness projection drift for %s / %s", reg.Name, route.Dataset)
			}
		}
	}
}

func TestAdaptiveProviderRegistryRegistrationIsGenericAndDoesNotNeedParallelMaps(t *testing.T) {
	route := inheritedProductionRoute("Synthetic Adapter", "Synthetic Dataset", "Synthetic Capability", 1, "Synthetic cadence", "Synthetic Consumer")
	regs := []ProviderRegistration{{
		Name:       "Synthetic Adapter",
		QuotaLabel: "Synthetic quota",
		CostClass:  "Synthetic cost",
		Configured: func(Settings, Secrets) bool { return true },
		Routes:     []ProviderDatasetContract{route},
		Diagnostics: []ProviderCapabilityDiagnosticRegistration{{
			Capability: "Synthetic diagnostics", Detail: "Synthetic detail", Uses: []string{"Synthetic Consumer"},
		}},
	}}

	snapshot := adaptiveProviderRegistrySnapshotFromRegistrations(regs, defaultState().Settings, Secrets{})
	manifest, ok := findAdaptiveProviderManifest(snapshot, "Synthetic Adapter")
	if !ok || !manifest.Configured || manifest.Configuration != "CONFIGURED" {
		t.Fatalf("synthetic adapter did not self-register through generic descriptor: %+v", manifest)
	}
	if len(manifest.Routes) != 1 || len(manifest.Diagnostics) != 1 {
		t.Fatalf("generic registration projection lost route/diagnostic metadata: %+v", manifest)
	}
	if got := routeChainsFromProviderRegistrations(regs)["Synthetic Dataset"]; !reflect.DeepEqual(got, []string{"Synthetic Adapter"}) {
		t.Fatalf("existing route construction did not discover synthetic registration: %v", got)
	}
}

func TestAdaptiveProviderRegistryCapabilityTruthIsBoundedBeforeRouterEligibility(t *testing.T) {
	settings := defaultState().Settings
	if got := adaptiveProviderRegistryCapabilityState("Finnhub", "News", settings, Secrets{}); got != providerCapabilityNotConfigured {
		t.Fatalf("unconfigured registered capability truth mismatch: %s", got)
	}
	if got := adaptiveProviderRegistryCapabilityState("Finnhub", "Options Chain", settings, Secrets{Finnhub: "configured"}); got != providerCapabilityNotSupported {
		t.Fatalf("unsupported capability truth mismatch: %s", got)
	}
	if got := adaptiveProviderRegistryCapabilityState("Unknown Provider", "News", settings, Secrets{}); got != providerCapabilityNotSupported {
		t.Fatalf("unknown provider must fail closed as unsupported: %s", got)
	}
	if got := adaptiveProviderRegistryCapabilityState("Finnhub", "News", settings, Secrets{Finnhub: "configured"}); got != providerRegistryCapabilityRegistered {
		t.Fatalf("configured governed registration should be REGISTERED, not a routing decision: %s", got)
	}
}

func TestAdaptiveProviderRegistryNeverPromotesLifecycleIntoProduction(t *testing.T) {
	route := inheritedProductionRoute("Synthetic Shadow", "Synthetic Dataset", "Synthetic Capability", 1, "Synthetic cadence", "Synthetic Consumer")
	route.Lifecycle = providerlifecycle.Shadow
	regs := []ProviderRegistration{{Name: "Synthetic Shadow", Configured: func(Settings, Secrets) bool { return true }, Routes: []ProviderDatasetContract{route}}}
	settings := defaultState().Settings

	if got := adaptiveProviderRegistryCapabilityStateFromRegistrations(regs, "Synthetic Shadow", "Synthetic Dataset", settings, Secrets{}); got != providerRegistryCapabilityRegisteredNotProduction {
		t.Fatalf("SHADOW route was not held below production: %s", got)
	}
	if got := routeChainsFromProviderRegistrations(regs)["Synthetic Dataset"]; len(got) != 0 {
		t.Fatalf("registry projection must not promote SHADOW route into canonical production chain: %v", got)
	}
	manifest, ok := findAdaptiveProviderManifest(adaptiveProviderRegistrySnapshotFromRegistrations(regs, settings, Secrets{}), "Synthetic Shadow")
	if !ok || len(manifest.Routes) != 1 || manifest.Routes[0].ProductionReady {
		t.Fatalf("manifest falsely reported SHADOW route production-ready: %+v", manifest)
	}
}

func TestAdaptiveProviderRegistryPreservesExplicitAuthorityClasses(t *testing.T) {
	snapshot := adaptiveProviderRegistrySnapshot(defaultState().Settings, Secrets{})
	sec, ok := findAdaptiveProviderManifest(snapshot, "SEC EDGAR")
	if !ok {
		t.Fatal("SEC EDGAR registration missing")
	}
	secRoute, ok := findAdaptiveProviderRoute(sec, "SEC", "Direct SEC/EDGAR filings and Form 4 authority")
	if !ok || secRoute.AuthorityClass != providerAuthorityDirectAuthority {
		t.Fatalf("direct SEC/EDGAR authority was not preserved: %+v", secRoute)
	}
	yf, ok := findAdaptiveProviderManifest(snapshot, "yfinance")
	if !ok {
		t.Fatal("yfinance registration missing")
	}
	for _, route := range yf.Routes {
		if route.AuthorityClass != providerAuthorityFallback {
			t.Fatalf("yfinance must remain fallback-only in registry projection: %+v", route)
		}
	}
	cboe, ok := findAdaptiveProviderManifest(snapshot, "CBOE")
	if !ok || len(cboe.Routes) == 0 || cboe.Routes[0].AuthorityClass != providerAuthorityCorroborative {
		t.Fatalf("CBOE corroborative role was not preserved: %+v", cboe)
	}
}

func TestAdaptiveProviderRegistryProjectionNeverContainsSecretMaterial(t *testing.T) {
	secret := "APR01-super-secret-material"
	snapshot := adaptiveProviderRegistrySnapshot(defaultState().Settings, Secrets{Finnhub: secret, TwelveData: secret, Marketaux: secret})
	for _, manifest := range snapshot.Providers {
		values := []string{manifest.ProviderID, manifest.Name, manifest.ContractVersion, manifest.Configuration, manifest.QuotaLabel, manifest.CostClass}
		for _, route := range manifest.Routes {
			values = append(values, route.Dataset, route.Capability, route.InstrumentClass, route.AuthorityClass, route.Lifecycle, route.ContractVersion, route.CanonicalOwner, route.SchemaContract, route.TimestampContract, route.FreshnessContract, route.RightsContract, route.ExpectedDelay)
			values = append(values, route.Uses...)
		}
		for _, diagnostic := range manifest.Diagnostics {
			values = append(values, diagnostic.Capability, diagnostic.Detail)
			values = append(values, diagnostic.Uses...)
		}
		if strings.Contains(strings.Join(values, "\x00"), secret) {
			t.Fatalf("registry projection leaked provider credential material for %s", manifest.Name)
		}
	}
	finnhub, ok := findAdaptiveProviderManifest(snapshot, "Finnhub")
	if !ok || !finnhub.Configured || finnhub.Configuration != "CONFIGURED" {
		t.Fatalf("redacted configuration presence was not preserved: %+v", finnhub)
	}
}
