package main

import (
	"encoding/json"
	"net/http"
	"strings"
)

const (
	maxHostedDeviceLabelLength       = 128
	maxHostedDeviceFingerprintLength = 512
)

func (a *Application) registerHostedIdentityRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/api/auth/device/register", a.auth(a.requireRecentAuthentication(postOnly(a.handleRegisterHostedDevice))))
	mux.HandleFunc("/api/auth/device/status", a.auth(a.requireRecentAuthentication(postOnly(a.handleHostedDeviceStatus))))
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
