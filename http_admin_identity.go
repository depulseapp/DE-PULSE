package main

import (
	"encoding/json"
	"net/http"
	"strings"
)

func adminIdentityErrorStatus(err error) int {
	if err == nil {
		return http.StatusOK
	}
	msg := strings.ToLower(err.Error())
	if strings.Contains(msg, "authority") || strings.Contains(msg, "insufficient") || strings.Contains(msg, "own account") || strings.Contains(msg, "current session") {
		return http.StatusForbidden
	}
	if strings.Contains(msg, "not found") {
		return http.StatusNotFound
	}
	return http.StatusBadRequest
}

func (a *Application) handleAdminIdentity(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}
	p, ok := principalFromContext(r.Context())
	if !ok {
		writeError(w, http.StatusUnauthorized, "Authentication required.")
		return
	}
	users, sessions, err := a.identity.adminSnapshot(p)
	if err != nil {
		writeError(w, adminIdentityErrorStatus(err), err.Error())
		return
	}
	securityEvents, err := a.identity.adminSecurityEvents(p)
	if err != nil {
		writeError(w, adminIdentityErrorStatus(err), err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"users": users, "sessions": sessions, "securityEvents": securityEvents, "presence": map[string]any{"activeWindowSeconds": int(adminPresenceActiveWindow.Seconds()), "source": "authenticated sessions", "states": []string{"ACTIVE", "IDLE", "OFFLINE"}}})
}

func (a *Application) handleAdminUserCreate(w http.ResponseWriter, r *http.Request) {
	p, _ := principalFromContext(r.Context())
	var req struct {
		Username          string   `json:"username"`
		DisplayName       string   `json:"displayName"`
		Role              UserRole `json:"role"`
		TemporaryPassword string   `json:"temporaryPassword"`
	}
	if err := json.NewDecoder(http.MaxBytesReader(w, r.Body, 16<<10)).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "Invalid user request.")
		return
	}
	user, err := a.identity.adminCreateUser(p, req.Username, req.DisplayName, req.Role, req.TemporaryPassword)
	if err != nil {
		writeError(w, adminIdentityErrorStatus(err), err.Error())
		return
	}
	writeJSON(w, http.StatusCreated, map[string]any{"ok": true, "user": user})
}

func (a *Application) handleAdminUserRole(w http.ResponseWriter, r *http.Request) {
	p, _ := principalFromContext(r.Context())
	var req struct {
		UserID string   `json:"userId"`
		Role   UserRole `json:"role"`
	}
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, "Invalid role request.")
		return
	}
	if err := a.identity.adminSetUserRole(p, req.UserID, req.Role); err != nil {
		writeError(w, adminIdentityErrorStatus(err), err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"ok": true})
}

func (a *Application) handleAdminUserStatus(w http.ResponseWriter, r *http.Request) {
	p, _ := principalFromContext(r.Context())
	var req struct {
		UserID string     `json:"userId"`
		Status UserStatus `json:"status"`
	}
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, "Invalid status request.")
		return
	}
	if err := a.identity.adminSetUserStatus(p, req.UserID, req.Status); err != nil {
		writeError(w, adminIdentityErrorStatus(err), err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"ok": true})
}

func (a *Application) handleAdminPasswordReset(w http.ResponseWriter, r *http.Request) {
	p, _ := principalFromContext(r.Context())
	var req struct {
		UserID            string `json:"userId"`
		TemporaryPassword string `json:"temporaryPassword"`
	}
	if err := json.NewDecoder(http.MaxBytesReader(w, r.Body, 16<<10)).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "Invalid password reset request.")
		return
	}
	if err := a.identity.adminResetPassword(p, req.UserID, req.TemporaryPassword); err != nil {
		writeError(w, adminIdentityErrorStatus(err), err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"ok": true, "mustSetPassword": true})
}

func (a *Application) handleAdminSessionRevoke(w http.ResponseWriter, r *http.Request) {
	p, _ := principalFromContext(r.Context())
	var req struct {
		SessionID string `json:"sessionId"`
	}
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, "Invalid session request.")
		return
	}
	if err := a.identity.adminRevokeSession(p, req.SessionID); err != nil {
		writeError(w, adminIdentityErrorStatus(err), err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"ok": true})
}
