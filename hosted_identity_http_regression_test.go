package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func hostedIdentityHTTPMux(app *Application) *http.ServeMux {
	mux := http.NewServeMux()
	app.registerHealthRoutes(mux)
	return mux
}

func TestHOST006HostedDeviceHTTPRegistrationBindsCurrentSession(t *testing.T) {
	_, s := newIdentityTestService(t)
	base := time.Unix(2_410_000_000, 0)
	token, p, _ := v184CredentialedOwner(t, s, base)
	app := &Application{identity: s}
	mux := hostedIdentityHTTPMux(app)
	csrf := "hosted-device-register-csrf"
	fingerprint := "sha256:http-registered-device"

	res := httptest.NewRecorder()
	mux.ServeHTTP(res, v184AuthenticatedPOST(
		"/api/auth/device/register",
		token,
		csrf,
		`{"label":"primary browser","fingerprintHash":"`+fingerprint+`"}`,
	))
	if res.Code != http.StatusOK {
		t.Fatalf("device registration failed: code=%d body=%s", res.Code, res.Body.String())
	}
	if strings.Contains(res.Body.String(), fingerprint) {
		t.Fatalf("device fingerprint leaked in response: %s", res.Body.String())
	}
	var body struct {
		OK bool `json:"ok"`
		Device struct {
			ID     string       `json:"id"`
			Label  string       `json:"label"`
			Status DeviceStatus `json:"status"`
		} `json:"device"`
		SessionBound bool `json:"sessionBound"`
	}
	if err := json.Unmarshal(res.Body.Bytes(), &body); err != nil {
		t.Fatal(err)
	}
	if !body.OK || !body.SessionBound || body.Device.ID == "" || body.Device.Status != DeviceActive {
		t.Fatalf("unexpected registration contract: %+v", body)
	}
	resolved, err := s.resolve(token, false)
	if err != nil {
		t.Fatalf("bound session stopped resolving: %v", err)
	}
	if resolved.DeviceID != body.Device.ID {
		t.Fatalf("session was not bound to registered device: principal=%+v device=%+v", resolved, body.Device)
	}
	decision := s.authorizeHostedIdentity(p, HostedIdentityRequirement{
		TenantID:                localTenantID,
		Capability:              hostedCapabilityDeviceManage,
		RequireRegisteredDevice: true,
	})
	if !decision.Allowed || decision.DeviceID != body.Device.ID {
		t.Fatalf("registered HTTP device did not become canonical trust evidence: %+v", decision)
	}
}

func TestHOST006HostedDeviceHTTPStatusRequiresRegisteredCurrentDevice(t *testing.T) {
	_, s := newIdentityTestService(t)
	base := time.Unix(2_420_000_000, 0)
	token, _, _ := v184CredentialedOwner(t, s, base)
	app := &Application{identity: s}
	mux := hostedIdentityHTTPMux(app)
	csrf := "hosted-device-status-csrf"

	res := httptest.NewRecorder()
	mux.ServeHTTP(res, v184AuthenticatedPOST(
		"/api/auth/device/status",
		token,
		csrf,
		`{"deviceId":"dev-untrusted","status":"REVOKED"}`,
	))
	if res.Code != http.StatusForbidden {
		t.Fatalf("unregistered session crossed device-management trust boundary: code=%d body=%s", res.Code, res.Body.String())
	}
	if _, err := s.resolve(token, false); err != nil {
		t.Fatalf("denied status request mutated current session: %v", err)
	}
}

func TestHOST006HostedDeviceHTTPStatusRevokesCurrentBoundSession(t *testing.T) {
	_, s := newIdentityTestService(t)
	base := time.Unix(2_430_000_000, 0)
	token, _, _ := v184CredentialedOwner(t, s, base)
	app := &Application{identity: s}
	mux := hostedIdentityHTTPMux(app)
	csrf := "hosted-device-lost-csrf"

	registered := httptest.NewRecorder()
	mux.ServeHTTP(registered, v184AuthenticatedPOST(
		"/api/auth/device/register",
		token,
		csrf,
		`{"label":"travel laptop","fingerprintHash":"sha256:http-lost-device"}`,
	))
	if registered.Code != http.StatusOK {
		t.Fatalf("device registration failed: code=%d body=%s", registered.Code, registered.Body.String())
	}
	var registration struct {
		Device struct {
			ID string `json:"id"`
		} `json:"device"`
	}
	if err := json.Unmarshal(registered.Body.Bytes(), &registration); err != nil || registration.Device.ID == "" {
		t.Fatalf("missing registered device id: body=%s err=%v", registered.Body.String(), err)
	}

	lost := httptest.NewRecorder()
	mux.ServeHTTP(lost, v184AuthenticatedPOST(
		"/api/auth/device/status",
		token,
		csrf,
		`{"deviceId":"`+registration.Device.ID+`","status":"LOST"}`,
	))
	if lost.Code != http.StatusOK {
		t.Fatalf("trusted device status update failed: code=%d body=%s", lost.Code, lost.Body.String())
	}
	var body struct {
		OK                    bool         `json:"ok"`
		Status                DeviceStatus `json:"status"`
		CurrentSessionRevoked bool         `json:"currentSessionRevoked"`
	}
	if err := json.Unmarshal(lost.Body.Bytes(), &body); err != nil {
		t.Fatal(err)
	}
	if !body.OK || body.Status != DeviceLost || !body.CurrentSessionRevoked {
		t.Fatalf("unexpected lost-device contract: %+v", body)
	}
	if _, err := s.resolve(token, false); err == nil {
		t.Fatal("LOST current device did not revoke its bound session")
	}
	clearedSessionCookie := false
	for _, cookie := range lost.Result().Cookies() {
		if cookie.Name == sessionCookieName && cookie.MaxAge < 0 {
			clearedSessionCookie = true
			break
		}
	}
	if !clearedSessionCookie {
		t.Fatal("revoked current device did not clear the browser session cookie")
	}
}

func TestHOST006HostedDeviceHTTPDoesNotReactivateRevokedTrust(t *testing.T) {
	_, s := newIdentityTestService(t)
	base := time.Unix(2_440_000_000, 0)
	token, _, _ := v184CredentialedOwner(t, s, base)
	app := &Application{identity: s}
	mux := hostedIdentityHTTPMux(app)
	csrf := "hosted-device-active-csrf"

	registered := httptest.NewRecorder()
	mux.ServeHTTP(registered, v184AuthenticatedPOST(
		"/api/auth/device/register",
		token,
		csrf,
		`{"label":"browser","fingerprintHash":"sha256:http-active-reject"}`,
	))
	if registered.Code != http.StatusOK {
		t.Fatalf("device registration failed: code=%d body=%s", registered.Code, registered.Body.String())
	}
	var registration struct {
		Device struct {
			ID string `json:"id"`
		} `json:"device"`
	}
	if err := json.Unmarshal(registered.Body.Bytes(), &registration); err != nil {
		t.Fatal(err)
	}

	active := httptest.NewRecorder()
	mux.ServeHTTP(active, v184AuthenticatedPOST(
		"/api/auth/device/status",
		token,
		csrf,
		`{"deviceId":"`+registration.Device.ID+`","status":"ACTIVE"}`,
	))
	if active.Code != http.StatusBadRequest {
		t.Fatalf("HTTP surface allowed direct ACTIVE mutation: code=%d body=%s", active.Code, active.Body.String())
	}
	if _, err := s.resolve(token, false); err != nil {
		t.Fatalf("rejected ACTIVE request mutated current session: %v", err)
	}
}
