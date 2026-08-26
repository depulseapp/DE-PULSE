# CURRENT Adaptive Delivery Process

**Canonical machine state:** `governance/current-state.json`  
**Certified Stable:** `v18.10.0` — immutable  
**Retained process-control authority (historical COMPLETE):** #107 / `ADAPT-PROVIDER-PROFESSIONAL-CLOSURE-001` / `adapt-provider-professional-closure-001`  
**Active v19 G1:** #148 / `ADAPT-HOSTED-TRUST-FOUNDATION-001` / `adapt-hosted-trust-foundation-001`  
**Closure ledger:** `governance/work-slices/ADAPT-HOSTED-TRUST-FOUNDATION-001/closure.json`.

v18.10.0 delivery remains authoritative and unchanged: **canonical Fast exact-head PASS**, **Qualified exact-head PASS**, canonical Release G11–G16, macOS Apple Silicon + Windows x64 native evidence, G15 provenance/SBOM and no-rebuild publication.

## Delivery / budget strategy

The v19.0.x Hosted Trust Foundation uses one branch and one long-lived Draft PR for `HOST-001..HOST-023`. Internal engineering packets remain small enough to prevent implementation gaps, but they do not automatically trigger a public release or Qualified run. Exact-head Fast validates coherent changed candidates. Default Qualified checkpoints are after `HOST-001..009`, after `HOST-010..020`, and final G10/`HOST-001..023`; add another only when the changed risk surface materially warrants it.

Do not merge unrelated packets merely to save Actions minutes. Do not split a coherent feature merely to create more release numbers. CI budget is optimized by batching coherent edits, keeping metadata with the implementation it describes, using Impact Planner, cancelling superseded runs where supported, and reserving heavy qualification/release workflows for meaningful risk/release boundaries.

## Active delivery rules

- Keep the three canonical workflow families: CI Fast, CI Qualified and Release. Do not create requirement-sized retry/certification/promotion workflows or branches.
- Each internal packet retains row-level closure/evidence ownership, so small packets do not reduce zero-gap assurance.
- Frozen v18 T1–T10 remain a conserved historical baseline; deeper historical assurance is impact-triggered without weakening changed product tests.
- Hosted Web is a real later v19 deployment/runtime target, not inferred from renderer browser qualification.
- PostgreSQL hosted activation remains blocked until tenant/account isolation, migrations, recovery/PITR, privacy lifecycle and cross-tenant adverse evidence pass.
- Cross-platform hosted authorization will eventually share account/session/device/RBAC/entitlement semantics; backend authorization stays authoritative.
- Security evidence includes negative tenant/cache/fan-out/stream/secret/mixed-client scenarios where applicable.
- Point-in-time/revision truth precedes adaptive evaluation and hosted evidence reuse.
- A `main` push may run continuity/branch hygiene, but that is not a second PR Fast.
- Public version identity advances at coherent band/release boundaries; internal packet checkpoints do not require installers/tags/releases.

Permanent product boundaries remain Smart Provider Router v2 sole general routing/admission authority, direct SEC/EDGAR filing/Form 4 authority, canonical freshness/degradation/cache/persistence/subscription/telemetry/reconciliation/lifecycle/identity/session ownership, U.S. Equities Processing with GLD/SLV/USO actionable exceptions, and No Execution.

## Exactly one next action

Finish `HOST-004..007` production reachability and coherent Fast validation on PR #149; continue the v19.0.x band without creating a separate release for that packet.
