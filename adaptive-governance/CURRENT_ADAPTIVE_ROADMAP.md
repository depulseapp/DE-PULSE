# CURRENT Adaptive Roadmap

**Canonical rebaseline:** `governance/V19_V20_REBASELINE.md`  
**Backlog map:** `governance/programs/V19-V20-REBASELINE/backlog-version-matrix.json`  
**HOST map:** `governance/programs/V19-V20-REBASELINE/host-requirement-version-map.json`  
**Legacy future-commitment conservation:** `governance/programs/V19-V20-REBASELINE/legacy-future-commitment-conservation.json`  
**Cross-integration map:** `governance/programs/V19-V20-REBASELINE/cross-integration-matrix.json`  
**Whole-product surface map:** `governance/programs/V19-V20-REBASELINE/whole-product-surface-rebaseline.json`  
**Certified Stable:** `v18.10.0` — immutable  
**Active version:** `v19.0.0` — Hosted Trust & Identity Foundation  
**Active issue/PR:** #148 / PR #149  
**Active branch:** `adapt-hosted-trust-foundation-001`

The prior future micro-version/packet scheduling is superseded. Requirements remain individually conserved and evidenced, but current planning is by coherent **version/build**.

The rebaseline now conserves four distinct requirement families: the 180 certified v18 responsibilities, all 72 HOST requirements, all 22 testing/future backlog issues, and legacy future-roadmap/build-plan commitments that are not fully represented by HOST/backlog IDs. No family may be assumed covered merely because a version title sounds related.

## Current v19 sequence

1. `v19.0.0` — Hosted Trust & Identity Foundation — HOST-001..023 + core #164/#156 security/auth overlap.
2. `v19.1.0` — Canonical Data Runtime & Global Symbol Processing — #150/#151/#153/#154/#155 runtime truth/#160 integration/#167 core.
3. `v19.2.0` — Hosted Gateway, Shared Serving & Sync Core — HOST-024..039.
4. `v19.3.0` — Cross-Platform Product, Roles & Information Architecture — HOST-040..047/053 + #152/#156/#159/#160 UI verification/#167 Admin/#171 + #164 UX parity + `LEGACY-TRADER-SETUP-SHORT-001` two-sided deterministic Desk setup contract.
5. `v19.4.0` — Market Intelligence & Research Workflow Quality — HOST-049 + #158/#161/#162/#171.
6. `v19.4.1` — Discovery & Opportunity Radar Effectiveness — HOST-048 + #163/#171.
7. `v19.5.0` — Price/Volume & Event-Anchored Intelligence — #168/#169.
8. `v19.5.1` — Options Structure & GEX Intelligence — #157.
9. `v19.6.0` — Point-in-Time Evidence & Outcome-Ready Foundation — HOST-057..064 + deterministic #165 foundation + institutional/two-sided thesis substrate.
10. `v19.6.1` — Hosted Reliability, Economics & Adaptive Readiness — HOST-050..056/065..071 + ADR-GDI/trader-quality readiness + final #170/#171 reconciliation.
11. `v19.7.0` — v19 Major Closure — HOST-072; no feature scope.

### Two-sided Desk setup conservation

Current source audit proves this is **not fully implemented today**. `computePlan` can label a low score `BEARISH`, but it still builds long-oriented plan geometry: targets above price/entry, invalidation below entry, and long-oriented `actionState` comparisons. Therefore the existing Entry Zone / Trim-Target / Invalidation / Setup Score strip is not a true SHORT plan.

`v19.3.0` must therefore establish one deterministic two-sided setup contract across applicable Day/Swing/Long horizons while preserving the approved Desk look-and-feel:
- explicit `LONG / SHORT / NO_SETUP-WAIT` side;
- setup-quality score independent from bullish/bearish direction, so a strong SHORT setup may have a high quality score;
- LONG entry -> target above -> invalidation below;
- SHORT entry -> cover/target below -> invalidation above;
- direction-aware positive R-multiple and action-state comparisons;
- truthful SHORT labels such as `Cover / Target` where useful;
- Research/Discovery/replay/outcome consumers use the same canonical side/geometry;
- bearish evidence alone never forces a SHORT setup when readiness/quality is insufficient;
- No Execution remains permanent.

This is distinct from v19.6/v20.3 institutional **two-sided thesis/TDTI** evidence.

## Current v20 sequence

1. `v20.0.0` — Outcome Learning & Adaptive Control Plane — #165 learning foundation plus conserved calibration/guardrails.
2. `v20.1.0` — Adaptive Chart Pattern & Similarity Intelligence — #166.
3. `v20.2.0` — Adaptive Market Synthesis, Regime & Discovery Learning — includes conserved ASBI normalization/synthesis/contradiction/outcome behavior on top of v20.0 controls.
4. `v20.3.0` — Adaptive Institutional & Two-Sided Thesis Intelligence.
5. `v20.3.1` — AODR Adaptive Opportunity Intelligence.
6. `v20.4.0` — Agent Orchestration & Controlled MCP/API — remaining #165.
7. `v20.5.0` — Adaptive Operations.
8. `v20.6.0` — Professional Adaptive Closure; no feature scope.

## Permanent intelligence direction

`canonical evidence -> deterministic intelligence -> cross-feature synthesis -> point-in-time outcomes -> bounded adaptive learning -> optional AI/agent explanation/orchestration`

#170 is a mandatory cross-integration/Market-Regime gate for every intelligence-bearing version. #171 is a mandatory UI/data-density/intelligence-maturity audit for every changed user-visible surface. Neither is a standalone isolated feature.

Hard UI protections remain: keep Day/Swing/Long Desk look-and-feel/workflow; keep Dashboard Market Regime and Desk Control materially intact; preserve Data Engine look-and-feel; preserve the current AI Copilot engine/header visual treatment unless separately justified.

## Current exact state / next action

The earlier product head `c5d0713d16f95522fd013123a78bc7cc58dc2422` was not qualified: Fast #1141 failed source-health because five hosted identity/session helpers lacked production references. Rebaseline governance work does not close that product defect. Fetch the live PR head/current checks, continue `v19.0.0` by wiring those capabilities through existing authenticated identity/HTTP owners (or removing/reworking helpers proven unnecessary), then obtain exact-head Fast. Do not start `v19.1.0` before `v19.0.0` exit criteria pass.
