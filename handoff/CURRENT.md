# DE.PULSE Current Handoff

**SUPERSEDES ALL PRIOR CHAT HANDOFFS**

**Canonical machine state:** `governance/current-state.json`  
**Certified Stable:** `v18.9.1-stable`  
**Certified Stable candidate:** `e55d8d25b15cec2ffb0f5411bc358bc40b359cf9`  
**Certified Stable source fingerprint:** `0062f46dea5690d0b3fcd8a9ed3b1f71ebe1522c7dee2cb218e9d36b9e0076ff`  
**Certified Stable build ID:** `v18.9.1-stable-20260821`  
**Immutable predecessor resume checkpoint release:** `v18.9.0` / `v18.9.0-stable`  
**Retained completed process authority:** #73 / `ADAPT-ROOT-CONVERGENCE-001` / `adapt-root-convergence-001`  
**Canonical retained process closure ledger:** `governance/work-slices/ADAPT-ROOT-CONVERGENCE-001/closure.json`  
**Completed Data Health baseline:** #80 / `ADAPT-DATAHEALTH-BASELINE-001` / candidate `1d8638acb06c7ce90719fc3b959d37f188eb8b40` / Fast #859 / Qualified #182 / PR #86 / merge `c75a5f1467920f57fa23c3dbc400e51edc5275c8`  
**Active product slice:** #81 / `ADAPT-PROVIDER-ROUTER-PRODUCTION-001` / `adapt-provider-router-production-001`  
**Parent program:** #79 / `ADAPT-PROVIDER-PRODUCTION-001`

## Current authority
Issue #79 is the cross-session authority for the all-provider Adaptive Data Health program. #80 is complete and its machine-readable provider/capability/SLO/fetch-path baseline remains authoritative. Current dependency order is **#81/#82 → #83 + #78 → #84**.

## #81 objective
Make Smart Provider Router v2 the real executable production authority for every routable provider capability identified as `MIGRATE` by #80. Preserve provider-specific HTTP/normalization loaders, canonical cache/single-flight/coalescing, direct authority boundaries, bounded fallback, provider workload/backpressure controls, capability-scoped circuit truth and non-secret decision reasons. Do not create a second router.

Known #80 migration debt begins with:
- Alpaca read-only U.S. market calendar;
- Alpaca canonical U.S. asset universe, with bounded live/paper same-provider endpoint fallback inside the Alpaca attempt;
- Alpaca corporate-actions announcements;
- Twelve Data direct global indices/futures acquisition.

The exact runtime owners include `provider_router.go`, `smart_router_v2.go`, `routed_refresh.go`, `symbol_universe.go`, `preparation_types_liquidity.go`, `market_activity_corporate.go`, and `global_market_providers.go`. Official public/authority fetches are not to be rank-swapped merely to eliminate a direct call.

## Permanent boundaries
- U.S. equities processing only.
- No Execution/order routing.
- Smart Provider Router v2 is the sole general routing/admission authority.
- Direct SEC/EDGAR remains authoritative for Form 4.
- GLD, SLV and USO remain actionable live-priority exceptions.
- TradeInsight remains shadow-first where governed.
- Reuse canonical freshness/cache/persistence/telemetry/state/validation owners.
- GitHub executable evidence outranks chat memory.

## Exactly one next action
On `adapt-provider-router-production-001`, implement the #81 Router v2 migration for the #80 `MIGRATE` paths using existing owners, extend existing Smart Router/Data Health regression protection, update the fetch-path ledger only after the bypass is actually removed, then obtain exact-head Fast and Qualified before merge. Keep #82 health/recovery/load shedding separate.

## Resume rule
1. Fetch live `main` and live `adapt-provider-router-production-001` first; another session may have advanced them.
2. Read this file, `governance/current-state.json`, issue #79 latest comments, issue #81 and its comments, the #81 work-slice/closure ledger, and #80 final qualification evidence.
3. Inspect commits since `c75a5f1467920f57fa23c3dbc400e51edc5275c8` before changing code so implemented work is never duplicated.
4. Continue actual #81 implementation from the exact current head; do not restart planning.
5. Preserve Smart Provider Router v2, direct SEC/EDGAR, canonical freshness/cache/persistence/telemetry/state owners, No Execution and U.S.-equities boundaries.
6. Use only canonical Fast, Qualified and Release workflows for qualification. No temporary workflow family and no gate weakening.
7. Do not merge #81 without exact-head Fast + impact-selected Qualified; do not start #83/#78 before required #81/#82 closure.
