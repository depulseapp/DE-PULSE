# DE.PULSE Full Product Audit Coverage Reconciliation — 2026-08-27

**Audit head:** `87baad3a90357b9af8e040a5ff34203a269fd861`  
**Stable comparison:** `v18.10.0`  
**Live branch at rebaseline start:** `adapt-hosted-trust-foundation-001`  
**Live head at rebaseline start:** `d49e54c7361fb4478ad8340bc4d98fe8f6c323f1`  
**Active version / PR:** `v19.0.0` / PR #149  
**Lifecycle:** `DEVELOPMENT`  
**Commercial activation:** OFF unless the Owner explicitly activates it.

This file is the human-readable audit coverage companion to:
- `governance/PRODUCT_AUDIT_REBASELINE_2026_08_27.md`;
- `governance/programs/V19-V20-REBASELINE/product-audit-5x5-target.json`;
- `governance/programs/V19-V20-REBASELINE/product-audit-finding-register.json`.

It exists to prove that the rebaseline did not stop at the Executive Assessment or maturity table. The full 54-page audit, including major gaps, P0/P1/P2/P3 backlog, 30/90-day plans, long-term strategy, Things We Are Not Thinking About Yet, red-team scale, adaptive operating model, temporal/event/pattern/AI architecture, cloud/client boundary, repository reorganization and five-year risks is conserved below.

## 1. Reconciliation rule

The audit is an evidence-based architecture/product baseline, not a frozen statement of live implementation forever. Its CURRENT/GAP/TARGET/UNVERIFIED labels are reconciled against newer source and executable evidence.

Since the audit head, the branch added materially stronger HOST-012 managed-recovery evidence machinery and the 5/5 audit rebaseline overlay. Those additions improve readiness but do not erase the audit's architectural gaps. In particular:
- real managed backup/PITR/operator recovery evidence remains external and OPEN;
- `SymbolIntelligenceSnapshot` is still a target, not current authority;
- a shared Opportunity Lifecycle is still a target;
- first-class Watchlist intelligence is still not implemented;
- renderer-owned authoritative desk/technical calculations remain a critical extraction target;
- Postgres v2 tenant isolation/outbox/normalized workspace remains target architecture;
- Web remains a future first-class client;
- adaptive learned production behavior remains intentionally unimplemented until point-in-time evidence is sufficient.

## 2. Ten Executive Assessment findings — all mandatory

| ID | Finding | Rebaseline status | Mandatory disposition |
|---|---|---|---|
| AUDIT-EXEC-01 | Shared intelligence exists, but the `Engine.Snapshot()`/`RuntimeSnapshot` boundary is wrong | OPEN | Versioned `SymbolIntelligenceSnapshot`, typed Evidence/events/deltas, field-level provenance/freshness; compatibility migration only |
| AUDIT-EXEC-02 | Opportunity Radar is real and reusable | OPEN | Preserve detector; feed one shared Opportunity Lifecycle; no cloned Watchlist scorer/lifecycle |
| AUDIT-EXEC-03 | Watchlists are membership containers, not ranked intelligence | OPEN | First-class Watchlist projection with shared lifecycle, explanations, transitions, Research handoff, alerts and sync |
| AUDIT-EXEC-04 | Adaptive evidence/outcomes exist but no learned production behavior | OPEN BY DESIGN | Point-in-time corpus, censoring, model registry, shadow/challenger, drift, approval and rollback before learned influence |
| AUDIT-EXEC-05 | Renderer owns material authoritative domain logic | CRITICAL OPEN | Golden characterization first; extract technical/horizon/side/geometry/scoring authority into versioned Go packages |
| AUDIT-EXEC-06 | Hosted trust is advanced but not production-ready | PARTIAL ACTIVE | Finish tenant persistence, secrets, IaC/service trust, real PITR/recovery, audit and hosted operations |
| AUDIT-EXEC-07 | Desktop assurance exceeds distribution maturity | OPEN | macOS Developer ID/notarization; Windows signing/installer; updater/channel/rollback; secure storage/deep links |
| AUDIT-EXEC-08 | Commercial data rights are a P0 public/commercial gate | CONTROL EXISTS / APPROVAL EXTERNAL | Maintain executable rights enforcement; obtain signed rights before relevant activation; credentials never imply rights |
| AUDIT-EXEC-09 | Documentation/current-state truth is duplicated and drift-prone | PARTIAL ACTIVE | Machine state precedence, canonical docs, generated current status, archive/demote duplicate narratives |
| AUDIT-EXEC-10 | Modular monolith is the correct next architecture | ADOPTED TARGET | Go modular monolith + Postgres/outbox + typed APIs/clients; no Kafka/Kubernetes/microservice program without measured need |

No Executive finding may be dropped because another finding appears to overlap it.

## 3. Maturity 5/5 objective

All eleven audit scorecard domains are governed by `product-audit-5x5-target.json`. The score is evidence, not aspiration. A domain cannot reach 5/5 while a contradictory Critical/High gap remains open.

The target is **architecture capable of 5/5 now, maturity earned from evidence over time**. Adaptive Intelligence in particular cannot be truthfully accelerated to 5/5 merely by adding a learner; it needs adequate point-in-time cohorts, walk-forward evaluation, selection-bias controls, drift monitoring and reversible production governance.

## 4. Audit-wide risks that must not be lost

The audit's “Things We Are Not Thinking About Yet” list is promoted into mandatory conservation rather than treated as optional ideas:

1. stable instrument/listing identity beyond ticker;
2. point-in-time fundamentals and revised macro vintages;
3. exchange calendar, DST and half-day truth;
4. clock skew and late/out-of-order events;
5. raw versus adjusted price basis;
6. options OI/Greeks quality and no unsupported dealer-sign inference;
7. alert fatigue, causal incident correlation and dedupe;
8. correlated providers/upstream-dependency truth;
9. revision/supersedes policy for corrected bars/events/filings/earnings;
10. privacy deletion versus immutable audit retention/pseudonymization;
11. cross-device Watchlist/preference conflict resolution;
12. explicit offline/last-known semantics;
13. AI egress residency/provider-right implications;
14. prompt injection through news/filings/social/community inputs;
15. API/event compatibility across old desktop clients and forced-upgrade policy;
16. personalization must not silently mutate shared calibration truth;
17. adaptive selection bias from hydrating only promoted symbols;
18. censored outcomes for halts/missing bars/delistings/gaps;
19. bandwidth/data-exposure cost of full snapshot fanout;
20. named operational ownership/on-call/escalation/provider incident playbooks.

Two product-specific unresolved semantics are added because they were explicitly surfaced by the audit and the intended Watchlist design:
21. **Long King / Short King** have no formal current definition. Do not build labels/columns until a deterministic evidence/horizon/outcome contract is approved.
22. **Call Wall / Put Wall** are planned terms, not current product semantics. Define expiry-aware quality-gated cluster semantics, OI as-of truth and rights first; do not imply signed dealer positioning.

Machine tracking: `product-audit-finding-register.json`.

## 5. Surface ownership / information architecture

Rebaseline dispositions from the tab-by-tab audit are mandatory:
- **Dashboard — IMPROVE:** operating summary only; do not become second Market Intelligence/Discovery.
- **Market Intelligence — KEEP/IMPROVE:** canonical market/regime/liquidity context consumed by symbol decisions.
- **Day/Swing/Long — KEEP / REWRITE DOMAIN BOUNDARY:** preserve user horizon/workflow; server owns deterministic policy and two-sided geometry.
- **Discovery — KEEP / MERGE PRODUCT MODEL:** broad-universe projection of shared lifecycle.
- **Opportunity Radar — KEEP DETECTOR / CONSOLIDATE LIFECYCLE:** not a competing product scorer/state owner.
- **Watchlist — ADD FIRST CLASS:** user-selected-universe projection of the same lifecycle; never a second scanner.
- **Research — KEEP / PROMOTE TO DECISION BRIEF:** immutable-as-of explanation destination.
- **AI Copilot — IMPROVE / CONSIDER CONSOLIDATION:** preserve bounded synthesis service; Research/global command may absorb common workflows; AI never becomes market truth.
- **Administration/Maintenance/Settings — KEEP ROLE-GATED:** operator machinery/secrets stay out of ordinary-user experience.
- **News/Earnings/Filings — MERGE PRESENTATION:** preserve services/evidence, use Research for symbol-specific detail; optional event explorer only when proven useful.

## 6. Canonical intelligence and Watchlist architecture

Canonical target flow:

`Market observations -> rights/quality -> deterministic features + evidence graph -> versioned SymbolIntelligenceSnapshot -> rules/approved models -> Opportunity Lifecycle -> product projections -> frozen Decision Brief -> outcomes -> governed adaptation`

Watchlist rule:
- membership/user intent changes the eligible universe;
- shared intelligence changes attention priority;
- promotion/demotion means **attention priority changed**, not BUY/SELL;
- every transition explains why, contradictions, freshness and trust;
- Rapid Move, RVOL/volatility, patterns, levels/Fibonacci, catalysts, earnings/SEC/insiders, future Call/Put Wall and formally defined Long/Short King semantics all enter through canonical evidence/features, not Watchlist-specific calculations;
- Discovery and Watchlist differ by universe/intent, not by scoring engines.

## 7. Temporal / point-in-time contract

Every material fact/event must be able to carry, where applicable:
`source_at`, `observed_at`, `ingested_at`, `effective_from`, `effective_to`, `as_of`, expiration/half-life, market session/exchange timezone, revision/supersedes, provider/dataset, rights version and quality.

Hard truth rules:
- retrieval/cache time is not evidence time;
- unknown observation/source time stays unknown;
- fundamentals/macros require vintage/as-of truth;
- corrected facts create revisions rather than silently rewriting historical belief;
- raw/adjusted price basis is explicit;
- options live quote/IV time is distinct from OI as-of date;
- learning/outcomes use strict point-in-time joins.

## 8. Adaptive / pattern / AI control

Three authorities remain separate:
1. **Deterministic rules** calculate market/technical facts, safety/rights/quality and versioned lifecycle policy.
2. **Statistical/adaptive models** estimate registered relationships; they cannot enter production without evaluation/approval or bypass safety/rights floors.
3. **LLM synthesis** explains structured evidence/contradictions with evidence IDs; it does not calculate canonical indicators/option exposure/weights/routing/rights/lifecycle transitions.

Controlled adaptive loop:
`freeze feature/evidence/regime/rights/policy snapshot -> decision/transition -> explicit outcome -> point-in-time evaluation -> time-split challenger -> shadow compare -> drift/subgroup checks -> sample/metric gates -> human approval -> gradual production -> automatic rollback threshold -> historical reproducibility`.

Pattern learning begins with structured split-safe multi-timeframe features and simple interpretable baselines/nearest-neighbor/tree/linear challengers. Chart-image/deep sequence models are later only if they beat simpler baselines out of sample.

## 9. Provider/data architecture

Permanent requirements:
- provider-specific code ends at adapter/capability boundary;
- Smart Provider Router v2 remains sole general route/admission owner;
- direct SEC/EDGAR authority remains explicit;
- rights, capability/entitlement, freshness, health, latency, quota/cost, fallback/recovery and usefulness are separate dimensions;
- provider agreement must not overstate independence when upstream feeds correlate;
- one SEC scheduler/global fair-access budget; event/accession dedupe;
- repeated low-value polling moves toward release/event-triggered work where possible;
- options/raw payload retention is quality/cost governed;
- dataset utility, consumer, materiality, freshness, retention, point-in-time and rights remain explicit.

## 10. Persistence and hosted architecture

Target hosted persistence is not the current blob-compatible adapter. It is Postgres v2 with tenant-owned/scoped keys, normalized identity/workspaces/Watchlists, RLS/isolation disposition, foreign keys, revisions, outbox, audit, retention/partitioning, conflict/tombstone semantics and real backup/PITR/restore/DR evidence.

SQLite remains appropriate for local/desktop bridge/cache where its contract is truthful. Desktop offline data must be encrypted/authorized and explicitly last-known/degraded.

Near-term hosted processing uses Postgres transactional outbox + worker leases. Kafka/service extraction is allowed only after measured outbox throughput/replay/independent scaling/team ownership proves the need.

## 11. Security, distribution and operations

Do not call hosted/public distribution ready until applicable evidence covers:
- managed hosted secrets/KMS and OS secure desktop storage;
- external identity/OIDC/PKCE/deep-link/device/session contract;
- tenant/cache/stream negative isolation;
- durable admin/config/rights/security audit with privacy-aware redaction/pseudonymization;
- HTTP timeout/body/header/drain hardening and hosted CSP/egress policy;
- dependency/SCA/secret-history/supply-chain assurance;
- macOS Developer ID/hardened runtime/notarization/stapling;
- Windows trusted signing and installer/update strategy;
- signed update manifests, staged channels and rollback;
- external telemetry/tracing/crash/health pipeline;
- SLO/error budgets, capacity/cost budgets, runbooks, support and on-call.

Commercial/public activation remains a separate explicit Owner decision even when technical readiness reaches 5/5.

## 12. Scale red-team conservation

Do not overengineer, but design boundaries that survive scale:
- **~100 users:** rights, full-snapshot/SSE cost, local secrets, manual support and unsigned updates become first pain points.
- **~1,000 users:** provider economics/rate limits, duplicate ingestion, JSON identity/workspace locking, sessions and alert fanout dominate.
- **~10,000 users:** licensing economics, time-series partitions, event backlogs, DB hotspots, adaptive selection bias and operational load dominate.

The approved safe path is modular monolith + Postgres + outbox + stateless APIs/workers + managed deployment and measured SLOs, not speculative Kafka/Kubernetes.

## 13. Repository/client boundary

Target responsibility split:
- **Core/cloud:** identity/entitlements/tenant policy, provider secrets/rights/shared ingestion, canonical observations/events/snapshots, deterministic intelligence/opportunities/briefs, alerts/outcomes/adaptive jobs, Postgres/outbox/retention/recovery/audit, AI gateway/eval versions.
- **Web/macOS/Windows clients:** UI/tables/charts/interactions, typed API/event validation, deterministic server-output rendering, local view models/accessibility, preferences/commands, citation navigation.
- **Desktop edge:** encrypted authorized cache, OS notifications/deep links, device credential in OS secure storage, explicit last-known read-only degradation, redacted export/support bundle.

Incremental target organization remains `/apps`, `/internal/domain|application|providers|persistence|platform`, `/packages/contracts|ui`, `/data`, `/tests`, `/docs`, `/governance`, `/release`. Move files only when dependency direction improves and characterization evidence protects behavior.

## 14. ADR conservation

All seventeen audit ADR decisions are required planning artifacts, whether implemented as individual ADR files or consolidated decisions with equivalent durable IDs:
1. Canonical Symbol Intelligence
2. Shared Opportunity Lifecycle
3. Server-Owned Deterministic Policy
4. Hosted Modular Monolith
5. Provider Capability Gateway
6. Postgres v2 / Tenant Isolation
7. Web/Desktop Client Contract
8. Desktop Framework
9. Hosted Authentication
10. Secrets and Key Management
11. Alert Delivery
12. Adaptive Control Plane
13. Point-in-Time Data
14. AI Model Gateway and Evals
15. Desktop Distribution
16. Observability and Privacy
17. Commercial Data Rights

## 15. Rebaselined version intent

The active v19.0 work is not restarted or discarded. The audit overlays dependency-correct architecture around it:
- **Immediate governance:** conserve all audit findings/risks, 5/5 criteria, surface dispositions, ADR decisions and compatibility strategy.
- **v19.0:** finish Hosted Trust & Identity technical foundation; current earliest OPEN dependency remains HOST-010..012 and real managed recovery evidence.
- **v19.1:** canonical intelligence foundation — golden vectors, server-owned domain boundaries, Observation/Evidence/Snapshot/Transition contracts, SymbolIntelligenceSnapshot compatibility path; existing provider-registry/Market Data scope must fit this architecture rather than become a one-off.
- **v19.2:** hosted serving/sync + tenant-aware Postgres v2/outbox/revision/conflict/tombstone + versioned user-scoped APIs/deltas.
- **v19.3:** shared Opportunity Lifecycle, cross-platform product/auth/role contract and two-sided deterministic setup authority; Radar/Rapid Move/Discovery become evidence producers/projections.
- **v19.4/v19.4.1:** Research Decision Brief + first-class Watchlist + Discovery convergence; no second scorer.
- **v19.5.x:** price/volume/event/options enrichment through canonical evidence; formal options walls only after quality/rights semantics.
- **v19.6.x:** point-in-time/outcomes, reliability/economics, observability/recovery, adaptive readiness and full 5/5 residual review.
- **v19.7:** deterministic/hosted major closure only after zero unexplained audit/responsibility rows.
- **v20:** governed adaptation only after point-in-time outcome quality is adequate.

## 16. Exact execution order

1. Rebaseline truth/enforcement and correct machine-state drift.
2. Finish the currently active v19.0 dependency band without bypassing real infrastructure evidence.
3. Secure rights/secrets boundaries appropriate to the next target and keep commercial activation separate.
4. Freeze current behavior with golden characterization vectors.
5. Extract canonical domain contracts/server authority.
6. Introduce shared Opportunity Lifecycle in shadow/dual-read mode.
7. Build tenant-aware persistence/outbox and user-scoped versioned APIs/events.
8. Ship Watchlist as a lifecycle projection and frozen-as-of Research Brief.
9. Establish shared Web/macOS/Windows platform foundation, secure distribution and hosted observability.
10. Launch Web on the same domain truth; only then advance statistically proven reversible adaptation.

## 17. Zero-miss completion rule

G10/G16 and any 5/5 claim must reconcile:
- all 180 certified v18 responsibilities;
- HOST-001..072;
- mapped v19/v20 backlog and legacy commitments;
- all ten Executive findings;
- all audit-wide risk rows;
- surface dispositions;
- temporal/point-in-time rules;
- adaptive/AI authority rules;
- provider/data utility/rights/freshness rules;
- persistence/security/distribution/operations requirements;
- scale and repository/client-boundary constraints;
- all ADR decisions;
- every applicable Mac/Windows/Web responsibility;
- positive/adverse/recovery/load/compatibility evidence and durable regression ownership.

Anything not evidenced remains `OPEN`, `PARTIAL`, `UNVERIFIED` or `EXTERNAL_BLOCKED`. It is never converted to complete by roadmap text.
