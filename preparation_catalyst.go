package main

import (
	"context"
	"fmt"
	"math"
	"reflect"
	"sort"
	"strings"
	"time"
)

func (e *Engine) runMarketOpenPrep(reason string) {
	late := strings.Contains(strings.ToLower(reason), "late") || strings.Contains(strings.ToLower(reason), "catch-up")
	e.setPreparationRich("market-open-prep", "RUNNING", reason, false, late, "", nil, nil, nil)
	now := time.Now()
	e.mu.RLock()
	previousCheckpoint := clone(e.marketOpenCheckpoint)
	previousFlags := clone(e.marketOpenFlags)
	quotes := clone(e.quotes)
	bars := clone(e.bars)
	earnings := clone(e.earnings)
	filings := clone(e.filings)
	news := clone(e.news)
	options := clone(e.options)
	direct := clone(e.globalDirect)
	metrics := clone(e.macroMetrics)
	liquidityBaselines := clone(e.liquidityBaselines)
	e.mu.RUnlock()
	e.app.mu.RLock()
	settings := clone(e.app.state.Settings)
	st := clone(e.app.state)
	e.app.mu.RUnlock()

	liq := deriveLiquidityStatesWithContext(quotes, bars, liquidityBaselines, now)
	flags := map[string][]string{}
	today := now.In(easternLocation()).Format("2006-01-02")
	premarket := map[string]PremarketSnapshot{}
	for sym, q := range quotes {
		if sym == "VIX" || q.Price <= 0 {
			continue
		}
		var f []string
		if l, ok := liq[sym]; ok && currentLiquidityMarketRisk(l, now) {
			f = append(f, "LIQUIDITY RISK")
		}
		pm, hasPM := premarketSnapshotFromBars(sym, bars[sym]["intraday"], q, now)
		if hasPM {
			premarket[sym] = pm
		}
		if extendedAtMarketOpen(q, bars[sym]["daily"], bars[sym]["intraday"], pm, now) {
			f = append(f, "EXTENDED")
		}
		for _, er := range earnings {
			if strings.EqualFold(er.Symbol, sym) && er.Date == today {
				f = append(f, "EVENT RISK")
				break
			}
		}
		if !containsString(f, "EVENT RISK") {
			for _, fi := range filings {
				if strings.EqualFold(fi.Symbol, sym) && isRecentDate(fi.FiledAt, 1) && materialSECFilingForTradingRisk(fi) {
					f = append(f, "EVENT RISK")
					break
				}
			}
		}
		if !containsString(f, "EVENT RISK") {
			for _, n := range news {
				if containsSymbol(n.Symbols, sym) && n.Datetime > 0 && time.Since(time.Unix(n.Datetime, 0)) < 12*time.Hour && materialText(n.Headline+" "+n.Summary) {
					f = append(f, "EVENT RISK")
					break
				}
			}
		}
		if len(f) > 0 {
			flags[sym] = uniqueStrings(f)
		}
	}

	global := deriveGlobalMarketContext(quotes, direct, metrics, settings.GlobalProviderMode)
	derived := deriveIntelligenceStates(metrics)
	macroStates := map[string]string{}
	for k, v := range derived {
		macroStates[k] = v.State
	}
	vix, vixState := 0.0, "UNAVAILABLE"
	if q := quotes["VIX"]; q.Price > 0 {
		vix = q.Price
		vixState = "NORMAL"
		if vix >= 30 {
			vixState = "HIGH"
		} else if vix >= 20 {
			vixState = "ELEVATED"
		}
	}
	optionsCount := 0
	for _, o := range options {
		if o.Symbol != "" && o.State != "" {
			optionsCount++
		}
	}
	daySet := symbolSetForWatchlist(st, st.Settings.DayWatchlistID, st.Settings.DayEnabled)
	swingSet := symbolSetForWatchlist(st, st.Settings.SwingWatchlistID, st.Settings.SwingEnabled)
	longSet := symbolSetForWatchlist(st, st.Settings.LongWatchlistID, st.Settings.LongEnabled)
	var swingChanges, longChanges []string
	exceptions := []CheckpointException{}
	for sym, fs := range flags {
		if len(fs) == 0 {
			continue
		}
		severity := "MEDIUM"
		if containsString(fs, "EVENT RISK") || containsString(fs, "LIQUIDITY RISK") {
			severity = "HIGH"
		}
		exceptions = append(exceptions, CheckpointException{Symbol: sym, Reason: strings.Join(fs, " · "), Severity: severity, Target: "research", Source: "Market Open Prep", UpdatedAt: now.UnixMilli()})
		if swingSet[sym] && (containsString(fs, "EVENT RISK") || containsString(fs, "EXTENDED")) {
			swingChanges = append(swingChanges, sym)
		}
		if longSet[sym] && containsString(fs, "EVENT RISK") {
			longChanges = append(longChanges, sym)
		}
	}
	sort.Slice(exceptions, func(i, j int) bool {
		if exceptions[i].Severity == exceptions[j].Severity {
			return exceptions[i].Symbol < exceptions[j].Symbol
		}
		return exceptions[i].Severity == "HIGH"
	})
	sort.Strings(swingChanges)
	sort.Strings(longChanges)

	freshSnap := e.Snapshot()
	_, degraded, freshnessExceptions := readinessFreshnessGate(freshSnap.Freshness, []string{"Quotes", "VIX", "Intraday Bars", "News", "Earnings", "SEC Filings"}, now)
	exceptions = append(exceptions, freshnessExceptions...)
	attention := checkpointAttention(exceptions, degraded)
	changed := []string{}
	if previousCheckpoint.RunAt > 0 {
		if previousCheckpoint.VIXState != vixState || math.Abs(previousCheckpoint.VIX-vix) >= 1 {
			changed = append(changed, fmt.Sprintf("VIX changed %.2f/%s → %.2f/%s", previousCheckpoint.VIX, previousCheckpoint.VIXState, vix, vixState))
		}
		if previousCheckpoint.GlobalTone != global.Tone {
			changed = append(changed, "Global tone changed "+defaultString(previousCheckpoint.GlobalTone, "—")+" → "+defaultString(global.Tone, "—"))
		}
	}
	for sym, fs := range flags {
		if strings.Join(previousFlags[sym], "|") != strings.Join(fs, "|") {
			changed = append(changed, sym+" · "+strings.Join(fs, " · "))
		}
	}
	if len(changed) > 8 {
		changed = changed[:8]
	}
	checkpoint := MarketOpenCheckpoint{State: attention, Attention: attention, RunAt: now.UnixMilli(), Late: late, VIX: vix, VIXState: vixState, GlobalTone: global.Tone, MacroStates: macroStates, Premarket: premarket, OptionsContexts: optionsCount, DayCandidates: len(daySet), SwingContextChanges: swingChanges, LongContextChanges: longChanges, Changed: changed, Exceptions: exceptions,
		Detail: []string{
			fmt.Sprintf("Premarket session reconciled for %d symbols", len(premarket)),
			fmt.Sprintf("VIX %s · global %s · %d macro states", vixState, defaultString(global.Tone, "NEUTRAL"), len(macroStates)),
			fmt.Sprintf("%d options contexts · %d symbols require review", optionsCount, len(flags)),
		}}

	e.mu.Lock()
	e.marketOpenFlags = flags
	e.marketOpenCheckpoint = checkpoint
	e.lastUpdated["market-open-prep"] = now.UnixMilli()
	e.mu.Unlock()
	detail := fmt.Sprintf("%s · VIX/global/macro/premarket/options/readiness reconciled · %d flagged", reason, len(flags))
	if late {
		detail = "LATE · " + detail
	}
	e.setPreparationRich("market-open-prep", attention, detail, !degraded, late, attention, checkpoint.Detail, changed, exceptions)
	e.evaluateCatalystWatch(now)
	_ = e.saveCache()
	e.app.broadcastRuntime()
}

func (e *Engine) marketOpenPrepLoop(ctx context.Context) {
	t := time.NewTicker(time.Minute)
	defer t.Stop()
	run := func(now time.Time) {
		if !e.isTradingDay(now) {
			return
		}
		et := now.In(easternLocation())
		mins := et.Hour()*60 + et.Minute()
		e.mu.RLock()
		p := e.preparations["market-open-prep"]
		e.mu.RUnlock()
		if preparationRanForTradingDay(p, now) {
			return
		}
		if marketOpenPrepWindow(now) {
			if preparationRetryDue(p, now, 2*time.Minute) {
				e.runMarketOpenPrep("scheduled")
			}
			return
		}
		if mins > 9*60+25 && mins <= 10*60+15 && preparationRetryDue(p, now, 5*time.Minute) {
			e.runMarketOpenPrep("missed-window catch-up")
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

func containsString(xs []string, s string) bool {
	for _, x := range xs {
		if x == s {
			return true
		}
	}
	return false
}
func uniqueStrings(xs []string) []string {
	seen := map[string]bool{}
	out := []string{}
	for _, x := range xs {
		if x != "" && !seen[x] {
			seen[x] = true
			out = append(out, x)
		}
	}
	return out
}
func containsSymbol(xs []string, s string) bool {
	for _, x := range xs {
		if strings.EqualFold(x, s) {
			return true
		}
	}
	return false
}
func isRecentDate(raw string, days int) bool {
	if raw == "" {
		return false
	}
	t, err := time.Parse(time.RFC3339, raw)
	if err != nil {
		t, err = time.Parse("2006-01-02", raw)
	}
	return err == nil && time.Since(t) <= time.Duration(days)*24*time.Hour
}

var materialTerms = []string{"earnings", "guidance", "raises guidance", "cuts guidance", "offering", "dilution", "acquisition", "merger", "fda", "sec investigation", "lawsuit", "settlement", "ceo", "cfo", "contract", "8-k", "bankruptcy", "restructuring", "investor day", "regulatory", "recall"}

func materialText(s string) bool {
	s = strings.ToLower(s)
	for _, t := range materialTerms {
		if strings.Contains(s, t) {
			return true
		}
	}
	return false
}

func earningsWindowMatches(er EarningsItem, session string) bool {
	hour := strings.ToLower(er.Hour)
	if strings.Contains(hour, "bmo") || strings.Contains(hour, "before") || strings.Contains(hour, "pre") {
		return session == "pre-market"
	}
	if strings.Contains(hour, "amc") || strings.Contains(hour, "after") || strings.Contains(hour, "post") {
		return session == "after-hours"
	}
	return session == "pre-market" || session == "after-hours"
}

func earningsReleased(er EarningsItem) bool {
	return er.EPSActual != nil || er.RevenueActual != nil
}

func catalystCompletionAt(triggerAt int64) time.Time {
	if triggerAt <= 0 {
		return time.Time{}
	}
	loc := easternLocation()
	trigger := time.UnixMilli(triggerAt).In(loc)
	day := time.Date(trigger.Year(), trigger.Month(), trigger.Day(), 0, 0, 0, 0, loc)
	isTradingDay := func(d time.Time) bool {
		return d.Weekday() >= time.Monday && d.Weekday() <= time.Friday && !isUSMarketHoliday(d)
	}
	closeFor := func(d time.Time) time.Time {
		hour := 16
		if isUSEarlyClose(d) {
			hour = 13
		}
		return time.Date(d.Year(), d.Month(), d.Day(), hour, 15, 0, 0, loc)
	}

	regularClose := closeFor(day).Add(-15 * time.Minute)
	if isTradingDay(day) && !trigger.After(regularClose.Add(-60*time.Minute)) {
		return closeFor(day)
	}
	for i := 1; i <= 10; i++ {
		d := day.AddDate(0, 0, i)
		if isTradingDay(d) {
			return closeFor(d)
		}
	}
	return trigger.Add(36 * time.Hour)
}

func catalystComplete(now time.Time, triggerAt int64) bool {
	until := catalystCompletionAt(triggerAt)
	return !until.IsZero() && !now.Before(until)
}

func catalystPhase(now time.Time, triggerAt int64) string {
	if triggerAt <= 0 {
		return "ARMED"
	}
	et := now.In(easternLocation())
	session := marketSessionET(now)
	age := now.Sub(time.UnixMilli(triggerAt))
	if age < 0 {
		age = 0
	}
	if catalystComplete(now, triggerAt) {
		return "COMPLETE"
	}
	if session == "pre-market" {
		return "PREMARKET REACTION"
	}
	if session == "regular" {
		trigger := time.UnixMilli(triggerAt).In(easternLocation())
		open := time.Date(et.Year(), et.Month(), et.Day(), 9, 30, 0, 0, easternLocation())
		mins := now.Sub(open).Minutes()

		if trigger.Year() == et.Year() && trigger.YearDay() == et.YearDay() && marketSessionET(trigger) == "regular" {
			mins = now.Sub(trigger).Minutes()
		}
		if mins < 5 {
			return "OPENING REACTION"
		}
		if mins < 15 {
			return "5m"
		}
		if mins < 30 {
			return "15m"
		}
		if mins < 60 {
			return "30m"
		}
		if mins < 120 {
			return "60m"
		}
		return "SESSION REACTION"
	}
	if session == "after-hours" {
		return "SESSION REACTION"
	}
	return "TRIGGERED"
}

func reactionAt(history []HistoryPoint, targetMs int64, triggerPrice float64) float64 {
	if targetMs <= 0 || triggerPrice <= 0 {
		return 0
	}
	for _, p := range history {
		if p.T >= targetMs && p.P > 0 {
			return (p.P/triggerPrice - 1) * 100
		}
	}
	return 0
}

func relativeVolumeFromDaily(q Quote, daily []Bar) float64 {
	if q.Volume <= 0 || len(daily) == 0 {
		return 0
	}
	start := len(daily) - 20
	if start < 0 {
		start = 0
	}
	total := 0.0
	count := 0
	for _, b := range daily[start:] {
		if b.V > 0 {
			total += b.V
			count++
		}
	}
	if count == 0 || total <= 0 {
		return 0
	}
	return q.Volume / (total / float64(count))
}

func openingRangeState(price float64, rows []Bar, now time.Time) string {
	if price <= 0 {
		return "UNKNOWN"
	}
	loc := easternLocation()
	day := now.In(loc).Format("2006-01-02")
	h, l := 0.0, 0.0
	for _, b := range rows {
		t := time.Unix(b.T, 0).In(loc)
		mins := t.Hour()*60 + t.Minute()
		if t.Format("2006-01-02") != day || mins < 9*60+30 || mins > 10*60 {
			continue
		}
		if b.H > h {
			h = b.H
		}
		if l == 0 || (b.L > 0 && b.L < l) {
			l = b.L
		}
	}
	if h <= 0 || l <= 0 {
		return "PENDING"
	}
	if price > h {
		return "ABOVE OPENING RANGE"
	}
	if price < l {
		return "BELOW OPENING RANGE"
	}
	return "INSIDE OPENING RANGE"
}

func finalizeCompletedCatalysts(existing map[string]CatalystReactionState, now time.Time) map[string]CatalystReactionState {
	updated := clone(existing)
	for sym, old := range updated {
		if old.TriggerAt > 0 && old.CompletedAt == 0 && catalystComplete(now, old.TriggerAt) {
			old.Phase = "COMPLETE"
			old.CompletedAt = now.UnixMilli()
			old.UpdatedAt = now.UnixMilli()
			updated[sym] = old
		}
	}
	return updated
}

func catalystReactionSemanticCopy(in map[string]CatalystReactionState) map[string]CatalystReactionState {
	out := clone(in)
	for sym, r := range out {

		r.UpdatedAt = 0
		out[sym] = r
	}
	return out
}

func catalystWatchMateriallyChanged(before, after map[string]CatalystReactionState, beforePrep, afterPrep PreparationJobStatus) bool {
	if beforePrep.State != afterPrep.State || beforePrep.Attention != afterPrep.Attention ||
		beforePrep.TradingDay != afterPrep.TradingDay || beforePrep.Detail != afterPrep.Detail ||
		!reflect.DeepEqual(beforePrep.Summary, afterPrep.Summary) ||
		!reflect.DeepEqual(beforePrep.Exceptions, afterPrep.Exceptions) {
		return true
	}
	return !reflect.DeepEqual(catalystReactionSemanticCopy(before), catalystReactionSemanticCopy(after))
}

func (e *Engine) evaluateCatalystWatch(now time.Time) {
	et := now.In(easternLocation())
	day := et.Format("2006-01-02")
	yesterday := et.AddDate(0, 0, -1).Format("2006-01-02")
	session := marketSessionET(now)
	e.mu.RLock()
	earnings := clone(e.earnings)
	news := clone(e.news)
	filings := clone(e.filings)
	quotes := clone(e.quotes)
	bars := clone(e.bars)
	history := clone(e.history)
	liquidityBaselines := clone(e.liquidityBaselines)
	existing := clone(e.catalystReactions)
	prep := e.preparations["catalyst-watch"]
	prepBefore := prep
	e.mu.RUnlock()

	type triggerDef struct {
		typ, trigger string
		at           int64
	}
	armed := map[string]triggerDef{}
	candidates := map[string]triggerDef{}
	for _, er := range earnings {
		h := strings.ToLower(er.Hour)
		isAMC := strings.Contains(h, "amc") || strings.Contains(h, "after") || strings.Contains(h, "post")
		relevantDay := er.Date == day || (isAMC && er.Date == yesterday && (session == "pre-market" || session == "regular"))
		if !relevantDay {
			continue
		}
		sym := normalizeSymbol(er.Symbol)
		if sym == "" {
			continue
		}
		if er.Date == day && earningsWindowMatches(er, session) && !earningsReleased(er) {
			armed[sym] = triggerDef{"EARNINGS", "Scheduled earnings · awaiting release confirmation", 0}
		}
		if earningsReleased(er) {
			triggerAt := now.UnixMilli()
			if old, ok := existing[sym]; ok && old.TriggerType == "EARNINGS" && old.TriggerAt > 0 {
				triggerAt = old.TriggerAt
			}
			candidates[sym] = triggerDef{"EARNINGS", "Earnings release confirmed", triggerAt}
		}
	}
	cutoff := now.Add(-2 * time.Hour).Unix()
	for _, n := range news {
		if n.Datetime < cutoff || !materialText(n.Headline+" "+n.Summary) {
			continue
		}
		for _, s := range n.Symbols {
			sym := normalizeSymbol(s)
			if sym == "" {
				continue
			}
			typ := "NEWS"
			if _, ok := armed[sym]; ok && (strings.Contains(strings.ToLower(n.Headline+" "+n.Summary), "earn") || strings.Contains(strings.ToLower(n.Headline+" "+n.Summary), "result")) {
				typ = "EARNINGS"
			}
			candidates[sym] = triggerDef{typ, n.Headline, n.Datetime * 1000}
		}
	}
	for _, f := range filings {
		if !isRecentDate(f.FiledAt, 1) || !materialSECFilingForTradingRisk(f) {
			continue
		}
		sym := normalizeSymbol(f.Symbol)
		if sym == "" {
			continue
		}
		candidates[sym] = triggerDef{"SEC", strings.TrimSpace(f.Form + " · " + f.Meaning), now.UnixMilli()}
	}

	for sym, old := range existing {
		if old.TriggerAt <= 0 || old.CompletedAt > 0 || catalystComplete(now, old.TriggerAt) {
			continue
		}
		if _, ok := candidates[sym]; !ok {
			candidates[sym] = triggerDef{old.TriggerType, old.Trigger, old.TriggerAt}
		}
	}

	liq := deriveLiquidityStatesWithContext(quotes, bars, liquidityBaselines, now)
	updated := finalizeCompletedCatalysts(existing, now)
	activeCount := 0
	for sym, tr := range candidates {
		q := quotes[sym]
		old, exists := existing[sym]
		triggerPrice, triggerAt := q.Price, tr.at
		if exists && old.TriggerPrice > 0 {
			triggerPrice, triggerAt = old.TriggerPrice, old.TriggerAt
		}
		if q.Price <= 0 || triggerPrice <= 0 {
			continue
		}
		move := (q.Price/triggerPrice - 1) * 100
		phase := catalystPhase(now, triggerAt)
		state := "UNCONFIRMED"
		if move >= 1 {
			state = "POSITIVE"
		} else if move <= -1 {
			state = "NEGATIVE"
		}
		volatility := "NORMAL"
		if math.Abs(move) >= 8 {
			state = "EXTREME VOLATILITY"
			volatility = "EXTREME"
		} else if math.Abs(move) >= 4 {
			volatility = "EXPANDED"
		}
		l := liq[sym]
		flags := []string{}
		if currentLiquidityMarketRisk(l, now) {
			flags = append(flags, "LIQUIDITY RISK")
		}
		pm, _ := premarketSnapshotFromBars(sym, bars[sym]["intraday"], q, now)
		if extendedAtMarketOpen(q, bars[sym]["daily"], bars[sym]["intraday"], pm, now) {
			flags = append(flags, "EXTENDED")
		}
		gap := 0.0
		prev := q.PreviousClose
		if prev <= 0 {
			prev = q.PriorSessionClose
		}
		if q.Open > 0 && prev > 0 {
			gap = (q.Open/prev - 1) * 100
		}
		vw := vwapFromCurrentSession(bars[sym]["intraday"], now)
		vwapState := "VWAP PENDING"
		if vw > 0 {
			if q.Price >= vw {
				vwapState = "ABOVE VWAP"
			} else {
				vwapState = "BELOW VWAP"
			}
		}
		holdFade := "MIXED"
		if math.Abs(move) < .5 {
			holdFade = "FLAT"
		} else if move > 0 && (vw == 0 || q.Price >= vw) {
			holdFade = "HOLDING"
		} else if move < 0 && (vw == 0 || q.Price <= vw) {
			holdFade = "HOLDING"
		} else {
			holdFade = "FADING / REVERSING"
		}
		reactions := map[string]float64{}
		for label, d := range map[string]time.Duration{"5m": 5 * time.Minute, "15m": 15 * time.Minute, "30m": 30 * time.Minute, "60m": 60 * time.Minute} {
			if now.Sub(time.UnixMilli(triggerAt)) >= d {
				reactions[label] = reactionAt(history[sym], triggerAt+d.Milliseconds(), triggerPrice)
			}
		}
		reactions["session"] = move
		completedAt := int64(0)
		if phase == "COMPLETE" {
			completedAt = now.UnixMilli()
		} else {
			activeCount++
		}
		updated[sym] = CatalystReactionState{Symbol: sym, TriggerType: tr.typ, Trigger: tr.trigger, Session: session, State: state, Phase: phase, TriggerAt: triggerAt, TriggerPrice: triggerPrice, LatestPrice: q.Price, MovePercent: move, GapPercent: gap, Volume: q.Volume, RelativeVolume: relativeVolumeFromDaily(q, bars[sym]["daily"]), SpreadPct: l.SpreadPct, Liquidity: l.State, VWAPState: vwapState, OpeningRangeState: openingRangeState(q.Price, bars[sym]["intraday"], now), HoldFadeState: holdFade, VolatilityState: volatility, ReactionPercent: reactions, Flags: uniqueStrings(flags), Detail: "Event-confirmed reaction lifecycle · contextual evidence only; deterministic desk scores unchanged.", UpdatedAt: now.UnixMilli(), CompletedAt: completedAt}
	}

	state := "READY"
	detail := "Event-driven · no material trigger active"
	attention := "READY"
	if len(armed) > 0 && activeCount == 0 {
		state = "ARMED"
		attention = "READY WITH CAUTION"
		detail = fmt.Sprintf("%d scheduled earnings event(s) armed · waiting for sourced release confirmation", len(armed))
	} else if activeCount > 0 {

		state = "TRIGGERED"
		for _, r := range updated {
			if r.CompletedAt == 0 && math.Abs(r.MovePercent) >= 1 {
				state = "REACTION"
				break
			}
		}
		attention = "REVIEW REQUIRED"
		detail = fmt.Sprintf("Event-driven · %d confirmed material trigger(s) active", activeCount)
	} else if len(updated) > 0 {
		state = "COMPLETE"
		detail = "Latest confirmed catalyst reaction lifecycle completed"
	}
	prep.State = state
	prep.Attention = attention
	prep.LastAttempt = now.UnixMilli()
	prep.LastSuccess = prep.LastAttempt
	prep.TradingDay = day
	prep.Detail = detail
	prep.Summary = []string{fmt.Sprintf("%d active confirmed catalyst(s)", activeCount), fmt.Sprintf("%d scheduled event(s) armed", len(armed))}

	materialChange := catalystWatchMateriallyChanged(existing, updated, prepBefore, prep)
	e.mu.Lock()
	e.catalystReactions = updated
	e.preparations["catalyst-watch"] = prep
	e.health["catalyst-watch"] = prep.State + " · " + prep.Detail
	e.lastUpdated["catalyst-watch"] = now.UnixMilli()
	e.mu.Unlock()
	if materialChange {
		_ = e.saveCache()
	}
}
