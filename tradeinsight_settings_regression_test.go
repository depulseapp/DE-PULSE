package main

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestTradeInsightAPIKeyPrecedenceAndEnvironmentFallback(t *testing.T) {
	t.Setenv("TIDATA_API_KEY", "env-primary")
	t.Setenv("TRADEINSIGHT_API_KEY", "env-legacy")
	if got := tradeInsightAPIKey("persisted-settings"); got != "persisted-settings" {
		t.Fatalf("persisted Settings key must take precedence, got %q", got)
	}
	if got := tradeInsightAPIKey(); got != "env-primary" {
		t.Fatalf("TIDATA_API_KEY must remain canonical env fallback, got %q", got)
	}
	t.Setenv("TIDATA_API_KEY", "")
	if got := tradeInsightAPIKey(); got != "env-legacy" {
		t.Fatalf("legacy TRADEINSIGHT_API_KEY compatibility fallback lost, got %q", got)
	}
}

func TestTradeInsightEngineResolverAndRouterUsePersistedSettings(t *testing.T) {
	t.Setenv("TIDATA_API_KEY", "env-primary")
	app := &Application{secrets: Secrets{TradeInsight: "persisted-settings"}}
	e := &Engine{app: app}
	if got := e.tradeInsightResolvedAPIKey(); got != "persisted-settings" {
		t.Fatalf("engine must resolve persisted key first, got %q", got)
	}
	if !e.tradeInsightConfigured() {
		t.Fatal("engine must report TradeInsight configured from persisted Settings")
	}
	if !e.providerConfigured("tradeinsight", app.secrets, Settings{}) {
		t.Fatal("Smart Provider Router v2 configuration truth must recognize persisted TradeInsight key")
	}
}

func TestTradeInsightPublicStateRedactsPersistedSecret(t *testing.T) {
	const secret = "tradeinsight-super-secret"
	app := &Application{state: defaultState(), secrets: Secrets{TradeInsight: secret}}
	public := app.publicStateLockedForUser(bootstrapOwnerID)
	if !public.HasTradeInsightKey {
		t.Fatal("public state must expose configured boolean")
	}
	payload, err := json.Marshal(public)
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(string(payload), secret) {
		t.Fatal("TradeInsight secret leaked into public state")
	}
}

func TestTradeInsightProviderTestMissingKeyDoesNotCallNetwork(t *testing.T) {
	got := testTradeInsight(context.Background(), "")
	if got.OK || got.Status != "missing" || got.Provider != "tradeinsight" {
		t.Fatalf("unexpected missing-key result: %+v", got)
	}
}

func TestTradeInsightRendererSettingsWiring(t *testing.T) {
	raw, err := os.ReadFile(filepath.Join("renderer", "renderer.js"))
	if err != nil {
		t.Fatal(err)
	}
	src := string(raw)
	for _, needle := range []string{
		"tradeinsight-key",
		"state.hasTradeInsightKey",
		"data-test-provider=\"tradeinsight\"",
		"data-clear-secret=\"tradeinsight\"",
		"tradeInsightKey:$('#tradeinsight-key')",
		"tradeInsightKey:body.tradeInsightKey",
	} {
		if !strings.Contains(src, needle) {
			t.Fatalf("renderer missing TradeInsight Settings wiring %q", needle)
		}
	}
}
