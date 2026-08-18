package main

import (
	"context"
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

type UserRecord struct {
	ID              string     `json:"id"`
	Username        string     `json:"username"`
	DisplayName     string     `json:"displayName,omitempty"`
	Role            UserRole   `json:"role"`
	Status          UserStatus `json:"status"`
	PasswordHash    string     `json:"passwordHash,omitempty"`
	MustSetPassword bool       `json:"mustSetPassword,omitempty"`
	CreatedAt       int64      `json:"createdAt"`
	UpdatedAt       int64      `json:"updatedAt"`
	LastLoginAt     int64      `json:"lastLoginAt,omitempty"`
}

type SessionRecord struct {
	ID                string `json:"id"`
	TokenHash         string `json:"tokenHash"`
	UserID            string `json:"userId"`
	CreatedAt         int64  `json:"createdAt"`
	AuthenticatedAt   int64  `json:"authenticatedAt,omitempty"`
	LastSeenAt        int64  `json:"lastSeenAt"`
	IdleExpiresAt     int64  `json:"idleExpiresAt"`
	AbsoluteExpiresAt int64  `json:"absoluteExpiresAt"`
	RevokedAt         int64  `json:"revokedAt,omitempty"`
	RotatedFrom       string `json:"rotatedFrom,omitempty"`
}

type IdentityPersistentState struct {
	Version   int             `json:"version"`
	Users     []UserRecord    `json:"users"`
	Sessions  []SessionRecord `json:"sessions"`
	UpdatedAt int64           `json:"updatedAt"`
}

type Principal struct {
	UserID      string   `json:"userId"`
	Username    string   `json:"username"`
	DisplayName string   `json:"displayName,omitempty"`
	Role        UserRole `json:"role"`
	SessionID   string   `json:"sessionId"`
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

const (
	defaultSessionIdleTTL     = 2 * time.Hour
	defaultSessionAbsoluteTTL = 24 * time.Hour
	defaultSensitiveReauthTTL = 15 * time.Minute
	sessionCookieName         = "depulse_session"
	csrfCookieName            = "depulse_csrf"
	bootstrapOwnerID          = "bootstrap-owner"
)

type IdentityService struct {
	mu          sync.Mutex
	persistence *PersistenceManager
	state       IdentityPersistentState
	now         func() time.Time
	idleTTL     time.Duration
	absoluteTTL time.Duration
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
	if err := s.ensureBootstrapOwnerLocked(); err != nil {
		return nil, err
	}
	return s, nil
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
		ID: bootstrapOwnerID, Username: "owner", DisplayName: "Local Owner", Role: RoleOwner,
		Status: UserActive, MustSetPassword: true, CreatedAt: now, UpdatedAt: now,
	})
	return s.persistLocked()
}

func (s *IdentityService) persistLocked() error {
	s.state.Version = 1
	s.state.UpdatedAt = s.now().UnixMilli()
	return s.persistence.SaveIdentityState(context.Background(), s.state)
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

	// Copy only the credential material needed for verification, then release the
	// identity mutex before running Argon2id. Password verification is deliberately
	// expensive and must not block unrelated session resolution.
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

	// Re-resolve after verification so a concurrent password change, disable, role
	// change or user replacement cannot issue a session from stale account state.
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
	// Credential changes revoke every prior session for the account before issuing a new token.
	for i := range s.state.Sessions {
		if s.state.Sessions[i].UserID == userID && s.state.Sessions[i].RevokedAt == 0 {
			s.state.Sessions[i].RevokedAt = now
		}
	}
	return s.createSessionLocked(s.state.Users[idx], "")
}

func (s *IdentityService) updateDisplayName(userID, displayName string) error {
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
			break
		}
	}
	if idx < 0 || s.state.Users[idx].Status != UserActive {
		return errors.New("user unavailable")
	}
	previousName := s.state.Users[idx].DisplayName
	previousUpdatedAt := s.state.Users[idx].UpdatedAt
	s.state.Users[idx].DisplayName = displayName
	s.state.Users[idx].UpdatedAt = s.now().UnixMilli()
	if err := s.persistLocked(); err != nil {
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
	raw := randomID("ses") + randomID("")
	now := s.now()
	nowMillis := now.UnixMilli()
	if authenticatedAt <= 0 || authenticatedAt > nowMillis {
		authenticatedAt = nowMillis
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
	rec := SessionRecord{ID: randomID("sid"), TokenHash: sessionTokenHash(raw), UserID: u.ID, CreatedAt: nowMillis, AuthenticatedAt: authenticatedAt, LastSeenAt: nowMillis, IdleExpiresAt: idleExpiresAt, AbsoluteExpiresAt: absoluteExpiresAt, RotatedFrom: rotatedFrom}
	s.state.Sessions = append(s.state.Sessions, rec)
	if err := s.persistLocked(); err != nil {
		return "", Principal{}, err
	}
	return raw, Principal{UserID: u.ID, Username: u.Username, DisplayName: u.DisplayName, Role: u.Role, SessionID: rec.ID}, nil
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
	if touch {
		// Sliding idle expiry never exceeds absolute expiry and is persisted only when a minute elapsed.
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
	return Principal{UserID: u.ID, Username: u.Username, DisplayName: u.DisplayName, Role: u.Role, SessionID: rec.ID}, nil
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
	for i := range s.state.Sessions {
		rec := &s.state.Sessions[i]
		if subtle.ConstantTimeCompare([]byte(rec.TokenHash), []byte(oldHash)) == 1 && rec.RevokedAt == 0 {
			oldID = rec.ID
			authenticatedAt = sessionAuthenticationTime(*rec)
			absoluteExpiresAt = rec.AbsoluteExpiresAt
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
	return s.createSessionWithSecurityLocked(u, oldID, authenticatedAt, absoluteExpiresAt)
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
