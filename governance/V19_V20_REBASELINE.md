# DE.PULSE v19/v20 Zero-Miss Rebaseline

**Status:** AUTHORITATIVE FUTURE REBASELINE  
**Stable baseline:** `v18.10.0` — immutable  
**Audit baseline main:** `7c8d0c6614ff4e8c14fc1fabb6aeadcf28a9e92c`  
**Active implementation:** issue #148 / PR #149 / `adapt-hosted-trust-foundation-001`  
**Testing/future backlog reconciled:** #150–#171 / 22 of 22  
**Hosted requirement ledger reconciled:** HOST-001–HOST-072 / 72 of 72

This document supersedes the prior future scheduling model that assigned a public version to almost every requirement row. Requirements remain conserved and individually evidenced, but they are **not releases**. The current planning and delivery unit is a coherent product **version/build**.

Machine companions:
- `governance/programs/V19-V20-REBASELINE/backlog-version-matrix.json`
- `governance/programs/V19-V20-REBASELINE/host-requirement-version-map.json`

## 1. Audit conclusion

The v18.10.0 closure remains the certified product baseline and is not reopened. Its exhaustive feature-assurance inventory, tests, canonical-owner map, UI/UX assurance and packaged evidence are reused as the starting truth. The rebaseline therefore audits **current source delta + active v19 work + all 22 new backlog issues + all 72 hosted requirements**, rather than pretending the certified v18 product is unknown.

The previous future roadmap was too granular: it used dozens of requirement-sized version reservations. That increased planning/CI/release overhead without improving requirement conservation. The new sequence keeps every requirement but groups implementation by real architecture/product boundaries.

### Current active v19 reality

The existing #148 / PR #149 work is retained and becomes **v19.0.0 Hosted Trust & Identity Foundation**. It is not restarted.

At the audit point:
- HOST-001..003 provider-rights implementation exists but final-version verification is still pending;
- HOST-004..007 identity/device/session/reauth work is in progress;
- PR #149 head `c5d0713d16f95522fd013123a78bc7cc58dc2422` is **not qualified**;
- Fast run #1141 failed source-health because five newly added hosted identity/session helpers have no production reference yet;
- the correct next product action is to wire those helpers through the existing authenticated identity/HTTP owners, then obtain exact-head Fast evidence.

No governance rebaseline may describe this unfinished head as completed.

## 2. Permanent product boundaries

All versions preserve:
- U.S. Equities Processing only, with GLD/SLV/USO as governed actionable exceptions;
- No Execution;
- Smart Provider Router v2 as sole general provider routing/admission authority;
- direct SEC/EDGAR authoritative where already governed, including Form 4 truth;
- canonical Data Health/freshness/degradation, subscription, cache/persistence/state, reconciliation, telemetry, identity/session/calendar owners;
- one canonical ticker/capability computation reused across consumers where rights permit;
- Mac Apple Silicon + Windows x64 required; shared hosted capability also requires Web parity unless explicitly justified N/A;
- provider lifecycle `SHADOW -> VALIDATED -> APPROVED -> PRODUCTION`; no automatic promotion;
- missing/stale/withheld evidence is never fabricated into healthy/neutral truth;
- deterministic market truth remains available when AI/adaptive systems are disabled.

## 3. New permanent intelligence rule

DE.PULSE must become smarter by **connecting and learning from canonical evidence**, not by adding more isolated cards or chat boxes.

Target maturity:

`canonical evidence -> deterministic intelligence -> cross-feature synthesis -> point-in-time outcome accumulation -> bounded adaptive learning -> optional AI/agent explanation or orchestration`

For every new or materially changed data point, feature, engine or surface, G2/G3/G10 must explicitly resolve:
1. canonical owner;
2. upstream canonical evidence;
3. user/trader decision purpose;
4. derived intelligence produced;
5. downstream consumers;
6. duplication/isolation risk;
7. Market Regime contribution: `YES / CONDITIONAL / NO`;
8. Outcome Learning contribution: `YES / CONDITIONAL / NO`;
9. intelligence maturity: `DETERMINISTIC_ONLY / ADAPTIVE_CANDIDATE / LEARNING_ENABLED / AI_ASSISTED / NOT_USEFUL`;
10. point-in-time/no-lookahead reproducibility;
11. freshness/degradation/recovery re-evaluation behavior;
12. role/product-entitlement/provider-right implications;
13. UI disposition where visible: `KEEP / IMPROVE / CONSOLIDATE / COLLAPSE / MOVE / REMOVE`;
14. positive, negative, failure, persistence, performance and platform evidence as applicable.

No new top-level release gate is introduced. These are mandatory dimensions inside existing G0–G16 assurance.

### Cross-integration contract

Each intelligence capability must explicitly classify integration with:
- Market Regime / Market Intelligence;
- Tradeability / Readiness;
- Discovery / Opportunity Radar;
- Watchlist / Global Symbols;
- Research;
- Day / Swing / Long Desks;
- Pre-Market / Market Open Prep;
- alerts / catalyst / rapid-move workflows;
- Pattern / Outcome Learning;
- Data Health / Maintenance;
- adaptive processing priority.

Use `REQUIRED / CONDITIONAL / NOT_USEFUL`. No integration is left accidental.

Specialized engines produce normalized evidence; they do **not** become competing Market Regime owners. One symbol, Fib level, options wall, pattern or catalyst cannot directly flip the global regime. Symbol evidence may promote to sector/global evidence only through explicit breadth/materiality/freshness aggregation.

## 4. Information architecture rule

The app currently contains substantial useful information. The goal is not indiscriminate removal; it is **decision density**.

For every visible section/card/metric/control:
- KEEP if it materially helps the page's primary workflow;
- IMPROVE if useful but poorly expressed;
- CONSOLIDATE if the same conclusion is repeated;
- COLLAPSE if useful only as deep evidence;
- MOVE if another canonical surface owns the workflow;
- REMOVE if low-value, obsolete, debug-like or misleading.

A UI removal does not delete useful canonical evidence. The evidence may continue feeding synthesis, Market Regime, learning or diagnostics.

Hard preserved UI decisions:
- Day/Swing/Long Desk current visual design and workflow remain materially intact;
- Dashboard Market Regime remains materially intact;
- Dashboard Desk Control remains materially intact;
- Data Engine current look-and-feel is preserved except genuine correctness/containment/accessibility defects;
- current AI Copilot engine/header visual treatment is preserved unless a separately proven defect or later product decision changes its role.

## 5. v19 rebaseline — deterministic hosted product + intelligent evidence foundation

### v19.0.0 — Hosted Trust & Identity Foundation

**Current active version.** Existing #148 / PR #149 remains the implementation vehicle.

Scope:
- HOST-001..023 in full;
- provider commercial/multi-user/proxy/cache/redistribution/display/AI rights evidence and fail-closed enforcement;
- tenant/account identity, role/capability truth, device lifecycle, session/refresh/revocation/reauth;
- product entitlement and quota policy distinct from RBAC/provider rights;
- privacy/data lifecycle;
- environment/IaC/service trust;
- PostgreSQL tenancy/HA/PITR foundation;
- managed secrets/KMS + supply-chain provenance;
- provider scorecards;
- point-in-time source/observed/revision/no-lookahead primitives;
- backlog #164 core authentication/session lifecycle audit;
- backend/security portions of #156 where required to prove capability authority.

Exit: HOST-001..023 all evidenced or explicitly externally blocked; exact-head Fast + impact-selected Qualified on the final candidate; HOST-024 remains blocked until exit.

### v19.1.0 — Canonical Data Runtime & Global Symbol Processing

Scope:
- #150 canonical data-path traceability;
- #151 global ticker/capability processing + reserved SPY/QQQ live priority;
- #153 Router v2 adoption/cross-integration across all providers/capabilities;
- #154 capability recovery -> canonical refresh -> dependent re-evaluation;
- #155 canonical Maintenance truth/diagnostic semantics;
- #160 Data Engine integration truth while preserving its UI;
- #167 Global Symbol Store core, user-symbol membership separation, lifecycle/retention and adaptive processing priority.

Architectural result:
`one global instrument + one capability processing owner -> canonical state -> many users/pages/Desks`.

This version must close duplicate page/Desk/user fetch paths before hosted multi-user fan-out is built.

### v19.2.0 — Hosted Gateway, Shared Serving & Sync Core

Scope: HOST-024..039.

Includes authenticated/versioned Hosted Provider Gateway, serving authorization composition, rights/entitlement-safe shared cache, lawful live fan-out through the existing subscription manager, stream revocation, API lifecycle, durable outbox/idempotency/revision/checkpoint/bootstrap/conflict/tombstone semantics, local account isolation, protected-session sync scheduling and tenant-aware observability.

Must reuse v19.1 Global Symbol/runtime processing; it may not create a per-user market-data pipeline.

### v19.3.0 — Cross-Platform Product, Roles & Information Architecture

Scope:
- HOST-040..047 and HOST-053;
- #152 app-shell/header cleanup;
- #156 full Owner/Admin/User page/action composition + canonical company/instrument name display;
- #159 Maintenance layered information architecture;
- #155 final visual/IA reconciliation after truth is fixed;
- #160 cross-platform presentation verification, preserving Data Engine UI;
- #167 Owner `Global Symbols` Administration experience;
- #171 whole-product tab-by-tab UI/data/intelligence audit for Dashboard, Desks, Settings, Administration, Documentation and shared navigation;
- #164 cross-platform login/session/reauth/logout UX verification.

Cross-platform account/session/settings/preferences/watchlists/desks/research state must converge without role or authorization forks.

### v19.4.0 — Market Intelligence & Research Workflow Quality

Scope:
- HOST-049 Market State/Modes/Readiness parity;
- #158 Market Intelligence usefulness/density/information architecture;
- #161 Research freshness defect + whole Research evidence hierarchy;
- #162 Research Target -> deterministic Decision Brief -> optional evidence-bound AI second opinion;
- relevant #171 whole-product audit rows.

Research supplies sourced evidence; Desks remain horizon decision workflows. AI remains optional and cannot independently fetch market data.

### v19.4.1 — Discovery & Opportunity Radar Effectiveness

Separate patch because Discovery is functionally heavy.

Scope:
- HOST-048 cross-platform Discovery parity;
- #163 complete universe -> detection -> admission -> rank -> actionability -> Research/Desk handoff audit/correction;
- premarket ACTIVE, regular ACTIVE-HIGH, postmarket ACTIVE, overnight REDUCED/EVENT-WAKE behavior;
- extreme material-mover admission, false-positive/bad-tick/split/liquidity handling, provider recovery and explainable omission reasons;
- relevant #171 audit rows.

### v19.5.0 — Price/Volume & Event-Anchored Intelligence

Scope:
- #168 Event-Anchored Levels + Reaction Intelligence;
- #169 Price & Volume Structure Intelligence.

Canonical deterministic evidence includes, where data supports it: POC/HVN/LVN/Value Area, governed swing levels, anchored VWAP, Fib from governed anchors, buyer/seller response zones, confluence, event anchors/reclaims/losses/reactions for NFP/CPI/PCE/FOMC/earnings/material catalysts.

No new page by default. Feed Watchlist, Discovery, Research, Prep, Desks, alerts and bounded Market Regime aggregation.

### v19.5.1 — Options Structure & GEX Intelligence

Separate patch because provider entitlement, schema/completeness and analytical validation are heavy.

Scope: #157.

Canonical `OPTIONS_CHAIN` and `OPTIONS_GREEKS_OI` through Router v2; provider adapters; truthful completeness/freshness; GEX structural proxy; Call Wall, Put Wall, Gamma Magnet/Pin, Gamma Flip, strike/expiry clusters and optional consumer heatmap. Withhold when evidence is insufficient. No provider-specific consumer coupling.

### v19.6.0 — Point-in-Time Evidence & Outcome-Ready Foundation

Scope:
- HOST-057..064;
- deterministic foundation portion of #165.

Includes 13F/institutional point-in-time source/revision/query, two-sided Long/Short evidence substrate, AODR candidate/rank/outcome lineage, and one canonical outcome-event/evidence snapshot contract usable by Discovery, Research, Desks, event/structure/options intelligence and future learning.

This version records outcomes; it does **not** grant learned production influence.

### v19.6.1 — Hosted Reliability, Economics & Adaptive Readiness

Scope:
- HOST-050..056 and HOST-065..071;
- #170 final v19 cross-integration/Market Regime contribution matrix;
- #171 final v19 whole-product intelligence/density reconciliation.

Includes tenant observability/fairness/security, mixed-client/recovery/load assurance, SLO/error budgets/runbooks, capacity/cost economics, provider licensing gaps, adaptive dataset readiness and explicit v20 readiness.

Exit requires every v19 capability to have canonical owner, consumer, cross-integration disposition, Market Regime disposition, Outcome Learning disposition, degradation/recovery behavior and evidence ownership.

### v19.7.0 — v19 Major Closure

**No feature scope.** Owns HOST-072.

Requires:
- all 72 HOST rows reconciled;
- every backlog requirement assigned to v19 complete or explicitly conserved to its mapped v20 version;
- zero unexplained P0/P1 implementation, security, data-truth, role, rights, persistence, UX/IA, cross-integration, performance or platform gaps;
- Mac/Windows/Web parity for shared v19 capabilities;
- exact-head G0–G16 qualification/publication.

## 6. v20 rebaseline — governed adaptive intelligence

v20 begins only after v19.7.0 and v19.6.1 readiness evidence. Adaptive production influence remains shadow-first and bounded.

### v20.0.0 — Outcome Learning & Adaptive Control Plane

Scope: #165 adaptive-learning portion.

Build immutable experiment/evidence snapshots, model/prompt/version governance, champion/challenger/shadow evaluation, historical analogue retrieval, regime-conditioned Outcome Store, calibration/confidence reliability, contradiction/FP/FN/miss/drift evidence and bounded promotion/rollback rules.

No uncontrolled self-modifying production logic.

### v20.1.0 — Adaptive Chart Pattern & Similarity Intelligence

Scope: #166.

Historical bootstrap from canonical OHLCV; multi-timeframe fingerprints; named and unlabeled recurring structures; point-in-time outcome labeling; similarity retrieval; Watchlist/Discovery/Research/Desk consumers; protected live-vs-background scheduling; `OBSERVED -> LEARNING -> SHADOW -> VALIDATED -> TRUSTED` lifecycle.

### v20.2.0 — Adaptive Market Synthesis, Regime & Discovery Learning

Use v19 point-in-time outcomes to learn bounded usefulness of:
- Market Regime contributors;
- Discovery admission/ranking/miss evidence;
- event reaction evidence;
- price/volume structural evidence;
- options structural evidence;
- alert/catalyst usefulness;
- horizon-specific readiness evidence.

Learning never lets one ticker directly control global Market Regime and never bypasses freshness/rights/authority.

### v20.3.0 — Adaptive Institutional & Two-Sided Thesis Intelligence

Institutional/13F feature extraction, revision/lag handling, two-sided thesis features, regime conditioning and outcome calibration. Preserve point-in-time/no-lookahead source truth.

### v20.3.1 — AODR Adaptive Opportunity Intelligence

Separate heavy patch: adaptive Opportunity Radar ranking/scoring, why/why-not, outcome/miss learning, stability/drift/fairness and bounded promotion. No execution semantics.

### v20.4.0 — Agent Orchestration & Controlled MCP/API

Remaining #165 orchestration/external-capability scope.

Agents may orchestrate canonical Research/Discovery/Prep capabilities and LLMs may synthesize/explain bounded evidence. MCP/API exposes versioned canonical capabilities only. Enforce auth/RBAC/product entitlement/provider rights/provenance/rate limits/audit. No raw provider credentials, unrestricted DB access, parallel market-data routing or single-vendor dependency.

### v20.5.0 — Adaptive Operations

Adaptive provider/evidence selection in SHADOW, quality/cost/freshness utility evidence, bounded budget/backpressure recommendations, governed policy promotion, rollback/drift/kill controls. Smart Provider Router v2 and provider lifecycle authority remain intact.

### v20.6.0 — Professional Adaptive Closure

**No feature scope.** Full zero-gap adaptive/AI/runtime/cross-platform closure through G0–G16. Every adaptive influence must be reproducible, calibrated, bounded, explainable, rollback-capable and traceable to point-in-time evidence.

## 7. Version dependency chain

`v19.0.0 -> v19.1.0 -> v19.2.0 -> v19.3.0 -> v19.4.0 -> v19.4.1 -> v19.5.0 -> v19.5.1 -> v19.6.0 -> v19.6.1 -> v19.7.0`

`v19.7.0 -> v20.0.0 -> v20.1.0 -> v20.2.0 -> v20.3.0 -> v20.3.1 -> v20.4.0 -> v20.5.0 -> v20.6.0`

A future G0 source-overlap audit may combine work that is already implemented or externally blocked, but it may not silently drop an acceptance row. A heavy version may be split into an additional patch version only when actual implementation/risk evidence justifies it; do not create requirement-sized versions.

## 8. Zero-miss control

Before a version reaches G10:
- every assigned backlog acceptance requirement is mapped to an executable owner/test or explicit N/A/external blocker;
- every applicable HOST row is mapped through the rebaseline overlay;
- every current/source-discovered surface and background owner affected by the version is included;
- no old page/provider/user-specific parallel owner remains unexplained;
- #170 cross-integration matrix is complete for changed intelligence;
- #171 UI/intelligence matrix is complete for changed user-visible surfaces;
- role/rights/entitlement/security negatives are present where applicable;
- provider/data failures, stale/missing/partial evidence and recovery are tested;
- persistence/restart/migration, load/resource and required-platform evidence are owned;
- AI/adaptive functionality has deterministic fallback and point-in-time provenance;
- implementation and documentation/handoff agree.

At G16, audit the **whole version**, not just the last code diff, and update GitHub handoff/current state so another assistant/account/model can resume without chat memory.

## 9. Immediate next action

Continue **v19.0.0** on existing branch/PR. Fix the current source-health failure by production-wiring the hosted identity/session helpers through the existing canonical authenticated identity/HTTP path or removing/reworking helpers proven unnecessary. Then run exact-head Fast. Do not begin v19.1.0 until v19.0.0 satisfies its trust-foundation exit criteria.
