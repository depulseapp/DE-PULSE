# DE.PULSE — Adaptive Delivery Process

**Status:** PERMANENT DELIVERY CONTRACT / AUDIT-REBASELINED  
**Authority:** `governance/ADAPTIVE-OPERATING-CONTRACT.md` + immutable candidate evidence  
**Live state:** `governance/current-state.json` + active closure ledger + `handoff/CURRENT.md`

Delivery flow:

`qualified source -> immutable RC -> native/Web build -> provenance/SBOM -> actual artifact/runtime audit -> signing/notarization -> release assurance -> authorized promotion -> health checks -> G16 handoff`

## 1. Delivery truth

A source test does not prove a package, installer, update, hosted deployment or user workflow. Delivery claims bind to the exact source fingerprint, toolchain, artifact SHA, configuration, migration state, platform and test evidence.

Historical Stable tags, artifacts and provenance are immutable. A failed TEST/RC cannot overwrite or destabilize the known-good Stable installation. Publication uses certified artifacts without rebuilding them.

## 2. Version-oriented delivery

- one coherent development branch/PR per version unless explicitly superseded;
- exact-head Fast for the coherent development candidate;
- impact-selected exact-head Qualified at material risk boundaries and G10;
- Release workflow only for an actual G11–G16 candidate;
- related small changes may share a version; heavy/high-risk behavior gets a real version split;
- commits, issue rows, audit findings and test cases remain internal traceability—not pseudo-releases;
- obsolete runs may be cancelled/superseded, but evidence from a different head is never reused as current.

The canonical workflow families remain Fast, Qualified and Release. Fast never dispatches Release.

## 3. Audit closure at delivery

A version is not delivery-complete until every applicable audit finding/risk, 5/5 row, surface disposition, ADR and compatibility migration is reconciled against real evidence.

Where applicable prove:

- server-owned deterministic policy and no client truth fork;
- versioned symbol/evidence/event/delta contracts;
- one Opportunity identity/lifecycle across detectors, Discovery, Watchlist, Alerts, Research and Desks;
- Watchlist as selected-universe projection with no duplicate scanner/scorer/provider loop;
- frozen-as-of Decision Brief lineage from the transition that caused attention;
- point-in-time/vintage/revision/raw-adjusted/censoring truth;
- adaptive selection-bias controls and deterministic fallback;
- tenant persistence/outbox/isolation/conflicts/tombstones/retention and real recovery;
- managed/OS secrets, audit, signed updates and package trust;
- old-client compatibility, deprecation and forced-upgrade behavior;
- offline/last-known cache truth;
- alert causal dedupe, delivery, quiet-hours and acknowledgement;
- provider correlation/revision/rights/freshness truth;
- SLO, runbook, recovery and operational ownership.

Documentation or a feature happy path cannot satisfy these requirements.

## 4. Compatibility-first delivery

Authority migration is delivered only when evidence proves:

`old behavior characterized -> new owner present -> dual/shadow comparison -> consumers migrated -> parity/approved improvement -> old authority retired -> rollback/regression retained`

This applies to renderer calculations, `RuntimeSnapshot`, Radar/Discovery state, Watchlists/workspaces, full-snapshot SSE, provider direct paths, identity/session and persistence schemas.

Golden vectors must distinguish known defects from intended behavior. Existing client/capability loss is prohibited unless an explicit decision approves removal and migration/rollback is clear.

## 5. Provider and Data Health delivery

Provider delivery conserves #80/#81/#82/#83/#78/#84 and the canonical provider matrix/SLO/fetch-path contracts.

Prove applicable:

- Smart Provider Router v2 remains sole general routing/admission owner;
- SEC/EDGAR and other direct authorities remain protected;
- Registry adapters self-register without duplicating Router/health/cache/persistence/subscription/lifecycle/state;
- capability, entitlement, freshness, history, quota, cost and serving reason are truthful;
- Settings/Maintenance/Data Health agree on configured/eligible/serving/fallback/recovery states;
- delayed data is not represented as live;
- outage/auth/rate-limit/stale/malformed/downgrade/recovery behavior is scoped and non-flapping;
- provider agreement does not overstate independent corroboration;
- corrected facts preserve revision/supersedes lineage;
- rights coordinates remain distinct from credentials and technical reachability;
- required consumers receive canonical routed state rather than page-specific fetches;
- secrets never leak into clients, logs, telemetry, fixtures, artifacts or GitHub evidence.

A provider Test button proves connectivity only, never implementation or Commercial/Public rights.

## 6. Adaptive and AI delivery

Learned production influence requires point-in-time outcomes, registered feature/model/policy versions, time-split evaluation, sample/uncertainty gates, shadow/champion-challenger comparison, drift/subgroup checks, explicit approval, bounded rollout and rollback thresholds.

AI/LLM/agent delivery proves rights filtering, evidence-ID grounding, contradiction/uncertainty preservation, prompt-injection resistance, source isolation, audit and deterministic non-AI operation. LLM output cannot overwrite canonical prices, indicators, routing, rights, lifecycle or trade-plan geometry.

SHADOW never silently becomes production.

## 7. Database, privacy and migration delivery

Applicable delivery evidence includes:

- expand/contract migration and backward/forward compatibility;
- tenant-key/isolation/RLS disposition and cross-tenant denial;
- transactional outbox/retry/idempotency and conflict/tombstone semantics;
- indexed/bounded query and retention behavior;
- restart, rollback and partial-outage reconciliation;
- real managed backup, PITR, restore/failover and measured RPO/RTO when claimed;
- privacy export/deactivation/deletion and backup-restore anti-resurrection;
- security/audit retention separated from personal-data retention;
- migration dry run against representative data and downgrade/old-client behavior.

Schema success on an empty local SQLite database is insufficient for hosted or upgrade claims.

## 8. Cross-platform delivery

Shared capability is complete only when every REQUIRED client is equivalent in domain meaning, authentication/roles, rights, freshness/provenance, intelligence conclusions and durable state.

### macOS Apple Silicon

- native artifact built from the certified source;
- Developer ID signing, hardened runtime, notarization and stapling when public distribution is targeted;
- installer/upgrade/rollback and local-data migration tested;
- Keychain-backed credentials and redacted crash diagnostics;
- actual packaged launch and critical workflows audited.

### Windows x64

- native artifact built from the certified source;
- trusted Authenticode signing and installer/update strategy when public distribution is targeted;
- install/upgrade/uninstall/rollback and data migration tested;
- Windows credential-vault storage and redacted diagnostics;
- actual packaged launch and critical workflows audited.

### Web

- deployed versioned API/event schema, OIDC/session/CSRF/reauth and tenant isolation tested;
- CSP, secure cookies/tokens, cache/privacy behavior and browser compatibility verified;
- streaming reconnect/backpressure/stale-state behavior verified;
- accessibility and critical user workflows tested against the deployed environment;
- no provider secret reaches browser state.

Desktop may keep encrypted authorized last-known data but must label offline/stale truth. Platform-specific mechanics may differ; domain decisions may not.

## 9. G13–G16 delivery gates

### G13 — Native Packaging & Provenance

Build all required artifacts from frozen source. Record source fingerprint, dependency lock/toolchain, SBOM/provenance, artifact identity and reproducibility inputs.

### G14 — Actual Artifact Runtime Audit

Install/launch the real artifact or deploy the real candidate. Verify identity, configuration, auth, migration/persistence/restart, provider degradation/recovery, core intelligence, accessibility and No Execution. Never fabricate unavailable native/runtime evidence.

### G15 — Release Assurance & Promotion

Confirm exact-artifact checks, signing/notarization, rollback, schema safeguards, release notes, support/runbooks, security/rights status and promotion authorization. Generate final hashes/manifests last. Commercial/Public activation remains separate.

### G16 — Adaptive Retrospective & Handoff

Aggregate build AIPLC, CI failures, defects, platform/runtime evidence, provider utility/cost, data freshness, alert/intelligence outcomes, UX observations and operational incidents. Convert learning into canonical prevention, named future work or an evidence-backed rejection. Reconcile machine state and leave exactly one next action.

## 10. Commercial/Public boundary

Development Production Ready may prove technical security, persistence, provider-rights enforcement, cross-platform behavior, signing/update capability and operations. It does not grant provider redistribution rights, approve public-user legal/compliance, enable subscriptions or authorize commercial activation.

Commercial/Public Ready requires external provider-specific rights evidence, public privacy/legal/support/compliance, commercial controls and an explicit Owner activation decision. Public-serving paths fail closed where required rights are absent or expired.

## 11. Scale and operations acceptance

Delivery evidence must match the promised scale:

- around 100 users: lawful serving, user-scoped deltas, managed secrets, signed beta, privacy/support runbooks;
- around 1,000 users: shared ingestion, tenant-normalized DB, outbox/workers, quotas/cost budgets and durable alert fanout;
- around 10,000 users: measured partition/retention/queue/DB-hotspot behavior, SLO/error budgets, unbiased adaptive sampling and sustainable on-call ownership.

Do not adopt Kafka, Kubernetes or microservices merely because a future count is large. Postgres + outbox + stateless API/workers remains the default until measurements justify more.

## 12. Post-release verification and rollback

After authorized promotion:

- verify published hashes/signatures and release identity;
- perform bounded install/deploy smoke and health checks;
- monitor crash, stale-data, provider, DB, queue, alert and migration health;
- stop rollout or roll back on predefined thresholds;
- preserve failed-candidate evidence without contaminating Stable;
- record incidents and prevention through AIPLC/G16.

## 13. Delivery exit criteria

Delivery closes only when source, tests, actual artifacts/deployment, migrations, security, data truth, rights state, platforms, observability, rollback and user workflow evidence agree on the same candidate. Maturity labels remain evidence-backed, and `handoff/CURRENT.md` plus machine state identify the exact next action.
