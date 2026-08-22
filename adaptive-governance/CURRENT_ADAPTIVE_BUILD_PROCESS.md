# DE.PULSE — Current Adaptive Build Process

**Certified Stable:** `v18.9.0-stable`  
**Master corrective program:** #65 / `ADAPT-PROVIDER-INTELLIGENCE-010`  
**Hosted architecture program:** #66 / `ADAPT-HOSTED-SYNC-001`  
**Immediate product scope:** #64 / `ADAPT-RUNTIME-CRASH-001` -> `v18.9.1`  
**Governance alignment:** draft PR #67.

## 1. Permanent engineering loop

`Understand -> impact-map -> reuse -> decompose -> design -> implement -> observe -> qualify -> recover/fix -> certify -> deliver -> learn`

G0-G16 is permanent. One primary responsibility per patch. Canonical owners are reused; known misses are fixed or durably assigned; exactly one next action closes each patch.

## 2. Cross-Platform Lockstep Process

DE.PULSE is one product across macOS, Windows and Web. Shared capability work is decomposed by capability, not by client.

At G1 freeze:

`Capability -> canonical owner -> macOS REQUIRED/N/A -> Windows REQUIRED/N/A -> Web REQUIRED/N/A`

Rules:
- all REQUIRED client adapters/surfaces are part of the same capability scope;
- one canonical domain/API/state contract drives all clients;
- platform mechanics may differ, but business logic, intelligence, authorization, product entitlement, provider-right decisions, state semantics, provenance/freshness and explanation meaning may not;
- a single platform may be used for diagnosis/adapter validation, but that is not a product pilot or delivery milestone;
- unresolved material parity debt blocks G10/G15 and blocks the next shared domain;
- temporary platform exceptions require an external blocker, explicit waiver/expiry, no GA claim and a durable recovery release;
- platform-specific corrective releases remain valid where the defect itself is platform-specific.

Hosted invariant:

`rights/identity/RBAC/product entitlement/privacy/IaC/recovery/secrets/supply-chain -> provider boundary -> serving policy -> sync protocol -> cross-platform account/state capability -> equivalence proof -> next capability -> mixed-client assurance`

## 3. G0-G16 operating process

### G0 — Exact Baseline
Re-fetch live GitHub head, handoff, issues/comments, release identity and branch/PR state. Compare since certified baseline to prevent duplication.

### G1 — Immutable Scope
Freeze one responsibility, explicit non-goals, target version, issue/scope ID, acceptance, rollback and platform applicability matrix. A REQUIRED client cannot be deferred to a later version for convenience.

### G2 — Architecture / Data Utility
Map canonical owners, tenant/data/privacy classification, trust boundaries and platform adapters. No client becomes an independent business-truth owner.

### G3 — Design / Dependency Readiness
Freeze one domain/API/state contract, platform adapter contracts, compatibility/version/deprecation policy, equivalence tests, migration/rollback, retention/export/deletion, IaC/service trust, dependency/provenance, SLO/observability, conflict/idempotency/retry and negative/load/failure tests.

### G4 — Development Exit
Canonical implementation plus every REQUIRED client adapter/surface exists and has unit/contract tests. Single-platform completion cannot satisfy a shared capability.

### G5 — FAST Qualification
Run affected deterministic tests. Share common fixtures/contracts; test only platform-specific deltas separately.

### G6 — Integration / MEDIUM Qualification
Prove cross-owner and cross-platform state/API integration, persistence/cache reuse, provider fallback/coverage and convergence/conflict behavior.

### G7 — Data / Security / Adaptive Intelligence
Prove equivalent tenant/RBAC/product-entitlement/provider-right/privacy/data outcomes across clients, including downgrade/revocation/denial and SHADOW/promotion rules.

### G8 — Performance / Capacity / Stability
Exercise mixed-client load where relevant: fairness/noisy-neighbor, provider headroom, DB/pool behavior, queues/backpressure, restart/warm-start, failover and protected-session reserves.

### G9 — Cross-Module / UI / UX
Audit all REQUIRED clients. Responsive/native interaction may differ; functionality, state meaning, role composition, explanations, freshness/provenance and authorization may not materially drift.

### G10 — Pre-Freeze Qualification
Unresolved P0 security/privacy/rights/environment/supply-chain/recovery/compatibility/duplicate-owner issue or material REQUIRED-platform parity gap blocks freeze.

### G11 — Immutable Release Candidate
Freeze exact source/fingerprint and required-platform matrix.

### G12 — Full Certification
Replay required system-wide certification on immutable RC.

### G13 — Native/Hosted Packaging / Provenance
Produce required macOS/Windows artifacts and Web/server deployments from certified source with hashes/provenance.

### G14 — Actual Artifact Runtime Audit
Validate each REQUIRED actual artifact/deployment. One platform PASS never substitutes for another REQUIRED platform.

### G15 — Release Assurance / Promotion
No shared capability is GA/Delivered until every REQUIRED client passes. Canary/cohort activation is allowed; platform lag is not normal canary strategy.

### G16 — Adaptive Retrospective / Handoff
Audit implementation misses, parity drift, incidents, metrics, privacy/supply-chain findings, cleanup and any temporary waivers. Handoff/issues must agree with executable evidence.

## 4. Failure and CI discipline

Classify before rerun: `PRODUCT_FAIL`, `GATE_TEST_FAIL`, `CI_HARNESS_FAIL`, `INFRA_FAIL`, `EXPECTED_NOOP`, `SUPERSEDED`.

`inspect -> isolate smallest affected owner/adapter -> preserve unrelated PASS -> repair -> rerun affected/dependent work only`

Never weaken a gate or remove a REQUIRED platform to make a candidate green. Use one product branch/PR per patch; no retry/certification branch families.

## 5. Protected-session execution

Pre-market, regular market and after-hours are Tier-0. Live/current work has first provider/runtime/DB/network/worker claim. Maintenance/sync is bounded, preemptible/checkpointed and cannot flood protected sessions.

## 6. v18.9.x process sequence

1. `v18.9.1` macOS runtime crash corrective
2. `v18.9.2` TradeInsight Settings/API-key UX
3. `v18.9.3` coverage-aware persistence-first Smart Provider Router
4. `v18.9.4` canonical company/instrument identity
5. `v18.9.5` Market Data Modes/capability diagnostics
6. `v18.9.6` provider observability/Adaptive telemetry
7. `v18.9.7` TradeInsight Form 4 SHADOW enrichment
8. `v18.9.8` TradeInsight symbol/company search
9. `v18.9.9` TradeInsight movers/ranking SHADOW evidence
10. `v18.9.10` remaining useful capability admission
11. `v18.9.11` session-aware Data Readiness Maintenance
12. `v18.9.12` professional closure

## 7. v19 process sequence

### v19.0.x — Hosted Foundations
1. `v19.0.0` provider legal-rights registry
2. `v19.0.1` tenant/identity/device/session control plane
3. `v19.0.2` product entitlement/metering policy
4. `v19.0.3` account data governance/privacy lifecycle
5. `v19.0.4` hosted environment/IaC/service trust
6. `v19.0.5` PostgreSQL tenancy/schema/pool/HA/PITR
7. `v19.0.6` managed secrets/KMS
8. `v19.0.7` software supply-chain/artifact/dependency assurance
9. `v19.0.8` provider SLO/cost/coverage scorecards
10. `v19.0.9` reconciliation/revision/point-in-time quality

### v19.1.x — Hosted Data Plane + Cross-Platform Account/State
1. `v19.1.0` zero-key Hosted Provider Gateway
2. `v19.1.1` unified serving policy/live fan-out
3. `v19.1.2` sync protocol foundation
4. `v19.1.3` **Mac + Windows + Web account/session client foundation**
5. `v19.1.4` **Mac + Windows + Web preferences**
6. `v19.1.5` **Mac + Windows + Web watchlists/master symbols**
7. `v19.1.6` **Mac + Windows + Web desks/workspaces**

No Mac product pilot exists.

### v19.2.x — Cross-Platform Shared Product + Assurance
1. `v19.2.0` **Mac + Windows + Web Research/durable state**
2. `v19.2.1` **Mac + Windows + Web Discovery/Opportunity Radar**
3. `v19.2.2` **Mac + Windows + Web Market State/Market Modes/Readiness/Explanations**
4. `v19.2.3` **Mac + Windows + Web Settings/Account/Device Controls**
5. `v19.2.4` **Mac + Windows + Web RBAC/Product-Entitlement UX**
6. `v19.2.5` tenant-aware metering/cost/usage observability
7. `v19.2.6` mixed-client security/abuse/noisy-neighbor/capacity hardening
8. `v19.2.7` #66 cross-platform adversarial/failure/recovery closure

### v19.3.x — Point-in-Time Evidence
`v19.3.0` 13F infrastructure -> `v19.3.1` two-sided evidence -> `v19.3.2` AODR lineage. Any user-facing output follows lockstep.

### v19.4.x — Reliability / Economics / Readiness
`v19.4.0` ADR-GDI reliability/capacity -> `v19.4.1` paid-provider gap evaluation -> `v19.4.2` v20 research-readiness audit.

### v19.5.0 — Major Closure
No feature scope. Zero unresolved P0 and zero material cross-platform parity debt for shared capabilities.

## 8. v20 process sequence

`v20.0.x` control/model governance -> `v20.1.x` ASBI -> `v20.2.x` adaptive Institutional/TDTI -> `v20.3.x` AODR -> `v20.4.0` adaptive operations -> `v20.5.0` closure.

Every shared adaptive user-facing capability follows lockstep. No Execution remains permanent.

## 9. Permanent owners / boundaries

Smart Provider Router v2 sole routing authority; canonical freshness/recovery sole freshness owner; existing multi-feed allocator sole subscription owner; BroadSnapshotBroker canonical reuse owner; canonical persistence/cache/state/session/calendar/identity owners; direct SEC/EDGAR authoritative; deterministic Day/Swing/Long protected; U.S. Equities Processing; GLD/SLV/USO actionable exceptions; No Execution.

## Exactly one next action

Run #64 / `v18.9.1` G0 from complete macOS crash evidence or deterministic reproduction and freeze its narrow G1. Do not start `v18.9.2` or v19 implementation until ordering permits it.