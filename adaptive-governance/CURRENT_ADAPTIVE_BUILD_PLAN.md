# CURRENT Adaptive Build Plan

**Canonical machine state:** `governance/current-state.json`  
**Certified Stable:** `v18.9.1-stable`  
**Completed provider-intelligence program:** #65 / #107  
**Planning rebaseline:** #110 / `ADAPT-V19-ZERO-MISS-PLAN-001`  
**Future hosted umbrella:** #66 — planned but no product work slice is reserved.

## Zero-miss v19 build contract

The canonical v19 plan is `governance/programs/ADAPT-HOSTED-SYNC-001/V19_ZERO_MISS_PLAN.md`; its machine conservation ledger is `governance/programs/ADAPT-HOSTED-SYNC-001/requirement-conservation.json`.

Every future patch:
1. owns one primary implementation responsibility;
2. starts with live-source overlap classification: `INHERITED`, `EXTEND_EXISTING_OWNER`, `REPLACE_CONSOLIDATE`, `NEW_RESIDUAL`, or `EXTERNAL_BLOCKED`;
3. binds all applicable #66 requirement IDs before G1;
4. records dependencies and exact completion/evidence criteria;
5. cannot close with an applicable unassigned/unevidenced row;
6. cannot advance to the next dependency band until the band zero-gap closure passes;
7. implements every shared user-facing capability in Mac + Windows + Web lockstep for all G1 `REQUIRED` clients.

The inherited Data Health build inputs remain `governance/data-health/provider-capability-matrix.json`, `governance/data-health/data-health-slo.json`, and `governance/data-health/provider-fetch-paths.json`. Their recurrence and ownership remain projected consistently through the Adaptive Roadmap, Build Plan, Build Process and Delivery Process. v19 extends existing owners; it does not create parallel Router, health, freshness, cache, persistence, telemetry, reconciliation, subscription, identity or lifecycle systems.

## Dependency bands

- `v19.0.0`–`v19.0.18`: rights, tenant/control, privacy, environment/trust, PostgreSQL, secrets, supply chain, provider quality and point-in-time primitives; `v19.0.18` is no-feature foundation closure.
- `v19.1.0`–`v19.1.15`: authenticated gateway, serving policy, lawful reuse/live fan-out and typed sync primitives; `v19.1.15` is no-feature data-plane/sync closure.
- `v19.2.0`–`v19.2.16`: Mac + Windows + Web capability enablement plus multi-user/DR/load assurance; `v19.2.16` is #66 closure.
- `v19.3.0`–`v19.3.7`: institutional/13F, two-sided evidence and AODR lineage; `v19.3.7` is no-feature evidence closure.
- `v19.4.0`–`v19.4.6`: operational SLO/runbooks/economics/adaptive readiness and zero-gap sweep.
- `v19.5.0`: no-feature Major Closure.

v20 remains provisional and cannot begin before `v19.4.5` readiness plus `v19.5.0` Major Closure.
