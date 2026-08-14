package main

import (
	"context"
	"net/http"
	"net/url"
	"strings"
	"time"
)

func (e *Engine) ensureDemoSymbol(symbol string) {
	symbol = normalizeSymbol(symbol)
	if symbol == "" {
		return
	}
	e.app.mu.RLock()
	st := clone(e.app.state)
	e.app.mu.RUnlock()
	daySet := symbolSetForWatchlist(st, st.Settings.DayWatchlistID, st.Settings.DayEnabled)
	swingSet := symbolSetForWatchlist(st, st.Settings.SwingWatchlistID, st.Settings.SwingEnabled)
	longSet := symbolSetForWatchlist(st, st.Settings.LongWatchlistID, st.Settings.LongEnabled)
	discoverySet := symbolSetForWatchlist(st, st.Settings.DiscoveryWatchlistID, true)
	if !daySet[symbol] && !swingSet[symbol] && !longSet[symbol] && !discoverySet[symbol] && !contains(generalSymbols, symbol) && normalizeSymbol(st.UI.SelectedTicker) != symbol {
		return
	}
	e.mu.RLock()
	existing := e.quotes[symbol]
	e.mu.RUnlock()
	price := existing.Price
	if price <= 0 {

		seed := 0
		for _, r := range symbol {
			seed += int(r)
		}
		price = 40 + float64(seed%260)
		pc := price * .995
		e.updateQuote(symbol, Quote{Price: price, PreviousClose: pc, SessionClose: pc, PriorSessionClose: pc * .995, Open: pc, High: price * 1.006, Low: price * .994, ProviderTimestamp: time.Now().UnixMilli()}, "demo")
	}
	e.seedDemoBars(symbol, price, daySet[symbol], swingSet[symbol] || longSet[symbol], swingSet[symbol] || longSet[symbol], longSet[symbol], daySet[symbol] || swingSet[symbol] || longSet[symbol] || discoverySet[symbol])
	e.app.broadcastRuntime()
}

func (e *Engine) refreshSingleFinnhubSnapshot(ctx context.Context, key, symbol string) {
	var q finnhubQuoteResponse
	if e.finnhubJSONForSymbol(ctx, key, symbol, "/quote?symbol="+url.QueryEscape(symbol), &q) == nil && q.Current > 0 {
		e.mergeFinnhubSnapshot(symbol, q)
	}
}

func (e *Engine) refreshSingleAlpacaSnapshot(ctx context.Context, key, secret, symbol string) {
	session := marketSessionET(time.Now())
	if session == "closed" || session == "weekend" {
		return
	}
	feed := "iex"
	if session == "overnight" {
		e.app.mu.RLock()
		mode := strings.ToLower(strings.TrimSpace(e.app.state.Settings.OvernightDataMode))
		e.app.mu.RUnlock()
		if mode == "indicative" {
			feed = "overnight"
		} else {
			feed = "boats"
		}
	}
	fetch := func(useFeed string) (map[string]alpacaLiveSnapshot, error) {
		raw := "https://data.alpaca.markets/v2/stocks/snapshots?symbols=" + url.QueryEscape(symbol) + "&feed=" + url.QueryEscape(useFeed)
		var payload map[string]alpacaLiveSnapshot
		err := getJSON(ctx, &http.Client{Timeout: 10 * time.Second}, raw, map[string]string{"APCA-API-KEY-ID": key, "APCA-API-SECRET-KEY": secret}, &payload)
		return payload, err
	}
	payload, err := fetch(feed)
	if err != nil && session == "overnight" && feed == "boats" {
		e.app.mu.RLock()
		mode := strings.ToLower(strings.TrimSpace(e.app.state.Settings.OvernightDataMode))
		e.app.mu.RUnlock()
		if mode == "auto" {
			feed = "overnight"
			payload, err = fetch(feed)
		}
	}
	if err != nil {
		return
	}
	snap, ok := payload[symbol]
	if !ok {
		return
	}
	price, stamp, kind := alpacaSnapshotPrice(snap, feed, session)
	if price > 0 {
		e.mergeAlpacaLiveSnapshot(symbol, price, stamp, snap, feed, kind)
	}
}

// Finnhub is v15's secondary live stream. The desired subscription set includes
// normal overflow plus dynamic failover for stale/unavailable Alpaca-primary symbols.
func (e *Engine) liveSymbols() []string {
	return append([]string{}, e.effectiveFinnhubSymbols()...)
}

func nthWeekday(year int, month time.Month, weekday time.Weekday, n int, loc *time.Location) time.Time {
	d := time.Date(year, month, 1, 12, 0, 0, 0, loc)
	offset := (int(weekday) - int(d.Weekday()) + 7) % 7
	return d.AddDate(0, 0, offset+(n-1)*7)
}

func lastWeekday(year int, month time.Month, weekday time.Weekday, loc *time.Location) time.Time {
	d := time.Date(year, month+1, 0, 12, 0, 0, 0, loc)
	for d.Weekday() != weekday {
		d = d.AddDate(0, 0, -1)
	}
	return d
}

func observedFixedHoliday(year int, month time.Month, day int, loc *time.Location) time.Time {
	d := time.Date(year, month, day, 12, 0, 0, 0, loc)
	if d.Weekday() == time.Saturday {
		return d.AddDate(0, 0, -1)
	}
	if d.Weekday() == time.Sunday {
		return d.AddDate(0, 0, 1)
	}
	return d
}

func easterSunday(year int, loc *time.Location) time.Time {

	a := year % 19
	b := year / 100
	c := year % 100
	d := b / 4
	e := b % 4
	f := (b + 8) / 25
	g := (b - f + 1) / 3
	h := (19*a + b - d - g + 15) % 30
	i := c / 4
	k := c % 4
	l := (32 + 2*e + 2*i - h - k) % 7
	m := (a + 11*h + 22*l) / 451
	month := (h + l - 7*m + 114) / 31
	day := ((h + l - 7*m + 114) % 31) + 1
	return time.Date(year, time.Month(month), day, 12, 0, 0, 0, loc)
}

func sameDate(a, b time.Time) bool {
	ay, am, ad := a.Date()
	by, bm, bd := b.Date()
	return ay == by && am == bm && ad == bd
}

func isUSMarketHoliday(et time.Time) bool {
	loc := et.Location()
	y := et.Year()
	holidays := []time.Time{
		observedFixedHoliday(y, time.January, 1, loc),
		nthWeekday(y, time.January, time.Monday, 3, loc),
		nthWeekday(y, time.February, time.Monday, 3, loc),
		easterSunday(y, loc).AddDate(0, 0, -2),
		lastWeekday(y, time.May, time.Monday, loc),
		observedFixedHoliday(y, time.June, 19, loc),
		observedFixedHoliday(y, time.July, 4, loc),
		nthWeekday(y, time.September, time.Monday, 1, loc),
		nthWeekday(y, time.November, time.Thursday, 4, loc),
		observedFixedHoliday(y, time.December, 25, loc),
	}

	holidays = append(holidays, observedFixedHoliday(y+1, time.January, 1, loc))
	for _, h := range holidays {
		if sameDate(et, h) {
			return true
		}
	}
	return false
}

func isUSEarlyClose(et time.Time) bool {
	loc := et.Location()
	y := et.Year()
	thanksgiving := nthWeekday(y, time.November, time.Thursday, 4, loc)
	if sameDate(et, thanksgiving.AddDate(0, 0, 1)) {
		return true
	}

	if et.Month() == time.July && et.Day() == 3 && et.Weekday() >= time.Monday && et.Weekday() <= time.Friday && !isUSMarketHoliday(et) {
		return true
	}
	if et.Month() == time.December && et.Day() == 24 && et.Weekday() >= time.Monday && et.Weekday() <= time.Friday && !isUSMarketHoliday(et) {
		return true
	}
	return false
}

func marketSessionET(at time.Time) string {
	loc, err := time.LoadLocation("America/New_York")
	if err != nil {
		loc = time.FixedZone("ET", -5*60*60)
	}
	et := at.In(loc)
	minutes := et.Hour()*60 + et.Minute()
	weekday := et.Weekday()

	if weekday == time.Saturday {
		return "weekend"
	}
	if weekday == time.Sunday {
		if minutes >= 20*60 {

			if isUSMarketHoliday(et.AddDate(0, 0, 1)) {
				return "closed"
			}
			return "overnight"
		}
		return "weekend"
	}
	if weekday == time.Friday && minutes >= 20*60 {
		return "weekend"
	}
	if isUSMarketHoliday(et) {
		return "closed"
	}
	regularClose := 16 * 60
	if isUSEarlyClose(et) {
		regularClose = 13 * 60
	}
	switch {
	case minutes < 4*60:
		return "overnight"
	case minutes < 9*60+30:
		return "pre-market"
	case minutes < regularClose:
		return "regular"
	case minutes < 20*60:
		return "after-hours"
	default:
		return "overnight"
	}
}

func marketSessionBoundaryET(at time.Time) (int64, string) {
	loc, err := time.LoadLocation("America/New_York")
	if err != nil {
		loc = time.FixedZone("ET", -5*60*60)
	}
	et := at.In(loc)
	session := marketSessionET(at)
	boundary := func(day time.Time, hour, minute int) time.Time {
		return time.Date(day.Year(), day.Month(), day.Day(), hour, minute, 0, 0, loc)
	}
	nextTradingDay := func(from time.Time) time.Time {
		d := time.Date(from.Year(), from.Month(), from.Day(), 12, 0, 0, 0, loc)
		for i := 0; i < 10; i++ {
			wd := d.Weekday()
			if wd >= time.Monday && wd <= time.Friday && !isUSMarketHoliday(d) {
				return d
			}
			d = d.AddDate(0, 0, 1)
		}
		return d
	}
	switch session {
	case "pre-market":
		return boundary(et, 9, 30).UnixMilli(), "Opens"
	case "regular":
		hour := 16
		if isUSEarlyClose(et) {
			hour = 13
		}
		return boundary(et, hour, 0).UnixMilli(), "Closes"
	case "after-hours":
		return boundary(et, 20, 0).UnixMilli(), "Ends"
	case "overnight":
		day := et
		if et.Hour() >= 20 {
			day = et.AddDate(0, 0, 1)
		}
		day = nextTradingDay(day)
		return boundary(day, 4, 0).UnixMilli(), "Pre-market"
	case "closed", "weekend":

		day := nextTradingDay(et)
		overnightStart := boundary(day.AddDate(0, 0, -1), 20, 0)
		if !overnightStart.After(et) {
			day = nextTradingDay(day.AddDate(0, 0, 1))
			overnightStart = boundary(day.AddDate(0, 0, -1), 20, 0)
		}
		return overnightStart.UnixMilli(), "Next overnight"
	default:
		return 0, ""
	}
}
