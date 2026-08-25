package main

import "time"

// sampleProviderUsefulness reuses the canonical Provider Router snapshot and
// buildProviderReconciliation truth. It records advisory observations only;
// it never mutates route order, capability state, freshness, subscriptions or
// any provider lifecycle decision.
func (e *Engine) sampleProviderUsefulness(now time.Time) []ProviderUsefulnessDiagnostic {
	if e == nil || e.app == nil {
		return canonicalProviderUsefulness.diagnostics()
	}
	e.app.mu.RLock()
	settings := clone(e.app.processingStateLocked().Settings)
	secrets := clone(e.app.secrets)
	persistence := e.app.persistence
	e.app.mu.RUnlock()

	// buildProviderRouterSnapshot is a read-only snapshot builder but, like the
	// canonical Engine.Snapshot path, expects Engine state to remain read-locked
	// while it inspects circuit/health/session fields.
	e.mu.RLock()
	quotes := clone(e.quotes)
	lastUpdated := clone(e.lastUpdated)
	providerQuotes := clone(e.providerQuotes)
	router := e.buildProviderRouterSnapshot(settings, secrets, quotes, lastUpdated)
	e.mu.RUnlock()

	decisions := buildProviderReconciliation(router, providerQuotes, quotes, now.UnixMilli())
	canonicalProviderUsefulness.bindPersistence(persistence)
	canonicalProviderUsefulness.observe(decisions, now.UnixMilli())
	canonicalProviderUsefulness.persist(persistence, now.UnixMilli())
	return canonicalProviderUsefulness.diagnostics()
}
