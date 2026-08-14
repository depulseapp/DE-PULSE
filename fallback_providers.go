package main

import (
	"context"
	"fmt"
	"math"
	"net/http"
	"net/url"
	"sort"
	"strconv"
	"strings"
	"time"
)

// Public Yahoo chart recovery is intentionally fallback-only.
// custom response avoids the awkward repeated JSON tags above.
type yahooChartEnvelope struct {
	Chart struct {
		Result []struct {
			Meta struct {
				Symbol             string  `json:"symbol"`
				RegularMarketPrice float64 `json:"regularMarketPrice"`
				PreviousClose      float64 `json:"chartPreviousClose"`
				RegularMarketTime  int64   `json:"regularMarketTime"`
				ExchangeName       string  `json:"exchangeName"`
				FullExchangeName   string  `json:"fullExchangeName"`
				InstrumentType     string  `json:"instrumentType"`
				Currency           string  `json:"currency"`
			} `json:"meta"`
			Timestamp  []int64 `json:"timestamp"`
			Indicators struct {
				Quote []struct {
					Open   []*float64 `json:"open"`
					High   []*float64 `json:"high"`
					Low    []*float64 `json:"low"`
					Close  []*float64 `json:"close"`
					Volume []*float64 `json:"volume"`
				} `json:"quote"`
			} `json:"indicators"`
		} `json:"result"`
		Error any `json:"error"`
	} `json:"chart"`
}

func fetchYahooChart(ctx context.Context, symbol, rangeText, interval string) (yahooChartEnvelope, error) {
	var out yahooChartEnvelope
	raw := "https://query1.finance.yahoo.com/v8/finance/chart/" + url.PathEscape(symbol) + "?range=" + url.QueryEscape(rangeText) + "&interval=" + url.QueryEscape(interval) + "&includePrePost=true&events=div%2Csplits"
	err := getJSON(ctx, &http.Client{Timeout: 15 * time.Second}, raw, map[string]string{"User-Agent": "Mozilla/5.0 DE.PULSE/15.0"}, &out)
	if err != nil {
		return out, err
	}
	if len(out.Chart.Result) == 0 {
		return out, fmt.Errorf("yfinance chart returned no result")
	}
	return out, nil
}
func floatPtrValue(xs []*float64, i int) float64 {
	if i < 0 || i >= len(xs) || xs[i] == nil {
		return 0
	}
	return *xs[i]
}

func (e *Engine) refreshYahooVIX(ctx context.Context) bool {
	if !e.providerAllowed("yfinance") {
		return false
	}
	out, err := fetchYahooChart(ctx, "^VIX", "5d", "5m")
	if err != nil {
		e.recordProviderFailure("yfinance", err)
		return false
	}
	r := out.Chart.Result[0]
	px := r.Meta.RegularMarketPrice
	if px < 5 || px > 200 {
		e.recordProviderFailure("yfinance", fmt.Errorf("invalid VIX price %.4f", px))
		return false
	}
	ts := normalizeObservationMs(r.Meta.RegularMarketTime)
	e.updateQuote("VIX", Quote{Price: px, PreviousClose: r.Meta.PreviousClose, ProviderTimestamp: ts, FeedType: "index-recovery", DataState: "delayed"}, "yfinance-vix:^VIX")
	e.recordProviderSuccess("yfinance")
	e.setHealth("vix", "Delayed VIX recovery · yfinance · ^VIX")
	return true
}

func historyRouteSymbols(e *Engine, only []string) []string {
	symbols := uniqueSymbols(only)
	if len(symbols) == 0 {
		e.app.mu.RLock()
		symbols = activeDeskSymbolsFromState(e.app.state)
		e.app.mu.RUnlock()
	}
	return symbols
}

func parseTwelveHistoryRows(values []struct{ Datetime, Open, High, Low, Close, Volume string }, layout string, loc *time.Location) []Bar {
	bars := []Bar{}
	for i := len(values) - 1; i >= 0; i-- {
		v := values[i]
		var dt time.Time
		var err error
		if loc != nil {
			dt, err = time.ParseInLocation(layout, v.Datetime, loc)
		} else {
			dt, err = time.Parse(layout, v.Datetime)
		}
		if err != nil {
			continue
		}
		op, _ := strconv.ParseFloat(v.Open, 64)
		hi, _ := strconv.ParseFloat(v.High, 64)
		lo, _ := strconv.ParseFloat(v.Low, 64)
		cl, _ := strconv.ParseFloat(v.Close, 64)
		vol, _ := strconv.ParseFloat(v.Volume, 64)
		if cl <= 0 {
			continue
		}
		bars = append(bars, Bar{T: dt.Unix(), O: op, H: hi, L: lo, C: cl, V: vol})
	}
	return bars
}

func (e *Engine) refreshTwelveHistoryMode(ctx context.Context, only []string, mode string) int {
	symbols := historyRouteSymbols(e, only)
	e.app.mu.RLock()
	td := strings.TrimSpace(e.app.secrets.TwelveData)
	e.app.mu.RUnlock()
	if td == "" || !e.providerAllowed("Twelve Data") {
		return 0
	}
	mode = strings.ToLower(strings.TrimSpace(mode))
	if mode == "" {
		mode = "all"
	}
	type querySpec struct {
		name, interval, outputsize, layout string
		loc                                *time.Location
	}
	loc, _ := time.LoadLocation("America/New_York")
	var specs []querySpec
	if mode == "all" || mode == "intraday" {
		specs = append(specs, querySpec{name: "intraday", interval: "15min", outputsize: "260", layout: "2006-01-02 15:04:05", loc: loc})
	}
	if mode == "all" || mode == "daily" {
		specs = append(specs,
			querySpec{name: "daily", interval: "1day", outputsize: "520", layout: "2006-01-02"},
			querySpec{name: "weekly", interval: "1week", outputsize: "320", layout: "2006-01-02"},
		)
	}
	loaded := 0
	intradayLoaded, dailyLoaded := false, false
	for _, sym := range symbols {
		for _, sp := range specs {
			// Preserve the same single canonical history owner used by Alpaca while
			// giving SPY/QQQ enough daily depth for v16.3 seasonality when Twelve
			// Data is the routed fallback.
			if sp.name == "daily" && (sym == "SPY" || sym == "QQQ") {
				sp.outputsize = "3000"
			}
			var payload struct {
				Values  []struct{ Datetime, Open, High, Low, Close, Volume string } `json:"values"`
				Status  string                                                      `json:"status"`
				Message string                                                      `json:"message"`
			}
			raw := "https://api.twelvedata.com/time_series?symbol=" + url.QueryEscape(sym) + "&interval=" + url.QueryEscape(sp.interval) + "&outputsize=" + sp.outputsize + "&apikey=" + url.QueryEscape(td)
			if err := getJSON(ctx, &http.Client{Timeout: 12 * time.Second}, raw, nil, &payload); err != nil || len(payload.Values) == 0 {
				if err != nil {
					e.recordProviderFailure("Twelve Data", err)
				}
				continue
			}
			bars := parseTwelveHistoryRows(payload.Values, sp.layout, sp.loc)
			if len(bars) == 0 {
				continue
			}
			e.mu.Lock()
			if e.bars[sym] == nil {
				e.bars[sym] = map[string][]Bar{}
			}
			e.bars[sym][sp.name] = bars
			e.mu.Unlock()
			if sp.name == "daily" {
				prior := 0.0
				if len(bars) > 1 {
					prior = bars[len(bars)-2].C
				}
				e.updateCanonicalSessionClose(sym, bars[len(bars)-1].C, bars[len(bars)-1].T*1000, prior)
			}
			loaded += len(bars)
			if sp.name == "intraday" {
				intradayLoaded = true
			} else {
				dailyLoaded = true
			}
			e.recordProviderSuccess("Twelve Data")
		}
	}
	if loaded > 0 {
		nowMs := time.Now().UnixMilli()
		e.mu.Lock()
		e.lastUpdated["history"] = nowMs
		if intradayLoaded {
			e.lastUpdated["history-intraday"] = nowMs
		}
		if dailyLoaded {
			e.lastUpdated["history-daily"] = nowMs
		}
		e.health["history"] = "healthy · Twelve Data fallback"
		e.mu.Unlock()
	}
	return loaded
}

func (e *Engine) refreshYahooHistoryMode(ctx context.Context, only []string, mode string) int {
	if !e.providerAllowed("yfinance") {
		return 0
	}
	symbols := historyRouteSymbols(e, only)
	mode = strings.ToLower(strings.TrimSpace(mode))
	if mode == "" {
		mode = "all"
	}
	type yahooHistorySpec struct{ name, rangeText, interval string }
	var specs []yahooHistorySpec
	if mode == "all" || mode == "intraday" {
		specs = append(specs, yahooHistorySpec{name: "intraday", rangeText: "5d", interval: "15m"})
	}
	if mode == "all" || mode == "daily" {
		specs = append(specs,
			yahooHistorySpec{name: "daily", rangeText: "2y", interval: "1d"},
			yahooHistorySpec{name: "weekly", rangeText: "5y", interval: "1wk"},
		)
	}
	loaded := 0
	intradayLoaded, dailyLoaded := false, false
	for _, sym := range symbols {
		for _, sp := range specs {
			if sp.name == "daily" && (sym == "SPY" || sym == "QQQ") {
				sp.rangeText = "10y"
			}
			out, err := fetchYahooChart(ctx, sym, sp.rangeText, sp.interval)
			if err != nil {
				e.recordProviderFailure("yfinance", err)
				continue
			}
			r := out.Chart.Result[0]
			if len(r.Indicators.Quote) == 0 {
				continue
			}
			qv := r.Indicators.Quote[0]
			bars := []Bar{}
			for i, ts := range r.Timestamp {
				cl := floatPtrValue(qv.Close, i)
				if cl <= 0 {
					continue
				}
				bars = append(bars, Bar{T: ts, O: floatPtrValue(qv.Open, i), H: floatPtrValue(qv.High, i), L: floatPtrValue(qv.Low, i), C: cl, V: floatPtrValue(qv.Volume, i)})
			}
			if len(bars) == 0 {
				continue
			}
			e.mu.Lock()
			if e.bars[sym] == nil {
				e.bars[sym] = map[string][]Bar{}
			}
			e.bars[sym][sp.name] = bars
			e.mu.Unlock()
			if sp.name == "daily" {
				prior := 0.0
				if len(bars) > 1 {
					prior = bars[len(bars)-2].C
				}
				e.updateCanonicalSessionClose(sym, bars[len(bars)-1].C, bars[len(bars)-1].T*1000, prior)
			}
			loaded += len(bars)
			if sp.name == "intraday" {
				intradayLoaded = true
			} else {
				dailyLoaded = true
			}
			e.recordProviderSuccess("yfinance")
		}
	}
	if loaded > 0 {
		nowMs := time.Now().UnixMilli()
		e.mu.Lock()
		e.lastUpdated["history"] = nowMs
		if intradayLoaded {
			e.lastUpdated["history-intraday"] = nowMs
		}
		if dailyLoaded {
			e.lastUpdated["history-daily"] = nowMs
		}
		e.health["history"] = "healthy · yfinance recovery"
		e.mu.Unlock()
	}
	return loaded
}

func (e *Engine) refreshHistoryRoutedMode(ctx context.Context, only []string, mode string) bool {
	e.app.mu.RLock()
	ak := strings.TrimSpace(e.app.secrets.AlpacaKey)
	as := strings.TrimSpace(e.app.secrets.AlpacaSecret)
	e.app.mu.RUnlock()
	label := "Historical Bars"
	switch strings.ToLower(strings.TrimSpace(mode)) {
	case "intraday":
		label = "Intraday Bars"
	case "daily":
		label = "Daily / Weekly History"
	}
	_, ok := e.executeProviderRoute(ctx, label, map[string]providerRouteAttempt{
		"Alpaca":      func(ctx context.Context) bool { return e.refreshAlpacaHistoryScopedMode(ctx, ak, as, only, mode) > 0 },
		"Twelve Data": func(ctx context.Context) bool { return e.refreshTwelveHistoryMode(ctx, only, mode) > 0 },
		"yfinance":    func(ctx context.Context) bool { return e.refreshYahooHistoryMode(ctx, only, mode) > 0 },
	})
	if ok {
		// Outcome measurement is event-driven from the canonical history commit.
		// This persists resolved ordering/MFE/MAE so rolling intraday windows cannot
		// later erase a previously established outcome.
		e.refreshValidationOutcomeState()
	}
	if !ok {
		key := "history"
		if mode == "intraday" {
			key = "history-intraday"
		} else if mode == "daily" {
			key = "history-daily"
		}
		e.setHealth(key, "degraded · provider route unavailable")
	}
	return ok
}

func (e *Engine) refreshHistoryRouted(ctx context.Context, only []string) bool {
	return e.refreshHistoryRoutedMode(ctx, only, "all")
}

func (e *Engine) refreshMarketauxNews(ctx context.Context) bool {
	e.app.mu.RLock()
	key := strings.TrimSpace(e.app.secrets.Marketaux)
	symbols := analysisSymbolsFromState(e.app.processingStateLocked())
	e.app.mu.RUnlock()
	if key == "" || !e.providerAllowed("Marketaux") {
		return false
	}
	if len(symbols) > 20 {
		symbols = symbols[:20]
	}
	raw := "https://api.marketaux.com/v1/news/all?api_token=" + url.QueryEscape(key) + "&language=en&limit=3"
	if len(symbols) > 0 {
		raw += "&symbols=" + url.QueryEscape(strings.Join(symbols, ","))
	}
	var p struct {
		Data []struct {
			UUID        string `json:"uuid"`
			Title       string `json:"title"`
			Description string `json:"description"`
			URL         string `json:"url"`
			Source      string `json:"source"`
			PublishedAt string `json:"published_at"`
			Entities    []struct {
				Symbol string `json:"symbol"`
			} `json:"entities"`
		} `json:"data"`
	}
	if err := getJSON(ctx, &http.Client{Timeout: 15 * time.Second}, raw, nil, &p); err != nil {
		e.recordProviderFailure("Marketaux", err)
		return false
	}
	items := []NewsItem{}
	for _, x := range p.Data {
		publishedAt := parseRFC3339Unix(x.PublishedAt)
		syms := []string{}
		for _, ent := range x.Entities {
			syms = append(syms, ent.Symbol)
		}
		items = append(items, NewsItem{ID: x.UUID, Datetime: publishedAt, Headline: x.Title, Summary: x.Description, Source: "Marketaux · " + x.Source, URL: x.URL, Symbols: uniqueSymbols(syms), Scope: "watchlist"})
	}
	if len(items) == 0 {

		e.mu.Lock()
		e.lastUpdated["news"] = time.Now().UnixMilli()
		e.health["news"] = "healthy · Marketaux fallback · no new articles"
		e.mu.Unlock()
		e.recordProviderSuccess("Marketaux")
		return true
	}
	e.mu.Lock()
	e.news = dedupeNews(append(items, e.news...))
	if len(e.news) > 150 {
		e.news = e.news[:150]
	}
	e.lastUpdated["news"] = time.Now().UnixMilli()
	e.health["news"] = "healthy · Marketaux fallback"
	e.mu.Unlock()
	e.recordProviderSuccess("Marketaux")
	return true
}

// SEC companyfacts is the middle fundamentals recovery hop. It supplies only
// facts that are directly derivable from filed XBRL data and never invents
// valuation metrics that SEC does not publish.
func (e *Engine) refreshSECFundamentals(ctx context.Context, symbols []string) int {
	e.app.mu.RLock()
	email := strings.TrimSpace(e.app.state.Settings.SECEmail)
	e.app.mu.RUnlock()
	if !strings.Contains(email, "@") {
		return 0
	}
	client := &http.Client{Timeout: 18 * time.Second}
	headers := map[string]string{"User-Agent": appName + "/" + appVersion + " " + email}
	var tickerMap map[string]struct {
		CIK    int    `json:"cik_str"`
		Ticker string `json:"ticker"`
	}
	if err := getJSON(ctx, client, "https://www.sec.gov/files/company_tickers.json", headers, &tickerMap); err != nil {
		e.recordProviderFailure("SEC", err)
		return 0
	}
	byTicker := map[string]string{}
	for _, row := range tickerMap {
		byTicker[strings.ToUpper(row.Ticker)] = fmt.Sprintf("%010d", row.CIK)
	}
	type factRow struct {
		End   string  `json:"end"`
		Filed string  `json:"filed"`
		Form  string  `json:"form"`
		Val   float64 `json:"val"`
	}
	type fact struct {
		Units map[string][]factRow `json:"units"`
	}
	latest := func(facts map[string]fact, tag string) (float64, bool) {
		f, ok := facts[tag]
		if !ok {
			return 0, false
		}
		rows := f.Units["USD"]
		best := ""
		v := 0.0
		found := false
		for _, r := range rows {
			form := strings.ToUpper(r.Form)
			if form != "10-Q" && form != "10-K" && form != "20-F" && form != "6-K" {
				continue
			}
			k := r.End + "|" + r.Filed
			if k > best {
				best = k
				v = r.Val
				found = true
			}
		}
		return v, found
	}
	latestTwo := func(facts map[string]fact, tag string) (float64, float64, bool) {
		f, ok := facts[tag]
		if !ok {
			return 0, 0, false
		}
		rows := append([]factRow(nil), f.Units["USD"]...)
		sort.Slice(rows, func(i, j int) bool { return rows[i].End > rows[j].End })
		vals := []float64{}
		seen := map[string]bool{}
		for _, r := range rows {
			form := strings.ToUpper(r.Form)
			if form != "10-K" && form != "20-F" {
				continue
			}
			if seen[r.End] {
				continue
			}
			seen[r.End] = true
			vals = append(vals, r.Val)
			if len(vals) == 2 {
				break
			}
		}
		if len(vals) < 2 {
			return 0, 0, false
		}
		return vals[0], vals[1], true
	}
	updated := 0
	for _, sym := range uniqueSymbols(symbols) {
		cik := byTicker[sym]
		if cik == "" {
			continue
		}
		var data struct {
			Facts map[string]map[string]fact `json:"facts"`
		}
		if err := getJSON(ctx, client, "https://data.sec.gov/api/xbrl/companyfacts/CIK"+cik+".json", headers, &data); err != nil {
			e.recordProviderFailure("SEC", err)
			continue
		}
		us := data.Facts["us-gaap"]
		if len(us) == 0 {
			continue
		}
		f := FundamentalSnapshot{Symbol: sym, UpdatedAt: time.Now().UnixMilli(), Source: "sec-companyfacts-recovery"}
		populated := false
		ac, aok := latest(us, "AssetsCurrent")
		lc, lok := latest(us, "LiabilitiesCurrent")
		if aok && lok && lc != 0 {
			f.CurrentRatio = ac / lc
			populated = true
		}
		eq, eok := latest(us, "StockholdersEquity")
		if !eok {
			eq, eok = latest(us, "StockholdersEquityIncludingPortionAttributableToNoncontrollingInterest")
		}
		debt := 0.0
		debtFound := false
		for _, tag := range []string{"LongTermDebtCurrent", "LongTermDebtNoncurrent", "LongTermDebtAndFinanceLeaseObligationsCurrent", "LongTermDebtAndFinanceLeaseObligationsNoncurrent"} {
			if v, ok := latest(us, tag); ok {
				debt += v
				debtFound = true
			}
		}
		if eok && debtFound && eq != 0 {
			f.DebtToEquity = debt / eq * 100
			populated = true
		}
		cfo, cok := latest(us, "NetCashProvidedByUsedInOperatingActivities")
		capex, xok := latest(us, "PaymentsToAcquirePropertyPlantAndEquipment")
		if cok {
			f.FreeCashFlow = cfo
			if xok {
				f.FreeCashFlow -= capex
			}
			populated = true
		}
		for _, tag := range []string{"Revenues", "RevenueFromContractWithCustomerExcludingAssessedTax", "SalesRevenueNet"} {
			if cur, prev, ok := latestTwo(us, tag); ok && prev != 0 {
				f.RevenueGrowth = (cur - prev) / prev * 100
				populated = true
				break
			}
		}
		if populated {
			e.mu.Lock()
			e.fundamentals[sym] = f
			e.mu.Unlock()
			updated++
			e.recordProviderSuccess("SEC")
		}
		if !sleepContext(ctx, 120*time.Millisecond) {
			break
		}
	}
	if updated > 0 {
		e.mu.Lock()
		e.lastUpdated["fundamentals"] = time.Now().UnixMilli()
		e.health["fundamentals"] = "healthy · SEC companyfacts recovery"
		e.mu.Unlock()
		_ = e.saveCache()
	}
	return updated
}

// Yahoo quoteSummary recovery fills only missing/failed fundamentals. It never
// supersedes a successful Finnhub or SEC observation.
func yahooPercentValue(v float64) float64 {

	if v != 0 && math.Abs(v) <= 2 {
		return v * 100
	}
	return v
}

func (e *Engine) refreshYahooFundamentals(ctx context.Context, symbols []string) int {
	if !e.providerAllowed("yfinance") {
		return 0
	}
	updated := 0
	for _, sym := range symbols {
		raw := "https://query1.finance.yahoo.com/v10/finance/quoteSummary/" + url.PathEscape(sym) + "?modules=summaryDetail%2CdefaultKeyStatistics%2CfinancialData"
		var p struct {
			QuoteSummary struct {
				Result []struct {
					SummaryDetail        map[string]any `json:"summaryDetail"`
					DefaultKeyStatistics map[string]any `json:"defaultKeyStatistics"`
					FinancialData        map[string]any `json:"financialData"`
				} `json:"result"`
			} `json:"quoteSummary"`
		}
		if err := getJSON(ctx, &http.Client{Timeout: 12 * time.Second}, raw, map[string]string{"User-Agent": "Mozilla/5.0 DE.PULSE/15.0"}, &p); err != nil || len(p.QuoteSummary.Result) == 0 {
			if err != nil {
				e.recordProviderFailure("yfinance", err)
			}
			continue
		}
		r := p.QuoteSummary.Result[0]
		rawNum := func(m map[string]any, k string) float64 {
			if v, ok := m[k].(map[string]any); ok {
				return toFloat(v["raw"])
			}
			return toFloat(m[k])
		}
		f := FundamentalSnapshot{Symbol: sym, MarketCap: rawNum(r.SummaryDetail, "marketCap"), PERatio: rawNum(r.SummaryDetail, "trailingPE"), ForwardPERatio: rawNum(r.SummaryDetail, "forwardPE"), DividendYield: yahooPercentValue(rawNum(r.SummaryDetail, "dividendYield")), FiftyTwoWeekHigh: rawNum(r.SummaryDetail, "fiftyTwoWeekHigh"), FiftyTwoWeekLow: rawNum(r.SummaryDetail, "fiftyTwoWeekLow"), DebtToEquity: rawNum(r.FinancialData, "debtToEquity"), CurrentRatio: rawNum(r.FinancialData, "currentRatio"), FreeCashFlow: rawNum(r.FinancialData, "freeCashflow"), ROE: yahooPercentValue(rawNum(r.FinancialData, "returnOnEquity")), GrossMargin: yahooPercentValue(rawNum(r.FinancialData, "grossMargins")), OperatingMargin: yahooPercentValue(rawNum(r.FinancialData, "operatingMargins")), RevenueGrowth: yahooPercentValue(rawNum(r.FinancialData, "revenueGrowth")), UpdatedAt: time.Now().UnixMilli(), Source: "yfinance-recovery"}
		e.mu.Lock()
		e.fundamentals[sym] = f
		e.mu.Unlock()
		updated++
		e.recordProviderSuccess("yfinance")
	}
	if updated > 0 {
		e.mu.Lock()
		e.lastUpdated["fundamentals"] = time.Now().UnixMilli()
		e.health["fundamentals"] = "healthy · yfinance recovery"
		e.mu.Unlock()
	}
	return updated
}

// SEC v15 enrichment: include OTHER transactions and recent rich transaction detail.
func testMarketauxProvider(ctx context.Context, key string) ProviderTestResult {
	r := ProviderTestResult{Provider: "marketaux", CheckedAt: time.Now().UTC().Format(time.RFC3339)}
	key = strings.TrimSpace(key)
	if key == "" {
		r.Status = "missing"
		r.Message = "Enter a Marketaux API key."
		return r
	}
	var out struct {
		Data []any `json:"data"`
		Meta any   `json:"meta"`
	}
	raw := "https://api.marketaux.com/v1/news/all?api_token=" + url.QueryEscape(key) + "&language=en&limit=1"
	if err := getJSON(ctx, &http.Client{Timeout: 12 * time.Second}, raw, nil, &out); err != nil {
		r.Status = "failed"
		r.Message = err.Error()
		return r
	}
	r.OK = true
	r.Status = "connected"
	r.Message = "Marketaux news fallback available."
	r.Details = []string{"News fallback only; Finnhub remains primary."}
	return r
}

func (e *Engine) refreshYahooEarnings(ctx context.Context, symbols []string) int {
	if !e.providerAllowed("yfinance") {
		return 0
	}
	items := []EarningsItem{}
	for _, sym := range uniqueSymbols(symbols) {
		raw := "https://query1.finance.yahoo.com/v10/finance/quoteSummary/" + url.PathEscape(sym) + "?modules=calendarEvents"
		var p struct {
			QuoteSummary struct {
				Result []struct {
					CalendarEvents struct {
						Earnings struct {
							EarningsDate []struct {
								Raw int64  `json:"raw"`
								Fmt string `json:"fmt"`
							} `json:"earningsDate"`
							EarningsAverage map[string]any `json:"earningsAverage"`
							RevenueAverage  map[string]any `json:"revenueAverage"`
						} `json:"earnings"`
					} `json:"calendarEvents"`
				} `json:"result"`
			} `json:"quoteSummary"`
		}
		if err := getJSON(ctx, &http.Client{Timeout: 12 * time.Second}, raw, map[string]string{"User-Agent": "Mozilla/5.0 DE.PULSE/15.0"}, &p); err != nil {
			e.recordProviderFailure("yfinance", err)
			continue
		}
		if len(p.QuoteSummary.Result) == 0 {
			continue
		}
		er := p.QuoteSummary.Result[0].CalendarEvents.Earnings
		for _, d := range er.EarningsDate {
			if d.Raw <= 0 {
				continue
			}
			t := time.Unix(d.Raw, 0)
			est := func(m map[string]any) *float64 {
				v := 0.0
				if x, ok := m["raw"]; ok {
					v = toFloat(x)
				}
				if v == 0 {
					return nil
				}
				return &v
			}
			items = append(items, EarningsItem{Symbol: sym, Date: t.Format("2006-01-02"), Hour: "TBD", EPSEstimate: est(er.EarningsAverage), RevenueEstimate: est(er.RevenueAverage)})
		}
		e.recordProviderSuccess("yfinance")
	}
	if len(items) == 0 {
		return 0
	}
	sort.Slice(items, func(i, j int) bool { return items[i].Date < items[j].Date })
	e.mu.Lock()
	e.earnings = items
	e.lastUpdated["earnings"] = time.Now().UnixMilli()
	e.health["earnings"] = "healthy · yfinance recovery"
	e.mu.Unlock()
	e.app.broadcastEarnings(clone(items))
	return len(items)
}
