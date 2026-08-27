package main

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"os"
	"strings"
	"testing"
)

func TestV184CommercialReadinessDefaultsFailClosedForEveryProvider(t *testing.T) {
	for _, registration := range providerRegistrations() {
		provider := registration.Name
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

func TestHOST002PublicTermsReviewCannotGrantExecutableProviderRights(t *testing.T) {
	const path = "governance/work-slices/ADAPT-HOSTED-TRUST-FOUNDATION-001/provider-public-terms-review.json"
	raw, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	var review struct {
		Schema             string `json:"schema"`
		PolicyVersion      string `json:"policyVersion"`
		ApprovalEligible   bool   `json:"approvalEligible"`
		ProductionDecision string `json:"productionDecision"`
		Providers          []struct {
			Provider                string   `json:"provider"`
			ApprovalEligible        bool     `json:"approvalEligible"`
			ProductionDecision      string   `json:"productionDecision"`
			MissingApprovalEvidence []string `json:"missingApprovalEvidence"`
			Evidence                []struct {
				SourceType string `json:"sourceType"`
				URL        string `json:"url"`
			} `json:"evidence"`
		} `json:"providers"`
	}
	if err := json.Unmarshal(raw, &review); err != nil {
		t.Fatal(err)
	}
	if review.Schema != "DE.PULSE-PROVIDER-PUBLIC-TERMS-REVIEW-1" || review.PolicyVersion == providerRightsBundlePolicyVersion {
		t.Fatalf("public terms review is not isolated from executable provider-rights policy: %+v", review)
	}
	if review.ApprovalEligible || review.ProductionDecision != providerCommercialBlocked {
		t.Fatalf("public terms review became approval-eligible: eligible=%v decision=%q", review.ApprovalEligible, review.ProductionDecision)
	}

	expected := map[string]bool{}
	for _, registration := range providerRegistrations() {
		if _, exists := expected[registration.Name]; exists {
			t.Fatalf("duplicate provider registration name: %q", registration.Name)
		}
		expected[registration.Name] = false
	}
	if len(review.Providers) != len(expected) {
		t.Fatalf("public terms provider count=%d; canonical registration count=%d", len(review.Providers), len(expected))
	}
	for _, provider := range review.Providers {
		if _, ok := expected[provider.Provider]; !ok {
			t.Fatalf("unexpected provider in public terms review: %q", provider.Provider)
		}
		if expected[provider.Provider] {
			t.Fatalf("duplicate provider in public terms review: %q", provider.Provider)
		}
		expected[provider.Provider] = true
		if provider.ApprovalEligible || provider.ProductionDecision != providerCommercialBlocked {
			t.Fatalf("%s public terms entry became approval-eligible: %+v", provider.Provider, provider)
		}
		if len(provider.Evidence) == 0 || len(provider.MissingApprovalEvidence) == 0 {
			t.Fatalf("%s public terms entry lacks evidence or explicit missing approval evidence", provider.Provider)
		}
		for _, evidence := range provider.Evidence {
			if evidence.SourceType == "" || !strings.HasPrefix(evidence.URL, "https://") {
				t.Fatalf("%s has non-reviewable public evidence reference: %+v", provider.Provider, evidence)
			}
		}
	}
	for provider, found := range expected {
		if !found {
			t.Fatalf("registered provider missing from public terms review: %s", provider)
		}
	}

	// Even if an operator accidentally points the executable bundle environment
	// at this reviewed public-terms artifact and pins its exact digest, the
	// different policy contract must fail closed rather than promote rights.
	sum := sha256.Sum256(raw)
	t.Setenv(providerRightsBundlePathEnv, path)
	t.Setenv(providerRightsBundleSHA256Env, "sha256:"+hex.EncodeToString(sum[:]))
	for provider := range expected {
		rights := providerDataRightsMetadata(provider)
		if rights.EvidenceBound || rights.ReviewState != providerRightsUnreviewed || rights.CommercialReadiness.State != providerCommercialBlocked {
			t.Fatalf("%s public terms review leaked into executable rights: %+v", provider, rights)
		}
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
