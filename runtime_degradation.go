package main

import (
	"sort"
	"strings"
)

type RuntimeDegradationState struct {
	// Code is the existing concise display label retained for renderer/API compatibility.
	Code string `json:"code,omitempty"`
	// ReasonCode is the canonical ADR-GDI reason taxonomy used for diagnosis,
	// persistence and future adaptive reliability learning.
	ReasonCode        string   `json:"reasonCode,omitempty"`
	PressureState     string   `json:"pressureState,omitempty"`
	Detail            string   `json:"detail,omitempty"`
	DecisionImpact    string   `json:"decisionImpact,omitempty"`
	CriticalUsable    bool     `json:"criticalUsable"`
	Abstain           bool     `json:"abstain,omitempty"`
	Affected          []string `json:"affected,omitempty"`
	AffectedConsumers []string `json:"affectedConsumers,omitempty"`
}

func criticalDecisionDataUsable(freshness []FreshnessDiagnostic, session string) bool {
	quoteUsable := false
	vixUsable := false
	for _, row := range freshness {
		state := strings.ToUpper(strings.TrimSpace(row.State))
		switch row.Dataset {
		case "Quotes":
			quoteUsable = state == "LIVE" || state == "FRESH" || state == "DUE SOON" || (session != "regular" && (state == "DELAYED" || state == "IDLE"))
		case "VIX":
			vixUsable = state == "LIVE" || state == "FRESH" || state == "DUE SOON" || state == "DELAYED" || state == "IDLE"
		}
	}
	return quoteUsable && vixUsable
}

func reliabilityConsumersFor(affected []string) []string {
	set := map[string]bool{}
	add := func(values ...string) {
		for _, value := range values {
			if strings.TrimSpace(value) != "" {
				set[value] = true
			}
		}
	}
	for _, item := range affected {
		s := strings.ToLower(strings.TrimSpace(item))
		switch {
		case strings.Contains(s, "live equities"), strings.Contains(s, "quote"):
			add("Day", "Swing", "Decision Queue", "Rapid Move", "Market Open Prep")
		case strings.Contains(s, "vix"), strings.Contains(s, "indice"):
			add("Market Regime", "Day", "Swing", "Pre-Market Prep", "Market Open Prep")
		case strings.Contains(s, "intraday"), strings.Contains(s, "bar"), strings.Contains(s, "history"):
			add("Day", "Swing", "Long", "Research")
		case strings.Contains(s, "news"):
			add("Research", "Event Intelligence", "Catalyst Reaction")
		case strings.Contains(s, "earning"):
			add("Research", "Catalyst Reaction", "Pre-Market Prep", "Market Open Prep")
		case strings.Contains(s, "filing"), strings.Contains(s, "sec"):
			add("Research", "Event Intelligence", "Long")
		case strings.Contains(s, "fundamental"):
			add("Long", "Research")
		case strings.Contains(s, "option"):
			add("Swing", "Research")
		case strings.Contains(s, "macro"):
			add("Market Regime", "Day", "Swing", "Long", "Research")
		case strings.Contains(s, "broad-discovery"), strings.Contains(s, "background"):
			add("Opportunity Radar")
		case strings.Contains(s, "provider refresh"):
			add("Research", "Opportunity Radar")
		case strings.Contains(s, "subscription promotions"), strings.Contains(s, "failover"):
			add("Day", "Swing", "Decision Queue", "Rapid Move")
		case strings.Contains(s, "provider request budget"):
			add("Provider-backed refresh consumers")
		}
	}
	out := make([]string, 0, len(set))
	for value := range set {
		out = append(out, value)
	}
	sort.Strings(out)
	return out
}

func finalizeRuntimeDegradation(out RuntimeDegradationState) RuntimeDegradationState {
	if strings.TrimSpace(out.Code) == "" {
		out.PressureState = "HEALTHY"
		out.DecisionImpact = "No active decision-relevant degradation is detected."
		return out
	}
	out.AffectedConsumers = reliabilityConsumersFor(out.Affected)
	if out.CriticalUsable {
		out.PressureState = "PROTECTED"
		out.DecisionImpact = "Critical decision evidence remains usable; degradation is isolated to the listed capability/consumer scope."
		return out
	}
	out.PressureState = "DEGRADED"
	out.Abstain = true
	out.DecisionImpact = "Required decision evidence is not trustworthy enough; affected current/readiness conclusions must ABSTAIN until recovery is proven."
	return out
}

func deriveRuntimeDegradation(status, mode string, feed FeedDiagnostics, freshness []FreshnessDiagnostic, router ProviderRouterSnapshot, load RuntimeLoadDiagnostics) RuntimeDegradationState {
	critical := criticalDecisionDataUsable(freshness, feed.MarketSession)
	out := RuntimeDegradationState{CriticalUsable: critical}
	if mode == "demo" || status == "stopped" {
		return finalizeRuntimeDegradation(out)
	}

	for _, class := range load.Workload {
		if class.Class != "provider-rest" || class.Queued <= 0 {
			continue
		}
		queueSaturated := class.MaxQueue > 0 && class.Queued >= class.MaxQueue
		if queueSaturated || class.OldestQueueAgeMs >= 2000 {
			out.Code = "LOCAL LOAD"
			out.ReasonCode = "LOCAL_OVERLOAD"
			if queueSaturated {
				out.ReasonCode = "QUEUE_SATURATED"
			}
			out.Detail = "Provider work is queued behind the bounded shared request budget"
			out.Affected = []string{"non-streaming provider refreshes"}
			return finalizeRuntimeDegradation(out)
		}
	}

	if strings.EqualFold(status, "degraded") {
		for _, sub := range load.LiveSubscriptions {
			if sub.Saturated {
				out.Code = "LIVE CAPACITY SATURATED"
				out.ReasonCode = "LOCAL_OVERLOAD"
				out.Detail = sub.Provider + " live subscription capacity is fully allocated; reserved headroom is exhausted"
				out.Affected = []string{"live subscription promotions/failover"}
				return finalizeRuntimeDegradation(out)
			}
		}
	}

	for _, provider := range load.ProviderRequests {
		if provider.RateLimited > 0 && provider.RequestsLastMin > 0 {
			out.Code = "RATE LIMITED"
			out.ReasonCode = "RATE_LIMITED"
			out.Detail = provider.Provider + " returned rate-limit pressure in the current request window"
			out.Affected = []string{"provider request budget"}
			return finalizeRuntimeDegradation(out)
		}
	}
	if load.Goroutines > 600 || load.HeapAllocBytes > 768*1024*1024 {
		out.Code = "LOCAL LOAD"
		out.ReasonCode = "LOCAL_OVERLOAD"
		out.Detail = "Local runtime pressure crossed the defensive v17 process threshold; optional/background work should be shed first"
		out.Affected = []string{"background and broad-discovery work"}
		return finalizeRuntimeDegradation(out)
	}

	for _, route := range router.Routes {
		for _, hop := range route.Route {
			if strings.EqualFold(hop.RateLimit, "RATE LIMITED") || strings.EqualFold(hop.Circuit, "RATE LIMITED") {
				out.Code = "RATE LIMITED"
				out.ReasonCode = "RATE_LIMITED"
				out.Detail = hop.Provider + " request budget is temporarily rate limited"
				out.Affected = append(out.Affected, route.Dataset)
				return finalizeRuntimeDegradation(out)
			}
		}
	}

	if feed.FeedState == "reconnecting" {
		out.Code = "NETWORK"
		out.ReasonCode = "NETWORK_FAILURE"
		out.Detail = "Primary and fallback live feeds are reconnecting"
		out.Affected = []string{"live equities"}
		return finalizeRuntimeDegradation(out)
	}
	if feed.FeedState == "finnhub-fallback" {
		out.Code = "PROVIDER DEGRADED"
		out.ReasonCode = "PROVIDER_DOWN"
		out.Detail = "Primary Alpaca feed is unavailable or quiet; Finnhub fallback is carrying live equity updates"
		out.Affected = []string{"primary live-equity route"}
		return finalizeRuntimeDegradation(out)
	}

	bad := []string{}
	good := 0
	for _, row := range freshness {
		state := strings.ToUpper(strings.TrimSpace(row.State))
		switch state {
		case "LIVE", "FRESH", "DUE SOON", "IDLE":
			good++
		case "STALE", "ERROR", "UNAVAILABLE":
			bad = append(bad, row.Dataset)
		}
	}
	if len(bad) > 0 && good > 0 {
		out.Code = "PARTIAL COVERAGE"
		out.ReasonCode = "LOW_COVERAGE"
		out.Detail = strings.Join(bad, ", ") + " need recovery while other canonical datasets remain usable"
		out.Affected = bad
		return finalizeRuntimeDegradation(out)
	}

	for _, route := range router.Routes {
		if strings.EqualFold(route.State, "UNAVAILABLE") {
			out.Code = "PROVIDER DEGRADED"
			out.ReasonCode = "PROVIDER_DOWN"
			out.Detail = route.Dataset + " provider route is unavailable"
			out.Affected = []string{route.Dataset}
			return finalizeRuntimeDegradation(out)
		}
	}

	// Fail closed when required live decision evidence is unusable but the
	// narrower network/provider/load classifiers cannot yet attribute a cause.
	// UNKNOWN is truthful here; absence of a diagnosis must never look healthy.
	if !critical {
		out.Code = "DATA DEGRADED"
		out.ReasonCode = "UNKNOWN"
		out.Detail = "Required decision evidence is stale, unavailable, or otherwise insufficient; no narrower recovery cause has been proven yet"
		out.Affected = bad
		if len(out.Affected) == 0 {
			out.Affected = []string{"Quotes", "VIX"}
		}
		return finalizeRuntimeDegradation(out)
	}

	if status == "degraded" {
		out.Code = "PROVIDER DEGRADED"
		out.ReasonCode = "UNKNOWN"
		out.Detail = "Runtime health is degraded; inspect Provider Router and Data Freshness for the active cause"
	}
	return finalizeRuntimeDegradation(out)
}
