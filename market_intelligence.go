package main

import (
	"fmt"
	"math"
	"sort"
	"strings"
	"time"
)

// Market Intelligence is a context-only professional decision layer. It reuses
// canonical quotes, bars, liquidity, Global/Macro and Provider Router state and
// never owns a provider fetch or mutates deterministic Day/Swing/Long scores.

type MarketStructureState struct {
	Horizon     string  `json:"horizon"`
	State       string  `json:"state"`
	Trend       string  `json:"trend,omitempty"`
	Price       float64 `json:"price,omitempty"`
	MomentumPct float64 `json:"momentumPct,omitempty"`
	DrawdownPct float64 `json:"drawdownPct,omitempty"`
	Coverage    int     `json:"coverage"`
	UpdatedAt   int64   `json:"updatedAt,omitempty"`
	Detail      string  `json:"detail,omitempty"`
}

type MarketBreadthInternalsState struct {
	State                  string   `json:"state"`
	Above20MAPct           *float64 `json:"above20MaPct,omitempty"`
	Above50MAPct           *float64 `json:"above50MaPct,omitempty"`
	Above200MAPct          *float64 `json:"above200MaPct,omitempty"`
	Above20Denominator     int      `json:"above20Denominator"`
	Above50Denominator     int      `json:"above50Denominator"`
	Above200Denominator    int      `json:"above200Denominator"`
	NewHighs20             int      `json:"newHighs20"`
	NewLows20              int      `json:"newLows20"`
	HighLowDenominator     int      `json:"highLowDenominator"`
	SectorAdvancers        int      `json:"sectorAdvancers"`
	SectorDecliners        int      `json:"sectorDecliners"`
	SectorUnchanged        int      `json:"sectorUnchanged"`
	SectorExpected         int      `json:"sectorExpected"`
	SectorParticipationPct *float64 `json:"sectorParticipationPct,omitempty"`
	UpdatedAt              int64    `json:"updatedAt,omitempty"`
	Detail                 string   `json:"detail,omitempty"`
}

type MarketBreadthState struct {
	Label            string                      `json:"label"`
	State            string                      `json:"state"`
	Expected         int                         `json:"expected"`
	Denominator      int                         `json:"denominator"`
	Fresh            int                         `json:"fresh"`
	Stale            int                         `json:"stale"`
	Advancers        int                         `json:"advancers"`
	Decliners        int                         `json:"decliners"`
	Unchanged        int                         `json:"unchanged"`
	CoveragePct      float64                     `json:"coveragePct"`
	ParticipationPct *float64                    `json:"participationPct,omitempty"`
	Internals        MarketBreadthInternalsState `json:"internals"`
	UpdatedAt        int64                       `json:"updatedAt,omitempty"`
	Detail           string                      `json:"detail,omitempty"`
}

type MarketTradeabilityState struct {
	State      string         `json:"state"`
	Score      int            `json:"score"`
	Drivers    []string       `json:"drivers,omitempty"`
	Blockers   []string       `json:"blockers,omitempty"`
	Components map[string]int `json:"components,omitempty"`
	UpdatedAt  int64          `json:"updatedAt,omitempty"`
	Detail     string         `json:"detail,omitempty"`
}

type MarketTradeabilityContext struct {
	EventMode  EventModeState
	Freshness  []FreshnessDiagnostic
	Liquidity  map[string]LiquidityState
	Scanner    ScannerState
	Options    map[string]OptionsContext
	Structure  map[string]MarketStructureState
	DaySymbols []string
}

type SymbolClassification struct {
	Symbol            string `json:"symbol"`
	Sector            string `json:"sector,omitempty"`
	Industry          string `json:"industry,omitempty"`
	SectorBenchmark   string `json:"sectorBenchmark,omitempty"`
	IndustryBenchmark string `json:"industryBenchmark,omitempty"`
}

type RelativeStrengthState struct {
	Symbol             string  `json:"symbol"`
	Horizon            string  `json:"horizon"`
	Benchmark          string  `json:"benchmark"`
	State              string  `json:"state"`
	SymbolReturnPct    float64 `json:"symbolReturnPct,omitempty"`
	BenchmarkReturnPct float64 `json:"benchmarkReturnPct,omitempty"`
	RelativePct        float64 `json:"relativePct,omitempty"`
	UpdatedAt          int64   `json:"updatedAt,omitempty"`
	Detail             string  `json:"detail,omitempty"`
}

type MarketRegimeState struct {
	Level          string   `json:"level"`
	Name           string   `json:"name,omitempty"`
	Benchmark      string   `json:"benchmark,omitempty"`
	State          string   `json:"state"`
	MomentumPct    float64  `json:"momentumPct,omitempty"`
	RelativePct    float64  `json:"relativePct,omitempty"`
	BreadthPct     *float64 `json:"breadthPct,omitempty"`
	MemberCount    int      `json:"memberCount,omitempty"`
	MemberExpected int      `json:"memberExpected,omitempty"`
	CoveragePct    float64  `json:"coveragePct,omitempty"`
	UpdatedAt      int64    `json:"updatedAt,omitempty"`
	Detail         string   `json:"detail,omitempty"`
}

type LiquiditySlippageState struct {
	Symbol    string  `json:"symbol"`
	State     string  `json:"state"`
	SpreadPct float64 `json:"spreadPct,omitempty"`
	UpdatedAt int64   `json:"updatedAt,omitempty"`
	Detail    string  `json:"detail,omitempty"`
}

type MarketIntelligenceSnapshot struct {
	Tradeability     MarketTradeabilityState         `json:"tradeability"`
	Structure        map[string]MarketStructureState `json:"structure"`
	Breadth          MarketBreadthState              `json:"breadth"`
	SelectedSymbol   string                          `json:"selectedSymbol,omitempty"`
	Classification   SymbolClassification            `json:"classification"`
	RelativeStrength []RelativeStrengthState         `json:"relativeStrength"`
	SectorRegime     MarketRegimeState               `json:"sectorRegime"`
	IndustryRegime   MarketRegimeState               `json:"industryRegime"`
	Liquidity        LiquiditySlippageState          `json:"liquidity"`
	UpdatedAt        int64                           `json:"updatedAt,omitempty"`
}

var marketIntelligenceBreadthUniverse = []string{
	"SPY", "QQQ", "DIA", "IWM", "XLK", "XLC", "XLY", "XLP", "XLE", "XLF", "XLV", "XLI", "XLB", "XLRE", "XLU",
}

// canonicalSymbolClassifications is the single ticker -> industry -> sector mapping
// owner for Market Intelligence. Unknown symbols stay explicitly unmapped.
var canonicalSymbolClassifications = map[string]SymbolClassification{
	"AAPL":  {Symbol: "AAPL", Sector: "Technology", Industry: "Technology Hardware", SectorBenchmark: "XLK"},
	"MSFT":  {Symbol: "MSFT", Sector: "Technology", Industry: "Software", SectorBenchmark: "XLK"},
	"ORCL":  {Symbol: "ORCL", Sector: "Technology", Industry: "Software", SectorBenchmark: "XLK"},
	"PLTR":  {Symbol: "PLTR", Sector: "Technology", Industry: "Software", SectorBenchmark: "XLK"},
	"NVDA":  {Symbol: "NVDA", Sector: "Technology", Industry: "Semiconductors", SectorBenchmark: "XLK", IndustryBenchmark: "SMH"},
	"AMD":   {Symbol: "AMD", Sector: "Technology", Industry: "Semiconductors", SectorBenchmark: "XLK", IndustryBenchmark: "SMH"},
	"AVGO":  {Symbol: "AVGO", Sector: "Technology", Industry: "Semiconductors", SectorBenchmark: "XLK", IndustryBenchmark: "SMH"},
	"INTC":  {Symbol: "INTC", Sector: "Technology", Industry: "Semiconductors", SectorBenchmark: "XLK", IndustryBenchmark: "SMH"},
	"MU":    {Symbol: "MU", Sector: "Technology", Industry: "Semiconductors", SectorBenchmark: "XLK", IndustryBenchmark: "SMH"},
	"QCOM":  {Symbol: "QCOM", Sector: "Technology", Industry: "Semiconductors", SectorBenchmark: "XLK", IndustryBenchmark: "SMH"},
	"ARM":   {Symbol: "ARM", Sector: "Technology", Industry: "Semiconductors", SectorBenchmark: "XLK", IndustryBenchmark: "SMH"},
	"TSM":   {Symbol: "TSM", Sector: "Technology", Industry: "Semiconductors", SectorBenchmark: "XLK", IndustryBenchmark: "SMH"},
	"ASML":  {Symbol: "ASML", Sector: "Technology", Industry: "Semiconductors", SectorBenchmark: "XLK", IndustryBenchmark: "SMH"},
	"AMAT":  {Symbol: "AMAT", Sector: "Technology", Industry: "Semiconductor Equipment", SectorBenchmark: "XLK", IndustryBenchmark: "SMH"},
	"LRCX":  {Symbol: "LRCX", Sector: "Technology", Industry: "Semiconductor Equipment", SectorBenchmark: "XLK", IndustryBenchmark: "SMH"},
	"META":  {Symbol: "META", Sector: "Communication Services", Industry: "Interactive Media", SectorBenchmark: "XLC"},
	"GOOGL": {Symbol: "GOOGL", Sector: "Communication Services", Industry: "Interactive Media", SectorBenchmark: "XLC"},
	"GOOG":  {Symbol: "GOOG", Sector: "Communication Services", Industry: "Interactive Media", SectorBenchmark: "XLC"},
	"NFLX":  {Symbol: "NFLX", Sector: "Communication Services", Industry: "Entertainment", SectorBenchmark: "XLC"},
	"AMZN":  {Symbol: "AMZN", Sector: "Consumer Discretionary", Industry: "Broadline Retail", SectorBenchmark: "XLY"},
	"TSLA":  {Symbol: "TSLA", Sector: "Consumer Discretionary", Industry: "Automobiles", SectorBenchmark: "XLY"},
	"JPM":   {Symbol: "JPM", Sector: "Financials", Industry: "Banks", SectorBenchmark: "XLF"},
	"BAC":   {Symbol: "BAC", Sector: "Financials", Industry: "Banks", SectorBenchmark: "XLF"},
	"SOFI":  {Symbol: "SOFI", Sector: "Financials", Industry: "Financial Services", SectorBenchmark: "XLF"},
	"XOM":   {Symbol: "XOM", Sector: "Energy", Industry: "Integrated Oil & Gas", SectorBenchmark: "XLE"},
	"CVX":   {Symbol: "CVX", Sector: "Energy", Industry: "Integrated Oil & Gas", SectorBenchmark: "XLE"},
	"USO":   {Symbol: "USO", Sector: "Energy", Industry: "Energy Commodity ETF", SectorBenchmark: "XLE"},
	"LLY":   {Symbol: "LLY", Sector: "Health Care", Industry: "Pharmaceuticals", SectorBenchmark: "XLV"},
	"UNH":   {Symbol: "UNH", Sector: "Health Care", Industry: "Managed Care", SectorBenchmark: "XLV"},
	"CAT":   {Symbol: "CAT", Sector: "Industrials", Industry: "Machinery", SectorBenchmark: "XLI"},
	"WMT":   {Symbol: "WMT", Sector: "Consumer Staples", Industry: "Consumer Staples Distribution", SectorBenchmark: "XLP"},
	"COST":  {Symbol: "COST", Sector: "Consumer Staples", Industry: "Consumer Staples Distribution", SectorBenchmark: "XLP"},
	"GLD":   {Symbol: "GLD", Sector: "Materials / Alternative", Industry: "Gold ETF", SectorBenchmark: "XLB"},
	"SLV":   {Symbol: "SLV", Sector: "Materials / Alternative", Industry: "Silver ETF", SectorBenchmark: "XLB"},
}

func classificationForSymbol(symbol string) SymbolClassification {
	symbol = normalizeSymbol(symbol)
	if c, ok := canonicalSymbolClassifications[symbol]; ok {
		return c
	}
	return SymbolClassification{Symbol: symbol}
}

// marketIntelligenceBarEvidence is the shared historical-evidence truth boundary.
// WHY: calculations may run at any time, but they may only publish an active
// market context when the exact requested timeframe has a plausible, current bar.
func marketIntelligenceBarEvidence(rows []Bar, key string, now time.Time) (int64, bool, string) {
	latest := int64(0)
	for _, b := range rows {
		if b.T > latest {
			latest = b.T
		}
	}
	if latest <= 0 {
		return 0, false, "Historical-bar timestamp unavailable."
	}
	stamp := time.Unix(latest, 0)
	if stamp.After(now.Add(5 * time.Minute)) {
		return latest * 1000, false, "Historical-bar timestamp is materially future-skewed."
	}
	age := now.Sub(stamp)
	limit := 7 * 24 * time.Hour
	switch key {
	case "intraday":
		limit = 45 * time.Minute
		switch marketSessionET(now) {
		case "overnight":
			limit = 24 * time.Hour
		case "closed", "weekend":
			limit = 96 * time.Hour
		}
	case "weekly":
		limit = 14 * 24 * time.Hour
	}
	if age > limit {
		return latest * 1000, false, fmt.Sprintf("Latest canonical %s bar is stale (%s old; limit %s).", key, age.Round(time.Minute), limit)
	}
	return latest * 1000, true, ""
}

func miBarReturn(rows []Bar, periods int) (float64, bool) {
	if periods <= 0 || len(rows) < periods+1 {
		return 0, false
	}
	a, b := rows[len(rows)-1-periods].C, rows[len(rows)-1].C
	if a <= 0 || b <= 0 {
		return 0, false
	}
	return (b/a - 1) * 100, true
}
func miSMA(rows []Bar, n int) (float64, bool) {
	if n <= 0 || len(rows) < n {
		return 0, false
	}
	sum := 0.0
	for _, b := range rows[len(rows)-n:] {
		if b.C <= 0 {
			return 0, false
		}
		sum += b.C
	}
	return sum / float64(n), true
}

func marketStructureFor(bars map[string]map[string][]Bar, horizon string, now time.Time) MarketStructureState {
	h := strings.ToLower(strings.TrimSpace(horizon))
	state := MarketStructureState{Horizon: horizon, State: "UNAVAILABLE", Trend: "UNKNOWN"}
	key, lookback := "daily", 20
	switch h {
	case "day":
		key, lookback = "intraday", 12
	case "long", "long-term":
		key, lookback = "weekly", 26
	}
	rows := bars["SPY"][key]
	state.Coverage = len(rows)
	if len(rows) < lookback+1 {
		state.Detail = fmt.Sprintf("Requires at least %d canonical %s bars; %d available.", lookback+1, key, len(rows))
		return state
	}
	stamp, ok, why := marketIntelligenceBarEvidence(rows, key, now)
	state.UpdatedAt = stamp
	if !ok {
		state.Detail = why
		return state
	}
	price := rows[len(rows)-1].C
	if price <= 0 {
		state.Detail = "Latest canonical bar close is unavailable."
		return state
	}
	mom, _ := miBarReturn(rows, lookback)
	smaN := lookback
	if smaN > 20 {
		smaN = 20
	}
	sma, _ := miSMA(rows, smaN)
	prior := rows[:len(rows)-1]
	if len(prior) > lookback {
		prior = prior[len(prior)-lookback:]
	}
	priorHigh, allHigh := 0.0, 0.0
	for _, b := range rows {
		if b.H > allHigh {
			allHigh = b.H
		}
	}
	for _, b := range prior {
		if b.H > priorHigh {
			priorHigh = b.H
		}
	}
	dd := 0.0
	if allHigh > 0 {
		dd = (price/allHigh - 1) * 100
	}
	dist := 0.0
	if sma > 0 {
		dist = (price/sma - 1) * 100
	}
	st, tr := "RANGE", "SIDEWAYS"
	switch {
	case priorHigh > 0 && price > priorHigh*1.002 && mom > 1.2:
		st, tr = "BREAKOUT", "UP"
	case dist > 6 && mom > 6:
		st, tr = "EXTENDED", "UP"
	case dd <= -10 && mom < 0:
		st, tr = "CORRECTION", "DOWN / CORRECTING"
	case price < sma && mom < -4:
		st, tr = "DOWNTREND", "DOWN"
	case price >= sma && mom > 1 && dd < -4:
		st, tr = "REVERSAL ATTEMPT", "IMPROVING"
	case price >= sma && mom > 0:
		st, tr = "PULLBACK", "UP"
	case price < sma && dd > -12:
		st, tr = "DEEP PULLBACK", "TESTING"
	case math.Abs(mom) < 2 && math.Abs(dist) < 2:
		st, tr = "BASE", "SIDEWAYS"
	}
	state.State, state.Trend, state.Price, state.MomentumPct, state.DrawdownPct = st, tr, price, mom, dd
	state.Detail = fmt.Sprintf("Canonical %s evidence · %d-period momentum %.1f%% · drawdown %.1f%% · %.1f%% vs %d-bar mean.", key, lookback, mom, dd, dist, smaN)
	return state
}

func quoteIsCurrentForMarketIntelligence(q Quote, now time.Time) (int64, bool, string) {
	if q.Price <= 0 {
		return 0, false, "price unavailable"
	}
	providerAge, receiptAge, valid, detail := quoteEvidenceTimestampTruth(q, now.UnixMilli())
	stamp := q.ProviderTimestamp
	if stamp == 0 {
		stamp = q.UpdatedAt
	}
	if !valid {
		return stamp, false, detail
	}
	limit := reconciliationAgeLimit(now.UnixMilli())
	if providerAge > limit || receiptAge > limit {
		return stamp, false, fmt.Sprintf("Quote evidence is stale (market age %s; receipt age %s; limit %s).", time.Duration(providerAge)*time.Millisecond, time.Duration(receiptAge)*time.Millisecond, time.Duration(limit)*time.Millisecond)
	}
	return stamp, true, ""
}

var marketIntelligenceSectorUniverse = []string{"XLK", "XLC", "XLY", "XLP", "XLE", "XLF", "XLV", "XLI", "XLB", "XLRE", "XLU"}

func pctPointer(n, d int) *float64 {
	if d <= 0 {
		return nil
	}
	v := float64(n) / float64(d) * 100
	return &v
}

func marketBreadthInternals(quotes map[string]Quote, bars map[string]map[string][]Bar, now time.Time) MarketBreadthInternalsState {
	out := MarketBreadthInternalsState{State: "UNAVAILABLE", SectorExpected: len(marketIntelligenceSectorUniverse)}
	above20, above50, above200 := 0, 0, 0
	latest := int64(0)
	for _, sym := range marketIntelligenceBreadthUniverse {
		q, ok := quotes[sym]
		if !ok || q.Price <= 0 {
			continue
		}
		qStamp, current, _ := quoteIsCurrentForMarketIntelligence(q, now)
		if !current {
			continue
		}
		rows := bars[sym]["daily"]
		if len(rows) == 0 {
			continue
		}
		bStamp, currentBars, _ := marketIntelligenceBarEvidence(rows, "daily", now)
		if !currentBars {
			continue
		}
		if qStamp > latest {
			latest = qStamp
		}
		if bStamp > latest {
			latest = bStamp
		}
		if len(rows) >= 20 {
			out.Above20Denominator++
			if ma, ok := miSMA(rows, 20); ok && q.Price > ma {
				above20++
			}
		}
		if len(rows) >= 50 {
			out.Above50Denominator++
			if ma, ok := miSMA(rows, 50); ok && q.Price > ma {
				above50++
			}
		}
		if len(rows) >= 200 {
			out.Above200Denominator++
			if ma, ok := miSMA(rows, 200); ok && q.Price > ma {
				above200++
			}
		}
		if len(rows) >= 20 {
			prior := rows[len(rows)-20:]
			high, low := 0.0, math.MaxFloat64
			for _, row := range prior {
				if row.H > high {
					high = row.H
				}
				if row.L > 0 && row.L < low {
					low = row.L
				}
			}
			if high > 0 && low < math.MaxFloat64 {
				out.HighLowDenominator++
				if q.Price >= high {
					out.NewHighs20++
				}
				if q.Price <= low {
					out.NewLows20++
				}
			}
		}
	}
	out.Above20MAPct = pctPointer(above20, out.Above20Denominator)
	out.Above50MAPct = pctPointer(above50, out.Above50Denominator)
	out.Above200MAPct = pctPointer(above200, out.Above200Denominator)
	for _, sym := range marketIntelligenceSectorUniverse {
		q := quotes[sym]
		stamp, current, _ := quoteIsCurrentForMarketIntelligence(q, now)
		if !current {
			continue
		}
		if stamp > latest {
			latest = stamp
		}
		switch {
		case q.ChangePercent > .05:
			out.SectorAdvancers++
		case q.ChangePercent < -.05:
			out.SectorDecliners++
		default:
			out.SectorUnchanged++
		}
	}
	sectorDenom := out.SectorAdvancers + out.SectorDecliners + out.SectorUnchanged
	out.SectorParticipationPct = pctPointer(out.SectorAdvancers, sectorDenom)
	out.UpdatedAt = latest
	usableInternals := out.Above20Denominator >= 6 || out.Above50Denominator >= 6 || out.HighLowDenominator >= 6
	sectorCoverage := 0.0
	if out.SectorExpected > 0 {
		sectorCoverage = float64(sectorDenom) / float64(out.SectorExpected) * 100
	}
	switch {
	case !usableInternals && sectorDenom == 0:
		out.State = "UNAVAILABLE"
	case usableInternals && sectorCoverage >= 60:
		out.State = "AVAILABLE"
	default:
		out.State = "DEGRADED"
	}
	out.Detail = fmt.Sprintf("Tracked-universe internals only: MA participation uses current canonical quotes against completed daily bars; 20-session highs/lows use %d eligible symbols; sector participation uses %d/%d current SPDR sector benchmarks. Not exchange-wide NYSE/Nasdaq breadth.", out.HighLowDenominator, sectorDenom, out.SectorExpected)
	return out
}

func marketBreadthStateWithBars(quotes map[string]Quote, bars map[string]map[string][]Bar, now time.Time) MarketBreadthState {
	b := MarketBreadthState{Label: "Tracked Broad-Market Breadth", State: "UNAVAILABLE", Expected: len(marketIntelligenceBreadthUniverse)}
	latest := int64(0)
	for _, sym := range marketIntelligenceBreadthUniverse {
		q, exists := quotes[sym]
		if !exists || q.Price <= 0 {
			b.Stale++
			continue
		}
		stamp, current, _ := quoteIsCurrentForMarketIntelligence(q, now)
		if !current {
			b.Stale++
			continue
		}
		b.Denominator++
		b.Fresh++
		if stamp > latest {
			latest = stamp
		}
		switch {
		case q.ChangePercent > .05:
			b.Advancers++
		case q.ChangePercent < -.05:
			b.Decliners++
		default:
			b.Unchanged++
		}
	}
	b.UpdatedAt = latest
	if b.Expected > 0 {
		b.CoveragePct = float64(b.Fresh) / float64(b.Expected) * 100
	}
	b.Internals = marketBreadthInternals(quotes, bars, now)
	if b.Fresh == 0 {
		b.Detail = "No current tracked-universe observations; breadth participation unavailable."
		return b
	}
	p := float64(b.Advancers) / float64(b.Fresh) * 100
	if b.CoveragePct < 60 {
		b.State = "DEGRADED"
		b.Detail = fmt.Sprintf("%d/%d current tracked broad/sector ETF observations (%.0f%% coverage); directional participation withheld.", b.Fresh, b.Expected, b.CoveragePct)
		return b
	}
	b.ParticipationPct = &p
	switch {
	case p >= 60:
		b.State = "BROAD"
	case p <= 40:
		b.State = "WEAK"
	default:
		b.State = "MIXED"
	}
	b.Detail = fmt.Sprintf("%d advancing / %d declining / %d unchanged across %d current observations; %.0f%% of the explicit %d-symbol tracked universe is current. Internals add MA participation, 20-session highs/lows and sector participation where canonical history is sufficient.", b.Advancers, b.Decliners, b.Unchanged, b.Fresh, b.CoveragePct, b.Expected)
	return b
}

func relativeStrengthFor(bars map[string]map[string][]Bar, symbol, benchmark, horizon string, now time.Time) RelativeStrengthState {
	symbol, benchmark = normalizeSymbol(symbol), normalizeSymbol(benchmark)
	out := RelativeStrengthState{Symbol: symbol, Benchmark: benchmark, Horizon: horizon, State: "UNAVAILABLE"}
	if symbol == "" || benchmark == "" {
		out.Detail = "Ticker or benchmark mapping unavailable."
		return out
	}
	key, periods := "daily", 20
	switch strings.ToLower(horizon) {
	case "day":
		key, periods = "intraday", 12
	case "long", "long-term":
		key, periods = "weekly", 26
	}
	sr, br := bars[symbol][key], bars[benchmark][key]
	if len(sr) < periods+1 || len(br) < periods+1 {
		out.Detail = fmt.Sprintf("Requires aligned canonical %s history for %s and %s at the exact requested horizon.", key, symbol, benchmark)
		return out
	}
	sStamp, sOK, sWhy := marketIntelligenceBarEvidence(sr, key, now)
	bStamp, bOK, bWhy := marketIntelligenceBarEvidence(br, key, now)
	out.UpdatedAt = sStamp
	if bStamp > 0 && (out.UpdatedAt == 0 || bStamp < out.UpdatedAt) {
		out.UpdatedAt = bStamp
	}
	if !sOK {
		out.Detail = sWhy
		return out
	}
	if !bOK {
		out.Detail = bWhy
		return out
	}
	sRet, sok := miBarReturn(sr, periods)
	bRet, bok := miBarReturn(br, periods)
	if !sok || !bok {
		out.Detail = "Aligned return evidence unavailable."
		return out
	}
	rel := sRet - bRet
	st := "IN LINE"
	if rel > 2 {
		st = "LEADING"
	} else if rel < -2 {
		st = "LAGGING"
	}
	out.State, out.SymbolReturnPct, out.BenchmarkReturnPct, out.RelativePct = st, sRet, bRet, rel
	out.Detail = fmt.Sprintf("%s %d-period relative performance · %s %.1f%% vs %s %.1f%% = %+.1f%%.", key, periods, symbol, sRet, benchmark, bRet, rel)
	return out
}

func regimeMembers(level, name string) []string {
	level = strings.ToLower(strings.TrimSpace(level))
	out := []string{}
	for sym, c := range canonicalSymbolClassifications {
		match := level == "sector" && c.Sector == name || level == "industry" && c.Industry == name
		if match {
			out = append(out, sym)
		}
	}
	sort.Strings(out)
	return uniqueSymbols(out)
}

func benchmarkRegime(level, name, benchmark string, bars map[string]map[string][]Bar, now time.Time) MarketRegimeState {
	out := MarketRegimeState{Level: level, Name: name, Benchmark: benchmark, State: "UNAVAILABLE"}
	if benchmark == "" {
		out.Detail = fmt.Sprintf("Canonical %s benchmark mapping unavailable.", strings.ToLower(level))
		return out
	}
	rows := bars[benchmark]["daily"]
	if len(rows) < 21 {
		out.Detail = fmt.Sprintf("Requires at least 21 canonical daily bars for %s.", benchmark)
		return out
	}
	stamp, ok, why := marketIntelligenceBarEvidence(rows, "daily", now)
	out.UpdatedAt = stamp
	if !ok {
		out.Detail = why
		return out
	}
	mom, _ := miBarReturn(rows, 20)
	out.MomentumPct = mom
	if benchmark != "SPY" {
		if rel := relativeStrengthFor(bars, benchmark, "SPY", "swing", now); rel.State != "UNAVAILABLE" {
			out.RelativePct = rel.RelativePct
		}
	}
	members := regimeMembers(level, name)
	out.MemberExpected = len(members)
	above, usable := 0, 0
	for _, sym := range members {
		memberRows := bars[sym]["daily"]
		if len(memberRows) < 20 {
			continue
		}
		memberStamp, current, _ := marketIntelligenceBarEvidence(memberRows, "daily", now)
		if !current {
			continue
		}
		last := memberRows[len(memberRows)-1].C
		ma, ok := miSMA(memberRows, 20)
		if !ok || last <= 0 || ma <= 0 {
			continue
		}
		usable++
		if last > ma {
			above++
		}
		if memberStamp > 0 && (out.UpdatedAt == 0 || memberStamp < out.UpdatedAt) {
			out.UpdatedAt = memberStamp
		}
	}
	out.MemberCount = usable
	if out.MemberExpected > 0 {
		out.CoveragePct = float64(usable) / float64(out.MemberExpected) * 100
	}
	if usable >= 2 {
		out.BreadthPct = pctPointer(above, usable)
	}
	state := "MIXED"
	breadth := 50.0
	breadthKnown := out.BreadthPct != nil
	if breadthKnown {
		breadth = *out.BreadthPct
	}
	switch {
	case mom > 2 && out.RelativePct >= 0 && (!breadthKnown || breadth >= 60):
		state = "LEADING"
	case mom < -2 && out.RelativePct <= 0 && (!breadthKnown || breadth <= 40):
		state = "LAGGING"
	case mom > 4 && breadthKnown && breadth < 40:
		state = "DIVERGING"
	case mom < -4 && breadthKnown && breadth > 60:
		state = "DIVERGING"
	}
	out.State = state
	breadthText := "member breadth unavailable"
	if out.BreadthPct != nil {
		breadthText = fmt.Sprintf("%.0f%% of %d/%d mapped members above 20D MA", *out.BreadthPct, usable, out.MemberExpected)
	}
	out.Detail = fmt.Sprintf("%s 20-session momentum %+.1f%% · relative vs SPY %+.1f%% · %s. Regime reconciles benchmark trend, relative strength and mapped-member participation without creating a second sector engine.", benchmark, mom, out.RelativePct, breadthText)
	return out
}

func liquiditySlippageState(symbol string, raw map[string]LiquidityState) LiquiditySlippageState {
	symbol = normalizeSymbol(symbol)
	out := LiquiditySlippageState{Symbol: symbol, State: "UNKNOWN"}
	x, ok := raw[symbol]
	if !ok {
		out.Detail = "No current liquidity evidence."
		return out
	}
	out.SpreadPct, out.UpdatedAt, out.Detail = x.SpreadPct, x.UpdatedAt, x.Detail
	switch x.State {
	case "HEALTHY":
		if x.SpreadPct > 0 && x.SpreadPct < .05 {
			out.State = "LOW RISK"
		} else {
			out.State = "NORMAL"
		}
	case "CAUTION":
		out.State = "ELEVATED"
	case "RISK":
		out.State = "HIGH"
	default:
		out.State = "UNKNOWN"
	}
	return out
}

// marketTradeability intentionally refuses a normal/selective state unless SPY,
// QQQ and VIX are all current. WHY: a partial core market picture can make
// otherwise reasonable context look safer than it is; degradation is explicit.
func freshnessRiskForTradeability(rows []FreshnessDiagnostic) (int, []string) {
	penalty := 0
	drivers := []string{}
	critical := map[string]bool{"vix": true, "global": true, "macro": true, "intraday bars": true}
	for _, row := range rows {
		name := strings.ToLower(strings.TrimSpace(row.Dataset))
		if !critical[name] {
			continue
		}
		state := strings.ToUpper(strings.TrimSpace(row.State))
		if state == "STALE" || state == "ERROR" || state == "UNAVAILABLE" {
			penalty += 6
			drivers = append(drivers, fmt.Sprintf("%s %s", row.Dataset, strings.ToLower(state)))
		}
	}
	if penalty > 18 {
		penalty = 18
	}
	return penalty, drivers
}

func liquidityRiskForTradeability(raw map[string]LiquidityState, daySymbols []string) (int, []string) {
	syms := uniqueSymbols(append([]string{"SPY", "QQQ"}, daySymbols...))
	risk, caution, usable := 0, 0, 0
	widest := 0.0
	for _, sym := range syms {
		x, ok := raw[sym]
		if !ok {
			continue
		}
		usable++
		if x.SpreadPct > widest {
			widest = x.SpreadPct
		}
		switch strings.ToUpper(x.State) {
		case "RISK":
			risk++
		case "CAUTION":
			caution++
		}
	}
	if usable == 0 {
		return 0, []string{"liquidity context unavailable"}
	}
	penalty := 0
	if risk > 0 {
		penalty += 15
	} else if caution > 0 {
		penalty += 7
	}
	if widest >= .75 {
		penalty += 8
	} else if widest >= .35 {
		penalty += 4
	}
	return penalty, []string{fmt.Sprintf("liquidity %d risk / %d caution across %d current symbols · widest spread %.2f%%", risk, caution, usable, widest)}
}

func optionsRiskForTradeability(options map[string]OptionsContext) (int, []string) {
	bearish, available := 0, 0
	ivTotal := 0.0
	ivN := 0
	for _, sym := range []string{"SPY", "QQQ"} {
		x, ok := options[sym]
		if !ok || strings.EqualFold(x.State, "UNAVAILABLE") {
			continue
		}
		available++
		if strings.EqualFold(x.Bias, "BEARISH") {
			bearish++
		}
		if x.AverageIV > 0 {
			ivTotal += x.AverageIV
			ivN++
		}
	}
	if available == 0 {
		return 0, []string{"index options context unavailable · no penalty"}
	}
	penalty := 0
	if bearish == 2 {
		penalty += 8
	} else if bearish == 1 {
		penalty += 3
	}
	avgIV := 0.0
	if ivN > 0 {
		avgIV = ivTotal / float64(ivN)
		if avgIV >= .50 {
			penalty += 5
		}
	}
	return penalty, []string{fmt.Sprintf("SPY/QQQ options %d/%d bearish%s", bearish, available, func() string {
		if avgIV > 0 {
			return fmt.Sprintf(" · avg IV %.0f%%", avgIV*100)
		}
		return ""
	}())}
}

// marketTradeability reconciles the originally approved market-level contract:
// Regime/VIX + breadth/internals + macro/event risk + liquidity/spreads +
// freshness + setup availability/extension + valid options context. Every input
// is reused from a canonical owner; this function performs no provider fetch and
// never mutates deterministic Day/Swing/Long Score/Action formulas.
func marketTradeabilityWithContext(quotes map[string]Quote, breadth MarketBreadthState, global GlobalMarketContext, ctx MarketTradeabilityContext, now time.Time) MarketTradeabilityState {
	out := MarketTradeabilityState{State: "DATA DEGRADED", Score: 0, Components: map[string]int{}, UpdatedAt: now.UnixMilli()}
	missing := []string{}
	stamps := []int64{}
	for _, sym := range []string{"SPY", "QQQ", "VIX"} {
		q := quotes[sym]
		stamp, ok, _ := quoteIsCurrentForMarketIntelligence(q, now)
		if !ok {
			missing = append(missing, sym)
		} else {
			stamps = append(stamps, stamp)
		}
	}
	if len(missing) > 0 {
		out.Blockers = append(out.Blockers, "Core context unavailable/stale: "+strings.Join(missing, ", "))
		out.Detail = "SPY + QQQ + VIX current evidence is required before DE.PULSE assigns a market Tradeability state."
		return out
	}
	score := 80
	drivers := []string{}
	vix := quotes["VIX"].Price
	volPenalty := 0
	switch {
	case vix >= 35:
		volPenalty = 45
	case vix >= 25:
		volPenalty = 25
	case vix >= 20:
		volPenalty = 10
	}
	score -= volPenalty
	out.Components["volatility"] = 100 - minInt(volPenalty*2, 100)
	drivers = append(drivers, fmt.Sprintf("VIX %.1f", vix))

	breadthPenalty := 0
	if breadth.State == "DEGRADED" || breadth.State == "UNAVAILABLE" {
		breadthPenalty += 15
		drivers = append(drivers, "breadth degraded")
	} else if breadth.ParticipationPct != nil {
		p := *breadth.ParticipationPct
		if p < 35 {
			breadthPenalty += 25
		} else if p < 45 {
			breadthPenalty += 10
		} else if p > 65 {
			score += 5
		}
		drivers = append(drivers, fmt.Sprintf("tracked breadth %.0f%%", p))
	}
	if x := breadth.Internals.Above50MAPct; x != nil && *x < 40 {
		breadthPenalty += 8
		drivers = append(drivers, fmt.Sprintf("above-50D participation %.0f%%", *x))
	}
	if x := breadth.Internals.SectorParticipationPct; x != nil && *x < 40 {
		breadthPenalty += 6
		drivers = append(drivers, fmt.Sprintf("sector participation %.0f%%", *x))
	}
	score -= breadthPenalty
	out.Components["breadth"] = 100 - minInt(breadthPenalty*2, 100)

	eventPenalty := 0
	if ctx.EventMode.Active {
		eventPenalty = 20
		drivers = append(drivers, fmt.Sprintf("high-impact event %s · %s", ctx.EventMode.Name, strings.ToLower(ctx.EventMode.Phase)))
	}
	score -= eventPenalty
	out.Components["eventRisk"] = 100 - eventPenalty*3

	liqPenalty, liqDrivers := liquidityRiskForTradeability(ctx.Liquidity, ctx.DaySymbols)
	score -= liqPenalty
	drivers = append(drivers, liqDrivers...)
	out.Components["liquidity"] = 100 - minInt(liqPenalty*3, 100)

	freshPenalty, freshDrivers := freshnessRiskForTradeability(ctx.Freshness)
	score -= freshPenalty
	drivers = append(drivers, freshDrivers...)
	out.Components["freshness"] = 100 - minInt(freshPenalty*3, 100)

	setupPenalty := 0
	if day, ok := ctx.Structure["day"]; ok && day.State == "EXTENDED" {
		setupPenalty += 10
		drivers = append(drivers, "SPY intraday structure extended")
	}
	if strings.EqualFold(ctx.Scanner.Status, "complete") && strings.EqualFold(ctx.Scanner.Mode, "day") && now.UnixMilli()-ctx.Scanner.UpdatedAt < int64(30*time.Minute/time.Millisecond) {
		qualified := 0
		for _, row := range ctx.Scanner.Results {
			if row.Score >= 55 {
				qualified++
			}
		}
		if qualified == 0 {
			setupPenalty += 10
			drivers = append(drivers, "recent Day scan found no qualified setups")
		} else {
			drivers = append(drivers, fmt.Sprintf("recent Day scan %d qualified setups", qualified))
		}
	}
	score -= setupPenalty
	out.Components["setups"] = 100 - minInt(setupPenalty*4, 100)

	optionsPenalty, optionDrivers := optionsRiskForTradeability(ctx.Options)
	score -= optionsPenalty
	drivers = append(drivers, optionDrivers...)
	out.Components["options"] = 100 - minInt(optionsPenalty*4, 100)

	globalPenalty := 0
	tone := strings.ToUpper(global.Tone)
	if strings.Contains(tone, "RISK-OFF") {
		globalPenalty = 15
		drivers = append(drivers, "global risk-off")
	} else if strings.Contains(tone, "CONSTRUCTIVE") || strings.Contains(tone, "RISK-ON") {
		score += 5
	}
	score -= globalPenalty
	out.Components["global"] = 100 - globalPenalty*4

	if score < 0 {
		score = 0
	}
	if score > 100 {
		score = 100
	}
	state := "TRADE NORMALLY"
	switch {
	case score < 35:
		state = "WAIT"
	case score < 50:
		state = "REDUCE SIZE"
	case score < 70:
		state = "SELECTIVE"
	}
	out.State, out.Score, out.Drivers = state, score, uniqueSymbols(drivers)
	out.Detail = fmt.Sprintf("Non-AI market-wide context score %d/100 reconciles volatility, breadth/internals, event risk, liquidity, freshness, setup availability/extension, options and global context. Tradeability is separate from ticker Trade Readiness and never mutates deterministic desk scores.", score)
	if len(stamps) > 0 {
		sort.Slice(stamps, func(i, j int) bool { return stamps[i] < stamps[j] })
		out.UpdatedAt = stamps[0]
	}
	return out
}

func appendRelativeStrengthRows(out []RelativeStrengthState, bars map[string]map[string][]Bar, symbol string, class SymbolClassification, horizon string, now time.Time) []RelativeStrengthState {
	benchmarks := []string{"SPY", "QQQ"}
	if class.SectorBenchmark != "" {
		benchmarks = append(benchmarks, class.SectorBenchmark)
	}
	if class.IndustryBenchmark != "" && class.IndustryBenchmark != class.SectorBenchmark {
		benchmarks = append(benchmarks, class.IndustryBenchmark)
	}
	for _, benchmark := range uniqueSymbols(benchmarks) {
		out = append(out, relativeStrengthFor(bars, symbol, benchmark, horizon, now))
	}
	return out
}

func buildMarketIntelligenceSnapshotWithContext(st AppState, quotes map[string]Quote, bars map[string]map[string][]Bar, rawLiquidity map[string]LiquidityState, global GlobalMarketContext, tradeCtx MarketTradeabilityContext, now time.Time) MarketIntelligenceSnapshot {
	selected := normalizeSymbol(st.UI.SelectedTicker)
	if selected == "" {
		selected = "SPY"
	}
	breadth := marketBreadthStateWithBars(quotes, bars, now)
	class := classificationForSymbol(selected)
	symbols := append([]string{}, activeDeskSymbolsFromState(st)...)
	symbols = append(symbols, discoverySymbolsFromState(st)...)
	symbols = append(symbols, selected)
	rs := []RelativeStrengthState{}
	for _, symbol := range uniqueSymbols(symbols) {
		c := classificationForSymbol(symbol)
		rs = appendRelativeStrengthRows(rs, bars, symbol, c, "day", now)
		rs = appendRelativeStrengthRows(rs, bars, symbol, c, "swing", now)
	}
	structure := map[string]MarketStructureState{"day": marketStructureFor(bars, "day", now), "swing": marketStructureFor(bars, "swing", now), "long": marketStructureFor(bars, "long", now)}
	sector := benchmarkRegime("Sector", class.Sector, class.SectorBenchmark, bars, now)
	industry := benchmarkRegime("Industry", class.Industry, class.IndustryBenchmark, bars, now)
	liq := liquiditySlippageState(selected, rawLiquidity)
	tradeCtx.Structure = structure
	tradeCtx.Liquidity = rawLiquidity
	if wl, ok := watchlistValueByID(st.Watchlists, st.Settings.DayWatchlistID); ok {
		tradeCtx.DaySymbols = userTradingSymbols(wl.Symbols)
	}
	return MarketIntelligenceSnapshot{Tradeability: marketTradeabilityWithContext(quotes, breadth, global, tradeCtx, now), Structure: structure, Breadth: breadth, SelectedSymbol: selected, Classification: class, RelativeStrength: rs, SectorRegime: sector, IndustryRegime: industry, Liquidity: liq, UpdatedAt: now.UnixMilli()}
}

// Compatibility wrapper preserves inherited tradeability regression ownership while v16.7 adds richer context through the explicit WithContext owner.
func marketTradeability(quotes map[string]Quote, breadth MarketBreadthState, global GlobalMarketContext, now time.Time) MarketTradeabilityState {
	return marketTradeabilityWithContext(quotes, breadth, global, MarketTradeabilityContext{}, now)
}
