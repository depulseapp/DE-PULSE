package main

import "testing"

// extreme30TestApp is the capability-owned in-memory Application fixture for
// the permanent Extreme-30 regression matrix. It preserves only the generic
// setup formerly supplied by the retired pre-v17 test stack.
func extreme30TestApp(t *testing.T) *Application {
	t.Helper()
	st := defaultState()
	ensureDedicatedDeskWatchlists(&st, defaultState())
	app := &Application{state: st, configDir: t.TempDir(), hub: NewHub(), sessionKey: "test"}
	app.engine = NewEngine(app)
	return app
}

// v15TestApp is a temporary call-site compatibility name inside the retained
// Extreme-30 matrix. Wave 4 removes this alias when the large matrix is
// capability-renamed/refactored without changing its assertions.
func v15TestApp(t *testing.T) *Application {
	return extreme30TestApp(t)
}

// setDeskBits provides deterministic desk-membership setup for the Extreme-30
// state-transition matrix.
func setDeskBits(t *testing.T, app *Application, sym, bits string) {
	t.Helper()
	if len(bits) != 3 {
		t.Fatalf("desk bit fixture %q must contain exactly three bits", bits)
	}
	app.mu.Lock()
	defer app.mu.Unlock()
	ensureDedicatedDeskWatchlists(&app.state, defaultState())
	for i, desk := range []string{"day", "swing", "long"} {
		setMembershipLocked(&app.state, desk, sym, bits[i] == '1')
	}
}
