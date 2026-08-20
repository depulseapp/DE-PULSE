# DE.PULSE — Current Adaptive Build Plan

**Immutable Stable:** `v18.8.0-stable`  
**Stable candidate:** `3a32d57dd4c74c6f812cc942a9d8049a7b517718`  
**Engineering branch:** `v18.8.1-development`  
**Next:** `v18.8.1 — Cost per Trustworthy Evidence + 10/10 Audit Hardening`.

## Entry condition

`ADAPT-REL-001` is CLOSED. The actual v18.8.0 Stable tag/Release is reconciled into the build checkpoint, release-evidence checkpoint, Stable evidence manifest and `handoff/CURRENT.md`. This was continuity/source-of-truth closure only; unchanged certified v18.8.0 binaries were not rebuilt.

## v18.8.1 mandatory work packets

| Order | Packet | State | Durable implementation/evidence |
| ---: | --- | --- | --- |
| 1 | `ADAPT-CI-001` Release State Coherence | **IMPLEMENTED** | `35531cd98465565c60cb2a26a3e066692d0f2168` plus CI-hardening invariant lock `fd58ce8949013f5e6ffd70b69a3a991f8d4453f1`. |
| 2 | `ADAPT-CI-002` Early G11 Target Preflight | **IMPLEMENTED** | Release-state coherence/G11 target preflight plus invariant coverage in the CI-hardening commits. |
| 3 | `ADAPT-CI-003` Cheap-First Fast Ordering | **IMPLEMENTED** | `89e78f0da57c4819c2ad818c73541fcc7713f269`; cheap governance/coherence/provenance checks precede expensive Go/Node/browser setup. |
| 4 | `ADAPT-CI-004` Safe Manual Dispatch Defaults | **IMPLEMENTED** | `02da4e1a560e16671188abdc69037116de4994c2`; manual Qualified defaults to adaptive/impact-safe qualification and full remains explicit. |
| 5 | `ADAPT-DATA-001` Universe Eligibility Contract | **IMPLEMENTED** | `3337386c5492a7af49e6e4dc49ef25dd23f94a44` + follow-up `d63d0753d4d321362934e0ad7dabc52b7dca9b32`; broad U.S.-equity eligibility is explicit and no longer silently constrained by `has_options`. |
| 6 | `ADAPT-DATA-002` Evidence-Time Truth | **IMPLEMENTED** | `c01bd4d683be6236ec884aaea0cc14f7adf4d0b8`; provider observation/evidence time drives Scanner/Radar freshness, retrieval time is bookkeeping, and unknown evidence remains unknown. |
| 7 | `ADAPT-ARCH-001` Shared-Universe Robustness | **CLOSED BY FRESH EVIDENCE** | Shared neutral universe owner is explicit in `symbol_universe.go`; deferred cleanup is cancellation/panic-safe; focused eligibility/retrieval-time/concurrency recovery tests live in `v18_8_1_universe_hardening_test.go`. Smart Provider Router v2 and BroadSnapshotBroker ownership are unchanged. |
| 8 | `ADAPT-UI-001` Renderer Modularization II | **IMPLEMENTED** | `f2f30d0c160f7bbf8e01f31271faf86d819808e8`; active Market Pulse/header ownership moved to release-neutral `renderer/market-header-ui.js`, behavior and `__v1851HeaderContracts` compatibility are preserved, and the renderer owner contract forbids the version-stacked header from remaining an active runtime owner. |
| 9 | `ADAPT-QA-001` Test/Gate Consolidation | **IMPLEMENTED** | `136af46f4f208edfd97fd728b3e3f1af61f2af31`; canonical `tests/renderer/market_header_owner_test.js` now proves Market Header owner/registry, ribbon order, idempotent ensure behavior, data-health updates, one base chrome update per wrapper call and legacy compatibility. It is transitively bound to the existing Qualified renderer owner regression, and `renderer_owner_contract.py` makes that binding mandatory without adding a workflow/job. Unique v18.5.1 browser hierarchy evidence remains conserved. |
| 10 | `ADAPT-GOV-001` Historical Reconciliation Identity | **IMPLEMENTED** | `bec6b8f0b4721e5eda891c22f72d51964ccd6590`; immutable `release/v18.5.1/HISTORICAL-RECONCILIATION-IDENTITY.json` now binds the historical ledger schema/release, incoming `v18.5.0-stable` tag/commit and `v18.5.1-development` branch independently of current `release_identity.json`; mutation self-tests are consumed by the existing reconciliation gate. |
| 11 | `ADAPT-COST-001` Cost per Trustworthy Evidence | **NEXT** | Measure avoided runs/setup/minutes and cost per trustworthy evidence without weakening required gates. |
| 12 | `ADAPT-RECON-001` Zero-Miss Reconciliation | **PENDING** | Fresh zero-miss reconciliation for every applicable v17→v18.8 conserved requirement plus post-ledger approved Adaptive commitments; bind current owner, source/behavior, regression/evidence and explicit disposition before v18 closure. |
| 13 | `ADAPT-UX-RESEARCH-001` Research Information Architecture | **PENDING** | Revalidate/fix Research ticker/input consistency, freshness badge, responsive top-area hierarchy, disabled/recovery states and containment. |
| 14 | `ADAPT-SYMBOL-001` Symbol/Desk Correctness | **PENDING** | Revalidate/fix DESKS truth, desk transitions, final-desk semantics, add/remove idempotency, Master Market Symbols, persistence/reload and Undo across Day/Swing/Long. |
| 15 | `ADAPT-READINESS-001` Prep/Readiness Semantics | **PENDING** | Revalidate/fix Pre-Market/Market Open preparation, missed-window catch-up, persisted job state, readiness evidence, SEC/event risk, EXTENDED and Catalyst Reaction lifecycle/measurements without redundant broad refetches. |
| 16 | `ADAPT-FRESHNESS-001` Freshness/Data Engine Correctness | **PENDING** | Revalidate/fix targeted Refresh/Age, dataset/session cadence, automatic stale recovery, priority refresh, source/reason/fallback truth and stable freshness UI using existing canonical freshness/degradation owners. |
| 17 | `ADAPT-RESEARCH-002` Research Correctness Closure | **PENDING** | Preserve all approved Research capabilities during Renderer Modularization II with behavior/browser/native equivalence; no capability loss through module extraction. |

A historical carry-forward that is already correct may close by fresh evidence. Do not rebuild working functionality solely because the requirement is old; do not leave a reproducible gap for a generic v18.10 catch-all.

## v18.9.0 mandatory TradeInsight packet

18. `ADAPT-TRADEINSIGHT-001` TradeInsight Full Capability Discovery, Utility Mapping & SHADOW Integration. Enumerate the complete capability surface actually available to the configured beta account/API. Congressional Trading, SEC Form 4 enrichment secondary to direct SEC, and historical OHLCV fallback/backfill are mandatory minimum roles, not a cap. Classify every available capability as `USE`, `CORROBORATE`, `FALLBACK`, `STORE_FOR_HISTORY`, `FUTURE`, or `NOT_USEFUL`; for every useful capability bind Smart Provider Router v2 ownership, consumer, authority, rights/entitlement, evidence-time/freshness, budget/rate limits, cache/retention, disagreement handling, Market Mode disposition, SHADOW telemetry, promotion criteria and graceful degradation. Full capability never means blindly calling every endpoint; canonical reuse and purpose-driven routing remain mandatory.

## Required execution order

`G0 continuity closure → G2/G3 owner/semantic design → CI hardening → data-truth/shared-universe hardening → bounded renderer modularization → test/gate consolidation → historical identity hardening → cost hardening → zero-miss user-trust reconciliation → G10 reconciliation → exact-head Fast/Qualified → one G11–G16 Release when release-capable → v18.9.0 TradeInsight full-capability SHADOW → v18.9.1 provider intelligence → v18.10 zero-gap closure`.

Use one version-development branch per release, coherent batches, one Draft PR, automatic Fast, same PR Ready, Qualified, exact-head merge and one Release. Never create retry/certification/promotion branches or duplicate workflows.

## Forward provider-intelligence completion

- **v18.9.0:** discover and SHADOW every useful TradeInsight beta capability through Smart Provider Router v2; no provider-to-UI silo.
- **v18.9.1:** use measured capability-level health, freshness, latency, disagreement, headroom, rights, cost and usefulness across TradeInsight and the wider provider pool for smarter routing and explicit Market Mode treatment.
- **v18.10.0:** prove zero unexplained provider/capability or v17→v18.8 requirement remainder before v18 closure.
- **v19:** mature professional provider/data quality, lineage, rights, cost and historical infrastructure.
- **v20:** mature adaptive provider/evidence intelligence from accumulated usefulness/outcome history; deterministic price truth and Day/Swing/Long formulas remain protected.

Permanent boundaries remain: U.S. Equities only, No Execution, G0–G16 only, Smart Provider Router v2 sole provider-routing owner, BroadSnapshotBroker sole broad snapshot reuse/coalescing owner, deterministic Day/Swing/Long truth protected.

## Exactly one next action

Execute `ADAPT-COST-001` on `v18.8.1-development`: measure avoided runs/setup/minutes and cost per trustworthy evidence using the existing CI telemetry and qualification architecture; preserve cheap-first, impact-routed, exact-head quality gates; add bounded telemetry/policy evidence only where a real gap exists; do not create a new workflow/job/branch or trigger qualification CI merely for packet development.