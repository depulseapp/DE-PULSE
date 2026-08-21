# DE.PULSE — Current Adaptive Build Process

**Certified Stable:** `v18.9.0-stable`  
**Current program:** #65 / `ADAPT-PROVIDER-INTELLIGENCE-010`  
**Immediate activity:** v18.9.1 corrective diagnosis for #64 / `ADAPT-RUNTIME-CRASH-001`  
**Active development branch:** none.

## Permanent process

**GitHub source of truth -> exact G0 baseline -> one bounded G1 scope -> G2 canonical-owner audit -> G3 dependency/contract readiness -> one version-development branch -> coherent code+tests -> one Draft PR -> exact-head Fast -> same PR Ready -> Qualified -> exact-head merge -> one canonical G11-G16 Release when release-capable -> post-release implementation-miss audit -> continuity reconciliation -> next patch.**

## Small-patch discipline

This discipline is permanent across v18.9.x, v19, v20 and later major versions.

1. One primary responsibility per patch.
2. Explicit non-goals are mandatory at G1.
3. Do not combine stability, routing architecture, provider admission, UX/domain redesign, data-infrastructure hardening and adaptive-model work in one patch.
4. Reuse existing owners first. Repair order remains `REUSE -> CONSOLIDATE -> REFACTOR -> DELETE/REPLACE -> ADD`.
5. Every patch must have deterministic acceptance tests tied directly to its G1 scope.
6. Before closure, re-audit implementation against the frozen scope and search for misses, bypasses, duplicate ownership and misleading UI/data truth.
7. Any newly discovered out-of-scope miss must receive a durable issue/target patch before closure; chat-only carry-forward is prohibited.
8. The next patch cannot begin until current patch evidence, issue state, handoff and checkpoints agree.
9. CI is intentionally economical: batch coherent source/test changes before PR, avoid duplicate/manual runs, classify failures before rerun, and never create retry/certification branch families.
10. If G0/G1 shows a planned packet is still too large, split it. A version number never justifies bundling.
11. G0-G16 remains the only gate model; no G17+.

## v18.9.x process sequence

- `v18.9.1`: #64 runtime crash only.
- `v18.9.2`: TradeInsight Settings/API-key UX only.
- `v18.9.3`: coverage-aware Smart Provider Router v2 core + persistence-first residual-gap fulfillment only.
- `v18.9.4`: canonical company identity/all-desk presentation only.
- `v18.9.5`: Market Data Modes/capability diagnostics only.
- `v18.9.6`: TradeInsight Form 4 enrichment only.
- `v18.9.7`: TradeInsight ticker/company search only.
- `v18.9.8`: TradeInsight movers/ranking evidence only.
- `v18.9.9`: remaining useful TradeInsight capability admission only.
- `v18.9.10`: provider efficiency/Adaptive Intelligence telemetry + protected-session headroom measurement only.
- `v18.9.11`: session-aware Data Readiness Maintenance only — light overnight + heavy weekend, with strict protection for pre-market/regular-market/after-hours.
- `v18.9.12`: professional closure audit only; no new feature scope.

## Persistence-first process contract

Provider work must begin from a consumer requirement and existing canonical evidence, not from a provider-first fetch loop.

Canonical decision sequence:

`consumer requirement -> in-memory canonical cache -> persisted canonical DB/state -> freshness/coverage/schema/provenance/rights validation -> residual gap -> Smart Provider Router -> targeted acquisition -> canonical merge/reconciliation -> persist -> serve`

Never refetch/recompute data already valid for the consumer solely because a provider is available. Revision-prone evidence must retain as-observed point-in-time history plus later revisions. Fast-changing evidence may be retained for history but cannot be reused as current truth beyond its freshness contract.

## Session-aware maintenance process contract

The existing canonical U.S. market calendar/session owner defines protected-session truth. Maintenance must not hard-code another market calendar.

### Protected Tier-0 sessions

**Pre-market, regular market and after-hours** are first-class decision-support sessions. During them:

- live/current data and intelligence outrank all maintenance and background synchronization;
- provider quota/headroom needed for current-session capability is reserved;
- maintenance external-provider calls are suspended unless directly required by a current/live consumer;
- CPU/memory/network/DB/background-worker capacity remains bounded and reserved for current-session work;
- maintenance and sync are preemptible/checkpointed and yield promptly to protected sessions, market shocks or high-priority current consumers;
- no deep historical fan-out, heavy reconciliation, compaction or broad background work may run.

### Light overnight process

After protected after-hours work ends and before the next protected pre-market window, the coordinator may run **small, high-value, gap-driven** work under conservative provider/runtime budgets: finalize completed-session data, fill small residual historical gaps, check incremental disclosures/revisions, reconcile small corporate-action/fundamental/macro/identity gaps, resolve bounded outcomes, run lightweight integrity/readiness checks and prepare warm canonical state for the next session.

It must stop/drain before the pre-market protection buffer, when provider headroom is inadequate or when runtime health deteriorates.

### Heavy weekend / extended market-closed process

Use larger non-trading windows for deeper but still bounded historical backfill/reconciliation, corporate actions, SEC/Form4/13F/congress/earnings/fundamental/macro history, point-in-time outcome resolution, provider-history consolidation, material feature preparation and DB/index/retention maintenance.

Heavy maintenance still requires named consumer/value, provider/data rights, rate/cost budgets and expected reuse. It is never a blind full-universe refetch.

### Catch-up and restart

If the app was not running during an eligible window, checkpointed/pending maintenance catches up only in a later eligible overnight/weekend period. Never dump maintenance backlog into pre-market/regular-market/after-hours. Restart/resume must avoid duplicate acquisition/work.

Manual maintenance uses the same budgets and cannot override protected-session safety.

Machine contract: `adaptive-governance/PERSISTENCE_REUSE_AND_OFF_HOURS_DATA_READINESS_CONTRACT.md`.

## v19 process objective

v19 is **Professional Data Infrastructure**, not another provider-router or maintenance rewrite. It consumes v18.9.x coverage-aware acquisition, persistence-first reuse, session-aware maintenance, identity and telemetry and adds measured capability/entitlement/rights, provider quality/SLO/cost scorecards, reconciliation/history/revision quality, 13F point-in-time infrastructure, TDTI/AODR evidence lineage, ADR-GDI professional reliability, DB/index/pool/capacity hardening, maintenance economics/protected-session reserve sizing, measured paid-provider gap analysis, and the mandatory hosted account/sync/provider-gateway program #66 / `ADAPT-HOSTED-SYNC-001`. v19 must end with a mandatory Major Closure before v20.

### Mandatory #66 process chain

#66 is executable v19 scope and MUST be represented in G0/G1 sequencing rather than left only in handoff/issue text. Its implementation is split into dependency-ordered small patches and cannot begin before the v18.9.x closure permits v19 work.

Required process ordering is:
1. reconcile #66 into the current Roadmap, Build Plan, Build Process and Delivery Process at the v19 transition and freeze exact dependencies;
2. prove canonical role/capability reuse (`SUPER_OWNER/OWNER/ADMIN/USER/DEMO`) and hosted account/user/device/session ownership; no duplicate identity/session truth;
3. prove PostgreSQL tenancy/schema/pool/HA/PITR/backup/restore/migration and recovery objectives before broad sync activation;
4. prove managed-secret/KMS lifecycle, environment isolation, rotation/rollback/compromise recovery and end-user zero-key deployment mode;
5. wrap the existing Smart Provider Router v2 and multi-feed subscription owner behind the authenticated hosted Provider Gateway; do not create a second provider or live-subscription stack;
6. make provider rights/entitlements machine-enforced across router/cache/persistence/REST/live fan-out before shared reuse is enabled;
7. implement sync transport as application-level typed protocol: bootstrap, SQLite atomic outbox, idempotency, authoritative server sequence, change log, checkpoint, tombstone, retention/compaction, stale-device re-bootstrap and mixed-version capability negotiation;
8. activate domains progressively: macOS preferences/watchlists first, then desks/workspaces, Windows parity, hosted-web parity, then rights-aware research/evidence;
9. prove local account isolation, user switching, offline/restart/reconnect, lost/revoked device behavior and no provider secret on commercial clients;
10. close with multi-user cost/usage/abuse/licensing/security/load/DR assurance and an implementation-miss audit.

Every #66 patch follows the normal G0-G16 model. G2 must map canonical owners; G3 must freeze rights/security/sync/conflict/protocol contracts; G7 must prove tenant/security/secret/rights isolation; G8 must prove DB/pool/load/recovery and protected-session safety; G9 must prove role-aware Mac/Windows/web behavior; G10 blocks freeze on any unresolved P0 audit item; G12 replays on the immutable RC; G14 validates affected native artifacts; G15 requires hosted-sync/provider-gateway assurance before promotion; G16 records incidents/misses and the next bounded packet.

Each v19 responsibility is a separate small packet unless G0 proves two items are inseparable. Point-in-time provenance, data rights, source independence, outcome lineage, hosted account isolation, provider-secret boundaries, sync correctness and protected-session reliability are release requirements, not documentation-only concerns.

## v20 process objective

v20 is **Adaptive Intelligence & Decision Research**. It may begin only after v19 Major Closure proves the evidence substrate and #66 hosted/synchronized account architecture are trustworthy enough for learning.

Adaptive work is split by responsibility: experiment/evidence ledger, historical analogues, calibration/drift, ASBI, adaptive 13F, TDTI, AODR, ADR-GDI adaptive optimization, model/prompt governance and final closure. Production influence follows `SHADOW -> VALIDATED -> APPROVED -> PRODUCTION`; no model or adaptive policy silently self-promotes. Deterministic Day/Swing/Long truth remains protected unless separately governed, and No Execution is permanent.

v20 may learn better provider/maintenance usefulness and reserve policies only through governed SHADOW/Champion-Challenger evidence. It may not reduce protected live-session safety, bypass data rights/provenance or sacrifice current truth for background learning.

## Cross-version dependency rule

`v18.9.x trustworthy acquisition/persistence/session readiness -> v19 measured professional evidence + hosted account/sync infrastructure -> v20 governed adaptive learning`.

v19 must not create a second router or maintenance coordinator to measure providers. #66 must reuse canonical provider/freshness/cache/persistence/subscription/session owners and must be operationally assured before synchronized evidence reaches v20. v20 must not learn from data lacking point-in-time provenance/rights or bypass canonical evidence owners. Weak/missing evidence remains UNKNOWN/ABSTAIN rather than being filled by model confidence.

## Adaptive provider process contract

The sole routing owner evaluates residual missing coverage and eligible providers, acquires only what is needed, merges with provenance, re-evaluates remaining gaps and stops only when the bounded requirement is met or eligible budget is exhausted.

A provider response marked successful does not imply consumer completeness. Static provider ordering is only a prior/tiebreaker. TradeInsight is never allowed to create its own router/cache/scanner/Market Mode/SEC truth/symbol/persistence/scheduler system. Hosted mode likewise may not create another provider router or subscription manager.

Provider validation lifecycle and runtime serving role are distinct concepts. SHADOW/VALIDATED/APPROVED describe evidence maturity; PRIMARY/FALLBACK/BACKFILL/ENRICH/CORROBORATE describe serving purpose. Promotion/demotion requires telemetry/evidence and must not silently alter deterministic Day/Swing/Long truth.

## Failure handling

Classify before action: `PRODUCT_FAIL`, `GATE_TEST_FAIL`, `CI_HARNESS_FAIL`, `INFRA_FAIL`, `EXPECTED_NOOP`, `SUPERSEDED`. Never weaken a gate to make a patch pass. A real post-certification escape becomes the next learning/corrective loop without rewriting prior release evidence.

## Permanent owners/boundaries

Smart Provider Router v2 sole routing authority; canonical freshness/recovery sole freshness owner; existing multi-feed allocator sole subscription owner; BroadSnapshotBroker canonical reuse owner; canonical persistence/cache owners; canonical U.S. market calendar/session owner; direct SEC/EDGAR authoritative; existing telemetry/symbol/state owners; deterministic Day/Swing/Long truth; U.S. Equities Processing; GLD/SLV/USO actionable exceptions; No Execution.

## Exactly one next action

Run #64 / v18.9.1 G0 from complete macOS crash evidence or deterministic reproduction and freeze the narrow G1 before any product-source change. Do not start v18.9.2 or any v19 implementation branch until current release ordering permits it.
