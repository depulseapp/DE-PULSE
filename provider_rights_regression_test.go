package main

import (
	"encoding/json"
	"os"
	"reflect"
	"strings"
	"testing"
	"time"
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

func approvedHostedRightsFixture() ProviderDataRightsMetadata {
	rights := ProviderDataRightsMetadata{
		PolicyVersion:       providerDataRightsPolicyVersion,
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
		EvidenceDigest:      "sha256:example-rights-evidence",
		ReviewedAt:          "2026-08-20T00:00:00Z",
		Detail:              "fixture",
	}
	rights.CommercialReadiness = evaluateProviderCommercialReadiness(rights)
	return rights
}

func TestHOST001ProviderRightsMetadataCoversHostedLegalDimensions(t *testing.T) {
	rights := providerDataRightsMetadata("Finnhub")
	for name, value := range map[string]string{
		"commercial": rights.CommercialUse,
		"multi-user": rights.MultiUserUse,
		"proxy":      rights.Proxying,
		"cache":      rights.CachingRetention,
		"redisplay":  rights.Redistribution,
		"display":    rights.Display,
		"AI":         rights.AIUse,
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
		func(r *ProviderDataRightsMetadata) { r.EvidenceBound = false },
		func(r *ProviderDataRightsMetadata) { r.EvidenceRef = "" },
		func(r *ProviderDataRightsMetadata) { r.EvidenceDigest = "" },
		func(r *ProviderDataRightsMetadata) { r.ReviewState = providerRightsUnreviewed },
	} {
		candidate := rights
		mutate(&candidate)
		decision := evaluateProviderHostedRightsDecision(candidate, providerHostedUseCommercialMultiUser, "prod", now)
		if decision.Allowed {
			t.Fatalf("missing review/provenance unexpectedly allowed: %+v", candidate)
		}
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
	if decision := evaluateProviderHostedRightsDecision(rights, providerHostedUseCacheRetention, "prod", beforeExpiry); !decision.Allowed {
		t.Fatalf("valid rights blocked before expiry: %+v", decision)
	}
	if decision := evaluateProviderHostedRightsDecision(rights, providerHostedUseCacheRetention, "prod", atExpiry); decision.Allowed {
		t.Fatalf("expired rights remained eligible: %+v", decision)
	}

	downgraded := rights
	downgraded.CachingRetention = providerRightsDenied
	if decision := evaluateProviderHostedRightsDecision(downgraded, providerHostedUseCacheRetention, "prod", beforeExpiry); decision.Allowed {
		t.Fatalf("rights downgrade did not block cache/persistence eligibility: %+v", decision)
	}
	if decision := evaluateProviderHostedRightsDecision(rights, providerHostedUseCacheRetention, "dev", beforeExpiry); decision.Allowed {
		t.Fatalf("unapproved environment unexpectedly eligible: %+v", decision)
	}
}

func TestHOST003HostedRightsDecisionDoesNotChangeSmartRouterScore(t *testing.T) {
	cap := ProviderCapabilityStateRecord{Provider: "Finnhub", Dataset: "News", State: providerCapabilitySupported}
	circuit := providerCircuit{}
	telemetry := ProviderRequestDiagnostics{Provider: "Finnhub", Successes: 10}
	before := smartRouteScore("Finnhub", "News", 1, WorkTierUserActionable, cap, circuit, telemetry, "regular")

	rights := approvedHostedRightsFixture()
	_ = evaluateProviderHostedRightsDecision(rights, providerHostedUseLiveFanout, "prod", time.Date(2026, 8, 25, 12, 0, 0, 0, time.UTC))
	after := smartRouteScore("Finnhub", "News", 1, WorkTierUserActionable, cap, circuit, telemetry, "regular")
	if !reflect.DeepEqual(before, after) {
		t.Fatalf("hosted rights decision mutated Smart Router score: before=%+v after=%+v", before, after)
	}
}
