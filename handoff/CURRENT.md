# DE.PULSE — Current Authoritative Handoff

**SUPERSEDES ALL PRIOR CHAT HANDOFFS. GitHub is authoritative; chat history and temporary workspaces are advisory only.**

## Current identity

**Certified Stable:** `v18.6.1-stable`  
**Stable promotion commit:** `42e8432f7530ae39cbfd6ceb0b0bd5f6311dc5cc`  
**Qualified v18.6.1 source head:** `5c3fae486f3e4b4a39a0b1d549916aea9e1295fd`  
**Certified Stable source fingerprint:** `b01c14e1d54b736785eab6c03407801c527edd7769ff6f3d41fd4b20dabebd75`  
**Canonical v18.6.1 Release run:** `32279232665`  
**Current merged engineering baseline:** `f1a9e0d0d76d4be565ac9355a09f77a33a3338a7` (`v18.6.7` / PR #51)  
**Engineering branch:** `v18.7.0-development`  
**Current slice:** `v18.7.0 — Runtime Reliability & Data Truth`  
**Scope:** G1 FROZEN  
**Candidate package identity:** `18.7.0` / `v18.7.0-stable-20260819`  
**Current PR:** resolve live GitHub state for `v18.7.0-development`; use exactly one Draft PR and never create a second PR to retrigger CI.  
**Repository:** `depulseapp/DE-PULSE`  
**Last updated:** 2026-08-19 America/Vancouver

**Important:** candidate source identity saying `STABLE` is package/release identity, not proof of promotion. The current certified Stable remains `v18.6.1-stable` until the v18.7 exact-head PR evidence passes, the exact source merges, and one G11–G16 Release run certifies/packages/audits/publishes it.

## Resume rule

Before doing any new work, read:

1. `AGENTS.md` or `CLAUDE.md`;
2. `governance/CI-EFFICIENCY-CONTRACT.md`;
3. this handoff;
4. `release_identity.json`;
5. `.depulse-certification/resume/build-checkpoint.json` and `.depulse-certification/resume/release-evidence-checkpoint.json`;
6. live GitHub branch/PR/status state;
7. `adaptive-governance/CURRENT_ADAPTIVE_ROADMAP.md`;
8. `adaptive-governance/CURRENT_ADAPTIVE_BUILD_PLAN.md`;
9. `adaptive-governance/CURRENT_ADAPTIVE_BUILD_PROCESS.md`;
10. `adaptive-governance/CURRENT_ADAPTIVE_DELIVERY_PROCESS.md`;
11. `v18_7_0_scope.json`;
12. `v18_7_0_g0_g3_contract.json`;
13. `adaptive-governance/V18_7_0_RUNTIME_RELIABILITY_AUDIT.md`;
14. `release/v18.7.0/release_contract.json`;
15. `release/v18.7.0/run_full_certification.sh`;
16. `release/v18.5.1/V17-V18-IMPLEMENTATION-RECONCILIATION.json`.

Never resume from model memory alone.

### Stable checkpoint rule

The two `.depulse-certification/resume/` checkpoints intentionally remain anchored to the immutable published v18.6.1 Stable until a new release is actually certified. Do **not** rewrite Stable PASS evidence to pretend the mutable v18.7 development head has already passed G11–G16. The current handoff + live branch/PR state represent in-flight engineering; the Stable checkpoints remain the last promoted recovery anchor.

## Immutable v18.6.1 Stable evidence

- Final exact-head Fast: run `32279055139` · PASS.
- Final exact-head Qualified: run `32279113032` · PASS.
- Product full Qualified: run `32276304863` · PASS including backend/race/randomized, renderer/deterministic and browser evidence.
- Canonical Release run `32279232665` · G11–G16 PASS.
- macOS Apple Silicon and Windows x64 packaged runtime evidence PASS.
- Stable tag `v18.6.1-stable` points to `42e8432f7530ae39cbfd6ceb0b0bd5f6311dc5cc`.

Permanent boundaries remain: U.S. Equities Processing, No Execution, deterministic Day/Swing/Long market truth, Smart Provider Router v2 sole routing ownership, provider/Market-Mode governance, GLD/SLV/USO actionable tradable exceptions and Adaptive Intelligence governance.

## v18.6.7 baseline closure

v18.6.7 / PR #51 is complete and merged to `main` at `f1a9e0d0d76d4be565ac9355a09f77a33a3338a7`.

Its exact PR source head `3f5b9fe3084adf2dd91aa5a1335ec94a3c042c0f` passed:

- `DE.PULSE/fast-head` · run `32307553267`;
- `DE.PULSE/qualified-head` · run `32307659764`.

It conserved the 296-row historical v17/v18 authority ledger, established current legacy test/gate inventory/ownership and completed the first safe capability-oriented test organization wave. It was engineering/process/test organization only and did not create a new Stable.

## v18.7 source-owner audit conclusion

v18.7 is a **hardening slice of existing owners, not an architecture expansion**.

Do not create parallel retry, provider-health, freshness, backpressure, coalescing, disagreement or degradation engines.

Canonical owners remain:

- Smart Provider Router v2 — executable provider/capability routing/health;
- `data_freshness.go` — dataset/session/provider freshness;
- `runtime_slo.go` / `RuntimeSLOTracker` — SLO/recovery state;
- `workload_controller.go` — bounded tiers/queues/reserved critical capacity/load shedding;
- `runtime_observability.go` — request/load/budget telemetry;
- `broad_snapshot_broker.go` — broad-snapshot canonical reuse/coalescing;
- `provider_reconciliation.go` — independent-source agreement/conflict;
- `Engine.Snapshot()` — canonical runtime aggregation boundary.

Full disposition is in `adaptive-governance/V18_7_0_RUNTIME_RELIABILITY_AUDIT.md`.

## v18.7 implemented branch scope

### A. Truthful degradation and blast radius

`runtime_degradation.go` now exposes:

- compatible display `Code`;
- canonical machine-readable `ReasonCode`;
- `PressureState`;
- `DecisionImpact`;
- `CriticalUsable`;
- `Abstain`;
- affected datasets/capabilities;
- affected downstream consumers.

Current canonical reason taxonomy:

- `QUEUE_SATURATED`
- `LOCAL_OVERLOAD`
- `RATE_LIMITED`
- `NETWORK_FAILURE`
- `PROVIDER_DOWN`
- `LOW_COVERAGE`
- `UNKNOWN`

Pressure semantics:

- `HEALTHY` — no active decision-relevant degradation;
- `PROTECTED` — scoped degradation exists but critical decision evidence remains usable;
- `DEGRADED` — required decision evidence is insufficient/untrustworthy; affected conclusions must abstain;
- `RECOVERING` — evidence is healthy again but recovery confirmation has not completed.

If critical live evidence is unusable and no narrower cause is proven, runtime now fails closed to `DATA DEGRADED` / `UNKNOWN` / `Abstain=true`. Absence of diagnosis is never treated as health.

A provider queue at its hard `MaxQueue` now becomes `QUEUE_SATURATED` immediately rather than waiting for queue age to cross a delay threshold.

### B. Recovery hysteresis

The existing `RuntimeSLOTracker` owns recovery hysteresis:

- degradation appears immediately;
- recovery requires 3 consecutive healthy observations **and** >=5 seconds continuous healthy stability;
- any relapse resets the streak;
- until confirmed, the prior degradation remains exposed as `RECOVERING`;
- SLO diagnostics expose pending confirmation state.

This prevents one transient provider success from falsely flipping overall runtime health to healthy.

### C. `UNKNOWN` / `ABSTAIN`

Required decision evidence now fails closed. Optional isolated dataset degradation may remain `PROTECTED`; missing critical evidence becomes `DEGRADED` + `Abstain=true` even when the exact lower-level cause is still unknown.

### D. Existing reliability controls retained/requalified

The audit found the following already materially implemented and integrated, so they are **not rebuilt**:

- provider/capability entitlement/health/circuit aggregation;
- finite Router fallback/retry chains and provider cooldowns;
- dataset/session/provider freshness truth;
- bounded request budgets and priority shedding;
- BroadSnapshotBroker fresh reuse + single-flight/coalescing;
- independent-source reconciliation and explicit `CONFLICT` without silent averaging;
- canonical reuse/provider-calls-avoided telemetry.

### E. Active-market reliability proof — STAGED, NOT YET PASS

`v18_7_0_active_market_reliability_test.go` uses production owners to test:

- 16 concurrent equivalent broad-snapshot requests collapsing to one provider fetch;
- observable coalesced waiters;
- hard-bounded provider concurrency + queue;
- immediate queue-saturation degradation truth;
- capacity SLO blocking while still allowing scoped `PROTECTED` critical-data semantics;
- no admission beyond hard provider capacity/queue limits.

Do **not** call this PASS until exact-head CI executes it.

### F. Release reproducibility closure

The previously deferred Release workflow hardening is now included in v18.7:

- Fast, Qualified and Release third-party Actions are immutable-SHA pinned and dependency-lock governed;
- Release G12 uses the same pinned Playwright requirements file + safe pip cache contract as Qualified;
- `tools/ci/reproducibility_gate.py` enforces all three workflows;
- no fourth workflow, no trigger branch and no duplicate release lane were added.

Release Action/browser dependency pinning is now **closed**, not deferred.

## v18.7 release identity / renderer strategy

`release_identity.json`, `VERSION.txt` and `app_bootstrap.go` bind candidate identity `18.7.0` / `v18.7.0-stable-20260819`.

To avoid rewriting the large legacy `renderer.js` only for version constants, v18.7 generalizes the existing small-overlay concept:

- `renderer/release-identity-v18.7.0.js` is loaded last;
- it owns only release/build integrity display + v18.7 QA entry;
- the existing `watchlist-desk-contract-v18.6.1.js` remains loaded underneath because it still owns valid watchlist behavior;
- `release_identity.py` now understands declared release identity overlays;
- historical `certification_plan.json` and `ci_pipeline_plan.json` remain conserved registries at their declared legacy version rather than being cosmetically relabeled.

Current execution authority is the three canonical GitHub workflows plus `release/v18.7.0/release_contract.json` and `release/v18.7.0/run_full_certification.sh`.

## v18.7 qualification status

**Not yet CI-certified.** Branch implementation/pre-PR preparation is staged; do not infer PASS from source existence.

Required before merge:

1. one Draft PR from `v18.7.0-development` to `main`;
2. exact-head `DE.PULSE/fast-head = success`;
3. same PR marked Ready only after Fast success;
4. full Qualified including:
   - backend/full Go;
   - race detector;
   - randomized order;
   - renderer/deterministic contracts;
   - Chrome broad behavior;
   - WebKit primary compatibility;
5. conserved 296-row ledger / G10 integrity;
6. release/reproducibility contracts green;
7. merge only on the exact current qualified head.

Any legitimate defect is fixed on this same branch/PR. A source-changing fix creates a new exact head and requires fresh exact-head evidence. No retry/certification/promotion branches.

## Post-merge Stable path

After exact-head merge:

- G11 proves exact source-head Fast + Qualified status and source-head → merged-candidate fingerprint equivalence;
- G12 runs `release/v18.7.0/run_full_certification.sh` on the immutable candidate;
- G13/G14 package/audit macOS Apple Silicon + Windows x64 in parallel;
- G15 verifies evidence graphs and exact artifacts;
- Publish uploads the exact certified artifacts without rebuild;
- G16 records final evidence and handoff.

Only then may `v18.7.0-stable` replace v18.6.1 as the certified Stable.

## External/open governance item

Do not claim main branch protection/ruleset is fixed unless live GitHub confirms an authorized protection/ruleset change. The engineering process remains fail-closed through exact-head workflow evidence even where repository protection is not externally enforced.

## Exactly one next action

Resolve the live PR state for `v18.7.0-development` after final branch-diff/prevalidation. If no PR exists and the exact diff is coherent, open **exactly one Draft PR** to `main`. If a PR already exists, reuse it. Require exact-head Fast; do not mark Ready until Fast succeeds. Then require full Qualified including backend/race/randomized, renderer/deterministic, Chrome and WebKit, and merge only on exact-head green evidence. After merge, allow exactly one G11–G16 Stable Release run.
