package main

import (
	"context"
	"encoding/csv"
	"fmt"
	"net/http"
	"net/url"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"
)

// v14.0.3 compliance layer. These providers feed canonical shared context only;
// they never mutate deterministic Day/Swing/Long score/action formulas.

type DirectGlobalProvider interface {
	Name() string
	Refresh(ctx context.Context) (map[string]GlobalDriver, error)
}

type DirectFuturesProvider interface {
	Name() string
	RefreshFutures(ctx context.Context) (map[string]GlobalDriver, error)
}

// Premium-ready capability contracts. Paid entitlements remain optional; these
// interfaces keep future providers out of the UI/domain model and let the free
// stack remain the truthful fallback.
type FastMacroRelease struct {
	EventID    string
	Consensus  *float64
	Actual     *float64
	Unit       string
	ReleasedAt int64
	Source     string
}
type FastMacroProvider interface {
	Name() string
	FetchRelease(ctx context.Context, event MacroEvent) (FastMacroRelease, error)
}
type OptionsIntelligenceProvider interface {
	Name() string
	Snapshot(ctx context.Context, symbol, mode string, underlying float64) (OptionsContext, error)
}
type alpacaOptionsProvider struct{ key, secret string }

func (p alpacaOptionsProvider) Name() string { return "Alpaca Options" }
func (p alpacaOptionsProvider) Snapshot(ctx context.Context, symbol, mode string, underlying float64) (OptionsContext, error) {
	return fetchOptionsContext(ctx, p.key, p.secret, symbol, mode, underlying)
}

type twelveDataProvider struct{ apiKey string }

func (p twelveDataProvider) Name() string { return "Twelve Data" }

var twelveDataBaseURL = "https://api.twelvedata.com"
var tdSymbolCache sync.Map // key: base URL + future flag + query; value: symbol|exchange

var twelveGlobalSearch = map[string]struct{ Label, Query, Group string }{
	"korea":     {"KOSPI / Korea", "KOSPI", "asia"},
	"taiwan":    {"TAIEX / Taiwan", "TAIEX", "asia-tech"},
	"japan":     {"Nikkei / Japan", "Nikkei 225", "asia"},
	"hong_kong": {"Hang Seng / Hong Kong", "Hang Seng", "china"},
	"china":     {"China Composite", "Shanghai Composite", "china"},
	"europe":    {"STOXX Europe", "STOXX Europe 600", "europe"},
	"dax":       {"DAX", "DAX", "europe"},
	"ftse":      {"FTSE 100", "FTSE 100", "europe"},
}
var twelveFutureSearch = map[string]struct{ Label, Query string }{
	"es_future":  {"S&P 500 E-mini Futures", "ES futures"},
	"nq_future":  {"Nasdaq 100 E-mini Futures", "NQ futures"},
	"rty_future": {"Russell 2000 Futures", "RTY futures"},
}

type tdSearchResponse struct {
	Data []struct {
		Symbol         string `json:"symbol"`
		InstrumentName string `json:"instrument_name"`
		Exchange       string `json:"exchange"`
		Country        string `json:"country"`
		Type           string `json:"instrument_type"`
	} `json:"data"`
	Status  string `json:"status"`
	Message string `json:"message"`
}
type tdQuoteResponse struct {
	Symbol        string `json:"symbol"`
	Name          string `json:"name"`
	Exchange      string `json:"exchange"`
	Currency      string `json:"currency"`
	Close         string `json:"close"`
	Open          string `json:"open"`
	High          string `json:"high"`
	Low           string `json:"low"`
	PreviousClose string `json:"previous_close"`
	PercentChange string `json:"percent_change"`
	Datetime      string `json:"datetime"`
	Timestamp     int64  `json:"timestamp"`
	Status        string `json:"status"`
	Message       string `json:"message"`
}

func tdPickSymbol(ctx context.Context, key, query string, future bool) (string, string, error) {
	cacheKey := fmt.Sprintf("%s|%t|%s", twelveDataBaseURL, future, strings.ToLower(strings.TrimSpace(query)))
	if v, ok := tdSymbolCache.Load(cacheKey); ok {
		parts := strings.SplitN(v.(string), "|", 2)
		if len(parts) == 2 {
			return parts[0], parts[1], nil
		}
	}
	var sr tdSearchResponse
	raw := twelveDataBaseURL + "/symbol_search?symbol=" + url.QueryEscape(query) + "&apikey=" + url.QueryEscape(key) + "&outputsize=20"
	if err := getJSON(ctx, &http.Client{Timeout: 12 * time.Second}, raw, nil, &sr); err != nil {
		return "", "", err
	}
	if strings.EqualFold(sr.Status, "error") {
		return "", "", fmt.Errorf("Twelve Data: %s", sr.Message)
	}
	for _, x := range sr.Data {
		typ := strings.ToLower(x.Type)
		name := strings.ToLower(x.InstrumentName + " " + x.Symbol)
		if future {
			if strings.Contains(typ, "future") || strings.Contains(name, "future") {
				tdSymbolCache.Store(cacheKey, x.Symbol+"|"+x.Exchange)
				return x.Symbol, x.Exchange, nil
			}
		} else if strings.Contains(typ, "index") || strings.Contains(name, "index") || strings.Contains(strings.ToLower(query), strings.ToLower(x.Symbol)) {
			tdSymbolCache.Store(cacheKey, x.Symbol+"|"+x.Exchange)
			return x.Symbol, x.Exchange, nil
		}
	}
	if len(sr.Data) > 0 {
		tdSymbolCache.Store(cacheKey, sr.Data[0].Symbol+"|"+sr.Data[0].Exchange)
		return sr.Data[0].Symbol, sr.Data[0].Exchange, nil
	}
	return "", "", fmt.Errorf("no matching instrument for %s", query)
}
func tdQuote(ctx context.Context, key, symbol string) (tdQuoteResponse, error) {
	var q tdQuoteResponse
	raw := twelveDataBaseURL + "/quote?symbol=" + url.QueryEscape(symbol) + "&apikey=" + url.QueryEscape(key)
	if err := getJSON(ctx, &http.Client{Timeout: 12 * time.Second}, raw, nil, &q); err != nil {
		return q, err
	}
	if strings.EqualFold(q.Status, "error") {
		return q, fmt.Errorf("Twelve Data: %s", q.Message)
	}
	return q, nil
}

// tdVIXQuote resolves only a canonical Cboe VIX index result. It deliberately
// rejects ETFs, futures and generic volatility products so a proxy can never be
// presented as the VIX spot index.
func tdVIXQuote(ctx context.Context, key string) (tdQuoteResponse, string, error) {
	if strings.TrimSpace(key) == "" {
		return tdQuoteResponse{}, "", fmt.Errorf("Twelve Data not configured")
	}
	var sr tdSearchResponse
	raw := twelveDataBaseURL + "/symbol_search?symbol=VIX&apikey=" + url.QueryEscape(key) + "&outputsize=20"
	if err := getJSON(ctx, &http.Client{Timeout: 12 * time.Second}, raw, nil, &sr); err != nil {
		return tdQuoteResponse{}, "", err
	}
	if strings.EqualFold(sr.Status, "error") {
		return tdQuoteResponse{}, "", fmt.Errorf("Twelve Data: %s", sr.Message)
	}
	for _, x := range sr.Data {
		typ := strings.ToLower(strings.TrimSpace(x.Type))
		name := strings.ToLower(strings.TrimSpace(x.InstrumentName))
		sym := strings.ToUpper(strings.TrimSpace(x.Symbol))
		ex := strings.ToUpper(strings.TrimSpace(x.Exchange))
		canonical := strings.Contains(typ, "index") && strings.Contains(name, "volatility") && (strings.Contains(name, "cboe") || sym == "VIX")
		if !canonical || strings.Contains(name, "future") || strings.Contains(name, "etf") || strings.Contains(typ, "future") || strings.Contains(typ, "etf") {
			continue
		}
		q, err := tdQuote(ctx, key, x.Symbol)
		if err != nil {
			continue
		}
		px, _ := strconv.ParseFloat(q.Close, 64)
		qname := strings.ToLower(q.Name)
		if px < 5 || px > 200 || (qname != "" && !strings.Contains(qname, "volatility") && !strings.Contains(strings.ToLower(q.Symbol), "vix")) {
			continue
		}
		providerSymbol := strings.TrimSpace(x.Symbol)
		if ex != "" {
			providerSymbol += " · " + ex
		}
		return q, providerSymbol, nil
	}
	return tdQuoteResponse{}, "", fmt.Errorf("canonical VIX index is not available for this Twelve Data entitlement")
}

var cboeVIXHistoryURL = "https://cdn.cboe.com/api/global/us_indices/daily_prices/VIX_History.csv"

// parseCboeVIXHistory parses Cboe's official daily VIX history. Cboe publishes
// DATE, OPEN, HIGH, LOW, CLOSE and updates the file daily. The returned bars are
// real index observations, never a volatility ETF/futures proxy.
func parseCboeVIXHistory(raw string) ([]Bar, error) {
	r := csv.NewReader(strings.NewReader(raw))
	records, err := r.ReadAll()
	if err != nil {
		return nil, err
	}
	if len(records) < 2 {
		return nil, fmt.Errorf("Cboe VIX history is empty")
	}
	out := make([]Bar, 0, len(records)-1)
	loc, _ := time.LoadLocation("America/New_York")
	for _, rec := range records[1:] {
		if len(rec) < 5 {
			continue
		}
		var d time.Time
		for _, layout := range []string{"01/02/2006", "2006-01-02", "1/2/2006"} {
			if t, e := time.ParseInLocation(layout, strings.TrimSpace(rec[0]), loc); e == nil {
				d = t
				break
			}
		}
		if d.IsZero() {
			continue
		}
		o, e1 := strconv.ParseFloat(strings.TrimSpace(rec[1]), 64)
		h, e2 := strconv.ParseFloat(strings.TrimSpace(rec[2]), 64)
		l, e3 := strconv.ParseFloat(strings.TrimSpace(rec[3]), 64)
		c, e4 := strconv.ParseFloat(strings.TrimSpace(rec[4]), 64)
		if e1 != nil || e2 != nil || e3 != nil || e4 != nil || c <= 0 {
			continue
		}
		closeTime := time.Date(d.Year(), d.Month(), d.Day(), 16, 0, 0, 0, loc)
		out = append(out, Bar{T: closeTime.Unix(), O: o, H: h, L: l, C: c})
	}
	if len(out) == 0 {
		return nil, fmt.Errorf("no valid Cboe VIX rows")
	}
	return out, nil
}

func fetchCboeVIXHistory(ctx context.Context) ([]Bar, error) {
	raw, err := fetchText(ctx, cboeVIXHistoryURL)
	if err != nil {
		return nil, err
	}
	return parseCboeVIXHistory(raw)
}
func tdDriver(ctx context.Context, key, k, label, query string, future bool) (GlobalDriver, error) {
	sym, ex, err := tdPickSymbol(ctx, key, query, future)
	if err != nil {
		return GlobalDriver{}, err
	}
	q, err := tdQuote(ctx, key, sym)
	if err != nil {
		return GlobalDriver{}, err
	}
	px, _ := strconv.ParseFloat(q.Close, 64)
	ch, _ := strconv.ParseFloat(q.PercentChange, 64)
	if px <= 0 {
		return GlobalDriver{}, fmt.Errorf("empty price for %s", sym)
	}
	ts := normalizeObservationMs(q.Timestamp)
	// Provider observation time stays unknown when Twelve Data does not supply it;
	// the engine records the independent DE.PULSE receipt time separately.
	session := "CURRENT"
	if future {
		session = "FUTURES"
	}
	return GlobalDriver{Key: k, Label: label, State: driverState(ch, false), Value: px, ChangePercent: ch, Source: "Twelve Data", Provenance: "DIRECT PROVIDER", UpdatedAt: ts, Confidence: 86, Detail: strings.TrimSpace(sym + " · " + ex), Session: session, ProviderSymbol: sym, IsProxy: false}, nil
}
func (p twelveDataProvider) Refresh(ctx context.Context) (map[string]GlobalDriver, error) {
	if strings.TrimSpace(p.apiKey) == "" {
		return nil, fmt.Errorf("not configured")
	}
	out := map[string]GlobalDriver{}
	errs := []string{}
	keys := make([]string, 0, len(twelveGlobalSearch))
	for k := range twelveGlobalSearch {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		s := twelveGlobalSearch[k]
		d, err := tdDriver(ctx, p.apiKey, k, s.Label, s.Query, false)
		if err != nil {
			errs = append(errs, k+": "+err.Error())
			continue
		}
		out[k] = d
	}
	if len(out) == 0 {
		return out, fmt.Errorf("direct global unavailable: %s", strings.Join(errs, "; "))
	}
	return out, nil
}
func (p twelveDataProvider) RefreshFutures(ctx context.Context) (map[string]GlobalDriver, error) {
	if strings.TrimSpace(p.apiKey) == "" {
		return nil, fmt.Errorf("not configured")
	}
	out := map[string]GlobalDriver{}
	errs := []string{}
	keys := make([]string, 0, len(twelveFutureSearch))
	for k := range twelveFutureSearch {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		s := twelveFutureSearch[k]
		d, err := tdDriver(ctx, p.apiKey, k, s.Label, s.Query, true)
		if err != nil {
			errs = append(errs, k+": "+err.Error())
			continue
		}
		d.State = driverState(d.ChangePercent, false)
		out[k] = d
	}
	if len(out) == 0 {
		return out, fmt.Errorf("direct futures unavailable: %s", strings.Join(errs, "; "))
	}
	return out, nil
}

func (e *Engine) refreshDirectGlobal(ctx context.Context, key string) {
	if strings.TrimSpace(key) == "" {
		e.setHealth("global-direct", "not configured · proxy/official fallback active")
		return
	}
	p := twelveDataProvider{apiKey: key}
	direct, err := p.Refresh(ctx)
	if direct == nil {
		direct = map[string]GlobalDriver{}
	}
	futures, ferr := p.RefreshFutures(ctx)
	for k, v := range futures {
		direct[k] = v
	}
	if err != nil && ferr != nil {
		e.setHealth("global-direct", "degraded · direct unavailable; real fallback active")
		return
	}
	if len(direct) == 0 {
		return
	}
	e.mu.Lock()
	for k, v := range direct {
		e.globalDirect[k] = v
	}
	e.lastUpdated["global-direct"] = time.Now().UnixMilli()
	e.health["global-direct"] = fmt.Sprintf("healthy · %d direct instruments", len(direct))
	e.mu.Unlock()
}

// Official/public macro actuals ------------------------------------------------
var treasuryXMLBaseURL = "https://home.treasury.gov/resource-center/data-chart-center/interest-rates/pages/xml"
var blsAPIBaseURL = "https://api.bls.gov/publicAPI/v2/timeseries/data"
var eiaAPIBaseURL = "https://api.eia.gov/v2/seriesid"
var beaNewsBaseURL = "https://www.bea.gov/news"

func metricWithChanges(key, label, unit, source, prov string, vals []float64, stamp int64) (MacroMetric, bool) {
	if len(vals) == 0 {
		return MacroMetric{}, false
	}
	m := MacroMetric{Key: key, Label: label, Value: vals[0], Unit: unit, Source: source, Provenance: prov, UpdatedAt: stamp, Status: "OFFICIAL"}
	at := func(i int) float64 {
		if i >= len(vals) {
			return vals[len(vals)-1]
		}
		return vals[i]
	}
	if len(vals) > 1 {
		m.Previous = vals[1]
	}
	m.Change5D = m.Value - at(5)
	m.Change20D = m.Value - at(20)
	m.Change1M = m.Value - at(22)
	m.Change3M = m.Value - at(66)
	return m, true
}
func (e *Engine) refreshTreasury(ctx context.Context) {
	yr := time.Now().Year()
	raw, err := fetchText(ctx, fmt.Sprintf("%s?data=daily_treasury_yield_curve&field_tdr_date_value=%d", treasuryXMLBaseURL, yr))
	if err != nil {
		e.setHealth("treasury", "degraded · cached/public fallback preserved")
		e.reconcileMacroRatesHealth()
		return
	}
	extract := func(tag string) []float64 {
		rx := regexp.MustCompile(`(?i)<d:` + regexp.QuoteMeta(tag) + `[^>]*>([^<]+)</d:` + regexp.QuoteMeta(tag) + `>`)
		ms := rx.FindAllStringSubmatch(raw, -1)
		vals := []float64{}
		for i := len(ms) - 1; i >= 0; i-- {
			v, er := strconv.ParseFloat(strings.TrimSpace(ms[i][1]), 64)
			if er == nil {
				vals = append(vals, v)
			}
		}
		return vals
	}
	defs := map[string]struct{ label, tag string }{"UST2Y": {"U.S. Treasury 2Y", "BC_2YEAR"}, "UST10Y": {"U.S. Treasury 10Y", "BC_10YEAR"}, "UST30Y": {"U.S. Treasury 30Y", "BC_30YEAR"}}
	out := map[string]MacroMetric{}
	now := time.Now().UnixMilli()
	for k, d := range defs {
		if m, ok := metricWithChanges(k, d.label, "%", "U.S. Treasury", "OFFICIAL", extract(d.tag), now); ok {
			out[k] = m
		}
	}
	// Real yield endpoint uses separate data key.
	if rr, er := fetchText(ctx, fmt.Sprintf("%s?data=daily_treasury_real_yield_curve&field_tdr_date_value=%d", treasuryXMLBaseURL, yr)); er == nil {
		rx := regexp.MustCompile(`(?i)<d:BC_10YEAR[^>]*>([^<]+)</d:BC_10YEAR>`)
		ms := rx.FindAllStringSubmatch(rr, -1)
		vals := []float64{}
		for i := len(ms) - 1; i >= 0; i-- {
			v, er := strconv.ParseFloat(strings.TrimSpace(ms[i][1]), 64)
			if er == nil {
				vals = append(vals, v)
			}
		}
		if m, ok := metricWithChanges("UST10Y_REAL", "10Y Real Yield", "%", "U.S. Treasury", "OFFICIAL", vals, now); ok {
			out[m.Key] = m
		}
	}
	if len(out) == 0 {
		e.setHealth("treasury", "degraded · no current official observations")
		e.reconcileMacroRatesHealth()
		return
	}
	e.mu.Lock()
	for k, v := range out {
		e.macroMetrics[k] = v
	}
	e.lastUpdated["treasury"] = now
	e.health["treasury"] = "healthy · official Treasury yields"
	e.mu.Unlock()
	e.reconcileMacroRatesHealth()
}

type blsResponse struct {
	Status  string   `json:"status"`
	Message []string `json:"message"`
	Results struct {
		Series []struct {
			SeriesID string `json:"seriesID"`
			Data     []struct {
				Year   string `json:"year"`
				Period string `json:"period"`
				Value  string `json:"value"`
			} `json:"data"`
		} `json:"series"`
	} `json:"Results"`
}

func blsSeries(ctx context.Context, series string) ([]float64, error) {
	var p blsResponse
	raw := blsAPIBaseURL + "/" + url.PathEscape(series)
	if err := getJSON(ctx, &http.Client{Timeout: 12 * time.Second}, raw, nil, &p); err != nil {
		return nil, err
	}
	if strings.ToUpper(p.Status) != "REQUEST_SUCCEEDED" {
		return nil, fmt.Errorf("BLS %s", strings.Join(p.Message, "; "))
	}
	vals := []float64{}
	if len(p.Results.Series) == 0 {
		return vals, nil
	}
	for _, d := range p.Results.Series[0].Data {
		if strings.HasPrefix(d.Period, "M") && d.Period != "M13" {
			v, er := strconv.ParseFloat(d.Value, 64)
			if er == nil {
				vals = append(vals, v)
			}
		}
	}
	return vals, nil
}
func (e *Engine) refreshBLSActuals(ctx context.Context) {
	defs := map[string]struct{ label, series, unit string }{
		"CPI_INDEX": {"CPI All Items", "CUUR0000SA0", "index"}, "CORE_CPI_INDEX": {"Core CPI", "CUUR0000SA0L1E", "index"},
		"CPI_SHELTER": {"CPI Shelter", "CUUR0000SAH1", "index"}, "CPI_SERVICES_X_ENERGY": {"CPI Services less Energy", "CUUR0000SASLE", "index"},
		"CPI_GOODS_X_FE": {"CPI Commodities less Food & Energy", "CUUR0000SACL1E", "index"}, "CPI_ENERGY": {"CPI Energy", "CUUR0000SA0E", "index"},
		"NONFARM": {"Total Nonfarm Payrolls", "CES0000000001", "thousands"}, "PRIVATE_PAYROLL": {"Total Private Payrolls", "CES0500000001", "thousands"},
		"UNEMP": {"Unemployment Rate", "LNS14000000", "%"}, "PARTICIPATION": {"Labor Force Participation", "LNS11300000", "%"},
		"AHE": {"Average Hourly Earnings", "CES0500000003", "USD/hour"}, "PPI_FINAL": {"PPI Final Demand", "WPUFD4", "index"},
		"PAYROLL_MANUFACTURING": {"Manufacturing Payrolls", "CES3000000001", "thousands"},
		"PAYROLL_CONSTRUCTION":  {"Construction Payrolls", "CES2000000001", "thousands"},
		"PAYROLL_FINANCIAL":     {"Financial Activities Payrolls", "CES5500000001", "thousands"},
		"PAYROLL_PRO_BUSINESS":  {"Professional & Business Services Payrolls", "CES6000000001", "thousands"},
	}
	now := time.Now().UnixMilli()
	out := map[string]MacroMetric{}
	for k, d := range defs {
		vals, err := blsSeries(ctx, d.series)
		if err != nil || len(vals) == 0 {
			continue
		}
		m := MacroMetric{Key: k, Label: d.label, Value: vals[0], Unit: d.unit, Source: "BLS", Provenance: "OFFICIAL", UpdatedAt: now, Status: "OFFICIAL"}
		if len(vals) > 1 {
			m.Previous = vals[1]
		}
		if strings.HasPrefix(k, "CPI_") || k == "CPI_INDEX" || k == "CORE_CPI_INDEX" || k == "PPI_FINAL" {
			if len(vals) > 12 && vals[12] != 0 {
				m.Value = (vals[0]/vals[12] - 1) * 100
				m.Unit = "% YoY"
				if len(vals) > 13 && vals[13] != 0 {
					m.Previous = (vals[1]/vals[13] - 1) * 100
				}
			}
		} else if (k == "NONFARM" || k == "PRIVATE_PAYROLL" || strings.HasPrefix(k, "PAYROLL_")) && len(vals) > 1 {
			m.Value = vals[0] - vals[1]
			m.Unit = "thousands change"
		}
		out[k] = m
	}
	if len(out) == 0 {
		e.setHealth("bls-actuals", "degraded · official BLS actuals unavailable")
		return
	}
	e.mu.Lock()
	for k, v := range out {
		e.macroMetrics[k] = v
	}
	e.lastUpdated["bls-actuals"] = now
	e.health["bls-actuals"] = "healthy · official BLS actuals"
	e.mu.Unlock()
}

type eiaLegacyResponse struct {
	Response struct {
		Data []map[string]any `json:"data"`
	} `json:"response"`
}

func eiaV2Rows(ctx context.Context, raw string) ([]map[string]any, error) {
	var p eiaLegacyResponse
	if err := getJSON(ctx, &http.Client{Timeout: 12 * time.Second}, raw, nil, &p); err != nil {
		return nil, err
	}
	return p.Response.Data, nil
}

func eiaValue(row map[string]any) (float64, bool) {
	for _, field := range []string{"value", "Value", "price"} {
		if x, ok := row[field]; ok {
			v, err := strconv.ParseFloat(fmt.Sprint(x), 64)
			if err == nil {
				return v, true
			}
		}
	}
	return 0, false
}
func (e *Engine) refreshEIAActuals(ctx context.Context, key string) {
	if strings.TrimSpace(key) == "" {
		e.setHealth("eia-actuals", "not configured · EIA key optional")
		return
	}
	out := map[string]MacroMetric{}
	now := time.Now().UnixMilli()
	defs := map[string]struct{ label, series, unit string }{
		"WTI_OFFICIAL": {"WTI Spot", "PET.RWTC.D", "USD/bbl"}, "BRENT_OFFICIAL": {"Brent Spot", "PET.RBRTE.D", "USD/bbl"},
		"CRUDE_STOCKS":      {"U.S. Crude Stocks ex SPR", "PET.WCESTUS1.W", "million barrels"},
		"GASOLINE_STOCKS":   {"U.S. Gasoline Stocks", "PET.WGTSTUS1.W", "million barrels"},
		"DISTILLATE_STOCKS": {"U.S. Distillate Stocks", "PET.WDISTUS1.W", "million barrels"},
		"CRUDE_PRODUCTION":  {"U.S. Crude Production", "PET.WCRFPUS2.W", "thousand b/d"},
		"CRUDE_IMPORTS":     {"U.S. Crude Imports", "PET.WCRIMUS2.W", "thousand b/d"},
		"REFINERY_UTIL":     {"U.S. Refinery Utilization", "PET.WPULEUS3.W", "%"},
		"NATGAS_STORAGE":    {"Lower 48 Working Gas Storage", "NG.NW2_EPG0_SWO_R48_BCF.W", "Bcf"},
	}
	for k, d := range defs {
		raw := eiaAPIBaseURL + "/" + url.PathEscape(d.series) + "?api_key=" + url.QueryEscape(key) + "&length=25"
		var p eiaLegacyResponse
		if er := getJSON(ctx, &http.Client{Timeout: 12 * time.Second}, raw, nil, &p); er != nil {
			continue
		}
		vals := []float64{}
		for _, row := range p.Response.Data {
			for _, field := range []string{"value", "Value", "price"} {
				if x, ok := row[field]; ok {
					v, er := strconv.ParseFloat(fmt.Sprint(x), 64)
					if er == nil {
						vals = append(vals, v)
						break
					}
				}
			}
		}
		if m, ok := metricWithChanges(k, d.label, d.unit, "EIA", "OFFICIAL", vals, now); ok {
			out[k] = m
		}
	}
	if strings.Contains(eiaAPIBaseURL, "api.eia.gov") {
		// P2 thematic context: hourly electric-system demand and STEO forecasts are
		// best-effort official context. They never drive deterministic setup scores.
		if rows, err := eiaV2Rows(ctx, "https://api.eia.gov/v2/electricity/rto/region-data/data/?api_key="+url.QueryEscape(key)+"&frequency=hourly&data[0]=value&facets[type][]=D&sort[0][column]=period&sort[0][direction]=desc&length=25"); err == nil && len(rows) > 0 {
			latestPeriod := fmt.Sprint(rows[0]["period"])
			total := 0.0
			count := 0
			for _, row := range rows {
				if fmt.Sprint(row["period"]) != latestPeriod {
					continue
				}
				if v, ok := eiaValue(row); ok {
					total += v
					count++
				}
			}
			if count > 0 {
				out["ELECTRICITY_DEMAND_SAMPLE"] = MacroMetric{Key: "ELECTRICITY_DEMAND_SAMPLE", Label: "Hourly Electric Demand · BA sample", Value: total, Unit: "MW sample", Source: "EIA", Provenance: "OFFICIAL", UpdatedAt: now, Status: "OFFICIAL", Period: latestPeriod}
			}
		}
		for keyName, seriesID := range map[string]string{"STEO_WTI_FORECAST": "WTIPUUS", "STEO_BRENT_FORECAST": "BREPUUS", "STEO_HENRY_HUB_FORECAST": "NGHHUUS"} {
			raw := "https://api.eia.gov/v2/steo/data/?api_key=" + url.QueryEscape(key) + "&frequency=monthly&data[0]=value&facets[seriesId][]=" + url.QueryEscape(seriesID) + "&sort[0][column]=period&sort[0][direction]=desc&length=2"
			if rows, err := eiaV2Rows(ctx, raw); err == nil && len(rows) > 0 {
				if v, ok := eiaValue(rows[0]); ok {
					label := map[string]string{"STEO_WTI_FORECAST": "STEO WTI Forecast", "STEO_BRENT_FORECAST": "STEO Brent Forecast", "STEO_HENRY_HUB_FORECAST": "STEO Henry Hub Forecast"}[keyName]
					out[keyName] = MacroMetric{Key: keyName, Label: label, Value: v, Unit: "forecast", Source: "EIA STEO", Provenance: "OFFICIAL FORECAST", UpdatedAt: now, Status: "OFFICIAL", Period: fmt.Sprint(rows[0]["period"])}
				}
			}
		}
	}

	if len(out) == 0 {
		e.setHealth("eia-actuals", "degraded · official EIA actuals unavailable")
		return
	}
	e.mu.Lock()
	for k, v := range out {
		e.macroMetrics[k] = v
	}
	e.lastUpdated["eia-actuals"] = now
	e.health["eia-actuals"] = "healthy · official EIA energy"
	e.mu.Unlock()
}

var percentRx = regexp.MustCompile(`(?i)([-+]?\d+(?:\.\d+)?)\s*percent`)

func beaMetricFromPage(ctx context.Context, path, key, label string) (MacroMetric, bool) {
	raw, err := fetchText(ctx, beaNewsBaseURL+path)
	if err != nil {
		return MacroMetric{}, false
	}
	text := visibleText(raw)
	m := percentRx.FindStringSubmatch(text)
	if len(m) < 2 {
		return MacroMetric{}, false
	}
	v, err := strconv.ParseFloat(m[1], 64)
	if err != nil {
		return MacroMetric{}, false
	}
	return MacroMetric{Key: key, Label: label, Value: v, Unit: "%", Source: "BEA", Provenance: "OFFICIAL PUBLIC RELEASE", UpdatedAt: time.Now().UnixMilli(), Status: "OFFICIAL", Period: "latest release"}, true
}
func (e *Engine) refreshBEAActuals(ctx context.Context) {
	out := map[string]MacroMetric{}
	if m, ok := beaMetricFromPage(ctx, "/glance/gdp", "BEA_GDP", "Real GDP"); ok {
		out[m.Key] = m
	}
	if m, ok := beaMetricFromPage(ctx, "/glance/pce", "BEA_PCE", "PCE Price Index"); ok {
		out[m.Key] = m
	}
	if len(out) == 0 {
		e.setHealth("bea-actuals", "degraded · public BEA actuals unavailable")
		return
	}
	e.mu.Lock()
	for k, v := range out {
		e.macroMetrics[k] = v
	}
	e.lastUpdated["bea-actuals"] = time.Now().UnixMilli()
	e.health["bea-actuals"] = "healthy · official BEA releases"
	e.mu.Unlock()
}

func (e *Engine) refreshOfficialMacroActuals(ctx context.Context, eiaKey string) {
	e.refreshTreasury(ctx)
	e.refreshBLSActuals(ctx)
	e.refreshBEAActuals(ctx)
	e.refreshEIAActuals(ctx, eiaKey)
}

func uniqueLabels(items []string) []string {
	seen := map[string]bool{}
	out := []string{}
	for _, v := range items {
		k := strings.ToLower(strings.TrimSpace(v))
		if k == "" || seen[k] {
			continue
		}
		seen[k] = true
		out = append(out, strings.TrimSpace(v))
	}
	return out
}
func affectedContext(ev MacroEvent, tracked []string) ([]string, []string) {
	n := strings.ToLower(ev.Name)
	sectors := []string{"Broad Market"}
	syms := []string{"SPY", "QQQ", "IWM", "VIX", "TLT", "UUP"}
	if strings.Contains(n, "oil") || strings.Contains(n, "energy") {
		sectors = append(sectors, "Energy")
		syms = append(syms, "XLE", "USO")
	}
	if strings.Contains(n, "employment") || strings.Contains(n, "payroll") || strings.Contains(n, "cpi") || strings.Contains(n, "consumer price") || strings.Contains(n, "producer price") || strings.Contains(n, "pce") || strings.Contains(n, "fomc") || strings.Contains(n, "fed") || strings.Contains(n, "inflation") {
		sectors = append(sectors, "Technology", "Financials", "Consumer")
		syms = append(syms, "XLK", "XLF", "XLY")
	}
	if ev.Region == "CN" {
		sectors = append(sectors, "Semiconductors", "Materials")
		// Global events propagate only into U.S.-market context/sector instruments;
		// do not promote a foreign-market proxy into the actionable ticker path.
		syms = append(syms, "SMH", "XLB")
	}
	if ev.Region == "JP" || ev.Region == "CN" {
		syms = append(syms, "SMH")
	}
	// tracked tickers are included for queue pre-context; no market direction is inferred from membership.
	for _, s := range tracked {
		if len(syms) >= 24 {
			break
		}
		syms = append(syms, s)
	}
	return uniqueSymbols(syms), uniqueLabels(sectors)
}
func updateEventLifecycles(events []MacroEvent, now time.Time) []MacroEvent {
	out := clone(events)
	for i := range out {
		if out[i].TimeKnown && out[i].StartsAt > 0 {
			out[i].Lifecycle = eventLifecycle(time.UnixMilli(out[i].StartsAt), now)
		}
	}
	return out
}

func (e *Engine) highImpactModeActive() bool {
	e.app.mu.RLock()
	enabled := e.app.state.Settings.MacroEventModeEnabled
	e.app.mu.RUnlock()
	e.mu.RLock()
	events := clone(e.macroEvents)
	e.mu.RUnlock()
	return eventModeFor(events, time.Now(), enabled).Active
}

// Conservative guidance extraction. It only accepts explicit ranges next to guidance/outlook language.
var guidanceAnchorRx = regexp.MustCompile(`(?i)guidance|outlook|expects?|forecast`)
var revenueRangeRx = regexp.MustCompile(`(?i)(?:revenue|sales)[^\n]{0,100}?\$?([0-9]+(?:\.[0-9]+)?)\s*(billion|million|b|m)?\s*(?:-|–|to)\s*\$?([0-9]+(?:\.[0-9]+)?)\s*(billion|million|b|m)?`)
var epsRangeRx = regexp.MustCompile(`(?i)(?:eps|earnings per share)[^\n]{0,100}?\$?([0-9]+(?:\.[0-9]+)?)\s*(?:-|–|to)\s*\$?([0-9]+(?:\.[0-9]+)?)`)

func scaleGuidance(v float64, u string) float64 {
	u = strings.ToLower(u)
	if u == "b" || strings.HasPrefix(u, "billion") {
		return v * 1e9
	}
	if u == "m" || strings.HasPrefix(u, "million") {
		return v * 1e6
	}
	return v
}
func parseGuidanceRanges(text string) (revLow, revHigh, epsLow, epsHigh float64, okRev, okEPS bool) {
	loc := guidanceAnchorRx.FindStringIndex(text)
	if loc == nil {
		return
	}
	end := loc[0] + 320
	if end > len(text) {
		end = len(text)
	}
	window := text[loc[0]:end]
	if m := revenueRangeRx.FindStringSubmatch(window); len(m) > 0 {
		a, _ := strconv.ParseFloat(m[1], 64)
		b, _ := strconv.ParseFloat(m[3], 64)
		a = scaleGuidance(a, m[2])
		b = scaleGuidance(b, m[4])
		if a > 0 && b >= a {
			revLow, revHigh, okRev = a, b, true
		}
	}
	if m := epsRangeRx.FindStringSubmatch(window); len(m) > 0 {
		a, _ := strconv.ParseFloat(m[1], 64)
		b, _ := strconv.ParseFloat(m[2], 64)
		if b >= a {
			epsLow, epsHigh, okEPS = a, b, true
		}
	}
	return
}

func (e *Engine) enrichEarningsGuidanceFromEvidence() {
	e.mu.Lock()
	defer e.mu.Unlock()
	bySym := map[string][]string{}
	for _, n := range e.news {
		if !strings.Contains(strings.ToLower(n.Headline+" "+n.Summary), "guidance") && !strings.Contains(strings.ToLower(n.Headline+" "+n.Summary), "outlook") {
			continue
		}
		for _, s := range n.Symbols {
			bySym[normalizeSymbol(s)] = append(bySym[normalizeSymbol(s)], n.Headline+". "+n.Summary)
		}
	}
	for _, f := range e.filings {
		if strings.EqualFold(f.Form, "8-K") {
			bySym[normalizeSymbol(f.Symbol)] = append(bySym[normalizeSymbol(f.Symbol)], f.Description)
		}
	}
	prev := map[string]EarningsItem{}
	sort.Slice(e.earnings, func(i, j int) bool { return e.earnings[i].Date < e.earnings[j].Date })
	for i := range e.earnings {
		x := &e.earnings[i]
		sym := normalizeSymbol(x.Symbol)
		for _, txt := range bySym[sym] {
			rl, rh, el, eh, okr, oke := parseGuidanceRanges(txt)
			if okr {
				x.GuidanceRevenueNewLow = &rl
				x.GuidanceRevenueNewHigh = &rh
				x.GuidanceSource = "Sourced company/SEC/news evidence"
			}
			if oke {
				x.GuidanceEPSNewLow = &el
				x.GuidanceEPSNewHigh = &eh
				x.GuidanceSource = "Sourced company/SEC/news evidence"
			}
			if okr || oke {
				if p, ok := prev[sym]; ok {
					x.GuidanceRevenuePrevLow = p.GuidanceRevenueNewLow
					x.GuidanceRevenuePrevHigh = p.GuidanceRevenueNewHigh
					x.GuidanceEPSPrevLow = p.GuidanceEPSNewLow
					x.GuidanceEPSPrevHigh = p.GuidanceEPSNewHigh
				}
				break
			}
		}
		if x.GuidanceRevenueNewLow != nil || x.GuidanceEPSNewLow != nil {
			prev[sym] = *x
		}
	}
}

// Official/public international close adapters. These are intentionally separate
// from live/direct providers so a completed local-market session remains truthful
// context rather than being mislabeled stale. The adapter only stores values it
// can parse from the official exchange response; failures preserve lower-priority
// real provider/proxy/cache routing.
var twsePublicCloseURL = "https://www.twse.com.tw/exchangeReport/MI_INDEX?response=html&type=ALLBUT0999"
var twseIndexRowRx = regexp.MustCompile(`發行量加權股價指數\s+([0-9,]+(?:\.[0-9]+)?)\s+([+-])?\s*([0-9,]+(?:\.[0-9]+)?)\s+([+-]?[0-9]+(?:\.[0-9]+)?)`)
var twseROCDateRx = regexp.MustCompile(`(\d{3})年(\d{2})月(\d{2})日`)

func twseOfficialClose(ctx context.Context) (GlobalDriver, error) {
	raw, err := fetchText(ctx, twsePublicCloseURL)
	if err != nil {
		return GlobalDriver{}, err
	}
	text := visibleText(raw)
	m := twseIndexRowRx.FindStringSubmatch(text)
	if len(m) < 5 {
		return GlobalDriver{}, fmt.Errorf("TWSE TAIEX row unavailable")
	}
	px, err := strconv.ParseFloat(strings.ReplaceAll(m[1], ",", ""), 64)
	if err != nil || px <= 0 {
		return GlobalDriver{}, fmt.Errorf("TWSE invalid close")
	}
	pct, err := strconv.ParseFloat(m[4], 64)
	if err != nil {
		return GlobalDriver{}, fmt.Errorf("TWSE invalid percent change")
	}
	stamp := time.Now().UnixMilli()
	dateLabel := "latest completed session"
	if dm := twseROCDateRx.FindStringSubmatch(text); len(dm) == 4 {
		y, _ := strconv.Atoi(dm[1])
		mo, _ := strconv.Atoi(dm[2])
		day, _ := strconv.Atoi(dm[3])
		y += 1911
		if loc, er := time.LoadLocation("Asia/Taipei"); er == nil {
			t := time.Date(y, time.Month(mo), day, 13, 30, 0, 0, loc)
			stamp = t.UnixMilli()
			dateLabel = t.Format("2006-01-02")
		}
	}
	return GlobalDriver{Key: "taiwan_official_close", Label: "TAIEX / Taiwan · official close", State: driverState(pct, false), Value: px, ChangePercent: pct, Source: "Taiwan Stock Exchange", Provenance: "OFFICIAL CLOSE", UpdatedAt: stamp, Confidence: 95, Detail: "TAIEX official close · " + dateLabel, Session: "OFFICIAL CLOSE", Underlying: "TAIEX", ProviderSymbol: "TAIEX", IsProxy: false}, nil
}

func (e *Engine) refreshOfficialGlobalCloses(ctx context.Context) {
	d, err := twseOfficialClose(ctx)
	if err != nil {
		e.setHealth("global-public", "degraded · official/public close unavailable; real fallback preserved")
		return
	}
	e.mu.Lock()
	// Direct/live providers are allowed to replace this same canonical key later;
	// otherwise the official completed session sits ahead of the ETF proxy.
	e.globalDirect[d.Key] = d
	e.lastUpdated["global-public"] = time.Now().UnixMilli()
	e.health["global-public"] = "healthy · TWSE official close"
	e.mu.Unlock()
}
