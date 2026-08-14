# DE.PULSE — Adaptive Build Plan

**Status:** ACTIVE / GOVERNED  
**Authority:** product sequencing comes from `governance/ROADMAP.md`; permanent operating rules come from `governance/ADAPTIVE-OPERATING-CONTRACT.md`. This file defines how an active release is planned efficiently inside those contracts.

## 1. Planning North Star

Every release plan follows:

**approved scope → impact map → reuse → smallest safe work packages → bounded execution → evidence → AIPLC learning → next action.**

The plan must maximize assurance **without** maximizing workload.

Permanent engineering order:

**REUSE → CONSOLIDATE → REFACTOR → DELETE/REPLACE → ADD**

Permanent gate model: **G0–G16 only.** No release plan may introduce G17+.

---

## 2. Durable resume state before G3 exit

Every release must define:
- exact incoming Stable tag/commit;
- one active development/release branch;
- canonical release identity/runtime profile;
- build checkpoint and release-evidence checkpoint locations;
- candidate/source fingerprint strategy;
- authoritative next step/blocker;
- G16 archive/cleanup and next-release seed.

Meaningful implementation work must be committed before it is considered resumable.

A checkpoint label never outranks actual GitHub branch HEAD, source fingerprint, workflow/artifact evidence or immutable RC/package identity.

---

## 3. Impact-First Planning

Before any expensive implementation/qualification, create an impact map:

**Git diff / requested change → canonical owner → dependency blast radius → affected tabs/features/data/roles/providers/runtime paths → reusable evidence → smallest required work/test set.**

Classify each responsibility as:
- **FRESH_REQUIRED** — affected and must be rerun/reviewed;
- **INHERITABLE** — unchanged and equivalent evidence may be reused;
- **SENTINEL_REQUIRED** — unchanged area needs a small dependency/non-regression probe;
- **NOT_APPLICABLE** — outside release impact.

Evidence inheritance is allowed only when relevant:
- source/artifact fingerprint is equivalent;
- canonical owner/dependency contract is unchanged;
- test definition/input semantics remain applicable;
- role/security/data-rights assumptions remain equivalent.

This classification must be durable enough for G10 coverage reconciliation.

---

## 4. Three-Depth AIPLC Plan

### Level 1 — every meaningful build/checkpoint
Run **Delta AIPLC** only for changed/affected areas plus dependency sentinels.

For each affected tab/feature/datum use:

**datum → purpose → canonical owner → consumer → freshness/materiality → independence/correlation → interpretation → decision value → explanation → outcome → learning.**

Every material challenge produces:
1. immediate fix/truthful disposition;
2. reusable prevention/cross-product pattern scan.

### Level 2 — G10
Run **Full Coverage Reconciliation**. Every required product/process responsibility must be either freshly evidenced or explicitly inherited from equivalent trustworthy evidence.

### Level 3 — G16 / Major Closure
Run the deep all-system Adaptive Retrospective, including architecture, source quality, data utility, UI/UX, reliability, performance, provider usefulness, adaptive intelligence, process failures and recurring defects.

Mechanically identical reruns with no new evidence may record `NO NEW LEARNING / EVIDENCE EQUIVALENT` instead of creating another heavy report.

---

## 5. Build Coordinator & bounded CI plan

Every release has one logical Build Coordinator and one authoritative G0–G16 dependency graph.

Plan independent lanes only where safe. Bounded parallelism must account for:
- runner/CPU/memory;
- browser instances;
- provider/API limits;
- database/storage contention;
- AI/LLM calls;
- artifact/publishing mutations.

One canonical owner exists for each expensive test/gate responsibility.

CI failures are classified before candidate health changes:
`PRODUCT_FAIL`, `GATE_TEST_FAIL`, `CI_HARNESS_FAIL`, `INFRA_FAIL`, `EXPECTED_NOOP`, `SUPERSEDED`.

Metadata/checkpoint-only paths must be isolated so they do not trigger unnecessary certification/native work.

Mutating workflows must be idempotent; a legitimate no-change operation is not a product failure.

---

## 6. Test/load efficiency plan

Plan the cheapest trustworthy evidence first:
- **G5:** changed-area FAST tests;
- **G6:** affected integration/MEDIUM tests;
- **G7/G8/G9:** bounded independent affected lanes;
- **G10:** complete coverage reconciliation;
- **G12:** one full certification on the immutable RC.

Provider/data testing:
- prefer deterministic fixtures, historical replay and canonical cached evidence when live behavior is not the subject under test;
- use bounded live-provider smoke only where actual live routing/fallback/freshness must be proven;
- share equivalent acquired evidence across test lanes rather than refetching it.

AI/LLM testing:
- prefer canonical evidence packages/fingerprints for grounding/regression;
- use small bounded true-model samples when actual model-runtime behavior must be certified;
- do not rerun materially identical synthesis per user/symbol/test when reusable evidence is equivalent.

Native:
- macOS Apple Silicon and Windows x64 evidence are independent lanes;
- preserve an unchanged platform PASS when exact RC/package identity remains unchanged.

---

## 7. Role-Aware Planning

Any release that changes tabs, cards, shell controls, administrative operations or role-sensitive data must define the role/capability composition before G3 exit.

Cover:
- SUPER_OWNER / OWNER;
- full-capability ADMIN where explicitly delegated;
- limited/delegated ADMIN;
- USER;
- DEMO;
- visible controls and backend-authorized actions;
- fields withheld server-side;
- information hierarchy/reflow;
- direct route/API denial;
- required desktop/tablet/narrow-browser coverage.

ADMIN authority is capability-based; USER/DEMO receive no implementation machinery.

Frontend visibility and backend authorization are one capability boundary.

Canonical details: `adaptive-governance/ROLE_AWARE_UI_COMPOSITION_CONTRACT.md` and `ROLE_AWARE_SESSION_SECURITY_CONTRACT.md`.

---

## 8. Functionality Utility / Minimal-Surface Planning

For every new/materially changed tab, card, engine, scanner, detector, job, watcher, scheduler, provider path, dataset, derived metric/model, alert, persistence field or administrative operation, record:
- purpose/consumer;
- canonical owner;
- reused implementation/data;
- shared/coalesced acquisition/computation/persistence opportunities;
- overlap/consolidation decision;
- required correlations;
- freshness/materiality/retention/rights/degraded behavior;
- performance cost;
- UI disposition;
- role/backend authorization;
- obsolete implementation/surface/job to retire where safe.

Default:

**one canonical intelligence owner → one deep-evidence home → concise contextual reuse elsewhere.**

Maintain `functionality_utility_registry.json`; use `functionality_utility_checkpoint_gate.py` before G10.

---

## 9. Governance-to-Implementation Closure Planning

Every applicable permanent requirement receives one disposition:
- `CURRENT_RELEASE_BLOCKER`;
- `CURRENT_RELEASE_PROCESS_HARDENING`;
- `NEXT_RELEASE_MANDATORY_ENTRY`;
- `FUTURE_STRATEGIC`.

A documented rule with no owner/evidence/completion plan is a planning defect.

For each applicable item record:
- implementation owner/target;
- dependencies;
- evidence gate/checkpoint;
- regression coverage;
- delivery impact;
- exact completion criteria.

Canonical details: `adaptive-governance/GOVERNANCE_IMPLEMENTATION_CLOSURE_CONTRACT.md`.

---

## 10. Shared Symbol Intelligence Planning

Any release affecting symbol demand, provider acquisition, Scanner/Radar, preparation/event work, research, user workspaces, adaptive synthesis, persistence or hosted scale must prove:

**shared canonical processing → material-change reuse → authorized personal composition.**

Plan:
- shared demand union;
- canonical processing keys;
- one owner for acquisition/normalization/state/calculation/reusable intelligence;
- cache/freshness/single-flight ownership;
- material-change invalidation graph;
- dynamic attention/provider budgets/backpressure/fairness;
- private-context/rights/entitlement partitions;
- multi-user overlapping/non-overlapping load scenarios;
- efficiency metrics.

Equivalent overlapping demand should scale primarily with **unique canonical demand**, not `users × symbols`.

Canonical details: `adaptive-governance/SHARED_SYMBOL_INTELLIGENCE_PROCESSING_CONTRACT.md`.

---

## 11. v18.2 required closure

Frozen v18.2 G1 remains authoritative and must not be expanded.

Before v18.2 promotion, the plan must close/evidence applicable current-release obligations including:
- capability-based ADMIN authorization shared by UI/backend;
- justified capability-scoped Administration composition;
- complete role × tab × viewport coverage using fresh evidence plus valid inheritance;
- authoritative Build State Ledger/checkpoint reconciliation from actual GitHub truth;
- one Build Coordinator with overlapping release workflows consolidated/subordinated;
- canonical active naming for gates/workflows/jobs/artifacts/checkpoints/capabilities;
- first v18.2 AIPLC before the next promotion decision.

AIPLC source-changing findings outside frozen v18.2 scope normally become named next-release entries unless they expose a genuine correctness/security/reliability release blocker.

---

## 12. v18.3 mandatory entry

All v18.2 audit source-changing utility/consolidation work remains mandatory input to v18.3 G1–G3, including the authoritative items in `functionality_utility_remediation.json`.

Key direction:
- shared/coalesced Scanner/Radar acquisition;
- Session Intelligence Coordinator for Pre-Market/Market Open Prep;
- Event Intelligence lifecycle ownership;
- deep-evidence/UI consolidation;
- shared hosted demand/canonical persistence/recovery;
- PostgreSQL must not create duplicate per-user market-wide computation.

These items cannot silently disappear.

---

## 13. G0–G16 planning map

- **G0:** exact baseline/checkpoint/open defects.
- **G1:** frozen release scope + applicable inherited governance.
- **G2:** architecture/data utility/canonical ownership/impact map.
- **G3:** dependency/readiness/test/decomposition plan.
- **G4:** development exit.
- **G5:** FAST affected-area qualification.
- **G6:** affected integration/MEDIUM qualification.
- **G7:** data/security/adaptive intelligence.
- **G8:** performance/capacity/stability.
- **G9:** cross-module/UI/UX.
- **G10:** full coverage reconciliation/pre-freeze.
- **G11:** immutable RC.
- **G12:** full RC certification.
- **G13:** native packaging/provenance.
- **G14:** actual artifact runtime audit.
- **G15:** release assurance/promotion.
- **G16:** deep adaptive retrospective/handoff.

No G17+.

---

## 14. Build Plan success criteria

The plan is 10/10 only when it can answer, with durable evidence:
1. What changed?
2. What depends on it?
3. What can be reused?
4. What must rerun?
5. Who owns each responsibility?
6. What is the smallest safe workload?
7. What learned/prevention action came from the build?
8. What remains and where is it explicitly placed?
9. Is G1 protected?
10. Can another conversation/runner resume from GitHub without guessing?
