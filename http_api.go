package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mime"
	"net/http"
	"net/url"
	"os"
	"path"
	"sort"
	"strings"
	"time"
)

func (a *Application) routes() http.Handler {
	mux := http.NewServeMux()
	a.registerHealthRoutes(mux)
	mux.HandleFunc("/api/auth/status", a.handleAuthStatus)
	mux.HandleFunc("/api/auth/login", a.handleLogin)
	mux.HandleFunc("/api/auth/logout", a.authAllowPasswordSetup(a.handleLogout))
	mux.HandleFunc("/api/auth/set-password", a.authAllowPasswordSetup(a.handleSetPassword))
	mux.HandleFunc("/api/auth/profile", a.auth(a.requireRecentAuthentication(postOnly(a.handleUpdateProfile))))
	mux.HandleFunc("/api/auth/rotate", a.auth(postOnly(a.handleRotateSession)))
	mux.HandleFunc("/api/auth/reauth", a.auth(postOnly(a.handleReauthenticate)))
	mux.HandleFunc("/api/bootstrap", a.auth(a.handleBootstrap))
	mux.HandleFunc("/api/events", a.auth(a.handleEvents))
	mux.HandleFunc("/api/admin/identity", a.requireRole(RoleAdmin, a.handleAdminIdentity))
	mux.HandleFunc("/api/admin/users/create", a.requireRole(RoleAdmin, a.requireRecentAuthentication(postOnly(a.handleAdminUserCreate))))
	mux.HandleFunc("/api/admin/users/role", a.requireRole(RoleAdmin, a.requireRecentAuthentication(postOnly(a.handleAdminUserRole))))
	mux.HandleFunc("/api/admin/users/status", a.requireRole(RoleAdmin, a.requireRecentAuthentication(postOnly(a.handleAdminUserStatus))))
	mux.HandleFunc("/api/admin/users/reset-password", a.requireRole(RoleAdmin, a.requireRecentAuthentication(postOnly(a.handleAdminPasswordReset))))
	mux.HandleFunc("/api/admin/sessions/revoke", a.requireRole(RoleAdmin, a.requireRecentAuthentication(postOnly(a.handleAdminSessionRevoke))))
	mux.HandleFunc("/api/settings/save", a.requireRole(RoleAdmin, a.requireRecentAuthentication(postOnly(a.handleSettingsSave))))
	mux.HandleFunc("/api/settings/ai-provider", a.requireRole(RoleAdmin, postOnly(a.handleAIProviderSelect)))
	mux.HandleFunc("/api/settings/clear-secret", a.requireRole(RoleAdmin, a.requireRecentAuthentication(postOnly(a.handleClearSecret))))
	mux.HandleFunc("/api/provider/test", a.requireRole(RoleAdmin, postOnly(a.handleProviderTest)))
	mux.HandleFunc("/api/cache/clear", a.requireRole(RoleAdmin, postOnly(a.handleCacheClear)))
	mux.HandleFunc("/api/cache/refresh", a.auth(postOnly(a.handleCacheRefresh)))
	mux.HandleFunc("/api/cache/pre-market-prep", a.auth(postOnly(a.handlePreMarketPrep)))
	mux.HandleFunc("/api/cache/market-open-prep", a.auth(postOnly(a.handleMarketOpenPrep)))
	mux.HandleFunc("/api/data-engine/catalyst-evaluate", a.auth(postOnly(a.handleCatalystEvaluate)))
	mux.HandleFunc("/api/data-engine/global-refresh", a.requireRole(RoleAdmin, postOnly(a.handleGlobalRefresh)))
	mux.HandleFunc("/api/data-engine/capabilities-recheck", a.requireRole(RoleAdmin, postOnly(a.handleCapabilityRecheck)))
	mux.HandleFunc("/api/data-engine/vix-refresh", a.auth(postOnly(a.handleVIXRefresh)))
	mux.HandleFunc("/api/data-engine/stream-reconnect", a.requireRole(RoleAdmin, postOnly(a.handleStreamReconnect)))
	mux.HandleFunc("/api/data-engine/refresh", a.auth(postOnly(a.handleTargetedRefresh)))
	mux.HandleFunc("/api/research/refresh", a.auth(postOnly(a.handleResearchRefresh)))
	mux.HandleFunc("/api/desk/membership", a.auth(postOnly(a.handleDeskMembership)))
	mux.HandleFunc("/api/master-symbol/add", a.auth(postOnly(a.handleMasterSymbolAdd)))
	mux.HandleFunc("/api/master-symbol/remove", a.auth(postOnly(a.handleMasterSymbolRemove)))
	mux.HandleFunc("/api/master-symbol/remove-all", a.auth(postOnly(a.handleMasterSymbolRemoveAll)))
	mux.HandleFunc("/api/master-symbol/restore", a.auth(postOnly(a.handleMasterSymbolRestore)))
	mux.HandleFunc("/api/cache/integrity", a.requireRole(RoleAdmin, postOnly(a.handleIntegrityCheck)))
	mux.HandleFunc("/api/maintenance/run", a.requireRole(RoleAdmin, postOnly(a.handleMaintenanceRun)))
	mux.HandleFunc("/api/discovery/scan", a.auth(postOnly(a.handleDiscoveryScan)))
	mux.HandleFunc("/api/engine/toggle", a.requireRole(RoleAdmin, postOnly(a.handleEngineToggle)))
	mux.HandleFunc("/api/runtime/start", a.requireRole(RoleAdmin, postOnly(a.handleRuntimeStart)))
	mux.HandleFunc("/api/runtime/stop", a.requireRole(RoleAdmin, postOnly(a.handleRuntimeStop)))
	mux.HandleFunc("/api/ui/scope", a.auth(postOnly(a.handleScope)))
	mux.HandleFunc("/api/ui/ticker", a.auth(postOnly(a.handleTicker)))
	mux.HandleFunc("/api/live-priority", a.auth(postOnly(a.handleLivePriority)))
	mux.HandleFunc("/api/watchlists/create", a.auth(postOnly(a.handleWatchlistCreate)))
	mux.HandleFunc("/api/watchlists/rename", a.auth(postOnly(a.handleWatchlistRename)))
	mux.HandleFunc("/api/watchlists/delete", a.auth(postOnly(a.handleWatchlistDelete)))
	mux.HandleFunc("/api/watchlists/add-symbol", a.auth(postOnly(a.handleAddSymbol)))
	mux.HandleFunc("/api/watchlists/remove-symbol", a.auth(postOnly(a.handleRemoveSymbol)))
	mux.HandleFunc("/api/ai/generate", a.auth(postOnly(a.handleAI)))
	mux.HandleFunc("/api/signal-validation/record", a.auth(postOnly(a.handleSignalValidationRecord)))
	mux.HandleFunc("/api/community/evidence", a.auth(postOnly(a.handleCommunityEvidence)))
	mux.HandleFunc("/api/profile/export", a.auth(a.handleExport))
	mux.HandleFunc("/api/profile/import", a.requireRole(RoleAdmin, postOnly(a.handleImport)))
	mux.HandleFunc("/api/app/open-url", a.auth(postOnly(a.handleOpenURL)))
	mux.HandleFunc("/api/app/quit", a.requireRole(RoleAdmin, postOnly(a.handleQuit)))
	mux.HandleFunc("/", a.handleStatic)
	return securityHeaders(a.observeHTTP(mux))
}

func securityHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("X-Frame-Options", "DENY")
		w.Header().Set("Referrer-Policy", "no-referrer")
		next.ServeHTTP(w, r)
	})
}
func (a *Application) handleStatic(w http.ResponseWriter, r *http.Request) {
	if a.identity == nil {
		if r.URL.Path == "/" || r.URL.Path == "/index.html" {
			http.SetCookie(w, &http.Cookie{Name: "pmt_session", Value: a.sessionKey, Path: "/", HttpOnly: true, SameSite: http.SameSiteStrictMode})
			r.URL.Path = "/renderer/index.html"
		} else {
			r.URL.Path = "/renderer" + r.URL.Path
		}
		clean := strings.TrimPrefix(path.Clean(r.URL.Path), "/")
		data, err := rendererFiles.ReadFile(clean)
		if err != nil {
			http.NotFound(w, r)
			return
		}
		if mt := mime.TypeByExtension(path.Ext(clean)); mt != "" {
			w.Header().Set("Content-Type", mt)
		}
		_, _ = w.Write(data)
		return
	}
	if r.URL.Path == "/" || r.URL.Path == "/index.html" {
		if p, err := a.identity.resolve(sessionTokenFromRequest(r), false); err == nil {
			if a.identity.userRequiresPassword(p.UserID) {
				r.URL.Path = "/renderer/setup.html"
			} else {
				r.URL.Path = "/renderer/index.html"
			}
		} else if a.identity.ownerNeedsPassword() {
			// One-time local migration compatibility: mint only a password-setup session.
			// auth() blocks this principal from all application APIs until Argon2id setup completes.
			token, _, bootErr := a.identity.bootstrapOwnerSession()
			if bootErr != nil {
				writeError(w, http.StatusServiceUnavailable, "Identity bootstrap unavailable.")
				return
			}
			setSessionCookie(w, r, token)
			r.URL.Path = "/renderer/setup.html"
		} else {
			r.URL.Path = "/renderer/login.html"
		}
	} else {
		r.URL.Path = "/renderer" + r.URL.Path
	}
	clean := strings.TrimPrefix(path.Clean(r.URL.Path), "/")
	data, err := rendererFiles.ReadFile(clean)
	if err != nil {
		http.NotFound(w, r)
		return
	}
	if mt := mime.TypeByExtension(path.Ext(clean)); mt != "" {
		w.Header().Set("Content-Type", mt)
	}
	w.Header().Set("Cache-Control", "no-store, no-cache, must-revalidate, max-age=0")
	w.Header().Set("Pragma", "no-cache")
	w.Header().Set("Expires", "0")
	_, _ = w.Write(data)
}

func (a *Application) handleBootstrap(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeError(w, 405, "Method not allowed")
		return
	}
	p, _ := principalFromContext(r.Context())
	writeJSON(w, 200, map[string]any{"state": a.publicStateForUser(p.UserID), "runtime": a.engine.SnapshotForUser(p.UserID), "principal": p})
}
func (a *Application) handleEvents(w http.ResponseWriter, r *http.Request) {
	flusher, ok := w.(http.Flusher)
	if !ok {
		writeError(w, 500, "Streaming unsupported")
		return
	}
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	p, _ := principalFromContext(r.Context())
	ch := a.hub.SubscribeForUser(p.UserID)
	defer a.hub.Unsubscribe(ch)
	initial, _ := json.Marshal(map[string]any{"type": "bootstrap", "state": a.publicStateForUser(p.UserID), "runtime": a.engine.SnapshotForUser(p.UserID)})
	fmt.Fprintf(w, "data: %s\n\n", initial)
	flusher.Flush()
	ticker := time.NewTicker(20 * time.Second)
	defer ticker.Stop()
	for {
		select {
		case <-r.Context().Done():
			return
		case data := <-ch:
			fmt.Fprintf(w, "data: %s\n\n", data)
			flusher.Flush()
		case <-ticker.C:
			if a.identity != nil {
				if !a.identity.sessionIDActive(p.SessionID) {
					return
				}
				a.identity.touchSessionID(p.SessionID)
			}
			fmt.Fprint(w, ": keepalive\n\n")
			flusher.Flush()
		}
	}
}

func decodeJSON(r *http.Request, out any) error {
	defer r.Body.Close()
	dec := json.NewDecoder(io.LimitReader(r.Body, 2<<20))
	if err := dec.Decode(out); err != nil {
		return err
	}
	// Reject trailing JSON/tokens. Mutating local-app endpoints must never
	// partially accept a valid prefix and silently ignore malformed suffixes.
	var extra any
	if err := dec.Decode(&extra); err != io.EOF {
		if err == nil {
			return errors.New("multiple JSON values are not allowed")
		}
		return err
	}
	return nil
}

func postOnly(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			writeError(w, http.StatusMethodNotAllowed, "Method not allowed")
			return
		}
		next(w, r)
	}
}
func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}
func writeError(w http.ResponseWriter, status int, msg string) {
	writeJSON(w, status, map[string]string{"error": msg})
}
func (a *Application) broadcastState() {
	a.broadcastStateForUser(bootstrapOwnerID)
}
func (a *Application) broadcastStateForUser(userID string) {
	a.hub.BroadcastUser(userID, map[string]any{"type": "state", "state": a.publicStateForUser(userID)})
}

// broadcastSharedState updates every authenticated workspace with the same
// shared operational Settings/status while preserving each user's private
// watchlists and UI state. Global configuration must never be broadcast using
// one user's public-state payload.
func (a *Application) broadcastSharedState() {
	a.mu.RLock()
	ids := make([]string, 0, len(a.workspaces))
	for userID := range a.workspaces {
		ids = append(ids, userID)
	}
	a.mu.RUnlock()
	if len(ids) == 0 {
		a.broadcastState()
		return
	}
	sort.Strings(ids)
	for _, userID := range ids {
		a.broadcastStateForUser(userID)
	}
}

func cleanCredential(v string) string {
	v = strings.TrimSpace(v)
	v = strings.Trim(v, "\"'` ")
	v = strings.ReplaceAll(v, "\r", "")
	v = strings.ReplaceAll(v, "\n", "")
	v = strings.ReplaceAll(v, "\t", "")
	return v
}

func parseRFC3339Unix(v string) int64 {
	t, err := time.Parse(time.RFC3339, strings.TrimSpace(v))
	if err != nil {
		return 0
	}
	return t.Unix()
}
func (a *Application) handleSettingsSave(w http.ResponseWriter, r *http.Request) {
	var in struct {
		Settings        Settings `json:"settings"`
		FinnhubKey      string   `json:"finnhubKey"`
		TradeInsightKey string   `json:"tradeInsightKey"`
		AlpacaKey       string   `json:"alpacaKey"`
		AlpacaSecret    string   `json:"alpacaSecret"`
		GroqKey         string   `json:"groqKey"`
		OpenRouterKey   string   `json:"openRouterKey"`
		GeminiKey       string   `json:"geminiKey"`
		FREDKey         string   `json:"fredKey"`
		BLSKey          string   `json:"blsKey"`
		EIAKey          string   `json:"eiaKey"`
		TwelveDataKey   string   `json:"twelveDataKey"`
		MarketauxKey    string   `json:"marketauxKey"`
	}
	if decodeJSON(r, &in) != nil {
		writeError(w, 400, "Invalid settings")
		return
	}
	userID := requestUserID(r.Context())
	a.mu.Lock()
	st := in.Settings
	if st.DataMode != "live" {
		st.DataMode = "demo"
	}
	if strings.TrimSpace(st.GroqModel) == "" {
		st.GroqModel = "openai/gpt-oss-120b"
	}
	if strings.TrimSpace(st.GeminiModel) == "" {
		st.GeminiModel = "gemini-3.1-flash-lite"
	}
	switch strings.ToLower(strings.TrimSpace(st.OpenRouterMode)) {
	case "fast", "balanced", "powerful", "specific":
	default:
		st.OpenRouterMode = "fast"
	}
	if !allowedOpenRouterModel(st.OpenRouterSpecificModel) {
		st.OpenRouterSpecificModel = "openai/gpt-5.6-sol"
	}
	switch strings.ToLower(strings.TrimSpace(st.AIProvider)) {
	case "groq", "openrouter", "gemini":
		st.AIProvider = strings.ToLower(strings.TrimSpace(st.AIProvider))
	default:
		st.AIProvider = "groq"
	}
	switch strings.ToLower(strings.TrimSpace(st.AIRoutingMode)) {
	case "manual", "efficient", "balanced", "deep":
		st.AIRoutingMode = strings.ToLower(strings.TrimSpace(st.AIRoutingMode))
	default:
		st.AIRoutingMode = "manual"
	}
	if st.SignalProfile == "" {
		st.SignalProfile = "balanced"
	}
	if st.MarketContext <= 0 {
		st.MarketContext = 15
	}
	if st.EarningsPenalty <= 0 {
		st.EarningsPenalty = 10
	}
	st.SECEmail = strings.TrimSpace(st.SECEmail)
	st.SwingWatchlistID = "swing"
	st.DayWatchlistID = "day"
	st.LongWatchlistID = "long"
	st.DiscoveryWatchlistID = "discovery"
	switch strings.ToLower(strings.TrimSpace(st.OvernightDataMode)) {
	case "auto", "indicative", "live":
	default:
		st.OvernightDataMode = "auto"
	}
	switch strings.ToLower(strings.TrimSpace(st.GlobalProviderMode)) {
	case "auto", "direct", "free-first", "proxy":
		st.GlobalProviderMode = strings.ToLower(strings.TrimSpace(st.GlobalProviderMode))
	default:
		st.GlobalProviderMode = "auto"
	}
	switch strings.ToLower(strings.TrimSpace(st.OptionsDataMode)) {
	case "auto", "opra", "indicative", "off":
		st.OptionsDataMode = strings.ToLower(strings.TrimSpace(st.OptionsDataMode))
	default:
		st.OptionsDataMode = "auto"
	}
	switch strings.ToLower(strings.TrimSpace(st.ResearchAIMode)) {
	case "manual", "automatic":
		st.ResearchAIMode = strings.ToLower(strings.TrimSpace(st.ResearchAIMode))
	default:
		st.ResearchAIMode = "manual"
	}
	a.state.Settings = st
	a.state.SettingsSavedAt = time.Now().UnixMilli()
	ensureDedicatedDeskWatchlists(&a.state, defaultState())
	if v := cleanCredential(in.FinnhubKey); v != "" {
		a.secrets.Finnhub = v
	}
	if v := cleanCredential(in.TradeInsightKey); v != "" {
		a.secrets.TradeInsight = v
	}
	if v := cleanCredential(in.AlpacaKey); v != "" {
		a.secrets.AlpacaKey = v
	}
	if v := cleanCredential(in.AlpacaSecret); v != "" {
		a.secrets.AlpacaSecret = v
	}
	if v := cleanCredential(in.GroqKey); v != "" {
		a.secrets.Groq = v
	}
	if v := cleanCredential(in.OpenRouterKey); v != "" {
		a.secrets.OpenRouter = v
	}
	if v := cleanCredential(in.GeminiKey); v != "" {
		a.secrets.Gemini = v
	}
	if v := cleanCredential(in.FREDKey); v != "" {
		a.secrets.FRED = v
	}
	if v := cleanCredential(in.BLSKey); v != "" {
		a.secrets.BLS = v
	}
	if v := cleanCredential(in.EIAKey); v != "" {
		a.secrets.EIA = v
	}
	if v := cleanCredential(in.TwelveDataKey); v != "" {
		a.secrets.TwelveData = v
	}
	if v := cleanCredential(in.MarketauxKey); v != "" {
		a.secrets.Marketaux = v
	}
	err := a.saveLocked()
	state := a.publicStateLockedForUser(userID)
	a.mu.Unlock()
	if err != nil {
		writeError(w, 500, err.Error())
		return
	}
	a.broadcastSharedState()
	writeJSON(w, 200, state)
}
func (a *Application) handleAIProviderSelect(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}
	var in struct {
		Provider string `json:"provider"`
	}
	if decodeJSON(r, &in) != nil {
		writeError(w, http.StatusBadRequest, "Invalid AI provider selection")
		return
	}
	provider := strings.ToLower(strings.TrimSpace(in.Provider))
	switch provider {
	case "groq", "openrouter", "gemini":
	default:
		writeError(w, http.StatusBadRequest, "Choose Groq, OpenRouter, or Google Gemini")
		return
	}
	userID := requestUserID(r.Context())
	a.mu.Lock()
	a.state.Settings.AIProvider = provider
	a.state.SettingsSavedAt = time.Now().UnixMilli()
	err := a.saveLocked()
	state := a.publicStateLockedForUser(userID)
	a.mu.Unlock()
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	a.broadcastSharedState()
	writeJSON(w, http.StatusOK, state)
}

func (a *Application) handleClearSecret(w http.ResponseWriter, r *http.Request) {
	var in struct {
		Name string `json:"name"`
	}
	if decodeJSON(r, &in) != nil {
		writeError(w, 400, "Invalid clear-secret request")
		return
	}
	userID := requestUserID(r.Context())
	a.mu.Lock()
	switch in.Name {
	case "finnhub":
		a.secrets.Finnhub = ""
	case "tradeinsight":
		a.secrets.TradeInsight = ""
	case "alpaca":
		a.secrets.AlpacaKey = ""
		a.secrets.AlpacaSecret = ""
	case "groq":
		a.secrets.Groq = ""
	case "openrouter":
		a.secrets.OpenRouter = ""
	case "gemini":
		a.secrets.Gemini = ""
	case "fred":
		a.secrets.FRED = ""
	case "bls":
		a.secrets.BLS = ""
	case "eia":
		a.secrets.EIA = ""
	case "twelvedata":
		a.secrets.TwelveData = ""
	case "marketaux":
		a.secrets.Marketaux = ""
	}
	_ = a.saveLocked()
	state := a.publicStateLockedForUser(userID)
	a.mu.Unlock()
	a.broadcastSharedState()
	writeJSON(w, 200, state)
}
func (a *Application) handleProviderTest(w http.ResponseWriter, r *http.Request) {
	var in struct {
		Provider        string `json:"provider"`
		FinnhubKey      string `json:"finnhubKey"`
		TradeInsightKey string `json:"tradeInsightKey"`
		AlpacaKey       string `json:"alpacaKey"`
		AlpacaSecret    string `json:"alpacaSecret"`
		GroqKey         string `json:"groqKey"`
		OpenRouterKey   string `json:"openRouterKey"`
		GeminiKey       string `json:"geminiKey"`
		FREDKey         string `json:"fredKey"`
		BLSKey          string `json:"blsKey"`
		EIAKey          string `json:"eiaKey"`
		TwelveDataKey   string `json:"twelveDataKey"`
		MarketauxKey    string `json:"marketauxKey"`
	}
	if decodeJSON(r, &in) != nil {
		writeError(w, 400, "Invalid provider-test request")
		return
	}
	a.mu.RLock()
	fk := a.secrets.Finnhub
	ti := a.secrets.TradeInsight
	ak := a.secrets.AlpacaKey
	as := a.secrets.AlpacaSecret
	gk := a.secrets.Groq
	ork := a.secrets.OpenRouter
	gmk := a.secrets.Gemini
	fred := a.secrets.FRED
	bls := a.secrets.BLS
	eia := a.secrets.EIA
	twelve := a.secrets.TwelveData
	marketaux := a.secrets.Marketaux
	settings := a.state.Settings
	a.mu.RUnlock()
	if v := cleanCredential(in.FinnhubKey); v != "" {
		fk = v
	}
	if v := cleanCredential(in.TradeInsightKey); v != "" {
		ti = v
	}
	if v := cleanCredential(in.AlpacaKey); v != "" {
		ak = v
	}
	if v := cleanCredential(in.AlpacaSecret); v != "" {
		as = v
	}
	if v := cleanCredential(in.GroqKey); v != "" {
		gk = v
	}
	if v := cleanCredential(in.OpenRouterKey); v != "" {
		ork = v
	}
	if v := cleanCredential(in.GeminiKey); v != "" {
		gmk = v
	}
	if v := cleanCredential(in.FREDKey); v != "" {
		fred = v
	}
	if v := cleanCredential(in.BLSKey); v != "" {
		bls = v
	}
	if v := cleanCredential(in.EIAKey); v != "" {
		eia = v
	}
	if v := cleanCredential(in.TwelveDataKey); v != "" {
		twelve = v
	}
	if v := cleanCredential(in.MarketauxKey); v != "" {
		marketaux = v
	}
	ctx, cancel := context.WithTimeout(r.Context(), 20*time.Second)
	defer cancel()
	provider := strings.ToLower(strings.TrimSpace(in.Provider))
	var result ProviderTestResult
	switch provider {
	case "tradeinsight":
		result = testTradeInsight(ctx, ti)
	case "alpaca":
		result = testAlpaca(ctx, ak, as)
	case "groq":
		result = testAIProvider(ctx, "groq", gk, settings.GroqModel)
	case "openrouter":
		result = testOpenRouter(ctx, ork, settings.OpenRouterMode, settings.OpenRouterSpecificModel)
	case "gemini":
		result = testAIProvider(ctx, "gemini", gmk, settings.GeminiModel)
	case "fred":
		result = testFREDProvider(ctx, fred)
	case "bls":
		result = testBLSProvider(ctx, bls)
	case "eia":
		result = testEIAProvider(ctx, eia)
	case "twelvedata":
		result = testTwelveDataProvider(ctx, twelve)
	case "marketaux":
		result = testMarketauxProvider(ctx, marketaux)
	case "options":
		result = testOptionsProvider(ctx, ak, as, settings.OptionsDataMode)
	default:
		result = testFinnhub(ctx, fk)
	}
	a.mu.Lock()
	if a.state.ProviderStatus == nil {
		a.state.ProviderStatus = map[string]ProviderTestResult{}
	}
	a.state.ProviderStatus[result.Provider] = result
	_ = a.saveLocked()
	a.mu.Unlock()
	a.broadcastSharedState()
	writeJSON(w, 200, result)
}
func testFinnhub(ctx context.Context, key string) ProviderTestResult {
	r := ProviderTestResult{Provider: "finnhub", CheckedAt: time.Now().UTC().Format(time.RFC3339)}
	if key == "" {
		r.Status = "missing"
		r.Message = "Enter a Finnhub API key."
		return r
	}
	var q map[string]any
	err := getJSON(ctx, &http.Client{Timeout: 12 * time.Second}, "https://finnhub.io/api/v1/quote?symbol=SPY&token="+url.QueryEscape(key), nil, &q)
	if err != nil {
		r.Status = "failed"
		r.Message = err.Error()
		return r
	}
	if toFloat(q["c"]) <= 0 {
		r.Status = "failed"
		r.Message = "Finnhub returned no SPY price."
		return r
	}
	r.OK = true
	r.Status = "connected"
	r.Message = "Finnhub quote access is working."
	r.Details = []string{"SPY quote authenticated", "WebSocket will be verified when the runtime starts"}
	return r
}
func testTradeInsight(ctx context.Context, key string) ProviderTestResult {
	r := ProviderTestResult{Provider: "tradeinsight", CheckedAt: time.Now().UTC().Format(time.RFC3339)}
	key = cleanCredential(key)
	if key == "" {
		r.Status = "missing"
		r.Message = "Enter a TradeInsight API key."
		return r
	}
	end := time.Now().UTC()
	params := url.Values{
		"ticker":        []string{"SPY"},
		"start":         []string{end.AddDate(0, 0, -10).Format("2006-01-02")},
		"end":           []string{end.AddDate(0, 0, 1).Format("2006-01-02")},
		"adjust_volume": []string{"true"},
	}
	rows, err := tradeInsightFetchRowsAt(ctx, &http.Client{Timeout: 12 * time.Second}, tradeInsightRESTBaseURL, key, "/ohlc", params)
	if err != nil {
		r.Status = "failed"
		r.Message = err.Error()
		return r
	}
	r.OK = true
	r.Status = "connected"
	r.Message = "TradeInsight historical-data access is working."
	r.Details = []string{fmt.Sprintf("%d SPY OHLC rows returned", len(rows)), "Shadow-first capability admission remains enforced"}
	return r
}

func testAlpaca(ctx context.Context, key, secret string) ProviderTestResult {
	r := ProviderTestResult{Provider: "alpaca", CheckedAt: time.Now().UTC().Format(time.RFC3339)}
	if key == "" || secret == "" {
		r.Status = "missing"
		r.Message = "Enter both Alpaca Key ID and Secret Key."
		return r
	}
	raw := "https://data.alpaca.markets/v2/stocks/SPY/bars?timeframe=1Day&limit=2&adjustment=all&feed=iex"
	var p map[string]any
	err := getJSON(ctx, &http.Client{Timeout: 12 * time.Second}, raw, map[string]string{"APCA-API-KEY-ID": key, "APCA-API-SECRET-KEY": secret}, &p)
	if err != nil {
		r.Status = "failed"
		r.Message = err.Error()
		return r
	}
	r.OK = true
	r.Status = "connected"
	r.Message = "Alpaca historical-bar access is working."
	r.Details = []string{"SPY daily bars authenticated", "Used for intraday, daily and weekly history"}
	return r
}
func testAIProvider(ctx context.Context, provider, key, model string) ProviderTestResult {
	r := ProviderTestResult{Provider: provider, CheckedAt: time.Now().UTC().Format(time.RFC3339)}
	if strings.TrimSpace(key) == "" {
		r.Status = "missing"
		r.Message = "Enter an API key."
		return r
	}
	provider = strings.ToLower(strings.TrimSpace(provider))
	if model == "" {
		if provider == "groq" {
			model = "openai/gpt-oss-120b"
		} else {
			model = "gemini-3.1-flash-lite"
		}
	}
	var text string
	var err error
	switch provider {
	case "groq":
		text, err = generateOpenAICompatibleResponse(ctx, "https://api.groq.com/openai/v1/responses", key, model, "Connectivity test only.", "Reply with exactly: OK", 32, false)
	case "gemini":
		text, err = generateGeminiResponse(ctx, key, model, "Connectivity test only.", "Reply with exactly: OK", 32)
	default:
		err = errors.New("unsupported AI provider")
	}
	if err != nil {
		r.Status = "failed"
		msg := strings.TrimSpace(err.Error())
		lower := strings.ToLower(msg)
		if strings.Contains(lower, "rate limit") || strings.Contains(lower, "resource_exhausted") || strings.Contains(lower, "429") {
			r.Status = "rate limited"
			r.Message = "API key reached its current provider rate limit."
		} else {
			r.Message = msg
		}
		return r
	}
	if strings.TrimSpace(text) == "" {
		r.Status = "failed"
		r.Message = "Provider returned no text."
		return r
	}
	r.OK = true
	r.Status = "connected"
	if provider == "groq" {
		r.Message = "Groq AI is ready."
	} else {
		r.Message = "Google Gemini is ready."
	}
	r.Details = []string{"Model: " + model}
	return r
}

func testOpenRouter(ctx context.Context, key, mode, specific string) ProviderTestResult {
	r := ProviderTestResult{Provider: "openrouter", CheckedAt: time.Now().UTC().Format(time.RFC3339)}
	if strings.TrimSpace(key) == "" {
		r.Status = "missing"
		r.Message = "Enter an OpenRouter API key."
		return r
	}
	result, err := generateOpenRouterResponse(ctx, key, mode, specific, "Connectivity test only.", "Reply with exactly: OK", 32)
	if err != nil {
		r.Status = "failed"
		msg := strings.TrimSpace(err.Error())
		lower := strings.ToLower(msg)
		if strings.Contains(lower, "credit") || strings.Contains(lower, "payment") || strings.Contains(lower, "insufficient") || strings.Contains(lower, "402") {
			r.Status = "billing required"
			r.Message = "OpenRouter key is valid, but credits may be required for the selected models."
		} else if strings.Contains(lower, "rate limit") || strings.Contains(lower, "429") {
			r.Status = "rate limited"
			r.Message = "OpenRouter reached its current rate limit."
		} else {
			r.Message = msg
		}
		return r
	}
	if strings.TrimSpace(result.Text) == "" {
		r.Status = "failed"
		r.Message = "OpenRouter returned no text."
		return r
	}
	r.OK = true
	r.Status = "connected"
	r.Message = "OpenRouter is ready."
	r.Details = []string{"Mode: " + strings.ToUpper(result.Mode), "Model: " + result.ActualModel}
	return r
}

func (a *Application) handleCacheClear(w http.ResponseWriter, r *http.Request) {
	before := int64(0)
	if st, err := os.Stat(a.cachePath()); err == nil {
		before = st.Size()
	}
	a.engine.mu.Lock()
	a.engine.quotes = map[string]Quote{}
	a.engine.history = map[string][]HistoryPoint{}
	a.engine.bars = map[string]map[string][]Bar{}
	a.engine.fundamentals = map[string]FundamentalSnapshot{}
	a.engine.news = []NewsItem{}
	a.engine.earnings = []EarningsItem{}
	a.engine.filings = []FilingItem{}
	a.engine.secIntelligence = map[string]SECIntelligenceSummary{}
	a.engine.scanner = ScannerState{Mode: "day", Status: "idle", Message: "Run a scan to discover candidates.", Results: []ScannerResult{}}
	a.engine.lastUpdated = map[string]int64{}
	a.engine.lastAlpacaAt = 0
	a.engine.lastAlpacaSymbol = ""
	a.engine.lastAlpacaFeed = ""
	a.engine.lastTradeAt = 0
	a.engine.lastTradeSymbol = ""
	a.engine.mu.Unlock()
	err := os.Remove(a.cachePath())
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		writeError(w, 500, "Unable to clear market cache: "+err.Error())
		return
	}
	clearedAt := time.Now().UnixMilli()
	a.mu.Lock()
	a.state.LastCacheCleared = clearedAt
	_ = a.saveLocked()
	a.mu.Unlock()
	a.broadcastRuntime()

	a.engine.requestHistoryHydration()
	a.broadcastSharedState()
	writeJSON(w, 200, map[string]any{"ok": true, "cachePath": a.cachePath(), "bytesBefore": before, "bytesAfter": 0, "clearedAt": clearedAt, "message": "Market cache cleared successfully. Required historical bars are rehydrating now while the runtime is active."})
}

func (a *Application) handleCacheRefresh(w http.ResponseWriter, r *http.Request) {
	a.engine.mu.RLock()
	last := a.engine.lastUpdated["cache-refresh"]
	status, mode := a.engine.status, a.engine.mode
	a.engine.mu.RUnlock()
	if status != "running" && status != "degraded" {
		writeError(w, 409, "Start the runtime before refreshing stale market data.")
		return
	}
	if mode != "live" {
		a.engine.setManualAction("refresh-due", "Refresh Due Data", "COMPLETE", "Demo mode uses generated local data; no provider refresh is required.", true)
		writeJSON(w, 200, map[string]any{"ok": true, "message": "Demo mode uses generated local data; no provider refresh is required."})
		return
	}
	if last > 0 && time.Since(time.UnixMilli(last)) < 2*time.Minute {
		a.engine.setManualAction("refresh-due", "Refresh Due Data", "COMPLETE", "Current cache is already fresh enough to reuse.", true)
		writeJSON(w, 200, map[string]any{"ok": true, "message": "A low-priority refresh completed recently; the current cache is already being reused."})
		return
	}
	a.mu.RLock()
	finnhubKey := strings.TrimSpace(a.secrets.Finnhub)
	alpacaKey := strings.TrimSpace(a.secrets.AlpacaKey)
	alpacaSecret := strings.TrimSpace(a.secrets.AlpacaSecret)
	a.mu.RUnlock()
	a.engine.setManualAction("refresh-due", "Refresh Due Data", "RUNNING", "Refreshing due/stale provider data without clearing good cache.", false)
	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), 4*time.Minute)
		defer cancel()
		a.engine.runLowPriorityRefresh(ctx, finnhubKey, alpacaKey, alpacaSecret, "manual")
		a.engine.mu.RLock()
		detail := a.engine.health["cache-refresh"]
		a.engine.mu.RUnlock()
		degraded := strings.Contains(strings.ToLower(detail), "degraded") || strings.Contains(strings.ToLower(detail), "failed")
		if degraded {
			a.engine.setManualAction("refresh-due", "Refresh Due Data", "DEGRADED", detail, false)
		} else {
			a.engine.setManualAction("refresh-due", "Refresh Due Data", "COMPLETE", detail, true)
		}
		a.broadcastRuntime()
	}()
	writeJSON(w, 202, map[string]any{"ok": true, "message": "Refreshing stale history, fundamentals, earnings, news and SEC data in the background. Live quote streams remain active."})
}

func (a *Application) handlePreMarketPrep(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed, "POST required")
		return
	}
	a.engine.mu.RLock()
	status, mode := a.engine.status, a.engine.mode
	a.engine.mu.RUnlock()
	if status != "running" && status != "degraded" {
		writeError(w, 409, "Start the runtime before running pre-market preparation.")
		return
	}
	a.mu.RLock()
	finnhubKey := strings.TrimSpace(a.secrets.Finnhub)
	alpacaKey := strings.TrimSpace(a.secrets.AlpacaKey)
	alpacaSecret := strings.TrimSpace(a.secrets.AlpacaSecret)
	a.mu.RUnlock()
	if mode != "live" {
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
		a.engine.runPreMarketPrep(ctx, finnhubKey, alpacaKey, alpacaSecret, "manual · demo context", false)
		writeJSON(w, 200, map[string]any{"ok": true, "message": "Pre-Market Prep completed using the current Demo context."})
		return
	}
	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), 4*time.Minute)
		defer cancel()
		a.engine.runPreMarketPrep(ctx, finnhubKey, alpacaKey, alpacaSecret, "manual", false)
	}()
	writeJSON(w, 202, map[string]any{"ok": true, "message": "Pre-market preparation started. Current data is reconciled without clearing good cache."})
}

func (a *Application) handleIntegrityCheck(w http.ResponseWriter, r *http.Request) {
	if a.engine == nil {
		writeError(w, 409, "Runtime engine is unavailable.")
		return
	}
	a.engine.setManualAction("integrity-check", "Weekly Integrity", "RUNNING", "Running non-destructive integrity checks.", false)
	a.engine.runWeeklyIntegrity("manual")
	a.engine.setManualAction("integrity-check", "Weekly Integrity", "COMPLETE", "Non-destructive cache integrity check completed. Watchlists and settings were not changed.", true)
	writeJSON(w, 200, map[string]any{"ok": true, "message": "Non-destructive cache integrity check completed. Watchlists and settings were not changed."})
}

func (a *Application) handleEngineToggle(w http.ResponseWriter, r *http.Request) {
	var in struct {
		Engine  string `json:"engine"`
		Enabled bool   `json:"enabled"`
	}
	if decodeJSON(r, &in) != nil {
		writeError(w, 400, "Invalid engine toggle")
		return
	}
	userID := requestUserID(r.Context())
	a.mu.Lock()
	switch strings.ToLower(strings.TrimSpace(in.Engine)) {
	case "day":
		a.state.Settings.DayEnabled = in.Enabled
	case "swing":
		a.state.Settings.SwingEnabled = in.Enabled
	case "long":
		a.state.Settings.LongEnabled = in.Enabled
	default:
		a.mu.Unlock()
		writeError(w, 400, "Unknown trading engine")
		return
	}
	_ = a.saveLocked()
	state := a.publicStateLockedForUser(userID)
	a.mu.Unlock()
	if a.engine != nil {
		a.engine.mu.RLock()
		ws := a.engine.ws
		running := a.engine.status == "running" || a.engine.status == "degraded"
		a.engine.mu.RUnlock()
		if running && ws != nil {
			a.engine.syncLiveSubscriptions(ws)
		}
		if running && in.Enabled {
			a.engine.requestHistoryHydration()
		}
	}
	a.broadcastSharedState()
	writeJSON(w, 200, state)
}

func (a *Application) handleRuntimeStart(w http.ResponseWriter, r *http.Request) {
	if err := a.engine.Start(); err != nil {
		writeError(w, 400, err.Error())
		return
	}
	writeJSON(w, 200, a.engine.Snapshot())
}
func (a *Application) handleRuntimeStop(w http.ResponseWriter, r *http.Request) {
	a.engine.Stop()
	writeJSON(w, 200, a.engine.Snapshot())
}
func (a *Application) handleLivePriority(w http.ResponseWriter, r *http.Request) {
	var in struct {
		Symbols []string `json:"symbols"`
	}
	if decodeJSON(r, &in) != nil {
		writeError(w, 400, "Invalid live-priority request")
		return
	}

	now := time.Now().UnixMilli()
	next := map[string]int64{}
	for _, raw := range in.Symbols {
		symbol, validSymbol := parseUserTicker(raw)
		if !validSymbol || len(next) >= 10 {
			continue
		}
		next[symbol] = now
	}
	a.engine.mu.Lock()
	a.engine.livePriorityHints = next
	ws := a.engine.ws
	alpacaWS := a.engine.alpacaWS
	a.engine.mu.Unlock()
	if ws != nil {
		a.engine.syncLiveSubscriptions(ws)
	}
	if alpacaWS != nil {
		a.engine.syncAlpacaSubscriptions(alpacaWS)
	}
	writeJSON(w, 200, map[string]any{"ok": true, "symbols": len(next)})
}

func (a *Application) handleScope(w http.ResponseWriter, r *http.Request) {
	var in struct {
		ScopeType   string `json:"scopeType"`
		WatchlistID string `json:"watchlistId"`
	}
	if decodeJSON(r, &in) != nil {
		writeError(w, 400, "Invalid scope")
		return
	}
	userID := requestUserID(r.Context())
	a.mu.Lock()
	workspaceState := a.workspaceStateLocked(userID)
	if in.ScopeType == "watchlist" {
		workspaceState.UI.ScopeType = "watchlist"
	} else {
		workspaceState.UI.ScopeType = "general"
	}
	if in.WatchlistID != "" && findWatchlistInState(&workspaceState, in.WatchlistID) != nil {
		workspaceState.UI.WatchlistID = in.WatchlistID
	}
	symbols := scopeSymbolsForState(&workspaceState)
	if !contains(symbols, workspaceState.UI.SelectedTicker) {
		if len(symbols) > 0 {
			workspaceState.UI.SelectedTicker = symbols[0]
		} else {
			workspaceState.UI.SelectedTicker = "SPY"
		}
	}
	if err := a.saveWorkspaceStateLocked(userID, workspaceState); err != nil {
		a.mu.Unlock()
		writeError(w, 500, err.Error())
		return
	}
	state := a.publicStateLockedForUser(userID)
	a.mu.Unlock()
	a.broadcastStateForUser(userID)
	writeJSON(w, 200, state)
}
func (a *Application) handleTicker(w http.ResponseWriter, r *http.Request) {
	var in struct {
		Symbol string `json:"symbol"`
	}
	if decodeJSON(r, &in) != nil {
		writeError(w, 400, "Invalid ticker")
		return
	}
	symbol, valid := parseSelectableTicker(in.Symbol)
	if !valid {
		writeError(w, 400, "Invalid ticker")
		return
	}
	userID := requestUserID(r.Context())
	a.mu.Lock()
	workspaceState := a.workspaceStateLocked(userID)
	changed := normalizeSymbol(workspaceState.UI.SelectedTicker) != symbol
	workspaceState.UI.SelectedTicker = symbol
	if err := a.saveWorkspaceStateLocked(userID, workspaceState); err != nil {
		a.mu.Unlock()
		writeError(w, 500, err.Error())
		return
	}
	state := a.publicStateLockedForUser(userID)
	a.mu.Unlock()
	if changed && a.engine != nil {
		a.engine.onSymbolSetChanged(symbol)
	}
	a.broadcastStateForUser(userID)
	writeJSON(w, 200, state)
}
func (a *Application) handleWatchlistCreate(w http.ResponseWriter, r *http.Request) {
	var in struct {
		Name string `json:"name"`
	}
	if decodeJSON(r, &in) != nil {
		writeError(w, 400, "Invalid watchlist request")
		return
	}
	name := strings.TrimSpace(in.Name)
	if name == "" {
		name = "New Watchlist"
	}
	wl := Watchlist{ID: randomID("watchlist"), Name: truncate(name, 60), Symbols: []string{}}
	userID := requestUserID(r.Context())
	a.mu.Lock()
	workspaceState := a.workspaceStateLocked(userID)
	workspaceState.Watchlists = append(workspaceState.Watchlists, wl)
	workspaceState.UI.WatchlistID = wl.ID
	workspaceState.UI.ScopeType = "watchlist"
	if err := a.saveWorkspaceStateLocked(userID, workspaceState); err != nil {
		a.mu.Unlock()
		writeError(w, 500, err.Error())
		return
	}
	a.mu.Unlock()
	a.broadcastStateForUser(userID)
	writeJSON(w, 200, wl)
}
func (a *Application) handleWatchlistRename(w http.ResponseWriter, r *http.Request) {
	var in struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	}
	if decodeJSON(r, &in) != nil {
		writeError(w, 400, "Invalid watchlist request")
		return
	}
	userID := requestUserID(r.Context())
	a.mu.Lock()
	workspaceState := a.workspaceStateLocked(userID)
	wl := findWatchlistInState(&workspaceState, in.ID)
	if wl == nil {
		a.mu.Unlock()
		writeError(w, 404, "Watchlist not found")
		return
	}
	if strings.TrimSpace(in.Name) != "" {
		wl.Name = truncate(strings.TrimSpace(in.Name), 60)
	}
	out := *wl
	if err := a.saveWorkspaceStateLocked(userID, workspaceState); err != nil {
		a.mu.Unlock()
		writeError(w, 500, err.Error())
		return
	}
	a.mu.Unlock()
	a.broadcastStateForUser(userID)
	writeJSON(w, 200, out)
}
func (a *Application) handleWatchlistDelete(w http.ResponseWriter, r *http.Request) {
	var in struct {
		ID string `json:"id"`
	}
	if decodeJSON(r, &in) != nil {
		writeError(w, 400, "Invalid watchlist request")
		return
	}
	userID := requestUserID(r.Context())
	a.mu.Lock()
	workspaceState := a.workspaceStateLocked(userID)
	if in.ID == "swing" || in.ID == "day" || in.ID == "long" || in.ID == "discovery" {
		a.mu.Unlock()
		writeError(w, 400, "Desk watchlists are permanent and are managed inside their trading desk")
		return
	}
	if len(workspaceState.Watchlists) <= 4 {
		a.mu.Unlock()
		writeError(w, 400, "Keep at least one watchlist")
		return
	}
	out := make([]Watchlist, 0, len(workspaceState.Watchlists)-1)
	for _, wl := range workspaceState.Watchlists {
		if wl.ID != in.ID {
			out = append(out, wl)
		}
	}
	workspaceState.Watchlists = out
	if findWatchlistInState(&workspaceState, workspaceState.UI.WatchlistID) == nil {
		workspaceState.UI.WatchlistID = workspaceState.Watchlists[0].ID
	}
	if err := a.saveWorkspaceStateLocked(userID, workspaceState); err != nil {
		a.mu.Unlock()
		writeError(w, 500, err.Error())
		return
	}
	state := a.publicStateLockedForUser(userID)
	a.mu.Unlock()
	a.broadcastStateForUser(userID)
	writeJSON(w, 200, state)
}
func (a *Application) handleAddSymbol(w http.ResponseWriter, r *http.Request) {
	var in struct {
		WatchlistID string `json:"watchlistId"`
		Symbol      string `json:"symbol"`
	}
	if decodeJSON(r, &in) != nil {
		writeError(w, 400, "Invalid watchlist request")
		return
	}
	symbol, validSymbol := parseUserTicker(in.Symbol)
	if !validSymbol {
		writeError(w, 400, "Enter a valid ticker symbol")
		return
	}
	userID := requestUserID(r.Context())
	a.mu.Lock()
	workspaceState := a.workspaceStateLocked(userID)
	wl := findWatchlistInState(&workspaceState, in.WatchlistID)
	if wl == nil {
		a.mu.Unlock()
		writeError(w, 404, "Watchlist not found")
		return
	}
	alreadyPresent := contains(wl.Symbols, symbol)
	if in.WatchlistID == "day" || in.WatchlistID == "swing" || in.WatchlistID == "long" {
		_, _, _ = applyDeskMembershipLocked(&workspaceState, in.WatchlistID, symbol, true)
		wl = findWatchlistInState(&workspaceState, in.WatchlistID)
	} else {
		wl.Symbols = uniqueSymbols(append(wl.Symbols, symbol))
	}
	workspaceState.UI.SelectedTicker = symbol
	out := *wl
	if err := a.saveWorkspaceStateLocked(userID, workspaceState); err != nil {
		a.mu.Unlock()
		writeError(w, 500, err.Error())
		return
	}
	a.mu.Unlock()
	a.broadcastStateForUser(userID)
	if !alreadyPresent && a.engine != nil {
		a.engine.onSymbolSetChanged(symbol)
	}
	writeJSON(w, 200, out)
}
func (a *Application) handleRemoveSymbol(w http.ResponseWriter, r *http.Request) {
	var in struct {
		WatchlistID string `json:"watchlistId"`
		Symbol      string `json:"symbol"`
	}
	if decodeJSON(r, &in) != nil {
		writeError(w, 400, "Invalid watchlist request")
		return
	}
	symbol, validSymbol := parseUserTicker(in.Symbol)
	if !validSymbol {
		writeError(w, 400, "Enter a valid ticker symbol")
		return
	}
	userID := requestUserID(r.Context())
	a.mu.Lock()
	workspaceState := a.workspaceStateLocked(userID)
	wl := findWatchlistInState(&workspaceState, in.WatchlistID)
	if wl == nil {
		a.mu.Unlock()
		writeError(w, 404, "Watchlist not found")
		return
	}
	if in.WatchlistID == "day" || in.WatchlistID == "swing" || in.WatchlistID == "long" {
		changed, protected, membership := applyDeskMembershipLocked(&workspaceState, in.WatchlistID, symbol, false)
		wl = findWatchlistInState(&workspaceState, in.WatchlistID)
		out := *wl
		if err := a.saveWorkspaceStateLocked(userID, workspaceState); err != nil {
			a.mu.Unlock()
			writeError(w, 500, err.Error())
			return
		}
		a.mu.Unlock()
		a.broadcastStateForUser(userID)
		if changed && a.engine != nil {
			a.engine.onSymbolSetChanged(symbol)
		}
		writeJSON(w, 200, map[string]any{"watchlist": out, "protected": protected, "changed": changed, "membership": membership})
		return
	}
	syms := make([]string, 0, len(wl.Symbols))
	for _, candidate := range wl.Symbols {
		if candidate != symbol {
			syms = append(syms, candidate)
		}
	}
	wl.Symbols = syms
	out := *wl
	if err := a.saveWorkspaceStateLocked(userID, workspaceState); err != nil {
		a.mu.Unlock()
		writeError(w, 500, err.Error())
		return
	}
	a.mu.Unlock()
	a.broadcastStateForUser(userID)
	if a.engine != nil {
		a.engine.onSymbolSetChanged(symbol)
	}
	writeJSON(w, 200, out)
}
