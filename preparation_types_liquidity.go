package main

import (
	"fmt"
	"math"
	"time"
)

var alpacaTradingBaseURL = "https://api.alpaca.markets"

// v14.3.0 Improvement Build adds contextual intelligence only. Deterministic
// Day/Swing/Long setup-score formulas remain untouched; these structures feed
// regime/readiness/queue/research and operational surfaces.
type CheckpointException struct {
	Symbol    string `json:"symbol,omitempty"`
	Reason    string `json:"reason"`
	Severity  string `json:"severity,omitempty"`
	Target    string `json:"target,omitempty"` // research / earnings / sec / day / swing / long
	Source    string `json:"source,omitempty"`
	UpdatedAt int64  `json:"updatedAt,omitempty"`
}

type PreparationJobStatus struct {
	Key          string                `json:"key"`
	Label        string                `json:"label"`
	State        string                `json:"state"`
	Detail       string                `json:"detail,omitempty"`
	LastAttempt  int64                 `json:"lastAttempt,omitempty"`
	LastSuccess  int64                 `json:"lastSuccess,omitempty"`
	AttemptCount int                   `json:"attemptCount,omitempty"`
	NextWindow   int64                 `json:"nextWindow,omitempty"`
	Window       string                `json:"window,omitempty"`
	TradingDay   string                `json:"tradingDay,omitempty"`
	Late         bool                  `json:"late,omitempty"`
	Attention    string                `json:"attention,omitempty"`
	Summary      []string              `json:"summary,omitempty"`
	Changed      []string              `json:"changed,omitempty"`
	Exceptions   []CheckpointException `json:"exceptions,omitempty"`
}

type LiquidityBaseline struct {
	NormalSpreadPct float64 `json:"normalSpreadPct,omitempty"`
	Samples         int     `json:"samples"`
	UpdatedAt       int64   `json:"updatedAt,omitempty"`
}

type LiquidityState struct {
	Symbol                     string   `json:"symbol"`
	State                      string   `json:"state"`
	Spread                     float64  `json:"spread,omitempty"`
	SpreadPct                  float64  `json:"spreadPercent,omitempty"`
	NormalSpreadPct            float64  `json:"normalSpreadPercent,omitempty"`
	RelativeSpreadMultiple     float64  `json:"relativeSpreadMultiple,omitempty"`
	SpreadBaselineSamples      int      `json:"spreadBaselineSamples,omitempty"`
	DollarVolume               float64  `json:"dollarVolume,omitempty"`
	AverageDollarVolume20D     float64  `json:"averageDollarVolume20d,omitempty"`
	VolumeBasis                string   `json:"volumeBasis,omitempty"`
	OpeningLiquidity           string   `json:"openingLiquidity,omitempty"`
	DailyRangeVolatilityPct    float64  `json:"dailyRangeVolatilityPct,omitempty"`
	VolatilityAdjustedSlippage float64  `json:"volatilityAdjustedSlippage,omitempty"`
	SlippageRisk               string   `json:"slippageRisk,omitempty"`
	BidSize                    float64  `json:"bidSize,omitempty"`
	AskSize                    float64  `json:"askSize,omitempty"`
	QuoteAgeMs                 int64    `json:"quoteAgeMs,omitempty"`
	Limitations                []string `json:"limitations,omitempty"`
	Detail                     string   `json:"detail,omitempty"`
	UpdatedAt                  int64    `json:"updatedAt,omitempty"`
}

type DerivedIntelligenceState struct {
	Key       string `json:"key"`
	Label     string `json:"label"`
	State     string `json:"state"`
	Detail    string `json:"detail"`
	Source    string `json:"source"`
	UpdatedAt int64  `json:"updatedAt,omitempty"`
}

type ProviderCapabilityEntry struct {
	Provider      string   `json:"provider"`
	Capability    string   `json:"capability"`
	Status        string   `json:"status"` // AVAILABLE / PLAN LIMITED / NOT ENTITLED / TEMPORARILY UNAVAILABLE
	Detail        string   `json:"detail,omitempty"`
	UpdatedAt     int64    `json:"updatedAt,omitempty"`
	Uses          []string `json:"uses,omitempty"`
	Priority      string   `json:"priority,omitempty"`
	Quota         string   `json:"quota,omitempty"`
	RateLimit     string   `json:"rateLimit,omitempty"`
	LatencyMs     int64    `json:"latencyMs,omitempty"`
	LastSuccess   int64    `json:"lastSuccess,omitempty"`
	LastFailure   int64    `json:"lastFailure,omitempty"`
	FailureCount  int      `json:"failureCount,omitempty"`
	Attempts      int64    `json:"attempts,omitempty"`
	LastError     string   `json:"lastError,omitempty"`
	Recovery      string   `json:"recovery,omitempty"`
	ExpectedDelay string   `json:"expectedDelay,omitempty"`
	CostClass     string   `json:"costClass,omitempty"`
}

type EarningsSurprisePoint struct {
	Period          string  `json:"period,omitempty"`
	Actual          float64 `json:"actual,omitempty"`
	Estimate        float64 `json:"estimate,omitempty"`
	Surprise        float64 `json:"surprise,omitempty"`
	SurprisePercent float64 `json:"surprisePercent,omitempty"`
}

type SymbolIntelligence struct {
	Symbol                 string                  `json:"symbol"`
	Peers                  []string                `json:"peers,omitempty"`
	EarningsSurprises      []EarningsSurprisePoint `json:"earningsSurprises,omitempty"`
	ConsecutiveBeats       int                     `json:"consecutiveBeats,omitempty"`
	ConsecutiveMisses      int                     `json:"consecutiveMisses,omitempty"`
	RecommendationTrend    string                  `json:"recommendationTrend,omitempty"`
	PriceTarget            float64                 `json:"priceTarget,omitempty"`
	InsiderNetShares       float64                 `json:"insiderNetShares,omitempty"`
	InstitutionalOwners    int                     `json:"institutionalOwners,omitempty"`
	InstitutionalNetChange float64                 `json:"institutionalNetChange,omitempty"`
	EntitlementNotes       []string                `json:"entitlementNotes,omitempty"`
	UpdatedAt              int64                   `json:"updatedAt,omitempty"`
}

type CatalystReactionState struct {
	Symbol            string             `json:"symbol"`
	TriggerType       string             `json:"triggerType"`
	Trigger           string             `json:"trigger"`
	Session           string             `json:"session"`
	State             string             `json:"state"`
	Phase             string             `json:"phase,omitempty"` // ARMED / TRIGGERED / PREMARKET REACTION / OPENING REACTION / 5m / 15m / 30m / 60m / SESSION REACTION / COMPLETE
	TriggerAt         int64              `json:"triggerAt"`
	TriggerPrice      float64            `json:"triggerPrice,omitempty"`
	LatestPrice       float64            `json:"latestPrice,omitempty"`
	MovePercent       float64            `json:"movePercent,omitempty"`
	GapPercent        float64            `json:"gapPercent,omitempty"`
	Volume            float64            `json:"volume,omitempty"`
	RelativeVolume    float64            `json:"relativeVolume,omitempty"`
	SpreadPct         float64            `json:"spreadPercent,omitempty"`
	Liquidity         string             `json:"liquidity,omitempty"`
	VWAPState         string             `json:"vwapState,omitempty"`
	OpeningRangeState string             `json:"openingRangeState,omitempty"`
	HoldFadeState     string             `json:"holdFadeState,omitempty"`
	VolatilityState   string             `json:"volatilityState,omitempty"`
	ReactionPercent   map[string]float64 `json:"reactionPercent,omitempty"` // 5m / 15m / 30m / 60m / session
	Flags             []string           `json:"flags,omitempty"`
	Detail            string             `json:"detail,omitempty"`
	UpdatedAt         int64              `json:"updatedAt,omitempty"`
	CompletedAt       int64              `json:"completedAt,omitempty"`
}

type AlpacaCalendarDay struct {
	Date  string `json:"date"`
	Open  string `json:"open"`
	Close string `json:"close"`
}

type PremarketSnapshot struct {
	Symbol       string  `json:"symbol"`
	Open         float64 `json:"open,omitempty"`
	High         float64 `json:"high,omitempty"`
	Low          float64 `json:"low,omitempty"`
	Last         float64 `json:"last,omitempty"`
	Volume       float64 `json:"volume,omitempty"`
	GapPercent   float64 `json:"gapPercent,omitempty"`
	RangePercent float64 `json:"rangePercent,omitempty"`
	Bars         int     `json:"bars,omitempty"`
	UpdatedAt    int64   `json:"updatedAt,omitempty"`
}

type MarketOpenCheckpoint struct {
	State               string                       `json:"state"`
	Attention           string                       `json:"attention,omitempty"`
	RunAt               int64                        `json:"runAt,omitempty"`
	Late                bool                         `json:"late,omitempty"`
	VIX                 float64                      `json:"vix,omitempty"`
	VIXState            string                       `json:"vixState,omitempty"`
	GlobalTone          string                       `json:"globalTone,omitempty"`
	MacroStates         map[string]string            `json:"macroStates,omitempty"`
	Premarket           map[string]PremarketSnapshot `json:"premarket,omitempty"`
	OptionsContexts     int                          `json:"optionsContexts,omitempty"`
	DayCandidates       int                          `json:"dayCandidates,omitempty"`
	SwingContextChanges []string                     `json:"swingContextChanges,omitempty"`
	LongContextChanges  []string                     `json:"longContextChanges,omitempty"`
	Detail              []string                     `json:"detail,omitempty"`
	Changed             []string                     `json:"changed,omitempty"`
	Exceptions          []CheckpointException        `json:"exceptions,omitempty"`
}

type MarketMover struct {
	Symbol        string  `json:"symbol"`
	ChangePercent float64 `json:"changePercent,omitempty"`
	Volume        float64 `json:"volume,omitempty"`
}

type MarketActivityState struct {
	MostActive []MarketMover `json:"mostActive,omitempty"`
	Gainers    []MarketMover `json:"gainers,omitempty"`
	Losers     []MarketMover `json:"losers,omitempty"`
	Status     string        `json:"status,omitempty"`
	UpdatedAt  int64         `json:"updatedAt,omitempty"`
}

type CorporateAction struct {
	ID               string  `json:"id,omitempty"`
	Symbol           string  `json:"symbol"`
	Type             string  `json:"type"`
	ProcessDate      string  `json:"processDate,omitempty"`
	ExDate           string  `json:"exDate,omitempty"`
	RecordDate       string  `json:"recordDate,omitempty"`
	PayableDate      string  `json:"payableDate,omitempty"`
	OldSymbol        string  `json:"oldSymbol,omitempty"`
	NewSymbol        string  `json:"newSymbol,omitempty"`
	Ratio            float64 `json:"ratio,omitempty"`
	CashAmount       float64 `json:"cashAmount,omitempty"`
	AdjustmentFactor float64 `json:"adjustmentFactor,omitempty"`
	Status           string  `json:"status,omitempty"`
	FirstSeenAt      int64   `json:"firstSeenAt,omitempty"`
	UpdatedAt        int64   `json:"updatedAt,omitempty"`
	Detail           string  `json:"detail,omitempty"`
	Source           string  `json:"source,omitempty"`
}

func easternLocation() *time.Location {
	loc, err := time.LoadLocation("America/New_York")
	if err != nil {
		return time.FixedZone("ET", -5*60*60)
	}
	return loc
}

func nextWindowAt(now time.Time, hour, minute int) int64 {
	loc := easternLocation()
	et := now.In(loc)
	for i := 0; i < 10; i++ {
		d := et.AddDate(0, 0, i)
		if d.Weekday() < time.Monday || d.Weekday() > time.Friday || isUSMarketHoliday(d) {
			continue
		}
		t := time.Date(d.Year(), d.Month(), d.Day(), hour, minute, 0, 0, loc)
		if t.After(et) {
			return t.UnixMilli()
		}
	}
	return 0
}

func initialPreparationJobs(now time.Time) map[string]PreparationJobStatus {
	return map[string]PreparationJobStatus{
		"pre-market-prep":  {Key: "pre-market-prep", Label: "Pre-Market Prep", State: "NOT RUN YET", Window: "3:15–3:50 AM ET", NextWindow: nextWindowAt(now, 3, 15), Detail: "Targeted stale-data preparation; live/overnight streams stay active."},
		"market-open-prep": {Key: "market-open-prep", Label: "Market Open Prep", State: "NOT RUN YET", Window: "9:20–9:25 AM ET", NextWindow: nextWindowAt(now, 9, 20), Detail: "Synchronized pre-bell trading-readiness checkpoint."},
		"catalyst-watch":   {Key: "catalyst-watch", Label: "Earnings & Material Catalyst Watch", State: "READY", Window: "Event-driven only", Detail: "Arms only for scheduled earnings or a material unexpected catalyst."},
	}
}

func marketOpenPrepWindow(now time.Time) bool {
	et := now.In(easternLocation())
	mins := et.Hour()*60 + et.Minute()
	return et.Weekday() >= time.Monday && et.Weekday() <= time.Friday && !isUSMarketHoliday(et) && mins >= 9*60+20 && mins <= 9*60+25
}

func liquidityHistoryMetrics(symbol string, bars map[string]map[string][]Bar, now time.Time) (dollarVolume, avgDollar20, rangePct float64, volumeBasis string) {
	frames := bars[symbol]
	daily := completedBarsBefore(frames["daily"], now.UnixMilli(), "daily")
	start := len(daily) - 20
	if start < 0 {
		start = 0
	}
	count := 0
	for _, b := range daily[start:] {
		if b.C <= 0 {
			continue
		}
		if b.V > 0 {
			avgDollar20 += b.C * b.V
		}
		if b.H > b.L && b.C > 0 {
			rangePct += (b.H - b.L) / b.C * 100
		}
		count++
	}
	if count > 0 {
		avgDollar20 /= float64(count)
		rangePct /= float64(count)
	}
	et := now.In(easternLocation())
	day := et.Format("2006-01-02")
	sessionVolume := 0.0
	lastPrice := 0.0
	for _, b := range frames["intraday"] {
		t := time.UnixMilli(normalizedBarTimestampMs(b.T)).In(easternLocation())
		if t.Format("2006-01-02") != day {
			continue
		}
		if b.V > 0 {
			sessionVolume += b.V
		}
		if b.C > 0 {
			lastPrice = b.C
		}
	}
	if sessionVolume > 0 && lastPrice > 0 {
		dollarVolume = sessionVolume * lastPrice
		volumeBasis = "Canonical intraday bars · current session"
	}
	return
}

func deriveLiquidityStatesWithContext(quotes map[string]Quote, bars map[string]map[string][]Bar, baselines map[string]LiquidityBaseline, now time.Time) map[string]LiquidityState {
	out := map[string]LiquidityState{}
	nowMs := now.UnixMilli()
	for sym, q := range quotes {
		if q.Price <= 0 {
			continue
		}
		spread, spreadPct := 0.0, 0.0
		if q.Ask > 0 && q.Bid > 0 && q.Ask >= q.Bid {
			spread = q.Ask - q.Bid
			mid := (q.Ask + q.Bid) / 2
			if mid > 0 {
				spreadPct = spread / mid * 100
			}
		}
		stamp := q.ProviderTimestamp
		if stamp == 0 {
			stamp = q.UpdatedAt
		}
		age := int64(0)
		if stamp > 0 {
			age = nowMs - stamp
			if age < 0 {
				age = 0
			}
		}
		validTime := stamp > 0 && stamp <= nowMs+30_000
		hasValidSpread := q.Bid > 0 && q.Ask > 0 && q.Ask >= q.Bid
		base := baselines[sym]
		relative := 0.0
		limits := []string{}
		if base.Samples >= 5 && base.NormalSpreadPct > 0 && spreadPct > 0 {
			relative = spreadPct / base.NormalSpreadPct
		} else {
			limits = append(limits, "Relative-to-normal spread is warming up until at least 5 canonical bid/ask observations exist.")
		}
		dollarVol, avgDollar, rangePct, volumeBasis := liquidityHistoryMetrics(sym, bars, now)
		if dollarVol == 0 && q.Volume > 0 {
			dollarVol = q.Price * q.Volume
			volumeBasis = "Provider quote volume; semantics may be trade-size or session volume depending on feed"
			limits = append(limits, "Dollar volume uses provider quote volume because current-session intraday aggregation is unavailable.")
		}
		volAdj := 0.0
		if spreadPct > 0 && rangePct > 0 {
			volAdj = spreadPct / rangePct
		}
		opening := "NOT IN OPENING WINDOW"
		et := now.In(easternLocation())
		mins := et.Hour()*60 + et.Minute()
		if marketSessionET(now) == "regular" && mins <= 10*60+30 {
			opening = "NORMAL"
			if !hasValidSpread || age > 120000 {
				opening = "DEGRADED"
			} else if spreadPct >= 0.30 || relative >= 2 {
				opening = "THIN / WIDE"
			}
		}
		state, detail := "UNKNOWN", "Bid/ask spread evidence unavailable; liquidity/slippage risk is unknown."
		slippage := "UNKNOWN"
		if !validTime {
			state, slippage, detail = "RISK", "HIGH", "Quote timestamp is missing or future-skewed; do not treat liquidity as current."
		} else if age > 120000 {
			state, slippage, detail = "RISK", "HIGH", "Quote is stale; do not treat liquidity as current."
		} else if !hasValidSpread {
			state, slippage, detail = "UNKNOWN", "UNKNOWN", "Valid bid/ask spread evidence is unavailable; liquidity/slippage risk is unknown."
		} else if spreadPct >= 0.75 || relative >= 3 || volAdj >= 0.20 {
			state, slippage, detail = "RISK", "HIGH", "Spread/slippage is materially wide versus current price, learned normal spread, or recent daily range."
		} else if spreadPct >= 0.30 || age > 30000 || relative >= 1.75 || volAdj >= 0.10 {
			state, slippage, detail = "CAUTION", "ELEVATED", "Spread/quote age or volatility-adjusted slippage needs confirmation."
		} else {
			state, slippage, detail = "HEALTHY", "LOW", "Quote current; spread/slippage is within available learned and volatility-adjusted guardrails."
		}
		out[sym] = LiquidityState{Symbol: sym, State: state, Spread: spread, SpreadPct: spreadPct, NormalSpreadPct: base.NormalSpreadPct, RelativeSpreadMultiple: relative, SpreadBaselineSamples: base.Samples, DollarVolume: dollarVol, AverageDollarVolume20D: avgDollar, VolumeBasis: volumeBasis, OpeningLiquidity: opening, DailyRangeVolatilityPct: rangePct, VolatilityAdjustedSlippage: volAdj, SlippageRisk: slippage, BidSize: q.BidSize, AskSize: q.AskSize, QuoteAgeMs: age, Limitations: limits, Detail: detail, UpdatedAt: stamp}
	}
	return out
}

// currentLiquidityMarketRisk separates actual current spread/slippage risk from
// stale/missing quote evidence. A stale quote is a DATA FRESHNESS problem and
// must not be multiplied into one apparent LIQUIDITY RISK per symbol.
func currentLiquidityMarketRisk(l LiquidityState, now time.Time) bool {
	if l.State != "RISK" || l.UpdatedAt <= 0 {
		return false
	}
	nowMs := now.UnixMilli()
	if l.UpdatedAt > nowMs+30_000 || l.QuoteAgeMs > 120_000 {
		return false
	}
	return true
}

func metric(metrics map[string]MacroMetric, key string) (MacroMetric, bool) {
	m, ok := metrics[key]
	return m, ok
}
func deriveIntelligenceStates(metrics map[string]MacroMetric) map[string]DerivedIntelligenceState {
	out := map[string]DerivedIntelligenceState{}

	ten, tenOK := metric(metrics, "DGS10")
	if !tenOK {
		ten, tenOK = metric(metrics, "UST10Y")
	}
	if tenOK {
		state := "STABLE"
		d := ten.Change5D
		if math.Abs(d) >= .25 {
			if d > 0 {
				state = "SHOCK HIGHER"
			} else {
				state = "SHOCK LOWER"
			}
		} else if d >= .08 {
			state = "RISING"
		} else if d <= -.08 {
			state = "FALLING"
		}
		curve := ""
		two, twoOK := metric(metrics, "DGS2")
		if !twoOK {
			two, twoOK = metric(metrics, "UST2Y")
		}
		if twoOK {
			curve = fmt.Sprintf(" · 10Y−2Y %.2fpp", ten.Value-two.Value)
		}
		source := defaultString(ten.Source, "FRED / Treasury")
		out["rates"] = DerivedIntelligenceState{Key: "rates", Label: "Rates State", State: state, Detail: fmt.Sprintf("10Y %.2f%% · 5D %+.2fpp%s", ten.Value, d, curve), Source: source, UpdatedAt: ten.UpdatedAt}
	}
	if hy, ok := metric(metrics, "BAMLH0A0HYM2"); ok {
		state := "HEALTHY"
		if hy.Change20D >= .45 {
			state = "STRESSED"
		} else if hy.Change20D >= .15 {
			state = "WIDENING"
		}
		out["credit"] = DerivedIntelligenceState{Key: "credit", Label: "Credit State", State: state, Detail: fmt.Sprintf("HY OAS %.2f%% · 20D %+.2fpp", hy.Value, hy.Change20D), Source: "FRED", UpdatedAt: hy.UpdatedAt}
	}
	if nfci, ok := metric(metrics, "NFCI"); ok {
		state := "NEUTRAL"
		if nfci.Value > .25 {
			state = "TIGHT"
		} else if nfci.Value < -.25 {
			state = "LOOSE"
		}
		out["financial-conditions"] = DerivedIntelligenceState{Key: "financial-conditions", Label: "Financial Conditions", State: state, Detail: fmt.Sprintf("NFCI %.2f", nfci.Value), Source: "FRED", UpdatedAt: nfci.UpdatedAt}
	}
	if usd, ok := metric(metrics, "DTWEXBGS"); ok {
		state := "NEUTRAL"
		if usd.Change20D > 1.5 {
			state = "STRONG"
		} else if usd.Change20D < -1.5 {
			state = "WEAK"
		}
		out["dollar"] = DerivedIntelligenceState{Key: "dollar", Label: "Dollar State", State: state, Detail: fmt.Sprintf("Broad USD %.2f · 20D %+.2f", usd.Value, usd.Change20D), Source: "FRED", UpdatedAt: usd.UpdatedAt}
	}
	if cpi, ok := metric(metrics, "CPI_INDEX"); ok {
		state := "MIXED"
		if cpi.Previous != 0 && cpi.Value < cpi.Previous-.1 {
			state = "COOLING"
		} else if cpi.Previous != 0 && cpi.Value > cpi.Previous+.1 {
			state = "HEATING"
		}
		detail := fmt.Sprintf("Headline %.2f%% YoY", cpi.Value)
		if core, ok2 := metric(metrics, "CORE_CPI_INDEX"); ok2 {
			detail += fmt.Sprintf(" · Core %.2f%%", core.Value)
		}
		out["inflation"] = DerivedIntelligenceState{Key: "inflation", Label: "Inflation State", State: state, Detail: detail, Source: "BLS", UpdatedAt: cpi.UpdatedAt}
	}
	if un, ok := metric(metrics, "UNEMP"); ok {
		state := "BALANCED"
		detail := fmt.Sprintf("Unemployment %.1f%%", un.Value)
		if pay, ok2 := metric(metrics, "NONFARM"); ok2 {
			detail += fmt.Sprintf(" · payroll %+.0fk", pay.Value)
			if pay.Value < 75 {
				state = "WEAKENING"
			} else if pay.Value > 250 {
				state = "STRONG"
			}
		}
		if wage, ok2 := metric(metrics, "AHE"); ok2 {
			detail += fmt.Sprintf(" · AHE %.2f", wage.Value)
		}
		out["labor"] = DerivedIntelligenceState{Key: "labor", Label: "Labor State", State: state, Detail: detail, Source: "BLS", UpdatedAt: un.UpdatedAt}
	}
	if crude, ok := metric(metrics, "CRUDE_STOCKS"); ok {
		state := "BALANCED"
		if crude.Previous != 0 {
			delta := crude.Value - crude.Previous
			if delta < -2 {
				state = "TIGHTENING"
			} else if delta > 2 {
				state = "LOOSENING"
			}
		}
		out["energy"] = DerivedIntelligenceState{Key: "energy", Label: "Energy State", State: state, Detail: fmt.Sprintf("Crude stocks %.1f · latest Δ %.1f", crude.Value, crude.Value-crude.Previous), Source: "EIA", UpdatedAt: crude.UpdatedAt}
	}
	return out
}
