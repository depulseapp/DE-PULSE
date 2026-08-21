# DE.PULSE — Current Adaptive Build Process

**Certified Stable:** `v18.9.0-stable`  
**Master corrective program:** #65 / `ADAPT-PROVIDER-INTELLIGENCE-010`  
**Hosted architecture program:** #66 / `ADAPT-HOSTED-SYNC-001`  
**Immediate product scope:** #64 / `ADAPT-RUNTIME-CRASH-001` -> `v18.9.1`  
**Governance alignment:** draft PR #67.

## 1. Permanent engineering loop

DE.PULSE engineering follows:

`Understand -> impact-map -> reuse -> decompose -> design -> implement -> observe -> qualify -> recover/fix -> certify -> deliver -> learn`

G0-G16 is permanent. No G17+ is permitted. Heavy work is decomposed into the smallest independently meaningful evidence units when that improves diagnosis, recovery, traceability or cost.

Every patch has:
- one primary responsibility;
- one canonical owner/dependency graph;
- exact predecessor/source identity;
- explicit non-goals;
- named data/provider/security/runtime/UI consumers;
- deterministic acceptance evidence;
- rollback/recovery/invalidation behavior;
- durable issue/handoff state;
- exactly one next action at closure.

## 2. Version/process alignment

- Major versions define strategic maturity generations.
- Minor bands define coherent dependency phases.
- Patch releases remain one-primary-responsibility units.
- Planned future numbers are reservations only and freeze at G1.
- G0/G1 may split a broad future item and shift unstarted reservations; shipped releases never move.
- Corrective/security issues preempt the roadmap when necessary.

General process invariant:

`stability -> canonical ownership -> observability -> SHADOW validation -> expansion -> operational hardening -> closure`

Hosted process invariant:

`tenant identity + RBAC + product entitlement + provider rights + persistence recovery + secret management -> hosted provider boundary -> unified serving policy -> sync transport -> client activation -> platform parity -> tenant-aware metering/capacity -> adversarial assurance -> evidence expansion`

Adaptive process invariant:

`experiment/control plane -> model/prompt governance -> research features -> calibration -> bounded promotion -> closure`

## 3. G0-G16 operating process

### G0 — Exact Baseline
Re-fetch live GitHub head, current handoff, issues/comments, release identity, predecessor evidence and existing PR/branch state. Compare commits since the last certified baseline so completed work is not duplicated.

### G1 — Immutable Scope
Freeze one primary responsibility, explicit non-goals, target version, issue/scope ID, affected contracts, acceptance and rollback disposition. Classify governance obligations as current blocker/process hardening/next mandatory/future strategic.

### G2 — Architecture / Data Utility
Map canonical owners and dependency blast radius. Run Functionality Utility checkpoint: purpose, consumer, reuse, correlation, freshness/materiality, retention/rights, cost, surface placement and retirement disposition.

For material hosted/security/data decisions, record an ADR or equivalent durable decision evidence including:
- authority/ownership boundary;
- tenant/account/data classification;
- trust boundaries and threat-model scope;
- failure/availability assumptions;
- alternatives rejected and migration implications.

Hosted designs must explicitly distinguish:
1. tenant/account identity;
2. RBAC/capability authorization;
3. DE.PULSE product-plan entitlement/quota policy;
4. upstream provider legal/data rights.

These may never collapse into one generic entitlement flag.

### G3 — Design / Dependency Readiness
Freeze executable contracts before coding:
- API/schema/protocol and version compatibility;
- API inventory/ownership/deprecation policy for hosted endpoints;
- tenant/RBAC/product-entitlement/provider-right serving order;
- migration and rollback/roll-forward plan;
- conflict/idempotency/retry semantics;
- SLO/error-budget/tenant-aware observability requirements;
- test fixtures/replay/live-provider boundaries;
- failure-injection/load/negative-test matrix;
- UI/role/capability composition;
- feature flag/kill switch/canary controls where risk justifies them.

### G4 — Development Exit
Implementation exists in the canonical owner, compiles and has unit/contract tests. No convenience duplication or hidden parallel state owner remains. Backward-compatible migrations use expand/contract where practical.

### G5 — FAST Qualification
Affected-area deterministic tests first. No broad expensive rerun when the delta is narrow.

### G6 — Integration / MEDIUM Qualification
Cross-owner integration, persistence/cache reuse, shared-symbol/single-flight behavior, provider fallback/coverage and state transitions for affected scope.

### G7 — Data / Security / Adaptive Intelligence
Prove data rights, provenance, point-in-time behavior, tenant isolation, object/function-level authorization, product-entitlement enforcement, provider-right enforcement, secret/redaction boundaries, revocation and SHADOW/promotion rules.

Negative tests are mandatory for security-sensitive scope, including:
- unauthorized route/API/SSE access;
- cross-account object access;
- role/capability downgrade;
- product-plan downgrade/suspension/quota exhaustion;
- revoked device/session;
- provider-right expiry/downgrade;
- secret leakage attempts;
- unsafe model influence.

### G8 — Performance / Capacity / Stability
Bounded load/soak/capacity evidence for affected responsibility. Where relevant include tenant-aware resource usage, noisy-neighbor/fairness tests, provider quota/headroom, DB/index/pool behavior, queues/backpressure, circuit breakers, memory/CPU/workers, restart/warm-start, failure injection, failover/recovery and protected-session reserve.

### G9 — Cross-Module / UI / UX
Audit affected surfaces across canonical role compositions and supported viewports. Direct-route/API authorization must match UI visibility. No blank role gaps, clipping/overlap, hidden privileged payload or duplicated deep-evidence surfaces.

### G10 — Pre-Freeze Qualification
Authoritative full coverage reconciliation. Every requirement is freshly evidenced or explicitly inherited from equivalent evidence. Any unresolved P0 security/tenant/product-entitlement/provider-right/recovery/compatibility/duplicate-owner issue blocks freeze.

### G11 — Immutable Release Candidate
Freeze exact source/fingerprint and candidate provenance.

### G12 — Full Certification
Replay mandatory system-wide certification on immutable RC. No source mutation.

### G13 — Native Packaging / Provenance
Produce required packages from certified RC with hashes/provenance; no hidden rebuild divergence.

### G14 — Actual Artifact Runtime Audit
Validate macOS Apple Silicon and Windows x64 independently where affected. Preserve an unchanged platform PASS only when exact RC/package identity and relevant assumptions remain equivalent.

### G15 — Release Assurance / Promotion
Consume the complete evidence graph. Hosted/high-risk changes use bounded/canary/progressive activation or equivalent controls where applicable, with explicit rollback/kill-switch readiness.

### G16 — Adaptive Retrospective / Handoff
Deep implementation-miss review, incident/failure classification, calls/reruns avoided, tenant/provider cost/usefulness metrics where applicable, obsolete machinery cleanup, scorecards, corrective learning and exact next release. Handoff/issues must agree with executable evidence.

## 4. Failure classification and recovery

Classify before rerun:
- `PRODUCT_FAIL`
- `GATE_TEST_FAIL`
- `CI_HARNESS_FAIL`
- `INFRA_FAIL`
- `EXPECTED_NOOP`
- `SUPERSEDED`

On failure:

`inspect actual state -> identify smallest affected package -> preserve unrelated PASS evidence -> repair -> rerun affected/dependent work only -> continue`

Never weaken a gate to make a candidate green. Repeated incidents become regressions/preflight controls where useful.

## 5. CI/resource efficiency

- one logical Build Coordinator owns authoritative release state;
- one product branch + PR per patch;
- no retry/certification branch families;
- coherent source/test batch before PR;
- exact-head FAST once per coherent candidate, Qualified once when ready, one G11-G16 release when release-capable;
- share fixtures/canonical evidence/provider acquisition across independent lanes;
- deterministic replay/cached evidence when live behavior is not under test;
- true provider/model calls only when live behavior must be certified;
- metadata/checkpoint-only changes do not wake full product qualification when fingerprints/contracts are unchanged;
- superseded runs stop consuming authoritative status/resources.

## 6. Protected-session execution contract

The canonical U.S. market calendar/session owner defines protected sessions.

**Pre-market, regular market and after-hours** always outrank maintenance and background synchronization.

During protected sessions:
- live/current intelligence has first provider/runtime/DB/network/worker claim;
- maintenance/sync uses bounded surplus capacity only;
- low-priority external acquisition suspends unless directly required by a live consumer;
- heavy reconciliation/compaction/backfill is prohibited;
- maintenance/sync is preemptible/checkpointed and yields promptly to current work/market shock.

Missed maintenance/sync catch-up waits for an eligible bounded window; it cannot flood a protected session.

## 7. v18.9.x process sequence

1. `v18.9.1` runtime crash
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
12. `v18.9.12` professional closure audit

**Audit change:** telemetry moved before provider capability expansion so subsequent SHADOW usefulness/promotion is measurable.

## 8. v19 process sequence

### v19.0.x — Governance / Control Plane / Data Foundation
Mandatory order:
1. `v19.0.0` provider capability/legal-rights registry;
2. `v19.0.1` tenant/identity/device/session control plane;
3. `v19.0.2` DE.PULSE product entitlement/metering policy;
4. `v19.0.3` PostgreSQL tenancy/schema/pool/HA/PITR foundation;
5. `v19.0.4` managed secrets/KMS;
6. `v19.0.5` provider quality/cost/coverage/SLO scorecards;
7. `v19.0.6` reconciliation/revision/point-in-time data quality.

No broad hosted provider/sync activation until required tenant/security/rights/product-entitlement/recovery foundations pass.

### v19.1.x — Hosted Data Plane / Sync Foundation
1. `v19.1.0` authenticated versioned zero-key Provider Gateway using existing Smart Provider Router v2;
2. `v19.1.1` unified serving policy + live fan-out isolation across tenant/RBAC/product-entitlement/provider-right dimensions;
3. `v19.1.2` typed sync protocol/bootstrap/outbox/idempotency/change-log/checkpoint/tombstone/compaction/mixed-version contract;
4. `v19.1.3` macOS preferences/watchlist pilot + local-account/offline/lost-device behavior;
5. `v19.1.4` desks/workspaces sync.

### v19.2.x — Cross-Platform Parity / #66 Assurance
1. `v19.2.0` Windows parity;
2. `v19.2.1` hosted web parity;
3. `v19.2.2` rights-aware research/state portability;
4. `v19.2.3` tenant-aware metering/cost/usage observability;
5. `v19.2.4` multi-user security/abuse/noisy-neighbor/capacity hardening;
6. `v19.2.5` #66 adversarial/failure/recovery closure audit.

### v19.3.x — Point-in-Time Evidence
1. `v19.3.0` institutional/13F infrastructure;
2. `v19.3.1` two-sided Long/Short evidence substrate;
3. `v19.3.2` AODR candidate/ranking/outcome lineage.

### v19.4.x — Reliability / Economics / Readiness
1. `v19.4.0` ADR-GDI professional reliability/capacity + operational runbook readiness;
2. `v19.4.1` specialized/paid-provider gap evaluation;
3. `v19.4.2` v20 research-readiness audit.

### v19.5.0 — Major Closure
No feature scope. Full data/security/tenant/product-entitlement/provider-right/operational/commercial/hosted/account/package zero-miss closure.

## 9. #66 hosted/sync process requirements

#66 is executable architecture scope, not a documentation project.

Canonical serving path:

`authenticated tenant/user/device -> RBAC -> DE.PULSE product entitlement -> provider legal/data rights -> canonical cache/persistence/freshness -> exact residual need -> Smart Provider Router v2 -> server-side provider credential -> authorized projection/fan-out`

Sync path:

`SQLite local mutation + atomic outbox -> authenticated idempotent server apply -> PostgreSQL authoritative revision/change sequence -> incremental pull -> transactional SQLite apply -> checkpoint advance`

Required controls:
- first-class tenant/account identity on hosted requests/data/metrics;
- server-side zero-key customer model;
- canonical roles/capabilities + revocation;
- separate product-plan entitlement/quota policy;
- separate provider legal/data rights enforcement;
- entitlement/right-aware cache/live fan-out;
- stable IDs/revisions/tombstones;
- new-device bootstrap + stale-checkpoint re-bootstrap;
- retention/compaction/inactive-device behavior;
- local account isolation/lost-device behavior;
- secret lifecycle/rotation/compromise recovery;
- API inventory/version/deprecation policy;
- mixed-version sync protocol;
- PostgreSQL HA/PITR/restore/migration safety;
- tenant-aware usage/cost/health/abuse telemetry;
- noisy-neighbor/fairness/resource controls;
- failure injection and recovery proof;
- no raw SQLite/PostgreSQL dual-master replication;
- no parallel provider/router/subscription/freshness/persistence/session owner.

## 10. v20 process sequence

### v20.0.x — Control & Governance
1. adaptive research/experiment ledger;
2. model/prompt governance + Champion/Challenger **before broad adaptive rollout**;
3. historical analogues/regime outcomes;
4. calibration/FP-FN/miss/contradiction/drift.

### v20.1.x — ASBI
Behavioral fingerprints/state transitions, then scenarios/probability momentum/calibration.

### v20.2.x — Institutional + TDTI
Adaptive 13F, then competing Long/Short/No Reliable Edge, then two-sided trade-plan validation.

### v20.3.x — AODR
Shared opportunity ranking, then diversity/opportunity-cost/personalized relevance after shared truth.

### v20.4.x — Adaptive Operations
ADR-GDI adaptive provider/recovery/workload/maintenance/reserve optimization remains SHADOW/Champion-Challenger until explicit promotion.

### v20.5.0 — Professional Closure
No feature scope. Calibrated utility, abstention, deterministic-boundary protection, privacy/security/data rights, reproducibility/rollback, actual artifacts and No Execution.

## 11. Permanent owners / boundaries

Smart Provider Router v2 sole routing authority; canonical freshness/recovery sole freshness owner; existing multi-feed allocator sole subscription owner; BroadSnapshotBroker canonical reuse owner; canonical persistence/cache/state owners; canonical U.S. market calendar/session owner; direct SEC/EDGAR authoritative; canonical tenant/identity/role/capability truth; deterministic Day/Swing/Long protected; U.S. Equities Processing; GLD/SLV/USO actionable exceptions; No Execution.

## Exactly one next action

Run #64 / `v18.9.1` G0 from complete macOS crash evidence or deterministic reproduction and freeze the narrow G1 before any product-source change. Do not start `v18.9.2` or any v19 implementation branch until current release ordering permits it.