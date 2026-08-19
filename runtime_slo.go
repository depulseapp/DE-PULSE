package main

import (
	"fmt"
	"sort"
	"strings"
	"sync"
	"time"
)

const (
	runtimeRecoveryHealthyObservationsRequired = 3
	runtimeRecoveryMinHealthyDuration          = 5 * time.Second
)

type RuntimeSLOCheck struct {
	Name     string `json:"name"`
	Status   string `json:"status"`
	Value    string `json:"value"`
	Limit    string `json:"limit"`
	Critical bool   `json:"critical"`
}

type RuntimeSLOAssessment struct {
	Status string            `json:"status"`
	Checks []RuntimeSLOCheck `json:"checks"`
}

type RuntimeStartupDiagnostics struct {
	BootstrapDurationMs    int64 `json:"bootstrapDurationMs"`
	CacheQuotesLoaded      int   `json:"cacheQuotesLoaded"`
	PersistedQuotesApplied int   `json:"persistedQuotesApplied"`
	WarmStartQuotes        int   `json:"warmStartQuotes"`
	WarmStartTargetQuotes  int   `json:"warmStartTargetQuotes"`
	WarmStartCoveragePct   int   `json:"warmStartCoveragePct"`
}

type RuntimeRecoveryDiagnostics struct {
	CurrentlyStaleDatasets          int   `json:"currentlyStaleDatasets"`
	StaleToCurrentEvents            int64 `json:"staleToCurrentEvents"`
	LastStaleToCurrentMs            int64 `json:"lastStaleToCurrentMs,omitempty"`
	DegradationRecoveryEvents       int64 `json:"degradationRecoveryEvents"`
	LastDegradationRecoveryMs       int64 `json:"lastDegradationRecoveryMs,omitempty"`
	DegradedSince                   int64 `json:"degradedSince,omitempty"`
	DegradationRecoveryPending      bool  `json:"degradationRecoveryPending,omitempty"`
	DegradationHealthyObservations  int   `json:"degradationHealthyObservations,omitempty"`
	DegradationHealthyRequired      int   `json:"degradationHealthyRequired,omitempty"`
	DegradationHealthySince         int64 `json:"degradationHealthySince,omitempty"`
}

type RuntimeSLOTracker struct {
	mu                          sync.Mutex
	staleSince                  map[string]time.Time
	staleToCurrentEvents        int64
	lastStaleToCurrentMs        int64
	degradedSince               time.Time
	degradationRecoveryEvents   int64
	lastDegradationRecoveryMs   int64
	lastDegradation             RuntimeDegradationState
	recoveryHealthySince        time.Time
	recoveryHealthyObservations int
}

func NewRuntimeSLOTracker() *RuntimeSLOTracker {
	return &RuntimeSLOTracker{staleSince: map[string]time.Time{}}
}

func freshnessNeedsRecovery(state string) bool {
	switch strings.ToUpper(strings.TrimSpace(state)) {
	case "STALE", "ERROR", "UNAVAILABLE":
		return true
	default:
		return false
	}
}

// StabilizeDegradation owns recovery hysteresis for the canonical runtime
// degradation payload. A real degraded condition becomes visible immediately,
// but recovery requires multiple consecutive healthy observations and a small
// stability window so one transient success cannot falsely declare recovery.
func (t *RuntimeSLOTracker) StabilizeDegradation(current RuntimeDegradationState, now time.Time) RuntimeDegradationState {
	if t == nil {
		return current
	}
	t.mu.Lock()
	defer t.mu.Unlock()

	if strings.TrimSpace(current.Code) != "" {
		t.lastDegradation = current
		t.recoveryHealthySince = time.Time{}
		t.recoveryHealthyObservations = 0
		return current
	}
	if strings.TrimSpace(t.lastDegradation.Code) == "" {
		return current
	}
	if t.recoveryHealthySince.IsZero() {
		t.recoveryHealthySince = now
	}
	t.recoveryHealthyObservations++
	stableFor := now.Sub(t.recoveryHealthySince)
	if stableFor < 0 {
		stableFor = 0
	}
	if t.recoveryHealthyObservations >= runtimeRecoveryHealthyObservationsRequired && stableFor >= runtimeRecoveryMinHealthyDuration {
		t.lastDegradation = RuntimeDegradationState{}
		t.recoveryHealthySince = time.Time{}
		t.recoveryHealthyObservations = 0
		return current
	}

	held := t.lastDegradation
	held.PressureState = "RECOVERING"
	held.Detail = strings.TrimSpace(held.Detail)
	if held.Detail != "" {
		held.Detail += " · "
	}
	held.Detail += fmt.Sprintf("Recovery validation in progress (%d/%d consecutive healthy observations; minimum %s stability).", t.recoveryHealthyObservations, runtimeRecoveryHealthyObservationsRequired, runtimeRecoveryMinHealthyDuration)
	return held
}

func (t *RuntimeSLOTracker) Observe(freshness []FreshnessDiagnostic, degradation RuntimeDegradationState, now time.Time) RuntimeRecoveryDiagnostics {
	if t == nil {
		return RuntimeRecoveryDiagnostics{}
	}
	t.mu.Lock()
	defer t.mu.Unlock()
	seen := map[string]bool{}
	for _, row := range freshness {
		key := strings.TrimSpace(row.Dataset)
		if key == "" {
			continue
		}
		seen[key] = true
		if freshnessNeedsRecovery(row.State) {
			if _, exists := t.staleSince[key]; !exists {
				t.staleSince[key] = now
			}
			continue
		}
		if started, exists := t.staleSince[key]; exists {
			ms := now.Sub(started).Milliseconds()
			if ms < 0 {
				ms = 0
			}
			t.lastStaleToCurrentMs = ms
			t.staleToCurrentEvents++
			delete(t.staleSince, key)
		}
	}
	for key := range t.staleSince {
		if !seen[key] {
			delete(t.staleSince, key)
		}
	}
	if strings.TrimSpace(degradation.Code) != "" {
		if t.degradedSince.IsZero() {
			t.degradedSince = now
		}
	} else if !t.degradedSince.IsZero() {
		ms := now.Sub(t.degradedSince).Milliseconds()
		if ms < 0 {
			ms = 0
		}
		t.lastDegradationRecoveryMs = ms
		t.degradationRecoveryEvents++
		t.degradedSince = time.Time{}
	}
	return t.diagnosticsLocked()
}

func (t *RuntimeSLOTracker) diagnosticsLocked() RuntimeRecoveryDiagnostics {
	d := RuntimeRecoveryDiagnostics{
		CurrentlyStaleDatasets:         len(t.staleSince),
		StaleToCurrentEvents:           t.staleToCurrentEvents,
		LastStaleToCurrentMs:           t.lastStaleToCurrentMs,
		DegradationRecoveryEvents:      t.degradationRecoveryEvents,
		LastDegradationRecoveryMs:      t.lastDegradationRecoveryMs,
		DegradationRecoveryPending:     strings.TrimSpace(t.lastDegradation.Code) != "" && t.recoveryHealthyObservations > 0,
		DegradationHealthyObservations: t.recoveryHealthyObservations,
		DegradationHealthyRequired:     runtimeRecoveryHealthyObservationsRequired,
	}
	if !t.degradedSince.IsZero() {
		d.DegradedSince = t.degradedSince.UnixMilli()
	}
	if !t.recoveryHealthySince.IsZero() {
		d.DegradationHealthySince = t.recoveryHealthySince.UnixMilli()
	}
	return d
}

func actionableSLOSymbols(st AppState) []string {
	syms := append([]string(nil), activeDeskSymbolsFromState(st)...)
	if selected, ok := parseUserTicker(st.UI.SelectedTicker); ok {
		syms = append(syms, selected)
	}
	// User-approved live tradable exceptions stay in the actionable freshness pool.
	syms = append(syms, "GLD", "SLV", "USO")
	return uniqueSymbols(syms)
}

func currentQuoteCoverage(quotes map[string]Quote, symbols []string, now int64) (current, total int, oldestMs int64) {
	for _, symbol := range uniqueSymbols(symbols) {
		symbol = normalizeSymbol(symbol)
		if symbol == "" {
			continue
		}
		total++
		q, ok := quotes[symbol]
		if !ok || q.Price <= 0 {
			continue
		}
		providerAge, receiptAge, isCurrent := quoteEvidenceAges(q, now)
		age := providerAge
		if receiptAge > age {
			age = receiptAge
		}
		if age > oldestMs {
			oldestMs = age
		}
		if isCurrent {
			current++
		}
	}
	return
}

func buildRuntimeSLOAssessmentWithContext(status, mode string, feed FeedDiagnostics, freshness []FreshnessDiagnostic, load RuntimeLoadDiagnostics, scanner ScannerState, appState AppState, quotes map[string]Quote, now time.Time) RuntimeSLOAssessment {
	checks := []RuntimeSLOCheck{}
	add := func(name, state, value, limit string, critical bool) {
		checks = append(checks, RuntimeSLOCheck{Name: name, Status: state, Value: value, Limit: limit, Critical: critical})
	}
	regular := feed.MarketSession == "regular" && mode != "demo" && status != "stopped"
	criticalUsable := criticalDecisionDataUsable(freshness, feed.MarketSession)
	criticalState := "PASS"
	if regular && !criticalUsable {
		criticalState = "BLOCK"
	}
	add("Critical decision datasets", criticalState, fmt.Sprintf("usable=%t", criticalUsable), "usable during regular session", true)

	if quotes != nil {
		selected := normalizeSymbol(appState.UI.SelectedTicker)
		if selected != "" {
			cur, total, age := currentQuoteCoverage(quotes, []string{selected}, now.UnixMilli())
			state := "PASS"
			if regular && (total == 0 || cur != total) {
				state = "BLOCK"
			}
			add("Selected-symbol freshness", state, fmt.Sprintf("%d/%d current · oldest %d ms", cur, total, age), "100% current during regular session", true)
		}
		syms := actionableSLOSymbols(appState)
		cur, total, age := currentQuoteCoverage(quotes, syms, now.UnixMilli())
		coverage := 100
		if total > 0 {
			coverage = cur * 100 / total
		}
		state := "PASS"
		if regular && coverage < 80 {
			state = "BLOCK"
		} else if regular && coverage < 95 {
			state = "WARN"
		}
		add("Actionable-watchlist freshness", state, fmt.Sprintf("%d/%d current · %d%% · oldest %d ms", cur, total, coverage, age), "≥95% target · <80% block during regular session", true)
	}

	p95 := load.HTTP.InteractiveP95Ms
	apiState := "PASS"
	if p95 > 1000 {
		apiState = "BLOCK"
	} else if p95 > 250 {
		apiState = "WARN"
	}
	add("Interactive API p95", apiState, fmt.Sprintf("%d ms", p95), "≤250 ms target · >1000 ms block", true)

	providerOldest, providerQueued := int64(0), 0
	maxQueueDepth, maxQueueAge := 0, int64(0)
	for _, row := range load.Workload {
		if row.Queued > maxQueueDepth {
			maxQueueDepth = row.Queued
		}
		if row.OldestQueueAgeMs > maxQueueAge {
			maxQueueAge = row.OldestQueueAgeMs
		}
		if row.Class == "provider-rest" {
			providerQueued, providerOldest = row.Queued, row.OldestQueueAgeMs
		}
	}
	providerState := "PASS"
	if providerOldest > 5000 || providerQueued >= 24 {
		providerState = "BLOCK"
	} else if providerOldest > 2000 || providerQueued >= 12 {
		providerState = "WARN"
	}
	add("Provider queue", providerState, fmt.Sprintf("depth=%d · oldest=%d ms", providerQueued, providerOldest), "depth <12 and age ≤2000 ms target · depth ≥24 or age >5000 ms block", true)
	add("All bounded queues", map[bool]string{true: "WARN", false: "PASS"}[maxQueueAge > 5000 || maxQueueDepth >= 24], fmt.Sprintf("max depth=%d · max age=%d ms", maxQueueDepth, maxQueueAge), "bounded; low-priority work sheds before critical tiers", false)

	persistAge := load.Persistence.OldestJobAgeMs
	persistState := "PASS"
	if persistAge > 10000 {
		persistState = "BLOCK"
	} else if persistAge > 2000 {
		persistState = "WARN"
	}
	add("Persistence queue age", persistState, fmt.Sprintf("%d ms", persistAge), "≤2000 ms target · >10000 ms block", false)
	writeState := "PASS"
	if load.Persistence.RowsWrittenLastMin > 3000 {
		writeState = "BLOCK"
	} else if load.Persistence.RowsWrittenLastMin > 600 {
		writeState = "WARN"
	}
	add("DB write rate", writeState, fmt.Sprintf("%d batches/min · %d rows/min", load.Persistence.WriteBatchesLastMin, load.Persistence.RowsWrittenLastMin), "≤600 rows/min target · >3000 rows/min block", false)

	cpuState := "PASS"
	if load.CPUUtilizationPct > 90 {
		cpuState = "BLOCK"
	} else if load.CPUUtilizationPct > 70 {
		cpuState = "WARN"
	}
	add("CPU utilization", cpuState, fmt.Sprintf("%.1f%% of GOMAXPROCS=%d", load.CPUUtilizationPct, load.GOMAXPROCS), "≤70% target · >90% block", true)
	processState := "PASS"
	if load.Goroutines > 600 || load.HeapAllocBytes > 768*1024*1024 {
		processState = "BLOCK"
	}
	add("Local process budget", processState, fmt.Sprintf("%d goroutines · %.1f MiB heap", load.Goroutines, float64(load.HeapAllocBytes)/(1024*1024)), "≤600 goroutines · ≤768 MiB heap", true)

	radarState := "PASS"
	if scanner.Radar.CadenceMs > 0 && scanner.Radar.DurationMs > scanner.Radar.CadenceMs {
		radarState = "WARN"
	}
	add("Opportunity Radar cycle", radarState, fmt.Sprintf("%d ms", scanner.Radar.DurationMs), "complete within configured cadence", false)

	maxBudget, knownBudgets := 0, 0
	for _, p := range load.ProviderRequests {
		if p.BudgetPerMinute > 0 {
			knownBudgets++
			if p.BudgetUtilizationPct > maxBudget {
				maxBudget = p.BudgetUtilizationPct
			}
		}
	}
	budgetState := "PASS"
	if maxBudget >= 100 {
		budgetState = "BLOCK"
	} else if maxBudget >= 80 {
		budgetState = "WARN"
	}
	add("Provider request budgets", budgetState, fmt.Sprintf("known=%d · max utilization=%d%%", knownBudgets, maxBudget), "<80% target · 100% block for known local budgets", true)

	maxSub := 0
	providers := make([]string, 0, len(load.LiveSubscriptions))
	for _, sub := range load.LiveSubscriptions {
		if sub.UtilizationPct > maxSub {
			maxSub = sub.UtilizationPct
		}
		providers = append(providers, fmt.Sprintf("%s %d/%d", sub.Provider, sub.Active, sub.Capacity))
	}
	sort.Strings(providers)
	subState := "PASS"
	if maxSub >= 100 {
		subState = "WARN"
	}
	add("Live subscription utilization", subState, strings.Join(providers, " · "), "reserved capacity protects Tier 0/1; saturation is explicit", true)

	startupState := "PASS"
	if load.Startup.BootstrapDurationMs > 5000 {
		startupState = "BLOCK"
	} else if load.Startup.BootstrapDurationMs > 1000 {
		startupState = "WARN"
	}
	add("Startup/warm-start time", startupState, fmt.Sprintf("%d ms · %d/%d warm quotes (%d%%)", load.Startup.BootstrapDurationMs, load.Startup.WarmStartQuotes, load.Startup.WarmStartTargetQuotes, load.Startup.WarmStartCoveragePct), "≤1000 ms target · >5000 ms block", false)
	add("Canonical reuse / provider calls avoided", "PASS", fmt.Sprintf("%d avoided · %d%% reuse hit rate", load.ProviderCallsAvoided, load.CanonicalReuseHitRatePct), "measure and trend upward without sacrificing freshness", false)

	storageState := "PASS"
	if load.StorageGrowthBytes > 1024*1024*1024 {
		storageState = "BLOCK"
	} else if load.StorageGrowthBytes > 250*1024*1024 {
		storageState = "WARN"
	}
	add("Storage growth", storageState, fmt.Sprintf("%d bytes since startup · %d bytes current", load.StorageGrowthBytes, load.Persistence.Store.StorageBytes), "≤250 MiB/startup target · >1 GiB block pending retention review", false)

	staleRecoveryState := "PASS"
	if load.Recovery.LastStaleToCurrentMs > 300000 {
		staleRecoveryState = "BLOCK"
	} else if load.Recovery.LastStaleToCurrentMs > 120000 {
		staleRecoveryState = "WARN"
	}
	add("Stale→current recovery", staleRecoveryState, fmt.Sprintf("events=%d · last=%d ms · currently stale=%d", load.Recovery.StaleToCurrentEvents, load.Recovery.LastStaleToCurrentMs, load.Recovery.CurrentlyStaleDatasets), "≤120s target · >300s block when recovery occurs", true)
	degradeState := "PASS"
	if load.Recovery.LastDegradationRecoveryMs > 300000 {
		degradeState = "BLOCK"
	} else if load.Recovery.LastDegradationRecoveryMs > 120000 {
		degradeState = "WARN"
	}
	if load.Recovery.DegradationRecoveryPending && degradeState == "PASS" {
		degradeState = "WARN"
	}
	add("Degradation recovery", degradeState, fmt.Sprintf("events=%d · last=%d ms · healthy confirmations=%d/%d", load.Recovery.DegradationRecoveryEvents, load.Recovery.LastDegradationRecoveryMs, load.Recovery.DegradationHealthyObservations, load.Recovery.DegradationHealthyRequired), "3 consecutive healthy confirmations and ≥5s stability before clear; ≤120s target · >300s block when recovery occurs", true)

	overall := "PASS"
	for _, c := range checks {
		if c.Status == "BLOCK" {
			overall = "BLOCK"
			break
		}
		if c.Status == "WARN" && overall == "PASS" {
			overall = "WARN"
		}
	}
	return RuntimeSLOAssessment{Status: overall, Checks: checks}
}
