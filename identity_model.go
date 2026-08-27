package main

import (
	"context"
	"crypto/ed25519"
	"crypto/rand"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/base64"
	"errors"
	"fmt"
	"strings"
	"sync"
	"time"
)

type UserRole string

const (
	RoleSuperOwner UserRole = "SUPER_OWNER"
	RoleOwner      UserRole = "OWNER"
	RoleAdmin      UserRole = "ADMIN"
	RoleUser       UserRole = "USER"
	RoleDemo       UserRole = "DEMO"
)

type UserStatus string

const (
	UserActive   UserStatus = "ACTIVE"
	UserDisabled UserStatus = "DISABLED"
)

type TenantStatus string

const (
	TenantActive   TenantStatus = "ACTIVE"
	TenantDisabled TenantStatus = "DISABLED"
	localTenantID               = "local-default"
)

type DeviceStatus string

const (
	DeviceActive  DeviceStatus = "ACTIVE"
	DeviceRevoked DeviceStatus = "REVOKED"
	DeviceLost    DeviceStatus = "LOST"
)

type TenantRecord struct {
	ID        string       `json:"id"`
	Name      string       `json:"name,omitempty"`
	Status    TenantStatus `json:"status"`
	CreatedAt int64        `json:"createdAt"`
	UpdatedAt int64        `json:"updatedAt"`
}

type DeviceRecord struct {
	ID              string       `json:"id"`
	TenantID        string       `json:"tenantId"`
	UserID          string       `json:"userId"`
	Label           string       `json:"label,omitempty"`
	FingerprintHash string       `json:"fingerprintHash"`
	Status          DeviceStatus `json:"status"`
	CreatedAt       int64        `json:"createdAt"`
	LastSeenAt      int64        `json:"lastSeenAt"`
	RevokedAt       int64        `json:"revokedAt,omitempty"`
}

type UserRecord struct {
	ID              string                      `json:"id"`
	TenantID        string                      `json:"tenantId,omitempty"`
	Username        string                      `json:"username"`
	DisplayName     string                      `json:"displayName,omitempty"`
	Role            UserRole                    `json:"role"`
	Status          UserStatus                  `json:"status"`
	PasswordHash    string                      `json:"passwordHash,omitempty"`
	MustSetPassword bool                        `json:"mustSetPassword,omitempty"`
	MFACredentials  []HostedMFACredentialRecord `json:"mfaCredentials,omitempty"`
	CreatedAt       int64                       `json:"createdAt"`
	UpdatedAt       int64                       `json:"updatedAt"`
	LastLoginAt     int64                       `json:"lastLoginAt,omitempty"`
}

type SessionRecord struct {
	ID                string                    `json:"id"`
	TokenHash         string                    `json:"tokenHash"`
	TenantID          string                    `json:"tenantId,omitempty"`
	UserID            string                    `json:"userId"`
	DeviceID          string                    `json:"deviceId,omitempty"`
	CreatedAt         int64                     `json:"createdAt"`
	AuthenticatedAt   int64                     `json:"authenticatedAt,omitempty"`
	MFAVerifiedAt     int64                     `json:"mfaVerifiedAt,omitempty"`
	MFAChallenge      *HostedMFAChallengeRecord `json:"mfaChallenge,omitempty"`
	LastSeenAt        int64                     `json:"lastSeenAt"`
	IdleExpiresAt     int64                     `json:"idleExpiresAt"`
	AbsoluteExpiresAt int64                     `json:"absoluteExpiresAt"`
	RevokedAt         int64                     `json:"revokedAt,omitempty"`
	RotatedFrom       string                    `json:"rotatedFrom,omitempty"`
}

type IdentitySecurityEventType string

const (
	IdentitySecurityDeviceRegistered  IdentitySecurityEventType = "DEVICE_REGISTERED"
	IdentitySecurityDeviceBound       IdentitySecurityEventType = "DEVICE_BOUND"
	IdentitySecurityDeviceStale       IdentitySecurityEventType = "DEVICE_STALE"
	IdentitySecurityDeviceLost        IdentitySecurityEventType = "DEVICE_LOST"
	IdentitySecurityDeviceRevoked     IdentitySecurityEventType = "DEVICE_REVOKED"
	IdentitySecurityDeviceReactivated IdentitySecurityEventType = "DEVICE_REACTIVATED"
	IdentitySecuritySessionRevoked    IdentitySecurityEventType = "SESSION_REVOKED"
	maxIdentitySecurityEvents                                   = 512
)

type IdentitySecurityEvent struct {
	ID        string                    `json:"id"`
	TenantID  string                    `json:"tenantId"`
	UserID    string                    `json:"userId,omitempty"`
	DeviceID  string                    `json:"deviceId,omitempty"`
	SessionID string                    `json:"sessionId,omitempty"`
	Type      IdentitySecurityEventType `json:"type"`
	CreatedAt int64                     `json:"createdAt"`
}

type IdentityPersistentState struct {
	Version             int                        `json:"version"`
	Tenants             []TenantRecord             `json:"tenants,omitempty"`
	Users               []UserRecord               `json:"users"`
	Devices             []DeviceRecord             `json:"devices,omitempty"`
	Sessions            []SessionRecord            `json:"sessions"`
	SecurityEvents      []IdentitySecurityEvent    `json:"securityEvents,omitempty"`
	ProductEntitlements []TenantProductEntitlement `json:"productEntitlements,omitempty"`
	UpdatedAt           int64                      `json:"updatedAt"`
}

type Principal struct {
	TenantID    string   `json:"tenantId,omitempty"`
	UserID      string   `json:"userId"`
	Username    string   `json:"username"`
	DisplayName string   `json:"displayName,omitempty"`
	Role        UserRole `json:"role"`
	SessionID   string   `json:"sessionId"`
	DeviceID    string   `json:"deviceId,omitempty"`
}

type identityContextKey struct{}

func principalFromContext(ctx context.Context) (Principal, bool) {
	p, ok := ctx.Value(identityContextKey{}).(Principal)
	return p, ok
}

func normalizeUsername(v string) string {
	return strings.ToLower(strings.TrimSpace(v))
}

func normalizeDisplayName(v string) (string, error) {
	v = strings.Join(strings.Fields(v), " ")
	runes := []rune(v)
	if len(runes) < 1 || len(runes) > 64 {
		return "", errors.New("display name must be between 1 and 64 characters")
	}
	for _, r := range runes {
		if r < 0x20 || r == 0x7f {
			return "", errors.New("display name contains unsupported characters")
		}
	}
	return v, nil
}

func validRole(r UserRole) bool {
	switch r {
	case RoleSuperOwner, RoleOwner, RoleAdmin, RoleUser, RoleDemo:
		return true
	default:
		return false
	}
}

func roleRank(r UserRole) int {
	switch r {
	case RoleSuperOwner:
		return 50
	case RoleOwner:
		return 40
	case RoleAdmin:
		return 30
	case RoleUser:
		return 20
	case RoleDemo:
		return 10
	default:
		return 0
	}
}

func roleAtLeast(actual, required UserRole) bool { return roleRank(actual) >= roleRank(required) }

func normalizedTenantID(value string) string {
	value = strings.TrimSpace(value)
	if value == "" {
		return localTenantID
	}
	return value
}

const (
	defaultSessionIdleTTL     = 2 * time.Hour
	defaultSessionAbsoluteTTL = 24 * time.Hour
	defaultSensitiveReauthTTL = 15 * time.Minute
	sessionCookieName         = "depulse_session"
	csrfCookieName            = "depulse_csrf"
	bootstrapOwnerID          = "bootstrap-owner"
)

type IdentityService struct {
	mu                sync.Mutex
	persistence       *PersistenceManager
	state             IdentityPersistentState
	persistedDevices  []DeviceRecord
	persistedSessions []SessionRecord
	now               func() time.Time
	idleTTL           time.Duration
	absoluteTTL       time.Duration
}

func cloneDeviceRecords(in []DeviceRecord) []DeviceRecord {
	return append([]DeviceRecord(nil), in...)
}

func cloneSessionRecords(in []SessionRecord) []SessionRecord {
	out := append([]SessionRecord(nil), in...)
	for i := range out {
		if in[i].MFAChallenge != nil {
			challenge := *in[i].MFAChallenge
			out[i].MFAChallenge = &challenge
		}
	}
	return out
}

func NewIdentityService(p *PersistenceManager) (*IdentityService, error) {
	if p == nil {
		return nil, errors.New("identity requires persistence")
	}
	s := &IdentityService{persistence: p, now: time.Now, idleTTL: defaultSessionIdleTTL, absoluteTTL: defaultSessionAbsoluteTTL}
	st, err := p.LoadIdentityState(context.Background())
	if err != nil {
		return nil, fmt.Errorf("load identity state: %w", err)
	}
	if st.Version == 0 {
		st.Version = 1
	}
	s.state = st
	s.persistedDevices = cloneDeviceRecords(st.Devices)
	s.persistedSessions = cloneSessionRecords(st.Sessions)
	if s.normalizeHostedIdentityState() {
		if err := s.persistLocked(); err != nil {
			return nil, fmt.Errorf("persist identity tenant migration: %w", err)
		}
	}
	if err := s.ensureHostedProductEntitlements(); err != nil {
		return nil, fmt.Errorf("persist identity product entitlement migration: %w", err)
	}
	if err := s.ensureBootstrapOwnerLocked(); err != nil {
		return nil, err
	}
	return s, nil
}

func (s *IdentityService) normalizeHostedIdentityState() bool {
	changed := false
	now := s.now().UnixMilli()
	localTenantFound := false
	for i := range s.state.Tenants {
		if s.state.Tenants[i].ID == localTenantID {
			localTenantFound = true
			if s.state.Tenants[i].Status == "" {
				s.state.Tenants[i].Status = TenantActive
				changed = true
			}
		}
	}
	if !localTenantFound {
		s.state.Tenants = append(s.state.Tenants, TenantRecord{ID: localTenantID, Name: "Local", Status: TenantActive, CreatedAt: now, UpdatedAt: now})
		changed = true
	}
	userTenant := make(map[string]string, len(s.state.Users))
	for i := range s.state.Users {
		tenantID := normalizedTenantID(s.state.Users[i].TenantID)
		if s.state.Users[i].TenantID != tenantID {
			s.state.Users[i].TenantID = tenantID
			changed = true
		}
		userTenant[s.state.Users[i].ID] = tenantID
	}
	for i := range s.state.Sessions {
		tenantID := strings.TrimSpace(s.state.Sessions[i].TenantID)
		if tenantID == "" {
			tenantID = normalizedTenantID(userTenant[s.state.Sessions[i].UserID])
			s.state.Sessions[i].TenantID = tenantID
			changed = true
		}
	}
	return changed
}

func (s *IdentityService) ensureBootstrapOwnerLocked() error {
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, u := range s.state.Users {
		if u.Role == RoleSuperOwner || u.Role == RoleOwner {
			return nil
		}
	}
	now := s.now().UnixMilli()
	s.state.Users = append(s.state.Users, UserRecord{
		ID: bootstrapOwnerID, TenantID: localTenantID, Username: "owner", DisplayName: "Local Owner", Role: RoleOwner,
		Status: UserActive, MustSetPassword: true, CreatedAt: now, UpdatedAt: now,
	})
	return s.persistLocked()
}

func (s *IdentityService) appendIdentitySecurityEventLocked(eventType IdentitySecurityEventType, tenantID, userID, deviceID, sessionID string, createdAt int64) {
	s.state.SecurityEvents = append(s.state.SecurityEvents, IdentitySecurityEvent{
		ID:        randomID("aud"),
		TenantID:  normalizedTenantID(tenantID),
		UserID:    strings.TrimSpace(userID),
		DeviceID:  strings.TrimSpace(deviceID),
		SessionID: strings.TrimSpace(sessionID),
		Type:      eventType,
		CreatedAt: createdAt,
	})
	if len(s.state.SecurityEvents) > maxIdentitySecurityEvents {
		start := len(s.state.SecurityEvents) - maxIdentitySecurityEvents
		trimmed := append([]IdentitySecurityEvent(nil), s.state.SecurityEvents[start:]...)
		s.state.SecurityEvents = trimmed
	}
}

func (s *IdentityService) deriveIdentitySecurityEventsLocked(createdAt int64) {
	previousDevices := make(map[string]DeviceRecord, len(s.persistedDevices))
	for _, device := range s.persistedDevices {
		previousDevices[device.ID] = device
	}
	for _, current := range s.state.Devices {
		previous, found := previousDevices[current.ID]
		if !found {
			s.appendIdentitySecurityEventLocked(IdentitySecurityDeviceRegistered, current.TenantID, current.UserID, current.ID, "", createdAt)
			continue
		}
		if current.Status == previous.Status {
			continue
		}
		eventType := IdentitySecurityEventType("")
		switch current.Status {
		case DeviceStale:
			eventType = IdentitySecurityDeviceStale
		case DeviceLost:
			eventType = IdentitySecurityDeviceLost
		case DeviceRevoked:
			eventType = IdentitySecurityDeviceRevoked
		case DeviceActive:
			eventType = IdentitySecurityDeviceReactivated
		}
		if eventType != "" {
			s.appendIdentitySecurityEventLocked(eventType, current.TenantID, current.UserID, current.ID, "", createdAt)
		}
	}

	previousSessions := make(map[string]SessionRecord, len(s.persistedSessions))
	for _, session := range s.persistedSessions {
		previousSessions[session.ID] = session
	}
	for _, current := range s.state.Sessions {
		previous, found := previousSessions[current.ID]
		if !found {
			continue
		}
		if strings.TrimSpace(current.DeviceID) != "" && current.DeviceID != previous.DeviceID {
			s.appendIdentitySecurityEventLocked(IdentitySecurityDeviceBound, current.TenantID, current.UserID, current.DeviceID, current.ID, createdAt)
		}
		if previous.RevokedAt == 0 && current.RevokedAt > 0 {
			s.appendIdentitySecurityEventLocked(IdentitySecuritySessionRevoked, current.TenantID, current.UserID, current.DeviceID, current.ID, createdAt)
		}
	}
}

func (s *IdentityService) persistLocked() error {
	previousEvents := append([]IdentitySecurityEvent(nil), s.state.SecurityEvents...)
	now := s.now().UnixMilli()
	s.deriveIdentitySecurityEventsLocked(now)
	s.state.Version = 1
	s.state.UpdatedAt = now
	if err := s.persistence.SaveIdentityState(context.Background(), s.state); err != nil {
		s.state.SecurityEvents = previousEvents
		return err
	}
	s.persistedDevices = cloneDeviceRecords(s.state.Devices)
	s.persistedSessions = cloneSessionRecords(s.state.Sessions)
	return nil
}

func (s *IdentityService) adminSecurityEvents(actor Principal) ([]IdentitySecurityEvent, error) {
	if !roleAtLeast(actor.Role, RoleAdmin) {
		return nil, errors.New("insufficient authority")
	}
	tenantID := normalizedTenantID(actor.TenantID)
	s.mu.Lock()
	defer s.mu.Unlock()
	out := make([]IdentitySecurityEvent, 0, len(s.state.SecurityEvents))
	for i := len(s.state.SecurityEvents) - 1; i >= 0; i-- {
		event := s.state.SecurityEvents[i]
		if normalizedTenantID(event.TenantID) == tenantID {
			out = append(out, event)
		}
	}
	return out, nil
}

func (s *IdentityService) ownerNeedsPassword() bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, u := range s.state.Users {
		if (u.Role == RoleOwner || u.Role == RoleSuperOwner) && u.Status == UserActive && (u.MustSetPassword || strings.TrimSpace(u.PasswordHash) == "") {
			return true
		}
	}
	return false
}

func (s *IdentityService) userRequiresPassword(userID string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, u := range s.state.Users {
		if u.ID == userID && u.Status == UserActive {
			return u.MustSetPassword || strings.TrimSpace(u.PasswordHash) == ""
		}
	}
	return false
}

func (s *IdentityService) bootstrapOwnerSession() (string, Principal, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	for i := range s.state.Users {
		u := &s.state.Users[i]
		if (u.Role == RoleOwner || u.Role == RoleSuperOwner) && u.Status == UserActive && (u.MustSetPassword || strings.TrimSpace(u.PasswordHash) == "") {
			return s.createSessionLocked(*u, "")
		}
	}
	return "", Principal{}, errors.New("bootstrap owner is not available")
}

func (s *IdentityService) authenticate(username, password string) (string, Principal, error) {
	username = normalizeUsername(username)
	if username == "" || password == "" {
		return "", Principal{}, errors.New("invalid username or password")
	}

	s.mu.Lock()
	var candidate UserRecord
	found := false
	for i := range s.state.Users {
		if normalizeUsername(s.state.Users[i].Username) == username {
			candidate = s.state.Users[i]
			found = true
			break
		}
	}
	s.mu.Unlock()
	if !found || candidate.Status != UserActive || strings.TrimSpace(candidate.PasswordHash) == "" {
		return "", Principal{}, errors.New("invalid username or password")
	}
	ok, err := verifyPasswordArgon2id(password, candidate.PasswordHash)
	if err != nil || !ok {
		return "", Principal{}, errors.New("invalid username or password")
	}

	s.mu.Lock()
	defer s.mu.Unlock()
	idx := -1
	for i := range s.state.Users {
		if s.state.Users[i].ID == candidate.ID {
			idx = i
			break
		}
	}
	if idx < 0 || s.state.Users[idx].Status != UserActive || s.state.Users[idx].PasswordHash != candidate.PasswordHash || normalizeUsername(s.state.Users[idx].Username) != username {
		return "", Principal{}, errors.New("invalid username or password")
	}
	now := s.now().UnixMilli()
	s.state.Users[idx].LastLoginAt = now
	s.state.Users[idx].UpdatedAt = now
	return s.createSessionLocked(s.state.Users[idx], "")
}

func (s *IdentityService) setPassword(userID, password string) (string, Principal, error) {
	hash, err := hashPasswordArgon2id(password)
	if err != nil {
		return "", Principal{}, err
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	idx := -1
	for i := range s.state.Users {
		if s.state.Users[i].ID == userID {
			idx = i
			break
		}
	}
	if idx < 0 || s.state.Users[idx].Status != UserActive {
		return "", Principal{}, errors.New("user unavailable")
	}
	now := s.now().UnixMilli()
	s.state.Users[idx].PasswordHash = hash
	s.state.Users[idx].MustSetPassword = false
	s.state.Users[idx].UpdatedAt = now
	for i := range s.state.Sessions {
		if s.state.Sessions[i].UserID == userID && s.state.Sessions[i].RevokedAt == 0 {
			s.state.Sessions[i].RevokedAt = now
		}
	}
	return s.createSessionLocked(s.state.Users[idx], "")
}

func (s *IdentityService) updateProfile(userID, username, displayName string) error {
	username = normalizeUsername(username)
	if !validAdminUsername(username) {
		return errors.New("username must be 3 to 64 characters using letters, numbers, dot, underscore or hyphen")
	}
	displayName, err := normalizeDisplayName(displayName)
	if err != nil {
		return err
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	idx := -1
	for i := range s.state.Users {
		if s.state.Users[i].ID == userID {
			idx = i
			continue
		}
		if normalizeUsername(s.state.Users[i].Username) == username {
			return errors.New("username already exists")
		}
	}
	if idx < 0 || s.state.Users[idx].Status != UserActive {
		return errors.New("user unavailable")
	}
	previousUsername := s.state.Users[idx].Username
	previousName := s.state.Users[idx].DisplayName
	previousUpdatedAt := s.state.Users[idx].UpdatedAt
	s.state.Users[idx].Username = username
	s.state.Users[idx].DisplayName = displayName
	s.state.Users[idx].UpdatedAt = s.now().UnixMilli()
	if err := s.persistLocked(); err != nil {
		s.state.Users[idx].Username = previousUsername
		s.state.Users[idx].DisplayName = previousName
		s.state.Users[idx].UpdatedAt = previousUpdatedAt
		return err
	}
	return nil
}

func sessionTokenHash(token string) string {
	sum := sha256.Sum256([]byte(token))
	return base64.RawURLEncoding.EncodeToString(sum[:])
}

func sessionAuthenticationTime(rec SessionRecord) int64 {
	if rec.AuthenticatedAt > 0 {
		return rec.AuthenticatedAt
	}
	return rec.CreatedAt
}

func (s *IdentityService) createSessionLocked(u UserRecord, rotatedFrom string) (string, Principal, error) {
	now := s.now()
	return s.createSessionWithSecurityLocked(u, rotatedFrom, now.UnixMilli(), now.Add(s.absoluteTTL).UnixMilli())
}

func (s *IdentityService) createSessionWithSecurityLocked(u UserRecord, rotatedFrom string, authenticatedAt, absoluteExpiresAt int64) (string, Principal, error) {
	return s.createSessionWithHostedSecurityLocked(u, rotatedFrom, authenticatedAt, absoluteExpiresAt, "", 0)
}

func (s *IdentityService) createSessionWithHostedSecurityLocked(u UserRecord, rotatedFrom string, authenticatedAt, absoluteExpiresAt int64, deviceID string, mfaVerifiedAt int64) (string, Principal, error) {
	raw := randomID("ses") + randomID("")
	now := s.now()
	nowMillis := now.UnixMilli()
	if authenticatedAt <= 0 || authenticatedAt > nowMillis {
		authenticatedAt = nowMillis
	}
	if mfaVerifiedAt < 0 || mfaVerifiedAt > nowMillis {
		mfaVerifiedAt = 0
	}
	if absoluteExpiresAt <= nowMillis {
		return "", Principal{}, errors.New("session lifetime exhausted")
	}
	cutoff := now.Add(-7 * 24 * time.Hour).UnixMilli()
	kept := s.state.Sessions[:0]
	for _, existing := range s.state.Sessions {
		if existing.RevokedAt > 0 && existing.RevokedAt < cutoff {
			continue
		}
		if existing.AbsoluteExpiresAt > 0 && existing.AbsoluteExpiresAt < cutoff {
			continue
		}
		kept = append(kept, existing)
	}
	s.state.Sessions = kept
	idleExpiresAt := now.Add(s.idleTTL).UnixMilli()
	if idleExpiresAt > absoluteExpiresAt {
		idleExpiresAt = absoluteExpiresAt
	}
	tenantID := normalizedTenantID(u.TenantID)
	rec := SessionRecord{ID: randomID("sid"), TokenHash: sessionTokenHash(raw), TenantID: tenantID, UserID: u.ID, DeviceID: strings.TrimSpace(deviceID), CreatedAt: nowMillis, AuthenticatedAt: authenticatedAt, MFAVerifiedAt: mfaVerifiedAt, LastSeenAt: nowMillis, IdleExpiresAt: idleExpiresAt, AbsoluteExpiresAt: absoluteExpiresAt, RotatedFrom: rotatedFrom}
	s.state.Sessions = append(s.state.Sessions, rec)
	if err := s.persistLocked(); err != nil {
		return "", Principal{}, err
	}
	return raw, Principal{TenantID: tenantID, UserID: u.ID, Username: u.Username, DisplayName: u.DisplayName, Role: u.Role, SessionID: rec.ID, DeviceID: rec.DeviceID}, nil
}

func (s *IdentityService) resolve(token string, touch bool) (Principal, error) {
	if strings.TrimSpace(token) == "" {
		return Principal{}, errors.New("missing session")
	}
	target := sessionTokenHash(token)
	s.mu.Lock()
	defer s.mu.Unlock()
	now := s.now().UnixMilli()
	si := -1
	for i := range s.state.Sessions {
		if subtle.ConstantTimeCompare([]byte(s.state.Sessions[i].TokenHash), []byte(target)) == 1 {
			si = i
			break
		}
	}
	if si < 0 {
		return Principal{}, errors.New("invalid session")
	}
	rec := &s.state.Sessions[si]
	if rec.RevokedAt != 0 || now >= rec.IdleExpiresAt || now >= rec.AbsoluteExpiresAt {
		return Principal{}, errors.New("expired session")
	}
	var u *UserRecord
	for i := range s.state.Users {
		if s.state.Users[i].ID == rec.UserID {
			u = &s.state.Users[i]
			break
		}
	}
	if u == nil || u.Status != UserActive || !validRole(u.Role) {
		return Principal{}, errors.New("inactive principal")
	}
	tenantID := normalizedTenantID(u.TenantID)
	if normalizedTenantID(rec.TenantID) != tenantID {
		return Principal{}, errors.New("session tenant mismatch")
	}
	if touch {
		if now-rec.LastSeenAt >= int64(time.Minute/time.Millisecond) {
			rec.LastSeenAt = now
			next := s.now().Add(s.idleTTL).UnixMilli()
			if next > rec.AbsoluteExpiresAt {
				next = rec.AbsoluteExpiresAt
			}
			rec.IdleExpiresAt = next
			_ = s.persistLocked()
		}
	}
	return Principal{TenantID: tenantID, UserID: u.ID, Username: u.Username, DisplayName: u.DisplayName, Role: u.Role, SessionID: rec.ID, DeviceID: rec.DeviceID}, nil
}

func (s *IdentityService) revokeToken(token string) error {
	hash := sessionTokenHash(token)
	s.mu.Lock()
	defer s.mu.Unlock()
	now := s.now().UnixMilli()
	changed := false
	for i := range s.state.Sessions {
		if subtle.ConstantTimeCompare([]byte(s.state.Sessions[i].TokenHash), []byte(hash)) == 1 && s.state.Sessions[i].RevokedAt == 0 {
			s.state.Sessions[i].RevokedAt = now
			changed = true
		}
	}
	if changed {
		return s.persistLocked()
	}
	return nil
}

func (s *IdentityService) rotate(token string) (string, Principal, error) {
	p, err := s.resolve(token, false)
	if err != nil {
		return "", Principal{}, err
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	now := s.now().UnixMilli()
	var u UserRecord
	found := false
	for i := range s.state.Users {
		if s.state.Users[i].ID == p.UserID {
			u = s.state.Users[i]
			found = true
			break
		}
	}
	if !found {
		return "", Principal{}, errors.New("principal unavailable")
	}
	oldHash := sessionTokenHash(token)
	oldID := ""
	authenticatedAt := int64(0)
	absoluteExpiresAt := int64(0)
	deviceID := ""
	mfaVerifiedAt := int64(0)
	for i := range s.state.Sessions {
		rec := &s.state.Sessions[i]
		if subtle.ConstantTimeCompare([]byte(rec.TokenHash), []byte(oldHash)) == 1 && rec.RevokedAt == 0 {
			oldID = rec.ID
			authenticatedAt = sessionAuthenticationTime(*rec)
			absoluteExpiresAt = rec.AbsoluteExpiresAt
			deviceID = rec.DeviceID
			mfaVerifiedAt = rec.MFAVerifiedAt
			if absoluteExpiresAt <= now {
				return "", Principal{}, errors.New("expired session")
			}
			rec.RevokedAt = now
			break
		}
	}
	if oldID == "" {
		return "", Principal{}, errors.New("invalid session")
	}
	return s.createSessionWithHostedSecurityLocked(u, oldID, authenticatedAt, absoluteExpiresAt, deviceID, mfaVerifiedAt)
}

func (s *IdentityService) sessionRecentlyAuthenticated(sessionID string, maxAge time.Duration) bool {
	if strings.TrimSpace(sessionID) == "" || maxAge <= 0 {
		return false
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	now := s.now().UnixMilli()
	for i := range s.state.Sessions {
		rec := s.state.Sessions[i]
		if rec.ID != sessionID || rec.RevokedAt != 0 || now >= rec.IdleExpiresAt || now >= rec.AbsoluteExpiresAt {
			continue
		}
		authenticatedAt := sessionAuthenticationTime(rec)
		return authenticatedAt > 0 && authenticatedAt <= now && now-authenticatedAt <= int64(maxAge/time.Millisecond)
	}
	return false
}

func (s *IdentityService) status(token string) map[string]any {
	out := map[string]any{"authenticated": false, "bootstrapRequired": s.ownerNeedsPassword(), "recentAuthentication": false}
	if p, err := s.resolve(token, false); err == nil {
		out["authenticated"] = true
		out["principal"] = p
		out["recentAuthentication"] = s.sessionRecentlyAuthenticated(p.SessionID, defaultSensitiveReauthTTL)
	}
	return out
}

const (
	hostedMFAAlgorithmEd25519 = "ed25519"
	hostedMFADomain           = "DE-PULSE-HOSTED-MFA-V1"
	hostedMFAChallengeTTL     = 5 * time.Minute
	maxHostedMFACredentials   = 8
	maxHostedMFALabelLength   = 128
)

const (
	IdentitySecurityMFAEnrolled           IdentitySecurityEventType = "MFA_CREDENTIAL_ENROLLED"
	IdentitySecurityMFAVerified           IdentitySecurityEventType = "MFA_VERIFIED"
	IdentitySecurityMFAVerificationFailed IdentitySecurityEventType = "MFA_VERIFICATION_FAILED"
	IdentitySecurityMFARevoked            IdentitySecurityEventType = "MFA_CREDENTIAL_REVOKED"
)

var (
	errHostedMFAInvalidProof        = errors.New("invalid MFA proof")
	errHostedMFACredentialMissing  = errors.New("MFA credential unavailable")
	errHostedMFAChallengeMissing   = errors.New("MFA challenge unavailable")
	errHostedMFARecentAuthRequired = errors.New("recent authentication required")
)

type HostedMFACredentialRecord struct {
	ID         string `json:"id"`
	Label      string `json:"label,omitempty"`
	Algorithm  string `json:"algorithm"`
	PublicKey  string `json:"publicKey"`
	CreatedAt  int64  `json:"createdAt"`
	LastUsedAt int64  `json:"lastUsedAt,omitempty"`
	RevokedAt  int64  `json:"revokedAt,omitempty"`
}

type HostedMFAChallengeRecord struct {
	ID            string `json:"id"`
	CredentialID  string `json:"credentialId"`
	ChallengeHash string `json:"challengeHash"`
	CreatedAt     int64  `json:"createdAt"`
	ExpiresAt     int64  `json:"expiresAt"`
}

type HostedMFACredentialView struct {
	ID         string `json:"id"`
	Label      string `json:"label,omitempty"`
	Algorithm  string `json:"algorithm"`
	CreatedAt  int64  `json:"createdAt"`
	LastUsedAt int64  `json:"lastUsedAt,omitempty"`
	RevokedAt  int64  `json:"revokedAt,omitempty"`
}

type HostedMFAChallenge struct {
	ID             string `json:"challengeId"`
	CredentialID   string `json:"credentialId"`
	Algorithm      string `json:"algorithm"`
	SigningPayload string `json:"signingPayload"`
	ExpiresAt      int64  `json:"expiresAt"`
}

func hostedMFACredentialView(rec HostedMFACredentialRecord) HostedMFACredentialView {
	return HostedMFACredentialView{ID: rec.ID, Label: rec.Label, Algorithm: rec.Algorithm, CreatedAt: rec.CreatedAt, LastUsedAt: rec.LastUsedAt, RevokedAt: rec.RevokedAt}
}

func decodeHostedMFAPublicKey(encoded string) (ed25519.PublicKey, string, error) {
	encoded = strings.TrimSpace(encoded)
	if encoded == "" || len(encoded) > 128 {
		return nil, "", errors.New("invalid MFA public key")
	}
	decoded, err := base64.RawURLEncoding.DecodeString(encoded)
	if err != nil || len(decoded) != ed25519.PublicKeySize {
		return nil, "", errors.New("invalid MFA public key")
	}
	canonical := base64.RawURLEncoding.EncodeToString(decoded)
	return ed25519.PublicKey(decoded), canonical, nil
}

func hostedMFASigningPayload(challengeID string, principal Principal, credentialID, nonce string) []byte {
	return []byte(strings.Join([]string{hostedMFADomain, strings.TrimSpace(challengeID), strings.TrimSpace(principal.SessionID), normalizedTenantID(principal.TenantID), strings.TrimSpace(principal.UserID), strings.TrimSpace(credentialID), strings.TrimSpace(nonce)}, "\n"))
}

func hostedMFAChallengeHash(payload []byte) string {
	sum := sha256.Sum256(payload)
	return base64.RawURLEncoding.EncodeToString(sum[:])
}

func hostedMFASessionMatches(rec SessionRecord, principal Principal, now int64) bool {
	return rec.ID == principal.SessionID && rec.UserID == principal.UserID && normalizedTenantID(rec.TenantID) == normalizedTenantID(principal.TenantID) && rec.RevokedAt == 0 && now < rec.IdleExpiresAt && now < rec.AbsoluteExpiresAt
}

func hostedMFASessionRecentlyAuthenticated(rec SessionRecord, now int64) bool {
	authenticatedAt := sessionAuthenticationTime(rec)
	return authenticatedAt > 0 && authenticatedAt <= now && now-authenticatedAt <= int64(defaultSensitiveReauthTTL/time.Millisecond)
}

func (s *IdentityService) enrollHostedMFACredential(principal Principal, label, encodedPublicKey string) (HostedMFACredentialRecord, error) {
	if s == nil {
		return HostedMFACredentialRecord{}, errors.New("identity unavailable")
	}
	_, canonicalKey, err := decodeHostedMFAPublicKey(encodedPublicKey)
	if err != nil {
		return HostedMFACredentialRecord{}, err
	}
	label = strings.Join(strings.Fields(label), " ")
	if len([]rune(label)) > maxHostedMFALabelLength {
		return HostedMFACredentialRecord{}, errors.New("MFA credential label too long")
	}
	if label == "" {
		label = "Security key"
	}

	s.mu.Lock()
	defer s.mu.Unlock()
	now := s.now().UnixMilli()
	sessionOK := false
	for i := range s.state.Sessions {
		if hostedMFASessionMatches(s.state.Sessions[i], principal, now) {
			if !hostedMFASessionRecentlyAuthenticated(s.state.Sessions[i], now) {
				return HostedMFACredentialRecord{}, errHostedMFARecentAuthRequired
			}
			sessionOK = true
			break
		}
	}
	if !sessionOK {
		return HostedMFACredentialRecord{}, errors.New("session unavailable for MFA enrollment")
	}

	for i := range s.state.Users {
		u := &s.state.Users[i]
		if u.ID != principal.UserID {
			continue
		}
		if u.Status != UserActive || normalizedTenantID(u.TenantID) != normalizedTenantID(principal.TenantID) {
			return HostedMFACredentialRecord{}, errors.New("tenant principal unavailable")
		}
		active := 0
		for _, existing := range u.MFACredentials {
			if existing.RevokedAt == 0 {
				active++
				if existing.Algorithm == hostedMFAAlgorithmEd25519 && existing.PublicKey == canonicalKey {
					return HostedMFACredentialRecord{}, errors.New("MFA credential already enrolled")
				}
			}
		}
		if active >= maxHostedMFACredentials {
			return HostedMFACredentialRecord{}, errors.New("MFA credential limit reached")
		}
		credential := HostedMFACredentialRecord{ID: randomID("mfa"), Label: label, Algorithm: hostedMFAAlgorithmEd25519, PublicKey: canonicalKey, CreatedAt: now}
		previousCredentials := append([]HostedMFACredentialRecord(nil), u.MFACredentials...)
		previousUpdatedAt := u.UpdatedAt
		previousEvents := append([]IdentitySecurityEvent(nil), s.state.SecurityEvents...)
		u.MFACredentials = append(u.MFACredentials, credential)
		u.UpdatedAt = now
		s.appendIdentitySecurityEventLocked(IdentitySecurityMFAEnrolled, principal.TenantID, principal.UserID, principal.DeviceID, principal.SessionID, now)
		if err := s.persistLocked(); err != nil {
			u.MFACredentials = previousCredentials
			u.UpdatedAt = previousUpdatedAt
			s.state.SecurityEvents = previousEvents
			return HostedMFACredentialRecord{}, err
		}
		return credential, nil
	}
	return HostedMFACredentialRecord{}, errors.New("tenant principal unavailable")
}

func (s *IdentityService) listHostedMFACredentials(principal Principal) ([]HostedMFACredentialView, bool, error) {
	if s == nil {
		return nil, false, errors.New("identity unavailable")
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	now := s.now().UnixMilli()
	sessionOK := false
	verified := false
	for i := range s.state.Sessions {
		rec := s.state.Sessions[i]
		if hostedMFASessionMatches(rec, principal, now) {
			sessionOK = true
			verified = rec.MFAVerifiedAt > 0 && rec.MFAVerifiedAt <= now && now-rec.MFAVerifiedAt <= int64(defaultSensitiveReauthTTL/time.Millisecond)
			break
		}
	}
	if !sessionOK {
		return nil, false, errors.New("session unavailable for MFA status")
	}
	for _, user := range s.state.Users {
		if user.ID != principal.UserID || user.Status != UserActive || normalizedTenantID(user.TenantID) != normalizedTenantID(principal.TenantID) {
			continue
		}
		out := make([]HostedMFACredentialView, 0, len(user.MFACredentials))
		for _, credential := range user.MFACredentials {
			out = append(out, hostedMFACredentialView(credential))
		}
		return out, verified, nil
	}
	return nil, false, errors.New("tenant principal unavailable")
}

func (s *IdentityService) createHostedMFAChallenge(principal Principal, credentialID string) (HostedMFAChallenge, error) {
	if s == nil {
		return HostedMFAChallenge{}, errors.New("identity unavailable")
	}
	nonceBytes := make([]byte, 32)
	if _, err := rand.Read(nonceBytes); err != nil {
		return HostedMFAChallenge{}, fmt.Errorf("generate MFA challenge: %w", err)
	}
	nonce := base64.RawURLEncoding.EncodeToString(nonceBytes)

	s.mu.Lock()
	defer s.mu.Unlock()
	now := s.now().UnixMilli()
	sessionIndex := -1
	for i := range s.state.Sessions {
		if hostedMFASessionMatches(s.state.Sessions[i], principal, now) {
			sessionIndex = i
			break
		}
	}
	if sessionIndex < 0 {
		return HostedMFAChallenge{}, errors.New("session unavailable for MFA challenge")
	}

	var selected HostedMFACredentialRecord
	found := false
	for _, user := range s.state.Users {
		if user.ID != principal.UserID || user.Status != UserActive || normalizedTenantID(user.TenantID) != normalizedTenantID(principal.TenantID) {
			continue
		}
		for _, credential := range user.MFACredentials {
			if credential.RevokedAt != 0 || credential.Algorithm != hostedMFAAlgorithmEd25519 {
				continue
			}
			if strings.TrimSpace(credentialID) == "" || credential.ID == strings.TrimSpace(credentialID) {
				if _, _, err := decodeHostedMFAPublicKey(credential.PublicKey); err != nil {
					continue
				}
				selected = credential
				found = true
				break
			}
		}
		break
	}
	if !found {
		return HostedMFAChallenge{}, errHostedMFACredentialMissing
	}

	challengeID := randomID("mfc")
	payload := hostedMFASigningPayload(challengeID, principal, selected.ID, nonce)
	expiresAt := s.now().Add(hostedMFAChallengeTTL).UnixMilli()
	rec := &s.state.Sessions[sessionIndex]
	if expiresAt > rec.IdleExpiresAt {
		expiresAt = rec.IdleExpiresAt
	}
	if expiresAt > rec.AbsoluteExpiresAt {
		expiresAt = rec.AbsoluteExpiresAt
	}
	if expiresAt <= now {
		return HostedMFAChallenge{}, errors.New("session expires before MFA challenge")
	}
	previous := rec.MFAChallenge
	rec.MFAChallenge = &HostedMFAChallengeRecord{ID: challengeID, CredentialID: selected.ID, ChallengeHash: hostedMFAChallengeHash(payload), CreatedAt: now, ExpiresAt: expiresAt}
	if err := s.persistLocked(); err != nil {
		rec.MFAChallenge = previous
		return HostedMFAChallenge{}, err
	}
	return HostedMFAChallenge{ID: challengeID, CredentialID: selected.ID, Algorithm: selected.Algorithm, SigningPayload: base64.RawURLEncoding.EncodeToString(payload), ExpiresAt: expiresAt}, nil
}

func (s *IdentityService) verifyHostedMFAChallenge(principal Principal, challengeID, credentialID, encodedPayload, encodedSignature string) error {
	if s == nil {
		return errors.New("identity unavailable")
	}
	payload, err := base64.RawURLEncoding.DecodeString(strings.TrimSpace(encodedPayload))
	if err != nil || len(payload) == 0 || len(payload) > 2048 {
		return errHostedMFAInvalidProof
	}
	signature, err := base64.RawURLEncoding.DecodeString(strings.TrimSpace(encodedSignature))
	if err != nil || len(signature) != ed25519.SignatureSize {
		return errHostedMFAInvalidProof
	}
	challengeID = strings.TrimSpace(challengeID)
	credentialID = strings.TrimSpace(credentialID)
	if challengeID == "" || credentialID == "" {
		return errHostedMFAInvalidProof
	}

	verificationErr := func() error {
		s.mu.Lock()
		defer s.mu.Unlock()
		now := s.now().UnixMilli()
		sessionIndex := -1
		for i := range s.state.Sessions {
			if hostedMFASessionMatches(s.state.Sessions[i], principal, now) {
				sessionIndex = i
				break
			}
		}
		if sessionIndex < 0 {
			return errHostedMFAInvalidProof
		}
		session := &s.state.Sessions[sessionIndex]
		challenge := session.MFAChallenge
		if challenge == nil || challenge.ID != challengeID || challenge.CredentialID != credentialID {
			return errHostedMFAChallengeMissing
		}

		userIndex := -1
		credentialIndex := -1
		var publicKey ed25519.PublicKey
		for i := range s.state.Users {
			u := &s.state.Users[i]
			if u.ID != principal.UserID || u.Status != UserActive || normalizedTenantID(u.TenantID) != normalizedTenantID(principal.TenantID) {
				continue
			}
			userIndex = i
			for j := range u.MFACredentials {
				credential := u.MFACredentials[j]
				if credential.ID == credentialID && credential.RevokedAt == 0 && credential.Algorithm == hostedMFAAlgorithmEd25519 {
					key, _, decodeErr := decodeHostedMFAPublicKey(credential.PublicKey)
					if decodeErr == nil {
						credentialIndex = j
						publicKey = key
					}
					break
				}
			}
			break
		}
		if userIndex < 0 || credentialIndex < 0 {
			return errHostedMFACredentialMissing
		}

		previousSession := *session
		previousCredentials := append([]HostedMFACredentialRecord(nil), s.state.Users[userIndex].MFACredentials...)
		previousEvents := append([]IdentitySecurityEvent(nil), s.state.SecurityEvents...)
		session.MFAChallenge = nil

		validHash := false
		if now <= challenge.ExpiresAt && challenge.ExpiresAt <= session.AbsoluteExpiresAt {
			sum := sha256.Sum256(payload)
			expected, decodeErr := base64.RawURLEncoding.DecodeString(challenge.ChallengeHash)
			if decodeErr == nil && len(expected) == len(sum) {
				validHash = subtle.ConstantTimeCompare(sum[:], expected) == 1
			}
		}
		validSignature := validHash && ed25519.Verify(publicKey, payload, signature)
		if !validSignature {
			s.appendIdentitySecurityEventLocked(IdentitySecurityMFAVerificationFailed, principal.TenantID, principal.UserID, principal.DeviceID, principal.SessionID, now)
			if err := s.persistLocked(); err != nil {
				*session = previousSession
				s.state.Users[userIndex].MFACredentials = previousCredentials
				s.state.SecurityEvents = previousEvents
				return err
			}
			return errHostedMFAInvalidProof
		}

		s.state.Users[userIndex].MFACredentials[credentialIndex].LastUsedAt = now
		s.state.Users[userIndex].UpdatedAt = now
		s.appendIdentitySecurityEventLocked(IdentitySecurityMFAVerified, principal.TenantID, principal.UserID, principal.DeviceID, principal.SessionID, now)
		if err := s.persistLocked(); err != nil {
			*session = previousSession
			s.state.Users[userIndex].MFACredentials = previousCredentials
			s.state.SecurityEvents = previousEvents
			return err
		}
		return nil
	}()
	if verificationErr != nil {
		return verificationErr
	}
	return s.recordHostedMFAVerification(principal)
}

func (s *IdentityService) revokeHostedMFACredential(principal Principal, credentialID string) error {
	if s == nil {
		return errors.New("identity unavailable")
	}
	credentialID = strings.TrimSpace(credentialID)
	if credentialID == "" {
		return errHostedMFACredentialMissing
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	now := s.now().UnixMilli()
	sessionOK := false
	for i := range s.state.Sessions {
		if hostedMFASessionMatches(s.state.Sessions[i], principal, now) {
			if !hostedMFASessionRecentlyAuthenticated(s.state.Sessions[i], now) {
				return errHostedMFARecentAuthRequired
			}
			sessionOK = true
			break
		}
	}
	if !sessionOK {
		return errors.New("session unavailable for MFA revocation")
	}

	for i := range s.state.Users {
		u := &s.state.Users[i]
		if u.ID != principal.UserID || u.Status != UserActive || normalizedTenantID(u.TenantID) != normalizedTenantID(principal.TenantID) {
			continue
		}
		credentialIndex := -1
		for j := range u.MFACredentials {
			if u.MFACredentials[j].ID == credentialID && u.MFACredentials[j].RevokedAt == 0 {
				credentialIndex = j
				break
			}
		}
		if credentialIndex < 0 {
			return errHostedMFACredentialMissing
		}
		previousCredentials := append([]HostedMFACredentialRecord(nil), u.MFACredentials...)
		previousSessions := append([]SessionRecord(nil), s.state.Sessions...)
		previousEvents := append([]IdentitySecurityEvent(nil), s.state.SecurityEvents...)
		previousUpdatedAt := u.UpdatedAt
		u.MFACredentials[credentialIndex].RevokedAt = now
		u.UpdatedAt = now
		for j := range s.state.Sessions {
			rec := &s.state.Sessions[j]
			if rec.UserID != principal.UserID || normalizedTenantID(rec.TenantID) != normalizedTenantID(principal.TenantID) {
				continue
			}
			if rec.MFAChallenge != nil && rec.MFAChallenge.CredentialID == credentialID {
				rec.MFAChallenge = nil
			}
			rec.MFAVerifiedAt = 0
		}
		s.appendIdentitySecurityEventLocked(IdentitySecurityMFARevoked, principal.TenantID, principal.UserID, principal.DeviceID, principal.SessionID, now)
		if err := s.persistLocked(); err != nil {
			u.MFACredentials = previousCredentials
			u.UpdatedAt = previousUpdatedAt
			s.state.Sessions = previousSessions
			s.state.SecurityEvents = previousEvents
			return err
		}
		return nil
	}
	return errors.New("tenant principal unavailable")
}
