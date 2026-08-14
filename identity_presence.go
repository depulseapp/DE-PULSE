package main

import "strings"

// sessionIDActive is the non-token presence/lifecycle check used by an already
// authenticated long-lived connection. It never creates authority: callers must
// have obtained the session ID from auth middleware first.
func (s *IdentityService) sessionIDActive(sessionID string) bool {
	if strings.TrimSpace(sessionID) == "" {
		return false
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	now := s.now().UnixMilli()
	for i := range s.state.Sessions {
		rec := s.state.Sessions[i]
		if rec.ID == sessionID {
			return rec.RevokedAt == 0 && now < rec.IdleExpiresAt && now < rec.AbsoluteExpiresAt
		}
	}
	return false
}
