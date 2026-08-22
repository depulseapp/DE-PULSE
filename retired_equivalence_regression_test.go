package main

import (
	"strings"
	"testing"
	"time"
)

func TestEventIntelligenceNewsClustersDuplicateSourcesAndRanksMateriality(t *testing.T) {
	now := time.Date(2026, 8, 12, 14, 0, 0, 0, time.UTC)
	news := []NewsItem{
		{ID: 1, Datetime: now.Add(-20 * time.Minute).Unix(), Headline: "NVDA raises guidance after earnings", Source: "Finnhub", Symbols: []string{"NVDA"}},
		{ID: 2, Datetime: now.Add(-18 * time.Minute).Unix(), Headline: "NVDA raises guidance after earnings!", Source: "Marketaux", Symbols: []string{"NVDA"}},
		{ID: 3, Datetime: now.Add(-10 * time.Minute).Unix(), Headline: "Markets open mixed", Source: "Finnhub"},
	}
	rows := buildEventNewsIntelligence(news, now)
	if len(rows) != 2 {
		t.Fatalf("expected duplicate story cluster, got %d: %+v", len(rows), rows)
	}
	if rows[0].Materiality != "HIGH" || rows[0].Category != "EARNINGS / GUIDANCE" || len(rows[0].SupportingSources) != 2 {
		t.Fatalf("material news clustering wrong: %+v", rows[0])
	}
	if rows[0].Freshness != "FRESH" {
		t.Fatalf("freshness wrong: %+v", rows[0])
	}
}

func TestEventIntelligenceNewsRejectsFutureSkew(t *testing.T) {
	now := time.Now()
	rows := buildEventNewsIntelligence([]NewsItem{{Datetime: now.Add(10 * time.Minute).Unix(), Headline: "Future headline", Source: "X"}}, now)
	if len(rows) != 0 {
		t.Fatalf("future-skewed news must not become event intelligence: %+v", rows)
	}
}

func TestEventIntelligenceMacroSurpriseRequiresActualAndForecast(t *testing.T) {
	now := time.Now()
	a, f, p := 3.1, 2.8, 2.9
	rows := buildEconomicCalendar([]MacroEvent{
		{ID: "cpi", Name: "Consumer Price Index", Impact: "HIGH", StartsAt: now.Add(-5 * time.Minute).UnixMilli(), Date: now.Format("2006-01-02"), TimeKnown: true, Actual: &a, Expected: &f, Previous: &p, Source: "BLS"},
		{ID: "gdp", Name: "GDP", Impact: "HIGH", StartsAt: now.Add(time.Hour).UnixMilli(), Date: now.Format("2006-01-02"), TimeKnown: true, Source: "BEA"},
	}, nil, now)
	var cpi, gdp EconomicCalendarEntry
	for _, r := range rows {
		if r.ID == "cpi" {
			cpi = r
		}
		if r.ID == "gdp" {
			gdp = r
		}
	}
	if cpi.Surprise == nil || *cpi.Surprise < 0.29 || *cpi.Surprise > 0.31 {
		t.Fatalf("surprise incorrect: %+v", cpi)
	}
	if gdp.Surprise != nil || !strings.Contains(gdp.Detail, "no synthetic consensus") {
		t.Fatalf("missing consensus must remain unavailable: %+v", gdp)
	}
}

func TestEventIntelligenceFedUsesSourcedTimeline(t *testing.T) {
	now := time.Now()
	date := now.Format("2006-01-02")
	cal := []EconomicCalendarEntry{
		{ID: "fed1", Name: "FOMC Meeting", Source: "Federal Reserve", Date: date, StartsAt: now.Add(30 * time.Minute).UnixMilli(), State: "UPCOMING"},
		{ID: "fed2", Name: "Press Conference", Source: "Federal Reserve", Date: date, StartsAt: now.Add(90 * time.Minute).UnixMilli(), State: "UPCOMING"},
	}
	x := buildFedIntelligence(cal, now)
	if x.State != "UPCOMING" || len(x.Timeline) != 2 || x.CountdownSec <= 0 {
		t.Fatalf("fed intelligence wrong: %+v", x)
	}
}

func TestEventIntelligenceDecisionCorrelatesWithoutScoreMutation(t *testing.T) {
	now := time.Now()
	last := map[string]int64{"macro-events": now.UnixMilli()}
	cal := []EconomicCalendarEntry{{ID: "cpi", Name: "CPI", Impact: "HIGH", StartsAt: now.Add(30 * time.Minute).UnixMilli()}}
	news := []EventNewsIntelligence{{ID: "n1", Headline: "NVDA raises guidance", Materiality: "HIGH", Freshness: "FRESH", Symbols: []string{"NVDA"}}}
	d := buildEventDecisionCorrelation(cal, news, []EarningsItem{{Symbol: "AAPL", Date: now.In(easternLocation()).Format("2006-01-02")}}, map[string]CatalystReactionState{}, EventModeState{}, last, now)
	if d.MarketRisk != "HIGH" || d.ReadinessActions["NVDA"] == "" || d.ReadinessActions["AAPL"] == "" {
		t.Fatalf("event decision correlation incomplete: %+v", d)
	}
	if !strings.HasPrefix(d.DeterministicScoreImpact, "NONE") {
		t.Fatalf("deterministic mutation guard missing: %+v", d)
	}
}

func TestEventIntelligenceNotificationsAreEventTriggeredNotConditionCards(t *testing.T) {
	now := time.Now()
	n := []EventNewsIntelligence{{ID: "n1", Headline: "Material", Materiality: "HIGH", PublishedAt: now.Add(-10 * time.Minute).Unix(), Symbols: []string{"NVDA"}}}
	cal := []EconomicCalendarEntry{{ID: "fomc", Name: "FOMC", Impact: "HIGH", StartsAt: now.Add(45 * time.Minute).UnixMilli(), Source: "Federal Reserve"}}
	rows := buildSmartNotifications(n, cal, nil, EventDecisionCorrelation{MarketRisk: "HIGH"}, now)
	if len(rows) != 2 {
		t.Fatalf("expected only actual event-change notifications, got %+v", rows)
	}
	for _, x := range rows {
		if x.Kind == "TRADEABILITY" {
			t.Fatalf("current-condition card leaked into notifications: %+v", x)
		}
	}
}

func TestEventIntelligenceReactionReusesCatalystAndMacroEvidence(t *testing.T) {
	now := time.Now().UnixMilli()
	cats := map[string]CatalystReactionState{"NVDA": {Trigger: "Earnings release", TriggerType: "EARNINGS", Phase: "15M", MovePercent: 5.2, RelativeVolume: 2.1, VWAPState: "ABOVE", UpdatedAt: now}}
	macro := []EventReaction{{EventID: "cpi", OffsetSec: 300, CapturedAt: now, Moves: map[string]float64{"SPY": -0.8, "QQQ": -1.1}}}
	rows := buildReactionIntelligence(cats, macro)
	if len(rows) != 2 {
		t.Fatalf("expected ticker + macro reactions: %+v", rows)
	}
	if rows[0].UpdatedAt == 0 || rows[1].UpdatedAt == 0 {
		t.Fatalf("reaction timestamps missing: %+v", rows)
	}
}

func TestEventIntelligenceMacroMissingForecastCannotManufactureSurprise(t *testing.T) {
	now := time.Now()
	actual := 4.0
	rows := buildEconomicCalendar([]MacroEvent{{ID: "x", Name: "Release", Impact: "HIGH", Actual: &actual, Source: "Official", StartsAt: now.UnixMilli()}}, nil, now)
	if len(rows) != 1 || rows[0].Surprise != nil || rows[0].SurprisePct != nil {
		t.Fatalf("false macro surprise: %+v", rows)
	}
}

func TestEventIntelligenceStaleCalendarDegradesDecisionTruth(t *testing.T) {
	now := time.Now()
	d := buildEventDecisionCorrelation(nil, nil, nil, nil, EventModeState{}, map[string]int64{"macro-events": now.Add(-25 * time.Hour).UnixMilli()}, now)
	if d.MarketRisk != "DATA DEGRADED" || d.State != "DATA DEGRADED" {
		t.Fatalf("stale calendar falsely healthy: %+v", d)
	}
}

func TestEventIntelligenceOldNewsDoesNotBecomeSmartNotification(t *testing.T) {
	now := time.Now()
	n := []EventNewsIntelligence{{ID: "old", Headline: "old", Materiality: "HIGH", PublishedAt: now.Add(-3 * time.Hour).Unix(), Symbols: []string{"NVDA"}}}
	if rows := buildSmartNotifications(n, nil, nil, EventDecisionCorrelation{}, now); len(rows) != 0 {
		t.Fatalf("old news became new notification: %+v", rows)
	}
}

func TestEventIntelligenceReactionDoesNotClaimCausation(t *testing.T) {
	now := time.Now().UnixMilli()
	rows := buildReactionIntelligence(nil, []EventReaction{{EventID: "cpi", OffsetSec: 300, CapturedAt: now, Moves: map[string]float64{"SPY": -1}}})
	if len(rows) != 1 || !strings.Contains(strings.ToLower(rows[0].Detail), "no causal claim") {
		t.Fatalf("reaction semantics overclaim causation: %+v", rows)
	}
}

func TestEventIntelligenceSnapshotReusesCanonicalSourceHealth(t *testing.T) {
	now := time.Now()
	x := buildEventIntelligenceSnapshot(nil, nil, EventModeState{}, nil, nil, nil, map[string]int64{}, map[string]string{"news": "healthy · Finnhub", "macro-events": "healthy · official"}, now)
	if x.SourceHealth["news"] != "healthy · Finnhub" || x.SourceHealth["macroEvents"] != "healthy · official" {
		t.Fatalf("source health not canonical: %+v", x.SourceHealth)
	}
}

func TestEventIntelligenceWindowNotificationTimestampIsStable(t *testing.T) {
	now := time.Now()
	starts := now.Add(45 * time.Minute).UnixMilli()
	cal := []EconomicCalendarEntry{{ID: "fomc", Name: "FOMC", Impact: "HIGH", StartsAt: starts, Source: "Federal Reserve"}}
	a := buildSmartNotifications(nil, cal, nil, EventDecisionCorrelation{}, now)
	b := buildSmartNotifications(nil, cal, nil, EventDecisionCorrelation{}, now.Add(time.Minute))
	if len(a) != 1 || len(b) != 1 || a[0].CreatedAt != b[0].CreatedAt || a[0].CreatedAt != starts-60*60_000 {
		t.Fatalf("event-window notification churned: a=%+v b=%+v", a, b)
	}
}

func TestEventIntelligenceOldFedCalendarDoesNotLookRecent(t *testing.T) {
	now := time.Now()
	old := now.Add(-10 * 24 * time.Hour).UnixMilli()
	x := buildFedIntelligence([]EconomicCalendarEntry{{ID: "fed-old", Name: "FOMC Meeting", Source: "Federal Reserve", Date: now.Add(-10 * 24 * time.Hour).Format("2006-01-02"), StartsAt: old}}, now)
	if x.State != "HISTORICAL" {
		t.Fatalf("old Fed event falsely recent: %+v", x)
	}
}

func TestEventIntelligencePastUnknownTimeEventCannotLookScheduled(t *testing.T) {
	now := time.Date(2026, 8, 12, 15, 0, 0, 0, time.UTC)
	ev := MacroEvent{ID: "old", Name: "Historical release", Impact: "HIGH", Date: "2026-08-10", TimeKnown: false, Lifecycle: "RESOLVED", Source: "Official"}
	rows := buildEconomicCalendar([]MacroEvent{ev}, nil, now)
	if len(rows) != 1 || rows[0].State != "HISTORICAL" {
		t.Fatalf("past unknown-time event falsely scheduled: %+v", rows)
	}
}

func TestEventIntelligenceECBMinutesCannotBecomeFedIntelligence(t *testing.T) {
	now := time.Now()
	cal := []EconomicCalendarEntry{{ID: "ecb", Name: "ECB Meeting Minutes", Source: "European Central Bank", StartsAt: now.Add(time.Hour).UnixMilli(), State: "UPCOMING"}}
	x := buildFedIntelligence(cal, now)
	if x.State != "UNAVAILABLE" {
		t.Fatalf("non-Fed minutes leaked into Fed intelligence: %+v", x)
	}
}

func TestEventIntelligenceFutureSkewedMacroHealthDegradesDecision(t *testing.T) {
	now := time.Now()
	d := buildEventDecisionCorrelation(nil, nil, nil, nil, EventModeState{}, map[string]int64{"macro-events": now.Add(30 * time.Minute).UnixMilli()}, now)
	if d.MarketRisk != "DATA DEGRADED" || !strings.Contains(strings.ToLower(strings.Join(d.Reasons, " ")), "future-skewed") {
		t.Fatalf("future-skewed macro health falsely trusted: %+v", d)
	}
}

func TestEventIntelligenceWindowNotificationPersistsThroughReactionWindow(t *testing.T) {
	start := time.Now().Add(-10 * time.Minute)
	cal := []EconomicCalendarEntry{{ID: "cpi", Name: "CPI", Impact: "HIGH", StartsAt: start.UnixMilli(), Source: "BLS"}}
	rows := buildSmartNotifications(nil, cal, nil, EventDecisionCorrelation{}, time.Now())
	if len(rows) != 1 || rows[0].ID != "window-cpi" || rows[0].CreatedAt != start.UnixMilli()-60*60_000 || rows[0].ExpiresAt != start.UnixMilli()+90*60_000 {
		t.Fatalf("event-window notification did not persist truthfully: %+v", rows)
	}
	if !strings.Contains(strings.ToLower(rows[0].Message), "reaction window active") {
		t.Fatalf("post-release notification message is misleading: %+v", rows[0])
	}
}

func TestEventIntelligenceEarningsDateDistanceIsCalendarBasedAcrossDST(t *testing.T) {
	loc := easternLocation()
	now := time.Date(2026, 3, 7, 12, 0, 0, 0, loc)
	d, ok := dateDistanceDays("2026-03-09", now)
	if !ok || d != 2 {
		t.Fatalf("DST distorted calendar-day distance: days=%d ok=%v", d, ok)
	}
}
