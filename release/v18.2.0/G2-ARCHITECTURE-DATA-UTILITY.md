# G2 Architecture / Data Utility — v18.2.0

- Identity ownership remains `IdentityService` + `IdentityPersistentState`.
- Presence is derived from persisted session `LastSeenAt`, idle/absolute expiry and revocation; no separate heartbeat database.
- Long-lived SSE reuses the same session truth and closes after administrative revocation/expiry.
- Administrative DTOs expose only operational fields required for user/session management. Password hashes, token hashes, opaque session tokens and secrets are excluded.
- Market providers, evidence, Router/Rapid Move, canonical quotes and deterministic scoring remain shared global intelligence.
- v18.1 `UserWorkspace` remains the owner of personal tracked symbols/watchlists/UI state.
- REUSE/CONSOLIDATE before ADD is mandatory; no parallel account/session/presence implementation.
