# DE.PULSE — Adaptive Build Process

## Permanent Resume Phase

A **Resume Reconciliation** phase is mandatory whenever work starts after interruption, handoff, new conversation, runner interruption, or uncertain prior execution state. It belongs inside the existing G0–G16 process and is not a new gate.

Resume Reconciliation performs:

**Detect active release → verify incoming Stable → read checkpoint → read actual branch SHA/release identity → inspect CI/artifacts → compare fingerprint/immutable RC identity → determine last trustworthy PASS → resume next required step.**

### Evidence reuse

Reuse a PASS only when its evidence is bound to the unchanged relevant source fingerprint or immutable RC/artifact identity. Do not infer PASS from chat text, a checkpoint label alone, or an older build with similar code.

### Invalidation

When source changes, determine the earliest affected G0–G16 gate and mark that gate and dependent downstream evidence for requalification. Do not invalidate unrelated earlier evidence, and do not retain later PASS evidence across a changed fingerprint.

### Persistence cadence

Commit meaningful implementation work before leaving its phase. Update checkpoint metadata after meaningful commits and gate transitions. Checkpoint-only commits live under the fingerprint-excluded `.depulse-certification/` area so recording recovery state does not mutate the product candidate fingerprint.

### Interrupted CI/native work

Inspect current GitHub run/job state before launching replacements. Resume/retry only incomplete or failed lanes where exact-source reuse is valid. A stopped conversation does not justify rerunning a green G12 or native platform audit.

## Permanent Adaptive CI execution

The release uses one logical Build Coordinator across G0–G16. Jobs may be parallelized only when dependencies and resource limits allow it, while each responsibility retains one canonical owner.

Before expensive qualification, run the preventative-learning preflight against known release-learning patterns. CI failures are classified as `PRODUCT_FAIL`, `GATE_TEST_FAIL`, `CI_HARNESS_FAIL`, `INFRA_FAIL`, `EXPECTED_NOOP`, or `SUPERSEDED` before candidate health is changed.

Mutating workflows are idempotent: a clean/no-change commit attempt is successful no-op behavior, not a failed build. Checkpoint/evidence-only changes must not wake unrelated full certification/native workflows.

The Build State Ledger is reconciled from actual branch HEAD, release identity, fingerprint, CI and artifact evidence. It is never trusted merely because a JSON checkpoint says PASS.

### Safety

Checkpoint/resume and adaptive CI never weaken G0–G16, native runtime requirements, Stable transformation/certification, source provenance, or the permanent No Execution Boundary.

Canonical protocols:
- `adaptive-governance/BUILD_RESUME_PROTOCOL.md`
- `adaptive-governance/ADAPTIVE_CI_OPERATING_CONTRACT.md`

## Permanent Role-Aware Composition Execution

Role/capability composition is executed as part of the normal G0–G16 process; it is not an optional polish pass and does not create a new gate.

Required execution pattern:

- **G1 — Immutable Scope:** freeze which tabs/shell controls change and the role/capability matrix for OWNER, selected/full-capability ADMIN, limited ADMIN, USER and DEMO.
- **G2 — Architecture / Data Utility:** identify the canonical capability owner, backend authorization boundary, sensitive-data suppression rules and reusable layout/composition primitives. Do not create parallel role-specific business logic or duplicate market pipelines.
- **G3 — Design / Dependency Readiness:** define page hierarchy, conditional sections, Administration-tab utility when applicable, navigation/direct-route behavior and role × viewport acceptance matrix before coding proceeds.
- **G4 — Development Exit:** frontend composition and backend authorization must be implemented together. A hidden button with an callable unauthorized API is a defect; an API-denied capability with awkward empty UI geometry is also a defect.
- **G7 — Data / Security / Adaptive Intelligence:** verify unauthorized roles are not sent secrets, token hashes, provider credentials, owner-only security state, raw implementation machinery or privileged administrative payloads merely to hide them client-side.
- **G9 — Cross-Module / UI / UX:** audit every tab plus global shell/navigation under each required role/capability profile and supported viewport family.
- **G10/G12:** role-aware regression evidence is mandatory before freeze/full certification.
- **G16:** record role/tab audit results, learn from recurring composition/security failures and consolidate them into reusable primitives/regressions.

G9 acceptance includes:

- correct visible and forbidden surfaces for each capability profile;
- SERVER-side denial/redirect for direct unauthorized route/API access;
- no implementation machinery for USER/DEMO;
- capability-scoped ADMIN behavior rather than blanket admin power;
- no blank grid cells, orphan headings/dividers, placeholder height, awkward whitespace or broken page endings after suppression;
- natural grid/flex reflow at desktop/tablet/narrow widths;
- preserved information hierarchy: privileged diagnostics do not automatically rise above primary market intelligence;
- correct keyboard/focus/tab order after conditional composition;
- no overlap, clipping, horizontal page scroll or unreadable compression;
- protected deterministic Day/Swing/Long formulas remain identical across roles;
- any new Administration or other tab is present only where justified by utility and capability ownership.

Permanent implementation rule: **compose from capabilities and canonical hierarchy; do not render the OWNER layout and subtract pieces afterward.**

Canonical role-aware rules: `adaptive-governance/ROLE_AWARE_UI_COMPOSITION_CONTRACT.md`.

## Permanent Functionality Utility, Reuse, Correlation & Surface Checkpoint

Every release must run the **Functionality Utility, Reuse, Correlation & Surface Checkpoint** for every new or materially changed tab, sub-tab, panel, engine, job, scheduler, watcher, dataset-facing workflow, metric, model, alert, API/provider acquisition path, persistence field and administrative operation.

This checkpoint is blocking inside existing G0–G16 and creates no additional gate.

Required execution pattern:

- **G1:** inventory every proposed change and record whether it should REUSE, EXTEND_EXISTING, CONSOLIDATE, remain INTERNAL/DRILLDOWN, justify a NEW_SURFACE, REMOVE/RETIRE an older owner, or DEFER.
- **G2:** prove canonical ownership, existing-data reuse, fetch-once/calculate-once behavior, required correlations, freshness/materiality, retention, rights/sensitivity and provider/runtime cost.
- **G3:** prove where the capability belongs in the existing product hierarchy. A new tab must have explicit workflow/security/clarity justification; organizational convenience is insufficient.
- **G4:** implement the approved owner/disposition without introducing convenience duplication or parallel acquisition/computation.
- **G7/G8:** verify data/evidence integrity, adaptive-governance boundaries, provider efficiency, capacity and bounded runtime/storage impact.
- **G9:** audit every tab and major function for repeated information, overlapping workflow ownership, duplicated deep-evidence homes, orphan data, unnecessary user surfaces and hierarchy drift.
- **G10:** rerun the machine-checkable registry gate against the actual candidate. Unresolved duplicate ownership, unjustified new surfaces or unowned/unconsumed functionality blocks freeze.
- **G16:** record what was reused, consolidated, retired or simplified and convert recurring duplication patterns into preventative tests/governance.

Permanent default: **one canonical intelligence owner → one deep-evidence home → concise contextual reuse elsewhere.**

Background work follows the same rule. A temporal preparation job, event watcher or integrity task does not automatically become a separate engine or tab. It must first reuse the canonical scheduler/router/cache and coalesce the same provider/dataset work where practical.

Every release maintains `functionality_utility_registry.json`; `functionality_utility_checkpoint_gate.py` verifies primary navigation coverage and required ownership/reuse/correlation/UI dispositions before G10.

Canonical checkpoint rules: `adaptive-governance/FUNCTIONALITY_UTILITY_INTEGRATION_CHECKPOINT.md`.

## Permanent Adaptive Work Decomposition & Gate Evolution

Large implementation, qualification, audit and delivery responsibilities must be decomposed when smaller independently verifiable work packages materially improve fault isolation, resumability, safe parallelism, evidence reuse, resource efficiency or learning.

Execution pattern:

**Understand → decompose → map dependencies → reuse → execute → checkpoint → evaluate evidence → adapt next work → integrate → certify → learn.**

The Build Coordinator owns the dependency graph. Each decomposed package has one canonical owner, explicit inputs/dependencies, completion criteria, evidence, downstream consumers and invalidation rules. Independent packages may run in parallel; dependent work remains ordered. Parallelism must be bounded so splitting work does not multiply runner load, provider/API calls, database work, browser load or duplicate expensive qualification.

Use meaningful intermediate checkpoints so a late failure does not force unrelated qualified work to rerun. Preserve unchanged-source PASS evidence and rerun the smallest affected/dependent set after failure or interruption.

Under the current gate map:

- G1/G2/G3 define decomposition, dependencies and evidence before heavy work;
- G4 may use independently testable implementation packages;
- G5/G6 provide early fast/medium checkpoints;
- G7/G8/G9 may shard independent data/security/adaptive, performance/capacity and role × viewport × surface work before producing one gate conclusion;
- G10 verifies all required sub-checkpoints/lanes are complete;
- G12 may internally decompose full certification but produces one authoritative RC-bound conclusion;
- G13/G14 keep macOS and Windows packaging/runtime evidence independent;
- G16 measures whether decomposition reduced cycle time/failure repetition and removes redundant packages/checks.

### Gate evolution

G0–G16 remains the canonical gate map by default, but it is no longer treated as structurally untouchable. If accumulated evidence shows a materially distinct assurance boundary cannot be represented cleanly as a checkpoint, sub-stage, parallel lane or strengthened existing gate, the gate model may evolve.

Before adding a new G-gate, pass the Gate Utility Test: distinct risk/responsibility, independent evidence, non-duplication, material value, canonical ownership, process-wide documentation/automation update, migration clarity for in-flight/future releases, and G16 post-use review.

No workflow may invent an isolated `G17`, `G18`, etc. without revising the canonical process. Prefer in order:

**reuse existing evidence → checkpoints → sharded/parallel lanes → strengthen/reassign an existing gate → add a new canonical gate only when materially justified.**

Canonical rules: `adaptive-governance/ADAPTIVE_WORK_DECOMPOSITION_CONTRACT.md`.

## Permanent Governance-to-Implementation Closure Execution

Governance adoption and implementation closure are separate states. The build process uses:

**Governed → Implemented → Enforced → Evidenced → Delivered → Learned.**

At G1 every applicable governance requirement is classified as `CURRENT_RELEASE_BLOCKER`, `CURRENT_RELEASE_PROCESS_HARDENING`, `NEXT_RELEASE_MANDATORY_ENTRY`, or `FUTURE_STRATEGIC`. No item may remain merely “documented.”

### Current-release enforcement

- **G2/G3:** bind every current blocker/process-hardening item to one canonical implementation owner, dependency graph, naming identity and evidence plan.
- **G4:** current product/security blockers must exist in actual source. For v18.2 this includes capability-based ADMIN authorization and the dedicated capability-scoped Administration composition.
- **G7:** backend authorization, capability delegation, payload redaction, identity/session/presence/admin audit data and direct API denial must prove the same capability truth.
- **G9:** execute the complete role × tab × viewport audit across every tab and global shell/navigation, not only changed surfaces.
- **G10:** refuse freeze if governance claims and implementation/evidence disagree. A stale Build State Ledger, overlapping authoritative workflows, noncanonical active names, blanket ADMIN power, or missing role-composition evidence blocks the affected qualification.
- **G11/G12:** consume only an authoritative reconciled Build State Ledger and one canonical orchestration graph for the immutable RC.
- **G16:** close completed items, explicitly carry `NEXT_RELEASE_MANDATORY_ENTRY` items into the named next release, remove obsolete workflow/naming machinery and add preventative regression for recurring governance-to-implementation gaps.

### Build State Ledger enforcement

The ledger must be derived/reconciled from GitHub branch HEAD, release identity, source fingerprint, CI/job evidence, artifacts, RC/native/package state and promotion/delivery state. Noncanonical status labels or stale source commits are corrected before the ledger is used for resume or certification.

### Build Coordinator enforcement

“One logical Build Coordinator” must be visible in actual workflow ownership. Separate workflows may exist as reusable/subordinate lanes, but they must not independently own overlapping release state mutations or authoritative gate conclusions. Path/event filters, concurrency/supersession and idempotent mutation behavior are mandatory where needed to prevent workflow storms.

### Naming enforcement

Canonical naming applies to active workflows/jobs/artifacts/checkpoints/gates/capabilities as well as docs. Obsolete temporary/numbered/ambiguous release machinery is renamed, consolidated, deprecated or removed before release closure where compatibility does not require it.

Canonical closure rules: `adaptive-governance/GOVERNANCE_IMPLEMENTATION_CLOSURE_CONTRACT.md`.

## Permanent Shared Symbol Intelligence Processing Execution

Shared symbol intelligence is enforced as a blocking architecture/performance checkpoint inside the existing gate model. It is not a cosmetic optimization and does not create a new gate by itself.

Required execution pattern:

- **G1 — Immutable Scope:** inventory every source of symbol demand and every affected consumer, including user tracking, selected symbols, Scanner/Radar, desks, Decision Queue, Rapid Move/Market Shock, preparation checkpoints, catalysts, Research and system context.
- **G2 — Architecture & Data Utility:** bind each dataset/capability to one canonical shared owner and define the processing identity used for safe reuse. Prove that user workspaces contribute demand/context rather than owning duplicate market pipelines.
- **G3 — Design & Dependency Readiness:** map the shared demand union, producer/consumer graph, freshness rules, in-flight ownership, material-change invalidation graph, rights/entitlement partitions, dynamic priority/backpressure and efficiency scorecard.
- **G4 — Development Exit:** reject unintended per-user or per-feature duplicate acquisition, subscriptions, calculations, canonical state or reusable synthesis. Shared provider/router/cache/state owners must exist in actual source where affected.
- **G6 — Integration & MEDIUM Qualification:** exercise overlapping consumers/users for the same symbol and prove one canonical result can safely serve multiple consumers; adding/removing one consumer must not corrupt the others.
- **G7 — Data, Security & Adaptive Intelligence:** verify provenance/point-in-time truth, independent-provider reconciliation, AI evidence-fingerprint reuse, private-context separation, entitlement/data-rights isolation and zero unauthorized cross-user leakage.
- **G8 — Performance, Capacity & Stability:** measure unique-demand scaling, provider calls/subscriptions, duplicate acquisition/calculation, cache/coalescing, reusable synthesis, fan-out, marginal overlapping-user cost, CPU/memory/storage, material-change latency, rate-limit pressure, fairness and long-running stability.
- **G9 — Cross-Module UI/UX:** ensure shared intelligence is composed into the right user surfaces without reproducing the same deep evidence or exposing implementation machinery.
- **G10 — Pre-Freeze Qualification:** unresolved unjustified duplicate processing or an architecture that materially scales as `users × symbols` blocks freeze for affected scope.
- **G12 — Full Certification:** replay applicable shared-processing/security/performance evidence on the immutable RC.
- **G16 — Adaptive Retrospective & Handoff:** record calls avoided, reuse/fan-out ratios, justified duplication, bottlenecks, regressions and new preventative controls.

### Shared execution invariants

1. The Global Symbol Registry is the canonical instrument/shared-processing membership owner.
2. Equivalent demand is represented once in the shared demand union.
3. A fresh equivalent acquisition/calculation/synthesis is reused across compatible consumers.
4. Simultaneous equivalent misses collapse to one in-flight owner where practical.
5. Material changes invalidate only affected downstream state rather than triggering blanket recomputation.
6. Dynamic attention is assigned by materiality/freshness/decision relevance/session/event risk/capacity.
7. Shared work is partitioned when rights, entitlement, tenant/security domain, private prompts/context or model/policy identity differ.
8. One user's burst or oversized symbol set cannot starve higher-priority shared work or other authorized users.
9. Adding an overlapping user/symbol consumer should primarily add authorization/composition/fan-out cost, not another market pipeline.
10. Intentional independent-provider reconciliation is allowed and recorded as validation, not mistaken for wasteful duplication.

### AI/LLM execution

Prefer:

**canonical evidence package → evidence fingerprint → correlation/synthesis → reusable intelligence → user-specific composition/question.**

Do not rerun materially identical AI synthesis solely because another user opens the same symbol. User-specific reasoning is reserved for private context, explicit questions, approved personalization, authorization/entitlement differences or workflow state that materially changes the answer. Cross-user caches must never contain private prompts/context or restricted outputs.

### v18.2–v18.5 progression

- **v18.2:** protect the no-per-user-market-pipeline invariant while closing multi-user identity/workspace scope.
- **v18.3:** implement shared Scanner/Radar acquisition, Session Intelligence Coordinator, Event Intelligence ownership and hosted shared demand union/canonical state.
- **v18.4:** close shared-cache/synthesis entitlement, security and data-rights isolation.
- **v18.5:** execute the full 10/10 overlapping-demand scorecard and treat material scaling/duplication/security shortfalls as closure blockers.

Canonical rules: `adaptive-governance/SHARED_SYMBOL_INTELLIGENCE_PROCESSING_CONTRACT.md`.
