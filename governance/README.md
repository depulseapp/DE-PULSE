# DE.PULSE Governance — Canonical Entry Point

**Status:** ACTIVE / AUTHORITATIVE FOR APPROVED INTENT  
**Repository:** `depulseapp/DE-PULSE`  
**Purpose:** Prevent approved product, architecture, roadmap, build-process, release-gate, and delivery decisions from being lost across chats, handoffs, or releases.

---

## 1. Source-of-Truth Hierarchy

Use the following hierarchy whenever DE.PULSE scope or a prior decision is questioned:

1. **Actual source / packaged release evidence** — truth about what is actually implemented and delivered.
2. **`governance/APPROVED-SCOPE.md`** — canonical truth about what has been approved as product/roadmap scope.
3. **`governance/ADAPTIVE-OPERATING-CONTRACT.md`** — canonical permanent engineering/build/release/delivery contract.
4. **`governance/ROADMAP.md`** — canonical placement and sequencing of approved future work.
5. **`governance/DECISION-LOG.md`** — append-only history of approved, superseded, rejected, or deferred material decisions.
6. **`release/<version>/G1-IMMUTABLE-SCOPE.md`** — immutable scope snapshot for one specific release after G1.
7. **G16 handoff / release evidence** — what actually happened in that release and what carries forward.
8. Chat memory or recollection — continuity aid only; never override the canonical GitHub records above.

A governance document never proves a feature exists in code. A release/code artifact never silently changes approved product intent.

---

## 2. Required Discussion Workflow

Before treating a proposed DE.PULSE idea as new, use:

**LOOKUP → COMPARE → CLASSIFY → DECIDE → UPDATE**

### LOOKUP
Read/search at minimum:
- `governance/APPROVED-SCOPE.md`
- `governance/ADAPTIVE-OPERATING-CONTRACT.md`
- `governance/ROADMAP.md`
- `governance/DECISION-LOG.md`

When release-specific, also inspect the current `release/<version>/G1-IMMUTABLE-SCOPE.md` and latest G16/handoff.

### COMPARE
Compare the proposal semantically, not only by exact wording. Check aliases, renamed features, overlapping responsibilities, and existing canonical owners.

### CLASSIFY
Every proposal must be classified as one of:

- **ALREADY APPROVED** — existing scope covers it.
- **PARTIALLY COVERED** — existing scope covers part; identify the delta.
- **REFINEMENT / HARDENING** — same approved purpose, stronger implementation or acceptance.
- **CONFLICT / SUPERSESSION** — contradicts or replaces an earlier decision.
- **NEW SCOPE** — materially new product/architecture behavior.
- **REJECT / NO CHANGE** — not useful or violates a permanent boundary.

### DECIDE
Do not silently add material scope. User approval is required for genuinely new material product decisions or supersession of an approved contract.

### UPDATE
After approval:
1. update the relevant canonical governance document;
2. append a Decision Log entry;
3. place the work in the Adaptive Roadmap/Build Plan;
4. when a release reaches G1, snapshot the applicable scope into its immutable G1 file.

---

## 3. Do Not Use Handoffs as the Only Memory

Handoffs are release continuation records, not the master product contract. They may summarize only what was relevant to that release.

Permanent decisions belong here in `governance/` and are referenced by future handoffs rather than repeatedly copied in full.

---

## 4. Release Relationship

The permanent model is:

`Governance approved intent`
→ `Adaptive Roadmap / Build Plan`
→ `G0 baseline`
→ `G1 immutable release scope snapshot`
→ `G2–G15 engineering / certification / delivery`
→ `G16 retrospective + handoff`
→ `governance reconciliation when a permanent decision changed`

No release may silently drop an approved permanent contract merely because it was omitted from a handoff summary.

---

## 5. Efficiency Rule

Do **not** rewrite all governance documents for every small patch.

Update governance only when:
- approved product scope changes;
- a permanent architecture/build/release/delivery contract changes;
- roadmap placement materially changes;
- a decision is superseded/rejected;
- G16 discovers a durable rule that must carry forward.

Release-specific implementation evidence remains under the release/QA evidence structure.

---

## 6. Permanent Review Question

Before every new build plan or substantial DE.PULSE discussion, ask:

> **Is this already approved, partially covered, a refinement, conflicting, or genuinely new — and where is the canonical GitHub evidence?**
