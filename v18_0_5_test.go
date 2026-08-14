package main

import (
	"encoding/json"
	"reflect"
	"testing"
)

func TestV1805TrackedSymbolCanonicalOwnerIncludesDisabledDesks(t *testing.T) {
	st := defaultState()
	st.Settings.DayEnabled = false
	st.Settings.SwingEnabled = true
	st.Settings.LongEnabled = false
	st.Watchlists = []Watchlist{
		{ID: "day", Name: "Day Trade Watchlist", Symbols: []string{"AAA", "DAYONLY"}},
		{ID: "swing", Name: "Swing Watchlist", Symbols: []string{"AAA", "SWONLY"}},
		{ID: "long", Name: "Long-Term Watchlist", Symbols: []string{"AAA", "LNGONLY"}},
		{ID: "discovery", Name: "Discovery Watchlist", Symbols: []string{"DISC"}},
	}

	got := trackedSymbolsLocked(&st)
	want := []string{"AAA", "DAYONLY", "SWONLY", "LNGONLY"}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("tracked symbols = %#v, want %#v", got, want)
	}

	removed := clearTrackedSymbolsLocked(&st)
	if len(removed) != 4 {
		t.Fatalf("removed membership entries = %d, want 4", len(removed))
	}
	for _, id := range deskIDs() {
		wl, ok := watchlistValueByID(st.Watchlists, id)
		if !ok {
			t.Fatalf("missing %s desk", id)
		}
		if wl.Symbols == nil || len(wl.Symbols) != 0 {
			t.Fatalf("%s desk must preserve explicit empty truth, got %#v", id, wl.Symbols)
		}
	}
	if got, ok := watchlistValueByID(st.Watchlists, "discovery"); !ok || !reflect.DeepEqual(got.Symbols, []string{"DISC"}) {
		t.Fatalf("discovery staging must be untouched, got %#v", got)
	}
}

func TestV1805TrackedSymbolsEmptySurvivesReloadMerge(t *testing.T) {
	st := defaultState()
	clearTrackedSymbolsLocked(&st)

	data, err := json.Marshal(st)
	if err != nil {
		t.Fatal(err)
	}
	var reloaded AppState
	if err := json.Unmarshal(data, &reloaded); err != nil {
		t.Fatal(err)
	}
	reloaded = mergeState(reloaded)

	if got := trackedSymbolsLocked(&reloaded); len(got) != 0 {
		t.Fatalf("empty tracked state was rehydrated unexpectedly: %#v", got)
	}
	for _, id := range deskIDs() {
		wl, ok := watchlistValueByID(reloaded.Watchlists, id)
		if !ok || wl.Symbols == nil || len(wl.Symbols) != 0 {
			t.Fatalf("%s explicit empty state did not survive merge: %#v", id, wl)
		}
	}
}

func TestV1805TrackedSymbolMutationIsIdempotent(t *testing.T) {
	st := defaultState()
	clearTrackedSymbolsLocked(&st)

	first := setTrackedSymbolLocked(&st, "nvda", true)
	setTrackedSymbolLocked(&st, "NVDA", true)
	if activeDeskCount(first) != len(deskIDs()) {
		t.Fatalf("first add membership = %#v", first)
	}
	if got := trackedSymbolsLocked(&st); !reflect.DeepEqual(got, []string{"NVDA"}) {
		t.Fatalf("tracked symbols after duplicate add = %#v", got)
	}
	for _, id := range deskIDs() {
		wl, _ := watchlistValueByID(st.Watchlists, id)
		if !reflect.DeepEqual(wl.Symbols, []string{"NVDA"}) {
			t.Fatalf("duplicate add duplicated %s membership: %#v", id, wl.Symbols)
		}
	}

	removed := setTrackedSymbolLocked(&st, "NVDA", false)
	again := setTrackedSymbolLocked(&st, "NVDA", false)
	if activeDeskCount(removed) != 0 || activeDeskCount(again) != 0 {
		t.Fatalf("remove must converge to empty membership: removed=%#v again=%#v", removed, again)
	}
	if got := trackedSymbolsLocked(&st); len(got) != 0 {
		t.Fatalf("tracked symbols after duplicate remove = %#v", got)
	}
}
