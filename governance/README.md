# DE.PULSE Governance — Canonical Entry Point

**Status:** ACTIVE / AUTHORITATIVE FOR APPROVED INTENT  
**Repository:** `depulseapp/DE-PULSE`  
**Purpose:** Prevent approved product, architecture, roadmap, build-process, release-gate, delivery, historical Stable truth, and unresolved defects from being lost across chats, handoffs, branches, refactors, or releases.

---

## 1. Source-of-Truth Hierarchy

Use the following hierarchy whenever DE.PULSE scope or a prior decision is questioned:

1. **Actual source / packaged release evidence** — truth about what is actually implemented and delivered.
2. **`governance/APPROVED-SCOPE.md`** — canonical base truth about approved product/roadmap scope.
3. **`governance/CONTINUITY-IMPLEMENTATION-CONTRACT.md`** — canonical carry-forward + governance-to-implementation contract for previously approved/certified capabilities, ownership/RBAC/platform rules, active release obligations, unresolved defect continuity, and future-loss prevention.
4. **`governance/ADAPTIVE-OPERATING-CONTRACT.md`** — canonical permanent engineering/build/release/delivery contract.
5. **`governance/ADAPTIVE-DATA-RELIABILITY-CONTRACT.md`** — canonical 10/10 Adaptive Data Reliability & Graceful Degradation Intelligence (ADR-GDI) contract across roadmap/build/delivery/runtime reliability.
6. **`governance/ROADMAP.md`** — canonical placement and sequencing of approved future work.
7. **`governance/DECISION-LOG.md`** — append-only history of approved, superseded, rejected, or deferred material decisions.
8. **`release/<version>/G1-IMMUTABLE-SCOPE.md`** — immutable scope snapshot for one specific release after G1.
9. **G16 handoff / release evidence** — what actually happened in that release and what carries forward.
10. **Historical certified Stable traceability / Major Closure evidence** — inherited implementation truth that must not silently vanish during governance compression or refactoring.
11. Chat memory or recollection — continuity aid only; never override the canonical GitHub records above.

A governance document never proves a feature exists in code. A release/code artifact never silently changes approved product intent. Absence from a newer condensed summary does not by itself retire an approved or certified capability.

---

## 2. Required Discussion Workflow

Before treating a proposed DE.PULSE idea as new, use:

**LOOKUP → COMPARE → CLASSIFY → DECIDE → UPDATE → IMPLEMENT/DISPOSITION → EVIDENCE → LEARN**

### LOOKUP
Read/search at minimum:
- `governance/APPROVED-SCOPE.md`;
- `governance/CONTINUITY-IMPLEMENTATION-CONTRACT.md`;
- `governance/ADAPTIVE-OPERATING-CONTRACT.md`;
- `governance/ADAPTIVE-DATA-RELIABILITY-CONTRACT.md` when reliability, data freshness, provider/runtime/database pressure, degraded-state semantics, backpressure or recovery is relevant;
- `governance/ROADMAP.md`;
- `governance/DECISION-LOG.md`.

When release-specific, also inspect the current `release/<version>/G1-IMMUTABLE-SCOPE.md`, active branch implementation evidence, latest G16/handoff, and relevant inherited Stable traceability.

### COMPARE
Compare the proposal semantically, not only by exact wording. Check aliases, renamed features, overlapping responsibilities, historical implementation, current canonical owners, and whether a newer summary merely omitted an older approved/certified capability.

### CLASSIFY
Every proposal/finding must be classified as one of:

- **ALREADY APPROVED** — existing scope covers it.
- **PARTIALLY COVERED** — existing scope covers part; identify the delta.
- **REFINEMENT / HARDENING** — same approved purpose, stronger implementation or acceptance.
- **CONFLICT / SUPERSESSION** — contradicts or replaces an earlier decision.
- **NEW SCOPE** — materially new product/architecture behavior.
- **REJECT / NO CHANGE** — not useful or violates a permanent boundary.

For implementation disposition, also classify as applicable:
- **CURRENT_RELEASE_BLOCKER**;
- **CURRENT_RELEASE_PROCESS_HARDENING**;
- **NEXT_RELEASE_MANDATORY_ENTRY**;
- **FUTURE_STRATEGIC**;
- **ALREADY_IMPLEMENTED / INHERITED**;
- **OPEN_DEFECT / RECONCILE**.

### DECIDE
Do not silently add material scope. User approval is required for genuinely new material product decisions or supersession of an approved contract.

### UPDATE
After approval:
1. update the relevant canonical governance document/contract;
2. append a Decision Log entry for material decisions;
3. place the work in the Adaptive Roadmap/Build Plan or explicit release disposition;
4. when a release reaches G1, snapshot the applicable scope into its immutable G1 file;
5. preserve inherited Stable capability/defect/platform/rights truth where applicable.

### IMPLEMENT / DISPOSITION
`APPROVED` is not `IMPLEMENTED`.

Applicable approved items must either become real source/runtime/test work or be explicitly carried to a named future release. Documentation-only closure is forbidden.

### EVIDENCE / LEARN
G10/G12/G14/G16 must be able to trace applicable requirements through source, tests, actual package/runtime and outcomes. G16 then feeds usefulness, defects, provider/data performance, reliability, UI findings and implementation gaps into the next Adaptive Build Plan.

---

## 3. Do Not Use Handoffs as the Only Memory

Handoffs are release continuation records, not the master product contract. They may summarize only what was relevant to that release.

Permanent decisions belong in `governance/` and are referenced by future handoffs rather than repeatedly copied in full.

Historical handoffs and certified traceability remain evidence sources for detecting an older approved/certified capability that was accidentally omitted from newer condensed governance.

---

## 4. Release Relationship

The permanent model is:

`Governance approved intent + continuity contract + inherited Stable truth`
→ `Adaptive Roadmap / Build Plan`
→ `G0 baseline + omission/divergence audit`
→ `G1 immutable release scope snapshot`
→ `G2–G15 engineering / certification / delivery`
→ `G16 retrospective + handoff`
→ `governance reconciliation when a permanent decision changed`

No release may silently drop an approved permanent contract, inherited certified capability, or unresolved user-reported defect merely because it was omitted from a handoff summary.

---

## 5. Executable Governance Rule

Permanent closure chain:

**Governed → Implemented → Enforced → Evidenced → Delivered → Learned**

At G10/G16 and Major Closure, explicitly detect:
- approved but never scheduled;
- scheduled but never implemented;
- implemented but not integrated;
- integrated but not tested;
- source-tested but not packaged/runtime-proven;
- runtime-proven but ineffective/noisy;
- useful and worth strengthening;
- redundant/obsolete and safe to consolidate/supersede;
- unresolved defect that has no closure evidence.

The full contract is `governance/CONTINUITY-IMPLEMENTATION-CONTRACT.md`.

---

## 6. Efficiency Rule

Do **not** rewrite all governance documents for every small patch.

Update governance only when:
- approved product scope changes;
- a permanent architecture/build/release/delivery contract changes;
- roadmap placement materially changes;
- a decision is superseded/rejected;
- G16 discovers a durable rule that must carry forward;
- a reconciliation proves an older approved/certified capability or unresolved defect is missing from current canonical continuity.

Release-specific implementation evidence remains under the release/QA evidence structure.

---

## 7. 10/10 Continuity Audit Rule

A reconciliation may be described as **10/10** only after all ten dimensions in `governance/CONTINUITY-IMPLEMENTATION-CONTRACT.md` pass: approval coverage, Stable inheritance, canonical ownership, release placement, executable traceability, defect continuity, security/rights/role continuity, performance/reliability continuity, delivery/platform truth, and future-loss prevention.

A score below 10 must be improved before the reconciliation itself is called 10/10. This does not mean every future roadmap capability is already implemented.

---

## 8. Permanent Review Question

Before every new build plan or substantial DE.PULSE discussion, ask:

> **Is this already approved, inherited/certified, partially covered, a refinement, conflicting, genuinely new, or still an unresolved defect — where is the canonical evidence, what release disposition applies, and how will implementation/runtime proof be obtained?**