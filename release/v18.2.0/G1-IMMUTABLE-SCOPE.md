# G1 Immutable Scope — v18.2.0

1. Extend canonical `IdentityService`; no second identity/session store.
2. Role-aware user creation for roles below actor authority.
3. Role/status lifecycle with active OWNER/SUPER_OWNER safety.
4. Temporary-password reset with forced change and session revocation.
5. Redacted admin user/session visibility; never expose credential/token hashes.
6. Presence states ACTIVE / IDLE / OFFLINE derived from canonical sessions.
7. Session revoke plus SSE termination when revoked/expired.
8. Compact OWNER/ADMIN-class Settings administration UI.
9. USER/DEMO do not see admin controls.
10. Preserve v18.1 per-user market-state isolation and shared fetch-once/calculate-once intelligence core.
11. No v18.3 hosted PostgreSQL/web architecture.
12. No execution scope and no deterministic Day/Swing/Long formula changes.
