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
Perform adaptive retrospective, cleanup, learning, and next-version handoff. Preserve anything required for reproduction, rollback, audit, provenance, or historical certification.

---

## 15. G16 Feedback Loop

After every Stable build, feed:

`runtime evidence + user screenshots + defects + performance telemetry + provider effectiveness + calls avoided + intelligence/signal outcomes + UI/UX observations + data usefulness + release-process failures`

into the Adaptive Retrospective.

Classify findings as:
1. FIX NOW / NEXT-BUILD BLOCKER
2. NEXT BUILD
3. NEXT COMPATIBLE BUILD
4. ROADMAP / LATER HARDENING
5. REJECT / NO CHANGE WITH EVIDENCE

No reported issue should silently disappear.

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
