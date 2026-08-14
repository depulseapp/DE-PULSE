package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

func (e *Engine) refreshSnapshots(ctx context.Context, key string) {

	e.app.mu.RLock()
	twelveKey := strings.TrimSpace(e.app.secrets.TwelveData)
	e.app.mu.RUnlock()
	tracked := e.trackedSymbols()
	candidates := make([]string, 0, len(tracked))
	avoided := int64(0)
	now := time.Now().UnixMilli()
	for _, symbol := range tracked {
		if symbol == "VIX" {
			continue
		}
		e.mu.RLock()
		q := e.quotes[symbol]
		e.mu.RUnlock()
		stamp := q.ProviderTimestamp
		if stamp == 0 {
			stamp = q.UpdatedAt
		}
		age := int64(1 << 62)
		if stamp > 0 {
			age = now - stamp
		}

		if q.Price > 0 && age <= 90_000 {
			avoided++
			continue
		}
		candidates = append(candidates, symbol)
	}
	if avoided > 0 {
		e.mu.Lock()
		e.providerCallsAvoided += avoided
		e.mu.Unlock()
	}
	updatedFinnhub, updatedTwelve, twelveAttempts := 0, 0, 0
	lastFinnhubCall := time.Time{}
	for _, symbol := range candidates {
		select {
		case <-ctx.Done():
			return
		default:
		}
		_, _ = e.executeProviderRoute(ctx, "US Live Equities", map[string]providerRouteAttempt{
			"Alpaca": func(context.Context) bool {

				return false
			},
			"Finnhub": func(ctx context.Context) bool {
				if strings.TrimSpace(key) == "" {
					return false
				}
				if !lastFinnhubCall.IsZero() {
					if wait := time.Until(lastFinnhubCall.Add(1050 * time.Millisecond)); wait > 0 && !sleepContext(ctx, wait) {
						return false
					}
				}
				lastFinnhubCall = time.Now()
				var q finnhubQuoteResponse
				err := e.finnhubJSON(ctx, key, "/quote?symbol="+url.QueryEscape(symbol), &q)
				if err != nil {
					e.recordProviderFailure("Finnhub", err)
					return false
				}
				if q.Current <= 0 {
					return false
				}
				e.mergeFinnhubSnapshot(symbol, q)
				e.recordProviderSuccess("Finnhub")
				updatedFinnhub++
				return true
			},
			"Twelve Data": func(ctx context.Context) bool {
				if twelveKey == "" || twelveAttempts >= 4 {
					return false
				}
				twelveAttempts++
				q, err := tdQuote(ctx, twelveKey, symbol)
				if err != nil {
					e.recordProviderFailure("Twelve Data", err)
					return false
				}
				px, _ := strconv.ParseFloat(q.Close, 64)
				if px <= 0 {
					return false
				}
				op, _ := strconv.ParseFloat(q.Open, 64)
				hi, _ := strconv.ParseFloat(q.High, 64)
				lo, _ := strconv.ParseFloat(q.Low, 64)
				prev, _ := strconv.ParseFloat(q.PreviousClose, 64)
				ts := normalizeObservationMs(q.Timestamp)
				e.updateQuote(symbol, Quote{Price: px, Open: op, High: hi, Low: lo, PreviousClose: prev, ProviderTimestamp: ts, FeedType: "rest-recovery", DataState: "snapshot"}, "twelvedata-equity-fallback:"+symbol)
				e.recordProviderSuccess("Twelve Data")
				updatedTwelve++
				return true
			},
		})
	}
	switch {
	case len(candidates) == 0:
		e.setHealth("quotes-rest", "idle · live/primary observations current")
	case updatedFinnhub+updatedTwelve > 0:
		e.setHealth("quotes-rest", fmt.Sprintf("recovery · Finnhub %d · Twelve Data %d", updatedFinnhub, updatedTwelve))
	default:
		e.setHealth("quotes-rest", "degraded · no recovery source updated stale symbols")
	}
}

func (e *Engine) mergeFinnhubSnapshot(symbol string, snap finnhubQuoteResponse) {
	symbol = normalizeSymbol(symbol)
	if symbol == "" || snap.Current <= 0 {
		return
	}
	now := time.Now().UnixMilli()
	e.mu.Lock()
	prev := e.quotes[symbol]
	q := prev
	recentStream := q.Source == "finnhub-websocket" && q.ProviderTimestamp > 0 && now-q.ProviderTimestamp <= 120_000
	if !recentStream || q.Price <= 0 {
		q.Symbol = symbol
		q.Price = snap.Current
		q.Source = "finnhub-rest"
		q.FeedType = "rest"
		q.DataState = "snapshot"
		q.ProviderTimestamp = snap.Timestamp * 1000
		q.UpdatedAt = now
	}
	if snap.Open > 0 {
		q.Open = snap.Open
	}
	if snap.High > 0 {
		q.High = snap.High
	}
	if snap.Low > 0 {
		q.Low = snap.Low
	}
	if snap.Previous > 0 {
		q.PreviousClose = snap.Previous
	}
	if q.PreviousClose > 0 && q.Price > 0 {
		q.Change = q.Price - q.PreviousClose
		q.ChangePercent = q.Change / q.PreviousClose * 100
	}
	e.quotes[symbol] = q
	e.lastUpdated["quotes-rest"] = now
	e.mu.Unlock()
	e.recordProviderQuoteObservation(symbol, q)
	e.evaluateRapidMoveObservation(symbol, q)
	e.propagateCanonicalQuoteChange(symbol, prev, q)
	e.app.broadcastSymbolEvent(q.Symbol, map[string]any{"type": "quote", "quote": q})
}

type finnhubQuoteResponse struct {
	Current   float64 `json:"c"`
	Change    float64 `json:"d"`
	Percent   float64 `json:"dp"`
	High      float64 `json:"h"`
	Low       float64 `json:"l"`
	Open      float64 `json:"o"`
	Previous  float64 `json:"pc"`
	Timestamp int64   `json:"t"`
}

func (e *Engine) refreshVIXSnapshot(ctx context.Context, key string) bool {

	_ = key
	e.app.mu.RLock()
	twelveKey := strings.TrimSpace(e.app.secrets.TwelveData)
	e.app.mu.RUnlock()
	active, ok := e.executeProviderRoute(ctx, "VIX / Indices", map[string]providerRouteAttempt{
		"Twelve Data": func(ctx context.Context) bool {
			q, providerSymbol, err := tdVIXQuote(ctx, twelveKey)
			if err != nil {
				e.recordProviderFailure("Twelve Data", err)
				return false
			}
			px, _ := strconv.ParseFloat(q.Close, 64)
			if px <= 0 {
				e.recordProviderFailure("Twelve Data", fmt.Errorf("invalid VIX price payload %q", q.Close))
				return false
			}
			prev, _ := strconv.ParseFloat(q.PreviousClose, 64)
			op, _ := strconv.ParseFloat(q.Open, 64)
			hi, _ := strconv.ParseFloat(q.High, 64)
			lo, _ := strconv.ParseFloat(q.Low, 64)
			ts := normalizeObservationMs(q.Timestamp)
			e.updateQuote("VIX", Quote{Price: px, Open: op, High: hi, Low: lo, PreviousClose: prev, ProviderTimestamp: ts, FeedType: "index-direct", DataState: "snapshot"}, "twelvedata-vix:"+providerSymbol)
			e.recordProviderSuccess("Twelve Data")
			e.setHealth("vix", "Direct VIX index · Twelve Data · "+providerSymbol)
			_ = e.refreshCboeVIXOfficial(ctx, false)
			return true
		},
		"yfinance": func(ctx context.Context) bool {
			ok := e.refreshYahooVIX(ctx)
			if ok {
				_ = e.refreshCboeVIXOfficial(ctx, false)
			}
			return ok
		},
		"CBOE": func(ctx context.Context) bool {
			ok := e.refreshCboeVIXOfficial(ctx, true)
			if ok {
				e.recordProviderSuccess("CBOE")
			}
			return ok
		},
	})
	if ok {
		e.setHealth("vix-route", "active · "+active)
		return true
	}
	e.mu.RLock()
	cached := e.quotes["VIX"].Price > 0
	e.mu.RUnlock()
	if cached {
		e.setHealth("vix", "Cached VIX · current providers unavailable")
	} else {
		e.setHealth("vix", "VIX unavailable · no validated source")
	}
	return false
}

func (e *Engine) refreshCboeVIXOfficial(ctx context.Context, updateSpot bool) bool {
	bars, err := fetchCboeVIXHistory(ctx)
	if err != nil {
		e.recordProviderFailure("CBOE", err)
		return false
	}
	if len(bars) == 0 {
		e.recordProviderFailure("CBOE", errors.New("CBOE VIX history returned no usable rows"))
		return false
	}
	e.recordProviderSuccess("CBOE")

	if len(bars) > 260 {
		bars = bars[len(bars)-260:]
	}
	e.mu.Lock()
	if e.bars["VIX"] == nil {
		e.bars["VIX"] = map[string][]Bar{}
	}
	e.bars["VIX"]["daily"] = append([]Bar(nil), bars...)
	e.lastUpdated["vix-history"] = time.Now().UnixMilli()
	e.mu.Unlock()
	if !updateSpot {
		return true
	}
	last := bars[len(bars)-1]
	prev := 0.0
	if len(bars) > 1 {
		prev = bars[len(bars)-2].C
	}
	e.updateQuote("VIX", Quote{Price: last.C, Open: last.O, High: last.H, Low: last.L, PreviousClose: prev, ProviderTimestamp: last.T * 1000, FeedType: "official-close", DataState: "snapshot"}, "cboe-vix-official-close")
	e.setHealth("vix", "Official VIX close · Cboe · daily")
	return true
}

func (e *Engine) finnhubJSONTier(ctx context.Context, key, endpoint string, tier WorkTier, out any) (err error) {
	tier = normalizeWorkTier(tier)
	if ok, reason := e.providerTelemetry.Allow("Finnhub", tier); !ok {
		return fmt.Errorf("%s Finnhub request deferred: %s", workTierLabel(tier), reason)
	}
	// Preserve Finnhub's existing paced request contract independently of the
	// shared provider concurrency budget.
	e.finnhubRateMu.Lock()
	wait := finnhubMinRequestInterval - time.Since(e.finnhubLastRequest)
	if wait > 0 {
		select {
		case <-ctx.Done():
			e.finnhubRateMu.Unlock()
			return ctx.Err()
		case <-time.After(wait):
		}
	}
	e.finnhubLastRequest = time.Now()
	e.finnhubRateMu.Unlock()

	release, ok := e.workload.AcquireTier(ctx, "provider-rest", tier)
	if !ok {
		if ctx.Err() != nil {
			return ctx.Err()
		}
		return fmt.Errorf("%s Finnhub request deferred by bounded provider capacity", workTierLabel(tier))
	}
	defer release()
	done := e.providerTelemetry.begin("Finnhub")
	defer func() { done(err) }()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, strings.TrimRight(finnhubAPIBaseURL, "/")+endpoint, nil)
	if err != nil {
		return err
	}
	req.Header.Set("X-Finnhub-Token", key)
	req.Header.Set("Accept", "application/json")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		body, _ := io.ReadAll(io.LimitReader(resp.Body, 300))
		return fmt.Errorf("HTTP %d: %s", resp.StatusCode, strings.TrimSpace(string(body)))
	}
	return json.NewDecoder(resp.Body).Decode(out)
}

func (e *Engine) finnhubJSONForSymbol(ctx context.Context, key, symbol, endpoint string, out any) error {
	return e.finnhubJSONTier(ctx, key, endpoint, e.workTierForSymbol(symbol), out)
}

func (e *Engine) finnhubJSON(ctx context.Context, key, endpoint string, out any) error {
	return e.finnhubJSONTier(ctx, key, endpoint, WorkTierUserActionable, out)
}

type alpacaLiveSnapshot struct {
	LatestTrade struct {
		Price float64 `json:"p"`
		Time  string  `json:"t"`
	} `json:"latestTrade"`
	LatestQuote struct {
		Ask     float64 `json:"ap"`
		Bid     float64 `json:"bp"`
		AskSize float64 `json:"as"`
		BidSize float64 `json:"bs"`
		Time    string  `json:"t"`
	} `json:"latestQuote"`
	MinuteBar struct {
		Open   float64 `json:"o"`
		High   float64 `json:"h"`
		Low    float64 `json:"l"`
		Close  float64 `json:"c"`
		Volume float64 `json:"v"`
		Time   string  `json:"t"`
	} `json:"minuteBar"`
	DailyBar struct {
		Open   float64 `json:"o"`
		High   float64 `json:"h"`
		Low    float64 `json:"l"`
		Close  float64 `json:"c"`
		Volume float64 `json:"v"`
		Time   string  `json:"t"`
	} `json:"dailyBar"`
	PrevDailyBar struct {
		Open   float64 `json:"o"`
		High   float64 `json:"h"`
		Low    float64 `json:"l"`
		Close  float64 `json:"c"`
		Volume float64 `json:"v"`
		Time   string  `json:"t"`
	} `json:"prevDailyBar"`
}

func providerTimeMillis(raw string) int64 {
	if strings.TrimSpace(raw) == "" {
		return 0
	}
	t, err := time.Parse(time.RFC3339Nano, raw)
	if err != nil {
		return 0
	}
	return t.UnixMilli()
}

// Alpaca's free overnight feed provides real-time indicative quotes while
// latest trades may be delayed. For overnight, prefer the quote midpoint so
// DE.PULSE does not present a delayed trade as the live overnight value.
func alpacaSnapshotPrice(snap alpacaLiveSnapshot, feed, session string) (float64, int64, string) {
	quoteMid := func() (float64, int64, bool) {
		if snap.LatestQuote.Ask > 0 && snap.LatestQuote.Bid > 0 {
			return (snap.LatestQuote.Ask + snap.LatestQuote.Bid) / 2, providerTimeMillis(snap.LatestQuote.Time), true
		}
		if snap.LatestQuote.Ask > 0 {
			return snap.LatestQuote.Ask, providerTimeMillis(snap.LatestQuote.Time), true
		}
		if snap.LatestQuote.Bid > 0 {
			return snap.LatestQuote.Bid, providerTimeMillis(snap.LatestQuote.Time), true
		}
		return 0, 0, false
	}
	if feed == "overnight" {
		if price, stamp, ok := quoteMid(); ok {
			return price, stamp, "indicative-quote"
		}
		if snap.MinuteBar.Close > 0 {
			return snap.MinuteBar.Close, providerTimeMillis(snap.MinuteBar.Time), "minute-bar"
		}
		if snap.LatestTrade.Price > 0 {
			return snap.LatestTrade.Price, providerTimeMillis(snap.LatestTrade.Time), "delayed-trade"
		}
		return 0, 0, ""
	}

	if session == "pre-market" || session == "after-hours" {
		if price, stamp, ok := quoteMid(); ok {
			return price, stamp, "quote"
		}
	}
	if snap.LatestTrade.Price > 0 {
		return snap.LatestTrade.Price, providerTimeMillis(snap.LatestTrade.Time), "trade"
	}
	if price, stamp, ok := quoteMid(); ok {
		return price, stamp, "quote"
	}
	if snap.MinuteBar.Close > 0 {
		return snap.MinuteBar.Close, providerTimeMillis(snap.MinuteBar.Time), "minute-bar"
	}
	return 0, 0, ""
}

func (e *Engine) refreshAlpacaLiveSnapshots(ctx context.Context, key, secret string) {
	session := marketSessionET(time.Now())
	if session == "closed" || session == "weekend" {
		e.setHealth("alpaca-live", "idle · market closed")
		return
	}
	feed := "iex"
	if session == "overnight" {
		e.app.mu.RLock()
		mode := strings.ToLower(strings.TrimSpace(e.app.state.Settings.OvernightDataMode))
		e.app.mu.RUnlock()
		switch mode {
		case "live":
			feed = "boats"
		case "indicative":
			feed = "overnight"
		default:

			feed = "boats"
		}
	}
	tracked := e.trackedSymbols()
	symbols := make([]string, 0, len(tracked))
	for _, symbol := range tracked {
		if symbol != "VIX" {
			symbols = append(symbols, symbol)
		}
	}
	if len(symbols) == 0 {
		return
	}
	client := &http.Client{Timeout: 15 * time.Second}
	updated := 0
	freshest := int64(0)
	freshestSymbol := ""
	for start := 0; start < len(symbols); start += 50 {
		end := minInt(start+50, len(symbols))
		batch := symbols[start:end]
		raw := "https://data.alpaca.markets/v2/stocks/snapshots?symbols=" + url.QueryEscape(strings.Join(batch, ",")) + "&feed=" + url.QueryEscape(feed)
		var payload map[string]alpacaLiveSnapshot
		if err := e.providerGetJSONTier(ctx, "Alpaca", e.workTierForSymbols(batch), client, raw, map[string]string{"APCA-API-KEY-ID": key, "APCA-API-SECRET-KEY": secret}, &payload); err != nil {
			e.recordProviderFailure("Alpaca", err)
			if session == "overnight" && feed == "boats" {
				e.app.mu.RLock()
				mode := strings.ToLower(strings.TrimSpace(e.app.state.Settings.OvernightDataMode))
				e.app.mu.RUnlock()
				if mode == "auto" {
					feed = "overnight"
					raw = "https://data.alpaca.markets/v2/stocks/snapshots?symbols=" + url.QueryEscape(strings.Join(batch, ",")) + "&feed=overnight"
					payload = map[string]alpacaLiveSnapshot{}
					if retryErr := e.providerGetJSONTier(ctx, "Alpaca", e.workTierForSymbols(batch), client, raw, map[string]string{"APCA-API-KEY-ID": key, "APCA-API-SECRET-KEY": secret}, &payload); retryErr == nil {
						e.setHealth("alpaca-live", "live overnight unavailable · using indicative quotes")
					} else {
						e.setError("Alpaca overnight", retryErr)
						e.setHealth("alpaca-live", "overnight unavailable")
						return
					}
				} else {
					e.setError("Alpaca live", err)
					e.setHealth("alpaca-live", "live overnight unavailable")
					return
				}
			} else {
				e.setError("Alpaca live", err)
				e.setHealth("alpaca-live", feed+" unavailable")
				return
			}
		}
		for symbol, snap := range payload {
			price, stamp, priceKind := alpacaSnapshotPrice(snap, feed, session)
			if price <= 0 {
				continue
			}
			if e.mergeAlpacaLiveSnapshot(symbol, price, stamp, snap, feed, priceKind) {
				updated++
			}
			if stamp > freshest {
				freshest = stamp
				freshestSymbol = normalizeSymbol(symbol)
			}
		}
	}
	if updated == 0 {
		e.setHealth("alpaca-live", feed+" · no current quotes")
		return
	}
	e.recordProviderSuccess("Alpaca")
	e.mu.Lock()
	if freshest > e.lastAlpacaAt {
		e.lastAlpacaAt = freshest
		e.lastAlpacaSymbol = freshestSymbol
	}
	e.lastAlpacaFeed = feed
	e.lastUpdated["alpaca-live"] = time.Now().UnixMilli()
	now := time.Now().UnixMilli()
	if freshest > 0 && now-freshest <= 30_000 {
		if feed == "overnight" {
			e.health["alpaca-live"] = fmt.Sprintf("overnight · real-time indicative · %d symbols", updated)
		} else {
			e.health["alpaca-live"] = fmt.Sprintf("%s · current · %d symbols", feed, updated)
		}
	} else {
		e.health["alpaca-live"] = fmt.Sprintf("%s · stale/quiet · %d symbols", feed, updated)
	}
	e.mu.Unlock()
}

func (e *Engine) mergeAlpacaLiveSnapshot(symbol string, price float64, providerAt int64, snap alpacaLiveSnapshot, feed, priceKind string) bool {
	symbol = normalizeSymbol(symbol)
	if symbol == "" || price <= 0 {
		return false
	}
	now := time.Now().UnixMilli()
	obsSource := "alpaca-" + feed + "-snapshot"
	if feed == "overnight" && priceKind == "indicative-quote" {
		obsSource = "alpaca-overnight-indicative-snapshot"
	}
	e.recordProviderQuoteObservation(symbol, Quote{Symbol: symbol, Price: price, Bid: snap.LatestQuote.Bid, Ask: snap.LatestQuote.Ask, ProviderTimestamp: providerAt, UpdatedAt: now, Source: obsSource, FeedType: feed, DataState: map[bool]string{true: "indicative", false: "snapshot"}[feed == "overnight" && priceKind == "indicative-quote"]})
	e.mu.Lock()
	prev := e.quotes[symbol]
	q := prev
	recentAlpacaStream := strings.HasPrefix(q.Source, "alpaca-iex-websocket") && q.ProviderTimestamp > 0 && now-q.ProviderTimestamp <= 20_000
	newer := providerAt == 0 || q.ProviderTimestamp == 0 || providerAt >= q.ProviderTimestamp

	applyPrice := !recentAlpacaStream && newer
	if feed == "overnight" && newer {
		applyPrice = true
	}
	if applyPrice {
		q.Symbol = symbol
		q.Price = price
		q.Bid = snap.LatestQuote.Bid
		q.Ask = snap.LatestQuote.Ask
		q.BidSize = snap.LatestQuote.BidSize
		q.AskSize = snap.LatestQuote.AskSize
		q.FeedType = feed
		if feed == "overnight" && priceKind == "indicative-quote" {
			q.Source = "alpaca-overnight-indicative-snapshot"
			q.DataState = "indicative"
		} else if feed == "boats" {
			q.Source = "alpaca-boats-live-snapshot"
			q.DataState = "live"
		} else {
			q.Source = "alpaca-" + feed + "-snapshot"
			q.DataState = "snapshot"
		}
		q.ProviderTimestamp = providerAt
		if q.ProviderTimestamp == 0 {
			q.ProviderTimestamp = now
		}
		q.UpdatedAt = now
	}
	if snap.DailyBar.Close > 0 {
		q.SessionClose = snap.DailyBar.Close
		q.SessionCloseAt = providerTimeMillis(snap.DailyBar.Time)
	}
	if snap.PrevDailyBar.Close > 0 {
		q.PriorSessionClose = snap.PrevDailyBar.Close
		q.PreviousClose = snap.PrevDailyBar.Close
	}
	if q.PreviousClose > 0 && q.Price > 0 {
		q.Change = q.Price - q.PreviousClose
		q.ChangePercent = q.Change / q.PreviousClose * 100
	}
	e.quotes[symbol] = q
	shouldBroadcast := applyPrice && time.Since(e.lastBroadcast[symbol]) >= 200*time.Millisecond
	if shouldBroadcast {
		e.lastBroadcast[symbol] = time.Now()
	}
	e.mu.Unlock()
	e.propagateCanonicalQuoteChange(symbol, prev, q)
	if shouldBroadcast {
		e.app.broadcastSymbolEvent(q.Symbol, map[string]any{"type": "quote", "quote": q})
	}
	return applyPrice
}
