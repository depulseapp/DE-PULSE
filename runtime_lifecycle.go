package main

import (
	"context"
	"errors"
	"math"

	mathrand "math/rand/v2"
	"strings"
	"time"
)

func liveEquityProviderConfigured(s Secrets) bool {
	if strings.TrimSpace(s.Finnhub) != "" {
		return true
	}
	if strings.TrimSpace(s.AlpacaKey) != "" && strings.TrimSpace(s.AlpacaSecret) != "" {
		return true
	}
	return strings.TrimSpace(s.TwelveData) != ""
}

var runtimeStopTimeout = 3 * time.Second

func (e *Engine) Start() error {
	e.mu.Lock()
	if e.status == "running" || e.status == "starting" {
		e.mu.Unlock()
		return nil
	}
	if e.status == "stopping" {
		e.mu.Unlock()
		return errors.New("market runtime is stopping; retry after it has stopped")
	}
	ctx, cancel := context.WithCancel(context.Background())
	e.cancel = cancel
	e.status = "starting"
	e.message = "Starting market runtime…"
	e.startedAt = time.Now().UTC().Format(time.RFC3339)
	e.lastError = ""
	e.mu.Unlock()

	e.app.mu.RLock()
	mode := e.app.state.Settings.DataMode
	key := e.app.secrets.Finnhub
	secrets := clone(e.app.secrets)
	e.app.mu.RUnlock()
	if mode == "live" && !liveEquityProviderConfigured(secrets) {
		cancel()
		e.mu.Lock()
		if e.cancel != nil {
			e.cancel = nil
		}
		e.status = "stopped"
		e.message = "Add Alpaca, Finnhub, or Twelve Data market credentials in Settings before starting Live mode."
		e.lastError = ""
		e.mu.Unlock()
		return errors.New("add Alpaca, Finnhub, or Twelve Data market credentials in Settings before starting Live mode")
	}
	if mode != "live" {
		mode = "demo"
	}
	e.mu.Lock()
	e.mode = mode
	e.mu.Unlock()
	e.wg.Add(1)
	go func() { defer e.wg.Done(); e.runtimeLoadProfileLoop(ctx) }()
	if mode == "demo" {
		e.startDemo(ctx)
	} else {
		e.startLive(ctx, key)
	}
	return nil
}

func (e *Engine) finishStop() {
	_ = e.saveCache()
	e.mu.Lock()

	e.status = "stopped"
	e.message = "Runtime is stopped"
	e.cancel = nil
	e.ws = nil
	e.alpacaWS = nil
	e.webSocketConnected = false
	e.subscribedSymbols = map[string]bool{}
	e.alpacaWebSocketConnected = false
	e.alpacaSubscribedSymbols = map[string]bool{}
	e.health = map[string]string{"quotes": "stopped", "history": "stopped", "fundamentals": "stopped", "vix": "stopped", "news": "stopped", "earnings": "stopped", "filings": "stopped", "global": "stopped", "macro-rates": "stopped", "fred-rates": "stopped", "macro-events": "stopped", "options": "stopped", "signal-validation": "ready", "scanner": "ready", "cache-refresh": "ready", "ai": "ready"}
	e.mu.Unlock()
	e.app.broadcastRuntime()
}

func (e *Engine) Stop() {
	e.mu.Lock()
	if e.status == "stopped" || e.status == "stopping" {
		e.mu.Unlock()
		return
	}
	e.status = "stopping"
	e.message = "Stopping market runtime…"
	cancel := e.cancel
	ws := e.ws
	alpacaWS := e.alpacaWS
	e.mu.Unlock()

	_ = e.saveCache()
	e.app.broadcastRuntime()
	if cancel != nil {
		cancel()
	}
	if ws != nil {
		_ = ws.Close()
	}
	if alpacaWS != nil {
		_ = alpacaWS.Close()
	}
	done := make(chan struct{})
	go func() { e.wg.Wait(); close(done) }()
	select {
	case <-done:
		e.finishStop()
	case <-time.After(runtimeStopTimeout):

		e.mu.Lock()
		e.message = "Runtime is still stopping…"
		e.mu.Unlock()
		e.app.broadcastRuntime()
		go func() {
			<-done
			e.finishStop()
		}()
	}
}

func (e *Engine) startDemo(ctx context.Context) {
	e.seedDemo()
	e.mu.Lock()
	e.status = "running"
	e.message = "Demo stream running"
	e.mu.Unlock()
	e.app.broadcastRuntime()
	e.wg.Add(1)
	go func() {
		defer e.wg.Done()
		ticker := time.NewTicker(time.Second)
		defer ticker.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				for _, symbol := range e.trackedSymbols() {
					e.mu.RLock()
					q, ok := e.quotes[symbol]
					e.mu.RUnlock()
					if !ok || q.Price <= 0 {
						continue
					}
					vol := 0.0008
					if symbol == "TSLA" || symbol == "PLTR" || symbol == "NVDA" {
						vol = 0.0017
					}
					next := math.Max(.01, q.Price+(mathrand.Float64()-.495)*q.Price*vol)
					e.updateQuote(symbol, Quote{Price: next, High: math.Max(q.High, next), Low: minPositive(q.Low, next), ProviderTimestamp: time.Now().UnixMilli()}, "demo")
				}
			}
		}
	}()
}

func symbolSetForWatchlist(st AppState, id string, enabled bool) map[string]bool {
	out := map[string]bool{}
	if !enabled {
		return out
	}
	if wl, ok := watchlistValueByID(st.Watchlists, id); ok {
		for _, symbol := range wl.Symbols {
			out[normalizeSymbol(symbol)] = true
		}
	}
	return out
}

func (e *Engine) seedDemo() {
	seeds := map[string]float64{"SPY": 608.42, "QQQ": 541.87, "DIA": 447.51, "IWM": 227.84, "VIX": 18.42, "GLD": 302.14, "SLV": 36.48, "TLT": 88.46, "USO": 78.22, "XLK": 266.19, "EWY": 74.80, "EWT": 61.40, "EWJ": 70.60, "EWH": 18.30, "MCHI": 57.20, "VGK": 74.10, "FEZ": 61.30, "UUP": 28.40, "HYG": 80.60, "LQD": 109.20, "SMH": 302.50, "NVDA": 193.62, "META": 784.18, "ORCL": 291.07, "PLTR": 178.55, "TSLA": 336.91, "SOFI": 24.82}
	e.app.mu.RLock()
	st := clone(e.app.state)
	e.app.mu.RUnlock()
	daySet := symbolSetForWatchlist(st, st.Settings.DayWatchlistID, st.Settings.DayEnabled)
	swingSet := symbolSetForWatchlist(st, st.Settings.SwingWatchlistID, st.Settings.SwingEnabled)
	longSet := symbolSetForWatchlist(st, st.Settings.LongWatchlistID, st.Settings.LongEnabled)
	discoverySet := symbolSetForWatchlist(st, st.Settings.DiscoveryWatchlistID, true)
	if st.Settings.DayEnabled {
		daySet["SPY"] = true
	}
	if st.Settings.SwingEnabled {
		swingSet["SPY"] = true
	}
	if st.Settings.LongEnabled {
		longSet["SPY"] = true
	}
	for _, symbol := range e.trackedSymbols() {
		price := seeds[symbol]
		if price == 0 {
			price = 50 + mathrand.Float64()*250
		}
		pc := price * (1 + (mathrand.Float64()-.5)*.025)
		e.updateQuote(symbol, Quote{Price: price, PreviousClose: pc, SessionClose: pc, PriorSessionClose: pc * (.99 + mathrand.Float64()*.02), Open: pc, High: math.Max(price, pc) * 1.006, Low: math.Min(price, pc) * .994, ProviderTimestamp: time.Now().UnixMilli()}, "demo")
		needIntraday := daySet[symbol]
		needDaily := swingSet[symbol] || longSet[symbol] || symbol == "VIX"
		needWeekly := longSet[symbol]
		needFundamentals := daySet[symbol] || swingSet[symbol] || longSet[symbol] || discoverySet[symbol]
		e.seedDemoBars(symbol, price, needIntraday, needDaily, swingSet[symbol] || longSet[symbol], needWeekly, needFundamentals)
	}
	e.mu.Lock()
	e.news = []NewsItem{
		{ID: "demo-1", Datetime: time.Now().Add(-10 * time.Minute).Unix(), Headline: "Demo mode is active — connect Finnhub in Settings for live market data", Summary: "The interface is using simulated prices. Your watchlists and settings are still saved locally.", Source: appName, Scope: "general", Symbols: e.trackedSymbols()},
		{ID: "demo-2", Datetime: time.Now().Add(-time.Hour).Unix(), Headline: "Watchlist briefing: technology names are showing elevated simulated activity", Summary: "This sample item demonstrates the watchlist-specific news and AI context experience.", Source: "Demo feed", Symbols: []string{"NVDA", "META", "ORCL", "PLTR"}},
	}
	tomorrow := time.Now().Add(24 * time.Hour)
	next := time.Now().Add(6 * 24 * time.Hour)
	eps1 := 1.24
	rev1 := 47200000000.0
	eps2 := 1.78
	rev2 := 16800000000.0
	e.earnings = []EarningsItem{{Symbol: "NVDA", Date: tomorrow.Format("2006-01-02"), Hour: "amc", EPSEstimate: &eps1, RevenueEstimate: &rev1, Quarter: 2, Year: tomorrow.Year()}, {Symbol: "ORCL", Date: next.Format("2006-01-02"), Hour: "amc", EPSEstimate: &eps2, RevenueEstimate: &rev2, Quarter: 1, Year: next.Year()}}
	e.filings = []FilingItem{{ID: "demo-filing", Symbol: "META", Company: "Meta Platforms, Inc.", Form: "8-K", FiledAt: time.Now().Add(-48 * time.Hour).Format("2006-01-02"), Description: "Sample current report shown in Demo mode"}}
	e.macroMetrics = map[string]MacroMetric{"DGS2": {Key: "DGS2", Label: "U.S. 2Y", Value: 4.02, Unit: "%", Source: "Demo", Provenance: "DEMO", UpdatedAt: time.Now().UnixMilli(), Status: "DEMO"}, "DGS10": {Key: "DGS10", Label: "U.S. 10Y", Value: 4.31, Unit: "%", Source: "Demo", Provenance: "DEMO", UpdatedAt: time.Now().UnixMilli(), Status: "DEMO"}, "DGS30": {Key: "DGS30", Label: "U.S. 30Y", Value: 4.88, Unit: "%", Source: "Demo", Provenance: "DEMO", UpdatedAt: time.Now().UnixMilli(), Status: "DEMO"}}
	demoEventAt := time.Now().Add(12 * time.Minute)
	e.macroEvents = []MacroEvent{{ID: "demo-high-impact", Region: "US", Name: "Demo High-Impact Macro Event", Impact: "HIGH", Lifecycle: "UPCOMING", StartsAt: demoEventAt.UnixMilli(), Date: demoEventAt.Format("2006-01-02"), TimeKnown: true, Source: "Demo", UpdatedAt: time.Now().UnixMilli()}}
	for _, sym := range []string{"NVDA", "TSLA", "PLTR", "META", "ORCL", "SPY", "QQQ"} {
		p := seeds[sym]
		e.options[sym] = OptionsContext{Symbol: sym, Provider: "Demo", Feed: "DEMO", State: "DEMO", Bias: map[string]string{"NVDA": "BULLISH", "TSLA": "CONFLICTING", "PLTR": "BULLISH"}[sym], CallContracts: 62, PutContracts: 47, CallVolume: 15400, PutVolume: 9800, PutCallVolume: .64, AverageIV: .42, ExpectedMove: p * .035, NearestExpiration: time.Now().AddDate(0, 0, 5).Format("2006-01-02"), UpdatedAt: time.Now().UnixMilli(), Provenance: "DEMO ONLY"}
		if e.options[sym].Bias == "" {
			o := e.options[sym]
			o.Bias = "NEUTRAL"
			e.options[sym] = o
		}
	}
	e.health = map[string]string{"quotes": "demo", "vix": "demo", "history": "demo", "news": "demo", "earnings": "demo", "filings": "demo", "fundamentals": "demo", "global": "demo", "macro-rates": "demo", "fred-rates": "demo", "macro-events": "demo", "options": "demo", "signal-validation": "demo fixtures only", "scanner": "ready", "cache-refresh": "ready", "ai": "ready"}
	now := time.Now().UnixMilli()
	for _, k := range []string{"quotes", "vix", "history", "news", "earnings", "filings", "fundamentals", "global", "macro-rates", "macro-events", "options"} {
		e.lastUpdated[k] = now
	}
	news := clone(e.news)
	earnings := clone(e.earnings)
	filings := clone(e.filings)
	e.mu.Unlock()
	e.app.broadcastNews(news)
	e.app.broadcastEarnings(earnings)
	e.app.broadcastFilings(filings, nil)
}

func (e *Engine) seedDemoBars(symbol string, base float64, needIntraday, needDaily, fullDaily, needWeekly, needFundamentals bool) {
	e.mu.Lock()
	defer e.mu.Unlock()
	if e.bars[symbol] == nil {
		e.bars[symbol] = map[string][]Bar{}
	}
	now := time.Now().UTC()
	makeBars := func(count int, step time.Duration, drift, vol float64) []Bar {
		out := make([]Bar, 0, count)
		price := base * (1 - drift*float64(count)/2)
		for i := count - 1; i >= 0; i-- {
			t := now.Add(-time.Duration(i) * step)
			move := drift + (mathrand.Float64()-.48)*vol
			open := price
			close := math.Max(.01, open*(1+move))
			high := math.Max(open, close) * (1 + mathrand.Float64()*vol*.35)
			low := math.Min(open, close) * (1 - mathrand.Float64()*vol*.35)
			volume := 1_000_000 + mathrand.Float64()*8_000_000
			out = append(out, Bar{T: t.Unix(), O: open, H: high, L: low, C: close, V: volume})
			price = close
		}
		if len(out) > 0 {
			scale := base / out[len(out)-1].C
			for i := range out {
				out[i].O *= scale
				out[i].H *= scale
				out[i].L *= scale
				out[i].C *= scale
			}
		}
		return out
	}
	if needIntraday {
		e.bars[symbol]["intraday"] = makeBars(96, 15*time.Minute, .00015, .006)
	}
	if needDaily {
		count := 30
		if fullDaily {
			count = 220
		}
		e.bars[symbol]["daily"] = makeBars(count, 24*time.Hour, .00045, .025)
	}
	if needWeekly {
		e.bars[symbol]["weekly"] = makeBars(80, 7*24*time.Hour, .0015, .055)
	}
	if needFundamentals {
		e.fundamentals[symbol] = FundamentalSnapshot{Symbol: symbol, MarketCap: 20_000 + mathrand.Float64()*2_000_000, PERatio: 12 + mathrand.Float64()*45, ForwardPERatio: 11 + mathrand.Float64()*42, PSRatio: 2 + mathrand.Float64()*18, PEGRatio: .5 + mathrand.Float64()*3.5, RevenueGrowth: -5 + mathrand.Float64()*35, EPSGrowth: -10 + mathrand.Float64()*50, GrossMargin: 25 + mathrand.Float64()*50, OperatingMargin: 5 + mathrand.Float64()*35, ROE: 5 + mathrand.Float64()*40, NetMargin: 4 + mathrand.Float64()*30, DebtToEquity: 20 + mathrand.Float64()*160, CurrentRatio: .7 + mathrand.Float64()*2.5, FreeCashFlow: mathrand.Float64() * 20_000_000_000, DividendYield: mathrand.Float64() * 3, FiftyTwoWeekHigh: base * 1.18, FiftyTwoWeekLow: base * .62, UpdatedAt: time.Now().UnixMilli(), Source: "demo"}
	}
}
