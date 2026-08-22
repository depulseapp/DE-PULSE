# DE.PULSE — Current Adaptive Delivery Process

**Certified Stable:** `v18.9.0-stable`  
**Master corrective program:** #65 / `ADAPT-PROVIDER-INTELLIGENCE-010`  
**Hosted architecture program:** #66 / `ADAPT-HOSTED-SYNC-001`  
**Immediate product scope:** #64 / `ADAPT-RUNTIME-CRASH-001` -> `v18.9.1`  
**Governance alignment:** draft PR #67.

## 1. Delivery invariant

A DE.PULSE release is not delivered merely because code compiles, CI is green, one platform works or governance is documented.

Delivery follows:

`Governed -> Implemented -> Enforced -> Evidenced -> Packaged/Deployed -> Cross-Platform Verified -> Delivered -> Learned`

Every patch proves frozen G1 scope, explicit non-goals, canonical-owner preservation, deterministic regression evidence, truthful degraded/UNKNOWN behavior, required runtime evidence and issue/handoff truth.

## 2. Cross-Platform Definition of Delivered

For each capability, G1 freezes `macOS / Windows / Web = REQUIRED or justified N/A`.

A **shared capability is Delivered only when**:
1. canonical domain/API/state implementation is complete;
2. every REQUIRED client implementation is complete;
3. all REQUIRED clients pass equivalent authorization, state, data, freshness/provenance and user-meaning tests;
4. required macOS/Windows artifacts and Web deployment/runtime evidence pass;
5. no material parity debt remains.

There is no normal delivery state called “Mac done, Windows/Web later.” A platform may be used internally for technical validation, but that does not satisfy delivery or GA.

Platform-specific releases remain allowed when the responsibility itself is platform-specific. #64 / `v18.9.1` is therefore valid as macOS-only crash corrective work.

Temporary platform exceptions require a proven external blocker, explicit waiver/expiry, no misleading Delivered/GA status and a named recovery release. They are emergency exceptions, not roadmap structure.

## 3. Release-train delivery model

Major versions are strategic maturity generations; minor bands are coherent dependency phases; patches remain one-primary-responsibility units. For a shared capability, its required platform adapters are part of that one responsibility.

No next shared domain begins while current required-platform parity debt is material. Planned version labels remain reservations until G1; shipped releases are immutable.

## 4. v18.9.x delivery train

1. `v18.9.1` runtime crash corrective — macOS-specific escaped defect.
2. `v18.9.2` TradeInsight Settings/API-key UX.
3. `v18.9.3` coverage-aware router.
4. `v18.9.4` canonical company/instrument identity.
5. `v18.9.5` Market Data Modes/diagnostics.
6. `v18.9.6` provider observability/Adaptive telemetry.
7. `v18.9.7` Form 4 SHADOW enrichment.
8. `v18.9.8` symbol/company search.
9. `v18.9.9` movers/ranking SHADOW evidence.
10. `v18.9.10` remaining useful capability admission.
11. `v18.9.11` Session-Aware Data Readiness Maintenance.
12. `v18.9.12` professional closure.

Where v18.9 shared functionality exists on both supported native clients, delivery preserves semantic equivalence. Web becomes a REQUIRED shared client when the v19 cross-platform client foundation is delivered.

## 5. Protected Tier-0 delivery contract

Pre-market, regular market and after-hours are protected Tier-0 decision-support sessions. Live/current work outranks background maintenance/sync; provider/runtime/DB/worker reserve is explicit; background work is bounded, preemptible and cannot flood protected sessions.

## 6. v19 delivery train — foundations first, then lockstep capabilities

### v19.0.x — Governance / Control Plane / Data Foundation

1. `v19.0.0` Provider Capability / Legal Rights Registry.
2. `v19.0.1` Hosted Tenant/Identity/Device/Session Control Plane.
3. `v19.0.2` DE.PULSE Product Entitlement / Metering Policy.
4. `v19.0.3` Account Data Governance / Privacy Lifecycle.
5. `v19.0.4` Hosted Environment / IaC / Service Trust Foundation.
6. `v19.0.5` PostgreSQL Tenancy / Schema / Pool / HA-PITR Foundation.
7. `v19.0.6` Managed Secrets / KMS Lifecycle.
8. `v19.0.7` Software Supply-Chain / Artifact & Dependency Assurance.
9. `v19.0.8` Provider Quality / Cost / Coverage / SLO Scorecards.
10. `v19.0.9` Data Reconciliation / Revision / Point-in-Time Quality.

Hosted shared capability delivery is blocked until applicable foundations are executable and evidenced.

### v19.1.x — Hosted Data Plane + Cross-Platform Account/Sync

1. **v19.1.0 Hosted Provider Gateway** — versioned zero-key server boundary around Smart Provider Router v2; shared canonical backend for all clients.
2. **v19.1.1 Unified Serving Policy + Live Fan-Out** — same tenant/RBAC/product-entitlement/provider-right/privacy decisions for REST/WebSocket/SSE/native/web consumers.
3. **v19.1.2 Sync Protocol Foundation** — one versioned state protocol; native SQLite/outbox behavior and Web hosted mutation behavior converge on PostgreSQL authority.
4. **v19.1.3 Cross-Platform Account/Session Client Foundation** — **Mac + Windows + Web in the same release responsibility**.
5. **v19.1.4 Cross-Platform Preferences + Watchlists** — **Mac + Windows + Web in the same release responsibility**.
6. **v19.1.5 Cross-Platform Desks/Workspaces** — **Mac + Windows + Web in the same release responsibility**.

No Mac product pilot exists. Phase exit requires all required clients for these capabilities to pass.

### v19.2.x — Cross-Platform Shared Product + Assurance

1. **v19.2.0 Cross-Platform Research / Durable State** — Mac + Windows + Web.
2. **v19.2.1 Cross-Platform Market Intelligence / Discovery / Market Modes** — Mac + Windows + Web; same canonical intelligence and explanation/source/freshness meaning.
3. **v19.2.2 Cross-Platform Settings / RBAC / Product-Entitlement UX** — Mac + Windows + Web; same role/capability truth and direct-route/API denial.
4. **v19.2.3 Tenant-Aware Metering / Cost / Usage Observability** — platform dimension is observable but does not become separate business truth.
5. **v19.2.4 Multi-User Security / Abuse / Capacity Hardening** — mixed-client concurrent load, fairness, rate/circuit/load-shedding and protected-session capacity.
6. **v19.2.5 #66 Cross-Platform Assurance Closure** — no feature scope; explicit equivalence + adversarial/failure/recovery matrix.

#66 is complete only after actual Mac + Windows + Web clients/services prove the architecture. Documentation or one-client success cannot close it.

### Required #66 cross-platform evidence

At minimum:
- same account/session/role/product-entitlement state across clients;
- same authorized watchlists/preferences/desks/workspaces/research state;
- same Market Intelligence/Discovery/Market Mode meaning for equivalent evidence;
- same provider/freshness/provenance labels and degraded/UNKNOWN behavior;
- same revocation/downgrade/denial outcomes;
- native offline/reconnect convergence without Web becoming a second authority;
- mixed Mac/Windows/Web concurrent edit/conflict/delete/tombstone behavior;
- cross-account and direct-route/API/stream denial;
- provider-right/product-plan downgrade behavior;
- account export/deactivation/deletion/retention behavior;
- secret rotation and DB failover/PITR/restore;
- mixed-version protocol behavior;
- environment/config drift and artifact/dependency provenance failure behavior;
- noisy-neighbor/fairness/backpressure/protected-session pressure;
- no provider secret or unnecessary user data leakage.

### v19.3.x — Point-in-Time Evidence

- `v19.3.0` Institutional/13F infrastructure.
- `v19.3.1` Two-sided Long/Short evidence substrate.
- `v19.3.2` AODR candidate/ranking/outcome lineage.

Any user-facing capability surfaced from these substrates follows the Cross-Platform Definition of Delivered.

### v19.4.x — Reliability / Economics / v20 Readiness

- `v19.4.0` ADR-GDI professional reliability/capacity/runbooks.
- `v19.4.1` Specialized/paid-provider gap evaluation.
- `v19.4.2` v20 research-readiness audit.

### v19.5.0 — Major Closure

No feature scope. Require #66 PASS, zero material shared-capability parity debt, tenant/product-entitlement/provider-right/privacy separation, data lifecycle, environment/IaC, supply-chain provenance, API compatibility, SLO/capacity, recovery/rollback and actual Mac/Windows/Web runtime/deployment evidence where applicable.

## 7. Delivery acceptance for shared capabilities

In addition to normal G0-G16 evidence:
- one canonical capability contract;
- required-platform matrix frozen at G1;
- platform adapters implemented in-scope, not deferred;
- cross-platform state/data/authorization/UX meaning equivalence;
- cross-platform conflict/reconnect/revocation behavior where stateful;
- actual supported artifact/deployment proof;
- no separate provider/router/freshness/persistence/session/intelligence owner per client;
- no GA until all REQUIRED platforms pass.

Responsive layout and native interaction patterns may differ. Product truth and capability meaning may not.

## 8. v20 delivery train — Adaptive Intelligence

v20 begins only after `v19.5.0` PASS. Model/prompt governance precedes broad adaptive rollout.

- `v20.0.x` adaptive control/governance/analogues/calibration.
- `v20.1.x` ASBI.
- `v20.2.x` adaptive Institutional + TDTI.
- `v20.3.x` AODR.
- `v20.4.0` ADR-GDI adaptive operations.
- `v20.5.0` Professional Closure.

Every shared adaptive user-facing capability follows Mac + Windows + Web lockstep and `SHADOW -> VALIDATED -> APPROVED -> PRODUCTION`. No Execution remains permanent.

## 9. Per-patch G15 promotion guidance

For shared capabilities, G15 promotion is capability-wide, not platform-by-platform. Canary/bounded cohorts may be used, but a Mac-only GA followed by Windows/Web catch-up is prohibited. Rollback/kill criteria must preserve a coherent supported-platform experience.

## 10. Major-version handoff contract

- **v18.9 -> v19:** trustworthy native runtime/data foundation.
- **v19 -> v20:** one canonical hosted/account/data/intelligence product with Mac/Windows/Web parity, point-in-time evidence, reliability/capacity and outcome lineage trustworthy enough for governed learning.

Major closure cannot be bypassed by starting the next major early.

## 11. Permanent boundaries

U.S. Equities Processing; GLD/SLV/USO actionable exceptions; No Execution; Smart Provider Router v2 sole routing authority; canonical freshness/recovery sole freshness owner; existing multi-feed allocator sole subscription owner; BroadSnapshotBroker canonical reuse owner; canonical persistence/cache/state/session/calendar/identity owners; direct SEC/EDGAR authoritative; G0-G16 only.

## Exactly one next action

Obtain/reproduce the #64 macOS crash and freeze the narrow `v18.9.1` G1. Do not create `v18.9.2` or v19 product implementation branches until ordering permits it.