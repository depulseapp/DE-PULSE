package main

import (
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
