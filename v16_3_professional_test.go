package main

import (
	"math"
	"strings"
	"testing"
	"time"
)

func v163DailyBar(day time.Time, o, h, l, c float64) Bar {
	return Bar{T: day.UTC().Unix(), O: o, H: h, L: l, C: c, V: 1_000_000}
}

func v163BaseSnapshot(ts int64) SignalSnapshot {
	return SignalSnapshot{
		ID:     "AAA-swing-evidence-1-" + validationFormulaVersion,
		Symbol: "AAA", Horizon: "swing", Timestamp: ts, Price: 100,
		Score: 72, Action: "APPROACHING ENTRY", Readiness: "READY",
		EvidenceSnapshotID: "evidence-1", FormulaVersion: validationFormulaVersion,
		SettingsFingerprint: "abc123", EarningsPenalty: 10, SignalProfile: "balanced",
		FamilyScores: map[string]float64{"Trend": 70, "Momentum": 70, "Participation": 70, "Structure": 70, "RelativeStrength": 70, "Market": 70},
		EntryLow:     98, EntryHigh: 100, TargetLow: 110, TargetHigh: 115, Invalidation: 94,
	}
}

func TestV163ProfessionalTargetBeforeEntryNeverCountsAsSuccess(t *testing.T) {
	base := time.Date(2026, 1, 2, 20, 0, 0, 0, time.UTC)
	x := v163BaseSnapshot(base.UnixMilli())
	bars := []Bar{
		v163DailyBar(base.Add(24*time.Hour), 101, 112, 101, 108), // target before any entry touch
		v163DailyBar(base.Add(48*time.Hour), 100, 101, 98, 99),   // entry later
		v163DailyBar(base.Add(72*time.Hour), 99, 111, 98, 110),   // target only after entry
	}
	evaluateOneSignalSnapshotWithActions(&x, map[string][]Bar{"daily": bars}, nil, base.Add(96*time.Hour))
	if x.OutcomeState != "TARGET_REACHED" {
		t.Fatalf("expected post-entry target success, got %s (%s)", x.OutcomeState, x.OutcomeDetail)
	}
	if x.TargetTouchedAt != normalizedBarTimestampMs(bars[2].T) {
		t.Fatalf("target-before-entry was incorrectly used; targetAt=%d want=%d", x.TargetTouchedAt, normalizedBarTimestampMs(bars[2].T))
	}
	if x.EntryTouchedAt != normalizedBarTimestampMs(bars[1].T) {
		t.Fatalf("entry timestamp mismatch: got %d", x.EntryTouchedAt)
	}
}

func TestV163ProfessionalInvalidationBeforeTarget(t *testing.T) {
	base := time.Date(2026, 2, 2, 20, 0, 0, 0, time.UTC)
	x := v163BaseSnapshot(base.UnixMilli())
	bars := []Bar{
		v163DailyBar(base.Add(24*time.Hour), 100, 103, 98, 101),
		v163DailyBar(base.Add(48*time.Hour), 99, 102, 93, 95),
		v163DailyBar(base.Add(72*time.Hour), 96, 112, 95, 111),
	}
	evaluateOneSignalSnapshotWithActions(&x, map[string][]Bar{"daily": bars}, nil, base.Add(96*time.Hour))
	if x.OutcomeState != "INVALIDATED" {
		t.Fatalf("expected invalidation first, got %s", x.OutcomeState)
	}
	if x.TargetTouchedAt != 0 || x.InvalidationAt == 0 {
		t.Fatalf("ordering evidence wrong: target=%d invalidation=%d", x.TargetTouchedAt, x.InvalidationAt)
	}
}

func TestV163ProfessionalEntryThenTargetAndExcursions(t *testing.T) {
	base := time.Date(2026, 3, 2, 20, 0, 0, 0, time.UTC)
	x := v163BaseSnapshot(base.UnixMilli())
	bars := []Bar{
		v163DailyBar(base.Add(24*time.Hour), 100, 103, 98.5, 101),
		v163DailyBar(base.Add(48*time.Hour), 101, 111, 97, 110),
	}
	evaluateOneSignalSnapshotWithActions(&x, map[string][]Bar{"daily": bars}, nil, base.Add(72*time.Hour))
	if x.OutcomeState != "TARGET_REACHED" || !x.EntryTouched {
		t.Fatalf("expected entry then target, got %+v", x)
	}
	if x.MFE <= 0 || x.MAE >= 0 {
		t.Fatalf("expected positive MFE and negative MAE, got mfe=%f mae=%f", x.MFE, x.MAE)
	}
	if x.ElapsedMinutes <= 0 {
		t.Fatalf("elapsed time should be recorded after resolution")
	}
}

func TestV163ProfessionalFirstEntryBarAmbiguityIsNotGuessed(t *testing.T) {
	base := time.Date(2026, 4, 2, 20, 0, 0, 0, time.UTC)
	x := v163BaseSnapshot(base.UnixMilli())
	bar := v163DailyBar(base.Add(24*time.Hour), 100, 111, 98, 109)
	evaluateOneSignalSnapshotWithActions(&x, map[string][]Bar{"daily": {bar}}, nil, base.Add(48*time.Hour))
	if x.OutcomeState != "AMBIGUOUS" {
		t.Fatalf("same-bar entry/target ordering must be ambiguous, got %s", x.OutcomeState)
	}
}

func TestV163ProfessionalNoFutureBarsDegradesTruthfully(t *testing.T) {
	base := time.Now().Add(-30 * 24 * time.Hour)
	x := v163BaseSnapshot(base.UnixMilli())
	evaluateOneSignalSnapshotWithActions(&x, map[string][]Bar{"daily": nil}, nil, time.Now())
	if x.OutcomeState != "UNAVAILABLE" {
		t.Fatalf("matured snapshot with no future evidence must be UNAVAILABLE, got %s", x.OutcomeState)
	}
}

func TestV163ProfessionalNoEntryUsesSupportedTenSessionWindow(t *testing.T) {
	base := time.Date(2026, 5, 1, 20, 0, 0, 0, time.UTC)
	x := v163BaseSnapshot(base.UnixMilli())
	bars := make([]Bar, 0, 12)
	for i := 1; i <= 10; i++ {
		bars = append(bars, v163DailyBar(base.Add(time.Duration(i)*24*time.Hour), 104, 108, 102, 105))
	}
	// A late entry after the supported 10-session outcome window must not retroactively convert NO_ENTRY.
	bars = append(bars, v163DailyBar(base.Add(11*24*time.Hour), 100, 102, 98, 99))
	evaluateOneSignalSnapshotWithActions(&x, map[string][]Bar{"daily": bars}, nil, base.Add(12*24*time.Hour))
	if x.OutcomeState != "NO_ENTRY" {
		t.Fatalf("late entry after 10 sessions must not rewrite supported outcome window, got %s", x.OutcomeState)
	}
	if x.EntryTouched {
		t.Fatalf("entry after the supported outcome window must not be marked as eligible entry")
	}
}

func TestV163ProfessionalWeekendGapDoesNotDelayCompletedDailyEvidence(t *testing.T) {
	friday := time.Date(2026, 5, 8, 0, 0, 0, 0, time.UTC)
	monday := friday.Add(72 * time.Hour)
	rows := []Bar{
		v163DailyBar(friday, 100, 103, 99, 102),
		v163DailyBar(monday, 103, 105, 102, 104),
	}
	cutoff := friday.Add(36 * time.Hour) // Saturday noon UTC; Friday daily bar is completed.
	got := completedBarsBefore(rows, cutoff.UnixMilli(), "daily")
	if len(got) != 1 || normalizedBarTimestampMs(got[0].T) != friday.UnixMilli() {
		t.Fatalf("weekend gap delayed completed Friday evidence: %+v", got)
	}
}

func TestV163ProfessionalSeasonalityUsesCanonicalAdjustedDailyOnly(t *testing.T) {
	start := time.Date(2015, 1, 2, 20, 0, 0, 0, time.UTC)
	daily := make([]Bar, 0, 420)
	price := 100.0
	for i := 0; i < 420; i++ {
		price *= 1.0005
		daily = append(daily, v163DailyBar(start.Add(time.Duration(i)*24*time.Hour), price*.99, price*1.01, price*.98, price))
	}
	rawWithSplitGap := append([]Bar(nil), daily...)
	rawWithSplitGap[len(rawWithSplitGap)/2].C = rawWithSplitGap[len(rawWithSplitGap)/2].C / 10
	bars := map[string]map[string][]Bar{
		"SPY": {"daily": daily, "daily-raw": rawWithSplitGap},
		"QQQ": {"daily": daily, "daily-raw": rawWithSplitGap},
	}
	got := buildSeasonalitySnapshot(bars, start.Add(500*24*time.Hour))
	if got.Symbols["SPY"].DailyBars != len(daily) {
		t.Fatalf("seasonality must use canonical daily history, got %d bars", got.Symbols["SPY"].DailyBars)
	}
	if !strings.Contains(got.Symbols["SPY"].Source, "Canonical daily history") {
		t.Fatalf("source truth missing: %q", got.Symbols["SPY"].Source)
	}
}

func TestV163ProfessionalSeasonalityInsufficientSampleDoesNotFabricateStatistics(t *testing.T) {
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

func TestV163ProfessionalCalibrationTinySampleAndNoProbabilityClaim(t *testing.T) {
	base := time.Now().Add(-90 * 24 * time.Hour).UnixMilli()
	s := SignalValidationState{Snapshots: []SignalSnapshot{{Symbol: "AAA", Horizon: "swing", Timestamp: base, Score: 82, Action: "APPROACHING ENTRY", OutcomeState: "TARGET_REACHED", EntryTouched: true, MFE: 8, MAE: -3, Outcomes: map[string]float64{"5D": 4.2}}}}
	got := buildCalibrationSnapshot(s, time.Now())
	if got.SetupScoreIsWinProbability {
		t.Fatalf("Setup Score must never become win probability")
	}
	if len(got.Groups) != 1 || got.Groups[0].State != "INSUFFICIENT" {
		t.Fatalf("tiny calibration sample must refuse confidence: %+v", got.Groups)
	}
	if strings.Contains(got.Message, "82%") || strings.Contains(got.Groups[0].Detail, "82%") {
		t.Fatalf("score-to-probability claim detected")
	}
}

func TestV163ProfessionalCalibrationDoesNotTreatNoEntryAsZeroExcursion(t *testing.T) {
	base := time.Now().Add(-90 * 24 * time.Hour).UnixMilli()
	s := SignalValidationState{Snapshots: []SignalSnapshot{
		{Symbol: "AAA", Horizon: "swing", Timestamp: base, Score: 72, Action: "APPROACHING ENTRY", OutcomeState: "NO_ENTRY", EntryTouched: false},
		{Symbol: "BBB", Horizon: "swing", Timestamp: base + 1, Score: 72, Action: "APPROACHING ENTRY", OutcomeState: "TARGET_REACHED", EntryTouched: true, MFE: 10, MAE: -4},
	}}
	got := buildCalibrationSnapshot(s, time.Now())
	if len(got.Groups) != 1 {
		t.Fatalf("expected one group, got %+v", got.Groups)
	}
	if got.Groups[0].AverageMFE == nil || math.Abs(*got.Groups[0].AverageMFE-10) > 1e-9 {
		t.Fatalf("NO_ENTRY must not inject synthetic zero MFE; got %+v", got.Groups[0].AverageMFE)
	}
	if got.Groups[0].AverageMAE == nil || math.Abs(*got.Groups[0].AverageMAE-(-4)) > 1e-9 {
		t.Fatalf("NO_ENTRY must not inject synthetic zero MAE; got %+v", got.Groups[0].AverageMAE)
	}
}

func TestV163ProfessionalCorrelationUsesAlignedReturnsAndFlagsSemiconductorConcentration(t *testing.T) {
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
			out = append(out, v163DailyBar(start.Add(time.Duration(i)*24*time.Hour), p, p*1.01, p*.99, p))
		}
		return out
	}
	st := AppState{Settings: Settings{DayEnabled: true, DayWatchlistID: "day"}, Watchlists: []Watchlist{{ID: "day", Symbols: []string{"NVDA", "AMD", "AVGO"}}}}
	bars := map[string]map[string][]Bar{
		"NVDA": {"daily": mk(1, 0)},
		"AMD":  {"daily": mk(1.05, 0)},
		"AVGO": {"daily": mk(.95, 11)},
	}
	got := buildCorrelationConcentrationSnapshot(st, ScannerState{}, bars, start.Add(120*24*time.Hour))
	foundPair := false
	for _, p := range got.Pairs {
		if (p.SymbolA == "NVDA" && p.SymbolB == "AMD") || (p.SymbolA == "AMD" && p.SymbolB == "NVDA") {
			foundPair = true
			if p.SampleCount < 60 || p.State != "HIGH" {
				t.Fatalf("expected aligned high-correlation semiconductor pair, got %+v", p)
			}
		}
	}
	if !foundPair {
		t.Fatalf("NVDA/AMD aligned correlation pair missing")
	}
	foundSector := false
	for _, g := range got.Concentrations {
		if g.Kind == "SECTOR" && g.Key == "Technology" && len(g.Symbols) >= 3 {
			foundSector = true
		}
	}
	if !foundSector {
		t.Fatalf("semiconductor candidate concentration not surfaced: %+v", got.Concentrations)
	}
}

func TestV163ProfessionalSignalSnapshotIdentityFreezesEvidenceAndSettings(t *testing.T) {
	e := &Engine{signalValidation: SignalValidationState{Snapshots: []SignalSnapshot{}}}
	base := v163BaseSnapshot(time.Now().UnixMilli())
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

func TestV163ProfessionalSplitAdjustmentPreservesFrozenOutcomeEconomics(t *testing.T) {
	base := time.Date(2025, 6, 2, 20, 0, 0, 0, time.UTC)
	x := v163BaseSnapshot(base.UnixMilli())
	// Current canonical adjustment=all history is on the post 10-for-1 split scale.
	bars := []Bar{
		v163DailyBar(base.Add(24*time.Hour), 9.9, 10.0, 9.8, 9.9),
		v163DailyBar(base.Add(48*time.Hour), 10.1, 11.1, 9.7, 11.0),
	}
	actions := []CorporateAction{{Symbol: "AAA", Type: "forward_split", ExDate: base.Add(36 * time.Hour).Format("2006-01-02"), Ratio: 10, AdjustmentFactor: 10, Status: "EFFECTIVE", Source: "Alpaca"}}
	evaluateOneSignalSnapshotWithActions(&x, map[string][]Bar{"daily": bars}, actions, base.Add(72*time.Hour))
	if x.OutcomeState != "TARGET_REACHED" {
		t.Fatalf("post-split adjusted history must preserve frozen target economics, got %s (%s)", x.OutcomeState, x.OutcomeDetail)
	}
	if math.Abs(x.OutcomeAdjustmentFactor-10) > 1e-9 || !strings.Contains(x.OutcomeAdjustmentDetail, "split evidence") {
		t.Fatalf("split normalization provenance missing: factor=%f detail=%q", x.OutcomeAdjustmentFactor, x.OutcomeAdjustmentDetail)
	}
	if got := x.Outcomes["1D"]; math.Abs(got-(-1.0)) > 0.05 { // 9.9 vs adjusted 10.0 anchor
		t.Fatalf("forward return must use split-adjusted frozen price anchor, got %.4f%%", got)
	}
}
