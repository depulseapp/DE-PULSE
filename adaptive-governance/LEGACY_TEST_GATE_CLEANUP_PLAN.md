# DE.PULSE — Legacy Test & Gate Cleanup Plan

**Status:** Bound for `v18.6.7-development`  
**Purpose:** Reduce version-stacked root-file debt without deleting unique regression protection, release evidence, or active compatibility coverage.

## Core rule

File age and versioned naming are never deletion criteria. Every candidate must have a proven execution/reference/coverage disposition before it is moved, renamed, consolidated, archived, or deleted.

## Required classification

Every versioned root test/gate file is classified as exactly one of:

- `ACTIVE_REQUIRED` — currently executed or referenced and protects unique behavior.
- `ACTIVE_DUPLICATE` — currently executed/referenced but its assertions are fully covered by a canonical capability-oriented owner.
- `UNREFERENCED_USEFUL` — not currently executed, but contains unique regression assertions that must be migrated before removal.
- `HISTORICAL_EVIDENCE` — retained as release/certification evidence, not treated as active product code.
- `SAFE_TO_REMOVE` — no active consumer, no unique assertion, no evidence obligation, and no governed reference.

No file may move directly from unknown state to deletion.

## Organization target

Use capability ownership rather than release-number ownership for long-lived active tests.

### Go tests

Go `*_test.go` files that test package-private `package main` behavior normally remain beside the package they test. Do **not** blindly move them into a generic `tests/` subfolder because that changes Go package boundaries and can break access to unexported symbols.

Where safe, consolidate/rename versioned files into capability-oriented peers such as:

- `market_intelligence_test.go`
- `runtime_reliability_test.go`
- `historical_replay_test.go`
- `provider_health_test.go`
- `watchlist_membership_test.go`

A later package modularization may move both implementation and tests together.

### Active Python CI/release gates

Long-lived active gates should converge under capability-oriented locations such as:

- `tools/ci/gates/`
- `tools/release/`

Release-specific gates remain with release evidence when their historical identity matters. Paths are changed only after all workflow/scripts/consumers are updated and proven.

### Renderer/browser tests

Long-lived active JavaScript/Python browser tests may converge under:

- `tests/renderer/`
- `tests/browser/`

Only move them when working-directory assumptions, imports, renderer paths, CI commands and release tooling have been updated and Chrome + WebKit evidence proves equivalence.

### Historical evidence

Historical release/certification assets stay under the relevant release/evidence location or remain in place when moving them would damage traceability. Historical evidence is not churned merely to make the repository root look smaller.

## v18.6.7 execution steps

1. Inventory all root-level versioned `v*_test.go`, `v*_professional_test.go`, `v*_renderer_test.js`, `v*_scope_gate.py` and related versioned QA/gate files.
2. Record whether each file is reached by Go discovery, Fast, Qualified, Release, another script/gate, or no current execution path.
3. Build an assertion/behavior coverage map so unique protections are visible before consolidation.
4. Classify every inventoried file using the five-state model above.
5. Perform only the first **proven-safe** consolidation/move/rename wave; do not force an arbitrary deletion count.
6. Migrate unique assertions from `UNREFERENCED_USEFUL` files into canonical active tests before removing the old file.
7. Update every affected consumer/path atomically in the same branch/PR.
8. Prove no coverage loss through the affected Fast/Qualified lanes; when renderer/browser files move, require Chrome + WebKit.
9. Preserve release identity unless a genuine release-capable product change requires otherwise.
10. Leave a durable inventory for any remaining legacy files with owner, reason retained and next deletion condition.

## Deletion gate

A legacy test/gate can be deleted only when all are true:

- no current runtime/CI/release/reference consumer requires its path;
- every unique assertion/contract is mapped to an active canonical owner or is intentionally retired with documented rationale;
- historical/certification traceability is not lost;
- affected deterministic, backend, renderer and browser qualification remains green;
- the cleanup is represented in current adaptive roadmap/build-plan/handoff evidence.

## Success criteria

`v18.6.7` is successful when:

- the version-stacked root inventory is complete and truthfully classified;
- actively useful tests have capability-oriented ownership;
- dead/duplicate files are removed only with evidence;
- folders are improved where technically safe rather than cosmetically forced;
- no regression protection is lost;
- the repository has a clear remaining-debt list instead of an unexplained pile of old version files.
