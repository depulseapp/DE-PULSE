# CURRENT Adaptive Build Process

**Canonical machine state:** `governance/current-state.json`  
**Certified Stable:** `v18.9.1-stable`  
**Retained completed process authority:** #73 / `ADAPT-ROOT-CONVERGENCE-001`  
**Completed provider/data-health program:** #79 with final #84 / PR #91 / merge `733d90ca125a4fe5abd38a2ea40de0623703dfd4`  
**Completed canonical identity:** #92 / PR #93 / merge `57d530e58bfb0b38cc108980cd5cd4a041014db8`  
**Active product work:** #95 / `ADAPT-PROVIDER-ONBOARDING-001` / `adapt-provider-onboarding-001`  
**Separate sibling residual:** #94 / `ADAPT-PROVIDER-TELEMETRY-001`  
**Parent program:** #65 / `ADAPT-PROVIDER-INTELLIGENCE-010`

For #95, implementation remains executable-first, overlap-first, reuse-first and canonical-owner-first:
- fetch live `main` and the live #95 branch before each mutation because concurrent sessions may advance either;
- inspect issue #65/#95 plus commits and existing source before adding a provider-onboarding owner; duplicated router/health/lifecycle/persistence architecture fails review;
- treat `ProviderRegistration` as admission description only. Smart Provider Router v2 owns route selection; provider lifecycle/readiness owns promotion; Data Health/freshness/degradation own current-truth health; rights governance remains separate;
- require a real adapter/normalizer and complete schema/timestamp/freshness/failure/rights/evidence/approval/invalidation contract before a capability can be production-routable;
- use production-route validation that must **fail closed** so a lifecycle string cannot bypass missing contract evidence;
- keep provider configuration fingerprints process-local and one-way; secrets or derived credential material must never be emitted to logs, snapshots, telemetry, persistence, tests, handoff or governance;
- revalidate entitlement only at safe canonical decision/action boundaries. Never mutate capability state while `Engine.Snapshot()` holds engine read locks;
- configuration changes may reopen only stale `NOT_ENTITLED`/`NOT_CONFIGURED` suppression for the affected provider; manual capability recheck may explicitly request a fresh entitlement probe for same-key plan changes;
- do not clear `SUPPORTED`, `NOT_SUPPORTED`, genuine outage/rate-limit/capacity evidence simply to make a provider eligible;
- preserve capability-scoped failure isolation and bounded fallback; provider plan failure on one dataset must not suppress unrelated datasets;
- preserve provider-specific diagnostic semantics while projecting them from the canonical registration; avoid incidental UI/order changes;
- regression-test synthetic registration adoption, production fail-closed behavior, free→paid and paid→restricted transitions, same-key recheck, route parity, capability isolation and current-provider behavior before canonical CI;
- keep #94 semantic usefulness telemetry separate and advisory/non-routing until its own governed validation/promotion evidence exists;
- update `governance/current-state.json`, the #95 work-slice/closure ledger, all four CURRENT Adaptive projections and `handoff/CURRENT.md` with the same active ownership and remaining evidence gaps;
- use canonical Fast exact-head PASS followed by deliberate impact-selected Qualified exact-head PASS; any source change after qualification creates a new candidate and invalidates earlier evidence;
- do not add workflow families, G17+ gates, parallel caches/databases/services, or weaken G0–G16/source-health/architecture gates.

## Retained Adaptive Data Health process contract
The completed sequence **#81/#82/#83/#78/#84** remains inherited executable history from the #80 baseline. Its canonical freshness, provider-capability classification, scoped degradation/recovery, Router v2, lifecycle/readiness and direct-authority boundaries remain active constraints on #95. New provider registration/adaptation must fail closed into those existing owners rather than replacing them.

## Acceptance discipline
Source commits are not proof of closure. The #95 ledger remains OPEN until executable evidence proves current-provider regression safety, entitlement adaptation, diagnostic compatibility and exact-head CI. Documentation may describe evidence but cannot substitute for it.

## Retained Adaptive Data Health invariants
The completed #79/#84 architecture remains authoritative under #95: truthful PARTIAL COVERAGE/DATA DEGRADED states, canonical freshness, Router v2, lifecycle/readiness, fault recovery, provider workload controls, cache/persistence reuse and provenance remain intact. Direct SEC/EDGAR Form 4 authority, U.S. equities, GLD/SLV/USO and No Execution are permanent boundaries.
