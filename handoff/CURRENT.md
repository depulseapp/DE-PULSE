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
**Active development branch:** none  
**Master corrective program:** #65 / `ADAPT-PROVIDER-INTELLIGENCE-010`  
**Immediate blocker/next patch:** #64 / `ADAPT-RUNTIME-CRASH-001` -> `v18.9.1`.

## v18.9.0 — COMPLETE / IMMUTABLE STABLE

Fast #481 / `32525637987`, Qualified #153 / `32525738828`, merged candidate `9ea81cddae4875ae15d3719ca028519a36c597b6`, and Release #32 / `32526121817` are the release authority. G11-G16, actual macOS Apple Silicon + Windows x64 packaged-runtime audits, G15 assurance, same-run no-rebuild publication and G16 evidence passed. Stable tag: `v18.9.0-stable`. Durable manifest: `release/v18.9.0/stable-evidence-manifest.json`.

## Post-Stable audit direction — issue #65

The architectural foundation remains valid: Smart Provider Router v2 is the sole executable router; canonical freshness/cache/telemetry/persistence/state owners remain authoritative; direct SEC/EDGAR remains authoritative; no TradeInsight-specific parallel system is allowed.

The audit identified corrective/product gaps: TradeInsight Settings/API-key UX, first-success vs coverage-aware fulfillment, canonical company-name identity/presentation, behavior-oriented Market Data Modes, Form 4 enrichment, ticker/company search, movers/ranking evidence, full useful-capability admission and stronger provider-efficiency/adaptive telemetry.

The user explicitly prefers **many small complete builds over heavy builds** to reduce implementation misses. This is now a permanent process rule: one primary responsibility per patch, explicit non-goals, implementation-miss audit before moving on, and no known miss left only in chat.

## Ordered v18.9.x patch train

1. `v18.9.1` — #64 runtime crash corrective ONLY.
2. `v18.9.2` — TradeInsight Settings/API-key UX ONLY.
3. `v18.9.3` — coverage-aware Smart Provider Router v2 + persistence-first residual-gap fulfillment ONLY.
4. `v18.9.4` — canonical company identity + all-desk presentation ONLY.
5. `v18.9.5` — Market Data Modes + capability diagnostics ONLY.
6. `v18.9.6` — TradeInsight SEC Form 4 enrichment ONLY.
7. `v18.9.7` — TradeInsight ticker/company search ONLY.
8. `v18.9.8` — TradeInsight movers/ranking evidence ONLY.
9. `v18.9.9` — remaining useful TradeInsight capability sweep ONLY.
10. `v18.9.10` — provider efficiency + Adaptive Intelligence telemetry + protected-session headroom measurement ONLY.
11. `v18.9.11` — Session-Aware Data Readiness Maintenance ONLY: light overnight + heavy weekend, using canonical persistence/router/session owners and strict protection for pre-market/regular-market/after-hours.
12. `v18.9.12` — whole v18.9.x professional closure audit ONLY; no new feature scope.

A later patch may be split further at G0/G1 if it is still too broad. Never merge patches merely to reduce version count.

## Permanent adaptive provider + persistence contract

`consumer requirement -> in-memory canonical cache -> persisted canonical DB/state -> validate freshness/coverage/schema/provenance/rights -> exact residual gap -> eligible-provider ranking -> targeted acquisition -> canonical merge/provenance -> coverage re-evaluation -> persist -> next provider only if still needed -> synthesized consumer state`

Provider success does not equal consumer completeness. No fixed provider chain is authoritative. Static ordering is at most a prior/tiebreaker. Validation lifecycle (`SHADOW/VALIDATED/...`) is separate from serving role (`PRIMARY/FALLBACK/BACKFILL/ENRICH/CORROBORATE/...`).

Never refetch/recompute trustworthy evidence already valid for the consumer solely because a provider is available. Revision-prone evidence preserves point-in-time/as-observed history plus later revisions. Live-sensitive values obey freshness TTLs and cannot be presented as current merely because they exist in the DB.

## Permanent session-aware maintenance contract

Pre-market, regular market and after-hours are protected Tier-0 decision-support sessions. They always receive first claim on provider quota/headroom, network, CPU, memory, DB and worker capacity.

- **Light overnight maintenance:** small, high-value, gap-driven readiness work only after protected after-hours and before the next protected pre-market window.
- **Heavy weekend/extended market-closed maintenance:** deeper but bounded backfill/reconciliation/index/retention/outcome work.
- Maintenance uses only bounded surplus capacity after protected-session reserves.
- External-provider maintenance acquisition suspends during protected sessions unless directly required by a current/live consumer.
- Maintenance must drain/preempt/checkpoint/resume around protected sessions or market shocks.
- Missed work catches up only in a later eligible overnight/weekend window, not during live sessions.
- No blind full-universe refetch and no parallel maintenance calendar/scheduler/router/cache/database owner.

Machine contract: `adaptive-governance/PERSISTENCE_REUSE_AND_OFF_HOURS_DATA_READINESS_CONTRACT.md`.

## Hosted / multi-device / account architecture — APPROVED DIRECTION

DE.PULSE is intentionally being kept compatible with a future single-account experience across macOS, Windows and hosted web without replacing the current native architecture.

- `SQLite` remains the fast local edge/offline store and warm working set for native clients.
- Future hosted `PostgreSQL` becomes the shared authority for sync-eligible account/device state and lawful hosted evidence.
- Synchronization is application-level, incremental, typed, idempotent and checkpointed using stable IDs/outbox/change events — never blind raw database replication or dual-master table sync.
- Native clients continue to work from SQLite during hosted/network outages and reconcile later from durable checkpoints.
- Watchlists, desk membership, preferences and other account state use explicit versions/optimistic concurrency and deterministic conflict handling.
- Provider secrets do not sync through ordinary SQLite/PostgreSQL data tables; they remain under the canonical secret/security owner.
- Current/live market truth still comes from canonical freshness/provider/state owners; neither SQLite nor PostgreSQL becomes a competing market-truth subsystem.
- v18.9.x must **not** introduce PostgreSQL; it must preserve clean repository/persistence boundaries, stable IDs, provenance and sync compatibility. Hosted sync architecture belongs to v19+ in small governed packets.
- Cross-platform session validity is server-side canonical truth for macOS, Windows and hosted web; clients consume it rather than inventing their own authentication/session truth.

Governing contracts:
- `adaptive-governance/SQLITE_POSTGRES_SYNC_CONTRACT.md`
- `adaptive-governance/ROLE_AWARE_SESSION_SECURITY_CONTRACT.md`
- `adaptive-governance/ROLE_AWARE_UI_COMPOSITION_CONTRACT.md`

## v19 / v20 continuity

v19 Professional Data Infrastructure inherits persistence-first reuse, session-aware maintenance, provider/data-rights governance and the approved SQLite/PostgreSQL hosted-sync contract. It measures provider/data rights, quality, cost, coverage, reliability, revision correctness, DB/index/pool/capacity behavior, calls avoided, maintenance value and protected-session reserve sizing. It must not recreate the router, canonical persistence owner or maintenance coordinator.

v20 Adaptive Intelligence consumes provenance-bound point-in-time evidence/outcomes and may learn provider/maintenance usefulness only through governed SHADOW/Champion-Challenger promotion. It may not reduce protected live-session safety, bypass rights/provenance or sacrifice current truth for background learning.

## Other continuity truth

Issue #57 / v18.8.1 Market Intelligence escape is now closed completed because `release/v18.8.2/stable-evidence-manifest.json` already records `RESOLVED_IN_V18.8.2_STABLE`. Its regression remains mandatory in final v18.9.x closure.

The real v18.9.0 macOS Apple Silicon `EXC_CRASH (SIGABRT)` remains unresolved under #64. Do not guess root cause or delete `PersonalMarketTerminal` state/API keys as a first troubleshooting step.

## Exactly one next action

Diagnose #64 using the complete macOS crash evidence or deterministic reproduction and freeze the narrow `v18.9.1` G1. Do not create a `v18.9.2` branch until `v18.9.1` is truthfully closed or the crash is proven external/non-product.

## Resume rule

Any ChatGPT account, Codex session, Claude or human maintainer must first fetch the **current live GitHub head** because another session/process may have advanced it. Then read `AGENTS.md`, `CLAUDE.md`, `governance/AI-ASSISTANT-PORTABILITY-CONTRACT.md`, this file, all four `adaptive-governance/CURRENT_ADAPTIVE_*` files, `adaptive-governance/PERSISTENCE_REUSE_AND_OFF_HOURS_DATA_READINESS_CONTRACT.md`, `adaptive-governance/SQLITE_POSTGRES_SYNC_CONTRACT.md`, `adaptive-governance/ROLE_AWARE_SESSION_SECURITY_CONTRACT.md`, `adaptive-governance/ROLE_AWARE_UI_COMPOSITION_CONTRACT.md`, `release_identity.json`, `release/v18.9.0/stable-evidence-manifest.json`, both `.depulse-certification/resume/` checkpoints, issue #65, issue #64 and their current comments before changing code. Inspect commits since the certified Stable/source baseline so completed work is not duplicated. GitHub objects and executable evidence outrank chat memory.
