# G1 Immutable Scope — v18.3.0

1. Add PostgreSQL repository parity beneath the existing canonical `PersistenceBackend`; do not create parallel business/data owners.
2. Preserve SQLite/local persistence as the default macOS/Windows desktop backend and preserve existing desktop runtime behavior unless explicitly entering hosted mode.
3. PostgreSQL selection is explicit and fail-closed; a requested hosted database must never silently fall back to local per-machine state.
4. Provide migration/schema parity for Global Symbol Registry, canonical quotes/history, evidence, decision lineage, outcomes, derived features, identity state, and per-user workspaces.
5. Use bounded transactions/concurrency and bounded connection pooling suitable for shared hosted state.
6. Preserve one shared canonical market-processing universe across all users; PostgreSQL must not create per-user provider, scanner, Router, Rapid Move, Opportunity Radar, research, or deterministic-scoring engines.
7. Preserve authenticated per-user workspace isolation and the canonical `IdentityService` ownership model.
8. Preserve truthful warm-start semantics: persisted quotes/evidence are restart context, never timeless LIVE truth.
9. Add explicit hosted browser/server runtime mode with separate liveness and persistence-backed readiness; default desktop native-window behavior remains intact.
10. Add DB/query/pool observability and bounded persistence-pressure diagnostics without exposing secrets.
11. Implement and prove backup/restore plus SQLite-to-PostgreSQL migration/export before G10/G11 freeze.
12. Run real PostgreSQL parity, migration concurrency, load/contention, restart and failure/recovery qualification before release freeze.
13. Harden dependency-compatible ADR-GDI behavior: indexed warm state, bounded persistence pressure, scoped degradation/readiness, and no database retry/fan-out storms.
14. Preserve protected deterministic Day/Swing/Long formulas and the permanent No Execution Boundary.
15. Do not expand into v18.4 commercial/data-rights/security scope except unavoidable safety corrections required by v18.3.
16. Required final native targets remain macOS Apple Silicon and Windows x64; hosted readiness is additive, not a replacement for required desktop certification.
