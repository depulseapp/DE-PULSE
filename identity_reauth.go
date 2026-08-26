package main

import (
	"errors"
	"strings"
	"time"
)

const (
	hostedCapabilityTenantManage  = "TENANT_MANAGE"
	hostedCapabilityAccountManage = "ACCOUNT_MANAGE"
	hostedCapabilityUserManage    = "USER_MANAGE"
	hostedCapabilitySessionManage = "SESSION_MANAGE"
	hostedCapabilityDeviceManage  = "DEVICE_MANAGE"
	hostedCapabilityStandardUse   = "STANDARD_USE"
	hostedCapabilityDemoUse       = "DEMO_USE"
	defaultHostedDeviceStaleTTL   = 30 * 24 * time.Hour
)

type HostedIdentityRequirement struct {
	TenantID                    string `json:"tenantId"`
	Capability                  string `json:"capability"`
	RequireRegisteredDevice     bool   `json:"requireRegisteredDevice"`
	RequireRecentAuthentication bool   `json:"requireRecentAuthentication"`
	RequireMFA                  bool   `json:"requireMfa"`
}

type HostedIdentityDecision struct {
	Allowed         bool     `json:"allowed"`
	TenantID        string   `json:"tenantId"`
	UserID          string   `json:"userId"`
	SessionID       string   `json:"sessionId"`
	DeviceID        string   `json:"deviceId,omitempty"`
	Capability      string   `json:"capability"`
	BlockingReasons []string `json:"blockingReasons,omitempty"`
}

func roleHasHostedCapability(role UserRole, capability string) bool {
	capability = strings.ToUpper(strings.TrimSpace(capability))
	switch role {
	case RoleSuperOwner:
		return capability == hostedCapabilityTenantManage || capability == hostedCapabilityAccountManage || capability == hostedCapabilityUserManage || capability == hostedCapabilitySessionManage || capability == hostedCapabilityDeviceManage || capability == hostedCapabilityStandardUse || capability == hostedCapabilityDemoUse
	case RoleOwner:
		return capability == hostedCapabilityTenantManage || capability == hostedCapabilityAccountManage || capability == hostedCapabilityUserManage || capability == hostedCapabilitySessionManage || capability == hostedCapabilityDeviceManage || capability == hostedCapabilityStandardUse
	case RoleAdmin:
		return capability == hostedCapabilityUserManage || capability == hostedCapabilitySessionManage || capability == hostedCapabilityDeviceManage || capability == hostedCapabilityStandardUse
	case RoleUser:
		return capability == hostedCapabilitySessionManage || capability == hostedCapabilityDeviceManage || capability == hostedCapabilityStandardUse
	case RoleDemo:
		return capability == hostedCapabilityDemoUse
	default:
		return false
	}
}

func (s *IdentityService) reauthenticateSession(sessionID, password string) (bool, error) {
	sessionID = strings.TrimSpace(sessionID)
	if sessionID == "" || password == "" {
		return false, nil
	}

	now := s.now().UnixMilli()
	s.mu.Lock()
	userID := ""
	passwordHash := ""
	for i := range s.state.Sessions {
		rec := s.state.Sessions[i]
		if rec.ID == sessionID && rec.RevokedAt == 0 && now < rec.IdleExpiresAt && now < rec.AbsoluteExpiresAt {
			userID = rec.UserID
			break
		}
	}
	if userID != "" {
		for i := range s.state.Users {
			u := s.state.Users[i]
			if u.ID == userID && u.Status == UserActive && strings.TrimSpace(u.PasswordHash) != "" {
				passwordHash = u.PasswordHash
				break
			}
		}
	}
	s.mu.Unlock()
	if userID == "" || passwordHash == "" {
		return false, nil
	}

	verified, err := verifyPasswordArgon2id(password, passwordHash)
	if err != nil {
		return false, err
	}
	if !verified {
		return false, nil
	}

	s.mu.Lock()
	defer s.mu.Unlock()
	now = s.now().UnixMilli()
	userCurrent := false
	for i := range s.state.Users {
		u := s.state.Users[i]
		if u.ID == userID && u.Status == UserActive && u.PasswordHash == passwordHash {
			userCurrent = true
			break
		}
	}
	if !userCurrent {
		return false, nil
	}
	for i := range s.state.Sessions {
		rec := &s.state.Sessions[i]
		if rec.ID != sessionID || rec.UserID != userID || rec.RevokedAt != 0 || now >= rec.IdleExpiresAt || now >= rec.AbsoluteExpiresAt {
			continue
		}
		previous := rec.AuthenticatedAt
		rec.AuthenticatedAt = now
		if err := s.persistLocked(); err != nil {
			rec.AuthenticatedAt = previous
			return false, err
		}
		return true, nil
	}
	return false, nil
}

func (s *IdentityService) registerHostedDevice(principal Principal, label, fingerprintHash string) (DeviceRecord, error) {
	fingerprintHash = strings.TrimSpace(fingerprintHash)
	if fingerprintHash == "" {
		return DeviceRecord{}, errors.New("device fingerprint hash required")
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	now := s.now().UnixMilli()
	tenantID := normalizedTenantID(principal.TenantID)
	for i := range s.state.Users {
		u := s.state.Users[i]
		if u.ID != principal.UserID {
			continue
		}
		if u.Status != UserActive || normalizedTenantID(u.TenantID) != tenantID {
			return DeviceRecord{}, errors.New("tenant principal unavailable")
		}
		for j := range s.state.Devices {
			d := &s.state.Devices[j]
			if d.UserID == u.ID && normalizedTenantID(d.TenantID) == tenantID && d.FingerprintHash == fingerprintHash && d.Status == DeviceActive {
				d.LastSeenAt = now
				if err := s.persistLocked(); err != nil {
					return DeviceRecord{}, err
				}
				return *d, nil
			}
		}
		d := DeviceRecord{ID: randomID("dev"), TenantID: tenantID, UserID: u.ID, Label: strings.TrimSpace(label), FingerprintHash: fingerprintHash, Status: DeviceActive, CreatedAt: now, LastSeenAt: now}
		s.state.Devices = append(s.state.Devices, d)
		if err := s.persistLocked(); err != nil {
			return DeviceRecord{}, err
		}
		return d, nil
	}
	return DeviceRecord{}, errors.New("tenant principal unavailable")
}

func (s *IdentityService) bindHostedDeviceToSession(principal Principal, deviceID string) error {
	deviceID = strings.TrimSpace(deviceID)
	if principal.SessionID == "" || deviceID == "" {
		return errors.New("session and device required")
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	now := s.now().UnixMilli()
	tenantID := normalizedTenantID(principal.TenantID)
	deviceOK := false
	for i := range s.state.Devices {
		d := &s.state.Devices[i]
		if d.ID == deviceID && d.UserID == principal.UserID && normalizedTenantID(d.TenantID) == tenantID && hostedDeviceActive(*d, now, defaultHostedDeviceStaleTTL) {
			d.LastSeenAt = now
			deviceOK = true
			break
		}
	}
	if !deviceOK {
		return errors.New("device unavailable for tenant principal")
	}
	for i := range s.state.Sessions {
		rec := &s.state.Sessions[i]
		if rec.ID == principal.SessionID && rec.UserID == principal.UserID && normalizedTenantID(rec.TenantID) == tenantID && rec.RevokedAt == 0 && now < rec.IdleExpiresAt && now < rec.AbsoluteExpiresAt {
			rec.DeviceID = deviceID
			return s.persistLocked()
		}
	}
	return errors.New("session unavailable for tenant principal")
}

func hostedDeviceActive(device DeviceRecord, now int64, staleTTL time.Duration) bool {
	if device.Status != DeviceActive || device.RevokedAt != 0 || device.LastSeenAt <= 0 || staleTTL <= 0 {
		return false
	}
	return now-device.LastSeenAt <= int64(staleTTL/time.Millisecond)
}

func (s *IdentityService) setHostedDeviceStatus(principal Principal, deviceID string, status DeviceStatus) error {
	if status != DeviceActive && status != DeviceRevoked && status != DeviceLost {
		return errors.New("invalid device status")
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	now := s.now().UnixMilli()
	tenantID := normalizedTenantID(principal.TenantID)
	for i := range s.state.Devices {
		d := &s.state.Devices[i]
		if d.ID != strings.TrimSpace(deviceID) {
			continue
		}
		if d.UserID != principal.UserID || normalizedTenantID(d.TenantID) != tenantID {
			return errors.New("cross-tenant device mutation denied")
		}
		d.Status = status
		if status == DeviceRevoked || status == DeviceLost {
			d.RevokedAt = now
			for j := range s.state.Sessions {
				rec := &s.state.Sessions[j]
				if rec.DeviceID == d.ID && rec.RevokedAt == 0 {
					rec.RevokedAt = now
				}
			}
		}
		return s.persistLocked()
	}
	return errors.New("device not found")
}

// recordHostedMFAVerification records externally established MFA-class proof.
// It never performs or infers MFA itself; callers must only invoke it after a
// separately verified MFA ceremony. Hosted authorization fails closed if proof
// is required but this timestamp is absent or stale.
func (s *IdentityService) recordHostedMFAVerification(principal Principal) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	now := s.now().UnixMilli()
	for i := range s.state.Sessions {
		rec := &s.state.Sessions[i]
		if rec.ID == principal.SessionID && rec.UserID == principal.UserID && normalizedTenantID(rec.TenantID) == normalizedTenantID(principal.TenantID) && rec.RevokedAt == 0 && now < rec.IdleExpiresAt && now < rec.AbsoluteExpiresAt {
			rec.MFAVerifiedAt = now
			return s.persistLocked()
		}
	}
	return errors.New("session unavailable for MFA proof")
}

func (s *IdentityService) authorizeHostedIdentity(principal Principal, requirement HostedIdentityRequirement) HostedIdentityDecision {
	decision := HostedIdentityDecision{
		TenantID:   normalizedTenantID(requirement.TenantID),
		UserID:     principal.UserID,
		SessionID:  principal.SessionID,
		DeviceID:   principal.DeviceID,
		Capability: strings.ToUpper(strings.TrimSpace(requirement.Capability)),
	}
	blocking := make([]string, 0, 8)
	if strings.TrimSpace(requirement.TenantID) == "" || normalizedTenantID(principal.TenantID) != normalizedTenantID(requirement.TenantID) {
		blocking = append(blocking, "tenant mismatch")
	}
	if decision.Capability == "" || !roleHasHostedCapability(principal.Role, decision.Capability) {
		blocking = append(blocking, "role capability denied")
	}

	s.mu.Lock()
	defer s.mu.Unlock()
	now := s.now().UnixMilli()
	tenantActive := false
	for i := range s.state.Tenants {
		t := s.state.Tenants[i]
		if t.ID == decision.TenantID && t.Status == TenantActive {
			tenantActive = true
			break
		}
	}
	if !tenantActive {
		blocking = append(blocking, "tenant inactive or unknown")
	}
	userActive := false
	for i := range s.state.Users {
		u := s.state.Users[i]
		if u.ID == principal.UserID && u.Status == UserActive && normalizedTenantID(u.TenantID) == decision.TenantID && u.Role == principal.Role {
			userActive = true
			break
		}
	}
	if !userActive {
		blocking = append(blocking, "tenant user inactive or mismatched")
	}

	var session *SessionRecord
	for i := range s.state.Sessions {
		rec := &s.state.Sessions[i]
		if rec.ID == principal.SessionID && rec.UserID == principal.UserID && normalizedTenantID(rec.TenantID) == decision.TenantID {
			session = rec
			break
		}
	}
	if session == nil || session.RevokedAt != 0 || now >= session.IdleExpiresAt || now >= session.AbsoluteExpiresAt {
		blocking = append(blocking, "session invalid expired or revoked")
	} else {
		decision.DeviceID = session.DeviceID
		if requirement.RequireRecentAuthentication {
			authenticatedAt := sessionAuthenticationTime(*session)
			if authenticatedAt <= 0 || authenticatedAt > now || now-authenticatedAt > int64(defaultSensitiveReauthTTL/time.Millisecond) {
				blocking = append(blocking, "recent authentication required")
			}
		}
		if requirement.RequireMFA {
			if session.MFAVerifiedAt <= 0 || session.MFAVerifiedAt > now || now-session.MFAVerifiedAt > int64(defaultSensitiveReauthTTL/time.Millisecond) {
				blocking = append(blocking, "recent MFA proof required")
			}
		}
	}

	if requirement.RequireRegisteredDevice {
		deviceOK := false
		if session != nil && strings.TrimSpace(session.DeviceID) != "" {
			for i := range s.state.Devices {
				d := s.state.Devices[i]
				if d.ID == session.DeviceID && d.UserID == principal.UserID && normalizedTenantID(d.TenantID) == decision.TenantID && hostedDeviceActive(d, now, defaultHostedDeviceStaleTTL) {
					deviceOK = true
					break
				}
			}
		}
		if !deviceOK {
			blocking = append(blocking, "registered active device required")
		}
	}

	decision.BlockingReasons = blocking
	decision.Allowed = len(blocking) == 0
	return decision
}

var errRecentAuthenticationRequired = errors.New("recent password authentication required")
