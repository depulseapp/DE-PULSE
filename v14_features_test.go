package main

import (
	"context"
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"
	"time"
)

func TestV14GlobalContextUsesTruthfulProxyAndDoesNotMutateQuotes(t *testing.T) {
	now := time.Now().UnixMilli()
	quotes := map[string]Quote{
		"SPY": {Symbol: "SPY", Price: 600, ChangePercent: 1.1, Source: "finnhub-websocket", DataState: "live", UpdatedAt: now},
		"QQQ": {Symbol: "QQQ", Price: 520, ChangePercent: 1.4, Source: "finnhub-websocket", DataState: "live", UpdatedAt: now},
		"IWM": {Symbol: "IWM", Price: 220, ChangePercent: .9, Source: "finnhub-websocket", DataState: "live", UpdatedAt: now},
		"EWT": {Symbol: "EWT", Price: 60, ChangePercent: 1.2, Source: "alpaca-iex", DataState: "live", UpdatedAt: now},
		"SMH": {Symbol: "SMH", Price: 300, ChangePercent: 1.5, Source: "alpaca-iex", DataState: "live", UpdatedAt: now},
		"VIX": {Symbol: "VIX", Price: 14.8, ChangePercent: -3, Source: "vix-index", DataState: "live", UpdatedAt: now},
	}
	before := clone(quotes)
	g := deriveGlobalMarketContext(quotes, nil, nil, "auto")
	if g.Tone != "RISK-ON" && g.Tone != "CONSTRUCTIVE" {
		t.Fatalf("unexpected tone: %+v", g)
	}
	d := g.Drivers["taiwan"]
	if d.Provenance != "LIVE PROXY" || !strings.Contains(d.Detail, "EWT proxy") {
		t.Fatalf("proxy provenance not truthful: %+v", d)
	}
	b1, _ := json.Marshal(before)
	b2, _ := json.Marshal(quotes)
	if string(b1) != string(b2) {
		t.Fatal("context derivation mutated canonical quote state")
	}
}

func TestV14OfficialEventParserKeepsUnknownTimeDateOnly(t *testing.T) {
	fixture := `<html><body><div>Consumer Price Index September 11, 2026</div><div>Employment Situation October 2, 2026</div></body></html>`
	got := parseOfficialEvents(fixture, "US", "BLS", "https://www.bls.gov", []string{"Consumer Price Index", "Employment Situation"})
	if len(got) != 2 {
		t.Fatalf("events=%d %+v", len(got), got)
	}
	for _, e := range got {
		if e.TimeKnown || e.StartsAt != 0 {
			t.Fatalf("invented exact release time: %+v", e)
		}
		if e.Source != "BLS" || e.Impact != "HIGH" {
			t.Fatalf("bad event: %+v", e)
		}
	}
}

func TestV14EventModeTMinus15AndReaction(t *testing.T) {
	now := time.Now()
	ev := MacroEvent{ID: "cpi", Name: "CPI", Impact: "HIGH", StartsAt: now.Add(10 * time.Minute).UnixMilli(), TimeKnown: true}
	pre := eventModeFor([]MacroEvent{ev}, now, true)
	if !pre.Active || pre.Phase != "PREP" || pre.CountdownS <= 0 {
		t.Fatalf("pre=%+v", pre)
	}
	react := eventModeFor([]MacroEvent{{ID: "cpi", Name: "CPI", Impact: "HIGH", StartsAt: now.Add(-30 * time.Second).UnixMilli(), TimeKnown: true}}, now, true)
	if !react.Active || react.Phase != "REACTION" {
		t.Fatalf("react=%+v", react)
	}
	if eventModeFor([]MacroEvent{ev}, now, false).Active {
		t.Fatal("disabled event mode became active")
	}
}

func TestV14OptionsAutoFallsBackToIndicativeAndNeverFabricatesOI(t *testing.T) {
	oldData, oldTrade := alpacaDataBaseURL, alpacaTradingBaseURL
	defer func() { alpacaDataBaseURL, alpacaTradingBaseURL = oldData, oldTrade }()
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Query().Get("feed") == "opra" {
			http.Error(w, "subscription required", http.StatusForbidden)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		exp := time.Now().AddDate(0, 1, 0).Format("060102")
		payload := fmt.Sprintf(`{"snapshots":{"SPY%sC00600000":{"dailyBar":{"v":1200},"impliedVolatility":0.20},"SPY%sP00600000":{"dailyBar":{"v":600},"impliedVolatility":0.24}}}`, exp, exp)
		_, _ = w.Write([]byte(payload))
	}))
	defer srv.Close()
	alpacaDataBaseURL, alpacaTradingBaseURL = srv.URL, srv.URL
	o, err := fetchOptionsContext(context.Background(), "k", "s", "SPY", "auto", 600)
	if err != nil {
		t.Fatal(err)
	}
	if o.Feed != "INDICATIVE" || o.State != "DELAYED/INDICATIVE" {
		t.Fatalf("fallback=%+v", o)
	}
	if o.CallVolume != 1200 || o.PutVolume != 600 || o.PutCallVolume != .5 {
		t.Fatalf("aggregation=%+v", o)
	}
	if len(o.Limitations) == 0 || !strings.Contains(strings.ToLower(strings.Join(o.Limitations, " ")), "open interest") {
		t.Fatal("missing truthful OI limitation")
	}
}

func TestV14OptionsAggregationHorizonInputsAreContextOnly(t *testing.T) {
	p := alpacaOptionChainResponse{Snapshots: map[string]struct {
		DailyBar *struct {
			V float64 `json:"v"`
		} `json:"dailyBar"`
		ImpliedVolatility float64 `json:"impliedVolatility"`
		Greeks            *struct {
			Gamma float64 `json:"gamma"`
		} `json:"greeks"`
	}{}}
	call := struct {
		DailyBar *struct {
			V float64 `json:"v"`
		} `json:"dailyBar"`
		ImpliedVolatility float64 `json:"impliedVolatility"`
		Greeks            *struct {
			Gamma float64 `json:"gamma"`
		} `json:"greeks"`
	}{ImpliedVolatility: .3}
	call.DailyBar = &struct {
		V float64 `json:"v"`
	}{V: 2000}
	put := call
	put.DailyBar = &struct {
		V float64 `json:"v"`
	}{V: 400}
	exp := time.Now().AddDate(0, 1, 0).Format("060102")
	p.Snapshots["NVDA"+exp+"C00200000"] = call
	p.Snapshots["NVDA"+exp+"P00200000"] = put
	far := call
	far.ImpliedVolatility = 2.5
	p.Snapshots["NVDA"+exp+"C00400000"] = far
	o := aggregateOptions("NVDA", "opra", 200, p, nil, fmt.Errorf("open interest unavailable in legacy aggregation test"))
	if o.Bias != "BULLISH" {
		t.Fatalf("bias=%s", o.Bias)
	}
	if o.ExpectedMove <= 0 || o.ExpectedMove > 40 {
		t.Fatalf("expected move should use nearest-expiry near-the-money IV, got %.2f", o.ExpectedMove)
	}
}

func TestV14SignalValidationDedupeAndRealOutcomeRule(t *testing.T) {
	e := &Engine{signalValidation: SignalValidationState{Snapshots: []SignalSnapshot{}}}
	ts := time.Now().Add(-12 * 24 * time.Hour).UnixMilli()
	x := SignalSnapshot{Symbol: "NVDA", Horizon: "swing", Timestamp: ts, Price: 100, Score: 75, Action: "BUY", Readiness: "READY"}
	e.recordSignalSnapshot(x)
	e.recordSignalSnapshot(SignalSnapshot{Symbol: "NVDA", Horizon: "swing", Timestamp: ts + 5*60*1000, Price: 101})
	if len(e.signalValidation.Snapshots) != 1 {
		t.Fatalf("dedupe failed: %d", len(e.signalValidation.Snapshots))
	}
	bars := map[string]map[string][]Bar{"NVDA": {"daily": []Bar{
		{T: time.UnixMilli(ts).Add(24 * time.Hour).Unix(), H: 106, L: 98, C: 104},
		{T: time.UnixMilli(ts).Add(2 * 24 * time.Hour).Unix(), H: 110, L: 102, C: 108},
		{T: time.UnixMilli(ts).Add(3 * 24 * time.Hour).Unix(), H: 112, L: 104, C: 110},
		{T: time.UnixMilli(ts).Add(4 * 24 * time.Hour).Unix(), H: 114, L: 105, C: 112},
		{T: time.UnixMilli(ts).Add(5 * 24 * time.Hour).Unix(), H: 116, L: 106, C: 115},
	}}}
	demo := evaluateSignalSnapshotsProfessionalWithActions(clone(e.signalValidation), bars, nil, "demo")
	if len(demo.Snapshots[0].Outcomes) != 0 {
		t.Fatal("demo outcomes leaked into validation evidence")
	}
	live := evaluateSignalSnapshotsProfessionalWithActions(clone(e.signalValidation), bars, nil, "live")
	if math.Abs(live.Snapshots[0].Outcomes["5D"]-15) > 1e-9 {
		t.Fatalf("5D=%v", live.Snapshots[0].Outcomes)
	}
	if math.Abs(live.Snapshots[0].MFE-16) > 1e-9 || math.Abs(live.Snapshots[0].MAE+2) > 1e-9 {
		t.Fatalf("excursions=%v/%v", live.Snapshots[0].MFE, live.Snapshots[0].MAE)
	}
}

func TestV14CompactAIEvidenceIsTickerSpecificAndBounded(t *testing.T) {
	snap := RuntimeSnapshot{Status: "running", Mode: "live", Quotes: map[string]Quote{"AAA": {Symbol: "AAA", Price: 100}, "BBB": {Symbol: "BBB", Price: 200}}, Fundamentals: map[string]FundamentalSnapshot{"AAA": {Symbol: "AAA", RevenueGrowth: 10}, "BBB": {Symbol: "BBB", RevenueGrowth: 99}}, Options: map[string]OptionsContext{"AAA": {Symbol: "AAA", Bias: "BULLISH"}}, News: []NewsItem{{Headline: strings.Repeat("A", 8000), Symbols: []string{"AAA"}}, {Headline: "BBB only", Symbols: []string{"BBB"}}}, Earnings: []EarningsItem{{Symbol: "AAA"}, {Symbol: "BBB"}}, Filings: []FilingItem{{Symbol: "AAA", Description: strings.Repeat("F", 8000)}, {Symbol: "BBB"}}}
	e := compactAIEvidence(AIRequest{Ticker: "AAA", Kind: "risk", ClientContext: map[string]any{"horizon": "day"}}, snap)
	if q := e["quote"].(Quote); q.Symbol != "AAA" {
		t.Fatal("wrong quote")
	}
	b := marshalBoundedContext(e, 12000)
	if len(b) > 12000 {
		t.Fatalf("context too large: %d", len(b))
	}
	if strings.Contains(string(b), "BBB only") {
		t.Fatal("cross-ticker evidence leaked")
	}
	if !isContextLimitError(assertErr("TPM token limit exceeded")) {
		t.Fatal("context-limit detector")
	}
}

type assertErr string

func (e assertErr) Error() string { return string(e) }

func TestV14SettingsMigrationAddsNewModesWithoutChangingLegacyDeskState(t *testing.T) {
	st := defaultState()
	st.Version = 16
	st.Settings.GlobalProviderMode = ""
	st.Settings.OptionsDataMode = ""
	st.Settings.MacroEventModeEnabled = false
	st.Settings.DayEnabled = false
	st.Settings.SwingEnabled = true
	st.Settings.LongEnabled = true
	got := mergeState(st)
	if got.Version != 17 || got.Settings.GlobalProviderMode != "auto" || got.Settings.OptionsDataMode != "auto" || !got.Settings.MacroEventModeEnabled {
		t.Fatalf("migration=%+v", got.Settings)
	}
	if got.Settings.DayEnabled {
		t.Fatal("legacy engine state changed")
	}
}

func TestV14OfficialEventParserUsesExactTimeOnlyWhenPresent(t *testing.T) {
	fixture := `<div>Consumer Price Index September 11, 2026 8:30 a.m. ET</div>`
	got := parseOfficialEvents(fixture, "US", "BLS", "https://www.bls.gov", []string{"Consumer Price Index"})
	if len(got) != 1 || !got[0].TimeKnown || got[0].StartsAt <= 0 {
		t.Fatalf("exact official time not preserved: %+v", got)
	}
}

func TestV14BroadBreadthUsesBroadUniverseNotWatchlist(t *testing.T) {
	q := map[string]Quote{}
	now := time.Now().UnixMilli()
	for _, s := range broadBreadthUniverse {
		q[s] = Quote{Symbol: s, Price: 100, ChangePercent: 1, UpdatedAt: now, DataState: "live"}
	}
	d, ok := broadBreadthDriver(q)
	if !ok || d.State != "SUPPORTIVE" || !strings.Contains(d.Detail, "not watchlist breadth") {
		t.Fatalf("breadth=%+v", d)
	}
}

func TestV14GlobalProviderModeHonored(t *testing.T) {
	q := map[string]Quote{"EWY": {Symbol: "EWY", Price: 50, ChangePercent: 2, DataState: "live"}}
	if _, ok := deriveGlobalMarketContext(q, nil, nil, "direct").Drivers["korea"]; ok {
		t.Fatal("Direct Only leaked ETF proxy")
	}
	if _, ok := deriveGlobalMarketContext(q, nil, nil, "proxy").Drivers["korea"]; !ok {
		t.Fatal("Proxy Only did not use real proxy")
	}
}

func TestV14EventModePrewarmIsDeduplicatedInImplementation(t *testing.T) {
	s := productionGoSourceForTest(t)
	if !strings.Contains(s, "warmed[mode.EventID]") {
		t.Fatal("high-impact prewarm can repeat and cause API storms")
	}
}

func TestV14FREDRefreshParsesOfficialValueField(t *testing.T) {
	old := fredAPIBaseURL
	defer func() { fredAPIBaseURL = old }()
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"observations":[{"date":"2026-08-08","value":"4.25"}]}`))
	}))
	defer srv.Close()
	fredAPIBaseURL = srv.URL
	e := &Engine{macroMetrics: map[string]MacroMetric{}, lastUpdated: map[string]int64{}, health: map[string]string{}}
	e.refreshFRED(context.Background(), "test-key")
	m, ok := e.macroMetrics["DGS10"]
	if !ok || math.Abs(m.Value-4.25) > 1e-9 {
		t.Fatalf("FRED value parse failed: %+v", e.macroMetrics)
	}
	if m.Provenance != "OFFICIAL" {
		t.Fatalf("bad provenance: %+v", m)
	}
}

func TestV14StressSixtyMovingSymbolsNoDroppedCanonicalState(t *testing.T) {
	a := newTestApplication(t)
	e := a.engine
	var wg sync.WaitGroup
	start := time.Now()
	for i := 0; i < 60; i++ {
		sym := fmt.Sprintf("S%02d", i)
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < 25; j++ {
				e.updateQuote(sym, Quote{Price: 100 + float64(j), PreviousClose: 99, ProviderTimestamp: time.Now().UnixMilli()}, "qa-stress")
			}
		}()
	}
	for i := 0; i < 20; i++ {
		_ = e.Snapshot()
	}
	wg.Wait()
	snap := e.Snapshot()
	count := 0
	for i := 0; i < 60; i++ {
		if snap.Quotes[fmt.Sprintf("S%02d", i)].Price > 0 {
			count++
		}
	}
	if count != 60 {
		t.Fatalf("canonical moving-symbol state dropped updates: %d/60", count)
	}
	limit := 5 * time.Second
	if raceBuild {
		// The race detector intentionally instruments every shared-memory access;
		// retain the original 5s production-path gate while allowing instrumentation overhead.
		limit = 15 * time.Second
	}
	if time.Since(start) > limit {
		t.Fatalf("stress path unexpectedly slow: %s (limit %s, race=%v)", time.Since(start), limit, raceBuild)
	}
}

func TestV14OptionsIVChangeUsesOnlyConsecutiveRealSnapshots(t *testing.T) {
	cur := applyOptionsIVChange(OptionsContext{}, OptionsContext{AverageIV: 0.42})
	if cur.IVChange != 0 {
		t.Fatalf("first snapshot must not invent IV change: %+v", cur)
	}
	cur = applyOptionsIVChange(OptionsContext{AverageIV: 0.42}, OptionsContext{AverageIV: 0.47})
	if math.Abs(cur.IVChange-5) > 0.0001 {
		t.Fatalf("expected +5.0 percentage points, got %.4f", cur.IVChange)
	}
}
