# DE.PULSE — Adaptive Operating Contract

**Status:** PERMANENT / GOVERNING  
**Rebaselined:** 2026-08-28 from the full-product audit  
**Scope:** Adaptive Roadmap, Build Plan, Build Process, Delivery Process, intelligence learning and product continuity

This contract defines the durable rules. Version scope may adapt through an approved Decision Log entry; a release note, handoff, `CURRENT_*` projection or chat cannot silently supersede it.

## 1. Product north star and boundaries

DE.PULSE is a market-intelligence operating system for U.S.-listed equities and approved U.S. ETFs. It continuously observes, correlates, prioritizes, explains, measures outcomes and adapts under explicit governance.

Permanent boundaries:

- research, intelligence and decision support only;
- no order execution, broker routing, paper/live trading, OMS/blotter, portfolio P&L or trade journal;
- ordinary users see intelligence, confidence, freshness and explanations—not Data Engine machinery;
- Web, macOS Apple Silicon and Windows x64 are first-class clients of one domain truth, not independent products;
- development readiness and Commercial/Public readiness remain separate as defined by `governance/PRODUCTION-READINESS-TIERS.md`;
- provider credentials and technical availability never imply display, storage, derivation, AI-use, multi-user or redistribution rights.

## 2. Four connected adaptive layers

1. **Adaptive Roadmap** — where the product is going and why.
2. **Adaptive Build Plan** — what should be built next and in what dependency order.
3. **Adaptive Build Process** — how a coherent build is designed, implemented and qualified.
4. **Adaptive Delivery Process** — how qualified source becomes a trustworthy artifact and durable next baseline.

Control loop:

`Observe -> Learn -> Prioritize -> Build -> Validate -> Deliver -> Measure -> Learn`

Product loop:

`Collect when useful -> Normalize -> Correlate -> Prioritize -> Explain -> Measure outcome -> Governed adaptation`

## 3. Canonical authority and documentation truth

The repository source-of-truth hierarchy is defined in `governance/README.md`. For these four layers, the only narrative authorities are:

- `governance/ROADMAP.md`;
- `adaptive-governance/ADAPTIVE_BUILD_PLAN.md`;
- `adaptive-governance/ADAPTIVE_BUILD_PROCESS.md`;
- `adaptive-governance/ADAPTIVE_DELIVERY_PROCESS.md`.

`adaptive-governance/CURRENT_ADAPTIVE_*.md` files are compatibility/status projections only. They may identify the active Stable/work slice/branch and point to machine state, but must not carry an independent roadmap, process or next-action narrative.

Current execution status comes from `governance/current-state.json`, the active closure ledger and `handoff/CURRENT.md`. Permanent intent comes from the canonical contracts and approved Decision Log. Historical release evidence remains immutable.

## 4. Canonical G0–G16 model

No top-level gates beyond G0–G16 may be created.

- **G0 — Exact Baseline**
- **G1 — Immutable Scope**
- **G2 — Architecture & Data Utility**
- **G3 — Design & Dependency Readiness**
- **G4 — Development Exit**
- **G5 — FAST Qualification**
- **G6 — Integration & MEDIUM Qualification**
- **G7 — Data, Security & Adaptive Intelligence**
- **G8 — Performance, Capacity & Stability**
- **G9 — Cross-Module UI/UX**
- **G10 — Pre-Freeze Qualification**
- **G11 — Immutable Release Candidate**
- **G12 — Full Certification**
- **G13 — Native Packaging & Provenance**
- **G14 — Actual Artifact Runtime Audit**
- **G15 — Release Assurance & Promotion**
- **G16 — Adaptive Retrospective & Handoff**

Governance order does not require wasteful serial execution. Independent checks may run with bounded concurrency, but evidence is reusable only when source fingerprint, test definition, inputs and environment assumptions remain equivalent.

## 5. Engineering and migration order

Before adding implementation, apply:

`REUSE -> CONSOLIDATE -> REFACTOR -> DELETE/REPLACE -> ADD`

Before moving authority, apply:

`CHARACTERIZE -> ADD NEW OWNER -> DUAL/SHADOW -> COMPARE -> MIGRATE CONSUMERS -> PROVE -> RETIRE OLD OWNER -> KEEP REGRESSION`

Consequences:

- prove an owner does not already exist before creating a service, scorer, router, store, event bus or state engine;
- preserve working capability during extraction;
- freeze golden vectors before moving renderer calculations or giant snapshot consumers;
- distinguish intentional defect correction from accidental behavior drift;
- retire duplicate authority only after parity or an explicitly approved improvement is proven.

## 6. Canonical intelligence architecture

Target flow:

`Provider observations/events -> rights + quality -> deterministic evidence/features -> SymbolIntelligenceSnapshot -> Opportunity Lifecycle -> product projections -> frozen Decision Brief -> outcomes -> governed adaptation`

Authority rules:

- one versioned `SymbolIntelligenceSnapshot` is the canonical per-symbol intelligence boundary;
- one shared Opportunity Lifecycle owns `DETECTED -> OBSERVED -> QUALIFIED -> PROMOTED -> HIGH_PRIORITY -> COOLING -> DEMOTED -> RESOLVED` transitions;
- Opportunity Radar and Rapid Move remain detectors/evidence producers, not parallel product state machines;
- Discovery is the broad-universe projection;
- Watchlist is the user-selected-universe projection and must not create another scanner, scorer, provider loop or lifecycle;
- Alerts consume material transitions/incidents and deduplicate causal chains rather than rescoring;
- Research consumes the frozen snapshot/transition that produced the attention change and owns the Decision Brief;
- Day/Swing/Long are horizon projections of shared evidence with server-owned deterministic side/geometry/policy;
- equivalent lawful data is fetched once, calculated once and consumed many times.

Every transition, promotion, demotion, confidence change, alert and learned adjustment must preserve evidence IDs, rule/model/policy version, timestamps, provider/dataset/rights lineage, contradictions and deterministic fallback.

## 7. Deterministic, learned and LLM authority

Keep three layers separate:

1. **Deterministic rules** own prices, indicators, normalization, routing, rights, eligibility, canonical state and protected lifecycle/policy invariants.
2. **Registered statistical/adaptive models** may rank or calibrate only through explicit versioned features, evaluation and governed promotion.
3. **LLM synthesis** may explain, summarize, compare and support research from rights-filtered structured evidence; it must not invent calculations, overwrite canonical state or independently fetch provider truth.

Adaptive production influence follows:

`SHADOW -> VALIDATED -> APPROVED -> PRODUCTION`

No silent self-modification or automatic authority/rights promotion is allowed. Every production model/policy must support reproducibility, audit, rollback and a deterministic no-change fallback.

## 8. Adaptive evidence and outcome contract

The durable loop is:

`feature/evidence snapshot -> decision/transition -> outcome window -> point-in-time evaluation -> challenger -> shadow comparison -> drift/subgroup checks -> approval -> bounded rollout -> rollback`

Required controls where applicable:

- feature, evidence, regime, rights, rule, model, prompt and policy versions;
- time-split/walk-forward evaluation and strict no-lookahead joins;
- sample floors, uncertainty and cohort counts;
- calibration, ranking lift, precision at attention capacity, false-promotion cost and stability—not generic accuracy alone;
- regime, liquidity, catalyst, evidence-coverage and user-personalization subgroup checks;
- unbiased/control sampling to counter adaptive selection bias;
- explicit censored outcomes for halts, missing bars, delistings or incomplete windows;
- decay, drift detection, promotion thresholds, human approval and rollback thresholds;
- explanation-grounding evaluation separate from decision-quality evaluation.

Pattern learning starts with structured multi-timeframe/event features and interpretable similarity, baseline, tree or linear challengers. Deep image/sequence models require later out-of-sample evidence; an LLM is not a learning system.

## 9. Temporal and point-in-time truth

Every material fact or event must disposition:

- `source_at`, `observed_at`, `ingested_at`;
- `effective_from`, `effective_to`, decision `as_of`;
- expiration/half-life and market session/exchange timezone;
- provider, dataset, instrument identity and rights version;
- quality, coverage, revision and supersedes identity;
- raw-versus-adjusted basis where price history is involved.

Retrieval/cache time is bookkeeping, not evidence time. Unknown source time stays unknown. Fundamentals and macro facts preserve vintages. Corrections create revisions. Options OI as-of time remains distinct from live quote/IV time. Historical replay cannot silently use later revisions.

## 10. Provider, data utility and canonical state

Smart Provider Router v2 remains the sole general capability-routing authority. Direct authorities such as SEC/EDGAR remain explicit. The Adaptive Provider Registry registers adapters/capabilities; it must not become Router v3 or another health/cache/persistence/lifecycle owner.

Every significant datum must answer:

1. What decision does it influence?
2. Which canonical owner and consumers use it?
3. Is it independent, correlated or duplicate evidence?
4. What freshness, retention and materiality apply?
5. What is its provider/API/storage cost?
6. What happens when it is absent, stale, partial, revised or contradictory?

Provider/data changes must conserve the #80/#81/#82/#83/#78/#84 Data Health program, including capability classification, canonical evidence-time freshness, lawful cache/fallback, minimally scoped `PARTIAL COVERAGE` / `DATA DEGRADED`, hysteresis, recovery, workload protection, upstream-correlation truth and revision lineage. Required but missing evidence fails closed.

## 11. Performance and scale discipline

Design for bounded CPU, memory, goroutines, queues, DB I/O, provider requests, subscription fanout, storage growth and client payloads. Protect critical current evidence before optional/background work.

Preferred default:

- Go modular monolith with explicit domain/application/provider/persistence/platform boundaries;
- memory-first live state with bounded caches;
- Postgres v2 plus transactional outbox for hosted durability and events;
- stateless API/workers where needed;
- versioned REST/event/delta contracts and thin typed clients;
- incremental/material-change propagation instead of full-snapshot fanout.

Kafka, Kubernetes and microservice extraction are not default roadmap goals. Adopt them only from measured throughput, replay, isolation, deployment or team-ownership need.

## 12. Security, privacy and client boundary

Core/cloud owns provider credentials, rights enforcement, ingestion, canonical intelligence, tenant state, outcomes, alert delivery, adaptive jobs, durable audit and AI egress policy. Clients own rendering, interaction, charts, accessibility and typed contract validation.

Desktop edge may keep an encrypted authorized last-known cache, OS notifications/deep links and device credentials in OS secure storage. It must not contain critical provider secrets or represent offline/stale state as live confidence.

Hosted/client changes require tenant isolation, least privilege, managed secrets, secure sessions/refresh, CSRF/reauth where applicable, role and product-entitlement separation, account privacy lifecycle, redacted telemetry/crash output, old-client compatibility and signed update provenance.

Privacy deletion and durable security/audit retention require an explicit policy; backup/PITR restoration must not silently resurrect deleted personal state.

## 13. UI and product coherence

**Complex inside -> intelligent synthesis -> simple outside.**

User surfaces prefer material change, priority, explanation, contradiction, confidence, freshness and a clear research path. Raw provider queues, databases, scheduler/circuit internals and Data Engine machinery belong in role-gated diagnostics.

Every visible element is classified `KEEP / IMPROVE / CONSOLIDATE / COLLAPSE / MOVE / REMOVE / ADD_FIRST_CLASS`. Preserve useful evidence even when presentation moves. Avoid tab proliferation and duplicate domain calculations.

`Long King / Short King` are undefined until an approved user purpose, horizon, evidence, exclusion, confidence, outcome and cross-integration contract exists. `Call Wall / Put Wall` require expiry-aware coverage, OI as-of, strike/cluster, quality and rights semantics; gamma multiplied by OI alone cannot establish signed dealer positioning.

## 14. Build, certification and delivery

- one coherent development branch/PR per version unless governance explicitly decides otherwise;
- Fast is exact-head development qualification and never dispatches Release;
- Qualified is risk-selected candidate qualification and required at G10;
- Release is the single G11–G16 certification/package/publication path;
- source changes after freeze invalidate affected evidence;
- actual packaged macOS Apple Silicon and Windows x64 artifacts must be tested when required;
- Web delivery must validate deployed runtime/API/auth/cache behavior, not source alone;
- macOS signing/notarization/stapling, Windows trusted signing/installer, secure update channels and rollback are required before first-class public desktop distribution;
- artifact hashes/SHA manifests are generated last, after contents are immutable;
- historical Stable artifacts and provenance are immutable.

## 15. AIPLC learning checkpoint

Run a concise Adaptive Intelligence & Product Learning Checkpoint after a meaningful build/checkpoint that changes or materially exercises product behavior, data, intelligence, routing, architecture, reliability, security, UX or delivery.

AIPLC asks:

`datum -> purpose -> owner -> consumer -> freshness/materiality -> interpretation -> decision value -> outcome -> learning`

For defects:

`symptom -> root cause -> owner -> fix -> recurrence prevention -> cross-product scan -> regression/gate -> measurable follow-up`

Score only applicable dimensions with evidence: AI/LLM smartness, data utility, canonical architecture, decision utility, reliability/freshness, performance, UX intelligence, outcome learning, testing/prevention and adaptive continuity. A below-target result remains explicit; never inflate it to claim 5/5 or 10/10.

AIPLC may re-prioritize the next build but cannot bypass G1, promote a model, authorize public/commercial use or change permanent scope by itself.

## 16. Two-sided thesis and trade-plan conservation

The approved TDTI scope remains a product-specific specialization of this contract, not a separate gate system. It reuses one canonical evidence/thesis owner and one side-aware deterministic contract:

- `LONG / SHORT / NO_SETUP-WAIT` side is separate from setup quality;
- LONG target/trim is above entry and invalidation below;
- SHORT cover/target is below entry and invalidation above;
- target/cover-before-entry cannot count as a successful outcome;
- unsupported shortability/borrow/short-interest evidence remains UNKNOWN;
- AI cannot invent levels or overwrite deterministic geometry;
- No Execution remains permanent.

Roadmap placement belongs only in `governance/ROADMAP.md`; obsolete version-specific placement in this permanent contract is superseded by DEC-2026-08-28-001.

## 17. Documentation, handoff and governance change

Update user/developer/capability/security/build/release documentation when behavior, architecture, providers, rights, persistence, platform packaging or operational truth materially changes. Do not rewrite many parallel current narratives for small patches.

Every meaningful checkpoint must leave GitHub, the active branch, machine state, closure ledger and `handoff/CURRENT.md` sufficient for a fresh authorized account to resume. Assistant memory and private workspaces are never required.

This contract may be materially changed only by an approved Decision Log entry. Major-version closure includes whole-scope reconstruction, architecture/data/security/rights/performance/UX/adaptive review, actual artifact evidence, professional acceptance and G16 handoff.
