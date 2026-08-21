# DE.PULSE — Current Adaptive Build Plan

**Immutable Stable:** `v18.9.0-stable`  
**Stable candidate:** `9ea81cddae4875ae15d3719ca028519a36c597b6`  
**Stable fingerprint:** `a8719090c341c874dbd1279cc31ad98e84075d5701c46a800bf951340780ecb9`  
**Active development branch:** none  
**Active PR:** none  
**Master corrective program:** #65 / `ADAPT-PROVIDER-INTELLIGENCE-010`  
**Immediate next patch:** `v18.9.1` / #64 / `ADAPT-RUNTIME-CRASH-001`.

## Build philosophy — completeness through small patches

DE.PULSE MUST NOT use heavy multi-domain builds in v18.9.x, v19 or v20. Each patch owns one primary responsibility and its directly necessary supporting work only. The purpose is to reduce blast radius, make review/test scope exact, catch implementation misses early and keep CI/release evidence interpretable.

Before any patch implementation:
- G0 re-fetches live GitHub and exact predecessor Stable/patch truth;
- G1 freezes one bounded scope plus explicit non-goals;
- G2 maps canonical owners and proves no parallel subsystem is needed;
- G3 freezes dependency/provider/data/model contracts and deterministic acceptance tests.

Before the next patch may begin:
- G4–G10 evidence for the current patch must be coherent;
- implementation-miss audit must be performed against G1 acceptance;
- all discovered out-of-scope findings must have a durable target patch/issue;
- open issues must be reconciled;
- handoff/checkpoints must identify exactly one next action.

Exact future patch numbers are provisional until each G1. If a packet is too large, split it rather than expanding G1.

## Ordered v18.9.x build packets

### v18.9.1 — Runtime Reliability
Scope owner #64. Crash diagnosis/fix only. No provider expansion, router redesign, Market Modes or company-name work. Exact root cause/reproduction, lifecycle regression, persisted-state/API-key continuity and packaged macOS Apple Silicon proof required.

### v18.9.2 — TradeInsight Settings/API Key
Settings integration only: existing Data Provider Settings + existing local secret owner, masked field, Save/Test/Clear, configured/connected/error/capability status, safe environment override. No provider-routing behavior change.

### v18.9.3 — Coverage-Aware Router Core
Smart Provider Router v2 evolution only. Add consumer requirement/coverage contracts, **in-memory + persisted canonical DB/state reuse before provider acquisition**, residual-gap calculation, eligible-provider ranking against remaining need, targeted acquisition, canonical merge/provenance/conflict handling, coverage re-evaluation and bounded stop criteria. Separate validation lifecycle from serving role. No UI/provider-feature expansion beyond what tests require.

### v18.9.4 — Company Identity
Canonical identity owner and presentation only. All Day/Swing/Long state headings show symbol + company name when known, e.g. `APP - AppLovin : In Entry Zone`; reuse in Research/Discovery/Add Symbol; symbol-only fallback. No TradeInsight search admission yet.

### v18.9.5 — Market Data Modes
Behavior/quality-oriented Adaptive modes and capability diagnostics only. Remove hard-coded provider-brand semantics where misleading. Surface actual provider contribution/freshness/coverage in diagnostics without creating a new Market Mode owner.

### v18.9.6 — TradeInsight Form 4
Contract-validated Corporate Insider/Form 4 enrichment only. SHADOW-first, direct SEC/EDGAR authoritative, source-family de-duplication, existing SEC/Ownership model reused, optional provider failure non-degrading.

### v18.9.7 — TradeInsight Symbol Search
Contract-validated ticker/company search only. Plug into canonical symbol validation/company identity as fallback/corroboration. U.S.-equity boundary final; GLD/SLV/USO exceptions preserved.

### v18.9.8 — TradeInsight Movers
Contract-validated mover/ranking evidence only. Opportunity Radar consumes it as SHADOW candidate evidence; existing scanner/ranker remains canonical. No undocumented REST endpoint assumptions.

### v18.9.9 — Remaining TradeInsight Capability Admission
Full useful-capability inventory/disposition only. Revalidate Congress, daily adjusted/raw OHLCV, corporate actions and bounded history under coverage-aware routing. Every useful entitlement needs a named consumer/disposition/freshness/retention/rate policy. No Python/MCP production dependency and no inferred intraday support.

### v18.9.10 — Provider Efficiency / Adaptive Telemetry
Measure and harden only: coverage completion, residual gaps, DB/cache reuse, provider calls avoided, provider usefulness, latency/errors/rate limits/backpressure, freshness failures, disagreement/corroboration, bounded fan-out, CPU/memory/goroutine stability, consumer materiality and the provider/runtime headroom required to protect pre-market/regular-market/after-hours. Feed evidence into adaptive ranking without changing deterministic trading truth opaquely.

### v18.9.11 — Session-Aware Data Readiness Maintenance
Single responsibility: implement one canonical maintenance coordinator using existing persistence/cache, Smart Provider Router v2, freshness, provider budgets, workload-priority and U.S. market calendar/session owners.

Required behavior:
- **Light overnight maintenance** only in eligible low-priority windows: bounded residual historical gaps, latest completed-session persistence, incremental SEC/disclosure/revision checks, small corporate-action/fundamental/macro/identity corrections, lightweight integrity/readiness work and bounded high-value precomputation.
- **Heavy weekend/extended-market-closed maintenance** for deeper historical backfill/reconciliation, corporate-action audits, SEC/Form4/13F/congress/earnings/fundamental/macro history, point-in-time outcome resolution, provider-history consolidation and bounded DB/index/retention work.
- **Protected Tier-0 sessions:** pre-market, regular market and after-hours. These sessions get first claim on provider quota/headroom, network, CPU, memory, DB and worker capacity.
- Maintenance uses only bounded surplus capacity after protected-session reserves.
- External-provider maintenance work suspends during protected sessions; data directly required by a current/live consumer becomes live fulfillment, not maintenance.
- Maintenance must drain/preempt/checkpoint/resume promptly when protected or market-shock/high-priority work begins.
- Missed overnight/weekend work catches up only in a later eligible window, never by dumping backlog into a protected session.
- Manual `Run Data Readiness Maintenance Now`, if exposed, obeys the same provider/runtime/session protections.
- No blind full-universe refetch and no second scheduler/calendar/router/cache/database owner.

Acceptance must prove overnight light maintenance, weekend heavy maintenance, protected-session quota/runtime reservation, preemption/checkpoint/resume, restart de-duplication, bounded catch-up and **no material impact on decision-critical freshness/latency/readiness**.

Machine contract: `adaptive-governance/PERSISTENCE_REUSE_AND_OFF_HOURS_DATA_READINESS_CONTRACT.md`.

### v18.9.12 — Professional Closure Audit
No new feature scope. Audit whole v18.9.x for misses/duplication/bypasses/orphan useful capabilities; retest #57 and #64 regressions; deterministic Day/Swing/Long equivalence; all-desk identity; DB-first reuse and residual-gap acquisition; provider failure/recovery; overnight/weekend maintenance; protected-session capacity reservation/preemption; actual macOS Apple Silicon + Windows x64 packages; Adaptive Intelligence Scorecard. Closure requires zero unexplained carry-forward and zero unowned useful capability.

## v19 build packets — Professional Data Infrastructure

**Entry:** v18 final closure PASS. v19 consumes the coverage-aware routing, persistence-first reuse, session-aware maintenance, canonical identity, adaptive Market Modes and provider telemetry created in v18.9.x; it does not recreate them.

Provisional small packets, each independently frozen at G1:
1. provider capability/entitlement/data-rights registry;
2. provider quality/reliability/latency/rate/cost/coverage/SLO scorecards including calls-avoided and maintenance value;
3. data reconciliation/source-disagreement/historical-adjustment/revision quality;
4. institutional/13F evidence and identity/mapping hardening;
5. two-sided Long/Short point-in-time evidence substrate;
6. AODR candidate/ranking/outcome lineage infrastructure;
7. ADR-GDI professional reliability/capacity/load-shedding hardening including protected-session reserve sizing, maintenance/preemption economics, DB/index/pool/capacity behavior and restart/warm-start;
8. measured specialized/paid-provider gap evaluation and migration readiness through the same router/persistence/session-priority contracts;
9. v20 point-in-time evidence/feature/outcome/provenance/research-readiness audit;
10. mandatory v19 Major Closure with no unexplained provider role/dataset/right/reliability gap.

## v20 build packets — Adaptive Intelligence & Decision Research

**Entry:** mandatory v19 Major Closure PASS. v20 learns only from point-in-time, rights-valid, provenance-bound evidence and outcomes.

Provisional small packets, split further where G0/G1 shows coupling/risk:
1. adaptive research control plane + immutable experiment/evidence ledger;
2. historical analogue + regime-conditioned outcome intelligence;
3. calibration / FP-FN / miss / contradiction / drift intelligence;
4. ASBI Behavioral Fingerprints + state-transition foundation;
5. ASBI scenarios / multi-horizon outlooks / probability momentum / calibration;
6. adaptive Institutional Holdings / 13F Intelligence;
7. TDTI competing Long / Short / No Reliable Edge thesis intelligence;
8. TDTI two-sided trade-plan/readiness/outcome validation, still No Execution;
9. AODR adaptive shared opportunity ranking;
10. AODR diversity/opportunity-cost/personalized relevance after shared truth;
11. ADR-GDI adaptive provider/recovery/workload/capacity optimization in SHADOW/Champion-Challenger, including evidence-based maintenance value and protected-session reserve optimization without self-promotion;
12. model/prompt governance + Champion/Challenger + explainable promotion/rollback;
13. v20 Professional Closure proving calibrated utility, abstention, deterministic boundaries, no silent self-modification and No Execution.

## Version dependency contract

`v18.9.x trustworthy acquisition/persistence/session readiness -> v19 measured professional evidence infrastructure -> v20 governed adaptive learning`.

v19 must not undo v18 Smart Provider Router or session-aware maintenance ownership. v20 must not bypass v19 provenance/rights/point-in-time lineage. No adaptive model may repair missing data by inventing confidence; missing/weak evidence remains UNKNOWN/ABSTAIN.

## Protected-session resource contract

Pre-market, regular market and after-hours are Tier-0 decision-support sessions. Their live/current workloads always outrank maintenance. Provider quota/headroom, network concurrency, CPU, memory, DB and worker capacity must retain explicit reserve for these sessions. Maintenance gets bounded surplus capacity and must yield/preempt before it can materially degrade live/current evidence.

Session boundaries come only from the existing canonical U.S. market calendar/session owner, including holidays/half-days/exceptional closures. Maintenance may not implement a second market calendar.

## CI/release efficiency rule

For each patch: develop the coherent code+test batch before opening the PR; one PR only; exact-head Fast once for the coherent candidate, Qualified once when Ready, then one canonical G11–G16 release when release-capable. Do not spend CI budget on avoidable duplicate runs, retry branches or certification branches. Failure classification remains mandatory before rerunning.

## Protected boundaries

Smart Provider Router v2 sole routing authority; canonical freshness/recovery sole freshness owner; existing multi-feed allocator sole subscription owner; BroadSnapshotBroker canonical reuse owner; canonical persistence/cache owners reused; canonical U.S. market calendar/session owner reused; direct SEC/EDGAR authoritative; deterministic Day/Swing/Long protected; U.S. Equities Processing; GLD/SLV/USO actionable exceptions; No Execution; G0–G16 only.

## Exactly one next action

Execute G0 for issue #64 / v18.9.1 using concrete macOS crash evidence or deterministic reproduction, then freeze its narrow G1. No v18.9.2 branch until v18.9.1 is truthfully closed.
