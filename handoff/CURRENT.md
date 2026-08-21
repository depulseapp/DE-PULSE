# DE.PULSE — Current Handoff

**SUPERSEDES ALL PRIOR CHAT HANDOFFS**

**GitHub source of truth:** `depulseapp/DE-PULSE`  
**Certified Stable:** `v18.9.0-stable`  
**Certified Stable candidate:** `9ea81cddae4875ae15d3719ca028519a36c597b6`  
**Certified Stable qualified source:** `9e86b5e731f7a585cc77c1521f3639fc7a208efc`  
**Certified Stable fingerprint:** `a8719090c341c874dbd1279cc31ad98e84075d5701c46a800bf951340780ecb9`  
**Certified Stable build ID:** `v18.9.0-stable-20260821`  
**Release PR:** #62 — merged  
**Completed release scope:** #61 / `ADAPT-TRADEINSIGHT-001` — closed completed  
**Active product development branch:** none  
**Active product PR:** none  
**Governance alignment PR:** #67 — draft  
**Master corrective program:** #65 / `ADAPT-PROVIDER-INTELLIGENCE-010`  
**Immediate blocker/next product patch:** #64 / `ADAPT-RUNTIME-CRASH-001` -> `v18.9.1`.

## v18.9.0 — COMPLETE / IMMUTABLE STABLE

Fast #481 / `32525637987`, Qualified #153 / `32525738828`, merged candidate `9ea81cddae4875ae15d3719ca028519a36c597b6`, and Release #32 / `32526121817` are immutable release authority. G11–G16, actual macOS Apple Silicon + Windows x64 packaged-runtime audits, G15 assurance and same-run publication passed. Stable tag: `v18.9.0-stable`. Durable manifest: `release/v18.9.0/stable-evidence-manifest.json`.

## Permanent release philosophy

DE.PULSE uses many small, dependency-ordered, independently certifiable patches rather than heavy multi-domain builds.

Permanent sequence:

`stabilize -> instrument -> validate -> expand -> harden -> close -> platform foundation -> scale -> learn`

- one primary responsibility per patch;
- G0–G16 only;
- canonical owners are reused rather than duplicated;
- observability needed for a capability exists before broad production admission;
- security/rights/identity/recovery foundations precede multi-user data-plane activation;
- point-in-time evidence precedes adaptive learning;
- model/prompt governance precedes broad adaptive production influence;
- known misses are fixed in-scope or durably assigned; no chat-only carry-forward.

## Version-alignment rule

- Major = strategic maturity stage.
- Minor band = coherent platform phase.
- Patch = one primary independently certifiable responsibility.
- Future version labels are planned reservations until G1; broad work may split and shift later unstarted reservations.
- Shipped versions are immutable.
- Corrective/security work preempts the planned train when necessary.

## Ordered v18.9.x patch train — architecture-audited alignment

The audit moved provider observability ahead of TradeInsight capability expansion so SHADOW admission is evidence-based.

1. `v18.9.1` — #64 runtime crash corrective ONLY.
2. `v18.9.2` — TradeInsight Settings/API-key UX ONLY.
3. `v18.9.3` — coverage-aware Smart Provider Router v2 + persistence-first residual-gap fulfillment ONLY.
4. `v18.9.4` — canonical company/instrument identity + all-desk presentation ONLY.
5. `v18.9.5` — Market Data Modes + capability diagnostics ONLY.
6. `v18.9.6` — provider observability + Adaptive Intelligence telemetry + protected-session headroom measurement ONLY.
7. `v18.9.7` — TradeInsight SEC Form 4 enrichment ONLY; SHADOW-first; direct SEC/EDGAR authoritative.
8. `v18.9.8` — TradeInsight ticker/company search ONLY.
9. `v18.9.9` — TradeInsight movers/ranking evidence ONLY.
10. `v18.9.10` — remaining useful TradeInsight capability admission ONLY.
11. `v18.9.11` — Session-Aware Data Readiness Maintenance ONLY: light overnight + heavy weekend; protected pre-market/regular/after-hours.
12. `v18.9.12` — whole v18.9.x professional closure audit ONLY; no new feature scope.

A future patch may split further at G0/G1 if still broad. Never merge unrelated work merely to preserve a planned number.

## Permanent adaptive provider + persistence contract

`consumer requirement -> in-memory canonical cache -> persisted canonical DB/state -> validate freshness/coverage/schema/provenance/rights -> exact residual gap -> eligible-provider ranking -> targeted acquisition -> canonical merge/reconciliation/provenance -> coverage re-evaluation -> persist -> next provider only if still needed -> synthesized consumer state`

Provider success does not equal consumer completeness. No fixed provider chain is authoritative. Validation lifecycle (`SHADOW/VALIDATED/...`) is separate from serving role (`PRIMARY/FALLBACK/BACKFILL/ENRICH/CORROBORATE/...`).

Never refetch/recompute trustworthy evidence already valid for the consumer. Revision-prone evidence preserves point-in-time/as-observed history plus later revisions. Live-sensitive values obey freshness truth and cannot become current merely because persisted.

Machine contract: `adaptive-governance/PERSISTENCE_REUSE_AND_OFF_HOURS_DATA_READINESS_CONTRACT.md`.

## Protected-session contract

Pre-market, regular market and after-hours are Tier-0 protected decision-support sessions with first claim on provider quota/headroom, network, CPU, memory, DB/pool and workers.

- light overnight work is small, high-value and gap-driven;
- heavy weekend/extended-closed work is deeper but bounded;
- maintenance/sync uses only surplus capacity after protected reserves;
- low-priority acquisition suspends/yields during protected sessions;
- maintenance/sync drains/preempts/checkpoints/resumes around live demand or market shocks;
- missed background work catches up only in later eligible windows;
- no blind full-universe refetch and no parallel calendar/scheduler/router/cache/database owner.

## Hosted / multi-device / zero-key architecture — issue #66

Approved v19 architecture:

```text
macOS / Windows
  local SQLite edge/offline working set
            │
            │ typed authenticated incremental sync
            ▼
      DE.PULSE hosted API/services
            │
            ├─ canonical account/session/capability truth
            ├─ Smart Provider Router v2
            ├─ canonical freshness/cache/persistence reuse
            ├─ rights/entitlement gate
            ├─ existing multi-feed subscription owner
            ├─ server-side managed provider secrets
            ▼
 PostgreSQL shared account/state authority
            ▲
            │
       hosted web
```

Permanent boundaries:
- SQLite remains native edge/offline store and warm working set.
- PostgreSQL is shared authority only for sync-eligible account/device state and lawful hosted evidence.
- Sync is application-level, typed, incremental, idempotent and checkpointed; never blind DB replication/dual-master tables.
- Normal commercial users authenticate only to DE.PULSE and have **zero provider API-key setup**.
- Platform provider credentials remain server-side in managed secrets/KMS and never appear in ordinary SQLite/Postgres sync rows, browser storage, client logs or normal telemetry.
- Hosted Provider Gateway wraps the existing Smart Provider Router v2; no second provider stack.
- REST/snapshot/WebSocket/SSE use the existing canonical provider/subscription owners and rights/entitlement enforcement.
- Cross-account/plan entitlement isolation applies before cache/persistence/fan-out reuse.
- Server-side canonical `SUPER_OWNER/OWNER/ADMIN/USER/DEMO` + capability truth applies across Mac/Windows/web/API/streams.
- New-device/stale-device bootstrap, outbox/idempotency/change-log/checkpoints/tombstones/compaction and mixed-version compatibility are mandatory.
- Lost/revoked devices, local account isolation, secret rotation, rights expiry, PostgreSQL PITR/restore and multi-user load/cost/security assurance are mandatory.
- Current/live market truth remains canonical freshness/provider/state-owned.

Governing contracts:
- `adaptive-governance/SQLITE_POSTGRES_SYNC_CONTRACT.md`
- `adaptive-governance/ROLE_AWARE_SESSION_SECURITY_CONTRACT.md`
- `adaptive-governance/ROLE_AWARE_UI_COMPOSITION_CONTRACT.md`
- issue #66 / `ADAPT-HOSTED-SYNC-001` and its architecture audit addenda.

## Planned v19 version train

Exact identity freezes at each G1; these are dependency-aligned reservations.

### v19.0.x — Governance / control plane / data foundation
- `v19.0.0` provider capability/entitlement/rights registry
- `v19.0.1` hosted identity/device/session control plane
- `v19.0.2` PostgreSQL tenancy/schema/pool/HA/PITR foundation
- `v19.0.3` managed secrets/KMS lifecycle
- `v19.0.4` provider quality/cost/coverage/SLO scorecards
- `v19.0.5` data reconciliation/revision/point-in-time quality

### v19.1.x — Zero-key provider data plane / native sync
- `v19.1.0` authenticated Hosted Provider Gateway
- `v19.1.1` runtime rights/entitlement + live fan-out isolation
- `v19.1.2` sync protocol foundation
- `v19.1.3` macOS preferences/watchlist pilot
- `v19.1.4` desks/workspaces sync

### v19.2.x — Cross-platform parity / #66 assurance
- `v19.2.0` Windows x64 parity
- `v19.2.1` hosted web parity
- `v19.2.2` rights-aware research/state portability
- `v19.2.3` multi-user security/cost/abuse/capacity hardening
- `v19.2.4` #66 hosted-sync/gateway adversarial/failure/recovery closure

### v19.3.x — Point-in-time evidence
- `v19.3.0` institutional/13F infrastructure
- `v19.3.1` two-sided Long/Short evidence substrate
- `v19.3.2` AODR candidate/ranking/outcome lineage

### v19.4.x — Reliability/economics/readiness
- `v19.4.0` ADR-GDI professional reliability/capacity
- `v19.4.1` specialized/paid-provider gap evaluation
- `v19.4.2` v20 research-readiness audit

### v19.5.0 — v19 Major Closure
No new feature scope. #66, rights/commercial posture, point-in-time quality, SLO/capacity, recovery/rollback and supported runtime/package evidence must PASS with zero unresolved P0 architecture gap.

## Planned v20 version train

v20 begins only after `v19.5.0` PASS and consumes point-in-time, rights-valid, provenance-bound evidence.

### v20.0.x — adaptive control/governance
- `v20.0.0` adaptive research control plane + immutable experiment ledger
- `v20.0.1` model/prompt governance + Champion/Challenger **before broad adaptive rollout**
- `v20.0.2` historical analogues/regime-conditioned outcomes
- `v20.0.3` calibration/FP-FN/miss/contradiction/drift

### v20.1.x — ASBI
- `v20.1.0` behavioral fingerprints/state transitions
- `v20.1.1` scenarios/probability momentum/calibration

### v20.2.x — institutional + TDTI
- `v20.2.0` adaptive 13F
- `v20.2.1` competing Long/Short/No Reliable Edge
- `v20.2.2` two-sided trade-plan/readiness/outcome validation; No Execution

### v20.3.x — AODR
- `v20.3.0` adaptive shared opportunity ranking
- `v20.3.1` diversity/opportunity-cost/personalized relevance after shared truth

### v20.4.x — adaptive operations
- `v20.4.0` ADR-GDI adaptive provider/recovery/workload/maintenance/reserve optimization under governed SHADOW/Champion-Challenger promotion

### v20.5.0 — v20 Professional Closure
No new feature scope. Calibration/utility/drift/abstention, deterministic boundaries, privacy/security/data rights, reproducibility/rollback, actual artifacts, zero silent self-modification and No Execution.

## Industry-strength controls inside G0–G16

Hosted/security/data/adaptive releases must absorb within existing gates:
- architecture decision + threat/data classification at G2;
- API/schema/protocol compatibility, SLO/observability, migration/rollback and failure-test plan at G3;
- contract tests + feature/kill-switch controls where appropriate at G4;
- negative authorization/tenant/rights/secret/adaptive tests at G7;
- load/soak/capacity/failure-injection/failover/protected-session proof at G8;
- role-aware cross-platform UX/direct-route/API consistency at G9;
- production-readiness/P0 blocking at G10;
- immutable RC/full certification G11/G12;
- package/provenance/runtime G13/G14;
- bounded/canary promotion for hosted-risk changes where applicable at G15;
- implementation-miss/incident/metric learning and cleanup at G16.

## Other continuity truth

Issue #57 / v18.8.1 Market Intelligence escape is closed completed because `release/v18.8.2/stable-evidence-manifest.json` records `RESOLVED_IN_V18.8.2_STABLE`; regression remains mandatory at v18.9.12 closure.

The real v18.9.0 macOS Apple Silicon `EXC_CRASH (SIGABRT)` remains unresolved under #64. Do not guess root cause or delete `PersonalMarketTerminal` state/API keys as a first troubleshooting step.

## Exactly one next action

Diagnose #64 using complete macOS crash evidence or deterministic reproduction and freeze the narrow `v18.9.1` G1. Do not create a `v18.9.2` or v19 product implementation branch until `v18.9.1` is truthfully closed or the crash is proven external/non-product.

## Resume rule

Any ChatGPT account, Codex session, Claude or human maintainer must first fetch the **current live GitHub head** because another session/process may have advanced it. Then read `AGENTS.md`, `CLAUDE.md`, `governance/AI-ASSISTANT-PORTABILITY-CONTRACT.md`, this file, all four `adaptive-governance/CURRENT_ADAPTIVE_*` files, `adaptive-governance/PERSISTENCE_REUSE_AND_OFF_HOURS_DATA_READINESS_CONTRACT.md`, `adaptive-governance/SQLITE_POSTGRES_SYNC_CONTRACT.md`, `adaptive-governance/ROLE_AWARE_SESSION_SECURITY_CONTRACT.md`, `adaptive-governance/ROLE_AWARE_UI_COMPOSITION_CONTRACT.md`, `adaptive-governance/ADAPTIVE_WORK_DECOMPOSITION_CONTRACT.md`, `adaptive-governance/ADAPTIVE_CI_OPERATING_CONTRACT.md`, `adaptive-governance/GOVERNANCE_IMPLEMENTATION_CLOSURE_CONTRACT.md`, `release_identity.json`, `release/v18.9.0/stable-evidence-manifest.json`, current resume checkpoints, issue #65, issue #64, issue #66 and current comments before changing code. Inspect commits since the certified Stable/source baseline so completed work is not duplicated. GitHub objects and executable evidence outrank chat memory.
