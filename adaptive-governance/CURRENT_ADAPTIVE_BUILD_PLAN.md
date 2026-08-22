# DE.PULSE — Current Adaptive Build Plan

**Immutable Stable:** `v18.9.0-stable`  
**Stable candidate:** `9ea81cddae4875ae15d3719ca028519a36c597b6`  
**Stable fingerprint:** `a8719090c341c874dbd1279cc31ad98e84075d5701c46a800bf951340780ecb9`  
**Active product development branch:** none  
**Active product PR:** none  
**Governance alignment PR:** #67 (draft)  
**Master corrective program:** #65 / `ADAPT-PROVIDER-INTELLIGENCE-010`  
**Hosted architecture program:** #66 / `ADAPT-HOSTED-SYNC-001`  
**Immediate next product patch:** `v18.9.1` / #64 / `ADAPT-RUNTIME-CRASH-001`.

## 1. Build philosophy

DE.PULSE uses one permanent G0-G16 release model and small, dependency-ordered, independently certifiable patches.

Canonical engineering order:

`stabilize -> establish canonical owners -> instrument -> validate in SHADOW -> expand -> operationalize -> close -> establish hosted control/privacy/environment plane -> establish recoverable persistence/secrets/supply chain -> activate hosted data plane -> expose shared capability cross-platform -> certify lockstep parity -> harden -> build evidence substrate -> learn`

Permanent rules:
- one primary responsibility per patch;
- `REUSE -> CONSOLIDATE -> REFACTOR -> DELETE/REPLACE -> ADD` before new machinery;
- canonical owners are extended, never forked for a provider/client/feature;
- observability required to judge a capability lands before broad capability admission;
- tenant identity, security, provider rights, product entitlement, data-governance/privacy lifecycle, hosted environment/IaC, persistence recovery, secret management and software-supply-chain assurance precede commercial hosted activation;
- migrations use backward-compatible expand/contract where practical and always have rollback/roll-forward/recovery disposition before activation;
- high-risk externally visible behavior uses bounded feature/capability flags, kill switches, canary/progressive activation or equivalent controls where useful;
- governance closure follows `Governed -> Implemented -> Enforced -> Evidenced -> Delivered -> Learned`;
- known misses are fixed in-scope or durably assigned to a named later release;
- one product branch + one product PR per patch; CI reruns only the smallest trustworthy affected set;
- G0-G16 is the only top-level release model.

## 2. Cross-Platform Lockstep Product Contract

DE.PULSE is one product across **macOS, Windows and Web**. Shared capabilities are designed once against canonical domain/API/state owners and delivered across every applicable supported client in the same governed release responsibility.

For every G1 scope, freeze a platform applicability matrix:
- `macOS: REQUIRED | N/A(reason)`
- `Windows: REQUIRED | N/A(reason)`
- `Web: REQUIRED | N/A(reason)`

Rules:
- a shared capability cannot be declared Delivered/GA while any REQUIRED supported platform is materially behind;
- there is no roadmap pattern of `Mac product pilot -> Windows catch-up -> Web catch-up`;
- a single platform may be used internally for technical diagnosis, adapter validation or failure isolation, but that is not delivery and cannot satisfy G10/G15;
- platform-specific corrective work is allowed when the defect/responsibility itself is platform-specific, for example the macOS-only `v18.9.1` crash;
- platform adapters may differ only where OS/browser behavior requires it; business logic, intelligence, state semantics, authorization, data provenance/freshness and user-facing meaning remain canonical;
- Mac and Windows may use native SQLite/secure-storage adapters while Web may use browser/session mechanisms, but account truth and shared domain behavior remain equivalent;
- no new shared capability starts while the current shared capability has unresolved material platform parity debt;
- a temporary platform exception requires a named external blocker, explicit waiver/expiry, no misleading GA claim and a durable recovery target. It is never normal roadmap sequencing.

**Definition of Delivered for a shared capability:** canonical backend/domain implementation PASS + all REQUIRED client adapters/surfaces PASS + cross-platform equivalence PASS + actual artifact/deployment proof PASS.

## 3. Version alignment

- **Major (`v19`, `v20`)** = strategic maturity generation.
- **Minor band (`v19.0.x`, `v19.1.x`, etc.)** = coherent dependency phase with explicit entry/exit criteria.
- **Patch (`x.y.z`)** = one primary independently certifiable responsibility. A shared responsibility may include all REQUIRED platform adapters because those adapters are part of one capability, not separate features.
- Future numbers are planning reservations only until G1.
- If a future packet is too broad, split it in the same minor band and shift later unstarted reservations.
- Shipped versions are immutable.
- Corrective/security work can preempt the train; dependencies are reconciled before planned work resumes.

## 4. v18.9.x build train — trustworthy runtime/data plane

### v18.9.1 — Runtime Reliability
Scope owner #64. Diagnose/fix only the real packaged macOS Apple Silicon SIGABRT from evidence/reproduction. Preserve SQLite/user state/API keys. Require warm-state/relaunch lifecycle regression and actual packaged runtime proof. This is a legitimate platform-specific corrective because the escaped defect is macOS-specific. No provider/router/Market Mode/identity scope.

### v18.9.2 — TradeInsight Settings / API-key UX
Reuse Data Provider Settings and canonical local secret owner. Masked Save/Test/Clear, truthful state, scroll/focus preservation, optional environment override only as developer/runtime fallback. Where the same Settings capability exists on supported native clients, keep behavior/semantics equivalent. No routing/capability expansion.

### v18.9.3 — Coverage-Aware Smart Provider Router Core
Evolve Smart Provider Router v2 from first-success to requirement/coverage-aware fulfillment. Memory + persisted canonical state first; validate freshness/schema/provenance/rights; calculate exact residual gap; rank eligible providers; acquire only the gap; merge/reconcile with provenance; re-evaluate coverage; stop deterministically. Validation lifecycle remains separate from serving role.

### v18.9.4 — Canonical Company / Instrument Identity
One identity owner and presentation contract across Day/Swing/Long, Research, Discovery and Add Symbol. Symbol-only fallback remains valid. No provider search admission yet.

### v18.9.5 — Market Data Modes + Capability Diagnostics
Behavior/quality-oriented Adaptive modes and diagnostics only. Show actual source contribution, freshness, coverage and fallback/backfill state. No provider-brand mode or duplicate Market Mode owner.

### v18.9.6 — Provider Observability / Adaptive Telemetry Foundation
**Reordered earlier by architecture audit.** Instrument before further SHADOW provider expansion.

Measure at minimum:
- requested vs fulfilled coverage and residual gaps;
- DB/cache/persistence reuse and calls avoided;
- provider contribution/usefulness, latency, errors, rate-limit/backpressure and freshness failures;
- disagreement/corroboration/source independence;
- shared-demand/coalescing/fan-out efficiency;
- CPU/memory/goroutine/worker/DB/network pressure;
- protected-session provider/runtime headroom;
- SHADOW usefulness/outcome evidence for promotion/demotion.

No deterministic Day/Swing/Long formula change.

### v18.9.7 — TradeInsight SEC Form 4 Enrichment
Contract-validated SHADOW-first enrichment/corroboration through existing SEC/Ownership. Direct SEC/EDGAR remains authoritative. Source-family de-duplication, rights/freshness/provenance and optional-provider non-degradation required. Promotion consumes v18.9.6 telemetry.

### v18.9.8 — TradeInsight Symbol / Company Search
Contract-validated search through canonical symbol validation/company identity as fallback/corroboration. U.S.-equity boundary remains final; GLD/SLV/USO exceptions preserved.

### v18.9.9 — TradeInsight Movers / Ranking Evidence
Contract-validated market-wide mover/ranking evidence enters Opportunity Radar only through existing scanner/ranker as SHADOW candidate evidence. No provider-specific ranking engine.

### v18.9.10 — Remaining Useful TradeInsight Capability Admission
Re-audit executable official REST/SDK/MCP surface and configured entitlement. Every useful capability gets explicit consumer, canonical owner, lifecycle, serving role, freshness, retention, rights, rate/cost and Market Mode disposition. No invented endpoints or hidden production Python/MCP dependency.

### v18.9.11 — Session-Aware Data Readiness Maintenance
One canonical light-overnight/heavy-weekend coordinator using existing persistence/cache, Smart Provider Router v2, freshness, provider budgets, telemetry, workload priority and U.S. session/calendar owners.

Protected pre-market, regular and after-hours workloads retain first claim on provider quota/headroom, CPU, memory, DB, network and workers. Maintenance is gap/value-driven, bounded, preemptible/checkpointed/resumable and never blind-refetches.

Machine contract: `adaptive-governance/PERSISTENCE_REUSE_AND_OFF_HOURS_DATA_READINESS_CONTRACT.md`.

### v18.9.12 — v18.9.x Professional Closure
No feature scope. Audit #57/#64 regressions, deterministic Day/Swing/Long equivalence, identity, Market Modes, DB-first residual-gap fulfillment, provider failure/recovery, SHADOW evidence, calls avoided, maintenance, protected-session capacity, macOS Apple Silicon + Windows x64 packages and Adaptive Intelligence Scorecard.

Exit requires zero unexplained carry-forward, zero orphan useful provider capability and zero duplicate router/freshness/persistence/session-scheduler/SEC/symbol/Market-Mode owner.

## 5. v19 build train — Professional Data Infrastructure + Hosted Account Platform

**Entry:** `v18.9.12` PASS.

Architecture target:

`macOS/Windows SQLite edge + Web client -> authenticated DE.PULSE hosted API/service -> PostgreSQL shared authority`

Normal commercial users are zero-key: provider credentials remain server-side only.

### Hosted serving decision order

`tenant/account -> session/device -> RBAC/capability -> DE.PULSE product entitlement -> provider legal/data rights -> privacy/data-class/retention policy -> canonical cache/persistence -> residual provider need -> Smart Provider Router v2 -> authorized projection/fan-out`

Provider rights, product entitlement, RBAC and privacy/data-governance policy remain distinct controls; no generic `entitled=true` shortcut is permitted.

### v19.0.x — Governance, Control Plane & Data Foundation

#### v19.0.0 — Provider Capability / Legal Rights Registry
Machine-readable provider/capability contract for lifecycle/serving role, commercial/multi-user use, server proxying, cache/retention, redistribution/display, derived/AI use, attribution, environment limits, concurrency/user limits and expiry. Unknown/expired rights fail closed.

#### v19.0.1 — Hosted Tenant / Identity / Device / Session Control Plane
Tenant context becomes first-class in identity and every hosted request. Canonical `SUPER_OWNER/OWNER/ADMIN/USER/DEMO`, account/user/device lifecycle, registration/revocation, session expiry/refresh, privileged re-authentication/MFA-class controls where applicable, audit and API/SSE/native/web consistency.

#### v19.0.2 — DE.PULSE Product Entitlement / Metering Policy
Billing-provider-agnostic product-plan truth: plan/status, feature/capability entitlement, quota dimensions, grace/suspension/disabled behavior and metering keys. Separate from RBAC and upstream provider rights. External checkout/invoicing may be deferred, but entitlement enforcement cannot be deferred past external multi-user activation.

#### v19.0.3 — Account Data Governance / Privacy Lifecycle
Define/enforce data inventory/classification, purpose/minimization, retention/deletion, review/export, deactivation/deletion, device retirement, SQLite/PostgreSQL/sync/cache/backup lifecycle, audit/log retention, operator/support access and data-residency disposition before shared schema/retention freezes.

#### v19.0.4 — Hosted Environment / IaC / Service Trust Foundation
Isolated dev/test/stage/prod, environment-specific identities/secrets/config, least-privilege service identities, network/TLS boundaries, IaC or equivalent reproducible desired state, drift detection, change review, reproducible deployment/rollback and environment observability. No production snowflake.

#### v19.0.5 — PostgreSQL Tenancy / Schema / Pool / Recovery Foundation
Shared authority only: tenant isolation, schema/migrations, indexes/pools/capacity, encrypted backup, HA/failover, PITR, restore drills and frozen RPO/RTO. Consume v19.0.3 privacy/retention policy. No broad sync activation.

#### v19.0.6 — Managed Secrets / KMS Lifecycle
Server-side provider/service-secret ownership, environment separation, least privilege, versioned refs, rotation/rollback/compromise revoke, redaction and audit. Commercial clients contain no platform provider secret.

#### v19.0.7 — Software Supply-Chain / Artifact & Dependency Assurance
Dependency/component inventory, direct/transitive visibility, vulnerability scanning, SBOM where applicable, source/license policy, approved sources, reproducible builds, artifact signing/attestation or equivalent integrity proof, provenance and vulnerable-component response.

#### v19.0.8 — Provider Quality / Cost / Coverage / SLO Scorecards
Turn v18.9 telemetry into capability/provider SLO and cost/usefulness evidence: freshness, completeness, latency, reliability, rate pressure, contribution, calls avoided, fallback quality, maintenance value and tenant-aware health.

#### v19.0.9 — Data Reconciliation / Revision / Point-in-Time Quality
Source independence, disagreement/reconciliation, corporate-action/adjustment correctness, historical gaps, revision preservation, observed/effective timestamps and provenance.

**v19.0 exit:** rights, tenant identity, product entitlement, privacy/data lifecycle, environment/IaC, PostgreSQL recovery, secrets, supply-chain provenance, SLOs and data quality must be executable/evidenced before shared hosted activation.

### v19.1.x — Hosted Data Plane + Cross-Platform Account/Sync Lockstep

#### v19.1.0 — Authenticated Hosted Provider Gateway
Expose existing Smart Provider Router v2 through a versioned hosted API boundary. Reuse canonical freshness/cache/persistence/state, resolve credentials server-side, add bounded rate/backpressure/circuit/kill-switch/degraded behavior and maintain API inventory/version/deprecation ownership. No second provider stack.

#### v19.1.1 — Unified Serving Policy + Live Fan-Out Isolation
Enforce tenant/RBAC/product-entitlement/provider-right/privacy checks across router/cache/persistence/REST/WebSocket/SSE. Existing multi-feed allocator remains upstream subscription owner. Prevent cross-tenant, tier, right or privacy leakage.

#### v19.1.2 — Sync Protocol Foundation
Stable IDs, schema/domain versions, snapshot/high-watermark bootstrap, SQLite atomic outbox for native clients, authenticated idempotent push, authoritative server sequence/change log, incremental pull, checkpoints, tombstones, retention/compaction, stale-device re-bootstrap and mixed-version negotiation. Web consumes the same canonical server state without becoming a second authority.

#### v19.1.3 — Cross-Platform Account / Session Client Foundation
**One shared capability across Mac + Windows + Web in the same release responsibility.** Implement login/session/device/role/capability/product-entitlement state and user-visible account/session behavior against the v19.0 control plane. Native secure-storage and browser session mechanisms may differ, but semantics and authorization outcomes must match. No platform is considered delivered separately.

#### v19.1.4 — Cross-Platform Preferences + Watchlists
**Mac + Windows + Web together.** Prove the same account-scoped preferences/watchlist truth, conflicts, reconnect behavior, user switching, device revocation, account deletion and resulting UI state. Native clients use local SQLite/offline semantics where applicable; Web uses canonical hosted state. Cross-platform convergence/equivalence is mandatory.

#### v19.1.5 — Cross-Platform Desks / Workspaces
**Mac + Windows + Web together.** Versioned membership/config/delete/history semantics through the same transport/state owners. No platform-specific desk engine. Same user-visible workspace meaning and authorization across all required clients.

**v19.1 exit:** gateway, serving policy, sync/account foundations and the first portable user domains must pass cross-platform lockstep. There is no macOS product pilot and no Windows/Web catch-up phase.

### v19.2.x — Cross-Platform Shared Product Expansion + Hosted Assurance

#### v19.2.0 — Cross-Platform Research / Durable State Portability
**Mac + Windows + Web together.** Portable lawful/provenance-bound research/state, same canonical Research structure, account state, authorization and retention semantics. Live market truth remains freshness/provider/state-owned.

#### v19.2.1 — Cross-Platform Market Intelligence / Discovery / Market Modes
**Mac + Windows + Web together.** Same canonical Market Intelligence, Opportunity Radar/Discovery outcomes, company identity, Market Modes, readiness/explanation semantics and source/freshness truth. UI may adapt responsively; intelligence meaning may not fork.

#### v19.2.2 — Cross-Platform Settings / RBAC / Product-Entitlement UX
**Mac + Windows + Web together.** Same role-aware navigation/visibility, direct-route/API denial, product-plan capability states, account/device controls and settings truth. Owner/Admin/User/Demo composition follows canonical role/capability contracts; UI hiding never substitutes for authorization.

#### v19.2.3 — Tenant-Aware Metering / Cost / Usage Observability
Per tenant/account/user/device/capability attribution, plan/quota consumption, cache/call avoidance, stream usage, provider cost where known, tenant health, throttling and anomalous consumption. Metrics distinguish client platform without creating per-platform business truth.

#### v19.2.4 — Multi-User Security / Abuse / Capacity Hardening
Object/function authorization negatives, sensitive-flow abuse protection, rate limits, fairness/noisy-neighbor isolation, queue/circuit/load shedding, edge protection, DB/pool limits, environment-boundary tests and protected-session capacity. Exercise mixed Mac/Windows/Web concurrent use.

#### v19.2.5 — Cross-Platform Hosted Sync / Gateway Assurance Closure (#66)
No feature scope. Full adversarial/failure/recovery matrix plus explicit **cross-platform equivalence matrix**. #66 closes only when Mac + Windows + Web all prove required account, state, research, intelligence, settings/authorization and zero-key behavior against the same canonical backend, with no unresolved material parity debt.

### v19.3.x — Professional Point-in-Time Evidence Substrate
- **v19.3.0 — Institutional / 13F Evidence Infrastructure.** Direct SEC truth, identity/mapping, amendments, filing lag, point-in-time holdings and outcome lineage.
- **v19.3.1 — Two-Sided Long / Short Thesis Evidence Substrate.** Point-in-time plans/theses, first-event ordering, side-aware outcomes and lawful short/crowding evidence with explicit UNKNOWN.
- **v19.3.2 — AODR Candidate / Ranking / Outcome Lineage.** Candidate/rank/reason transitions, shared-ranking efficiency, diversity/correlation metadata and surfaced-vs-missed outcomes.

Any user-facing v19.3 capability exposed before v19 closure follows the Cross-Platform Lockstep Product Contract.

### v19.4.x — Reliability, Economics & v20 Readiness
- **v19.4.0 — ADR-GDI Professional Reliability / Capacity.** SLO/error-budget/degradation history, provider/DB/runtime scorecards, restart/warm-start, query/index/pool tuning, load shedding, reserve sizing, maintenance economics and hosted runbooks.
- **v19.4.1 — Specialized / Paid Provider Gap Evaluation.** Only measured v19 gaps justify provider change; same router/rights/persistence/session contracts.
- **v19.4.2 — v20 Research-Readiness Dataset / Lineage Audit.** Prove point-in-time features/outcomes, provenance, rights/privacy compatibility, independence, leakage controls and reliability history.

### v19.5.0 — v19 Major Closure
No feature scope. Require #66 closure, zero material Mac/Windows/Web parity debt for shared capabilities, tenant isolation, product-entitlement/provider-right/privacy separation, data lifecycle, environment/IaC, supply-chain provenance, API compatibility, tested rollback/recovery, actual supported native/web artifacts/deployments, SLO/capacity truth and zero unresolved P0 architecture gap.

## 6. v20 build train — Adaptive Intelligence & Decision Research

**Entry:** `v19.5.0` PASS. Learn only from point-in-time, rights-valid, privacy-compatible, provenance-bound evidence/outcomes. Production influence remains `SHADOW -> VALIDATED -> APPROVED -> PRODUCTION`; No Execution remains permanent.

### v20.0.x — Adaptive Research Control & Governance
- **v20.0.0 — Adaptive Research Control Plane + Immutable Experiment Ledger.** Dataset/version lineage, feature/provenance snapshots, cohorts, reproducibility, leakage controls, retention/privacy compatibility and promotion/rollback evidence.
- **v20.0.1 — Model / Prompt Governance + Champion/Challenger.** Model/prompt identity, independent evaluation, explainability, drift, approval and rollback before broad adaptive rollout.
- **v20.0.2 — Historical Analogues + Regime-Conditioned Outcomes.**
- **v20.0.3 — Calibration / FP-FN / Miss / Contradiction / Drift.**

### v20.1.x — ASBI
- `v20.1.0` Behavioral Fingerprints + State Transitions.
- `v20.1.1` Scenarios / Probability Momentum / Calibration.

### v20.2.x — Institutional + TDTI
- `v20.2.0` Adaptive Institutional / 13F Intelligence.
- `v20.2.1` TDTI Competing Long / Short / No Reliable Edge.
- `v20.2.2` TDTI Two-Sided Trade-Plan Validation. No Execution.

### v20.3.x — AODR
- `v20.3.0` Adaptive Shared Opportunity Ranking.
- `v20.3.1` Diversity / Opportunity Cost / Personalized Relevance after shared truth.

### v20.4.x — Adaptive Operations
- `v20.4.0` ADR-GDI Adaptive Optimization under SHADOW/Champion-Challenger; no self-promotion.

### v20.5.0 — v20 Professional Closure
No feature scope. Calibrated utility, drift/abstention, deterministic-boundary protection, privacy/security/data rights, rollback, reproducibility, actual supported artifacts, **cross-platform parity for every shared adaptive user-facing capability**, zero silent self-modification and No Execution.

## 7. Cross-version dependency contract

`v18.9 stabilize/instrument/validate -> v19.0 hosted foundations -> v19.1 shared account/sync capabilities in lockstep -> v19.2 shared product capabilities in lockstep + assurance -> v19.3 evidence substrate -> v19.4 reliability/readiness -> v19.5 closure -> v20 governed adaptive research in lockstep`

No later phase may compensate for an unresolved earlier foundation or platform parity gap by creating a parallel owner or hiding uncertainty behind model confidence.

## 8. G0-G16 cross-platform checkpoints

- **G1:** platform applicability matrix REQUIRED/N/A is frozen for every scope.
- **G2:** canonical domain owner and platform-adapter boundaries mapped; no client owns duplicate business truth.
- **G3:** one API/domain/state contract plus platform-specific adapter contracts, compatibility matrix, UI semantics and equivalence tests.
- **G4:** implementation exists for all REQUIRED platforms before shared capability Development Exit; diagnostic-only single-platform work is not delivery.
- **G5/G6:** affected-platform deterministic and cross-platform integration tests.
- **G7:** authorization/data/privacy/provider-right outcomes equivalent across clients.
- **G8:** mixed-client load/capacity/recovery where applicable.
- **G9:** role-aware Mac/Windows/Web UX equivalence; responsive layout differences allowed, meaning/functionality drift is not.
- **G10:** unresolved material parity debt blocks freeze.
- **G11/G12:** immutable RC/full certification includes required platform matrix.
- **G13/G14:** actual macOS/Windows packages and Web deployment/runtime evidence where REQUIRED.
- **G15:** no GA/Delivered state for a shared capability until all REQUIRED platforms pass. Canary/bounded cohorts may be used, but platform lag is not a normal canary strategy.
- **G16:** implementation-miss review includes parity drift and exact next handoff.

All existing security/privacy/IaC/supply-chain/SLO requirements remain inside their appropriate G0-G16 gates; no G17+ is introduced.

## 9. Protected-session resource contract

Pre-market, regular market and after-hours remain Tier-0 decision-support sessions. Live/current workloads always outrank maintenance/background sync. Provider quota/headroom, network, CPU, memory, DB/pool and workers retain explicit reserve. Maintenance/sync uses bounded surplus capacity and yields/preempts before current truth is materially degraded.

## 10. CI/release efficiency

For each patch: coherent code+test batch before PR; one product PR; exact-head FAST once for coherent candidate; Qualified once when Ready; one canonical G11-G16 release when release-capable. Cross-platform scope should share fixtures/contracts/evidence where trustworthy rather than running three independent business-logic pipelines. Classify before rerun and preserve equivalent unaffected evidence.

## 11. Protected boundaries

Smart Provider Router v2 sole routing authority; canonical freshness/recovery sole freshness owner; existing multi-feed allocator sole subscription owner; BroadSnapshotBroker canonical reuse owner; canonical persistence/cache/state owners reused; canonical U.S. session calendar reused; direct SEC/EDGAR authoritative; canonical tenant/identity/role/capability truth; deterministic Day/Swing/Long protected; U.S. Equities Processing; GLD/SLV/USO actionable exceptions; No Execution; G0-G16 only.

## Exactly one next action

Execute G0 for issue #64 / `v18.9.1` using complete macOS crash evidence or deterministic reproduction, then freeze its narrow G1. Do not start `v18.9.2` or any v19 implementation branch until `v18.9.1` is truthfully closed or the crash is proven external/non-product.