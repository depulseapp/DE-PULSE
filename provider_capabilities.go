package main

import (
	"math"
	"sort"
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

// providerCapabilityLegacyDisplayOrder freezes only the pre-#95 Data Engine
// presentation order. It is intentionally not provider onboarding metadata: a
// future provider does not need to edit this list to become registered/routable.
// Unlisted diagnostics sort after the retained legacy rows by provider/capability.
func providerCapabilityLegacyDisplayOrder(provider, capability string) int {
	key := providerKey(provider) + "|" + strings.ToLower(strings.TrimSpace(capability))
	return map[string]int{
		providerKey("Finnhub") + "|primary u.s. equity + earnings/peers":                                 1,
		providerKey("Finnhub") + "|analyst / insider premium context":                                    2,
		providerKey("Alpaca") + "|iex quotes / snapshots / liquidity":                                    3,
		providerKey("Alpaca") + "|sip movers / most active":                                              4,
		providerKey(tradeInsightProviderName) + "|adjusted daily ohlcv / corporate-action corroboration": 5,
		providerKey("FRED") + "|rates / credit / conditions / usd":                                       6,
		providerKey("BLS") + "|inflation / labor / wages / ppi":                                          7,
		providerKey("EIA") + "|petroleum / natural gas / energy state":                                   8,
		providerKey("Twelve Data") + "|fx / direct global context":                                       9,
		providerKey("Twelve Data") + "|vix / indices / historical recovery":                              10,
		providerKey("yfinance") + "|recovery-only public market context":                                 11,
		providerKey("CBOE") + "|official vix validation / delayed close":                                 12,
		providerKey("Marketaux") + "|stock news fallback":                                                13,
	}[key]
}

func buildProviderCapabilityRegistryFromRegistrations(regs []ProviderRegistration, settings Settings, secrets Secrets, health map[string]string, intel map[string]SymbolIntelligence, direct map[string]GlobalDriver) []ProviderCapabilityEntry {
	type displayRow struct {
		entry       ProviderCapabilityEntry
		legacyOrder int
	}
	now := time.Now().UnixMilli()
	displayRows := make([]displayRow, 0)
	for _, reg := range regs {
		for _, diagnostic := range reg.Diagnostics {
			status := "TEMPORARILY UNAVAILABLE"
			if diagnostic.Status != nil {
				status = diagnostic.Status(settings, secrets, health, intel, direct)
			}
			entry := ProviderCapabilityEntry{
				Provider: reg.Name, Capability: diagnostic.Capability, Status: status,
				Detail: diagnostic.Detail, UpdatedAt: now, Uses: append([]string(nil), diagnostic.Uses...),
			}
			displayRows = append(displayRows, displayRow{entry: entry, legacyOrder: providerCapabilityLegacyDisplayOrder(reg.Name, diagnostic.Capability)})
		}
	}
	sort.SliceStable(displayRows, func(i, j int) bool {
		left, right := displayRows[i], displayRows[j]
		leftLegacy, rightLegacy := left.legacyOrder > 0, right.legacyOrder > 0
		if leftLegacy != rightLegacy {
			return leftLegacy
		}
		if leftLegacy && left.legacyOrder != right.legacyOrder {
			return left.legacyOrder < right.legacyOrder
		}
		if left.entry.Provider != right.entry.Provider {
			return strings.ToLower(left.entry.Provider) < strings.ToLower(right.entry.Provider)
		}
		return strings.ToLower(left.entry.Capability) < strings.ToLower(right.entry.Capability)
	})
	rows := make([]ProviderCapabilityEntry, 0, len(displayRows))
	for _, row := range displayRows {
		rows = append(rows, row.entry)
	}
	return rows
}

func buildProviderCapabilityRegistry(settings Settings, secrets Secrets, health map[string]string, intel map[string]SymbolIntelligence, direct map[string]GlobalDriver) []ProviderCapabilityEntry {
	return buildProviderCapabilityRegistryFromRegistrations(providerRegistrations(), settings, secrets, health, intel, direct)
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
