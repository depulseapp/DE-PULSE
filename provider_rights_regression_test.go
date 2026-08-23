package main

import (
	"encoding/json"
	"os"
	"reflect"
	"strings"
	"testing"
)

func TestV184ProviderDataRightsDefaultFailsClosed(t *testing.T) {
	providers := []string{"Alpaca", "Finnhub", "Twelve Data", "Marketaux", "FRED", "SEC EDGAR", "yfinance", "CBOE"}
	for _, provider := range providers {
		rights := providerDataRightsMetadata(provider)
		if rights.PolicyVersion != providerDataRightsPolicyVersion {
			t.Fatalf("%s policy version = %q", provider, rights.PolicyVersion)
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
	hop := ProviderRouteHop{
		Provider:    "Finnhub",
		Entitlement: providerCapabilitySupported,
		DataRights:  providerDataRightsMetadata("Finnhub"),
	}
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
	if !ok {
		t.Fatalf("dataRights missing from provider route hop: %s", raw)
	}
	if rights["reviewState"] != providerRightsUnreviewed || rights["commercialUse"] != providerRightsNotAsserted || rights["redistribution"] != providerRightsNotAsserted || rights["aiUse"] != providerRightsNotAsserted {
		t.Fatalf("rights metadata not conservative: %#v", rights)
	}
}

func TestV184ProviderRightsDoNotChangeSmartRouterScore(t *testing.T) {
	cap := ProviderCapabilityStateRecord{Provider: "Finnhub", Dataset: "News", State: providerCapabilitySupported}
	circuit := providerCircuit{}
	telemetry := ProviderRequestDiagnostics{Provider: "Finnhub", Successes: 10}
	before := smartRouteScore("Finnhub", "News", 1, WorkTierUserActionable, cap, circuit, telemetry, "regular")
	_ = providerDataRightsMetadata("Finnhub")
	after := smartRouteScore("Finnhub", "News", 1, WorkTierUserActionable, cap, circuit, telemetry, "regular")
	if !reflect.DeepEqual(before, after) {
		t.Fatalf("governance metadata changed route scoring: before=%+v after=%+v", before, after)
	}
}

func TestV184ProviderRightsStayOutOfExecutableRouting(t *testing.T) {
	smartRouter, err := os.ReadFile("smart_router_v2.go")
	if err != nil {
		t.Fatal(err)
	}
	for _, forbidden := range []string{"DataRights", "dataRights", "providerDataRights", "CommercialUse", "Redistribution", "AIUse"} {
		if strings.Contains(string(smartRouter), forbidden) {
			t.Fatalf("Smart Router eligibility/scoring now references governance-only rights marker %q", forbidden)
		}
	}

	router, err := os.ReadFile("provider_router.go")
	if err != nil {
		t.Fatal(err)
	}
	source := string(router)
	start := strings.Index(source, "func (e *Engine) executeProviderRoute")
	end := strings.Index(source, "func (e *Engine) buildProviderRouterSnapshot")
	if start < 0 || end <= start {
		t.Fatal("could not isolate executable provider routing authority")
	}
	executable := source[start:end]
	for _, forbidden := range []string{"DataRights", "dataRights", "providerDataRights", "CommercialUse", "Redistribution", "AIUse"} {
		if strings.Contains(executable, forbidden) {
			t.Fatalf("executable provider routing now references governance-only rights marker %q", forbidden)
		}
	}
}
