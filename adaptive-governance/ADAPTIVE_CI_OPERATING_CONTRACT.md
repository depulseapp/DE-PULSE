# DE.PULSE — Adaptive CI Operating Contract

Status: **Permanent governing contract**  
Applies to: Adaptive Roadmap, Adaptive Build Plan, Adaptive Build Process, Adaptive Delivery Process, Build Resume Protocol  
Effective: **v18.2 and all later releases**

## Purpose

DE.PULSE CI/CD must behave as an adaptive engineering system: **observe → classify → diagnose → generalize → prevent → execute → measure → learn**. Learning is not complete when a failure is merely fixed; the generalized lesson must prevent the same failure class from recurring where practical.

This contract operates entirely inside canonical G0–G16. It creates no additional release gate and never weakens an existing gate.

## 1. Single logical Build Coordinator

Every release has one logical orchestration owner. Individual jobs may run in parallel when safe, but they are coordinated as one dependency graph:

**G0–G4 → G5–G10 → G11 → G12 → G13/G14 macOS + Windows → G15 → G16**

Separate workflows must not independently mutate the same release state, duplicate expensive authoritative work, or create self-triggering push loops. A release may use reusable/sharded jobs, but each responsibility has one canonical owner.

## 2. Mandatory failure classification

Every failed or terminated CI lane is classified before release state is changed:

- `PRODUCT_FAIL` — reproducible DE.PULSE product defect or violated product contract.
- `GATE_TEST_FAIL` — gate/regression logic is invalid, stale, brittle, or requires correction before trustworthy evidence can exist.
- `CI_HARNESS_FAIL` — workflow/YAML/script/transport/path/tool invocation failure unrelated to product correctness.
- `INFRA_FAIL` — hosted runner, browser, network, toolchain, capacity, timeout, or external infrastructure failure.
- `EXPECTED_NOOP` — requested mutation had nothing to change; this is successful idempotent behavior and must not become a red build.
- `SUPERSEDED` — a newer source fingerprint/RC has replaced the run; superseded evidence is retained historically but cannot block the current candidate.

Only reproducible `PRODUCT_FAIL` evidence means the candidate itself is defective. Other classes remain blocking when they prevent trustworthy evidence, but they must not be mislabeled as product failure.

## 3. Preventative learning preflight

Before expensive qualification, the build coordinator consults the permanent release-learning registry and checks known failure patterns. At minimum, preflight protects against:

- generated debris such as `__pycache__`, logs, temp files, build outputs, or transient binaries entering canonical source/transport;
- platform-specific shell assumptions when a portable Python/Go implementation is appropriate;
- unguarded `git commit` or equivalent no-op mutations;
- workflow self-trigger loops and checkpoint-only commits waking full certification;
- duplicate ownership of expensive tests or gates;
- insufficient Git history for lineage/provenance checks;
- hard-coded predecessor/current release literals where canonical release identity should be derived;
- tests mutating frozen source;
- stale jobs or stale fingerprint evidence contaminating a newer candidate;
- package/runtime evidence being reused after source or packaging identity changes.

A new recurring failure pattern must be generalized into the learning registry and, where practical, added to this preventative layer instead of being remembered only as history.

## 4. Build State Ledger / Checkpoint v2

The checkpoint is a derived, reconcilable **Build State Ledger**, not an optimistic manual status note. It must reconcile against actual GitHub branch state, canonical release identity, source fingerprint, CI runs, artifacts, immutable RC and release objects.

Required state includes:

- release/channel and incoming Stable identity;
- active branch and actual current HEAD;
- canonical source fingerprint and fingerprint state;
- G0–G16 status using `PASS`, `FAIL`, `PENDING`, `BLOCKED`, `INVALIDATED`, or `SUPERSEDED` as appropriate;
- evidence reference for every claimed PASS;
- evidence-reuse eligibility;
- failure classification and exact blocker when applicable;
- last trustworthy PASS;
- earliest required resume gate;
- exactly one next action;
- macOS package name/hash/G14 status;
- Windows package name/hash/G14 status;
- G15 promotion status;
- G16 closure status;
- `userDelivery` state: `NOT_READY`, `READY`, or `DELIVERED`;
- linked release-learning incident ID when a new generalized lesson was created.

Checkpoint data never overrides actual GitHub evidence. If the ledger and GitHub disagree, reconciliation corrects the ledger.

## 5. Metadata isolation and idempotency

Checkpoint/evidence metadata remains fingerprint-excluded. Metadata-only commits must not trigger full product qualification or native packaging unless an evidence-relevant dependency actually changed.

All mutating CI steps must be idempotent. A no-change operation is `EXPECTED_NOOP`, not failure. Workflows that commit generated governance/evidence must explicitly detect whether a diff exists before committing.

Where practical, path filters, concurrency groups and supersession rules prevent obsolete runs from consuming resources or producing misleading red/green noise.

## 6. Evidence reuse and partial resume

Exact-source PASS evidence is reused according to the Build Resume Protocol. A failed Windows G14 does not invalidate an unchanged G12 or successful macOS G14. A harness/infra repair reruns only the affected lane plus any validation necessary to prove the repair did not alter the certified candidate.

## 7. Native user-delivery invariant

A completed Stable release is not reported as fully delivered until all of the following are true:

1. G13 produced the required native packages from the certified source/RC.
2. G14 macOS Apple Silicon actual-artifact runtime audit passed.
3. G14 Windows x64 actual-artifact runtime audit passed.
4. G15 Release Assurance passed and promotion is valid.
5. Permanent release assets exist for both native platforms.
6. Package hashes/provenance are recorded.
7. The user-facing release handoff explicitly surfaces both runnable package assets and their status.
8. `userDelivery = DELIVERED` is recorded at closure.

GitHub packaging without surfacing the assets to the user is `READY`, not `DELIVERED`.

Development commits do not require native user delivery. Qualified RC/TEST deliveries and Stable releases that reach native delivery stages must follow their prescribed G13–G15 artifact requirements.

## 8. Canonical concise release status

Normal release reporting uses:

**Code | G0–G12 | macOS G13/G14 | Windows G13/G14 | G15 | G16 | User Delivery**

A completed native delivery also names the Mac and Windows assets and their hashes/provenance reference.

## 9. G16 adaptive closure

G16 must review CI incidents from the release and determine:

- which were product failures vs gate/harness/infra/no-op/superseded events;
- whether each genuinely new lesson has a root cause and generalized rule;
- which owning G0–G16 gate or preventative preflight now protects against recurrence;
- whether duplicate/obsolete workflows or checks can be consolidated or removed;
- whether branch/workflow hygiene is restored to the preferred minimal state;
- whether both native assets were surfaced and `userDelivery` was completed.

The goal is fewer repeated failures and less CI churn over time without reducing coverage or release assurance.

## Permanent rule

**DE.PULSE does not merely recover from CI failures. It classifies them, learns from them, converts useful lessons into preventative controls, and reduces the probability of recurrence while preserving full G0–G16 assurance.**
