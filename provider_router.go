package main

import (
	"context"
	"fmt"
	"sort"
	"strings"
	"time"
)

type ProviderRouteHop struct {
	Provider      string                     `json:"provider"`
	Role          string                     `json:"role"`
	Configured    bool                       `json:"configured"`
	Health        string                     `json:"health"`
	Circuit       string                     `json:"circuit"`
	Priority      int                        `json:"priority"`
	Quota         string                     `json:"quota,omitempty"`
	RateLimit     string                     `json:"rateLimit,omitempty"`
	LatencyMs     int64                      `json:"latencyMs,omitempty"`
	LastSuccess   int64                      `json:"lastSuccess,omitempty"`
	LastFailure   int64                      `json:"lastFailure,omitempty"`
	FailureCount  int                        `json:"failureCount,omitempty"`
	Attempts      int64                      `json:"attempts,omitempty"`
	LastError     string                     `json:"lastError,omitempty"`
	Recovery      string                     `json:"recovery,omitempty"`
	ExpectedDelay string                     `json:"expectedDelay,omitempty"`
	CostClass     string                     `json:"costClass,omitempty"`
	Entitlement   string                     `json:"entitlement,omitempty"`
	DataRights    ProviderDataRightsMetadata `json:"dataRights"`
	Score         float64                    `json:"score,omitempty"`
	ScoreReasons  []string                   `json:"scoreReasons,omitempty"`
}

type ProviderRouteState struct {
	Dataset       string             `json:"dataset"`
	Primary       string             `json:"primary"` // compatibility alias for Preferred
	Active        string             `json:"active"`  // compatibility alias for Serving
	Preferred     string             `json:"preferred"`
	Serving       string             `json:"serving"`
	Reason        string             `json:"reason,omitempty"`
	Capability    string             `json:"capability,omitempty"`
	PolicyVersion string             `json:"policyVersion,omitempty"`
	State         string             `json:"state"`
	Route         []ProviderRouteHop `json:"route"`
	LastSuccess   int64              `json:"lastSuccess,omitempty"`
	Detail        string             `json:"detail,omitempty"`
}

type ProviderRouterSnapshot struct {
	UpdatedAt     int64                `json:"updatedAt"`
	PolicyVersion string               `json:"policyVersion,omitempty"`
	Scorecard     SmartRouterScorecard `json:"scorecard"`
	Routes        []ProviderRouteState `json:"routes"`
}

type FreshnessDiagnostic struct {
	Dataset           string   `json:"dataset"`
	State             string   `json:"state"`
	Provider          string   `json:"provider"`
	ProviderTimestamp int64    `json:"providerTimestamp,omitempty"`
	ReceivedAt        int64    `json:"receivedAt,omitempty"`
	DataTimestamp     int64    `json:"dataTimestamp,omitempty"`
	CacheAt           int64    `json:"cacheAt,omitempty"`
	TimestampBasis    string   `json:"timestampBasis,omitempty"`
	FreshnessBasis    string   `json:"freshnessBasis,omitempty"`
	AgeMs             int64    `json:"ageMs,omitempty"` // compatibility: primary freshness/check age
	CheckAgeMs        int64    `json:"checkAgeMs,omitempty"`
	DataAgeMs         int64    `json:"dataAgeMs,omitempty"`
	ExpectedCadenceMs int64    `json:"expectedCadenceMs"`
	FreshLimitMs      int64    `json:"freshLimitMs"`
	StaleLimitMs      int64    `json:"staleLimitMs"`
	NextExpectedAt    int64    `json:"nextExpectedAt,omitempty"`
	Reason            string   `json:"reason,omitempty"`
	Fallback          string   `json:"fallback,omitempty"`
	Affected          []string `json:"affected,omitempty"`
	Session           string   `json:"session"`
	Action            string   `json:"action,omitempty"`
}

type FreshnessSummary struct {
	Live        int `json:"live"`
	Fresh       int `json:"fresh"`
	DueSoon     int `json:"dueSoon"`
	Delayed     int `json:"delayed"`
	Stale       int `json:"stale"`
	Error       int `json:"error"`
	Unavailable int `json:"unavailable"`
	Idle        int `json:"idle"`
}

type providerCircuit struct {
	Failures         int
	OpenUntil        int64
	RateLimitedUntil int64
	LastFailure      int64
	LastSuccess      int64
	LastError        string
	LatencyMs        int64
	Attempts         int64
}

func routeChains() map[string][]string {
	return map[string][]string{
		"US Live Equities":            {"Alpaca", "Finnhub", "Twelve Data"},
		"VIX / Indices":               {"Twelve Data", "yfinance", "CBOE"},
		canonicalHistoricalBarsDataset: {"Alpaca", tradeInsightProviderName, "Twelve Data", "yfinance"},
		"News":                        {"Finnhub", "Marketaux"},
		"Earnings":                    {"Finnhub", "yfinance"},
		"Fundamentals":                {"Finnhub", "SEC", "yfinance"},
		"SEC":                         {"SEC EDGAR"},
		"Macro":                       {"FRED"},
	}
}

func providerKey(provider string) string {
	p := strings.ToLower(strings.TrimSpace(provider))
	p = strings.ReplaceAll(p, " ", "-")
	return p
}

func (e *Engine) providerConfigured(provider string, secrets Secrets, settings Settings) bool {
	switch strings.ToLower(provider) {
	case "alpaca":
		return strings.TrimSpace(secrets.AlpacaKey) != "" && strings.TrimSpace(secrets.AlpacaSecret) != ""
	case "finnhub":
		return strings.TrimSpace(secrets.Finnhub) != ""
	case "tradeinsight":
		return tradeInsightConfigured()
	case "twelve data":
		return strings.TrimSpace(secrets.TwelveData) != ""
	case "marketaux":
		return strings.TrimSpace(secrets.Marketaux) != ""
	case "fred":
		return strings.TrimSpace(secrets.FRED) != ""
	case "sec", "sec edgar":
		return strings.TrimSpace(settings.SECEmail) != ""
	case "yfinance", "cboe":
		return true
	}
	return false
}

func providerQuotaLabel(provider string) string {
	switch provider {
	case "Alpaca":
		return "Entitlement / feed dependent"
	case "Finnhub":
		return "API plan / endpoint dependent"
	case tradeInsightProviderName:
		return "Runtime tier / response headers"
	case "Twelve Data":
		return "Credit based"
	case "Marketaux":
		return "Request quota"
	case "FRED":
		return "Free API key"
	case "SEC", "SEC EDGAR":
		return "Fair-access policy"
	case "yfinance":
		return "Public recovery · best effort"
	case "CBOE":
		return "Public official/delayed"
	}
	return "Provider dependent"
}

func providerCostClass(provider string) string {
	switch provider {
	case "yfinance", "CBOE", "SEC", "SEC EDGAR":
		return "Public / no API fee"
	case tradeInsightProviderName:
		return "Beta / free tier"
	case "FRED", "Marketaux", "Twelve Data", "Finnhub":
		return "Free tier / optional paid upgrade"
	case "Alpaca":
		return "Broker/data entitlement"
	}
	return "Provider dependent"
}

func expectedProviderDelay(dataset, provider string) string {
	if dataset == "VIX / Indices" && provider == "CBOE" {
		return "Official delayed/close validation"
	}
	if provider == "yfinance" {
		return "Recovery-only; may be delayed"
	}
	if dataset == canonicalHistoricalBarsDataset {
		return "Completed-bar cadence"
	}
	if dataset == "SEC" || provider == "SEC EDGAR" {
		return "Filing dissemination cadence"
	}
	if dataset == "Macro" {
		return "Series release cadence"
	}
	return "Near-live when entitled"
}

func (e *Engine) circuitStatusLocked(provider string, now int64) string {
	c := e.providerCircuits[providerKey(provider)]
	if c.RateLimitedUntil > now {
		return "RATE LIMITED"
	}
	if c.OpenUntil > now {
		return "OPEN"
	}
	if c.Failures > 0 {
		return "PROBING"
	}
	return "CLOSED"
}

func (e *Engine) providerAllowed(provider string) bool {
	e.mu.RLock()
	defer e.mu.RUnlock()
	st := e.circuitStatusLocked(provider, time.Now().UnixMilli())
	return st != "OPEN" && st != "RATE LIMITED"
}

func (e *Engine) recordProviderSuccess(provider string) {
	e.mu.Lock()
	defer e.mu.Unlock()
	k := providerKey(provider)
	c := e.providerCircuits[k]
	c.Failures = 0
	c.OpenUntil = 0
	c.RateLimitedUntil = 0
	c.LastSuccess = time.Now().UnixMilli()
	c.LastError = ""
	e.providerCircuits[k] = c
}

func (e *Engine) recordProviderLatency(provider string, started time.Time) {
	e.mu.Lock()
	defer e.mu.Unlock()
	k := providerKey(provider)
	c := e.providerCircuits[k]
	c.Attempts++
	c.LatencyMs = time.Since(started).Milliseconds()
	e.providerCircuits[k] = c
}

func (e *Engine) recordProviderFailure(provider string, err error) {
	e.mu.Lock()
	defer e.mu.Unlock()
	k := providerKey(provider)
	c := e.providerCircuits[k]
	c.Failures++
	now := time.Now()
	c.LastFailure = now.UnixMilli()
	if err != nil {
		c.LastError = err.Error()
		low := strings.ToLower(c.LastError)
		if strings.Contains(low, "http 429") || strings.Contains(low, "rate limit") || strings.Contains(low, "too many requests") {
			c.RateLimitedUntil = now.Add(5 * time.Minute).UnixMilli()
			c.OpenUntil = c.RateLimitedUntil
		}
	}
	if c.Failures >= 3 && c.OpenUntil <= now.UnixMilli() {
		c.OpenUntil = now.Add(2 * time.Minute).UnixMilli()
	}
	e.providerCircuits[k] = c
}

func sourceProvider(source string) string {
	s := strings.ToLower(source)
	switch {
	case strings.Contains(s, "twelvedata"):
		return "Twelve Data"
	case strings.Contains(s, "tradeinsight"):
		return tradeInsightProviderName
	case strings.Contains(s, "yfinance"), strings.Contains(s, "yahoo"):
		return "yfinance"
	case strings.Contains(s, "cboe"):
		return "CBOE"
	case strings.Contains(s, "alpaca"):
		return "Alpaca"
	case strings.Contains(s, "finnhub"):
		return "Finnhub"
	case strings.Contains(s, "marketaux"):
		return "Marketaux"
	case strings.Contains(s, "sec"):
		return "SEC EDGAR"
	case strings.Contains(s, "fred"):
		return "FRED"
	}
	return "—"
}

type providerRouteAttempt func(context.Context) bool

// executeProviderRoute is the single executable routing authority. Provider-
// specific loaders know how to fetch/normalize one source; only this function
// decides which provider is attempted and in what order.
func (e *Engine) executeProviderRoute(ctx context.Context, dataset string, attempts map[string]providerRouteAttempt) (string, bool) {
	e.app.mu.RLock()
	settings := clone(e.app.state.Settings)
	secrets := clone(e.app.secrets)
	e.app.mu.RUnlock()
	tier := workTierFromContext(ctx, WorkTierUserActionable)
	ranked := e.rankedProviderRoute(dataset, tier, settings, secrets, time.Now())
	for _, candidate := range ranked {
		provider := candidate.Provider
		attempt := attempts[provider]
		if attempt == nil || !candidate.Eligible {
			continue
		}
		if !e.providerAllowedFor(dataset, provider) {
			continue
		}

		e.mu.RLock()
		before := e.providerCircuits[providerKey(provider)]
		e.mu.RUnlock()
		started := time.Now()
		ok := attempt(ctx)
		e.recordProviderLatency(provider, started)
		e.recordProviderCapabilityLatency(dataset, provider, started)
		if ok {
			e.recordProviderCapabilityCircuitSuccess(dataset, provider)
			e.mu.Lock()
			e.smartRouterScorecard.RouteDecisions++
			e.smartRouterScorecard.LastDecisionAt = time.Now().UnixMilli()
			if preferred := routeChains()[dataset]; len(preferred) > 0 && !strings.EqualFold(provider, preferred[0]) {
				e.smartRouterScorecard.FallbackDecisions++
			}
			e.mu.Unlock()
			return provider, true
		}
		e.mu.RLock()
		after := e.providerCircuits[providerKey(provider)]
		e.mu.RUnlock()
		if after.LastFailure > before.LastFailure || (after.LastFailure == before.LastFailure && after.LastError != "" && after.LastError != before.LastError) {
			e.recordProviderCapabilityCircuitFailure(dataset, provider, fmt.Errorf("%s", after.LastError))
		}
	}
	return "", false
}

func (e *Engine) buildProviderRouterSnapshot(settings Settings, secrets Secrets, quotes map[string]Quote, last map[string]int64) ProviderRouterSnapshot {
	nowTime := time.Now()
	now := nowTime.UnixMilli()
	session := marketSessionET(nowTime)
	routes := routeChains()
	keys := make([]string, 0, len(routes))
	for k := range routes {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	telemetry := e.providerTelemetry.Diagnostics()
	scorecard := e.smartRouterScorecard
	scorecard.PolicyVersion = smartRouterPolicyVersion
	out := ProviderRouterSnapshot{UpdatedAt: now, PolicyVersion: smartRouterPolicyVersion, Scorecard: scorecard}
	for _, dataset := range keys {
		chain := routes[dataset]
		active := ""
		detail := ""
		lastSuccess := int64(0)
		switch dataset {
		case "VIX / Indices":
			active = sourceProvider(quotes["VIX"].Source)
			lastSuccess = quotes["VIX"].UpdatedAt
			detail = quotes["VIX"].Source
		case canonicalHistoricalBarsDataset:
			lastSuccess = last["history"]
			active = sourceProvider(e.health["history"])
		case "News":
			lastSuccess = last["news"]
			active = sourceProvider(e.health["news"])
		case "Earnings":
			lastSuccess = last["earnings"]
			active = sourceProvider(e.health["earnings"])
		case "Fundamentals":
			lastSuccess = last["fundamentals"]
			active = sourceProvider(e.health["fundamentals"])
			if strings.Contains(strings.ToLower(e.health["fundamentals"]), "sec") {
				active = "SEC"
			}
		case "SEC":
			lastSuccess = last["filings"]
			active = "SEC EDGAR"
		case "Macro":
			lastSuccess = last["macro"]
			active = "FRED"
		case "US Live Equities":
			alpacaLast := maxInt64(e.lastAlpacaStreamAt, e.lastAlpacaAt)
			lastSuccess = maxInt64(e.lastTradeAt, alpacaLast)
			if alpacaLast > 0 && alpacaLast >= e.lastTradeAt {
				active = "Alpaca"
			} else if e.lastTradeAt > 0 {
				active = "Finnhub"
			}
		}
		if active == "—" {
			active = ""
		}

		hops := make([]ProviderRouteHop, 0, len(chain))
		preferred := ""
		preferredScore := -1e18
		preferredReason := ""
		for i, provider := range chain {
			configured := e.providerConfigured(provider, secrets, settings)
			capKey := providerCapabilityKey(provider, dataset, session)
			cap := e.providerCapabilityStates[capKey]
			if cap.Key == "" {
				cap = ProviderCapabilityStateRecord{Key: capKey, Provider: provider, Dataset: dataset, InstrumentClass: providerInstrumentClass(dataset), Session: session, State: providerCapabilityUnknown, PolicyVersion: smartRouterPolicyVersion}
			}
			capCircuit := e.providerCapabilityCircuits[providerCapabilityCircuitKey(provider, dataset)]
			globalCircuit := e.providerCircuits[providerKey(provider)]
			if capCircuit.LastSuccess == 0 && capCircuit.LastFailure == 0 && capCircuit.Failures == 0 {
				capCircuit = globalCircuit
			}
			score := smartRouteScore(provider, dataset, i+1, WorkTierUserActionable, cap, capCircuit, telemetryForProvider(telemetry, provider), session)
			if !configured {
				score.Eligible = false
				score.State = providerCapabilityNotConfigured
				score.Reasons = append(score.Reasons, "not configured; entitlement not probed")
			}
			if score.Eligible && score.Score > preferredScore {
				preferred = provider
				preferredScore = score.Score
				preferredReason = strings.Join(score.Reasons, " · ")
			}

			circuit := e.capabilityCircuitStatusLocked(provider, dataset, now)
			if circuit == "CLOSED" {
				gc := e.circuitStatusLocked(provider, now)
				if gc == "OPEN" || gc == "RATE LIMITED" || gc == "PROBING" {
					circuit = gc
				}
			}
			health := "AVAILABLE"
			entitlement := cap.State
			if entitlement == "" {
				entitlement = providerCapabilityUnknown
			}
			if !configured {
				health = "NOT CONFIGURED"
				entitlement = providerCapabilityUnknown
			} else if providerCapabilityStateActive(cap, now) {
				health = cap.State
			} else if circuit == "OPEN" || circuit == "RATE LIMITED" {
				health = "DEGRADED"
			}
			role := "fallback"
			if i == 0 {
				role = "configured primary"
			}
			if provider == "CBOE" {
				role = "validation / delayed fallback"
			}
			if provider == "yfinance" {
				role = "recovery fallback"
			}
			c := capCircuit
			rate := "NORMAL"
			if c.RateLimitedUntil > now {
				rate = "RATE LIMITED"
			}
			recovery := "READY"
			if circuit == "OPEN" || circuit == "RATE LIMITED" || providerCapabilityStateActive(cap, now) {
				recovery = "SUPPRESSED"
			} else if circuit == "PROBING" {
				recovery = "PROBING"
			} else if c.Failures == 0 && c.LastFailure > 0 && c.LastSuccess > c.LastFailure {
				recovery = "RECOVERED"
			}
			hops = append(hops, ProviderRouteHop{
				Provider: provider, Role: role, Configured: configured, Health: health, Circuit: circuit, Priority: i + 1,
				Quota: providerQuotaLabel(provider), RateLimit: rate, LatencyMs: c.LatencyMs, LastSuccess: c.LastSuccess, LastFailure: c.LastFailure,
				FailureCount: c.Failures, Attempts: c.Attempts, LastError: defaultString(cap.Reason, c.LastError), Recovery: recovery,
				ExpectedDelay: expectedProviderDelay(dataset, provider), CostClass: providerCostClass(provider), Entitlement: entitlement,
				DataRights: providerDataRightsMetadata(provider), Score: score.Score, ScoreReasons: append([]string(nil), score.Reasons...),
			})
		}
		if preferred == "" && len(chain) > 0 {
			preferred = chain[0]
			preferredReason = "configured route order; no eligible provider has enough current evidence"
		}
		if active == "" {
			for _, h := range hops {
				if h.Configured && h.Recovery != "SUPPRESSED" && h.Circuit != "OPEN" && h.Circuit != "RATE LIMITED" {
					active = h.Provider
					break
				}
		}

		state := "READY"
		reason := "Preferred provider is serving canonical data."
		if active == "" {
			state = "UNAVAILABLE"
			reason = "No eligible provider is currently serving canonical data."
		} else if preferred != "" && !strings.EqualFold(active, preferred) {
			state = "FALLBACK"
			reason = "Serving provider differs from current preferred route."
			for _, h := range hops {
				if !strings.EqualFold(h.Provider, preferred) {
					continue
				}
				switch {
				case h.Entitlement == providerCapabilityNotEntitled:
					reason = "Preferred provider is NOT_ENTITLED for this capability/session."
				case h.Circuit == "RATE LIMITED":
					reason = "Preferred provider is RATE_LIMITED for this capability."
				case h.Circuit == "OPEN":
					reason = "Preferred provider capability circuit is OPEN."
				case h.Health == providerCapabilityTemporarilyUnavailable || h.Health == providerCapabilityDegraded:
					reason = "Preferred provider capability is temporarily degraded."
				default:
					reason = "Serving provider currently has the usable canonical evidence for this dataset."
				}
				break
			}
		}
		if preferredReason != "" {
			detail = strings.TrimSpace(strings.Join([]string{detail, "router=" + preferredReason}, " · "))
		}
		primary := ""
		if len(chain) > 0 {
			primary = chain[0]
		}
		out.Routes = append(out.Routes, ProviderRouteState{
			Dataset: dataset, Primary: primary, Active: active, Preferred: preferred, Serving: active, Reason: reason,
			Capability: dataset, PolicyVersion: smartRouterPolicyVersion, State: state, Route: hops, LastSuccess: lastSuccess, Detail: detail,
		})
	}
	return out
}
