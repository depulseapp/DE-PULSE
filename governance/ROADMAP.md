# DE.PULSE — Canonical Adaptive Roadmap

**Status:** APPROVED / ADAPTIVE  
**Authority:** canonical product sequencing and approved strategic workstreams  
**Rule:** mission/workstreams are durable; future minor/patch placement may adapt when dependency, safety, performance, evidence, commercial rights, or operational reality justify movement. Shipped releases and certified evidence are immutable truth.

---

## 1. Canonical roadmap truth model

This roadmap distinguishes **approved direction** from **certified implementation reality**.

Permanent rules:
- a historical roadmap label does not prove that capability shipped in that version;
- actual Stable tags, release evidence, source/artifact provenance and current handoff define what was truly delivered;
- future placement may be re-ordered only with durable rationale and synchronized Build Plan / Build Process / Delivery Process / handoff updates;
- no roadmap edit may retroactively claim an unimplemented capability was delivered;
- corrective/security/reliability work may preempt planned feature ordering;
- every known miss is implemented in-scope or durably assigned to a named later release;
- G0–G16 remains the only release model.

### Historical placement reconciliation

Earlier versions of this roadmap provisionally placed **PostgreSQL / Hosted Shared State** in `v18.3` and broader hosted/commercial security in `v18.4`. The certified product evolved through later v18.x work without establishing that hosted PostgreSQL/multi-device architecture as production authority. Those historical labels are therefore **superseded as future placement, not rewritten as delivered history**.

Current certified truth is:
- `v18.9.0-stable` is the immutable Stable baseline;
- `v18.9.x` remains native-first and **must not introduce hosted PostgreSQL synchronization**;
- PostgreSQL shared authority, zero-key hosted Provider Gateway and Mac/Windows/web account parity are explicitly planned for `v19` under issue #66 / `ADAPT-HOSTED-SYNC-001`;
- all previously approved hosted/security/data-rights objectives remain durable, but now land in the dependency-correct v19 control-plane/data-plane train below.

Historical v18 recovery/reconciliation detail remains preserved in its dedicated governance records, including `governance/V18_ADAPTIVE_RECOVERY_AND_CLOSURE_PROGRAM.md`, release evidence and archived reconciliation files.

---

## 2. Permanent product / architecture boundaries

DE.PULSE is a **U.S. equities research, intelligence and decision-support system**.

Permanent constraints:
- **No Execution** — no broker routing, order execution, paper/live trading, OMS/blotter or autonomous/semi-autonomous execution;
- U.S. Equities Processing remains the product market boundary; `GLD`, `SLV`, `USO` remain explicit actionable tradable exceptions;
- Smart Provider Router v2 is the sole executable provider-routing authority;
- canonical freshness/recovery is the sole freshness truth owner;
- the existing multi-feed allocator/subscription owner remains canonical;
- BroadSnapshotBroker/shared acquisition remains canonical reuse ownership where applicable;
- canonical persistence/cache/state, symbol identity, session/calendar and telemetry owners are reused;
- direct SEC/EDGAR remains authoritative for filing truth/Form 4 authority;
- equivalent lawful symbol/evidence work is processed canonically once and fan-outs to authorized consumers rather than `users × symbols` duplicate pipelines;
- normal UI shows synthesized intelligence, not provider/runtime plumbing;
- deterministic Day/Swing/Long truth remains protected unless separately governed;
- adaptive production influence follows **SHADOW → VALIDATED → APPROVED → PRODUCTION** with evidence, rollback and explicit approval;
- no silent self-modification and no model confidence used to hide missing/weak evidence.

---

## 3. Durable approved strategic workstreams

The following workstreams remain approved across the current roadmap. Exact implementation version follows the dependency-aligned train in sections 4–6.

### Smart Intelligent Provider Router v2
Capability/entitlement-aware, coverage-aware provider routing with Preferred vs Serving semantics, deterministic cooldown/circuits, provider budgets, protected-session headroom, source disagreement handling, persistence-first reuse, residual-gap acquisition, calls avoided and usefulness telemetry.

### Shared Symbol Intelligence / Multi-User Demand Union
Global Symbol Registry and shared canonical intelligence serve overlapping authorized demand. Per-user symbols/watchlists/preferences remain isolated, while equivalent provider acquisition/calculation/synthesis is reused where lawful.

### Adaptive Opportunity Discovery & Recommendations (AODR) Foundation
Reuse Global Symbol Registry, Reliable Actionable Universe, Shared Symbol Intelligence and Opportunity Radar to support:
- My Market vs Global Opportunities;
- bounded candidate ranking;
- NOW / WATCH / PASS / ABSTAIN;
- point-in-time rank/reason/bucket lineage and outcomes;
- shared ranking efficiency;
- diversity/correlation context;
- concise `why now / confirms / invalidates / what changes the view` synthesis;
- ADR-GDI suppression/demotion under weak/degraded required evidence.

AODR is not a second scanner/ranker/provider silo.

### TradeInsight — SHADOW / Secondary Intelligence
TradeInsight operates only through canonical provider/data owners. Approved candidate roles include insider/congressional context, historical OHLCV backfill/reconciliation, corporate actions, symbol metadata/search corroboration, Opportunity Radar candidate evidence and future controlled AI/MCP research where rights permit. No TradeInsight-specific router/cache/scanner/scheduler/Market Mode/SEC truth/symbol registry/persistence system.

### Provider → Market Mode Adaptive Integration
Every provider capability gets an explicit disposition: `INTEGRATED`, `CONTEXTUAL_ONLY`, `NOT_RELEVANT`, or `INTENTIONALLY_HIDDEN`. Provider count never changes a mode. Deterministic/statistical code owns numeric market truth; LLM/adaptive reasoning may extract/correlate/compare/explain but cannot directly set a Market Mode or silently reweight protected formulas.

Machine authority remains `provider_market_mode_integration_registry.json` inside the Functionality Utility & Integration checkpoint.

### Institutional Holdings / 13F Evidence
Build lawful point-in-time SEC 13F evidence with manager/CIK/security mappings, amendments/restatements, report/acceptance time, normalized disclosed holdings, quarter-over-quarter states, corporate-action reconciliation, limitations and subsequent outcomes. Direct SEC EDGAR remains canonical filing truth.

### Two-Sided Thesis Evidence / TDTI Foundation
Preserve one canonical ticker/horizon evidence snapshot, Entry/Target/Invalidation truth, short-entry context where supportable, first-event ordering, side-aware MFE/MAE, ASBI/regime/catalyst/liquidity/institutional context and lawful short-interest/crowding/borrow/SSR evidence with explicit UNKNOWN/ABSTAIN when unavailable.

TDTI is not a separate short-trading product and remains No Execution.

### Adaptive Data Reliability & Graceful Degradation Intelligence (ADR-GDI)
Reliability is foundational, not deferred to v20. Durable objectives include:
- capability-level health and degradation reasons;
- consumer/dependency blast radius;
- dataset/horizon/session freshness SLOs;
- provider retry/circuit discipline;
- duplicate-work elimination and single-flight/coalescing;
- bounded queues/backpressure/load shedding;
- warm canonical persistence/restart recovery;
- truthful UNKNOWN/degraded/ABSTAIN;
- scoped user messaging and deeper Maintenance diagnostics;
- reliability/outcome history for later governed adaptive optimization.

---

## 4. v18.9.x — Trustworthy Runtime, Acquisition, Observability & Operational Closure

**Certified entry:** `v18.9.0-stable`  
**Master program:** #65 / `ADAPT-PROVIDER-INTELLIGENCE-010`  
**Immediate blocker:** #64 / `ADAPT-RUNTIME-CRASH-001` → `v18.9.1`

The architecture audit deliberately places **observability before broad provider capability expansion** so SHADOW admission/promotion is evidence-based.

1. **v18.9.1 — Runtime crash corrective ONLY.** Evidence-based packaged macOS Apple Silicon SIGABRT diagnosis/fix; preserve SQLite/user state/API keys; lifecycle/relaunch regression and packaged proof.
2. **v18.9.2 — TradeInsight Settings/API-key UX ONLY.** Canonical Settings/local secret owner; masked Save/Test/Clear; truthful connection/capability status; preserve scroll/focus.
3. **v18.9.3 — Coverage-aware Smart Provider Router core ONLY.** Memory/persisted evidence first; validate freshness/schema/provenance/rights; exact residual-gap acquisition; deterministic reconciliation/provenance/coverage re-evaluation; lifecycle separate from serving role.
4. **v18.9.4 — Canonical company/instrument identity ONLY.** One identity owner reused by Desks/Research/Discovery/Add Symbol with symbol-only fallback.
5. **v18.9.5 — Market Data Modes + capability diagnostics ONLY.** Behavior/quality modes and truthful source/freshness/coverage/fallback diagnostics; no provider-brand mode owner.
6. **v18.9.6 — Provider Observability / Adaptive Telemetry Foundation ONLY.** Coverage, residual gaps, DB/cache reuse, calls avoided, contribution/usefulness, latency/errors/rate/backpressure, freshness failures, disagreement/corroboration, shared-work/fan-out efficiency, CPU/memory/DB/worker pressure and protected-session headroom.
7. **v18.9.7 — TradeInsight SEC Form 4 enrichment ONLY.** Contract-validated SHADOW-first enrichment; direct SEC/EDGAR authoritative; source-family de-duplication; promotion/usefulness measured by v18.9.6 telemetry.
8. **v18.9.8 — TradeInsight ticker/company search ONLY.** Canonical symbol validation/company identity fallback/corroboration.
9. **v18.9.9 — TradeInsight movers/ranking evidence ONLY.** SHADOW candidate evidence through existing Opportunity Radar/ranker; no parallel scanner/ranking engine.
10. **v18.9.10 — Remaining useful TradeInsight capability admission ONLY.** Every useful entitlement gets owner/consumer/lifecycle/serving-role/freshness/rights/retention/rate/cost/Market-Mode disposition; no invented endpoint/intraday capability or uncontrolled Python/MCP production dependency.
11. **v18.9.11 — Session-Aware Data Readiness Maintenance ONLY.** One canonical light-overnight + heavy-weekend/extended-closed coordinator; protected pre-market/regular/after-hours reserves; bounded value-driven work; preemption/checkpoint/resume; no blind refetch.
12. **v18.9.12 — v18.9.x Professional Closure ONLY.** No feature scope; implementation-miss/duplicate-owner audit, #57/#64 regressions, deterministic equivalence, provider/persistence/maintenance/protected-session proof, Adaptive Intelligence Scorecard, actual macOS Apple Silicon + Windows x64 runtime/package evidence.

### v18.9.x exit contract

v18.9.12 must produce a zero-gap-enough native foundation for v19:
- runtime stability;
- canonical secrets/settings/identity;
- persistence-first coverage-aware acquisition;
- useful provider telemetry;
- evidence-based TradeInsight admission;
- session-aware maintenance with no protected-session degradation;
- no unexplained useful capability or duplicate owner;
- exact artifact/runtime provenance.

---

## 5. v19 — Professional Data Infrastructure + Hosted Account Platform

**Entry:** `v18.9.12` closure PASS.  
**Hosted program:** #66 / `ADAPT-HOSTED-SYNC-001`.

v19 professionalizes data rights/quality/provenance/reliability/economics and delivers one DE.PULSE account across macOS, Windows and hosted web.

### Target hosted architecture

```text
macOS / Windows
  SQLite edge/offline working set
        │
        │ typed authenticated incremental sync
        ▼
DE.PULSE hosted API / services
        │
        ├─ canonical identity/session/capability truth
        ├─ Smart Provider Router v2
        ├─ canonical freshness/cache/persistence reuse
        ├─ rights/entitlement gate
        ├─ existing multi-feed subscription owner
        ├─ server-side managed provider secrets
        ▼
PostgreSQL shared account/state authority
        ▲
        │
 hosted web through same APIs
```

Commercial normal users are **zero-key**: they authenticate only to DE.PULSE and never receive/configure platform provider credentials.

### v19.0.x — Governance / Control Plane / Data Foundation

- **v19.0.0 — Provider Capability / Entitlement / Rights Registry.** Runtime-enforceable commercial/multi-user, proxying, cache/retention, redistribution/display, derived/AI-use, attribution, limits/environment/expiry policy.
- **v19.0.1 — Hosted Identity / Device / Session Control Plane.** Canonical `SUPER_OWNER/OWNER/ADMIN/USER/DEMO` + capabilities, account/device lifecycle, revocation, server-side API/SSE/native/web session truth and privileged re-authentication/strong-auth controls where appropriate.
- **v19.0.2 — PostgreSQL Tenancy / Schema / Pool / HA-PITR Foundation.** Shared authority, tenancy isolation, migrations, indexing/pool/capacity, encrypted backup, restore/failover, RPO/RTO; no broad sync activation.
- **v19.0.3 — Managed Secrets / KMS Lifecycle.** Environment isolation, least privilege, rotation/rollback/compromise revoke, redaction/audit; zero platform provider secret in commercial clients.
- **v19.0.4 — Provider Quality / Cost / Coverage / SLO Scorecards.** Measured freshness/completeness/latency/reliability/rate/cost/usefulness/calls-avoided/fallback/maintenance value and error-budget signals.
- **v19.0.5 — Data Reconciliation / Revision / Point-in-Time Quality.** Source independence/conflicts, historical corrections/adjustments, observed/effective/revision truth and provenance.

**Phase invariant:** multi-user hosted provider/sync activation waits for applicable rights, identity, DB/recovery and secret foundations.

### v19.1.x — Zero-Key Provider Data Plane + Native Sync Foundation

- **v19.1.0 — Authenticated Hosted Provider Gateway.** Wrap existing Smart Provider Router v2 behind hosted boundary; server-side credentials, cache/persistence reuse, bounded circuits/backpressure/kill controls; no second provider stack.
- **v19.1.1 — Rights/Entitlement + Live Fan-Out Isolation.** Machine-enforced router/cache/persistence/REST/WebSocket/SSE authorization and provider-right boundaries using the existing multi-feed owner; no premium/realtime/right leakage across accounts/plans.
- **v19.1.2 — Sync Protocol Foundation.** Snapshot/high-watermark bootstrap, stable IDs, SQLite atomic outbox, authenticated idempotent push, authoritative server sequence/change log, incremental pull, per-device checkpoints, tombstones, retention/compaction, stale-device re-bootstrap and mixed-version negotiation; never raw DB replication.
- **v19.1.3 — macOS Preferences + Watchlist Pilot.** Offline/restart/reconnect/conflicts, local account isolation/user switching, lost/revoked device and non-destructive state proof.
- **v19.1.4 — Desks / Workspaces Sync.** Versioned membership/configuration/delete/history semantics through the same transport/state owners.

### v19.2.x — Cross-Platform Account Parity + #66 Assurance

- **v19.2.0 — Windows x64 Sync Parity.** Same account/protocol/security semantics; no Windows-specific truth.
- **v19.2.1 — Hosted Web Account Parity.** Same authenticated APIs/session/capability/PostgreSQL account truth; browser cache non-authoritative.
- **v19.2.2 — Rights-Aware Research / Durable State Portability.** Only lawful, entitled, provenance-bound product-owned durable artifacts; live market truth remains provider/freshness/state-owned.
- **v19.2.3 — Multi-User Security / Cost / Abuse / Capacity Hardening.** Per-account/user/device/capability attribution/quotas/throttling, fairness/starvation prevention, edge limits, streaming/cost/cache-call efficiency, load shedding and protected-session reserve.
- **v19.2.4 — #66 Hosted Sync / Gateway Assurance Closure.** No feature scope; adversarial/failure/recovery matrix covering isolation/revocation, replay/network interruption, bootstrap/re-bootstrap, rights expiry, secret rotation, DB failover/PITR/restore, mixed clients, backpressure and protected-session pressure.

### v19.3.x — Professional Point-in-Time Evidence Substrate

- **v19.3.0 — Institutional / 13F Evidence Infrastructure.** Direct SEC truth, manager/security mapping, amendments, filing lag, point-in-time holdings, storage/indexing and outcome lineage.
- **v19.3.1 — Two-Sided Long / Short Thesis Evidence Substrate.** Point-in-time plan/thesis snapshots, first-event ordering, side-aware outcomes and lawful short/crowding evidence with explicit UNKNOWN.
- **v19.3.2 — AODR Candidate / Ranking / Outcome Lineage.** My Market/Global truth, point-in-time rank/reason/transition history, shared-ranking efficiency, diversity metadata and surfaced-vs-missed outcomes.

### v19.4.x — Reliability / Economics / v20 Readiness

- **v19.4.0 — ADR-GDI Professional Reliability / Capacity.** Capability SLO/degradation history, provider/DB/runtime scorecards, warm start, indexes/pools, bounded operating limits, load shedding, protected-session reserves and maintenance/preemption economics with soak/failure evidence.
- **v19.4.1 — Specialized / Paid Provider Gap Evaluation.** Add/replace only when measured quality/capability/rights/cost evidence proves a material gap; same canonical router/rights/persistence/session contracts.
- **v19.4.2 — v20 Research-Readiness Audit.** No model scope; prove point-in-time features/outcomes, provenance, rights, independence, synchronized-evidence safety, leakage controls and reliability history.

### v19.5.0 — v19 Major Closure

No new feature scope. Principal Engineer + security/data/operational/commercial audit. Require:
- #66 PASS across Mac/Windows/web;
- zero-key provider-secret boundary;
- machine-enforced provider rights/entitlements;
- deterministic sync/offline/conflict/delete/bootstrap behavior;
- PostgreSQL/secret/provider-right recovery drills;
- truthful SLO/error-budget/capacity/economics;
- point-in-time data/provenance readiness;
- actual supported runtime/package/service provenance;
- zero unowned dataset/provider role/duplicate owner and zero unresolved P0 architecture gap.

**Mandatory v19 Major Closure before v20.**

---

## 6. v20 — Adaptive Intelligence & Decision Research

**Entry:** `v19.5.0` PASS.

Purpose: use structured point-in-time evidence/outcome history to improve decision support without creating a silent self-modifying trading system.

The architecture audit moves **model/prompt governance before broad adaptive rollout** so identity/evaluation/rollback/promotion controls exist before adaptive intelligence gains production influence.

### v20.0.x — Adaptive Research Control & Governance

- **v20.0.0 — Adaptive Research Control Plane + Immutable Experiment Ledger.** Dataset/feature/provenance/cohort/version lineage, reproducibility, leakage controls and promotion/rollback evidence.
- **v20.0.1 — Model / Prompt Governance + Champion/Challenger.** Stable model/prompt/policy identity, independent evaluation, explainability, drift, approval/rollback and evidence-bound promotion.
- **v20.0.2 — Historical Analogues + Regime-Conditioned Outcomes.** Point-in-time analogue retrieval and conditioned outcome distributions without changing deterministic truth.
- **v20.0.3 — Calibration / FP-FN / Miss / Contradiction / Drift Intelligence.** Calibration, false positives/negatives, missed opportunities, contradictions, drift and abstention thresholds.

### v20.1.x — Adaptive Stock Behavior Intelligence (ASBI)

- **v20.1.0 — Behavioral Fingerprints + State Transitions.** Canonical behavior features, hierarchical symbol/peer/sector/regime context and immutable behavior ledger.
- **v20.1.1 — Scenarios / Probability Momentum / Calibration.** Competing paths, multi-horizon outlooks, expected-move distributions, evidence sufficiency/conflict and ABSTAIN/NO RELIABLE EDGE.

### v20.2.x — Adaptive Institutional + TDTI

- **v20.2.0 — Adaptive Institutional / 13F Intelligence.** Manager fingerprints, persistence/concentration, consensus/crowding, rotation, stale-data penalties and calibrated outcomes.
- **v20.2.1 — TDTI Competing Long / Short / No Reliable Edge.** Same canonical snapshot; separate Direction Probability, Thesis Strength, Confidence and Opportunity Quality with cause-aware confirmation/invalidation.
- **v20.2.2 — TDTI Two-Sided Trade-Plan Validation.** Long/Short entry/target/invalidation/R:R, readiness/time-to-resolution/MFE/MAE/risk and historical calibration; still No Execution.

### v20.3.x — Adaptive Opportunity Discovery & Recommendations

- **v20.3.0 — Adaptive Shared Opportunity Ranking.** Cross-candidate ASBI/TDTI-aware ranking with expected magnitude/time-to-resolution, extension/chase/R:R/degradation penalties and candidate-vs-surfaced-vs-missed outcomes.
- **v20.3.1 — Diversity / Opportunity Cost / Personalized Relevance.** Correlation/theme/catalyst diversity and user relevance only after shared canonical truth; ABSTAIN/no strong opportunities is valid.

### v20.4.x — Adaptive Operations

- **v20.4.0 — ADR-GDI Adaptive Optimization.** Governed SHADOW/Champion-Challenger learning for provider recovery, cooldown/backoff, workload priority, maintenance value, protected-session reserve sizing, fallback usefulness and capacity policy. No self-promotion or safety reduction without evidence/approval.

### v20.5.0 — v20 Professional Closure

No feature scope. Require calibration/decision utility/drift/abstention, deterministic-boundary protection, privacy/security/data rights, reproducibility, rollback, actual supported artifacts, zero silent self-modification and No Execution.

---

## 7. Detailed adaptive workstream placement

### ASBI
**v18/v19 preparation:** preserve lawful point-in-time structured evidence, behavior features, sequence/outcome context, provider provenance/freshness, regime/sector/catalyst/SEC/earnings/institutional context, historical adjustments and data-rights metadata.  
**v20 implementation:** behavior fingerprints, state transitions, competing scenarios, multi-horizon outlooks, probability momentum, hierarchical learning, catalyst-aware analogues, outcome distributions, confidence/data sufficiency, conflict/independence, ABSTAIN and calibration/drift under Champion/Challenger.

### Adaptive Institutional / 13F
**v19 substrate:** point-in-time SEC filing/report/acceptance identity, manager/security mappings, amendments/restatements, disclosed holdings, QoQ states, corporate-action reconciliation, filing-lag/confidentiality/incompleteness limitations and outcomes from public availability time.  
**v20 intelligence:** manager behavioral fingerprints, persistence/concentration, accumulation/reduction breadth, consensus/crowding, thematic/sector rotation, usefulness by regime/stock type, convergence/contradiction, stale-data penalties, calibrated outcomes and ABSTAIN.

### TDTI
**v19 substrate:** one canonical evidence snapshot; Long/Short plan/outcome lineage; first-event ordering; MFE/MAE; ASBI/regime/catalyst/liquidity/institutional/insider/congressional/options context; lawful short-specific data; explicit UNKNOWN.  
**v20 intelligence:** competing Long/Short/No Reliable Edge theses, direction probability vs thesis strength/confidence/opportunity quality, structural invalidation, per-side readiness/probability momentum, multiple paths/time-to-resolution, risk intelligence, cross-horizon reconciliation, concise WHY/CONFIRMS/INVALIDATES/WATCH synthesis and side-aware calibration.

### AODR
**v19 substrate:** My Market vs Global buckets, Global Symbol Registry/Actionable Universe eligibility, shared ranking inputs, Opportunity Radar promotion/demotion lifecycle, NOW/WATCH/PASS/ABSTAIN, point-in-time rank/reason/outcome history, diversity metadata, user relevance isolated from shared truth and ADR-GDI dependency.  
**v20 intelligence:** adaptive cross-candidate ranking, opportunity cost, extension/chase/R:R/degradation penalties, diversity-aware sets, personalized relevance after shared truth, candidate-vs-surfaced-vs-missed calibration and Champion/Challenger.

### ADR-GDI
**v18.9:** observability, freshness/headroom/usefulness telemetry, shared-work efficiency, session-aware maintenance, runtime/protected-session reliability.  
**v19:** professional SLO/degradation/provider/DB/runtime/capacity/economics hardening.  
**v20:** governed adaptive reliability optimization only after SHADOW/Champion-Challenger proof.

---

## 8. Industry-strength controls inside G0–G16

For material hosted/security/data/adaptive work, existing gates absorb the following without creating G17+:
- **G2:** canonical owner map, architecture decision record/equivalent, trust/tenant/data classification and threat-model scope;
- **G3:** API/schema/protocol compatibility, migration/rollback/roll-forward, SLO/error-budget, observability, negative/failure/load test plans and feature/kill-switch/canary strategy where useful;
- **G4:** implementation + unit/contract tests + backward-compatible migration behavior;
- **G7:** negative authorization/tenant isolation, rights/entitlement, secret/redaction and adaptive promotion evidence;
- **G8:** load/soak/capacity, queues/backpressure, circuits, provider/DB failure, failover/recovery and protected-session proof;
- **G9:** role-aware Mac/Windows/web UX and direct-route/API authorization consistency;
- **G10:** production-readiness reconciliation; unresolved P0 security/rights/recovery/compatibility/duplicate-owner issue blocks freeze;
- **G11/G12:** immutable RC + full certification;
- **G13/G14:** package/provenance/actual-artifact runtime proof where applicable;
- **G15:** bounded/canary/progressive promotion and explicit rollback/kill criteria for hosted-risk changes where applicable;
- **G16:** implementation-miss/incident/metric learning, obsolete machinery cleanup and authoritative next handoff.

---

# Permanent Adaptive Placement Rule

Before assigning any unshipped approved workstream to a release, evaluate:

**value + dependency + architecture + safety + performance + defects + data/provider readiness + commercial/data rights + test evidence + operational recovery + accumulated outcomes**

If moving a workstream is safer/more correct, update this canonical roadmap, the four CURRENT Adaptive governance files, relevant issue/handoff and `governance/DECISION-LOG.md` rather than silently changing scope.
