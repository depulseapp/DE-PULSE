package main

import (
	"sort"
	"strings"
	"time"
)

type AdaptiveDataPolicyState struct {
	Session                string   `json:"session"`
	HotSymbols             []string `json:"hotSymbols"`
	RadarCadenceMs         int64    `json:"radarCadenceMs"`
	IntradayHistoryCadence int64    `json:"intradayHistoryCadenceMs"`
	CachePersistCadence    int64    `json:"cachePersistCadenceMs"`
	ProviderState          string   `json:"providerState"`
	Policy                 string   `json:"policy"`
	UpdatedAt              int64    `json:"updatedAt"`
}

type ShadowExperiment struct {
	Key                 string `json:"key"`
	Stage               string `json:"stage"`
	ProductionValue     string `json:"productionValue"`
	ShadowValue         string `json:"shadowValue"`
	Observation         string `json:"observation"`
	CanMutateProduction bool   `json:"canMutateProduction"`
}

type ShadowControlState struct {
	PromotionPath string             `json:"promotionPath"`
	Experiments   []ShadowExperiment `json:"experiments"`
	UpdatedAt     int64              `json:"updatedAt"`
}

func adaptiveIntradayHistoryCadence(session string, hot bool) time.Duration {
	switch session {
	case "regular":
		if hot {
			return 2 * time.Minute
		}
		return 5 * time.Minute
	case "pre-market", "after-hours":
		if hot {
			return 3 * time.Minute
		}
		return 5 * time.Minute
	case "overnight":
		return 15 * time.Minute
	default:
		return 30 * time.Minute
	}
}

func adaptiveCachePersistCadence(session string, hot bool) time.Duration {
	switch session {
	case "regular", "pre-market", "after-hours":
		if hot {
			return time.Minute
		}
		return 2 * time.Minute
	case "overnight":
		return 5 * time.Minute
	default:
		return 10 * time.Minute
	}
}

func adaptiveProviderDegraded(health map[string]string) bool {
	h := strings.ToLower(health["alpaca-live"] + " " + health["alpaca-stream"] + " " + health["quotes"])
	return strings.Contains(h, "error") || strings.Contains(h, "failed") || strings.Contains(h, "reconnecting") || strings.Contains(h, "degraded")
}

func buildAdaptiveDataPolicyState(scanner ScannerState, health map[string]string, now time.Time) AdaptiveDataPolicyState {
	session := marketSessionET(now)
	hot := make([]string, 0, len(scanner.Radar.Promotions))
	for _, p := range scanner.Radar.Promotions {
		if sym := normalizeSymbol(p.Symbol); sym != "" {
			hot = append(hot, sym)
		}
	}
	hot = uniqueSymbols(hot)
	sort.Strings(hot)
	degraded := adaptiveProviderDegraded(health)
	radar := opportunityRadarCadence(session, len(hot) > 0, degraded)
	history := adaptiveIntradayHistoryCadence(session, len(hot) > 0)
	cache := adaptiveCachePersistCadence(session, len(hot) > 0)
	provider := "HEALTHY"
	if degraded {
		provider = "DEGRADED"
	}
	return AdaptiveDataPolicyState{
		Session:                session,
		HotSymbols:             hot,
		RadarCadenceMs:         int64(radar / time.Millisecond),
		IntradayHistoryCadence: int64(history / time.Millisecond),
		CachePersistCadence:    int64(cache / time.Millisecond),
		ProviderState:          provider,
		Policy:                 "Session + symbol role + Opportunity Radar state + provider health; bounded and observable.",
		UpdatedAt:              now.UnixMilli(),
	}
}

func buildShadowControlState(scanner ScannerState, now time.Time) ShadowControlState {
	shadowMatches := 0
	prodMatches := len(scanner.Radar.Promotions)
	for _, x := range scanner.Radar.Candidates {
		if opportunityPromotionEligible(x, opportunityShadowFloor) {
			shadowMatches++
		}
	}
	return ShadowControlState{
		PromotionPath: "SHADOW → VALIDATED → APPROVED → PRODUCTION",
		Experiments: []ShadowExperiment{
			{Key: "opportunity-promotion-floor", Stage: "SHADOW", ProductionValue: "78", ShadowValue: "72", Observation: formatShadowObservation(prodMatches, shadowMatches), CanMutateProduction: false},
			{Key: "adaptive-freshness-tightening", Stage: "SHADOW", ProductionValue: "bounded v16.10 policy", ShadowValue: "more aggressive hot-symbol cadence", Observation: "Observe provider cost/freshness benefit before any tighter production policy.", CanMutateProduction: false},
		},
		UpdatedAt: now.UnixMilli(),
	}
}

func formatShadowObservation(prod, shadow int) string {
	return "Production promotions " + itoa(prod) + " · shadow-only candidates " + itoa(shadow) + "; shadow cannot change live allocation."
}

func itoa(v int) string {
	if v == 0 {
		return "0"
	}
	const digits = "0123456789"
	buf := [20]byte{}
	i := len(buf)
	for v > 0 {
		i--
		buf[i] = digits[v%10]
		v /= 10
	}
	return string(buf[i:])
}

func (e *Engine) currentAdaptivePolicy(now time.Time) AdaptiveDataPolicyState {
	e.mu.RLock()
	scanner := clone(e.scanner)
	health := clone(e.health)
	e.mu.RUnlock()
	return buildAdaptiveDataPolicyState(scanner, health, now)
}
