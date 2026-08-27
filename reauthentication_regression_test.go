package main

import (
	"crypto/ed25519"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"errors"
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

func TestHOST006StaleRegisteredDevicePersistsLifecycleAndRevokesBoundSession(t *testing.T) {
	store, s := newIdentityTestService(t)
	base := time.Unix(2_350_000_000, 0)
	token, p, _ := v184CredentialedOwner(t, s, base)
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
	if err := s.persistLocked(); err != nil {
		s.mu.Unlock()
		t.Fatal(err)
	}
	s.mu.Unlock()

	decision := s.authorizeHostedIdentity(p, HostedIdentityRequirement{TenantID: localTenantID, Capability: hostedCapabilityStandardUse, RequireRegisteredDevice: true})
	if decision.Allowed {
		t.Fatalf("stale registered device remained authorized: %+v", decision)
	}
	if _, err := s.resolve(token, false); err == nil {
		t.Fatal("stale device transition did not revoke its bound session")
	}

	s.mu.Lock()
	status := DeviceStatus("")
	for _, candidate := range s.state.Devices {
		if candidate.ID == device.ID {
			status = candidate.Status
			break
		}
	}
	s.mu.Unlock()
	if status != DeviceStale {
		t.Fatalf("stale device was denied without durable STALE lifecycle state: status=%s", status)
	}

	reloaded, err := NewIdentityService(store)
	if err != nil {
		t.Fatal(err)
	}
	reloaded.mu.Lock()
	defer reloaded.mu.Unlock()
	persistedStatus := DeviceStatus("")
	boundSessionRevoked := false
	for _, candidate := range reloaded.state.Devices {
		if candidate.ID == device.ID {
			persistedStatus = candidate.Status
			break
		}
	}
	for _, session := range reloaded.state.Sessions {
		if session.ID == p.SessionID {
			boundSessionRevoked = session.RevokedAt > 0
			break
		}
	}
	if persistedStatus != DeviceStale || !boundSessionRevoked {
		t.Fatalf("stale lifecycle did not survive restart: status=%s sessionRevoked=%v", persistedStatus, boundSessionRevoked)
	}
}

func TestHOST007SensitiveHostedActionRequiresFreshReauthAndCryptographicMFAProof(t *testing.T) {
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
	publicKey, privateKey := hostedMFATestKey(t)
	credential, err := s.enrollHostedMFACredential(p, "security key", publicKey)
	if err != nil {
		t.Fatal(err)
	}
	challenge, err := s.createHostedMFAChallenge(p, credential.ID)
	if err != nil {
		t.Fatal(err)
	}
	payload, err := base64.RawURLEncoding.DecodeString(challenge.SigningPayload)
	if err != nil {
		t.Fatal(err)
	}
	signature := base64.RawURLEncoding.EncodeToString(ed25519.Sign(privateKey, payload))
	if err := s.verifyHostedMFAChallenge(p, challenge.ID, challenge.CredentialID, challenge.SigningPayload, signature); err != nil {
		t.Fatal(err)
	}
	if decision := s.authorizeHostedIdentity(p, requirement); !decision.Allowed {
		t.Fatalf("fresh reauth + cryptographic MFA proof blocked: %+v", decision)
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
	challenge, err = s.createHostedMFAChallenge(p, credential.ID)
	if err != nil {
		t.Fatal(err)
	}
	payload, err = base64.RawURLEncoding.DecodeString(challenge.SigningPayload)
	if err != nil {
		t.Fatal(err)
	}
	signature = base64.RawURLEncoding.EncodeToString(ed25519.Sign(privateKey, payload))
	if err := s.verifyHostedMFAChallenge(p, challenge.ID, challenge.CredentialID, challenge.SigningPayload, signature); err != nil {
		t.Fatal(err)
	}
	if decision := s.authorizeHostedIdentity(p, requirement); !decision.Allowed {
		t.Fatalf("refreshed sensitive assurance blocked: %+v", decision)
	}
}

func hostedIdentityHTTPMux(app *Application) *http.ServeMux {
	mux := http.NewServeMux()
	app.registerHealthRoutes(mux)
	return mux
}

func TestHOST006HostedDeviceHTTPRegistrationBindsCurrentSession(t *testing.T) {
	_, s := newIdentityTestService(t)
	base := time.Unix(2_410_000_000, 0)
	token, p, _ := v184CredentialedOwner(t, s, base)
	app := &Application{identity: s}
	mux := hostedIdentityHTTPMux(app)
	csrf := "hosted-device-register-csrf"
	fingerprint := "sha256:http-registered-device"

	res := httptest.NewRecorder()
	mux.ServeHTTP(res, v184AuthenticatedPOST("/api/auth/device/register", token, csrf, `{"label":"primary browser","fingerprintHash":"`+fingerprint+`"}`))
	if res.Code != http.StatusOK {
		t.Fatalf("device registration failed: code=%d body=%s", res.Code, res.Body.String())
	}
	if strings.Contains(res.Body.String(), fingerprint) {
		t.Fatalf("device fingerprint leaked in response: %s", res.Body.String())
	}
	var body struct {
		OK     bool `json:"ok"`
		Device struct {
			ID     string       `json:"id"`
			Label  string       `json:"label"`
			Status DeviceStatus `json:"status"`
		} `json:"device"`
		SessionBound bool `json:"sessionBound"`
	}
	if err := json.Unmarshal(res.Body.Bytes(), &body); err != nil {
		t.Fatal(err)
	}
	if !body.OK || !body.SessionBound || body.Device.ID == "" || body.Device.Status != DeviceActive {
		t.Fatalf("unexpected registration contract: %+v", body)
	}
	resolved, err := s.resolve(token, false)
	if err != nil {
		t.Fatalf("bound session stopped resolving: %v", err)
	}
	if resolved.DeviceID != body.Device.ID {
		t.Fatalf("session was not bound to registered device: principal=%+v device=%+v", resolved, body.Device)
	}
	decision := s.authorizeHostedIdentity(p, HostedIdentityRequirement{TenantID: localTenantID, Capability: hostedCapabilityDeviceManage, RequireRegisteredDevice: true})
	if !decision.Allowed || decision.DeviceID != body.Device.ID {
		t.Fatalf("registered HTTP device did not become canonical trust evidence: %+v", decision)
	}
}

func TestHOST006HostedDeviceHTTPStatusRequiresRegisteredCurrentDevice(t *testing.T) {
	_, s := newIdentityTestService(t)
	base := time.Unix(2_420_000_000, 0)
	token, _, _ := v184CredentialedOwner(t, s, base)
	app := &Application{identity: s}
	mux := hostedIdentityHTTPMux(app)
	csrf := "hosted-device-status-csrf"

	res := httptest.NewRecorder()
	mux.ServeHTTP(res, v184AuthenticatedPOST("/api/auth/device/status", token, csrf, `{"deviceId":"dev-untrusted","status":"REVOKED"}`))
	if res.Code != http.StatusForbidden {
		t.Fatalf("unregistered session crossed device-management trust boundary: code=%d body=%s", res.Code, res.Body.String())
	}
	if _, err := s.resolve(token, false); err != nil {
		t.Fatalf("denied status request mutated current session: %v", err)
	}
}

func TestHOST006HostedDeviceHTTPStatusRevokesCurrentBoundSession(t *testing.T) {
	_, s := newIdentityTestService(t)
	base := time.Unix(2_430_000_000, 0)
	token, _, _ := v184CredentialedOwner(t, s, base)
	app := &Application{identity: s}
	mux := hostedIdentityHTTPMux(app)
	csrf := "hosted-device-lost-csrf"

	registered := httptest.NewRecorder()
	mux.ServeHTTP(registered, v184AuthenticatedPOST("/api/auth/device/register", token, csrf, `{"label":"travel laptop","fingerprintHash":"sha256:http-lost-device"}`))
	if registered.Code != http.StatusOK {
		t.Fatalf("device registration failed: code=%d body=%s", registered.Code, registered.Body.String())
	}
	var registration struct {
		Device struct {
			ID string `json:"id"`
		} `json:"device"`
	}
	if err := json.Unmarshal(registered.Body.Bytes(), &registration); err != nil || registration.Device.ID == "" {
		t.Fatalf("missing registered device id: body=%s err=%v", registered.Body.String(), err)
	}

	lost := httptest.NewRecorder()
	mux.ServeHTTP(lost, v184AuthenticatedPOST("/api/auth/device/status", token, csrf, `{"deviceId":"`+registration.Device.ID+`","status":"LOST"}`))
	if lost.Code != http.StatusOK {
		t.Fatalf("trusted device status update failed: code=%d body=%s", lost.Code, lost.Body.String())
	}
	var body struct {
		OK                    bool         `json:"ok"`
		Status                DeviceStatus `json:"status"`
		CurrentSessionRevoked bool         `json:"currentSessionRevoked"`
	}
	if err := json.Unmarshal(lost.Body.Bytes(), &body); err != nil {
		t.Fatal(err)
	}
	if !body.OK || body.Status != DeviceLost || !body.CurrentSessionRevoked {
		t.Fatalf("unexpected lost-device contract: %+v", body)
	}
	if _, err := s.resolve(token, false); err == nil {
		t.Fatal("LOST current device did not revoke its bound session")
	}
	clearedSessionCookie := false
	for _, cookie := range lost.Result().Cookies() {
		if cookie.Name == sessionCookieName && cookie.MaxAge < 0 {
			clearedSessionCookie = true
			break
		}
	}
	if !clearedSessionCookie {
		t.Fatal("revoked current device did not clear the browser session cookie")
	}
}

func TestHOST006HostedDeviceHTTPDoesNotReactivateRevokedTrust(t *testing.T) {
	_, s := newIdentityTestService(t)
	base := time.Unix(2_440_000_000, 0)
	token, _, _ := v184CredentialedOwner(t, s, base)
	app := &Application{identity: s}
	mux := hostedIdentityHTTPMux(app)
	csrf := "hosted-device-active-csrf"

	registered := httptest.NewRecorder()
	mux.ServeHTTP(registered, v184AuthenticatedPOST("/api/auth/device/register", token, csrf, `{"label":"browser","fingerprintHash":"sha256:http-active-reject"}`))
	if registered.Code != http.StatusOK {
		t.Fatalf("device registration failed: code=%d body=%s", registered.Code, registered.Body.String())
	}
	var registration struct {
		Device struct {
			ID string `json:"id"`
		} `json:"device"`
	}
	if err := json.Unmarshal(registered.Body.Bytes(), &registration); err != nil {
		t.Fatal(err)
	}

	active := httptest.NewRecorder()
	mux.ServeHTTP(active, v184AuthenticatedPOST("/api/auth/device/status", token, csrf, `{"deviceId":"`+registration.Device.ID+`","status":"ACTIVE"}`))
	if active.Code != http.StatusBadRequest {
		t.Fatalf("HTTP surface allowed direct ACTIVE mutation: code=%d body=%s", active.Code, active.Body.String())
	}
	if _, err := s.resolve(token, false); err != nil {
		t.Fatalf("rejected ACTIVE request mutated current session: %v", err)
	}
}

func hostedMFATestKey(t *testing.T) (string, ed25519.PrivateKey) {
	t.Helper()
	publicKey, privateKey, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		t.Fatal(err)
	}
	return base64.RawURLEncoding.EncodeToString(publicKey), privateKey
}

func hostedMFAJSONPOST(t *testing.T, path, token, csrf string, payload any) *http.Request {
	t.Helper()
	body, err := json.Marshal(payload)
	if err != nil {
		t.Fatal(err)
	}
	return v184AuthenticatedPOST(path, token, csrf, string(body))
}

func hostedMFASessionDecision(s *IdentityService, p Principal) HostedIdentityDecision {
	return s.authorizeHostedIdentity(p, HostedIdentityRequirement{TenantID: p.TenantID, Capability: hostedCapabilityTenantManage, RequireRecentAuthentication: true, RequireMFA: true})
}

func mustJSONForHOST007(t *testing.T, value any) string {
	t.Helper()
	encoded, err := json.Marshal(value)
	if err != nil {
		t.Fatal(err)
	}
	return string(encoded)
}

func TestHOST007PublicKeyMFACeremonyEnforcesSensitiveAuthorizationAndRejectsReplay(t *testing.T) {
	resetV184LoginLimiter(t)
	_, s := newIdentityTestService(t)
	base := time.Unix(2_450_000_000, 0)
	token, p, _ := v184CredentialedOwner(t, s, base)
	publicKey, privateKey := hostedMFATestKey(t)
	app := &Application{identity: s}
	mux := hostedIdentityHTTPMux(app)
	csrf := "host007-mfa-csrf"

	enrolled := httptest.NewRecorder()
	mux.ServeHTTP(enrolled, hostedMFAJSONPOST(t, "/api/auth/mfa/credential/enroll", token, csrf, map[string]any{"label": "primary security key", "algorithm": hostedMFAAlgorithmEd25519, "publicKey": publicKey}))
	if enrolled.Code != http.StatusOK {
		t.Fatalf("MFA enrollment failed: code=%d body=%s", enrolled.Code, enrolled.Body.String())
	}
	if strings.Contains(enrolled.Body.String(), publicKey) {
		t.Fatalf("MFA public key leaked from credential response: %s", enrolled.Body.String())
	}
	var enrollment struct {
		Credential HostedMFACredentialView `json:"credential"`
	}
	if err := json.Unmarshal(enrolled.Body.Bytes(), &enrollment); err != nil || enrollment.Credential.ID == "" {
		t.Fatalf("missing enrolled credential: body=%s err=%v", enrolled.Body.String(), err)
	}
	if decision := hostedMFASessionDecision(s, p); decision.Allowed {
		t.Fatalf("sensitive hosted action allowed before MFA ceremony: %+v", decision)
	}

	challenged := httptest.NewRecorder()
	mux.ServeHTTP(challenged, hostedMFAJSONPOST(t, "/api/auth/mfa/challenge", token, csrf, map[string]any{"credentialId": enrollment.Credential.ID}))
	if challenged.Code != http.StatusOK {
		t.Fatalf("MFA challenge failed: code=%d body=%s", challenged.Code, challenged.Body.String())
	}
	var challengeBody struct {
		Challenge HostedMFAChallenge `json:"challenge"`
	}
	if err := json.Unmarshal(challenged.Body.Bytes(), &challengeBody); err != nil {
		t.Fatal(err)
	}
	challenge := challengeBody.Challenge
	payload, err := base64.RawURLEncoding.DecodeString(challenge.SigningPayload)
	if err != nil || challenge.ID == "" || challenge.CredentialID != enrollment.Credential.ID {
		t.Fatalf("invalid MFA challenge contract: %+v err=%v", challenge, err)
	}
	signature := base64.RawURLEncoding.EncodeToString(ed25519.Sign(privateKey, payload))
	verification := map[string]any{"challengeId": challenge.ID, "credentialId": challenge.CredentialID, "signingPayload": challenge.SigningPayload, "signature": signature}

	missingCSRF := httptest.NewRequest(http.MethodPost, "/api/auth/mfa/verify", strings.NewReader(mustJSONForHOST007(t, verification)))
	missingCSRF.AddCookie(&http.Cookie{Name: sessionCookieName, Value: token})
	missingCSRF.Header.Set("Content-Type", "application/json")
	blocked := httptest.NewRecorder()
	mux.ServeHTTP(blocked, missingCSRF)
	if blocked.Code != http.StatusForbidden {
		t.Fatalf("MFA verify crossed CSRF boundary: code=%d body=%s", blocked.Code, blocked.Body.String())
	}

	verified := httptest.NewRecorder()
	mux.ServeHTTP(verified, hostedMFAJSONPOST(t, "/api/auth/mfa/verify", token, csrf, verification))
	if verified.Code != http.StatusOK {
		t.Fatalf("valid MFA signature rejected: code=%d body=%s", verified.Code, verified.Body.String())
	}
	if decision := hostedMFASessionDecision(s, p); !decision.Allowed {
		t.Fatalf("cryptographically verified MFA proof did not satisfy sensitive authorization: %+v", decision)
	}
	proofAt := v184SessionByID(t, s, p.SessionID).MFAVerifiedAt
	if proofAt != base.UnixMilli() {
		t.Fatalf("unexpected MFA proof time: got=%d want=%d", proofAt, base.UnixMilli())
	}

	replay := httptest.NewRecorder()
	mux.ServeHTTP(replay, hostedMFAJSONPOST(t, "/api/auth/mfa/verify", token, csrf, verification))
	if replay.Code != http.StatusForbidden {
		t.Fatalf("MFA challenge replay was accepted: code=%d body=%s", replay.Code, replay.Body.String())
	}
	if got := v184SessionByID(t, s, p.SessionID).MFAVerifiedAt; got != proofAt {
		t.Fatalf("replay mutated MFA proof timestamp: before=%d after=%d", proofAt, got)
	}
}

func TestHOST007InvalidSignatureAndCrossSessionChallengeFailClosed(t *testing.T) {
	resetV184LoginLimiter(t)
	_, s := newIdentityTestService(t)
	base := time.Unix(2_460_000_000, 0)
	_, p1, password := v184CredentialedOwner(t, s, base)
	publicKey, _ := hostedMFATestKey(t)
	credential, err := s.enrollHostedMFACredential(p1, "primary", publicKey)
	if err != nil {
		t.Fatal(err)
	}
	challenge, err := s.createHostedMFAChallenge(p1, credential.ID)
	if err != nil {
		t.Fatal(err)
	}
	payload, err := base64.RawURLEncoding.DecodeString(challenge.SigningPayload)
	if err != nil {
		t.Fatal(err)
	}
	_, wrongPrivateKey := hostedMFATestKey(t)
	wrongSignature := base64.RawURLEncoding.EncodeToString(ed25519.Sign(wrongPrivateKey, payload))

	_, p2, err := s.authenticate("owner", password)
	if err != nil {
		t.Fatal(err)
	}
	if err := s.verifyHostedMFAChallenge(p2, challenge.ID, challenge.CredentialID, challenge.SigningPayload, wrongSignature); !errors.Is(err, errHostedMFAChallengeMissing) {
		t.Fatalf("cross-session challenge was not denied: err=%v", err)
	}
	if decision := hostedMFASessionDecision(s, p2); decision.Allowed {
		t.Fatalf("cross-session attempt created MFA proof: %+v", decision)
	}
	if err := s.verifyHostedMFAChallenge(p1, challenge.ID, challenge.CredentialID, challenge.SigningPayload, wrongSignature); !errors.Is(err, errHostedMFAInvalidProof) {
		t.Fatalf("invalid signature was not rejected: err=%v", err)
	}
	if decision := hostedMFASessionDecision(s, p1); decision.Allowed {
		t.Fatalf("invalid signature created MFA proof: %+v", decision)
	}
	if err := s.verifyHostedMFAChallenge(p1, challenge.ID, challenge.CredentialID, challenge.SigningPayload, wrongSignature); !errors.Is(err, errHostedMFAChallengeMissing) {
		t.Fatalf("failed verification did not consume one-time challenge: err=%v", err)
	}

	events, err := s.adminSecurityEvents(p1)
	if err != nil {
		t.Fatal(err)
	}
	foundFailure := false
	for _, event := range events {
		if event.Type == IdentitySecurityMFAVerificationFailed && event.UserID == p1.UserID && event.SessionID == p1.SessionID {
			foundFailure = true
			break
		}
	}
	if !foundFailure {
		t.Fatal("MFA verification failure was not observable in canonical security events")
	}
}

func TestHOST007And164MFACredentialChallengeProofRotationAndRevocationPersist(t *testing.T) {
	resetV184LoginLimiter(t)
	store, s := newIdentityTestService(t)
	base := time.Unix(2_470_000_000, 0)
	token, p, _ := v184CredentialedOwner(t, s, base)
	publicKey, privateKey := hostedMFATestKey(t)
	credential, err := s.enrollHostedMFACredential(p, "persistent key", publicKey)
	if err != nil {
		t.Fatal(err)
	}
	challenge, err := s.createHostedMFAChallenge(p, credential.ID)
	if err != nil {
		t.Fatal(err)
	}

	reloaded, err := NewIdentityService(store)
	if err != nil {
		t.Fatal(err)
	}
	reloaded.now = func() time.Time { return base }
	credentials, recentMFA, err := reloaded.listHostedMFACredentials(p)
	if err != nil || recentMFA || len(credentials) != 1 || credentials[0].ID != credential.ID {
		t.Fatalf("MFA credential state did not survive restart: credentials=%+v recent=%v err=%v", credentials, recentMFA, err)
	}
	payload, err := base64.RawURLEncoding.DecodeString(challenge.SigningPayload)
	if err != nil {
		t.Fatal(err)
	}
	signature := base64.RawURLEncoding.EncodeToString(ed25519.Sign(privateKey, payload))
	if err := reloaded.verifyHostedMFAChallenge(p, challenge.ID, challenge.CredentialID, challenge.SigningPayload, signature); err != nil {
		t.Fatalf("persisted challenge could not complete after restart: %v", err)
	}
	if decision := hostedMFASessionDecision(reloaded, p); !decision.Allowed {
		t.Fatalf("persisted ceremony did not create canonical MFA assurance: %+v", decision)
	}

	rotatedToken, rotatedPrincipal, err := reloaded.rotate(token)
	if err != nil || rotatedToken == "" || rotatedPrincipal.SessionID == p.SessionID {
		t.Fatalf("session rotation failed after MFA: token=%q principal=%+v err=%v", rotatedToken, rotatedPrincipal, err)
	}
	if decision := hostedMFASessionDecision(reloaded, rotatedPrincipal); !decision.Allowed {
		t.Fatalf("valid MFA assurance did not survive secure session rotation: %+v", decision)
	}
	if err := reloaded.revokeHostedMFACredential(rotatedPrincipal, credential.ID); err != nil {
		t.Fatal(err)
	}
	if decision := hostedMFASessionDecision(reloaded, rotatedPrincipal); decision.Allowed {
		t.Fatalf("revoked MFA credential left session assurance active: %+v", decision)
	}
	if _, err := reloaded.createHostedMFAChallenge(rotatedPrincipal, credential.ID); !errors.Is(err, errHostedMFACredentialMissing) {
		t.Fatalf("revoked MFA credential remained challengeable: err=%v", err)
	}

	persistedAgain, err := NewIdentityService(store)
	if err != nil {
		t.Fatal(err)
	}
	persistedAgain.now = func() time.Time { return base }
	credentials, recentMFA, err = persistedAgain.listHostedMFACredentials(rotatedPrincipal)
	if err != nil || recentMFA || len(credentials) != 1 || credentials[0].RevokedAt == 0 {
		t.Fatalf("MFA revocation did not survive restart: credentials=%+v recent=%v err=%v", credentials, recentMFA, err)
	}
}
