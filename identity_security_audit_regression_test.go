package main

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func identitySecurityEventTypes(events []IdentitySecurityEvent) map[IdentitySecurityEventType]bool {
	out := make(map[IdentitySecurityEventType]bool, len(events))
	for _, event := range events {
		out[event.Type] = true
	}
	return out
}

func TestHOST006IdentitySecurityAuditPersistsStaleRetirementWithoutSecrets(t *testing.T) {
	store, s := newIdentityTestService(t)
	base := time.Unix(2_450_000_000, 0)
	token, p, _ := v184CredentialedOwner(t, s, base)
	fingerprint := "sha256:security-audit-fingerprint"
	device, err := s.registerHostedDevice(p, "audited browser", fingerprint)
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
		t.Fatalf("stale audited device remained authorized: %+v", decision)
	}
	if _, err := s.resolve(token, false); err == nil {
		t.Fatal("stale audited device did not revoke its bound session")
	}

	events, err := s.adminSecurityEvents(p)
	if err != nil {
		t.Fatal(err)
	}
	types := identitySecurityEventTypes(events)
	for _, required := range []IdentitySecurityEventType{IdentitySecurityDeviceRegistered, IdentitySecurityDeviceBound, IdentitySecurityDeviceStale, IdentitySecuritySessionRevoked} {
		if !types[required] {
			t.Fatalf("missing security event %s in %+v", required, events)
		}
	}
	encoded, err := json.Marshal(events)
	if err != nil {
		t.Fatal(err)
	}
	text := string(encoded)
	if strings.Contains(text, fingerprint) || strings.Contains(text, token) || strings.Contains(text, "tokenHash") || strings.Contains(text, "fingerprintHash") {
		t.Fatalf("security audit leaked secret/fingerprint material: %s", text)
	}

	reloaded, err := NewIdentityService(store)
	if err != nil {
		t.Fatal(err)
	}
	persisted, err := reloaded.adminSecurityEvents(p)
	if err != nil {
		t.Fatal(err)
	}
	persistedTypes := identitySecurityEventTypes(persisted)
	if !persistedTypes[IdentitySecurityDeviceStale] || !persistedTypes[IdentitySecuritySessionRevoked] {
		t.Fatalf("security audit did not survive identity restart: %+v", persisted)
	}
}

func TestHOST006IdentitySecurityAuditIsBoundedTenantScopedAndPrivileged(t *testing.T) {
	_, s := newIdentityTestService(t)
	base := time.Unix(2_460_000_000, 0)
	_, p, _ := v184CredentialedOwner(t, s, base)

	s.mu.Lock()
	for i := 0; i < maxIdentitySecurityEvents+8; i++ {
		s.appendIdentitySecurityEventLocked(IdentitySecurityDeviceRegistered, localTenantID, p.UserID, randomID("dev"), "", base.Add(time.Duration(i)*time.Second).UnixMilli())
	}
	s.appendIdentitySecurityEventLocked(IdentitySecurityDeviceRegistered, "tenant-other", "other-user", "other-device", "", base.Add(time.Hour).UnixMilli())
	if err := s.persistLocked(); err != nil {
		s.mu.Unlock()
		t.Fatal(err)
	}
	s.mu.Unlock()

	events, err := s.adminSecurityEvents(p)
	if err != nil {
		t.Fatal(err)
	}
	if len(events) > maxIdentitySecurityEvents {
		t.Fatalf("security audit exceeded bound: got=%d max=%d", len(events), maxIdentitySecurityEvents)
	}
	for _, event := range events {
		if normalizedTenantID(event.TenantID) != localTenantID {
			t.Fatalf("cross-tenant security event leaked to owner: %+v", event)
		}
	}
	user := p
	user.Role = RoleUser
	if _, err := s.adminSecurityEvents(user); err == nil {
		t.Fatal("USER role read privileged security audit")
	}
}

func TestHOST006AdminIdentityHTTPProjectsTenantSecurityAudit(t *testing.T) {
	_, s := newIdentityTestService(t)
	base := time.Unix(2_470_000_000, 0)
	_, p, _ := v184CredentialedOwner(t, s, base)
	device, err := s.registerHostedDevice(p, "http audit browser", "sha256:http-audit-device")
	if err != nil {
		t.Fatal(err)
	}
	if err := s.bindHostedDeviceToSession(p, device.ID); err != nil {
		t.Fatal(err)
	}

	app := &Application{identity: s}
	req := httptest.NewRequest(http.MethodGet, "/api/admin/identity", nil)
	req = req.WithContext(context.WithValue(req.Context(), identityContextKey{}, p))
	res := httptest.NewRecorder()
	app.handleAdminIdentity(res, req)
	if res.Code != http.StatusOK {
		t.Fatalf("admin identity audit projection failed: code=%d body=%s", res.Code, res.Body.String())
	}
	var body struct {
		SecurityEvents []IdentitySecurityEvent `json:"securityEvents"`
	}
	if err := json.Unmarshal(res.Body.Bytes(), &body); err != nil {
		t.Fatal(err)
	}
	if len(body.SecurityEvents) == 0 || !identitySecurityEventTypes(body.SecurityEvents)[IdentitySecurityDeviceRegistered] {
		t.Fatalf("admin identity response omitted security audit: %s", res.Body.String())
	}
	if strings.Contains(res.Body.String(), "sha256:http-audit-device") || strings.Contains(res.Body.String(), "fingerprintHash") {
		t.Fatalf("admin identity audit projection leaked device fingerprint: %s", res.Body.String())
	}
}
