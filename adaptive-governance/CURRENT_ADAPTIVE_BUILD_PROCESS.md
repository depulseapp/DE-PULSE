# CURRENT Adaptive Build Process

**Canonical rebaseline:** `governance/V19_V20_REBASELINE.md`  
**Provider-registry additive rebaseline:** `governance/V19_V20_PROVIDER_REGISTRY_REBASELINE.md`  
**Provider-registry permanent contract:** `adaptive-governance/ADAPTIVE_PROVIDER_REGISTRY_CONTRACT.md`  
**Machine current-state authority:** `governance/current-state.json`  
**Provider-registry machine map:** `governance/programs/V19-V20-REBASELINE/adaptive-provider-registry.json`  
**Certified Stable:** `v18.10.0` — immutable  
**Active version:** `v19.0.0`  
**Active work slice:** `ADAPT-HOSTED-TRUST-FOUNDATION-001`  
**Active issue/PR/branch:** #148 / PR #149 / `adapt-hosted-trust-foundation-001`  
**Canonical closure ledger:** `governance/work-slices/ADAPT-HOSTED-TRUST-FOUNDATION-001/closure.json`

The execution loop remains source-driven and exact-head:

`LOOKUP -> COMPARE -> CLASSIFY -> DECIDE -> UPDATE -> Fast -> Qualified -> G11-G16 when a release is produced`

## Version-first operating rule

The planning and closure unit is a coherent version/build. Do not use requirement packets as the current roadmap/build abstraction.

Inside a version:
- implement dependency-correct changes in coherent commits;
- use focused local/unit/static evidence while editing;
- batch a coherent candidate for exact-head Fast;
- run Qualified at material risk boundaries and at G10 according to Impact Planner/current governance;
- do not create requirement-sized branches, PRs, public versions or workflows;
- do not weaken evidence simply to reduce CI use.

A feature-heavy version may be split into an actual patch version when source/risk evidence proves that is safer. Current planned heavy splits are `v19.4.1`, `v19.5.1` and `v20.3.1`.

## Mandatory G2/G3/G10 audit dimensions

For each changed responsibility:
1. prove or assign the canonical owner;
2. prove current source overlap and remove/consolidate duplicate owners;
3. map upstream evidence -> derived state -> downstream consumers;
4. identify professional trader/investor decision utility;
5. apply #170 cross-integration and Market Regime disposition;
6. apply #171 UI/data-density/intelligence-maturity disposition when visible;
7. prove stale/missing/partial/contradictory evidence truth and recovery re-evaluation;
8. preserve point-in-time/no-lookahead truth where outcomes/history/adaptation are involved;
9. prove role/RBAC/product-entitlement/provider-right separation and negative authorization where applicable;
10. prove persistence/restart/migration, load/backpressure and required platforms where applicable;
11. bind durable regression ownership.

## Conserved Data Health process

Provider/data changes use the canonical #80 baseline and the dependency chain #81/#82/#83/#78/#84. This chain remains active as a process invariant even though its original issues are complete.

- Smart Provider Router v2 remains the executable authority for general routable provider capabilities; explicit direct-authority evidence such as SEC/EDGAR is preserved rather than forced through a rank-swappable route.
- Every new or changed provider/capability/fetch path is reconciled against `provider-capability-matrix.json`, `data-health-slo.json` and `provider-fetch-paths.json`; unclassified production network/provider behavior must fail closed.
- `canonical freshness` is based on provider observation/event/publication/filing time. Retrieval/cache timestamps are bookkeeping, and unknown observation time stays unknown.
- Before declaring degradation, reuse policy-valid canonical warm/cache evidence and eligible fallback; do not invent freshness or hide a genuinely missing required input.
- Scope health at the smallest truthful capability/symbol/consumer level before escalation. Optional-provider failures do not create false app-global degradation.
- Recovery is automatic when required canonical evidence becomes healthy, with hysteresis, anti-flapping and authority rules preserved.
- Capability lifecycle remains `SHADOW -> VALIDATED -> APPROVED -> PRODUCTION`; reachability, API-key presence or transient success never auto-promotes authority.
- Under provider/runtime pressure, protect critical decision evidence first and shed optional/background work before core Data Health deteriorates.

## Adaptive Provider Registry execution process

Every new provider must be implemented through `adaptive-governance/ADAPTIVE_PROVIDER_REGISTRY_CONTRACT.md` and issue #153 when part of v19.1 provider adoption.

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
- failure, recovery, load and platform test plan.

### G4/G5/G6 — implementation and integration
Prove:
- adapter self-registration occurs in actual runtime;
- configured secret is redacted and resolved canonically;
- shared provider-test path works without secret leakage;
- Registry candidate projection feeds the existing Router;
- consumers reach canonical routed state rather than provider-specific fetches;
- duplicate fetch/subscription paths are absent;
- provider transport quirks are normalized at the adapter boundary.

### G7 — data/security/adaptive intelligence
Prove:
- delayed/live/history/freshness truth;
- entitlement and rights are distinct;
- no automatic lifecycle promotion;
- direct-authority protection;
- source/provenance/observation-time truth;
- optional-provider failures remain scoped;
- provider quality/utility evidence is bounded/auditable;
- AI/LLM/agents cannot fetch provider truth independently or bypass Router/rights/authority.

### G8 — load/capacity/provider economics
Prove applicable:
- quota/credits/rate-limit headroom behavior;
- batching/cache reuse/calls avoided;
- backpressure and bounded fallback depth;
- no duplicate provider fan-out;
- provider outage and recovery without flapping;
- protected decision-support workloads survive provider/runtime pressure.

### G9/G10 — UX, cross-integration and zero-miss reconciliation
Prove:
- Settings, Maintenance and Data Health agree on configured/eligible/serving/fallback/freshness states;
- all required cross-integration consumers use canonical routed state;
- provider-specific UI/page fetches do not exist without governed authority exception;
- required Mac/Windows/Web provider Settings/admin surfaces consume one generic metadata contract;
- no capability or adapter is considered complete solely because a Test button succeeds.

## Market Data execution specifics

Market Data (`marketdata.app`) is the first provider that must prove this process end-to-end.

Current official vendor semantics audited 2026-08-26 must be represented as changeable provider-contract evidence, not eternal assumptions:
- Bearer authentication;
- `MARKETDATA_TOKEN` environment fallback;
- HTTP 200 and 203 successful response handling;
- trial data currently delayed;
- current Trader Trial 100,000 credits/day and one-year non-AAPL historical limit;
- current paid Trader real-time stock/options capability subject to endpoint freshness.

Mandatory regression scenarios:
- missing token;
- invalid token / 401;
- 403 including provider operational restrictions;
- 429/quota exhaustion;
- 5xx/timeout;
- HTTP 203 success;
- malformed/partial/schema drift;
- stale/delayed data must not be labeled live;
- trial -> Free downgrade;
- trial/free -> paid/live entitlement expansion;
- capability disappears after downgrade;
- Router fallback and recovery;
- no automatic lifecycle promotion;
- no duplicate active subscriptions/fetches;
- cross-integration dependent re-evaluation;
- Settings/Maintenance/Data Health semantic agreement;
- secret never appears in renderer, logs, URLs, telemetry, fixtures or GitHub evidence.

Effective capability, not provider marketing plan name, is the technical routing contract. A provider account upgrade may automatically expand **technical eligibility** after reprobe, but lifecycle/authority promotion remains explicit/governed.

## Intelligence process

The target maturity is:

`canonical evidence -> deterministic intelligence -> cross-feature synthesis -> point-in-time outcome accumulation -> bounded adaptive learning -> optional AI/agent explanation/orchestration`

Rules:
- deterministic logic remains authoritative for market truth unless a governed adaptive promotion explicitly changes a bounded influence;
- adaptive learning starts with point-in-time outcomes and SHADOW evidence;
- one specialized feature may contribute evidence to Market Regime but cannot create a competing regime engine;
- one symbol-level observation cannot directly flip global regime;
- missing evidence is not neutral evidence;
- recovery/new canonical evidence must re-evaluate dependent consumers;
- AI/agents consume canonical rights-filtered capabilities; they do not independently fetch market truth, invent evidence, rewrite canonical rules or require one model vendor.

Provider intelligence follows the same rule. Provider availability, freshness, latency, quota pressure, fallback success, schema integrity, disagreement, calls avoided, cost and usefulness may adapt Router eligibility/selection within governed policy and produce promotion recommendations. They do not silently promote lifecycle authority or infer rights.

## UI / information architecture process

For every visible element choose `KEEP / IMPROVE / CONSOLIDATE / COLLAPSE / MOVE / REMOVE`. Preserve access to useful deep evidence even if it leaves the primary page.

Hard constraints:
- preserve Day/Swing/Long Desk look-and-feel and workflow;
- preserve Dashboard Market Regime and Desk Control materially;
- preserve Data Engine look-and-feel except proven defects;
- preserve current AI Copilot engine/header visual treatment unless separately justified.

Provider Settings should be metadata-driven and concise. Detailed provider rankings, cooldowns, headroom, latency and recovery evidence belong in Maintenance/Data Health rather than ordinary user surfaces.

## Historical and current authority

Frozen v18 T1-T10 remains baseline conservation evidence; do not rerun history mechanically when unchanged. Current changed code and new requirements receive fresh impact-selected evidence. Smart Provider Router v2, Data Health/freshness, cache/persistence, subscription, telemetry/reconciliation/lifecycle, identity/session and direct-authority boundaries remain canonical.

## Exactly one next action

Continue `v19.0.0` on PR #149 from current live GitHub/executable evidence. The Adaptive Provider Registry / Market Data rebaseline is approved future v19.1 scope and must not bypass unfinished v19.0 Development Production Ready work. Governance-only rebaseline commits require fresh exact-head Fast before dependency-band advancement or release qualification.
