package main

import (
	"fmt"
	"math"
	"strings"
	"time"
)

var proxyDriverMap = map[string]struct {
	Label  string
	Symbol string
	Group  string
}{
	"us_broad":         {"U.S. Broad Market", "SPY", "risk"},
	"us_growth":        {"U.S. Growth", "QQQ", "risk"},
	"us_small":         {"U.S. Small Caps", "IWM", "risk"},
	"korea":            {"Korea / KOSPI proxy", "EWY", "asia"},
	"taiwan":           {"Taiwan / TAIEX proxy", "EWT", "asia-tech"},
	"japan":            {"Japan proxy", "EWJ", "asia"},
	"hong_kong":        {"Hong Kong proxy", "EWH", "china"},
	"china":            {"China proxy", "MCHI", "china"},
	"europe":           {"Europe proxy", "VGK", "europe"},
	"eurozone":         {"Eurozone proxy", "FEZ", "europe"},
	"usd":              {"U.S. Dollar proxy", "UUP", "fx"},
	"usd_jpy":          {"USD/JPY proxy", "FXY", "fx"},
	"usd_cnh":          {"USD/CNH proxy", "CYB", "fx"},
	"eur_usd":          {"EUR/USD proxy", "FXE", "fx"},
	"high_yield":       {"High Yield Credit proxy", "HYG", "credit"},
	"investment_grade": {"Investment Grade Credit proxy", "LQD", "credit"},
	"semiconductors":   {"Global Semiconductor read-through", "SMH", "asia-tech"},
	"rates":            {"Long Duration / Rates proxy", "TLT", "rates"},
	"oil":              {"WTI Oil proxy", "USO", "commodity"},
	"brent":            {"Brent Oil proxy", "BNO", "commodity"},
	"copper":           {"Copper proxy", "CPER", "commodity"},
	"natural_gas":      {"Natural Gas proxy", "UNG", "commodity"},
	"gold":             {"Gold", "GLD", "commodity"},
	"silver":           {"Silver", "SLV", "commodity"},
}

func driverState(change float64, inverse bool) string {
	if inverse {
		change = -change
	}
	switch {
	case change >= 0.8:
		return "SUPPORTIVE"
	case change <= -0.8:
		return "HEADWIND"
	case math.Abs(change) < 0.25:
		return "NEUTRAL"
	default:
		return "CONFLICTING"
	}
}

func quoteProvenance(q Quote) string {
	st := strings.ToUpper(strings.TrimSpace(q.DataState))
	switch st {
	case "LIVE", "CURRENT", "DELAYED", "CACHED", "CACHE", "INDICATIVE":
		if st == "CACHE" {
			return "CACHED"
		}
		return st
	}
	if strings.EqualFold(q.FeedType, "cache") {
		return "CACHED"
	}
	if strings.EqualFold(q.FeedType, "demo") {
		return "DEMO"
	}
	if q.Price > 0 {
		return "CURRENT"
	}
	return "UNAVAILABLE"
}

var broadBreadthUniverse = []string{"SPY", "QQQ", "DIA", "IWM", "XLK", "XLC", "XLY", "XLP", "XLE", "XLF", "XLV", "XLI", "XLB", "XLRE", "XLU"}

func broadBreadthDriver(quotes map[string]Quote) (GlobalDriver, bool) {
	up, down, total := 0, 0, 0
	latest := int64(0)
	for _, sym := range broadBreadthUniverse {
		q, ok := quotes[sym]
		if !ok || q.Price <= 0 {
			continue
		}
		total++
		if q.ChangePercent > .15 {
			up++
		} else if q.ChangePercent < -.15 {
			down++
		}
		if q.UpdatedAt > latest {
			latest = q.UpdatedAt
		}
	}
	if total < 6 {
		return GlobalDriver{}, false
	}
	participation := float64(up-down) / float64(total) * 100
	state := "NEUTRAL"
	if participation >= 30 {
		state = "SUPPORTIVE"
	} else if participation <= -30 {
		state = "HEADWIND"
	} else if math.Abs(participation) >= 12 {
		state = "CONFLICTING"
	}
	return GlobalDriver{Key: "breadth", Label: "Broad Market Breadth", State: state, Value: participation, Source: "Broad ETF universe", Provenance: "REAL PROXY BREADTH", UpdatedAt: latest, Confidence: minInt(88, 55+total*2), Detail: fmt.Sprintf("%d up / %d down across %d broad/sector ETFs; not watchlist breadth", up, down, total)}, true
}

func sectorToneDriver(quotes map[string]Quote) (GlobalDriver, bool) {
	sectors := []string{"XLK", "XLC", "XLY", "XLP", "XLE", "XLF", "XLV", "XLI", "XLB", "XLRE", "XLU"}
	vals := []float64{}
	latest := int64(0)
	for _, sym := range sectors {
		if q, ok := quotes[sym]; ok && q.Price > 0 {
			vals = append(vals, q.ChangePercent)
			if q.UpdatedAt > latest {
				latest = q.UpdatedAt
			}
		}
	}
	if len(vals) < 5 {
		return GlobalDriver{}, false
	}
	avg := 0.0
	for _, v := range vals {
		avg += v
	}
	avg /= float64(len(vals))
	return GlobalDriver{Key: "sectors", Label: "Sector Participation", State: driverState(avg, false), Value: avg, ChangePercent: avg, Source: "SPDR sector ETFs", Provenance: "REAL PROXY", UpdatedAt: latest, Confidence: minInt(88, 58+len(vals)*2), Detail: fmt.Sprintf("Average %.2f%% across %d sectors", avg, len(vals))}, true
}

func deriveGlobalMarketContext(quotes map[string]Quote, direct map[string]GlobalDriver, metrics map[string]MacroMetric, mode string) GlobalMarketContext {
	now := time.Now().UnixMilli()
	providerMode := strings.ToLower(strings.TrimSpace(mode))
	if providerMode == "" {
		providerMode = "auto"
	}
	drivers := map[string]GlobalDriver{}

	for k, d := range direct {
		if d.Value <= 0 {
			continue
		}
		prov := strings.ToUpper(d.Provenance)
		isOfficial := strings.Contains(prov, "OFFICIAL") || strings.Contains(prov, "PUBLIC")
		isDirect := strings.Contains(prov, "DIRECT")
		include := false
		switch providerMode {
		case "proxy":
			include = false
		case "direct":
			include = isDirect
		case "free-first":
			include = isOfficial
		default:
			include = isDirect || isOfficial
		}
		if include {
			d.IsProxy = false
			if d.Provenance == "" {
				d.Provenance = "DIRECT PROVIDER"
			}
			drivers[k] = d
		}
	}

	if providerMode != "direct" {
		if d, ok := broadBreadthDriver(quotes); ok {
			drivers[d.Key] = d
		}
		if d, ok := sectorToneDriver(quotes); ok {
			drivers[d.Key] = d
		}
		for key, spec := range proxyDriverMap {
			q, ok := quotes[spec.Symbol]
			if !ok || q.Price <= 0 {
				continue
			}
			inverse := key == "usd"
			state := driverState(q.ChangePercent, inverse)
			conf := 72
			prov := quoteProvenance(q)
			if prov == "CACHED" {
				conf = 45
			}
			if prov == "DEMO" {
				conf = 35
			}
			proxyKey := key
			label := spec.Label

			if _, exists := drivers[key]; exists {
				proxyKey = key + "_proxy"
				label = spec.Label + " · live proxy"
			}
			d := GlobalDriver{Key: proxyKey, Label: label, State: state, Value: q.Price, ChangePercent: q.ChangePercent, Source: q.Source, Provenance: "LIVE PROXY", UpdatedAt: q.UpdatedAt, Confidence: conf, Detail: spec.Symbol + " proxy · " + prov, Underlying: strings.TrimSuffix(spec.Label, " proxy"), ProviderSymbol: spec.Symbol, IsProxy: true, Session: "CURRENT PROXY"}
			if q.UpdatedAt == 0 {
				d.UpdatedAt = now
			}
			drivers[proxyKey] = d
		}
	}

	if v, ok := quotes["VIX"]; ok && v.Price > 0 {
		state := "NEUTRAL"
		if v.Price >= 25 || v.ChangePercent >= 10 {
			state = "HEADWIND"
		} else if v.Price <= 15 && v.ChangePercent <= 0 {
			state = "SUPPORTIVE"
		} else if math.Abs(v.ChangePercent) >= 5 {
			state = "CONFLICTING"
		}
		drivers["volatility"] = GlobalDriver{Key: "volatility", Label: "Volatility Tone", State: state, Value: v.Price, ChangePercent: v.ChangePercent, Source: v.Source, Provenance: quoteProvenance(v), UpdatedAt: v.UpdatedAt, Confidence: 90, Detail: "Dedicated VIX path"}
	}

	for _, id := range []string{"UST10Y", "DGS10"} {
		if m, ok := metrics[id]; ok && m.Status != "UNAVAILABLE" && m.Value > 0 {
			state := "NEUTRAL"
			delta := m.Change5D
			if delta >= .12 {
				state = "HEADWIND"
			} else if delta <= -.12 {
				state = "SUPPORTIVE"
			} else if m.Value >= 5 {
				state = "HEADWIND"
			}
			drivers["rates_10y"] = GlobalDriver{Key: "rates_10y", Label: "U.S. 10Y Yield", State: state, Value: m.Value, ChangePercent: delta, Source: m.Source, Provenance: m.Provenance, UpdatedAt: m.UpdatedAt, Confidence: 90, Detail: fmt.Sprintf("%.2f%% · 5D %+0.2f", m.Value, delta)}
			break
		}
	}

	// Score by logical driver family so direct + official-close + proxy evidence can all
	// remain visible without being counted as three independent markets.
	type familyVote struct{ sum, weight float64 }
	families := map[string]familyVote{}
	supportive, headwind := 0, 0
	for k, d := range drivers {
		base := strings.TrimSuffix(strings.TrimSuffix(k, "_official_close"), "_proxy")
		w := 1.0
		switch base {
		case "us_broad", "us_growth", "us_small", "taiwan", "semiconductors", "nq_future", "es_future", "volatility":
			w = 1.35
		case "breadth", "sectors", "rates_10y", "high_yield":
			w = 1.2
		case "silver", "natural_gas":
			w = .55
		}

		if strings.HasSuffix(k, "_official_close") || strings.HasSuffix(k, "_proxy") {
			w *= .65
		}
		vote := 0.0
		switch strings.ToUpper(d.State) {
		case "SUPPORTIVE":
			vote = 1
			supportive++
		case "HEADWIND":
			vote = -1
			headwind++
		}
		f := families[base]
		f.sum += vote * w
		f.weight += w
		families[base] = f
	}
	score, weighted := 0.0, 0.0
	for base, f := range families {
		if f.weight == 0 {
			continue
		}
		familyValue := f.sum / f.weight
		fw := 1.0
		switch base {
		case "us_broad", "us_growth", "us_small", "taiwan", "semiconductors", "nq_future", "es_future", "volatility":
			fw = 1.35
		case "breadth", "sectors", "rates_10y", "high_yield":
			fw = 1.2
		case "silver", "natural_gas":
			fw = .55
		}
		score += familyValue * fw
		weighted += fw
	}
	norm := 0.0
	if weighted > 0 {
		norm = score / weighted
	}
	tone := "NEUTRAL"
	switch {
	case norm >= .38:
		tone = "RISK-ON"
	case norm >= .12:
		tone = "CONSTRUCTIVE"
	case norm <= -.38:
		tone = "RISK-OFF"
	case norm <= -.12:
		tone = "CAUTIOUS"
	}
	confidence := 0
	if len(drivers) > 0 {
		sum := 0
		for _, d := range drivers {
			sum += d.Confidence
		}
		confidence = minInt(95, sum/len(drivers))
	}
	summary := fmt.Sprintf("%s · %d supportive / %d headwind evidence points", tone, supportive, headwind)
	return GlobalMarketContext{Tone: tone, Confidence: confidence, Drivers: drivers, UpdatedAt: now, Mode: strings.ToUpper(defaultString(mode, "AUTO")), Summary: summary}
}
func defaultString(v, def string) string {
	if strings.TrimSpace(v) == "" {
		return def
	}
	return v
}

func buildCapabilities(s Settings, sec Secrets, mode string, global GlobalMarketContext, options map[string]OptionsContext, metrics map[string]MacroMetric, events []MacroEvent) []CapabilityStatus {
	now := time.Now().UnixMilli()
	directCount, publicCount, proxyCount, futureCount := 0, 0, 0, 0
	for k, d := range global.Drivers {
		if strings.HasSuffix(k, "_future") {
			futureCount++
		}
		prov := strings.ToUpper(d.Provenance)
		if d.IsProxy || strings.Contains(prov, "PROXY") {
			proxyCount++
		} else if strings.Contains(prov, "OFFICIAL") || strings.Contains(prov, "PUBLIC") {
			publicCount++
		} else if strings.Contains(prov, "DIRECT") {
			directCount++
		}
	}
	metricActive := func(keys ...string) bool {
		for _, k := range keys {
			if m, ok := metrics[k]; ok && m.Status != "UNAVAILABLE" && m.Value != 0 {
				return true
			}
		}
		return false
	}
	caps := []CapabilityStatus{
		{Capability: "U.S. Equities", Source: "Alpaca + Finnhub + Twelve Data", Mode: "ROUTED PRIMARY + FALLBACK", Status: "ACTIVE", UpdatedAt: now, Details: []string{"Alpaca IEX/snapshots primary", "Finnhub secondary live/recovery", "Twelve Data tertiary recovery", "VIX remains a dedicated index path"}},
		{Capability: "Global Indices", Source: "Direct provider + official/public + U.S. ETF proxies", Mode: strings.ToUpper(defaultString(s.GlobalProviderMode, "auto")), Status: capStatus(len(global.Drivers) > 0), UpdatedAt: global.UpdatedAt, Details: []string{fmt.Sprintf("Direct instruments active: %d", directCount), fmt.Sprintf("Official/public closes active: %d", publicCount), fmt.Sprintf("Real proxy instruments active: %d", proxyCount), "Direct → official/public → real proxy → real cache → unavailable", "Proxy is never presented as the underlying index"}},
		{Capability: "U.S. Futures", Source: "Direct provider interface (ES/NQ/RTY)", Mode: "PREMIUM/DIRECT READY", Status: func() string {
			if futureCount > 0 {
				return "ACTIVE"
			}
			if sec.TwelveData != "" {
				return "DEGRADED"
			}
			return "NOT CONFIGURED"
		}(), UpdatedAt: now, Details: []string{fmt.Sprintf("Direct futures currently active: %d", futureCount), "SPY/QQQ/IWM overnight architecture remains the real fallback context", "Premium/licensed entitlement can be added without UI redesign"}},
		{Capability: "Macro Events", Source: "Fed / BLS / BEA / Census / DOL / ISM / Eurostat / ECB / BOJ / China NBS / PBOC / Customs", Mode: "OFFICIAL/PUBLIC", Status: capStatus(len(events) > 0), UpdatedAt: now, Details: []string{"UPCOMING → RELEASED → MARKET REACTION → RESOLVED", "High-impact Event Mode prepares context at approximately T−15", "Unknown official times remain date-only"}},
		{Capability: "Macro Actuals", Source: "Treasury / BLS / BEA / EIA / FRED", Mode: "FREE/OFFICIAL", Status: capStatus(metricActive("UST10Y", "DGS10", "CPI_INDEX", "BEA_GDP", "WTI_OFFICIAL")), UpdatedAt: now, Details: []string{"Treasury yields and real-yield path", "BLS CPI/labor actuals", "BEA GDP/PCE public releases", "EIA energy actuals when free key configured"}},
		{Capability: "Fast Macro Consensus + Actual", Source: "Future premium provider", Mode: "PROVIDER INTERFACE READY", Status: "NOT CONFIGURED", UpdatedAt: now, Details: []string{"Premium upgrade available", "Free/official sources are used for correctness/context", "No institutional first-second latency is claimed"}},
		{Capability: "Options Intelligence", Source: "Alpaca Options", Mode: strings.ToUpper(defaultString(s.OptionsDataMode, "auto")), Status: capStatus(len(options) > 0), UpdatedAt: now, Details: []string{"OPRA when entitled; indicative fallback in AUTO", "IV / ΔIV / put-call volume / expected move where sourced", "Open interest/unusual flow remains unavailable unless a real provider supplies it", "Context only — deterministic Scores unchanged"}},
		{Capability: "Signal Validation", Source: "Canonical snapshots + real historical bars", Mode: "VALIDATION", Status: "ACTIVE", UpdatedAt: now, Details: []string{"Research + Queue + Global/Macro + Options + Readiness context recorded", "1D/3D/5D/10D outcomes and MFE/MAE when later real bars exist", "Never auto-reweights trading formulas"}},
	}
	if sec.FRED != "" {
		caps = append(caps, CapabilityStatus{Capability: "Rates & Credit", Source: "FRED + Treasury", Mode: "OFFICIAL API/PUBLIC", Status: capStatus(metricActive("DGS10", "BAMLH0A0HYM2", "UST10Y")), UpdatedAt: now, Details: []string{"5D/20D and 1M/3M trend deltas retained for horizon interpretation"}})
	} else {
		caps = append(caps, CapabilityStatus{Capability: "Rates & Credit", Source: "Treasury + optional FRED", Mode: "FREE/OFFICIAL", Status: func() string {
			if metricActive("UST10Y") {
				return "ACTIVE"
			}
			return "NOT CONFIGURED"
		}(), UpdatedAt: now, Details: []string{"Treasury is managed internally", "Optional free FRED API key expands rates/credit history"}})
	}
	if sec.TwelveData != "" {
		caps = append(caps, CapabilityStatus{Capability: "Direct Global Provider", Source: "Twelve Data", Mode: "OPTIONAL", Status: func() string {
			if directCount > 0 {
				return "ACTIVE"
			}
			return "DEGRADED"
		}(), UpdatedAt: now, Details: []string{"Capability and entitlement are verified by representative requests", "Real ETF proxies remain fallback in AUTO"}})
	} else {
		caps = append(caps, CapabilityStatus{Capability: "Direct Global Provider", Source: "Twelve Data / future licensed feed", Mode: "OPTIONAL", Status: "NOT CONFIGURED", UpdatedAt: now, Details: []string{"Premium/direct upgrade available", "Provider interfaces are implemented now; paid entitlement is optional"}})
	}
	return caps
}
func capStatus(ok bool) string {
	if ok {
		return "ACTIVE"
	}
	return "DEGRADED"
}
