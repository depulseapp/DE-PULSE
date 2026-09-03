# DE.PULSE Governance — Canonical Entry Point

**Status:** ACTIVE / AUTHORITATIVE FOR APPROVED INTENT  
**Repository:** `depulseapp/DE-PULSE`  
**Purpose:** Prevent approved product, architecture, roadmap, build-process, release-gate, delivery, historical Stable truth, and unresolved defects from being lost across chats, handoffs, branches, refactors, or releases.

---

## 1. Source-of-Truth Hierarchy

Use the following hierarchy whenever DE.PULSE scope or a prior decision is questioned:

1. **Actual source / packaged release evidence** — truth about what is actually implemented and delivered.
2. **`governance/APPROVED-SCOPE.md`** — canonical base truth about approved product/roadmap scope.
3. **`governance/AI-ASSISTANT-PORTABILITY-CONTRACT.md`** — permanent GitHub-backed, vendor/account-independent resume and handoff contract.
4. **`governance/CONTINUITY-IMPLEMENTATION-CONTRACT.md`** — canonical carry-forward + governance-to-implementation contract for previously approved/certified capabilities, ownership/RBAC/platform rules, active release obligations, unresolved defect continuity, and future-loss prevention.
5. **`governance/ADAPTIVE-OPERATING-CONTRACT.md`** — canonical permanent engineering/build/release/delivery contract.
6. **Full-product audit rebaseline and machine registers** — `governance/PRODUCT_AUDIT_REBASELINE_2026_08_27.md`, `governance/PRODUCT_AUDIT_COVERAGE_2026_08_27.md`, and the finding/coverage/5×5 JSON under `governance/programs/V19-V20-REBASELINE/`; these conserve the audit without claiming implementation.
7. **Canonical adaptive narratives** — `governance/ROADMAP.md`, `adaptive-governance/ADAPTIVE_BUILD_PLAN.md`, `adaptive-governance/ADAPTIVE_BUILD_PROCESS.md`, and `adaptive-governance/ADAPTIVE_DELIVERY_PROCESS.md`; their disposition index is `adaptive-governance/README.md`.
8. **`governance/ADAPTIVE-CI-QUALIFICATION-CONTRACT.md`** — permanent checkpoint-based CI/CD execution contract across the four adaptive layers; optimizes execution timing without reducing evidence or release quality.
9. **`governance/ADAPTIVE-DATA-RELIABILITY-CONTRACT.md`** — canonical Adaptive Data Reliability & Graceful Degradation Intelligence contract across roadmap/build/delivery/runtime reliability.
10. **`governance/DECISION-LOG.md` plus approved material decision records under `governance/DECISION-*.md`** — append-only history of approved, superseded, rejected, or deferred material decisions.
11. **Current execution state** — live GitHub PR/head/check/artifact evidence, `governance/current-state.json`, the active closure ledger, and `handoff/CURRENT.md`; these own exact progress and the one next action, not permanent scope.
12. **`release/<version>/G1-IMMUTABLE-SCOPE.md`** — immutable scope snapshot for one specific release after G1.
13. **Historical certified Stable traceability / Major Closure evidence** — inherited implementation truth that must not silently vanish during governance compression or refactoring.
14. Chat memory or recollection — continuity aid only; never override the canonical GitHub records above.

A governance document never proves a feature exists in code. A release/code artifact never silently changes approved product intent. Absence from a newer condensed summary does not by itself retire an approved or certified capability.

---

## 2. Required Discussion Workflow

Before treating a proposed DE.PULSE idea as new, use:

**LOOKUP → COMPARE → CLASSIFY → DECIDE → UPDATE → IMPLEMENT/DISPOSITION → EVIDENCE → LEARN**

### LOOKUP
Read/search at minimum:
- `governance/APPROVED-SCOPE.md`;
- `governance/AI-ASSISTANT-PORTABILITY-CONTRACT.md` when resuming, changing assistant/account, or producing a handoff;
- `governance/CONTINUITY-IMPLEMENTATION-CONTRACT.md`;
- `governance/ADAPTIVE-OPERATING-CONTRACT.md`;
- `adaptive-governance/README.md` and its single-authority map;
- `governance/PRODUCT_AUDIT_REBASELINE_2026_08_27.md` plus the finding register when audit scope is relevant;
- `governance/ADAPTIVE-CI-QUALIFICATION-CONTRACT.md` when build/CI/CD/checkpoint/cost/qualification/release execution is relevant;
- `governance/ADAPTIVE-DATA-RELIABILITY-CONTRACT.md` when reliability, data freshness, provider/runtime/database pressure, degraded-state semantics, backpressure or recovery is relevant;
- `governance/ROADMAP.md`;
- `adaptive-governance/ADAPTIVE_BUILD_PLAN.md`, `adaptive-governance/ADAPTIVE_BUILD_PROCESS.md`, and `adaptive-governance/ADAPTIVE_DELIVERY_PROCESS.md` when planning or executing a build;
- `governance/DECISION-LOG.md` and applicable `governance/DECISION-*.md` material decision records.

When release-specific, also inspect the current `release/<version>/G1-IMMUTABLE-SCOPE.md`, active branch implementation evidence, latest G16/handoff, relevant qualification checkpoint, and relevant inherited Stable traceability.

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
2. append/register a material Decision Log record;
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

`adaptive-governance/CURRENT_ADAPTIVE_*.md` files are compatibility/status projections only. They must not contain an independent roadmap, build plan, process, delivery contract, gap ledger or next action. Update exact execution state in `governance/current-state.json`, the active closure ledger and `handoff/CURRENT.md`; update durable intent in the single canonical narrative named by `adaptive-governance/README.md`.

Historical handoffs and certified traceability remain evidence sources for detecting an older approved/certified capability that was accidentally omitted from newer condensed governance.

---

## 4. Release Relationship

The permanent model is:

`Governance approved intent + continuity contract + inherited Stable truth`
→ `Adaptive Roadmap / Build Plan`
→ `development batches with tests implemented alongside source`
→ `deliberate exact-SHA Development Checkpoint`
→ `G0–G5 FAST qualification`
→ `G6–G10 qualification only after FAST green`
→ `G11 immutable RC + G12–G15 certification/delivery`
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

Do **not** rewrite all governance documents or run exhaustive CI for every small patch.

Update governance only when:
- approved product scope changes;
- a permanent architecture/build/release/delivery contract changes;
- roadmap placement materially changes;
- a decision is superseded/rejected;
- G16 discovers a durable rule that must carry forward;
- a reconciliation proves an older approved/certified capability or unresolved defect is missing from current canonical continuity.

CI/CD follows `governance/ADAPTIVE-CI-QUALIFICATION-CONTRACT.md`: normal development commits are quiet by default, qualification is checkpoint-based and exact-SHA bound, cheap gates run before expensive gates, and release/native certification remains exhaustive.

Release-specific implementation evidence remains under the release/QA evidence structure.

---

## 7. 10/10 Continuity Audit Rule

A reconciliation may be described as **10/10** only after all ten dimensions in `governance/CONTINUITY-IMPLEMENTATION-CONTRACT.md` pass: approval coverage, Stable inheritance, canonical ownership, release placement, executable traceability, defect continuity, security/rights/role continuity, performance/reliability continuity, delivery/platform truth, and future-loss prevention.

A score below 10 must be improved before the reconciliation itself is called 10/10. This does not mean every future roadmap capability is already implemented.

---

## 8. Permanent Review Question

Before every new build plan or substantial DE.PULSE discussion, ask:

> **Is this already approved, inherited/certified, partially covered, a refinement, conflicting, genuinely new, or still an unresolved defect — where is the canonical evidence, what release disposition applies, what checkpoint/qualification evidence is required, and how will packaged-runtime proof be obtained?**
