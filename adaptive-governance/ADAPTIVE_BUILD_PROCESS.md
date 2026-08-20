# DE.PULSE — Adaptive Build Process

**Status:** ACTIVE / GOVERNED  
**Authority:** `governance/ADAPTIVE-OPERATING-CONTRACT.md` controls the permanent process. This file operationalizes it for active release work.

Permanent model:

**Resume → impact-map → reuse → execute smallest affected set → checkpoint → AIPLC → reconcile coverage → certify → deliver → learn.**

**G0–G16 is the permanent top-level gate model. No G17+.**

---

## 1. Resume Reconciliation

Whenever work resumes after interruption, handoff, new conversation, runner interruption or uncertain state:

**detect active release → verify incoming Stable → read checkpoint → read actual branch HEAD/release identity → inspect CI/artifacts → compare source/artifact identity → determine last trustworthy PASS → resume smallest required next step.**

Actual GitHub evidence outranks chat summaries or stale checkpoint labels.

Checkpoint-only metadata may advance branch HEAD without changing the candidate fingerprint when the fingerprint contract explicitly excludes those paths.

---

## 2. Change Impact Triage before expensive work

Before running qualification, classify the delta:

**Git diff → canonical owner → dependencies → affected tabs/features/data/roles/providers/runtime paths → evidence inheritance → rerun set.**

Each responsibility becomes:
- `FRESH_REQUIRED`;
- `INHERITABLE`;
- `SENTINEL_REQUIRED`;
- `NOT_APPLICABLE`.

Reuse PASS evidence only when relevant source/artifact fingerprint, dependency/input contract, test definition and security/role assumptions remain equivalent.

When source/tooling changes, invalidate the earliest affected G0–G16 responsibility and its dependent evidence only. Do not invalidate unrelated evidence.

---

## 3. One Build Coordinator

One logical Build Coordinator owns:
- G0–G16 dependency graph;
- authoritative release state conclusions;
- lane dependencies/concurrency;
- evidence inheritance/invalidation;
- checkpoint reconciliation;
- workflow supersession/idempotency.

Separate workflows may exist only as subordinate/reusable lanes. They must not independently mutate conflicting authoritative release state.

CI failures are classified before product health changes:
`PRODUCT_FAIL`, `GATE_TEST_FAIL`, `CI_HARNESS_FAIL`, `INFRA_FAIL`, `EXPECTED_NOOP`, `SUPERSEDED`.

---

## 4. Bounded adaptive decomposition

Heavy work follows:

**Understand → impact-map → decompose → reuse → execute → checkpoint → evaluate → adapt → integrate → certify → learn.**

Split work only when it improves fault isolation, recovery, evidence reuse, speed or clarity.

Independent lanes may run in parallel only within actual CPU/memory/browser/provider/API/DB/model capacity. Splitting must not multiply expensive setup, provider acquisition, AI calls or release mutations.

Canonical details: `adaptive-governance/ADAPTIVE_WORK_DECOMPOSITION_CONTRACT.md`.

---

## 5. Three-depth AIPLC execution

### Level 1 — Delta AIPLC after every meaningful build/checkpoint
Audit changed/affected areas plus dependency sentinels.

For each affected tab/feature/material datum evaluate:

**datum → purpose → canonical owner → consumer → freshness/materiality → independence/correlation → interpretation → decision value → explanation → outcome → learning.**

Disposition:
- KEEP / STRENGTHEN;
- SYNTHESIZE / CORRELATE;
- CONSOLIDATE / REUSE;
- DEMOTE / DRILL-DOWN;
- SUPPRESS / REMOVE;
- ABSTAIN / UNKNOWN.

Every material challenge produces:
1. immediate fix/truthful disposition;
2. reusable recurrence prevention and cross-product pattern scan.

### Level 2 — G10 Full Coverage Reconciliation
Every required tab/feature/data/security/performance/process responsibility must be freshly evidenced or explicitly inherited from equivalent trustworthy evidence.

### Level 3 — G16 / Major Closure
Run the deep system retrospective across product intelligence, architecture, source quality, data utility, providers, reliability, performance, UI/UX, security, build/CI and outcome learning.

Mechanically identical evidence may use `NO NEW LEARNING / EVIDENCE EQUIVALENT` rather than duplicate reports.

---

## 6. Efficient qualification execution

### G5 — FAST
Run changed-area unit/static/schema/smoke/regression tests plus required sentinels.

### G6 — Integration / MEDIUM
Run affected dependency/cross-module workflows only.

### G7 — Data / Security / Adaptive Intelligence
Run impacted data/provenance/authorization/adaptive-governance lanes. Preserve point-in-time/no-lookahead and SHADOW→VALIDATED→APPROVED→PRODUCTION boundaries.

### G8 — Performance / Capacity / Stability
Measure changed hot paths and affected capacity/runtime behavior. Full supported-load closure belongs at the appropriate release closure (including mandatory v18.5 coverage), not every commit.

### G9 — Cross-Module / UI / UX
Per build, review changed/affected surfaces and dependency sentinels. G10 reconciles complete tab/role/viewport coverage using fresh + valid inherited evidence.

### G10 — Pre-Freeze
No missing/unowned coverage. Full coverage ledger must reconcile to actual candidate fingerprint.

### G11/G12
Freeze immutable RC, then run one authoritative full certification on that RC. Do not repeat G12 because of metadata/checkpoint-only commits that do not alter RC/test contract.

---

## 7. Provider/data/model load discipline

Provider/data tests:
- use fixtures, historical replay or canonical cached evidence when live provider behavior is not the test subject;
- bounded live smoke only for actual live routing/fallback/freshness/entitlement behavior;
- equivalent test lanes reuse canonical acquisition instead of refetching.

AI/LLM tests:
- use deterministic canonical evidence packages/fingerprints for grounding/regression where possible;
- use bounded true-model samples only when model-runtime behavior is the actual requirement;
- reusable synthesis is reused across equivalent consumers/tests;
- private context/rights/entitlements remain isolated.

---

## 8. Role-aware execution

When role-sensitive behavior is affected:
- frontend composition and backend authorization are implemented/tested together;
- ADMIN is capability-based, never blanket power by label alone;
- USER/DEMO receive no implementation machinery or privileged payloads;
- direct unauthorized routes/APIs are denied server-side;
- suppressed content reflows naturally without dead geometry;
- deterministic Day/Swing/Long logic remains role-invariant.

G9/G10 coverage may reuse unchanged role×surface×viewport evidence when dependency/security contracts are equivalent.

Canonical details: `ROLE_AWARE_UI_COMPOSITION_CONTRACT.md` and `ROLE_AWARE_SESSION_SECURITY_CONTRACT.md`.

---

## 9. Functionality Utility / reuse execution

Every new/materially changed tab, card, engine, scanner, detector, job, watcher, scheduler, provider path, dataset, metric/model, alert, persistence field or admin operation must prove:
- purpose/consumer;
- one canonical owner;
- reuse/consolidation before addition;
- relevant correlation;
- freshness/materiality/retention/rights behavior;
- bounded provider/runtime/storage cost;
- correct UI disposition;
- obsolete path/surface/job retirement where safe.

Permanent default:

**one canonical intelligence owner → one deep-evidence home → concise contextual reuse elsewhere.**

Run `functionality_utility_checkpoint_gate.py` before G10.

---

## 10. Shared Symbol Intelligence execution

Equivalent lawful symbol/dataset/intelligence work is canonical once and fanned out to authorized consumers.

Required invariants:
1. Global Symbol Registry/shared demand union owns processing membership.
2. Equivalent acquisition/calculation/synthesis is reused.
3. Simultaneous equivalent misses coalesce where practical.
4. Material changes invalidate only affected downstream state.
5. Dynamic attention protects decision-critical evidence.
6. User/private/tenant/rights/entitlement boundaries partition sharing when required.
7. One user cannot starve higher-value shared work.
8. Overlapping-user cost should scale primarily with unique canonical demand.

Canonical details: `SHARED_SYMBOL_INTELLIGENCE_PROCESSING_CONTRACT.md`.

---

## 11. Governance-to-Implementation execution

Lifecycle:

**Governed → Implemented → Enforced → Evidenced → Delivered → Learned.**

At G1 classify applicable governance as:
- `CURRENT_RELEASE_BLOCKER`;
- `CURRENT_RELEASE_PROCESS_HARDENING`;
- `NEXT_RELEASE_MANDATORY_ENTRY`;
- `FUTURE_STRATEGIC`.

Documentation alone cannot close a release requirement.

For v18.2, product/security blockers and evidence-trust hardening must close before promotion; v18.3 source-changing consolidation items remain named mandatory entry work and must not be smuggled into frozen v18.2 G1.

---

## 12. G0–G16 execution map

- **G0:** exact baseline/repository/checkpoint truth.
- **G1:** immutable scope + applicable governance disposition.
- **G2:** architecture/data utility/canonical ownership/impact graph.
- **G3:** dependency/readiness/decomposition/evidence plan.
- **G4:** development exit.
- **G5:** FAST changed-area qualification.
- **G6:** affected integration/MEDIUM qualification.
- **G7:** data/security/adaptive intelligence.
- **G8:** performance/capacity/stability.
- **G9:** cross-module/UI/UX.
- **G10:** complete coverage reconciliation/pre-freeze.
- **G11:** immutable RC.
- **G12:** full RC certification.
- **G13:** native packaging/provenance.
- **G14:** actual artifact runtime audit.
- **G15:** release assurance/promotion.
- **G16:** deep Adaptive Retrospective/handoff.

No G17+.

---

## 12A. Adaptive CI decision and execution loop

Every change passes this loop before hosted execution is selected:

`classify change → map canonical owners/dependencies → invalidate affected evidence → add historical-risk lanes → enforce mandatory gates → estimate runtime/cost → execute cheap-first → classify result → reuse or rerun smallest lane → learn at G16`.

Process rules:

1. The planner is advisory for optimization and additive risk detection; deterministic policy owns mandatory lanes.
2. Required CI runs whenever it supplies necessary independent, platform, security, database, browser, immutable-artifact or release evidence. Budget does not overrule quality.
3. `ci-fast` blocks `ci-qualified`; G10/G12 blocks native G13/G14; G13/G14 blocks G15; G15 blocks no-rebuild publication and G16.
4. Independent lanes run in bounded parallel. Failed or cancelled prerequisites prevent dependent jobs from starting.
5. Each job consumes a fixed source/evidence fingerprint and emits requirement IDs, selected reason, result, failure class, runtime, runner OS, cache result, retry relationship and artifact provenance.
6. Unchanged PASS evidence is reused only when its source, dependencies, gate semantics and platform requirements still match.
7. A retry is forbidden until the failure is classified as product, test/gate, harness, infrastructure, expected no-op or superseded. Retry only the failed and invalidated dependents.
8. Workflow and native-harness changes receive their own tests: action/workflow lint, UTF-8 reads, OS-aware permission assertions and readiness probes based on actual process/port signals rather than fragile log timing.
9. CI cannot edit/push product source or create/delete temporary committed workflows. Release publication receives the only narrowly scoped write capability.
10. G16 reviews cost per trustworthy evidence, not raw spend alone. The $5 amount is a soft anomaly signal; prevention and quality remain binding.

`CI-ADAPTIVE-18.5.1-001` and `governance/GITHUB_ACTIONS_EFFICIENCY_CONTRACT.md` are mandatory inputs to G1, G3, G10 and G16.

---

## 12B. Independent-audit execution contract

Applies to: `AUDIT-18-UI-001`, `AUDIT-18-AI-001`, `AUDIT-18-AI-RIGHTS-001`, `AUDIT-18-PROVIDER-001`, `AUDIT-18-ARCH-001`, `AUDIT-18-CI-001`, `AUDIT-18-SECURITY-001`, `AUDIT-18-PROVENANCE-001`, `AUDIT-18-TRADER-001` and `AUDIT-18-QA-001`.

Every build/checkpoint begins by loading the ten `AUDIT-18-*` records from the reconciliation ledger together with their existing reopened-defect and CI links. The audit corpus is never reconstructed from conversation memory.

1. **Classify and conserve.** Map each change to audit IDs, canonical owners, dependencies, user impact, rights and invalidated evidence. G1 must reconcile 295 identities.
2. **Behavior before markers.** For UI/UX changes, write the real event sequence first: pointer/focus/scroll state plus repeated live updates and structural mutations. A class/token/source marker is supporting evidence only.
3. **Incremental live rendering.** Quote/SSE events may patch keyed values; full surface remounts require an explicit structural reason. Capture semantic anchor, focus and selection before any unavoidable remount and restore them after layout stabilizes.
4. **AI correctness and safety.** Freeze prompt/safety/schema/model/provider fingerprints; build the actual bounded evidence envelope; enforce schema capability or abstain; run golden, citation, contradiction, missing-evidence and injection evals on every affected change.
5. **Rights-aware egress.** Resolve provider×dataset AI-use permission before external-model calls. Unknown or unbound rights deny egress and expose a truthful diagnostic; no API key or commercial plan infers permission.
6. **Provider experiments.** Production uses one canonical route; bounded rotating shadow samples collect comparative evidence. Experiments cannot exceed budget/rights/backpressure controls or influence production before promotion approval.
7. **Learning experiments.** Outcome/provider/model policies are versioned proposals. Use cutoff-safe data, sample-depth labels, calibration/drift measures and explicit rollback; protected formulas never self-modify.
8. **Strangler architecture cleanup.** Move one canonical owner at a time behind stable interfaces, prove deterministic equivalence and retire the old owner in the same bounded packet. Do not perform a broad rewrite.
9. **Security/supply chain.** Applicable changes run secret, dependency, vulnerability, SBOM and provenance lanes cheap-first. Desktop credential migration and HTTP timeout changes require failure/recovery tests.
10. **CI truth.** The active branch must always match at least one durable workflow trigger. A workflow-trigger gap is a release blocker, not an infrastructure footnote.
11. **Semantic test review.** When expected behavior changes—such as final-desk symbol removal—update requirement, UI copy, backend, undo/recovery and tests together. A historical test cannot override current frozen intent.
12. **Evidence invalidation.** Any prompt/schema/router/rights/UI-store/workflow/security/provenance change invalidates its dependent PASS evidence even when the visible version string is unchanged.

A checkpoint cannot report PASS while an applicable audit subfinding is absent, collapsed into a generic note, or supported only by historical/static evidence.

---

## 13. 10/10 Process acceptance

The Build Process is 10/10 only when it proves:
- no stale state is trusted over actual GitHub truth;
- no unnecessary full reruns occur;
- no affected responsibility escapes coverage;
- inherited evidence is equivalence-bound;
- bounded parallelism reduces rather than multiplies load;
- failures rerun the smallest dependent set;
- AIPLC creates both fixes and prevention;
- frozen G1 is protected;
- one authoritative coordinator/state exists;
- another runner/conversation can resume without guessing.

## 14. Cross-assistant/account resume execution

At G0 and after interruption, load `AGENTS.md` or `CLAUDE.md`, then the vendor-neutral portability contract, `handoff/CURRENT.md`, actual GitHub branch/PR/tag/check/artifact state and the Build State Ledger. Run `adaptive_resume_gate.py` before implementation or evidence reuse.

The process must produce the same result regardless of assistant vendor or ChatGPT/Claude user account. If GitHub is inaccessible, stop for repository connection; do not infer status from model memory. If continuity artifacts disagree with GitHub, classify the stale artifact, correct it and resume from the earliest invalidated G0–G16 owner.

After meaningful work, commit candidate changes first and checkpoint metadata second so `candidateSourceCommit` can name the exact candidate while fingerprint-excluded metadata may advance branch HEAD safely.

---

## 15. Canonical CI retry/resume and continuous-improvement process

The adaptive build process must operate one permanent CI system that evolves in place.

### 15.1 Canonical workflow invariant

Normal active CI is limited to:

1. `ci-fast.yml` — cheap automatic preflight/impact classification;
2. `ci-qualified.yml` — parameterized affected-area qualification through G10;
3. `release.yml` — parameterized G11–G16 release orchestration.

Release versions, RC identities, gates, lanes, OS targets and retry/resume points are parameters. A new release must not create a new workflow family.

A workflow-allowlist gate must reject unapproved `.github/workflows/*.yml` additions unless the change explicitly modifies the permanent CI architecture and carries its own migration/rollback evidence.

### 15.2 Same-workflow failure handling

After a failure:

- `INFRA_FAIL`: use GitHub failed-job rerun on the same workflow/SHA where possible;
- `CI_HARNESS_FAIL`: fix the shared harness/tooling, add regression coverage, then dispatch/rerun the same workflow and affected lane;
- `GATE_TEST_FAIL`: correct the canonical test contract and rerun that gate plus invalidated dependents in the same workflow family;
- `PRODUCT_FAIL`: fix source and rerun the same workflow on the new SHA from the earliest invalidated gate;
- `EXPECTED_NOOP`: record and continue without inventing a retry;
- `SUPERSEDED`: preserve history but do not spend compute recreating obsolete evidence.

The process must never respond to an ordinary lane failure by committing `*-retry.yml`, `*-monitor.yml`, `*-probe.yml`, `*-recovery.yml`, `*-certification.yml` or `*-publish.yml` workflow files.

### 15.3 Independent PASS preservation

Evidence reuse is lane-specific. If macOS G13/G14 passes and Windows fails, the macOS PASS remains authoritative when exact source, package contract, test semantics and artifact identity are unchanged. Repair and rerun Windows only; G15 consumes both matching evidence graphs.

This principle applies to any independent lane whose dependency fingerprint remains valid.

### 15.4 Learning must change the canonical system

G16 process learning is incomplete until every material repeated CI/process failure has either:

- changed the canonical workflow/tool/test and added regression/prevention evidence; or
- received an explicit evidence-backed `NO_IMPLEMENTATION_CHANGE_REQUIRED` disposition.

A workaround that succeeds once but leaves the canonical implementation unchanged is temporary recovery, not adaptive maturity.

The v18.5.2 Windows provenance lesson is the reference pattern: platform-specific filesystem materialization must not define canonical source identity; shared provenance tooling must derive identity from platform-neutral canonical Git object bytes and test this across supported runner families.

### 15.5 Workflow and branch hygiene inside G16

Before G16 closes:

- inventory active workflows against the allowlist;
- delete obsolete version-specific/diagnostic workflow files from the active branch after evidence is safely preserved in Git history/releases;
- inventory branches and classify each as active, unique/unmerged, protected evidence dependency or obsolete;
- preserve/reconcile unique content, then delete merged/obsolete RC/retry/certification/promotion/old-development branches;
- converge to `main` + one active release development branch + short-lived feature/fix branches;
- treat RC as immutable SHA/checkpoint and Stable as tag + GitHub Release;
- record pre/post workflow count, branch count, removed items, retained exceptions and evidence/cost avoided.

A release cannot claim 10/10 adaptive process maturity while known orphaned workflows or obsolete operational branches are intentionally left without disposition.

This section permanently operationalizes the CI convergence rules from the Adaptive Roadmap/Build Plan and `governance/GITHUB_ACTIONS_EFFICIENCY_CONTRACT.md` without adding a new top-level gate.

---

## 16. v18.8 audit-derived 10/10 execution process

The v18.8.0 audit findings are now a permanent process contract. They apply first to v18.8.1 and recur wherever the same responsibility is affected.

### 16.1 Resume must converge actual Stable truth

At G0 and after any Stable publication:

1. read the actual GitHub Stable tag/Release and completed Release workflow;
2. compare release identity, build/checkpoint/evidence files, Stable manifest, handoff and CURRENT Adaptive overlays;
3. classify every mismatch in one Release State Coherence result;
4. correct stale continuity metadata before normal next-release product implementation;
5. prove a fresh assistant/account resolves the same current Stable and next action.

A successful publication with stale handoff/checkpoints is **publication PASS but continuity/G16 incomplete**. Do not infer the old checkpoint is authoritative merely because it is durable.

### 16.2 Cheap-first Release State Coherence

Fast execution order must prefer:

`checkout → Python/impact → release-state coherence → workflow/ledger/portability/provenance/release identity → Python syntax → Go/Node setup when still required → code/renderer tests → exact-head status`.

The coherence validator must aggregate errors. If VERSION, checkpoint, manifest and handoff are all stale, the same run reports all four; it must not intentionally force four sequential red/green cycles.

Cheap governance failure stops expensive setup where possible. This is optimization, not evidence reduction.

### 16.3 G11 early publication-feasibility check

Before G12:

- derive the target tag from canonical release identity;
- verify branch/release version compatibility;
- verify build ID and predecessor/baseline semantics;
- verify required release scaffold/certification script exists;
- inspect the target Stable tag;
- if it already exists at a different immutable candidate, fail G11 before G12/native work;
- if it is the exact already-published candidate, classify idempotently rather than rebuilding;
- retain the publication-time collision guard as a second independent safety barrier.

Release #27 is the recurrence reference: the late guard protected v18.7 correctly, but the same impossible publication should now be rejected before expensive certification.

### 16.4 Data-time truth process

For every market datum used in freshness/decision truth:

- preserve `providerEvidenceAt` separately from `retrievedAt`/observation time;
- never substitute current wall-clock time merely because provider evidence time is missing or unparsable;
- missing evidence time becomes explicit UNKNOWN/degraded/ABSTAIN according to consumer requirements;
- tests must distinguish fresh provider evidence, stale evidence, missing timestamp, retrieval time and cached observation time.

### 16.5 Discovery-universe semantic process

Universe acquisition and universe eligibility are separate responsibilities:

1. acquisition obtains provider rows through Smart Provider Router v2-compatible provider policy;
2. an explicit eligibility function/policy owns exchange, tradability, symbol form, asset-type and any liquidity/optionability restriction;
3. any `has_options` or similar provider filter is intentional, named and tested; it cannot silently redefine the advertised broad universe;
4. Scanner/Radar continue to share one canonical neutral universe owner;
5. BroadSnapshotBroker remains the snapshot reuse/coalescing owner; no parallel broker/router/freshness engine is introduced.

### 16.6 Shared-universe robustness process

The canonical shared universe must guarantee:

- neutral health/diagnostic naming for all consumers;
- refresh ownership released on success, provider failure, context cancellation and unexpected loader failure;
- waiters cannot remain permanently blocked by a failed refresh owner;
- original provider evidence time is preserved through failed refresh;
- bounded retry suppression remains distinct from freshness TTL;
- deterministic fallback is visibly fallback, never silently provider-fresh.

### 16.7 Renderer and test ownership migration

v18.8.1 follows a strangler process:

- choose one bounded capability owner at a time;
- establish stable module/test names independent of release number;
- prove equivalence/behavior;
- move callers/tests;
- retire the superseded owner only when unique evidence is preserved;
- do not mass-rename or mass-delete the renderer/test corpus;
- historical release-number tests may remain as provenance until their unique contract is represented in capability-oriented evidence.

### 16.8 Safe manual CI

Automatic PR lifecycle remains canonical. Manual dispatch exists for controlled diagnostics/recovery only. Default inputs must choose the smallest safe/adaptive lane; selecting `full` or release-capable work requires explicit operator intent.

### 16.9 Failure and learning rule

The #413–#415 and #418 sequences are recurrence evidence. A release-metadata or cache-identity failure is not considered adaptively learned until the canonical coherence/identity tooling can detect the same class in one preflight. G16 records the reduction in duplicate runs/setup and verifies quality gates remain unchanged.

No new top-level gate is introduced; these checks live inside G0/G3/G5/G10/G11/G16 as appropriate.