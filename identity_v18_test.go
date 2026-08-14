package main

import (
	"bytes"
	"context"
	"encoding/hex"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	vendoredargon2 "depulse/internal/vendorcrypto/argon2"
)

func newIdentityTestService(t *testing.T) (*PersistenceManager, *IdentityService) {
	t.Helper()
	p := NewPersistenceManager(t.TempDir())
	if !p.Diagnostics().Ready {
		t.Fatalf("persistence not ready: %+v", p.Diagnostics())
	}
	s, err := NewIdentityService(p)
	if err != nil {
		_ = p.Close()
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = p.Close() })
	return p, s
}

func TestV180VendoredArgon2idKnownVector(t *testing.T) {
	got := vendoredargon2.IDKey([]byte("password"), []byte("somesalt"), 1, 64, 1, 24)
	want, _ := hex.DecodeString("655ad15eac652dc59f7170a7332bf49b8469be1fdb9c28bb")
	if !bytes.Equal(got, want) {
		t.Fatalf("Argon2id vector mismatch: got %x want %x", got, want)
	}
}

func TestV180Argon2idPasswordPolicyAndVerification(t *testing.T) {
	if _, err := hashPasswordArgon2id("short"); err == nil {
		t.Fatal("short password accepted")
	}
	hash, err := hashPasswordArgon2id("correct horse battery staple")
	if err != nil {
		t.Fatal(err)
	}
	if !strings.HasPrefix(hash, "$argon2id$v=19$m=65536,t=3,p=4$") {
		t.Fatalf("unexpected PHC: %s", hash)
	}
	ok, err := verifyPasswordArgon2id("correct horse battery staple", hash)
	if err != nil || !ok {
		t.Fatalf("valid password rejected: %v", err)
	}
	ok, err = verifyPasswordArgon2id("wrong password", hash)
	if err != nil {
		t.Fatal(err)
	}
	if ok {
		t.Fatal("wrong password accepted")
	}
	if _, err = verifyPasswordArgon2id("x", "$argon2id$v=19$m=999999999,t=3,p=4$c29tZXNhbHQ$YWJjZGVmZ2hpamtsbW5vcA"); err == nil {
		t.Fatal("unbounded attacker-controlled memory parameter accepted")
	}
}

func TestV180BootstrapOwnerCredentialMigrationPersistenceAndNoPlaintext(t *testing.T) {
	dir := t.TempDir()
	p := NewPersistenceManager(dir)
	if !p.Diagnostics().Ready {
		t.Fatal(p.Diagnostics().LastError)
	}
	defer p.Close()
	s, err := NewIdentityService(p)
	if err != nil {
		t.Fatal(err)
	}
	if !s.ownerNeedsPassword() {
		t.Fatal("fresh install must require bootstrap owner credential")
	}
	oldToken, principal, err := s.bootstrapOwnerSession()
	if err != nil {
		t.Fatal(err)
	}
	if principal.Role != RoleOwner || principal.Username != "owner" {
		t.Fatalf("wrong bootstrap principal: %+v", principal)
	}
	password := "v18 secure owner passphrase"
	newToken, np, err := s.setPassword(principal.UserID, password)
	if err != nil {
		t.Fatal(err)
	}
	if np.Role != RoleOwner || s.ownerNeedsPassword() {
		t.Fatal("bootstrap migration did not complete")
	}
	if _, err := s.resolve(oldToken, false); err == nil {
		t.Fatal("pre-credential bootstrap token remained valid")
	}
	if _, err := s.resolve(newToken, false); err != nil {
		t.Fatalf("rotated credential session invalid: %v", err)
	}
	s2, err := NewIdentityService(p)
	if err != nil {
		t.Fatal(err)
	}
	loginToken, lp, err := s2.authenticate("OWNER", password)
	if err != nil {
		t.Fatal(err)
	}
	if lp.UserID != principal.UserID {
		t.Fatal("persistent identity changed across reload")
	}
	if _, err := s2.resolve(loginToken, false); err != nil {
		t.Fatal(err)
	}
	artifact := "depulse-v17.db"
	if p.Diagnostics().Backend == "file-fallback" {
		artifact = "persistent-intelligence-fallback.json"
	}
	raw, err := os.ReadFile(filepath.Join(dir, artifact))
	if err != nil {
		t.Fatalf("read %s persistence artifact: %v", p.Diagnostics().Backend, err)
	}
	if bytes.Contains(raw, []byte(password)) {
		t.Fatal("plaintext password found in persistence")
	}
	if bytes.Contains(raw, []byte(loginToken)) || bytes.Contains(raw, []byte(newToken)) {
		t.Fatal("raw session token found in persistence")
	}
}

func TestV180SessionRotationRevocationIdleAndAbsoluteExpiry(t *testing.T) {
	_, s := newIdentityTestService(t)
	base := time.Unix(2_000_000_000, 0)
	s.now = func() time.Time { return base }
	s.idleTTL = 10 * time.Minute
	s.absoluteTTL = time.Hour
	token, p, err := s.bootstrapOwnerSession()
	if err != nil {
		t.Fatal(err)
	}
	rotated, rp, err := s.rotate(token)
	if err != nil {
		t.Fatal(err)
	}
	if rp.UserID != p.UserID {
		t.Fatal("rotation changed principal")
	}
	if _, err := s.resolve(token, false); err == nil {
		t.Fatal("rotated token remained valid")
	}
	if _, err := s.resolve(rotated, false); err != nil {
		t.Fatal(err)
	}
	s.now = func() time.Time { return base.Add(11 * time.Minute) }
	if _, err := s.resolve(rotated, false); err == nil {
		t.Fatal("idle-expired session accepted")
	}

	s.now = func() time.Time { return base }
	s.idleTTL = 2 * time.Hour
	s.absoluteTTL = time.Hour
	token2, _, err := s.bootstrapOwnerSession()
	if err != nil {
		t.Fatal(err)
	}
	// Session absolute expiry was minted at +1h while idle expiry extends beyond it.
	s.now = func() time.Time { return base.Add(61 * time.Minute) }
	if _, err := s.resolve(token2, false); err == nil {
		t.Fatal("absolute-expired session accepted")
	}
}

func TestV180RoleEnforcementDeniesUserAllowsAdminHierarchy(t *testing.T) {
	_, s := newIdentityTestService(t)
	now := time.Now().UnixMilli()
	s.mu.Lock()
	user := UserRecord{ID: "usr_user", Username: "u", PasswordHash: "test-only-present", Role: RoleUser, Status: UserActive, CreatedAt: now, UpdatedAt: now}
	admin := UserRecord{ID: "usr_admin", Username: "a", PasswordHash: "test-only-present", Role: RoleAdmin, Status: UserActive, CreatedAt: now, UpdatedAt: now}
	s.state.Users = append(s.state.Users, user, admin)
	_ = s.persistLocked()
	userToken, _, err := s.createSessionLocked(user, "")
	if err != nil {
		s.mu.Unlock()
		t.Fatal(err)
	}
	adminToken, _, err := s.createSessionLocked(admin, "")
	s.mu.Unlock()
	if err != nil {
		t.Fatal(err)
	}
	app := &Application{identity: s}
	called := false
	h := app.requireRole(RoleAdmin, func(w http.ResponseWriter, r *http.Request) { called = true; w.WriteHeader(204) })
	rr := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.AddCookie(&http.Cookie{Name: sessionCookieName, Value: userToken})
	req.AddCookie(&http.Cookie{Name: csrfCookieName, Value: "csrf-user"})
	req.Header.Set("X-DE-PULSE-CSRF", "csrf-user")
	h(rr, req)
	if rr.Code != http.StatusForbidden || called {
		t.Fatalf("USER crossed ADMIN boundary: code=%d called=%v", rr.Code, called)
	}
	rr = httptest.NewRecorder()
	req = httptest.NewRequest(http.MethodGet, "/", nil)
	req.AddCookie(&http.Cookie{Name: sessionCookieName, Value: adminToken})
	req.AddCookie(&http.Cookie{Name: csrfCookieName, Value: "csrf-admin"})
	req.Header.Set("X-DE-PULSE-CSRF", "csrf-admin")
	h(rr, req)
	if rr.Code != http.StatusNoContent || !called {
		t.Fatalf("ADMIN blocked: code=%d called=%v", rr.Code, called)
	}
}

func TestV180HTTPLoginLogoutLifecycle(t *testing.T) {
	_, s := newIdentityTestService(t)
	app := &Application{identity: s}
	bootToken, p, err := s.bootstrapOwnerSession()
	if err != nil {
		t.Fatal(err)
	}
	_, _, err = s.setPassword(p.UserID, "owner password for login test")
	if err != nil {
		t.Fatal(err)
	}
	_ = s.revokeToken(bootToken)
	body := strings.NewReader(`{"username":"owner","password":"owner password for login test"}`)
	rr := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/api/auth/login", body)
	app.handleLogin(rr, req)
	if rr.Code != 200 {
		t.Fatalf("login failed %d %s", rr.Code, rr.Body.String())
	}
	var sessionCookie *http.Cookie
	for _, c := range rr.Result().Cookies() {
		if c.Name == sessionCookieName {
			sessionCookie = c
		}
	}
	if sessionCookie == nil || sessionCookie.Value == "" || !sessionCookie.HttpOnly || sessionCookie.SameSite != http.SameSiteStrictMode {
		t.Fatalf("bad session cookie: %+v", sessionCookie)
	}
	if _, err := s.resolve(sessionCookie.Value, false); err != nil {
		t.Fatal(err)
	}
	rr2 := httptest.NewRecorder()
	req2 := httptest.NewRequest(http.MethodPost, "/api/auth/logout", nil)
	req2.AddCookie(sessionCookie)
	app.handleLogout(rr2, req2)
	if rr2.Code != 200 {
		t.Fatalf("logout failed: %d", rr2.Code)
	}
	if _, err := s.resolve(sessionCookie.Value, false); err == nil {
		t.Fatal("logout did not revoke server-side session")
	}
}

func TestV180IdentityStateSurvivesMarketPersistenceOperations(t *testing.T) {
	p, s := newIdentityTestService(t)
	tok, pr, err := s.bootstrapOwnerSession()
	if err != nil {
		t.Fatal(err)
	}
	p.EnqueueSymbols([]SymbolRegistryRecord{{Symbol: "NVDA", Active: true, LastSeenAt: time.Now().UnixMilli(), ProviderEligible: true}})
	p.EnqueueQuotes(map[string]Quote{"NVDA": {Symbol: "NVDA", Price: 123.45, Source: "test", UpdatedAt: time.Now().UnixMilli()}})
	time.Sleep(100 * time.Millisecond)
	st, err := p.LoadIdentityState(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if len(st.Users) == 0 {
		t.Fatal("identity lost during market persistence")
	}
	if _, err := s.resolve(tok, false); err != nil {
		t.Fatalf("session lost during market writes for %s: %v", pr.Username, err)
	}
}

func cookieNamed(resp *http.Response, name string) *http.Cookie {
	for _, c := range resp.Cookies() {
		if c.Name == name {
			return c
		}
	}
	return nil
}

func TestV180BootstrapSetupBoundaryAndCSRFRoutes(t *testing.T) {
	_, s := newIdentityTestService(t)
	app := &Application{identity: s}
	h := app.routes()

	// Fresh v17 migration lands on the credential-setup surface, not the app.
	rr := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	h.ServeHTTP(rr, req)
	if rr.Code != http.StatusOK || !strings.Contains(rr.Body.String(), "Secure the Owner account") {
		t.Fatalf("fresh bootstrap did not land on setup: code=%d body=%q", rr.Code, rr.Body.String())
	}
	session := cookieNamed(rr.Result(), sessionCookieName)
	csrf := cookieNamed(rr.Result(), csrfCookieName)
	if session == nil || csrf == nil || session.Value == "" || csrf.Value == "" {
		t.Fatalf("bootstrap cookies missing: session=%+v csrf=%+v", session, csrf)
	}

	// Bootstrap-owner compatibility is deliberately setup-only.
	blocked := httptest.NewRecorder()
	blockedReq := httptest.NewRequest(http.MethodGet, "/protected", nil)
	blockedReq.AddCookie(session)
	app.auth(func(w http.ResponseWriter, r *http.Request) { w.WriteHeader(http.StatusNoContent) })(blocked, blockedReq)
	if blocked.Code != http.StatusPreconditionRequired {
		t.Fatalf("bootstrap session reached normal API: code=%d", blocked.Code)
	}

	// State-changing setup requires the double-submit CSRF token.
	noCSRF := httptest.NewRecorder()
	noCSRFReq := httptest.NewRequest(http.MethodPost, "/api/auth/set-password", strings.NewReader(`{"password":"bootstrap route password"}`))
	noCSRFReq.AddCookie(session)
	h.ServeHTTP(noCSRF, noCSRFReq)
	if noCSRF.Code != http.StatusForbidden {
		t.Fatalf("set-password without CSRF was not denied: %d", noCSRF.Code)
	}

	set := httptest.NewRecorder()
	setReq := httptest.NewRequest(http.MethodPost, "/api/auth/set-password", strings.NewReader(`{"password":"bootstrap route password"}`))
	setReq.AddCookie(session)
	setReq.AddCookie(csrf)
	setReq.Header.Set("X-DE-PULSE-CSRF", csrf.Value)
	h.ServeHTTP(set, setReq)
	if set.Code != http.StatusOK {
		t.Fatalf("credential setup failed: %d %s", set.Code, set.Body.String())
	}
	newSession := cookieNamed(set.Result(), sessionCookieName)
	newCSRF := cookieNamed(set.Result(), csrfCookieName)
	if newSession == nil || newCSRF == nil || newSession.Value == session.Value {
		t.Fatalf("credential setup did not rotate session/csrf: old=%+v new=%+v csrf=%+v", session, newSession, newCSRF)
	}
	if _, err := s.resolve(session.Value, false); err == nil {
		t.Fatal("bootstrap token remained valid after credential setup")
	}

	// Logout is also a state-changing operation and must reject missing CSRF.
	logoutDenied := httptest.NewRecorder()
	logoutDeniedReq := httptest.NewRequest(http.MethodPost, "/api/auth/logout", nil)
	logoutDeniedReq.AddCookie(newSession)
	h.ServeHTTP(logoutDenied, logoutDeniedReq)
	if logoutDenied.Code != http.StatusForbidden {
		t.Fatalf("logout without CSRF was not denied: %d", logoutDenied.Code)
	}
	if _, err := s.resolve(newSession.Value, false); err != nil {
		t.Fatalf("denied logout revoked session: %v", err)
	}

	logout := httptest.NewRecorder()
	logoutReq := httptest.NewRequest(http.MethodPost, "/api/auth/logout", nil)
	logoutReq.AddCookie(newSession)
	logoutReq.AddCookie(newCSRF)
	logoutReq.Header.Set("X-DE-PULSE-CSRF", newCSRF.Value)
	h.ServeHTTP(logout, logoutReq)
	if logout.Code != http.StatusOK {
		t.Fatalf("valid logout failed: %d %s", logout.Code, logout.Body.String())
	}
	if _, err := s.resolve(newSession.Value, false); err == nil {
		t.Fatal("logout route did not revoke session")
	}
}

func TestV180HTTPSessionRotateEndpoint(t *testing.T) {
	_, s := newIdentityTestService(t)
	app := &Application{identity: s}
	boot, p, err := s.bootstrapOwnerSession()
	if err != nil {
		t.Fatal(err)
	}
	active, _, err := s.setPassword(p.UserID, "owner password for rotate endpoint")
	if err != nil {
		t.Fatal(err)
	}
	_ = s.revokeToken(boot)

	csrf := "csrf-test-token"
	rr := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/api/auth/rotate", nil)
	req.AddCookie(&http.Cookie{Name: sessionCookieName, Value: active})
	req.AddCookie(&http.Cookie{Name: csrfCookieName, Value: csrf})
	req.Header.Set("X-DE-PULSE-CSRF", csrf)
	app.auth(postOnly(app.handleRotateSession))(rr, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("rotate failed: %d %s", rr.Code, rr.Body.String())
	}
	var next *http.Cookie
	for _, c := range rr.Result().Cookies() {
		if c.Name == sessionCookieName && c.Value != "" {
			next = c
		}
	}
	if next == nil {
		t.Fatal("rotated session cookie missing")
	}
	if _, err := s.resolve(active, false); err == nil {
		t.Fatal("old token survived HTTP rotation")
	}
	if _, err := s.resolve(next.Value, false); err != nil {
		t.Fatalf("rotated token invalid: %v", err)
	}
}
