package main

import (
	"context"
	"runtime"
	"runtime/metrics"
	"time"
)

type LiveSubscriptionBudgetDiagnostics struct {
	Provider          string `json:"provider"`
	Capacity          int    `json:"capacity"`
	NormalTarget      int    `json:"normalTarget"`
	ReservedCapacity int    `json:"reservedCapacity"`
	Active            int    `json:"active"`
	Available         int    `json:"available"`
	ReserveUsed       int    `json:"reserveUsed"`
	UtilizationPct    int    `json:"utilizationPct"`
	Saturated         bool   `json:"saturated"`
	Connected         bool   `json:"connected"`
}

type RuntimeLoadDiagnostics struct {
	SampledAt                int64                               `json:"sampledAt"`
	Goroutines               int                                 `json:"goroutines"`
	GOMAXPROCS               int                                 `json:"gomaxprocs"`
	CPUUtilizationPct        float64                             `json:"cpuUtilizationPct"`
	HeapAllocBytes           uint64                              `json:"heapAllocBytes"`
	HeapInuseBytes           uint64                              `json:"heapInuseBytes"`
	HeapObjects              uint64                              `json:"heapObjects"`
	TotalAllocBytes          uint64                              `json:"totalAllocBytes"`
	NumGC                    uint32                              `json:"numGC"`
	GCPauseTotalNs           uint64                              `json:"gcPauseTotalNs"`
	Persistence              PersistenceDiagnostics              `json:"persistence"`
	Workload                 []WorkClassDiagnostics              `json:"workload"`
	ProviderRequests         []ProviderRequestDiagnostics        `json:"providerRequests,omitempty"`
	LiveSubscriptions        []LiveSubscriptionBudgetDiagnostics `json:"liveSubscriptions,omitempty"`
	BroadSnapshots           BroadSnapshotBrokerDiagnostics      `json:"broadSnapshots"`
	HTTP                     HTTPRuntimeDiagnostics              `json:"http"`
	ProviderCallsAvoided     int64                               `json:"providerCallsAvoided"`
	CanonicalReuseHitRatePct int                                 `json:"canonicalReuseHitRatePct"`
	StorageGrowthBytes       int64                               `json:"storageGrowthBytes"`
	Startup                  RuntimeStartupDiagnostics           `json:"startup"`
	Recovery                 RuntimeRecoveryDiagnostics          `json:"recovery"`
	CanonicalPipeline        CanonicalPipelineDiagnostics        `json:"canonicalPipeline"`
}

func runtimeCPUTotalSeconds() float64 {
	samples := []metrics.Sample{{Name: "/cpu/classes/total:cpu-seconds"}}
	metrics.Read(samples)
	if samples[0].Value.Kind() != metrics.KindFloat64 {
		return 0
	}
	return samples[0].Value.Float64()
}

func liveSubscriptionBudget(provider string, capacity, normalTarget, active int, connected bool) LiveSubscriptionBudgetDiagnostics {
	available := capacity - active
	if available < 0 {
		available = 0
	}
	reserveUsed := active - normalTarget
	if reserveUsed < 0 {
		reserveUsed = 0
	}
	utilization := 0
	if capacity > 0 {
		utilization = (active * 100) / capacity
	}
	return LiveSubscriptionBudgetDiagnostics{
		Provider: provider, Capacity: capacity, NormalTarget: normalTarget,
		ReservedCapacity: capacity - normalTarget, Active: active, Available: available,
		ReserveUsed: reserveUsed, UtilizationPct: utilization,
		Saturated: capacity > 0 && active >= capacity, Connected: connected,
	}
}

func (e *Engine) sampleRuntimeLoad() {
	var mem runtime.MemStats
	runtime.ReadMemStats(&mem)
	now := time.Now()
	cpuTotal := runtimeCPUTotalSeconds()
	cpuPct := 0.0
	e.mu.Lock()
	if !e.lastRuntimeCPUSampledAt.IsZero() && cpuTotal >= e.lastRuntimeCPUTotal {
		elapsed := now.Sub(e.lastRuntimeCPUSampledAt).Seconds()
		capacity := elapsed * float64(runtime.GOMAXPROCS(0))
		if capacity > 0 {
			cpuPct = (cpuTotal - e.lastRuntimeCPUTotal) * 100 / capacity
			if cpuPct < 0 {
				cpuPct = 0
			}
		}
	}
	e.lastRuntimeCPUTotal = cpuTotal
	e.lastRuntimeCPUSampledAt = now
	startup := e.startupProfile
	initialStorage := e.initialPersistenceBytes
	e.mu.Unlock()
	persistence := PersistenceDiagnostics{Backend: "disabled"}
	if e.app != nil && e.app.persistence != nil {
		persistence = e.app.persistence.Diagnostics()
	}
	workload := []WorkClassDiagnostics(nil)
	if e.workload != nil {
		workload = e.workload.Diagnostics()
	}
	providers := []ProviderRequestDiagnostics(nil)
	if e.providerTelemetry != nil {
		providers = e.providerTelemetry.Diagnostics()
	}
	e.mu.RLock()
	finnhubActive := len(e.subscribedSymbols)
	alpacaActive := len(e.alpacaSubscribedSymbols)
	finnhubConnected := e.webSocketConnected
	alpacaConnected := e.alpacaWebSocketConnected
	e.mu.RUnlock()
	liveSubscriptions := []LiveSubscriptionBudgetDiagnostics{
		liveSubscriptionBudget("Alpaca IEX", alpacaPlanMaxSymbols, alpacaActiveTarget, alpacaActive, alpacaConnected),
		liveSubscriptionBudget("Finnhub", finnhubPlanMaxSymbols, finnhubActiveTarget, finnhubActive, finnhubConnected),
	}
	broadSnapshots := e.broadSnapshotDiagnostics()
	httpDiag := HTTPRuntimeDiagnostics{}
	if e.app != nil && e.app.httpTelemetry != nil {
		httpDiag = e.app.httpTelemetry.Diagnostics()
	}
	e.mu.RLock()
	avoided := e.providerCallsAvoided
	canonicalPipeline := e.canonicalPipeline
	e.mu.RUnlock()
	providerRequests := int64(0)
	for _, row := range providers {
		providerRequests += row.Requests
	}
	reuseRate := 0
	if total := avoided + providerRequests; total > 0 {
		reuseRate = int(avoided * 100 / total)
	}
	storageGrowth := persistence.Store.StorageBytes - initialStorage
	sample := RuntimeLoadDiagnostics{
		SampledAt:                now.UnixMilli(),
		Goroutines:               runtime.NumGoroutine(),
		GOMAXPROCS:               runtime.GOMAXPROCS(0),
		CPUUtilizationPct:        cpuPct,
		HeapAllocBytes:           mem.HeapAlloc,
		HeapInuseBytes:           mem.HeapInuse,
		HeapObjects:              mem.HeapObjects,
		TotalAllocBytes:          mem.TotalAlloc,
		NumGC:                    mem.NumGC,
		GCPauseTotalNs:           mem.PauseTotalNs,
		Persistence:              persistence,
		Workload:                 workload,
		ProviderRequests:         providers,
		LiveSubscriptions:        liveSubscriptions,
		BroadSnapshots:           broadSnapshots,
		HTTP:                     httpDiag,
		ProviderCallsAvoided:     avoided,
		CanonicalReuseHitRatePct: reuseRate,
		StorageGrowthBytes:       storageGrowth,
		Startup:                  startup,
		CanonicalPipeline:        canonicalPipeline,
	}
	e.mu.Lock()
	e.runtimeLoad = sample
	e.lastUpdated["runtime-load-profile"] = sample.SampledAt
	e.mu.Unlock()
}

func (e *Engine) runtimeLoadProfileLoop(ctx context.Context) {
	e.sampleRuntimeLoad()
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			e.sampleRuntimeLoad()
			return
		case <-ticker.C:
			e.sampleRuntimeLoad()
		}
	}
}
