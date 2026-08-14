package main

import (
	"fmt"
	"sort"
	"strings"
	"time"
)

func freshnessLimits(dataset, provider, session string) (cadence, fresh, stale time.Duration) {
	active := session == "regular" || session == "pre-market" || session == "after-hours"
	switch dataset {
	case "Quotes":
		if active {
			return 30 * time.Second, 90 * time.Second, 2 * time.Minute
		}
		return 5 * time.Minute, 15 * time.Minute, 30 * time.Minute
	case "VIX":
		if provider == "CBOE" {
			return 10 * time.Minute, 30 * time.Minute, 60 * time.Minute
		}
		if provider == "yfinance" {
			return 2 * time.Minute, 10 * time.Minute, 20 * time.Minute
		}
		return 2 * time.Minute, 5 * time.Minute, 10 * time.Minute
	case "Intraday Bars":
		if active {
			return 5 * time.Minute, 7 * time.Minute, 10 * time.Minute
		}
		return 15 * time.Minute, 30 * time.Minute, 60 * time.Minute
	case "Daily / Weekly History":
		return 24 * time.Hour, 30 * time.Hour, 72 * time.Hour
	case "News":
		if session == "regular" {
			return 5 * time.Minute, 7 * time.Minute, 15 * time.Minute
		}
		if session == "pre-market" || session == "after-hours" {
			return 10 * time.Minute, 15 * time.Minute, 25 * time.Minute
		}
		return 30 * time.Minute, 40 * time.Minute, 90 * time.Minute
	case "Earnings":
		return 2 * time.Hour, 150 * time.Minute, 4 * time.Hour
	case "SEC Filings":
		if active {
			return 15 * time.Minute, 20 * time.Minute, 35 * time.Minute
		}
		return 30 * time.Minute, 45 * time.Minute, 90 * time.Minute
	case "Fundamentals":
		return 24 * time.Hour, 26 * time.Hour, 36 * time.Hour
	case "Global":
		if active {
			return 7 * time.Minute, 10 * time.Minute, 20 * time.Minute
		}
		return 20 * time.Minute, 30 * time.Minute, 60 * time.Minute
	case "Macro":
		return 6 * time.Hour, 12 * time.Hour, 30 * time.Hour
	case "Options":
		if active {
			return 3 * time.Minute, 5 * time.Minute, 8 * time.Minute
		}
		return 10 * time.Minute, 15 * time.Minute, 30 * time.Minute
	}
	return 15 * time.Minute, 30 * time.Minute, time.Hour
}

func freshnessStateWithLimits(dataset, provider, session string, ts int64, health string, now int64, cadence, fresh, stale time.Duration) (string, string, int64, int64, int64) {
	age := int64(0)
	if ts > 0 {
		age = now - ts
		if age < 0 {
			age = 0
		}
	}
	lower := strings.ToLower(health)
	if strings.Contains(lower, "error") || strings.Contains(lower, "failed") {
		return "ERROR", "Provider/refresh error", age, int64(cadence / time.Millisecond), int64(stale / time.Millisecond)
	}
	if ts <= 0 {
		if strings.Contains(lower, "idle") || strings.Contains(lower, "stopped") {
			return "IDLE", "Not scheduled / inactive", 0, int64(cadence / time.Millisecond), int64(stale / time.Millisecond)
		}
		return "UNAVAILABLE", "No successful check/observation", 0, int64(cadence / time.Millisecond), int64(stale / time.Millisecond)
	}
	if dataset == "Quotes" && (session == "closed" || session == "weekend") {

		if age <= int64(96*time.Hour/time.Millisecond) {
			return "IDLE", "Market closed; recent valid quote retained", age, int64(cadence / time.Millisecond), int64(stale / time.Millisecond)
		}
		return "STALE", "Market closed but retained quote is older than the last plausible trading session", age, int64(cadence / time.Millisecond), int64(stale / time.Millisecond)
	}
	if dataset == "VIX" && provider == "CBOE" {

		if session != "regular" && age <= int64(96*time.Hour/time.Millisecond) {
			return "DELAYED", "Official CBOE delayed/close reference; not live", age, int64(cadence / time.Millisecond), int64(96 * time.Hour / time.Millisecond)
		}
		if age <= int64(stale/time.Millisecond) {
			return "DELAYED", "Valid delayed VIX source", age, int64(cadence / time.Millisecond), int64(stale / time.Millisecond)
		}
	}
	if dataset == "VIX" && provider == "yfinance" && age <= int64(stale/time.Millisecond) {
		return "DELAYED", "Valid delayed VIX recovery source", age, int64(cadence / time.Millisecond), int64(stale / time.Millisecond)
	}
	dueSoon := int64(float64(cadence/time.Millisecond) * .8)
	if age >= dueSoon && age <= int64(fresh/time.Millisecond) {
		return "DUE SOON", "Approaching next targeted refresh", age, int64(cadence / time.Millisecond), int64(stale / time.Millisecond)
	}
	if age <= int64(fresh/time.Millisecond) {
		if dataset == "Quotes" || dataset == "VIX" {
			return "LIVE", "Within live freshness threshold", age, int64(cadence / time.Millisecond), int64(stale / time.Millisecond)
		}
		return "FRESH", "Within expected freshness threshold", age, int64(cadence / time.Millisecond), int64(stale / time.Millisecond)
	}
	if age <= int64(stale/time.Millisecond) {
		return "DELAYED", "Past ideal cadence; targeted recovery is due", age, int64(cadence / time.Millisecond), int64(stale / time.Millisecond)
	}
	return "STALE", "Successful check/observation older than dataset-specific threshold", age, int64(cadence / time.Millisecond), int64(stale / time.Millisecond)
}

func freshnessState(dataset, provider, session string, ts int64, health string, now int64) (string, string, int64, int64, int64) {
	cadence, fresh, stale := freshnessLimits(dataset, provider, session)
	return freshnessStateWithLimits(dataset, provider, session, ts, health, now, cadence, fresh, stale)
}

func normalizeObservationMs(ts int64) int64 {
	if ts <= 0 {
		return 0
	}
	if ts < 10_000_000_000 {
		return ts * 1000
	}
	return ts
}

func safeFreshnessTimestamp(ts, receipt, now int64) (int64, bool) {
	ts = normalizeObservationMs(ts)
	receipt = normalizeObservationMs(receipt)
	const maxFutureSkew = int64(30 * time.Second / time.Millisecond)
	if ts > now+maxFutureSkew {
		if receipt > 0 && receipt <= now+maxFutureSkew {
			return receipt, true
		}
		return 0, true
	}
	return ts, false
}
func parseObservationDateMs(s string) int64 {
	s = strings.TrimSpace(s)
	if s == "" {
		return 0
	}
	for _, layout := range []string{time.RFC3339, time.RFC3339Nano, "2006-01-02"} {
		if v, err := time.Parse(layout, s); err == nil {
			return v.UnixMilli()
		}
	}
	return 0
}

func (e *Engine) buildFreshnessDiagnostics(quotes map[string]Quote, last map[string]int64, health map[string]string, scopes ...[]string) ([]FreshnessDiagnostic, FreshnessSummary) {
	nowTime := time.Now()
	now := nowTime.UnixMilli()
	session := marketSessionET(nowTime)
	cacheAt := last["cache"]
	var quoteScope, intradayScope []string
	useScopedQuotes := len(scopes) > 0
	useScopedIntraday := len(scopes) > 1
	if useScopedQuotes {
		quoteScope = uniqueSymbols(scopes[0])
	}
	if useScopedIntraday {
		intradayScope = uniqueSymbols(scopes[1])
	}
	quoteProvider, quoteReceived := int64(0), last["quotes"]
	quoteMissing := 0
	if useScopedQuotes {
		quoteReceived = 0
		for _, sym := range quoteScope {
			q, ok := quotes[sym]
			p := normalizeObservationMs(q.ProviderTimestamp)
			rc := normalizeObservationMs(q.UpdatedAt)
			if !ok || q.Price <= 0 || (p == 0 && rc == 0) {
				quoteMissing++
				continue
			}
			if p == 0 {
				p = rc
			}
			if quoteProvider == 0 || (p > 0 && p < quoteProvider) {
				quoteProvider = p
			}
			if rc > 0 && (quoteReceived == 0 || rc < quoteReceived) {
				quoteReceived = rc
			}
		}
		if quoteReceived == 0 {
			quoteReceived = last["quotes"]
		}
	} else {
		for sym, q := range quotes {
			if sym == "VIX" {
				continue
			}
			quoteProvider = maxInt64(quoteProvider, normalizeObservationMs(q.ProviderTimestamp))
			quoteReceived = maxInt64(quoteReceived, normalizeObservationMs(q.UpdatedAt))
		}
	}
	vq := quotes["VIX"]
	vProvider := normalizeObservationMs(vq.ProviderTimestamp)
	vReceived := normalizeObservationMs(vq.UpdatedAt)
	if vReceived == 0 {
		vReceived = last["vix"]
	}
	intradayProvider, dailyProvider := int64(0), int64(0)
	intradayMissing := 0
	if useScopedIntraday {
		for _, sym := range intradayScope {
			sets := e.bars[sym]
			rows := sets["intraday"]
			if len(rows) == 0 {
				intradayMissing++
				continue
			}
			ts := normalizeObservationMs(rows[len(rows)-1].T)
			if intradayProvider == 0 || (ts > 0 && ts < intradayProvider) {
				intradayProvider = ts
			}
		}
	} else {
		for sym, sets := range e.bars {
			if sym == "VIX" {
				continue
			}
			if rows := sets["intraday"]; len(rows) > 0 {
				intradayProvider = maxInt64(intradayProvider, normalizeObservationMs(rows[len(rows)-1].T))
			}
		}
	}
	for sym, sets := range e.bars {
		if sym == "VIX" {
			continue
		}
		for _, name := range []string{"daily", "weekly"} {
			if rows := sets[name]; len(rows) > 0 {
				dailyProvider = maxInt64(dailyProvider, normalizeObservationMs(rows[len(rows)-1].T))
			}
		}
	}
	newsProvider := int64(0)
	for _, n := range e.news {
		newsProvider = maxInt64(newsProvider, normalizeObservationMs(n.Datetime))
	}
	filingProvider := int64(0)
	for _, f := range e.filings {
		filingProvider = maxInt64(filingProvider, parseObservationDateMs(f.FiledAt))
	}
	globalProvider := int64(0)
	for _, g := range e.globalDirect {
		globalProvider = maxInt64(globalProvider, normalizeObservationMs(g.UpdatedAt))
	}
	macroProvider := int64(0)
	for _, m := range e.macroMetrics {
		macroProvider = maxInt64(macroProvider, normalizeObservationMs(m.UpdatedAt))
	}
	optionsProvider := int64(0)
	vp := sourceProvider(vq.Source)
	type def struct {
		Dataset, key, provider, fallback, action, basis, freshnessBasis string
		affected                                                        []string
		dataAt, receivedAt                                              int64
		sparse                                                          bool
	}
	defs := []def{
		{"Quotes", "quotes", "Live equity router", "Finnhub → Twelve Data", "quotes", "provider market timestamp", "market observation", []string{"Dashboard", "Day", "Swing", "Long", "Decision Queue"}, quoteProvider, quoteReceived, false},
		{"VIX", "vix", vp, "yfinance → CBOE", "vix", "provider/index observation timestamp", "market observation", []string{"Market Regime", "Dashboard", "Trade Readiness", "Pre-Market Prep", "Market Open Prep"}, vProvider, vReceived, false},
		{"Intraday Bars", "history", sourceProvider(health["history"]), "Twelve Data → yfinance", "history-intraday", "latest intraday bar", "completed-bar observation", []string{"Intraday Charts", "Day", "Trade Readiness"}, intradayProvider, maxInt64(last["history-intraday"], last["history"]), false},
		{"Daily / Weekly History", "history", sourceProvider(health["history"]), "Twelve Data → yfinance", "history-daily", "latest daily/weekly bar", "last successful history reconciliation", []string{"Swing", "Long", "Signal Validation"}, dailyProvider, maxInt64(last["history-daily"], last["history"]), true},
		{"News", "news", sourceProvider(health["news"]), "Marketaux", "news", "latest article publication", "last successful news check", []string{"Dashboard", "Research", "Catalyst Watch"}, newsProvider, last["news"], true},
		{"Earnings", "earnings", sourceProvider(health["earnings"]), "yfinance", "earnings", "latest sourced earnings event", "last successful earnings check", []string{"Dashboard", "Research", "Readiness", "Catalyst Watch"}, 0, last["earnings"], true},
		{"SEC Filings", "filings", "SEC EDGAR", "—", "filings", "latest filing date", "last successful EDGAR check", []string{"Research", "Long", "Day/Swing Catalyst Risk"}, filingProvider, last["filings"], true},
		{"Fundamentals", "fundamentals", sourceProvider(health["fundamentals"]), "SEC → yfinance", "fundamentals", "latest provider financial snapshot", "last successful fundamentals reconciliation", []string{"Research", "Long-Term"}, 0, last["fundamentals"], true},
		{"Global", "global-direct", "Provider Router", "Official/public/proxy", "global", "source observation timestamp", "market observation", []string{"Dashboard", "Market Regime", "Readiness"}, globalProvider, maxInt64(last["global-direct"], last["global"]), false},
		{"Macro", "macro", "FRED", "BLS/EIA context", "macro", "latest economic observation", "release-aware successful reconciliation", []string{"Global Context", "Regime", "Readiness"}, macroProvider, maxInt64(last["macro"], last["macro-rates"], last["fred-rates"], last["treasury"], last["bls-actuals"], last["eia-actuals"], last["bea-actuals"], last["macro-events"]), true},
		{"Options", "options", "Alpaca", "Indicative fallback", "options", "provider timestamp unavailable", "last successful options receipt", []string{"Research", "Readiness", "Decision Queue"}, optionsProvider, last["options"], true},
	}
	rows := make([]FreshnessDiagnostic, 0, len(defs))
	sum := FreshnessSummary{}
	for _, d := range defs {
		provider := d.provider
		if provider == "" || provider == "—" {
			provider = "—"
		}
		safeDataAt, providerClockAnomaly := safeFreshnessTimestamp(d.dataAt, d.receivedAt, now)
		safeReceivedAt, receiptClockAnomaly := safeFreshnessTimestamp(d.receivedAt, 0, now)
		stateTs := safeDataAt
		if d.sparse {
			stateTs = safeReceivedAt
		}
		if stateTs == 0 {
			stateTs = safeReceivedAt
		}
		var st, reason string
		var age, cadence, stale int64
		if d.Dataset == "Earnings" {
			ci := earningsRefreshIntervalFrom(e.earnings, nowTime)
			st, reason, age, cadence, stale = freshnessStateWithLimits(d.Dataset, provider, session, stateTs, health[d.key], now, ci, time.Duration(float64(ci)*1.25), ci*2)
		} else {
			st, reason, age, cadence, stale = freshnessState(d.Dataset, provider, session, stateTs, health[d.key], now)
		}

		if d.Dataset == "Quotes" && useScopedQuotes && quoteMissing > 0 && st != "ERROR" {
			st = "STALE"
			reason = fmt.Sprintf("%d of %d active/Research symbol quote(s) missing; targeted recovery required", quoteMissing, len(quoteScope))
		}
		if d.Dataset == "Intraday Bars" && useScopedIntraday && intradayMissing > 0 && st != "ERROR" {
			st = "STALE"
			reason = fmt.Sprintf("%d of %d Day symbol intraday history set(s) missing; targeted recovery required", intradayMissing, len(intradayScope))
		}
		dataAge := int64(-1)
		if safeDataAt > 0 {
			dataAge = now - safeDataAt
			if dataAge < 0 {
				dataAge = 0
			}
		}
		checkAge := int64(-1)
		if safeReceivedAt > 0 {
			checkAge = now - safeReceivedAt
			if checkAge < 0 {
				checkAge = 0
			}
		}
		if d.sparse && d.receivedAt > 0 && d.dataAt > 0 && dataAge > checkAge+int64(5*time.Minute/time.Millisecond) {
			switch d.Dataset {
			case "News":
				reason = "Successfully checked; no newer mapped article"
			case "SEC Filings":
				reason = "Successfully checked EDGAR; no newer filing"
			case "Macro":
				reason = "Current for latest scheduled release; awaiting next release"
			case "Daily / Weekly History":
				reason = "History reconciliation current; latest completed bar retained"
			default:
				reason = "Successful reconciliation current; source observation unchanged"
			}
		}
		if providerClockAnomaly || receiptClockAnomaly {
			if reason != "" {
				reason += " · "
			}
			if providerClockAnomaly && safeDataAt > 0 {
				reason += "provider timestamp is ahead of local clock; freshness uses DE.PULSE receipt"
			} else {
				reason += "future timestamp rejected; reconciliation required"
			}
		}
		_, f, _ := freshnessLimits(d.Dataset, provider, session)
		if d.Dataset == "Earnings" {
			f = time.Duration(float64(earningsRefreshIntervalFrom(e.earnings, nowTime)) * 1.25)
		}
		fresh := int64(f / time.Millisecond)
		next := int64(0)
		if d.receivedAt > 0 {
			next = d.receivedAt + cadence
		} else if stateTs > 0 {
			next = stateTs + cadence
		}
		if stateTs <= 0 {
			age = -1
		}
		row := FreshnessDiagnostic{Dataset: d.Dataset, State: st, Provider: provider, ProviderTimestamp: d.dataAt, ReceivedAt: d.receivedAt, DataTimestamp: d.dataAt, CacheAt: cacheAt, TimestampBasis: d.basis, FreshnessBasis: d.freshnessBasis, AgeMs: age, CheckAgeMs: checkAge, DataAgeMs: dataAge, ExpectedCadenceMs: cadence, FreshLimitMs: fresh, StaleLimitMs: stale, NextExpectedAt: next, Reason: reason, Fallback: d.fallback, Affected: d.affected, Session: session, Action: d.action}
		rows = append(rows, row)
		switch st {
		case "LIVE":
			sum.Live++
		case "FRESH":
			sum.Fresh++
		case "DUE SOON":
			sum.DueSoon++
		case "DELAYED":
			sum.Delayed++
		case "STALE":
			sum.Stale++
		case "ERROR":
			sum.Error++
		case "UNAVAILABLE":
			sum.Unavailable++
		case "IDLE":
			sum.Idle++
		}
	}
	return rows, sum
}

func enrichProviderCapabilityRegistry(rows []ProviderCapabilityEntry, router ProviderRouterSnapshot) []ProviderCapabilityEntry {
	by := map[string][]struct {
		dataset string
		hop     ProviderRouteHop
	}{}
	for _, r := range router.Routes {
		for _, h := range r.Route {
			by[h.Provider] = append(by[h.Provider], struct {
				dataset string
				hop     ProviderRouteHop
			}{r.Dataset, h})
		}
	}
	for i := range rows {
		items := by[rows[i].Provider]
		if len(items) == 0 && rows[i].Provider == "SEC EDGAR" {
			items = by["SEC"]
		}
		if len(items) == 0 {
			continue
		}
		sort.Slice(items, func(a, b int) bool { return items[a].hop.Priority < items[b].hop.Priority })
		h := items[0].hop
		pr := []string{}
		delays := []string{}
		for _, x := range items {
			pr = append(pr, fmt.Sprintf("%s #%d", x.dataset, x.hop.Priority))
			if x.hop.ExpectedDelay != "" {
				delays = append(delays, x.hop.ExpectedDelay)
			}
		}
		rows[i].Priority = strings.Join(pr, " · ")
		rows[i].Quota = h.Quota
		rows[i].RateLimit = h.RateLimit
		rows[i].LatencyMs = h.LatencyMs
		rows[i].LastSuccess = h.LastSuccess
		rows[i].LastFailure = h.LastFailure
		rows[i].FailureCount = h.FailureCount
		rows[i].Attempts = h.Attempts
		rows[i].LastError = h.LastError
		rows[i].Recovery = h.Recovery
		seenDelay := map[string]bool{}
		uniqDelay := []string{}
		for _, d := range delays {
			if d != "" && !seenDelay[d] {
				seenDelay[d] = true
				uniqDelay = append(uniqDelay, d)
			}
		}
		rows[i].ExpectedDelay = strings.Join(uniqDelay, " · ")
		rows[i].CostClass = h.CostClass
	}
	return rows
}
