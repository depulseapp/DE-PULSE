package main

import (
	"encoding/json"
	"os"
	"testing"
)

type lifecycleGovernanceRegistry struct {
	Schema       string `json:"schema"`
	ScopeID      string `json:"scopeId"`
	Issue        int    `json:"issue"`
	Capabilities []struct {
		Provider       string `json:"provider"`
		Capability     string `json:"capability"`
		Lifecycle      string `json:"lifecycle"`
		Authority      string `json:"authority"`
		EvidenceStatus string `json:"evidenceStatus"`
		Readiness      string `json:"readiness"`
	} `json:"capabilities"`
}

type baselineProviderMatrix struct {
	Providers []struct {
		Provider     string `json:"provider"`
		Capabilities []struct {
			Capability string `json:"capability"`
		} `json:"capabilities"`
	} `json:"providers"`
}

func loadLifecycleJSON[T any](t *testing.T, path string) T {
	t.Helper()
	raw, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	var out T
	if err := json.Unmarshal(raw, &out); err != nil {
		t.Fatalf("invalid %s: %v", path, err)
	}
	return out
}

func TestADAPTProviderLifecycleRegistryExhaustsBaselineCapabilityMatrix(t *testing.T) {
	matrix := loadLifecycleJSON[baselineProviderMatrix](t, "governance/data-health/provider-capability-matrix.json")
	registry := loadLifecycleJSON[lifecycleGovernanceRegistry](t, "governance/data-health/provider-lifecycle-readiness.json")
	if registry.Schema != "DE.PULSE-PROVIDER-CAPABILITY-LIFECYCLE-READINESS-1" || registry.ScopeID != "ADAPT-PROVIDER-LIFECYCLE-001" || registry.Issue != 83 {
		t.Fatalf("lifecycle registry identity drift: %+v", registry)
	}
	want := map[string]bool{}
	for _, provider := range matrix.Providers {
		for _, cap := range provider.Capabilities {
			want[provider.Provider+"::"+cap.Capability] = true
		}
	}
	got := map[string]bool{}
	for _, row := range registry.Capabilities {
		key := row.Provider + "::" + row.Capability
		if got[key] {
			t.Fatalf("duplicate lifecycle capability: %s", key)
		}
		got[key] = true
		if row.Authority == "" || row.EvidenceStatus == "" || row.Readiness == "" {
			t.Fatalf("lifecycle row missing authority/evidence/readiness: %+v", row)
		}
		switch row.Lifecycle {
		case providerLifecycleShadow, providerLifecycleValidated, providerLifecycleApproved, providerLifecycleProduction:
		default:
			t.Fatalf("unsupported governed lifecycle %q for %s", row.Lifecycle, key)
		}
	}
	if len(want) != 26 || len(got) != len(want) {
		t.Fatalf("baseline/lifecycle row count drift: baseline=%d lifecycle=%d", len(want), len(got))
	}
	for key := range want {
		if !got[key] {
			t.Fatalf("baseline capability missing lifecycle/readiness classification: %s", key)
		}
	}
	for key := range got {
		if !want[key] {
			t.Fatalf("lifecycle registry contains unowned capability outside #80 baseline: %s", key)
		}
	}
}

func TestADAPTProviderLifecycleTradeInsightRowsRemainShadowFirstUntilGovernedPromotion(t *testing.T) {
	registry := loadLifecycleJSON[lifecycleGovernanceRegistry](t, "governance/data-health/provider-lifecycle-readiness.json")
	count := 0
	for _, row := range registry.Capabilities {
		if row.Provider != tradeInsightProviderName {
			continue
		}
		count++
		if row.Lifecycle != providerLifecycleShadow || row.EvidenceStatus != "COMMON_READINESS_REQUIRED" {
			t.Fatalf("TradeInsight must remain on common SHADOW readiness path before #78 promotion: %+v", row)
		}
	}
	if count != 3 {
		t.Fatalf("expected all three #80 TradeInsight rows, got %d", count)
	}
}
