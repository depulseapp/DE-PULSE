# DE.PULSE — Persistence Reuse & Off-Hours Data Readiness Contract

**Status:** APPROVED / PERMANENT ADAPTIVE DATA CONTRACT  
**Applies to:** v18.9.x, v19, v20 and later releases  
**Primary owners:** existing canonical persistence/cache owners + Smart Provider Router v2  
**No new parallel cache, database, router, scheduler or provider-specific data silo is permitted.**

## 1. Purpose

DE.PULSE must avoid paying provider/API/runtime cost for data it already owns in trustworthy canonical state. This applies to **all useful data**, not only historical OHLCV.

The system should accumulate a durable, point-in-time evidence base over time so current consumers can reuse valid data immediately and future adaptive intelligence can learn from provenance-bound history.

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

## 5. Off-Hours / Weekend Data Readiness Maintenance

DE.PULSE should have one canonical **Data Readiness Maintenance** activity for low-priority, non-time-critical data preparation. Weekend/off-hours windows are preferred because they reduce contention with live market work and allow provider budgets to be used more efficiently.

This is not a broad refetch job. It is **gap-driven maintenance**.

### Maintenance sequence

`inventory canonical data -> identify missing/expired/revision-due evidence -> prioritize by consumer value and future-use probability -> apply rights/rate/cost budgets -> Smart Router acquire residual gaps -> reconcile/persist -> validate DB integrity/readiness -> record telemetry`

### Eligible work

Where useful, lawful and contract-supported, maintenance may include:

- historical OHLCV gap/backfill and adjustment reconciliation;
- corporate actions and split/dividend reconciliation;
- SEC filings, Form 4 normalization/enrichment, 13F filing history and amendments;
- congressional/insider disclosure history;
- earnings events/results and historical reaction evidence;
- fundamentals snapshots and revision-aware history;
- symbol/company identity gaps and authoritative corrections;
- macro/FRED/BLS/EIA historical observations and revisions;
- historical Market Mode/regime evidence required for research;
- point-in-time Opportunity Radar / Research / thesis / readiness lineage where approved;
- subsequent outcome resolution for stored historical decisions/evidence;
- provider usefulness, coverage and reliability history;
- material derived features that are safe to precompute and whose inputs/version are known;
- DB integrity checks, index/statistics maintenance, bounded compaction/cleanup and retention enforcement where appropriate.

### Explicit exclusions

Maintenance must not:

- refetch everything merely because it is the weekend;
- treat stale quotes/spreads/live VIX as current Monday truth;
- overwrite original point-in-time evidence with later revisions;
- exceed provider entitlement/rate/cost budgets;
- fetch or persist data without a useful consumer/retention reason;
- create a TradeInsight-specific or provider-specific scheduler/data store;
- compete with higher-priority market-open/live work;
- change deterministic Day/Swing/Long formulas;
- create execution capability.

## 6. Scheduling and catch-up behavior

- Prefer bounded weekend/off-hours windows on U.S. non-trading periods.
- When the local app is not running during the preferred window, maintenance must remain resumable and may perform a bounded catch-up on the next eligible startup/off-hours period rather than assuming work happened.
- A future hosted deployment may execute the same canonical maintenance contract continuously without changing ownership or semantics.
- Maintenance must be pauseable/preemptible when a higher-priority market/session workload starts.
- Manual `Run Data Readiness Maintenance Now` may be exposed through Maintenance when implementation reaches the relevant patch, reusing the same coordinator and budgets.

## 7. Priority model

Maintenance priority should be adaptive and bounded. Prefer evidence that:

1. is required by current My Market / actionable symbols;
2. closes a known current consumer gap;
3. is expensive or rate-limited to obtain during live hours;
4. is likely to be reused by Research/Opportunity Radar/Market Modes;
5. improves point-in-time history needed by v20 learning;
6. resolves revisions/amendments/corporate-action correctness;
7. improves provider-quality/rights/reliability measurement.

Do not perform low-value full-universe work when the expected reuse/value does not justify provider/runtime cost.

## 8. Telemetry / acceptance

Record at minimum:

- database/cache hits;
- provider calls avoided;
- residual gaps detected;
- gaps filled;
- records reused vs fetched;
- new records vs revisions/amendments;
- provider selected and why;
- rate/cost budget used;
- rights-blocked/skipped work;
- stale/invalid records detected;
- maintenance duration and resource usage;
- failures/deferred work and next eligible retry;
- resulting readiness by capability/symbol cohort.

## 9. Patch placement

### v18.9.3 — Coverage-Aware Smart Provider Router
Must make persisted canonical evidence part of fulfillment before external provider acquisition. The router computes the **residual gap**, not the full theoretical request.

### v18.9.10 — Provider Efficiency / Adaptive Telemetry
Must measure DB/cache reuse, provider calls avoided, residual-gap acquisition efficiency and persistence usefulness.

### v18.9.11 — Off-Hours Data Readiness Maintenance
New single-responsibility patch. Implement the canonical bounded weekend/off-hours gap audit/backfill/reconciliation/catch-up coordinator using existing persistence, Smart Provider Router, freshness, rights, telemetry and workload-priority owners.

### v18.9.12 — Whole v18.9.x Professional Closure Audit
The former v18.9.11 closure moves to v18.9.12 so maintenance is not bundled into the closure audit. Closure must test DB-first reuse, partial-gap acquisition, weekend/catch-up behavior, provider-call avoidance, revision preservation and live-priority preemption.

Exact patch numbers remain subject to G1, but the **separation of responsibilities is mandatory**.

### v19
Professionalize retention/rights/revision policies, database/index/capacity behavior, point-in-time lineage, provider-quality scorecards and evidence-store readiness. A measured gap may justify specialized/paid data, but it plugs into the same router/persistence contract.

### v20
Adaptive intelligence consumes the accumulated point-in-time evidence/outcome store. It must not bypass provenance/freshness/rights or repair missing history by inventing confidence.

## 10. Permanent principle

**Fetch once when useful -> normalize once -> persist truthfully -> reuse many times -> refresh only what changed or expired -> learn from accumulated point-in-time evidence.**
