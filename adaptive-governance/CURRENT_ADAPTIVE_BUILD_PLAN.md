# CURRENT Adaptive Build Plan

**Canonical machine state:** `governance/current-state.json`  
**Certified Stable:** `v18.10.0` — immutable  
**Retained process-control authority (historical COMPLETE):** #107 / `ADAPT-PROVIDER-PROFESSIONAL-CLOSURE-001` / `adapt-provider-professional-closure-001`  
**Active v19 G1:** #148 / `ADAPT-HOSTED-TRUST-FOUNDATION-001` / `adapt-hosted-trust-foundation-001`  
**Work-slice:** `governance/work-slices/ADAPT-HOSTED-TRUST-FOUNDATION-001/work-slice.json`  
**G1 scope:** `governance/work-slices/ADAPT-HOSTED-TRUST-FOUNDATION-001/g1-scope.json`  
**Closure ledger:** `governance/work-slices/ADAPT-HOSTED-TRUST-FOUNDATION-001/closure.json`  
**Parent:** #66 / `ADAPT-HOSTED-SYNC-001`.

The v18.10.0 build contract remains closed: T1–T10, exact-head Fast/Qualified, macOS Apple Silicon + Windows x64 lifecycle, G15 provenance, SBOM, no-rebuild publication and G16 all passed. #148 owns `HOST-001..HOST-023` as one coherent **v19.0.x Hosted Trust Foundation band**, not 23 micro-releases.

## Packet sizing and version rule

Engineering packets stay small enough to make omissions visible. Release/qualification boundaries stay large enough to avoid Actions waste. Do not trade implementation completeness for fewer CI runs, and do not turn every small packet into a new public version.

Internal packets:
- P1 `HOST-001..003` provider rights.
- P2 `HOST-004..007` tenant/account identity/device/session/reauth.
- P3 `HOST-008..009` product entitlement/quota.
- P4 `HOST-010..014` privacy/environment/service trust.
- P5 `HOST-015..016` PostgreSQL tenancy/recovery via existing `persistence_backend_postgres.go`.
- P6 `HOST-017..020` managed secrets/KMS + supply-chain/deploy provenance.
- P7 `HOST-021..022` provider scorecards/point-in-time truth.
- P8 `HOST-023` zero-gap closure.

Per-packet acceptance is strict: canonical owner, consumers, reuse/consolidation disposition, positive evidence, adverse evidence, persistence/security/UI applicability, point-in-time/freshness semantics where data-bearing, and closure acceptance must all be bound before the packet can advance. Findings are classified as `PRODUCT_BEHAVIOR_GAP`, `TEST_OR_EVIDENCE_GAP`, `OWNERSHIP_BINDING_GAP`, or `NOT_APPLICABLE` before changing product code.

Default Qualified checkpoints are: Q1 after P1–P3 (`HOST-001..009`) because provider rights + identity + entitlement form one authorization boundary; Q2 after P4–P6 (`HOST-010..020`) because privacy + PostgreSQL + secrets form the hosted trust/persistence boundary; Q3 final after P7–P8/G10 (`HOST-001..023`). Extra Qualified runs are permitted only for materially new risk surfaces, not because a packet number changed. Coherent implementation candidates still receive exact-head Fast; batch pure closure/evidence metadata with its implementation commit when safe.

## Conserved Data Health build inputs

The inherited Data Health build inputs remain `governance/data-health/provider-capability-matrix.json`, `governance/data-health/data-health-slo.json`, and `governance/data-health/provider-fetch-paths.json`. Their ownership and recurrence continue to project consistently through the **Adaptive Roadmap**, **Adaptive Build Plan**, **Build Process**, and **Delivery Process**. Future changes reuse the existing Router, Data Health, canonical freshness, cache, persistence, telemetry, reconciliation, subscription, identity and lifecycle owners.

Do not reduce evidence for changed candidates. Reduce candidate churn and historical gate accumulation instead. Frozen v18 T1–T10 remain a conservation baseline; exact-head Fast/Qualified and canonical Release remain mandatory where applicable.

## Exactly one next action

Finish P2 `HOST-004..007` production reachability on PR #149, obtain coherent exact-head Fast, then continue P3 on the same branch/PR. Do not create a separate public version or release for P2.
