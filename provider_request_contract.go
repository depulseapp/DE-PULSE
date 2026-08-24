package main

import (
	"context"
	"strings"
)

// providerRequestFailureIsLocalNeutral identifies request outcomes that are
// caused by the caller or DE.PULSE's own admission/backpressure controls rather
// than by provider health. Router migrations must not turn these outcomes into
// provider or capability circuit failures.
func providerRequestFailureIsLocalNeutral(ctx context.Context, err error) bool {
	if err == nil {
		return false
	}
	if ctx != nil && ctx.Err() != nil {
		return true
	}
	low := strings.ToLower(strings.TrimSpace(err.Error()))
	return strings.Contains(low, "request deferred") ||
		strings.Contains(low, "deferred: provider ") ||
		strings.Contains(low, "provider request rejected by bounded workload budget") ||
		strings.Contains(low, "bounded provider capacity")
}
