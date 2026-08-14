# G3 Design / Dependency Readiness — v18.2.0

- Existing Argon2id, opaque HttpOnly session, SameSite Strict, CSRF and RBAC middleware are reused.
- New admin mutations remain authenticated, CSRF-protected and role-gated.
- Admin UI is modular (`admin-v18.2.js` / `admin-v18.2.css`) and loaded after the canonical renderer.
- Responsive rules cover desktop, tablet and narrow browser widths.
- Native packaging targets remain macOS Apple Silicon and Windows x64.
- Existing certification, responsive, HTTP, deterministic, security and native lifecycle harnesses remain blocking regressions.
