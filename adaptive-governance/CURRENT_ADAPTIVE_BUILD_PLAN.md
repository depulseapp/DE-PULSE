# DE.PULSE — Current Adaptive Build Plan

**Operational overlay date:** 2026-08-19  
**Immutable Stable baseline:** `v18.6.1-stable`  
**Current engineering slice:** `v18.7.0-development` — Runtime Reliability & Data Truth  
**Scope status:** G1 FROZEN  
**Authority:** current execution overlay. Permanent contracts and historical release evidence remain intact.

## 1. Normal target lifecycle

`1 development branch → batch coherent work → 1 Draft PR → Fast → same PR Ready → full Qualified → merge → one Release G11–G16 run when release-capable → Stable`

Rules:

- Batch coherent source/governance/release preparation before PR whenever practical.
- Never create trigger/retry/certification/promotion branches.
- Never create a second PR merely to retrigger CI.
- Fix legitimate source/test/gate defects on the same branch/PR.
- Same-SHA infrastructure failure reruns only affected work when justified.
- Main push performs hygiene only.
- Publication uses exact same-run certified native artifacts; no post-certification rebuild.
- GitHub is the continuation authority across ChatGPT/Claude/accounts.

## 2. G0–G16 execution map

- **G0 Exact Stable Intake:** immutable Stable identity, current main, source/fingerprint, conserved ledger, CI state, open defects and dependencies.
- **G1 Immutable Scope:** exact committed scope and boundaries; no silent additions/removals.
- **G2 Architecture / Data Utility:** canonical owner, consumers, provider/rights, reuse, freshness, retention and duplication review.
- **G3 Design / Dependency / Impact Readiness:** Impact Planner risk classes, expected Fast/Qualified/browser/release lanes and CI cost shape.
- **G4 Development Exit:** one development branch, coherent source, frozen scope, release preparation consistent.
- **G5 Fast Qualification:** exact-head syntax/format/unit/contract/governance checks.
- **G6 Integration / Medium Qualification:** affected cross-module evidence.
- **G7 Data / Security / Adaptive Intelligence:** provider/data-rights/security/adaptive contracts.
- **G8 Performance / Capacity / Stability:** load, queues, coalescing, recovery, stability and latency evidence.
- **G9 Cross-Module / UI / UX:** renderer and browser behavior where affected.
- **G10 Pre-Freeze Qualified Candidate:** exact-head full Qualified; Chrome + WebKit primary proof; conserved-ledger/release-tooling evidence.
- **G11 Immutable Release Candidate:** merge bound to exact Fast + Qualified source head and equal source fingerprint.
- **G12 Full Certification:** current-source v18.7 certification from immutable candidate.
- **G13 Native Packaging / Provenance:** candidate-native artifacts.
- **G14 Actual Artifact Runtime Audit:** macOS Apple Silicon + Windows x64 packaged behavior/provenance.
- **G15 Release Assurance / Promotion:** evidence graph and exact artifact hashes verified.
- **Publish:** exact certified artifacts only; no rebuild.
- **G16 Adaptive Retrospective / Handoff:** durable release evidence, current source of truth and next intake.

No top-level gates beyond G0–G16.

## 3. Completed engineering baseline

- PRs #46–#50: CI/release process hardening, reproducibility, Chrome + WebKit, telemetry and Documentation ownership.
- **v18.6.7 / PR #51:** merged to `main` at `f1a9e0d0d76d4be565ac9355a09f77a33a3338a7` after exact-head Fast + Qualified success; conserved 296-row authority ledger, legacy test/gate inventory and first safe capability-oriented test organization wave complete.
- No Stable was manufactured for v18.6.7; certified Stable remains `v18.6.1-stable`.

## 4. v18.7.0 frozen scope and owners

Primary contracts:

- `v18_7_0_scope.json`
- `v18_7_0_g0_g3_contract.json`
- `adaptive-governance/V18_7_0_RUNTIME_RELIABILITY_AUDIT.md`
- `release/v18.7.0/release_contract.json`

### Preserve existing owners

- Smart Provider Router v2 — executable provider/capability routing and health.
- `data_freshness.go` — session/provider/dataset freshness truth.
- `RuntimeSLOTracker` — runtime SLO/recovery owner.
- `WorkloadController` — bounded queues, tiers, reserved critical capacity and shedding.
- `ProviderTelemetry` / runtime observability — request budget/load metrics.
- `BroadSnapshotBroker` — reuse/coalescing/single-flight for broad snapshots.
- Provider Reconciliation — independent-source agreement/conflict owner.
- `Engine.Snapshot()` — canonical runtime aggregation boundary.

No parallel retry, health, freshness, reconciliation, degradation or routing engine may be introduced.

## 5. v18.7 implementation packets

### Packet A — Degradation truth — IMPLEMENTED ON BRANCH

- canonical `ReasonCode` taxonomy while keeping compatible display `Code`;
- `HEALTHY` / `PROTECTED` / `DEGRADED` / `RECOVERING` pressure states;
- explicit affected datasets and downstream consumers;
- `DecisionImpact` semantics;
- fail-closed `UNKNOWN` + `ABSTAIN` when required decision evidence is insufficient;
- immediate `QUEUE_SATURATED` when provider queue reaches its hard bound.

### Packet B — Recovery hysteresis — IMPLEMENTED ON BRANCH

- active degradation surfaces immediately;
- clear only after 3 consecutive healthy observations **and** >=5 seconds stability;
- relapse resets the recovery streak;
- runtime SLO exposes pending confirmations rather than falsely reporting recovery.

### Packet C — Existing reliability controls — AUDITED / REUSE

Requalify rather than rebuild:

- provider/capability state, entitlement, circuits and fallback;
- bounded provider route/retry behavior;
- session-aware freshness;
- request budgets and low-tier shedding;
- cross-provider disagreement;
- snapshot reuse/coalescing.

### Packet D — Active-market reliability proof — STAGED, NOT YET PASS

`v18_7_0_active_market_reliability_test.go` uses production owners to prove:

- 16 equivalent concurrent broad-snapshot requests collapse to one provider fetch;
- coalescing is observable;
- provider concurrency + queue remain hard bounded;
- queue saturation becomes explicit degradation immediately;
- critical current evidence can remain `PROTECTED` while capacity SLO blocks;
- work beyond the hard capacity/queue boundary is not admitted.

Do not call this PASS until CI executes it.

### Packet E — Release reproducibility — IMPLEMENTED ON BRANCH

- Fast, Qualified and Release external Actions immutable-SHA pinned and dependency-lock governed;
- Release G12 Playwright uses the same pinned requirements + safe pip-cache contract as Qualified;
- reproducibility gate enforces all three workflows;
- no new workflow or trigger architecture introduced.

### Packet F — Release identity/certification — IMPLEMENTED ON BRANCH, NOT STABLE

- candidate identity `18.7.0` / `v18.7.0-stable-20260819`;
- small `renderer/release-identity-v18.7.0.js` last-loaded overlay avoids rewriting the legacy renderer monolith solely for identity;
- v18.6.1 watchlist extension remains behavior owner underneath;
- `release/v18.7.0/run_full_certification.sh` is current-source aware;
- exact-head Qualified remains mandatory WebKit proof before merge; G12 does not manufacture a redundant second qualification workflow.

## 6. Pre-PR completion checklist

Before opening the single Draft PR:

1. source-owner audit durable and aligned with frozen G1 scope;
2. release identity internally consistent;
3. release script references current test paths;
4. Release workflow/dependency lock/reproducibility gate agree;
5. four Adaptive overlays + authoritative handoff point to v18.7;
6. exact branch diff contains no unrelated product scope;
7. no existing v18.7 PR exists;
8. current head recorded for candidate intake.

Static/pre-PR inspection is not a substitute for CI. Fast begins when the Draft PR opens.

## 7. PR / qualification sequence

1. Open **exactly one Draft PR** from `v18.7.0-development` to `main`.
2. Require exact-head `DE.PULSE/fast-head = success`.
3. Fix any legitimate failure on the same branch/PR; do not create retry branches.
4. Mark the same PR Ready only after Fast is green.
5. Require **full Qualified** because v18.7 changes backend reliability, tests, release identity and Release workflow.
6. Full Qualified must include backend/full Go, race, randomized, renderer/deterministic, Chrome and WebKit.
7. Run/retain G10 conserved-ledger and release/reproducibility evidence.
8. Merge only when the exact current head is green.

## 8. G11–G16 Stable sequence

After merge of the exact qualified source head:

- G11 verifies exact source-head Fast + Qualified statuses and source-head → merged-candidate fingerprint equivalence.
- G12 runs `release/v18.7.0/run_full_certification.sh` from the immutable merged candidate.
- G13/G14 package/audit macOS Apple Silicon + Windows x64 in parallel.
- G15 binds both native evidence graphs and exact artifact hashes.
- Publish uploads the exact certified artifacts from that run; no rebuild.
- G16 records final v18.7 handoff/evidence.

Only then may `v18.7.0-stable` be called the current Stable.

## 9. Failure handling

- `PRODUCT_FAIL`: source fix, same branch/PR.
- `GATE_TEST_FAIL`: fix defective source/test/gate without weakening a valid contract.
- `CI_HARNESS_FAIL`: same branch/PR.
- `INFRA_FAIL`: bounded unchanged-SHA retry only when recovery signal exists.
- `EXPECTED_NOOP`: record intentional skip/idempotency.
- `SUPERSEDED`: cancel/ignore obsolete exact-head work.

A source-changing fix creates a new exact head and requires new exact-head evidence.

## 10. Next version-visible sequence

- `v18.7.1` — only if qualification finds a focused user-trust reliability patch.
- `v18.8.0` — Shared Intelligence Consolidation.
- `v18.8.1` — Renderer Modularization II.
- `v18.9.0` — TradeInsight SHADOW Integration.
- `v18.9.1` — Provider Intelligence & Market-Mode Hardening.
- `v18.10.0` — v18 Major Closure Candidate.
- `v18.10.1` — closure patch only if justified.
- v19 Professional Data Infrastructure; v20 Adaptive Intelligence maturation.

## 11. Quality floor

Efficiency may reduce duplicate work/files/runs; it may never reduce exact-source provenance, deterministic market truth, unique regression coverage, Chrome + WebKit evidence when required, provider/data/security/rights controls, conserved requirement traceability, required native Stable proof, same-artifact publication, No Execution or any permanent product boundary.
