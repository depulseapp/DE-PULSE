# DE.PULSE — Current Adaptive Roadmap

**Operational overlay date:** 2026-08-19  
**Stable baseline:** `v18.6.1-stable`  
**Authority:** This file is the current-state operational overlay. It does not replace permanent product contracts, `governance/ROADMAP.md`, immutable release evidence, or historical adaptive-governance records. Where an older adaptive document describes v18.2/v18.5.1 as the active release, that statement is historical and this overlay is authoritative for current status.

## Permanent invariants

- U.S. Equities Processing only.
- No Execution: no order routing, paper trading, P&L, portfolio or journal execution features.
- G0–G16 remains the only top-level release-gate model.
- macOS Apple Silicon and Windows x64 are required Stable platforms.
- Smart Provider Router v2 is the sole provider-routing owner; no duplicate router.
- GLD, SLV and USO remain explicit actionable tradable exceptions.
- Provider count never changes Market Mode by itself. Every provider capability must be `INTEGRATED`, `CONTEXTUAL_ONLY`, `NOT_RELEVANT` or `INTENTIONALLY_HIDDEN` for Market Mode.
- Adaptive Intelligence improves interpretation, prioritization and evidence use without silently rewriting deterministic market truth.
- GitHub is the source of truth across ChatGPT, Claude or another authorized engineering agent.

## Phase 0 — Post-v18.6.1 engineering hardening

Purpose: eliminate remaining process debt exposed by the v18.6.1 trial before starting a broad feature slice.

### Packet A — CI decision intelligence and release rehearsal — CURRENT

- Impact Planner v2 with explicit change classes.
- Fail-closed lane selection: uncertain/product changes use full qualification.
- Side-effect-free pre-merge Release Rehearsal.
- Formal failure taxonomy: `PRODUCT_FAIL`, `GATE_TEST_FAIL`, `CI_HARNESS_FAIL`, `INFRA_FAIL`, `EXPECTED_NOOP`, `SUPERSEDED`.
- Current roadmap/build-plan/reconciliation overlays.
- Validate this packet itself with one development branch, one PR, Fast, then process-only Qualified CI-harness/portability.

### Packet B — Reproducible CI dependencies — PLANNED

- Pin GitHub Actions to immutable commit SHAs after exact upstream verification.
- Pin Playwright and browser revision used by browser gates.
- Add safe dependency/browser caching where it materially reduces runner time.
- Add generic workflow linting (`actionlint` or equivalent) in addition to DE.PULSE-specific policy checks.
- Review and tighten per-job permissions without weakening required exact-status/publication behavior.

### Packet C — Risk-directed browser coverage — PLANNED

- Use Impact Planner v2 `RENDERER_UI` classification to activate a focused WebKit lane for Safari-sensitive renderer/UI changes.
- Keep Chromium as the broad browser gate.
- Do not double the entire browser matrix for backend-only work.

### Packet D — Durable evidence and CI observability — PLANNED

- Emit a compact immutable Stable evidence manifest: source SHA/fingerprint, Fast/Qualified/Release run IDs, native artifact hashes, G12/G15/G16 state and tool versions.
- Track lane runtime, queue time, cache hit/miss and failure classification per build.
- Retain large transient logs only as long as useful; retain compact release truth durably.

### Packet E — Renderer maintainability — PLANNED, INCREMENTAL

Do not rewrite the monolithic renderer. Use a strangler approach by canonical capability owner:

1. Watchlist.
2. Session/header.
3. Research.
4. Documentation.
5. Admin.
6. Opportunity/Discovery.
7. Shared UI state.

For each extraction: extract → deterministic/equivalence proof → browser proof → native proof where relevant → remove the old duplicate owner. Historical/versioned assets remain until proven unreferenced and safely removable.

## Adaptive v18.x product work after Phase 0

Selection is evidence-driven at G0–G3; version numbers are assigned only after coherent scope is chosen.

### Workstream A — User-trust defects

Prioritize reproducible issues involving stale/freshness truth, refresh/focus behavior, visual instability, membership/state, misleading degraded states and readiness truth. Real-money decision-support trust issues outrank cosmetic additions.

### Workstream B — ADR-GDI / runtime reliability

Continue capability-level health, freshness SLOs, provider circuits, bounded retries, single-flight/coalescing, queue limits, backpressure, load shedding, warm-restart state, provider disagreement, consumer blast radius, `UNKNOWN`/`ABSTAIN`, and recovery hysteresis. A low-value source failure must not create misleading global degradation.

### Workstream C — Shared intelligence utility consolidation

- Scanner + Opportunity Radar: one canonical acquisition/cache owner.
- Pre-Market Prep + Market Open Prep: one Session Intelligence Coordinator.
- Earnings + Catalyst Reaction: one Event Intelligence lifecycle.
- Research: canonical deep-evidence destination.

Remove duplicate internal pipelines, not user-facing capabilities.

### Workstream D — Renderer architecture consolidation

Gradually replace version-stacked implementation ownership with capability-oriented modules such as `watchlist`, `session-intelligence-ui`, `research-ui`, `admin-ui`, and shared state. Release version belongs in release identity, not permanently in every module filename.

### Workstream E — TradeInsight controlled integration

TradeInsight remains `SHADOW` / `SECONDARY` through Smart Provider Router v2 only. Each capability requires explicit entitlement, rights, freshness, coverage, quality, role, consumer, storage/redistribution/AI-use disposition and Market Mode treatment. No provider-to-UI silo.

## v18 Major Closure

Do not close v18 until the conserved requirement ledger has zero unexplained applicable rows, no reopened release blocker remains, key workflows have fresh evidence, Day/Swing/Long deterministic equivalence passes, browser/UI/security/provider failure/rate-limit/degradation/restart/load behavior is qualified, macOS + Windows native artifacts pass, and G16 records the final retrospective/handoff.

## v19 — Professional Data Infrastructure

Mature provider capability scorecards, quality/latency/cost/entitlement/rights, disagreement handling, historical completeness, institutional/13F infrastructure, point-in-time evidence, recommendation/outcome history, provider usefulness and data-reliability telemetry.

## v20 — Adaptive Intelligence

Use accumulated validated evidence for historical analogues, ASBI, calibration, false-positive/false-negative learning, regime-conditioned behavior, provider/evidence usefulness, Adaptive Opportunity Discovery, Two-Sided Thesis intelligence, institutional behavior, contradiction learning, outcome tracking and Champion/Challenger evaluation. Deterministic numeric truth remains deterministic.

## Roadmap decision rule

Every future slice starts from immutable Stable, current reconciliation state and current CI evidence. Scope is selected for user value and risk reduction, not to make a version number look large or to touch old files. File age alone is never a cleanup criterion.
