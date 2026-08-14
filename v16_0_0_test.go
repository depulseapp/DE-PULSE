package main

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestV1600IdentityUsesVersionAndRuntimeContract(t *testing.T) {
	if appVersion == "" {
		t.Fatal("runtime version identity must not be empty")
	}
	if buildID == "" {
		t.Fatal("runtime build identity must not be empty")
	}
}

func TestV180ApplicationUsesReleaseChannelConfigIsolation(t *testing.T) {
	base := isolateUserConfig(t)
	app, err := NewApplication()
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = app.persistence.Close() })
	want := filepath.Join(base, stableRuntimeConfigDirName)
	if releaseChannel != "STABLE" {
		want = filepath.Join(base, v18TestRuntimeConfigDirName)
	}
	if app.configDir != want {
		t.Fatalf("v18 %s configDir=%s want=%s", releaseChannel, app.configDir, want)
	}
}

func TestV180TestProfileHelperMigratesStableSecretsWithoutMutatingStable(t *testing.T) {
	base := isolateUserConfig(t)
	dir := filepath.Join(base, stableRuntimeConfigDirName)
	if err := os.MkdirAll(dir, 0700); err != nil {
		t.Fatal(err)
	}
	want := Secrets{Finnhub: "fin-old", AlpacaKey: "alp-key", AlpacaSecret: "alp-secret", Groq: "groq-old", OpenRouter: "or-old", Gemini: "gem-old", FRED: "fred-old", BLS: "bls-old", EIA: "eia-old", TwelveData: "td-old", Marketaux: "mx-old"}
	b, _ := json.Marshal(want)
	if err := os.WriteFile(filepath.Join(dir, "secrets.json"), b, 0600); err != nil {
		t.Fatal(err)
	}
	target, err := prepareV18TestConfig(base)
	if err != nil {
		t.Fatal(err)
	}
	if target != filepath.Join(base, v18TestRuntimeConfigDirName) {
		t.Fatalf("unexpected v18 TEST helper config: %s", target)
	}
	targetBytes, err := os.ReadFile(filepath.Join(target, "secrets.json"))
	if err != nil || !bytes.Equal(targetBytes, b) {
		t.Fatalf("prior Stable secrets not migrated by TEST helper: err=%v", err)
	}
	sourceBytes, err := os.ReadFile(filepath.Join(dir, "secrets.json"))
	if err != nil || !bytes.Equal(sourceBytes, b) {
		t.Fatalf("Stable secrets were modified during TEST helper migration: err=%v", err)
	}
}

func testRouter(active string) ProviderRouterSnapshot {
	return ProviderRouterSnapshot{Routes: []ProviderRouteState{{Dataset: "US Live Equities", Primary: "Alpaca", Active: active, State: "READY", Route: []ProviderRouteHop{{Provider: "Alpaca", Priority: 1}, {Provider: "Finnhub", Priority: 2}, {Provider: "Twelve Data", Priority: 3}}}}}
}

func TestV1600ProviderReconciliationUsesRouterWinnerNotPriceOrder(t *testing.T) {
	now := time.Now().UnixMilli()
	obs := map[string]map[string]Quote{"AAPL": {
		"Alpaca":  {Symbol: "AAPL", Price: 200.00, ProviderTimestamp: now - 200, UpdatedAt: now - 100, Source: "alpaca-iex-websocket-trade"},
		"Finnhub": {Symbol: "AAPL", Price: 201.00, ProviderTimestamp: now - 150, UpdatedAt: now - 90, Source: "finnhub-websocket"},
	}}
	canonical := map[string]Quote{"AAPL": {Symbol: "AAPL", Price: 200.00, Source: "alpaca-iex-websocket-trade", UpdatedAt: now}}
	rows := buildProviderReconciliation(testRouter("Alpaca"), obs, canonical, now)
	if len(rows) != 1 {
		t.Fatalf("rows=%+v", rows)
	}
	r := rows[0]
	if r.State != "CONFLICT" || r.CanonicalProvider != "Alpaca" || r.CanonicalValue != 200 {
		t.Fatalf("bad reconciliation: %+v", r)
	}
	if !strings.Contains(r.Reason, "Canonical winner is Alpaca") {
		t.Fatalf("reason=%s", r.Reason)
	}
}

func TestV1600ProviderReconciliationDoesNotUseNonContemporaneousObservationAsConflict(t *testing.T) {
	// Saturday deliberately exercises the long closed-market freshness allowance.
	// A prior-close observation can remain usable, but a provider snapshot from a
	// different market moment must not manufacture a cross-provider conflict.
	now := time.Date(2026, 8, 15, 16, 0, 0, 0, time.UTC).UnixMilli()
	fresh := time.Date(2026, 8, 14, 20, 0, 0, 0, time.UTC).UnixMilli()
	obs := map[string]map[string]Quote{"AAPL": {
		"Alpaca":  {Symbol: "AAPL", Price: 200, ProviderTimestamp: fresh, UpdatedAt: fresh, Source: "alpaca-iex-websocket-trade"},
		"Finnhub": {Symbol: "AAPL", Price: 210, ProviderTimestamp: fresh - int64(5*time.Minute/time.Millisecond), UpdatedAt: fresh - int64(5*time.Minute/time.Millisecond), Source: "finnhub-websocket"},
	}}
	rows := buildProviderReconciliation(testRouter("Alpaca"), obs, nil, now)
	if len(rows) != 1 || rows[0].State != "SINGLE SOURCE" || len(rows[0].Observations) != 1 {
		t.Fatalf("non-contemporaneous observation manufactured conflict: %+v", rows)
	}
}

func TestV1600ProviderReconciliationAllowsContemporaneousPriorCloseEvidence(t *testing.T) {
	now := time.Date(2026, 8, 15, 16, 0, 0, 0, time.UTC).UnixMilli()
	closeAt := time.Date(2026, 8, 14, 20, 0, 0, 0, time.UTC).UnixMilli()
	obs := map[string]map[string]Quote{"AAPL": {
		"Alpaca":  {Symbol: "AAPL", Price: 200, ProviderTimestamp: closeAt, UpdatedAt: closeAt, Source: "alpaca-iex-websocket-trade"},
		"Finnhub": {Symbol: "AAPL", Price: 200.05, ProviderTimestamp: closeAt - 30_000, UpdatedAt: closeAt - 29_000, Source: "finnhub-websocket"},
	}}
	rows := buildProviderReconciliation(testRouter("Alpaca"), obs, nil, now)
	if len(rows) != 1 || rows[0].State != "AGREED" || len(rows[0].Observations) != 2 {
		t.Fatalf("contemporaneous prior-close observations should reconcile: %+v", rows)
	}
}

func completeResearchFixture(now time.Time) (AppState, map[string]Quote, map[string]map[string][]Bar, map[string]FundamentalSnapshot, map[string]int64, map[string]string) {
	st := defaultState()
	st.UI.SelectedTicker = "AAPL"
	ms := now.UnixMilli()
	quotes := map[string]Quote{"AAPL": {Symbol: "AAPL", Price: 200, Source: "alpaca-iex-websocket-trade", ProviderTimestamp: ms - 1000, UpdatedAt: ms - 800}}
	bars := map[string]map[string][]Bar{"AAPL": {"daily": {{T: now.Add(-24 * time.Hour).Unix(), O: 195, H: 202, L: 194, C: 200, V: 1_000_000}}}}
	fundamentals := map[string]FundamentalSnapshot{"AAPL": {Symbol: "AAPL", MarketCap: 3e12, PERatio: 30, UpdatedAt: ms - 1000, Source: "Finnhub"}}
	last := map[string]int64{"research-history:AAPL": ms - 10_000, "research-fundamentals:AAPL": ms - 10_000, "research-news:AAPL": ms - 10_000, "research-earnings:AAPL": ms - 10_000, "research-sec:AAPL": ms - 10_000}
	health := map[string]string{"research-sec:AAPL": "healthy"}
	return st, quotes, bars, fundamentals, last, health
}

func TestV1600ResearchPackageIsSelectedTickerSpecific(t *testing.T) {
	now := time.Date(2026, 8, 11, 14, 0, 0, 0, time.UTC)
	st, quotes, bars, fundamentals, last, health := completeResearchFixture(now)
	rec := []ProviderReconciliationDecision{{Dataset: "US Live Equities", Symbol: "AAPL", State: "AGREED", CanonicalProvider: "Alpaca", CanonicalValue: 200, UpdatedAt: now.UnixMilli(), Observations: []ProviderQuoteObservation{{Symbol: "AAPL", Provider: "Alpaca", Price: 200, ProviderTimestamp: now.UnixMilli() - 1000, ReceivedAt: now.UnixMilli() - 800}, {Symbol: "AAPL", Provider: "Finnhub", Price: 200.05, ProviderTimestamp: now.UnixMilli() - 900, ReceivedAt: now.UnixMilli() - 700}}, Reason: "two current providers agree"}}
	pkg := buildResearchPackageTruth(st, quotes, bars, fundamentals, nil, nil, nil, last, health, rec, now)
	if pkg.State != "FRESH" {
		t.Fatalf("expected fresh target package, got %+v", pkg)
	}
	// Global dataset success must not make a selected ticker current.
	delete(last, "research-news:AAPL")
	last["news"] = now.UnixMilli()
	pkg = buildResearchPackageTruth(st, quotes, bars, fundamentals, nil, nil, nil, last, health, rec, now)
	if pkg.State != "PARTIAL" {
		t.Fatalf("global News check falsely satisfied target: %+v", pkg)
	}
	found := false
	for _, c := range pkg.Components {
		if c.Dataset == "News" && c.State == "PARTIAL" {
			found = true
		}
	}
	if !found {
		t.Fatal("selected-ticker News component did not expose partial state")
	}
}

func TestV1600ResearchPackageBlocksMissingCriticalTickerEvidence(t *testing.T) {
	now := time.Date(2026, 8, 11, 14, 0, 0, 0, time.UTC)
	st, quotes, bars, fundamentals, last, health := completeResearchFixture(now)
	delete(quotes, "AAPL")
	pkg := buildResearchPackageTruth(st, quotes, bars, fundamentals, nil, nil, nil, last, health, nil, now)
	if pkg.State != "BLOCKED" {
		t.Fatalf("missing quote must block: %+v", pkg)
	}
}

func TestV1600ResearchPackageInheritsMaterialProviderConflict(t *testing.T) {
	now := time.Date(2026, 8, 11, 14, 0, 0, 0, time.UTC)
	st, quotes, bars, fundamentals, last, health := completeResearchFixture(now)
	rec := []ProviderReconciliationDecision{{Dataset: "US Live Equities", Symbol: "AAPL", State: "CONFLICT", DifferencePct: 1.2, Reason: "Independent prices materially disagree"}}
	pkg := buildResearchPackageTruth(st, quotes, bars, fundamentals, nil, nil, nil, last, health, rec, now)
	if pkg.State != "BLOCKED" {
		t.Fatalf("large selected-ticker conflict must block: %+v", pkg)
	}
}

func TestV1600EvidenceSnapshotStableAcrossClockAging(t *testing.T) {
	pkg := ResearchPackageTruth{Symbol: "AAPL", State: "FRESH", GeneratedAt: 1000, Components: []ResearchEvidenceComponent{{Dataset: "Quote", Symbol: "AAPL", State: "FRESH", Required: true, Critical: true, Source: "Alpaca", CheckAt: 900, DataAt: 800, CheckAgeMs: 100, DataAgeMs: 200, Detail: "current"}}}
	a := buildEvidenceSnapshot(pkg, nil, CorporateActionTruth{})
	pkg.GeneratedAt = 9000
	pkg.Components[0].CheckAgeMs = 8100
	pkg.Components[0].DataAgeMs = 8200
	b := buildEvidenceSnapshot(pkg, nil, CorporateActionTruth{})
	if a.ID != b.ID {
		t.Fatalf("clock aging changed evidence ID: %s vs %s", a.ID, b.ID)
	}
	pkg.Components[0].State = "STALE"
	c := buildEvidenceSnapshot(pkg, nil, CorporateActionTruth{})
	if c.ID == a.ID {
		t.Fatal("material evidence state change did not change evidence ID")
	}
}

func TestV1600CorporateActionParserCapturesSplitAndSymbolLineage(t *testing.T) {
	payload := map[string]any{
		"forward_splits": []any{map[string]any{"symbol": "NVDA", "ex_date": "2024-06-10", "old_rate": 1.0, "new_rate": 10.0}},
		"name_changes":   []any{map[string]any{"old_symbol": "FB", "new_symbol": "META", "effective_date": "2022-06-09"}},
	}
	tracked := map[string]bool{"NVDA": true, "FB": true, "META": true}
	rows := parseAlpacaCorporateActionResponse(payload, tracked)
	if len(rows) != 2 {
		b, _ := json.Marshal(rows)
		t.Fatalf("actions=%s", b)
	}
	var split CorporateAction
	for _, r := range rows {
		if r.Symbol == "NVDA" {
			split = r
		}
	}
	if split.Ratio != 10 || split.AdjustmentFactor != 10 {
		t.Fatalf("split=%+v", split)
	}
	truth := buildCorporateActionTruth(rows, map[string]map[string][]Bar{"NVDA": {"daily-raw": {{T: 1, C: 100}}}}, time.Now().UnixMilli(), map[string]RawHistoryCoverage{"NVDA": {Symbol: "NVDA", State: "COMPLETE", BarCount: 1, PageCount: 1, PaginationComplete: true}})
	if truth.SymbolLineage["FB"] != "META" {
		t.Fatalf("lineage=%+v actions=%+v", truth.SymbolLineage, rows)
	}
	if !truth.RawHistoryAvailable["NVDA"] {
		t.Fatal("raw-history provenance not exposed")
	}
}

func TestV1600RawCorporateHistoryUsesRawAdjustmentAndSeparateStore(t *testing.T) {
	oldBase := alpacaDataBaseURL
	defer func() { alpacaDataBaseURL = oldBase }()
	var query string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		query = r.URL.RawQuery
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"bars":{"NVDA":[{"c":100,"h":101,"l":99,"o":100,"v":1000,"t":"2026-08-10T20:00:00Z"}]}}`))
	}))
	defer srv.Close()
	alpacaDataBaseURL = srv.URL
	st := defaultState()
	app := &Application{state: st, configDir: t.TempDir(), hub: NewHub(), sessionKey: "x"}
	app.engine = NewEngine(app)
	n := app.engine.refreshAlpacaRawHistoryForCorporateActions(context.Background(), "k", "s", []CorporateAction{{Symbol: "NVDA", Type: "forward_split"}})
	if n != 1 {
		t.Fatalf("loaded=%d", n)
	}
	if !strings.Contains(query, "adjustment=raw") {
		t.Fatalf("query=%s", query)
	}
	if len(app.engine.bars["NVDA"]["daily-raw"]) != 1 {
		t.Fatalf("bars=%+v", app.engine.bars["NVDA"])
	}
}

func TestV180RuntimeUsesReleaseChannelResolver(t *testing.T) {
	s := productionGoSourceForTest(t)
	if !strings.Contains(s, `resolveV18RuntimeConfig(base)`) {
		t.Fatal("v18 release-channel runtime resolver missing")
	}
}
