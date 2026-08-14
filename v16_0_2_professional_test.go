package main

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func currentQuote(sym, source string, price float64, now int64) Quote {
	return Quote{Symbol: sym, Source: source, Price: price, ProviderTimestamp: now - 500, UpdatedAt: now - 250}
}

func TestV1602ProfessionalNoEvidenceNeverMeansAgreement(t *testing.T) {
	now := time.Now().UnixMilli()
	rows := buildProviderReconciliation(testRouter("Alpaca"), nil, nil, now)
	if len(rows) != 0 {
		t.Fatalf("no symbols should not manufacture reconciliation rows: %+v", rows)
	}
	st, quotes, bars, fundamentals, last, health := completeResearchFixture(time.UnixMilli(now))
	pkg := buildResearchPackageTruth(st, quotes, bars, fundamentals, nil, nil, nil, last, health, nil, time.UnixMilli(now))
	for _, c := range pkg.Components {
		if c.Dataset == "Provider Reconciliation" {
			if c.State == "FRESH" || !strings.Contains(strings.ToLower(c.Detail), "no current") {
				t.Fatalf("absence of comparison evidence was treated as healthy: %+v", c)
			}
			return
		}
	}
	t.Fatal("Provider Reconciliation component missing")
}

func TestV1602ProfessionalSingleSourceIsNotCrossProviderAgreement(t *testing.T) {
	now := time.Now().UnixMilli()
	canonical := map[string]Quote{"AAPL": currentQuote("AAPL", "alpaca-iex-websocket-trade", 200, now)}
	rows := buildProviderReconciliation(testRouter("Alpaca"), nil, canonical, now)
	if len(rows) != 1 || rows[0].State != "SINGLE SOURCE" || len(rows[0].Observations) != 1 {
		t.Fatalf("single canonical source must remain SINGLE SOURCE: %+v", rows)
	}
	if strings.Contains(strings.ToLower(rows[0].Reason), "providers agree") || strings.Contains(strings.ToLower(rows[0].Reason), "observations agree") {
		t.Fatalf("single source made an agreement claim: %s", rows[0].Reason)
	}
}

func TestV1602ProfessionalMarketAgeAndReceiptAgeBothMustBeCurrent(t *testing.T) {
	now := time.Date(2026, 8, 11, 15, 30, 0, 0, time.UTC).UnixMilli() // regular session
	obs := map[string]map[string]Quote{"AAPL": {
		"Alpaca":      {Symbol: "AAPL", Price: 200, Source: "alpaca", ProviderTimestamp: now - int64(10*time.Minute/time.Millisecond), UpdatedAt: now - 100},
		"Finnhub":     {Symbol: "AAPL", Price: 200.1, Source: "finnhub", ProviderTimestamp: now - 200, UpdatedAt: now - int64(10*time.Minute/time.Millisecond)},
		"Twelve Data": {Symbol: "AAPL", Price: 200.2, Source: "twelve-data", ProviderTimestamp: now - 500, UpdatedAt: now - 400},
	}}
	rows := buildProviderReconciliation(testRouter("Twelve Data"), obs, nil, now)
	if len(rows) != 1 || len(rows[0].Observations) != 1 || rows[0].Observations[0].Provider != "Twelve Data" {
		t.Fatalf("stale market/transport observations participated as current: %+v", rows)
	}
}

func TestV1602EdgeClockSkewCannotCreateCurrentEvidence(t *testing.T) {
	now := time.Date(2026, 8, 11, 15, 30, 0, 0, time.UTC).UnixMilli()
	obs := map[string]map[string]Quote{"AAPL": {
		"Alpaca": {Symbol: "AAPL", Price: 200, Source: "alpaca", ProviderTimestamp: now + int64(5*time.Minute/time.Millisecond), UpdatedAt: now - 100},
	}}
	rows := buildProviderReconciliation(testRouter("Alpaca"), obs, nil, now)
	if len(rows) != 1 || rows[0].State != "STALE" || len(rows[0].Observations) != 0 {
		t.Fatalf("future-skewed provider timestamp treated as current: %+v", rows)
	}
}

func TestV1602EdgeClosedMarketAllowsTruthfulRecentCloseEvidence(t *testing.T) {
	now := time.Date(2026, 8, 15, 16, 0, 0, 0, time.UTC).UnixMilli() // Saturday
	friday := time.Date(2026, 8, 14, 20, 0, 0, 0, time.UTC).UnixMilli()
	canonical := map[string]Quote{"AAPL": {Symbol: "AAPL", Price: 200, Source: "alpaca", ProviderTimestamp: friday, UpdatedAt: friday}}
	rows := buildProviderReconciliation(testRouter("Alpaca"), nil, canonical, now)
	if len(rows) != 1 || rows[0].State != "SINGLE SOURCE" {
		t.Fatalf("valid recent closed-market evidence was incorrectly stale: %+v", rows)
	}
}

func TestV1602ProfessionalResearchSingleSourceDegradesConfidence(t *testing.T) {
	now := time.Date(2026, 8, 11, 14, 0, 0, 0, time.UTC)
	st, quotes, bars, fundamentals, last, health := completeResearchFixture(now)
	rec := []ProviderReconciliationDecision{{Dataset: "US Live Equities", Symbol: "AAPL", State: "SINGLE SOURCE", CanonicalProvider: "Alpaca", CanonicalValue: 200, UpdatedAt: now.UnixMilli(), Observations: []ProviderQuoteObservation{{Symbol: "AAPL", Provider: "Alpaca", Price: 200, ProviderTimestamp: now.UnixMilli() - 500, ReceivedAt: now.UnixMilli() - 300}}, Reason: "Only one current provider observation is available."}}
	pkg := buildResearchPackageTruth(st, quotes, bars, fundamentals, nil, nil, nil, last, health, rec, now)
	if pkg.State != "PARTIAL" {
		t.Fatalf("single-source research should be PARTIAL rather than fully cross-checked: %+v", pkg)
	}
}

func TestV1602ProfessionalMaterialConflictBlocksResearch(t *testing.T) {
	now := time.Date(2026, 8, 11, 14, 0, 0, 0, time.UTC)
	st, quotes, bars, fundamentals, last, health := completeResearchFixture(now)
	rec := []ProviderReconciliationDecision{{Dataset: "US Live Equities", Symbol: "AAPL", State: "CONFLICT", DifferencePct: 1.25, CanonicalProvider: "Alpaca", CanonicalValue: 200, UpdatedAt: now.UnixMilli(), Reason: "Material independent provider conflict."}}
	pkg := buildResearchPackageTruth(st, quotes, bars, fundamentals, nil, nil, nil, last, health, rec, now)
	if pkg.State != "BLOCKED" {
		t.Fatalf("material provider conflict must block decision-grade research: %+v", pkg)
	}
}

func TestV1602EdgeRawHistoryPartialFailureExposesCoverageAndRecovers(t *testing.T) {
	oldBase := alpacaDataBaseURL
	defer func() { alpacaDataBaseURL = oldBase }()
	failSecond := true
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v2/stocks/bars" {
			http.NotFound(w, r)
			return
		}
		if r.URL.Query().Get("page_token") == "p2" {
			if failSecond {
				http.Error(w, "temporary failure", 500)
				return
			}
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`{"bars":{"AAPL":[{"o":101,"h":102,"l":100,"c":101,"v":11,"t":"2026-08-11T00:00:00Z"}]}}`))
			return
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"bars":{"AAPL":[{"o":100,"h":101,"l":99,"c":100,"v":10,"t":"2026-08-10T00:00:00Z"}]},"next_page_token":"p2"}`))
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
	actions := []CorporateAction{{ID: "ca", Symbol: "AAPL", Type: "forward_split"}}
	app.engine.refreshAlpacaRawHistoryForCorporateActions(context.Background(), "k", "s", actions)
	cov := app.engine.rawHistoryCoverage["AAPL"]
	if cov.State != "PARTIAL" || cov.PaginationComplete || cov.BarCount != 1 || !strings.Contains(strings.ToLower(app.engine.health["corporate-actions-raw-history"]), "degraded") {
		t.Fatalf("partial page failure was not surfaced truthfully: cov=%+v health=%q", cov, app.engine.health["corporate-actions-raw-history"])
	}
	failSecond = false
	app.engine.refreshAlpacaRawHistoryForCorporateActions(context.Background(), "k", "s", actions)
	cov = app.engine.rawHistoryCoverage["AAPL"]
	if cov.State != "COMPLETE" || !cov.PaginationComplete || cov.BarCount != 2 || len(app.engine.bars["AAPL"]["daily-raw"]) != 2 || !strings.Contains(strings.ToLower(app.engine.health["corporate-actions-raw-history"]), "healthy") {
		t.Fatalf("successful recovery did not close raw-history coverage: cov=%+v health=%q bars=%+v", cov, app.engine.health["corporate-actions-raw-history"], app.engine.bars["AAPL"]["daily-raw"])
	}
}

func TestV1602EdgeRawHistoryNoBarsIsUnavailableNotComplete(t *testing.T) {
	oldBase := alpacaDataBaseURL
	defer func() { alpacaDataBaseURL = oldBase }()
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"bars":{}}`))
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
	app.engine.refreshAlpacaRawHistoryForCorporateActions(context.Background(), "k", "s", []CorporateAction{{Symbol: "AAPL", Type: "forward_split"}})
	cov := app.engine.rawHistoryCoverage["AAPL"]
	if cov.State != "UNAVAILABLE" || app.engine.bars["AAPL"]["daily-raw"] != nil {
		t.Fatalf("empty provider history falsely reported complete: %+v", cov)
	}
}

func TestV1602EdgeRawHistoryCoveragePersistsAcrossRestart(t *testing.T) {
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
	app.engine.rawHistoryCoverage["AAPL"] = RawHistoryCoverage{Symbol: "AAPL", State: "COMPLETE", BarCount: 1234, PageCount: 2, PaginationComplete: true, FirstBarAt: 1, LastBarAt: 2, UpdatedAt: 3}
	app.engine.mu.Unlock()
	if err := app.engine.saveCache(); err != nil {
		t.Fatal(err)
	}
	e2 := NewEngine(app)
	cov, ok := e2.rawHistoryCoverage["AAPL"]
	if !ok || cov.State != "COMPLETE" || cov.BarCount != 1234 || !cov.PaginationComplete {
		t.Fatalf("raw-history coverage provenance lost across restart: %+v", e2.rawHistoryCoverage)
	}
}

func TestV1602ProfessionalCorporateTruthOnlyCallsRawCompleteWithProof(t *testing.T) {
	actions := []CorporateAction{{ID: "x", Symbol: "AAPL", Type: "forward_split"}}
	bars := map[string]map[string][]Bar{"AAPL": {"daily-raw": {{T: 1, C: 100}}}}
	partial := buildCorporateActionTruth(actions, bars, 10, map[string]RawHistoryCoverage{"AAPL": {Symbol: "AAPL", State: "PARTIAL", BarCount: 1, PaginationComplete: false}})
	if partial.RawHistoryAvailable["AAPL"] {
		t.Fatal("partial raw history exposed as complete")
	}
	complete := buildCorporateActionTruth(actions, bars, 10, map[string]RawHistoryCoverage{"AAPL": {Symbol: "AAPL", State: "COMPLETE", BarCount: 1, PageCount: 1, PaginationComplete: true}})
	if !complete.RawHistoryAvailable["AAPL"] {
		t.Fatal("complete paginated raw history not exposed")
	}
}

func TestV1602ProfessionalEvidenceSnapshotTracksCorporateCoverageTruth(t *testing.T) {
	pkg := ResearchPackageTruth{Symbol: "AAPL", State: "FRESH", GeneratedAt: 1, Components: []ResearchEvidenceComponent{{Dataset: "Quote", Symbol: "AAPL", State: "FRESH", Required: true}}}
	corp := CorporateActionTruth{RawHistoryCoverage: map[string]RawHistoryCoverage{"AAPL": {Symbol: "AAPL", State: "PARTIAL", BarCount: 10, PaginationComplete: false}}}
	a := buildEvidenceSnapshot(pkg, nil, corp)
	corp.RawHistoryCoverage["AAPL"] = RawHistoryCoverage{Symbol: "AAPL", State: "COMPLETE", BarCount: 100, PageCount: 2, PaginationComplete: true}
	b := buildEvidenceSnapshot(pkg, nil, corp)
	if a.ID == b.ID || len(b.CorporateHistoryCoverage) != 1 || b.CorporateHistoryCoverage[0].State != "COMPLETE" {
		t.Fatalf("material corporate-history truth did not change evidence snapshot: a=%+v b=%+v", a, b)
	}
}

func TestV1602EdgeCorporateActionPageFailureCannotEndHealthy(t *testing.T) {
	oldBase := alpacaDataBaseURL
	defer func() { alpacaDataBaseURL = oldBase }()
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v1/corporate-actions" {
			http.NotFound(w, r)
			return
		}
		if r.URL.Query().Get("page_token") == "p2" {
			http.Error(w, "temporary page failure", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"forward_splits":[{"id":"ca1","symbol":"AAPL","process_date":"2024-06-01","old_rate":1,"new_rate":4}],"next_page_token":"p2"}`))
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
	app.state.Watchlists = []Watchlist{{ID: "day", Symbols: []string{"AAPL"}}, {ID: "swing"}, {ID: "long"}, {ID: "discovery"}}
	app.mu.Unlock()
	app.engine.refreshAlpacaCorporateActions(context.Background(), "k", "s")
	h := strings.ToLower(app.engine.health["corporate-actions"])
	if strings.Contains(h, "healthy") || !strings.Contains(h, "partial") {
		t.Fatalf("partial corporate-action pagination became healthy: %q", h)
	}
}

func TestV1602EdgeCorporateActionEntitlementFailureCannotReportHealthy(t *testing.T) {
	oldBase := alpacaDataBaseURL
	defer func() { alpacaDataBaseURL = oldBase }()
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { http.Error(w, "forbidden", http.StatusForbidden) }))
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
	app.state.Watchlists = []Watchlist{{ID: "day", Symbols: []string{"AAPL"}}, {ID: "swing"}, {ID: "long"}, {ID: "discovery"}}
	app.mu.Unlock()
	app.engine.refreshAlpacaCorporateActions(context.Background(), "k", "s")
	h := strings.ToLower(app.engine.health["corporate-actions"])
	if strings.Contains(h, "healthy") || !strings.Contains(h, "not entitled") {
		t.Fatalf("entitlement failure became healthy/opaque: %q", h)
	}
}
