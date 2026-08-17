# DE.PULSE — Developer documentation

## v18.4.0 TEST — Security / commercial-readiness architecture

v18.4 extends the existing HTTP security/auth boundary and runtime telemetry instead of creating parallel owners. Fresh re-authentication protects only high-impact mutations. Hosted mutation/expensive-work quotas are keyed from verified identity after canonical auth/CSRF resolution, use bounded in-memory windows, return deterministic HTTP 429/`Retry-After`, and expose aggregate counters without raw user/session/IP identifiers.

Provider `DataRights` / `CommercialReadiness` metadata is governance-only and fails closed unless evidence plus explicit commercial/redistribution/AI-use approval are bound. Structural tests prevent those fields from entering Smart Router eligibility/scoring or executable routing. v18.3 PostgreSQL/shared-state owners and protected deterministic formulas are unchanged.

## v18.3.0 STABLE — PostgreSQL / hosted shared-state architecture

Stable promotion switches the canonical desktop runtime/config to `PersonalMarketTerminal` while hosted mode continues to require explicit PostgreSQL selection. The v18.3 TEST profile remains historical isolation only; no persistence owner, provider pipeline, or deterministic scoring logic changes during promotion.

`PersistenceBackend` remains the single repository contract. Local desktop selection delegates to the existing SQLite/file fallback implementations; an explicit `postgres`/`postgresql` selection binds the hosted target to a PostgreSQL implementation under the `postgres` build tag. Requested PostgreSQL is fail-closed: a missing driver build or DSN produces an unavailable backend rather than local fallback.

PostgreSQL mirrors SQLite migrations 1–4 for Global Symbol Registry, canonical quotes/history, evidence, decision lineage, outcomes, derived features, `IdentityPersistentState`, and `UserWorkspace`. It uses bounded `database/sql` pooling, transactional upserts, a migration advisory lock, query/pool diagnostics, 5-minute quote-history buckets and the existing 30-day history retention. Warm-start quotes are always relabelled `persisted`, never timeless LIVE truth.

Hosted runtime selection is explicit and bypasses desktop instance-focus/native-window behavior while preserving the existing desktop path by default. `/api/health` is liveness; `/api/ready` requires canonical persistence plus IdentityService readiness. `processingStateLocked()` remains the one shared union owner, so hosted persistence cannot create per-user provider/scanner/scoring pipelines.

## v18.2.0 STABLE — Admin / presence / session architecture

Stable promotion switches the canonical release identity to `STABLE` and runtime/config to `PersonalMarketTerminal`; the v18.2 TEST profile remains historical isolation logic only.

v18.2 extends the existing `IdentityService`; it does not introduce a second user/session/presence store. `AdminUserView` / `AdminSessionView` expose only operational fields, while password hashes, token hashes and opaque tokens remain server-side. Role hierarchy and active-owner safety are centralized in the identity owner. Role/status/password changes revoke affected sessions.

Presence is derived from persisted `SessionRecord.LastSeenAt`, idle/absolute expiry and revocation. SSE keepalive revalidates the same canonical session and touches `LastSeenAt`; revoked or expired streams terminate. The modular Settings admin UI is role-gated to SUPER_OWNER / OWNER / ADMIN and uses existing CSRF protection for mutations. v18.1 `UserWorkspace` remains the personal-market-state owner, and shared provider/evidence/Router/Rapid Move/deterministic scoring pipelines remain unchanged. PostgreSQL/hosted web stays v18.3.

## v18.1.0 STABLE — Multi-user ownership architecture

`UserWorkspace` is the single durable owner for personal watchlists and UI state. It is persisted by the existing `PersistenceBackend` abstraction and keyed by immutable authenticated `UserID`; the pre-v18.1 owner state migrates once, then shared `state.json` retains operational settings rather than a private copy of the owner's symbols. Auth materializes an empty workspace for a new account before normal request handling.

The engine remains one canonical backend intelligence core. `processingStateLocked()` forms a deduplicated union of user symbols for provider/scoring work; it does not create per-user Router, quote, Radar, Rapid Move or deterministic scoring pipelines. `SnapshotForUser` and targeted Hub fanout are presentation/privacy filters only. Shared provider/API configuration and runtime policy remain ADMIN-controlled. v18.2 owns admin/presence/session lifecycle UI. The **No Execution** boundary and protected formulas are unchanged.

## v18.0.6 STABLE — Smart Provider Router v2 + Rapid Move hardening

Build `v18.0.6-test-smart-provider-router-rapid-move-market-shock-hardening-20260814` reuses the canonical provider router, reconciliation owner, Rapid Move state machine, Radar promotion path and durable outcome memory. It adds canonical source-disagreement scorecard truth, MARKET_SHOCK classification, explicit hysteresis, provider-time outcome-anchor preservation, and governed SHADOW learning (`SHADOW → VALIDATED → APPROVED → PRODUCTION`) with no automatic promotion and no protected formula impact.

## v18.0.5 STABLE — UI/UX + Symbol Management Hardening

v18.0.5 Stable is based on **v18.0.4 STABLE** with v17.5.1 retained as the v18 major-family provenance anchor. Cross-desk tracked-symbol mutations are consolidated behind one canonical owner while existing durable watchlist/symbol-registry persistence is reused; explicit empty desk lists must not rehydrate defaults after Remove All. Role-aware surface suppression keeps provider/runtime machinery out of USER/DEMO without deleting privileged diagnostics. Opportunity Radar and Research Target changes reuse existing data pipelines; no new dataset/provider/model/background task is introduced. The Provider Router, protected deterministic formulas and **No Execution** boundary remain unchanged.

## v18.0.4 STABLE — Windows persistence lifecycle / native closure

Real-Application tests register persistence cleanup before `t.TempDir` teardown so Windows SQLite handles are closed deterministically. Native G14 verifies canonical `/api/health` `buildId`, exact source SHA/fingerprint, native SQLite/package execution, Stable migration/isolation, bootstrap OWNER, CSRF/password setup, restart login, Smart Router v2, Rapid Move, and Coverage Truth.

## v18.0.3 TEST — Native cross-platform runtime/test portability

The runtime embedded-FS path owner now uses URL/embedded-FS semantics (`path.Clean` / `path.Ext`) instead of host-filesystem semantics (`filepath.Clean`), preventing Windows backslashes from producing embedded-resource 404s. Real-Application tests isolate HOME/XDG_CONFIG_HOME/APPDATA and derive the effective base with `os.UserConfigDir()`. Windows `os.FileMode` is not treated as an ACL source; native G14 owns Windows security/isolation evidence. No Smart Router v2, Rapid Move, provider, scoring, or deterministic Day/Swing/Long behavior changes.

## v18.0.2 TEST — Canonical source fingerprint

The canonical source fingerprint owner is `source_fingerprint.py`. Relative paths are sorted and hashed with platform-neutral POSIX separators, and CI, pre-freeze qualification, and full certification reuse that owner. G14 must reject source SHA or canonical fingerprint drift on every native target. Runtime Smart Router v2 / Rapid Move behavior is inherited from v18.0.1 unchanged.


## v18.0.1 TEST — Smart Router v2 + Rapid Move foundation

Build `v18.0.1-test-smart-router-v2-rapid-move-foundation-20260813` extends the **existing canonical Provider Router** rather than introducing a parallel owner. Smart Router v2 adds versioned provider×dataset×instrument×session capability state, deterministic route scoring, per-capability circuits, Preferred-vs-Serving reasons, persistent NOT_ENTITLED suppression/cooldown, p50/p95 latency and calls-avoided telemetry. Legacy route-chain contracts remain compatible while the executable owner ranks eligible providers dynamically.

`rapid_move_intelligence.go` consumes the canonical quote/history/provider-observation stores; it does not run a duplicate scanner. Bounded 15s/30s/60s/2m/5m windows are computed from in-memory history, then correlated with source agreement, spread/liquidity, corporate actions, catalyst/news/SEC/earnings and SPY/QQQ context. Material events reuse Opportunity Radar promotions and `livePriorityHints`. Point-in-time event/provider/receipt/detection/outcome timestamps, a trace ID, policy version and 1m/5m/20m outcomes feed the existing Evidence → Decision Lineage → Outcome persistence path. Only material/learning evidence is persisted; raw ticks are not warehoused.

Production detection remains deterministic and bounded. Adaptive improvement is SHADOW-first under **SHADOW → VALIDATED → APPROVED → PRODUCTION**; no learned policy silently self-modifies production. Data Utility registration covers provider capability state, Rapid Move events and outcome-learning evidence.

## v18.0.0 TEST — Identity & Secure Session Foundation

v18.0 extends the existing v17 repository/application boundary rather than adding a parallel authentication stack. `Application.auth(...)` resolves a canonical `Principal`; roles are `SUPER_OWNER`, `OWNER`, `ADMIN`, `USER`, and `DEMO`. Shared/system-mutating routes are centrally role-gated while user research/desk operations remain authenticated.

Local credentials use **Argon2id** only. The credential verifier bounds PHC parameters before allocation, session cookies contain opaque random tokens only, persistence stores token hashes rather than raw session tokens, credential changes revoke prior sessions, and idle expiry never extends beyond absolute expiry. State-changing authenticated requests use same-origin double-submit CSRF validation. Password verification runs outside the identity mutex and account state is revalidated before a session is issued.

Persistence schema v3 adds durable identity/session state while retaining all v17 market tables. The v18 TEST runtime is `PersonalMarketTerminal-v18-TEST`: first-run migration clones the compatible Stable profile transactionally through a temporary directory, excludes instance/log artifacts and refuses to clone while the detected Stable instance is live. This is deliberate side-by-side TEST isolation; v17.5.1 Stable remains untouched.

v18.0 does **not** introduce per-user watchlist ownership, hosted PostgreSQL activation, presence/admin operations, Smart Provider Router v2, TradeInsight, or any execution capability. Those remain later v18 workstreams. Deterministic market logic remains an inherited regression contract.


## v17.5.1 STABLE — Release identity & documentation hardening

This patch changes **release identity/documentation only**. Canonical runtime identity is `17.5.1 / STABLE / v17.5.1-stable-release-identity-documentation-hardening-20260813`. No Provider Router, persistence schema, workload policy, deterministic decision formula or executable product boundary is changed. G16 now protects Stable-promotion reconciliation so current docs cannot retain pre-promotion RC/TEST claims and the complete v17.0→v17.5 history cannot silently disappear.

## v17.5.0 Stable — Major Closure engineering result

The v17 family closed with the frozen 20-item scope mapped to current owners/evidence, inherited approved scope preserved, native macOS SQLite runtime acceptance, Principal Engineer review, Professional Trader/Investor review and GO v17→v18. Exact-RC-bit promotion remains historical provenance; v17.5.1 corrects the user-visible Stable identity without changing the v17 architecture.

## v17.4.0 — Completed UX + Operational Hardening checkpoint

Canonical renderer ownership was consolidated for Master Market Symbols. Preparation exception truth now distinguishes shared stale/freshness root causes from genuine per-symbol liquidity/event/extended risks, while retaining bounded groups and per-symbol drill-down.

## v17.3.0 — Completed Performance + Reliability Hardening checkpoint

Runtime SLOs consume canonical freshness, provider/workload telemetry and persistence diagnostics. CPU/memory/goroutine, interactive latency, provider/subscription utilization, warm-start, provider-call avoidance, DB write/storage growth and degradation-recovery evidence are owned by the runtime observability/SLO path.

## v17.2.0 — Completed Canonical Intelligence + Data Efficiency checkpoint

`canonical_pipeline.go`/persistence intelligence own material-change propagation, immutable evidence/decision lineage, separate outcomes and incremental derived features. Heavy downstream work is not recomputed on every quote tick. The Alpaca IEX allocation/self-deadlock correction is protected by regression/race evidence.

## v17.1.0 — Completed Runtime Load + Backpressure checkpoint

`workload_controller.go`, runtime priority and Provider Router integration own Tier 0–4 priority, bounded queues/concurrency, critical reserve, low-priority shedding and attributable load/capacity degradation. Provider-plan entitlements remain measured/configured rather than guessed constants.

## v17.0.0 — Completed Persistence Foundation checkpoint

`PersistenceBackend` provides the storage abstraction; SQLite migrations/indexed canonical tables, retained Global Symbol Registry, persisted canonical warm start, async/coalesced material writes, structured evidence/history and runtime/provider/persistence observability form the foundation. Freshness labels are re-derived at read time; persisted LIVE/CURRENT labels are never trusted as timeless facts.

## v16.11.0 Stable — Major Closure engineering contract

v16.11 is a no-feature-by-default **Major Closure & Release Assurance** build. It reconstructs the v16.1→v16.10 family from current source and immutable acceptance contracts, maps every capability to canonical code owners/fresh executable evidence, then runs full G0–G16 certification. Historical reports corroborate but do not substitute for fresh current-source proof.

New permanent controls include the **Major Closure & Release Assurance Contract** (RL-029), freshness/readiness cross-surface safety (RL-030), and source-package hygiene (RL-031). `release_identity.json` remains the canonical current-release owner. Provider Router remains the sole executable provider authority. `source_health_baseline.json` is now an active maintainability policy rather than stale historical metadata, and the principal-engineer review consumes the same policy.

The real-money closure fix is intentionally renderer/readiness-only: `tradeReadiness()` prevents stale/cached/history-only current-price evidence from presenting clean Day readiness. It does **not** change deterministic Day/Swing/Long Score/Action formulas or introduce a second decision engine.

CI/CD remains resource-class aware under RL-028: Go/build-cache-heavy jobs use a controlled lane; independent HTTP/professional/adversarial work may run bounded parallel lanes. Heavy inherited audits remain independent evidence owners under RL-027 and are not recursively nested in release scope gates.

No v17 database/repository, v18 authentication/multi-user, Portfolio, Journal, Paper Trading or execution scope is pulled into v16.11. v17 begins only after an explicit closure GO decision.

## v16.10.0 Stable — Opportunity Radar / Adaptive Data Policy / Shadow Foundation

v16.10 keeps **Discovery/Scanner as the only scanner owner**. `opportunity_radar.go` extends the existing canonical snapshot/scanner/live-priority architecture; it does not introduce another Provider Router, stream owner or deterministic engine. Broad U.S.-equity candidates are snapshot/batch qualified and at most the bounded promotion reserve is inserted through existing `livePriorityHints`. `adaptive_data_policy.go` tightens targeted history/cache cadence for hot promoted symbols and relaxes work when cold/provider-degraded. `shadow_experiments.json` and `ShadowControlState` are observational only: `CanMutateProduction=false`.

The mandatory Information Architecture/Data Efficiency audit is `v16_10_information_architecture_data_efficiency_audit.md`. Data Utility now explicitly consumes Market Activity and Opportunity Radar datasets. `release_identity.json` remains the canonical current-release owner; Provider Router authority, U.S.-equities boundary, original roadmap **30 FULL / 0 PARTIAL / 0 MISSING**, No Execution Boundary and deterministic formulas remain protected. No v17 database or v18 auth/multi-user architecture is pulled forward. RL-027 forbids recursively nesting expensive inherited audits inside release scope gates; unique evidence is independently checkpointed.

## v16.9.0 Stable — Final Original Roadmap / Evidence Fusion / Replay Closure

v16.9 closes only original items **#10 Community Intelligence, #11 Oil / Energy and #20 Historical Replay** and reaches **30 FULL / 0 PARTIAL / 0 MISSING**. Community is implemented as a canonical evidence-fusion extension over existing Context/Event/AI safety owners: source policy → normalization → U.S.-ticker resolution → dedupe/cluster → diversity/velocity → authoritative corroboration → materiality. `community_source_policy.json` governs access, retention and AI eligibility; untrusted community evidence cannot silently mutate deterministic decisions.

Oil/Energy extends the existing EIA/Context owner with WTI/Brent truth, trend/change, Brent–WTI spread and U.S.-market relevance while explicitly withholding unsupported CL/BZ continuous-contract/roll claims. Replay extends Professional Validation & Learning and reuses the production deterministic engine plus cutoff-filtered historical state; the reusable scenario catalog never fabricates missing evidence. **Provider Router** remains the sole executable provider authority, `release_identity.json` remains the canonical release identity owner, Adaptive Build Process v2/Pre-Freeze Qualification remains mandatory, and no v17 database or v18 auth/multi-user architecture is pulled forward.

## v16.8.1 Stable — US Scope / Data Utility / Build Process v2

`release_identity.json` is the canonical current-release metadata owner. Actionable user symbols are constrained to U.S.-listed equities/ETFs; the existing Provider Router/scanner/history owners are reused. Economic Calendar keeps U.S. events primary and selective foreign evidence as `GLOBAL_CONTEXT`; only `US_MARKET_CRITICAL` events activate full Event Mode. SPY/QQQ + selected ticker are Tier-0 live candidates. `data_utility_registry.json` and `data_health_policy.json` make downstream utility/freshness/cache obligations machine-checkable. Pre-freeze qualification is required before G12 heavy certification. No new canonical provider, deterministic engine or historical store is introduced.


## v16.8.0 Stable — Engineering / Governance Update

v16.8.0 closes #6/#8/#9/#21/#27 through canonical-owner REFACTOR/EXTEND. Provider Router remains the sole executable route authority. G1 freezes scope; G11 proves immutable clauses; G16 governs documentation. The original-professional-roadmap-acceptance.json contract remains immutable. RL-020/RL-021/RL-022/RL-023 remain ACTIVE; RL-022 was learned during G12 documentation-title validation. Adaptive Performance/Scalability, Testing and Intelligence contracts are inherited. `De-Pulse.app` and canonical Stable runtime continuity are preserved.


> **v16.7.0 Stable architecture:** this is a closure/refactor release, not a new engine generation. #3 extends existing Event Intelligence; #12/#13/#14/#15 extend existing Market Intelligence and Scanner owners. No new market-data provider path was added. `marketTradeabilityWithContext` consumes existing EventMode, Freshness, Liquidity, Scanner, Structure, Options and Global evidence; Breadth internals derive from existing canonical quotes/bars; scanner RS reuses already-fetched snapshots/history. Deterministic desk formulas remain frozen.

> **RL-020 / immutable roadmap acceptance:** `renderer/qa/original-professional-roadmap-acceptance.json` is now the clause-level original contract. G1 freezes it, G11 requires behavioral proof, and G16 reports FULL/PARTIAL/MISSING/DEFERRED/EXPLICITLY_AMENDED. An ID/name cannot silently replace its original acceptance depth.

> **RL-021 / G12 infrastructure classification:** `bounded_race_gate.py` distinguishes a `go test -list` enumeration timeout from an actual race/test failure, preserves the incident and retries enumeration once. A race PASS still requires actual `-race` execution.

> **RL-022 / G12/G16 documentation identity:** the canonical `# DE.PULSE — User documentation` and `# DE.PULSE — Developer documentation` headings remain the first Markdown line. Release-specific notes are nested below those titles.

> **Runtime identity:** standard app remains `De-Pulse.app`; Stable configuration remains `PersonalMarketTerminal`, preserving compatible prior Stable API keys/settings. G2 retains zero orphan production helpers and the REUSE → CONSOLIDATE → REFACTOR → DELETE/REPLACE → ADD rule.

## Core Developer Guide

> **v16.6.0 Stable — Full Professional Integration & Acceptance architecture:** no second engine/store/provider path is introduced. All 30 professional requirements are reconciled through their existing canonical owners. Confirmed defect `V16.6-MS-01` is fixed inside canonical desk-membership persistence: `nil` remains a legacy/missing sentinel, while non-nil empty slices remain intentionally empty across normalization, persistence and reload. `/api/master-symbol/remove-all` mutates the same DAY/SWING/LONG state; protected System Market Context remains in `requiredSymbolsFromState` without being misrepresented as user-managed membership.

> **Adaptive learning RL-017:** persisted collection semantics must preserve **missing vs explicitly empty**. G2/G11 own this protection through normalization contracts plus last-item/remove-all reload regressions. No G17+ gate was added.

> **v16.5.0 Stable — Context & Alternative Intelligence architecture:** `context_alternative_intelligence.go` derives one `ContextAlternativeIntelligenceSnapshot` from existing Market Intelligence, canonical quotes/options, EIA Macro and persisted user-authorized community evidence. No second Provider Router, quote store, options scheduler, historical store or deterministic engine was added.

> **GEX ownership:** `options_intelligence.go` extends the existing Alpaca options owner with snapshot gamma plus contract open-interest reconciliation through `fetchOptionOpenInterest`. GEX is withheld unless matched gamma/OI coverage and OI recency are defensible. The published value is explicitly a structural signed gamma×OI proxy; open interest cannot identify dealer long/short ownership.

> **Community ownership/safety — UNTRUSTED COMMUNITY INTELLIGENCE:** `community_intelligence_api.go` provides local-session authenticated manual add/remove only; there is no automatic community scraping. Stored text is fed to the existing v16.4 AI untrusted-content boundary with `untrustedExternalContent=true`. `ContextAlternativeIntelligenceSnapshot` is derived during the normal Engine Snapshot and has deterministic impact **NONE**.

> **v16.4.0 Stable — Professional Research AI architecture:** `ai_research_v2.go` derives one `AIResearchPackage` from existing canonical Evidence Snapshot, Research Package, Market Intelligence and Event Intelligence owners. `GenerateAI` was refactored in place; there is no second deterministic engine, Evidence Snapshot, provider store, or research truth path. Evidence IDs are deterministic and response citations are restricted to supplied IDs.

> **Safety boundary:** provider calls receive privileged system instructions separately from the user/evidence envelope. News, SEC/web/provider text is marked untrusted data; instruction-like content cannot become system/developer instructions. Response sanitization rejects unknown evidence IDs and trade/order-style next actions. Setup Score remains explicitly **not win probability**.

> **Routing/cache:** `resolveAIRouting` supports manual/efficient/balanced/deep policies using only configured providers with bounded output tokens. `aiCacheKey` uses the material Evidence Package identity, task/question and routing policy. Ordinary quote ticks do not invalidate the deep-analysis package; material Evidence Snapshot changes do.

> **Adaptive Release Learning Loop:** G0 loads `release_learning_registry.json` and proves ACTIVE prior lessons still have owning G0–G16 gates/permanent protections. G16 records new lessons, adapts the correct existing gate, and consolidates obsolete checks. The pipeline must not append G17+ merely to remember failures. Build-system changes follow **REUSE → CONSOLIDATE → REFACTOR → DELETE/REPLACE → ADD**. Randomized seeds are owned by the certification plan; the generic G12 runner is version-agnostic and accepts bounded positive reproducible seeds. G16 also validates the top-level README current-release identity so repository continuation truth cannot lag the in-app Documentation tabs.

> **Release/app identity:** active engineering states are Working Candidate / RC / Stable. Final macOS Stable bundle is **`De-Pulse.app`**. Canonical runtime/config remains `PersonalMarketTerminal` for prior Stable credential continuity. G16 is **Stable Release Stamp + Handoff + Learning Loop**.

> **Stable runtime/config continuity:** `NewApplication()` uses canonical `PersonalMarketTerminal`. This intentionally supersedes the former v16 TEST-only directory and preserves the v15.1.2-compatible `state.json` / `secrets.json` schema. Stable promotion tests load the prior Secrets schema and verify Finnhub, Alpaca, Groq, OpenRouter, Gemini, FRED, BLS, EIA, Twelve Data and Marketaux credentials unchanged.

> **v16.3.0 Stable — Professional Validation & Learning architecture:** `validation_learning.go` is the canonical derived analytics owner over the existing `SignalValidationState`, `Engine.bars`, corporate-action truth and candidate/watchlist/scanner context. `signal_validation.go` was refactored in place to freeze evidence ID, deterministic family inputs, settings fingerprint, Entry/Target/Invalidation levels and post-decision outcome truth. There is no second validation store, historical-bar store, corporate-action adjustment path or Evidence Snapshot model.

> **Replay invariant:** renderer replay calls the **same production deterministic scoring primitive** used by `computePlan`; v16.3 does not introduce a replay-only formula. Replay cutoff filtering uses completed-bar time, with future News/SEC/Macro/Earnings results withheld and current state/caches restored after every run. New frozen snapshots can be replayed `EXACT`; legacy rows are explicitly `SCENARIO/PARTIAL`. Post-snapshot split normalization reuses canonical Corporate Action truth.

> **Data/efficiency:** SPY/QQQ daily history is deepened to roughly ten years for defensible seasonality while remaining inside the existing Provider Router/canonical `Engine.bars` ownership. Ordinary symbols retain the normal history budget. Alpaca remains primary; Twelve Data/yfinance follow the routed fallback policy. Outcome resolution is persisted after canonical history commits so rolling intraday windows cannot erase an already established result. Deterministic Day/Swing/Long Score/Action formulas remain v14.3.7-compatible and unchanged.

> **v16.2.0 TEST — Professional Event Intelligence architecture:** `event_intelligence.go` is a derived read-only consolidation layer. It performs no provider fetches. Canonical input owners remain `research_providers.go` (News/Earnings), `macro_events.go` (official calendar/Event Mode), SEC/Catalyst Watch, and existing reaction stores. `Engine.Snapshot()` builds `EventIntelligenceSnapshot`, which is consumed by Market Intelligence, Research, Trade Readiness, Decision Queue, Maintenance and AI context. This preserves the permanent architecture `fetch once → canonical store → derived context → page-aware render`.

> **v16.2 inherited ownership:** #1 News Intelligence = source clustering/materiality/freshness; #2 Macro v2 = sourced Actual/Forecast/Previous/Surprise; #3 Economic Calendar = timing/impact/lifecycle/reaction rows; #4 FOMC/Fed = sourced Fed timeline; #21 Smart Notifications = stable event IDs + expiry, no generic Alerts workspace; #28 Reaction Intelligence = Catalyst/Event Mode reuse; #29 Event → Decision = context-only Readiness/Queue integration. This describes the inherited v16.2 layer; v16.3 consumes it as historical/validation context without duplicating its canonical owners.

> **Developer invariant:** Event Intelligence may change `Trade Readiness` and `Decision Queue` attention only. It must not mutate deterministic Day/Swing/Long Score/Action. Any future event provider adapter must enter through the canonical Provider Router/source owner rather than being called directly from `event_intelligence.go` or renderer code.


> **v16.0.6:** Final v16.0 consolidation/delivery closure. Carries all v16.0.5 professional-truth fixes, keeps the Stable-style primary header, restores the normal native app launch path with binaries prepackaged, and makes the complete release/handoff artifact set mandatory. No v16.1 Market Intelligence/Tradeability feature is included. Deterministic Day/Swing/Long Score/Action formulas remain frozen; v15.1.2 remains the protected Stable baseline.

> **v15.1.2 Stable:** Adds Research v2, canonical detailed Form 4 intelligence, Data Freshness v2/check-vs-data age semantics, session/event-aware refresh scheduling, stale auto-recovery, history cadence separation, trading-readiness checkpoint persistence/catch-up, and the permanent 48-item Approved Scope Completeness gate. Deterministic Day/Swing/Long score math remains protected.


> **v15.0.1 Stable:** Built on the frozen v14.3.7 baseline. Adds centralized provider routing/circuit breaking, dataset-aware freshness diagnostics, canonical desk-membership mutation APIs, Master Symbol global remove/restore, fallback-only yfinance, CBOE VIX validation, Marketaux fallback, richer SEC Form 4 transaction rendering, and v15 regression coverage. Deterministic trading formulas remain equivalent to v14.3.7.

> **v14.3.7 Stable:** Exact v14.3.6 baseline plus a rendered font-fit audit. Changes are CSS-only font sizing for verified overflow cases (Decision Queue compact fields, Event & Data Context heading, Discovery numeric cells); deterministic/data/provider logic is unchanged.


> **v14.3.6 Stable:** Exact v14.3.5 baseline. Adds Form 4 transaction-code classification, canonical desk-membership rendering, SPY/QQQ regime context reuse, Discovery layout hardening, and aggregate Macro Rates health reconciliation (Treasury core + optional FRED enrichment) while keeping FRED provider capability health separate. Deterministic Day/Swing/Long score behavior is unchanged.

> **v14.3.4 Stable:** Exact v14.3.3 baseline plus UI/operational corrections only. Adds blocking regressions for latest-five QA metadata, header long-message start preservation, Decision Queue default evidence expansion, Data Engine action geometry, Auto-Start header placement, and canonical VIX manual refresh reuse.


> **v14.3.3 Stable:** Replaces static first-50 Finnhub allocation with a canonical priority allocator (Finnhub primary → Alpaca IEX overflow → snapshot pool), reserves live headroom, publishes per-symbol Live Coverage State, and centralizes generic manual operations in Maintenance with persistent action lifecycle state. Gate 3 correlation/reuse, Gate 5 copy consistency, and Gate 10/11 build-identity checks are blocking.

> Final UI acceptance: Data Engine status values are never clipped or ellipsized. Frequently changing rows reserve a minimum height; all rows may expand for real operational detail.

> **v14.3.2 Completion Final:** Finalizes the audited v14.3 line with header-contained transient notifications and contextual Data Engine manual actions. Manual actions are thin calls into the same scheduler/refresh/evaluation paths; conditional Primary Live Stream Reconnect closes only the existing Finnhub socket so the production reconnect loop handles re-authentication/subscription while fallback remains active. No manual action bypasses entitlement, materiality, canonical-store, or deterministic-score protections.

> **v14.3.1 Completion Patch:** Immutable completion over v14.3.0. Adds verified provider capability state, Twelve fallback adapters/tests, full Market Open checkpoint correlation, exact entry-zone readiness integration, catalyst phase gating, cadence efficiency, and removes the sidebar pseudo-element overlay. Deterministic score math remains frozen.

> **v14.3.0 Improvement Build:** Built from the exact v14.2.2 Stable source. Adds canonical P1/P2 provider intelligence, Market Open Prep, event-driven material-catalyst reaction monitoring, structured provider/preparation state, and separate blocking integration/correlation plus tab-placement audits. Deterministic Day/Swing/Long Setup Score math remains frozen.


This document describes the v14 architecture inherited through v14.2.2 and extended by v14.3.0 and the invariants required to keep the terminal deterministic, efficient, integrated, and understandable. The validated deterministic Day/Swing/Long Action/Score behavior remains a protected baseline; v14 adds global/macro/options context, Trade Readiness, Signal Validation, provider-capability architecture, and cross-module correlation without silently changing that model.

> **v14.2 Stable Hotfix:** Exact v14.0.3 Stable architecture plus three UI-only corrections: notification layout/cleanup, natural-height panel cleanup, and removal of the sidebar Data Engine Details row. The v14.0.3 handoff-compliance + cross-module architecture remains unchanged. All contextual engines feed canonical shared state and are evaluated across relevant modules under Gate 3; deterministic Day/Swing/Long Score/Action math remains protected.


## Permanent G0–G16 Build / Release Process

Effective with v16.0.4, DE.PULSE uses a blocking **Find → Fix → Regression Test → Retest** process. A gate is not report-only: a fixable in-scope or professional-correctness defect is corrected in the same Working Candidate/RC rather than creating patch churn.

- **G0 — Previous Build Closure + Build-Environment Readiness:** protected Stable/source integrity, supported compiler/toolchain, security/package tooling, required network/native capability. Do not start an RC on a host that cannot complete certification.
- **G1 — Master Scope + Product Architecture + Risk:** freeze original wording, IDs, exclusions, canonical owners, surfaces, failure states and acceptance rules.
- **G2 — Source Health + Architecture Fit + Reuse Audit:** inspect the existing tree before coding. Prefer **REUSE → CONSOLIDATE → REFACTOR → DELETE/REPLACE → ADD**. Block orphan production helpers, duplicate implementations, retired runtime paths, parallel canonical owners and silent monolith growth. `source_health_architecture_gate.py` is permanent.
- **G3 — Verification Design:** positive, negative, boundary, recovery, adversarial, fuzz/fault, professional, content/visual and native cases are designed before implementation.
- **G4 — Development + Continuous Fast Tests:** implement incrementally; touched code must leave the codebase cleaner or no worse.
- **G5 — Feature Completion + Failure/Recovery Truth:** prove Data → Validation → Logic → State → Why → UI → Failure → Recovery.
- **G6 — Cross-Module + Canonical Ownership + Decision Truth:** fetch once / calculate once / canonical store; no tab-owned business truth, stale copies or AI/UI divergence.
- **G7 — Code Quality + Maintainability + Performance Health:** formatting, vet/static correctness, resource lifecycle, dead/duplicate code, complexity/performance hot paths, bounded provider/cache work and developer WHY/invariant comments. `code_quality_performance_gate.py` is permanent.
- **G8 — Security + Abuse + Fuzz + Resilience.**
- **G9 — Content + Terminology + Responsive + Visual Fit.**
- **G10 — Professional Trader / Investor Acceptance.**
- **G11 — Independent Original-Scope + Fresh Adversarial Audit:** reread G1, not the closure summary.
- **G12 — RC Freeze + Heavy Certification:** exact frozen RC; normal, vet, race, randomized, coverage, Extreme, deterministic, HTTP, renderer, responsive, security and professional suites. Timing-sensitive/heavy suites run sequentially or in bounded isolated shards to avoid host-contention false failures.
- **G13 — RC Artifacts + SBOM + Provenance + Package Integrity:** packages/provenance/SBOM; SHA manifest is generated last and every entry is immediately rehashed.
- **G14 — Native RC Acceptance:** real Mac/Windows where required, including Stable+RC isolation, persistence, sleep/wake, real credentials and physical UI workflows.
- **G15 — Clean Artifact + Final Release Authorization:** re-extract source, exact-file/hash/provenance/SBOM verification, critical smoke and Master Scope reconciliation.
- **G16 — Stable Release Stamp + Handoff + Learning Loop:** only after authorization; stamp the completed release Stable, record residual risk, escaped-defect learning and explicitly authorize or block the next minor. Working Candidate/RC identity remains internal and never becomes the primary app brand.

### Source-health policy

Source cleanup is continuous, not a once-per-major-release project. Every build cleans touched code; every minor release audits the full affected module path; periodic whole-repository audits address legacy module size, dependencies, caches/provider paths, tests and documentation. Developer comments explain **why, invariants, assumptions, provider/freshness/trading/concurrency decisions and non-obvious optimizations** rather than narrating obvious syntax.

## Architecture Principles

DE.PULSE follows these rules:

- Fetch once, normalize once, reuse everywhere.
- Equities/ETFs/market instruments share a canonical Master Symbol Store.
- VIX uses a separate special-index path.
- v15 routing uses Alpaca IEX/snapshots as the preferred U.S.-equity source, Finnhub as secondary live/recovery, and Twelve Data as tertiary recovery; Alpaca also provides historical and overnight capabilities.
- Derived Day/Swing/Long analyses are horizon-specific caches, not three separate market-data stacks.
- Quote freshness/provenance is independent from derived-analysis caching.
- Research consumes canonical state and stable service data; it does not create a new live subscription system.
- Slow datasets are cached according to information-change rate and refreshed atomically.
- AI is advisory only; Research Confirmation, Risk Review, Catalyst Review, and Custom Question never rewrite deterministic desk outputs.
- Alerts, Portfolio, Journal, order entry, positions/P&L, and autonomous execution remain intentionally absent.

[[diagram:overall]]


## v14 Shared Context Architecture

The v14 rule is: `Fetch once → canonical state → horizon interpreter → relevant surfaces`. Global, macro and options capabilities must not create page-owned provider clients or duplicate quote subscriptions.

`Global/official/free/direct providers → validation/provenance → Engine canonical context → Dashboard / Day / Swing / Long / Queue / Research / AI / Maintenance`

### Global Market Driver Engine
`deriveGlobalMarketContext` consumes canonical quotes, optional direct drivers and macro metrics. Direct data wins when present in AUTO. Real U.S.-listed proxies are permitted fallbacks and are explicitly labeled `LIVE PROXY`/real proxy rather than masquerading as the underlying index. `Direct Only` does not leak ETF proxies; `Proxy Only` does not pretend a direct provider is active. VIX remains separate and is only interpreted by the driver engine.

Broad breadth uses `broadBreadthUniverse`; it must never use a personal watchlist as the market universe. Sector participation is likewise derived from broad sector ETFs.

### Macro Event Engine
Official calendar parsing is conservative. `MacroEvent.TimeKnown=false` requires `StartsAt=0`; date-only events cannot trigger Event Mode. Only an explicit official time is converted using the relevant market-region timezone. Event Mode de-duplicates T−15 history pre-warming per event and captures reaction snapshots at supported offsets without involving AI.

### Options Intelligence
The options layer consumes real provider snapshots, normalizes them to `OptionsContext`, and stores them in the runtime snapshot. In AUTO, Alpaca OPRA is attempted when entitled and real indicative data is the fallback. Expected move uses nearest-expiry near-the-money IV, not an indiscriminate whole-chain IV average. Open interest is not fabricated when the current source does not expose it. Polling is intentionally scoped to the selected symbol plus a small horizon-priority set to avoid chain/API storms.

### Trade Readiness and Signal Validation
Trade Readiness is renderer/context orchestration and never mutates `planFor`/deterministic Score/Action. Signal snapshots are de-duplicated and evaluated only against real future historical bars in live/replay evidence. Demo fixtures never become performance evidence.

### AI compact evidence
The backend builds ticker-specific, review-specific evidence and bounds the payload before provider submission. A context/token-limit response gets one smaller-context retry. Renderer cache identity includes ticker + horizon + review type + evidence fingerprint (and custom question where applicable), so Risk Review and Catalyst Review cannot accidentally reuse Research Confirmation output.

### Cross-Module Integration Gate
Every new or enhanced capability is mapped across `Dashboard | Day | Swing | Long | Decision Queue | Discovery | Research | AI | Settings | Maintenance | Documentation`. Each surface is marked Implemented, Not Relevant, or Intentionally Hidden. Integration must reuse canonical state, preserve horizon-specific relevance, and surface meaningful confirmations/contradictions rather than copy every metric everywhere. This is blocking Gate 11.

## Project Structure

The source package contains:

```text
main.go
main_test.go
v14_features.go
v14_features_test.go
renderer/
  index.html
  renderer.js
  styles.css
  docs/
    user.md
    developer.md
    limitations.md
  qa/
    manifest.json
    v*.txt
platform-icons/
VERSION.txt
README.md
```

The Go process owns persistence, provider access, runtime state, HTTP endpoints, migrations, cache files, and packaging entry behavior. The embedded renderer owns presentation, client-side derived display caches, navigation, symbol handoff, and progressive disclosure.

## Backend Responsibilities

The backend is responsible for:

- Application state/settings/watchlists.
- Provider credentials.
- Finnhub WebSocket/REST lifecycle.
- Alpaca IEX/snapshot/history/overnight lifecycle.
- SEC/filing collection.
- News, earnings, fundamentals, scanner state.
- Canonical RuntimeSnapshot assembly.
- Market-cache load/save and atomic writes.
- Manual stale-data refresh, targeted Pre-Market Preparation, and Saturday non-destructive integrity maintenance.
- AI provider routing and structured response normalization.
- Maintenance/health APIs.

The backend must not add duplicate requests merely because the user navigates from a desk to Research.

## Renderer Responsibilities

The renderer is responsible for:

- Dashboard/Decision Queue presentation.
- Discovery candidate qualification and explicit promotion actions.
- Day/Swing/Long desk presentation.
- Horizon-specific contextual evidence.
- Research symbol/origin state and back-navigation.
- Structured AI Research Confirmation display.
- Progressive disclosure of optional metrics.
- Truthful LIVE/CURRENT/STALE/CACHED/INDICATIVE labels from backend/runtime provenance.
- Page-aware rerendering and derived-analysis cache reuse.


## v14.0.3 Compliance Architecture

The v14.0.3 contextual pipeline is:

`provider adapters → canonical runtime state → provenance/freshness → global/macro/options interpretation → horizon interpreter → Dashboard/Day/Swing/Long/Queue/Research/AI/Readiness/Validation`.

Provider contracts now include `DirectGlobalProvider`, `DirectFuturesProvider`, `OptionsIntelligenceProvider`, and `FastMacroProvider`. Paid entitlements are optional. Twelve Data is an optional direct-global/futures implementation; Alpaca is the current Options implementation. The free/official stack includes Treasury, BLS, BEA, EIA (key where required), FRED (optional key), official/public event sources, and an official TWSE TAIEX close adapter. Direct, official-close and live-proxy evidence can coexist in canonical context; logical-family aggregation prevents those representations from being treated as three independent markets.

Global provider modes are explicit: AUTO, Direct Only, Free First, and Proxy Only. Production never creates a synthetic fallback. Completed official local-market closes retain `OFFICIAL CLOSE` provenance/session status.

### Cross-module release invariant
Every new/enhanced capability is evaluated against Dashboard/Regime, Day, Swing, Long, Queue, Discovery, Research, AI, Trade Readiness, Signal Validation, Settings, Maintenance and Documentation. Each surface is `Integrated`, `Intentionally Hidden`, or `Not Relevant`. Parallel provider subscriptions or tab-owned copies of canonical market state fail Gate 11.

### UI state / shell invariants
`details` expansion state is captured/restored across runtime rerenders where disclosures remain. Dashboard More Market Context is intentionally rendered open with no collapse control. The sidebar uses one scrollable region that includes the detailed Data Engine health surface, while the creator footer remains outside it. Maintenance retains deeper troubleshooting diagnostics. Build Information is sourced from backend `appVersion/buildID`; the renderer displays a `BUILD IDENTITY VERIFIED`/mismatch guard against stale package combinations.

### Event-mode invariants
At T−15 for a known HIGH event, the runtime builds affected symbol/sector context, prewarms history and market-critical providers once, and defers nonessential background fundamentals/filing work. Event reactions use a frozen pre-release price baseline and capture +5s/+30s/+1m/+5m/+15m/+1h observations where runtime cadence permits, including internal capture latency. AI is not invoked by the critical event path.

## Provider and Live-Feed Architecture

[[diagram:feeds]]

Preferred live behavior is session-aware:

- Pre-market / regular / after-hours: Alpaca primary, Finnhub secondary, Twelve Data tertiary recovery.
- Alpaca IEX may be used as controlled fallback without overwriting a fresher protected primary quote.
- Overnight: Alpaca Auto/Indicative/Live logic according to entitlement and configured mode.
- REST snapshots are lower-priority fallback data.

Provider/source/freshness fields are retained in canonical Quote records. Disk-restored quotes are explicitly changed to `DataState=cache` and `FeedType=cache` until replaced.

## Market Session State Machine

[[diagram:sessions]]

Session logic is based on America/New_York and defines overnight, pre-market, regular, after-hours, closed, and weekend states. UI labels may additionally show America/Los_Angeles times.

Day context uses the current session deliberately:

- Overnight: display overnight move vs previous/prior session close when available.
- Pre-market: display pre-market move.
- Regular: emphasize the regular-session open gap and current VWAP.
- After-hours: retain the regular open-gap context plus current after-hours quote provenance.

The session transition changes the contextual lens; it does not destroy historical bars or watchlist state.

## Master Symbol Store

[[diagram:stores]]

`requiredSymbolsFromState` and the runtime subscription manager deduplicate the universe. A ticker appearing in Market Instruments, Day, Swing, Long, or Discovery still maps to one canonical equity/ETF record.

The canonical record may contain:

- Quote.
- Intraday/daily/weekly bars.
- Fundamental snapshot.
- Provider/source timestamps.
- Event/freshness versions used by derived caches.

### VIX Special Index Store

VIX is intentionally excluded from normal equity subscription assumptions. It has distinct provider/health semantics and can be consumed by Market Regime without pretending it is an ordinary Alpaca-IEX equity symbol.

## Subscription Management and Deduplication

The provider lifecycle owns subscriptions. Pages are consumers.

Rules:

- Navigation never creates a second subscription for the same symbol/provider.
- Adding a new watchlist symbol updates the centralized symbol set.
- Engine OFF does not delete a shared subscription that another surface still needs.
- Research/AI navigation does not subscribe independently.
- Immediate symbol hydration is allowed, but must feed the same canonical store.

## Historical Bars and Hydration

Historical bars have two paths:

- Normal Live refresh approximately every 15 minutes.
- Immediate required hydration when missing after symbol add, cache clear, runtime start, or engine re-enable.

This prevents an empty chart from waiting for the periodic timer.

Bar datasets are persisted to disk and restored across restarts. Derived desk calculations version their relevant bar inputs so a changed final bar invalidates only the appropriate horizon result.

## Derived Analysis Cache

[[diagram:cache]]

The renderer maintains derived Day/Swing/Long caches keyed by a calculation-relevant signature. The signature uses quote values, bar versions, earnings/event version, regime version, settings, and fundamentals for Long where applicable.

Cache rules:

- Signature unchanged: reuse the exact deterministic output.
- Signature changed: recompute only the affected horizon/symbol.
- Freshness labels remain dynamic and are not frozen merely because the mathematical result was reused.
- Research/AI confirmation is not part of the deterministic Action/Score signature.

## Market Regime Cache

Day, Swing, and Long may consume different regime inputs/timeframes while sharing the same underlying canonical market symbols. VIX and TLT dependencies remain explicit where used. Regime calculations should be invalidated by their own compact signature rather than every quote event globally.

## Event-Driven Rendering

Quote events update shared chrome/ticker state incrementally. Full-page render requests are suppressed when the changed symbol cannot affect the visible page. Dashboard intentionally remains broad because many tracked symbols may affect the Decision Queue.

Any v13.3.0 visual refinement must preserve this behavior; do not solve information-density problems by creating high-frequency full rerenders.

## Trading Desk Architecture

[[diagram:desks]]

The desk contract remains:

- Compact watchlist for scan speed.
- Detailed deterministic setup table.
- Selected-symbol chart/plan/context detail.
- Empty Add Symbol row at bottom of detailed table.
- Engine lifecycle control.
- Direct Research action preserving symbol and origin.

### Day context contract

Must-have visible metrics:

- Session + provenance.
- Overnight/Pre-Market Move before regular open, Open Gap during regular/after-hours.
- VWAP.
- RVOL.
- Spread.
- Dollar Volume.
- ATR %.

Optional progressive-disclosure metrics:

- RSI.
- EMA 9 / 21.
- Support / Resistance.

Do not add long-term moving averages/fundamental ratios to the Day first-glance lens.

### Swing context contract

Must-have visible metrics:

- SMA 20/50/200 trend structure.
- Relative Strength vs SPY.
- ATR %.
- RVOL.
- Support / Resistance.
- Earnings proximity.

Optional:

- RSI.

Do not duplicate Day-only VWAP/spread/open-gap metrics in Swing.

### Long-Term context contract

Must-have visible metrics:

- Revenue Growth.
- EPS Growth.
- Operating Margin.
- Free Cash Flow.
- Debt / Equity.
- P/E / Forward P/E.
- SMA 50 / 200.
- Relative Strength vs SPY.

Optional:

- ROE.
- 52-week position.
- Current Ratio.
- Dividend Yield.

Do not place intraday spread, VWAP, gap, EMA 9/21, or RVOL in the Long-Term first-glance lens.

## Engine Lifecycle

A desk engine OFF state pauses that horizon analysis but preserves its watchlist and shared canonical market data. Turning an engine ON invalidates/re-enters the relevant analysis lifecycle and requests missing required history immediately.

Provider subscription ownership remains centralized; an engine toggle must not multiply subscriptions.

## Setup Calculations and Equivalence Rule

The deterministic `computePlan` / `planFor` behavior inherited from v13.2.7 is a release invariant for v13.3.0 Stable.

Context metrics and Research Confirmation are explanatory layers. They must not feed back into:

- Action.
- Score.
- Entry.
- Target.
- Invalidation.
- Existing risk/reward result.

Any future formula change requires an explicitly versioned model change plus before/after equivalence/behavior tests. It must not be smuggled into a UI/refinement patch.

## Discovery Architecture

Discovery is candidate generation, not another trading desk.

It owns:

- Scan mode: Day/Swing/Long.
- Filters.
- Candidate reasons.
- Discovery Rank.
- Stage action.
- Explicit Add to current-horizon desk action.
- Research action.

It does not own deterministic desk Action/Score. The UI must continue to label the candidate number **Discovery Rank** to avoid semantic collision.

A scanner-only result can exist before the full symbol is in the tracked Master Symbol Store. Staging/promotion hydrates the canonical state once.

## Dashboard and Decision Queue

Dashboard remains an executive market-intelligence page. It must not become a portfolio/account/P&L screen.

Decision Queue is retained as the primary action funnel. It aggregates useful horizon candidates and provides two main destinations:

- Open the relevant desk.
- Open Research.

Catalyst Pulse provides only a compact preview. Full symbol evidence belongs in Research; full tracked-universe lists remain secondary routes.

## Symbol Research Context Architecture

[[diagram:research]]

Research stores client-side context:

- selected symbol.
- origin page.
- current Research subtab.
- optional AI Research Confirmation for the symbol.

Primary Research subtabs are:

- Overview.
- Fundamentals.
- Catalysts.
- Filings.

A standalone Technicals Research tab is intentionally removed. Technical setup decisions belong in the horizon desk, reducing duplicate charts/metrics and making the flow easier to understand.

Research is optional. It can be entered from Dashboard, Discovery, or any desk and returns to the originating surface with the selected symbol preserved.

## Research Fundamentals

Fundamental data is consumed from the canonical runtime snapshot and persisted market cache. Research displays must-have fields first and optional fields under progressive disclosure.

Long-Term can reuse the same fundamental record; Research does not fetch a second copy because its tab was opened.

## Catalysts, News, Earnings, and SEC

News, Earnings, and SEC/Filings remain backend services and secondary renderer routes, but they are not primary sidebar tabs in v13.3.0 Stable. Dashboard Catalyst Pulse and Research provide the primary contextual entry points.

Primary integration:

- Dashboard Catalyst Pulse previews them.
- Research Catalysts provides selected-symbol News + Earnings.
- Research Filings provides selected-symbol SEC evidence.
- All News / All Earnings / All Filings buttons access secondary broad views.

This preserves capability while removing primary-navigation duplication.

## AI Routing and Structured Research Confirmation

The frontend sends an `AIRequest` with the selected symbol plus `ClientContext`, including current session, quote state, deterministic Day/Swing/Long outputs, and each desk's must-have contextual evidence.

The backend prompt explicitly forbids trade execution and score mutation. AI providers must return valid JSON with:

- `verdict`: FAVORABLE, MIXED, CAUTION, INFORMATIONAL.
- `confidence`: 0–100.
- `reasons`: max 3.
- `risks`: max 3.
- `catalyst`.
- `bestFitHorizon`: day/swing/long/none.
- `nextAction`: user-controlled next step.
- `summary`: one sentence.
- `details`: bounded evidence text.

`parseAIStructuredPayload` normalizes/clamps this output. If a provider returns unstructured text, the server falls back to INFORMATIONAL and exposes the raw content as bounded details rather than treating it as a reliable verdict.

The renderer stores the structured result per symbol in `aiResearchResults` for the current UI session and displays it in AI Copilot and Research Overview.

AI can expose explicit buttons such as Open Swing or Add to Swing. It cannot perform those actions without the user's click.

## Local HTTP API and Security

Important endpoints include existing runtime/settings/watchlist/provider routes plus:

- `/api/cache/clear` for explicit destructive market-cache clearing.
- `/api/cache/refresh` for non-destructive low-priority stable-data refresh.
- `/api/ai/generate` for structured AI Research Confirmation.

Local endpoint authentication/session protections remain unchanged. External URLs are restricted to valid HTTP/HTTPS targets.

## Persistence and Schema Migration

Application state and market cache are separate.

Application/profile state includes settings, watchlists, selected UI state, maintenance timestamps, and non-secret configuration. Secrets remain separately handled and are excluded from profile export.

The market cache is version-tolerant JSON. v13.3.0 persists:

- Quotes.
- History.
- Intraday/daily/weekly bars.
- Fundamentals.
- News.
- Earnings.
- Filings.
- SEC intelligence summaries.
- Scanner state.
- Per-dataset `lastUpdated` timestamps.
- `savedAt`.

Older cache files lacking the new optional fields must still load safely.

## Cache Lifecycle and Freshness Policy

Cache behavior is data-class specific.

### Live quotes

Event-driven. Disk persistence is for restart continuity only. Restored quotes are marked CACHED until replaced by a provider.

### Historical bars

Normal Live refresh: 15 minutes. Immediate hydration when required.

### News

Live runtime refresh: approximately 10 minutes. Persist across restarts.

### Earnings

Refresh: approximately 2 hours. Persist across restarts.

### SEC filings

Refresh: approximately 30 minutes. Persist across restarts.

### Fundamentals

Refresh: approximately 24 hours, with event-driven eligibility when new earnings/financial filing evidence makes fundamentals stale. Persist across restarts.

### Disk cache save cadence

`cachePersistLoop` continues to provide bounded persistence, but `saveCache` fingerprints cache content with the volatile `SavedAt` field neutralized. If canonical cache content has not changed, the physical cache file is not rewritten. When content changes, the cache is written through `atomicWrite`, so a replacement file is committed only after the new JSON has been prepared. Runtime stop still performs a final persistence attempt.

### Pre-Market Preparation

`preMarketPrepLoop` checks periodically and runs once per eligible trading day during approximately 3:15–3:50 AM ET. It skips weekends and known US market holidays.

Crucial invariant: **Pre-Market Preparation does not stop live/overnight quote streams.** It calls the low-priority stable-data refresh pipeline for stale/missing history, fundamentals, earnings, news, and SEC, then records preparation health/timestamps. It is targeted refresh, not destructive cache clearing.

### Weekly Integrity

`weeklyIntegrityLoop` performs a Saturday non-destructive integrity pass. It normalizes cache structures, verifies persisted-state health, records the integrity result, and saves only if content changed. It must not change watchlists, settings, credentials, selected ticker, trading-engine switches, or provider subscriptions.

If a stable provider refresh fails, last-known-good data remains available with older provenance. New data replaces existing state only after a successful response/validation path.

### Manual stale-data refresh

Maintenance calls `/api/cache/refresh`. In Live mode the backend starts the same low-priority stable-data pipeline in the background, rate-limited so repeated clicks do not start an immediate storm. Demo mode returns a no-provider-refresh-needed message.

## Maintenance Architecture

Settings configures. Maintenance inspects and repairs.

Maintenance includes:

- Weekly non-destructive diagnostic run.
- Central Data Engine health.
- Cache metrics/activity.
- Freshness & Cache Policy.
- Last stale/manual refresh timestamp, Pre-Market Preparation status, and Weekly Integrity status.
- Refresh Stale Data Now.
- Explicit Clear Market Cache.
- Version/build identity.
- Latest five full QA reports.

The `cache-refresh` health key is part of the same centralized health model; Maintenance does not create a second provider checker.

## QA and Release History

The package contains only the latest five full QA reports:

- v13.3.0 Stable.
- v13.2.9 Test.
- v13.2.8 Test.
- v13.2.7 Stable.
- v13.2.6.

Historical report text is preserved as historical evidence; old reports are removed from the package rather than rewritten.

## Documentation System

Documentation is bundled locally/offline. Markdown supports headings, lists, code blocks, brand normalization, and native responsive architecture diagrams through tokens such as `[[diagram:overall]]`.

User and Developer titles use the visual DE.PULSE brand component followed by ` - User Documentation` or ` - Developer Documentation`.

## Global Visual Brand Component

Every visible product-name treatment uses the reusable DE.PULSE component:

- uppercase DE and PULSE.
- warm ivory letters.
- independent silver center dot.
- display and inline alignment variants.

Plain-text internal/log/accessibility values may remain `DE.PULSE`.

## Global Dropdown Affordance

Every actual `select` must retain a visible right-side chevron, sufficient right padding, and consistent hover/focus state. Plain text inputs do not display the chevron.

## Build and Packaging

The release process builds:

- macOS ARM64.
- macOS Intel.
- a universal macOS package containing both.
- Windows x64 GUI binary/package.
- Source ZIP.
- QA report.
- SHA-256 manifest.

Approved platform icon assets must remain byte-identical unless an icon change was explicitly requested.

## Release Regression Gates

Before handing off a v14.3.0 Stable package:

- `go test ./...`.
- `go test -race ./...`.
- `go vet ./...`.
- renderer JavaScript syntax check.
- renderer logic/static tests.
- deterministic trading-result equivalence against the Stable model inputs.
- engine OFF/ON behavior for all three desks.
- no duplicate subscription regressions.
- cache load/save round-trip including new research datasets.
- cache clear preserves settings/watchlists and rehydrates history.
- manual cache refresh does not stop live streams.
- structured AI parser normalization/fallback tests.
- Discovery → Desk → Research → Back handoff tests.
- Decision Queue navigation tests.
- responsive UI matrix with no page-level horizontal overflow.
- previous test/stable state migration.
- final Source ZIP re-extraction and retest.
- package-content/build-identity validation.

Native macOS WKWebView/Gatekeeper/Finder-icon behavior, Windows shell/Explorer behavior, and authenticated provider behavior with private user credentials still require real-machine acceptance testing.

## Extension Rules

When adding future data:

- Decide whether it is canonical raw state, derived horizon state, Research evidence, or UI-only presentation state.
- Do not create a new provider loop if an existing canonical dataset can serve the feature.
- Match refresh cadence to how quickly the information can change.
- Preserve source and last-updated metadata.
- Prefer stale-while-revalidate to deleting known-good slow data.
- Put must-have metrics first and optional metrics behind progressive disclosure.
- Do not mix Discovery Rank, deterministic desk Score, and AI Research Confirmation into one opaque number.
- Any future scoring-model modification must be explicitly versioned and tested separately from UI/cache work.

> **v14.0.3 Stable:** The Stable release adds global/macro/options intelligence, Trade Readiness, Signal Validation, premium-ready provider interfaces and cross-module correlation to the horizon-aware v13.3 architecture while preserving the validated deterministic setup model.


### Options IV change
`OptionsContext.IVChange` is derived only from consecutive real normalized option snapshots and is stored as percentage-point change (`(newIV-oldIV)*100`). No synthetic prior IV is created.

## v14.0.3 UI consistency / responsive-stability notes

The patch is renderer-only apart from release identity. Global Driver cards use a two-row internal grid so title/status/detail cannot occupy the same line. The Market Instruments tape preserves identical normal/hover dimensions and no longer pauses its transform animation on pointer hover. Event/Data detail context uses a scoped grid rather than modifying the generic compact-row pattern. SEC summary typography is similarly scoped. A >30 minute cache age receives a renderer warning class; underlying quote/fallback behavior is unchanged. Focused Playwright geometry tests are blocking alongside the standard responsive matrix.

## v14.3.0 Engineering Notes

### Canonical reuse and correlation
New provider fields are normalized into canonical raw/derived state and reused rather than fetched/calculated independently per page. `RuntimeSnapshot` exposes preparation state, Liquidity Health, derived intelligence, provider capability registry, symbol intelligence, catalyst reactions, Market Open flags, market-activity context, and relevant corporate actions. AI receives those same canonical facts and remains non-authoritative.

### Event scheduling
- Market Open Prep is market-calendar aware and scheduled for the pre-bell 9:20–9:25 AM ET window on U.S. trading days.
- Earnings & Material Catalyst Reaction Watch is selective/event-driven. Scheduled earnings arm only the applicable BMO/AMC date/session. Unexpected material News/SEC events may arm it through the normal event pipeline. Low-impact events remain on normal cadence.
- Slow official macro/energy datasets are cached according to source cadence; they are not polled as if they were tick data.

### Deterministic-score protection
The existing Day/Swing/Long Setup Score/Action formula path is frozen in this build. New liquidity, macro, catalyst, provider-intelligence, and Market Open evidence may change Trade Readiness/context presentation but not the deterministic Setup Score. Deterministic-equivalence regression against exact v14.2.2 is a blocking test.

### v14.3.0 blocking release-gate order
1. Unit Validation
2. Functional & Workflow Validation
3. Cross-Module Integration, Correlation & Reuse
4. Tab-by-Tab UI Placement & Information Hierarchy Audit
5. UI/UX & Responsive Validation
6. Performance & Data-Efficiency Validation
7. Source Cleanup & Engineering Quality
8. Engineer Acceptance Test
9. Professional Trader Acceptance Test
10. Documentation & Operational Readiness
11. Platform & Package Validation
12. Exact-Source Freeze & Final Retest

Gate 3 checks canonical reuse, correlation, and relevance across surfaces. Gate 4 is intentionally separate: after integration, each tab is audited for decision priority, grouping, spacing, duplication, and horizon-appropriate placement.


## v14.3.6 Stable engineering notes

- `enrichForm4` reads canonical Form 4 XML transaction codes and promotes only P/S to BUY/SELL. All other supported SEC transaction codes remain OTHER. The enriched canonical `FilingItem` is reused by SEC Filings, Research, horizon SEC context and SEC Intelligence.
- `deskMembershipStrip` reads the permanent Day/Swing/Long watchlists directly; Discovery and desk watchlists reuse the same membership fact.
- `regimeQuotePill` reuses canonical SPY/QQQ quote state. `regimeDirectionalState` is a presentation-only mapping of the existing horizon Market Regime score and never writes back into deterministic scoring.
- Global Market Drivers state/dot is rendered from `GlobalMarketContext.Tone`, whose backend aggregation already combines logical driver families and avoids double-counting direct/official/proxy variants.
- `reconcileMacroRatesHealth` separates aggregate Macro Rates health from source-specific `fred-rates` health. Usable current Treasury core rates are sufficient for healthy aggregate core-rate availability.
- Catalyst Watch stores READY / ARMED / TRIGGERED / REACTION between evaluations; manual-action RUNNING state is kept separate.
- Discovery uses explicit responsive column ownership. Lower-width containment is internal to the scanner table; document-level horizontal overflow remains release-blocking.
- Gate 4/5 includes a fit/placement rule: changed UI is incomplete if it overlaps, clips, wraps awkwardly, creates unnecessary dead space, or is visually detached from its workflow context.


## Extreme Production / Professional Trader Acceptance gate

Starting with v15.1.2, `extreme_30_matrix_test.go` and `extreme_30_gate.py` are permanent release tests. They cover the 30 production-abuse categories documented in `renderer/qa/v15.1.2-extreme30-matrix.md` and run alongside unit, race, randomized-order, HTTP workflow, deterministic-equivalence, Professional Trader, responsive/layout, content, version, approved-scope, exact-source, and platform-package gates. Native target-OS GUI launch and authenticated live-provider entitlement acceptance remain environment-dependent final checks.

## Permanent G0-G16 Build Process and Checkpointed Certification

Starting with the v16.0.4 source-health retrofit, DE.PULSE uses the permanent G0-G16 build/release process. Source health and architecture fit are inspected before feature addition; code quality and performance health are rechecked after implementation; heavy RC certification is checkpointed so successful checks are not repeated after runner interruption.

### Gate sequence

- **G0 - Previous Build Closure & Build-Environment Readiness.** Protect Stable, validate previous closure/runtime isolation, approved Go/toolchain, security tools, packaging capability and native-test capability. Environment/tooling gaps are **BLOCKED**, not product defects.
- **G1 - Master Scope + Product Architecture + Risk.** Freeze original wording, requirement IDs, exclusions, canonical ownership, surfaces, failure states and acceptance rules.
- **G2 - Source Health + Architecture Fit + Reuse Audit.** Decide REUSE / CONSOLIDATE / REFACTOR / DELETE / REPLACE before ADD. Block orphan helpers, retired runtime paths, duplicate ownership and uncontrolled legacy-module growth.
- **G3 - Verification Design.** Define positive, negative, boundary, recovery, adversarial, fuzz, fault-injection, content/visual and native-UAT cases before coding.
- **G4 - Development + Continuous Fast Tests.** Implement incrementally; touched code must leave the source no worse than before.
- **G5 - Feature Completion + Failure/Recovery Truth.** Prove Data -> Validation -> Logic -> State -> Why -> UI -> Failure -> Recovery.
- **G6 - Cross-Module + Canonical Ownership + Decision Truth.** Fetch/calculate once, reuse canonical stores, prevent competing freshness/state/AI/UI truth.
- **G7 - Code Quality + Maintainability + Performance Health.** Enforce formatting/static quality, resource lifecycle, WHY/invariant comments and focused performance/concurrency regressions.
- **G8 - Security + Abuse + Fuzz + Resilience.** Vulnerability/dependency scan, secrets/session checks, malformed/untrusted input, concurrency/network/provider chaos and fault injection.
- **G9 - Content + Terminology + Responsive + Visual Fit.** Fix wording/casing/wrapping/clipping/density/alignment using realistic worst-case content.
- **G10 - Professional Trader / Investor Acceptance.** Verify evidence/state/why/uncertainty is trustworthy enough to influence real capital.
- **G11 - Independent Original-Scope + Fresh Adversarial Audit.** Re-read the G1 contract and invent fresh scenarios rather than validating only the implementation summary.
- **G12 - RC Freeze + Heavy Certification.** Run the exact frozen source through full/race/vet/randomized/coverage/Extreme/deterministic/HTTP/renderer/responsive/professional suites. CPU-heavy checks run as bounded independent processes.
- **G13 - RC Artifacts + SBOM + Provenance + Package Integrity.** Build target packages and Source ZIP; generate SBOM/provenance; create the SHA manifest last and immediately reverify every delivered file.
- **G14 - Native RC Acceptance.** Real macOS/Windows lifecycle, persistence, Stable+RC isolation, sleep/wake, real credentials/entitlements and physical visual-fit acceptance where required.
- **G15 - Clean Artifact + Final Release Authorization.** Re-extract Source ZIP, exact-file hash check, identity/provenance/SBOM checks, critical truth/deterministic/security smoke and exact Master Scope reconciliation.
- **G16 - Stable Release Stamp + Handoff + Learning Loop.** Only after authorization: stamp the completed release Stable, record residual risk/escaped-defect learning and explicitly authorize or block the next minor. Working Candidate/RC identity remains internal.

### Checkpointed certification runner

`certification_runner.py` is the permanent orchestration authority for automated release checks. `certification_plan.json` defines the ordered checks. Every expensive item runs as a separate process with its own timeout, log, duration and durable result in an external certification directory beside the source tree (`DEPULSE_CERT_DIR` can override it). Certification output is never a source input.

Result states are intentionally distinct:

- **PASS** - the exact-source check succeeded.
- **PRODUCT FAIL** - DE.PULSE/test behavior failed and the Working Candidate must be fixed.
- **INFRA FAIL** - the runner/host interrupted or could not complete the check; do not diagnose it as a product defect.
- **BLOCKED** - a mandatory prerequisite such as approved toolchain, vulnerability tooling or native target-host acceptance is unavailable.

PASS checkpoints are reused only while the exact source fingerprint is unchanged. Any source/document/gate modification changes the fingerprint and automatically invalidates prior PASS results. Generated certification logs/results are excluded from the fingerprint.

Common commands:

```text
python3 certification_runner.py --status
python3 certification_runner.py --phase G12
python3 certification_runner.py --check g12_random_1600405
python3 certification_runner.py --all
python3 certification_runner.py --reset
```

Randomized-order seeds and race-detector shards are individual checks. If a session stops after seed 4 or race shard 7, the next run resumes at the first non-PASS check instead of repeating the already-clean work. Timing-sensitive/heavy suites should not be run concurrently on the same release host.

### Source cleanup and comments

Every build performs touched-code cleanup. Every minor release performs a broader affected-module audit, and periodic full-repository health audits are required as architecture debt accumulates. Developer comments explain **WHY, invariants, trading assumptions, canonical ownership, freshness/fallback behavior, concurrency and non-obvious performance decisions**; comments should not restate obvious syntax.

### v16.0.6 final delivery and launch rule

The production app package must contain the actual native DE.PULSE binaries before delivery. Release certification, vulnerability scanning, and evidence generation are separate release operations and must never replace or block the normal app-launch path. The primary header remains Stable-style; release channel/status belongs in release artifacts and Application & Build Identity.

### v16.0.5 final-closure identity rule

Runtime/UI identity carries only semantic version + build ID. Working Candidate / RC / Stable is release-artifact status, not primary brand-header content. Completed user-facing releases are Stable; Working Candidate/RC labels remain internal certification states. Promotion must never require a cosmetic header source change that invalidates exact-source certification.

> **RL-023 / G11/G12/G16 current-release identity:** HTTP/integration fixtures derive current version/build expectations from `VERSION.txt`; predecessor release literals are used only for historical compatibility assertions.

## v18.3 persistence archive, migration and recovery contract

The persistence archive is a versioned backend-neutral snapshot of the canonical persistence repository: Global Symbol Registry, canonical quotes, quote history, evidence, decision lineage, outcomes, derived features, IdentityPersistentState and per-user UserWorkspace state. Archives are wrapped in a SHA-256 integrity envelope and written atomically with private `0600` permissions because identity password/session hashes are security-sensitive. Provider API secrets remain outside this archive.

Operational export uses `DEPULSE_PERSISTENCE_EXPORT_PATH`. Hosted migration/restore uses `DEPULSE_PERSISTENCE_RESTORE_PATH`; restore defaults to `empty` and rejects a non-empty target. `DEPULSE_PERSISTENCE_RESTORE_MODE=replace` is the explicit destructive restore mode. Restore occurs after repository initialization/migrations but before IdentityService/workspace bootstrap so imported identity/workspace truth remains canonical.

Runtime database readiness is distinct from initialization. `PersistenceManager` retains/coalesces bounded pending work when PostgreSQL becomes unavailable, marks readiness `DEGRADED`, and uses one exponential health-probe retry lane (250 ms to 5 s cap) with two-success recovery hysteresis. `/api/ready` performs a rate-limited persistence health probe and returns 503 while canonical persistence is unavailable; `/api/health` remains process liveness. The intelligence persistence queue deduplicates immutable IDs/feature keys and has a hard 50,000-record ceiling, shedding reproducible derived features first with explicit diagnostics rather than allowing unbounded memory growth.
