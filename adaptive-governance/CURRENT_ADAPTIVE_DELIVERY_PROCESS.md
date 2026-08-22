# DE.PULSE — Current Adaptive Delivery Process

**Certified Stable:** `v18.9.0-stable`  
**Active product branch/PR:** `v18.9.1-development` / Draft PR #69  
**Active product scope:** #64 / `ADAPT-RUNTIME-CRASH-001` -> `v18.9.1`  
**Mandatory next delivery/process hardening:** #70 / `ADAPT-CI-CONVERGENCE-001` after truthful v18.9.1 closure  
**Master corrective program:** #65 / `ADAPT-PROVIDER-INTELLIGENCE-010`  
**Hosted architecture program:** #66 / `ADAPT-HOSTED-SYNC-001`  
**Governance alignment:** merged PR #67.

## 1. Delivery invariant

A DE.PULSE release is not delivered because code compiles, CI is green, one platform works, or governance exists.

`Governed -> Implemented -> Enforced -> Evidenced -> Packaged/Deployed -> Cross-Platform Verified -> Delivered -> Learned`

Release engineering is part of delivery quality. CI efficiency may remove redundant computation, but may not remove required evidence.

For process/repository-hardening work, documentation does not equal implementation. A scope is Delivered only when the repository/workflow/tooling changes exist and the enforcement/evidence proves them.

## 2. Cross-Platform Definition of Delivered

For every G1 capability, freeze `macOS / Windows / Web = REQUIRED or justified N/A`.

A shared capability is Delivered only when:
1. canonical domain/API/state implementation passes;
2. every REQUIRED client adapter/surface passes;
3. equivalent authorization, state, data, freshness/provenance and user-meaning tests pass;
4. required macOS/Windows artifacts and Web/server deployment evidence pass;
5. no material parity debt remains.

There is no normal state of `Mac done -> Windows later -> Web later`. A single platform may be used for diagnosis or technical validation, but not as a product pilot or GA milestone.

Platform-specific releases remain valid when the actual responsibility is platform-specific; #64 / `v18.9.1` is the current example.

Temporary platform exceptions require a proven external blocker, explicit waiver/expiry, no misleading Delivered/GA status and a named recovery release.

## 3. Public release vs work-slice delivery

DE.PULSE must not use public product versions merely as internal work-item counters.

After #70, delivery tracks independently:
- public `productVersion`;
- `workSliceId`/scope IDs included in that release;
- exact candidate SHA and source fingerprint;
- immutable build ID;
- evidence schema version;
- monotonic platform build number.

Multiple independently governed work slices may be grouped into one coherent public release when compatibility, dependency order and release risk allow it. Urgent reliability/security/bug corrections may still ship independently as PATCH releases. Shipped versions/tags/evidence are immutable.

## 4. Immediate delivery order

### A. Complete `v18.9.1` / #64

G0/root-cause isolation is complete. Delivery remains blocked on trustworthy exact-head qualification.

Required before PR #69 can leave Draft:
- hosted-runner execution restored/confirmed;
- exact-head Fast PASS;
- full exact-head Qualified PASS;
- backend, renderer, Chrome and WebKit required evidence PASS;
- actual packaged macOS fresh native startup and bounded dwell PASS;
- warm relaunch using the same profile PASS;
- SQLite/profile reuse PASS;
- no `protocol does not exist` evidence;
- deterministic native cleanup PASS.

Merge/release still requires explicit authorization and canonical G11-G16.

### B. Execute #70 / `ADAPT-CI-CONVERGENCE-001`

After truthful v18.9.1 closure, #70 is the mandatory next delivery/process workstream before any v18.9.2+ product implementation.

#70 is Delivered only when executable evidence proves:
- exactly three routine workflows remain;
- adaptive job selection is dependency-aware and fail-closed;
- manual/reusable runs bind the complete intended delta;
- targeted native rehearsals occur before final Release when impact requires them;
- WebKit/browser and native-lifecycle responsibilities are explicit and reuse canonical packaging owners;
- future G12 has one version-neutral executor;
- active version-named tests are migrated under governed capability ownership without regression-coverage loss;
- impossible publication fails at G11 before expensive work;
- Stable publication is digest-immutable and same-digest idempotent;
- Stable releases are globally serialized;
- actual repository protection/ruleset evidence exists;
- product/work-slice/source/build/evidence identities are separated;
- future version/tag semantics are explicitly SemVer or explicitly custom;
- platform build numbering is monotonic and collision-safe;
- final runnable artifact provenance/toolchain identity is durable;
- handoff and Adaptive overlays share actual current-state truth;
- G16 reports runner minutes/reruns avoided and confirms no assurance was removed;
- a safe repository-root `.gitignore` exists without hiding tracked continuity checkpoints;
- source-health recursively covers every production Go package before any package relocation;
- one canonical root-layout inventory/allowlist owner exists;
- stale root `certification_plan.json` / `certification_runner.py` / `ci_pipeline.py` / `ci_pipeline_plan.json` no longer act as a competing current CI/release model;
- reusable root tooling is owned under stable `tools/` paths with consumers migrated atomically;
- historical/version-scoped root material is moved/retired only with evidence mapping and without altering shipped Stable history;
- package-local Go test coverage remains intact through capability-oriented renames and any later package extraction;
- retained assets/policies/registries have explicit stable owners and passing consumers;
- a CI root allowlist prevents the version-stacked root problem from returning;
- before/after root inventory and moved/renamed/removed file counts are recorded.

**Documentation-only completion is explicitly insufficient.** #70 cannot be marked Delivered or closed until those repository/code/workflow/tooling changes exist and applicable Fast/Qualified/full evidence passes.

Detailed contract: `adaptive-governance/CURRENT_ADAPTIVE_CI_CONVERGENCE.md` and issue #70.

### C. Re-baseline future product release grouping

Only after #70 passes should the remaining v18.9.x reservations be confirmed, grouped or shifted prospectively. Dependency/work-slice order remains authoritative; historical releases are untouched.

### D. Resume product delivery

Then resume the next product capability in dependency order.

## 5. Stable release delivery contract after #70

Permanent workflow topology:

`Fast -> Qualified -> Release`

No fourth workflow family.

Delivery sequence:

`work slice(s) -> exact-head qualification -> immutable release candidate -> G11 publication feasibility -> canonical G12 -> G13/G14 actual artifacts/deployments -> G15 exact-artifact assurance -> immutable publication -> G16 retrospective/handoff`

Rules:
- G11 rejects impossible version/predecessor/tag/publication states before expensive G12/G13/G14;
- final publication uses exact same-run certified artifacts and never rebuilds;
- existing Stable asset with same digest = idempotent no-op/reuse;
- differing existing Stable digest = release integrity failure, never overwrite;
- only one Stable G11-G16 graph may publish at a time;
- candidate/toolchain/artifact provenance is retained durably;
- exact-head-invalid evidence is not reusable;
- repository-layout enforcement remains part of permanent CI so version-stacked root clutter cannot recur.

## 6. Remaining v18.9.x product work — current reservations

After #70 re-baselining, preserve this dependency order:

1. TradeInsight Settings/API-key UX — currently reserved `v18.9.2`
2. coverage-aware Smart Provider Router — reserved `v18.9.3`
3. canonical company/instrument identity — reserved `v18.9.4`
4. Market Data Modes/diagnostics — reserved `v18.9.5`
5. provider observability/Adaptive telemetry — reserved `v18.9.6`
6. Form 4 SHADOW enrichment — reserved `v18.9.7`
7. symbol/company search — reserved `v18.9.8`
8. movers/ranking SHADOW evidence — reserved `v18.9.9`
9. remaining useful capability admission — reserved `v18.9.10`
10. Session-Aware Data Readiness Maintenance — reserved `v18.9.11`
11. Professional Closure — reserved `v18.9.12`

These are work/release reservations, not a requirement for eleven separate future Stable releases.

## 7. Protected Tier-0 contract

Pre-market, regular market and after-hours are protected. Live/current work has first provider/runtime/DB/worker claim. Maintenance/sync is bounded, preemptible/checkpointed and cannot flood protected sessions.

## 8. v19 delivery train

**Entry:** v18 native/data-plane closure + #70 executable CI/versioning/repository-hygiene PASS.

### v19.0.x — Hosted Foundations
Provider legal-rights registry; tenant/identity/device/session control plane; product entitlement/metering; account data governance/privacy; IaC/service trust; PostgreSQL HA/PITR; managed secrets/KMS; software supply-chain assurance; provider quality/cost/SLO scorecards; point-in-time reconciliation.

### v19.1.x — Hosted Data Plane + Cross-Platform Account/State
Hosted Provider Gateway; Unified Serving Policy + Live Fan-Out; Sync Protocol Foundation; Mac + Windows + Web Account/Session, Preferences, Watchlists/Master Symbols and Desks/Workspaces.

No Mac product pilot exists. Each shared capability closes across all REQUIRED clients before the next shared capability proceeds.

### v19.2.x — Cross-Platform Shared Product + Assurance
Mac + Windows + Web Research, Discovery/Opportunity Radar, Market State/Modes/Readiness/Explanations, Settings/Account/Device Controls and RBAC/Product-Entitlement UX; tenant-aware metering; mixed-client security/abuse/capacity; #66 Cross-Platform Assurance Closure.

### Required #66 evidence

At minimum:
- same account/session/role/product-entitlement truth across clients;
- same authorized preferences/watchlists/desks/workspaces/research state;
- same Discovery/Opportunity Radar and Market Mode/readiness/explanation meaning for equivalent evidence;
- same provider/freshness/provenance/degraded/UNKNOWN truth;
- same revocation/downgrade/denial outcomes;
- native offline/reconnect convergence without Web becoming a second authority;
- mixed-client concurrent edit/conflict/delete/tombstone behavior;
- provider-right/product-plan downgrade behavior;
- account export/deactivation/deletion/retention behavior;
- secret rotation and DB failover/PITR/restore;
- mixed-version protocol behavior;
- environment/config drift and artifact/dependency provenance failure behavior;
- noisy-neighbor/fairness/backpressure/protected-session pressure;
- no provider secret or unnecessary user-data leakage.

### v19.3.x — Point-in-Time Evidence
Institutional/13F -> Two-sided Long/Short evidence -> AODR lineage.

### v19.4.x — Reliability / Economics / v20 Readiness
ADR-GDI reliability/capacity -> specialized/paid-provider gap evaluation -> v20 research-readiness audit.

### v19.5.0 — Major Closure
No feature scope. Require #66 PASS, zero material shared-capability parity debt, rights/privacy/security/IaC/supply-chain/API/recovery/SLO/capacity and actual required Mac/Windows/Web runtime/deployment evidence.

## 9. G0-G16 delivery enforcement

- G1 freezes scope/work-slice identity, product-version disposition and REQUIRED/N/A platform matrix.
- G4 requires every REQUIRED client implementation for shared scope and requires actual executable changes for process/repository-hardening work.
- G5/G6 may select affected evidence adaptively, but ambiguity fails closed.
- G10 blocks freeze on material parity debt, insufficient evidence mapping, competing CI/release owners, lost package/test coverage or unenforced root-layout migration.
- G11 freezes source/build/release identity and proves publication feasibility.
- G12 uses canonical full certification.
- G13/G14 validate actual required native/Web artifacts/deployments.
- G15 forbids shared-capability GA until all REQUIRED clients pass and forbids differing Stable-byte overwrite.
- G16 audits parity drift, CI efficiency, provenance, root before/after inventory, moved/renamed/removed paths and temporary waivers.

## 10. v20 delivery train

`adaptive control/model governance -> ASBI -> adaptive Institutional/TDTI -> AODR -> adaptive operations -> Professional Closure`

Every shared adaptive user-facing capability follows Mac + Windows + Web lockstep and `SHADOW -> VALIDATED -> APPROVED -> PRODUCTION`. No Execution remains permanent.

## 11. Permanent boundaries

U.S. Equities Processing; GLD/SLV/USO actionable exceptions; No Execution; Smart Provider Router v2 sole routing authority; canonical freshness/recovery sole freshness owner; existing multi-feed allocator sole subscription owner; BroadSnapshotBroker canonical reuse owner; canonical persistence/cache/state/session/calendar/identity owners; direct SEC/EDGAR authoritative; G0-G16 only.

## Exactly one next action

Restore/confirm GitHub Actions hosted-runner execution and finish exact-head #64 / `v18.9.1` qualification. After truthful closure, #70 is the mandatory executable next delivery/process workstream—including actual repository-root cleanup and enforcement—before any next product implementation.
