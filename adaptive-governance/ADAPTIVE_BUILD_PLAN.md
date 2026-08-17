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

## 12A. v18.5.1 escaped-defect and implementation-reconciliation plan

v18.5.1 combines repository-archetype closure with mandatory recovery of escaped v18.5 behavior. The recovery work is not optional cleanup and may not inherit a previous PASS for an affected surface.

### Required recovery lanes

| Lane | Scope | Required outputs |
|---|---|---|
| A — release/profile truth | COPY-18.5.1-001 and VERSION-18.5.1-002 | Version-aware/release-neutral copy, explicit migration-identifier classification, fresh-install/upgrade/TEST-profile matrix and packaged proof. |
| B — symbol/watchlist semantics | SYMBOL-18.5.1-001 and SYMBOL-18.5.1-002 across Day/Swing/Long | Explicit row-X contract, final-membership behavior, state/persistence tests and current-desk visual-state evidence. |
| C — interaction continuity | NAV-18.5.1-001 | Root-cause removal for live/SSE rerender jumps, scroll/focus preservation tests and dwell/reload retest. |
| D — Research requirement recovery | Reopened v15.1.0 approved items 17–19 | Responsive layout, truthful labels/freshness/actions/disabled states, missing-evidence recovery and requirement-linked tests. |
| E — full promise reconciliation | 48 inherited requirements, 20 v17 items, all v18 workstreams/release entries, 13 orphaned utility remediations, conversational commitments and accepted defects | Durable implementation ledger, unexplained-gap list, owners, decisions, regression IDs and actual-package evidence. |
| F — confirmed orphan recovery | IMPL-18-TRADEINSIGHT-001 | Canonical Smart Router SHADOW adapter, configuration/entitlement, rights/provenance controls, tests and actual-package evidence. |
| G — acquisition/coordination remediation | IMPL-18-UTILITY-001 and IMPL-18-UTILITY-002 | Shared snapshot broker/cache, one Session Intelligence Coordinator, coalescing/provider-budget/restart/catch-up tests and package proof. |
| H — surface/route remediation | IMPL-18-UTILITY-003 and IMPL-18-UTILITY-004 | Minimal Discovery hierarchy, canonical evidence homes, safe redirects/deep links and responsive package proof. |
| I — documentation governance | IMPL-18-DOC-001 | Server-authoritative audience policy, role-composed UI/direct paths, Documentation Impact Manifest and cross-role security/package tests. |
| J — dependency readiness | IMPL-17-DEPS-001 | Canonical dependency/readiness registry, User Action Required records, gate binding and role-safe operational proof. |
| K — remaining utility carry-forward | Dashboard, Market Intelligence, three Desks, Catalyst Watch and Maintenance records | Fresh design/behavior/performance/package disposition; no registry-only or target-release-only closure. |

Each ledger row must contain:

`origin ID/source → user-visible promise → current source owner → current observed behavior → disposition → defect/implementation owner → fix or approved placement → regression test ID → browser evidence → macOS package evidence → Windows package evidence → closure approver`.

Allowed dispositions are only `FRESH_PASS`, `REOPENED`, `NOT_IMPLEMENTED`, `INTENTIONALLY_SUPERSEDED`, `NOT_APPLICABLE` and explicitly `ROADMAP_PLACED_FUTURE`. An inherited or historical PASS is not a current disposition.

### Recurrence and duplicate handling

1. Search prior scope, roadmap, release evidence and issue history before creating an item.
2. A matching report is linked to its original requirement/defect; it is not silently discarded as a duplicate.
3. If the behavior reproduces, reopen the item, increment a recurrence counter and invalidate affected/dependent evidence.
4. A defect reported again after claimed closure is a release blocker and requires explicit G16 root-cause plus prevention evidence.
5. Closure always follows **Find → Fix → Regression Test → Retest**. Source inspection, a unit test or a screenshot alone is insufficient.

### Gate binding

- **G0:** freeze the exact v18.5.0 Stable/package baseline and the complete inherited + v17 + v18 + 13-item utility remediation + conversational commitment + known-defect inventory.
- **G1:** freeze repository-archetype work together with the six escaped-defect groups, all seven confirmed implementation-gap groups, the 13-item utility carry-forward and full v17/v18 reconciliation; exclude unrelated redesign.
- **G2–G3:** record canonical owner, blast radius, migration/state/API/UI impact, test IDs and package matrix for every row.
- **G4:** implement fixes atomically with their tests; do not hide semantic changes inside structural moves.
- **G5–G6:** pass source/unit and functional/integration coverage, including persistence, membership combinations and live-refresh behavior.
- **G9:** run direct UI/UX acceptance across Day/Swing/Long and Research at desktop, tablet and narrow widths, including scroll/focus continuity.
- **G10:** reconcile all 48 inherited requirements, all 20 v17 items, every v18 workstream/release entry, all 13 utility remediations, every accepted conversational commitment and defect to fresh current evidence; no blank, orphaned or unexplained rows.
- **G11–G12:** freeze one immutable RC only after reconciliation is clean, then run full certification.
- **G13–G14:** build and directly exercise macOS Apple Silicon and Windows x64 artifacts; affected workflows cannot inherit v18.5 package evidence.
- **G15:** block promotion for any reopened, not-implemented or unexplained item unless an intentional supersession/placement is explicitly approved.
- **G16:** record root cause, why previous gates missed it, the new regression/prevention control and the next-release handoff.

---

## 12B. Requirement-Controlled Build Slicing

This section supersedes plan-only slicing. A build slice is executable only when it is derived from the reconciliation ledger and carries its parent requirement IDs.

### Mandatory work-packet schema

Every work packet must declare:

- immutable requirement/defect IDs and exact origins;
- current user-visible promise;
- target release/slice and implementation owner;
- source modules, API/state/persistence/UI/package impacts;
- dependency and evidence-invalidation graph;
- security/role/data-rights implications;
- performance/provider/storage budget;
- regression test IDs and behavioral scenarios;
- browser/API/runtime evidence requirements;
- macOS and Windows package evidence requirements;
- completion status and unresolved blockers.

A work packet with no parent IDs cannot enter G4. A parent ID with no work packet, approved disposition or explicit evidence-only lane blocks G3.

### Slice-conservation checkpoints

| Gate | Required continuity proof |
|---|---|
| G0 | Canonical baseline plus complete ledger count/ID fingerprint. |
| G1 | Every incoming ID assigned; no silent omission, duplicate ownership or unapproved deferral. |
| G3 | Every assigned ID has owner, impact map, tests, evidence and dependency readiness. |
| G4 | Implementation commits map back to IDs; no unrelated hidden semantic work. |
| G5–G9 | Test and behavior results attach to IDs and invalidate stale predecessor evidence. |
| G10 | Full ledger reconciliation: zero missing/open/unowned/untested/evidence-stale applicable rows. |
| G11 | RC freezes exact ledger, source and evidence fingerprints. |
| G12–G14 | Full certification and both native artifact audits attach to the frozen IDs/source. |
| G15 | Promotion rechecks ID conservation and exact artifact provenance. |
| G16 | Recurrences, root causes, prevention and remaining future placements are explicit. |

### Hardened final-closure execution board

| Wave | Primary work | Parallelizable lanes | Exit condition |
|---|---|---|---|
| 0 | Freeze 282-row scope and dependency graph | inventory, source mapping, test mapping, documentation impact | zero missing/duplicate/unowned/unplaced IDs |
| 1 | Escaped user defects | copy/version; symbol semantics; navigation; Research UI | focused tests + direct browser proof |
| 2 | Confirmed implementation misses | TradeInsight; snapshot broker; session coordinator; routes/surfaces; docs; dependency register | canonical ownership + functional/security/performance proof |
| 3 | Remaining utility remediations | Dashboard, Market Intelligence, Desks, Catalyst, Maintenance | all 13 remediation rows final |
| 4 | Complete v17/v18 proof | backend, renderer, persistence, provider, security, docs lanes | every applicable row current and evidenced |
| 5 | Cross-cutting hardening | roles, desks, viewports, restart, failure, load, data rights | G7–G10 matrices clean |
| 6 | Immutable release | G12 full suite; macOS and Windows in parallel | exact packages pass G13/G14/G15 |
| 7 | Retrospective/handoff | documentation, learning, cleanup, provenance | G16 reconstructable handoff |

### Evidence freshness and invalidation

Evidence records must carry source commit and dependency fingerprint. Changes invalidate at least:

- renderer change → affected UI, responsive, keyboard and native workflow evidence;
- API/state/persistence change → functional, isolation, migration, restart and package evidence;
- provider/router change → entitlement, rights, fallback, rate-limit, performance and decision-truth evidence;
- build/package change → G13/G14/G15 artifact evidence;
- requirement wording or ownership change → G1/G3/G10 traceability evidence;
- repeated defect → original row, linked regressions and dependent package proof.

Only unaffected evidence with a matching fingerprint may be reused.

### Quantitative final-closure scorecard

All values must be achieved before G11:

- requirement identity conservation: **282/282**;
- applicable rows with owner: **100%**;
- applicable rows with regression mapping: **100%**;
- applicable user workflows with browser evidence: **100%**;
- applicable rows with current source/evidence fingerprint: **100%**;
- open/reopened/not-implemented applicable rows: **0**;
- unexplained deferrals: **0**;
- repeated defects without root-cause/prevention evidence: **0**;
- required native platforms passing actual-artifact audit: **2/2**;
- protected Day/Swing/Long formula drift: **0**;
- No Execution or U.S. Equities boundary violations: **0**.

The detailed execution authority is `governance/V18_5_1_HARDENED_FINAL_CLOSURE_BUILD_PLAN.md`.

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
11. Are all input requirement IDs conserved through every slice?
12. Does every PASS have current behavior and exact-package evidence rather than historical/static evidence?
