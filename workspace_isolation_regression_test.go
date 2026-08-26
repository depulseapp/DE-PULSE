package main

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"
)

func v181TestApp(t *testing.T) *Application {
	t.Helper()
	app := newTestApplication(t)
	neutral := defaultEmptyWorkspaceState()
	app.state.Watchlists = neutral.Watchlists
	app.state.UI = neutral.UI
	app.workspaces = map[string]UserWorkspace{}
	return app
}

func v181WorkspaceWithSymbols(userID string, symbols ...string) UserWorkspace {
	st := defaultEmptyWorkspaceState()
	for _, sym := range symbols {
		setTrackedSymbolLocked(&st, sym, true)
	}
	if len(symbols) > 0 {
		st.UI.SelectedTicker = symbols[0]
	}
	return workspaceFromState(userID, st, time.UnixMilli(1_700_000_000_000))
}

func v181Request(method, target, body string, principal Principal) *http.Request {
	req := httptest.NewRequest(method, target, strings.NewReader(body))
	return req.WithContext(context.WithValue(req.Context(), identityContextKey{}, principal))
}

func TestV181WorkspaceIsolationAndEmptyNewUser(t *testing.T) {
	app := v181TestApp(t)
	app.workspaces["user-a"] = v181WorkspaceWithSymbols("user-a", "NVDA")
	app.workspaces["user-b"] = v181WorkspaceWithSymbols("user-b", "MSFT")

	a := app.publicStateForUser("user-a")
	b := app.publicStateForUser("user-b")
	c := app.publicStateForUser("user-c")
	if !contains(trackedSymbolsFromWatchlists(a.Watchlists), "NVDA") || contains(trackedSymbolsFromWatchlists(a.Watchlists), "MSFT") {
		t.Fatalf("user A workspace leaked: %+v", a.Watchlists)
	}
	if !contains(trackedSymbolsFromWatchlists(b.Watchlists), "MSFT") || contains(trackedSymbolsFromWatchlists(b.Watchlists), "NVDA") {
		t.Fatalf("user B workspace leaked: %+v", b.Watchlists)
	}
	if got := trackedSymbolsFromWatchlists(c.Watchlists); len(got) != 0 {
		t.Fatalf("new user inherited another user's symbols: %v", got)
	}
}

func trackedSymbolsFromWatchlists(watchlists []Watchlist) []string {
	st := defaultEmptyWorkspaceState()
	st.Watchlists = clone(watchlists)
	return trackedSymbolsLocked(&st)
}

func TestV181ProcessingUnionDeduplicatesAndRetainsOtherOwner(t *testing.T) {
	app := v181TestApp(t)
	app.workspaces["user-a"] = v181WorkspaceWithSymbols("user-a", "NVDA", "AMD")
	app.workspaces["user-b"] = v181WorkspaceWithSymbols("user-b", "NVDA", "MSFT")

	app.mu.RLock()
	processing := app.processingStateLocked()
	app.mu.RUnlock()
	active := activeDeskSymbolsFromState(processing)
	if count := countSymbol(active, "NVDA"); count != 1 {
		t.Fatalf("shared processing universe must deduplicate NVDA, count=%d active=%v", count, active)
	}
	for _, sym := range []string{"AMD", "MSFT"} {
		if !contains(active, sym) {
			t.Fatalf("shared processing universe missing %s: %v", sym, active)
		}
	}

	principalA := Principal{UserID: "user-a", Username: "a", Role: RoleUser}
	rr := httptest.NewRecorder()
	app.handleMasterSymbolRemove(rr, v181Request(http.MethodPost, "/api/master-symbol/remove", `{"symbol":"NVDA"}`, principalA))
	if rr.Code != http.StatusOK {
		t.Fatalf("remove failed: %d %s", rr.Code, rr.Body.String())
	}
	if contains(app.workspaceSymbols("user-a"), "NVDA") {
		t.Fatal("user A still owns NVDA after removal")
	}
	if !contains(app.workspaceSymbols("user-b"), "NVDA") {
		t.Fatal("user A removal mutated user B workspace")
	}
	app.mu.RLock()
	processing = app.processingStateLocked()
	app.mu.RUnlock()
	if !contains(activeDeskSymbolsFromState(processing), "NVDA") {
		t.Fatal("NVDA dropped from shared processing universe while user B still owns it")
	}
}

func (a *Application) workspaceSymbols(userID string) []string {
	a.mu.RLock()
	defer a.mu.RUnlock()
	return trackedSymbolsLocked(ptrAppState(a.workspaceStateLocked(userID)))
}

func ptrAppState(st AppState) *AppState { return &st }

func countSymbol(items []string, symbol string) int {
	n := 0
	for _, item := range items {
		if normalizeSymbol(item) == normalizeSymbol(symbol) {
			n++
		}
	}
	return n
}

func TestV181MasterSymbolMutationIsUserScoped(t *testing.T) {
	app := v181TestApp(t)
	app.workspaces["user-a"] = defaultUserWorkspace("user-a")
	app.workspaces["user-b"] = defaultUserWorkspace("user-b")

	principalA := Principal{UserID: "user-a", Username: "a", Role: RoleUser}
	rr := httptest.NewRecorder()
	app.handleMasterSymbolAdd(rr, v181Request(http.MethodPost, "/api/master-symbol/add", `{"symbol":"NVDA"}`, principalA))
	if rr.Code != http.StatusOK {
		t.Fatalf("add failed: %d %s", rr.Code, rr.Body.String())
	}
	if !contains(app.workspaceSymbols("user-a"), "NVDA") {
		t.Fatal("requesting user did not receive NVDA")
	}
	if contains(app.workspaceSymbols("user-b"), "NVDA") {
		t.Fatal("master-symbol add leaked to another user")
	}
	if got := trackedSymbolsLocked(&app.state); len(got) != 0 {
		t.Fatalf("shared state retained personal tracked symbols: %v", got)
	}
}

func TestV181HubTargetsPersonalStateEvents(t *testing.T) {
	h := NewHub()
	a := h.SubscribeForUser("user-a")
	b := h.SubscribeForUser("user-b")
	defer h.Unsubscribe(a)
	defer h.Unsubscribe(b)
	h.BroadcastUser("user-a", map[string]any{"type": "state", "owner": "user-a"})
	select {
	case payload := <-a:
		var got map[string]any
		if err := json.Unmarshal(payload, &got); err != nil || got["owner"] != "user-a" {
			t.Fatalf("unexpected targeted payload: %s", payload)
		}
	case <-time.After(500 * time.Millisecond):
		t.Fatal("user A did not receive its state event")
	}
	select {
	case payload := <-b:
		t.Fatalf("user B received user A state event: %s", payload)
	case <-time.After(50 * time.Millisecond):
	}
}

func TestV181LegacyOwnerMigrationAndWorkspacePersistence(t *testing.T) {
	dir := t.TempDir()
	p := NewPersistenceManager(dir)
	defer p.Close()
	app := &Application{configDir: dir, hub: NewHub(), state: defaultState(), persistence: p}
	legacyTicker := app.state.UI.SelectedTicker
	legacySymbols := trackedSymbolsLocked(&app.state)
	if len(legacySymbols) == 0 {
		t.Fatal("fixture must contain legacy Stable symbols")
	}
	if err := app.initializeUserWorkspaces(); err != nil {
		t.Fatal(err)
	}
	owner := app.publicStateForUser(bootstrapOwnerID)
	if owner.UI.SelectedTicker != legacyTicker || len(trackedSymbolsFromWatchlists(owner.Watchlists)) == 0 {
		t.Fatalf("legacy owner workspace not preserved: %+v", owner)
	}
	app.mu.RLock()
	if got := trackedSymbolsLocked(&app.state); len(got) != 0 {
		app.mu.RUnlock()
		t.Fatalf("shared state still owns legacy personal symbols: %v", got)
	}
	app.mu.RUnlock()
	loaded, err := p.LoadUserWorkspaces(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	found := false
	for _, ws := range loaded {
		if ws.UserID == bootstrapOwnerID {
			found = true
			if len(trackedSymbolsFromWatchlists(ws.Watchlists)) == 0 {
				t.Fatal("persisted owner workspace lost legacy symbols")
			}
		}
	}
	if !found {
		t.Fatal("bootstrap owner workspace was not persisted")
	}
}

func TestV181RuntimeSnapshotFiltersOtherWorkspaceSymbols(t *testing.T) {
	app := v181TestApp(t)
	app.workspaces["user-a"] = v181WorkspaceWithSymbols("user-a", "NVDA")
	app.workspaces["user-b"] = v181WorkspaceWithSymbols("user-b", "MSFT")
	snap := RuntimeSnapshot{
		Quotes: map[string]Quote{
			"SPY":  {Symbol: "SPY", Price: 600},
			"NVDA": {Symbol: "NVDA", Price: 180},
			"MSFT": {Symbol: "MSFT", Price: 500},
		},
		History: map[string][]HistoryPoint{"NVDA": {{T: 1, P: 180}}, "MSFT": {{T: 1, P: 500}}},
		Feed: FeedDiagnostics{
			SubscribedSymbols:       []string{"SPY", "NVDA", "MSFT"},
			AlpacaSubscribedSymbols: []string{"NVDA", "MSFT"},
			SnapshotSymbols:         []string{"NVDA", "MSFT"},
			LastTradeSymbol:         "MSFT",
		},
		News: []NewsItem{
			{Headline: "NVDA item", Symbols: []string{"NVDA"}},
			{Headline: "MSFT item", Symbols: []string{"MSFT"}},
			{Headline: "Market item"},
		},
		Earnings: []EarningsItem{{Symbol: "NVDA"}, {Symbol: "MSFT"}},
		SECIntelligence: map[string]SECIntelligenceSummary{
			"NVDA": {Symbol: "NVDA"}, "MSFT": {Symbol: "MSFT"},
		},
		AdaptiveDataPolicy: AdaptiveDataPolicyState{HotSymbols: []string{"NVDA", "MSFT"}},
	}
	got := app.runtimeSnapshotForUserFrom("user-a", snap)
	if got.Quotes["NVDA"].Symbol != "NVDA" || got.Quotes["SPY"].Symbol != "SPY" {
		t.Fatalf("allowed runtime symbols missing: %+v", got.Quotes)
	}
	if _, ok := got.Quotes["MSFT"]; ok {
		t.Fatal("user A runtime leaked user B quote")
	}
	if _, ok := got.History["MSFT"]; ok {
		t.Fatal("user A runtime leaked user B history")
	}
	if _, ok := got.SECIntelligence["MSFT"]; ok {
		t.Fatal("user A runtime leaked user B SEC intelligence")
	}
	if contains(got.Feed.SubscribedSymbols, "MSFT") || contains(got.Feed.SnapshotSymbols, "MSFT") || got.Feed.LastTradeSymbol == "MSFT" {
		t.Fatalf("feed diagnostics leaked user B subscription: %+v", got.Feed)
	}
	if contains(got.AdaptiveDataPolicy.HotSymbols, "MSFT") {
		t.Fatalf("adaptive hot-symbol diagnostics leaked user B: %+v", got.AdaptiveDataPolicy.HotSymbols)
	}
	if len(got.News) != 2 || got.News[0].Headline != "NVDA item" || got.News[1].Headline != "Market item" {
		t.Fatalf("unexpected filtered news: %+v", got.News)
	}
	if len(got.Earnings) != 1 || got.Earnings[0].Symbol != "NVDA" {
		t.Fatalf("unexpected filtered earnings: %+v", got.Earnings)
	}
}

func TestV181SymbolEventFanoutTargetsWorkspaceOwners(t *testing.T) {
	app := v181TestApp(t)
	app.workspaces["user-a"] = v181WorkspaceWithSymbols("user-a", "NVDA")
	app.workspaces["user-b"] = v181WorkspaceWithSymbols("user-b", "MSFT")
	aCh := app.hub.SubscribeForUser("user-a")
	bCh := app.hub.SubscribeForUser("user-b")
	defer app.hub.Unsubscribe(aCh)
	defer app.hub.Unsubscribe(bCh)

	app.broadcastSymbolEvent("NVDA", map[string]any{"type": "quote", "symbol": "NVDA"})
	select {
	case <-aCh:
	case <-time.After(500 * time.Millisecond):
		t.Fatal("NVDA owner did not receive quote event")
	}
	select {
	case payload := <-bCh:
		t.Fatalf("non-owner received private-symbol quote event: %s", payload)
	case <-time.After(50 * time.Millisecond):
	}

	app.broadcastSymbolEvent("SPY", map[string]any{"type": "quote", "symbol": "SPY"})
	for name, ch := range map[string]chan []byte{"a": aCh, "b": bCh} {
		select {
		case <-ch:
		case <-time.After(500 * time.Millisecond):
			t.Fatalf("%s did not receive shared market-context quote", name)
		}
	}
}

func TestV181EnsureUserWorkspaceCreatesDurableEmptyWorkspace(t *testing.T) {
	dir := t.TempDir()
	p := NewPersistenceManager(dir)
	defer p.Close()
	app := &Application{configDir: dir, hub: NewHub(), state: defaultEmptyWorkspaceState(), persistence: p, identity: &IdentityService{}, workspaces: map[string]UserWorkspace{}}
	if err := app.ensureUserWorkspace("user-new"); err != nil {
		t.Fatal(err)
	}
	state := app.publicStateForUser("user-new")
	if got := trackedSymbolsFromWatchlists(state.Watchlists); len(got) != 0 {
		t.Fatalf("new durable workspace must be empty, got %v", got)
	}
	loaded, err := p.LoadUserWorkspaces(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	found := false
	for _, ws := range loaded {
		if ws.UserID == "user-new" {
			found = true
			if len(trackedSymbolsFromWatchlists(ws.Watchlists)) != 0 {
				t.Fatalf("persisted new workspace inherited symbols: %+v", ws.Watchlists)
			}
		}
	}
	if !found {
		t.Fatal("new authenticated workspace was not persisted")
	}
}

func hostPrivacyUserFixture(t *testing.T) (*PersistenceManager, *IdentityService, *Application, string, Principal, string) {
	t.Helper()
	p, identity := newIdentityTestService(t)
	base := time.Unix(2_430_000_000, 0)
	_, owner, _ := v184CredentialedOwner(t, identity, base)
	created, err := identity.adminCreateUser(owner, "privacy-user", "Privacy User", RoleUser, "temporary privacy password")
	if err != nil {
		t.Fatal(err)
	}
	token, principal, err := identity.setPassword(created.ID, "durable privacy password")
	if err != nil {
		t.Fatal(err)
	}
	app := &Application{configDir: t.TempDir(), hub: NewHub(), state: defaultEmptyWorkspaceState(), persistence: p, identity: identity, workspaces: map[string]UserWorkspace{}}
	workspace := v181WorkspaceWithSymbols(created.ID, "NVDA", "AMD")
	app.workspaces[created.ID] = workspace
	if err := p.SaveUserWorkspace(context.Background(), workspace); err != nil {
		t.Fatal(err)
	}
	return p, identity, app, token, principal, "durable privacy password"
}

func hostPrivacyGET(path, token string) *http.Request {
	req := httptest.NewRequest(http.MethodGet, path, nil)
	req.AddCookie(&http.Cookie{Name: sessionCookieName, Value: token})
	return req
}

func hostPrivacyDELETE(path, token, csrf, body string) *http.Request {
	req := httptest.NewRequest(http.MethodDelete, path, strings.NewReader(body))
	req.AddCookie(&http.Cookie{Name: sessionCookieName, Value: token})
	req.AddCookie(&http.Cookie{Name: csrfCookieName, Value: csrf})
	req.Header.Set("X-DE-PULSE-CSRF", csrf)
	req.Header.Set("Content-Type", "application/json")
	return req
}

func TestHOST010AccountExportIncludesPortableCategoriesAndExcludesSecrets(t *testing.T) {
	_, identity, app, token, principal, _ := hostPrivacyUserFixture(t)
	fingerprint := "sha256:privacy-device-fingerprint-secret"
	device, err := identity.registerHostedDevice(principal, "privacy browser", fingerprint)
	if err != nil {
		t.Fatal(err)
	}
	if err := identity.bindHostedDeviceToSession(principal, device.ID); err != nil {
		t.Fatal(err)
	}

	identity.mu.Lock()
	passwordHash := ""
	tokenHash := ""
	fingerprintHash := ""
	for _, user := range identity.state.Users {
		if user.ID == principal.UserID {
			passwordHash = user.PasswordHash
		}
	}
	for _, session := range identity.state.Sessions {
		if session.UserID == principal.UserID {
			tokenHash = session.TokenHash
		}
	}
	for _, storedDevice := range identity.state.Devices {
		if storedDevice.UserID == principal.UserID {
			fingerprintHash = storedDevice.FingerprintHash
		}
	}
	identity.mu.Unlock()
	if passwordHash == "" || tokenHash == "" || fingerprintHash == "" {
		t.Fatal("fixture did not contain credential material to prove export exclusion")
	}

	rr := httptest.NewRecorder()
	hostedIdentityHTTPMux(app).ServeHTTP(rr, hostPrivacyGET("/api/auth/account/export", token))
	if rr.Code != http.StatusOK {
		t.Fatalf("account export failed: code=%d body=%s", rr.Code, rr.Body.String())
	}
	if rr.Header().Get("Cache-Control") != "no-store" {
		t.Fatalf("account export is cacheable: headers=%v", rr.Header())
	}
	body := rr.Body.String()
	for label, secret := range map[string]string{"password hash": passwordHash, "session token hash": tokenHash, "device fingerprint": fingerprintHash, "raw token": token} {
		if secret != "" && strings.Contains(body, secret) {
			t.Fatalf("%s leaked in account export: %s", label, body)
		}
	}
	var exported AccountDataExport
	if err := json.Unmarshal(rr.Body.Bytes(), &exported); err != nil {
		t.Fatal(err)
	}
	if exported.SchemaVersion != 1 || exported.Profile.ID != principal.UserID || exported.Profile.Username != principal.Username {
		t.Fatalf("identity/profile missing from export: %+v", exported.Profile)
	}
	if !contains(trackedSymbolsFromWatchlists(exported.Watchlists), "NVDA") || exported.Preferences.SelectedTicker != "NVDA" {
		t.Fatalf("workspace preferences/watchlists missing from export: %+v", exported)
	}
	if exported.ResearchMetadata == nil || exported.Feedback == nil || exported.CategoryStatus["researchMetadata"] != "not-persisted-per-user" || exported.CategoryStatus["feedback"] != "not-persisted-per-user" {
		t.Fatalf("non-persisted categories were not represented truthfully: %+v", exported)
	}
	if len(exported.HostedSyncMetadata.Devices) != 1 || len(exported.HostedSyncMetadata.Sessions) == 0 {
		t.Fatalf("safe hosted sync metadata missing: %+v", exported.HostedSyncMetadata)
	}
}

func TestHOST011AccountDeletionRevokesIdentityAndScrubsWorkspace(t *testing.T) {
	p, identity, app, token, principal, _ := hostPrivacyUserFixture(t)
	device, err := identity.registerHostedDevice(principal, "delete browser", "sha256:delete-device-secret")
	if err != nil {
		t.Fatal(err)
	}
	if err := identity.bindHostedDeviceToSession(principal, device.ID); err != nil {
		t.Fatal(err)
	}
	csrf := "privacy-delete-csrf"
	rr := httptest.NewRecorder()
	hostedIdentityHTTPMux(app).ServeHTTP(rr, hostPrivacyDELETE("/api/auth/account", token, csrf, `{"reason":"USER_REQUESTED"}`))
	if rr.Code != http.StatusOK {
		t.Fatalf("account deletion failed: code=%d body=%s", rr.Code, rr.Body.String())
	}
	if _, err := identity.resolve(token, false); err == nil {
		t.Fatal("deleted account session remained usable")
	}

	identity.mu.Lock()
	var tombstone UserRecord
	found := false
	for _, user := range identity.state.Users {
		if user.ID == principal.UserID {
			tombstone = user
			found = true
		}
	}
	for _, session := range identity.state.Sessions {
		if session.UserID == principal.UserID {
			identity.mu.Unlock()
			t.Fatalf("deleted account session retained: %+v", session)
		}
	}
	for _, storedDevice := range identity.state.Devices {
		if storedDevice.UserID == principal.UserID {
			identity.mu.Unlock()
			t.Fatalf("deleted account device retained: %+v", storedDevice)
		}
	}
	identity.mu.Unlock()
	if !found || !isAccountDeletionTombstone(tombstone) {
		t.Fatalf("minimal deletion tombstone missing: %+v", tombstone)
	}
	if tombstone.TenantID == "" || tombstone.DisplayName != accountDeletionReasonUserRequest || tombstone.PasswordHash != "" || tombstone.Role != "" || tombstone.CreatedAt != 0 || tombstone.LastLoginAt != 0 {
		t.Fatalf("tombstone retained more than deletion proof: %+v", tombstone)
	}

	workspaces, err := p.LoadUserWorkspaces(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	workspaceFound := false
	for _, workspace := range workspaces {
		if workspace.UserID != principal.UserID {
			continue
		}
		workspaceFound = true
		if got := trackedSymbolsFromWatchlists(workspace.Watchlists); len(got) != 0 {
			t.Fatalf("deleted account workspace retained personal symbols: %v", got)
		}
	}
	if !workspaceFound {
		t.Fatal("deleted account did not retain the expected minimal empty workspace row")
	}
}

func TestHOST011LastCriticalOwnerCannotSelfDelete(t *testing.T) {
	_, identity := newIdentityTestService(t)
	base := time.Unix(2_440_000_000, 0)
	token, owner, _ := v184CredentialedOwner(t, identity, base)
	if _, err := identity.deleteAccount(owner.UserID, accountDeletionReasonUserRequest); err == nil || !strings.Contains(err.Error(), "last active owner") {
		t.Fatalf("sole critical owner deletion was not rejected: %v", err)
	}
	if _, err := identity.resolve(token, false); err != nil {
		t.Fatalf("rejected owner deletion damaged current identity: %v", err)
	}
}

func TestHOST011ReplaceRestoreCarriesForwardDeletionTombstone(t *testing.T) {
	deletedAt := int64(2_450_000_000_000)
	currentTombstone := UserRecord{ID: "privacy-deleted", TenantID: localTenantID, Username: accountDeletionTombstoneUsername, DisplayName: accountDeletionReasonUserRequest, Status: UserDisabled, UpdatedAt: deletedAt}
	archive := PersistenceArchive{
		HasIdentity: true,
		Identity: IdentityPersistentState{
			Version: 1,
			Users: []UserRecord{{ID: "privacy-deleted", TenantID: localTenantID, Username: "old-profile", DisplayName: "Old Profile", Role: RoleUser, Status: UserActive, PasswordHash: "old-password-hash"}},
			Devices: []DeviceRecord{{ID: "old-device", TenantID: localTenantID, UserID: "privacy-deleted", FingerprintHash: "old-fingerprint", Status: DeviceActive}},
			Sessions: []SessionRecord{{ID: "old-session", UserID: "privacy-deleted", TokenHash: "old-token-hash"}},
		},
		UserWorkspaces: []UserWorkspace{v181WorkspaceWithSymbols("privacy-deleted", "NVDA", "TSLA")},
	}
	got := enforceArchiveAccountDeletionPrivacy(archive, IdentityPersistentState{Users: []UserRecord{currentTombstone}})
	if len(got.Identity.Users) != 1 || !isAccountDeletionTombstone(got.Identity.Users[0]) || got.Identity.Users[0].UpdatedAt != deletedAt {
		t.Fatalf("replace restore resurrected deleted identity: %+v", got.Identity.Users)
	}
	if len(got.Identity.Devices) != 0 || len(got.Identity.Sessions) != 0 {
		t.Fatalf("replace restore resurrected deleted security records: devices=%+v sessions=%+v", got.Identity.Devices, got.Identity.Sessions)
	}
	if len(got.UserWorkspaces) != 1 || len(trackedSymbolsFromWatchlists(got.UserWorkspaces[0].Watchlists)) != 0 {
		t.Fatalf("replace restore resurrected deleted workspace: %+v", got.UserWorkspaces)
	}
}

func TestHOST012PrivacyDataHandlingInventoryCoversPersistentStoresAndREST(t *testing.T) {
	raw, err := os.ReadFile("governance/programs/ADAPT-HOSTED-SYNC-001/privacy-data-handling.json")
	if err != nil {
		t.Fatal(err)
	}
	var inventory struct {
		Requirements []string `json:"requirements"`
		Stores       []struct {
			ID        string `json:"id"`
			Deletion  string `json:"deletion"`
			Retention string `json:"retention"`
			Export    string `json:"export"`
		} `json:"stores"`
		REST struct {
			Export struct {
				Method string `json:"method"`
				Path   string `json:"path"`
			} `json:"export"`
			Delete struct {
				Method string `json:"method"`
				Path   string `json:"path"`
			} `json:"delete"`
		} `json:"rest"`
		OperatorAccess map[string]string `json:"operatorAccess"`
	}
	if err := json.Unmarshal(raw, &inventory); err != nil {
		t.Fatal(err)
	}
	for _, requirement := range []string{"HOST-010", "HOST-011", "HOST-012"} {
		if !contains(inventory.Requirements, requirement) {
			t.Fatalf("privacy inventory lost requirement %s", requirement)
		}
	}
	stores := map[string]bool{}
	for _, store := range inventory.Stores {
		stores[store.ID] = store.Deletion != "" && store.Retention != "" && store.Export != ""
	}
	for _, required := range []string{"sqlite", "postgres", "cache", "pitr-backup", "audit"} {
		if !stores[required] {
			t.Fatalf("privacy inventory missing complete %s behavior: %+v", required, stores)
		}
	}
	if inventory.REST.Export.Method != http.MethodGet || inventory.REST.Export.Path != "/api/auth/account/export" || inventory.REST.Delete.Method != http.MethodDelete || inventory.REST.Delete.Path != "/api/auth/account" {
		t.Fatalf("privacy REST contract drifted: %+v", inventory.REST)
	}
	if inventory.OperatorAccess["restore"] == "" || inventory.OperatorAccess["deletedAccount"] == "" {
		t.Fatalf("operator privacy behavior missing: %+v", inventory.OperatorAccess)
	}
}
