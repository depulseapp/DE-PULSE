package main

import (
	"strings"
	"testing"
	"time"
)

func v167DailyBars(symbol string, count int, start, step float64, now time.Time) []Bar {
	rows := make([]Bar, 0, count)
	base := now.AddDate(0, 0, -count+1).Truncate(24 * time.Hour)
	for i := 0; i < count; i++ {
		c := start + float64(i)*step
		rows = append(rows, Bar{T: base.AddDate(0, 0, i).Unix(), O: c - step*.25, H: c + 1, L: c - 1, C: c, V: 1_000_000})
	}
	return rows
}

func v167CurrentQuote(symbol string, price, change float64, now time.Time) Quote {
	return Quote{Symbol: symbol, Price: price, ChangePercent: change, Source: "Alpaca", ProviderTimestamp: now.UnixMilli(), UpdatedAt: now.UnixMilli()}
}

func TestV167EconomicCalendarOriginalAcceptanceCategoryAndImpactTruth(t *testing.T) {
	now := time.Now()
	actual, forecast, previous := 3.1, 3.0, 3.2
	events := []MacroEvent{
		{ID: "cpi", Region: "US", Name: "Core CPI YoY", Impact: "HIGH", StartsAt: now.Add(time.Hour).UnixMilli(), Date: now.Format("2006-01-02"), TimeKnown: true, Expected: &forecast, Actual: &actual, Previous: &previous, Unit: "%", Source: "BLS"},
		{ID: "nfp", Region: "US", Name: "Nonfarm Payrolls", Impact: "MEDIUM", StartsAt: now.Add(2 * time.Hour).UnixMilli(), Date: now.Format("2006-01-02"), TimeKnown: true, Source: "BLS"},
		{ID: "fomc", Region: "US", Name: "FOMC Rate Decision", Impact: "HIGH", StartsAt: now.Add(3 * time.Hour).UnixMilli(), Date: now.Format("2006-01-02"), TimeKnown: true, Source: "Federal Reserve"},
	}
	rows := buildEconomicCalendar(events, nil, now)
	if len(rows) != 3 {
		t.Fatalf("calendar rows=%d", len(rows))
	}
	got := map[string]EconomicCalendarEntry{}
	for _, row := range rows {
		got[row.ID] = row
	}
	if got["cpi"].Category != "INFLATION" || got["nfp"].Category != "LABOR" || got["fomc"].Category != "CENTRAL BANK" {
		t.Fatalf("derived category truth lost: %+v", got)
	}
	if got["cpi"].Actual == nil || got["cpi"].Forecast == nil || got["cpi"].Previous == nil || got["cpi"].Surprise == nil {
		t.Fatalf("actual/forecast/previous/surprise incomplete: %+v", got["cpi"])
	}
}

func TestV167BreadthInternalsOriginalAcceptance(t *testing.T) {
	now := time.Now()
	quotes := map[string]Quote{}
	bars := map[string]map[string][]Bar{}
	for i, sym := range marketIntelligenceBreadthUniverse {
		up := i%3 != 0
		step := .20
		if !up {
			step = -.08
		}
		rows := v167DailyBars(sym, 220, 80, step, now)
		price := rows[len(rows)-1].C
		quotes[sym] = v167CurrentQuote(sym, price, func() float64 {
			if up {
				return .8
			}
			return -.5
		}(), now)
		bars[sym] = map[string][]Bar{"daily": rows}
	}
	b := marketBreadthStateWithBars(quotes, bars, now)
	in := b.Internals
	if in.State != "AVAILABLE" {
		t.Fatalf("internals state=%+v", in)
	}
	if in.Above20MAPct == nil || in.Above50MAPct == nil || in.Above200MAPct == nil {
		t.Fatalf("MA participation missing=%+v", in)
	}
	if in.HighLowDenominator < 10 {
		t.Fatalf("high/low coverage weak=%+v", in)
	}
	if in.SectorExpected != 11 || in.SectorParticipationPct == nil {
		t.Fatalf("sector participation missing=%+v", in)
	}
	if !strings.Contains(strings.ToLower(in.Detail), "not exchange-wide") {
		t.Fatalf("coverage limitation not truthful=%q", in.Detail)
	}
}

func TestV167RelativeStrengthIncludesDayAndSwingAgainstCanonicalBenchmarks(t *testing.T) {
	now := time.Now()
	st := defaultState()
	st.UI.SelectedTicker = "NVDA"
	bars := map[string]map[string][]Bar{}
	for _, sym := range []string{"NVDA", "SPY", "QQQ", "XLK", "SMH"} {
		step := .15
		if sym == "NVDA" {
			step = .45
		}
		bars[sym] = map[string][]Bar{
			"daily":    v167DailyBars(sym, 80, 100, step, now),
			"intraday": v161BarsEnding(sym, 40, 100, .2, now, 5*time.Minute),
		}
	}
	snap := buildMarketIntelligenceSnapshotWithContext(st, nil, bars, nil, GlobalMarketContext{}, MarketTradeabilityContext{}, now)
	seen := map[string]bool{}
	for _, row := range snap.RelativeStrength {
		if row.Symbol == "NVDA" {
			seen[strings.ToUpper(row.Horizon)+":"+row.Benchmark] = true
		}
	}
	for _, key := range []string{"DAY:SPY", "DAY:QQQ", "DAY:XLK", "DAY:SMH", "SWING:SPY", "SWING:QQQ", "SWING:XLK", "SWING:SMH"} {
		if !seen[key] {
			t.Fatalf("missing canonical RS %s; seen=%v", key, seen)
		}
	}
}

func TestV167SectorIndustryRegimeCombinesMomentumRelativeStrengthAndMemberBreadth(t *testing.T) {
	now := time.Now()
	bars := map[string]map[string][]Bar{}
	for _, sym := range []string{"SPY", "XLK", "NVDA", "AMD", "AVGO", "INTC", "MU", "QCOM", "ARM", "TSM", "ASML", "AMAT", "LRCX"} {
		step := .15
		if sym == "XLK" || sym == "NVDA" || sym == "AMD" || sym == "AVGO" {
			step = .35
		}
		bars[sym] = map[string][]Bar{"daily": v167DailyBars(sym, 80, 100, step, now)}
	}
	r := benchmarkRegime("Sector", "Technology", "XLK", bars, now)
	if r.State == "UNAVAILABLE" || r.BreadthPct == nil || r.MemberCount == 0 || r.CoveragePct <= 0 {
		t.Fatalf("sector regime does not reconcile member breadth=%+v", r)
	}
	if !strings.Contains(strings.ToLower(r.Detail), "relative") || !strings.Contains(strings.ToLower(r.Detail), "participation") {
		t.Fatalf("regime detail not explainable=%q", r.Detail)
	}
}

func TestV167TradeabilityConsumesOriginalContextWithoutChangingDeskFormula(t *testing.T) {
	now := time.Now()
	quotes := map[string]Quote{
		"SPY": v167CurrentQuote("SPY", 500, -.4, now), "QQQ": v167CurrentQuote("QQQ", 450, -.7, now), "VIX": v167CurrentQuote("VIX", 24, 7, now),
	}
	p := 32.0
	breadth := MarketBreadthState{State: "AVAILABLE", ParticipationPct: &p, Internals: MarketBreadthInternalsState{State: "AVAILABLE", Above50MAPct: &p, SectorParticipationPct: &p}}
	ctx := MarketTradeabilityContext{
		EventMode: EventModeState{Active: true, Name: "CPI", Phase: "PRE"},
		Freshness: []FreshnessDiagnostic{{Dataset: "news", State: "STALE"}},
		Liquidity: map[string]LiquidityState{"SPY": {Symbol: "SPY", State: "RISK", SpreadPct: .35}, "QQQ": {Symbol: "QQQ", State: "RISK", SpreadPct: .30}},
		Scanner:   ScannerState{Mode: "day", Status: "complete", UpdatedAt: now.UnixMilli(), Results: nil},
		Options:   map[string]OptionsContext{"SPY": {Symbol: "SPY", State: "AVAILABLE", Bias: "BEARISH", AverageIV: .55}, "QQQ": {Symbol: "QQQ", State: "AVAILABLE", Bias: "BEARISH", AverageIV: .52}},
		Structure: map[string]MarketStructureState{"day": {State: "EXTENDED"}},
	}
	out := marketTradeabilityWithContext(quotes, breadth, GlobalMarketContext{Tone: "RISK-OFF"}, ctx, now)
	if out.State == "TRADE NORMALLY" || out.Score >= 70 {
		t.Fatalf("risk context not reconciled=%+v", out)
	}
	for _, key := range []string{"volatility", "breadth", "eventRisk", "liquidity", "freshness", "setups", "options", "global"} {
		if _, ok := out.Components[key]; !ok {
			t.Fatalf("component %s missing: %+v", key, out.Components)
		}
	}
	if !strings.Contains(out.Detail, "never mutates deterministic desk scores") {
		t.Fatalf("deterministic protection not explicit=%q", out.Detail)
	}
}

func TestV167TradeabilityFeedsPreparationSurfacesWithoutPersistingDerivedState(t *testing.T) {
	now := time.Now()
	jobs := initialPreparationJobs(now)
	tradeability := MarketTradeabilityState{State: "WAIT", Detail: "Macro event and liquidity risk require patience.", UpdatedAt: now.UnixMilli()}
	out := preparationsWithMarketTradeability(clone(jobs), tradeability, now)
	for _, key := range []string{"pre-market-prep", "market-open-prep"} {
		p := out[key]
		if !strings.Contains(strings.Join(p.Summary, " | "), "Market Tradeability · WAIT") {
			t.Fatalf("%s missing Tradeability summary: %+v", key, p.Summary)
		}
		found := false
		for _, x := range p.Exceptions {
			if x.Reason == "Market Tradeability WAIT" && x.Source == "Market Intelligence" {
				found = true
			}
		}
		if !found {
			t.Fatalf("%s missing Tradeability review exception: %+v", key, p.Exceptions)
		}
	}
	if len(jobs["pre-market-prep"].Summary) != 0 {
		t.Fatalf("source preparation state was mutated: %+v", jobs["pre-market-prep"].Summary)
	}
}

func TestV167DayDiscoveryRankingUsesSharedSessionRelativeStrength(t *testing.T) {
	rows := []ScannerResult{
		{Symbol: "SPY", Mode: "day", ChangePercent: .5, Score: 60},
		{Symbol: "QQQ", Mode: "day", ChangePercent: .7, Score: 60},
		{Symbol: "NVDA", Mode: "day", ChangePercent: 3.0, Score: 60},
		{Symbol: "TSLA", Mode: "day", ChangePercent: -1.0, Score: 60},
	}
	got := applyScannerSessionRelativeStrength(rows, "day")
	by := map[string]ScannerResult{}
	for _, row := range got {
		by[row.Symbol] = row
	}
	if by["NVDA"].RelativeStrength <= 2 || by["NVDA"].Score <= 60 || by["NVDA"].RSBenchmark != "SPY/QQQ SESSION" {
		t.Fatalf("Day discovery RS not feeding rank: %+v", by["NVDA"])
	}
	if by["TSLA"].RelativeStrength >= 0 || by["TSLA"].Score >= 60 {
		t.Fatalf("lagging Day discovery RS not reflected: %+v", by["TSLA"])
	}
}
