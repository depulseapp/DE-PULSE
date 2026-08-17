package main

import (
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"
	"time"
)

func resetV184LoginLimiter(t *testing.T) {
	t.Helper()
	old := loginLimiter
	loginLimiter = newLoginAbuseLimiter(128, 5, 5*time.Minute, 5*time.Minute)
	t.Cleanup(func() { loginLimiter = old })
}

func TestV184ForwardedHeadersFailClosedUntilExplicitlyTrusted(t *testing.T) {
	t.Setenv(runtimeModeEnv, "desktop")
	t.Setenv(hostedTrustProxyHeadersEnv, "true")
	r := httptest.NewRequest(http.MethodGet, "http://local.test/", nil)
	r.RemoteAddr = "10.0.0.2:4000"
	r.Header.Set("X-Forwarded-Proto", "https")
	r.Header.Set("X-Forwarded-For", "203.0.113.9")
	if requestIsSecure(r) {
		t.Fatal("desktop trusted spoofed forwarded proto")
	}
	if got := effectiveClientAddress(r); got != "10.0.0.2" {
		t.Fatalf("desktop trusted spoofed forwarded client: %q", got)
	}

	t.Setenv(runtimeModeEnv, "hosted")
	t.Setenv(hostedTrustProxyHeadersEnv, "")
	if requestIsSecure(r) {
		t.Fatal("hosted mode trusted forwarded proto without explicit proxy trust")
	}
	if got := effectiveClientAddress(r); got != "10.0.0.2" {
		t.Fatalf("hosted mode trusted forwarded client without opt-in: %q", got)
	}

	t.Setenv(hostedTrustProxyHeadersEnv, "true")
	if !requestIsSecure(r) {
		t.Fatal("explicit trusted proxy did not honor https scheme")
	}
	if got := effectiveClientAddress(r); got != "203.0.113.9" {
		t.Fatalf("explicit trusted proxy client mismatch: %q", got)
	}
}

func TestV184SecureCookiesHonorTLSOrExplicitTrustedProxy(t *testing.T) {
	t.Setenv(runtimeModeEnv, "hosted")
	t.Setenv(hostedTrustProxyHeadersEnv, "true")
	r := httptest.NewRequest(http.MethodPost, "http://internal.test/api/auth/login", nil)
	r.Header.Set("X-Forwarded-Proto", "https")
	w := httptest.NewRecorder()
	setSessionCookie(w, r, "token")
	for _, c := range w.Result().Cookies() {
		if (c.Name == sessionCookieName || c.Name == csrfCookieName) && !c.Secure {
			t.Fatalf("trusted https proxy emitted insecure cookie: %+v", c)
		}
	}

	t.Setenv(runtimeModeEnv, "desktop")
	t.Setenv(hostedTrustProxyHeadersEnv, "")
	tlsReq := httptest.NewRequest(http.MethodPost, "https://local.test/api/auth/login", nil)
	tlsRec := httptest.NewRecorder()
	setSessionCookie(tlsRec, tlsReq, "token")
	for _, c := range tlsRec.Result().Cookies() {
		if (c.Name == sessionCookieName || c.Name == csrfCookieName) && !c.Secure {
			t.Fatalf("TLS request emitted insecure cookie: %+v", c)
		}
	}
}

func TestV184SecurityPerimeterOriginAndHeaders(t *testing.T) {
	t.Setenv(runtimeModeEnv, "hosted")
	called := 0
	h := securityPerimeter(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called++
		w.WriteHeader(http.StatusNoContent)
	}))

	foreign := httptest.NewRecorder()
	foreignReq := httptest.NewRequest(http.MethodPost, "https://app.example/api", nil)
	foreignReq.Header.Set("Origin", "https://evil.example")
	h.ServeHTTP(foreign, foreignReq)
	if foreign.Code != http.StatusForbidden || called != 0 {
		t.Fatalf("foreign origin crossed perimeter: code=%d called=%d", foreign.Code, called)
	}

	same := httptest.NewRecorder()
	sameReq := httptest.NewRequest(http.MethodPost, "https://app.example/api", nil)
	sameReq.Header.Set("Origin", "https://app.example")
	h.ServeHTTP(same, sameReq)
	if same.Code != http.StatusNoContent || called != 1 {
		t.Fatalf("same origin blocked: code=%d called=%d", same.Code, called)
	}
	for name := range map[string]bool{
		"Content-Security-Policy": true,
		"Permissions-Policy": true,
		"Cross-Origin-Opener-Policy": true,
		"X-Permitted-Cross-Domain-Policies": true,
		"Strict-Transport-Security": true,
	} {
		if strings.TrimSpace(same.Header().Get(name)) == "" {
			t.Fatalf("missing security header %s", name)
		}
	}

	missing := httptest.NewRecorder()
	missingReq := httptest.NewRequest(http.MethodPost, "https://app.example/api", nil)
	h.ServeHTTP(missing, missingReq)
	if missing.Code != http.StatusNoContent || called != 2 {
		t.Fatalf("native/originless request blocked: code=%d called=%d", missing.Code, called)
	}

	t.Setenv(runtimeModeEnv, "desktop")
	desktop := httptest.NewRecorder()
	h.ServeHTTP(desktop, httptest.NewRequest(http.MethodGet, "https://local.example/", nil))
	if desktop.Header().Get("Strict-Transport-Security") != "" {
		t.Fatal("desktop response emitted hosted HSTS policy")
	}
}

func TestV184ConfiguredPublicOriginControlsHostedBrowserMutations(t *testing.T) {
	t.Setenv(runtimeModeEnv, "hosted")
	t.Setenv(hostedPublicOriginEnv, "https://pulse.example")
	h := securityPerimeter(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { w.WriteHeader(http.StatusNoContent) }))

	allowed := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPost, "http://127.0.0.1:8080/api", nil)
	r.Header.Set("Origin", "https://pulse.example")
	h.ServeHTTP(allowed, r)
	if allowed.Code != http.StatusNoContent {
		t.Fatalf("configured public origin blocked: %d", allowed.Code)
	}

	denied := httptest.NewRecorder()
	r2 := httptest.NewRequest(http.MethodPost, "http://127.0.0.1:8080/api", nil)
	r2.Header.Set("Origin", "http://127.0.0.1:8080")
	h.ServeHTTP(denied, r2)
	if denied.Code != http.StatusForbidden {
		t.Fatalf("internal origin bypassed configured public origin: %d", denied.Code)
	}
}

func TestV184LoginAbuseLimiterIsBoundedAndKeyedWithoutRawIdentity(t *testing.T) {
	resetV184LoginLimiter(t)
	_, s := newIdentityTestService(t)
	app := &Application{identity: s}

	for i := 0; i < 5; i++ {
		rr := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodPost, "/api/auth/login", strings.NewReader(`{"username":"missing-user","password":"wrong"}`))
		r.RemoteAddr = "192.0.2.10:5000"
		app.handleLogin(rr, r)
		if rr.Code != http.StatusUnauthorized || strings.Contains(strings.ToLower(rr.Body.String()), "missing-user") {
			t.Fatalf("failure %d leaked identity or wrong status: %d %s", i+1, rr.Code, rr.Body.String())
		}
	}

	blocked := httptest.NewRecorder()
	blockedReq := httptest.NewRequest(http.MethodPost, "/api/auth/login", strings.NewReader(`{"username":"missing-user","password":"wrong"}`))
	blockedReq.RemoteAddr = "192.0.2.10:5000"
	app.handleLogin(blocked, blockedReq)
	if blocked.Code != http.StatusTooManyRequests {
		t.Fatalf("sixth login was not throttled: %d", blocked.Code)
	}
	if n, err := strconv.Atoi(blocked.Header().Get("Retry-After")); err != nil || n < 1 {
		t.Fatalf("invalid Retry-After: %q", blocked.Header().Get("Retry-After"))
	}

	separate := httptest.NewRecorder()
	separateReq := httptest.NewRequest(http.MethodPost, "/api/auth/login", strings.NewReader(`{"username":"different-user","password":"wrong"}`))
	separateReq.RemoteAddr = "192.0.2.10:5000"
	app.handleLogin(separate, separateReq)
	if separate.Code != http.StatusUnauthorized {
		t.Fatalf("separate login key inherited throttle: %d", separate.Code)
	}

	loginLimiter.mu.Lock()
	defer loginLimiter.mu.Unlock()
	if len(loginLimiter.entries) > loginLimiter.maxEntries {
		t.Fatalf("login limiter exceeded bounded capacity: %d", len(loginLimiter.entries))
	}
	for key := range loginLimiter.entries {
		if strings.Contains(strings.ToLower(key), "missing-user") || strings.Contains(key, "192.0.2.10") {
			t.Fatal("raw login identity retained in limiter key")
		}
	}
}
