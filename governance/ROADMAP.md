# DE.PULSE — Canonical Adaptive Roadmap

**Status:** ACTIVE / AUTHORITATIVE FOR PRODUCT PLACEMENT  
**Rebaselined:** 2026-08-28 from the full-product audit  
**Current execution status:** `governance/current-state.json` + active closure ledger + `handoff/CURRENT.md`

This file owns durable product sequencing. It deliberately does not duplicate a live branch SHA or next action; those change faster than roadmap intent.

## 1. Canonical inputs

The roadmap conserves:

- `governance/APPROVED-SCOPE.md` and the permanent contracts;
- certified v18 responsibilities and immutable Stable evidence;
- `governance/V19_V20_REBASELINE.md` and its machine maps;
- `governance/PRODUCT_AUDIT_REBASELINE_2026_08_27.md`;
- `governance/PRODUCT_AUDIT_COVERAGE_2026_08_27.md`;
- `governance/programs/V19-V20-REBASELINE/product-audit-finding-register.json`;
- `governance/programs/V19-V20-REBASELINE/product-audit-5x5-target.json`;
- backlog, HOST, legacy-commitment, cross-integration and whole-product-surface maps under `governance/programs/V19-V20-REBASELINE/`;
- `governance/programs/ADAPT-HOSTED-SYNC-001/requirement-conservation.json` and the **Zero-Miss Future-Version Conservation** rule;
- approved material decisions in `governance/DECISION-LOG.md`.

Documentation never proves implementation. Code/runtime/package evidence defines CURRENT; this roadmap defines TARGET and dependency placement.

## 2. North star

DE.PULSE becomes one market-intelligence operating system for U.S.-listed equities and approved U.S. ETFs:

`observations -> rights/quality -> deterministic evidence -> SymbolIntelligenceSnapshot -> Opportunity Lifecycle -> product projections -> frozen Decision Brief -> outcomes -> governed adaptation`

Web, macOS Apple Silicon and Windows x64 consume the same domain truth. DE.PULSE remains decision support only; No Execution and hidden Data Engine internals remain permanent.

## 3. Audit rebaseline decisions

The full-product audit established ten mandatory directions:

1. Replace giant `RuntimeSnapshot` as the domain boundary with versioned symbol intelligence, typed evidence/events and deltas.
2. Keep Opportunity Radar as a detector while one shared Opportunity Lifecycle owns state.
3. Add first-class Watchlist as a selected-universe projection, never another scorer/scanner.
4. Build point-in-time outcomes and controlled challenger/shadow learning before learned production influence.
5. Golden-characterize and move renderer-owned technical/desk/side/geometry/scoring authority into Go domain owners.
6. Finish hosted tenant persistence, managed secrets, service trust/IaC, recovery/PITR, audit and operations.
7. Finish secure macOS/Windows distribution, updates, rollback and OS credential storage.
8. Preserve executable provider-rights controls and require actual rights evidence before Commercial/Public activation.
9. Keep one machine-led current truth and one canonical narrative per adaptive layer.
10. Evolve as a Go modular monolith with Postgres/outbox/versioned APIs/thin clients; extract services only from measured need.

The complete risk register remains mandatory, including instrument identity, bitemporal/vintage facts, exchange calendars/DST, clock skew/late events, raw/adjusted basis, OI as-of quality, provider correlation/revisions, alert causal dedupe, privacy-vs-audit retention, cross-device conflicts, offline truth, AI egress/prompt injection, client-schema compatibility, personalization separation, adaptive selection bias, censored outcomes, snapshot fanout and operational ownership.

## 4. Product surface model

| Surface | Roadmap purpose | Disposition |
|---|---|---|
| Dashboard | Concise attention and operating summary | IMPROVE |
| Market Intelligence | Shared regime, liquidity, macro and market context | KEEP / IMPROVE |
| Day / Swing / Long | Horizon projections of shared evidence | KEEP UX / MOVE DOMAIN AUTHORITY |
| Discovery | Broad-universe opportunity projection | KEEP / MERGE MODEL |
| Opportunity Radar | Detection and evidence production | KEEP DETECTOR / CONSOLIDATE LIFECYCLE |
| Watchlist | Selected-universe opportunity projection | ADD FIRST CLASS |
| Research | Frozen-as-of Decision Brief and deep evidence | KEEP / PROMOTE |
| Alerts | Material lifecycle transition/incident delivery | IMPROVE / NO RESCORING |
| AI | Evidence-bounded synthesis and research interaction | IMPROVE / NO PARALLEL TRUTH |
| Admin / Maintenance / Settings | Role-gated operation and configuration | KEEP ROLE-GATED |
| News / Earnings / Filings | Symbol evidence and catalyst history | CONSOLIDATE PRESENTATION INTO RESEARCH |

Every feature must answer a distinct user question. Presentation may consolidate without deleting useful evidence or domain services.

## 5. Shared opportunity and Watchlist direction

One lifecycle serves Discovery, Watchlist, Radar, Alerts, Research and Desks:

`DETECTED -> OBSERVED -> QUALIFIED -> PROMOTED -> HIGH_PRIORITY -> COOLING -> DEMOTED -> RESOLVED`

- Discovery evaluates the approved broad universe.
- Watchlist evaluates only user-selected symbols.
- Radar/Rapid Move contribute evidence and triggers.
- Alerts deliver material transitions with causal dedupe.
- Research opens the same frozen snapshot/transition as a Decision Brief.

Promotion/demotion uses multiple evidence families, temporal decay, contradictions, quality/freshness and market context. A numeric rank is never sufficient without explanation and lineage.

`Long King / Short King` remain undefined until an approved evidence/horizon/outcome contract exists. `Call Wall / Put Wall` remain planned until expiry, coverage, OI-as-of, cluster, quality and rights semantics are formalized.

## 6. v19 — deterministic hosted product foundation

### v19.0.0 — Hosted Trust & Identity Foundation

HOST-001..023 plus the applicable core security/auth scope in #164/#156. Close technical Development Production Ready evidence for provider-rights controls, tenant identity/session/device/MFA, product entitlement, privacy lifecycle, managed environment/service trust, tenant persistence/recovery, secrets, supply chain, provider scorecards and point-in-time truth. Commercial/Public activation remains separate and OFF.

### v19.1.0 — Canonical Intelligence & Provider Foundation

Existing #150/#151/#153/#154/#155/#160/#167 core plus audit canonical-boundary/renderer-extraction foundations:

- freeze golden vectors for `computePlan`, technical state, Rapid Move, Radar and Research;
- define Observation, Evidence, `SymbolIntelligenceSnapshot`, Transition and DecisionBrief schemas;
- introduce snapshot compatibility behind current consumers;
- establish server-owned technical/horizon package boundaries;
- adopt the Adaptive Provider Registry and Market Data through the generic capability/Router/Data Health path.

### v19.2.0 — Hosted Serving, Sync & Postgres v2

HOST-024..039. Add tenant-aware normalized Postgres v2, explicit RLS/isolation disposition, revisions, transactional outbox, conflict/tombstone policy, user-scoped versioned APIs/deltas, lawful hosted fanout and cross-device/offline foundations.

### v19.3.0 — Shared Opportunity Lifecycle & Cross-Platform Contract

HOST-040..047/053, #152/#156/#159/#160/#167/#171/#164 UX and `LEGACY-TRADER-SETUP-SHORT-001`:

- introduce the shared Opportunity aggregate/state machine behind shadow/dual-read comparison;
- adapt Radar/Rapid Move/Discovery instead of duplicating them;
- move two-sided deterministic setup side/geometry/policy to the server;
- establish shared authentication, roles, information architecture, API/event compatibility and minimum-client rules.

### v19.4.0 — Market Intelligence, Research Brief & Watchlist Foundation

HOST-049, #158/#161/#162/#171 and audit Watchlist scope:

- frozen-as-of Research Decision Brief identity;
- first-class Watchlist selected-universe projection;
- ranked attention, promotion/demotion explanation, contradictions, confidence, freshness and Research handoff;
- no Watchlist scorer, provider loop or duplicate persistence owner.

### v19.4.1 — Discovery / Watchlist / Radar Convergence

HOST-048, #163/#171 and conserved halt/LULD/pause/resume behavior. Prove both universe projections use one lifecycle, Radar remains a detector, and lifecycle transitions/alerts are durable and causally deduplicated.

### v19.5.0 — Price/Volume & Event-Anchored Intelligence

#168/#169. Add canonical event identity, revisions, temporal reaction context and incident correlation through shared snapshot/lifecycle/Brief owners.

### v19.5.1 — Options Structure & GEX Intelligence

#157. Formalize expiry-aware coverage, OI-as-of, strike clusters, quality and rights before Call/Put Wall. Do not infer signed dealer positioning from gamma multiplied by OI.

### v19.6.0 — Point-in-Time Evidence & Outcome-Ready Foundation

HOST-057..064, deterministic #165 and institutional/two-sided substrate. Add stable instrument identity, bitemporal/vintage facts, raw/adjusted lineage, feature snapshot IDs, explicit censoring and unbiased/control outcome sampling.

### v19.6.1 — Reliability, Economics, Observability & 5/5 Readiness

HOST-050..056/065..071, ADR-GDI, provider-gap, #170/#171 reconciliation, provider reliability/economics scorecards, external SLO/on-call/recovery/scale evidence and desktop distribution readiness. Review every 5/5 maturity residual without score inflation.

### v19.7.0 — v19 Major Closure

HOST-072; no feature scope. Zero unexplained audit/responsibility rows, compatibility migrations reconciled, exact-head G0–G16 evidence and Commercial/Public activation still OFF unless separately authorized.

## 7. v20 — governed adaptive intelligence

1. **v20.0.0 — Outcome Learning & Adaptive Control Plane:** point-in-time evaluation, model/policy registry, sample floors, challenger/shadow, drift, approval and rollback.
2. **v20.1.0 — Pattern & Similarity Intelligence:** structured split-safe features and interpretable baselines before deep/image models.
3. **v20.2.0 — Adaptive Market Synthesis, Regime & Discovery Learning:** conserved ASBI normalization, contradiction, abstention and outcomes.
4. **v20.3.0 — Adaptive Institutional & Two-Sided Thesis Intelligence.**
5. **v20.3.1 — AODR Adaptive Opportunity Intelligence.**
6. **v20.4.0 — Agent Orchestration & Controlled MCP/API:** evidence-scoped, rights-aware and auditable.
7. **v20.5.0 — Adaptive Operations:** bounded provider utility/cost/reliability priors inside Smart Provider Router v2; no parallel router or automatic authority/rights promotion.
8. **v20.6.0 — Professional Adaptive Closure:** no feature scope.

## 8. Platform and cloud boundary

Core/cloud owns provider secrets/rights, shared ingestion, canonical intelligence, lifecycle/Briefs, tenant persistence, alert delivery, outcomes/adaptive jobs, recovery/audit and AI gateway policy.

Web/macOS/Windows clients own interaction, rendering, charts, accessibility and typed schema validation. Desktop may use encrypted authorized last-known cache, OS deep links/notifications and secure credential storage. Offline state is visibly stale/degraded, never simulated as live.

Default hosted topology remains Postgres + outbox + stateless API/workers. Kafka, Kubernetes and service extraction require measured throughput, isolation, replay or ownership evidence—not projected user counts alone.

## 9. Maturity and commercial boundary

The eleven-domain 5/5 target is evidence-backed. A domain closes only when its machine criteria are satisfied and no Critical/High gap contradicts the claim. Adaptive maturity is earned over time; documentation or version naming cannot award it.

Development Production Ready means technically robust, secure, persistent, cross-platform, tested and full provider/data capable. Commercial/Public Ready adds provider licensing/rights, public-user legal/compliance and a commercial activation audit. Development closure never grants public data rights or activates commerce.

## 10. Planning horizons

- **Immediate / active:** finish the machine-current v19.0 dependency band; never use future audit work to bypass it.
- **Next 30–90 days:** close v19.0 technical trust, then golden-characterize and establish v19.1 canonical contracts/provider boundary.
- **6 months:** Postgres v2/outbox/versioned serving plus shared Opportunity Lifecycle and server-owned deterministic policy.
- **12 months:** first-class Watchlist, frozen Research Brief, Discovery/Radar convergence, event/options/temporal foundations and cross-platform beta distribution.
- **24 months:** governed outcome learning, pattern/similarity challengers, adaptive synthesis and operational/commercial readiness based on evidence.
- **Long term:** one explainable, lineage-preserving market-intelligence system whose shared evidence improves every product surface and client.

The exact active dependency and next action always come from `governance/current-state.json`, the active closure ledger and `handoff/CURRENT.md`.

## 11. Zero-miss and roadmap change rule

No version closes with an unexplained applicable certified responsibility, HOST/backlog/legacy commitment, audit finding/risk, surface disposition, ADR, role/right/platform case, compatibility migration or regression owner.

Roadmap changes require source-overlap review, explicit material Decision Log entry when scope or sequencing materially changes, and synchronized machine maps. A `CURRENT_*` projection or handoff cannot change this roadmap.
