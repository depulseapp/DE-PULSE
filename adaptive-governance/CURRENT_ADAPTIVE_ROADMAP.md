# CURRENT Adaptive Roadmap

**Canonical rebaseline:** `governance/V19_V20_REBASELINE.md`  
**Product-audit rebaseline:** `governance/PRODUCT_AUDIT_REBASELINE_2026_08_27.md`  
**Full audit coverage:** `governance/PRODUCT_AUDIT_COVERAGE_2026_08_27.md`  
**5/5 maturity target:** `governance/programs/V19-V20-REBASELINE/product-audit-5x5-target.json`  
**Audit finding register:** `governance/programs/V19-V20-REBASELINE/product-audit-finding-register.json`  
**Provider-registry additive rebaseline:** `governance/V19_V20_PROVIDER_REGISTRY_REBASELINE.md`  
**Provider-registry permanent contract:** `adaptive-governance/ADAPTIVE_PROVIDER_REGISTRY_CONTRACT.md`  
**Machine current-state authority:** `governance/current-state.json`  
**Backlog map:** `governance/programs/V19-V20-REBASELINE/backlog-version-matrix.json`  
**HOST map:** `governance/programs/V19-V20-REBASELINE/host-requirement-version-map.json`  
**Provider-registry machine map:** `governance/programs/V19-V20-REBASELINE/adaptive-provider-registry.json`  
**Legacy future-commitment conservation:** `governance/programs/V19-V20-REBASELINE/legacy-future-commitment-conservation.json`  
**Cross-integration map:** `governance/programs/V19-V20-REBASELINE/cross-integration-matrix.json`  
**Whole-product surface map:** `governance/programs/V19-V20-REBASELINE/whole-product-surface-rebaseline.json`  
**Certified Stable:** `v18.10.0` — immutable  
**Active version:** `v19.0.0` — Hosted Trust & Identity Foundation  
**Active work slice:** `ADAPT-HOSTED-TRUST-FOUNDATION-001`  
**Active issue/PR:** #148 / PR #149  
**Active branch:** `adapt-hosted-trust-foundation-001`  
**Canonical closure ledger:** `governance/work-slices/ADAPT-HOSTED-TRUST-FOUNDATION-001/closure.json`

## North Star after full-product audit

DE.PULSE becomes one market-intelligence operating system for U.S.-listed equities and approved ETFs:

`market observations -> rights/quality -> deterministic features/evidence -> SymbolIntelligenceSnapshot -> Opportunity Lifecycle -> Dashboard/MI/Desks/Discovery/Watchlist/Alerts/Research projections -> frozen Decision Brief -> outcomes -> governed adaptation`

Web, macOS and Windows are first-class clients of one domain truth. DE.PULSE remains decision support only and **No Execution** remains permanent.

The architecture target is a Go **modular monolith** with strong domain/application/provider/persistence/platform boundaries, tenant-aware Postgres + transactional outbox, versioned APIs/events and thin typed clients. Kafka/Kubernetes/microservices are not roadmap objectives; extraction occurs only after measured need.

## Product-audit Executive findings — roadmap conservation

All ten Executive findings are now roadmap requirements:

1. **Canonical boundary:** replace giant RuntimeSnapshot as domain authority with versioned `SymbolIntelligenceSnapshot`, typed Evidence/events/deltas and field-level freshness/provenance.
2. **Radar:** preserve detector value; move lifecycle authority into one shared Opportunity aggregate/state machine.
3. **Watchlist:** add first-class selected-universe intelligence projection with ranked attention, promotion/demotion transitions, explanations, contradictions, Research handoff, alerts and sync; never create another scorer/scanner.
4. **Adaptive:** point-in-time outcomes, censoring, model registry, time-split evaluation, shadow/champion-challenger, drift/approval/rollback before learned production behavior.
5. **Renderer authority:** golden-characterize and move technical/desk/side/geometry/scoring truth into versioned server-owned Go packages.
6. **Hosted trust:** finish tenant persistence, managed secrets, IaC/service trust, real PITR/recovery, audit, observability and operations.
7. **Desktop distribution:** secure storage, macOS Developer ID/notarization, Windows trusted signing/installer, update channels and rollback.
8. **Commercial data rights:** executable rights enforcement plus signed external evidence before applicable public/commercial activation; credentials never imply rights.
9. **Documentation:** one machine-led current truth; duplicate base/CURRENT narratives are consolidated/demoted and historical evidence remains immutable.
10. **Modular monolith:** preserve Go/core reuse and avoid premature distributed-system complexity.

The roadmap also conserves the full audit-risk register: stable instrument identity; point-in-time fundamentals/macros; exchange calendar/DST; clock skew/late events; raw/adjusted basis; OI as-of quality; alert causal dedupe; provider correlation; revision policy; privacy-vs-audit retention; cross-device conflicts; offline truth; AI egress; prompt injection; old-client schema compatibility; personalization/calibration separation; adaptive selection bias; censored outcomes; full-snapshot fanout cost; operational ownership; undefined Long King/Short King semantics; and planned-only Call/Put Wall semantics.

## 5/5 maturity objective

The audit scorecard becomes an evidence-backed target across all eleven domains. 5/5 is not a deadline label or commercial activation switch. Each domain closes only when all machine-target criteria are evidenced and no Critical/High gap contradicts the score.

Time-dependent maturity is earned rather than forced. Adaptive Intelligence does not reach 5/5 until point-in-time cohorts, walk-forward evaluation, selection-bias controls, drift/rollback and production-governance evidence exist.

## Surface roadmap

- **Dashboard — IMPROVE:** concise operating/attention summary only.
- **Market Intelligence — KEEP/IMPROVE:** shared market/regime/liquidity context.
- **Day/Swing/Long — KEEP / REWRITE DOMAIN BOUNDARY:** preserve horizon UX while server owns deterministic two-sided policy.
- **Discovery — KEEP / MERGE MODEL:** broad-universe projection of shared lifecycle.
- **Opportunity Radar — KEEP DETECTOR / CONSOLIDATE LIFECYCLE.**
- **Watchlist — ADD FIRST CLASS:** selected-universe projection of same lifecycle.
- **Research — KEEP / PROMOTE TO FROZEN DECISION BRIEF.**
- **AI Copilot — IMPROVE / CONSIDER CONSOLIDATION:** keep evidence-bounded synthesis; common workflows may move into Research/global command; never parallel truth.
- **Administration/Maintenance/Settings — KEEP ROLE-GATED.**
- **News/Earnings/Filings — MERGE SYMBOL PRESENTATION INTO RESEARCH** unless a separate event explorer later proves value.

## Current v19 sequence — rebaselined

1. `v19.0.0` — **Hosted Trust & Identity Foundation** — HOST-001..023 + core #164/#156 security/auth overlap. Current earliest OPEN dependency is HOST-010..012 real managed backup/PITR/operator recovery evidence through HOST-016. Product-audit architecture does not bypass this band.
2. `v19.1.0` — **Canonical Intelligence & Provider Foundation** — existing #150/#151/#153/#154/#155/#160/#167 core plus audit canonical-boundary/renderer-extraction foundations. Golden vectors; Observation/Evidence/Snapshot/Transition/DecisionBrief schemas; `SymbolIntelligenceSnapshot` compatibility adapter; server-owned technical/horizon package boundaries; Adaptive Provider Registry + Market Data first adoption through generic capability path.
3. `v19.2.0` — **Hosted Serving, Sync & Postgres v2** — HOST-024..039. Tenant-aware normalized persistence, RLS/isolation disposition, revisions/outbox/conflict/tombstones, user-scoped versioned APIs/deltas, hosted cache/fanout, cross-device/offline schema foundations.
4. `v19.3.0` — **Shared Opportunity Lifecycle & Cross-Platform Product Contract** — HOST-040..047/053 + #152/#156/#159/#160/#167/#171/#164 UX + `LEGACY-TRADER-SETUP-SHORT-001`. Shared Opportunity state machine; Radar/Rapid Move/Discovery adapters; server-owned two-sided deterministic setup; cross-platform auth/role/IA/API compatibility.
5. `v19.4.0` — **Market Intelligence, Research Decision Brief & Watchlist Foundation** — HOST-049 + #158/#161/#162/#171 + audit Watchlist requirement. Frozen-as-of Brief identity and first-class Watchlist selected-universe projection with ranked attention/explanation/contradiction/freshness/trust/Research handoff; no new scorer.
6. `v19.4.1` — **Discovery / Watchlist / Radar Convergence** — HOST-048 + #163/#171 + halt/LULD/pause/resume. Broad-universe Discovery and selected-universe Watchlist use one lifecycle; Radar remains detector; durable transition alerts/causal dedupe are reconciled.
7. `v19.5.0` — **Price/Volume & Event-Anchored Intelligence** — #168/#169. Canonical event identity/revisions, reaction context and incident correlation.
8. `v19.5.1` — **Options Structure & GEX Intelligence** — #157. Call/Put Wall only after expiry-aware quality/coverage/OI-as-of/rights semantics. No signed-dealer claim from gamma×OI. Long King/Short King stays blocked until an approved evidence/horizon/outcome definition exists.
9. `v19.6.0` — **Point-in-Time Evidence & Outcome-Ready Foundation** — HOST-057..064 + deterministic #165 + institutional/two-sided thesis substrate. Stable instrument identity, bitemporal/vintage facts, raw/adjusted lineage, explicit censoring, unbiased/control outcome sampling, feature snapshot IDs.
10. `v19.6.1` — **Hosted Reliability, Economics, Observability & 5/5 Readiness** — HOST-050..056/065..071 + ADR-GDI + provider-gap + #170/#171 reconciliation + registry scorecards + SLO/on-call/recovery/scale + desktop distribution readiness. Full audit/maturity residual review; domains below 5/5 stay explicit.
11. `v19.7.0` — **v19 Major Deterministic/Hosted Closure** — HOST-072; no feature scope. Zero unexplained audit/responsibility rows, compatibility migration evidence and exact-head Fast/Qualified/G0-G16. Commercial activation remains OFF unless Owner explicitly changes it.

### Two-sided deterministic Desk conservation

Current source is not a true SHORT plan: bearish labels can coexist with long-oriented target/invalidation/action geometry. v19.3 must provide one canonical contract:
- `LONG / SHORT / NO_SETUP-WAIT` side;
- direction separate from setup-quality score;
- LONG target above/invalidation below;
- SHORT cover/target below/invalidation above;
- positive direction-aware R/action-state/entry-distance/sort/chart/replay semantics;
- Research/Discovery/Watchlist consume same canonical side/geometry;
- strong bearish evidence does not force a SHORT setup when readiness is insufficient;
- No Execution permanent.

`Long King / Short King` is not automatically this contract. It needs its own formal evidence/user-purpose definition before product use.

## Conserved provider / Data Health program

The completed #80/#81/#82/#83/#78/#84 chain remains permanent across versions:
- executable provider/capability baseline, SLOs, ownership and fetch-path classification;
- Smart Provider Router v2 for routable capabilities while direct authority remains explicit;
- canonical Data Health evaluation/fallback/recovery/hysteresis/workload protection;
- lifecycle `SHADOW -> VALIDATED -> APPROVED -> PRODUCTION` with no automatic authority promotion;
- TradeInsight/Market Data are provider-adoption cases, not alternate routers/health/caches/lifecycles;
- zero-gap fault/native/professional evidence.

Truthful unresolved required evidence remains visible at the smallest valid scope as **PARTIAL COVERAGE** or **DATA DEGRADED**. These states must not be suppressed or normalized away merely to make the product appear healthy; valid cache/fallback and automatic recovery may clear them only when canonical evidence is actually restored.

New audit additions:
- provider corroboration includes upstream-independence/correlation truth;
- corrections use revision/supersedes lineage;
- provider observation/publication/filing time remains evidence time;
- rights coordinates remain explicit for display/cache/persistence/derived/AI/multi-user/redistribution.

## Adaptive Provider Registry / future-proof onboarding

Permanent direction:

`Provider Adapter -> Adaptive Provider Registry -> capability/entitlement probes -> rights/authority + Data Health -> Smart Provider Router v2 -> canonical state -> all useful consumers`

Provider-specific implementation ends at adapter; consumers request capabilities, not provider names. Technical eligibility/fallback/recovery may change automatically from observed capability; authority lifecycle and public/commercial rights never auto-promote.

Market Data remains the first reusable Registry adopter in v19.1, but it must fit the canonical Snapshot/provider boundary and cannot become a one-off architecture.

## Current v20 sequence — rebaselined

1. `v20.0.0` — **Outcome Learning & Adaptive Control Plane** — point-in-time evaluation, model/policy registry, sample floors, shadow/champion-challenger, drift/approval/rollback.
2. `v20.1.0` — **Adaptive Chart Pattern & Similarity Intelligence** — structured split-safe features and interpretable baselines/nearest-neighbor/tree/linear challengers before deep/image models.
3. `v20.2.0` — **Adaptive Market Synthesis, Regime & Discovery Learning** — conserved ASBI normalization/synthesis/contradiction/abstention/outcome behavior.
4. `v20.3.0` — **Adaptive Institutional & Two-Sided Thesis Intelligence**.
5. `v20.3.1` — **AODR Adaptive Opportunity Intelligence**.
6. `v20.4.0` — **Agent Orchestration & Controlled MCP/API**.
7. `v20.5.0` — **Adaptive Operations** — bounded provider utility/cost/reliability priors inside Smart Provider Router v2; no parallel router/authority/rights promotion.
8. `v20.6.0` — **Professional Adaptive Closure**; no feature scope.

## Permanent temporal/adaptive direction

Every material fact gets an explicit disposition for source/observed/ingested/effective/as-of/expiry/session/revision/provider/dataset/rights/quality truth. Historical evaluation uses point-in-time joins; revised facts do not rewrite history silently; raw/adjusted price basis is explicit; OI as-of is separate from live option quote/IV time; censored outcomes are first-class.

Adaptive authority remains:
`deterministic rules -> registered statistical challenger -> shadow/evaluation -> approved bounded production influence`, with LLM synthesis last and non-authoritative for canonical calculations/routing/rights/lifecycle.

## Scale and client roadmap

Near-term architecture must survive roughly 100 users without rights/full-snapshot/local-secret/unsigned-update failures; 1,000 users without duplicate-ingestion/JSON-lock/session/alert-fanout collapse; and 10,000 users with partitioning/event queues/DB hotspots/adaptive selection bias/on-call load managed. Postgres + outbox + stateless API/workers remains default until measured service extraction need.

Core owns secrets/rights/ingestion/canonical state/intelligence/opportunities/briefs/outcomes/adaptive/persistence/recovery/audit. Clients own interaction/rendering/charts/accessibility/typed API validation. Desktop edge may keep encrypted authorized last-known cache and OS notifications/deep links but never simulate live confidence.

## Permanent intelligence direction

`canonical evidence -> deterministic intelligence -> shared Opportunity Lifecycle -> product projections -> point-in-time outcomes -> bounded adaptive learning -> optional AI/agent explanation/orchestration`

#170 remains mandatory cross-integration/Market-Regime gate. #171 remains mandatory UI/data-density/intelligence-maturity audit. The product-audit finding register is now an equally mandatory cross-version conservation family.

## Current exact state / next action

The active release remains `v19.0.0` on #148 / PR #149. Finish the audit-governance rebaseline with fresh exact-head Fast, then return to the machine-current earliest open dependency: **HOST-010..012 real managed PostgreSQL backup/PITR/operator recovery deletion-retention evidence through HOST-016**. Do not start v19.1, Watchlist, Opportunity Lifecycle or HOST-013+ while that dependency remains unresolved unless governance explicitly reclassifies it.
