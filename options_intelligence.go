package main

import (
	"context"
	"fmt"
	"math"
	"net/http"
	"net/url"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"
)

type alpacaOptionChainResponse struct {
	Snapshots map[string]struct {
		DailyBar *struct {
			V float64 `json:"v"`
		} `json:"dailyBar"`
		ImpliedVolatility float64 `json:"impliedVolatility"`
		Greeks            *struct {
			Gamma float64 `json:"gamma"`
		} `json:"greeks"`
	} `json:"snapshots"`
	NextPageToken string `json:"next_page_token"`
}

type alpacaOptionContractsResponse struct {
	OptionContracts []struct {
		Symbol           string `json:"symbol"`
		Type             string `json:"type"`
		OpenInterest     string `json:"open_interest"`
		OpenInterestDate string `json:"open_interest_date"`
	} `json:"option_contracts"`
	NextPageToken string `json:"next_page_token"`
}

type optionOpenInterestRecord struct {
	OpenInterest float64
	Date         time.Time
}

func fetchOptionOpenInterest(ctx context.Context, key, secret, symbol string, now time.Time) (map[string]optionOpenInterestRecord, error) {
	from := now.Format("2006-01-02")
	to := now.AddDate(0, 0, 45).Format("2006-01-02")
	base := alpacaTradingBaseURL + "/v2/options/contracts?underlying_symbols=" + url.QueryEscape(normalizeSymbol(symbol)) +
		"&expiration_date_gte=" + url.QueryEscape(from) + "&expiration_date_lte=" + url.QueryEscape(to) + "&limit=10000"
	out := map[string]optionOpenInterestRecord{}
	pageToken := ""
	for page := 0; page < 20; page++ {
		raw := base
		if pageToken != "" {
			raw += "&page_token=" + url.QueryEscape(pageToken)
		}
		var resp alpacaOptionContractsResponse
		if err := getJSON(ctx, &http.Client{Timeout: 15 * time.Second}, raw, map[string]string{"APCA-API-KEY-ID": key, "APCA-API-SECRET-KEY": secret}, &resp); err != nil {
			return nil, err
		}
		for _, c := range resp.OptionContracts {
			oi, err := strconv.ParseFloat(strings.TrimSpace(c.OpenInterest), 64)
			if err != nil || oi < 0 {
				continue
			}
			d, err := time.Parse("2006-01-02", strings.TrimSpace(c.OpenInterestDate))
			if err != nil {
				continue
			}
			out[strings.ToUpper(strings.TrimSpace(c.Symbol))] = optionOpenInterestRecord{OpenInterest: oi, Date: d}
		}
		next := strings.TrimSpace(resp.NextPageToken)
		if next == "" {
			break
		}
		if next == pageToken {
			return nil, fmt.Errorf("option open-interest pagination did not advance")
		}
		pageToken = next
		if page == 19 {
			return nil, fmt.Errorf("option open-interest pagination exceeded safety limit")
		}
	}
	if len(out) == 0 {
		return nil, fmt.Errorf("no option open-interest records returned")
	}
	return out, nil
}

var occRx = regexp.MustCompile(`^([A-Z]{1,6})(\d{6})([CP])(\d{8})$`)

func parseOCCContractDetails(s string) (kind string, expiry time.Time, strike float64, ok bool) {
	m := occRx.FindStringSubmatch(strings.ToUpper(strings.TrimSpace(s)))
	if len(m) == 0 {
		return
	}
	t, err := time.Parse("060102", m[2])
	if err != nil {
		return
	}
	rawStrike, err := strconv.Atoi(m[4])
	if err != nil {
		return
	}
	return m[3], t, float64(rawStrike) / 1000.0, true
}

type gexAccumulator struct {
	Call, Put, OI float64
	Contracts     int
}

func finalizeGEXStructure(strikes map[float64]gexAccumulator, expirations map[string]gexAccumulator, underlying float64) (major []GEXStrikeLevel, zones []GEXConcentrationZone, expiry []GEXExpirationLevel, flip *float64, flipMethod string) {
	levels := make([]GEXStrikeLevel, 0, len(strikes))
	totalAbs := 0.0
	for strike, a := range strikes {
		net := a.Call + a.Put
		abs := math.Abs(net)
		totalAbs += abs
		levels = append(levels, GEXStrikeLevel{Strike: strike, CallGEX: a.Call, PutGEX: a.Put, NetGEX: net, AbsoluteGEX: abs, OpenInterest: a.OI, Contracts: a.Contracts})
	}
	sort.Slice(levels, func(i, j int) bool { return levels[i].Strike < levels[j].Strike })
	// Structural flip is only published when adjacent strike-level signed proxy
	// values cross zero. Choose the crossing nearest spot; this is not dealer positioning.
	bestDist := math.MaxFloat64
	for i := 1; i < len(levels); i++ {
		a, b := levels[i-1], levels[i]
		if a.NetGEX == 0 {
			x := a.Strike
			d := math.Abs(x - underlying)
			if d < bestDist {
				flip = &x
				bestDist = d
			}
		} else if (a.NetGEX < 0 && b.NetGEX > 0) || (a.NetGEX > 0 && b.NetGEX < 0) {
			denom := math.Abs(a.NetGEX) + math.Abs(b.NetGEX)
			if denom > 0 {
				x := a.Strike + (b.Strike-a.Strike)*(math.Abs(a.NetGEX)/denom)
				d := math.Abs(x - underlying)
				if d < bestDist {
					xx := x
					flip = &xx
					bestDist = d
				}
			}
		}
	}
	if flip != nil {
		flipMethod = "Adjacent strike sign-change interpolation in structural signed gamma×OI proxy; not measured dealer gamma."
	}
	major = append(major, levels...)
	sort.Slice(major, func(i, j int) bool {
		if major[i].AbsoluteGEX == major[j].AbsoluteGEX {
			return major[i].Strike < major[j].Strike
		}
		return major[i].AbsoluteGEX > major[j].AbsoluteGEX
	})
	if len(major) > 8 {
		major = major[:8]
	}
	// Concentration zones are non-overlapping clusters around the largest retained
	// strikes using a spot-scaled width. This keeps computation bounded/linear.
	used := map[float64]bool{}
	width := math.Max(underlying*0.015, 1)
	for _, m := range major {
		if len(zones) >= 3 || used[m.Strike] {
			continue
		}
		z := GEXConcentrationZone{LowStrike: m.Strike - width, HighStrike: m.Strike + width}
		for _, lv := range levels {
			if used[lv.Strike] || lv.Strike < z.LowStrike || lv.Strike > z.HighStrike {
				continue
			}
			used[lv.Strike] = true
			z.NetGEX += lv.NetGEX
			z.AbsoluteGEX += lv.AbsoluteGEX
			z.StrikeCount++
		}
		if totalAbs > 0 {
			z.SharePct = z.AbsoluteGEX / totalAbs * 100
		}
		zones = append(zones, z)
	}
	for exp, a := range expirations {
		net := a.Call + a.Put
		expiry = append(expiry, GEXExpirationLevel{Expiration: exp, CallGEX: a.Call, PutGEX: a.Put, NetGEX: net, AbsoluteGEX: math.Abs(net), OpenInterest: a.OI, Contracts: a.Contracts})
	}
	sort.Slice(expiry, func(i, j int) bool { return expiry[i].Expiration < expiry[j].Expiration })
	if len(expiry) > 8 {
		expiry = expiry[:8]
	}
	return
}

func aggregateOptions(symbol, feed string, underlying float64, p alpacaOptionChainResponse, oi map[string]optionOpenInterestRecord, oiErr error) OptionsContext {
	now := time.Now()
	o := OptionsContext{Symbol: normalizeSymbol(symbol), Provider: "Alpaca Options", Feed: strings.ToUpper(feed), State: "CURRENT", Bias: "NEUTRAL", UpdatedAt: now.UnixMilli(), Provenance: "REAL OPTIONS SNAPSHOT", GEXState: "UNAVAILABLE", UnderlyingPrice: underlying}
	strikeGEX := map[float64]gexAccumulator{}
	expirationGEX := map[string]gexAccumulator{}
	var ivSum float64
	ivN := 0
	var nearest time.Time
	var oldestOIDate time.Time
	type ivPoint struct {
		exp        time.Time
		strike, iv float64
	}
	points := []ivPoint{}
	for contract, snap := range p.Snapshots {
		kind, exp, strike, ok := parseOCCContractDetails(contract)
		if !ok {
			continue
		}
		vol := 0.0
		if snap.DailyBar != nil {
			vol = snap.DailyBar.V
		}
		if kind == "C" {
			o.CallContracts++
			o.CallVolume += vol
		} else {
			o.PutContracts++
			o.PutVolume += vol
		}
		if snap.ImpliedVolatility > 0 && snap.ImpliedVolatility < 10 {
			ivSum += snap.ImpliedVolatility
			ivN++
			if exp.After(now) {
				points = append(points, ivPoint{exp: exp, strike: strike, iv: snap.ImpliedVolatility})
			}
		}
		if exp.After(now) && (nearest.IsZero() || exp.Before(nearest)) {
			nearest = exp
		}
		if snap.Greeks != nil && snap.Greeks.Gamma != 0 {
			o.GammaContracts++
			if rec, found := oi[strings.ToUpper(contract)]; found {
				o.OpenInterestContracts++
				o.GammaOIContracts++
				if oldestOIDate.IsZero() || rec.Date.Before(oldestOIDate) {
					oldestOIDate = rec.Date
				}
				if underlying > 0 && rec.OpenInterest > 0 {
					gex := snap.Greeks.Gamma * rec.OpenInterest * 100 * underlying * underlying * 0.01
					sa := strikeGEX[strike]
					ea := expirationGEX[exp.Format("2006-01-02")]
					sa.OI += rec.OpenInterest
					ea.OI += rec.OpenInterest
					sa.Contracts++
					ea.Contracts++
					if kind == "C" {
						o.CallGEX += gex
						sa.Call += gex
						ea.Call += gex
					} else {
						o.PutGEX -= gex
						sa.Put -= gex
						ea.Put -= gex
					}
					strikeGEX[strike] = sa
					expirationGEX[exp.Format("2006-01-02")] = ea
				}
			}
		}
	}
	if len(oi) > 0 {
		o.OpenInterestContracts = len(oi)
	}
	if o.GammaContracts > 0 {
		o.GammaOICoveragePct = float64(o.GammaOIContracts) / float64(o.GammaContracts) * 100
	}
	if !oldestOIDate.IsZero() {
		o.OpenInterestDate = oldestOIDate.Format("2006-01-02")
	}
	if oiErr != nil {
		o.Limitations = append(o.Limitations, "GEX unavailable because Alpaca contract open interest could not be reconciled: "+oiErr.Error())
	} else if o.GammaContracts == 0 {
		o.Limitations = append(o.Limitations, "GEX unavailable because the option snapshot did not contain usable gamma values.")
	} else {
		oiAge := 999.0
		futureSkew := false
		if !oldestOIDate.IsZero() {
			oiAge = now.Sub(oldestOIDate).Hours() / 24
			futureSkew = oldestOIDate.After(now.Add(24 * time.Hour))
		}
		if o.GammaOIContracts >= 30 && o.GammaOICoveragePct >= 60 && oiAge >= -1 && oiAge <= 7 && !futureSkew && underlying > 0 {
			o.NetGEX = o.CallGEX + o.PutGEX
			o.MajorGammaStrikes, o.GammaZones, o.ExpirationGEX, o.GammaFlip, o.GammaFlipMethod = finalizeGEXStructure(strikeGEX, expirationGEX, underlying)
			o.GEXState = "AVAILABLE"
			o.GEXQuality = "MODERATE"
			if strings.EqualFold(feed, "opra") && o.GammaOICoveragePct >= 80 {
				o.GEXQuality = "HIGH"
			}
			o.Limitations = append(o.Limitations, "GEX is a structural signed gamma × open-interest proxy. Open interest does not reveal dealer long/short positioning, so this is not measured dealer exposure.")
		} else {
			o.Limitations = append(o.Limitations, fmt.Sprintf("GEX withheld: matched gamma+OI %d/%d (%.0f%%); minimum is 30 contracts, 60%% coverage, current OI and a valid underlying price.", o.GammaOIContracts, o.GammaContracts, o.GammaOICoveragePct))
		}
	}
	if o.CallVolume > 0 {
		o.PutCallVolume = o.PutVolume / o.CallVolume
	}
	if ivN > 0 {
		o.AverageIV = ivSum / float64(ivN)
	}
	if !nearest.IsZero() {
		o.NearestExpiration = nearest.Format("2006-01-02")

		atmSum, atmN := 0.0, 0
		band := math.Max(underlying*.05, 1)
		for _, pt := range points {
			if pt.exp.Equal(nearest) && underlying > 0 && math.Abs(pt.strike-underlying) <= band {
				atmSum += pt.iv
				atmN++
			}
		}
		atmIV := 0.0
		if atmN > 0 {
			atmIV = atmSum / float64(atmN)
		}
		days := math.Max(1, nearest.Sub(now).Hours()/24)
		if underlying > 0 && atmIV > 0 {
			o.ExpectedMove = underlying * atmIV * math.Sqrt(days/365)
		} else if underlying > 0 {
			o.Limitations = append(o.Limitations, "Expected move unavailable because nearest-expiry near-the-money IV was not present.")
		}
	}
	switch {
	case o.PutCallVolume >= 1.35:
		o.Bias = "BEARISH"
	case o.PutCallVolume > 0 && o.PutCallVolume <= .72:
		o.Bias = "BULLISH"
	case o.CallVolume+o.PutVolume == 0:
		o.Bias = "INCOMPLETE"
	}
	if strings.EqualFold(feed, "indicative") {
		o.State = "DELAYED/INDICATIVE"
		o.Provenance = "ALPACA INDICATIVE"
	}
	return o
}
func fetchOptionsContext(ctx context.Context, key, secret, symbol, mode string, underlying float64) (OptionsContext, error) {
	if key == "" || secret == "" {
		return OptionsContext{}, fmt.Errorf("Alpaca credentials are required")
	}
	feeds := []string{"opra"}
	mode = strings.ToLower(strings.TrimSpace(mode))
	if mode == "indicative" {
		feeds = []string{"indicative"}
	} else if mode == "auto" || mode == "" {
		feeds = []string{"opra", "indicative"}
	}
	var last error
	for _, feed := range feeds {
		raw := alpacaDataBaseURL + "/v1beta1/options/snapshots/" + url.PathEscape(symbol) + "?feed=" + url.QueryEscape(feed) + "&limit=1000"
		var p alpacaOptionChainResponse
		err := getJSON(ctx, &http.Client{Timeout: 15 * time.Second}, raw, map[string]string{"APCA-API-KEY-ID": key, "APCA-API-SECRET-KEY": secret}, &p)
		if err == nil && len(p.Snapshots) > 0 {
			oi, oiErr := fetchOptionOpenInterest(ctx, key, secret, symbol, time.Now())
			return aggregateOptions(symbol, feed, underlying, p, oi, oiErr), nil
		}
		if err == nil {
			err = fmt.Errorf("no option snapshots returned")
		}
		last = err
	}
	return OptionsContext{}, last
}

func applyOptionsIVChange(prev, cur OptionsContext) OptionsContext {
	if prev.AverageIV > 0 && cur.AverageIV > 0 {
		cur.IVChange = (cur.AverageIV - prev.AverageIV) * 100
	}
	return cur
}

func (e *Engine) refreshOptions(ctx context.Context, key, secret string) {
	e.app.mu.RLock()
	st := clone(e.app.state)
	e.app.mu.RUnlock()
	mode := st.Settings.OptionsDataMode
	if mode == "off" {
		e.setHealth("options", "off")
		return
	}
	syms := []string{}
	selected := normalizeSymbol(st.UI.SelectedTicker)
	if selected != "" {
		syms = append(syms, selected)
	}
	for _, deskID := range []string{"day", "swing", "long"} {
		for _, wl := range st.Watchlists {
			if wl.ID != deskID {
				continue
			}
			limit := minInt(2, len(wl.Symbols))
			for _, s := range wl.Symbols[:limit] {
				syms = append(syms, normalizeSymbol(s))
			}
		}
	}
	syms = uniqueSymbols(syms)
	okCount := 0
	var provider OptionsIntelligenceProvider = alpacaOptionsProvider{key: key, secret: secret}
	for _, s := range syms {
		if s == "" || s == "VIX" {
			continue
		}
		e.mu.RLock()
		under := e.quotes[s].Price
		e.mu.RUnlock()
		o, err := provider.Snapshot(ctx, s, mode, under)
		if err != nil {
			continue
		}
		e.mu.Lock()
		if prev, ok := e.options[s]; ok {
			o = applyOptionsIVChange(prev, o)
		}
		e.options[s] = o
		e.lastUpdated["options"] = o.UpdatedAt
		e.mu.Unlock()
		okCount++
	}
	if okCount > 0 {
		e.setHealth("options", fmt.Sprintf("healthy · %d symbols", okCount))
	} else {
		e.setHealth("options", "degraded · real options unavailable")
	}
}
