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
**Active branch:** `main` after this continuity reconciliation merges.  
**Next release line:** `v18.8.2 — Market Intelligence Reliability / ADAPT-FRESHNESS-001 REOPENED`.

## v18.8.1 Stable closure

PR #56 qualified exact source head `07624965519cdd406c6db1e19771cf75dec825b4`, then merged to candidate `410679ba0d6459f66a44db15a0a55f30741a7c53`. Fast #425, Qualified #146 and Release #29 passed. Release #29 completed G11, G12, macOS Apple Silicon + Windows x64 G13/G14 actual packaged-runtime audits, G15 Release Assurance, exact same-run no-rebuild publication and G16 workflow evidence.

Durable Stable evidence is now reconciled into:
- `.depulse-certification/resume/build-checkpoint.json`;
- `.depulse-certification/resume/release-evidence-checkpoint.json`;
- `release/v18.8.1/stable-evidence-manifest.json`;
- this handoff and the four CURRENT Adaptive overlays.

The continuity commit is fingerprint-excluded metadata/process hardening. It does not rebuild, republish or redefine the immutable v18.8.1 binaries/tag.

## Historical v18.8.1 G1 authority

v18.8.1 used `release/v18.8.1/release_contract.json` plus `adaptive-governance/V18.8.1-ZERO-MISS-RECONCILIATION.json` as its actual frozen scope/closure authority. No retrospective `G1-IMMUTABLE-SCOPE.md` is fabricated. Future release lines must freeze an explicit G1 scope authority before qualification.

## Escaped runtime finding after Stable

GitHub issue **#57** is OPEN: Market Intelligence can remain broadly `DATA DEGRADED` in pre-market, regular-session and after-hours operation, with Market Tradeability unavailable/`0/100`, SPY/QQQ/VIX evidence absent and tracked breadth `0/15`.

Classification: **OPEN_DEFECT / NEXT_RELEASE_MANDATORY_ENTRY** for v18.8.2. `ADAPT-FRESHNESS-001` is reopened for the affected Market Intelligence consumer-demand/freshness path. This is an ADR-GDI/Smart Provider Router v2/canonical freshness refinement, not a new data engine.

Required architecture:
- Smart Provider Router v2 remains the sole provider-routing authority;
- canonical freshness/recovery remains the sole freshness owner;
- existing live/snapshot allocation remains canonical;
- Market Intelligence benchmark/breadth demand becomes first-class canonical demand;
- SPY/QQQ/VIX receive protected market-context priority and breadth symbols receive bounded lower priority;
- unavailable/unknown evidence must not masquerade as observed numeric zero;
- true provider outages/rate limits/staleness remain truthfully degraded.

Minimum v18.8.2 evidence: all-session behavior, partial breadth, VIX-only failure, total acquisition failure, fallback/recovery, truthful unavailable-vs-zero UI, browser/runtime retest and actual macOS Apple Silicon + Windows x64 package proof.

## Roadmap continuity

1. v18.8.2 — bounded Market Intelligence reliability repair for issue #57 / reopened `ADAPT-FRESHNESS-001`.
2. v18.9.0 — `ADAPT-TRADEINSIGHT-001` full beta capability discovery, utility mapping and SHADOW integration through Smart Provider Router v2.
3. v18.9.1 — Provider Intelligence & Market-Mode hardening.
4. v18.10.0 — zero-gap v18 Major Closure Candidate.

Permanent boundaries remain U.S. Equities Processing, No Execution, deterministic Day/Swing/Long truth, Smart Provider Router v2 sole routing authority, BroadSnapshotBroker canonical reuse owner and GLD/SLV/USO actionable tradable exceptions.

## Post-Stable continuity prevention

A cheap post-Stable continuity sentinel now treats `main` carrying a later STABLE package identity than the durable Stable checkpoint as a failure. The sentinel also requires aligned Stable evidence manifest, current handoff and CURRENT Adaptive overlays when repository identity is Stable-aligned. This prevents a successful release workflow from silently leaving GitHub resume truth on the prior Stable.

## Resume rule

Read `AGENTS.md` or `CLAUDE.md`, `governance/AI-ASSISTANT-PORTABILITY-CONTRACT.md`, this handoff, both `.depulse-certification/resume/` checkpoints, `release/v18.8.1/stable-evidence-manifest.json`, the four CURRENT Adaptive overlays, issue #57 and actual GitHub branch/PR/workflow/release state. Run `python3 adaptive_resume_gate.py` and `python3 tools/ci/post_stable_continuity_gate.py`. Never resume from model memory alone.

## Exactly one next action

Create `v18.8.2-development` from the reconciled `main` baseline and execute G0 exact-baseline diagnosis for issue #57 before changing product source or beginning v18.9.0 TradeInsight work.
