package main

import (
	"net/http"
	"path"
	"strings"
)

type documentationAudience int

const (
	documentationAudienceNone documentationAudience = iota
	documentationAudienceAuthenticated
	documentationAudienceDeveloper
)

func documentationAudienceForPath(rawPath string) documentationAudience {
	clean := path.Clean("/" + strings.TrimPrefix(strings.TrimSpace(rawPath), "/"))
	if clean == "/docs/developer.md" {
		return documentationAudienceDeveloper
	}
	if clean == "/docs" || strings.HasPrefix(clean, "/docs/") {
		return documentationAudienceAuthenticated
	}
	return documentationAudienceNone
}

// protectDocumentationHTTP keeps documentation audience policy server-authoritative.
// The renderer may hide privileged documentation for clarity, but direct HTTP paths
// are always checked against the same canonical session and role hierarchy used by
// protected application APIs. Local identity-disabled desktop compatibility is
// intentionally preserved; role-aware enforcement applies when IdentityService is active.
func (a *Application) protectDocumentationHTTP(next http.Handler) http.Handler {
	if next == nil {
		next = http.NotFoundHandler()
	}
	if a == nil || a.identity == nil {
		return next
	}
	authenticated := a.auth(func(w http.ResponseWriter, r *http.Request) {
		next.ServeHTTP(w, r)
	})
	developer := a.requireRole(RoleAdmin, func(w http.ResponseWriter, r *http.Request) {
		next.ServeHTTP(w, r)
	})
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch documentationAudienceForPath(r.URL.Path) {
		case documentationAudienceDeveloper:
			developer(w, r)
		case documentationAudienceAuthenticated:
			authenticated(w, r)
		default:
			next.ServeHTTP(w, r)
		}
	})
}
