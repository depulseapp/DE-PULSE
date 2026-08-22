# DE.PULSE — Current Handoff

**SUPERSEDES ALL PRIOR CHAT HANDOFFS**

**GitHub source of truth:** `depulseapp/DE-PULSE`  
**Certified Stable:** `v18.9.0-stable`  
**Certified Stable candidate:** `9ea81cddae4875ae15d3719ca028519a36c597b6`  
**Certified Stable qualified source:** `9e86b5e731f7a585cc77c1521f3639fc7a208efc`  
**Certified Stable fingerprint:** `a8719090c341c874dbd1279cc31ad98e84075d5701c46a800bf951340780ecb9`  
**Certified Stable build ID:** `v18.9.0-stable-20260821`  
**Release PR:** #62 merged  
**Completed scope:** #61 / `ADAPT-TRADEINSIGHT-001`  
**Active product branch/PR:** `v18.9.1-development` / Draft PR #69  
**Active corrective scope:** #64 / `ADAPT-RUNTIME-CRASH-001` / `v18.9.1`  
**Approved mandatory next process-hardening workstream:** #70 / `ADAPT-CI-CONVERGENCE-001`  
**Governance alignment PR:** #67 merged  
**Master corrective program:** #65 / `ADAPT-PROVIDER-INTELLIGENCE-010`  
**Hosted architecture program:** #66 / `ADAPT-HOSTED-SYNC-001`

## Immediate execution rule

Do not start the next product capability until #64 / v18.9.1 is truthfully closed **and** #70 CI/versioning convergence has been executed/re-baselined as required below.

First re-fetch live `main`, `v18.9.1-development`, PR #69, issue #64/comments and issue #70 because concurrent sessions/processes may advance governance or source.

Do **not** restart #64 G0. Deterministic reproduction is complete and the canonical macOS native-window root cause is isolated. Continue from the exact current branch/PR state and qualification blocker.

## v18.9.1 current execution state

The canonical v18.9.0 package audit proved the Stable artifact uses `DePulseLauncher` -> `DePulse-arm64`; the original user crash path `/Applications/De-Pulse.app/Contents/MacOS/De-Pulse` is a distinct/non-matching executable signature and must not be conflated with the canonical reproduction.

Qualified #157 (`32541591746`) executed the real packaged macOS JXA/Cocoa/WKWebView path on macOS 15.7.7 Apple Silicon with pinned Go 1.26.6. Packaged backend/readiness/root and SQLite migrations `[1,2,3,4]` passed, then the native JXA child failed with `protocol does not exist (-2700)`. The only explicit protocol lookup was formal `protocols:['NSApplicationDelegate']` in the JXA registered subclass.

The v18.9.1 corrective candidate removes only that unnecessary formal protocol declaration while preserving the NSObject delegate selectors and `app.delegate=delegate` behavior. Existing macOS packaging and WebKit qualification owners are hardened rather than replaced.

Last runtime/test/harness implementation head before governance/handoff-only updates: `083b69c6772bb7a0fa14a7cdea70f4bd695a10bb`. Re-fetch the branch before relying on any current head SHA.

Current corrective implementation includes:
- root-cause JXA correction in `desktop_lifecycle.go`;
- `macOSWindowScript(...)` extraction plus `v18_9_1_desktop_lifecycle_test.go` regression coverage preventing formal `NSApplicationDelegate` reintroduction while protecting Cocoa/WebKit/delegate/url/icon behavior;
- canonical executable identity checks (`DePulseLauncher`, `DePulse-arm64`, legacy `Contents/MacOS/De-Pulse` absent);
- real packaged non-headless fresh startup, 3-second liveness dwell and warm relaunch on the same profile;
- fresh/warm SQLite integrity and profile reuse checks;
- explicit rejection of retained `protocol does not exist` evidence;
- bounded deterministic TERM/INT/KILL native cleanup, including instance/identity/health probe failure paths;
- additive evidence checks while preserving established `DE.PULSE-G13-G14-NATIVE-2` schema compatibility;
- existing macOS WebKit lane pinned to Go 1.26.6 and required to consume the full native lifecycle evidence contract;
- frozen v18.9.1 G0-G3/scope files: `v18_9_1_g0_g3_contract.json` and `v18_9_1_scope.json`.

### Qualification truth

Last executable Fast: #557 (`32541805955`) on earlier head `89c95e4ae156cd77d66cfbd96d3911375fdee940` — PASS including gofmt, go vet and full Go suite.

Later source/test/harness hardening moved the branch beyond that SHA. Current-head Fast/Qualified attempts subsequently failed **before runner execution** with `steps=null`, including Fast #559/#560/#565/#566 and controlled required-lane retries. This is classified `INFRA_FAIL`, not product-test failure. Therefore the current branch is **not claimed Fast- or Qualified-passed**.

Do not burn repeated Actions retries while zero-step failures persist. The connected repository API cannot prove the account-level reason for hosted-runner refusal.

### Mandatory v18.9.1 resume sequence once runner execution is available

1. Re-fetch live `main`, development branch head, PR #69 and #64 comments.
2. Confirm no concurrent source changes invalidate the frozen v18.9.1 scope.
3. Run exact-head Fast and require PASS.
4. Re-arm/run full exact-head Qualified and require backend, renderer, Chrome and WebKit lanes PASS.
5. In WebKit/macOS proof require actual packaged fresh native startup + 3-second dwell, warm relaunch on the same profile, SQLite/profile reuse, no `protocol does not exist`, and deterministic cleanup evidence PASS.
6. Only after exact-head evidence is green may PR #69 leave Draft and proceed to normal review/readiness assessment.
7. Merge/release remains prohibited without explicit authorization; canonical G11-G16 is the only promotion path.

`ADAPT-RUNTIME-CRASH-001` remains **OPEN / NOT CLOSED**.

## Mandatory post-v18.9.1 checkpoint — #70

#70 / `ADAPT-CI-CONVERGENCE-001` is approved and is now part of the Adaptive Roadmap, Build Plan, Build Process and Delivery Process.

It is a **process/release-engineering workstream, not a product feature/version**. It executes immediately after truthful v18.9.1 closure and before the next product implementation so DE.PULSE does not continue accumulating version-specific CI/release machinery.

Canonical detailed contract:
`adaptive-governance/CURRENT_ADAPTIVE_CI_CONVERGENCE.md`

Required direction:
- keep exactly three routine workflows: Fast, Qualified, Release;
- Planner v3 dependency-aware job/evidence selection with fail-closed full fallback;
- manual/reusable runs bind trustworthy full-delta base/merge-base;
- targeted native rehearsal inside Qualified, with WebKit/browser and native-lifecycle ownership separated while canonical packaging scripts are reused;
- one canonical version-neutral G12 executor plus declarative release/capability manifests;
- incremental migration of active version-named tests to capability-oriented paths with equivalence proof;
- early G11 version/predecessor/tag/publication-feasibility checks before expensive certification/native work;
- Stable asset digest immutability and repository-wide release serialization;
- verify/harden actual repository rulesets/branch protection;
- separate public product version, work-slice ID, candidate SHA/fingerprint, build ID, evidence schema and monotonic platform build number;
- make a prospective SemVer-vs-explicit-custom versioning decision without rewriting shipped history;
- canonical toolchain identity plus final runnable-artifact provenance/attestation/SBOM at the supply-chain milestone;
- converge handoff and CURRENT Adaptive overlays on one actual current-state truth;
- G16 reports runner minutes/reruns avoided and proves no quality evidence was removed.

## Version/release planning rule after #70

Do **not** make every independently governed work slice or CI build a public version.

Preferred model:

`work slice -> commits -> Fast -> deliberate Qualified checkpoint -> exact candidate SHA/fingerprint -> optional coherent release grouping -> one public product version -> one canonical G11-G16/native publication`

Current post-v18.9.1 `v18.9.2` ... `v18.9.12` labels are planning reservations pending #70. Preserve their dependency/work-slice intent; #70 may regroup future public releases prospectively. Shipped versions/tags/evidence remain immutable.

## Remaining v18 product dependency order after #70

1. TradeInsight Settings/API-key UX — currently reserved `v18.9.2`
2. Coverage-aware Smart Provider Router — reserved `v18.9.3`
3. Canonical company/instrument identity — reserved `v18.9.4`
4. Market Data Modes/capability diagnostics — reserved `v18.9.5`
5. Provider observability/Adaptive telemetry — reserved `v18.9.6`
6. TradeInsight Form 4 SHADOW enrichment — reserved `v18.9.7`
7. TradeInsight symbol/company search — reserved `v18.9.8`
8. TradeInsight movers/ranking SHADOW evidence — reserved `v18.9.9`
9. Remaining useful capability admission — reserved `v18.9.10`
10. Session-Aware Data Readiness Maintenance — reserved `v18.9.11`
11. Professional Closure — reserved `v18.9.12`

Dependency order remains authoritative even if future public release grouping changes.

## Permanent release philosophy

Small dependency-ordered work slices, one primary responsibility each, G0-G16 only, canonical owners reused, observability before broad capability admission, point-in-time evidence before adaptive learning, model/prompt governance before broad adaptive influence, and durable issue/handoff truth.

Release engineering follows the same rule: reuse/consolidate permanent workflow and test owners, remove duplicate/version-specific machinery only after equivalence proof, and optimize runner cost through dependency awareness/evidence reuse rather than reduced assurance.

## Permanent Cross-Platform Lockstep Rule

DE.PULSE is **one product across macOS, Windows and Web**.

For every shared capability G1 freezes Mac/Windows/Web as REQUIRED or justified N/A. One canonical domain/API/state contract drives all REQUIRED clients. Platform adapters may differ only for OS/browser mechanics; business logic, intelligence, account/state semantics, authorization, product entitlement, provider-right decisions, freshness/provenance and explanation meaning may not fork.

A single platform may be used for diagnosis/technical validation, but there is **no product pilot** and no normal `Mac -> Windows -> Web catch-up` sequence. Shared capability GA/Delivered state and G10/G15 require all REQUIRED clients. No next shared domain begins while material parity debt remains. Temporary exceptions require an external blocker, waiver/expiry, no misleading GA claim and a named recovery release.

Platform-specific corrective work is valid where the defect itself is platform-specific; #64/v18.9.1 is the current example.

## Provider/persistence contract

`consumer requirement -> memory cache -> persisted canonical state -> validate freshness/coverage/schema/provenance/rights -> residual gap -> Smart Provider Router v2 -> targeted acquisition -> canonical merge/reconciliation -> persist -> serve`

Direct SEC/EDGAR remains authoritative for Form 4. TradeInsight validation-required capabilities start SHADOW-first. No provider-specific router/cache/scanner/scheduler/Market Mode/SEC/symbol/persistence subsystem.

## Hosted / zero-key architecture

`Mac/Windows native SQLite edge + Web -> DE.PULSE hosted APIs -> tenant/RBAC/product-entitlement/provider-right/privacy checks -> canonical router/freshness/cache/state -> PostgreSQL shared authority`

Normal commercial users authenticate only to DE.PULSE; platform provider credentials remain server-side in managed secrets/KMS.

Hosted serving keeps five dimensions separate:
1. tenant/account identity;
2. RBAC/capabilities;
3. DE.PULSE product entitlement/plan/quota;
4. upstream provider legal/data rights;
5. privacy/data-governance policy.

## Planned v19 train — current authority

Entry requires v18 native/data-plane closure plus #70 PASS.

### v19.0.x — Hosted Foundations
Provider legal-rights registry; tenant/identity/device/session control plane; product entitlement/metering; account data governance/privacy; hosted environment/IaC/service trust; PostgreSQL HA/PITR; managed secrets/KMS; software supply-chain/artifact/dependency assurance; provider SLO/cost/coverage scorecards; reconciliation/revision/point-in-time quality.

### v19.1.x — Hosted Data Plane + Cross-Platform Account/State
Hosted Provider Gateway; Unified Serving Policy + Live Fan-Out; Sync Protocol Foundation; **Mac + Windows + Web** Account/Session, Preferences, Watchlists/Master Symbols and Desks/Workspaces.

### v19.2.x — Cross-Platform Shared Product + Assurance
**Mac + Windows + Web** Research/Durable State, Discovery/Opportunity Radar, Market State/Modes/Readiness/Explanations, Settings/Account/Device Controls and RBAC/Product-Entitlement UX; tenant-aware metering; mixed-client security/abuse/capacity; #66 Cross-Platform Assurance Closure.

### v19.3.x — Point-in-Time Evidence
Institutional/13F -> Two-sided Long/Short evidence -> AODR lineage.

### v19.4.x — Reliability / Economics / Readiness
ADR-GDI reliability/capacity -> specialized/paid-provider gap evaluation -> v20 research-readiness audit.

### v19.5.0 — Major Closure
No feature scope. #66 must PASS with zero material Mac/Windows/Web parity debt plus rights/privacy/security/IaC/supply-chain/API/recovery/SLO/capacity proof.

## Planned v20 train

`adaptive control/model governance -> ASBI -> adaptive Institutional/TDTI -> AODR -> adaptive operations -> Professional Closure`

Every shared adaptive user-facing capability follows Mac + Windows + Web lockstep. No Execution remains permanent.

## G0-G16 lockstep/release enforcement

G1 scope/platform/version disposition; G2 canonical owner/adapters; G3 one contract + equivalence/dependency plan; G4 all REQUIRED implementations; G5/G6 affected evidence with fail-closed uncertainty; G7 equivalent security/data outcomes; G8 mixed-client/runtime/CI efficiency; G9 function/meaning equivalence; G10 parity/evidence sufficiency blocks freeze; G11 exact candidate + publication feasibility; G12 canonical system certification; G13/G14 actual artifacts/deployments; G15 no GA until all REQUIRED clients pass and no differing Stable-byte overwrite; G16 parity/CI-efficiency/provenance/handoff audit.

## Permanent boundaries

U.S. Equities Processing; GLD/SLV/USO actionable exceptions; No Execution; Smart Provider Router v2 sole routing authority; canonical freshness/recovery sole freshness owner; existing multi-feed allocator sole subscription owner; BroadSnapshotBroker canonical reuse owner; canonical persistence/cache/state/session/calendar/identity owners; direct SEC/EDGAR authoritative; canonical `SUPER_OWNER/OWNER/ADMIN/USER/DEMO`; no parallel client/provider truth.

## Exactly one next action

Restore/confirm GitHub Actions hosted-runner execution, then re-fetch the exact current `v18.9.1-development` head and run exact-head Fast followed by full Qualified. Do not move PR #69 out of Draft, merge or release until that evidence passes. After truthful v18.9.1 closure, execute #70 before any next product implementation.