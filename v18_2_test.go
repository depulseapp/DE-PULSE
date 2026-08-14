package main

import (
	"encoding/json"
	"strings"
	"testing"
	"time"
)

func TestV182AdminRoleHierarchyAndLifecycle(t *testing.T) {
	_, s := newIdentityTestService(t)
	ownerToken, owner, err := s.bootstrapOwnerSession()
	if err != nil {
		t.Fatal(err)
	}
	if owner.Role != RoleOwner {
		t.Fatalf("expected OWNER, got %s", owner.Role)
	}

	admin, err := s.adminCreateUser(owner, "ops.admin", "Ops Admin", RoleAdmin, "temporary admin password")
	if err != nil {
		t.Fatal(err)
	}
	user, err := s.adminCreateUser(Principal{UserID: admin.ID, Username: admin.Username, Role: RoleAdmin}, "analyst.one", "Analyst One", RoleUser, "temporary user password")
	if err != nil {
		t.Fatal(err)
	}
	if _, err := s.adminCreateUser(Principal{UserID: admin.ID, Username: admin.Username, Role: RoleAdmin}, "bad.owner", "Bad Owner", RoleOwner, "temporary owner password"); err == nil {
		t.Fatal("ADMIN created OWNER")
	}
	if err := s.adminSetUserStatus(Principal{UserID: admin.ID, Username: admin.Username, Role: RoleAdmin}, owner.UserID, UserDisabled); err == nil {
		t.Fatal("ADMIN disabled OWNER")
	}
	if err := s.adminSetUserRole(Principal{UserID: admin.ID, Username: admin.Username, Role: RoleAdmin}, owner.UserID, RoleUser); err == nil {
		t.Fatal("ADMIN changed OWNER role")
	}
	if err := s.adminSetUserRole(owner, user.ID, RoleAdmin); err != nil {
		t.Fatalf("OWNER could not promote USER to ADMIN: %v", err)
	}
	if err := s.adminSetUserStatus(owner, admin.ID, UserDisabled); err != nil {
		t.Fatalf("OWNER could not disable ADMIN: %v", err)
	}
	if _, err := s.resolve(ownerToken, false); err != nil {
		t.Fatalf("owner session changed by lower-user lifecycle: %v", err)
	}
}

func TestV182PresenceUsesCanonicalSessionAndSSETouch(t *testing.T) {
	_, s := newIdentityTestService(t)
	base := time.Unix(2_100_000_000, 0)
	s.now = func() time.Time { return base }
	_, owner, err := s.bootstrapOwnerSession()
	if err != nil {
		t.Fatal(err)
	}
	user, err := s.adminCreateUser(owner, "presence.user", "Presence User", RoleUser, "presence temporary password")
	if err != nil {
		t.Fatal(err)
	}
	_, up, err := s.authenticate(user.Username, "presence temporary password")
	if err != nil {
		t.Fatal(err)
	}
	users, sessions, err := s.adminSnapshot(owner)
	if err != nil {
		t.Fatal(err)
	}
	found := false
	for _, u := range users {
		if u.ID == user.ID {
			found = true
			if u.Presence != "ACTIVE" || u.ActiveSessionCount != 1 {
				t.Fatalf("unexpected initial presence: %+v", u)
			}
		}
	}
	if !found {
		t.Fatal("user missing from admin snapshot")
	}
	for _, rec := range sessions {
		if rec.UserID == user.ID && (rec.Presence != "ACTIVE" || rec.ID != up.SessionID) {
			t.Fatalf("unexpected session view: %+v", rec)
		}
	}

	s.now = func() time.Time { return base.Add(2 * time.Minute) }
	users, _, _ = s.adminSnapshot(owner)
	for _, u := range users {
		if u.ID == user.ID && u.Presence != "IDLE" {
			t.Fatalf("expected IDLE after inactivity: %+v", u)
		}
	}
	s.touchSessionID(up.SessionID)
	users, _, _ = s.adminSnapshot(owner)
	for _, u := range users {
		if u.ID == user.ID && u.Presence != "ACTIVE" {
			t.Fatalf("SSE/session touch did not restore ACTIVE: %+v", u)
		}
	}
}

func TestV182PasswordResetDisableAndSessionRevocation(t *testing.T) {
	_, s := newIdentityTestService(t)
	_, owner, err := s.bootstrapOwnerSession()
	if err != nil {
		t.Fatal(err)
	}
	user, err := s.adminCreateUser(owner, "lifecycle.user", "Lifecycle User", RoleUser, "initial temporary password")
	if err != nil {
		t.Fatal(err)
	}
	token, principal, err := s.authenticate(user.Username, "initial temporary password")
	if err != nil {
		t.Fatal(err)
	}
	if err := s.adminResetPassword(owner, user.ID, "replacement temporary password"); err != nil {
		t.Fatal(err)
	}
	if _, err := s.resolve(token, false); err == nil {
		t.Fatal("password reset did not revoke prior session")
	}
	newToken, newPrincipal, err := s.authenticate(user.Username, "replacement temporary password")
	if err != nil {
		t.Fatal(err)
	}
	if !s.userRequiresPassword(user.ID) {
		t.Fatal("admin reset must require password change")
	}
	if err := s.adminRevokeSession(owner, newPrincipal.SessionID); err != nil {
		t.Fatal(err)
	}
	if _, err := s.resolve(newToken, false); err == nil {
		t.Fatal("admin session revoke did not invalidate session")
	}

	token2, _, err := s.authenticate(user.Username, "replacement temporary password")
	if err != nil {
		t.Fatal(err)
	}
	if err := s.adminSetUserStatus(owner, user.ID, UserDisabled); err != nil {
		t.Fatal(err)
	}
	if _, err := s.resolve(token2, false); err == nil {
		t.Fatal("disabled user session remained valid")
	}
	if _, _, err := s.authenticate(user.Username, "replacement temporary password"); err == nil {
		t.Fatal("disabled user could authenticate")
	}
	_ = principal
}

func TestV182AdminViewsRedactCredentialMaterial(t *testing.T) {
	_, s := newIdentityTestService(t)
	_, owner, err := s.bootstrapOwnerSession()
	if err != nil {
		t.Fatal(err)
	}
	user, err := s.adminCreateUser(owner, "redaction.user", "Redaction User", RoleUser, "redaction temporary password")
	if err != nil {
		t.Fatal(err)
	}
	_, _, err = s.authenticate(user.Username, "redaction temporary password")
	if err != nil {
		t.Fatal(err)
	}
	users, sessions, err := s.adminSnapshot(owner)
	if err != nil {
		t.Fatal(err)
	}
	raw, err := json.Marshal(map[string]any{"users": users, "sessions": sessions})
	if err != nil {
		t.Fatal(err)
	}
	text := strings.ToLower(string(raw))
	for _, forbidden := range []string{"passwordhash", "tokenhash", "argon2id", "redaction temporary password"} {
		if strings.Contains(text, forbidden) {
			t.Fatalf("admin view leaked credential material %q", forbidden)
		}
	}
}

func TestV182CannotRevokeOwnCurrentSessionThroughAdminPath(t *testing.T) {
	_, s := newIdentityTestService(t)
	_, owner, err := s.bootstrapOwnerSession()
	if err != nil {
		t.Fatal(err)
	}
	if err := s.adminRevokeSession(owner, owner.SessionID); err == nil {
		t.Fatal("admin path revoked actor current session")
	}
}
