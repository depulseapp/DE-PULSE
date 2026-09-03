package main

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math"
	"sort"
	"strings"
	"time"
)

const rapidMovePolicyVersion = "rapid-move-v1.1.0"
const rapidMoveLearningPolicyVersion = "rapid-move-learning-v1.0.0"

var rapidMoveWindows = []struct {
	Label     string
	Duration  time.Duration
	BaseLimit float64
}{
	{"15s", 15 * time.Second, 3.0},
	{"30s", 30 * time.Second, 4.0},
	{"60s", 60 * time.Second, 5.0},
	{"2m", 2 * time.Minute, 6.5},
	{"5m", 5 * time.Minute, 8.0},
}

type RapidMoveWindowMetric struct {
	Window       string  `json:"window"`
	WindowMs     int64   `json:"windowMs"`
	StartAt      int64   `json:"startAt"`
	StartPrice   float64 `json:"startPrice"`
	ReturnPct    float64 `json:"returnPct"`
	ThresholdPct float64 `json:"thresholdPct"`
	TriggerRatio float64 `json:"triggerRatio"`
}

type RapidMoveCoverage struct {
	State              string `json:"state"`
	LiveSymbols        int    `json:"liveSymbols"`
	SubscribedSymbols  int    `json:"subscribedSymbols"`
	WindowReadySymbols int    `json:"windowReadySymbols"`
	BroadSeedProvider  string `json:"broadSeedProvider,omitempty"`
	BroadSeedState     string `json:"broadSeedState,omitempty"`
	Detail             string `json:"detail"`
	UpdatedAt          int64  `json:"updatedAt"`
}

type RapidMovePolicyGovernance struct {
	DetectionPolicyVersion string   `json:"detectionPolicyVersion"`
	DetectionStage         string   `json:"detectionStage"`
	LearningPolicyVersion  string   `json:"learningPolicyVersion"`
	LearningStage          string   `json:"learningStage"`
	PromotionPath          []string `json:"promotionPath"`
	AutoPromotion          bool     `json:"autoPromotion"`
	ProtectedFormulaImpact string   `json:"protectedFormulaImpact"`
}

func rapidMovePolicyGovernance() RapidMovePolicyGovernance {
	return RapidMovePolicyGovernance{
		DetectionPolicyVersion: rapidMovePolicyVersion, DetectionStage: "PRODUCTION",
		LearningPolicyVersion: rapidMoveLearningPolicyVersion, LearningStage: "SHADOW",
		PromotionPath: []string{"SHADOW", "VALIDATED", "APPROVED", "PRODUCTION"}, AutoPromotion: false,
		ProtectedFormulaImpact: "NONE · Day/Swing/Long deterministic formulas remain unchanged",
	}
}

type RapidMoveEvent struct {
	ID                    string   `json:"id"`
	TraceID               string   `json:"traceId"`
	Symbol                string   `json:"symbol"`
	Direction             string   `json:"direction"`
	State                 string   `json:"state"`          // EARLY / VALIDATING / CONFIRMED / EXTENDED / FADING / RESOLVED
	Classification        string   `json:"classification"` // RAPID_MOVE / MARKET_SHOCK
	Severity              string   `json:"severity"`
	Alerted               bool     `json:"alerted"`
	ShadowWouldAlert      bool     `json:"shadowWouldAlert,omitempty"`
	PolicyVersion         string   `json:"policyVersion"`
	LearningPolicyVersion string   `json:"learningPolicyVersion"`
	AdaptiveStage         string   `json:"adaptiveStage"`
	Session               string   `json:"session"`
	Window                string   `json:"window"`
	WindowMs              int64    `json:"windowMs"`
	WindowStartAt         int64    `json:"windowStartAt"`
	EventProviderAt       int64    `json:"eventProviderAt,omitempty"`
	ReceivedAt            int64    `json:"receivedAt,omitempty"`
	DetectedAt            int64    `json:"detectedAt"`
	UpdatedAt             int64    `json:"updatedAt"`
	StartPrice            float64  `json:"startPrice"`
	Price                 float64  `json:"price"`
	MovePct               float64  `json:"movePct"`
	ThresholdPct          float64  `json:"thresholdPct"`
	TriggerRatio          float64  `json:"triggerRatio"`
	MaterialityScore      float64  `json:"materialityScore"`
	RelativeVolume        float64  `json:"relativeVolume,omitempty"`
	VolumeState           string   `json:"volumeState"`
	SpreadPct             float64  `json:"spreadPct,omitempty"`
	LiquidityState        string   `json:"liquidityState"`
	SourceAgreement       string   `json:"sourceAgreement"`
	SourceAgreementDetail string   `json:"sourceAgreementDetail,omitempty"`
	CatalystState         string   `json:"catalystState"`
	CatalystSummary       string   `json:"catalystSummary,omitempty"`
	CatalystPublishedAt   int64    `json:"catalystPublishedAt,omitempty"`
	IndependentSources    int      `json:"independentSources,omitempty"`
	MarketContext         string   `json:"marketContext,omitempty"`
	MechanicalRisk        string   `json:"mechanicalRisk,omitempty"`
	Reasons               []string `json:"reasons,omitempty"`
	Limitations           []string `json:"limitations,omitempty"`
	MFE                   float64  `json:"mfe,omitempty"`
	MAE                   float64  `json:"mae,omitempty"`
	Outcome1mPct          *float64 `json:"outcome1mPct,omitempty"`
	Outcome5mPct          *float64 `json:"outcome5mPct,omitempty"`
	Outcome20mPct         *float64 `json:"outcome20mPct,omitempty"`
	OutcomeState          string   `json:"outcomeState,omitempty"`
}

type RapidMoveScorecard struct {
	PolicyVersion            string `json:"policyVersion"`
	Observations             int64  `json:"observations"`
	CandidateEvaluations     int64  `json:"candidateEvaluations"`
	ProductionAlerts         int64  `json:"productionAlerts"`
	ShadowWouldAlert         int64  `json:"shadowWouldAlert"`
	MechanicalSuppressed     int64  `json:"mechanicalSuppressed"`
	LiquiditySuppressed      int64  `json:"liquiditySuppressed"`
	SourceConflictSuppressed int64  `json:"sourceConflictSuppressed"`
	DuplicateUpdates         int64  `json:"duplicateUpdates"`
	ConfirmedCatalysts       int64  `json:"confirmedCatalysts"`
	UnexplainedMoves         int64  `json:"unexplainedMoves"`
	OutcomesResolved         int64  `json:"outcomesResolved"`
	MarketShockAlerts        int64  `json:"marketShockAlerts"`
	HysteresisRetained       int64  `json:"hysteresisRetained"`
	ContinuedOutcomes        int64  `json:"continuedOutcomes"`
	ReversedOutcomes         int64  `json:"reversedOutcomes"`
	FadedMixedOutcomes       int64  `json:"fadedMixedOutcomes"`
	LastEventAt              int64  `json:"lastEventAt,omitempty"`
}

type RapidMoveState struct {
	Status        string                    `json:"status"`
	PolicyVersion string                    `json:"policyVersion"`
	AdaptiveStage string                    `json:"adaptiveStage"`
	Coverage      RapidMoveCoverage         `json:"coverage"`
	Active        []RapidMoveEvent          `json:"active"`
	Recent        []RapidMoveEvent          `json:"recent"`
	Scorecard     RapidMoveScorecard        `json:"scorecard"`
	Governance    RapidMovePolicyGovernance `json:"governance"`
	Detail        string                    `json:"detail"`
	UpdatedAt     int64                     `json:"updatedAt"`
}

func rapidMoveSpreadPct(q Quote) float64 {
	if q.Bid <= 0 || q.Ask < q.Bid {
		return 0
	}
	mid := (q.Bid + q.Ask) / 2
	if mid <= 0 {
		return 0
	}
	return (q.Ask - q.Bid) / mid * 100
}

func rapidMoveHistoryMetric(history []HistoryPoint, nowMs int64, current float64, window time.Duration) (RapidMoveWindowMetric, bool) {
	if current <= 0 || len(history) < 2 || nowMs <= 0 {
		return RapidMoveWindowMetric{}, false
	}
	target := nowMs - window.Milliseconds()
	tolerance := int64(math.Max(5000, float64(window.Milliseconds())*.35))
	bestIdx := -1
	bestDistance := int64(1<<62 - 1)
	for i := len(history) - 1; i >= 0; i-- {
		ts := normalizeObservationMs(history[i].T)
		if ts <= 0 || history[i].P <= 0 {
			continue
		}
		d := ts - target
		if d < 0 {
			d = -d
		}
		if d < bestDistance {
			bestDistance = d
			bestIdx = i
		}
		if ts < target-tolerance {
			break
		}
	}
	if bestIdx < 0 || bestDistance > tolerance {
		return RapidMoveWindowMetric{}, false
	}
	p := history[bestIdx]
	ts := normalizeObservationMs(p.T)
	ret := (current/p.P - 1) * 100
	return RapidMoveWindowMetric{WindowMs: window.Milliseconds(), StartAt: ts, StartPrice: p.P, ReturnPct: ret}, true
}

func rapidMoveThreshold(base float64, symbol string, q Quote, daily []Bar, now time.Time) float64 {
	factor := 1.0
	atr := atrPercentFromBars(daily)
	if atr > 0 {
		factor *= math.Max(.82, math.Min(1.45, .82+atr/12))
	}
	session := marketSessionET(now)
	switch session {
	case "pre-market", "after-hours":
		factor *= 1.12
	case "overnight":
		factor *= 1.25
	}
	et := now.In(easternLocation())
	mins := et.Hour()*60 + et.Minute()
	if session == "regular" && mins >= 9*60+30 && mins < 9*60+40 {
		factor *= 1.15
	}
	if q.Price > 0 && q.Price < 5 {
		factor *= 1.25
	}
	return base * factor
}

func rapidMoveSourceAgreement(observations map[string]Quote, nowMs, maxSkewMs int64) (string, string) {
	type obs struct {
		provider string
		price    float64
	}
	if maxSkewMs <= 0 {
		maxSkewMs = 5_000
	}
	rows := []obs{}
	excludedLagging := 0
	for provider, q := range observations {
		if q.Price <= 0 {
			continue
		}
		_, _, current := quoteEvidenceAges(q, nowMs)
		if !current {
			continue
		}
		stamp := normalizeObservationMs(q.ProviderTimestamp)
		if stamp <= 0 {
			stamp = normalizeObservationMs(q.UpdatedAt)
		}
		if stamp <= 0 || nowMs-stamp > maxSkewMs || stamp-nowMs > 2_000 {
			excludedLagging++
			continue
		}
		rows = append(rows, obs{provider, q.Price})
	}
	if len(rows) < 2 {
		detail := "Only one contemporaneous provider observation is available."
		if excludedLagging > 0 {
			detail = fmt.Sprintf("Only one contemporaneous provider observation is available; %d lagging/non-contemporaneous observation(s) excluded.", excludedLagging)
		}
		return "SINGLE_SOURCE", detail
	}
	minP, maxP := rows[0].price, rows[0].price
	for _, r := range rows[1:] {
		if r.price < minP {
			minP = r.price
		}
		if r.price > maxP {
			maxP = r.price
		}
	}
	if minP <= 0 {
		return "UNKNOWN", "Provider price comparison is unavailable."
	}
	diff := (maxP/minP - 1) * 100
	names := make([]string, 0, len(rows))
	for _, r := range rows {
		names = append(names, r.provider)
	}
	sort.Strings(names)
	if diff >= 2.0 {
		return "CONFLICT", fmt.Sprintf("Current providers disagree by %.2f%% (%s).", diff, strings.Join(names, ", "))
	}
	if diff >= .75 {
		return "MIXED", fmt.Sprintf("Provider spread is %.2f%% (%s).", diff, strings.Join(names, ", "))
	}
	return "AGREED", fmt.Sprintf("Current providers agree within %.2f%% (%s).", diff, strings.Join(names, ", "))
}

func rapidMoveMechanicalRisk(symbol string, actions []CorporateAction, now time.Time) string {
	symbol = normalizeSymbol(symbol)
	day := now.In(easternLocation()).Format("2006-01-02")
	for _, a := range actions {
		if normalizeSymbol(a.Symbol) != symbol {
			continue
		}
		typ := strings.ToLower(strings.TrimSpace(a.Type))
		if !strings.Contains(typ, "split") && !strings.Contains(typ, "merger") && !strings.Contains(typ, "name") && !strings.Contains(typ, "spinoff") {
			continue
		}
		if a.ExDate == day || a.ProcessDate == day || (strings.EqualFold(a.Status, "EFFECTIVE") && a.UpdatedAt > now.Add(-24*time.Hour).UnixMilli()) {
			return fmt.Sprintf("%s corporate action may mechanically explain the price discontinuity", strings.ReplaceAll(strings.ToUpper(typ), "_", " "))
		}
	}
	return ""
}

func rapidMoveCatalyst(symbol string, news []NewsItem, filings []FilingItem, earnings []EarningsItem, reactions map[string]CatalystReactionState, now time.Time) (state, summary string, publishedAt int64, independent int) {
	symbol = normalizeSymbol(symbol)
	nowMs := now.UnixMilli()
	if c, ok := reactions[symbol]; ok && c.TriggerAt > 0 && c.TriggerAt <= nowMs+2_000 && c.CompletedAt == 0 {
		return "CONFIRMED", defaultString(c.Trigger, c.TriggerType), c.TriggerAt, 1
	}
	clusters := buildEventNewsIntelligence(news, now)
	for _, n := range clusters {
		found := false
		for _, s := range n.Symbols {
			if normalizeSymbol(s) == symbol {
				found = true
				break
			}
		}
		if !found || n.PublishedAt <= 0 || n.PublishedAt > now.Unix()+2 || now.Unix()-n.PublishedAt > 20*60 {
			continue
		}
		sources := map[string]bool{}
		if strings.TrimSpace(n.Source) != "" {
			sources[strings.ToLower(strings.TrimSpace(n.Source))] = true
		}
		for _, s := range n.SupportingSources {
			if strings.TrimSpace(s) != "" {
				sources[strings.ToLower(strings.TrimSpace(s))] = true
			}
		}
		material := strings.EqualFold(n.Materiality, "HIGH") || materialText(n.Headline+" "+n.Summary)
		if material {
			return "CONFIRMED", n.Headline, n.PublishedAt * 1000, len(sources)
		}
		return "POSSIBLE", n.Headline, n.PublishedAt * 1000, len(sources)
	}
	for _, f := range filings {
		if normalizeSymbol(f.Symbol) == symbol && isRecentDate(f.FiledAt, 1) && materialSECFilingForTradingRisk(f) {
			return "POSSIBLE", strings.TrimSpace(f.Form + " · " + defaultString(f.Meaning, f.Description)), notificationTimestamp(f.FiledAt, now.UnixMilli()), 1
		}
	}
	day := now.In(easternLocation()).Format("2006-01-02")
	for _, er := range earnings {
		if normalizeSymbol(er.Symbol) == symbol && er.Date == day {
			return "POSSIBLE", "Earnings-day catalyst context", now.UnixMilli(), 1
		}
	}
	return "UNEXPLAINED", "No fresh material catalyst is confirmed yet.", 0, 0
}

func rapidMoveMarketContext(history map[string][]HistoryPoint, nowMs int64, window time.Duration, direction string) string {
	best := 0.0
	name := ""
	for _, sym := range []string{"SPY", "QQQ"} {
		h := history[sym]
		if len(h) == 0 {
			continue
		}
		last := h[len(h)-1]
		lastAt := normalizeObservationMs(last.T)
		tolerance := int64(math.Max(5_000, float64(window.Milliseconds())*.35))
		if lastAt <= 0 || nowMs-lastAt > tolerance || lastAt-nowMs > 2_000 {
			continue
		}
		cur := last.P
		m, ok := rapidMoveHistoryMetric(h, nowMs, cur, window)
		if !ok {
			continue
		}
		if math.Abs(m.ReturnPct) > math.Abs(best) {
			best, name = m.ReturnPct, sym
		}
	}
	if name == "" {
		return "MARKET_CONTEXT_PENDING"
	}
	same := (best > 0 && direction == "UP") || (best < 0 && direction == "DOWN")
	if math.Abs(best) >= .8 && same {
		return fmt.Sprintf("MARKET_WIDE · %s %.2f%% in comparable window", name, best)
	}
	return fmt.Sprintf("IDIOSYNCRATIC · %s %.2f%% in comparable window", name, best)
}

func rapidMoveRelativeVolume(q Quote, daily []Bar) (float64, string) {
	// Finnhub websocket v is trade size rather than trustworthy cumulative-day volume.
	if strings.Contains(strings.ToLower(q.Source), "finnhub-websocket") {
		return 0, "PENDING"
	}
	rv := relativeVolumeFromDaily(q, daily)
	if rv <= 0 {
		return 0, "PENDING"
	}
	if rv >= 2 {
		return rv, "ABNORMAL"
	}
	if rv >= 1.25 {
		return rv, "ELEVATED"
	}
	return rv, "NORMAL"
}

func rapidMoveClassification(marketContext string) string {
	if strings.HasPrefix(strings.TrimSpace(marketContext), "MARKET_WIDE") {
		return "MARKET_SHOCK"
	}
	return "RAPID_MOVE"
}

func rapidMoveApplyHysteresis(existing, candidate RapidMoveEvent) (RapidMoveEvent, bool) {
	if !existing.Alerted || existing.State == "RESOLVED" || existing.Direction != candidate.Direction {
		return candidate, false
	}
	before := candidate.State
	switch existing.State {
	case "FADING":
		// A fading event may re-accelerate only when it again satisfies a production state.
		if candidate.State != "CONFIRMED" && candidate.State != "EXTENDED" {
			candidate.State = "FADING"
		}
	case "EXTENDED":
		// Small retracements do not manufacture a downgrade; the outcome path owns FADING.
		if candidate.State == "EARLY" || candidate.State == "VALIDATING" || candidate.State == "CONFIRMED" {
			candidate.State = "EXTENDED"
		}
	case "CONFIRMED":
		if candidate.State == "EARLY" || candidate.State == "VALIDATING" {
			candidate.State = "CONFIRMED"
		}
	case "VALIDATING":
		if candidate.State == "EARLY" {
			candidate.State = "VALIDATING"
		}
	}
	return candidate, candidate.State != before
}

func rapidMoveMateriality(metric RapidMoveWindowMetric, q Quote, relativeVolume float64, volumeState, agreement, catalyst, marketContext, mechanical string) (float64, []string) {
	score := math.Min(55, metric.TriggerRatio*42)
	reasons := []string{fmt.Sprintf("%+.2f%% in %s vs %.2f%% context threshold", metric.ReturnPct, metric.Window, metric.ThresholdPct)}
	spread := rapidMoveSpreadPct(q)
	if spread > 0 && spread <= .35 {
		score += 12
		reasons = append(reasons, fmt.Sprintf("healthy %.2f%% spread", spread))
	} else if spread > .75 {
		score -= 18
		reasons = append(reasons, fmt.Sprintf("wide %.2f%% spread", spread))
	}
	if relativeVolume >= 2 {
		score += 12
		reasons = append(reasons, fmt.Sprintf("%.2fx relative volume", relativeVolume))
	} else if volumeState == "PENDING" {
		reasons = append(reasons, "volume confirmation pending")
	}
	switch agreement {
	case "AGREED":
		score += 10
		reasons = append(reasons, "multi-source price agreement")
	case "MIXED":
		score += 3
		reasons = append(reasons, "provider prices mixed")
	case "CONFLICT":
		score -= 35
		reasons = append(reasons, "source disagreement requires validation")
	}
	switch catalyst {
	case "CONFIRMED":
		score += 16
		reasons = append(reasons, "fresh material catalyst confirmed")
	case "POSSIBLE":
		score += 7
		reasons = append(reasons, "possible fresh catalyst")
	case "UNEXPLAINED":
		reasons = append(reasons, "move is unexplained; catalyst search remains active")
	}
	if strings.HasPrefix(marketContext, "IDIOSYNCRATIC") {
		score += 5
		reasons = append(reasons, "move is not explained by SPY/QQQ")
	} else if strings.HasPrefix(marketContext, "MARKET_WIDE") {
		score -= 4
		reasons = append(reasons, "market-wide move reduces company-specificity")
	}
	if mechanical != "" {
		score -= 70
		reasons = append(reasons, mechanical)
	}
	if score < 0 {
		score = 0
	}
	if score > 100 {
		score = 100
	}
	return score, uniqueStrings(reasons)
}

func (e *Engine) rapidMoveCoverageLocked(now time.Time) RapidMoveCoverage {
	subscribed := map[string]bool{}
	for s := range e.subscribedSymbols {
		subscribed[normalizeSymbol(s)] = true
	}
	for s := range e.alpacaSubscribedSymbols {
		subscribed[normalizeSymbol(s)] = true
	}
	windowReady := 0
	nowMs := now.UnixMilli()
	for symbol := range subscribed {
		q := e.quotes[symbol]
		if q.Price <= 0 || len(e.history[symbol]) < 2 {
			continue
		}
		_, _, current := quoteEvidenceAges(q, nowMs)
		if current {
			windowReady++
		}
	}
	broadProvider := ""
	broadState := "UNAVAILABLE"
	if strings.EqualFold(e.marketActivity.Status, "AVAILABLE") {
		broadProvider = "Alpaca market activity"
		broadState = "SEED ONLY"
	}
	detail := fmt.Sprintf("Tiered coverage: %d subscribed symbol(s), %d currently short-window-ready. Full U.S.-market 15s/30s/60s coverage is not claimed without an entitled market-wide event feed; market-activity seeds can promote additional symbols into higher-fidelity monitoring.", len(subscribed), windowReady)
	return RapidMoveCoverage{State: "TIERED_PARTIAL", LiveSymbols: len(subscribed), SubscribedSymbols: len(subscribed), WindowReadySymbols: windowReady, BroadSeedProvider: broadProvider, BroadSeedState: broadState, Detail: detail, UpdatedAt: now.UnixMilli()}
}

func (e *Engine) buildRapidMoveStateLocked(now time.Time) RapidMoveState {
	active := []RapidMoveEvent{}
	for _, ev := range e.rapidMoveEvents {
		if ev.State != "RESOLVED" && ev.UpdatedAt >= now.Add(-30*time.Minute).UnixMilli() {
			active = append(active, ev)
		}
	}
	sort.Slice(active, func(i, j int) bool {
		if active[i].Severity != active[j].Severity {
			return active[i].Severity == "HIGH"
		}
		return active[i].UpdatedAt > active[j].UpdatedAt
	})
	recent := append([]RapidMoveEvent(nil), e.rapidMoveRecent...)
	sort.Slice(recent, func(i, j int) bool { return recent[i].UpdatedAt > recent[j].UpdatedAt })
	if len(recent) > 30 {
		recent = recent[:30]
	}
	status := "READY"
	if len(active) > 0 {
		status = "ACTIVE"
	}
	scorecard := e.rapidMoveScorecard
	scorecard.PolicyVersion = rapidMovePolicyVersion
	return RapidMoveState{Status: status, PolicyVersion: rapidMovePolicyVersion, AdaptiveStage: "PRODUCTION deterministic v1.1 · adaptive learning SHADOW", Coverage: e.rapidMoveCoverageLocked(now), Active: active, Recent: recent, Scorecard: scorecard, Governance: rapidMovePolicyGovernance(), Detail: "Event-first rapid-move and market-shock intelligence; protected Day/Swing/Long formulas are unchanged.", UpdatedAt: now.UnixMilli()}
}

func rapidMoveBestMetric(symbol string, q Quote, history []HistoryPoint, daily []Bar, now time.Time) (RapidMoveWindowMetric, bool) {
	nowMs := q.ProviderTimestamp
	if nowMs <= 0 {
		nowMs = q.UpdatedAt
	}
	nowMs = normalizeObservationMs(nowMs)
	best := RapidMoveWindowMetric{}
	okBest := false
	for _, w := range rapidMoveWindows {
		m, ok := rapidMoveHistoryMetric(history, nowMs, q.Price, w.Duration)
		if !ok {
			continue
		}
		m.Window = w.Label
		m.ThresholdPct = rapidMoveThreshold(w.BaseLimit, symbol, q, daily, now)
		if m.ThresholdPct <= 0 {
			continue
		}
		m.TriggerRatio = math.Abs(m.ReturnPct) / m.ThresholdPct
		if !okBest || m.TriggerRatio > best.TriggerRatio {
			best, okBest = m, true
		}
	}
	return best, okBest
}

func rapidMoveEventID(symbol string, startAt int64) string {
	bucket := startAt / 5000
	return fmt.Sprintf("rapid-%s-%d", normalizeSymbol(symbol), bucket)
}

func (e *Engine) promoteRapidMoveToRadarLocked(ev RapidMoveEvent) {
	if !ev.Alerted || ev.Symbol == "" {
		return
	}
	promotionState := "RAPID MOVE"
	if ev.Classification == "MARKET_SHOCK" {
		promotionState = "MARKET SHOCK"
	}
	p := OpportunityPromotion{Symbol: ev.Symbol, Score: ev.MaterialityScore, State: promotionState, Reasons: append([]string(nil), ev.Reasons...), PromotedAt: ev.DetectedAt, LastConfirmedAt: ev.UpdatedAt, ExpiresAt: ev.UpdatedAt + 15*60_000, ShadowWouldMatch: true}
	found := false
	for i := range e.scanner.Radar.Promotions {
		if normalizeSymbol(e.scanner.Radar.Promotions[i].Symbol) == ev.Symbol {
			old := e.scanner.Radar.Promotions[i]
			if old.PromotedAt > 0 {
				p.PromotedAt = old.PromotedAt
			}
			e.scanner.Radar.Promotions[i] = p
			found = true
			break
		}
	}
	if !found {
		e.scanner.Radar.Promotions = append(e.scanner.Radar.Promotions, p)
	}
	sort.SliceStable(e.scanner.Radar.Promotions, func(i, j int) bool { return e.scanner.Radar.Promotions[i].Score > e.scanner.Radar.Promotions[j].Score })
	if len(e.scanner.Radar.Promotions) > opportunityMaxPromotions {
		e.scanner.Radar.Promotions = e.scanner.Radar.Promotions[:opportunityMaxPromotions]
	}
	e.scanner.Radar.Status = "ACTIVE"
	if ev.Classification == "MARKET_SHOCK" {
		e.scanner.Radar.Message = "Market Shock intelligence promoted a market-wide material event for immediate investigation; normal Opportunity Radar scanning continues."
	} else {
		e.scanner.Radar.Message = "Rapid Move intelligence promoted a material event for immediate investigation; normal Opportunity Radar scanning continues."
	}
	e.livePriorityHints[ev.Symbol] = ev.UpdatedAt
}

func rapidMoveEventAnchorMs(ev RapidMoveEvent) int64 {
	if ev.EventProviderAt > 0 {
		return normalizeObservationMs(ev.EventProviderAt)
	}
	if ev.UpdatedAt > 0 {
		return normalizeObservationMs(ev.UpdatedAt)
	}
	return ev.DetectedAt
}

func rapidMoveOutcomePct(detectedPrice, current float64) float64 {
	if detectedPrice <= 0 || current <= 0 {
		return 0
	}
	return (current/detectedPrice - 1) * 100
}

func (e *Engine) updateRapidMoveOutcomeLocked(ev RapidMoveEvent, q Quote, nowMs int64) (RapidMoveEvent, bool) {
	if ev.DetectedAt <= 0 || q.Price <= 0 {
		return ev, false
	}
	changed := false
	move := rapidMoveOutcomePct(ev.Price, q.Price)
	signed := move
	if ev.Direction == "DOWN" {
		signed = -move
	}
	if signed > ev.MFE {
		ev.MFE = signed
		changed = true
	}
	adverse := -signed
	if adverse > ev.MAE {
		ev.MAE = adverse
		changed = true
	}
	anchor := rapidMoveEventAnchorMs(ev)
	elapsed := nowMs - anchor
	// A material event should de-escalate when most of the original directional
	// move has been given back. Use provider/event time so replay is point-in-time
	// correct and do not wait for the 20-minute outcome resolution.
	if ev.Alerted && ev.State != "RESOLVED" && elapsed >= 45_000 && ev.StartPrice > 0 && ev.Price > 0 {
		detectedMove := rapidMoveOutcomePct(ev.StartPrice, ev.Price)
		currentMove := rapidMoveOutcomePct(ev.StartPrice, q.Price)
		if ev.Direction == "DOWN" {
			detectedMove = -detectedMove
			currentMove = -currentMove
		}
		if detectedMove > 0 {
			retained := currentMove / detectedMove
			if retained < .60 && ev.State != "FADING" {
				ev.State = "FADING"
				changed = true
			}
		}
	}
	set := func(target int64, dst **float64) {
		if elapsed < target || *dst != nil {
			return
		}
		v := move
		*dst = &v
		changed = true
	}
	set(60_000, &ev.Outcome1mPct)
	set(5*60_000, &ev.Outcome5mPct)
	set(20*60_000, &ev.Outcome20mPct)
	if ev.Outcome20mPct != nil && ev.State != "RESOLVED" {
		dir := *ev.Outcome20mPct
		if ev.Direction == "DOWN" {
			dir = -dir
		}
		switch {
		case dir >= 2:
			ev.OutcomeState = "CONTINUED"
		case dir <= -1:
			ev.OutcomeState = "REVERSED"
		default:
			ev.OutcomeState = "FADED / MIXED"
		}
		ev.State = "RESOLVED"
		e.rapidMoveScorecard.OutcomesResolved++
		switch ev.OutcomeState {
		case "CONTINUED":
			e.rapidMoveScorecard.ContinuedOutcomes++
		case "REVERSED":
			e.rapidMoveScorecard.ReversedOutcomes++
		default:
			e.rapidMoveScorecard.FadedMixedOutcomes++
		}
		changed = true
	}
	return ev, changed
}

func (e *Engine) broadcastRapidMoveUpdate(ev RapidMoveEvent) {
	if e == nil || e.app == nil || e.app.hub == nil || !ev.Alerted {
		return
	}
	e.mu.RLock()
	scanner := clone(e.scanner)
	rapidState := e.buildRapidMoveStateLocked(time.Now())
	e.mu.RUnlock()
	e.app.hub.Broadcast(map[string]any{"type": "rapid-move", "event": ev, "rapidMove": rapidState, "scanner": scanner})
}

func (e *Engine) evaluateRapidMoveObservation(symbol string, q Quote) {
	symbol = normalizeSymbol(symbol)
	if e == nil || symbol == "" || symbol == "VIX" || q.Price <= 0 {
		return
	}
	if _, ok := parseUserTicker(symbol); !ok {
		return
	}
	now := time.Now()
	qTime := normalizeObservationMs(q.ProviderTimestamp)
	if qTime <= 0 {
		qTime = q.UpdatedAt
	}
	if qTime <= 0 {
		qTime = now.UnixMilli()
	}
	eventNow := time.UnixMilli(qTime)

	var outcomeUpdate *RapidMoveEvent
	outcomeStateChanged := false
	e.mu.Lock()
	e.rapidMoveScorecard.Observations++
	if existing, ok := e.rapidMoveEvents[symbol]; ok && existing.DetectedAt > 0 {
		updated, changed := e.updateRapidMoveOutcomeLocked(existing, q, qTime)
		if changed {
			updated.UpdatedAt = qTime
			e.rapidMoveEvents[symbol] = updated
			outcomeStateChanged = updated.State != existing.State
			if updated.State == "RESOLVED" && existing.State != "RESOLVED" {
				e.rapidMoveRecent = append(e.rapidMoveRecent, updated)
				if len(e.rapidMoveRecent) > 60 {
					e.rapidMoveRecent = append([]RapidMoveEvent(nil), e.rapidMoveRecent[len(e.rapidMoveRecent)-60:]...)
				}
			}
			copy := updated
			outcomeUpdate = &copy
		}
	}
	history := append([]HistoryPoint(nil), e.history[symbol]...)
	daily := append([]Bar(nil), e.bars[symbol]["daily"]...)
	observations := clone(e.providerQuotes[symbol])
	actions := append([]CorporateAction(nil), e.corporateActions...)
	news := append([]NewsItem(nil), e.news...)
	filings := append([]FilingItem(nil), e.filings...)
	earnings := append([]EarningsItem(nil), e.earnings...)
	reactions := clone(e.catalystReactions)
	allHistory := map[string][]HistoryPoint{"SPY": append([]HistoryPoint(nil), e.history["SPY"]...), "QQQ": append([]HistoryPoint(nil), e.history["QQQ"]...)}
	e.mu.Unlock()
	if outcomeUpdate != nil {
		e.persistRapidMoveEvent(*outcomeUpdate, false)
		if outcomeStateChanged {
			e.broadcastRapidMoveUpdate(*outcomeUpdate)
		}
	}

	metric, ok := rapidMoveBestMetric(symbol, q, history, daily, eventNow)
	if !ok || metric.TriggerRatio < .72 {
		return
	}
	e.mu.Lock()
	e.rapidMoveScorecard.CandidateEvaluations++
	e.mu.Unlock()

	maxSourceSkew := metric.WindowMs / 3
	if maxSourceSkew > 10_000 {
		maxSourceSkew = 10_000
	}
	if maxSourceSkew < 2_000 {
		maxSourceSkew = 2_000
	}
	agreement, agreementDetail := rapidMoveSourceAgreement(observations, qTime, maxSourceSkew)
	mechanical := rapidMoveMechanicalRisk(symbol, actions, eventNow)
	catalystState, catalystSummary, catalystAt, independent := rapidMoveCatalyst(symbol, news, filings, earnings, reactions, eventNow)
	direction := "UP"
	if metric.ReturnPct < 0 {
		direction = "DOWN"
	}
	marketContext := rapidMoveMarketContext(allHistory, qTime, time.Duration(metric.WindowMs)*time.Millisecond, direction)
	relativeVolume, volumeState := rapidMoveRelativeVolume(q, daily)
	spread := rapidMoveSpreadPct(q)
	materiality, reasons := rapidMoveMateriality(metric, q, relativeVolume, volumeState, agreement, catalystState, marketContext, mechanical)

	shadowWouldAlert := metric.TriggerRatio >= .82 && materiality >= 55 && mechanical == "" && q.Price >= 2
	productionAlert := metric.TriggerRatio >= 1.0 && materiality >= 65 && mechanical == "" && q.Price >= 2 && spread <= 1.25 && agreement != "CONFLICT"
	// 5%/60s is a strong baseline trigger, but still requires basic data-quality safeguards.
	if metric.Window == "60s" && math.Abs(metric.ReturnPct) >= 5 && mechanical == "" && q.Price >= 2 && spread <= 1.25 && agreement != "CONFLICT" {
		productionAlert = materiality >= 58
	}

	state := "EARLY"
	severity := "MEDIUM"
	if metric.TriggerRatio >= 1 {
		state = "VALIDATING"
	}
	if productionAlert && catalystState == "CONFIRMED" {
		state = "CONFIRMED"
	}
	if productionAlert && metric.TriggerRatio >= 1.45 {
		state = "EXTENDED"
	}
	if productionAlert && materiality >= 78 {
		severity = "HIGH"
	}
	if productionAlert && metric.Window == "60s" && math.Abs(metric.ReturnPct) >= 5 {
		severity = "HIGH"
	}

	limitations := []string{}
	if volumeState == "PENDING" {
		limitations = append(limitations, "Reliable live relative-volume confirmation is pending for this feed observation.")
	}
	if agreement == "SINGLE_SOURCE" {
		limitations = append(limitations, "Only one current provider price is available; alert remains provenance-labeled.")
	}
	if strings.Contains(marketContext, "PENDING") {
		limitations = append(limitations, "SPY/QQQ short-window comparison is not currently available.")
	}
	limitations = append(limitations, "No full-market 15s/30s/60s coverage is claimed unless an entitled broad event feed is active.")

	id := rapidMoveEventID(symbol, metric.StartAt)
	ev := RapidMoveEvent{
		ID: id, TraceID: id, Symbol: symbol, Direction: direction, State: state, Classification: rapidMoveClassification(marketContext), Severity: severity, Alerted: productionAlert,
		ShadowWouldAlert: shadowWouldAlert, PolicyVersion: rapidMovePolicyVersion, LearningPolicyVersion: rapidMoveLearningPolicyVersion, AdaptiveStage: "PRODUCTION deterministic v1.1 · SHADOW learning",
		Session: marketSessionET(eventNow), Window: metric.Window, WindowMs: metric.WindowMs, WindowStartAt: metric.StartAt,
		EventProviderAt: qTime, ReceivedAt: normalizeObservationMs(q.UpdatedAt), DetectedAt: now.UnixMilli(), UpdatedAt: qTime,
		StartPrice: metric.StartPrice, Price: q.Price, MovePct: metric.ReturnPct, ThresholdPct: metric.ThresholdPct, TriggerRatio: metric.TriggerRatio,
		MaterialityScore: materiality, RelativeVolume: relativeVolume, VolumeState: volumeState, SpreadPct: spread,
		LiquidityState: func() string {
			if spread == 0 {
				return "UNKNOWN"
			}
			if spread <= .35 {
				return "HEALTHY"
			}
			if spread <= .75 {
				return "CAUTION"
			}
			return "RISK"
		}(),
		SourceAgreement: agreement, SourceAgreementDetail: agreementDetail, CatalystState: catalystState, CatalystSummary: catalystSummary,
		CatalystPublishedAt: catalystAt, IndependentSources: independent, MarketContext: marketContext, MechanicalRisk: mechanical,
		Reasons: reasons, Limitations: uniqueStrings(limitations),
	}

	firstAlert := false
	materialUpdate := false
	newCatalystConfirmation := false
	e.mu.Lock()
	existing, exists := e.rapidMoveEvents[symbol]
	if exists && existing.State != "RESOLVED" && qTime-rapidMoveEventAnchorMs(existing) <= 12*60_000 && existing.Direction == ev.Direction {
		ev.ID, ev.TraceID, ev.DetectedAt = existing.ID, existing.TraceID, existing.DetectedAt
		// Preserve the first provider/event anchor across deduped updates so outcome
		// horizons do not slide forward under a busy live stream.
		if existing.EventProviderAt > 0 {
			ev.EventProviderAt = existing.EventProviderAt
		}
		if existing.WindowStartAt > 0 {
			ev.WindowStartAt = existing.WindowStartAt
		}
		wasAlerted := existing.Alerted
		ev.Alerted = existing.Alerted || ev.Alerted
		firstAlert = ev.Alerted && !wasAlerted
		if existing.Classification == "MARKET_SHOCK" && ev.Classification != "MARKET_SHOCK" {
			ev.Classification = existing.Classification
		}
		var retained bool
		ev, retained = rapidMoveApplyHysteresis(existing, ev)
		if retained {
			e.rapidMoveScorecard.HysteresisRetained++
		}
		newCatalystConfirmation = ev.Alerted && existing.CatalystState != "CONFIRMED" && ev.CatalystState == "CONFIRMED"
		materialUpdate = ev.Alerted && (firstAlert || ev.State != existing.State || newCatalystConfirmation || ev.Severity != existing.Severity || ev.Classification != existing.Classification)
		ev.MFE, ev.MAE = existing.MFE, existing.MAE
		ev.Outcome1mPct, ev.Outcome5mPct, ev.Outcome20mPct = existing.Outcome1mPct, existing.Outcome5mPct, existing.Outcome20mPct
		ev.OutcomeState = existing.OutcomeState
		e.rapidMoveScorecard.DuplicateUpdates++
	} else if ev.Alerted {
		firstAlert = true
		materialUpdate = true
	}
	if ev.ShadowWouldAlert && (!exists || !existing.ShadowWouldAlert) {
		e.rapidMoveScorecard.ShadowWouldAlert++
	}
	if mechanical != "" {
		e.rapidMoveScorecard.MechanicalSuppressed++
	}
	if spread > 1.25 || q.Price < 2 {
		e.rapidMoveScorecard.LiquiditySuppressed++
	}
	if agreement == "CONFLICT" {
		e.rapidMoveScorecard.SourceConflictSuppressed++
	}
	if firstAlert {
		e.rapidMoveScorecard.ProductionAlerts++
		if ev.Classification == "MARKET_SHOCK" {
			e.rapidMoveScorecard.MarketShockAlerts++
		}
	}
	if (firstAlert && catalystState == "CONFIRMED") || newCatalystConfirmation {
		e.rapidMoveScorecard.ConfirmedCatalysts++
	}
	if firstAlert && catalystState == "UNEXPLAINED" {
		e.rapidMoveScorecard.UnexplainedMoves++
	}
	e.rapidMoveScorecard.LastEventAt = qTime
	e.rapidMoveEvents[symbol] = ev
	e.promoteRapidMoveToRadarLocked(ev)
	e.lastUpdated["rapid-move"] = qTime
	e.health["rapid-move"] = fmt.Sprintf("%s · %s %+.2f%%/%s · score %.0f", ev.State, symbol, ev.MovePct, ev.Window, ev.MaterialityScore)
	e.mu.Unlock()

	if ev.Alerted || ev.ShadowWouldAlert {
		e.persistRapidMoveEvent(ev, firstAlert)
	}
	if firstAlert {
		e.requestHistoryHydration(symbol)
	}
	if materialUpdate {
		e.broadcastRapidMoveUpdate(ev)
	}
}

func (e *Engine) persistRapidMoveEvent(ev RapidMoveEvent, firstDecision bool) {
	if e == nil || e.app == nil || e.app.persistence == nil || ev.ID == "" {
		return
	}
	raw, err := json.Marshal(ev)
	if err != nil {
		return
	}
	sum := sha256.Sum256(raw)
	batch := PersistenceIntelligenceBatch{Features: []DerivedFeatureRecord{{
		Symbol: ev.Symbol, FeatureKey: "rapid-move-current", FeatureVersion: rapidMovePolicyVersion, AsOf: ev.UpdatedAt,
		SourceHash: hex.EncodeToString(sum[:])[:24], Payload: raw,
	}}}
	if firstDecision {
		evidenceID := ev.ID + ":evidence"
		anchor := rapidMoveEventAnchorMs(ev)
		batch.Evidence = append(batch.Evidence, EvidenceRecord{ID: evidenceID, Symbol: ev.Symbol, Kind: "rapid-move-event", SourceAt: anchor, ObservedAt: ev.DetectedAt, IngestedAt: ev.DetectedAt, KnownAt: ev.DetectedAt, EffectiveFrom: anchor, RevisionID: evidenceID, AmendmentState: "ORIGINAL", Source: "canonical-live-pipeline", Provenance: rapidMovePolicyVersion, FreshnessState: "POINT_IN_TIME", RightsEvidenceRef: "internal:canonical-live-pipeline", RetentionClass: "DECISION_LINEAGE", Payload: raw})
		batch.Decisions = append(batch.Decisions, DecisionLineageRecord{ID: ev.ID + ":decision", Symbol: ev.Symbol, Horizon: "INTRADAY", EvidenceID: evidenceID, DecisionKind: "rapid-move-materiality", DecisionValue: ev.State, FormulaVersion: rapidMovePolicyVersion, CreatedAt: ev.DetectedAt, Payload: raw})
	}
	if ev.State == "RESOLVED" && ev.Outcome20mPct != nil {
		batch.Outcomes = append(batch.Outcomes, OutcomeHistoryRecord{ID: ev.ID + ":outcome-20m", DecisionID: ev.ID + ":decision", Symbol: ev.Symbol, Horizon: "INTRADAY", ObservedAt: ev.UpdatedAt, OutcomeLabel: ev.OutcomeState, Payload: raw})
	}
	e.app.persistence.EnqueueIntelligence(batch)
}

func rapidMoveNotifications(state RapidMoveState, now time.Time) []SmartNotification {
	out := []SmartNotification{}
	nowMs := now.UnixMilli()
	for _, ev := range state.Active {
		if !ev.Alerted || ev.DetectedAt <= 0 || nowMs-ev.UpdatedAt < 0 || nowMs-ev.UpdatedAt > 15*60_000 {
			continue
		}
		kindPrefix := "RAPID MOVE"
		if ev.Classification == "MARKET_SHOCK" {
			kindPrefix = "MARKET SHOCK"
		}
		kind := kindPrefix + " · " + ev.State
		title := fmt.Sprintf("%s · %+.2f%% / %s", ev.Symbol, ev.MovePct, ev.Window)
		message := fmt.Sprintf("Score %.0f/100 · %s", ev.MaterialityScore, ev.CatalystState)
		if ev.CatalystSummary != "" {
			message += " · " + ev.CatalystSummary
		}
		if ev.SourceAgreement != "" {
			message += " · " + ev.SourceAgreement
		}
		out = append(out, SmartNotification{ID: "rapid-notify-" + ev.ID, Severity: ev.Severity, Kind: kind, Symbol: ev.Symbol, Title: title, Message: message, CreatedAt: ev.DetectedAt, ExpiresAt: ev.UpdatedAt + 15*60_000, EventID: ev.TraceID})
	}
	return out
}
