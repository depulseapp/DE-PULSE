# CURRENT Adaptive Build Plan

**Canonical rebaseline:** `governance/V19_V20_REBASELINE.md`  
**Product-audit rebaseline:** `governance/PRODUCT_AUDIT_REBASELINE_2026_08_27.md`  
**Full audit coverage:** `governance/PRODUCT_AUDIT_COVERAGE_2026_08_27.md`  
**5/5 maturity target:** `governance/programs/V19-V20-REBASELINE/product-audit-5x5-target.json`  
**Audit finding register:** `governance/programs/V19-V20-REBASELINE/product-audit-finding-register.json`  
**Provider-registry additive rebaseline:** `governance/V19_V20_PROVIDER_REGISTRY_REBASELINE.md`  
**Provider-registry permanent contract:** `adaptive-governance/ADAPTIVE_PROVIDER_REGISTRY_CONTRACT.md`  
**Machine current-state authority:** `governance/current-state.json`  
**Backlog/version matrix:** `governance/programs/V19-V20-REBASELINE/backlog-version-matrix.json`  
**HOST/version map:** `governance/programs/V19-V20-REBASELINE/host-requirement-version-map.json`  
**Provider-registry machine map:** `governance/programs/V19-V20-REBASELINE/adaptive-provider-registry.json`  
**Legacy future-commitment conservation:** `governance/programs/V19-V20-REBASELINE/legacy-future-commitment-conservation.json`  
**Cross-integration matrix:** `governance/programs/V19-V20-REBASELINE/cross-integration-matrix.json`  
**Whole-product surface map:** `governance/programs/V19-V20-REBASELINE/whole-product-surface-rebaseline.json`  
**Certified Stable:** `v18.10.0` — immutable  
**Active build:** `v19.0.0` — Hosted Trust & Identity Foundation  
**Active work slice:** `ADAPT-HOSTED-TRUST-FOUNDATION-001`  
**Active issue/PR/branch:** #148 / PR #149 / `adapt-hosted-trust-foundation-001`  
**Canonical closure ledger:** `governance/work-slices/ADAPT-HOSTED-TRUST-FOUNDATION-001/closure.json`

## Product-audit 5/5 build overlay — mandatory

The 27-Aug-2026 full product audit is now a cumulative build authority. The plan must conserve and disposition every applicable row in the machine audit finding register, not only the ten Executive findings.

All ten Executive findings are mandatory build responsibilities:
1. replace `RuntimeSnapshot` as the de-facto domain boundary with versioned `SymbolIntelligenceSnapshot` + typed Evidence/events/deltas;
2. preserve Opportunity Radar as a detector/evidence producer but establish one shared Opportunity Lifecycle;
3. build first-class Watchlist as a user-selected-universe projection of shared lifecycle, never a second scorer/scanner;
4. build point-in-time outcomes and governed shadow/champion-challenger adaptation before any learned production influence;
5. characterize and move authoritative technical/desk/side/geometry/scoring logic from renderer to versioned Go domain owners;
6. finish hosted tenant persistence, managed secrets, service trust/IaC, real recovery/PITR, audit and operations;
7. close desktop distribution gaps: secure storage, signing/notarization/Authenticode, installer/update channels and rollback;
8. preserve executable provider/data-rights controls and require real rights evidence before applicable commercial/public activation;
9. remove/demote duplicate current-state narratives and generate status from machine/executable truth;
10. evolve as a Go modular monolith with Postgres/outbox/versioned APIs/thin clients; no speculative microservice/Kafka/Kubernetes program.

The full audit also adds mandatory cross-cutting rows for instrument identity, point-in-time/vintage data, calendar/DST/clock skew, raw/adjusted basis, options OI as-of quality, alert causal dedupe, provider correlation/revisions, privacy-vs-audit retention, cross-device conflicts, offline truth, AI egress/prompt injection, old-client schema compatibility, personalization/calibration separation, adaptive selection bias, censored outcomes, full-snapshot fanout cost and operational ownership. `Long King / Short King` remain undefined and may not become product labels/Watchlist columns until an explicit evidence/horizon/outcome contract exists. `Call Wall / Put Wall` remain planned semantics until expiry-aware quality/rights/OI-as-of rules are approved.

### 5/5 claim rule

Every version that materially affects a maturity domain must identify which `product-audit-5x5-target.json` closure rows it advances. A domain may move to 5/5 only from current objective evidence and only when no unresolved Critical/High audit gap contradicts the claim. Architecture can become 5/5-capable before time-dependent maturity such as Adaptive Intelligence is earned.

## Build sizing rule

Plan and communicate work by **real version/build**, not requirement packets. Requirement rows, audit findings, issue acceptance bullets and CI evidence rows remain granular traceability units inside the version.

- Combine small, related changes when they share canonical owners and acceptance evidence.
- Give feature-heavy/risk-heavy work its own version or patch version. Current deliberate heavy splits include `v19.4.1` Discovery/Watchlist convergence, `v19.5.1` Options/GEX and `v20.3.1` AODR.
- Do not create a version for every requirement, audit row, card, page defect or CI checkpoint.
- Do not enlarge a version so much that ownership, adverse testing or rollback becomes unclear.
- Use commits and risk-based CI checkpoints for implementation progress; they are not product planning units.

## Required build matrix for every version

Before G3, every assigned requirement/issue/legacy-conservation/audit row must resolve:
- source-overlap disposition: `INHERITED / EXTEND_EXISTING_OWNER / REPLACE_CONSOLIDATE / NEW_RESIDUAL / EXTERNAL_BLOCKED / N_A`;
- canonical owner and upstream evidence;
- actual consumers and user/trader decision purpose;
- data freshness/provenance/point-in-time semantics;
- positive + negative/failure evidence;
- persistence/restart/migration applicability;
- role/auth/product-entitlement/provider-right applicability;
- load/resource/backpressure applicability;
- Mac/Windows/Web requirement or justified N/A;
- durable regression owner;
- applicable Executive-audit finding and audit-wide-risk IDs;
- applicable 5/5 maturity-domain closure rows;
- compatibility/shadow/dual-read/rollback strategy when authority moves.

For intelligence-bearing work also record:
- intelligence maturity: `DETERMINISTIC_ONLY / ADAPTIVE_CANDIDATE / LEARNING_ENABLED / AI_ASSISTED / NOT_USEFUL`;
- downstream cross-integration `REQUIRED / CONDITIONAL / NOT_USEFUL` for Market Regime, Tradeability, Discovery, Watchlist, Research, Desks, Prep, alerts, Outcome/Pattern, Data Health/Maintenance and processing priority;
- Market Regime contribution `YES / CONDITIONAL / NO`;
- Outcome Learning contribution `YES / CONDITIONAL / NO`;
- duplicate/isolation conflict and consolidation decision;
- re-evaluation behavior after canonical recovery/new evidence;
- temporal fields/as-of/revision/raw-adjusted/censoring applicability;
- personalization-vs-shared-policy boundary;
- selection-bias/unbiased-sampling requirement where outcomes/adaptation are involved.

For visible work also record #171 and audit surface disposition: `KEEP / IMPROVE / CONSOLIDATE / COLLAPSE / MOVE / REMOVE / ADD_FIRST_CLASS`.

## Compatibility-first authority migration

Every architecture extraction follows:

`CHARACTERIZE -> ADD NEW OWNER -> DUAL WRITE/DUAL READ OR SHADOW -> COMPARE -> MIGRATE CONSUMERS -> PROVE EQUIVALENCE/IMPROVEMENT -> REMOVE OLD AUTHORITY -> KEEP REGRESSION`

This is mandatory for renderer scoring, RuntimeSnapshot, Opportunity Radar/Discovery lifecycle, watchlists/workspaces, provider direct paths, full-snapshot SSE, identity/session state and persistence schema changes. A rewrite is not permission to delete working capability first.

## Canonical product/intelligence target

Target flow:

`Providers/events -> canonical observations/evidence -> deterministic features -> SymbolIntelligenceSnapshot -> Opportunity Lifecycle -> user/horizon projections -> frozen Decision Brief -> outcomes -> governed adaptation`

Surface rules:
- Dashboard = concise attention/operating summary;
- Market Intelligence = shared market/regime/liquidity context;
- Day/Swing/Long = horizon projections of shared evidence with server-owned two-sided deterministic policy;
- Discovery = broad-universe projection;
- Radar = detector/evidence producer;
- Watchlist = user-selected-universe projection;
- Alerts = material lifecycle transitions/incidents, not rescoring;
- Research = authoritative frozen-as-of Decision Brief;
- AI = evidence-bounded explanation/command surface, never parallel truth;
- Admin/Maintenance/Settings = role-gated operator surfaces;
- News/Earnings/Filings = preserve services, consolidate symbol presentation into Research unless a separate explorer proves necessary.

## Conserved Data Health build owners

Every version that changes provider access, freshness, cache/persistence, consumer health, routing, lifecycle, recovery or load behavior must reuse the completed #80/#81/#82/#83/#78/#84 program rather than inventing a local substitute.

The machine-readable build authorities are:
- `governance/data-health/provider-capability-matrix.json` — provider/capability owner, consumers, authority, freshness, fallback, materiality, degraded impact, router dataset and lifecycle;
- `governance/data-health/data-health-slo.json` — canonical evidence-time semantics, healthy coverage, freshness/fallback/degradation/recovery/load-protection metrics and truth rules;
- `governance/data-health/provider-fetch-paths.json` — every production external fetch path classified and owned as `MIGRATE`, `DIRECT_AUTHORITY` or justified `N/A`.

For applicable changes the Build Plan must reconcile these artifacts with the current Adaptive Roadmap, Build Process, Delivery Process and full audit finding register. Smart Provider Router v2 remains the general routable authority; direct-authority rows are not silently rank-swapped. Any new provider/capability/fetch path must be classified in the canonical artifacts before the build can claim zero-miss coverage.

Provider corroboration must also account for common upstream feeds. Two vendors agreeing is not automatically independent evidence. Corrections/revisions must create durable revision lineage rather than silently rewriting prior decision truth.

## Adaptive Provider Registry build plan

Every new provider must use `adaptive-governance/ADAPTIVE_PROVIDER_REGISTRY_CONTRACT.md`. Provider-specific logic should stop at one adapter; all later routing and cross-integration must be generic.

Mandatory provider work sequence:
1. source-overlap audit against Router/Data Health/secrets/cache/subscription/telemetry/canonical-state owners;
2. define one provider adapter + manifest/schema;
3. define provider auth/secret metadata without storing secrets in the Registry;
4. integrate the existing Data Providers Settings/secrets/test/clear/redaction path instead of creating provider-specific Settings services;
5. deterministic endpoint/schema/transport fixtures and vendor-quirk normalization;
6. bounded configured-key capability/entitlement probe where appropriate;
7. update provider-capability/fetch-path/rights/authority artifacts;
8. self-register adapter in the Adaptive Provider Registry;
9. prove Smart Provider Router v2 is the only general selection path;
10. define `REQUIRED / CONDITIONAL / NOT_USEFUL` cross-integration for every applicable consumer;
11. begin in SHADOW/validation-first unless existing evidence justifies another governed state;
12. collect health/freshness/latency/quota/cost/usefulness evidence;
13. prove outage, auth, rate-limit, stale, malformed, downgrade, recovery and subscription-change behavior;
14. prove no duplicate acquisition/subscription/canonical state;
15. apply Mac/Windows/Web and role/admin applicability;
16. lifecycle promotion remains an explicit governed decision after evidence;
17. classify provider-upstream correlation, revision behavior, rights dimensions and observation-time truth.

### Market Data first-adopter build requirements

`v19.1.0` / #153 remains the first generic Registry provider adoption, but it must fit inside the canonical-intelligence extraction rather than become a competing v19.1 architecture.

Minimum build responsibilities:
- adapter/manifest self-registration;
- canonical `MARKETDATA_TOKEN` environment fallback;
- Bearer-header authentication;
- generic provider Settings card using the TradeInsight-derived preserve/replace/clear/test/redaction secret behavior;
- token never returned to renderer or logged;
- HTTP `200` and `203` accepted as successful Market Data responses under current vendor semantics;
- delayed trial/free data represented truthfully and never as live;
- effective capability/entitlement probing so a later paid/live entitlement can change technical Router eligibility without a DE.PULSE source change solely because the account plan changed;
- no hard-coded `Trader => primary/live` routing shortcut;
- SHADOW-first evidence for capability quality/usefulness;
- quota/credit headroom, health, latency, freshness, rights and utility feed normal Router eligibility/scoring;
- trial/free downgrade and paid/live expansion regression scenarios;
- provider single-IP/account restriction failures classified truthfully when encountered;
- mandatory cross-integration through canonical state to applicable Research, Discovery/Radar, Desks, Prep, Market Intelligence/Regime contribution, alerts, history/options, Data Health/Maintenance and future Outcome Learning consumers.

A successful provider test button is not implementation closure.

## Rebaselined v19 build plan

- `v19.0.0` — **Hosted Trust & Identity Foundation**: HOST-001..023, #164 core auth/session, security portions of #156. Finish the current dependency band; current earliest OPEN band is HOST-010..012 and its real managed backup/PITR/operator recovery dependency. Product-audit architecture work must not be used to bypass unfinished technical trust evidence.
- `v19.1.0` — **Canonical Intelligence & Provider Foundation**: existing #150/#151/#153/#154/#155/#160/#167 core plus `AUDIT-EXEC-01`, `AUDIT-EXEC-05` foundation. Freeze golden vectors for current `computePlan`, technical state, Radar, Rapid Move and Research; define Observation/Evidence/Snapshot/Transition/DecisionBrief schemas; introduce `SymbolIntelligenceSnapshot` behind RuntimeSnapshot compatibility; establish server-owned technical/horizon package boundaries; continue Adaptive Provider Registry/Market Data through the shared capability architecture.
- `v19.2.0` — **Hosted Serving, Sync & Postgres v2**: HOST-024..039 plus audit persistence/API risks. Tenant-aware normalized Postgres v2, RLS/isolation disposition, revisions/outbox/conflict/tombstones, versioned user-scoped APIs/deltas, hosted cache/fanout rights isolation, cross-device/offline schema foundations.
- `v19.3.0` — **Shared Opportunity Lifecycle & Cross-Platform Product Contract**: HOST-040..047/053 + #152/#156/#159/#160 presentation/#167 Admin/#171/#164 UX + `LEGACY-TRADER-SETUP-SHORT-001` + `AUDIT-EXEC-02` + server-owned two-sided deterministic setup authority. Rapid Move/Radar/Discovery feed the shared Opportunity state machine in shadow/dual-read before old lifecycle authority retires. Establish shared auth/role/IA/client compatibility rules.
- `v19.4.0` — **Market Intelligence, Research Decision Brief & Watchlist Foundation**: HOST-049 + #158/#161/#162/#171 + `AUDIT-EXEC-03`. Research gains frozen-as-of Decision Brief identity; first-class Watchlist becomes the selected-universe projection of shared lifecycle with ranked attention/explanation/contradiction/freshness/trust and Research handoff; no Watchlist scorer.
- `v19.4.1` — **Discovery / Watchlist / Radar Convergence**: HOST-048 + #163/#171 + conserved halt/LULD/pause/resume behavior. Discovery remains broad-universe projection, Watchlist user-selected universe, Radar detector; durable transition alerts and lifecycle consistency proven across both projections.
- `v19.5.0` — **Price/Volume & Event-Anchored Intelligence**: #168/#169. New evidence feeds canonical snapshot/lifecycle/Decision Brief; event identity/revisions/causal alert correlation apply.
- `v19.5.1` — **Options Structure & GEX Intelligence**: #157. Formal Call/Put Wall semantics may enter only after expiry-aware coverage/OI-as-of/quality/rights rules. No unsupported signed dealer-positioning claim. Long King/Short King remains blocked until a formal deterministic domain definition is approved and mapped to horizon/side/outcome semantics.
- `v19.6.0` — **Point-in-Time Evidence & Outcome-Ready Foundation**: HOST-057..064 + deterministic #165 + institutional/two-sided thesis substrate + temporal/audit risks. Bitemporal/vintage facts, stable instrument identity, raw/adjusted lineage, explicit censored outcomes, unbiased/control outcome sampling and feature snapshot IDs.
- `v19.6.1` — **Hosted Reliability, Economics, Observability & 5/5 Readiness**: HOST-050..056/065..071 + ADR-GDI + provider-gap + final #170/#171 reconciliation + provider-registry reliability/coverage/economics/readiness scorecards + real operational SLO/on-call/recovery/scale evidence + desktop signing/update readiness + maturity-domain residual review. Any domain below 5/5 creates explicit residual work; no score inflation.
- `v19.7.0` — **v19 Major Deterministic/Hosted Closure**: HOST-072; no new feature scope. Zero unexplained audit/responsibility rows, compatibility migrations reconciled, exact-head Fast/Qualified/G0-G16. Commercial activation remains OFF unless Owner explicitly changes it.

### v19.3.0 two-sided deterministic Desk setup acceptance

Current source is not a true SHORT implementation: bearish labels are possible, but plan geometry/action-state logic remains long-oriented. `v19.3.0` must correct the **contract**, not merely change text or colors.

Required implementation behavior:
- explicit side: `LONG / SHORT / NO_SETUP-WAIT`;
- separate directional evidence from 0–100 setup-quality score;
- a strong SHORT setup can score highly for quality without being encoded as a low numeric score;
- LONG: target/trim above entry and invalidation below;
- SHORT: cover/target below entry and invalidation above;
- R-multiple/action-state/entry-distance/sort/chart overlays/replay snapshots work correctly for both sides;
- Research and Discovery consume the same canonical setup side/geometry rather than reconstructing it;
- existing Day/Swing/Long Desk look-and-feel remains materially unchanged; only truthful side-aware labels/content change where necessary;
- No Execution remains permanent.

The institutional/TDTI two-sided thesis substrate remains separate in v19.6.0/v20.3.0.

## Rebaselined v20 build plan

- `v20.0.0` Outcome Learning & Adaptive Control Plane: point-in-time evaluation contract, model/policy registry, sample floors, shadow/champion-challenger, drift and rollback; no autonomous promotion.
- `v20.1.0` Adaptive Chart Pattern & Similarity Intelligence (#166): begin with structured split-safe feature vectors and interpretable similarity/baselines before deep/image models.
- `v20.2.0` Adaptive Market Synthesis, Market Regime & Discovery Learning: conserved ASBI normalization/synthesis/contradiction/abstention/outcome behavior on top of v20.0 controls.
- `v20.3.0` Adaptive Institutional / Two-Sided Thesis.
- `v20.3.1` AODR Adaptive Opportunity Intelligence.
- `v20.4.0` Agent Orchestration & Controlled MCP/API.
- `v20.5.0` Adaptive Operations: bounded provider utility/cost/reliability priors inside Smart Provider Router v2 policy; no parallel router or automatic authority/rights promotion.
- `v20.6.0` Professional Adaptive Closure; no feature scope.

## Zero-miss audit rule

No version closes with an unassigned applicable certified-v18 responsibility, backlog row, HOST requirement, legacy future-roadmap/build-plan commitment, source-discovered responsibility, audit Executive finding, audit-wide risk, surface disposition, temporal/point-in-time rule, ADR, cross-integration, role/right case, UI disposition or regression owner.

For providers specifically, zero-miss includes the Registry manifest, Settings/secret path, capability/entitlement truth, Router eligibility, Data Health, rights/authority, upstream-correlation/revision truth, cross-integration, adverse/recovery evidence, platform applicability and lifecycle disposition.

For adaptive/outcome work, zero-miss includes time-split evaluation, point-in-time joins, censoring, selection-bias controls, cohort/sample counts, uncertainty, model/policy version, deterministic baseline/fallback, approval and rollback.

## Exactly one next action

Finish this governance rebaseline checkpoint with exact-head Fast. Then remain in `v19.0.0` / PR #149 at the earliest open dependency: **HOST-010..012 real managed PostgreSQL backup/PITR/operator recovery deletion-retention evidence through the existing HOST-016 recovery owner**. Do not start `v19.1.0`, Watchlist implementation, Opportunity Lifecycle extraction or HOST-013+ until the active dependency is truthfully resolved or explicitly reclassified by governance.
