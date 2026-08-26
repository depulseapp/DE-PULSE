# CURRENT Adaptive Delivery Process

**Canonical machine state:** `governance/current-state.json`  
**Certified Stable:** `v18.10.0` — immutable  
**Post-v18 audit:** #145 / `adapt-post-v18-overlap-rebaseline-001` — PASS candidate pending exact-head qualification/merge.

v18.10.0 delivery remains authoritative and unchanged: **canonical Fast exact-head PASS**, **Qualified exact-head PASS**, canonical Release G11–G16, macOS Apple Silicon + Windows x64 native evidence, G15 provenance/SBOM and no-rebuild publication.

## Rebaselined delivery rules

- Keep the three canonical workflow families: CI Fast, CI Qualified and Release. Do not create version-specific retry/certification/promotion workflow families.
- A `main` push may run continuity and branch hygiene, but those jobs are not a second PR Fast. Documentation/run naming should make that distinction explicit.
- Minimize candidate-SHA amplification by batching coherent changes. Never reuse Fast/Qualified evidence after the candidate changes.
- Frozen v18 T1–T10 become a conserved historical baseline. Do not let every future major version append another permanently unconditional chain of historical gate scripts; deeper historical assurance should be impact-triggered while baseline conservation remains fail closed.
- `HOST-001..072` are traceability requirements. Delivery groups them into coherent implementation/release bands; row labels do not require individual release runs.
- Hosted Web is a real v19 deployment/runtime target, not inferred from renderer browser qualification. Hosted release evidence must eventually include deploy identity, migration compatibility, tenant isolation, secret/service trust, rollback/recovery and production SLO evidence.
- PostgreSQL hosted activation remains blocked until tenant/account schema and authorization isolation, recovery/PITR, migration strategy, privacy lifecycle and adverse cross-tenant evidence pass.
- Cross-platform parity means Mac + Windows + Web share the same account/session/device/RBAC/entitlement/domain semantics. Backend authorization is always authoritative.
- Security delivery evidence includes negative tenant/cache/fan-out/stream-revocation/secret/mixed-client scenarios, not only happy-path authentication.
- Point-in-time/revision truth and decision-outcome calibration are part of future trader-quality delivery evidence where affected.
- Repository `main` protection should be enabled when plan capabilities support canonical required statuses; internal release provenance remains mandatory regardless.

Permanent product boundaries remain Smart Provider Router v2 sole general routing/admission authority, direct SEC/EDGAR filing/Form 4 authority, canonical freshness/degradation/cache/persistence/subscription/telemetry/reconciliation/lifecycle/identity/session ownership, U.S. Equities Processing with GLD/SLV/USO actionable exceptions, and No Execution.

## Exactly one next action

Qualify and merge #145; then update canonical machine state to PASS and only afterward reserve the first coherent v19 G1 band.
