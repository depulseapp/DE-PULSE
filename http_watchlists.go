package main

import (
	"net/http"
	"strings"
)

// Watchlist HTTP mutation handlers are kept in this owner so http_api.go
// remains focused on shared API/runtime responsibilities and within budget.
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
