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
