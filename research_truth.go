package main

import (
	"fmt"
	"strings"
	"time"
)

func componentAge(now, ts int64) int64 {
	if ts <= 0 {
		return -1
	}
	age := now - ts
	if age < 0 {
		return 0
	}
	return age
}

// evidenceAge returns both age and timestamp validity. A materially future-dated
// source timestamp is not "age zero" evidence; it is clock-skewed/invalid truth.
func evidenceAge(now, ts int64, allowedFutureSkew time.Duration) (age int64, valid bool, anomaly string) {
	if ts <= 0 {
		return -1, false, "missing timestamp"
	}
	delta := now - ts
	limit := int64(allowedFutureSkew / time.Millisecond)
	if delta < -limit {
		return delta, false, "future timestamp exceeds allowed clock skew"
	}
	if delta < 0 {
		delta = 0
	}
	return delta, true, ""
}

func researchTargetStamp(last map[string]int64, kind, symbol string) int64 {
	return last["research-"+kind+":"+normalizeSymbol(symbol)]
}

func latestNewsForSymbol(news []NewsItem, symbol string) int64 {
	symbol = normalizeSymbol(symbol)
	latest := int64(0)
	for _, n := range news {
		match := false
		for _, s := range n.Symbols {
			if normalizeSymbol(s) == symbol {
				match = true
				break
			}
		}
		if !match && strings.Contains(","+strings.ToUpper(n.Related)+",", ","+symbol+",") {
			match = true
		}
		if match && n.Datetime > 0 {
			t := n.Datetime
			if t < 1_000_000_000_000 {
				t *= 1000
			}
			if t > latest {
				latest = t
			}
		}
	}
	return latest
}

func latestFilingForSymbol(filings []FilingItem, symbol string) int64 {
	symbol = normalizeSymbol(symbol)
	latest := int64(0)
	for _, f := range filings {
		if normalizeSymbol(f.Symbol) != symbol {
			continue
		}

		if strings.TrimSpace(f.FiledAt) != "" {
			if t, err := time.Parse("2006-01-02", strings.TrimSpace(f.FiledAt)); err == nil && t.UnixMilli() > latest {
				latest = t.UnixMilli()
			}
		}
	}
	return latest
}

func minPositiveInt64(values ...int64) int64 {
	min := int64(0)
	for _, v := range values {
		if v > 0 && (min == 0 || v < min) {
			min = v
		}
	}
	return min
}

func maxPositive(values ...int64) int64 {
	max := int64(0)
	for _, v := range values {
		if v > max {
			max = v
		}
	}
	return max
}

func researchCheckState(nowMs, stamp, maxAge int64) (string, string) {
	age, valid, anomaly := evidenceAge(nowMs, stamp, 30*time.Second)
	if !valid {
		if stamp <= 0 {
			return "PARTIAL", "required check is missing"
		}
		return "STALE", "required check timestamp is invalid: " + anomaly
	}
	if age > maxAge {
		return "PARTIAL", fmt.Sprintf("required check is older than %d ms", maxAge)
	}
	return "FRESH", "current"
}

func newsMatchesResearchSymbol(n NewsItem, symbol string) bool {
	symbol = normalizeSymbol(symbol)
	for _, s := range n.Symbols {
		if normalizeSymbol(s) == symbol {
			return true
		}
	}
	return strings.Contains(","+strings.ToUpper(n.Related)+",", ","+symbol+",")
}

// selectedTickerEarningsEvidence resolves only the selected Research symbol and keeps
// scheduled-versus-released semantics explicit. It must not infer a release merely
// because an earnings date exists; that would manufacture catalyst certainty.
func selectedTickerEarningsEvidence(earnings []EarningsItem, symbol string, now time.Time) (string, bool, bool) {
	symbol = normalizeSymbol(symbol)
	loc := easternLocation()
	today := now.In(loc)
	todayDate := today.Format("2006-01-02")
	type candidate struct {
		er    EarningsItem
		when  time.Time
		score int64
	}
	var best *candidate
	for _, er := range earnings {
		if normalizeSymbol(er.Symbol) != symbol || strings.TrimSpace(er.Date) == "" {
			continue
		}
		when, err := time.ParseInLocation("2006-01-02", strings.TrimSpace(er.Date), loc)
		if err != nil {
			continue
		}

		days := int64(when.Sub(time.Date(today.Year(), today.Month(), today.Day(), 0, 0, 0, 0, loc)).Hours() / 24)
		score := int64(2_000_000)
		if days == 0 {
			score = 0
		} else if days > 0 {
			score = 10_000 + days
		} else {
			score = 100_000 + (-days)
		}
		if best == nil || score < best.score {
			c := candidate{er: er, when: when, score: score}
			best = &c
		}
	}
	if best == nil {
		return "No selected-ticker earnings event loaded.", false, false
	}
	er := best.er
	released := earningsReleased(er)
	status := "scheduled"
	if released {
		status = "released"
	}
	detail := fmt.Sprintf("Earnings %s %s · %s", strings.TrimSpace(er.Date), defaultString(strings.TrimSpace(er.Hour), "time not supplied"), status)
	if er.EPSActual != nil || er.EPSEstimate != nil {
		actual, estimate := "—", "—"
		if er.EPSActual != nil {
			actual = fmt.Sprintf("%.4g", *er.EPSActual)
		}
		if er.EPSEstimate != nil {
			estimate = fmt.Sprintf("%.4g", *er.EPSEstimate)
		}
		detail += " · EPS " + actual + " vs " + estimate
	}
	materialNow := strings.TrimSpace(er.Date) == todayDate || (released && best.when.After(today.Add(-48*time.Hour)) && !best.when.After(today.Add(24*time.Hour)))
	return detail, materialNow, released
}

func latestMaterialNewsEvidence(news []NewsItem, symbol string, now time.Time) (string, int64, bool) {
	cutoff := now.Add(-2 * time.Hour).UnixMilli()
	bestAt := int64(0)
	best := ""
	for _, n := range news {
		if !newsMatchesResearchSymbol(n, symbol) || !materialText(n.Headline+" "+n.Summary) {
			continue
		}
		at := n.Datetime
		if at > 0 && at < 1_000_000_000_000 {
			at *= 1000
		}
		if at < cutoff || at > now.Add(30*time.Second).UnixMilli() {
			continue
		}
		if at >= bestAt {
			bestAt = at
			best = "Material news · " + defaultString(strings.TrimSpace(n.Headline), "headline unavailable") + " · " + defaultString(strings.TrimSpace(n.Source), "source unavailable")
		}
	}
	return best, bestAt, best != ""
}

func latestMaterialSECEvidence(filings []FilingItem, symbol string, now time.Time) (string, int64, bool) {
	symbol = normalizeSymbol(symbol)
	bestAt := int64(0)
	best := ""
	for _, f := range filings {
		if normalizeSymbol(f.Symbol) != symbol || !materialSECFilingForTradingRisk(f) || strings.TrimSpace(f.FiledAt) == "" {
			continue
		}
		t, err := time.ParseInLocation("2006-01-02", strings.TrimSpace(f.FiledAt), easternLocation())
		if err != nil {
			continue
		}
		at := t.UnixMilli()
		if now.Sub(t) < 0 || now.Sub(t) > 48*time.Hour {
			continue
		}
		if at >= bestAt {
			bestAt = at
			best = "Material SEC · " + defaultString(strings.TrimSpace(f.Form), "filing") + " · " + defaultString(strings.TrimSpace(f.Meaning), strings.TrimSpace(f.Description))
		}
	}
	return best, bestAt, best != ""
}

// catalystMaterialEventComponent is the canonical selected-ticker event boundary for
// Research freshness. A "no active catalyst" result is Fresh only when Catalyst
// Watch plus News, Earnings, and SEC checks are all current and internally consistent.
func catalystMaterialEventComponent(symbol string, news []NewsItem, earnings []EarningsItem, filings []FilingItem, last map[string]int64, health map[string]string, reactions map[string]CatalystReactionState, now time.Time) ResearchEvidenceComponent {
	nowMs := now.UnixMilli()
	ncheck := researchTargetStamp(last, "news", symbol)
	echeck := researchTargetStamp(last, "earnings", symbol)
	scheck := researchTargetStamp(last, "sec", symbol)
	wcheck := last["catalyst-watch"]
	checks := []struct {
		name         string
		stamp, limit int64
	}{
		{"News", ncheck, int64((30 * time.Minute) / time.Millisecond)},
		{"Earnings", echeck, int64((4 * time.Hour) / time.Millisecond)},
		{"SEC", scheck, int64((60 * time.Minute) / time.Millisecond)},
		{"Catalyst Watch", wcheck, int64((60 * time.Minute) / time.Millisecond)},
	}
	state := "FRESH"
	reasons := []string{}
	for _, c := range checks {
		cs, why := researchCheckState(nowMs, c.stamp, c.limit)
		if cs == "STALE" {
			state = "STALE"
		} else if cs == "PARTIAL" && state == "FRESH" {
			state = "PARTIAL"
		}
		if cs != "FRESH" {
			reasons = append(reasons, c.name+" "+why)
		}
	}
	for _, key := range []string{"research-sec:" + symbol, "catalyst-watch"} {
		h := strings.ToLower(strings.TrimSpace(health[key]))
		if strings.Contains(h, "error") || strings.Contains(h, "failed") || strings.Contains(h, "unavailable") || strings.Contains(h, "degraded") {
			if state == "FRESH" {
				state = "PARTIAL"
			}
			reasons = append(reasons, key+" is degraded")
		}
	}

	earningsDetail, earningsMaterialNow, earningsReleasedNow := selectedTickerEarningsEvidence(earnings, symbol, now)
	newsDetail, materialNewsAt, hasMaterialNews := latestMaterialNewsEvidence(news, symbol, now)
	secDetail, materialSECAt, hasMaterialSEC := latestMaterialSECEvidence(filings, symbol, now)
	materialPieces := []string{}
	if earningsMaterialNow {
		materialPieces = append(materialPieces, earningsDetail)
	}
	if hasMaterialNews {
		materialPieces = append(materialPieces, newsDetail)
	}
	if hasMaterialSEC {
		materialPieces = append(materialPieces, secDetail)
	}

	r, active := reactions[symbol]
	dataAt := maxPositive(latestNewsForSymbol(news, symbol), latestFilingForSymbol(filings, symbol), materialNewsAt, materialSECAt)
	detail := "No active catalyst is reported; Catalyst Watch plus required selected-ticker News, Earnings and SEC checks are current and valid. " + earningsDetail
	source := "Catalyst Watch · News/Earnings/SEC"
	if active && strings.TrimSpace(r.State) != "" {
		triggerAt := r.TriggerAt
		updatedAt := r.UpdatedAt
		if _, valid, anomaly := evidenceAge(nowMs, maxPositive(triggerAt, updatedAt), 30*time.Second); !valid {
			state = "STALE"
			reasons = append(reasons, "Catalyst timestamp is invalid: "+anomaly)
		}
		dataAt = maxPositive(dataAt, triggerAt, updatedAt)
		detail = fmt.Sprintf("Active catalyst: %s · %s · %s. %s", defaultString(r.TriggerType, "material event"), defaultString(r.State, "TRIGGERED"), defaultString(r.Phase, "active"), strings.TrimSpace(r.Detail))
		if len(materialPieces) > 0 {
			detail += " Evidence: " + strings.Join(materialPieces, " · ") + "."
		}
	} else if len(materialPieces) > 0 {

		detail = "Material event context: " + strings.Join(materialPieces, " · ") + "."
		if earningsReleasedNow || hasMaterialNews || hasMaterialSEC {
			if state == "FRESH" {
				state = "PARTIAL"
			}
			reasons = append(reasons, "material evidence exists without a matching Catalyst Watch reaction state")
		} else {
			detail += " Scheduled earnings risk is armed; no release/reaction is recorded yet."
		}
	}
	if len(reasons) > 0 {
		detail += " Required-evidence issue: " + strings.Join(reasons, "; ") + "."
	}
	return ResearchEvidenceComponent{Dataset: "Catalyst & Material Event Context", Required: true, State: state, Source: source, CheckAt: minPositiveInt64(ncheck, echeck, scheck, wcheck), DataAt: dataAt, Detail: detail}
}

// requiredMarketContextComponent prevents ticker Research from appearing fully Fresh
// when the broad decision context is stale or missing. Session-aware limits are
// intentional: closed/weekend evidence is judged differently from live-session data.
func requiredMarketContextComponent(quotes map[string]Quote, global GlobalMarketContext, now time.Time) ResearchEvidenceComponent {
	nowMs := now.UnixMilli()
	session := marketSessionET(now)
	limit := int64((5 * time.Minute) / time.Millisecond)
	if session == "overnight" {
		limit = int64((90 * time.Minute) / time.Millisecond)
	}
	if session == "closed" || session == "weekend" {
		limit = int64((96 * time.Hour) / time.Millisecond)
	}
	state := "FRESH"
	reasons := []string{}
	sources := []string{}
	evidence := []string{}
	checks := []int64{}
	dataTimes := []int64{}
	for _, sym := range []string{"SPY", "QQQ", "VIX"} {
		q := quotes[sym]
		dataAt := q.ProviderTimestamp
		if dataAt <= 0 {
			dataAt = q.UpdatedAt
		}
		sources = append(sources, sym+"="+providerFromQuoteSource(q.Source))
		evidence = append(evidence, fmt.Sprintf("%s %.4g (%s)", sym, q.Price, defaultString(q.DataState, q.FeedType)))
		checks = append(checks, q.UpdatedAt)
		dataTimes = append(dataTimes, dataAt)
		if q.Price <= 0 || dataAt <= 0 {
			if state == "FRESH" {
				state = "PARTIAL"
			}
			reasons = append(reasons, sym+" unavailable")
			continue
		}
		providerAge, receiptAge, valid, why := quoteEvidenceTimestampTruth(q, nowMs)
		if !valid {
			state = "STALE"
			reasons = append(reasons, sym+" "+why)
		} else if providerAge > limit || receiptAge > limit {
			state = "STALE"
			reasons = append(reasons, fmt.Sprintf("%s exceeds session-aware freshness (%d/%d ms)", sym, providerAge, receiptAge))
		}
		if strings.Contains(strings.ToLower(q.DataState+" "+q.FeedType), "delay") {
			reasons = append(reasons, sym+" is delayed (truthfully labeled)")
		}
	}
	checks = append(checks, global.UpdatedAt)
	dataTimes = append(dataTimes, global.UpdatedAt)
	if global.UpdatedAt <= 0 || strings.TrimSpace(global.Mode) == "" {
		if state == "FRESH" {
			state = "PARTIAL"
		}
		reasons = append(reasons, "Global Context unavailable")
	} else if age, valid, anomaly := evidenceAge(nowMs, global.UpdatedAt, 30*time.Second); !valid {
		state = "STALE"
		reasons = append(reasons, "Global Context timestamp invalid: "+anomaly)
	} else if age > limit {
		state = "STALE"
		reasons = append(reasons, fmt.Sprintf("Global Context exceeds session-aware freshness (%d ms)", age))
	}
	sources = append(sources, "Global="+defaultString(global.Mode, "unknown"))
	evidence = append(evidence, "Global "+defaultString(global.Tone, "UNKNOWN")+" / "+defaultString(global.Mode, "unknown"))
	detail := "Required broad market context is current: " + strings.Join(evidence, " · ") + "."
	if len(reasons) > 0 {
		detail = "Required market context is degraded: " + strings.Join(reasons, "; ") + ". Evidence: " + strings.Join(evidence, " · ") + "."
	}
	return ResearchEvidenceComponent{Dataset: "Required Market Context", Required: true, State: state, Source: strings.Join(sources, " · "), CheckAt: minPositiveInt64(checks...), DataAt: minPositiveInt64(dataTimes...), Detail: detail}
}

// buildResearchPackageTruth owns page-level Research truth. Required components use
// worst-dependency semantics: a recent check cannot hide stale underlying evidence,
// and optional context must never upgrade a degraded required component.
func buildResearchPackageTruth(st AppState, quotes map[string]Quote, bars map[string]map[string][]Bar, fundamentals map[string]FundamentalSnapshot, news []NewsItem, earnings []EarningsItem, filings []FilingItem, last map[string]int64, health map[string]string, reconciliations []ProviderReconciliationDecision, now time.Time, contexts ...ResearchPackageContext) ResearchPackageTruth {
	symbol := normalizeSymbol(st.UI.SelectedTicker)
	nowMs := now.UnixMilli()
	result := ResearchPackageTruth{Symbol: symbol, State: "BLOCKED", GeneratedAt: nowMs, Components: []ResearchEvidenceComponent{}}
	if symbol == "" || symbol == "VIX" {
		result.BlockingReasons = []string{"No valid equity/ETF Research target selected"}
		return result
	}
	ctx := ResearchPackageContext{}
	if len(contexts) > 0 {
		ctx = contexts[0]
	}
	add := func(c ResearchEvidenceComponent) {
		c.Symbol = symbol
		if age, valid, _ := evidenceAge(nowMs, c.CheckAt, 30*time.Second); valid {
			c.CheckAgeMs = age
		} else {
			c.CheckAgeMs = -1
		}
		if age, valid, _ := evidenceAge(nowMs, c.DataAt, 30*time.Second); valid {
			c.DataAgeMs = age
		} else {
			c.DataAgeMs = -1
		}
		result.Components = append(result.Components, c)
		if c.State == "BLOCKED" || c.State == "STALE" {
			result.BlockingReasons = append(result.BlockingReasons, c.Dataset+" "+c.State)
		}
	}
	session := marketSessionET(now)
	q := quotes[symbol]
	qData := q.ProviderTimestamp
	if qData <= 0 {
		qData = q.UpdatedAt
	}
	qLimit := int64((3 * time.Minute) / time.Millisecond)
	if session == "overnight" {
		qLimit = int64((90 * time.Minute) / time.Millisecond)
	}
	if session == "closed" || session == "weekend" {
		qLimit = int64((96 * time.Hour) / time.Millisecond)
	}
	qs := "FRESH"
	qdetail := "Selected-ticker quote is current."
	if q.Price <= 0 || qData <= 0 {
		qs = "BLOCKED"
		qdetail = "Selected-ticker quote is unavailable."
	} else {
		providerAge, receiptAge, validTimestamp, timestampDetail := quoteEvidenceTimestampTruth(q, nowMs)
		if !validTimestamp {
			qs = "STALE"
			qdetail = timestampDetail
		} else if providerAge > qLimit || receiptAge > qLimit {
			qs = "STALE"
			qdetail = fmt.Sprintf("Selected-ticker quote exceeds the Research session-aware freshness limit (provider/data age %d ms; receipt age %d ms).", providerAge, receiptAge)
		}
	}
	add(ResearchEvidenceComponent{Dataset: "Quote", Required: true, Critical: true, State: qs, Source: providerFromQuoteSource(q.Source), CheckAt: q.UpdatedAt, DataAt: qData, Detail: qdetail})

	daily := bars[symbol]["daily"]
	hs := "FRESH"
	hdetail := "Canonical adjusted daily history is available."
	hdata := int64(0)
	if len(daily) > 0 {
		hdata = daily[len(daily)-1].T * 1000
	}
	hcheck := researchTargetStamp(last, "history", symbol)
	if len(daily) == 0 {
		hs = "BLOCKED"
		hdetail = "Selected-ticker adjusted daily history is unavailable."
	} else {
		maxData := int64((120 * time.Hour) / time.Millisecond)
		if session == "regular" || session == "pre-market" || session == "after-hours" {
			maxData = int64((96 * time.Hour) / time.Millisecond)
		}
		if componentAge(nowMs, hdata) > maxData {
			hs = "STALE"
			hdetail = "Latest selected-ticker daily bar is too old for Research."
		} else if hcheck == 0 || componentAge(nowMs, hcheck) > int64((48*time.Hour)/time.Millisecond) {
			hs = "PARTIAL"
			hdetail = "Daily bars exist, but this selected ticker's history route has not been reconciled recently."
		}
	}
	add(ResearchEvidenceComponent{Dataset: "Daily History", Required: true, Critical: true, State: hs, Source: "Provider Router · adjustment=all", CheckAt: hcheck, DataAt: hdata, Detail: hdetail})

	f := fundamentals[symbol]
	fcheck := researchTargetStamp(last, "fundamentals", symbol)
	fs := "FRESH"
	fdetail := "Selected-ticker fundamentals have current data and target reconciliation evidence."
	if !fundamentalSnapshotUsable(f) {
		fs = "PARTIAL"
		fdetail = "Selected-ticker fundamentals are unavailable/incomplete."
	} else {
		fAge, fValid, fAnomaly := evidenceAge(nowMs, f.UpdatedAt, 30*time.Second)
		fDataLimit := int64((72 * time.Hour) / time.Millisecond)
		if session == "closed" || session == "weekend" {
			fDataLimit = int64((120 * time.Hour) / time.Millisecond)
		}
		if !fValid {
			fs = "STALE"
			fdetail = "Fundamentals timestamp is invalid: " + fAnomaly + "."
		} else if fAge > fDataLimit {
			fs = "STALE"
			fdetail = fmt.Sprintf("Selected-ticker fundamental snapshot is too old (%d ms) even though the route may have been checked recently.", fAge)
		} else if fcheck == 0 {
			fs = "PARTIAL"
			fdetail = "Fundamentals exist, but selected-ticker target reconciliation has not been proven."
		} else if checkAge, ok, anomaly := evidenceAge(nowMs, fcheck, 30*time.Second); !ok {
			fs = "PARTIAL"
			fdetail = "Fundamentals reconciliation timestamp is invalid: " + anomaly + "."
		} else if checkAge > int64((36*time.Hour)/time.Millisecond) {
			fs = "PARTIAL"
			fdetail = "Fundamentals data is usable, but selected-ticker target reconciliation is not current."
		}
	}
	add(ResearchEvidenceComponent{Dataset: "Fundamentals", Required: true, State: fs, Source: f.Source, CheckAt: fcheck, DataAt: f.UpdatedAt, Detail: fdetail})

	ncheck := researchTargetStamp(last, "news", symbol)
	ns := "FRESH"
	ndetail := "Selected-ticker news route was checked recently; sparse/no-news results remain valid checks."
	if ncheck == 0 || componentAge(nowMs, ncheck) > int64((30*time.Minute)/time.Millisecond) {
		ns = "PARTIAL"
		ndetail = "Selected-ticker news has not been target-reconciled recently."
	}
	add(ResearchEvidenceComponent{Dataset: "News", Required: true, State: ns, Source: "Finnhub → Marketaux", CheckAt: ncheck, DataAt: latestNewsForSymbol(news, symbol), Detail: ndetail})

	echeck := researchTargetStamp(last, "earnings", symbol)
	es := "FRESH"
	edetail := "Selected-ticker earnings route was target-reconciled."
	if echeck == 0 || componentAge(nowMs, echeck) > int64((4*time.Hour)/time.Millisecond) {
		es = "PARTIAL"
		edetail = "Selected-ticker earnings has not been target-reconciled recently."
	}
	add(ResearchEvidenceComponent{Dataset: "Earnings", Required: true, State: es, Source: "Finnhub → yfinance", CheckAt: echeck, DataAt: 0, Detail: edetail})

	scheck := researchTargetStamp(last, "sec", symbol)
	ss := "FRESH"
	sdetail := "Selected-ticker SEC/Ownership route was target-reconciled."
	shealth := strings.ToLower(health["research-sec:"+symbol])
	if strings.Contains(shealth, "error") || strings.Contains(shealth, "failed") || strings.Contains(shealth, "unavailable") || strings.Contains(shealth, "degraded") {
		ss = "PARTIAL"
		sdetail = "Selected-ticker SEC/Ownership reconciliation is degraded."
	} else if scheck == 0 || componentAge(nowMs, scheck) > int64((60*time.Minute)/time.Millisecond) {
		ss = "PARTIAL"
		sdetail = "Selected-ticker SEC/Ownership has not been target-reconciled recently."
	}
	add(ResearchEvidenceComponent{Dataset: "SEC & Ownership", Required: true, State: ss, Source: "SEC EDGAR", CheckAt: scheck, DataAt: latestFilingForSymbol(filings, symbol), Detail: sdetail})

	if len(contexts) > 0 {
		add(catalystMaterialEventComponent(symbol, news, earnings, filings, last, health, ctx.CatalystReactions, now))
		add(requiredMarketContextComponent(quotes, ctx.Global, now))
	}

	ps := "PARTIAL"
	pdetail := "No current provider reconciliation evidence exists for the selected ticker; no agreement claim is made."
	pcheck, pdata := int64(0), int64(0)
	foundRec := false
	for _, r := range reconciliations {
		if r.Symbol != symbol {
			continue
		}
		foundRec = true
		pcheck = r.UpdatedAt
		pdetail = r.Reason
		for _, o := range r.Observations {
			t := o.ProviderTimestamp
			if t <= 0 {
				t = o.ReceivedAt
			}
			if t > pdata {
				pdata = t
			}
		}
		switch r.State {
		case "AGREED":
			ps = "FRESH"
		case "SINGLE SOURCE":
			ps = "PARTIAL"
		case "STALE":
			ps = "STALE"
		case "CONFLICT":
			ps = "PARTIAL"
			if r.DifferencePct >= 1.0 {
				ps = "BLOCKED"
			}
		default:
			ps = "PARTIAL"
		}
		break
	}
	if !foundRec {
		pcheck, pdata = 0, 0
	}
	add(ResearchEvidenceComponent{Dataset: "Provider Reconciliation", Required: true, Critical: true, State: ps, Source: "Provider Router", CheckAt: pcheck, DataAt: pdata, Detail: pdetail})

	state := "FRESH"
	hasPartial := false
	hasStale := false
	hasBlocked := false
	for _, c := range result.Components {
		switch c.State {
		case "BLOCKED":
			hasBlocked = true
		case "STALE":
			hasStale = true
		case "PARTIAL":
			hasPartial = true
		}
	}
	if hasBlocked {
		state = "BLOCKED"
	} else if hasStale {
		state = "STALE"
	} else if hasPartial {
		state = "PARTIAL"
	}
	result.State = state
	return result
}
