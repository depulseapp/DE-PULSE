# DE.PULSE — Current Adaptive Delivery Process

**Certified Stable:** `v18.9.0-stable`  
**Master corrective program:** #65 / `ADAPT-PROVIDER-INTELLIGENCE-010`  
**Hosted architecture program:** #66 / `ADAPT-HOSTED-SYNC-001`  
**Immediate product scope:** #64 / `ADAPT-RUNTIME-CRASH-001` -> `v18.9.1`  
**Governance alignment:** draft PR #67.

## 1. Delivery invariant

A DE.PULSE release is not delivered because code compiles, CI is green, one platform works, or governance exists.

`Governed -> Implemented -> Enforced -> Evidenced -> Packaged/Deployed -> Cross-Platform Verified -> Delivered -> Learned`

## 2. Cross-Platform Definition of Delivered

For every G1 capability, freeze `macOS / Windows / Web = REQUIRED or justified N/A`.

A shared capability is Delivered only when:
1. canonical domain/API/state implementation passes;
2. every REQUIRED client adapter/surface passes;
3. equivalent authorization, state, data, freshness/provenance and user-meaning tests pass;
4. required macOS/Windows artifacts and Web/server deployment evidence pass;
5. no material parity debt remains.

There is no normal state of `Mac done -> Windows later -> Web later`. A single platform may be used for diagnosis or technical validation, but not as a product pilot or GA milestone.

Platform-specific releases remain valid when the actual responsibility is platform-specific; #64 / `v18.9.1` is the current example.

Temporary platform exceptions require a proven external blocker, explicit waiver/expiry, no misleading Delivered/GA status and a named recovery release.

## 3. Release-train delivery model

Major = strategic maturity generation. Minor band = dependency phase. Patch = one primary responsibility. Required platform adapters for one shared capability belong to that same patch.

No next shared domain begins while the current capability has material REQUIRED-platform parity debt.

## 4. v18.9.x delivery train

1. `v18.9.1` macOS runtime crash corrective
2. `v18.9.2` TradeInsight Settings/API-key UX
3. `v18.9.3` coverage-aware router
4. `v18.9.4` canonical company/instrument identity
5. `v18.9.5` Market Data Modes/diagnostics
6. `v18.9.6` provider observability/Adaptive telemetry
7. `v18.9.7` Form 4 SHADOW enrichment
8. `v18.9.8` symbol/company search
9. `v18.9.9` movers/ranking SHADOW evidence
10. `v18.9.10` remaining useful capability admission
11. `v18.9.11` Session-Aware Data Readiness Maintenance
12. `v18.9.12` professional closure

Shared native behavior remains semantically equivalent where applicable. Web becomes a REQUIRED shared client once the v19 cross-platform client foundation lands.

## 5. Protected Tier-0 contract

Pre-market, regular market and after-hours are protected. Live/current work has first provider/runtime/DB/worker claim. Maintenance/sync is bounded, preemptible/checkpointed and cannot flood protected sessions.

## 6. v19 delivery train

### v19.0.x — Hosted Foundations

- `v19.0.0` Provider Capability / Legal Rights Registry
- `v19.0.1` Hosted Tenant / Identity / Device / Session Control Plane
- `v19.0.2` DE.PULSE Product Entitlement / Metering Policy
- `v19.0.3` Account Data Governance / Privacy Lifecycle
- `v19.0.4` Hosted Environment / IaC / Service Trust Foundation
- `v19.0.5` PostgreSQL Tenancy / Schema / Pool / HA-PITR Foundation
- `v19.0.6` Managed Secrets / KMS Lifecycle
- `v19.0.7` Software Supply-Chain / Artifact & Dependency Assurance
- `v19.0.8` Provider Quality / Cost / Coverage / SLO Scorecards
- `v19.0.9` Data Reconciliation / Revision / Point-in-Time Quality

Shared hosted activation is blocked until applicable foundations pass.

### v19.1.x — Hosted Data Plane + Cross-Platform Account/State

- `v19.1.0` Hosted Provider Gateway
- `v19.1.1` Unified Serving Policy + Live Fan-Out
- `v19.1.2` Sync Protocol Foundation
- `v19.1.3` **Cross-Platform Account/Session Client Foundation — Mac + Windows + Web**
- `v19.1.4` **Cross-Platform Preferences — Mac + Windows + Web**
- `v19.1.5` **Cross-Platform Watchlists/Master Symbols — Mac + Windows + Web**
- `v19.1.6` **Cross-Platform Desks/Workspaces — Mac + Windows + Web**

No Mac product pilot exists. Each capability closes across all REQUIRED clients before the next shared capability proceeds.

### v19.2.x — Cross-Platform Shared Product + Assurance

- `v19.2.0` **Cross-Platform Research/Durable State — Mac + Windows + Web**
- `v19.2.1` **Cross-Platform Discovery/Opportunity Radar — Mac + Windows + Web**
- `v19.2.2` **Cross-Platform Market State/Market Modes/Readiness/Explanations — Mac + Windows + Web**
- `v19.2.3` **Cross-Platform Settings/Account/Device Controls — Mac + Windows + Web**
- `v19.2.4` **Cross-Platform RBAC/Product-Entitlement UX — Mac + Windows + Web**
- `v19.2.5` Tenant-Aware Metering / Cost / Usage Observability
- `v19.2.6` Mixed-Client Multi-User Security / Abuse / Capacity Hardening
- `v19.2.7` #66 Cross-Platform Hosted Sync / Gateway Assurance Closure

#66 is complete only after actual Mac + Windows + Web clients/services prove the architecture. Documentation or one-client success cannot close it.

### Required #66 evidence

At minimum:
- same account/session/role/product-entitlement truth across clients;
- same authorized preferences/watchlists/desks/workspaces/research state;
- same Discovery/Opportunity Radar and Market Mode/readiness/explanation meaning for equivalent evidence;
- same provider/freshness/provenance/degraded/UNKNOWN truth;
- same revocation/downgrade/denial outcomes;
- native offline/reconnect convergence without Web becoming a second authority;
- mixed-client concurrent edit/conflict/delete/tombstone behavior;
- provider-right/product-plan downgrade behavior;
- account export/deactivation/deletion/retention behavior;
- secret rotation and DB failover/PITR/restore;
- mixed-version protocol behavior;
- environment/config drift and artifact/dependency provenance failure behavior;
- noisy-neighbor/fairness/backpressure/protected-session pressure;
- no provider secret or unnecessary user-data leakage.

### v19.3.x — Point-in-Time Evidence

- `v19.3.0` Institutional/13F infrastructure
- `v19.3.1` Two-sided Long/Short evidence substrate
- `v19.3.2` AODR candidate/ranking/outcome lineage

Any user-facing output from these substrates follows the lockstep definition of Delivered.

### v19.4.x — Reliability / Economics / v20 Readiness

- `v19.4.0` ADR-GDI professional reliability/capacity/runbooks
- `v19.4.1` specialized/paid-provider gap evaluation
- `v19.4.2` v20 research-readiness audit

### v19.5.0 — Major Closure

No feature scope. Require #66 PASS, zero material shared-capability parity debt, tenant/RBAC/product-entitlement/provider-right/privacy separation, data lifecycle, environment/IaC, supply-chain provenance, API compatibility, SLO/capacity, recovery/rollback and actual required Mac/Windows/Web/runtime/deployment evidence.

## 7. G0-G16 delivery enforcement

- G1 freezes REQUIRED/N/A platform matrix.
- G4 requires every REQUIRED client implementation for shared scope.
- G6 proves cross-platform state/API integration.
- G7 proves equivalent authorization/data/privacy/right outcomes.
- G8 proves mixed-client load/recovery where relevant.
- G9 proves functionality/meaning equivalence.
- G10 blocks freeze on material parity debt.
- G13/G14 validate actual required native/Web artifacts/deployments.
- G15 forbids shared-capability GA until all REQUIRED clients pass.
- G16 audits parity drift and temporary waivers.

## 8. v20 delivery train

`v20.0.x` control/model governance -> `v20.1.x` ASBI -> `v20.2.x` adaptive Institutional/TDTI -> `v20.3.x` AODR -> `v20.4.0` adaptive operations -> `v20.5.0` closure.

Every shared adaptive user-facing capability follows Mac + Windows + Web lockstep and `SHADOW -> VALIDATED -> APPROVED -> PRODUCTION`. No Execution remains permanent.

## 9. G15 promotion guidance

Promotion for shared capabilities is capability-wide, not platform-by-platform. Canary/bounded cohorts are allowed, but a Mac-only GA followed by Windows/Web catch-up is prohibited.

## 10. Permanent boundaries

U.S. Equities Processing; GLD/SLV/USO actionable exceptions; No Execution; Smart Provider Router v2 sole routing authority; canonical freshness/recovery sole freshness owner; existing multi-feed allocator sole subscription owner; BroadSnapshotBroker canonical reuse owner; canonical persistence/cache/state/session/calendar/identity owners; direct SEC/EDGAR authoritative; G0-G16 only.

## Exactly one next action

Obtain/reproduce #64 macOS crash evidence and freeze narrow `v18.9.1` G1. Do not create `v18.9.2` or v19 product implementation branches until ordering permits it.