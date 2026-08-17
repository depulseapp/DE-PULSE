package main

const (
	providerDataRightsPolicyVersion = "provider-data-rights-v18.4.0"
	providerRightsUnreviewed        = "UNREVIEWED"
	providerRightsNotAsserted       = "NOT_ASSERTED"
)

// ProviderDataRightsMetadata is governance truth only. It deliberately does not
// participate in provider eligibility, ranking, fallback, or deterministic
// Day/Swing/Long logic. Operational API entitlement is tracked separately by
// Smart Router v2.
type ProviderDataRightsMetadata struct {
	PolicyVersion string `json:"policyVersion"`
	ReviewState   string `json:"reviewState"`
	CommercialUse string `json:"commercialUse"`
	Redistribution string `json:"redistribution"`
	AIUse         string `json:"aiUse"`
	EvidenceBound bool   `json:"evidenceBound"`
	Detail        string `json:"detail"`
}

// providerDataRightsMetadata fails closed for commercialization readiness.
// A configured/working API key or a public endpoint never implies contractual
// permission for commercial use, redistribution, or AI/LLM use. Provider-
// specific approvals require separately bound evidence and are intentionally
// not inferred in this release slice.
func providerDataRightsMetadata(provider string) ProviderDataRightsMetadata {
	return ProviderDataRightsMetadata{
		PolicyVersion: providerDataRightsPolicyVersion,
		ReviewState:   providerRightsUnreviewed,
		CommercialUse: providerRightsNotAsserted,
		Redistribution: providerRightsNotAsserted,
		AIUse:         providerRightsNotAsserted,
		EvidenceBound: false,
		Detail:        "Operational entitlement is separate from licensing/data-rights approval; provider-specific rights evidence is not bound.",
	}
}
