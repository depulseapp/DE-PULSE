# DE.PULSE — Current Adaptive Delivery Process

**Certified Stable:** `v18.9.0-stable`  
**Stable candidate:** `9ea81cddae4875ae15d3719ca028519a36c597b6`  
**Stable fingerprint:** `a8719090c341c874dbd1279cc31ad98e84075d5701c46a800bf951340780ecb9`  
**Master corrective program:** #65 / `ADAPT-PROVIDER-INTELLIGENCE-010`  
**Immediate patch:** `v18.9.1` / #64 / `ADAPT-RUNTIME-CRASH-001`.

## v18.9.0 delivery — COMPLETE / IMMUTABLE

Fast #481, Qualified #153 and Release #32 are the release authority for v18.9.0. G11-G16 passed, including actual macOS Apple Silicon + Windows x64 packaged-runtime audits, G15 assurance, no-rebuild Stable publication and G16 handoff evidence. Later corrective work does not rewrite that evidence.

## Permanent small-patch delivery policy

v18.9.x, v19 and v20 ship **small, independently auditable patches** rather than heavy releases.

For each patch:
1. establish exact predecessor/baseline at G0;
2. freeze one primary responsibility at G1 with explicit non-goals;
3. complete G2/G3 ownership/contract readiness before source edits;
4. use one development branch and one PR;
5. batch a coherent code+test candidate before opening/advancing CI;
6. exact-head Fast -> same PR Ready -> Qualified;
7. merge only the qualified exact head;
8. run the single canonical G11-G16 release path when release-capable;
9. publish only certified same-run artifacts without rebuild;
10. perform a post-release implementation-miss/open-issue audit and reconcile handoff/checkpoints before the next patch begins.

No duplicate release workflows, retry branches, certification branches or avoidable duplicate CI runs. Failure classification is required before rerun. If a planned patch is too large, split it before G1 rather than weakening auditability.

## v18.9.x ordered delivery train

- **v18.9.1** — runtime crash corrective only (#64).
- **v18.9.2** — TradeInsight Settings/API-key UX only.
- **v18.9.3** — coverage-aware Smart Provider Router v2 + persistence-first residual-gap fulfillment only.
- **v18.9.4** — canonical company identity/all-desk presentation only.
- **v18.9.5** — Market Data Modes + capability diagnostics only.
- **v18.9.6** — TradeInsight SEC Form 4 enrichment only.
- **v18.9.7** — TradeInsight ticker/company search only.
- **v18.9.8** — TradeInsight movers/ranking evidence only.
- **v18.9.9** — remaining useful TradeInsight capability sweep only.
- **v18.9.10** — provider efficiency + Adaptive Intelligence telemetry + protected-session headroom measurement only.
- **v18.9.11** — Session-Aware Data Readiness Maintenance only: light overnight + heavy weekend, with strict provider/runtime protection for pre-market, regular market and after-hours.
- **v18.9.12** — professional closure audit only; no new feature scope.

The train is dependency-ordered, but a later patch may be split further if its G0/G1 proves it is still too large. Patch numbers are not a reason to bundle scope.

## Persistence-first delivery contract

For affected data paths, delivery proof must show the runtime checks canonical memory/cache and persisted DB/state before external acquisition, validates freshness/coverage/schema/provenance/rights, computes only the residual gap, fetches only that gap when provider contracts permit, merges/reconciles deterministically and persists truthful point-in-time/revision state.

A green provider call is not sufficient evidence. Delivery must prove the consumer requirement is complete enough or truthfully partial/UNKNOWN.

## Session-aware maintenance delivery contract

The canonical U.S. market calendar/session owner defines all session boundaries, holidays, half-days and exceptional closures. No second calendar/scheduler truth may be created.

### Protected Tier-0 sessions

**Pre-market, regular market and after-hours** must remain top-priority operating sessions. Release evidence for v18.9.11 and later reliability work must demonstrate:

- explicit provider quota/headroom reservation for current-session capabilities;
- live/current quote, spread/liquidity, VIX/market context, news/catalyst/SEC, Opportunity Radar, Research/readiness and other decision-critical work outrank maintenance;
- maintenance external-provider acquisition suspends during protected sessions unless the acquisition is directly required by a live consumer;
- bounded CPU, memory, network, DB and background-worker usage so maintenance cannot materially increase current-session latency or cause self-inflicted `DATA DEGRADED`;
- prompt preemption/drain/checkpoint/resume when protected sessions, market shocks or higher-priority current work starts;
- no heavy compaction, deep reconciliation or historical fan-out during protected sessions.

### Light overnight delivery proof

Prove bounded daily readiness work after protected after-hours and before protected pre-market: residual history gaps, completed-session finalization, incremental disclosures/revisions, small corporate-action/fundamental/macro/identity corrections, bounded outcome resolution, lightweight integrity checks and warm-state preparation. Prove maintenance stops/drains before the pre-market protection buffer or when provider/runtime headroom is insufficient.

### Heavy weekend delivery proof

Prove deeper but bounded market-closed work: historical backfill/reconciliation, corporate actions, SEC/Form4/13F/congress/earnings/fundamental/macro history, outcome lineage, provider-history consolidation, material feature preparation and bounded DB/index/retention maintenance. Prove there is no blind full-universe refetch and all work obeys rights/rate/cost/value constraints.

### Catch-up delivery proof

If eligible maintenance windows are missed because the local app is off, prove checkpointed work resumes only in the next eligible overnight/weekend window and does not flood pre-market/regular-market/after-hours. Restart must not duplicate completed provider calls/work.

Machine contract: `adaptive-governance/PERSISTENCE_REUSE_AND_OFF_HOURS_DATA_READINESS_CONTRACT.md`.

## v19 delivery train — Professional Data Infrastructure

v19 begins only after v18 final closure. Deliver as small packets for: provider capability/entitlement/rights; provider quality/cost/coverage/SLOs; source disagreement/history/revision quality; 13F infrastructure; two-sided evidence substrate; AODR point-in-time outcome lineage; ADR-GDI professional reliability/capacity including protected-session reserve sizing and maintenance/preemption economics; measured specialized/paid-provider gap evaluation; v20 research-readiness evidence audit; then mandatory v19 Major Closure.

No v19 packet may recreate Smart Provider Router v2, persistence/cache ownership, session-aware maintenance, canonical identity, freshness or other owners established in v18. v19 measures and hardens them.

## v20 delivery train — Adaptive Intelligence & Decision Research

v20 begins only after v19 Major Closure. Deliver separately: adaptive research/experiment ledger; historical analogue/regime outcomes; calibration/FP-FN/miss/contradiction/drift; ASBI fingerprints/state transitions; ASBI scenarios/probability momentum; adaptive 13F; TDTI competing Long/Short/No Reliable Edge; TDTI two-sided trade-plan validation; AODR adaptive shared ranking; AODR diversity/personalized relevance after shared truth; ADR-GDI adaptive optimization including governed provider/maintenance usefulness and reserve-policy learning; model/prompt + Champion/Challenger governance; then v20 Professional Closure.

Adaptive production influence remains `SHADOW -> VALIDATED -> APPROVED -> PRODUCTION`. No silent self-promotion, no silent deterministic formula change and No Execution. Adaptive maintenance/provider optimization may not reduce protected live-session safety without governed evidence and explicit promotion.

## Delivery acceptance per patch

A patch is not complete merely because code compiles or CI is green. Delivery evidence must show:
- frozen G1 scope fully implemented;
- explicit non-goals untouched;
- deterministic regression tests for the changed responsibility;
- truthful degraded/partial/UNKNOWN/ABSTAIN states;
- no parallel owner introduced;
- no newly discovered implementation miss left without a durable target;
- packaged-runtime proof where the patch affects native/runtime behavior;
- current issue state and handoff agree with executable evidence.

For persistence/router changes, acceptance also requires DB/cache hit and residual-gap proof. For session-aware maintenance, acceptance additionally requires live-session non-degradation, provider reserve, preemption/checkpoint/resume and missed-window catch-up proof. For v19, acceptance additionally requires point-in-time provenance/data-rights/reliability truth where applicable. For v20, acceptance additionally requires leakage-safe evaluation, calibration/utility evidence, reproducibility, SHADOW/Champion-Challenger evidence and explicit rollback before promotion.

## Major-version handoff contract

- **v18 -> v19:** acquisition, persistence-first reuse, session-aware overnight/weekend readiness, protected live-session capacity, canonical identity, adaptive Market Modes, useful provider capabilities and provider telemetry are trustworthy and zero-gap enough for professional measurement.
- **v19 -> v20:** provider/data quality, rights, provenance, point-in-time feature/outcome lineage, reliability history and capacity/maintenance measurements are sufficient for governed adaptive research.
- A major-version closure cannot be bypassed by starting the next major's implementation early.

## Final v18.9.x closure

v18.9.12 must prove zero unexplained carry-forward, zero orphan useful provider capability, zero duplicate routing/freshness/persistence/session-scheduler/SEC/symbol/Market-Mode owner, deterministic Day/Swing/Long equivalence, #57 and #64 regressions, DB-first reuse, partial-gap acquisition, overnight/weekend maintenance, protected-session capacity reservation/preemption/recovery, Adaptive Intelligence Scorecard results and actual macOS Apple Silicon + Windows x64 packaged-runtime evidence.

## Permanent boundaries

U.S. Equities Processing; GLD/SLV/USO actionable exceptions; No Execution; Smart Provider Router v2 sole routing authority; canonical freshness/recovery sole freshness owner; existing multi-feed allocator sole subscription owner; BroadSnapshotBroker canonical reuse owner; canonical persistence/cache owners; canonical U.S. market calendar/session owner; direct SEC/EDGAR authoritative; existing telemetry/symbol/state owners; G0-G16 only.

## Exactly one next action

Obtain/reproduce the #64 macOS crash and freeze the narrow v18.9.1 scope. Do not create the v18.9.2 branch until v18.9.1 is truthfully closed or the crash is proven external/non-product.
