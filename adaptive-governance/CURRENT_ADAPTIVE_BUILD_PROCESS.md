# CURRENT Adaptive Build Process

**Canonical machine state:** `governance/current-state.json`  
**Certified Stable:** `v18.9.1-stable`  
**Completed process foundation:** #73 / `ADAPT-ROOT-CONVERGENCE-001` / `adapt-root-convergence-001`

#76 / `ADAPT-TRADEINSIGHT-SETTINGS-001` completed through PR-first implementation, exact-head Fast #848, deliberate Qualified #181, and expected-head PR #77 merge `a171ce2258632bd4bd6aa737176f2d6dffb44689`. Its closure evidence is retained under `governance/work-slices/ADAPT-TRADEINSIGHT-SETTINGS-001/`.

## Process rule for #79 children

Start #80 / `ADAPT-DATAHEALTH-BASELINE-001` from live `main` only after this closure lands. Register it as a `PROCESS_RELEASE_ENGINEERING` work slice with a non-empty executable closure ledger from the first coherent commit, then synchronize all CURRENT projections.

Implementation/audit rules:
- inventory code and runtime evidence before adding architecture;
- Smart Provider Router v2 remains the sole general provider routing/admission owner;
- direct-authority paths must be explicitly classified instead of rank-swapped;
- new provider/capability/runtime external-fetch paths must fail CI if absent from the machine registry/bypass audit;
- missing credential, untested provider, connectivity failure, runtime degradation and production-readiness are distinct states;
- tests must prove authority, freshness, fallback, degradation scope, recovery and no secret leakage where applicable;
- avoid parallel provider-specific health, quota, cache, telemetry or lifecycle systems;
- exact-head Fast precedes deliberate Qualified; only qualified exact heads may merge; post-merge source-of-truth closure remains separate when required.

Authorized sequence remains #80 → #81/#82 → #83 + #78 → #84. Documentation is projection of executable owners/evidence, never a substitute for them.
