# DE.PULSE v18.7.0 — Runtime Reliability & Data Truth Audit

**Engineering branch:** `v18.7.0-development`  
**Engineering baseline:** `f1a9e0d0d76d4be565ac9355a09f77a33a3338a7`  
**Certified Stable baseline:** `v18.6.1-stable`  
**Scope authority:** `v18_7_0_scope.json` + `v18_7_0_g0_g3_contract.json`  
**Audit principle:** reuse current canonical owners; implement only verified remaining gaps.

## Executive result

v18.7.0 does **not** need a new retry engine, provider-health engine, freshness engine, provider router, reconciliation engine, queue manager or duplicate-work broker. The source audit found that those controls already exist and are materially integrated into the canonical `Engine.Snapshot()` path.

Verified gaps were narrower and user-trust oriented:

1. degradation reasons needed a canonical machine-readable taxonomy while preserving existing concise display labels;
2. degradation needed explicit affected-consumer/blast-radius truth;
3. required decision evidence needed fail-closed `UNKNOWN` + `ABSTAIN` semantics when no narrower cause was proven;
4. recovery needed hysteresis so one transient healthy sample could not clear a real degraded state;
5. a full provider queue needed immediate `QUEUE_SATURATED` truth instead of waiting for queue-age pressure;
6. Release G11-G16 still used mutable third-party Action tags and G12 installed Playwright outside the locked browser requirements contract;
7. realistic active-market reliability proof needed to combine production coalescing, bounded queues and degradation/SLO behavior.

Those are the v18.7.0 implementation targets. Existing controls remain the owners.

## Source-owner disposition

| v18.7.0 concern | Canonical owner | Audit disposition | v18.7 action |
|---|---|---|---|
| `DATA DEGRADED` semantics | `runtime_degradation.go` + `Engine.Snapshot()` | **VERIFIED GAP** | Canonical reason codes, pressure state, affected consumers and decision impact added. |
| Critical decision truth / abstention | freshness + `runtime_degradation.go` | **VERIFIED GAP** | Required live decision evidence now fails closed to `UNKNOWN` / `ABSTAIN` when no narrower cause is proven. |
| Provider/capability health aggregation | `provider_capabilities.go`, `provider_router.go`, `smart_router_v2.go` | **INHERITED / PROVEN** | Preserve. Human-readable registry and executable Router health remain distinct layers. No new aggregator. |
| Dataset/session freshness | `data_freshness.go` | **INHERITED / PROVEN** | Preserve provider/receipt timestamps, session-aware limits, delayed-feed truth and consumer scope. |
| Runtime freshness/SLO | `runtime_slo.go` | **INHERITED + HARDENED** | Preserve existing SLOs; add explicit recovery-pending observations and hysteresis. |
| Duplicate work / single-flight | `broad_snapshot_broker.go` | **INHERITED / PROVEN** | Preserve current broker; add higher-fanout active-market proof. |
| Canonical reuse | snapshot broker + runtime reuse telemetry | **INHERITED / PROVEN** | Preserve reuse; continue measuring provider calls avoided / hit rate. |
| Bounded retries / fallback | `executeProviderRoute`, provider loaders | **INHERITED / PROVEN** | Preserve finite Router chain, bounded provider-specific attempts and cooldowns. No retry loop added. |
| Provider/capability circuits | `provider_router.go`, `smart_router_v2.go` | **INHERITED / PROVEN** | Preserve cooldown / revalidation / half-open probe behavior. Runtime hysteresis protects user-facing recovery truth. |
| Request budget / rate limit | `runtime_observability.go` + Router | **INHERITED / PROVEN** | Preserve local pacing, budget pressure and low-tier shedding. |
| Backpressure / load shedding | `workload_controller.go` | **INHERITED + HARDENED** | Preserve tier/reserved-critical model; recognize hard provider-queue saturation immediately. |
| Source disagreement | `provider_reconciliation.go` | **INHERITED / PROVEN** | Preserve timestamp-valid contemporaneous comparison, Router-priority winner and explicit `CONFLICT`; never average silently. |
| Recovery hysteresis | `RuntimeSLOTracker` + `Engine.Snapshot()` | **VERIFIED GAP** | Require 3 consecutive healthy observations and >=5s stability; relapse resets recovery streak. |
| Active-market reliability evidence | production broker/workload/degradation/SLO owners | **VERIFIED EVIDENCE GAP** | Add deterministic synthetic burst/pressure integration tests; actual CI result required before claiming PASS. |
| Release Action/browser reproducibility | `release.yml` + CI dependency lock/gate | **DEFERRED PROCESS GAP** | Immutable SHA pins + locked Playwright requirements/cache contract added and gate-enforced. |

## Provider/capability health ownership

The audit specifically rejected adding another health service.

`buildProviderCapabilityRegistry(...)` provides human-readable capability rows with provider, capability, entitlement/availability status, detail and the consuming DE.PULSE surfaces.

`buildProviderRouterSnapshot(...)` remains the executable truth owner. For each dataset/provider it combines:

- configured/not-configured state;
- persisted capability entitlement/state;
- capability circuit and provider-wide circuit state;
- request success/error/latency/budget telemetry;
- preferred versus serving provider;
- route score and score reasons;
- fallback reason;
- rate-limit state;
- last success/failure;
- recovery state;
- expected delay, cost class and data-rights metadata.

This is sufficient health aggregation for v18.7.0. A parallel aggregator would create conflicting truth.

## Freshness and decision truth

The current freshness owner already distinguishes check/receipt/provider/data timestamps, session-specific cadence and provider-specific delay semantics. v18.7.0 therefore does not replace freshness calculation.

The verified gap was downstream fail-closed behavior: a live runtime could have unusable required decision evidence without a narrower network/provider/load cause and therefore risk falling through without a degradation code. v18.7.0 now emits:

- display: `DATA DEGRADED`;
- canonical reason: `UNKNOWN`;
- `CriticalUsable=false`;
- `Abstain=true`;
- explicit affected datasets/consumers.

Absence of diagnosis is no longer treated as evidence of health.

## Degradation taxonomy

Canonical v18.7 reason codes currently used by runtime degradation:

- `QUEUE_SATURATED`
- `LOCAL_OVERLOAD`
- `RATE_LIMITED`
- `NETWORK_FAILURE`
- `PROVIDER_DOWN`
- `LOW_COVERAGE`
- `UNKNOWN`

Existing concise `Code` values remain for renderer/API compatibility. Canonical `ReasonCode` is the durable diagnosis field.

Pressure states:

- `HEALTHY` — no active decision-relevant degradation;
- `PROTECTED` — degradation exists but critical decision evidence remains usable and the blast radius is isolated;
- `DEGRADED` — required decision evidence is not trustworthy enough; affected conclusions must abstain;
- `RECOVERING` — underlying samples are healthy but recovery hysteresis has not yet been satisfied.

## Recovery rule

A real degradation appears immediately.

Recovery requires:

1. three consecutive healthy observations; **and**
2. at least five seconds of continuous healthy stability.

Any relapse resets the healthy streak. Until both conditions are satisfied, the prior degradation remains exposed as `RECOVERING`. This prevents a single transient provider success from falsely declaring the runtime healthy.

## Load / backpressure truth

The existing workload controller remains authoritative:

- market-critical, user-actionable, promoted-radar, broad-discovery and background tiers;
- bounded class concurrency;
- bounded queues;
- reserved critical capacity;
- lower-priority shedding before protected work.

v18.7 additionally treats `Queued >= MaxQueue` as immediate `QUEUE_SATURATED` truth. This is intentionally independent from queue age: a hard-full queue is already capacity pressure even at the instant it becomes full.

## Duplicate-work truth

The existing broad snapshot broker already:

- reuses fresh canonical observations;
- fetches only missing symbols;
- canonicalizes symbol sets;
- coalesces overlapping in-flight requests;
- bounds cache entries;
- expires stale observations;
- exposes provider-fetch/coalescing diagnostics.

v18.7 adds a synthetic active-market burst test using 16 concurrent equivalent requests and production broker code. It must prove a single provider fetch plus observable coalesced waiters in CI before this evidence is marked PASS.

## Release workflow reproducibility closure

The deferred Release hardening is part of this release because `release.yml` participates directly in the v18.7 Stable publication path.

v18.7 now requires:

- immutable 40-hex SHAs for all external Actions in Fast, Qualified and Release;
- the exact pinned Action set in `tools/ci/ci_dependency_lock.json`;
- Release G12 browser setup through `tools/ci/browser-requirements.txt` instead of an unpinned direct Playwright install;
- the same safe pip cache dependency contract used by Qualified;
- `tools/ci/reproducibility_gate.py` to reject drift in any of these three canonical workflows.

This is pin/reproducibility hardening only. It does not add a fourth workflow or change the G11-G16 release architecture.

## Qualification requirements before v18.7.0 can be called complete

The following are **required evidence, not assumed PASS**:

1. exact-head Fast;
2. Go format/vet/full suite including v17 degradation/backpressure and v18.7 regressions;
3. active-market synthetic coalescing/pressure proof;
4. full Qualified backend;
5. race detector;
6. randomized test order;
7. deterministic equivalence/market-truth gates;
8. renderer regressions affected by public runtime payload changes;
9. Chrome broad behavior;
10. WebKit primary compatibility;
11. full G10 reconciliation against the conserved authority ledger;
12. immutable merged candidate;
13. one G11-G16 release run with macOS Apple Silicon + Windows x64 packaged runtime evidence and same-run publication.

No item above is considered passed merely because its test or workflow exists.

## Architectural conclusion

The v18.7 implementation remains a **hardening slice of existing owners**, not an architecture expansion. Smart Provider Router v2 remains sole routing authority; canonical freshness remains sole freshness authority; RuntimeSLOTracker remains recovery/SLO owner; WorkloadController remains backpressure owner; BroadSnapshotBroker remains duplicate snapshot suppression owner; Provider Reconciliation remains disagreement owner; `Engine.Snapshot()` remains the canonical aggregation boundary.

This preserves DE.PULSE's Adaptive Intelligence direction while keeping deterministic market truth auditable and fail-closed.
