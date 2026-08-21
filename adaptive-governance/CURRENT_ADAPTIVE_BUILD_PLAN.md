# DE.PULSE — Current Adaptive Build Plan

**Immutable Stable:** `v18.9.0-stable`  
**Stable candidate:** `9ea81cddae4875ae15d3719ca028519a36c597b6`  
**Stable fingerprint:** `a8719090c341c874dbd1279cc31ad98e84075d5701c46a800bf951340780ecb9`  
**Active product development branch:** none  
**Active product PR:** none  
**Governance alignment PR:** #67 (draft)  
**Master corrective program:** #65 / `ADAPT-PROVIDER-INTELLIGENCE-010`  
**Immediate next product patch:** `v18.9.1` / #64 / `ADAPT-RUNTIME-CRASH-001`.

## 1. Build philosophy — industry-aligned small release trains

DE.PULSE uses one permanent G0–G16 release model and prefers small, dependency-ordered, independently certifiable patches over heavy multi-domain builds.

Permanent engineering sequence:

`stabilize -> instrument -> validate -> expand -> harden -> close -> platform-foundation -> scale -> learn`

Rules:
- one primary responsibility per patch; only tightly coupled support work may accompany it;
- `REUSE -> CONSOLIDATE -> REFACTOR -> DELETE/REPLACE -> ADD` before creating new machinery;
- canonical owners are extended rather than duplicated;
- observability needed to judge a capability must exist before broad production admission of that capability;
- security, rights, identity and recovery foundations precede multi-user data-plane activation;
- migrations are backward-compatible/expand-contract where practical and have rollback/recovery plans before activation;
- externally visible or high-risk behavior uses feature/capability flags, kill switches or equivalent bounded rollout controls where useful;
- no governance requirement is complete at documentation-only state: `Governed -> Implemented -> Enforced -> Evidenced -> Delivered -> Learned`;
- every known miss is fixed in-scope or assigned durably to a named later release; no chat-only carry-forward;
- one product branch + one PR per patch; CI runs the smallest trustworthy affected set and reuses equivalent evidence;
- G0–G16 remains the only top-level gate model.

## 2. Version alignment contract

DE.PULSE version numbers communicate dependency phases without weakening small-patch discipline.

- **Major (`v19`, `v20`)** = strategic architecture generation / product maturity stage.
- **Minor band (`v19.0.x`, `v19.1.x`, ...)** = a coherent platform phase with a clear entry/exit contract.
- **Patch (`x.y.z`)** = one primary independently certifiable responsibility.
- A phase-ending `.0` or closure release may be audit/assurance-only; it does not need feature scope.
- Version labels below are **planned reservations/alignment**, not immutable scope. Exact release identity freezes only at that release's G1.
- If G0/G1 proves a planned packet is too broad, split it inside the same minor band and shift later unstarted reservations. Never bundle unrelated work merely to preserve a number.
- Already shipped versions are immutable and never renumbered.
- Corrective/security hotfixes take priority over planned feature order; after the corrective closes, the dependency graph is re-reconciled before resuming.

## 3. v18.9.x — trustworthy acquisition, observability and operational closure

`v18.9.0-stable` is immutable. The remaining v18.9.x line is corrective/professional hardening under #65.

### v18.9.1 — Runtime Reliability
Scope owner #64. Diagnose/fix only the real packaged macOS Apple Silicon SIGABRT from evidence/reproduction. Preserve SQLite/user state/API keys. Require lifecycle/relaunch regression and actual packaged runtime proof. No provider/router/Market Mode/identity changes.

### v18.9.2 — TradeInsight Settings / API Key UX
Reuse the existing Data Provider Settings and canonical local secret owner. Masked Save/Test/Clear, truthful configured/connected/error/capability state, scroll/focus preservation, optional environment override only as developer/runtime fallback. No routing/provider-capability expansion.

### v18.9.3 — Coverage-Aware Smart Provider Router Core
Evolve Smart Provider Router v2 from first-success to requirement/coverage-aware fulfillment. Use memory + persisted canonical state first, validate freshness/schema/provenance/rights, calculate exact residual gap, rank eligible providers for the remaining need, acquire only the gap, merge/reconcile with provenance, re-evaluate coverage and stop deterministically. Validation lifecycle remains separate from serving role. No provider-specific router path.

### v18.9.4 — Canonical Company / Instrument Identity
One identity owner and presentation contract. Day/Swing/Long, Research, Discovery and Add Symbol reuse the same canonical symbol/company identity; symbol-only fallback remains valid. No provider search admission yet.

### v18.9.5 — Market Data Modes + Capability Diagnostics
Behavior/quality-oriented Adaptive modes and capability diagnostics only. Diagnostics expose actual source contribution, freshness, coverage and fallback/backfill state without provider-brand modes or a second Market Mode owner.

### v18.9.6 — Provider Observability / Adaptive Telemetry Foundation
**Reordered earlier by architecture audit.** Instrument before further SHADOW provider expansion so later admission decisions are evidence-based.

Measure at minimum:
- requested vs fulfilled coverage and residual gaps;
- DB/cache/persistence reuse and provider calls avoided;
- provider contribution/usefulness, latency, errors, rate-limit/backpressure and freshness failures;
- disagreement/corroboration/source-independence evidence;
- shared-demand/coalescing/fan-out efficiency where applicable;
- CPU, memory, goroutine/worker, DB and network pressure;
- protected-session provider/runtime headroom;
- SHADOW usefulness/outcome evidence needed for provider promotion/demotion.

This patch instruments/observes; it does not silently change deterministic Day/Swing/Long truth.

### v18.9.7 — TradeInsight SEC Form 4 Enrichment
Contract-validated SHADOW-first enrichment/corroboration through the existing SEC/Ownership model. Direct SEC/EDGAR remains authoritative. Source-family de-duplication, rights/freshness/provenance and optional-provider non-degradation required. Promotion uses v18.9.6 telemetry.

### v18.9.8 — TradeInsight Symbol / Company Search
Contract-validated ticker/company search through canonical symbol validation/company identity as fallback/corroboration. U.S.-equity boundary final; GLD/SLV/USO actionable exceptions preserved. No parallel registry.

### v18.9.9 — TradeInsight Movers / Ranking Evidence
Contract-validated market-wide mover/ranking evidence enters Opportunity Radar only through the existing scanner/ranker as SHADOW candidate evidence. No undocumented endpoint assumptions and no provider-specific ranking engine.

### v18.9.10 — Remaining Useful TradeInsight Capability Admission
Re-audit executable official REST/SDK/MCP surface and configured entitlement. Congress, daily adjusted/raw OHLCV, corporate actions, bounded history and any other useful capability receive explicit consumer, owner, lifecycle, serving role, freshness, retention, rights, rate/cost and Market Mode disposition. No Python/MCP production dependency unless separately governed; no invented intraday capability.

### v18.9.11 — Session-Aware Data Readiness Maintenance
Implement one canonical coordinator for light overnight + heavier bounded weekend/extended-closed maintenance using existing persistence/cache, Smart Provider Router v2, freshness, provider budgets, telemetry, workload-priority and canonical U.S. session/calendar owners.

Protected pre-market, regular-market and after-hours workloads retain first claim on provider quota/headroom, CPU, memory, DB, network and workers. Maintenance is gap/value-driven, bounded, preemptible/checkpointed/resumable, cannot blind-refetch and cannot create another calendar/router/cache/database/scheduler owner.

Acceptance includes overnight/weekend behavior, reserve sizing, preemption/drain/checkpoint/resume, missed-window bounded catch-up, restart de-duplication and no material decision-critical degradation.

Machine contract: `adaptive-governance/PERSISTENCE_REUSE_AND_OFF_HOURS_DATA_READINESS_CONTRACT.md`.

### v18.9.12 — v18.9.x Professional Closure Audit
No new feature scope. Deep implementation-miss and architecture audit across the entire line: #57/#64 regressions, deterministic Day/Swing/Long equivalence, identity, Market Modes, DB-first residual-gap fulfillment, provider failure/recovery, SHADOW evidence, calls avoided, maintenance, protected-session capacity, macOS Apple Silicon + Windows x64 packages and Adaptive Intelligence Scorecard.

Closure requires zero unexplained carry-forward, zero orphan useful provider capability and zero duplicate routing/freshness/persistence/session-scheduler/SEC/symbol/Market-Mode owner.

## 4. v19 — Professional Data Infrastructure + Hosted Account Platform

**Entry:** `v18.9.12` Major-Line Closure PASS. v19 consumes v18.9 coverage-aware routing, persistence-first reuse, telemetry, identity, Market Modes and session-aware maintenance. It does not recreate them.

Architecture target:

`macOS/Windows SQLite edge -> authenticated DE.PULSE hosted API/service -> PostgreSQL shared authority`

Hosted web uses the same service/API. Commercial normal users are **zero-key**: they authenticate only to DE.PULSE; platform provider credentials stay server-side in managed secrets/KMS. Issue #66 / `ADAPT-HOSTED-SYNC-001` is executable v19 scope.

### v19.0.x — Governance, Control Plane & Data Foundation

#### v19.0.0 — Provider Capability / Entitlement / Rights Registry
Machine-readable provider/capability truth for lifecycle, serving role, commercial/multi-user entitlement, server-side proxying, caching/retention, redistribution/display, derived/AI use, attribution, environment, limits and expiry. Runtime consumers fail closed when rights are unknown or expired.

#### v19.0.1 — Hosted Identity / Device / Session Control Plane
Canonical `SUPER_OWNER/OWNER/ADMIN/USER/DEMO` + capability authorization, account/user/device identity, registration/revocation, session expiry/refresh, privileged re-authentication/MFA-class assurance where applicable, account lifecycle and audit. Server-side truth applies consistently to API/SSE/native/web.

#### v19.0.2 — PostgreSQL Tenancy / Schema / Pool / Recovery Foundation
PostgreSQL shared-authority foundation only: tenancy/account isolation, schema/migration ownership, connection-pool/capacity boundaries, indexes, encrypted backup, HA/failover disposition, PITR, restore drills and frozen RPO/RTO behavior. Use backward-compatible expand/contract migrations and rollback/roll-forward plans. No broad sync domain activation yet.

#### v19.0.3 — Managed Secrets / KMS Lifecycle
Server-side provider-secret ownership, environment separation, least privilege, versioned references, rotation, rollback, compromise revoke, redaction and auditable administration. Commercial clients contain no platform provider secret. No Provider Gateway activation until this passes.

#### v19.0.4 — Provider Quality / Cost / Coverage / SLO Scorecards
Turn v18.9 telemetry into measured capability/provider SLOs and usefulness/cost evidence: freshness, completeness, latency, reliability, rate pressure, contribution, calls avoided, fallback quality and maintenance value. Define SLO/error-budget signals that later gateway/load-shedding decisions can consume.

#### v19.0.5 — Data Reconciliation / Revision / Point-in-Time Quality
Source independence, disagreement/reconciliation policy, corporate-action/adjustment correctness, historical depth/gaps, revision preservation, observed/effective timestamps and point-in-time provenance. No silent overwrite of history.

**v19.0.x exit:** rights, identity, PostgreSQL, secrets, SLOs and data-quality foundations are executable and evidenced. Only then may shared hosted provider/sync activation begin.

### v19.1.x — Zero-Key Provider Data Plane & Native Sync Foundation

#### v19.1.0 — Authenticated Hosted Provider Gateway
Wrap the existing Smart Provider Router v2 behind the hosted boundary. REST/snapshot/provider requests use server-side credentials only. Reuse canonical freshness/cache/persistence/state. Add bounded rate/backpressure/circuit behavior and kill-switch/degraded controls. Do not create another provider stack.

#### v19.1.1 — Machine-Enforced Rights / Entitlement + Live Fan-Out Isolation
Rights/entitlement gates actively control router/cache/persistence/REST/WebSocket/SSE reuse. Reuse the existing multi-feed subscription owner for upstream subscriptions and authorized downstream fan-out. Prevent premium/realtime/right-restricted evidence from leaking across accounts/plans.

#### v19.1.2 — Sync Protocol Foundation
Application-level typed sync only: stable IDs, schema/domain capability versions, new-device snapshot + high-watermark bootstrap, SQLite atomic outbox, authenticated idempotent push, authoritative server sequence/change log, incremental pull, per-device durable checkpoint, tombstones, retention/compaction, stale-device re-bootstrap and mixed-version negotiation. Client wall clock is never sole ordering authority. No raw DB replication.

#### v19.1.3 — macOS Preferences + Watchlist Pilot
First real sync-domain activation. Prove portable preference/ticker convergence, offline writes, restart/reconnect, conflicts, user switching, local account isolation, lost/revoked device behavior and no destructive reset of `PersonalMarketTerminal`.

#### v19.1.4 — Desks / Workspaces Sync
Versioned membership/configuration/delete/conflict/history semantics using the same sync transport and canonical state owners. No desk-specific sync engine.

**v19.1.x exit:** zero-key gateway, rights-safe live serving and a proven native sync transport/macOS pilot exist without compromising protected-session resources.

### v19.2.x — Cross-Platform Account Parity & Hosted Assurance

#### v19.2.0 — Windows x64 Sync Parity
Same protocol/state/security semantics as macOS; no Windows-specific account truth or provider path.

#### v19.2.1 — Hosted Web Account Parity
Same authenticated hosted APIs/session/capability truth and PostgreSQL-backed account state. Browser SQLite is not authoritative. Security headers/session/cookie controls and direct-route/API authorization are evidenced.

#### v19.2.2 — Rights-Aware Research / Durable State Portability
Sync only product-owned, lawful, entitled, provenance-bound durable research/state. Preserve as-observed/revision lineage where required. Live market truth remains canonical freshness/provider/state-owned.

#### v19.2.3 — Multi-User Security / Cost / Abuse / Capacity Hardening
Per account/user/device/capability attribution, quotas, throttling, abuse controls, cache/call avoidance, streaming usage, provider cost where known, fairness/starvation prevention, edge limits, circuit breakers and protected-session capacity. Threat-model/adversarial tests are mandatory.

#### v19.2.4 — Hosted Sync / Gateway Assurance Closure (#66)
No new feature scope. Failure injection and recovery matrix: cross-account denial, role/session/device revocation, duplicate replay, network loss during apply, checkpoint expiry/re-bootstrap, provider-right downgrade/expiry, secret rotation rollback, DB failover/PITR/restore, mixed-version clients, queue/backpressure and protected-session pressure. #66 closes only on executable Mac/Windows/web parity and zero-key/rights/sync/DR proof.

### v19.3.x — Professional Point-in-Time Evidence Substrate

#### v19.3.0 — Institutional / 13F Evidence Infrastructure
Direct SEC truth, manager/CIK/security mapping, amendments/restatements, filing-lag truth, point-in-time holdings, storage/indexing and outcome lineage.

#### v19.3.1 — Two-Sided Long / Short Thesis Evidence Substrate
Point-in-time thesis/plan snapshots, target/invalidation ordering, side-aware MFE/MAE and lawful/trustworthy short-interest/crowding/borrow/SSR context; explicit UNKNOWN where evidence is insufficient.

#### v19.3.2 — AODR Candidate / Ranking / Outcome Lineage
My Market vs Global truth, candidate/rank/reason snapshots, NOW/WATCH/PASS/ABSTAIN transitions, shared-ranking efficiency, diversity/correlation metadata, surfaced-vs-missed outcomes and recommendation usefulness.

### v19.4.x — Reliability, Economics & v20 Readiness

#### v19.4.0 — ADR-GDI Professional Reliability / Capacity
Capability SLOs, degradation history, provider/DB/runtime scorecards, restart/warm-start, query/index/pool tuning, load shedding, operating limits, protected-session reserve sizing and maintenance/preemption economics. Use controlled failure/soak/capacity evidence.

#### v19.4.1 — Specialized / Paid Provider Gap Evaluation
Only after measured v19 evidence proves material capability/quality/rights/cost gaps. Any new provider enters through the same router/rights/persistence/session-priority contracts; no special path.

#### v19.4.2 — v20 Research-Readiness Dataset / Lineage Audit
No new model scope. Prove point-in-time features/outcomes, provenance, rights, independence, synchronized-evidence safety, leakage controls and reliability history are adequate for governed adaptive research.

### v19.5.0 — v19 Major Closure
No new feature scope. Whole-system Principal Engineer / security / data / operational-readiness audit. Require zero unexplained provider role, zero unowned dataset, truthful commercial/data-rights posture, tested rollback/recovery, actual supported runtime/package evidence, #66 closure, SLO/capacity truth and zero unresolved P0 architecture gap before v20.

## 5. v20 — Adaptive Intelligence & Decision Research

**Entry:** `v19.5.0` Major Closure PASS. v20 learns only from point-in-time, rights-valid, provenance-bound evidence/outcomes. Production influence remains `SHADOW -> VALIDATED -> APPROVED -> PRODUCTION`; No Execution remains permanent.

### v20.0.x — Adaptive Research Control & Governance Foundation
- **v20.0.0 — Adaptive Research Control Plane + Immutable Experiment Ledger.** Dataset/version lineage, feature/provenance snapshots, cohorts, reproducibility, leakage controls, promotion/rollback evidence.
- **v20.0.1 — Model / Prompt Governance + Champion/Challenger.** Moved before broad adaptive rollout: model/prompt identities, independent evaluation, explainability, approval/rollback, drift and evidence-bound promotion.
- **v20.0.2 — Historical Analogues + Regime-Conditioned Outcomes.** Point-in-time analogue retrieval and conditioned outcome distributions without altering deterministic truth.
- **v20.0.3 — Calibration / FP-FN / Miss / Contradiction / Drift Intelligence.** Calibration, false positive/negative, missed opportunity, contradiction, drift and abstention thresholds.

### v20.1.x — ASBI
- **v20.1.0 — ASBI Behavioral Fingerprints + State Transitions.** Canonical behavior features, hierarchical context and immutable Behavior Intelligence Ledger.
- **v20.1.1 — ASBI Scenarios / Probability Momentum / Calibration.** Competing paths, multi-horizon outlooks, expected-move distributions, evidence sufficiency and ABSTAIN/NO RELIABLE EDGE.

### v20.2.x — Institutional + Two-Sided Thesis Intelligence
- **v20.2.0 — Adaptive Institutional / 13F Intelligence.** Manager fingerprints, persistence/concentration, consensus/crowding, rotation and calibrated usefulness.
- **v20.2.1 — TDTI Competing Long / Short / No Reliable Edge.** Same canonical snapshot; separate direction probability, thesis strength/confidence/opportunity quality and cause-aware confirmation/invalidation.
- **v20.2.2 — TDTI Two-Sided Trade-Plan Validation.** Long/Short entry/target/invalidation/R:R, readiness, time-to-resolution, MFE/MAE and historical calibration; still No Execution.

### v20.3.x — AODR
- **v20.3.0 — Adaptive Shared Opportunity Ranking.** Cross-candidate ranking using governed ASBI/TDTI readiness/quality with extension/chase/R:R/degradation penalties and surfaced-vs-missed outcomes.
- **v20.3.1 — Diversity / Opportunity Cost / Personalized Relevance.** Shared truth first; user relevance second. Correlation/theme/catalyst diversity and ABSTAIN/no-strong-opportunity remain valid outcomes.

### v20.4.x — Adaptive Operations
- **v20.4.0 — ADR-GDI Adaptive Optimization.** SHADOW/Champion-Challenger learning for provider recovery, cooldown/backoff, workload priority, maintenance value, reserve sizing and capacity policy. Cannot self-promote or reduce live-session protection without evidence/approval.

### v20.5.0 — v20 Professional Closure
No new feature scope. Calibrated utility, drift/abstention, deterministic-boundary protection, privacy/security/data rights, rollback, reproducibility, actual supported artifacts, zero silent self-modification and No Execution.

## 6. Cross-version dependency contract

`v18.9.x stabilize + instrument + validate -> v19.0 control plane -> v19.1 hosted data plane/sync -> v19.2 platform parity/assurance -> v19.3 evidence substrate -> v19.4 reliability/readiness -> v19.5 closure -> v20 governed adaptive research`

No later phase may compensate for an earlier unresolved foundation by creating a parallel owner or by hiding uncertainty behind model confidence.

## 7. G0–G16 industry-strength checkpoints for hosted/adaptive work

Without creating G17+, affected releases must include within existing gates:
- **G2:** canonical-owner map + Architecture Decision Record for material irreversible decisions; tenant/data-classification and threat-model scope;
- **G3:** API/schema/protocol contracts, migration/rollback plan, compatibility matrix, SLO/error-budget targets, observability and failure-injection plan;
- **G4:** implementation + unit/contract tests + feature/kill-switch controls where useful;
- **G7:** negative authorization/tenant isolation, rights, secret/redaction and adaptive-governance evidence;
- **G8:** load/soak/capacity, backpressure, chaos/failure recovery, DB/pool/provider limits and protected-session proof;
- **G9:** role-aware Mac/Windows/web UX and direct-route/API denial consistency;
- **G10:** production-readiness reconciliation; unresolved P0 threat/rights/recovery/compatibility item blocks freeze;
- **G11/G12:** immutable RC + full certification;
- **G13/G14:** native package/provenance/runtime proof where applicable;
- **G15:** controlled promotion/canary or bounded rollout evidence for hosted-risk changes where applicable;
- **G16:** implementation-miss review, incident/metric learning, obsolete machinery cleanup and exact next handoff.

## 8. Protected-session resource contract

Pre-market, regular market and after-hours remain Tier-0 decision-support sessions. Their live/current workloads always outrank maintenance and background synchronization. Provider quota/headroom, network concurrency, CPU, memory, DB/pool and workers retain explicit reserve. Maintenance/sync uses bounded surplus capacity and yields/preempts before it can materially degrade current truth.

Session boundaries come only from the canonical U.S. market calendar/session owner.

## 9. CI / release efficiency

For each patch: coherent code+test batch before PR; one product PR; exact-head FAST once for coherent candidate; Qualified once when Ready; one canonical G11–G16 release when release-capable. Classify failure before rerun. Reuse unchanged evidence when fingerprints/dependencies remain equivalent. No retry/certification branch families.

## 10. Protected boundaries

Smart Provider Router v2 sole routing authority; canonical freshness/recovery sole freshness owner; existing multi-feed allocator sole subscription owner; BroadSnapshotBroker canonical reuse owner; canonical persistence/cache/state owners reused; canonical U.S. session calendar reused; direct SEC/EDGAR authoritative; `SUPER_OWNER/OWNER/ADMIN/USER/DEMO` capability truth canonical; deterministic Day/Swing/Long protected; U.S. Equities Processing; GLD/SLV/USO actionable exceptions; No Execution; G0–G16 only.

## Exactly one next action

Execute G0 for issue #64 / `v18.9.1` using complete macOS crash evidence or deterministic reproduction, then freeze its narrow G1. Do not start `v18.9.2` or any v19 implementation branch until `v18.9.1` is truthfully closed or the crash is proven external/non-product.
