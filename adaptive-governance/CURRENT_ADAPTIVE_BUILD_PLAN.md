# CURRENT Adaptive Build Plan

**Canonical machine state:** `governance/current-state.json`  
**Certified Stable:** `v18.9.1-stable`  
**Retained completed process authority:** #73 / `ADAPT-ROOT-CONVERGENCE-001` / `adapt-root-convergence-001`  
**Completed foundation:** #80 / `ADAPT-DATAHEALTH-BASELINE-001`  
**Completed Router adoption:** #81 / `ADAPT-PROVIDER-ROUTER-PRODUCTION-001` / PR #87 / merge `1870dd3881dbe7f6463f242e35fdc19e70d9ae15`  
**Active product work:** #82 / `ADAPT-DATAHEALTH-RUNTIME-001` / `adapt-datahealth-runtime-001`  
**Parent program:** #79 / `ADAPT-PROVIDER-PRODUCTION-001`

#80 produced the executable `provider-capability-matrix.json`, `data-health-slo.json`, `provider-fetch-paths.json` and recurrence protection. #81 consumed that migration ledger and landed Router v2 adoption for every general `MIGRATE` path without replacing provider-specific loaders, canonical cache/coalescing owners or direct authorities.

#82 required build outputs are:
1. Consolidate runtime health decisions into one canonical evaluation path using existing freshness/provider/cache/persistence/telemetry/state owners.
2. Scope health by capability/consumer/symbol/session before desk, Market Mode or app-global escalation.
3. Reuse valid warm/cached/persisted evidence before declaring degradation where freshness policy permits.
4. Make fallback/revalidation/recovery automatic when eligible, with hysteresis-protected transitions and stale warning unlatch.
5. Preserve genuine degraded truth when required evidence is stale, missing, contradictory or below quality.
6. Protect critical decision-support work before optional/background work and shed lower-value work under provider/runtime pressure.
7. Bound scanner, prep, event, research and background fan-out; prevent local overload or one optional capability failure from becoming false global degradation.
8. Record non-secret reason/scope/freshness/provider/fallback/recovery/duration telemetry and execute the #82 fault matrix.
9. Maintain Adaptive Roadmap, Build Plan, Build Process and Delivery Process alignment.

The prior sequence `#81/#82/#83/#78/#84` is now advanced through completed #81; remaining execution is #82 → #83 + #78 → #84. No parallel health subsystem and no routing-authority change is permitted.
