package main

import (
	"context"
	"fmt"
	"math"
	"net/http"
	"net/url"
	"strings"
	"time"
)

func (e *Engine) refreshTwelveFX(ctx context.Context, key string) {
	pairs := map[string]string{"USD/JPY": "USD / JPY", "EUR/USD": "EUR / USD", "GBP/USD": "GBP / USD", "USD/CAD": "USD / CAD", "USD/CNH": "USD / CNH"}
	updated := 0
	for pair, label := range pairs {
		raw := strings.TrimRight(twelveDataBaseURL, "/") + "/quote?symbol=" + url.QueryEscape(pair) + "&apikey=" + url.QueryEscape(key)
		var p map[string]any
		if err := getJSON(ctx, &http.Client{Timeout: 10 * time.Second}, raw, nil, &p); err != nil {
			continue
		}
		price := toFloat(p["close"])
		if price <= 0 {
			price = toFloat(p["price"])
		}
		if price <= 0 {
			continue
		}
		change := toFloat(p["percent_change"])
		k := "fx_" + strings.NewReplacer("/", "", "-", "").Replace(strings.ToLower(pair))
		e.mu.Lock()
		e.globalDirect[k] = GlobalDriver{Key: k, Label: label, State: func() string {
			if math.Abs(change) < .15 {
				return "NEUTRAL"
			}
			if change > 0 {
				return "UP"
			}
			return "DOWN"
		}(), Value: price, ChangePercent: change, Source: "Twelve Data", Provenance: "DIRECT FX", UpdatedAt: time.Now().UnixMilli(), Confidence: 90, Detail: fmt.Sprintf("Direct FX · %+.2f%%", change), Underlying: pair, ProviderSymbol: pair}
		e.mu.Unlock()
		updated++
	}
	e.mu.Lock()
	if updated > 0 {
		e.lastUpdated["fx-direct"] = time.Now().UnixMilli()
		e.health["fx-direct"] = fmt.Sprintf("healthy · Twelve Data direct FX · %d pairs", updated)
	} else {
		e.health["fx-direct"] = "plan limited or unavailable · proxy/global fallbacks retained"
	}
	e.mu.Unlock()
}

func (a *Application) handleMarketOpenPrep(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed, "POST required")
		return
	}
	a.engine.mu.RLock()
	status, mode := a.engine.status, a.engine.mode
	a.engine.mu.RUnlock()
	if status != "running" && status != "degraded" {
		writeError(w, 409, "Start the runtime before running Market Open Prep.")
		return
	}
	if mode != "live" {
		a.engine.runMarketOpenPrep("manual · demo context")
		writeJSON(w, 200, map[string]any{"ok": true, "message": "Market Open Prep completed using the current Demo context."})
		return
	}
	go a.engine.runMarketOpenPrep("manual")
	writeJSON(w, 202, map[string]any{"ok": true, "message": "Market Open Prep started. Current overnight/premarket evidence is being reconciled without clearing good cache."})
}

// v14.3.2 Data Engine manual actions are thin controls over the same production
// refresh/evaluation paths used by scheduled/runtime processing. They never bypass
// provider entitlement, materiality, canonical-store, or deterministic-score rules.
func (a *Application) dataEngineRuntimeReady() (string, string, bool) {
	a.engine.mu.RLock()
	defer a.engine.mu.RUnlock()
	status, mode := a.engine.status, a.engine.mode
	return status, mode, status == "running" || status == "degraded"
}

func (a *Application) handleCatalystEvaluate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed, "POST required")
		return
	}
	_, _, ready := a.dataEngineRuntimeReady()
	if !ready {
		writeError(w, 409, "Start the runtime before evaluating earnings/material catalysts.")
		return
	}
	now := time.Now()
	a.engine.setManualAction("catalyst-evaluate", "Earnings & Material Catalyst Watch", "RUNNING", "Evaluating scheduled earnings and material News/SEC triggers.", false)
	a.engine.evaluateCatalystWatch(now)
	a.engine.mu.RLock()
	p := a.engine.preparations["catalyst-watch"]
	reactionCount := len(a.engine.catalystReactions)
	a.engine.mu.RUnlock()
	message := "No qualifying material catalyst detected. Event-driven watch remains ready."
	tone := "info"
	switch strings.ToUpper(strings.TrimSpace(p.State)) {
	case "ARMED":
		message = "Scheduled earnings are armed, but sourced release evidence has not arrived yet."
	case "TRIGGERED", "REACTION", "COMPLETE":
		if reactionCount > 0 {
			message = fmt.Sprintf("Catalyst reaction context evaluated for %d symbol(s).", reactionCount)
		}
	case "DEGRADED":
		message, tone = "Catalyst evaluation completed with degraded evidence; no reaction is promoted without material confirmation.", "warning"
	}
	manualState, success := "COMPLETE", true
	if strings.EqualFold(p.State, "DEGRADED") {
		manualState, success = "DEGRADED", false
	}
	a.engine.setManualAction("catalyst-evaluate", "Earnings & Material Catalyst Watch", manualState, message, success)
	writeJSON(w, 200, map[string]any{"ok": true, "message": message, "tone": tone, "state": p.State, "refreshAfterMs": 350})
}

func (a *Application) handleGlobalRefresh(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed, "POST required")
		return
	}
	_, mode, ready := a.dataEngineRuntimeReady()
	if !ready {
		writeError(w, 409, "Start the runtime before refreshing Global / FX context.")
		return
	}
	if mode != "live" {
		a.engine.setManualAction("global-refresh", "Global / FX Context", "COMPLETE", "Demo mode already has local Global / FX context; no provider refresh is required.", true)
		writeJSON(w, 200, map[string]any{"ok": true, "message": "Demo mode already has local Global / FX context; no provider refresh is required.", "tone": "info", "refreshAfterMs": 250})
		return
	}
	a.mu.RLock()
	twelveKey := strings.TrimSpace(a.secrets.TwelveData)
	a.mu.RUnlock()
	a.engine.setManualAction("global-refresh", "Global / FX Context", "RUNNING", "Refreshing direct Global / FX evidence and official-close fallbacks.", false)
	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), 75*time.Second)
		defer cancel()
		if twelveKey != "" {
			a.engine.refreshDirectGlobal(ctx, twelveKey)
			a.engine.refreshTwelveFX(ctx, twelveKey)
		}
		a.engine.refreshOfficialGlobalCloses(ctx)
		a.engine.mu.Lock()
		a.engine.lastUpdated["manual-global-refresh"] = time.Now().UnixMilli()
		a.engine.mu.Unlock()
		a.engine.mu.RLock()
		g, fx := a.engine.health["global-direct"], a.engine.health["fx-direct"]
		a.engine.mu.RUnlock()
		degraded := strings.Contains(strings.ToLower(g+" "+fx), "degraded") || strings.Contains(strings.ToLower(g+" "+fx), "unavailable")
		if degraded {
			a.engine.setManualAction("global-refresh", "Global / FX Context", "DEGRADED", g+" · "+fx, false)
		} else {
			a.engine.setManualAction("global-refresh", "Global / FX Context", "COMPLETE", g+" · "+fx, true)
		}
		a.broadcastRuntime()
	}()
	writeJSON(w, 202, map[string]any{"ok": true, "message": "Refreshing direct Global / FX evidence and official-close fallbacks without clearing cache.", "tone": "info", "refreshAfterMs": 1200})
}

func (a *Application) handleCapabilityRecheck(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed, "POST required")
		return
	}
	_, mode, ready := a.dataEngineRuntimeReady()
	if !ready {
		writeError(w, 409, "Start the runtime before rechecking provider capabilities.")
		return
	}
	if mode != "live" {
		a.engine.setManualAction("capability-recheck", "Provider Capabilities", "COMPLETE", "Demo capability state rechecked; live entitlements require Live mode.", true)
		a.engine.mu.Lock()
		a.engine.health["provider-capabilities"] = "ready · demo capability state"
		a.engine.lastUpdated["provider-capabilities"] = time.Now().UnixMilli()
		a.engine.mu.Unlock()
		writeJSON(w, 200, map[string]any{"ok": true, "message": "Demo capability state rechecked. Live entitlements are verified only in Live mode.", "tone": "info", "refreshAfterMs": 250})
		return
	}
	a.mu.RLock()
	finnhubKey := strings.TrimSpace(a.secrets.Finnhub)
	alpacaKey := strings.TrimSpace(a.secrets.AlpacaKey)
	alpacaSecret := strings.TrimSpace(a.secrets.AlpacaSecret)
	fredKey := strings.TrimSpace(a.secrets.FRED)
	eiaKey := strings.TrimSpace(a.secrets.EIA)
	twelveKey := strings.TrimSpace(a.secrets.TwelveData)
	a.mu.RUnlock()
	a.engine.setManualAction("capability-recheck", "Provider Capabilities", "RUNNING", "Rechecking provider capabilities and entitlements.", false)
	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
		defer cancel()
		a.engine.mu.Lock()
		a.engine.health["provider-capabilities"] = "checking · manual recheck"
		a.engine.mu.Unlock()
		if finnhubKey != "" {
			a.engine.refreshFinnhubIntelligence(ctx, finnhubKey)
		}
		if alpacaKey != "" && alpacaSecret != "" {
			a.engine.refreshAlpacaMarketCalendar(ctx, alpacaKey, alpacaSecret)
			a.engine.refreshAlpacaMarketActivity(ctx, alpacaKey, alpacaSecret)
			a.engine.refreshAlpacaCorporateActions(ctx, alpacaKey, alpacaSecret)
		}
		if fredKey != "" {
			a.engine.refreshFRED(ctx, fredKey)
		}
		a.engine.refreshOfficialMacroActuals(ctx, eiaKey)
		if twelveKey != "" {

			a.engine.refreshTwelveFX(ctx, twelveKey)
		}
		a.engine.mu.Lock()
		a.engine.health["provider-capabilities"] = "complete · verified from current provider evidence"
		a.engine.lastUpdated["provider-capabilities"] = time.Now().UnixMilli()
		a.engine.mu.Unlock()
		a.engine.setManualAction("capability-recheck", "Provider Capabilities", "COMPLETE", "Provider capability evidence rechecked; unavailable entitlements remain truthfully labeled.", true)
		a.broadcastRuntime()
	}()
	writeJSON(w, 202, map[string]any{"ok": true, "message": "Provider capability recheck started. Entitlements remain unavailable until positively verified.", "tone": "info", "refreshAfterMs": 1500})
}

func (a *Application) handleVIXRefresh(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed, "POST required")
		return
	}
	_, mode, ready := a.dataEngineRuntimeReady()
	if !ready {
		writeError(w, 409, "Start the runtime before refreshing VIX.")
		return
	}
	if mode != "live" {
		a.engine.setManualAction("vix-refresh", "VIX", "COMPLETE", "Demo VIX is generated locally; no provider refresh is required.", true)
		writeJSON(w, 200, map[string]any{"ok": true, "message": "Demo VIX is generated locally; no provider refresh is required.", "tone": "info", "refreshAfterMs": 250})
		return
	}
	a.mu.RLock()
	finnhubKey := strings.TrimSpace(a.secrets.Finnhub)
	a.mu.RUnlock()
	a.engine.setManualAction("vix-refresh", "VIX", "RUNNING", "Refreshing VIX through the validated direct → official → cache path.", false)
	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), 45*time.Second)
		defer cancel()
		a.engine.refreshVIXSnapshot(ctx, finnhubKey)
		a.engine.mu.Lock()
		a.engine.lastUpdated["manual-vix-refresh"] = time.Now().UnixMilli()
		a.engine.mu.Unlock()
		a.engine.mu.RLock()
		detail := a.engine.health["vix"]
		a.engine.mu.RUnlock()
		degraded := strings.Contains(strings.ToLower(detail), "unavailable") || strings.Contains(strings.ToLower(detail), "degraded") || strings.Contains(strings.ToLower(detail), "failed")
		if degraded {
			a.engine.setManualAction("vix-refresh", "VIX", "DEGRADED", detail, false)
		} else {
			a.engine.setManualAction("vix-refresh", "VIX", "COMPLETE", detail, true)
		}
		a.broadcastRuntime()
	}()
	writeJSON(w, 202, map[string]any{"ok": true, "message": "VIX refresh started using the validated direct → official → cache path.", "tone": "info", "refreshAfterMs": 1000})
}

func (a *Application) handleStreamReconnect(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed, "POST required")
		return
	}
	_, mode, ready := a.dataEngineRuntimeReady()
	if !ready {
		writeError(w, 409, "Start the runtime before reconnecting the live stream.")
		return
	}
	if mode != "live" {
		a.engine.setManualAction("stream-reconnect", "Primary Live Stream", "COMPLETE", "Demo mode does not use a provider WebSocket; no reconnect is required.", true)
		writeJSON(w, 200, map[string]any{"ok": true, "message": "Demo mode does not use a provider WebSocket; no reconnect is required.", "tone": "info", "refreshAfterMs": 250})
		return
	}

	a.engine.mu.Lock()
	alpacaWS, alpacaConnected := a.engine.alpacaWS, a.engine.alpacaWebSocketConnected
	finnhubWS, finnhubConnected := a.engine.ws, a.engine.webSocketConnected
	provider := ""
	var ws *WSClient
	if alpacaConnected && alpacaWS != nil {
		provider, ws = "Alpaca", alpacaWS
	} else if finnhubConnected && finnhubWS != nil {
		provider, ws = "Finnhub", finnhubWS
	}
	if ws != nil {
		a.engine.health["quotes"] = "manual reconnect requested · alternate provider remains available"
		a.engine.message = "Primary routed live stream reconnect requested"
	}
	a.engine.mu.Unlock()
	if ws == nil {
		a.engine.setManualAction("stream-reconnect", "Primary Live Stream", "DEGRADED", "No connected routed live stream is available to recycle; automatic reconnect/fallback remains active.", false)
		writeJSON(w, 200, map[string]any{"ok": true, "message": "Live providers are already disconnected/reconnecting automatically; fallback remains active where available.", "tone": "warn", "refreshAfterMs": 500})
		return
	}
	a.engine.setManualAction("stream-reconnect", "Primary Live Stream", "RUNNING", provider+" reconnect requested; alternate routed providers remain active.", false)
	_ = ws.Close()
	go func(target string) {
		deadline := time.Now().Add(35 * time.Second)
		for time.Now().Before(deadline) {
			time.Sleep(500 * time.Millisecond)
			a.engine.mu.RLock()
			ok := (target == "Alpaca" && a.engine.alpacaWebSocketConnected) || (target == "Finnhub" && a.engine.webSocketConnected)
			a.engine.mu.RUnlock()
			if ok {
				a.engine.setManualAction("stream-reconnect", "Primary Live Stream", "COMPLETE", target+" stream reconnected and routed subscriptions were restored.", true)
				a.broadcastRuntime()
				return
			}
		}
		a.engine.setManualAction("stream-reconnect", "Primary Live Stream", "DEGRADED", target+" reconnect is still pending; alternate provider recovery remains available.", false)
		a.broadcastRuntime()
	}(provider)
	a.broadcastRuntime()
	writeJSON(w, 202, map[string]any{"ok": true, "message": provider + " routed stream reconnect requested. The existing production loop will re-authenticate and resubscribe while fallback remains active.", "tone": "info", "refreshAfterMs": 1000})
}
