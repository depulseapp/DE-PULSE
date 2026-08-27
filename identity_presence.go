package main

import (
	"crypto/ed25519"
	"crypto/rand"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/base64"
	"errors"
	"fmt"
	"strings"
	"time"
)

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
