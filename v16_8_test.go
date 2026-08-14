package main

import (
	"math"
	"strings"
	"testing"
	"time"
)

func v168Quote(sym string, price, change, bid, ask float64, now time.Time) Quote {
	return Quote{Symbol: sym, Price: price, ChangePercent: change, Bid: bid, Ask: ask, Source: "Alpaca", ProviderTimestamp: now.UnixMilli(), UpdatedAt: now.UnixMilli()}
}

func TestV168HeatMapOriginalAcceptanceModesReuseCanonicalQuotes(t *testing.T) {
	now := time.Now()
	st := defaultState()
	st.UI.SelectedTicker = "NVDA"
	quotes := map[string]Quote{}
	for _, sym := range []string{"SPY", "QQQ", "NVDA", "META", "ORCL", "PLTR", "TSLA", "SOFI", "XLK", "XLC", "XLY", "XLP", "XLE", "XLF", "XLV", "XLI", "XLB", "XLRE", "XLU"} {
		quotes[sym] = v168Quote(sym, 100, 0.5, 99.9, 100.1, now)
	}
	h := v165HeatMapForState(st, quotes, now)
	for _, k := range []string{"sector", "watchlist", "broad"} {
		if _, ok := h.Modes[k]; !ok {
			t.Fatalf("missing heat mode %s: %+v", k, h.Modes)
		}
	}
	if h.Modes["watchlist"].Expected == 0 || h.Modes["watchlist"].Fresh == 0 {
		t.Fatalf("watchlist heat not populated: %+v", h.Modes["watchlist"])
	}
	if !strings.Contains(strings.ToLower(h.Modes["broad"].CoverageBasis), "not a claim of full 500") {
		t.Fatalf("broad/S&P proxy coverage truth missing: %q", h.Modes["broad"].CoverageBasis)
	}
	if h.Universe != h.Modes["sector"].Universe {
		t.Fatal("legacy sector owner drifted")
	}
}

func TestV168SeasonalityTenYearStatisticsAndCurrentYearComparison(t *testing.T) {
	now := time.Date(2026, time.August, 12, 20, 0, 0, 0, time.UTC)
	bars := map[string]map[string][]Bar{}
	for _, sym := range []string{"SPY", "QQQ"} {
		rows := []Bar{}
		price := 100.0
		start := time.Date(2015, time.December, 1, 21, 0, 0, 0, time.UTC)
		for d := start; d.Before(now); d = d.AddDate(0, 0, 1) {
			if d.Weekday() == time.Saturday || d.Weekday() == time.Sunday {
				continue
			}
			price *= 1 + 0.00015 + float64(int(d.Month())-6)*0.00001
			rows = append(rows, Bar{T: d.Unix(), O: price * .999, H: price * 1.01, L: price * .99, C: price, V: 2_000_000})
		}
		bars[sym] = map[string][]Bar{"daily": rows}
	}
	snap := buildSeasonalitySnapshot(bars, now)
	spy := snap.Symbols["SPY"]
	if spy.State != "AVAILABLE" {
		t.Fatalf("seasonality unavailable: %+v", spy)
	}
	jan := spy.Monthly[0]
	if jan.SampleCount < 8 || jan.AverageReturnPct == nil || jan.MedianReturnPct == nil || jan.PositiveFrequencyPct == nil || jan.BestReturnPct == nil || jan.WorstReturnPct == nil {
		t.Fatalf("10-year stats incomplete: %+v", jan)
	}
	aug := spy.Monthly[7]
	if aug.CurrentYearReturnPct == nil || aug.CurrentYearObservation != "MONTH-TO-DATE" {
		t.Fatalf("current-year comparison missing: %+v", aug)
	}
	if jan.HistoricalYears > 10 {
		t.Fatalf("published sample exceeds trailing 10 years: %+v", jan)
	}
}

func TestV168GEXMajorStrikesZonesExpirationAndFlipTruth(t *testing.T) {
	strikes := map[float64]gexAccumulator{95: {Call: 100, Put: -500, OI: 1000, Contracts: 10}, 100: {Call: 200, Put: -300, OI: 1200, Contracts: 12}, 105: {Call: 700, Put: -100, OI: 1500, Contracts: 15}, 110: {Call: 800, Put: -50, OI: 1300, Contracts: 12}}
	exps := map[string]gexAccumulator{"2026-08-21": {Call: 900, Put: -600, OI: 3000, Contracts: 25}, "2026-09-18": {Call: 900, Put: -350, OI: 2000, Contracts: 20}}
	major, zones, expiry, flip, method := finalizeGEXStructure(strikes, exps, 102)
	if len(major) == 0 || len(zones) == 0 || len(expiry) != 2 {
		t.Fatalf("GEX structure incomplete major=%v zones=%v expiry=%v", major, zones, expiry)
	}
	if flip == nil || *flip <= 100 || *flip >= 105 || !strings.Contains(strings.ToLower(method), "not measured dealer") {
		t.Fatalf("flip-style truth incomplete flip=%v method=%q", flip, method)
	}
}

func TestV168LiquidityRelativeSpreadDollarVolumeOpeningAndVolatilityAdjustedRisk(t *testing.T) {
	loc := easternLocation()
	now := time.Date(2026, time.August, 12, 9, 45, 0, 0, loc)
	sym := "NVDA"
	q := v168Quote(sym, 100, 0, 99.90, 100.10, now)
	daily := []Bar{}
	for i := 25; i > 0; i-- {
		d := now.AddDate(0, 0, -i)
		daily = append(daily, Bar{T: d.Unix(), O: 100, H: 103, L: 97, C: 100, V: 2_000_000})
	}
	intraday := []Bar{{T: now.Add(-15 * time.Minute).Unix(), O: 99, H: 101, L: 99, C: 100, V: 300_000}}
	out := deriveLiquidityStatesWithContext(map[string]Quote{sym: q}, map[string]map[string][]Bar{sym: {"daily": daily, "intraday": intraday}}, map[string]LiquidityBaseline{sym: {NormalSpreadPct: .05, Samples: 30, UpdatedAt: now.UnixMilli()}}, now)[sym]
	if out.RelativeSpreadMultiple < 3.5 || out.DollarVolume <= 0 || out.AverageDollarVolume20D <= 0 || out.DailyRangeVolatilityPct <= 0 || out.VolatilityAdjustedSlippage <= 0 {
		t.Fatalf("liquidity metrics incomplete: %+v", out)
	}
	if out.OpeningLiquidity == "NOT IN OPENING WINDOW" || out.SlippageRisk == "UNKNOWN" {
		t.Fatalf("opening/slippage classification missing: %+v", out)
	}
}

func TestV168SmartNotificationsMaterialStateChangesOnly(t *testing.T) {
	now := time.Now()
	base := now.Add(-5 * time.Minute).UnixMilli()
	sym := "NVDA"
	validation := SignalValidationState{Snapshots: []SignalSnapshot{
		{Symbol: sym, Horizon: "day", Timestamp: base - 60_000, Price: 95, EntryLow: 99, EntryHigh: 101, Readiness: "READY", MarketTradeability: "TRADE NORMALLY", LiquidityState: "HEALTHY"},
		{Symbol: sym, Horizon: "day", Timestamp: base, Price: 100, EntryLow: 99, EntryHigh: 101, Readiness: "CAUTION", MarketTradeability: "WAIT", LiquidityState: "RISK"},
	}}
	sec := map[string]SECIntelligenceSummary{sym: {UpdatedAt: base, RecentInsiderTransactions: []SECInsiderTransaction{{Actor: "Director", Classification: "BUY", Code: "P", TransactionDate: now.Format("2006-01-02"), FiledAt: now.Format("2006-01-02"), Shares: 1000, Value: 100000}}}}
	fresh := []FreshnessDiagnostic{{Dataset: "Quotes", State: "STALE", Reason: "quote overdue", ReceivedAt: base}}
	router := ProviderRouterSnapshot{Routes: []ProviderRouteState{{Dataset: "VIX", Route: []ProviderRouteHop{{Provider: "yfinance", Recovery: "RECOVERED", LastFailure: base - 120_000, LastSuccess: base}}}}}
	notes := buildSmartNotifications(nil, nil, nil, EventDecisionCorrelation{}, now, SmartNotificationContext{Validation: validation, SEC: sec, Freshness: fresh, ProviderRouter: router})
	kinds := map[string]bool{}
	for _, n := range notes {
		kinds[n.Kind] = true
	}
	for _, k := range []string{"READINESS CHANGED", "ENTRY ZONE REACHED", "TRADEABILITY DETERIORATED", "LIQUIDITY DETERIORATED", "INSIDER PURCHASE", "DATA DEGRADED", "DATA RECOVERED"} {
		if !kinds[k] {
			t.Fatalf("missing smart notification %s: %+v", k, notes)
		}
	}
	for _, n := range notes {
		if math.IsNaN(float64(n.CreatedAt)) {
			t.Fatal("invalid notification time")
		}
	}
}
