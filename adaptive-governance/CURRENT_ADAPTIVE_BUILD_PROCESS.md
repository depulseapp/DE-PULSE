# DE.PULSE — Current Adaptive Build Process

**Certified Stable:** `v18.9.0-stable`  
**Active product branch/PR:** `v18.9.1-development` / Draft PR #69  
**Active product scope:** #64 / `ADAPT-RUNTIME-CRASH-001` -> `v18.9.1`  
**Mandatory next process-hardening workstream:** #70 / `ADAPT-CI-CONVERGENCE-001` after truthful v18.9.1 closure  
**Master corrective program:** #65 / `ADAPT-PROVIDER-INTELLIGENCE-010`  
**Hosted architecture program:** #66 / `ADAPT-HOSTED-SYNC-001`  
**Governance alignment:** merged PR #67.

## 1. Permanent engineering loop

`Understand -> impact-map -> reuse -> decompose -> design -> implement -> observe -> qualify -> recover/fix -> certify -> deliver -> learn`

G0-G16 is permanent. One primary responsibility per work slice. Canonical owners are reused; known misses are fixed or durably assigned; exactly one next action closes each work slice.

The CI/release process follows the same adaptive principle as the product: preserve trustworthy evidence, remove redundant work, fail closed on uncertainty, and never add parallel owners merely for convenience.

## 2. Cross-Platform Lockstep Process

DE.PULSE is one product across macOS, Windows and Web. Shared capability work is decomposed by capability, not by client.

At G1 freeze:

`Capability -> canonical owner -> macOS REQUIRED/N/A -> Windows REQUIRED/N/A -> Web REQUIRED/N/A`

Rules:
- all REQUIRED client adapters/surfaces are part of the same capability scope;
- one canonical domain/API/state contract drives all clients;
- platform mechanics may differ, but business logic, intelligence, authorization, product entitlement, provider-right decisions, state semantics, provenance/freshness and explanation meaning may not;
- single-platform diagnosis is not a product pilot or delivery milestone;
- unresolved material parity debt blocks G10/G15 and the next shared domain;
- temporary platform exceptions require an external blocker, explicit waiver/expiry, no GA claim and a durable recovery release;
- platform-specific corrective releases remain valid where the defect itself is platform-specific.

Hosted invariant:

`rights/identity/RBAC/product entitlement/privacy/IaC/recovery/secrets/supply-chain -> provider boundary -> serving policy -> sync protocol -> cross-platform account/state capability -> equivalence proof -> next capability -> mixed-client assurance`

## 3. G0-G16 operating process

### G0 — Exact Baseline
Re-fetch live GitHub head, handoff, issues/comments, release identity and branch/PR state. Compare since certified baseline to prevent duplication.

### G1 — Immutable Scope
Freeze one responsibility, explicit non-goals, product-version disposition, work-slice ID, acceptance, rollback and platform applicability matrix. A REQUIRED client cannot be deferred for convenience.

### G2 — Architecture / Data Utility
Map canonical owners, tenant/data/privacy classification, trust boundaries and platform adapters. No client or provider becomes an independent business-truth owner.

### G3 — Design / Dependency Readiness
Freeze one domain/API/state contract, platform adapter contracts, compatibility/version/deprecation policy, equivalence tests, migration/rollback, retention/export/deletion, IaC/service trust, dependency/provenance, SLO/observability, conflict/idempotency/retry and negative/load/failure tests.

For CI/release tooling, also freeze:
- exact changed-path/change-class mapping;
- selected evidence graph and dependency invalidation rules;
- explicit fail-closed fallback;
- trusted base SHA/merge-base for manual/reusable runs;
- release/publication feasibility checks;
- immutable artifact/evidence ownership.

### G4 — Development Exit
Canonical implementation plus every REQUIRED client adapter/surface exists and has unit/contract tests. Single-platform completion cannot satisfy a shared capability.

For CI tooling, permanent executors/owners must be changed rather than adding version-specific workflow families.

### G5 — FAST Qualification
Run affected deterministic tests. Share common fixtures/contracts; test only platform-specific deltas separately.

After #70, Fast/Qualified selection must be dependency-aware and explicit. Unknown/mixed risk fails closed to broader evidence.

### G6 — Integration / MEDIUM Qualification
Prove cross-owner and cross-platform state/API integration, persistence/cache reuse, provider fallback/coverage and convergence/conflict behavior.

### G7 — Data / Security / Adaptive Intelligence
Prove equivalent tenant/RBAC/product-entitlement/provider-right/privacy/data outcomes across clients, including downgrade/revocation/denial and SHADOW/promotion rules.

### G8 — Performance / Capacity / Stability
Exercise mixed-client load where relevant: fairness/noisy-neighbor, provider headroom, DB/pool behavior, queues/backpressure, restart/warm-start, failover and protected-session reserves.

CI/release process changes also measure runner-minutes, duplicate runs, cancellations/supersession, queue pressure and avoided reruns without treating cost as permission to skip required evidence.

### G9 — Cross-Module / UI / UX
Audit all REQUIRED clients. Responsive/native interaction may differ; functionality, state meaning, role composition, explanations, freshness/provenance and authorization may not materially drift.

### G10 — Pre-Freeze Qualification
Unresolved P0 security/privacy/rights/environment/supply-chain/recovery/compatibility/duplicate-owner issue or material REQUIRED-platform parity gap blocks freeze.

For CI/release changes, G10 must also prove that selected evidence is complete for the dependency graph and that no permanent gate was weakened.

### G11 — Immutable Release Candidate + Publication Feasibility
Freeze exact source/fingerprint, product version, work-slice set, build identity and required-platform matrix.

Before expensive G12/G13/G14 work, verify:
- version/build/predecessor compatibility;
- current Stable expectation;
- target tag/release state;
- existing asset/tag collision semantics;
- publication can proceed or be classified idempotently.

Impossible publication fails here.

### G12 — Full Certification
Replay required system-wide certification on immutable RC. After #70, use one canonical version-neutral G12 executor plus declarative release/capability manifests; do not create a new certification executor per version.

### G13 — Native/Hosted Packaging / Provenance
Produce required macOS/Windows artifacts and Web/server deployments from certified source with hashes/provenance. Targeted pre-release native rehearsals may occur in Qualified when impact requires them; G13 remains authoritative packaging evidence.

### G14 — Actual Artifact Runtime Audit
Validate each REQUIRED actual artifact/deployment. One platform PASS never substitutes for another REQUIRED platform.

### G15 — Release Assurance / Promotion
No shared capability is GA/Delivered until every REQUIRED client passes. Publication consumes the exact same-run certified artifacts, does not rebuild, and may not overwrite differing Stable bytes. Same digest may be treated as idempotent no-op.

### G16 — Adaptive Retrospective / Handoff
Audit implementation misses, parity drift, incidents, metrics, privacy/supply-chain findings, cleanup and any temporary waivers. Handoff/issues must agree with executable evidence.

After #70, G16 also records:
- runner minutes and reruns avoided;
- evidence reused vs rerun and why it was still valid;
- toolchain/runner identity;
- release provenance/attestation status;
- confirmation that no required evidence was removed;
- residual versioned test/gate cleanup debt.

## 4. Failure and CI discipline

Classify before rerun: `PRODUCT_FAIL`, `GATE_TEST_FAIL`, `CI_HARNESS_FAIL`, `INFRA_FAIL`, `EXPECTED_NOOP`, `SUPERSEDED`.

`inspect -> isolate smallest affected owner/adapter -> preserve unrelated valid PASS -> repair -> rerun affected/dependent work only`

Rules:
- never weaken a gate or remove a REQUIRED platform to make a candidate green;
- exactly three routine workflows remain: Fast, Qualified, Release;
- no retry/certification/promotion branch families;
- manual/reusable qualification must bind a trustworthy complete delta;
- stale/exact-head-invalid evidence is never reused;
- Stable release execution is globally serialized;
- CI cost is optimized through dependency awareness, caching, evidence reuse and fewer duplicate builds, not by lowering assurance.

## 5. Current corrective state

#64 / `v18.9.1` G0 deterministic reproduction/root cause is complete. Do not restart G0.

Current blocker is hosted-runner execution. Once execution is available:
1. re-fetch live branch/PR/#64 state;
2. exact-head Fast PASS;
3. full exact-head Qualified PASS;
4. required backend/renderer/Chrome/WebKit PASS;
5. actual packaged macOS fresh native startup + dwell PASS;
6. warm relaunch on the same persisted profile PASS;
7. SQLite/profile reuse, protocol-resolution and deterministic cleanup PASS;
8. only then may PR #69 leave Draft.

Merge/release still requires explicit authorization and canonical G11-G16.

## 6. Mandatory process-hardening checkpoint — #70

After truthful v18.9.1 closure, execute `ADAPT-CI-CONVERGENCE-001` before the next product implementation.

The detailed executable target is `adaptive-governance/CURRENT_ADAPTIVE_CI_CONVERGENCE.md` plus issue #70.

Core process changes:
- Planner v3 dependency-aware job selection;
- explicit/manual full-delta base binding;
- targeted native rehearsal inside Qualified;
- separate WebKit/native-lifecycle owners with canonical script reuse;
- one canonical G12 executor;
- capability-oriented active test ownership;
- early G11 publication feasibility;
- digest-immutable Stable publication;
- global release serialization;
- verified repository protection/rulesets;
- public version/work-slice/source/build/evidence identity separation;
- prospective SemVer-or-custom decision;
- monotonic platform build number;
- canonical toolchain/provenance evidence;
- common machine-readable current-state truth for handoff/Adaptive overlays.

## 7. Product sequence after #70

The current remaining v18.9.x version labels are reservations pending #70 re-baselining. Dependency order is preserved even if public release grouping changes:

`TradeInsight Settings/API-key UX -> coverage-aware router -> canonical company/instrument identity -> Market Data Modes/diagnostics -> provider observability/Adaptive telemetry -> Form 4 SHADOW enrichment -> symbol/company search -> movers/ranking SHADOW evidence -> remaining useful capability admission -> Session-Aware Data Readiness Maintenance -> Professional Closure`

Each remains independently governed by scope/work-slice ID. Do not require each work slice to become a public Stable version merely for traceability.

## 8. Protected-session execution

Pre-market, regular market and after-hours are Tier-0. Live/current work has first provider/runtime/DB/network/worker claim. Maintenance/sync is bounded, preemptible/checkpointed and cannot flood protected sessions.

## 9. v19 / v20 process inheritance

v19 entry requires v18 native/data-plane closure plus #70 PASS. The existing hosted-foundation -> provider gateway/serving/sync -> cross-platform account/state -> shared product -> mixed-client assurance -> point-in-time evidence -> reliability ordering remains unchanged.

v20 remains governed adaptive intelligence under `SHADOW -> VALIDATED -> APPROVED -> PRODUCTION`, immutable experiment/model/prompt governance and No Execution.

## 10. Permanent owners / boundaries

Smart Provider Router v2 sole routing authority; canonical freshness/recovery sole freshness owner; existing multi-feed allocator sole subscription owner; BroadSnapshotBroker canonical reuse owner; canonical persistence/cache/state/session/calendar/identity owners; direct SEC/EDGAR authoritative; deterministic Day/Swing/Long protected; U.S. Equities Processing; GLD/SLV/USO actionable exceptions; No Execution.

## Exactly one next action

Restore/confirm GitHub Actions hosted-runner execution and finish exact-head #64 / `v18.9.1` qualification. After truthful v18.9.1 closure, execute #70 before any next product implementation.