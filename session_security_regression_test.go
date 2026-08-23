package main

import (
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
