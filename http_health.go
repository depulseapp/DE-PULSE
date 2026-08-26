package main

import (
	"context"
	"net/http"
	"time"
)

func (a *Application) registerHealthRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/api/health", a.handleHealth)
	mux.HandleFunc("/api/ready", a.handleReady)
	a.registerHostedIdentityRoutes(mux)
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
