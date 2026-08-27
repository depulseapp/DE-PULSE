# CURRENT Adaptive Build Process

**Canonical rebaseline:** `governance/V19_V20_REBASELINE.md`  
**Product-audit rebaseline:** `governance/PRODUCT_AUDIT_REBASELINE_2026_08_27.md`  
**Full audit coverage:** `governance/PRODUCT_AUDIT_COVERAGE_2026_08_27.md`  
**5/5 maturity target:** `governance/programs/V19-V20-REBASELINE/product-audit-5x5-target.json`  
**Audit finding register:** `governance/programs/V19-V20-REBASELINE/product-audit-finding-register.json`  
**Provider-registry additive rebaseline:** `governance/V19_V20_PROVIDER_REGISTRY_REBASELINE.md`  
**Provider-registry permanent contract:** `adaptive-governance/ADAPTIVE_PROVIDER_REGISTRY_CONTRACT.md`  
**Machine current-state authority:** `governance/current-state.json`  
**Provider-registry machine map:** `governance/programs/V19-V20-REBASELINE/adaptive-provider-registry.json`  
**Certified Stable:** `v18.10.0` — immutable  
**Active version:** `v19.0.0`  
**Active work slice:** `ADAPT-HOSTED-TRUST-FOUNDATION-001`  
**Active issue/PR/branch:** #148 / PR #149 / `adapt-hosted-trust-foundation-001`  
**Canonical closure ledger:** `governance/work-slices/ADAPT-HOSTED-TRUST-FOUNDATION-001/closure.json`

The execution loop remains source-driven and exact-head, now extended by audit conservation:

`LOOKUP -> COMPARE -> CLASSIFY -> AUDIT-DISPOSITION -> DECIDE -> CHARACTERIZE -> UPDATE/SHADOW -> RETEST -> Fast -> Qualified -> RECONCILE -> G11-G16 when a release is produced`

## Product-audit process overlay — mandatory

Every changed capability must be reconciled against:
- all 180 certified v18 responsibilities;
- HOST-001..072;
- mapped backlog and legacy commitments;
- all ten `AUDIT-EXEC-*` findings;
- all `AUDIT-RISK-*` rows;
- applicable surface disposition;
- applicable 5/5 maturity closure rows;
- applicable ADR decision(s);
- current source/executable evidence.

No audit finding is considered addressed merely because a version title sounds related. No finding is considered fixed by documentation alone.

### Architecture-migration process

Authority moves must use:

`CHARACTERIZE -> ADD NEW OWNER -> DUAL WRITE/DUAL READ OR SHADOW -> COMPARE -> MIGRATE CONSUMERS -> PROVE EQUIVALENCE/IMPROVEMENT -> REMOVE OLD AUTHORITY -> KEEP REGRESSION`

Mandatory characterization targets before extraction include current `computePlan`, technical features, Day/Swing/Long horizon behavior, Rapid Move, Opportunity Radar, Discovery handoff and Research outputs. Golden vectors must freeze current behavior and known defects separately so we do not accidentally certify a defect as the target contract.

A migration may intentionally improve a known defect only when the changed expected behavior is explicitly versioned and adversely tested.

## Version-first operating rule

The planning and closure unit is a coherent version/build. Do not use requirement packets or audit rows as separate roadmap/build abstractions.

Inside a version:
- implement dependency-correct changes in coherent commits;
- use focused local/unit/static evidence while editing;
- batch a coherent candidate for exact-head Fast;
- run Qualified at material risk boundaries and at G10 according to Impact Planner/current governance;
- do not create requirement-sized branches, PRs, public versions or workflows;
- do not weaken evidence simply to reduce CI use;
- do not pull future audit architecture into the current active dependency band when it would bypass unresolved technical evidence.

## Mandatory G2/G3/G10 audit dimensions

For each changed responsibility:
1. prove or assign the canonical owner;
2. prove current source overlap and consolidate duplicate owners;
3. map upstream evidence -> deterministic features/state -> downstream consumers;
4. identify professional trader/investor decision utility;
5. apply #170 cross-integration and Market Regime disposition;
6. apply #171 plus the audit surface disposition when visible;
7. prove stale/missing/partial/contradictory evidence truth and recovery re-evaluation;
8. preserve point-in-time/no-lookahead truth where history/outcomes/adaptation are involved;
9. prove role/RBAC/product-entitlement/provider-right separation and negative authorization where applicable;
10. prove persistence/restart/migration, load/backpressure and required platforms where applicable;
11. bind durable regression ownership;
12. record applicable `AUDIT-EXEC-*`, `AUDIT-RISK-*`, ADR and 5/5-domain rows;
13. define compatibility/dual/shadow/rollback behavior whenever authority or schema moves;
14. define temporal truth: source/observed/ingested/effective/as-of/expiry/session/revision/provider/dataset/rights/quality where material;
15. define raw-vs-adjusted basis, instrument identity and revision semantics for price/history-sensitive work;
16. define censoring and unbiased/control outcome sampling for learning/evaluation work;
17. define personalization versus shared model/policy truth;
18. define correlated-provider/upstream independence when corroboration or failover claims depend on multiple providers.

## Canonical intelligence process

Target processing order:

`observations/events -> rights/quality -> deterministic features/evidence -> SymbolIntelligenceSnapshot -> Opportunity Lifecycle -> projections/Decision Brief -> outcomes -> governed adaptation`

Hard authority rules:
- deterministic server-owned calculations are the canonical market/technical truth;
- Watchlist changes selected universe/user intent, not scoring authority;
- Discovery is broad-universe projection; Radar/Rapid Move are evidence/detection inputs, not parallel product state machines;
- Research consumes the same frozen snapshot/transition that produced promotion and remains the authoritative Decision Brief;
- alerts consume material transitions/incidents and dedupe causal chains rather than independently rescoring;
- AI/LLMs explain rights-filtered structured evidence with evidence IDs; they do not calculate canonical prices/indicators/options exposure/weights/routing/rights or lifecycle transitions.

### Undefined product semantics

`Long King / Short King` are not current formal domain concepts. Before implementation, define:
- user purpose and horizon;
- required evidence families and exclusions;
- relationship to explicit `LONG / SHORT / NO_SETUP-WAIT` side;
- quality/confidence/freshness/contradiction semantics;
- outcome definition and evaluation;
- Watchlist/Discovery/Research/Desk cross-integration;
- whether the concept adds information beyond existing setup/lifecycle fields.

`Call Wall / Put Wall` require expiry-aware options quality, OI as-of date, strike/cluster semantics, coverage and rights. Never infer signed dealer positioning from gamma×OI alone.

## Temporal / point-in-time process

Every material fact/event must have an explicit disposition for:
- `source_at`;
- `observed_at`;
- `ingested_at`;
- `effective_from` / `effective_to`;
- `as_of` decision cutoff;
- expiration or half-life;
- market session + exchange timezone;
- revision / supersedes identity;
- provider/dataset;
- rights version;
- quality/coverage.

Rules:
- retrieval/cache time is bookkeeping, not evidence time;
- unknown source/observation time remains unknown;
- fundamentals and macro facts preserve vintages;
- corrections create new revisions;
- raw/adjusted series are explicit;
- OI time is separate from live option quote/IV time;
- historical evaluation uses strict point-in-time joins;
- halts/missing bars/delistings/gaps can yield censored outcomes rather than false PASS/FAIL.

## Adaptive operating process

Three authorities stay separate:
1. deterministic rules;
2. registered statistical/adaptive models;
3. LLM synthesis.

Controlled adaptive sequence:
`freeze feature/evidence/regime/rights/policy -> decision/transition -> explicit outcome -> point-in-time evaluation -> time-split challenger -> shadow compare -> drift/subgroup checks -> sample/metric gates -> human approval -> gradual production -> rollback threshold -> reproducibility`.

Mandatory evaluation considerations:
- calibration, ranking lift, precision at limited attention capacity, false-promotion cost and stability rather than generic accuracy;
- time-split/walk-forward evaluation;
- deterministic/no-change baselines;
- uncertainty and sample counts;
- regime/liquidity/catalyst/evidence-coverage subgroup stability;
- adaptive selection bias caused by hydrating/promoting only existing-policy candidates;
- explicit censored outcomes;
- model/prompt/schema/policy/evidence versioning;
- explanation grounding evaluated separately from decision quality.

Pattern work starts with structured split-safe multi-timeframe feature vectors and interpretable similarity/baseline/tree/linear challengers. Deep chart-image/sequence models are later only if out-of-sample evidence justifies them.

## Conserved Data Health process

Provider/data changes use the canonical #80 baseline and dependency chain #81/#82/#83/#78/#84.

- Smart Provider Router v2 remains the executable authority for general routable provider capabilities; explicit direct-authority evidence such as SEC/EDGAR is preserved.
- Every new/changed provider/capability/fetch path is reconciled against `provider-capability-matrix.json`, `data-health-slo.json` and `provider-fetch-paths.json`; unclassified production network/provider behavior fails closed.
- canonical freshness is based on provider observation/event/publication/filing time; retrieval/cache timestamps are bookkeeping.
- before declaring degradation, reuse policy-valid canonical warm/cache evidence and eligible fallback.
- scope health at the smallest truthful capability/symbol/consumer level before escalation.
- recovery is automatic with hysteresis/anti-flapping/authority rules.
- capability lifecycle remains `SHADOW -> VALIDATED -> APPROVED -> PRODUCTION`; reachability/key presence/transient success never auto-promotes authority.
- under pressure, protect critical decision evidence and shed optional/background work first.
- provider agreement carries upstream-independence/correlation truth; two vendors can share one upstream failure.
- corrected provider facts/events create revision/supersedes lineage.

## Adaptive Provider Registry execution process

Every new provider uses `adaptive-governance/ADAPTIVE_PROVIDER_REGISTRY_CONTRACT.md` and issue #153 when part of v19.1 provider adoption.

Required process:

`source overlap -> adapter/manifest -> secret metadata -> Settings integration -> deterministic fixtures -> capability/entitlement probe -> provider matrix/fetch-path/rights classification -> Registry self-registration -> Router-only proof -> Data Health integration -> cross-integration -> SHADOW evidence -> adverse/recovery tests -> governed promotion decision`

### G2 — architecture / canonical ownership
Prove:
- provider-specific code ends at one adapter boundary;
- Adaptive Provider Registry is registration/capability projection only, never Router v3;
- Smart Provider Router v2 remains sole general selection authority;
- no provider-specific freshness/cache/persistence/subscription/telemetry/canonical-state owner is introduced;
- direct authorities remain explicit;
- Settings uses canonical secret/config owners;
- consumers request capabilities rather than provider names;
- cross-integration consumers are identified before implementation.

### G3 — provider/readiness design
Define:
- adapter manifest/schema version;
- auth/secret-reference and environment fallback;
- capability, freshness, history, entitlement, quota/rate-limit and authority semantics;
- deterministic unsupported/unentitled/withheld states;
- source/provider observation-time semantics;
- capability probes and reprobe triggers;
- plan/subscription upgrade/downgrade behavior;
- runtime eligibility vs governed lifecycle distinction;
- rights/public-production behavior;
- cross-integration matrix;
- failure, recovery, load and platform test plan;
- provider correlation/upstream dependency and revision behavior.

### G4/G5/G6 — implementation and integration
Prove:
- adapter self-registration occurs in actual runtime;
- configured secret is redacted and resolved canonically;
- shared provider-test path works without leakage;
- Registry candidate projection feeds existing Router;
- consumers reach canonical routed state rather than provider-specific fetches;
- duplicate fetch/subscription paths are absent;
- provider transport quirks are normalized at adapter boundary.

### G7 — data/security/adaptive intelligence
Prove:
- delayed/live/history/freshness truth;
- entitlement and rights are distinct;
- no automatic lifecycle promotion;
- direct-authority protection;
- source/provenance/observation-time/revision truth;
- optional-provider failures remain scoped;
- provider quality/utility evidence is bounded/auditable;
- AI/LLM/agents cannot fetch provider truth independently or bypass Router/rights/authority.

### G8 — load/capacity/provider economics
Prove applicable:
- quota/credits/rate-limit headroom;
- batching/cache reuse/calls avoided;
- backpressure and bounded fallback depth;
- no duplicate provider fan-out;
- outage/recovery without flapping;
- protected decision-support workloads survive pressure;
- cost/capacity budgets at expected user scale.

### G9/G10 — UX, cross-integration and zero-miss reconciliation
Prove:
- Settings, Maintenance and Data Health agree on configured/eligible/serving/fallback/freshness states;
- required consumers use canonical routed state;
- provider-specific UI/page fetches do not exist without governed authority exception;
- required Mac/Windows/Web Settings/admin surfaces consume one generic metadata contract;
- no provider/capability is complete solely because a Test button succeeds;
- applicable audit and 5/5 rows are reconciled.

## Market Data execution specifics

Market Data (`marketdata.app`) remains the first provider that must prove the generic process end-to-end, while fitting the v19.1 canonical-intelligence boundary rather than defining its own architecture.

Current official vendor semantics audited 2026-08-26 are changeable provider-contract evidence:
- Bearer authentication;
- `MARKETDATA_TOKEN` environment fallback;
- HTTP 200 and 203 success handling;
- trial data currently delayed;
- current trial/plan limits treated as observations, not eternal constants.

Mandatory regressions include missing/invalid token, 401/403, 429/quota exhaustion, 5xx/timeout, HTTP 203, malformed/partial/schema drift, delayed-vs-live truth, subscription downgrade/upgrade, capability loss/expansion, Router fallback/recovery, no lifecycle auto-promotion, no duplicate active fetch/subscription, dependent re-evaluation, UI/health semantic agreement and zero secret leakage.

## UI / information architecture process

Every visible element gets `KEEP / IMPROVE / CONSOLIDATE / COLLAPSE / MOVE / REMOVE / ADD_FIRST_CLASS`. Preserve useful deep evidence even if presentation moves.

Audit surface intent is authoritative planning input:
- Dashboard stays concise;
- Market Intelligence remains market-context owner;
- Day/Swing/Long keep workflows while server authority moves;
- Discovery and Watchlist are different universe projections of one lifecycle;
- Radar remains detector/evidence producer;
- Research becomes frozen Decision Brief owner;
- AI may consolidate common workflows into Research/global command but keeps bounded synthesis capability;
- Admin/Maintenance/Settings remain role-gated;
- News/Earnings/Filings symbol presentation consolidates into Research where useful.

## Hosted/client/platform process

Core owns authoritative data acquisition, secrets/rights, deterministic intelligence, opportunity/briefs, outcomes/adaptive jobs, tenant persistence/outbox/recovery/audit and AI gateway policy. Clients own rendering/interactions/charts/accessibility/typed schema validation. Desktop edge may hold encrypted authorized last-known cache, OS notifications/deep links and device credentials in OS secure storage.

No old desktop client may silently consume an incompatible event/API schema. Define versioning/deprecation/minimum-client/forced-upgrade rules before Web/client migration.

## Exactly one next action

Finish the audit-governance rebaseline checkpoint and exact-head Fast. Then continue `v19.0.0` on PR #149 at the machine-current earliest open dependency: **HOST-010..012 real managed backup/PITR/operator recovery evidence through HOST-016**. Do not begin HOST-013+, v19.1 canonical extraction, Market Data, Watchlist or Opportunity Lifecycle implementation while that dependency remains unresolved unless governance explicitly reclassifies it.
