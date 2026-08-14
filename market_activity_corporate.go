package main

import (
	"context"
	"fmt"
	"math"
	"net/http"
	"net/url"
	"sort"
	"strconv"
	"strings"
	"time"
)

func (e *Engine) refreshAlpacaMarketCalendar(ctx context.Context, key, secret string) {
	start := time.Now().In(easternLocation()).AddDate(0, 0, -2).Format("2006-01-02")
	end := time.Now().In(easternLocation()).AddDate(0, 2, 0).Format("2006-01-02")
	raw := strings.TrimRight(alpacaTradingBaseURL, "/") + "/v2/calendar?start=" + url.QueryEscape(start) + "&end=" + url.QueryEscape(end)
	var rows []AlpacaCalendarDay
	err := getJSON(ctx, &http.Client{Timeout: 12 * time.Second}, raw, map[string]string{"APCA-API-KEY-ID": key, "APCA-API-SECRET-KEY": secret}, &rows)
	if err != nil {
		e.setHealth("market-calendar", "degraded · using built-in U.S. calendar fallback")
		return
	}
	cal := map[string]AlpacaCalendarDay{}
	for _, r := range rows {
		cal[r.Date] = r
	}
	e.mu.Lock()
	e.alpacaCalendar = cal
	e.lastUpdated["market-calendar"] = time.Now().UnixMilli()
	e.health["market-calendar"] = fmt.Sprintf("healthy · Alpaca calendar · %d sessions", len(rows))
	e.mu.Unlock()
}

func parseMoverRows(v any) []MarketMover {
	arr, ok := v.([]any)
	if !ok {
		return nil
	}
	out := []MarketMover{}
	for _, x := range arr {
		m, ok := x.(map[string]any)
		if !ok {
			continue
		}
		sym := normalizeSymbol(fmt.Sprint(m["symbol"]))
		if sym == "" {
			continue
		}
		mv := MarketMover{Symbol: sym}
		for _, k := range []string{"percent_change", "change_percent", "changePercentage"} {
			if z, ok := m[k]; ok {
				mv.ChangePercent, _ = strconv.ParseFloat(fmt.Sprint(z), 64)
				break
			}
		}
		if z, ok := m["volume"]; ok {
			mv.Volume, _ = strconv.ParseFloat(fmt.Sprint(z), 64)
		}
		out = append(out, mv)
	}
	return out
}
func (e *Engine) refreshAlpacaMarketActivity(ctx context.Context, key, secret string) {
	headers := map[string]string{"APCA-API-KEY-ID": key, "APCA-API-SECRET-KEY": secret}
	client := &http.Client{Timeout: 12 * time.Second}
	state := MarketActivityState{UpdatedAt: time.Now().UnixMilli()}
	success := 0
	var active map[string]any
	if getJSON(ctx, client, strings.TrimRight(alpacaDataBaseURL, "/")+"/v1beta1/screener/stocks/most-actives?top=20&by=volume", headers, &active) == nil {
		if v, ok := active["most_actives"]; ok {
			state.MostActive = parseMoverRows(v)
			success++
		}
	}
	var movers map[string]any
	if getJSON(ctx, client, strings.TrimRight(alpacaDataBaseURL, "/")+"/v1beta1/screener/stocks/movers?top=20", headers, &movers) == nil {
		if v, ok := movers["gainers"]; ok {
			state.Gainers = parseMoverRows(v)
		}
		if v, ok := movers["losers"]; ok {
			state.Losers = parseMoverRows(v)
		}
		success++
	}
	e.mu.Lock()
	defer e.mu.Unlock()
	if success == 0 {
		e.health["market-activity"] = "plan limited or unavailable · Discovery does not depend on it"
		return
	}
	state.Status = "AVAILABLE"
	e.marketActivity = state
	e.lastUpdated["market-activity"] = state.UpdatedAt
	e.health["market-activity"] = "available · Alpaca market activity"
}

func corporateActionSymbol(row map[string]any, tracked map[string]bool) string {

	for _, k := range []string{"symbol", "old_symbol", "source_symbol", "acquiree_symbol", "acquirer_symbol", "new_symbol"} {
		s := normalizeSymbol(fmt.Sprint(row[k]))
		if s != "" && tracked[s] {
			return s
		}
	}
	return ""
}

func corporateActionDetail(kind string, row map[string]any) string {
	parts := []string{}
	if oldRate, newRate := toFloat(row["old_rate"]), toFloat(row["new_rate"]); oldRate > 0 && newRate > 0 {
		parts = append(parts, fmt.Sprintf("rate %.4g:%.4g", newRate, oldRate))
	} else if rate := toFloat(row["rate"]); rate != 0 {
		parts = append(parts, fmt.Sprintf("rate %.4g", rate))
	}
	if cashRate := toFloat(row["cash_rate"]); cashRate != 0 {
		parts = append(parts, fmt.Sprintf("cash %.4g", cashRate))
	}
	for _, k := range []string{"new_symbol", "acquirer_symbol"} {
		if s := normalizeSymbol(fmt.Sprint(row[k])); s != "" {
			parts = append(parts, "related "+s)
			break
		}
	}
	if len(parts) == 0 {
		parts = append(parts, strings.ReplaceAll(kind, "_", " "))
	}
	return strings.Join(parts, " · ")
}

func parseAlpacaCorporateActionResponse(payload map[string]any, tracked map[string]bool) []CorporateAction {
	out := []CorporateAction{}
	for bucket, rawRows := range payload {
		if bucket == "next_page_token" || bucket == "nextPageToken" {
			continue
		}
		rows, ok := rawRows.([]any)
		if !ok {
			continue
		}
		kind := strings.TrimSuffix(bucket, "s")
		for _, raw := range rows {
			r, ok := raw.(map[string]any)
			if !ok {
				continue
			}
			sym := corporateActionSymbol(r, tracked)
			if sym == "" {
				continue
			}
			id := strings.TrimSpace(fmt.Sprint(r["id"]))
			if id == "<nil>" {
				id = ""
			}
			processDate := strings.TrimSpace(fmt.Sprint(r["process_date"]))
			if processDate == "<nil>" {
				processDate = ""
			}
			exDate := strings.TrimSpace(fmt.Sprint(r["ex_date"]))
			if exDate == "<nil>" {
				exDate = ""
			}
			if exDate == "" {
				exDate = strings.TrimSpace(fmt.Sprint(r["effective_date"]))
				if exDate == "<nil>" {
					exDate = ""
				}
			}
			recordDate := strings.TrimSpace(fmt.Sprint(r["record_date"]))
			if recordDate == "<nil>" {
				recordDate = ""
			}
			payableDate := strings.TrimSpace(fmt.Sprint(r["payable_date"]))
			if payableDate == "<nil>" {
				payableDate = ""
			}
			oldSym := normalizeSymbol(fmt.Sprint(r["old_symbol"]))
			newSym := normalizeSymbol(fmt.Sprint(r["new_symbol"]))
			if oldSym == "" && kind == "name_change" {
				oldSym = normalizeSymbol(fmt.Sprint(r["source_symbol"]))
			}
			if newSym == "" && kind == "name_change" {
				newSym = normalizeSymbol(fmt.Sprint(r["symbol"]))
			}
			oldRate, newRate := toFloat(r["old_rate"]), toFloat(r["new_rate"])
			ratio := toFloat(r["rate"])
			if oldRate > 0 && newRate > 0 {
				ratio = newRate / oldRate
			}
			status := "UPCOMING"
			if exDate != "" {
				if d, err := time.Parse("2006-01-02", exDate); err == nil && !d.After(time.Now()) {
					status = "EFFECTIVE"
				}
			}
			nowMs := time.Now().UnixMilli()
			out = append(out, CorporateAction{
				ID: id, Symbol: sym, Type: kind, ProcessDate: processDate, ExDate: exDate, RecordDate: recordDate, PayableDate: payableDate,
				OldSymbol: oldSym, NewSymbol: newSym, Ratio: ratio, CashAmount: toFloat(r["cash_rate"]), AdjustmentFactor: ratio,
				Status: status, FirstSeenAt: nowMs, UpdatedAt: nowMs, Detail: corporateActionDetail(kind, r), Source: "Alpaca",
			})
		}
	}
	return out
}

func (e *Engine) refreshAlpacaCorporateActions(ctx context.Context, key, secret string) {
	now := time.Now()
	endDate := now.AddDate(0, 3, 0).Format("2006-01-02")
	headers := map[string]string{"APCA-API-KEY-ID": key, "APCA-API-SECRET-KEY": secret}
	tracked := map[string]bool{}
	syms := []string{}
	for _, s := range e.trackedSymbols() {
		if s == "VIX" || tracked[s] {
			continue
		}
		tracked[s] = true
		syms = append(syms, s)
	}
	if len(syms) == 0 {
		return
	}
	sort.Strings(syms)
	client := &http.Client{Timeout: 18 * time.Second}
	fresh := []CorporateAction{}
	truncated := false
	anyFailure := false
	entitlementFailure := false
	backfilled := []string{}
	for i := 0; i < len(syms); i += 50 {
		j := i + 50
		if j > len(syms) {
			j = len(syms)
		}
		batch := syms[i:j]
		startDate := now.AddDate(0, 0, -30).Format("2006-01-02")
		needsBackfill := false
		e.mu.RLock()
		for _, sym := range batch {
			if e.lastUpdated["corporate-actions-backfill:"+sym] == 0 {
				needsBackfill = true
				break
			}
		}
		e.mu.RUnlock()
		if needsBackfill {
			startDate = now.AddDate(-15, 0, 0).Format("2006-01-02")
		}
		pageToken := ""
		seenTokens := map[string]bool{}
		success := true
		batchComplete := false
		batchTruncated := false
		for page := 0; ; page++ {
			if page >= 100 {
				truncated = true
				batchTruncated = true
				break
			}
			q := url.Values{}
			q.Set("symbols", strings.Join(batch, ","))
			q.Set("start", startDate)
			q.Set("end", endDate)
			q.Set("limit", "1000")
			q.Set("sort", "asc")
			if pageToken != "" {
				q.Set("page_token", pageToken)
			}
			raw := strings.TrimRight(alpacaDataBaseURL, "/") + "/v1/corporate-actions?" + q.Encode()
			var payload map[string]any
			if err := getJSON(ctx, client, raw, headers, &payload); err != nil {
				success = false
				anyFailure = true
				msg := strings.ToLower(err.Error())
				if strings.Contains(msg, "403") || strings.Contains(msg, "401") {
					entitlementFailure = true
				}
				break
			}
			fresh = append(fresh, parseAlpacaCorporateActionResponse(payload, tracked)...)
			next := strings.TrimSpace(fmt.Sprint(payload["next_page_token"]))
			if next == "" || next == "<nil>" {
				batchComplete = true
				break
			}
			if seenTokens[next] {
				truncated = true
				batchTruncated = true
				break
			}
			seenTokens[next] = true
			pageToken = next
		}
		if success && batchComplete && !batchTruncated && needsBackfill {
			backfilled = append(backfilled, batch...)
		}
	}
	nowMs := time.Now().UnixMilli()
	e.mu.Lock()
	merged := mergeCorporateActionLedger(e.corporateActions, fresh, nowMs)
	e.corporateActions = merged
	e.lastUpdated["corporate-actions"] = nowMs
	for _, sym := range backfilled {
		e.lastUpdated["corporate-actions-backfill:"+sym] = nowMs
	}
	if entitlementFailure {
		e.health["corporate-actions"] = fmt.Sprintf("plan limited or not entitled · persistent ledger retained · %d total actions", len(merged))
	} else if anyFailure {
		e.health["corporate-actions"] = fmt.Sprintf("partial · persistent ledger retained · %d total actions · latest provider refresh incomplete", len(merged))
	} else if truncated {
		e.health["corporate-actions"] = fmt.Sprintf("partial · persistent ledger · %d total actions · pagination safety stop", len(merged))
	} else {
		e.health["corporate-actions"] = fmt.Sprintf("healthy · persistent ledger · %d total actions · %d refreshed", len(merged), len(fresh))
	}
	e.mu.Unlock()

	if len(merged) > 0 {
		_ = e.refreshAlpacaRawHistoryForCorporateActions(ctx, key, secret, merged)
	}
}

func (e *Engine) refreshFinnhubIntelligence(ctx context.Context, key string) {
	syms := e.trackedSymbols()
	if len(syms) > 20 {
		syms = syms[:20]
	}
	updated := 0
	for _, sym := range syms {
		if sym == "VIX" {
			continue
		}
		si := SymbolIntelligence{Symbol: sym, UpdatedAt: time.Now().UnixMilli()}
		notes := []string{}
		var earnings []struct {
			Period                                      string `json:"period"`
			Actual, Estimate, Surprise, SurprisePercent float64
		}
		if err := e.finnhubJSONForSymbol(ctx, key, sym, "/stock/earnings?symbol="+url.QueryEscape(sym)+"&limit=4", &earnings); err == nil {
			for _, x := range earnings {
				si.EarningsSurprises = append(si.EarningsSurprises, EarningsSurprisePoint{Period: x.Period, Actual: x.Actual, Estimate: x.Estimate, Surprise: x.Surprise, SurprisePercent: x.SurprisePercent})
			}
			for _, x := range si.EarningsSurprises {
				if x.Actual > x.Estimate {
					if si.ConsecutiveMisses > 0 {
						break
					}
					si.ConsecutiveBeats++
				} else if x.Actual < x.Estimate {
					if si.ConsecutiveBeats > 0 {
						break
					}
					si.ConsecutiveMisses++
				} else {
					break
				}
			}
		} else {
			notes = append(notes, "Earnings surprise endpoint unavailable")
		}
		var peers []string
		if err := e.finnhubJSONForSymbol(ctx, key, sym, "/stock/peers?symbol="+url.QueryEscape(sym), &peers); err == nil {
			si.Peers = uniqueSymbols(peers)
		} else {
			notes = append(notes, "Peer endpoint unavailable")
		}
		var recRaw []map[string]any
		if err := e.finnhubJSONForSymbol(ctx, key, sym, "/stock/recommendation?symbol="+url.QueryEscape(sym), &recRaw); err == nil && len(recRaw) > 0 {
			r := recRaw[0]
			buy := toInt(r["buy"]) + toInt(r["strongBuy"])
			sell := toInt(r["sell"]) + toInt(r["strongSell"])
			hold := toInt(r["hold"])
			if buy > sell+hold/2 {
				si.RecommendationTrend = "BULLISH"
			} else if sell > buy {
				si.RecommendationTrend = "BEARISH"
			} else {
				si.RecommendationTrend = "MIXED"
			}
		} else {
			notes = append(notes, "Recommendation endpoint plan-limited/unavailable")
		}
		var pt map[string]any
		if err := e.finnhubJSONForSymbol(ctx, key, sym, "/stock/price-target?symbol="+url.QueryEscape(sym), &pt); err == nil {
			for _, k := range []string{"targetMean", "targetMedian"} {
				if v, ok := pt[k]; ok {
					si.PriceTarget = toFloat(v)
					if si.PriceTarget > 0 {
						break
					}
				}
			}
		} else {
			notes = append(notes, "Price-target endpoint plan-limited/unavailable")
		}
		var ins struct {
			Data []map[string]any `json:"data"`
		}
		if err := e.finnhubJSONForSymbol(ctx, key, sym, "/stock/insider-transactions?symbol="+url.QueryEscape(sym), &ins); err == nil {
			for _, r := range ins.Data {
				si.InsiderNetShares += toFloat(r["change"])
				if si.InsiderNetShares == 0 {
					si.InsiderNetShares += toFloat(r["share"])
				}
			}
		} else {
			notes = append(notes, "Insider endpoint plan-limited/unavailable")
		}
		var own map[string]any
		if err := e.finnhubJSONForSymbol(ctx, key, sym, "/stock/fund-ownership?symbol="+url.QueryEscape(sym)+"&limit=20", &own); err == nil {
			rows := []any{}
			for _, field := range []string{"ownership", "data"} {
				if x, ok := own[field].([]any); ok {
					rows = x
					break
				}
			}
			for _, raw := range rows {
				r, ok := raw.(map[string]any)
				if !ok {
					continue
				}
				si.InstitutionalOwners++
				for _, f := range []string{"change", "positionChange", "changeShares"} {
					if v, ok := r[f]; ok {
						si.InstitutionalNetChange += toFloat(v)
						break
					}
				}
			}
		} else {
			notes = append(notes, "Fund ownership endpoint plan-limited/unavailable")
		}
		si.EntitlementNotes = notes
		e.mu.Lock()
		e.symbolIntelligence[sym] = si
		e.mu.Unlock()
		updated++
	}
	e.mu.Lock()
	e.lastUpdated["symbol-intelligence"] = time.Now().UnixMilli()
	if updated > 0 {
		e.health["symbol-intelligence"] = fmt.Sprintf("healthy · %d symbols · premium fields entitlement-aware", updated)
	} else {
		e.health["symbol-intelligence"] = "degraded · no intelligence refreshed"
	}
	e.mu.Unlock()
}
func toInt(v any) int { return int(math.Round(toFloat(v))) }
