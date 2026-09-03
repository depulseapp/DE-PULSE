package main

import (
	"strings"
	"sync"
	"time"
)

// providerConfigurationObservations is process-local only. Values are one-way
// fingerprints produced by the canonical registration contract; credentials
// never enter Router snapshots, persistence, telemetry, logs or diagnostics.
var providerConfigurationObservations sync.Map // map[*Engine]map[string][32]byte

func providerConfigurationSnapshot(settings Settings, secrets Secrets) map[string][32]byte {
	out := map[string][32]byte{}
	registry := adaptiveProviderRegistrySnapshot(settings, secrets)
	registered := make(map[string]bool, len(registry.Providers))
	for _, manifest := range registry.Providers {
		registered[manifest.ProviderID] = true
	}
	for _, reg := range providerRegistrations() {
		if !registered[providerKey(reg.Name)] || reg.ConfigurationFingerprint == nil {
			continue
		}
		out[reg.Name] = reg.ConfigurationFingerprint(settings, secrets)
	}
	return out
}

func entitlementConfigurationState(state string) bool {
	return state == providerCapabilityNotEntitled || state == providerCapabilityNotConfigured
}

func (e *Engine) reopenProviderEntitlementStates(changed map[string]bool, reason string) []string {
	if e == nil || len(changed) == 0 {
		return nil
	}
	now := time.Now().UnixMilli()
	persist := make([]ProviderCapabilityStateRecord, 0)
	changedProviders := make([]string, 0, len(changed))

	e.mu.Lock()
	for key, record := range e.providerCapabilityStates {
		if !changed[providerKey(record.Provider)] || !entitlementConfigurationState(record.State) {
			continue
		}
		record.State = providerCapabilityUnknown
		record.Reason = reason
		record.LastObservedAt = now
		record.RevalidateAt = 0
		record.PolicyVersion = smartRouterPolicyVersion
		e.providerCapabilityStates[key] = record
		persist = append(persist, record)

		circuitKey := providerCapabilityCircuitKey(record.Provider, record.Dataset)
		circuit := e.providerCapabilityCircuits[circuitKey]
		circuit.Failures = 0
		circuit.OpenUntil = 0
		circuit.RateLimitedUntil = 0
		circuit.LastError = ""
		e.providerCapabilityCircuits[circuitKey] = circuit
	}

	for _, reg := range providerRegistrations() {
		if !changed[providerKey(reg.Name)] {
			continue
		}
		changedProviders = append(changedProviders, reg.Name)
		globalKey := providerKey(reg.Name)
		circuit := e.providerCircuits[globalKey]
		if strings.TrimSpace(circuit.LastError) == "" {
			continue
		}
		state, _, _ := classifyProviderCapabilityFailure(&providerConfigurationError{message: circuit.LastError})
		if state == providerCapabilityNotEntitled {
			circuit.Failures = 0
			circuit.OpenUntil = 0
			circuit.RateLimitedUntil = 0
			circuit.LastError = ""
			e.providerCircuits[globalKey] = circuit
		}
	}
	e.mu.Unlock()

	for _, record := range persist {
		e.persistProviderCapabilityState(record)
	}
	return changedProviders
}

// refreshProviderConfigurationEntitlements makes credential/configuration
// changes visible to Smart Provider Router v2 without polling. It is called at
// canonical route-decision time before ranking. Only stale entitlement/config
// suppression is reopened; SUPPORTED, NOT_SUPPORTED and transient health states
// are retained because they represent different evidence.
//
// On the first route decision after process start, resettable persisted records
// are reopened once. This safely handles an offline credential change between
// runs while bounding the extra provider work to at most one fresh probe per
// affected capability after restart.
func (e *Engine) refreshProviderConfigurationEntitlements(settings Settings, secrets Secrets) []string {
	if e == nil {
		return nil
	}
	current := providerConfigurationSnapshot(settings, secrets)
	previousAny, loaded := providerConfigurationObservations.Load(e)
	previous, _ := previousAny.(map[string][32]byte)
	providerConfigurationObservations.Store(e, current)

	changed := map[string]bool{}
	for provider, fingerprint := range current {
		if !loaded || previous[provider] != fingerprint {
			changed[providerKey(provider)] = true
		}
	}
	return e.reopenProviderEntitlementStates(changed, "provider configuration changed; entitlement requires fresh evidence")
}

// forceProviderEntitlementRevalidation is the bounded same-key plan-upgrade
// path used by the explicit capability recheck action. It does not poll and it
// does not erase health evidence: only configured providers with resettable
// NOT_ENTITLED/NOT_CONFIGURED state are reopened for one fresh canonical probe.
func (e *Engine) forceProviderEntitlementRevalidation(settings Settings, secrets Secrets) []string {
	if e == nil {
		return nil
	}
	providerConfigurationObservations.Store(e, providerConfigurationSnapshot(settings, secrets))
	changed := map[string]bool{}
	registry := adaptiveProviderRegistrySnapshot(settings, secrets)
	for _, manifest := range registry.Providers {
		if !manifest.Configured {
			continue
		}
		consistent := true
		for _, route := range manifest.Routes {
			if adaptiveProviderRegistryCapabilityState(manifest.Name, route.Dataset, settings, secrets) == providerCapabilityNotConfigured {
				consistent = false
				break
			}
		}
		if consistent {
			changed[manifest.ProviderID] = true
		}
	}
	return e.reopenProviderEntitlementStates(changed, "manual capability recheck requested; entitlement requires fresh evidence")
}

type providerConfigurationError struct{ message string }

func (e *providerConfigurationError) Error() string { return e.message }
