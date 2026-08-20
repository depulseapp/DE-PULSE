package main

import (
	"context"
	"fmt"
	"math"
	"net/http"
	"net/url"
	"sort"
	"strings"
	"sync"
	"time"
)

type alpacaAsset struct {
	Symbol   string `json:"symbol"`
	Status   string `json:"status"`
	Tradable bool   `json:"tradable"`
	Exchange string `json:"exchange"`
}

func (e *Engine) scannerUniverse(ctx context.Context, key, secret string) []string {
	return e.canonicalUSSymbolUniverse(ctx, key, secret, time.Now())
}

var discoverySeedUniverse = []string{"AAPL", "MSFT", "NVDA", "AMZN", "META", "GOOGL", "TSLA", "AMD", "AVGO", "NFLX", "ORCL", "PLTR", "CRM", "ADBE", "INTC", "MU", "QCOM", "SMCI", "ARM", "SOFI", "JPM", "BAC", "GS", "V", "MA", "XOM", "CVX", "LLY", "UNH", "COST", "WMT", "HD", "NKE", "DIS", "UBER", "ABNB", "SHOP", "PYPL", "COIN", "MSTR", "RBLX", "SNOW", "CRWD", "PANW", "DDOG", "NET", "MDB", "NOW", "IBM", "DELL", "MRVL", "AMAT", "LRCX", "KLAC", "TSM", "ASML", "GE", "CAT", "BA", "GM", "F", "RIVN", "LCID", "HOOD", "ROKU", "SQ", "AFRM", "DKNG", "CELH", "CAVA", "RDDT", "APP", "HIMS", "DUOL", "TTD", "CVNA", "ARM", "SMR", "OKLO", "NNE", "IONQ", "RGTI", "QBTS", "RKLB", "LUNR", "ASTS", "ACHR", "JOBY", "CRWV", "NBIS", "TEM", "SOUN", "PATH", "AI", "BBAI", "MARA", "RIOT", "CLSK", "IREN", "WULF", "GLD", "SLV", "TLT", "IWM", "QQQ", "SPY"}

func scannerScoreFromSnapshot(symbol, mode string, snap alpacaLiveSnapshot) ScannerResult {
	price := snap.LatestTrade.Price
	if price <= 0 {
		price = snap.MinuteBar.Close
	}
	if price <= 0 && snap.LatestQuote.Bid > 0 && snap.LatestQuote.Ask > 0 {
		price = (snap.LatestQuote.Bid + snap.LatestQuote.Ask) / 2
	}
	prev := snap.PrevDailyBar.Close
	dayClose := snap.DailyBar.Close
	if dayClose <= 0 {
		dayClose = price
	}
	dayOpen := snap.DailyBar.Open
	if dayOpen <= 0 {
		dayOpen = dayClose
	}
	changePct, gapPct := 0.0, 0.0
	if prev > 0 && price > 0 {
		changePct = (price - prev) / prev * 100
	}
	if prev > 0 && dayOpen > 0 {
		gapPct = (dayOpen - prev) / prev * 100
	}
	spreadPct := 0.0
	if price > 0 && snap.LatestQuote.Ask > 0 && snap.LatestQuote.Bid > 0 {
		spreadPct = (snap.LatestQuote.Ask - snap.LatestQuote.Bid) / price * 100
	}
	dollarVol := price * snap.DailyBar.Volume
	relVol := 0.0
	if snap.PrevDailyBar.Volume > 0 && snap.DailyBar.Volume > 0 {
		relVol = snap.DailyBar.Volume / snap.PrevDailyBar.Volume
	}
	trendScore := clampScanner(changePct*10, -100, 100)
	momentumScore := clampScanner((changePct+gapPct*.35)*12, -100, 100)
	score := 50.0
	reasons := []string{}
	switch mode {
	case "day":
		score += math.Min(24, math.Abs(changePct)*4)
		if spreadPct > 0 && spreadPct < .25 {
			score += 8
			reasons = append(reasons, "Tight current spread")
		}
		if relVol >= 1.25 {
			score += math.Min(12, (relVol-1)*12)
			reasons = append(reasons, fmt.Sprintf("%.2fx volume vs prior session", relVol))
		}
		if math.Abs(changePct) >= 2 {
			reasons = append(reasons, fmt.Sprintf("%.1f%% session move", changePct))
		}
		if math.Abs(gapPct) >= 1.5 {
			score += math.Min(10, math.Abs(gapPct)*2)
			reasons = append(reasons, "Meaningful opening gap")
		}
	case "swing":
		if changePct > 0 {
			score += math.Min(18, changePct*3)
			reasons = append(reasons, "Positive daily momentum")
		} else {
			score += math.Max(-16, changePct*2)
		}
		if relVol >= 1.15 {
			score += math.Min(8, (relVol-1)*10)
			reasons = append(reasons, "Above-prior-session volume")
		}
		if spreadPct > 0 && spreadPct < .35 {
			score += 5
		}
	case "long":
		if changePct > 0 {
			score += math.Min(8, changePct)
		}
		score += 3
		reasons = append(reasons, "Queued for trend and fundamental enrichment")
	}
	if spreadPct > 1.5 {
		score -= 20
		reasons = append(reasons, "Wide spread")
	}
	if price < 2 {
		score -= 25
		reasons = append(reasons, "Low-priced security")
	}
	if score < 0 {
		score = 0
	}
	if score > 100 {
		score = 100
	}
	stamp := latestProviderEvidenceMillis(snap.LatestTrade.Time, snap.LatestQuote.Time)
	return ScannerResult{Symbol: normalizeSymbol(symbol), Mode: mode, Price: price, ChangePercent: changePct, GapPercent: gapPct, RelativeVolume: relVol, DollarVolume: dollarVol, SpreadPercent: spreadPct, TrendScore: trendScore, MomentumScore: momentumScore, Score: score, Reasons: reasons, Provider: "alpaca", UpdatedAt: stamp}
}

func clampScanner(v, lo, hi float64) float64 {
	if v < lo {
		return lo
	}
	if v > hi {
		return hi
	}
	return v
}

func scannerSMA(values []float64, n int) float64 {
	if n <= 0 || len(values) < n {
		return 0
	}
	var total float64
	for _, v := range values[len(values)-n:] {
		total += v
	}
	return total / float64(n)
}

func scannerReturn(values []float64, periods int) float64 {
	if len(values) < 2 {
		return 0
	}
	if periods <= 0 || periods >= len(values) {
		periods = len(values) - 1
	}
	base := values[len(values)-1-periods]
	if base <= 0 {
		return 0
	}
	return (values[len(values)-1]/base - 1) * 100
}

func scannerRSI(values []float64, n int) float64 {
	if n <= 0 || len(values) < n+1 {
		return 50
	}
	var gains, losses float64
	start := len(values) - n
	for i := start; i < len(values); i++ {
		d := values[i] - values[i-1]
		if d >= 0 {
			gains += d
		} else {
			losses -= d
		}
	}
	if losses == 0 {
		return 100
	}
	rs := gains / losses
	return 100 - 100/(1+rs)
}

// applyScannerSessionRelativeStrength reuses the already-fetched discovery
// snapshots. It introduces no provider call and only adjusts the non-
// deterministic Discovery rank. Day RS is the current-session move relative
// to SPY/QQQ; deterministic desk Score/Action formulas remain untouched.
func applyScannerSessionRelativeStrength(results []ScannerResult, mode string) []ScannerResult {
	if len(results) == 0 || (mode != "day" && mode != "swing") {
		return results
	}
	bench := []float64{}
	for _, row := range results {
		if row.Symbol == "SPY" || row.Symbol == "QQQ" {
			bench = append(bench, row.ChangePercent)
		}
	}
	if len(bench) == 0 {
		return results
	}
	avg := 0.0
	for _, v := range bench {
		avg += v
	}
	avg /= float64(len(bench))
	for i := range results {
		rs := results[i].ChangePercent - avg
		results[i].RelativeStrength = rs
		results[i].RSBenchmark = "SPY/QQQ SESSION"
		if mode == "day" {
			results[i].Score = clampScanner(results[i].Score+clampScanner(rs*1.8, -8, 8), 0, 100)
			if math.Abs(rs) >= 1 {
				results[i].Reasons = append(results[i].Reasons, fmt.Sprintf("Session RS %+.1f%% vs SPY/QQQ", rs))
			}
		}
	}
	return results
}

func (e *Engine) enrichScannerHistory(ctx context.Context, key, secret, mode string, results []ScannerResult) []ScannerResult {
	if len(results) == 0 || key == "" || secret == "" {
		return results
	}
	limit := minInt(len(results), 60)
	syms := make([]string, 0, limit)
	for _, x := range results[:limit] {
		syms = append(syms, x.Symbol)
	}
	// SPY/QQQ are benchmark evidence for Swing RS. Reuse the same historical
	// request rather than creating per-symbol benchmark fetches.
	syms = uniqueSymbols(append(syms, "SPY", "QQQ"))
	start := time.Now().AddDate(-2, 0, 0).UTC().Format(time.RFC3339)
	raw := "https://data.alpaca.markets/v2/stocks/bars?symbols=" + url.QueryEscape(strings.Join(syms, ",")) + "&timeframe=1Day&start=" + url.QueryEscape(start) + "&limit=10000&adjustment=all&feed=iex&sort=asc"
	var payload struct {
		Bars map[string][]struct {
			C float64 `json:"c"`
			V float64 `json:"v"`
		} `json:"bars"`
	}
	client := &http.Client{Timeout: 25 * time.Second}
	if err := e.providerGetJSONTier(ctx, "Alpaca", WorkTierBroadDiscovery, client, raw, map[string]string{"APCA-API-KEY-ID": key, "APCA-API-SECRET-KEY": secret}, &payload); err != nil {
		return results
	}
	index := map[string]int{}
	for i := range results {
		index[results[i].Symbol] = i
	}
	benchmarkReturns := []float64{}
	for _, benchmark := range []string{"SPY", "QQQ"} {
		rows := payload.Bars[benchmark]
		closes := make([]float64, 0, len(rows))
		for _, row := range rows {
			if row.C > 0 {
				closes = append(closes, row.C)
			}
		}
		if len(closes) >= 21 {
			benchmarkReturns = append(benchmarkReturns, scannerReturn(closes, 20))
		}
	}
	benchmark20 := 0.0
	if len(benchmarkReturns) > 0 {
		for _, v := range benchmarkReturns {
			benchmark20 += v
		}
		benchmark20 /= float64(len(benchmarkReturns))
	}
	for sym, rows := range payload.Bars {
		i, ok := index[normalizeSymbol(sym)]
		if !ok || len(rows) < 5 {
			continue
		}
		closes := make([]float64, 0, len(rows))
		vols := make([]float64, 0, len(rows))
		for _, row := range rows {
			if row.C > 0 {
				closes = append(closes, row.C)
				vols = append(vols, row.V)
			}
		}
		if len(closes) < 5 {
			continue
		}
		last := closes[len(closes)-1]
		s20, s50, s200 := scannerSMA(closes, 20), scannerSMA(closes, 50), scannerSMA(closes, 200)
		r20, r60, r126 := scannerReturn(closes, 20), scannerReturn(closes, 60), scannerReturn(closes, 126)
		rsi := scannerRSI(closes, 14)
		trend := 0.0
		votes := 0.0
		addVote := func(cond bool, weight float64) {
			votes += weight
			if cond {
				trend += weight
			} else {
				trend -= weight
			}
		}
		if s20 > 0 {
			addVote(last > s20, 1)
		}
		if s50 > 0 {
			addVote(last > s50, 1.2)
		}
		if s20 > 0 && s50 > 0 {
			addVote(s20 > s50, 1)
		}
		if s200 > 0 {
			addVote(last > s200, 1.4)
		}
		if s50 > 0 && s200 > 0 {
			addVote(s50 > s200, 1.4)
		}
		trendScore := 0.0
		if votes > 0 {
			trendScore = clampScanner(trend/votes*100, -100, 100)
		}
		momentumScore := clampScanner((r20*.55+r60*.25+r126*.20)*3+(rsi-50)*1.1, -100, 100)
		rv := results[i].RelativeVolume
		if len(vols) >= 21 {
			var total float64
			for _, v := range vols[len(vols)-21 : len(vols)-1] {
				total += v
			}
			avg20 := total / 20
			if avg20 > 0 {
				rv = vols[len(vols)-1] / avg20
			}
		}
		results[i].RSI = rsi
		results[i].TrendScore = trendScore
		results[i].MomentumScore = momentumScore
		results[i].RelativeVolume = rv
		reasons := results[i].Reasons
		swingRSAdjust := 0.0
		if mode == "swing" && len(benchmarkReturns) > 0 && len(closes) >= 21 {
			rs := r20 - benchmark20
			results[i].RelativeStrength = rs
			results[i].RSBenchmark = "SPY/QQQ 20D"
			swingRSAdjust = clampScanner(rs*1.25, -10, 10)
			if math.Abs(rs) >= 2 {
				reasons = append(reasons, fmt.Sprintf("20D RS %+.1f%% vs SPY/QQQ", rs))
			}
		}
		if trendScore >= 35 {
			reasons = append(reasons, "Trend aligned above key moving averages")
		}
		if trendScore <= -35 {
			reasons = append(reasons, "Trend below key moving averages")
		}
		if rsi >= 55 && rsi <= 75 {
			reasons = append(reasons, fmt.Sprintf("RSI %.0f supports momentum", rsi))
		}
		if rv >= 1.3 {
			reasons = append(reasons, fmt.Sprintf("%.2fx 20-day volume", rv))
		}
		results[i].Reasons = reasons
		switch mode {
		case "day":
			results[i].Score = clampScanner(results[i].Score*.62+(50+momentumScore*.30+math.Min(20, math.Max(0, (rv-1)*18)))*.38, 0, 100)
		case "swing":
			results[i].Score = clampScanner(50+trendScore*.28+momentumScore*.18+math.Min(12, math.Max(-12, r20*.55))+math.Min(8, math.Max(0, (rv-1)*8))+swingRSAdjust, 0, 100)
		case "long":
			results[i].Score = clampScanner(50+trendScore*.34+clampScanner(r126*1.2, -18, 18)+clampScanner(r60*.8, -10, 10), 0, 100)
		}
	}
	return results
}

func scannerFundamentalScore(metric map[string]any) float64 {
	rev := clampScanner(toFloat(metric["revenueGrowthTTMYoy"]), -50, 80)
	eps := clampScanner(toFloat(metric["epsGrowthTTMYoy"]), -80, 120)
	roe := clampScanner(toFloat(metric["roeTTM"]), -30, 60)
	margin := clampScanner(toFloat(metric["netProfitMarginTTM"]), -30, 50)
	de := toFloat(metric["totalDebt/totalEquityAnnual"])
	dePenalty := clampScanner((de-100)*.20, 0, 30)
	growth := clampScanner(50+rev*.55+eps*.30, 0, 100)
	quality := clampScanner(50+roe*.75+margin*.55-dePenalty, 0, 100)
	return clampScanner(growth*.52+quality*.48, 0, 100)
}

func enrichScannerFundamentals(a *Application, ctx context.Context, key string, results []ScannerResult) []ScannerResult {
	if len(results) == 0 || strings.TrimSpace(key) == "" {
		return results
	}
	limit := minInt(len(results), 12)
	type scored struct {
		idx      int
		score    float64
		snapshot FundamentalSnapshot
	}
	ch := make(chan scored, limit)
	sem := make(chan struct{}, 3)
	var wg sync.WaitGroup
	for i := 0; i < limit; i++ {
		i := i
		wg.Add(1)
		go func() {
			defer wg.Done()
			sem <- struct{}{}
			defer func() { <-sem }()
			var payload struct {
				Metric map[string]any `json:"metric"`
			}
			if err := a.engine.finnhubJSONForSymbol(ctx, key, results[i].Symbol, "/stock/metric?symbol="+url.QueryEscape(results[i].Symbol)+"&metric=all", &payload); err != nil {
				return
			}
			m := payload.Metric
			f := FundamentalSnapshot{Symbol: results[i].Symbol, MarketCap: toFloat(m["marketCapitalization"]), PERatio: toFloat(m["peTTM"]), ForwardPERatio: toFloat(m["forwardPEAnnual"]), PSRatio: toFloat(m["psTTM"]), PEGRatio: toFloat(m["pegRatio"]), RevenueGrowth: toFloat(m["revenueGrowthTTMYoy"]), EPSGrowth: toFloat(m["epsGrowthTTMYoy"]), GrossMargin: toFloat(m["grossMarginTTM"]), OperatingMargin: toFloat(m["operatingMarginTTM"]), ROE: toFloat(m["roeTTM"]), NetMargin: toFloat(m["netProfitMarginTTM"]), DebtToEquity: toFloat(m["totalDebt/totalEquityAnnual"]), CurrentRatio: toFloat(m["currentRatioAnnual"]), FreeCashFlow: toFloat(m["freeCashFlowTTM"]), DividendYield: toFloat(m["dividendYieldIndicatedAnnual"]), FiftyTwoWeekHigh: toFloat(m["52WeekHigh"]), FiftyTwoWeekLow: toFloat(m["52WeekLow"]), UpdatedAt: time.Now().UnixMilli(), Source: "finnhub"}
			ch <- scored{idx: i, score: scannerFundamentalScore(m), snapshot: f}
		}()
	}
	wg.Wait()
	close(ch)
	for x := range ch {
		results[x.idx].FundamentalScore = x.score
		results[x.idx].Score = clampScanner(results[x.idx].Score*.68+x.score*.32, 0, 100)
		results[x.idx].Reasons = append(results[x.idx].Reasons, fmt.Sprintf("Fundamental quality %.0f/100", x.score))
		a.engine.mu.Lock()
		a.engine.fundamentals[x.snapshot.Symbol] = x.snapshot
		a.engine.mu.Unlock()
	}
	return results
}

func demoScannerResults(mode string) []ScannerResult {
	base := []struct {
		S               string
		P, C, G, Sp, RV float64
		Trend, Momentum float64
		Fundamental     float64
	}{
		{"AMD", 178.42, 3.8, 2.1, .08, 1.82, 72, 78, 68}, {"AVGO", 337.80, 2.4, 1.2, .05, 1.44, 88, 69, 84}, {"CRWD", 471.24, -2.2, -1.5, .10, 1.38, 64, -32, 78}, {"HOOD", 118.55, 5.6, 3.0, .09, 2.24, 58, 90, 61}, {"COIN", 331.18, 4.4, 2.2, .12, 1.95, 67, 84, 55}, {"RDDT", 221.37, 3.1, 1.4, .11, 1.57, 74, 71, 66}, {"APP", 498.33, 2.7, .9, .08, 1.32, 91, 62, 86}, {"HIMS", 63.82, -3.5, -2.2, .18, 1.71, 49, -58, 73}, {"RKLB", 49.62, 6.2, 4.1, .22, 2.63, 57, 93, 48}, {"CAVA", 92.48, 2.0, .8, .16, 1.26, 77, 55, 72}, {"PANW", 207.24, 1.8, .5, .06, 1.12, 86, 51, 82}, {"MU", 154.17, 3.4, 1.6, .07, 1.68, 79, 76, 71}, {"UBER", 94.81, 1.5, .4, .06, 1.08, 83, 43, 76}, {"SNOW", 226.73, -1.4, -.8, .10, 1.16, 54, -18, 69}, {"MSTR", 421.86, 4.8, 2.6, .13, 1.91, 63, 88, 41},
	}
	out := make([]ScannerResult, 0, len(base))
	for i, x := range base {
		score := 88 - float64(i)*1.6
		reasons := []string{"Liquidity filter passed"}
		switch mode {
		case "day":
			score = clampScanner(48+math.Abs(x.C)*4+x.RV*9-math.Max(0, x.Sp-.15)*18, 0, 100)
			reasons = append(reasons, fmt.Sprintf("%.2fx relative volume", x.RV), "Intraday momentum candidate")
		case "swing":
			score = clampScanner(50+x.Trend*.28+x.Momentum*.15, 0, 100)
			reasons = append(reasons, "Daily trend and momentum screen", fmt.Sprintf("Trend %.0f/100", x.Trend))
		case "long":
			score = clampScanner(45+x.Trend*.28+x.Fundamental*.34, 0, 100)
			reasons = append(reasons, "Long-term trend screen", fmt.Sprintf("Fundamental quality %.0f/100", x.Fundamental))
		}
		out = append(out, ScannerResult{Symbol: x.S, Mode: mode, Price: x.P, ChangePercent: x.C, GapPercent: x.G, RelativeVolume: x.RV, SpreadPercent: x.Sp, RSI: clampScanner(50+x.Momentum*.28, 20, 85), TrendScore: x.Trend, MomentumScore: x.Momentum, FundamentalScore: x.Fundamental, Score: score, Reasons: reasons, Provider: "demo", UpdatedAt: time.Now().UnixMilli()})
	}
	sort.Slice(out, func(i, j int) bool { return out[i].Score > out[j].Score })
	return out
}

func (a *Application) handleDiscoveryScan(w http.ResponseWriter, r *http.Request) {
	var in struct {
		Mode string `json:"mode"`
	}
	if decodeJSON(r, &in) != nil {
		writeError(w, 400, "Invalid discovery request")
		return
	}
	mode := strings.ToLower(strings.TrimSpace(in.Mode))
	if mode != "day" && mode != "swing" && mode != "long" {
		mode = "day"
	}
	started := time.Now()
	a.engine.mu.Lock()
	radar := clone(a.engine.scanner.Radar)
	a.engine.scanner = ScannerState{Mode: mode, Status: "running", Message: "Scanning the market…", Results: []ScannerResult{}, Radar: radar}
	a.engine.health["scanner"] = "manual scan running · Opportunity Radar preserved"
	a.engine.mu.Unlock()
	a.hub.Broadcast(map[string]any{"type": "scanner", "scanner": a.engine.Snapshot().Scanner})
	a.mu.RLock()
	dataMode := a.state.Settings.DataMode
	key := strings.TrimSpace(a.secrets.AlpacaKey)
	secret := strings.TrimSpace(a.secrets.AlpacaSecret)
	finnhubKey := strings.TrimSpace(a.secrets.Finnhub)
	a.mu.RUnlock()
	var results []ScannerResult
	scanned := 0
	if dataMode != "live" || key == "" || secret == "" {
		results = demoScannerResults(mode)
		scanned = len(discoverySeedUniverse)
	} else {

		releaseScanner, acquired := a.engine.workload.Acquire(r.Context(), "scanner")
		if !acquired {
			writeError(w, 499, "Discovery scan canceled before scanner capacity became available")
			return
		}
		defer releaseScanner()
		client := &http.Client{Timeout: 20 * time.Second}
		universe := a.engine.scannerUniverse(r.Context(), key, secret)
		feed := "iex"
		session := marketSessionET(time.Now())
		if session == "overnight" {
			feed = "overnight"
		}
		for start := 0; start < len(universe); start += 50 {
			end := minInt(start+50, len(universe))
			batch := universe[start:end]
			raw := "https://data.alpaca.markets/v2/stocks/snapshots?symbols=" + url.QueryEscape(strings.Join(batch, ",")) + "&feed=" + url.QueryEscape(feed)
			var payload map[string]alpacaLiveSnapshot
			if err := a.engine.providerGetJSONTier(r.Context(), "Alpaca", WorkTierBroadDiscovery, client, raw, map[string]string{"APCA-API-KEY-ID": key, "APCA-API-SECRET-KEY": secret}, &payload); err != nil {
				continue
			}
			scanned += len(batch)
			for symbol, snap := range payload {
				x := scannerScoreFromSnapshot(symbol, mode, snap)
				if x.Price >= 2 && x.SpreadPercent <= 2.5 {
					results = append(results, x)
				}
			}
		}
		results = applyScannerSessionRelativeStrength(results, mode)
		sort.Slice(results, func(i, j int) bool { return results[i].Score > results[j].Score })
		if len(results) > 80 {
			results = results[:80]
		}
		results = a.engine.enrichScannerHistory(r.Context(), key, secret, mode, results)
		if mode == "long" && finnhubKey != "" {
			results = enrichScannerFundamentals(a, r.Context(), finnhubKey, results)
		}
		sort.Slice(results, func(i, j int) bool { return results[i].Score > results[j].Score })
		if len(results) > 60 {
			results = results[:60]
		}
	}
	dur := time.Since(started).Milliseconds()
	now := time.Now().UnixMilli()
	a.engine.mu.RLock()
	radar = clone(a.engine.scanner.Radar)
	a.engine.mu.RUnlock()
	state := ScannerState{Mode: mode, Status: "complete", Message: fmt.Sprintf("%d candidates ranked from %d scanned symbols.", len(results), scanned), Results: results, Scanned: scanned, DurationMs: dur, UpdatedAt: now, Radar: radar}
	a.engine.mu.Lock()
	a.engine.scanner = state
	a.engine.health["scanner"] = "complete"
	a.engine.lastUpdated["scanner"] = now
	a.engine.mu.Unlock()
	a.hub.Broadcast(map[string]any{"type": "scanner", "scanner": state})
	writeJSON(w, 200, state)
}

type MaintenanceCheck struct {
	ID         string `json:"id"`
	Label      string `json:"label"`
	Status     string `json:"status"`
	Detail     string `json:"detail"`
	DurationMs int64  `json:"durationMs,omitempty"`
}
