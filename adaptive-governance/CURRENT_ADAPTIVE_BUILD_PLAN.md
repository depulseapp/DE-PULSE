# DE.PULSE — Current Adaptive Build Plan

**Immutable Stable:** `v18.9.0-stable`  
**Stable candidate:** `9ea81cddae4875ae15d3719ca028519a36c597b6`  
**Stable fingerprint:** `a8719090c341c874dbd1279cc31ad98e84075d5701c46a800bf951340780ecb9`  
**Active product development branch:** none  
**Active product PR:** none  
**Governance alignment PR:** #67 (merged)  
**Master corrective program:** #65 / `ADAPT-PROVIDER-INTELLIGENCE-010`  
**Hosted architecture program:** #66 / `ADAPT-HOSTED-SYNC-001`  
**Immediate next product patch:** `v18.9.1` / #64 / `ADAPT-RUNTIME-CRASH-001`.

## 1. Build philosophy

DE.PULSE uses one permanent G0-G16 release model and small, dependency-ordered, independently certifiable patches.

`stabilize -> canonicalize -> instrument -> validate -> expand -> operationalize -> close -> hosted rights/identity/entitlement/privacy/IaC/data foundations -> provider gateway/serving/sync -> cross-platform capability -> equivalence certification -> next cross-platform capability -> mixed-client hardening -> point-in-time evidence -> reliability -> governed learning`

Permanent rules:
- one primary responsibility per patch;
- `REUSE -> CONSOLIDATE -> REFACTOR -> DELETE/REPLACE -> ADD` before new machinery;
- canonical owners are extended, never forked for a provider/client/feature;
- observability needed to judge a capability lands before broad admission;
- hosted identity/security/rights/product-entitlement/privacy/IaC/recovery/secrets/supply-chain foundations precede commercial activation;
- migrations have rollback/roll-forward/recovery disposition before activation;
- known misses are fixed in-scope or durably assigned;
- G0-G16 is the only top-level release model.

## 2. Cross-Platform Lockstep Product Contract

DE.PULSE is one product across **macOS, Windows and Web**.

Every G1 freezes:
- `macOS: REQUIRED | N/A(reason)`
- `Windows: REQUIRED | N/A(reason)`
- `Web: REQUIRED | N/A(reason)`

For a shared capability:
- one canonical domain/API/state contract drives all REQUIRED clients;
- all REQUIRED platform adapters/surfaces are part of the same release responsibility;
- platform adapters may differ only where OS/browser mechanics require it;
- business logic, intelligence, authorization, product entitlement, provider-right decisions, state semantics, provenance/freshness and explanation meaning may not fork by platform;
- one-platform technical validation is diagnostic only, never a product pilot;
- no shared capability is Delivered/GA until all REQUIRED platforms pass;
- no next shared domain starts while the current shared domain has material parity debt;
- temporary platform exceptions require an external blocker, explicit waiver/expiry, no misleading GA status and a durable recovery release.

Platform-specific corrective work is allowed when the responsibility itself is platform-specific, such as the macOS `v18.9.1` crash.

## 3. Version semantics

- Major = strategic maturity generation.
- Minor band = coherent dependency phase.
- Patch = one primary independently certifiable responsibility. Required platform adapters for one shared capability are part of that responsibility.
- Future reservations may shift/split before G1; shipped versions are immutable.

## 4. v18.9.x — Trustworthy Native Runtime/Data Plane

1. **v18.9.1 — Runtime Reliability** — #64 macOS SIGABRT corrective only; preserve SQLite/user state/API keys; warm-state/relaunch regression; packaged Apple Silicon proof.
2. **v18.9.2 — TradeInsight Settings/API-key UX** — canonical Settings/local secret owner; masked Save/Test/Clear; truthful status; no routing expansion.
3. **v18.9.3 — Coverage-Aware Smart Provider Router Core** — memory/DB reuse first, residual-gap acquisition, deterministic merge/provenance/coverage re-evaluation.
4. **v18.9.4 — Canonical Company/Instrument Identity** — one identity owner across Desks/Research/Discovery/Add Symbol.
5. **v18.9.5 — Market Data Modes + Capability Diagnostics** — behavior/quality modes with truthful source/freshness/coverage/fallback state.
6. **v18.9.6 — Provider Observability / Adaptive Telemetry** — coverage, calls avoided, usefulness, latency/errors/rate/freshness, disagreement, runtime pressure and protected-session headroom.
7. **v18.9.7 — TradeInsight SEC Form 4 Enrichment** — SHADOW-first; direct SEC/EDGAR authoritative.
8. **v18.9.8 — TradeInsight Symbol/Company Search** — canonical symbol-validation/company-identity fallback/corroboration.
9. **v18.9.9 — TradeInsight Movers/Ranking Evidence** — SHADOW evidence through Opportunity Radar; no parallel scanner/ranker.
10. **v18.9.10 — Remaining Useful TradeInsight Capability Admission** — explicit consumer/owner/lifecycle/serving role/freshness/retention/rights/rate-cost disposition.
11. **v18.9.11 — Session-Aware Data Readiness Maintenance** — bounded overnight/weekend maintenance with protected-session reserves and preemption/checkpoint/resume.
12. **v18.9.12 — Professional Closure** — no feature scope; implementation-miss/duplicate-owner/provider/runtime/package closure.

## 5. v19 — Professional Hosted Product

**Entry:** `v18.9.12` PASS.

Architecture target:

`Mac/Windows native SQLite edge + Web -> DE.PULSE hosted APIs -> canonical policy/router/state -> PostgreSQL shared authority`

Normal commercial users are zero-key; platform provider credentials remain server-side.

Hosted serving order:

`tenant/account -> session/device -> RBAC/capability -> DE.PULSE product entitlement -> provider legal/data rights -> privacy/data policy -> canonical cache/persistence/freshness -> residual need -> Smart Provider Router v2 -> authorized projection`

### v19.0.x — Governance / Control Plane / Data Foundation

- `v19.0.0` Provider Capability / Legal Rights Registry
- `v19.0.1` Hosted Tenant / Identity / Device / Session Control Plane
- `v19.0.2` DE.PULSE Product Entitlement / Metering Policy
- `v19.0.3` Account Data Governance / Privacy Lifecycle
- `v19.0.4` Hosted Environment / IaC / Service Trust Foundation
- `v19.0.5` PostgreSQL Tenancy / Schema / Pool / HA-PITR Foundation
- `v19.0.6` Managed Secrets / KMS Lifecycle
- `v19.0.7` Software Supply-Chain / Artifact & Dependency Assurance
- `v19.0.8` Provider Quality / Cost / Coverage / SLO Scorecards
- `v19.0.9` Data Reconciliation / Revision / Point-in-Time Quality

**Exit:** all required hosted foundations are executable/evidenced before shared client activation.

### v19.1.x — Hosted Data Plane + Cross-Platform Account/State Foundation

- `v19.1.0` **Authenticated Hosted Provider Gateway** — versioned zero-key server boundary around Smart Provider Router v2; canonical freshness/cache/persistence reuse; rate/backpressure/circuit/kill controls; API inventory/version/deprecation ownership.
- `v19.1.1` **Unified Serving Policy + Live Fan-Out Isolation** — enforce tenant/RBAC/product-entitlement/provider-right/privacy checks across REST/WebSocket/SSE/cache/persistence; existing multi-feed owner remains canonical.
- `v19.1.2` **Sync Protocol Foundation** — stable IDs, bootstrap/high-watermark, native outbox, idempotent push, authoritative server sequence/change log, incremental pull, checkpoints, tombstones, compaction/re-bootstrap and mixed-version negotiation.
- `v19.1.3` **Cross-Platform Account / Session Client Foundation — Mac + Windows + Web** — login/session/device/role/capability/product-entitlement client semantics against one server truth.
- `v19.1.4` **Cross-Platform Preferences — Mac + Windows + Web** — account-scoped preferences, offline/reconnect where applicable, user switching and account lifecycle equivalence.
- `v19.1.5` **Cross-Platform Watchlists / Master Symbols — Mac + Windows + Web** — same membership/add/remove/remove-all semantics, conflicts, sync/convergence and canonical symbol identity.
- `v19.1.6` **Cross-Platform Desks / Workspaces — Mac + Windows + Web** — versioned membership/configuration/delete/history semantics through the same canonical owners.

**Exit:** no macOS pilot and no Windows/Web catch-up phase. All REQUIRED clients pass each portable account/state capability before the next shared domain proceeds.

### v19.2.x — Cross-Platform Shared Product + Hosted Assurance

- `v19.2.0` **Cross-Platform Research / Durable State — Mac + Windows + Web** — same canonical Research structure, lawful/provenance-bound durable artifacts, retention and authorization.
- `v19.2.1` **Cross-Platform Discovery / Opportunity Radar — Mac + Windows + Web** — same candidate/ranking/evidence/outcome meaning; no client-specific scanner/ranker.
- `v19.2.2` **Cross-Platform Market State / Market Modes / Readiness / Explanations — Mac + Windows + Web** — same regime/mode/readiness/explanation semantics and source/freshness truth.
- `v19.2.3` **Cross-Platform Settings / Account / Device Controls — Mac + Windows + Web** — same account/device/settings truth; client-specific secure storage/session mechanics allowed.
- `v19.2.4` **Cross-Platform RBAC / Product-Entitlement UX — Mac + Windows + Web** — same role-aware navigation, capability visibility, downgrade/suspension behavior and direct-route/API denial; UI hiding never authorization.
- `v19.2.5` **Tenant-Aware Metering / Cost / Usage Observability** — per tenant/account/user/device/capability attribution, quota consumption, calls avoided, stream usage, provider cost where known and platform dimension without platform-specific business truth.
- `v19.2.6` **Mixed-Client Multi-User Security / Abuse / Capacity Hardening** — concurrent Mac/Windows/Web load, object/function authorization negatives, fairness/noisy-neighbor isolation, rate/circuit/load shedding, DB/pool limits and protected-session capacity.
- `v19.2.7` **#66 Cross-Platform Hosted Sync / Gateway Assurance Closure** — no feature scope; explicit cross-platform equivalence plus adversarial/failure/recovery/privacy/environment/supply-chain matrix.

### v19.3.x — Point-in-Time Evidence Substrate

- `v19.3.0` Institutional / 13F Evidence Infrastructure
- `v19.3.1` Two-Sided Long / Short Thesis Evidence Substrate
- `v19.3.2` AODR Candidate / Ranking / Outcome Lineage

Any user-facing capability surfaced from these substrates obeys the lockstep contract.

### v19.4.x — Reliability / Economics / v20 Readiness

- `v19.4.0` ADR-GDI Professional Reliability / Capacity + operator runbooks
- `v19.4.1` Specialized / Paid Provider Gap Evaluation
- `v19.4.2` v20 Research-Readiness Dataset / Lineage Audit

### v19.5.0 — v19 Major Closure

No feature scope. Require #66 PASS, zero material Mac/Windows/Web parity debt for shared capabilities, tenant/RBAC/product-entitlement/provider-right/privacy separation, data lifecycle, IaC/environment assurance, supply-chain provenance, API compatibility, recovery/rollback, SLO/capacity and actual supported artifact/deployment evidence.

## 6. v20 — Adaptive Intelligence & Decision Research

**Entry:** `v19.5.0` PASS.

- `v20.0.0` Adaptive Research Control Plane + Immutable Experiment Ledger
- `v20.0.1` Model / Prompt Governance + Champion/Challenger
- `v20.0.2` Historical Analogues + Regime-Conditioned Outcomes
- `v20.0.3` Calibration / FP-FN / Miss / Contradiction / Drift
- `v20.1.0` ASBI Behavioral Fingerprints + State Transitions
- `v20.1.1` ASBI Scenarios / Probability Momentum / Calibration
- `v20.2.0` Adaptive Institutional / 13F Intelligence
- `v20.2.1` TDTI Competing Long / Short / No Reliable Edge
- `v20.2.2` TDTI Two-Sided Trade-Plan Validation — No Execution
- `v20.3.0` AODR Adaptive Shared Opportunity Ranking
- `v20.3.1` AODR Diversity / Opportunity Cost / Personalized Relevance after shared truth
- `v20.4.0` ADR-GDI Adaptive Optimization under SHADOW/Champion-Challenger
- `v20.5.0` Professional Closure

Every shared adaptive user-facing capability is implemented and promoted Mac + Windows + Web in lockstep.

## 7. G0-G16 Lockstep Enforcement

- **G1:** freeze platform applicability matrix.
- **G2:** map canonical owner + platform adapters; no duplicate client business truth.
- **G3:** one domain/API/state contract + adapter contracts + equivalence tests.
- **G4:** all REQUIRED client implementations for shared scope before Development Exit.
- **G5/G6:** affected-platform deterministic + cross-platform integration tests.
- **G7:** equivalent authorization/data/privacy/provider-right outcomes.
- **G8:** mixed-client load/capacity/recovery where applicable.
- **G9:** role-aware Mac/Windows/Web function/meaning equivalence; responsive/native interaction differences allowed.
- **G10:** material parity debt blocks freeze.
- **G13/G14:** actual required native/Web artifacts/deployments/runtime proof.
- **G15:** no GA/Delivered state until all REQUIRED clients pass.
- **G16:** parity-drift/waiver audit and exact handoff.

## 8. Protected boundaries

Smart Provider Router v2 sole routing authority; canonical freshness/recovery sole freshness owner; existing multi-feed allocator sole subscription owner; BroadSnapshotBroker canonical reuse owner; canonical persistence/cache/state/session/calendar/identity owners; direct SEC/EDGAR authoritative; deterministic Day/Swing/Long protected; U.S. Equities Processing; GLD/SLV/USO actionable exceptions; No Execution; G0-G16 only.

## Exactly one next action

Execute G0 for #64 / `v18.9.1` using complete macOS crash evidence or deterministic reproduction, then freeze its narrow G1. Do not start `v18.9.2` or v19 implementation until ordering permits it.