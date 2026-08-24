package main

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"sort"
	"strconv"
	"strings"
	"time"
)

const (
	tradeInsightCongressPageSize = 200
	tradeInsightCongressMaxPages = 5
)

type tradeInsightCongressAPIRow struct {
	TraderName      string  `json:"trader_name"`
	TraderUUID      string  `json:"trader_uuid"`
	Ticker          string  `json:"ticker"`
	IssuerName      string  `json:"issuer_name"`
	TransactionType string  `json:"transaction_type"`
	Value           float64 `json:"value"`
	TradedDate      string  `json:"traded_date"`
	FilingDate      string  `json:"filing_date"`
	Chamber         string  `json:"chamber"`
}

type tradeInsightCongressFetchResult struct {
	Rows      []tradeInsightCongressAPIRow
	Total     int
	Truncated bool
}

type tradeInsightCongressEvidence struct {
	ID                 string  `json:"id"`
	TraderName         string  `json:"traderName"`
	TraderUUID         string  `json:"traderUUID"`
	Ticker             string  `json:"ticker"`
	IssuerName         string  `json:"issuerName"`
	TransactionType    string  `json:"transactionType"`
	Value              float64 `json:"value"`
	TradedDate         string  `json:"tradedDate"`
	FilingDate         string  `json:"filingDate"`
	Chamber            string  `json:"chamber"`
	DisclosureLagDays  int     `json:"disclosureLagDays"`
	DisclosureLagValid bool    `json:"disclosureLagValid"`
	AmendmentState     string  `json:"amendmentState"`
}

type tradeInsightCongressEvidencePayload struct {
	Trade      tradeInsightCongressEvidence `json:"trade"`
	Lifecycle  string                       `json:"lifecycle"`
	Role       string                       `json:"role"`
	Truncated  bool                         `json:"truncated"`
	Total      int                          `json:"total"`
	Provenance string                       `json:"provenance"`
}

func tradeInsightCongressSemanticKey(row tradeInsightCongressAPIRow) string {
	parts := []string{
		strings.TrimSpace(row.TraderUUID),
		strings.TrimSpace(row.TraderName),
		normalizeSymbol(row.Ticker),
		strings.TrimSpace(row.IssuerName),
		strings.ToLower(strings.TrimSpace(row.TransactionType)),
		strconv.FormatFloat(row.Value, 'g', -1, 64),
		strings.TrimSpace(row.TradedDate),
		strings.TrimSpace(row.FilingDate),
		strings.TrimSpace(row.Chamber),
	}
	sum := sha256.Sum256([]byte(strings.Join(parts, "|")))
	return hex.EncodeToString(sum[:])
}

func normalizeTradeInsightCongressRows(symbol string, rows []tradeInsightCongressAPIRow) []tradeInsightCongressEvidence {
	symbol = normalizeSymbol(symbol)
	if symbol == "" {
		return nil
	}
	seen := map[string]bool{}
	out := make([]tradeInsightCongressEvidence, 0, len(rows))
	for _, row := range rows {
		row.Ticker = normalizeSymbol(row.Ticker)
		row.TraderName = strings.TrimSpace(row.TraderName)
		row.TraderUUID = strings.TrimSpace(row.TraderUUID)
		row.IssuerName = strings.TrimSpace(row.IssuerName)
		row.TransactionType = strings.ToLower(strings.TrimSpace(row.TransactionType))
		row.TradedDate = strings.TrimSpace(row.TradedDate)
		row.FilingDate = strings.TrimSpace(row.FilingDate)
		row.Chamber = strings.TrimSpace(row.Chamber)
		if row.Ticker != symbol || row.TransactionType == "" || row.TradedDate == "" || row.FilingDate == "" || (row.TraderName == "" && row.TraderUUID == "") {
			continue
		}
		tradedAt, tradedErr := time.Parse("2006-01-02", row.TradedDate)
		filedAt, filedErr := time.Parse("2006-01-02", row.FilingDate)
		if tradedErr != nil || filedErr != nil {
			continue
		}
		key := tradeInsightCongressSemanticKey(row)
		if seen[key] {
			continue
		}
		seen[key] = true
		lag := int(filedAt.Sub(tradedAt).Hours() / 24)
		out = append(out, tradeInsightCongressEvidence{
			ID:                 "tradeinsight-congress-" + key,
			TraderName:         row.TraderName,
			TraderUUID:         row.TraderUUID,
			Ticker:             row.Ticker,
			IssuerName:         row.IssuerName,
			TransactionType:    row.TransactionType,
			Value:              row.Value,
			TradedDate:         row.TradedDate,
			FilingDate:         row.FilingDate,
			Chamber:            row.Chamber,
			DisclosureLagDays:  lag,
			DisclosureLagValid: lag >= 0,
			AmendmentState:     "NOT_INFERRED_NO_VENDOR_ID",
		})
	}
	sort.Slice(out, func(i, j int) bool {
		if out[i].FilingDate == out[j].FilingDate {
			return out[i].ID < out[j].ID
		}
		return out[i].FilingDate > out[j].FilingDate
	})
	return out
}

func tradeInsightFetchCongressAtObserved(ctx context.Context, client *http.Client, baseURL, key, symbol string, begin func() func(error)) (tradeInsightCongressFetchResult, error) {
	var result tradeInsightCongressFetchResult
	key = strings.TrimSpace(key)
	if key == "" {
		return result, fmt.Errorf("TradeInsight not configured: TIDATA_API_KEY missing")
	}
	symbol = normalizeSymbol(symbol)
	if !validUserTicker(symbol) {
		return result, fmt.Errorf("TradeInsight congressional trades invalid ticker")
	}
	if client == nil {
		client = &http.Client{Timeout: 18 * time.Second}
	}
	baseURL = strings.TrimRight(strings.TrimSpace(baseURL), "/")
	if baseURL == "" {
		baseURL = tradeInsightRESTBaseURL
	}

	for page := 0; page < tradeInsightCongressMaxPages; page++ {
		offset := page * tradeInsightCongressPageSize
		q := url.Values{}
		q.Set("ticker", symbol)
		q.Set("limit", strconv.Itoa(tradeInsightCongressPageSize))
		q.Set("offset", strconv.Itoa(offset))
		raw := baseURL + tradeInsightCongressTradesPath + "?" + q.Encode()
		req, err := http.NewRequestWithContext(ctx, http.MethodGet, raw, nil)
		if err != nil {
			return result, err
		}
		req.Header.Set("Authorization", "Bearer "+key)
		req.Header.Set("Accept", "application/json")

		done := func(error) {}
		if begin != nil {
			if observed := begin(); observed != nil {
				done = observed
			}
		}
		resp, err := client.Do(req)
		if err != nil {
			wrapped := fmt.Errorf("TradeInsight congressional request failed: %w", err)
			done(wrapped)
			return result, wrapped
		}
		body, readErr := io.ReadAll(io.LimitReader(resp.Body, 8<<20))
		_ = resp.Body.Close()
		if readErr != nil {
			wrapped := fmt.Errorf("TradeInsight congressional response read failed: %w", readErr)
			done(wrapped)
			return result, wrapped
		}
		if resp.StatusCode < 200 || resp.StatusCode >= 300 {
			detail := tradeInsightSafeError(body, key)
			if retry := strings.TrimSpace(resp.Header.Get("Retry-After")); retry != "" {
				detail += " · retry-after=" + retry
			}
			wrapped := fmt.Errorf("TradeInsight congressional HTTP %d: %s", resp.StatusCode, detail)
			done(wrapped)
			return result, wrapped
		}
		var envelope struct {
			Data   json.RawMessage `json:"data"`
			Total  int             `json:"total"`
			Limit  int             `json:"limit"`
			Offset int             `json:"offset"`
		}
		if err := json.Unmarshal(body, &envelope); err != nil {
			wrapped := fmt.Errorf("TradeInsight congressional malformed JSON: %w", err)
			done(wrapped)
			return result, wrapped
		}
		if len(envelope.Data) == 0 || string(envelope.Data) == "null" {
			wrapped := fmt.Errorf("TradeInsight congressional response missing data array")
			done(wrapped)
			return result, wrapped
		}
		var rows []tradeInsightCongressAPIRow
		if err := json.Unmarshal(envelope.Data, &rows); err != nil {
			wrapped := fmt.Errorf("TradeInsight congressional invalid data array: %w", err)
			done(wrapped)
			return result, wrapped
		}
		done(nil)
		result.Rows = append(result.Rows, rows...)
		if envelope.Total > result.Total {
			result.Total = envelope.Total
		}
		if len(rows) < tradeInsightCongressPageSize || offset+len(rows) >= envelope.Total {
			return result, nil
		}
	}
	if result.Total > len(result.Rows) {
		result.Truncated = true
	}
	return result, nil
}

func tradeInsightCongressEvidenceRecords(rows []tradeInsightCongressEvidence, result tradeInsightCongressFetchResult) []EvidenceRecord {
	records := make([]EvidenceRecord, 0, len(rows))
	for _, row := range rows {
		filedAt, err := time.Parse("2006-01-02", row.FilingDate)
		if err != nil {
			continue
		}
		payload, err := json.Marshal(tradeInsightCongressEvidencePayload{
			Trade:      row,
			Lifecycle:  "SHADOW",
			Role:       "alternative-evidence",
			Truncated:  result.Truncated,
			Total:      result.Total,
			Provenance: "TradeInsight /congress/v1/trades · official REST schema · SHADOW",
		})
		if err != nil {
			continue
		}
		records = append(records, EvidenceRecord{
			ID:             row.ID,
			Symbol:         row.Ticker,
			Kind:           "congressional-trade-disclosure",
			Source:         tradeInsightProviderName,
			ObservedAt:     filedAt.UTC().UnixMilli(),
			FreshnessState: "SHADOW",
			Provenance:     "TradeInsight /congress/v1/trades · normalized · SHADOW",
			Payload:        payload,
		})
	}
	return records
}

func (e *Engine) refreshTradeInsightCongressResearchSymbol(ctx context.Context, symbol string) int {
	symbol = normalizeSymbol(symbol)
	if !validUserTicker(symbol) {
		return 0
	}
	admission, registered := tradeInsightCapabilityAdmissionLookup("congressional-trades")
	if !registered || !admission.runtimeAdmitted() {
		return 0
	}
	key := e.tradeInsightResolvedAPIKey()
	if key == "" || !e.providerAllowed(tradeInsightProviderName) {
		return 0
	}
	var begin func() func(error)
	if e.providerTelemetry != nil {
		begin = func() func(error) { return e.providerTelemetry.begin(tradeInsightProviderName) }
	}
	result, err := tradeInsightFetchCongressAtObserved(ctx, &http.Client{Timeout: 18 * time.Second}, tradeInsightRESTBaseURL, key, symbol, begin)
	if err != nil {
		e.recordProviderFailure(tradeInsightProviderName, err)
		e.setHealth("research-congress:"+symbol, "optional degraded · TradeInsight SHADOW")
		return 0
	}
	normalized := normalizeTradeInsightCongressRows(symbol, result.Rows)
	records := tradeInsightCongressEvidenceRecords(normalized, result)
	if len(records) > 0 && e.app != nil && e.app.persistence != nil {
		e.app.persistence.EnqueueIntelligence(PersistenceIntelligenceBatch{Evidence: records})
	}
	e.recordProviderSuccess(tradeInsightProviderName)
	e.mu.Lock()
	if e.lastUpdated != nil {
		e.lastUpdated["research-congress:"+symbol] = time.Now().UnixMilli()
	}
	if e.health != nil {
		state := fmt.Sprintf("healthy · TradeInsight · SHADOW · %d disclosures", len(records))
		if result.Truncated {
			state += " · bounded result"
		}
		e.health["research-congress:"+symbol] = state
	}
	e.mu.Unlock()
	return len(records)
}
