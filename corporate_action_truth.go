package main

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"sort"
	"strings"
	"time"
)

func lifecycleSymbol(v string) string {
	s := normalizeSymbol(v)
	switch s {
	case "", "NIL", "NULL", "NONE", "NA", "N-A", "UNKNOWN", "UNDEFINED":
		return ""
	}
	return s
}

func lifecycleActionKind(v string) string {
	s := strings.ToLower(strings.TrimSpace(v))
	s = strings.NewReplacer("-", "_", " ", "_").Replace(s)
	return s
}

func lifecycleActionEffective(status string) bool {
	s := strings.ToUpper(strings.TrimSpace(status))
	return s == "EFFECTIVE" || s == "COMPLETE" || s == "COMPLETED" || s == "PROCESSED" || s == "CLOSED"
}

func lifecycleRank(state string) int {
	switch state {
	case "DELISTED":
		return 90
	case "DELISTING PENDING":
		return 80
	case "MERGED":
		return 70
	case "MERGER PENDING":
		return 60
	case "NAME OR TICKER CHANGE":
		return 50
	case "NEW LISTING / IPO":
		return 40
	case "CORPORATE ACTION PENDING":
		return 30
	case "ACTIVE":
		return 20
	default:
		return 10
	}
}

func lifecycleStateFromAction(a CorporateAction) (string, bool) {
	k := lifecycleActionKind(a.Type)
	effective := lifecycleActionEffective(a.Status)
	switch {
	case strings.Contains(k, "delist") || strings.Contains(k, "worthless") || strings.Contains(k, "removal"):
		if effective {
			return "DELISTED", true
		}
		return "DELISTING PENDING", true
	case strings.Contains(k, "merger") || strings.Contains(k, "acquisition"):
		if effective {
			return "MERGED", true
		}
		return "MERGER PENDING", true
	case strings.Contains(k, "name_change") || strings.Contains(k, "ticker_change") || strings.Contains(k, "symbol_change"):
		return "NAME OR TICKER CHANGE", true
	case strings.Contains(k, "ipo") || strings.Contains(k, "new_listing") || k == "listing":
		return "NEW LISTING / IPO", true
	case !effective && strings.TrimSpace(a.Status) != "":
		return "CORPORATE ACTION PENDING", true
	}
	return "", false
}

func lifecycleEffectiveDate(a CorporateAction) string {
	for _, v := range []string{a.ExDate, a.ProcessDate, a.RecordDate, a.PayableDate} {
		if strings.TrimSpace(v) != "" {
			return strings.TrimSpace(v)
		}
	}
	return ""
}

func lifecycleDateMillis(v string) int64 {
	v = strings.TrimSpace(v)
	if v == "" {
		return 0
	}
	for _, layout := range []string{"2006-01-02", time.RFC3339, "2006-01-02T15:04:05"} {
		if t, err := time.Parse(layout, v); err == nil {
			return t.UnixMilli()
		}
	}
	return 0
}

func lifecycleStateOrderTime(s SymbolLifecycleState) int64 {
	if t := lifecycleDateMillis(s.EffectiveDate); t > 0 {
		return t
	}
	if t := lifecycleDateMillis(s.ProcessDate); t > 0 {
		return t
	}
	return 0
}

func lifecycleCandidateSupersedes(current, candidate SymbolLifecycleState) bool {
	curEvent, candEvent := lifecycleStateOrderTime(current), lifecycleStateOrderTime(candidate)
	if curEvent > 0 && candEvent > 0 && curEvent != candEvent {
		return candEvent > curEvent
	}
	if candEvent > 0 && curEvent == 0 {
		return true
	}
	if lifecycleRank(candidate.State) != lifecycleRank(current.State) {
		return lifecycleRank(candidate.State) > lifecycleRank(current.State)
	}
	return candidate.UpdatedAt >= current.UpdatedAt
}

// buildCorporateActionTruth is the single lifecycle/history authority. Newer effective
// evidence may supersede an older severe event (for example, a legitimate relisting),
// while raw-vs-adjusted history coverage remains independent from lifecycle labels.
func buildCorporateActionTruth(actions []CorporateAction, bars map[string]map[string][]Bar, now int64, coverageMaps ...map[string]RawHistoryCoverage) CorporateActionTruth {
	coverage := map[string]RawHistoryCoverage{}
	if len(coverageMaps) > 0 && coverageMaps[0] != nil {
		coverage = clone(coverageMaps[0])
	}
	out := CorporateActionTruth{Actions: clone(actions), SymbolLineage: map[string]string{}, Lifecycle: map[string]SymbolLifecycleState{}, RawHistoryAvailable: map[string]bool{}, RawHistoryCoverage: coverage, CanonicalHistory: "Adjusted OHLCV (Alpaca adjustment=all when Alpaca is active); raw comparison bars are retained as daily-raw with explicit completeness provenance when entitled.", UpdatedAt: now}
	seen := map[string]bool{}

	for rawSym, series := range bars {
		sym := lifecycleSymbol(rawSym)
		if sym == "" {
			continue
		}
		state := SymbolLifecycleState{Symbol: sym, State: "UNKNOWN", Source: "Canonical adjusted daily history", UpdatedAt: now, Reason: "Current listing status is not proven by lifecycle evidence."}
		if daily := series["daily"]; len(daily) > 0 {
			lastAt := daily[len(daily)-1].T * 1000
			if lastAt > 0 && now >= lastAt && now-lastAt <= int64((10*24*time.Hour)/time.Millisecond) {
				state.State = "ACTIVE"
				state.EffectiveDate = time.UnixMilli(lastAt).UTC().Format("2006-01-02")
				state.Reason = "Recent canonical adjusted daily history supports an active listing; no stronger lifecycle event overrides it."
			}
		}
		out.Lifecycle[sym] = state
	}
	for _, a := range actions {
		syms := uniqueSymbols([]string{lifecycleSymbol(a.Symbol), lifecycleSymbol(a.OldSymbol), lifecycleSymbol(a.NewSymbol)})
		for _, sym := range syms {
			if sym == "" {
				continue
			}
			if !seen[sym] {
				seen[sym] = true
				out.AffectedSymbols = append(out.AffectedSymbols, sym)
			}
		}
		oldS, newS := lifecycleSymbol(a.OldSymbol), lifecycleSymbol(a.NewSymbol)
		if strings.Contains(lifecycleActionKind(a.Type), "name_change") && oldS != "" && newS != "" && oldS != newS {
			out.SymbolLineage[oldS] = newS
		}
		state, material := lifecycleStateFromAction(a)
		if !material {
			continue
		}
		reason := fmt.Sprintf("%s evidence from %s", defaultString(a.Type, "corporate action"), defaultString(a.Source, "provider"))
		if oldS != "" && newS != "" && oldS != newS {
			reason += " · " + oldS + " → " + newS
		}
		for _, sym := range syms {
			if sym == "" {
				continue
			}
			candidate := SymbolLifecycleState{Symbol: sym, State: state, Source: defaultString(a.Source, "Corporate action provider"), EvidenceID: a.ID, EffectiveDate: lifecycleEffectiveDate(a), ProcessDate: a.ProcessDate, FirstSeenAt: a.FirstSeenAt, UpdatedAt: a.UpdatedAt, Reason: reason}
			if candidate.UpdatedAt <= 0 {
				candidate.UpdatedAt = now
			}
			current, ok := out.Lifecycle[sym]
			if !ok || lifecycleCandidateSupersedes(current, candidate) {
				out.Lifecycle[sym] = candidate
			}
		}
	}
	sort.Strings(out.AffectedSymbols)

	for _, sym := range out.AffectedSymbols {
		cov, ok := out.RawHistoryCoverage[sym]
		if !ok {
			rows := bars[sym]["daily-raw"]
			if len(rows) > 0 {
				cov = RawHistoryCoverage{Symbol: sym, State: "PARTIAL", BarCount: len(rows), PaginationComplete: false, Detail: "Raw bars exist but completeness provenance is unavailable."}
			} else {
				cov = RawHistoryCoverage{Symbol: sym, State: "UNAVAILABLE", Detail: "Raw comparison history has not been reconciled."}
			}
			out.RawHistoryCoverage[sym] = cov
		}
		out.RawHistoryAvailable[sym] = strings.EqualFold(cov.State, "COMPLETE") && cov.PaginationComplete && cov.BarCount > 0
	}
	return out
}

// buildEvidenceSnapshot fingerprints material decision evidence, not wall-clock aging.
// Snapshot IDs change when lifecycle, catalyst, required-market-context, provider, or
// history truth changes; ordinary passage of time alone must not churn the identity.
func buildEvidenceSnapshot(pkg ResearchPackageTruth, reconciliations []ProviderReconciliationDecision, corp CorporateActionTruth) EvidenceSnapshot {
	type componentFingerprint struct {
		Dataset, State, Source, Detail string
		Required, Critical             bool
		CheckAt, DataAt                int64
	}
	type obsFingerprint struct {
		Provider                      string
		Price, Bid, Ask               float64
		ProviderTimestamp, ReceivedAt int64
		Source, FeedType, DataState   string
	}
	type recFingerprint struct {
		Dataset, Symbol, State, CanonicalProvider, Reason string
		CanonicalValue, DifferencePct                     float64
		Observations                                      []obsFingerprint
	}
	type actionFingerprint struct {
		ID, Symbol, Type, ProcessDate, ExDate, RecordDate, PayableDate, OldSymbol, NewSymbol, Status, Detail, Source string
		Ratio, CashAmount, AdjustmentFactor                                                                          float64
	}
	type lifecycleFingerprint struct {
		Symbol, State, Source, EvidenceID, EffectiveDate, ProcessDate, Reason string
		FirstSeenAt, UpdatedAt                                                int64
	}
	components := make([]componentFingerprint, 0, len(pkg.Components))
	for _, c := range pkg.Components {
		components = append(components, componentFingerprint{c.Dataset, c.State, c.Source, c.Detail, c.Required, c.Critical, c.CheckAt, c.DataAt})
	}
	rel := []ProviderReconciliationDecision{}
	recfp := []recFingerprint{}
	for _, r := range reconciliations {
		if r.Symbol != pkg.Symbol {
			continue
		}
		rel = append(rel, r)
		x := recFingerprint{Dataset: r.Dataset, Symbol: r.Symbol, State: r.State, CanonicalProvider: r.CanonicalProvider, Reason: r.Reason, CanonicalValue: r.CanonicalValue, DifferencePct: r.DifferencePct}
		for _, o := range r.Observations {
			x.Observations = append(x.Observations, obsFingerprint{o.Provider, o.Price, o.Bid, o.Ask, o.ProviderTimestamp, o.ReceivedAt, o.Source, o.FeedType, o.DataState})
		}
		recfp = append(recfp, x)
	}
	acts := []CorporateAction{}
	actfp := []actionFingerprint{}
	for _, a := range corp.Actions {
		if normalizeSymbol(a.Symbol) != pkg.Symbol && normalizeSymbol(a.OldSymbol) != pkg.Symbol && normalizeSymbol(a.NewSymbol) != pkg.Symbol {
			continue
		}
		acts = append(acts, a)
		actfp = append(actfp, actionFingerprint{a.ID, a.Symbol, a.Type, a.ProcessDate, a.ExDate, a.RecordDate, a.PayableDate, a.OldSymbol, a.NewSymbol, a.Status, a.Detail, a.Source, a.Ratio, a.CashAmount, a.AdjustmentFactor})
	}
	coverage := []RawHistoryCoverage{}
	if c, ok := corp.RawHistoryCoverage[pkg.Symbol]; ok {
		coverage = append(coverage, c)
	}
	var lifecycle *SymbolLifecycleState
	var lifecycleFP *lifecycleFingerprint
	if l, ok := corp.Lifecycle[pkg.Symbol]; ok {
		lc := l
		lifecycle = &lc
		fp := lifecycleFingerprint{l.Symbol, l.State, l.Source, l.EvidenceID, l.EffectiveDate, l.ProcessDate, l.Reason, l.FirstSeenAt, l.UpdatedAt}
		lifecycleFP = &fp
	}
	semantic := struct {
		Symbol, State  string
		Components     []componentFingerprint
		Reconciliation []recFingerprint
		Actions        []actionFingerprint
		Coverage       []RawHistoryCoverage
		Lifecycle      *lifecycleFingerprint
	}{pkg.Symbol, pkg.State, components, recfp, actfp, coverage, lifecycleFP}
	raw, _ := json.Marshal(semantic)
	sum := sha256.Sum256(raw)
	id := hex.EncodeToString(sum[:12])
	return EvidenceSnapshot{ID: id, Symbol: pkg.Symbol, GeneratedAt: pkg.GeneratedAt, ResearchState: pkg.State, Components: clone(pkg.Components), ProviderReconciliation: rel, CorporateActions: acts, CorporateHistoryCoverage: clone(coverage), SymbolLifecycle: lifecycle}
}

func corporateActionKey(a CorporateAction) string {
	if id := strings.TrimSpace(a.ID); id != "" {
		return "id:" + id
	}
	parts := []string{strings.ToLower(strings.TrimSpace(a.Type)), normalizeSymbol(a.Symbol), normalizeSymbol(a.OldSymbol), normalizeSymbol(a.NewSymbol), strings.TrimSpace(a.ProcessDate), strings.TrimSpace(a.ExDate), strings.TrimSpace(a.RecordDate), strings.TrimSpace(a.PayableDate), fmt.Sprintf("%.8g", a.Ratio), fmt.Sprintf("%.8g", a.CashAmount)}
	return strings.Join(parts, "|")
}

func mergeCorporateActionLedger(existing, fresh []CorporateAction, now int64) []CorporateAction {
	if now <= 0 {
		now = time.Now().UnixMilli()
	}
	m := map[string]CorporateAction{}
	for _, a := range existing {
		k := corporateActionKey(a)
		if k == "" {
			continue
		}
		if a.FirstSeenAt <= 0 {
			if a.UpdatedAt > 0 {
				a.FirstSeenAt = a.UpdatedAt
			} else {
				a.FirstSeenAt = now
			}
		}
		m[k] = a
	}
	for _, a := range fresh {
		k := corporateActionKey(a)
		if k == "" {
			continue
		}
		if old, ok := m[k]; ok {
			a.FirstSeenAt = old.FirstSeenAt
			if a.FirstSeenAt <= 0 {
				a.FirstSeenAt = now
			}
		} else if a.FirstSeenAt <= 0 {
			a.FirstSeenAt = now
		}
		if a.UpdatedAt <= 0 {
			a.UpdatedAt = now
		}
		m[k] = a
	}
	out := make([]CorporateAction, 0, len(m))
	for _, a := range m {
		out = append(out, a)
	}
	sort.Slice(out, func(i, j int) bool {
		di := out[i].ProcessDate
		if di == "" {
			di = out[i].ExDate
		}
		dj := out[j].ProcessDate
		if dj == "" {
			dj = out[j].ExDate
		}
		if di != dj {
			return di < dj
		}
		if out[i].Symbol != out[j].Symbol {
			return out[i].Symbol < out[j].Symbol
		}
		return out[i].Type < out[j].Type
	})
	return out
}

func (e *Engine) refreshResearchHistory(ctx context.Context, symbol string) bool {
	symbol = normalizeSymbol(symbol)
	if symbol == "" || symbol == "VIX" {
		return false
	}
	dailyOK := e.refreshHistoryRoutedMode(ctx, []string{symbol}, "daily")

	_ = e.refreshHistoryRoutedMode(ctx, []string{symbol}, "intraday")
	if dailyOK {
		now := time.Now().UnixMilli()
		e.mu.Lock()
		e.lastUpdated["research-history:"+symbol] = now
		e.health["research-history:"+symbol] = "healthy · selected-ticker daily history reconciled"
		e.mu.Unlock()
	}
	return dailyOK
}

func mergeBarsByTime(existing, fresh []Bar) []Bar {
	m := map[int64]Bar{}
	for _, b := range existing {
		if b.T > 0 && b.C > 0 {
			m[b.T] = b
		}
	}
	for _, b := range fresh {
		if b.T > 0 && b.C > 0 {
			m[b.T] = b
		}
	}
	out := make([]Bar, 0, len(m))
	for _, b := range m {
		out = append(out, b)
	}
	sort.Slice(out, func(i, j int) bool { return out[i].T < out[j].T })
	return out
}

// refreshAlpacaRawHistoryForCorporateActions backfills unadjusted history required to
// preserve corporate-action lineage. Pagination or transport failure remains PARTIAL;
// it must never silently mark raw history COMPLETE.
func (e *Engine) refreshAlpacaRawHistoryForCorporateActions(ctx context.Context, key, secret string, actions []CorporateAction) int {
	if strings.TrimSpace(key) == "" || strings.TrimSpace(secret) == "" {
		return 0
	}
	set := map[string]bool{}
	for _, a := range actions {
		for _, s := range []string{a.Symbol, a.OldSymbol, a.NewSymbol} {
			s = normalizeSymbol(s)
			if s != "" && s != "VIX" {
				set[s] = true
			}
		}
	}
	symbols := make([]string, 0, len(set))
	for s := range set {
		symbols = append(symbols, s)
	}
	sort.Strings(symbols)
	if len(symbols) == 0 {
		return 0
	}
	client := &http.Client{Timeout: 40 * time.Second}
	loaded := 0
	startAt := time.Now().AddDate(-15, 0, 0).UTC().Format(time.RFC3339)
	endAt := time.Now().UTC().Format(time.RFC3339)
	headers := map[string]string{"APCA-API-KEY-ID": key, "APCA-API-SECRET-KEY": secret}
	anyFailure := false
	for i := 0; i < len(symbols); i += 50 {
		j := i + 50
		if j > len(symbols) {
			j = len(symbols)
		}
		batch := symbols[i:j]
		acc := map[string][]Bar{}
		pageToken := ""
		seenTokens := map[string]bool{}
		pages := 0
		complete := true
		failDetail := ""
		for {
			if pages >= 100 {
				complete = false
				failDetail = "pagination safety stop"
				anyFailure = true
				break
			}
			q := url.Values{}
			q.Set("symbols", strings.Join(batch, ","))
			q.Set("timeframe", "1Day")
			q.Set("start", startAt)
			q.Set("end", endAt)
			q.Set("limit", "10000")
			q.Set("adjustment", "raw")
			q.Set("feed", "iex")
			q.Set("sort", "asc")
			if pageToken != "" {
				q.Set("page_token", pageToken)
			}
			raw := strings.TrimRight(alpacaDataBaseURL, "/") + "/v2/stocks/bars?" + q.Encode()
			var payload struct {
				Bars map[string][]struct {
					C, H, L, O, V float64
					T             string `json:"t"`
				} `json:"bars"`
				NextPageToken string `json:"next_page_token"`
			}
			// anonymous numeric fields require explicit tags through map fallback below; use a map-shaped decode for correctness.
			var rawPayload struct {
				Bars          map[string][]map[string]any `json:"bars"`
				NextPageToken string                      `json:"next_page_token"`
			}
			_ = payload
			if err := getJSON(ctx, client, raw, headers, &rawPayload); err != nil {
				complete = false
				failDetail = "provider page fetch failed: " + err.Error()
				anyFailure = true
				break
			}
			pages++
			for sym, rows := range rawPayload.Bars {
				sym = normalizeSymbol(sym)
				for _, r := range rows {
					ts := strings.TrimSpace(fmt.Sprint(r["t"]))
					t, err := time.Parse(time.RFC3339Nano, ts)
					c := toFloat(r["c"])
					if err != nil || c <= 0 {
						continue
					}
					acc[sym] = append(acc[sym], Bar{T: t.Unix(), O: toFloat(r["o"]), H: toFloat(r["h"]), L: toFloat(r["l"]), C: c, V: toFloat(r["v"])})
				}
			}
			next := strings.TrimSpace(rawPayload.NextPageToken)
			if next == "" {
				break
			}
			if seenTokens[next] {
				complete = false
				failDetail = "repeated pagination token"
				anyFailure = true
				break
			}
			seenTokens[next] = true
			pageToken = next
		}
		nowMs := time.Now().UnixMilli()
		e.mu.Lock()
		if e.rawHistoryCoverage == nil {
			e.rawHistoryCoverage = map[string]RawHistoryCoverage{}
		}
		for _, sym := range batch {
			freshBars := mergeBarsByTime(nil, acc[sym])
			existing := e.bars[sym]["daily-raw"]
			rows := freshBars
			if !complete {
				rows = mergeBarsByTime(existing, freshBars)
			}
			if e.bars[sym] == nil {
				e.bars[sym] = map[string][]Bar{}
			}
			if len(rows) > 0 {
				e.bars[sym]["daily-raw"] = rows
				e.lastUpdated["history-raw:"+sym] = nowMs
				loaded += len(freshBars)
			}
			cov := RawHistoryCoverage{Symbol: sym, RequestedStart: startAt, RequestedEnd: endAt, BarCount: len(rows), PageCount: pages, PaginationComplete: complete, UpdatedAt: nowMs}
			if len(rows) == 0 && complete {
				cov.State = "UNAVAILABLE"
				cov.Detail = "Provider returned no raw daily bars for the requested range."
			} else if complete {
				cov.State = "COMPLETE"
				cov.Detail = "All provider pages reconciled for the requested raw-history range."
			} else {
				cov.State = "PARTIAL"
				cov.Detail = failDetail
			}
			if len(rows) > 0 {
				cov.FirstBarAt = rows[0].T * 1000
				cov.LastBarAt = rows[len(rows)-1].T * 1000
			}
			e.rawHistoryCoverage[sym] = cov
		}
		e.mu.Unlock()
	}
	if anyFailure {
		e.setHealth("corporate-actions-raw-history", "degraded · raw comparison history partial; coverage details retained per symbol")
	} else if loaded > 0 {
		e.setHealth("corporate-actions-raw-history", fmt.Sprintf("healthy · %d raw daily bars reconciled with complete pagination", loaded))
	} else {
		e.setHealth("corporate-actions-raw-history", "unavailable · provider returned no raw comparison bars")
	}
	return loaded
}
