package main

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"
)

// Research refreshes are target-scoped so selecting one ticker never fans out
// expensive REST calls across the entire watchlist universe. The canonical
// stores are still updated, but only the selected symbol is reconciled here.
func (e *Engine) refreshResearchFundamentals(ctx context.Context, key, symbol string) bool {
	symbol = normalizeSymbol(symbol)
	if symbol == "" {
		return false
	}
	e.mu.RLock()
	priorGlobalAt := e.lastUpdated["fundamentals"]
	priorGlobalHealth := e.health["fundamentals"]
	e.mu.RUnlock()
	active, ok := e.executeProviderRoute(ctx, "Fundamentals", map[string]providerRouteAttempt{
		"Finnhub": func(ctx context.Context) bool {
			if strings.TrimSpace(key) == "" || e.highImpactModeActive() {
				return false
			}
			var payload struct {
				Metric map[string]any `json:"metric"`
			}
			if err := e.finnhubJSONForSymbol(ctx, key, symbol, "/stock/metric?symbol="+url.QueryEscape(symbol)+"&metric=all", &payload); err != nil {
				e.recordProviderFailure("Finnhub", err)
				return false
			}
			m := payload.Metric
			if len(m) == 0 {
				return false
			}
			f := FundamentalSnapshot{Symbol: symbol, MarketCap: toFloat(m["marketCapitalization"]), PERatio: toFloat(m["peTTM"]), ForwardPERatio: toFloat(m["forwardPEAnnual"]), PSRatio: toFloat(m["psTTM"]), PEGRatio: toFloat(m["pegRatio"]), RevenueGrowth: toFloat(m["revenueGrowthTTMYoy"]), EPSGrowth: toFloat(m["epsGrowthTTMYoy"]), GrossMargin: toFloat(m["grossMarginTTM"]), OperatingMargin: toFloat(m["operatingMarginTTM"]), ROE: toFloat(m["roeTTM"]), NetMargin: toFloat(m["netProfitMarginTTM"]), DebtToEquity: toFloat(m["totalDebt/totalEquityAnnual"]), CurrentRatio: toFloat(m["currentRatioAnnual"]), FreeCashFlow: toFloat(m["freeCashFlowTTM"]), DividendYield: toFloat(m["dividendYieldIndicatedAnnual"]), FiftyTwoWeekHigh: toFloat(m["52WeekHigh"]), FiftyTwoWeekLow: toFloat(m["52WeekLow"]), UpdatedAt: time.Now().UnixMilli(), Source: "finnhub"}
			e.mu.Lock()
			e.fundamentals[symbol] = f
			e.mu.Unlock()
			e.recordProviderSuccess("Finnhub")
			return true
		},
		"SEC":      func(ctx context.Context) bool { return e.refreshSECFundamentals(ctx, []string{symbol}) > 0 },
		"yfinance": func(ctx context.Context) bool { return e.refreshYahooFundamentals(ctx, []string{symbol}) > 0 },
	})
	e.mu.Lock()

	e.lastUpdated["fundamentals"] = priorGlobalAt
	e.health["fundamentals"] = priorGlobalHealth
	if ok {
		now := time.Now().UnixMilli()
		e.lastUpdated["research-fundamentals:"+symbol] = now
		e.health["research-fundamentals:"+symbol] = "healthy · " + active
	}
	e.mu.Unlock()
	return ok
}

func (e *Engine) refreshResearchEarnings(ctx context.Context, key, symbol string) bool {
	symbol = normalizeSymbol(symbol)
	if symbol == "" {
		return false
	}
	checked := false
	active, ok := e.executeProviderRoute(ctx, "Earnings", map[string]providerRouteAttempt{
		"Finnhub": func(ctx context.Context) bool {
			if strings.TrimSpace(key) == "" {
				return false
			}
			from := time.Now().Add(-365 * 24 * time.Hour).Format("2006-01-02")
			to := time.Now().Add(120 * 24 * time.Hour).Format("2006-01-02")
			var payload struct {
				EarningsCalendar []EarningsItem `json:"earningsCalendar"`
			}
			endpoint := fmt.Sprintf("/calendar/earnings?from=%s&to=%s&symbol=%s", from, to, url.QueryEscape(symbol))
			if err := e.finnhubJSONForSymbol(ctx, key, symbol, endpoint, &payload); err != nil {
				e.recordProviderFailure("Finnhub", err)
				return false
			}
			checked = true
			e.recordProviderSuccess("Finnhub")
			e.mu.Lock()
			merged := make([]EarningsItem, 0, len(e.earnings)+len(payload.EarningsCalendar))
			for _, x := range e.earnings {
				if !strings.EqualFold(x.Symbol, symbol) {
					merged = append(merged, x)
				}
			}
			merged = append(merged, payload.EarningsCalendar...)
			sort.Slice(merged, func(i, j int) bool { return merged[i].Date < merged[j].Date })
			e.earnings = merged
			e.mu.Unlock()
			return true
		},
		"yfinance": func(ctx context.Context) bool {
			e.mu.RLock()
			before := clone(e.earnings)
			priorAt := e.lastUpdated["earnings"]
			priorHealth := e.health["earnings"]
			e.mu.RUnlock()
			n := e.refreshYahooEarnings(ctx, []string{symbol})
			if n <= 0 {
				return false
			}
			e.mu.Lock()
			selected := clone(e.earnings)
			merged := make([]EarningsItem, 0, len(before)+len(selected))
			for _, x := range before {
				if !strings.EqualFold(x.Symbol, symbol) {
					merged = append(merged, x)
				}
			}
			merged = append(merged, selected...)
			sort.Slice(merged, func(i, j int) bool { return merged[i].Date < merged[j].Date })
			e.earnings = merged
			e.lastUpdated["earnings"] = priorAt
			e.health["earnings"] = priorHealth
			e.mu.Unlock()
			checked = true
			return true
		},
	})
	if ok || checked {
		now := time.Now().UnixMilli()
		e.mu.Lock()
		e.lastUpdated["research-earnings:"+symbol] = now
		e.health["research-earnings:"+symbol] = "healthy · " + defaultString(active, "provider check")
		e.mu.Unlock()
		e.enrichEarningsGuidanceFromEvidence()
		e.evaluateCatalystWatch(time.Now())
		return true
	}
	return false
}

func (e *Engine) refreshResearchNews(ctx context.Context, key, symbol string) bool {
	symbol = normalizeSymbol(symbol)
	if symbol == "" {
		return false
	}
	active, ok := e.executeProviderRoute(ctx, "News", map[string]providerRouteAttempt{
		"Finnhub": func(ctx context.Context) bool {
			if strings.TrimSpace(key) == "" {
				return false
			}
			from := time.Now().Add(-7 * 24 * time.Hour).Format("2006-01-02")
			to := time.Now().Format("2006-01-02")
			var items []NewsItem
			endpoint := fmt.Sprintf("/company-news?symbol=%s&from=%s&to=%s", url.QueryEscape(symbol), from, to)
			if err := e.finnhubJSONForSymbol(ctx, key, symbol, endpoint, &items); err != nil {
				e.recordProviderFailure("Finnhub", err)
				return false
			}
			for i := range items {
				items[i].Symbols = uniqueSymbols(append([]string{symbol}, strings.Split(items[i].Related, ",")...))
			}
			e.mu.Lock()
			kept := make([]NewsItem, 0, len(e.news)+len(items))
			for _, n := range e.news {
				if !containsSymbol(n.Symbols, symbol) {
					kept = append(kept, n)
				}
			}
			kept = append(kept, items...)
			e.news = dedupeNews(kept)
			if len(e.news) > 150 {
				e.news = e.news[:150]
			}
			e.mu.Unlock()
			e.recordProviderSuccess("Finnhub")
			return true
		},
		"Marketaux": func(ctx context.Context) bool {
			e.app.mu.RLock()
			k := strings.TrimSpace(e.app.secrets.Marketaux)
			e.app.mu.RUnlock()
			if k == "" || !e.providerAllowed("Marketaux") {
				return false
			}
			raw := "https://api.marketaux.com/v1/news/all?api_token=" + url.QueryEscape(k) + "&language=en&limit=3&symbols=" + url.QueryEscape(symbol)
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
				syms := []string{symbol}
				for _, ent := range x.Entities {
					syms = append(syms, ent.Symbol)
				}
				items = append(items, NewsItem{ID: x.UUID, Datetime: publishedAt, Headline: x.Title, Summary: x.Description, Source: "Marketaux · " + x.Source, URL: x.URL, Symbols: uniqueSymbols(syms), Scope: "watchlist"})
			}
			e.mu.Lock()
			kept := make([]NewsItem, 0, len(e.news)+len(items))
			for _, n := range e.news {
				if !containsSymbol(n.Symbols, symbol) {
					kept = append(kept, n)
				}
			}
			kept = append(kept, items...)
			e.news = dedupeNews(kept)
			if len(e.news) > 150 {
				e.news = e.news[:150]
			}
			e.mu.Unlock()
			e.recordProviderSuccess("Marketaux")
			return true
		},
	})
	if ok {
		now := time.Now().UnixMilli()
		e.mu.Lock()
		e.lastUpdated["research-news:"+symbol] = now
		e.health["research-news:"+symbol] = "healthy · " + active
		e.mu.Unlock()
		e.enrichEarningsGuidanceFromEvidence()
		e.evaluateCatalystWatch(time.Now())
	}
	return ok
}

func (e *Engine) refreshFundamentals(ctx context.Context, key string) {
	e.app.mu.RLock()
	symbols := analysisSymbolsFromState(e.app.processingStateLocked())
	e.app.mu.RUnlock()
	active, ok := e.executeProviderRoute(ctx, "Fundamentals", map[string]providerRouteAttempt{
		"Finnhub":  func(ctx context.Context) bool { return e.refreshFinnhubFundamentals(ctx, key) },
		"SEC":      func(ctx context.Context) bool { return e.refreshSECFundamentals(ctx, symbols) > 0 },
		"yfinance": func(ctx context.Context) bool { return e.refreshYahooFundamentals(ctx, symbols) > 0 },
	})
	if !ok {
		e.setHealth("fundamentals", "unavailable · provider route exhausted")
		return
	}
	e.setHealth("fundamentals-route", "active · "+active)
}

func (e *Engine) refreshFinnhubFundamentals(ctx context.Context, key string) bool {
	if strings.TrimSpace(key) == "" {
		return false
	}
	if e.highImpactModeActive() {
		e.setHealth("fundamentals", "deferred · high-impact event mode")
		return false
	}
	e.setHealth("fundamentals", "loading")
	e.app.mu.RLock()
	symbols := analysisSymbolsFromState(e.app.processingStateLocked())
	e.app.mu.RUnlock()
	if len(symbols) > 30 {
		symbols = symbols[:30]
	}
	sem := make(chan struct{}, 4)
	var wg sync.WaitGroup
	var mu sync.Mutex
	updated := 0
	for _, sym := range symbols {
		sym := sym
		wg.Add(1)
		go func() {
			defer wg.Done()
			sem <- struct{}{}
			defer func() { <-sem }()
			var payload struct {
				Metric map[string]any `json:"metric"`
			}
			if err := e.finnhubJSONForSymbol(ctx, key, sym, "/stock/metric?symbol="+url.QueryEscape(sym)+"&metric=all", &payload); err != nil {
				e.recordProviderFailure("Finnhub", err)
				return
			}
			e.recordProviderSuccess("Finnhub")
			m := payload.Metric
			f := FundamentalSnapshot{Symbol: sym, MarketCap: toFloat(m["marketCapitalization"]), PERatio: toFloat(m["peTTM"]), ForwardPERatio: toFloat(m["forwardPEAnnual"]), PSRatio: toFloat(m["psTTM"]), PEGRatio: toFloat(m["pegRatio"]), RevenueGrowth: toFloat(m["revenueGrowthTTMYoy"]), EPSGrowth: toFloat(m["epsGrowthTTMYoy"]), GrossMargin: toFloat(m["grossMarginTTM"]), OperatingMargin: toFloat(m["operatingMarginTTM"]), ROE: toFloat(m["roeTTM"]), NetMargin: toFloat(m["netProfitMarginTTM"]), DebtToEquity: toFloat(m["totalDebt/totalEquityAnnual"]), CurrentRatio: toFloat(m["currentRatioAnnual"]), FreeCashFlow: toFloat(m["freeCashFlowTTM"]), DividendYield: toFloat(m["dividendYieldIndicatedAnnual"]), FiftyTwoWeekHigh: toFloat(m["52WeekHigh"]), FiftyTwoWeekLow: toFloat(m["52WeekLow"]), UpdatedAt: time.Now().UnixMilli(), Source: "finnhub"}
			mu.Lock()
			e.mu.Lock()
			e.fundamentals[sym] = f
			e.lastUpdated["fundamentals"] = time.Now().UnixMilli()
			e.mu.Unlock()
			updated++
			mu.Unlock()
		}()
	}
	wg.Wait()
	if updated > 0 {
		e.setHealth("fundamentals", "healthy · Finnhub")
		_ = e.saveCache()
		e.app.broadcastRuntime()
		return true
	}
	e.setHealth("fundamentals", "degraded · Finnhub primary returned no fundamentals")
	return false
}

func (e *Engine) refreshNews(ctx context.Context, key string) {
	active, ok := e.executeProviderRoute(ctx, "News", map[string]providerRouteAttempt{
		"Finnhub":   func(ctx context.Context) bool { return e.refreshFinnhubNews(ctx, key) },
		"Marketaux": func(ctx context.Context) bool { return e.refreshMarketauxNews(ctx) },
	})
	if !ok {
		e.setHealth("news", "unavailable · provider route exhausted")
		return
	}
	e.setHealth("news-route", "active · "+active)
	e.mu.RLock()
	out := clone(e.news)
	e.mu.RUnlock()
	e.app.broadcastNews(out)
	// v16.2 Event Intelligence is derived from canonical stores. Publish a full
	// runtime snapshot after the canonical News route commits so event context
	// updates immediately without introducing another fetch loop.
	e.app.broadcastRuntime()
}

func (e *Engine) refreshFinnhubNews(ctx context.Context, key string) bool {
	e.setHealth("news", "loading")
	if strings.TrimSpace(key) == "" {
		return false
	}
	var all []NewsItem
	var mu sync.Mutex
	successful := 0
	var general []NewsItem
	var generalErr error
	if err := e.finnhubJSON(ctx, key, "/news?category=general&minId=0", &general); err == nil {
		successful++
		for i := range general {
			general[i].Scope = "general"
		}
		if len(general) > 40 {
			general = general[:40]
		}
		all = append(all, general...)
	} else {
		generalErr = err
	}
	e.app.mu.RLock()
	symbols := analysisSymbolsFromState(e.app.processingStateLocked())
	e.app.mu.RUnlock()
	if len(symbols) > 20 {
		symbols = symbols[:20]
	}
	sem := make(chan struct{}, 4)
	var wg sync.WaitGroup
	from := time.Now().Add(-72 * time.Hour).Format("2006-01-02")
	to := time.Now().Format("2006-01-02")
	for _, symbol := range symbols {
		symbol := symbol
		wg.Add(1)
		go func() {
			defer wg.Done()
			sem <- struct{}{}
			defer func() { <-sem }()
			var items []NewsItem
			endpoint := fmt.Sprintf("/company-news?symbol=%s&from=%s&to=%s", url.QueryEscape(symbol), from, to)
			if e.finnhubJSONForSymbol(ctx, key, symbol, endpoint, &items) == nil {
				mu.Lock()
				successful++
				mu.Unlock()
				if len(items) > 15 {
					items = items[:15]
				}
				for i := range items {
					items[i].Symbols = uniqueSymbols(append([]string{symbol}, strings.Split(items[i].Related, ",")...))
				}
				mu.Lock()
				all = append(all, items...)
				mu.Unlock()
			}
		}()
	}
	wg.Wait()
	if successful == 0 {
		if generalErr == nil {
			generalErr = fmt.Errorf("Finnhub returned no usable news responses")
		}
		e.recordProviderFailure("Finnhub", generalErr)
		e.setHealth("news", "degraded · Finnhub primary returned no news")
		return false
	}

	e.recordProviderSuccess("Finnhub")
	all = dedupeNews(all)
	if len(all) > 150 {
		all = all[:150]
	}
	e.mu.Lock()
	e.news = all
	e.health["news"] = "healthy · Finnhub"
	e.lastUpdated["news"] = time.Now().UnixMilli()
	e.mu.Unlock()
	e.enrichEarningsGuidanceFromEvidence()
	e.evaluateCatalystWatch(time.Now())
	e.app.broadcastNews(clone(all))
	return true
}

func dedupeNews(items []NewsItem) []NewsItem {
	seen := map[string]bool{}
	out := make([]NewsItem, 0, len(items))
	for _, item := range items {
		key := fmt.Sprintf("%v", item.ID)
		if key == "<nil>" || key == "" {
			key = strings.ToLower(item.URL)
		}
		if key == "" {
			key = strings.ToLower(item.Headline) + "|" + strconv.FormatInt(item.Datetime, 10)
		}
		if seen[key] {
			continue
		}
		seen[key] = true
		out = append(out, item)
	}
	sort.Slice(out, func(i, j int) bool { return out[i].Datetime > out[j].Datetime })
	return out
}

func reportedEarningsFingerprint(items []EarningsItem) string {
	parts := make([]string, 0, len(items))
	for _, it := range items {
		if it.EPSActual == nil && it.RevenueActual == nil {
			continue
		}
		eps, rev := "", ""
		if it.EPSActual != nil {
			eps = fmt.Sprintf("%.8f", *it.EPSActual)
		}
		if it.RevenueActual != nil {
			rev = fmt.Sprintf("%.2f", *it.RevenueActual)
		}
		parts = append(parts, fmt.Sprintf("%s|%s|%d|%d|%s|%s", strings.ToUpper(it.Symbol), it.Date, it.Quarter, it.Year, eps, rev))
	}
	sort.Strings(parts)
	return strings.Join(parts, ";")
}

func financialFilingFingerprint(items []FilingItem) string {
	parts := make([]string, 0, len(items))
	for _, it := range items {
		base := strings.ToUpper(strings.TrimSpace(strings.TrimSuffix(it.Form, "/A")))
		if base != "10-Q" && base != "10-K" && base != "20-F" && base != "6-K" {
			continue
		}
		parts = append(parts, strings.ToUpper(it.Symbol)+"|"+base+"|"+it.FiledAt+"|"+it.ReportDate)
	}
	sort.Strings(parts)
	return strings.Join(parts, ";")
}

func (e *Engine) refreshEarnings(ctx context.Context, key string) {
	symbols := e.trackedSymbols()
	active, ok := e.executeProviderRoute(ctx, "Earnings", map[string]providerRouteAttempt{
		"Finnhub":  func(ctx context.Context) bool { return e.refreshFinnhubEarnings(ctx, key) },
		"yfinance": func(ctx context.Context) bool { return e.refreshYahooEarnings(ctx, symbols) > 0 },
	})
	if !ok {
		e.setHealth("earnings", "unavailable · provider route exhausted")
		return
	}
	e.setHealth("earnings-route", "active · "+active)
	e.app.broadcastRuntime()
}

func (e *Engine) refreshFinnhubEarnings(ctx context.Context, key string) bool {
	e.setHealth("earnings", "loading")
	if strings.TrimSpace(key) == "" {
		return false
	}
	e.mu.RLock()
	oldReportedFingerprint := reportedEarningsFingerprint(e.earnings)
	e.mu.RUnlock()
	symbols := e.trackedSymbols()
	from := time.Now().Add(-21 * 24 * time.Hour).Format("2006-01-02")
	to := time.Now().Add(60 * 24 * time.Hour).Format("2006-01-02")
	var all []EarningsItem
	var mu sync.Mutex
	successful := 0
	sem := make(chan struct{}, 4)
	var wg sync.WaitGroup
	for _, symbol := range symbols {
		symbol := symbol
		wg.Add(1)
		go func() {
			defer wg.Done()
			sem <- struct{}{}
			defer func() { <-sem }()
			var payload struct {
				EarningsCalendar []EarningsItem `json:"earningsCalendar"`
			}
			endpoint := fmt.Sprintf("/calendar/earnings?from=%s&to=%s&symbol=%s", from, to, url.QueryEscape(symbol))
			if e.finnhubJSONForSymbol(ctx, key, symbol, endpoint, &payload) == nil {
				mu.Lock()
				successful++
				all = append(all, payload.EarningsCalendar...)
				mu.Unlock()
			}
		}()
	}
	wg.Wait()
	if successful == 0 && len(symbols) > 0 {
		e.recordProviderFailure("Finnhub", fmt.Errorf("Finnhub returned no usable earnings responses"))
		e.setHealth("earnings", "degraded · Finnhub primary returned no earnings")
		return false
	}
	if successful > 0 {
		e.recordProviderSuccess("Finnhub")
	}
	seen := map[string]bool{}
	var out []EarningsItem
	for _, it := range all {
		k := fmt.Sprintf("%s|%s|%d|%d", it.Symbol, it.Date, it.Quarter, it.Year)
		if seen[k] {
			continue
		}
		seen[k] = true
		out = append(out, it)
	}
	sort.Slice(out, func(i, j int) bool { return out[i].Date < out[j].Date })
	if len(out) > 200 {
		out = out[:200]
	}
	newReportedFingerprint := reportedEarningsFingerprint(out)
	e.mu.Lock()
	e.earnings = out
	e.health["earnings"] = "healthy · Finnhub"
	e.lastUpdated["earnings"] = time.Now().UnixMilli()
	if oldReportedFingerprint != "" && newReportedFingerprint != "" && oldReportedFingerprint != newReportedFingerprint {
		e.lastUpdated["fundamentals"] = 0
		e.health["fundamentals"] = "stale · new reported earnings"
	}
	e.mu.Unlock()
	e.app.broadcastEarnings(clone(out))
	e.enrichEarningsGuidanceFromEvidence()
	e.evaluateCatalystWatch(time.Now())
	return true
}
