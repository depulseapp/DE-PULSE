# DE.PULSE Current Handoff

**Canonical machine state:** `governance/current-state.json`  
**Certified Stable:** `v18.9.1-stable`  
**Certified Stable candidate:** `e55d8d25b15cec2ffb0f5411bc358bc40b359cf9`  
**Certified Stable source fingerprint:** `0062f46dea5690d0b3fcd8a9ed3b1f71ebe1522c7dee2cb218e9d36b9e0076ff`  
**Certified Stable build ID:** `v18.9.1-stable-20260821`  
**Retained completed process authority:** #73 / `ADAPT-ROOT-CONVERGENCE-001` / `adapt-root-convergence-001`  
**Process closure ledger:** `governance/work-slices/ADAPT-ROOT-CONVERGENCE-001/closure.json`  
**Completed product slice:** #76 / `ADAPT-TRADEINSIGHT-SETTINGS-001` / PR #77 / merge `a171ce2258632bd4bd6aa737176f2d6dffb44689`  
**Active product slice:** #80 / `ADAPT-DATAHEALTH-BASELINE-001` / `adapt-datahealth-baseline-001`  
**Parent program:** #79 / `ADAPT-PROVIDER-PRODUCTION-001`

## Current authority
Issue #79 is the cross-session authority for the all-provider Adaptive Data Health program. Issue #80 is the current executable foundation slice. The required dependency order remains **#80 → #81/#82 → #83 + #78 → #84**.

PR #77 is landed and issue #76 is complete. Its immutable implementation evidence is retained at `governance/work-slices/ADAPT-TRADEINSIGHT-SETTINGS-001/final-qualification-evidence.json`; no public Stable release was created by that product merge.

## #80 objective
Build an executable, exhaustive provider/capability/data-source baseline for the whole application: provider/capability matrix, authority/freshness/cache/fallback/consumer/materiality ownership, Data Health SLOs, runtime fetch-path bypass dispositions, scoped truthful degradation/recovery rules, and a recurrence gate that rejects unclassified providers/capabilities/external fetch paths.

Smart Provider Router v2 remains the sole **general** routing/admission authority. Direct SEC/EDGAR retains Form 4 authority. Legitimate direct-authority/public sources must be classified explicitly rather than rank-swapped through the general router. Reuse canonical freshness, cache, persistence, telemetry, state and validation owners; do not create parallel provider-specific Data Health subsystems.

## Permanent boundaries
- U.S. equities processing only.
- No execution/order routing.
- GLD, SLV and USO remain actionable live-priority exceptions.
- TradeInsight remains shadow-first where governed; direct SEC/EDGAR remains authoritative for Form 4.
- GitHub executable evidence outranks chat memory.

## Resume rule
1. Fetch live `main` and the live head of `adapt-datahealth-baseline-001` first; another session may have advanced them.
2. Read this file, `governance/current-state.json`, issue #79 latest comments, issue #80 and its comments, and the #80 work-slice/closure ledger.
3. Inspect commits since `a171ce2258632bd4bd6aa737176f2d6dffb44689` before changing code so implemented work is never duplicated.
4. Continue actual #80 implementation from the exact current head. Do not restart planning from scratch.
5. Use only canonical Fast, Qualified and Release workflows for qualification. No temporary workflow family and no gate weakening.
6. Do not start #81/#82 until #80 executable acceptance and exact-head evidence are complete.
