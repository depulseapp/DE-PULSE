package main

import (
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestExtremeProbeRejectMalformedUserTickers(t *testing.T) {
	bad := []string{"AA@PL", "AAPL<script>", "AAPL/US", "AAPL \\n MSFT", "💥AAPL", "AAPL$"}
	for _, raw := range bad {
		app := v15TestApp(t)
		rr := httptest.NewRecorder()
		req := httptest.NewRequest("POST", "/api/master-symbol/add", strings.NewReader(`{"symbol":"`+strings.ReplaceAll(raw, "\\n", " ")+`"}`))
		req.Header.Set("Content-Type", "application/json")
		app.handleMasterSymbolAdd(rr, req)
		if rr.Code == 200 {
			t.Fatalf("master add accepted malformed raw ticker %q => %s", raw, rr.Body.String())
		}
	}
}

func TestExtremeProbeFutureTimestampNotFresh(t *testing.T) {
	now := time.Now().UnixMilli()
	ts, anomaly := safeFreshnessTimestamp(now+10*60*1000, 0, now)
	if !anomaly || ts != 0 {
		t.Fatalf("future timestamp not quarantined: ts=%d anomaly=%v", ts, anomaly)
	}
}

func TestExtremeProbeMutatingGETDoesNotMutate(t *testing.T) {
	app := v15TestApp(t)
	before := len(app.state.Watchlists)
	rr := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/api/watchlists/create", nil)
	app.handleWatchlistCreate(rr, req)
	if rr.Code == 200 && len(app.state.Watchlists) != before {
		t.Fatalf("GET mutated watchlists: before=%d after=%d body=%s", before, len(app.state.Watchlists), rr.Body.String())
	}
}
