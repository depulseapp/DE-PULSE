package main

import (
	"net/http"
)

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
