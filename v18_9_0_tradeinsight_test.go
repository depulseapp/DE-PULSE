package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
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
