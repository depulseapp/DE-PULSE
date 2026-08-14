package main

import (
	"fmt"
	"math"
	"sort"
	"strings"
	"time"
)

// v16.5 Context & Alternative Intelligence is one derived context layer over
// existing canonical stores. It does not own provider polling and cannot mutate
// deterministic Day/Swing/Long Score/Action.

type SentimentComponent struct {
	Name       string   `json:"name"`
	State      string   `json:"state"`
	Score      *float64 `json:"score,omitempty"`
	Weight     float64  `json:"weight"`
	Confidence int      `json:"confidence"`
	Source     string   `json:"source"`
	UpdatedAt  int64    `json:"updatedAt,omitempty"`
	Detail     string   `json:"detail,omitempty"`
}

type SentimentCompositeState struct {
	State              string               `json:"state"`
	Score              *float64             `json:"score,omitempty"`
	Confidence         int                  `json:"confidence"`
	ComponentsExpected int                  `json:"componentsExpected"`
	ComponentsUsed     int                  `json:"componentsUsed"`
	Components         []SentimentComponent `json:"components"`
	Missing            []string             `json:"missing,omitempty"`
	UpdatedAt          int64                `json:"updatedAt,omitempty"`
	Detail             string               `json:"detail"`
}

type HeatMapCell struct {
	Symbol        string  `json:"symbol"`
	Name          string  `json:"name"`
	Sector        string  `json:"sector,omitempty"`
	ChangePercent float64 `json:"changePercent,omitempty"`
	State         string  `json:"state"`
	Freshness     string  `json:"freshness"`
	Source        string  `json:"source,omitempty"`
	UpdatedAt     int64   `json:"updatedAt,omitempty"`
}

type HeatMapMode struct {
	Key            string        `json:"key"`
	Label          string        `json:"label"`
	Universe       string        `json:"universe"`
	CoverageBasis  string        `json:"coverageBasis"`
	Expected       int           `json:"expected"`
	Fresh          int           `json:"fresh"`
	Stale          int           `json:"stale"`
	Missing        int           `json:"missing"`
	CoveragePct    float64       `json:"coveragePct"`
	DirectionalPct *float64      `json:"directionalPct,omitempty"`
	Cells          []HeatMapCell `json:"cells"`
	UpdatedAt      int64         `json:"updatedAt,omitempty"`
	Detail         string        `json:"detail"`
}

type MarketSectorHeatMap struct {
	// The legacy top-level fields remain the canonical sector mode so existing
	// sentiment and clients keep a single stable owner. v16.8 adds explicit
	// context modes without duplicating quote fetch or signal ownership.
	Universe       string                 `json:"universe"`
	Expected       int                    `json:"expected"`
	Fresh          int                    `json:"fresh"`
	Stale          int                    `json:"stale"`
	Missing        int                    `json:"missing"`
	CoveragePct    float64                `json:"coveragePct"`
	DirectionalPct *float64               `json:"directionalPct,omitempty"`
	Cells          []HeatMapCell          `json:"cells"`
	UpdatedAt      int64                  `json:"updatedAt,omitempty"`
	Detail         string                 `json:"detail"`
	Modes          map[string]HeatMapMode `json:"modes,omitempty"`
}

type GEXContextState struct {
	Symbol              string                 `json:"symbol"`
	State               string                 `json:"state"`
	Quality             string                 `json:"quality,omitempty"`
	NetGEX              float64                `json:"netGex,omitempty"`
	CallGEX             float64                `json:"callGex,omitempty"`
	PutGEX              float64                `json:"putGex,omitempty"`
	CoveragePct         float64                `json:"coveragePct,omitempty"`
	GammaOIContracts    int                    `json:"gammaOiContracts,omitempty"`
	GammaContracts      int                    `json:"gammaContracts,omitempty"`
	TotalContracts      int                    `json:"totalContracts,omitempty"`
	OpenInterestDate    string                 `json:"openInterestDate,omitempty"`
	Feed                string                 `json:"feed,omitempty"`
	Source              string                 `json:"source,omitempty"`
	UpdatedAt           int64                  `json:"updatedAt,omitempty"`
	Detail              string                 `json:"detail"`
	Limitations         []string               `json:"limitations,omitempty"`
	UnderlyingPrice     float64                `json:"underlyingPrice,omitempty"`
	MajorGammaStrikes   []GEXStrikeLevel       `json:"majorGammaStrikes,omitempty"`
	GammaZones          []GEXConcentrationZone `json:"gammaZones,omitempty"`
	ExpirationGEX       []GEXExpirationLevel   `json:"expirationGex,omitempty"`
	GammaFlip           *float64               `json:"gammaFlip,omitempty"`
	GammaFlipMethod     string                 `json:"gammaFlipMethod,omitempty"`
	DeterministicImpact string                 `json:"deterministicImpact"`
}

type CommunityEvidenceItem struct {
	ID              string `json:"id"`
	Symbol          string `json:"symbol,omitempty"`
	Source          string `json:"source"`
	Platform        string `json:"platform,omitempty"`
	IngestionMode   string `json:"ingestionMode,omitempty"`
	RightsStatus    string `json:"rightsStatus,omitempty"`
	AIEligibility   string `json:"aiEligibility,omitempty"`
	RetentionPolicy string `json:"retentionPolicy,omitempty"`
	Fingerprint     string `json:"fingerprint,omitempty"`
	Stance          string `json:"stance"`
	Text            string `json:"text"`
	URL             string `json:"url,omitempty"`
	ObservedAt      int64  `json:"observedAt,omitempty"`
	SubmittedAt     int64  `json:"submittedAt"`
}

type CommunitySourcePolicy struct {
	Platform        string `json:"platform"`
	AccessMode      string `json:"accessMode"`
	RightsStatus    string `json:"rightsStatus"`
	AIEligibility   string `json:"aiEligibility"`
	RetentionPolicy string `json:"retentionPolicy"`
	Status          string `json:"status"`
	Detail          string `json:"detail"`
}

type CommunityEvidenceCluster struct {
	ID                string   `json:"id"`
	Symbol            string   `json:"symbol,omitempty"`
	Representative    string   `json:"representative"`
	Sources           []string `json:"sources"`
	Platforms         []string `json:"platforms"`
	URLs              []string `json:"urls,omitempty"`
	SourceDiversity   int      `json:"sourceDiversity"`
	Mentions          int      `json:"mentions"`
	MentionVelocity1H int      `json:"mentionVelocity1h"`
	Bullish           int      `json:"bullish"`
	Bearish           int      `json:"bearish"`
	Mixed             int      `json:"mixed"`
	Unknown           int      `json:"unknown"`
	Corroborated      bool     `json:"corroborated"`
	Corroboration     []string `json:"corroboration,omitempty"`
	AIEligibleItems   int      `json:"aiEligibleItems"`
	Materiality       string   `json:"materiality"`
	LatestAt          int64    `json:"latestAt,omitempty"`
	Detail            string   `json:"detail"`
}

type CommunityIntelligenceState struct {
	Label               string                     `json:"label"`
	State               string                     `json:"state"`
	Total               int                        `json:"total"`
	UniqueEvents        int                        `json:"uniqueEvents"`
	SourceDiversity     int                        `json:"sourceDiversity"`
	MentionVelocity1H   int                        `json:"mentionVelocity1h"`
	AIEligible          int                        `json:"aiEligible"`
	Bullish             int                        `json:"bullish"`
	Bearish             int                        `json:"bearish"`
	Mixed               int                        `json:"mixed"`
	Unknown             int                        `json:"unknown"`
	Items               []CommunityEvidenceItem    `json:"items"`
	Clusters            []CommunityEvidenceCluster `json:"clusters"`
	Policies            []CommunitySourcePolicy    `json:"policies"`
	UpdatedAt           int64                      `json:"updatedAt,omitempty"`
	Detail              string                     `json:"detail"`
	Untrusted           bool                       `json:"untrustedExternalContent"`
	DeterministicImpact string                     `json:"deterministicImpact"`
}

type EnergyInstrumentTruth struct {
	Name               string   `json:"name"`
	Contract           string   `json:"contract"`
	State              string   `json:"state"`
	ReferenceType      string   `json:"referenceType"`
	Value              *float64 `json:"value,omitempty"`
	Change5D           *float64 `json:"change5d,omitempty"`
	Change20D          *float64 `json:"change20d,omitempty"`
	Trend              string   `json:"trend"`
	Source             string   `json:"source,omitempty"`
	Delayed            bool     `json:"delayed"`
	ContinuousContract bool     `json:"continuousContract"`
	RollAdjusted       bool     `json:"rollAdjusted"`
	UpdatedAt          int64    `json:"updatedAt,omitempty"`
	Semantics          string   `json:"semantics"`
}

type OilEnergyContextState struct {
	State               string                `json:"state"`
	WTI                 *MacroMetric          `json:"wti,omitempty"`
	Brent               *MacroMetric          `json:"brent,omitempty"`
	WTIContext          EnergyInstrumentTruth `json:"wtiContext"`
	BrentContext        EnergyInstrumentTruth `json:"brentContext"`
	BrentWTISpread      *float64              `json:"brentWtiSpread,omitempty"`
	SpreadState         string                `json:"spreadState,omitempty"`
	CrudeStocks         *MacroMetric          `json:"crudeStocks,omitempty"`
	Production          *MacroMetric          `json:"production,omitempty"`
	RefineryUtil        *MacroMetric          `json:"refineryUtil,omitempty"`
	USO                 *Quote                `json:"uso,omitempty"`
	XLE                 *Quote                `json:"xle,omitempty"`
	EnergySectorState   string                `json:"energySectorState,omitempty"`
	USMarketRelevance   string                `json:"usMarketRelevance,omitempty"`
	ComponentsUsed      int                   `json:"componentsUsed"`
	UpdatedAt           int64                 `json:"updatedAt,omitempty"`
	Sources             []string              `json:"sources"`
	Limitations         []string              `json:"limitations"`
	Detail              string                `json:"detail"`
	DeterministicImpact string                `json:"deterministicImpact"`
}

type ContextAlternativeIntelligenceSnapshot struct {
	Sentiment SentimentCompositeState    `json:"sentiment"`
	HeatMap   MarketSectorHeatMap        `json:"heatMap"`
	GEX       map[string]GEXContextState `json:"gex"`
	Community CommunityIntelligenceState `json:"community"`
	OilEnergy OilEnergyContextState      `json:"oilEnergy"`
	UpdatedAt int64                      `json:"updatedAt,omitempty"`
}

type heatMapSymbolDef struct{ Symbol, Name, Sector string }

var v165SectorHeatUniverse = []heatMapSymbolDef{
	{"XLK", "Technology", "Technology"}, {"XLC", "Communication Services", "Communication Services"},
	{"XLY", "Consumer Discretionary", "Consumer Discretionary"}, {"XLP", "Consumer Staples", "Consumer Staples"},
	{"XLE", "Energy", "Energy"}, {"XLF", "Financials", "Financials"}, {"XLV", "Health Care", "Health Care"},
	{"XLI", "Industrials", "Industrials"}, {"XLB", "Materials", "Materials"}, {"XLRE", "Real Estate", "Real Estate"},
	{"XLU", "Utilities", "Utilities"},
}

func v165Clamp(v, lo, hi float64) float64 {
	if v < lo {
		return lo
	}
	if v > hi {
		return hi
	}
	return v
}

func v165FloatPtr(v float64) *float64 { x := v; return &x }

func v165QuoteStamp(q Quote) int64 {
	if q.ProviderTimestamp > 0 {
		return normalizeObservationMs(q.ProviderTimestamp)
	}
	return normalizeObservationMs(q.UpdatedAt)
}

func v168BuildHeatMode(key, label, universe, basis string, defs []heatMapSymbolDef, quotes map[string]Quote, now time.Time) HeatMapMode {
	out := HeatMapMode{Key: key, Label: label, Universe: universe, CoverageBasis: basis, Expected: len(defs), Cells: []HeatMapCell{}}
	adv := 0
	for _, def := range defs {
		q, ok := quotes[def.Symbol]
		cell := HeatMapCell{Symbol: def.Symbol, Name: def.Name, Sector: def.Sector, State: "UNAVAILABLE", Freshness: "UNAVAILABLE"}
		if !ok || q.Price <= 0 {
			out.Missing++
			out.Cells = append(out.Cells, cell)
			continue
		}
		stamp := v165QuoteStamp(q)
		freshness, _, _, _, _ := freshnessState("Quotes", providerFromQuoteSource(q.Source), marketSessionET(now), stamp, "", now.UnixMilli())
		cell.ChangePercent, cell.Freshness, cell.Source, cell.UpdatedAt = q.ChangePercent, freshness, q.Source, stamp
		switch freshness {
		case "LIVE", "FRESH", "DUE SOON", "IDLE", "DELAYED":
			out.Fresh++
			if q.ChangePercent >= 1.0 {
				cell.State = "STRONG"
			} else if q.ChangePercent >= 0.15 {
				cell.State = "POSITIVE"
			} else if q.ChangePercent <= -1.0 {
				cell.State = "WEAK"
			} else if q.ChangePercent <= -0.15 {
				cell.State = "NEGATIVE"
			} else {
				cell.State = "FLAT"
			}
			if q.ChangePercent > 0 {
				adv++
			}
			if stamp > out.UpdatedAt {
				out.UpdatedAt = stamp
			}
		default:
			out.Stale++
			cell.State = "STALE"
		}
		out.Cells = append(out.Cells, cell)
	}
	if out.Expected > 0 {
		out.CoveragePct = float64(out.Fresh) / float64(out.Expected) * 100
	}
	if out.Fresh > 0 {
		out.DirectionalPct = v165FloatPtr(float64(adv) / float64(out.Fresh) * 100)
	}
	out.Detail = fmt.Sprintf("%d/%d canonical constituents current; stale/missing cells do not vote. %s", out.Fresh, out.Expected, basis)
	return out
}

func v168HeatDefs(symbols []string) []heatMapSymbolDef {
	seen := map[string]bool{}
	out := []heatMapSymbolDef{}
	for _, raw := range symbols {
		sym := normalizeSymbol(raw)
		if sym == "" || sym == "VIX" || seen[sym] {
			continue
		}
		seen[sym] = true
		c := classificationForSymbol(sym)
		name, sector := sym, c.Sector
		if sector == "" {
			sector = "Unclassified"
		}
		out = append(out, heatMapSymbolDef{Symbol: sym, Name: name, Sector: sector})
	}
	sort.Slice(out, func(i, j int) bool { return out[i].Symbol < out[j].Symbol })
	return out
}

func v165HeatMapForState(st AppState, quotes map[string]Quote, now time.Time) MarketSectorHeatMap {
	sector := v168BuildHeatMode("sector", "Sector Benchmarks", "11 SPDR sector benchmark ETFs", "Explicit sector benchmark universe.", v165SectorHeatUniverse, quotes, now)
	watchDefs := v168HeatDefs(analysisSymbolsFromState(st))
	watch := v168BuildHeatMode("watchlist", "Watchlist Heat", "Active desk + discovery/selected canonical symbols", "Reuses the existing application/watchlist quote universe; no separate fetch path.", watchDefs, quotes, now)
	broadSymbols := append([]string{}, marketIntelligenceBreadthUniverse...)
	for sym := range canonicalSymbolClassifications {
		if _, ok := quotes[sym]; ok {
			broadSymbols = append(broadSymbols, sym)
		}
	}
	broad := v168BuildHeatMode("broad", "Broad Market / S&P Proxy", "SPY + canonical broad/sector ETFs + classified covered equities", "This is a coverage-truthful S&P/broad-market proxy over symbols already in the canonical quote store, not a claim of full 500-constituent coverage.", v168HeatDefs(broadSymbols), quotes, now)
	return MarketSectorHeatMap{Universe: sector.Universe, Expected: sector.Expected, Fresh: sector.Fresh, Stale: sector.Stale, Missing: sector.Missing, CoveragePct: sector.CoveragePct, DirectionalPct: sector.DirectionalPct, Cells: sector.Cells, UpdatedAt: sector.UpdatedAt, Detail: sector.Detail, Modes: map[string]HeatMapMode{"sector": sector, "watchlist": watch, "broad": broad}}
}

func v165OptionsSentiment(o OptionsContext) (SentimentComponent, bool) {
	if o.PutCallVolume <= 0 || strings.EqualFold(o.Bias, "INCOMPLETE") || o.UpdatedAt <= 0 {
		return SentimentComponent{}, false
	}
	score := v165Clamp((1-o.PutCallVolume)/(1+o.PutCallVolume)*100, -100, 100)
	confidence := 70
	if strings.Contains(strings.ToUpper(o.State), "INDICATIVE") || strings.Contains(strings.ToUpper(o.State), "DELAYED") {
		confidence = 45
	}
	return SentimentComponent{Name: "Options Positioning · " + o.Symbol, State: o.Bias, Score: v165FloatPtr(score), Weight: 1, Confidence: confidence, Source: o.Provider + " / " + o.Feed, UpdatedAt: o.UpdatedAt, Detail: fmt.Sprintf("Put/call volume %.2f; options are context only.", o.PutCallVolume)}, true
}

func v165Sentiment(mi MarketIntelligenceSnapshot, heat MarketSectorHeatMap, quotes map[string]Quote, options map[string]OptionsContext, now time.Time) SentimentCompositeState {
	out := SentimentCompositeState{State: "UNAVAILABLE", ComponentsExpected: 5, Components: []SentimentComponent{}, Missing: []string{}}
	add := func(c SentimentComponent) {
		out.Components = append(out.Components, c)
		if c.UpdatedAt > out.UpdatedAt {
			out.UpdatedAt = c.UpdatedAt
		}
	}
	if mi.Breadth.ParticipationPct != nil && mi.Breadth.Fresh > 0 {
		score := v165Clamp((*mi.Breadth.ParticipationPct-50)*2, -100, 100)
		add(SentimentComponent{Name: "Tracked Breadth", State: mi.Breadth.State, Score: v165FloatPtr(score), Weight: 1.25, Confidence: int(v165Clamp(mi.Breadth.CoveragePct, 0, 100)), Source: "DE.PULSE Market Intelligence", UpdatedAt: mi.Breadth.UpdatedAt, Detail: mi.Breadth.Detail})
	} else {
		out.Missing = append(out.Missing, "Tracked breadth unavailable")
	}
	if heat.DirectionalPct != nil && heat.Fresh >= 6 {
		score := v165Clamp((*heat.DirectionalPct-50)*2, -100, 100)
		add(SentimentComponent{Name: "Sector Participation", State: "SECTOR BREADTH", Score: v165FloatPtr(score), Weight: 1.25, Confidence: int(v165Clamp(heat.CoveragePct, 0, 100)), Source: heat.Universe, UpdatedAt: heat.UpdatedAt, Detail: heat.Detail})
	} else {
		out.Missing = append(out.Missing, "Sector heat-map coverage insufficient")
	}
	marketScores := []float64{}
	marketStamp := int64(0)
	marketSources := []string{}
	for _, sym := range []string{"SPY", "QQQ"} {
		if q, ok := quotes[sym]; ok && q.Price > 0 {
			stamp := v165QuoteStamp(q)
			f, _, _, _, _ := freshnessState("Quotes", providerFromQuoteSource(q.Source), marketSessionET(now), stamp, "", now.UnixMilli())
			if f != "STALE" && f != "UNAVAILABLE" && f != "ERROR" {
				marketScores = append(marketScores, v165Clamp(q.ChangePercent*20, -100, 100))
				if stamp > marketStamp {
					marketStamp = stamp
				}
				marketSources = append(marketSources, sym+" "+q.Source)
			}
		}
	}
	if len(marketScores) > 0 {
		sum := 0.0
		for _, x := range marketScores {
			sum += x
		}
		score := sum / float64(len(marketScores))
		add(SentimentComponent{Name: "SPY / QQQ Tape", State: "OBSERVED", Score: v165FloatPtr(score), Weight: 1.5, Confidence: 80, Source: strings.Join(marketSources, " · "), UpdatedAt: marketStamp, Detail: "Observed broad-index price change; no text sentiment inference."})
	} else {
		out.Missing = append(out.Missing, "SPY/QQQ current tape unavailable")
	}
	optAdded := 0
	for _, sym := range []string{"SPY", "QQQ"} {
		if c, ok := v165OptionsSentiment(options[sym]); ok {
			add(c)
			optAdded++
			break
		}
	}
	if optAdded == 0 {
		out.Missing = append(out.Missing, "Broad-index options positioning unavailable")
	}
	globalState := strings.ToUpper(strings.TrimSpace(mi.Tradeability.State))
	if globalState != "" && globalState != "DATA DEGRADED" {
		score := 0.0
		switch globalState {
		case "TRADE NORMALLY":
			score = 35
		case "SELECTIVE":
			score = 10
		case "REDUCE SIZE":
			score = -20
		case "WAIT":
			score = -45
		}
		add(SentimentComponent{Name: "Market Tradeability Context", State: globalState, Score: v165FloatPtr(score), Weight: .75, Confidence: 65, Source: "DE.PULSE Market Intelligence", UpdatedAt: mi.Tradeability.UpdatedAt, Detail: "Context state contributes conservatively; it is not a directional trading signal."})
	} else {
		out.Missing = append(out.Missing, "Market Tradeability context degraded")
	}
	out.ComponentsUsed = len(out.Components)
	weighted, weights, conf := 0.0, 0.0, 0.0
	for _, c := range out.Components {
		if c.Score != nil {
			weighted += *c.Score * c.Weight
			weights += c.Weight
			conf += float64(c.Confidence) * c.Weight
		}
	}
	if out.ComponentsUsed < 3 || weights == 0 {
		out.Confidence = int(v165Clamp(conf/math.Max(weights, 1), 0, 45))
		out.Detail = fmt.Sprintf("Only %d/%d required sentiment components are usable; no composite direction is published.", out.ComponentsUsed, out.ComponentsExpected)
		return out
	}
	score := v165Clamp(weighted/weights, -100, 100)
	out.Score = v165FloatPtr(score)
	coverage := float64(out.ComponentsUsed) / float64(out.ComponentsExpected)
	out.Confidence = int(v165Clamp(conf/weights*coverage, 0, 100))
	switch {
	case score >= 25:
		out.State = "BULLISH"
	case score <= -25:
		out.State = "BEARISH"
	case math.Abs(score) <= 10:
		out.State = "NEUTRAL"
	default:
		out.State = "MIXED"
	}
	out.Detail = fmt.Sprintf("Transparent composite from %d/%d available market components; missing components reduce confidence and never become neutral zeros.", out.ComponentsUsed, out.ComponentsExpected)
	return out
}

func v165GEX(options map[string]OptionsContext) map[string]GEXContextState {
	out := map[string]GEXContextState{}
	for sym, o := range options {
		total := o.CallContracts + o.PutContracts
		s := GEXContextState{Symbol: normalizeSymbol(sym), State: "UNAVAILABLE", CoveragePct: o.GammaOICoveragePct, GammaOIContracts: o.GammaOIContracts, GammaContracts: o.GammaContracts, TotalContracts: total, OpenInterestDate: o.OpenInterestDate, Feed: o.Feed, Source: o.Provider, UpdatedAt: o.UpdatedAt, DeterministicImpact: "NONE", Limitations: append([]string{}, o.Limitations...)}
		if o.GEXState == "AVAILABLE" {
			s.State, s.Quality, s.NetGEX, s.CallGEX, s.PutGEX = o.GEXState, o.GEXQuality, o.NetGEX, o.CallGEX, o.PutGEX
			s.UnderlyingPrice, s.MajorGammaStrikes, s.GammaZones, s.ExpirationGEX, s.GammaFlip, s.GammaFlipMethod = o.UnderlyingPrice, clone(o.MajorGammaStrikes), clone(o.GammaZones), clone(o.ExpirationGEX), o.GammaFlip, o.GammaFlipMethod
			s.Detail = fmt.Sprintf("Estimated structural gamma exposure proxy from %d/%d gamma-bearing contracts with matched open interest (%.0f%% coverage); %d major strikes and %d expiration summaries retained; full chain contains %d contracts.", o.GammaOIContracts, o.GammaContracts, o.GammaOICoveragePct, len(o.MajorGammaStrikes), len(o.ExpirationGEX), total)
		} else {
			s.Detail = "GEX withheld because gamma/open-interest completeness or recency is insufficient."
		}
		out[s.Symbol] = s
	}
	return out
}

func v165MetricPtr(metrics map[string]MacroMetric, key string) *MacroMetric {
	if m, ok := metrics[key]; ok && m.UpdatedAt > 0 {
		x := m
		return &x
	}
	return nil
}
func v165QuotePtr(quotes map[string]Quote, key string) *Quote {
	if q, ok := quotes[key]; ok && q.Price > 0 {
		x := q
		return &x
	}
	return nil
}

func v169EnergyTrend(m *MacroMetric) string {
	if m == nil {
		return "UNAVAILABLE"
	}
	change := m.Change5D
	if math.Abs(change) < 0.01 {
		change = m.Change20D
	}
	switch {
	case change >= 3:
		return "RISING"
	case change <= -3:
		return "FALLING"
	default:
		return "STABLE / MIXED"
	}
}

func v169EnergyInstrument(name, contract string, spot *MacroMetric, quote *Quote) EnergyInstrumentTruth {
	out := EnergyInstrumentTruth{Name: name, Contract: contract, State: "UNAVAILABLE", ReferenceType: "OFFICIAL SPOT REFERENCE", Trend: "UNAVAILABLE", Delayed: true, ContinuousContract: false, RollAdjusted: false, Semantics: contract + " continuous-futures context is not subscribed from the current free routed provider set; official spot reference is used when available and is never mislabeled as a futures contract."}
	if quote != nil && quote.Price > 0 {
		v := quote.Price
		c5 := quote.ChangePercent
		out.State = "AVAILABLE"
		out.ReferenceType = "FUTURES/CONTEXT QUOTE"
		out.Value = &v
		out.Change5D = &c5
		out.Trend = "CURRENT QUOTE"
		out.Source = quote.Source
		out.Delayed = strings.Contains(strings.ToUpper(quote.DataState), "DELAY") || strings.Contains(strings.ToUpper(quote.FeedType), "DELAY")
		out.UpdatedAt = v165QuoteStamp(*quote)
		out.Semantics = contract + " context is provider-labelled quote evidence. DE.PULSE does not infer roll adjustment or continuous-contract construction unless the provider explicitly supplies those semantics. No futures execution."
		return out
	}
	if spot != nil {
		v, c5, c20 := spot.Value, spot.Change5D, spot.Change20D
		out.State = "AVAILABLE · SPOT REFERENCE"
		out.Value, out.Change5D, out.Change20D = &v, &c5, &c20
		out.Trend = v169EnergyTrend(spot)
		out.Source = spot.Source
		out.UpdatedAt = spot.UpdatedAt
		out.Semantics = contract + " continuous futures are unavailable from the configured free provider set; " + name + " official spot/reference data is used for contextual trend only. Roll/curve behavior is therefore not fabricated."
	}
	return out
}

func v165OilEnergy(metrics map[string]MacroMetric, quotes map[string]Quote, now time.Time) OilEnergyContextState {
	out := OilEnergyContextState{State: "UNAVAILABLE", WTI: v165MetricPtr(metrics, "WTI_OFFICIAL"), Brent: v165MetricPtr(metrics, "BRENT_OFFICIAL"), CrudeStocks: v165MetricPtr(metrics, "CRUDE_STOCKS"), Production: v165MetricPtr(metrics, "CRUDE_PRODUCTION"), RefineryUtil: v165MetricPtr(metrics, "REFINERY_UTIL"), USO: v165QuotePtr(quotes, "USO"), XLE: v165QuotePtr(quotes, "XLE"), Sources: []string{}, Limitations: []string{"USO is a tradable futures-based ETF, not WTI spot; basis, futures-curve and roll effects can diverge from crude spot prices.", "EIA publication timing differs by dataset; UpdatedAt/Period remain the source truth.", "CL/BZ continuous-contract roll/curve semantics are never inferred when the configured provider does not supply them."}, DeterministicImpact: "NONE"}
	out.WTIContext = v169EnergyInstrument("WTI", "CL", out.WTI, v165QuotePtr(quotes, "CL"))
	out.BrentContext = v169EnergyInstrument("Brent", "BZ", out.Brent, v165QuotePtr(quotes, "BZ"))
	for _, m := range []*MacroMetric{out.WTI, out.Brent, out.CrudeStocks, out.Production, out.RefineryUtil} {
		if m != nil {
			out.ComponentsUsed++
			out.Sources = appendUniqueString(out.Sources, m.Source)
			if m.UpdatedAt > out.UpdatedAt {
				out.UpdatedAt = m.UpdatedAt
			}
		}
	}
	for _, q := range []*Quote{out.USO, out.XLE} {
		if q != nil {
			out.ComponentsUsed++
			out.Sources = appendUniqueString(out.Sources, q.Source)
			if st := v165QuoteStamp(*q); st > out.UpdatedAt {
				out.UpdatedAt = st
			}
		}
	}
	if out.WTI != nil && out.Brent != nil {
		spread := out.Brent.Value - out.WTI.Value
		out.BrentWTISpread = &spread
		switch {
		case spread >= 8:
			out.SpreadState = "WIDE BRENT PREMIUM"
		case spread <= 1:
			out.SpreadState = "NARROW"
		default:
			out.SpreadState = "NORMAL RANGE"
		}
	}
	wtiTrend := out.WTIContext.Trend
	xleChange := 0.0
	if out.XLE != nil {
		xleChange = out.XLE.ChangePercent
	}
	switch {
	case wtiTrend == "RISING" && xleChange > 0.5:
		out.EnergySectorState = "SUPPORTIVE"
		out.USMarketRelevance = "Rising crude plus positive XLE confirms an energy-sector tailwind; also monitor inflation/rates sensitivity for the broader U.S. market."
	case wtiTrend == "FALLING" && xleChange < -0.5:
		out.EnergySectorState = "HEADWIND"
		out.USMarketRelevance = "Falling crude plus weak XLE is an energy-sector headwind and may ease inflation pressure; context only."
	case out.ComponentsUsed > 0:
		out.EnergySectorState = "MIXED"
		out.USMarketRelevance = "Oil, inventories and U.S. energy-equity evidence are mixed; no broad market conclusion is forced."
	default:
		out.EnergySectorState = "UNAVAILABLE"
		out.USMarketRelevance = "Insufficient energy evidence for U.S.-market interpretation."
	}
	switch {
	case out.WTI != nil && out.Brent != nil && out.USO != nil:
		out.State = "AVAILABLE"
	case out.ComponentsUsed >= 2:
		out.State = "PARTIAL"
	}
	if out.State == "UNAVAILABLE" {
		out.Detail = "Official EIA WTI/Brent and current USO/XLE evidence are unavailable."
	} else {
		spreadText := "spread unavailable"
		if out.BrentWTISpread != nil {
			spreadText = fmt.Sprintf("Brent-WTI %.2f USD/bbl · %s", *out.BrentWTISpread, out.SpreadState)
		}
		out.Detail = fmt.Sprintf("%d energy components available · WTI %s · Brent %s · %s. %s", out.ComponentsUsed, out.WTIContext.Trend, out.BrentContext.Trend, spreadText, out.USMarketRelevance)
	}
	return out
}

func buildContextAlternativeIntelligenceSnapshot(st AppState, quotes map[string]Quote, options map[string]OptionsContext, metrics map[string]MacroMetric, mi MarketIntelligenceSnapshot, news []NewsItem, filings []FilingItem, now time.Time) ContextAlternativeIntelligenceSnapshot {
	heat := v165HeatMapForState(st, quotes, now)
	return ContextAlternativeIntelligenceSnapshot{Sentiment: v165Sentiment(mi, heat, quotes, options, now), HeatMap: heat, GEX: v165GEX(options), Community: buildCommunityEvidenceFusion(st.CommunityEvidence, news, filings, now), OilEnergy: v165OilEnergy(metrics, quotes, now), UpdatedAt: now.UnixMilli()}
}
