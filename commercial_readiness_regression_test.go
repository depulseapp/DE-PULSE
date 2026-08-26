package main

import (
	"encoding/json"
	"os"
	"strings"
	"testing"
)

func TestV184CommercialReadinessDefaultsFailClosedForEveryProvider(t *testing.T) {
	providers := []string{"Alpaca", "Finnhub", "Twelve Data", "Marketaux", "FRED", "SEC EDGAR", "yfinance", "CBOE"}
	for _, provider := range providers {
		rights := providerDataRightsMetadata(provider)
		ready := rights.CommercialReadiness
		if ready.PolicyVersion != providerCommercialReadinessPolicyVersion {
			t.Fatalf("%s commercial-readiness policy = %q", provider, ready.PolicyVersion)
		}
		if ready.State != providerCommercialBlocked || !ready.RequiresReview {
			t.Fatalf("%s was not fail-closed: %+v", provider, ready)
		}
		if ready.CommercialUseReady || ready.RedistributionReady || ready.AIUseReady {
			t.Fatalf("%s received implicit commercial rights: %+v", provider, ready)
		}
		if len(ready.BlockingReasons) < 3 {
			t.Fatalf("%s missing explicit blocking reasons: %+v", provider, ready)
		}
	}
}

func TestV184CommercialReadinessRequiresEvidenceAndExplicitApproval(t *testing.T) {
	base := ProviderDataRightsMetadata{
		PolicyVersion:  providerDataRightsPolicyVersion,
		Provider:       "Finnhub",
		ReviewState:    providerRightsApproved,
		CommercialUse:  providerRightsApproved,
		Redistribution: providerRightsApproved,
		AIUse:          providerRightsApproved,
		EvidenceRef:    "rights/provider/finnhub/test-fixture",
		EvidenceDigest: "sha256:" + strings.Repeat("a", 64),
	}
	withoutEvidence := evaluateProviderCommercialReadiness(base)
	if withoutEvidence.State != providerCommercialBlocked || withoutEvidence.CommercialUseReady || withoutEvidence.RedistributionReady || withoutEvidence.AIUseReady {
		t.Fatalf("approval labels without bound evidence became ready: %+v", withoutEvidence)
	}

	base.EvidenceBound = true
	ready := evaluateProviderCommercialReadiness(base)
	if ready.State != providerCommercialReady || ready.RequiresReview || !ready.CommercialUseReady || !ready.RedistributionReady || !ready.AIUseReady {
		t.Fatalf("synthetic fully evidence-bound approval did not become ready: %+v", ready)
	}
	if len(ready.BlockingReasons) != 0 {
		t.Fatalf("ready state retained blockers: %+v", ready.BlockingReasons)
	}

	base.AIUse = providerRightsNotAsserted
	partial := evaluateProviderCommercialReadiness(base)
	if partial.State != providerCommercialBlocked || partial.AIUseReady || !partial.CommercialUseReady || !partial.RedistributionReady {
		t.Fatalf("partial rights review did not fail closed by dimension: %+v", partial)
	}
}

func TestV184CommercialReadinessIsSerializedBesideOperationalEntitlement(t *testing.T) {
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
		t.Fatalf("operational entitlement drifted: %#v", got["entitlement"])
	}
	rights, ok := got["dataRights"].(map[string]any)
	if !ok {
		t.Fatalf("dataRights missing: %s", raw)
	}
	commercial, ok := rights["commercialReadiness"].(map[string]any)
	if !ok {
		t.Fatalf("commercialReadiness missing from governance metadata: %s", raw)
	}
	if commercial["state"] != providerCommercialBlocked || commercial["requiresReview"] != true {
		t.Fatalf("serialized commercial readiness not fail-closed: %#v", commercial)
	}
}

func TestV184CommercialReadinessNeverEntersExecutableRouting(t *testing.T) {
	smartRouter, err := os.ReadFile("smart_router_v2.go")
	if err != nil {
		t.Fatal(err)
	}
	for _, forbidden := range []string{"CommercialReadiness", "commercialReadiness", "evaluateProviderCommercialReadiness", "providerCommercialReady", "providerCommercialBlocked"} {
		if strings.Contains(string(smartRouter), forbidden) {
			t.Fatalf("Smart Router references governance-only commercial-readiness marker %q", forbidden)
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
	for _, forbidden := range []string{"CommercialReadiness", "commercialReadiness", "evaluateProviderCommercialReadiness", "providerCommercialReady", "providerCommercialBlocked"} {
		if strings.Contains(executable, forbidden) {
			t.Fatalf("executable provider routing references governance-only commercial-readiness marker %q", forbidden)
		}
	}
}
