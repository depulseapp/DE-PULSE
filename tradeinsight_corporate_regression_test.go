package main

import (
	"context"
	"net/http"
	"strings"
	"testing"
	"time"
)

func TestV189TradeInsightCorporateActionsNormalizeHistoryEvidence(t *testing.T) {
	now := time.Date(2026, 8, 21, 12, 0, 0, 0, time.UTC).UnixMilli()
	rows := []map[string]any{
		{"date": "2026-08-20", "dividend": 0.25, "split_ratio": 1.0},
		{"date": "2026-08-21", "dividend": 0.0, "split_ratio": 2.0},
		{"date": "bad-date", "dividend": 9.0, "split_ratio": 9.0},
	}

	actions := tradeInsightCorporateActions("aapl", rows, now)
	if len(actions) != 2 {
		t.Fatalf("actions=%d want=2: %+v", len(actions), actions)
	}
	dividend, split := actions[0], actions[1]
	if dividend.Symbol != "AAPL" || dividend.Type != "cash_dividend" || dividend.ExDate != "2026-08-20" || dividend.CashAmount != 0.25 || dividend.Source != tradeInsightProviderName || dividend.Status != "EFFECTIVE" {
		t.Fatalf("dividend normalization mismatch: %+v", dividend)
	}
	if split.Symbol != "AAPL" || split.Type != "split" || split.ExDate != "2026-08-21" || split.Ratio != 2 || split.AdjustmentFactor != 2 || split.Source != tradeInsightProviderName || split.Status != "EFFECTIVE" {
		t.Fatalf("split normalization mismatch: %+v", split)
	}
}

func TestV189TradeInsightCorporateActionsRejectInvalidOrNeutralRows(t *testing.T) {
	rows := []map[string]any{
		{"date": "2026-08-20", "dividend": 0.0, "split_ratio": 1.0},
		{"date": "", "dividend": 1.0, "split_ratio": 2.0},
	}
	if got := tradeInsightCorporateActions("AAPL", rows, time.Now().UnixMilli()); len(got) != 0 {
		t.Fatalf("unexpected actions: %+v", got)
	}
	if got := tradeInsightCorporateActions("VIX", []map[string]any{{"date": "2026-08-20", "dividend": 1.0}}, time.Now().UnixMilli()); len(got) != 0 {
		t.Fatalf("VIX must not emit corporate actions: %+v", got)
	}
}

func TestV189TradeInsightSupplementalMergePreservesExistingCanonicalAction(t *testing.T) {
	now := time.Date(2026, 8, 21, 12, 0, 0, 0, time.UTC).UnixMilli()
	existing := []CorporateAction{{
		ID: "alpaca-authoritative-id", Symbol: "AAPL", Type: "split", ExDate: "2026-08-20", Ratio: 2, AdjustmentFactor: 2,
		Status: "EFFECTIVE", FirstSeenAt: now - 1000, UpdatedAt: now - 500, Detail: "rate 2:1", Source: "Alpaca",
	}}
	supplemental := []CorporateAction{{
		Symbol: "AAPL", Type: "split", ExDate: "2026-08-20", Ratio: 2, AdjustmentFactor: 2,
		Status: "EFFECTIVE", FirstSeenAt: now, UpdatedAt: now, Detail: "ratio 2 · supplemental adjusted-history evidence", Source: tradeInsightProviderName,
	}}

	merged := mergeSupplementalCorporateActions(existing, supplemental, now)
	if len(merged) != 1 {
		t.Fatalf("semantic duplicate must not grow ledger: %+v", merged)
	}
	if merged[0].Source != "Alpaca" || merged[0].ID != "alpaca-authoritative-id" {
		t.Fatalf("supplemental evidence overwrote existing canonical action: %+v", merged[0])
	}
}

func TestV189TradeInsightSupplementalMergeAddsNovelAction(t *testing.T) {
	now := time.Date(2026, 8, 21, 12, 0, 0, 0, time.UTC).UnixMilli()
	existing := []CorporateAction{{Symbol: "AAPL", Type: "split", ExDate: "2026-08-20", Ratio: 2, Source: "Alpaca"}}
	supplemental := []CorporateAction{{Symbol: "AAPL", Type: "cash_dividend", ExDate: "2026-08-20", CashAmount: 0.25, Source: tradeInsightProviderName}}
	merged := mergeSupplementalCorporateActions(existing, supplemental, now)
	if len(merged) != 2 {
		t.Fatalf("novel supplemental action not added: %+v", merged)
	}
}

func TestADAPTProviderCorporateActionsProductionUsesRouterV2(t *testing.T) {
	e := newV1801Engine(t)
	configureAdaptProviderUniverseAlpaca(e)
	calls := 0
	client := adaptProviderUniverseClient(func(r *http.Request) (*http.Response, error) {
		calls++
		if r.URL.Path != "/v1/corporate-actions" || r.URL.Query().Get("symbols") == "" || r.URL.Query().Get("start") == "" || r.URL.Query().Get("end") == "" {
			t.Fatalf("unexpected corporate-actions request: %s", r.URL.String())
		}
		if r.Header.Get("APCA-API-KEY-ID") != "test-key" || r.Header.Get("APCA-API-SECRET-KEY") != "test-secret" {
			t.Fatal("Alpaca corporate-actions auth headers missing")
		}
		return adaptProviderUniverseResponse(http.StatusOK, `{}`), nil
	})

	e.mu.RLock()
	beforeDecisions := e.smartRouterScorecard.RouteDecisions
	e.mu.RUnlock()
	if !e.refreshAlpacaCorporateActionsWithClient(context.Background(), "test-key", "test-secret", client) {
		t.Fatal("production corporate-actions route failed")
	}
	if calls == 0 {
		t.Fatal("corporate-actions Router attempt did not reach provider transport")
	}

	e.mu.RLock()
	updatedAt := e.lastUpdated["corporate-actions"]
	health := e.health["corporate-actions"]
	capability := e.providerCapabilityCircuits[providerCapabilityCircuitKey("Alpaca", canonicalUSCorporateActionsDataset)]
	afterDecisions := e.smartRouterScorecard.RouteDecisions
	e.mu.RUnlock()
	if updatedAt == 0 || !strings.Contains(health, "healthy · persistent ledger") {
		t.Fatalf("corporate-actions canonical state not completed: updatedAt=%d health=%q", updatedAt, health)
	}
	if capability.LastSuccess == 0 || afterDecisions != beforeDecisions+1 {
		t.Fatalf("corporate-actions Router success evidence missing: capability=%+v decisions=%d->%d", capability, beforeDecisions, afterDecisions)
	}
}

func TestADAPTProviderCorporateActionsFailureIsCapabilityScopedAndPreservesLedgerFreshness(t *testing.T) {
	e := newV1801Engine(t)
	configureAdaptProviderUniverseAlpaca(e)
	const staleAt int64 = 246813579
	e.mu.Lock()
	e.corporateActions = []CorporateAction{{ID: "prior-action", Symbol: "AAPL", Type: "split", ExDate: "2026-08-20", Ratio: 2, Source: "Alpaca"}}
	e.lastUpdated["corporate-actions"] = staleAt
	e.health["corporate-actions"] = "healthy · persistent ledger · prior"
	e.mu.Unlock()
	client := adaptProviderUniverseClient(func(r *http.Request) (*http.Response, error) {
		return adaptProviderUniverseResponse(http.StatusInternalServerError, `{"message":"corporate actions unavailable"}`), nil
	})
	if e.refreshAlpacaCorporateActionsWithClient(context.Background(), "test-key", "test-secret", client) {
		t.Fatal("terminal corporate-actions provider failure unexpectedly succeeded")
	}

	e.mu.RLock()
	ledger := append([]CorporateAction(nil), e.corporateActions...)
	updatedAt := e.lastUpdated["corporate-actions"]
	health := e.health["corporate-actions"]
	global := e.providerCircuits[providerKey("Alpaca")]
	corporateCircuit := e.providerCapabilityCircuits[providerCapabilityCircuitKey("Alpaca", canonicalUSCorporateActionsDataset)]
	universeCircuit := e.providerCapabilityCircuits[providerCapabilityCircuitKey("Alpaca", canonicalUSAssetUniverseDataset)]
	calendarCircuit := e.providerCapabilityCircuits[providerCapabilityCircuitKey("Alpaca", canonicalUSMarketCalendarDataset)]
	liveCircuit := e.providerCapabilityCircuits[providerCapabilityCircuitKey("Alpaca", "US Live Equities")]
	e.mu.RUnlock()
	if len(ledger) != 1 || ledger[0].ID != "prior-action" || updatedAt != staleAt || !strings.Contains(health, "partial · persistent ledger retained") {
		t.Fatalf("failed corporate-actions refresh corrupted ledger/freshness: ledger=%+v updatedAt=%d health=%q", ledger, updatedAt, health)
	}
	if global.Failures != 0 || global.LastFailure != 0 {
		t.Fatalf("corporate-actions failure leaked into global Alpaca circuit: %+v", global)
	}
	if corporateCircuit.Failures != 1 || corporateCircuit.LastFailure == 0 {
		t.Fatalf("corporate-actions capability failure not recorded: %+v", corporateCircuit)
	}
	if universeCircuit.Failures != 0 || calendarCircuit.Failures != 0 || liveCircuit.Failures != 0 || !e.providerAllowedFor(canonicalUSAssetUniverseDataset, "Alpaca") || !e.providerAllowedFor(canonicalUSMarketCalendarDataset, "Alpaca") || !e.providerAllowedFor("US Live Equities", "Alpaca") {
		t.Fatalf("corporate-actions failure suppressed unrelated Alpaca capabilities: universe=%+v calendar=%+v live=%+v", universeCircuit, calendarCircuit, liveCircuit)
	}
}

func TestADAPTProviderCorporateActionsCancellationIsNeutralAndPreservesLedgerState(t *testing.T) {
	e := newV1801Engine(t)
	configureAdaptProviderUniverseAlpaca(e)
	const staleAt int64 = 135792468
	const priorHealth = "healthy · persistent ledger · prior"
	e.mu.Lock()
	e.corporateActions = []CorporateAction{{ID: "prior-action", Symbol: "AAPL", Type: "cash_dividend", ExDate: "2026-08-20", CashAmount: 0.25, Source: "Alpaca"}}
	e.lastUpdated["corporate-actions"] = staleAt
	e.health["corporate-actions"] = priorHealth
	e.mu.Unlock()
	calls := 0
	client := adaptProviderUniverseClient(func(r *http.Request) (*http.Response, error) {
		calls++
		return nil, context.Canceled
	})
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	if e.refreshAlpacaCorporateActionsWithClient(ctx, "test-key", "test-secret", client) {
		t.Fatal("canceled corporate-actions refresh unexpectedly succeeded")
	}
	if calls > 1 {
		t.Fatalf("caller cancellation must stop corporate-actions provider fanout, calls=%d", calls)
	}

	e.mu.RLock()
	ledger := append([]CorporateAction(nil), e.corporateActions...)
	updatedAt := e.lastUpdated["corporate-actions"]
	health := e.health["corporate-actions"]
	global := e.providerCircuits[providerKey("Alpaca")]
	corporateCircuit := e.providerCapabilityCircuits[providerCapabilityCircuitKey("Alpaca", canonicalUSCorporateActionsDataset)]
	e.mu.RUnlock()
	if len(ledger) != 1 || ledger[0].ID != "prior-action" || updatedAt != staleAt || health != priorHealth {
		t.Fatalf("canceled corporate-actions refresh corrupted canonical state: ledger=%+v updatedAt=%d health=%q", ledger, updatedAt, health)
	}
	if global.Failures != 0 || global.LastFailure != 0 || corporateCircuit.Failures != 0 || corporateCircuit.LastFailure != 0 {
		t.Fatalf("caller cancellation poisoned corporate-actions/provider health: global=%+v corporate=%+v", global, corporateCircuit)
	}
}
