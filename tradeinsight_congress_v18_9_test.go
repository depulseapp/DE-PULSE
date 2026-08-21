package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"sync/atomic"
	"testing"
)

func tradeInsightCongressFixture() tradeInsightCongressAPIRow {
	return tradeInsightCongressAPIRow{
		TraderName:      "Nancy Pelosi",
		TraderUUID:      "f2b9e3a1-4c8d-5e6f-a7b8-c9d0e1f2a3b4",
		Ticker:          "NVDA",
		IssuerName:      "NVIDIA Corp",
		TransactionType: "buy",
		Value:           1500000,
		TradedDate:      "2024-01-08",
		FilingDate:      "2024-01-15",
		Chamber:         "House",
	}
}

func TestV189TradeInsightCongressNormalizationUsesOfficialNumericValueAndDisclosureLag(t *testing.T) {
	got := normalizeTradeInsightCongressRows("NVDA", []tradeInsightCongressAPIRow{tradeInsightCongressFixture()})
	if len(got) != 1 {
		t.Fatalf("expected one normalized disclosure, got %d", len(got))
	}
	row := got[0]
	if row.Value != 1500000 {
		t.Fatalf("expected official numeric value, got %v", row.Value)
	}
	if row.DisclosureLagDays != 7 || !row.DisclosureLagValid {
		t.Fatalf("expected valid 7-day disclosure lag, got days=%d valid=%v", row.DisclosureLagDays, row.DisclosureLagValid)
	}
	if row.AmendmentState != "NOT_INFERRED_NO_VENDOR_ID" {
		t.Fatalf("unexpected amendment state %q", row.AmendmentState)
	}
	if !strings.HasPrefix(row.ID, "tradeinsight-congress-") {
		t.Fatalf("unexpected deterministic id %q", row.ID)
	}
}

func TestV189TradeInsightCongressExactDuplicatesCollapseButChangesRemainDistinct(t *testing.T) {
	base := tradeInsightCongressFixture()
	changed := base
	changed.Value = 1600000
	got := normalizeTradeInsightCongressRows("NVDA", []tradeInsightCongressAPIRow{base, base, changed})
	if len(got) != 2 {
		t.Fatalf("expected exact duplicate collapse with changed row preserved, got %d", len(got))
	}
	if got[0].AmendmentState != "NOT_INFERRED_NO_VENDOR_ID" || got[1].AmendmentState != "NOT_INFERRED_NO_VENDOR_ID" {
		t.Fatal("vendor schema has no amendment identity; non-identical disclosures must not be labeled amendments")
	}
}

func TestV189TradeInsightCongressMalformedPartialRowsAreSkipped(t *testing.T) {
	valid := tradeInsightCongressFixture()
	missingTrader := valid
	missingTrader.TraderName, missingTrader.TraderUUID = "", ""
	missingDate := valid
	missingDate.FilingDate = ""
	badDate := valid
	badDate.TradedDate = "not-a-date"
	wrongTicker := valid
	wrongTicker.Ticker = "AAPL"
	got := normalizeTradeInsightCongressRows("NVDA", []tradeInsightCongressAPIRow{missingTrader, missingDate, badDate, wrongTicker, valid})
	if len(got) != 1 {
		t.Fatalf("expected only the verified complete row, got %d", len(got))
	}
}

func TestV189TradeInsightCongressNegativeLagIsPreservedButInvalid(t *testing.T) {
	row := tradeInsightCongressFixture()
	row.TradedDate = "2024-01-15"
	row.FilingDate = "2024-01-08"
	got := normalizeTradeInsightCongressRows("NVDA", []tradeInsightCongressAPIRow{row})
	if len(got) != 1 {
		t.Fatalf("expected row preserved, got %d", len(got))
	}
	if got[0].DisclosureLagDays != -7 || got[0].DisclosureLagValid {
		t.Fatalf("expected preserved invalid -7 day lag, got days=%d valid=%v", got[0].DisclosureLagDays, got[0].DisclosureLagValid)
	}
}

func TestV189TradeInsightCongressFetchUsesExactDocumentedRequestAndObservation(t *testing.T) {
	var calls atomic.Int32
	var began atomic.Int32
	var done atomic.Int32
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		calls.Add(1)
		if r.Method != http.MethodGet {
			t.Errorf("method = %s", r.Method)
		}
		if r.URL.Path != tradeInsightCongressTradesPath {
			t.Errorf("path = %s", r.URL.Path)
		}
		if r.URL.Query().Get("ticker") != "NVDA" || r.URL.Query().Get("limit") != "200" || r.URL.Query().Get("offset") != "0" {
			t.Errorf("unexpected query %q", r.URL.RawQuery)
		}
		if r.Header.Get("Authorization") != "Bearer secret-key" {
			t.Errorf("authorization header not set")
		}
		if r.Header.Get("Accept") != "application/json" {
			t.Errorf("accept header = %q", r.Header.Get("Accept"))
		}
		_ = json.NewEncoder(w).Encode(map[string]any{
			"data":   []tradeInsightCongressAPIRow{tradeInsightCongressFixture()},
			"total":  1,
			"limit":  200,
			"offset": 0,
		})
	}))
	defer server.Close()

	begin := func() func(error) {
		began.Add(1)
		return func(err error) {
			if err != nil {
				t.Errorf("observer received error: %v", err)
			}
			done.Add(1)
		}
	}
	result, err := tradeInsightFetchCongressAtObserved(context.Background(), server.Client(), server.URL, "secret-key", "nvda", begin)
	if err != nil {
		t.Fatal(err)
	}
	if calls.Load() != 1 || began.Load() != 1 || done.Load() != 1 {
		t.Fatalf("calls=%d began=%d done=%d", calls.Load(), began.Load(), done.Load())
	}
	if len(result.Rows) != 1 || result.Total != 1 || result.Truncated {
		t.Fatalf("unexpected result %+v", result)
	}
}

func TestV189TradeInsightCongressPaginationHonorsTwoHundredRowLimitAndTotal(t *testing.T) {
	var calls atomic.Int32
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		calls.Add(1)
		if r.URL.Query().Get("limit") != "200" {
			t.Errorf("limit must stay within documented Congress max: %q", r.URL.Query().Get("limit"))
		}
		offset, _ := strconv.Atoi(r.URL.Query().Get("offset"))
		rows := make([]tradeInsightCongressAPIRow, 0)
		count := 200
		if offset == 200 {
			count = 1
		}
		for i := 0; i < count; i++ {
			row := tradeInsightCongressFixture()
			row.TraderUUID = fmt.Sprintf("uuid-%03d", offset+i)
			rows = append(rows, row)
		}
		_ = json.NewEncoder(w).Encode(map[string]any{"data": rows, "total": 201, "limit": 200, "offset": offset})
	}))
	defer server.Close()

	result, err := tradeInsightFetchCongressAtObserved(context.Background(), server.Client(), server.URL, "secret-key", "NVDA", nil)
	if err != nil {
		t.Fatal(err)
	}
	if calls.Load() != 2 || len(result.Rows) != 201 || result.Total != 201 || result.Truncated {
		t.Fatalf("calls=%d rows=%d total=%d truncated=%v", calls.Load(), len(result.Rows), result.Total, result.Truncated)
	}
}

func TestV189TradeInsightCongressMissingKeyIsGracefulAndDoesNotCallProvider(t *testing.T) {
	var calls atomic.Int32
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		calls.Add(1)
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer server.Close()

	_, err := tradeInsightFetchCongressAtObserved(context.Background(), server.Client(), server.URL, "", "NVDA", nil)
	if err == nil || !strings.Contains(err.Error(), "TIDATA_API_KEY missing") {
		t.Fatalf("expected missing-key error, got %v", err)
	}
	if calls.Load() != 0 {
		t.Fatalf("provider called %d times", calls.Load())
	}
}

func TestV189TradeInsightCongressMalformedOrMissingDataIsRejected(t *testing.T) {
	for name, body := range map[string]string{
		"malformed":    `{`,
		"missing-data": `{"total":1,"limit":200,"offset":0}`,
		"null-data":    `{"data":null,"total":1,"limit":200,"offset":0}`,
		"object-data":  `{"data":{},"total":1,"limit":200,"offset":0}`,
	} {
		t.Run(name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				_, _ = w.Write([]byte(body))
			}))
			defer server.Close()
			if _, err := tradeInsightFetchCongressAtObserved(context.Background(), server.Client(), server.URL, "secret", "NVDA", nil); err == nil {
				t.Fatal("expected invalid response error")
			}
		})
	}
}

func TestV189TradeInsightCongressRateLimitRedactsKeyAndPreservesRetryAfter(t *testing.T) {
	const secret = "do-not-leak-this-key"
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Retry-After", "17")
		w.WriteHeader(http.StatusTooManyRequests)
		_, _ = w.Write([]byte("key=" + secret))
	}))
	defer server.Close()

	_, err := tradeInsightFetchCongressAtObserved(context.Background(), server.Client(), server.URL, secret, "NVDA", nil)
	if err == nil {
		t.Fatal("expected HTTP 429 error")
	}
	msg := err.Error()
	if strings.Contains(msg, secret) || !strings.Contains(msg, "[REDACTED]") || !strings.Contains(msg, "retry-after=17") {
		t.Fatalf("unsafe or incomplete error: %s", msg)
	}
}

func TestV189TradeInsightCongressEvidenceRecordsAreShadowAlternativeEvidence(t *testing.T) {
	normalized := normalizeTradeInsightCongressRows("NVDA", []tradeInsightCongressAPIRow{tradeInsightCongressFixture()})
	records := tradeInsightCongressEvidenceRecords(normalized, tradeInsightCongressFetchResult{Total: 42, Truncated: true})
	if len(records) != 1 {
		t.Fatalf("expected one record, got %d", len(records))
	}
	record := records[0]
	if record.Kind != "congressional-trade-disclosure" || record.Source != tradeInsightProviderName || record.FreshnessState != "SHADOW" {
		t.Fatalf("unexpected evidence metadata %+v", record)
	}
	if !strings.Contains(record.Provenance, "/congress/v1/trades") || !strings.Contains(record.Provenance, "SHADOW") {
		t.Fatalf("unexpected provenance %q", record.Provenance)
	}
	var payload tradeInsightCongressEvidencePayload
	if err := json.Unmarshal(record.Payload, &payload); err != nil {
		t.Fatal(err)
	}
	if payload.Trade.DisclosureLagDays != 7 || !payload.Trade.DisclosureLagValid || payload.Trade.AmendmentState != "NOT_INFERRED_NO_VENDOR_ID" {
		t.Fatalf("unexpected normalized semantics %+v", payload.Trade)
	}
	if !payload.Truncated || payload.Total != 42 || payload.Role != "alternative-evidence" || payload.Lifecycle != "SHADOW" {
		t.Fatalf("unexpected payload metadata %+v", payload)
	}
}
