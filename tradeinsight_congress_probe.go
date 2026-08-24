package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sort"
	"strings"
	"time"
)

const tradeInsightCongressTradesPath = "/congress/v1/trades"

// tradeInsightSchemaField records only structural schema evidence. It never
// retains response values, which keeps configured-key diagnostics safe to log
// without leaking provider payload contents.
type tradeInsightSchemaField struct {
	Name          string   `json:"name"`
	Types         []string `json:"types"`
	PresentInRows int      `json:"present_in_rows"`
}

// tradeInsightSchemaFingerprint is intentionally value-free. It captures the
// top-level container and observed field/type contract for runtime diagnostics
// without changing capability admission or lifecycle state.
type tradeInsightSchemaFingerprint struct {
	Container     string                    `json:"container"`
	RowsObserved  int                       `json:"rows_observed"`
	NonObjectRows int                       `json:"non_object_rows"`
	Fields        []tradeInsightSchemaField `json:"fields"`
}

func tradeInsightJSONType(value any) string {
	switch value.(type) {
	case nil:
		return "null"
	case bool:
		return "boolean"
	case float64:
		return "number"
	case string:
		return "string"
	case []any:
		return "array"
	case map[string]any:
		return "object"
	default:
		return "unknown"
	}
}

func tradeInsightSchemaFingerprintFromPayload(body []byte) (tradeInsightSchemaFingerprint, error) {
	var payload any
	if err := json.Unmarshal(body, &payload); err != nil {
		return tradeInsightSchemaFingerprint{}, fmt.Errorf("TradeInsight Congress malformed JSON: %w", err)
	}

	fingerprint := tradeInsightSchemaFingerprint{}
	var records []any
	switch top := payload.(type) {
	case []any:
		fingerprint.Container = "array"
		records = top
	case map[string]any:
		raw, ok := top["data"]
		if !ok {
			return tradeInsightSchemaFingerprint{}, fmt.Errorf("TradeInsight Congress schema probe: object response has no data array")
		}
		rows, ok := raw.([]any)
		if !ok {
			return tradeInsightSchemaFingerprint{}, fmt.Errorf("TradeInsight Congress schema probe: data is not an array")
		}
		fingerprint.Container = "data"
		records = rows
	default:
		return tradeInsightSchemaFingerprint{}, fmt.Errorf("TradeInsight Congress schema probe: unsupported top-level JSON type %s", tradeInsightJSONType(payload))
	}

	typesByField := map[string]map[string]bool{}
	presentByField := map[string]int{}
	for _, record := range records {
		row, ok := record.(map[string]any)
		if !ok {
			fingerprint.NonObjectRows++
			continue
		}
		fingerprint.RowsObserved++
		for rawName, value := range row {
			name := strings.TrimSpace(rawName)
			if name == "" {
				continue
			}
			if typesByField[name] == nil {
				typesByField[name] = map[string]bool{}
			}
			typesByField[name][tradeInsightJSONType(value)] = true
			presentByField[name]++
		}
	}

	if fingerprint.RowsObserved == 0 {
		return fingerprint, fmt.Errorf("TradeInsight Congress schema probe returned no object rows")
	}

	fieldNames := make([]string, 0, len(typesByField))
	for name := range typesByField {
		fieldNames = append(fieldNames, name)
	}
	sort.Strings(fieldNames)
	fingerprint.Fields = make([]tradeInsightSchemaField, 0, len(fieldNames))
	for _, name := range fieldNames {
		types := make([]string, 0, len(typesByField[name]))
		for typ := range typesByField[name] {
			types = append(types, typ)
		}
		sort.Strings(types)
		fingerprint.Fields = append(fingerprint.Fields, tradeInsightSchemaField{
			Name:          name,
			Types:         types,
			PresentInRows: presentByField[name],
		})
	}
	return fingerprint, nil
}

func tradeInsightCongressSchemaProbeAt(ctx context.Context, client *http.Client, baseURL, key string) (tradeInsightSchemaFingerprint, error) {
	return tradeInsightCongressSchemaProbeAtObserved(ctx, client, baseURL, key, nil)
}

// tradeInsightCongressSchemaProbeAtObserved is an explicit value-free runtime
// diagnostic, not a production data route. The request intentionally remains
// query-free so it does not depend on the production Congressional adapter's
// filtering or pagination behavior.
func tradeInsightCongressSchemaProbeAtObserved(ctx context.Context, client *http.Client, baseURL, key string, begin func() func(error)) (tradeInsightSchemaFingerprint, error) {
	key = strings.TrimSpace(key)
	if key == "" {
		return tradeInsightSchemaFingerprint{}, fmt.Errorf("TradeInsight Congress schema probe not configured: TIDATA_API_KEY missing")
	}

	admission, ok := tradeInsightCapabilityAdmissionLookup("congressional-trades")
	if !ok {
		return tradeInsightSchemaFingerprint{}, fmt.Errorf("TradeInsight Congress admission row missing")
	}
	expectedEvidence := "/trading-data/v1" + tradeInsightCongressTradesPath
	if !strings.Contains(admission.EndpointEvidence, expectedEvidence) {
		return tradeInsightSchemaFingerprint{}, fmt.Errorf("TradeInsight Congress endpoint evidence mismatch: got %q", admission.EndpointEvidence)
	}

	if client == nil {
		client = &http.Client{Timeout: 18 * time.Second}
	}
	baseURL = strings.TrimRight(strings.TrimSpace(baseURL), "/")
	if baseURL == "" {
		baseURL = tradeInsightRESTBaseURL
	}
	raw := baseURL + tradeInsightCongressTradesPath
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, raw, nil)
	if err != nil {
		return tradeInsightSchemaFingerprint{}, err
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
		wrapped := fmt.Errorf("TradeInsight Congress schema probe request failed: %w", err)
		done(wrapped)
		return tradeInsightSchemaFingerprint{}, wrapped
	}
	body, readErr := io.ReadAll(io.LimitReader(resp.Body, 8<<20))
	_ = resp.Body.Close()
	if readErr != nil {
		wrapped := fmt.Errorf("TradeInsight Congress schema probe response read failed: %w", readErr)
		done(wrapped)
		return tradeInsightSchemaFingerprint{}, wrapped
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		detail := tradeInsightSafeError(body, key)
		if retry := strings.TrimSpace(resp.Header.Get("Retry-After")); retry != "" {
			detail += " · retry-after=" + retry
		}
		wrapped := fmt.Errorf("TradeInsight Congress schema probe HTTP %d: %s", resp.StatusCode, detail)
		done(wrapped)
		return tradeInsightSchemaFingerprint{}, wrapped
	}

	fingerprint, err := tradeInsightSchemaFingerprintFromPayload(body)
	done(err)
	return fingerprint, err
}

// probeTradeInsightCongressSchema participates in shared provider telemetry but
// does not write freshness, health, cache, Research, Event Intelligence or
// canonical persistence state. A failed optional diagnostic therefore cannot
// create a DATA DEGRADED cascade or change the Congressional SHADOW lifecycle.
func (e *Engine) probeTradeInsightCongressSchema(ctx context.Context) (tradeInsightSchemaFingerprint, error) {
	var begin func() func(error)
	if e != nil && e.providerTelemetry != nil {
		begin = func() func(error) { return e.providerTelemetry.begin(tradeInsightProviderName) }
	}
	return tradeInsightCongressSchemaProbeAtObserved(ctx, &http.Client{Timeout: 18 * time.Second}, tradeInsightRESTBaseURL, e.tradeInsightResolvedAPIKey(), begin)
}
