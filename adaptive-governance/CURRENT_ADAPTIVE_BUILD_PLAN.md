# DE.PULSE — Current Adaptive Build Plan

**Immutable Stable:** `v18.8.1-stable`  
**Stable candidate:** `410679ba0d6459f66a44db15a0a55f30741a7c53`  
**Stable qualified source:** `07624965519cdd406c6db1e19771cf75dec825b4`  
**Current engineering branch:** `v18.8.2-development`  
**Open defect:** GitHub issue #57  
**Current line:** `v18.8.2 — Market Intelligence Reliability`.

## Entry condition — PASS

v18.8.1 post-Stable continuity is aligned across checkpoints, Stable evidence manifest, handoff and CURRENT Adaptive overlays. The v18.8.1 Stable binary/tag remains immutable and is not rebuilt by v18.8.2 development.

## v18.8.2 mandatory work

`ADAPT-FRESHNESS-001` is **REOPENED** for issue #57. Preserve its original requirement lineage rather than creating a disconnected reliability engine.

### G0 — COMPLETE

Exact-baseline diagnosis proved that the canonical live/snapshot allocator already owns SPY/QQQ and all 15 Market Intelligence breadth symbols. The escaped defect is the narrower freshness/recovery accountability gap: Market Intelligence quote demand was not included in the sole scoped quote-freshness row after transient missing/stale evidence.

### G1 — FROZEN

Scope is issue #57 only. No TradeInsight/provider expansion, no execution, no Day/Swing/Long semantic change, no second router, no second freshness system, no second data engine and no duplicate subscription manager.

### G2/G3 — FROZEN

Canonical owner chain remains:

`Existing Market Intelligence/master-market demand → existing multi-feed live/snapshot allocation → Smart Provider Router v2 → canonical quote/evidence-time/freshness/recovery → Market Intelligence consumers`.

SPY/QQQ remain protected by existing allocator priority; VIX remains on its canonical special-index path. Breadth remains the existing 15-symbol bounded market-context universe and is already admitted by the canonical allocator as live or snapshot demand. The repair therefore changes freshness accountability, not provider routing or subscription priority.

### G4 — IMPLEMENTATION COMPLETE; EXIT PENDING FAST

Bounded implementation:
- add existing `broadBreadthUniverse` to canonical quote freshness scope in `Engine.Snapshot()`, deduped with active desk/Research symbols;
- leave `multiFeedAllocationWithHints`, Smart Provider Router v2 and routed refresh ownership unchanged;
- add deterministic regressions for allocator ownership, missing/recovered breadth freshness, SPY/QQQ/VIX individually missing, stale evidence and 0/15 unavailable truth;
- add presentation-only truth reconciliation so `DATA DEGRADED`/`UNAVAILABLE` does not render as a meaningful `0/100`;
- integrate renderer assertions into the existing Fast Node lane rather than creating a new workflow/job.

No post-PR branch push is planned unless Fast exposes a real defect.

## Remaining qualification matrix

G5 Fast must prove Go formatting/vet/full suite, renderer syntax and existing renderer contract including v18.8.2 unavailable-vs-zero assertions.

G6-G10 must then prove integration/data truth/provider fallback/recovery/performance-backpressure/browser behavior and exact-head readiness. Required scenarios include pre-market, regular session, post-market, partial breadth, VIX-only failure, total acquisition failure, stale evidence, Smart Router fallback and recovery-to-current.

G11-G16 remains the normal immutable RC/certification/native-package/actual-artifact/assurance/publication/retrospective path. Actual macOS Apple Silicon and Windows x64 package proof is mandatory.

## After v18.8.2

v18.9.0 `ADAPT-TRADEINSIGHT-001` → v18.9.1 provider/Market-Mode hardening → v18.10 zero-gap closure.

## Exactly one next action

Open exactly one Draft PR from `v18.8.2-development` to `main` and inspect the automatically triggered CI Fast result on that exact head. No retry/certification branches and no duplicate manual Fast run.
