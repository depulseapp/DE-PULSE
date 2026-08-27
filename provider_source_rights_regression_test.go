package main

import (
	"testing"
	"time"
)

func TestHOST002ProviderRightsSourceResolverCoversCanonicalRegistrations(t *testing.T) {
	seen := map[string]bool{}
	for _, registration := range providerRegistrations() {
		if seen[registration.Name] {
			t.Fatalf("duplicate provider registration name: %q", registration.Name)
		}
		seen[registration.Name] = true
		if got := providerRightsSourceProvider(registration.Name); got != registration.Name {
			t.Fatalf("canonical provider source %q resolved to %q", registration.Name, got)
		}
	}
}

func TestHOST002ProviderRightsSourceResolverPreservesAliasesAndSECBoundary(t *testing.T) {
	cases := map[string]string{
		"Alpaca IEX":                         "Alpaca",
		"Finnhub quote":                      "Finnhub",
		"TradeInsight historical":            "TradeInsight",
		"TwelveData historical":              "Twelve Data",
		"Marketaux news":                     "Marketaux",
		"FRED macro":                         "FRED",
		"SEC fundamentals":                   "SEC",
		"SEC EDGAR filing":                   "SEC EDGAR",
		"EDGAR Form 4":                       "SEC EDGAR",
		"Yahoo Finance":                      "yfinance",
		"CBOE VIX":                           "CBOE",
		"BLS.gov CPI":                        "BLS",
		"Bureau of Labor Statistics CPI":     "BLS",
		"EIA.gov petroleum":                  "EIA",
		"Energy Information Administration":  "EIA",
		"DE.PULSE semantic evidence":         "—",
		"millisecond internal timing evidence": "—",
	}
	for source, want := range cases {
		if got := providerRightsSourceProvider(source); got != want {
			t.Fatalf("source %q resolved to %q; want %q", source, got, want)
		}
	}
}

func TestHOST002RegisteredExternalEvidenceCannotMasqueradeAsInternal(t *testing.T) {
	t.Setenv(runtimeModeEnv, "hosted")
	t.Setenv(hostedEnvironmentEnv, "prod")
	clearProviderRightsBundleForTest(t)
	now := time.Date(2026, 8, 26, 12, 0, 0, 0, time.UTC)

	for _, registration := range providerRegistrations() {
		record := EvidenceRecord{Source: registration.Name}
		if hostedRightsExternalEvidenceAllowed(record, now) {
			t.Fatalf("unreviewed registered provider %q crossed hosted persistence as internal evidence", registration.Name)
		}
	}
}

func TestHOST003ReviewedNonStandardProviderEvidenceCanReachCanonicalPersistence(t *testing.T) {
	providers := []string{"TradeInsight", "SEC", "SEC EDGAR", "BLS", "EIA"}
	records := make([]ProviderDataRightsMetadata, 0, len(providers))
	for _, provider := range providers {
		records = append(records, approvedHostedRightsFixtureFor(provider))
	}
	bindHostedRightsBundleForTest(t, records...)
	now := time.Date(2026, 8, 26, 12, 0, 0, 0, time.UTC)

	for _, provider := range providers {
		record := EvidenceRecord{Source: provider}
		if !hostedRightsExternalEvidenceAllowed(record, now) {
			t.Fatalf("reviewed provider %q could not reach canonical hosted persistence", provider)
		}
	}
}
