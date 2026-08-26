# CURRENT Adaptive Build Process

**Canonical rebaseline:** `governance/V19_V20_REBASELINE.md`  
**Certified Stable:** `v18.10.0` — immutable  
**Active version:** `v19.0.0`  
**Active issue/PR/branch:** #148 / PR #149 / `adapt-hosted-trust-foundation-001`

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

## UI / information architecture process

For every visible element choose `KEEP / IMPROVE / CONSOLIDATE / COLLAPSE / MOVE / REMOVE`. Preserve access to useful deep evidence even if it leaves the primary page.

Hard constraints:
- preserve Day/Swing/Long Desk look-and-feel and workflow;
- preserve Dashboard Market Regime and Desk Control materially;
- preserve Data Engine look-and-feel except proven defects;
- preserve current AI Copilot engine/header visual treatment unless separately justified.

## Historical and current authority

Frozen v18 T1-T10 remains baseline conservation evidence; do not rerun history mechanically when unchanged. Current changed code and new requirements receive fresh impact-selected evidence. Smart Provider Router v2, Data Health/freshness, cache/persistence, subscription, telemetry/reconciliation/lifecycle, identity/session and direct-authority boundaries remain canonical.

## Exactly one next action

Continue `v19.0.0` on PR #149. Exact-head Fast #1165 on `f36417fda84e063d5a9cafcc31c464b051f5b3af` cleared recursive helper/source ownership with zero unregistered orphan production Go helpers, but the recursive source-health lane remained red because CURRENT Adaptive Data Health conservation had drifted. Restore that CURRENT governance contract and obtain fresh exact-head Fast. Do not start `v19.1.0` while v19.0.0 remains unqualified.
