package main

import (
	"math"
	"strings"
	"time"
)

func capabilityStatusFromHealth(configured bool, h string) string {
	if !configured {
		return "NOT ENTITLED"
	}
	s := strings.ToLower(strings.TrimSpace(h))
	if strings.Contains(s, "entitlement") || strings.Contains(s, "subscription") || strings.Contains(s, "plan limited") || strings.Contains(s, "403") || strings.Contains(s, "payment required") {
		return "PLAN LIMITED"
	}

	for _, bad := range []string{"unavailable", "degraded", "failed", "error", "offline", "disconnected", "not configured", "checking", "loading"} {
		if strings.Contains(s, bad) {
			return "TEMPORARILY UNAVAILABLE"
		}
	}

	for _, good := range []string{"healthy", "available", "connected", "direct fx", "official"} {
		if strings.Contains(s, good) {
			return "AVAILABLE"
		}
	}
	return "TEMPORARILY UNAVAILABLE"
}

func buildProviderCapabilityRegistry(settings Settings, secrets Secrets, health map[string]string, intel map[string]SymbolIntelligence, direct map[string]GlobalDriver) []ProviderCapabilityEntry {
	now := time.Now().UnixMilli()
	usesCore := []string{"Market Regime", "Day", "Swing", "Long", "Discovery", "Research", "Decision Queue", "Trade Readiness"}
	hasPremiumIntel := false
	for _, v := range intel {
		if v.RecommendationTrend != "" || v.PriceTarget > 0 || v.InsiderNetShares != 0 {
			hasPremiumIntel = true
			break
		}
	}
	directFX := false
	for k := range direct {
		if strings.HasPrefix(k, "fx_") {
			directFX = true
			break
		}
	}
	rows := []ProviderCapabilityEntry{
		{Provider: "Finnhub", Capability: "Primary U.S. equity + earnings/peers", Status: capabilityStatusFromHealth(secrets.Finnhub != "", health["quotes-rest"]), Detail: "Live/REST plus earnings surprise and peer context.", UpdatedAt: now, Uses: usesCore},
		{Provider: "Finnhub", Capability: "Analyst / insider premium context", Status: func() string {
			if secrets.Finnhub == "" {
				return "NOT ENTITLED"
			}
			if hasPremiumIntel {
				return "AVAILABLE"
			}
			return "PLAN LIMITED"
		}(), Detail: "Used only when endpoint entitlement returns real data; never required for deterministic scoring.", UpdatedAt: now, Uses: []string{"Research", "Swing", "Long", "Decision Queue"}},
		{Provider: "Alpaca", Capability: "IEX quotes / snapshots / liquidity", Status: capabilityStatusFromHealth(secrets.AlpacaKey != "" && secrets.AlpacaSecret != "", health["alpaca-live"]), Detail: "Bid/ask, size, spread, quote age and snapshot hydration.", UpdatedAt: now, Uses: []string{"Day", "Market Open Prep", "Trade Readiness", "Decision Queue"}},
		{Provider: "Alpaca", Capability: "SIP movers / most active", Status: func() string {
			if secrets.AlpacaKey == "" || secrets.AlpacaSecret == "" {
				return "NOT ENTITLED"
			}
			if strings.Contains(strings.ToLower(health["market-activity"]), "available") {
				return "AVAILABLE"
			}
			return "PLAN LIMITED"
		}(), Detail: "Discovery seed only when account entitlement permits.", UpdatedAt: now, Uses: []string{"Discovery", "Dashboard"}},
		{Provider: tradeInsightProviderName, Capability: "Adjusted daily OHLCV / corporate-action corroboration", Status: func() string {
			if !tradeInsightConfigured(secrets.TradeInsight) {
				return "NOT ENTITLED"
			}
			history := strings.ToLower(strings.TrimSpace(health["history"]))
			actions := strings.ToLower(strings.TrimSpace(health["tradeinsight-corporate-actions"]))
			if strings.Contains(history, "tradeinsight") || strings.Contains(actions, "healthy") {
				return "AVAILABLE"
			}
			return "SHADOW"
		}(), Detail: "Smart Router v2 member of canonical Historical Bars. Daily adjusted OHLCV only; dividends/splits are supplemental evidence merged into the canonical corporate-action ledger.", UpdatedAt: now, Uses: []string{"Historical Bars", "Research", "Corporate Actions", "Adaptive Intelligence"}},
		{Provider: "FRED", Capability: "Rates / credit / conditions / USD", Status: capabilityStatusFromHealth(secrets.FRED != "", health["fred-rates"]), Detail: "Slow macro state; cadence-aware cache.", UpdatedAt: now, Uses: []string{"Market Regime", "Swing", "Long", "Research", "Trade Readiness"}},
		{Provider: "BLS", Capability: "Inflation / labor / wages / PPI", Status: capabilityStatusFromHealth(true, health["bls-actuals"]), Detail: "Official release-triggered actuals; no invented consensus.", UpdatedAt: now, Uses: []string{"Market Regime", "Swing", "Long", "Research"}},
		{Provider: "EIA", Capability: "Petroleum / natural gas / energy state", Status: capabilityStatusFromHealth(secrets.EIA != "", health["eia-actuals"]), Detail: "Official energy release context.", UpdatedAt: now, Uses: []string{"Market Regime", "Research", "Swing", "Long"}},
		{Provider: "Twelve Data", Capability: "FX / direct global context", Status: func() string {
			if secrets.TwelveData == "" {
				return "NOT ENTITLED"
			}
			if directFX {
				return "AVAILABLE"
			}
			return capabilityStatusFromHealth(true, health["global-direct"])
		}(), Detail: "Direct FX/global where entitled; official/proxy/cache fallback remains truthful.", UpdatedAt: now, Uses: []string{"Market Regime", "Dashboard", "Research"}},
		{Provider: "Twelve Data", Capability: "VIX / indices / historical recovery", Status: capabilityStatusFromHealth(secrets.TwelveData != "", health["vix"]), Detail: "v15 primary VIX/index route and first historical fallback after Alpaca.", UpdatedAt: now, Uses: []string{"Market Regime", "Dashboard", "Day", "Swing", "Long"}},
		{Provider: "yfinance", Capability: "Recovery-only public market context", Status: "AVAILABLE", Detail: "Fallback only for VIX, historical bars, earnings and fundamentals; never the primary live production feed.", UpdatedAt: now, Uses: []string{"Data Freshness", "Research", "Recovery"}},
		{Provider: "CBOE", Capability: "Official VIX validation / delayed close", Status: "AVAILABLE", Detail: "Authoritative VIX validation and delayed/official fallback only; not a general stock provider.", UpdatedAt: now, Uses: []string{"VIX", "Market Regime", "Data Freshness"}},
		{Provider: "Marketaux", Capability: "Stock news fallback", Status: func() string {
			if secrets.Marketaux != "" {
				return "AVAILABLE"
			}
			return "NOT ENTITLED"
		}(), Detail: "Supplemental/fallback company news when Finnhub is unavailable.", UpdatedAt: now, Uses: []string{"News", "Dashboard", "Research", "Catalyst Watch"}},
	}
	return rows
}

func (e *Engine) setPreparationRich(key, state, detail string, success, late bool, attention string, summary, changed []string, exceptions []CheckpointException) {
	e.mu.Lock()
	p := e.preparations[key]
	if p.Key == "" {
		p.Key = key
		p.Label = key
	}
	now := time.Now()
	day := now.In(easternLocation()).Format("2006-01-02")
	if p.TradingDay != day {
		p.AttemptCount = 0
	}

	if strings.EqualFold(state, "RUNNING") && automaticPreparationTrigger(detail) {
		p.AttemptCount++
	}
	p.State = state
	p.Detail = detail
	p.LastAttempt = now.UnixMilli()
	p.TradingDay = day
	p.Late = late
	if attention != "" {
		p.Attention = attention
	}
	if summary != nil {
		p.Summary = append([]string(nil), summary...)
	}
	if changed != nil {
		p.Changed = append([]string(nil), changed...)
	}
	if exceptions != nil {
		p.Exceptions = append([]CheckpointException(nil), exceptions...)
	}
	if success {
		p.LastSuccess = p.LastAttempt
	}
	if key == "pre-market-prep" {
		p.Window = "3:15–3:50 AM ET"
		p.NextWindow = nextWindowAt(now, 3, 15)
	}
	if key == "market-open-prep" {
		p.Window = "9:20–9:25 AM ET"
		p.NextWindow = nextWindowAt(now, 9, 20)
	}
	e.preparations[key] = p
	e.health[key] = state + " · " + detail
	e.mu.Unlock()
}

func automaticPreparationTrigger(detail string) bool {
	detail = strings.ToLower(strings.TrimSpace(detail))
	return detail == "scheduled" || strings.Contains(detail, "missed-window catch-up")
}

func preparationRanForTradingDay(p PreparationJobStatus, now time.Time) bool {
	day := now.In(easternLocation()).Format("2006-01-02")
	sameDay := func(ts int64) bool {
		return ts > 0 && time.UnixMilli(ts).In(easternLocation()).Format("2006-01-02") == day
	}

	if sameDay(p.LastSuccess) {
		return true
	}
	if strings.EqualFold(p.State, "RUNNING") && sameDay(p.LastAttempt) {
		return true
	}
	return p.TradingDay == day && p.AttemptCount >= 3
}

func preparationRetryDue(p PreparationJobStatus, now time.Time, minGap time.Duration) bool {
	if p.LastAttempt <= 0 {
		return true
	}
	return now.Sub(time.UnixMilli(p.LastAttempt)) >= minGap
}

func readinessFreshnessGate(rows []FreshnessDiagnostic, required []string, now time.Time) (usable int, degraded bool, exceptions []CheckpointException) {
	by := map[string]FreshnessDiagnostic{}
	for _, row := range rows {
		by[row.Dataset] = row
	}
	for _, name := range required {
		row, ok := by[name]
		if !ok {
			degraded = true
			exceptions = append(exceptions, CheckpointException{Reason: name + " · freshness diagnostic missing", Severity: "HIGH", Target: "maintenance", Source: name, UpdatedAt: now.UnixMilli()})
			continue
		}
		state := strings.ToUpper(strings.TrimSpace(row.State))
		switch state {
		case "LIVE", "FRESH", "DUE SOON":
			usable++
		case "DELAYED":

			if name == "VIX" {
				usable++
				exceptions = append(exceptions, CheckpointException{Reason: name + " · DELAYED · " + defaultString(row.Reason, "past ideal cadence"), Severity: "MEDIUM", Target: "maintenance", Source: defaultString(row.Provider, name), UpdatedAt: now.UnixMilli()})
			} else {
				degraded = true
				exceptions = append(exceptions, CheckpointException{Reason: name + " · DELAYED · " + defaultString(row.Reason, "critical evidence is delayed"), Severity: "HIGH", Target: "maintenance", Source: defaultString(row.Provider, name), UpdatedAt: now.UnixMilli()})
			}
		default:
			degraded = true
			exceptions = append(exceptions, CheckpointException{Reason: name + " · " + defaultString(state, "UNAVAILABLE") + " · " + defaultString(row.Reason, "critical evidence not current"), Severity: "HIGH", Target: "maintenance", Source: defaultString(row.Provider, name), UpdatedAt: now.UnixMilli()})
		}
	}
	return
}

func checkpointAttention(exceptions []CheckpointException, degraded bool) string {
	if degraded {
		return "DATA DEGRADED"
	}
	critical := false
	for _, x := range exceptions {
		if strings.EqualFold(x.Severity, "HIGH") || strings.EqualFold(x.Severity, "CRITICAL") {
			critical = true
			break
		}
	}
	if critical {
		return "REVIEW REQUIRED"
	}
	if len(exceptions) > 0 {
		return "READY WITH CAUTION"
	}
	return "READY"
}

func (e *Engine) isTradingDay(at time.Time) bool {
	et := at.In(easternLocation())
	day := et.Format("2006-01-02")
	e.mu.RLock()
	cal, ok := e.alpacaCalendar[day]
	e.mu.RUnlock()
	if ok {
		return cal.Open != "" && cal.Close != ""
	}
	return et.Weekday() >= time.Monday && et.Weekday() <= time.Friday && !isUSMarketHoliday(et)
}

func premarketSnapshotFromBars(sym string, bars []Bar, q Quote, now time.Time) (PremarketSnapshot, bool) {
	loc := easternLocation()
	day := now.In(loc).Format("2006-01-02")
	var rows []Bar
	for _, b := range bars {
		t := time.Unix(b.T, 0).In(loc)
		mins := t.Hour()*60 + t.Minute()
		if t.Format("2006-01-02") == day && mins >= 4*60 && mins < 9*60+30 {
			rows = append(rows, b)
		}
	}
	if len(rows) == 0 {
		return PremarketSnapshot{}, false
	}
	pm := PremarketSnapshot{Symbol: sym, Open: rows[0].O, High: rows[0].H, Low: rows[0].L, Last: rows[len(rows)-1].C, Bars: len(rows), UpdatedAt: rows[len(rows)-1].T * 1000}
	for _, b := range rows {
		if b.H > pm.High {
			pm.High = b.H
		}
		if pm.Low == 0 || (b.L > 0 && b.L < pm.Low) {
			pm.Low = b.L
		}
		pm.Volume += b.V
	}
	base := q.PreviousClose
	if base <= 0 {
		base = q.PriorSessionClose
	}
	if base > 0 && pm.Last > 0 {
		pm.GapPercent = (pm.Last/base - 1) * 100
	}
	if pm.Low > 0 && pm.High >= pm.Low {
		pm.RangePercent = (pm.High/pm.Low - 1) * 100
	}
	return pm, true
}

func materialSECFilingForTradingRisk(f FilingItem) bool {
	form := strings.ToUpper(strings.TrimSpace(f.Form))
	category := strings.ToLower(strings.TrimSpace(f.Category))
	text := strings.ToLower(strings.TrimSpace(f.Description + " " + f.Meaning + " " + f.Items))
	if category == "offering" || strings.Contains(form, "S-1") || strings.Contains(form, "424B") || strings.Contains(form, "F-1") {
		return true
	}
	if category == "material" || form == "8-K" || strings.HasPrefix(form, "8-K/") {
		return true
	}
	if form == "10-Q" || form == "10-K" {
		return materialText(text) || strings.Contains(text, "earn") || strings.Contains(text, "guidance") || strings.Contains(text, "going concern")
	}
	if strings.HasPrefix(form, "4") || category == "insider" {

		sig := strings.ToUpper(strings.TrimSpace(f.Signal))
		return (sig == "BUY" || sig == "SELL") && (math.Abs(f.Value) >= 100000 || math.Abs(f.Shares) >= 10000)
	}
	if form == "13D" || form == "13D/A" {
		return true
	}
	return materialText(text)
}

func atrPercentFromBars(rows []Bar) float64 {
	if len(rows) < 3 {
		return 0
	}
	start := len(rows) - 15
	if start < 1 {
		start = 1
	}
	total := 0.0
	count := 0
	for i := start; i < len(rows); i++ {
		h, l, prev := rows[i].H, rows[i].L, rows[i-1].C
		tr := math.Max(h-l, math.Max(math.Abs(h-prev), math.Abs(l-prev)))
		if tr > 0 {
			total += tr
			count++
		}
	}
	price := rows[len(rows)-1].C
	if count == 0 || price <= 0 {
		return 0
	}
	return (total / float64(count)) / price * 100
}

func vwapFromCurrentSession(rows []Bar, now time.Time) float64 {
	loc := easternLocation()
	day := now.In(loc).Format("2006-01-02")
	totalPV, totalV := 0.0, 0.0
	for _, b := range rows {
		t := time.Unix(b.T, 0).In(loc)
		if t.Format("2006-01-02") != day || b.V <= 0 {
			continue
		}
		typ := (b.H + b.L + b.C) / 3
		if typ <= 0 {
			continue
		}
		totalPV += typ * b.V
		totalV += b.V
	}
	if totalV <= 0 {
		return 0
	}
	return totalPV / totalV
}

func extendedAtMarketOpen(q Quote, daily, intraday []Bar, pm PremarketSnapshot, now time.Time) bool {
	if math.Abs(q.ChangePercent) >= 8 {
		return true
	}
	atrPct := atrPercentFromBars(daily)
	if atrPct <= 0 {
		atrPct = 2.5
	}
	if math.Abs(pm.GapPercent) >= math.Max(3, 1.6*atrPct) || pm.RangePercent >= math.Max(4, 1.8*atrPct) {
		return true
	}
	vw := vwapFromCurrentSession(intraday, now)
	if vw > 0 && q.Price > 0 {
		dist := math.Abs(q.Price/vw-1) * 100
		if dist >= math.Max(2.5, 1.25*atrPct) {
			return true
		}
	}
	return false
}
