package main

import (
	"context"
	"errors"
	"sort"
	"strings"
	"time"
)

const userWorkspaceVersion = 1

// defaultUserWorkspace is intentionally empty of personal trading symbols.
// Shared market-context symbols (SPY/QQQ/etc.) remain owned by the canonical
// runtime core and therefore do not need to be copied into every user account.
func defaultUserWorkspace(userID string) UserWorkspace {
	st := defaultState()
	for i := range st.Watchlists {
		st.Watchlists[i].Symbols = []string{}
	}
	st.UI = UIState{ScopeType: "watchlist", WatchlistID: "swing", SelectedTicker: "SPY"}
	return workspaceFromState(userID, st, time.Now())
}

func workspaceFromState(userID string, st AppState, now time.Time) UserWorkspace {
	ensureDedicatedDeskWatchlists(&st, defaultEmptyWorkspaceState())
	if _, ok := parseSelectableTicker(st.UI.SelectedTicker); !ok {
		st.UI.SelectedTicker = "SPY"
	}
	return UserWorkspace{
		Version:    userWorkspaceVersion,
		UserID:     strings.TrimSpace(userID),
		Watchlists: clone(st.Watchlists),
		UI:         st.UI,
		UpdatedAt:  now.UnixMilli(),
	}
}

func defaultEmptyWorkspaceState() AppState {
	st := defaultState()
	for i := range st.Watchlists {
		st.Watchlists[i].Symbols = []string{}
	}
	st.UI = UIState{ScopeType: "watchlist", WatchlistID: "swing", SelectedTicker: "SPY"}
	return st
}

func normalizeUserWorkspace(workspace UserWorkspace) UserWorkspace {
	userID := strings.TrimSpace(workspace.UserID)
	base := defaultEmptyWorkspaceState()
	if len(workspace.Watchlists) > 0 {
		base.Watchlists = clone(workspace.Watchlists)
	}
	if workspace.UI.ScopeType != "" || workspace.UI.WatchlistID != "" || workspace.UI.SelectedTicker != "" {
		base.UI = workspace.UI
	}
	ensureDedicatedDeskWatchlists(&base, defaultEmptyWorkspaceState())
	if base.UI.ScopeType != "general" {
		base.UI.ScopeType = "watchlist"
	}
	if base.UI.ScopeType == "watchlist" {
		if _, ok := watchlistValueByID(base.Watchlists, base.UI.WatchlistID); !ok {
			base.UI.WatchlistID = "swing"
		}
	}
	if selected, ok := parseSelectableTicker(base.UI.SelectedTicker); ok {
		base.UI.SelectedTicker = selected
	} else {
		base.UI.SelectedTicker = "SPY"
	}
	workspace.Version = userWorkspaceVersion
	workspace.UserID = userID
	workspace.Watchlists = base.Watchlists
	workspace.UI = base.UI
	if workspace.UpdatedAt <= 0 {
		workspace.UpdatedAt = time.Now().UnixMilli()
	}
	return workspace
}

func legacyOwnerWorkspace(userID string, legacy AppState) UserWorkspace {
	legacy = mergeState(legacy)
	return normalizeUserWorkspace(UserWorkspace{
		Version:    userWorkspaceVersion,
		UserID:     userID,
		Watchlists: clone(legacy.Watchlists),
		UI:         legacy.UI,
		UpdatedAt:  time.Now().UnixMilli(),
	})
}

// initializeUserWorkspaces migrates the pre-v18.1 single-user state exactly
// once. The legacy Stable owner's personal market state becomes the bootstrap
// owner's workspace. New users start empty instead of inheriting that data.
func (a *Application) initializeUserWorkspaces() error {
	if a.persistence == nil {
		return errors.New("workspace persistence unavailable")
	}
	loaded, err := a.persistence.LoadUserWorkspaces(context.Background())
	if err != nil {
		return err
	}
	a.workspaces = make(map[string]UserWorkspace, len(loaded)+1)
	for _, workspace := range loaded {
		workspace = normalizeUserWorkspace(workspace)
		if workspace.UserID == "" {
			continue
		}
		a.workspaces[workspace.UserID] = workspace
	}
	if _, ok := a.workspaces[bootstrapOwnerID]; !ok {
		workspace := legacyOwnerWorkspace(bootstrapOwnerID, a.state)
		if err := a.persistence.SaveUserWorkspace(context.Background(), workspace); err != nil {
			return err
		}
		a.workspaces[bootstrapOwnerID] = workspace
	}
	// Shared state must not retain a private copy of the owner's symbols. Keep
	// only the operational Settings/provider status and a neutral UI shell.
	neutral := defaultEmptyWorkspaceState()
	a.state.Watchlists = neutral.Watchlists
	a.state.UI = neutral.UI
	return a.saveLocked()
}

func (a *Application) ensureUserWorkspace(userID string) error {
	userID = strings.TrimSpace(userID)
	if userID == "" {
		return errors.New("workspace user id is required")
	}
	a.mu.Lock()
	defer a.mu.Unlock()
	if a.workspaces == nil {
		// Legacy direct-handler tests intentionally remain single-user.
		if a.identity == nil {
			return nil
		}
		a.workspaces = map[string]UserWorkspace{}
	}
	if _, ok := a.workspaces[userID]; ok {
		return nil
	}
	workspace := defaultUserWorkspace(userID)
	a.workspaces[userID] = workspace
	if a.persistence != nil {
		if err := a.persistence.SaveUserWorkspace(context.Background(), workspace); err != nil {
			delete(a.workspaces, userID)
			return err
		}
		a.persistence.EnqueueSymbols(symbolRegistryRecords(a.processingStateLocked(), time.Now()))
	}
	return nil
}

func (a *Application) workspaceStateLocked(userID string) AppState {
	state := clone(a.state)
	if a.workspaces == nil {
		return state
	}
	workspace, ok := a.workspaces[strings.TrimSpace(userID)]
	if !ok {
		workspace = defaultUserWorkspace(userID)
	}
	workspace = normalizeUserWorkspace(workspace)
	state.Watchlists = clone(workspace.Watchlists)
	state.UI = workspace.UI
	return state
}

func (a *Application) saveWorkspaceStateLocked(userID string, state AppState) error {
	// Legacy direct-handler fixtures predate IdentityService/workspace persistence.
	// Keep their single-user behavior without creating a production anonymous workspace.
	if a.workspaces == nil && a.persistence == nil {
		a.state.Watchlists = clone(state.Watchlists)
		a.state.UI = state.UI
		return a.saveLocked()
	}
	userID = strings.TrimSpace(userID)
	if userID == "" {
		return errors.New("workspace user id is required")
	}
	workspace := workspaceFromState(userID, state, time.Now())
	workspace = normalizeUserWorkspace(workspace)
	if a.workspaces == nil {
		a.workspaces = map[string]UserWorkspace{}
	}
	a.workspaces[userID] = workspace
	if a.persistence != nil {
		if err := a.persistence.SaveUserWorkspace(context.Background(), workspace); err != nil {
			return err
		}
		a.persistence.EnqueueSymbols(symbolRegistryRecords(a.processingStateLocked(), time.Now()))
	}
	return nil
}

// processingStateLocked builds one canonical market-processing universe from
// all user workspaces. It is a union, never a per-user provider pipeline.
func (a *Application) processingStateLocked() AppState {
	state := clone(a.state)
	if a.workspaces == nil {
		return state
	}
	// Preserve any in-memory legacy/test mutation as a compatibility input, but
	// production v18.1 initialization strips personal symbols from shared state.
	ensureDedicatedDeskWatchlists(&state, defaultEmptyWorkspaceState())
	byID := map[string]*Watchlist{}
	for i := range state.Watchlists {
		byID[state.Watchlists[i].ID] = &state.Watchlists[i]
	}
	ids := make([]string, 0, len(a.workspaces))
	for userID := range a.workspaces {
		ids = append(ids, userID)
	}
	sort.Strings(ids)
	for _, userID := range ids {
		workspace := normalizeUserWorkspace(a.workspaces[userID])
		for _, wl := range workspace.Watchlists {
			if target := byID[wl.ID]; target != nil {
				target.Symbols = uniqueSymbols(append(target.Symbols, wl.Symbols...))
			}
		}
		if selected, ok := parseUserTicker(workspace.UI.SelectedTicker); ok {
			if target := byID["discovery"]; target != nil {
				target.Symbols = uniqueSymbols(append(target.Symbols, selected))
			}
		}
	}
	return state
}

func requestUserID(rContext context.Context) string {
	if p, ok := principalFromContext(rContext); ok && strings.TrimSpace(p.UserID) != "" {
		return p.UserID
	}
	// Legacy direct-handler tests do not install an identity context. Treat them
	// as the bootstrap owner rather than inventing a second anonymous workspace.
	return bootstrapOwnerID
}

func findWatchlistInState(state *AppState, id string) *Watchlist {
	if state == nil {
		return nil
	}
	for i := range state.Watchlists {
		if state.Watchlists[i].ID == id {
			return &state.Watchlists[i]
		}
	}
	return nil
}

func scopeSymbolsForState(state *AppState) []string {
	if state == nil {
		return nil
	}
	if state.UI.ScopeType == "general" {
		return generalSymbols
	}
	if wl := findWatchlistInState(state, state.UI.WatchlistID); wl != nil {
		return wl.Symbols
	}
	return nil
}
