# DE.PULSE — Current Adaptive Roadmap

**Certified Stable:** `v18.9.0-stable`  
**Certified candidate:** `9ea81cddae4875ae15d3719ca028519a36c597b6`  
**Certified fingerprint:** `a8719090c341c874dbd1279cc31ad98e84075d5701c46a800bf951340780ecb9`  
**Build ID:** `v18.9.0-stable-20260821`  
**Master corrective program:** issue #65 / `ADAPT-PROVIDER-INTELLIGENCE-010`  
**Immediate blocker / next patch:** issue #64 / `ADAPT-RUNTIME-CRASH-001` -> `v18.9.1`.

## 1. Roadmap north star

DE.PULSE evolves through three ordered maturity stages:

`v18.9.x trustworthy runtime/data acquisition -> v19 professional hosted/data infrastructure -> v20 governed adaptive intelligence`

The roadmap follows permanent principles:
- U.S. Equities Processing only; GLD/SLV/USO remain explicit actionable exceptions;
- No Execution;
- one canonical owner per responsibility;
- shared/canonical symbol intelligence first, authorized personal composition second;
- persistence-first/reuse-first before external acquisition;
- Smart Provider Router v2 remains sole routing authority;
- canonical freshness, cache, persistence, subscription, SEC, session/calendar, identity and state owners are reused;
- provider/API availability never justifies collection/display by itself;
- data rights/entitlements are enforceable runtime inputs, not documentation-only metadata;
- adaptive production influence follows `SHADOW -> VALIDATED -> APPROVED -> PRODUCTION` with rollback;
- no silent deterministic Day/Swing/Long change, no silent self-promotion and no confidence invented around missing evidence;
- G0–G16 remains the permanent release model.

## 2. Release-train / version-alignment rule

Starting with v18.9, DE.PULSE uses version numbers to communicate dependency phases while retaining small-patch delivery.

- **Major** = strategic maturity generation (`v19` infrastructure, `v20` adaptive intelligence).
- **Minor band** = coherent platform phase (`v19.0.x` control plane, `v19.1.x` hosted data plane/sync, etc.).
- **Patch** = one primary independently certifiable responsibility.
- Planned future numbers are alignment reservations only; exact identity freezes at G1.
- G0/G1 may split a broad packet and shift later unstarted reservations inside the same band.
- Shipped versions are immutable.
- Corrective/security work interrupts the planned train when needed; dependency truth is reconciled before resuming.

This produces a recognizable industry-style release train without turning minor versions into heavy bundles.

## 3. v18.9.0 — COMPLETE / IMMUTABLE STABLE

Issue #61 / `ADAPT-TRADEINSIGHT-001` is complete. Qualified source `9e86b5e731f7a585cc77c1521f3639fc7a208efc`, certified candidate `9ea81cddae4875ae15d3719ca028519a36c597b6`, Release #32 and `release/v18.9.0/stable-evidence-manifest.json` remain immutable authority.

Post-Stable findings do not rewrite v18.9.0; they enter #65 as small corrective/professional patches.

## 4. v18.9.x — Stabilize, instrument, validate, operationalize

Architecture audit reordered observability **before** broad TradeInsight capability admission so SHADOW/production decisions are measured rather than assumed.

1. **v18.9.1 — Runtime crash corrective ONLY** — #64; evidence-based macOS Apple Silicon SIGABRT diagnosis/fix, persisted-state/API-key continuity, lifecycle/relaunch regression, packaged runtime proof.
2. **v18.9.2 — TradeInsight Settings/API-key UX ONLY** — canonical Settings/local secret owner; masked Save/Test/Clear; truthful status; save preserves context.
3. **v18.9.3 — Coverage-aware Smart Provider Router core ONLY** — memory/DB reuse first, residual-gap acquisition, deterministic merge/provenance/coverage re-evaluation; lifecycle != serving role.
4. **v18.9.4 — Canonical company/instrument identity ONLY** — one identity owner reused by desks/Research/Discovery/Add Symbol.
5. **v18.9.5 — Market Data Modes + capability diagnostics ONLY** — behavior/quality modes, source/freshness/coverage diagnostics, no provider-brand mode owner.
6. **v18.9.6 — Provider observability / Adaptive telemetry foundation ONLY** — coverage, calls avoided, contribution/usefulness, freshness, rate/backpressure, disagreement, shared-work efficiency, runtime load and protected-session headroom.
7. **v18.9.7 — TradeInsight SEC Form 4 enrichment ONLY** — SHADOW-first; direct SEC/EDGAR authoritative; source-family de-duplication; measured using v18.9.6 telemetry.
8. **v18.9.8 — TradeInsight ticker/company search ONLY** — canonical symbol validation/identity fallback/corroboration.
9. **v18.9.9 — TradeInsight movers/ranking evidence ONLY** — SHADOW candidate evidence into canonical Opportunity Radar/ranker.
10. **v18.9.10 — Remaining useful TradeInsight capability admission ONLY** — every useful entitlement gets named consumer/owner/lifecycle/serving role/freshness/rights/retention/rate/Market Mode disposition.
11. **v18.9.11 — Session-Aware Data Readiness Maintenance ONLY** — one canonical light-overnight/heavy-weekend coordinator, protected pre-market/regular/after-hours reserves, preemption/checkpoint/resume and no blind refetch.
12. **v18.9.12 — v18.9.x Professional Closure ONLY** — no new feature scope; zero-miss/duplicate-owner audit, runtime/provider/persistence/session/package proof and Adaptive Intelligence Scorecard.

### v18.9.x exit contract

v18.9.12 must prove:
- current native runtime reliability;
- canonical settings/secrets behavior;
- coverage-aware persistence-first provider fulfillment;
- stable canonical identity;
- behavior-oriented Market Data Modes;
- telemetry sufficient to evaluate provider usefulness and operating cost;
- TradeInsight capabilities admitted only through canonical owners and lifecycle rules;
- session-aware maintenance that does not degrade protected live sessions;
- actual macOS Apple Silicon + Windows x64 package/runtime evidence;
- zero unexplained carry-forward, orphan useful capability or duplicate owner.

## 5. v19 — Professional Data Infrastructure + Hosted Account Platform

**Entry:** v18.9.12 PASS.

v19 has two equally important missions:
1. professionalize data quality, rights, provenance, reliability, economics and point-in-time evidence; and
2. deliver issue #66 / `ADAPT-HOSTED-SYNC-001`: one DE.PULSE account across macOS, Windows and hosted web with SQLite native edge/offline state, PostgreSQL shared authority and a zero-key hosted Provider Gateway.

### Target architecture

```text
macOS / Windows
  SQLite edge/offline working set
        │
        │ typed authenticated sync
        ▼
DE.PULSE hosted API / services
        │
        ├─ canonical identity/session/capability truth
        ├─ Smart Provider Router v2 + freshness/cache/persistence reuse
        ├─ rights/entitlement gate
        ├─ existing multi-feed subscription owner
        ├─ server-side managed provider secrets
        ▼
PostgreSQL shared account/state authority
        ▲
        │
 hosted web through same APIs
```

Commercial normal users authenticate only to DE.PULSE and never receive/configure DE.PULSE-owned provider keys.

### v19.0.x — Governance, Control Plane & Data Foundation

- **v19.0.0 — Provider Capability / Entitlement / Rights Registry.** Runtime-enforceable provider/capability rights, commercial/multi-user use, caching, redistribution, derived/AI use, attribution, limits and expiry.
- **v19.0.1 — Hosted Identity / Device / Session Control Plane.** Canonical roles/capabilities, account/device lifecycle, revocation, privileged re-authentication and server-side API/SSE/native/web truth.
- **v19.0.2 — PostgreSQL Tenancy / Schema / Pool / HA-PITR Foundation.** Shared authority, tenancy isolation, migrations, capacity, encrypted backup, restore/failover and RPO/RTO. No broad sync activation.
- **v19.0.3 — Managed Secrets / KMS Lifecycle.** Environment-isolated provider secrets, rotation/rollback/compromise recovery and redaction. Zero provider secret on commercial clients.
- **v19.0.4 — Provider Quality / Cost / Coverage / SLO Scorecards.** Measured capability SLO/error-budget/usefulness/cost evidence from v18.9 telemetry.
- **v19.0.5 — Data Reconciliation / Revision / Point-in-Time Quality.** Source independence, revisions, historical correctness and provenance.

**Phase rule:** security/rights/identity/recovery/observability foundations must pass before shared hosted data-plane activation.

### v19.1.x — Zero-Key Provider Data Plane & Native Sync Foundation

- **v19.1.0 — Authenticated Hosted Provider Gateway.** Existing Smart Provider Router v2 behind hosted boundary; cache/persistence reuse, bounded circuits/backpressure and server-side credentials.
- **v19.1.1 — Rights/Entitlement + Live Fan-Out Isolation.** Machine-enforced cache/persistence/REST/WebSocket/SSE serving boundaries; existing multi-feed owner reused; no entitlement leakage.
- **v19.1.2 — Sync Protocol Foundation.** Snapshot/high-watermark bootstrap, SQLite atomic outbox, idempotency, authoritative server sequence/change log, incremental pull, checkpoints, tombstones, compaction/retention, stale-device re-bootstrap and mixed-version negotiation.
- **v19.1.3 — macOS Preferences + Watchlist Pilot.** Offline/restart/reconnect/conflict/local-account/lost-device proof.
- **v19.1.4 — Desks / Workspaces Sync.** Versioned membership/configuration/delete/history semantics through the same transport.

### v19.2.x — Cross-Platform Account Parity & Hosted Assurance

- **v19.2.0 — Windows x64 Sync Parity.** Same protocol/security/account truth; no Windows fork.
- **v19.2.1 — Hosted Web Account Parity.** Same APIs/session/capability/PostgreSQL account state; browser cache non-authoritative.
- **v19.2.2 — Rights-Aware Research / Durable State Portability.** Only lawful, entitled, provenance-bound durable artifacts; live market truth remains freshness/provider/state-owned.
- **v19.2.3 — Multi-User Security / Cost / Abuse / Capacity Hardening.** Per-account attribution/quotas, fairness, edge limits, cost, throttling, load shedding and protected-session capacity.
- **v19.2.4 — #66 Hosted Sync / Gateway Assurance Closure.** Adversarial/failure/recovery matrix; Mac/Windows/web parity; zero-key secrets; entitlement isolation; mixed-version sync; DB/secret/provider-right recovery.

### v19.3.x — Professional Point-in-Time Evidence Substrate

- **v19.3.0 — Institutional / 13F Evidence Infrastructure.** Direct SEC truth, identity/mapping, amendments, filing lag, point-in-time holdings and outcomes.
- **v19.3.1 — Two-Sided Long/Short Thesis Evidence Substrate.** Point-in-time plans, validation order, MFE/MAE and lawful short/crowding evidence with explicit UNKNOWN.
- **v19.3.2 — AODR Candidate / Ranking / Outcome Lineage.** Candidate/rank/reason/transition/history, shared-ranking efficiency and surfaced-vs-missed outcomes.

### v19.4.x — Reliability, Economics & v20 Readiness

- **v19.4.0 — ADR-GDI Professional Reliability / Capacity.** SLO/degradation history, provider/DB/runtime scorecards, warm start, indexes/pools, load shedding, reserve sizing and maintenance economics.
- **v19.4.1 — Specialized/Paid Provider Gap Evaluation.** Add/replace only when measured evidence proves a material gap; same router/rights/persistence/session contracts.
- **v19.4.2 — v20 Research-Readiness Audit.** Point-in-time evidence/features/outcomes/provenance/rights/independence/leakage/reliability readiness.

### v19.5.0 — v19 Major Closure

No new feature scope. Principal Engineer + security/data/operational audit. Require #66 PASS, truthful rights/commercial posture, recovery/rollback drills, SLO/capacity evidence, supported runtime/package proof, zero unowned dataset/provider role and zero unresolved P0 architecture gap.

## 6. v20 — Adaptive Intelligence & Decision Research

**Entry:** v19.5.0 PASS.

The v20 audit moved model/prompt governance **before** broad adaptive-model rollout. Adaptive systems may only learn from the trustworthy point-in-time substrate and may not become a repair layer for unreliable data.

### v20.0.x — Adaptive Research Control & Governance
- **v20.0.0 — Adaptive Research Control Plane + Immutable Experiment Ledger.** Reproducible datasets/features/cohorts/lineage, leakage controls, promotion/rollback evidence.
- **v20.0.1 — Model/Prompt Governance + Champion/Challenger.** Model/prompt identity, independent evaluation, explainability, drift, approval and rollback before broad production influence.
- **v20.0.2 — Historical Analogues + Regime-Conditioned Outcomes.** Point-in-time analogue/outcome research.
- **v20.0.3 — Calibration / FP-FN / Miss / Contradiction / Drift.** Calibration and abstention controls.

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

## 7. Industry-strength release controls inside G0–G16

For hosted/security/data/adaptive scope, existing gates absorb the following controls rather than creating G17+:
- architecture decisions/owner map and threat/data classification at G2;
- API/schema/protocol compatibility, SLOs, observability, migration/rollback and failure-injection plan at G3;
- contract tests and bounded feature/kill-switch controls at G4;
- negative authorization/tenant/rights/secret/adaptive tests at G7;
- load/soak/capacity/chaos/failover/backpressure/protected-session proof at G8;
- cross-platform role-aware UX and direct-route/API denial at G9;
- production-readiness review and unresolved-P0 blocking at G10;
- immutable RC/full certification G11/G12;
- package/provenance/runtime proof G13/G14;
- controlled/canary promotion for hosted-risk changes where applicable at G15;
- implementation-miss, incident/metric learning and cleanup at G16.

## 8. Why this ordering is intentional

- **Observability precedes provider expansion** so SHADOW usefulness/promotion is measured.
- **Rights/identity/PostgreSQL/secrets precede hosted gateway/sync** so multi-user activation does not become a security/data-rights retrofit.
- **Gateway + sync foundation precedes Windows/web/research parity** so every client consumes one protocol and account truth.
- **Hosted assurance closes before wider evidence dependence** so v20 never learns from unproven synchronized state.
- **Point-in-time evidence precedes adaptive learning** so backtests/evaluation avoid look-ahead/revision leakage.
- **Model governance precedes broad adaptive rollout** so evaluation, rollback and promotion exist before intelligence can influence production.

## Exactly one next action

Diagnose #64 from complete macOS crash evidence or deterministic reproduction and freeze the narrow `v18.9.1` G1. Do not create `v18.9.2` or any v19 product implementation branch until `v18.9.1` is truthfully closed or proven external/non-product.
