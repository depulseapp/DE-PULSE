package main

import (
	"context"
	"crypto/subtle"
	"encoding/json"
	"net/http"
	"strings"
	"time"
)

func sessionTokenFromRequest(r *http.Request) string {
	c, err := r.Cookie(sessionCookieName)
	if err != nil {
		return ""
	}
	return strings.TrimSpace(c.Value)
}

func setSessionCookie(w http.ResponseWriter, r *http.Request, token string) {
	secure := requestIsSecure(r)
	http.SetCookie(w, &http.Cookie{Name: sessionCookieName, Value: token, Path: "/", HttpOnly: true, SameSite: http.SameSiteStrictMode, Secure: secure, MaxAge: int(defaultSessionAbsoluteTTL / time.Second)})
	// Double-submit CSRF token: readable by same-origin JS, never accepted from a form body.
	csrf := randomID("csrf") + randomID("")
	http.SetCookie(w, &http.Cookie{Name: csrfCookieName, Value: csrf, Path: "/", HttpOnly: false, SameSite: http.SameSiteStrictMode, Secure: secure, MaxAge: int(defaultSessionAbsoluteTTL / time.Second)})
}
func clearSessionCookie(w http.ResponseWriter, r *http.Request) {
	secure := requestIsSecure(r)
	http.SetCookie(w, &http.Cookie{Name: sessionCookieName, Value: "", Path: "/", HttpOnly: true, SameSite: http.SameSiteStrictMode, Secure: secure, MaxAge: -1, Expires: time.Unix(1, 0)})
	http.SetCookie(w, &http.Cookie{Name: csrfCookieName, Value: "", Path: "/", HttpOnly: false, SameSite: http.SameSiteStrictMode, Secure: secure, MaxAge: -1, Expires: time.Unix(1, 0)})
}

func validRequestCSRF(r *http.Request) bool {
	switch r.Method {
	case http.MethodGet, http.MethodHead, http.MethodOptions:
		return true
	}
	cookie, err := r.Cookie(csrfCookieName)
	if err != nil || strings.TrimSpace(cookie.Value) == "" {
		return false
	}
	header := strings.TrimSpace(r.Header.Get("X-DE-PULSE-CSRF"))
	if header == "" || len(header) != len(cookie.Value) {
		return false
	}
	return subtle.ConstantTimeCompare([]byte(header), []byte(cookie.Value)) == 1
}

func decodeJSONBody(w http.ResponseWriter, r *http.Request, dst any) error {
	return json.NewDecoder(http.MaxBytesReader(w, r.Body, 8<<10)).Decode(dst)
}

func (a *Application) handleAuthStatus(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeError(w, 405, "Method not allowed")
		return
	}
	if a.identity == nil {
		writeJSON(w, 200, map[string]any{"authenticated": false, "bootstrapRequired": false})
		return
	}
	writeJSON(w, 200, a.identity.status(sessionTokenFromRequest(r)))
}

func (a *Application) handleLogin(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, 405, "Method not allowed")
		return
	}
	var req struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	if err := json.NewDecoder(http.MaxBytesReader(w, r.Body, 8<<10)).Decode(&req); err != nil {
		writeError(w, 400, "Invalid login request.")
		return
	}
	abuseKey := loginAbuseKey(req.Username, r)
	if allowed, retryAfter := loginLimiter.Allow(abuseKey); !allowed {
		rejectThrottledLogin(w, retryAfter)
		return
	}
	token, p, err := a.identity.authenticate(req.Username, req.Password)
	if err != nil {
		loginLimiter.RecordFailure(abuseKey)
		time.Sleep(120 * time.Millisecond)
		writeError(w, http.StatusUnauthorized, "Invalid username or password.")
		return
	}
	loginLimiter.Reset(abuseKey)
	setSessionCookie(w, r, token)
	writeJSON(w, 200, map[string]any{"ok": true, "principal": p})
}

func (a *Application) handleLogout(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, 405, "Method not allowed")
		return
	}
	token := sessionTokenFromRequest(r)
	if token != "" {
		_ = a.identity.revokeToken(token)
	}
	clearSessionCookie(w, r)
	writeJSON(w, 200, map[string]any{"ok": true})
}

func (a *Application) handleRotateSession(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}
	token := sessionTokenFromRequest(r)
	rotated, p, err := a.identity.rotate(token)
	if err != nil {
		clearSessionCookie(w, r)
		writeError(w, http.StatusUnauthorized, "Authentication required.")
		return
	}
	setSessionCookie(w, r, rotated)
	writeJSON(w, http.StatusOK, map[string]any{"ok": true, "principal": p})
}

func (a *Application) handleSetPassword(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, 405, "Method not allowed")
		return
	}
	p, ok := principalFromContext(r.Context())
	if !ok {
		writeError(w, http.StatusUnauthorized, "Authentication required.")
		return
	}
	var req struct {
		Password string `json:"password"`
	}
	if err := json.NewDecoder(http.MaxBytesReader(w, r.Body, 8<<10)).Decode(&req); err != nil {
		writeError(w, 400, "Invalid password request.")
		return
	}
	token, np, err := a.identity.setPassword(p.UserID, req.Password)
	if err != nil {
		writeError(w, 400, err.Error())
		return
	}
	setSessionCookie(w, r, token)
	writeJSON(w, 200, map[string]any{"ok": true, "principal": np})
}

func (a *Application) handleUpdateProfile(w http.ResponseWriter, r *http.Request) {
	p, ok := principalFromContext(r.Context())
	if !ok {
		writeError(w, http.StatusUnauthorized, "Authentication required.")
		return
	}
	var req struct {
		Username    string `json:"username"`
		DisplayName string `json:"displayName"`
	}
	if err := json.NewDecoder(http.MaxBytesReader(w, r.Body, 8<<10)).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "Invalid profile request.")
		return
	}
	if err := a.identity.updateProfile(p.UserID, req.Username, req.DisplayName); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	np, err := a.identity.resolve(sessionTokenFromRequest(r), false)
	if err != nil {
		writeError(w, http.StatusUnauthorized, "Authentication required.")
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"ok": true, "principal": np})
}

func (a *Application) authResolved(allowPasswordSetup, enforceProductAccess bool, next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if a.identity == nil {
			c, err := r.Cookie("pmt_session")
			if err != nil || a.sessionKey == "" || c.Value != a.sessionKey {
				writeError(w, http.StatusForbidden, "Invalid local app session.")
				return
			}
			next(w, r)
			return
		}
		p, err := a.identity.resolve(sessionTokenFromRequest(r), true)
		if err != nil {
			clearSessionCookie(w, r)
			writeError(w, http.StatusUnauthorized, "Authentication required.")
			return
		}
		if !allowPasswordSetup && a.identity.userRequiresPassword(p.UserID) {
			writeError(w, http.StatusPreconditionRequired, "Owner password setup required.")
			return
		}
		if !validRequestCSRF(r) {
			writeError(w, http.StatusForbidden, "Security token validation failed.")
			return
		}
		if enforceProductAccess {
			if !a.enforceHostedProductAccess(w, r, p) {
				return
			}
			if !a.enforceHostedRequestQuota(w, r, p) {
				return
			}
			if !a.consumeHostedProductRequest(w, r, p) {
				return
			}
		}
		if err := a.ensureUserWorkspace(p.UserID); err != nil {
			writeError(w, http.StatusInternalServerError, "User workspace unavailable.")
			return
		}
		next(w, r.WithContext(context.WithValue(r.Context(), identityContextKey{}, p)))
	}
}
func (a *Application) auth(next http.HandlerFunc) http.HandlerFunc {
	return a.authResolved(false, true, next)
}
func (a *Application) authAllowPasswordSetup(next http.HandlerFunc) http.HandlerFunc {
	return a.authResolved(true, true, next)
}
func (a *Application) authAccountLifecycle(next http.HandlerFunc) http.HandlerFunc {
	return a.authResolved(false, false, next)
}

func (a *Application) requireRole(required UserRole, next http.HandlerFunc) http.HandlerFunc {
	return a.auth(func(w http.ResponseWriter, r *http.Request) {
		if a.identity == nil {
			next(w, r)
			return
		}
		p, ok := principalFromContext(r.Context())
		if !ok || !roleAtLeast(p.Role, required) {
			writeError(w, http.StatusForbidden, "Insufficient role for this operation.")
			return
		}
		next(w, r)
	})
}

func (a *Application) handleReauthenticate(w http.ResponseWriter, r *http.Request) {
	p, ok := principalFromContext(r.Context())
	if !ok || a.identity == nil {
		writeError(w, http.StatusUnauthorized, "Authentication required.")
		return
	}
	var req struct {
		Password string `json:"password"`
	}
	if err := json.NewDecoder(http.MaxBytesReader(w, r.Body, 8<<10)).Decode(&req); err != nil || strings.TrimSpace(req.Password) == "" {
		writeError(w, http.StatusBadRequest, "Current password is required.")
		return
	}
	abuseKey := "reauth:" + loginAbuseKey(p.Username, r)
	if allowed, retryAfter := loginLimiter.Allow(abuseKey); !allowed {
		rejectThrottledLogin(w, retryAfter)
		return
	}
	verified, err := a.identity.reauthenticateSession(p.SessionID, req.Password)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Password verification unavailable.")
		return
	}
	if !verified {
		loginLimiter.RecordFailure(abuseKey)
		time.Sleep(120 * time.Millisecond)
		writeError(w, http.StatusForbidden, "Password verification failed.")
		return
	}
	loginLimiter.Reset(abuseKey)
	writeJSON(w, http.StatusOK, map[string]any{"ok": true, "recentAuthentication": true})
}

func writeReauthRequired(w http.ResponseWriter) {
	writeJSON(w, http.StatusPreconditionRequired, map[string]any{
		"error": "Recent password authentication is required for this security-sensitive change.",
		"code":  "REAUTH_REQUIRED",
	})
}

func (a *Application) requireRecentAuthentication(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if a.identity == nil {
			next(w, r)
			return
		}
		p, ok := principalFromContext(r.Context())
		if !ok {
			writeError(w, http.StatusUnauthorized, "Authentication required.")
			return
		}
		if !a.identity.sessionRecentlyAuthenticated(p.SessionID, defaultSensitiveReauthTTL) {
			writeReauthRequired(w)
			return
		}
		next(w, r)
	}
}

const (
	maxHostedDeviceLabelLength       = 128
	maxHostedDeviceFingerprintLength = 512
)

func (a *Application) registerHostedIdentityRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/api/auth/device/register", a.auth(a.requireRecentAuthentication(postOnly(a.handleRegisterHostedDevice))))
	mux.HandleFunc("/api/auth/device/status", a.auth(a.requireRecentAuthentication(postOnly(a.handleHostedDeviceStatus))))
	mux.HandleFunc("/api/auth/account/export", a.authAccountLifecycle(a.requireRecentAuthentication(a.handleAccountDataExport)))
	mux.HandleFunc("/api/auth/account", a.authAccountLifecycle(a.requireRecentAuthentication(a.handleDeleteAccount)))
}

func (a *Application) authorizeHostedDeviceManagement(p Principal, requireRegisteredDevice bool) HostedIdentityDecision {
	if a.identity == nil {
		return HostedIdentityDecision{Allowed: false, BlockingReasons: []string{"identity unavailable"}}
	}
	return a.identity.authorizeHostedIdentity(p, HostedIdentityRequirement{
		TenantID:                    p.TenantID,
		Capability:                  hostedCapabilityDeviceManage,
		RequireRegisteredDevice:     requireRegisteredDevice,
		RequireRecentAuthentication: true,
	})
}

func (a *Application) handleRegisterHostedDevice(w http.ResponseWriter, r *http.Request) {
	p, ok := principalFromContext(r.Context())
	if !ok || a.identity == nil {
		writeError(w, http.StatusUnauthorized, "Authentication required.")
		return
	}
	var req struct {
		Label           string `json:"label"`
		FingerprintHash string `json:"fingerprintHash"`
	}
	if err := json.NewDecoder(http.MaxBytesReader(w, r.Body, 8<<10)).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "Invalid device registration request.")
		return
	}
	req.Label = strings.TrimSpace(req.Label)
	req.FingerprintHash = strings.TrimSpace(req.FingerprintHash)
	if req.FingerprintHash == "" || len(req.FingerprintHash) > maxHostedDeviceFingerprintLength || len(req.Label) > maxHostedDeviceLabelLength {
		writeError(w, http.StatusBadRequest, "Invalid device registration request.")
		return
	}
	if decision := a.authorizeHostedDeviceManagement(p, false); !decision.Allowed {
		writeError(w, http.StatusForbidden, "Hosted identity authorization failed.")
		return
	}
	device, err := a.identity.registerHostedDevice(p, req.Label, req.FingerprintHash)
	if err != nil {
		writeError(w, http.StatusBadRequest, "Device registration failed.")
		return
	}
	if err := a.identity.bindHostedDeviceToSession(p, device.ID); err != nil {
		writeError(w, http.StatusConflict, "Device could not be bound to the current session.")
		return
	}
	decision := a.authorizeHostedDeviceManagement(p, true)
	if !decision.Allowed || decision.DeviceID != device.ID {
		// Binding succeeded but the resulting trust decision failed. Revoke the
		// device and every bound session rather than leave partially trusted state.
		_ = a.identity.setHostedDeviceStatus(p, device.ID, DeviceRevoked)
		clearSessionCookie(w, r)
		writeError(w, http.StatusForbidden, "Hosted device trust verification failed.")
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{
		"ok": true,
		"device": map[string]any{
			"id":     device.ID,
			"label":  device.Label,
			"status": device.Status,
		},
		"sessionBound": true,
	})
}

func (a *Application) handleHostedDeviceStatus(w http.ResponseWriter, r *http.Request) {
	p, ok := principalFromContext(r.Context())
	if !ok || a.identity == nil {
		writeError(w, http.StatusUnauthorized, "Authentication required.")
		return
	}
	var req struct {
		DeviceID string       `json:"deviceId"`
		Status   DeviceStatus `json:"status"`
	}
	if err := json.NewDecoder(http.MaxBytesReader(w, r.Body, 8<<10)).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "Invalid device status request.")
		return
	}
	req.DeviceID = strings.TrimSpace(req.DeviceID)
	req.Status = DeviceStatus(strings.ToUpper(strings.TrimSpace(string(req.Status))))
	if req.DeviceID == "" || (req.Status != DeviceLost && req.Status != DeviceRevoked) {
		// Reactivation is intentionally not exposed here. A lost/revoked device
		// must establish a fresh registration instead of mutating back to ACTIVE.
		writeError(w, http.StatusBadRequest, "Device status must be LOST or REVOKED.")
		return
	}
	decision := a.authorizeHostedDeviceManagement(p, true)
	if !decision.Allowed {
		writeError(w, http.StatusForbidden, "A trusted registered device is required for this operation.")
		return
	}
	if err := a.identity.setHostedDeviceStatus(p, req.DeviceID, req.Status); err != nil {
		writeError(w, http.StatusForbidden, "Device status update denied.")
		return
	}
	currentSessionRevoked := decision.DeviceID == req.DeviceID
	if currentSessionRevoked {
		clearSessionCookie(w, r)
	}
	writeJSON(w, http.StatusOK, map[string]any{
		"ok":                    true,
		"deviceId":              req.DeviceID,
		"status":                req.Status,
		"currentSessionRevoked": currentSessionRevoked,
	})
}
