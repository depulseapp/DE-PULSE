# DE.PULSE — Current Adaptive Delivery Process

**Certified Stable:** `v18.9.0-stable`  
**Master corrective program:** #65 / `ADAPT-PROVIDER-INTELLIGENCE-010`  
**Hosted architecture program:** #66 / `ADAPT-HOSTED-SYNC-001`  
**Immediate product scope:** #64 / `ADAPT-RUNTIME-CRASH-001` -> `v18.9.1`  
**Governance alignment:** draft PR #67.

## 1. Delivery invariant

A DE.PULSE release is not delivered merely because code compiles, CI is green or governance is documented.

Delivery follows:

`Governed -> Implemented -> Enforced -> Evidenced -> Packaged/Deployed -> Verified -> Delivered -> Learned`

Every patch proves:
- frozen G1 scope fully implemented;
- explicit non-goals untouched;
- canonical owners preserved;
- deterministic regression evidence for changed responsibility;
- truthful partial/degraded/UNKNOWN/ABSTAIN behavior;
- no known implementation miss left without a durable target;
- required native/hosted runtime evidence;
- issue/handoff/checkpoint truth agrees with executable evidence.

## 2. Release-train delivery model

Major versions are strategic maturity generations; minor bands are coherent dependency phases; patch releases remain small independent units.

A minor band never authorizes one large bundle. Each patch inside the band passes its own G0-G16 lifecycle and can be held/rolled back independently.

Planned future version labels are reservations until G1. G0/G1 may split broad work and shift unstarted reservations. Shipped versions are immutable. Corrective/security work can preempt the planned train.

## 3. v18.9.x delivery train

1. **v18.9.1 — Runtime crash corrective** — packaged macOS reproduction/root cause/fix, preserved user state/API keys, warm-state/relaunch regression.
2. **v18.9.2 — TradeInsight Settings/API-key UX** — canonical secret owner, masked controls, truthful status, scroll/focus preservation.
3. **v18.9.3 — Coverage-aware router** — DB/cache reuse first, residual-gap acquisition, provenance/coverage proof.
4. **v18.9.4 — Canonical company/instrument identity** — shared identity/presentation.
5. **v18.9.5 — Market Data Modes/diagnostics** — behavior-oriented modes and capability freshness/coverage/source truth.
6. **v18.9.6 — Provider observability/Adaptive telemetry** — deliberately before provider expansion; measurable SHADOW usefulness, provider/call-avoidance/runtime/headroom evidence.
7. **v18.9.7 — Form 4 enrichment** — SHADOW-first, SEC authoritative, measured using v18.9.6 telemetry.
8. **v18.9.8 — Symbol/company search** — canonical fallback/corroboration.
9. **v18.9.9 — Movers/ranking evidence** — SHADOW candidate evidence through Opportunity Radar.
10. **v18.9.10 — Remaining useful capability admission** — explicit disposition/consumer/rights/freshness/rate/retention/lifecycle for each useful entitlement.
11. **v18.9.11 — Session-Aware Data Readiness Maintenance** — overnight/weekend bounded maintenance with protected-session priority/preemption/checkpoints.
12. **v18.9.12 — Professional closure audit** — no new feature scope; zero-miss/duplicate-owner/provider/runtime/package closure.

No later patch begins until the current patch handoff identifies exactly one next action and every discovered miss is durably dispositioned.

## 4. Protected Tier-0 delivery contract

Pre-market, regular market and after-hours are protected Tier-0 decision-support sessions.

Maintenance, hosted gateway, synchronization, DB/pool and adaptive workload releases must show:
- provider/runtime/DB/worker reserve for current-session capabilities;
- live/current decision support outranks background work;
- maintenance/sync acquisition suspends/yields when not required by live consumers;
- background queues/concurrency are bounded and preemptible;
- no heavy compaction/reconciliation/backfill during protected sessions;
- recovery/catch-up occurs in eligible lower-priority windows without request storms.

Machine contract: `adaptive-governance/PERSISTENCE_REUSE_AND_OFF_HOURS_DATA_READINESS_CONTRACT.md`.

## 5. v19 delivery train — Professional Data Infrastructure + Hosted Account Platform

v19 begins only after `v18.9.12` PASS.

### v19.0.x — Governance, Control Plane & Data Foundation

Delivery is dependency-enforced:
1. **v19.0.0 Provider Capability / Legal Rights Registry** — runtime-consumable provider/capability rights; unknown/expired disallowed behavior fails closed.
2. **v19.0.1 Hosted Tenant/Identity/Device/Session** — first-class tenant context plus canonical role/capability/session/device truth across API/SSE/native/web; revocation and privileged re-authentication evidenced.
3. **v19.0.2 DE.PULSE Product Entitlement / Metering Policy** — billing-provider-agnostic plan/status/feature/quota/grace/suspension policy. Separate from RBAC and upstream provider rights.
4. **v19.0.3 PostgreSQL Foundation** — tenancy/schema/pool/index/HA/PITR/backup/restore/migration/RPO/RTO; no broad sync activation.
5. **v19.0.4 Managed Secrets/KMS** — environment isolation, rotation/rollback/compromise recovery/redaction; zero platform provider secret on commercial clients.
6. **v19.0.5 Provider SLO/Cost/Coverage Scorecards** — measured freshness/completeness/latency/reliability/rate/cost/usefulness/calls-avoided/error-budget and tenant-aware health evidence.
7. **v19.0.6 Data Reconciliation/Revision Quality** — independence/conflict/revisions/adjustments/point-in-time provenance.

**Phase exit:** hosted provider/sync activation is blocked until tenant identity, RBAC, product entitlement, provider rights, secret and DB/recovery controls are executable and evidenced.

### v19.1.x — Zero-Key Provider Data Plane & Native Sync

1. **v19.1.0 Hosted Provider Gateway** — authenticated/versioned server boundary around existing Smart Provider Router v2; canonical freshness/cache/persistence reuse; bounded circuits/backpressure/kill switch; API inventory/version/deprecation ownership.
2. **v19.1.1 Unified Serving Policy + Live Fan-Out** — tenant/RBAC/product-entitlement/provider-right gating at router/cache/persistence/REST/WebSocket/SSE; existing multi-feed owner; no cross-tenant/tier/right leakage.
3. **v19.1.2 Sync Protocol Foundation** — snapshot/high-watermark bootstrap, SQLite atomic outbox, idempotent push, authoritative server sequence/change log, incremental pull, checkpoints/tombstones/retention/compaction/re-bootstrap/mixed-version negotiation.
4. **v19.1.3 macOS Pilot** — preferences/watchlists, offline/restart/reconnect/conflicts/local-account isolation/user switching/lost-device proof.
5. **v19.1.4 Desks/Workspaces** — versioned membership/config/delete/history convergence using the same transport.

### v19.2.x — Cross-Platform Parity & #66 Assurance

1. **v19.2.0 Windows x64 parity** — same account/sync/security semantics.
2. **v19.2.1 Hosted web parity** — same API/session/capability/product-entitlement/PostgreSQL state; browser cache non-authoritative.
3. **v19.2.2 Rights-aware research/state portability** — only lawful/provenance-bound durable artifacts.
4. **v19.2.3 Tenant-Aware Metering / Cost / Usage Observability** — tenant/account/user/device/capability attribution, plan/quota consumption, cache/call avoidance, streaming usage, provider cost where known and tenant-health signals.
5. **v19.2.4 Multi-User Security / Abuse / Capacity Hardening** — object/function authorization negatives, sensitive-flow abuse protection, rate limits, fairness/noisy-neighbor isolation, edge limits, circuits/load shedding and protected-session capacity.
6. **v19.2.5 #66 assurance closure** — no feature scope; full adversarial/failure/recovery matrix and implementation-miss audit.

#66 is complete only after actual supported clients/services prove the architecture. Documentation alone does not close it.

### Required #66 negative/adversarial evidence

At minimum:
- cross-account object/API/stream access denial;
- role/capability downgrade;
- product-plan downgrade/suspension/quota exhaustion;
- revoked session/device behavior;
- stale/unsupported client protocol behavior without losing local outbox mutations;
- duplicate replay/idempotency;
- network loss during server/client apply;
- new-device and checkpoint-expired bootstrap/re-bootstrap;
- concurrent edit/delete/tombstone behavior;
- provider-right downgrade/expiry;
- secret rotation failure, rollback and compromise revoke;
- DB failover/PITR/restore/migration rollback;
- long-lived stream reauthorization/revocation;
- API inventory/version/deprecation compatibility;
- provider outage vs DB/gateway outage distinction;
- queue/backpressure/circuit/load shedding;
- noisy-neighbor/fairness pressure;
- protected-session load pressure/background-yield behavior;
- user switching/local SQLite account isolation;
- no provider secret in clients/logs/telemetry/sync payloads.

### v19.3.x — Point-in-Time Evidence

- **v19.3.0 Institutional/13F infrastructure** — direct SEC truth, manager/security identity, amendments/filing lag/point-in-time holdings/outcome lineage.
- **v19.3.1 Two-sided Long/Short evidence substrate** — point-in-time plan/thesis/outcome evidence and explicit UNKNOWN for missing lawful data.
- **v19.3.2 AODR candidate/ranking/outcome lineage** — point-in-time candidate/rank/reason transitions and surfaced-vs-missed outcomes.

### v19.4.x — Reliability / Economics / v20 Readiness

- **v19.4.0 ADR-GDI professional reliability/capacity** — SLO/error-budget/degradation, warm start, indexes/pools, load shedding, bounds/reserves, maintenance/preemption economics, controlled failure/soak proof and hosted operational runbooks.
- **v19.4.1 Specialized/paid-provider gap evaluation** — only measured gaps justify provider change; same canonical router/rights/persistence/session contracts.
- **v19.4.2 v20 research-readiness audit** — no new model scope; point-in-time/provenance/rights/independence/leakage/reliability readiness.

### v19.5.0 — Major Closure

No feature scope. Require:
- #66 PASS;
- tenant isolation and negative cross-tenant evidence;
- canonical RBAC/session/device revocation;
- DE.PULSE product entitlement separated from provider legal/data rights;
- provider/data-rights/commercial posture truthful and machine-enforced;
- API inventory/version/deprecation truth;
- point-in-time/revision/provenance quality;
- SLO/error-budget/reliability/capacity evidence;
- DB/secret/provider-right recovery drills;
- tenant-aware metering/cost/usage/health evidence;
- multi-user abuse/noisy-neighbor/capacity proof;
- no unowned dataset/provider role/parallel owner;
- actual supported packages/runtimes/services and provenance;
- zero unresolved P0 security/rights/product-entitlement/recovery/compatibility issue.

## 6. Delivery acceptance for hosted/security/data patches

In addition to normal G0-G16 evidence, affected patches require:
- ADR/equivalent durable decision for material authority/trust-boundary changes;
- threat model + data/tenant classification;
- API/schema/protocol compatibility and inventory/deprecation contract;
- backward-compatible migration + rollback/roll-forward plan;
- SLO/error-budget + tenant-aware observability before broad activation;
- contract and negative authorization tests;
- load/soak/capacity/fairness/failure-injection/failover evidence where relevant;
- bounded feature flag/kill switch/circuit breaker/canary controls for risky hosted activation where appropriate;
- operator runbook/recovery steps not dependent on the original developer;
- no secrets in source/client/sync/ordinary logs/traces/crash output;
- exact source/artifact/deployment provenance.

## 7. v20 delivery train — Adaptive Intelligence

v20 begins only after `v19.5.0` PASS.

### v20.0.x — Adaptive control/governance before broad rollout
- **v20.0.0 Adaptive research control plane + immutable experiment ledger**;
- **v20.0.1 Model/prompt governance + Champion/Challenger**;
- **v20.0.2 Historical analogues/regime outcomes**;
- **v20.0.3 Calibration/FP-FN/miss/contradiction/drift**.

### v20.1.x — ASBI
- `v20.1.0` behavioral fingerprints/state transitions;
- `v20.1.1` scenarios/probability momentum/calibration.

### v20.2.x — Institutional + TDTI
- `v20.2.0` adaptive 13F;
- `v20.2.1` competing Long/Short/No Reliable Edge;
- `v20.2.2` two-sided trade-plan/readiness/outcome validation; No Execution.

### v20.3.x — AODR
- `v20.3.0` adaptive shared ranking;
- `v20.3.1` diversity/opportunity cost/personalized relevance after shared truth.

### v20.4.x — Adaptive operations
- `v20.4.0` ADR-GDI adaptive provider/recovery/workload/maintenance/reserve optimization under SHADOW/Champion-Challenger until explicit promotion.

### v20.5.0 — Professional Closure
No feature scope. Require calibration/utility/drift/abstention, deterministic-boundary protection, privacy/security/data rights, reproducibility, rollback, actual supported artifacts, zero silent self-modification and No Execution.

Adaptive production influence remains `SHADOW -> VALIDATED -> APPROVED -> PRODUCTION`.

## 8. Per-patch G15 delivery/promotion guidance

For local/native-only low-risk patches, normal certified artifact promotion is sufficient.

For hosted/server/security/data-plane changes where partial rollout is possible, G15 prefers:
- disabled-by-default or controlled capability activation until evidence is complete;
- canary/bounded cohort or owner-only pilot where practical;
- explicit rollback/kill-switch criteria;
- health/SLO/error-budget/tenant-health monitoring during activation;
- no irreversible schema/protocol/client dependency before compatibility is proven.

This remains inside G15; no additional gate is created.

## 9. Major-version handoff contract

- **v18.9 -> v19:** runtime, coverage-aware acquisition, persistence reuse, identity, Market Modes, provider telemetry and session-aware maintenance are trustworthy enough for hosted/professional infrastructure.
- **v19 -> v20:** tenant/RBAC/product-entitlement/provider-right boundaries, point-in-time/provenance data, hosted account/sync/gateway security, reliability/capacity and outcome lineage are trustworthy enough for governed learning.
- Major closure cannot be bypassed by starting the next major early.

## 10. Permanent boundaries

U.S. Equities Processing; GLD/SLV/USO actionable exceptions; No Execution; Smart Provider Router v2 sole routing authority; canonical freshness/recovery sole freshness owner; existing multi-feed allocator sole subscription owner; BroadSnapshotBroker canonical reuse owner; canonical persistence/cache/state owners; canonical U.S. market calendar/session owner; direct SEC/EDGAR authoritative; canonical tenant/identity/role/capability truth; G0-G16 only.

## Exactly one next action

Obtain/reproduce the #64 macOS crash and freeze the narrow `v18.9.1` scope. Do not create `v18.9.2` or v19 product implementation branches until `v18.9.1` is truthfully closed or the crash is proven external/non-product.