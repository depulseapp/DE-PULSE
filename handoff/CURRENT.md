# DE.PULSE — Current Authoritative Handoff

**SUPERSEDES ALL PRIOR CHAT HANDOFFS.** GitHub is authoritative; chat history and temporary workspaces are advisory only.

**Release:** `v18.6.1`  
**Active branch:** `main`  
**Certified Stable:** `v18.6.1-stable`  
**Stable target / certified merged candidate:** `42e8432f7530ae39cbfd6ceb0b0bd5f6311dc5cc`  
**Qualified v18.6.1 source head:** `5c3fae486f3e4b4a39a0b1d549916aea9e1295fd`  
**Certified source fingerprint:** `b01c14e1d54b736785eab6c03407801c527edd7769ff6f3d41fd4b20dabebd75`  
**Build ID:** `v18.6.1-stable-20260819`  
**Canonical v18.6.1 certification/publication run:** `32279232665` (Release #24)  
**Repository:** `depulseapp/DE-PULSE`  
**Latest merged engineering main before this branch:** `b3ca18c14b1e53069a6736e29ad9e3b09f87bda5` (PR #50 / Packet E)  
**Engineering branch:** `v18.6.7-development`  
**Current slice:** Fresh Reconciliation, Scope Bind & Legacy Test/Gate Hygiene  
**Current PR:** resolve live GitHub PR state for `v18.6.7-development`; do not rely on a hard-coded PR number in this handoff  
**Last updated:** 2026-08-19 America/Vancouver

`Release` and `Active branch` intentionally mirror the immutable build-resume checkpoint (`v18.6.1` / `main`). `Engineering branch` records in-flight engineering without mutating Stable checkpoint identity.

## Resume rule

Read `AGENTS.md` or `CLAUDE.md`, `governance/CI-EFFICIENCY-CONTRACT.md`, this handoff, `release_identity.json`, both `.depulse-certification/resume/` checkpoints, current GitHub PR/issue state and immutable `v18.6.1-stable`. Never resume from model memory alone.

Also read:

- `adaptive-governance/CURRENT_ADAPTIVE_ROADMAP.md`
- `adaptive-governance/CURRENT_ADAPTIVE_BUILD_PLAN.md`
- `adaptive-governance/CURRENT_ADAPTIVE_BUILD_PROCESS.md`
- `adaptive-governance/CURRENT_ADAPTIVE_DELIVERY_PROCESS.md`
- `adaptive-governance/V18_6_7_CURRENT_RECONCILIATION.md`
- `adaptive-governance/LEGACY_TEST_GATE_CLEANUP_PLAN.md`
- `adaptive-governance/LEGACY_TEST_GATE_INVENTORY.md`
- `release/v18.5.1/V17-V18-IMPLEMENTATION-RECONCILIATION.json`
- `release/v18.6.1/stable-evidence-manifest.json`

## Immutable v18.6.1 Stable evidence

- Final exact-head Fast: run `32279055139` · PASS.
- Final exact-head Qualified: run `32279113032` · PASS.
- Product candidate full Qualified: run `32276304863` · backend/race/randomized, renderer, deterministic 2403/2403 and browser/global-remove coverage PASS.
- Canonical Release #24: run `32279232665` · G11–G16 PASS including macOS Apple Silicon and Windows x64 native evidence and same-run publication.
- Stable tag: `v18.6.1-stable` → `42e8432f7530ae39cbfd6ceb0b0bd5f6311dc5cc`.

Permanent boundaries remain unchanged: U.S. Equities Processing, No Execution, deterministic Day/Swing/Long truth, Smart Provider Router ownership, provider/Market-Mode rules, GLD/SLV/USO tradable exceptions and Adaptive Intelligence governance.

## Phase 0 hardening — COMPLETE

- Packet A / PR #46 — Impact Planner v2, Release Rehearsal and current adaptive overlays.
- Packet B / PR #47 — reproducible Fast/Qualified dependencies, Action pins, Playwright pin, permissions/reproducibility gate.
- Packet C / PR #48 — Chrome + WebKit co-primary browser coverage and risk-directed routing.
- Packet D / PR #49 — durable Stable evidence, CI observability/telemetry and structural workflow lint.
- Packet E / PR #50 — Documentation capability-oriented renderer owner foundation with explicit legacy fallback, owner/decorator metadata and direct Chrome + WebKit proof.

Packet E merged to `main` at `b3ca18c14b1e53069a6736e29ad9e3b09f87bda5`.

## v18.6.7 reconciliation state

The historical conserved requirement ledger is conclusively located at:

`release/v18.5.1/V17-V18-IMPLEMENTATION-RECONCILIATION.json`

Intake blob: `2a32b3f93203d61b1aca55172530652d736bbf55`. The ledger declares **296 tracked authority rows**. This means 296 conserved records, not 296 current defects/open items.

The old v18.6.1 reconciliation file remains historical evidence of the earlier uncertainty and is not rewritten. `adaptive-governance/V18_6_7_CURRENT_RECONCILIATION.md` is the current overlay.

Fast now executes `v18_5_1_v17_v18_reconciliation_gate.py` in inventory mode so schema, canonical scope alignment, 296-row conservation, immutable IDs and status vocabulary cannot silently drift. Historical v18.5.1 statuses are not treated as current evidence.

## Legacy Test & Gate Hygiene

Governed by:

- `adaptive-governance/LEGACY_TEST_GATE_CLEANUP_PLAN.md`
- `adaptive-governance/LEGACY_TEST_GATE_INVENTORY.md`
- `tools/ci/legacy_test_gate_inventory.py`

Classification model:

- `ACTIVE_REQUIRED`
- `ACTIVE_DUPLICATE`
- `UNREFERENCED_USEFUL`
- `HISTORICAL_EVIDENCE`
- `SAFE_TO_REMOVE`

`SAFE_TO_REMOVE` is never automatically inferred. Go `*_test.go` remains active through `go test ./...`; versioned Python/JS executables with direct current CI/certification consumers remain active; unreferenced executable evidence defaults conservatively to `UNREFERENCED_USEFUL` until its assertions are mapped.

### First safe capability-oriented wave — STAGED

Exact test-body preservation:

- `v18_6_ai_hardening_test.go` → `ai_hardening_test.go`
- `v18_6_broad_snapshot_broker_test.go` → `broad_snapshot_broker_test.go`
- `v18_6_documentation_access_test.go` → `documentation_access_test.go`
- `v18_6_session_intelligence_coordinator_test.go` → `session_intelligence_coordinator_test.go`
- `v18_6_surface_consolidation_test.js` → `tests/renderer/surface_consolidation_test.js`
- `v18_6_documentation_access_test.js` → `tests/renderer/documentation_access_test.js`

The Go tests remain beside package `main`; the renderer tests were verified to read production resources from repository working-directory paths, so moving the scripts does not change their file resolution.

Fast consumes the new renderer paths. Impact Planner now treats `tests/renderer/` and `tests/browser/` as `RENDERER_UI`, preserving full Qualified + primary WebKit signaling.

The current certification plan still explicitly consumes inherited/versioned v16/v17/v18 gates and focused test prefixes. Do not move/delete those merely to clean the root; migrate assertion ownership and certification consumers first.

## v18.6.7 expected qualification

Because the branch changes Go test paths, renderer test paths and CI policy, the intended qualification is fail-closed **full**:

- Fast exact head, including workflow policy, legacy test/gate inventory and conserved ledger integrity;
- Go formatting/vet/full suite;
- Qualified backend full/race/randomized;
- renderer/deterministic/owner regressions;
- Chrome broad behavior;
- WebKit primary compatibility;
- telemetry retained.

No Stable Release is expected because neither `release_identity.json` nor `.github/workflows/release.yml` is being changed.

## Version-visible sequence after v18.6.7

Subject to G1 scope freeze:

1. **v18.7.0 — Runtime Reliability & Data Truth** — provisional next unless a higher-severity fresh user-trust blocker leads/is bundled.
2. **v18.7.1 — User-Trust Reliability Closure** — only if needed.
3. **v18.8.0 — Shared Intelligence Consolidation**.
4. **v18.8.1 — Renderer Modularization II**.
5. **v18.9.0 — TradeInsight SHADOW Integration**.
6. **v18.9.1 — Provider Intelligence & Market-Mode Hardening**.
7. **v18.10.0 — v18 Major Closure Candidate**.
8. **v18.10.1 — closure patch only if needed**.
9. **v19.x — Professional Data Infrastructure**.
10. **v20.x — Adaptive Intelligence**.

## Repository/source-age rule

Do not modify or delete files because GitHub says they are several days old. Historical release evidence, governance and actively loaded compatibility assets stay until references/consumers/evidence requirements prove safe removal.

## External/open governance item

Repository `main` branch protection/ruleset remains disabled as of the latest checked state. Do not report it as fixed unless GitHub confirms an authorized branch-protection/ruleset change.

## Exactly one next action

Resolve the live PR state for `v18.6.7-development`. If no PR exists, open exactly one Draft PR to `main`. If the PR exists, continue the same PR only: require exact-head Fast; mark Ready only after Fast passes; require full Qualified including backend/race/randomized, renderer/deterministic, Chrome and WebKit; then merge only on exact-head green evidence. Fix legitimate defects on the same branch/PR and do not create retry/certification/promotion branches.
