# DE.PULSE v18.5.1 — Repository Archetype Refactor Build Plan

**Release type:** behavior-preserving patch / developer-proofing  
**Dependency:** starts only after v18.5.0 Stable promotion  
**Product scope:** no new user-facing capability; no trading/execution scope; no decision-semantic changes.

## Mission

Transform the current flat/mixed repository into a predictable archetype-style structure so a developer can navigate DE.PULSE by ownership rather than by scanning hundreds of root filenames.

This patch implements `governance/REPOSITORY_STRUCTURE_CONTRACT.md`.

## Non-negotiable rule

**Structure changes must not change product behavior.**

If a move exposes a real defect or requires behavior/semantic modification, that change is separated, documented, and requalified from the earliest affected gate.

## Workstreams

### 1. Inventory + dependency graph

Generate a machine-readable inventory of every tracked file containing:

- current path;
- category;
- owning subsystem;
- production/test/tool/config/doc/release classification;
- direct callers/importers/path references where determinable;
- disposition: KEEP / MOVE / DELETE / REVIEW;
- target path;
- risk class.

No file moves until path-dependent callers have been enumerated.

### 2. Remove dead repository clutter

Delete from the active branch after evidence/dependency verification:

- completed one-shot/discovery workflows;
- version-specific workflows superseded by current canonical workflows;
- generated/transient artifacts accidentally committed;
- obsolete duplicate gates/scripts;
- dead reference outputs no longer used by docs/tests/runtime.

Git history and immutable Stable tags remain the historical recovery mechanism.

Uncertain items go temporarily to `archive/to-review/<date>/` with a disposition manifest; the folder is not a permanent dumping ground.

### 3. Move non-Go artifacts first

Preferred destinations:

```text
tools/gates/
tools/ci/
tools/release/
tools/dev/
tests/integration/
tests/acceptance/
tests/ui/
tests/performance/
tests/degradation/
tests/fixtures/
config/policies/
config/schemas/
docs/architecture/
docs/governance/
docs/operations/
docs/releases/
assets/reference/
```

All CI/workflow/script references are changed atomically with each move group.

### 4. Go package decomposition

Move away from the broad root `package main` layout toward:

```text
cmd/depulse/
internal/app/
internal/domain/
internal/intelligence/
internal/market/
internal/providers/
internal/persistence/
internal/runtime/
internal/security/
internal/transport/
```

Decomposition rules:

- package boundaries follow ownership/invariants;
- avoid circular dependencies;
- keep interfaces narrow;
- do not export symbols merely to make a move compile when a cleaner ownership boundary is possible;
- preserve canonical shared state, provider routing, persistence ownership, fetch-once/calculate-once, bounded work/backpressure, and truthful degradation semantics;
- preserve user isolation and hosted/native parity;
- keep UI/transport DTO concerns out of core market/intelligence packages where practical.

### 5. Developer experience

Add/update:

- concise root `README` developer map;
- architecture/module ownership diagram;
- build/test commands;
- package dependency rules;
- repository structure gate preventing renewed root sprawl;
- active-workflow manifest;
- CODEOWNERS or ownership guidance if useful later.

## Test strategy — heavy qualification required

The refactor must pass all applicable existing qualification plus explicit structural-regression checks.

### Static / compile

- repository structure gate;
- broken-path/reference scan;
- `gofmt`;
- `go vet ./...`;
- full `go test -count=1 ./...`;
- CGO-disabled full test;
- Windows x64 compile;
- macOS Apple Silicon compile/package.

### Concurrency / performance

- focused and bounded race detector;
- v16.11+ performance/capacity gates;
- queue/backpressure/load-shedding suites;
- duplicate-work/single-flight assertions;
- multi-user/multi-symbol fanout;
- background-job pressure;
- long-running stability where existing harness permits.

### Data / persistence / degradation

- real PostgreSQL integration using ephemeral CI PostgreSQL;
- restart/warm-start persistence;
- PostgreSQL pressure/unavailability;
- provider failure/rate limit/fallback;
- stale evidence/source disagreement;
- degradation blast-radius containment;
- recovery hysteresis;
- truthful UNKNOWN/ABSTAIN/readiness.

### Product regression

- deterministic equivalence;
- functional workflow;
- Professional Trader/Investor acceptance;
- Extreme-30;
- randomized test order;
- browser HTTP/runtime;
- responsive/UI/renderer/accessibility;
- approved-scope traceability;
- security runtime/docs;
- data utility/architecture/adaptive-governance gates.

### Native / provenance

- macOS Apple Silicon actual package/runtime audit;
- Windows x64 actual package/runtime audit;
- source fingerprint and provenance;
- exact artifact SHA manifest;
- GitHub Release containing runnable macOS + Windows + source + checksums + certification evidence.

## Gate behavior

Use the canonical G0–G16 model; no new top-level gate is created.

Repository restructuring is evaluated mainly through:

- G0 exact v18.5 Stable baseline;
- G1 immutable no-feature scope;
- G2 architecture/package ownership;
- G3 dependency/path readiness;
- G4 source-quality/developer-proofing;
- G5–G10 full regression/medium/performance/security/UI closure;
- G11 immutable RC;
- G12 full certification;
- G13–G15 native/provenance/runtime/promotion;
- G16 cleanup/retrospective/handoff.

## Definition of done

v18.5.1 is complete only when:

1. repository root contains only intentionally allowed entry-point files/directories;
2. active workflows are current, minimal and named by purpose rather than historical accident;
3. tests/tools/config/docs are categorized;
4. production Go package ownership is understandable from the directory tree;
5. no dead/duplicate file remains without an explicit reason;
6. all structural gates and full product/native qualification pass;
7. no approved user-visible functionality or decision semantics changed;
8. a developer unfamiliar with the history can locate the major subsystem they need without searching the entire root.
