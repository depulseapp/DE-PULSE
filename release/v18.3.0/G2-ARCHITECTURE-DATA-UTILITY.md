# G2 Architecture / Data Utility — v18.3.0

- `PersistenceBackend` remains the single storage contract for symbol registry, quotes, structured intelligence, identity state, workspaces and store statistics.
- `PersistenceManager` remains the canonical batching/coalescing/material-change owner. PostgreSQL is a backend implementation, not a second queue, cache or business layer.
- Local desktop backend selection continues through the existing SQLite implementation (or the existing non-CGO file fallback where applicable).
- Hosted PostgreSQL is explicit via runtime configuration and a `postgres` build-tagged driver. If PostgreSQL is requested but unavailable/unconfigured, persistence readiness fails closed rather than falling back locally.
- PostgreSQL schema versions 1–4 mirror SQLite ownership: symbols/quotes, history/intelligence, identity state, and per-user workspaces.
- `UserWorkspace` remains personal-market-state ownership; `processingStateLocked()` remains the one canonical deduplicated union used for market/provider/scanner/intelligence work.
- `IdentityService` remains the only identity/session lifecycle owner. PostgreSQL persists canonical identity state through the same repository contract; it does not create another auth model.
- Warm-start canonical quotes retain values/provenance but are relabelled `persisted`; runtime freshness/live truth must be re-established by normal evidence owners.
- Hosted runtime selection is an outer application/runtime concern. Handlers, intelligence modules and deterministic scoring do not branch on database type.
- `/api/health` is process liveness. `/api/ready` is dependency readiness and must reflect canonical persistence plus IdentityService availability.
- PostgreSQL pool and operation diagnostics are exposed through existing runtime persistence diagnostics; DSNs, credentials, password hashes, token hashes and provider secrets are never emitted.
- Database work remains bounded and coalesced; provider/data acquisition is not triggered by persistence reads/writes.
