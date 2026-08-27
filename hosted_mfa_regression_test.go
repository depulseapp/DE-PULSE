package main

import (
	"crypto/ed25519"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func hostedMFATestKey(t *testing.T) (string, ed25519.PrivateKey) {
	t.Helper()
	publicKey, privateKey, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		t.Fatal(err)
	}
	return base64.RawURLEncoding.EncodeToString(publicKey), privateKey
}

func hostedMFAJSONPOST(t *testing.T, path, token, csrf string, payload any) *http.Request {
	t.Helper()
	body, err := json.Marshal(payload)
	if err != nil {
		t.Fatal(err)
	}
	return v184AuthenticatedPOST(path, token, csrf, string(body))
}

func hostedMFASessionDecision(s *IdentityService, p Principal) HostedIdentityDecision {
	return s.authorizeHostedIdentity(p, HostedIdentityRequirement{
		TenantID:                    p.TenantID,
		Capability:                  hostedCapabilityTenantManage,
		RequireRecentAuthentication: true,
		RequireMFA:                  true,
	})
}

func TestHOST007PublicKeyMFACeremonyEnforcesSensitiveAuthorizationAndRejectsReplay(t *testing.T) {
	resetV184LoginLimiter(t)
	_, s := newIdentityTestService(t)
	base := time.Unix(2_450_000_000, 0)
	token, p, _ := v184CredentialedOwner(t, s, base)
	publicKey, privateKey := hostedMFATestKey(t)
	app := &Application{identity: s}
	mux := hostedIdentityHTTPMux(app)
	csrf := "host007-mfa-csrf"

	enrolled := httptest.NewRecorder()
	mux.ServeHTTP(enrolled, hostedMFAJSONPOST(t, "/api/auth/mfa/credential/enroll", token, csrf, map[string]any{
		"label": "primary security key", "algorithm": hostedMFAAlgorithmEd25519, "publicKey": publicKey,
	}))
	if enrolled.Code != http.StatusOK {
		t.Fatalf("MFA enrollment failed: code=%d body=%s", enrolled.Code, enrolled.Body.String())
	}
	if strings.Contains(enrolled.Body.String(), publicKey) {
		t.Fatalf("MFA public key leaked from credential response: %s", enrolled.Body.String())
	}
	var enrollment struct {
		Credential HostedMFACredentialView `json:"credential"`
	}
	if err := json.Unmarshal(enrolled.Body.Bytes(), &enrollment); err != nil || enrollment.Credential.ID == "" {
		t.Fatalf("missing enrolled credential: body=%s err=%v", enrolled.Body.String(), err)
	}

	if decision := hostedMFASessionDecision(s, p); decision.Allowed {
		t.Fatalf("sensitive hosted action allowed before MFA ceremony: %+v", decision)
	}

	challenged := httptest.NewRecorder()
	mux.ServeHTTP(challenged, hostedMFAJSONPOST(t, "/api/auth/mfa/challenge", token, csrf, map[string]any{
		"credentialId": enrollment.Credential.ID,
	}))
	if challenged.Code != http.StatusOK {
		t.Fatalf("MFA challenge failed: code=%d body=%s", challenged.Code, challenged.Body.String())
	}
	var challengeBody struct {
		Challenge HostedMFAChallenge `json:"challenge"`
	}
	if err := json.Unmarshal(challenged.Body.Bytes(), &challengeBody); err != nil {
		t.Fatal(err)
	}
	challenge := challengeBody.Challenge
	payload, err := base64.RawURLEncoding.DecodeString(challenge.SigningPayload)
	if err != nil || challenge.ID == "" || challenge.CredentialID != enrollment.Credential.ID {
		t.Fatalf("invalid MFA challenge contract: %+v err=%v", challenge, err)
	}
	signature := base64.RawURLEncoding.EncodeToString(ed25519.Sign(privateKey, payload))
	verification := map[string]any{
		"challengeId": challenge.ID, "credentialId": challenge.CredentialID,
		"signingPayload": challenge.SigningPayload, "signature": signature,
	}

	missingCSRF := httptest.NewRequest(http.MethodPost, "/api/auth/mfa/verify", strings.NewReader(mustJSONForHOST007(t, verification)))
	missingCSRF.AddCookie(&http.Cookie{Name: sessionCookieName, Value: token})
	missingCSRF.Header.Set("Content-Type", "application/json")
	blocked := httptest.NewRecorder()
	mux.ServeHTTP(blocked, missingCSRF)
	if blocked.Code != http.StatusForbidden {
		t.Fatalf("MFA verify crossed CSRF boundary: code=%d body=%s", blocked.Code, blocked.Body.String())
	}

	verified := httptest.NewRecorder()
	mux.ServeHTTP(verified, hostedMFAJSONPOST(t, "/api/auth/mfa/verify", token, csrf, verification))
	if verified.Code != http.StatusOK {
		t.Fatalf("valid MFA signature rejected: code=%d body=%s", verified.Code, verified.Body.String())
	}
	if decision := hostedMFASessionDecision(s, p); !decision.Allowed {
		t.Fatalf("cryptographically verified MFA proof did not satisfy sensitive authorization: %+v", decision)
	}
	proofAt := v184SessionByID(t, s, p.SessionID).MFAVerifiedAt
	if proofAt != base.UnixMilli() {
		t.Fatalf("unexpected MFA proof time: got=%d want=%d", proofAt, base.UnixMilli())
	}

	replay := httptest.NewRecorder()
	mux.ServeHTTP(replay, hostedMFAJSONPOST(t, "/api/auth/mfa/verify", token, csrf, verification))
	if replay.Code != http.StatusForbidden {
		t.Fatalf("MFA challenge replay was accepted: code=%d body=%s", replay.Code, replay.Body.String())
	}
	if got := v184SessionByID(t, s, p.SessionID).MFAVerifiedAt; got != proofAt {
		t.Fatalf("replay mutated MFA proof timestamp: before=%d after=%d", proofAt, got)
	}
}

func mustJSONForHOST007(t *testing.T, value any) string {
	t.Helper()
	encoded, err := json.Marshal(value)
	if err != nil {
		t.Fatal(err)
	}
	return string(encoded)
}

func TestHOST007InvalidSignatureAndCrossSessionChallengeFailClosed(t *testing.T) {
	resetV184LoginLimiter(t)
	_, s := newIdentityTestService(t)
	base := time.Unix(2_460_000_000, 0)
	_, p1, password := v184CredentialedOwner(t, s, base)
	publicKey, _ := hostedMFATestKey(t)
	credential, err := s.enrollHostedMFACredential(p1, "primary", publicKey)
	if err != nil {
		t.Fatal(err)
	}
	challenge, err := s.createHostedMFAChallenge(p1, credential.ID)
	if err != nil {
		t.Fatal(err)
	}
	payload, err := base64.RawURLEncoding.DecodeString(challenge.SigningPayload)
	if err != nil {
		t.Fatal(err)
	}
	_, wrongPrivateKey := hostedMFATestKey(t)
	wrongSignature := base64.RawURLEncoding.EncodeToString(ed25519.Sign(wrongPrivateKey, payload))

	_, p2, err := s.authenticate("owner", password)
	if err != nil {
		t.Fatal(err)
	}
	if err := s.verifyHostedMFAChallenge(p2, challenge.ID, challenge.CredentialID, challenge.SigningPayload, wrongSignature); !errors.Is(err, errHostedMFAChallengeMissing) {
		t.Fatalf("cross-session challenge was not denied: err=%v", err)
	}
	if decision := hostedMFASessionDecision(s, p2); decision.Allowed {
		t.Fatalf("cross-session attempt created MFA proof: %+v", decision)
	}

	if err := s.verifyHostedMFAChallenge(p1, challenge.ID, challenge.CredentialID, challenge.SigningPayload, wrongSignature); !errors.Is(err, errHostedMFAInvalidProof) {
		t.Fatalf("invalid signature was not rejected: err=%v", err)
	}
	if decision := hostedMFASessionDecision(s, p1); decision.Allowed {
		t.Fatalf("invalid signature created MFA proof: %+v", decision)
	}
	if err := s.verifyHostedMFAChallenge(p1, challenge.ID, challenge.CredentialID, challenge.SigningPayload, wrongSignature); !errors.Is(err, errHostedMFAChallengeMissing) {
		t.Fatalf("failed verification did not consume one-time challenge: err=%v", err)
	}

	events, err := s.adminSecurityEvents(p1)
	if err != nil {
		t.Fatal(err)
	}
	foundFailure := false
	for _, event := range events {
		if event.Type == IdentitySecurityMFAVerificationFailed && event.UserID == p1.UserID && event.SessionID == p1.SessionID {
			foundFailure = true
			break
		}
	}
	if !foundFailure {
		t.Fatal("MFA verification failure was not observable in canonical security events")
	}
}

func TestHOST007And164MFACredentialChallengeProofRotationAndRevocationPersist(t *testing.T) {
	resetV184LoginLimiter(t)
	store, s := newIdentityTestService(t)
	base := time.Unix(2_470_000_000, 0)
	token, p, _ := v184CredentialedOwner(t, s, base)
	publicKey, privateKey := hostedMFATestKey(t)
	credential, err := s.enrollHostedMFACredential(p, "persistent key", publicKey)
	if err != nil {
		t.Fatal(err)
	}
	challenge, err := s.createHostedMFAChallenge(p, credential.ID)
	if err != nil {
		t.Fatal(err)
	}

	reloaded, err := NewIdentityService(store)
	if err != nil {
		t.Fatal(err)
	}
	reloaded.now = func() time.Time { return base }
	credentials, recentMFA, err := reloaded.listHostedMFACredentials(p)
	if err != nil || recentMFA || len(credentials) != 1 || credentials[0].ID != credential.ID {
		t.Fatalf("MFA credential state did not survive restart: credentials=%+v recent=%v err=%v", credentials, recentMFA, err)
	}
	payload, err := base64.RawURLEncoding.DecodeString(challenge.SigningPayload)
	if err != nil {
		t.Fatal(err)
	}
	signature := base64.RawURLEncoding.EncodeToString(ed25519.Sign(privateKey, payload))
	if err := reloaded.verifyHostedMFAChallenge(p, challenge.ID, challenge.CredentialID, challenge.SigningPayload, signature); err != nil {
		t.Fatalf("persisted challenge could not complete after restart: %v", err)
	}
	if decision := hostedMFASessionDecision(reloaded, p); !decision.Allowed {
		t.Fatalf("persisted ceremony did not create canonical MFA assurance: %+v", decision)
	}

	rotatedToken, rotatedPrincipal, err := reloaded.rotate(token)
	if err != nil || rotatedToken == "" || rotatedPrincipal.SessionID == p.SessionID {
		t.Fatalf("session rotation failed after MFA: token=%q principal=%+v err=%v", rotatedToken, rotatedPrincipal, err)
	}
	if decision := hostedMFASessionDecision(reloaded, rotatedPrincipal); !decision.Allowed {
		t.Fatalf("valid MFA assurance did not survive secure session rotation: %+v", decision)
	}

	if err := reloaded.revokeHostedMFACredential(rotatedPrincipal, credential.ID); err != nil {
		t.Fatal(err)
	}
	if decision := hostedMFASessionDecision(reloaded, rotatedPrincipal); decision.Allowed {
		t.Fatalf("revoked MFA credential left session assurance active: %+v", decision)
	}
	if _, err := reloaded.createHostedMFAChallenge(rotatedPrincipal, credential.ID); !errors.Is(err, errHostedMFACredentialMissing) {
		t.Fatalf("revoked MFA credential remained challengeable: err=%v", err)
	}

	persistedAgain, err := NewIdentityService(store)
	if err != nil {
		t.Fatal(err)
	}
	persistedAgain.now = func() time.Time { return base }
	credentials, recentMFA, err = persistedAgain.listHostedMFACredentials(rotatedPrincipal)
	if err != nil || recentMFA || len(credentials) != 1 || credentials[0].RevokedAt == 0 {
		t.Fatalf("MFA revocation did not survive restart: credentials=%+v recent=%v err=%v", credentials, recentMFA, err)
	}
}
