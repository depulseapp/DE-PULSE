# DE.PULSE — Current Adaptive Build Plan

**Immutable Stable:** `v18.9.0-stable`  
**Stable candidate:** `9ea81cddae4875ae15d3719ca028519a36c597b6`  
**Stable fingerprint:** `a8719090c341c874dbd1279cc31ad98e84075d5701c46a800bf951340780ecb9`  
**Active product development branch:** `v18.9.1-development`  
**Active product PR:** Draft PR #69  
**Active product scope:** #64 / `ADAPT-RUNTIME-CRASH-001` -> `v18.9.1`  
**Approved mandatory process-hardening workstream:** #70 / `ADAPT-CI-CONVERGENCE-001` — immediately after truthful v18.9.1 closure  
**Governance alignment PR:** #67 (merged)  
**Master corrective program:** #65 / `ADAPT-PROVIDER-INTELLIGENCE-010`  
**Hosted architecture program:** #66 / `ADAPT-HOSTED-SYNC-001`

## 1. Build philosophy

DE.PULSE uses one permanent G0-G16 release model and small, dependency-ordered, independently governed work slices.

`stabilize -> harden build/release machinery -> canonicalize -> instrument -> validate -> expand -> operationalize -> close -> hosted foundations -> provider gateway/serving/sync -> cross-platform capability -> equivalence certification -> mixed-client hardening -> point-in-time evidence -> reliability -> governed learning`

Permanent rules:
- one primary responsibility per work slice;
- `REUSE -> CONSOLIDATE -> REFACTOR -> DELETE/REPLACE -> ADD` before new machinery;
- canonical owners are extended, never forked for a provider/client/feature;
- observability needed to judge a capability lands before broad admission;
- migrations have rollback/roll-forward/recovery disposition before activation;
- known misses are fixed in-scope or durably assigned;
- G0-G16 is the only top-level release model;
- exactly three routine workflow families remain: Fast, Qualified and Release;
- CI efficiency may reduce redundant work, never required evidence.

## 2. Cross-Platform Lockstep Product Contract

DE.PULSE is one product across **macOS, Windows and Web**.

Every G1 freezes:
- `macOS: REQUIRED | N/A(reason)`
- `Windows: REQUIRED | N/A(reason)`
- `Web: REQUIRED | N/A(reason)`

For a shared capability:
- one canonical domain/API/state contract drives all REQUIRED clients;
- all REQUIRED platform adapters/surfaces are part of the same governed capability responsibility;
- platform adapters may differ only where OS/browser mechanics require it;
- business logic, intelligence, authorization, product entitlement, provider-right decisions, state semantics, provenance/freshness and explanation meaning may not fork by platform;
- one-platform technical validation is diagnostic only, never a product pilot;
- no shared capability is Delivered/GA until all REQUIRED platforms pass;
- no next shared domain starts while the current shared domain has material parity debt;
- temporary platform exceptions require an external blocker, explicit waiver/expiry, no misleading GA status and a durable recovery release.

Platform-specific corrective work is valid where the responsibility itself is platform-specific, such as `v18.9.1`.

## 3. Build identity and version semantics

Shipped DE.PULSE versions/tags/evidence are immutable.

For future work, #70 must separate:
- `productVersion` — public/user-visible compatibility version;
- `workSliceId` — independently governed engineering responsibility;
- `candidateSha` — exact Git source;
- `sourceFingerprint` — exact relevant source identity;
- `buildId` — immutable build/evidence identity;
- `evidenceSchemaVersion` — evidence contract identity;
- platform build number — monotonic packaging identity.

Current v18.9.x version numbers after v18.9.1 are **planning reservations**, not a requirement that every work slice become its own public Stable release.

Before the next product implementation, #70 must make one prospective policy decision:
1. adopt Semantic Versioning for future releases; or
2. explicitly declare a custom DE.PULSE release-train versioning scheme.

If SemVer is adopted, backward-compatible features normally group into MINOR releases and PATCH remains corrective/fix scope. Historical tags including `-stable` are not rewritten.

## 4. Mandatory immediate build sequence

### Step A — Complete `v18.9.1` / #64

G0 deterministic reproduction/root cause is already complete. Do not restart it.

Remaining requirements:
- restore/confirm GitHub Actions hosted-runner execution;
- re-fetch exact live branch head and #64/PR #69 state;
- exact-head Fast PASS;
- full exact-head Qualified PASS;
- backend, renderer, Chrome and WebKit required evidence PASS;
- actual packaged macOS fresh native startup + startup dwell PASS;
- warm relaunch using the same persisted profile PASS;
- SQLite/profile reuse PASS;
- no `protocol does not exist` evidence;
- deterministic cleanup PASS;
- PR #69 stays Draft until the required evidence is green.

### Step B — Execute #70 / `ADAPT-CI-CONVERGENCE-001`

This is now a **mandatory build-plan checkpoint**, not optional backlog. It begins only after truthful v18.9.1 closure unless a discovered CI defect makes current evidence structurally untrustworthy.

Required build outcomes:
1. Planner v3 selects deterministic affected jobs/lane requirements with explicit reasons and fail-closed full fallback.
2. Manual/reusable runs bind an explicit/authoritative full-delta base.
3. Native/release harness changes trigger targeted real rehearsal inside existing Qualified.
4. WebKit browser and native lifecycle have explicit owners while reusing canonical native scripts.
5. One canonical version-neutral G12 executor replaces future per-version certification executors.
6. Release identity/version-only source churn is reduced.
7. Active version-named tests migrate incrementally to capability-oriented paths with equivalence proof.
8. G11 performs version/predecessor/tag/publication feasibility before expensive G12/G13/G14.
9. Stable assets are digest-immutable; same bytes are idempotent, differing bytes fail.
10. Stable G11-G16 execution is repository-wide serialized.
11. Actual `main` ruleset/branch-protection enforcement is verified and hardened if absent.
12. Public version/work-slice/source/build/evidence identities are formally separated.
13. Future version/tag semantics are explicitly SemVer or explicitly custom; history remains immutable.
14. Native build numbers are monotonic and collision-tested.
15. Final runnable artifact provenance/attestation/SBOM is added at the supply-chain milestone.
16. One canonical toolchain manifest records resolved toolchain/runner identity.
17. Handoff and four CURRENT Adaptive overlays consume/check one actual current-state truth.
18. G16 reports runner minutes/reruns avoided and confirms no quality requirement was removed.

Detailed contract: `adaptive-governance/CURRENT_ADAPTIVE_CI_CONVERGENCE.md` and issue #70.

### Step C — Re-baseline the remaining v18.9.x product reservations

After #70 passes, review the remaining planned capabilities as **work slices first** and public releases second. Preserve dependency order and scope IDs even if multiple slices are grouped under one future public version.

### Step D — Resume product work

Only after Steps A-C may the next product capability begin.

## 5. Remaining v18.9.x product work — current reservations

Dependency order remains:

1. TradeInsight Settings/API-key UX — currently reserved `v18.9.2`
2. Coverage-Aware Smart Provider Router Core — reserved `v18.9.3`
3. Canonical Company/Instrument Identity — reserved `v18.9.4`
4. Market Data Modes + Capability Diagnostics — reserved `v18.9.5`
5. Provider Observability / Adaptive Telemetry — reserved `v18.9.6`
6. TradeInsight SEC Form 4 Enrichment — reserved `v18.9.7`; SHADOW-first; direct SEC/EDGAR authoritative
7. TradeInsight Symbol/Company Search — reserved `v18.9.8`
8. TradeInsight Movers/Ranking Evidence — reserved `v18.9.9`; Opportunity Radar remains canonical
9. Remaining Useful TradeInsight Capability Admission — reserved `v18.9.10`
10. Session-Aware Data Readiness Maintenance — reserved `v18.9.11`
11. Professional Closure — reserved `v18.9.12`; no feature scope

Do not interpret this reservation list as a requirement for eleven more public Stable releases. #70 owns that prospective normalization.

## 6. Adaptive CI build contract after #70

The permanent topology remains:

`Fast -> Qualified -> Release`

Planner v3 should emit affected evidence requirements such as:
- backend/core;
- renderer;
- Chrome;
- WebKit;
- persistence/DB/restart/migration;
- security/rights;
- provider/router/data;
- portability/CI harness;
- macOS native rehearsal;
- Windows native rehearsal;
- full fallback.

Unknown, mixed or dependency-ambiguous impact fails closed to full. Cost never suppresses required evidence.

Preferred build flow:

`work slice -> commits -> Fast -> deliberate Qualified checkpoint -> exact candidate SHA/fingerprint -> optional release grouping -> one public product release -> one G11-G16/native publication`

## 7. v19 — Professional Hosted Product

**Entry:** v18 native/data-plane closure + #70 CI/versioning convergence PASS.

Architecture target:

`Mac/Windows native SQLite edge + Web -> DE.PULSE hosted APIs -> canonical policy/router/state -> PostgreSQL shared authority`

Normal commercial users are zero-key; platform provider credentials remain server-side.

Hosted serving order:

`tenant/account -> session/device -> RBAC/capability -> DE.PULSE product entitlement -> provider legal/data rights -> privacy/data policy -> canonical cache/persistence/freshness -> residual need -> Smart Provider Router v2 -> authorized projection`

### v19.0.x — Governance / Control Plane / Data Foundation
- provider legal-rights registry
- tenant/identity/device/session control plane
- product entitlement/metering policy
- account data governance/privacy lifecycle
- hosted environment/IaC/service trust
- PostgreSQL tenancy/schema/pool/HA-PITR
- managed secrets/KMS
- software supply-chain/artifact/dependency assurance
- provider quality/cost/coverage/SLO scorecards
- reconciliation/revision/point-in-time quality

### v19.1.x — Hosted Data Plane + Cross-Platform Account/State
- Hosted Provider Gateway
- Unified Serving Policy + Live Fan-Out
- Sync Protocol Foundation
- Mac + Windows + Web Account/Session Client Foundation
- Mac + Windows + Web Preferences
- Mac + Windows + Web Watchlists/Master Symbols
- Mac + Windows + Web Desks/Workspaces

### v19.2.x — Cross-Platform Shared Product + Assurance
- Mac + Windows + Web Research/Durable State
- Mac + Windows + Web Discovery/Opportunity Radar
- Mac + Windows + Web Market State/Market Modes/Readiness/Explanations
- Mac + Windows + Web Settings/Account/Device Controls
- Mac + Windows + Web RBAC/Product-Entitlement UX
- tenant-aware metering/cost/usage observability
- mixed-client multi-user security/abuse/capacity hardening
- #66 Cross-Platform Hosted Sync/Gateway Assurance Closure

### v19.3.x — Point-in-Time Evidence
Institutional/13F -> Two-Sided Long/Short evidence -> AODR lineage.

### v19.4.x — Reliability / Economics / v20 Readiness
ADR-GDI reliability/capacity -> specialized/paid-provider gap evaluation -> v20 research-readiness audit.

### v19.5.0 — Major Closure
No feature scope. Require #66 PASS, zero material Mac/Windows/Web parity debt, rights/privacy/security/IaC/supply-chain/API/recovery/SLO/capacity and actual supported artifact/deployment evidence.

## 8. v20 — Adaptive Intelligence & Decision Research

`adaptive control/model governance -> ASBI -> adaptive Institutional/TDTI -> AODR -> adaptive operations -> Professional Closure`

Every shared adaptive user-facing capability is implemented and promoted Mac + Windows + Web in lockstep. No Execution remains permanent.

## 9. G0-G16 enforcement additions

- G5/G6 use impact/dependency-aware evidence selection but fail closed when uncertain.
- G10 must prove the candidate's selected evidence graph is sufficient and exact-head.
- G11 must fail impossible publication before expensive certification/native work.
- G12 uses the canonical version-neutral executor after #70.
- G13/G14 validate actual required native/Web artifacts/deployments.
- G15 cannot mutate differing Stable bytes.
- G16 records parity drift, CI efficiency, evidence reuse, toolchain/provenance and remaining cleanup debt.

## 10. Protected boundaries

Smart Provider Router v2 sole routing authority; canonical freshness/recovery sole freshness owner; existing multi-feed allocator sole subscription owner; BroadSnapshotBroker canonical reuse owner; canonical persistence/cache/state/session/calendar/identity owners; direct SEC/EDGAR authoritative; deterministic Day/Swing/Long protected; U.S. Equities Processing; GLD/SLV/USO actionable exceptions; No Execution; G0-G16 only.

## Exactly one next action

Restore/confirm GitHub Actions hosted-runner execution and finish exact-head #64 / `v18.9.1` qualification. After truthful closure, #70 / `ADAPT-CI-CONVERGENCE-001` is the mandatory next build-plan workstream before any v18.9.2+ product implementation.