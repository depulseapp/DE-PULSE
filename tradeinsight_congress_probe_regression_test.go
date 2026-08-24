package main

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"sync/atomic"
	"testing"
)

func tradeInsightSchemaFieldByName(t *testing.T, fingerprint tradeInsightSchemaFingerprint, name string) tradeInsightSchemaField {
	t.Helper()
	for _, field := range fingerprint.Fields {
		if field.Name == name {
			return field
		}
	}
	t.Fatalf("schema field %q missing from %+v", name, fingerprint.Fields)
	return tradeInsightSchemaField{}
}

func TestV189TradeInsightCongressSchemaFingerprintIsDeterministicAndValueFree(t *testing.T) {
	body := []byte(`[
		{"ticker":"AAPL","amount":123.45,"nested":{"secret":"do-not-log"},"nullable":null},
		{"ticker":"MSFT","amount":"$1,001-$15,000","extra":true}
	]`)
	fingerprint, err := tradeInsightSchemaFingerprintFromPayload(body)
	if err != nil {
		t.Fatal(err)
	}
	if fingerprint.Container != "array" || fingerprint.RowsObserved != 2 || fingerprint.NonObjectRows != 0 {
		t.Fatalf("unexpected fingerprint summary: %+v", fingerprint)
	}
	gotNames := make([]string, 0, len(fingerprint.Fields))
	for _, field := range fingerprint.Fields {
		gotNames = append(gotNames, field.Name)
	}
	if strings.Join(gotNames, ",") != "amount,extra,nested,nullable,ticker" {
		t.Fatalf("field ordering = %v", gotNames)
	}
	amount := tradeInsightSchemaFieldByName(t, fingerprint, "amount")
	if strings.Join(amount.Types, ",") != "number,string" || amount.PresentInRows != 2 {
		t.Fatalf("amount schema = %+v", amount)
	}
	encoded, err := json.Marshal(fingerprint)
	if err != nil {
		t.Fatal(err)
	}
	for _, rawValue := range []string{"AAPL", "MSFT", "123.45", "$1,001-$15,000", "do-not-log"} {
		if strings.Contains(string(encoded), rawValue) {
			t.Fatalf("schema fingerprint leaked response value %q: %s", rawValue, encoded)
		}
	}
}

func TestV189TradeInsightCongressSchemaFingerprintHandlesEnvelopeAndPartialRows(t *testing.T) {
	body := []byte(`{"data":[{"ticker":"AAPL"},null,"partial",42,{"ticker":null,"chamber":"house"}]}`)
	fingerprint, err := tradeInsightSchemaFingerprintFromPayload(body)
	if err != nil {
		t.Fatal(err)
	}
	if fingerprint.Container != "data" || fingerprint.RowsObserved != 2 || fingerprint.NonObjectRows != 3 {
		t.Fatalf("unexpected partial fingerprint: %+v", fingerprint)
	}
	ticker := tradeInsightSchemaFieldByName(t, fingerprint, "ticker")
	if strings.Join(ticker.Types, ",") != "null,string" || ticker.PresentInRows != 2 {
		t.Fatalf("ticker schema = %+v", ticker)
	}
	chamber := tradeInsightSchemaFieldByName(t, fingerprint, "chamber")
	if chamber.PresentInRows != 1 {
		t.Fatalf("chamber presence = %+v", chamber)
	}
}

func TestV189TradeInsightCongressSchemaFingerprintRejectsMalformedOrUnverifiedContainers(t *testing.T) {
	cases := [][]byte{
		[]byte(`{"unexpected":[]}`),
		[]byte(`{"data":{}}`),
		[]byte(`[]`),
		[]byte(`not-json`),
	}
	for _, body := range cases {
		if _, err := tradeInsightSchemaFingerprintFromPayload(body); err == nil {
			t.Fatalf("expected schema-probe rejection for %q", body)
		}
	}
}

func TestV189TradeInsightCongressSchemaProbeUsesExactDocumentedRequestAndSharedObservation(t *testing.T) {
	const key = "test-congress-key"
	var observedBegin atomic.Int32
	var observedDone atomic.Int32
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Fatalf("method = %s", r.Method)
		}
		if r.URL.Path != tradeInsightCongressTradesPath {
			t.Fatalf("path = %q", r.URL.Path)
		}
		if r.URL.RawQuery != "" {
			t.Fatalf("schema probe must stay value-minimal: %q", r.URL.RawQuery)
		}
		if got := r.Header.Get("Authorization"); got != "Bearer "+key {
			t.Fatalf("authorization = %q", got)
		}
		if got := r.Header.Get("Accept"); got != "application/json" {
			t.Fatalf("accept = %q", got)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`[{
			"opaque_field":"opaque-value",
			"optional_number":7
		}]`))
	}))
	defer server.Close()

	begin := func() func(error) {
		observedBegin.Add(1)
		return func(err error) {
			if err != nil {
				t.Errorf("unexpected observed error: %v", err)
			}
			observedDone.Add(1)
		}
	}
	fingerprint, err := tradeInsightCongressSchemaProbeAtObserved(context.Background(), server.Client(), server.URL, key, begin)
	if err != nil {
		t.Fatal(err)
	}
	if fingerprint.RowsObserved != 1 || observedBegin.Load() != 1 || observedDone.Load() != 1 {
		t.Fatalf("fingerprint/observation mismatch: fp=%+v begin=%d done=%d", fingerprint, observedBegin.Load(), observedDone.Load())
	}
	encoded, _ := json.Marshal(fingerprint)
	if strings.Contains(string(encoded), "opaque-value") {
		t.Fatalf("probe retained response values: %s", encoded)
	}
}

func TestV189TradeInsightCongressSchemaProbeMissingKeyIsGracefulAndDoesNotCallProvider(t *testing.T) {
	var calls atomic.Int32
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		calls.Add(1)
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer server.Close()

	if _, err := tradeInsightCongressSchemaProbeAt(context.Background(), server.Client(), server.URL, "   "); err == nil || !strings.Contains(err.Error(), "not configured") {
		t.Fatalf("missing-key error = %v", err)
	}
	if calls.Load() != 0 {
		t.Fatalf("provider calls with missing key = %d", calls.Load())
	}
}

func TestV189TradeInsightCongressSchemaProbeRateLimitIsRedactedAndBackpressureVisible(t *testing.T) {
	const key = "super-secret-provider-key"
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Retry-After", "17")
		w.WriteHeader(http.StatusTooManyRequests)
		_, _ = w.Write([]byte("rate limited for " + key))
	}))
	defer server.Close()

	_, err := tradeInsightCongressSchemaProbeAt(context.Background(), server.Client(), server.URL, key)
	if err == nil {
		t.Fatal("expected 429 error")
	}
	text := err.Error()
	if strings.Contains(text, key) {
		t.Fatalf("provider key leaked in error: %s", text)
	}
	if !strings.Contains(text, "429") || !strings.Contains(text, "retry-after=17") || !strings.Contains(text, "[REDACTED]") {
		t.Fatalf("rate-limit/backpressure truth missing: %s", text)
	}
}

func TestV189TradeInsightCongressProbeCannotPromoteBeyondShadow(t *testing.T) {
	before, ok := tradeInsightCapabilityAdmissionLookup("congressional-trades")
	if !ok {
		t.Fatal("Congressional admission row missing")
	}
	beforeLifecycle := tradeInsightCapabilityLifecycleTruth(before.ID)
	if !before.SchemaVerified || !before.RuntimeEnabled || !before.runtimeAdmitted() || beforeLifecycle != "PRODUCTION" {
		t.Fatalf("Congress must enter probe diagnostics already governed PRODUCTION after #84 promotion: %+v lifecycle=%s", before, beforeLifecycle)
	}
	_ = tradeInsightSchemaFingerprint{Container: "data", RowsObserved: 1}
	after, _ := tradeInsightCapabilityAdmissionLookup("congressional-trades")
	afterLifecycle := tradeInsightCapabilityLifecycleTruth(after.ID)
	if after != before || afterLifecycle != beforeLifecycle {
		t.Fatalf("schema probe must not mutate or auto-promote/demote Congress: before=%+v after=%+v lifecycle before=%s after=%s", before, after, beforeLifecycle, afterLifecycle)
	}
}

func TestV189TradeInsightCongressSchemaProbeLiveOptIn(t *testing.T) {
	if os.Getenv("DE_PULSE_TRADEINSIGHT_CONGRESS_SCHEMA_PROBE") != "1" {
		t.Skip("set DE_PULSE_TRADEINSIGHT_CONGRESS_SCHEMA_PROBE=1 with TIDATA_API_KEY for optional value-free runtime schema diagnostics")
	}
	if !tradeInsightConfigured() {
		t.Fatal("TIDATA_API_KEY or TRADEINSIGHT_API_KEY must be configured for the live schema diagnostic")
	}
	fingerprint, err := tradeInsightCongressSchemaProbeAt(context.Background(), nil, tradeInsightRESTBaseURL, tradeInsightAPIKey())
	if err != nil {
		t.Fatal(err)
	}
	encoded, err := json.Marshal(fingerprint)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("TradeInsight Congress value-free schema fingerprint: %s", encoded)
}
