package main

import (
	"encoding/json"
	"math"
	"strings"
	"testing"
	"time"
)

func v163Bar(at time.Time, o, h, l, c float64) Bar {
	return Bar{T: at.UTC().Unix(), O: o, H: h, L: l, C: c, V: 1_000_000}
}

func v163Snapshot(at time.Time) SignalSnapshot {
	return SignalSnapshot{
		ID:                 "AAA-swing-evidence-v163",
		Symbol:             "AAA",
		Horizon:            "swing",
		Timestamp:          at.UnixMilli(),
		Price:              100,
		Score:              72,
		Action:             "BUY",
		Readiness:          "READY",
		EvidenceSnapshotID: "evidence-v163",
		FormulaVersion:     validationFormulaVersion,
		EntryLow:           95,
		EntryHigh:          100,
		TargetLow:          110,
		TargetHigh:         115,
		Invalidation:       90,
	}
}

func TestV163OutcomeTargetBeforeEntryNeverCountsAsSuccess(t *testing.T) {
	t0 := time.Date(2026, 1, 2, 16, 0, 0, 0, time.UTC)
	x := v163Snapshot(t0)
	bars := map[string][]Bar{"daily": {
		v163Bar(t0.Add(24*time.Hour), 112, 120, 111, 118), // target region before entry; must not count
		v163Bar(t0.Add(48*time.Hour), 99, 100, 96, 98),    // entry touch
		v163Bar(t0.Add(72*time.Hour), 101, 108, 96, 105),  // no later target
	}}
	evaluateOneSignalSnapshotWithActions(&x, bars, nil, t0.Add(4*24*time.Hour))
	if x.OutcomeState == "TARGET_REACHED" {
		t.Fatalf("target before entry was incorrectly counted as success: %+v", x)
	}
	if !x.EntryTouched {
		t.Fatalf("expected later entry touch to be recorded: %+v", x)
	}
}

func TestV163OutcomeInvalidationBeforeTarget(t *testing.T) {
	t0 := time.Date(2026, 1, 2, 16, 0, 0, 0, time.UTC)
	x := v163Snapshot(t0)
	bars := map[string][]Bar{"daily": {
		v163Bar(t0.Add(24*time.Hour), 99, 100, 96, 98),
		v163Bar(t0.Add(48*time.Hour), 97, 106, 88, 91),
		v163Bar(t0.Add(72*time.Hour), 100, 112, 98, 111),
	}}
	evaluateOneSignalSnapshotWithActions(&x, bars, nil, t0.Add(4*24*time.Hour))
	if x.OutcomeState != "INVALIDATED" {
		t.Fatalf("expected INVALIDATED first, got %s (%s)", x.OutcomeState, x.OutcomeDetail)
	}
	if x.InvalidationAt == 0 || x.TargetTouchedAt != 0 {
		t.Fatalf("unexpected resolution timestamps: %+v", x)
	}
}

func TestV163OutcomeEntryThenTarget(t *testing.T) {
	t0 := time.Date(2026, 1, 2, 16, 0, 0, 0, time.UTC)
	x := v163Snapshot(t0)
	bars := map[string][]Bar{"daily": {
		v163Bar(t0.Add(24*time.Hour), 99, 100, 96, 98),
		v163Bar(t0.Add(48*time.Hour), 101, 112, 96, 111),
	}}
	evaluateOneSignalSnapshotWithActions(&x, bars, nil, t0.Add(3*24*time.Hour))
	if x.OutcomeState != "TARGET_REACHED" {
		t.Fatalf("expected TARGET_REACHED, got %s (%s)", x.OutcomeState, x.OutcomeDetail)
	}
	if x.TargetTouchedAt <= x.EntryTouchedAt {
		t.Fatalf("target must resolve after entry: %+v", x)
	}
}

func TestV163OutcomeEntryNeverTouched(t *testing.T) {
	t0 := time.Date(2026, 1, 2, 16, 0, 0, 0, time.UTC)
	x := v163Snapshot(t0)
	rows := make([]Bar, 10)
	for i := range rows {
		rows[i] = v163Bar(t0.Add(time.Duration(i+1)*24*time.Hour), 105, 109, 101, 106)
	}
	evaluateOneSignalSnapshotWithActions(&x, map[string][]Bar{"daily": rows}, nil, t0.Add(12*24*time.Hour))
	if x.OutcomeState != "NO_ENTRY" {
		t.Fatalf("expected NO_ENTRY after supported 10-session window, got %s", x.OutcomeState)
	}
}

func TestV163OutcomeMissingMaturedHistoryUnavailable(t *testing.T) {
	t0 := time.Date(2025, 1, 2, 16, 0, 0, 0, time.UTC)
	x := v163Snapshot(t0)
	evaluateOneSignalSnapshotWithActions(&x, map[string][]Bar{}, nil, t0.Add(40*24*time.Hour))
	if x.OutcomeState != "UNAVAILABLE" {
		t.Fatalf("old snapshot without post-snapshot canonical bars must be UNAVAILABLE, got %s", x.OutcomeState)
	}
}

func TestV163OutcomeCoarseBarOrderingAmbiguous(t *testing.T) {
	t0 := time.Date(2026, 1, 2, 16, 0, 0, 0, time.UTC)
	x := v163Snapshot(t0)
	bars := map[string][]Bar{"daily": {
		v163Bar(t0.Add(24*time.Hour), 99, 100, 96, 98),
		v163Bar(t0.Add(48*time.Hour), 100, 112, 88, 101),
	}}
	evaluateOneSignalSnapshotWithActions(&x, bars, nil, t0.Add(3*24*time.Hour))
	if x.OutcomeState != "AMBIGUOUS" {
		t.Fatalf("target+invalidation inside one coarse post-entry bar must be AMBIGUOUS, got %s", x.OutcomeState)
	}
}

func TestV163LegacySnapshotRemainsPartial(t *testing.T) {
	t0 := time.Date(2025, 1, 2, 16, 0, 0, 0, time.UTC)
	x := SignalSnapshot{Symbol: "AAA", Horizon: "swing", Timestamp: t0.UnixMilli(), Price: 100}
	evaluateOneSignalSnapshotWithActions(&x, map[string][]Bar{"daily": {v163Bar(t0.Add(24*time.Hour), 100, 103, 99, 102)}}, nil, t0.Add(48*time.Hour))
	if !x.LegacyPartial || x.OutcomeState != "PARTIAL" {
		t.Fatalf("legacy snapshot must remain explicit PARTIAL, got %+v", x)
	}
}

func TestV163FrozenReplayFieldsSurviveAPIRoundTrip(t *testing.T) {
	days := 4
	in := v163Snapshot(time.Now().Add(-time.Hour))
	in.FamilyScores = map[string]float64{"Trend": 80, "Momentum": 70, "Participation": 75, "Structure": 72, "RelativeStrength": 68, "Market": 65}
	in.EarningsDays = &days
	in.SettingsFingerprint = "settings-abc"
	raw, err := json.Marshal(in)
	if err != nil {
		t.Fatal(err)
	}
	var out SignalSnapshot
	if err := json.Unmarshal(raw, &out); err != nil {
		t.Fatal(err)
	}
	if out.SettingsFingerprint != in.SettingsFingerprint || out.EarningsDays == nil || *out.EarningsDays != days || len(out.FamilyScores) != len(in.FamilyScores) {
		t.Fatalf("frozen replay identity/features were lost: in=%+v out=%+v", in, out)
	}
}

func TestV163SeasonalityInsufficientSampleDoesNotFabricateSignal(t *testing.T) {
	start := time.Date(2026, 1, 2, 0, 0, 0, 0, time.UTC)
	rows := make([]Bar, 200)
	p := 100.0
	for i := range rows {
		p *= 1.0005
		rows[i] = v163Bar(start.Add(time.Duration(i)*24*time.Hour), p, p*1.01, p*.99, p)
	}
	out := buildSeasonalitySnapshot(map[string]map[string][]Bar{"SPY": {"daily": rows}, "QQQ": {"daily": rows}}, start.Add(201*24*time.Hour))
	for _, sym := range []string{"SPY", "QQQ"} {
		st := out.Symbols[sym]
		if st.State != "INSUFFICIENT" {
			t.Fatalf("%s must degrade with short history, got %s", sym, st.State)
		}
		for _, m := range st.DayOfWeek {
			if m.State == "AVAILABLE" || m.PositiveFrequencyPct != nil {
				t.Fatalf("%s weekday sample fabricated descriptive confidence: %+v", sym, m)
			}
		}
	}
}

func TestV163CalibrationTinySampleAndScoreNotProbability(t *testing.T) {
	t0 := time.Date(2025, 1, 1, 16, 0, 0, 0, time.UTC)
	s := SignalValidationState{}
	for i := 0; i < 12; i++ {
		x := v163Snapshot(t0.Add(time.Duration(i) * 24 * time.Hour))
		x.Score = 82
		x.OutcomeState = "TARGET_REACHED"
		x.Outcomes = map[string]float64{"5D": 3.5}
		s.Snapshots = append(s.Snapshots, x)
	}
	out := buildCalibrationSnapshot(s, time.Now())
	if out.SetupScoreIsWinProbability {
		t.Fatal("Setup Score must never become win probability")
	}
	if len(out.Groups) != 1 || out.Groups[0].State != "INSUFFICIENT" {
		t.Fatalf("tiny sample must be INSUFFICIENT: %+v", out.Groups)
	}
	if strings.Contains(strings.ToLower(out.Groups[0].Detail), "82% chance") {
		t.Fatalf("probability claim leaked into calibration detail: %s", out.Groups[0].Detail)
	}
}

func v163CorrelatedDaily(start time.Time, n int, scale float64) []Bar {
	rows := make([]Bar, 0, n+1)
	p := 100.0 * scale
	for i := 0; i <= n; i++ {
		if i > 0 {
			move := (float64((i % 7) - 3)) / 1000
			p *= 1 + move
		}
		rows = append(rows, v163Bar(start.Add(time.Duration(i)*24*time.Hour), p, p*1.005, p*.995, p))
	}
	return rows
}

func TestV163CorrelationRequiresAlignedReturnWindow(t *testing.T) {
	start := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)
	a, _, _ := dailyReturnMap(v163CorrelatedDaily(start, 80, 1), start.Add(90*24*time.Hour).UnixMilli())
	b, _, _ := dailyReturnMap(v163CorrelatedDaily(start.Add(30*24*time.Hour), 50, 2), start.Add(90*24*time.Hour).UnixMilli())
	_, n := pearsonAligned(a, b)
	if n >= 60 {
		t.Fatalf("misaligned windows must not be treated as 60+ aligned samples, got %d", n)
	}
}

func TestV163SemiconductorConcentrationIsAttentionContextOnly(t *testing.T) {
	start := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)
	state := AppState{
		Settings:   Settings{DayEnabled: true, DayWatchlistID: "day"},
		Watchlists: []Watchlist{{ID: "day", Name: "Day", Symbols: []string{"NVDA", "AMD", "AVGO"}}},
		UI:         UIState{SelectedTicker: "NVDA"},
	}
	bars := map[string]map[string][]Bar{
		"NVDA": {"daily": v163CorrelatedDaily(start, 140, 1)},
		"AMD":  {"daily": v163CorrelatedDaily(start, 140, 2)},
		"AVGO": {"daily": v163CorrelatedDaily(start, 140, 3)},
	}
	out := buildCorrelationConcentrationSnapshot(state, ScannerState{}, bars, start.Add(150*24*time.Hour))
	foundIndustry := false
	foundHigh := false
	for _, g := range out.Concentrations {
		if g.Kind == "INDUSTRY" && g.Key == "Semiconductors" && len(g.Symbols) >= 2 {
			foundIndustry = true
			if !strings.Contains(strings.ToLower(g.Detail), "attention") {
				t.Fatalf("concentration must remain attention context: %+v", g)
			}
		}
	}
	for _, p := range out.Pairs {
		if p.SampleCount >= 126 && (p.State == "HIGH" || p.State == "ELEVATED") {
			foundHigh = true
			if math.Abs(p.Correlation) < .70 {
				t.Fatalf("elevated correlation below threshold: %+v", p)
			}
		}
	}
	if !foundIndustry || !foundHigh {
		t.Fatalf("expected semiconductor concentration and correlated-pair context, got %+v", out)
	}
}

func TestV163CanonicalSeasonalityIgnoresDailyRawComparisonPath(t *testing.T) {
	start := time.Date(2016, 1, 1, 0, 0, 0, 0, time.UTC)
	adjusted := v163CorrelatedDaily(start, 2600, 1)
	raw := v163CorrelatedDaily(start, 2600, 100)
	// Make raw comparison history intentionally absurd; canonical analytics must still use daily.
	for i := range raw {
		raw[i].C *= 1000
		raw[i].H *= 1000
		raw[i].L *= 1000
	}
	out := buildSeasonalitySnapshot(map[string]map[string][]Bar{
		"SPY": {"daily": adjusted, "daily-raw": raw},
		"QQQ": {"daily": adjusted, "daily-raw": raw},
	}, start.Add(2700*24*time.Hour))
	if out.Symbols["SPY"].State != "AVAILABLE" {
		t.Fatalf("expected deep canonical adjusted history to be AVAILABLE: %+v", out.Symbols["SPY"])
	}
	for _, m := range out.Symbols["SPY"].Monthly {
		if m.AverageReturnPct != nil && math.Abs(*m.AverageReturnPct) > 100 {
			t.Fatalf("daily-raw comparison path leaked into canonical seasonality: %+v", m)
		}
	}
}

func TestV163OutcomeIgnoresInProgressDailyBar(t *testing.T) {
	t0 := time.Date(2026, 1, 2, 0, 0, 0, 0, time.UTC)
	x := v163Snapshot(t0)
	// The first post-snapshot bar is completed and stays above the Entry Zone.
	// The second bar is still in progress at `now` and would touch entry/target if leaked.
	bars := map[string][]Bar{"daily": {
		v163Bar(t0.Add(24*time.Hour), 105, 108, 103, 106),
		v163Bar(t0.Add(48*time.Hour), 99, 112, 96, 111),
	}}
	now := t0.Add(48*time.Hour + 12*time.Hour)
	evaluateOneSignalSnapshotWithActions(&x, bars, nil, now)
	if x.EntryTouched || x.TargetTouchedAt != 0 {
		t.Fatalf("in-progress current daily bar leaked into outcome evidence: %+v", x)
	}
	if x.OutcomeState == "TARGET_REACHED" || x.OutcomeState == "AMBIGUOUS" {
		t.Fatalf("in-progress bar must not resolve an outcome, got %s (%s)", x.OutcomeState, x.OutcomeDetail)
	}
}

func TestV163DiscoveryCandidatesParticipateInConcentration(t *testing.T) {
	start := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)
	state := AppState{UI: UIState{SelectedTicker: "NVDA"}}
	scanner := ScannerState{Results: []ScannerResult{
		{Symbol: "AMD", Score: 91},
		{Symbol: "AVGO", Score: 88},
	}}
	bars := map[string]map[string][]Bar{
		"NVDA": {"daily": v163CorrelatedDaily(start, 140, 1)},
		"AMD":  {"daily": v163CorrelatedDaily(start, 140, 2)},
		"AVGO": {"daily": v163CorrelatedDaily(start, 140, 3)},
	}
	out := buildCorrelationConcentrationSnapshot(state, scanner, bars, start.Add(150*24*time.Hour))
	foundIndustry := false
	for _, g := range out.Concentrations {
		if g.Kind == "INDUSTRY" && g.Key == "Semiconductors" {
			members := map[string]bool{}
			for _, sym := range g.Symbols {
				members[sym] = true
			}
			if members["NVDA"] && members["AMD"] && members["AVGO"] {
				foundIndustry = true
				break
			}
		}
	}
	if !foundIndustry {
		t.Fatalf("Discovery/scanner candidates were not included in semiconductor concentration context: %+v", out.Concentrations)
	}
}
