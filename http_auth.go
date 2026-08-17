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

func (a *Application) authResolved(allowPasswordSetup bool, next http.HandlerFunc) http.HandlerFunc {
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
		if err := a.ensureUserWorkspace(p.UserID); err != nil {
			writeError(w, http.StatusInternalServerError, "User workspace unavailable.")
			return
		}
		next(w, r.WithContext(context.WithValue(r.Context(), identityContextKey{}, p)))
	}
}
func (a *Application) auth(next http.HandlerFunc) http.HandlerFunc {
	return a.authResolved(false, next)
}
func (a *Application) authAllowPasswordSetup(next http.HandlerFunc) http.HandlerFunc {
	return a.authResolved(true, next)
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
