package main

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math"
	"sort"
	"strings"
	"time"
)

const smartRouterPolicyVersion = "smart-router-v2.0.0"

const (
	providerCapabilitySupported              = "SUPPORTED"
	providerCapabilityNotSupported           = "NOT_SUPPORTED"
	providerCapabilityNotConfigured          = "NOT_CONFIGURED"
	providerCapabilityNotEntitled            = "NOT_ENTITLED"
	providerCapabilityUnknown                = "UNKNOWN"
	providerCapabilityTemporarilyUnavailable = "TEMPORARILY_UNAVAILABLE"
	providerCapabilityRateLimited            = "RATE_LIMITED"
	providerCapabilitySaturated              = "SATURATED"
	providerCapabilityDegraded               = "DEGRADED"
)

type ProviderCapabilityStateRecord struct {
	Key             string `json:"key"`
	Provider        string `json:"provider"`
	Dataset         string `json:"dataset"`
	InstrumentClass string `json:"instrumentClass"`
	Session         string `json:"session"`
	State           string `json:"state"`
	Reason          string `json:"reason,omitempty"`
	FirstObservedAt int64  `json:"firstObservedAt,omitempty"`
	LastObservedAt  int64  `json:"lastObservedAt,omitempty"`
	RevalidateAt    int64  `json:"revalidateAt,omitempty"`
	FailureCount    int64  `json:"failureCount,omitempty"`
	SuccessCount    int64  `json:"successCount,omitempty"`
	CallsAvoided    int64  `json:"callsAvoided,omitempty"`
	PolicyVersion   string `json:"policyVersion"`
}

type ProviderRouteScore struct {
	Provider string   `json:"provider"`
	Dataset  string   `json:"dataset"`
	Score    float64  `json:"score"`
	Eligible bool     `json:"eligible"`
	State    string   `json:"state"`
	Reasons  []string `json:"reasons,omitempty"`
}

type SmartRouterScorecard struct {
	PolicyVersion       string `json:"policyVersion"`
	RouteDecisions      int64  `json:"routeDecisions"`
	FallbackDecisions   int64  `json:"fallbackDecisions"`
	EntitlementAvoided  int64  `json:"entitlementAvoided"`
	CircuitAvoided      int64  `json:"circuitAvoided"`
	CapacityAvoided     int64  `json:"capacityAvoided"`
	SourceDisagreements int64  `json:"sourceDisagreements"`
	LastDecisionAt      int64  `json:"lastDecisionAt,omitempty"`
}

func providerInstrumentClass(dataset string) string {
	switch dataset {
	case "VIX / Indices":
		return "INDEX"
	case "Macro":
		return "MACRO"
	case "SEC":
		return "FILING"
	default:
		return "US_EQUITY"
	}
}

func providerCapabilityKey(provider, dataset, session string) string {
	clean := func(s string) string {
		s = strings.ToLower(strings.TrimSpace(s))
		s = strings.NewReplacer(" ", "-", "/", "-", "_", "-").Replace(s)
		for strings.Contains(s, "--") {
			s = strings.ReplaceAll(s, "--", "-")
		}
		return strings.Trim(s, "-")
	}
	if session == "" {
		session = "any"
	}
	return strings.Join([]string{clean(provider), clean(dataset), clean(providerInstrumentClass(dataset)), clean(session)}, "|")
}

func providerCapabilityStateActive(r ProviderCapabilityStateRecord, now int64) bool {
	if strings.TrimSpace(r.State) == "" || r.State == providerCapabilitySupported || r.State == providerCapabilityUnknown {
		return false
	}
	return r.RevalidateAt <= 0 || now < r.RevalidateAt
}

func classifyProviderCapabilityFailure(err error) (state string, ttl time.Duration, reason string) {
	if err == nil {
		return providerCapabilityUnknown, 0, ""
	}
	reason = strings.TrimSpace(err.Error())
	low := strings.ToLower(reason)
	// Specific transient/capacity conditions must win before generic plan/subscription wording.
	if strings.Contains(low, "http 429") || strings.Contains(low, "rate limit") || strings.Contains(low, "too many requests") {
		return providerCapabilityRateLimited, 5 * time.Minute, reason
	}
	if strings.Contains(low, "saturated") || strings.Contains(low, "capacity") || strings.Contains(low, "subscription limit") {
		return providerCapabilitySaturated, 2 * time.Minute, reason
	}
	// Higher-plan/entitlement failures are deterministic capability truth, not outages.
	if strings.Contains(low, "not entitled") || strings.Contains(low, "plan limited") || strings.Contains(low, "payment required") ||
		strings.Contains(low, "subscription") || strings.Contains(low, "available starting with pro") || strings.Contains(low, "available starting with venture") ||
		strings.Contains(low, "higher plan") || strings.Contains(low, "upgrade your plan") || strings.Contains(low, "http 402") || strings.Contains(low, "http 403") ||
		(strings.Contains(low, "http 404") && (strings.Contains(low, "pro") || strings.Contains(low, "venture") || strings.Contains(low, "plan"))) {
		return providerCapabilityNotEntitled, 12 * time.Hour, reason
	}
	if strings.Contains(low, "timeout") || strings.Contains(low, "deadline exceeded") || strings.Contains(low, "temporarily unavailable") || strings.Contains(low, "connection reset") {
		return providerCapabilityTemporarilyUnavailable, 90 * time.Second, reason
	}
	return providerCapabilityDegraded, 2 * time.Minute, reason
}

func (e *Engine) providerCapabilityRecordLocked(provider, dataset string, now time.Time) ProviderCapabilityStateRecord {
	if e.providerCapabilityStates == nil {
		e.providerCapabilityStates = map[string]ProviderCapabilityStateRecord{}
	}
	key := providerCapabilityKey(provider, dataset, marketSessionET(now))
	r := e.providerCapabilityStates[key]
	if r.Key == "" {
		r = ProviderCapabilityStateRecord{
			Key: key, Provider: provider, Dataset: dataset, InstrumentClass: providerInstrumentClass(dataset), Session: marketSessionET(now),
			State: providerCapabilityUnknown, PolicyVersion: smartRouterPolicyVersion,
		}
	}
	if r.RevalidateAt > 0 && now.UnixMilli() >= r.RevalidateAt && r.State != providerCapabilitySupported {
		r.State = providerCapabilityUnknown
		r.Reason = "cooldown expired; capability eligible for revalidation"
		r.RevalidateAt = 0
		e.providerCapabilityStates[key] = r
	}
	return r
}

func (e *Engine) providerCapabilityAllowed(dataset, provider string, now time.Time) (bool, string) {
	e.mu.Lock()
	r := e.providerCapabilityRecordLocked(provider, dataset, now)
	if providerCapabilityStateActive(r, now.UnixMilli()) {
		r.CallsAvoided++
		e.providerCapabilityStates[r.Key] = r
		e.providerCallsAvoided++
		switch r.State {
		case providerCapabilityNotConfigured, providerCapabilityNotEntitled, providerCapabilityNotSupported:
			e.smartRouterScorecard.EntitlementAvoided++
		case providerCapabilitySaturated:
			e.smartRouterScorecard.CapacityAvoided++
		case providerCapabilityRateLimited, providerCapabilityTemporarilyUnavailable, providerCapabilityDegraded:
			e.smartRouterScorecard.CircuitAvoided++
		}
		e.mu.Unlock()
		return false, r.State
	}
	e.mu.Unlock()
	return true, ""
}

func (e *Engine) recordProviderCapabilityFailure(dataset, provider string, err error) {
	state, ttl, reason := classifyProviderCapabilityFailure(err)
	if err == nil {
		return
	}
	now := time.Now()
	e.mu.Lock()
	r := e.providerCapabilityRecordLocked(provider, dataset, now)
	if r.FirstObservedAt == 0 {
		r.FirstObservedAt = now.UnixMilli()
	}
	r.LastObservedAt = now.UnixMilli()
	r.State = state
	r.Reason = reason
	r.FailureCount++
	r.PolicyVersion = smartRouterPolicyVersion
	if ttl > 0 {
		r.RevalidateAt = now.Add(ttl).UnixMilli()
	}
	e.providerCapabilityStates[r.Key] = r
	e.mu.Unlock()
	e.persistProviderCapabilityState(r)
}

func (e *Engine) recordProviderCapabilitySuccess(dataset, provider string) {
	now := time.Now()
	e.mu.Lock()
	r := e.providerCapabilityRecordLocked(provider, dataset, now)
	if r.FirstObservedAt == 0 {
		r.FirstObservedAt = now.UnixMilli()
	}
	r.LastObservedAt = now.UnixMilli()
	r.State = providerCapabilitySupported
	r.Reason = "capability served canonical route successfully"
	r.RevalidateAt = 0
	r.SuccessCount++
	r.PolicyVersion = smartRouterPolicyVersion
	e.providerCapabilityStates[r.Key] = r
	e.mu.Unlock()
	e.persistProviderCapabilityState(r)
}

func (e *Engine) persistProviderCapabilityState(r ProviderCapabilityStateRecord) {
	if e == nil || e.app == nil || e.app.persistence == nil || r.Key == "" {
		return
	}
	raw, err := json.Marshal(r)
	if err != nil {
		return
	}
	sum := sha256.Sum256(raw)
	e.app.persistence.EnqueueIntelligence(PersistenceIntelligenceBatch{Features: []DerivedFeatureRecord{{
		Symbol: "__GLOBAL__", FeatureKey: "provider-capability:" + r.Key, FeatureVersion: smartRouterPolicyVersion,
		AsOf: r.LastObservedAt, SourceHash: hex.EncodeToString(sum[:])[:24], Payload: raw,
	}}})
}

func telemetryForProvider(rows []ProviderRequestDiagnostics, provider string) ProviderRequestDiagnostics {
	for _, row := range rows {
		if strings.EqualFold(row.Provider, provider) {
			return row
		}
	}
	return ProviderRequestDiagnostics{}
}

func smartRouteScore(provider, dataset string, basePriority int, tier WorkTier, cap ProviderCapabilityStateRecord, circuit providerCircuit, telemetry ProviderRequestDiagnostics, session string) ProviderRouteScore {
	score := 1000.0 - float64(basePriority-1)*120
	reasons := []string{fmt.Sprintf("base route priority %d", basePriority)}
	eligible := true
	state := cap.State
	if state == "" {
		state = providerCapabilityUnknown
	}
	switch state {
	case providerCapabilityNotEntitled, providerCapabilityNotSupported:
		score -= 10000
		eligible = false
		reasons = append(reasons, state)
	case providerCapabilityRateLimited, providerCapabilitySaturated:
		score -= 6000
		eligible = false
		reasons = append(reasons, state)
	case providerCapabilityTemporarilyUnavailable:
		score -= 3000
		eligible = false
		reasons = append(reasons, state)
	case providerCapabilityDegraded:
		score -= 350
		reasons = append(reasons, "capability degraded")
	case providerCapabilitySupported:
		score += 60
		reasons = append(reasons, "capability previously served")
	}
	now := time.Now().UnixMilli()
	if circuit.RateLimitedUntil > now || circuit.OpenUntil > now {
		score -= 6000
		eligible = false
		reasons = append(reasons, "circuit suppressed")
	} else if circuit.Failures > 0 {
		score -= math.Min(240, float64(circuit.Failures*60))
		reasons = append(reasons, "recent capability failures")
	}
	completed := telemetry.Successes + telemetry.Errors
	if completed > 0 {
		successPct := float64(telemetry.Successes) / float64(completed) * 100
		score += (successPct - 85) * 1.5
		reasons = append(reasons, fmt.Sprintf("%.0f%% provider success", successPct))
	}
	latency := telemetry.P95LatencyMs
	if latency <= 0 {
		latency = telemetry.AverageLatencyMs
	}
	if latency > 0 {
		penalty := math.Min(180, float64(latency)/30)
		score -= penalty
		reasons = append(reasons, fmt.Sprintf("latency p95/avg %dms", latency))
	}
	if telemetry.BudgetState == "PRESSURE" {
		score -= 90
		reasons = append(reasons, "request budget pressure")
	} else if telemetry.BudgetState == "EXHAUSTED" || telemetry.BudgetState == "RATE LIMITED" {
		if tier >= WorkTierRadarPromoted {
			eligible = false
			score -= 5000
		} else {
			score -= 250
		}
		reasons = append(reasons, strings.ToLower(telemetry.BudgetState))
	}
	// Tier 0/1 values freshness/reliability over cost; lower tiers become more cost sensitive.
	if tier >= WorkTierBroadDiscovery {
		switch providerCostClass(provider) {
		case "Broker/data entitlement":
			score -= 15
		case "Free tier / optional paid upgrade":
			score -= 8
		}
	}
	if dataset == "US Live Equities" && (session == "regular" || session == "pre-market" || session == "after-hours") {
		if strings.EqualFold(provider, "Finnhub") {
			score += 220
			reasons = append(reasons, "preferred live-stream role when suitable")
		} else if strings.EqualFold(provider, "Alpaca") {
			score += 40
			reasons = append(reasons, "IEX live/snapshot overflow role")
		}
	}
	return ProviderRouteScore{Provider: provider, Dataset: dataset, Score: score, Eligible: eligible, State: state, Reasons: reasons}
}

func (e *Engine) rankedProviderRoute(dataset string, tier WorkTier, settings Settings, secrets Secrets, now time.Time) []ProviderRouteScore {
	chain := routeChains()[dataset]
	telemetry := e.providerTelemetry.Diagnostics()
	session := marketSessionET(now)
	out := make([]ProviderRouteScore, 0, len(chain))
	e.mu.Lock()
	for i, provider := range chain {
		cap := e.providerCapabilityRecordLocked(provider, dataset, now)
		circuit := e.providerCapabilityCircuits[providerCapabilityCircuitKey(provider, dataset)]
		score := smartRouteScore(provider, dataset, i+1, tier, cap, circuit, telemetryForProvider(telemetry, provider), session)
		if !e.providerConfigured(provider, secrets, settings) {
			score.Eligible = false
			score.Score -= 10000
			score.State = providerCapabilityNotConfigured
			score.Reasons = append(score.Reasons, "not configured; entitlement not probed")
		}
		out = append(out, score)
	}
	e.mu.Unlock()
	sort.SliceStable(out, func(i, j int) bool {
		if out[i].Eligible != out[j].Eligible {
			return out[i].Eligible
		}
		if out[i].Score != out[j].Score {
			return out[i].Score > out[j].Score
		}
		return out[i].Provider < out[j].Provider
	})
	return out
}

func providerCapabilityCircuitKey(provider, dataset string) string {
	return providerKey(provider) + "|" + strings.ToLower(strings.TrimSpace(dataset))
}

func (e *Engine) capabilityCircuitStatusLocked(provider, dataset string, now int64) string {
	c := e.providerCapabilityCircuits[providerCapabilityCircuitKey(provider, dataset)]
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

func (e *Engine) providerAllowedFor(dataset, provider string) bool {
	now := time.Now()
	if ok, _ := e.providerCapabilityAllowed(dataset, provider, now); !ok {
		return false
	}
	e.mu.RLock()
	st := e.capabilityCircuitStatusLocked(provider, dataset, now.UnixMilli())
	global := e.circuitStatusLocked(provider, now.UnixMilli())
	e.mu.RUnlock()
	return st != "OPEN" && st != "RATE LIMITED" && global != "OPEN" && global != "RATE LIMITED"
}

func (e *Engine) recordProviderCapabilityCircuitSuccess(dataset, provider string) {
	e.mu.Lock()
	k := providerCapabilityCircuitKey(provider, dataset)
	c := e.providerCapabilityCircuits[k]
	c.Failures = 0
	c.OpenUntil = 0
	c.RateLimitedUntil = 0
	c.LastSuccess = time.Now().UnixMilli()
	c.LastError = ""
	e.providerCapabilityCircuits[k] = c
	e.mu.Unlock()
	e.recordProviderCapabilitySuccess(dataset, provider)
}

func (e *Engine) recordProviderCapabilityCircuitFailure(dataset, provider string, err error) {
	now := time.Now()
	e.mu.Lock()
	k := providerCapabilityCircuitKey(provider, dataset)
	c := e.providerCapabilityCircuits[k]
	c.Failures++
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
	e.providerCapabilityCircuits[k] = c
	e.mu.Unlock()
	e.recordProviderCapabilityFailure(dataset, provider, err)
}

func (e *Engine) recordProviderCapabilityLatency(dataset, provider string, started time.Time) {
	e.mu.Lock()
	defer e.mu.Unlock()
	k := providerCapabilityCircuitKey(provider, dataset)
	c := e.providerCapabilityCircuits[k]
	c.Attempts++
	c.LatencyMs = time.Since(started).Milliseconds()
	e.providerCapabilityCircuits[k] = c
}
