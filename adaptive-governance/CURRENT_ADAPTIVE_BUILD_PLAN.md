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

`stabilize -> establish canonical owners -> instrument -> validate in SHADOW -> expand -> operationalize -> close -> establish hosted control plane -> activate hosted data plane -> synchronize -> prove parity -> harden -> build evidence substrate -> learn`

Permanent rules:
- one primary responsibility per patch;
- `REUSE -> CONSOLIDATE -> REFACTOR -> DELETE/REPLACE -> ADD` before new machinery;
- canonical owners are extended, never forked for a provider/client/feature;
- observability required to judge a capability lands before broad capability admission;
- tenant identity, security, provider rights, product entitlement, persistence recovery and secret management precede commercial hosted activation;
- migrations use backward-compatible expand/contract where practical and always have rollback/roll-forward/recovery disposition before activation;
- high-risk externally visible behavior uses bounded feature/capability flags, kill switches, canary/progressive activation or equivalent controls where useful;
- governance closure follows `Governed -> Implemented -> Enforced -> Evidenced -> Delivered -> Learned`;
- known misses are fixed in-scope or durably assigned to a named later release;
- one product branch + one product PR per patch; CI reruns only the smallest trustworthy affected set;
- G0-G16 is the only top-level release model.

## 2. Version alignment

- **Major (`v19`, `v20`)** = strategic maturity generation.
- **Minor band (`v19.0.x`, `v19.1.x`, etc.)** = coherent dependency phase with explicit entry/exit criteria.
- **Patch (`x.y.z`)** = one primary independently certifiable responsibility.
- Future numbers are planning reservations only until G1.
- If a future packet is too broad, split it in the same minor band and shift later unstarted reservations.
- Shipped versions are immutable.
- Corrective/security work can preempt the train; dependencies are reconciled before planned work resumes.

## 3. v18.9.x build train — trustworthy runtime/data plane

### v18.9.1 — Runtime Reliability
Scope owner #64. Diagnose/fix only the real packaged macOS Apple Silicon SIGABRT from evidence/reproduction. Preserve SQLite/user state/API keys. Require warm-state/relaunch lifecycle regression and actual packaged runtime proof. No provider/router/Market Mode/identity scope.

### v18.9.2 — TradeInsight Settings / API-key UX
Reuse Data Provider Settings and canonical local secret owner. Masked Save/Test/Clear, truthful state, scroll/focus preservation, optional environment override only as developer/runtime fallback. No routing/capability expansion.

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

## 4. v19 build train — Professional Data Infrastructure + Hosted Account Platform

**Entry:** `v18.9.12` PASS.

Architecture target:

`macOS/Windows SQLite edge -> authenticated DE.PULSE hosted API/service -> PostgreSQL shared authority`

Hosted web uses the same API. Normal commercial users are zero-key: provider credentials remain server-side only.

### Hosted serving decision order

`tenant/account -> session/device -> RBAC/capability -> DE.PULSE product entitlement -> provider legal/data rights -> freshness/data-class policy -> canonical cache/persistence -> residual provider need -> Smart Provider Router v2 -> authorized projection/fan-out`

Provider rights, product entitlement and RBAC are never represented as one generic `entitled=true` flag.

### v19.0.x — Governance, Control Plane & Data Foundation

#### v19.0.0 — Provider Capability / Legal Rights Registry
Machine-readable provider/capability contract for lifecycle/serving role, commercial/multi-user use, server proxying, cache/retention, redistribution/display, derived/AI use, attribution, environment limits, concurrency/user limits and expiry. Unknown/expired rights fail closed.

#### v19.0.1 — Hosted Tenant / Identity / Device / Session Control Plane
Tenant context becomes first-class in identity and every hosted request. Canonical `SUPER_OWNER/OWNER/ADMIN/USER/DEMO`, account/user/device lifecycle, registration/revocation, session expiry/refresh, privileged re-authentication/MFA-class controls where applicable, audit and API/SSE/native/web consistency.

#### v19.0.2 — DE.PULSE Product Entitlement / Metering Policy
Billing-provider-agnostic product-plan truth: plan/status, feature/capability entitlement, quota dimensions, grace/suspension/disabled behavior and metering keys. Separate from RBAC and upstream provider rights. External payment checkout/invoicing may be deferred, but entitlement enforcement cannot be deferred past external multi-user activation.

#### v19.0.3 — PostgreSQL Tenancy / Schema / Pool / Recovery Foundation
Shared authority only: tenant/account isolation, schema/migration ownership, connection pools/capacity, indexes, encrypted backup, HA/failover, PITR, restore drills and frozen RPO/RTO. Use expand/contract migration strategy where practical. No broad sync activation.

#### v19.0.4 — Managed Secrets / KMS Lifecycle
Server-side provider-secret ownership, environment separation, least privilege, versioned references, rotation, rollback, compromise revoke, redaction and auditable administration. Commercial clients contain no platform provider secret.

#### v19.0.5 — Provider Quality / Cost / Coverage / SLO Scorecards
Turn v18.9 telemetry into capability/provider SLO and cost/usefulness evidence: freshness, completeness, latency, reliability, rate pressure, contribution, calls avoided, fallback quality, maintenance value and tenant-aware health signals.

#### v19.0.6 — Data Reconciliation / Revision / Point-in-Time Quality
Source independence, disagreement/reconciliation, corporate-action/adjustment correctness, historical depth/gaps, revision preservation, observed/effective timestamps and provenance.

**v19.0 exit:** tenant identity, provider rights, product entitlement, PostgreSQL recovery, secrets, SLOs and data quality are executable/evidenced before shared hosted activation.

### v19.1.x — Zero-Key Provider Data Plane + Native Sync Foundation

#### v19.1.0 — Authenticated Hosted Provider Gateway
Expose the existing Smart Provider Router v2 through a versioned hosted API boundary. Reuse canonical freshness/cache/persistence/state, resolve provider credentials server-side only, add bounded rate/backpressure/circuit/kill-switch/degraded behavior, and maintain API inventory/version/deprecation ownership. No second provider stack.

#### v19.1.1 — Unified Serving Policy + Live Fan-Out Isolation
Enforce tenant/RBAC/product-entitlement/provider-right checks across router/cache/persistence/REST/WebSocket/SSE. Existing multi-feed allocator remains upstream subscription owner. Prevent cross-tenant, premium-tier and legally restricted data leakage.

#### v19.1.2 — Sync Protocol Foundation
Stable IDs, schema/domain capability versions, snapshot/high-watermark bootstrap, SQLite atomic outbox, authenticated idempotent push, authoritative server sequence/change log, incremental pull, durable per-device checkpoint, tombstones, retention/compaction, stale-device re-bootstrap and mixed-version negotiation. Client wall clock never decides conflicts alone. No raw SQLite/PostgreSQL replication.

#### v19.1.3 — macOS Preferences + Watchlist Pilot
First real sync-domain activation. Prove portable preferences/ticker convergence, offline writes, restart/reconnect, conflicts, user switching, local account isolation, lost/revoked-device behavior and no destructive reset of `PersonalMarketTerminal`.

#### v19.1.4 — Desks / Workspaces Sync
Versioned membership/configuration/delete/conflict/history semantics using the same transport and canonical state owners. No desk-specific sync engine.

**v19.1 exit:** zero-key gateway, unified serving policy and proven native sync transport/macOS pilot without protected-session degradation.

### v19.2.x — Cross-Platform Account Parity & Hosted Assurance

#### v19.2.0 — Windows x64 Sync Parity
Same protocol/state/security semantics as macOS; no Windows-specific account truth/provider path.

#### v19.2.1 — Hosted Web Account Parity
Same authenticated APIs, session/capability/product-entitlement truth and PostgreSQL account state. Browser cache is non-authoritative. Direct-route/API authorization and web session/cookie/security-header controls are evidenced.

#### v19.2.2 — Rights-Aware Research / Durable State Portability
Sync only product-owned, lawful, entitled, provenance-bound durable research/state. Preserve as-observed/revision history where required. Live market truth remains canonical freshness/provider/state-owned.

#### v19.2.3 — Tenant-Aware Metering / Cost / Usage Observability
Per tenant/account/user/device/capability attribution, plan/quota consumption, cache/call avoidance, stream usage, provider cost where known, tenant health, throttling and anomalous consumption visibility. This is the operational evidence for product-plan policy and noisy-neighbor controls.

#### v19.2.4 — Multi-User Security / Abuse / Capacity Hardening
Object/function-level authorization negatives, sensitive-flow abuse protection, rate limits, fairness, noisy-neighbor isolation, queue/circuit/load shedding, edge protection, DB/pool limits and protected-session capacity.

#### v19.2.5 — Hosted Sync / Gateway Assurance Closure (#66)
No feature scope. Failure/adversarial matrix: cross-account denial, role/session/device revocation, duplicate replay, network loss during apply, checkpoint expiry/re-bootstrap, provider-right downgrade/expiry, product-plan downgrade/suspension, secret rotation rollback, DB failover/PITR/restore, mixed-version clients, queue pressure and protected-session load. #66 closes only on executable Mac/Windows/web parity plus zero-key/rights/product-entitlement/sync/DR proof.

### v19.3.x — Professional Point-in-Time Evidence Substrate

#### v19.3.0 — Institutional / 13F Evidence Infrastructure
Direct SEC truth, manager/CIK/security mapping, amendments/restatements, filing-lag truth, point-in-time holdings, storage/indexing and outcome lineage.

#### v19.3.1 — Two-Sided Long / Short Thesis Evidence Substrate
Point-in-time thesis/plan snapshots, target/invalidation ordering, side-aware MFE/MAE and lawful/trustworthy short-interest/crowding/borrow/SSR context; explicit UNKNOWN where evidence is insufficient.

#### v19.3.2 — AODR Candidate / Ranking / Outcome Lineage
My Market vs Global truth, candidate/rank/reason snapshots, NOW/WATCH/PASS/ABSTAIN transitions, shared-ranking efficiency, diversity/correlation metadata, surfaced-vs-missed outcomes and recommendation usefulness.

### v19.4.x — Reliability, Economics & v20 Readiness

#### v19.4.0 — ADR-GDI Professional Reliability / Capacity
Capability SLO/error-budget/degradation history, provider/DB/runtime scorecards, restart/warm-start, query/index/pool tuning, load shedding, operating limits, reserve sizing, maintenance/preemption economics and hosted incident/runbook readiness.

#### v19.4.1 — Specialized / Paid Provider Gap Evaluation
Only after measured v19 evidence proves material capability/quality/rights/cost gaps. Any provider plugs into the same router/rights/persistence/session-priority contracts.

#### v19.4.2 — v20 Research-Readiness Dataset / Lineage Audit
No model scope. Prove point-in-time features/outcomes, provenance, rights, independence, synchronized-evidence safety, leakage controls and reliability history are adequate for adaptive research.

### v19.5.0 — v19 Major Closure
No feature scope. Whole-system Principal Engineer/security/data/operational/commercial-readiness audit. Require #66 closure, tenant isolation, product-entitlement/provider-right separation, API inventory/version compatibility, tested rollback/recovery, actual supported runtime/package evidence, SLO/capacity truth and zero unresolved P0 architecture gap.

## 5. v20 build train — Adaptive Intelligence & Decision Research

**Entry:** `v19.5.0` PASS. Learn only from point-in-time, rights-valid, provenance-bound evidence/outcomes. Production influence remains `SHADOW -> VALIDATED -> APPROVED -> PRODUCTION`; No Execution remains permanent.

### v20.0.x — Adaptive Research Control & Governance
- **v20.0.0 — Adaptive Research Control Plane + Immutable Experiment Ledger.** Dataset/version lineage, feature/provenance snapshots, cohorts, reproducibility, leakage controls, promotion/rollback evidence.
- **v20.0.1 — Model / Prompt Governance + Champion/Challenger.** Model/prompt identity, independent evaluation, explainability, drift, approval and rollback **before broad adaptive rollout**.
- **v20.0.2 — Historical Analogues + Regime-Conditioned Outcomes.** Point-in-time analogue/outcome research.
- **v20.0.3 — Calibration / FP-FN / Miss / Contradiction / Drift.** Calibration, missed-opportunity and abstention controls.

### v20.1.x — ASBI
- **v20.1.0 — Behavioral Fingerprints + State Transitions.**
- **v20.1.1 — Scenarios / Probability Momentum / Calibration.**

### v20.2.x — Institutional + TDTI
- **v20.2.0 — Adaptive Institutional / 13F Intelligence.**
- **v20.2.1 — TDTI Competing Long / Short / No Reliable Edge.**
- **v20.2.2 — TDTI Two-Sided Trade-Plan Validation.** No Execution.

### v20.3.x — AODR
- **v20.3.0 — Adaptive Shared Opportunity Ranking.**
- **v20.3.1 — Diversity / Opportunity Cost / Personalized Relevance.** Shared truth first; user relevance second.

### v20.4.x — Adaptive Operations
- **v20.4.0 — ADR-GDI Adaptive Optimization.** SHADOW/Champion-Challenger provider/recovery/workload/maintenance/reserve learning; no self-promotion.

### v20.5.0 — v20 Professional Closure
No feature scope. Calibrated utility, drift/abstention, deterministic-boundary protection, privacy/security/data rights, rollback, reproducibility, actual supported artifacts, zero silent self-modification and No Execution.

## 6. Cross-version dependency contract

`v18.9 stabilize/instrument/validate -> v19.0 control plane -> v19.1 hosted data plane/sync -> v19.2 parity/assurance -> v19.3 evidence substrate -> v19.4 reliability/readiness -> v19.5 closure -> v20 governed adaptive research`

No later phase may compensate for an unresolved earlier foundation by creating a parallel owner or hiding uncertainty behind model confidence.

## 7. G0-G16 industry-strength checkpoints

Hosted/security/data/adaptive patches add responsibilities inside existing gates, never G17+:
- **G2:** canonical-owner map, tenant/data classification, threat-model scope and ADR for material irreversible decisions;
- **G3:** API/schema/protocol contracts, API inventory/deprecation, migration/rollback, compatibility matrix, SLO/error-budget targets, observability, failure-injection and rollout plan;
- **G4:** implementation + unit/contract tests + feature/kill-switch controls where useful;
- **G7:** negative authorization/tenant isolation, product-entitlement, provider-right, secret/redaction and adaptive-governance proof;
- **G8:** load/soak/capacity, fairness/noisy-neighbor, backpressure, chaos/failure recovery, DB/pool/provider limits and protected-session proof;
- **G9:** role-aware Mac/Windows/web UX and direct-route/API denial consistency;
- **G10:** production-readiness reconciliation; unresolved P0 threat/rights/recovery/compatibility/duplicate-owner issue blocks freeze;
- **G11/G12:** immutable RC/full certification;
- **G13/G14:** native package/provenance/runtime proof where applicable;
- **G15:** controlled/canary/progressive promotion evidence for hosted-risk changes where applicable;
- **G16:** implementation-miss review, incidents/metrics, cleanup and exact next handoff.

## 8. Protected-session resource contract

Pre-market, regular market and after-hours remain Tier-0 decision-support sessions. Live/current workloads always outrank maintenance/background sync. Provider quota/headroom, network, CPU, memory, DB/pool and workers retain explicit reserve. Maintenance/sync uses bounded surplus capacity and yields/preempts before current truth is materially degraded.

## 9. CI/release efficiency

For each patch: coherent code+test batch before PR; one product PR; exact-head FAST once for coherent candidate; Qualified once when Ready; one canonical G11-G16 release when release-capable. Classify before rerun, reuse unchanged evidence when fingerprints/dependencies remain equivalent and never create retry/certification branch families.

## 10. Protected boundaries

Smart Provider Router v2 sole routing authority; canonical freshness/recovery sole freshness owner; existing multi-feed allocator sole subscription owner; BroadSnapshotBroker canonical reuse owner; canonical persistence/cache/state owners reused; canonical U.S. session calendar reused; direct SEC/EDGAR authoritative; canonical tenant/identity/role/capability truth; deterministic Day/Swing/Long protected; U.S. Equities Processing; GLD/SLV/USO actionable exceptions; No Execution; G0-G16 only.

## Exactly one next action

Execute G0 for issue #64 / `v18.9.1` using complete macOS crash evidence or deterministic reproduction, then freeze its narrow G1. Do not start `v18.9.2` or any v19 implementation branch until `v18.9.1` is truthfully closed or the crash is proven external/non-product.