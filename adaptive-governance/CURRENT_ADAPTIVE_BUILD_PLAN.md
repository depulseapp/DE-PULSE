# CURRENT Adaptive Build Plan

**Canonical machine state:** `governance/current-state.json`  
**Certified Stable:** `v18.10.0` — immutable  
**Retained process-control authority (historical COMPLETE):** #107 / `ADAPT-PROVIDER-PROFESSIONAL-CLOSURE-001` / `adapt-provider-professional-closure-001`  
**Post-v18 audit:** #145 / `adapt-post-v18-overlap-rebaseline-001` — PASS candidate pending exact-head qualification/merge  
**Future hosted umbrella:** #66 — `PLANNED_UNSTARTED` until #145 is merged.

The v18.10.0 build contract remains closed: T1–T10, exact-head Fast/Qualified, macOS Apple Silicon + Windows x64 lifecycle, G15 provenance, SBOM, no-rebuild publication and G16 all passed. Future work must conserve those responsibilities without reproducing a new historical closure stack on every version.

## Conserved Data Health build inputs

The inherited Data Health build inputs remain `governance/data-health/provider-capability-matrix.json`, `governance/data-health/data-health-slo.json`, and `governance/data-health/provider-fetch-paths.json`. Their ownership and recurrence continue to project consistently through the **Adaptive Roadmap**, **Adaptive Build Plan**, **Build Process**, and **Delivery Process**. Future changes reuse the existing Router, Data Health, canonical freshness, cache, persistence, telemetry, reconciliation, subscription, identity and lifecycle owners.

## Rebaselined build rules

### Requirement/evidence binding at implementation time
Every changed requirement must enter the slice with: canonical owner, consumers, reuse/consolidation disposition, positive evidence owner, adverse evidence owner, persistence/security/UI applicability, point-in-time/freshness semantics where data-bearing, and closure acceptance. Classify findings as `PRODUCT_BEHAVIOR_GAP`, `TEST_OR_EVIDENCE_GAP`, `OWNERSHIP_BINDING_GAP`, or `NOT_APPLICABLE` before changing product code.

### Coherent HOST implementation bands
`HOST-001..HOST-072` remain row-level traceability requirements, not mandatory one-row/one-version release events. Build coherent dependency bands and batch related mutations before opening/advancing a PR where practical. Every band still closes with zero unexplained applicable row.

Recommended high-level sequencing:
- **Hosted trust foundation:** provider rights + tenant/account identity + session/device + product entitlement/quota + privacy/IaC/service trust + PostgreSQL tenant/recovery + managed secrets/supply chain + point-in-time primitives.
- **Hosted serving/sync foundation:** gateway/serving policy + entitlement-safe cache/coalescing + shared live fan-out + API lifecycle + transactional outbox/mutation/pull/bootstrap/conflict/retention + tenant observability/backpressure.
- **Cross-platform hosted parity:** Mac/Windows/Web account/device/RBAC/entitlement + portable domain sync + shared Discovery/Market State + fairness/adversarial security + mixed-client/recovery/load evidence.
- **Evidence/trader-quality:** institutional point-in-time evidence + two-sided research + AODR lineage/outcomes + reliability/cost/rights + adaptive research-readiness audit.
- **Major closure:** zero-gap G0–G16 only; no feature scope.

### PostgreSQL boundary
Reuse `persistence_backend_postgres.go`; do not create a second hosted persistence layer. Connection success is not hosted readiness. Tenant/account schema ownership, isolation/RLS disposition, migrations, recovery/PITR, capacity, privacy lifecycle, authorization and adverse cross-tenant evidence are required before shared activation.

### Trading-quality build requirements
Point-in-time/no-lookahead truth is prerequisite to historical evaluation. Extend canonical decision/outcome owners for MFE/MAE, horizon return, false-positive/miss rate, confidence calibration and usefulness by regime/freshness/contradiction. Provider usefulness remains advisory until lifecycle-approved. Resolve exchange halt/LULD/volatility-pause/resume semantics under the canonical tradeability/event owner before professional hosted closure.

### CI/build efficiency
Do not reduce evidence for changed candidates. Reduce candidate churn and historical gate accumulation instead. Frozen v18 T1–T10 become a conservation baseline; deeper historical assurance runs only when affected owners/contracts change. Keep exact-head Fast/Qualified and canonical Release.

## Exactly one next action

Qualify and merge #145. After the audit is PASS on `main`, reserve one coherent first v19 G1 band—not a chain of requirement-sized micro-releases.
