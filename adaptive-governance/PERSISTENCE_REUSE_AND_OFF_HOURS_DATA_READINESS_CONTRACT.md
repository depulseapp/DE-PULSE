# DE.PULSE — Persistence Reuse & Session-Aware Data Readiness Contract

**Status:** APPROVED / PERMANENT ADAPTIVE DATA CONTRACT  
**Applies to:** v18.9.x, v19, v20 and later releases  
**Primary owners:** existing canonical persistence/cache owners + Smart Provider Router v2 + existing U.S. market calendar/session/workload-priority owners  
**No new parallel cache, database, router, scheduler or provider-specific data silo is permitted.**

## 1. Purpose

DE.PULSE must avoid paying provider/API/runtime cost for data it already owns in trustworthy canonical state. This applies to **all useful data**, not only historical OHLCV.

The system should accumulate a durable, point-in-time evidence base over time so current consumers can reuse valid data immediately and future adaptive intelligence can learn from provenance-bound history.

DE.PULSE must also prepare useful non-live data proactively without compromising the three protected decision-support sessions: **pre-market, regular market and after-hours**. These sessions always receive first claim on provider quota/headroom, network concurrency, CPU, memory, database capacity and background-worker capacity.

## 2. Canonical request order

Every data consumer should follow this decision sequence where applicable:

`consumer requirement -> in-memory canonical cache -> persisted canonical DB/state -> validate coverage/freshness/schema/provenance/rights -> compute residual gap -> Smart Provider Router v2 ranks eligible providers for only the gap -> acquire -> normalize/reconcile -> persist -> serve`

A provider must not be called merely because it exists or because a fixed fallback order names it. Existing trustworthy data is a first-class source of fulfillment.

## 3. Persistence-first / reuse-first rules

1. Reuse valid stored evidence before external fetch or recomputation.
2. Fetch only the missing, expired, revised or materially insufficient portion where provider contracts permit targeted acquisition.
3. A successful provider response does not imply complete consumer coverage.
4. Immutable or effectively immutable historical records should be retained and reused where provider/data rights allow.
5. Revision-prone records must preserve point-in-time/as-observed truth and later revisions rather than silently overwriting history.
6. Current/live-sensitive values require freshness TTLs; an old quote, spread, VIX value or market state may be useful as history but cannot be presented as current truth.
7. Expensive/material derived intelligence may be persisted when useful, together with input fingerprint, algorithm/model/version, generated-at time, provenance, freshness and invalidation rules.
8. Canonical de-duplication prevents the same underlying evidence from becoming multiple provider-specific truths.
9. Every retained dataset requires a named consumer, recovery/audit/learning value, or material provider-call-saving value.
10. Provider persistence, redistribution, commercial-use and AI-use rights override retention preferences.
11. Secrets remain in the canonical secret owner and never enter ordinary market/evidence persistence.

## 4. Examples

### Historical price data

If APP requires 500 daily bars and canonical persistence already contains 495 valid bars:

`DB: 495 valid -> gap: 5 -> Smart Router acquires only missing 5 where supported -> merge -> persist -> consumer gets 500`

The system must not download all 500 again solely because an API call is available.

### SEC / disclosure evidence

A previously acquired filing/disclosure is reused from canonical persistence. Subsequent work checks only for newer filings, amendments, corrections or newly required normalized fields. Direct SEC/EDGAR remains authoritative for filing truth.

### Fundamentals / earnings / macro

Historical snapshots are retained point-in-time. Refresh cadence follows the capability's real change/revision behavior. Old observations remain available for research even when current truth has advanced.

### Company identity

Known symbol/company/exchange/asset identity is reused until a material change or authoritative correction occurs; provider search is a gap/corroboration path, not a mandatory repeated call.

## 5. Session-Aware Data Readiness Maintenance

DE.PULSE should have one canonical **Data Readiness Maintenance** activity for low-priority, non-time-critical data preparation. It has two operating tiers:

- **LIGHT OVERNIGHT MAINTENANCE** on normal trading-day cycles, after protected after-hours work has ended and before the next protected pre-market session begins.
- **HEAVY WEEKEND / MARKET-CLOSED MAINTENANCE** during larger non-trading windows, with deeper but still bounded backfill/reconciliation work.

The maintenance coordinator must use the existing canonical U.S. market calendar/session engine. It must not hard-code a second definition of pre-market, regular market, after-hours, holidays, half-days or exceptional closures.

This is not a broad refetch job. It is **gap-driven, value-driven and resource-budgeted maintenance**.

### Maintenance sequence

`inventory canonical data -> identify missing/expired/revision-due evidence -> prioritize by consumer value and future-use probability -> reserve protected live-session capacity -> apply rights/rate/cost budgets -> Smart Router acquire residual gaps -> reconcile/persist -> validate DB integrity/readiness -> record telemetry -> checkpoint/resume`

## 6. Protected-session priority contract

### Tier 0 — Pre-Market / Regular Market / After-Hours

These are protected high-priority operating sessions. During them:

- live/current quote, spread, liquidity, VIX/market context, catalyst/news/SEC/event, Opportunity Radar, Research/readiness and other decision-critical workloads outrank maintenance;
- provider quota/headroom needed for current-session capability must be reserved and unavailable to low-priority maintenance;
- background maintenance external-provider calls are suspended unless the same acquisition is directly required by a live/current consumer, in which case it is live fulfillment rather than maintenance;
- maintenance worker concurrency, CPU, memory, DB writes and network usage must be bounded so they cannot materially increase current-session latency or produce self-inflicted `DATA DEGRADED`;
- queued maintenance must be preemptible/checkpointed and yield promptly when a protected session or market-shock/high-priority workload begins;
- no heavy compaction, deep reconciliation or large historical fan-out runs are allowed.

### Tier 1 — Light Overnight Maintenance

Purpose: keep the app ready for the next pre-market session with small, high-value work only.

Preferred work includes:

- finalize/fill the latest completed daily/intraday historical gaps needed by actionable/My Market symbols;
- persist completed session summaries/outcomes and resolve bounded pending historical outcomes;
- check incremental SEC/filing/amendment/disclosure changes;
- update revision-due fundamentals/earnings/macro data when their real cadence warrants it;
- reconcile small corporate-action/adjustment deltas;
- fill high-value company/symbol identity gaps;
- compact/checkpoint small canonical persistence queues;
- perform lightweight integrity/readiness checks;
- precompute only bounded material derived features with known input/version provenance;
- calculate next-session readiness gaps so pre-market starts from warm canonical state.

Overnight work must have conservative provider-call, concurrency, CPU, memory and DB-write budgets. It must stop early if provider headroom, runtime health or the approaching pre-market protection buffer says to stop.

### Tier 2 — Heavy Weekend / Extended Market-Closed Maintenance

Purpose: use the larger non-trading window for deeper work that is useful but inappropriate during daily live operations.

Eligible work includes:

- deeper historical OHLCV backfill/gap repair and adjustment reconciliation;
- corporate-action/split/dividend audits;
- deeper SEC/Form 4/13F history, amendments and identity reconciliation;
- congressional/insider disclosure historical gap repair;
- earnings/fundamental history backfill and revision reconciliation;
- macro historical revision reconciliation;
- historical Market Mode/regime evidence preparation;
- point-in-time Opportunity Radar / Research / thesis / readiness lineage and subsequent outcome resolution;
- provider usefulness/coverage/reliability history consolidation;
- material feature generation whose inputs/version/rights are known;
- bounded DB integrity checks, index/statistics maintenance, compaction/cleanup and retention enforcement;
- pre-building high-value research/evidence coverage for current actionable/My Market symbols and other bounded high-probability consumers.

Heavy maintenance is still not permission for blind full-universe acquisition. Expected reuse, materiality, provider budgets, storage rights and runtime cost remain mandatory.

## 7. Provider Capacity Reservation

Maintenance must never consume provider capability needed to keep DE.PULSE first-class during protected sessions.

For every capability/provider, Smart Provider Router/workload policy should account for:

- entitlement and per-minute/day/month quotas;
- current rate-limit headroom;
- reset timing;
- expected next protected-session demand;
- provider latency/error/circuit state;
- capability criticality;
- alternative-provider availability;
- maintenance value per provider call;
- historical calls-avoided benefit.

The maintenance budget is **surplus/bounded capacity after protected-session reserve**, not the other way around. If there is uncertainty, preserve capacity for live/current work.

A provider with scarce quota may receive zero overnight/weekend maintenance traffic even if technically available. A provider with abundant cheap historical entitlement may receive more gap-backfill work when useful and lawful.

## 8. Explicit exclusions

Maintenance must not:

- refetch everything merely because the market is closed;
- consume provider quota/headroom reserved for pre-market, regular market or after-hours;
- treat stale quotes/spreads/live VIX as current next-session truth;
- overwrite original point-in-time evidence with later revisions;
- exceed provider entitlement/rate/cost budgets;
- fetch or persist data without a useful consumer/retention reason;
- create a TradeInsight-specific or provider-specific scheduler/data store;
- compete with higher-priority current-session or market-shock work;
- create large background DB/CPU/network pressure near a protected-session boundary;
- change deterministic Day/Swing/Long formulas;
- create execution capability.

## 9. Scheduling, preemption and catch-up behavior

- Session windows come from the existing canonical U.S. market calendar/session owner.
- Light overnight maintenance runs only when the app is in an eligible low-priority period and runtime/provider health permits it.
- Heavy weekend maintenance runs only in sufficiently large eligible market-closed windows.
- Use a configurable **pre-market protection buffer** so maintenance drains/checkpoints before pre-market work must become fully responsive; do not hard-code a competing calendar definition.
- If the local app was not running, missed work remains resumable. Perform bounded catch-up only in the next eligible overnight/weekend window; do not dump missed maintenance into a protected live session.
- A future hosted deployment may execute the same canonical maintenance contract continuously while still obeying session priority and provider reserves.
- Maintenance must be pauseable/preemptible and checkpoint progress so partial work is not repeated unnecessarily.
- Manual `Run Data Readiness Maintenance Now` may be exposed through Maintenance, but it must still obey rights, provider reserves, runtime protection and session-priority rules; manual action cannot override live-safety limits.

## 10. Priority model

Maintenance priority should be adaptive and bounded. Prefer evidence that:

1. is required by current My Market / actionable symbols;
2. closes a known current consumer gap;
3. improves readiness for the next pre-market/regular/after-hours cycle;
4. is expensive or rate-limited to obtain during protected live sessions;
5. is likely to be reused by Research/Opportunity Radar/Market Modes;
6. improves point-in-time history needed by v20 learning;
7. resolves revisions/amendments/corporate-action correctness;
8. improves provider-quality/rights/reliability measurement;
9. has high expected calls-avoided benefit per maintenance call.

Do not perform low-value full-universe work when expected reuse/value does not justify provider/runtime cost.

## 11. Telemetry / acceptance

Record at minimum:

- database/cache hits;
- provider calls avoided;
- residual gaps detected;
- gaps filled;
- records reused vs fetched;
- new records vs revisions/amendments;
- provider selected and why;
- provider quota/headroom reserved for protected sessions;
- maintenance provider/rate/cost budget used;
- rights-blocked/skipped work;
- stale/invalid records detected;
- overnight vs weekend work performed;
- maintenance duration, concurrency, CPU/memory/DB/network resource use;
- preemption count/reason and checkpoint/resume success;
- failures/deferred work and next eligible retry;
- resulting readiness by capability/symbol cohort;
- effect on protected-session latency/freshness/degradation metrics.

Acceptance requires proof that maintenance produces **no material degradation** to pre-market, regular-market or after-hours freshness/latency/readiness under realistic provider/rate-limit/runtime pressure.

## 12. Patch placement

### v18.9.3 — Coverage-Aware Smart Provider Router
Must make persisted canonical evidence part of fulfillment before external provider acquisition. The router computes the **residual gap**, not the full theoretical request.

### v18.9.10 — Provider Efficiency / Adaptive Telemetry
Must measure DB/cache reuse, provider calls avoided, residual-gap acquisition efficiency, persistence usefulness and the provider/runtime headroom required to protect live sessions.

### v18.9.11 — Session-Aware Data Readiness Maintenance
Single-responsibility patch. Implement the canonical **light overnight + heavy weekend** gap audit/backfill/reconciliation/catch-up coordinator using existing persistence, Smart Provider Router, freshness, provider budgets, U.S. market-session calendar, telemetry and workload-priority owners.

Acceptance includes:

- overnight light-maintenance proof;
- weekend heavy-maintenance proof;
- pre-market/regular/after-hours provider-capacity reservation;
- preemption/drain/checkpoint/resume proof;
- missed-window bounded catch-up proof;
- no duplicate work after restart;
- no measurable decision-critical live-session degradation;
- no broad blind refetch.

### v18.9.12 — Whole v18.9.x Professional Closure Audit
Closure must test DB-first reuse, partial-gap acquisition, overnight/weekend behavior, provider-call avoidance, revision preservation, protected-session capacity reservation, live-priority preemption and recovery.

Exact patch numbers remain subject to G1, but the **separation of responsibilities is mandatory**.

### v19
Professionalize retention/rights/revision policies, database/index/capacity behavior, point-in-time lineage, provider-quality scorecards and evidence-store readiness. Measure maintenance value, quota economics, capacity reservations and whether paid/specialized providers are justified by observed gaps. New providers still plug into the same router/persistence/session-priority contract.

### v20
Adaptive intelligence consumes the accumulated point-in-time evidence/outcome store and may learn better maintenance/provider usefulness policies through governed SHADOW/Champion-Challenger evaluation. It must not bypass provenance/freshness/rights, self-promote, or sacrifice live-session truth for background learning.

## 13. Permanent principle

**Live/current decision support first. Fetch once when useful -> normalize once -> persist truthfully -> reuse many times -> fill small gaps overnight -> perform deeper bounded repair on weekends -> always reserve provider/runtime capacity for pre-market, regular market and after-hours -> learn from accumulated point-in-time evidence.**
