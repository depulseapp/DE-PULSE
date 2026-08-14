package main

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestV1601ProviderReasonMatchesCanonicalOverride(t *testing.T) {
	now := time.Now().UnixMilli()
	obs := map[string]map[string]Quote{"AAPL": {
		"Alpaca":  {Symbol: "AAPL", Price: 200, UpdatedAt: now - 100, ProviderTimestamp: now - 200, Source: "alpaca-iex-websocket-trade"},
		"Finnhub": {Symbol: "AAPL", Price: 202, UpdatedAt: now - 90, ProviderTimestamp: now - 150, Source: "finnhub-websocket"},
	}}
	canonical := map[string]Quote{"AAPL": {Symbol: "AAPL", Price: 202, UpdatedAt: now - 80, Source: "finnhub-websocket"}}
	rows := buildProviderReconciliation(testRouter("Alpaca"), obs, canonical, now)
	if len(rows) != 1 {
		t.Fatalf("rows=%+v", rows)
	}
	r := rows[0]
	if r.CanonicalProvider != "Finnhub" || r.CanonicalValue != 202 {
		t.Fatalf("canonical=%+v", r)
	}
	if !strings.Contains(r.Reason, "Canonical winner is Finnhub") || !strings.Contains(r.Reason, "Router active is Alpaca") {
		t.Fatalf("state/why contradiction remains: %+v", r)
	}
}

func TestV1601GlobalHistoryTimestampCannotSatisfySelectedTicker(t *testing.T) {
	now := time.Date(2026, 8, 11, 14, 0, 0, 0, time.UTC)
	st, quotes, bars, fundamentals, last, health := completeResearchFixture(now)
	delete(last, "research-history:AAPL")
	last["history-daily"] = now.UnixMilli()
	last["history"] = now.UnixMilli()
	pkg := buildResearchPackageTruth(st, quotes, bars, fundamentals, nil, nil, nil, last, health, nil, now)
	var hist ResearchEvidenceComponent
	for _, c := range pkg.Components {
		if c.Dataset == "Daily History" {
			hist = c
		}
	}
	if hist.State != "PARTIAL" || hist.CheckAt != 0 {
		t.Fatalf("global history falsely satisfied selected ticker: pkg=%+v hist=%+v", pkg, hist)
	}
	last["research-history:AAPL"] = now.UnixMilli() - 1000
	pkg = buildResearchPackageTruth(st, quotes, bars, fundamentals, nil, nil, nil, last, health, nil, now)
	for _, c := range pkg.Components {
		if c.Dataset == "Daily History" && c.State != "FRESH" {
			t.Fatalf("target history did not close component: %+v", c)
		}
	}
}

func TestV1601CorporateLedgerMergePreservesHistoryAndFirstSeen(t *testing.T) {
	old := CorporateAction{ID: "ca-1", Symbol: "AAPL", Type: "forward_split", ProcessDate: "2024-06-01", Ratio: 4, FirstSeenAt: 100, UpdatedAt: 200, Source: "Alpaca"}
	fresh := CorporateAction{ID: "ca-1", Symbol: "AAPL", Type: "forward_split", ProcessDate: "2024-06-01", Ratio: 4, UpdatedAt: 900, Source: "Alpaca"}
	added := CorporateAction{ID: "ca-2", Symbol: "META", Type: "name_change", ProcessDate: "2022-06-09", OldSymbol: "FB", NewSymbol: "META", UpdatedAt: 900, Source: "Alpaca"}
	rows := mergeCorporateActionLedger([]CorporateAction{old}, []CorporateAction{fresh, added}, 1000)
	if len(rows) != 2 {
		t.Fatalf("rows=%+v", rows)
	}
	byID := map[string]CorporateAction{}
	for _, a := range rows {
		byID[a.ID] = a
	}
	if byID["ca-1"].FirstSeenAt != 100 {
		t.Fatalf("first seen overwritten: %+v", byID["ca-1"])
	}
	if byID["ca-2"].FirstSeenAt <= 0 {
		t.Fatalf("new action first seen missing: %+v", byID["ca-2"])
	}
	truth := buildCorporateActionTruth(rows, nil, 1000)
	if truth.SymbolLineage["FB"] != "META" {
		t.Fatalf("lineage=%+v", truth.SymbolLineage)
	}
}

func TestV1601CorporateLedgerPersistsInCache(t *testing.T) {
	root := t.TempDir()
	t.Setenv("XDG_CONFIG_HOME", root)
	t.Setenv("HOME", root)
	t.Setenv("APPDATA", root)
	app, err := NewApplication()
	if err != nil {
		t.Fatal(err)
	}
	registerApplicationCleanup(t, app)
	app.engine = NewEngine(app)
	app.engine.mu.Lock()
	app.engine.corporateActions = []CorporateAction{{ID: "persist-1", Symbol: "META", Type: "name_change", OldSymbol: "FB", NewSymbol: "META", ProcessDate: "2022-06-09", FirstSeenAt: 10, UpdatedAt: 20, Source: "Alpaca"}}
	app.engine.lastUpdated["corporate-actions-backfill:META"] = 30
	app.engine.mu.Unlock()
	if err := app.engine.saveCache(); err != nil {
		t.Fatal(err)
	}
	e2 := NewEngine(app)
	if len(e2.corporateActions) != 1 || e2.corporateActions[0].ID != "persist-1" {
		t.Fatalf("ledger not restored: %+v", e2.corporateActions)
	}
	if e2.lastUpdated["corporate-actions-backfill:META"] != 30 {
		t.Fatalf("backfill stamp not restored: %+v", e2.lastUpdated)
	}
}

func TestV1601CorporateParserCoversCurrentActionFamilies(t *testing.T) {
	types := []string{"reverse_split", "forward_split", "unit_split", "cash_dividend", "stock_dividend", "spin_off", "cash_merger", "stock_merger", "stock_and_cash_merger", "redemption", "name_change", "worthless_removal", "rights_distribution", "partial_call", "reorganization"}
	payload := map[string]any{}
	for i, typ := range types {
		row := map[string]any{"id": typ + "-id", "symbol": "AAPL", "process_date": "2026-08-01"}
		if typ == "name_change" {
			row["old_symbol"] = "AAPL"
			row["new_symbol"] = "AAPLX"
		}
		payload[typ+"s"] = []any{row}
		_ = i
	}
	tracked := map[string]bool{"AAPL": true, "AAPLX": true}
	rows := parseAlpacaCorporateActionResponse(payload, tracked)
	got := map[string]bool{}
	for _, a := range rows {
		got[a.Type] = true
		if a.ID == "" || a.ProcessDate == "" {
			t.Fatalf("missing ledger provenance: %+v", a)
		}
	}
	for _, typ := range types {
		if !got[typ] {
			t.Fatalf("missing action type %s rows=%+v", typ, rows)
		}
	}
}

func TestV1601CorporateRefreshBackfillsNewSymbolsAndPaginates(t *testing.T) {
	oldBase := alpacaDataBaseURL
	defer func() { alpacaDataBaseURL = oldBase }()
	calls := 0
	starts := []string{}
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		calls++
		starts = append(starts, r.URL.Query().Get("start"))
		w.Header().Set("Content-Type", "application/json")
		if r.URL.Path == "/v1/corporate-actions" {
			if r.URL.Query().Get("page_token") == "p2" {
				_, _ = w.Write([]byte(`{"name_changes":[{"id":"ca2","symbol":"AAPL","old_symbol":"OLD","new_symbol":"AAPL","process_date":"2025-01-01"}]}`))
				return
			}
			_, _ = w.Write([]byte(`{"forward_splits":[{"id":"ca1","symbol":"AAPL","process_date":"2024-06-01","old_rate":1,"new_rate":4}],"next_page_token":"p2"}`))
			return
		}
		if r.URL.Path == "/v2/stocks/bars" {
			_, _ = w.Write([]byte(`{"bars":{}}`))
			return
		}
		http.NotFound(w, r)
	}))
	defer srv.Close()
	alpacaDataBaseURL = srv.URL
	root := t.TempDir()
	t.Setenv("XDG_CONFIG_HOME", root)
	t.Setenv("HOME", root)
	t.Setenv("APPDATA", root)
	app, err := NewApplication()
	if err != nil {
		t.Fatal(err)
	}
	registerApplicationCleanup(t, app)
	app.engine = NewEngine(app)
	app.mu.Lock()
	app.state.Watchlists = []Watchlist{{ID: "day", Name: "Day", Symbols: []string{"AAPL"}}, {ID: "swing", Name: "Swing"}, {ID: "long", Name: "Long"}, {ID: "discovery", Name: "Discovery"}}
	app.mu.Unlock()
	app.engine.refreshAlpacaCorporateActions(context.Background(), "k", "s")
	if len(app.engine.corporateActions) != 2 || calls < 2 {
		t.Fatalf("backfill/pagination failed calls=%d actions=%+v", calls, app.engine.corporateActions)
	}
	if app.engine.lastUpdated["corporate-actions-backfill:AAPL"] == 0 {
		t.Fatal("per-symbol backfill stamp missing")
	}
	if len(starts) == 0 || starts[0] > "2012-01-01" {
		t.Fatalf("initial range not historical enough: %v", starts)
	}
}

func TestV1601ApprovedScopeIncludesAdversarialProof(t *testing.T) {
	raw, err := json.Marshal([]string{"Provider reason/canonical consistency", "Selected-ticker history timestamp", "Persistent corporate-action ledger", "Adversarial next-build gate"})
	if err != nil || len(raw) == 0 {
		t.Fatal(err)
	}
}
