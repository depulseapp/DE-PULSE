package main

import (
	"fmt"
	"sort"
	"strings"
)

// hostedTenantIdentityPartitions projects the canonical in-memory identity
// aggregate into tenant-owned persistence records. Application identity remains
// the authorization owner; this function only makes that ownership explicit at
// the hosted persistence boundary.
func hostedTenantIdentityPartitions(state IdentityPersistentState, allowLegacy bool) (map[string]IdentityPersistentState, error) {
	partitions := map[string]IdentityPersistentState{}
	tenantSeen := map[string]struct{}{}
	userTenant := map[string]string{}
	deviceSeen := map[string]struct{}{}
	sessionSeen := map[string]struct{}{}
	eventSeen := map[string]struct{}{}
	productSeen := map[string]struct{}{}

	ensureTenant := func(record TenantRecord) (string, error) {
		tenantID := strings.TrimSpace(record.ID)
		if tenantID == "" {
			return "", fmt.Errorf("hosted tenant persistence: tenant id is required")
		}
		if _, exists := tenantSeen[tenantID]; exists {
			return "", fmt.Errorf("hosted tenant persistence: duplicate tenant %q", tenantID)
		}
		record.ID = tenantID
		tenantSeen[tenantID] = struct{}{}
		partitions[tenantID] = IdentityPersistentState{
			Version:   state.Version,
			Tenants:   []TenantRecord{record},
			UpdatedAt: state.UpdatedAt,
		}
		return tenantID, nil
	}

	for _, tenant := range state.Tenants {
		if _, err := ensureTenant(tenant); err != nil {
			return nil, err
		}
	}

	ensureLegacyLocalTenant := func() {
		if _, exists := tenantSeen[localTenantID]; exists {
			return
		}
		tenantSeen[localTenantID] = struct{}{}
		partitions[localTenantID] = IdentityPersistentState{
			Version: state.Version,
			Tenants: []TenantRecord{{
				ID:        localTenantID,
				Name:      "Local",
				Status:    TenantActive,
				CreatedAt: state.UpdatedAt,
				UpdatedAt: state.UpdatedAt,
			}},
			UpdatedAt: state.UpdatedAt,
		}
	}

	for _, user := range state.Users {
		userID := strings.TrimSpace(user.ID)
		if userID == "" {
			return nil, fmt.Errorf("hosted tenant persistence: user id is required")
		}
		if _, exists := userTenant[userID]; exists {
			return nil, fmt.Errorf("hosted tenant persistence: duplicate user %q", userID)
		}
		tenantID := strings.TrimSpace(user.TenantID)
		if tenantID == "" && allowLegacy {
			tenantID = localTenantID
			ensureLegacyLocalTenant()
		}
		if tenantID == "" {
			return nil, fmt.Errorf("hosted tenant persistence: user %q has no tenant", userID)
		}
		if _, exists := tenantSeen[tenantID]; !exists {
			return nil, fmt.Errorf("hosted tenant persistence: user %q references unknown tenant %q", userID, tenantID)
		}
		user.ID = userID
		user.TenantID = tenantID
		userTenant[userID] = tenantID
		part := partitions[tenantID]
		part.Users = append(part.Users, user)
		partitions[tenantID] = part
	}

	resolveUserTenant := func(userID, explicitTenant, kind, recordID string) (string, error) {
		userID = strings.TrimSpace(userID)
		if userID == "" {
			return "", fmt.Errorf("hosted tenant persistence: %s %q has no user", kind, recordID)
		}
		tenantID, exists := userTenant[userID]
		if !exists {
			return "", fmt.Errorf("hosted tenant persistence: %s %q references unknown user %q", kind, recordID, userID)
		}
		explicitTenant = strings.TrimSpace(explicitTenant)
		if explicitTenant == "" && allowLegacy {
			explicitTenant = tenantID
		}
		if explicitTenant == "" {
			return "", fmt.Errorf("hosted tenant persistence: %s %q has no tenant", kind, recordID)
		}
		if explicitTenant != tenantID {
			return "", fmt.Errorf("hosted tenant persistence: %s %q crosses tenant boundary %q -> %q", kind, recordID, explicitTenant, tenantID)
		}
		return tenantID, nil
	}

	for _, device := range state.Devices {
		deviceID := strings.TrimSpace(device.ID)
		if deviceID == "" {
			return nil, fmt.Errorf("hosted tenant persistence: device id is required")
		}
		if _, exists := deviceSeen[deviceID]; exists {
			return nil, fmt.Errorf("hosted tenant persistence: duplicate device %q", deviceID)
		}
		deviceSeen[deviceID] = struct{}{}
		tenantID, err := resolveUserTenant(device.UserID, device.TenantID, "device", deviceID)
		if err != nil {
			return nil, err
		}
		device.ID = deviceID
		device.UserID = strings.TrimSpace(device.UserID)
		device.TenantID = tenantID
		part := partitions[tenantID]
		part.Devices = append(part.Devices, device)
		partitions[tenantID] = part
	}

	for _, session := range state.Sessions {
		sessionID := strings.TrimSpace(session.ID)
		if sessionID == "" {
			return nil, fmt.Errorf("hosted tenant persistence: session id is required")
		}
		if _, exists := sessionSeen[sessionID]; exists {
			return nil, fmt.Errorf("hosted tenant persistence: duplicate session %q", sessionID)
		}
		sessionSeen[sessionID] = struct{}{}
		tenantID, err := resolveUserTenant(session.UserID, session.TenantID, "session", sessionID)
		if err != nil {
			return nil, err
		}
		session.ID = sessionID
		session.UserID = strings.TrimSpace(session.UserID)
		session.TenantID = tenantID
		part := partitions[tenantID]
		part.Sessions = append(part.Sessions, session)
		partitions[tenantID] = part
	}

	for _, event := range state.SecurityEvents {
		eventID := strings.TrimSpace(event.ID)
		if eventID == "" {
			return nil, fmt.Errorf("hosted tenant persistence: security event id is required")
		}
		if _, exists := eventSeen[eventID]; exists {
			return nil, fmt.Errorf("hosted tenant persistence: duplicate security event %q", eventID)
		}
		eventSeen[eventID] = struct{}{}
		tenantID := strings.TrimSpace(event.TenantID)
		if strings.TrimSpace(event.UserID) != "" {
			resolved, err := resolveUserTenant(event.UserID, tenantID, "security event", eventID)
			if err != nil {
				return nil, err
			}
			tenantID = resolved
		} else {
			if tenantID == "" {
				return nil, fmt.Errorf("hosted tenant persistence: security event %q has no tenant", eventID)
			}
			if _, exists := tenantSeen[tenantID]; !exists {
				return nil, fmt.Errorf("hosted tenant persistence: security event %q references unknown tenant %q", eventID, tenantID)
			}
		}
		event.ID = eventID
		event.TenantID = tenantID
		event.UserID = strings.TrimSpace(event.UserID)
		part := partitions[tenantID]
		part.SecurityEvents = append(part.SecurityEvents, event)
		partitions[tenantID] = part
	}

	for _, product := range state.ProductEntitlements {
		tenantID := strings.TrimSpace(product.TenantID)
		if tenantID == "" {
			return nil, fmt.Errorf("hosted tenant persistence: product entitlement has no tenant")
		}
		if _, exists := tenantSeen[tenantID]; !exists {
			return nil, fmt.Errorf("hosted tenant persistence: product entitlement references unknown tenant %q", tenantID)
		}
		if _, exists := productSeen[tenantID]; exists {
			return nil, fmt.Errorf("hosted tenant persistence: duplicate product entitlement for tenant %q", tenantID)
		}
		productSeen[tenantID] = struct{}{}
		product.TenantID = tenantID
		part := partitions[tenantID]
		part.ProductEntitlements = append(part.ProductEntitlements, product)
		partitions[tenantID] = part
	}

	return partitions, nil
}

func hostedTenantIdentityFromPartitions(partitions map[string]IdentityPersistentState) (IdentityPersistentState, error) {
	if len(partitions) == 0 {
		return IdentityPersistentState{}, nil
	}
	keys := make([]string, 0, len(partitions))
	for tenantID := range partitions {
		keys = append(keys, tenantID)
	}
	sort.Strings(keys)
	merged := IdentityPersistentState{}
	for _, tenantID := range keys {
		part := partitions[tenantID]
		validated, err := hostedTenantIdentityPartitions(part, false)
		if err != nil {
			return IdentityPersistentState{}, fmt.Errorf("tenant %q payload invalid: %w", tenantID, err)
		}
		if len(validated) != 1 {
			return IdentityPersistentState{}, fmt.Errorf("tenant %q payload contains %d tenant partitions", tenantID, len(validated))
		}
		if _, ok := validated[tenantID]; !ok {
			return IdentityPersistentState{}, fmt.Errorf("tenant row %q payload belongs to another tenant", tenantID)
		}
		if part.Version > merged.Version {
			merged.Version = part.Version
		}
		if part.UpdatedAt > merged.UpdatedAt {
			merged.UpdatedAt = part.UpdatedAt
		}
		merged.Tenants = append(merged.Tenants, part.Tenants...)
		merged.Users = append(merged.Users, part.Users...)
		merged.Devices = append(merged.Devices, part.Devices...)
		merged.Sessions = append(merged.Sessions, part.Sessions...)
		merged.SecurityEvents = append(merged.SecurityEvents, part.SecurityEvents...)
		merged.ProductEntitlements = append(merged.ProductEntitlements, part.ProductEntitlements...)
	}
	// Repartitioning the aggregate is a final cross-row uniqueness and
	// ownership check (duplicate users/devices/sessions cannot hide in separate
	// tenant rows).
	if _, err := hostedTenantIdentityPartitions(merged, false); err != nil {
		return IdentityPersistentState{}, err
	}
	return merged, nil
}

func hostedTenantUserIndex(state IdentityPersistentState) (map[string]string, error) {
	if _, err := hostedTenantIdentityPartitions(state, false); err != nil {
		return nil, err
	}
	out := make(map[string]string, len(state.Users))
	for _, user := range state.Users {
		userID := strings.TrimSpace(user.ID)
		tenantID := strings.TrimSpace(user.TenantID)
		if userID == "" || tenantID == "" {
			return nil, fmt.Errorf("hosted tenant persistence: incomplete user ownership")
		}
		out[userID] = tenantID
	}
	return out, nil
}
