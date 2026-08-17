package main

import (
	"errors"
	"strings"
)

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

var errRecentAuthenticationRequired = errors.New("recent password authentication required")
