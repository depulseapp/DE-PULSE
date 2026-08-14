package main

import (
	"fmt"
	"hash/fnv"
	"math"
	"sort"
	"strings"
	"time"
	"unicode"
)

// Event Intelligence is a derived/contextual layer over canonical News, Macro,
// Earnings, SEC/Catalyst and reaction stores. It never fetches data and never
// mutates deterministic Day/Swing/Long Score/Action formulas.
type EventNewsIntelligence struct {
	ID                string   `json:"id"`
	Headline          string   `json:"headline"`
	Summary           string   `json:"summary,omitempty"`
	Category          string   `json:"category"`
	Materiality       string   `json:"materiality"`
	Freshness         string   `json:"freshness"`
	Symbols           []string `json:"symbols,omitempty"`
	Source            string   `json:"source"`
	SupportingSources []string `json:"supportingSources,omitempty"`
	URL               string   `json:"url,omitempty"`
	PublishedAt       int64    `json:"publishedAt,omitempty"`
	AgeMs             int64    `json:"ageMs,omitempty"`
	Detail            string   `json:"detail,omitempty"`
}

type EconomicCalendarEntry struct {
	ID              string             `json:"id"`
	Name            string             `json:"name"`
	Region          string             `json:"region"`
	Scope           string             `json:"scope"`
	ProcessingClass string             `json:"processingClass,omitempty"`
	Category        string             `json:"category"`
	Impact          string             `json:"impact"`
	State           string             `json:"state"`
	StartsAt        int64              `json:"startsAt,omitempty"`
	Date            string             `json:"date"`
	TimeKnown       bool               `json:"timeKnown"`
	Actual          *float64           `json:"actual,omitempty"`
	Forecast        *float64           `json:"forecast,omitempty"`
	Previous        *float64           `json:"previous,omitempty"`
	Surprise        *float64           `json:"surprise,omitempty"`
	SurprisePct     *float64           `json:"surprisePct,omitempty"`
	Unit            string             `json:"unit,omitempty"`
	Source          string             `json:"source"`
	SourceURL       string             `json:"sourceUrl,omitempty"`
	UpdatedAt       int64              `json:"updatedAt,omitempty"`
	ReactionOffset  int                `json:"reactionOffsetSec,omitempty"`
	ReactionAt      int64              `json:"reactionAt,omitempty"`
	ReactionMoves   map[string]float64 `json:"reactionMoves,omitempty"`
	Detail          string             `json:"detail,omitempty"`
}

type FedIntelligenceState struct {
	State        string   `json:"state"`
	EventID      string   `json:"eventId,omitempty"`
	Name         string   `json:"name,omitempty"`
	StartsAt     int64    `json:"startsAt,omitempty"`
	CountdownSec int64    `json:"countdownSec,omitempty"`
	Phase        string   `json:"phase,omitempty"`
	Source       string   `json:"source,omitempty"`
	Timeline     []string `json:"timeline,omitempty"`
	Detail       string   `json:"detail"`
}

type ReactionIntelligenceItem struct {
	ID              string             `json:"id"`
	Scope           string             `json:"scope"`
	Symbol          string             `json:"symbol,omitempty"`
	Event           string             `json:"event"`
	EventType       string             `json:"eventType"`
	Phase           string             `json:"phase"`
	State           string             `json:"state"`
	MovePct         float64            `json:"movePct,omitempty"`
	RelativeVolume  float64            `json:"relativeVolume,omitempty"`
	VWAPState       string             `json:"vwapState,omitempty"`
	OpeningRange    string             `json:"openingRange,omitempty"`
	HoldFadeState   string             `json:"holdFadeState,omitempty"`
	ReactionWindows map[string]float64 `json:"reactionWindows,omitempty"`
	CrossAssetMoves map[string]float64 `json:"crossAssetMoves,omitempty"`
	UpdatedAt       int64              `json:"updatedAt,omitempty"`
	Detail          string             `json:"detail,omitempty"`
}

type SmartNotification struct {
	ID        string `json:"id"`
	Severity  string `json:"severity"`
	Kind      string `json:"kind"`
	Symbol    string `json:"symbol,omitempty"`
	Title     string `json:"title"`
	Message   string `json:"message"`
	CreatedAt int64  `json:"createdAt"`
	ExpiresAt int64  `json:"expiresAt,omitempty"`
	EventID   string `json:"eventId,omitempty"`
}

type EventDecisionCorrelation struct {
	State                    string            `json:"state"`
	MarketRisk               string            `json:"marketRisk"`
	AffectedSymbols          []string          `json:"affectedSymbols,omitempty"`
	ReadinessActions         map[string]string `json:"readinessActions,omitempty"`
	QueueActions             []string          `json:"queueActions,omitempty"`
	Reasons                  []string          `json:"reasons,omitempty"`
	DeterministicScoreImpact string            `json:"deterministicScoreImpact"`
	UpdatedAt                int64             `json:"updatedAt,omitempty"`
}

type EventIntelligenceSnapshot struct {
	News          []EventNewsIntelligence    `json:"news"`
	Calendar      []EconomicCalendarEntry    `json:"calendar"`
	Fed           FedIntelligenceState       `json:"fed"`
	Reactions     []ReactionIntelligenceItem `json:"reactions"`
	Notifications []SmartNotification        `json:"notifications"`
	Decision      EventDecisionCorrelation   `json:"decision"`
	SourceHealth  map[string]string          `json:"sourceHealth"`
	UpdatedAt     int64                      `json:"updatedAt,omitempty"`
}

func eventIntelHash(s string) string {
	h := fnv.New64a()
	_, _ = h.Write([]byte(s))
	return fmt.Sprintf("%016x", h.Sum64())
}

func normalizeEventHeadline(s string) string {
	s = strings.ToLower(strings.TrimSpace(s))
	var b strings.Builder
	space := false
	for _, r := range s {
		if unicode.IsLetter(r) || unicode.IsDigit(r) {
			b.WriteRune(r)
			space = false
		} else if !space {
			b.WriteByte(' ')
			space = true
		}
	}
	words := strings.Fields(b.String())
	stop := map[string]bool{"the": true, "a": true, "an": true, "to": true, "of": true, "and": true, "for": true, "on": true, "in": true, "says": true}
	out := words[:0]
	for _, w := range words {
		if !stop[w] {
			out = append(out, w)
		}
	}
	if len(out) > 12 {
		out = out[:12]
	}
	return strings.Join(out, " ")
}

func newsMaterialityAndCategory(n NewsItem) (string, string) {
	txt := strings.ToLower(n.Headline + " " + n.Summary)
	high := []string{"earnings", "guidance", "acquisition", "merger", "bankruptcy", "offering", "fda", "sec investigation", "doj", "antitrust", "ceo resign", "cfo resign", "recall", "data breach", "contract award", "rating downgrade", "rating upgrade", "restatement", "delisting"}
	medium := []string{"partnership", "launch", "approval", "settlement", "lawsuit", "layoff", "restructur", "forecast", "outlook", "dividend", "buyback"}
	category := "MARKET NEWS"
	switch {
	case strings.Contains(txt, "earn") || strings.Contains(txt, "guidance"):
		category = "EARNINGS / GUIDANCE"
	case strings.Contains(txt, "merger") || strings.Contains(txt, "acquisition"):
		category = "M&A"
	case strings.Contains(txt, "sec ") || strings.Contains(txt, "filing") || strings.Contains(txt, "restatement"):
		category = "REGULATORY / FILING"
	case strings.Contains(txt, "fda"):
		category = "REGULATORY"
	case strings.Contains(txt, "offering") || strings.Contains(txt, "buyback") || strings.Contains(txt, "dividend"):
		category = "CAPITAL ACTION"
	case len(n.Symbols) > 0:
		category = "TICKER CATALYST"
	}
	score := 0
	for _, x := range high {
		if strings.Contains(txt, x) {
			score += 2
		}
	}
	for _, x := range medium {
		if strings.Contains(txt, x) {
			score++
		}
	}
	if len(n.Symbols) > 0 {
		score++
	}
	if score >= 3 {
		return "HIGH", category
	}
	if score >= 1 {
		return "MEDIUM", category
	}
	return "LOW", category
}

func newsFreshness(publishedSec int64, now time.Time) (string, int64) {
	if publishedSec <= 0 {
		return "UNAVAILABLE", 0
	}
	stamp := time.Unix(publishedSec, 0)
	if stamp.After(now.Add(5 * time.Minute)) {
		return "INVALID", stamp.Sub(now).Milliseconds()
	}
	age := now.Sub(stamp)
	switch {
	case age <= 6*time.Hour:
		return "FRESH", age.Milliseconds()
	case age <= 24*time.Hour:
		return "CURRENT", age.Milliseconds()
	case age <= 72*time.Hour:
		return "AGING", age.Milliseconds()
	default:
		return "STALE", age.Milliseconds()
	}
}

func buildEventNewsIntelligence(news []NewsItem, now time.Time) []EventNewsIntelligence {
	type cluster struct {
		item    EventNewsIntelligence
		sources map[string]bool
	}
	clusters := map[string]*cluster{}
	for _, n := range news {
		headline := strings.TrimSpace(n.Headline)
		if headline == "" {
			continue
		}
		freshness, age := newsFreshness(n.Datetime, now)
		if freshness == "INVALID" {
			continue
		}
		materiality, category := newsMaterialityAndCategory(n)
		norm := normalizeEventHeadline(headline)
		if norm == "" {
			continue
		}
		// Bucket by 6h so separately recurring stories on later days do not collapse.
		bucket := n.Datetime / (6 * 3600)
		key := eventIntelHash(norm + fmt.Sprintf("|%d", bucket))
		syms := uniqueSymbols(append([]string{}, n.Symbols...))
		src := strings.TrimSpace(n.Source)
		if src == "" {
			src = "Unknown source"
		}
		if c, ok := clusters[key]; ok {
			c.sources[src] = true
			c.item.Symbols = uniqueSymbols(append(c.item.Symbols, syms...))
			if n.Datetime > c.item.PublishedAt {
				c.item.PublishedAt = n.Datetime
				c.item.AgeMs = age
				c.item.Freshness = freshness
				c.item.Headline = headline
				c.item.Summary = n.Summary
				c.item.URL = n.URL
				c.item.Source = src
			}
			if materiality == "HIGH" || (materiality == "MEDIUM" && c.item.Materiality == "LOW") {
				c.item.Materiality, c.item.Category = materiality, category
			}
			continue
		}
		clusters[key] = &cluster{item: EventNewsIntelligence{ID: key, Headline: headline, Summary: strings.TrimSpace(n.Summary), Category: category, Materiality: materiality, Freshness: freshness, Symbols: syms, Source: src, URL: n.URL, PublishedAt: n.Datetime, AgeMs: age}, sources: map[string]bool{src: true}}
	}
	out := make([]EventNewsIntelligence, 0, len(clusters))
	for _, c := range clusters {
		for s := range c.sources {
			c.item.SupportingSources = append(c.item.SupportingSources, s)
		}
		sort.Strings(c.item.SupportingSources)
		if len(c.item.SupportingSources) > 1 {
			c.item.Detail = fmt.Sprintf("%d supporting sources clustered as one event.", len(c.item.SupportingSources))
		}
		out = append(out, c.item)
	}
	rank := map[string]int{"HIGH": 3, "MEDIUM": 2, "LOW": 1}
	sort.Slice(out, func(i, j int) bool {
		if rank[out[i].Materiality] != rank[out[j].Materiality] {
			return rank[out[i].Materiality] > rank[out[j].Materiality]
		}
		return out[i].PublishedAt > out[j].PublishedAt
	})
	if len(out) > 40 {
		out = out[:40]
	}
	return out
}

func economicCalendarState(ev MacroEvent, now time.Time) string {
	if ev.StartsAt <= 0 {
		// An official event can have a known calendar date but no reliable clock
		// time. Preserve lifecycle/date truth so an already-resolved event never
		// reappears as a future "TIME TBD" item.
		if strings.EqualFold(strings.TrimSpace(ev.Lifecycle), "RESOLVED") || strings.EqualFold(strings.TrimSpace(ev.Lifecycle), "HISTORICAL") {
			return "HISTORICAL"
		}
		if ev.Date != "" {
			if d, err := time.ParseInLocation("2006-01-02", ev.Date, now.Location()); err == nil {
				today, _ := time.ParseInLocation("2006-01-02", now.Format("2006-01-02"), now.Location())
				if d.Before(today) {
					return "HISTORICAL"
				}
			}
			return "SCHEDULED · TIME TBD"
		}
		return "UNAVAILABLE"
	}
	t := time.UnixMilli(ev.StartsAt)
	d := t.Sub(now)
	switch {
	case d > 0:
		return "UPCOMING"
	case d >= -90*time.Minute:
		return "REACTION WINDOW"
	case d >= -24*time.Hour:
		return "RELEASED / RECENT"
	default:
		return "HISTORICAL"
	}
}

func economicEventCategory(ev MacroEvent) string {
	s := strings.ToLower(strings.TrimSpace(ev.Name + " " + ev.Source))
	switch {
	case strings.Contains(s, "fomc") || strings.Contains(s, "federal reserve") || strings.Contains(s, "fed rate") || strings.Contains(s, "powell") || strings.Contains(s, "ecb") || strings.Contains(s, "boj") || strings.Contains(s, "central bank"):
		return "CENTRAL BANK"
	case strings.Contains(s, "cpi") || strings.Contains(s, "ppi") || strings.Contains(s, "inflation") || strings.Contains(s, "pce") || strings.Contains(s, "consumer price") || strings.Contains(s, "producer price"):
		return "INFLATION"
	case strings.Contains(s, "nonfarm") || strings.Contains(s, "payroll") || strings.Contains(s, "employment") || strings.Contains(s, "unemployment") || strings.Contains(s, "jobless") || strings.Contains(s, "jolts") || strings.Contains(s, "jobs"):
		return "LABOR"
	case strings.Contains(s, "gdp") || strings.Contains(s, "gross domestic") || strings.Contains(s, "industrial production"):
		return "GROWTH"
	case strings.Contains(s, "retail sales") || strings.Contains(s, "personal income") || strings.Contains(s, "personal spending") || strings.Contains(s, "consumer confidence") || strings.Contains(s, "consumer sentiment"):
		return "CONSUMER"
	case strings.Contains(s, "ism") || strings.Contains(s, "pmi") || strings.Contains(s, "manufactur") || strings.Contains(s, "services"):
		return "BUSINESS ACTIVITY"
	case strings.Contains(s, "housing") || strings.Contains(s, "home sales") || strings.Contains(s, "building permit") || strings.Contains(s, "housing start"):
		return "HOUSING"
	case strings.Contains(s, "eia") || strings.Contains(s, "crude") || strings.Contains(s, "oil") || strings.Contains(s, "natural gas"):
		return "ENERGY"
	case strings.Contains(s, "trade balance") || strings.Contains(s, "imports") || strings.Contains(s, "exports"):
		return "TRADE"
	default:
		return "OTHER"
	}
}

func macroEventScope(ev MacroEvent) string {
	region := strings.TrimSpace(ev.Region)
	if region == "" || strings.EqualFold(region, "US") {
		return "US"
	}
	return "GLOBAL CONTEXT"
}

func macroEventProcessingClass(ev MacroEvent) string {
	if strings.TrimSpace(ev.ProcessingClass) != "" {
		return strings.ToUpper(strings.TrimSpace(ev.ProcessingClass))
	}
	region := strings.TrimSpace(ev.Region)
	if region == "" || strings.EqualFold(region, "US") {
		if strings.EqualFold(strings.TrimSpace(ev.Impact), "HIGH") {
			return "US_MARKET_CRITICAL"
		}
		return "US_CONTEXT"
	}
	return "GLOBAL_CONTEXT"
}

func buildEconomicCalendar(events []MacroEvent, reactions []EventReaction, now time.Time) []EconomicCalendarEntry {
	latest := map[string]EventReaction{}
	for _, r := range reactions {
		if cur, ok := latest[r.EventID]; !ok || r.CapturedAt > cur.CapturedAt {
			latest[r.EventID] = r
		}
	}
	out := make([]EconomicCalendarEntry, 0, len(events))
	for _, ev := range events {
		row := EconomicCalendarEntry{ID: ev.ID, Name: ev.Name, Region: ev.Region, Scope: macroEventScope(ev), ProcessingClass: macroEventProcessingClass(ev), Category: economicEventCategory(ev), Impact: defaultString(ev.Impact, "UNKNOWN"), State: economicCalendarState(ev, now), StartsAt: ev.StartsAt, Date: ev.Date, TimeKnown: ev.TimeKnown, Actual: ev.Actual, Forecast: ev.Expected, Previous: ev.Previous, Unit: ev.Unit, Source: ev.Source, SourceURL: ev.SourceURL, UpdatedAt: ev.UpdatedAt}
		if ev.Actual != nil && ev.Expected != nil {
			x := *ev.Actual - *ev.Expected
			row.Surprise = &x
			if math.Abs(*ev.Expected) > 1e-9 {
				p := x / math.Abs(*ev.Expected) * 100
				row.SurprisePct = &p
			}
		}
		if r, ok := latest[ev.ID]; ok {
			row.ReactionOffset, row.ReactionAt, row.ReactionMoves = r.OffsetSec, r.CapturedAt, clone(r.Moves)
		}
		if ev.Expected == nil || ev.Actual == nil {
			row.Detail = "Actual/forecast surprise remains unavailable until both values are sourced; no synthetic consensus."
		}
		out = append(out, row)
	}
	sort.Slice(out, func(i, j int) bool {
		ai, aj := out[i].StartsAt, out[j].StartsAt
		if ai == 0 {
			ai = 1 << 62
		}
		if aj == 0 {
			aj = 1 << 62
		}
		futureI, futureJ := ai >= now.UnixMilli(), aj >= now.UnixMilli()
		if futureI != futureJ {
			return futureI
		}
		if futureI {
			return ai < aj
		}
		return ai > aj
	})
	if len(out) > 60 {
		out = out[:60]
	}
	return out
}

func isFedCalendarEntry(e EconomicCalendarEntry) bool {
	s := strings.ToLower(e.Name + " " + e.Source)
	// "Minutes" and "press conference" are not uniquely Federal Reserve
	// concepts. Require an explicit Fed/FOMC anchor so ECB/BOJ/etc. events do
	// not leak into the Fed intelligence lane.
	fedAnchor := strings.Contains(s, "fomc") || strings.Contains(s, "federal reserve") || strings.Contains(s, "fed rate") || strings.Contains(s, "fed minutes") || strings.Contains(s, "powell")
	if fedAnchor {
		return true
	}
	return strings.Contains(s, "press conference") && strings.Contains(s, "fed")
}

func buildFedIntelligence(calendar []EconomicCalendarEntry, now time.Time) FedIntelligenceState {
	candidates := []EconomicCalendarEntry{}
	for _, e := range calendar {
		if isFedCalendarEntry(e) {
			candidates = append(candidates, e)
		}
	}
	if len(candidates) == 0 {
		return FedIntelligenceState{State: "UNAVAILABLE", Detail: "No sourced Federal Reserve/FOMC calendar event is available."}
	}
	sort.Slice(candidates, func(i, j int) bool {
		di, dj := math.Abs(float64(candidates[i].StartsAt-now.UnixMilli())), math.Abs(float64(candidates[j].StartsAt-now.UnixMilli()))
		if candidates[i].StartsAt == 0 {
			di = math.MaxFloat64
		}
		if candidates[j].StartsAt == 0 {
			dj = math.MaxFloat64
		}
		return di < dj
	})
	best := candidates[0]
	phase := best.State
	countdown := int64(0)
	if best.StartsAt > 0 {
		countdown = (best.StartsAt - now.UnixMilli()) / 1000
	}
	state := "AVAILABLE"
	if best.StartsAt > 0 && countdown < -7*24*3600 {
		state = "HISTORICAL"
	} else if best.StartsAt > 0 && countdown > 0 {
		state = "UPCOMING"
	} else if best.StartsAt > 0 && countdown >= -5400 {
		state = "REACTION"
	} else {
		state = "RECENT"
	}
	timeline := []string{}
	for _, e := range candidates {
		if best.Date != "" && e.Date != best.Date {
			continue
		}
		label := e.Name
		if e.StartsAt > 0 {
			label += " · " + time.UnixMilli(e.StartsAt).Format("15:04 MST")
		}
		timeline = append(timeline, label)
	}
	if len(timeline) > 8 {
		timeline = timeline[:8]
	}
	return FedIntelligenceState{State: state, EventID: best.ID, Name: best.Name, StartsAt: best.StartsAt, CountdownSec: countdown, Phase: phase, Source: best.Source, Timeline: timeline, Detail: "Sourced Fed lifecycle only. Decision, statement, press conference and minutes are separated only when the canonical calendar provides distinct evidence."}
}

func reactionState(move float64, phase string) string {
	if strings.Contains(strings.ToUpper(phase), "COMPLETE") {
		if move >= 1 {
			return "POSITIVE HOLD"
		}
		if move <= -1 {
			return "NEGATIVE HOLD"
		}
		return "MIXED / MUTED"
	}
	if move >= 1 {
		return "POSITIVE"
	}
	if move <= -1 {
		return "NEGATIVE"
	}
	return "PRICE DISCOVERY"
}

func buildReactionIntelligence(cats map[string]CatalystReactionState, macro []EventReaction) []ReactionIntelligenceItem {
	out := []ReactionIntelligenceItem{}
	for s, c := range cats {
		out = append(out, ReactionIntelligenceItem{ID: "catalyst-" + s, Scope: "TICKER", Symbol: s, Event: c.Trigger, EventType: defaultString(c.TriggerType, "CATALYST"), Phase: c.Phase, State: reactionState(c.MovePercent, c.Phase), MovePct: c.MovePercent, RelativeVolume: c.RelativeVolume, VWAPState: c.VWAPState, OpeningRange: c.OpeningRangeState, HoldFadeState: c.HoldFadeState, ReactionWindows: clone(c.ReactionPercent), UpdatedAt: c.UpdatedAt, Detail: "Canonical Catalyst Watch lifecycle; context only."})
	}
	for _, r := range macro {
		move := r.Moves["SPY"]
		out = append(out, ReactionIntelligenceItem{ID: fmt.Sprintf("macro-%s-%d", r.EventID, r.OffsetSec), Scope: "MARKET", Event: r.EventID, EventType: "MACRO", Phase: fmt.Sprintf("%ds", r.OffsetSec), State: reactionState(move, fmt.Sprintf("%ds", r.OffsetSec)), MovePct: move, CrossAssetMoves: clone(r.Moves), UpdatedAt: r.CapturedAt, Detail: "Captured cross-asset reaction from canonical Event Mode; no causal claim beyond observed timing."})
	}
	sort.Slice(out, func(i, j int) bool { return out[i].UpdatedAt > out[j].UpdatedAt })
	if len(out) > 30 {
		out = out[:30]
	}
	return out
}

func dateDistanceDays(raw string, now time.Time) (int, bool) {
	if strings.TrimSpace(raw) == "" {
		return 0, false
	}
	loc := easternLocation()
	t, err := time.ParseInLocation("2006-01-02", raw, loc)
	if err != nil {
		return 0, false
	}
	nowET := now.In(loc)
	// Compare calendar dates in UTC rather than elapsed local hours so DST
	// 23/25-hour days cannot turn "two calendar days away" into one day.
	a := time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, time.UTC)
	b := time.Date(nowET.Year(), nowET.Month(), nowET.Day(), 0, 0, 0, 0, time.UTC)
	return int(a.Sub(b) / (24 * time.Hour)), true
}

func buildEventDecisionCorrelation(calendar []EconomicCalendarEntry, news []EventNewsIntelligence, earnings []EarningsItem, cats map[string]CatalystReactionState, mode EventModeState, last map[string]int64, now time.Time) EventDecisionCorrelation {
	out := EventDecisionCorrelation{State: "NORMAL", MarketRisk: "NORMAL", ReadinessActions: map[string]string{}, DeterministicScoreImpact: "NONE — event context may change Trade Readiness / Decision Queue attention but never Day/Swing/Long Score/Action formulas.", UpdatedAt: now.UnixMilli()}
	affected := []string{}
	setRisk := func(r string, reason string) {
		rank := map[string]int{"NORMAL": 0, "ELEVATED": 1, "HIGH": 2, "DATA DEGRADED": 3}
		if rank[r] > rank[out.MarketRisk] {
			out.MarketRisk = r
		}
		if reason != "" {
			out.Reasons = append(out.Reasons, reason)
		}
	}
	macroCheck := last["macro-events"]
	macroAge := now.UnixMilli() - macroCheck
	if macroCheck <= 0 || macroAge > 24*3600_000 || macroCheck > now.Add(5*time.Minute).UnixMilli() {
		reason := "Economic-calendar check is unavailable or stale."
		if macroCheck > now.Add(5*time.Minute).UnixMilli() {
			reason = "Economic-calendar check timestamp is future-skewed and invalid."
		}
		setRisk("DATA DEGRADED", reason)
	}
	if mode.Active {
		setRisk("HIGH", fmt.Sprintf("%s is in an active high-impact event window.", defaultString(mode.Name, "Macro event")))
		affected = append(affected, mode.AffectedSymbols...)
	}
	for _, e := range calendar {
		if strings.ToUpper(e.Impact) != "HIGH" || e.StartsAt <= 0 {
			continue
		}
		d := time.UnixMilli(e.StartsAt).Sub(now)
		if d >= 0 && d <= 60*time.Minute {
			setRisk("HIGH", fmt.Sprintf("%s due within 60 minutes.", e.Name))
		} else if d > 60*time.Minute && d <= 24*time.Hour {
			setRisk("ELEVATED", fmt.Sprintf("%s due within 24 hours.", e.Name))
		}
	}
	for _, n := range news {
		if n.Materiality != "HIGH" || (n.Freshness != "FRESH" && n.Freshness != "CURRENT") {
			continue
		}
		for _, s := range n.Symbols {
			s = normalizeSymbol(s)
			if s == "" {
				continue
			}
			affected = append(affected, s)
			out.ReadinessActions[s] = "EVENT RISK · material news requires review"
			out.QueueActions = append(out.QueueActions, fmt.Sprintf("%s · review material news: %s", s, n.Headline))
		}
	}
	for _, er := range earnings {
		s := normalizeSymbol(er.Symbol)
		if s == "" {
			continue
		}
		d, ok := dateDistanceDays(er.Date, now)
		if !ok || d < 0 || d > 1 {
			continue
		}
		if er.EPSActual == nil && er.RevenueActual == nil {
			affected = append(affected, s)
			out.ReadinessActions[s] = "EVENT RISK · earnings due within 1 day"
			out.QueueActions = append(out.QueueActions, fmt.Sprintf("%s · earnings due %s", s, er.Date))
		}
	}
	for s, c := range cats {
		if c.TriggerAt <= 0 || c.CompletedAt > 0 {
			continue
		}
		affected = append(affected, s)
		out.ReadinessActions[s] = "EVENT RISK · catalyst reaction still price-discovering"
		out.QueueActions = append(out.QueueActions, fmt.Sprintf("%s · catalyst reaction %s", s, defaultString(c.Phase, c.State)))
	}
	affected = uniqueSymbols(affected)
	sort.Strings(affected)
	out.AffectedSymbols = affected
	out.QueueActions = uniqueStrings(out.QueueActions)
	if len(out.QueueActions) > 12 {
		out.QueueActions = out.QueueActions[:12]
	}
	out.Reasons = uniqueStrings(out.Reasons)
	switch out.MarketRisk {
	case "HIGH":
		out.State = "EVENT RISK HIGH"
	case "ELEVATED":
		out.State = "EVENT RISK ELEVATED"
	case "DATA DEGRADED":
		out.State = "DATA DEGRADED"
	default:
		if len(affected) > 0 {
			out.State = "TICKER EVENTS ACTIVE"
		}
	}
	return out
}

type SmartNotificationContext struct {
	Validation     SignalValidationState
	SEC            map[string]SECIntelligenceSummary
	Freshness      []FreshnessDiagnostic
	ProviderRouter ProviderRouterSnapshot
	Scanner        ScannerState
	RapidMove      RapidMoveState
}

func notificationTimestamp(raw string, fallback int64) int64 {
	for _, layout := range []string{time.RFC3339, "2006-01-02T15:04:05Z", "2006-01-02"} {
		if t, err := time.Parse(layout, strings.TrimSpace(raw)); err == nil {
			return t.UnixMilli()
		}
	}
	return fallback
}

func notificationRank(state string) int {
	s := strings.ToUpper(strings.TrimSpace(state))
	switch {
	case strings.Contains(s, "BLOCK") || strings.Contains(s, "RISK") || strings.Contains(s, "POOR") || strings.Contains(s, "STALE") || strings.Contains(s, "ERROR"):
		return 4
	case strings.Contains(s, "CAUTION") || strings.Contains(s, "DEGRADED") || strings.Contains(s, "EXTENDED"):
		return 3
	case strings.Contains(s, "WAIT") || strings.Contains(s, "WATCH") || strings.Contains(s, "NEUTRAL"):
		return 2
	case strings.Contains(s, "READY") || strings.Contains(s, "HEALTHY") || strings.Contains(s, "TRADEABLE") || strings.Contains(s, "FRESH") || strings.Contains(s, "LIVE"):
		return 1
	}
	return 0
}

func signalStateChangeNotifications(v SignalValidationState, now time.Time) []SmartNotification {
	groups := map[string][]SignalSnapshot{}
	for _, x := range v.Snapshots {
		if x.Symbol == "" || x.Horizon == "" || x.Timestamp <= 0 {
			continue
		}
		k := normalizeSymbol(x.Symbol) + "|" + strings.ToLower(x.Horizon)
		groups[k] = append(groups[k], x)
	}
	out := []SmartNotification{}
	nowms := now.UnixMilli()
	for _, rows := range groups {
		sort.Slice(rows, func(i, j int) bool { return rows[i].Timestamp > rows[j].Timestamp })
		if len(rows) < 2 {
			continue
		}
		cur, prev := rows[0], rows[1]
		if nowms-cur.Timestamp < 0 || nowms-cur.Timestamp > 24*3600_000 {
			continue
		}
		sym := normalizeSymbol(cur.Symbol)
		horizon := strings.ToUpper(cur.Horizon)
		if cur.Readiness != "" && prev.Readiness != "" && !strings.EqualFold(cur.Readiness, prev.Readiness) {
			sev := "MEDIUM"
			if notificationRank(cur.Readiness) > notificationRank(prev.Readiness) {
				sev = "HIGH"
			}
			out = append(out, SmartNotification{ID: fmt.Sprintf("readiness-%s-%s-%d", sym, strings.ToLower(cur.Horizon), cur.Timestamp), Severity: sev, Kind: "READINESS CHANGED", Symbol: sym, Title: sym + " · " + horizon + " readiness", Message: prev.Readiness + " → " + cur.Readiness, CreatedAt: cur.Timestamp, ExpiresAt: cur.Timestamp + 12*3600_000})
		}
		prevIn := prev.EntryLow > 0 && prev.EntryHigh >= prev.EntryLow && prev.Price >= prev.EntryLow && prev.Price <= prev.EntryHigh
		curIn := cur.EntryLow > 0 && cur.EntryHigh >= cur.EntryLow && cur.Price >= cur.EntryLow && cur.Price <= cur.EntryHigh
		if prevIn != curIn {
			kind, msg := "ENTRY ZONE LOST", fmt.Sprintf("Price %.2f left %.2f–%.2f", cur.Price, cur.EntryLow, cur.EntryHigh)
			sev := "MEDIUM"
			if curIn {
				kind = "ENTRY ZONE REACHED"
				msg = fmt.Sprintf("Price %.2f entered %.2f–%.2f", cur.Price, cur.EntryLow, cur.EntryHigh)
			}
			out = append(out, SmartNotification{ID: fmt.Sprintf("entry-%s-%s-%d", sym, strings.ToLower(cur.Horizon), cur.Timestamp), Severity: sev, Kind: kind, Symbol: sym, Title: sym + " · " + horizon + " entry zone", Message: msg, CreatedAt: cur.Timestamp, ExpiresAt: cur.Timestamp + 8*3600_000})
		}
		if cur.MarketTradeability != "" && prev.MarketTradeability != "" && !strings.EqualFold(cur.MarketTradeability, prev.MarketTradeability) {
			kind, sev := "TRADEABILITY RECOVERED", "MEDIUM"
			if notificationRank(cur.MarketTradeability) > notificationRank(prev.MarketTradeability) {
				kind = "TRADEABILITY DETERIORATED"
				sev = "HIGH"
			}
			out = append(out, SmartNotification{ID: fmt.Sprintf("tradeability-%s-%d", sym, cur.Timestamp), Severity: sev, Kind: kind, Symbol: sym, Title: sym + " · Tradeability", Message: prev.MarketTradeability + " → " + cur.MarketTradeability, CreatedAt: cur.Timestamp, ExpiresAt: cur.Timestamp + 12*3600_000})
		}
		if cur.LiquidityState != "" && prev.LiquidityState != "" && !strings.EqualFold(cur.LiquidityState, prev.LiquidityState) && notificationRank(cur.LiquidityState) > notificationRank(prev.LiquidityState) {
			out = append(out, SmartNotification{ID: fmt.Sprintf("liquidity-%s-%d", sym, cur.Timestamp), Severity: "HIGH", Kind: "LIQUIDITY DETERIORATED", Symbol: sym, Title: sym + " · Liquidity risk", Message: prev.LiquidityState + " → " + cur.LiquidityState, CreatedAt: cur.Timestamp, ExpiresAt: cur.Timestamp + 8*3600_000})
		}
	}
	return out
}

func insiderPurchaseNotifications(sec map[string]SECIntelligenceSummary, now time.Time) []SmartNotification {
	out := []SmartNotification{}
	cutoff := now.AddDate(0, 0, -7).UnixMilli()
	for sym, sum := range sec {
		for _, tx := range sum.RecentInsiderTransactions {
			if !strings.EqualFold(tx.Classification, "BUY") {
				continue
			}
			at := notificationTimestamp(defaultString(tx.FiledAt, tx.TransactionDate), sum.UpdatedAt)
			if at < cutoff || at > now.Add(24*time.Hour).UnixMilli() {
				continue
			}
			actor := defaultString(tx.Actor, "Insider")
			msg := fmt.Sprintf("%s · %.0f shares", actor, tx.Shares)
			if tx.Value > 0 {
				msg += fmt.Sprintf(" · $%.0f", tx.Value)
			}
			id := eventIntelHash(strings.Join([]string{normalizeSymbol(sym), actor, tx.TransactionDate, tx.Code, fmt.Sprintf("%.0f", tx.Shares)}, "|"))
			out = append(out, SmartNotification{ID: "insider-buy-" + id, Severity: "HIGH", Kind: "INSIDER PURCHASE", Symbol: normalizeSymbol(sym), Title: normalizeSymbol(sym) + " · Insider purchase", Message: msg, CreatedAt: at, ExpiresAt: at + 7*24*3600_000})
		}
	}
	return out
}

func dataHealthNotifications(fresh []FreshnessDiagnostic, router ProviderRouterSnapshot, now time.Time) []SmartNotification {
	out := []SmartNotification{}
	nowms := now.UnixMilli()
	for _, r := range fresh {
		st := strings.ToUpper(r.State)
		if st != "STALE" && st != "ERROR" && st != "UNAVAILABLE" && st != "DELAYED" {
			continue
		}
		at := r.ReceivedAt
		if at <= 0 {
			at = r.ProviderTimestamp
		}
		if at <= 0 {
			at = nowms
		}
		sev := "MEDIUM"
		if st == "STALE" || st == "ERROR" || st == "UNAVAILABLE" {
			sev = "HIGH"
		}
		out = append(out, SmartNotification{ID: "data-degraded-" + strings.ToLower(strings.ReplaceAll(r.Dataset, " ", "-")) + "-" + st, Severity: sev, Kind: "DATA DEGRADED", Title: r.Dataset + " · " + st, Message: defaultString(r.Reason, "Canonical dataset is degraded; recovery/fallback rules apply."), CreatedAt: at, ExpiresAt: nowms + 2*3600_000})
	}
	for _, route := range router.Routes {
		for _, hop := range route.Route {
			if !strings.EqualFold(hop.Recovery, "RECOVERED") || hop.LastSuccess <= 0 || hop.LastSuccess < hop.LastFailure || nowms-hop.LastSuccess > 2*3600_000 {
				continue
			}
			out = append(out, SmartNotification{ID: fmt.Sprintf("data-recovered-%s-%s-%d", strings.ToLower(route.Dataset), strings.ToLower(hop.Provider), hop.LastSuccess), Severity: "MEDIUM", Kind: "DATA RECOVERED", Title: route.Dataset + " · Recovered", Message: hop.Provider + " recovered through the canonical Provider Router.", CreatedAt: hop.LastSuccess, ExpiresAt: hop.LastSuccess + 2*3600_000})
		}
	}
	return out
}

func opportunityRadarNotifications(scanner ScannerState, now time.Time) []SmartNotification {
	out := []SmartNotification{}
	nowms := now.UnixMilli()
	for _, p := range scanner.Radar.Promotions {
		if p.Symbol == "" || p.PromotedAt <= 0 || nowms-p.PromotedAt < 0 || nowms-p.PromotedAt > 10*60_000 {
			continue
		}
		msg := fmt.Sprintf("Opportunity Radar promoted %s at %.0f/100", p.Symbol, p.Score)
		if len(p.Reasons) > 0 {
			msg += " · " + strings.Join(p.Reasons[:minInt(2, len(p.Reasons))], " · ")
		}
		out = append(out, SmartNotification{ID: fmt.Sprintf("opportunity-%s-%d", p.Symbol, p.PromotedAt), Severity: "MEDIUM", Kind: "OPPORTUNITY RADAR", Symbol: p.Symbol, Title: p.Symbol + " · unusual activity", Message: msg, CreatedAt: p.PromotedAt, ExpiresAt: p.ExpiresAt})
	}
	return out
}

func buildEventNotifications(news []EventNewsIntelligence, calendar []EconomicCalendarEntry, reactions []ReactionIntelligenceItem, decision EventDecisionCorrelation, now time.Time) []SmartNotification {
	out := []SmartNotification{}
	nowms := now.UnixMilli()
	for _, n := range news {
		if n.Materiality != "HIGH" || n.PublishedAt <= 0 || nowms-n.PublishedAt*1000 < 0 || nowms-n.PublishedAt*1000 > 2*3600_000 {
			continue
		}
		sym := ""
		if len(n.Symbols) > 0 {
			sym = n.Symbols[0]
		}
		out = append(out, SmartNotification{ID: "news-" + n.ID, Severity: "HIGH", Kind: "NEW MATERIAL NEWS", Symbol: sym, Title: defaultString(sym, "Market") + " · Material news", Message: n.Headline, CreatedAt: n.PublishedAt * 1000, ExpiresAt: n.PublishedAt*1000 + 6*3600_000, EventID: n.ID})
	}
	for _, e := range calendar {
		if strings.ToUpper(e.Impact) != "HIGH" || e.StartsAt <= 0 {
			continue
		}
		d := e.StartsAt - nowms
		// Once the event window has been entered, retain the same stable
		// notification through the 90-minute reaction window. The ID/CreatedAt
		// remain stable, so recomputation does not manufacture a new alert.
		if d >= -90*60_000 && d <= 60*60_000 {
			message := fmt.Sprintf("High-impact event due in %d minutes · %s", d/60_000, e.Source)
			if d < 0 {
				message = fmt.Sprintf("High-impact event released %d minutes ago · reaction window active · %s", (-d)/60_000, e.Source)
			}
			out = append(out, SmartNotification{ID: "window-" + e.ID, Severity: "HIGH", Kind: "EVENT WINDOW ENTERED", Title: e.Name, Message: message, CreatedAt: e.StartsAt - 60*60_000, ExpiresAt: e.StartsAt + 90*60_000, EventID: e.ID})
		}
	}
	for _, r := range reactions {
		if r.UpdatedAt <= 0 || nowms-r.UpdatedAt < 0 || nowms-r.UpdatedAt > 60*60_000 {
			continue
		}
		material := math.Abs(r.MovePct) >= 1 || r.RelativeVolume >= 1.5
		if r.Scope == "MARKET" {
			material = math.Abs(r.MovePct) >= 0.5
		}
		if !material {
			continue
		}
		out = append(out, SmartNotification{ID: "reaction-" + r.ID, Severity: "MEDIUM", Kind: "REACTION UPDATE", Symbol: r.Symbol, Title: defaultString(r.Symbol, "Market") + " · Reaction update", Message: fmt.Sprintf("%s · %s · move %.2f%%", r.Event, r.Phase, r.MovePct), CreatedAt: r.UpdatedAt, ExpiresAt: r.UpdatedAt + 2*3600_000, EventID: r.ID})
	}
	sort.Slice(out, func(i, j int) bool {
		if out[i].Severity != out[j].Severity {
			return out[i].Severity == "HIGH"
		}
		return out[i].CreatedAt > out[j].CreatedAt
	})
	seen := map[string]bool{}
	dedup := out[:0]
	for _, n := range out {
		if seen[n.ID] {
			continue
		}
		seen[n.ID] = true
		dedup = append(dedup, n)
	}
	out = dedup
	if len(out) > 12 {
		out = out[:12]
	}
	return out
}

func buildSmartNotifications(news []EventNewsIntelligence, calendar []EconomicCalendarEntry, reactions []ReactionIntelligenceItem, decision EventDecisionCorrelation, now time.Time, contexts ...SmartNotificationContext) []SmartNotification {
	out := buildEventNotifications(news, calendar, reactions, decision, now)
	if len(contexts) > 0 {
		c := contexts[0]
		out = append(out, signalStateChangeNotifications(c.Validation, now)...)
		out = append(out, insiderPurchaseNotifications(c.SEC, now)...)
		out = append(out, dataHealthNotifications(c.Freshness, c.ProviderRouter, now)...)
		out = append(out, opportunityRadarNotifications(c.Scanner, now)...)
		out = append(out, rapidMoveNotifications(c.RapidMove, now)...)
	}
	sort.Slice(out, func(i, j int) bool {
		if out[i].Severity != out[j].Severity {
			return out[i].Severity == "HIGH"
		}
		return out[i].CreatedAt > out[j].CreatedAt
	})
	seen := map[string]bool{}
	dedup := out[:0]
	for _, n := range out {
		if seen[n.ID] {
			continue
		}
		seen[n.ID] = true
		dedup = append(dedup, n)
	}
	if len(dedup) > 20 {
		dedup = dedup[:20]
	}
	return dedup
}

func buildEventIntelligenceSnapshot(news []NewsItem, events []MacroEvent, eventMode EventModeState, macroReactions []EventReaction, cats map[string]CatalystReactionState, earnings []EarningsItem, last map[string]int64, health map[string]string, now time.Time, contexts ...SmartNotificationContext) EventIntelligenceSnapshot {
	newsIntel := buildEventNewsIntelligence(news, now)
	calendar := buildEconomicCalendar(events, macroReactions, now)
	reactions := buildReactionIntelligence(cats, macroReactions)
	decision := buildEventDecisionCorrelation(calendar, newsIntel, earnings, cats, eventMode, last, now)
	sourceHealth := map[string]string{"news": defaultString(health["news"], "unknown"), "macroEvents": defaultString(health["macro-events"], "unknown"), "eventMode": defaultString(health["event-mode"], "idle"), "catalystWatch": defaultString(health["catalyst-watch"], "unknown")}
	return EventIntelligenceSnapshot{News: newsIntel, Calendar: calendar, Fed: buildFedIntelligence(calendar, now), Reactions: reactions, Notifications: buildSmartNotifications(newsIntel, calendar, reactions, decision, now, contexts...), Decision: decision, SourceHealth: sourceHealth, UpdatedAt: now.UnixMilli()}
}
