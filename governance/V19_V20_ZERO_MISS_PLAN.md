# DE.PULSE — v19/v20 Zero-Miss Program Plan

**Status:** PLANNING AUTHORITY CANDIDATE — issue #110  
**Baseline main:** `6aef3806d5684cc75daec0a2274bbf51fe135201`  
**Certified Stable:** `v18.9.1-stable`  
**Primary hosted umbrella:** #66 / `ADAPT-HOSTED-SYNC-001`  
**Planning rule:** one primary responsibility per slice; no unexplained carry-forward; no product implementation is authorized by this planning document alone.

## 1. Why the train is being re-baselined

The v18 program is closed by executable evidence. Future v19/v20 work was already directionally correct, but several old version labels bundled too many independent contracts into one patch. That creates a repeatable implementation-miss risk: a broad slice can look complete while identity, device, replay, retention, privacy, recovery, rights, or cross-platform parity subcontracts remain unfinished.

This plan therefore uses smaller dependency-ordered slices. A later slice may not absorb an unfinished prerequisite merely to preserve version numbering.

Every future implementation G0/G1 must classify existing source first as one of:

- `INHERITED` — already implemented and still valid by current executable evidence;
- `EXTEND_EXISTING_OWNER` — canonical owner exists but the planned contract has a real residual;
- `REPLACE_OR_CONSOLIDATE` — duplicate/obsolete ownership must be consolidated before adding behavior;
- `NEW_RESIDUAL` — no adequate owner exists and a new bounded owner is actually required.

Permanent engineering order remains **REUSE → CONSOLIDATE → REFACTOR → DELETE/REPLACE → ADD**.

## 2. v18 closure guard

v18 is treated as closed unless fresh executable evidence proves otherwise.

Retained completed v18 foundations include #57 Market Intelligence freshness repair, #61 TradeInsight, #64 v18.9.1 runtime corrective, #70 CI convergence, #76 TradeInsight Settings/API-key UX, #79/#84 Provider/Data Health, #92 canonical identity, #94 provider observability/usefulness, #95 provider onboarding, #102 Post-Stable continuity and #107 Professional Closure.

A historical issue label or old roadmap reservation does not reopen v18. Any newly discovered v18 regression must satisfy all of the following before it can interrupt v19:

1. current `main` reproduces the defect or current source/evidence proves the requirement is genuinely unimplemented;
2. existing later work does not already supersede or close it;
3. it is recorded in `governance/v19-v20-requirement-conservation.json` as a named corrective blocker;
4. it receives its own bounded corrective slice and executable closure evidence.

No unexplained v18 carry-forward may be hidden inside v19.

## 3. Permanent cross-platform rule

DE.PULSE is one product across **macOS Apple Silicon + Windows x64 + Web**. Shared capability work is planned by capability, not by client.

For every shared slice:

- G1 freezes Mac / Windows / Web as `REQUIRED` or justified `N/A`;
- one canonical domain/API/state contract is implemented;
- all REQUIRED adapters/surfaces are part of that same responsibility;
- platform differences are limited to OS/browser mechanics;
- business logic, account/state semantics, authorization, product entitlement, provider-right decisions, freshness/provenance and explanation meaning cannot fork;
- one-platform success is diagnostic only, never GA/Delivered proof;
- the next shared domain cannot start while material parity debt remains.

## 4. v19 — Professional Hosted Product

### v19.0.x — Governance, identity, privacy, infrastructure and data foundations

| Version | Primary responsibility | Depends on | Not complete until |
|---|---|---|---|
| `v19.0.0` | Provider legal/data-rights registry extension | v18 provider registration + rights owners | Commercial/multi-user, proxying, cache/retention, display/redistribution, derived/AI, concurrency, attribution, plan/environment, effective/expiry/renewal metadata are evidence-bound and fail closed. |
| `v19.0.1` | Provider-right runtime policy contract | v19.0.0 | Canonical cache/persistence/serving interfaces can deterministically enforce right downgrade/expiry without creating a second router or storage owner. |
| `v19.0.2` | Tenant/account identity boundary | v19.0.1 | Every hosted request/data owner has explicit tenant/account context and isolation semantics. |
| `v19.0.3` | Device registration/revocation/lifecycle | v19.0.2 | Device identity, naming, last-seen, revoke/lost/stale retirement and device attribution are canonical. |
| `v19.0.4` | Session lifecycle + privileged re-authentication | v19.0.2-v19.0.3 | Expiry/refresh/revocation, sensitive-action re-authentication, privileged MFA/passkey-class policy and recovery governance are explicit. |
| `v19.0.5` | Hosted RBAC/capability authorization adapter | v19.0.2-v19.0.4 | `SUPER_OWNER / OWNER / ADMIN / USER / DEMO` and capability truth are enforced consistently at API/stream/UI boundaries. |
| `v19.0.6` | DE.PULSE product-entitlement policy | v19.0.2, v19.0.5 | Customer plan/status/feature entitlement is a canonical control separate from RBAC and provider rights. |
| `v19.0.7` | Product quota/grace/suspension/downgrade semantics | v19.0.6 | Quota keys, grace, suspension, disabled, upgrade/downgrade and deterministic enforcement/metering semantics are frozen. |
| `v19.0.8` | Account-data classification/minimization/retention | v19.0.2 | Product/account data classes have purpose, consumer, minimization and retention ownership. |
| `v19.0.9` | Account export/deactivation/deletion/tombstone lifecycle | v19.0.8 | Review/export, deactivate/delete and durable deletion/tombstone behavior is executable across owned account data. |
| `v19.0.10` | Operator access, audit/log retention and residency disposition | v19.0.8-v19.0.9 | Support/operator access is auditable, logs are minimized/retained by policy, and region/residency assumptions are explicit. |
| `v19.0.11` | Hosted environment/IaC/service-trust foundation | v19.0.2 | Dev/test/stage/prod isolation, versioned desired state, service identities, ingress/egress/network/TLS trust and environment-specific secret references exist. |
| `v19.0.12` | Deployment drift/rollback/provenance | v19.0.11 | Reproducible bootstrap/deploy, drift detection, rollback/fail-forward and environment deployment-state provenance are executable. |
| `v19.0.13` | PostgreSQL tenant/schema/migration authority | v19.0.2, v19.0.8-v19.0.12 | Shared hosted authority is tenant-scoped, migration-owned and does not create raw SQLite↔PostgreSQL replication or dual-master tables. |
| `v19.0.14` | PostgreSQL pool/index/capacity/transaction limits | v19.0.13 | Bounded pools, indexes, transaction limits and capacity controls are measured and fail safely. |
| `v19.0.15` | PostgreSQL HA/backup/PITR/RPO-RTO/restore | v19.0.13-v19.0.14 | HA/failover disposition, encrypted backup, PITR, restore drills, migration rollback/roll-forward and server-rollback reconciliation are proven. |
| `v19.0.16` | Managed secret/KMS canonical store/reference model | v19.0.11 | Hosted provider secrets are server-side only, environment-scoped, least-privilege and version-referenced. |
| `v19.0.17` | Secret rotation/revoke/health/audit/no-leak lifecycle | v19.0.16 | Staged rotation, rollback, compromise revoke, health cutover, audit and no plaintext in logs/traces/backups are proven without client redeploy. |
| `v19.0.18` | Dependency inventory/SBOM/license/vulnerability policy | v19.0.11-v19.0.12 | Direct/transitive components, trusted sources, licenses, vulnerabilities and SBOM policy are executable. |
| `v19.0.19` | Artifact integrity/attestation/patch-revocation | v19.0.18 | Source→build→artifact→environment traceability, integrity/signing/attestation where applicable and vulnerable-component patch/revoke process exist. |
| `v19.0.20` | Provider quality/cost/coverage/SLO scorecards | v18 telemetry/Data Health, v19.0.0-v19.0.1 | Scorecards combine measured health/freshness/latency/disagreement/headroom/rights/cost/usefulness without creating routing authority. |
| `v19.0.21` | Reconciliation/revision/point-in-time data quality | v19.0.20 | Revision history, reconciliation, provenance and no-lookahead point-in-time quality rules are canonical. |
| `v19.0.22` | **v19.0 Zero-Gap Foundation Closure** | v19.0.0-v19.0.21 | No feature scope. Every v19.0 conservation row is evidenced or explicitly blocked externally; no unexplained prerequisite remains before hosted activation. |

### v19.1.x — Hosted zero-key data plane and sync foundation

| Version | Primary responsibility | Depends on | Not complete until |
|---|---|---|---|
| `v19.1.0` | Hosted API/stream inventory + lifecycle | v19.0.22 | API/protocol versioning, compatibility matrix, deprecation/removal, supported range and truthful unsupported-client behavior are frozen. |
| `v19.1.1` | Authenticated zero-key Provider Gateway boundary | v19.1.0, v19.0.16-v19.0.17 | Normal clients authenticate only to DE.PULSE; server-side gateway wraps existing canonical owners and never exposes platform provider keys. |
| `v19.1.2` | Unified hosted serving authorization chain | v19.0.0-v19.0.7, v19.1.1 | Tenant → session/device → RBAC → product entitlement → provider rights → freshness/data class → cache/persistence → Router v2 ordering is executable and distinguishable. |
| `v19.1.3` | REST/snapshot projection + entitlement-safe cache reuse | v19.1.2 | Shared cached/provider evidence is returned only when tenant/product/right/data-class rules permit it; identical lawful demand coalesces. |
| `v19.1.4` | WebSocket/SSE live fan-out isolation | v19.1.2-v19.1.3 | Existing multi-feed subscription owner remains upstream authority; downstream fan-out re-authorizes and terminates promptly on revocation/downgrade. |
| `v19.1.5` | Rights/entitlement downgrade/expiry serving behavior | v19.1.3-v19.1.4 | Cache/persistence/REST/live behavior after provider-right or product-plan downgrade is deterministic: retain/delete/quarantine/derived-only/block as governed. |
| `v19.1.6` | Per-account/user/device rate, abuse and protected-session capacity | v19.1.1-v19.1.5 | Quotas, throttling, fairness, abuse controls and Tier-0 market-session protection are enforced without coupling billing to Router ranking. |
| `v19.1.7` | Typed sync protocol envelope/domain registry/version negotiation | v19.0.13-v19.0.15, v19.1.0 | Stable IDs, tenant/user/device context, schema/protocol/domain versions and compatibility negotiation are explicit. |
| `v19.1.8` | SQLite outbox + idempotency/replay | v19.1.7 | Local mutation+outbox atomicity, idempotency scope/lifetime, request fingerprint and duplicate semantic replay are restart-safe. |
| `v19.1.9` | Server revision/change-sequence + durable checkpoint | v19.1.7-v19.1.8 | PostgreSQL applies transactionally, server sequence is authoritative and device checkpoints advance only after durable local apply. |
| `v19.1.10` | New/stale-device bootstrap + high-watermark/re-bootstrap | v19.1.9 | Snapshot/page stream, high-watermark, delta continuation and non-destructive checkpoint-expired re-bootstrap converge. |
| `v19.1.11` | Tombstone/change-log retention, compaction and inactive-device expiry | v19.1.9-v19.1.10, v19.0.8-v19.0.10 | Offline-safe retention/compaction is bounded and deletion is never inferred from absence. |
| `v19.1.12` | Domain conflict policy + mixed-version unsupported-client behavior | v19.1.7-v19.1.11 | Optimistic concurrency and domain-specific conflict rules replace universal LWW; unsupported clients preserve pending outbox and receive truthful upgrade/degraded behavior. |
| `v19.1.13` | **v19.1 Zero-Gap Data-Plane/Sync Closure** | v19.1.0-v19.1.12 | No feature scope. Gateway, REST/live serving and sync foundation have zero unexplained contract gaps before cross-platform domain activation. |

### v19.2.x — Cross-platform account/state product and #66 assurance

| Version | Primary responsibility | Depends on | Not complete until |
|---|---|---|---|
| `v19.2.0` | Mac + Windows + Web account/session client foundation | v19.1.13 | All REQUIRED clients implement equivalent account/session semantics against the same APIs. |
| `v19.2.1` | Native local-account isolation, secure credentials, lost-device/offline policy | v19.2.0, v19.0.3-v19.0.4 | Mac/Windows local SQLite and OS credentials cannot leak across users; lost/revoked device behavior and truthful offline rules are proven. |
| `v19.2.2` | Cross-platform portable preferences | v19.2.0-v19.2.1 | Preferences converge across Mac/Windows/Web with deterministic conflict/version rules. |
| `v19.2.3` | Cross-platform Watchlists/Master Symbols | v19.2.0-v19.2.1 | Membership/add/remove/tombstone semantics converge with no duplicate canonical symbol owner. |
| `v19.2.4` | Cross-platform Desks/Workspaces | v19.2.0-v19.2.1 | Desk membership/config/delete semantics converge and preserve canonical desk intelligence ownership. |
| `v19.2.5` | Saved searches/notes/product-owned research state | v19.2.0-v19.2.1 | Sync-safe user-owned research state has revision/history/conflict behavior that cannot silently lose work. |
| `v19.2.6` | Lawful durable research/evidence portability | v19.0.0-v19.0.1, v19.2.5 | Only entitled, provenance-bound, retention-safe evidence is shared; immutable market evidence is not accidentally deleted with account state. |
| `v19.2.7` | Cross-platform Discovery/Opportunity Radar | v19.2.0, v19.2.3, v19.1.3-v19.1.4 | Mac/Windows/Web expose equivalent Discovery/AODR-foundation meaning through canonical shared processing. |
| `v19.2.8` | Cross-platform Market State/Modes/Readiness/Explanations | v19.2.0, v19.1.3-v19.1.4 | Equivalent market-state truth and explanations across clients without changing protected deterministic ownership. |
| `v19.2.9` | Cross-platform Settings/Account/Device controls | v19.2.0-v19.2.1 | User/device/account controls are equivalent; provider-secret controls remain privileged and zero-key users never receive raw provider credentials. |
| `v19.2.10` | Cross-platform RBAC/Product-Entitlement UX | v19.0.5-v19.0.7, v19.2.0 | UI composition and backend authorization agree for all canonical roles/capabilities/plans. |
| `v19.2.11` | Tenant-aware usage/cost/health observability | v19.1.6-v19.1.12 | Outbox/lag/retry/conflict/DB/auth/device/plan/provider/cost/cache/stream/protected-session metrics are tenant-aware and distinguish failure domains. |
| `v19.2.12` | Multi-user security/abuse/noisy-neighbor/capacity hardening | v19.2.0-v19.2.11 | Object/function authorization negatives, bounded resources, fairness, backpressure, load shedding and provider/DB/gateway failure containment are proven. |
| `v19.2.13` | Mixed-client adversarial/failure/recovery drills | v19.2.12, v19.0.15, v19.0.17 | Revocation, downgrade, secret compromise, DB failover/PITR, old-client rejection, reconnect/replay and protected-session scenarios pass across applicable clients. |
| `v19.2.14` | **#66 Cross-Platform Zero-Gap Closure** | v19.2.0-v19.2.13 | No feature scope. #66 may close only with zero material Mac/Windows/Web parity debt and zero unexplained identity/RBAC/product-right/privacy/sync/recovery/security/capacity gap. |

### v19.3.x — Point-in-time evidence foundation

| Version | Primary responsibility | Depends on | Not complete until |
|---|---|---|---|
| `v19.3.0` | Institutional/13F ingestion + provenance | v19.2.14 | 13F source acquisition is canonical, provenance-bound and point-in-time safe. |
| `v19.3.1` | 13F revision/history/query model | v19.3.0 | Amendments/revisions and historical queries cannot leak future knowledge. |
| `v19.3.2` | Two-sided thesis evidence substrate / TDTI foundation | v19.3.1 | Bull/bear evidence is independently represented with source/time lineage. |
| `v19.3.3` | TDTI evidence quality/contradiction/explanation | v19.3.2 | Contradictory evidence and unknowns remain explicit; deterministic truth is not overwritten. |
| `v19.3.4` | AODR candidate/ranking lineage | v19.3.0-v19.3.3 | Candidate/ranking inputs, transformations and reasons are reproducible and provenance-bound. |
| `v19.3.5` | AODR outcome/evaluation lineage substrate | v19.3.4 | Outcome labels and evaluation windows are point-in-time safe and ready for later adaptive use without execution features. |
| `v19.3.6` | **v19.3 Zero-Gap Point-in-Time Evidence Closure** | v19.3.0-v19.3.5 | No feature scope. Institutional/TDTI/AODR evidence lineage is complete enough for reliability and later adaptive work. |

### v19.4.x — Professional reliability, economics and v20 readiness

| Version | Primary responsibility | Depends on | Not complete until |
|---|---|---|---|
| `v19.4.0` | ADR-GDI SLO/capacity/failure taxonomy | v19.2.14, v19.3.6 | Professional SLOs and provider/gateway/DB/sync/model failure classes are measured with protected-session budgets. |
| `v19.4.1` | ADR-GDI degradation/failover/runbooks/canary controls | v19.4.0 | Graceful degradation, failover, bounded queues, kill-switch/canary and operator runbooks are executable and assistant-independent. |
| `v19.4.2` | Specialized/paid-provider gap inventory | v19.0.20-v19.0.21, v19.4.0 | Remaining professional data-quality gaps are measured by capability, not provider brand hype. |
| `v19.4.3` | Provider economics/license-suitability/upgrade thresholds | v19.4.2 | Cost, rights, quality and reliability evidence define when paid/mature providers are justified. |
| `v19.4.4` | v20 adaptive-research readiness audit | v19.3.6-v19.4.3 | Evidence/outcome history, lineage, calibration prerequisites and protected deterministic boundaries are proven ready. |
| `v19.4.5` | **v19.4 Zero-Gap Reliability/Readiness Closure** | v19.4.0-v19.4.4 | No feature scope. No unexplained reliability/economics/readiness gap remains before Major Closure. |

### `v19.5.0` — Major Closure

No feature scope. Requires all v19 band closures, #66 PASS, zero material Mac/Windows/Web parity debt, identity/RBAC/product-entitlement/provider-right/privacy separation, privacy/data lifecycle, API compatibility, PostgreSQL/secret recovery, IaC/environment/supply-chain assurance, SLO/capacity/DR proof, actual supported artifacts/deployments and a fresh implementation-miss audit. Any applicable unassigned conservation row blocks v19 closure.

## 5. v20 — Governed Adaptive Intelligence

v20 may learn/adapt from trustworthy evidence, but deterministic market truth, provider routing authority, rights, sync and execution boundaries remain protected. No silent self-modification is allowed. Promotion remains `SHADOW → VALIDATED → APPROVED → PRODUCTION` with explicit evidence and rollback.

### v20.0.x — Adaptive research/model governance control plane

- `v20.0.0` — Adaptive Research Control Plane boundaries and canonical consumers.
- `v20.0.1` — Immutable experiment/evidence/outcome ledger.
- `v20.0.2` — Model/prompt/version registry and reproducibility.
- `v20.0.3` — Champion/challenger + SHADOW evaluation framework.
- `v20.0.4` — Human-approved promotion/rollback/kill-switch controls.
- `v20.0.5` — Calibration, false-positive/false-negative, miss, contradiction and drift baseline.
- `v20.0.6` — **v20.0 Zero-Gap Adaptive-Control Closure** — no feature scope.

### v20.1.x — ASBI

The historical roadmap name `ASBI` remains authoritative until its G0 source/issue audit freezes expanded semantics.

- `v20.1.0` — ASBI data/feature eligibility and rights/provenance contract.
- `v20.1.1` — ASBI SHADOW inference through canonical evidence owners.
- `v20.1.2` — ASBI evaluation/calibration/outcome linkage.
- `v20.1.3` — ASBI explanation/correlation/contradiction/abstention behavior.
- `v20.1.4` — ASBI bounded controlled promotion.
- `v20.1.5` — **ASBI Zero-Gap Closure** — no feature scope.

### v20.2.x — Adaptive Institutional/13F + TDTI

- `v20.2.0` — Adaptive Institutional/13F SHADOW features.
- `v20.2.1` — Institutional outcome/regime calibration.
- `v20.2.2` — Adaptive TDTI two-sided evidence synthesis in SHADOW.
- `v20.2.3` — TDTI contradiction/confidence/abstention calibration.
- `v20.2.4` — Controlled production promotion of approved Institutional/TDTI influence.
- `v20.2.5` — **Institutional/TDTI Zero-Gap Closure** — no feature scope.

### v20.3.x — AODR adaptive opportunity intelligence

- `v20.3.0` — Adaptive candidate/ranking SHADOW layer over canonical Opportunity Radar.
- `v20.3.1` — Outcome-linked opportunity prioritization learning.
- `v20.3.2` — Regime-aware contextual weighting without replacing deterministic market truth.
- `v20.3.3` — Explanation/contradiction/abstention and false-positive control.
- `v20.3.4` — Controlled approved production influence.
- `v20.3.5` — **AODR Zero-Gap Closure** — no feature scope.

### v20.4.x — ADR-GDI adaptive operations

- `v20.4.0` — Adaptive degradation/fallback recommendations in SHADOW only.
- `v20.4.1` — Provider/evidence usefulness weighting advisory layer; Smart Provider Router v2 remains sole routing owner.
- `v20.4.2` — Adaptive resource/cost/capacity scheduling recommendations under protected-session budgets.
- `v20.4.3` — Human-approved bounded operational policy promotion/rollback; no silent self-modification.
- `v20.4.4` — Drift/incident/outage learning and runbook feedback.
- `v20.4.5` — **ADR-GDI Adaptive Operations Zero-Gap Closure** — no feature scope.

### `v20.5.0` — Professional Closure

No feature scope. Requires all v20 band closures, immutable experiment/model/prompt/outcome provenance, calibrated evidence, explicit abstention/contradiction behavior, approved promotion history, rollback/kill evidence, zero protected deterministic-owner regression, zero Router/rights/sync ownership leakage, cross-platform lockstep for all shared user-facing adaptive capabilities and zero unexplained conservation-ledger gap.

## 6. Requirement-conservation rule

`governance/v19-v20-requirement-conservation.json` is the machine-readable map from GitHub issue requirements to these versions.

Every applicable requirement must have exactly one current disposition:

- `INHERITED` — already closed by trustworthy executable evidence;
- `IMPLEMENT_IN` — owned by exactly one planned version;
- `FUTURE_BLOCKED` — deliberately deferred with an explicit blocker/reason and named future owner.

`UNASSIGNED`, `OPEN_WITHOUT_TARGET`, duplicate primary ownership, unexplained carry-forward or a version missing from this plan is a planning failure.

Each dependency-band closure must reconcile every row targeted at or before that band. A zero-gap closure may not be used to implement omitted feature work; it verifies, fixes only discovered closure defects through their original owners, and blocks advancement until conservation is clean.

## 7. Implementation-start rule

This planning slice does **not** reserve or start a v19 product implementation. After #110 is merged and closed:

1. fetch live `main` and current issue/source truth;
2. inspect the first planned version against current owners/evidence;
3. classify overlap as `INHERITED`, `EXTEND_EXISTING_OWNER`, `REPLACE_OR_CONSOLIDATE`, or `NEW_RESIDUAL`;
4. if the first version is already satisfied, close it by evidence or advance to the first real residual rather than rewriting equivalent code;
5. reserve exactly one bounded product slice;
6. use exact-head Fast → impact-selected Qualified → expected-head merge and existing G0–G16 only.
