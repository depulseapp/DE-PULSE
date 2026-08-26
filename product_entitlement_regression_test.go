package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func productQuotaByDimension(t *testing.T, snapshot ProductEntitlementSnapshot, dimension string) ProductQuotaDecision {
	t.Helper()
	dimension = normalizeProductQuotaDimension(dimension)
	for _, quota := range snapshot.Quotas {
		if quota.Dimension == dimension {
			return quota
		}
	}
	t.Fatalf("missing product quota dimension %s in %+v", dimension, snapshot.Quotas)
	return ProductQuotaDecision{}
}

func authenticatedProductOwner(t *testing.T, identity *IdentityService) (string, Principal) {
	t.Helper()
	_, bootstrap, err := identity.bootstrapOwnerSession()
	if err != nil {
		t.Fatal(err)
	}
	token, principal, err := identity.setPassword(bootstrap.UserID, "v19 product entitlement passphrase")
	if err != nil {
		t.Fatal(err)
	}
	return token, principal
}

func setProductUsageForTest(t *testing.T, identity *IdentityService, tenantID, dimension string, used int64) {
	t.Helper()
	identity.mu.Lock()
	defer identity.mu.Unlock()
	state, index, ok := identity.tenantProductStateLocked(tenantID)
	if !ok {
		t.Fatal("tenant product entitlement missing")
	}
	if state.Usage == nil {
		state.Usage = map[string]int64{}
	}
	state.WindowStartedAt = identity.now().UnixMilli()
	state.Usage[normalizeProductQuotaDimension(dimension)] = used
	state.UpdatedAt = identity.now().UnixMilli()
	identity.state.ProductEntitlements[index] = state
	if err := identity.persistLocked(); err != nil {
		t.Fatal(err)
	}
}

func TestHOST008ProductEntitlementIsIndependentFromRBACAndProviderRights(t *testing.T) {
	_, identity := newIdentityTestService(t)
	_, principal := authenticatedProductOwner(t, identity)
	if err := identity.setTenantProductPolicy(principal.TenantID, ProductPlanFoundationExtended, ProductStatusActive); err != nil {
		t.Fatal(err)
	}

	ownerDecision := identity.authorizeHostedProduct(principal, productCapabilityExtendedIntelligence)
	if !ownerDecision.Allowed {
		t.Fatalf("extended product capability unexpectedly denied: %+v", ownerDecision)
	}
	forgedRoleView := principal
	forgedRoleView.Role = RoleDemo
	demoRoleDecision := identity.authorizeHostedProduct(forgedRoleView, productCapabilityExtendedIntelligence)
	if !demoRoleDecision.Allowed {
		t.Fatalf("product entitlement incorrectly depended on RBAC role: %+v", demoRoleDecision)
	}
	if ownerDecision.Plan != demoRoleDecision.Plan || ownerDecision.Status != demoRoleDecision.Status {
		t.Fatalf("product policy changed with RBAC role: owner=%+v demo=%+v", ownerDecision, demoRoleDecision)
	}

	identityDecision := identity.authorizeHostedIdentity(principal, HostedIdentityRequirement{TenantID: principal.TenantID, Capability: hostedCapabilityStandardUse})
	if !identityDecision.Allowed {
		t.Fatalf("canonical RBAC decision unexpectedly failed: %+v", identityDecision)
	}
	if productCapabilityHostedServing == hostedCapabilityStandardUse {
		t.Fatal("product and RBAC capability namespaces collapsed")
	}
}

func TestHOST008ProductUpgradeDowngradeGraceAndSuspensionAreDeterministic(t *testing.T) {
	_, identity := newIdentityTestService(t)
	_, principal := authenticatedProductOwner(t, identity)
	if err := identity.setTenantProductPolicy(principal.TenantID, ProductPlanFoundationExtended, ProductStatusActive); err != nil {
		t.Fatal(err)
	}
	if decision := identity.authorizeHostedProduct(principal, productCapabilityExtendedIntelligence); !decision.Allowed {
		t.Fatalf("extended capability denied before downgrade: %+v", decision)
	}
	if _, err := identity.consumeHostedProductQuota(principal, productCapabilityHostedServing, []ProductQuotaCharge{{Dimension: productQuotaHostedMutationUnits, Units: 3}}); err != nil {
		t.Fatal(err)
	}
	before, err := identity.productEntitlementSnapshot(principal.TenantID)
	if err != nil {
		t.Fatal(err)
	}
	beforeUsage := productQuotaByDimension(t, before, productQuotaHostedMutationUnits).Used

	if err := identity.setTenantProductPolicy(principal.TenantID, ProductPlanFoundationCore, ProductStatusActive); err != nil {
		t.Fatal(err)
	}
	if decision := identity.authorizeHostedProduct(principal, productCapabilityExtendedIntelligence); decision.Allowed {
		t.Fatalf("downgrade retained extended capability: %+v", decision)
	}
	afterDowngrade, err := identity.productEntitlementSnapshot(principal.TenantID)
	if err != nil {
		t.Fatal(err)
	}
	if got := productQuotaByDimension(t, afterDowngrade, productQuotaHostedMutationUnits).Used; got != beforeUsage {
		t.Fatalf("downgrade reset metering usage: before=%d after=%d", beforeUsage, got)
	}

	if err := identity.setTenantProductPolicy(principal.TenantID, ProductPlanFoundationCore, ProductStatusGrace); err != nil {
		t.Fatal(err)
	}
	if decision := identity.authorizeHostedProduct(principal, productCapabilityHostedServing); !decision.Allowed || decision.Status != ProductStatusGrace {
		t.Fatalf("grace did not preserve current plan serving: %+v", decision)
	}
	if err := identity.setTenantProductPolicy(principal.TenantID, ProductPlanFoundationCore, ProductStatusSuspended); err != nil {
		t.Fatal(err)
	}
	if decision := identity.authorizeHostedProduct(principal, productCapabilityHostedServing); decision.Allowed {
		t.Fatalf("suspended product served: %+v", decision)
	}
	if err := identity.setTenantProductPolicy(principal.TenantID, ProductPlanFoundationCore, ProductStatusDisabled); err != nil {
		t.Fatal(err)
	}
	if decision := identity.authorizeHostedProduct(principal, productCapabilityHostedServing); decision.Allowed {
		t.Fatalf("disabled product served: %+v", decision)
	}

	if err := identity.setTenantProductPolicy(principal.TenantID, ProductPlanFoundationExtended, ProductStatusActive); err != nil {
		t.Fatal(err)
	}
	if decision := identity.authorizeHostedProduct(principal, productCapabilityExtendedIntelligence); !decision.Allowed {
		t.Fatalf("upgrade did not restore extended capability: %+v", decision)
	}
	afterUpgrade, err := identity.productEntitlementSnapshot(principal.TenantID)
	if err != nil {
		t.Fatal(err)
	}
	if got := productQuotaByDimension(t, afterUpgrade, productQuotaHostedMutationUnits).Used; got != beforeUsage {
		t.Fatalf("upgrade reset metering usage: before=%d after=%d", beforeUsage, got)
	}
}

func TestHOST009ProductQuotaExhaustionPersistsAndFailsClosed(t *testing.T) {
	persistence, identity := newIdentityTestService(t)
	_, principal := authenticatedProductOwner(t, identity)
	if err := identity.setTenantProductPolicy(principal.TenantID, ProductPlanFoundationCore, ProductStatusActive); err != nil {
		t.Fatal(err)
	}
	policy, ok := productPlanPolicyFor(ProductPlanFoundationCore)
	if !ok {
		t.Fatal("core product plan missing")
	}
	limit := policy.quotas[normalizeProductQuotaDimension(productQuotaHostedMutationUnits)]
	setProductUsageForTest(t, identity, principal.TenantID, productQuotaHostedMutationUnits, limit-1)

	decision, err := identity.consumeHostedProductQuota(principal, productCapabilityHostedServing, []ProductQuotaCharge{{Dimension: productQuotaHostedMutationUnits, Units: 1}})
	if err != nil || !decision.Allowed {
		t.Fatalf("last available product unit was not consumable: decision=%+v err=%v", decision, err)
	}
	blocked, err := identity.consumeHostedProductQuota(principal, productCapabilityHostedServing, []ProductQuotaCharge{{Dimension: productQuotaHostedMutationUnits, Units: 1}})
	if err != nil {
		t.Fatal(err)
	}
	if blocked.Allowed || blocked.RetryAfterSeconds < 1 {
		t.Fatalf("exhausted product quota did not fail closed: %+v", blocked)
	}

	reloaded, err := NewIdentityService(persistence)
	if err != nil {
		t.Fatal(err)
	}
	snapshot, err := reloaded.productEntitlementSnapshot(principal.TenantID)
	if err != nil {
		t.Fatal(err)
	}
	if got := productQuotaByDimension(t, snapshot, productQuotaHostedMutationUnits).Used; got != limit {
		t.Fatalf("product metering did not survive reload: got=%d want=%d", got, limit)
	}
}

func TestHOST009ProductPolicyDeniesBeforeProtectedProjection(t *testing.T) {
	t.Setenv(runtimeModeEnv, "hosted")
	resetV184HostedQuotaLimiter(t, 16, 10, 10)
	persistence, identity := newIdentityTestService(t)
	token, principal := authenticatedProductOwner(t, identity)
	app := &Application{
		identity:      identity,
		persistence:   persistence,
		state:         defaultState(),
		workspaces:    map[string]UserWorkspace{principal.UserID: defaultUserWorkspace(principal.UserID)},
		httpTelemetry: NewRequestTelemetry(),
	}
	called := 0
	handler := app.auth(func(w http.ResponseWriter, _ *http.Request) {
		called++
		w.WriteHeader(http.StatusNoContent)
	})

	if err := identity.setTenantProductPolicy(principal.TenantID, ProductPlanFoundationCore, ProductStatusSuspended); err != nil {
		t.Fatal(err)
	}
	blocked := httptest.NewRecorder()
	blockedReq := httptest.NewRequest(http.MethodGet, "/api/state", nil)
	blockedReq.AddCookie(&http.Cookie{Name: sessionCookieName, Value: token})
	handler.ServeHTTP(blocked, blockedReq)
	if blocked.Code != http.StatusForbidden || called != 0 {
		t.Fatalf("suspended product reached protected projection: code=%d called=%d body=%s", blocked.Code, called, blocked.Body.String())
	}

	if err := identity.setTenantProductPolicy(principal.TenantID, ProductPlanFoundationCore, ProductStatusGrace); err != nil {
		t.Fatal(err)
	}
	allowed := httptest.NewRecorder()
	allowedReq := httptest.NewRequest(http.MethodGet, "/api/state", nil)
	allowedReq.AddCookie(&http.Cookie{Name: sessionCookieName, Value: token})
	handler.ServeHTTP(allowed, allowedReq)
	if allowed.Code != http.StatusNoContent || called != 1 {
		t.Fatalf("grace product did not reach protected projection: code=%d called=%d body=%s", allowed.Code, called, allowed.Body.String())
	}
}

func TestHOST009AntiAbuseThrottleDoesNotDoubleChargeProductMetering(t *testing.T) {
	t.Setenv(runtimeModeEnv, "hosted")
	resetV184HostedQuotaLimiter(t, 16, 1, 1)
	persistence, identity := newIdentityTestService(t)
	token, principal := authenticatedProductOwner(t, identity)
	app := &Application{
		identity:      identity,
		persistence:   persistence,
		state:         defaultState(),
		workspaces:    map[string]UserWorkspace{principal.UserID: defaultUserWorkspace(principal.UserID)},
		httpTelemetry: NewRequestTelemetry(),
	}
	handler := app.auth(func(w http.ResponseWriter, _ *http.Request) { w.WriteHeader(http.StatusNoContent) })
	request := func() *http.Request {
		r := httptest.NewRequest(http.MethodPost, "/api/watchlists/create", nil)
		r.AddCookie(&http.Cookie{Name: sessionCookieName, Value: token})
		r.AddCookie(&http.Cookie{Name: csrfCookieName, Value: "csrf-product"})
		r.Header.Set("X-DE-PULSE-CSRF", "csrf-product")
		return r
	}

	first := httptest.NewRecorder()
	handler.ServeHTTP(first, request())
	if first.Code != http.StatusNoContent {
		t.Fatalf("first hosted mutation failed: %d %s", first.Code, first.Body.String())
	}
	before, err := identity.productEntitlementSnapshot(principal.TenantID)
	if err != nil {
		t.Fatal(err)
	}
	usedBefore := productQuotaByDimension(t, before, productQuotaHostedMutationUnits).Used

	second := httptest.NewRecorder()
	handler.ServeHTTP(second, request())
	if second.Code != http.StatusTooManyRequests {
		t.Fatalf("anti-abuse throttle did not block second mutation: %d %s", second.Code, second.Body.String())
	}
	after, err := identity.productEntitlementSnapshot(principal.TenantID)
	if err != nil {
		t.Fatal(err)
	}
	if got := productQuotaByDimension(t, after, productQuotaHostedMutationUnits).Used; got != usedBefore {
		t.Fatalf("anti-abuse rejection consumed product metering: before=%d after=%d", usedBefore, got)
	}
}

func TestHOST008AuthorizedProductStatusProjection(t *testing.T) {
	persistence, identity := newIdentityTestService(t)
	token, principal := authenticatedProductOwner(t, identity)
	app := &Application{
		identity:    identity,
		persistence: persistence,
		state:       defaultState(),
		workspaces:  map[string]UserWorkspace{principal.UserID: defaultUserWorkspace(principal.UserID)},
	}
	mux := http.NewServeMux()
	app.registerHealthRoutes(mux)

	r := httptest.NewRequest(http.MethodGet, "/api/auth/product/status", nil)
	r.AddCookie(&http.Cookie{Name: sessionCookieName, Value: token})
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, r)
	if w.Code != http.StatusOK {
		t.Fatalf("authorized product status failed: %d %s", w.Code, w.Body.String())
	}
	if !strings.Contains(w.Body.String(), string(ProductPlanFoundationCore)) || !strings.Contains(w.Body.String(), string(ProductStatusActive)) {
		t.Fatalf("authorized product status omitted canonical plan/status: %s", w.Body.String())
	}
}

func TestHOST009MeteringWindowResetsWithoutPlanTransitionReset(t *testing.T) {
	_, identity := newIdentityTestService(t)
	_, principal := authenticatedProductOwner(t, identity)
	base := time.Unix(2_400_000_000, 0)
	identity.now = func() time.Time { return base }
	if err := identity.setTenantProductPolicy(principal.TenantID, ProductPlanFoundationExtended, ProductStatusActive); err != nil {
		t.Fatal(err)
	}
	if _, err := identity.consumeHostedProductQuota(principal, productCapabilityHostedServing, []ProductQuotaCharge{{Dimension: productQuotaHostedMutationUnits, Units: 2}}); err != nil {
		t.Fatal(err)
	}
	before, err := identity.productEntitlementSnapshot(principal.TenantID)
	if err != nil {
		t.Fatal(err)
	}
	if got := productQuotaByDimension(t, before, productQuotaHostedMutationUnits).Used; got != 2 {
		t.Fatalf("unexpected initial product usage: %d", got)
	}
	if err := identity.setTenantProductPolicy(principal.TenantID, ProductPlanFoundationCore, ProductStatusActive); err != nil {
		t.Fatal(err)
	}
	afterTransition, err := identity.productEntitlementSnapshot(principal.TenantID)
	if err != nil {
		t.Fatal(err)
	}
	if got := productQuotaByDimension(t, afterTransition, productQuotaHostedMutationUnits).Used; got != 2 {
		t.Fatalf("plan transition reset active metering window: %d", got)
	}

	identity.now = func() time.Time { return base.Add(productMeteringWindow + time.Second) }
	afterWindow, err := identity.productEntitlementSnapshot(principal.TenantID)
	if err != nil {
		t.Fatal(err)
	}
	if got := productQuotaByDimension(t, afterWindow, productQuotaHostedMutationUnits).Used; got != 0 {
		t.Fatalf("expired metering window did not reset effective usage: %d", got)
	}
}
