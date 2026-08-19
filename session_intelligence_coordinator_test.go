package main

import (
	"os"
	"strings"
	"testing"
	"time"
)

func v186ET(year int, month time.Month, day, hour, minute int) time.Time {
	return time.Date(year, month, day, hour, minute, 0, 0, easternLocation())
}

func TestV186SessionCoordinatorPreservesPreMarketWindowAndCatchup(t *testing.T) {
	scheduledAt := v186ET(2026, time.August, 18, 3, 20)
	plan := preMarketCheckpointPlan(scheduledAt, PreparationJobStatus{})
	if !plan.Run || plan.Reason != "scheduled" || plan.Late {
		t.Fatalf("expected scheduled pre-market checkpoint, got %+v", plan)
	}

	catchupAt := v186ET(2026, time.August, 18, 4, 5)
	plan = preMarketCheckpointPlan(catchupAt, PreparationJobStatus{})
	if !plan.Run || plan.Reason != "missed-window catch-up" || !plan.Late {
		t.Fatalf("expected late pre-market catch-up, got %+v", plan)
	}

	cutoff := v186ET(2026, time.August, 18, 9, 20)
	if plan := preMarketCheckpointPlan(cutoff, PreparationJobStatus{}); plan.Run {
		t.Fatalf("pre-market catch-up must stop at Market Open Prep handoff, got %+v", plan)
	}
}

func TestV186SessionCoordinatorPreservesMarketOpenWindowAndCatchup(t *testing.T) {
	scheduledAt := v186ET(2026, time.August, 18, 9, 22)
	plan := marketOpenCheckpointPlan(scheduledAt, PreparationJobStatus{})
	if !plan.Run || plan.Reason != "scheduled" || plan.Late {
		t.Fatalf("expected scheduled market-open checkpoint, got %+v", plan)
	}

	catchupAt := v186ET(2026, time.August, 18, 9, 31)
	plan = marketOpenCheckpointPlan(catchupAt, PreparationJobStatus{})
	if !plan.Run || plan.Reason != "missed-window catch-up" || !plan.Late {
		t.Fatalf("expected market-open catch-up, got %+v", plan)
	}

	afterCutoff := v186ET(2026, time.August, 18, 10, 16)
	if plan := marketOpenCheckpointPlan(afterCutoff, PreparationJobStatus{}); plan.Run {
		t.Fatalf("market-open catch-up must stop after 10:15 ET, got %+v", plan)
	}
}

func TestV186SessionCoordinatorHonorsRetryAndSameDayCompletion(t *testing.T) {
	now := v186ET(2026, time.August, 18, 9, 22)
	tooSoon := PreparationJobStatus{LastAttempt: now.Add(-time.Minute).UnixMilli()}
	if plan := marketOpenCheckpointPlan(now, tooSoon); plan.Run {
		t.Fatalf("market-open retry guard was bypassed: %+v", plan)
	}

	preNow := v186ET(2026, time.August, 18, 3, 20)
	preTooSoon := PreparationJobStatus{LastAttempt: preNow.Add(-5 * time.Minute).UnixMilli()}
	if plan := preMarketCheckpointPlan(preNow, preTooSoon); plan.Run {
		t.Fatalf("pre-market retry guard was bypassed: %+v", plan)
	}

	completed := PreparationJobStatus{
		TradingDay:  "2026-08-18",
		LastSuccess: now.Add(-2 * time.Minute).UnixMilli(),
	}
	if plan := marketOpenCheckpointPlan(now, completed); plan.Run {
		t.Fatalf("completed trading-day checkpoint must not rerun: %+v", plan)
	}
}

func TestV186StartLiveUsesOneAutomaticSessionCoordinator(t *testing.T) {
	raw, err := os.ReadFile("runtime_jobs.go")
	if err != nil {
		t.Fatal(err)
	}
	src := string(raw)
	if strings.Count(src, "e.sessionIntelligenceCoordinatorLoop(ctx, key, alpacaKey, alpacaSecret)") != 1 {
		t.Fatalf("startLive must start exactly one Session Intelligence Coordinator")
	}
	if strings.Contains(src, "e.preMarketPrepLoop(ctx, key, alpacaKey, alpacaSecret)") {
		t.Fatalf("legacy Pre-Market Prep scheduler is still started in parallel")
	}
	if strings.Contains(src, "e.marketOpenPrepLoop(ctx)") {
		t.Fatalf("legacy Market Open Prep scheduler is still started in parallel")
	}
	if !strings.Contains(src, "session-intelligence-coordinator") {
		t.Fatalf("coordinator health state must be initialized")
	}
}
