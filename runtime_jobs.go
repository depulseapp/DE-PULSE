package main

import (
	"context"
	"fmt"
	"math"
	"strings"
	"time"
)

func minPositive(a, b float64) float64 {
	if a <= 0 {
		return b
	}
	return math.Min(a, b)
}

func (e *Engine) startLive(ctx context.Context, key string) {
	e.mu.Lock()
	e.health = map[string]string{"quotes": "connecting", "quotes-rest": "loading", "alpaca-stream": "checking", "alpaca-live": "checking", "vix": "detecting", "history": "loading", "news": "loading", "earnings": "loading", "filings": "loading", "fundamentals": "loading", "global": "loading", "global-direct": "checking", "fx-direct": "checking", "macro-rates": "loading", "fred-rates": "checking", "treasury": "loading", "bls-actuals": "loading", "bea-actuals": "loading", "eia-actuals": "checking", "macro-events": "loading", "event-mode": "ready", "options": "checking", "signal-validation": "ready", "cache-refresh": "ready", "pre-market-prep": "READY · awaiting scheduled window", "market-open-prep": "READY · awaiting scheduled window", "session-intelligence-coordinator": "READY · scheduler active", "catalyst-watch": "READY · event driven", "market-calendar": "checking", "market-activity": "checking", "corporate-actions": "checking", "symbol-intelligence": "checking", "provider-capabilities": "checking", "ai": "ready"}
	e.mu.Unlock()
	e.app.broadcastRuntime()
	e.app.mu.RLock()
	alpacaKey := strings.TrimSpace(e.app.secrets.AlpacaKey)
	alpacaSecret := strings.TrimSpace(e.app.secrets.AlpacaSecret)
	fredKey := strings.TrimSpace(e.app.secrets.FRED)
	eiaKey := strings.TrimSpace(e.app.secrets.EIA)
	twelveKey := strings.TrimSpace(e.app.secrets.TwelveData)
	e.app.mu.RUnlock()
	if strings.TrimSpace(key) != "" {
		e.wg.Add(1)
		go func() { defer e.wg.Done(); e.liveQuoteLoop(ctx, key) }()
	} else {
		e.setHealth("quotes", "Finnhub not configured · Alpaca/Twelve routing remains available")
	}
	e.wg.Add(1)
	go func() { defer e.wg.Done(); e.periodic(ctx, time.Minute, func() { e.refreshQuotesRouted(ctx, key) }) }()
	e.wg.Add(1)
	go func() { defer e.wg.Done(); e.periodic(ctx, 2*time.Minute, func() { e.refreshVIXSnapshot(ctx, key) }) }()
	e.wg.Add(1)
	go func() { defer e.wg.Done(); e.sessionAwareNewsLoop(ctx, key) }()
	e.wg.Add(1)
	go func() { defer e.wg.Done(); e.sessionAwareEarningsLoop(ctx, key) }()
	e.wg.Add(1)
	go func() { defer e.wg.Done(); e.sessionAwareSECLoop(ctx) }()
	e.wg.Add(1)
	go func() { defer e.wg.Done(); e.periodic(ctx, 24*time.Hour, func() { e.refreshFundamentals(ctx, key) }) }()
	e.wg.Add(1)
	go func() {
		defer e.wg.Done()
		e.periodic(ctx, 6*time.Hour, func() { e.refreshMacroRouted(ctx, fredKey, eiaKey) })
	}()
	e.wg.Add(1)
	go func() {
		defer e.wg.Done()
		e.periodic(ctx, 30*time.Minute, func() { e.refreshOfficialGlobalCloses(ctx) })
	}()
	e.wg.Add(1)
	go func() {
		defer e.wg.Done()
		e.periodic(ctx, 7*time.Minute, func() { e.refreshDirectGlobal(ctx, twelveKey) })
	}()
	e.wg.Add(1)
	go func() { defer e.wg.Done(); e.periodic(ctx, 6*time.Hour, func() { e.refreshMacroEvents(ctx) }) }()
	e.wg.Add(1)
	go func() { defer e.wg.Done(); e.eventModeLoop(ctx, alpacaKey, alpacaSecret) }()

	e.wg.Add(1)
	go func() { defer e.wg.Done(); e.sessionAwareHistoryLoop(ctx) }()
	if alpacaKey != "" && alpacaSecret != "" {
		e.wg.Add(1)
		go func() { defer e.wg.Done(); e.opportunityRadarLoop(ctx, alpacaKey, alpacaSecret) }()
		e.wg.Add(1)
		go func() { defer e.wg.Done(); e.alpacaIEXLoop(ctx, alpacaKey, alpacaSecret) }()
		e.wg.Add(1)
		go func() {
			defer e.wg.Done()
			e.periodic(ctx, 15*time.Second, func() { e.refreshAlpacaLiveSnapshots(ctx, alpacaKey, alpacaSecret) })
		}()
		e.wg.Add(1)
		go func() { defer e.wg.Done(); e.sessionAwareOptionsLoop(ctx, alpacaKey, alpacaSecret) }()
	} else {
		e.setHealth("history", "router ready · Alpaca unavailable; Twelve Data / yfinance recovery eligible")
		e.setHealth("alpaca-stream", "setup required · IEX fallback unavailable")
		e.setHealth("alpaca-live", "setup required · overnight unavailable")
	}
	e.wg.Add(1)
	go func() {
		defer e.wg.Done()
		t := time.NewTicker(5 * time.Second)
		defer t.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			case <-t.C:
				e.refreshFeedStatus()
			}
		}
	}()
	e.wg.Add(1)
	go func() { defer e.wg.Done(); e.autoFreshnessRecoveryLoop(ctx) }()
	e.wg.Add(1)
	go func() { defer e.wg.Done(); e.cachePersistLoop(ctx) }()
	if alpacaKey != "" && alpacaSecret != "" {
		e.wg.Add(1)
		go func() {
			defer e.wg.Done()
			e.periodic(ctx, 12*time.Hour, func() { e.refreshAlpacaMarketCalendar(ctx, alpacaKey, alpacaSecret) })
		}()
		e.wg.Add(1)
		go func() {
			defer e.wg.Done()
			e.periodic(ctx, 15*time.Minute, func() { e.refreshAlpacaMarketActivity(ctx, alpacaKey, alpacaSecret) })
		}()
		e.wg.Add(1)
		go func() {
			defer e.wg.Done()
			e.periodic(ctx, 12*time.Hour, func() { e.refreshAlpacaCorporateActions(ctx, alpacaKey, alpacaSecret) })
		}()
	}
	if key != "" {
		e.wg.Add(1)
		go func() {
			defer e.wg.Done()
			e.periodic(ctx, 6*time.Hour, func() { e.refreshFinnhubIntelligence(ctx, key) })
		}()
	}
	if twelveKey != "" {
		e.wg.Add(1)
		go func() {
			defer e.wg.Done()
			e.periodic(ctx, 10*time.Minute, func() { e.refreshTwelveFX(ctx, twelveKey) })
		}()

	}
	e.wg.Add(1)
	go func() { defer e.wg.Done(); e.sessionIntelligenceCoordinatorLoop(ctx, key, alpacaKey, alpacaSecret) }()
	e.wg.Add(1)
	go func() { defer e.wg.Done(); e.weeklyIntegrityLoop(ctx) }()
}

func (e *Engine) cachePersistLoop(ctx context.Context) {
	for {
		policy := e.currentAdaptivePolicy(time.Now())
		wait := time.Duration(policy.CachePersistCadence) * time.Millisecond
		if wait <= 0 {
			wait = 2 * time.Minute
		}
		t := time.NewTimer(wait)
		select {
		case <-ctx.Done():
			t.Stop()
			_ = e.saveCache()
			return
		case now := <-t.C:
			if e.saveCache() == nil {
				e.mu.Lock()
				e.lastUpdated["cache"] = now.UnixMilli()
				e.lastUpdated["adaptive-cache-policy"] = now.UnixMilli()
				e.mu.Unlock()
			}
		}
	}
}

func capWaitToSessionBoundary(now time.Time, desired time.Duration) time.Duration {
	if desired <= 0 {
		return time.Second
	}
	boundaryAt, _ := marketSessionBoundaryET(now)
	if boundaryAt <= 0 {
		return desired
	}
	until := time.UnixMilli(boundaryAt).Sub(now)
	if until <= 0 {
		return minDuration(desired, time.Second)
	}

	until += 250 * time.Millisecond
	if until < desired {
		return until
	}
	return desired
}

func minDuration(a, b time.Duration) time.Duration {
	if a < b {
		return a
	}
	return b
}

// sessionAwareNewsLoop keeps headlines fresh without treating overnight as a fully idle period.
// Active sessions use a 10-minute cadence; overnight/closed periods use 30 minutes.
func (e *Engine) sessionAwareNewsLoop(ctx context.Context, finnhubKey string) {
	e.refreshNews(ctx, finnhubKey)
	for {
		session := marketSessionET(time.Now())
		wait := 30 * time.Minute
		if session == "regular" {
			wait = 5 * time.Minute
		} else if session == "pre-market" || session == "after-hours" {
			wait = 10 * time.Minute
		}
		wait = capWaitToSessionBoundary(time.Now(), wait)
		t := time.NewTimer(wait)
		select {
		case <-ctx.Done():
			t.Stop()
			return
		case <-t.C:
			e.refreshNews(ctx, finnhubKey)
		}
	}
}

func (e *Engine) sessionAwareHistoryLoop(ctx context.Context) {
	e.refreshHistoryRoutedMode(ctx, nil, "all")
	lastFullIntraday := time.Now()
	for {
		now := time.Now()
		policy := e.currentAdaptivePolicy(now)
		wait := time.Duration(policy.IntradayHistoryCadence) * time.Millisecond
		if wait <= 0 {
			wait = 5 * time.Minute
		}
		wait = capWaitToSessionBoundary(now, wait)
		t := time.NewTimer(wait)
		select {
		case <-ctx.Done():
			t.Stop()
			return
		case <-t.C:
			current := time.Now()
			currentSession := marketSessionET(current)
			active := currentSession == "regular" || currentSession == "pre-market" || currentSession == "after-hours"
			if active {
				currentPolicy := e.currentAdaptivePolicy(current)
				if len(currentPolicy.HotSymbols) > 0 && current.Sub(lastFullIntraday) < 5*time.Minute {
					e.refreshHistoryRoutedMode(ctx, currentPolicy.HotSymbols, "intraday")
				} else {
					e.refreshHistoryRoutedMode(ctx, nil, "intraday")
					lastFullIntraday = current
				}
			}

			if e.refreshDue("history-daily", 24*time.Hour, current) || currentSession == "after-hours" && e.refreshDue("history-daily", 4*time.Hour, current) {
				e.refreshHistoryRoutedMode(ctx, nil, "daily")
			}
		}
	}
}

func (e *Engine) optionsRefreshInterval(now time.Time) time.Duration {
	session := marketSessionET(now)
	base := 10 * time.Minute
	if session == "regular" || session == "pre-market" || session == "after-hours" {
		base = 3 * time.Minute
	}
	e.mu.RLock()
	defer e.mu.RUnlock()
	for _, r := range e.catalystReactions {
		if r.TriggerAt > 0 && r.CompletedAt == 0 {
			return time.Minute
		}
	}
	today := now.In(easternLocation()).Format("2006-01-02")
	for _, er := range e.earnings {
		if er.Date == today {
			return time.Minute
		}
	}
	return base
}

func (e *Engine) sessionAwareOptionsLoop(ctx context.Context, alpacaKey, alpacaSecret string) {
	e.refreshOptions(ctx, alpacaKey, alpacaSecret)
	for {
		now := time.Now()
		wait := capWaitToSessionBoundary(now, e.optionsRefreshInterval(now))
		t := time.NewTimer(wait)
		select {
		case <-ctx.Done():
			t.Stop()
			return
		case <-t.C:
			e.refreshOptions(ctx, alpacaKey, alpacaSecret)
		}
	}
}

func (e *Engine) secRefreshInterval(now time.Time) time.Duration {
	session := marketSessionET(now)
	wait := 30 * time.Minute
	if session == "regular" || session == "pre-market" || session == "after-hours" {
		wait = 15 * time.Minute
	}

	e.mu.RLock()
	defer e.mu.RUnlock()
	for _, r := range e.catalystReactions {
		if r.CompletedAt == 0 && r.TriggerAt > 0 && now.Sub(time.UnixMilli(r.TriggerAt)) < 6*time.Hour {
			return 5 * time.Minute
		}
	}
	if p, ok := e.preparations["catalyst-watch"]; ok && (p.State == "ARMED" || p.State == "TRIGGERED" || p.State == "REACTION") {
		return 5 * time.Minute
	}
	return wait
}

func (e *Engine) sessionAwareSECLoop(ctx context.Context) {
	e.refreshSECRouted(ctx)
	for {
		now := time.Now()
		wait := capWaitToSessionBoundary(now, e.secRefreshInterval(now))
		t := time.NewTimer(wait)
		select {
		case <-ctx.Done():
			t.Stop()
			return
		case <-t.C:
			e.refreshSECRouted(ctx)
		}
	}
}

func earningsRefreshIntervalFrom(events []EarningsItem, now time.Time) time.Duration {
	interval := 2 * time.Hour
	today := now.In(easternLocation())
	for _, er := range events {
		d, err := time.ParseInLocation("2006-01-02", er.Date, easternLocation())
		if err != nil {
			continue
		}
		delta := d.Sub(time.Date(today.Year(), today.Month(), today.Day(), 0, 0, 0, 0, easternLocation()))
		if delta >= 0 && delta <= 24*time.Hour {
			interval = 30 * time.Minute
		}
		if er.Date == today.Format("2006-01-02") && !earningsReleased(er) {
			h := strings.ToLower(er.Hour)
			session := marketSessionET(now)
			if (strings.Contains(h, "bmo") && session == "pre-market") || (strings.Contains(h, "amc") && session == "after-hours") {
				interval = 7 * time.Minute
			}
		}
	}
	return interval
}

func (e *Engine) earningsRefreshInterval(now time.Time) time.Duration {
	e.mu.RLock()
	events := clone(e.earnings)
	e.mu.RUnlock()
	return earningsRefreshIntervalFrom(events, now)
}

func (e *Engine) sessionAwareEarningsLoop(ctx context.Context, finnhubKey string) {
	e.refreshEarnings(ctx, finnhubKey)
	for {
		now := time.Now()
		wait := capWaitToSessionBoundary(now, e.earningsRefreshInterval(now))
		t := time.NewTimer(wait)
		select {
		case <-ctx.Done():
			t.Stop()
			return
		case <-t.C:
			e.refreshEarnings(ctx, finnhubKey)
		}
	}
}

func freshnessPriority(dataset string) int {
	switch dataset {
	case "Quotes", "VIX", "Intraday Bars":
		return 1
	case "News", "SEC Filings", "Earnings":
		return 2
	case "Daily / Weekly History", "Fundamentals", "Global", "Macro", "Options":
		return 3
	}
	return 4
}

func freshnessRecoveryDue(row FreshnessDiagnostic, now int64) bool {
	state := strings.ToUpper(strings.TrimSpace(row.State))
	if row.Action == "" || state == "IDLE" {
		return false
	}
	switch state {
	case "STALE", "ERROR", "UNAVAILABLE":
		return true
	case "DELAYED":

		return row.NextExpectedAt > 0 && now >= row.NextExpectedAt
	}
	return row.NextExpectedAt > 0 && now >= row.NextExpectedAt
}

func (e *Engine) autoFreshnessRecoveryLoop(ctx context.Context) {
	t := time.NewTicker(30 * time.Second)
	defer t.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-t.C:
			snap := e.Snapshot()
			e.app.mu.RLock()
			sec := clone(e.app.secrets)
			e.app.mu.RUnlock()
			for _, row := range snap.Freshness {
				priority := freshnessPriority(row.Dataset)

				if !freshnessRecoveryDue(row, time.Now().UnixMilli()) {
					continue
				}
				minGap := time.Minute
				if priority == 2 {
					minGap = 3 * time.Minute
				}
				if priority >= 3 {
					minGap = 10 * time.Minute
				}
				key := "auto-recovery:" + row.Action
				if !e.refreshDue(key, minGap, time.Now()) {
					continue
				}
				recoveryCtx, tier := freshnessRecoveryWorkContext(ctx, priority)
				if e.workload.ShouldShed(tier) {
					e.mu.Lock()
					e.health["auto-recovery"] = "deferred · LOCAL LOAD · " + row.Dataset + " · " + workTierLabel(tier)
					e.mu.Unlock()
					continue
				}
				e.mu.Lock()
				e.lastUpdated[key] = time.Now().UnixMilli()
				e.health["auto-recovery"] = "refreshing · " + row.Dataset + " · " + workTierLabel(tier)
				e.mu.Unlock()
				rctx, cancel := context.WithTimeout(recoveryCtx, 40*time.Second)
				ok := e.refreshDatasetRouted(rctx, row.Action, sec)
				cancel()
				e.mu.Lock()
				if ok {
					e.health["auto-recovery"] = "recovered · " + row.Dataset + " · " + workTierLabel(tier)
				} else {
					e.health["auto-recovery"] = "stale · recovery failed · " + row.Dataset + " · " + workTierLabel(tier)
				}
				e.mu.Unlock()
				if priority <= 2 {
					break
				}
			}
		}
	}
}

// preMarketPrepLoop performs one targeted, low-priority preparation pass on US
// trading days between 3:15 and 3:50 AM ET. Overnight quote streams remain live;
// this job only refreshes stale research/history inputs and atomically persists
// validated results. It never clears the market cache.
func (e *Engine) runPreMarketPrep(ctx context.Context, finnhubKey, alpacaKey, alpacaSecret, reason string, late bool) {
	nowTime := time.Now()
	e.setPreparationRich("pre-market-prep", "RUNNING", reason, false, late, "", nil, nil, nil)
	e.mu.RLock()
	before := clone(e.preparations["pre-market-prep"])
	e.mu.RUnlock()

	if e.mode != "demo" {
		e.runLowPriorityRefresh(ctx, finnhubKey, alpacaKey, alpacaSecret, reason)
	}
	e.mu.RLock()
	quotes := clone(e.quotes)
	tracked := len(e.trackedSymbols())
	e.mu.RUnlock()

	critical := []string{"Quotes", "VIX", "Intraday Bars", "News", "Earnings", "SEC Filings"}
	freshSnap := e.Snapshot()
	healthy, degraded, exceptions := readinessFreshnessGate(freshSnap.Freshness, critical, nowTime)
	vixText := "VIX unavailable"
	if q := quotes["VIX"]; q.Price > 0 {
		vixText = fmt.Sprintf("VIX %.2f · %s", q.Price, defaultString(q.Source, q.FeedType))
	}
	summary := []string{
		fmt.Sprintf("%d/%d critical datasets usable", healthy, len(critical)),
		fmt.Sprintf("%d canonical symbols prepared", tracked),
		vixText,
		fmt.Sprintf("%d exception(s) require attention", len(exceptions)),
	}
	changed := []string{}
	if len(before.Summary) > 0 && strings.Join(before.Summary, "|") != strings.Join(summary, "|") {
		changed = append(changed, "Critical data/freshness summary changed since the prior checkpoint")
	}
	degraded = degraded || healthy < len(critical)
	attention := checkpointAttention(exceptions, degraded)
	detail := fmt.Sprintf("%s · %d/%d critical datasets usable · %d exception(s)", reason, healthy, len(critical), len(exceptions))
	if late {
		detail = "LATE · " + detail
	}
	e.mu.Lock()
	e.lastUpdated["pre-market-prep"] = time.Now().UnixMilli()
	e.mu.Unlock()
	e.setPreparationRich("pre-market-prep", attention, detail, !degraded, late, attention, summary, changed, exceptions)
	_ = e.saveCache()
	e.app.broadcastRuntime()
}

// preMarketPrepLoop performs one targeted, low-priority preparation pass on US
// trading days between 3:15 and 3:50 AM ET. If the app starts later that morning
// and today's checkpoint did not run, one catch-up pass runs before Market Open Prep.
func (e *Engine) preMarketPrepLoop(ctx context.Context, finnhubKey, alpacaKey, alpacaSecret string) {
	t := time.NewTicker(5 * time.Minute)
	defer t.Stop()
	run := func(now time.Time) {
		if !e.isTradingDay(now) {
			return
		}
		et := now.In(easternLocation())
		mins := et.Hour()*60 + et.Minute()
		e.mu.RLock()
		p := e.preparations["pre-market-prep"]
		e.mu.RUnlock()
		if preparationRanForTradingDay(p, now) {
			return
		}
		if preMarketPrepWindow(now) {
			if preparationRetryDue(p, now, 10*time.Minute) {
				e.runPreMarketPrep(ctx, finnhubKey, alpacaKey, alpacaSecret, "scheduled", false)
			}
			return
		}
		if mins > 3*60+50 && mins < 9*60+20 && preparationRetryDue(p, now, 20*time.Minute) {
			e.runPreMarketPrep(ctx, finnhubKey, alpacaKey, alpacaSecret, "missed-window catch-up", true)
		}
	}
	run(time.Now())
	for {
		select {
		case <-ctx.Done():
			return
		case now := <-t.C:
			run(now)
		}
	}
}

// weeklyIntegrityLoop performs non-destructive cache/integrity housekeeping on
// Saturday. It validates in-memory structures and persists the current good state;
// user watchlists/settings and live provider configuration are never modified.
func (e *Engine) weeklyIntegrityLoop(ctx context.Context) {
	t := time.NewTicker(30 * time.Minute)
	defer t.Stop()
	lastWeek := ""
	run := func(now time.Time) {
		loc, err := time.LoadLocation("America/New_York")
		if err != nil {
			loc = time.FixedZone("ET", -5*60*60)
		}
		et := now.In(loc)
		y, w := et.ISOWeek()
		weekKey := fmt.Sprintf("%04d-%02d", y, w)
		if !weeklyIntegrityWindow(now) || weekKey == lastWeek {
			return
		}
		if e.shouldShedBackground() {
			e.setHealth("weekly-integrity", "deferred · LOCAL LOAD")
			return
		}
		release, ok := e.workload.TryAcquireTier("background", WorkTierBackground)
		if !ok {
			e.setHealth("weekly-integrity", "deferred · background worker busy")
			return
		}
		lastWeek = weekKey
		e.runWeeklyIntegrity("scheduled")
		release()
	}
	for {
		select {
		case <-ctx.Done():
			return
		case now := <-t.C:
			run(now)
		}
	}
}

func (e *Engine) runWeeklyIntegrity(reason string) {
	e.mu.Lock()
	if e.quotes == nil {
		e.quotes = map[string]Quote{}
	}
	if e.history == nil {
		e.history = map[string][]HistoryPoint{}
	}
	if e.bars == nil {
		e.bars = map[string]map[string][]Bar{}
	}
	if e.fundamentals == nil {
		e.fundamentals = map[string]FundamentalSnapshot{}
	}
	if e.secIntelligence == nil {
		e.secIntelligence = map[string]SECIntelligenceSummary{}
	}
	now := time.Now().UnixMilli()
	e.lastUpdated["weekly-integrity"] = now
	e.health["weekly-integrity"] = "healthy · non-destructive · " + reason
	e.mu.Unlock()
	_ = e.saveCache()
	e.app.broadcastRuntime()
}

func (e *Engine) refreshDue(key string, maxAge time.Duration, now time.Time) bool {
	e.mu.RLock()
	last := e.lastUpdated[key]
	e.mu.RUnlock()
	if last <= 0 {
		return true
	}
	return now.Sub(time.UnixMilli(last)) >= maxAge
}

func preMarketPrepWindow(now time.Time) bool {
	loc, err := time.LoadLocation("America/New_York")
	if err != nil {
		loc = time.FixedZone("ET", -5*60*60)
	}
	et := now.In(loc)
	mins := et.Hour()*60 + et.Minute()
	return et.Weekday() >= time.Monday && et.Weekday() <= time.Friday && !isUSMarketHoliday(et) && mins >= 3*60+15 && mins <= 3*60+50
}

func weeklyIntegrityWindow(now time.Time) bool {
	loc, err := time.LoadLocation("America/New_York")
	if err != nil {
		loc = time.FixedZone("ET", -5*60*60)
	}
	et := now.In(loc)
	return et.Weekday() == time.Saturday && et.Hour() >= 10
}

func (e *Engine) runLowPriorityRefresh(ctx context.Context, finnhubKey, alpacaKey, alpacaSecret, reason string) {
	ctx = withWorkTier(ctx, WorkTierBackground)
	if e.shouldShedBackground() {
		e.setHealth("cache-refresh", "deferred · LOCAL LOAD · critical work protected")
		return
	}
	release, ok := e.workload.TryAcquireTier("background", WorkTierBackground)
	if !ok {
		e.setHealth("cache-refresh", "deferred · background worker busy")
		return
	}
	defer release()
	e.setHealth("cache-refresh", "refreshing · "+reason)
	now := time.Now()
	refreshed := 0
	skipped := 0
	if e.refreshDue("history", 5*time.Minute, now) {
		e.refreshHistoryRouted(ctx, nil)
		refreshed++
	} else {
		skipped++
	}

	if e.refreshDue("fundamentals", 24*time.Hour, now) {
		e.refreshFundamentals(ctx, finnhubKey)
		refreshed++
	} else {
		skipped++
	}
	if e.refreshDue("earnings", 2*time.Hour, now) {
		e.refreshEarnings(ctx, finnhubKey)
		refreshed++
	} else {
		skipped++
	}
	if finnhubKey != "" || strings.TrimSpace(e.app.secrets.Marketaux) != "" {
		newsTTL := 30 * time.Minute
		session := marketSessionET(now)
		if session == "regular" || session == "pre-market" || session == "after-hours" {
			newsTTL = 10 * time.Minute
		}
		if e.refreshDue("news", newsTTL, now) {
			e.refreshNews(ctx, finnhubKey)
			refreshed++
		} else {
			skipped++
		}
	}
	if e.refreshDue("filings", 30*time.Minute, now) {
		e.refreshSECRouted(ctx)
		refreshed++
	} else {
		skipped++
	}

	e.app.mu.RLock()
	fredKey := strings.TrimSpace(e.app.secrets.FRED)
	eiaKey := strings.TrimSpace(e.app.secrets.EIA)
	e.app.mu.RUnlock()
	if e.refreshDue("fred-rates", 12*time.Hour, now) || e.refreshDue("macro-rates", 12*time.Hour, now) || e.refreshDue("bls-actuals", 24*time.Hour, now) || e.refreshDue("eia-actuals", 24*time.Hour, now) {
		e.refreshMacroRouted(ctx, fredKey, eiaKey)
		refreshed++
	} else {
		skipped++
	}
	_ = e.saveCache()
	stamp := time.Now().UnixMilli()
	e.mu.Lock()
	e.lastUpdated["cache-refresh"] = stamp
	warnings := false
	for _, key := range []string{"history", "fundamentals", "earnings", "news", "filings"} {
		v := strings.ToLower(e.health[key])
		if strings.Contains(v, "degraded") || strings.Contains(v, "error") || strings.Contains(v, "failed") {
			warnings = true
			break
		}
	}
	if warnings {
		e.health["cache-refresh"] = fmt.Sprintf("degraded · cached fallbacks · %s · %d refreshed · %d current", reason, refreshed, skipped)
	} else {
		e.health["cache-refresh"] = fmt.Sprintf("healthy · %s · %d refreshed · %d already current", reason, refreshed, skipped)
	}
	e.mu.Unlock()
	e.app.broadcastRuntime()
}
