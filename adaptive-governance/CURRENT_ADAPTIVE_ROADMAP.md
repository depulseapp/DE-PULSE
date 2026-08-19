# DE.PULSE — Current Adaptive Roadmap

**Operational overlay date:** 2026-08-19  
**Certified Stable baseline:** `v18.6.1-stable`  
**Current engineering slice:** `v18.7.0-development` — Runtime Reliability & Data Truth  
**Authority:** current-state operational overlay. Permanent product contracts, `governance/ROADMAP.md`, immutable release evidence and historical adaptive-governance records remain preserved.

## Permanent invariants

- U.S. Equities Processing only.
- No Execution: no order routing, paper trading, P&L, portfolio or journal execution features.
- G0–G16 remains the only top-level release-gate model.
- macOS Apple Silicon and Windows x64 are required Stable platforms.
- Smart Provider Router v2 is the sole provider-routing owner; no duplicate router.
- GLD, SLV and USO remain explicit actionable tradable exceptions.
- Provider count never changes Market Mode by itself. Every provider capability requires explicit treatment.
- Adaptive Intelligence improves synthesis/prioritization/evidence use without silently rewriting deterministic market truth.
- GitHub is the source of truth across ChatGPT, Claude or another authorized engineering agent.
- File age/version naming is never a deletion criterion; ownership, consumers, evidence and unique regression coverage decide cleanup.

## Post-v18.6.1 engineering baseline — COMPLETE

Phase 0 hardening PRs #46–#50 completed Impact Planner/rehearsal, reproducible Fast/Qualified dependencies, Chrome + WebKit policy, CI telemetry and Documentation renderer ownership.

**v18.6.7 / PR #51 is COMPLETE.** It merged to `main` at `f1a9e0d0d76d4be565ac9355a09f77a33a3338a7` after exact-head Fast and Qualified success. It conserved the 296-row authority ledger, established the legacy test/gate inventory and completed the first safe capability-oriented test organization wave. It did not manufacture a Stable release.

Historical details remain in:

- `adaptive-governance/V18_6_7_CURRENT_RECONCILIATION.md`
- `adaptive-governance/LEGACY_TEST_GATE_CLEANUP_PLAN.md`
- `adaptive-governance/LEGACY_TEST_GATE_INVENTORY.md`
- `release/v18.5.1/V17-V18-IMPLEMENTATION-RECONCILIATION.json`

## v18.7.0 — Runtime Reliability & Data Truth — ACTIVE / G1 FROZEN

Authoritative scope/evidence:

- `v18_7_0_scope.json`
- `v18_7_0_g0_g3_contract.json`
- `adaptive-governance/V18_7_0_RUNTIME_RELIABILITY_AUDIT.md`
- `release/v18.7.0/release_contract.json`
- `release/v18.7.0/run_full_certification.sh`

### Source-audit conclusion

Do **not** build parallel reliability infrastructure. Existing canonical owners are retained:

- Smart Provider Router v2 for provider selection/capability state;
- `data_freshness.go` for dataset/session/provider freshness;
- `runtime_slo.go` for SLO/recovery tracking;
- `workload_controller.go` for bounded queues/backpressure/load shedding;
- `runtime_observability.go` for request/load telemetry and budgets;
- `broad_snapshot_broker.go` for broad-snapshot reuse/coalescing;
- `provider_reconciliation.go` for cross-provider disagreement truth;
- `Engine.Snapshot()` as the canonical aggregation boundary.

### Verified v18.7 gaps being closed

1. Canonical degradation reason taxonomy while retaining compatible display labels.
2. Explicit affected-consumer/blast-radius truth.
3. Fail-closed `UNKNOWN` + `ABSTAIN` when required decision evidence is insufficient and no narrower cause is proven.
4. Recovery hysteresis: three consecutive healthy observations plus >=5 seconds stability; relapse resets the streak.
5. Immediate `QUEUE_SATURATED` truth at a hard-full provider queue.
6. Release workflow immutable Action pinning and locked G12 browser dependency setup.
7. Active-market reliability proof using production coalescing, bounded workload and degradation/SLO owners.

### Existing controls retained/requalified rather than rebuilt

- provider/capability health aggregation and entitlement/circuit truth;
- session-aware freshness SLOs;
- duplicate-work suppression/single-flight/coalescing;
- bounded routing/retry/circuit/cooldown behavior;
- tiered request budgets/backpressure/load shedding;
- cross-provider disagreement/reconciliation;
- canonical reuse telemetry.

### Candidate/release identity

The development branch is being prepared as the single release-capable v18.7 candidate. `release_identity.json` now binds `18.7.0` / `v18.7.0-stable-20260819`, but **this does not mean v18.7.0 is Stable**. Stable remains `v18.6.1-stable` until exact-head Fast + full Qualified pass, the PR merges, and one G11–G16 Release run certifies/packages/audits/publishes the exact candidate.

Frontend identity uses `renderer/release-identity-v18.7.0.js`, a small last-loaded identity overlay. This avoids rewriting the legacy renderer monolith solely for version constants; the existing v18.6.1 watchlist behavior extension remains intact underneath.

### Required qualification before v18.7.0 completion

- exact-head Fast;
- Go formatting/vet/full suite;
- v17 reliability/backpressure regressions;
- v18.7 reliability + active-market synthetic proof;
- race detector;
- randomized Go order;
- deterministic/renderer contracts;
- Chrome broad behavior;
- WebKit primary compatibility;
- conserved 296-row ledger integrity / G10 reconciliation;
- Release workflow/reproducibility contracts;
- merge only on exact-head green evidence;
- one G11–G16 Stable run with macOS Apple Silicon + Windows x64 actual packaged runtime proof and no-rebuild publication.

No staged test is called PASS before CI executes it.

## Version-visible sequence after v18.7.0

Each exact scope remains subject to G1 freeze.

### v18.7.1 — User-Trust Reliability Closure — ONLY IF NEEDED

Focused patch only for still-reproducible trust defects discovered during v18.7 qualification: stale/readiness/refresh/focus/state/UI reliability. Skip when there is no justified patch scope.

### v18.8.0 — Shared Intelligence Consolidation

- Scanner + Opportunity Radar → canonical acquisition/cache owner.
- Pre-Market Prep + Market Open Prep → Session Intelligence Coordinator.
- Earnings + Catalyst Reaction → Event Intelligence lifecycle.
- Research → canonical deep-evidence destination.

Remove duplicate pipelines, not useful user-facing capabilities.

### v18.8.1 — Renderer Modularization II

Continue the strangler model after Documentation. Candidate owners include Watchlist, Session, Research, Admin and shared UI state. Deletion follows equivalence + Chrome/WebKit proof, not file age.

### v18.9.0 — TradeInsight SHADOW Integration

Congressional Trading Intelligence, SEC Form 4 enrichment and historical OHLCV fallback/backfill through Smart Provider Router v2 only, with explicit rights/entitlement/freshness/consumer/Market-Mode disposition.

### v18.9.1 — Provider Intelligence & Market-Mode Hardening

Provider usefulness/freshness/disagreement/latency/headroom scoring, explicit Market Mode disposition and provider-role truth.

### v18.10.0 — v18 Major Closure Candidate

Fresh row-level reconciliation, zero unexplained applicable rows, closure blockers resolved, deterministic/browser/provider/security/degradation/restart/load qualification, required native proof and G16 final v18 handoff.

### v18.10.1 — Closure Patch — ONLY IF NEEDED

Certification/closure defects only; no feature expansion.

## v19 — Professional Data Infrastructure

- **v19.0.0:** mature provider capability/quality/latency/cost/entitlement/rights infrastructure, historical completeness, lineage and reliability telemetry.
- **v19.1.0:** institutional/13F and point-in-time evidence infrastructure, recommendation/outcome history and institutional behavior intelligence.
- **v19.2.0:** professional research infrastructure including Two-Sided Thesis, entry/target/invalidation history, MFE/MAE and deeper evidence usefulness validation.

## v20 — Adaptive Intelligence

- **v20.0.0:** historical analogues, regime-conditioned learning, calibration, false-positive/miss learning, provider/evidence usefulness and outcome tracking.
- **v20.1.0:** adaptive opportunity ranking, contradiction learning, evidence weighting and workflow intelligence.
- **v20.2.0:** champion/challenger intelligence with measurable promotion criteria; AI never silently rewrites deterministic market truth.

## Roadmap decision rule

Every future slice starts from immutable Stable, current reconciliation state and current CI evidence. Scope is selected for user value and risk reduction—not to make a version number look large, refresh file dates, duplicate an existing owner or remove historical filenames cosmetically.
