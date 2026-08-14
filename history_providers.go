package main

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"sort"
	"strings"
	"time"
)

type alpacaHistorySpec struct {
	Name, Timeframe string
	Start           time.Time
	Symbols         []string
}

func historySpecsForStateMode(st AppState, onlySymbols []string, mode string) []alpacaHistorySpec {
	daySet := symbolSetForWatchlist(st, st.Settings.DayWatchlistID, st.Settings.DayEnabled)
	swingSet := symbolSetForWatchlist(st, st.Settings.SwingWatchlistID, st.Settings.SwingEnabled)
	longSet := symbolSetForWatchlist(st, st.Settings.LongWatchlistID, st.Settings.LongEnabled)

	focus := normalizeSymbol(st.UI.SelectedTicker)
	if focus != "" && focus != "VIX" {
		if st.Settings.DayEnabled {
			daySet[focus] = true
		}
		if st.Settings.SwingEnabled {
			swingSet[focus] = true
		}
		if st.Settings.LongEnabled {
			longSet[focus] = true
		}
	}

	if st.Settings.DayEnabled {
		for _, symbol := range []string{"SPY", "QQQ", "IWM", "XLK"} {
			daySet[symbol] = true
		}
	}
	if st.Settings.SwingEnabled {
		for _, symbol := range []string{"SPY", "QQQ", "IWM", "XLK"} {
			swingSet[symbol] = true
		}
	}
	if st.Settings.LongEnabled {
		for _, symbol := range []string{"SPY", "QQQ", "IWM", "TLT"} {
			longSet[symbol] = true
		}
	}

	filter := map[string]bool{}
	for _, symbol := range onlySymbols {
		symbol = normalizeSymbol(symbol)
		if symbol != "" {
			filter[symbol] = true
		}
	}
	setList := func(set map[string]bool) []string {
		out := make([]string, 0, len(set))
		for symbol := range set {
			if symbol == "VIX" || symbol == "" {
				continue
			}
			if len(filter) > 0 && !filter[symbol] {
				continue
			}
			out = append(out, symbol)
		}
		sort.Strings(out)
		return out
	}
	dailySet := map[string]bool{}
	for symbol := range swingSet {
		dailySet[symbol] = true
	}
	for symbol := range longSet {
		dailySet[symbol] = true
	}
	mode = strings.ToLower(strings.TrimSpace(mode))
	if mode == "" {
		mode = "all"
	}
	var specs []alpacaHistorySpec
	if mode == "all" || mode == "intraday" {
		if syms := setList(daySet); len(syms) > 0 {
			specs = append(specs, alpacaHistorySpec{"intraday", "15Min", time.Now().AddDate(0, -1, 0), syms})
		}
	}
	if mode == "all" || mode == "daily" {
		if syms := setList(dailySet); len(syms) > 0 {
			specs = append(specs, alpacaHistorySpec{"daily", "1Day", time.Now().AddDate(-2, 0, 0), syms})
		}
		if syms := setList(longSet); len(syms) > 0 {
			specs = append(specs, alpacaHistorySpec{"weekly", "1Week", time.Now().AddDate(-6, 0, 0), syms})
		}
	}
	return specs
}

func historySpecsForState(st AppState, onlySymbols []string) []alpacaHistorySpec {
	return historySpecsForStateMode(st, onlySymbols, "all")
}

// requestHistoryHydration triggers a one-shot bar refresh without changing the
// normal 15-minute cadence. Passing no symbols hydrates the complete active
// universe (used after cache clear or engine-enable); passing symbols limits the
// request to those newly required by a watchlist change.
func (e *Engine) requestHistoryHydration(symbols ...string) {
	e.mu.RLock()
	status, mode := e.status, e.mode
	e.mu.RUnlock()
	if status != "running" && status != "degraded" {
		return
	}

	if mode == "demo" {
		requested := uniqueSymbols(symbols)
		if len(requested) == 0 {
			e.app.mu.RLock()
			st := clone(e.app.state)
			e.app.mu.RUnlock()
			for _, sp := range historySpecsForState(st, nil) {
				requested = append(requested, sp.Symbols...)
			}
			requested = uniqueSymbols(requested)
		}
		for _, symbol := range requested {
			e.ensureDemoSymbol(symbol)
		}
		e.mu.Lock()
		e.lastUpdated["history"] = time.Now().UnixMilli()
		e.health["history"] = fmt.Sprintf("demo · %d symbols rehydrated", len(requested))
		e.mu.Unlock()
		return
	}
	if mode != "live" {
		return
	}
	e.app.mu.RLock()
	key := strings.TrimSpace(e.app.secrets.AlpacaKey)
	secret := strings.TrimSpace(e.app.secrets.AlpacaSecret)
	e.app.mu.RUnlock()
	if key == "" || secret == "" {
		return
	}
	requested := append([]string(nil), symbols...)
	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), 75*time.Second)
		defer cancel()
		e.refreshAlpacaHistoryScoped(ctx, key, secret, requested)
	}()
}

func (e *Engine) refreshAlpacaHistoryScopedMode(ctx context.Context, key, secret string, onlySymbols []string, mode string) int {
	e.historyRefreshMu.Lock()
	defer e.historyRefreshMu.Unlock()
	e.setHealth("history", "loading")
	e.app.mu.RLock()
	st := clone(e.app.state)
	e.app.mu.RUnlock()
	specs := historySpecsForStateMode(st, onlySymbols, mode)
	if len(specs) == 0 {
		if len(onlySymbols) == 0 {
			e.setHealth("history", "idle · no enabled desk symbols")
		}
		return 0
	}
	client := &http.Client{Timeout: 40 * time.Second}
	loaded := 0
	intradayLoaded, dailyLoaded := false, false
	for _, sp := range specs {
		type historyBatch struct {
			symbols []string
			start   time.Time
		}
		groups := []historyBatch{{symbols: sp.Symbols, start: sp.Start}}
		if sp.Name == "daily" {
			regular, benchmarks := []string{}, []string{}
			for _, sym := range sp.Symbols {
				if sym == "SPY" || sym == "QQQ" {
					benchmarks = append(benchmarks, sym)
				} else {
					regular = append(regular, sym)
				}
			}
			groups = nil
			if len(regular) > 0 {
				groups = append(groups, historyBatch{symbols: regular, start: sp.Start})
			}
			if len(benchmarks) > 0 {
				groups = append(groups, historyBatch{symbols: benchmarks, start: time.Now().AddDate(-10, 0, 0)})
			}
		}
		for _, group := range groups {
			for start := 0; start < len(group.symbols); start += 50 {
				end := minInt(start+50, len(group.symbols))
				batch := group.symbols[start:end]
				raw := strings.TrimRight(alpacaDataBaseURL, "/") + "/v2/stocks/bars?symbols=" + url.QueryEscape(strings.Join(batch, ",")) + "&timeframe=" + url.QueryEscape(sp.Timeframe) + "&start=" + url.QueryEscape(group.start.UTC().Format(time.RFC3339)) + "&limit=10000&adjustment=all&feed=iex&sort=asc"
				var payload struct {
					Bars map[string][]struct {
						C float64 `json:"c"`
						H float64 `json:"h"`
						L float64 `json:"l"`
						O float64 `json:"o"`
						V float64 `json:"v"`
						T string  `json:"t"`
					} `json:"bars"`
				}
				if err := e.providerGetJSONTier(ctx, "Alpaca", e.workTierForSymbols(batch), client, raw, map[string]string{"APCA-API-KEY-ID": key, "APCA-API-SECRET-KEY": secret}, &payload); err != nil {
					e.recordProviderFailure("Alpaca", err)
					e.setError("Alpaca history", err)
					e.setHealth("history", "degraded")
					continue
				}
				type sessionCloseUpdate struct {
					symbol string
					close  float64
					at     int64
					prior  float64
				}
				closeUpdates := make([]sessionCloseUpdate, 0, len(payload.Bars))
				e.mu.Lock()
				for sym, rows := range payload.Bars {
					sym = normalizeSymbol(sym)
					if e.bars[sym] == nil {
						e.bars[sym] = map[string][]Bar{}
					}
					out := make([]Bar, 0, len(rows))
					for _, r := range rows {
						t, err := time.Parse(time.RFC3339Nano, r.T)
						if err != nil {
							continue
						}
						out = append(out, Bar{T: t.Unix(), O: r.O, H: r.H, L: r.L, C: r.C, V: r.V})
					}
					if len(out) > 0 {
						e.bars[sym][sp.Name] = out
						loaded += len(out)
						if sp.Name == "intraday" {
							intradayLoaded = true
						}
						if sp.Name == "daily" || sp.Name == "weekly" {
							dailyLoaded = true
						}
					}
					if sp.Name == "daily" && len(out) > 0 {
						prior := 0.0
						if len(out) > 1 {
							prior = out[len(out)-2].C
						}
						closeUpdates = append(closeUpdates, sessionCloseUpdate{symbol: sym, close: out[len(out)-1].C, at: out[len(out)-1].T * 1000, prior: prior})
					}
				}
				e.mu.Unlock()
				for _, update := range closeUpdates {
					e.updateCanonicalSessionClose(update.symbol, update.close, update.at, update.prior)
				}
			}
		}
	}
	if loaded > 0 {
		e.mu.Lock()
		nowMs := time.Now().UnixMilli()
		e.lastUpdated["history"] = nowMs
		if intradayLoaded {
			e.lastUpdated["history-intraday"] = nowMs
		}
		if dailyLoaded {
			e.lastUpdated["history-daily"] = nowMs
		}
		e.mu.Unlock()
		e.recordProviderSuccess("Alpaca")
		e.setHealth("history", fmt.Sprintf("healthy · Alpaca · %d bars · desk-scoped", loaded))
	} else {
		e.setHealth("history", "degraded · Alpaca primary returned no bars")
	}
	return loaded
}

func (e *Engine) refreshAlpacaHistoryScoped(ctx context.Context, key, secret string, onlySymbols []string) int {
	return e.refreshAlpacaHistoryScopedMode(ctx, key, secret, onlySymbols, "all")
}
