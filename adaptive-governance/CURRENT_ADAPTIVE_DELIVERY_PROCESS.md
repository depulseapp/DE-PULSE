# DE.PULSE — Current Adaptive Delivery Process

**Operational overlay date:** 2026-08-19  
**Stable baseline:** `v18.6.1-stable`

This document is the current delivery overlay for `ADAPTIVE_DELIVERY_PROCESS.md`. Historical v18.5.x/v18.6 trial evidence remains preserved, but it does not describe the current normal delivery path.

## Delivery objective

Deliver the smallest coherent high-value change with complete affected-risk evidence, predictable runner cost, immutable provenance and a handoff another authorized AI/account can resume from immediately.

## Normal delivery path

`development branch → Draft PR → Fast → same PR Ready → Qualified → merge → Release only when required → Stable`

Expected normal counts for a clean slice:

- 1 development branch;
- 1 PR;
- 1 Fast candidate run;
- 1 Qualified candidate run;
- 1 Release G11–G16 run only for a release-capable merge.

Legitimate source fixes may add Fast/Qualified attempts on the same branch/PR. Infrastructure retries do not create branches or new product identities.

## Delivery checkpoints

### Before PR

- Stable/resume evidence read.
- Scope and permanent boundaries confirmed.
- Impact/risk classes identified.
- No duplicate architecture owner introduced.

### Draft PR / Fast

- Draft communicates that source may still move.
- Fast provides cheap exact-head feedback.
- Superseded Fast work may be cancelled.

### Ready / Qualified

- Ready means the exact source head is a candidate.
- Qualified lane is selected deterministically by Impact Planner.
- Process-only changes use harness + portability.
- Product/mixed changes use full affected product qualification.
- Release Rehearsal is mandatory for CI/release-tooling risk.
- Targeted WebKit becomes mandatory when the dedicated hardening packet connects the current `webkit_required` signal to execution.

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

## Efficiency telemetry target

For each completed slice, G16 should increasingly report:

- number of branches and PRs;
- Fast/Qualified/Release run counts;
- queue time and execution time by lane;
- native-runner use;
- cache hit/miss where meaningful;
- failure category counts;
- avoided duplicate work;
- end-to-end time from first Draft PR to completion/Stable.

The desired trend is fewer duplicate events and lower irrelevant runner use, not fewer quality checks for affected risk.

## Artifact/evidence retention

Transient CI logs may expire. Stable truth must become durable through a compact release evidence manifest containing source SHA/fingerprint, canonical run IDs, artifact hashes, required gate states and tool/dependency versions. This is a planned hardening packet; existing immutable release evidence is preserved until then.

## Repository hygiene after delivery

- merged development branch removed by governed hygiene;
- no retry/certification/promotion branches;
- no disposable temp files;
- no deletion of historical evidence or active compatibility assets merely because they are old;
- current handoff points to the real next action.

## Delivery quality rule

Optimization may change scheduling, lane selection, caching or runner choice. It may never bypass exact-source provenance, affected functional evidence, deterministic truth checks, provider/data/security/rights controls, required native Stable proof, same-artifact publication, G0–G16 governance or permanent product boundaries.
