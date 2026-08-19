# DE.PULSE — Legacy Test & Gate Inventory

**Slice:** `v18.6.7-development`  
**Inventory authority:** `tools/ci/legacy_test_gate_inventory.py`  
**Cleanup contract:** `adaptive-governance/LEGACY_TEST_GATE_CLEANUP_PLAN.md`

## Why this is a live inventory instead of a one-time spreadsheet

The repository contains a large version-stacked root history spanning v14, v15, v16, v17 and v18. Some files are still active through Go package discovery or explicit certification-plan consumers; others are no longer directly invoked but may contain unique assertions or historical release knowledge. A static hand-maintained list would become stale quickly.

`tools/ci/legacy_test_gate_inventory.py` therefore inventories the current root on every governed candidate and classifies every version-stacked executable test/gate under the approved five-state vocabulary.

## Current classification rules

1. Root `*_test.go` files are `ACTIVE_REQUIRED` because the canonical Fast/Qualified suites use `go test ./...`; Go discovers them regardless of filename.
2. Versioned Python/JavaScript tests or gates explicitly named by current Fast, Qualified, Release, `certification_plan.json`, or the canonical workflow policy are `ACTIVE_REQUIRED`.
3. A versioned executable with no direct current control-plane consumer defaults conservatively to `UNREFERENCED_USEFUL`, **not** `SAFE_TO_REMOVE`. Its unique assertions/evidence must be mapped first.
4. `SAFE_TO_REMOVE` is never inferred automatically. It requires deliberate assertion/evidence review and governed removal.
5. Historical JSON/scope matrices/audits are outside the executable cleanup inference. They remain historical evidence unless a separate traceability migration proves otherwise.

## Confirmed certification-bound legacy examples

The current certification plan still deliberately consumes inherited/versioned evidence, including:

- `v17_0_cross_platform_persistence_gate.py`;
- `v17_4_renderer_test.js`;
- `v16_11_performance_stability_gate.py`;
- `v16_10_performance_gate.py`;
- `v18_0_5_renderer_test.js`;
- `v18_3_principal_engineer_gate.py`;
- `v18_documentation_typography_gate.py`;
- focused Go suites selected by test-name prefixes such as `TestV170`, `TestV171`, `TestV172`, `TestV174`, `TestV1801`, `TestV1805`, `TestV1806`, and `TestV183`.

These are not candidates for cosmetic deletion. Their behavior must first move into canonical capability-oriented owners and the certification plan must be updated atomically.

## First safe consolidation wave — v18.6.7

This wave changes organization without changing test bodies.

| Old root path | New capability-oriented path | Reason |
|---|---|---|
| `v18_6_ai_hardening_test.go` | `ai_hardening_test.go` | Go discovery is filename-agnostic; keep package-private coverage beside package. |
| `v18_6_broad_snapshot_broker_test.go` | `broad_snapshot_broker_test.go` | Same exact Go test blob, capability name. |
| `v18_6_documentation_access_test.go` | `documentation_access_test.go` | Same exact Go test blob, capability name. |
| `v18_6_session_intelligence_coordinator_test.go` | `session_intelligence_coordinator_test.go` | Same exact Go test blob, capability name. |
| `v18_6_surface_consolidation_test.js` | `tests/renderer/surface_consolidation_test.js` | Current Fast consumer redirected to a capability-oriented renderer test folder. |
| `v18_6_documentation_access_test.js` | `tests/renderer/documentation_access_test.js` | Current Fast consumer redirected to a capability-oriented renderer test folder. |

The move/rename uses the existing Git blob for each test, so the first wave is byte-for-byte content preservation rather than a behavioral rewrite.

## CI protection for the new organization

- Fast executes the two organized renderer tests from `tests/renderer/`.
- Impact Planner treats both `tests/renderer/` and `tests/browser/` as `RENDERER_UI`, forcing full qualification and primary WebKit evidence for future moves/changes in those areas.
- The live inventory fails if a first-wave legacy path reappears or a new path disappears.
- Go Fast/Qualified continues `go test ./...`, so the four renamed Go tests remain automatically discovered.

## Remaining cleanup strategy

The old pile is intentionally **not** moved wholesale into `tests/legacy/`. Doing that would only hide technical debt and can break package-private Go tests, scope-gate working-directory assumptions, certification-plan paths, or historical traceability.

Next waves should proceed by capability:

1. map unique assertions in unreferenced v16/v17 renderer tests into canonical `tests/renderer/` owners;
2. migrate useful static scope-gate assertions into current capability/CI contracts;
3. update current certification consumers before removing version-specific executable paths;
4. rename active Go files in-place by capability when filename references are absent;
5. leave historical release/scope evidence where traceability is stronger than cosmetic root cleanup.

## Exit rule

v18.6.7 does not need to eliminate every old filename. It must leave the repository in a **known state**: every version-stacked executable is dynamically classified, the first safe capability-oriented wave is complete, useful coverage is preserved, and every remaining file has a governed migration/deletion condition rather than an assumption based on age.
