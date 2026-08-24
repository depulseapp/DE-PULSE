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
**Completed Data Health baseline:** #80 / `ADAPT-DATAHEALTH-BASELINE-001` / PR #86 / merge `c75a5f1467920f57fa23c3dbc400e51edc5275c8`  
**Completed Router production adoption:** #81 / `ADAPT-PROVIDER-ROUTER-PRODUCTION-001` / candidate `60d63bb862525edea5a5ea8be7469778f44afc54` / Fast #878 / Qualified #184 / PR #87 / merge `1870dd3881dbe7f6463f242e35fdc19e70d9ae15`  
**Active product slice:** #82 / `ADAPT-DATAHEALTH-RUNTIME-001` / `adapt-datahealth-runtime-001`  
**Parent program:** #79 / `ADAPT-PROVIDER-PRODUCTION-001`

## Current authority
Issue #79 is the cross-session authority for the all-provider Adaptive Data Health program. #80 and #81 are complete. Their machine-readable provider/capability/SLO/fetch-path baseline, Router v2 adoption, regression evidence and immutable qualification evidence remain authoritative inputs to #82. Current dependency order is **#82 → #83 + #78 → #84**.

## #81 completed outcome
Smart Provider Router v2 is now the executable general routing/admission authority for the former #80 `MIGRATE` paths:
- Alpaca canonical U.S. asset universe, with bounded live/paper same-provider endpoint fallback inside the Alpaca attempt;
- Alpaca read-only U.S. market calendar;
- Alpaca corporate-actions announcements;
- Twelve Data global indices/futures context.

Capability-scoped failure/cancellation isolation is regression-protected, existing cache/single-flight/coalescing and provider workload owners remain in place, and the fetch-path ledger contains no surviving #80 general-routing `MIGRATE` bypass. Direct SEC/EDGAR Form 4 and other explicit official-authority paths remain unchanged.

## #82 objective
Make `PARTIAL COVERAGE` / `DATA DEGRADED` truthful, capability-scoped and short-lived when recoverable. Reuse the existing canonical freshness/provider/cache/persistence/telemetry/state owners and Smart Provider Router v2; do not create a parallel health subsystem.

Required runtime outcomes:
- one canonical data-health evaluation path scoped by capability/consumer/symbol/session before wider escalation;
- reuse defensible warm/cached/persisted evidence when freshness policy permits;
- automatic eligible fallback, revalidation and hysteresis-protected recovery/unlatch;
- critical decision-support work protected before optional/background work;
- bounded scanner/prep/event/research/background fan-out under provider/runtime pressure;
- lower-value load shedding before core data health deteriorates;
- no optional provider/capability failure causing false app-global degradation;
- truthful degraded state when required evidence is stale, missing, contradictory or below quality;
- non-secret telemetry for reason, scope, freshness, serving/preferred provider context, fallback/recovery state and duration.

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
On `adapt-datahealth-runtime-001`, inspect the existing runtime degradation/reliability/workload owners and implement #82 by consolidation/reuse: first bind one canonical scoped data-health evaluation and fault matrix, then close warm-state/fallback/recovery/hysteresis and pressure/load-shedding gaps without changing Router v2 authority. Obtain exact-head Fast and impact-selected Qualified before merge. Do not start #83/#78 until #82 closes.

## Resume rule
1. Fetch live `main` and live `adapt-datahealth-runtime-001` first; another session may have advanced them.
2. Read this file, `governance/current-state.json`, issue #79 latest comments, issue #82 and its comments, the #82 work-slice/closure ledger, and #81 final qualification evidence.
3. Inspect commits since `1870dd3881dbe7f6463f242e35fdc19e70d9ae15` before changing code so implemented work is never duplicated.
4. Continue actual #82 implementation from the exact current head; do not restart planning.
5. Preserve Smart Provider Router v2, direct SEC/EDGAR, canonical freshness/cache/persistence/telemetry/state owners, No Execution and U.S.-equities boundaries.
6. Use only canonical Fast, Qualified and Release workflows for qualification. No temporary workflow family and no gate weakening.
7. Do not merge #82 without exact-head Fast + impact-selected Qualified; do not start #83/#78 before required #82 closure.
