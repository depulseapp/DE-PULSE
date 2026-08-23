package adaptivepolicy

import (
	"strings"
	"time"
)

// IntradayHistoryCadence returns the bounded history-refresh cadence for the
// current market session and symbol heat state.
func IntradayHistoryCadence(session string, hot bool) time.Duration {
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

// CachePersistCadence returns the bounded persistence cadence for the current
// market session and symbol heat state.
func CachePersistCadence(session string, hot bool) time.Duration {
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

// ProviderDegraded classifies the canonical provider-health fields used by
// adaptive policy. The caller remains responsible for canonical state ownership.
func ProviderDegraded(health map[string]string) bool {
	h := strings.ToLower(health["alpaca-live"] + " " + health["alpaca-stream"] + " " + health["quotes"])
	return strings.Contains(h, "error") || strings.Contains(h, "failed") || strings.Contains(h, "reconnecting") || strings.Contains(h, "degraded")
}
