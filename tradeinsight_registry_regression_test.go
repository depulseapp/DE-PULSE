package main

import "testing"

func v189TradeInsightRegistryStatus(t *testing.T, health map[string]string) string {
	t.Helper()
	rows := buildProviderCapabilityRegistry(Settings{}, Secrets{}, health, map[string]SymbolIntelligence{}, map[string]GlobalDriver{})
	for _, row := range rows {
		if row.Provider == tradeInsightProviderName && row.Capability == "Adjusted daily OHLCV / corporate-action corroboration" {
			return row.Status
		}
	}
	t.Fatal("TradeInsight capability registry row is missing")
	return ""
}

func TestV189TradeInsightRegistryNotEntitledWithoutRuntimeSecret(t *testing.T) {
	t.Setenv("TIDATA_API_KEY", "")
	t.Setenv("TRADEINSIGHT_API_KEY", "")
	if got := v189TradeInsightRegistryStatus(t, map[string]string{}); got != "NOT ENTITLED" {
		t.Fatalf("status = %q, want NOT ENTITLED", got)
	}
}

func TestV189TradeInsightRegistryStartsShadowWhenConfigured(t *testing.T) {
	t.Setenv("TIDATA_API_KEY", "fixture")
	t.Setenv("TRADEINSIGHT_API_KEY", "")
	if got := v189TradeInsightRegistryStatus(t, map[string]string{}); got != "SHADOW" {
		t.Fatalf("status = %q, want SHADOW", got)
	}
}

func TestV189TradeInsightRegistryDoesNotBorrowAnotherHistoryProviderHealth(t *testing.T) {
	t.Setenv("TIDATA_API_KEY", "fixture")
	t.Setenv("TRADEINSIGHT_API_KEY", "")
	if got := v189TradeInsightRegistryStatus(t, map[string]string{"history": "healthy · Alpaca daily history"}); got != "SHADOW" {
		t.Fatalf("status = %q, want SHADOW when canonical history is healthy from another provider", got)
	}
}

func TestV189TradeInsightRegistryAvailableOnlyFromTradeInsightEvidence(t *testing.T) {
	t.Setenv("TIDATA_API_KEY", "fixture")
	t.Setenv("TRADEINSIGHT_API_KEY", "")
	if got := v189TradeInsightRegistryStatus(t, map[string]string{"history": "healthy · TradeInsight daily fallback · adjusted OHLCV"}); got != "AVAILABLE" {
		t.Fatalf("status = %q, want AVAILABLE", got)
	}
	if got := v189TradeInsightRegistryStatus(t, map[string]string{"tradeinsight-corporate-actions": "healthy · 2 supplemental dividend/split events normalized into canonical ledger"}); got != "AVAILABLE" {
		t.Fatalf("corporate-action status = %q, want AVAILABLE", got)
	}
}
