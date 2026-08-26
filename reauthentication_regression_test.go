package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func v184CredentialedOwner(t *testing.T, s *IdentityService, base time.Time) (string, Principal, string) {
	t.Helper()
	s.now = func() time.Time { return base }
	_, boot, err := s.bootstrapOwnerSession()
	if err != nil {
		t.Fatal(err)
	}
	password := "v18.4 owner security passphrase"
	token, principal, err := s.setPassword(boot.UserID, password)
	if err != nil {
		t.Fatal(err)
	}
	return token, principal, password
}

func v184AuthenticatedPOST(path, token, csrf, body string) *http.Request {
	r := httptest.NewRequest(http.MethodPost, path, strings.NewReader(body))
	r.AddCookie(&http.Cookie{Name: sessionCookieName, Value: token})
	r.AddCookie(&http.Cookie{Name: csrfCookieName, Value: csrf})
	r.Header.Set("X-DE-PULSE-CSRF", csrf)
	r.Header.Set("Content-Type", "application/json")
	return r
}

func TestV184ReauthenticationRefreshesOnlyAuthenticationAge(t *testing.T) {
	_, s := newIdentityTestService(t)
	base := time.Unix(2_300_000_000, 0)
	token, p, password := v184CredentialedOwner(t, s, base)
	original := v184SessionByID(t, s, p.SessionID)

	s.now = func() time.Time { return base.Add(20 * time.Minute) }
	if s.sessionRecentlyAuthenticated(p.SessionID, defaultSensitiveReauthTTL) {
		t.Fatal("aged session unexpectedly remained recent")
	}
	ok, err := s.reauthenticateSession(p.SessionID, password)
	if err != nil || !ok {
		t.Fatalf("reauthentication failed: ok=%v err=%v", ok, err)
	}
	current := v184SessionByID(t, s, p.SessionID)
	if current.AuthenticatedAt != base.Add(20*time.Minute).UnixMilli() {
		t.Fatalf("authentication age was not refreshed: %+v", current)
	}
	if current.AbsoluteExpiresAt != original.AbsoluteExpiresAt || current.IdleExpiresAt != original.IdleExpiresAt || current.CreatedAt != original.CreatedAt {
		t.Fatalf("reauthentication changed session lifetime: before=%+v after=%+v", original, current)
	}
	if _, err := s.resolve(token, false); err != nil {
		t.Fatalf("reauthentication invalidated current token: %v", err)
	}
}

func TestV184WrongReauthenticationDoesNotRefreshOrRevokeSession(t *testing.T) {
	resetV184LoginLimiter(t)
	_, s := newIdentityTestService(t)
	base := time.Unix(2_310_000_000, 0)
	token, p, _ := v184CredentialedOwner(t, s, base)
	s.now = func() time.Time { return base.Add(20 * time.Minute) }
	before := v184SessionByID(t, s, p.SessionID)
	ok, err := s.reauthenticateSession(p.SessionID, "wrong password")
	if err != nil || ok {
		t.Fatalf("wrong password result: ok=%v err=%v", ok, err)
	}
	after := v184SessionByID(t, s, p.SessionID)
	if after.AuthenticatedAt != before.AuthenticatedAt || after.RevokedAt != 0 {
		t.Fatalf("wrong password mutated session: before=%+v after=%+v", before, after)
	}
	if _, err := s.resolve(token, false); err != nil {
		t.Fatalf("wrong reauth password revoked current session: %v", err)
	}
}

func TestV184ProtectedMutationRequiresReauthThenSucceeds(t *testing.T) {
	resetV184LoginLimiter(t)
	_, s := newIdentityTestService(t)
	base := time.Unix(2_320_000_000, 0)
	token, p, password := v184CredentialedOwner(t, s, base)
	app := &Application{identity: s}
	s.now = func() time.Time { return base.Add(20 * time.Minute) }

	called := 0
	protected := app.requireRole(RoleAdmin, app.requireRecentAuthentication(postOnly(func(w http.ResponseWriter, r *http.Request) {
		called++
		w.WriteHeader(http.StatusNoContent)
	})))
	csrf := "v184-reauth-csrf"
	blocked := httptest.NewRecorder()
	protected(blocked, v184AuthenticatedPOST("/sensitive", token, csrf, `{}`))
	if blocked.Code != http.StatusPreconditionRequired || called != 0 {
		t.Fatalf("stale privileged session crossed recent-auth boundary: code=%d called=%d body=%s", blocked.Code, called, blocked.Body.String())
	}
	var blockedBody map[string]any
	if err := json.Unmarshal(blocked.Body.Bytes(), &blockedBody); err != nil || blockedBody["code"] != "REAUTH_REQUIRED" {
		t.Fatalf("missing machine-readable reauth contract: %s err=%v", blocked.Body.String(), err)
	}

	reauth := app.auth(postOnly(app.handleReauthenticate))
	wrong := httptest.NewRecorder()
	reauth(wrong, v184AuthenticatedPOST("/api/auth/reauth", token, csrf, `{"password":"wrong password"}`))
	if wrong.Code != http.StatusForbidden {
		t.Fatalf("wrong password did not stay inside reauth dialog contract: code=%d body=%s", wrong.Code, wrong.Body.String())
	}
	if called != 0 {
		t.Fatal("wrong reauth password allowed protected mutation")
	}

	verified := httptest.NewRecorder()
	reauth(verified, v184AuthenticatedPOST("/api/auth/reauth", token, csrf, `{"password":"`+password+`"}`))
	if verified.Code != http.StatusOK {
		t.Fatalf("valid reauth failed: code=%d body=%s", verified.Code, verified.Body.String())
	}
	if !s.sessionRecentlyAuthenticated(p.SessionID, defaultSensitiveReauthTTL) {
		t.Fatal("valid reauth did not refresh recent-auth state")
	}

	allowed := httptest.NewRecorder()
	protected(allowed, v184AuthenticatedPOST("/sensitive", token, csrf, `{}`))
	if allowed.Code != http.StatusNoContent || called != 1 {
		t.Fatalf("recently reauthenticated principal blocked: code=%d called=%d body=%s", allowed.Code, called, allowed.Body.String())
	}
}

func TestHOST004LegacyIdentityMigratesToExplicitLocalTenant(t *testing.T) {
	p, s := newIdentityTestService(t)
	s.mu.Lock()
	if len(s.state.Tenants) == 0 || s.state.Tenants[0].ID != localTenantID || s.state.Tenants[0].Status != TenantActive {
		s.mu.Unlock()
		t.Fatalf("local tenant migration missing: %+v", s.state.Tenants)
	}
	for _, u := range s.state.Users {
		if normalizedTenantID(u.TenantID) != localTenantID {
			s.mu.Unlock()
			t.Fatalf("legacy user not tenant-bound: %+v", u)
		}
	}
	s.mu.Unlock()
	reloaded, err := NewIdentityService(p)
	if err != nil {
		t.Fatal(err)
	}
	reloaded.mu.Lock()
	defer reloaded.mu.Unlock()
	if len(reloaded.state.Tenants) == 0 || reloaded.state.Tenants[0].ID != localTenantID {
		t.Fatalf("tenant migration did not persist: %+v", reloaded.state.Tenants)
	}
}

func TestHOST005RoleCapabilitiesAreExplicitAndIndependent(t *testing.T) {
	cases := []struct {
		role       UserRole
		capability string
		want       bool
	}{
		{RoleSuperOwner, hostedCapabilityTenantManage, true},
		{RoleOwner, hostedCapabilityAccountManage, true},
		{RoleAdmin, hostedCapabilityUserManage, true},
		{RoleAdmin, hostedCapabilityTenantManage, false},
		{RoleUser, hostedCapabilityStandardUse, true},
		{RoleUser, hostedCapabilityUserManage, false},
		{RoleDemo, hostedCapabilityDemoUse, true},
		{RoleDemo, hostedCapabilityStandardUse, false},
	}
	for _, tc := range cases {
		if got := roleHasHostedCapability(tc.role, tc.capability); got != tc.want {
			t.Fatalf("role=%s capability=%s got=%v want=%v", tc.role, tc.capability, got, tc.want)
		}
	}
}

func TestHOST004HostedIdentityDeniesCrossTenantAndInactiveTenant(t *testing.T) {
	_, s := newIdentityTestService(t)
	base := time.Unix(2_330_000_000, 0)
	_, p, _ := v184CredentialedOwner(t, s, base)

	cross := s.authorizeHostedIdentity(p, HostedIdentityRequirement{TenantID: "tenant-other", Capability: hostedCapabilityStandardUse})
	if cross.Allowed {
		t.Fatalf("cross-tenant principal was allowed: %+v", cross)
	}

	s.mu.Lock()
	for i := range s.state.Tenants {
		if s.state.Tenants[i].ID == localTenantID {
			s.state.Tenants[i].Status = TenantDisabled
		}
	}
	_ = s.persistLocked()
	s.mu.Unlock()
	inactive := s.authorizeHostedIdentity(p, HostedIdentityRequirement{TenantID: localTenantID, Capability: hostedCapabilityStandardUse})
	if inactive.Allowed {
		t.Fatalf("disabled tenant was allowed: %+v", inactive)
	}
}

func TestHOST006DeviceLifecycleRevokesBoundSessionsAndBlocksCrossTenantMutation(t *testing.T) {
	_, s := newIdentityTestService(t)
	base := time.Unix(2_340_000_000, 0)
	token, p, _ := v184CredentialedOwner(t, s, base)
	device, err := s.registerHostedDevice(p, "browser", "sha256:test-device")
	if err != nil {
		t.Fatal(err)
	}
	if err := s.bindHostedDeviceToSession(p, device.ID); err != nil {
		t.Fatal(err)
	}
	allowed := s.authorizeHostedIdentity(p, HostedIdentityRequirement{TenantID: localTenantID, Capability: hostedCapabilityStandardUse, RequireRegisteredDevice: true})
	if !allowed.Allowed || allowed.DeviceID != device.ID {
		t.Fatalf("active registered device blocked: %+v", allowed)
	}

	crossTenant := p
	crossTenant.TenantID = "tenant-other"
	if err := s.setHostedDeviceStatus(crossTenant, device.ID, DeviceLost); err == nil {
		t.Fatal("cross-tenant device mutation was allowed")
	}
	if err := s.setHostedDeviceStatus(p, device.ID, DeviceLost); err != nil {
		t.Fatal(err)
	}
	if _, err := s.resolve(token, false); err == nil {
		t.Fatal("lost device did not revoke its bound session")
	}
	blocked := s.authorizeHostedIdentity(p, HostedIdentityRequirement{TenantID: localTenantID, Capability: hostedCapabilityStandardUse, RequireRegisteredDevice: true})
	if blocked.Allowed {
		t.Fatalf("lost device/session remained authorized: %+v", blocked)
	}
}

func TestHOST006StaleRegisteredDeviceFailsClosed(t *testing.T) {
	_, s := newIdentityTestService(t)
	base := time.Unix(2_350_000_000, 0)
	_, p, _ := v184CredentialedOwner(t, s, base)
	device, err := s.registerHostedDevice(p, "browser", "sha256:stale-device")
	if err != nil {
		t.Fatal(err)
	}
	if err := s.bindHostedDeviceToSession(p, device.ID); err != nil {
		t.Fatal(err)
	}
	s.mu.Lock()
	for i := range s.state.Devices {
		if s.state.Devices[i].ID == device.ID {
			s.state.Devices[i].LastSeenAt = base.Add(-31 * 24 * time.Hour).UnixMilli()
		}
	}
	_ = s.persistLocked()
	s.mu.Unlock()
	decision := s.authorizeHostedIdentity(p, HostedIdentityRequirement{TenantID: localTenantID, Capability: hostedCapabilityStandardUse, RequireRegisteredDevice: true})
	if decision.Allowed {
		t.Fatalf("stale registered device remained authorized: %+v", decision)
	}
}

func TestHOST007SensitiveHostedActionRequiresFreshReauthAndExplicitMFAProof(t *testing.T) {
	_, s := newIdentityTestService(t)
	base := time.Unix(2_360_000_000, 0)
	_, p, password := v184CredentialedOwner(t, s, base)
	device, err := s.registerHostedDevice(p, "browser", "sha256:mfa-device")
	if err != nil {
		t.Fatal(err)
	}
	if err := s.bindHostedDeviceToSession(p, device.ID); err != nil {
		t.Fatal(err)
	}
	requirement := HostedIdentityRequirement{TenantID: localTenantID, Capability: hostedCapabilityTenantManage, RequireRegisteredDevice: true, RequireRecentAuthentication: true, RequireMFA: true}
	if decision := s.authorizeHostedIdentity(p, requirement); decision.Allowed {
		t.Fatalf("sensitive action allowed without MFA proof: %+v", decision)
	}
	if err := s.recordHostedMFAVerification(p); err != nil {
		t.Fatal(err)
	}
	if decision := s.authorizeHostedIdentity(p, requirement); !decision.Allowed {
		t.Fatalf("fresh reauth + explicit MFA proof blocked: %+v", decision)
	}

	s.now = func() time.Time { return base.Add(16 * time.Minute) }
	if decision := s.authorizeHostedIdentity(p, requirement); decision.Allowed {
		t.Fatalf("stale sensitive assurance remained authorized: %+v", decision)
	}
	ok, err := s.reauthenticateSession(p.SessionID, password)
	if err != nil || !ok {
		t.Fatalf("reauth failed: ok=%v err=%v", ok, err)
	}
	if decision := s.authorizeHostedIdentity(p, requirement); decision.Allowed {
		t.Fatalf("password reauth alone incorrectly satisfied MFA: %+v", decision)
	}
	if err := s.recordHostedMFAVerification(p); err != nil {
		t.Fatal(err)
	}
	if decision := s.authorizeHostedIdentity(p, requirement); !decision.Allowed {
		t.Fatalf("refreshed sensitive assurance blocked: %+v", decision)
	}
}
