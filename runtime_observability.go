package main

import (
	"context"
	"fmt"
	"net/http"
	"sort"
	"strings"
	"sync"
	"time"
)

type ProviderRequestDiagnostics struct {
	Provider             string  `json:"provider"`
	Requests             int64   `json:"requests"`
	RequestsLastMin      int     `json:"requestsLastMinute"`
	Successes            int64   `json:"successes"`
	Errors               int64   `json:"errors"`
	RateLimited          int64   `json:"rateLimited"`
	InFlight             int     `json:"inFlight"`
	PeakInFlight         int     `json:"peakInFlight"`
	LastLatencyMs        int64   `json:"lastLatencyMs,omitempty"`
	AverageLatencyMs     int64   `json:"averageLatencyMs,omitempty"`
	P50LatencyMs         int64   `json:"p50LatencyMs,omitempty"`
	P95LatencyMs         int64   `json:"p95LatencyMs,omitempty"`
	SuccessPct           float64 `json:"successPct,omitempty"`
	LastRequestAt        int64   `json:"lastRequestAt,omitempty"`
	LastError            string  `json:"lastError,omitempty"`
	BudgetPerMinute      int     `json:"budgetPerMinute,omitempty"`
	BudgetRemaining      int     `json:"budgetRemaining,omitempty"`
	BudgetUtilizationPct int     `json:"budgetUtilizationPct,omitempty"`
	BudgetState          string  `json:"budgetState,omitempty"`
	CooldownUntil        int64   `json:"cooldownUntil,omitempty"`
	BudgetShed           int64   `json:"budgetShed"`
}

type providerRequestState struct {
	provider      string
	requests      int64
	successes     int64
	errors        int64
	rateLimited   int64
	inFlight      int
	peakInFlight  int
	totalLatency  int64
	lastLatency   int64
	lastRequestAt int64
	lastError     string
	recent        []time.Time
	latencies     []int64
	cooldownUntil time.Time
	budgetShed    int64
}

type ProviderTelemetry struct {
	mu       sync.Mutex
	provider map[string]*providerRequestState
}

func NewProviderTelemetry() *ProviderTelemetry {
	return &ProviderTelemetry{provider: map[string]*providerRequestState{}}
}

// providerBudgetPerMinute only returns limits DE.PULSE can justify from its own
// pacing policy. Entitlement-dependent provider limits remain observable rather
// than being guessed or hardcoded.
func providerBudgetPerMinute(provider string) int {
	if strings.EqualFold(strings.TrimSpace(provider), "Finnhub") && finnhubMinRequestInterval > 0 {
		n := int(time.Minute / finnhubMinRequestInterval)
		if n > 1 {
			n--
		} // keep one-request safety margin beneath local pacing
		return n
	}
	return 0
}

func (p *ProviderTelemetry) stateLocked(provider string) *providerRequestState {
	provider = strings.TrimSpace(provider)
	if provider == "" {
		provider = "Unknown"
	}
	st := p.provider[provider]
	if st == nil {
		st = &providerRequestState{provider: provider}
		p.provider[provider] = st
	}
	return st
}

func pruneProviderRecent(now time.Time, st *providerRequestState) {
	cutoff := now.Add(-time.Minute)
	keep := st.recent[:0]
	for _, at := range st.recent {
		if at.After(cutoff) {
			keep = append(keep, at)
		}
	}
	st.recent = keep
}

// Allow sheds only Tier 2-4 requests when a justified local budget is exhausted
// or a provider has just returned rate-limit pressure. Tier 0/1 remains eligible
// so market-critical/actionable truth can recover or fail over.
func (p *ProviderTelemetry) Allow(provider string, tier WorkTier) (bool, string) {
	if p == nil || tier <= WorkTierUserActionable {
		return true, ""
	}
	now := time.Now()
	p.mu.Lock()
	defer p.mu.Unlock()
	st := p.stateLocked(provider)
	pruneProviderRecent(now, st)
	if !st.cooldownUntil.IsZero() && now.Before(st.cooldownUntil) {
		st.budgetShed++
		return false, "provider rate-limit cooldown"
	}
	if budget := providerBudgetPerMinute(provider); budget > 0 && len(st.recent) >= budget {
		st.budgetShed++
		return false, "provider request budget exhausted"
	}
	return true, ""
}

func (p *ProviderTelemetry) begin(provider string) func(error) {
	if p == nil {
		return func(error) {}
	}
	provider = strings.TrimSpace(provider)
	if provider == "" {
		provider = "Unknown"
	}
	started := time.Now()
	p.mu.Lock()
	st := p.stateLocked(provider)
	pruneProviderRecent(started, st)
	st.requests++
	st.inFlight++
	if st.inFlight > st.peakInFlight {
		st.peakInFlight = st.inFlight
	}
	st.lastRequestAt = started.UnixMilli()
	st.recent = append(st.recent, started)
	p.mu.Unlock()
	return func(err error) {
		elapsed := time.Since(started).Milliseconds()
		p.mu.Lock()
		defer p.mu.Unlock()
		st := p.stateLocked(provider)
		if st.inFlight > 0 {
			st.inFlight--
		}
		st.totalLatency += elapsed
		st.lastLatency = elapsed
		st.latencies = append(st.latencies, elapsed)
		if len(st.latencies) > 256 {
			st.latencies = append([]int64(nil), st.latencies[len(st.latencies)-256:]...)
		}
		if err == nil {
			st.successes++
			st.lastError = ""
			return
		}
		st.errors++
		st.lastError = err.Error()
		lowered := strings.ToLower(err.Error())
		if strings.Contains(lowered, "http 429") || strings.Contains(lowered, "rate limit") || strings.Contains(lowered, "too many requests") {
			st.rateLimited++
			st.cooldownUntil = time.Now().Add(30 * time.Second)
		}
	}
}

func latencyPercentiles(values []int64) (p50, p95 int64) {
	if len(values) == 0 {
		return 0, 0
	}
	rows := append([]int64(nil), values...)
	sort.Slice(rows, func(i, j int) bool { return rows[i] < rows[j] })
	at := func(q float64) int64 {
		idx := int(float64(len(rows)-1) * q)
		if idx < 0 {
			idx = 0
		}
		if idx >= len(rows) {
			idx = len(rows) - 1
		}
		return rows[idx]
	}
	return at(0.50), at(0.95)
}

func (p *ProviderTelemetry) Diagnostics() []ProviderRequestDiagnostics {
	if p == nil {
		return nil
	}
	now := time.Now()
	p.mu.Lock()
	defer p.mu.Unlock()
	out := make([]ProviderRequestDiagnostics, 0, len(p.provider))
	for _, st := range p.provider {
		pruneProviderRecent(now, st)
		avg := int64(0)
		completed := st.successes + st.errors
		successPct := 0.0
		if completed > 0 {
			avg = st.totalLatency / completed
			successPct = float64(st.successes) / float64(completed) * 100
		}
		p50, p95 := latencyPercentiles(st.latencies)
		budget := providerBudgetPerMinute(st.provider)
		remaining, utilization := 0, 0
		budgetState := "ENTITLEMENT DEPENDENT"
		if budget > 0 {
			remaining = budget - len(st.recent)
			if remaining < 0 {
				remaining = 0
			}
			utilization = len(st.recent) * 100 / budget
			budgetState = "AVAILABLE"
			if utilization >= 80 {
				budgetState = "PRESSURE"
			}
			if len(st.recent) >= budget {
				budgetState = "EXHAUSTED"
			}
		}
		cooldown := int64(0)
		if !st.cooldownUntil.IsZero() && now.Before(st.cooldownUntil) {
			cooldown = st.cooldownUntil.UnixMilli()
			budgetState = "RATE LIMITED"
		}
		out = append(out, ProviderRequestDiagnostics{
			Provider: st.provider, Requests: st.requests, RequestsLastMin: len(st.recent), Successes: st.successes,
			Errors: st.errors, RateLimited: st.rateLimited, InFlight: st.inFlight, PeakInFlight: st.peakInFlight,
			LastLatencyMs: st.lastLatency, AverageLatencyMs: avg, P50LatencyMs: p50, P95LatencyMs: p95, SuccessPct: successPct, LastRequestAt: st.lastRequestAt, LastError: st.lastError,
			BudgetPerMinute: budget, BudgetRemaining: remaining, BudgetUtilizationPct: utilization, BudgetState: budgetState,
			CooldownUntil: cooldown, BudgetShed: st.budgetShed,
		})
	}
	sort.Slice(out, func(i, j int) bool { return out[i].Provider < out[j].Provider })
	return out
}

func (e *Engine) providerGetJSONTier(ctx context.Context, provider string, tier WorkTier, client *http.Client, rawURL string, headers map[string]string, out any) error {
	if e == nil {
		return getJSON(ctx, client, rawURL, headers, out)
	}
	tier = workTierFromContext(ctx, tier)
	if ok, reason := e.providerTelemetry.Allow(provider, tier); !ok {
		return fmt.Errorf("%s deferred: %s", workTierLabel(tier), reason)
	}
	release, ok := e.workload.AcquireTier(ctx, "provider-rest", tier)
	if !ok {
		if err := ctx.Err(); err != nil {
			return err
		}
		return fmt.Errorf("provider request rejected by bounded workload budget")
	}
	defer release()
	done := e.providerTelemetry.begin(provider)
	err := getJSON(ctx, client, rawURL, headers, out)
	done(err)
	return err
}

func (e *Engine) shouldShedTier(tier WorkTier) bool {
	if e == nil || e.workload == nil {
		return false
	}
	return e.workload.ShouldShed(tier)
}

func (e *Engine) shouldShedBackground() bool {
	return e.shouldShedTier(WorkTierBackground)
}

type HTTPRuntimeDiagnostics struct {
	Requests                 int64  `json:"requests"`
	InFlight                 int    `json:"inFlight"`
	LastLatencyMs            int64  `json:"lastLatencyMs,omitempty"`
	AverageLatencyMs         int64  `json:"averageLatencyMs,omitempty"`
	InteractiveP95Ms         int64  `json:"interactiveP95Ms,omitempty"`
	InteractiveMaxMs         int64  `json:"interactiveMaxMs,omitempty"`
	SlowInteractive          int64  `json:"slowInteractive"`
	LastPath                 string `json:"lastPath,omitempty"`
	HostedMutationAllowed    int64  `json:"hostedMutationAllowed"`
	HostedMutationThrottled  int64  `json:"hostedMutationThrottled"`
	HostedExpensiveAllowed   int64  `json:"hostedExpensiveAllowed"`
	HostedExpensiveThrottled int64  `json:"hostedExpensiveThrottled"`
}

type RequestTelemetry struct {
	mu                       sync.Mutex
	requests                 int64
	inFlight                 int
	totalMs                  int64
	lastMs                   int64
	lastPath                 string
	slow                     int64
	interactive              []int64
	hostedMutationAllowed    int64
	hostedMutationThrottled  int64
	hostedExpensiveAllowed   int64
	hostedExpensiveThrottled int64
}

func NewRequestTelemetry() *RequestTelemetry { return &RequestTelemetry{} }

func isExpensiveAPIPath(path string) bool {
	for _, heavy := range []string{"refresh", "scan", "generate", "maintenance", "integrity", "provider/test", "pre-market", "market-open", "catalyst", "stream-reconnect"} {
		if strings.Contains(path, heavy) {
			return true
		}
	}
	return false
}

func isInteractiveAPIPath(path string) bool {
	if !strings.HasPrefix(path, "/api/") || path == "/api/events" {
		return false
	}
	return !isExpensiveAPIPath(path)
}

func (t *RequestTelemetry) begin(path string) func() {
	if t == nil {
		return func() {}
	}
	start := time.Now()
	interactive := isInteractiveAPIPath(path)
	t.mu.Lock()
	t.requests++
	t.inFlight++
	t.lastPath = path
	t.mu.Unlock()
	return func() {
		ms := time.Since(start).Milliseconds()
		t.mu.Lock()
		defer t.mu.Unlock()
		if t.inFlight > 0 {
			t.inFlight--
		}
		t.totalMs += ms
		t.lastMs = ms
		if interactive {
			if ms > 250 {
				t.slow++
			}
			t.interactive = append(t.interactive, ms)
			if len(t.interactive) > 128 {
				t.interactive = append([]int64(nil), t.interactive[len(t.interactive)-128:]...)
			}
		}
	}
}

func (t *RequestTelemetry) RecordHostedQuota(expensive, allowed bool) {
	if t == nil {
		return
	}
	t.mu.Lock()
	defer t.mu.Unlock()
	if allowed {
		t.hostedMutationAllowed++
		if expensive {
			t.hostedExpensiveAllowed++
		}
		return
	}
	t.hostedMutationThrottled++
	if expensive {
		t.hostedExpensiveThrottled++
	}
}

func (t *RequestTelemetry) Diagnostics() HTTPRuntimeDiagnostics {
	if t == nil {
		return HTTPRuntimeDiagnostics{}
	}
	t.mu.Lock()
	defer t.mu.Unlock()
	avg := int64(0)
	if t.requests > 0 {
		avg = t.totalMs / t.requests
	}
	values := append([]int64(nil), t.interactive...)
	sort.Slice(values, func(i, j int) bool { return values[i] < values[j] })
	p95, max := int64(0), int64(0)
	if len(values) > 0 {
		idx := (95*len(values)+99)/100 - 1
		if idx < 0 {
			idx = 0
		}
		if idx >= len(values) {
			idx = len(values) - 1
		}
		p95 = values[idx]
		max = values[len(values)-1]
	}
	return HTTPRuntimeDiagnostics{
		Requests:                 t.requests,
		InFlight:                 t.inFlight,
		LastLatencyMs:            t.lastMs,
		AverageLatencyMs:         avg,
		InteractiveP95Ms:         p95,
		InteractiveMaxMs:         max,
		SlowInteractive:          t.slow,
		LastPath:                 t.lastPath,
		HostedMutationAllowed:    t.hostedMutationAllowed,
		HostedMutationThrottled:  t.hostedMutationThrottled,
		HostedExpensiveAllowed:   t.hostedExpensiveAllowed,
		HostedExpensiveThrottled: t.hostedExpensiveThrottled,
	}
}

func (a *Application) observeHTTP(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if a == nil || a.httpTelemetry == nil || !strings.HasPrefix(r.URL.Path, "/api/") || r.URL.Path == "/api/events" {
			next.ServeHTTP(w, r)
			return
		}
		done := a.httpTelemetry.begin(r.URL.Path)
		defer done()
		next.ServeHTTP(w, r)
	})
}
