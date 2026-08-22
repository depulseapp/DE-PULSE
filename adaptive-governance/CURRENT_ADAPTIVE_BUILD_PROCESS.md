# DE.PULSE — Current Adaptive Build Process

**Certified Stable:** `v18.9.0-stable`  
**Master corrective program:** #65 / `ADAPT-PROVIDER-INTELLIGENCE-010`  
**Hosted architecture program:** #66 / `ADAPT-HOSTED-SYNC-001`  
**Immediate product scope:** #64 / `ADAPT-RUNTIME-CRASH-001` -> `v18.9.1`  
**Governance alignment:** draft PR #67.

## 1. Permanent engineering loop

`Understand -> impact-map -> reuse -> decompose -> design -> implement -> observe -> qualify -> recover/fix -> certify -> deliver -> learn`

G0-G16 is permanent. No G17+ is permitted. Every patch has one primary responsibility, one canonical owner/dependency graph, exact predecessor identity, explicit non-goals, deterministic acceptance evidence, rollback/recovery behavior, durable issue/handoff state and exactly one next action at closure.

## 2. Cross-Platform Lockstep Process

DE.PULSE is one product across macOS, Windows and Web. Shared capability work is not decomposed by client unless the responsibility is genuinely platform-specific.

At G1 every patch freezes:

`Capability -> canonical owner -> macOS REQUIRED/N/A -> Windows REQUIRED/N/A -> Web REQUIRED/N/A`

Process invariants:
- one canonical domain/API/state contract drives all required clients;
- client adapters may differ for Keychain/Credential Manager/browser sessions, native SQLite/browser behavior, packaging, windowing, notifications and other platform concerns;
- business logic, intelligence, authorization, product entitlement, provider-right decisions, state transitions, provenance/freshness and user-visible meaning may not fork by platform;
- all REQUIRED client implementations belong to the same shared capability responsibility and must reach Development Exit before the shared capability can freeze;
- a single platform may be used temporarily for debugging or adapter validation, but it does not satisfy delivery and may not be called a product pilot;
- no new shared domain begins while a required client has unresolved material parity debt in the current domain;
- temporary platform exceptions require an external blocker, explicit waiver/expiry, no GA/Delivered claim and a durable recovery scope;
- platform-specific corrective releases remain valid where the actual defect is platform-specific, such as #64 / v18.9.1.

Hosted process invariant:

`rights + tenant identity + RBAC + product entitlement + privacy + environment/IaC + persistence recovery + secrets + supply-chain assurance -> provider boundary -> serving policy -> sync protocol -> cross-platform account/session -> cross-platform domain capability -> equivalence certification -> multi-user hardening -> closure`

Adaptive process invariant:

`experiment/control plane -> model/prompt governance -> adaptive capability implemented against canonical shared truth -> Mac/Windows/Web lockstep where user-facing -> calibration -> bounded promotion -> closure`

## 3. G0-G16 operating process

### G0 — Exact Baseline
Re-fetch live GitHub head, current handoff, issues/comments, release identity, predecessor evidence and existing PR/branch state. Compare commits since the last certified baseline so work is not duplicated.

### G1 — Immutable Scope
Freeze one primary responsibility, explicit non-goals, target version, scope ID, affected contracts, acceptance, rollback and the Mac/Windows/Web applicability matrix. A shared capability cannot defer a REQUIRED client to a later roadmap version merely for convenience.

### G2 — Architecture / Data Utility
Map canonical owners and dependency blast radius. Run Functionality Utility checkpoint. For hosted/security/data work record authority, tenant/data/privacy classification, trust boundaries, threat-model scope, environment/service boundaries, failure assumptions and rejected alternatives. No client becomes an independent business-truth owner.

Hosted serving distinguishes:
1. tenant/account identity;
2. RBAC/capability authorization;
3. DE.PULSE product-plan entitlement/quota;
4. upstream provider legal/data rights;
5. privacy/data-governance policy.

### G3 — Design / Dependency Readiness
Freeze:
- one canonical API/domain/state contract;
- platform adapter contracts for every REQUIRED client;
- compatibility/version/deprecation policy;
- cross-platform equivalence tests;
- migration/rollback/roll-forward;
- data retention/export/deletion behavior;
- environment/IaC/service-trust design;
- dependency/SBOM/provenance policy;
- SLO/error-budget/observability;
- conflict/idempotency/retry semantics;
- negative/load/failure matrix;
- role/capability UI composition;
- feature flag/kill/canary controls where risk justifies them.

### G4 — Development Exit
Canonical implementation exists and all REQUIRED platform adapters/surfaces for the frozen capability compile and have unit/contract tests. A Mac-only implementation of a shared scope is not G4 complete. Platform-specific diagnostic branches/tests may exist transiently but cannot redefine scope.

### G5 — FAST Qualification
Run affected deterministic tests per changed owner/adapter. Reuse common fixtures/contracts to avoid three independent business-logic suites.

### G6 — Integration / MEDIUM Qualification
Cross-owner integration plus cross-platform state/API equivalence for REQUIRED clients. Prove persistence/cache reuse, shared-symbol/single-flight behavior and affected provider fallback/coverage.

### G7 — Data / Security / Adaptive Intelligence
Prove rights, provenance, point-in-time truth, tenant isolation, object/function authorization, product entitlement, provider-right enforcement, privacy/data lifecycle, secret boundaries, revocation and SHADOW/promotion rules. Equivalent requests from required clients must produce equivalent authorization/data outcomes.

### G8 — Performance / Capacity / Stability
Bounded load/soak/capacity evidence. Where relevant exercise concurrent Mac/Windows/Web usage, noisy-neighbor/fairness, provider headroom, DB/index/pool behavior, queues/backpressure, circuits, restart/warm-start, environment/config failure and protected-session reserve.

### G9 — Cross-Module / UI / UX
Audit all REQUIRED clients. Responsive/native interaction differences are allowed; functionality, state meaning, role composition, explanations, source/freshness truth and direct-route/API authorization may not materially drift. No clipping/overlap/hidden privileged payload.

### G10 — Pre-Freeze Qualification
Authoritative reconciliation. Any unresolved P0 security/privacy/rights/environment/supply-chain/recovery/compatibility/duplicate-owner issue **or material REQUIRED-platform parity gap** blocks freeze.

### G11 — Immutable Release Candidate
Freeze exact source/fingerprint plus required-platform matrix.

### G12 — Full Certification
Replay mandatory system-wide certification on immutable RC.

### G13 — Native/Hosted Packaging / Provenance
Produce required macOS/Windows artifacts and Web/server deployments from certified source with hashes/provenance. No hidden rebuild divergence.

### G14 — Actual Artifact Runtime Audit
Validate actual macOS Apple Silicon and Windows x64 artifacts where REQUIRED, plus actual Web/hosted deployment where REQUIRED. A PASS on one client cannot substitute for another REQUIRED client.

### G15 — Release Assurance / Promotion
A shared capability cannot be GA/Delivered until every REQUIRED client passes. Canary/bounded cohorts may be used across the release, but platform lag is not normal canary strategy. Rollback/kill-switch readiness remains mandatory where applicable.

### G16 — Adaptive Retrospective / Handoff
Review implementation misses, incidents, parity drift, calls/reruns avoided, cost/usefulness, privacy/supply-chain findings, cleanup and exact next release. Handoff/issues must state any temporary platform waiver explicitly; silent parity debt is forbidden.

## 4. Failure classification and recovery

Classify before rerun:
- `PRODUCT_FAIL`
- `GATE_TEST_FAIL`
- `CI_HARNESS_FAIL`
- `INFRA_FAIL`
- `EXPECTED_NOOP`
- `SUPERSEDED`

On failure:

`inspect actual state -> identify smallest affected owner/adapter -> preserve unrelated PASS evidence -> repair -> rerun affected/dependent work only -> continue`

Never weaken a gate or drop a REQUIRED platform to make a candidate green.

## 5. CI/resource efficiency

- one logical Build Coordinator owns release state;
- one product branch + PR per patch;
- no retry/certification branch families;
- shared domain fixtures/contracts once, platform-adapter tests where they differ;
- exact-head FAST once per coherent candidate, Qualified once when ready, one G11-G16 release when release-capable;
- metadata/checkpoint-only changes do not wake full qualification when fingerprints/contracts are equivalent;
- superseded runs stop consuming resources.

## 6. Protected-session execution contract

Pre-market, regular market and after-hours always outrank maintenance/background sync. Live/current intelligence has first provider/runtime/DB/network/worker claim. Maintenance/sync uses bounded surplus capacity, is checkpointed/preemptible and cannot flood protected sessions.

## 7. v18.9.x process sequence

1. `v18.9.1` runtime crash — macOS-specific corrective.
2. `v18.9.2` TradeInsight Settings/API-key UX.
3. `v18.9.3` coverage-aware persistence-first Smart Provider Router.
4. `v18.9.4` canonical company/instrument identity.
5. `v18.9.5` Market Data Modes/capability diagnostics.
6. `v18.9.6` provider observability/Adaptive telemetry.
7. `v18.9.7` TradeInsight Form 4 SHADOW enrichment.
8. `v18.9.8` TradeInsight symbol/company search.
9. `v18.9.9` TradeInsight movers/ranking SHADOW evidence.
10. `v18.9.10` remaining useful capability admission.
11. `v18.9.11` session-aware Data Readiness Maintenance.
12. `v18.9.12` professional closure.

Shared native capabilities remain canonical and equivalent across supported native clients where applicable; Web becomes a REQUIRED client for shared product work once the v19 cross-platform client foundation lands.

## 8. v19 process sequence

### v19.0.x — Governance / Control Plane / Data Foundation
1. `v19.0.0` provider legal-rights registry;
2. `v19.0.1` tenant/identity/device/session control plane;
3. `v19.0.2` product entitlement/metering policy;
4. `v19.0.3` account data governance/privacy lifecycle;
5. `v19.0.4` hosted environment/IaC/service-trust foundation;
6. `v19.0.5` PostgreSQL tenancy/schema/pool/HA/PITR;
7. `v19.0.6` managed secrets/KMS;
8. `v19.0.7` software supply-chain/artifact/dependency assurance;
9. `v19.0.8` provider SLO/cost/coverage scorecards;
10. `v19.0.9` reconciliation/revision/point-in-time quality.

### v19.1.x — Hosted Data Plane + Cross-Platform Account/Sync
1. `v19.1.0` zero-key Hosted Provider Gateway;
2. `v19.1.1` unified serving policy/live fan-out;
3. `v19.1.2` sync protocol foundation;
4. `v19.1.3` **Mac + Windows + Web account/session client foundation**;
5. `v19.1.4` **Mac + Windows + Web preferences/watchlists**;
6. `v19.1.5` **Mac + Windows + Web desks/workspaces**.

No Mac product pilot exists. Technical single-platform validation is internal only and cannot satisfy phase exit.

### v19.2.x — Cross-Platform Shared Product + Assurance
1. `v19.2.0` **Mac + Windows + Web Research/durable state**;
2. `v19.2.1` **Mac + Windows + Web Market Intelligence/Discovery/Market Modes**;
3. `v19.2.2` **Mac + Windows + Web Settings/RBAC/product-entitlement UX**;
4. `v19.2.3` tenant-aware metering/cost/usage observability;
5. `v19.2.4` mixed-client security/abuse/noisy-neighbor/capacity hardening;
6. `v19.2.5` #66 cross-platform adversarial/failure/recovery closure.

### v19.3.x — Point-in-Time Evidence
1. `v19.3.0` institutional/13F infrastructure;
2. `v19.3.1` two-sided Long/Short evidence substrate;
3. `v19.3.2` AODR candidate/ranking/outcome lineage.

Any user-facing capability follows lockstep delivery.

### v19.4.x — Reliability / Economics / Readiness
1. `v19.4.0` ADR-GDI reliability/capacity/runbooks;
2. `v19.4.1` specialized/paid-provider gap evaluation;
3. `v19.4.2` v20 research-readiness audit.

### v19.5.0 — Major Closure
No feature scope. Zero unresolved P0 and zero material cross-platform parity debt for shared capabilities.

## 9. #66 hosted/sync process requirements

#66 is executable architecture scope, not documentation. Mac, Windows and Web consume the same account, authorization, provider, state and intelligence truth. Native offline mechanics may differ from Web, but semantic equivalence is required.

Canonical serving path:

`authenticated tenant/user/device -> RBAC -> product entitlement -> provider rights -> privacy policy -> canonical cache/persistence/freshness -> residual need -> Smart Provider Router v2 -> server-side credential -> authorized projection`

Sync/state path:

`native SQLite mutation + outbox / Web hosted mutation -> authenticated idempotent server apply -> PostgreSQL authoritative revision/change sequence -> authorized client projection/checkpoint`

No raw DB replication, no client-specific provider/router/session truth and no platform-specific business logic fork.

## 10. v20 process sequence

### v20.0.x — Control & Governance
Research/experiment ledger, then model/prompt governance + Champion/Challenger, then analogue/calibration work.

### v20.1.x — ASBI
Behavioral fingerprints/state transitions, then scenarios/probability momentum/calibration.

### v20.2.x — Institutional + TDTI
Adaptive 13F, competing Long/Short/No Reliable Edge, two-sided trade-plan validation.

### v20.3.x — AODR
Shared opportunity ranking, then diversity/opportunity-cost/personal relevance after shared truth.

### v20.4.x — Adaptive Operations
ADR-GDI adaptive optimization remains SHADOW/Champion-Challenger until explicit promotion.

### v20.5.0 — Professional Closure
No feature scope. Cross-platform parity is mandatory for every shared adaptive user-facing capability. No Execution remains permanent.

## 11. Permanent owners / boundaries

Smart Provider Router v2 sole routing authority; canonical freshness/recovery sole freshness owner; existing multi-feed allocator sole subscription owner; BroadSnapshotBroker canonical reuse owner; canonical persistence/cache/state/session/calendar/identity owners; direct SEC/EDGAR authoritative; deterministic Day/Swing/Long protected; U.S. Equities Processing; GLD/SLV/USO actionable exceptions; No Execution.

## Exactly one next action

Run #64 / `v18.9.1` G0 from complete macOS crash evidence or deterministic reproduction and freeze the narrow G1 before product-source change. Do not start `v18.9.2` or v19 implementation until ordering permits it.