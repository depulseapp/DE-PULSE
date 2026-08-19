# DE.PULSE — Current Adaptive Build Plan

**Operational overlay date:** 2026-08-19  
**Immutable Stable baseline:** `v18.6.1-stable`  
**Current engineering slice:** Phase 0 Packet C — Chrome + WebKit primary browser-risk routing  
**Authority:** current execution overlay. Permanent contracts and historical release evidence remain intact.

## 1. Normal target lifecycle

`1 development branch → 1 Draft PR → Fast → same PR Ready → Qualified → merge → Release only when release identity/release workflow requires it → Stable`

Rules:

- Never create trigger/retry/certification/promotion branches.
- Never create a second PR just to retrigger CI.
- Fix source/test/gate defects on the same branch and same PR; return to Draft when appropriate, then obtain new exact-head Fast and Qualified evidence.
- Same-SHA infrastructure failure reruns only affected work when possible.
- Main push performs hygiene only.
- Publication uses exact same-run certified native artifacts; no post-certification rebuild.

## 2. G0–G16 execution map

- **G0 Exact Stable Intake:** immutable Stable identity, source SHA/fingerprint, open defects/issues, CI state and dependencies.
- **G1 Immutable Scope:** every committed scope item has stable traceability; no silent additions/removals.
- **G2 Architecture / Data Utility:** owner, consumer, provider, entitlement/rights, source of truth, reuse, freshness, retention and duplication.
- **G3 Design / Dependency / Impact Readiness:** Impact Planner classifies affected surfaces, tests, portability, browser risk and expected CI cost.
- **G4 Development Exit:** one version-development branch, one Draft PR, clean source and scope traceability.
- **G5 Fast Qualification:** cheap exact-head syntax/format/unit/contract checks for affected risk.
- **G6 Integration / Medium Qualification:** affected integration and cross-module evidence.
- **G7 Data / Security / Adaptive Intelligence:** provider/data-rights/security/adaptive evidence when applicable.
- **G8 Performance / Capacity / Stability:** load/runtime/backpressure/stability evidence when applicable.
- **G9 Cross-Module / UI / UX:** affected renderer/browser/interaction evidence.
- **G10 Pre-Freeze Qualified Candidate:** exact-head Qualified success; Release Rehearsal for CI/release changes; Chrome + WebKit primary browser evidence whenever the selected risk/lane requires browser qualification.
- **G11 Immutable Release Candidate:** merged candidate bound to exact Fast + Qualified source head and equal source fingerprint.
- **G12 Full Certification:** authoritative full certification from immutable candidate.
- **G13 Native Packaging / Provenance:** required native packages from candidate.
- **G14 Actual Artifact Runtime Audit:** packaged macOS Apple Silicon and Windows x64 behavior/provenance.
- **G15 Release Assurance / Promotion:** native evidence graphs and exact artifact hashes verified.
- **Publish:** exact certified artifacts only; no rebuild.
- **G16 Adaptive Retrospective / Handoff:** current source of truth, evidence, defects, CI performance and next intake.

No new top-level gates beyond G0–G16.

## 3. Impact Planner v2

Change classes:

- `CI_HARNESS`
- `RELEASE_TOOLING`
- `BACKEND`
- `RENDERER_UI`
- `AUTH_SECURITY`
- `PROVIDER_ROUTER`
- `DATA_RIGHTS`
- `PERSISTENCE`
- `RELIABILITY_PERFORMANCE`
- `CERTIFICATION_GOVERNANCE`

Unknown non-process content fails closed to full qualification.

### Lane selection

- Process/governance/CI-only: `ci-harness` + portability.
- Product/mixed/uncertain: `full` Qualified.
- Normal `full` and explicit `browser` qualification: both primary browser engines must pass.
- `RENDERER_UI`: WebKit evidence is mandatory even if the lane is narrowed.
- WebKit harness/routing changes: real WebKit evidence is mandatory while remaining process-only.
- Backend/provider-only narrowed work: no unnecessary browser-engine runtime.
- `CI_HARNESS` or `RELEASE_TOOLING`: Release Rehearsal required through workflow policy.

### Browser policy

- **Chrome and WebKit are co-primary browser engines.**
- Chrome carries the broad behavioral regression suite.
- WebKit carries the primary cross-engine compatibility suite for core UI/interaction contracts.
- `full` and `browser` candidates require both.
- Renderer/UI changes require WebKit through `webkit_required`.
- Other engines, including Firefox if introduced later, remain secondary/risk-directed unless evidence justifies promotion.
- The design intentionally avoids a blind N-browser matrix: primary engines get mandatory evidence; secondary engines are added only where risk/value warrants it.

## 4. Failure taxonomy

- `PRODUCT_FAIL`
- `GATE_TEST_FAIL`
- `CI_HARNESS_FAIL`
- `INFRA_FAIL`
- `EXPECTED_NOOP`
- `SUPERSEDED`

Failure classification never permits bypassing required quality gates.

## 5. Phase 0 status

### Packet A — COMPLETE

- Impact Planner v2 + self-test.
- Side-effect-free Release Rehearsal.
- Current roadmap/build-plan/process/delivery overlays.
- Honest v18.6.1 reconciliation baseline.
- Portable authoritative handoff.

Merged PR #46; final Fast/Qualified process-only path passed.

### Packet B — COMPLETE

- Fast/Qualified third-party Actions pinned to immutable SHAs.
- Canonical CI dependency lock.
- Playwright `1.62.0` pinned.
- Safe pip cache keyed by browser lock.
- Reproducibility + least-privilege gate.
- `release.yml` Action pinning intentionally deferred to the next genuine release-capable product slice.

Merged PR #47; final Fast #384 and Qualified #134 passed using CI-harness + Ubuntu/macOS/Windows portability. No Stable release was triggered.

Generic workflow linting remains useful and is carried into Packet D unless a safe zero-duplication implementation is completed earlier.

### Packet C — ACTIVE

- Propagate `webkit_required` into Qualified.
- Make Chrome and WebKit the two primary browser engines.
- Keep Chrome as the broad behavioral regression suite.
- Add primary WebKit compatibility proof for watchlist/global-remove/DESKS/no-CURRENT semantics, failure handling, short-height Settings save-bar behavior and centered alert layout.
- Require WebKit for `full`, `browser`, `RENDERER_UI`, and WebKit-harness changes.
- Keep backend/provider-only narrowed work free of unnecessary browser cost.
- Bind WebKit success into exact-head Qualified evidence whenever required.
- Keep other browser engines secondary by default.

Expected Packet C self-validation: process-only → Fast → Qualified CI-harness/portability **plus one real WebKit job** because Packet C changes its own WebKit compatibility harness. Backend/renderer/Chrome product suites remain skipped. This validates the new WebKit execution path without manufacturing a product build.

### Packet D — Durable evidence / telemetry

- Compact durable Stable evidence manifest.
- Lane runtime/queue/cache/failure-class telemetry.
- CI cost/runtime trend detection.
- Generic workflow linting if still outstanding.
- Preserve evidence sufficient to detect workflow amplification/regression early.

### Packet E — Renderer modularization foundation

Incrementally extract capability owners from the large renderer/compatibility stack. Each migration requires deterministic equivalence plus required Chrome/WebKit/native evidence before removing former owners. Never delete based on file age alone.

## 6. Source and repository hygiene

A file being several days old is not a defect. Do not touch unchanged files merely to refresh GitHub dates. Remove a file only when references/consumers/evidence needs are proven absent and protected history is unaffected.

Historical certification, governance, approved reference assets and actively loaded compatibility layers stay until explicit cleanup proof says otherwise.

## 7. Quality floor

Efficiency may reduce duplicate work, runner choice or irrelevant lanes; it may not reduce:

- exact-source provenance;
- deterministic tests;
- Chrome + WebKit primary evidence when browser qualification is required;
- data/security/rights controls;
- macOS Apple Silicon + Windows x64 Stable certification;
- same-artifact publication;
- conserved requirement traceability;
- No Execution and other permanent product boundaries.

## 8. Product intake after Phase 0

Run fresh G0–G3 against current reconciliation and select the highest-value coherent v18.x slice. Provisional priority order remains: user-trust defects → runtime/ADR-GDI reliability → shared intelligence utility consolidation → renderer maintainability → controlled TradeInsight SHADOW integration.
