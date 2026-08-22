# DE.PULSE — Current Adaptive Roadmap

**Certified Stable:** `v18.9.0-stable`  
**Certified candidate:** `9ea81cddae4875ae15d3719ca028519a36c597b6`  
**Certified fingerprint:** `a8719090c341c874dbd1279cc31ad98e84075d5701c46a800bf951340780ecb9`  
**Build ID:** `v18.9.0-stable-20260821`  
**Master corrective program:** #65 / `ADAPT-PROVIDER-INTELLIGENCE-010`  
**Hosted architecture program:** #66 / `ADAPT-HOSTED-SYNC-001`  
**Immediate next product patch:** #64 / `ADAPT-RUNTIME-CRASH-001` -> `v18.9.1`.

## 1. Roadmap north star

DE.PULSE evolves through:

`v18.9.x trustworthy runtime/data plane -> v19 professional hosted multi-tenant product with Mac/Windows/Web lockstep -> v20 governed adaptive intelligence with the same lockstep product contract`

Permanent boundaries:
- U.S. Equities Processing only; GLD/SLV/USO actionable exceptions;
- No Execution;
- G0-G16 only;
- one canonical owner per responsibility;
- Smart Provider Router v2 remains sole routing authority;
- canonical freshness, cache, persistence, subscription, SEC, identity, session/calendar and state owners are reused;
- provider legal rights, DE.PULSE product entitlement, RBAC and privacy/data-governance remain distinct controls;
- adaptive production influence follows `SHADOW -> VALIDATED -> APPROVED -> PRODUCTION`;
- no silent deterministic Day/Swing/Long change or silent self-promotion.

## 2. Permanent Cross-Platform Lockstep Roadmap Rule

DE.PULSE is **one product**, not separate Mac/Windows/Web roadmaps.

For every shared capability:
- one canonical domain/API/state contract is planned first;
- G1 freezes whether macOS, Windows and Web are REQUIRED or justified N/A;
- all REQUIRED clients are part of the same capability release responsibility;
- platform-specific adapters may differ only where the OS/browser requires it;
- business logic, intelligence, account/state semantics, authorization, product entitlement, provider rights, freshness/provenance and explanation meaning may not fork by client;
- a shared capability is not Delivered/GA until all REQUIRED clients pass;
- there is no normal sequence of `Mac first -> Windows later -> Web later`;
- a single platform may be used for internal technical validation, but that is not a product pilot or roadmap milestone;
- no next shared domain begins while the current domain has material required-platform parity debt;
- temporary exceptions require a proven external blocker, explicit expiry, no misleading GA status and a named recovery release.

Platform-specific corrective work remains valid where the responsibility itself is platform-specific, such as `v18.9.1` for the macOS crash.

## 3. Version alignment

- **Major** = strategic maturity generation.
- **Minor band** = coherent dependency phase.
- **Patch** = one primary independently certifiable responsibility; required platform adapters for one shared capability belong to that responsibility.
- Future numbers are reservations until G1 and may shift/split before implementation.
- Shipped versions are immutable.
- Corrective/security work may preempt the planned train.

## 4. v18.9.x — Stabilize -> instrument -> validate -> operationalize -> close

1. `v18.9.1` Runtime crash corrective ONLY — #64; macOS-specific escaped defect.
2. `v18.9.2` TradeInsight Settings/API-key UX ONLY.
3. `v18.9.3` Coverage-aware Smart Provider Router core ONLY.
4. `v18.9.4` Canonical company/instrument identity ONLY.
5. `v18.9.5` Market Data Modes + capability diagnostics ONLY.
6. `v18.9.6` Provider observability / Adaptive telemetry ONLY.
7. `v18.9.7` TradeInsight SEC Form 4 SHADOW enrichment ONLY.
8. `v18.9.8` TradeInsight ticker/company search ONLY.
9. `v18.9.9` TradeInsight movers/ranking SHADOW evidence ONLY.
10. `v18.9.10` Remaining useful TradeInsight capability admission ONLY.
11. `v18.9.11` Session-Aware Data Readiness Maintenance ONLY.
12. `v18.9.12` Professional Closure ONLY.

The v18.9 line remains native-first. Shared native behavior must stay semantically equivalent where applicable. Web becomes a REQUIRED shared client when the v19 cross-platform client foundation is delivered.

## 5. v19 — Professional Data Infrastructure + Hosted Account Platform

**Entry:** `v18.9.12` PASS.

Target architecture:

```text
Mac / Windows native clients        Hosted Web
       |                                |
       | SQLite edge/offline            | hosted client state
       +---------------+----------------+
                       |
              DE.PULSE hosted APIs
                       |
       tenant/RBAC/product entitlement/
       provider rights/privacy policy
                       |
        Smart Provider Router v2 +
        canonical freshness/cache/state
                       |
                PostgreSQL authority
```

Normal commercial users are zero-key. Provider credentials remain server-side.

### v19.0.x — Governance / Control Plane / Data Foundation

- `v19.0.0` Provider Capability / Legal Rights Registry.
- `v19.0.1` Hosted Tenant / Identity / Device / Session Control Plane.
- `v19.0.2` DE.PULSE Product Entitlement / Metering Policy.
- `v19.0.3` Account Data Governance / Privacy Lifecycle.
- `v19.0.4` Hosted Environment / IaC / Service Trust Foundation.
- `v19.0.5` PostgreSQL Tenancy / Schema / Pool / HA-PITR Foundation.
- `v19.0.6` Managed Secrets / KMS Lifecycle.
- `v19.0.7` Software Supply-Chain / Artifact & Dependency Assurance.
- `v19.0.8` Provider Quality / Cost / Coverage / SLO Scorecards.
- `v19.0.9` Data Reconciliation / Revision / Point-in-Time Quality.

**Exit:** hosted activation waits until these foundations are executable and evidenced.

### v19.1.x — Hosted Data Plane + Cross-Platform Account/Sync Lockstep

- `v19.1.0` Authenticated Hosted Provider Gateway.
- `v19.1.1` Unified Serving Policy + Live Fan-Out Isolation.
- `v19.1.2` Sync Protocol Foundation.
- `v19.1.3` **Cross-Platform Account / Session Client Foundation — Mac + Windows + Web together.**
- `v19.1.4` **Cross-Platform Preferences + Watchlists — Mac + Windows + Web together.**
- `v19.1.5` **Cross-Platform Desks / Workspaces — Mac + Windows + Web together.**

There is no macOS product pilot. A platform can be used internally for technical validation, but phase exit requires every REQUIRED client.

### v19.2.x — Cross-Platform Shared Product + Assurance

- `v19.2.0` **Cross-Platform Research / Durable State — Mac + Windows + Web.**
- `v19.2.1` **Cross-Platform Market Intelligence / Discovery / Market Modes — Mac + Windows + Web.**
- `v19.2.2` **Cross-Platform Settings / RBAC / Product-Entitlement UX — Mac + Windows + Web.**
- `v19.2.3` Tenant-Aware Metering / Cost / Usage Observability.
- `v19.2.4` Mixed-Client Multi-User Security / Abuse / Capacity Hardening.
- `v19.2.5` **#66 Cross-Platform Hosted Sync / Gateway Assurance Closure.**

#66 cannot close with material parity debt. Equivalent Mac/Windows/Web requests must observe the same canonical account, authorization, state, intelligence, freshness/provenance and degraded/UNKNOWN semantics.

### v19.3.x — Professional Point-in-Time Evidence Substrate

- `v19.3.0` Institutional / 13F Evidence Infrastructure.
- `v19.3.1` Two-Sided Long / Short Thesis Evidence Substrate.
- `v19.3.2` AODR Candidate / Ranking / Outcome Lineage.

Any user-facing capability surfaced from these substrates follows the lockstep contract.

### v19.4.x — Reliability / Economics / v20 Readiness

- `v19.4.0` ADR-GDI Professional Reliability / Capacity.
- `v19.4.1` Specialized / Paid Provider Gap Evaluation.
- `v19.4.2` v20 Research-Readiness Audit.

### v19.5.0 — v19 Major Closure

No new feature scope. Require #66 PASS, zero material Mac/Windows/Web parity debt for shared capabilities, tenant isolation, product-entitlement/provider-right/privacy separation, data lifecycle, environment/IaC, supply-chain provenance, API compatibility, recovery/rollback, SLO/capacity and actual supported native/Web runtime/deployment evidence.

## 6. v20 — Adaptive Intelligence & Decision Research

**Entry:** `v19.5.0` PASS.

### v20.0.x — Adaptive Research Control & Governance
- `v20.0.0` Adaptive Research Control Plane + Immutable Experiment Ledger.
- `v20.0.1` Model / Prompt Governance + Champion/Challenger.
- `v20.0.2` Historical Analogues + Regime-Conditioned Outcomes.
- `v20.0.3` Calibration / FP-FN / Miss / Contradiction / Drift.

### v20.1.x — ASBI
- `v20.1.0` Behavioral Fingerprints + State Transitions.
- `v20.1.1` Scenarios / Probability Momentum / Calibration.

### v20.2.x — Institutional + TDTI
- `v20.2.0` Adaptive Institutional / 13F Intelligence.
- `v20.2.1` TDTI Competing Long / Short / No Reliable Edge.
- `v20.2.2` TDTI Two-Sided Trade-Plan Validation. No Execution.

### v20.3.x — AODR
- `v20.3.0` Adaptive Shared Opportunity Ranking.
- `v20.3.1` Diversity / Opportunity Cost / Personalized Relevance after shared truth.

### v20.4.x — Adaptive Operations
- `v20.4.0` ADR-GDI Adaptive Optimization under SHADOW/Champion-Challenger.

### v20.5.0 — Professional Closure

No feature scope. Require calibration/utility/drift/abstention, deterministic-boundary protection, privacy/security/data rights, reproducibility/rollback, actual supported artifacts, **cross-platform lockstep for every shared user-facing adaptive capability**, zero silent self-modification and No Execution.

## 7. Industry-strength controls inside G0-G16

No G17+.

- **G1:** required-platform matrix for every capability.
- **G2:** canonical owner + client-adapter boundaries, tenant/data/privacy classification, threat model.
- **G3:** one API/domain/state contract, platform adapter contracts, compatibility/equivalence tests, migration/rollback, SLOs, IaC/supply-chain requirements.
- **G4:** all REQUIRED platform implementations for shared scope before Development Exit.
- **G7:** equivalent authorization/data/privacy/provider-right outcomes across clients.
- **G8:** mixed-client load/capacity/failure evidence where relevant.
- **G9:** Mac/Windows/Web UX/function/meaning equivalence; responsive/native interaction differences allowed.
- **G10:** material parity debt blocks freeze.
- **G13/G14:** actual required platform artifacts/deployments/runtime proof.
- **G15:** no shared-capability GA until all REQUIRED clients pass.
- **G16:** parity-drift audit and exact handoff.

## 8. Why this ordering is intentional

- observability precedes provider expansion;
- rights/security/privacy/IaC/recovery precede hosted activation;
- one gateway/serving/sync foundation precedes client capability delivery;
- **shared client capabilities are delivered in lockstep, not through platform catch-up phases**;
- multi-user assurance follows real mixed-client use;
- point-in-time evidence precedes adaptive learning;
- model governance precedes broad adaptive influence.

## Exactly one next action

Diagnose #64 from complete macOS crash evidence or deterministic reproduction and freeze narrow `v18.9.1` G1. Do not start `v18.9.2` or v19 implementation until ordering permits it.