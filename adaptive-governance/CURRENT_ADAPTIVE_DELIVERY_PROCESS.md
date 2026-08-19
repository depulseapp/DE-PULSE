# DE.PULSE — Current Adaptive Delivery Process

**Operational overlay date:** 2026-08-19  
**Stable baseline:** `v18.6.1-stable`

This document is the current delivery overlay for `ADAPTIVE_DELIVERY_PROCESS.md`. Historical v18.5.x/v18.6 trial evidence remains preserved, but it does not describe the current normal delivery path.

## Delivery objective

Deliver the smallest coherent high-value change with complete affected-risk evidence, predictable runner consumption, immutable provenance and a handoff another authorized AI/account can resume from immediately.

## Normal delivery path

`development branch → batch coherent work → Draft PR → Fast → same PR Ready → Qualified → merge → Release only when required → Stable`

Expected normal counts for a clean slice:

- 1 development branch;
- 1 PR;
- 1 Fast candidate run;
- 1 Qualified candidate run;
- 1 Release G11–G16 run only for a release-capable merge.

Legitimate source fixes may add Fast/Qualified attempts on the same branch/PR. Infrastructure retries do not create branches or new product identities. Preparatory file writes should be batched before opening the PR whenever practical so they do not create avoidable `synchronize` runs.

## Delivery checkpoints

### Before PR

- Stable/resume evidence read.
- Scope and permanent boundaries confirmed.
- Impact/risk classes identified.
- No duplicate architecture owner introduced.
- Coherent branch work assembled before PR open where practical.

### Draft PR / Fast

- Draft communicates that source may still move.
- Fast provides cheap exact-head feedback.
- Superseded Fast work may be cancelled.
- One logical correction should normally produce one new candidate SHA and one Fast synchronization.

### Ready / Qualified

- Ready means the exact source head is a candidate.
- Qualified lane is selected deterministically by Impact Planner.
- Process-only changes use harness + portability.
- Product/mixed changes use full affected product qualification.
- Release Rehearsal is mandatory for CI/release-tooling risk.
- Chrome + WebKit are co-primary browser engines; both are required for full/browser qualification, and WebKit is also required for renderer/UI or WebKit-harness risk.
- Qualified records compact telemetry for queue/runtime/platform runner use, browser dependency setup/cache signals and workflow amplification.

### Merge

- Merge only after exact-head Qualified success.
- No manual trigger branch.
- Main push performs hygiene rather than duplicate product testing.

### Stable delivery

When release identity/canonical release workflow makes the merge release-capable:

- G11 binds merged candidate to exact source-head Fast/Qualified and fingerprint.
- G12 certifies the immutable candidate.
- macOS Apple Silicon and Windows x64 native jobs run in parallel.
- G15 verifies evidence/artifacts.
- publication uploads those exact artifacts without rebuild.
- G16 writes retrospective/handoff evidence.

A process-only governance merge is complete after its governed PR/Qualified/hygiene path; it does not create an artificial Stable release.

## Failure delivery behavior

Every failure must be classified before choosing a retry action:

- `PRODUCT_FAIL`: fix source, same branch/PR.
- `GATE_TEST_FAIL`: investigate assertion vs contract; fix the defective side without weakening quality.
- `CI_HARNESS_FAIL`: fix harness on same branch/PR.
- `INFRA_FAIL`: rerun unchanged-SHA failed work only.
- `EXPECTED_NOOP`: retain as intentional evidence.
- `SUPERSEDED`: obsolete run may be cancelled/ignored in favor of newer exact head.

Repeated same-SHA infrastructure retries are bounded: retry only when there is a reasonable recovery signal; otherwise stop and investigate rather than consume Actions in a loop.

## Efficiency telemetry

For each Qualified candidate, retain compact operational evidence for:

- exact candidate SHA and selected lane;
- Fast/Qualified/Release counts on the PR branch;
- queue time and execution time per completed job;
- Linux/macOS/Windows runner minutes;
- Chrome/WebKit dependency setup duration and pip cache-hit state where applicable;
- amplification warnings above conservative thresholds.

Telemetry records platform consumption, not invented dollar values. Actual currency cost remains the billing system's authority. Telemetry warnings are diagnostic: they surface event amplification but never justify skipping required quality checks.

## Artifact/evidence retention

Transient CI logs may expire. Two layers are retained:

1. **Stable truth:** repo-durable compact Stable evidence manifest bound to the authoritative release evidence checkpoint, including source/fingerprint, canonical run IDs, artifact hashes and gate states.
2. **Operational CI telemetry:** compact Qualified JSON artifact retained for 30 days plus a human-readable job summary.

A retrospective Stable manifest does not redefine an already published immutable tag/binary.

## Repository hygiene after delivery

- merged development branch removed by governed hygiene;
- no retry/certification/promotion branches;
- no disposable temp files;
- no deletion of historical evidence or active compatibility assets merely because they are old;
- current handoff points to the real next action.

## Delivery quality rule

Optimization may change scheduling, lane selection, caching or runner choice. It may never bypass exact-source provenance, affected functional evidence, Chrome + WebKit browser proof when required, deterministic truth checks, provider/data/security/rights controls, required native Stable proof, same-artifact publication, G0–G16 governance or permanent product boundaries.
