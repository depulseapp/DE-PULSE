package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"sort"
	"strings"
	"time"
)

func (e *Engine) periodic(ctx context.Context, interval time.Duration, fn func()) {
	fn()
	ticker := time.NewTicker(interval)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			fn()
		}
	}
}

func (e *Engine) liveQuoteLoop(ctx context.Context, key string) {
	backoff := 2 * time.Second
	for {
		select {
		case <-ctx.Done():
			return
		default:
		}
		if !e.providerAllowed("Finnhub") {
			if !sleepContext(ctx, 5*time.Second) {
				return
			}
			continue
		}
		ws, err := DialWebSocket(ctx, "wss://ws.finnhub.io?token="+url.QueryEscape(key))
		if err != nil {
			e.recordProviderFailure("Finnhub", err)
			e.setError("Quotes", err)
			e.mu.Lock()
			e.webSocketConnected = false
			e.mu.Unlock()
			e.setStatus("degraded", "Finnhub stream could not connect. Retrying…")
			if !sleepContext(ctx, backoff) {
				return
			}
			if backoff < 30*time.Second {
				backoff *= 2
			}
			continue
		}
		e.recordProviderSuccess("Finnhub")
		e.mu.Lock()
		e.ws = ws
		e.webSocketConnected = true
		e.wsConnectedAt = time.Now().UnixMilli()
		e.subscribedSymbols = map[string]bool{}
		e.health["quotes"] = "connected"
		e.status = "running"
		e.message = "Finnhub connected · waiting for market trades"
		e.mu.Unlock()
		e.syncLiveSubscriptions(ws)
		e.app.broadcastRuntime()
		backoff = 2 * time.Second

		connCtx, connCancel := context.WithCancel(ctx)
		syncDone := make(chan struct{})
		go func() {
			defer close(syncDone)
			syncTicker := time.NewTicker(3 * time.Second)
			pingTicker := time.NewTicker(20 * time.Second)
			defer syncTicker.Stop()
			defer pingTicker.Stop()
			for {
				select {
				case <-connCtx.Done():
					return
				case <-syncTicker.C:
					e.syncLiveSubscriptions(ws)
				case <-pingTicker.C:
					_ = ws.WritePing()
				}
			}
		}()

		for {
			msg, err := ws.ReadText(ctx)
			if err != nil {
				e.recordProviderFailure("Finnhub", err)
				connCancel()
				<-syncDone
				_ = ws.Close()
				e.mu.Lock()
				e.ws = nil
				e.webSocketConnected = false
				e.wsConnectedAt = 0
				e.subscribedSymbols = map[string]bool{}
				e.health["quotes"] = "reconnecting"
				e.mu.Unlock()
				break
			}
			now := time.Now().UnixMilli()
			e.mu.Lock()
			e.lastMessageAt = now
			e.mu.Unlock()
			var payload struct {
				Type string `json:"type"`
				Data []struct {
					Price     float64 `json:"p"`
					Symbol    string  `json:"s"`
					Timestamp int64   `json:"t"`
					Volume    float64 `json:"v"`
				} `json:"data"`
				Msg string `json:"msg"`
			}
			if json.Unmarshal(msg, &payload) != nil {
				continue
			}
			if payload.Type == "trade" {
				for _, tr := range payload.Data {
					sym := normalizeSymbol(tr.Symbol)
					if tr.Price <= 0 || sym == "" {
						continue
					}
					e.mu.Lock()
					e.lastTradeAt = now
					e.lastTradeSymbol = sym
					e.health["quotes"] = "secondary streaming · Finnhub"
					prev := e.quotes[sym]
					e.mu.Unlock()

					e.recordProviderQuoteObservation(sym, Quote{Symbol: sym, Price: tr.Price, Volume: tr.Volume, ProviderTimestamp: tr.Timestamp, UpdatedAt: now, Source: "finnhub-websocket", FeedType: "websocket", DataState: "live"})

					if quoteIsRecentAlpacaLive(prev, now) || (strings.Contains(strings.ToLower(prev.Source), "alpaca") && prev.ProviderTimestamp > 0 && now-prev.ProviderTimestamp <= 30_000) {
						continue
					}
					e.updateQuote(sym, Quote{Price: tr.Price, Volume: tr.Volume, ProviderTimestamp: tr.Timestamp, FeedType: "websocket", DataState: "live"}, "finnhub-websocket")
				}
			} else if payload.Type == "error" {
				e.setError("Quotes", errors.New(payload.Msg))
			}
		}
		select {
		case <-ctx.Done():
			return
		default:
		}
		e.setStatus("degraded", "Finnhub stream disconnected. Retrying…")
		if !sleepContext(ctx, backoff) {
			return
		}
	}
}

func (e *Engine) syncLiveSubscriptions(ws *WSClient) {
	desired := map[string]bool{}
	for _, symbol := range e.liveSymbols() {
		desired[symbol] = true
	}
	e.mu.RLock()
	current := clone(e.subscribedSymbols)
	e.mu.RUnlock()
	for symbol := range current {
		if desired[symbol] {
			continue
		}
		if ws.WriteText(fmt.Sprintf(`{"type":"unsubscribe","symbol":"%s"}`, symbol)) == nil {
			e.mu.Lock()
			delete(e.subscribedSymbols, symbol)
			e.mu.Unlock()
		}
	}
	for symbol := range desired {
		if current[symbol] {
			continue
		}
		if ws.WriteText(fmt.Sprintf(`{"type":"subscribe","symbol":"%s"}`, symbol)) == nil {
			e.mu.Lock()
			e.subscribedSymbols[symbol] = true
			e.mu.Unlock()
		}
	}
}

const alpacaIEXPoolLimit = 25

// alpacaIEXSymbols returns DE.PULSE v15's preferred live pool. Alpaca IEX is
// kept intentionally below the 30-symbol plan ceiling: 25 normal symbols plus
// five reserve/failover slots. Day -> Swing -> Long-Term -> Discovery priority
// and pinned tradables are retained. VIX is excluded because it is a Cboe index.
func (e *Engine) alpacaIEXSymbols() []string {
	return append([]string{}, e.effectiveAlpacaIEXSymbols()...)
}

func (e *Engine) syncAlpacaSubscriptions(ws *WSClient) {
	if ws == nil {
		return
	}
	desired := map[string]bool{}
	for _, symbol := range e.alpacaIEXSymbols() {
		desired[symbol] = true
	}
	e.mu.RLock()
	current := clone(e.alpacaSubscribedSymbols)
	e.mu.RUnlock()
	adds, removes := []string{}, []string{}
	for symbol := range current {
		if !desired[symbol] {
			removes = append(removes, symbol)
		}
	}
	for symbol := range desired {
		if !current[symbol] {
			adds = append(adds, symbol)
		}
	}
	sort.Strings(adds)
	sort.Strings(removes)
	write := func(action string, symbols []string) error {
		if len(symbols) == 0 {
			return nil
		}
		payload, _ := json.Marshal(map[string]any{"action": action, "trades": symbols, "quotes": symbols})
		return ws.WriteText(string(payload))
	}
	if err := write("unsubscribe", removes); err == nil && len(removes) > 0 {
		e.mu.Lock()
		for _, symbol := range removes {
			delete(e.alpacaSubscribedSymbols, symbol)
		}
		e.mu.Unlock()
	}
	if err := write("subscribe", adds); err == nil && len(adds) > 0 {
		e.mu.Lock()
		for _, symbol := range adds {
			e.alpacaSubscribedSymbols[symbol] = true
		}
		e.mu.Unlock()
	}
}

func (e *Engine) mergeAlpacaIEXStream(symbol string, price, bid, ask float64, providerAt int64, kind string) {
	e.mergeAlpacaIEXStreamAt(symbol, price, bid, ask, providerAt, kind, time.Now())
}

func (e *Engine) mergeAlpacaIEXStreamAt(symbol string, price, bid, ask float64, providerAt int64, kind string, observedAt time.Time) {
	symbol = normalizeSymbol(symbol)
	if symbol == "" || price <= 0 {
		return
	}
	session := marketSessionET(observedAt)
	if session == "overnight" || session == "closed" || session == "weekend" {
		return
	}
	now := observedAt.UnixMilli()
	if providerAt <= 0 {
		providerAt = now
	}
	e.recordProviderQuoteObservation(symbol, Quote{Symbol: symbol, Price: price, Bid: bid, Ask: ask, ProviderTimestamp: providerAt, UpdatedAt: now, Source: "alpaca-iex-websocket-" + kind, FeedType: "iex", DataState: "live"})

	// Resolve allocation before taking the engine write lock. multiFeedAllocation
	// reads engine state and therefore must never be called while e.mu is held.
	alloc := e.multiFeedAllocation()
	e.mu.Lock()
	prev := e.quotes[symbol]
	primaryAlpaca := containsLiveSymbol(alloc.Alpaca, symbol)
	recentFinnhub := prev.Source == "finnhub-websocket" && prev.ProviderTimestamp > 0 && now-prev.ProviderTimestamp <= 20_000
	newer := prev.ProviderTimestamp == 0 || providerAt >= prev.ProviderTimestamp
	if providerAt >= e.lastAlpacaStreamAt {
		e.lastAlpacaStreamAt = providerAt
		e.lastAlpacaStreamSymbol = symbol
	}
	e.health["alpaca-stream"] = fmt.Sprintf("IEX · %d active · %d max · %d reserve", len(e.alpacaSubscribedSymbols), alpacaPlanMaxSymbols, alpacaPlanMaxSymbols-alpacaActiveTarget)
	e.lastUpdated["alpaca-stream"] = now
	e.mu.Unlock()

	if (!primaryAlpaca && recentFinnhub) || !newer {
		if newer && (bid > 0 || ask > 0) {
			e.mu.Lock()
			current := e.quotes[symbol]
			next := current
			if bid > 0 {
				next.Bid = bid
			}
			if ask > 0 {
				next.Ask = ask
			}
			e.quotes[symbol] = next
			e.mu.Unlock()
			e.propagateCanonicalQuoteChange(symbol, current, next)
		}
		return
	}
	source := "alpaca-iex-websocket-trade"
	if kind == "quote" {
		source = "alpaca-iex-websocket-quote"
	}
	e.updateQuote(symbol, Quote{Price: price, Bid: bid, Ask: ask, ProviderTimestamp: providerAt, FeedType: "iex", DataState: "live"}, source)
}

func (e *Engine) alpacaIEXLoop(ctx context.Context, key, secret string) {
	backoff := 3 * time.Second
	for {
		select {
		case <-ctx.Done():
			return
		default:
		}
		if !e.providerAllowed("Alpaca") {
			if !sleepContext(ctx, 5*time.Second) {
				return
			}
			continue
		}
		ws, err := DialWebSocket(ctx, "wss://stream.data.alpaca.markets/v2/iex")
		if err != nil {
			e.recordProviderFailure("Alpaca", err)
			e.setHealth("alpaca-stream", "IEX unavailable · retrying")
			if !sleepContext(ctx, backoff) {
				return
			}
			if backoff < 45*time.Second {
				backoff *= 2
			}
			continue
		}
		auth, _ := json.Marshal(map[string]string{"action": "auth", "key": key, "secret": secret})
		if err := ws.WriteText(string(auth)); err != nil {
			e.recordProviderFailure("Alpaca", err)
			_ = ws.Close()
			e.setHealth("alpaca-stream", "IEX authentication write failed")
			if !sleepContext(ctx, backoff) {
				return
			}
			continue
		}
		authenticated := false
		for !authenticated {
			msg, readErr := ws.ReadText(ctx)
			if readErr != nil {
				err = readErr
				break
			}
			var events []struct {
				Type string `json:"T"`
				Msg  string `json:"msg"`
			}
			if json.Unmarshal(msg, &events) != nil {
				continue
			}
			for _, event := range events {
				if event.Type == "success" && strings.EqualFold(event.Msg, "authenticated") {
					authenticated = true
					break
				}
				if event.Type == "error" {
					err = errors.New(event.Msg)
					break
				}
			}
			if err != nil {
				break
			}
		}
		if !authenticated {
			_ = ws.Close()
			if err != nil {
				e.setHealth("alpaca-stream", "IEX authentication unavailable · snapshots active")
			}
			if !sleepContext(ctx, backoff) {
				return
			}
			if backoff < 45*time.Second {
				backoff *= 2
			}
			continue
		}
		e.recordProviderSuccess("Alpaca")
		e.mu.Lock()
		e.alpacaWS = ws
		e.alpacaWebSocketConnected = true
		e.alpacaSubscribedSymbols = map[string]bool{}
		e.health["alpaca-stream"] = fmt.Sprintf("IEX connected · %d normal · %d reserve", alpacaActiveTarget, alpacaPlanMaxSymbols-alpacaActiveTarget)
		e.mu.Unlock()
		e.syncAlpacaSubscriptions(ws)
		e.app.broadcastRuntime()
		backoff = 3 * time.Second

		connCtx, connCancel := context.WithCancel(ctx)
		syncDone := make(chan struct{})
		go func() {
			defer close(syncDone)
			syncTicker := time.NewTicker(3 * time.Second)
			pingTicker := time.NewTicker(20 * time.Second)
			defer syncTicker.Stop()
			defer pingTicker.Stop()
			for {
				select {
				case <-connCtx.Done():
					return
				case <-syncTicker.C:
					e.syncAlpacaSubscriptions(ws)
				case <-pingTicker.C:
					_ = ws.WritePing()
				}
			}
		}()

		for {
			msg, readErr := ws.ReadText(ctx)
			if readErr != nil {
				e.recordProviderFailure("Alpaca", readErr)
				connCancel()
				<-syncDone
				_ = ws.Close()
				e.mu.Lock()
				e.alpacaWS = nil
				e.alpacaWebSocketConnected = false
				e.alpacaSubscribedSymbols = map[string]bool{}
				e.health["alpaca-stream"] = "IEX reconnecting · snapshots active"
				e.mu.Unlock()
				break
			}
			var events []struct {
				Type   string  `json:"T"`
				Msg    string  `json:"msg"`
				Symbol string  `json:"S"`
				Price  float64 `json:"p"`
				Bid    float64 `json:"bp"`
				Ask    float64 `json:"ap"`
				Time   string  `json:"t"`
			}
			if json.Unmarshal(msg, &events) != nil {
				continue
			}
			for _, event := range events {
				stamp := providerTimeMillis(event.Time)
				switch event.Type {
				case "t":
					if event.Price > 0 {
						e.mergeAlpacaIEXStream(event.Symbol, event.Price, 0, 0, stamp, "trade")
					}
				case "q":
					price := 0.0
					if event.Bid > 0 && event.Ask > 0 {
						price = (event.Bid + event.Ask) / 2
					} else if event.Bid > 0 {
						price = event.Bid
					} else if event.Ask > 0 {
						price = event.Ask
					}
					if price > 0 {
						e.mergeAlpacaIEXStream(event.Symbol, price, event.Bid, event.Ask, stamp, "quote")
					}
				case "error":
					if strings.TrimSpace(event.Msg) != "" {
						if strings.Contains(strings.ToLower(event.Msg), "symbol limit") {
							e.setHealth("alpaca-stream", fmt.Sprintf("IEX capacity reached · %d-symbol pool · snapshots active", alpacaIEXPoolLimit))
						} else {
							e.setHealth("alpaca-stream", "IEX warning · "+event.Msg)
						}
					}
				}
			}
		}
		select {
		case <-ctx.Done():
			return
		default:
		}
		if !sleepContext(ctx, backoff) {
			return
		}
	}
}

func (e *Engine) refreshFeedStatus() {
	e.mu.Lock()
	if e.mode != "live" || (e.status != "running" && e.status != "degraded") {
		e.mu.Unlock()
		return
	}
	session := marketSessionET(time.Now())
	now := time.Now().UnixMilli()
	recentAlpacaStream := e.lastAlpacaStreamAt > 0 && now-e.lastAlpacaStreamAt <= 30_000
	recentAlpaca := e.lastAlpacaAt > 0 && now-e.lastAlpacaAt <= 30_000

	base := e.lastMessageAt
	if base == 0 {
		base = e.wsConnectedAt
	}
	staleSocket := e.webSocketConnected && session != "overnight" && session != "closed" && session != "weekend" && base > 0 && now-base > 120_000
	if staleSocket {
		ws := e.ws
		e.health["quotes"] = "stale stream · reconnecting"
		e.message = "Finnhub stream stopped moving · reconnecting automatically"
		e.mu.Unlock()
		if ws != nil {
			_ = ws.Close()
		}
		e.app.broadcastRuntime()
		return
	}

	switch {
	case session == "weekend":
		e.health["quotes"] = "weekend"
		e.message = "US equities weekend · showing last available trade"
	case session == "closed":
		e.health["quotes"] = "market closed"
		e.message = "US market closed · showing last available trade"
	case session == "overnight":
		if recentAlpaca && (e.lastAlpacaFeed == "overnight" || e.lastAlpacaFeed == "boats") {
			e.health["quotes"] = "alpaca overnight · live"
			e.message = "Overnight session · Alpaca overnight quotes updating"
		} else {
			e.health["quotes"] = "overnight · waiting"
			e.message = "Overnight session · waiting for Alpaca overnight quotes"
		}
	case recentAlpacaStream:
		e.health["quotes"] = "streaming · Alpaca IEX primary"
		e.message = "Alpaca IEX primary stream is updating"
	case recentAlpaca:
		e.health["quotes"] = "current · Alpaca snapshot primary"
		e.message = "Alpaca primary snapshots are current"
	case e.lastTradeAt > 0 && now-e.lastTradeAt <= 30_000:
		e.health["quotes"] = "fallback · Finnhub streaming"
		e.message = "Alpaca is quiet · Finnhub secondary stream is keeping prices live"
	case e.alpacaWebSocketConnected:
		e.health["quotes"] = "connected · Alpaca quiet"
		e.message = "Alpaca primary connected · waiting for a subscribed update"
	case e.webSocketConnected:
		e.health["quotes"] = "fallback connected · Finnhub quiet"
		e.message = "Finnhub secondary connected · waiting for a subscribed update"
	default:
		e.health["quotes"] = "reconnecting"
		e.message = "Primary and fallback equity streams are reconnecting"
	}
	e.mu.Unlock()
	e.app.broadcastRuntime()
}

func sleepContext(ctx context.Context, d time.Duration) bool {
	select {
	case <-ctx.Done():
		return false
	case <-time.After(d):
		return true
	}
}

// catalystQuoteReactionActive keeps quote-level reaction work event-driven.
// A dormant or completed catalyst must not make every quote tick recompute/persist
// reaction state; only a sourced, triggered, incomplete catalyst is decision-active.
func catalystQuoteReactionActive(reactions map[string]CatalystReactionState) bool {
	for _, r := range reactions {
		if r.TriggerAt > 0 && r.CompletedAt == 0 {
			return true
		}
	}
	return false
}

func (e *Engine) updateQuote(symbol string, patch Quote, source string) {
	symbol = normalizeSymbol(symbol)
	if symbol == "" || patch.Price <= 0 {
		return
	}
	e.mu.Lock()
	prev := e.quotes[symbol]
	q := prev
	q.Symbol = symbol
	q.Price = patch.Price
	q.Source = source
	if patch.Bid != 0 {
		q.Bid = patch.Bid
	}
	if patch.Ask != 0 {
		q.Ask = patch.Ask
	}
	if patch.BidSize != 0 {
		q.BidSize = patch.BidSize
	}
	if patch.AskSize != 0 {
		q.AskSize = patch.AskSize
	}
	if patch.FeedType != "" {
		q.FeedType = patch.FeedType
	}
	if patch.DataState != "" {
		q.DataState = patch.DataState
	} else if source == "finnhub-websocket" {
		q.DataState = "live"
	} else if source == "demo" {
		q.DataState = "demo"
	} else {
		q.DataState = "snapshot"
	}
	q.UpdatedAt = time.Now().UnixMilli()
	if patch.Open != 0 {
		q.Open = patch.Open
	}
	if patch.High != 0 {
		q.High = patch.High
	} else if q.High == 0 || patch.Price > q.High {
		q.High = patch.Price
	}
	if patch.Low != 0 {
		q.Low = patch.Low
	} else if q.Low == 0 || patch.Price < q.Low {
		q.Low = patch.Price
	}
	if patch.PreviousClose != 0 {
		q.PreviousClose = patch.PreviousClose
	}
	if patch.SessionClose != 0 {
		q.SessionClose = patch.SessionClose
	}
	if patch.SessionCloseAt != 0 {
		q.SessionCloseAt = patch.SessionCloseAt
	}
	if patch.PriorSessionClose != 0 {
		q.PriorSessionClose = patch.PriorSessionClose
	}
	if patch.Volume != 0 {
		q.Volume = patch.Volume
	}
	if patch.ProviderTimestamp != 0 {
		q.ProviderTimestamp = normalizeObservationMs(patch.ProviderTimestamp)
	}

	if q.PreviousClose != 0 {
		q.Change = q.Price - q.PreviousClose
		q.ChangePercent = q.Change / q.PreviousClose * 100
	}
	h := e.history[symbol]
	now := q.ProviderTimestamp
	if now <= 0 {
		now = q.UpdatedAt
	}
	if len(h) > 0 && now-h[len(h)-1].T < 700 {
		h[len(h)-1] = HistoryPoint{T: now, P: q.Price}
	} else {
		h = append(h, HistoryPoint{T: now, P: q.Price})
	}
	if len(h) > 500 {
		h = h[len(h)-500:]
	}
	e.history[symbol] = h
	e.quotes[symbol] = q
	if q.Bid > 0 && q.Ask >= q.Bid {
		mid := (q.Bid + q.Ask) / 2
		if mid > 0 {
			sp := (q.Ask - q.Bid) / mid * 100
			if sp >= 0 && sp <= 5 {
				b := e.liquidityBaselines[symbol]
				if b.Samples == 0 {
					b.NormalSpreadPct = sp
				} else {
					alpha := 0.08
					b.NormalSpreadPct = b.NormalSpreadPct*(1-alpha) + sp*alpha
				}
				if b.Samples < 10000 {
					b.Samples++
				}
				b.UpdatedAt = q.UpdatedAt
				e.liquidityBaselines[symbol] = b
			}
		}
	}
	e.lastUpdated["quotes"] = q.UpdatedAt
	if symbol == "VIX" {
		e.lastUpdated["vix"] = q.UpdatedAt
	}
	shouldBroadcast := time.Since(e.lastBroadcast[symbol]) >= 200*time.Millisecond
	if shouldBroadcast {
		e.lastBroadcast[symbol] = time.Now()
	}
	e.mu.Unlock()
	e.recordProviderQuoteObservation(symbol, q)
	e.evaluateRapidMoveObservation(symbol, q)
	e.propagateCanonicalQuoteChange(symbol, prev, q)
	if shouldBroadcast {

		historyCopy := append([]HistoryPoint(nil), h...)
		e.app.broadcastSymbolEvent(q.Symbol, map[string]any{"type": "quote", "quote": q, "history": historyCopy})
	}
}
