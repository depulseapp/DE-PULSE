package main

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestV1511CatalystAfterHoursPersistsThroughNextTradingSession(t *testing.T) {
	loc := easternLocation()
	trigger := time.Date(2026, 8, 10, 16, 5, 0, 0, loc) // Monday AMC
	if got := catalystPhase(time.Date(2026, 8, 11, 10, 10, 0, 0, loc), trigger.UnixMilli()); got == "COMPLETE" {
		t.Fatalf("AMC catalyst completed before next-session reaction finished: %s", got)
	}
	if got := catalystPhase(time.Date(2026, 8, 11, 10, 35, 0, 0, loc), trigger.UnixMilli()); got != "60m" {
		t.Fatalf("expected 60m next-session phase, got %s", got)
	}
	if got := catalystPhase(time.Date(2026, 8, 11, 16, 16, 0, 0, loc), trigger.UnixMilli()); got != "COMPLETE" {
		t.Fatalf("expected completion after next session close, got %s", got)
	}
}

func TestV1511CatalystFridayAfterHoursPersistsThroughMonday(t *testing.T) {
	loc := easternLocation()
	trigger := time.Date(2026, 8, 7, 16, 10, 0, 0, loc) // Friday AMC
	sunday := time.Date(2026, 8, 9, 18, 0, 0, 0, loc)
	if catalystComplete(sunday, trigger.UnixMilli()) {
		t.Fatal("Friday AMC catalyst must survive weekend")
	}
	mondayOpen := time.Date(2026, 8, 10, 9, 40, 0, 0, loc)
	if got := catalystPhase(mondayOpen, trigger.UnixMilli()); got != "5m" {
		t.Fatalf("expected Monday opening reaction, got %s", got)
	}
	mondayDone := time.Date(2026, 8, 10, 16, 16, 0, 0, loc)
	if !catalystComplete(mondayDone, trigger.UnixMilli()) {
		t.Fatal("Friday AMC catalyst should complete after Monday session")
	}
}

func TestV1511ReadinessFreshnessGateRejectsStaleCritical(t *testing.T) {
	now := time.Now()
	rows := []FreshnessDiagnostic{
		{Dataset: "Quotes", State: "LIVE", Provider: "Alpaca"},
		{Dataset: "VIX", State: "FRESH", Provider: "Twelve Data"},
		{Dataset: "Intraday Bars", State: "STALE", Provider: "Alpaca", Reason: "too old"},
	}
	usable, degraded, ex := readinessFreshnessGate(rows, []string{"Quotes", "VIX", "Intraday Bars"}, now)
	if usable != 2 || !degraded || len(ex) != 1 || ex[0].Severity != "HIGH" {
		t.Fatalf("unexpected gate result usable=%d degraded=%v ex=%+v", usable, degraded, ex)
	}
}

func TestV1511ReadinessFreshnessGateAllowsTruthfulDelayedWithCaution(t *testing.T) {
	now := time.Now()
	rows := []FreshnessDiagnostic{
		{Dataset: "Quotes", State: "LIVE", Provider: "Alpaca"},
		{Dataset: "VIX", State: "DELAYED", Provider: "CBOE", Reason: "official delayed quote"},
	}
	usable, degraded, ex := readinessFreshnessGate(rows, []string{"Quotes", "VIX"}, now)
	if usable != 2 || degraded || len(ex) != 1 || ex[0].Severity != "MEDIUM" || checkpointAttention(ex, degraded) != "READY WITH CAUTION" {
		t.Fatalf("unexpected delayed gate result usable=%d degraded=%v ex=%+v", usable, degraded, ex)
	}
}

func TestV1511FreshnessRecoveryIncludesUnavailableAndError(t *testing.T) {
	now := time.Now().UnixMilli()
	for _, st := range []string{"STALE", "ERROR", "UNAVAILABLE"} {
		if !freshnessRecoveryDue(FreshnessDiagnostic{State: st, Action: "news"}, now) {
			t.Fatalf("%s should trigger auto recovery", st)
		}
	}
	if freshnessRecoveryDue(FreshnessDiagnostic{State: "DELAYED", Action: "vix", NextExpectedAt: now + 60000}, now) {
		t.Fatal("truthful delayed source must wait until its next expected check")
	}
	if freshnessRecoveryDue(FreshnessDiagnostic{State: "IDLE", Action: "news"}, now) {
		t.Fatal("IDLE must not trigger auto recovery")
	}
	if !freshnessRecoveryDue(FreshnessDiagnostic{State: "FRESH", Action: "news", NextExpectedAt: now - 1}, now) {
		t.Fatal("past due FRESH row should proactively reconcile")
	}
}

func TestV1511Form4ParsesNonDerivativeAndDerivativeSemantics(t *testing.T) {
	xmlBody := `<?xml version="1.0"?><ownershipDocument>
      <reportingOwner><reportingOwnerId><rptOwnerName>Jane Trader</rptOwnerName></reportingOwnerId><reportingOwnerRelationship><isOfficer>1</isOfficer><officerTitle>CFO</officerTitle></reportingOwnerRelationship></reportingOwner>
      <nonDerivativeTable><nonDerivativeTransaction><transactionDate><value>2026-08-07</value></transactionDate><transactionCoding><transactionCode>P</transactionCode></transactionCoding><transactionAmounts><transactionShares><value>1000</value></transactionShares><transactionPricePerShare><value>25.50</value></transactionPricePerShare><transactionAcquiredDisposedCode><value>A</value></transactionAcquiredDisposedCode></transactionAmounts><postTransactionAmounts><sharesOwnedFollowingTransaction><value>12000</value></sharesOwnedFollowingTransaction></postTransactionAmounts></nonDerivativeTransaction></nonDerivativeTable>
      <derivativeTable><derivativeTransaction><transactionDate><value>2026-08-07</value></transactionDate><transactionCoding><transactionCode>M</transactionCode></transactionCoding><transactionAmounts><transactionShares><value>500</value></transactionShares><transactionPricePerShare><value>10.00</value></transactionPricePerShare><transactionAcquiredDisposedCode><value>A</value></transactionAcquiredDisposedCode></transactionAmounts><postTransactionAmounts><valueOwnedFollowingTransaction><value>2000</value></valueOwnedFollowingTransaction></postTransactionAmounts></derivativeTransaction></derivativeTable>
    </ownershipDocument>`
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { _, _ = w.Write([]byte(xmlBody)) }))
	defer srv.Close()
	item := FilingItem{Symbol: "AAA", Form: "4", FiledAt: "2026-08-08", URL: srv.URL}
	enrichForm4(context.Background(), srv.Client(), nil, srv.URL, &item)
	if len(item.Transactions) != 2 {
		t.Fatalf("expected two economic legs, got %+v", item.Transactions)
	}
	if item.Transactions[0].Classification != "BUY" || item.Transactions[0].Actor != "Jane Trader" || item.Transactions[0].Shares != 1000 {
		t.Fatalf("purchase leg not enriched: %+v", item.Transactions[0])
	}
	if item.Transactions[1].Classification != "OTHER" || item.Transactions[1].Code != "M" || !strings.Contains(item.Transactions[1].Meaning, "exercise") {
		t.Fatalf("derivative exercise must remain OTHER: %+v", item.Transactions[1])
	}
	if !strings.EqualFold(item.Signal, "Buy") || !strings.Contains(item.Meaning, "Other") {
		t.Fatalf("filing-level summary should retain genuine buy plus OTHER context: signal=%s meaning=%s", item.Signal, item.Meaning)
	}
}

func TestV1511SECIntelligenceHasUpdatedAtAndLargerTransactionWindow(t *testing.T) {
	tx := make([]SECInsiderTransaction, 0, 40)
	for i := 0; i < 40; i++ {
		tx = append(tx, SECInsiderTransaction{Classification: "OTHER", TransactionDate: "2026-08-01", Shares: 1})
	}
	got := buildSECIntelligence("AAA", []FilingItem{{Symbol: "AAA", Form: "4", Category: "insider", FiledAt: "2026-08-02", Transactions: tx}})
	if got.UpdatedAt == 0 {
		t.Fatal("SEC summary must expose reconciliation timestamp")
	}
	if len(got.RecentInsiderTransactions) != 40 {
		t.Fatalf("expected all 40 enriched rows retained, got %d", len(got.RecentInsiderTransactions))
	}
}

func TestV1511ResearchReadyUsesTruthfulTargetEvidence(t *testing.T) {
	app := newTestApplication(t)
	sym := "AAA"
	now := time.Date(2026, 8, 10, 11, 0, 0, 0, easternLocation())
	ms := now.UnixMilli()
	app.engine.mu.Lock()
	app.engine.quotes[sym] = Quote{Symbol: sym, Price: 100, ProviderTimestamp: ms, UpdatedAt: ms, Source: "Alpaca"}
	app.engine.quotes["VIX"] = Quote{Symbol: "VIX", Price: 16, ProviderTimestamp: ms, UpdatedAt: ms, Source: "Twelve Data"}
	app.engine.bars[sym] = map[string][]Bar{
		"intraday": {{T: now.Unix(), O: 99, H: 101, L: 98, C: 100, V: 1000}},
		"daily":    {{T: now.Unix(), O: 95, H: 101, L: 94, C: 100, V: 10000}},
		"weekly":   {{T: now.Unix(), O: 90, H: 102, L: 89, C: 100, V: 50000}},
	}
	app.engine.fundamentals[sym] = FundamentalSnapshot{Symbol: sym, RevenueGrowth: 10, UpdatedAt: ms, Source: "Finnhub"}
	app.engine.lastUpdated["quotes"] = ms
	app.engine.lastUpdated["history"] = ms
	app.engine.lastUpdated["history-intraday"] = ms
	app.engine.lastUpdated["history-daily"] = ms
	app.engine.lastUpdated["research-news:"+sym] = ms
	app.engine.lastUpdated["research-earnings:"+sym] = ms
	app.engine.lastUpdated["research-fundamentals:"+sym] = ms
	app.engine.lastUpdated["research-sec:"+sym] = ms
	app.engine.health["history"] = "healthy · Alpaca"
	app.engine.health["news"] = "stale · universe overdue"
	app.engine.health["earnings"] = "stale · universe overdue"
	app.engine.health["fundamentals"] = "stale · universe overdue"
	app.engine.health["filings"] = "stale · universe overdue"
	app.engine.mu.Unlock()
	ready, issues := app.engine.researchPackageReadinessAt(sym, now)
	if !ready || len(issues) != 0 {
		t.Fatalf("target-scoped current evidence should be research-ready: ready=%v issues=%v", ready, issues)
	}
	app.engine.mu.Lock()
	app.engine.bars[sym]["intraday"][0].T = now.Add(-3 * time.Hour).Unix()
	app.engine.lastUpdated["history-intraday"] = now.Add(-3 * time.Hour).UnixMilli()
	app.engine.mu.Unlock()
	ready, issues = app.engine.researchPackageReadinessAt(sym, now)
	if ready || len(issues) == 0 {
		t.Fatalf("stale critical history must block AI-ready state: ready=%v issues=%v", ready, issues)
	}
}

func TestV1511CompletedCatalystIsActuallyFinalized(t *testing.T) {
	loc := easternLocation()
	trigger := time.Date(2026, 8, 10, 16, 5, 0, 0, loc)
	now := time.Date(2026, 8, 11, 16, 16, 0, 0, loc)
	got := finalizeCompletedCatalysts(map[string]CatalystReactionState{"AAA": {Symbol: "AAA", TriggerAt: trigger.UnixMilli(), Phase: "SESSION REACTION"}}, now)
	if got["AAA"].Phase != "COMPLETE" || got["AAA"].CompletedAt == 0 {
		t.Fatalf("completed catalyst remained active: %+v", got["AAA"])
	}
}

func TestV1511CBOECloseRemainsTruthfulDelayedOutsideSession(t *testing.T) {
	now := time.Now().UnixMilli()
	st, reason, _, _, _ := freshnessStateWithLimits("VIX", "CBOE", "after-hours", now-int64(2*time.Hour/time.Millisecond), "healthy", now, 10*time.Minute, 30*time.Minute, 60*time.Minute)
	if st != "DELAYED" || !strings.Contains(reason, "not live") {
		t.Fatalf("expected official close delayed semantics, got %s %s", st, reason)
	}
}

func TestV1511CatalystIntradayTriggerMeasuresFromTriggerNotOpen(t *testing.T) {
	loc := easternLocation()
	trigger := time.Date(2026, 8, 10, 13, 0, 0, 0, loc)
	if got := catalystPhase(time.Date(2026, 8, 10, 13, 3, 0, 0, loc), trigger.UnixMilli()); got != "OPENING REACTION" {
		t.Fatalf("intraday catalyst should still be in opening reaction 3m after trigger, got %s", got)
	}
	if got := catalystPhase(time.Date(2026, 8, 10, 13, 7, 0, 0, loc), trigger.UnixMilli()); got != "5m" {
		t.Fatalf("intraday catalyst should use 5m-from-trigger phase, got %s", got)
	}
	if got := catalystPhase(time.Date(2026, 8, 10, 13, 20, 0, 0, loc), trigger.UnixMilli()); got != "15m" {
		t.Fatalf("intraday catalyst should use 15m-from-trigger phase, got %s", got)
	}
}

func TestV1511ResearchReadinessRejectsStaleTargetQuoteEvenWhenUniverseFresh(t *testing.T) {
	app := newTestApplication(t)
	sym := "AAA"
	now := time.Date(2026, 8, 10, 11, 0, 0, 0, easternLocation())
	fresh := now.UnixMilli()
	stale := now.Add(-4 * time.Hour).UnixMilli()
	app.engine.mu.Lock()
	app.engine.quotes[sym] = Quote{Symbol: sym, Price: 100, ProviderTimestamp: stale, UpdatedAt: stale, Source: "Alpaca"}
	app.engine.quotes["SPY"] = Quote{Symbol: "SPY", Price: 500, ProviderTimestamp: fresh, UpdatedAt: fresh, Source: "Alpaca"}
	app.engine.quotes["VIX"] = Quote{Symbol: "VIX", Price: 16, ProviderTimestamp: fresh, UpdatedAt: fresh, Source: "Twelve Data"}
	app.engine.bars[sym] = map[string][]Bar{
		"intraday": {{T: now.Unix(), O: 99, H: 101, L: 98, C: 100, V: 1000}},
		"daily":    {{T: now.Unix(), O: 95, H: 101, L: 94, C: 100, V: 10000}},
		"weekly":   {{T: now.Unix(), O: 90, H: 102, L: 89, C: 100, V: 50000}},
	}
	app.engine.fundamentals[sym] = FundamentalSnapshot{Symbol: sym, RevenueGrowth: 10, UpdatedAt: fresh, Source: "Finnhub"}
	for _, k := range []string{"quotes", "history", "history-intraday", "history-daily", "research-news:" + sym, "research-earnings:" + sym, "research-fundamentals:" + sym, "research-sec:" + sym} {
		app.engine.lastUpdated[k] = fresh
	}
	app.engine.health["quotes"] = "healthy · Alpaca"
	app.engine.health["history"] = "healthy · Alpaca"
	app.engine.health["news"] = "healthy · Finnhub"
	app.engine.health["earnings"] = "healthy · Finnhub"
	app.engine.health["fundamentals"] = "healthy · Finnhub"
	app.engine.health["filings"] = "healthy · SEC EDGAR"
	app.engine.mu.Unlock()
	ready, issues := app.engine.researchPackageReadinessAt(sym, now)
	if ready || len(issues) == 0 {
		t.Fatalf("stale target quote must block research even when another symbol keeps universe Quotes fresh: ready=%v issues=%v", ready, issues)
	}
}

func TestV1511ClosedSessionVeryOldQuoteIsNotIndefinitelyIdle(t *testing.T) {
	now := time.Now().UnixMilli()
	state, _, _, _, _ := freshnessStateWithLimits("Quotes", "Alpaca", "weekend", now-int64(10*24*time.Hour/time.Millisecond), "healthy", now, 5*time.Minute, 15*time.Minute, 30*time.Minute)
	if state == "IDLE" {
		t.Fatalf("a 10-day-old quote must not remain IDLE merely because market is closed")
	}
}

func TestV1511TargetedRefreshReportsProviderFailureTruthfully(t *testing.T) {
	app := newTestApplication(t)
	app.engine.mu.Lock()
	app.engine.mode = "live"
	app.engine.mu.Unlock()
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	if app.engine.refreshDatasetRouted(ctx, "news", Secrets{}) {
		t.Fatal("news routed refresh without Finnhub/Marketaux credentials must report failure")
	}
}
