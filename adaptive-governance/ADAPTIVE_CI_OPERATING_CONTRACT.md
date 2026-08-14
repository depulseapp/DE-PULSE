# DE.PULSE — Adaptive CI Operating Contract

Status: **Permanent governing contract**  
Applies to: Adaptive Roadmap, Adaptive Build Plan, Adaptive Build Process, Adaptive Delivery Process, Build Resume Protocol  
Effective: **v18.2 and all later releases**

## Purpose

DE.PULSE CI/CD must behave as an adaptive engineering system: **observe → classify → diagnose → generalize → prevent → execute → measure → learn**. Learning is not complete when a failure is merely fixed; the generalized lesson must prevent the same failure class from recurring where practical.

This contract operates under the current canonical G0–G16 map. That map remains the default but may evolve only through the Gate Utility Test defined in `adaptive-governance/ADAPTIVE_WORK_DECOMPOSITION_CONTRACT.md`; ad-hoc workflow-local gates are forbidden.

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
- package/runtime evidence being reused after source or packaging identity changes;
- a new primary navigation surface missing from `functionality_utility_registry.json`;
- a new/materially changed capability without purpose, canonical owner, active consumer, existing-data reuse analysis, correlation targets and UI-placement disposition;
- a proposed new tab without explicit workflow/security/clarity separation justification;
- a new provider acquisition/computation path that duplicates fresh canonical evidence or equivalent in-flight work without an approved bounded exception;
- new data/persistence classes without Data Utility ownership, consumers, retention and security/sensitivity treatment;
- obsolete or superseded surfaces/jobs remaining silently after a new canonical owner is introduced.

`functionality_utility_checkpoint_gate.py` is a mandatory pre-freeze preflight from v18.2 onward. It complements, and never replaces, `data_utility_gate.py`, source-health, security, performance, UI and full certification gates.

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

Functionality-utility evidence is source-bound. If a source change adds/moves a surface, engine, job, dataset, acquisition path, metric or administrative operation, the affected functionality-utility disposition and downstream G9/G10 evidence must be revalidated. Pure implementation changes that do not alter these relationships may reuse earlier scope/architecture reasoning only when the registry and actual candidate still agree.

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
- whether the functionality utility registry still matches actual tabs/engines/jobs and whether any implementation introduced avoidable duplicate acquisition, computation, storage or user surfaces;
- what was reused, correlated, consolidated, retired or moved to drill-down/internal-only during the release;
- whether branch/workflow hygiene is restored to the preferred minimal state;
- whether both native assets were surfaced and `userDelivery` was completed.

The goal is fewer repeated failures, less duplicated product machinery and less CI churn over time without reducing coverage or release assurance.

## 10. Adaptive decomposition and parallel execution

Heavy CI responsibilities should be split into smaller independently verifiable work packages, checkpoints, shards or parallel lanes when doing so materially improves fault isolation, resumability, evidence reuse, speed or resource efficiency.

Every decomposed lane must define:

- canonical responsibility owner;
- dependencies and source/fingerprint/artifact inputs;
- PASS/FAIL/BLOCKED criteria;
- evidence output;
- downstream consumers;
- invalidation/reuse rules.

Independent lanes may run in parallel. Dependent work remains ordered. Shared setup, caches, artifacts, data acquisition and test fixtures should be reused so sharding does not multiply work unnecessarily.

Concurrency is bounded by real capacity. More jobs are not automatically better; total runner minutes, CPU, memory, browser load, artifact volume, provider/API calls and database pressure must be considered.

After a lane fails, classify the failure and rerun the smallest affected/dependent set whose evidence is invalid. Do not restart an entire heavy phase when unchanged-source evidence remains trustworthy.

Examples of appropriate decomposition include:

- role × viewport × surface-family G9 UI shards;
- independent performance/load scenarios with one aggregated G8 conclusion;
- separate macOS and Windows G13/G14 evidence;
- fast/medium/full qualification layers that reuse earlier evidence where valid;
- platform/toolchain-specific compile/runtime lanes;
- independent source/data/security checks that share one candidate fingerprint.

## 11. Gate Utility Test and canonical gate evolution

The current G0–G16 map is the default. A new canonical G-gate may be added only when a recurring/material release risk cannot be owned cleanly by an existing gate and cannot be handled more simply through checkpointing, sharding/parallel lanes or strengthening/reassigning an existing gate.

Before any gate addition, prove:

- distinct risk/responsibility;
- independent durable evidence and clear PASS/FAIL criteria;
- non-duplication with existing gates/checkpoints;
- material assurance value;
- exactly one canonical owner after consolidation;
- coordinated updates to Roadmap, Build Plan, Build Process, Delivery Process, CI, ledger/checkpoint schema and automation;
- migration clarity for in-flight/future releases while historical Stable evidence remains immutable;
- required G16 review after the new gate is used.

No workflow-local `G17`/`G18` is valid by itself.

Preferred process evolution order:

**reuse existing evidence → checkpoints → sharded/parallel lanes → strengthen/reassign an existing gate → add a new canonical gate only when materially justified.**

Canonical rules: `adaptive-governance/ADAPTIVE_WORK_DECOMPOSITION_CONTRACT.md`.

## 12. Governance-to-Implementation CI Closure

CI is responsible for proving that governance is executable rather than merely documented.

### Build Coordinator closure

The repository must converge from overlapping authoritative release workflows to one actual Build Coordinator/dependency graph. Separate workflow files are allowed only as reusable/subordinate lanes with one canonical upstream owner. No two workflows may independently own the same release mutation or authoritative gate conclusion.

Before affected release evidence is trusted, CI must identify and consolidate, subordinate, path-scope, event-scope, disable or retire overlapping formalization/docs/pre-freeze/certification/promotion workflows that create duplicate work or self-triggering storms.

### Build State Ledger closure

The Build State Ledger must be generated or deterministically reconciled from actual GitHub state. CI must reject/correct stale source commits, noncanonical evidence-state labels, outdated gate evidence, artifact mismatches or contradictory next-step/blocker data before resume/certification consumes the ledger.

A manual checkpoint may remain as a persisted representation, but it is never the source of truth and must be updated from authoritative branch/fingerprint/CI/artifact evidence.

### Canonical naming closure

`canonical_naming_gate.py` and the canonical naming registry apply to active CI/release machinery as well as documentation. Temporary numbered/ambiguous workflows, jobs and artifacts must be renamed/consolidated/retired where no compatibility need remains. One release responsibility has one canonical name and one canonical owner.

### Governance implementation status

Applicable governance requirements are tracked through:

**Governed → Implemented → Enforced → Evidenced → Delivered → Learned.**

CI may not turn `Governed` into PASS without corresponding implementation/evidence. Current-release blockers and process-hardening items are defined in `adaptive-governance/GOVERNANCE_IMPLEMENTATION_CLOSURE_CONTRACT.md`.

## Permanent rule

**DE.PULSE does not merely recover from CI failures. It classifies them, learns from them, converts useful lessons into preventative controls, and reduces the probability of recurrence while preserving full release assurance. Product growth and engineering execution follow the same discipline: understand, decompose, reuse, correlate, execute, checkpoint, evaluate, adapt, integrate and learn; add new machinery or gates only when their utility is proven.**
