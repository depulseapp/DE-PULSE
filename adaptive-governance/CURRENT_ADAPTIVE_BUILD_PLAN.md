# DE.PULSE — Current Adaptive Build Plan

**Operational overlay date:** 2026-08-19  
**Immutable Stable baseline:** `v18.6.1-stable`  
**Current engineering slice:** post-v18.6.1 process hardening  
**Authority:** current execution overlay. Permanent contracts and historical release evidence remain intact.

## 1. Normal target lifecycle

The normal successful build path is intentionally small:

`1 development branch → 1 Draft PR → Fast → same PR Ready → Qualified → merge → Release only when release identity/release workflow requires it → Stable`

Rules:

- Never create trigger/retry/certification/promotion branches.
- Never create a second PR just to retrigger CI.
- A source/test defect is fixed on the same branch and same PR; return PR to Draft when appropriate, then Fast and Ready again.
- A same-SHA infrastructure failure reruns only the affected failed job/run when possible.
- A main push performs hygiene only; it is not a duplicate product Fast event.
- Publication uploads the exact same-run certified native artifacts; no post-certification rebuild.

## 2. G0–G16 execution map

- **G0 Exact Stable Intake:** immutable Stable identity, source SHA/fingerprint, open defects/issues, CI state and dependencies.
- **G1 Immutable Scope:** every committed scope item has stable traceability; no silent additions/removals.
- **G2 Architecture / Data Utility:** owner, consumer, provider, entitlement/rights, source of truth, reuse, freshness, retention and duplication are reviewed.
- **G3 Design / Dependency / Impact Readiness:** Impact Planner classifies affected surfaces, tests, portability, browser risk and expected CI cost.
- **G4 Development Exit:** one version-development branch, one Draft PR, clean source and scope traceability.
- **G5 Fast Qualification:** cheap exact-head syntax/format/unit/contract checks needed for the affected change.
- **G6 Integration / Medium Qualification:** affected integration and cross-module evidence.
- **G7 Data / Security / Adaptive Intelligence:** provider/data-rights/security/adaptive-intelligence evidence when applicable.
- **G8 Performance / Capacity / Stability:** load/runtime/backpressure/stability evidence when applicable.
- **G9 Cross-Module / UI / UX:** affected renderer/browser/interaction evidence.
- **G10 Pre-Freeze Qualified Candidate:** exact-head Qualified success; Release Rehearsal for CI/release changes; targeted WebKit when renderer risk warrants it.
- **G11 Immutable Release Candidate:** merged candidate is bound to exact Fast + Qualified source head and equal source fingerprint.
- **G12 Full Certification:** authoritative full certification from the immutable candidate.
- **G13 Native Packaging / Provenance:** required native packages from the candidate.
- **G14 Actual Artifact Runtime Audit:** packaged macOS Apple Silicon and Windows x64 behavior/provenance.
- **G15 Release Assurance / Promotion:** both native evidence graphs and exact artifact hashes verified.
- **Publish:** exact certified artifacts only; no rebuild.
- **G16 Adaptive Retrospective / Handoff:** current source of truth, evidence, defects, CI performance and next intake recorded.

No new top-level gates beyond G0–G16.

## 3. Impact Planner v2 change classes

The planner uses these explicit classes:

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

Classification is conservative. Unknown non-process content fails closed to full qualification.

### Lane selection

- Process/governance/CI-harness-only change: `ci-harness` Qualified lane + portability.
- Any product/mixed/uncertain change: `full` Qualified lane.
- `RENDERER_UI`: emits a WebKit-required planning signal; targeted WebKit execution is a follow-up hardening packet.
- `CI_HARNESS` or `RELEASE_TOOLING`: Release Rehearsal is required through workflow policy.

## 4. Failure taxonomy and retry matrix

- `PRODUCT_FAIL`: source behavior is wrong. Fix same branch/PR; run new exact-head Fast and Qualified as required.
- `GATE_TEST_FAIL`: test/gate correctly detects a contract mismatch or the gate itself needs governed correction. Preserve evidence; fix same branch/PR.
- `CI_HARNESS_FAIL`: CI/release tooling contract is defective. Correct same governed branch/PR; do not manufacture product identity changes.
- `INFRA_FAIL`: transient runner/network/service failure with unchanged source. Rerun affected work only.
- `EXPECTED_NOOP`: intentionally skipped/inapplicable lane or idempotent release action.
- `SUPERSEDED`: older run cancelled because a newer PR head exists.

A failure category never justifies bypassing required quality gates.

## 5. Current Phase 0 packets

### Packet A — Implement now

- Impact Planner v2 classification and self-test.
- Side-effect-free Release Rehearsal executed by canonical workflow policy.
- Current roadmap/build-plan/process/delivery overlays.
- Honest v18.6.1 reconciliation baseline without inventing missing ledger rows.
- Update `handoff/CURRENT.md` so another AI/account resumes correctly.

Expected CI behavior for Packet A: process-only → one Fast → one Qualified CI-harness/portability. No release publication because `release_identity.json` and canonical `release.yml` are intentionally unchanged.

### Packet B — Next hardening

- Immutable SHA pins for third-party GitHub Actions after upstream verification.
- Pinned Playwright + browser revision and safe cache strategy.
- Generic workflow linting.
- Per-job least-privilege permission review.

### Packet C — Risk-directed browser coverage

- Connect `webkit_required` to a focused Qualified WebKit job.
- Scope WebKit to Safari-sensitive renderer/UI changes; avoid duplicating full browser work for backend-only changes.

### Packet D — Evidence / telemetry

- Durable compact Stable evidence manifest.
- Lane runtime/queue/cache/failure-class telemetry.
- Trend CI cost and identify regressions before workflow amplification returns.

### Packet E — Renderer modularization

Incrementally extract capability owners from the large renderer/compatibility stack. Each migration requires equivalence + browser/native evidence before deleting the former owner. Never delete based on age alone.

## 6. Source and repository hygiene

A file being several days old is not a defect. Do not touch unchanged files merely to refresh GitHub dates. Remove a file only when all references/consumers/evidence needs are proven absent and protected history is not affected.

Historical certification, governance, approved reference assets and actively loaded compatibility layers are preserved until an explicit cleanup proof says otherwise.

## 7. Quality floor

Efficiency may reduce duplicate work, runner choice or irrelevant lanes; it may not reduce:

- exact-source provenance;
- deterministic tests;
- browser behavior required by the affected risk;
- data/security/rights controls;
- macOS Apple Silicon + Windows x64 Stable certification;
- same-artifact publication;
- conserved requirement traceability;
- No Execution and other permanent product boundaries.

## 8. Next product intake after Phase 0

Run fresh G0–G3 against the current reconciliation state and select the highest-value coherent v18.x slice. Provisional priority order is user-trust defects, runtime/ADR-GDI reliability, shared intelligence utility consolidation, renderer maintainability and controlled TradeInsight SHADOW integration. The current roadmap overlay is authoritative for sequencing.
