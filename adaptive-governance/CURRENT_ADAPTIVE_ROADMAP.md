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

Purpose: eliminate remaining process and maintainability debt exposed by the v18.6.1 trial before starting a broad feature slice.

### Packet A — CI decision intelligence and release rehearsal — COMPLETE

Merged PR #46. Delivered Impact Planner v2, fail-closed lane selection, side-effect-free Release Rehearsal, formal failure taxonomy, current governance overlays and the honest v18.6.1 reconciliation baseline.

### Packet B — Reproducible CI dependencies — COMPLETE

Merged PR #47. Fast/Qualified third-party Actions are immutable-SHA pinned, Playwright is pinned, browser pip caching is deterministic, permissions are gated, and dependency drift fails closed. Release-workflow Action pinning remains intentionally deferred to the next genuine release-capable product slice.

### Packet C — Chrome + WebKit primary browser coverage — COMPLETE

Merged PR #48 at `23ecb71f60e1658d68bcef6248044ce53b6dd851`.

- Chrome and WebKit are the two primary engines.
- Chrome remains the broad behavior suite.
- WebKit is the co-primary compatibility proof and is mandatory for `full`, `browser`, renderer/UI and WebKit-harness risk.
- WebKit runs on `macos-15` without Linux apt dependency amplification.
- Backend/provider/process work avoids browser runtime when unaffected.
- Other engines remain secondary/risk-directed.
- Final exact-head Fast #393 and Qualified #138 passed, including real WebKit compatibility and Ubuntu/macOS/Windows portability.

### Packet D — Durable evidence and CI observability — COMPLETE

Merged PR #49 at `2885de409c86f771d582f09f54e0f6c564f6c59d`.

- repo-durable v18.6.1 Stable evidence manifest bound to the authoritative checkpoint;
- Stable evidence drift gate;
- Qualified per-job queue/runtime/platform telemetry;
- Linux/macOS/Windows runner-consumption visibility without fabricated currency rates;
- Chrome/WebKit dependency setup duration and pip-cache hit signals;
- PR-level Fast/Qualified/Release amplification warnings;
- compact 30-day operational telemetry artifact plus job summary;
- zero-network generic workflow structural lint in addition to DE.PULSE semantic workflow policy.

Final Packet D Fast #396 and Qualified #139 passed. Qualified correctly selected CI-harness + Ubuntu/macOS/Windows portability and skipped product/browser lanes. Telemetry reported Fast 2 / Qualified 1 / Release 0 as `OK`, correctly recognizing the one legitimate lint correction rather than a loop.

### Packet E — Renderer modularization foundation — ACTIVE

Fresh source inventory shows the current classic renderer is still layered: `renderer.js` is about 425 KB, `styles.css` about 316 KB, and `index.html` loads the monolith before compatibility/feature layers. File age is not a cleanup signal; active ownership is.

Packet E uses a strangler model rather than a monolith rewrite. The first bound capability is **Documentation** because it is lower trading-risk, already has an explicit role-access decorator, and its rendering/hydration/Markdown lifecycle is embedded in the monolith.

Current Packet E scope:

1. Add capability-oriented `renderer/documentation-ui.js` as the active runtime owner for Documentation Markdown, hydration and view rendering.
2. Load it after `renderer.js` and before `documentation-access-v18.6.js`, so the existing role-access policy decorates the new owner rather than competing with it.
3. Register renderer ownership metadata and the access decorator explicitly.
4. Retain the monolith Documentation implementations as an **inactive legacy fallback** for this first strangler step; do not pretend physical deletion is complete.
5. Keep the transitional dependency on monolith `architectureDiagram` explicit; it is a later extraction target rather than hidden coupling.
6. Add a renderer owner contract that prevents version-stacked new owner naming, duplicate loads, load-order drift, untracked fallback deletion, or loss of primary-engine evidence.
7. Extend the existing Fast role-access regression so it proves the new owner.
8. Require direct owner evidence in Qualified renderer plus both primary engines: Chrome and WebKit.
9. Preserve deterministic market truth; Documentation extraction must not modify scoring/market math.

Packet E is renderer/product source and therefore must receive full Qualified coverage plus WebKit. It does not change release identity or the canonical Release workflow, so it does not manufacture a new Stable release by itself.

### Physical-deletion rule after Packet E

A runtime owner can move before the legacy source is deleted. Physical removal from `renderer.js` occurs only when:

- no runtime/reference consumer depends on the fallback;
- direct equivalence evidence exists;
- Chrome + WebKit owner behavior passes;
- renderer logic/deterministic tests stay green;
- the capability owner contract is updated from `ACTIVE_OWNER_WITH_LEGACY_FALLBACK` to a no-fallback state.

This prevents a risky big-bang renderer rewrite while still reducing active ownership ambiguity immediately.

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

Continue capability-oriented ownership after the Documentation proof. Candidate sequence is evidence-selected, not fixed; higher-risk Watchlist/Session/Research owners require stronger transition evidence than Documentation. New long-lived owners use capability names such as `watchlist.js`, `session-intelligence-ui.js`, `research-ui.js`, `admin-ui.js`, and shared UI state rather than accumulating release-version filenames.

### Workstream E — TradeInsight controlled integration

TradeInsight remains `SHADOW` / `SECONDARY` through Smart Provider Router v2 only. Each capability requires explicit entitlement, rights, freshness, coverage, quality, role, consumer, storage/redistribution/AI-use disposition and Market Mode treatment. No provider-to-UI silo.

## v18 Major Closure

Do not close v18 until the conserved requirement ledger has zero unexplained applicable rows, no reopened release blocker remains, key workflows have fresh evidence, Day/Swing/Long deterministic equivalence passes, Chrome + WebKit UI behavior, security/provider failure/rate-limit/degradation/restart/load behavior are qualified, macOS + Windows native artifacts pass, and G16 records the final retrospective/handoff.

## v19 — Professional Data Infrastructure

Mature provider capability scorecards, quality/latency/cost/entitlement/rights, disagreement handling, historical completeness, institutional/13F infrastructure, point-in-time evidence, recommendation/outcome history, provider usefulness and data-reliability telemetry.

## v20 — Adaptive Intelligence

Use accumulated validated evidence for historical analogues, ASBI, calibration, false-positive/false-negative learning, regime-conditioned behavior, provider/evidence usefulness, Adaptive Opportunity Discovery, Two-Sided Thesis intelligence, institutional behavior, contradiction learning, outcome tracking and Champion/Challenger evaluation. Deterministic numeric truth remains deterministic.

## Roadmap decision rule

Every future slice starts from immutable Stable, current reconciliation state and current CI evidence. Scope is selected for user value and risk reduction, not to make a version number look large or to touch old files. File age alone is never a cleanup criterion.
