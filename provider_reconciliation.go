package main

import (
	"fmt"
	"sort"
	"strings"
	"time"
)

type ProviderQuoteObservation struct {
	Symbol            string  `json:"symbol"`
	Provider          string  `json:"provider"`
	Price             float64 `json:"price"`
	Bid               float64 `json:"bid,omitempty"`
	Ask               float64 `json:"ask,omitempty"`
	ProviderTimestamp int64   `json:"providerTimestamp,omitempty"`
	ReceivedAt        int64   `json:"receivedAt"`
	Source            string  `json:"source,omitempty"`
	FeedType          string  `json:"feedType,omitempty"`
	DataState         string  `json:"dataState,omitempty"`
}

type ProviderReconciliationDecision struct {
	Dataset           string                     `json:"dataset"`
	Symbol            string                     `json:"symbol,omitempty"`
	State             string                     `json:"state"` // AGREED, CONFLICT, SINGLE SOURCE, STALE
	Observations      []ProviderQuoteObservation `json:"observations,omitempty"`
	CanonicalProvider string                     `json:"canonicalProvider,omitempty"`
	CanonicalValue    float64                    `json:"canonicalValue,omitempty"`
	DifferencePct     float64                    `json:"differencePct,omitempty"`
	Reason            string                     `json:"reason"`
	UpdatedAt         int64                      `json:"updatedAt"`
}

type ResearchEvidenceComponent struct {
	Dataset    string `json:"dataset"`
	Symbol     string `json:"symbol"`
	State      string `json:"state"` // FRESH, PARTIAL, STALE, BLOCKED
	Required   bool   `json:"required"`
	Critical   bool   `json:"critical,omitempty"`
	Source     string `json:"source,omitempty"`
	CheckAt    int64  `json:"checkAt,omitempty"`
	DataAt     int64  `json:"dataAt,omitempty"`
	CheckAgeMs int64  `json:"checkAgeMs,omitempty"`
	DataAgeMs  int64  `json:"dataAgeMs,omitempty"`
	Detail     string `json:"detail,omitempty"`
}

type ResearchPackageTruth struct {
	Symbol             string                      `json:"symbol"`
	State              string                      `json:"state"` // FRESH, PARTIAL, STALE, BLOCKED
	Components         []ResearchEvidenceComponent `json:"components"`
	BlockingReasons    []string                    `json:"blockingReasons,omitempty"`
	EvidenceSnapshotID string                      `json:"evidenceSnapshotId,omitempty"`
	GeneratedAt        int64                       `json:"generatedAt"`
}

type RawHistoryCoverage struct {
	Symbol             string `json:"symbol"`
	State              string `json:"state"` // COMPLETE, PARTIAL, UNAVAILABLE, ERROR
	RequestedStart     string `json:"requestedStart,omitempty"`
	RequestedEnd       string `json:"requestedEnd,omitempty"`
	FirstBarAt         int64  `json:"firstBarAt,omitempty"`
	LastBarAt          int64  `json:"lastBarAt,omitempty"`
	BarCount           int    `json:"barCount"`
	PageCount          int    `json:"pageCount"`
	PaginationComplete bool   `json:"paginationComplete"`
	UpdatedAt          int64  `json:"updatedAt,omitempty"`
	Detail             string `json:"detail,omitempty"`
}

type SymbolLifecycleState struct {
	Symbol        string `json:"symbol"`
	State         string `json:"state"` // ACTIVE, NEW LISTING / IPO, NAME OR TICKER CHANGE, MERGER PENDING / MERGED, DELISTING PENDING / DELISTED, CORPORATE ACTION PENDING, UNKNOWN
	Source        string `json:"source,omitempty"`
	EvidenceID    string `json:"evidenceId,omitempty"`
	EffectiveDate string `json:"effectiveDate,omitempty"`
	ProcessDate   string `json:"processDate,omitempty"`
	FirstSeenAt   int64  `json:"firstSeenAt,omitempty"`
	UpdatedAt     int64  `json:"updatedAt,omitempty"`
	Reason        string `json:"reason"`
}

type CorporateActionTruth struct {
	Actions             []CorporateAction               `json:"actions,omitempty"`
	SymbolLineage       map[string]string               `json:"symbolLineage,omitempty"`
	Lifecycle           map[string]SymbolLifecycleState `json:"lifecycle,omitempty"`
	AffectedSymbols     []string                        `json:"affectedSymbols,omitempty"`
	CanonicalHistory    string                          `json:"canonicalHistory"`
	RawHistoryAvailable map[string]bool                 `json:"rawHistoryAvailable,omitempty"` // compatibility: true only when coverage is COMPLETE
	RawHistoryCoverage  map[string]RawHistoryCoverage   `json:"rawHistoryCoverage,omitempty"`
	UpdatedAt           int64                           `json:"updatedAt,omitempty"`
}

type EvidenceSnapshot struct {
	ID                       string                           `json:"id"`
	Symbol                   string                           `json:"symbol"`
	GeneratedAt              int64                            `json:"generatedAt"`
	ResearchState            string                           `json:"researchState"`
	Components               []ResearchEvidenceComponent      `json:"components"`
	ProviderReconciliation   []ProviderReconciliationDecision `json:"providerReconciliation,omitempty"`
	CorporateActions         []CorporateAction                `json:"corporateActions,omitempty"`
	CorporateHistoryCoverage []RawHistoryCoverage             `json:"corporateHistoryCoverage,omitempty"`
	SymbolLifecycle          *SymbolLifecycleState            `json:"symbolLifecycle,omitempty"`
}

type ResearchPackageContext struct {
	CatalystReactions map[string]CatalystReactionState
	Global            GlobalMarketContext
}

func providerFromQuoteSource(source string) string {
	s := strings.ToLower(strings.TrimSpace(source))
	switch {
	case strings.Contains(s, "alpaca"):
		return "Alpaca"
	case strings.Contains(s, "finnhub"):
		return "Finnhub"
	case strings.Contains(s, "twelve"):
		return "Twelve Data"
	case strings.Contains(s, "yahoo"), strings.Contains(s, "yfinance"):
		return "yfinance"
	case strings.Contains(s, "cboe"):
		return "CBOE"
	case strings.Contains(s, "demo"):
		return "Demo"
	}
	if strings.TrimSpace(source) == "" {
		return "Unknown"
	}
	return strings.TrimSpace(source)
}

func observationFromQuote(symbol string, q Quote) ProviderQuoteObservation {
	return ProviderQuoteObservation{Symbol: normalizeSymbol(symbol), Provider: providerFromQuoteSource(q.Source), Price: q.Price, Bid: q.Bid, Ask: q.Ask, ProviderTimestamp: q.ProviderTimestamp, ReceivedAt: q.UpdatedAt, Source: q.Source, FeedType: q.FeedType, DataState: q.DataState}
}

func (e *Engine) recordProviderQuoteObservation(symbol string, q Quote) {
	symbol = normalizeSymbol(symbol)
	if symbol == "" || q.Price <= 0 {
		return
	}
	if q.UpdatedAt <= 0 {
		q.UpdatedAt = time.Now().UnixMilli()
	}
	provider := providerFromQuoteSource(q.Source)
	e.mu.Lock()
	defer e.mu.Unlock()
	if e.providerQuotes == nil {
		e.providerQuotes = map[string]map[string]Quote{}
	}
	if e.providerQuotes[symbol] == nil {
		e.providerQuotes[symbol] = map[string]Quote{}
	}
	e.providerQuotes[symbol][provider] = q
}

func routeProviderPriority(route ProviderRouteState, provider string) int {
	for _, h := range route.Route {
		if strings.EqualFold(h.Provider, provider) {
			return h.Priority
		}
	}
	return 999
}

func reconciliationRouteForSymbol(router ProviderRouterSnapshot, symbol string) ProviderRouteState {
	dataset := "US Live Equities"
	if normalizeSymbol(symbol) == "VIX" {
		dataset = "VIX / Indices"
	}
	for _, r := range router.Routes {
		if r.Dataset == dataset {
			return r
		}
	}
	return ProviderRouteState{Dataset: dataset}
}

func reconciliationAgeLimit(now int64) int64 {
	session := marketSessionET(time.UnixMilli(now))
	switch session {
	case "overnight":
		return int64((90 * time.Minute) / time.Millisecond)
	case "closed", "weekend":
		return int64((96 * time.Hour) / time.Millisecond)
	default:
		return int64((2 * time.Minute) / time.Millisecond)
	}
}

const quoteMaxFutureSkewMs int64 = 30_000

// quoteEvidenceTimestampTruth is the shared validity layer. It separates the age
// of market/provider evidence from local receipt age and rejects material clock skew.
// Consumers may apply different freshness limits only after timestamp validity passes.
func quoteEvidenceTimestampTruth(q Quote, now int64) (providerAge, receiptAge int64, valid bool, detail string) {
	marketAt := normalizeObservationMs(q.ProviderTimestamp)
	receivedAt := normalizeObservationMs(q.UpdatedAt)
	if marketAt <= 0 {
		marketAt = receivedAt
	}
	if receivedAt <= 0 {
		receivedAt = marketAt
	}
	if marketAt <= 0 || receivedAt <= 0 {
		return -1, -1, false, "Quote evidence timestamp is unavailable."
	}
	providerAge = now - marketAt
	receiptAge = now - receivedAt
	if providerAge < -quoteMaxFutureSkewMs {
		return providerAge, receiptAge, false, "Provider market timestamp is materially in the future."
	}
	if receiptAge < -quoteMaxFutureSkewMs {
		return providerAge, receiptAge, false, "Local receipt timestamp is materially in the future."
	}
	if providerAge < 0 {
		providerAge = 0
	}
	if receiptAge < 0 {
		receiptAge = 0
	}
	return providerAge, receiptAge, true, "Quote evidence timestamps are valid."
}

func quoteEvidenceAges(q Quote, now int64) (providerAge, receiptAge int64, current bool) {
	providerAge, receiptAge, valid, _ := quoteEvidenceTimestampTruth(q, now)
	if !valid {
		return providerAge, receiptAge, false
	}
	limit := reconciliationAgeLimit(now)
	return providerAge, receiptAge, providerAge <= limit && receiptAge <= limit
}

const reconciliationContemporaneousWindowMs int64 = int64((2 * time.Minute) / time.Millisecond)

// reconciliationObservationTime returns the best comparable evidence timestamp.
// Provider time is authoritative when available; receipt time is only the fallback
// for providers that do not expose a market timestamp.
func reconciliationObservationTime(o ProviderQuoteObservation) int64 {
	if o.ProviderTimestamp > 0 {
		return normalizeObservationMs(o.ProviderTimestamp)
	}
	return normalizeObservationMs(o.ReceivedAt)
}

// keepContemporaneousProviderObservations prevents observations from different
// market moments from manufacturing AGREED/CONFLICT truth. Closed/weekend session
// rules may legitimately keep an old close usable as SINGLE SOURCE evidence, but
// cross-provider reconciliation still requires observations to describe roughly
// the same market moment.
func keepContemporaneousProviderObservations(rows []ProviderQuoteObservation) []ProviderQuoteObservation {
	if len(rows) < 2 {
		return rows
	}
	freshest := int64(0)
	for _, row := range rows {
		ts := reconciliationObservationTime(row)
		if ts > freshest {
			freshest = ts
		}
	}
	if freshest <= 0 {
		return rows
	}
	out := make([]ProviderQuoteObservation, 0, len(rows))
	for _, row := range rows {
		ts := reconciliationObservationTime(row)
		if ts > 0 && freshest-ts <= reconciliationContemporaneousWindowMs {
			out = append(out, row)
		}
	}
	return out
}

// buildProviderReconciliation compares valid contemporaneous provider observations
// without treating absence of a second source as agreement. Canonical provider choice
// follows the Provider Router priority so reconciliation cannot create a parallel owner.
func providerReconciliationConflictCount(rows []ProviderReconciliationDecision) int64 {
	var count int64
	for _, row := range rows {
		if strings.EqualFold(strings.TrimSpace(row.State), "CONFLICT") {
			count++
		}
	}
	return count
}

func buildProviderReconciliation(router ProviderRouterSnapshot, observations map[string]map[string]Quote, canonical map[string]Quote, now int64) []ProviderReconciliationDecision {
	set := map[string]bool{}
	for s := range observations {
		set[normalizeSymbol(s)] = true
	}
	for s := range canonical {
		set[normalizeSymbol(s)] = true
	}
	symbols := make([]string, 0, len(set))
	for s := range set {
		if s != "" {
			symbols = append(symbols, s)
		}
	}
	sort.Strings(symbols)
	out := make([]ProviderReconciliationDecision, 0, len(symbols))
	for _, symbol := range symbols {
		route := reconciliationRouteForSymbol(router, symbol)
		rows := []ProviderQuoteObservation{}
		for provider, q := range observations[symbol] {
			if q.Price <= 0 {
				continue
			}
			_, _, current := quoteEvidenceAges(q, now)
			if !current {
				continue
			}
			x := observationFromQuote(symbol, q)
			x.Provider = provider
			rows = append(rows, x)
		}

		if cq, ok := canonical[symbol]; ok && cq.Price > 0 {
			_, _, current := quoteEvidenceAges(cq, now)
			if current {
				cp := providerFromQuoteSource(cq.Source)
				found := false
				for _, r := range rows {
					if strings.EqualFold(r.Provider, cp) {
						found = true
						break
					}
				}
				if !found {
					rows = append(rows, observationFromQuote(symbol, cq))
				}
			}
		}
		rows = keepContemporaneousProviderObservations(rows)
		sort.Slice(rows, func(i, j int) bool {
			pi, pj := routeProviderPriority(route, rows[i].Provider), routeProviderPriority(route, rows[j].Provider)
			if pi != pj {
				return pi < pj
			}
			ti, tj := rows[i].ProviderTimestamp, rows[j].ProviderTimestamp
			if ti == 0 {
				ti = rows[i].ReceivedAt
			}
			if tj == 0 {
				tj = rows[j].ReceivedAt
			}
			return ti > tj
		})
		if len(rows) == 0 {
			cp, cv := "", 0.0
			reason := "No current provider observation is eligible for cross-provider reconciliation; no agreement claim is made."
			if cq, ok := canonical[symbol]; ok && cq.Price > 0 {
				cp, cv = providerFromQuoteSource(cq.Source), cq.Price
				pa, ra, _ := quoteEvidenceAges(cq, now)
				reason = fmt.Sprintf("Canonical %s quote exists but is not current enough for reconciliation (provider age %d ms; receipt age %d ms).", cp, pa, ra)
			}
			out = append(out, ProviderReconciliationDecision{Dataset: route.Dataset, Symbol: symbol, State: "STALE", CanonicalProvider: cp, CanonicalValue: cv, Reason: reason, UpdatedAt: now})
			continue
		}
		winner := rows[0]
		selection := "highest-priority current provider observation"
		for _, r := range rows {
			if route.Active != "" && strings.EqualFold(r.Provider, route.Active) {
				winner = r
				selection = "Provider Router active provider"
				break
			}
		}
		if cq, ok := canonical[symbol]; ok && cq.Price > 0 {
			cp := providerFromQuoteSource(cq.Source)
			for _, r := range rows {
				if strings.EqualFold(r.Provider, cp) {
					winner = r
					if route.Active != "" && !strings.EqualFold(route.Active, cp) {
						selection = fmt.Sprintf("canonical runtime quote source (Router active is %s)", route.Active)
					} else {
						selection = "canonical runtime quote source / Router route"
					}
					break
				}
			}
		}
		state, diff := "SINGLE SOURCE", 0.0
		reason := fmt.Sprintf("Only one current provider observation is available. Canonical winner is %s from the %s; no cross-provider agreement claim is made.", winner.Provider, selection)
		if len(rows) >= 2 {
			lo, hi := rows[0].Price, rows[0].Price
			for _, r := range rows[1:] {
				if r.Price < lo {
					lo = r.Price
				}
				if r.Price > hi {
					hi = r.Price
				}
			}
			denom := winner.Price
			if denom <= 0 {
				denom = (hi + lo) / 2
			}
			if denom > 0 {
				diff = (hi - lo) / denom * 100
			}
			if diff > 0.10 {
				state = "CONFLICT"
				reason = fmt.Sprintf("Current independent provider observations differ by %.3f%%. Canonical winner is %s from the %s; values are never averaged silently.", diff, winner.Provider, selection)
			} else {
				state = "AGREED"
				reason = fmt.Sprintf("%d current provider observations agree within 0.10%%. Canonical winner is %s from the %s.", len(rows), winner.Provider, selection)
			}
		}
		out = append(out, ProviderReconciliationDecision{Dataset: route.Dataset, Symbol: symbol, State: state, Observations: rows, CanonicalProvider: winner.Provider, CanonicalValue: winner.Price, DifferencePct: diff, Reason: reason, UpdatedAt: now})
	}
	return out
}
