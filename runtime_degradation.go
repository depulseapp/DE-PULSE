package main

import "strings"

type RuntimeDegradationState struct {
	Code           string   `json:"code,omitempty"`
	Detail         string   `json:"detail,omitempty"`
	CriticalUsable bool     `json:"criticalUsable"`
	Affected       []string `json:"affected,omitempty"`
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

func deriveRuntimeDegradation(status, mode string, feed FeedDiagnostics, freshness []FreshnessDiagnostic, router ProviderRouterSnapshot, load RuntimeLoadDiagnostics) RuntimeDegradationState {
	critical := criticalDecisionDataUsable(freshness, feed.MarketSession)
	out := RuntimeDegradationState{CriticalUsable: critical}
	if mode == "demo" || status == "stopped" {
		return out
	}

	for _, class := range load.Workload {
		if class.Class == "provider-rest" && class.Queued > 0 && class.OldestQueueAgeMs >= 2000 {
			out.Code = "LOCAL LOAD"
			out.Detail = "Provider work is queued behind the bounded shared request budget"
			out.Affected = []string{"non-streaming provider refreshes"}
			return out
		}
	}

	if strings.EqualFold(status, "degraded") {
		for _, sub := range load.LiveSubscriptions {
			if sub.Saturated {
				out.Code = "LIVE CAPACITY SATURATED"
				out.Detail = sub.Provider + " live subscription capacity is fully allocated; reserved headroom is exhausted"
				out.Affected = []string{"live subscription promotions/failover"}
				return out
			}
		}
	}

	for _, provider := range load.ProviderRequests {
		if provider.RateLimited > 0 && provider.RequestsLastMin > 0 {
			out.Code = "RATE LIMITED"
			out.Detail = provider.Provider + " returned rate-limit pressure in the current request window"
			out.Affected = []string{"provider request budget"}
			return out
		}
	}
	if load.Goroutines > 600 || load.HeapAllocBytes > 768*1024*1024 {
		out.Code = "LOCAL LOAD"
		out.Detail = "Local runtime pressure crossed the defensive v17 process threshold; optional/background work should be shed first"
		out.Affected = []string{"background and broad-discovery work"}
		return out
	}

	for _, route := range router.Routes {
		for _, hop := range route.Route {
			if strings.EqualFold(hop.RateLimit, "RATE LIMITED") || strings.EqualFold(hop.Circuit, "RATE LIMITED") {
				out.Code = "RATE LIMITED"
				out.Detail = hop.Provider + " request budget is temporarily rate limited"
				out.Affected = append(out.Affected, route.Dataset)
				return out
			}
		}
	}

	if feed.FeedState == "reconnecting" {
		out.Code = "NETWORK"
		out.Detail = "Primary and fallback live feeds are reconnecting"
		out.Affected = []string{"live equities"}
		return out
	}
	if feed.FeedState == "finnhub-fallback" {
		out.Code = "PROVIDER DEGRADED"
		out.Detail = "Primary Alpaca feed is unavailable or quiet; Finnhub fallback is carrying live equity updates"
		out.Affected = []string{"primary live-equity route"}
		return out
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
		out.Detail = strings.Join(bad, ", ") + " need recovery while other canonical datasets remain usable"
		out.Affected = bad
		return out
	}

	for _, route := range router.Routes {
		if strings.EqualFold(route.State, "UNAVAILABLE") {
			out.Code = "PROVIDER DEGRADED"
			out.Detail = route.Dataset + " provider route is unavailable"
			out.Affected = []string{route.Dataset}
			return out
		}
	}
	if status == "degraded" {
		out.Code = "PROVIDER DEGRADED"
		out.Detail = "Runtime health is degraded; inspect Provider Router and Data Freshness for the active cause"
	}
	return out
}
