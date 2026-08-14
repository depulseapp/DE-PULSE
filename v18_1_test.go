package main

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
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
