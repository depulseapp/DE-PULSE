# DE.PULSE — Material Decision Log

**Policy:** Append-only history of material approved, superseded, rejected, deferred, or SHADOW decisions.

Current canonical wording lives in:
- `governance/APPROVED-SCOPE.md`
- `governance/CONTINUITY-IMPLEMENTATION-CONTRACT.md`
- `governance/ADAPTIVE-OPERATING-CONTRACT.md`
- `governance/ADAPTIVE-DATA-RELIABILITY-CONTRACT.md`
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

## DEC-2026-08-14-003 — Institutional Holdings / 13F Intelligence

**Status:** APPROVED  
**Date:** 2026-08-14  
**Placement:** v18/v19 evidence/data foundation; v20 adaptive institutional intelligence

### Decision
Adopt **Institutional Holdings / 13F Intelligence** as a permanent Smart-Money / Institutional Intelligence capability.

Direct SEC EDGAR is the canonical filing source. DE.PULSE will preserve report-period time separately from public filing/acceptance time and must never treat quarterly 13F disclosures as live ownership or leak filings into historical analysis before they became public.

The capability includes normalized holdings/reconciliation, amendments/restatements, truthful quarter-over-quarter position states, manager/institutional Behavioral Fingerprints, disclosed-position conviction/persistence, consensus/crowding, sector/thematic rotation, outcome learning, and adaptive correlation with insider/congressional activity, ASBI, Rapid Move, Opportunity Radar, earnings/SEC/news catalysts, market/sector regime, price/volume/relative strength, and options context where useful.

13F limitations—including filing lag, absence of short positions, potentially omitted small positions, confidential treatment, and incomplete representation of a manager's total real-time portfolio—must remain explicit in data semantics, UI, backtesting and adaptive learning.

Production influence remains:

**SHADOW → VALIDATED → APPROVED → PRODUCTION**

Full scope: `governance/APPROVED-SCOPE.md` INT-006 through INT-016 and `governance/ROADMAP.md` Institutional Holdings / 13F Intelligence placement.

---

## DEC-2026-08-14-004 — 10/10 Two-Sided Directional Thesis & Trade Plan Intelligence

**Status:** APPROVED  
**Date:** 2026-08-14  
**Affects:** Adaptive Roadmap, Adaptive Build Plan, Adaptive Build Process, Adaptive Delivery Process, Day/Swing/Long, ASBI, Research, Opportunity Radar, Decision Queue, Historical Validation  
**Placement:** v18/v19 evidence/data foundation; v20 major adaptive implementation/validation

### Decision
Adopt **Two-Sided Directional Thesis & Trade Plan Intelligence (TDTI)** as a permanent 10/10 AI/LLM-style decision-intelligence contract.

For each eligible ticker/horizon, DE.PULSE evaluates Long and Short as competing theses from the same canonical point-in-time evidence rather than mirroring a long plan into a short plan. Valid outputs include LONG THESIS, SHORT THESIS, BOTH DEVELOPING, and NO RELIABLE EDGE.

The system must distinguish Direction Probability, Thesis Strength, Confidence, Opportunity Quality / Decision Utility, and Readiness. Long planning includes Entry, Trim/Target, structural Invalidation and R:R. Short planning includes Short Entry, Cover/Trim, downside targets, structural Short Invalidation and R:R, plus short-specific squeeze/extension/catalyst/liquidity/crowding/borrow/shortability/SSR context only where trustworthy and available.

TDTI integrates ASBI behavioral states/path transitions, probability momentum, regime/sector/catalyst context, evidence independence/contradictions, cross-horizon reasoning, immutable point-in-time thesis/trade-plan lineage, side-aware outcome measurement, Champion/Challenger learning and ABSTAIN / NO RELIABLE EDGE.

The product should feel AI/LLM-style through concise synthesis of `WHY / CONFIRMS / INVALIDATES / WHAT CHANGES THE VIEW / WHAT TO WATCH NEXT`, while LLMs remain grounded research/synthesis components rather than canonical market-truth owners. No unsupported level invention, no silent formula/model self-modification, and no duplicate Long/Short data or execution engines.

TDTI does not create a Short Desk by default and does not change the permanent No Execution Boundary. Trade plans remain research/decision-support outputs only.

Production adaptive influence remains:

**SHADOW → VALIDATED → APPROVED → PRODUCTION**

Full scope: `governance/APPROVED-SCOPE.md` TDTI-001 through TDTI-023, `governance/ROADMAP.md` Two-Sided Directional Thesis & Trade Plan Intelligence placement, and `governance/ADAPTIVE-OPERATING-CONTRACT.md` section 20.

---

## DEC-2026-08-14-005 — 10/10 Adaptive Data Reliability & Graceful Degradation Intelligence

**Status:** APPROVED  
**Date:** 2026-08-14  
**Affects:** Adaptive Roadmap, Adaptive Build Plan, Adaptive Build Process, Adaptive Delivery Process, Provider Router, Shared Symbol Intelligence, SQLite/PostgreSQL, Day/Swing/Long, ASBI, TDTI, Research, Opportunity Radar, Decision Queue, Maintenance/Data Engine  
**Placement:** dependency-compatible v18 reliability work; v18.3 persistent/shared-state foundation; mandatory v18.5 reliability closure; v19 professional hardening; v20 governed adaptive optimization

### Decision
Adopt **Adaptive Data Reliability & Graceful Degradation Intelligence (ADR-GDI)** as a permanent 10/10 reliability architecture and operating contract.

The system must stop treating `DATA DEGRADED` as a broad generic state. Reliability is capability-specific, freshness-aware, consumer/dependency-aware and impact-aware. Optional/stale context must not contaminate unrelated decision surfaces; missing/unreliable required evidence must cause scoped degradation or ABSTAIN / NO RELIABLE EDGE rather than false confidence.

ADR-GDI includes:
- capability-level health;
- dataset/horizon/session freshness SLOs;
- canonical degradation reason codes;
- consumer dependency graphs and blast-radius intelligence;
- `NORMAL → PRESSURE → PROTECTED → DEGRADED → RECOVERING → HEALTHY` runtime state;
- workload prioritization/backpressure/load shedding;
- persistent warm canonical SQLite/PostgreSQL state;
- fetch-once/calculate-once, single-flight/coalescing and material-change propagation;
- Provider Router circuit/cooldown/fallback discipline and calls avoided;
- DB/query/pool/runtime observability so PostgreSQL does not become the new bottleneck;
- graceful confidence reduction/UNKNOWN/ABSTAIN semantics;
- degradation/recovery event ledger;
- adaptive provider/recovery/workload learning under SHADOW → VALIDATED → APPROVED → PRODUCTION governance;
- concise impact-aware USER messaging with deeper Maintenance/Data Engine diagnostics;
- G0–G16 failure-injection, load, restart, DB, queue, fallback, recovery and actual packaged-runtime acceptance.

The success criterion is **not zero truthful degradation events**. It is that DE.PULSE rarely degrades because of its own architecture, isolates unavoidable failures to the smallest truthful blast radius, protects high-value decision-critical evidence, explains impact clearly, recovers safely with hysteresis, and learns from real reliability outcomes without silent production self-modification.

PostgreSQL is one component of the solution, not the solution by itself.

Full contract: `governance/ADAPTIVE-DATA-RELIABILITY-CONTRACT.md` and roadmap placement in `governance/ROADMAP.md`.

---

## DEC-2026-08-14-006 — 10/10 Adaptive Opportunity Discovery & Recommendations

**Status:** APPROVED  
**Date:** 2026-08-14  
**Affects:** Adaptive Roadmap, Adaptive Build Plan, Adaptive Build Process, Adaptive Delivery Process, Global Symbol Registry, Reliable Actionable Universe, Shared Symbol Intelligence, Opportunity Radar, My Market, ASBI, TDTI, ADR-GDI, Dashboard, Decision Queue, Research, Historical Validation  
**Placement:** dependency-compatible v18/v19 productization and evidence foundation; v19 hardening; v20 mature adaptive ranking/personalization

### Decision
Adopt **Adaptive Opportunity Discovery & Recommendations (AODR)** as a permanent 10/10 AI/LLM-style opportunity-prioritization capability.

AODR converts market intelligence DE.PULSE already understands into two user-facing opportunity groups:
- **My Market Opportunities** — strongest material opportunities among symbols already in that user's watchlists/My Market;
- **Global Opportunities** — strongest material opportunities from the eligible broader market universe that are not already in that user's My Market.

AODR does not create another scanner, symbol registry, market-data pipeline, chart engine, ASBI/TDTI clone, or independent recommendation-truth silo. It reuses Global Symbol Registry eligibility, Reliable Actionable Universe quality, Shared Symbol Intelligence, Opportunity Radar broad observation/PROMOTE/DEMOTE, ASBI behavior intelligence, TDTI Long/Short opportunity quality/readiness, and ADR-GDI reliability truth.

Underlying market evidence remains canonical and shared where semantically valid. User preferences may affect relevance/ranking/presentation, but must not fabricate different market truth for different users. The permanent model is:

**shared canonical market truth → user-specific relevance/ranking/presentation**

Ranking is not a highest-score leaderboard. It must consider applicable Opportunity Quality, Readiness, ASBI state/path/probability momentum, TDTI Long/Short thesis, expected magnitude/time-to-resolution, regime, relative strength, liquidity, catalysts, evidence independence/contradiction, freshness/data sufficiency, extension/chase and squeeze risk, historical usefulness, materiality and opportunity cost.

Candidate lifecycle supports **NOW / WATCH / PASS / ABSTAIN / NO RELIABLE EDGE**. AODR must be allowed to show no strong opportunities rather than manufacture recommendations. Strong direction with poor entry location, poor R:R, excessive extension or degraded required evidence must be demoted even when the underlying thesis remains strong.

AODR preserves staged scalable processing:

**broad low-cost observation → PROMOTE → deeper shared analysis → rank/surface → DEMOTE**

It must not run expensive deep analysis across the entire Global Symbol Registry continuously. Correlation/diversity logic should avoid presenting many copies of the same sector/theme/catalyst as independent opportunities while remaining a recommendation-quality feature, not portfolio construction.

Normal UX should be concise and AI/LLM-style: why now, Long/Short/horizon, readiness/quality, confirms, invalidates, key contradiction/risk, extension/chase state and what to watch next. LLM output must remain grounded in canonical structured evidence and must not invent prices, levels, catalysts, probabilities or recommendations.

Material surfaced recommendations must be recorded point-in-time before outcomes are known so DE.PULSE can measure recommendation usefulness, misses, redundancy, staleness, extension/chase errors, diversity value, degradation effects and decision utility. Learned ranking/personalization follows **SHADOW → VALIDATED → APPROVED → PRODUCTION** with Champion/Challenger evaluation and no silent production self-modification.

AODR does not change the permanent No Execution Boundary. Recommendations are research/decision-support prioritization only.

Full scope: `governance/APPROVED-SCOPE.md` AODR-001 through AODR-018 and `governance/ROADMAP.md` Adaptive Opportunity Discovery & Recommendations placement.

---

## DEC-2026-08-14-007 — 10/10 Seven-Day Continuity & Governance-to-Implementation Reconciliation

**Status:** APPROVED  
**Date:** 2026-08-14  
**Affects:** all canonical governance, inherited Stable functionality, v18.2 active implementation, v18.3 mandatory carry-forward work, RBAC/ownership, documentation, native delivery, unresolved defects, G0–G16 traceability  
**Placement:** permanent governance hardening; current v18.2 obligations preserved; specified source-changing remediation enters v18.3; inherited Stable capability remains protected across future releases

### Decision
Adopt `governance/CONTINUITY-IMPLEMENTATION-CONTRACT.md` as the permanent 10/10 continuity and governance-to-implementation closure contract.

The reconciliation was explicitly audited as a quality problem rather than treated as a naming exercise. A simple “we remembered the features” pass was insufficient because it did not guarantee executable traceability, Stable inheritance, unresolved-defect continuity, ownership/RBAC/documentation continuity, platform truth or safe synchronization between newer canonical governance and unique active implementation.

The permanent closure chain is:

**DISCUSSED → EXPLICITLY APPROVED → CANONICALIZED → ROADMAP-PLACED → G1-FROZEN WHEN APPLICABLE → IMPLEMENTED → ENFORCED → EVIDENCED → PACKAGED/RUNTIME-PROVEN → DELIVERED → OUTCOME-MEASURED → LEARNED → ADAPTED**

The contract explicitly carries forward previously approved/certified truth including:
- U.S. Equities Processing Boundary and U.S.-listed actionable semantics;
- U.S.-primary Economic Calendar/Event Intelligence with US_MARKET_CRITICAL / US_CONTEXT / GLOBAL_CONTEXT behavior;
- Community Intelligence / Global Evidence Fusion Hub with rights, provenance, dedupe, corroboration and prompt-injection safeguards;
- Unusual Volume & Volatility Intelligence as Opportunity Radar/Discovery capability rather than another scanner;
- Historical Replay/point-in-time no-lookahead validation;
- exactly 1 active SUPER OWNER + up to 5 OWNERs, with capability-based ADMIN authority;
- Ownership Transfer / Transferable Asset Registry requirements;
- role-aware Documentation and matching AI retrieval authorization;
- required macOS Apple Silicon and Windows x64 native delivery, with hosted Linux/web when hosted architecture is delivered and no Intel-Mac investment unless explicitly re-approved;
- governance-to-implementation disposition and traceability;
- active v18.2 capability-based Administration, Administration composition, Role × Tab × Viewport audit, authoritative build-state ledger, Build Coordinator and canonical naming obligations;
- v18.3 mandatory utility/remediation carry-forward for shared Discovery acquisition, Session Intelligence, Event Intelligence and information-architecture/deep-evidence consolidation;
- unresolved Master Market Symbols layout and Prep `Requires Review` repetition defects, which remain OPEN_DEFECT / RECONCILE until actual closure evidence exists;
- inherited Stable capability non-regression and branch/governance synchronization protections.

`APPROVED` is not `IMPLEMENTED`. Documentation alone cannot close a requirement. At G10/G12/G14/G16 and Major Closure, applicable requirements must be traceable through owner/source/tests/package/runtime/outcomes or explicitly carried to a named future release.

No new top-level gates or duplicate product engines are created. Existing canonical owners and G0–G16 remain authoritative.

The reconciliation itself may be called 10/10 only when all ten continuity quality dimensions in the contract pass. This rating describes governance-continuity quality, not a claim that future ASBI/TDTI/AODR/ADR-GDI roadmap features are already implemented.

---

## DEC-2026-08-14-008 — Adaptive Intelligence & Product Learning Checkpoint (AIPLC)

**Status:** APPROVED / PERMANENT  
**Date:** 2026-08-14  
**Affects:** Adaptive Build Plan, Adaptive Build Process, G0–G16 feedback, all tabs/surfaces, data utility, canonical architecture, UI/UX, providers, performance/reliability, testing, defects, adaptive intelligence and next-build prioritization  
**Placement:** after every meaningful build/checkpoint; aggregated into every Stable G16 Adaptive Retrospective

### Decision
Adopt **Adaptive Intelligence & Product Learning Checkpoint (AIPLC)** as a permanent DE.PULSE operating rule inside the existing Adaptive Build Process and G16 feedback loop. It does **not** add G17 or another top-level release gate.

DE.PULSE must learn from each meaningful build rather than waiting only for Stable closure. The permanent cycle is:

**BUILD → TEST → OBSERVE → TAB × FEATURE × DATA INTELLIGENCE AUDIT → CHALLENGE / ROOT-CAUSE REVIEW → LEARN → IMPROVE → PREVENT → MEASURE → FEED NEXT BUILD**

Every affected tab, feature and material datum must be evaluated against the AI/LLM-style North Star using:

**datum → purpose → canonical owner → consumer → freshness/materiality → independence/correlation → interpretation → decision value → explanation → outcome → learning**

The checkpoint must determine whether evidence should be kept/strengthened, synthesized/correlated, consolidated/reused, demoted to drill-down, suppressed/removed, or explicitly treated as UNKNOWN/ABSTAIN. Raw availability never justifies UI prominence or processing by itself.

Every material defect/challenge should produce both:
1. an immediate fix or truthful disposition;
2. a reusable prevention/learning improvement that reduces recurrence, including cross-product pattern scanning and stronger architecture/test/process protection where appropriate.

AIPLC uses a ten-dimension evidence-based quality score covering AI/LLM Smartness, Data Utility, Canonical Architecture, Decision Utility/User Value, Reliability/Freshness Truth, Performance/Efficiency, UI/UX Intelligence, Learning/Outcome Measurement, Testing/Prevention, and Adaptive Continuity.

**10/10 is a target and quality threshold, not an automatic label.** If an applicable dimension is below 10, improve it immediately when safely in scope and re-audit. If a current build cannot safely close the gap, record the exact residual, evidence, risk and named disposition. Never inflate a rating merely to call the build 10/10.

AIPLC is allowed to influence the next Adaptive Build Plan using validated evidence: strengthen useful providers/intelligence, demote low-value inputs, consolidate redundant code/data/UI, improve workload priority, keep weak learned logic in SHADOW, strengthen root-cause prevention, and simplify user-facing surfaces. It does not authorize silent self-modifying production logic or bypass frozen G1 scope.

Checkpoint evidence should be concise/durable and include build identity, affected surfaces, findings/root causes, changes/prevention, score/residuals, next-build dispositions, metrics to watch and any promotion blocker.

For the currently active v18.2 line, AIPLC is process hardening only: preserve the frozen v18.2 G1, resume qualification from the synchronized active branch, and run the first applicable AIPLC before the next promotion decision. Future source-changing findings outside frozen scope are carried to the named next build unless they are genuine release blockers requiring governed correction/requalification.

Full wording: `governance/ADAPTIVE-OPERATING-CONTRACT.md` section 15.

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

## DEC-2026-08-28-001 — Canonical Adaptive Documentation Rebaseline

**Status:** APPROVED / IMPLEMENTED IN GOVERNANCE

Use one narrative authority per adaptive concern: `governance/ADAPTIVE-OPERATING-CONTRACT.md`, `governance/ROADMAP.md`, and the canonical adaptive Build Plan, Build Process and Delivery Process under `adaptive-governance/`.

`CURRENT_ADAPTIVE_*` and `ADAPTIVE_ROADMAP.md` remain compatibility/status projections only. Exact progress and the one next action belong in machine current state, the active closure ledger and `handoff/CURRENT.md`. Preserve historical evidence; do not treat planned audit gaps as implemented.

Full record: `governance/DECISION-2026-08-28-ADAPTIVE-DOCUMENTATION-REBASELINE.md`.

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
