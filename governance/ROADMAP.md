# DE.PULSE — Canonical Adaptive Roadmap

**Status:** APPROVED / ADAPTIVE  
**Rule:** Mission and approved workstreams are durable; exact minor placement may adapt when dependency, safety, performance, evidence, or data-rights constraints justify movement.

---

## v18 — Secure Multi-User Platform + Smart Provider Intelligence

Approved v18 workstreams include:

### v18.0 — Identity & Secure Sessions
Identity/authentication/session foundation, Argon2id credentials, secure server-side sessions, role enforcement, migration compatibility, and Stable-safe native delivery.

### Smart Intelligent Provider Router v2
Capability/entitlement-aware provider routing, Preferred vs Serving semantics, deterministic cooldown/circuits, provider budgets, Tier protection, source disagreement handling, calls avoided, and provider usefulness telemetry.

### v18.1 — Multi-User / My Market Symbols
Per-user symbols/watchlists/preferences with shared Global Symbol Registry and Shared Symbol Intelligence. Avoid duplicate canonical provider/calculation work for the same symbol.

### Adaptive Opportunity Discovery & Recommendations Foundation
Begin dependency-compatible productization of the market intelligence DE.PULSE already computes so users can see **My Market Opportunities** and **Global Opportunities** without creating another scanner or recommendation data silo.

v18/v19 foundation should reuse Global Symbol Registry, Reliable Actionable Universe, Shared Symbol Intelligence and Opportunity Radar to support:
- user-specific My Market vs outside-My-Market bucketing;
- shared canonical opportunity inputs rather than duplicate fetch/calculation work;
- bounded candidate ranking and NOW/WATCH/PASS/ABSTAIN semantics;
- Long/Short/horizon-aware opportunity context where current canonical TDTI/ASBI evidence supports it;
- point-in-time recommendation/ranking lineage and outcome collection;
- simple user-facing `why now / what confirms / what changes the view` synthesis grounded in canonical evidence;
- ADR-GDI-aware suppression/demotion when required evidence is degraded.

Do not disrupt the frozen v18.2 Admin / Presence / Session Operations scope to force mature AODR intelligence into that release. Implement foundational pieces only when dependency-compatible with the active build plan.

### TradeInsight — SHADOW / SECONDARY Intelligence
Use through the canonical Smart Router only. Approved contextual roles include insider/congressional intelligence, historical OHLCV backfill/reconciliation, corporate actions, symbol metadata enrichment, Opportunity Radar corroboration, and future controlled AI/MCP research where rights permit.

### Institutional Holdings / 13F Evidence Foundation
Begin/continue canonical point-in-time SEC 13F evidence collection where useful and lawful. Preserve filing/report-period identity, manager/security mappings, amendments/restatements, normalized disclosed holdings, quarter-over-quarter deltas, provenance/freshness, limitations, and subsequent outcomes. Direct SEC EDGAR remains canonical filing truth; provider enrichment must flow through the canonical Provider Router.

Do not treat 13F as live ownership or force mature adaptive manager intelligence into v18.2–v18.5.

### Two-Sided Thesis Evidence Foundation
Preserve and improve the point-in-time evidence/outcome lineage required for future Long-vs-Short thesis validation, including existing Entry/Target/Invalidation plan truth, previously approved short-entry context where supportable, first-event ordering, MFE/MAE, ASBI-relevant behavior evidence, catalyst/regime/liquidity context, and reliable short-specific data provenance when available.

Do **not** force the mature adaptive Two-Sided Directional Thesis engine into v18.2–v18.5 as scope creep. v18 should ensure current architecture and evidence are not designed in a way that blocks truthful two-sided learning later.

### Adaptive Data Reliability & Graceful Degradation — v18 Foundation
Implement the dependency-compatible foundation of the approved 10/10 ADR-GDI contract rather than waiting for v20 or assuming PostgreSQL alone will solve degradation.

v18 priorities include:
- capability-level health rather than one broad provider/app flag;
- canonical degradation reason codes;
- consumer/dependency-aware blast radius;
- dataset/horizon/session freshness SLOs;
- provider circuit/retry discipline;
- duplicate-work elimination, single-flight/coalescing and calls avoided;
- workload priority, bounded queues/backpressure and graceful load shedding;
- warm canonical persistence/restart recovery;
- truthful `UNKNOWN`, degraded and ABSTAIN semantics;
- concise impact-aware user messaging plus deeper Maintenance diagnostics.

The core target is to eliminate **self-inflicted, overly broad or unexplained** `DATA DEGRADED` states while preserving truthful warnings when required evidence is genuinely unavailable/stale/unreliable.

### v18.2 — Admin / Presence / Session Operations
SUPER OWNER / OWNER / ADMIN operations, user/session visibility and lifecycle, bounded presence, enable/disable and revocation controls, role-aware operations, and approved session-policy controls.

### v18.3 — PostgreSQL / Hosted Shared State
PostgreSQL repository parity, migrations, transactions/concurrency, shared canonical runtime state, browser/hosted architecture, pooling, backup/restore, migration/export, load/contention testing, and hosted health/readiness.

v18.3 must also use the persistence transition to harden ADR-GDI foundations where dependency-compatible: indexed warm canonical state, restart recovery, shared-symbol reuse, DB/query/pool observability, capability-health persistence where appropriate, bounded persistence pressure, and protection against PostgreSQL becoming a new bottleneck.

AODR shared ranking/recommendation state should reuse the same canonical shared-symbol and per-user isolation architecture where dependency-compatible; PostgreSQL must not cause each user to trigger duplicate market-wide recommendation computation.

### v18.4 — Security + Commercial / Data-Rights Hardening
Secrets/security hardening, auth/session/CSRF/cookie review, adversarial authorization testing, provider entitlement/data-rights metadata, hosted/commercial-use readiness, quota/abuse safeguards, observability, and licensing/redistribution/AI-use suitability.

### v18.5 — Major Closure & Release Assurance
**MANDATORY before v19.** Reconstruct full v18 scope and run fresh architecture, source quality, performance/capacity, security, data-utility, UI/UX, adaptive-intelligence, native/runtime, Principal Engineer, Professional Trader/Investor, and release-assurance closure.

ADR-GDI is a mandatory v18.5 closure dimension. Prove under realistic supported load that local/runtime architecture does not materially cause broad `DATA DEGRADED`; test provider failures/rate limits, stale data, source disagreement, DB pressure/unavailability, queue saturation, restart/warm-start, multi-user/symbol fan-out, background-job pressure, load shedding, recovery hysteresis, blast-radius correctness and actual packaged-runtime degradation UX.

If dependency-compatible AODR foundation is present by v18.5, closure must verify it reuses canonical symbol intelligence, preserves My Market vs Global separation, does not manufacture recommendations under weak/degraded evidence, and does not create unbounded full-universe/provider load.

If self-inflicted overload can delay or misstate decision-critical live/current evidence, it is a release blocker until fixed or explicitly constrained with truthful operating limits.

#### v18.5.1 — Audit, Containment & Urgent Recovery

**CURRENT RELEASE ENTRY SLICE.** v18.5.1 reconstructs the full v17/v18 ledger, installs anti-slicing enforcement and executes the safest urgent recovery work. It is not automatically the final v18 closure release. Open applicable work must be explicitly assigned to v18.6 or a later evidence-selected v18.x slice; it may not remain unowned or disappear.

Mandatory v18.x recovery-program set:

- **COPY-18.5.1-001 — version/profile preservation truth:** remove stale `v17` preservation copy and derive migration/preservation messaging from the actual source and target release context.
- **SYMBOL-18.5.1-001 — removal contract:** the visible row remove control must work consistently in Day, Swing and Long, including the final desk membership; counts, selection, persistence and reload behavior must reconcile.
- **SYMBOL-18.5.1-002 — active desk state:** a symbol's membership in the current desk must be unmistakably distinguished from its memberships in other desks across Day, Swing and Long.
- **NAV-18.5.1-001 — interaction continuity:** live/SSE refresh must not unexpectedly move a user from the middle of a section to the top or discard relevant focus/selection state.
- **RESEARCH-v15.1.0-17-19-REOPENED — Research top area:** reopen the approved Research ticker-input consistency, freshness-badge placement and responsive-placement requirements. Correct layout density/alignment, truthful freshness and action semantics, disabled-state clarity, and evidence-incomplete recovery guidance.
- **IMPL-18-TRADEINSIGHT-001 — orphaned committed workstream:** implement and certify the approved TradeInsight SHADOW / SECONDARY roles through the canonical Smart Router. Current source contains no adapter, endpoint/configuration, router capability, provider-rights record or integration test; this is a confirmed v18 implementation miss.
- **IMPL-18-UTILITY-001 — shared snapshot acquisition:** replace independent Scanner/Radar Alpaca snapshot fetch paths with one bounded, freshness-aware, in-flight-coalesced canonical broker while preserving distinct ranking.
- **IMPL-18-UTILITY-002 — Session Intelligence Coordinator:** place Pre-Market and Market Open checkpoints under one canonical temporal coordinator with shared scheduler/router/cache ownership.
- **IMPL-18-UTILITY-003 — Market Activity surface retirement:** retain seeds as an input or drill-down, not a prominent equal-level Discovery card.
- **IMPL-18-UTILITY-004 — legacy evidence routes:** redirect standalone News/Earnings/Filings routes to Research or Market/Event Intelligence with regression-safe deep links.
- **IMPL-18-DOC-001 — role-aware documentation:** enforce the frozen OWNER/ADMIN/USER/DEMO documentation composition and create the required Documentation Impact Manifest.
- **IMPL-17-DEPS-001 — dependency readiness:** implement the External Dependency & Provider Readiness checkpoint and durable User Action Required register over existing capability/rights foundations.
- **VERSION-18.5.1-002 — active version drift:** remove or explicitly classify/test obsolete v17, v18.0.4 and v18.4.0 identity strings in current user/runtime/TEST-profile paths.
- **UTILITY-v18.3-CARRY-FORWARD — 13 orphaned remediations:** disposition and close every item in `functionality_utility_remediation.json`; six have confirmed current-source failures and seven require fresh behavioral/design/package proof.

v18.5.1 must run a **full v17 + v18 approved-implementation reconciliation**. Reconstruct the 48 inherited requirements that v17/v18 had to preserve, all 20 frozen v17 items, every v18 major workstream, every v18.0–v18.5 release clause, all 13 orphaned functionality/utility remediations, accepted conversational commitment, defect history and issue record. Classify current source plus actual packaged behavior as `FRESH_PASS`, `REOPENED`, `NOT_IMPLEMENTED`, `INTENTIONALLY_SUPERSEDED`, `NOT_APPLICABLE` or explicitly roadmap-placed future scope. Nothing is silently assumed complete because an earlier release once recorded PASS.

**Adaptive v18.x Recovery & Closure Program:** the program executes inventory freeze → urgent defects → confirmed implementation misses → all 13 utility remediations → complete v17/v18 evidence reconciliation → cross-cutting hardening → zero-gap major-closure G10 → one immutable RC → actual macOS/Windows audit → G16 prevention/handoff. These phases may span v18.5.1, v18.6, v18.7 or later coherent v18.x slices.

The reconciliation ledger is the parent authority. Release slices must conserve the same immutable IDs; they cannot independently drop or close work. Slice G1/G3/G10/G15 blocks on missing ownership or unexplained placement. The final major-closure G10/G15 additionally blocks on any open/reopened/not-implemented applicable row.

After each slice, use `Observe → Reconcile → Prioritize → Slice → Build → Validate → Measure → Learn → Replan` to select the smallest safe next v18.x scope. Do not predeclare a final version or compress unresolved work to fit one.

v19 planning may continue, but implementation cannot dilute or share the active v18 source/evidence lane before the evidence-selected v18 closure reaches G16.

Detailed authority: `governance/V18_ADAPTIVE_RECOVERY_AND_CLOSURE_PROGRAM.md`.

Permanent anti-miss rules:

1. Before adding a report, search and compare it with prior requirements/defects. A reproducible recurrence reopens and escalates the original item; it is not discarded as a duplicate.
2. Every defect and implementation promise must preserve one durable chain: **origin → current observation → owner → fix/disposition → regression test → actual package proof → closure evidence**.
3. A source, dependency, renderer or package change invalidates earlier PASS evidence for the affected surface and its dependents.
4. Closure requires direct browser evidence plus actual macOS Apple Silicon and Windows x64 packaged-runtime evidence for affected user workflows.
5. A v18.x slice is blocked by unresolved items assigned to that slice or any unexplained/unowned remainder. The final v18 closure is blocked until all applicable defects, misses and reconciliation rows have current final evidence.


#### v18.6+ — Adaptive Recovery, Implementation & Hardening Capacity

No fixed feature allocation is assumed. At each prior-slice G16 / next-slice G0–G3, choose from the still-open ledger using user impact, recurrence, dependency readiness, rights/security, performance/freshness risk, coupling, evidence invalidation and learning value.

Provisional evidence-selected sequence:

- **v18.5.1 candidate:** control plane plus repeated user-trust defects and exact next-slice placement.
- **v18.6 candidate:** canonical utility/ownership recovery, route/surface cleanup, role-aware documentation and dependency readiness.
- **v18.7 candidate:** TradeInsight/provider qualification plus remaining intelligence-surface consolidation.
- **v18.8+ candidate:** full evidence convergence, cross-cutting hardening and major closure only if zero-gap ready.

This sequence is not frozen product scope. Each minor release becomes binding only at its own G1 and may be split or reordered from measured evidence. The final v18 closure release number is assigned only when the remaining applicable scope is ready for one immutable RC and 2/2 native artifact audit.

---

## v19 — Professional Data Infrastructure

Purpose: make provider/data quality, rights, cost, reliability, and suitability measurable rather than assumption-driven.

Approved direction:
- provider capability matrix;
- entitlement and commercial suitability;
- quality/reliability/latency/rate/cost/coverage;
- redistribution/persistent-storage/AI-use rights;
- data reconciliation and source disagreement;
- historical depth/adjustment quality;
- specialized/paid providers only when measured capability gaps justify them;
- formal long-term role classification for SHADOW candidate providers such as TradeInsight;
- harden institutional/13F ingestion, manager identity, security/CUSIP/FIGI mapping, combination/notice reports, amendments/restatements, corporate-action reconciliation, point-in-time semantics, filing-lag/freshness truth, storage/indexing, outcome lineage, and data-rights/provenance;
- harden two-sided thesis evidence infrastructure: point-in-time Long/Short plan snapshots, target/invalidation ordering, side-aware MFE/MAE, behavior/regime/catalyst lineage, short-interest/crowding/shortability/borrow/SSR data only where trustworthy and lawful, and explicit UNKNOWN semantics when unavailable;
- harden AODR opportunity infrastructure: point-in-time ranking/recommendation lineage, My Market vs Global bucket truth, user-preference isolation, shared-ranking efficiency, diversity/correlation metadata, NOW/WATCH/PASS/ABSTAIN transitions, recommendation usefulness outcomes, missed-opportunity analysis, staleness/degradation effects and explainable ranking provenance;
- harden ADR-GDI with measured capability SLOs, degradation history, fallback quality, provider/DB/runtime reliability scorecards, query/index tuning, capacity limits, restart behavior, load-shedding effectiveness and commercial/hosted operating limits.

v19 must also ensure sufficient point-in-time historical evidence, feature history, outcome history, provenance, rights and reliability lineage for v20 adaptive research, including ASBI, adaptive Institutional Holdings / 13F Intelligence, 10/10 Two-Sided Directional Thesis & Trade Plan Intelligence, Adaptive Opportunity Discovery & Recommendations, and adaptive reliability optimization.

**Mandatory v19 Major Closure before v20.**

---

## v20 — Adaptive Intelligence & Decision Research

Purpose: use structured evidence/outcome history accumulated since v17 to improve decision support without creating a silent self-modifying trading system.

Approved capabilities include:
- historical analogues;
- pattern/state outcomes;
- regime-conditioned evidence;
- calibration;
- false-positive / false-negative analysis;
- contradiction and drift analysis;
- provider/evidence usefulness;
- model/prompt governance;
- explainable adaptive ranking;
- controlled Champion/Challenger evaluation;
- **Adaptive Stock Behavior Intelligence (ASBI) — 10/10 Contract**;
- **Adaptive Institutional Holdings / 13F Intelligence**;
- **Two-Sided Directional Thesis & Trade Plan Intelligence (TDTI) — 10/10 Contract**;
- **Adaptive Opportunity Discovery & Recommendations (AODR) — 10/10 Contract**;
- **Adaptive Data Reliability & Graceful Degradation Intelligence (ADR-GDI) — governed adaptive optimization using reliability history accumulated earlier**.

Production promotion remains:

**SHADOW → VALIDATED → APPROVED → PRODUCTION**

No execution and no silent self-modification.

---

# Adaptive Stock Behavior Intelligence — Roadmap Placement

ASBI is a core v20 intelligence workstream, but its trustworthy data foundation begins earlier.

### v18 / v19 Preparation
Collect and preserve, where useful and lawful:
- point-in-time structured evidence;
- behavior features and sequence context;
- subsequent outcomes;
- symbol eligibility/reliability metadata;
- provider provenance/freshness;
- market/sector/regime context;
- catalyst/news/SEC/earnings context;
- institutional/13F context where available and point-in-time valid;
- historical depth/adjustments;
- decision lineage;
- data-rights metadata.

Do not force the mature ASBI prediction engine into v18.2–v18.5 as scope creep.

### v20 Major Implementation
Build/validate:
- Behavioral Fingerprints;
- state-transition modeling;
- competing scenario/path generation;
- multi-horizon Behavior Outlooks;
- Behavior Probability Momentum;
- hierarchical symbol/peer/sector/regime learning;
- catalyst-aware analogue retrieval;
- expected-move/outcome distributions;
- Behavior Probability vs Opportunity Quality;
- confidence/data sufficiency;
- evidence conflict/independence;
- ABSTAIN / NO RELIABLE EDGE;
- Early Warning vs Confirmation;
- immutable Behavior Intelligence Ledger;
- calibration/drift/outcome measurement;
- Champion/Challenger SHADOW evaluation.

---

# Institutional Holdings / 13F Intelligence — Roadmap Placement

13F is a durable Smart-Money / Institutional Intelligence dataset and adaptive context layer, not a live-trading signal.

### v18 / v19 Preparation
Collect/reconcile applicable public SEC evidence with point-in-time truth:
- 13F-HR / 13F-HR/A / 13F-NT / 13F-NT/A;
- manager CIK / Form 13F identity / accession;
- report period and filing/acceptance timestamp;
- disclosed holdings and relevant information-table fields;
- amendments/restatements and later-added holdings;
- CUSIP/FIGI/security mapping into the Global Symbol Registry;
- quarter-over-quarter states such as NEW / INCREASED / REDUCED / REPORTED EXIT / UNCHANGED / NOT COMPARABLE;
- split/corporate-action/reorganization reconciliation;
- explicit filing-lag, short-position, small-position, confidential-treatment, and incomplete-portfolio limitations;
- subsequent outcome history from the time information became public.

### v20 Adaptive Institutional Intelligence
Build/validate:
- Manager / Institutional Behavioral Fingerprints;
- persistence and disclosed-concentration behavior;
- accumulation/reduction breadth;
- consensus vs crowding;
- passive/index/common-factor adjustment where feasible;
- sector/thematic institutional rotation;
- manager/cohort usefulness by stock type, sector and regime;
- correlation with insider, congressional, ASBI, Rapid Move, Opportunity Radar, earnings/SEC/news, market/sector regime, price/volume/relative strength and options context where useful;
- convergence and contradiction reasoning;
- adaptive stale-data penalties;
- historical outcome distributions and calibration;
- Champion/Challenger evaluation for learned institutional models/rankings;
- ABSTAIN when filing completeness, mapping, history, independence, or relevance is insufficient.

13F-derived production influence remains **SHADOW → VALIDATED → APPROVED → PRODUCTION** and must never silently alter protected deterministic Day/Swing/Long formulas.

---

# Two-Sided Directional Thesis & Trade Plan Intelligence — Roadmap Placement

The earlier product direction already allowed `entry / short-entry zone context when supportable`. The 10/10 TDTI contract turns that into a complete AI/LLM-style competing-thesis system without creating a separate short-trading product or execution engine.

### v18 / v19 Preparation
Preserve/build the evidence and validation substrate needed to compare Long and Short fairly:
- one canonical ticker/horizon evidence snapshot;
- existing Entry/Target/Invalidation truth and short-entry context where supportable;
- structural trigger/confirmation/invalidation evidence;
- side-aware outcome ordering;
- actual eligible-entry anchoring for MFE/MAE where measurable;
- ASBI state/path context;
- market/sector regime, catalyst, liquidity and relative-strength lineage;
- institutional/insider/congressional context with freshness/independence;
- options context where useful;
- short-interest/crowding and borrow/shortability/SSR context only when reliable, lawful and clearly sourced;
- explicit UNKNOWN/ABSTAIN when short-specific evidence is unavailable.

### v20 Major Two-Sided Intelligence
Build/validate:
- competing Long / Short / No Reliable Edge thesis reasoning from the same evidence snapshot;
- separate Direction Probability, Thesis Strength, Confidence, Opportunity Quality and Readiness;
- Long Entry / Trim-Target / Invalidation / R:R;
- Short Entry / Cover-Trim / Downside Targets / Short Invalidation / R:R;
- structural rather than naïvely mirrored invalidation;
- per-side readiness lifecycle and probability momentum;
- ASBI-conditioned short-chase, squeeze, rebound, bull-trap and continuation reasoning;
- multiple expected paths, outcome distributions and time-to-resolution;
- long-specific and short-specific risk intelligence;
- cross-horizon contradiction/reconciliation;
- evidence independence and cause-aware reasoning;
- AI/LLM-style concise `WHY / CONFIRMS / INVALIDATES / WHAT CHANGES THE VIEW / WHAT TO WATCH NEXT` synthesis;
- immutable two-sided thesis/trade-plan ledger;
- side-aware historical calibration, MFE/MAE, false positives/misses and Decision Utility;
- Champion/Challenger evaluation;
- Professional Trader/Investor acceptance;
- ABSTAIN / NO RELIABLE EDGE as a first-class valid outcome.

TDTI must reuse ASBI, Day/Swing/Long, Opportunity Radar, Decision Queue, Research, Historical Validation and canonical evidence ownership. It is **not** a new duplicate intelligence silo.

Production influence remains **SHADOW → VALIDATED → APPROVED → PRODUCTION**. No silent formula/model self-modification and no execution capability.

---

# Adaptive Opportunity Discovery & Recommendations — Roadmap Placement

AODR turns existing global/shared market intelligence into useful user-facing prioritization. It is not a replacement for Opportunity Radar, the Global Symbol Registry, ASBI or TDTI.

### v18 / v19 Preparation
Build/harden dependency-compatible foundations:
- My Market Opportunities vs Global Opportunities bucketing with no per-user duplicates;
- Global Symbol Registry / Reliable Actionable Universe eligibility truth;
- Shared Symbol Intelligence reuse so ranking does not trigger duplicate provider/calculation work;
- Opportunity Radar broad-observe → PROMOTE → deeper shared analysis → DEMOTE lifecycle;
- bounded canonical ranking inputs and material-change propagation;
- Long/Short/horizon labels using current canonical TDTI/ASBI truth where available;
- NOW / WATCH / PASS / ABSTAIN states;
- point-in-time candidate/rank/reason/bucket lineage before outcomes are known;
- recommendation usefulness, miss, redundancy, extension/chase, staleness and degradation outcomes;
- correlation/sector/theme/common-catalyst metadata for diversity-aware surfacing;
- user preference/relevance state isolated from shared canonical market truth;
- grounded concise why-now / confirms / invalidates / what-to-watch-next presentation;
- ADR-GDI dependency and freshness integration.

Do not force mature adaptive ranking/personalization into frozen v18.2 scope.

### v20 Major Adaptive Opportunity Intelligence
Build/validate:
- adaptive cross-candidate ranking using TDTI Opportunity Quality/Readiness and ASBI state/path/probability momentum;
- expected-magnitude/time-to-resolution and opportunity-cost-aware prioritization;
- strong penalties for extension/chase, poor R:R, degraded required evidence, stale data, contradictions and low independence;
- diversity/correlation-aware recommendation sets rather than many copies of the same factor/catalyst;
- personalized relevance layered **after** shared market truth;
- historical candidate-vs-surfaced-vs-missed comparison;
- adaptive recommendation utility/calibration and false-positive/miss analysis;
- Champion/Challenger ranking/personalization evaluation;
- concise AI/LLM-style synthesis grounded in canonical evidence;
- ABSTAIN / no strong opportunities as a valid high-quality outcome.

Production ranking/personalization influence remains **SHADOW → VALIDATED → APPROVED → PRODUCTION**. No silent self-modification and no execution capability.

---

# Adaptive Data Reliability & Graceful Degradation Intelligence — Roadmap Placement

ADR-GDI is a cross-cutting reliability architecture and adaptive operating contract. Basic reliability is not deferred to v20.

### v18.3 Foundation
Implement/harden:
- capability-level health;
- consumer/dependency blast radius;
- dataset/horizon/session freshness SLOs;
- reason taxonomy and recovery state;
- warm SQLite/PostgreSQL canonical persistence;
- shared-symbol reuse and request coalescing;
- workload priorities, bounded queues/backpressure and load shedding;
- Provider Router circuits/cooldowns/fallback discipline;
- DB/query/pool/runtime observability;
- graceful degraded/UNKNOWN/ABSTAIN semantics;
- impact-aware UI plus Maintenance diagnostics.

Dependency-compatible fixes for retry storms, duplicate work, broad false degradation or local overload should be implemented earlier rather than waiting for v18.3 if they are safe and clearly evidenced.

### v18.5 Mandatory Reliability Closure
Prove:
- supported-load SLO attainment for decision-critical capabilities;
- self-inflicted degradation is eliminated or bounded by explicit truthful operating limits;
- optional failures do not contaminate unrelated consumers;
- required failures produce scoped degradation/ABSTAIN;
- provider/DB/runtime pressure does not create retry/fetch storms;
- restart/warm-start avoids unnecessary full rebuilds;
- load shedding protects high-value live work;
- recovery uses hysteresis and does not flap;
- packaged runtime can explain reason, impact, fallback and recovery.

### v19 Professional Reliability Hardening
Measure and optimize provider/DB/runtime reliability, SLOs, fallback quality, degradation history, query/index/capacity behavior, commercial operating limits and cost/value efficiency.

### v20 Adaptive Reliability Optimization
Use accumulated reliability history to improve provider recovery prediction, cooldown/backoff selection, workload prioritization, fallback usefulness and capacity policy through governed SHADOW/Champion-Challenger evaluation.

No adaptive reliability policy may silently self-promote to production.

Full permanent contract: `governance/ADAPTIVE-DATA-RELIABILITY-CONTRACT.md`.

---

# Permanent Adaptive Placement Rule

Before assigning any approved workstream to a specific minor release, evaluate:

**value + dependency + architecture + safety + performance + defects + data/provider readiness + rights + test evidence + accumulated outcomes**

If moving a workstream is safer or more correct, update this roadmap and record the decision in `governance/DECISION-LOG.md` rather than silently changing scope.