# DE.PULSE — Canonical Adaptive Roadmap

**Status:** APPROVED / ADAPTIVE  
**Authority:** canonical product sequencing and approved strategic workstreams  
**Rule:** shipped releases are immutable truth; future reservations may adapt only with durable rationale and synchronized governance/handoff updates.

## 1. Canonical truth model

- Stable tags, release evidence, source/artifact provenance and current handoff define what actually shipped.
- Historical provisional version labels do not prove implementation.
- Corrective/security/reliability/privacy/supply-chain work may preempt future sequencing.
- Known misses are fixed or durably assigned.
- G0-G16 is the only release model.

Historical provisional `v18.3 PostgreSQL / v18.4 hosted-security` placement is superseded as future placement only. Hosted PostgreSQL/shared-account authority belongs to the dependency-correct v19 train.

## 2. Permanent product/architecture boundaries

DE.PULSE is a U.S.-equities research/intelligence/decision-support system.

Permanent constraints:
- No Execution;
- U.S. Equities Processing with GLD/SLV/USO actionable exceptions;
- Smart Provider Router v2 sole executable routing authority;
- canonical freshness, subscription, persistence/cache/state, identity, session/calendar and telemetry owners reused;
- direct SEC/EDGAR authoritative for filing truth;
- equivalent lawful evidence processed canonically once and reused/fan-out where permitted;
- hosted serving keeps tenant identity, RBAC, DE.PULSE product entitlement, provider legal/data rights and privacy/data-governance policy distinct;
- production infrastructure/configuration is versioned/reproducible and artifact/dependency provenance auditable;
- deterministic Day/Swing/Long protected unless separately governed;
- adaptive influence follows SHADOW -> VALIDATED -> APPROVED -> PRODUCTION;
- no silent self-modification or invented confidence for missing evidence.

## 3. Permanent Cross-Platform Lockstep Contract

DE.PULSE is **one product across macOS, Windows and Web**.

For every shared capability:
1. one canonical domain/API/state contract;
2. G1 freezes each client as REQUIRED or justified N/A;
3. all REQUIRED client adapters/surfaces are part of the same capability release responsibility;
4. platform mechanics may differ only where OS/browser behavior requires it;
5. business logic, intelligence, account/state semantics, authorization, product entitlement, provider-right decisions, freshness/provenance and explanation meaning may not fork;
6. one-platform technical validation is diagnostic only and is not a product pilot;
7. no Delivered/GA state and no next shared domain while material REQUIRED-platform parity debt remains;
8. temporary platform exceptions require an external blocker, explicit waiver/expiry and named recovery release.

Platform-specific corrective work is permitted where the defect/responsibility itself is platform-specific.

## 4. Zero-Miss Future-Version Conservation

Starting with v19, every future patch has one primary implementation responsibility and is bound to the machine requirement-conservation ledger for its parent program. Before G1 implementation, source overlap is classified `INHERITED`, `EXTEND_EXISTING_OWNER`, `REPLACE_CONSOLIDATE`, `NEW_RESIDUAL`, or `EXTERNAL_BLOCKED`. A roadmap label never proves missing implementation.

Each dependency band ends with a no-feature zero-gap closure. The next band may not start while an applicable requirement is unassigned, unexplained, unevidenced or silently carried forward. Canonical detail for #66 is `governance/programs/ADAPT-HOSTED-SYNC-001/V19_ZERO_MISS_PLAN.md` and `governance/programs/ADAPT-HOSTED-SYNC-001/requirement-conservation.json`.

## 5. Durable strategic workstreams

- Smart Intelligent Provider Router v2 / coverage-aware residual-gap routing.
- Shared Symbol Intelligence / multi-user demand union.
- Opportunity Radar / AODR foundation.
- TradeInsight as SHADOW/secondary intelligence through canonical owners only.
- Provider -> Market Mode adaptive integration without provider-brand ownership.
- Institutional/13F point-in-time evidence.
- Two-sided thesis/TDTI evidence.
- ADR-GDI reliability/graceful degradation.
- Hosted account/zero-key Provider Gateway/sync architecture under #66.
- Governed Adaptive Intelligence after trustworthy evidence infrastructure.

## 6. v18.9.x — Trustworthy Native Runtime/Data Plane

Historical v18.9.x labels are dispositioned by completed GitHub evidence and must not be reopened from roadmap prose. The immutable Stable remains `v18.9.1-stable`; the completed provider-intelligence program is #65/#107.

1. `v18.9.1` runtime crash corrective
2. `v18.9.2` TradeInsight Settings/API-key UX
3. `v18.9.3` coverage-aware Smart Provider Router core
4. `v18.9.4` canonical company/instrument identity
5. `v18.9.5` Market Data Modes/capability diagnostics
6. `v18.9.6` provider observability/Adaptive telemetry
7. `v18.9.7` TradeInsight SEC Form 4 SHADOW enrichment
8. `v18.9.8` TradeInsight symbol/company search
9. `v18.9.9` TradeInsight movers/ranking SHADOW evidence
10. `v18.9.10` remaining useful capability admission
11. `v18.9.11` Session-Aware Data Readiness Maintenance
12. `v18.9.12` Professional Closure

Their final dispositions are governed by live GitHub evidence; several later labels were delivered/inherited/superseded without separate public Stable publication.

## 7. v19 — Professional Hosted Product — granular rebaseline

All labels below are planned reservations until their own G0/G1 source-overlap audit. #66 is not broadly started merely because the roadmap is planned.

### v19.0.x — Governance / Control Plane / Trust / Data Foundation
- `v19.0.0` Provider Legal/Data Rights Registry + Evidence Binding
- `v19.0.1` Provider Rights Lifecycle / Downgrade Enforcement
- `v19.0.2` Hosted Tenant/Account + Canonical Role Context
- `v19.0.3` Device Registry / Lifecycle / Revocation
- `v19.0.4` Session Lifecycle / Re-auth / Revocation
- `v19.0.5` DE.PULSE Product Entitlement Policy
- `v19.0.6` Quota / Metering / Plan Transition Policy
- `v19.0.7` Account Data Classification / Minimization / Retention
- `v19.0.8` Account Export / Deletion / Residency Lifecycle
- `v19.0.9` Hosted Environment / IaC / Drift Foundation
- `v19.0.10` Service Identity / Network / TLS Trust Foundation
- `v19.0.11` PostgreSQL Tenant Schema / Migration Foundation
- `v19.0.12` PostgreSQL Capacity / HA / PITR / Restore
- `v19.0.13` Managed Secrets / KMS Storage + Resolution
- `v19.0.14` Secret Rotation / Revoke / Rollback / Audit
- `v19.0.15` Software Supply Chain / SBOM / Artifact Provenance
- `v19.0.16` Provider Quality / Cost / Coverage / SLO Scorecards
- `v19.0.17` Revision / Reconciliation / Point-in-Time Quality Primitives
- `v19.0.18` v19.0 Foundation Zero-Gap Closure

### v19.1.x — Hosted Gateway / Serving / Sync Primitives
- `v19.1.0` Authenticated Versioned Hosted Provider Gateway
- `v19.1.1` Unified Hosted Serving Authorization Policy
- `v19.1.2` Rights/Entitlement-Safe Cache + Persistence Reuse
- `v19.1.3` Hosted Live Subscription Reuse + Authorized Fan-Out
- `v19.1.4` Long-Lived Stream Revocation / Re-authorization
- `v19.1.5` API / Protocol Compatibility + Deprecation Lifecycle
- `v19.1.6` SQLite Durable Outbox + Typed Mutation Envelope
- `v19.1.7` Server Idempotency / Revision / Change Sequence
- `v19.1.8` Client Pull / Durable Apply / Checkpoint
- `v19.1.9` New-Device Bootstrap / High-Watermark
- `v19.1.10` Stale Checkpoint / Tombstone / Compaction Recovery
- `v19.1.11` Domain Conflict / Version / Delete Semantics
- `v19.1.12` Local Account Isolation / Lost-Device Behavior
- `v19.1.13` Sync Retry / Backpressure / Protected-Session Scheduling
- `v19.1.14` Gateway + Sync Tenant-Aware Observability
- `v19.1.15` v19.1 Data-Plane + Sync Zero-Gap Closure

### v19.2.x — Cross-Platform Shared Product + #66 Assurance
- `v19.2.0` Cross-Platform Account / Session / Device Client Foundation
- `v19.2.1` Cross-Platform Settings / Account / Device Controls
- `v19.2.2` Cross-Platform RBAC / Product-Entitlement UX
- `v19.2.3` Cross-Platform Portable Preferences
- `v19.2.4` Cross-Platform Watchlists / Master Symbols
- `v19.2.5` Cross-Platform Desks / Workspaces
- `v19.2.6` Cross-Platform Saved Searches / Notes / Research State
- `v19.2.7` Cross-Platform Rights-Aware Durable Research / Evidence
- `v19.2.8` Cross-Platform Discovery / Opportunity Radar
- `v19.2.9` Cross-Platform Market State / Modes / Readiness / Explanations
- `v19.2.10` Tenant Usage / Cost / Rights / Entitlement Observability
- `v19.2.11` Multi-User Fairness / Rate Limits / Noisy-Neighbor Controls
- `v19.2.12` Multi-User Security / Abuse / Tenant-Isolation Hardening
- `v19.2.13` Mixed-Client Compatibility Enforcement
- `v19.2.14` Hosted Recovery / DR / Secret-Rotation Drill
- `v19.2.15` Protected-Session Load / Capacity / Outage Assurance
- `v19.2.16` #66 Cross-Platform Assurance Closure

### v19.3.x — Point-in-Time Evidence Foundation
- `v19.3.0` Institutional / 13F Source Model + Provenance
- `v19.3.1` 13F Ingest / Backfill / Amendment / Revision
- `v19.3.2` 13F Point-in-Time Query / Snapshot
- `v19.3.3` Two-Sided Long / Short Evidence Substrate
- `v19.3.4` AODR Candidate Lineage
- `v19.3.5` AODR Ranking / Explanation Lineage
- `v19.3.6` AODR Outcome / Miss Lineage
- `v19.3.7` v19.3 Point-in-Time Evidence Zero-Gap Closure

### v19.4.x — Reliability / Economics / Adaptive Readiness
- `v19.4.0` Hosted SLO / Error Budget / Failure Classification
- `v19.4.1` Operational Runbooks / Incident / Rollback / Kill Readiness
- `v19.4.2` Measured Capacity / Cost Economics
- `v19.4.3` Provider License / Plan / Paid-Gap Evaluation
- `v19.4.4` Adaptive Evidence / Provenance Readiness
- `v19.4.5` v20 Research-Readiness Audit
- `v19.4.6` v19 Pre-Closure Zero-Gap Sweep

### v19.5.0 — Major Closure
- `v19.5.0` no-feature Major Closure. Require #66 PASS, every conservation row reconciled, zero material Mac/Windows/Web parity debt, rights/identity/RBAC/product-entitlement/privacy separation, lifecycle/IaC/supply-chain/API/recovery/SLO/capacity proof and actual supported artifact/deployment evidence.

## 8. v20 — Governed Adaptive Intelligence — provisional granular reservations

v20 cannot start before `v19.4.5` research-readiness and `v19.5.0` Major Closure. These labels remain provisional and must be re-audited against v19 evidence before G1.

### v20.0.x — Adaptive control-plane foundation
- `v20.0.0` Adaptive Research Control Plane Boundary
- `v20.0.1` Immutable Experiment Ledger
- `v20.0.2` Evidence Snapshot / Dataset Provenance Reproducibility
- `v20.0.3` Model Registry / Version / Approval Governance
- `v20.0.4` Prompt / Template Version Governance
- `v20.0.5` Champion / Challenger / Shadow Evaluation
- `v20.0.6` Historical Analogue Retrieval
- `v20.0.7` Regime-Conditioned Outcome Store
- `v20.0.8` Calibration / Confidence Reliability
- `v20.0.9` FP / FN / Miss / Contradiction / Drift Registry
- `v20.0.10` v20.0 Adaptive Foundation Zero-Gap Closure

### v20.1.x — ASBI
- `v20.1.0` ASBI Evidence Feature Normalization
- `v20.1.1` ASBI Candidate Synthesis
- `v20.1.2` ASBI Contradiction / Abstention Handling
- `v20.1.3` ASBI Confidence / Explanation Contract
- `v20.1.4` ASBI Outcome Feedback
- `v20.1.5` ASBI Learning Guardrails / Promotion Boundaries
- `v20.1.6` ASBI Zero-Gap Closure

### v20.2.x — Adaptive Institutional / TDTI
- `v20.2.0` Adaptive Institutional / 13F Feature Extraction
- `v20.2.1` Institutional Revision / Lag / Amendment Semantics
- `v20.2.2` TDTI Two-Sided Thesis Feature Model
- `v20.2.3` Institutional/TDTI Regime Conditioning
- `v20.2.4` Institutional/TDTI Outcome Calibration
- `v20.2.5` v20.2 Institutional/TDTI Zero-Gap Closure

### v20.3.x — AODR
- `v20.3.0` AODR Candidate Ranking Policy
- `v20.3.1` AODR Opportunity Scoring
- `v20.3.2` AODR Why / Why-Not Explanation
- `v20.3.3` AODR Outcome / Miss Learning
- `v20.3.4` AODR Stability / Drift / Fairness Guardrails
- `v20.3.5` v20.3 AODR Zero-Gap Closure

### v20.4.x — Adaptive operations
- `v20.4.0` Adaptive Provider/Evidence Selection — SHADOW
- `v20.4.1` Quality / Cost / Freshness Utility Weighting — SHADOW
- `v20.4.2` Dynamic Budget / Backpressure Recommendations
- `v20.4.3` Validated Policy Promotion Framework
- `v20.4.4` Adaptive Rollback / Drift / Kill Controls
- `v20.4.5` v20.4 Adaptive Operations Zero-Gap Closure

### v20.5.0
- `v20.5.0` v20 Professional Closure — no feature scope.

Every shared adaptive user-facing capability follows Mac + Windows + Web lockstep. Deterministic market truth, No Execution, provider lifecycle promotion and Router authority remain protected.

## 9. G0-G16 enforcement

G1 requirement conservation + platform matrix; G2 canonical owner/adapters; G3 one contract + equivalence/negative/failure tests; G4 all required implementations; G6 cross-platform integration; G7 security/data-rights/privacy; G8 load/capacity/recovery; G9 function/meaning equivalence; G10 zero unowned coverage; G11/G12 immutable RC certification; G13/G14 actual artifacts/deployments; G15 no GA with unresolved applicable rows; G16 zero-gap learning/handoff.

## 10. Permanent principle

**Build shared truth once -> expose each capability across all applicable supported clients together -> prove requirement conservation + equivalence -> only then advance the product.**
