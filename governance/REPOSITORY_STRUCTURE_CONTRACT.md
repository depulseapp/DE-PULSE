# DE.PULSE Repository Structure Contract

**Status:** Permanent engineering contract from v18.5 Major Closure onward  
**Goal:** Make the repository navigable, predictable, and developer-proof in the same spirit as a well-defined Maven archetype, without changing DE.PULSE product behavior.

## Principles

1. The repository root is an entry point, not a storage area.
2. Production code, tests, tooling, configuration, documentation, release evidence, and assets must have distinct homes.
3. Historical release material must not remain mixed with active build machinery.
4. Temporary one-shot CI/discovery workflows must be retired after use.
5. Files are deleted when Git history already provides sufficient historical preservation; an `archive/` area is reserved only for material that remains operationally useful.
6. Moving Go files across package boundaries is an architectural refactor and must not be disguised as housekeeping.
7. Structural cleanup must be behavior-preserving and heavily regression-tested.

## Target archetype

```text
DE-PULSE/
├── cmd/
│   └── depulse/                    # executable entry point (post-v18.5 refactor)
├── internal/
│   ├── app/                        # application orchestration/state
│   ├── domain/                     # core domain models and invariants
│   ├── intelligence/               # adaptive/research/decision-support logic
│   ├── market/                     # market state, scanners, readiness, catalysts
│   ├── providers/                  # provider adapters/router/subscriptions
│   ├── persistence/                # repository/database/PostgreSQL
│   ├── runtime/                    # lifecycle, scheduling, load/backpressure
│   ├── security/                   # auth/RBAC/session/security policy
│   └── transport/                  # HTTP/API handlers and DTO boundaries
├── renderer/                       # desktop/web renderer assets
├── tests/
│   ├── integration/
│   ├── acceptance/
│   ├── ui/
│   ├── performance/
│   ├── degradation/
│   └── fixtures/
├── tools/
│   ├── gates/                      # Python/other qualification gates
│   ├── ci/                         # CI orchestration helpers
│   ├── release/                    # packaging, provenance, release identity
│   └── dev/                        # local developer utilities
├── config/
│   ├── policies/                   # JSON policy/registry inputs
│   └── schemas/
├── docs/
│   ├── architecture/
│   ├── governance/
│   ├── operations/
│   └── releases/
├── release/<version>/              # immutable/version-scoped release evidence
├── assets/
│   └── reference/                  # approved screenshots/reference media
├── .depulse-certification/         # resumable machine certification state
├── .github/workflows/              # active workflows only
├── README.md
├── VERSION.txt
├── go.mod
└── go.sum
```

## Root allowlist after Repository Archetype Refactor

The normal repository root should contain only:

- `README.md`, `VERSION.txt`, `go.mod`, `go.sum`, `.gitignore` and similarly universal repo metadata;
- top-level architecture directories listed above;
- no version-specific test scripts, qualification gates, screenshots, temporary workflows, or policy JSON files.

## v18.5 safe cleanup scope

Before v18.5 Stable freeze:

- retire temporary v18.5 discovery/one-shot workflows after their evidence is consumed;
- keep only workflows still required to certify/promote v18.5;
- record this structure contract and the post-Stable migration plan;
- do **not** move production Go files across package boundaries;
- do **not** perform broad path moves of gates/tests until their callers/import/path assumptions have been inventoried and mechanically updated.

This preserves the validity of current Major Closure work while stopping additional repository sprawl.

## v18.5.1 — Repository Archetype Refactor

v18.5.1 is reserved as a behavior-preserving repository/source-structure patch immediately after v18.5 Stable.

### Phase A — Inventory and dependency graph

Classify every tracked file as one of:

- production runtime;
- Go unit/package test;
- integration/degradation/performance test;
- UI/renderer/acceptance test;
- CI/gate/release tool;
- config/policy/schema;
- documentation/governance;
- release/certification evidence;
- asset/reference;
- obsolete/deletable.

For every non-Go script, identify callers and relative-path assumptions before moving it.

### Phase B — Non-Go deterministic relocation

Move path-stable categories first:

- `*_gate.py` and qualification helpers → `tools/gates/` or `tools/ci/`;
- release/version/provenance helpers → `tools/release/`;
- JS/Python acceptance/UI harnesses → `tests/acceptance/` or `tests/ui/`;
- performance/degradation harnesses → corresponding `tests/` subtree;
- policy/registry JSON → `config/policies/`;
- screenshots/reference media → `assets/reference/`;
- long-lived governance/docs → `docs/` hierarchy where compatible;
- obsolete version-specific workflows → delete from the active branch after validation; Git history/tags retain historical copies.

All workflow/script references must be updated in the same atomic change set.

### Phase C — Go package decomposition

Only after Phase B is green, migrate from the current broad `package main` layout toward `cmd/depulse` + `internal/...` packages.

Rules:

- split by ownership and dependency direction, not arbitrary file count;
- avoid cyclic packages;
- keep provider, persistence, adaptive-intelligence, runtime/backpressure, and transport boundaries explicit;
- preserve fetch-once/calculate-once and canonical-state ownership;
- preserve No Execution and US Equities Processing boundaries;
- no behavior/API/decision-semantic change is permitted as part of structural movement;
- any unavoidable semantic change is treated as a separate defect/change and requalified from the earliest invalidated gate.

## Mandatory validation for structural moves

The refactor is not accepted based on compilation alone. At minimum it must pass:

1. repository-structure/path-reference gate;
2. `gofmt`, `go vet`, full `go test ./...`;
3. focused race detector plus bounded/full race coverage appropriate to the change;
4. CGO-disabled fallback;
5. real PostgreSQL integration/readiness path;
6. ADR-GDI provider failure, stale evidence, backpressure/load shedding, fanout, restart, recovery-hysteresis and truthful UNKNOWN/ABSTAIN coverage;
7. deterministic equivalence and professional acceptance;
8. browser HTTP/responsive/UI regression;
9. source-quality/architecture/data-utility/adaptive-governance gates;
10. macOS Apple Silicon and Windows x64 packaging/runtime certification;
11. exact source fingerprint/provenance comparison proving the refactor did not silently alter approved product scope.

## Deletion policy

Do not keep a permanent `delete/` dumping ground. For files confirmed obsolete:

- delete them from the active branch after dependency verification;
- rely on Git history and immutable Stable tags for historical recovery.

If deletion confidence is not yet sufficient, place the item temporarily under `archive/to-review/<date>/` with a manifest containing original path, reason, owner, and planned disposition. Temporary archive items must be resolved at the next Major Closure.

## Definition of done

A new developer should be able to open the repository and answer, within a few minutes:

- where production code lives;
- where provider/persistence/intelligence/runtime code lives;
- where tests live and how they are grouped;
- where CI/release gates live;
- where policies/config live;
- where release evidence lives;
- which workflows are currently active;
- how to build, test, certify and package DE.PULSE.

If those answers require scanning the root for filenames, the structure contract is not yet satisfied.
