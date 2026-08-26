package main

import (
	"errors"
	"sort"
	"strings"
	"time"
)

type ProductPlan string

type ProductStatus string

const (
	ProductPlanFoundationCore     ProductPlan = "FOUNDATION_CORE"
	ProductPlanFoundationExtended ProductPlan = "FOUNDATION_EXTENDED"

	ProductStatusActive    ProductStatus = "ACTIVE"
	ProductStatusGrace     ProductStatus = "GRACE"
	ProductStatusSuspended ProductStatus = "SUSPENDED"
	ProductStatusDisabled  ProductStatus = "DISABLED"

	productCapabilityHostedServing       = "HOSTED_SERVING"
	productCapabilityExtendedIntelligence = "EXTENDED_INTELLIGENCE"
	productQuotaHostedMutationUnits      = "HOSTED_MUTATION_UNITS"
	productQuotaHostedExpensiveUnits     = "HOSTED_EXPENSIVE_UNITS"
	productMeteringWindow                 = 24 * time.Hour
)

const (
	foundationCoreMutationUnitsPerWindow      int64 = 10000
	foundationCoreExpensiveUnitsPerWindow     int64 = 1000
	foundationExtendedMutationUnitsPerWindow  int64 = 50000
	foundationExtendedExpensiveUnitsPerWindow int64 = 5000
)

type TenantProductEntitlement struct {
	TenantID        string            `json:"tenantId"`
	Plan            ProductPlan       `json:"plan"`
	Status          ProductStatus     `json:"status"`
	WindowStartedAt int64             `json:"windowStartedAt"`
	Usage           map[string]int64  `json:"usage,omitempty"`
	StatusChangedAt int64             `json:"statusChangedAt,omitempty"`
	UpdatedAt       int64             `json:"updatedAt"`
}

type ProductQuotaCharge struct {
	Dimension string `json:"dimension"`
	Units     int64  `json:"units"`
}

type ProductQuotaDecision struct {
	Dimension string `json:"dimension"`
	Limit     int64  `json:"limit"`
	Used      int64  `json:"used"`
	Remaining int64  `json:"remaining"`
}

type ProductEntitlementDecision struct {
	Allowed           bool                   `json:"allowed"`
	TenantID          string                 `json:"tenantId"`
	Plan              ProductPlan            `json:"plan,omitempty"`
	Status            ProductStatus          `json:"status,omitempty"`
	Capability        string                 `json:"capability,omitempty"`
	Quotas            []ProductQuotaDecision `json:"quotas,omitempty"`
	BlockingReasons   []string               `json:"blockingReasons,omitempty"`
	RetryAfterSeconds int                    `json:"retryAfterSeconds,omitempty"`
}

type ProductEntitlementSnapshot struct {
	TenantID        string                 `json:"tenantId"`
	Plan            ProductPlan            `json:"plan"`
	Status          ProductStatus          `json:"status"`
	Capabilities    []string               `json:"capabilities"`
	Quotas          []ProductQuotaDecision `json:"quotas"`
	WindowStartedAt int64                  `json:"windowStartedAt"`
	WindowEndsAt    int64                  `json:"windowEndsAt"`
	UpdatedAt       int64                  `json:"updatedAt"`
}

type productPlanPolicy struct {
	capabilities map[string]struct{}
	quotas       map[string]int64
}

func validProductPlan(plan ProductPlan) bool {
	return plan == ProductPlanFoundationCore || plan == ProductPlanFoundationExtended
}

func validProductStatus(status ProductStatus) bool {
	switch status {
	case ProductStatusActive, ProductStatusGrace, ProductStatusSuspended, ProductStatusDisabled:
		return true
	default:
		return false
	}
}

func productPlanPolicyFor(plan ProductPlan) (productPlanPolicy, bool) {
	standard := strings.ToUpper(productCapabilityHostedServing)
	extended := strings.ToUpper(productCapabilityExtendedIntelligence)
	mutation := strings.ToUpper(productQuotaHostedMutationUnits)
	expensive := strings.ToUpper(productQuotaHostedExpensiveUnits)
	switch plan {
	case ProductPlanFoundationCore:
		return productPlanPolicy{
			capabilities: map[string]struct{}{standard: {}},
			quotas: map[string]int64{
				mutation:  foundationCoreMutationUnitsPerWindow,
				expensive: foundationCoreExpensiveUnitsPerWindow,
			},
		}, true
	case ProductPlanFoundationExtended:
		return productPlanPolicy{
			capabilities: map[string]struct{}{standard: {}, extended: {}},
			quotas: map[string]int64{
				mutation:  foundationExtendedMutationUnitsPerWindow,
				expensive: foundationExtendedExpensiveUnitsPerWindow,
			},
		}, true
	default:
		return productPlanPolicy{}, false
	}
}

func cloneTenantProductEntitlements(in []TenantProductEntitlement) []TenantProductEntitlement {
	out := make([]TenantProductEntitlement, len(in))
	for i := range in {
		out[i] = in[i]
		if in[i].Usage != nil {
			out[i].Usage = make(map[string]int64, len(in[i].Usage))
			for key, value := range in[i].Usage {
				out[i].Usage[key] = value
			}
		}
	}
	return out
}

func (s *IdentityService) ensureHostedProductEntitlements() error {
	if s == nil {
		return errors.New("identity unavailable")
	}
	s.mu.Lock()
	missing := make([]string, 0, len(s.state.Tenants))
	for _, tenant := range s.state.Tenants {
		tenantID := strings.TrimSpace(tenant.ID)
		if tenantID == "" {
			continue
		}
		found := false
		for _, product := range s.state.ProductEntitlements {
			if strings.TrimSpace(product.TenantID) == tenantID {
				found = true
				break
			}
		}
		if !found {
			missing = append(missing, tenantID)
		}
	}
	s.mu.Unlock()
	for _, tenantID := range missing {
		if err := s.setTenantProductPolicy(tenantID, ProductPlanFoundationCore, ProductStatusActive); err != nil {
			return err
		}
	}
	return nil
}

func (s *IdentityService) setTenantProductPolicy(tenantID string, plan ProductPlan, status ProductStatus) error {
	if s == nil {
		return errors.New("identity unavailable")
	}
	tenantID = strings.TrimSpace(tenantID)
	if tenantID == "" || !validProductPlan(plan) || !validProductStatus(status) {
		return errors.New("invalid tenant product policy")
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	tenantExists := false
	for _, tenant := range s.state.Tenants {
		if strings.TrimSpace(tenant.ID) == tenantID {
			tenantExists = true
			break
		}
	}
	if !tenantExists {
		return errors.New("tenant unavailable for product policy")
	}
	previous := cloneTenantProductEntitlements(s.state.ProductEntitlements)
	now := s.now().UnixMilli()
	index := -1
	matches := 0
	for i := range s.state.ProductEntitlements {
		if strings.TrimSpace(s.state.ProductEntitlements[i].TenantID) == tenantID {
			index = i
			matches++
		}
	}
	if matches > 1 {
		return errors.New("ambiguous tenant product policy")
	}
	if index < 0 {
		s.state.ProductEntitlements = append(s.state.ProductEntitlements, TenantProductEntitlement{
			TenantID: tenantID, Plan: plan, Status: status, WindowStartedAt: now,
			Usage: map[string]int64{}, StatusChangedAt: now, UpdatedAt: now,
		})
	} else {
		current := &s.state.ProductEntitlements[index]
		if current.WindowStartedAt <= 0 {
			current.WindowStartedAt = now
		}
		if current.Usage == nil {
			current.Usage = map[string]int64{}
		}
		if current.Status != status {
			current.StatusChangedAt = now
		}
		current.Plan = plan
		current.Status = status
		current.UpdatedAt = now
	}
	if err := s.persistLocked(); err != nil {
		s.state.ProductEntitlements = previous
		return err
	}
	return nil
}

func productStatusAllowsServing(status ProductStatus) bool {
	return status == ProductStatusActive || status == ProductStatusGrace
}

func normalizeProductCapability(value string) string {
	return strings.ToUpper(strings.TrimSpace(value))
}

func normalizeProductQuotaDimension(value string) string {
	return strings.ToUpper(strings.TrimSpace(value))
}

func (s *IdentityService) tenantProductStateLocked(tenantID string) (TenantProductEntitlement, int, bool) {
	tenantID = strings.TrimSpace(tenantID)
	index := -1
	matches := 0
	for i := range s.state.ProductEntitlements {
		if strings.TrimSpace(s.state.ProductEntitlements[i].TenantID) == tenantID {
			index = i
			matches++
		}
	}
	if matches != 1 || index < 0 {
		return TenantProductEntitlement{}, -1, false
	}
	return s.state.ProductEntitlements[index], index, true
}

func effectiveProductMetering(state TenantProductEntitlement, now int64) (int64, map[string]int64) {
	windowMillis := int64(productMeteringWindow / time.Millisecond)
	if state.WindowStartedAt <= 0 || now < state.WindowStartedAt || now-state.WindowStartedAt >= windowMillis {
		return now, map[string]int64{}
	}
	usage := make(map[string]int64, len(state.Usage))
	for key, value := range state.Usage {
		if value > 0 {
			usage[normalizeProductQuotaDimension(key)] += value
		}
	}
	return state.WindowStartedAt, usage
}

func normalizedProductCharges(charges []ProductQuotaCharge) (map[string]int64, error) {
	out := make(map[string]int64, len(charges))
	for _, charge := range charges {
		dimension := normalizeProductQuotaDimension(charge.Dimension)
		if dimension == "" || charge.Units <= 0 {
			return nil, errors.New("invalid product quota charge")
		}
		if out[dimension] > 0 && charge.Units > int64(^uint64(0)>>1)-out[dimension] {
			return nil, errors.New("product quota charge overflow")
		}
		out[dimension] += charge.Units
	}
	return out, nil
}

func evaluateProductDecision(state TenantProductEntitlement, capability string, charges map[string]int64, now int64) ProductEntitlementDecision {
	capability = normalizeProductCapability(capability)
	decision := ProductEntitlementDecision{
		TenantID: strings.TrimSpace(state.TenantID), Plan: state.Plan, Status: state.Status, Capability: capability,
	}
	blocking := make([]string, 0, 4)
	policy, planOK := productPlanPolicyFor(state.Plan)
	if !planOK {
		blocking = append(blocking, "unknown product plan")
	}
	if !validProductStatus(state.Status) {
		blocking = append(blocking, "unknown product status")
	} else if !productStatusAllowsServing(state.Status) {
		blocking = append(blocking, "product status blocks serving")
	}
	if capability == "" {
		blocking = append(blocking, "product capability required")
	} else if planOK {
		if _, ok := policy.capabilities[capability]; !ok {
			blocking = append(blocking, "product capability not entitled")
		}
	}
	windowStartedAt, usage := effectiveProductMetering(state, now)
	if planOK {
		dimensions := make([]string, 0, len(policy.quotas))
		for dimension := range policy.quotas {
			dimensions = append(dimensions, dimension)
		}
		sort.Strings(dimensions)
		for _, dimension := range dimensions {
			limit := policy.quotas[dimension]
			used := usage[dimension]
			remaining := limit - used
			if remaining < 0 {
				remaining = 0
			}
			decision.Quotas = append(decision.Quotas, ProductQuotaDecision{Dimension: dimension, Limit: limit, Used: used, Remaining: remaining})
		}
		for dimension, units := range charges {
			limit, ok := policy.quotas[dimension]
			if !ok {
				blocking = append(blocking, "product quota dimension not entitled")
				continue
			}
			used := usage[dimension]
			if units > limit || used > limit-units {
				blocking = append(blocking, "product quota exhausted")
				windowEndsAt := windowStartedAt + int64(productMeteringWindow/time.Millisecond)
				retry := int((windowEndsAt - now + int64(time.Second/time.Millisecond) - 1) / int64(time.Second/time.Millisecond))
				if retry < 1 {
					retry = 1
				}
				if decision.RetryAfterSeconds == 0 || retry < decision.RetryAfterSeconds {
					decision.RetryAfterSeconds = retry
				}
			}
		}
	}
	decision.BlockingReasons = blocking
	decision.Allowed = len(blocking) == 0
	return decision
}

func (s *IdentityService) authorizeHostedProduct(principal Principal, capability string) ProductEntitlementDecision {
	decision := ProductEntitlementDecision{TenantID: normalizedTenantID(principal.TenantID), Capability: normalizeProductCapability(capability)}
	if s == nil || strings.TrimSpace(principal.UserID) == "" || strings.TrimSpace(principal.SessionID) == "" {
		decision.BlockingReasons = []string{"verified authentication context required"}
		return decision
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	state, _, ok := s.tenantProductStateLocked(decision.TenantID)
	if !ok {
		decision.BlockingReasons = []string{"tenant product entitlement unavailable"}
		return decision
	}
	return evaluateProductDecision(state, decision.Capability, nil, s.now().UnixMilli())
}

func (s *IdentityService) consumeHostedProductQuota(principal Principal, capability string, charges []ProductQuotaCharge) (ProductEntitlementDecision, error) {
	decision := ProductEntitlementDecision{TenantID: normalizedTenantID(principal.TenantID), Capability: normalizeProductCapability(capability)}
	if s == nil || strings.TrimSpace(principal.UserID) == "" || strings.TrimSpace(principal.SessionID) == "" {
		decision.BlockingReasons = []string{"verified authentication context required"}
		return decision, nil
	}
	normalizedCharges, err := normalizedProductCharges(charges)
	if err != nil {
		return decision, err
	}
	if len(normalizedCharges) == 0 {
		return s.authorizeHostedProduct(principal, capability), nil
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	state, index, ok := s.tenantProductStateLocked(decision.TenantID)
	if !ok {
		decision.BlockingReasons = []string{"tenant product entitlement unavailable"}
		return decision, nil
	}
	now := s.now().UnixMilli()
	decision = evaluateProductDecision(state, capability, normalizedCharges, now)
	if !decision.Allowed {
		return decision, nil
	}
	previous := cloneTenantProductEntitlements(s.state.ProductEntitlements)
	windowStartedAt, usage := effectiveProductMetering(state, now)
	current := &s.state.ProductEntitlements[index]
	current.WindowStartedAt = windowStartedAt
	current.Usage = usage
	for dimension, units := range normalizedCharges {
		current.Usage[dimension] += units
	}
	current.UpdatedAt = now
	if err := s.persistLocked(); err != nil {
		s.state.ProductEntitlements = previous
		return ProductEntitlementDecision{}, err
	}
	return evaluateProductDecision(*current, capability, nil, now), nil
}

func (s *IdentityService) productEntitlementSnapshot(tenantID string) (ProductEntitlementSnapshot, error) {
	if s == nil {
		return ProductEntitlementSnapshot{}, errors.New("identity unavailable")
	}
	tenantID = strings.TrimSpace(tenantID)
	if tenantID == "" {
		return ProductEntitlementSnapshot{}, errors.New("tenant product entitlement unavailable")
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	state, _, ok := s.tenantProductStateLocked(tenantID)
	if !ok {
		return ProductEntitlementSnapshot{}, errors.New("tenant product entitlement unavailable")
	}
	policy, ok := productPlanPolicyFor(state.Plan)
	if !ok || !validProductStatus(state.Status) {
		return ProductEntitlementSnapshot{}, errors.New("tenant product entitlement invalid")
	}
	now := s.now().UnixMilli()
	windowStartedAt, usage := effectiveProductMetering(state, now)
	capabilities := make([]string, 0, len(policy.capabilities))
	for capability := range policy.capabilities {
		capabilities = append(capabilities, capability)
	}
	sort.Strings(capabilities)
	dimensions := make([]string, 0, len(policy.quotas))
	for dimension := range policy.quotas {
		dimensions = append(dimensions, dimension)
	}
	sort.Strings(dimensions)
	quotas := make([]ProductQuotaDecision, 0, len(dimensions))
	for _, dimension := range dimensions {
		limit := policy.quotas[dimension]
		used := usage[dimension]
		remaining := limit - used
		if remaining < 0 {
			remaining = 0
		}
		quotas = append(quotas, ProductQuotaDecision{Dimension: dimension, Limit: limit, Used: used, Remaining: remaining})
	}
	return ProductEntitlementSnapshot{
		TenantID: tenantID, Plan: state.Plan, Status: state.Status, Capabilities: capabilities, Quotas: quotas,
		WindowStartedAt: windowStartedAt,
		WindowEndsAt:    windowStartedAt + int64(productMeteringWindow/time.Millisecond),
		UpdatedAt:       state.UpdatedAt,
	}, nil
}
