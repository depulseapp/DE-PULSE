package main

import (
	"strings"
	"testing"
	"time"
)

func TestV162ProfessionalMacroMissingForecastCannotManufactureSurprise(t *testing.T) {
	now := time.Now()
	actual := 4.0
	rows := buildEconomicCalendar([]MacroEvent{{ID: "x", Name: "Release", Impact: "HIGH", Actual: &actual, Source: "Official", StartsAt: now.UnixMilli()}}, nil, now)
	if len(rows) != 1 || rows[0].Surprise != nil || rows[0].SurprisePct != nil {
		t.Fatalf("false macro surprise: %+v", rows)
	}
}
func TestV162ProfessionalStaleCalendarDegradesDecisionTruth(t *testing.T) {
	now := time.Now()
	d := buildEventDecisionCorrelation(nil, nil, nil, nil, EventModeState{}, map[string]int64{"macro-events": now.Add(-25 * time.Hour).UnixMilli()}, now)
	if d.MarketRisk != "DATA DEGRADED" || d.State != "DATA DEGRADED" {
		t.Fatalf("stale calendar falsely healthy: %+v", d)
	}
}
func TestV162ProfessionalOldNewsDoesNotBecomeSmartNotification(t *testing.T) {
	now := time.Now()
	n := []EventNewsIntelligence{{ID: "old", Headline: "old", Materiality: "HIGH", PublishedAt: now.Add(-3 * time.Hour).Unix(), Symbols: []string{"NVDA"}}}
	if rows := buildSmartNotifications(n, nil, nil, EventDecisionCorrelation{}, now); len(rows) != 0 {
		t.Fatalf("old news became new notification: %+v", rows)
	}
}
func TestV162ProfessionalReactionDoesNotClaimCausation(t *testing.T) {
	now := time.Now().UnixMilli()
	rows := buildReactionIntelligence(nil, []EventReaction{{EventID: "cpi", OffsetSec: 300, CapturedAt: now, Moves: map[string]float64{"SPY": -1}}})
	if len(rows) != 1 || !strings.Contains(strings.ToLower(rows[0].Detail), "no causal claim") {
		t.Fatalf("reaction semantics overclaim causation: %+v", rows)
	}
}
func TestV162ProfessionalSnapshotReusesCanonicalSourceHealth(t *testing.T) {
	now := time.Now()
	x := buildEventIntelligenceSnapshot(nil, nil, EventModeState{}, nil, nil, nil, map[string]int64{}, map[string]string{"news": "healthy · Finnhub", "macro-events": "healthy · official"}, now)
	if x.SourceHealth["news"] != "healthy · Finnhub" || x.SourceHealth["macroEvents"] != "healthy · official" {
		t.Fatalf("source health not canonical: %+v", x.SourceHealth)
	}
}

func TestV162ProfessionalEventWindowNotificationTimestampIsStable(t *testing.T) {
	now := time.Now()
	starts := now.Add(45 * time.Minute).UnixMilli()
	cal := []EconomicCalendarEntry{{ID: "fomc", Name: "FOMC", Impact: "HIGH", StartsAt: starts, Source: "Federal Reserve"}}
	a := buildSmartNotifications(nil, cal, nil, EventDecisionCorrelation{}, now)
	b := buildSmartNotifications(nil, cal, nil, EventDecisionCorrelation{}, now.Add(time.Minute))
	if len(a) != 1 || len(b) != 1 || a[0].CreatedAt != b[0].CreatedAt || a[0].CreatedAt != starts-60*60_000 {
		t.Fatalf("event-window notification churned: a=%+v b=%+v", a, b)
	}
}
func TestV162ProfessionalOldFedCalendarDoesNotLookRecent(t *testing.T) {
	now := time.Now()
	old := now.Add(-10 * 24 * time.Hour).UnixMilli()
	x := buildFedIntelligence([]EconomicCalendarEntry{{ID: "fed-old", Name: "FOMC Meeting", Source: "Federal Reserve", Date: now.Add(-10 * 24 * time.Hour).Format("2006-01-02"), StartsAt: old}}, now)
	if x.State != "HISTORICAL" {
		t.Fatalf("old Fed event falsely recent: %+v", x)
	}
}

func TestV162ProfessionalPastUnknownTimeEventCannotLookScheduled(t *testing.T) {
	now := time.Date(2026, 8, 12, 15, 0, 0, 0, time.UTC)
	ev := MacroEvent{ID: "old", Name: "Historical release", Impact: "HIGH", Date: "2026-08-10", TimeKnown: false, Lifecycle: "RESOLVED", Source: "Official"}
	rows := buildEconomicCalendar([]MacroEvent{ev}, nil, now)
	if len(rows) != 1 || rows[0].State != "HISTORICAL" {
		t.Fatalf("past unknown-time event falsely scheduled: %+v", rows)
	}
}

func TestV162ProfessionalECBMinutesCannotBecomeFedIntelligence(t *testing.T) {
	now := time.Now()
	cal := []EconomicCalendarEntry{{ID: "ecb", Name: "ECB Meeting Minutes", Source: "European Central Bank", StartsAt: now.Add(time.Hour).UnixMilli(), State: "UPCOMING"}}
	x := buildFedIntelligence(cal, now)
	if x.State != "UNAVAILABLE" {
		t.Fatalf("non-Fed minutes leaked into Fed intelligence: %+v", x)
	}
}

func TestV162ProfessionalFutureSkewedMacroHealthDegradesDecision(t *testing.T) {
	now := time.Now()
	d := buildEventDecisionCorrelation(nil, nil, nil, nil, EventModeState{}, map[string]int64{"macro-events": now.Add(30 * time.Minute).UnixMilli()}, now)
	if d.MarketRisk != "DATA DEGRADED" || !strings.Contains(strings.ToLower(strings.Join(d.Reasons, " ")), "future-skewed") {
		t.Fatalf("future-skewed macro health falsely trusted: %+v", d)
	}
}

func TestV162ProfessionalEventWindowNotificationPersistsThroughReactionWindow(t *testing.T) {
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

func TestV162ProfessionalEarningsDateDistanceIsCalendarBasedAcrossDST(t *testing.T) {
	loc := easternLocation()
	now := time.Date(2026, 3, 7, 12, 0, 0, 0, loc)
	d, ok := dateDistanceDays("2026-03-09", now)
	if !ok || d != 2 {
		t.Fatalf("DST distorted calendar-day distance: days=%d ok=%v", d, ok)
	}
}
