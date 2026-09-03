package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func v184SessionByID(t *testing.T, s *IdentityService, sessionID string) SessionRecord {
	t.Helper()
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, rec := range s.state.Sessions {
		if rec.ID == sessionID {
			return rec
		}
	}
	t.Fatalf("session %s not found", sessionID)
	return SessionRecord{}
}

func TestV184RotationPreservesAuthenticationAndAbsoluteLifetime(t *testing.T) {
	_, s := newIdentityTestService(t)
	base := time.Unix(2_100_000_000, 0)
	s.now = func() time.Time { return base }
	s.idleTTL = 45 * time.Minute
	s.absoluteTTL = 2 * time.Hour

	token, principal, err := s.bootstrapOwnerSession()
	if err != nil {
		t.Fatal(err)
	}
	original := v184SessionByID(t, s, principal.SessionID)
	if original.AuthenticatedAt != base.UnixMilli() {
		t.Fatalf("fresh session auth timestamp mismatch: %+v", original)
	}
	wantAbsolute := base.Add(2 * time.Hour).UnixMilli()
	if original.AbsoluteExpiresAt != wantAbsolute {
		t.Fatalf("fresh absolute expiry mismatch: got %d want %d", original.AbsoluteExpiresAt, wantAbsolute)
	}

	s.now = func() time.Time { return base.Add(20 * time.Minute) }
	rotated, rp, err := s.rotate(token)
	if err != nil {
		t.Fatal(err)
	}
	if rotated == token {
		t.Fatal("rotation reused raw session token")
	}
	next := v184SessionByID(t, s, rp.SessionID)
	if next.CreatedAt != base.Add(20*time.Minute).UnixMilli() {
		t.Fatalf("rotation creation time mismatch: %+v", next)
	}
	if next.AuthenticatedAt != original.AuthenticatedAt {
		t.Fatalf("rotation incorrectly refreshed authentication: old=%d new=%d", original.AuthenticatedAt, next.AuthenticatedAt)
	}
	if next.AbsoluteExpiresAt != original.AbsoluteExpiresAt {
		t.Fatalf("rotation extended absolute lifetime: old=%d new=%d", original.AbsoluteExpiresAt, next.AbsoluteExpiresAt)
	}
	if s.sessionRecentlyAuthenticated(rp.SessionID, defaultSensitiveReauthTTL) {
		t.Fatal("20-minute-old authentication incorrectly counted as recent")
	}
	if !s.sessionRecentlyAuthenticated(rp.SessionID, 30*time.Minute) {
		t.Fatal("20-minute-old authentication not recognized within wider policy")
	}

	s.now = func() time.Time { return base.Add(30 * time.Minute) }
	if _, err := s.resolve(rotated, true); err != nil {
		t.Fatal(err)
	}
	if s.sessionRecentlyAuthenticated(rp.SessionID, defaultSensitiveReauthTTL) {
		t.Fatal("passive activity refreshed authentication recency")
	}

	s.now = func() time.Time { return base.Add(121 * time.Minute) }
	if _, err := s.resolve(rotated, false); err == nil {
		t.Fatal("rotated token exceeded original absolute lifetime")
	}
}

func TestV184LegacySessionAuthenticationFallsBackToCreationTime(t *testing.T) {
	_, s := newIdentityTestService(t)
	base := time.Unix(2_200_000_000, 0)
	s.now = func() time.Time { return base }
	_, p, err := s.bootstrapOwnerSession()
	if err != nil {
		t.Fatal(err)
	}
	s.mu.Lock()
	for i := range s.state.Sessions {
		if s.state.Sessions[i].ID == p.SessionID {
			s.state.Sessions[i].AuthenticatedAt = 0
		}
	}
	s.mu.Unlock()

	s.now = func() time.Time { return base.Add(4 * time.Minute) }
	if !s.sessionRecentlyAuthenticated(p.SessionID, 5*time.Minute) {
		t.Fatal("legacy session did not fall back to createdAt")
	}
	s.now = func() time.Time { return base.Add(6 * time.Minute) }
	if s.sessionRecentlyAuthenticated(p.SessionID, 5*time.Minute) {
		t.Fatal("legacy createdAt fallback exceeded recent-auth policy")
	}
}

func TestHOST007And164AuthStatusDiscoversPersistedSessionAndRejectsExpiry(t *testing.T) {
	store, s := newIdentityTestService(t)
	base := time.Unix(2_490_000_000, 0)
	s.now = func() time.Time { return base }
	s.idleTTL = 2 * time.Hour
	s.absoluteTTL = 2 * time.Hour

	_, owner, err := s.bootstrapOwnerSession()
	if err != nil {
		t.Fatal(err)
	}
	token, principal, err := s.setPassword(owner.UserID, "v19 persisted discovery passphrase")
	if err != nil {
		t.Fatal(err)
	}

	reloaded, err := NewIdentityService(store)
	if err != nil {
		t.Fatal(err)
	}
	reloaded.now = func() time.Time { return base.Add(30 * time.Minute) }
	app := &Application{identity: reloaded}

	statusRequest := func() *http.Request {
		req := httptest.NewRequest(http.MethodGet, "/api/auth/status", nil)
		req.AddCookie(&http.Cookie{Name: sessionCookieName, Value: token})
		return req
	}
	readStatus := func(rr *httptest.ResponseRecorder) struct {
		Authenticated bool      `json:"authenticated"`
		Principal     Principal `json:"principal"`
	} {
		var body struct {
			Authenticated bool      `json:"authenticated"`
			Principal     Principal `json:"principal"`
		}
		if err := json.Unmarshal(rr.Body.Bytes(), &body); err != nil {
			t.Fatalf("decode auth status: %v body=%s", err, rr.Body.String())
		}
		return body
	}

	valid := httptest.NewRecorder()
	app.handleAuthStatus(valid, statusRequest())
	if valid.Code != http.StatusOK {
		t.Fatalf("persisted session discovery failed: code=%d body=%s", valid.Code, valid.Body.String())
	}
	validStatus := readStatus(valid)
	if !validStatus.Authenticated || validStatus.Principal.UserID != principal.UserID || validStatus.Principal.Role != principal.Role {
		t.Fatalf("persisted session was not rediscovered with canonical principal: %+v", validStatus)
	}

	reloaded.now = func() time.Time { return base.Add(121 * time.Minute) }
	expired := httptest.NewRecorder()
	app.handleAuthStatus(expired, statusRequest())
	if expired.Code != http.StatusOK {
		t.Fatalf("expired session discovery status failed: code=%d body=%s", expired.Code, expired.Body.String())
	}
	expiredStatus := readStatus(expired)
	if expiredStatus.Authenticated {
		t.Fatalf("expired persisted session was reported authenticated: %+v", expiredStatus)
	}
}

func TestHOST007And164FailedSessionRenewalFailsClosedAndClearsAuthCookies(t *testing.T) {
	_, s := newIdentityTestService(t)
	base := time.Unix(2_500_000_000, 0)
	s.now = func() time.Time { return base }
	s.idleTTL = 2 * time.Hour
	s.absoluteTTL = time.Hour

	_, owner, err := s.bootstrapOwnerSession()
	if err != nil {
		t.Fatal(err)
	}
	token, _, err := s.setPassword(owner.UserID, "v19 renewal failure passphrase")
	if err != nil {
		t.Fatal(err)
	}

	s.now = func() time.Time { return base.Add(61 * time.Minute) }
	app := &Application{identity: s}
	csrf := "host164-renewal-csrf"
	rr := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/api/auth/rotate", nil)
	req.AddCookie(&http.Cookie{Name: sessionCookieName, Value: token})
	req.AddCookie(&http.Cookie{Name: csrfCookieName, Value: csrf})
	req.Header.Set("X-DE-PULSE-CSRF", csrf)
	app.auth(postOnly(app.handleRotateSession))(rr, req)

	if rr.Code != http.StatusUnauthorized {
		t.Fatalf("expired session renewal did not require authentication: code=%d body=%s", rr.Code, rr.Body.String())
	}
	sessionCookie := cookieNamed(rr.Result(), sessionCookieName)
	csrfCookie := cookieNamed(rr.Result(), csrfCookieName)
	if sessionCookie == nil || sessionCookie.Value != "" || sessionCookie.MaxAge >= 0 {
		t.Fatalf("failed renewal did not clear session cookie: %+v", sessionCookie)
	}
	if csrfCookie == nil || csrfCookie.Value != "" || csrfCookie.MaxAge >= 0 {
		t.Fatalf("failed renewal did not clear CSRF cookie: %+v", csrfCookie)
	}
	if _, err := s.resolve(token, false); err == nil {
		t.Fatal("expired renewal token remained usable")
	}
}

func TestHOST007And164RoleDowngradeRevokesActiveSessionAndFreshLoginUsesNewRole(t *testing.T) {
	_, s := newIdentityTestService(t)
	_, owner, err := s.bootstrapOwnerSession()
	if err != nil {
		t.Fatal(err)
	}
	admin, err := s.adminCreateUser(owner, "host164.admin", "HOST 164 Admin", RoleAdmin, "temporary host164 admin password")
	if err != nil {
		t.Fatal(err)
	}
	activeToken, activePrincipal, err := s.setPassword(admin.ID, "host164 active admin passphrase")
	if err != nil {
		t.Fatal(err)
	}
	if activePrincipal.Role != RoleAdmin || !roleHasHostedCapability(activePrincipal.Role, hostedCapabilityUserManage) {
		t.Fatalf("active ADMIN principal missing expected capability: %+v", activePrincipal)
	}

	if err := s.adminSetUserRole(owner, admin.ID, RoleUser); err != nil {
		t.Fatal(err)
	}
	if _, err := s.resolve(activeToken, false); err == nil {
		t.Fatal("role downgrade left the old privileged session valid")
	}

	freshToken, freshPrincipal, err := s.authenticate(admin.Username, "host164 active admin passphrase")
	if err != nil {
		t.Fatal(err)
	}
	defer s.revokeToken(freshToken)
	if freshPrincipal.Role != RoleUser {
		t.Fatalf("fresh login retained stale role after downgrade: %+v", freshPrincipal)
	}
	if roleHasHostedCapability(freshPrincipal.Role, hostedCapabilityUserManage) {
		t.Fatalf("downgraded USER retained admin capability: %+v", freshPrincipal)
	}
	if !roleHasHostedCapability(freshPrincipal.Role, hostedCapabilityStandardUse) {
		t.Fatalf("downgraded USER lost normal product capability: %+v", freshPrincipal)
	}
}
