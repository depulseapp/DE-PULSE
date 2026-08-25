# DE.PULSE — Canonical Adaptive Roadmap

**Status:** APPROVED / ADAPTIVE — v19/v20 rebaseline candidate under #110  
**Authority:** canonical product sequencing and approved strategic workstreams  
**Rule:** shipped releases are immutable truth; future reservations may adapt only with durable rationale, requirement conservation and synchronized governance/handoff updates.

## 1. Canonical truth model

- Stable tags, release evidence, source/artifact provenance and current handoff define what actually shipped.
- Historical provisional version labels do not prove implementation.
- Corrective/security/reliability/privacy/supply-chain work may preempt future sequencing.
- Known misses are fixed or durably assigned.
- G0-G16 is the only release model.
- Future G0/G1 performs source-overlap classification before implementation: `INHERITED`, `EXTEND_EXISTING_OWNER`, `REPLACE_OR_CONSOLIDATE`, or `NEW_RESIDUAL`.
- `governance/v19-v20-requirement-conservation.json` prevents issue requirements from disappearing between versions.

Historical provisional `v18.3 PostgreSQL / v18.4 hosted-security` placement remains superseded. Hosted PostgreSQL/shared-account authority belongs to the dependency-correct v19 train.

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

DE.PULSE is **one product across macOS Apple Silicon, Windows x64 and Web**.

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

## 4. Durable strategic workstreams

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

## 5. v18 — CLOSED / immutable history

v18 is closed by executable evidence. Do not reopen old labels from roadmap prose alone.

Retained completed foundations include:
- #57 / v18.8.2 Market Intelligence freshness/recovery repair;
- #61 TradeInsight;
- #64 / `v18.9.1-stable` runtime corrective;
- #70 CI convergence;
- #76 TradeInsight Settings/API-key UX;
- #79/#84 Provider/Data Health and Smart Provider Router v2 production closure;
- #92 canonical identity;
- #94 provider observability/usefulness;
- #95 provider onboarding/adaptation;
- #102 Post-Stable continuity;
- #107 Professional Closure and parent #65 completion.

Any newly discovered v18 regression must be proven against current `main`, checked against later superseding work, inserted into the requirement-conservation ledger, and resolved as a separate bounded corrective before v19 proceeds. There is currently no such open corrective.

## 6. v19 — Professional Hosted Product

Detailed responsibilities/dependencies/closure criteria are authoritative in `governance/V19_V20_ZERO_MISS_PLAN.md` and machine-mapped in `governance/v19-v20-requirement-conservation.json`.

### v19.0.x — Governance / identity / privacy / infrastructure / data foundations

1. `v19.0.0` Provider Legal/Data-Rights Registry Extension
2. `v19.0.1` Provider-Right Runtime Policy Contract
3. `v19.0.2` Tenant/Account Identity Boundary
4. `v19.0.3` Device Registration/Revocation/Lifecycle
5. `v19.0.4` Session Lifecycle + Privileged Re-authentication
6. `v19.0.5` Hosted RBAC/Capability Authorization Adapter
7. `v19.0.6` DE.PULSE Product-Entitlement Policy
8. `v19.0.7` Product Quota/Grace/Suspension/Downgrade Semantics
9. `v19.0.8` Account-Data Classification/Minimization/Retention
10. `v19.0.9` Account Export/Deactivation/Deletion/Tombstone Lifecycle
11. `v19.0.10` Operator Access/Audit-Log Retention/Residency Disposition
12. `v19.0.11` Hosted Environment/IaC/Service-Trust Foundation
13. `v19.0.12` Deployment Drift/Rollback/Provenance
14. `v19.0.13` PostgreSQL Tenant/Schema/Migration Authority
15. `v19.0.14` PostgreSQL Pool/Index/Capacity/Transaction Limits
16. `v19.0.15` PostgreSQL HA/Backup/PITR/RPO-RTO/Restore
17. `v19.0.16` Managed Secret/KMS Canonical Store/Reference Model
18. `v19.0.17` Secret Rotation/Revoke/Health/Audit/No-Leak Lifecycle
19. `v19.0.18` Dependency Inventory/SBOM/License/Vulnerability Policy
20. `v19.0.19` Artifact Integrity/Attestation/Patch-Revocation
21. `v19.0.20` Provider Quality/Cost/Coverage/SLO Scorecards
22. `v19.0.21` Reconciliation/Revision/Point-in-Time Data Quality
23. `v19.0.22` **Zero-Gap Foundation Closure — no feature scope**

### v19.1.x — Hosted zero-key data plane + sync foundation

1. `v19.1.0` Hosted API/Stream Inventory + Lifecycle
2. `v19.1.1` Authenticated Zero-Key Provider Gateway Boundary
3. `v19.1.2` Unified Hosted Serving Authorization Chain
4. `v19.1.3` REST/Snapshot Projection + Entitlement-Safe Cache Reuse
5. `v19.1.4` WebSocket/SSE Live Fan-Out Isolation
6. `v19.1.5` Rights/Entitlement Downgrade/Expiry Serving Behavior
7. `v19.1.6` Per-Account/User/Device Rate, Abuse + Protected-Session Capacity
8. `v19.1.7` Typed Sync Protocol Envelope/Domain Registry/Version Negotiation
9. `v19.1.8` SQLite Outbox + Idempotency/Replay
10. `v19.1.9` Server Revision/Change-Sequence + Durable Checkpoint
11. `v19.1.10` New/Stale-Device Bootstrap + High-Watermark/Re-bootstrap
12. `v19.1.11` Tombstone/Change-Log Retention, Compaction + Inactive-Device Expiry
13. `v19.1.12` Domain Conflict Policy + Mixed-Version Unsupported-Client Behavior
14. `v19.1.13` **Zero-Gap Data-Plane/Sync Closure — no feature scope**

### v19.2.x — Cross-platform account/state product + #66 assurance

1. `v19.2.0` Mac + Windows + Web Account/Session Client Foundation
2. `v19.2.1` Native Local-Account Isolation/Secure Credentials/Lost-Device/Offline Policy
3. `v19.2.2` Cross-Platform Portable Preferences
4. `v19.2.3` Cross-Platform Watchlists/Master Symbols
5. `v19.2.4` Cross-Platform Desks/Workspaces
6. `v19.2.5` Saved Searches/Notes/Product-Owned Research State
7. `v19.2.6` Lawful Durable Research/Evidence Portability
8. `v19.2.7` Cross-Platform Discovery/Opportunity Radar
9. `v19.2.8` Cross-Platform Market State/Modes/Readiness/Explanations
10. `v19.2.9` Cross-Platform Settings/Account/Device Controls
11. `v19.2.10` Cross-Platform RBAC/Product-Entitlement UX
12. `v19.2.11` Tenant-Aware Usage/Cost/Health Observability
13. `v19.2.12` Multi-User Security/Abuse/Noisy-Neighbor/Capacity Hardening
14. `v19.2.13` Mixed-Client Adversarial/Failure/Recovery Drills
15. `v19.2.14` **#66 Cross-Platform Zero-Gap Closure — no feature scope**

#66 cannot close on one-client success. Closure requires zero material REQUIRED-platform parity debt plus all zero-key, sync, security, privacy, rights, product-entitlement, IaC, supply-chain, DR and protected-session evidence.

### v19.3.x — Point-in-Time Evidence

1. `v19.3.0` Institutional/13F Ingestion + Provenance
2. `v19.3.1` 13F Revision/History/Query Model
3. `v19.3.2` Two-Sided Thesis Evidence Substrate / TDTI Foundation
4. `v19.3.3` TDTI Evidence Quality/Contradiction/Explanation
5. `v19.3.4` AODR Candidate/Ranking Lineage
6. `v19.3.5` AODR Outcome/Evaluation Lineage Substrate
7. `v19.3.6` **Zero-Gap Point-in-Time Evidence Closure — no feature scope**

### v19.4.x — Professional Reliability / Economics / v20 Readiness

1. `v19.4.0` ADR-GDI SLO/Capacity/Failure Taxonomy
2. `v19.4.1` ADR-GDI Degradation/Failover/Runbooks/Canary Controls
3. `v19.4.2` Specialized/Paid-Provider Gap Inventory
4. `v19.4.3` Provider Economics/License-Suitability/Upgrade Thresholds
5. `v19.4.4` v20 Adaptive-Research Readiness Audit
6. `v19.4.5` **Zero-Gap Reliability/Readiness Closure — no feature scope**

### `v19.5.0` — Major Closure

No feature scope. Requires every prior v19 band closure, #66 PASS, zero material Mac/Windows/Web parity debt, rights/identity/RBAC/product-entitlement/privacy separation, data lifecycle, API compatibility, PostgreSQL and secret recovery, IaC/environment and supply-chain assurance, SLO/capacity/DR proof, actual supported artifact/deployment evidence and a fresh implementation-miss audit.

## 7. v20 — Governed Adaptive Intelligence

Every shared adaptive user-facing capability follows Mac + Windows + Web lockstep. Production influence remains `SHADOW -> VALIDATED -> APPROVED -> PRODUCTION`. No silent self-modification; No Execution.

### v20.0.x — Adaptive Research / Model Governance Control Plane

1. `v20.0.0` Adaptive Research Control Plane Boundaries/Consumers
2. `v20.0.1` Immutable Experiment/Evidence/Outcome Ledger
3. `v20.0.2` Model/Prompt/Version Registry + Reproducibility
4. `v20.0.3` Champion/Challenger + SHADOW Evaluation
5. `v20.0.4` Human-Approved Promotion/Rollback/Kill Controls
6. `v20.0.5` Calibration/FP-FN/Miss/Contradiction/Drift Baseline
7. `v20.0.6` **Zero-Gap Adaptive-Control Closure — no feature scope**

### v20.1.x — ASBI

`ASBI` remains the historical roadmap name until its G0 source/issue audit freezes expanded semantics.

1. `v20.1.0` ASBI Data/Feature Eligibility + Rights/Provenance
2. `v20.1.1` ASBI SHADOW Inference
3. `v20.1.2` ASBI Evaluation/Calibration/Outcome Linkage
4. `v20.1.3` ASBI Explanation/Correlation/Contradiction/Abstention
5. `v20.1.4` ASBI Controlled Promotion
6. `v20.1.5` **ASBI Zero-Gap Closure — no feature scope**

### v20.2.x — Adaptive Institutional/13F + TDTI

1. `v20.2.0` Adaptive Institutional/13F SHADOW Features
2. `v20.2.1` Institutional Outcome/Regime Calibration
3. `v20.2.2` Adaptive TDTI Two-Sided Evidence Synthesis in SHADOW
4. `v20.2.3` TDTI Contradiction/Confidence/Abstention Calibration
5. `v20.2.4` Controlled Approved Institutional/TDTI Production Influence
6. `v20.2.5` **Institutional/TDTI Zero-Gap Closure — no feature scope**

### v20.3.x — AODR Adaptive Opportunity Intelligence

1. `v20.3.0` Adaptive Candidate/Ranking SHADOW Layer over Opportunity Radar
2. `v20.3.1` Outcome-Linked Opportunity Prioritization Learning
3. `v20.3.2` Regime-Aware Contextual Weighting
4. `v20.3.3` Explanation/Contradiction/Abstention + False-Positive Control
5. `v20.3.4` Controlled Approved Production Influence
6. `v20.3.5` **AODR Zero-Gap Closure — no feature scope**

### v20.4.x — ADR-GDI Adaptive Operations

1. `v20.4.0` Adaptive Degradation/Fallback Recommendations in SHADOW
2. `v20.4.1` Provider/Evidence Usefulness Weighting Advisory Layer
3. `v20.4.2` Adaptive Resource/Cost/Capacity Scheduling Recommendations
4. `v20.4.3` Human-Approved Operational Policy Promotion/Rollback
5. `v20.4.4` Drift/Incident/Outage Learning + Runbook Feedback
6. `v20.4.5` **ADR-GDI Adaptive Operations Zero-Gap Closure — no feature scope**

### `v20.5.0` — Professional Closure

No feature scope. Requires all v20 band closures, immutable experiment/model/prompt/outcome provenance, calibration and abstention evidence, explicit contradiction behavior, approved promotion history, rollback/kill proof, zero deterministic-owner regression, zero Router/rights/sync ownership leakage, cross-platform lockstep for shared adaptive capabilities and zero unexplained conservation-ledger gap.

## 8. Zero-miss sequencing and closure rule

Every planned implementation version has one primary responsibility. Necessary adapters/tests for REQUIRED platforms are part of that responsibility, but unrelated domains are not bundled simply to preserve a version number.

Each dependency band ends with a no-feature closure checkpoint. A closure version:
- cannot be used to hide omitted feature implementation;
- reconciles every conservation row targeted at or before the band;
- blocks advancement for `UNASSIGNED`, `OPEN_WITHOUT_TARGET`, duplicate owner, unexplained carry-forward or missing evidence;
- may send a discovered defect back to its owning slice/contract rather than inventing another parallel owner.

## 9. Permanent principle

**Build shared truth once -> conserve every approved requirement -> implement one bounded responsibility -> expose shared capabilities across all applicable supported clients together -> prove equivalence -> close the dependency band with zero unexplained gaps -> only then advance.**
