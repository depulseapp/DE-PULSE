package main

import (
	"strings"
	"testing"
	"time"
)

func TestV162NewsClustersDuplicateSourcesAndRanksMateriality(t *testing.T) {
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

func TestV162NewsRejectsFutureSkew(t *testing.T) {
	now := time.Now()
	rows := buildEventNewsIntelligence([]NewsItem{{Datetime: now.Add(10 * time.Minute).Unix(), Headline: "Future headline", Source: "X"}}, now)
	if len(rows) != 0 {
		t.Fatalf("future-skewed news must not become event intelligence: %+v", rows)
	}
}

func TestV162MacroSurpriseRequiresActualAndForecast(t *testing.T) {
	now := time.Now()
	a, f, p := 3.1, 2.8, 2.9
	rows := buildEconomicCalendar([]MacroEvent{{ID: "cpi", Name: "Consumer Price Index", Impact: "HIGH", StartsAt: now.Add(-5 * time.Minute).UnixMilli(), Date: now.Format("2006-01-02"), TimeKnown: true, Actual: &a, Expected: &f, Previous: &p, Source: "BLS"}, {ID: "gdp", Name: "GDP", Impact: "HIGH", StartsAt: now.Add(time.Hour).UnixMilli(), Date: now.Format("2006-01-02"), TimeKnown: true, Source: "BEA"}}, nil, now)
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

func TestV162FedIntelligenceUsesSourcedTimeline(t *testing.T) {
	now := time.Now()
	date := now.Format("2006-01-02")
	cal := []EconomicCalendarEntry{{ID: "fed1", Name: "FOMC Meeting", Source: "Federal Reserve", Date: date, StartsAt: now.Add(30 * time.Minute).UnixMilli(), State: "UPCOMING"}, {ID: "fed2", Name: "Press Conference", Source: "Federal Reserve", Date: date, StartsAt: now.Add(90 * time.Minute).UnixMilli(), State: "UPCOMING"}}
	x := buildFedIntelligence(cal, now)
	if x.State != "UPCOMING" || len(x.Timeline) != 2 || x.CountdownSec <= 0 {
		t.Fatalf("fed intelligence wrong: %+v", x)
	}
}

func TestV162EventDecisionCorrelatesWithoutScoreMutation(t *testing.T) {
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

func TestV162SmartNotificationsAreEventTriggeredNotConditionCards(t *testing.T) {
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

func TestV162ReactionIntelligenceReusesCatalystAndMacroEvidence(t *testing.T) {
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
