# CURRENT Adaptive Roadmap

**Canonical rebaseline:** `governance/V19_V20_REBASELINE.md`  
**Provider-registry additive rebaseline:** `governance/V19_V20_PROVIDER_REGISTRY_REBASELINE.md`  
**Provider-registry permanent contract:** `adaptive-governance/ADAPTIVE_PROVIDER_REGISTRY_CONTRACT.md`  
**Machine current-state authority:** `governance/current-state.json`  
**Backlog map:** `governance/programs/V19-V20-REBASELINE/backlog-version-matrix.json`  
**HOST map:** `governance/programs/V19-V20-REBASELINE/host-requirement-version-map.json`  
**Provider-registry machine map:** `governance/programs/V19-V20-REBASELINE/adaptive-provider-registry.json`  
**Legacy future-commitment conservation:** `governance/programs/V19-V20-REBASELINE/legacy-future-commitment-conservation.json`  
**Cross-integration map:** `governance/programs/V19-V20-REBASELINE/cross-integration-matrix.json`  
**Whole-product surface map:** `governance/programs/V19-V20-REBASELINE/whole-product-surface-rebaseline.json`  
**Certified Stable:** `v18.10.0` — immutable  
**Active version:** `v19.0.0` — Hosted Trust & Identity Foundation  
**Active work slice:** `ADAPT-HOSTED-TRUST-FOUNDATION-001`  
**Active issue/PR:** #148 / PR #149  
**Active branch:** `adapt-hosted-trust-foundation-001`  
**Canonical closure ledger:** `governance/work-slices/ADAPT-HOSTED-TRUST-FOUNDATION-001/closure.json`

The prior future micro-version/packet scheduling is superseded. Requirements remain individually conserved and evidenced, but current planning is by coherent **version/build**.

The rebaseline now conserves four distinct requirement families: the 180 certified v18 responsibilities, all 72 HOST requirements, all 22 testing/future backlog issues, and legacy future-roadmap/build-plan commitments that are not fully represented by HOST/backlog IDs. No family may be assumed covered merely because a version title sounds related.

## Current v19 sequence

1. `v19.0.0` — Hosted Trust & Identity Foundation — HOST-001..023 + core #164/#156 security/auth overlap.
2. `v19.1.0` — Canonical Data Runtime & Global Symbol Processing — #150/#151/#153/#154/#155 runtime truth/#160 integration/#167 core + Adaptive Provider Registry / Market Data first adoption.
3. `v19.2.0` — Hosted Gateway, Shared Serving & Sync Core — HOST-024..039.
4. `v19.3.0` — Cross-Platform Product, Roles & Information Architecture — HOST-040..047/053 + #152/#156/#159/#160 UI verification/#167 Admin/#171 + #164 UX parity + `LEGACY-TRADER-SETUP-SHORT-001` two-sided deterministic Desk setup contract + role-aware Mac/Windows/Web provider Settings presentation using the generic registry metadata contract.
5. `v19.4.0` — Market Intelligence & Research Workflow Quality — HOST-049 + #158/#161/#162/#171.
6. `v19.4.1` — Discovery & Opportunity Radar Effectiveness — HOST-048 + #163/#171.
7. `v19.5.0` — Price/Volume & Event-Anchored Intelligence — #168/#169.
8. `v19.5.1` — Options Structure & GEX Intelligence — #157.
9. `v19.6.0` — Point-in-Time Evidence & Outcome-Ready Foundation — HOST-057..064 + deterministic #165 foundation + institutional/two-sided thesis substrate.
10. `v19.6.1` — Hosted Reliability, Economics & Adaptive Readiness — HOST-050..056/065..071 + ADR-GDI/trader-quality readiness + final #170/#171 reconciliation + provider reliability/coverage/economics/readiness scorecards across the full registry including Market Data.
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

## Conserved provider / Data Health program

The completed Adaptive Data Health/provider-production chain remains a permanent cross-version contract and must not disappear from CURRENT planning merely because its original issues are closed:

- #80 — executable provider/capability baseline, SLOs, ownership and fetch-path classification.
- #81 — Smart Provider Router v2 production adoption for every routable capability, while preserving explicitly classified direct-authority paths.
- #82 — one canonical runtime Data Health evaluation path with capability/symbol/consumer scoping, eligible cache/fallback recovery, hysteresis, workload priority and load shedding.
- #83 — common capability lifecycle/readiness evidence (`SHADOW -> VALIDATED -> APPROVED -> PRODUCTION`) with no automatic authority promotion.
- #78 — TradeInsight production readiness is one provider-adoption case inside #81/#83, not a parallel router, freshness, cache, telemetry or lifecycle subsystem.
- #84 — zero-gap fault injection, native proof and professional closure for the program.

Canonical routing authority remains **Smart Provider Router v2** for general routable provider capabilities. Explicit direct authorities remain explicit; in particular direct SEC/EDGAR remains authoritative for Form 4 and other source-authority cases defined by the provider matrix/fetch-path contracts.

Data-health truth is evidence-time based: provider observation/event/publication/filing time is authoritative; retrieval/cache time is bookkeeping; unknown observation time remains unknown. Valid warm/cache evidence and eligible fallback are reused before avoidable degradation. Optional-provider failure cannot contaminate unrelated symbols, desks, Market Modes or application-global health.

`PARTIAL COVERAGE` and `DATA DEGRADED` remain truthful, attributable and minimally scoped states for genuine unresolved required-evidence gaps. Recovery is automatic when canonical evidence becomes healthy under policy, with hysteresis and authority rules preserved; the product must never suppress a genuine data problem merely to appear healthy.

## Adaptive Provider Registry / future-proof provider onboarding

The approved provider-onboarding direction is now a permanent extension of the conserved Data Health program:

`Provider Adapter -> Adaptive Provider Registry -> capability/entitlement probes -> rights/authority + Data Health -> Smart Provider Router v2 -> canonical state -> all useful consumers`

Permanent roadmap rules:
- provider-specific implementation ends at one standards-compliant adapter;
- the adapter self-registers with the Registry;
- capability/configuration/effective-entitlement/freshness/history/quota/health observations become generic Router eligibility inputs where technically observable;
- consumers request capabilities, never provider names;
- technical eligibility, fallback, demotion, cooldown, recovery and subscription-plan reprobe may be automatic;
- `SHADOW -> VALIDATED -> APPROVED -> PRODUCTION` authority promotion remains governed and never automatic;
- direct-authority sources and public/commercial rights cannot be replaced/inferred automatically;
- every provider capability receives mandatory #170 cross-integration dispositions across applicable Research, Discovery/Radar, Desks, Prep, Market Intelligence/Regime contribution, alerts, history/options, Data Health/Maintenance and future Outcome Learning consumers;
- provider-card/credential Settings UX becomes metadata-driven enough that adding another token/API-key adapter does not require another Settings architecture;
- Market Data (`marketdata.app`) is the first concrete adopter, using the reusable secret/config/test/redaction pattern learned from TradeInsight #76;
- effective Market Data entitlement changes (for example delayed trial/free to future paid/live access) should update technical eligibility through reprobe rather than require a DE.PULSE code release solely because the provider subscription changed;
- v19.6.1 measures provider reliability/coverage/economics/readiness; v20.5 may add bounded learned provider-utility/cost priors, but Smart Provider Router v2 remains the sole routing authority.

This approved future scope is implemented primarily in v19.1/#153 and may not be pulled forward to bypass unfinished v19.0 Development Production Ready work.

## Current v20 sequence

1. `v20.0.0` — Outcome Learning & Adaptive Control Plane — #165 learning foundation plus conserved calibration/guardrails.
2. `v20.1.0` — Adaptive Chart Pattern & Similarity Intelligence — #166.
3. `v20.2.0` — Adaptive Market Synthesis, Regime & Discovery Learning — includes conserved ASBI normalization/synthesis/contradiction/outcome behavior on top of v20.0 controls.
4. `v20.3.0` — Adaptive Institutional & Two-Sided Thesis Intelligence.
5. `v20.3.1` — AODR Adaptive Opportunity Intelligence.
6. `v20.4.0` — Agent Orchestration & Controlled MCP/API — remaining #165.
7. `v20.5.0` — Adaptive Operations — includes bounded provider usefulness/cost/reliability priors inside Router policy; no parallel router or automatic lifecycle/rights promotion.
8. `v20.6.0` — Professional Adaptive Closure; no feature scope.

## Permanent intelligence direction

`canonical evidence -> deterministic intelligence -> cross-feature synthesis -> point-in-time outcomes -> bounded adaptive learning -> optional AI/agent explanation/orchestration`

#170 is a mandatory cross-integration/Market-Regime gate for every intelligence-bearing version. #171 is a mandatory UI/data-density/intelligence-maturity audit for every changed user-visible surface. Neither is a standalone isolated feature.

Hard UI protections remain: keep Day/Swing/Long Desk look-and-feel/workflow; keep Dashboard Market Regime and Desk Control materially intact; preserve Data Engine look-and-feel; preserve the current AI Copilot engine/header visual treatment unless separately justified.

## Current exact state / next action

The active release remains `v19.0.0` on #148 / PR #149. The Adaptive Provider Registry / Market Data rebaseline is approved **future** scope and does not change the current dependency band. Continue from live GitHub/executable evidence and close the current v19.0 exact next action before starting `v19.1.0`. Governance-only provider-rebaseline commits require fresh exact-head CI before any dependency-band advancement or release qualification.
