# DE.PULSE — Material Decision Log

**Policy:** Append-only history of material approved, superseded, rejected, deferred, or SHADOW decisions.

Current canonical wording lives in:
- `governance/APPROVED-SCOPE.md`
- `governance/ADAPTIVE-OPERATING-CONTRACT.md`
- `governance/ROADMAP.md`

Do not delete old decisions when direction changes. Mark them **SUPERSEDED** and reference the replacing decision.

---

## DEC-2026-08-14-001 — GitHub Canonical Governance Layer

**Status:** APPROVED  
**Date:** 2026-08-14

### Decision
Use the DE.PULSE GitHub repository as the durable canonical source of truth for approved product scope, permanent operating contracts, roadmap placement, and material decision history.

Before treating a proposed feature/change as new, follow:

**LOOKUP → COMPARE → CLASSIFY → DECIDE → UPDATE**

Classify as one of:
- ALREADY APPROVED
- PARTIALLY COVERED
- REFINEMENT / HARDENING
- CONFLICT / SUPERSESSION
- NEW SCOPE
- REJECT / NO CHANGE

Release-specific `G1-IMMUTABLE-SCOPE.md` remains the immutable build snapshot for that release. G16 remains the retrospective/handoff and must reconcile durable changes back into governance when appropriate.

---

## DEC-2026-08-14-002 — Adaptive Stock Behavior Intelligence

**Status:** APPROVED  
**Date:** 2026-08-14  
**Placement:** v18/v19 evidence preparation; v20 major implementation/validation

### Decision
Adopt Adaptive Stock Behavior Intelligence (ASBI) as a permanent major adaptive-intelligence workstream.

ASBI uses behavioral states, transitions, competing paths, multi-horizon probabilities, confidence/uncertainty, expected outcome distributions, historical analogues, Behavioral Fingerprints, catalyst/regime context, ABSTAIN behavior, immutable forecast/outcome logging, and Champion/Challenger evaluation.

Production influence remains:

**SHADOW → VALIDATED → APPROVED → PRODUCTION**

Full scope: `governance/APPROVED-SCOPE.md` ASBI-001 through ASBI-022.

---

## DEC-2026-08-13-001 — Four-Layer Adaptive Operating Model

**Status:** APPROVED / CARRIED FORWARD

Keep four connected governing layers:
1. Adaptive Roadmap — where/why;
2. Adaptive Build Plan — what next;
3. Adaptive Build Process — how engineered/qualified;
4. Adaptive Delivery Process — how certified work becomes trustworthy Stable delivery.

---

## DEC-2026-08-13-002 — Canonical G0–G16 Release Model

**Status:** APPROVED / CARRIED FORWARD

Use only G0–G16 as top-level release gates. New responsibilities belong inside the existing gate model.

---

## DEC-2026-08-13-003 — Dependency-Aware Bounded Parallelism

**Status:** APPROVED / CARRIED FORWARD

Use CI/CD-style bounded parallel execution for dependency-safe independent work/tests while preserving governance order and source-fingerprint validity of evidence.

---

## DEC-2026-08-13-004 — Smart Intelligent Provider Router v2

**Status:** APPROVED / COMMITTED

Maintain one canonical Provider Router with dataset/capability/entitlement-aware routing, Preferred vs Serving semantics, deterministic cooldown/circuits, provider budgets/headroom, source disagreement handling, calls avoided, and provider usefulness telemetry.

---

## DEC-2026-08-13-005 — TradeInsight SHADOW / SECONDARY Role

**Status:** EXPERIMENTAL / SHADOW / COMMITTED

TradeInsight may enrich DE.PULSE only through the canonical Provider Router and starts as SHADOW/SECONDARY. Its long-term provider role is reassessed under v19 Professional Data Infrastructure based on measured value, quality, reliability, cost, and rights.

---

## DEC-2026-08-13-006 — Major Roadmap Sequence

**Status:** APPROVED

**v18 Secure Multi-User Platform + Smart Provider Intelligence**  
→ mandatory v18 Major Closure  
→ **v19 Professional Data Infrastructure**  
→ mandatory v19 Major Closure  
→ **v20 Adaptive Intelligence & Decision Research**

Exact compatible minor placement may adapt when dependency/risk/evidence justify it.

---

## DEC-2026-08-13-007 — Permanent No Execution Boundary

**Status:** APPROVED / PERMANENT

DE.PULSE remains research/intelligence/decision support. No order execution, broker routing, automated/semi-automated trading, paper/live trading product, OMS/blotter, portfolio/P&L, or journal scope unless explicitly superseded by a future approved material product decision.

---

## DEC-2026-08-13-008 — Autonomous v18 Continuation Through v18.5

**Status:** APPROVED

Routine implementation, testing, certification, packaging, cleanup, and build operations through v18.5 may proceed autonomously. Stop only for unavoidable credentials/secrets, new financial commitments, legal/licensing/data-rights decisions, irreversible/high-impact external actions requiring approval, or genuinely new material product decisions.

---

# Append Template

```text
## DEC-YYYY-MM-DD-NNN — Title
Status: APPROVED | SUPERSEDED | DEFERRED | REJECTED | EXPERIMENTAL/SHADOW
Date: YYYY-MM-DD
Affects: <areas>

Decision:
<what was decided>

Why:
<why it matters>

Supersedes / Superseded by:
<decision id if applicable>
```
