# DE.PULSE — Current Authoritative Handoff

**SUPERSEDES ALL PRIOR CHAT HANDOFFS. GitHub is authoritative; model memory is advisory only.**

## Current identity

**Certified Stable:** `v18.8.1-stable`  
**Certified candidate:** `410679ba0d6459f66a44db15a0a55f30741a7c53`  
**Qualified source head:** `07624965519cdd406c6db1e19771cf75dec825b4`  
**Certified source fingerprint:** `bfefa3605ab29b4678275936a3e60e45133d0b592b91298551731f6d629a9d92`  
**Certified Build ID:** `v18.8.1-stable-20260820`  
**Stable Fast:** #425 / `32414640774`  
**Stable Qualified:** #146 / `32415313893`  
**Stable Release:** #29 / `32415750821`  
**Stable source PR:** #56  
**Active development branch:** `v18.8.2-development`  
**Active release line:** `v18.8.2 — Market Intelligence Reliability / ADAPT-FRESHNESS-001 REOPENED`.

## v18.8.1 Stable authority

v18.8.1 remains immutable. PR #56 qualified source head `07624965519cdd406c6db1e19771cf75dec825b4`, then merged to candidate `410679ba0d6459f66a44db15a0a55f30741a7c53`. Fast #425, Qualified #146 and Release #29 passed, including macOS Apple Silicon + Windows x64 actual packaged-runtime audits, G15 Release Assurance, same-run no-rebuild publication and G16 workflow evidence.

Durable Stable evidence remains in:
- `.depulse-certification/resume/build-checkpoint.json`;
- `.depulse-certification/resume/release-evidence-checkpoint.json`;
- `release/v18.8.1/stable-evidence-manifest.json`.

The v18.8.2 branch does not redefine or rebuild v18.8.1 Stable.

## v18.8.2 issue #57 — G0–G4 status

GitHub issue **#57** is OPEN. Runtime symptom: Market Intelligence can remain broadly `DATA DEGRADED` across sessions, with SPY/QQQ/VIX evidence absent and breadth `0/15`.

### G0 exact-baseline diagnosis — COMPLETE

The live/snapshot allocator already owns the required market context. `multiFeedAllocationWithHints` protects SPY/QQQ in live Tier 0 and assigns every canonical master market symbol to either live or snapshot demand. The 15-symbol breadth universe is already contained in canonical master market symbols. VIX remains on its existing canonical special-index path.

The escaped defect is therefore **not missing provider demand and not a need for another router/data engine**. The gap is canonical freshness/recovery accountability: `Engine.Snapshot()` previously scoped quote freshness mostly to active desk symbols + selected Research. A transient missing/stale Market Intelligence quote could therefore fall outside the sole freshness row and never become targeted recovery demand.

### G1–G3 — FROZEN ON ISSUE #57

Bounded scope and ownership are frozen in issue #57 comments:
- Smart Provider Router v2 remains sole provider-routing authority;
- existing multi-feed live/snapshot allocation remains sole subscription allocation owner;
- canonical freshness/recovery + existing routed refresh remain sole recovery owners;
- no Market Intelligence-specific provider fetch loop;
- no second data engine/freshness system/subscription manager;
- deterministic Day/Swing/Long truth and No Execution remain unchanged;
- TradeInsight is excluded from v18.8.2.

### G4 implementation — COMPLETE; exit validation pending Fast

Current product/test diff versus `main` is intentionally bounded:
- `engine_core.go`: canonical quote freshness scope now includes the existing `broadBreadthUniverse`, deduped with active desk/Research demand. This makes missing/stale Market Intelligence breadth quotes visible to the existing canonical freshness/recovery loop.
- `renderer/market-intelligence-truth.js`: presentation-only truth reconciliation renders Market Tradeability score as `UNAVAILABLE` when state is `DATA DEGRADED`/`UNAVAILABLE`; evaluated states can still legitimately show numeric `0/100`.
- `renderer/index.html`: loads the v18.8.2 truth layer after the existing renderer; no navigation or other UI structure change remains.
- `v18_8_2_market_intelligence_freshness_test.go`: deterministic regressions cover allocator ownership of all 15 breadth symbols, protected SPY/QQQ, missing/recovered breadth quote freshness, SPY/QQQ/VIX missing individually, stale core evidence, and 0/15 remaining `UNAVAILABLE` rather than directional numeric truth.
- `tests/renderer/surface_consolidation_test.js`: existing Fast Node lane now validates unavailable-vs-zero truth and script wiring. No new workflow/job was added.

A temporary standalone renderer test was deleted after its coverage was consolidated into the existing Fast renderer lane. No retry/certification branch exists and no product workflow has been duplicated.

Important implementation refinement: the public Market Intelligence breadth contract remains the canonical 15-symbol universe. All 15 are already allocator-owned live or snapshot demand, so the repair does **not** alter allocator/provider priority or create unconditional independent polling. It only makes this existing canonical demand freshness-accountable and recoverable through existing owners.

## Required remaining v18.8.2 evidence

Before closure, prove:
- Fast exact-head qualification of the bounded Go + renderer regressions;
- pre-market, regular, after-hours behavior;
- partial breadth and complete acquisition failure truth;
- VIX-only failure;
- Smart Provider Router fallback/recovery behavior without bypass;
- browser/runtime unavailable-vs-zero presentation;
- performance/backpressure/provider-budget integrity;
- actual macOS Apple Silicon + Windows x64 packaged-runtime proof through normal G11–G16 delivery.

## Roadmap continuity

1. v18.8.2 — bounded Market Intelligence reliability repair for issue #57 / reopened `ADAPT-FRESHNESS-001`.
2. v18.9.0 — `ADAPT-TRADEINSIGHT-001` full beta capability discovery, utility mapping and SHADOW integration through Smart Provider Router v2.
3. v18.9.1 — Provider Intelligence & Market-Mode hardening.
4. v18.10.0 — zero-gap v18 Major Closure Candidate.

Permanent boundaries remain U.S. Equities Processing, No Execution, deterministic Day/Swing/Long truth, Smart Provider Router v2 sole routing authority, BroadSnapshotBroker canonical reuse owner and GLD/SLV/USO actionable tradable exceptions.

## Resume rule

Read `AGENTS.md` or `CLAUDE.md`, `governance/AI-ASSISTANT-PORTABILITY-CONTRACT.md`, this handoff, both `.depulse-certification/resume/` checkpoints, `release/v18.8.1/stable-evidence-manifest.json`, the four CURRENT Adaptive overlays, issue #57 and actual GitHub branch/PR/workflow/release state. Never resume from model memory alone.

## Exactly one next action

Open exactly one **Draft** PR from `v18.8.2-development` to `main` after the final diff check, then evaluate the automatically triggered CI Fast result on that exact PR head. Do not push another branch commit after PR creation unless Fast identifies a real defect that requires correction; do not create retry/certification branches or manually duplicate CI.
