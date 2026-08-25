package main

import (
	"reflect"
	"strings"
	"testing"
	"time"

	"depulse/internal/providerlifecycle"
)

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
