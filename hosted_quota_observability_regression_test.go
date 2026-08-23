package main

import (
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"
	"time"
)

func resetV184HostedQuotaLimiter(t *testing.T, maxEntries, mutationLimit, expensiveLimit int) {
	t.Helper()
	old := hostedQuotaLimiter
	hostedQuotaLimiter = newHostedRequestQuotaLimiter(maxEntries, mutationLimit, expensiveLimit, time.Minute)
	t.Cleanup(func() { hostedQuotaLimiter = old })
}

func TestV184HostedQuotaAppliesOnlyToHostedAuthenticatedMutations(t *testing.T) {
	t.Setenv(runtimeModeEnv, "desktop")
	mutation := httptest.NewRequest(http.MethodPost, "/api/watchlists/create", nil)
	if hostedQuotaApplies(mutation) {
		t.Fatal("desktop mutation entered hosted quota")
	}

	t.Setenv(runtimeModeEnv, "hosted")
	if !hostedQuotaApplies(mutation) {
		t.Fatal("hosted authenticated mutation bypassed quota")
	}
	if hostedQuotaApplies(httptest.NewRequest(http.MethodGet, "/api/bootstrap", nil)) {
		t.Fatal("read-only hosted request entered mutation quota")
	}
	if hostedQuotaApplies(httptest.NewRequest(http.MethodPost, "/api/auth/reauth", nil)) {
		t.Fatal("authentication endpoint entered hosted request quota")
	}
}

func TestV184HostedQuotaWindowsAreDeterministicAndBounded(t *testing.T) {
	limiter := newHostedRequestQuotaLimiter(2, 2, 1, time.Minute)
	now := time.Unix(2_300_000_000, 0)
	limiter.now = func() time.Time { return now }

	if ok, _ := limiter.Allow("user-a", true); !ok {
		t.Fatal("first expensive mutation was blocked")
	}
	if ok, retry := limiter.Allow("user-a", true); ok || retry != 60 {
		t.Fatalf("expensive quota mismatch: ok=%v retry=%d", ok, retry)
	}
	if ok, _ := limiter.Allow("user-a", false); !ok {
		t.Fatal("remaining broad mutation budget was not available")
	}
	if ok, retry := limiter.Allow("user-a", false); ok || retry != 60 {
		t.Fatalf("mutation quota mismatch: ok=%v retry=%d", ok, retry)
	}

	now = now.Add(61 * time.Second)
	if ok, retry := limiter.Allow("user-a", true); !ok || retry != 0 {
		t.Fatalf("quota window did not reset: ok=%v retry=%d", ok, retry)
	}

	if ok, _ := limiter.Allow("user-b", false); !ok {
		t.Fatal("second user unexpectedly blocked")
	}
	if ok, _ := limiter.Allow("user-c", false); !ok {
		t.Fatal("bounded eviction rejected a new user")
	}
	limiter.mu.Lock()
	defer limiter.mu.Unlock()
	if len(limiter.entries) > limiter.maxEntries {
		t.Fatalf("hosted quota limiter exceeded bounded capacity: %d", len(limiter.entries))
	}
}

func TestV184HostedQuotaUsesHashedVerifiedIdentityAndAggregateTelemetry(t *testing.T) {
	t.Setenv(runtimeModeEnv, "hosted")
	resetV184HostedQuotaLimiter(t, 8, 3, 1)
	app := &Application{httpTelemetry: NewRequestTelemetry()}
	principal := Principal{UserID: "verified-user-a", Username: "user-a"}

	first := httptest.NewRecorder()
	firstReq := httptest.NewRequest(http.MethodPost, "/api/cache/refresh", nil)
	if !app.enforceHostedRequestQuota(first, firstReq, principal) {
		t.Fatalf("first expensive request blocked: %d %s", first.Code, first.Body.String())
	}

	blocked := httptest.NewRecorder()
	blockedReq := httptest.NewRequest(http.MethodPost, "/api/cache/refresh", nil)
	if app.enforceHostedRequestQuota(blocked, blockedReq, principal) {
		t.Fatal("second expensive request exceeded quota but was allowed")
	}
	if blocked.Code != http.StatusTooManyRequests {
		t.Fatalf("wrong throttle status: %d", blocked.Code)
	}
	if n, err := strconv.Atoi(blocked.Header().Get("Retry-After")); err != nil || n < 1 {
		t.Fatalf("invalid Retry-After: %q", blocked.Header().Get("Retry-After"))
	}

	other := httptest.NewRecorder()
	otherReq := httptest.NewRequest(http.MethodPost, "/api/cache/refresh", nil)
	if !app.enforceHostedRequestQuota(other, otherReq, Principal{UserID: "verified-user-b", Username: "user-b"}) {
		t.Fatal("separate verified user inherited another user's quota")
	}

	diag := app.httpTelemetry.Diagnostics()
	if diag.HostedMutationAllowed != 2 || diag.HostedMutationThrottled != 1 || diag.HostedExpensiveAllowed != 2 || diag.HostedExpensiveThrottled != 1 {
		t.Fatalf("aggregate hosted quota telemetry mismatch: %+v", diag)
	}

	hostedQuotaLimiter.mu.Lock()
	defer hostedQuotaLimiter.mu.Unlock()
	for key := range hostedQuotaLimiter.entries {
		if strings.Contains(key, principal.UserID) || strings.Contains(key, principal.Username) {
			t.Fatal("raw identity retained in hosted quota key")
		}
	}
}

func TestV184HostedQuotaRunsAfterCanonicalAuthAndCSRF(t *testing.T) {
	t.Setenv(runtimeModeEnv, "hosted")
	resetV184HostedQuotaLimiter(t, 8, 1, 1)
	persistence, identity := newIdentityTestService(t)
	_, bootstrap, err := identity.bootstrapOwnerSession()
	if err != nil {
		t.Fatal(err)
	}
	token, principal, err := identity.setPassword(bootstrap.UserID, "v18.4 hosted quota passphrase")
	if err != nil {
		t.Fatal(err)
	}
	app := &Application{
		identity:      identity,
		persistence:   persistence,
		state:         defaultState(),
		workspaces:    map[string]UserWorkspace{principal.UserID: defaultUserWorkspace(principal.UserID)},
		httpTelemetry: NewRequestTelemetry(),
	}
	called := 0
	handler := app.auth(func(w http.ResponseWriter, r *http.Request) {
		called++
		w.WriteHeader(http.StatusNoContent)
	})
	request := func(sessionToken, csrf string) *http.Request {
		r := httptest.NewRequest(http.MethodPost, "/api/watchlists/create", nil)
		if sessionToken != "" {
			r.AddCookie(&http.Cookie{Name: sessionCookieName, Value: sessionToken})
		}
		if csrf != "" {
			r.AddCookie(&http.Cookie{Name: csrfCookieName, Value: csrf})
			r.Header.Set("X-DE-PULSE-CSRF", csrf)
		}
		return r
	}

	unauthorized := httptest.NewRecorder()
	handler.ServeHTTP(unauthorized, request("fake-session", "csrf"))
	if unauthorized.Code != http.StatusUnauthorized || called != 0 {
		t.Fatalf("unverified session reached quota-protected handler: code=%d called=%d", unauthorized.Code, called)
	}
	if got := app.httpTelemetry.Diagnostics().HostedMutationAllowed; got != 0 {
		t.Fatalf("unverified session affected hosted quota telemetry: %d", got)
	}

	first := httptest.NewRecorder()
	handler.ServeHTTP(first, request(token, "csrf-valid"))
	if first.Code != http.StatusNoContent || called != 1 {
		t.Fatalf("authenticated CSRF-valid request blocked: code=%d called=%d", first.Code, called)
	}

	second := httptest.NewRecorder()
	handler.ServeHTTP(second, request(token, "csrf-valid"))
	if second.Code != http.StatusTooManyRequests || called != 1 {
		t.Fatalf("quota did not stop second authenticated mutation: code=%d called=%d", second.Code, called)
	}
}
