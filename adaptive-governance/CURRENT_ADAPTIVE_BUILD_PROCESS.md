# DE.PULSE — Current Adaptive Build Process

**Operational overlay date:** 2026-08-19  
**Stable baseline:** `v18.6.1-stable`

This document is the current execution overlay for `ADAPTIVE_BUILD_PROCESS.md`. Permanent contracts and historical gate evidence remain authoritative; stale older statements about the active release are historical.

## Process principles

1. **GitHub source of truth.** Resume from repository contracts/evidence, never model memory alone.
2. **One change stream.** One version-development branch and one PR for a coherent slice.
3. **Exact-head evidence.** Fast and Qualified statuses must bind to the exact PR source head used by Release G11.
4. **Risk-directed validation.** Run all tests necessary for the affected risk, but do not spend premium/native runners on unrelated process-only changes.
5. **Fail closed.** Unknown/mixed product impact uses full qualification.
6. **No CI event manufacturing.** No trigger/retry/certification/promotion branches or PRs.
7. **No quality-by-date cleanup.** A file is changed or removed only for a functional/governance reason, not because its GitHub timestamp looks old.
8. **Immutable Stable publication.** Certified binaries are published without rebuilding.

## Development flow

### Intake / G0–G3

- Resolve current Stable tag, source SHA/fingerprint and handoff.
- Reconcile open issues/requirements against current source.
- Freeze a coherent scope with traceability.
- Review architecture, data utility, provider/rights implications and duplicate ownership.
- Run Impact Planner v2 classification and identify expected CI lanes.

### G4–G5

- Create one `v<next>-development` branch from current main/Stable baseline as appropriate.
- Open one Draft PR.
- Fast runs on PR open/synchronize/reopen and writes `DE.PULSE/fast-head` only after the selected exact-head checks pass.
- Main push does branch hygiene only.

### G6–G10

- When the exact head is ready, mark the same PR Ready.
- Qualified runs once for that candidate event.
- Process-only scope uses `ci-harness` + portability.
- Product/mixed scope uses full backend/renderer/browser qualification.
- Release-tooling/CI changes are exercised by side-effect-free Release Rehearsal before merge.
- Renderer/UI risk emits a targeted-WebKit planning signal; the execution lane is introduced in the dedicated browser-hardening packet.

### Failure handling

Use the canonical taxonomy:

- `PRODUCT_FAIL`: source fix on same branch/PR.
- `GATE_TEST_FAIL`: governed source/test correction on same branch/PR; never weaken a valid assertion just to pass CI.
- `CI_HARNESS_FAIL`: repair CI/release harness on same branch/PR.
- `INFRA_FAIL`: rerun failed unchanged-SHA work only.
- `EXPECTED_NOOP`: record intentional skip/idempotent outcome.
- `SUPERSEDED`: older exact-head run is obsolete and may be cancelled.

A source-changing fix produces a new exact head and therefore requires new exact-head evidence.

### G11–G16 / release-capable changes

A release-capable merged PR must satisfy:

- G11 exact Fast + Qualified statuses from the source head;
- source-head → merged-candidate fingerprint equivalence;
- G12 full certification;
- G13/G14 macOS Apple Silicon + Windows x64 native/package runtime evidence in parallel;
- G15 evidence graph/artifact assurance;
- same-run publication with no rebuild;
- G16 retrospective/handoff.

Process-only merges that do not change release identity/canonical release workflow must not manufacture a Stable release.

## CI efficiency guardrails

- Keep exactly three active workflow trust stages: Fast, Qualified, Release.
- Concurrency may cancel superseded PR-head work, never a canonical release in progress.
- Prefer Ubuntu for broad deterministic/product qualification; native runners are used only when portability/native proof is required.
- Cache only reproducible dependencies; a cache miss must never change correctness.
- Every substantial job must be bounded by an intentional timeout.
- Permissions should trend toward per-job least privilege.
- Third-party Actions and browser dependencies should be made immutable/pinned in the next reproducibility packet after upstream verification.

## Source maintainability process

Renderer consolidation follows a strangler migration, not a rewrite. For each capability owner:

`inventory → establish canonical owner → extract → deterministic/equivalence proof → browser proof → native proof if relevant → remove proven duplicate`

Historical compatibility layers are not deleted until runtime references and certification dependencies prove they are no longer active.

## Exit condition

A build/process change is complete only when the chosen CI lane passes on the exact source head, the PR is merged through the canonical lifecycle, branch hygiene is clean, current handoff is updated, and no hidden release or workflow amplification was created.
