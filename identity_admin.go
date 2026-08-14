package main

import (
	"errors"
	"sort"
	"strings"
	"time"
)

const adminPresenceActiveWindow = 90 * time.Second

type AdminUserView struct {
	ID                 string     `json:"id"`
	Username           string     `json:"username"`
	DisplayName        string     `json:"displayName,omitempty"`
	Role               UserRole   `json:"role"`
	Status             UserStatus `json:"status"`
	MustSetPassword    bool       `json:"mustSetPassword"`
	CreatedAt          int64      `json:"createdAt"`
	UpdatedAt          int64      `json:"updatedAt"`
	LastLoginAt        int64      `json:"lastLoginAt,omitempty"`
	LastSeenAt         int64      `json:"lastSeenAt,omitempty"`
	Presence           string     `json:"presence"`
	ActiveSessionCount int        `json:"activeSessionCount"`
	Manageable         bool       `json:"manageable"`
}

type AdminSessionView struct {
	ID                string `json:"id"`
	UserID            string `json:"userId"`
	Username          string `json:"username"`
	CreatedAt         int64  `json:"createdAt"`
	LastSeenAt        int64  `json:"lastSeenAt"`
	IdleExpiresAt     int64  `json:"idleExpiresAt"`
	AbsoluteExpiresAt int64  `json:"absoluteExpiresAt"`
	RevokedAt         int64  `json:"revokedAt,omitempty"`
	Presence          string `json:"presence"`
	Current           bool   `json:"current"`
	Revokable         bool   `json:"revokable"`
}

func sessionPresence(rec SessionRecord, now int64) string {
	if rec.RevokedAt != 0 || now >= rec.IdleExpiresAt || now >= rec.AbsoluteExpiresAt {
		return "OFFLINE"
	}
	if now-rec.LastSeenAt <= int64(adminPresenceActiveWindow/time.Millisecond) {
		return "ACTIVE"
	}
	return "IDLE"
}

func canManageRole(actor, target UserRole) bool {
	return validRole(actor) && validRole(target) && roleRank(actor) > roleRank(target)
}

func validAdminUsername(v string) bool {
	v = normalizeUsername(v)
	if len(v) < 3 || len(v) > 64 {
		return false
	}
	for _, r := range v {
		if (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') || r == '.' || r == '_' || r == '-' {
			continue
		}
		return false
	}
	return true
}

func (s *IdentityService) adminSnapshot(actor Principal) ([]AdminUserView, []AdminSessionView, error) {
	if !roleAtLeast(actor.Role, RoleAdmin) {
		return nil, nil, errors.New("insufficient role")
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	now := s.now().UnixMilli()
	usersByID := make(map[string]UserRecord, len(s.state.Users))
	for _, u := range s.state.Users {
		usersByID[u.ID] = u
	}
	views := make([]AdminUserView, 0, len(s.state.Users))
	for _, u := range s.state.Users {
		v := AdminUserView{ID: u.ID, Username: u.Username, DisplayName: u.DisplayName, Role: u.Role, Status: u.Status, MustSetPassword: u.MustSetPassword, CreatedAt: u.CreatedAt, UpdatedAt: u.UpdatedAt, LastLoginAt: u.LastLoginAt, Presence: "OFFLINE", Manageable: actor.UserID != u.ID && canManageRole(actor.Role, u.Role)}
		for _, rec := range s.state.Sessions {
			if rec.UserID != u.ID {
				continue
			}
			if rec.LastSeenAt > v.LastSeenAt {
				v.LastSeenAt = rec.LastSeenAt
			}
			p := sessionPresence(rec, now)
			if p != "OFFLINE" {
				v.ActiveSessionCount++
			}
			if p == "ACTIVE" || (p == "IDLE" && v.Presence == "OFFLINE") {
				v.Presence = p
			}
		}
		views = append(views, v)
	}
	sort.Slice(views, func(i, j int) bool {
		if roleRank(views[i].Role) != roleRank(views[j].Role) {
			return roleRank(views[i].Role) > roleRank(views[j].Role)
		}
		return views[i].Username < views[j].Username
	})

	sessions := make([]AdminSessionView, 0, len(s.state.Sessions))
	for _, rec := range s.state.Sessions {
		u, ok := usersByID[rec.UserID]
		if !ok {
			continue
		}
		p := sessionPresence(rec, now)
		if p == "OFFLINE" && rec.RevokedAt == 0 && now >= rec.AbsoluteExpiresAt+int64(7*24*time.Hour/time.Millisecond) {
			continue
		}
		sessions = append(sessions, AdminSessionView{ID: rec.ID, UserID: rec.UserID, Username: u.Username, CreatedAt: rec.CreatedAt, LastSeenAt: rec.LastSeenAt, IdleExpiresAt: rec.IdleExpiresAt, AbsoluteExpiresAt: rec.AbsoluteExpiresAt, RevokedAt: rec.RevokedAt, Presence: p, Current: rec.ID == actor.SessionID, Revokable: rec.RevokedAt == 0 && rec.ID != actor.SessionID && canManageRole(actor.Role, u.Role)})
	}
	sort.Slice(sessions, func(i, j int) bool { return sessions[i].LastSeenAt > sessions[j].LastSeenAt })
	return views, sessions, nil
}

func (s *IdentityService) adminCreateUser(actor Principal, username, displayName string, role UserRole, temporaryPassword string) (AdminUserView, error) {
	username = normalizeUsername(username)
	displayName = strings.TrimSpace(displayName)
	if !roleAtLeast(actor.Role, RoleAdmin) || !validAdminUsername(username) {
		return AdminUserView{}, errors.New("invalid user request")
	}
	if role == RoleSuperOwner || !canManageRole(actor.Role, role) {
		return AdminUserView{}, errors.New("role is outside actor authority")
	}
	hash, err := hashPasswordArgon2id(temporaryPassword)
	if err != nil {
		return AdminUserView{}, err
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, u := range s.state.Users {
		if normalizeUsername(u.Username) == username {
			return AdminUserView{}, errors.New("username already exists")
		}
	}
	now := s.now().UnixMilli()
	u := UserRecord{ID: randomID("usr"), Username: username, DisplayName: displayName, Role: role, Status: UserActive, PasswordHash: hash, MustSetPassword: true, CreatedAt: now, UpdatedAt: now}
	s.state.Users = append(s.state.Users, u)
	if err := s.persistLocked(); err != nil {
		s.state.Users = s.state.Users[:len(s.state.Users)-1]
		return AdminUserView{}, err
	}
	return AdminUserView{ID: u.ID, Username: u.Username, DisplayName: u.DisplayName, Role: u.Role, Status: u.Status, MustSetPassword: true, CreatedAt: u.CreatedAt, UpdatedAt: u.UpdatedAt, Presence: "OFFLINE", Manageable: true}, nil
}

func (s *IdentityService) activeCriticalOwnersLocked(exceptID string, exceptRole UserRole, exceptStatus UserStatus) int {
	count := 0
	for _, u := range s.state.Users {
		role, status := u.Role, u.Status
		if u.ID == exceptID {
			role, status = exceptRole, exceptStatus
		}
		if status == UserActive && (role == RoleOwner || role == RoleSuperOwner) {
			count++
		}
	}
	return count
}

func (s *IdentityService) adminSetUserRole(actor Principal, userID string, role UserRole) error {
	if !validRole(role) || role == RoleSuperOwner || strings.TrimSpace(userID) == "" || actor.UserID == userID {
		return errors.New("invalid role change")
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
	if idx < 0 || !canManageRole(actor.Role, s.state.Users[idx].Role) || !canManageRole(actor.Role, role) {
		return errors.New("role change outside actor authority")
	}
	if s.activeCriticalOwnersLocked(userID, role, s.state.Users[idx].Status) == 0 {
		return errors.New("at least one active owner is required")
	}
	now := s.now().UnixMilli()
	s.state.Users[idx].Role = role
	s.state.Users[idx].UpdatedAt = now
	for i := range s.state.Sessions {
		if s.state.Sessions[i].UserID == userID && s.state.Sessions[i].RevokedAt == 0 {
			s.state.Sessions[i].RevokedAt = now
		}
	}
	return s.persistLocked()
}

func (s *IdentityService) adminSetUserStatus(actor Principal, userID string, status UserStatus) error {
	if status != UserActive && status != UserDisabled {
		return errors.New("invalid user status")
	}
	if strings.TrimSpace(userID) == "" || actor.UserID == userID {
		return errors.New("cannot change own account status")
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
	if idx < 0 || !canManageRole(actor.Role, s.state.Users[idx].Role) {
		return errors.New("status change outside actor authority")
	}
	if s.activeCriticalOwnersLocked(userID, s.state.Users[idx].Role, status) == 0 {
		return errors.New("at least one active owner is required")
	}
	now := s.now().UnixMilli()
	s.state.Users[idx].Status = status
	s.state.Users[idx].UpdatedAt = now
	if status == UserDisabled {
		for i := range s.state.Sessions {
			if s.state.Sessions[i].UserID == userID && s.state.Sessions[i].RevokedAt == 0 {
				s.state.Sessions[i].RevokedAt = now
			}
		}
	}
	return s.persistLocked()
}

func (s *IdentityService) adminResetPassword(actor Principal, userID, temporaryPassword string) error {
	if strings.TrimSpace(userID) == "" || actor.UserID == userID {
		return errors.New("invalid password reset target")
	}
	hash, err := hashPasswordArgon2id(temporaryPassword)
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
	if idx < 0 || !canManageRole(actor.Role, s.state.Users[idx].Role) || s.state.Users[idx].Status != UserActive {
		return errors.New("password reset outside actor authority")
	}
	now := s.now().UnixMilli()
	s.state.Users[idx].PasswordHash = hash
	s.state.Users[idx].MustSetPassword = true
	s.state.Users[idx].UpdatedAt = now
	for i := range s.state.Sessions {
		if s.state.Sessions[i].UserID == userID && s.state.Sessions[i].RevokedAt == 0 {
			s.state.Sessions[i].RevokedAt = now
		}
	}
	return s.persistLocked()
}

func (s *IdentityService) adminRevokeSession(actor Principal, sessionID string) error {
	if strings.TrimSpace(sessionID) == "" || sessionID == actor.SessionID {
		return errors.New("cannot revoke current session here")
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	now := s.now().UnixMilli()
	for i := range s.state.Sessions {
		rec := &s.state.Sessions[i]
		if rec.ID != sessionID {
			continue
		}
		var target *UserRecord
		for j := range s.state.Users {
			if s.state.Users[j].ID == rec.UserID {
				target = &s.state.Users[j]
				break
			}
		}
		if target == nil || !canManageRole(actor.Role, target.Role) {
			return errors.New("session revoke outside actor authority")
		}
		if rec.RevokedAt == 0 {
			rec.RevokedAt = now
			return s.persistLocked()
		}
		return nil
	}
	return errors.New("session not found")
}

func (s *IdentityService) touchSessionID(sessionID string) {
	if strings.TrimSpace(sessionID) == "" {
		return
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	now := s.now().UnixMilli()
	for i := range s.state.Sessions {
		rec := &s.state.Sessions[i]
		if rec.ID != sessionID || rec.RevokedAt != 0 || now >= rec.AbsoluteExpiresAt || now >= rec.IdleExpiresAt {
			continue
		}
		if now-rec.LastSeenAt < int64(time.Minute/time.Millisecond) {
			return
		}
		rec.LastSeenAt = now
		next := s.now().Add(s.idleTTL).UnixMilli()
		if next > rec.AbsoluteExpiresAt {
			next = rec.AbsoluteExpiresAt
		}
		rec.IdleExpiresAt = next
		_ = s.persistLocked()
		return
	}
}
