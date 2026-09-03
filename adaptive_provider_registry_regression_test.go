package main

import (
	"encoding/json"
	"reflect"
	"strings"
	"testing"

	"depulse/internal/providerlifecycle"
)

func findAdaptiveProviderManifest(snapshot AdaptiveProviderRegistrySnapshot, provider string) (AdaptiveProviderManifest, bool) {
	want := providerKey(provider)
	for _, row := range snapshot.Providers {
		if row.ProviderID == want {
			return row, true
		}
	}
	return AdaptiveProviderManifest{}, false
}

func findAdaptiveProviderRoute(manifest AdaptiveProviderManifest, dataset, capability string) (AdaptiveProviderRouteManifest, bool) {
	for _, row := range manifest.Routes {
		if row.Dataset == dataset && row.Capability == capability {
			return row, true
		}
	}
	return AdaptiveProviderRouteManifest{}, false
}

func TestAdaptiveProviderRegistryProjectsEveryCanonicalRegistration(t *testing.T) {
	settings := defaultState().Settings
	settings.SECEmail = "depulse@example.com"
	secrets := Secrets{
		Finnhub: "fh", TradeInsight: "ti", AlpacaKey: "ak", AlpacaSecret: "as",
		FRED: "fred", EIA: "eia", TwelveData: "td", Marketaux: "ma",
	}
	regs := providerRegistrations()
	snapshot := adaptiveProviderRegistrySnapshotFromRegistrations(regs, settings, secrets)
	if snapshot.ContractVersion != adaptiveProviderRegistryContractVersion {
		t.Fatalf("registry contract version mismatch: %q", snapshot.ContractVersion)
	}
	if len(snapshot.Providers) != len(regs) {
		t.Fatalf("registry projection lost registrations: got=%d want=%d", len(snapshot.Providers), len(regs))
	}
	seen := map[string]bool{}
	for _, reg := range regs {
		manifest, ok := findAdaptiveProviderManifest(snapshot, reg.Name)
		if !ok {
			t.Fatalf("registration missing from adaptive registry projection: %s", reg.Name)
		}
		if seen[manifest.ProviderID] {
			t.Fatalf("duplicate provider id in registry: %s", manifest.ProviderID)
		}
		seen[manifest.ProviderID] = true
		if manifest.Routable != (len(reg.Routes) > 0) {
			t.Fatalf("routable projection mismatch for %s: %+v", reg.Name, manifest)
		}
		for _, route := range reg.Routes {
			projected, ok := findAdaptiveProviderRoute(manifest, route.Dataset, route.Capability)
			if !ok {
				t.Fatalf("route missing from registry projection: %s / %s / %s", reg.Name, route.Dataset, route.Capability)
			}
			if projected.Lifecycle != route.Lifecycle || projected.CanonicalOwner != route.CanonicalOwner || projected.ContractVersion != route.ContractVersion {
				t.Fatalf("registry rewrote governed route contract for %s / %s: got=%+v source=%+v", reg.Name, route.Dataset, projected, route)
			}
			if projected.ProductionReady != providerDatasetContractProductionReady(route) {
				t.Fatalf("production readiness projection drift for %s / %s", reg.Name, route.Dataset)
			}
		}
	}
}

func TestAdaptiveProviderRegistryRegistrationIsGenericAndDoesNotNeedParallelMaps(t *testing.T) {
	route := inheritedProductionRoute("Synthetic Adapter", "Synthetic Dataset", "Synthetic Capability", 1, "Synthetic cadence", "Synthetic Consumer")
	regs := []ProviderRegistration{{
		Name:       "Synthetic Adapter",
		QuotaLabel: "Synthetic quota",
		CostClass:  "Synthetic cost",
		Configured: func(Settings, Secrets) bool { return true },
		Routes:     []ProviderDatasetContract{route},
		Diagnostics: []ProviderCapabilityDiagnosticRegistration{{
			Capability: "Synthetic diagnostics", Detail: "Synthetic detail", Uses: []string{"Synthetic Consumer"},
		}},
	}}

	snapshot := adaptiveProviderRegistrySnapshotFromRegistrations(regs, defaultState().Settings, Secrets{})
	manifest, ok := findAdaptiveProviderManifest(snapshot, "Synthetic Adapter")
	if !ok || !manifest.Configured || manifest.Configuration != "CONFIGURED" {
		t.Fatalf("synthetic adapter did not self-register through generic descriptor: %+v", manifest)
	}
	if len(manifest.Routes) != 1 || len(manifest.Diagnostics) != 1 {
		t.Fatalf("generic registration projection lost route/diagnostic metadata: %+v", manifest)
	}
	if got := routeChainsFromProviderRegistrations(regs)["Synthetic Dataset"]; !reflect.DeepEqual(got, []string{"Synthetic Adapter"}) {
		t.Fatalf("existing route construction did not discover synthetic registration: %v", got)
	}
}

func TestAdaptiveProviderRegistryCapabilityTruthIsBoundedBeforeRouterEligibility(t *testing.T) {
	settings := defaultState().Settings
	if got := adaptiveProviderRegistryCapabilityState("Finnhub", "News", settings, Secrets{}); got != providerCapabilityNotConfigured {
		t.Fatalf("unconfigured registered capability truth mismatch: %s", got)
	}
	if got := adaptiveProviderRegistryCapabilityState("Finnhub", "Options Chain", settings, Secrets{Finnhub: "configured"}); got != providerCapabilityNotSupported {
		t.Fatalf("unsupported capability truth mismatch: %s", got)
	}
	if got := adaptiveProviderRegistryCapabilityState("Unknown Provider", "News", settings, Secrets{}); got != providerCapabilityNotSupported {
		t.Fatalf("unknown provider must fail closed as unsupported: %s", got)
	}
	if got := adaptiveProviderRegistryCapabilityState("Finnhub", "News", settings, Secrets{Finnhub: "configured"}); got != providerRegistryCapabilityRegistered {
		t.Fatalf("configured governed registration should be REGISTERED, not a routing decision: %s", got)
	}
}

func TestAdaptiveProviderRegistryNeverPromotesLifecycleIntoProduction(t *testing.T) {
	route := inheritedProductionRoute("Synthetic Shadow", "Synthetic Dataset", "Synthetic Capability", 1, "Synthetic cadence", "Synthetic Consumer")
	route.Lifecycle = providerlifecycle.Shadow
	regs := []ProviderRegistration{{Name: "Synthetic Shadow", Configured: func(Settings, Secrets) bool { return true }, Routes: []ProviderDatasetContract{route}}}
	settings := defaultState().Settings

	if got := adaptiveProviderRegistryCapabilityStateFromRegistrations(regs, "Synthetic Shadow", "Synthetic Dataset", settings, Secrets{}); got != providerRegistryCapabilityRegisteredNotProduction {
		t.Fatalf("SHADOW route was not held below production: %s", got)
	}
	if got := routeChainsFromProviderRegistrations(regs)["Synthetic Dataset"]; len(got) != 0 {
		t.Fatalf("registry projection must not promote SHADOW route into canonical production chain: %v", got)
	}
	manifest, ok := findAdaptiveProviderManifest(adaptiveProviderRegistrySnapshotFromRegistrations(regs, settings, Secrets{}), "Synthetic Shadow")
	if !ok || len(manifest.Routes) != 1 || manifest.Routes[0].ProductionReady {
		t.Fatalf("manifest falsely reported SHADOW route production-ready: %+v", manifest)
	}
}

func TestAdaptiveProviderRegistryPreservesExplicitAuthorityClasses(t *testing.T) {
	snapshot := adaptiveProviderRegistrySnapshot(defaultState().Settings, Secrets{})
	sec, ok := findAdaptiveProviderManifest(snapshot, "SEC EDGAR")
	if !ok {
		t.Fatal("SEC EDGAR registration missing")
	}
	secRoute, ok := findAdaptiveProviderRoute(sec, "SEC", "Direct SEC/EDGAR filings and Form 4 authority")
	if !ok || secRoute.AuthorityClass != providerAuthorityDirectAuthority {
		t.Fatalf("direct SEC/EDGAR authority was not preserved: %+v", secRoute)
	}
	yf, ok := findAdaptiveProviderManifest(snapshot, "yfinance")
	if !ok {
		t.Fatal("yfinance registration missing")
	}
	for _, route := range yf.Routes {
		if route.AuthorityClass != providerAuthorityFallback {
			t.Fatalf("yfinance must remain fallback-only in registry projection: %+v", route)
		}
	}
	cboe, ok := findAdaptiveProviderManifest(snapshot, "CBOE")
	if !ok || len(cboe.Routes) == 0 || cboe.Routes[0].AuthorityClass != providerAuthorityCorroborative {
		t.Fatalf("CBOE corroborative role was not preserved: %+v", cboe)
	}
}

func TestAdaptiveProviderRegistryProjectionNeverContainsSecretMaterial(t *testing.T) {
	secret := "APR01-super-secret-material"
	snapshot := adaptiveProviderRegistrySnapshot(defaultState().Settings, Secrets{Finnhub: secret, TwelveData: secret, Marketaux: secret})
	raw, err := json.Marshal(snapshot)
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(string(raw), secret) {
		t.Fatalf("registry projection leaked provider credential material: %s", raw)
	}
	finnhub, ok := findAdaptiveProviderManifest(snapshot, "Finnhub")
	if !ok || !finnhub.Configured || finnhub.Configuration != "CONFIGURED" {
		t.Fatalf("redacted configuration presence was not preserved: %+v", finnhub)
	}
}
