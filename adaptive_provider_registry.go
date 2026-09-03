package main

import (
	"sort"
	"strings"
)

const adaptiveProviderRegistryContractVersion = "adaptive-provider-registry-v19.1.0"

const (
	providerAuthorityRoutable        = "ROUTABLE"
	providerAuthorityFallback        = "FALLBACK"
	providerAuthorityCorroborative   = "CORROBORATIVE"
	providerAuthorityDirectAuthority = "DIRECT_AUTHORITY"
	providerAuthorityNonRoutable     = "NON_ROUTABLE"

	providerRegistryCapabilityRegistered              = "REGISTERED"
	providerRegistryCapabilityRegisteredNotProduction = "REGISTERED_NOT_PRODUCTION"
)

// AdaptiveProviderRouteManifest is a read-only projection of the existing
// ProviderRegistration route contract. It does not own routing, health,
// freshness, lifecycle, rights, cache, persistence, subscriptions or canonical
// state. Smart Provider Router v2 remains the sole general selection authority.
type AdaptiveProviderRouteManifest struct {
	Dataset           string   `json:"dataset"`
	Capability        string   `json:"capability"`
	InstrumentClass   string   `json:"instrumentClass"`
	AuthorityClass    string   `json:"authorityClass"`
	Lifecycle         string   `json:"lifecycle"`
	ContractVersion   string   `json:"contractVersion"`
	CanonicalOwner    string   `json:"canonicalOwner"`
	SchemaContract    string   `json:"schemaContract"`
	TimestampContract string   `json:"timestampContract"`
	FreshnessContract string   `json:"freshnessContract"`
	RightsContract    string   `json:"rightsContract"`
	ExpectedDelay     string   `json:"expectedDelay,omitempty"`
	Uses              []string `json:"uses,omitempty"`
	ProductionReady   bool     `json:"productionReady"`
}

// AdaptiveProviderDiagnosticManifest exposes registered diagnostic capability
// metadata without exposing the diagnostic function or any credential material.
type AdaptiveProviderDiagnosticManifest struct {
	Capability string   `json:"capability"`
	Detail     string   `json:"detail,omitempty"`
	Uses       []string `json:"uses,omitempty"`
}

// AdaptiveProviderManifest is the generic manifest projected from the existing
// canonical ProviderRegistration descriptor. Configuration is represented only
// as presence/absence; secrets are never copied into this projection.
type AdaptiveProviderManifest struct {
	ProviderID       string                               `json:"providerId"`
	Name             string                               `json:"name"`
	ContractVersion  string                               `json:"contractVersion"`
	Configured       bool                                 `json:"configured"`
	Configuration    string                               `json:"configuration"`
	QuotaLabel       string                               `json:"quotaLabel,omitempty"`
	CostClass        string                               `json:"costClass,omitempty"`
	Routable         bool                                 `json:"routable"`
	Routes           []AdaptiveProviderRouteManifest     `json:"routes,omitempty"`
	Diagnostics      []AdaptiveProviderDiagnosticManifest `json:"diagnostics,omitempty"`
}

// AdaptiveProviderRegistrySnapshot is intentionally a computed snapshot rather
// than a second runtime state owner. All entries come from providerRegistrations
// (or an injected registration slice in tests); Router eligibility remains in
// Smart Provider Router v2 and Data Health/freshness remains separately owned.
type AdaptiveProviderRegistrySnapshot struct {
	ContractVersion string                     `json:"contractVersion"`
	Providers       []AdaptiveProviderManifest `json:"providers"`
}

func adaptiveProviderAuthorityClass(provider string, route ProviderDatasetContract) string {
	switch {
	case strings.EqualFold(provider, "SEC EDGAR") && strings.EqualFold(route.Dataset, "SEC"):
		return providerAuthorityDirectAuthority
	case strings.EqualFold(provider, "SEC"):
		return providerAuthorityCorroborative
	case strings.EqualFold(provider, "CBOE"):
		return providerAuthorityCorroborative
	case strings.EqualFold(provider, "yfinance"):
		return providerAuthorityFallback
	default:
		return providerAuthorityRoutable
	}
}

func adaptiveProviderRouteManifest(provider string, route ProviderDatasetContract) AdaptiveProviderRouteManifest {
	return AdaptiveProviderRouteManifest{
		Dataset:           route.Dataset,
		Capability:        route.Capability,
		InstrumentClass:   providerInstrumentClass(route.Dataset),
		AuthorityClass:    adaptiveProviderAuthorityClass(provider, route),
		Lifecycle:         route.Lifecycle,
		ContractVersion:   route.ContractVersion,
		CanonicalOwner:    route.CanonicalOwner,
		SchemaContract:    route.SchemaContract,
		TimestampContract: route.TimestampContract,
		FreshnessContract: route.FreshnessContract,
		RightsContract:    route.RightsContract,
		ExpectedDelay:     route.ExpectedDelay,
		Uses:              append([]string(nil), route.Uses...),
		ProductionReady:   providerDatasetContractProductionReady(route),
	}
}

func adaptiveProviderManifest(reg ProviderRegistration, settings Settings, secrets Secrets) AdaptiveProviderManifest {
	configured := reg.Configured != nil && reg.Configured(settings, secrets)
	manifest := AdaptiveProviderManifest{
		ProviderID:      providerKey(reg.Name),
		Name:            reg.Name,
		ContractVersion: adaptiveProviderRegistryContractVersion,
		Configured:      configured,
		Configuration:   providerCapabilityNotConfigured,
		QuotaLabel:      reg.QuotaLabel,
		CostClass:       reg.CostClass,
		Routable:        len(reg.Routes) > 0,
	}
	if configured {
		manifest.Configuration = "CONFIGURED"
	}
	for _, route := range reg.Routes {
		manifest.Routes = append(manifest.Routes, adaptiveProviderRouteManifest(reg.Name, route))
	}
	for _, diagnostic := range reg.Diagnostics {
		manifest.Diagnostics = append(manifest.Diagnostics, AdaptiveProviderDiagnosticManifest{
			Capability: diagnostic.Capability,
			Detail:     diagnostic.Detail,
			Uses:       append([]string(nil), diagnostic.Uses...),
		})
	}
	if len(manifest.Routes) == 0 {
		manifest.Routable = false
	}
	sort.SliceStable(manifest.Routes, func(i, j int) bool {
		if manifest.Routes[i].Dataset != manifest.Routes[j].Dataset {
			return manifest.Routes[i].Dataset < manifest.Routes[j].Dataset
		}
		if manifest.Routes[i].Capability != manifest.Routes[j].Capability {
			return manifest.Routes[i].Capability < manifest.Routes[j].Capability
		}
		return manifest.Routes[i].Lifecycle < manifest.Routes[j].Lifecycle
	})
	sort.SliceStable(manifest.Diagnostics, func(i, j int) bool {
		return strings.ToLower(manifest.Diagnostics[i].Capability) < strings.ToLower(manifest.Diagnostics[j].Capability)
	})
	return manifest
}

func adaptiveProviderRegistrySnapshotFromRegistrations(regs []ProviderRegistration, settings Settings, secrets Secrets) AdaptiveProviderRegistrySnapshot {
	out := AdaptiveProviderRegistrySnapshot{ContractVersion: adaptiveProviderRegistryContractVersion}
	for _, reg := range regs {
		if strings.TrimSpace(reg.Name) == "" {
			continue
		}
		out.Providers = append(out.Providers, adaptiveProviderManifest(reg, settings, secrets))
	}
	sort.SliceStable(out.Providers, func(i, j int) bool {
		return out.Providers[i].ProviderID < out.Providers[j].ProviderID
	})
	return out
}

func adaptiveProviderRegistrySnapshot(settings Settings, secrets Secrets) AdaptiveProviderRegistrySnapshot {
	return adaptiveProviderRegistrySnapshotFromRegistrations(providerRegistrations(), settings, secrets)
}

// adaptiveProviderRegistryCapabilityState answers only registration/configuration
// truth. It deliberately does not answer Router eligibility or provider health.
// Those remain Smart Provider Router v2/Data Health responsibilities.
func adaptiveProviderRegistryCapabilityStateFromRegistrations(regs []ProviderRegistration, provider, dataset string, settings Settings, secrets Secrets) string {
	wantProvider := providerKey(provider)
	wantDataset := strings.TrimSpace(dataset)
	for _, reg := range regs {
		if providerKey(reg.Name) != wantProvider {
			continue
		}
		for _, route := range reg.Routes {
			if !strings.EqualFold(strings.TrimSpace(route.Dataset), wantDataset) {
				continue
			}
			if reg.Configured == nil || !reg.Configured(settings, secrets) {
				return providerCapabilityNotConfigured
			}
			if !providerDatasetContractProductionReady(route) {
				return providerRegistryCapabilityRegisteredNotProduction
			}
			return providerRegistryCapabilityRegistered
		}
		return providerCapabilityNotSupported
	}
	return providerCapabilityNotSupported
}

func adaptiveProviderRegistryCapabilityState(provider, dataset string, settings Settings, secrets Secrets) string {
	return adaptiveProviderRegistryCapabilityStateFromRegistrations(providerRegistrations(), provider, dataset, settings, secrets)
}
