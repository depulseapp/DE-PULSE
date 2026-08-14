package main

import (
	"testing"
	"time"
)

func TestV151DataFreshnessCadenceMatrix(t *testing.T) {
	cases := []struct {
		dataset, provider, session string
		cadence, fresh, stale      time.Duration
	}{
		{"Quotes", "Alpaca", "regular", 30 * time.Second, 90 * time.Second, 2 * time.Minute},
		{"VIX", "Twelve Data", "regular", 2 * time.Minute, 5 * time.Minute, 10 * time.Minute},
		{"VIX", "CBOE", "after-hours", 10 * time.Minute, 30 * time.Minute, time.Hour},
		{"Intraday Bars", "Alpaca", "regular", 5 * time.Minute, 7 * time.Minute, 10 * time.Minute},
		{"Daily / Weekly History", "Alpaca", "regular", 24 * time.Hour, 30 * time.Hour, 72 * time.Hour},
		{"News", "Finnhub", "regular", 5 * time.Minute, 7 * time.Minute, 15 * time.Minute},
		{"News", "Finnhub", "after-hours", 10 * time.Minute, 15 * time.Minute, 25 * time.Minute},
		{"News", "Finnhub", "closed", 30 * time.Minute, 40 * time.Minute, 90 * time.Minute},
		{"SEC Filings", "SEC EDGAR", "regular", 15 * time.Minute, 20 * time.Minute, 35 * time.Minute},
		{"Fundamentals", "Finnhub", "regular", 24 * time.Hour, 26 * time.Hour, 36 * time.Hour},
		{"Global", "Twelve Data", "regular", 7 * time.Minute, 10 * time.Minute, 20 * time.Minute},
		{"Macro", "FRED", "regular", 6 * time.Hour, 12 * time.Hour, 30 * time.Hour},
		{"Options", "Alpaca", "regular", 3 * time.Minute, 5 * time.Minute, 8 * time.Minute},
		{"Options", "Alpaca", "closed", 10 * time.Minute, 15 * time.Minute, 30 * time.Minute},
	}
	for _, tc := range cases {
		c, f, s := freshnessLimits(tc.dataset, tc.provider, tc.session)
		if c != tc.cadence || f != tc.fresh || s != tc.stale {
			t.Fatalf("%s/%s cadence=%v/%v/%v want=%v/%v/%v", tc.dataset, tc.session, c, f, s, tc.cadence, tc.fresh, tc.stale)
		}
	}
}

func TestV151SparseFreshnessUsesCheckAgeNotLatestItemAge(t *testing.T) {
	app := v15TestApp(t)
	e := app.engine
	now := time.Now().UnixMilli()
	e.mu.Lock()
	// Latest article is old, but Finnhub reconciliation just succeeded.
	e.news = []NewsItem{{Headline: "No newer headline", Source: "Finnhub", Datetime: (now - 77*60*1000) / 1000}}
	e.filings = []FilingItem{{Symbol: "AAPL", Form: "10-Q", FiledAt: time.Now().Add(-48 * time.Hour).UTC().Format(time.RFC3339)}}
	e.lastUpdated["news"] = now - time.Minute.Milliseconds()
	e.lastUpdated["filings"] = now - 8*time.Minute.Milliseconds()
	e.lastUpdated["history-intraday"] = now - 2*time.Minute.Milliseconds()
	e.lastUpdated["history-daily"] = now - 3*time.Hour.Milliseconds()
	e.bars["AAPL"] = map[string][]Bar{
		"intraday": {{T: time.Now().Add(-2 * time.Minute).Unix(), O: 100, H: 101, L: 99, C: 100, V: 1000}},
		"daily":    {{T: time.Now().Add(-24 * time.Hour).Unix(), O: 99, H: 101, L: 98, C: 100, V: 10000}},
	}
	e.mu.Unlock()
	rows, _ := e.buildFreshnessDiagnostics(map[string]Quote{}, clone(e.lastUpdated), map[string]string{
		"news": "healthy · Finnhub", "filings": "healthy · SEC EDGAR", "history": "healthy · Alpaca",
	})
	by := map[string]FreshnessDiagnostic{}
	for _, r := range rows {
		by[r.Dataset] = r
	}
	news := by["News"]
	if news.State == "STALE" || news.CheckAgeMs > 2*time.Minute.Milliseconds() || news.DataAgeMs < 70*time.Minute.Milliseconds() {
		t.Fatalf("news check/data age semantics wrong: %+v", news)
	}
	if news.Action != "news" || news.FreshnessBasis != "last successful news check" {
		t.Fatalf("news semantics metadata wrong: %+v", news)
	}
	sec := by["SEC Filings"]
	if sec.State == "STALE" || sec.CheckAgeMs > 10*time.Minute.Milliseconds() || sec.DataAgeMs < 24*time.Hour.Milliseconds() {
		t.Fatalf("SEC check/data age semantics wrong: %+v", sec)
	}
	if by["Intraday Bars"].Action != "history-intraday" || by["Daily / Weekly History"].Action != "history-daily" {
		t.Fatalf("history refresh actions are not independently targeted: intraday=%q daily=%q", by["Intraday Bars"].Action, by["Daily / Weekly History"].Action)
	}
}

func TestV151HistoryModesAndOptionsEventCadence(t *testing.T) {
	app := v15TestApp(t)
	e := app.engine
	app.mu.Lock()
	ensureDedicatedDeskWatchlists(&app.state, defaultState())
	setMembershipLocked(&app.state, "day", "AAPL", true)
	setMembershipLocked(&app.state, "swing", "AAPL", true)
	setMembershipLocked(&app.state, "long", "AAPL", true)
	st := clone(app.state)
	app.mu.Unlock()
	intraday := historySpecsForStateMode(st, nil, "intraday")
	daily := historySpecsForStateMode(st, nil, "daily")
	if len(intraday) == 0 || len(daily) == 0 {
		t.Fatalf("history mode specs missing: intraday=%v daily=%v", intraday, daily)
	}
	for _, sp := range intraday {
		if sp.Name != "intraday" {
			t.Fatalf("intraday mode leaked %s", sp.Name)
		}
	}
	for _, sp := range daily {
		if sp.Name == "intraday" {
			t.Fatalf("daily mode leaked intraday spec")
		}
	}
	now := time.Now()
	e.mu.Lock()
	e.catalystReactions["AAPL"] = CatalystReactionState{Symbol: "AAPL", TriggerAt: now.Add(-time.Minute).UnixMilli()}
	e.mu.Unlock()
	if got := e.optionsRefreshInterval(now); got != time.Minute {
		t.Fatalf("active catalyst options cadence=%v want 1m", got)
	}
}

func TestV151FreshnessPriorityOrdering(t *testing.T) {
	if freshnessPriority("Quotes") != 1 || freshnessPriority("VIX") != 1 || freshnessPriority("Intraday Bars") != 1 {
		t.Fatal("priority-1 live datasets changed")
	}
	if freshnessPriority("News") != 2 || freshnessPriority("SEC Filings") != 2 || freshnessPriority("Earnings") != 2 {
		t.Fatal("priority-2 catalyst datasets changed")
	}
	if freshnessPriority("Daily / Weekly History") != 3 || freshnessPriority("Fundamentals") != 3 || freshnessPriority("Macro") != 3 {
		t.Fatal("priority-3 contextual datasets changed")
	}
}
