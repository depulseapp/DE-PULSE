# DE.PULSE — Current Adaptive Roadmap

**Operational overlay date:** 2026-08-19  
**Stable baseline:** `v18.6.1-stable`  
**Current engineering slice:** `v18.6.7-development` — Fresh Reconciliation, Scope Bind & Legacy Test/Gate Hygiene  
**Authority:** This file is the current-state operational overlay. It does not replace permanent product contracts, `governance/ROADMAP.md`, immutable release evidence, or historical adaptive-governance records.

## Permanent invariants

- U.S. Equities Processing only.
- No Execution: no order routing, paper trading, P&L, portfolio or journal execution features.
- G0–G16 remains the only top-level release-gate model.
- macOS Apple Silicon and Windows x64 are required Stable platforms.
- Smart Provider Router v2 is the sole provider-routing owner; no duplicate router.
- GLD, SLV and USO remain explicit actionable tradable exceptions.
- Provider count never changes Market Mode by itself. Every provider capability must have an explicit Market Mode treatment.
- Adaptive Intelligence improves interpretation, prioritization and evidence use without silently rewriting deterministic market truth.
- GitHub is the source of truth across ChatGPT, Claude or another authorized engineering agent.
- File age/version naming is never a deletion criterion; ownership, consumers, evidence and unique regression coverage decide cleanup.

## Phase 0 — Post-v18.6.1 engineering hardening — COMPLETE

- **Packet A / PR #46:** Impact Planner v2, Release Rehearsal, failure taxonomy, current governance overlays and honest reconciliation baseline.
- **Packet B / PR #47:** reproducible Fast/Qualified dependencies, immutable Action pins, Playwright pin, safe cache and least-privilege gate.
- **Packet C / PR #48:** Chrome + WebKit co-primary browser coverage and risk-directed execution.
- **Packet D / PR #49:** durable Stable evidence, CI runtime/platform/browser-setup telemetry, amplification warnings and workflow structural lint.
- **Packet E / PR #50:** Documentation capability-oriented renderer ownership with explicit legacy fallback and direct Chrome + WebKit proof.

Packet E merged to `main` at `b3ca18c14b1e53069a6736e29ad9e3b09f87bda5`.

## v18.6.7 — Fresh Reconciliation, Scope Bind & Legacy Test/Gate Hygiene — ACTIVE

Purpose: establish a clean, evidence-backed engineering baseline after Phase 0 before broad v18.x product expansion.

### A. Fresh reconciliation and scope bind

The historical conserved authority ledger is now conclusively located at:

`release/v18.5.1/V17-V18-IMPLEMENTATION-RECONCILIATION.json`

Its v18.6.7 intake blob is `2a32b3f93203d61b1aca55172530652d736bbf55` and it declares **296 tracked authority rows**. The 296 count is conservation scope, not a count of open defects or unfinished features.

Current reconciliation is governed by `adaptive-governance/V18_6_7_CURRENT_RECONCILIATION.md`.

Rules:

1. Preserve every original ID/history; do not regenerate a replacement ledger.
2. Historical v18.5.1 statuses remain historical evidence and are never assumed current.
3. Fast runs the existing reconciliation gate in inventory mode to preserve schema, row count, IDs and canonical scope alignment.
4. Current row disposition uses `FRESH_PASS`, `REOPENED`, `NOT_IMPLEMENTED`, `INTENTIONALLY_SUPERSEDED`, `NOT_APPLICABLE`, or `ROADMAP_FUTURE_SCOPE` only when current evidence supports it.
5. Prioritize still-reproducible user-trust/runtime/provider/shared-intelligence/renderer risks for fresh row-level mapping.
6. Bind the next product slice at G0–G3; G1 remains the immutable scope-freeze point.

### B. Legacy Test & Gate Hygiene

Governed by:

- `adaptive-governance/LEGACY_TEST_GATE_CLEANUP_PLAN.md`
- `adaptive-governance/LEGACY_TEST_GATE_INVENTORY.md`
- `tools/ci/legacy_test_gate_inventory.py`

Every versioned executable test/gate is classified as one of:

- `ACTIVE_REQUIRED`
- `ACTIVE_DUPLICATE`
- `UNREFERENCED_USEFUL`
- `HISTORICAL_EVIDENCE`
- `SAFE_TO_REMOVE`

`SAFE_TO_REMOVE` is never inferred automatically. Useful assertions must be migrated before deletion.

### First safe organization wave — STAGED

The first wave preserves test bodies while replacing release-number ownership with capability ownership where technically safe:

- `v18_6_ai_hardening_test.go` → `ai_hardening_test.go`
- `v18_6_broad_snapshot_broker_test.go` → `broad_snapshot_broker_test.go`
- `v18_6_documentation_access_test.go` → `documentation_access_test.go`
- `v18_6_session_intelligence_coordinator_test.go` → `session_intelligence_coordinator_test.go`
- `v18_6_surface_consolidation_test.js` → `tests/renderer/surface_consolidation_test.js`
- `v18_6_documentation_access_test.js` → `tests/renderer/documentation_access_test.js`

The Go tests stay beside the current Go package so package-private access and `go test ./...` discovery remain intact. The renderer tests moved only after verifying they resolve production files from repository working-directory paths rather than script-relative imports.

Fast now consumes the new renderer paths. Impact Planner treats `tests/renderer/` and `tests/browser/` as `RENDERER_UI`, preserving full qualification and WebKit risk signaling for future organization work.

The current certification plan still explicitly consumes inherited/versioned evidence such as v16.10/v16.11 performance gates, v17.4 renderer acceptance and focused v17/v18 Go test prefixes. Those remain until their assertions and consumers are deliberately migrated; they are not cosmetic deletion candidates.

### v18.6.7 exit

- located 296-row authority ledger is structurally conserved and Fast-bound;
- current reconciliation truth is documented without inherited-status assumptions;
- targeted version-stacked root executable tests/gates have a live classification/consumer model;
- first proven-safe capability-oriented organization wave is complete;
- remaining retained legacy files have an explicit migration/deletion condition rather than an age-based assumption;
- no regression coverage is lost;
- affected Fast/Qualified evidence is green;
- next product scope is bound at G0–G3/G1 as appropriate.

No Stable release is manufactured for cleanup/governance/test-organization work alone. Stable release remains tied to genuine release-capable product/release-identity scope.

## Version-visible build train after v18.6.7

These are roadmap allocations. Each version becomes immutable only at its G1 scope freeze; evidence may combine or defer a patch when appropriate.

### v18.7.0 — Runtime Reliability & Data Truth — PROVISIONAL NEXT

Unless fresh reconciliation finds a higher-severity user-trust blocker that should lead or be bundled into this slice:

- exact `DATA DEGRADED` semantics and blast-radius isolation;
- provider/capability health aggregation;
- freshness SLO truth;
- duplicate-work suppression, single-flight and request coalescing;
- bounded retries and provider circuits;
- queue/backpressure/load shedding and runtime overload behavior;
- disagreement handling and recovery hysteresis;
- truthful `UNKNOWN` / `ABSTAIN` when evidence is insufficient;
- realistic active-market load/latency proof.

Do not rebuild controls current source already has; close only verified remaining gaps.

### v18.7.1 — User-Trust Reliability Closure — IF NEEDED

Focused patch for still-reproducible trust defects discovered during v18.7.0 qualification: stale/freshness/readiness truth, refresh/focus/state integrity or material UI reliability. Skip this version when there is no justified patch scope.

### v18.8.0 — Shared Intelligence Consolidation

- Scanner + Opportunity Radar → canonical acquisition/cache owner.
- Pre-Market Prep + Market Open Prep → Session Intelligence Coordinator.
- Earnings + Catalyst Reaction → Event Intelligence lifecycle.
- Research → canonical deep-evidence destination.

Remove duplicate internal pipelines, not useful user-facing capabilities.

### v18.8.1 — Renderer Modularization II

Continue the strangler model after Documentation. Candidate owners include Watchlist, Session, Research, Admin and shared UI state. Physical deletion follows direct equivalence and Chrome + WebKit proof, not file age.

### v18.9.0 — TradeInsight SHADOW Integration

Congressional Trading Intelligence, SEC Form 4 enrichment and historical OHLCV fallback/backfill only through Smart Provider Router v2 with explicit rights/entitlement/freshness/consumer/Market-Mode disposition. No provider-to-UI silo.

### v18.9.1 — Provider Intelligence & Market-Mode Hardening

Capability usefulness/freshness/disagreement/latency/headroom scoring, explicit Market Mode disposition, rate-limit/backpressure behavior and provider-role truth.

### v18.10.0 — v18 Major Closure Candidate

Fresh row-level reconciliation, zero unexplained applicable rows, closure of release blockers, deterministic/browser/provider/security/degradation/restart/load qualification, required native proof and G16 final v18 handoff.

### v18.10.1 — Closure Patch — ONLY IF NEEDED

Certification/closure defects only. No feature expansion. Skip when v18.10.0 closes cleanly.

## v19 — Professional Data Infrastructure

### v19.0.0
Mature provider capability/quality/latency/cost/entitlement/rights infrastructure, historical completeness, data lineage and reliability telemetry.

### v19.1.0
Institutional/13F and point-in-time evidence infrastructure, recommendation/outcome history and institutional behavior intelligence.

### v19.2.0
Professional research infrastructure: Two-Sided Thesis, entry/target/invalidation history, MFE/MAE, first-event ordering and deeper evidence usefulness validation.

## v20 — Adaptive Intelligence

### v20.0.0
Adaptive Intelligence foundation: historical analogues, regime-conditioned learning, calibration, false-positive/miss learning, provider/evidence usefulness and outcome tracking.

### v20.1.0
Adaptive Opportunity Intelligence: ranking maturation, contradiction learning, evidence weighting and workflow intelligence.

### v20.2.0
Champion/Challenger intelligence: controlled model/rule comparison with measurable promotion criteria; AI never silently rewrites deterministic market truth.

## Roadmap decision rule

Every future slice starts from immutable Stable, current reconciliation state and current CI evidence. Scope is selected for user value and risk reduction—not to make a version number look large, refresh file dates, or remove old filenames cosmetically.
