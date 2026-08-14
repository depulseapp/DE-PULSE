package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

// v16.6 defect closure: an explicitly empty desk is valid state. The final
// removal must not be serialized as nil because nil is reserved for legacy
// "missing watchlist" migration and would restore defaults on normalization.
func TestV166LastMasterSymbolRemovalPersistsExplicitEmpty(t *testing.T) {
	app := newTestApplication(t)
	const sym = "AAPL"
	app.mu.Lock()
	for _, id := range deskIDs() {
		wl, _ := watchlistValueByID(app.state.Watchlists, id)
		wl.Symbols = []string{sym}
		for i := range app.state.Watchlists {
			if app.state.Watchlists[i].ID == id {
				app.state.Watchlists[i] = wl
			}
		}
	}
	app.state.UI.SelectedTicker = sym
	app.mu.Unlock()

	rr := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/api/master-symbol/remove", strings.NewReader(`{"symbol":"AAPL"}`))
	req.Header.Set("Content-Type", "application/json")
	app.handleMasterSymbolRemove(rr, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("remove status=%d body=%s", rr.Code, rr.Body.String())
	}

	app.mu.RLock()
	st := clone(app.state)
	app.mu.RUnlock()
	for _, id := range deskIDs() {
		wl, ok := watchlistValueByID(st.Watchlists, id)
		if !ok || wl.Symbols == nil || len(wl.Symbols) != 0 {
			t.Fatalf("%s must persist as explicit empty slice, got %#v", id, wl.Symbols)
		}
	}

	// Re-normalization must not restore default symbols.
	ensureDedicatedDeskWatchlists(&st, defaultState())
	for _, id := range deskIDs() {
		wl, _ := watchlistValueByID(st.Watchlists, id)
		if wl.Symbols == nil || len(wl.Symbols) != 0 {
			t.Fatalf("%s repopulated after normalization: %#v", id, wl.Symbols)
		}
	}

	raw, err := os.ReadFile(app.statePath())
	if err != nil {
		t.Fatal(err)
	}
	var persisted AppState
	if err := json.Unmarshal(raw, &persisted); err != nil {
		t.Fatal(err)
	}
	for _, id := range deskIDs() {
		wl, _ := watchlistValueByID(persisted.Watchlists, id)
		if wl.Symbols == nil || len(wl.Symbols) != 0 {
			t.Fatalf("%s persisted as nil/repopulated: %#v", id, wl.Symbols)
		}
	}

	// Simulate application restart using the same config directory.
	app.load()
	for _, id := range deskIDs() {
		wl, _ := watchlistValueByID(app.state.Watchlists, id)
		if wl.Symbols == nil || len(wl.Symbols) != 0 {
			t.Fatalf("%s repopulated after reload: %#v", id, wl.Symbols)
		}
	}
}

func TestV166MasterRemoveAllClearsEveryUserDeskSymbolAndSurvivesReload(t *testing.T) {
	app := newTestApplication(t)
	app.mu.Lock()
	for _, id := range deskIDs() {
		wl, _ := watchlistValueByID(app.state.Watchlists, id)
		wl.Symbols = []string{"AAPL", "MSFT", "NVDA"}
		for i := range app.state.Watchlists {
			if app.state.Watchlists[i].ID == id {
				app.state.Watchlists[i] = wl
			}
		}
	}
	app.state.UI.SelectedTicker = "NVDA"
	app.mu.Unlock()

	rr := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/api/master-symbol/remove-all", strings.NewReader(`{}`))
	req.Header.Set("Content-Type", "application/json")
	app.handleMasterSymbolRemoveAll(rr, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("remove-all status=%d body=%s", rr.Code, rr.Body.String())
	}
	var out struct {
		OK           bool `json:"ok"`
		RemovedCount int  `json:"removedCount"`
	}
	if err := json.Unmarshal(rr.Body.Bytes(), &out); err != nil {
		t.Fatal(err)
	}
	if !out.OK || out.RemovedCount != 3 {
		t.Fatalf("unexpected remove-all response: %+v body=%s", out, rr.Body.String())
	}

	app.load()
	if got := activeDeskSymbolsFromState(app.state); len(got) != 0 {
		t.Fatalf("user Master Symbol Store repopulated after reload: %v", got)
	}
	for _, id := range deskIDs() {
		wl, _ := watchlistValueByID(app.state.Watchlists, id)
		if wl.Symbols == nil || len(wl.Symbols) != 0 {
			t.Fatalf("%s not explicitly empty after reload: %#v", id, wl.Symbols)
		}
	}
}

func TestV166ClearingUserMasterStoreKeepsProtectedMarketContext(t *testing.T) {
	app := newTestApplication(t)
	app.mu.Lock()
	for _, id := range deskIDs() {
		wl, _ := watchlistValueByID(app.state.Watchlists, id)
		wl.Symbols = []string{}
		for i := range app.state.Watchlists {
			if app.state.Watchlists[i].ID == id {
				app.state.Watchlists[i] = wl
			}
		}
	}
	app.mu.Unlock()

	tracked := app.engine.trackedSymbols()
	for _, required := range []string{"SPY", "QQQ", "GLD", "SLV", "USO", "VIX"} {
		if !contains(tracked, required) {
			t.Fatalf("protected market context %s disappeared after user store clear: %v", required, tracked)
		}
	}
	if got := activeDeskSymbolsFromState(app.state); len(got) != 0 {
		t.Fatalf("protected context leaked back into user-managed desks: %v", got)
	}
}
