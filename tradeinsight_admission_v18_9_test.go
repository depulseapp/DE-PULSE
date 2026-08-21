package main

import (
	"strings"
	"testing"
)

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
	allowed := map[string]bool{"daily-history": true, "adjusted-history": true, "corporate-actions": true}
	for _, row := range tradeInsightCapabilityAdmissionRegistry() {
		if got := row.runtimeAdmitted(); got != allowed[row.ID] {
			t.Fatalf("runtimeAdmitted(%s) = %v, want %v (row=%+v)", row.ID, got, allowed[row.ID], row)
		}
		if row.RuntimeEnabled && !row.SchemaVerified {
			t.Fatalf("unverified schema must never be runtime enabled: %+v", row)
		}
	}
}

func TestV189TradeInsightCongressEndpointKnownButSchemaGated(t *testing.T) {
	row := v189TradeInsightAdmissionByID(t, "congressional-trades")
	if row.EndpointEvidence != "/trading-data/v1/congress/v1/trades" {
		t.Fatalf("congress endpoint evidence = %q", row.EndpointEvidence)
	}
	if row.SchemaVerified || row.RuntimeEnabled || row.runtimeAdmitted() {
		t.Fatalf("Congress must remain schema-gated before configured-key contract proof: %+v", row)
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

func TestV189TradeInsightBulkHistoryDoesNotInventServerBulkEndpoint(t *testing.T) {
	row := v189TradeInsightAdmissionByID(t, "bulk-history")
	if !row.SchemaVerified {
		t.Fatal("bounded bulk history should inherit the verified per-ticker daily history schema")
	}
	if !strings.Contains(strings.ToLower(row.EndpointEvidence), "client-side") || !strings.Contains(strings.ToLower(row.EndpointEvidence), "no server-side bulk endpoint") {
		t.Fatalf("bulk semantics must record official SDK fan-out truth: %q", row.EndpointEvidence)
	}
	if row.RuntimeEnabled || row.runtimeAdmitted() {
		t.Fatalf("bulk history must stay gated until bounded canonical job wiring exists: %+v", row)
	}
}
