package main

import (
	"depulse/internal/providerlifecycle"
	"encoding/json"
	"os"
	"testing"
)

type tradeInsightProductionReadinessArtifact struct {
	Schema               string `json:"schema"`
	ScopeID              string `json:"scopeId"`
	Issue                int    `json:"issue"`
	CommonLifecycleOwner string `json:"commonLifecycleOwner"`
	Approval             struct {
		Status          string `json:"status"`
		TargetLifecycle string `json:"targetLifecycle"`
		ApprovalRef     string `json:"approvalRef"`
		ProductionGate  string `json:"productionGate"`
	} `json:"approval"`
	FixtureEvidence providerlifecycle.Evidence `json:"fixtureEvidence"`
	Candidates      []struct {
		ID                string   `json:"id"`
		Qualification     string   `json:"qualification"`
		ApprovedLifecycle string   `json:"approvedLifecycle"`
		RegressionIDs     []string `json:"regressionIds"`
	} `json:"candidates"`
	HardGated []struct {
		ID string `json:"id"`
	} `json:"hardGated"`
}

func loadTradeInsightProductionReadiness(t *testing.T) tradeInsightProductionReadinessArtifact {
	t.Helper()
	raw, err := os.ReadFile("governance/data-health/tradeinsight-production-readiness.json")
	if err != nil {
		t.Fatal(err)
	}
	var out tradeInsightProductionReadinessArtifact
	if err := json.Unmarshal(raw, &out); err != nil {
		t.Fatalf("invalid TradeInsight readiness artifact: %v", err)
	}
	return out
}

func TestTradeInsight78QualificationUsesCommonLifecycleThresholds(t *testing.T) {
	artifact := loadTradeInsightProductionReadiness(t)
	if artifact.Schema != "DE.PULSE-TRADEINSIGHT-PRODUCTION-READINESS-1" || artifact.ScopeID != "ADAPT-TRADEINSIGHT-PRODUCTION-001" || artifact.Issue != 78 {
		t.Fatalf("TradeInsight readiness identity drift: %+v", artifact)
	}
	if artifact.CommonLifecycleOwner != "internal/providerlifecycle" {
		t.Fatalf("TradeInsight must reuse common lifecycle owner, got %q", artifact.CommonLifecycleOwner)
	}
	if artifact.Approval.Status != "EXPLICIT_GOVERNED_APPROVAL" || artifact.Approval.TargetLifecycle != providerlifecycle.Approved || artifact.Approval.ProductionGate != "#84/ADAPT-DATAHEALTH-CLOSURE-001" {
		t.Fatalf("unexpected governed approval record: %+v", artifact.Approval)
	}

	expected := map[string]bool{
		"daily-history": true,
		"adjusted-history": true,
		"corporate-actions": true,
		"bulk-history": true,
		"congressional-trades": true,
	}
	if len(artifact.Candidates) != len(expected) {
		t.Fatalf("qualification candidates=%d want=%d", len(artifact.Candidates), len(expected))
	}
	for _, candidate := range artifact.Candidates {
		if !expected[candidate.ID] {
			t.Fatalf("unexpected qualification candidate %q", candidate.ID)
		}
		row, ok := tradeInsightCapabilityAdmissionLookup(candidate.ID)
		if !ok || !row.runtimeAdmitted() {
			t.Fatalf("candidate %q is not runtime-admitted through the existing contract: %+v", candidate.ID, row)
		}
		if row.Lifecycle != providerlifecycle.Shadow {
			t.Fatalf("#78 source admission must remain SHADOW until governed approval is consumed: %+v", row)
		}
		policy, ok := row.readinessPolicy()
		if !ok {
			t.Fatalf("candidate %q missing common readiness policy", candidate.ID)
		}
		diagnostic := providerlifecycle.Evaluate(policy, artifact.FixtureEvidence)
		if diagnostic.Readiness != providerlifecycle.ReadinessReadyForPromotion || diagnostic.Lifecycle != providerlifecycle.Shadow {
			t.Fatalf("candidate %q did not qualify without lifecycle mutation: %+v", candidate.ID, diagnostic)
		}
		if candidate.Qualification != providerlifecycle.ReadinessReadyForPromotion || candidate.ApprovedLifecycle != providerlifecycle.Approved || len(candidate.RegressionIDs) == 0 {
			t.Fatalf("candidate %q missing qualification/approval/regression evidence: %+v", candidate.ID, candidate)
		}
		delete(expected, candidate.ID)
	}
	if len(expected) != 0 {
		t.Fatalf("missing qualification candidates: %v", expected)
	}
}

func TestTradeInsight78HardGatesRemainNonExecutable(t *testing.T) {
	artifact := loadTradeInsightProductionReadiness(t)
	want := map[string]bool{"sec-form4": true, "top-movers": true, "symbol-search": true}
	if len(artifact.HardGated) != len(want) {
		t.Fatalf("hard-gated count=%d want=%d", len(artifact.HardGated), len(want))
	}
	for _, gated := range artifact.HardGated {
		if !want[gated.ID] {
			t.Fatalf("unexpected hard-gated capability %q", gated.ID)
		}
		row, ok := tradeInsightCapabilityAdmissionLookup(gated.ID)
		if !ok || row.runtimeAdmitted() || row.SchemaVerified || row.RuntimeEnabled {
			t.Fatalf("hard-gated capability escaped admission boundary: %+v", row)
		}
		if _, ok := row.readinessPolicy(); ok {
			t.Fatalf("hard-gated capability %q must not receive a promotion policy", gated.ID)
		}
		delete(want, gated.ID)
	}
	if len(want) != 0 {
		t.Fatalf("missing hard-gated capabilities: %v", want)
	}
}

func TestTradeInsight78RuntimePressureCannotBecomeGovernedApproval(t *testing.T) {
	row, _ := tradeInsightCapabilityAdmissionLookup("daily-history")
	policy, ok := row.readinessPolicy()
	if !ok {
		t.Fatal("daily history readiness policy missing")
	}
	evidence := loadTradeInsightProductionReadiness(t).FixtureEvidence
	evidence.RateLimited = 8
	evidence.Errors = 8
	evidence.Successes = 32
	got := providerlifecycle.Evaluate(policy, evidence)
	if got.Readiness == providerlifecycle.ReadinessReadyForPromotion || got.Lifecycle != providerlifecycle.Shadow {
		t.Fatalf("runtime pressure must block readiness without mutating governed lifecycle: %+v", got)
	}
}
