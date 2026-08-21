# DE.PULSE — Persistence Reuse & Session-Aware Data Readiness Contract

**Status:** APPROVED / PERMANENT ADAPTIVE DATA CONTRACT  
**Applies to:** v18.9.x, v19, v20 and later releases  
**Primary owners:** canonical persistence/cache owners + Smart Provider Router v2 + canonical U.S. market calendar/session/workload-priority owners  
**No parallel cache, database, router, scheduler or provider-specific data silo is permitted.**

## 1. Purpose

DE.PULSE must avoid paying provider/API/runtime cost for trustworthy evidence it already owns. Persistence and reuse apply to **all useful data**, not only OHLCV.

The system should accumulate a durable point-in-time evidence base so current consumers can reuse valid state immediately and future adaptive intelligence can learn from provenance-bound history.

Background preparation must never compromise the protected decision-support sessions: **pre-market, regular market and after-hours**. Those sessions retain first claim on provider quota/headroom, network, CPU, memory, DB/pool and workers.

## 2. Canonical request order

Where applicable:

`consumer requirement -> memory canonical cache -> persisted canonical DB/state -> validate coverage/freshness/schema/provenance/rights -> exact residual gap -> Smart Provider Router v2 ranks eligible providers for the gap -> targeted acquire -> normalize/reconcile -> persist -> serve`

Provider availability is not a reason to refetch. A successful provider response is not proof that the consumer requirement is complete.

## 3. Persistence-first / reuse-first rules

1. Reuse valid stored evidence before external fetch or recomputation.
2. Fetch only missing, expired, revised or materially insufficient evidence where provider semantics/rights allow targeted acquisition.
3. Immutable/effectively immutable history should be retained/reused where lawful.
4. Revision-prone records preserve original as-observed truth and later revisions; never silently rewrite point-in-time history.
5. Live/current-sensitive values require freshness semantics; stale quote/spread/VIX/market state can remain history but cannot be current truth.
6. Expensive/material derived intelligence may persist with input fingerprint, algorithm/model/prompt version, generated-at, provenance, freshness and invalidation rules.
7. Canonical de-duplication prevents provider-specific copies from becoming multiple truths.
8. Every retained dataset requires a named consumer, recovery/audit/learning value or material calls-avoided value.
9. Provider persistence/redistribution/commercial/AI-use rights override retention preferences.
10. Secrets remain in canonical secret owners and never enter ordinary market/evidence persistence.
11. Shared multi-user reuse occurs only inside compatible entitlement, rights and tenant/security domains.

## 4. Examples

### Historical price data
If a consumer needs 500 bars and 495 trustworthy bars already exist:

`DB 495 valid -> residual gap 5 -> acquire only missing 5 where supported -> reconcile -> persist -> serve 500`

### SEC / disclosures
Reuse already-acquired filings; check only for new/amended/corrected evidence or newly required fields. Direct SEC/EDGAR remains authoritative for filing truth.

### Fundamentals / earnings / macro
Retain point-in-time snapshots and later revisions. Refresh cadence follows real change/revision behavior.

### Company identity
Reuse canonical symbol/company/exchange/asset identity until material change or authoritative correction. Provider search is gap/corroboration, not mandatory repeated work.

## 5. Session-Aware Data Readiness Maintenance

DE.PULSE has one canonical **Data Readiness Maintenance** responsibility for low-priority, non-time-critical preparation.

### Tier 0 — Protected sessions
**Pre-market / Regular Market / After-Hours**

During protected sessions:
- current quote/liquidity/VIX/market context/catalyst/news/SEC/Opportunity Radar/Research/readiness work outranks background work;
- provider quota/headroom needed by current-session capabilities is reserved;
- low-priority provider acquisition suspends unless directly required by a live consumer;
- maintenance worker/CPU/memory/DB/network use is bounded;
- queued work is preemptible/checkpointed and yields promptly to current or market-shock work;
- no heavy compaction, deep reconciliation or broad historical fan-out.

### Tier 1 — Light Overnight Maintenance
Small, high-value work after protected after-hours and before the next protected pre-market window:
- latest completed-session persistence and bounded history gaps;
- incremental SEC/filing/amendment/disclosure checks;
- revision-due fundamentals/earnings/macro updates;
- small corporate-action/adjustment/identity corrections;
- bounded outcome resolution;
- lightweight integrity/readiness checks;
- bounded high-value derived precomputation with known provenance/version;
- next-session readiness-gap preparation.

Overnight work has conservative provider-call, concurrency, CPU, memory and DB budgets and drains before the pre-market protection buffer.

### Tier 2 — Heavy Weekend / Extended Market-Closed Maintenance
Deeper but bounded work:
- historical OHLCV gap/adjustment repair;
- corporate-action/split/dividend audits;
- SEC/Form 4/13F/congress history/amendment/identity reconciliation;
- earnings/fundamental/macro history and revisions;
- point-in-time Market Mode/Opportunity/Research/thesis/outcome lineage;
- provider usefulness/coverage/reliability history consolidation;
- material feature generation with known inputs/version/rights;
- DB integrity/index/statistics/compaction/retention work;
- high-value readiness coverage for actionable/My Market symbols.

Weekend does **not** mean full-universe blind acquisition.

## 6. Maintenance sequence

`inventory canonical state -> identify missing/expired/revision-due evidence -> prioritize by consumer/material value -> reserve protected-session capacity -> apply rights/rate/cost budgets -> acquire residual gaps through Smart Router -> reconcile/persist -> validate integrity/readiness -> record telemetry -> checkpoint/resume`

The existing canonical U.S. session/calendar owner defines all session/holiday/half-day boundaries. Maintenance cannot create another calendar truth.

## 7. Provider capacity reservation

For every provider/capability, workload policy should account for:
- entitlement and per-minute/day/month limits;
- current headroom/reset timing;
- expected protected-session demand;
- provider latency/error/circuit state;
- capability criticality and alternatives;
- maintenance value per call;
- historical calls-avoided benefit.

Maintenance budget is **bounded surplus after live reserve**. If uncertain, preserve capacity for current decision support.

## 8. Preemption, catch-up and restart

- light/heavy maintenance runs only in eligible windows and when provider/runtime health permits;
- use a configurable pre-market protection buffer from canonical session truth;
- missed work resumes only in a later eligible low-priority period;
- no backlog dump into pre-market/regular/after-hours;
- work is pauseable/preemptible/checkpointed;
- restart/resume must avoid duplicate acquisition/work;
- manual `Run Data Readiness Maintenance Now`, if exposed, obeys the same rights/budget/session protections;
- hosted deployments may run the same contract continuously but still obey protected-session priority and provider reserves.

## 9. Explicit exclusions

Maintenance/persistence must not:
- blind-refetch because market is closed;
- consume reserved live-session provider/runtime capacity;
- present stale live-sensitive data as current;
- overwrite original point-in-time evidence with later revisions;
- exceed provider entitlement/rate/cost/retention rights;
- persist/fetch data without a useful consumer/retention reason;
- create provider-specific scheduler/store/router/cache owners;
- create large background DB/CPU/network pressure near protected sessions;
- change protected deterministic Day/Swing/Long formulas;
- create execution capability.

## 10. Adaptive priority model

Prefer work that:
1. serves current actionable/My Market symbols;
2. closes a known consumer gap;
3. improves next-session readiness;
4. is expensive/rate-limited during protected sessions;
5. is likely to be reused by Research/Opportunity Radar/Market Modes;
6. improves point-in-time evidence needed by future learning;
7. resolves revisions/amendments/corporate-action correctness;
8. improves provider quality/rights/reliability measurement;
9. has high expected calls-avoided value per maintenance call.

## 11. Telemetry / acceptance

Record at minimum:
- DB/cache hits and records reused vs fetched;
- provider calls avoided;
- residual gaps detected/filled;
- new records vs revisions;
- provider selected and why;
- provider quota/headroom reserved;
- rate/cost/maintenance budget used;
- rights-blocked/skipped work;
- stale/invalid evidence detected;
- overnight/weekend work volume;
- duration/concurrency/CPU/memory/DB/network usage;
- preemption/reason/checkpoint/resume success;
- failures/deferred work and next eligible retry;
- readiness by capability/symbol cohort;
- effect on protected-session freshness/latency/degradation metrics.

Acceptance requires proof of **no material degradation** to protected decision-support sessions under realistic provider/rate/runtime pressure.

## 12. Version placement — aligned from v18.9 onward

### v18.9.3 — Coverage-Aware Smart Provider Router
Persisted canonical evidence becomes part of fulfillment before provider acquisition. The router computes residual gap, not the full theoretical request.

### v18.9.6 — Provider Observability / Adaptive Telemetry Foundation
**Architecture-audit alignment:** telemetry is intentionally earlier than provider capability expansion. It measures DB/cache reuse, calls avoided, residual-gap efficiency, provider usefulness, rate/freshness/runtime pressure and protected-session headroom so v18.9.7–v18.9.10 SHADOW provider admission is evidence-based.

### v18.9.11 — Session-Aware Data Readiness Maintenance
Implement the canonical light-overnight + heavy-weekend gap audit/backfill/reconciliation/catch-up coordinator using existing persistence, Smart Provider Router, freshness, provider budgets, telemetry and session/workload owners.

Acceptance includes overnight/weekend proof, provider/runtime reservation, preemption/drain/checkpoint/resume, missed-window catch-up, restart de-duplication, no broad blind refetch and no measurable decision-critical degradation.

### v18.9.12 — Professional Closure
Audit DB-first reuse, residual-gap acquisition, calls avoided, revision preservation, provider telemetry/admission, overnight/weekend behavior, protected-session reserves/preemption/recovery and actual supported packages.

### v19
Professionalize retention/rights/revisions, provider scorecards, PostgreSQL/index/pool/capacity behavior, hosted/shared persistence, sync eligibility, point-in-time lineage, maintenance economics and protected-session reserve sizing. New providers still use the same router/persistence/session-priority contract.

### v20
Adaptive intelligence consumes accumulated point-in-time evidence/outcomes and may learn provider/maintenance usefulness only through governed SHADOW/Champion-Challenger evaluation and explicit promotion. It cannot bypass provenance/freshness/rights or reduce live-session safety.

## 13. Permanent principle

**Live/current decision support first. Fetch once when useful -> normalize once -> persist truthfully -> reuse many times -> acquire only residual gaps -> prepare small gaps overnight -> perform deeper bounded repair on weekends -> reserve provider/runtime capacity for protected sessions -> learn only from provenance-bound point-in-time evidence.**
