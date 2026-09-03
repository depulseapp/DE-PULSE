package main

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"sort"
	"strings"
	"sync"
	"time"
)

const (
	providerUsefulnessFeatureKey     = "provider-semantic-usefulness"
	providerUsefulnessFeatureVersion = "v1"
	providerUsefulnessMinCrossSource = 5
	providerUsefulnessSeenLimit      = 512
	providerUsefulnessPersistEveryMs = int64((30 * time.Second) / time.Millisecond)
)

type ProviderUsefulnessDiagnostic struct {
	Provider            string   `json:"provider"`
	State               string   `json:"state"` // INSUFFICIENT or OBSERVING
	EligibleSamples     int64    `json:"eligibleSamples"`
	CrossSourceSamples  int64    `json:"crossSourceSamples"`
	AgreementSamples    int64    `json:"agreementSamples"`
	ConflictSamples     int64    `json:"conflictSamples"`
	SingleSourceSamples int64    `json:"singleSourceSamples"`
	CanonicalSelections int64    `json:"canonicalSelections"`
	ExcludedSamples     int64    `json:"excludedSamples"`
	AgreementPct        *float64 `json:"agreementPct,omitempty"`
	RoutingImpact       string   `json:"routingImpact"`
	LastObservedAt      int64    `json:"lastObservedAt,omitempty"`
}

// ProviderOperationalScorecard is a read-only composition of measurements
// already owned by Smart Router v2, Data Freshness, provider transport
// telemetry, live-subscription budgets, provider rights and semantic
// usefulness. It is deliberately not an input to routing. Pointer-valued
// metrics preserve the distinction between a measured zero and no evidence.
type ProviderOperationalScorecard struct {
	Provider                   string   `json:"provider"`
	State                      string   `json:"state"`
	Configured                 bool     `json:"configured"`
	Capabilities               []string `json:"capabilities,omitempty"`
	CoverageDatasets           []string `json:"coverageDatasets,omitempty"`
	ServingDatasets            []string `json:"servingDatasets,omitempty"`
	HealthMeasurementState     string   `json:"healthMeasurementState"`
	HealthStates               []string `json:"healthStates,omitempty"`
	FreshnessMeasurementState  string   `json:"freshnessMeasurementState"`
	FreshnessStates            []string `json:"freshnessStates,omitempty"`
	TransportMeasurementState  string   `json:"transportMeasurementState"`
	CompletedRequests          int64    `json:"completedRequests"`
	SuccessPct                 *float64 `json:"successPct,omitempty"`
	P50LatencyMs               *int64   `json:"p50LatencyMs,omitempty"`
	P95LatencyMs               *int64   `json:"p95LatencyMs,omitempty"`
	RateLimited                int64    `json:"rateLimited"`
	HeadroomMeasurementState   string   `json:"headroomMeasurementState"`
	RequestBudgetRemaining     *int     `json:"requestBudgetRemaining,omitempty"`
	LiveSubscriptionAvailable  *int     `json:"liveSubscriptionAvailable,omitempty"`
	RightsMeasurementState     string   `json:"rightsMeasurementState"`
	RightsReviewState          string   `json:"rightsReviewState"`
	RightsEvidenceBound        bool     `json:"rightsEvidenceBound"`
	CommercialReadinessState   string   `json:"commercialReadinessState"`
	CostClass                  string   `json:"costClass"`
	CostMeasurementState       string   `json:"costMeasurementState"`
	ObservedCostUSD            *float64 `json:"observedCostUsd,omitempty"`
	UsefulnessMeasurementState string   `json:"usefulnessMeasurementState"`
	CrossSourceSamples         int64    `json:"crossSourceSamples"`
	AgreementPct               *float64 `json:"agreementPct,omitempty"`
	RoutingImpact              string   `json:"routingImpact"`
	UpdatedAt                  int64    `json:"updatedAt"`
}

type providerOperationalScorecardBuilder struct {
	row           ProviderOperationalScorecard
	healthSeen    bool
	freshSeen     bool
	transportSeen bool
	headroomSeen  bool
	degraded      bool
}

func providerScorecardSameProvider(left, right string) bool {
	left = providerKey(left)
	right = providerKey(right)
	return left != "" && right != "" && (left == right || strings.HasPrefix(left, right+"-") || strings.HasPrefix(right, left+"-"))
}

func providerScorecardAppendUnique(values []string, value string) []string {
	value = strings.TrimSpace(value)
	if value == "" {
		return values
	}
	for _, existing := range values {
		if strings.EqualFold(existing, value) {
			return values
		}
	}
	return append(values, value)
}

func providerScorecardInt64Pointer(value int64) *int64     { v := value; return &v }
func providerScorecardIntPointer(value int) *int           { v := value; return &v }
func providerScorecardFloatPointer(value float64) *float64 { v := value; return &v }

// buildProviderOperationalScorecards never scores or orders providers. It only
// projects canonical measurements and explicit unknowns for privileged
// operational consumers.
func buildProviderOperationalScorecards(registrations []ProviderRegistration, settings Settings, secrets Secrets, router ProviderRouterSnapshot, freshness []FreshnessDiagnostic, requests []ProviderRequestDiagnostics, usefulness []ProviderUsefulnessDiagnostic, subscriptions []LiveSubscriptionBudgetDiagnostics, now int64) []ProviderOperationalScorecard {
	builders := map[string]*providerOperationalScorecardBuilder{}
	ensure := func(provider string) *providerOperationalScorecardBuilder {
		key := providerKey(provider)
		if key == "" {
			return nil
		}
		if builders[key] == nil {
			rights := providerDataRightsMetadata(provider)
			rightsState := "UNBOUND"
			if rights.EvidenceBound {
				rightsState = "BOUND"
			}
			builders[key] = &providerOperationalScorecardBuilder{row: ProviderOperationalScorecard{
				Provider: provider, State: "UNOBSERVED", HealthMeasurementState: "UNKNOWN",
				FreshnessMeasurementState: "UNKNOWN", TransportMeasurementState: "UNKNOWN",
				HeadroomMeasurementState: "UNKNOWN", RightsMeasurementState: rightsState,
				RightsReviewState: defaultString(rights.ReviewState, "UNKNOWN"), RightsEvidenceBound: rights.EvidenceBound,
				CommercialReadinessState: defaultString(rights.CommercialReadiness.State, "UNKNOWN"),
				CostClass:                providerCostFromRegistration(provider), CostMeasurementState: "DECLARED_CLASS_ONLY",
				UsefulnessMeasurementState: "UNKNOWN", RoutingImpact: "OBSERVABILITY_ONLY", UpdatedAt: now,
			}}
		}
		return builders[key]
	}

	for _, registration := range registrations {
		builder := ensure(registration.Name)
		if builder == nil {
			continue
		}
		builder.row.Configured = registration.Configured != nil && registration.Configured(settings, secrets)
		builder.row.CostClass = defaultString(registration.CostClass, builder.row.CostClass)
		for _, route := range registration.Routes {
			builder.row.Capabilities = providerScorecardAppendUnique(builder.row.Capabilities, route.Capability)
			builder.row.CoverageDatasets = providerScorecardAppendUnique(builder.row.CoverageDatasets, route.Dataset)
		}
		for _, diagnostic := range registration.Diagnostics {
			builder.row.Capabilities = providerScorecardAppendUnique(builder.row.Capabilities, diagnostic.Capability)
		}
	}

	for _, route := range router.Routes {
		for _, hop := range route.Route {
			builder := ensure(hop.Provider)
			if builder == nil {
				continue
			}
			builder.row.Configured = builder.row.Configured || hop.Configured
			builder.row.CoverageDatasets = providerScorecardAppendUnique(builder.row.CoverageDatasets, route.Dataset)
			builder.row.HealthStates = providerScorecardAppendUnique(builder.row.HealthStates, route.Dataset+": "+defaultString(hop.Health, "UNKNOWN"))
			if strings.EqualFold(route.Serving, hop.Provider) {
				builder.row.ServingDatasets = providerScorecardAppendUnique(builder.row.ServingDatasets, route.Dataset)
			}
			if hop.Attempts > 0 || hop.LastSuccess > 0 || hop.LastFailure > 0 {
				builder.healthSeen = true
				builder.row.HealthMeasurementState = "MEASURED"
			} else if builder.row.HealthMeasurementState == "UNKNOWN" && hop.Configured {
				builder.row.HealthMeasurementState = "DECLARED_ONLY"
			}
			health := strings.ToUpper(strings.TrimSpace(hop.Health + " " + hop.Circuit))
			if strings.Contains(health, "DEGRADED") || strings.Contains(health, "ERROR") || strings.Contains(health, "OPEN") || strings.Contains(health, "RATE LIMITED") || strings.Contains(health, "BLOCKED") {
				builder.degraded = true
			}
		}
	}

	for _, diagnostic := range freshness {
		for _, builder := range builders {
			if !providerScorecardSameProvider(builder.row.Provider, diagnostic.Provider) {
				continue
			}
			builder.freshSeen = true
			builder.row.FreshnessMeasurementState = "MEASURED"
			builder.row.FreshnessStates = providerScorecardAppendUnique(builder.row.FreshnessStates, diagnostic.Dataset+": "+defaultString(diagnostic.State, "UNKNOWN"))
			state := strings.ToUpper(strings.TrimSpace(diagnostic.State))
			if state == "STALE" || state == "ERROR" || state == "UNAVAILABLE" {
				builder.degraded = true
			}
		}
	}

	for _, request := range requests {
		builder := ensure(request.Provider)
		if builder == nil {
			continue
		}
		completed := request.Successes + request.Errors
		builder.row.CompletedRequests = completed
		builder.row.RateLimited = request.RateLimited
		if completed > 0 {
			builder.transportSeen = true
			builder.row.TransportMeasurementState = "MEASURED"
			builder.row.SuccessPct = providerScorecardFloatPointer(request.SuccessPct)
			builder.row.P50LatencyMs = providerScorecardInt64Pointer(request.P50LatencyMs)
			builder.row.P95LatencyMs = providerScorecardInt64Pointer(request.P95LatencyMs)
		}
		if request.BudgetPerMinute > 0 {
			builder.headroomSeen = true
			builder.row.HeadroomMeasurementState = "MEASURED"
			builder.row.RequestBudgetRemaining = providerScorecardIntPointer(request.BudgetRemaining)
		}
		if request.RateLimited > 0 || request.Errors > request.Successes {
			builder.degraded = true
		}
	}

	for _, subscription := range subscriptions {
		for _, builder := range builders {
			if !providerScorecardSameProvider(builder.row.Provider, subscription.Provider) {
				continue
			}
			builder.headroomSeen = true
			builder.row.HeadroomMeasurementState = "MEASURED"
			builder.row.LiveSubscriptionAvailable = providerScorecardIntPointer(subscription.Available)
			if subscription.Saturated {
				builder.degraded = true
			}
		}
	}

	for _, diagnostic := range usefulness {
		builder := ensure(diagnostic.Provider)
		if builder == nil {
			continue
		}
		builder.row.UsefulnessMeasurementState = defaultString(diagnostic.State, "UNKNOWN")
		builder.row.CrossSourceSamples = diagnostic.CrossSourceSamples
		if diagnostic.AgreementPct != nil {
			builder.row.AgreementPct = providerScorecardFloatPointer(*diagnostic.AgreementPct)
		}
	}

	out := make([]ProviderOperationalScorecard, 0, len(builders))
	for _, builder := range builders {
		sort.Strings(builder.row.Capabilities)
		sort.Strings(builder.row.CoverageDatasets)
		sort.Strings(builder.row.ServingDatasets)
		sort.Strings(builder.row.HealthStates)
		sort.Strings(builder.row.FreshnessStates)
		measured := 0
		for _, present := range []bool{builder.healthSeen, builder.freshSeen, builder.transportSeen, builder.headroomSeen, builder.row.AgreementPct != nil} {
			if present {
				measured++
			}
		}
		if measured > 0 {
			builder.row.State = "PARTIAL"
		}
		if measured == 5 {
			builder.row.State = "MEASURED"
		}
		if builder.degraded {
			builder.row.State = "DEGRADED"
		}
		out = append(out, builder.row)
	}
	sort.Slice(out, func(i, j int) bool { return strings.ToUpper(out[i].Provider) < strings.ToUpper(out[j].Provider) })
	return out
}

type providerUsefulnessAggregate struct {
	EligibleSamples     int64 `json:"eligibleSamples"`
	CrossSourceSamples  int64 `json:"crossSourceSamples"`
	AgreementSamples    int64 `json:"agreementSamples"`
	ConflictSamples     int64 `json:"conflictSamples"`
	SingleSourceSamples int64 `json:"singleSourceSamples"`
	CanonicalSelections int64 `json:"canonicalSelections"`
	ExcludedSamples     int64 `json:"excludedSamples"`
	LastObservedAt      int64 `json:"lastObservedAt,omitempty"`
}

type providerUsefulnessPersistedState struct {
	Providers map[string]providerUsefulnessAggregate `json:"providers"`
	SeenOrder []string                               `json:"seenOrder,omitempty"`
	UpdatedAt int64                                  `json:"updatedAt,omitempty"`
}

type providerUsefulnessTracker struct {
	mu            sync.Mutex
	persistence   *PersistenceManager
	restored      bool
	providers     map[string]providerUsefulnessAggregate
	seen          map[string]bool
	seenOrder     []string
	dirty         bool
	lastPersistAt int64
}

var canonicalProviderUsefulness = newProviderUsefulnessTracker()

func newProviderUsefulnessTracker() *providerUsefulnessTracker {
	return &providerUsefulnessTracker{
		providers: map[string]providerUsefulnessAggregate{},
		seen:      map[string]bool{},
	}
}

func providerUsefulnessProvider(provider string) string {
	provider = strings.TrimSpace(provider)
	if provider == "" || strings.EqualFold(provider, "Unknown") {
		return ""
	}
	return provider
}

func providerUsefulnessDecisionSignature(row ProviderReconciliationDecision) string {
	parts := []string{
		strings.ToUpper(strings.TrimSpace(row.Dataset)),
		normalizeSymbol(row.Symbol),
		strings.ToUpper(strings.TrimSpace(row.State)),
		strings.ToUpper(strings.TrimSpace(row.CanonicalProvider)),
		jsonFloat(row.CanonicalValue),
	}
	observations := append([]ProviderQuoteObservation(nil), row.Observations...)
	sort.Slice(observations, func(i, j int) bool {
		a := strings.ToUpper(strings.TrimSpace(observations[i].Provider)) + "|" + strings.TrimSpace(observations[i].Source)
		b := strings.ToUpper(strings.TrimSpace(observations[j].Provider)) + "|" + strings.TrimSpace(observations[j].Source)
		if a != b {
			return a < b
		}
		return reconciliationObservationTime(observations[i]) < reconciliationObservationTime(observations[j])
	})
	for _, observation := range observations {
		parts = append(parts, strings.Join([]string{
			strings.ToUpper(strings.TrimSpace(observation.Provider)),
			strings.TrimSpace(observation.Source),
			jsonFloat(observation.Price),
			jsonInt64(reconciliationObservationTime(observation)),
		}, "|"))
	}
	sum := sha256.Sum256([]byte(strings.Join(parts, "\n")))
	return hex.EncodeToString(sum[:])
}

func jsonFloat(v float64) string {
	raw, _ := json.Marshal(v)
	return string(raw)
}

func jsonInt64(v int64) string {
	raw, _ := json.Marshal(v)
	return string(raw)
}

func (t *providerUsefulnessTracker) bindPersistence(p *PersistenceManager) {
	if t == nil || p == nil {
		return
	}
	t.mu.Lock()
	if t.persistence == p && t.restored {
		t.mu.Unlock()
		return
	}
	if t.persistence != p {
		t.persistence = p
		t.providers = map[string]providerUsefulnessAggregate{}
		t.seen = map[string]bool{}
		t.seenOrder = nil
		t.restored = false
		t.dirty = false
		t.lastPersistAt = 0
	}
	t.mu.Unlock()

	state, ok := loadProviderUsefulnessState(p)
	t.mu.Lock()
	defer t.mu.Unlock()
	if t.persistence != p || t.restored {
		return
	}
	if ok {
		t.restoreStateLocked(state)
	}
	t.restored = true
}

func loadProviderUsefulnessState(p *PersistenceManager) (providerUsefulnessPersistedState, bool) {
	if p == nil || p.backend == nil {
		return providerUsefulnessPersistedState{}, false
	}
	backend, ok := p.backend.(persistenceArchiveBackend)
	if !ok {
		return providerUsefulnessPersistedState{}, false
	}
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	archive, err := backend.ExportPersistenceArchive(ctx)
	if err != nil {
		return providerUsefulnessPersistedState{}, false
	}
	for _, feature := range archive.Features {
		if feature.FeatureKey != providerUsefulnessFeatureKey || feature.FeatureVersion != providerUsefulnessFeatureVersion {
			continue
		}
		var state providerUsefulnessPersistedState
		if json.Unmarshal(feature.Payload, &state) == nil {
			return state, true
		}
	}
	return providerUsefulnessPersistedState{}, false
}

func (t *providerUsefulnessTracker) restoreStateLocked(state providerUsefulnessPersistedState) {
	if state.Providers == nil {
		state.Providers = map[string]providerUsefulnessAggregate{}
	}
	t.providers = state.Providers
	t.seen = map[string]bool{}
	t.seenOrder = nil
	start := 0
	if len(state.SeenOrder) > providerUsefulnessSeenLimit {
		start = len(state.SeenOrder) - providerUsefulnessSeenLimit
	}
	for _, signature := range state.SeenOrder[start:] {
		signature = strings.TrimSpace(signature)
		if signature == "" || t.seen[signature] {
			continue
		}
		t.seen[signature] = true
		t.seenOrder = append(t.seenOrder, signature)
	}
	t.lastPersistAt = state.UpdatedAt
	t.dirty = false
}

func (t *providerUsefulnessTracker) markSeenLocked(signature string) bool {
	if signature == "" || t.seen[signature] {
		return false
	}
	t.seen[signature] = true
	t.seenOrder = append(t.seenOrder, signature)
	for len(t.seenOrder) > providerUsefulnessSeenLimit {
		oldest := t.seenOrder[0]
		t.seenOrder = t.seenOrder[1:]
		delete(t.seen, oldest)
	}
	return true
}

func (t *providerUsefulnessTracker) updateProviderLocked(provider string, now int64, apply func(*providerUsefulnessAggregate)) {
	provider = providerUsefulnessProvider(provider)
	if provider == "" {
		return
	}
	row := t.providers[provider]
	apply(&row)
	if now > row.LastObservedAt {
		row.LastObservedAt = now
	}
	t.providers[provider] = row
}

func uniqueUsefulnessProviders(rows []ProviderQuoteObservation) []string {
	seen := map[string]bool{}
	out := make([]string, 0, len(rows))
	for _, row := range rows {
		provider := providerUsefulnessProvider(row.Provider)
		if provider == "" || seen[strings.ToUpper(provider)] {
			continue
		}
		seen[strings.ToUpper(provider)] = true
		out = append(out, provider)
	}
	return out
}

func (t *providerUsefulnessTracker) observe(rows []ProviderReconciliationDecision, now int64) bool {
	if t == nil || len(rows) == 0 {
		return false
	}
	if now <= 0 {
		now = time.Now().UnixMilli()
	}
	t.mu.Lock()
	defer t.mu.Unlock()
	changed := false
	for _, decision := range rows {
		signature := providerUsefulnessDecisionSignature(decision)
		if !t.markSeenLocked(signature) {
			continue
		}
		changed = true
		providers := uniqueUsefulnessProviders(decision.Observations)
		state := strings.ToUpper(strings.TrimSpace(decision.State))
		switch state {
		case "AGREED", "CONFLICT":
			for _, provider := range providers {
				t.updateProviderLocked(provider, now, func(row *providerUsefulnessAggregate) {
					row.EligibleSamples++
					row.CrossSourceSamples++
					if state == "AGREED" {
						row.AgreementSamples++
					} else {
						row.ConflictSamples++
					}
				})
			}
		case "SINGLE SOURCE":
			for _, provider := range providers {
				t.updateProviderLocked(provider, now, func(row *providerUsefulnessAggregate) {
					row.EligibleSamples++
					row.SingleSourceSamples++
				})
			}
		case "STALE":
			// Reconciliation deliberately drops stale, non-contemporaneous and
			// invalid/future observations. Attribute an exclusion only when the
			// canonical decision still identifies a provider; never guess a source.
			t.updateProviderLocked(decision.CanonicalProvider, now, func(row *providerUsefulnessAggregate) {
				row.ExcludedSamples++
			})
		}
		canonical := providerUsefulnessProvider(decision.CanonicalProvider)
		if canonical != "" && state != "STALE" {
			for _, provider := range providers {
				if strings.EqualFold(provider, canonical) {
					t.updateProviderLocked(provider, now, func(row *providerUsefulnessAggregate) {
						row.CanonicalSelections++
					})
					break
				}
			}
		}
	}
	if changed {
		t.dirty = true
	}
	return changed
}

func (t *providerUsefulnessTracker) diagnostics() []ProviderUsefulnessDiagnostic {
	if t == nil {
		return nil
	}
	t.mu.Lock()
	defer t.mu.Unlock()
	out := make([]ProviderUsefulnessDiagnostic, 0, len(t.providers))
	for provider, aggregate := range t.providers {
		state := "INSUFFICIENT"
		var agreementPct *float64
		if aggregate.CrossSourceSamples >= providerUsefulnessMinCrossSource {
			state = "OBSERVING"
			pct := float64(aggregate.AgreementSamples) / float64(aggregate.CrossSourceSamples) * 100
			agreementPct = &pct
		}
		out = append(out, ProviderUsefulnessDiagnostic{
			Provider: provider, State: state,
			EligibleSamples: aggregate.EligibleSamples, CrossSourceSamples: aggregate.CrossSourceSamples,
			AgreementSamples: aggregate.AgreementSamples, ConflictSamples: aggregate.ConflictSamples,
			SingleSourceSamples: aggregate.SingleSourceSamples, CanonicalSelections: aggregate.CanonicalSelections,
			ExcludedSamples: aggregate.ExcludedSamples, AgreementPct: agreementPct,
			RoutingImpact: "ADVISORY_ONLY", LastObservedAt: aggregate.LastObservedAt,
		})
	}
	sort.Slice(out, func(i, j int) bool { return strings.ToUpper(out[i].Provider) < strings.ToUpper(out[j].Provider) })
	return out
}

func (t *providerUsefulnessTracker) stateLocked(now int64) providerUsefulnessPersistedState {
	providers := make(map[string]providerUsefulnessAggregate, len(t.providers))
	for provider, aggregate := range t.providers {
		providers[provider] = aggregate
	}
	return providerUsefulnessPersistedState{
		Providers: providers,
		SeenOrder: append([]string(nil), t.seenOrder...),
		UpdatedAt: now,
	}
}

func (t *providerUsefulnessTracker) persist(p *PersistenceManager, now int64) {
	if t == nil || p == nil {
		return
	}
	if now <= 0 {
		now = time.Now().UnixMilli()
	}
	t.mu.Lock()
	if !t.dirty || (t.lastPersistAt > 0 && now-t.lastPersistAt < providerUsefulnessPersistEveryMs) {
		t.mu.Unlock()
		return
	}
	state := t.stateLocked(now)
	t.dirty = false
	t.lastPersistAt = now
	t.mu.Unlock()
	raw, err := json.Marshal(state)
	if err != nil {
		t.mu.Lock()
		t.dirty = true
		t.mu.Unlock()
		return
	}
	sum := sha256.Sum256(raw)
	p.EnqueueIntelligence(PersistenceIntelligenceBatch{Features: []DerivedFeatureRecord{{
		Symbol: "", FeatureKey: providerUsefulnessFeatureKey, FeatureVersion: providerUsefulnessFeatureVersion,
		AsOf: now, SourceHash: hex.EncodeToString(sum[:]), Payload: raw,
	}}})
}
