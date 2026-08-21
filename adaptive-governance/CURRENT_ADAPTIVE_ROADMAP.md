# DE.PULSE — Current Adaptive Roadmap

**Certified Stable:** `v18.9.0-stable`  
**Certified candidate:** `9ea81cddae4875ae15d3719ca028519a36c597b6`  
**Certified fingerprint:** `a8719090c341c874dbd1279cc31ad98e84075d5701c46a800bf951340780ecb9`  
**Build ID:** `v18.9.0-stable-20260821`  
**Master corrective program:** issue #65 / `ADAPT-PROVIDER-INTELLIGENCE-010`  
**Hosted architecture program:** issue #66 / `ADAPT-HOSTED-SYNC-001`  
**Immediate next product patch:** issue #64 / `ADAPT-RUNTIME-CRASH-001` -> `v18.9.1`.

## 1. Roadmap north star

DE.PULSE evolves through three ordered maturity stages:

`v18.9.x trustworthy runtime/data plane -> v19 professional multi-tenant hosted/data infrastructure -> v20 governed adaptive intelligence`

Permanent boundaries:
- U.S. Equities Processing only; GLD/SLV/USO remain explicit actionable exceptions;
- No Execution;
- G0-G16 remains the only top-level release model;
- one canonical owner per responsibility;
- shared/canonical symbol intelligence first, authorized personal composition second;
- persistence-first/reuse-first before external acquisition;
- Smart Provider Router v2 remains the sole provider-routing authority;
- canonical freshness, cache, persistence, subscription, SEC, market-session/calendar, identity and state owners are reused;
- provider/API availability never justifies fetching, storing or displaying data by itself;
- provider legal rights, DE.PULSE product entitlement and RBAC authorization are separate controls and all must pass before data is served;
- adaptive production influence follows `SHADOW -> VALIDATED -> APPROVED -> PRODUCTION` with explicit rollback;
- no silent deterministic Day/Swing/Long change, silent self-promotion or confidence invented around missing evidence.

## 2. Version-alignment contract

Starting from v18.9, version numbers communicate dependency phases while preserving small-patch delivery.

- **Major** = strategic maturity generation (`v19` infrastructure, `v20` adaptive intelligence).
- **Minor band** = coherent platform phase with a clear entry/exit contract (`v19.0.x`, `v19.1.x`, etc.).
- **Patch** = one primary independently certifiable responsibility.
- Planned future numbers are reservations only; exact scope/version identity freezes at G1.
- G0/G1 may split an oversized unstarted packet and shift later unstarted reservations inside the same band.
- Shipped versions are immutable and are never renumbered.
- Corrective/security work may interrupt the train; dependency truth is reconciled before planned work resumes.
- Minor bands are not permission to bundle unrelated work. Every patch still traverses G0-G16 independently unless an approved closure release is explicitly evidence-only.

## 3. v18.9.0 — COMPLETE / IMMUTABLE STABLE

Issue #61 / `ADAPT-TRADEINSIGHT-001` is complete. Qualified source `9e86b5e731f7a585cc77c1521f3639fc7a208efc`, certified candidate `9ea81cddae4875ae15d3719ca028519a36c597b6`, Release #32 and `release/v18.9.0/stable-evidence-manifest.json` remain immutable authority.

Post-Stable findings do not rewrite v18.9.0; they enter #65 as bounded corrective/professional patches.

## 4. v18.9.x — Stabilize -> instrument -> validate -> operationalize -> close

Architecture re-audit moved provider observability **before** broad TradeInsight capability admission. This is intentional: SHADOW promotion and provider usefulness must be measured rather than assumed.

1. **v18.9.1 — Runtime crash corrective ONLY** — #64. Evidence-based macOS Apple Silicon SIGABRT diagnosis/fix, persisted-state/API-key continuity, warm-state/relaunch regression and packaged-runtime proof.
2. **v18.9.2 — TradeInsight Settings/API-key UX ONLY** — canonical Settings/local secret owner; masked Save/Test/Clear; truthful state; save preserves scroll/focus/context.
3. **v18.9.3 — Coverage-aware Smart Provider Router core ONLY** — memory/DB reuse first, exact residual-gap acquisition, deterministic merge/provenance/coverage re-evaluation; validation lifecycle remains separate from serving role.
4. **v18.9.4 — Canonical company/instrument identity ONLY** — one identity owner reused by desks/Research/Discovery/Add Symbol.
5. **v18.9.5 — Market Data Modes + capability diagnostics ONLY** — behavior/quality modes and source/freshness/coverage diagnostics; no provider-brand mode owner.
6. **v18.9.6 — Provider observability / Adaptive telemetry foundation ONLY** — coverage, calls avoided, provider contribution/usefulness, freshness, rate/backpressure, disagreement, shared-work efficiency, runtime pressure and protected-session headroom.
7. **v18.9.7 — TradeInsight SEC Form 4 enrichment ONLY** — contract-validated SHADOW-first enrichment/corroboration; direct SEC/EDGAR authoritative; source-family de-duplication; measured via v18.9.6 telemetry.
8. **v18.9.8 — TradeInsight ticker/company search ONLY** — canonical symbol validation/company identity fallback/corroboration.
9. **v18.9.9 — TradeInsight movers/ranking evidence ONLY** — contract-validated SHADOW candidate evidence into canonical Opportunity Radar/ranker.
10. **v18.9.10 — Remaining useful TradeInsight capability admission ONLY** — every useful entitlement gets named consumer/owner/lifecycle/serving-role/freshness/rights/retention/rate/cost/Market-Mode disposition.
11. **v18.9.11 — Session-Aware Data Readiness Maintenance ONLY** — one canonical light-overnight/heavy-weekend coordinator with strict protected pre-market/regular/after-hours reserves, preemption/checkpoint/resume and no blind refetch.
12. **v18.9.12 — v18.9.x Professional Closure ONLY** — no new feature scope; zero-miss/duplicate-owner audit, #57/#64 regression, runtime/provider/persistence/session/package proof and Adaptive Intelligence Scorecard.

### v18.9.x exit contract

v18.9.12 must prove:
- native runtime reliability and warm-state/relaunch safety;
- canonical Settings/secret behavior;
- coverage-aware persistence-first fulfillment;
- stable canonical instrument/company identity;
- behavior-oriented Market Data Modes;
- telemetry sufficient to evaluate provider usefulness, cost and protected-session headroom;
- TradeInsight capabilities admitted only through canonical owners/lifecycle rules;
- session-aware maintenance with no material protected-session degradation;
- actual macOS Apple Silicon + Windows x64 package/runtime evidence;
- zero unexplained carry-forward, orphan useful capability or duplicate owner.

## 5. v19 — Professional Data Infrastructure + Multi-Tenant Hosted Account Platform

**Entry:** `v18.9.12` PASS.

v19 has two coupled missions:
1. professionalize data quality, rights, provenance, reliability, economics and point-in-time evidence; and
2. deliver #66: one authorized DE.PULSE identity across macOS, Windows and hosted web with SQLite as native edge/offline working state, PostgreSQL as shared hosted authority and a zero-key hosted Provider Gateway.

### Hosted authority model

```text
macOS / Windows
  SQLite edge/offline working set
        │
        │ typed authenticated incremental sync
        ▼
DE.PULSE hosted API / services
        │
        ├─ tenant/account + user/device/session truth
        ├─ RBAC/capabilities
        ├─ DE.PULSE product-plan entitlement
        ├─ upstream provider legal-rights gate
        ├─ Smart Provider Router v2 + canonical freshness/cache/persistence
        ├─ existing multi-feed subscription owner
        └─ server-side managed provider secrets
        │
        ▼
PostgreSQL shared account/state authority
        ▲
        │
 hosted web through same APIs
```

Normal commercial users authenticate only to DE.PULSE and never receive/configure DE.PULSE-owned provider credentials.

### Four separate hosted authorization dimensions

These may interact but must not be collapsed into one field:
1. **Tenant/account identity** — which DE.PULSE tenant/account owns the request/data.
2. **RBAC/capability authorization** — what this user/device/session may do.
3. **DE.PULSE product entitlement** — what the customer plan/status/quota allows DE.PULSE to expose.
4. **Provider legal/data rights** — what DE.PULSE is contractually permitted to acquire/cache/retain/derive/display/redistribute.

Every hosted read/write/stream projection applies the relevant dimensions before returning data.

### v19.0.x — Governance, Control Plane & Data Foundation

- **v19.0.0 — Provider Capability / Legal Rights Registry.** Machine-readable provider/capability rights for server proxying, commercial/multi-user use, caching/retention, redistribution/display, derived/AI use, attribution, environment/limits and expiry. Fail closed on unknown/expired rights.
- **v19.0.1 — Hosted Tenant/Identity/Device/Session Control Plane.** Tenant context becomes first-class; canonical `SUPER_OWNER/OWNER/ADMIN/USER/DEMO`, account/user/device lifecycle, revocation, privileged re-authentication, audit and API/SSE/native/web session truth.
- **v19.0.2 — DE.PULSE Product Entitlement / Metering Policy.** Billing-provider-agnostic plan/status/feature/quota/grace/suspension model and metering dimensions. This is separate from provider rights and RBAC. External checkout/invoicing integration is not required for internal pilot, but entitlement enforcement is required before external commercial multi-user activation.
- **v19.0.3 — PostgreSQL Tenancy / Schema / Pool / HA-PITR Foundation.** Shared-authority schema, tenant isolation, migration ownership, indexes/pools/capacity, encrypted backup, restore/failover, PITR and frozen RPO/RTO. No broad sync activation.
- **v19.0.4 — Managed Secrets / KMS Lifecycle.** Environment isolation, least privilege, rotation/rollback/compromise revoke, redaction and auditable provider-secret administration. No platform provider secret on commercial clients.
- **v19.0.5 — Provider Quality / Cost / Coverage / SLO Scorecards.** Measured freshness/completeness/latency/reliability/rate pressure/usefulness/cost/calls-avoided/fallback quality and tenant-aware operational signals.
- **v19.0.6 — Data Reconciliation / Revision / Point-in-Time Quality.** Source independence, disagreement/reconciliation, historical/corporate-action correctness, revision preservation and observed/effective timestamp lineage.

**Exit:** rights, tenant identity, product entitlement, PostgreSQL recovery, secrets, observability and data-quality foundations are executable/evidenced before shared hosted activation.

### v19.1.x — Zero-Key Hosted Provider Data Plane + Native Sync Foundation

- **v19.1.0 — Authenticated Hosted Provider Gateway.** Existing Smart Provider Router v2 behind a versioned hosted API boundary; canonical cache/persistence/freshness reuse; server-side provider credentials; bounded circuits/backpressure/kill-switch/degraded behavior; API inventory/deprecation ownership.
- **v19.1.1 — Unified Serving Policy + Live Fan-Out Isolation.** Tenant/RBAC/product-entitlement/provider-rights checks gate cache/persistence/REST/WebSocket/SSE projections. Existing multi-feed owner remains upstream subscription authority. Prevent cross-tenant or higher-tier/right-restricted leakage.
- **v19.1.2 — Sync Protocol Foundation.** Snapshot/high-watermark bootstrap, SQLite atomic outbox, authenticated idempotent push, authoritative server sequence/change log, incremental pull, checkpoints, tombstones, retention/compaction, stale-device re-bootstrap and mixed-version negotiation. No raw DB replication.
- **v19.1.3 — macOS Preferences + Watchlist Pilot.** Offline/restart/reconnect/conflict/user-switch/local-account/lost-device proof.
- **v19.1.4 — Desks / Workspaces Sync.** Versioned membership/configuration/delete/history semantics through the same transport.

### v19.2.x — Cross-Platform Account Parity & Hosted Assurance

- **v19.2.0 — Windows x64 Sync Parity.** Same protocol/security/account truth; no Windows fork.
- **v19.2.1 — Hosted Web Account Parity.** Same APIs/session/capability/product-entitlement/PostgreSQL state; browser cache non-authoritative.
- **v19.2.2 — Rights-Aware Research / Durable State Portability.** Only product-owned, lawful, entitled, provenance-bound durable artifacts; live market truth remains freshness/provider/state-owned.
- **v19.2.3 — Tenant-Aware Metering / Cost / Usage Observability.** Tenant/account/user/device/capability attribution, plan/quota consumption, cache/call avoidance, streaming usage, provider cost where known and tenant-health signals.
- **v19.2.4 — Multi-User Security / Abuse / Capacity Hardening.** Rate limits, fairness/noisy-neighbor controls, sensitive-flow abuse protection, circuit/load shedding, edge/API protection and protected-session capacity.
- **v19.2.5 — #66 Hosted Sync / Gateway Assurance Closure.** Failure/adversarial/recovery matrix; Mac/Windows/web parity; zero-key secrets; tenant/product-entitlement/provider-right isolation; mixed-version sync; DB/secret/provider-right recovery.

### v19.3.x — Professional Point-in-Time Evidence Substrate

- **v19.3.0 — Institutional / 13F Evidence Infrastructure.** Direct SEC truth, identity/mapping, amendments, filing lag, point-in-time holdings and outcome lineage.
- **v19.3.1 — Two-Sided Long/Short Thesis Evidence Substrate.** Point-in-time plans, target/invalidation ordering, side-aware MFE/MAE and lawful short/crowding evidence with explicit UNKNOWN.
- **v19.3.2 — AODR Candidate / Ranking / Outcome Lineage.** Candidate/rank/reason/transition history, shared-ranking efficiency, diversity/correlation metadata and surfaced-vs-missed outcomes.

### v19.4.x — Reliability, Economics & v20 Readiness

- **v19.4.0 — ADR-GDI Professional Reliability / Capacity.** SLO/error-budget/degradation history, provider/DB/runtime scorecards, warm start, indexes/pools, load shedding, reserve sizing, maintenance economics and hosted operational runbooks.
- **v19.4.1 — Specialized/Paid Provider Gap Evaluation.** Add/replace only where measured evidence proves a material capability/quality/rights/cost gap; same router/rights/persistence/session contracts.
- **v19.4.2 — v20 Research-Readiness Audit.** Prove point-in-time evidence/features/outcomes/provenance/rights/independence/leakage/reliability readiness.

### v19.5.0 — v19 Major Closure

No new feature scope. Principal Engineer + security/data/operational/commercial-readiness audit. Require #66 PASS, tenant isolation, product-entitlement/provider-right separation, truthful rights/commercial posture, API inventory/version compatibility, recovery/rollback drills, SLO/capacity evidence, supported runtime/package proof, zero unowned dataset/provider role and zero unresolved P0 architecture gap.

## 6. v20 — Adaptive Intelligence & Decision Research

**Entry:** `v19.5.0` PASS.

Model/prompt governance is established **before** broad adaptive-model rollout so evaluation, rollback and promotion exist before model influence grows.

### v20.0.x — Adaptive Research Control & Governance
- **v20.0.0 — Adaptive Research Control Plane + Immutable Experiment Ledger.**
- **v20.0.1 — Model/Prompt Governance + Champion/Challenger.**
- **v20.0.2 — Historical Analogues + Regime-Conditioned Outcomes.**
- **v20.0.3 — Calibration / FP-FN / Miss / Contradiction / Drift.**

### v20.1.x — ASBI
- **v20.1.0 — Behavioral Fingerprints + State Transitions.**
- **v20.1.1 — Scenarios / Probability Momentum / Calibration.**

### v20.2.x — Institutional + TDTI
- **v20.2.0 — Adaptive Institutional / 13F Intelligence.**
- **v20.2.1 — TDTI Competing Long / Short / No Reliable Edge.**
- **v20.2.2 — TDTI Two-Sided Trade-Plan Validation.** No Execution.

### v20.3.x — AODR
- **v20.3.0 — Adaptive Shared Opportunity Ranking.**
- **v20.3.1 — Diversity / Opportunity Cost / Personalized Relevance after Shared Truth.**

### v20.4.x — Adaptive Operations
- **v20.4.0 — ADR-GDI Adaptive Optimization.** SHADOW/Champion-Challenger provider/recovery/workload/maintenance/reserve learning; no self-promotion.

### v20.5.0 — v20 Professional Closure

No new feature scope. Calibrated utility/drift/abstention, deterministic-boundary protection, privacy/security/data rights, reproducibility/rollback, actual supported artifacts, zero silent self-modification and No Execution.

## 7. Industry-strength controls inside G0-G16

No G17+ is introduced. Hosted/security/data/adaptive controls are absorbed into existing gates:
- **G2:** canonical owner map; tenant/data classification; threat model; ADR for material irreversible choices;
- **G3:** API/schema/protocol compatibility, inventory/deprecation, migration/rollback, SLO/error budget, observability, failure-injection and rollout plan;
- **G4:** contract tests and bounded feature/kill-switch controls where useful;
- **G7:** negative authorization/tenant/product-entitlement/provider-right/secret/adaptive tests;
- **G8:** load/soak/capacity, noisy-neighbor/fairness, backpressure, chaos/failover and protected-session proof;
- **G9:** cross-platform role-aware UX and direct-route/API denial;
- **G10:** production-readiness reconciliation; unresolved P0 threat/rights/recovery/compatibility gap blocks freeze;
- **G11/G12:** immutable RC/full certification;
- **G13/G14:** package/provenance/runtime proof;
- **G15:** controlled/canary activation for hosted-risk changes where applicable with rollback/kill switch;
- **G16:** implementation-miss, incident/metric learning, cleanup and exact next handoff.

## 8. Why this ordering is intentional

- Observability precedes provider expansion so SHADOW usefulness/promotion is measurable.
- Tenant identity, provider rights, product entitlement, PostgreSQL recovery and secrets precede hosted activation so security/commercial controls are not retrofits.
- Gateway + unified serving policy precede broad sync/client parity so all clients consume one trust boundary.
- Sync foundation precedes Mac/Windows/web/research parity so there is one protocol/state authority.
- Tenant-aware metering and capacity controls precede #66 closure so commercial scale is measurable, not inferred.
- Point-in-time evidence precedes adaptive learning so research avoids look-ahead/revision leakage.
- Model governance precedes broad adaptive rollout so promotion/rollback exist before production influence.

## Exactly one next action

Diagnose #64 from complete macOS crash evidence or deterministic reproduction and freeze the narrow `v18.9.1` G1. Do not start `v18.9.2` or any v19 product implementation branch until `v18.9.1` is truthfully closed or proven external/non-product.