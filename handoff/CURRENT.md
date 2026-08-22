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
**Governance alignment PR:** #67 merged  
**Master corrective program:** #65 / `ADAPT-PROVIDER-INTELLIGENCE-010`  
**Hosted architecture program:** #66 / `ADAPT-HOSTED-SYNC-001`

## Immediate execution rule

Do not start v18.9.2 or v19 implementation until #64 / v18.9.1 is truthfully closed. First refetch live `main`, `v18.9.1-development`, PR #69, issue #64 and current comments because concurrent sessions/processes may advance the branch.

Do **not** restart G0. G0 deterministic reproduction is complete and the canonical macOS native-window root cause is isolated. Continue from the exact current branch/PR state and the current qualification blocker described below.

## v18.9.1 current execution state

The canonical v18.9.0 package audit proved the Stable artifact uses `DePulseLauncher` -> `DePulse-arm64`; the original user crash path `/Applications/De-Pulse.app/Contents/MacOS/De-Pulse` is a distinct/non-matching executable signature and must not be conflated with the canonical reproduction.

Qualified #157 (`32541591746`) executed the real packaged macOS JXA/Cocoa/WKWebView path on macOS 15.7.7 Apple Silicon with pinned Go 1.26.6. Packaged backend/readiness/root and SQLite migrations `[1,2,3,4]` passed, then the native JXA child failed with `protocol does not exist (-2700)`. The only explicit protocol lookup was formal `protocols:['NSApplicationDelegate']` in the JXA registered subclass.

The v18.9.1 corrective candidate therefore removes only that unnecessary formal protocol declaration while preserving the NSObject delegate selectors and `app.delegate=delegate` behavior. The existing macOS packaging and WebKit qualification owners are hardened rather than replaced.

Implementation head immediately before this handoff-only update: `083b69c6772bb7a0fa14a7cdea70f4bd695a10bb`. Re-fetch the branch before relying on that SHA.

Current implementation includes:
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

Later source/test/harness hardening moved the branch beyond that SHA. Current-head Fast/Qualified attempts subsequently fail **before runner execution** with `steps=null`, including Fast #559/#560/#565/#566 and controlled required-lane retries. This is classified `INFRA_FAIL`, not product-test failure. Therefore the current branch head is **not claimed Fast- or Qualified-passed**.

Do not burn repeated Actions retries while zero-step failures persist. The connected repository API cannot prove the account-level reason for hosted-runner refusal. Public GitHub status was operational when checked; private-repository usage/budget/payment/quota constraints remain a possible but unproven account-specific class.

### Mandatory resume sequence once runner execution is available

1. Re-fetch live `main`, development branch head, PR #69 and #64 comments.
2. Confirm no concurrent source changes invalidate the frozen v18.9.1 scope.
3. Run exact-head Fast and require PASS.
4. Re-arm/run full exact-head Qualified and require backend, renderer, Chrome and WebKit lanes PASS.
5. In WebKit/macOS proof require actual packaged fresh native startup + 3-second dwell, warm relaunch on the same profile, SQLite/profile reuse, no `protocol does not exist`, and deterministic cleanup evidence PASS.
6. Only after exact-head evidence is green may PR #69 leave Draft and proceed to normal review/readiness assessment.
7. Merge/release remains prohibited without explicit authorization; canonical G11-G16 is the only promotion path.

`ADAPT-RUNTIME-CRASH-001` remains **OPEN / NOT CLOSED**.

## Permanent release philosophy

Small dependency-ordered patches, one primary responsibility each, G0-G16 only, canonical owners reused, observability before broad capability admission, point-in-time evidence before adaptive learning, model/prompt governance before broad adaptive influence, and durable issue/handoff truth.

## Permanent Cross-Platform Lockstep Rule

DE.PULSE is **one product across macOS, Windows and Web**.

For every shared capability G1 freezes Mac/Windows/Web as REQUIRED or justified N/A. One canonical domain/API/state contract drives all REQUIRED clients. Platform adapters may differ only for OS/browser mechanics; business logic, intelligence, account/state semantics, authorization, product entitlement, provider-right decisions, freshness/provenance and explanation meaning may not fork.

A single platform may be used for diagnosis/technical validation, but there is **no product pilot** and no normal `Mac -> Windows -> Web catch-up` sequence. Shared capability GA/Delivered state and G10/G15 require all REQUIRED clients. No next shared domain begins while material parity debt remains. Temporary exceptions require an external blocker, waiver/expiry, no misleading GA claim and a named recovery release.

Platform-specific corrective work is valid where the defect itself is platform-specific; #64/v18.9.1 is the current example.

## Ordered v18.9.x train

1. `v18.9.1` runtime crash corrective
2. `v18.9.2` TradeInsight Settings/API-key UX
3. `v18.9.3` coverage-aware Smart Provider Router
4. `v18.9.4` canonical company/instrument identity
5. `v18.9.5` Market Data Modes/capability diagnostics
6. `v18.9.6` provider observability/Adaptive telemetry
7. `v18.9.7` TradeInsight Form 4 SHADOW enrichment
8. `v18.9.8` TradeInsight symbol/company search
9. `v18.9.9` TradeInsight movers/ranking SHADOW evidence
10. `v18.9.10` remaining useful capability admission
11. `v18.9.11` Session-Aware Data Readiness Maintenance
12. `v18.9.12` Professional Closure

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

### v19.0.x — Hosted Foundations
- `v19.0.0` provider legal-rights registry
- `v19.0.1` tenant/identity/device/session control plane
- `v19.0.2` product entitlement/metering policy
- `v19.0.3` account data governance/privacy lifecycle
- `v19.0.4` hosted environment/IaC/service trust
- `v19.0.5` PostgreSQL tenancy/schema/pool/HA-PITR
- `v19.0.6` managed secrets/KMS
- `v19.0.7` software supply-chain/artifact/dependency assurance
- `v19.0.8` provider SLO/cost/coverage scorecards
- `v19.0.9` reconciliation/revision/point-in-time quality

### v19.1.x — Hosted Data Plane + Cross-Platform Account/State
- `v19.1.0` Hosted Provider Gateway
- `v19.1.1` Unified Serving Policy + Live Fan-Out
- `v19.1.2` Sync Protocol Foundation
- `v19.1.3` **Mac + Windows + Web Account/Session Client Foundation**
- `v19.1.4` **Mac + Windows + Web Preferences**
- `v19.1.5` **Mac + Windows + Web Watchlists/Master Symbols**
- `v19.1.6` **Mac + Windows + Web Desks/Workspaces**

### v19.2.x — Cross-Platform Shared Product + Assurance
- `v19.2.0` **Mac + Windows + Web Research/Durable State**
- `v19.2.1` **Mac + Windows + Web Discovery/Opportunity Radar**
- `v19.2.2` **Mac + Windows + Web Market State/Market Modes/Readiness/Explanations**
- `v19.2.3` **Mac + Windows + Web Settings/Account/Device Controls**
- `v19.2.4` **Mac + Windows + Web RBAC/Product-Entitlement UX**
- `v19.2.5` tenant-aware metering/cost/usage observability
- `v19.2.6` mixed-client security/abuse/noisy-neighbor/capacity hardening
- `v19.2.7` **#66 Cross-Platform Assurance Closure**

### v19.3.x — Point-in-Time Evidence
- `v19.3.0` institutional/13F infrastructure
- `v19.3.1` two-sided Long/Short evidence substrate
- `v19.3.2` AODR candidate/ranking/outcome lineage

### v19.4.x — Reliability / Economics / Readiness
- `v19.4.0` ADR-GDI professional reliability/capacity/runbooks
- `v19.4.1` specialized/paid-provider gap evaluation
- `v19.4.2` v20 research-readiness audit

### v19.5.0 — v19 Major Closure
No feature scope. #66 must PASS with zero material Mac/Windows/Web parity debt plus rights/privacy/security/IaC/supply-chain/API/recovery/SLO/capacity proof.

## Planned v20 train

`v20.0.x` adaptive control/model governance -> `v20.1.x` ASBI -> `v20.2.x` adaptive Institutional/TDTI -> `v20.3.x` AODR -> `v20.4.0` adaptive operations -> `v20.5.0` Professional Closure.

Every shared adaptive user-facing capability follows Mac + Windows + Web lockstep. No Execution remains permanent.

## G0-G16 lockstep enforcement

G1 platform matrix; G2 canonical owner/adapters; G3 one contract + equivalence tests; G4 all REQUIRED implementations; G6 cross-platform integration; G7 equivalent security/data outcomes; G8 mixed-client load; G9 function/meaning equivalence; G10 parity debt blocks freeze; G13/G14 actual artifacts/deployments; G15 no GA until all REQUIRED clients pass; G16 parity-drift audit.

## Permanent boundaries

U.S. Equities Processing; GLD/SLV/USO actionable exceptions; No Execution; Smart Provider Router v2 sole routing authority; canonical freshness/recovery sole freshness owner; existing multi-feed allocator sole subscription owner; BroadSnapshotBroker canonical reuse owner; canonical persistence/cache/state/session/calendar/identity owners; direct SEC/EDGAR authoritative; canonical `SUPER_OWNER/OWNER/ADMIN/USER/DEMO`; no parallel client/provider truth.

## Exactly one next action

Restore/confirm GitHub Actions hosted-runner execution, then re-fetch the exact current `v18.9.1-development` head and run exact-head Fast followed by full Qualified. Do not move PR #69 out of Draft, merge, release, or start v18.9.2 until that evidence passes.
