package main

import (
	"context"
	"net/http"
	"strings"
	"time"
)

func maxInt64(v ...int64) int64 {
	m := int64(0)
	for _, x := range v {
		if x > m {
			m = x
		}
	}
	return m
}

func (a *Application) cacheInfo() CacheInfo {
	a.mu.RLock()
	defer a.mu.RUnlock()
	return a.cacheInfoLocked()
}

// Canonical cross-desk membership state. Membership is stored once in the
// permanent Day/Swing/Long watchlists and rendered everywhere from that state.
func deskIDs() []string { return []string{"day", "swing", "long"} }
func deskMembershipsLocked(st *AppState, symbol string) map[string]bool {
	out := map[string]bool{"day": false, "swing": false, "long": false}
	symbol = normalizeSymbol(symbol)
	for _, id := range deskIDs() {
		if wl, ok := watchlistValueByID(st.Watchlists, id); ok {
			out[id] = contains(wl.Symbols, symbol)
		}
	}
	return out
}
func activeDeskCount(m map[string]bool) int {
	n := 0
	for _, id := range deskIDs() {
		if m[id] {
			n++
		}
	}
	return n
}

// trackedSymbolsLocked is the single canonical owner for the user-managed
// cross-desk symbol set. Day/Swing/Long watchlists remain durable membership
// views for compatibility, but cross-desk add/remove/clear operations must go
// through the helpers below rather than mutating each desk independently.
func trackedSymbolsLocked(st *AppState) []string {
	var out []string
	for _, id := range deskIDs() {
		if wl, ok := watchlistValueByID(st.Watchlists, id); ok {
			out = append(out, wl.Symbols...)
		}
	}
	return userTradingSymbols(out)
}

func setTrackedSymbolLocked(st *AppState, sym string, on bool) map[string]bool {
	sym = normalizeSymbol(sym)
	for _, id := range deskIDs() {
		setMembershipLocked(st, id, sym, on)
	}
	return deskMembershipsLocked(st, sym)
}

func clearTrackedSymbolsLocked(st *AppState) map[string]map[string]bool {
	removed := map[string]map[string]bool{}
	for _, sym := range trackedSymbolsLocked(st) {
		removed[sym] = deskMembershipsLocked(st, sym)
	}
	for _, id := range deskIDs() {
		wl, ok := watchlistValueByID(st.Watchlists, id)
		if !ok {
			continue
		}
		// Explicit empty slices are durable truth. A nil slice is reserved for
		// legacy/missing state and may legitimately hydrate defaults.
		wl.Symbols = []string{}
		for i := range st.Watchlists {
			if st.Watchlists[i].ID == id {
				st.Watchlists[i] = wl
				break
			}
		}
	}
	return removed
}
func setMembershipLocked(st *AppState, id, symbol string, on bool) {
	wl, ok := watchlistValueByID(st.Watchlists, id)
	if !ok {
		return
	}
	if on {
		wl.Symbols = uniqueSymbols(append(wl.Symbols, symbol))
	} else {
		out := wl.Symbols[:0]
		for _, s := range wl.Symbols {
			if normalizeSymbol(s) != symbol {
				out = append(out, s)
			}
		}
		// Preserve an explicit empty slice. A nil slice means "legacy/missing" to
		// ensureDedicatedDeskWatchlists and would repopulate defaults on the next
		// normalization pass, making the final Master Symbol removal reappear.
		wl.Symbols = append([]string{}, out...)
	}
	for i := range st.Watchlists {
		if st.Watchlists[i].ID == id {
			st.Watchlists[i] = wl
			break
		}
	}
}

func applyDeskMembershipLocked(st *AppState, desk, sym string, desired bool) (changed, protected bool, membership map[string]bool) {
	ensureDedicatedDeskWatchlists(st, defaultState())
	membership = deskMembershipsLocked(st, sym)
	if desired && !membership[desk] {
		setMembershipLocked(st, desk, sym, true)
		changed = true
	} else if !desired && membership[desk] {
		if activeDeskCount(membership) > 1 {
			setMembershipLocked(st, desk, sym, false)
			changed = true
		} else {
			protected = true
		}
	}
	membership = deskMembershipsLocked(st, sym)
	return
}

func (a *Application) handleDeskMembership(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, 405, "Method not allowed")
		return
	}
	var in struct {
		Symbol string `json:"symbol"`
		Desk   string `json:"desk"`
		Active *bool  `json:"active,omitempty"`
	}
	if decodeJSON(r, &in) != nil {
		writeError(w, 400, "Invalid request")
		return
	}
	sym, validSymbol := parseUserTicker(in.Symbol)
	desk := strings.ToLower(strings.TrimSpace(in.Desk))
	if !validSymbol || !(desk == "day" || desk == "swing" || desk == "long") {
		writeError(w, 400, "Invalid symbol or desk")
		return
	}
	userID := requestUserID(r.Context())
	a.mu.Lock()
	workspaceState := a.workspaceStateLocked(userID)
	m := deskMembershipsLocked(&workspaceState, sym)
	desired := !m[desk]
	if in.Active != nil {
		desired = *in.Active
	}
	changed, protected, membership := applyDeskMembershipLocked(&workspaceState, desk, sym, desired)
	workspaceState.UI.SelectedTicker = sym
	if err := a.saveWorkspaceStateLocked(userID, workspaceState); err != nil {
		a.mu.Unlock()
		writeError(w, 500, err.Error())
		return
	}
	state := a.publicStateLockedForUser(userID)
	a.mu.Unlock()
	a.broadcastStateForUser(userID)
	if changed && a.engine != nil {
		a.engine.onSymbolSetChanged(sym)
	}
	writeJSON(w, 200, map[string]any{"symbol": sym, "membership": membership, "changed": changed, "protected": protected, "state": state})
}

var foreignListingSuffixes = []string{
	".TO", ".V", ".CN", ".NE", ".L", ".AX", ".HK", ".T", ".KS", ".KQ",
	".SI", ".NS", ".BO", ".PA", ".DE", ".F", ".MI", ".AS", ".BR", ".SW",
	".ST", ".OL", ".CO", ".HE", ".IC", ".IS", ".SA", ".MX",
}

func hasForeignListingSuffix(sym string) bool {
	upper := strings.ToUpper(strings.TrimSpace(sym))
	for _, suffix := range foreignListingSuffixes {
		if strings.HasSuffix(upper, suffix) {
			return true
		}
	}
	return false
}

func validUserTicker(sym string) bool {
	if sym == "" || sym == "VIX" || len(sym) > 8 || hasForeignListingSuffix(sym) {
		return false
	}
	for i, r := range sym {
		if (r >= 'A' && r <= 'Z') || (i > 0 && r >= '0' && r <= '9') || (i > 0 && (r == '.' || r == '-')) {
			continue
		}
		return false
	}
	return true
}

// parseUserTicker validates the user's literal stock/ETF input instead of first
// stripping unsupported characters. This prevents malformed input such as
// "AA/PL" from silently becoming a different valid-looking ticker ("AAPL").
// VIX remains a dedicated index-context feed and is intentionally excluded from
// user trading-desk membership.
func parseUserTicker(raw string) (string, bool) {
	sym := strings.ToUpper(strings.TrimSpace(raw))
	if !validUserTicker(sym) {
		return "", false
	}

	if normalizeSymbol(sym) != sym {
		return "", false
	}
	return sym, true
}

func yahooMetaIsUSActionable(meta struct {
	Symbol             string  `json:"symbol"`
	RegularMarketPrice float64 `json:"regularMarketPrice"`
	PreviousClose      float64 `json:"chartPreviousClose"`
	RegularMarketTime  int64   `json:"regularMarketTime"`
	ExchangeName       string  `json:"exchangeName"`
	FullExchangeName   string  `json:"fullExchangeName"`
	InstrumentType     string  `json:"instrumentType"`
	Currency           string  `json:"currency"`
}) bool {
	ex := strings.ToUpper(strings.TrimSpace(meta.ExchangeName + " " + meta.FullExchangeName))
	instrument := strings.ToUpper(strings.TrimSpace(meta.InstrumentType))
	if instrument != "" && instrument != "EQUITY" && instrument != "ETF" {
		return false
	}
	if meta.Currency != "" && !strings.EqualFold(meta.Currency, "USD") {
		return false
	}
	for _, token := range []string{"NASDAQ", "NYSE", "NYSEARCA", "ARCA", "AMEX", "NYSE AMERICAN", "CBOE", "BATS", "NMS", "NGM", "NCM", "NYQ", "PCX", "ASE", "BTS"} {
		if strings.Contains(ex, token) {
			return true
		}
	}
	return false
}

func parseSelectableTicker(raw string) (string, bool) {
	sym := strings.ToUpper(strings.TrimSpace(raw))
	if sym == "VIX" {
		return sym, true
	}
	return parseUserTicker(sym)
}

func (a *Application) handleMasterSymbolAdd(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, 405, "Method not allowed")
		return
	}
	var in struct {
		Symbol string `json:"symbol"`
	}
	if decodeJSON(r, &in) != nil {
		writeError(w, 400, "Invalid request")
		return
	}
	sym, validSymbol := parseUserTicker(in.Symbol)
	if !validSymbol {
		writeError(w, 400, "Enter a valid tradable stock/ETF ticker. VIX remains a dedicated index context feed.")
		return
	}
	userID := requestUserID(r.Context())
	a.mu.RLock()
	mode := a.state.Settings.DataMode
	workspaceState := a.workspaceStateLocked(userID)
	alreadyTracked := false
	for _, candidate := range trackedSymbolsLocked(&workspaceState) {
		if normalizeSymbol(candidate) == sym {
			alreadyTracked = true
			break
		}
	}
	a.mu.RUnlock()
	if !alreadyTracked && mode == "live" {
		ctx, cancel := context.WithTimeout(r.Context(), 12*time.Second)
		defer cancel()
		chart, err := fetchYahooChart(ctx, sym, "5d", "1d")
		if err != nil || len(chart.Chart.Result) == 0 || (chart.Chart.Result[0].Meta.RegularMarketPrice <= 0 && len(chart.Chart.Result[0].Timestamp) == 0) {
			writeError(w, 422, "Ticker could not be validated by the current recovery provider. Nothing was added.")
			return
		}
		if !yahooMetaIsUSActionable(chart.Chart.Result[0].Meta) {
			writeError(w, 422, "DE.PULSE actionable ticker processing is limited to U.S.-listed stocks and ETFs. Nothing was added.")
			return
		}
	}

	a.mu.Lock()
	workspaceState = a.workspaceStateLocked(userID)
	ensureDedicatedDeskWatchlists(&workspaceState, defaultEmptyWorkspaceState())
	before := deskMembershipsLocked(&workspaceState, sym)
	membership := setTrackedSymbolLocked(&workspaceState, sym, true)
	changed := activeDeskCount(before) != activeDeskCount(membership)
	workspaceState.UI.SelectedTicker = sym
	if err := a.saveWorkspaceStateLocked(userID, workspaceState); err != nil {
		a.mu.Unlock()
		writeError(w, 500, err.Error())
		return
	}
	state := a.publicStateLockedForUser(userID)
	a.mu.Unlock()
	a.broadcastStateForUser(userID)
	if a.engine != nil {
		a.engine.onSymbolSetChanged(sym)
	}
	message := "Tracked across Day, Swing and Long-Term."
	if !changed {
		message = "Already tracked across Day, Swing and Long-Term."
	} else if activeDeskCount(before) > 0 {
		message = "Missing desk memberships repaired; ticker is now tracked across Day, Swing and Long-Term."
	}
	writeJSON(w, 200, map[string]any{"ok": true, "symbol": sym, "changed": changed, "membership": membership, "message": message, "state": state})
}
func (a *Application) handleMasterSymbolRemove(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, 405, "Method not allowed")
		return
	}
	var in struct {
		Symbol string `json:"symbol"`
	}
	if decodeJSON(r, &in) != nil {
		writeError(w, 400, "Invalid request")
		return
	}
	sym, validSymbol := parseUserTicker(in.Symbol)
	if !validSymbol {
		writeError(w, 400, "Invalid symbol")
		return
	}
	userID := requestUserID(r.Context())
	a.mu.Lock()
	workspaceState := a.workspaceStateLocked(userID)
	ensureDedicatedDeskWatchlists(&workspaceState, defaultEmptyWorkspaceState())
	before := deskMembershipsLocked(&workspaceState, sym)
	setTrackedSymbolLocked(&workspaceState, sym, false)
	if normalizeSymbol(workspaceState.UI.SelectedTicker) == sym {
		replacement := ""
		for _, candidate := range activeDeskSymbolsFromState(workspaceState) {
			if candidate != sym {
				replacement = candidate
				break
			}
		}
		if replacement == "" {
			for _, candidate := range discoverySymbolsFromState(workspaceState) {
				if candidate != sym {
					replacement = candidate
					break
				}
			}
		}
		if replacement == "" {
			replacement = "SPY"
		}
		workspaceState.UI.SelectedTicker = replacement
	}
	if err := a.saveWorkspaceStateLocked(userID, workspaceState); err != nil {
		a.mu.Unlock()
		writeError(w, 500, err.Error())
		return
	}
	state := a.publicStateLockedForUser(userID)
	a.mu.Unlock()
	a.broadcastStateForUser(userID)
	if a.engine != nil {
		a.engine.onSymbolSetChanged(sym)
	}
	writeJSON(w, 200, map[string]any{"symbol": sym, "removed": before, "state": state})
}
func (a *Application) handleMasterSymbolRemoveAll(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, 405, "Method not allowed")
		return
	}
	userID := requestUserID(r.Context())
	a.mu.Lock()
	workspaceState := a.workspaceStateLocked(userID)
	ensureDedicatedDeskWatchlists(&workspaceState, defaultEmptyWorkspaceState())
	removed := clearTrackedSymbolsLocked(&workspaceState)
	if _, ok := parseUserTicker(workspaceState.UI.SelectedTicker); !ok || len(removed) > 0 {
		workspaceState.UI.SelectedTicker = "SPY"
	}
	if err := a.saveWorkspaceStateLocked(userID, workspaceState); err != nil {
		a.mu.Unlock()
		writeError(w, 500, err.Error())
		return
	}
	state := a.publicStateLockedForUser(userID)
	a.mu.Unlock()
	a.broadcastStateForUser(userID)
	if a.engine != nil {
		a.engine.onSymbolSetChanged("")
	}
	writeJSON(w, 200, map[string]any{"ok": true, "removedCount": len(removed), "removed": removed, "state": state})
}
func (a *Application) handleMasterSymbolRestore(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, 405, "Method not allowed")
		return
	}
	var in struct {
		Symbol     string          `json:"symbol"`
		Membership map[string]bool `json:"membership"`
	}
	if decodeJSON(r, &in) != nil {
		writeError(w, 400, "Invalid request")
		return
	}
	sym, validSymbol := parseUserTicker(in.Symbol)
	if !validSymbol {
		writeError(w, 400, "Invalid symbol")
		return
	}
	userID := requestUserID(r.Context())
	a.mu.Lock()
	workspaceState := a.workspaceStateLocked(userID)
	ensureDedicatedDeskWatchlists(&workspaceState, defaultEmptyWorkspaceState())
	for _, id := range deskIDs() {
		if in.Membership[id] {
			setMembershipLocked(&workspaceState, id, sym, true)
		}
	}
	workspaceState.UI.SelectedTicker = sym
	if err := a.saveWorkspaceStateLocked(userID, workspaceState); err != nil {
		a.mu.Unlock()
		writeError(w, 500, err.Error())
		return
	}
	state := a.publicStateLockedForUser(userID)
	a.mu.Unlock()
	a.broadcastStateForUser(userID)
	if a.engine != nil {
		a.engine.onSymbolSetChanged(sym)
	}
	writeJSON(w, 200, map[string]any{"ok": true, "symbol": sym, "membership": in.Membership, "state": state})
}

// Generic watchlist mutations belong with the canonical desk/watchlist membership owner.
// The HTTP route contract remains unchanged; this split keeps shared http_api.go within its responsibility budget.
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
