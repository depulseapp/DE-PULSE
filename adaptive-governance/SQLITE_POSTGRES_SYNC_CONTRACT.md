# DE.PULSE — SQLite ↔ PostgreSQL Sync Contract

**Status:** APPROVED / FUTURE HOSTED-DATA CONTRACT  
**Applies to:** v19+ hosted/multi-device evolution; v18.9.x must remain compatible  
**Primary owners:** existing canonical persistence/repository owners + future hosted sync service  
**No blind database replication, dual-master raw-table sync, or provider-specific storage silo is permitted.**

## 1. Purpose

DE.PULSE currently relies on local SQLite for native persistence. When hosted/shared PostgreSQL is introduced, SQLite and PostgreSQL must cooperate through one canonical logical data model so the native app remains fast/offline-capable while hosted/multi-device state becomes durable and shareable.

The goal is **controlled incremental synchronization**, not copying whole databases back and forth.

## 2. Authority model

- **SQLite = local edge/offline store and warm working set.** It supports low-latency local reads, restart recovery, offline continuity, and bounded local writes.
- **PostgreSQL = shared hosted authority for sync-eligible server/shared data once hosted mode exists.** It is the canonical shared copy for multi-device/account-scoped state and hosted evidence that is lawful to persist.
- Neither database is allowed to become an independent competing truth owner.
- Current/live market truth still comes from canonical freshness/provider/state owners; a database row is not current merely because it exists.

## 3. Sync mechanism

Use application-level change synchronization, not raw SQLite/Postgres replication:

1. transactionally write the local/server domain change;
2. record an immutable/idempotent change event or outbox entry with stable record ID, data class, version/revision, timestamp, provenance and schema version;
3. push unsynced eligible changes in bounded batches;
4. server validates rights/schema/identity/version and applies idempotently;
5. client pulls changes since its last durable checkpoint/watermark;
6. apply deterministically to SQLite;
7. advance sync checkpoint only after transactional success;
8. retry safely after restart/network loss without duplicate records.

## 4. Data-class sync policy

### Append-only / evidence history
SEC filings, Form 4 evidence, congressional disclosures, 13F snapshots, historical bars, point-in-time research/opportunity/thesis/outcome evidence and other immutable/as-observed records should use stable IDs + append/revision semantics. Conflicts are normally additive rather than destructive.

### Revision-prone reference/history
Corporate actions, fundamentals, macro series, amendments/restatements and other corrected data retain original as-observed versions plus later revision lineage. Sync must not silently overwrite point-in-time history.

### User/account state
Watchlists, desk membership, user preferences and other sync-eligible user state use explicit record versions/optimistic concurrency and deterministic conflict handling. Do not use unqualified last-write-wins where it can silently discard meaningful edits.

### Derived intelligence
Persist/sync only when useful. Include input fingerprint, algorithm/model/prompt version where applicable, generated-at time, provenance and invalidation rules so stale derivations are not mistaken for current truth.

### Live/high-frequency data
Do **not** sync every tick/quote/spread merely because it exists locally. Persist/sync selected snapshots only when a named historical, research, recovery or audit consumer requires them.

### Secrets
API keys/tokens/passwords never sync through ordinary SQLite/Postgres market-data tables. Use the canonical secret owner and a separately governed hosted-secret design if hosted credentials are introduced.

## 5. Persistence-first request behavior

For native clients:

`consumer -> memory -> SQLite -> validate coverage/freshness -> residual gap -> provider/hosted source -> merge -> SQLite -> async eligible sync`

For hosted services:

`consumer -> memory/cache -> PostgreSQL -> validate coverage/freshness -> residual gap -> Smart Provider Router -> merge -> PostgreSQL -> distribute eligible deltas`

The presence of PostgreSQL must reduce duplicate external provider requests, not cause each device/service to independently refetch the same immutable history.

## 6. Provider/data-rights boundary

Every synchronized data class must be allowed by provider/storage/redistribution/commercial/AI-use rights. If a provider allows local cache but not hosted persistence or redistribution, the record remains local-only or stores only derived/metadata forms permitted by the contract. v19 provider-rights registry is authoritative for this decision.

## 7. Session-aware sync priority

Pre-market, regular market and after-hours remain Tier-0 protected sessions.

- Foreground/live fulfillment outranks SQLite↔Postgres synchronization.
- Sync uses bounded background bandwidth, DB connections, CPU and worker capacity.
- Large upload/download reconciliation is deferred to overnight/weekend Data Readiness windows.
- Small critical account-state changes may sync promptly when cheap, but must not starve market-data work.
- Sync must checkpoint and resume after network loss, shutdown or preemption.

## 8. Schema and migration discipline

- Maintain one canonical logical schema/domain contract with dialect-specific migrations where required.
- Every persisted/synced record carries a compatible schema/version contract.
- SQLite and PostgreSQL migrations must be independently testable and forward/backward compatibility explicitly governed during rolling upgrades.
- No feature may require dropping/recreating the user's SQLite database as a normal upgrade path.

## 9. Integrity and conflict requirements

Required protections include:

- globally stable IDs/idempotency keys;
- revision/version numbers or monotonic sequence where appropriate;
- created-at / observed-at / effective-at / updated-at distinction;
- source/provenance and input fingerprint;
- tombstone/delete semantics where deletion is allowed;
- per-device/server sync checkpoint;
- duplicate detection;
- deterministic conflict resolution by data class;
- checksums/integrity validation for material batches;
- audit trail for rejected/conflicted changes;
- safe replay after crash/restart.

## 10. Availability behavior

- The native app must continue operating from SQLite when PostgreSQL/network/hosted service is unavailable, subject to normal freshness truth.
- Hosted outage cannot erase local data or force a reset.
- When connectivity returns, bounded catch-up resumes from the last durable checkpoint.
- A sync backlog may reduce cloud freshness but must not falsely mark stale shared data as current.

## 11. Patch / roadmap placement

### v18.9.x
Do not introduce PostgreSQL. Preserve clean canonical repository/persistence boundaries and stable IDs/provenance needed for later sync. v18.9.3 persistence-first routing and v18.9.11 session-aware maintenance must remain sync-compatible.

### v19 — Professional Data Infrastructure
Add a dedicated small packet for **SQLite/PostgreSQL sync architecture + hosted persistence readiness** after data-rights/schema/quality contracts are mature enough. Scope includes domain IDs, outbox/change log, checkpoints, conflict policies, migration parity, sync eligibility by data class, offline/restart proof, and protected-session resource budgets. If hosted Postgres is not yet deployed, implement/test the abstraction and contract without forcing cloud infrastructure prematurely.

### v20
Adaptive intelligence may consume synchronized point-in-time evidence only when provenance, versioning, rights and lineage are intact. Sync must not alter historical observations or permit model outputs to overwrite source truth.

## 12. Permanent principle

**SQLite is the fast local edge; PostgreSQL is the shared hosted authority; synchronization is incremental, typed, idempotent, provenance-bound and session-aware — never blind dual-master replication.**
