package main

import (
	"strings"
	"testing"
	"time"
)

func v1604ContextFixture(now time.Time) (AppState, map[string]Quote, map[string]map[string][]Bar, map[string]FundamentalSnapshot, map[string]int64, map[string]string, ResearchPackageContext) {
	st, quotes, bars, fundamentals, last, health := completeResearchFixture(now)
	ms := now.UnixMilli()
	for _, sym := range []string{"SPY", "QQQ", "VIX"} {
		price := 500.0
		if sym == "QQQ" {
			price = 450
		}
		if sym == "VIX" {
			price = 18
		}
		quotes[sym] = Quote{Symbol: sym, Price: price, Source: "alpaca-iex", ProviderTimestamp: ms - 1200, UpdatedAt: ms - 800, DataState: "LIVE"}
	}
	last["catalyst-watch"] = ms - 1000
	ctx := ResearchPackageContext{Global: GlobalMarketContext{Tone: "NEUTRAL", Confidence: 80, Mode: "direct+fallback", Summary: "Canonical global context", UpdatedAt: ms - 1500}}
	return st, quotes, bars, fundamentals, last, health, ctx
}

func TestV1604LifecycleCanonicalStatesAndProvenance(t *testing.T) {
	now := time.Date(2026, 8, 11, 15, 30, 0, 0, time.UTC)
	ms := now.UnixMilli()
	bars := map[string]map[string][]Bar{
		"AAPL": {"daily": {{T: now.Add(-24 * time.Hour).Unix(), C: 200}}},
		"OLD":  {"daily": {{T: now.Add(-24 * time.Hour).Unix(), C: 10}}},
	}
	actions := []CorporateAction{
		{ID: "name-1", Symbol: "OLD", OldSymbol: "OLD", NewSymbol: "NEW", Type: "name_change", Status: "EFFECTIVE", ProcessDate: "2026-08-10", Source: "Alpaca", FirstSeenAt: ms - 5000, UpdatedAt: ms - 1000},
		{ID: "merge-1", Symbol: "M1", Type: "merger", Status: "UPCOMING", ProcessDate: "2026-08-20", Source: "Alpaca", UpdatedAt: ms - 900},
		{ID: "merge-2", Symbol: "M2", Type: "merger", Status: "EFFECTIVE", ProcessDate: "2026-08-10", Source: "Alpaca", UpdatedAt: ms - 800},
		{ID: "delist-1", Symbol: "D1", Type: "delisting", Status: "UPCOMING", ProcessDate: "2026-08-15", Source: "Alpaca", UpdatedAt: ms - 700},
		{ID: "delist-2", Symbol: "D2", Type: "delisting", Status: "EFFECTIVE", ProcessDate: "2026-08-09", Source: "Alpaca", UpdatedAt: ms - 600},
		{ID: "ipo-1", Symbol: "IPOX", Type: "new_listing", Status: "EFFECTIVE", ProcessDate: "2026-08-11", Source: "Alpaca", UpdatedAt: ms - 500},
	}
	truth := buildCorporateActionTruth(actions, bars, ms)
	want := map[string]string{"AAPL": "ACTIVE", "OLD": "NAME OR TICKER CHANGE", "NEW": "NAME OR TICKER CHANGE", "M1": "MERGER PENDING", "M2": "MERGED", "D1": "DELISTING PENDING", "D2": "DELISTED", "IPOX": "NEW LISTING / IPO"}
	for sym, state := range want {
		got, ok := truth.Lifecycle[sym]
		if !ok || got.State != state {
			t.Fatalf("%s lifecycle=%+v want %s", sym, got, state)
		}
		if got.Reason == "" || got.Source == "" || got.UpdatedAt == 0 {
			t.Fatalf("%s lifecycle lacks provenance/timing/reason: %+v", sym, got)
		}
	}
	if truth.SymbolLineage["OLD"] != "NEW" {
		t.Fatalf("lineage not preserved: %+v", truth.SymbolLineage)
	}
}

func TestV1604LifecycleUnknownAndMalformedSymbolsNeverFabricated(t *testing.T) {
	now := time.Date(2026, 8, 11, 15, 30, 0, 0, time.UTC).UnixMilli()
	bars := map[string]map[string][]Bar{"STALE": {"daily": {{T: time.UnixMilli(now).Add(-30 * 24 * time.Hour).Unix(), C: 1}}}}
	actions := []CorporateAction{{ID: "bad", Symbol: "NULL", OldSymbol: "NIL", NewSymbol: "NONE", Type: "name_change", Status: "EFFECTIVE", Source: "Alpaca"}}
	truth := buildCorporateActionTruth(actions, bars, now)
	if truth.Lifecycle["STALE"].State != "UNKNOWN" {
		t.Fatalf("old evidence should stay UNKNOWN: %+v", truth.Lifecycle["STALE"])
	}
	for _, bad := range []string{"NULL", "NIL", "NONE"} {
		if _, ok := truth.Lifecycle[bad]; ok {
			t.Fatalf("fabricated placeholder lifecycle symbol %s", bad)
		}
	}
}

func TestV1604LifecycleLedgerCorrectionAndCoverageTruthSurvive(t *testing.T) {
	now := time.Date(2026, 8, 11, 15, 30, 0, 0, time.UTC).UnixMilli()
	existing := []CorporateAction{{ID: "m", Symbol: "ABC", Type: "merger", Status: "UPCOMING", Source: "Alpaca", FirstSeenAt: now - 10000, UpdatedAt: now - 9000}}
	fresh := []CorporateAction{{ID: "m", Symbol: "ABC", Type: "merger", Status: "EFFECTIVE", Source: "Alpaca", UpdatedAt: now - 1000}}
	merged := mergeCorporateActionLedger(existing, fresh, now)
	cov := map[string]RawHistoryCoverage{"ABC": {Symbol: "ABC", State: "PARTIAL", BarCount: 2, PageCount: 1, PaginationComplete: false}}
	truth := buildCorporateActionTruth(merged, map[string]map[string][]Bar{"ABC": {"daily-raw": {{T: 1, C: 1}, {T: 2, C: 2}}}}, now, cov)
	if len(merged) != 1 || truth.Lifecycle["ABC"].State != "MERGED" {
		t.Fatalf("corrected action not canonical: merged=%+v lifecycle=%+v", merged, truth.Lifecycle["ABC"])
	}
	if truth.RawHistoryAvailable["ABC"] || truth.RawHistoryCoverage["ABC"].State != "PARTIAL" {
		t.Fatalf("partial history falsely complete: %+v", truth.RawHistoryCoverage["ABC"])
	}
}

func TestV1604ResearchRequiresCatalystAndMarketContext(t *testing.T) {
	now := time.Date(2026, 8, 11, 15, 30, 0, 0, time.UTC)
	st, quotes, bars, fundamentals, last, health, ctx := v1604ContextFixture(now)
	pkg := buildResearchPackageTruth(st, quotes, bars, fundamentals, nil, nil, nil, last, health, agreedResearchRec(now), now, ctx)
	for _, name := range []string{"Catalyst & Material Event Context", "Required Market Context"} {
		c, ok := researchComponent(pkg, name)
		if !ok || c.State != "FRESH" || !c.Required {
			t.Fatalf("%s missing/not fresh: %+v", name, c)
		}
	}
	if pkg.State != "FRESH" {
		t.Fatalf("complete v16.0.4 package should be FRESH: %+v", pkg)
	}
}

func TestV1604NoActiveCatalystCannotBeFreshWithStaleNewsCheck(t *testing.T) {
	now := time.Date(2026, 8, 11, 15, 30, 0, 0, time.UTC)
	st, quotes, bars, fundamentals, last, health, ctx := v1604ContextFixture(now)
	last["research-news:AAPL"] = now.Add(-2 * time.Hour).UnixMilli()
	pkg := buildResearchPackageTruth(st, quotes, bars, fundamentals, nil, nil, nil, last, health, agreedResearchRec(now), now, ctx)
	c, _ := researchComponent(pkg, "Catalyst & Material Event Context")
	if c.State == "FRESH" || pkg.State == "FRESH" || !strings.Contains(c.Detail, "News") {
		t.Fatalf("stale News allowed false no-catalyst freshness: %+v pkg=%s", c, pkg.State)
	}
}

func TestV1604ActiveCatalystCarriesMaterialEvidenceAndFutureSkewDegrades(t *testing.T) {
	now := time.Date(2026, 8, 11, 15, 30, 0, 0, time.UTC)
	st, quotes, bars, fundamentals, last, health, ctx := v1604ContextFixture(now)
	ctx.CatalystReactions = map[string]CatalystReactionState{"AAPL": {Symbol: "AAPL", TriggerType: "EARNINGS", Trigger: "Q3", State: "REACTION", Phase: "15m", TriggerAt: now.Add(-10 * time.Minute).UnixMilli(), UpdatedAt: now.Add(-time.Minute).UnixMilli(), Detail: "Gap holding above VWAP"}}
	pkg := buildResearchPackageTruth(st, quotes, bars, fundamentals, nil, nil, nil, last, health, agreedResearchRec(now), now, ctx)
	c, _ := researchComponent(pkg, "Catalyst & Material Event Context")
	if c.State != "FRESH" || !strings.Contains(c.Detail, "EARNINGS") || c.DataAt == 0 {
		t.Fatalf("active catalyst not exposed: %+v", c)
	}
	r := ctx.CatalystReactions["AAPL"]
	r.UpdatedAt = now.Add(5 * time.Minute).UnixMilli()
	ctx.CatalystReactions["AAPL"] = r
	pkg = buildResearchPackageTruth(st, quotes, bars, fundamentals, nil, nil, nil, last, health, agreedResearchRec(now), now, ctx)
	c, _ = researchComponent(pkg, "Catalyst & Material Event Context")
	if c.State != "STALE" || !strings.Contains(strings.ToLower(c.Detail), "future") {
		t.Fatalf("future catalyst timestamp accepted: %+v", c)
	}
}

func TestV1604MarketContextWorstRequiredDependencyAndWeekendSemantics(t *testing.T) {
	now := time.Date(2026, 8, 11, 15, 30, 0, 0, time.UTC)
	st, quotes, bars, fundamentals, last, health, ctx := v1604ContextFixture(now)
	delete(quotes, "VIX")
	pkg := buildResearchPackageTruth(st, quotes, bars, fundamentals, nil, nil, nil, last, health, agreedResearchRec(now), now, ctx)
	c, _ := researchComponent(pkg, "Required Market Context")
	if c.State == "FRESH" || pkg.State == "FRESH" || !strings.Contains(c.Detail, "VIX") {
		t.Fatalf("missing VIX allowed fully Fresh research: %+v pkg=%s", c, pkg.State)
	}

	weekend := time.Date(2026, 8, 15, 16, 0, 0, 0, time.UTC)
	st, quotes, bars, fundamentals, last, health, ctx = v1604ContextFixture(weekend)
	friday := time.Date(2026, 8, 14, 20, 0, 0, 0, time.UTC).UnixMilli()
	for _, sym := range []string{"AAPL", "SPY", "QQQ", "VIX"} {
		q := quotes[sym]
		q.ProviderTimestamp = friday
		q.UpdatedAt = friday
		quotes[sym] = q
	}
	ctx.Global.UpdatedAt = friday
	pkg = buildResearchPackageTruth(st, quotes, bars, fundamentals, nil, nil, nil, last, health, agreedResearchRec(weekend), weekend, ctx)
	c, _ = researchComponent(pkg, "Required Market Context")
	if c.State != "FRESH" {
		t.Fatalf("truthful recent Friday close was incorrectly stale on weekend: %+v", c)
	}
}

func TestV1604MarketContextFutureTimestampCannotBeFreshAndSourceIsVisible(t *testing.T) {
	now := time.Date(2026, 8, 11, 15, 30, 0, 0, time.UTC)
	st, quotes, bars, fundamentals, last, health, ctx := v1604ContextFixture(now)
	q := quotes["QQQ"]
	q.ProviderTimestamp = now.Add(5 * time.Minute).UnixMilli()
	quotes["QQQ"] = q
	pkg := buildResearchPackageTruth(st, quotes, bars, fundamentals, nil, nil, nil, last, health, agreedResearchRec(now), now, ctx)
	c, _ := researchComponent(pkg, "Required Market Context")
	if c.State != "STALE" || !strings.Contains(strings.ToLower(c.Detail), "future") || !strings.Contains(c.Source, "QQQ=") || !strings.Contains(c.Source, "Global=") {
		t.Fatalf("market context timestamp/source truth failed: %+v", c)
	}
}

func TestV1604EvidenceSnapshotChangesOnlyForMaterialEvidence(t *testing.T) {
	now := time.Date(2026, 8, 11, 15, 30, 0, 0, time.UTC)
	st, quotes, bars, fundamentals, last, health, ctx := v1604ContextFixture(now)
	pkg := buildResearchPackageTruth(st, quotes, bars, fundamentals, nil, nil, nil, last, health, agreedResearchRec(now), now, ctx)
	corp := buildCorporateActionTruth(nil, bars, now.UnixMilli())
	a := buildEvidenceSnapshot(pkg, agreedResearchRec(now), corp)
	// Wall-clock aging alone does not change semantic evidence because GeneratedAt is excluded from fingerprint.
	pkg2 := pkg
	pkg2.GeneratedAt = now.Add(time.Minute).UnixMilli()
	b := buildEvidenceSnapshot(pkg2, agreedResearchRec(now), corp)
	if a.ID != b.ID {
		t.Fatalf("clock aging churned snapshot: %s != %s", a.ID, b.ID)
	}

	acts := []CorporateAction{{ID: "d", Symbol: "AAPL", Type: "delisting", Status: "UPCOMING", Source: "Alpaca", UpdatedAt: now.UnixMilli()}}
	corp2 := buildCorporateActionTruth(acts, bars, now.UnixMilli())
	c := buildEvidenceSnapshot(pkg, agreedResearchRec(now), corp2)
	if c.ID == a.ID || c.SymbolLifecycle == nil || c.SymbolLifecycle.State != "DELISTING PENDING" {
		t.Fatalf("lifecycle transition did not change snapshot: a=%s c=%+v", a.ID, c)
	}

	ctx.CatalystReactions = map[string]CatalystReactionState{"AAPL": {Symbol: "AAPL", TriggerType: "NEWS", State: "TRIGGERED", TriggerAt: now.Add(-time.Minute).UnixMilli(), UpdatedAt: now.UnixMilli(), Detail: "material headline"}}
	pkg3 := buildResearchPackageTruth(st, quotes, bars, fundamentals, nil, nil, nil, last, health, agreedResearchRec(now), now, ctx)
	d := buildEvidenceSnapshot(pkg3, agreedResearchRec(now), corp)
	if d.ID == a.ID {
		t.Fatal("material catalyst did not change snapshot")
	}

	q := quotes["VIX"]
	q.Price = 35
	quotes["VIX"] = q
	pkg4 := buildResearchPackageTruth(st, quotes, bars, fundamentals, nil, nil, nil, last, health, agreedResearchRec(now), now, ctx)
	e := buildEvidenceSnapshot(pkg4, agreedResearchRec(now), corp)
	if e.ID == d.ID {
		t.Fatal("material required-market evidence did not change snapshot")
	}
}

func TestV1604NoTradeabilityOrDeterministicMutationInNewResearchComponents(t *testing.T) {
	now := time.Date(2026, 8, 11, 15, 30, 0, 0, time.UTC)
	st, quotes, bars, fundamentals, last, health, ctx := v1604ContextFixture(now)
	pkg := buildResearchPackageTruth(st, quotes, bars, fundamentals, nil, nil, nil, last, health, agreedResearchRec(now), now, ctx)
	for _, c := range pkg.Components {
		if strings.Contains(strings.ToLower(c.Dataset+" "+c.Detail), "tradeability") {
			t.Fatalf("v16.1 Tradeability leaked into v16.0.4: %+v", c)
		}
	}
}

func TestV1604AuthorizationLaterRelistingSupersedesOldDelisting(t *testing.T) {
	now := time.Date(2026, 8, 11, 15, 30, 0, 0, time.UTC)
	ms := now.UnixMilli()
	actions := []CorporateAction{
		{ID: "old-delist", Symbol: "ABC", Type: "delisting", Status: "EFFECTIVE", ProcessDate: "2025-06-01", Source: "Alpaca", UpdatedAt: ms - 1000},
		{ID: "relist", Symbol: "ABC", Type: "new_listing", Status: "EFFECTIVE", ProcessDate: "2026-08-01", Source: "Alpaca", UpdatedAt: ms - 900},
	}
	truth := buildCorporateActionTruth(actions, nil, ms)
	if got := truth.Lifecycle["ABC"]; got.State != "NEW LISTING / IPO" || got.EvidenceID != "relist" {
		t.Fatalf("later relisting did not supersede historical delisting: %+v", got)
	}
}

func TestV1604AuthorizationRecentTradingEvidenceSupersedesOlderCompletedEvent(t *testing.T) {
	now := time.Date(2026, 8, 11, 15, 30, 0, 0, time.UTC)
	ms := now.UnixMilli()
	bars := map[string]map[string][]Bar{"ABC": {"daily": {{T: now.Add(-24 * time.Hour).Unix(), C: 25}}}}
	actions := []CorporateAction{{ID: "old-delist", Symbol: "ABC", Type: "delisting", Status: "EFFECTIVE", ProcessDate: "2025-06-01", Source: "Alpaca", UpdatedAt: ms - 500}}
	truth := buildCorporateActionTruth(actions, bars, ms)
	if got := truth.Lifecycle["ABC"]; got.State != "ACTIVE" {
		t.Fatalf("newer affirmative trading evidence should supersede historical completed event: %+v", got)
	}
}

func TestV1604AuthorizationSelectedTickerEarningsEvidenceIsExplicit(t *testing.T) {
	now := time.Date(2026, 8, 11, 15, 30, 0, 0, time.UTC)
	st, quotes, bars, fundamentals, last, health, ctx := v1604ContextFixture(now)
	est := 1.25
	earnings := []EarningsItem{{Symbol: "AAPL", Date: "2026-08-11", Hour: "AMC", EPSEstimate: &est}}
	pkg := buildResearchPackageTruth(st, quotes, bars, fundamentals, nil, earnings, nil, last, health, agreedResearchRec(now), now, ctx)
	c, _ := researchComponent(pkg, "Catalyst & Material Event Context")
	if c.State != "FRESH" || !strings.Contains(c.Detail, "Earnings 2026-08-11 AMC") || !strings.Contains(strings.ToLower(c.Detail), "scheduled earnings risk") {
		t.Fatalf("selected-ticker earnings evidence not explicit: %+v", c)
	}
}

func TestV1604AuthorizationUnexpectedMaterialNewsWithoutReactionIsNotFresh(t *testing.T) {
	now := time.Date(2026, 8, 11, 15, 30, 0, 0, time.UTC)
	st, quotes, bars, fundamentals, last, health, ctx := v1604ContextFixture(now)
	news := []NewsItem{{Datetime: now.Add(-5 * time.Minute).Unix(), Headline: "AAPL announces acquisition", Source: "Finnhub", Symbols: []string{"AAPL"}}}
	pkg := buildResearchPackageTruth(st, quotes, bars, fundamentals, news, nil, nil, last, health, agreedResearchRec(now), now, ctx)
	c, _ := researchComponent(pkg, "Catalyst & Material Event Context")
	if c.State == "FRESH" || !strings.Contains(c.Detail, "Material news") || !strings.Contains(c.Detail, "without a matching Catalyst Watch reaction") {
		t.Fatalf("material news mismatch allowed false Fresh catalyst truth: %+v", c)
	}
}

func TestV1604AuthorizationCatalystWatchCheckIsRequiredForNoCatalystFresh(t *testing.T) {
	now := time.Date(2026, 8, 11, 15, 30, 0, 0, time.UTC)
	st, quotes, bars, fundamentals, last, health, ctx := v1604ContextFixture(now)
	last["catalyst-watch"] = now.Add(-3 * time.Hour).UnixMilli()
	pkg := buildResearchPackageTruth(st, quotes, bars, fundamentals, nil, nil, nil, last, health, agreedResearchRec(now), now, ctx)
	c, _ := researchComponent(pkg, "Catalyst & Material Event Context")
	if c.State == "FRESH" || !strings.Contains(c.Detail, "Catalyst Watch") {
		t.Fatalf("stale Catalyst Watch check allowed false Fresh no-catalyst state: %+v", c)
	}
}

func TestV1604StopTimeoutPersistsCurrentCacheBeforeReturning(t *testing.T) {
	app := newTestApplication(t)
	now := time.Now().UnixMilli()
	app.engine.mu.Lock()
	app.engine.quotes["SPY"] = Quote{Symbol: "SPY", Price: 600, PreviousClose: 599, UpdatedAt: now, ProviderTimestamp: now}
	app.engine.bars["SPY"] = map[string][]Bar{"daily": {{T: time.Now().Unix(), O: 598, H: 601, L: 597, C: 600, V: 1000}}}
	app.engine.status = "running"
	app.engine.cancel = func() {}
	app.engine.mu.Unlock()

	oldTimeout := runtimeStopTimeout
	runtimeStopTimeout = 15 * time.Millisecond
	defer func() { runtimeStopTimeout = oldTimeout }()

	// Simulate a worker that outlives the synchronous stop timeout.
	app.engine.wg.Add(1)
	app.engine.Stop()
	app.engine.mu.RLock()
	status := app.engine.status
	app.engine.mu.RUnlock()
	if status != "stopping" {
		t.Fatalf("timed-out stop status=%q want stopping", status)
	}

	reloaded := NewEngine(app)
	if q := reloaded.quotes["SPY"]; q.Price != 600 {
		t.Fatalf("pre-timeout cache snapshot did not preserve quote: %+v", q)
	}
	if len(reloaded.bars["SPY"]["daily"]) != 1 {
		t.Fatalf("pre-timeout cache snapshot did not preserve bars: %+v", reloaded.bars["SPY"])
	}

	app.engine.wg.Done()
	deadline := time.Now().Add(time.Second)
	for time.Now().Before(deadline) {
		app.engine.mu.RLock()
		status = app.engine.status
		app.engine.mu.RUnlock()
		if status == "stopped" {
			return
		}
		time.Sleep(5 * time.Millisecond)
	}
	t.Fatalf("runtime did not finalize after simulated worker exit; status=%q", status)
}

func TestV1604QuoteLevelCatalystTrackingIsEventDriven(t *testing.T) {
	if catalystQuoteReactionActive(nil) {
		t.Fatal("quote-level Catalyst Watch must remain dormant with no confirmed active trigger")
	}
	completed := map[string]CatalystReactionState{
		"AAA": {Symbol: "AAA", TriggerType: "NEWS", TriggerAt: 1000, CompletedAt: 2000},
	}
	if catalystQuoteReactionActive(completed) {
		t.Fatal("completed catalyst must not keep quote-level reaction tracking active")
	}
	active := map[string]CatalystReactionState{
		"AAA": {Symbol: "AAA", TriggerType: "EARNINGS", TriggerAt: 1000},
	}
	if !catalystQuoteReactionActive(active) {
		t.Fatal("confirmed incomplete catalyst must keep quote-level reaction tracking active")
	}
}
