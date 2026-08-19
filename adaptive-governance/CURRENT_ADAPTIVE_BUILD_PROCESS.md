# DE.PULSE — Current Adaptive Build Process

**Operational overlay date:** 2026-08-19  
**Stable baseline:** `v18.6.1-stable`

This document is the current execution overlay for `ADAPTIVE_BUILD_PROCESS.md`. Permanent contracts and historical gate evidence remain authoritative; stale older statements about the active release are historical.

## Process principles

1. **GitHub source of truth.** Resume from repository contracts/evidence, never model memory alone.
2. **One change stream.** One version-development branch and one PR for a coherent slice.
3. **Batch before PR.** Assemble coherent preparatory branch work before opening the PR whenever practical; do not manufacture Fast runs with incomplete micro-commits.
4. **Exact-head evidence.** Fast and Qualified statuses bind to the exact PR source head used by Release G11.
5. **Risk-directed validation.** Run all tests necessary for affected risk, but do not spend premium/native runners on unrelated work.
6. **Fail closed.** Unknown/mixed product impact uses full qualification.
7. **No CI event manufacturing.** No trigger/retry/certification/promotion branches or PRs.
8. **No quality-by-date cleanup.** Change/remove files only for functional/governance reasons, never because timestamps look old.
9. **Immutable Stable publication.** Certified binaries publish without rebuilding.
10. **Evidence is measurable.** Qualified records compact runtime/platform/amplification telemetry without replacing functional proof.
11. **Conserve requirement identity.** Historical requirement IDs/ledgers are never regenerated or silently dropped; current status is a fresh evidence layer, not an overwrite of historical truth.
12. **Capability-oriented test ownership.** Version-stacked tests/gates are consolidated only after consumers and unique assertions are mapped; cleanup may organize evidence but may never reduce regression coverage.

## Development flow

### Intake / G0–G3

- Resolve current Stable tag, source SHA/fingerprint and handoff.
- Resolve the conserved requirement ledger/source before inventing replacement IDs or counts.
- Reconcile open requirements against current source and distinguish historical status from current evidence.
- Freeze coherent scope with exact requirement traceability.
- Review architecture, data utility, provider/rights implications and duplicate ownership.
- Inventory legacy test/gate consumers when a slice reorganizes test/evidence paths.
- Run Impact Planner v2 classification and identify expected CI lanes.

### G4–G5

- Create one `v<next>-development` branch from current main/Stable baseline as appropriate.
- Batch coherent branch changes before opening the PR when possible.
- Open one Draft PR.
- Fast runs on PR open/synchronize/reopen and writes `DE.PULSE/fast-head` only after selected exact-head checks pass.
- For reconciliation/cleanup work, Fast must prove conserved-ledger integrity and legacy test/gate inventory integrity before the candidate advances.
- Main push does branch hygiene only.

### G6–G10

- When the exact head is ready, mark the same PR Ready.
- Qualified runs for that candidate event.
- Process-only scope uses `ci-harness` + portability.
- Product/mixed scope uses full backend/renderer/browser qualification.
- Moving/renaming active Go or renderer/browser test paths is treated as product-test risk, not cosmetic documentation work.
- Release-tooling/CI changes are exercised by side-effect-free Release Rehearsal before merge.
- Chrome and WebKit are co-primary browser engines. `full`/`browser` requires both; renderer/UI and WebKit-harness risk requires WebKit; unaffected backend/provider/process work does not.
- Qualified emits compact telemetry for queue/runtime/platform use, browser dependency setup/cache signals and PR workflow amplification.

### Failure handling

Use the canonical taxonomy:

- `PRODUCT_FAIL`: source fix on same branch/PR.
- `GATE_TEST_FAIL`: governed source/test correction on same branch/PR; never weaken a valid assertion just to pass CI.
- `CI_HARNESS_FAIL`: repair CI/release harness on same branch/PR.
- `INFRA_FAIL`: rerun failed unchanged-SHA work only.
- `EXPECTED_NOOP`: record intentional skip/idempotent outcome.
- `SUPERSEDED`: older exact-head run is obsolete and may be cancelled.

A source-changing fix produces a new exact head and therefore requires new exact-head evidence. Repeated unexplained same-SHA retries are not acceptable; stop and investigate infrastructure instead of burning runner capacity.

### G11–G16 / release-capable changes

A release-capable merged PR must satisfy:

- G11 exact Fast + Qualified statuses from source head;
- source-head → merged-candidate fingerprint equivalence;
- G12 full certification;
- G13/G14 macOS Apple Silicon + Windows x64 native/package runtime evidence in parallel;
- G15 evidence graph/artifact assurance;
- same-run publication with no rebuild;
- G16 retrospective/handoff.

Process-only/test-organization merges that do not change release identity/canonical release workflow must not manufacture a Stable release.

## Conserved requirement reconciliation process

The current conserved v17/v18 authority ledger is `release/v18.5.1/V17-V18-IMPLEMENTATION-RECONCILIATION.json`. Its 296 tracked rows are a conservation boundary, not a claim that 296 defects remain open.

For current reconciliation:

`historical ID/status → current source owner → current consumer → current evidence → fresh disposition`

Rules:

- preserve the historical ledger artifact and original IDs;
- run its structural reconciliation gate in inventory mode during v18.6.7 Fast;
- never promote an old PASS/status to current evidence automatically;
- never mechanically mark all rows reopened merely because the historical snapshot is old;
- use current dispositions only with current evidence;
- v18 Major Closure requires every applicable conserved row to be explicitly explained.

## Legacy test/gate hygiene process

Use the five-state model:

- `ACTIVE_REQUIRED`
- `ACTIVE_DUPLICATE`
- `UNREFERENCED_USEFUL`
- `HISTORICAL_EVIDENCE`
- `SAFE_TO_REMOVE`

Required migration sequence:

`inventory → identify consumer → map unique assertions → establish capability-oriented owner/path → atomically update consumers → run affected Fast/Qualified/browser/native evidence → remove old path only if proven safe`

Special rules:

- Go package-private `*_test.go` normally stays beside the package; filename consolidation is preferred before package relocation.
- `go test ./...` discovery means an old-looking Go test can still be active even with no explicit workflow filename reference.
- unreferenced Python/JS executable evidence defaults to `UNREFERENCED_USEFUL`, not dead.
- renderer/browser organization under `tests/renderer/` or `tests/browser/` remains `RENDERER_UI` and therefore keeps WebKit protection.
- historical release/scope/audit evidence is not moved merely to make the root visually clean.
- `SAFE_TO_REMOVE` is an evidence conclusion, never an automatic inventory result.

## Stable evidence process

Each published Stable should have a compact durable evidence manifest that indexes the authoritative checkpoint rather than replacing it. The manifest binds Stable tag/candidate/fingerprint, Fast/Qualified/Release run IDs, required gate states and native/G15/G16 artifact digests.

A retrospective manifest committed after Stable must explicitly state that it does not redefine the immutable tagged artifact. Impact Planner may treat only the exact `release/v*/stable-evidence-manifest.json` path as process evidence; executable release scripts remain full-qualification scope.

## CI efficiency guardrails

- Keep exactly three active workflow trust stages: Fast, Qualified, Release.
- Concurrency may cancel superseded PR-head work, never a canonical release in progress.
- Prefer Ubuntu for broad deterministic/product qualification; native runners are used only when portability/native proof is required.
- Cache only reproducible dependencies; a cache miss must never change correctness.
- Every substantial job is bounded by an intentional timeout.
- Permissions use per-job least privilege.
- Fast/Qualified third-party Actions and Playwright are pinned; Release Action pinning remains scheduled with the next genuine release-capable slice.
- Workflow structure is checked by a zero-network structural lint plus DE.PULSE semantic workflow policy.
- Qualified telemetry warnings expose abnormal workflow amplification but do not block legitimate defect-correction attempts.

## Source maintainability process

Renderer consolidation follows a strangler migration, not a rewrite. For each capability owner:

`inventory → establish canonical owner → extract → deterministic/equivalence proof → Chrome + WebKit proof → native proof if relevant → remove proven duplicate`

Test/gate consolidation follows the analogous evidence-first process defined above. Historical compatibility layers and certification assets are not deleted until runtime/reference/certification dependencies prove they are no longer required.

## Exit condition

A build/process change is complete only when the chosen CI lane passes on the exact source head, the PR is merged through the canonical lifecycle, branch hygiene is clean, current handoff is updated, telemetry/evidence contracts remain truthful, conserved requirement/test coverage remains explainable, and no hidden release or workflow amplification was created.
