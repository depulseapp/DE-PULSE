package main

import (
	"strings"
	"testing"
)

// tradeInsightCanonicalOwnerDependency is certification-only metadata for the
// durable G3 canonical-owner map from scope issue #61. Keeping it in test scope
// prevents provider-specific ownership metadata from becoming runtime state.
type tradeInsightCanonicalOwnerDependency struct {
	Concern        string
	CanonicalOwner string
	Dependencies   []string
	Contract       string
}

func tradeInsightCanonicalOwnerDependencies() []tradeInsightCanonicalOwnerDependency {
	return []tradeInsightCanonicalOwnerDependency{
		{
			Concern:        "congressional-trading",
			CanonicalOwner: "Research/Event alternative-evidence boundary",
			Dependencies:   []string{"research_truth.go", "event_intelligence.go", "evidence_alternative.go", "EvidenceRecord"},
			Contract:       "Disclosure lag and provenance remain explicit; Event Intelligence stays fetch-free; Congressional evidence is governed production alternative evidence after #78 approval/#84 closure and never deterministic trade truth.",
		},
		{
			Concern:        "sec-form4-enrichment",
			CanonicalOwner: "sec_intelligence.go",
			Dependencies:   []string{"research_truth.go", "evidence_alternative.go", "edgar_enrichment.go", "EvidenceRecord"},
			Contract:       "Direct SEC/EDGAR remains authoritative; TradeInsight can only corroborate or enrich after endpoint/schema admission and source-family evidence is not double-counted.",
		},
		{
			Concern:        "opportunity-radar-movers",
			CanonicalOwner: "market_activity_corporate.go -> opportunity_radar.go / Discovery",
			Dependencies:   []string{"existing market-event candidate flow", "shared source metadata", "shared freshness metadata"},
			Contract:       "Provider movers enter the canonical Discovery/Opportunity Radar path as candidate evidence; no TradeInsight-specific scanner or ranking owner is allowed.",
		},
		{
			Concern:        "symbol-validation-fallback",
			CanonicalOwner: "PersistenceManager / Global Symbol Registry + existing symbol validation",
			Dependencies:   []string{"data.go", "router_extreme_v18_8.go"},
			Contract:       "Provider lookup is fallback-only after canonical/local misses; final validation remains U.S.-equity-only while preserving GLD/SLV/USO as explicit tradable exceptions.",
		},
		{
			Concern:        "provider-telemetry-usefulness",
			CanonicalOwner: "ProviderTelemetry / RuntimeLoad diagnostics + Smart Router capability records",
			Dependencies:   []string{"shared provider health", "loading", "staleness", "usefulness"},
			Contract:       "TradeInsight capability metrics attach to shared telemetry and usefulness owners; no provider-specific telemetry store is allowed.",
		},
		{
			Concern:        "freshness-cache-persistence",
			CanonicalOwner: "existing freshness/cache owners + generic EvidenceRecord / DerivedFeatureRecord persistence",
			Dependencies:   []string{"freshness_adaptive.go", "shared caches", "persistence_repository.go", "persistence_intelligence.go"},
			Contract:       "Reuse canonical freshness, cache and persistence ownership; no TradeInsight-specific shadow tables or duplicate canonical state are allowed.",
		},
		{
			Concern:        "quota-rate-limit-backpressure",
			CanonicalOwner: "provider capability registry + router health/backoff path",
			Dependencies:   []string{"shared provider health", "shared backoff", "Smart Provider Router v2 capability state"},
			Contract:       "HTTP 429/backpressure may temporarily de-rank the affected capability but must never reorder the fixed historical provider route.",
		},
		{
			Concern:        "shadow-lifecycle-promotion",
			CanonicalOwner: "Adaptive Intelligence / Smart Provider Router capability state + evidence logs",
			Dependencies:   []string{"shared capability state", "validation evidence", "review/promotion gate"},
			Contract:       "Lifecycle is shadow_registered -> shadow_running -> proposed_promotion -> approved_live; promotion requires evidence and explicit approval, never automatic promotion.",
		},
	}
}

func v189TradeInsightAdmissionByID(t *testing.T, id string) tradeInsightCapabilityAdmission {
	t.Helper()
	for _, row := range tradeInsightCapabilityAdmissionRegistry() {
		if row.ID == id {
			return row
		}
	}
	t.Fatalf("TradeInsight admission %q is missing", id)
	return tradeInsightCapabilityAdmission{}
}

func TestV189TradeInsightAdmissionRegistryCoversIssue61Matrix(t *testing.T) {
	expected := []string{
		"daily-history", "adjusted-history", "corporate-actions", "bulk-history",
		"congressional-trades", "sec-form4", "top-movers", "symbol-search",
		"generic-market-price", "mcp-interface", "python-sdk", "vendor-derived-scores",
	}
	rows := tradeInsightCapabilityAdmissionRegistry()
	if len(rows) != len(expected) {
		t.Fatalf("admission rows = %d, want %d", len(rows), len(expected))
	}
	seen := map[string]bool{}
	for _, row := range rows {
		if row.ID == "" || row.Capability == "" || row.Disposition == "" || row.Consumer == "" || row.Authority == "" || row.Lifecycle == "" {
			t.Fatalf("incomplete admission row: %+v", row)
		}
		if seen[row.ID] {
			t.Fatalf("duplicate admission id %q", row.ID)
		}
		seen[row.ID] = true
	}
	for _, id := range expected {
		if !seen[id] {
			t.Fatalf("missing admission id %q", id)
		}
	}
}

func TestV189TradeInsightOnlyVerifiedWiredCapabilitiesAreRuntimeAdmitted(t *testing.T) {
	allowed := map[string]bool{"daily-history": true, "adjusted-history": true, "corporate-actions": true, "bulk-history": true, "congressional-trades": true}
	for _, row := range tradeInsightCapabilityAdmissionRegistry() {
		if got := row.runtimeAdmitted(); got != allowed[row.ID] {
			t.Fatalf("runtimeAdmitted(%s) = %v, want %v (row=%+v)", row.ID, got, allowed[row.ID], row)
		}
		if row.RuntimeEnabled && !row.SchemaVerified {
			t.Fatalf("unverified schema must never be runtime enabled: %+v", row)
		}
	}
}

func TestV189TradeInsightLifecycleTruthNeverAdvertisesGatedCapability(t *testing.T) {
	for _, id := range []string{"daily-history", "adjusted-history", "corporate-actions", "bulk-history", "congressional-trades"} {
		if got := tradeInsightCapabilityLifecycleTruth(id); got != "PRODUCTION" {
			t.Fatalf("lifecycle truth for %s = %q, want PRODUCTION after #78 approval/#84 closure candidate", id, got)
		}
	}
	for _, id := range []string{"sec-form4", "top-movers", "symbol-search", "generic-market-price", "mcp-interface", "python-sdk", "vendor-derived-scores", "unknown-capability"} {
		if got := tradeInsightCapabilityLifecycleTruth(id); got != "GATED" {
			t.Fatalf("lifecycle truth for %s = %q, want GATED", id, got)
		}
	}
}

func TestV189TradeInsightCanonicalOwnerMapCoversRemainingG3Concerns(t *testing.T) {
	expected := []string{
		"congressional-trading",
		"sec-form4-enrichment",
		"opportunity-radar-movers",
		"symbol-validation-fallback",
		"provider-telemetry-usefulness",
		"freshness-cache-persistence",
		"quota-rate-limit-backpressure",
		"shadow-lifecycle-promotion",
	}
	rows := tradeInsightCanonicalOwnerDependencies()
	if len(rows) != len(expected) {
		t.Fatalf("canonical-owner rows = %d, want %d", len(rows), len(expected))
	}
	seen := map[string]tradeInsightCanonicalOwnerDependency{}
	for _, row := range rows {
		if strings.TrimSpace(row.Concern) == "" || strings.TrimSpace(row.CanonicalOwner) == "" || len(row.Dependencies) == 0 || strings.TrimSpace(row.Contract) == "" {
			t.Fatalf("incomplete canonical-owner row: %+v", row)
		}
		if _, exists := seen[row.Concern]; exists {
			t.Fatalf("duplicate canonical-owner concern %q", row.Concern)
		}
		seen[row.Concern] = row
	}
	for _, concern := range expected {
		if _, ok := seen[concern]; !ok {
			t.Fatalf("missing canonical-owner concern %q", concern)
		}
	}

	sec := strings.ToLower(seen["sec-form4-enrichment"].Contract)
	if !strings.Contains(sec, "direct sec/edgar remains authoritative") {
		t.Fatalf("SEC owner contract must preserve direct SEC/EDGAR authority: %q", seen["sec-form4-enrichment"].Contract)
	}
	congress := strings.ToLower(seen["congressional-trading"].Contract)
	if !strings.Contains(congress, "disclosure lag") || !strings.Contains(congress, "fetch-free") || !strings.Contains(congress, "never deterministic") {
		t.Fatalf("Congress owner contract must retain disclosure lag, fetch-free Event Intelligence and non-deterministic truth: %q", seen["congressional-trading"].Contract)
	}
	symbols := seen["symbol-validation-fallback"].Contract
	for _, symbol := range []string{"GLD", "SLV", "USO"} {
		if !strings.Contains(symbols, symbol) {
			t.Fatalf("symbol validation contract must preserve %s tradable exception: %q", symbol, symbols)
		}
	}
	quota := strings.ToLower(seen["quota-rate-limit-backpressure"].Contract)
	if !strings.Contains(quota, "429") || !strings.Contains(quota, "must never reorder the fixed historical provider route") {
		t.Fatalf("quota/backpressure contract must protect fixed history order: %q", seen["quota-rate-limit-backpressure"].Contract)
	}
	promotion := strings.ToLower(seen["shadow-lifecycle-promotion"].Contract)
	if !strings.Contains(promotion, "shadow_registered") || !strings.Contains(promotion, "approved_live") || !strings.Contains(promotion, "never automatic promotion") {
		t.Fatalf("SHADOW lifecycle contract must preserve explicit promotion gates: %q", seen["shadow-lifecycle-promotion"].Contract)
	}
}

func TestV189TradeInsightCongressOfficialSchemaIsShadowAdmitted(t *testing.T) {
	row := v189TradeInsightAdmissionByID(t, "congressional-trades")
	if !strings.Contains(row.EndpointEvidence, "/trading-data/v1/congress/v1/trades") || !strings.Contains(row.EndpointEvidence, "insight-data.mdx") {
		t.Fatalf("congress endpoint/schema evidence = %q", row.EndpointEvidence)
	}
	if !row.SchemaVerified || !row.RuntimeEnabled || !row.runtimeAdmitted() || row.Lifecycle != "PRODUCTION" {
		t.Fatalf("Congress verified schema must be governed PRODUCTION after #78 approval/#84 closure candidate: %+v", row)
	}
	if !strings.Contains(strings.ToLower(row.GateReason), "#84") || !strings.Contains(strings.ToLower(row.Authority), "never becomes deterministic") {
		t.Fatalf("Congress production admission must retain #84 gate provenance and non-deterministic authority: %+v", row)
	}
}

func TestV189TradeInsightDoesNotInventUnverifiedRESTEndpoints(t *testing.T) {
	for _, id := range []string{"sec-form4", "generic-market-price", "vendor-derived-scores"} {
		row := v189TradeInsightAdmissionByID(t, id)
		if strings.TrimSpace(row.EndpointEvidence) != "" {
			t.Fatalf("%s must not invent endpoint evidence: %q", id, row.EndpointEvidence)
		}
		if row.runtimeAdmitted() {
			t.Fatalf("%s must not be runtime admitted", id)
		}
	}
	for _, id := range []string{"top-movers", "symbol-search"} {
		row := v189TradeInsightAdmissionByID(t, id)
		if !strings.Contains(row.EndpointEvidence, "MCP") || !strings.Contains(row.EndpointEvidence, "REST path/output schema not verified") {
			t.Fatalf("%s evidence must distinguish documented MCP capability from unverified REST contract: %q", id, row.EndpointEvidence)
		}
		if row.runtimeAdmitted() {
			t.Fatalf("%s must remain non-executable until REST/schema proof", id)
		}
	}
}

func TestV189TradeInsightBulkHistoryUsesOnlyBoundedCanonicalFanout(t *testing.T) {
	row := v189TradeInsightAdmissionByID(t, "bulk-history")
	if !row.SchemaVerified {
		t.Fatal("bounded bulk history should inherit the verified per-ticker daily history schema")
	}
	evidence := strings.ToLower(row.EndpointEvidence)
	if !strings.Contains(evidence, "client-side") || !strings.Contains(evidence, "no server-side bulk endpoint") {
		t.Fatalf("bulk semantics must record official SDK fan-out truth: %q", row.EndpointEvidence)
	}
	if !row.RuntimeEnabled || !row.runtimeAdmitted() || row.Lifecycle != "PRODUCTION" {
		t.Fatalf("bounded canonical bulk history must be governed PRODUCTION after #78 approval/#84 closure candidate: %+v", row)
	}
	reason := strings.ToLower(row.GateReason)
	for _, required := range []string{"canonical historical bars", "deduplicated", "sequential", "50 symbols", "#78 approval", "#84"} {
		if !strings.Contains(reason, required) {
			t.Fatalf("bulk history admission must document %q: %q", required, row.GateReason)
		}
	}
}
