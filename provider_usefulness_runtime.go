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

	e.mu.RLock()
	quotes := clone(e.quotes)
	lastUpdated := clone(e.lastUpdated)
	providerQuotes := clone(e.providerQuotes)
	e.mu.RUnlock()

	router := e.buildProviderRouterSnapshot(settings, secrets, quotes, lastUpdated)
	decisions := buildProviderReconciliation(router, providerQuotes, quotes, now.UnixMilli())
	canonicalProviderUsefulness.bindPersistence(persistence)
	canonicalProviderUsefulness.observe(decisions, now.UnixMilli())
	canonicalProviderUsefulness.persist(persistence, now.UnixMilli())
	return canonicalProviderUsefulness.diagnostics()
}
