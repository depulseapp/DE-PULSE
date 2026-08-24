package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"sync/atomic"
	"testing"
	"time"
)

func TestV189HistoryRouterUsesOneCanonicalDataset(t *testing.T) {
	routes := routeChains()
	chain := routes[canonicalHistoricalBarsDataset]
	if len(chain) == 0 {
		t.Fatalf("canonical %q route is missing", canonicalHistoricalBarsDataset)
	}
	if _, exists := routes["Intraday Bars"]; exists {
		t.Fatal("intraday history must not create a second router dataset")
	}
	if _, exists := routes["Daily / Weekly History"]; exists {
		t.Fatal("daily history must not create a second router dataset")
	}
	seen := false
	for _, provider := range chain {
		if provider == tradeInsightProviderName {
			seen = true
			break
		}
	}
	if !seen {
		t.Fatal("TradeInsight must be a member of the canonical Historical Bars route")
	}
}

func TestV189TradeInsightConfiguredFromSecretEnvironment(t *testing.T) {
	t.Setenv("TIDATA_API_KEY", "")
	t.Setenv("TRADEINSIGHT_API_KEY", "")
	if tradeInsightConfigured() {
		t.Fatal("TradeInsight must be NOT_CONFIGURED when no runtime secret is present")
	}
	t.Setenv("TIDATA_API_KEY", "ti_fixture_secret")
	if !tradeInsightConfigured() {
		t.Fatal("TIDATA_API_KEY should configure TradeInsight")
	}
	if !(&Engine{}).providerConfigured(tradeInsightProviderName, Secrets{}, Settings{}) {
		t.Fatal("Smart Provider Router must see the TradeInsight runtime secret")
	}
}

func TestV189TradeInsightDailyOnly(t *testing.T) {
	t.Setenv("TIDATA_API_KEY", "ti_fixture_secret")
	var e *Engine
	if got := e.refreshTradeInsightHistoryMode(context.Background(), []string{"SPY"}, "intraday"); got != 0 {
		t.Fatalf("TradeInsight must not claim intraday support, got %d bars", got)
	}
	if got := e.refreshTradeInsightHistoryMode(context.Background(), []string{"SPY"}, "all"); got != 0 {
		t.Fatalf("TradeInsight must not satisfy mixed all-mode history with daily-only data, got %d bars", got)
	}
}

func TestV189TradeInsightPaginationAuthAndEnvelope(t *testing.T) {
	const key = "ti_fixture_secret"
	var calls atomic.Int32
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		calls.Add(1)
		if got := r.Header.Get("Authorization"); got != "Bearer "+key {
			t.Fatalf("authorization header = %q", got)
		}
		offset := r.URL.Query().Get("offset")
		w.Header().Set("Content-Type", "application/json")
		if offset == "0" {
			rows := make([]map[string]any, tradeInsightPageSize)
			for i := range rows {
				rows[i] = map[string]any{"date": "2026-08-20", "adj_close": 100 + i}
			}
			_ = json.NewEncoder(w).Encode(map[string]any{"data": rows})
			return
		}
		_ = json.NewEncoder(w).Encode(map[string]any{"data": []map[string]any{{"date": "2026-08-21", "adj_close": 200.0}}})
	}))
	defer server.Close()

	rows, err := tradeInsightFetchRowsAt(context.Background(), server.Client(), server.URL, key, "/ohlc", url.Values{"ticker": []string{"SPY"}})
	if err != nil {
		t.Fatal(err)
	}
	if len(rows) != tradeInsightPageSize+1 {
		t.Fatalf("rows = %d, want %d", len(rows), tradeInsightPageSize+1)
	}
	if calls.Load() != 2 {
		t.Fatalf("calls = %d, want 2", calls.Load())
	}
}

func TestV189TradeInsightErrorsAreRouterClassifiableAndRedacted(t *testing.T) {
	const key = "ti_do_not_leak"
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Retry-After", "17")
		w.WriteHeader(http.StatusTooManyRequests)
		_, _ = fmt.Fprintf(w, `{"message":"rate limit for %s"}`, key)
	}))
	defer server.Close()

	_, err := tradeInsightFetchRowsAt(context.Background(), server.Client(), server.URL, key, "/ohlc", url.Values{"ticker": []string{"SPY"}})
	if err == nil {
		t.Fatal("expected rate-limit error")
	}
	msg := err.Error()
	if !strings.Contains(strings.ToLower(msg), "http 429") || !strings.Contains(strings.ToLower(msg), "rate limit") {
		t.Fatalf("error is not classifiable by Smart Provider Router: %q", msg)
	}
	if strings.Contains(msg, key) {
		t.Fatalf("API key leaked in error: %q", msg)
	}
	if !strings.Contains(msg, "retry-after=17") {
		t.Fatalf("Retry-After evidence missing: %q", msg)
	}
}

func TestV189TradeInsightAdjustedRowsAndDerivedWeekly(t *testing.T) {
	rows := []map[string]any{
		{"date": "2026-08-17", "open": 9.0, "high": 11.0, "low": 8.0, "close": 10.0, "volume": 90.0, "adj_open": 99.0, "adj_high": 101.0, "adj_low": 98.0, "adj_close": 100.0, "adj_volume": 900.0, "dividend": 0.5, "split_ratio": 0.0},
		{"date": "2026-08-18", "open": 10.0, "high": 12.0, "low": 9.0, "close": 11.0, "volume": 100.0, "adj_open": 100.0, "adj_high": 102.0, "adj_low": 99.0, "adj_close": 101.0, "adj_volume": 1000.0, "dividend": 0.0, "split_ratio": 2.0},
	}
	bars := tradeInsightBars(rows, true)
	if len(bars) != 2 {
		t.Fatalf("bars = %d", len(bars))
	}
	if bars[0].O != 99 || bars[1].C != 101 || bars[1].V != 1000 {
		t.Fatalf("adjusted fields were not used: %+v", bars)
	}
	weekly := aggregateDailyBarsToWeekly(bars)
	if len(weekly) != 1 {
		t.Fatalf("weekly bars = %d, want 1", len(weekly))
	}
	if weekly[0].O != 99 || weekly[0].H != 102 || weekly[0].L != 98 || weekly[0].C != 101 || weekly[0].V != 1900 {
		t.Fatalf("derived weekly aggregation is wrong: %+v", weekly[0])
	}
	parsed := tradeInsightHistoryRows(rows)
	if parsed[0].Dividend != 0.5 || parsed[1].SplitRatio != 2 {
		t.Fatalf("corporate-action evidence was not retained by normalization: %+v", parsed)
	}
}

func TestV189TradeInsightDateOrdering(t *testing.T) {
	rows := []map[string]any{
		{"date": "2026-08-21", "adj_close": 102.0},
		{"date": "2026-08-19", "adj_close": 100.0},
		{"date": "2026-08-20", "adj_close": 101.0},
	}
	bars := tradeInsightBars(rows, true)
	if len(bars) != 3 {
		t.Fatalf("bars = %d", len(bars))
	}
	if !time.Unix(bars[0].T, 0).Before(time.Unix(bars[1].T, 0)) || !time.Unix(bars[1].T, 0).Before(time.Unix(bars[2].T, 0)) {
		t.Fatal("TradeInsight history must be canonicalized in ascending time order")
	}
}

func TestV189TradeInsightHistoryFanoutObeysBulkAdmissionAndCap(t *testing.T) {
	input := []string{"VIX", " aapl ", "AAPL"}
	for i := 0; i < 60; i++ {
		input = append(input, fmt.Sprintf("T%02d", i))
	}

	gated := tradeInsightHistorySymbolsForRefresh(nil, input, false)
	if len(gated) != 1 || gated[0] != "AAPL" {
		t.Fatalf("gated bulk history must collapse to one eligible canonical symbol, got %v", gated)
	}

	admitted := tradeInsightHistorySymbolsForRefresh(nil, input, true)
	if len(admitted) != tradeInsightHistoryFanoutMaxSymbols {
		t.Fatalf("admitted fan-out = %d symbols, want hard cap %d", len(admitted), tradeInsightHistoryFanoutMaxSymbols)
	}
	seen := map[string]bool{}
	for _, sym := range admitted {
		if sym == "" || sym == "VIX" {
			t.Fatalf("ineligible symbol escaped bounded history selection: %q", sym)
		}
		if seen[sym] {
			t.Fatalf("duplicate symbol escaped canonical history selection: %q", sym)
		}
		seen[sym] = true
	}
	if !seen["AAPL"] {
		t.Fatalf("normalized canonical symbol missing from admitted fan-out: %v", admitted)
	}
}

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

func TestTradeInsightCongressUsesPersistedEngineCredential(t *testing.T) {
	t.Setenv("TIDATA_API_KEY", "")
	t.Setenv("TRADEINSIGHT_API_KEY", "")
	app := &Application{secrets: Secrets{TradeInsight: "persisted-congress-key"}}
	e := &Engine{app: app}
	if got := e.tradeInsightResolvedAPIKey(); got != "persisted-congress-key" {
		t.Fatalf("Congress runtime must be able to resolve persisted Settings key, got %q", got)
	}

	raw, err := os.ReadFile("tradeinsight_congress.go")
	if err != nil {
		t.Fatal(err)
	}
	src := string(raw)
	for _, needle := range []string{
		"key := e.tradeInsightResolvedAPIKey()",
		"tradeInsightRESTBaseURL, key, symbol, begin",
	} {
		if !strings.Contains(src, needle) {
			t.Fatalf("Congress runtime missing persisted credential wiring %q", needle)
		}
	}
	for _, forbidden := range []string{
		"!tradeInsightConfigured()",
		"tradeInsightAPIKey(), symbol, begin",
	} {
		if strings.Contains(src, forbidden) {
			t.Fatalf("Congress runtime regressed to env-only credential path %q", forbidden)
		}
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
