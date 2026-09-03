package main

import (
	"strings"
	"testing"
)

func host015TenantFixture() IdentityPersistentState {
	return IdentityPersistentState{
		Version: 1,
		Tenants: []TenantRecord{
			{ID: "tenant-a", Name: "A", Status: TenantActive, CreatedAt: 1, UpdatedAt: 10},
			{ID: "tenant-b", Name: "B", Status: TenantActive, CreatedAt: 2, UpdatedAt: 10},
		},
		Users: []UserRecord{
			{ID: "user-a", TenantID: "tenant-a", Username: "a", Role: RoleUser, Status: UserActive, CreatedAt: 1, UpdatedAt: 10},
			{ID: "user-b", TenantID: "tenant-b", Username: "b", Role: RoleUser, Status: UserActive, CreatedAt: 2, UpdatedAt: 10},
		},
		Devices: []DeviceRecord{
			{ID: "device-a", TenantID: "tenant-a", UserID: "user-a", Status: DeviceActive, CreatedAt: 3, LastSeenAt: 10},
			{ID: "device-b", TenantID: "tenant-b", UserID: "user-b", Status: DeviceActive, CreatedAt: 4, LastSeenAt: 10},
		},
		Sessions: []SessionRecord{
			{ID: "session-a", TenantID: "tenant-a", UserID: "user-a", CreatedAt: 3, LastSeenAt: 10, IdleExpiresAt: 100, AbsoluteExpiresAt: 200},
			{ID: "session-b", TenantID: "tenant-b", UserID: "user-b", CreatedAt: 4, LastSeenAt: 10, IdleExpiresAt: 100, AbsoluteExpiresAt: 200},
		},
		SecurityEvents: []IdentitySecurityEvent{
			{ID: "event-a", TenantID: "tenant-a", UserID: "user-a", Type: IdentitySecuritySessionRevoked, CreatedAt: 8},
			{ID: "event-b", TenantID: "tenant-b", UserID: "user-b", Type: IdentitySecuritySessionRevoked, CreatedAt: 9},
		},
		ProductEntitlements: []TenantProductEntitlement{
			{TenantID: "tenant-a", Plan: ProductPlanFoundationCore, Status: ProductStatusActive, WindowStartedAt: 1, UpdatedAt: 10},
			{TenantID: "tenant-b", Plan: ProductPlanFoundationCore, Status: ProductStatusActive, WindowStartedAt: 1, UpdatedAt: 10},
		},
		UpdatedAt: 10,
	}
}

func TestHOST015TenantPersistencePartitionsEveryIdentityRecord(t *testing.T) {
	state := host015TenantFixture()
	parts, err := hostedTenantIdentityPartitions(state, false)
	if err != nil {
		t.Fatal(err)
	}
	if len(parts) != 2 {
		t.Fatalf("tenant partitions=%d want=2", len(parts))
	}
	for tenantID, expectedUser := range map[string]string{"tenant-a": "user-a", "tenant-b": "user-b"} {
		part := parts[tenantID]
		if len(part.Tenants) != 1 || part.Tenants[0].ID != tenantID {
			t.Fatalf("tenant row not isolated: tenant=%s payload=%+v", tenantID, part.Tenants)
		}
		if len(part.Users) != 1 || part.Users[0].ID != expectedUser || part.Users[0].TenantID != tenantID {
			t.Fatalf("user row not isolated: tenant=%s users=%+v", tenantID, part.Users)
		}
		if len(part.Devices) != 1 || part.Devices[0].TenantID != tenantID || len(part.Sessions) != 1 || part.Sessions[0].TenantID != tenantID {
			t.Fatalf("device/session row crossed tenant boundary: tenant=%s devices=%+v sessions=%+v", tenantID, part.Devices, part.Sessions)
		}
		if len(part.SecurityEvents) != 1 || part.SecurityEvents[0].TenantID != tenantID || len(part.ProductEntitlements) != 1 || part.ProductEntitlements[0].TenantID != tenantID {
			t.Fatalf("audit/product row crossed tenant boundary: tenant=%s events=%+v products=%+v", tenantID, part.SecurityEvents, part.ProductEntitlements)
		}
	}

	merged, err := hostedTenantIdentityFromPartitions(parts)
	if err != nil {
		t.Fatal(err)
	}
	if len(merged.Tenants) != 2 || len(merged.Users) != 2 || len(merged.Devices) != 2 || len(merged.Sessions) != 2 || len(merged.SecurityEvents) != 2 || len(merged.ProductEntitlements) != 2 {
		t.Fatalf("merged tenant state lost records: %+v", merged)
	}
	index, err := hostedTenantUserIndex(merged)
	if err != nil {
		t.Fatal(err)
	}
	if index["user-a"] != "tenant-a" || index["user-b"] != "tenant-b" {
		t.Fatalf("tenant user index incorrect: %+v", index)
	}
}

func TestHOST015TenantPersistenceRejectsCrossTenantSession(t *testing.T) {
	state := host015TenantFixture()
	state.Sessions[0].TenantID = "tenant-b"
	_, err := hostedTenantIdentityPartitions(state, false)
	if err == nil || !strings.Contains(err.Error(), "crosses tenant boundary") {
		t.Fatalf("expected cross-tenant session denial, got %v", err)
	}
}

func TestHOST015TenantPersistenceRejectsDuplicateUserAcrossTenants(t *testing.T) {
	state := host015TenantFixture()
	state.Users[1].ID = "user-a"
	_, err := hostedTenantIdentityPartitions(state, false)
	if err == nil || !strings.Contains(err.Error(), "duplicate user") {
		t.Fatalf("expected duplicate user denial, got %v", err)
	}
}

func TestHOST015TenantPersistenceLegacyMigrationIsExplicitlyLocalOnly(t *testing.T) {
	legacy := IdentityPersistentState{
		Version:   1,
		Users:     []UserRecord{{ID: "legacy-user", Username: "legacy", Role: RoleOwner, Status: UserActive, CreatedAt: 1, UpdatedAt: 2}},
		Sessions:  []SessionRecord{{ID: "legacy-session", UserID: "legacy-user", CreatedAt: 1, LastSeenAt: 2, IdleExpiresAt: 100, AbsoluteExpiresAt: 200}},
		UpdatedAt: 2,
	}
	parts, err := hostedTenantIdentityPartitions(legacy, true)
	if err != nil {
		t.Fatal(err)
	}
	part, ok := parts[localTenantID]
	if !ok || len(parts) != 1 || len(part.Tenants) != 1 || part.Tenants[0].ID != localTenantID {
		t.Fatalf("legacy migration did not bind only to canonical local tenant: %+v", parts)
	}
	if part.Users[0].TenantID != localTenantID || part.Sessions[0].TenantID != localTenantID {
		t.Fatalf("legacy dependent ownership not normalized: users=%+v sessions=%+v", part.Users, part.Sessions)
	}
	if _, err := hostedTenantIdentityPartitions(legacy, false); err == nil {
		t.Fatal("strict tenant persistence unexpectedly accepted unowned legacy identity")
	}
}

func TestHOST015TenantMergeRejectsTamperedTenantPayload(t *testing.T) {
	parts, err := hostedTenantIdentityPartitions(host015TenantFixture(), false)
	if err != nil {
		t.Fatal(err)
	}
	tampered := parts["tenant-a"]
	tampered.Users[0].TenantID = "tenant-b"
	parts["tenant-a"] = tampered
	if _, err := hostedTenantIdentityFromPartitions(parts); err == nil {
		t.Fatal("tampered tenant payload unexpectedly merged")
	}
}
