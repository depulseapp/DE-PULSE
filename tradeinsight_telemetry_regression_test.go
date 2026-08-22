package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

func tradeInsightTelemetryRow(t *testing.T, telemetry *ProviderTelemetry) ProviderRequestDiagnostics {
	t.Helper()
	for _, row := range telemetry.Diagnostics() {
		if row.Provider == tradeInsightProviderName {
			return row
		}
	}
	t.Fatal("TradeInsight provider telemetry row is missing")
	return ProviderRequestDiagnostics{}
}

func TestV189TradeInsightObservedFetchReportsEveryPage(t *testing.T) {
	telemetry := NewProviderTelemetry()
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if r.URL.Query().Get("offset") == "0" {
			rows := make([]map[string]any, tradeInsightPageSize)
			for i := range rows {
				rows[i] = map[string]any{"date": "2026-08-20", "adj_close": float64(i + 1)}
			}
			_ = json.NewEncoder(w).Encode(map[string]any{"data": rows})
			return
		}
		_ = json.NewEncoder(w).Encode(map[string]any{"data": []map[string]any{{"date": "2026-08-21", "adj_close": 101.0}}})
	}))
	defer server.Close()

	begin := func() func(error) { return telemetry.begin(tradeInsightProviderName) }
	rows, err := tradeInsightFetchRowsAtObserved(context.Background(), server.Client(), server.URL, "fixture", "/ohlc", url.Values{"ticker": []string{"SPY"}}, begin)
	if err != nil {
		t.Fatal(err)
	}
	if len(rows) != tradeInsightPageSize+1 {
		t.Fatalf("rows = %d, want %d", len(rows), tradeInsightPageSize+1)
	}
	diag := tradeInsightTelemetryRow(t, telemetry)
	if diag.Requests != 2 || diag.Successes != 2 || diag.Errors != 0 {
		t.Fatalf("unexpected TradeInsight telemetry: %+v", diag)
	}
}

func TestV189TradeInsightObservedFetchReportsRateLimit(t *testing.T) {
	telemetry := NewProviderTelemetry()
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Retry-After", "11")
		w.WriteHeader(http.StatusTooManyRequests)
		_, _ = fmt.Fprint(w, `{"message":"rate limit"}`)
	}))
	defer server.Close()

	begin := func() func(error) { return telemetry.begin(tradeInsightProviderName) }
	_, err := tradeInsightFetchRowsAtObserved(context.Background(), server.Client(), server.URL, "fixture", "/ohlc", url.Values{"ticker": []string{"SPY"}}, begin)
	if err == nil {
		t.Fatal("expected TradeInsight 429 error")
	}
	diag := tradeInsightTelemetryRow(t, telemetry)
	if diag.Requests != 1 || diag.Errors != 1 || diag.RateLimited != 1 {
		t.Fatalf("rate-limit telemetry not recorded: %+v", diag)
	}
	if diag.BudgetState != "RATE LIMITED" {
		t.Fatalf("budget state = %q, want RATE LIMITED", diag.BudgetState)
	}
}
