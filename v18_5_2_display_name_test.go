package main

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestV1852ConfigurableDisplayNamePersistsWithoutChangingRoleOrSession(t *testing.T) {
	persistence, identity := newIdentityTestService(t)
	token, principal, err := identity.bootstrapOwnerSession()
	if err != nil {
		t.Fatal(err)
	}
	originalSession := principal.SessionID
	if err := identity.updateDisplayName(principal.UserID, "  Deivaram   Venkatachalapathy  "); err != nil {
		t.Fatal(err)
	}
	resolved, err := identity.resolve(token, false)
	if err != nil {
		t.Fatalf("display-name update invalidated session: %v", err)
	}
	if resolved.DisplayName != "Deivaram Venkatachalapathy" {
		t.Fatalf("display name was not normalized: %+v", resolved)
	}
	if resolved.Role != RoleOwner || resolved.Username != "owner" || resolved.SessionID != originalSession {
		t.Fatalf("display-name update changed identity authority or session: %+v", resolved)
	}

	reloaded, err := NewIdentityService(persistence)
	if err != nil {
		t.Fatal(err)
	}
	persisted, err := reloaded.resolve(token, false)
	if err != nil {
		t.Fatalf("persisted session unavailable after identity reload: %v", err)
	}
	if persisted.DisplayName != "Deivaram Venkatachalapathy" || persisted.Role != RoleOwner {
		t.Fatalf("display name or role did not persist correctly: %+v", persisted)
	}

	app := &Application{identity: reloaded}
	body := strings.NewReader(`{"displayName":"DV Market Desk"}`)
	req := httptest.NewRequest(http.MethodPost, "/api/auth/profile", body)
	req.AddCookie(&http.Cookie{Name: sessionCookieName, Value: token})
	req = req.WithContext(context.WithValue(req.Context(), identityContextKey{}, persisted))
	rr := httptest.NewRecorder()
	app.handleUpdateProfile(rr, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("profile endpoint failed: %d %s", rr.Code, rr.Body.String())
	}
	var response struct {
		Principal Principal `json:"principal"`
	}
	if err := json.Unmarshal(rr.Body.Bytes(), &response); err != nil {
		t.Fatal(err)
	}
	if response.Principal.DisplayName != "DV Market Desk" || response.Principal.Role != RoleOwner || response.Principal.SessionID != originalSession {
		t.Fatalf("profile response changed role/session or omitted name: %+v", response.Principal)
	}

	if err := reloaded.updateDisplayName(principal.UserID, "   "); err == nil {
		t.Fatal("blank display name accepted")
	}
	still, err := reloaded.resolve(token, false)
	if err != nil {
		t.Fatal(err)
	}
	if still.DisplayName != "DV Market Desk" {
		t.Fatalf("invalid update changed the saved display name: %+v", still)
	}
}
