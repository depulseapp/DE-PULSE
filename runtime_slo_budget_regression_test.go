package main

import (
	"testing"
	"time"
)

func t8RuntimeSLOCheck(t *testing.T, assessment RuntimeSLOAssessment, name string) RuntimeSLOCheck {
	t.Helper()
	for _, check := range assessment.Checks {
		if check.Name == name {
			return check
		}
	}
	t.Fatalf("runtime SLO check %q missing from canonical assessment: %+v", name, assessment.Checks)
	return RuntimeSLOCheck{}
}

func t8RuntimeSLOAssessment(load RuntimeLoadDiagnostics) RuntimeSLOAssessment {
	return buildRuntimeSLOAssessmentWithContext(
		"running",
		"live",
		FeedDiagnostics{MarketSession: "regular", FeedState: "streaming"},
		[]FreshnessDiagnostic{{Dataset: "Quotes", State: "LIVE"}, {Dataset: "VIX", State: "FRESH"}},
		load,
		ScannerState{},
		AppState{},
		nil,
		time.Date(2026, 8, 25, 14, 0, 0, 0, easternLocation()),
	)
}

func TestT8CanonicalRuntimeSLOBudgetsFailClosedAtOwnedThresholds(t *testing.T) {
	var load RuntimeLoadDiagnostics

	load.HTTP.InteractiveP95Ms = 250
	if got := t8RuntimeSLOCheck(t, t8RuntimeSLOAssessment(load), "Interactive API p95"); got.Status != "PASS" || got.Limit != "≤250 ms target · >1000 ms block" {
		t.Fatalf("interactive API target drifted: %+v", got)
	}
	load.HTTP.InteractiveP95Ms = 251
	if got := t8RuntimeSLOCheck(t, t8RuntimeSLOAssessment(load), "Interactive API p95"); got.Status != "WARN" {
		t.Fatalf("interactive API must warn above target: %+v", got)
	}
	load.HTTP.InteractiveP95Ms = 1001
	if got := t8RuntimeSLOCheck(t, t8RuntimeSLOAssessment(load), "Interactive API p95"); got.Status != "BLOCK" {
		t.Fatalf("interactive API must block above release limit: %+v", got)
	}
	load.HTTP.InteractiveP95Ms = 0

	load.Workload = []WorkClassDiagnostics{{Class: "provider-rest", Queued: 12, OldestQueueAgeMs: 2001}}
	if got := t8RuntimeSLOCheck(t, t8RuntimeSLOAssessment(load), "Provider queue"); got.Status != "WARN" {
		t.Fatalf("provider pressure must warn at owned target boundary: %+v", got)
	}
	load.Workload = []WorkClassDiagnostics{{Class: "provider-rest", Queued: 24, OldestQueueAgeMs: 0}}
	if got := t8RuntimeSLOCheck(t, t8RuntimeSLOAssessment(load), "Provider queue"); got.Status != "BLOCK" {
		t.Fatalf("provider queue depth must block at hard release boundary: %+v", got)
	}
	load.Workload = nil

	load.Persistence.OldestJobAgeMs = 2001
	if got := t8RuntimeSLOCheck(t, t8RuntimeSLOAssessment(load), "Persistence queue age"); got.Status != "WARN" {
		t.Fatalf("persistence queue age must warn above target: %+v", got)
	}
	load.Persistence.OldestJobAgeMs = 10001
	if got := t8RuntimeSLOCheck(t, t8RuntimeSLOAssessment(load), "Persistence queue age"); got.Status != "BLOCK" {
		t.Fatalf("persistence queue age must block above hard limit: %+v", got)
	}
	load.Persistence.OldestJobAgeMs = 0
	load.Persistence.RowsWrittenLastMin = 601
	if got := t8RuntimeSLOCheck(t, t8RuntimeSLOAssessment(load), "DB write rate"); got.Status != "WARN" {
		t.Fatalf("DB write rate must warn above target: %+v", got)
	}
	load.Persistence.RowsWrittenLastMin = 3001
	if got := t8RuntimeSLOCheck(t, t8RuntimeSLOAssessment(load), "DB write rate"); got.Status != "BLOCK" {
		t.Fatalf("DB write rate must block above hard limit: %+v", got)
	}
	load.Persistence.RowsWrittenLastMin = 0

	load.CPUUtilizationPct = 71
	if got := t8RuntimeSLOCheck(t, t8RuntimeSLOAssessment(load), "CPU utilization"); got.Status != "WARN" {
		t.Fatalf("CPU must warn above owned target: %+v", got)
	}
	load.CPUUtilizationPct = 91
	if got := t8RuntimeSLOCheck(t, t8RuntimeSLOAssessment(load), "CPU utilization"); got.Status != "BLOCK" {
		t.Fatalf("CPU must block above owned release limit: %+v", got)
	}
	load.CPUUtilizationPct = 0

	load.Goroutines = 601
	if got := t8RuntimeSLOCheck(t, t8RuntimeSLOAssessment(load), "Local process budget"); got.Status != "BLOCK" {
		t.Fatalf("goroutine budget must block above hard limit: %+v", got)
	}
	load.Goroutines = 0
	load.HeapAllocBytes = 769 * 1024 * 1024
	if got := t8RuntimeSLOCheck(t, t8RuntimeSLOAssessment(load), "Local process budget"); got.Status != "BLOCK" {
		t.Fatalf("heap budget must block above hard limit: %+v", got)
	}
	load.HeapAllocBytes = 0

	load.Startup.BootstrapDurationMs = 1001
	if got := t8RuntimeSLOCheck(t, t8RuntimeSLOAssessment(load), "Startup/warm-start time"); got.Status != "WARN" {
		t.Fatalf("startup must warn above target: %+v", got)
	}
	load.Startup.BootstrapDurationMs = 5001
	if got := t8RuntimeSLOCheck(t, t8RuntimeSLOAssessment(load), "Startup/warm-start time"); got.Status != "BLOCK" {
		t.Fatalf("startup must block above hard limit: %+v", got)
	}
	load.Startup.BootstrapDurationMs = 0

	load.StorageGrowthBytes = 251 * 1024 * 1024
	if got := t8RuntimeSLOCheck(t, t8RuntimeSLOAssessment(load), "Storage growth"); got.Status != "WARN" {
		t.Fatalf("storage growth must warn above target: %+v", got)
	}
	load.StorageGrowthBytes = 1024*1024*1024 + 1
	if got := t8RuntimeSLOCheck(t, t8RuntimeSLOAssessment(load), "Storage growth"); got.Status != "BLOCK" {
		t.Fatalf("storage growth must block above hard limit: %+v", got)
	}
}
