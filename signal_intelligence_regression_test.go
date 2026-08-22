package main

import (
	"encoding/json"
	"math"
	"strings"
	"testing"
	"time"
)

func signalEquivalenceBar(at time.Time, o, h, l, c float64) Bar {
	return Bar{T: at.UTC().Unix(), O: o, H: h, L: l, C: c, V: 1_000_000}
}

func signalEquivalenceSnapshot(at time.Time) SignalSnapshot {
	return SignalSnapshot{
		ID:                 "AAA-swing-equivalence",
		Symbol:             "AAA",
		Horizon:            "swing",
		Timestamp:          at.UnixMilli(),
		Price:              100,
		Score:              72,
		Action:             "BUY",
		Readiness:          "READY",
		EvidenceSnapshotID: "evidence-equivalence",
		FormulaVersion:     validationFormulaVersion,
		EntryLow:           95,
		EntryHigh:          100,
		TargetLow:          110,
		TargetHigh:         115,
		Invalidation:       90,
	}
}

func signalProfessionalSnapshot(ts int64) SignalSnapshot {
	return SignalSnapshot{
		ID:     "AAA-swing-evidence-1-" + validationFormulaVersion,
		Symbol: "AAA", Horizon: "swing", Timestamp: ts, Price: 100,
		Score: 72, Action: "APPROACHING ENTRY", Readiness: "READY",
		EvidenceSnapshotID: "evidence-1", FormulaVersion: validationFormulaVersion,
		SettingsFingerprint: "abc123", EarningsPenalty: 10, SignalProfile: "balanced",
		FamilyScores: map[string]float64{"Trend": 70, "Momentum": 70, "Participation": 70, "Structure": 70, "RelativeStrength": 70, "Market": 70},
		EntryLow: 98, EntryHigh: 100, TargetLow: 110, TargetHigh: 115, Invalidation: 94,
	}
}

func signalCorrelatedDaily(start time.Time, n int, scale float64) []Bar {
	rows := make([]Bar, 0, n+1)
	p := 100.0 * scale
	for i := 0; i <= n; i++ {
		if i > 0 {
			move := float64((i%7)-3) / 1000
			p *= 1 + move
		}
		rows = append(rows, signalEquivalenceBar(start.Add(time.Duration(i)*24*time.Hour), p, p*1.005, p*.995, p))
	}
	return rows
}

func TestSignalValidationOutcomeTargetBeforeEntryNeverCountsAsSuccess(t *testing.T) {
	t0 := time.Date(2026, 1, 2, 16, 0, 0, 0, time.UTC)
	x := signalEquivalenceSnapshot(t0)
	bars := map[string][]Bar{"daily": {
		signalEquivalenceBar(t0.Add(24*time.Hour), 112, 120, 111, 118),
		signalEquivalenceBar(t0.Add(48*time.Hour), 99, 100, 96, 98),
		signalEquivalenceBar(t0.Add(72*time.Hour), 101, 108, 96, 105),
	}}
	evaluateOneSignalSnapshotWithActions(&x, bars, nil, t0.Add(4*24*time.Hour))
	if x.OutcomeState == "TARGET_REACHED" {
		t.Fatalf("target before entry was incorrectly counted as success: %+v", x)
	}
	if !x.EntryTouched {
		t.Fatalf("expected later entry touch to be recorded: %+v", x)
	}
}

func TestSignalValidationOutcomeInvalidationBeforeTarget(t *testing.T) {
	t0 := time.Date(2026, 1, 2, 16, 0, 0, 0, time.UTC)
	x := signalEquivalenceSnapshot(t0)
	bars := map[string][]Bar{"daily": {
		signalEquivalenceBar(t0.Add(24*time.Hour), 99, 100, 96, 98),
		signalEquivalenceBar(t0.Add(48*time.Hour), 97, 106, 88, 91),
		signalEquivalenceBar(t0.Add(72*time.Hour), 100, 112, 98, 111),
	}}
	evaluateOneSignalSnapshotWithActions(&x, bars, nil, t0.Add(4*24*time.Hour))
	if x.OutcomeState != "INVALIDATED" {
		t.Fatalf("expected INVALIDATED first, got %s (%s)", x.OutcomeState, x.OutcomeDetail)
	}
	if x.InvalidationAt == 0 || x.TargetTouchedAt != 0 {
		t.Fatalf("unexpected resolution timestamps: %+v", x)
	}
}

func TestSignalValidationOutcomeEntryThenTarget(t *testing.T) {
	t0 := time.Date(2026, 1, 2, 16, 0, 0, 0, time.UTC)
	x := signalEquivalenceSnapshot(t0)
	bars := map[string][]Bar{"daily": {
		signalEquivalenceBar(t0.Add(24*time.Hour), 99, 100, 96, 98),
		signalEquivalenceBar(t0.Add(48*time.Hour), 101, 112, 96, 111),
	}}
	evaluateOneSignalSnapshotWithActions(&x, bars, nil, t0.Add(3*24*time.Hour))
	if x.OutcomeState != "TARGET_REACHED" {
		t.Fatalf("expected TARGET_REACHED, got %s (%s)", x.OutcomeState, x.OutcomeDetail)
	}
	if x.TargetTouchedAt <= x.EntryTouchedAt {
		t.Fatalf("target must resolve after entry: %+v", x)
	}
}

func TestSignalValidationOutcomeEntryNeverTouched(t *testing.T) {
	t0 := time.Date(2026, 1, 2, 16, 0, 0, 0, time.UTC)
	x := signalEquivalenceSnapshot(t0)
	rows := make([]Bar, 10)
	for i := range rows {
		rows[i] = signalEquivalenceBar(t0.Add(time.Duration(i+1)*24*time.Hour), 105, 109, 101, 106)
	}
	evaluateOneSignalSnapshotWithActions(&x, map[string][]Bar{"daily": rows}, nil, t0.Add(12*24*time.Hour))
	if x.OutcomeState != "NO_ENTRY" {
		t.Fatalf("expected NO_ENTRY after supported 10-session window, got %s", x.OutcomeState)
	}
}

func TestSignalValidationOutcomeMissingMaturedHistoryUnavailable(t *testing.T) {
	t0 := time.Date(2025, 1, 2, 16, 0, 0, 0, time.UTC)
	x := signalEquivalenceSnapshot(t0)
	evaluateOneSignalSnapshotWithActions(&x, map[string][]Bar{}, nil, t0.Add(40*24*time.Hour))
	if x.OutcomeState != "UNAVAILABLE" {
		t.Fatalf("old snapshot without post-snapshot canonical bars must be UNAVAILABLE, got %s", x.OutcomeState)
	}
}

func TestSignalValidationOutcomeCoarseBarOrderingAmbiguous(t *testing.T) {
	t0 := time.Date(2026, 1, 2, 16, 0, 0, 0, time.UTC)
	x := signalEquivalenceSnapshot(t0)
	bars := map[string][]Bar{"daily": {
		signalEquivalenceBar(t0.Add(24*time.Hour), 99, 100, 96, 98),
		signalEquivalenceBar(t0.Add(48*time.Hour), 100, 112, 88, 101),
	}}
	evaluateOneSignalSnapshotWithActions(&x, bars, nil, t0.Add(3*24*time.Hour))
	if x.OutcomeState != "AMBIGUOUS" {
		t.Fatalf("target+invalidation inside one coarse post-entry bar must be AMBIGUOUS, got %s", x.OutcomeState)
	}
}

func TestSignalValidationLegacySnapshotRemainsPartial(t *testing.T) {
	t0 := time.Date(2025, 1, 2, 16, 0, 0, 0, time.UTC)
	x := SignalSnapshot{Symbol: "AAA", Horizon: "swing", Timestamp: t0.UnixMilli(), Price: 100}
	evaluateOneSignalSnapshotWithActions(&x, map[string][]Bar{"daily": {signalEquivalenceBar(t0.Add(24*time.Hour), 100, 103, 99, 102)}}, nil, t0.Add(48*time.Hour))
	if !x.LegacyPartial || x.OutcomeState != "PARTIAL" {
		t.Fatalf("legacy snapshot must remain explicit PARTIAL, got %+v", x)
	}
}

func TestSignalValidationFrozenReplayFieldsSurviveRoundTrip(t *testing.T) {
	days := 4
	in := signalEquivalenceSnapshot(time.Now().Add(-time.Hour))
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

func TestSignalValidationSeasonalityInsufficientSampleDoesNotFabricateSignal(t *testing.T) {
	start := time.Date(2026, 1, 2, 0, 0, 0, 0, time.UTC)
	rows := make([]Bar, 200)
	p := 100.0
	for i := range rows {
		p *= 1.0005
		rows[i] = signalEquivalenceBar(start.Add(time.Duration(i)*24*time.Hour), p, p*1.01, p*.99, p)
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

func TestSignalValidationCalibrationTinySampleAndScoreNotProbability(t *testing.T) {
	t0 := time.Date(2025, 1, 1, 16, 0, 0, 0, time.UTC)
	s := SignalValidationState{}
	for i := 0; i < 12; i++ {
		x := signalEquivalenceSnapshot(t0.Add(time.Duration(i) * 24 * time.Hour))
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

func TestSignalValidationCorrelationRequiresAlignedReturnWindow(t *testing.T) {
	start := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)
	a, _, _ := dailyReturnMap(signalCorrelatedDaily(start, 80, 1), start.Add(90*24*time.Hour).UnixMilli())
	b, _, _ := dailyReturnMap(signalCorrelatedDaily(start.Add(30*24*time.Hour), 50, 2), start.Add(90*24*time.Hour).UnixMilli())
	_, n := pearsonAligned(a, b)
	if n >= 60 {
		t.Fatalf("misaligned windows must not be treated as 60+ aligned samples, got %d", n)
	}
}

func TestSignalValidationSemiconductorConcentrationIsAttentionContextOnly(t *testing.T) {
	start := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)
	state := AppState{Settings: Settings{DayEnabled: true, DayWatchlistID: "day"}, Watchlists: []Watchlist{{ID: "day", Name: "Day", Symbols: []string{"NVDA", "AMD", "AVGO"}}}, UI: UIState{SelectedTicker: "NVDA"}}
	bars := map[string]map[string][]Bar{"NVDA": {"daily": signalCorrelatedDaily(start, 140, 1)}, "AMD": {"daily": signalCorrelatedDaily(start, 140, 2)}, "AVGO": {"daily": signalCorrelatedDaily(start, 140, 3)}}
	out := buildCorrelationConcentrationSnapshot(state, ScannerState{}, bars, start.Add(150*24*time.Hour))
	foundIndustry, foundHigh := false, false
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

func TestSignalValidationCanonicalSeasonalityIgnoresRawComparisonPath(t *testing.T) {
	start := time.Date(2016, 1, 1, 0, 0, 0, 0, time.UTC)
	adjusted := signalCorrelatedDaily(start, 2600, 1)
	raw := signalCorrelatedDaily(start, 2600, 100)
	for i := range raw {
		raw[i].C *= 1000
		raw[i].H *= 1000
		raw[i].L *= 1000
	}
	out := buildSeasonalitySnapshot(map[string]map[string][]Bar{"SPY": {"daily": adjusted, "daily-raw": raw}, "QQQ": {"daily": adjusted, "daily-raw": raw}}, start.Add(2700*24*time.Hour))
	if out.Symbols["SPY"].State != "AVAILABLE" {
		t.Fatalf("expected deep canonical adjusted history to be AVAILABLE: %+v", out.Symbols["SPY"])
	}
	for _, m := range out.Symbols["SPY"].Monthly {
		if m.AverageReturnPct != nil && math.Abs(*m.AverageReturnPct) > 100 {
			t.Fatalf("daily-raw comparison path leaked into canonical seasonality: %+v", m)
		}
	}
}

func TestSignalValidationOutcomeIgnoresInProgressDailyBar(t *testing.T) {
	t0 := time.Date(2026, 1, 2, 0, 0, 0, 0, time.UTC)
	x := signalEquivalenceSnapshot(t0)
	bars := map[string][]Bar{"daily": {
		signalEquivalenceBar(t0.Add(24*time.Hour), 105, 108, 103, 106),
		signalEquivalenceBar(t0.Add(48*time.Hour), 99, 112, 96, 111),
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

func TestSignalValidationDiscoveryCandidatesParticipateInConcentration(t *testing.T) {
	start := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)
	state := AppState{UI: UIState{SelectedTicker: "NVDA"}}
	scanner := ScannerState{Results: []ScannerResult{{Symbol: "AMD", Score: 91}, {Symbol: "AVGO", Score: 88}}}
	bars := map[string]map[string][]Bar{"NVDA": {"daily": signalCorrelatedDaily(start, 140, 1)}, "AMD": {"daily": signalCorrelatedDaily(start, 140, 2)}, "AVGO": {"daily": signalCorrelatedDaily(start, 140, 3)}}
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

func TestSignalValidationProfessionalTargetBeforeEntryNeverCountsAsSuccess(t *testing.T) {
	base := time.Date(2026, 1, 2, 20, 0, 0, 0, time.UTC)
	x := signalProfessionalSnapshot(base.UnixMilli())
	bars := []Bar{
		signalEquivalenceBar(base.Add(24*time.Hour), 101, 112, 101, 108),
		signalEquivalenceBar(base.Add(48*time.Hour), 100, 101, 98, 99),
		signalEquivalenceBar(base.Add(72*time.Hour), 99, 111, 98, 110),
	}
	evaluateOneSignalSnapshotWithActions(&x, map[string][]Bar{"daily": bars}, nil, base.Add(96*time.Hour))
	if x.OutcomeState != "TARGET_REACHED" {
		t.Fatalf("expected post-entry target success, got %s (%s)", x.OutcomeState, x.OutcomeDetail)
	}
	if x.TargetTouchedAt != normalizedBarTimestampMs(bars[2].T) || x.EntryTouchedAt != normalizedBarTimestampMs(bars[1].T) {
		t.Fatalf("target/entry timestamp ordering mismatch: %+v", x)
	}
}

func TestSignalValidationProfessionalInvalidationBeforeTarget(t *testing.T) {
	base := time.Date(2026, 2, 2, 20, 0, 0, 0, time.UTC)
	x := signalProfessionalSnapshot(base.UnixMilli())
	bars := []Bar{
		signalEquivalenceBar(base.Add(24*time.Hour), 100, 103, 98, 101),
		signalEquivalenceBar(base.Add(48*time.Hour), 99, 102, 93, 95),
		signalEquivalenceBar(base.Add(72*time.Hour), 96, 112, 95, 111),
	}
	evaluateOneSignalSnapshotWithActions(&x, map[string][]Bar{"daily": bars}, nil, base.Add(96*time.Hour))
	if x.OutcomeState != "INVALIDATED" || x.TargetTouchedAt != 0 || x.InvalidationAt == 0 {
		t.Fatalf("expected invalidation before target, got %+v", x)
	}
}

func TestSignalValidationProfessionalEntryThenTargetAndExcursions(t *testing.T) {
	base := time.Date(2026, 3, 2, 20, 0, 0, 0, time.UTC)
	x := signalProfessionalSnapshot(base.UnixMilli())
	bars := []Bar{signalEquivalenceBar(base.Add(24*time.Hour), 100, 103, 98.5, 101), signalEquivalenceBar(base.Add(48*time.Hour), 101, 111, 97, 110)}
	evaluateOneSignalSnapshotWithActions(&x, map[string][]Bar{"daily": bars}, nil, base.Add(72*time.Hour))
	if x.OutcomeState != "TARGET_REACHED" || !x.EntryTouched || x.MFE <= 0 || x.MAE >= 0 || x.ElapsedMinutes <= 0 {
		t.Fatalf("expected entry then target with excursions/elapsed evidence, got %+v", x)
	}
}

func TestSignalValidationProfessionalFirstEntryBarAmbiguityIsNotGuessed(t *testing.T) {
	base := time.Date(2026, 4, 2, 20, 0, 0, 0, time.UTC)
	x := signalProfessionalSnapshot(base.UnixMilli())
	bar := signalEquivalenceBar(base.Add(24*time.Hour), 100, 111, 98, 109)
	evaluateOneSignalSnapshotWithActions(&x, map[string][]Bar{"daily": {bar}}, nil, base.Add(48*time.Hour))
	if x.OutcomeState != "AMBIGUOUS" {
		t.Fatalf("same-bar entry/target ordering must be ambiguous, got %s", x.OutcomeState)
	}
}

func TestSignalValidationProfessionalNoFutureBarsDegradesTruthfully(t *testing.T) {
	base := time.Now().Add(-30 * 24 * time.Hour)
	x := signalProfessionalSnapshot(base.UnixMilli())
	evaluateOneSignalSnapshotWithActions(&x, map[string][]Bar{"daily": nil}, nil, time.Now())
	if x.OutcomeState != "UNAVAILABLE" {
		t.Fatalf("matured snapshot with no future evidence must be UNAVAILABLE, got %s", x.OutcomeState)
	}
}

func TestSignalValidationProfessionalNoEntryUsesSupportedTenSessionWindow(t *testing.T) {
	base := time.Date(2026, 5, 1, 20, 0, 0, 0, time.UTC)
	x := signalProfessionalSnapshot(base.UnixMilli())
	bars := make([]Bar, 0, 12)
	for i := 1; i <= 10; i++ {
		bars = append(bars, signalEquivalenceBar(base.Add(time.Duration(i)*24*time.Hour), 104, 108, 102, 105))
	}
	bars = append(bars, signalEquivalenceBar(base.Add(11*24*time.Hour), 100, 102, 98, 99))
	evaluateOneSignalSnapshotWithActions(&x, map[string][]Bar{"daily": bars}, nil, base.Add(12*24*time.Hour))
	if x.OutcomeState != "NO_ENTRY" || x.EntryTouched {
		t.Fatalf("late entry after supported window must not rewrite NO_ENTRY: %+v", x)
	}
}

func TestSignalValidationProfessionalWeekendGapDoesNotDelayCompletedDailyEvidence(t *testing.T) {
	friday := time.Date(2026, 5, 8, 0, 0, 0, 0, time.UTC)
	monday := friday.Add(72 * time.Hour)
	rows := []Bar{signalEquivalenceBar(friday, 100, 103, 99, 102), signalEquivalenceBar(monday, 103, 105, 102, 104)}
	cutoff := friday.Add(36 * time.Hour)
	got := completedBarsBefore(rows, cutoff.UnixMilli(), "daily")
	if len(got) != 1 || normalizedBarTimestampMs(got[0].T) != friday.UnixMilli() {
		t.Fatalf("weekend gap delayed completed Friday evidence: %+v", got)
	}
}

func TestSignalValidationProfessionalSeasonalityUsesCanonicalAdjustedDailyOnly(t *testing.T) {
	start := time.Date(2015, 1, 2, 20, 0, 0, 0, time.UTC)
	daily := make([]Bar, 0, 420)
	price := 100.0
	for i := 0; i < 420; i++ {
		price *= 1.0005
		daily = append(daily, signalEquivalenceBar(start.Add(time.Duration(i)*24*time.Hour), price*.99, price*1.01, price*.98, price))
	}
	rawWithSplitGap := append([]Bar(nil), daily...)
	rawWithSplitGap[len(rawWithSplitGap)/2].C /= 10
	bars := map[string]map[string][]Bar{"SPY": {"daily": daily, "daily-raw": rawWithSplitGap}, "QQQ": {"daily": daily, "daily-raw": rawWithSplitGap}}
	got := buildSeasonalitySnapshot(bars, start.Add(500*24*time.Hour))
	if got.Symbols["SPY"].DailyBars != len(daily) || !strings.Contains(got.Symbols["SPY"].Source, "Canonical daily history") {
		t.Fatalf("seasonality must use canonical adjusted daily history: %+v", got.Symbols["SPY"])
	}
}

func TestSignalValidationProfessionalSeasonalityInsufficientSampleDoesNotFabricateStatistics(t *testing.T) {
	now := time.Now()
	bars := map[string]map[string][]Bar{"SPY": {"daily": {{T: now.Add(-48 * time.Hour).Unix(), C: 100}, {T: now.Add(-24 * time.Hour).Unix(), C: 101}}}}
	got := buildSeasonalitySnapshot(bars, now)
	spy := got.Symbols["SPY"]
	if spy.State != "INSUFFICIENT" {
		t.Fatalf("small sample must be INSUFFICIENT, got %s", spy.State)
	}
	for _, m := range append(append([]SeasonalityMetric{}, spy.Monthly...), spy.DayOfWeek...) {
		if m.State == "AVAILABLE" || m.AverageReturnPct != nil || m.PositiveFrequencyPct != nil {
			t.Fatalf("insufficient sample fabricated statistics: %+v", m)
		}
	}
}

func TestSignalValidationProfessionalCalibrationTinySampleAndNoProbabilityClaim(t *testing.T) {
	base := time.Now().Add(-90 * 24 * time.Hour).UnixMilli()
	s := SignalValidationState{Snapshots: []SignalSnapshot{{Symbol: "AAA", Horizon: "swing", Timestamp: base, Score: 82, Action: "APPROACHING ENTRY", OutcomeState: "TARGET_REACHED", EntryTouched: true, MFE: 8, MAE: -3, Outcomes: map[string]float64{"5D": 4.2}}}}
	got := buildCalibrationSnapshot(s, time.Now())
	if got.SetupScoreIsWinProbability || len(got.Groups) != 1 || got.Groups[0].State != "INSUFFICIENT" || strings.Contains(got.Message, "82%") || strings.Contains(got.Groups[0].Detail, "82%") {
		t.Fatalf("tiny calibration sample must refuse probability confidence: %+v", got)
	}
}

func TestSignalValidationProfessionalCalibrationDoesNotTreatNoEntryAsZeroExcursion(t *testing.T) {
	base := time.Now().Add(-90 * 24 * time.Hour).UnixMilli()
	s := SignalValidationState{Snapshots: []SignalSnapshot{
		{Symbol: "AAA", Horizon: "swing", Timestamp: base, Score: 72, Action: "APPROACHING ENTRY", OutcomeState: "NO_ENTRY", EntryTouched: false},
		{Symbol: "BBB", Horizon: "swing", Timestamp: base + 1, Score: 72, Action: "APPROACHING ENTRY", OutcomeState: "TARGET_REACHED", EntryTouched: true, MFE: 10, MAE: -4},
	}}
	got := buildCalibrationSnapshot(s, time.Now())
	if len(got.Groups) != 1 || got.Groups[0].AverageMFE == nil || math.Abs(*got.Groups[0].AverageMFE-10) > 1e-9 || got.Groups[0].AverageMAE == nil || math.Abs(*got.Groups[0].AverageMAE-(-4)) > 1e-9 {
		t.Fatalf("NO_ENTRY must not inject synthetic zero excursions: %+v", got.Groups)
	}
}

func TestSignalValidationProfessionalCorrelationUsesAlignedReturnsAndFlagsConcentration(t *testing.T) {
	start := time.Date(2025, 1, 2, 20, 0, 0, 0, time.UTC)
	mk := func(scale float64, skipEvery int) []Bar {
		out := []Bar{}
		p := 100.0
		for i := 0; i < 90; i++ {
			if skipEvery > 0 && i%skipEvery == 0 {
				continue
			}
			r := 0.002 * math.Sin(float64(i)/4)
			p *= 1 + r*scale
			out = append(out, signalEquivalenceBar(start.Add(time.Duration(i)*24*time.Hour), p, p*1.01, p*.99, p))
		}
		return out
	}
	st := AppState{Settings: Settings{DayEnabled: true, DayWatchlistID: "day"}, Watchlists: []Watchlist{{ID: "day", Symbols: []string{"NVDA", "AMD", "AVGO"}}}}
	bars := map[string]map[string][]Bar{"NVDA": {"daily": mk(1, 0)}, "AMD": {"daily": mk(1.05, 0)}, "AVGO": {"daily": mk(.95, 11)}}
	got := buildCorrelationConcentrationSnapshot(st, ScannerState{}, bars, start.Add(120*24*time.Hour))
	foundPair, foundSector := false, false
	for _, p := range got.Pairs {
		if (p.SymbolA == "NVDA" && p.SymbolB == "AMD") || (p.SymbolA == "AMD" && p.SymbolB == "NVDA") {
			foundPair = true
			if p.SampleCount < 60 || p.State != "HIGH" {
				t.Fatalf("expected aligned high-correlation pair, got %+v", p)
			}
		}
	}
	for _, g := range got.Concentrations {
		if g.Kind == "SECTOR" && g.Key == "Technology" && len(g.Symbols) >= 3 {
			foundSector = true
		}
	}
	if !foundPair || !foundSector {
		t.Fatalf("aligned correlation/concentration evidence missing: %+v", got)
	}
}

func TestSignalValidationProfessionalSnapshotIdentityFreezesEvidenceAndSettings(t *testing.T) {
	e := &Engine{signalValidation: SignalValidationState{Snapshots: []SignalSnapshot{}}}
	base := signalProfessionalSnapshot(time.Now().UnixMilli())
	first := e.recordSignalSnapshot(base)
	changed := base
	changed.Timestamp += 60_000
	changed.Score = 10
	changed.Action = "INVALIDATED"
	second := e.recordSignalSnapshot(changed)
	if first.ID != second.ID || second.Score != first.Score || second.Action != first.Action {
		t.Fatalf("same evidence/settings/formula identity must be immutable; first=%+v second=%+v", first, second)
	}
}

func TestSignalValidationProfessionalSplitAdjustmentPreservesFrozenOutcomeEconomics(t *testing.T) {
	base := time.Date(2025, 6, 2, 20, 0, 0, 0, time.UTC)
	x := signalProfessionalSnapshot(base.UnixMilli())
	bars := []Bar{
		signalEquivalenceBar(base.Add(24*time.Hour), 9.9, 10.0, 9.8, 9.9),
		signalEquivalenceBar(base.Add(48*time.Hour), 10.1, 11.1, 9.7, 11.0),
	}
	actions := []CorporateAction{{Symbol: "AAA", Type: "forward_split", ExDate: base.Add(36 * time.Hour).Format("2006-01-02"), Ratio: 10, AdjustmentFactor: 10, Status: "EFFECTIVE", Source: "Alpaca"}}
	evaluateOneSignalSnapshotWithActions(&x, map[string][]Bar{"daily": bars}, actions, base.Add(72*time.Hour))
	if x.OutcomeState != "TARGET_REACHED" {
		t.Fatalf("post-split adjusted history must preserve frozen target economics, got %s (%s)", x.OutcomeState, x.OutcomeDetail)
	}
	if math.Abs(x.OutcomeAdjustmentFactor-10) > 1e-9 || !strings.Contains(x.OutcomeAdjustmentDetail, "split evidence") {
		t.Fatalf("split normalization provenance missing: factor=%f detail=%q", x.OutcomeAdjustmentFactor, x.OutcomeAdjustmentDetail)
	}
	if got := x.Outcomes["1D"]; math.Abs(got-(-1.0)) > 0.05 {
		t.Fatalf("forward return must use split-adjusted frozen price anchor, got %.4f%%", got)
	}
}
