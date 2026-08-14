# DE.PULSE — 10/10 Adaptive Data Reliability & Graceful Degradation Intelligence (ADR-GDI)

**Status:** APPROVED / PERMANENT / GOVERNING  
**Applies to:** Adaptive Roadmap, Adaptive Build Plan, Adaptive Build Process, Adaptive Delivery Process, Provider Router, Shared Symbol Intelligence, SQLite/PostgreSQL persistence, Day/Swing/Long, Research, Opportunity Radar, Decision Queue, ASBI, TDTI, Maintenance/Data Engine  
**Primary goal:** `DATA DEGRADED` must mean that decision-relevant evidence is genuinely degraded—not that DE.PULSE overloaded itself, retried wastefully, rebuilt already-known state, or allowed a non-critical dataset to contaminate the whole product.

This contract extends the permanent `governance/ADAPTIVE-OPERATING-CONTRACT.md`. It does **not** add another top-level gate; all implementation and certification responsibilities remain inside G0–G16.

---

# 1. North Star

DE.PULSE must be resilient, truthful, selective and self-protecting.

Preferred reliability loop:

**Observe → classify capability health → measure freshness/impact → protect critical work → reuse warm canonical state → route/fallback intelligently → shed low-value load → explain degradation precisely → recover with hysteresis → record outcome → learn → validate → adapt**

The system should remain useful under partial failure. It must not turn every provider hiccup, stale low-priority dataset, background job, or optional-context gap into a broad `DATA DEGRADED` state.

The target standard is:

> **Critical decision evidence stays fast and current; failures are isolated to the smallest truthful capability/consumer scope; optional evidence reduces context/confidence rather than collapsing the whole ticker; local overload is prevented before it reaches users; genuine decision-critical degradation produces explicit ABSTAIN / NO RELIABLE EDGE where required.**

---

# 2. ADR-001 — Capability-Level Health, Not One Global Provider Flag

Health must be tracked at the smallest useful capability/dataset level, not only by provider or app-wide state.

Examples:
- `EQUITY_LIVE_QUOTE`;
- `INTRADAY_BARS`;
- `DAILY_HISTORY`;
- `BID_ASK_LIQUIDITY`;
- `NEWS`;
- `EARNINGS`;
- `SEC_FILINGS`;
- `FUNDAMENTALS`;
- `OPTIONS_CONTEXT`;
- `MACRO`;
- `VIX`;
- `13F`;
- `SHORT_INTEREST / BORROW / SHORTABILITY` where available;
- `AI_SYNTHESIS`;
- persistence/query capability.

A provider may be healthy for one capability and limited, rate-limited, stale or not entitled for another. DE.PULSE must not collapse these into a misleading single provider state.

---

# 3. ADR-002 — Consumer / Decision Dependency Graph

Every decision surface must declare which capabilities are:

- **REQUIRED** — absence/staleness can make the decision unsafe or invalid;
- **MATERIAL WHEN PRESENT** — relevant context that may change confidence/quality but is not always mandatory;
- **OPTIONAL CONTEXT** — useful enrichment whose absence should not degrade the core decision;
- **NOT RELEVANT** — must not affect that surface.

Dependency is horizon/context aware.

Illustrative Day example:

`Day Thesis`
→ live price `REQUIRED`
→ intraday bars `REQUIRED`
→ liquidity/spread `REQUIRED` when necessary for tradeability
→ current material catalyst/news `MATERIAL WHEN PRESENT`
→ fundamentals `OPTIONAL / LOW PRIORITY`
→ 13F `CONTEXT ONLY`

Illustrative Long-Term example may invert several priorities.

A capability can degrade only the consumers that actually depend on it.

---

# 4. ADR-003 — Freshness SLOs Are Dataset + Horizon + Session Specific

There is no universal stale timeout.

Define explicit freshness/SLO expectations by:
- dataset/capability;
- horizon/consumer;
- market session;
- provider capability;
- event sensitivity;
- materiality.

Examples:
- live quotes may tolerate seconds;
- intraday bars may tolerate a small bounded delay depending on use;
- fundamentals may remain useful for hours;
- 13F is inherently quarterly/lagged and must be judged by filing/report-period truth rather than live-quote freshness.

The UI must never label naturally slow-moving data as degraded merely because it is older than a live-feed threshold.

---

# 5. ADR-004 — Canonical Degradation Reason Taxonomy

Every degraded/limited state must have an explicit reason code. Approved examples include:

- `STALE`;
- `MISSING`;
- `NOT_ENTITLED`;
- `RATE_LIMITED`;
- `PROVIDER_DOWN`;
- `NETWORK_FAILURE`;
- `SOURCE_DISAGREEMENT`;
- `MALFORMED_RESPONSE`;
- `CACHE_MISS`;
- `WARMING`;
- `LOCAL_OVERLOAD`;
- `QUEUE_SATURATED`;
- `DB_SLOW`;
- `DB_UNAVAILABLE`;
- `MISSING_HISTORY`;
- `LOW_COVERAGE`;
- `RIGHTS_RESTRICTED`;
- `UNKNOWN`.

Do not show generic `DATA DEGRADED` when a more precise reason is known.

Each event should preserve capability, preferred/serving source, timestamps, age, impact, fallback/recovery state and affected consumers where relevant.

---

# 6. ADR-005 — Blast-Radius / Impact Intelligence

The system must calculate **what the degradation actually affects**.

Examples:
- `OPTIONS_CONTEXT delayed → Swing context reduced → deterministic core plan unaffected`;
- `FUNDAMENTALS stale → Long confidence may reduce → Day unaffected`;
- `LIVE_QUOTE stale → Day/Swing current readiness unsafe → affected horizons degrade/ABSTAIN`;
- `13F unavailable → institutional context unavailable → no broad market-data failure`.

Impact should propagate through the dependency graph rather than through a global red flag.

---

# 7. ADR-006 — Health / Pressure State Machine

Runtime reliability should use an explicit state model such as:

**NORMAL → PRESSURE → PROTECTED → DEGRADED → RECOVERING → HEALTHY**

Meaning:
- **NORMAL/HEALTHY:** within SLO/headroom;
- **PRESSURE:** capacity or latency is worsening; start reducing low-value work before users are affected;
- **PROTECTED:** critical workloads remain protected while lower-priority/background work is throttled/deferred;
- **DEGRADED:** required capability/consumer SLO cannot be met truthfully;
- **RECOVERING:** service has resumed but hysteresis/stability proof is still required;
- **HEALTHY:** recovery is sustained and SLO restored.

A single successful response must not instantly erase a genuine degraded state.

---

# 8. ADR-007 — Workload Priority / Critical Evidence Protection

Under pressure, prioritize by decision value and time sensitivity.

Illustrative priority order:
1. selected ticker / current user decision-critical live price and liquidity;
2. active Day candidates, Decision Queue and Rapid Move material events;
3. immediate earnings/material catalysts and Market Open/Pre-Market readiness;
4. Swing/current Research consumers;
5. Long-Term/current contextual work;
6. background historical hydration/reconciliation;
7. 13F bulk history, calibration, replay, maintenance and other deferrable heavy analytics.

The exact tiers may adapt, but low-value/background work must yield before critical live/current evidence is allowed to miss SLO because of local contention.

---

# 9. ADR-008 — Persistent Warm Canonical State

SQLite desktop and PostgreSQL hosted/shared state are reliability tools, not goals by themselves.

Persist only useful canonical state that reduces restart storms, duplicate provider work, repeated recalculation or loss of decision lineage.

Warm-start expectations include, where lawful/useful:
- canonical latest evidence + freshness/provenance;
- historical bars/features needed for current decisions;
- capability/entitlement/circuit memory where safe;
- relevant normalized filings/news/fundamentals/13F history;
- outcome/validation lineage;
- provider usefulness/reliability telemetry;
- degradation/recovery state needed for truthful restart behavior.

Restart should preferentially restore useful warm state and selectively refresh what actually needs refreshing rather than blindly refetching everything.

---

# 10. ADR-009 — Fetch Once / Calculate Once / Single-Flight / Coalescing

Equivalent requests/calculations for the same canonical symbol/dataset/time window must be coalesced where semantically valid and lawful.

Required principles:
- one canonical owner;
- single-flight for in-progress equivalent requests;
- request coalescing;
- shared result fan-out;
- material-change propagation;
- incremental recomputation;
- bounded caches;
- affected-consumer invalidation only.

If five modules simultaneously discover NVDA history is stale, DE.PULSE should normally execute one canonical refresh, not five provider requests.

---

# 11. ADR-010 — Smart Router / Circuit / Retry Discipline

Provider Router remains the sole executable provider-routing authority.

Reliability behavior must use:
- capability/entitlement memory;
- Preferred vs Serving semantics;
- deterministic cooldowns for `NOT_ENTITLED` and repeated hard failures;
- 429/rate-limit awareness;
- per-capability circuit breakers;
- provider saturation/headroom;
- source disagreement handling;
- exponential/adaptive backoff where appropriate;
- jitter where useful;
- provider calls avoided telemetry.

Never create retry storms by repeatedly probing a provider/capability already known to be unavailable or rate limited.

---

# 12. ADR-011 — Graceful Load Shedding / Backpressure

Queues, goroutines/workers, subscriptions, DB work and background jobs must be bounded.

When pressure rises:
- slow/defer low-priority refreshes;
- pause deferrable reconciliation/calibration/history jobs;
- reduce redundant UI recomputation;
- preserve critical live subscriptions/calls;
- reject/drop only low-value work using explicit policy;
- retain truthful freshness/impact metadata.

Load shedding must never fabricate health or silently discard decision-critical evidence.

---

# 13. ADR-012 — Graceful Intelligence Degradation / ABSTAIN

Missing optional evidence should not automatically invalidate a decision. Instead, where appropriate:
- reduce confidence;
- mark the specific context unavailable/stale;
- explain impact;
- continue using independent healthy evidence.

Missing or unreliable **required** evidence must cause the affected decision surface to degrade or return **ABSTAIN / NO RELIABLE EDGE** rather than produce false confidence.

ASBI and TDTI must consume capability/dependency health so probabilities, confidence, opportunity quality and readiness never pretend evidence is current when it is not.

---

# 14. ADR-013 — Source Disagreement Is Not a Provider Outage

When trustworthy providers disagree materially:
- preserve both observations and timestamps;
- apply dataset-specific tolerance;
- consult approved validation sources when justified;
- identify canonical winner/reason where defensible;
- otherwise retain `SOURCE_DISAGREEMENT` and conservative decision semantics.

Do not convert unresolved disagreement into an arbitrary first-response-wins result or a misleading whole-provider failure.

---

# 15. ADR-014 — Database Reliability / PostgreSQL Must Not Become the New Bottleneck

Database introduction does not by itself solve `DATA DEGRADED`.

For SQLite/PostgreSQL track applicable:
- query p50/p95/p99 latency;
- slow queries;
- connection-pool utilization/saturation;
- lock/contention/deadlock behavior;
- transaction latency/failures;
- index usage/effectiveness;
- cache/buffer effectiveness where observable;
- write amplification/background write pressure;
- storage growth;
- checkpoint/vacuum/maintenance impact where applicable;
- retry/failover/recovery behavior;
- migration/backup/restore correctness.

Persistence must remain asynchronous/background where safe and must not block current live evidence unnecessarily.

---

# 16. ADR-015 — Reliability / Degradation Event Ledger

Persist meaningful reliability events so the system can learn from actual operations.

Record applicable:
- event/reason ID;
- capability/dataset;
- provider/source;
- affected symbol(s)/consumer(s);
- start/end/recovery timestamps;
- freshness/SLO breach;
- pressure/load state;
- fallback used;
- queue/DB/runtime metrics;
- provider calls avoided;
- user/decision impact;
- recovery result;
- repeated-pattern identity.

Do not indiscriminately store noisy telemetry forever; retain enough structured evidence for diagnosis, trend analysis, capacity planning and adaptive improvement.

---

# 17. ADR-016 — Adaptive Reliability Learning

DE.PULSE should learn from reliability outcomes, but production behavior must remain governed.

Potential adaptive learning includes:
- provider/capability recovery-time distributions;
- failure recurrence patterns;
- useful cooldown/backoff ranges;
- dataset/provider usefulness;
- fallback quality;
- calls avoided;
- workload cost vs decision value;
- cache usefulness;
- DB query/index hot spots;
- which low-priority jobs contribute to pressure;
- which degradation alerts were actionable vs noisy.

Adaptive routing/scheduling changes with production impact follow **SHADOW → VALIDATED → APPROVED → PRODUCTION**. No silent self-modifying reliability policy.

---

# 18. ADR-017 — User-Facing Reliability UX

Normal users should see concise impact-oriented status, not infrastructure noise.

Preferred examples:

> **Market data healthy**  
> Fundamentals refresh delayed by provider rate limit · cached data 14h old · Day/Swing live evidence unaffected.

> **Options Context delayed**  
> Core price/history/news/SEC evidence current · Swing confidence reduced slightly.

> **Decision data degraded**  
> Live quote freshness exceeded required SLO · no trustworthy fallback currently available · Day readiness withheld.

Rules:
- capability + reason + freshness + impact + serving fallback when relevant;
- do not use red/global severity for optional context failures;
- no raw queue/database/circuit machinery in ordinary USER/DEMO surfaces;
- deeper diagnostics belong in Maintenance/Data Engine/admin views;
- do not oscillate status due to one transient success/failure.

---

# 19. ADR-018 — Maintenance / Data Engine Diagnostics

Detailed diagnostics should support operators without overwhelming normal users.

Expose applicable:
- capability health;
- Preferred vs Serving provider;
- entitlement;
- reason/circuit/cooldown;
- freshness/SLO age;
- p50/p95/p99 provider latency;
- error/429/malformed rates;
- request/subscription budget/headroom;
- queue depth/age;
- pressure/load-shed state;
- DB pool/query/lock indicators;
- cache hit/usefulness;
- calls avoided;
- last successful canonical contribution;
- affected consumers/blast radius;
- recovery state/history.

---

# 20. ADR-019 — Cross-Module Integration

ADR-GDI must integrate with existing canonical owners rather than create a parallel health engine.

Consumers include as relevant:
- Dashboard;
- Day / Swing / Long;
- Research;
- Opportunity Radar;
- Decision Queue;
- Market/Event Intelligence;
- Rapid Move;
- ASBI;
- TDTI;
- AI synthesis;
- Pre-Market / Market Open Prep;
- Maintenance/Data Engine;
- hosted/multi-user services.

Each consumer receives capability/freshness/impact truth from shared canonical health state.

---

# 21. ADR-020 — Reliability Scorecard / SLOs

Measure at least the applicable:
- decision-critical data SLO attainment;
- capability availability by session;
- stale-data incidence/duration;
- local-overload degradation incidence;
- provider-originated degradation incidence;
- time in NORMAL/PRESSURE/PROTECTED/DEGRADED/RECOVERING;
- mean/p50/p95 time-to-recovery;
- fallback success/usefulness;
- queue saturation/oldest-work age;
- DB p50/p95/p99 and pool saturation;
- duplicate requests/calculations prevented;
- calls avoided;
- cache usefulness;
- restart warm-up time/provider-call burst;
- false degradation / noisy alert rate;
- blast-radius accuracy;
- ABSTAIN correctness when required data is missing;
- critical-work latency during background load;
- long-run/soak stability;
- storage growth;
- adaptive policy Champion/Challenger evidence where used.

The objective is not to eliminate truthful degradation messages. It is to eliminate **self-inflicted, overly broad, unexplained or unnecessary degradation**.

---

# 22. ADR-021 — Adaptive Roadmap Placement

### v18.3 — Persistent / Hosted Shared-State Foundation
Use the PostgreSQL/shared-state work to harden:
- warm canonical state;
- indexed persistence/query paths;
- connection pooling/concurrency;
- shared symbol reuse;
- capability health persistence where appropriate;
- workload/DB observability;
- restart recovery;
- pressure/backpressure integration.

Do not wait for PostgreSQL before fixing obvious provider retry storms, duplicate work, capability-health semantics or priority/load-shedding defects that are dependency-compatible earlier.

### v18.5 — Mandatory Reliability Closure
Before v19, prove that local/runtime architecture does not materially cause broad `DATA DEGRADED` under realistic supported load.

Required closure includes provider-load, concurrency, restart/warm-start, DB parity, queue/backpressure, load shedding, long-run/soak, storage growth, capability blast radius, degraded/recovery UX, and actual packaged runtime evidence.

If local overload can delay/misstate live/current decision-critical evidence, it is a **release blocker** until fixed or explicitly constrained with truthful limits.

### v19 — Professional Data Infrastructure
Harden provider/database reliability with measured capability SLOs, data-quality/rights/cost matrices, long-run capacity, query/index tuning, degradation history, fallback quality and commercial/hosted operating limits.

### v20 — Adaptive Reliability Optimization
Use accumulated reliability history for governed adaptive provider/recovery/workload optimization, Champion/Challenger evaluation, improved failure prediction and smarter prioritization. v20 is **not** the first time basic reliability is fixed; it is where learned optimization matures.

---

# 23. ADR-022 — Adaptive Build Plan

When selecting ADR-GDI work, prefer this dependency order:

1. instrument and classify current `DATA DEGRADED` causes before guessing;
2. establish capability-level health and canonical reason codes;
3. map decision-consumer dependencies and freshness SLOs;
4. eliminate duplicate fetch/calculation/retry storms using canonical single-flight/coalescing;
5. enforce priority/backpressure/load shedding;
6. add persistent warm-state/restart recovery and DB observability;
7. implement blast-radius-aware UI/consumer semantics;
8. persist meaningful reliability/degradation outcomes;
9. optimize provider/DB/runtime capacity;
10. introduce adaptive recovery/routing/scheduling only with SHADOW evidence;
11. validate under realistic supported load before promotion.

Build order may adapt when a measured bottleneck or safety issue makes another dependency more urgent.

---

# 24. ADR-023 — Adaptive Build Process — G0–G16 Responsibilities

ADR-GDI uses the existing G0–G16 model only.

- **G0 Exact Baseline:** capture exact current degradation symptoms, reason coverage, provider/runtime/DB/load metrics, open reliability defects and source fingerprint.
- **G1 Immutable Scope:** freeze the reliability slice, affected capabilities/consumers, SLOs and acceptance criteria.
- **G2 Architecture / Data Utility:** prove canonical health/data owners, dependency graph, no duplicate router/cache/store/health engine, and justify persisted telemetry/data value.
- **G3 Design / Dependency Readiness:** freeze reason taxonomy, SLO/freshness semantics, priority/load-shed policy, fallback/recovery behavior, DB/provider dependencies and test oracle.
- **G4 Development Exit:** unit/static/schema tests for reason mapping, dependency impact, freshness, circuit/retry, queue bounds, recovery hysteresis, persistence and restart semantics.
- **G5 FAST Qualification:** verify HEALTHY/PRESSURE/DEGRADED/RECOVERING paths, optional-vs-required dependency behavior and basic fallback.
- **G6 Integration / MEDIUM Qualification:** verify Provider Router, Shared Symbol Intelligence, SQLite/PostgreSQL, Day/Swing/Long, ASBI/TDTI, Research/Queue/Radar and Data Engine integration without duplicate work.
- **G7 Data / Security / Adaptive Intelligence:** enforce provenance, truthful stale/missing/UNKNOWN, rights, no false recovery, ABSTAIN, event-ledger integrity and governed adaptive reliability policy.
- **G8 Performance / Capacity / Stability:** profile provider calls, duplicate work, CPU/memory/GC/locks, queue depth, DB p95/p99, pool saturation, UI/API latency, restart burst, load shedding, long-run soak and storage growth under realistic load.
- **G9 Cross-Module / UI / UX:** prove capability-specific impact messaging, no broad false degradation, accessible reason/freshness/impact language and separation of USER vs diagnostic machinery.
- **G10 Pre-Freeze Qualification:** run failure injection and supported-load scenarios before RC freeze; unresolved self-inflicted decision-critical degradation blocks freeze.
- **G11 Immutable RC:** freeze exact reliability policy/config/SLO/reason taxonomy/source identity.
- **G12 Full Certification:** run full regression plus provider 429/outage/slow response, stale cache, source disagreement, DB slow/unavailable, queue saturation, restart/warm-start, multi-user/symbol fan-out, background-job pressure and long-run recovery cases.
- **G13 Native Packaging / Provenance:** package the exact certified policy/config/schema/migrations and provenance.
- **G14 Actual Artifact Runtime Audit:** verify the packaged app under real runtime conditions, including restart, persistence, provider failure, degraded UI, fallback, recovery and no false HEALTHY state.
- **G15 Release Assurance / Promotion:** verify rollback/reproducibility, supported operating limits, production-vs-SHADOW reliability policy identity and no silent adaptive promotion.
- **G16 Adaptive Retrospective / Handoff:** feed degradation cause/duration, provider/DB/runtime bottlenecks, calls avoided, alert usefulness, recovery performance and capacity findings into the next Adaptive Build Plan.

---

# 25. ADR-024 — Adaptive Delivery Process / Stable Acceptance

A Stable artifact satisfies ADR-GDI only when the **actual packaged runtime** proves, as applicable:
- critical decision evidence retains priority under realistic supported load;
- optional capability failures do not incorrectly degrade unrelated decisions;
- required evidence failures cause truthful affected-scope degradation/ABSTAIN;
- capability/reason/freshness/impact/fallback semantics are correct;
- provider retries/circuits do not create storms;
- equivalent work is coalesced/reused;
- warm restart avoids unnecessary provider rebuild storms;
- queues/caches/workers are bounded;
- load shedding protects critical work;
- DB pooling/query/index/persistence paths stay within defined operating limits;
- recovery uses hysteresis and does not flap;
- degradation/recovery history survives restart where required;
- AI/ASBI/TDTI never treat stale/missing evidence as healthy/current;
- normal users receive concise impact-aware messaging;
- Maintenance/Data Engine can explain the underlying reason;
- no release claims `HEALTHY` when required decision evidence is not trustworthy.

---

# 26. Permanent Success Criterion

The long-term goal is **not zero degradation events**. Real providers, networks, entitlements and markets fail.

The 10/10 standard is:

> **DE.PULSE should rarely degrade because of its own architecture; when external or unavoidable degradation occurs, it should isolate the smallest truthful blast radius, preserve the most valuable current intelligence, explain the cause/impact clearly, recover safely, and learn from the event without silently changing production behavior.**
