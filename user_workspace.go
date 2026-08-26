package main

import (
	"context"
	"errors"
	"net/http"
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

const (
	accountDeletionTombstoneUsername = "__deleted_account__"
	accountDeletionReasonUserRequest = "USER_REQUESTED"
)

type AccountDeletionTombstone struct {
	UserID    string `json:"userId"`
	TenantID  string `json:"tenantId"`
	DeletedAt int64  `json:"deletedAt"`
	Reason    string `json:"reason"`
}

type AccountPrivacyProfile struct {
	ID          string     `json:"id"`
	TenantID    string     `json:"tenantId"`
	Username    string     `json:"username"`
	DisplayName string     `json:"displayName,omitempty"`
	Role        UserRole   `json:"role"`
	Status      UserStatus `json:"status"`
	CreatedAt   int64      `json:"createdAt"`
	UpdatedAt   int64      `json:"updatedAt"`
	LastLoginAt int64      `json:"lastLoginAt,omitempty"`
}

type AccountPrivacyDevice struct {
	ID         string       `json:"id"`
	Label      string       `json:"label,omitempty"`
	Status     DeviceStatus `json:"status"`
	CreatedAt  int64        `json:"createdAt"`
	LastSeenAt int64        `json:"lastSeenAt"`
	RevokedAt  int64        `json:"revokedAt,omitempty"`
}

type AccountPrivacySession struct {
	ID                string `json:"id"`
	DeviceID          string `json:"deviceId,omitempty"`
	CreatedAt         int64  `json:"createdAt"`
	AuthenticatedAt   int64  `json:"authenticatedAt,omitempty"`
	MFAVerifiedAt     int64  `json:"mfaVerifiedAt,omitempty"`
	LastSeenAt        int64  `json:"lastSeenAt"`
	IdleExpiresAt     int64  `json:"idleExpiresAt"`
	AbsoluteExpiresAt int64  `json:"absoluteExpiresAt"`
	RevokedAt         int64  `json:"revokedAt,omitempty"`
}

type AccountHostedSyncMetadata struct {
	Devices  []AccountPrivacyDevice  `json:"devices"`
	Sessions []AccountPrivacySession `json:"sessions"`
}

type AccountDataExport struct {
	SchemaVersion      int                       `json:"schemaVersion"`
	ExportedAt         int64                     `json:"exportedAt"`
	Profile            AccountPrivacyProfile     `json:"profile"`
	Preferences        UIState                   `json:"preferences"`
	Watchlists         []Watchlist               `json:"watchlists"`
	ResearchMetadata   []map[string]any          `json:"researchMetadata"`
	Feedback           []map[string]any          `json:"feedback"`
	HostedSyncMetadata AccountHostedSyncMetadata `json:"hostedSyncMetadata"`
	CategoryStatus     map[string]string         `json:"categoryStatus"`
}

func normalizeAccountDeletionReason(raw string) (string, error) {
	reason := strings.ToUpper(strings.TrimSpace(raw))
	if reason == "" || reason == accountDeletionReasonUserRequest {
		return accountDeletionReasonUserRequest, nil
	}
	return "", errors.New("unsupported account deletion reason")
}

func isAccountDeletionTombstone(user UserRecord) bool {
	return user.Username == accountDeletionTombstoneUsername && user.Status == UserDisabled && user.Role == "" && user.PasswordHash == "" && !user.MustSetPassword
}

func accountDeletionTombstoneFromUser(user UserRecord) (AccountDeletionTombstone, bool) {
	if !isAccountDeletionTombstone(user) {
		return AccountDeletionTombstone{}, false
	}
	return AccountDeletionTombstone{UserID: user.ID, TenantID: normalizedTenantID(user.TenantID), DeletedAt: user.UpdatedAt, Reason: user.DisplayName}, true
}

func (s *IdentityService) accountPrivacyExport(userID string) (AccountPrivacyProfile, AccountHostedSyncMetadata, error) {
	if s == nil {
		return AccountPrivacyProfile{}, AccountHostedSyncMetadata{}, errors.New("identity unavailable")
	}
	userID = strings.TrimSpace(userID)
	s.mu.Lock()
	defer s.mu.Unlock()
	var user UserRecord
	found := false
	for _, candidate := range s.state.Users {
		if candidate.ID == userID && !isAccountDeletionTombstone(candidate) {
			user = candidate
			found = true
			break
		}
	}
	if !found || user.Status != UserActive {
		return AccountPrivacyProfile{}, AccountHostedSyncMetadata{}, errors.New("account unavailable")
	}
	profile := AccountPrivacyProfile{ID: user.ID, TenantID: normalizedTenantID(user.TenantID), Username: user.Username, DisplayName: user.DisplayName, Role: user.Role, Status: user.Status, CreatedAt: user.CreatedAt, UpdatedAt: user.UpdatedAt, LastLoginAt: user.LastLoginAt}
	metadata := AccountHostedSyncMetadata{Devices: []AccountPrivacyDevice{}, Sessions: []AccountPrivacySession{}}
	for _, device := range s.state.Devices {
		if device.UserID != userID {
			continue
		}
		metadata.Devices = append(metadata.Devices, AccountPrivacyDevice{ID: device.ID, Label: device.Label, Status: device.Status, CreatedAt: device.CreatedAt, LastSeenAt: device.LastSeenAt, RevokedAt: device.RevokedAt})
	}
	for _, session := range s.state.Sessions {
		if session.UserID != userID {
			continue
		}
		metadata.Sessions = append(metadata.Sessions, AccountPrivacySession{ID: session.ID, DeviceID: session.DeviceID, CreatedAt: session.CreatedAt, AuthenticatedAt: session.AuthenticatedAt, MFAVerifiedAt: session.MFAVerifiedAt, LastSeenAt: session.LastSeenAt, IdleExpiresAt: session.IdleExpiresAt, AbsoluteExpiresAt: session.AbsoluteExpiresAt, RevokedAt: session.RevokedAt})
	}
	return profile, metadata, nil
}

func (s *IdentityService) deleteAccount(userID, rawReason string) (AccountDeletionTombstone, error) {
	if s == nil {
		return AccountDeletionTombstone{}, errors.New("identity unavailable")
	}
	reason, err := normalizeAccountDeletionReason(rawReason)
	if err != nil {
		return AccountDeletionTombstone{}, err
	}
	userID = strings.TrimSpace(userID)
	if userID == "" {
		return AccountDeletionTombstone{}, errors.New("account user id is required")
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	idx := -1
	var user UserRecord
	for i := range s.state.Users {
		if s.state.Users[i].ID == userID {
			idx = i
			user = s.state.Users[i]
			break
		}
	}
	if idx < 0 || isAccountDeletionTombstone(user) || user.Status != UserActive {
		return AccountDeletionTombstone{}, errors.New("account unavailable")
	}
	if (user.Role == RoleOwner || user.Role == RoleSuperOwner) && s.activeCriticalOwnersLocked(userID, "", UserDisabled) == 0 {
		return AccountDeletionTombstone{}, errors.New("cannot delete the last active owner")
	}
	now := s.now().UnixMilli()
	previousUsers := append([]UserRecord(nil), s.state.Users...)
	previousDevices := append([]DeviceRecord(nil), s.state.Devices...)
	previousSessions := append([]SessionRecord(nil), s.state.Sessions...)
	s.state.Users[idx] = UserRecord{ID: user.ID, TenantID: normalizedTenantID(user.TenantID), Username: accountDeletionTombstoneUsername, DisplayName: reason, Status: UserDisabled, UpdatedAt: now}
	devices := make([]DeviceRecord, 0, len(s.state.Devices))
	for _, device := range s.state.Devices {
		if device.UserID != userID {
			devices = append(devices, device)
		}
	}
	sessions := make([]SessionRecord, 0, len(s.state.Sessions))
	for _, session := range s.state.Sessions {
		if session.UserID != userID {
			sessions = append(sessions, session)
		}
	}
	s.state.Devices = devices
	s.state.Sessions = sessions
	if err := s.persistLocked(); err != nil {
		s.state.Users = previousUsers
		s.state.Devices = previousDevices
		s.state.Sessions = previousSessions
		return AccountDeletionTombstone{}, err
	}
	tombstone, _ := accountDeletionTombstoneFromUser(s.state.Users[idx])
	return tombstone, nil
}

func privacyBlankWorkspace(userID string, at int64) UserWorkspace {
	if at <= 0 {
		at = time.Now().UnixMilli()
	}
	return workspaceFromState(userID, defaultEmptyWorkspaceState(), time.UnixMilli(at))
}

func (a *Application) accountDataExport(userID string) (AccountDataExport, error) {
	if a == nil || a.identity == nil {
		return AccountDataExport{}, errors.New("identity unavailable")
	}
	profile, hosted, err := a.identity.accountPrivacyExport(userID)
	if err != nil {
		return AccountDataExport{}, err
	}
	a.mu.RLock()
	workspace, ok := a.workspaces[strings.TrimSpace(userID)]
	if !ok {
		workspace = defaultUserWorkspace(userID)
	}
	workspace = normalizeUserWorkspace(workspace)
	a.mu.RUnlock()
	return AccountDataExport{SchemaVersion: 1, ExportedAt: time.Now().UTC().UnixMilli(), Profile: profile, Preferences: workspace.UI, Watchlists: clone(workspace.Watchlists), ResearchMetadata: []map[string]any{}, Feedback: []map[string]any{}, HostedSyncMetadata: hosted, CategoryStatus: map[string]string{"researchMetadata": "not-persisted-per-user", "feedback": "not-persisted-per-user"}}, nil
}

func (a *Application) deleteAccountData(userID, rawReason string) (AccountDeletionTombstone, error) {
	if a == nil || a.identity == nil || a.persistence == nil {
		return AccountDeletionTombstone{}, errors.New("durable account deletion unavailable")
	}
	userID = strings.TrimSpace(userID)
	if userID == "" {
		return AccountDeletionTombstone{}, errors.New("account user id is required")
	}
	reason, err := normalizeAccountDeletionReason(rawReason)
	if err != nil {
		return AccountDeletionTombstone{}, err
	}
	a.mu.Lock()
	if a.workspaces == nil {
		a.workspaces = map[string]UserWorkspace{}
	}
	previous, existed := a.workspaces[userID]
	if !existed {
		previous = defaultUserWorkspace(userID)
	}
	blank := privacyBlankWorkspace(userID, time.Now().UnixMilli())
	a.workspaces[userID] = blank
	if err := a.persistence.SaveUserWorkspace(context.Background(), blank); err != nil {
		if existed {
			a.workspaces[userID] = previous
		} else {
			delete(a.workspaces, userID)
		}
		a.mu.Unlock()
		return AccountDeletionTombstone{}, err
	}
	a.mu.Unlock()

	tombstone, err := a.identity.deleteAccount(userID, reason)
	if err != nil {
		restoreErr := a.persistence.SaveUserWorkspace(context.Background(), previous)
		a.mu.Lock()
		if existed {
			a.workspaces[userID] = previous
		} else {
			delete(a.workspaces, userID)
		}
		a.mu.Unlock()
		if restoreErr != nil {
			return AccountDeletionTombstone{}, errors.New("account deletion failed and workspace rollback was unavailable")
		}
		return AccountDeletionTombstone{}, err
	}
	a.mu.RLock()
	processing := a.processingStateLocked()
	a.mu.RUnlock()
	a.persistence.EnqueueSymbols(symbolRegistryRecords(processing, time.Now()))
	return tombstone, nil
}

func accountDeletionTombstones(state IdentityPersistentState) map[string]UserRecord {
	out := map[string]UserRecord{}
	for _, user := range state.Users {
		if _, ok := accountDeletionTombstoneFromUser(user); !ok {
			continue
		}
		current, ok := out[user.ID]
		if !ok || user.UpdatedAt >= current.UpdatedAt {
			out[user.ID] = user
		}
	}
	return out
}

func mergeAccountDeletionTombstones(restored, current IdentityPersistentState) IdentityPersistentState {
	tombstones := accountDeletionTombstones(restored)
	for userID, tombstone := range accountDeletionTombstones(current) {
		prior, ok := tombstones[userID]
		if !ok || tombstone.UpdatedAt >= prior.UpdatedAt {
			tombstones[userID] = tombstone
		}
	}
	if len(tombstones) == 0 {
		return restored
	}
	users := make([]UserRecord, 0, len(restored.Users)+len(tombstones))
	for _, user := range restored.Users {
		if _, deleted := tombstones[user.ID]; deleted {
			continue
		}
		users = append(users, user)
	}
	ids := make([]string, 0, len(tombstones))
	for userID := range tombstones {
		ids = append(ids, userID)
	}
	sort.Strings(ids)
	for _, userID := range ids {
		users = append(users, tombstones[userID])
	}
	devices := make([]DeviceRecord, 0, len(restored.Devices))
	for _, device := range restored.Devices {
		if _, deleted := tombstones[device.UserID]; !deleted {
			devices = append(devices, device)
		}
	}
	sessions := make([]SessionRecord, 0, len(restored.Sessions))
	for _, session := range restored.Sessions {
		if _, deleted := tombstones[session.UserID]; !deleted {
			sessions = append(sessions, session)
		}
	}
	restored.Users = users
	restored.Devices = devices
	restored.Sessions = sessions
	return restored
}

func enforceArchiveAccountDeletionPrivacy(archive PersistenceArchive, current IdentityPersistentState) PersistenceArchive {
	archive.Identity = mergeAccountDeletionTombstones(archive.Identity, current)
	tombstones := accountDeletionTombstones(archive.Identity)
	if len(tombstones) == 0 {
		return archive
	}
	archive.HasIdentity = true
	for i := range archive.UserWorkspaces {
		tombstone, deleted := tombstones[archive.UserWorkspaces[i].UserID]
		if !deleted {
			continue
		}
		archive.UserWorkspaces[i] = privacyBlankWorkspace(tombstone.ID, tombstone.UpdatedAt)
	}
	return archive
}

func (a *Application) handleAccountDataExport(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}
	principal, ok := principalFromContext(r.Context())
	if !ok {
		writeError(w, http.StatusUnauthorized, "Authentication required.")
		return
	}
	export, err := a.accountDataExport(principal.UserID)
	if err != nil {
		writeError(w, http.StatusServiceUnavailable, "Account export unavailable.")
		return
	}
	w.Header().Set("Cache-Control", "no-store")
	w.Header().Set("Content-Disposition", `attachment; filename="depulse-account-export.json"`)
	writeJSON(w, http.StatusOK, export)
}

func (a *Application) handleDeleteAccount(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		writeError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}
	principal, ok := principalFromContext(r.Context())
	if !ok {
		writeError(w, http.StatusUnauthorized, "Authentication required.")
		return
	}
	var req struct {
		Reason string `json:"reason"`
	}
	if err := decodeJSONBody(w, r, &req); err != nil {
		writeError(w, http.StatusBadRequest, "Invalid account deletion request.")
		return
	}
	tombstone, err := a.deleteAccountData(principal.UserID, req.Reason)
	if err != nil {
		switch {
		case strings.Contains(err.Error(), "unsupported account deletion reason"):
			writeError(w, http.StatusBadRequest, "Unsupported account deletion reason.")
		case strings.Contains(err.Error(), "last active owner"):
			writeError(w, http.StatusConflict, "Another active owner is required before deleting this account.")
		default:
			writeError(w, http.StatusServiceUnavailable, "Account deletion could not be completed.")
		}
		return
	}
	clearSessionCookie(w, r)
	w.Header().Set("Cache-Control", "no-store")
	writeJSON(w, http.StatusOK, map[string]any{"ok": true, "tombstone": tombstone})
}
