package main

import (
	"context"
	"errors"
	"net/http"
	"strings"
	"time"
)

func (a *Application) registerHealthRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/api/health", a.handleHealth)
	mux.HandleFunc("/api/ready", a.handleReady)
	a.registerHostedIdentityRoutes(mux)
	a.registerHostedMFARoutes(mux)
	a.registerHostedProductRoutes(mux)
}

func (a *Application) handleHealth(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, map[string]any{
		"ok": true, "version": appVersion, "buildId": buildID, "runtimeMode": runtimeMode(),
	})
}

func (a *Application) handleReady(w http.ResponseWriter, r *http.Request) {
	diagnostics := PersistenceDiagnostics{Backend: "disabled"}
	if a.persistence != nil {
		probeCtx, cancel := context.WithTimeout(r.Context(), 750*time.Millisecond)
		a.persistence.ProbeReady(probeCtx)
		cancel()
		diagnostics = a.persistence.Diagnostics()
	}
	ready := a.identity != nil && diagnostics.Ready
	status := http.StatusOK
	if !ready {
		status = http.StatusServiceUnavailable
	}
	writeJSON(w, status, map[string]any{
		"ok": ready, "version": appVersion, "buildId": buildID, "runtimeMode": runtimeMode(),
		"persistence": map[string]any{"backend": diagnostics.Backend, "ready": diagnostics.Ready, "healthState": diagnostics.HealthState, "lastHealthyAt": diagnostics.LastHealthyAt, "retryScheduled": diagnostics.RetryScheduled, "retryBackoffMs": diagnostics.RetryBackoffMs, "pool": diagnostics.Pool, "database": diagnostics.Database},
	})
}

func (a *Application) registerHostedProductRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/api/auth/product/status", a.auth(a.handleHostedProductStatus))
}

func (a *Application) handleHostedProductStatus(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}
	p, ok := principalFromContext(r.Context())
	if !ok || a.identity == nil {
		writeError(w, http.StatusUnauthorized, "Authentication required.")
		return
	}
	snapshot, err := a.identity.productEntitlementSnapshot(p.TenantID)
	if err != nil {
		writeError(w, http.StatusServiceUnavailable, "Product entitlement unavailable.")
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"ok": true, "product": snapshot})
}

func (a *Application) registerHostedMFARoutes(mux *http.ServeMux) {
	mux.HandleFunc("/api/auth/mfa/credentials", a.authAccountLifecycle(a.handleHostedMFACredentials))
	mux.HandleFunc("/api/auth/mfa/credential/enroll", a.authAccountLifecycle(a.requireRecentAuthentication(postOnly(a.handleEnrollHostedMFACredential))))
	mux.HandleFunc("/api/auth/mfa/credential/revoke", a.authAccountLifecycle(a.requireRecentAuthentication(postOnly(a.handleRevokeHostedMFACredential))))
	mux.HandleFunc("/api/auth/mfa/challenge", a.authAccountLifecycle(postOnly(a.handleCreateHostedMFAChallenge)))
	mux.HandleFunc("/api/auth/mfa/verify", a.authAccountLifecycle(postOnly(a.handleVerifyHostedMFAChallenge)))
}

func (a *Application) handleHostedMFACredentials(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}
	p, ok := principalFromContext(r.Context())
	if !ok || a.identity == nil {
		writeError(w, http.StatusUnauthorized, "Authentication required.")
		return
	}
	credentials, verified, err := a.identity.listHostedMFACredentials(p)
	if err != nil {
		writeError(w, http.StatusServiceUnavailable, "MFA status unavailable.")
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"ok": true, "credentials": credentials, "recentMfa": verified})
}

func (a *Application) handleEnrollHostedMFACredential(w http.ResponseWriter, r *http.Request) {
	p, ok := principalFromContext(r.Context())
	if !ok || a.identity == nil {
		writeError(w, http.StatusUnauthorized, "Authentication required.")
		return
	}
	var req struct {
		Label     string `json:"label"`
		Algorithm string `json:"algorithm"`
		PublicKey string `json:"publicKey"`
	}
	if err := decodeJSONBody(w, r, &req); err != nil || (strings.TrimSpace(req.Algorithm) != "" && !strings.EqualFold(strings.TrimSpace(req.Algorithm), hostedMFAAlgorithmEd25519)) {
		writeError(w, http.StatusBadRequest, "Invalid MFA credential request.")
		return
	}
	credential, err := a.identity.enrollHostedMFACredential(p, req.Label, req.PublicKey)
	if err != nil {
		if errors.Is(err, errHostedMFARecentAuthRequired) {
			writeReauthRequired(w)
			return
		}
		writeError(w, http.StatusBadRequest, "MFA credential enrollment failed.")
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"ok": true, "credential": hostedMFACredentialView(credential)})
}

func (a *Application) handleCreateHostedMFAChallenge(w http.ResponseWriter, r *http.Request) {
	p, ok := principalFromContext(r.Context())
	if !ok || a.identity == nil {
		writeError(w, http.StatusUnauthorized, "Authentication required.")
		return
	}
	var req struct {
		CredentialID string `json:"credentialId"`
	}
	if err := decodeJSONBody(w, r, &req); err != nil {
		writeError(w, http.StatusBadRequest, "Invalid MFA challenge request.")
		return
	}
	challenge, err := a.identity.createHostedMFAChallenge(p, req.CredentialID)
	if err != nil {
		if errors.Is(err, errHostedMFACredentialMissing) {
			writeError(w, http.StatusConflict, "No active MFA credential is available.")
			return
		}
		writeError(w, http.StatusServiceUnavailable, "MFA challenge unavailable.")
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"ok": true, "challenge": challenge})
}

func (a *Application) handleVerifyHostedMFAChallenge(w http.ResponseWriter, r *http.Request) {
	p, ok := principalFromContext(r.Context())
	if !ok || a.identity == nil {
		writeError(w, http.StatusUnauthorized, "Authentication required.")
		return
	}
	var req struct {
		ChallengeID    string `json:"challengeId"`
		CredentialID   string `json:"credentialId"`
		SigningPayload string `json:"signingPayload"`
		Signature      string `json:"signature"`
	}
	if err := decodeJSONBody(w, r, &req); err != nil {
		writeError(w, http.StatusBadRequest, "Invalid MFA verification request.")
		return
	}
	abuseKey := "mfa:" + loginAbuseKey(p.Username, r)
	if allowed, retryAfter := loginLimiter.Allow(abuseKey); !allowed {
		rejectThrottledLogin(w, retryAfter)
		return
	}
	if err := a.identity.verifyHostedMFAChallenge(p, req.ChallengeID, req.CredentialID, req.SigningPayload, req.Signature); err != nil {
		if errors.Is(err, errHostedMFAInvalidProof) || errors.Is(err, errHostedMFAChallengeMissing) || errors.Is(err, errHostedMFACredentialMissing) {
			loginLimiter.RecordFailure(abuseKey)
			writeError(w, http.StatusForbidden, "MFA verification failed.")
			return
		}
		writeError(w, http.StatusServiceUnavailable, "MFA verification unavailable.")
		return
	}
	loginLimiter.Reset(abuseKey)
	writeJSON(w, http.StatusOK, map[string]any{"ok": true, "recentMfa": true})
}

func (a *Application) handleRevokeHostedMFACredential(w http.ResponseWriter, r *http.Request) {
	p, ok := principalFromContext(r.Context())
	if !ok || a.identity == nil {
		writeError(w, http.StatusUnauthorized, "Authentication required.")
		return
	}
	var req struct {
		CredentialID string `json:"credentialId"`
	}
	if err := decodeJSONBody(w, r, &req); err != nil || strings.TrimSpace(req.CredentialID) == "" {
		writeError(w, http.StatusBadRequest, "Invalid MFA credential revocation request.")
		return
	}
	if err := a.identity.revokeHostedMFACredential(p, req.CredentialID); err != nil {
		if errors.Is(err, errHostedMFARecentAuthRequired) {
			writeReauthRequired(w)
			return
		}
		if errors.Is(err, errHostedMFACredentialMissing) {
			writeError(w, http.StatusNotFound, "MFA credential not found.")
			return
		}
		writeError(w, http.StatusServiceUnavailable, "MFA credential revocation unavailable.")
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"ok": true, "credentialId": strings.TrimSpace(req.CredentialID), "recentMfa": false})
}
