package main

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// ---- Signal validation ----------------------------------------------------
func (e *Engine) recordSignalSnapshot(in SignalSnapshot) SignalSnapshot {
	in.Symbol = normalizeSymbol(in.Symbol)
	in.Horizon = strings.ToLower(strings.TrimSpace(in.Horizon))
	if in.Timestamp == 0 {
		in.Timestamp = time.Now().UnixMilli()
	}
	if strings.TrimSpace(in.FormulaVersion) == "" {
		in.FormulaVersion = validationFormulaVersion
	}
	// v16.3 identity prefers immutable evidence/setup identity. Legacy clients without
	// EvidenceSnapshotID retain the previous time-window dedupe behavior.
	if in.ID == "" {
		if strings.TrimSpace(in.EvidenceSnapshotID) != "" {
			in.ID = fmt.Sprintf("%s-%s-%s-%s", in.Symbol, in.Horizon, strings.TrimSpace(in.EvidenceSnapshotID), strings.TrimSpace(in.FormulaVersion))
		} else {
			in.ID = fmt.Sprintf("%s-%s-%d", in.Symbol, in.Horizon, in.Timestamp)
			in.LegacyPartial = true
		}
	}
	e.mu.Lock()
	for i := len(e.signalValidation.Snapshots) - 1; i >= 0; i-- {
		x := e.signalValidation.Snapshots[i]
		if x.ID == in.ID {
			e.mu.Unlock()
			return x
		}
		if strings.TrimSpace(in.EvidenceSnapshotID) == "" && x.Symbol == in.Symbol && x.Horizon == in.Horizon && x.Timestamp >= in.Timestamp-30*60*1000 {
			e.mu.Unlock()
			return x
		}
	}
	e.signalValidation.Snapshots = append(e.signalValidation.Snapshots, in)
	if len(e.signalValidation.Snapshots) > 1000 {
		e.signalValidation.Snapshots = e.signalValidation.Snapshots[len(e.signalValidation.Snapshots)-1000:]
	}
	e.signalValidation.UpdatedAt = time.Now().UnixMilli()
	e.signalValidation.Message = "Frozen validation snapshots preserve deterministic setup/evidence truth; outcomes use only later canonical bars and never auto-change formulas."
	e.mu.Unlock()
	e.enqueueSignalSnapshotPersistence(in)
	return in
}

func evaluateSignalSnapshotsWithActions(s SignalValidationState, bars map[string]map[string][]Bar, actions []CorporateAction, mode string) SignalValidationState {
	return evaluateSignalSnapshotsProfessionalWithActions(s, bars, actions, mode)
}

// compactAIEvidence deliberately excludes entire watchlists/scanners. It keeps only the
// target symbol and review-type-relevant evidence to stay below provider context limits.
func compactAIEvidence(req AIRequest, snap RuntimeSnapshot) map[string]any {
	sym := normalizeSymbol(req.Ticker)
	kind := strings.ToLower(strings.TrimSpace(req.Kind))
	out := map[string]any{"generatedAt": time.Now().UTC().Format(time.RFC3339), "dataMode": snap.Mode, "runtimeStatus": snap.Status, "selectedTicker": sym, "marketSession": snap.Feed.MarketSession, "feedState": snap.Feed.FeedState, "quote": snap.Quotes[sym], "deskContext": req.ClientContext, "global": snap.Global, "eventMode": snap.EventMode, "options": snap.Options[sym]}
	if f, ok := snap.Fundamentals[sym]; ok {
		out["fundamentals"] = f
	}
	news := []NewsItem{}
	for _, n := range snap.News {
		match := false
		for _, s := range n.Symbols {
			if normalizeSymbol(s) == sym {
				match = true
				break
			}
		}
		if match {
			news = append(news, n)
			if len(news) >= 4 {
				break
			}
		}
	}
	earnings := []EarningsItem{}
	for _, x := range snap.Earnings {
		if normalizeSymbol(x.Symbol) == sym {
			earnings = append(earnings, x)
			if len(earnings) >= 3 {
				break
			}
		}
	}
	filings := []FilingItem{}
	for _, x := range snap.Filings {
		if normalizeSymbol(x.Symbol) == sym {
			filings = append(filings, x)
			if len(filings) >= 4 {
				break
			}
		}
	}
	switch kind {
	case "risk":
		out["earnings"] = earnings
		out["filings"] = filings
		out["recentNews"] = news
		out["macroEvents"] = nearestEvents(snap.MacroEvents, 4)
	case "news":
		out["recentNews"] = news
		out["earnings"] = earnings
		out["filings"] = filings
	case "ticker":
		out["recentNews"] = news
		out["earnings"] = earnings
		out["filings"] = filings
		out["secIntelligence"] = snap.SECIntelligence[sym]
	default:
		out["recentNews"] = news
		out["earnings"] = earnings
	}
	return out
}
func nearestEvents(events []MacroEvent, n int) []MacroEvent {
	now := time.Now().Add(-24 * time.Hour).UnixMilli()
	out := []MacroEvent{}
	for _, e := range events {
		if e.StartsAt >= now {
			out = append(out, e)
			if len(out) >= n {
				break
			}
		}
	}
	return out
}
func isContextLimitError(err error) bool {
	if err == nil {
		return false
	}
	s := strings.ToLower(err.Error())
	return strings.Contains(s, "token") || strings.Contains(s, "context") || strings.Contains(s, "tpm") || strings.Contains(s, "too large") || strings.Contains(s, "request too")
}

func (a *Application) handleSignalValidationRecord(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}
	var in SignalSnapshot
	if decodeJSON(r, &in) != nil {
		writeError(w, http.StatusBadRequest, "Invalid signal snapshot")
		return
	}
	sym, ok := parseUserTicker(in.Symbol)
	if !ok {
		writeError(w, http.StatusBadRequest, "Invalid signal snapshot")
		return
	}
	in.Symbol = sym
	in.Horizon = strings.ToLower(strings.TrimSpace(in.Horizon))
	if in.Horizon != "day" && in.Horizon != "swing" && in.Horizon != "long" {
		writeError(w, http.StatusBadRequest, "Invalid signal snapshot")
		return
	}
	saved := a.engine.recordSignalSnapshot(in)
	_ = a.engine.saveCache()
	writeJSON(w, http.StatusOK, map[string]any{"ok": true, "snapshot": saved})
}

func testFREDProvider(ctx context.Context, key string) ProviderTestResult {
	r := ProviderTestResult{Provider: "fred", CheckedAt: time.Now().UTC().Format(time.RFC3339)}
	if strings.TrimSpace(key) == "" {
		r.Status = "missing"
		r.Message = "Enter a free FRED API key."
		return r
	}
	var p struct {
		Observations []struct {
			Value string `json:"value"`
		} `json:"observations"`
	}
	raw := "https://api.stlouisfed.org/fred/series/observations?series_id=DGS10&api_key=" + url.QueryEscape(key) + "&file_type=json&sort_order=desc&limit=3"
	if err := getJSON(ctx, &http.Client{Timeout: 12 * time.Second}, raw, nil, &p); err != nil {
		r.Status = "failed"
		r.Message = err.Error()
		return r
	}
	r.OK = true
	r.Status = "connected"
	r.Message = "FRED official rates access is working."
	r.Details = []string{"DGS10 representative request succeeded", "Used for rates/real-yield/credit context"}
	return r
}
func testBLSProvider(ctx context.Context, key string) ProviderTestResult {
	r := ProviderTestResult{Provider: "bls", CheckedAt: time.Now().UTC().Format(time.RFC3339)}
	body := strings.NewReader(`{"seriesid":["CUUR0000SA0"],"startyear":"2025","endyear":"2026"}`)
	req, _ := http.NewRequestWithContext(ctx, http.MethodPost, strings.TrimRight(blsAPIBaseURL, "/")+"/", body)
	req.Header.Set("Content-Type", "application/json")
	resp, err := (&http.Client{Timeout: 12 * time.Second}).Do(req)
	if err != nil {
		r.Status = "failed"
		r.Message = err.Error()
		return r
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		r.Status = "failed"
		r.Message = fmt.Sprintf("BLS HTTP %d", resp.StatusCode)
		return r
	}
	r.OK = true
	r.Status = "connected"
	r.Message = "BLS official CPI/labor API is reachable."
	if strings.TrimSpace(key) == "" {
		r.Details = []string{"Public API mode active", "Registration key is optional for higher limits"}
	} else {
		r.Details = []string{"Registered-key field is stored securely", "Public representative request succeeded"}
	}
	return r
}
func testEIAProvider(ctx context.Context, key string) ProviderTestResult {
	r := ProviderTestResult{Provider: "eia", CheckedAt: time.Now().UTC().Format(time.RFC3339)}
	if strings.TrimSpace(key) == "" {
		r.Status = "missing"
		r.Message = "Enter a free EIA API key."
		return r
	}
	raw := "https://api.eia.gov/v2/petroleum/pri/spt/data/?api_key=" + url.QueryEscape(key) + "&length=1"
	var p map[string]any
	if err := getJSON(ctx, &http.Client{Timeout: 12 * time.Second}, raw, nil, &p); err != nil {
		r.Status = "failed"
		r.Message = err.Error()
		return r
	}
	r.OK = true
	r.Status = "connected"
	r.Message = "EIA official energy-data access is working."
	return r
}
func testTwelveDataProvider(ctx context.Context, key string) ProviderTestResult {
	r := ProviderTestResult{Provider: "twelvedata", CheckedAt: time.Now().UTC().Format(time.RFC3339)}
	if strings.TrimSpace(key) == "" {
		r.Status = "missing"
		r.Message = "Enter a Twelve Data key only if you want the optional direct-global provider."
		return r
	}
	var p map[string]any
	raw := "https://api.twelvedata.com/quote?symbol=SPY&apikey=" + url.QueryEscape(key)
	if err := getJSON(ctx, &http.Client{Timeout: 12 * time.Second}, raw, nil, &p); err != nil {
		r.Status = "failed"
		r.Message = err.Error()
		return r
	}
	if strings.TrimSpace(fmt.Sprint(p["status"])) == "error" {
		r.Status = "failed"
		r.Message = fmt.Sprint(p["message"])
		return r
	}
	r.OK = true
	r.Status = "connected"
	r.Message = "Twelve Data representative quote request succeeded."
	r.Details = []string{"Direct/global interface configured", "Actual international entitlement still depends on the provider plan"}
	return r
}
func testOptionsProvider(ctx context.Context, key, secret, mode string) ProviderTestResult {
	r := ProviderTestResult{Provider: "options", CheckedAt: time.Now().UTC().Format(time.RFC3339)}
	if key == "" || secret == "" {
		r.Status = "missing"
		r.Message = "Options Intelligence uses the existing Alpaca credentials."
		return r
	}
	o, err := fetchOptionsContext(ctx, key, secret, "SPY", mode, 0)
	if err != nil {
		r.Status = "failed"
		r.Message = err.Error()
		return r
	}
	r.OK = true
	r.Status = "connected"
	r.Message = "Options Intelligence access is working."
	r.Details = []string{"Feed: " + o.Feed, "Real snapshot data only", "OPRA entitlement is tested; AUTO may fall back to indicative"}
	return r
}

func validationSnapshotKey(x SignalSnapshot) string {
	if strings.TrimSpace(x.ID) != "" {
		return "id:" + strings.TrimSpace(x.ID)
	}
	return fmt.Sprintf("legacy:%s|%s|%d", normalizeSymbol(x.Symbol), strings.ToLower(strings.TrimSpace(x.Horizon)), x.Timestamp)
}

// refreshValidationOutcomeState persists only post-decision analytics fields back
// into the canonical Signal Validation state. Frozen decision/evidence inputs are
// never overwritten by later provider data. The merge prevents a concurrent new
// snapshot from being lost while outcome evaluation is running.
func (e *Engine) refreshValidationOutcomeState() {
	e.mu.RLock()
	if e.mode != "live" || len(e.signalValidation.Snapshots) == 0 {
		e.mu.RUnlock()
		return
	}
	base := clone(e.signalValidation)
	bars := clone(e.bars)
	actions := clone(e.corporateActions)
	e.mu.RUnlock()

	evaluated := evaluateSignalSnapshotsProfessionalWithActions(base, bars, actions, "live")
	byKey := make(map[string]SignalSnapshot, len(evaluated.Snapshots))
	for _, x := range evaluated.Snapshots {
		byKey[validationSnapshotKey(x)] = x
	}

	e.mu.Lock()
	changed := false
	changedSnapshots := make([]SignalSnapshot, 0)
	for i := range e.signalValidation.Snapshots {
		cur := &e.signalValidation.Snapshots[i]
		x, ok := byKey[validationSnapshotKey(*cur)]
		if !ok {
			continue
		}
		rowChanged := cur.OutcomeState != x.OutcomeState || cur.OutcomeUpdatedAt != x.OutcomeUpdatedAt || cur.EntryTouchedAt != x.EntryTouchedAt || cur.TargetTouchedAt != x.TargetTouchedAt || cur.InvalidationAt != x.InvalidationAt || cur.MFE != x.MFE || cur.MAE != x.MAE || len(cur.Outcomes) != len(x.Outcomes)
		if rowChanged {
			changed = true
		}
		cur.Outcomes = clone(x.Outcomes)
		cur.MFE = x.MFE
		cur.MAE = x.MAE
		cur.OutcomeState = x.OutcomeState
		cur.OutcomeDetail = x.OutcomeDetail
		cur.OutcomeUpdatedAt = x.OutcomeUpdatedAt
		cur.EntryTouched = x.EntryTouched
		cur.EntryTouchedAt = x.EntryTouchedAt
		cur.TargetTouchedAt = x.TargetTouchedAt
		cur.InvalidationAt = x.InvalidationAt
		cur.ElapsedMinutes = x.ElapsedMinutes
		cur.OutcomeAdjustmentFactor = x.OutcomeAdjustmentFactor
		cur.OutcomeAdjustmentDetail = x.OutcomeAdjustmentDetail
		cur.LegacyPartial = x.LegacyPartial
		if rowChanged {
			changedSnapshots = append(changedSnapshots, clone(*cur))
		}
	}
	if changed {
		e.signalValidation.UpdatedAt = evaluated.UpdatedAt
		e.signalValidation.Message = evaluated.Message
	}
	e.mu.Unlock()
	for _, snapshot := range changedSnapshots {
		e.enqueueSignalSnapshotPersistence(snapshot)
	}
}
