# DE.PULSE — Adaptive Operating Contract

**Status:** PERMANENT / GOVERNING  
**Scope:** Adaptive Roadmap, Adaptive Build Plan, Adaptive Build Process, Adaptive Delivery Process

This document defines the permanent operating contract. Release-specific scope may adapt, but these rules remain unless explicitly superseded by an approved Decision Log entry.

---

## 1. Four Connected Adaptive Layers

### Adaptive Roadmap
Defines **where DE.PULSE is going and why**.

### Adaptive Build Plan
Defines **what should be built next**, based on value, dependency, risk, defects, architecture, performance, test evidence, provider/data readiness, user feedback, and accumulated outcomes.

### Adaptive Build Process
Defines **how a build is engineered and qualified**.

### Adaptive Delivery Process
Defines **how certified source becomes a trustworthy Stable artifact and authoritative next-version baseline**.

Controlling loop:

**Observe → Learn → Prioritize → Build → Validate → Deliver → Measure → Learn again**

Product loop:

**Collect only when useful → Understand → Correlate → Filter → Rank → Explain → Measure outcome → Learn**

---

## 2. Canonical G0–G16 Mapping

No additional top-level gates may be created beyond G0–G16. New checks belong inside the appropriate gate.

- **G0 — Exact Baseline**
- **G1 — Immutable Scope**
- **G2 — Architecture / Data Utility**
- **G3 — Design / Dependency Readiness**
- **G4 — Development Exit**
- **G5 — FAST Qualification**
- **G6 — Integration / MEDIUM Qualification**
- **G7 — Data / Security / Adaptive Intelligence**
- **G8 — Performance / Capacity / Stability**
- **G9 — Cross-Module / UI / UX**
- **G10 — Pre-Freeze Qualification**
- **G11 — Immutable Release Candidate**
- **G12 — Full Certification**
- **G13 — Native Packaging / Provenance**
- **G14 — Actual Artifact Runtime Audit**
- **G15 — Release Assurance / Promotion**
- **G16 — Adaptive Retrospective / Handoff**

---

## 3. Build Engineering Order

Before adding new implementation, apply:

**REUSE → CONSOLIDATE → REFACTOR → DELETE/REPLACE → ADD**

Required consequences:
- prove no canonical owner already exists before creating another;
- no append-only architecture growth;
- remove obsolete/dead/superseded code when safe;
- preserve developer readability and testability;
- prefer one canonical data owner and one calculation owner;
- fetch once / calculate once / reuse/fan out where semantically equivalent and lawful.

---

## 4. Dependency-Aware CI/CD Execution

G0–G16 define governance order, not unnecessary serialization of every test.

Safe independent work and tests should run in bounded parallel when dependencies permit. Use a dependency-aware DAG and avoid repeating expensive equivalent evidence.

Evidence reuse is valid only when the relevant:
- frozen source fingerprint;
- test definition;
- required environment/toolchain assumptions;
- input contract

remain equivalent.

Any source/tooling change after freeze invalidates affected certification evidence. Never reuse stale PASS evidence across a materially different fingerprint.

A wrapper/orchestrator timeout is not automatically a product PASS or FAIL; inspect durable process/checkpoint evidence.

---

## 5. G0–G3 Before Implementation

### G0
Verify exact predecessor/baseline, repository hygiene, version identity, source fingerprint/provenance, and known defects/open issues.

### G1
Freeze the release-specific immutable scope. Approved permanent governance remains inherited unless explicitly excluded with evidence and approval.

### G2
Audit architecture/data utility before implementation. Required reasoning includes:

`dataset → canonical owner → consumers → purpose → refresh → retention → cost → fallback`

Reject duplicate owners, orphan data, unnecessary persistence, unnecessary provider calls, or data with no meaningful consumer.

### G3
Confirm design/dependency/provider/rights/test readiness and identify blockers before expensive implementation.

---

## 6. Adaptive Intelligence Contract

DE.PULSE is a smart, continuously improving research and decision-support system, not a static information terminal.

Permanent learning loop:

**Collect only when useful → Understand → Normalize → Correlate → Integrate → Remember → Compare → Reason → Prioritize → Explain → Measure outcome → Learn → Validate → Adapt**

Rules:
- deterministic critical paths remain bounded, fast, and explainable;
- adaptive intelligence enriches/calibrates without silently replacing protected production logic;
- adaptive production influence follows **SHADOW → VALIDATED → APPROVED → PRODUCTION**;
- no silent self-modifying production decision logic;
- historical replay/backtesting uses point-in-time truth and no lookahead;
- retain meaningful evidence/events/outcomes, not indiscriminate raw-tick warehousing;
- measure provider usefulness, false positives, misses, latency, freshness, calibration, outcome quality, and decision utility where relevant.

---

## 7. Adaptive Intelligence Scorecard

For intelligent capabilities, measure applicable dimensions such as:
- detection/decision latency;
- false positives;
- misses / false negatives;
- useful-alert rate;
- provider calls avoided;
- provider usefulness;
- freshness failures;
- evidence independence and contradiction quality;
- calibration;
- outcome quality;
- lead time;
- signal/model drift;
- decision utility;
- unnecessary processing avoided.

Do not optimize solely for headline accuracy.

---

## 8. Performance & Scalability Contract

Always consider:
- provider/API requests and subscription use;
- duplicate fetches/calculations and symbol fan-out;
- CPU, memory, allocations, GC, goroutines, locks/contention;
- queues/backpressure and bounded concurrency;
- DB reads/writes, indexes, storage growth;
- cache bounds/hit value;
- UI/API latency;
- background-job impact;
- long-running stability;
- restart/recovery;
- provider rate-limit pressure.

Prefer:
- memory-first live state;
- fetch once / calculate once;
- incremental computation;
- material-change propagation;
- asynchronous persistence where safe;
- bounded caches/queues;
- indexed DB access;
- background/overnight heavy computation;
- workload priority and graceful load shedding.

Critical live/current evidence must be protected before lower-priority work under pressure.

---

## 9. Data Utility / Evidence Value Contract

No datum is fetched, stored, computed, displayed, or retained merely because it is available.

Every significant datum must answer:
1. What decision/question does it help?
2. Who consumes it?
3. Is it independent information or duplicate evidence?
4. Should it be correlated/derived/aggregated instead?
5. Is it fresh and material?
6. Is it temporary, durable, or discardable?
7. What happens when absent, stale, degraded, or contradictory?

UI should surface material intelligence, changes, risks, contradictions, prioritized opportunities, confidence/freshness, and concise explanations—not implementation plumbing.

---

## 10. Canonical Provider / Shared Intelligence Contract

Maintain one canonical Provider Router and canonical dataset ownership.

Provider routing must be capability/entitlement-aware and distinguish:
- Preferred provider;
- Serving provider;
- routing/fallback reason;
- capability;
- entitlement;
- circuit/cooldown/rate-limit/saturation state.

Material cross-provider disagreement must be detected and handled conservatively, with provenance/freshness preserved.

Shared Symbol Intelligence principle:
- equivalent lawful processing for the same unique symbol should be canonical once;
- authorized users consume/fan out from shared canonical intelligence rather than causing duplicate provider/calculation work;
- user-specific membership/preferences remain isolated from shared symbol intelligence.

---

## 11. UI / UX Contract

**Complex inside → intelligent synthesis → simple outside.**

Normal USER/DEMO experiences should prefer:
- conclusions;
- material intelligence;
- confidence;
- freshness;
- meaningful degraded-data warnings;
- concise why/explanation;
- clear next research context.

Raw provider plumbing, queues, capacity, databases, circuit internals, and scheduler machinery belong in appropriate OWNER/SUPER_OWNER/ADMIN/Maintenance diagnostics rather than ordinary user surfaces.

Responsive/layout/typography integrity is a permanent G9 responsibility.

---

## 12. Product Boundary — No Execution

DE.PULSE is a **research + intelligence + decision-support system**.

Permanent No Execution Boundary:
- no order execution;
- no broker routing;
- no live/paper trading;
- no OMS/blotter;
- no autonomous or semi-autonomous trading;
- no portfolio/P&L or journal product scope unless explicitly superseded by a future approved product decision.

Options data may provide contextual intelligence but does not itself convert DE.PULSE into an options-trading product and must not silently alter protected deterministic Day/Swing/Long formulas.

---

## 13. Stable / RC Protection

- Historical Stable artifacts/tags are immutable provenance.
- TEST/RC work must not overwrite or destabilize the known-good Stable installation.
- Failed TEST/RC must leave Stable unaffected.
- Do not weaken tests or gates to obtain PASS.
- Never fabricate native/runtime evidence when a platform is unavailable.

---

## 14. Freeze, Packaging, Runtime, Provenance

### G10
Complete pre-freeze qualification before immutable RC creation.

### G11
Freeze the RC source identity. Any affected source/tooling change requires requalification of the new fingerprint.

### G12
Perform full certification on the frozen source with applicable unit/integration/E2E/race/randomized/security/performance/data/UI/professional acceptance coverage.

### G13
Build native artifacts and provenance from the certified frozen source.

### G14
Audit the **actual packaged artifact**: launch/run behavior, runtime identity, persistence/migration, platform behavior, configuration, critical flows, and packaging integrity.

### G15
Confirm release assurance, rollback, reproducibility, authoritative provenance, required evidence retention, and promotion truth before Stable promotion.

Generate final artifact hashes/SHA manifest **last**, after final contents are immutable.

### G16
Perform the full Stable-level Adaptive Retrospective, cleanup, learning, and next-version handoff. Aggregate the per-build AIPLC evidence defined below and preserve anything required for reproduction, rollback, audit, provenance, or historical certification.

---

## 15. Adaptive Intelligence & Product Learning Checkpoints

DE.PULSE must learn from **every meaningful build**, not only after a Stable release. The permanent checkpoint is **Adaptive Intelligence & Product Learning Checkpoint (AIPLC)**.

AIPLC is part of the existing Adaptive Build Process and G16 feedback loop. It does **not** create a new G17+ gate.

Permanent cycle:

**BUILD → TEST → OBSERVE → TAB × FEATURE × DATA INTELLIGENCE AUDIT → CHALLENGE / ROOT-CAUSE REVIEW → LEARN → IMPROVE → PREVENT → MEASURE → FEED NEXT BUILD**

### 15.1 When AIPLC runs

Run a concise AIPLC after every **meaningful build/checkpoint** that changes or materially exercises one or more of:
- user-visible behavior or information hierarchy;
- data acquisition, normalization, persistence, freshness or provenance;
- intelligence, ranking, recommendation, signal or explanation behavior;
- provider routing/fallback/cost/capacity;
- architecture/canonical ownership;
- performance/reliability/concurrency/background work;
- security/RBAC/session behavior;
- testing/certification/build/release machinery;
- a defect class whose root cause can teach a reusable lesson.

A mechanically identical rerun with no new evidence may reuse the latest applicable AIPLC with an explicit **NO NEW LEARNING / EVIDENCE EQUIVALENT** disposition; do not manufacture redundant reports.

### 15.2 Tab × Feature × Datum Intelligence Audit

For every affected tab/surface and every material feature/data element, ask whether DE.PULSE is behaving like a smart AI/LLM-style research system or merely displaying available information.

Use the chain:

**datum → purpose → canonical owner → consumer → freshness/materiality → independence/correlation → interpretation → decision value → explanation → outcome → learning**

Audit applicable surfaces including Dashboard, Day, Swing, Long, Discovery / Opportunity Radar, Research, Decision Queue, Market Intelligence, Administration, Settings, Maintenance and future surfaces.

For each affected element determine whether to:
- **KEEP / STRENGTHEN** — clear material decision value;
- **SYNTHESIZE / CORRELATE** — raw evidence should become higher-level intelligence;
- **CONSOLIDATE / REUSE** — duplicate data/calculation/UI ownership exists;
- **DEMOTE / DRILL-DOWN** — useful support but too prominent;
- **SUPPRESS / REMOVE** — low-value, stale, redundant or confusing;
- **ABSTAIN / UNKNOWN** — evidence is insufficient for a reliable conclusion.

Raw availability is never sufficient justification for collection, computation, persistence or UI prominence.

### 15.3 Challenge → Root Cause → Prevention

Every material defect/challenge should produce **two outputs**:

1. the immediate fix or truthful disposition;
2. the reusable learning/prevention that makes the same class of failure less likely.

Required reasoning:

**observed symptom → root cause → affected canonical owner → immediate fix → recurrence prevention → cross-product pattern scan → test/gate/contract hardening → measurable follow-up**

Do not stop at symptom repair when the same root cause can exist elsewhere.

Examples of reusable learning include:
- duplicate calculations → consolidate the canonical calculation owner and scan equivalent paths;
- runtime overload → improve workload priority/backpressure/provider budgeting and protect critical evidence;
- repeated UI clutter → improve materiality/information hierarchy rules across related surfaces;
- provider failures → improve capability routing/fallback/reason semantics rather than patching one consumer;
- recurring test escape → strengthen the architecture invariant or reusable test oracle, not only one regression case.

### 15.4 Ten-Dimension AIPLC Quality Score

Score applicable dimensions with evidence; **10/10 is a target and quality threshold, never an automatic label**:

1. **AI/LLM Smartness** — synthesis, reasoning, prioritization and useful explanation vs raw display;
2. **Data Utility** — every material datum has justified purpose/consumer/value;
3. **Canonical Architecture** — reuse, ownership clarity, no duplicate engines/calculations;
4. **Decision Utility / User Value** — materially improves research/decision support;
5. **Reliability / Freshness Truth** — graceful scoped degradation, no false confidence;
6. **Performance / Efficiency** — bounded cost, calls avoided, no unnecessary processing;
7. **UI/UX Intelligence** — complex inside, simple outside, correct hierarchy;
8. **Learning / Outcome Measurement** — useful outcomes, misses, false positives, drift and provider value are measurable where applicable;
9. **Testing / Prevention** — root-cause regression protection and cross-module coverage;
10. **Adaptive Continuity** — finding is fixed, explicitly carried forward, or rejected with evidence; nothing silently disappears.

If any applicable dimension is below 10:
- improve it immediately when safe, dependency-compatible and within frozen scope;
- rerun the affected checks and rescore;
- if it cannot safely reach 10 in the current release, record the exact gap, evidence, risk and named disposition (`NEXT BUILD`, `NEXT COMPATIBLE BUILD`, `ROADMAP`, or `REJECT / NO CHANGE WITH EVIDENCE`);
- never inflate the score merely to call the build 10/10.

A build may still proceed when a below-10 finding is legitimately deferred and non-blocking, but the checkpoint itself must remain truthful about the residual.

### 15.5 AIPLC Can Change the Next Build

AIPLC is not ceremonial. Its validated learning must be allowed to alter the next Adaptive Build Plan.

Examples:
- reduce/demote a provider or dataset that produces little unique value;
- strengthen a provider/capability that measurably improves usefulness/reliability;
- remove or consolidate redundant data/UI/code;
- change workload priority based on measured pressure;
- keep weak intelligence in SHADOW or demote it;
- strengthen intelligence only after outcome/calibration evidence;
- promote a repeated defect pattern into architecture/test/process prevention;
- prioritize a user-facing simplification when internal complexity leaks outward.

Production adaptive intelligence still follows **SHADOW → VALIDATED → APPROVED → PRODUCTION**. AIPLC does not authorize silent self-modifying production logic or unapproved product scope.

### 15.6 Scope / Freeze Safety

AIPLC may identify work at any time, but it does not bypass G1.

- safe fixes/hardening already inside frozen scope may close in the current build;
- genuine escaped defects that materially threaten correctness, security, reliability or release truth may block promotion and require governed correction/requalification;
- source-changing new scope discovered after G1 normally becomes a named next-build/next-compatible-build entry;
- future strategic intelligence remains roadmap work until dependency-ready and approved where required.

### 15.7 Checkpoint Evidence

Keep AIPLC evidence concise and durable. Prefer a machine-readable/structured checkpoint plus a short human summary rather than a large repetitive report.

Minimum evidence:
- build/source identity;
- affected surfaces/capabilities;
- challenges/findings;
- root causes;
- changes made;
- reusable prevention added;
- 10-dimension scores and residuals;
- next-build dispositions;
- metrics/outcomes to watch;
- explicit statement of whether any current promotion blocker remains.

The next build should consume unresolved AIPLC items during G0/G1/G2 rather than rediscovering them from memory.

### 15.8 Stable G16 Deep Retrospective

After every Stable build, G16 aggregates all per-build AIPLC checkpoints plus:

`runtime evidence + user screenshots + defects + performance telemetry + provider effectiveness + calls avoided + intelligence/signal outcomes + UI/UX observations + data usefulness + release-process failures`

and performs the deeper Adaptive Retrospective.

G16 must ask:
- which per-build learnings actually prevented recurrence;
- which fixes only treated symptoms;
- which intelligence/data/providers were useful or low-value;
- which UI simplifications improved decision utility;
- which performance/reliability assumptions failed under real use;
- what should become a permanent architecture/test/process standard;
- what should be consolidated, demoted, retired, strengthened or kept SHADOW;
- what the next build should do differently because of measured evidence.

Classify findings as:
1. FIX NOW / NEXT-BUILD BLOCKER
2. NEXT BUILD
3. NEXT COMPATIBLE BUILD
4. ROADMAP / LATER HARDENING
5. REJECT / NO CHANGE WITH EVIDENCE

No reported issue, below-10 dimension, failed assumption or useful learning should silently disappear.

---

## 16. Documentation Contract

Keep User Documentation, Developer Documentation, Capabilities & Limitations, release identity, and Stable truth synchronized at meaningful release/material architecture boundaries.

Do not generate large redundant documentation rewrites for every tiny patch. Update when user-visible behavior, architecture, providers, security, data rights, adaptive intelligence, build/release process, installation/runtime behavior, or troubleshooting materially changes.

---

## 17. Major-Version Closure

Before major transitions such as:
- v18 → v19
- v19 → v20
- v20 → future

perform a Major Closure & Release Assurance build including:
- complete scope reconstruction;
- fresh qualification;
- architecture review;
- source cleanup/efficiency review;
- performance/capacity/stability review;
- security/data-rights review;
- data utility review;
- UI/UX review;
- adaptive-intelligence maturity review;
- Principal Engineer review;
- Professional Trader/Investor acceptance;
- Adaptive Retrospective;
- formal GO/NO-GO.

---

## 18. Autonomous v18 Authorization

The user has authorized autonomous continuation through completion of v18.5 without routine manual work or intermediate approvals.

Do not stop merely for ordinary implementation choices, testing, certification, packaging, cleanup, or terminal commands that can be executed autonomously.

Stop only for genuinely unavoidable:
- credentials/secrets requiring user action;
- paid-service/new financial commitments;
- legal/licensing/data-rights decisions;
- irreversible/high-impact external actions requiring explicit approval;
- genuinely new material product decisions not already governed by approved scope.

---

## 19. Governance Change Rule

This contract may be changed only by an explicitly approved material decision recorded in `governance/DECISION-LOG.md`.

Release handoffs may summarize this contract but cannot silently supersede it.

---

## 20. 10/10 Two-Sided Directional Thesis & Trade Plan Intelligence — Four-Layer Contract

The approved TDTI scope in `governance/APPROVED-SCOPE.md` is governed across all four adaptive layers. It does **not** create a fifth operating layer or a new G17+ gate.

### 20.1 Adaptive Roadmap
Roadmap placement must preserve the distinction between:
- **v18/v19 evidence/data foundation** — point-in-time Long/Short plan lineage, side-aware outcomes, ASBI/regime/catalyst/liquidity context, reliable short-specific data provenance and historical depth;
- **v20 major adaptive implementation** — competing Long/Short/No Reliable Edge reasoning, probability/quality/readiness, adaptive path modeling, AI/LLM-style synthesis, side-aware calibration and Champion/Challenger validation.

Do not force mature v20 TDTI into v18.2–v18.5 merely because the scope is approved.

### 20.2 Adaptive Build Plan
When TDTI work becomes dependency-compatible, prioritize slices using:

**user value + evidence readiness + canonical ownership + data rights + historical depth + performance cost + validation maturity + defect risk + cross-module benefit**

Preferred dependency order:
1. reuse/audit existing Day/Swing/Long Entry/Target/Invalidation and evidence-snapshot owners;
2. freeze a canonical two-sided thesis/trade-plan schema and semantics;
3. harden side-aware outcome measurement and no-lookahead history;
4. add short-specific contextual datasets only where useful/reliable/lawful;
5. integrate ASBI state/path and cross-horizon reasoning;
6. introduce learned probability/quality/readiness in SHADOW;
7. add AI/LLM-style structured synthesis and contradiction reasoning;
8. validate Champion/Challenger before production influence;
9. expose concise contextual UI only after truth/performance is proven.

Do not create a separate Short Desk, second market-data pipeline, second historical store, or duplicate deterministic horizon engine merely to implement the short side.

### 20.3 Adaptive Build Process — G0–G16 Responsibilities
TDTI uses the existing gate model:

- **G0 Exact Baseline:** identify current Long-plan formulas, prior short-entry direction, validation/outcome owners, open defects, and exact source fingerprint.
- **G1 Immutable Scope:** freeze which TDTI slice is in the release and what remains roadmap-only.
- **G2 Architecture / Data Utility:** prove one canonical evidence/thesis owner; classify every proposed component REUSE / CONSOLIDATE / REFACTOR / DELETE-REPLACE / ADD; justify short-interest/borrow/SSR or other data by consumer/value/rights.
- **G3 Design / Dependency Readiness:** freeze Long/Short semantics, structural trigger/invalidation rules, UNKNOWN/ABSTAIN behavior, data/provider/rights dependencies, and test oracle before implementation.
- **G4 Development Exit:** unit/static/schema tests must prove directional arithmetic/sign semantics, structural levels, no unsupported LLM level invention, and backward compatibility for protected outputs where no approved change applies.
- **G5 FAST Qualification:** smoke/regression on Long, Short and NO RELIABLE EDGE states, including missing/stale/degraded evidence.
- **G6 Integration / MEDIUM Qualification:** verify ASBI, Day/Swing/Long, Research, Opportunity Radar, Decision Queue, Event/Market Intelligence, Historical Validation and Smart-Money integration without duplicate computation.
- **G7 Data / Security / Adaptive Intelligence:** enforce point-in-time truth, no lookahead, evidence independence, provenance/freshness, ABSTAIN, LLM grounding, side-aware immutable ledger, SHADOW/Champion-Challenger governance and no silent self-modification.
- **G8 Performance / Capacity / Stability:** measure incremental CPU/memory/provider/storage/latency cost, duplicate calculations, material-change propagation, bounded background learning and long-running stability. Shared Long/Short reasoning should reuse the same canonical evidence rather than double workload naïvely.
- **G9 Cross-Module / UI / UX:** prove clear Long-vs-Short labels, no mirrored/ambiguous terminology, responsive trade-plan geometry, concise strongest-thesis + competing-thesis presentation, accessible uncertainty/UNKNOWN and no execution-looking controls.
- **G10 Pre-Freeze Qualification:** require side-aware deterministic/invariant tests, no-lookahead replay, outcome-ordering tests, calibration evidence where claimed, contradiction/ABSTAIN tests and Professional Trader/Investor review before RC freeze.
- **G11 Immutable RC:** freeze the exact thesis/model/prompt/rule version and evidence contract.
- **G12 Full Certification:** run full regression plus walk-forward/point-in-time evaluation, Long vs Short calibration, false positives/misses, target/invalidation ordering, MFE/MAE, Decision Utility, regime robustness, squeeze/bounce/bull-trap failure cases, cross-horizon cases, AI grounding and ABSTAIN quality as applicable.
- **G13 Native Packaging / Provenance:** package the exact certified model/rule/prompt/schema identities and required provenance.
- **G14 Actual Artifact Runtime Audit:** launch the packaged app and verify real rendered Long/Short plans, state transitions, UNKNOWN/degraded behavior, restart persistence and no execution capability.
- **G15 Release Assurance / Promotion:** prove rollback/reproducibility, production-vs-SHADOW truth and that no Challenger or LLM logic was silently promoted.
- **G16 Adaptive Retrospective / Handoff:** feed side-specific outcomes, false positives/misses, MFE/MAE, thesis flips, ABSTAIN quality, alert usefulness, short-specific data usefulness, provider cost and UI findings into the next Adaptive Build Plan.

### 20.4 Adaptive Delivery Process
A delivered TDTI-capable artifact is trustworthy only when the **actual package**, not just source tests, proves applicable:
- Long and Short are derived from the same point-in-time canonical evidence contract;
- Entry/Short Entry, Trim/Target/Cover and Invalidation semantics are directionally correct and structurally explainable;
- unsupported/missing shortability, borrow, short-interest or similar evidence is UNKNOWN rather than guessed;
- target/cover-before-entry cannot count as success in validation;
- historical replay does not leak future news, filings, earnings, fundamentals, market state or provider corrections;
- AI/LLM synthesis is grounded in available evidence, preserves contradictions/uncertainty, and cannot silently overwrite canonical truth;
- NO RELIABLE EDGE / ABSTAIN works in normal and degraded states;
- Champion/Challenger and SHADOW/production identities are visible in provenance/diagnostics where applicable;
- packaged runtime preserves restart/persistence truth and does not expose order-entry, borrow, portfolio, P&L or execution controls;
- UI remains **complex inside → intelligent synthesis → simple outside**.

### 20.5 TDTI Adaptive Scorecard
Measure applicable:
- Long vs Short calibration separately;
- thesis-strength usefulness;
- Opportunity Quality / Decision Utility;
- readiness precision and lead time;
- target/cover/invalidation first-event truth;
- MFE/MAE and time-to-resolution;
- false-positive and missed-material-move cost by side;
- squeeze/rebound/bull-trap/failed-breakdown failure modes;
- cross-horizon contradiction quality;
- ASBI state-transition usefulness;
- evidence independence/contradiction handling;
- LLM grounding / unsupported-claim rate;
- ABSTAIN quality;
- alert burden/usefulness;
- provider/data usefulness and calls avoided;
- latency/resource overhead;
- regime robustness/drift;
- Champion/Challenger evidence;
- Professional Trader/Investor acceptance.

Headline directional accuracy alone is insufficient for promotion.