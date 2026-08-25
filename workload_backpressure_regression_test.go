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

func t8WorkClassDiag(w *WorkloadController, class string) WorkClassDiagnostics {
	for _, row := range w.Diagnostics() {
		if row.Class == class {
			return row
		}
	}
	return WorkClassDiagnostics{}
}

func waitT8WorkClassState(t *testing.T, w *WorkloadController, class string, wantQueued, wantInFlight int) WorkClassDiagnostics {
	t.Helper()
	deadline := time.Now().Add(time.Second)
	for time.Now().Before(deadline) {
		d := t8WorkClassDiag(w, class)
		if d.Class == class && d.Queued == wantQueued && d.InFlight == wantInFlight {
			return d
		}
		time.Sleep(2 * time.Millisecond)
	}
	d := t8WorkClassDiag(w, class)
	t.Fatalf("%s did not converge to queued=%d inFlight=%d: %+v", class, wantQueued, wantInFlight, d)
	return d
}

func TestT8RepeatedWorkloadPressureConvergesWithoutLeakedPermitsOrQueue(t *testing.T) {
	cases := []struct {
		class  string
		rounds int
	}{
		{class: "provider-rest", rounds: 12},
		{class: "scanner", rounds: 32},
		{class: "background", rounds: 32},
		{class: "ai", rounds: 32},
	}

	for _, tc := range cases {
		t.Run(tc.class, func(t *testing.T) {
			w := NewWorkloadController()
			baseline := t8WorkClassDiag(w, tc.class)
			if baseline.Class != tc.class || baseline.Capacity < 1 || baseline.MaxQueue < 1 {
				t.Fatalf("%s workload contract unavailable: %+v", tc.class, baseline)
			}

			var expectedCompleted int64
			var expectedCanceled int64
			for round := 0; round < tc.rounds; round++ {
				holds := make([]func(), 0, baseline.Capacity)
				for slot := 0; slot < baseline.Capacity; slot++ {
					release, ok := w.TryAcquireTier(tc.class, WorkTierMarketCritical)
					if !ok {
						t.Fatalf("round %d could not establish %s pressure at slot %d/%d: %+v", round, tc.class, slot+1, baseline.Capacity, t8WorkClassDiag(w, tc.class))
					}
					holds = append(holds, release)
				}

				ctx, cancel := context.WithCancel(context.Background())
				done := make(chan struct{}, baseline.MaxQueue)
				for queued := 0; queued < baseline.MaxQueue; queued++ {
					go func() {
						release, ok := w.AcquireTier(ctx, tc.class, WorkTierUserActionable)
						if ok {
							release()
						}
						done <- struct{}{}
					}()
				}

				pressured := waitT8WorkClassState(t, w, tc.class, baseline.MaxQueue, baseline.Capacity)
				if pressured.Queued > pressured.MaxQueue || pressured.InFlight > pressured.Capacity {
					t.Fatalf("round %d exceeded bounded %s resources: %+v", round, tc.class, pressured)
				}

				cancel()
				for worker := 0; worker < baseline.MaxQueue; worker++ {
					select {
					case <-done:
					case <-time.After(time.Second):
						t.Fatalf("round %d %s queued worker did not cancel", round, tc.class)
					}
				}
				expectedCanceled += int64(baseline.MaxQueue)

				for _, release := range holds {
					release()
				}
				expectedCompleted += int64(baseline.Capacity)

				idle := waitT8WorkClassState(t, w, tc.class, 0, 0)
				if idle.Canceled != expectedCanceled || idle.Completed != expectedCompleted {
					t.Fatalf("round %d %s resource accounting did not converge: completed=%d/%d canceled=%d/%d diag=%+v", round, tc.class, idle.Completed, expectedCompleted, idle.Canceled, expectedCanceled, idle)
				}

				recovery, ok := w.TryAcquireTier(tc.class, WorkTierBackground)
				if !ok {
					t.Fatalf("round %d %s did not recover admission after pressure: %+v", round, tc.class, t8WorkClassDiag(w, tc.class))
				}
				recovery()
				expectedCompleted++
				idle = waitT8WorkClassState(t, w, tc.class, 0, 0)
				if idle.Completed != expectedCompleted {
					t.Fatalf("round %d %s recovery permit was not released exactly once: %+v", round, tc.class, idle)
				}
			}

			final := waitT8WorkClassState(t, w, tc.class, 0, 0)
			if final.Queued != 0 || final.InFlight != 0 {
				t.Fatalf("%s soak left runtime work behind: %+v", tc.class, final)
			}
		})
	}
}
