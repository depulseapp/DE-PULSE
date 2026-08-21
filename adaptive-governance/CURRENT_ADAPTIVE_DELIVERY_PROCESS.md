# DE.PULSE — Current Adaptive Delivery Process

**Certified Stable:** `v18.9.0-stable`  
**Master corrective program:** #65 / `ADAPT-PROVIDER-INTELLIGENCE-010`  
**Immediate product scope:** #64 / `ADAPT-RUNTIME-CRASH-001` -> `v18.9.1`  
**Governance alignment:** draft PR #67.

## 1. Delivery invariant

A DE.PULSE release is not delivered merely because code compiles, CI is green or governance is documented.

Delivery follows:

`Governed -> Implemented -> Enforced -> Evidenced -> Packaged/Deployed -> Verified -> Delivered -> Learned`

Every patch must prove:
- frozen G1 scope fully implemented;
- explicit non-goals untouched;
- canonical owners preserved;
- deterministic regression evidence for changed responsibility;
- truthful partial/degraded/UNKNOWN/ABSTAIN behavior;
- no known implementation miss left without a durable target;
- required native/hosted runtime evidence;
- issue/handoff/checkpoint truth agrees with executable evidence.

## 2. Release-train delivery model

Major versions are strategic maturity stages; minor bands are coherent dependency phases; patch releases remain small independent units.

A minor band does **not** authorize one large bundle. Each patch inside the band passes its own G0–G16 lifecycle and can be rolled back/held independently.

Planned future version labels are reservations until G1. G0/G1 may split broad work and shift unstarted reservations. Shipped versions are immutable.

## 3. v18.9.x delivery train

1. **v18.9.1 — Runtime crash corrective** — packaged macOS reproduction/root cause/fix, preserved user state/API keys, lifecycle/relaunch regression.
2. **v18.9.2 — TradeInsight Settings/API-key UX** — canonical secret owner, masked controls, truthful status, scroll/focus preservation.
3. **v18.9.3 — Coverage-aware router** — DB/cache reuse first, residual-gap acquisition, provenance/coverage proof.
4. **v18.9.4 — Canonical identity** — shared company/instrument identity and all-desk presentation.
5. **v18.9.5 — Market Data Modes/diagnostics** — behavior-oriented modes and capability freshness/coverage/source truth.
6. **v18.9.6 — Provider observability/Adaptive telemetry** — moved before provider expansion; must produce measurable SHADOW usefulness, provider/call-avoidance/runtime/headroom evidence.
7. **v18.9.7 — Form 4 enrichment** — SHADOW-first, SEC authoritative, measured with v18.9.6 telemetry.
8. **v18.9.8 — Symbol/company search** — canonical fallback/corroboration.
9. **v18.9.9 — Movers/ranking evidence** — SHADOW candidate evidence through Opportunity Radar.
10. **v18.9.10 — Remaining useful capability admission** — explicit disposition/consumer/rights/freshness/rate/retention/lifecycle for every useful entitlement.
11. **v18.9.11 — Session-Aware Data Readiness Maintenance** — overnight/weekend bounded maintenance with protected-session priority/preemption/checkpoints.
12. **v18.9.12 — Professional closure audit** — no new feature scope; zero-miss/duplicate-owner/provider/runtime/package closure.

No later patch begins until current patch handoff identifies exactly one next action and any discovered miss is durably dispositioned.

## 4. Protected Tier-0 delivery contract

Pre-market, regular market and after-hours are protected Tier-0 decision-support sessions.

Release evidence for maintenance, hosted gateway, synchronization, DB/pool or adaptive workload changes must demonstrate:
- explicit provider/runtime/DB/worker reserve for current-session capabilities;
- live/current quote/liquidity/VIX/news/catalyst/SEC/Opportunity Radar/Research/readiness outrank background work;
- maintenance/sync acquisition suspends/yields when not directly required by live consumers;
- background queues/concurrency are bounded and preemptible;
- no heavy compaction/reconciliation/backfill near or during protected sessions;
- recovery/catch-up occurs in eligible low-priority windows without request storms.

Machine contract: `adaptive-governance/PERSISTENCE_REUSE_AND_OFF_HOURS_DATA_READINESS_CONTRACT.md`.

## 5. v19 delivery train — Professional Data Infrastructure + Hosted Account Platform

v19 begins only after `v18.9.12` closure PASS.

### v19.0.x — Governance, Control Plane & Data Foundation

Delivery order is dependency-enforced:

1. **v19.0.0 Rights Registry** — machine-readable and runtime-consumable provider/capability rights/entitlements; unknown/expired disallowed behavior fails closed.
2. **v19.0.1 Hosted Identity/Device/Session** — canonical role/capability/session/device truth across API/SSE/native/web; revocation and privileged re-authentication evidenced.
3. **v19.0.2 PostgreSQL Foundation** — tenancy/schema/pool/index/HA/PITR/backup/restore/migration/RPO/RTO; no broad sync activation.
4. **v19.0.3 Managed Secrets/KMS** — environment isolation, rotation/rollback/compromise recovery/redaction; zero platform provider secret on commercial clients.
5. **v19.0.4 Provider SLO/Cost/Coverage Scorecards** — measured freshness/completeness/latency/reliability/rate/cost/usefulness/calls-avoided/error-budget evidence.
6. **v19.0.5 Data Reconciliation/Revision Quality** — independence/conflict/revisions/adjustments/point-in-time provenance.

**Phase exit gate:** hosted provider/sync activation is blocked until required rights, identity, secret and DB/recovery controls are executable and evidenced.

### v19.1.x — Zero-Key Provider Data Plane & Native Sync

1. **v19.1.0 Hosted Provider Gateway** — authenticated server boundary around existing Smart Provider Router v2; canonical freshness/cache/persistence reuse; bounded circuits/backpressure/kill switch.
2. **v19.1.1 Rights/Entitlement + Live Fan-Out** — runtime gating at router/cache/persistence/REST/WebSocket/SSE; existing multi-feed owner; no entitlement leakage.
3. **v19.1.2 Sync Protocol Foundation** — snapshot/high-watermark bootstrap, SQLite atomic outbox, idempotent push, authoritative server sequence/change log, incremental pull, checkpoints/tombstones/retention/compaction/re-bootstrap/mixed-version negotiation.
4. **v19.1.3 macOS Pilot** — preferences/watchlists, offline/restart/reconnect/conflicts/local account isolation/user switching/lost-device proof.
5. **v19.1.4 Desks/Workspaces** — versioned membership/config/delete/history convergence using same transport.

### v19.2.x — Cross-Platform Parity & #66 Assurance

1. **v19.2.0 Windows x64 parity** — same account/sync/security semantics.
2. **v19.2.1 Hosted web parity** — same API/session/capability/PostgreSQL state; browser cache non-authoritative.
3. **v19.2.2 Rights-aware research/state portability** — only lawful/provenance-bound durable artifacts.
4. **v19.2.3 Multi-user security/cost/abuse/capacity hardening** — account/user/device/capability usage attribution, quotas, throttling, edge limits, fairness, cache/call avoidance and protected-session load shedding.
5. **v19.2.4 #66 assurance closure** — no new feature scope; full adversarial/failure/recovery matrix and implementation-miss audit.

#66 / `ADAPT-HOSTED-SYNC-001` is complete only after actual supported clients/services prove the architecture. Documentation does not close it.

### Required #66 negative/adversarial evidence

At minimum:
- cross-account read/write denial;
- role/capability downgrade and revoked session/device behavior;
- stale/unsupported client protocol behavior without losing local outbox mutations;
- duplicate replay/idempotency;
- network loss during server/client apply;
- new-device and checkpoint-expired bootstrap/re-bootstrap;
- concurrent edit/delete/tombstone behavior;
- provider right/entitlement downgrade/expiry;
- secret rotation failure, rollback and compromise revoke;
- DB failover/PITR/restore/server rollback;
- long-lived stream reauthorization/revocation;
- queue/backpressure/provider outage vs DB/gateway outage distinction;
- protected-session load pressure and background-yield behavior;
- user switching/local SQLite account isolation;
- no provider secret in clients/logs/telemetry/sync payloads.

### v19.3.x — Point-in-Time Evidence

- **v19.3.0 Institutional/13F infrastructure** — direct SEC truth, manager/security identity, amendments/filing lag/point-in-time holdings/outcome lineage.
- **v19.3.1 Two-sided Long/Short evidence substrate** — point-in-time plan/thesis/outcome evidence and explicit UNKNOWN for missing lawful data.
- **v19.3.2 AODR candidate/ranking/outcome lineage** — point-in-time candidate/rank/reason transitions and surfaced-vs-missed outcomes.

### v19.4.x — Reliability / Economics / v20 Readiness

- **v19.4.0 ADR-GDI professional reliability/capacity** — SLO/degradation, warm start, indexes/pools, load shedding, bounds/reserves, maintenance/preemption economics, controlled failure/soak proof.
- **v19.4.1 Specialized/paid-provider gap evaluation** — only measured gaps justify provider change; same canonical router/rights/persistence/session contracts.
- **v19.4.2 v20 research-readiness audit** — no new model scope; prove point-in-time/provenance/rights/independence/leakage/reliability readiness.

### v19.5.0 — Major Closure

No new feature scope. Require:
- #66 closure PASS;
- provider/data-rights/commercial posture truthful and machine enforced;
- point-in-time/revision/provenance quality;
- SLO/error-budget/reliability/capacity evidence;
- DB/secret/provider-right recovery drills;
- no unowned dataset/provider role/parallel owner;
- actual supported packages/runtimes/services and provenance;
- zero unresolved P0 security/rights/recovery/compatibility issue.

## 6. Delivery acceptance for hosted/security/data patches

In addition to normal G0–G16 evidence, affected patches require industry-strength operational readiness:
- Architecture Decision Record/equivalent durable decision for material authority/trust-boundary changes;
- explicit threat model + data/tenant classification;
- API/schema/protocol compatibility contract;
- backward-compatible migration + rollback/roll-forward plan;
- SLO/error-budget + observability before broad activation;
- contract and negative authorization tests;
- load/soak/capacity/failure-injection/failover evidence where relevant;
- bounded feature flag/kill switch/circuit breaker/canary controls for risky hosted activation where appropriate;
- runbook/recovery steps sufficient for an operator other than the original developer;
- no secrets in source/client/sync/ordinary logs/traces/crash output;
- exact source/artifact/deployment provenance.

## 7. v20 delivery train — Adaptive Intelligence

v20 begins only after `v19.5.0` PASS.

### v20.0.x — Adaptive control/governance before broad model rollout
- **v20.0.0 Adaptive research control plane + immutable experiment ledger**;
- **v20.0.1 Model/prompt governance + Champion/Challenger** — deliberately moved before broad adaptive intelligence;
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
- `v20.4.0` ADR-GDI adaptive provider/recovery/workload/maintenance/reserve optimization under SHADOW/Champion-Challenger and explicit promotion.

### v20.5.0 — Professional Closure
No new feature scope. Require calibration/utility/drift/abstention, deterministic-boundary protection, privacy/security/data rights, reproducibility, rollback, actual supported artifacts, zero silent self-modification and No Execution.

## 8. Per-patch G15 delivery/promotion guidance

For local/native-only low-risk patches, normal certified artifact promotion is sufficient.

For hosted/server/security/data-plane changes where partial rollout is possible, G15 should prefer:
- disabled-by-default or controlled capability activation until evidence is complete;
- canary/bounded cohort or owner-only pilot where practical;
- explicit rollback/kill-switch criteria;
- health/SLO/error-budget monitoring during activation;
- no irreversible schema/protocol/client dependency before compatibility is proven.

Promotion mechanics must stay inside the canonical release process; this does not create another gate.

## 9. Major-version handoff contract

- **v18.9 -> v19:** runtime, coverage-aware acquisition, persistence reuse, identity, Market Modes, provider telemetry and session-aware maintenance are trustworthy and zero-gap enough for hosted/professional infrastructure.
- **v19 -> v20:** rights/provenance/point-in-time data, hosted account/sync/gateway security, reliability/capacity and outcome lineage are trustworthy enough for governed learning.
- A major closure cannot be bypassed by starting the next major early.

## 10. Permanent boundaries

U.S. Equities Processing; GLD/SLV/USO actionable exceptions; No Execution; Smart Provider Router v2 sole routing authority; canonical freshness/recovery sole freshness owner; existing multi-feed allocator sole subscription owner; BroadSnapshotBroker canonical reuse owner; canonical persistence/cache/state owners; canonical U.S. market calendar/session owner; direct SEC/EDGAR authoritative; canonical identity/role/capability truth; G0–G16 only.

## Exactly one next action

Obtain/reproduce the #64 macOS crash and freeze the narrow `v18.9.1` scope. Do not create `v18.9.2` or v19 product implementation branches until `v18.9.1` is truthfully closed or the crash is proven external/non-product.
