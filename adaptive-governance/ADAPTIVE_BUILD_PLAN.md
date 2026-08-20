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

Provider → Market Mode assessment is a permanent companion method. Maintain `provider_market_mode_integration_registry.json`; every Router/capability provider must receive `INTEGRATED`, `CONTEXTUAL_ONLY`, `NOT_RELEVANT` or `INTENTIONALLY_HIDDEN` for named Market Mode scopes. Provider count never changes a mode; evidence flows through the Smart Router/canonical state and production influence remains `SHADOW → VALIDATED → APPROVED → PRODUCTION`. The existing functionality/utility checkpoint gate owns machine enforcement at G2/G10, avoiding a new top-level gate.

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

## 12A. v18.5.1 audit, containment and urgent-recovery entry plan

v18.5.1 establishes the complete reconciliation control plane, prevents further silent loss and executes the safest urgent recovery work. It does not have to compress every remaining v18 item into one patch. Any item not selected for v18.5.1 must be explicitly assigned to the next evidence-selected v18.x slice with owner, reason, dependency and user impact; it cannot remain an unbound backlog note.

### Required recovery lanes

| Lane | Scope | Required outputs |
|---|---|---|
| A — release/profile truth | COPY-18.5.1-001 and VERSION-18.5.1-002 | Version-aware/release-neutral copy, explicit migration-identifier classification, fresh-install/upgrade/TEST-profile matrix and packaged proof. |
| B — symbol/watchlist semantics | SYMBOL-18.5.1-001 and SYMBOL-18.5.1-002 across Day/Swing/Long | Explicit row-X contract, final-membership behavior, state/persistence tests and current-desk visual-state evidence. |
| C — interaction/render stability | NAV-18.5.1-001 and HOVER-18.5.1-001 | Root-cause removal for live/SSE rerender jumps and hover-triggered blinking; stable DOM/geometry/paint, scroll/focus preservation, pointer-dwell/repeated-update and reload tests. |
| D — Research requirement recovery | Reopened v15.1.0 approved items 17–19 | Responsive layout, truthful labels/freshness/actions/disabled states, missing-evidence recovery and requirement-linked tests. |
| E — full promise reconciliation | 48 inherited requirements, 20 v17 items, all v18 workstreams/release entries, 13 orphaned utility remediations, conversational commitments and accepted defects | Durable implementation ledger, unexplained-gap list, owners, decisions, regression IDs and actual-package evidence. |
| F — confirmed orphan recovery | IMPL-18-TRADEINSIGHT-001 + CONVO-V18-003 | Canonical Smart Router SHADOW adapter, configuration/entitlement, rights/provenance controls, provider → Market Mode disposition for every capability, contextual/integrated/not-relevant/hidden enforcement, tests and actual-package evidence. |
| G — acquisition/coordination remediation | IMPL-18-UTILITY-001 and IMPL-18-UTILITY-002 | Shared snapshot broker/cache, one Session Intelligence Coordinator, coalescing/provider-budget/restart/catch-up tests and package proof. |
| H — surface/route remediation | IMPL-18-UTILITY-003 and IMPL-18-UTILITY-004 | Minimal Discovery hierarchy, canonical evidence homes, safe redirects/deep links and responsive package proof. |
| I — documentation governance | IMPL-18-DOC-001 | Server-authoritative audience policy, role-composed UI/direct paths, Documentation Impact Manifest and cross-role security/package tests. |
| J — dependency readiness | IMPL-17-DEPS-001 | Canonical dependency/readiness registry, User Action Required records, gate binding and role-safe operational proof. |
| K — remaining utility carry-forward | Dashboard, Market Intelligence, three Desks, Catalyst Watch and Maintenance records | Fresh design/behavior/performance/package disposition; no registry-only or target-release-only closure. |
| L — header information hierarchy | HEADER-18.5.1-001 | Preserve the existing ET/PT module unchanged; create a stable session-aware secondary Market Pulse Ribbon for market session, time, coverage and data control; move market instruments tertiary; place Sign Out in the Local Owner menu; prove responsive/accessibility/state truth with no animation-induced flicker or layout shift. |
| M — adaptive CI control plane | CI-ADAPTIVE-18.5.1-001 | Durable fast/qualified/release workflows; impact-and-risk lane planner; deterministic mandatory-gate guardrails; cache/concurrency/retention/security hardening; native harness portability; failure taxonomy and per-run cost/evidence telemetry. The $5 budget is a soft alert, not a quality cap. |

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
- **G1:** freeze the exact v18.5.1 assignments from the seven escaped-defect groups, seven confirmed implementation-gap groups and 13-item utility carry-forward; explicitly place every unselected applicable row into a named/candidate next v18.x lane with owner and reason.
- **G2–G3:** record canonical owner, blast radius, migration/state/API/UI impact, test IDs and package matrix for every row.
- **G4:** implement fixes atomically with their tests; do not hide semantic changes inside structural moves.
- **G5–G6:** pass source/unit and functional/integration coverage, including persistence, membership combinations and live-refresh behavior.
- **G9:** run direct UI/UX acceptance across Day/Swing/Long and Research at desktop, tablet and narrow widths, including scroll/focus continuity.
- **G10:** for a slice release, every item assigned to that slice is final and every remaining applicable row has explicit next-v18.x placement; for the major closure release, all applicable 48 inherited, 20 v17, v18 workstream/release, 13 remediation, conversation and defect rows have fresh final evidence.
- **G11–G12:** freeze one immutable RC only after reconciliation is clean, then run full certification.
- **G13–G14:** build and directly exercise macOS Apple Silicon and Windows x64 artifacts; affected workflows cannot inherit v18.5 package evidence.
- **G15:** block slice promotion for any unresolved current-slice item or unexplained/unowned remainder. Block major closure for any reopened, not-implemented or next-slice placement.
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

### Release selection and feedback loop

After every v18.x slice:

1. ingest defects, misses, performance/SLO results, provider usefulness, data-rights/dependency changes and package findings;
2. invalidate evidence affected by the completed changes;
3. reprioritize the still-open ledger by user impact, recurrence, dependency readiness, risk, cost and learning value;
4. select the smallest coherent next slice;
5. freeze that slice at G1;
6. do not name a release “final closure” until the major-closure readiness conditions are already true.

Candidate version numbers such as v18.6 or v18.7 are capacity, not fixed content. Their exact scope is chosen adaptively at G0–G3.

### Adaptive v18.x execution board

| Wave | Primary work | Parallelizable lanes | Exit condition |
|---|---|---|---|
| 0 | Freeze 296-row scope and dependency graph | inventory, source mapping, test mapping, documentation impact | zero missing/duplicate/unowned/unplaced IDs |
| 1 | Escaped user defects | copy/version; symbol semantics; hover/render stability; navigation; Research UI | focused tests + direct browser proof |
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

- requirement identity conservation: **296/296**;
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

The detailed execution authority is `governance/V18_ADAPTIVE_RECOVERY_AND_CLOSURE_PROGRAM.md`.

---

## 12C. Adaptive CI build lane

`CI-ADAPTIVE-18.5.1-001` is a ledger-backed work packet and cannot be sliced out as “process only.”

### Planned implementation

1. Retire obsolete or one-shot workflows from the active branch after inventory and replace them with a minimal durable set:
   - `ci-fast.yml` for automatic cheap preflight on affected source/config/workflow changes;
   - `ci-qualified.yml` as reusable impact-selected qualification;
   - `release.yml` for G10/G12 through G16 orchestration and exact-candidate publication.
2. Implement a deterministic planner input/output record containing changed paths, canonical owners, dependency fingerprint, affected requirements, prior failure classes, selected/skipped lanes with reasons, runtime/cost estimate and mandatory overrides.
3. Run cheap preflight before expensive jobs: schema/gate tests, source checks, workflow lint, PowerShell lint, UTF-8/path/permission portability checks and native readiness-probe unit tests.
4. Use dependency caches, concurrency cancellation and bounded matrices. Parallelize independent native lanes only after shared prerequisites are green; parallelism reduces elapsed time but does not itself reduce billed minutes.
5. Keep workflow permissions read-only by default. Grant narrowly scoped write permission only to the final publication job.
6. Prohibit committed diagnostic workflows, self-deleting workflows, CI-authored source commits and recursive push loops. Use parameterized `workflow_dispatch` inputs and reusable workflows for diagnostics.
7. Retain development artifacts 3–7 days, failure diagnostics 7–14 days and RC evidence 30 days; place certified deliverables and provenance in the immutable GitHub Release.
8. Treat optional controlled self-hosted runners as iterative native-debug support only. Final macOS Apple Silicon and Windows x64 proof remains clean, exact-candidate, independently reproducible evidence.
9. Classify every failure before retry: `PRODUCT_FAIL`, `GATE_TEST_FAIL`, `CI_HARNESS_FAIL`, `INFRA_FAIL`, `EXPECTED_NOOP` or `SUPERSEDED`.
10. Feed run telemetry into G16. Planner changes remain proposals until shadow replay, validation, approval and production promotion.

### Build-lane acceptance

The lane passes only when:

- a source change cannot avoid required affected-area CI;
- a documentation-only or proven-unaffected change cannot accidentally launch native certification;
- mandatory gates cannot be skipped by the adaptive planner;
- each selected/skipped lane has a machine-readable reason and evidence fingerprint;
- the $5 signal can warn and replan but cannot reduce required evidence;
- repeated harness/infra failures create prevention tests;
- exact final-candidate G13/G14/G15 evidence remains mandatory and publication performs no rebuild.

Detailed authority: `governance/GITHUB_ACTIONS_EFFICIENCY_CONTRACT.md`.

---

## 12D. Independent-audit controlled build packets

The ten `AUDIT-18-*` records raised the conserved authority to 295 rows; `CONVO-V18-003` adds the permanent provider → Market Mode method and raises it to **296 rows**. Each packet must retain its subfindings; closing one subfinding cannot close the parent.

| Packet | Mandatory implementation and acceptance | Required proof |
|---|---|---|
| `AUDIT-18-UI-001` | Patch quote/timestamp/status cells in place; structural renders only for real membership/layout change; keyed hover/focus/selection; semantic scroll anchor capture/restore; explicitly resolve final-desk removal semantics and update labels/backend/tests together; implement the unchanged-size ET/PT Market Pulse Ribbon and tertiary instruments strip in stable containers. | Playwright hover+dwell under live/SSE events; mid-section scroll/focus; add/remove/undo; Day/Swing/Long highlight contrast; Research/header responsive; reduced motion; CLS/repaint budget; macOS/Windows package retest. |
| `AUDIT-18-AI-001` | Bound the actual nested evidence envelope by bytes/tokens; rank compaction by materiality/contradiction/role instead of first-N order; cache key includes provider, requested/actual model, prompt, safety, schema and routing versions; TTL/invalidation; strict schema capability with safe abstention. | Unit/property tests for byte bounds and cache isolation; golden evidence corpus; citation/contradiction/missing-evidence scores; injection/schema-escape corpus; continuous model/prompt eval with cost/latency. |
| `AUDIT-18-AI-RIGHTS-001` | Canonical provider×dataset rights registry; AI-egress allow/deny decision bound to evidence; fail closed when AI-use rights are unknown; redacted diagnostics. | Denied-egress integration tests; approved-evidence fixtures; cross-provider/fallback proof; no secret/raw-rights leakage; exact-package behavior. |
| `AUDIT-18-PROVIDER-001` | Persist provider observations and downstream consumption/outcome links; rotating low-cost shadow sampling; disagreement/truth anchors; usefulness and cost-per-useful-evidence scorecards; paid-provider promotion proposal only after measured lift and rights approval; keep every provider/capability bound to an explicit Market Mode disposition. | Provider/Market Mode registry gate; dataset/session/regime fixtures; bounded call-budget/load tests; avoided-call telemetry; no production auto-promotion; sample-depth/confidence report; rollback. |
| `AUDIT-18-ARCH-001` | Strangler extraction into platform, market data/router, intelligence, persistence and web store/component domains; shared persistence semantics with thin OS adapters; tokens/components/page CSS; stable Makefile/Taskfile commands and release-scaffold manifest. | Deterministic equivalence; dependency graph; source-size/override trend; no duplicate owner; incremental package tests; no wholesale rewrite or protected-formula drift. |
| `AUDIT-18-CI-001` | Current branch cannot bypass CI; durable fast/qualified/release workflows; caches, cancellation, bounded matrices, least privilege, short dev retention, exact-candidate native lanes and no-rebuild publication; retire active historical copies after evidence inventory. | Workflow lint/fixture tests; branch/path trigger matrix; selected/skipped reason record; cost/runtime/cache metrics; macOS+Windows final-candidate provenance. |
| `AUDIT-18-SECURITY-001` | Keychain/Credential Manager on desktop and injected hosted secrets; explicit per-provider clients; response-header/idle/connection/request deadlines; SSE-compatible server timeout policy; external-link interstitial/allow policy; CSP class migration; CodeQL, `govulncheck`, dependency/secret review and SBOM. | Security unit/integration tests, abuse/fault injection, secret migration/rollback, dependency scan outputs, SBOM bound to candidate, no provider data/credential disclosure. |
| `AUDIT-18-PROVENANCE-001` | Sign the immutable Stable source tag and/or publish keyless artifact attestations; preserve existing hashes, source binding and no-rebuild publication. | Verification command/output for source tag and each artifact; identity mismatch negative test. |
| `AUDIT-18-TRADER-001` | Add regime/session/liquidity/sector/catalyst conditioning only with sample sufficiency; measure false, duplicate, late and missed alerts; evidence-thesis change log without positions/P&L; SHADOW competing next-state probabilities with calibration/abstention; consume only approved canonical provider evidence from the provider/Market Mode registry. | Cutoff-safe replay; provider-disposition and source-disagreement fixtures; Brier/log-loss/calibration/coverage only for true probability outputs; sample/confidence labels; no setup-score probability conversion; no automatic formula/policy promotion. |
| `AUDIT-18-QA-001` | Every source marker has a behavioral owner; user-reported interaction defects receive realistic event sequences and dwell; test contracts are reviewed when product semantics change; fresh source/behavior/package fingerprints. | G5 focused, G6 integration, G7 data/security/AI, G8 performance, G9 browser/UX, G10 conservation, G13–G15 exact packages. |

### Slice and dependency rules

- v18.5.1 owns the reopened defect set plus `AUDIT-18-UI-001`, `AUDIT-18-CI-001` containment and applicable `AUDIT-18-QA-001`.
- v18.6 cannot start TradeInsight AI consumption until `AUDIT-18-AI-RIGHTS-001` is designed and fail-closed; AI routing work must jointly close `AUDIT-18-AI-001`.
- v18.7 architecture extraction cannot change protected formulas, No Execution, U.S. Equities scope or user-visible semantics without a separately frozen requirement.
- `AUDIT-18-PROVIDER-001` and `AUDIT-18-TRADER-001` enter production only through `SHADOW → VALIDATED → APPROVED → PRODUCTION`; v18 may close their required foundations while explicitly placing mature policy stages in v19/v20.
- Security, rights, QA and provenance packets are cross-cutting dependencies, not optional cleanup.
- G1 conservation is **296/296**. G10/G15 require all applicable audit subfindings, tests and evidence fingerprints to reconcile.

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

## 15. Assistant-independent build planning

Every release plan must include a provider-neutral continuity packet owned by GitHub:

- `AGENTS.md` and `CLAUDE.md` remain thin adapters to one vendor-neutral contract;
- `handoff/CURRENT.md` is the one current human-readable continuation record;
- the Build State Ledger is the machine-readable release state and names exactly one next action;
- the actual branch/PR/head/tag/check/artifact state is reconciled before either file is trusted;
- permanent or active-release changes are reflected in the canonical Roadmap, Build Plan, Build Process and Delivery Process;
- secrets and account-specific session material are never stored in continuity artifacts.

The plan is not resumable until a newly authorized ChatGPT/Codex or Claude account can determine the same last trustworthy PASS and earliest resume gate without asking the user to reconstruct prior conversations.

---

## 16. Permanent CI convergence and repository-hygiene build plan

Every active release plan must use **one evolving CI implementation**, not create release-specific workflow copies.

### Canonical workflow surface

The only normal active workflow entry points are:

- `.github/workflows/ci-fast.yml` — automatic cheap preflight and impact classification;
- `.github/workflows/ci-qualified.yml` — parameterized affected-area G4–G10 qualification and same-lane retry/resume;
- `.github/workflows/release.yml` — parameterized G11–G16 orchestration, independent native lanes, G15 assurance and no-rebuild publication.

Version, source SHA, release identity, build ID, gate, lane, platform and resume point are inputs/configuration. They must not create `vX.Y-*`, `*-retry`, `*-monitor`, `*-probe`, `*-recovery`, `*-certification` or `*-publish` workflow families.

### Failure/retry planning

For every failed lane, record the classification and exact evidence invalidation before any rerun:

- `INFRA_FAIL` → rerun failed job(s) on the same workflow and same SHA;
- `CI_HARNESS_FAIL` → fix the shared harness/workflow, add regression coverage, rerun the same affected lane;
- `GATE_TEST_FAIL` → fix the canonical test contract, rerun that gate plus only its invalidated dependents;
- `PRODUCT_FAIL` → correct source, then run the same workflow on the new SHA from the earliest invalidated gate;
- independent PASS lanes remain reusable when source/package/test/platform fingerprints still match.

A workaround workflow is not an acceptable prevention action when the canonical workflow can own the fix.

### Mandatory hygiene work packet

Before normal v18.6 product implementation, plan and close a behavior-neutral CI/repository consolidation packet that:

1. inventories every active workflow and branch;
2. preserves unique unmerged content/evidence before removal;
3. implements the three canonical workflows using reusable scripts/actions rather than duplicated YAML logic;
4. folds v18.5.2 lessons—especially canonical Git-object provenance and platform-independent source identity—into shared tooling;
5. adds workflow/harness regression tests for Linux, macOS and Windows assumptions;
6. adds a machine-enforced workflow allowlist preventing unapproved workflow sprawl;
7. removes obsolete version-specific workflow files from active `main` after validation;
8. reconciles and removes merged/obsolete RC, retry, certification, promotion-tooling and old development branches after immutable tag/release/evidence safety is proven;
9. converges the branch model to `main` + one active release development branch + short-lived feature/fix branches;
10. treats RC as an immutable SHA/checkpoint and Stable as tag + GitHub Release rather than long-lived branches;
11. records workflow/branch counts and cleanup disposition in G16.

### Learning-to-implementation rule

A CI/process finding is not complete merely because G16 documents it. The build plan must schedule one of:

- a canonical workflow/tool/test improvement with regression evidence; or
- an explicit `NO_IMPLEMENTATION_CHANGE_REQUIRED` disposition with supporting evidence.

Repeated harness/process defects without a canonical prevention change are release-process blockers.

This section is permanent and strengthens `CI-ADAPTIVE-18.5.1-001`, `AUDIT-18-CI-001`, `governance/GITHUB_ACTIONS_EFFICIENCY_CONTRACT.md` and `governance/REPOSITORY_STRUCTURE_CONTRACT.md` without adding G17+.

---

## 17. v18.8 audit-derived 10/10 hardening build plan

**Classification:** `NEXT_RELEASE_MANDATORY_ENTRY` for v18.8.1 plus permanent recurrence prevention. The packet strengthens G0–G16; it does not create a new gate.

### 17.1 Mandatory work packets

| ID | Build owner / work | Minimum evidence |
|---|---|---|
| `ADAPT-REL-001` | Post-Stable metadata/handoff convergence owner. Reconcile actual `v18.8.0-stable` release truth into checkpoints, Stable evidence manifest, handoff and CURRENT overlays before normal v18.8.1 product implementation. | Actual tag/release/run lookup + coherence check + cross-assistant resume proof. |
| `ADAPT-CI-001` | Implement one `release_state_coherence`-style canonical validator covering release identity, VERSION/app/renderer/cache-bust identity, checkpoints, Stable manifest, handoff, predecessor/baseline and target tag state. | Fixture matrix proving multiple simultaneous mismatches are reported in one run and a coherent state passes. |
| `ADAPT-CI-002` | Add G11 target-release preflight: target tag, semantic version/build ID, predecessor, release scaffold and existing-tag collision before G12 starts. Retain publication collision guard. | Existing-conflicting-tag negative test + valid-new-tag positive test. |
| `ADAPT-CI-003` | Reorder Fast to cheap-first governance/coherence/identity/provenance before Go/Node setup; preserve all existing required tests. | Workflow policy/impact fixtures + run summary showing expensive setup skipped after cheap failure. |
| `ADAPT-CI-004` | Make manual dispatch safe-by-default. `full` qualification is explicit; normal Draft Fast → Ready Qualified remains automatic canonical path. | Workflow input regression proving no accidental full/native run from defaults. |
| `ADAPT-DATA-001` | Extract an explicit discovery-universe eligibility policy. Decide and encode whether provider `has_options` filtering is intentional; do not silently call an optionable subset “broad U.S. equities.” | Provider-query and eligibility fixtures including ETFs/class shares/special symbols and documented exclusions. |
| `ADAPT-DATA-002` | Separate `providerEvidenceAt` from `retrievedAt`; remove any freshness fallback that converts missing provider time to wall-clock “now.” | Missing/invalid timestamp tests proving UNKNOWN/ABSTAIN rather than fabricated freshness; normal timestamp tests preserved. |
| `ADAPT-ARCH-001` | Harden canonical shared-universe owner: neutral health naming, panic-safe/single-flight cleanup, waiter cancellation and exact retry/TTL semantics; preserve Smart Provider Router v2 and BroadSnapshotBroker as sole owners of their responsibilities. | Concurrency, cancellation, failure, panic/cleanup and no-duplicate-owner tests. |
| `ADAPT-UI-001` | Renderer Modularization II: migrate bounded capability owners away from release-number filenames; release identity moves to metadata/cache keys. Preserve layout/behavior unless separately scoped. | Equivalence/browser tests; no duplicate owner; cache/identity regression. |
| `ADAPT-QA-001` | Inventory root version-stacked executable gates/tests; classify ACTIVE_REQUIRED / UNIQUE_HISTORICAL / CONSOLIDATE / RETIRE_CANDIDATE; migrate capability-by-capability only with evidence conservation. | Before/after inventory, behavior equivalence and workflow-policy proof. |
| `ADAPT-GOV-001` | Rename/reframe historical reconciliation output so v18.5.1 remains provenance rather than current-state identity. | Gate output fixture showing historical baseline + actual current release separately. |
| `ADAPT-COST-001` | Measure avoided setup/runs/minutes and cost per trustworthy evidence; optimization cannot weaken mandatory lanes. | G16 telemetry comparing baseline and hardened CI. |

### 17.2 Required execution order

1. **G0 continuity closure:** reconcile v18.8.0 actual Stable truth (`ADAPT-REL-001`).
2. **G2/G3 design:** freeze release-state owner map, discovery-universe semantics, evidence-time semantics and renderer/test migration boundaries.
3. **Cheap CI hardening first:** `ADAPT-CI-001..004` before using v18.8.1 for broad implementation so later failures are cheaper and diagnostic-complete.
4. **Data-truth hardening:** `ADAPT-DATA-001..002` + `ADAPT-ARCH-001` with deterministic fixtures and no hidden behavior expansion.
5. **Renderer/test modularization:** bounded `ADAPT-UI-001` + first `ADAPT-QA-001` capability migration; do not rewrite the whole UI/test corpus.
6. **G10:** reconcile every packet with current evidence and explicit remainder placement.
7. **G11–G16:** normal exact-head one-PR release path; early tag collision guard must be green before full G12/native work.

### 17.3 Cheap-first Fast target order

Preferred preflight order:

`checkout → Python → deterministic impact plan → Release State Coherence → workflow/ledger/portability/provenance/release identity → Python syntax → setup Go/Node only when still required and green → gofmt/vet/tests → renderer tests → exact-head status`.

This ordering is an efficiency improvement only. It removes no required evidence.

### 17.4 10/10 build-plan acceptance additions

The plan is not 10/10 unless it can additionally prove:
- all release-state mismatches are surfaced in one cheap diagnostic pass;
- impossible/conflicting Stable publication is rejected before G12/native spend;
- current Stable metadata is atomically converged after publication;
- market evidence time is never synthesized from retrieval time;
- discovery-universe inclusion/exclusion is explicit and tested;
- renderer/test version stacking has a measured bounded migration path;
- historical ledgers cannot misidentify the current release;
- CI savings come from avoided duplicate/late work, not reduced quality.