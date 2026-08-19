# DE.PULSE — Current Adaptive Build Plan

**Operational overlay date:** 2026-08-19  
**Immutable Stable baseline:** `v18.6.1-stable`  
**Current engineering slice:** Phase 0 Packet D — durable CI evidence and telemetry  
**Authority:** current execution overlay. Permanent contracts and historical release evidence remain intact.

## 1. Normal target lifecycle

`1 development branch → batch coherent branch work → 1 Draft PR → Fast → same PR Ready → Qualified → merge → Release only when release identity/release workflow requires it → Stable`

Rules:

- Build a coherent branch before opening the PR whenever practical so preparatory file writes do not manufacture `synchronize` events.
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
- **G10 Pre-Freeze Qualified Candidate:** exact-head Qualified success; Release Rehearsal for CI/release changes; Chrome + WebKit primary browser evidence whenever selected risk requires browser qualification; compact CI telemetry retained.
- **G11 Immutable Release Candidate:** merged candidate bound to exact Fast + Qualified source head and equal source fingerprint.
- **G12 Full Certification:** authoritative full certification from immutable candidate.
- **G13 Native Packaging / Provenance:** required native packages from candidate.
- **G14 Actual Artifact Runtime Audit:** packaged macOS Apple Silicon and Windows x64 behavior/provenance.
- **G15 Release Assurance / Promotion:** native evidence graphs and exact artifact hashes verified.
- **Publish:** exact certified artifacts only; no rebuild.
- **G16 Adaptive Retrospective / Handoff:** current source of truth, durable release evidence, defects, CI performance and next intake.

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
- A file matching `release/v*/stable-evidence-manifest.json` is a retrospective evidence index and is process-only, but remains `RELEASE_TOOLING` governed. Other release scripts/artifacts remain full-qualification scope.

### Browser policy

- **Chrome and WebKit are co-primary browser engines.**
- Chrome carries the broad behavioral regression suite.
- WebKit carries the primary cross-engine compatibility suite for core UI/interaction contracts.
- `full` and `browser` candidates require both.
- Renderer/UI changes require WebKit through `webkit_required`.
- Primary WebKit executes on `macos-15` with exact pinned Playwright plus `playwright install webkit`; Linux `--with-deps webkit` is prohibited.
- Other engines, including Firefox if introduced later, remain secondary/risk-directed unless evidence justifies promotion.

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

Merged PR #46. Delivered Impact Planner v2, Release Rehearsal, current governance overlays, honest v18.6.1 reconciliation and portable handoff.

### Packet B — COMPLETE

Merged PR #47. Delivered immutable Fast/Qualified Action pins, dependency lock, Playwright `1.62.0` pin, deterministic pip caching and reproducibility/permission gate. Release-workflow Action pinning remains deferred to the next genuine release-capable product slice.

### Packet C — COMPLETE

Merged PR #48 at `23ecb71f60e1658d68bcef6248044ce53b6dd851`.

- Chrome + WebKit co-primary policy implemented.
- Primary WebKit runs on `macos-15` without Linux apt amplification.
- Core watchlist/global-remove, membership semantics, short-height Settings save bar and centered header compatibility are covered.
- Final Fast #393 PASS.
- Final Qualified #138 PASS with real WebKit + Ubuntu/macOS/Windows portability; unrelated backend/renderer/Chrome product lanes were correctly skipped for the process-only packet.

### Packet D — ACTIVE

Deliver:

1. `release/v18.6.1/stable-evidence-manifest.json` bound to the authoritative checkpoint and explicitly non-redefining.
2. `stable_evidence_gate.py` to fail on Stable run/artifact/fingerprint drift.
3. `ci_telemetry.py` + self-test for per-job queue/runtime/platform consumption and workflow amplification.
4. Qualified telemetry artifact retained 30 days plus human-readable job summary.
5. Linux/macOS/Windows runner-minute visibility; no fabricated currency estimates.
6. Chrome/WebKit dependency setup duration and setup-python pip cache-hit signals when those lanes run.
7. warning thresholds for abnormal per-PR Fast/Qualified/Release run counts; warnings inform investigation rather than blocking legitimate defect correction.
8. zero-network workflow structural lint complementing DE.PULSE semantic policy checks.
9. workflow policy integration so telemetry/evidence/lint cannot silently disappear.

Packet D itself remains `ci-harness` + portability. It must not run WebKit merely because telemetry observes browser jobs, and it must not create a Stable release.

### Packet E — Renderer modularization foundation

Incrementally extract capability owners from the large renderer/compatibility stack. Each migration requires deterministic equivalence plus required Chrome/WebKit/native evidence before removing former owners. Never delete based on file age alone.

## 6. CI telemetry contract

Qualified emits compact operational evidence with:

- exact candidate SHA and selected lane;
- per-job queue seconds and execution seconds;
- consumed runner minutes split into Linux, macOS, Windows and unknown;
- Chrome/WebKit dependency setup seconds and pip cache-hit state when applicable;
- current-PR Fast/Qualified/Release run counts;
- amplification warnings above conservative thresholds;
- explicit `actualCurrencyCost: null` because GitHub billing remains the authority for financial cost.

Telemetry is diagnostic/operational evidence. It cannot replace functional qualification, exact-head statuses, native release evidence or the immutable Stable evidence manifest.

## 7. Source and repository hygiene

A file being several days old is not a defect. Do not touch unchanged files merely to refresh GitHub dates. Remove a file only when references/consumers/evidence needs are proven absent and protected history is unaffected.

Historical certification, governance, approved reference assets and actively loaded compatibility layers stay until explicit cleanup proof says otherwise.

## 8. Quality floor

Efficiency may reduce duplicate work, runner choice or irrelevant lanes; it may not reduce:

- exact-source provenance;
- deterministic tests;
- Chrome + WebKit primary evidence when browser qualification is required;
- data/security/rights controls;
- macOS Apple Silicon + Windows x64 Stable certification;
- same-artifact publication;
- conserved requirement traceability;
- No Execution and other permanent product boundaries.

## 9. Product intake after Phase 0

Run fresh G0–G3 against current reconciliation and select the highest-value coherent v18.x slice. Provisional priority order remains: user-trust defects → runtime/ADR-GDI reliability → shared intelligence utility consolidation → renderer maintainability → controlled TradeInsight SHADOW integration.
