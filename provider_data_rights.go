package main

import "strings"

const (
	providerDataRightsPolicyVersion          = "provider-data-rights-v18.4.0"
	providerCommercialReadinessPolicyVersion = "provider-commercial-readiness-v18.4.0"
	providerRightsUnreviewed                 = "UNREVIEWED"
	providerRightsNotAsserted                = "NOT_ASSERTED"
	providerRightsApproved                   = "APPROVED"
	providerCommercialBlocked                = "BLOCKED"
	providerCommercialReady                  = "READY"
)

// ProviderCommercialReadiness is a fail-closed release-suitability view over
// provider-specific rights evidence. It is governance truth only and must never
// participate in provider eligibility, ranking, fallback, or deterministic
// Day/Swing/Long logic.
type ProviderCommercialReadiness struct {
	PolicyVersion       string   `json:"policyVersion"`
	State               string   `json:"state"`
	CommercialUseReady  bool     `json:"commercialUseReady"`
	RedistributionReady bool     `json:"redistributionReady"`
	AIUseReady          bool     `json:"aiUseReady"`
	RequiresReview      bool     `json:"requiresReview"`
	BlockingReasons     []string `json:"blockingReasons,omitempty"`
}

// ProviderDataRightsMetadata is governance truth only. It deliberately does not
// participate in provider eligibility, ranking, fallback, or deterministic
// Day/Swing/Long logic. Operational API entitlement is tracked separately by
// Smart Router v2.
type ProviderDataRightsMetadata struct {
	PolicyVersion       string                      `json:"policyVersion"`
	ReviewState         string                      `json:"reviewState"`
	CommercialUse       string                      `json:"commercialUse"`
	Redistribution      string                      `json:"redistribution"`
	AIUse               string                      `json:"aiUse"`
	EvidenceBound       bool                        `json:"evidenceBound"`
	Detail              string                      `json:"detail"`
	CommercialReadiness ProviderCommercialReadiness `json:"commercialReadiness"`
}

func providerRightExplicitlyApproved(rights ProviderDataRightsMetadata, value string) bool {
	return rights.EvidenceBound &&
		strings.EqualFold(strings.TrimSpace(rights.ReviewState), providerRightsApproved) &&
		strings.EqualFold(strings.TrimSpace(value), providerRightsApproved)
}

// evaluateProviderCommercialReadiness treats anything other than explicit,
// evidence-bound approval as not release-ready. This is intentionally stricter
// than operational entitlement: a working API key never grants commercial,
// redistribution, or AI/LLM-use approval by inference.
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

// providerDataRightsMetadata fails closed for commercialization readiness.
// A configured/working API key or a public endpoint never implies contractual
// permission for commercial use, redistribution, or AI/LLM use. Provider-
// specific approvals require separately bound evidence and are intentionally
// not inferred in v18.4.
func providerDataRightsMetadata(provider string) ProviderDataRightsMetadata {
	rights := ProviderDataRightsMetadata{
		PolicyVersion:  providerDataRightsPolicyVersion,
		ReviewState:    providerRightsUnreviewed,
		CommercialUse:  providerRightsNotAsserted,
		Redistribution: providerRightsNotAsserted,
		AIUse:          providerRightsNotAsserted,
		EvidenceBound:  false,
		Detail:         "Operational entitlement is separate from licensing/data-rights approval; provider-specific rights evidence is not bound.",
	}
	rights.CommercialReadiness = evaluateProviderCommercialReadiness(rights)
	return rights
}
