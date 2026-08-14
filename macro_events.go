package main

import (
	"context"
	"fmt"
	"html"
	"io"
	"math"
	"net/http"
	"net/url"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"
)

var fredAPIBaseURL = "https://api.stlouisfed.org"

// reconcileMacroRatesHealth publishes aggregate Macro Rates health from the
// canonical official sources. U.S. Treasury is sufficient for core rates; FRED
// enriches the state with credit, financial conditions and USD. This prevents a
// temporary/optional FRED issue from falsely degrading usable official rates.
func (e *Engine) reconcileMacroRatesHealth() {
	e.mu.Lock()
	defer e.mu.Unlock()
	metricOK := func(keys ...string) bool {
		for _, k := range keys {
			if m, ok := e.macroMetrics[k]; ok && m.Status != "UNAVAILABLE" && m.Value != 0 {
				return true
			}
		}
		return false
	}
	treasuryCore := metricOK("UST10Y", "UST2Y", "UST30Y")
	fredCore := metricOK("DGS10", "DGS2", "DGS30")
	treasuryHealth := strings.ToLower(strings.TrimSpace(e.health["treasury"]))
	fredHealth := strings.ToLower(strings.TrimSpace(e.health["fred-rates"]))
	treasuryCurrent := treasuryCore && strings.Contains(treasuryHealth, "healthy")
	fredCurrent := fredCore && strings.Contains(fredHealth, "healthy")

	stamp := e.lastUpdated["treasury"]
	if e.lastUpdated["fred-rates"] > stamp {
		stamp = e.lastUpdated["fred-rates"]
	}
	if stamp > 0 && (treasuryCore || fredCore) {
		e.lastUpdated["macro-rates"] = stamp
	}

	switch {
	case treasuryCurrent && fredCurrent:
		e.health["macro-rates"] = "healthy · official Treasury core + FRED rates/credit enrichment"
	case treasuryCurrent:
		if strings.Contains(fredHealth, "not configured") || fredHealth == "" {
			e.health["macro-rates"] = "healthy · official Treasury core rates · FRED enrichment optional"
		} else {
			e.health["macro-rates"] = "healthy · official Treasury core rates · FRED enrichment temporarily unavailable"
		}
	case fredCurrent:
		e.health["macro-rates"] = "healthy · official FRED rates/credit · Treasury fallback temporarily unavailable"
	case treasuryCore || fredCore:
		e.health["macro-rates"] = "degraded · cached official rates only"
	default:
		e.health["macro-rates"] = "temporarily unavailable · no usable official rates yet"
	}
}

// FRED official API latest observation refresh. Missing credentials leave the
// official Treasury core-rates path active; FRED remains optional enrichment.
func (e *Engine) refreshFRED(ctx context.Context, key string) {
	if strings.TrimSpace(key) == "" {
		e.setHealth("fred-rates", "not configured · FRED key optional")
		e.reconcileMacroRatesHealth()
		return
	}
	series := map[string]struct{ label, unit string }{
		"DGS2": {"U.S. 2Y", "%"}, "DGS10": {"U.S. 10Y", "%"}, "DGS30": {"U.S. 30Y", "%"},
		"DFII10": {"10Y Real Yield", "%"}, "DFF": {"Effective Fed Funds", "%"}, "SOFR": {"SOFR", "%"},
		"T10Y2Y": {"10Y−2Y Treasury Spread", "pp"}, "BAMLH0A0HYM2": {"High Yield Spread", "%"},
		"BAMLC0A0CM": {"Investment Grade Spread", "%"}, "NFCI": {"Chicago Fed Financial Conditions", "index"},
		"STLFSI4": {"St. Louis Fed Financial Stress", "index"}, "DTWEXBGS": {"Broad Trade-Weighted U.S. Dollar", "index"},
		"WALCL": {"Federal Reserve Total Assets", "millions USD"},
	}
	out := map[string]MacroMetric{}
	client := &http.Client{Timeout: 12 * time.Second}
	for id, def := range series {
		raw := fredAPIBaseURL + "/fred/series/observations?series_id=" + url.QueryEscape(id) + "&api_key=" + url.QueryEscape(key) + "&file_type=json&sort_order=desc&limit=90"
		var payload struct {
			Observations []struct {
				Date  string `json:"date"`
				Value string `json:"value"`
			} `json:"observations"`
		}
		if err := getJSON(ctx, client, raw, nil, &payload); err != nil {
			continue
		}
		vals := []float64{}
		stamp := int64(0)
		for _, ob := range payload.Observations {
			if ob.Value == "." {
				continue
			}
			v, err := strconv.ParseFloat(ob.Value, 64)
			if err != nil {
				continue
			}
			vals = append(vals, v)
			if stamp == 0 {
				if t, er := time.Parse("2006-01-02", ob.Date); er == nil {
					stamp = t.UnixMilli()
				}
			}
		}
		if m, ok := metricWithChanges(id, def.label, def.unit, "FRED", "OFFICIAL", vals, stamp); ok {
			out[id] = m
		}
	}
	if len(out) == 0 {
		e.setHealth("fred-rates", "temporarily unavailable · keeping cached FRED enrichment")
		e.reconcileMacroRatesHealth()
		return
	}
	now := time.Now().UnixMilli()
	e.mu.Lock()
	for k, v := range out {
		e.macroMetrics[k] = v
	}
	e.lastUpdated["fred-rates"] = now
	if _, ten := out["DGS10"]; ten && len(out) >= 6 {
		e.health["fred-rates"] = fmt.Sprintf("healthy · FRED official · %d/%d series", len(out), len(series))
	} else {
		e.health["fred-rates"] = fmt.Sprintf("degraded · partial FRED enrichment · %d/%d series", len(out), len(series))
	}
	e.mu.Unlock()
	e.reconcileMacroRatesHealth()
}
func fetchText(ctx context.Context, raw string) (string, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, raw, nil)
	if err != nil {
		return "", err
	}
	req.Header.Set("User-Agent", "DE.PULSE/14.0 market-intelligence terminal")
	resp, err := (&http.Client{Timeout: 15 * time.Second}).Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		b, _ := io.ReadAll(io.LimitReader(resp.Body, 300))
		return "", fmt.Errorf("HTTP %d: %s", resp.StatusCode, strings.TrimSpace(string(b)))
	}
	b, err := io.ReadAll(io.LimitReader(resp.Body, 3<<20))
	return string(b), err
}

var tagRx = regexp.MustCompile(`(?s)<[^>]+>`)
var wsRx = regexp.MustCompile(`\s+`)

func visibleText(raw string) string {
	return strings.TrimSpace(wsRx.ReplaceAllString(html.UnescapeString(tagRx.ReplaceAllString(raw, " ")), " "))
}

var monthRx = regexp.MustCompile(`(?i)\b(January|February|March|April|May|June|July|August|September|October|November|December)\s+(\d{1,2})(?:\s*[-–]\s*(\d{1,2}))?,?\s+(20\d{2})\b`)
var isoDateRx = regexp.MustCompile(`\b(20\d{2})[-/](\d{1,2})[-/](\d{1,2})\b`)
var ampmTimeRx = regexp.MustCompile(`(?i)\b(1[0-2]|0?[1-9]):([0-5]\d)\s*(a\.?m\.?|p\.?m\.?)\s*(ET|EST|EDT|CET|CEST|JST|CST)?\b`)
var explicit24TimeRx = regexp.MustCompile(`(?i)\b([01]?\d|2[0-3]):([0-5]\d)\s*(ET|EST|EDT|CET|CEST|JST|China Standard Time)\b`)

func officialDateTime(window, region string, date time.Time) (time.Time, bool) {
	locName := map[string]string{"US": "America/New_York", "EU": "Europe/Berlin", "JP": "Asia/Tokyo", "CN": "Asia/Shanghai"}[region]
	if locName == "" {
		return time.Time{}, false
	}
	loc, err := time.LoadLocation(locName)
	if err != nil {
		return time.Time{}, false
	}
	if m := ampmTimeRx.FindStringSubmatch(window); len(m) > 0 {
		h, _ := strconv.Atoi(m[1])
		min, _ := strconv.Atoi(m[2])
		ap := strings.ToLower(strings.ReplaceAll(m[3], ".", ""))
		if ap == "pm" && h < 12 {
			h += 12
		}
		if ap == "am" && h == 12 {
			h = 0
		}
		return time.Date(date.Year(), date.Month(), date.Day(), h, min, 0, 0, loc), true
	}
	if m := explicit24TimeRx.FindStringSubmatch(window); len(m) > 0 {
		h, _ := strconv.Atoi(m[1])
		min, _ := strconv.Atoi(m[2])
		return time.Date(date.Year(), date.Month(), date.Day(), h, min, 0, 0, loc), true
	}
	return time.Time{}, false
}

func parseOfficialEvents(raw, region, source, sourceURL string, names []string) []MacroEvent {
	text := visibleText(raw)
	lower := strings.ToLower(text)
	now := time.Now()
	events := []MacroEvent{}
	for _, name := range names {
		needle := strings.ToLower(name)
		pos := 0
		for {
			i := strings.Index(lower[pos:], needle)
			if i < 0 {
				break
			}
			i += pos
			start := maxInt(0, i-180)
			end := minInt2(len(text), i+len(name)+220)
			window := text[start:end]
			dt, ok := dateFromWindow(window, now)
			if ok && dt.After(now.AddDate(0, -2, 0)) && dt.Before(now.AddDate(1, 6, 0)) {
				id := strings.ToLower(strings.ReplaceAll(region+"-"+name+"-"+dt.Format("20060102"), " ", "-"))
				impact := "HIGH"
				if strings.Contains(strings.ToLower(name), "minutes") || strings.Contains(strings.ToLower(name), "claims") {
					impact = "MEDIUM"
				}
				eventTime, timeKnown := officialDateTime(window, region, dt)
				startsAt := int64(0)
				lifecycle := "UPCOMING"
				if timeKnown {
					startsAt = eventTime.UnixMilli()
					lifecycle = eventLifecycle(eventTime, now)
				} else if dt.Before(time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, dt.Location())) {
					lifecycle = "RESOLVED"
				}
				ev := MacroEvent{ID: id, Region: region, Name: name, Impact: impact, Lifecycle: lifecycle, StartsAt: startsAt, Date: dt.Format("2006-01-02"), TimeKnown: timeKnown, Source: source, SourceURL: sourceURL, UpdatedAt: now.UnixMilli()}
				ev.ProcessingClass = macroEventProcessingClass(ev)
				events = append(events, ev)
			}
			pos = i + len(needle)
			if pos >= len(lower) {
				break
			}
		}
	}

	seen := map[string]bool{}
	out := events[:0]
	for _, ev := range events {
		if seen[ev.ID] {
			continue
		}
		seen[ev.ID] = true
		out = append(out, ev)
	}
	sort.Slice(out, func(i, j int) bool { return out[i].Date < out[j].Date })
	return out
}
func maxInt(a, b int) int {
	if a > b {
		return a
	}
	return b
}
func minInt2(a, b int) int {
	if a < b {
		return a
	}
	return b
}
func dateFromWindow(s string, now time.Time) (time.Time, bool) {
	if m := monthRx.FindStringSubmatch(s); len(m) > 0 {
		month, _ := time.Parse("January", m[1])
		day, _ := strconv.Atoi(m[2])
		year, _ := strconv.Atoi(m[4])
		return time.Date(year, month.Month(), day, 0, 0, 0, 0, time.Local), true
	}
	if m := isoDateRx.FindStringSubmatch(s); len(m) > 0 {
		y, _ := strconv.Atoi(m[1])
		mo, _ := strconv.Atoi(m[2])
		d, _ := strconv.Atoi(m[3])
		return time.Date(y, time.Month(mo), d, 0, 0, 0, 0, time.Local), true
	}
	return time.Time{}, false
}
func eventLifecycle(t, now time.Time) string {
	if t.After(now) {
		return "UPCOMING"
	}
	age := now.Sub(t)
	if age < 30*time.Second {
		return "RELEASED"
	}
	if age <= time.Hour {
		return "MARKET REACTION"
	}
	return "RESOLVED"
}

func (e *Engine) refreshMacroEvents(ctx context.Context) {
	sources := []struct {
		region, source, url string
		names               []string
	}{
		{"US", "Federal Reserve", "https://www.federalreserve.gov/monetarypolicy/fomccalendars.htm", []string{"FOMC Meeting", "Federal Open Market Committee", "Minutes"}},
		{"US", "BLS", "https://www.bls.gov/schedule/news_release/", []string{"Consumer Price Index", "Employment Situation", "Producer Price Index"}},
		{"US", "BEA", "https://www.bea.gov/news/schedule", []string{"Gross Domestic Product", "Personal Income and Outlays"}},
		{"US", "Census", "https://www.census.gov/economic-indicators/calendar-listview.html", []string{"Advance Monthly Sales for Retail and Food Services", "Retail Sales"}},
		{"US", "Department of Labor", "https://www.dol.gov/ui/data.pdf", []string{"Unemployment Insurance Weekly Claims", "Initial Claims"}},
		{"US", "ISM", "https://www.ismworld.org/supply-management-news-and-reports/reports/ism-report-on-business/", []string{"Manufacturing PMI", "Services PMI", "ISM Manufacturing", "ISM Services"}},
		// Global Context is deliberately selective: only major central-bank/growth
		// events with credible U.S.-market transmission remain in the active calendar.
		{"EU", "ECB", "https://www.ecb.europa.eu/press/calendars/mgcgc/html/index.en.html", []string{"Monetary policy meeting", "Governing Council"}},
		{"JP", "BOJ", "https://www.boj.or.jp/en/mopo/mpmsche_minu/index.htm", []string{"Monetary Policy Meeting", "Outlook for Economic Activity and Prices"}},
		{"CN", "China NBS", "https://www.stats.gov.cn/english/PressRelease/", []string{"Purchasing Managers Index", "Gross Domestic Product", "National Economy"}},
		{"CN", "PBOC", "http://www.pbc.gov.cn/en/3688006/index.html", []string{"Loan Prime Rate", "Reserve Requirement"}},
	}
	combined := []MacroEvent{}
	success := 0
	for _, src := range sources {
		raw, err := fetchText(ctx, src.url)
		if err != nil {
			continue
		}
		success++
		combined = append(combined, parseOfficialEvents(raw, src.region, src.source, src.url, src.names)...)
	}
	if success == 0 {
		e.setHealth("macro-events", "degraded · official calendars unavailable; cached events preserved")
		return
	}
	sort.Slice(combined, func(i, j int) bool {
		if combined[i].Date == combined[j].Date {
			return combined[i].StartsAt < combined[j].StartsAt
		}
		return combined[i].Date < combined[j].Date
	})
	if len(combined) > 80 {
		combined = combined[:80]
	}
	e.mu.Lock()
	e.macroEvents = combined
	e.lastUpdated["macro-events"] = time.Now().UnixMilli()
	e.health["macro-events"] = fmt.Sprintf("healthy · %d official calendars", success)
	e.mu.Unlock()
	// Calendar interpretation/Fed/event-risk are derived from this canonical
	// store; publish the updated snapshot without adding another scheduler.
	e.app.broadcastRuntime()
}

func eventModeFor(events []MacroEvent, now time.Time, enabled bool) EventModeState {
	if !enabled {
		return EventModeState{}
	}
	var best *MacroEvent
	for i := range events {
		ev := events[i]
		if ev.Impact != "HIGH" || !ev.TimeKnown || ev.StartsAt <= 0 || macroEventProcessingClass(ev) != "US_MARKET_CRITICAL" {
			continue
		}
		d := time.UnixMilli(ev.StartsAt).Sub(now)
		if d < -time.Hour || d > 15*time.Minute {
			continue
		}
		if best == nil || math.Abs(d.Seconds()) < math.Abs(time.UnixMilli(best.StartsAt).Sub(now).Seconds()) {
			x := ev
			best = &x
		}
	}
	if best == nil {
		return EventModeState{}
	}
	d := time.UnixMilli(best.StartsAt).Sub(now)
	phase := "PREP"
	if d <= 0 {
		phase = "REACTION"
	}
	syms, sectors := affectedContext(*best, nil)
	exp := "NOT AVAILABLE"
	if best.Expected != nil {
		exp = "AVAILABLE"
	}
	return EventModeState{Active: true, EventID: best.ID, Name: best.Name, StartsAt: best.StartsAt, CountdownS: int64(d.Seconds()), Phase: phase, AffectedSymbols: syms, AffectedSectors: sectors, Prepared: phase == "REACTION", QueuePrepared: true, MetadataReady: true, ExpectationStatus: exp, DataPriority: "MARKET_CRITICAL"}
}
func (e *Engine) eventModeLoop(ctx context.Context, alpacaKey, alpacaSecret string) {
	t := time.NewTicker(5 * time.Second)
	defer t.Stop()
	captured := map[string]map[int]bool{}
	warmed := map[string]bool{}
	baselines := map[string]map[string]float64{}
	baselineAt := map[string]int64{}
	for {
		select {
		case <-ctx.Done():
			return
		case now := <-t.C:
			e.app.mu.RLock()
			enabled := e.app.state.Settings.MacroEventModeEnabled
			fredKey := strings.TrimSpace(e.app.secrets.FRED)
			eiaKey := strings.TrimSpace(e.app.secrets.EIA)
			tdKey := strings.TrimSpace(e.app.secrets.TwelveData)
			e.app.mu.RUnlock()
			e.mu.RLock()
			events := clone(e.macroEvents)
			quotes := clone(e.quotes)
			e.mu.RUnlock()
			mode := eventModeFor(events, now, enabled)
			if !mode.Active {
				continue
			}
			var ev MacroEvent
			for _, x := range events {
				if x.ID == mode.EventID {
					ev = x
					break
				}
			}
			syms, sectors := affectedContext(ev, e.trackedSymbols())
			mode.AffectedSymbols = syms
			mode.AffectedSectors = sectors

			if mode.CountdownS > 0 {
				b := map[string]float64{}
				for _, s := range []string{"SPY", "QQQ", "IWM", "VIX", "TLT", "UUP", "GLD", "USO", "SMH"} {
					if q, ok := quotes[s]; ok && q.Price > 0 {
						b[s] = q.Price
					}
				}
				if len(b) > 0 {
					baselines[mode.EventID] = b
					baselineAt[mode.EventID] = now.UnixMilli()
				}
			}
			if mode.Phase == "PREP" && !warmed[mode.EventID] {
				warmed[mode.EventID] = true
				if alpacaKey != "" && alpacaSecret != "" {
					go e.refreshAlpacaHistoryScoped(ctx, alpacaKey, alpacaSecret, syms)
				}
				go e.refreshMacroRouted(ctx, fredKey, eiaKey)
				go e.refreshDirectGlobal(ctx, tdKey)
				e.setHealth("event-mode", "active · market-critical data prioritized; nonessential background refresh deferred")
			}
			if mode.CountdownS > 0 {
				continue
			}
			base := baselines[mode.EventID]
			if len(base) == 0 {
				base = map[string]float64{}
				for s, q := range quotes {
					if q.Price > 0 {
						base[s] = q.Price
					}
				}
				baselines[mode.EventID] = base
				baselineAt[mode.EventID] = now.UnixMilli()
			}
			elapsed := int(-mode.CountdownS)
			for _, off := range []int{5, 30, 60, 300, 900, 3600} {
				if elapsed < off || elapsed > off+6 {
					continue
				}
				if captured[mode.EventID] == nil {
					captured[mode.EventID] = map[int]bool{}
				}
				if captured[mode.EventID][off] {
					continue
				}
				moves := map[string]float64{}
				for s, bp := range base {
					if bp <= 0 {
						continue
					}
					if q, ok := quotes[s]; ok && q.Price > 0 {
						moves[s] = (q.Price/bp - 1) * 100
					}
				}
				target := ev.StartsAt + int64(off)*1000
				lat := now.UnixMilli() - target
				if lat < 0 {
					lat = 0
				}
				e.mu.Lock()
				e.eventReactions = append(e.eventReactions, EventReaction{EventID: mode.EventID, OffsetSec: off, CapturedAt: now.UnixMilli(), BaselineAt: baselineAt[mode.EventID], Moves: moves, Baseline: clone(base), InternalLatencyMs: lat})
				if len(e.eventReactions) > 100 {
					e.eventReactions = e.eventReactions[len(e.eventReactions)-100:]
				}
				e.lastUpdated["event-reactions"] = now.UnixMilli()
				e.mu.Unlock()
				captured[mode.EventID][off] = true
			}
		}
	}
}
