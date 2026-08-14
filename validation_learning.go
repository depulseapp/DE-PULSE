package main

import (
	"fmt"
	"math"
	"sort"
	"strings"
	"time"
)

const validationFormulaVersion = "deterministic-v14.3.7-compatible-v16.3"

func ptrFloat64(v float64) *float64 { return &v }

func normalizedBarTimestampMs(t int64) int64 {
	if t <= 0 {
		return 0
	}
	if t < 1_000_000_000_000 {
		return t * 1000
	}
	return t
}

func completedBarsBefore(bars []Bar, cutoffMs int64, timeframe string) []Bar {
	if cutoffMs <= 0 || len(bars) == 0 {
		return nil
	}
	fallback := int64(24 * time.Hour / time.Millisecond)
	switch strings.ToLower(strings.TrimSpace(timeframe)) {
	case "intraday":
		fallback = int64(15 * time.Minute / time.Millisecond)
	case "weekly":
		fallback = int64(7 * 24 * time.Hour / time.Millisecond)
	}
	out := make([]Bar, 0, len(bars))
	for _, b := range bars {
		start := normalizedBarTimestampMs(b.T)
		if start <= 0 {
			continue
		}
		// Use the canonical timeframe duration rather than the next observed bar.
		// Weekend/holiday gaps must not keep an already completed Friday daily bar
		// artificially "in progress" until the next trading session.
		end := start + fallback
		if end <= cutoffMs {
			out = append(out, b)
		}
	}
	return out
}

func futureBarsStrictlyAfter(bars []Bar, snapshotMs int64) []Bar {
	out := make([]Bar, 0, len(bars))
	for _, b := range bars {
		if normalizedBarTimestampMs(b.T) > snapshotMs {
			out = append(out, b)
		}
	}
	return out
}

func barTouchesRange(b Bar, low, high float64) bool {
	if low <= 0 || high <= 0 || high < low {
		return false
	}
	return b.L <= high && b.H >= low
}

func evaluateSignalSnapshotsProfessionalWithActions(s SignalValidationState, bars map[string]map[string][]Bar, actions []CorporateAction, mode string) SignalValidationState {
	if mode != "live" {
		return s
	}
	now := time.Now()
	for i := range s.Snapshots {
		x := &s.Snapshots[i]
		evaluateOneSignalSnapshotWithActions(x, bars[x.Symbol], actions, now)
	}
	s.UpdatedAt = time.Now().UnixMilli()
	s.Message = "Frozen deterministic snapshots are evaluated only from later canonical bars; target-before-entry never counts, ambiguous OHLC ordering stays ambiguous, and Setup Score is not probability."
	return s
}

func postSnapshotSplitAdjustment(symbol string, snapshotMs int64, actions []CorporateAction, now time.Time) (float64, string) {
	factor := 1.0
	used := []string{}
	for _, a := range actions {
		if normalizeSymbol(a.Symbol) != normalizeSymbol(symbol) && normalizeSymbol(a.OldSymbol) != normalizeSymbol(symbol) && normalizeSymbol(a.NewSymbol) != normalizeSymbol(symbol) {
			continue
		}
		kind := lifecycleActionKind(a.Type)
		if !strings.Contains(kind, "split") || a.Ratio <= 0 {
			continue
		}
		effective := lifecycleDateMillis(lifecycleEffectiveDate(a))
		if effective <= snapshotMs || effective <= 0 || effective > now.UnixMilli() {
			continue
		}
		factor *= a.Ratio
		used = append(used, fmt.Sprintf("%s %.4gx", defaultString(lifecycleEffectiveDate(a), "effective date unknown"), a.Ratio))
	}
	if factor <= 0 || math.IsNaN(factor) || math.IsInf(factor, 0) {
		return 1, ""
	}
	if math.Abs(factor-1) < 1e-12 {
		return 1, ""
	}
	return factor, "Frozen pre-action price levels were normalized to the current canonical adjusted-history scale using post-snapshot split evidence: " + strings.Join(used, " · ")
}

func evaluateOneSignalSnapshotWithActions(x *SignalSnapshot, symbolBars map[string][]Bar, actions []CorporateAction, now time.Time) {
	if x == nil || x.Timestamp <= 0 || x.Price <= 0 {
		return
	}
	adjustmentFactor, adjustmentDetail := postSnapshotSplitAdjustment(x.Symbol, x.Timestamp, actions, now)
	priceAnchor := x.Price
	entryLow, entryHigh := x.EntryLow, x.EntryHigh
	targetLow, targetHigh := x.TargetLow, x.TargetHigh
	invalidation := x.Invalidation
	if adjustmentFactor != 1 {
		priceAnchor /= adjustmentFactor
		entryLow /= adjustmentFactor
		entryHigh /= adjustmentFactor
		targetLow /= adjustmentFactor
		targetHigh /= adjustmentFactor
		invalidation /= adjustmentFactor
		x.OutcomeAdjustmentFactor = adjustmentFactor
		x.OutcomeAdjustmentDetail = adjustmentDetail
	}
	// Daily forward returns remain supported for 1D/3D/5D/10D analytics. Only completed
	// post-snapshot bars are eligible; an in-progress current session is not outcome truth.
	dailyCompleted := completedBarsBefore(symbolBars["daily"], now.UnixMilli(), "daily")
	daily := futureBarsStrictlyAfter(dailyCompleted, x.Timestamp)
	if len(daily) > 0 {
		if x.Outcomes == nil {
			x.Outcomes = map[string]float64{}
		}
		for _, n := range []int{1, 3, 5, 10} {
			if len(daily) >= n && daily[n-1].C > 0 {
				x.Outcomes[fmt.Sprintf("%dD", n)] = (daily[n-1].C/priceAnchor - 1) * 100
			}
		}
	}

	// Once a resolution has been established from canonical post-snapshot evidence,
	// preserve that frozen analytics result even if short intraday history later rolls
	// out of the live cache. Daily forward returns above may continue to fill in.
	switch x.OutcomeState {
	case "TARGET_REACHED", "INVALIDATED", "NO_ENTRY", "AMBIGUOUS":
		if x.OutcomeUpdatedAt > 0 {
			return
		}
	}

	levelsValid := entryLow > 0 && entryHigh >= entryLow && targetLow > 0 && invalidation > 0
	if !levelsValid {
		x.LegacyPartial = true
		if x.OutcomeState == "" {
			x.OutcomeState = "PARTIAL"
			x.OutcomeDetail = "Legacy snapshot has no frozen Entry Zone / target / invalidation levels; forward returns are retained but success/failure is not fabricated."
		}
		// Preserve the older MFE/MAE behavior for legacy snapshots so migration does not erase useful analytics.
		if len(daily) > 0 {
			maxH, minL := priceAnchor, priceAnchor
			limit := minInt(len(daily), 10)
			for _, b := range daily[:limit] {
				if b.H > maxH {
					maxH = b.H
				}
				if b.L < minL {
					minL = b.L
				}
			}
			x.MFE = (maxH/priceAnchor - 1) * 100
			x.MAE = (minL/priceAnchor - 1) * 100
			x.OutcomeUpdatedAt = normalizedBarTimestampMs(daily[limit-1].T)
		}
		return
	}

	resolutionKey := "daily"
	if x.Horizon == "day" && len(symbolBars["intraday"]) > 0 {
		resolutionKey = "intraday"
	}
	resolutionCompleted := completedBarsBefore(symbolBars[resolutionKey], now.UnixMilli(), resolutionKey)
	future := futureBarsStrictlyAfter(resolutionCompleted, x.Timestamp)
	// v16.3 outcome ordering is bounded to the supported 10-session analytics window.
	// A later entry must not retroactively convert a completed NO_ENTRY result.
	if resolutionKey == "daily" && len(future) > 10 {
		future = future[:10]
	}
	if len(future) == 0 {
		maxPendingAge := 16 * 24 * time.Hour
		if x.Horizon == "day" {
			maxPendingAge = 36 * time.Hour
		}
		snapshotAt := time.UnixMilli(x.Timestamp)
		if !snapshotAt.IsZero() && now.Sub(snapshotAt) > maxPendingAge {
			x.OutcomeState = "UNAVAILABLE"
			x.OutcomeDetail = "No post-snapshot canonical bars are available for an outcome window that should already have matured; missing evidence is not treated as zero or success."
		} else {
			x.OutcomeState = "PENDING"
			x.OutcomeDetail = "No completed post-snapshot canonical bars are available yet."
		}
		return
	}

	entryIndex := -1
	for i, b := range future {
		if barTouchesRange(b, entryLow, entryHigh) {
			entryIndex = i
			x.EntryTouched = true
			x.EntryTouchedAt = normalizedBarTimestampMs(b.T)
			// We cannot infer whether a target or invalidation happened before the first entry touch inside one OHLC bar.
			targetSame := b.H >= targetLow
			invalidSame := b.L <= invalidation
			if targetSame || invalidSame {
				x.OutcomeState = "AMBIGUOUS"
				x.OutcomeDetail = "Entry and a resolution level are both inside the first entry OHLC bar; finer canonical evidence is required to establish ordering."
				if targetSame {
					x.TargetTouchedAt = normalizedBarTimestampMs(b.T)
				}
				if invalidSame {
					x.InvalidationAt = normalizedBarTimestampMs(b.T)
				}
				x.OutcomeUpdatedAt = normalizedBarTimestampMs(b.T)
				return
			}
			break
		}
	}
	if entryIndex < 0 {
		if len(daily) >= 10 {
			x.OutcomeState = "NO_ENTRY"
			x.OutcomeDetail = "Entry Zone was not touched during the supported 10-session validation window."
		} else {
			x.OutcomeState = "PENDING"
			x.OutcomeDetail = "Entry Zone has not been touched; validation window is still incomplete."
		}
		x.OutcomeUpdatedAt = normalizedBarTimestampMs(future[len(future)-1].T)
		return
	}

	entryRef := (entryLow + entryHigh) / 2
	if entryRef <= 0 {
		entryRef = x.Price
	}
	maxH, minL := entryRef, entryRef
	for j := entryIndex; j < len(future); j++ {
		b := future[j]
		if b.H > maxH {
			maxH = b.H
		}
		if b.L < minL {
			minL = b.L
		}
		// The entry bar was already checked above. Resolution starts on later bars.
		if j == entryIndex {
			continue
		}
		hitTarget := b.H >= targetLow
		hitInvalid := b.L <= invalidation
		if hitTarget && hitInvalid {
			x.OutcomeState = "AMBIGUOUS"
			x.TargetTouchedAt = normalizedBarTimestampMs(b.T)
			x.InvalidationAt = normalizedBarTimestampMs(b.T)
			x.OutcomeDetail = "Target and invalidation were both inside one OHLC bar; ordering cannot be inferred without finer evidence."
			x.OutcomeUpdatedAt = normalizedBarTimestampMs(b.T)
			break
		}
		if hitInvalid {
			x.OutcomeState = "INVALIDATED"
			x.InvalidationAt = normalizedBarTimestampMs(b.T)
			x.OutcomeDetail = "Entry Zone was touched before invalidation; invalidation resolved first."
			x.OutcomeUpdatedAt = normalizedBarTimestampMs(b.T)
			break
		}
		if hitTarget {
			x.OutcomeState = "TARGET_REACHED"
			x.TargetTouchedAt = normalizedBarTimestampMs(b.T)
			x.OutcomeDetail = "Entry Zone was touched before target; target resolved first."
			x.OutcomeUpdatedAt = normalizedBarTimestampMs(b.T)
			break
		}
	}
	if x.OutcomeState == "" || x.OutcomeState == "PENDING" {
		x.OutcomeState = "PENDING"
		x.OutcomeDetail = "Entry Zone was touched; target/invalidation resolution is still pending."
		x.OutcomeUpdatedAt = normalizedBarTimestampMs(future[len(future)-1].T)
	}
	if entryRef > 0 {
		x.MFE = (maxH/entryRef - 1) * 100
		x.MAE = (minL/entryRef - 1) * 100
	}
	resolvedAt := x.TargetTouchedAt
	if x.InvalidationAt > 0 && (resolvedAt == 0 || x.InvalidationAt < resolvedAt) {
		resolvedAt = x.InvalidationAt
	}
	if resolvedAt > x.EntryTouchedAt && x.EntryTouchedAt > 0 {
		x.ElapsedMinutes = (resolvedAt - x.EntryTouchedAt) / int64(time.Minute/time.Millisecond)
	}
	if adjustmentDetail != "" && !strings.Contains(x.OutcomeDetail, adjustmentDetail) {
		x.OutcomeDetail = strings.TrimSpace(x.OutcomeDetail + " " + adjustmentDetail)
	}
}

func medianFloat64(v []float64) float64 {
	if len(v) == 0 {
		return 0
	}
	x := append([]float64(nil), v...)
	sort.Float64s(x)
	m := len(x) / 2
	if len(x)%2 == 1 {
		return x[m]
	}
	return (x[m-1] + x[m]) / 2
}

func averageFloat64(v []float64) float64 {
	if len(v) == 0 {
		return 0
	}
	t := 0.0
	for _, x := range v {
		t += x
	}
	return t / float64(len(v))
}

func positiveFrequency(v []float64) float64 {
	if len(v) == 0 {
		return 0
	}
	pos := 0
	for _, x := range v {
		if x > 0 {
			pos++
		}
	}
	return float64(pos) / float64(len(v)) * 100
}

func dateFromUnixBar(t int64) string {
	ms := normalizedBarTimestampMs(t)
	if ms <= 0 {
		return ""
	}
	return time.UnixMilli(ms).UTC().Format("2006-01-02")
}

func buildSeasonalitySnapshot(bars map[string]map[string][]Bar, now time.Time) SeasonalitySnapshot {
	out := SeasonalitySnapshot{Symbols: map[string]SeasonalitySymbolState{}, UpdatedAt: now.UnixMilli(), Message: "Trailing 10-year SPY/QQQ descriptive context only; current-year observations are shown separately and never mutate predictive scores."}
	for _, sym := range []string{"SPY", "QQQ"} {
		daily := completedBarsBefore(bars[sym]["daily"], now.UnixMilli(), "daily")
		st := SeasonalitySymbolState{Symbol: sym, State: "UNAVAILABLE", Source: "Canonical daily history · Provider Router · trailing 10 years", DailyBars: len(daily), Monthly: []SeasonalityMetric{}, DayOfWeek: []SeasonalityMetric{}}
		if len(daily) < 2 {
			st.Limitations = []string{"Missing completed canonical daily history."}
			out.Symbols[sym] = st
			continue
		}
		cutoff := now.AddDate(-10, 0, 0)
		filtered := make([]Bar, 0, len(daily))
		for _, b := range daily {
			t := time.UnixMilli(normalizedBarTimestampMs(b.T)).UTC()
			// Keep the prior month/year boundary needed to calculate the first
			// in-window monthly return without widening the published sample.
			if !t.Before(cutoff.AddDate(0, -2, 0)) {
				filtered = append(filtered, b)
			}
		}
		if len(filtered) < 2 {
			filtered = daily
		}
		st.SampleFrom = dateFromUnixBar(filtered[0].T)
		st.SampleTo = dateFromUnixBar(filtered[len(filtered)-1].T)

		type monthClose struct {
			year  int
			month time.Month
			close float64
			date  string
		}
		months := []monthClose{}
		for _, b := range filtered {
			t := time.UnixMilli(normalizedBarTimestampMs(b.T)).UTC()
			if b.C <= 0 {
				continue
			}
			if len(months) == 0 || months[len(months)-1].year != t.Year() || months[len(months)-1].month != t.Month() {
				months = append(months, monthClose{year: t.Year(), month: t.Month(), close: b.C, date: t.Format("2006-01-02")})
			} else {
				months[len(months)-1].close = b.C
				months[len(months)-1].date = t.Format("2006-01-02")
			}
		}
		monthReturns := map[time.Month][]float64{}
		monthDates := map[time.Month][]string{}
		currentYear := map[time.Month]float64{}
		currentObservation := map[time.Month]string{}
		for i := 1; i < len(months); i++ {
			if months[i-1].close <= 0 || months[i].close <= 0 {
				continue
			}
			r := (months[i].close/months[i-1].close - 1) * 100
			if months[i].year == now.Year() {
				currentYear[months[i].month] = r
				if months[i].month == now.Month() {
					currentObservation[months[i].month] = "MONTH-TO-DATE"
				} else {
					currentObservation[months[i].month] = "COMPLETED MONTH"
				}
				continue
			}
			if months[i].year < now.Year()-10 || months[i].year >= now.Year() {
				continue
			}
			monthReturns[months[i].month] = append(monthReturns[months[i].month], r)
			monthDates[months[i].month] = append(monthDates[months[i].month], months[i].date)
		}
		for m := time.January; m <= time.December; m++ {
			vals := monthReturns[m]
			metric := SeasonalityMetric{Key: fmt.Sprintf("month-%02d", int(m)), Label: m.String(), State: "INSUFFICIENT", SampleCount: len(vals), HistoricalYears: len(vals), Detail: "Requires at least 8 completed observations inside the trailing 10-year window."}
			if len(monthDates[m]) > 0 {
				metric.DateFrom = monthDates[m][0]
				metric.DateTo = monthDates[m][len(monthDates[m])-1]
			}
			if cur, ok := currentYear[m]; ok {
				metric.CurrentYearReturnPct = ptrFloat64(cur)
				metric.CurrentYearObservation = currentObservation[m]
			}
			if len(vals) >= 8 {
				metric.State = "AVAILABLE"
				metric.AverageReturnPct = ptrFloat64(averageFloat64(vals))
				metric.MedianReturnPct = ptrFloat64(medianFloat64(vals))
				metric.PositiveFrequencyPct = ptrFloat64(positiveFrequency(vals))
				best, worst := vals[0], vals[0]
				for _, v := range vals[1:] {
					if v > best {
						best = v
					}
					if v < worst {
						worst = v
					}
				}
				metric.BestReturnPct = ptrFloat64(best)
				metric.WorstReturnPct = ptrFloat64(worst)
				metric.Detail = "Trailing 10-year month-of-year distribution; current year is a separate comparison, not part of the historical average."
			}
			st.Monthly = append(st.Monthly, metric)
		}

		weekdayReturns := map[time.Weekday][]float64{}
		weekdayDates := map[time.Weekday][]string{}
		for i := 1; i < len(filtered); i++ {
			if filtered[i-1].C <= 0 || filtered[i].C <= 0 {
				continue
			}
			t := time.UnixMilli(normalizedBarTimestampMs(filtered[i].T)).UTC()
			if t.Before(cutoff) {
				continue
			}
			r := (filtered[i].C/filtered[i-1].C - 1) * 100
			weekdayReturns[t.Weekday()] = append(weekdayReturns[t.Weekday()], r)
			weekdayDates[t.Weekday()] = append(weekdayDates[t.Weekday()], t.Format("2006-01-02"))
		}
		for _, wd := range []time.Weekday{time.Monday, time.Tuesday, time.Wednesday, time.Thursday, time.Friday} {
			vals := weekdayReturns[wd]
			metric := SeasonalityMetric{Key: "weekday-" + strings.ToLower(wd.String()), Label: wd.String(), State: "INSUFFICIENT", SampleCount: len(vals), Detail: "Requires at least 252 completed sessions overall and 40 observations for this weekday."}
			if len(weekdayDates[wd]) > 0 {
				metric.DateFrom = weekdayDates[wd][0]
				metric.DateTo = weekdayDates[wd][len(weekdayDates[wd])-1]
			}
			if len(filtered) >= 252 && len(vals) >= 40 {
				metric.State = "AVAILABLE"
				metric.AverageReturnPct = ptrFloat64(averageFloat64(vals))
				metric.MedianReturnPct = ptrFloat64(medianFloat64(vals))
				metric.PositiveFrequencyPct = ptrFloat64(positiveFrequency(vals))
				metric.Detail = "Trailing 10-year descriptive day-of-week history; not a directional trading instruction."
			}
			st.DayOfWeek = append(st.DayOfWeek, metric)
		}
		available := false
		for _, m := range st.Monthly {
			available = available || m.State == "AVAILABLE"
		}
		for _, m := range st.DayOfWeek {
			available = available || m.State == "AVAILABLE"
		}
		if available {
			st.State = "AVAILABLE"
		} else {
			st.State = "INSUFFICIENT"
			st.Limitations = append(st.Limitations, "Current canonical history depth does not yet meet professional sample floors.")
		}
		st.Limitations = append(st.Limitations, "Event-window seasonality remains UNAVAILABLE unless event-aligned historical evidence has defensible timestamps and sample depth; no proxy is fabricated.")
		out.Symbols[sym] = st
	}
	return out
}

func scoreBand(score float64) string {
	switch {
	case score < 40:
		return "0-39"
	case score < 55:
		return "40-54"
	case score < 70:
		return "55-69"
	case score < 85:
		return "70-84"
	default:
		return "85-100"
	}
}

func calibrationSampleState(n int) string {
	switch {
	case n < 30:
		return "INSUFFICIENT"
	case n < 100:
		return "LOW SAMPLE CONFIDENCE"
	case n < 300:
		return "MEDIUM SAMPLE CONFIDENCE"
	default:
		return "HIGHER SAMPLE DEPTH"
	}
}

func buildCalibrationSnapshot(s SignalValidationState, now time.Time) CalibrationSnapshot {
	type acc struct {
		group CalibrationGroup
		ret5  []float64
		mfe   []float64
		mae   []float64
	}
	m := map[string]*acc{}
	eligible := 0
	for _, x := range s.Snapshots {
		if x.LegacyPartial || x.OutcomeState == "" || x.OutcomeState == "PENDING" || x.OutcomeState == "PARTIAL" || x.OutcomeState == "UNAVAILABLE" {
			continue
		}
		eligible++
		band := scoreBand(x.Score)
		key := x.Horizon + "|" + band + "|" + x.Action
		a := m[key]
		if a == nil {
			a = &acc{group: CalibrationGroup{Horizon: x.Horizon, ScoreBand: band, Action: x.Action, DateFrom: x.Timestamp, DateTo: x.Timestamp}}
			m[key] = a
		}
		a.group.SampleCount++
		if x.Timestamp < a.group.DateFrom || a.group.DateFrom == 0 {
			a.group.DateFrom = x.Timestamp
		}
		if x.Timestamp > a.group.DateTo {
			a.group.DateTo = x.Timestamp
		}
		switch x.OutcomeState {
		case "TARGET_REACHED":
			a.group.TargetReached++
		case "INVALIDATED":
			a.group.Invalidated++
		case "NO_ENTRY":
			a.group.NoEntry++
		case "AMBIGUOUS":
			a.group.Ambiguous++
		}
		if v, ok := x.Outcomes["5D"]; ok {
			a.ret5 = append(a.ret5, v)
		}
		// Excursion statistics are meaningful only after an eligible Entry Zone touch.
		// NO_ENTRY rows must not silently inject synthetic 0% MFE/MAE into calibration.
		if x.EntryTouched {
			a.mfe = append(a.mfe, x.MFE)
			a.mae = append(a.mae, x.MAE)
		}
	}
	groups := make([]CalibrationGroup, 0, len(m))
	for _, a := range m {
		g := a.group
		g.State = calibrationSampleState(g.SampleCount)
		resolved := g.TargetReached + g.Invalidated
		if resolved > 0 {
			g.TargetReachedPct = ptrFloat64(float64(g.TargetReached) / float64(resolved) * 100)
		}
		if len(a.ret5) > 0 {
			g.Average5DReturn = ptrFloat64(averageFloat64(a.ret5))
		}
		if len(a.mfe) > 0 {
			g.AverageMFE = ptrFloat64(averageFloat64(a.mfe))
		}
		if len(a.mae) > 0 {
			g.AverageMAE = ptrFloat64(averageFloat64(a.mae))
		}
		g.Detail = "Descriptive historical outcomes only. Setup Score is not win probability and this group cannot auto-change production formulas."
		groups = append(groups, g)
	}
	sort.Slice(groups, func(i, j int) bool {
		if groups[i].Horizon != groups[j].Horizon {
			return groups[i].Horizon < groups[j].Horizon
		}
		if groups[i].ScoreBand != groups[j].ScoreBand {
			return groups[i].ScoreBand < groups[j].ScoreBand
		}
		return groups[i].Action < groups[j].Action
	})
	return CalibrationSnapshot{Groups: groups, EligibleSnapshots: eligible, SetupScoreIsWinProbability: false, UpdatedAt: now.UnixMilli(), Message: "Calibration is descriptive validation, never probability conversion or automatic formula learning."}
}

func dailyReturnMap(bars []Bar, nowMs int64) (map[int64]float64, string, string) {
	completed := completedBarsBefore(bars, nowMs, "daily")
	out := map[int64]float64{}
	for i := 1; i < len(completed); i++ {
		if completed[i-1].C <= 0 || completed[i].C <= 0 {
			continue
		}
		out[completed[i].T] = completed[i].C/completed[i-1].C - 1
	}
	from, to := "", ""
	if len(completed) > 0 {
		from = dateFromUnixBar(completed[0].T)
		to = dateFromUnixBar(completed[len(completed)-1].T)
	}
	return out, from, to
}

func pearsonAligned(a, b map[int64]float64) (float64, int) {
	keys := make([]int64, 0)
	for k := range a {
		if _, ok := b[k]; ok {
			keys = append(keys, k)
		}
	}
	if len(keys) < 2 {
		return 0, len(keys)
	}
	sort.Slice(keys, func(i, j int) bool { return keys[i] < keys[j] })
	ma, mb := 0.0, 0.0
	for _, k := range keys {
		ma += a[k]
		mb += b[k]
	}
	ma /= float64(len(keys))
	mb /= float64(len(keys))
	num, da, db := 0.0, 0.0, 0.0
	for _, k := range keys {
		x, y := a[k]-ma, b[k]-mb
		num += x * y
		da += x * x
		db += y * y
	}
	if da <= 0 || db <= 0 {
		return 0, len(keys)
	}
	return num / math.Sqrt(da*db), len(keys)
}

func validationCandidates(st AppState, scanner ScannerState) []string {
	out := activeDeskSymbolsFromState(st)
	if sym, ok := parseUserTicker(st.UI.SelectedTicker); ok {
		out = append(out, sym)
	}
	// Scanner/Discovery candidates are attention context, not positions. Keep only a bounded
	// high-ranked slice so correlation work cannot grow quadratically with an unbounded universe.
	for i, row := range scanner.Results {
		if i >= 20 {
			break
		}
		if sym, ok := parseUserTicker(row.Symbol); ok {
			out = append(out, sym)
		}
	}
	return uniqueSymbols(out)
}

func buildCorrelationConcentrationSnapshot(st AppState, scanner ScannerState, bars map[string]map[string][]Bar, now time.Time) CorrelationConcentrationSnapshot {
	candidates := validationCandidates(st, scanner)
	if len(candidates) > 30 {
		candidates = candidates[:30]
	}
	returns := map[string]map[int64]float64{}
	ranges := map[string][2]string{}
	for _, sym := range candidates {
		r, from, to := dailyReturnMap(bars[sym]["daily"], now.UnixMilli())
		returns[sym] = r
		ranges[sym] = [2]string{from, to}
	}
	pairs := []CorrelationPair{}
	for i := 0; i < len(candidates); i++ {
		for j := i + 1; j < len(candidates); j++ {
			a, b := candidates[i], candidates[j]
			corr, n := pearsonAligned(returns[a], returns[b])
			p := CorrelationPair{SymbolA: a, SymbolB: b, SampleCount: n, State: "INSUFFICIENT", Detail: "Requires at least 60 aligned completed daily returns."}
			p.WindowFrom = ranges[a][0]
			if ranges[b][0] > p.WindowFrom {
				p.WindowFrom = ranges[b][0]
			}
			p.WindowTo = ranges[a][1]
			if p.WindowTo == "" || (ranges[b][1] != "" && ranges[b][1] < p.WindowTo) {
				p.WindowTo = ranges[b][1]
			}
			if n >= 60 {
				p.Correlation = corr
				switch {
				case corr >= 0.85:
					p.State = "HIGH"
				case corr >= 0.70:
					p.State = "ELEVATED"
				default:
					p.State = "NORMAL"
				}
				if n < 126 {
					p.Detail = "Low-sample correlation context (60–125 aligned returns); attention only."
				} else {
					p.Detail = "Aligned-return correlation context; attention only, never automatic trade rejection."
				}
			}
			pairs = append(pairs, p)
		}
	}
	sort.Slice(pairs, func(i, j int) bool {
		rank := func(s string) int {
			switch s {
			case "HIGH":
				return 3
			case "ELEVATED":
				return 2
			case "NORMAL":
				return 1
			default:
				return 0
			}
		}
		ri, rj := rank(pairs[i].State), rank(pairs[j].State)
		if ri != rj {
			return ri > rj
		}
		return math.Abs(pairs[i].Correlation) > math.Abs(pairs[j].Correlation)
	})
	if len(pairs) > 120 {
		pairs = pairs[:120]
	}

	sectorMembers := map[string][]string{}
	industryMembers := map[string][]string{}
	for _, sym := range candidates {
		c, ok := canonicalSymbolClassifications[sym]
		if !ok {
			continue
		}
		if c.Sector != "" {
			sectorMembers[c.Sector] = append(sectorMembers[c.Sector], sym)
		}
		if c.Industry != "" {
			industryMembers[c.Industry] = append(industryMembers[c.Industry], sym)
		}
	}
	groups := []ConcentrationGroup{}
	for k, members := range sectorMembers {
		members = uniqueSymbols(members)
		if len(members) >= 2 {
			groups = append(groups, ConcentrationGroup{Key: k, Kind: "SECTOR", Symbols: members, State: "CONCENTRATED", Detail: "Multiple current candidates share the same canonical sector classification; attention context only."})
		}
	}
	for k, members := range industryMembers {
		members = uniqueSymbols(members)
		if len(members) >= 2 {
			groups = append(groups, ConcentrationGroup{Key: k, Kind: "INDUSTRY", Symbols: members, State: "CONCENTRATED", Detail: "Multiple current candidates share the same canonical industry classification; attention context only."})
		}
	}
	// Add pair-derived factor concentration when high/elevated correlation exists.
	for _, p := range pairs {
		if p.State == "HIGH" || p.State == "ELEVATED" {
			groups = append(groups, ConcentrationGroup{Key: p.SymbolA + "/" + p.SymbolB, Kind: "RETURN CORRELATION", Symbols: []string{p.SymbolA, p.SymbolB}, State: p.State, Detail: fmt.Sprintf("%.2f correlation across %d aligned daily returns; attention context only.", p.Correlation, p.SampleCount)})
		}
	}
	sort.Slice(groups, func(i, j int) bool {
		if groups[i].Kind != groups[j].Kind {
			return groups[i].Kind < groups[j].Kind
		}
		return groups[i].Key < groups[j].Key
	})
	if len(groups) > 40 {
		groups = groups[:40]
	}
	return CorrelationConcentrationSnapshot{Pairs: pairs, Concentrations: groups, CandidateCount: len(candidates), UpdatedAt: now.UnixMilli(), Message: "Correlation/concentration uses aligned canonical returns and candidate context only; no Portfolio, positions/P&L or automatic rejection."}
}

func buildValidationLearningSnapshot(st AppState, validation SignalValidationState, scanner ScannerState, bars map[string]map[string][]Bar, events []MacroEvent, earnings []EarningsItem, now time.Time) ValidationLearningSnapshot {
	return ValidationLearningSnapshot{
		Seasonality:   buildSeasonalitySnapshot(bars, now),
		Calibration:   buildCalibrationSnapshot(validation, now),
		Concentration: buildCorrelationConcentrationSnapshot(st, scanner, bars, now),
		ReplayState:   "AVAILABLE WITH TRUTH LIMITS",
		ReplayDetail:  "Replay uses the same deterministic renderer core against an isolated cutoff-filtered runtime. Reusable CPI/FOMC/earnings-gap/high-VIX/dislocation scenarios are cataloged only when canonical retained evidence exists. Exact snapshot replay requires frozen evidence/settings; otherwise it is explicitly SCENARIO/PARTIAL/UNAVAILABLE.",
		ReplayCatalog: replayScenarioCatalog(st, bars, events, earnings, now),
		UpdatedAt:     now.UnixMilli(),
	}
}

func replayScenarioTimestampForDate(date string, hour string) int64 {
	date = strings.TrimSpace(date)
	if len(date) < 10 {
		return 0
	}
	clock := "12:00:00"
	h := strings.ToLower(strings.TrimSpace(hour))
	if strings.Contains(h, "bmo") {
		clock = "08:00:00"
	}
	if strings.Contains(h, "amc") {
		clock = "16:30:00"
	}
	t, err := time.Parse("2006-01-02 15:04:05", date[:10]+" "+clock)
	if err != nil {
		return 0
	}
	return t.UTC().UnixMilli()
}

func replayScenarioBarGap(prev, cur Bar) float64 {
	if prev.C <= 0 || cur.O <= 0 {
		return 0
	}
	return (cur.O/prev.C - 1) * 100
}

func replayScenarioCatalog(st AppState, bars map[string]map[string][]Bar, events []MacroEvent, earnings []EarningsItem, now time.Time) ReplayScenarioCatalog {
	out := ReplayScenarioCatalog{Scenarios: []ReplayScenarioDescriptor{}, Kinds: []string{"CPI_SHOCK", "FOMC_SHOCK", "EARNINGS_GAP", "HIGH_VIX", "MARKET_DISLOCATION"}, UpdatedAt: now.UnixMilli(), Message: "Reusable historical scenario descriptors are derived only from canonical retained evidence. Running a scenario reuses the existing cutoff-filtered deterministic replay engine; it never creates mock-live/paper-trading state."}
	seen := map[string]bool{}
	add := func(x ReplayScenarioDescriptor) {
		if x.ID == "" || seen[x.ID] {
			return
		}
		seen[x.ID] = true
		if x.Cutoff > 0 && x.Cutoff <= now.UnixMilli() {
			x.State = "AVAILABLE"
			out.Available++
		} else if x.State == "" {
			x.State = "UNAVAILABLE"
		}
		out.Scenarios = append(out.Scenarios, x)
	}
	for _, e := range events {
		if e.StartsAt <= 0 || e.StartsAt > now.UnixMilli() {
			continue
		}
		name := strings.ToUpper(e.Name)
		switch {
		case strings.Contains(name, "CPI") || strings.Contains(name, "CONSUMER PRICE"):
			add(ReplayScenarioDescriptor{ID: "macro-cpi-" + e.ID, Kind: "CPI_SHOCK", Label: "CPI · " + e.Date, Symbol: "SPY", Horizon: "day", Cutoff: e.StartsAt, Source: e.Source, Evidence: e.Name, Detail: "Historical CPI event cutoff; only evidence knowable at or before the event timestamp is eligible."})
		case strings.Contains(name, "FOMC") || strings.Contains(name, "FEDERAL RESERVE") || strings.Contains(name, "FED RATE"):
			add(ReplayScenarioDescriptor{ID: "macro-fomc-" + e.ID, Kind: "FOMC_SHOCK", Label: "FOMC/Fed · " + e.Date, Symbol: "QQQ", Horizon: "day", Cutoff: e.StartsAt, Source: e.Source, Evidence: e.Name, Detail: "Historical Fed-event cutoff; replay strips present-day caches and future event corrections."})
		}
	}
	// Earnings-gap scenarios require an actual earnings row and matching canonical daily bars.
	for _, er := range earnings {
		sym := normalizeSymbol(er.Symbol)
		if sym == "" {
			continue
		}
		cut := replayScenarioTimestampForDate(er.Date, er.Hour)
		if cut <= 0 || cut > now.UnixMilli() {
			continue
		}
		daily := bars[sym]["daily"]
		bestGap := 0.0
		for i := 1; i < len(daily); i++ {
			bt := normalizedBarTimestampMs(daily[i].T)
			if absInt64(bt-cut) > int64(36*time.Hour/time.Millisecond) {
				continue
			}
			gap := replayScenarioBarGap(daily[i-1], daily[i])
			if math.Abs(gap) > math.Abs(bestGap) {
				bestGap = gap
			}
		}
		if math.Abs(bestGap) >= 3 {
			add(ReplayScenarioDescriptor{ID: fmt.Sprintf("earnings-gap-%s-%s", sym, strings.ReplaceAll(er.Date, "-", "")), Kind: "EARNINGS_GAP", Label: fmt.Sprintf("%s earnings gap %.1f%%", sym, bestGap), Symbol: sym, Horizon: "swing", Cutoff: cut, Source: "Canonical earnings + daily bars", Evidence: fmt.Sprintf("earnings date %s · observed opening gap %.1f%%", er.Date, bestGap), Detail: "Earnings-linked gap scenario. Replay uses cutoff-safe price history; later earnings values/results are not backfilled into the historical decision."})
		}
	}
	// High-VIX sessions come only from retained canonical VIX history.
	vix := bars["VIX"]["daily"]
	for i := len(vix) - 1; i >= 0; i-- {
		if vix[i].C >= 30 {
			cut := normalizedBarTimestampMs(vix[i].T) + int64(12*time.Hour/time.Millisecond)
			add(ReplayScenarioDescriptor{ID: fmt.Sprintf("high-vix-%d", normalizedBarTimestampMs(vix[i].T)), Kind: "HIGH_VIX", Label: fmt.Sprintf("High VIX %.1f", vix[i].C), Symbol: "SPY", Horizon: "swing", Cutoff: cut, Source: "Canonical VIX daily history", Evidence: fmt.Sprintf("VIX close %.1f", vix[i].C), Detail: "High-volatility regime scenario from retained VIX history; no future bars are eligible."})
			break
		}
	}
	// A broad dislocation scenario uses the largest retained SPY daily move, not a hard-coded historical claim.
	spy := bars["SPY"]["daily"]
	best := 0.0
	bestIdx := -1
	for i := 1; i < len(spy); i++ {
		if spy[i-1].C <= 0 {
			continue
		}
		r := (spy[i].C/spy[i-1].C - 1) * 100
		if math.Abs(r) > math.Abs(best) {
			best = r
			bestIdx = i
		}
	}
	if bestIdx > 0 && math.Abs(best) >= 2.5 {
		cut := normalizedBarTimestampMs(spy[bestIdx].T) + int64(12*time.Hour/time.Millisecond)
		add(ReplayScenarioDescriptor{ID: fmt.Sprintf("spy-dislocation-%d", normalizedBarTimestampMs(spy[bestIdx].T)), Kind: "MARKET_DISLOCATION", Label: fmt.Sprintf("SPY dislocation %.1f%%", best), Symbol: "SPY", Horizon: "swing", Cutoff: cut, Source: "Canonical SPY daily history", Evidence: fmt.Sprintf("largest retained daily move %.1f%%", best), Detail: "Reusable market-dislocation scenario selected from retained canonical SPY history."})
	}
	// Truthful unavailable archetypes remain visible when the current cache has no eligible retained evidence.
	for _, kind := range out.Kinds {
		found := false
		for _, x := range out.Scenarios {
			if x.Kind == kind {
				found = true
				break
			}
		}
		if !found {
			add(ReplayScenarioDescriptor{ID: "unavailable-" + strings.ToLower(kind), Kind: kind, Label: strings.ReplaceAll(kind, "_", " "), State: "UNAVAILABLE", Source: "Canonical retained evidence", Detail: "No eligible retained historical evidence is currently available for this scenario class; DE.PULSE does not fabricate a case."})
		}
	}
	sort.SliceStable(out.Scenarios, func(i, j int) bool {
		if out.Scenarios[i].State != out.Scenarios[j].State {
			return out.Scenarios[i].State == "AVAILABLE"
		}
		return out.Scenarios[i].Cutoff > out.Scenarios[j].Cutoff
	})
	if len(out.Scenarios) > 20 {
		out.Scenarios = out.Scenarios[:20]
	}
	return out
}
