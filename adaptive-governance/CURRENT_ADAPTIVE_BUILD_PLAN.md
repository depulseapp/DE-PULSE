# DE.PULSE — Current Adaptive Build Plan

**Operational overlay date:** 2026-08-19  
**Immutable Stable baseline:** `v18.6.1-stable`  
**Current engineering slice:** `v18.6.7-development` — Fresh Reconciliation, Scope Bind & Legacy Test/Gate Hygiene  
**Authority:** current execution overlay. Permanent contracts and historical release evidence remain intact.

## 1. Normal target lifecycle

`1 development branch → batch coherent branch work → 1 Draft PR → Fast → same PR Ready → Qualified → merge → Release only when release identity/release workflow requires it → Stable`

Rules:

- Build a coherent branch before opening the PR whenever practical so preparatory writes do not manufacture synchronize events.
- Never create trigger/retry/certification/promotion branches.
- Never create a second PR just to retrigger CI.
- Fix source/test/gate defects on the same branch and same PR.
- Same-SHA infrastructure failure reruns only affected work when possible.
- Main push performs hygiene only.
- Publication uses exact same-run certified native artifacts; no post-certification rebuild.

## 2. G0–G16 execution map

- **G0 Exact Stable Intake:** immutable Stable identity, source SHA/fingerprint, current requirement ledger, open defects/issues, CI state and dependencies.
- **G1 Immutable Scope:** every committed scope item has stable traceability; no silent additions/removals.
- **G2 Architecture / Data Utility:** owner, consumer, provider, entitlement/rights, source of truth, reuse, freshness, retention and duplication.
- **G3 Design / Dependency / Impact Readiness:** Impact Planner classifies affected surfaces, tests, portability, browser risk and expected CI cost.
- **G4 Development Exit:** one version-development branch, one Draft PR, clean source and scope traceability.
- **G5 Fast Qualification:** cheap exact-head syntax/format/unit/contract checks for affected risk, including conserved-ledger and legacy-test inventory integrity when applicable.
- **G6 Integration / Medium Qualification:** affected integration and cross-module evidence.
- **G7 Data / Security / Adaptive Intelligence:** provider/data-rights/security/adaptive evidence when applicable.
- **G8 Performance / Capacity / Stability:** load/runtime/backpressure/stability evidence when applicable.
- **G9 Cross-Module / UI / UX:** affected renderer/browser/interaction evidence.
- **G10 Pre-Freeze Qualified Candidate:** exact-head Qualified success; Release Rehearsal for CI/release changes; Chrome + WebKit primary evidence whenever selected risk requires browser qualification; compact CI telemetry retained.
- **G11 Immutable Release Candidate:** merged candidate bound to exact Fast + Qualified source head and equal source fingerprint.
- **G12 Full Certification:** authoritative full certification from immutable candidate.
- **G13 Native Packaging / Provenance:** required native packages from candidate.
- **G14 Actual Artifact Runtime Audit:** packaged macOS Apple Silicon and Windows x64 behavior/provenance.
- **G15 Release Assurance / Promotion:** native evidence graphs and exact artifact hashes verified.
- **Publish:** exact certified artifacts only; no rebuild.
- **G16 Adaptive Retrospective / Handoff:** current source of truth, durable release evidence, defects, CI performance and next intake.

No new top-level gates beyond G0–G16.

## 3. Phase 0 status — COMPLETE

- Packet A / PR #46 — Impact Planner v2, Release Rehearsal, governance overlays, reconciliation baseline.
- Packet B / PR #47 — immutable Fast/Qualified Action pins, dependency lock, Playwright pin, reproducibility/permission gate.
- Packet C / PR #48 — Chrome + WebKit co-primary browser policy and risk-directed execution.
- Packet D / PR #49 — durable Stable evidence, CI telemetry, amplification warnings and workflow structural lint.
- Packet E / PR #50 — Documentation capability-oriented renderer ownership with explicit legacy fallback and direct Chrome + WebKit proof.

Packet E merged to `main` at `b3ca18c14b1e53069a6736e29ad9e3b09f87bda5`.

## 4. v18.6.7 bound scope

### 4.1 Conserved requirement reconciliation

The original authority ledger is located at:

`release/v18.5.1/V17-V18-IMPLEMENTATION-RECONCILIATION.json`

Intake blob: `2a32b3f93203d61b1aca55172530652d736bbf55`. Declared tracked rows: **296**.

Governed by `adaptive-governance/V18_6_7_CURRENT_RECONCILIATION.md`.

Execution:

1. Preserve all original IDs/history and leave the v18.5.1 artifact immutable as historical evidence.
2. Fast executes `v18_5_1_v17_v18_reconciliation_gate.py` in inventory mode to prove row conservation, canonical scope alignment, ID uniqueness and status-vocabulary integrity.
3. Do not equate the 296 tracked records with 296 current defects/open items.
4. Do not inherit old statuses as current truth.
5. Freshly map exact IDs for still-relevant user-trust/runtime/provider/shared-intelligence/renderer risks to current source owners and evidence.
6. Use current dispositions `FRESH_PASS`, `REOPENED`, `NOT_IMPLEMENTED`, `INTENTIONALLY_SUPERSEDED`, `NOT_APPLICABLE`, `ROADMAP_FUTURE_SCOPE` only with supporting evidence.
7. Bind the next product slice from the verified remaining gap set; G1 freezes its exact scope.

### 4.2 Legacy test/gate inventory

Governed by:

- `adaptive-governance/LEGACY_TEST_GATE_CLEANUP_PLAN.md`
- `adaptive-governance/LEGACY_TEST_GATE_INVENTORY.md`
- `tools/ci/legacy_test_gate_inventory.py`

The live inventory scans root version-stacked executable tests/gates and records current execution/reference consumers. Classification vocabulary:

- `ACTIVE_REQUIRED`
- `ACTIVE_DUPLICATE`
- `UNREFERENCED_USEFUL`
- `HISTORICAL_EVIDENCE`
- `SAFE_TO_REMOVE`

Safety rules:

- Go root `*_test.go` defaults `ACTIVE_REQUIRED` because `go test ./...` discovers it.
- Explicit current Fast/Qualified/Release/certification-plan consumers make Python/JS gates/tests `ACTIVE_REQUIRED`.
- An unreferenced executable defaults conservatively to `UNREFERENCED_USEFUL`, never automatically `SAFE_TO_REMOVE`.
- Historical JSON/scope/audit evidence is not cosmetically moved merely to shrink the root.
- No arbitrary deletion target.

### 4.3 First safe organization wave — STAGED

Byte-for-byte test-body preservation:

| Old path | New path |
|---|---|
| `v18_6_ai_hardening_test.go` | `ai_hardening_test.go` |
| `v18_6_broad_snapshot_broker_test.go` | `broad_snapshot_broker_test.go` |
| `v18_6_documentation_access_test.go` | `documentation_access_test.go` |
| `v18_6_session_intelligence_coordinator_test.go` | `session_intelligence_coordinator_test.go` |
| `v18_6_surface_consolidation_test.js` | `tests/renderer/surface_consolidation_test.js` |
| `v18_6_documentation_access_test.js` | `tests/renderer/documentation_access_test.js` |

The Go files remain in the package root so package-private access and Go discovery remain unchanged. The renderer test bodies use repository-working-directory file reads, so moving them under `tests/renderer/` does not alter resource resolution.

Fast now executes the organized renderer tests. Impact Planner classifies `tests/renderer/` and `tests/browser/` as `RENDERER_UI`, preserving full qualification and WebKit signaling for future test organization work.

### 4.4 Retained legacy files

The current certification plan still deliberately invokes inherited/versioned evidence including v16.10/v16.11 performance gates, v17 persistence/readiness checks, v17.4 renderer acceptance, v18 scope/principal-engineer/typography gates, and focused `TestV17*`/`TestV18*` Go subsets.

These remain until:

1. their unique assertions are mapped;
2. a capability-oriented active owner exists;
3. every CI/certification consumer is updated atomically;
4. affected Fast/Qualified/Chrome/WebKit/native evidence is green as required.

### 4.5 v18.6.7 qualification

Because this slice changes Go test paths, renderer test paths and CI policy, Impact Planner should fail closed to **full qualification** with WebKit required.

Required evidence:

- exact-head Fast PASS;
- workflow policy + legacy inventory PASS;
- conserved 296-row ledger inventory PASS;
- Go formatting/vet/full suite PASS with renamed tests still discovered;
- Qualified backend full/race/randomized PASS;
- Qualified renderer/deterministic/owner regressions PASS;
- Chrome broad suite PASS;
- WebKit primary compatibility PASS;
- no release identity or `release.yml` change;
- CI telemetry retained and amplification remains normal.

### 4.6 v18.6.7 exit criteria

All must be true:

1. current reconciliation truth is documented using the located conserved ledger without rewriting historical statuses;
2. targeted root executable tests/gates have a live classification/consumer model;
3. first safe capability-oriented organization wave is complete;
4. all unique assertions in removed/moved files remain preserved or explicitly governed;
5. retained legacy files have a clear migration/deletion condition;
6. exact-head Fast and full Qualified evidence pass;
7. next coherent product slice is selected and its G0–G3/G1 scope is bound.

No Stable release is expected from this test/governance organization slice because release identity and canonical Release workflow remain untouched.

## 5. Version-visible build sequence after v18.6.7

Roadmap allocation; each exact scope freezes at G1.

1. **v18.7.0 — Runtime Reliability & Data Truth** — provisional next: `DATA DEGRADED` truth, provider/capability health, freshness SLO, coalescing/single-flight, bounded retries/circuits, backpressure/load shedding, disagreement/hysteresis, `UNKNOWN`/`ABSTAIN`, active-market load evidence. A higher-severity fresh user-trust blocker may be bundled or may lead this slice.
2. **v18.7.1 — User-Trust Reliability Closure, if needed** — focused stale/readiness/focus/state/UI reliability patch; skip if no justified patch scope.
3. **v18.8.0 — Shared Intelligence Consolidation** — canonical Scanner/Opportunity Radar acquisition/cache, Session Intelligence Coordinator, Event Intelligence lifecycle, Research evidence reuse.
4. **v18.8.1 — Renderer Modularization II** — next capability owners after Documentation, strangler/equivalence/Chrome+WebKit proof.
5. **v18.9.0 — TradeInsight SHADOW Integration** — controlled Router-only secondary/shadow capability integration with rights and utility governance.
6. **v18.9.1 — Provider Intelligence & Market-Mode Hardening** — provider usefulness/freshness/disagreement/headroom and explicit Market Mode disposition.
7. **v18.10.0 — v18 Major Closure Candidate** — zero unexplained applicable rows and full closure evidence.
8. **v18.10.1 — closure patch only if needed** — no feature expansion.
9. **v19.0.0 / v19.1.0 / v19.2.0 — Professional Data Infrastructure**.
10. **v20.0.0 / v20.1.0 / v20.2.0 — Adaptive Intelligence maturation**.

## 6. Impact Planner / browser policy

Change classes remain: `CI_HARNESS`, `RELEASE_TOOLING`, `BACKEND`, `RENDERER_UI`, `AUTH_SECURITY`, `PROVIDER_ROUTER`, `DATA_RIGHTS`, `PERSISTENCE`, `RELIABILITY_PERFORMANCE`, `CERTIFICATION_GOVERNANCE`.

Unknown non-process content fails closed to full qualification. Chrome and WebKit are co-primary. `RENDERER_UI`, full and browser candidates require WebKit. Backend/provider-only narrowed work avoids unnecessary browser runtime.

## 7. Source and repository hygiene

A file being several days old is not a defect. Delete/move only after consumer/reference/evidence and unique-assertion proof.

The desired end state is a repository whose active tests/gates are capability-oriented and understandable, while historical evidence remains traceable and the regression safety net is at least as strong as before cleanup.

## 8. Quality floor

Efficiency/cleanup may reduce duplicate files, duplicate work and root clutter; it may not reduce:

- exact-source provenance;
- deterministic tests;
- active unique regression coverage;
- Chrome + WebKit primary evidence when browser qualification is required;
- data/security/rights controls;
- macOS Apple Silicon + Windows x64 Stable certification;
- same-artifact publication;
- conserved requirement traceability;
- No Execution and other permanent product boundaries.
