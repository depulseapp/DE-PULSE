package main

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"
)

func (e *Engine) refreshSECRouted(ctx context.Context) bool {
	before := int64(0)
	e.mu.RLock()
	before = e.lastUpdated["filings"]
	e.mu.RUnlock()
	active, ok := e.executeProviderRoute(ctx, "SEC", map[string]providerRouteAttempt{
		"SEC EDGAR": func(ctx context.Context) bool {
			e.refreshFilings(ctx)
			e.mu.RLock()
			after := e.lastUpdated["filings"]
			health := strings.ToLower(e.health["filings"])
			e.mu.RUnlock()
			if after > before || strings.Contains(health, "healthy") {
				e.recordProviderSuccess("SEC EDGAR")
				return true
			}
			if strings.Contains(health, "degraded") || strings.Contains(health, "error") {
				e.recordProviderFailure("SEC EDGAR", fmt.Errorf("%s", health))
			}
			return false
		},
	})
	if ok {
		e.setHealth("sec-route", "active · "+active)
	}
	return ok
}

func (e *Engine) refreshMacroRouted(ctx context.Context, fredKey, eiaKey string) bool {
	active, ok := e.executeProviderRoute(ctx, "Macro", map[string]providerRouteAttempt{
		"FRED": func(ctx context.Context) bool {
			before := int64(0)
			e.mu.RLock()
			before = e.lastUpdated["fred-rates"]
			e.mu.RUnlock()
			e.refreshFRED(ctx, fredKey)
			e.mu.RLock()
			after := e.lastUpdated["fred-rates"]
			health := strings.ToLower(e.health["fred-rates"])
			e.mu.RUnlock()
			if after > before || strings.Contains(health, "healthy") {
				e.recordProviderSuccess("FRED")
				return true
			}
			if strings.Contains(health, "unavailable") || strings.Contains(health, "degraded") {
				e.recordProviderFailure("FRED", fmt.Errorf("%s", health))
			}
			return false
		},
	})

	e.refreshOfficialMacroActuals(ctx, eiaKey)
	if ok {
		e.setHealth("macro-route", "active · "+active)
	}
	return ok
}

func (e *Engine) refreshDemoDataset(dataset string) bool {
	dataset = strings.ToLower(strings.TrimSpace(dataset))
	valid := map[string]string{
		"vix": "vix", "history": "history", "history-intraday": "history-intraday", "history-daily": "history-daily", "news": "news", "earnings": "earnings",
		"filings": "filings", "fundamentals": "fundamentals", "global": "global",
		"macro": "macro", "options": "options", "quotes": "quotes",
	}
	key, ok := valid[dataset]
	if !ok {
		return false
	}
	if dataset == "vix" {
		e.ensureDemoSymbol("VIX")
	}
	if dataset == "history" || dataset == "history-intraday" || dataset == "history-daily" {
		e.requestHistoryHydration()
	}
	now := time.Now().UnixMilli()
	e.mu.Lock()
	e.lastUpdated[key] = now
	if dataset == "global" {
		e.lastUpdated["global-direct"] = now
	}
	e.health[key] = "demo · targeted refresh"
	e.mu.Unlock()
	return true
}

func (e *Engine) refreshQuotesRouted(ctx context.Context, finnhubKey string) {
	e.app.mu.RLock()
	ak := strings.TrimSpace(e.app.secrets.AlpacaKey)
	as := strings.TrimSpace(e.app.secrets.AlpacaSecret)
	e.app.mu.RUnlock()
	_, _ = e.executeProviderRoute(ctx, "US Live Equities", map[string]providerRouteAttempt{
		"Alpaca": func(ctx context.Context) bool {
			before := int64(0)
			e.mu.RLock()
			before = maxInt64(e.lastAlpacaStreamAt, e.lastAlpacaAt)
			e.mu.RUnlock()
			e.refreshAlpacaLiveSnapshots(ctx, ak, as)
			e.mu.RLock()
			after := maxInt64(e.lastAlpacaStreamAt, e.lastAlpacaAt)
			e.mu.RUnlock()
			return after > before
		},
	})

	e.refreshSnapshots(ctx, finnhubKey)
}

func fundamentalSnapshotUsable(f FundamentalSnapshot) bool {
	return f.MarketCap > 0 || f.PERatio != 0 || f.ForwardPERatio != 0 || f.RevenueGrowth != 0 || f.EPSGrowth != 0 || f.GrossMargin != 0 || f.OperatingMargin != 0 || f.ROE != 0 || f.NetMargin != 0 || f.DebtToEquity != 0 || f.CurrentRatio != 0 || f.FreeCashFlow != 0
}

func (e *Engine) researchPackageReadiness(symbol string) (bool, []string) {
	return e.researchPackageReadinessAt(symbol, time.Now())
}

// researchPackageReadinessAt contains the deterministic readiness rules while
// allowing session-boundary tests to supply a fixed clock. Production callers
// continue to use researchPackageReadiness, so this refactor changes no UI or
// runtime behavior.
func (e *Engine) researchPackageReadinessAt(symbol string, nowTime time.Time) (bool, []string) {
	symbol = normalizeSymbol(symbol)
	snap := e.Snapshot()
	issues := []string{}
	session := marketSessionET(nowTime)

	for _, name := range []string{"News", "Earnings", "Fundamentals"} {

		targetKey := ""
		targetMaxAge := 30 * time.Minute
		switch name {
		case "News":
			targetKey, targetMaxAge = "research-news:"+symbol, 15*time.Minute
		case "Earnings":
			targetKey, targetMaxAge = "research-earnings:"+symbol, 2*time.Hour
		case "Fundamentals":
			targetKey, targetMaxAge = "research-fundamentals:"+symbol, 24*time.Hour
		}
		ts := snap.LastUpdated[targetKey]
		checkAge, valid, _ := evidenceAge(nowTime.UnixMilli(), ts, 30*time.Second)
		if valid && checkAge <= int64(targetMaxAge/time.Millisecond) {
			continue
		}
		issues = append(issues, name+" selected-ticker reconciliation STALE/INVALID")
	}

	secCurrent := false
	targetSECHealth := strings.ToLower(strings.TrimSpace(snap.Health["research-sec:"+symbol]))
	targetSECFailed := strings.Contains(targetSECHealth, "setup required") || strings.Contains(targetSECHealth, "degraded") || strings.Contains(targetSECHealth, "unavailable") || strings.Contains(targetSECHealth, "error") || strings.Contains(targetSECHealth, "failed")
	if !targetSECFailed {
		if age, valid, _ := evidenceAge(nowTime.UnixMilli(), snap.LastUpdated["research-sec:"+symbol], 30*time.Second); valid && age <= int64((30*time.Minute)/time.Millisecond) {
			secCurrent = true
		}
	}
	if !secCurrent {
		issues = append(issues, "SEC Filings STALE/UNAVAILABLE")
	}

	q := snap.Quotes[symbol]
	providerAge, receiptAge, quoteCurrent := quoteEvidenceAges(q, nowTime.UnixMilli())
	quoteMaxAge := 3 * time.Minute
	switch session {
	case "overnight":
		quoteMaxAge = 90 * time.Minute
	case "weekend", "closed":
		quoteMaxAge = 96 * time.Hour
	}
	if q.Price <= 0 || q.ProviderTimestamp <= 0 && q.UpdatedAt <= 0 {
		issues = append(issues, "Quote unavailable for "+symbol)
	} else if !quoteCurrent || providerAge > int64(quoteMaxAge/time.Millisecond) || receiptAge > int64(quoteMaxAge/time.Millisecond) {
		issues = append(issues, "Quote stale/clock-skewed for "+symbol)
	}

	sets := snap.Bars[symbol]
	dailyRows := sets["daily"]
	if len(dailyRows) == 0 && len(snap.History[symbol]) == 0 {
		issues = append(issues, "Daily history unavailable for "+symbol)
	}
	intradayRows := sets["intraday"]
	if len(intradayRows) == 0 {
		issues = append(issues, "Intraday history unavailable for "+symbol)
	} else {
		barTs, _ := safeFreshnessTimestamp(intradayRows[len(intradayRows)-1].T, 0, nowTime.UnixMilli())
		barMaxAge := 30 * time.Minute
		switch session {
		case "overnight":
			barMaxAge = 90 * time.Minute
		case "weekend", "closed":
			barMaxAge = 96 * time.Hour
		}
		if barTs <= 0 || nowTime.Sub(time.UnixMilli(barTs)) > barMaxAge {
			issues = append(issues, "Intraday history stale for "+symbol)
		}
	}
	if f, ok := snap.Fundamentals[symbol]; !ok || !fundamentalSnapshotUsable(f) {
		issues = append(issues, "Fundamentals unavailable for "+symbol)
	} else {
		fAge, valid, _ := evidenceAge(nowTime.UnixMilli(), f.UpdatedAt, 30*time.Second)
		fLimit := int64((72 * time.Hour) / time.Millisecond)
		if session == "closed" || session == "weekend" {
			fLimit = int64((120 * time.Hour) / time.Millisecond)
		}
		if !valid || fAge > fLimit {
			issues = append(issues, "Fundamentals stale/clock-skewed for "+symbol)
		}
	}
	return len(issues) == 0, uniqueStrings(issues)
}

func (a *Application) handleResearchRefresh(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, 405, "Method not allowed")
		return
	}
	var in struct {
		Symbol string `json:"symbol"`
	}
	if decodeJSON(r, &in) != nil {
		writeError(w, 400, "Invalid request")
		return
	}
	sym, validSymbol := parseUserTicker(in.Symbol)
	if !validSymbol {
		writeError(w, 400, "Invalid research ticker")
		return
	}
	userID := requestUserID(r.Context())
	a.mu.Lock()
	workspaceState := a.workspaceStateLocked(userID)
	workspaceState.UI.SelectedTicker = sym
	if err := a.saveWorkspaceStateLocked(userID, workspaceState); err != nil {
		a.mu.Unlock()
		writeError(w, 500, err.Error())
		return
	}
	a.mu.Unlock()
	a.broadcastStateForUser(userID)

	a.engine.hydrateFocusQuote(sym)

	a.engine.mu.RLock()
	mode := a.engine.mode
	status := a.engine.status
	a.engine.mu.RUnlock()
	if mode == "demo" {
		a.engine.ensureDemoSymbol(sym)
		now := time.Now().UnixMilli()
		a.engine.mu.Lock()
		a.engine.lastUpdated["research:"+sym] = now
		for _, kind := range []string{"history", "fundamentals", "news", "earnings", "sec"} {
			a.engine.lastUpdated["research-"+kind+":"+sym] = now
		}
		a.engine.health["research-sec:"+sym] = "healthy · demo evidence"
		a.engine.mu.Unlock()
		writeJSON(w, 200, map[string]any{"ok": true, "symbol": sym, "ready": true, "runtime": a.engine.Snapshot(), "message": "Demo research package hydrated."})
		return
	}
	if status != "running" && status != "degraded" {
		writeJSON(w, 200, map[string]any{"ok": true, "symbol": sym, "ready": false, "runtime": a.engine.Snapshot(), "message": "Research target selected. Start the runtime to refresh sourced evidence."})
		return
	}
	ctx, cancel := context.WithTimeout(r.Context(), 55*time.Second)
	defer cancel()
	a.mu.RLock()
	sec := clone(a.secrets)
	a.mu.RUnlock()

	a.engine.refreshResearchHistory(ctx, sym)
	a.engine.refreshResearchFundamentals(ctx, sec.Finnhub, sym)
	a.engine.refreshResearchEarnings(ctx, sec.Finnhub, sym)
	a.engine.refreshResearchNews(ctx, sec.Finnhub, sym)
	a.engine.refreshSECResearchSymbol(ctx, sym)
	// Congressional disclosures are optional SHADOW alternative evidence. SEC
	// remains authoritative, and this result never participates in readiness.
	_ = a.engine.refreshTradeInsightCongressResearchSymbol(ctx, sym)
	ready, issues := a.engine.researchPackageReadiness(sym)
	now := time.Now().UnixMilli()
	a.engine.mu.Lock()
	a.engine.lastUpdated["research:"+sym] = now
	if ready {
		a.engine.lastUpdated["research-ready:"+sym] = now
	}
	a.engine.mu.Unlock()
	_ = a.engine.saveCache()
	snap := a.engine.Snapshot()
	a.broadcastRuntime()
	message := "Sourced Research package reconciled and ready for AI second opinion."
	if !ready {
		message = "Research reconciliation completed, but critical evidence is not fully current: " + strings.Join(issues, "; ")
	}
	writeJSON(w, 200, map[string]any{"ok": true, "symbol": sym, "ready": ready, "issues": issues, "runtime": snap, "message": message})
}

func (e *Engine) refreshDatasetRouted(ctx context.Context, dataset string, s Secrets) bool {
	dataset = strings.ToLower(strings.TrimSpace(dataset))
	updated := func(keys []string, fn func()) bool {
		e.mu.RLock()
		before := int64(0)
		for _, k := range keys {
			before = maxInt64(before, e.lastUpdated[k])
		}
		e.mu.RUnlock()
		fn()
		e.mu.RLock()
		after := int64(0)
		for _, k := range keys {
			after = maxInt64(after, e.lastUpdated[k])
		}
		e.mu.RUnlock()
		return after > before
	}
	switch dataset {
	case "vix":
		return e.refreshVIXSnapshot(ctx, s.Finnhub)
	case "history":
		return e.refreshHistoryRouted(ctx, nil)
	case "history-intraday":
		return e.refreshHistoryRoutedMode(ctx, nil, "intraday")
	case "history-daily":
		return e.refreshHistoryRoutedMode(ctx, nil, "daily")
	case "news":
		return updated([]string{"news"}, func() { e.refreshNews(ctx, s.Finnhub) })
	case "earnings":
		return updated([]string{"earnings"}, func() { e.refreshEarnings(ctx, s.Finnhub) })
	case "filings":
		return e.refreshSECRouted(ctx)
	case "fundamentals":
		return updated([]string{"fundamentals"}, func() { e.refreshFundamentals(ctx, s.Finnhub) })
	case "global":
		return updated([]string{"global-direct", "global-public", "global"}, func() {
			e.refreshDirectGlobal(ctx, s.TwelveData)
			e.refreshOfficialGlobalCloses(ctx)
		})
	case "macro":
		return e.refreshMacroRouted(ctx, s.FRED, s.EIA)
	case "options":
		return updated([]string{"options"}, func() { e.refreshOptions(ctx, s.AlpacaKey, s.AlpacaSecret) })
	case "quotes":
		return updated([]string{"quotes"}, func() { e.refreshQuotesRouted(ctx, s.Finnhub) })
	default:
		return false
	}
}

func freshnessDiagnosticForAction(rows []FreshnessDiagnostic, action string) (FreshnessDiagnostic, bool) {
	action = strings.ToLower(strings.TrimSpace(action))
	for _, r := range rows {
		if strings.ToLower(strings.TrimSpace(r.Action)) == action {
			return r, true
		}
	}
	return FreshnessDiagnostic{}, false
}

func targetedRefreshOutcome(before, after FreshnessDiagnostic, foundAfter, success bool) (bool, string) {
	providerChanged := foundAfter && after.ProviderTimestamp > 0 && after.ProviderTimestamp != before.ProviderTimestamp
	receivedChanged := foundAfter && after.ReceivedAt > 0 && after.ReceivedAt != before.ReceivedAt
	changed := providerChanged
	message := "Refresh completed; provider observation is unchanged."
	if providerChanged {
		message = "Refresh completed with a newer provider observation."
	} else if receivedChanged {
		basis := strings.ToLower(after.FreshnessBasis)
		if strings.Contains(basis, "successful") || strings.Contains(basis, "receipt") || strings.Contains(basis, "reconciliation") {
			message = "Refresh checked successfully; source observation is unchanged. Check Age reset while Data Age remains truthful."
		} else if before.ProviderTimestamp == 0 {
			changed = true
			message = "Refresh completed; provider timestamp was unavailable, so DE.PULSE receipt time advanced."
		} else {
			message = "Refresh completed; receipt time advanced while the provider observation remained unchanged."
		}
	}
	if !success || (foundAfter && (strings.EqualFold(after.State, "ERROR") || strings.EqualFold(after.State, "UNAVAILABLE"))) {
		return false, "Refresh failed or remains unavailable; the previous valid observation and age were preserved."
	}
	return changed, message
}

func (a *Application) handleTargetedRefresh(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, 405, "Method not allowed")
		return
	}
	var in struct {
		Dataset string `json:"dataset"`
	}
	if decodeJSON(r, &in) != nil {
		writeError(w, 400, "Invalid refresh request")
		return
	}
	dataset := strings.ToLower(strings.TrimSpace(in.Dataset))
	beforeSnap := a.engine.Snapshot()
	before, _ := freshnessDiagnosticForAction(beforeSnap.Freshness, dataset)
	a.engine.mu.RLock()
	mode := a.engine.mode
	a.engine.mu.RUnlock()
	if mode == "demo" {
		if !a.engine.refreshDemoDataset(dataset) {
			writeError(w, 400, "Unknown dataset")
			return
		}
		a.broadcastRuntime()
		afterSnap := a.engine.Snapshot()
		after, _ := freshnessDiagnosticForAction(afterSnap.Freshness, dataset)
		changed := after.ProviderTimestamp != before.ProviderTimestamp || after.ReceivedAt != before.ReceivedAt
		writeJSON(w, 200, map[string]any{"ok": true, "dataset": dataset, "changed": changed, "message": "Demo targeted refresh completed.", "before": before, "after": after, "runtime": afterSnap})
		return
	}
	ctx, cancel := context.WithTimeout(r.Context(), 50*time.Second)
	defer cancel()
	a.mu.RLock()
	s := clone(a.secrets)
	a.mu.RUnlock()
	success := a.engine.refreshDatasetRouted(ctx, dataset, s)
	if !success && !containsString([]string{"vix", "history", "history-intraday", "history-daily", "news", "earnings", "filings", "fundamentals", "global", "macro", "options", "quotes"}, dataset) {
		writeError(w, 400, "Unknown dataset")
		return
	}

	afterSnap := a.engine.Snapshot()
	after, foundAfter := freshnessDiagnosticForAction(afterSnap.Freshness, dataset)
	changed, message := targetedRefreshOutcome(before, after, foundAfter, success)
	a.broadcastRuntime()
	writeJSON(w, 200, map[string]any{"ok": success, "dataset": dataset, "changed": changed, "message": message, "before": before, "after": after, "runtime": afterSnap})
}
