package main

import (
	"context"
	"fmt"
	"math"
	"net/http"
	"net/url"
	"sort"
	"strings"
	"time"
)

const (
	opportunityProductionFloor = 78.0
	opportunityRetainFloor     = 70.0
	opportunityShadowFloor     = 72.0
	opportunityMaxPromotions   = 5
	opportunityPromotionTTL    = 4 * time.Minute
	opportunityUniverseTTL     = 12 * time.Hour
	opportunityRotateCount     = 160
)

// regularSessionProgress returns the fraction of the regular U.S. session that
// has elapsed. It is used only to normalize volume participation, never as a
// trading signal or deterministic score input.
func regularSessionProgress(now time.Time) float64 {
	et := now.In(easternLocation())
	open := time.Date(et.Year(), et.Month(), et.Day(), 9, 30, 0, 0, easternLocation())
	closeAt := time.Date(et.Year(), et.Month(), et.Day(), 16, 0, 0, 0, easternLocation())
	if !et.After(open) {
		return 0
	}
	if !et.Before(closeAt) {
		return 1
	}
	return clampScanner(et.Sub(open).Seconds()/closeAt.Sub(open).Seconds(), 0, 1)
}

func sessionRelativeVolumeFromSnapshot(snap alpacaLiveSnapshot, session string, now time.Time) float64 {
	cur, prev := snap.DailyBar.Volume, snap.PrevDailyBar.Volume
	if cur <= 0 || prev <= 0 {
		return 0
	}
	if session == "regular" {
		// A small floor prevents the first few minutes from generating absurd RVOL
		// values while still making early participation visible.
		fraction := math.Max(.08, regularSessionProgress(now))
		return cur / (prev * fraction)
	}
	return cur / prev
}

func rangeExpansionFromSnapshot(snap alpacaLiveSnapshot) float64 {
	prevClose := snap.PrevDailyBar.Close
	if prevClose <= 0 {
		return 0
	}
	currentRange := math.Max(0, snap.DailyBar.High-snap.DailyBar.Low) / prevClose * 100
	prevRange := math.Max(0, snap.PrevDailyBar.High-snap.PrevDailyBar.Low) / prevClose * 100
	if prevRange <= 0 {
		return 0
	}
	return currentRange / prevRange
}

func priceConfirmationFromSnapshot(change float64, snap alpacaLiveSnapshot) string {
	if math.Abs(change) < 1.25 || snap.DailyBar.High <= snap.DailyBar.Low {
		return "DEVELOPING"
	}
	price := snap.LatestTrade.Price
	if price <= 0 {
		price = snap.DailyBar.Close
	}
	mid := (snap.DailyBar.High + snap.DailyBar.Low) / 2
	if (change > 0 && price >= mid) || (change < 0 && price <= mid) {
		return "CONFIRMED"
	}
	return "MIXED"
}

// enrichOpportunityMetrics extends the existing Discovery score with a
// non-execution opportunity layer. It deliberately does not change protected
// Day/Swing/Long deterministic formulas.
func enrichOpportunityMetrics(base ScannerResult, snap alpacaLiveSnapshot, session string, now time.Time, persistent, catalyst bool) ScannerResult {
	x := base
	x.SessionRelativeVolume = sessionRelativeVolumeFromSnapshot(snap, session, now)
	x.RangeExpansion = rangeExpansionFromSnapshot(snap)
	x.PriceConfirmation = priceConfirmationFromSnapshot(x.ChangePercent, snap)

	rvol := x.SessionRelativeVolume
	if rvol <= 0 {
		rvol = x.RelativeVolume
	}
	x.UnusualVolumeScore = clampScanner((rvol-1)*48+math.Log10(math.Max(1, x.DollarVolume/1_000_000))*9, 0, 100)
	x.VolatilityScore = clampScanner(math.Abs(x.ChangePercent)*13+math.Abs(x.GapPercent)*7+math.Max(0, x.RangeExpansion-1)*38, 0, 100)
	liquidity := clampScanner(100-x.SpreadPercent*90, 0, 100)
	confirmation := 35.0
	if x.PriceConfirmation == "CONFIRMED" {
		confirmation = 100
	} else if x.PriceConfirmation == "MIXED" {
		confirmation = 55
	}
	x.OpportunityScore = .36*x.UnusualVolumeScore + .34*x.VolatilityScore + .20*liquidity + .10*confirmation
	if catalyst {
		x.OpportunityScore += 7
		x.Reasons = append(x.Reasons, "Material catalyst/news context")
	}
	if persistent {
		x.OpportunityScore += 5
		x.Reasons = append(x.Reasons, "Persistent across radar cycles")
	}
	x.OpportunityScore = clampScanner(x.OpportunityScore, 0, 100)
	if rvol >= 1.5 {
		x.Reasons = append(x.Reasons, fmt.Sprintf("%.2fx session-normalized volume", rvol))
	}
	if x.RangeExpansion >= 1.3 {
		x.Reasons = append(x.Reasons, fmt.Sprintf("%.2fx prior-day range expansion", x.RangeExpansion))
	}
	if math.Abs(x.ChangePercent) >= 2.5 {
		x.Reasons = append(x.Reasons, fmt.Sprintf("%.1f%% price dislocation", x.ChangePercent))
	}
	if x.PriceConfirmation == "CONFIRMED" {
		x.Reasons = append(x.Reasons, "Price confirms the active move")
	}
	x.Reasons = uniqueStrings(x.Reasons)
	return x
}

func opportunityPromotionEligible(x ScannerResult, floor float64) bool {
	if x.OpportunityScore < floor || x.Price < 2 || x.DollarVolume < 5_000_000 || x.SpreadPercent <= 0 || x.SpreadPercent > .75 {
		return false
	}
	rvol := x.SessionRelativeVolume
	if rvol <= 0 {
		rvol = x.RelativeVolume
	}
	return rvol >= 1.35 || x.RangeExpansion >= 1.25 || math.Abs(x.ChangePercent) >= 3
}

func previousPromotionMap(rows []OpportunityPromotion) map[string]OpportunityPromotion {
	out := map[string]OpportunityPromotion{}
	for _, p := range rows {
		if sym := normalizeSymbol(p.Symbol); sym != "" {
			out[sym] = p
		}
	}
	return out
}

func selectOpportunityPromotions(candidates []ScannerResult, previous []OpportunityPromotion, now time.Time) []OpportunityPromotion {
	prev := previousPromotionMap(previous)
	rows := append([]ScannerResult{}, candidates...)
	sort.SliceStable(rows, func(i, j int) bool { return rows[i].OpportunityScore > rows[j].OpportunityScore })
	out := make([]OpportunityPromotion, 0, opportunityMaxPromotions)
	nowms := now.UnixMilli()
	for _, x := range rows {
		old, wasPromoted := prev[x.Symbol]
		floor := opportunityProductionFloor
		if wasPromoted {
			floor = opportunityRetainFloor
		}
		if !opportunityPromotionEligible(x, floor) {
			continue
		}
		started := nowms
		if wasPromoted && old.PromotedAt > 0 {
			started = old.PromotedAt
		}
		out = append(out, OpportunityPromotion{
			Symbol:           x.Symbol,
			Score:            x.OpportunityScore,
			State:            "PROMOTED",
			Reasons:          append([]string{}, x.Reasons...),
			PromotedAt:       started,
			LastConfirmedAt:  nowms,
			ExpiresAt:        now.Add(opportunityPromotionTTL).UnixMilli(),
			ShadowWouldMatch: opportunityPromotionEligible(x, opportunityShadowFloor),
		})
		if len(out) >= opportunityMaxPromotions {
			break
		}
	}
	return out
}

func opportunityRadarCadence(session string, hasPromotion, degraded bool) time.Duration {
	var d time.Duration
	switch session {
	case "regular":
		if hasPromotion {
			d = 60 * time.Second
		} else {
			d = 2 * time.Minute
		}
	case "pre-market", "after-hours":
		d = 3 * time.Minute
	case "overnight":
		d = 5 * time.Minute
	default:
		d = 30 * time.Minute
	}
	if degraded {
		d *= 2
		if d > 30*time.Minute {
			d = 30 * time.Minute
		}
	}
	return d
}

func radarSessionActive(session string) bool {
	return session == "regular" || session == "pre-market" || session == "after-hours" || session == "overnight"
}

func (e *Engine) opportunityUniverse(ctx context.Context, key, secret string, now time.Time) []string {
	e.mu.RLock()
	cached := append([]string{}, e.radarUniverse...)
	cachedAt := e.radarUniverseAt
	e.mu.RUnlock()
	if len(cached) > 0 && now.UnixMilli()-cachedAt < int64(opportunityUniverseTTL/time.Millisecond) {
		return cached
	}
	rows := e.scannerUniverse(ctx, key, secret)
	if len(rows) == 0 {
		rows = append([]string{}, discoverySeedUniverse...)
	}
	e.mu.Lock()
	e.radarUniverse = append([]string{}, rows...)
	e.radarUniverseAt = now.UnixMilli()
	e.mu.Unlock()
	return rows
}

func radarSampleUniverse(universe []string, cursor int) ([]string, int) {
	seeds := uniqueSymbols(discoverySeedUniverse)
	if len(universe) == 0 {
		return seeds, 0
	}
	out := append([]string{}, seeds...)
	count := minInt(opportunityRotateCount, len(universe))
	for i := 0; i < count; i++ {
		out = append(out, universe[(cursor+i)%len(universe)])
	}
	return uniqueSymbols(out), (cursor + count) % len(universe)
}

func (e *Engine) fetchOpportunitySnapshots(ctx context.Context, key, secret, session string, universe []string) (map[string]alpacaLiveSnapshot, int, error) {
	feed := "iex"
	if session == "overnight" {
		feed = "overnight"
	}
	client := &http.Client{Timeout: 18 * time.Second}
	all := map[string]alpacaLiveSnapshot{}
	scanned := 0
	var lastErr error
	for start := 0; start < len(universe); start += 50 {
		end := minInt(start+50, len(universe))
		batch := universe[start:end]
		raw := "https://data.alpaca.markets/v2/stocks/snapshots?symbols=" + url.QueryEscape(strings.Join(batch, ",")) + "&feed=" + url.QueryEscape(feed)
		var payload map[string]alpacaLiveSnapshot
		if err := e.providerGetJSONTier(ctx, "Alpaca", WorkTierRadarPromoted, client, raw, map[string]string{"APCA-API-KEY-ID": key, "APCA-API-SECRET-KEY": secret}, &payload); err != nil {
			lastErr = err
			continue
		}
		scanned += len(batch)
		for sym, snap := range payload {
			all[normalizeSymbol(sym)] = snap
		}
	}
	if len(all) == 0 && lastErr != nil {
		return nil, scanned, lastErr
	}
	return all, scanned, nil
}

func (e *Engine) radarCatalystSymbols(now time.Time) map[string]bool {
	out := map[string]bool{}
	cutoff := now.Add(-2 * time.Hour).Unix()
	e.mu.RLock()
	news := clone(e.news)
	filings := clone(e.filings)
	for sym, c := range e.catalystReactions {
		if c.TriggerAt > 0 && (c.CompletedAt == 0 || c.CompletedAt >= now.Add(-2*time.Hour).UnixMilli()) {
			out[normalizeSymbol(sym)] = true
		}
	}
	for _, n := range news {
		if n.Datetime < cutoff {
			continue
		}
		for _, sym := range n.Symbols {
			if s := normalizeSymbol(sym); s != "" {
				out[s] = true
			}
		}
	}
	e.mu.RUnlock()

	// Community evidence is never authoritative, but fused/corroborated recent
	// evidence can increase Opportunity Radar context. It still has no direct
	// protected deterministic Action/Score/Readiness impact.
	e.app.mu.RLock()
	communityItems := clone(e.app.state.CommunityEvidence)
	e.app.mu.RUnlock()
	community := buildCommunityEvidenceFusion(communityItems, news, filings, now)
	recentCommunity := now.Add(-2 * time.Hour).UnixMilli()
	for _, c := range community.Clusters {
		if c.Symbol == "" || c.LatestAt < recentCommunity {
			continue
		}
		if c.Materiality == "HIGH" || c.Materiality == "ELEVATED" {
			out[normalizeSymbol(c.Symbol)] = true
		}
	}
	return out
}

func (e *Engine) runOpportunityRadarCycle(ctx context.Context, key, secret string, now time.Time) {
	session := marketSessionET(now)
	e.mu.RLock()
	prevRadar := clone(e.scanner.Radar)
	providerHealth := strings.ToLower(e.health["alpaca-live"] + " " + e.health["alpaca-stream"])
	cursor := e.radarCursor
	e.mu.RUnlock()
	degraded := strings.Contains(providerHealth, "error") || strings.Contains(providerHealth, "failed") || strings.Contains(providerHealth, "reconnecting")
	cadence := opportunityRadarCadence(session, len(prevRadar.Promotions) > 0, degraded)
	if !radarSessionActive(session) {
		state := OpportunityRadarState{Status: "IDLE", Session: session, Message: "U.S. market session inactive; broad opportunity scanning is paused to avoid wasteful provider work.", CadenceMs: int64(cadence / time.Millisecond), LastRun: prevRadar.LastRun, NextRun: now.Add(cadence).UnixMilli(), Provider: "Alpaca snapshots", ProductionFloor: opportunityProductionFloor, ShadowFloor: opportunityShadowFloor, ShadowOnly: true}
		e.mu.Lock()
		e.scanner.Radar = state
		e.health["scanner"] = "Opportunity Radar idle · market session inactive"
		e.mu.Unlock()
		return
	}
	started := time.Now()
	universe := e.opportunityUniverse(ctx, key, secret, now)
	sample, nextCursor := radarSampleUniverse(universe, cursor)
	// Reuse entitled market-activity data as low-cost dynamic seeds instead of
	// fetching it only to render a card. The normal U.S.-symbol boundary still applies.
	e.mu.RLock()
	activity := clone(e.marketActivity)
	e.mu.RUnlock()
	for _, row := range append(append(append([]MarketMover{}, activity.MostActive...), activity.Gainers...), activity.Losers...) {
		if sym, ok := parseUserTicker(row.Symbol); ok {
			sample = append(sample, sym)
		}
	}
	sample = uniqueSymbols(sample)
	payload, scanned, err := e.fetchOpportunitySnapshots(ctx, key, secret, session, sample)
	if err != nil {
		state := prevRadar
		state.Status = "DEGRADED"
		state.Session = session
		state.Message = "Opportunity Radar snapshot refresh degraded: " + err.Error()
		state.CadenceMs = int64(cadence / time.Millisecond)
		state.NextRun = now.Add(cadence).UnixMilli()
		state.Provider = "Alpaca snapshots"
		e.mu.Lock()
		e.scanner.Radar = state
		e.health["scanner"] = "Opportunity Radar degraded · snapshot provider"
		e.mu.Unlock()
		return
	}
	previous := map[string]ScannerResult{}
	for _, x := range prevRadar.Candidates {
		previous[x.Symbol] = x
	}
	catalysts := e.radarCatalystSymbols(now)
	candidates := make([]ScannerResult, 0, len(payload))
	for symbol, snap := range payload {
		base := scannerScoreFromSnapshot(symbol, "day", snap)
		if base.Price < 2 || base.SpreadPercent <= 0 || base.SpreadPercent > 2.5 || base.DollarVolume < 1_000_000 {
			continue
		}
		prior, persisted := previous[symbol]
		persisted = persisted && prior.OpportunityScore >= 60
		x := enrichOpportunityMetrics(base, snap, session, now, persisted, catalysts[symbol])
		if x.OpportunityScore >= 48 || x.SessionRelativeVolume >= 1.25 || math.Abs(x.ChangePercent) >= 2 {
			candidates = append(candidates, x)
		}
	}
	sort.SliceStable(candidates, func(i, j int) bool { return candidates[i].OpportunityScore > candidates[j].OpportunityScore })
	if len(candidates) > 40 {
		candidates = candidates[:40]
	}
	promotions := selectOpportunityPromotions(candidates, prevRadar.Promotions, now)
	prevPromoted := previousPromotionMap(prevRadar.Promotions)
	newlyPromoted := []string{}
	for _, p := range promotions {
		e.mu.Lock()
		e.livePriorityHints[p.Symbol] = now.UnixMilli()
		e.mu.Unlock()
		if _, existed := prevPromoted[p.Symbol]; !existed {
			newlyPromoted = append(newlyPromoted, p.Symbol)
		}
	}
	if len(newlyPromoted) > 0 {
		// Material promotion is the trigger for immediate targeted history hydration;
		// normal symbols keep the existing session cadence.
		e.requestHistoryHydration(newlyPromoted...)
	}
	cadence = opportunityRadarCadence(session, len(promotions) > 0, degraded)
	state := OpportunityRadarState{
		Status:          "ACTIVE",
		Session:         session,
		Message:         fmt.Sprintf("%d opportunities ranked from %d low-cost snapshot observations; %d promoted to bounded live priority.", len(candidates), scanned, len(promotions)),
		Candidates:      candidates,
		Promotions:      promotions,
		Scanned:         scanned,
		DurationMs:      time.Since(started).Milliseconds(),
		CadenceMs:       int64(cadence / time.Millisecond),
		LastRun:         now.UnixMilli(),
		NextRun:         now.Add(cadence).UnixMilli(),
		Provider:        "Alpaca snapshots",
		ProductionFloor: opportunityProductionFloor,
		ShadowFloor:     opportunityShadowFloor,
		ShadowOnly:      true,
	}
	e.mu.Lock()
	e.radarCursor = nextCursor
	e.scanner.Radar = state
	e.health["scanner"] = fmt.Sprintf("Opportunity Radar active · %d promoted", len(promotions))
	e.lastUpdated["scanner-radar"] = now.UnixMilli()
	ws, alpacaWS := e.ws, e.alpacaWS
	e.mu.Unlock()
	// Promotions reuse the existing live-allocation owner and reserve slots. No
	// second live engine or subscription owner is introduced.
	if ws != nil {
		e.syncLiveSubscriptions(ws)
	}
	if alpacaWS != nil {
		e.syncAlpacaSubscriptions(alpacaWS)
	}
	e.app.hub.Broadcast(map[string]any{"type": "scanner", "scanner": e.Snapshot().Scanner})
}

func (e *Engine) opportunityRadarLoop(ctx context.Context, key, secret string) {
	for {
		now := time.Now()
		if e.shouldShedBackground() {
			e.mu.Lock()
			state := e.scanner.Radar
			state.Status = "DEFERRED"
			state.Message = "Broad Opportunity Radar work deferred by v17 load shedding; selected/watchlist freshness is protected first."
			state.NextRun = now.Add(time.Minute).UnixMilli()
			e.scanner.Radar = state
			e.health["scanner"] = "deferred · LOCAL LOAD · critical symbol freshness protected"
			e.mu.Unlock()
		} else if release, ok := e.workload.TryAcquireTier("scanner", WorkTierRadarPromoted); ok {
			e.runOpportunityRadarCycle(ctx, key, secret, now)
			release()
		}
		e.mu.RLock()
		cadenceMs := e.scanner.Radar.CadenceMs
		e.mu.RUnlock()
		wait := time.Duration(cadenceMs) * time.Millisecond
		if wait <= 0 {
			wait = 2 * time.Minute
		}
		wait = capWaitToSessionBoundary(time.Now(), wait)
		t := time.NewTimer(wait)
		select {
		case <-ctx.Done():
			t.Stop()
			return
		case <-t.C:
		}
	}
}
