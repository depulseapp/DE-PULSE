package main

import (
	"strings"
	"time"
)

const (
	providerDataRightsPolicyVersion          = "provider-data-rights-v19.0.0"
	providerCommercialReadinessPolicyVersion = "provider-commercial-readiness-v18.4.0"
	providerHostedRightsPolicyVersion        = "provider-hosted-rights-v19.0.0"
	providerRightsUnreviewed                 = "UNREVIEWED"
	providerRightsNotAsserted                = "NOT_ASSERTED"
	providerRightsApproved                   = "APPROVED"
	providerRightsDenied                     = "DENIED"
	providerCommercialBlocked                = "BLOCKED"
	providerCommercialReady                  = "READY"
	providerHostedRightsBlocked              = "BLOCKED"
	providerHostedRightsAllowed              = "ALLOWED"
	providerRightsRenewalUnknown              = "UNKNOWN"
	providerRightsRenewalCurrent              = "CURRENT"
	providerRightsRenewalExpired              = "EXPIRED"
)

const (
	providerHostedUseCommercialMultiUser = "COMMERCIAL_MULTI_USER"
	providerHostedUseProxy               = "PROXY"
	providerHostedUseCacheRetention      = "CACHE_RETENTION"
	providerHostedUseRedisplay           = "REDISTRIBUTION_DISPLAY"
	providerHostedUseAI                  = "AI_DERIVED_USE"
	providerHostedUseLiveFanout          = "LIVE_FANOUT"
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
	ReviewState         string                      `json:"reviewState"`
	CommercialUse       string                      `json:"commercialUse"`
	MultiUserUse        string                      `json:"multiUserUse"`
	Proxying             string                      `json:"proxying"`
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

// ProviderHostedRightsDecision is the single fail-closed decision consumed by
// later hosted serving/cache/persistence/live-fanout composition. It is not a
// provider rank, score or fallback decision and must not become a second router.
type ProviderHostedRightsDecision struct {
	PolicyVersion string   `json:"policyVersion"`
	State         string   `json:"state"`
	Allowed       bool     `json:"allowed"`
	Purpose       string   `json:"purpose"`
	Environment   string   `json:"environment"`
	EvidenceRef   string   `json:"evidenceRef,omitempty"`
	EvidenceDigest string  `json:"evidenceDigest,omitempty"`
	BlockingReasons []string `json:"blockingReasons,omitempty"`
}

func providerRightExplicitlyApproved(rights ProviderDataRightsMetadata, value string) bool {
	return rights.EvidenceBound &&
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
	blocking := make([]string, 0, 10)

	if !rights.EvidenceBound || strings.TrimSpace(rights.EvidenceRef) == "" || strings.TrimSpace(rights.EvidenceDigest) == "" {
		blocking = append(blocking, "provider-specific rights evidence/provenance is not bound")
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
	default:
		blocking = append(blocking, "unknown hosted rights purpose")
	}

	state := providerHostedRightsBlocked
	if len(blocking) == 0 {
		state = providerHostedRightsAllowed
	}
	return ProviderHostedRightsDecision{
		PolicyVersion: providerHostedRightsPolicyVersion,
		State: state,
		Allowed: state == providerHostedRightsAllowed,
		Purpose: purpose,
		Environment: environment,
		EvidenceRef: strings.TrimSpace(rights.EvidenceRef),
		EvidenceDigest: strings.TrimSpace(rights.EvidenceDigest),
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

	blocking := make([]string, 0, 4)
	if !rights.EvidenceBound {
		blocking = append(blocking, "provider-specific rights evidence is not bound")
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

// providerDataRightsMetadata fails closed. Provider-specific approvals require
// separately bound, reviewable evidence and are never inferred from a configured
// credential, successful request, public endpoint, or inherited production route.
func providerDataRightsMetadata(provider string) ProviderDataRightsMetadata {
	rights := ProviderDataRightsMetadata{
		PolicyVersion:    providerDataRightsPolicyVersion,
		ReviewState:      providerRightsUnreviewed,
		CommercialUse:    providerRightsNotAsserted,
		MultiUserUse:     providerRightsNotAsserted,
		Proxying:          providerRightsNotAsserted,
		CachingRetention: providerRightsNotAsserted,
		Redistribution:   providerRightsNotAsserted,
		Display:           providerRightsNotAsserted,
		AIUse:             providerRightsNotAsserted,
		UsageLimits:       "",
		Attribution:       "",
		RenewalState:      providerRightsRenewalUnknown,
		EvidenceBound:     false,
		Detail:            "Operational entitlement is separate from licensing/data-rights approval; provider-specific commercial, multi-user, proxy, cache/retention, redistribution/display, AI/derived-use, environment and term evidence is not bound.",
	}
	rights.CommercialReadiness = evaluateProviderCommercialReadiness(rights)
	return rights
}
