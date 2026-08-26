package main

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"
)

const (
	providerDataRightsPolicyVersion          = "provider-data-rights-v19.0.0"
	providerCommercialReadinessPolicyVersion = "provider-commercial-readiness-v18.4.0"
	providerHostedRightsPolicyVersion        = "provider-hosted-rights-v19.0.0"
	providerRightsBundlePolicyVersion        = "provider-rights-bundle-v19.0.0"
	providerRightsBundlePathEnv              = "DEPULSE_PROVIDER_RIGHTS_BUNDLE_PATH"
	providerRightsBundleSHA256Env            = "DEPULSE_PROVIDER_RIGHTS_BUNDLE_SHA256"
	providerRightsUnreviewed                 = "UNREVIEWED"
	providerRightsNotAsserted                = "NOT_ASSERTED"
	providerRightsApproved                   = "APPROVED"
	providerRightsDenied                     = "DENIED"
	providerCommercialBlocked                = "BLOCKED"
	providerCommercialReady                  = "READY"
	providerHostedRightsBlocked              = "BLOCKED"
	providerHostedRightsAllowed              = "ALLOWED"
	providerRightsRenewalUnknown             = "UNKNOWN"
	providerRightsRenewalCurrent             = "CURRENT"
	providerRightsRenewalExpired             = "EXPIRED"
)

const (
	providerHostedUseCommercialMultiUser = "COMMERCIAL_MULTI_USER"
	providerHostedUseProxy               = "PROXY"
	providerHostedUseCacheRetention      = "CACHE_RETENTION"
	providerHostedUseRedisplay           = "REDISTRIBUTION_DISPLAY"
	providerHostedUseAI                  = "AI_DERIVED_USE"
	providerHostedUseLiveFanout          = "LIVE_FANOUT"
	// PRODUCTION_SERVING is the conservative hosted admission boundary used by
	// executable provider routing. It requires every right needed to acquire,
	// proxy, retain and redisplay shared hosted market data. AI use remains a
	// separate purpose because not every provider payload is sent to AI.
	providerHostedUseProductionServing = "PRODUCTION_SERVING"
)

// ProviderCommercialReadiness is the v18 fail-closed release-suitability view
// over provider-specific rights evidence. It remains separate from operational
// entitlement and from Smart Router scoring/routing.
type ProviderCommercialReadiness struct {
	PolicyVersion       string   `json:"policyVersion"`
	State               string   `json:"state"`
	CommercialUseReady  bool     `json:"commercialUseReady"`
	RedistributionReady bool     `json:"redistributionReady"`
	AIUseReady          bool     `json:"aiUseReady"`
	RequiresReview      bool     `json:"requiresReview"`
	BlockingReasons     []string `json:"blockingReasons,omitempty"`
}

// ProviderDataRightsMetadata is canonical provider legal/data-rights governance
// truth. Operational API entitlement is tracked separately by Smart Router v2.
// A working key or public endpoint never grants any field here by inference.
type ProviderDataRightsMetadata struct {
	PolicyVersion       string                      `json:"policyVersion"`
	Provider            string                      `json:"provider"`
	ReviewState         string                      `json:"reviewState"`
	CommercialUse       string                      `json:"commercialUse"`
	MultiUserUse        string                      `json:"multiUserUse"`
	Proxying            string                      `json:"proxying"`
	CachingRetention    string                      `json:"cachingRetention"`
	Redistribution      string                      `json:"redistribution"`
	Display             string                      `json:"display"`
	AIUse               string                      `json:"aiUse"`
	UsageLimits         string                      `json:"usageLimits"`
	Attribution         string                      `json:"attribution"`
	AllowedEnvironments []string                    `json:"allowedEnvironments,omitempty"`
	EffectiveAt         string                      `json:"effectiveAt,omitempty"`
	ExpiresAt           string                      `json:"expiresAt,omitempty"`
	RenewalState        string                      `json:"renewalState"`
	EvidenceBound       bool                        `json:"evidenceBound"`
	EvidenceRef         string                      `json:"evidenceRef,omitempty"`
	EvidenceDigest      string                      `json:"evidenceDigest,omitempty"`
	ReviewedAt          string                      `json:"reviewedAt,omitempty"`
	Detail              string                      `json:"detail"`
	CommercialReadiness ProviderCommercialReadiness `json:"commercialReadiness"`
}

type ProviderDataRightsBundle struct {
	PolicyVersion string                       `json:"policyVersion"`
	Records       []ProviderDataRightsMetadata `json:"records"`
}

// ProviderHostedRightsDecision is the single fail-closed decision consumed by
// hosted routing/serving/cache/persistence/live-fanout composition. It is not a
// provider rank, score or fallback decision and must not become a second router.
type ProviderHostedRightsDecision struct {
	PolicyVersion   string   `json:"policyVersion"`
	Provider        string   `json:"provider"`
	State           string   `json:"state"`
	Allowed         bool     `json:"allowed"`
	Purpose         string   `json:"purpose"`
	Environment     string   `json:"environment"`
	EvidenceRef     string   `json:"evidenceRef,omitempty"`
	EvidenceDigest  string   `json:"evidenceDigest,omitempty"`
	BlockingReasons []string `json:"blockingReasons,omitempty"`
}

func normalizeProviderRightsDigest(value string) (string, bool) {
	value = strings.ToLower(strings.TrimSpace(value))
	value = strings.TrimPrefix(value, "sha256:")
	if len(value) != sha256.Size*2 {
		return "", false
	}
	if _, err := hex.DecodeString(value); err != nil {
		return "", false
	}
	return value, true
}

func providerRightExplicitlyApproved(rights ProviderDataRightsMetadata, value string) bool {
	_, digestOK := normalizeProviderRightsDigest(rights.EvidenceDigest)
	return rights.EvidenceBound &&
		digestOK &&
		strings.TrimSpace(rights.Provider) != "" &&
		strings.EqualFold(strings.TrimSpace(rights.ReviewState), providerRightsApproved) &&
		strings.EqualFold(strings.TrimSpace(value), providerRightsApproved)
}

func parseProviderRightsTime(value string) (time.Time, bool) {
	value = strings.TrimSpace(value)
	if value == "" {
		return time.Time{}, false
	}
	t, err := time.Parse(time.RFC3339, value)
	return t, err == nil
}

func providerRightsEnvironmentAllowed(rights ProviderDataRightsMetadata, environment string) bool {
	environment = strings.ToLower(strings.TrimSpace(environment))
	if environment == "" || len(rights.AllowedEnvironments) == 0 {
		return false
	}
	for _, candidate := range rights.AllowedEnvironments {
		if strings.EqualFold(strings.TrimSpace(candidate), environment) {
			return true
		}
	}
	return false
}

// evaluateProviderHostedRightsDecision enforces HOST-001..003 independently of
// operational provider availability. The caller supplies now so expiry/downgrade
// behavior is deterministic and testable. Unknown purpose, missing provenance,
// malformed/effective-future/expired evidence, an unapproved environment, or an
// unapproved required right always blocks.
func evaluateProviderHostedRightsDecision(rights ProviderDataRightsMetadata, purpose, environment string, now time.Time) ProviderHostedRightsDecision {
	purpose = strings.ToUpper(strings.TrimSpace(purpose))
	environment = strings.ToLower(strings.TrimSpace(environment))
	blocking := make([]string, 0, 12)

	if strings.TrimSpace(rights.Provider) == "" {
		blocking = append(blocking, "provider identity is not bound to the rights record")
	}
	if !rights.EvidenceBound || strings.TrimSpace(rights.EvidenceRef) == "" {
		blocking = append(blocking, "provider-specific rights evidence/provenance is not bound")
	}
	if _, ok := normalizeProviderRightsDigest(rights.EvidenceDigest); !ok {
		blocking = append(blocking, "provider rights evidence digest is not a valid SHA-256 binding")
	}
	if !strings.EqualFold(strings.TrimSpace(rights.ReviewState), providerRightsApproved) {
		blocking = append(blocking, "provider rights review is not approved")
	}
	if strings.EqualFold(strings.TrimSpace(rights.RenewalState), providerRightsRenewalExpired) {
		blocking = append(blocking, "provider rights renewal state is expired")
	}

	effectiveAt, effectiveOK := parseProviderRightsTime(rights.EffectiveAt)
	if !effectiveOK {
		blocking = append(blocking, "provider rights effective time is missing or invalid")
	} else if now.Before(effectiveAt) {
		blocking = append(blocking, "provider rights are not yet effective")
	}
	expiresAt, expiryOK := parseProviderRightsTime(rights.ExpiresAt)
	if !expiryOK {
		blocking = append(blocking, "provider rights expiry time is missing or invalid")
	} else if !now.Before(expiresAt) {
		blocking = append(blocking, "provider rights evidence is expired")
	}
	if !providerRightsEnvironmentAllowed(rights, environment) {
		blocking = append(blocking, "provider rights do not approve this environment")
	}
	if strings.TrimSpace(rights.UsageLimits) == "" {
		blocking = append(blocking, "provider usage limits are not recorded")
	}
	if strings.TrimSpace(rights.Attribution) == "" {
		blocking = append(blocking, "provider attribution requirements are not recorded")
	}

	require := func(value, reason string) {
		if !providerRightExplicitlyApproved(rights, value) {
			blocking = append(blocking, reason)
		}
	}
	switch purpose {
	case providerHostedUseCommercialMultiUser:
		require(rights.CommercialUse, "commercial use is not approved")
		require(rights.MultiUserUse, "multi-user use is not approved")
	case providerHostedUseProxy:
		require(rights.CommercialUse, "commercial use is not approved")
		require(rights.Proxying, "proxying is not approved")
	case providerHostedUseCacheRetention:
		require(rights.CommercialUse, "commercial use is not approved")
		require(rights.CachingRetention, "caching/retention is not approved")
	case providerHostedUseRedisplay:
		require(rights.CommercialUse, "commercial use is not approved")
		require(rights.Redistribution, "redistribution is not approved")
		require(rights.Display, "display is not approved")
	case providerHostedUseAI:
		require(rights.CommercialUse, "commercial use is not approved")
		require(rights.AIUse, "AI/derived use is not approved")
	case providerHostedUseLiveFanout:
		require(rights.CommercialUse, "commercial use is not approved")
		require(rights.MultiUserUse, "multi-user use is not approved")
		require(rights.Proxying, "proxying is not approved")
		require(rights.Redistribution, "redistribution is not approved")
		require(rights.Display, "display is not approved")
	case providerHostedUseProductionServing:
		require(rights.CommercialUse, "commercial use is not approved")
		require(rights.MultiUserUse, "multi-user use is not approved")
		require(rights.Proxying, "proxying is not approved")
		require(rights.CachingRetention, "caching/retention is not approved")
		require(rights.Redistribution, "redistribution is not approved")
		require(rights.Display, "display is not approved")
	default:
		blocking = append(blocking, "unknown hosted rights purpose")
	}

	state := providerHostedRightsBlocked
	if len(blocking) == 0 {
		state = providerHostedRightsAllowed
	}
	return ProviderHostedRightsDecision{
		PolicyVersion:   providerHostedRightsPolicyVersion,
		Provider:        strings.TrimSpace(rights.Provider),
		State:           state,
		Allowed:         state == providerHostedRightsAllowed,
		Purpose:         purpose,
		Environment:     environment,
		EvidenceRef:     strings.TrimSpace(rights.EvidenceRef),
		EvidenceDigest:  strings.TrimSpace(rights.EvidenceDigest),
		BlockingReasons: blocking,
	}
}

// evaluateProviderCommercialReadiness preserves the v18 release-suitability
// contract. Hosted serving eligibility is deliberately stricter and is owned by
// evaluateProviderHostedRightsDecision.
func evaluateProviderCommercialReadiness(rights ProviderDataRightsMetadata) ProviderCommercialReadiness {
	commercialReady := providerRightExplicitlyApproved(rights, rights.CommercialUse)
	redistributionReady := providerRightExplicitlyApproved(rights, rights.Redistribution)
	aiReady := providerRightExplicitlyApproved(rights, rights.AIUse)

	blocking := make([]string, 0, 5)
	if strings.TrimSpace(rights.Provider) == "" {
		blocking = append(blocking, "provider identity is not bound")
	}
	if !rights.EvidenceBound {
		blocking = append(blocking, "provider-specific rights evidence is not bound")
	}
	if _, ok := normalizeProviderRightsDigest(rights.EvidenceDigest); !ok {
		blocking = append(blocking, "provider-specific evidence digest is invalid")
	}
	if !commercialReady {
		blocking = append(blocking, "commercial-use approval is not bound")
	}
	if !redistributionReady {
		blocking = append(blocking, "redistribution approval is not bound")
	}
	if !aiReady {
		blocking = append(blocking, "AI-use approval is not bound")
	}

	state := providerCommercialBlocked
	if commercialReady && redistributionReady && aiReady {
		state = providerCommercialReady
		blocking = nil
	}
	return ProviderCommercialReadiness{
		PolicyVersion:       providerCommercialReadinessPolicyVersion,
		State:               state,
		CommercialUseReady:  commercialReady,
		RedistributionReady: redistributionReady,
		AIUseReady:          aiReady,
		RequiresReview:      state != providerCommercialReady,
		BlockingReasons:     blocking,
	}
}

func defaultProviderDataRightsMetadata(provider string) ProviderDataRightsMetadata {
	rights := ProviderDataRightsMetadata{
		PolicyVersion:    providerDataRightsPolicyVersion,
		Provider:         strings.TrimSpace(provider),
		ReviewState:      providerRightsUnreviewed,
		CommercialUse:    providerRightsNotAsserted,
		MultiUserUse:     providerRightsNotAsserted,
		Proxying:         providerRightsNotAsserted,
		CachingRetention: providerRightsNotAsserted,
		Redistribution:   providerRightsNotAsserted,
		Display:          providerRightsNotAsserted,
		AIUse:            providerRightsNotAsserted,
		UsageLimits:      "",
		Attribution:      "",
		RenewalState:     providerRightsRenewalUnknown,
		EvidenceBound:    false,
		Detail:           "Operational entitlement is separate from licensing/data-rights approval; provider-specific commercial, multi-user, proxy, cache/retention, redistribution/display, AI/derived-use, environment and term evidence is not bound.",
	}
	rights.CommercialReadiness = evaluateProviderCommercialReadiness(rights)
	return rights
}

func loadProviderDataRightsBundle() (ProviderDataRightsBundle, error) {
	path := strings.TrimSpace(os.Getenv(providerRightsBundlePathEnv))
	if path == "" {
		return ProviderDataRightsBundle{}, fmt.Errorf("provider rights bundle path is not configured")
	}
	expected, ok := normalizeProviderRightsDigest(os.Getenv(providerRightsBundleSHA256Env))
	if !ok {
		return ProviderDataRightsBundle{}, fmt.Errorf("provider rights bundle SHA-256 pin is missing or invalid")
	}
	raw, err := os.ReadFile(path)
	if err != nil {
		return ProviderDataRightsBundle{}, fmt.Errorf("provider rights bundle cannot be read: %w", err)
	}
	sum := sha256.Sum256(raw)
	if hex.EncodeToString(sum[:]) != expected {
		return ProviderDataRightsBundle{}, fmt.Errorf("provider rights bundle SHA-256 pin mismatch")
	}
	var bundle ProviderDataRightsBundle
	if err := json.Unmarshal(raw, &bundle); err != nil {
		return ProviderDataRightsBundle{}, fmt.Errorf("provider rights bundle is invalid JSON: %w", err)
	}
	if bundle.PolicyVersion != providerRightsBundlePolicyVersion {
		return ProviderDataRightsBundle{}, fmt.Errorf("provider rights bundle policy version %q is not supported", bundle.PolicyVersion)
	}
	return bundle, nil
}

func providerDataRightsMetadata(provider string) ProviderDataRightsMetadata {
	fallback := defaultProviderDataRightsMetadata(provider)
	bundle, err := loadProviderDataRightsBundle()
	if err != nil {
		if isHostedRuntime() {
			fallback.Detail = "Hosted provider rights are blocked: " + err.Error()
		}
		return fallback
	}
	for _, candidate := range bundle.Records {
		if !strings.EqualFold(strings.TrimSpace(candidate.Provider), strings.TrimSpace(provider)) {
			continue
		}
		if candidate.PolicyVersion != providerDataRightsPolicyVersion {
			fallback.Detail = "Hosted provider rights are blocked: provider rights record policy version is invalid."
			return fallback
		}
		if strings.TrimSpace(candidate.Provider) == "" {
			fallback.Detail = "Hosted provider rights are blocked: provider identity is missing from the bound record."
			return fallback
		}
		candidate.Provider = strings.TrimSpace(candidate.Provider)
		candidate.CommercialReadiness = evaluateProviderCommercialReadiness(candidate)
		return candidate
	}
	if isHostedRuntime() {
		fallback.Detail = "Hosted provider rights are blocked: no provider-specific record exists in the SHA-pinned rights bundle."
	}
	return fallback
}

func hostedProviderRightsDecision(provider, purpose string, now time.Time) ProviderHostedRightsDecision {
	if !isHostedRuntime() {
		return ProviderHostedRightsDecision{
			PolicyVersion: providerHostedRightsPolicyVersion,
			Provider:      strings.TrimSpace(provider),
			State:         providerHostedRightsAllowed,
			Allowed:       true,
			Purpose:       strings.ToUpper(strings.TrimSpace(purpose)),
			Environment:   "desktop",
		}
	}
	environment := strings.ToLower(strings.TrimSpace(os.Getenv(hostedEnvironmentEnv)))
	return evaluateProviderHostedRightsDecision(providerDataRightsMetadata(provider), purpose, environment, now)
}

func hostedProviderRightsAllowed(provider, purpose string, now time.Time) bool {
	return hostedProviderRightsDecision(provider, purpose, now).Allowed
}

func providerQuoteHostedRightsAllowed(q Quote, purpose string, now time.Time) bool {
	if !isHostedRuntime() {
		return true
	}
	provider := sourceProvider(q.Source)
	if provider == "" || provider == "—" {
		return false
	}
	return hostedProviderRightsAllowed(provider, purpose, now)
}

// hostedRightsFilteredMarketCache extends the existing market-cache owner. In
// hosted mode only provider-attributed quotes with current production-serving
// and retention approval can cross the cache boundary. Legacy mixed collections
// without per-row provider provenance fail closed; desktop remains unchanged.
func hostedRightsFilteredMarketCache(c MarketCache) MarketCache {
	if !isHostedRuntime() {
		return c
	}
	out := MarketCache{
		Quotes:                   map[string]Quote{},
		ProviderCapabilityStates: c.ProviderCapabilityStates,
		SavedAt:                  c.SavedAt,
	}
	now := time.Now()
	for symbol, q := range c.Quotes {
		if providerQuoteHostedRightsAllowed(q, providerHostedUseProductionServing, now) {
			out.Quotes[symbol] = q
		}
	}
	return out
}

func (c MarketCache) MarshalJSON() ([]byte, error) {
	type marketCacheAlias MarketCache
	filtered := hostedRightsFilteredMarketCache(c)
	return json.Marshal(marketCacheAlias(filtered))
}

func (c *MarketCache) UnmarshalJSON(data []byte) error {
	type marketCacheAlias MarketCache
	var decoded marketCacheAlias
	if err := json.Unmarshal(data, &decoded); err != nil {
		return err
	}
	filtered := hostedRightsFilteredMarketCache(MarketCache(decoded))
	*c = filtered
	return nil
}
