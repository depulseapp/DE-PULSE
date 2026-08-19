# DE.PULSE — Current Adaptive Build Process

**Operational overlay date:** 2026-08-19  
**Certified Stable baseline:** `v18.6.1-stable`  
**Active engineering slice:** `v18.7.0-development` — Runtime Reliability & Data Truth

This is the current execution overlay for `ADAPTIVE_BUILD_PROCESS.md`. Permanent contracts and historical evidence remain authoritative; stale statements about older active releases are historical only.

## Process principles

1. **GitHub source of truth.** Resume from repository contracts/evidence, never model memory alone.
2. **One change stream.** One version-development branch and one PR for a coherent slice.
3. **Batch before PR.** Assemble coherent source/governance/release preparation before opening the PR whenever practical.
4. **Exact-head evidence.** Fast and Qualified bind to the exact PR source head later consumed by G11.
5. **Risk-directed validation.** Run all evidence required by affected risk; avoid unrelated premium/native work.
6. **Fail closed.** Unknown/mixed product impact uses full qualification; insufficient market evidence becomes `UNKNOWN/ABSTAIN`, not false health.
7. **No CI event manufacturing.** No trigger/retry/certification/promotion branches or second PRs.
8. **No cleanup by age.** Change/remove files only for functional/governance reasons.
9. **Immutable Stable publication.** Certified binaries publish without rebuilding.
10. **Conserve requirement identity.** Historical requirement IDs/ledgers are never regenerated or silently dropped.
11. **Reuse canonical owners.** A roadmap bullet does not justify a duplicate subsystem when current source already owns the concern.
12. **Capability-oriented test ownership.** Consolidation follows consumer/assertion evidence and cannot reduce coverage.

## v18.7 G0–G3 intake rule

The frozen v18.7 scope is governed by:

- `v18_7_0_scope.json`
- `v18_7_0_g0_g3_contract.json`
- `adaptive-governance/V18_7_0_RUNTIME_RELIABILITY_AUDIT.md`

Source audit precedes implementation. Every reliability concern is classified as:

- `VERIFIED_GAP_IMPLEMENT`;
- `INHERITED_HARDEN_ONLY`;
- `INHERITED_PROVEN_REQUALIFY`;
- `OUT_OF_SCOPE`.

For v18.7, existing Smart Provider Router, freshness, workload/backpressure, request telemetry, snapshot coalescing and provider reconciliation remain their canonical owners. No parallel retry/health/freshness/degradation/reconciliation engine may be introduced.

## Development flow

### G0–G3

- Resolve certified Stable, exact current `main`, handoff and checkpoints.
- Resolve the conserved 296-row v17/v18 ledger and current reconciliation overlay.
- Audit source owners before deciding what is missing.
- Freeze exact scope and permanent boundaries at G1.
- Classify impact, provider/data-rights implications, duplication risk and expected CI lanes.

### G4–G5

- Create one `v18.7.0-development` branch from exact `main` baseline.
- Batch coherent branch changes before PR.
- Prevalidate release identity, workflow contracts and exact diff statically.
- Open one Draft PR only after coherent preparation.
- Fast runs on PR open/synchronize/reopen and writes `DE.PULSE/fast-head` only after exact-head checks pass.
- Main push performs branch hygiene only.

### G6–G10

- Fix legitimate defects on the same branch/PR.
- Mark the same PR Ready only after Fast is green.
- v18.7 requires **full Qualified** because it changes backend runtime reliability, tests, release identity and Release workflow.
- Full Qualified requires backend/full Go, race, randomized order, renderer/deterministic, Chrome and WebKit primary evidence.
- Release-tooling risk must remain side-effect free before merge; no publication occurs from Fast/Qualified.
- G10 requires conserved-ledger integrity and current release/reproducibility contracts.

### Failure handling

Use the canonical taxonomy:

- `PRODUCT_FAIL`: fix source on same branch/PR.
- `GATE_TEST_FAIL`: correct defective source/test/gate; never weaken a valid assertion just to pass.
- `CI_HARNESS_FAIL`: repair harness on same branch/PR.
- `INFRA_FAIL`: bounded same-SHA retry only when a recovery signal exists.
- `EXPECTED_NOOP`: retain intentional skip/idempotent evidence.
- `SUPERSEDED`: older exact-head work is obsolete.

A source-changing fix creates a new exact head and requires new exact-head evidence. Repeated unexplained same-SHA retries are not acceptable.

## v18.7 reliability implementation rule

### Degradation

- Raw degradation remains deterministic.
- Canonical machine-readable reason lives in `ReasonCode`; concise `Code` remains renderer/API compatibility text.
- Affected dataset/capability and consumer blast radius must be explicit.
- If critical decision evidence is unusable, affected conclusions fail closed with `Abstain=true`.
- If no narrower cause is proven, use truthful `UNKNOWN`; absence of diagnosis is never health.

### Recovery

- Degradation appears immediately.
- Recovery is held as `RECOVERING` until 3 consecutive healthy observations and >=5 seconds stability.
- Relapse resets the streak.
- The existing `RuntimeSLOTracker` owns recovery state; no new recovery service.

### Load/backpressure

- Existing WorkloadController tiers/capacities/queues/reserved critical headroom remain canonical.
- A hard-full provider queue is immediately `QUEUE_SATURATED` even before queue age exceeds warning thresholds.
- Optional/background work sheds before protected critical work.

### Provider/capability health

- Human-readable capability registry remains descriptive.
- Smart Provider Router v2 remains executable health/routing truth, combining configuration, entitlement/capability state, circuits, telemetry, preferred/serving route, fallback reason and recovery.
- A successful half-open provider probe may close the provider/capability circuit; runtime recovery hysteresis prevents a single probe from falsely flipping user-facing overall health.

### Duplicate work

- BroadSnapshotBroker remains snapshot reuse/coalescing owner.
- v18.7 adds higher-fanout active-market proof; it does not create a second broker.

## Release identity process

v18.7 is a genuine release-capable slice, so release identity is prepared on the same development branch before PR qualification.

`release_identity.json` may identify the candidate as `18.7.0` / `STABLE` for package identity, but **Stable status is not achieved by editing identity files**. Certified Stable remains `v18.6.1-stable` until G11–G16 finishes successfully.

v18.7 uses `renderer/release-identity-v18.7.0.js` as a last-loaded identity overlay. This generalizes the small-overlay pattern so the legacy renderer monolith is not rewritten solely for two version constants. Existing behavioral extensions remain loaded underneath.

Historical `certification_plan.json` / `ci_pipeline_plan.json` are conserved registries. `release/v18.7.0/release_contract.json` declares their legacy versions; current v18.7 execution authority is the three canonical GitHub workflows + current release contract + `run_full_certification.sh`.

## G11–G16 / release-capable merge

A merged v18.7 PR must satisfy:

- G11 exact Fast + full Qualified statuses from the source head;
- source-head → merged-candidate fingerprint equivalence;
- G12 current-source certification via `release/v18.7.0/run_full_certification.sh`;
- G13/G14 macOS Apple Silicon + Windows x64 package/runtime evidence in parallel;
- G15 exact evidence graph/artifact assurance;
- same-run publication with no rebuild;
- G16 retrospective/handoff.

The G12 script re-runs current reliability, Go/race/randomized, renderer/deterministic, Chrome, governance and reconciliation contracts. WebKit remains a mandatory **pre-merge exact-head Qualified** proof verified by G11; G12 does not add a redundant fourth workflow/lane.

## Conserved requirement reconciliation

The authoritative historical ledger remains `release/v18.5.1/V17-V18-IMPLEMENTATION-RECONCILIATION.json` with 296 conserved rows.

Current evaluation always follows:

`historical ID/status → current source owner → current consumer → current evidence → fresh disposition`

Never assume old PASS is current evidence. Never mechanically reopen all rows because the snapshot is old. v18 Major Closure must explicitly explain every applicable row.

## Legacy test/gate hygiene

Five-state model remains:

- `ACTIVE_REQUIRED`
- `ACTIVE_DUPLICATE`
- `UNREFERENCED_USEFUL`
- `HISTORICAL_EVIDENCE`
- `SAFE_TO_REMOVE`

Migration order:

`inventory → consumer → unique assertions → capability-oriented owner/path → atomic consumer update → affected Fast/Qualified/Chrome/WebKit/native proof → remove only if safe`

## CI efficiency guardrails

- Exactly three workflow trust stages: Fast, Qualified, Release.
- Fast/Qualified/Release third-party Actions are immutable-SHA pinned and lock-file governed.
- Qualified and Release G12 use the same exact Playwright requirements file + safe pip cache contract.
- Release pinning is **closed in v18.7**, not deferred.
- Concurrency may cancel superseded PR-head work, never canonical Release in progress.
- Prefer Ubuntu for broad deterministic qualification; native runners only where required.
- Cache only reproducible dependencies; cache misses never change correctness.
- Every substantial job is timeout-bounded.
- Permissions use per-job least privilege.
- Structural lint + semantic workflow policy enforce the three-workflow model.
- Telemetry warnings expose amplification but never justify skipping required quality checks.

## Exit condition

v18.7 development is complete only when the one PR has exact-head Fast + full Qualified success, G10 reconciliation/release evidence is truthful, the exact head merges, and one G11–G16 run certifies/packages/audits/publishes the same candidate. Only then update Stable handoff/tag/evidence to v18.7.0.
