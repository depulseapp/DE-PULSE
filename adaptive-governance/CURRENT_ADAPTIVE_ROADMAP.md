# DE.PULSE — Current Adaptive Roadmap

**Certified Stable:** `v18.9.0-stable`  
**Certified candidate:** `9ea81cddae4875ae15d3719ca028519a36c597b6`  
**Certified fingerprint:** `a8719090c341c874dbd1279cc31ad98e84075d5701c46a800bf951340780ecb9`  
**Build ID:** `v18.9.0-stable-20260821`  
**Active product branch/PR:** `v18.9.1-development` / Draft PR #69  
**Active corrective:** #64 / `ADAPT-RUNTIME-CRASH-001` -> `v18.9.1`  
**Approved CI/versioning/repository convergence:** #70 / `ADAPT-CI-CONVERGENCE-001` — mandatory immediately after truthful v18.9.1 closure  
**Master corrective program:** #65 / `ADAPT-PROVIDER-INTELLIGENCE-010`  
**Hosted architecture program:** #66 / `ADAPT-HOSTED-SYNC-001`

## 1. North star

`v18.9 trustworthy native runtime/data plane + scalable release engineering + clean enforceable repository ownership -> v19 one hosted DE.PULSE product delivered Mac/Windows/Web in lockstep -> v20 governed adaptive intelligence under the same lockstep contract`

Permanent boundaries: U.S. Equities Processing; GLD/SLV/USO actionable exceptions; No Execution; G0-G16 only; Smart Provider Router v2 sole routing owner; canonical freshness/cache/persistence/subscription/session/SEC/identity owners; direct SEC/EDGAR authority; deterministic Day/Swing/Long protection; SHADOW -> VALIDATED -> APPROVED -> PRODUCTION.

## 2. Cross-Platform Lockstep Roadmap Rule

DE.PULSE does not maintain separate feature roadmaps for Mac, Windows and Web.

For each shared capability:
- one canonical domain/API/state contract;
- G1 freezes Mac/Windows/Web as REQUIRED or justified N/A;
- all REQUIRED clients ship within the same capability release responsibility;
- platform adapters may differ only for OS/browser mechanics;
- business logic, intelligence, state, authorization, product entitlement, provider-right decisions, freshness/provenance and explanation meaning do not fork;
- one-platform technical validation is internal only, not a product pilot;
- no GA and no next shared domain while material required-platform parity debt remains;
- temporary exceptions require an external blocker, explicit waiver/expiry and named recovery release.

Platform-specific corrective work is allowed when the actual responsibility is platform-specific, such as `v18.9.1`.

## 3. Immediate ordered roadmap

1. **Finish `v18.9.1` / #64 truthfully.** G0/root cause is complete. Restore hosted-runner execution, then exact-head Fast + full Qualified + required macOS fresh/warm native evidence. PR #69 stays Draft until that passes.
2. **Execute #70 / `ADAPT-CI-CONVERGENCE-001`.** This is a mandatory executable process/repository-hardening checkpoint before the next product implementation. Documentation alone cannot satisfy it. It includes actual CI/release changes and repository-root cleanup/enforcement.
3. **Re-baseline future v18.9.x reservations** under the approved separation of public product version, work-slice ID, source SHA/fingerprint, build ID and evidence schema. Shipped versions/tags/evidence remain immutable.
4. **Resume product capability work** only after #70 has established and proven the scalable CI/versioning/repository contract.

The detailed #70 contract is `adaptive-governance/CURRENT_ADAPTIVE_CI_CONVERGENCE.md` and GitHub issue #70. It is part of this roadmap, not optional backlog.

## 4. v18.9.x — Trustworthy Native Runtime/Data Plane

The following product identities remain planning reservations until #70 makes the prospective versioning decision. They may be regrouped into a smaller number of public releases without losing their independent work-slice/governance identity. Shipped versions are never renumbered.

1. `v18.9.1` Runtime crash corrective — active
2. `v18.9.2` TradeInsight Settings/API-key UX — reservation
3. `v18.9.3` Coverage-aware Smart Provider Router — reservation
4. `v18.9.4` Canonical company/instrument identity — reservation
5. `v18.9.5` Market Data Modes/capability diagnostics — reservation
6. `v18.9.6` Provider observability/Adaptive telemetry — reservation
7. `v18.9.7` TradeInsight Form 4 SHADOW enrichment — reservation
8. `v18.9.8` TradeInsight symbol/company search — reservation
9. `v18.9.9` TradeInsight movers/ranking SHADOW evidence — reservation
10. `v18.9.10` Remaining useful TradeInsight capability admission — reservation
11. `v18.9.11` Session-Aware Data Readiness Maintenance — reservation
12. `v18.9.12` Professional Closure — reservation/no feature scope

Dependency order remains authoritative even if future public release grouping changes:

`runtime reliability -> scalable CI/versioning/repository ownership -> settings/secret UX -> coverage-aware routing -> canonical identity -> diagnostics -> observability -> SHADOW capability validation -> operational readiness -> closure`

## 5. CI / versioning / repository roadmap checkpoint — #70

The permanent workflow topology stays exactly:

`Fast -> Qualified -> Release`

No fourth routine workflow and no retry/certification/promotion branch family.

Roadmap outcomes required from #70:
- Planner v3 selects the smallest trustworthy evidence graph with explicit reasons and fail-closed full fallback;
- manual/reusable runs bind a trustworthy full delta, never silently `HEAD^` only;
- targeted native rehearsal lives inside Qualified and remains separate from WebKit browser ownership while reusing canonical native scripts;
- future G12 uses one version-neutral executor plus declarative release/capability manifests;
- active version-named tests migrate incrementally to capability ownership paths with equivalence proof;
- G11 fails impossible publication before expensive certification/native work;
- Stable assets are digest-immutable and release execution is globally serialized;
- actual repository branch/ruleset enforcement is verified;
- public product version is separated from work-slice/build/source/evidence identity;
- future version/tag semantics are explicitly SemVer-compliant or explicitly custom; historical releases remain untouched;
- native build numbers are monotonic/collision-safe;
- toolchain identity and final runnable-artifact provenance/attestation become durable release evidence;
- G16 reports runner minutes/reruns avoided and proves no required evidence was removed;
- repository-root `.gitignore` safely excludes transient local/build/editor evidence without hiding tracked continuity state;
- source-health becomes recursive/package-aware before production Go files are moved;
- the existing legacy inventory becomes the canonical root-layout inventory/allowlist owner;
- reusable root CI/dev/release tooling migrates to stable `tools/` ownership with consumers changed atomically;
- stale `certification_plan.json`, `certification_runner.py`, `ci_pipeline.py`, and `ci_pipeline_plan.json` cease to act as a competing current CI/release model through consolidation/retirement with equivalence proof;
- historical/version-scoped non-Go root material moves to governed release/history ownership only after evidence mapping;
- active version-named Go tests migrate safely without losing package-local/full/race/randomized/focused coverage;
- production package decomposition occurs only after package-aware source-health guards exist and implementation/tests move together;
- policies, registries and retained assets move to stable canonical owners with all consumers preserved;
- a CI-enforced small root allowlist prevents version-stacked clutter from recurring;
- G16 records before/after root inventory and files moved/renamed/consolidated/removed.

**Roadmap completion rule:** #70 is not complete if only Markdown/issues were updated. Actual repository/code/workflow/tooling changes and applicable executable evidence are mandatory.

## 6. v19 — Professional Hosted Product

**Entry:** v18 native/data-plane closure plus #70 CI/versioning/repository-hygiene convergence PASS.

### v19.0.x — Governance / Control Plane / Data Foundation
- `v19.0.0` provider legal-rights registry
- `v19.0.1` tenant/identity/device/session control plane
- `v19.0.2` DE.PULSE product entitlement/metering policy
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

There is no macOS product pilot and no later Windows/Web catch-up phase.

### v19.2.x — Cross-Platform Shared Product + Assurance
- `v19.2.0` **Mac + Windows + Web Research/Durable State**
- `v19.2.1` **Mac + Windows + Web Discovery/Opportunity Radar**
- `v19.2.2` **Mac + Windows + Web Market State/Market Modes/Readiness/Explanations**
- `v19.2.3` **Mac + Windows + Web Settings/Account/Device Controls**
- `v19.2.4` **Mac + Windows + Web RBAC/Product-Entitlement UX**
- `v19.2.5` tenant-aware metering/cost/usage observability
- `v19.2.6` mixed-client multi-user security/abuse/capacity hardening
- `v19.2.7` **#66 Cross-Platform Assurance Closure**

### v19.3.x — Point-in-Time Evidence
- `v19.3.0` Institutional/13F infrastructure
- `v19.3.1` Two-sided Long/Short evidence substrate
- `v19.3.2` AODR candidate/ranking/outcome lineage

### v19.4.x — Reliability / Economics / v20 Readiness
- `v19.4.0` ADR-GDI professional reliability/capacity/runbooks
- `v19.4.1` specialized/paid-provider gap evaluation
- `v19.4.2` v20 research-readiness audit

### v19.5.0 — Major Closure
No feature scope. Require #66 PASS, zero material Mac/Windows/Web parity debt, rights/identity/RBAC/product-entitlement/privacy separation, IaC/environment and supply-chain assurance, data lifecycle, API compatibility, recovery/rollback, SLO/capacity and actual supported artifact/deployment proof.

## 7. v20 — Governed Adaptive Intelligence

- `v20.0.0` Adaptive Research Control Plane + Immutable Experiment Ledger
- `v20.0.1` Model/Prompt Governance + Champion/Challenger
- `v20.0.2` Historical Analogues/Regime Outcomes
- `v20.0.3` Calibration/FP-FN/Miss/Contradiction/Drift
- `v20.1.x` ASBI
- `v20.2.x` adaptive Institutional/13F + TDTI
- `v20.3.x` AODR
- `v20.4.0` ADR-GDI adaptive operations
- `v20.5.0` Professional Closure

Every shared adaptive user-facing capability follows Mac + Windows + Web lockstep.

## 8. G0-G16 lockstep enforcement

G1 platform matrix -> G2 canonical owner/adapters -> G3 one contract + equivalence tests -> G4 all REQUIRED client implementations/process changes -> G6 cross-platform/integration/path-consumer proof -> G7 equivalent security/data outcomes -> G8 mixed-client/cost/capacity -> G9 UX/function/meaning equivalence -> G10 parity debt/competing owner/lost coverage blocks freeze -> G11 exact candidate + publication feasibility -> G12 system certification -> G13/G14 actual artifacts/deployments -> G15 no GA until all REQUIRED clients pass -> G16 parity/CI-efficiency/root-inventory/handoff audit.

## Exactly one next action

Restore/confirm GitHub Actions hosted-runner execution, then finish exact-head #64 / `v18.9.1` Fast + full Qualified + native proof. After truthful v18.9.1 closure, #70—including actual repository-root implementation—is the mandatory next build/process workstream before any v18.9.2+ product implementation.
