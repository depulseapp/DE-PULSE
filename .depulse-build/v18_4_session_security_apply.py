#!/usr/bin/env python3
from pathlib import Path

p = Path('identity_model.go')
s = p.read_text()

def replace_exact(old, new, label):
    global s
    count = s.count(old)
    if count != 1:
        raise SystemExit(f'{label}: expected exactly one match, got {count}')
    s = s.replace(old, new, 1)

replace_exact('''type SessionRecord struct {
\tID                string `json:"id"`
\tTokenHash         string `json:"tokenHash"`
\tUserID            string `json:"userId"`
\tCreatedAt         int64  `json:"createdAt"`
\tLastSeenAt        int64  `json:"lastSeenAt"`
\tIdleExpiresAt     int64  `json:"idleExpiresAt"`
\tAbsoluteExpiresAt int64  `json:"absoluteExpiresAt"`
\tRevokedAt         int64  `json:"revokedAt,omitempty"`
\tRotatedFrom       string `json:"rotatedFrom,omitempty"`
}''','''type SessionRecord struct {
\tID                string `json:"id"`
\tTokenHash         string `json:"tokenHash"`
\tUserID            string `json:"userId"`
\tCreatedAt         int64  `json:"createdAt"`
\tAuthenticatedAt   int64  `json:"authenticatedAt,omitempty"`
\tLastSeenAt        int64  `json:"lastSeenAt"`
\tIdleExpiresAt     int64  `json:"idleExpiresAt"`
\tAbsoluteExpiresAt int64  `json:"absoluteExpiresAt"`
\tRevokedAt         int64  `json:"revokedAt,omitempty"`
\tRotatedFrom       string `json:"rotatedFrom,omitempty"`
}''','session record')

replace_exact('''const (
\tdefaultSessionIdleTTL     = 2 * time.Hour
\tdefaultSessionAbsoluteTTL = 24 * time.Hour
\tsessionCookieName         = "depulse_session"
\tcsrfCookieName            = "depulse_csrf"
\tbootstrapOwnerID          = "bootstrap-owner"
)''','''const (
\tdefaultSessionIdleTTL     = 2 * time.Hour
\tdefaultSessionAbsoluteTTL = 24 * time.Hour
\tdefaultSensitiveReauthTTL = 15 * time.Minute
\tsessionCookieName         = "depulse_session"
\tcsrfCookieName            = "depulse_csrf"
\tbootstrapOwnerID          = "bootstrap-owner"
)''','session constants')

replace_exact('''func (s *IdentityService) createSessionLocked(u UserRecord, rotatedFrom string) (string, Principal, error) {
\traw := randomID("ses") + randomID("")
\tnow := s.now()
\tcutoff := now.Add(-7 * 24 * time.Hour).UnixMilli()
\tkept := s.state.Sessions[:0]
\tfor _, existing := range s.state.Sessions {
\t\tif existing.RevokedAt > 0 && existing.RevokedAt < cutoff {
\t\t\tcontinue
\t\t}
\t\tif existing.AbsoluteExpiresAt > 0 && existing.AbsoluteExpiresAt < cutoff {
\t\t\tcontinue
\t\t}
\t\tkept = append(kept, existing)
\t}
\ts.state.Sessions = kept
\trec := SessionRecord{ID: randomID("sid"), TokenHash: sessionTokenHash(raw), UserID: u.ID, CreatedAt: now.UnixMilli(), LastSeenAt: now.UnixMilli(), IdleExpiresAt: now.Add(s.idleTTL).UnixMilli(), AbsoluteExpiresAt: now.Add(s.absoluteTTL).UnixMilli(), RotatedFrom: rotatedFrom}
\ts.state.Sessions = append(s.state.Sessions, rec)
\tif err := s.persistLocked(); err != nil {
\t\treturn "", Principal{}, err
\t}
\treturn raw, Principal{UserID: u.ID, Username: u.Username, DisplayName: u.DisplayName, Role: u.Role, SessionID: rec.ID}, nil
}''','''func sessionAuthenticationTime(rec SessionRecord) int64 {
\tif rec.AuthenticatedAt > 0 {
\t\treturn rec.AuthenticatedAt
\t}
\treturn rec.CreatedAt
}

func (s *IdentityService) createSessionLocked(u UserRecord, rotatedFrom string) (string, Principal, error) {
\tnow := s.now()
\treturn s.createSessionWithSecurityLocked(u, rotatedFrom, now.UnixMilli(), now.Add(s.absoluteTTL).UnixMilli())
}

func (s *IdentityService) createSessionWithSecurityLocked(u UserRecord, rotatedFrom string, authenticatedAt, absoluteExpiresAt int64) (string, Principal, error) {
\traw := randomID("ses") + randomID("")
\tnow := s.now()
\tnowMillis := now.UnixMilli()
\tif authenticatedAt <= 0 || authenticatedAt > nowMillis {
\t\tauthenticatedAt = nowMillis
\t}
\tif absoluteExpiresAt <= nowMillis {
\t\treturn "", Principal{}, errors.New("session lifetime exhausted")
\t}
\tcutoff := now.Add(-7 * 24 * time.Hour).UnixMilli()
\tkept := s.state.Sessions[:0]
\tfor _, existing := range s.state.Sessions {
\t\tif existing.RevokedAt > 0 && existing.RevokedAt < cutoff {
\t\t\tcontinue
\t\t}
\t\tif existing.AbsoluteExpiresAt > 0 && existing.AbsoluteExpiresAt < cutoff {
\t\t\tcontinue
\t\t}
\t\tkept = append(kept, existing)
\t}
\ts.state.Sessions = kept
\tidleExpiresAt := now.Add(s.idleTTL).UnixMilli()
\tif idleExpiresAt > absoluteExpiresAt {
\t\tidleExpiresAt = absoluteExpiresAt
\t}
\trec := SessionRecord{ID: randomID("sid"), TokenHash: sessionTokenHash(raw), UserID: u.ID, CreatedAt: nowMillis, AuthenticatedAt: authenticatedAt, LastSeenAt: nowMillis, IdleExpiresAt: idleExpiresAt, AbsoluteExpiresAt: absoluteExpiresAt, RotatedFrom: rotatedFrom}
\ts.state.Sessions = append(s.state.Sessions, rec)
\tif err := s.persistLocked(); err != nil {
\t\treturn "", Principal{}, err
\t}
\treturn raw, Principal{UserID: u.ID, Username: u.Username, DisplayName: u.DisplayName, Role: u.Role, SessionID: rec.ID}, nil
}''','session creation')

replace_exact('''func (s *IdentityService) rotate(token string) (string, Principal, error) {
\tp, err := s.resolve(token, false)
\tif err != nil {
\t\treturn "", Principal{}, err
\t}
\ts.mu.Lock()
\tdefer s.mu.Unlock()
\tnow := s.now().UnixMilli()
\tvar u UserRecord
\tfound := false
\tfor i := range s.state.Users {
\t\tif s.state.Users[i].ID == p.UserID {
\t\t\tu = s.state.Users[i]
\t\t\tfound = true
\t\t\tbreak
\t\t}
\t}
\tif !found {
\t\treturn "", Principal{}, errors.New("principal unavailable")
\t}
\toldHash := sessionTokenHash(token)
\toldID := ""
\tfor i := range s.state.Sessions {
\t\tif subtle.ConstantTimeCompare([]byte(s.state.Sessions[i].TokenHash), []byte(oldHash)) == 1 && s.state.Sessions[i].RevokedAt == 0 {
\t\t\ts.state.Sessions[i].RevokedAt = now
\t\t\toldID = s.state.Sessions[i].ID
\t\t\tbreak
\t\t}
\t}
\treturn s.createSessionLocked(u, oldID)
}''','''func (s *IdentityService) rotate(token string) (string, Principal, error) {
\tp, err := s.resolve(token, false)
\tif err != nil {
\t\treturn "", Principal{}, err
\t}
\ts.mu.Lock()
\tdefer s.mu.Unlock()
\tnow := s.now().UnixMilli()
\tvar u UserRecord
\tfound := false
\tfor i := range s.state.Users {
\t\tif s.state.Users[i].ID == p.UserID {
\t\t\tu = s.state.Users[i]
\t\t\tfound = true
\t\t\tbreak
\t\t}
\t}
\tif !found {
\t\treturn "", Principal{}, errors.New("principal unavailable")
\t}
\toldHash := sessionTokenHash(token)
\toldID := ""
\tauthenticatedAt := int64(0)
\tabsoluteExpiresAt := int64(0)
\tfor i := range s.state.Sessions {
\t\trec := &s.state.Sessions[i]
\t\tif subtle.ConstantTimeCompare([]byte(rec.TokenHash), []byte(oldHash)) == 1 && rec.RevokedAt == 0 {
\t\t\toldID = rec.ID
\t\t\tauthenticatedAt = sessionAuthenticationTime(*rec)
\t\t\tabsoluteExpiresAt = rec.AbsoluteExpiresAt
\t\t\tif absoluteExpiresAt <= now {
\t\t\t\treturn "", Principal{}, errors.New("expired session")
\t\t\t}
\t\t\trec.RevokedAt = now
\t\t\tbreak
\t\t}
\t}
\tif oldID == "" {
\t\treturn "", Principal{}, errors.New("invalid session")
\t}
\treturn s.createSessionWithSecurityLocked(u, oldID, authenticatedAt, absoluteExpiresAt)
}''','session rotation')

replace_exact('''func (s *IdentityService) status(token string) map[string]any {''','''func (s *IdentityService) sessionRecentlyAuthenticated(sessionID string, maxAge time.Duration) bool {
\tif strings.TrimSpace(sessionID) == "" || maxAge <= 0 {
\t\treturn false
\t}
\ts.mu.Lock()
\tdefer s.mu.Unlock()
\tnow := s.now().UnixMilli()
\tfor i := range s.state.Sessions {
\t\trec := s.state.Sessions[i]
\t\tif rec.ID != sessionID || rec.RevokedAt != 0 || now >= rec.IdleExpiresAt || now >= rec.AbsoluteExpiresAt {
\t\t\tcontinue
\t\t}
\t\tauthenticatedAt := sessionAuthenticationTime(rec)
\t\treturn authenticatedAt > 0 && authenticatedAt <= now && now-authenticatedAt <= int64(maxAge/time.Millisecond)
\t}
\treturn false
}

func (s *IdentityService) status(token string) map[string]any {''','recent auth query')

p.write_text(s)
Path('v18_4_session_security_test.go').write_text(Path('.depulse-build/v18_4_session_security_test.go.payload').read_text())
