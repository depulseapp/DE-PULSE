# DE.PULSE — Current Adaptive Delivery Process

**Operational overlay date:** 2026-08-19  
**Stable baseline:** `v18.6.1-stable`

This document is the current delivery overlay for `ADAPTIVE_DELIVERY_PROCESS.md`. Historical v18.5.x/v18.6 trial evidence remains preserved, but it does not describe the current normal delivery path.

## Delivery objective

Deliver the smallest coherent high-value change with complete affected-risk evidence, predictable runner consumption, immutable provenance, conserved requirement/test coverage and a handoff another authorized AI/account can resume from immediately.

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
- Conserved requirement ledger/source resolved; no replacement IDs or counts invented.
- Historical status distinguished from current evidence.
- Scope and permanent boundaries confirmed.
- Impact/risk classes identified.
- No duplicate architecture owner introduced.
- When tests/gates are reorganized, current consumers and unique assertions are inventoried first.
- Coherent branch work assembled before PR open where practical.

### Draft PR / Fast

- Draft communicates that source may still move.
- Fast provides cheap exact-head feedback.
- Superseded Fast work may be cancelled.
- One logical correction should normally produce one new candidate SHA and one Fast synchronization.
- Reconciliation/cleanup candidates must prove both the conserved requirement-ledger inventory and the legacy test/gate inventory before advancing.

### Ready / Qualified

- Ready means the exact source head is a candidate.
- Qualified lane is selected deterministically by Impact Planner.
- Process-only changes use harness + portability.
- Product/mixed changes use full affected product qualification.
- Moving/renaming active Go or renderer/browser test paths is not treated as harmless documentation churn; affected regression lanes must still run.
- Release Rehearsal is mandatory for CI/release-tooling risk.
- Chrome + WebKit are co-primary browser engines; both are required for full/browser qualification, and WebKit is also required for renderer/UI or WebKit-harness risk.
- Qualified records compact telemetry for queue/runtime/platform runner use, browser dependency setup/cache signals and workflow amplification.

### Merge

- Merge only after exact-head Qualified success.
- No manual trigger branch.
- Main push performs hygiene rather than duplicate product testing.
- Moved/deleted tests/gates must have their new owner/path and retained/deletion rationale represented in current governance/handoff evidence.

### Stable delivery

When release identity/canonical release workflow makes the merge release-capable:

- G11 binds merged candidate to exact source-head Fast/Qualified and fingerprint.
- G12 certifies the immutable candidate.
- macOS Apple Silicon and Windows x64 native jobs run in parallel.
- G15 verifies evidence/artifacts.
- publication uploads those exact artifacts without rebuild.
- G16 writes retrospective/handoff evidence.

A process/test-organization/governance merge that does not change release identity or the canonical Release workflow is complete after its governed PR/Qualified/hygiene path; it does not create an artificial Stable release.

## Requirement reconciliation delivery rule

The conserved historical v17/v18 ledger contains 296 tracked authority rows. That number is a conservation boundary, not an open-defect count.

Delivery must preserve:

- exact historical IDs/history;
- the historical artifact itself;
- structural row-conservation validation;
- current evidence as a separate fresh disposition layer.

No delivery may silently drop an applicable row, assume a historical PASS is current, or mark an unresolved row complete because its original version is old.

## Legacy test/gate delivery rule

For each moved/renamed/deleted version-stacked test/gate, delivery evidence must answer:

1. What currently executes/references it?
2. What unique assertions/contracts does it protect?
3. Where is the canonical capability-oriented owner after the change?
4. Were all path/working-directory/import/certification consumers updated atomically?
5. Which Fast/Qualified/Chrome/WebKit/native evidence proves no coverage loss?
6. Is the historical artifact still needed for release/certification traceability?

A file is removable only after that evidence supports `SAFE_TO_REMOVE`. Unreferenced executable tests default to `UNREFERENCED_USEFUL`, not junk.

Go package-private tests normally remain beside their package until implementation/package ownership moves with them. Renderer/browser tests may live under `tests/renderer/` / `tests/browser/`, but those paths remain `RENDERER_UI` and retain primary WebKit protection.

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

Historical requirement ledgers, release scope matrices and certification evidence remain traceable even when active executable tests/gates are later consolidated.

A retrospective Stable manifest does not redefine an already published immutable tag/binary.

## Repository hygiene after delivery

- merged development branch removed by governed hygiene;
- no retry/certification/promotion branches;
- no disposable temp files;
- no deletion of historical evidence or active compatibility assets merely because they are old;
- no unexplained legacy executable: retained version-stacked tests/gates have a classification/migration condition;
- current handoff points to the real next action.

## Delivery quality rule

Optimization may change scheduling, lane selection, caching, file organization or runner choice. It may never bypass exact-source provenance, conserved requirement traceability, unique regression coverage, Chrome + WebKit browser proof when required, deterministic truth checks, provider/data/security/rights controls, required native Stable proof, same-artifact publication, G0–G16 governance or permanent product boundaries.
