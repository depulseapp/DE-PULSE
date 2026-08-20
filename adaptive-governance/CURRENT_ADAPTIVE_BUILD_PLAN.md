# DE.PULSE — Current Adaptive Build Plan

**Immutable Stable:** `v18.8.1-stable`  
**Stable candidate:** `410679ba0d6459f66a44db15a0a55f30741a7c53`  
**Stable qualified source:** `07624965519cdd406c6db1e19771cf75dec825b4`  
**Next engineering branch:** `v18.8.2-development` (create only after post-Stable continuity reconciliation is on `main`)  
**Open defect:** GitHub issue #57  
**Current line:** `v18.8.2 — Market Intelligence Reliability`.

## Entry condition

v18.8.1 post-Stable continuity must be aligned across checkpoints, Stable evidence manifest, handoff and CURRENT Adaptive overlays. The v18.8.1 Stable binary/tag remains immutable and is not rebuilt by continuity metadata.

## v18.8.2 mandatory work

`ADAPT-FRESHNESS-001` is **REOPENED** for issue #57. Preserve its original requirement lineage rather than creating a disconnected reliability engine.

At G0 capture the exact all-session failure: Market Tradeability, SPY/QQQ/VIX evidence, tracked breadth 15-symbol coverage, provider/source timestamps, freshness acceptance, live/snapshot allocation and recovery state.

At G1 freeze the bounded affected surface and acceptance matrix. At G2/G3 prove one canonical owner chain:

`Market Intelligence demand → canonical demand/allocation → Smart Provider Router v2 → canonical quote/evidence-time/freshness/recovery → Market Intelligence consumers`.

Required implementation behavior:
- Market Intelligence required benchmark/breadth symbols participate in canonical demand/recovery;
- SPY/QQQ/VIX are protected market-context demand; breadth remains bounded/lower priority;
- no blind independent polling or duplicate provider path;
- unknown/unavailable does not become `0/100`, `0/15` or `0%` unless zero is genuinely observed;
- optional context cannot contaminate unrelated consumers;
- true required-evidence failure remains degraded/ABSTAIN as appropriate;
- provider budgets, circuits, backpressure and calls-avoided semantics remain intact.

Required tests include pre-market, regular session, post-market, partial breadth, VIX-only failure, total acquisition failure, stale evidence, provider fallback, recovery hysteresis and renderer unavailable-vs-zero semantics.

Qualification/delivery uses the normal single branch → one Draft PR/Fast → same PR Ready/Qualified → exact-head merge → one G11–G16 release path. Actual macOS Apple Silicon and Windows x64 package proof is mandatory for affected runtime behavior.

## After v18.8.2

v18.9.0 `ADAPT-TRADEINSIGHT-001` → v18.9.1 provider/Market-Mode hardening → v18.10 zero-gap closure.

## Exactly one next action

Create `v18.8.2-development` from reconciled `main` and execute G0 exact-baseline diagnosis for issue #57 before product code changes.
