package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func v186DocumentationSessionForRole(t *testing.T, s *IdentityService, role UserRole) string {
	t.Helper()
	now := time.Now().UnixMilli()
	id := "v186_docs_" + strings.ToLower(string(role))
	user := UserRecord{
		ID:           id,
		Username:     id,
		DisplayName:  "v18.6 docs " + string(role),
		PasswordHash: "test-only-present",
		Role:         role,
		Status:       UserActive,
		CreatedAt:    now,
		UpdatedAt:    now,
	}
	s.mu.Lock()
	s.state.Users = append(s.state.Users, user)
	if err := s.persistLocked(); err != nil {
		s.mu.Unlock()
		t.Fatal(err)
	}
	token, _, err := s.createSessionLocked(user, "")
	s.mu.Unlock()
	if err != nil {
		t.Fatal(err)
	}
	return token
}

func TestV186DocumentationAudiencePathPolicy(t *testing.T) {
	cases := map[string]documentationAudience{
		"/docs/user.md":               documentationAudienceAuthenticated,
		"/docs/limitations.md":        documentationAudienceAuthenticated,
		"/docs/developer.md":          documentationAudienceDeveloper,
		"/docs/./developer.md":        documentationAudienceDeveloper,
		"/docs/unknown.md":            documentationAudienceAuthenticated,
		"/api/bootstrap":              documentationAudienceNone,
		"/renderer/docs/developer.md": documentationAudienceNone,
	}
	for raw, want := range cases {
		if got := documentationAudienceForPath(raw); got != want {
			t.Fatalf("path %q audience=%v want=%v", raw, got, want)
		}
	}
}

func TestV186DocumentationDirectPathRoleEnforcement(t *testing.T) {
	_, s := newIdentityTestService(t)
	app := &Application{identity: s}
	called := 0
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called++
		w.WriteHeader(http.StatusNoContent)
	})
	h := app.protectDocumentationHTTP(next)

	request := func(path, token string) int {
		rr := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, path, nil)
		if token != "" {
			req.AddCookie(&http.Cookie{Name: sessionCookieName, Value: token})
		}
		h.ServeHTTP(rr, req)
		return rr.Code
	}

	if code := request("/docs/user.md", ""); code != http.StatusUnauthorized {
		t.Fatalf("anonymous documentation access=%d want=%d", code, http.StatusUnauthorized)
	}

	roles := []UserRole{RoleDemo, RoleUser, RoleAdmin, RoleOwner, RoleSuperOwner}
	for _, role := range roles {
		token := v186DocumentationSessionForRole(t, s, role)
		if code := request("/docs/user.md", token); code != http.StatusNoContent {
			t.Fatalf("%s user docs=%d", role, code)
		}
		if code := request("/docs/limitations.md", token); code != http.StatusNoContent {
			t.Fatalf("%s limitations docs=%d", role, code)
		}
		wantDeveloper := http.StatusForbidden
		if roleAtLeast(role, RoleAdmin) {
			wantDeveloper = http.StatusNoContent
		}
		if code := request("/docs/developer.md", token); code != wantDeveloper {
			t.Fatalf("%s developer docs=%d want=%d", role, code, wantDeveloper)
		}
	}

	if called == 0 {
		t.Fatal("authorized documentation requests never reached the static handler")
	}
}

func TestV186DocumentationLocalIdentityCompatibility(t *testing.T) {
	app := &Application{}
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { w.WriteHeader(http.StatusTeapot) })
	rr := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/docs/developer.md", nil)
	app.protectDocumentationHTTP(next).ServeHTTP(rr, req)
	if rr.Code != http.StatusTeapot {
		t.Fatalf("identity-disabled local docs changed: %s", fmt.Sprint(rr.Code))
	}
}
