package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"
	"time"
)

func TestV142ReleaseIdentity(t *testing.T) {
	if appVersion == "" {
		t.Fatal("runtime version identity must not be empty")
	}
	if buildID == "" {
		t.Fatal("runtime build identity must not be empty")
	}
}

func TestEventLifecycleIncludesMarketReaction(t *testing.T) {
	now := time.Now()
	cases := []struct {
		at   time.Time
		want string
	}{
		{now.Add(time.Minute), "UPCOMING"},
		{now.Add(-10 * time.Second), "RELEASED"},
		{now.Add(-2 * time.Minute), "MARKET REACTION"},
		{now.Add(-2 * time.Hour), "RESOLVED"},
	}
	for _, c := range cases {
		if got := eventLifecycle(c.at, now); got != c.want {
			t.Fatalf("%v got %q want %q", c.at, got, c.want)
		}
	}
}

func TestGuidanceParserRequiresExplicitGuidance(t *testing.T) {
	rl, rh, el, eh, okr, oke := parseGuidanceRanges("Company raises guidance and expects revenue $10.2 billion to $10.8 billion and EPS $4.10 to $4.30.")
	if !okr || !oke || rl != 10.2e9 || rh != 10.8e9 || el != 4.10 || eh != 4.30 {
		t.Fatalf("bad parse %v %v %v %v %v %v", rl, rh, el, eh, okr, oke)
	}
	_, _, _, _, okr, oke = parseGuidanceRanges("Revenue was $10.2 billion to $10.8 billion in a historical table.")
	if okr || oke {
		t.Fatal("must not infer guidance from unrelated ranges")
	}
}

func TestMetricHorizonChanges(t *testing.T) {
	vals := make([]float64, 70)
	for i := range vals {
		vals[i] = 100 - float64(i)
	}
	m, ok := metricWithChanges("X", "X", "%", "QA", "OFFICIAL", vals, time.Now().UnixMilli())
	if !ok {
		t.Fatal("no metric")
	}
	if m.Change5D != 5 || m.Change20D != 20 || m.Change1M != 22 || m.Change3M != 66 {
		t.Fatalf("changes %#v", m)
	}
}

func TestTwelveDataDirectProviderAndFutures(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/symbol_search", func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query().Get("symbol")
		typ := "Index"
		sym := "IDX"
		if strings.Contains(strings.ToLower(q), "future") {
			typ = "Futures"
			sym = "FUT"
		}
		_ = json.NewEncoder(w).Encode(map[string]any{"data": []map[string]any{{"symbol": sym, "instrument_name": q, "exchange": "QA", "instrument_type": typ}}})
	})
	mux.HandleFunc("/quote", func(w http.ResponseWriter, r *http.Request) {
		_ = json.NewEncoder(w).Encode(map[string]any{"symbol": r.URL.Query().Get("symbol"), "close": "123.45", "percent_change": "1.25", "timestamp": time.Now().Unix()})
	})
	ts := httptest.NewServer(mux)
	defer ts.Close()
	old := twelveDataBaseURL
	twelveDataBaseURL = ts.URL
	defer func() { twelveDataBaseURL = old }()
	p := twelveDataProvider{apiKey: "qa"}
	d, err := p.Refresh(context.Background())
	if err != nil || len(d) == 0 {
		t.Fatalf("direct %v %v", len(d), err)
	}
	f, err := p.RefreshFutures(context.Background())
	if err != nil || len(f) != 3 {
		t.Fatalf("futures %v %v", len(f), err)
	}
	if f["es_future"].Provenance != "DIRECT PROVIDER" {
		t.Fatalf("%#v", f["es_future"])
	}
}

func TestOfficialMacroAdapters(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/treasury", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, `<feed><entry><d:BC_2YEAR>4.10</d:BC_2YEAR><d:BC_10YEAR>4.30</d:BC_10YEAR><d:BC_30YEAR>4.90</d:BC_30YEAR></entry><entry><d:BC_2YEAR>4.00</d:BC_2YEAR><d:BC_10YEAR>4.20</d:BC_10YEAR><d:BC_30YEAR>4.80</d:BC_30YEAR></entry></feed>`)
	})
	mux.HandleFunc("/bls/CUUR0000SA0", blsFixture([]string{"320", "319", "318", "317", "316", "315", "314", "313", "312", "311", "310", "309", "308", "307"}))
	mux.HandleFunc("/bls/CUUR0000SA0L1E", blsFixture([]string{"330", "329", "328", "327", "326", "325", "324", "323", "322", "321", "320", "319", "318", "317"}))
	mux.HandleFunc("/bls/CES0000000001", blsFixture([]string{"160000", "159800"}))
	mux.HandleFunc("/bls/LNS14000000", blsFixture([]string{"4.2", "4.1"}))
	mux.HandleFunc("/eia/PET.RWTC.D", eiaFixture())
	mux.HandleFunc("/eia/PET.RBRTE.D", eiaFixture())
	mux.HandleFunc("/bea/glance/gdp", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "Real gross domestic product increased 2.8 percent in the latest release.")
	})
	mux.HandleFunc("/bea/glance/pce", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "The PCE price index increased 2.5 percent.")
	})
	ts := httptest.NewServer(mux)
	defer ts.Close()
	oldT, oldB, oldE, oldBEA := treasuryXMLBaseURL, blsAPIBaseURL, eiaAPIBaseURL, beaNewsBaseURL
	treasuryXMLBaseURL = ts.URL + "/treasury"
	blsAPIBaseURL = ts.URL + "/bls"
	eiaAPIBaseURL = ts.URL + "/eia"
	beaNewsBaseURL = ts.URL + "/bea"
	defer func() { treasuryXMLBaseURL, blsAPIBaseURL, eiaAPIBaseURL, beaNewsBaseURL = oldT, oldB, oldE, oldBEA }()
	app := newTestApplication(t)
	e := app.engine
	e.refreshOfficialMacroActuals(context.Background(), "qa")
	e.mu.RLock()
	defer e.mu.RUnlock()
	for _, k := range []string{"UST10Y", "CPI_INDEX", "NONFARM", "BEA_GDP", "WTI_OFFICIAL"} {
		if _, ok := e.macroMetrics[k]; !ok {
			t.Fatalf("missing %s %#v", k, e.macroMetrics)
		}
	}
	if got := e.health["macro-rates"]; !strings.HasPrefix(got, "healthy · official Treasury core rates") {
		t.Fatalf("official Treasury core rates should keep Macro Rates healthy: %q", got)
	}
}

func blsFixture(vals []string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		data := []map[string]string{}
		for i, v := range vals {
			data = append(data, map[string]string{"year": "2026", "period": fmt.Sprintf("M%02d", 12-(i%12)), "value": v})
		}
		_ = json.NewEncoder(w).Encode(map[string]any{"status": "REQUEST_SUCCEEDED", "Results": map[string]any{"series": []map[string]any{{"seriesID": "x", "data": data}}}})
	}
}
func eiaFixture() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		_ = json.NewEncoder(w).Encode(map[string]any{"response": map[string]any{"data": []map[string]any{{"value": 80.0}, {"value": 79.0}, {"value": 78.0}}}})
	}
}

func TestAffectedContext(t *testing.T) {
	s, sec := affectedContext(MacroEvent{Region: "US", Name: "Consumer Price Index"}, []string{"NVDA"})
	if !contains(s, "NVDA") || !contains(s, "TLT") || !contains(sec, "Technology") {
		t.Fatalf("%v %v", s, sec)
	}
}

func TestTWSEOfficialPublicClose(t *testing.T) {
	old := twsePublicCloseURL
	defer func() { twsePublicCloseURL = old }()
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, `<html><body>115年08月07日 價格指數 發行量加權股價指數 44,225.91 - 170.79 -0.38 </body></html>`)
	}))
	defer srv.Close()
	twsePublicCloseURL = srv.URL
	d, err := twseOfficialClose(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if d.Key != "taiwan_official_close" || d.Provenance != "OFFICIAL CLOSE" || d.IsProxy || d.Value != 44225.91 || d.ChangePercent != -0.38 {
		t.Fatalf("unexpected TWSE driver: %+v", d)
	}
	if d.Session != "OFFICIAL CLOSE" {
		t.Fatalf("session = %q", d.Session)
	}
}

func TestV142GlobalRoutingRetainsOfficialAndProxyWithoutDirectLeak(t *testing.T) {
	now := time.Now().UnixMilli()
	quotes := map[string]Quote{"EWT": {Symbol: "EWT", Price: 80, ChangePercent: 1.1, UpdatedAt: now, DataState: "live", Source: "qa"}}
	direct := map[string]GlobalDriver{
		"taiwan":                {Key: "taiwan", Label: "TAIEX direct", Value: 45000, ChangePercent: 1.4, Provenance: "DIRECT PROVIDER", Confidence: 90},
		"taiwan_official_close": {Key: "taiwan_official_close", Label: "TAIEX official close", Value: 44225, ChangePercent: -0.38, Provenance: "OFFICIAL CLOSE", Confidence: 95},
	}
	auto := deriveGlobalMarketContext(quotes, direct, nil, "auto")
	if _, ok := auto.Drivers["taiwan"]; !ok {
		t.Fatal("AUTO missing direct Taiwan")
	}
	if _, ok := auto.Drivers["taiwan_official_close"]; !ok {
		t.Fatal("AUTO dropped official close")
	}
	if _, ok := auto.Drivers["taiwan_proxy"]; !ok {
		t.Fatalf("AUTO dropped supplementary live proxy: %#v", auto.Drivers)
	}
	free := deriveGlobalMarketContext(quotes, direct, nil, "free-first")
	if d, ok := free.Drivers["taiwan"]; ok && strings.Contains(strings.ToUpper(d.Provenance), "DIRECT") {
		t.Fatal("FREE FIRST leaked paid/direct primary")
	}
	if _, ok := free.Drivers["taiwan_official_close"]; !ok {
		t.Fatal("FREE FIRST missing official close")
	}
	if d, ok := free.Drivers["taiwan"]; !ok || !d.IsProxy {
		t.Fatalf("FREE FIRST should retain current real proxy: %#v", free.Drivers)
	}
}

func TestV1422ParseCboeVIXHistory(t *testing.T) {
	raw := "DATE,OPEN,HIGH,LOW,CLOSE\n08/07/2026,18.10,18.80,17.90,18.42\n08/08/2026,18.40,19.10,18.20,18.90\n"
	bars, err := parseCboeVIXHistory(raw)
	if err != nil {
		t.Fatal(err)
	}
	if len(bars) != 2 || bars[0].C != 18.42 || bars[1].C != 18.90 || bars[1].T <= bars[0].T {
		t.Fatalf("unexpected Cboe VIX bars: %+v", bars)
	}
}

func TestV1422TwelveDataVIXRejectsProxyAndAcceptsCanonicalIndex(t *testing.T) {
	old := twelveDataBaseURL
	defer func() { twelveDataBaseURL = old; tdSymbolCache = sync.Map{} }()
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/symbol_search":
			_ = json.NewEncoder(w).Encode(map[string]any{"data": []map[string]any{
				{"symbol": "VIXY", "instrument_name": "ProShares VIX Short-Term Futures ETF", "exchange": "CBOE", "instrument_type": "ETF"},
				{"symbol": "VIX", "instrument_name": "Cboe Volatility Index", "exchange": "CBOE", "instrument_type": "Index"},
			}})
		case "/quote":
			if r.URL.Query().Get("symbol") != "VIX" {
				t.Fatalf("proxy symbol requested: %s", r.URL.Query().Get("symbol"))
			}
			_ = json.NewEncoder(w).Encode(map[string]any{"symbol": "VIX", "name": "Cboe Volatility Index", "exchange": "CBOE", "close": "18.50", "previous_close": "19.10", "open": "19.00", "high": "19.20", "low": "18.20", "timestamp": 1786320000})
		default:
			http.NotFound(w, r)
		}
	}))
	defer srv.Close()
	twelveDataBaseURL = srv.URL
	tdSymbolCache = sync.Map{}
	q, provider, err := tdVIXQuote(context.Background(), "test")
	if err != nil {
		t.Fatal(err)
	}
	if q.Symbol != "VIX" || q.Close != "18.50" || !strings.Contains(provider, "CBOE") {
		t.Fatalf("unexpected Twelve Data VIX result: %+v %q", q, provider)
	}
}
