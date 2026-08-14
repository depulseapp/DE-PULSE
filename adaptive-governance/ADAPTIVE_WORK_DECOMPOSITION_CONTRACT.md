# DE.PULSE — Adaptive Work Decomposition & G0–G16 Efficiency Contract

Status: **Permanent governing contract**  
Applies to: Adaptive Roadmap, Adaptive Build Plan, Adaptive Build Process, Adaptive Delivery Process, Adaptive CI Operating Contract  
Effective: **v18.2 and all later releases**

## Purpose

DE.PULSE engineering work follows the same adaptive-intelligence principle as the product itself. Heavy implementation, qualification, audit and delivery work should be decomposed into the smallest independently meaningful evidence units when that improves speed, recovery, diagnosis, reuse or learning.

Permanent execution loop:

**Understand → impact-map → decompose → reuse → execute → checkpoint → evaluate → adapt → integrate → certify → learn.**

This contract does **not** create or permit new top-level release gates. **G0–G16 is permanent.** New responsibilities must be represented as checkpoints, sub-stages, shards, bounded parallel lanes, or strengthened ownership inside an existing G0–G16 gate.

---

## 1. Decompose heavy work only when useful

Before a heavy task begins, the Build Coordinator decides whether it should remain one unit or be split into smaller work packages.

Decomposition is preferred when it materially improves:
- fault isolation/root-cause diagnosis;
- exact-source evidence reuse after partial failure;
- safe bounded parallel execution;
- interruption recovery/resumability;
- ownership clarity;
- test/evidence traceability;
- avoidance of one slow/flaky responsibility blocking unrelated work;
- reuse of already-qualified evidence;
- ability to learn from intermediate outcomes and adapt the next action.

Do not split work merely to create more jobs. A split that adds coordination, CPU, memory, provider/API, browser, DB or artifact load without better assurance is rejected.

---

## 2. Dependency-aware execution

Every decomposed package has:
- one canonical owner/responsibility;
- explicit inputs/dependencies;
- source/fingerprint/artifact identity where relevant;
- PASS / FAIL / BLOCKED or completion criteria;
- durable evidence location;
- downstream consumers;
- invalidation rules when inputs change.

Independent packages may execute in parallel. Dependent packages remain ordered. The Build Coordinator owns the dependency graph and prevents concurrent mutation of the same canonical release state.

Parallelism is always bounded by actual runner/provider/runtime/database/browser capacity.

---

## 3. Delta-first impact model

Before rerunning expensive work, classify what actually changed:

**Git diff → canonical owner → dependency blast radius → affected surfaces/data/features/roles/providers → reusable evidence → smallest required rerun set.**

Unchanged evidence may be inherited only when its relevant source fingerprint/artifact identity, dependency contract, test definition, role/security assumptions and input semantics remain equivalent.

Metadata/checkpoint-only changes must not trigger product requalification when they are explicitly excluded from the candidate fingerprint and do not change the relevant test contract.

---

## 4. Checkpoint-first recovery

A heavy phase should expose meaningful intermediate checkpoints whenever a later failure would otherwise force unnecessary repetition.

A checkpoint is reusable only when its relevant inputs remain unchanged. It is not a substitute for a required G0–G16 conclusion.

On failure/interruption:

**inspect actual state → identify smallest affected package → preserve unrelated PASS evidence → repair → rerun affected/dependent work only → continue.**

Conversation interruption alone never invalidates unchanged-source evidence.

---

## 5. Three-depth review model

### Level 1 — Every meaningful build: Delta AIPLC
Review changed/affected areas plus dependency sentinels only.

### Level 2 — G10: Full coverage reconciliation
Every required tab/feature/data/security/performance responsibility must be either:
- freshly evidenced on the current candidate; or
- explicitly inherited from equivalent trustworthy evidence.

### Level 3 — G16 / Major Closure: Deep system review
Perform broad cross-product learning, consolidation and process review.

This prevents repeated brute-force audits after every small change while preserving complete release coverage.

---

## 6. Efficient test/load rules

Prefer:
- G5 changed-area FAST tests;
- G6 affected integration tests;
- bounded independent G7/G8/G9 shards;
- deterministic fixtures/replay/canonical cached evidence when live provider behavior is not the subject under test;
- reusable evidence packages/fingerprints for AI/LLM grounding tests;
- small bounded true-model/runtime samples only when actual LLM behavior must be certified;
- independent macOS Apple Silicon and Windows x64 packaging/runtime lanes;
- reuse of an unchanged platform PASS when the exact RC/package identity is unchanged.

Do not:
- rerun full G12 for checkpoint/documentation metadata only;
- duplicate provider acquisition across test lanes when equivalent evidence can be shared;
- rerun materially identical AI synthesis solely because another test/user opens the same symbol;
- run full all-tab screenshot/manual-style review after every commit.

---

## 7. G0–G16 integration

- **G0:** exact baseline/repository/checkpoint truth.
- **G1:** immutable release scope and applicable permanent governance.
- **G2:** canonical ownership/data utility/impact graph.
- **G3:** dependencies, evidence plan, decomposition and bounded lanes.
- **G4:** implementation exit using independently testable packages where useful.
- **G5:** FAST affected-area qualification.
- **G6:** affected integration/MEDIUM qualification.
- **G7:** data/security/adaptive-intelligence affected lanes.
- **G8:** performance/capacity/stability affected lanes.
- **G9:** cross-module/UI/UX affected lanes plus required coverage sentinels.
- **G10:** authoritative full coverage reconciliation before freeze.
- **G11:** immutable RC identity.
- **G12:** full certification on immutable RC.
- **G13/G14:** independent native packaging/runtime lanes.
- **G15:** assurance/promotion consumes the complete evidence graph.
- **G16:** deep retrospective, consolidation, learning and next-release handoff.

No G17+ may be introduced.

---

## 8. Adaptive execution without scope drift

The Build Coordinator may adapt sequencing, test focus, parallelism and remediation based on evidence, provided it preserves:
- frozen G1 scope;
- permanent product/architecture/security/data-rights contracts;
- required G0–G16 assurance;
- No Execution boundary;
- truthful source/artifact provenance.

AIPLC findings outside frozen scope are carried to a named next build/next compatible build unless they reveal a genuine correctness/security/reliability blocker that requires governed correction and requalification.

---

## 9. Success scorecard

At G16, measure whether decomposition actually improved:
- cycle time;
- reruns avoided;
- provider/model calls avoided;
- runner/resource load;
- failure isolation;
- resumability;
- evidence reuse;
- duplicate work removed;
- recurrence prevention;
- clarity of ownership.

Redundant packages/checks are consolidated or removed.

---

## 10. Permanent rule

**DE.PULSE solves heavy engineering work by understanding what changed, reusing what is already trustworthy, splitting only where evidence benefits, executing independent work with bounded parallelism, rerunning the smallest affected set, integrating one authoritative result, and learning. G0–G16 remains the permanent top-level gate model.**
