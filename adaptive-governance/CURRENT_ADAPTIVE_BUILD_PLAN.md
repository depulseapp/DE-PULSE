# CURRENT Adaptive Build Plan

**Canonical machine state:** `governance/current-state.json`  
**Certified Stable:** `v18.10.0` — immutable  
**Retained process-control authority (historical COMPLETE):** #107 / `ADAPT-PROVIDER-PROFESSIONAL-CLOSURE-001` / `adapt-provider-professional-closure-001`  
**Active v19 G1:** #148 / `ADAPT-HOSTED-TRUST-FOUNDATION-001` / `adapt-hosted-trust-foundation-001`  
**Work-slice:** `governance/work-slices/ADAPT-HOSTED-TRUST-FOUNDATION-001/work-slice.json`  
**G1 scope:** `governance/work-slices/ADAPT-HOSTED-TRUST-FOUNDATION-001/g1-scope.json`  
**Closure ledger:** `governance/work-slices/ADAPT-HOSTED-TRUST-FOUNDATION-001/closure.json`  
**Parent:** #66 / `ADAPT-HOSTED-SYNC-001`.

The v18.10.0 build contract remains closed: T1–T10, exact-head Fast/Qualified, macOS Apple Silicon + Windows x64 lifecycle, G15 provenance, SBOM, no-rebuild publication and G16 all passed. #148 owns `HOST-001..HOST-023` as one coherent dependency-ordered architecture band rather than 23 micro-releases.

## Conserved Data Health build inputs

The inherited Data Health build inputs remain `governance/data-health/provider-capability-matrix.json`, `governance/data-health/data-health-slo.json`, and `governance/data-health/provider-fetch-paths.json`. Their ownership and recurrence continue to project consistently through the **Adaptive Roadmap**, **Adaptive Build Plan**, **Build Process**, and **Delivery Process**. Future changes reuse the existing Router, Data Health, canonical freshness, cache, persistence, telemetry, reconciliation, subscription, identity and lifecycle owners.

## Active build order

- HOST-001..003 provider rights.
- HOST-004..007 tenant/account identity/device/session.
- HOST-008..009 product entitlement/quota.
- HOST-010..014 privacy/environment/service trust.
- HOST-015..016 PostgreSQL tenancy/recovery via existing `persistence_backend_postgres.go`.
- HOST-017..020 managed secrets/supply chain.
- HOST-021..022 provider scorecards/point-in-time truth.
- HOST-023 zero-gap closure.

Every changed requirement carries canonical owner, consumers, reuse/consolidation disposition, positive evidence, adverse evidence, persistence/security/UI applicability, point-in-time/freshness semantics where data-bearing, and closure acceptance. Classify findings as `PRODUCT_BEHAVIOR_GAP`, `TEST_OR_EVIDENCE_GAP`, `OWNERSHIP_BINDING_GAP`, or `NOT_APPLICABLE` before changing product code.

Do not reduce evidence for changed candidates. Reduce candidate churn and historical gate accumulation instead. Frozen v18 T1–T10 remain a conservation baseline; exact-head Fast/Qualified and canonical Release remain mandatory where applicable.

## Exactly one next action

Pass exact-head G1 Fast for #148, then implement `HOST-001..HOST-003_PROVIDER_RIGHTS` through existing provider-rights owners on the same Draft PR.
