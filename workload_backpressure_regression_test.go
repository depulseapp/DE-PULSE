package main

import (
	"context"
	"testing"
	"time"
)

func providerWorkDiag(w *WorkloadController) WorkClassDiagnostics {
	for _, row := range w.Diagnostics() {
		if row.Class == "provider-rest" {
			return row
		}
	}
	return WorkClassDiagnostics{}
}

func waitProviderQueue(t *testing.T, w *WorkloadController, want int) WorkClassDiagnostics {
	t.Helper()
	deadline := time.Now().Add(750 * time.Millisecond)
	for time.Now().Before(deadline) {
		d := providerWorkDiag(w)
		if d.Queued >= want {
			return d
		}
		time.Sleep(5 * time.Millisecond)
	}
	d := providerWorkDiag(w)
	t.Fatalf("provider queue never reached %d: %+v", want, d)
	return d
}

func TestV171ReservedCapacityProtectsTierZeroAndOne(t *testing.T) {
	w := NewWorkloadController()
	d := providerWorkDiag(w)
	lowLimit := d.Capacity - d.ReservedCritical
	if d.ReservedCritical < 1 || lowLimit < 1 {
		t.Fatalf("provider critical reserve not configured: %+v", d)
	}

	lowReleases := make([]func(), 0, lowLimit)
	for i := 0; i < lowLimit; i++ {
		release, ok := w.TryAcquireTier("provider-rest", WorkTierBroadDiscovery)
		if !ok {
			t.Fatalf("Tier 3 should fill normal capacity slot %d/%d", i+1, lowLimit)
		}
		lowReleases = append(lowReleases, release)
	}
	if release, ok := w.TryAcquireTier("provider-rest", WorkTierBroadDiscovery); ok {
		release()
		t.Fatal("Tier 3 consumed capacity reserved for critical/actionable work")
	}

	criticalReleases := []func(){}
	for _, tier := range []WorkTier{WorkTierMarketCritical, WorkTierUserActionable} {
		release, ok := w.TryAcquireTier("provider-rest", tier)
		if !ok {
			t.Fatalf("tier %d could not consume protected reserve", tier)
		}
		criticalReleases = append(criticalReleases, release)
	}
	d = providerWorkDiag(w)
	if d.InFlight != d.Capacity || d.InFlightByTier[WorkTierBroadDiscovery] != lowLimit {
		t.Fatalf("tier allocation diagnostics incorrect: %+v", d)
	}
	if d.RejectedByTier[WorkTierBroadDiscovery] < 1 || d.Shed < 1 {
		t.Fatalf("shed/rejection pressure not observable: %+v", d)
	}
	for _, release := range criticalReleases {
		release()
	}
	for _, release := range lowReleases {
		release()
	}
}

func TestV171ProviderQueueIsHardBounded(t *testing.T) {
	w := NewWorkloadController()
	d := providerWorkDiag(w)
	if d.ReservedCriticalQueue < 1 || d.ReservedCriticalQueue >= d.MaxQueue {
		t.Fatalf("provider critical queue reserve not configured: %+v", d)
	}
	holds := make([]func(), 0, d.Capacity)
	for i := 0; i < d.Capacity; i++ {
		release, ok := w.TryAcquireTier("provider-rest", WorkTierMarketCritical)
		if !ok {
			t.Fatalf("could not fill provider slot %d", i)
		}
		holds = append(holds, release)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	done := make(chan struct{}, d.MaxQueue)
	lowQueueLimit := d.MaxQueue - d.ReservedCriticalQueue
	for i := 0; i < lowQueueLimit; i++ {
		go func() {
			release, ok := w.AcquireTier(ctx, "provider-rest", WorkTierBackground)
			if ok {
				release()
			}
			done <- struct{}{}
		}()
	}
	d = waitProviderQueue(t, w, lowQueueLimit)
	if d.Queued != lowQueueLimit {
		t.Fatalf("optional queue exceeded its bounded non-critical share: %+v", d)
	}

	rejectCtx, rejectCancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer rejectCancel()
	if release, ok := w.AcquireTier(rejectCtx, "provider-rest", WorkTierBackground); ok {
		release()
		t.Fatal("optional work consumed queue headroom reserved for critical recovery")
	}
	afterLowReject := providerWorkDiag(w)
	if afterLowReject.Queued != lowQueueLimit || afterLowReject.RejectedByTier[WorkTierBackground] < 1 || afterLowReject.Shed < 1 {
		t.Fatalf("optional queue shedding not observable: %+v", afterLowReject)
	}

	for i := 0; i < d.ReservedCriticalQueue; i++ {
		go func() {
			release, ok := w.AcquireTier(ctx, "provider-rest", WorkTierMarketCritical)
			if ok {
				release()
			}
			done <- struct{}{}
		}()
	}
	full := waitProviderQueue(t, w, d.MaxQueue)
	if full.Queued != full.MaxQueue || full.QueuedByTier[WorkTierMarketCritical] != full.ReservedCriticalQueue {
		t.Fatalf("critical work could not use protected queue headroom: %+v", full)
	}

	extraCriticalCtx, extraCriticalCancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer extraCriticalCancel()
	if release, ok := w.AcquireTier(extraCriticalCtx, "provider-rest", WorkTierMarketCritical); ok {
		release()
		t.Fatal("total provider queue exceeded hard maxQueue bound")
	}
	after := providerWorkDiag(w)
	if after.Queued > after.MaxQueue || after.Rejected < 2 {
		t.Fatalf("bounded queue rejection not observable: %+v", after)
	}

	cancel()
	for _, release := range holds {
		release()
	}
	for i := 0; i < d.MaxQueue; i++ {
		select {
		case <-done:
		case <-time.After(time.Second):
			t.Fatal("queued worker did not cancel cleanly")
		}
	}
}

func TestV171LowPriorityShedsBeforeActionable(t *testing.T) {
	w := NewWorkloadController()
	d := providerWorkDiag(w)
	holds := []func(){}
	for i := 0; i < d.Capacity-d.ReservedCritical; i++ {
		release, ok := w.TryAcquireTier("provider-rest", WorkTierBroadDiscovery)
		if !ok {
			t.Fatal("could not establish normal-pool pressure")
		}
		holds = append(holds, release)
	}
	if !w.ShouldShed(WorkTierBackground) || !w.ShouldShed(WorkTierBroadDiscovery) {
		t.Fatal("lower-priority work should shed when normal provider capacity is exhausted")
	}
	if w.ShouldShed(WorkTierUserActionable) || w.ShouldShed(WorkTierMarketCritical) {
		t.Fatal("Tier 0/1 work must not be shed by local backpressure policy")
	}
	release, ok := w.TryAcquireTier("provider-rest", WorkTierUserActionable)
	if !ok {
		t.Fatal("actionable work lost protected capacity")
	}
	release()
	for _, hold := range holds {
		hold()
	}
}

func TestADAPTDataHealthFreshnessRecoveryPriorityUsesProtectedWorkTiers(t *testing.T) {
	w := NewWorkloadController()
	for priority, want := range map[int]WorkTier{
		1: WorkTierMarketCritical,
		2: WorkTierUserActionable,
		3: WorkTierBroadDiscovery,
		4: WorkTierBroadDiscovery,
	} {
		ctx, tier := freshnessRecoveryWorkContext(context.Background(), priority)
		if tier != want || workTierFromContext(ctx, WorkTierBackground) != want {
			t.Fatalf("freshness recovery priority %d mapped to %v, want %v", priority, tier, want)
		}
	}

	d := providerWorkDiag(w)
	holds := make([]func(), 0, d.Capacity-d.ReservedCritical)
	for i := 0; i < d.Capacity-d.ReservedCritical; i++ {
		release, ok := w.TryAcquireTier("provider-rest", WorkTierBroadDiscovery)
		if !ok {
			t.Fatalf("could not establish optional provider pressure at slot %d: %+v", i, providerWorkDiag(w))
		}
		holds = append(holds, release)
	}
	defer func() {
		for _, release := range holds {
			release()
		}
	}()

	optionalCtx, optionalTier := freshnessRecoveryWorkContext(context.Background(), 3)
	if workTierFromContext(optionalCtx, WorkTierBackground) != WorkTierBroadDiscovery || optionalTier != WorkTierBroadDiscovery {
		t.Fatalf("optional recovery lost canonical broad-discovery tier: %v", optionalTier)
	}
	if !w.ShouldShed(optionalTier) {
		t.Fatal("optional stale-data recovery must shed while normal provider capacity is pressured")
	}

	protected := []struct {
		priority int
		want     WorkTier
	}{{1, WorkTierMarketCritical}, {2, WorkTierUserActionable}}
	protectedReleases := make([]func(), 0, len(protected))
	for _, tc := range protected {
		ctx, tier := freshnessRecoveryWorkContext(context.Background(), tc.priority)
		if tier != tc.want || w.ShouldShed(tier) {
			t.Fatalf("priority %d recovery must remain protected under local pressure: tier=%v", tc.priority, tier)
		}
		release, ok := w.TryAcquireTier("provider-rest", workTierFromContext(ctx, WorkTierBackground))
		if !ok {
			t.Fatalf("priority %d recovery could not use protected provider capacity: %+v", tc.priority, providerWorkDiag(w))
		}
		protectedReleases = append(protectedReleases, release)
	}
	for _, release := range protectedReleases {
		release()
	}
}

func TestV171LiveSubscriptionBudgetsExposeReservedHeadroom(t *testing.T) {
	alpaca := liveSubscriptionBudget("Alpaca IEX", 30, 25, 27, true)
	if alpaca.ReservedCapacity != 5 || alpaca.ReserveUsed != 2 || alpaca.Available != 3 || alpaca.Saturated {
		t.Fatalf("Alpaca live budget incorrect: %+v", alpaca)
	}
	finnhub := liveSubscriptionBudget("Finnhub", 50, 45, 50, true)
	if finnhub.ReservedCapacity != 5 || finnhub.ReserveUsed != 5 || finnhub.Available != 0 || !finnhub.Saturated {
		t.Fatalf("Finnhub live budget incorrect: %+v", finnhub)
	}
	degraded := deriveRuntimeDegradation(
		"degraded", "live", FeedDiagnostics{MarketSession: "regular"},
		[]FreshnessDiagnostic{{Dataset: "Quotes", State: "FRESH"}, {Dataset: "VIX", State: "FRESH"}},
		ProviderRouterSnapshot{}, RuntimeLoadDiagnostics{LiveSubscriptions: []LiveSubscriptionBudgetDiagnostics{finnhub}},
	)
	if degraded.Code != "LIVE CAPACITY SATURATED" || !degraded.CriticalUsable {
		t.Fatalf("live capacity saturation reason not attributable: %+v", degraded)
	}
}

func TestV171ProviderRateLimitCooldownShedsLowTierButProtectsActionable(t *testing.T) {
	telemetry := NewProviderTelemetry()
	done := telemetry.begin("Alpaca")
	done(&testHTTP429Error{})

	if ok, _ := telemetry.Allow("Alpaca", WorkTierBroadDiscovery); ok {
		t.Fatal("Tier 3 provider work should be shed during observed rate-limit cooldown")
	}
	if ok, reason := telemetry.Allow("Alpaca", WorkTierUserActionable); !ok {
		t.Fatalf("Tier 1 recovery must remain eligible during provider cooldown: %s", reason)
	}
	var row ProviderRequestDiagnostics
	for _, d := range telemetry.Diagnostics() {
		if d.Provider == "Alpaca" {
			row = d
		}
	}
	if row.BudgetState != "RATE LIMITED" || row.CooldownUntil <= time.Now().UnixMilli() || row.BudgetShed < 1 {
		t.Fatalf("rate-limit budget telemetry incomplete: %+v", row)
	}
}

func TestV171FinnhubBudgetComesFromLocalPacingNotGuessedEntitlement(t *testing.T) {
	finnhub := providerBudgetPerMinute("Finnhub")
	if finnhub <= 0 {
		t.Fatal("Finnhub local pacing should expose a justified request/minute budget")
	}
	if got := providerBudgetPerMinute("Alpaca"); got != 0 {
		t.Fatalf("Alpaca entitlement-dependent RPM must not be guessed, got %d", got)
	}
}
