# DE.PULSE — Current Adaptive Build Plan

**Immutable Stable:** `v18.8.0-stable`  
**Stable candidate:** `3a32d57dd4c74c6f812cc942a9d8049a7b517718`  
**Next:** `v18.8.1 — Renderer Modularization II + 10/10 Audit Hardening`.

## Entry condition

First close `ADAPT-REL-001`: reconcile the actual v18.8.0 Stable tag/Release into build checkpoint, release-evidence checkpoint, Stable evidence manifest and `handoff/CURRENT.md`. This is continuity/source-of-truth closure only; unchanged certified v18.8.0 binaries do not rebuild.

## v18.8.1 mandatory work packets

1. `ADAPT-CI-001` Release State Coherence validator covering release identity, VERSION/app/renderer/cache-bust identity, checkpoints, Stable manifest, handoff, predecessor/baseline and target tag state in one result.
2. `ADAPT-CI-002` G11 target-tag/version/build/predecessor/release-scaffold preflight before G12/native work, retaining the final publication collision guard.
3. `ADAPT-CI-003` cheap-first Fast ordering; expensive Go/Node/browser setup starts only after cheap governance/coherence/identity/provenance checks are green.
4. `ADAPT-CI-004` safe manual dispatch defaults; full qualification must be explicit.
5. `ADAPT-DATA-001` explicit discovery-universe eligibility and intentional treatment of provider filters such as `has_options`.
6. `ADAPT-DATA-002` provider evidence time separated from retrieval time; missing evidence time is UNKNOWN/degraded/ABSTAIN rather than wall-clock `now`.
7. `ADAPT-ARCH-001` neutral shared-universe diagnostics, panic/cancellation-safe single-flight cleanup and explicit eligibility/concurrency tests.
8. `ADAPT-UI-001` Renderer Modularization II via bounded capability-owner extraction, not a broad redesign.
9. `ADAPT-QA-001` inventory and gradual capability-based consolidation of version-stacked tests/gates; preserve unique historical evidence.
10. `ADAPT-GOV-001` separate historical reconciliation baseline identity from current release identity.
11. `ADAPT-COST-001` measure avoided runs/setup/minutes and cost per trustworthy evidence without weakening required gates.
12. `ADAPT-RECON-001` fresh zero-miss reconciliation for every applicable v17→v18.8 conserved requirement plus post-ledger approved Adaptive commitments; bind current owner, source/behavior, regression/evidence and explicit disposition before v18 closure.
13. `ADAPT-UX-RESEARCH-001` revalidate/fix Research ticker/input consistency, freshness badge, responsive top-area hierarchy, disabled/recovery states and containment.
14. `ADAPT-SYMBOL-001` revalidate/fix DESKS truth, desk transitions, final-desk semantics, add/remove idempotency, Master Market Symbols, persistence/reload and Undo across Day/Swing/Long.
15. `ADAPT-READINESS-001` revalidate/fix Pre-Market/Market Open preparation, missed-window catch-up, persisted job state, readiness evidence, SEC/event risk, EXTENDED and Catalyst Reaction lifecycle/measurements without redundant broad refetches.
16. `ADAPT-FRESHNESS-001` revalidate/fix targeted Refresh/Age, dataset/session cadence, automatic stale recovery, priority refresh, source/reason/fallback truth and stable freshness UI using existing canonical freshness/degradation owners.
17. `ADAPT-RESEARCH-002` preserve all approved Research capabilities during Renderer Modularization II with behavior/browser/native equivalence; no capability loss through module extraction.

A historical carry-forward that is already correct may close by fresh evidence. Do not rebuild working functionality solely because the requirement is old; do not leave a reproducible gap for a generic v18.10 catch-all.

## v18.9.0 mandatory TradeInsight packet

18. `ADAPT-TRADEINSIGHT-001` TradeInsight Full Capability Discovery, Utility Mapping & SHADOW Integration. Enumerate the complete capability surface actually available to the configured beta account/API. Congressional Trading, SEC Form 4 enrichment secondary to direct SEC, and historical OHLCV fallback/backfill are mandatory minimum roles, not a cap. Classify every available capability as `USE`, `CORROBORATE`, `FALLBACK`, `STORE_FOR_HISTORY`, `FUTURE`, or `NOT_USEFUL`; for every useful capability bind Smart Provider Router v2 ownership, consumer, authority, rights/entitlement, evidence-time/freshness, budget/rate limits, cache/retention, disagreement handling, Market Mode disposition, SHADOW telemetry, promotion criteria and graceful degradation. Full capability never means blindly calling every endpoint; canonical reuse and purpose-driven routing remain mandatory.

## Required execution order

`G0 continuity closure → G2/G3 owner/semantic design → CI hardening → data-truth/shared-universe hardening → zero-miss user-trust reconciliation → bounded renderer/test modularization → G10 reconciliation → exact-head Fast/Qualified → one G11–G16 Release when release-capable → v18.9.0 TradeInsight full-capability SHADOW → v18.9.1 provider intelligence → v18.10 zero-gap closure`.

Use one version-development branch per release, coherent batches, one Draft PR, automatic Fast, same PR Ready, Qualified, exact-head merge and one Release. Never create retry/certification/promotion branches or duplicate workflows.

## Forward provider-intelligence completion

- **v18.9.0:** discover and SHADOW every useful TradeInsight beta capability through Smart Provider Router v2; no provider-to-UI silo.
- **v18.9.1:** use measured capability-level health, freshness, latency, disagreement, headroom, rights, cost and usefulness across TradeInsight and the wider provider pool for smarter routing and explicit Market Mode treatment.
- **v18.10.0:** prove zero unexplained provider/capability or v17→v18.8 requirement remainder before v18 closure.
- **v19:** mature professional provider/data quality, lineage, rights, cost and historical infrastructure.
- **v20:** mature adaptive provider/evidence intelligence from accumulated usefulness/outcome history; deterministic price truth and Day/Swing/Long formulas remain protected.

Permanent boundaries remain: U.S. Equities only, No Execution, G0–G16 only, Smart Provider Router v2 sole provider-routing owner, BroadSnapshotBroker sole broad snapshot reuse/coalescing owner, deterministic Day/Swing/Long truth protected.