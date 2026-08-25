package main

import (
	"strings"
	"testing"
	"time"
)

func TestProviderRegistrationProjectsSyntheticCapabilityDiagnostic(t *testing.T) {
	regs := []ProviderRegistration{{
		Name: "Synthetic Provider",
		Diagnostics: []ProviderCapabilityDiagnosticRegistration{{
			Capability: "Synthetic Capability",
			Detail:     "Synthetic diagnostic proves one-contract onboarding projection.",
			Uses:       []string{"Synthetic Consumer"},
			Status: func(Settings, Secrets, map[string]string, map[string]SymbolIntelligence, map[string]GlobalDriver) string {
				return "AVAILABLE"
			},
		}},
	}}
	rows := buildProviderCapabilityRegistryFromRegistrations(regs, Settings{}, Secrets{}, nil, nil, nil)
	if len(rows) != 1 {
		t.Fatalf("synthetic provider diagnostic projection returned %d rows, want 1", len(rows))
	}
	row := rows[0]
	if row.Provider != "Synthetic Provider" || row.Capability != "Synthetic Capability" || row.Status != "AVAILABLE" {
		t.Fatalf("synthetic provider diagnostic was not projected from registration: %+v", row)
	}
	if len(row.Uses) != 1 || row.Uses[0] != "Synthetic Consumer" {
		t.Fatalf("synthetic provider downstream consumer metadata was lost: %+v", row.Uses)
	}
}

func TestRegisteredCapabilityDiagnosticsPreserveSpecializedTruth(t *testing.T) {
	settings := defaultState().Settings
	secrets := Secrets{
		Finnhub: "fh", TradeInsight: "ti", AlpacaKey: "ak", AlpacaSecret: "as",
		FRED: "fred", EIA: "eia", TwelveData: "td", Marketaux: "ma",
	}
	health := map[string]string{
		"quotes-rest": "healthy", "alpaca-live": "healthy", "market-activity": "plan limited",
		"history": "healthy · TradeInsight", "fred-rates": "healthy", "bls-actuals": "official",
		"eia-actuals": "healthy", "global-direct": "healthy", "vix": "healthy",
	}
	rows := buildProviderCapabilityRegistry(settings, secrets, health, nil, map[string]GlobalDriver{
		"fx_usdjpy": {Key: "fx_usdjpy", Source: "Twelve Data"},
	})
	statuses := map[string]string{}
	for _, row := range rows {
		statuses[row.Provider+"|"+row.Capability] = row.Status
	}
	checks := map[string]string{
		"Finnhub|Primary U.S. equity + earnings/peers":                    "AVAILABLE",
		"Finnhub|Analyst / insider premium context":                      "PLAN LIMITED",
		"Alpaca|IEX quotes / snapshots / liquidity":                      "AVAILABLE",
		"Alpaca|SIP movers / most active":                                "PLAN LIMITED",
		tradeInsightProviderName + "|Adjusted daily OHLCV / corporate-action corroboration": "AVAILABLE",
		"FRED|Rates / credit / conditions / USD":                         "AVAILABLE",
		"BLS|Inflation / labor / wages / PPI":                            "AVAILABLE",
		"EIA|Petroleum / natural gas / energy state":                     "AVAILABLE",
		"Twelve Data|FX / direct global context":                         "AVAILABLE",
		"Twelve Data|VIX / indices / historical recovery":                "AVAILABLE",
		"yfinance|Recovery-only public market context":                    "AVAILABLE",
		"CBOE|Official VIX validation / delayed close":                   "AVAILABLE",
		"Marketaux|Stock news fallback":                                  "AVAILABLE",
	}
	if len(statuses) != len(checks) {
		t.Fatalf("provider diagnostic count changed: got=%d want=%d statuses=%v", len(statuses), len(checks), statuses)
	}
	for key, want := range checks {
		if got := statuses[key]; got != want {
			t.Fatalf("diagnostic truth changed for %s: got=%q want=%q", key, got, want)
		}
	}
}

func TestManualEntitlementRevalidationProbesSameKeyPlanChangeWithoutErasingOtherEvidence(t *testing.T) {
	e := newV1801Engine(t)
	resetProviderConfigurationObservationForTest(e)
	settings := defaultState().Settings
	secrets := Secrets{Finnhub: "same-key", AlpacaKey: "ak", AlpacaSecret: "as"}
	_ = e.refreshProviderConfigurationEntitlements(settings, secrets)

	now := time.Now()
	finnhubKey := providerCapabilityKey("Finnhub", "News", marketSessionET(now))
	alpacaKey := providerCapabilityKey("Alpaca", "US Live Equities", marketSessionET(now))
	transientKey := providerCapabilityKey("Finnhub", "Earnings", marketSessionET(now))
	e.mu.Lock()
	e.providerCapabilityStates[finnhubKey] = ProviderCapabilityStateRecord{
		Key: finnhubKey, Provider: "Finnhub", Dataset: "News", InstrumentClass: providerInstrumentClass("News"), Session: marketSessionET(now),
		State: providerCapabilityNotEntitled, Reason: "HTTP 403 plan limited", LastObservedAt: now.UnixMilli(), RevalidateAt: now.Add(12 * time.Hour).UnixMilli(), PolicyVersion: smartRouterPolicyVersion,
	}
	e.providerCapabilityStates[alpacaKey] = ProviderCapabilityStateRecord{
		Key: alpacaKey, Provider: "Alpaca", Dataset: "US Live Equities", InstrumentClass: providerInstrumentClass("US Live Equities"), Session: marketSessionET(now),
		State: providerCapabilitySupported, Reason: "capability served canonical route successfully", LastObservedAt: now.UnixMilli(), PolicyVersion: smartRouterPolicyVersion,
	}
	e.providerCapabilityStates[transientKey] = ProviderCapabilityStateRecord{
		Key: transientKey, Provider: "Finnhub", Dataset: "Earnings", InstrumentClass: providerInstrumentClass("Earnings"), Session: marketSessionET(now),
		State: providerCapabilityTemporarilyUnavailable, Reason: "timeout", LastObservedAt: now.UnixMilli(), RevalidateAt: now.Add(90 * time.Second).UnixMilli(), PolicyVersion: smartRouterPolicyVersion,
	}
	e.mu.Unlock()

	changed := e.forceProviderEntitlementRevalidation(settings, secrets)
	foundFinnhub := false
	for _, provider := range changed {
		if provider == "Finnhub" { foundFinnhub = true }
	}
	if !foundFinnhub {
		t.Fatalf("manual recheck did not include configured Finnhub: %v", changed)
	}

	e.mu.RLock()
	finnhub := e.providerCapabilityStates[finnhubKey]
	alpaca := e.providerCapabilityStates[alpacaKey]
	transient := e.providerCapabilityStates[transientKey]
	e.mu.RUnlock()
	if finnhub.State != providerCapabilityUnknown || finnhub.RevalidateAt != 0 || !strings.Contains(finnhub.Reason, "manual capability recheck") {
		t.Fatalf("same-key plan recheck did not reopen stale entitlement: %+v", finnhub)
	}
	if alpaca.State != providerCapabilitySupported {
		t.Fatalf("manual recheck erased healthy evidence: %+v", alpaca)
	}
	if transient.State != providerCapabilityTemporarilyUnavailable {
		t.Fatalf("manual recheck erased transient health evidence: %+v", transient)
	}
}
