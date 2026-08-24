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
	FallbackActive    bool     `json:"fallbackActive,omitempty"`
	PreferredProvider string   `json:"preferredProvider,omitempty"`
	ServingProvider   string   `json:"servingProvider,omitempty"`
	FallbackStatus    string   `json:"fallbackStatus,omitempty"`
	FallbackDetail    string   `json:"fallbackDetail,omitempty"`
	// WarmStateActive means current canonical evidence is still within freshness
	// policy even though a provider/transport route is temporarily unavailable.
	// It is diagnostic context only and must never extend an evidence timestamp.
	WarmStateActive bool `json:"warmStateActive,omitempty"`
	// TransportIssues records non-secret provider/network context that has not
	// yet invalidated current canonical evidence. These issues become a data
	// health degradation only when required evidence is no longer usable.
	TransportIssues []string `json:"transportIssues,omitempty"`
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

func appendUniqueRuntimeIssue(values []string, value string) []string {
	value = strings.TrimSpace(value)
	if value == "" {
		return values
	}
	for _, existing := range values {
		if strings.EqualFold(strings.TrimSpace(existing), value) {
			return values
		}
	}
	return append(values, value)
}

func finalizeRuntimeDegradation(out RuntimeDegradationState) RuntimeDegradationState {
	if strings.TrimSpace(out.Code) == "" {
		out.PressureState = "HEALTHY"
		if out.WarmStateActive || out.FallbackActive || len(out.TransportIssues) > 0 {
			out.DecisionImpact = "Current canonical evidence remains usable within freshness policy; provider/transport pressure is isolated while fallback or warm evidence remains valid."
		} else {
			out.DecisionImpact = "No active decision-relevant degradation is detected."
		}
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

	for _, sub := range load.LiveSubscriptions {
		if sub.Saturated {
			out.Code = "LIVE CAPACITY SATURATED"
			out.ReasonCode = "LOCAL_OVERLOAD"
			out.Detail = sub.Provider + " live subscription capacity is fully allocated; reserved headroom is exhausted"
			out.Affected = []string{"live subscription promotions/failover"}
			return finalizeRuntimeDegradation(out)
		}
	}

	for _, provider := range load.ProviderRequests {
		if provider.RateLimited <= 0 || provider.RequestsLastMin <= 0 {
			continue
		}
		issue := provider.Provider + " request budget is rate limited"
		out.TransportIssues = appendUniqueRuntimeIssue(out.TransportIssues, issue)
		if critical {
			out.WarmStateActive = true
			continue
		}
		out.Code = "RATE LIMITED"
		out.ReasonCode = "RATE_LIMITED"
		out.Detail = provider.Provider + " returned rate-limit pressure in the current request window"
		out.Affected = []string{"provider request budget"}
		return finalizeRuntimeDegradation(out)
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
			if !strings.EqualFold(hop.RateLimit, "RATE LIMITED") && !strings.EqualFold(hop.Circuit, "RATE LIMITED") {
				continue
			}
			issue := route.Dataset + " · " + hop.Provider + " rate limited"
			out.TransportIssues = appendUniqueRuntimeIssue(out.TransportIssues, issue)
			if critical {
				out.WarmStateActive = true
				continue
			}
			out.Code = "RATE LIMITED"
			out.ReasonCode = "RATE_LIMITED"
			out.Detail = hop.Provider + " request budget is temporarily rate limited"
			out.Affected = append(out.Affected, route.Dataset)
			return finalizeRuntimeDegradation(out)
		}
	}

	if feed.FeedState == "reconnecting" {
		out.TransportIssues = appendUniqueRuntimeIssue(out.TransportIssues, "Primary and fallback live feeds are reconnecting")
		if critical {
			out.WarmStateActive = true
		} else {
			out.Code = "NETWORK"
			out.ReasonCode = "NETWORK_FAILURE"
			out.Detail = "Primary and fallback live feeds are reconnecting"
			out.Affected = []string{"live equities"}
			return finalizeRuntimeDegradation(out)
		}
	}
	if feed.FeedState == "finnhub-fallback" {
		// A functioning fallback is routing context, not by itself a data-health
		// failure. Keep preferred/serving provenance visible and let canonical
		// freshness/route checks below decide whether evidence is actually degraded.
		out.FallbackActive = true
		out.PreferredProvider = "Alpaca"
		out.ServingProvider = "Finnhub"
		out.FallbackStatus = "ACTIVE"
		out.FallbackDetail = "Primary Alpaca live feed is unavailable or quiet; Finnhub fallback is carrying live equity updates"
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

	unavailableRoutes := []string{}
	for _, route := range router.Routes {
		if !strings.EqualFold(route.State, "UNAVAILABLE") {
			continue
		}
		unavailableRoutes = appendUniqueRuntimeIssue(unavailableRoutes, route.Dataset)
		out.TransportIssues = appendUniqueRuntimeIssue(out.TransportIssues, route.Dataset+" provider route unavailable")
	}
	if len(unavailableRoutes) > 0 {
		if critical {
			// Do not stamp or extend evidence time here. Freshness remains the sole
			// authority for how long the already-valid warm state can be reused.
			out.WarmStateActive = true
		} else {
			out.Code = "PROVIDER DEGRADED"
			out.ReasonCode = "PROVIDER_DOWN"
			out.Detail = strings.Join(unavailableRoutes, ", ") + " provider route(s) are unavailable and required canonical evidence is no longer usable"
			out.Affected = unavailableRoutes
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

	// Do not re-create degradation from the mutable runtime status alone. The
	// canonical evidence/router/load evaluation above owns current truth, while
	// RuntimeSLOTracker.StabilizeDegradation owns hysteresis. This is what allows
	// a previously degraded runtime to enter RECOVERING and ultimately unlatch.
	return finalizeRuntimeDegradation(out)
}
