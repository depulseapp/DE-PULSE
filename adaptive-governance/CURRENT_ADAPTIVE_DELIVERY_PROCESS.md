# DE.PULSE — Current Adaptive Delivery Process

**Operational overlay date:** 2026-08-19  
**Certified Stable baseline:** `v18.6.1-stable`  
**Active engineering slice:** `v18.7.0-development` — Runtime Reliability & Data Truth

This is the current delivery overlay for `ADAPTIVE_DELIVERY_PROCESS.md`. Historical trial/release evidence remains preserved but does not override the current normal path.

## Delivery objective

Deliver the smallest coherent high-value change with complete affected-risk evidence, predictable runner consumption, immutable provenance, conserved requirement/test coverage and a GitHub handoff another authorized AI/account can resume from immediately.

## Normal delivery path

`development branch → batch coherent work → Draft PR → Fast → same PR Ready → full Qualified → merge → one Release G11–G16 run when release-capable → Stable`

Expected normal counts for a clean release-capable slice:

- 1 development branch;
- 1 PR;
- 1 Fast candidate run;
- 1 Qualified candidate run;
- 1 Release G11–G16 run.

Legitimate source fixes may add Fast/Qualified attempts on the same branch/PR. Infrastructure retries never create new branches or product identities. Preparatory writes are batched before PR whenever practical to avoid unnecessary synchronize events.

## v18.7 delivery checkpoints

### Before PR

- Stable/resume evidence read.
- Exact `main` baseline resolved as `f1a9e0d0d76d4be565ac9355a09f77a33a3338a7`.
- v18.6.7 / PR #51 confirmed merged with exact-head Fast + Qualified success.
- Conserved 296-row authority ledger retained; no replacement IDs/counts invented.
- v18.7 source-owner audit completed before implementation.
- G1 scope frozen in `v18_7_0_scope.json`.
- Existing Router/freshness/backpressure/coalescing/reconciliation owners preserved.
- Runtime reliability gaps implemented only where verified.
- Release Action/browser reproducibility deferral closed.
- v18.7 release contract, current-source G12 script and candidate identity prepared on the same branch.
- Four Adaptive overlays and handoff updated before PR.
- Exact diff audited for unrelated scope.
- No existing v18.7 PR before creating the single Draft PR.

Static prevalidation never counts as CI PASS.

### Draft PR / Fast

- Open exactly one Draft PR from `v18.7.0-development` to `main`.
- Fast provides cheap exact-head feedback and records `DE.PULSE/fast-head` only after selected checks pass.
- Fast must validate workflow policy, immutable dependency/reproducibility contract, release identity, conserved-ledger integrity, Go formatting/vet/full unit suite as selected, and current renderer organization as selected by impact.
- Superseded Fast work may be cancelled.
- Any source correction stays on the same branch/PR.

### Ready / full Qualified

- Mark Ready only after exact-head Fast success.
- v18.7 requires **full** qualification; it is not a process-only slice.
- Required lanes: backend/full Go, race, randomized order, renderer/deterministic, Chrome broad behavior and WebKit primary compatibility.
- Qualified exact-head status must be `DE.PULSE/qualified-head = success`.
- Active-market synthetic reliability tests are not called PASS until executed here/backend.
- Release workflow/reproducibility changes must remain side-effect free before merge.
- Qualified telemetry records queue/runtime/platform use, browser dependency setup/cache and workflow amplification.

### Merge

- Merge only after exact current head has Fast + full Qualified success.
- No manual trigger branch or second PR.
- G11 later proves source-head → merged-candidate fingerprint equivalence.
- Main push performs hygiene rather than duplicate product qualification.

### Stable delivery

Because v18.7 changes canonical release identity / Release workflow, the merge is release-capable:

1. **G11** binds the merged candidate to exact source-head Fast + Qualified and equal fingerprint.
2. **G12** runs `release/v18.7.0/run_full_certification.sh` on the immutable candidate.
3. **G13/G14** package and audit macOS Apple Silicon + Windows x64 in parallel.
4. **G15** verifies native evidence graphs and exact artifacts.
5. **Publish** uploads those exact certified artifacts without rebuilding.
6. **G16** records final evidence/retrospective/handoff.

Only successful completion of this path creates `v18.7.0-stable`. Candidate identity text in source is not itself Stable evidence.

## v18.7 evidence semantics

### Degradation / blast radius

A degradation record must say what failed, why, which datasets/capabilities are affected, which consumers are affected, whether critical decision evidence is usable and whether affected conclusions must abstain.

`PROTECTED` means a scoped degradation exists while critical decision evidence remains usable. `DEGRADED` means required decision evidence is not trustworthy enough. `RECOVERING` means healthy evidence exists but the hysteresis window is not yet satisfied.

### Recovery

Recovery is delivered as evidence only after three consecutive healthy observations and >=5 seconds stable health. A relapse resets the streak. One successful provider request is not enough to declare overall recovery.

### Load/backpressure

A hard-full provider queue is immediately visible as `QUEUE_SATURATED`. The existing WorkloadController remains the admission/shedding owner. Optional work may shed while critical evidence remains usable; capacity SLO may still block the candidate until pressure resolves.

### Active-market proof

The v18.7 synthetic reliability test intentionally uses production BroadSnapshotBroker, WorkloadController, degradation and SLO code. It proves burst coalescing and hard capacity behavior without requiring expensive/live provider traffic. Live provider smoke is only justified when provider behavior itself is the subject under test.

## Requirement reconciliation delivery rule

The historical v17/v18 ledger contains 296 conserved authority rows. The number is a conservation boundary, not an open-defect count.

Delivery preserves:

- exact historical IDs/history;
- the original artifact;
- structural row-conservation validation;
- current evidence as a separate fresh disposition layer.

No delivery silently drops an applicable row, assumes historical PASS is current, or marks an unresolved row complete because its original release is old.

## Legacy test/gate delivery rule

For any moved/renamed/deleted executable evidence, delivery must answer:

1. current consumers/references;
2. unique assertions/contracts;
3. new capability-oriented owner/path;
4. atomic path/import/certification updates;
5. affected Fast/Qualified/Chrome/WebKit/native proof;
6. historical traceability need.

A file is removable only after evidence supports `SAFE_TO_REMOVE`. Unreferenced executable tests default to `UNREFERENCED_USEFUL`.

## Failure delivery behavior

- `PRODUCT_FAIL`: source fix, same branch/PR.
- `GATE_TEST_FAIL`: investigate source vs assertion; fix the defective side without weakening quality.
- `CI_HARNESS_FAIL`: fix harness on same branch/PR.
- `INFRA_FAIL`: bounded unchanged-SHA rerun only with a recovery signal.
- `EXPECTED_NOOP`: retain intentional evidence.
- `SUPERSEDED`: obsolete run may be cancelled/ignored for the newer exact head.

Repeated same-SHA retries without diagnosis are not acceptable.

## CI/release reproducibility delivery rule

Fast, Qualified and Release are the only workflow trust stages. All third-party Actions in those workflows are immutable-SHA pinned and lock-file governed. Qualified and Release G12 use the exact Playwright requirements file and safe pip-cache contract. The reproducibility gate rejects mutable tags, lock drift, direct unpinned Playwright install or write-permission drift.

Release pinning is closed in v18.7; it is no longer a deferred item.

## Artifact/evidence retention

- Stable truth: durable repo evidence manifest bound to tag/candidate/fingerprint/run IDs/artifact hashes/gate states after publication.
- Operational CI telemetry: compact Qualified artifact + summary.
- Historical ledgers/scope/audits/certification evidence remain traceable after active test consolidation.
- A retrospective manifest never redefines an already published immutable tag/binary.

## Repository hygiene after delivery

- merged development branch removed by governed hygiene;
- no retry/certification/promotion branches;
- no disposable temp files;
- no deletion of historical evidence/compatibility assets merely because they are old;
- no unexplained legacy executable;
- current handoff points to the real next action.

## Delivery quality rule

Optimization may change scheduling, caching, lane selection, file organization or runner choice. It may never bypass exact-source provenance, conserved requirement traceability, unique regression coverage, Chrome + WebKit proof when required, deterministic market truth, provider/data/security/rights controls, native Stable proof, same-artifact publication, G0–G16 governance or permanent product boundaries.
