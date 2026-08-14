package main

import "time"

func symbolInWatchlist(st AppState, id, symbol string) bool {
	wl, ok := watchlistValueByID(st.Watchlists, id)
	if !ok {
		return false
	}
	for _, s := range wl.Symbols {
		if normalizeSymbol(s) == symbol {
			return true
		}
	}
	return false
}

func symbolInPromotion(rows []OpportunityPromotion, symbol string, now int64) bool {
	for _, p := range rows {
		if normalizeSymbol(p.Symbol) != symbol {
			continue
		}
		if p.ExpiresAt == 0 || p.ExpiresAt >= now {
			return true
		}
	}
	return false
}

// workTierForSymbol is the canonical runtime-priority owner for v17. It is
// intentionally independent from score/action formulas and only controls how
// scarce runtime/provider capacity is allocated.
func (e *Engine) workTierForSymbol(symbol string) WorkTier {
	symbol = normalizeSymbol(symbol)
	if symbol == "" {
		return WorkTierBackground
	}
	if symbol == "SPY" || symbol == "QQQ" {
		return WorkTierMarketCritical
	}

	var st AppState
	if e != nil && e.app != nil {
		e.app.mu.RLock()
		st = clone(e.app.state)
		e.app.mu.RUnlock()
	}
	if normalizeSymbol(st.UI.SelectedTicker) == symbol {
		return WorkTierMarketCritical
	}

	now := time.Now().UnixMilli()
	if e != nil {
		e.mu.RLock()
		_, catalyst := e.catalystReactions[symbol]
		_, openRisk := e.marketOpenFlags[symbol]
		hintAt := e.livePriorityHints[symbol]
		promoted := symbolInPromotion(e.scanner.Radar.Promotions, symbol, now)
		e.mu.RUnlock()
		if catalyst || openRisk {
			return WorkTierMarketCritical
		}
		// Decision Queue / user-focus promotion hints are actionable, not broad scan.
		if hintAt > 0 && now-hintAt <= livePriorityHintTTL.Milliseconds() {
			return WorkTierUserActionable
		}
		if promoted {
			return WorkTierRadarPromoted
		}
	}

	for _, s := range pinnedTradableLiveSymbols {
		if s == symbol {
			return WorkTierUserActionable
		}
	}
	if symbolInWatchlist(st, st.Settings.DayWatchlistID, symbol) ||
		symbolInWatchlist(st, st.Settings.SwingWatchlistID, symbol) ||
		symbolInWatchlist(st, st.Settings.LongWatchlistID, symbol) {
		return WorkTierUserActionable
	}
	if symbolInWatchlist(st, st.Settings.DiscoveryWatchlistID, symbol) {
		return WorkTierRadarPromoted
	}
	for _, s := range masterSymbolsFromState(st) {
		if normalizeSymbol(s) == symbol {
			return WorkTierBroadDiscovery
		}
	}
	return WorkTierBackground
}

func (e *Engine) workTierForSymbols(symbols []string) WorkTier {
	best := WorkTierBackground
	for _, symbol := range symbols {
		tier := e.workTierForSymbol(symbol)
		if tier < best {
			best = tier
			if best == WorkTierMarketCritical {
				break
			}
		}
	}
	return best
}

func workTierLabel(t WorkTier) string {
	switch normalizeWorkTier(t) {
	case WorkTierMarketCritical:
		return "TIER 0 · MARKET CRITICAL"
	case WorkTierUserActionable:
		return "TIER 1 · USER ACTIONABLE"
	case WorkTierRadarPromoted:
		return "TIER 2 · RADAR PROMOTED"
	case WorkTierBroadDiscovery:
		return "TIER 3 · BROAD DISCOVERY"
	default:
		return "TIER 4 · BACKGROUND"
	}
}
