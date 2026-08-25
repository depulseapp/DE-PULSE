package main

import (
	"context"
	"testing"
	"time"
)

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
