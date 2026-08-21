package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"sort"
	"strings"
	"time"
)

const (
	canonicalHistoricalBarsDataset = "Historical Bars"
	tradeInsightProviderName       = "TradeInsight"
	tradeInsightRESTBaseURL        = "https://api.tradeinsight.info/trading-data/v1"
	tradeInsightPageSize           = 1000
	tradeInsightMaxPages           = 10
)

type tradeInsightHistoryRow struct {
	Date       string
	Open       float64
	High       float64
	Low        float64
	Close      float64
	Volume     float64
	AdjOpen    float64
	AdjHigh    float64
	AdjLow     float64
	AdjClose   float64
	AdjVolume  float64
	Dividend   float64
	SplitRatio float64
}

func tradeInsightAPIKey() string {
	if key := strings.TrimSpace(os.Getenv("TIDATA_API_KEY")); key != "" {
		return key
	}
	return strings.TrimSpace(os.Getenv("TRADEINSIGHT_API_KEY"))
}

func tradeInsightConfigured() bool { return tradeInsightAPIKey() != "" }

func tradeInsightSafeError(body []byte, key string) string {
	msg := strings.TrimSpace(string(body))
	if msg == "" {
		return "request failed"
	}
	if key != "" {
		msg = strings.ReplaceAll(msg, key, "[REDACTED]")
	}
	if len(msg) > 300 {
		msg = msg[:300]
	}
	return msg
}

func tradeInsightRowsFromPayload(body []byte) ([]map[string]any, error) {
	var direct []map[string]any
	if err := json.Unmarshal(body, &direct); err == nil {
		return direct, nil
	}
	var envelope struct {
		Data []map[string]any `json:"data"`
	}
	if err := json.Unmarshal(body, &envelope); err != nil {
		return nil, fmt.Errorf("TradeInsight malformed JSON: %w", err)
	}
	return envelope.Data, nil
}

func tradeInsightFetchRowsAt(ctx context.Context, client *http.Client, baseURL, key, path string, params url.Values) ([]map[string]any, error) {
	return tradeInsightFetchRowsAtObserved(ctx, client, baseURL, key, path, params, nil)
}

// tradeInsightFetchRowsAtObserved preserves the deterministic standalone fetch
// helper while allowing runtime calls to report every HTTP page into DE.PULSE's
// shared ProviderTelemetry. Smart Router v2 can therefore learn from actual
// TradeInsight latency, errors and rate-limit pressure instead of treating this
// provider as telemetry-unknown.
func tradeInsightFetchRowsAtObserved(ctx context.Context, client *http.Client, baseURL, key, path string, params url.Values, begin func() func(error)) ([]map[string]any, error) {
	key = strings.TrimSpace(key)
	if key == "" {
		return nil, fmt.Errorf("TradeInsight not configured: TIDATA_API_KEY missing")
	}
	if client == nil {
		client = &http.Client{Timeout: 18 * time.Second}
	}
	baseURL = strings.TrimRight(strings.TrimSpace(baseURL), "/")
	if baseURL == "" {
		baseURL = tradeInsightRESTBaseURL
	}
	all := make([]map[string]any, 0, tradeInsightPageSize)
	for page := 0; page < tradeInsightMaxPages; page++ {
		q := url.Values{}
		for k, values := range params {
			for _, value := range values {
				q.Add(k, value)
			}
		}
		q.Set("limit", fmt.Sprint(tradeInsightPageSize))
		q.Set("offset", fmt.Sprint(page*tradeInsightPageSize))
		raw := baseURL + "/" + strings.TrimLeft(path, "/") + "?" + q.Encode()
		req, err := http.NewRequestWithContext(ctx, http.MethodGet, raw, nil)
		if err != nil {
			return nil, err
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
			wrapped := fmt.Errorf("TradeInsight request failed: %w", err)
			done(wrapped)
			return nil, wrapped
		}
		body, readErr := io.ReadAll(io.LimitReader(resp.Body, 8<<20))
		_ = resp.Body.Close()
		if readErr != nil {
			wrapped := fmt.Errorf("TradeInsight response read failed: %w", readErr)
			done(wrapped)
			return nil, wrapped
		}
		if resp.StatusCode < 200 || resp.StatusCode >= 300 {
			detail := tradeInsightSafeError(body, key)
			if retry := strings.TrimSpace(resp.Header.Get("Retry-After")); retry != "" {
				detail += " · retry-after=" + retry
			}
			wrapped := fmt.Errorf("TradeInsight HTTP %d: %s", resp.StatusCode, detail)
			done(wrapped)
			return nil, wrapped
		}
		rows, err := tradeInsightRowsFromPayload(body)
		if err != nil {
			done(err)
			return nil, err
		}
		done(nil)
		all = append(all, rows...)
		if len(rows) < tradeInsightPageSize {
			return all, nil
		}
	}
	return all, fmt.Errorf("TradeInsight pagination safety stop after %d pages", tradeInsightMaxPages)
}

func tradeInsightFetchRows(ctx context.Context, path string, params url.Values) ([]map[string]any, error) {
	return tradeInsightFetchRowsAt(ctx, &http.Client{Timeout: 18 * time.Second}, tradeInsightRESTBaseURL, tradeInsightAPIKey(), path, params)
}

func (e *Engine) tradeInsightFetchRows(ctx context.Context, path string, params url.Values) ([]map[string]any, error) {
	var begin func() func(error)
	if e != nil && e.providerTelemetry != nil {
		begin = func() func(error) { return e.providerTelemetry.begin(tradeInsightProviderName) }
	}
	return tradeInsightFetchRowsAtObserved(ctx, &http.Client{Timeout: 18 * time.Second}, tradeInsightRESTBaseURL, tradeInsightAPIKey(), path, params, begin)
}

func tradeInsightHistoryRows(rows []map[string]any) []tradeInsightHistoryRow {
	out := make([]tradeInsightHistoryRow, 0, len(rows))
	for _, row := range rows {
		date := strings.TrimSpace(fmt.Sprint(row["date"]))
		if date == "" || date == "<nil>" {
			continue
		}
		out = append(out, tradeInsightHistoryRow{
			Date:       date,
			Open:       toFloat(row["open"]),
			High:       toFloat(row["high"]),
			Low:        toFloat(row["low"]),
			Close:      toFloat(row["close"]),
			Volume:     toFloat(row["volume"]),
			AdjOpen:    toFloat(row["adj_open"]),
			AdjHigh:    toFloat(row["adj_high"]),
			AdjLow:     toFloat(row["adj_low"]),
			AdjClose:   toFloat(row["adj_close"]),
			AdjVolume:  toFloat(row["adj_volume"]),
			Dividend:   toFloat(row["dividend"]),
			SplitRatio: toFloat(row["split_ratio"]),
		})
	}
	return out
}

func tradeInsightCorporateActions(symbol string, rows []map[string]any, now int64) []CorporateAction {
	symbol = normalizeSymbol(symbol)
	if symbol == "" || symbol == "VIX" {
		return nil
	}
	if now <= 0 {
		now = time.Now().UnixMilli()
	}
	today := time.UnixMilli(now).UTC().Format("2006-01-02")
	actions := []CorporateAction{}
	for _, row := range tradeInsightHistoryRows(rows) {
		if _, err := time.Parse("2006-01-02", row.Date); err != nil {
			continue
		}
		status := "UPCOMING"
		if row.Date <= today {
			status = "EFFECTIVE"
		}
		if row.Dividend > 0 {
			actions = append(actions, CorporateAction{
				Symbol: symbol, Type: "cash_dividend", ExDate: row.Date, CashAmount: row.Dividend,
				Status: status, FirstSeenAt: now, UpdatedAt: now,
				Detail: fmt.Sprintf("cash %.4g · supplemental adjusted-history evidence", row.Dividend), Source: tradeInsightProviderName,
			})
		}
		if row.SplitRatio > 0 && row.SplitRatio != 1 {
			actions = append(actions, CorporateAction{
				Symbol: symbol, Type: "split", ExDate: row.Date, Ratio: row.SplitRatio, AdjustmentFactor: row.SplitRatio,
				Status: status, FirstSeenAt: now, UpdatedAt: now,
				Detail: fmt.Sprintf("ratio %.4g · supplemental adjusted-history evidence", row.SplitRatio), Source: tradeInsightProviderName,
			})
		}
	}
	return actions
}

func corporateActionSemanticKey(a CorporateAction) string {
	a.ID = ""
	return corporateActionKey(a)
}

func mergeSupplementalCorporateActions(existing, supplemental []CorporateAction, now int64) []CorporateAction {
	seen := map[string]bool{}
	for _, action := range existing {
		if key := corporateActionSemanticKey(action); key != "" {
			seen[key] = true
		}
	}
	fresh := make([]CorporateAction, 0, len(supplemental))
	for _, action := range supplemental {
		key := corporateActionSemanticKey(action)
		if key == "" || seen[key] {
			continue
		}
		seen[key] = true
		fresh = append(fresh, action)
	}
	return mergeCorporateActionLedger(existing, fresh, now)
}

func tradeInsightBars(rows []map[string]any, adjusted bool) []Bar {
	parsed := tradeInsightHistoryRows(rows)
	bars := make([]Bar, 0, len(parsed))
	for _, row := range parsed {
		t, err := time.Parse("2006-01-02", row.Date)
		if err != nil {
			continue
		}
		o, h, l, c, v := row.Open, row.High, row.Low, row.Close, row.Volume
		if adjusted {
			o, h, l, c, v = row.AdjOpen, row.AdjHigh, row.AdjLow, row.AdjClose, row.AdjVolume
		}
		if c <= 0 {
			continue
		}
		bars = append(bars, Bar{T: t.Unix(), O: o, H: h, L: l, C: c, V: v})
	}
	sort.Slice(bars, func(i, j int) bool { return bars[i].T < bars[j].T })
	return bars
}

func aggregateDailyBarsToWeekly(rows []Bar) []Bar {
	if len(rows) == 0 {
		return nil
	}
	out := []Bar{}
	lastYear, lastWeek := -1, -1
	for _, row := range rows {
		t := time.Unix(row.T, 0).UTC()
		year, week := t.ISOWeek()
		if len(out) == 0 || year != lastYear || week != lastWeek {
			out = append(out, Bar{T: row.T, O: row.O, H: row.H, L: row.L, C: row.C, V: row.V})
			lastYear, lastWeek = year, week
			continue
		}
		cur := &out[len(out)-1]
		if row.H > cur.H {
			cur.H = row.H
		}
		if cur.L == 0 || (row.L > 0 && row.L < cur.L) {
			cur.L = row.L
		}
		cur.C = row.C
		cur.V += row.V
	}
	return out
}

func (e *Engine) refreshTradeInsightHistoryMode(ctx context.Context, only []string, mode string) int {
	mode = strings.ToLower(strings.TrimSpace(mode))
	// The official tidata SDK documents only daily (1d) history. TradeInsight
	// therefore participates only in the daily/weekly canonical history refresh;
	// it must never masquerade as an intraday provider.
	if mode != "daily" {
		return 0
	}
	if !tradeInsightConfigured() || !e.providerAllowed(tradeInsightProviderName) {
		return 0
	}
	symbols := historyRouteSymbols(e, only)
	if len(symbols) > 50 {
		symbols = symbols[:50]
	}
	loaded := 0
	for _, sym := range symbols {
		if sym == "" || sym == "VIX" {
			continue
		}
		start := time.Now().AddDate(-2, 0, 0)
		if sym == "SPY" || sym == "QQQ" {
			start = time.Now().AddDate(-10, 0, 0)
		}
		params := url.Values{
			"ticker":        []string{sym},
			"start":         []string{start.Format("2006-01-02")},
			"end":           []string{time.Now().AddDate(0, 0, 1).Format("2006-01-02")},
			"adjust_volume": []string{"true"},
		}
		rows, err := e.tradeInsightFetchRows(ctx, "/ohlc", params)
		if err != nil {
			e.recordProviderFailure(tradeInsightProviderName, err)
			continue
		}
		daily := tradeInsightBars(rows, true)
		if len(daily) == 0 {
			continue
		}
		weekly := aggregateDailyBarsToWeekly(daily)
		nowMs := time.Now().UnixMilli()
		actions := tradeInsightCorporateActions(sym, rows, nowMs)
		e.mu.Lock()
		if e.bars[sym] == nil {
			e.bars[sym] = map[string][]Bar{}
		}
		e.bars[sym]["daily"] = daily
		if len(weekly) > 0 {
			e.bars[sym]["weekly"] = weekly
		}
		if len(actions) > 0 {
			e.corporateActions = mergeSupplementalCorporateActions(e.corporateActions, actions, nowMs)
			e.lastUpdated["tradeinsight-corporate-actions"] = nowMs
			e.health["tradeinsight-corporate-actions"] = fmt.Sprintf("healthy · %d supplemental dividend/split events normalized into canonical ledger", len(actions))
		}
		e.mu.Unlock()
		prior := 0.0
		if len(daily) > 1 {
			prior = daily[len(daily)-2].C
		}
		e.updateCanonicalSessionClose(sym, daily[len(daily)-1].C, daily[len(daily)-1].T*1000, prior)
		loaded += len(daily) + len(weekly)
		e.recordProviderSuccess(tradeInsightProviderName)
	}
	if loaded > 0 {
		nowMs := time.Now().UnixMilli()
		e.mu.Lock()
		e.lastUpdated["history"] = nowMs
		e.lastUpdated["history-daily"] = nowMs
		e.health["history"] = "healthy · TradeInsight daily fallback · adjusted OHLCV · weekly derived from canonical daily bars"
		e.mu.Unlock()
	}
	return loaded
}
