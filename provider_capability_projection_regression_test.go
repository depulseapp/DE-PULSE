package main

import (
	"reflect"
	"strings"
	"testing"
	"time"
)

func providerCapabilityProjectionKeys(rows []ProviderCapabilityEntry) []string {
	out := make([]string, 0, len(rows))
	for _, row := range rows {
		out = append(out, row.Provider+"|"+row.Capability)
	}
	return out
}

func TestProviderCapabilityProjectionPreservesLegacyDataEngineOrder(t *testing.T) {
	rows := buildProviderCapabilityRegistryFromRegistrations(
		providerRegistrations(),
		defaultState().Settings,
		Secrets{},
		map[string]string{},
		map[string]SymbolIntelligence{},
		map[string]GlobalDriver{},
	)
	want := []string{
		"Finnhub|Primary U.S. equity + earnings/peers",
		"Finnhub|Analyst / insider premium context",
		"Alpaca|IEX quotes / snapshots / liquidity",
		"Alpaca|SIP movers / most active",
		tradeInsightProviderName + "|Adjusted daily OHLCV / corporate-action corroboration",
		"FRED|Rates / credit / conditions / USD",
		"BLS|Inflation / labor / wages / PPI",
		"EIA|Petroleum / natural gas / energy state",
		"Twelve Data|FX / direct global context",
		"Twelve Data|VIX / indices / historical recovery",
		"yfinance|Recovery-only public market context",
		"CBOE|Official VIX validation / delayed close",
		"Marketaux|Stock news fallback",
	}
	if got := providerCapabilityProjectionKeys(rows); !reflect.DeepEqual(got, want) {
		t.Fatalf("provider registration changed Data Engine capability order:\n got=%#v\nwant=%#v", got, want)
	}
}

func TestProviderCapabilityProjectionPlacesFutureDiagnosticsAfterLegacyRowsDeterministically(t *testing.T) {
	regs := append([]ProviderRegistration(nil), providerRegistrations()...)
	regs = append(regs,
		ProviderRegistration{Name: "ZZZ Future", Diagnostics: []ProviderCapabilityDiagnosticRegistration{{Capability: "Beta Capability"}}},
		ProviderRegistration{Name: "AAA Future", Diagnostics: []ProviderCapabilityDiagnosticRegistration{{Capability: "Zulu Capability"}, {Capability: "Alpha Capability"}}},
	)
	rows := buildProviderCapabilityRegistryFromRegistrations(
		regs,
		defaultState().Settings,
		Secrets{},
		map[string]string{},
		map[string]SymbolIntelligence{},
		map[string]GlobalDriver{},
	)
	got := providerCapabilityProjectionKeys(rows)
	if len(got) != 16 {
		t.Fatalf("expected 13 retained diagnostics plus 3 future diagnostics, got %d: %v", len(got), got)
	}
	wantTail := []string{
		"AAA Future|Alpha Capability",
		"AAA Future|Zulu Capability",
		"ZZZ Future|Beta Capability",
	}
	if !reflect.DeepEqual(got[len(got)-3:], wantTail) {
		t.Fatalf("future provider diagnostics must sort after retained rows without onboarding-map edits: got tail=%v want=%v", got[len(got)-3:], wantTail)
	}
}

func TestProviderCapabilityProjectionRetainsSpecialStatusSemantics(t *testing.T) {
	rows := buildProviderCapabilityRegistryFromRegistrations(
		providerRegistrations(),
		defaultState().Settings,
		Secrets{},
		map[string]string{},
		map[string]SymbolIntelligence{},
		map[string]GlobalDriver{},
	)
	byKey := map[string]string{}
	for _, row := range rows {
		byKey[row.Provider+"|"+row.Capability] = row.Status
	}
	checks := map[string]string{
		"Finnhub|Analyst / insider premium context":                       "NOT ENTITLED",
		"Alpaca|SIP movers / most active":                                 "NOT ENTITLED",
		tradeInsightProviderName + "|Adjusted daily OHLCV / corporate-action corroboration": "NOT ENTITLED",
		"yfinance|Recovery-only public market context":                     "AVAILABLE",
		"CBOE|Official VIX validation / delayed close":                     "AVAILABLE",
	}
	for key, want := range checks {
		if got := byKey[key]; got != want {
			t.Fatalf("provider capability status semantics changed for %s: got=%q want=%q", key, got, want)
		}
	}
}

func TestManualSameKeyCapabilityRecheckReopensOnlyEntitlementSuppression(t *testing.T) {
	e := newV1801Engine(t)
	resetProviderConfigurationObservationForTest(e)
	settings := defaultState().Settings
	secrets := Secrets{Finnhub: "same-key"}
	_ = e.refreshProviderConfigurationEntitlements(settings, secrets)

	now := time.Now()
	newsKey := providerCapabilityKey("Finnhub", "News", marketSessionET(now))
	earningsKey := providerCapabilityKey("Finnhub", "Earnings", marketSessionET(now))
	e.mu.Lock()
	e.providerCapabilityStates[newsKey] = ProviderCapabilityStateRecord{
		Key: newsKey, Provider: "Finnhub", Dataset: "News", InstrumentClass: providerInstrumentClass("News"), Session: marketSessionET(now),
		State: providerCapabilityNotEntitled, Reason: "HTTP 403 plan limited", LastObservedAt: now.UnixMilli(), RevalidateAt: now.Add(12 * time.Hour).UnixMilli(), PolicyVersion: smartRouterPolicyVersion,
	}
	e.providerCapabilityStates[earningsKey] = ProviderCapabilityStateRecord{
		Key: earningsKey, Provider: "Finnhub", Dataset: "Earnings", InstrumentClass: providerInstrumentClass("Earnings"), Session: marketSessionET(now),
		State: providerCapabilityTemporarilyUnavailable, Reason: "timeout", LastObservedAt: now.UnixMilli(), RevalidateAt: now.Add(90 * time.Second).UnixMilli(), PolicyVersion: smartRouterPolicyVersion,
	}
	e.providerCapabilityCircuits[providerCapabilityCircuitKey("Finnhub", "News")] = providerCircuit{Failures: 3, OpenUntil: now.Add(time.Hour).UnixMilli(), LastError: "HTTP 403 plan limited"}
	e.mu.Unlock()

	changed := e.forceProviderEntitlementRevalidation(settings, secrets)
	foundFinnhub := false
	for _, provider := range changed {
		if provider == "Finnhub" {
			foundFinnhub = true
			break
		}
	}
	if !foundFinnhub {
		t.Fatalf("manual capability recheck did not include configured Finnhub: %v", changed)
	}

	e.mu.RLock()
	news := e.providerCapabilityStates[newsKey]
	earnings := e.providerCapabilityStates[earningsKey]
	newsCircuit := e.providerCapabilityCircuits[providerCapabilityCircuitKey("Finnhub", "News")]
	e.mu.RUnlock()
	if news.State != providerCapabilityUnknown || news.RevalidateAt != 0 || !strings.Contains(news.Reason, "manual capability recheck") {
		t.Fatalf("same-key manual recheck did not reopen entitlement evidence: %+v", news)
	}
	if earnings.State != providerCapabilityTemporarilyUnavailable || earnings.Reason != "timeout" {
		t.Fatalf("manual entitlement recheck erased genuine transient health evidence: %+v", earnings)
	}
	if newsCircuit.Failures != 0 || newsCircuit.OpenUntil != 0 || newsCircuit.LastError != "" {
		t.Fatalf("manual entitlement recheck did not clear entitlement-caused capability suppression: %+v", newsCircuit)
	}
}
