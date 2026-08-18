# DE.PULSE — User documentation

## v18.6.0 STABLE — Adaptive utility and intelligence hardening

DE.PULSE now shares broad Scanner/Opportunity Radar snapshot acquisition, coordinates Pre-Market Prep and Market Open Prep through one serialized session-intelligence scheduler, and keeps low-value Market Activity / legacy evidence routes in supporting drill-down paths. Documentation is role-aware and privileged developer content is enforced server-side. Provider/dependency readiness is explicit and fail-closed. External AI research uses bounded material context, strict structured outputs, expiring cache identity, continuous-evaluation telemetry, and provider×dataset rights filtering before egress. Deterministic Day/Swing/Long formulas, the U.S. Equities Processing Boundary, Smart Provider Router ownership, and the permanent No Execution Boundary remain unchanged.


## v18.5.2 STABLE — Hotfix recovery candidate

Day, Swing, and Long-Term desks render through a self-contained symbol normalizer. Complete ET/PT clocks live in the compact Market Pulse Ribbon, while Research keeps its full evidence hierarchy. Tracked-symbol typing survives live quote updates and Add Symbol succeeds on the first attempt. Settings now allows a password-reverified display name and sign-in username while OWNER remains a separate role. The full Save Settings row stays visible inside the application window. Deterministic formulas, U.S.-listed scope, and No Execution are unchanged. Promotion remains pending exact-source local and native package certification.

## v18.5.1 STABLE — Recovery patch

Incremental live updates preserve row identity/hover/focus/selection/scroll. Desk-row × removes a tracked symbol globally while desk pills retain final-membership protection; Undo restores exact memberships/selection. Research/header responsive hierarchy and version-neutral profile copy are hardened. Stable continuity uses the canonical PersonalMarketTerminal profile. No Execution is unchanged.


## v18.5.0 STABLE — Major Closure & Release Assurance

v18.5 STABLE is the Major Closure over **v18.4.0 STABLE**. Final publication is gated by exact Stable G11-G15 certification, including required native artifacts. It does not add unrelated trading features. It re-certifies the approved v18 platform and treats runtime overload / intermittent DATA DEGRADED / slow response as release-blocking if local load can delay or misstate decision-critical evidence. Readiness must stay conservative and truthful under provider, database, queue/backpressure, restart, multi-user/fan-out and recovery pressure. The permanent **No Execution** boundary is unchanged.

Release recovery is now explicit: final Stable runnable packages belong in the GitHub `depulseapp/DE-PULSE` Release for `v18.5.0-stable`, with source/SHA/certification evidence. Historical originals remain available under ChatGPT Library `/DE.PULSE/<version>/` while GitHub backfill is normalized.


## v18.4.0 STABLE — Security / commercial readiness hardening

v18.4 STABLE adds fresh password confirmation for security-sensitive settings/admin changes, hosted request-abuse safeguards and explicit provider commercial/data-rights readiness. Normal research/market workflows are not re-auth gated. Provider commercial, redistribution and AI-use rights are never inferred from a configured API key; unreviewed rights remain blocked/review-required in governance metadata.

The Stable build uses canonical `PersonalMarketTerminal` and preserves compatible prior Stable state; the historical `PersonalMarketTerminal-v18.4.0-TEST` profile remains isolated. Desktop storage remains SQLite by default; hosted PostgreSQL/shared-state behavior, Smart Provider Router execution, deterministic Day/Swing/Long formulas and the **No Execution** boundary remain unchanged.

## v18.3.0 STABLE — PostgreSQL / hosted shared state

v18.3 STABLE adds the certified hosted-state foundation without changing normal desktop storage. macOS and Windows desktop builds continue to use the existing local SQLite profile by default. A hosted build explicitly selects PostgreSQL; if that hosted database is unavailable or misconfigured, readiness fails instead of silently using a local substitute.

Personal watchlists and UI state remain per-user. Market-data/provider/scanner/intelligence work remains one shared canonical pipeline across the deduplicated union of required symbols, so adding users does not create one market engine per account. `/api/health` remains liveness; `/api/ready` reports whether canonical identity and persistence are actually ready. Migration/export, backup/restore, contention/recovery, PostgreSQL 17 hosted runtime, and required native release qualification passed G0–G15 before Stable promotion. The **No Execution** boundary is unchanged.

## v18.2.0 STABLE — Administration, Presence & Sessions

The Stable channel uses the canonical `PersonalMarketTerminal` runtime/profile and preserves compatible prior Stable state; TEST-profile isolation is not used after promotion.

Privileged OWNER/ADMIN-class users now have a compact Settings administration surface for creating lower-authority accounts, changing permitted roles/status, resetting temporary passwords, viewing ACTIVE / IDLE / OFFLINE presence and revoking eligible sessions. Presence comes from authenticated session truth; there is no separate heartbeat database. Credential hashes and opaque session tokens are never shown.

Role and lifecycle changes revoke affected sessions, and a revoked/expired long-lived SSE connection closes on its next keepalive. USER/DEMO do not see administration controls. v18.1 personal market workspaces remain isolated, while provider/evidence/scoring intelligence remains shared and efficient. Hosted PostgreSQL/browser deployment remains v18.3. The **No Execution** boundary is unchanged.

## v18.1.0 STABLE — My Market Symbols / personal workspace

Each authenticated user now has a durable personal market workspace keyed by the account's immutable `UserID`. Tracked Symbols, Day/Swing/Long memberships, personal watchlists, scope and selected ticker no longer share one global list. New users start with an empty personal symbol set; the existing Stable user's state is migrated once into the OWNER workspace.

Market evidence remains shared and efficient. DE.PULSE still fetches/calculates canonical evidence once for the deduplicated union of needed symbols, while Bootstrap, runtime diagnostics and symbol-specific streaming are filtered to the signed-in user's workspace. Provider configuration and other global operational settings remain ADMIN-owned. This release does not add the v18.2 user-administration/presence/session-management surface. The **No Execution** boundary is unchanged.

## v18.0.6 STABLE — Router / Market Shock hardening

v18.0.6 STABLE runs in the isolated `PersonalMarketTerminal` profile cloned from v18.0.5 Stable on first use. User-facing behavior remains intelligence-first: material Rapid Move events may now be identified as **MARKET SHOCK** when correlated broad-market context is present, while low-level provider and adaptive-policy machinery remains hidden from normal USER/DEMO views.

## v18.0.5 STABLE — UI/UX + Symbol Management Hardening

v18.0.5 is the promoted Stable patch over **v18.0.4 STABLE**. **Tracked Symbols** is the user-facing cross-desk control for Day / Swing / Long-Term; Add, Remove and Remove All now converge through one canonical mutation path and preserve an intentional empty state across reload/restart. Opportunity Radar and Research Target receive focused responsive/UI hardening. USER/DEMO surfaces show conclusions, freshness and material degraded-data meaning instead of provider/cache/scheduler/database internals; privileged Maintenance diagnostics remain available. Protected deterministic formulas and the **No Execution** boundary are unchanged.

## v18.0.4 STABLE — Native cross-platform closure

v18.0.4 is the promoted Stable v18.0.x foundation. Windows SQLite lifecycle cleanup, canonical health `buildId`, native restart persistence, authentication, Smart Router v2, Rapid Move, and Coverage Truth cleared macOS Apple Silicon + Windows x64 G14 and G15. Stable promotion changes release identity/runtime targeting only; market-intelligence behavior is unchanged.

## v18.0.3 TEST — Native runtime portability hardening

v18.0.3 fixes native cross-platform release blockers found by G14. The Windows app now resolves embedded setup/login/static resources with platform-neutral paths, while TEST-profile migration and test isolation use the correct OS user-config location. Market-intelligence behavior is unchanged from v18.0.2/v18.0.1. The build remains TEST until native Apple Silicon + Windows x64 G14 and G15 pass.

## v18.0.2 TEST — Native delivery hardening

v18.0.2 fixes a cross-platform release-verification defect discovered during native G14. The market-intelligence behavior from v18.0.1 is unchanged; this patch makes one exact source tree verify with one canonical fingerprint on macOS and Windows and uses a separate `PersonalMarketTerminal-v18.0.2-TEST` profile.


## v18.0.1 TEST — Smart Router + Rapid Move intelligence

v18.0.1 is a separate TEST build over **v17.5.1 STABLE** and the v18.0 identity/session foundation. It uses `PersonalMarketTerminal-v18.0.1-TEST`, so it does not write into Stable or the frozen v18.0 TEST profile.

When a currently observed U.S.-listed symbol makes an unusually fast, material move, DE.PULSE evaluates 15s/30s/60s/2m/5m velocity together with liquidity, source agreement, market/session context, corporate-action risk and fresh catalyst evidence. A material first alert appears immediately in the header, is retained in Smart Notifications, and can promote the symbol into Opportunity Radar/live priority for investigation. The event can evolve through `EARLY → VALIDATING → CONFIRMED → EXTENDED → FADING/RESOLVED`; an unexplained validated move can alert before a headline arrives.

Coverage is intentionally truthful: v18.0.1 provides exact short-window detection for symbols receiving current canonical observations and does not claim full-market rapid-move surveillance without an entitled broad event feed. This is decision support only; deterministic Day/Swing/Long formulas and the **No Execution** boundary are unchanged.

## v18.0.0 TEST — Identity & Secure Session Foundation

v18.0.0 is a **separate TEST build** over the authoritative v17.5.1 Stable baseline. On first launch, close v17.5.1 Stable once so DE.PULSE can clone the compatible Stable profile into `PersonalMarketTerminal-v18-TEST`; the Stable profile is not modified. After that one-time copy, Stable and v18 TEST can be reopened independently.

The first v18 TEST launch opens **Secure the Owner account** before the market workspace. Create a password of at least 12 characters. DE.PULSE stores an **Argon2id** password hash, uses opaque server-side sessions, rotates/revokes sessions on credential changes, applies idle and absolute expiry, and provides **Sign Out** in Settings. Existing compatible settings, watchlists, API keys and persisted market intelligence are carried into the isolated TEST profile.

This slice is the security foundation only. The role model exists, but v18.1 per-user market symbols and v18.2 user/admin/presence operations are not part of v18.0. Deterministic Day/Swing/Long Score/Action formulas and the permanent No Execution Boundary are unchanged.


## v17.5.1 STABLE — Release identity & documentation hardening

v17.5.1 is a **non-functional Stable hardening patch** over the fully certified v17.5.0 Major Closure bits. It corrects the runtime/renderer release identity from RC to STABLE, removes obsolete pre-promotion wording, and restores the complete v17.0→v17.5 delivery history in the in-app documentation. Trading/data logic, deterministic Day/Swing/Long Score/Action formulas, provider ownership and the No Execution Boundary are unchanged.

## v17.5.0 Stable — v17 Major Closure & Release Assurance

v17.5 completed the v17 Major Closure: the frozen 20-item v17 scope was reconstructed from current code, inherited v16 safety/decision contracts were rerun, Principal Engineer and Professional Trader/Investor acceptance passed, native macOS SQLite acceptance passed, and v17→v18 GO was authorized. v17.5 added no unrelated trading feature.

## v17.4.0 — Completed UX + Operational Hardening checkpoint

Master Market Symbols received balanced responsive controls. Preparation exceptions were grouped by real root cause with counts and expandable per-symbol Review. Stale quote evidence remains a freshness/data-health problem instead of being multiplied into misleading per-symbol liquidity warnings. Performance & Runtime Load diagnostics remain available in Maintenance.

## v17.3.0 — Completed Performance + Reliability Hardening checkpoint

v17.3 converted the v17 runtime architecture into measurable acceptance: selected/actionable freshness, API latency, workload queues, CPU/memory/goroutines, provider/subscription budgets, warm-start coverage, provider-call reuse, DB write rate, storage growth and recovery from stale/degraded states. Maintenance/Data Engine exposes these diagnostics.

## v17.2.0 — Completed Canonical Intelligence + Data Efficiency checkpoint

v17.2 strengthened the one-canonical-owner model: material-change propagation, structured evidence/history, immutable Decision Lineage, separate outcomes, incremental/source-hash derived features and provider-call reuse. It also closed a live Alpaca IEX self-deadlock risk by resolving allocation before taking the write lock.

## v17.1.0 — Completed Runtime Load + Backpressure checkpoint

v17.1 introduced Tier 0–4 workload priority, bounded provider queues, protected Tier 0/1 reserve, cancellation-safe waiting, low-priority shedding, provider-budget telemetry and live-subscription capacity diagnostics. Critical/actionable work is protected before broad discovery/background work under pressure.

## v17.0.0 — Completed Persistence Foundation checkpoint

v17.0 established the repository boundary, SQLite schema/migrations, retained Global Symbol Registry, persistent canonical quote warm start, async/material-change persistence, structured evidence/lineage foundations, runtime/provider observability and expanded Data Utility governance. The initial Windows SQLite packaging gap was subsequently closed in v17.1.

## v16.11.0 Stable — v16 Major Closure & Release Assurance

v16.11 is the mandatory final assurance pass before v17. It does not add a new trading surface. The complete v16 family is freshly rechecked from current code, including the original **30 FULL / 0 PARTIAL / 0 MISSING** professional roadmap and the v16.10 **10/10** Opportunity/Decision Intelligence scope.

For real-money decision support, current-price freshness now directly constrains **Trade Readiness**. During the regular session, `STALE`, `CACHED` or `HISTORY ONLY` current-price evidence cannot show a clean Day `READY` state: Day is `INCOMPLETE`; Swing/Long are `CONDITIONAL`. `LAST TRADE` in a Day context requires live-session confirmation. This is a safety/readiness overlay only; deterministic Day/Swing/Long Score/Action formulas remain unchanged.

The closure review also rechecks Opportunity Radar, liquidity/slippage, macro/event risk, Community Intelligence, replay/no-lookahead, provider degradation, responsive behavior and settings/API-key continuity. **Setup Score is not win probability**, Radar is investigation priority rather than a BUY/SELL signal, and DE.PULSE remains research/decision support only with no execution capability.

The v16 → v17 transition is authorized only after the Major Closure Scope Matrix, full blocking certification, senior-engineer review and professional trader/investor real-money review are green. This includes an independent **Professional Trader / Investor** acceptance track.

## v16.10.0 Stable — Opportunity & Decision Intelligence

**Discovery** now contains the Always-On **Opportunity Radar**. It watches the eligible U.S.-listed universe through bounded snapshot/batch scanning, reuses Most Active/Gainers/Losers and existing market evidence, and promotes only a small number of qualified symbols for deeper/live observation. Qualification considers session-normalized relative volume, volatility/range expansion, price confirmation, liquidity/spread, persistence and material catalyst context. A promotion is an investigation signal, **not a promise of profit or a BUY instruction**.

Radar promotions expire/demote when conditions normalize and cannot starve your active watchlist. **SPY/QQQ** retain market-critical priority and **GLD/SLV/USO** retain their approved tradable priority. Maintenance now exposes adaptive freshness/cache policy and read-only Shadow observations. Shadow follows **SHADOW → VALIDATED → APPROVED → PRODUCTION** and cannot change production behavior until explicitly validated/approved. Original professional roadmap status remains **30 FULL / 0 PARTIAL / 0 MISSING**. `De-Pulse.app` / `PersonalMarketTerminal` continuity is preserved.

## v16.9.0 Stable — Final Original Professional Roadmap Closure

DE.PULSE now closes the original professional roadmap at **30 FULL / 0 PARTIAL / 0 MISSING**. **Global Community Evidence Fusion** normalizes permitted/user-authorized X, Reddit, Telegram, Discord, WhatsApp and manual evidence where supported, keeps it separate from verified news, deduplicates repeated narratives, maps U.S.-listed tickers, exposes source diversity/mention velocity/corroboration/materiality, and applies source-specific AI eligibility. Community content remains **UNTRUSTED** and cannot directly change deterministic DAY/SWING/LONG Action/Score/Tradeability/Readiness.

**Oil / Energy Context** now uses truthful WTI and Brent reference evidence, trend/change and Brent–WTI spread context, plus U.S.-market relevance through existing energy context. When CL/BZ continuous-contract or roll-adjusted evidence is not actually available, DE.PULSE says so rather than fabricating futures semantics. **Historical Replay / Scenario Validation** adds reusable CPI/FOMC shock, earnings-gap, high-VIX and market-dislocation scenarios built only from retained canonical evidence and replayed through the same cutoff-safe deterministic engines with no lookahead. It is a validation tool, never paper trading. `De-Pulse.app` and `PersonalMarketTerminal` continuity are preserved.

## v16.8.1 Stable — US Market Scope & Data Utility Hardening

DE.PULSE now processes actionable per-symbol intelligence only for **U.S.-listed stocks and ETFs**. Foreign exchanges remain unavailable for watchlists/desks/Research/Discovery processing, while selected international macro/market evidence remains available as **Global Context** for understanding U.S. conditions. Economic Calendar defaults to **US** and exposes a separate **Global Context** scope. SPY/QQQ and the actively selected ticker receive market-critical live priority. Deterministic DAY/SWING/LONG scores are unchanged.


## v16.8.0 Stable — Professional Context & Readiness Closure

DE.PULSE v16.8.0 Stable adds professional Heat Map modes, 10-year SPY/QQQ Seasonality context, quality-gated structural GEX, selective Smart Notifications and richer Liquidity/Slippage context. These are context/readiness aids: deterministic DAY/SWING/LONG scores remain unchanged. `De-Pulse.app` continues to use the canonical `PersonalMarketTerminal` settings and API keys.


> **v16.7.0 Stable — Core Market Decision Intelligence Closure:** the original professional-roadmap acceptance depth is now fully closed for **#3 Economic Calendar, #12 Market Tradeability, #13 Breadth & Internals, #14 Relative Strength / Weakness and #15 Sector / Industry Regime**. The independent original-roadmap status moves from 17 FULL / 13 PARTIAL to **22 FULL / 8 PARTIAL / 0 MISSING**.

> **What changed in the app:** Market Intelligence now provides a full Economic Calendar filter bar for impact, region, category, date and sort order; Tradeability shows the canonical component breakdown; Breadth shows tracked-universe 20/50/200-day MA participation, 20-session high/low counts and sector participation; selected-ticker Relative Strength includes DAY and SWING comparisons; Sector/Industry Regime shows momentum, relative strength and mapped-member participation/coverage. Discovery also uses shared Day/Swing relative-strength evidence for ranking without additional per-symbol provider fetches.

> **Truth boundary:** Breadth internals are explicitly **tracked-universe** evidence, not a claim of exchange-wide NYSE/Nasdaq A/D breadth. Missing or stale evidence degrades/withholds context. Market Tradeability and Relative Strength remain contextual and **never change deterministic Day/Swing/Long Setup Score/Action formulas**.

> **Stable continuity:** `De-Pulse.app` continues to use `PersonalMarketTerminal`; compatible watchlists, settings and API keys remain preserved.

## Core User Guide

> **v16.6.0 Stable — Full Professional Integration & Acceptance:** all 30 professional requirements are reconciled end to end. The confirmed Master Symbol Store defect is fixed: you can remove the final user-managed ticker or use **Remove All**, and the empty DAY/SWING/LONG store stays empty after refresh/restart. Maintenance now distinguishes **User Master Symbol Store** from protected **System Market Context**; SPY/QQQ and required context feeds remaining available does not mean your user store failed to clear.

> **v16.6 authority:** deterministic Day/Swing/Long Score/Action formulas are unchanged. This is an integration/acceptance release, not a new feature expansion; no v16.7 scope is automatically created. `De-Pulse.app`, `PersonalMarketTerminal`, compatible saved API keys/settings and the existing untrusted AI/community safety boundaries remain intact.

> **v16.5.0 Stable — Context & Alternative Intelligence:** Market Intelligence now adds a transparent **Sentiment Composite**, explicit 11-sector **Market / Sector Heat Map**, defensible **GEX context**, **UNTRUSTED COMMUNITY INTELLIGENCE**, and official-source **Oil / Energy Context**. These are context layers only; deterministic Day/Swing/Long Action/Score remains unchanged.

> **Truth & coverage:** Missing sentiment components reduce confidence and never become neutral zeros. Heat-map cells expose universe/coverage/freshness and stale/missing cells do not vote. GEX is shown only when gamma + open-interest coverage and recency are sufficient; it is a structural proxy, not measured dealer positioning.

> **Community & energy safety:** Community evidence is entered manually by you, is always labeled **UNTRUSTED COMMUNITY INTELLIGENCE**, and is never treated as verified fact or a deterministic input. Oil/Energy keeps official EIA WTI/inventory evidence distinct from tradable `USO`/`XLE`; `USO` is a futures-based ETF, not WTI spot, so basis/curve/roll differences remain explicit. `De-Pulse.app` and existing compatible API keys/settings in `PersonalMarketTerminal` are preserved.

> **v16.4.0 Stable — Professional Research AI:** Research AI is now a structured, evidence-grounded **second opinion**. Results separate Bull / Base / Bear cases, contradictions, missing evidence, risks and catalysts, and expose the canonical Evidence Package / Evidence Snapshot identity used for the analysis. AI confidence is confidence in the analysis/evidence completeness — **Setup Score is not win probability**.

> **External-content safety:** News, SEC filings, provider/web text and other external content are treated as **untrusted data**, never instructions. Instruction-like text is isolated and may raise a safety warning. AI may cite only evidence IDs supplied by DE.PULSE and cannot fabricate missing evidence, freshness, provider agreement or causation.

> **Routing & cost control:** Settings → Research AI offers Manual, Efficient, Balanced and Deep routing. Manual uses only your selected configured provider. Automatic policies use only configured providers with bounded output tokens. Material Evidence Package caching prevents repeat paid/deep calls when only ordinary quote ticks changed.

> **Authority boundary:** AI cannot trade, change deterministic Action/Score, Market Tradeability, Trade Readiness, execution logic, or desk/watchlist membership. It remains advisory.

> **Stable app standard:** the macOS application is **`De-Pulse.app`**. Release version/channel belongs in build metadata, not the app filename. Stable runtime remains `~/Library/Application Support/PersonalMarketTerminal`, so compatible saved API keys/settings continue forward.

> **v16.3.0 Stable — Professional Validation & Learning:** DE.PULSE now separates **what it knew at decision time** from **what happened afterward**. Research can replay a frozen v16.3 decision through the same deterministic Day/Swing/Long scoring primitive, then show outcome evidence separately. Signal Outcome Measurement tracks Entry Zone touch, target/invalidation ordering, elapsed time, MFE/MAE and supported 1D/3D/5D/10D returns. A target reached before Entry Zone is never counted as success; unresolved coarse-bar ordering is `AMBIGUOUS`.

> **Stable upgrade continuity:** v16.3.0 Stable uses the same canonical `PersonalMarketTerminal` settings directory as v15.1.2 and earlier Stable builds. Compatible saved settings, watchlists and API keys are reused automatically; you do not need to re-enter preserved credentials. The installer/launch command removes only the stale instance lock, not `state.json` or `secrets.json`.

> **Where to use it:** Research shows Historical Replay, Calibration, SPY/QQQ Seasonality and selected-ticker correlation context. Market Intelligence shows only material market-level Seasonality and Candidate Concentration. Decision Queue may raise **attention priority** when candidates are highly correlated, but opportunity quality and deterministic Score/Action are unchanged. Maintenance shows sample/replay/concentration health.

> **Truth guardrails:** Seasonality and calibration are descriptive, not trading instructions. **Setup Score is not win probability.** Small samples are `INSUFFICIENT`; missing history is `UNAVAILABLE`. Replay is `EXACT` only when the frozen v16.3 evidence/features/settings identity is sufficient; legacy snapshots degrade to `SCENARIO` or `PARTIAL`. v16.1 Market Intelligence and v16.2 Event Intelligence remain canonical inputs and are not duplicated.

> **v16.2.0 TEST — Professional Event Intelligence:** Market Intelligence now activates its EVENTS section using one derived Event Intelligence snapshot over the existing canonical News, Macro, Earnings, SEC/Catalyst Watch and Event Mode reaction stores. News is clustered/deduplicated across supporting sources, classified by category/materiality and publication-time validity. The Economic Calendar shows sourced timing/impact/lifecycle plus Actual/Forecast/Previous and Surprise only when evidence exists. FOMC/Fed intelligence uses sourced Federal Reserve calendar evidence only. Smart Notifications are selective event/state-change items (material news, event-window entry, material reactions), not a restored generic Alerts tab. Reaction Intelligence reuses Catalyst Watch and Event Mode capture. Event → Decision correlation may raise Trade Readiness risk and Decision Queue attention but never changes deterministic Day/Swing/Long Score/Action.

> **How to use v16.2:** Market Intelligence → EVENTS is the market-level event view. Research → Catalysts contains ticker-specific Event Intelligence. Maintenance shows Event Intelligence source/job health. Trade Readiness and Decision Queue expose event risk when applicable. If event/calendar evidence is stale, missing or clock-invalid, DE.PULSE degrades/marks unavailable instead of inventing a clean event state.

> **v16.1 foundation retained:** Market Structure, Market Tradeability, tracked-universe Breadth, Relative Strength, Sector/Industry Regime and Liquidity/Slippage remain canonical and unchanged except for consuming event context through Readiness/Queue. v16.1.1 exact weekly/freshness/spread-evidence truth hardening remains permanent.


> **v16.0.6:** Final v16.0 consolidation/delivery closure. Carries all v16.0.5 professional-truth fixes, keeps the Stable-style primary header, restores the normal native app launch path with binaries prepackaged, and makes the complete release/handoff artifact set mandatory. No v16.1 Market Intelligence/Tradeability feature is included. Deterministic Day/Swing/Long Score/Action formulas remain frozen; v15.1.2 remains the protected Stable baseline.

> **v15.1.2 Stable:** Research v2, detailed SEC & Ownership insider transactions, reliable DAY/SWING/LONG membership, Master Market Symbols add-to-all, Data Freshness v2 with adaptive Check/Data Age and automatic recovery, and hardened Pre-Market/Market Open/Catalyst readiness workflows. Stable requires all 48 approved items and every blocking release gate to pass.


> **v15.0.1 Stable:** Central Provider Router and Data Freshness diagnostics, truthful VIX routing, synchronized DAY/SWING/LONG membership, Master Market Symbols controls, restored SEC Form 4 BUY/SELL/OTHER visibility, and a blocking all-tab layout/hover/scroll integrity audit. Stable status requires all 53 approved scoped items and every blocking release gate to pass.

> **v14.3.7 Stable:** Font-fit audit across every rendered tab/surface and supported viewport. Only overflowing compact text was reduced to fit; no content, placement, ordering, trading logic, or workflow was removed or changed.


> **v14.3.6 Stable:** Exact v14.3.5 stability baseline plus SEC Form 4 BUY/SELL/OTHER clarity, side readiness-dot relocation, Options settings-card spacing cleanup, canonical DAY/SWING/LONG membership pills, SPY/QQQ Market Regime pills, a full Discovery overlap audit, and corrected Macro Rates health using Treasury core rates with optional FRED enrichment. Deterministic Setup Scores are unchanged.

> **v14.3.4 Stable:** UI stability patch: cleaner Data Engine action alignment, smaller readiness-state typography, latest-five QA history with complete metadata, Decision Queue Why? expanded by default, content-driven cache-policy cards, robust header notification fitting, Settings-header Auto-Start control, and Maintenance VIX manual refresh. Deterministic scores are unchanged.


> **v14.3.3 Stable:** Adds priority-based multi-feed live allocation, canonical Live Coverage State, persistent manual-action completion feedback, Maintenance Data Engine Detail controls, adaptive, unclipped Data Engine row geometry/status typography, content/copy consistency auditing, and canonical build identity. GLD/SLV/USO remain pinned tradable live symbols. The side Data Engine keeps only the three trading-readiness controls. Deterministic Setup Scores remain unchanged.

> **v14.3.2 Completion Final:** Keeps the v14.3.1 provider/readiness completion fixes, moves transient notifications into the existing global header center lane with no header-height change/overlap, and adds small Data Engine manual Run/Refresh/Evaluate controls plus a conditional Primary Live Stream Reconnect that reuses the existing production reconnect loop. Setup Scores remain unchanged.

> **v14.3.1 Completion Patch:** Closes the v14.3.0 audit: full Market Open reconciliation, actual entry-zone validation, ARMED→TRIGGERED→REACTION catalyst phases, truthful capability verification, Twelve Data fallback-validation, slower official-macro cadence, and removal of the sidebar black scroll-fade/jitter. Setup Scores remain unchanged.

> **v14.3.0 Improvement Build:** Built from the exact v14.2.2 Stable source. Adds Market Open Prep, event-driven Earnings & Material Catalyst Reaction Watch, P1/P2 provider intelligence, canonical derived macro/liquidity states, initial notification/Save-scroll preservation work (the floating-toast approach is superseded by the v14.3.2 header slot), and audited UI placement. Deterministic Day/Swing/Long Setup Scores remain unchanged.


DE.PULSE is a personal market-intelligence terminal for discovering stocks, evaluating Day/Swing/Long-Term setups, researching evidence, and understanding market context. It is not a brokerage or order-entry system. The terminal keeps deterministic trading-desk analysis separate from optional Research and AI confirmation so you can see exactly which layer produced each conclusion.

> **v14.2 Stable Hotfix:** Built from the v14.0.3 Stable baseline. v14.2 makes three conservative UI corrections only: centered/no-gap transient notifications, less forced empty panel space, and the simpler sidebar Data Engine without a Details row. The v14.0.3 handoff-compliance and cross-module architecture remains the active product baseline. It completes the approved Global/Macro/Options architecture, deepens horizon-specific Trade Readiness and Signal Validation, fixes persistent expanders/sidebar layout/Options styling/build identity, and keeps the validated deterministic Day/Swing/Long Score/Action formulas unchanged.


## v14.0.3: Global, Macro, Options and Trade Readiness

v14 keeps the trading desks deterministic while adding shared context that is interpreted differently by horizon. Dashboard stays concise; deeper evidence is progressively disclosed in the desk, Research, AI, or Maintenance rather than copied everywhere.

### Global Market Drivers
DE.PULSE correlates U.S. risk, Asia, Europe, rates, FX, credit, commodities, volatility, sector participation and semiconductor read-through. In AUTO, a configured direct underlying/index feed is preferred, official/public completed-session evidence is retained separately where robustly available, and truthful real ETF proxies remain current fallback/context. TWSE TAIEX has an official public close adapter; other markets use direct-provider data when configured or clearly labeled proxies when no robust public underlying feed is available. A proxy is never presented as the underlying index. Completed official closes remain valid session evidence rather than being marked stale solely because the local exchange is closed.

Broad-market breadth is calculated from a broad/index/sector universe when available. Your personal watchlist is never presented as overall market breadth.

### Macro events and High-Impact Event Mode
Official/public Fed, BLS, BEA, ECB, BOJ and China calendars feed one shared macro-event layer. If the official source supplies only a date, DE.PULSE keeps it date-only and does not invent a release time. A known HIGH-impact event can activate preparation around T−15, prioritize market data and capture deterministic post-release reaction windows. AI is not on this critical path.

### Options Intelligence
Options are contextual intelligence, not an Options Trading Desk. DE.PULSE can use real option snapshots to show call/put volume balance, IV, and a nearest-expiry near-the-money expected move when the source contains enough data. Unsupported fields remain unavailable; they are never guessed.

- **Day:** immediate unusual/volume/IV/event confirmation or conflict.
- **Swing:** persistent positioning, IV/expected-move and catalyst context.
- **Long-Term:** only material longer-horizon positioning/risk; short-lived flow is intentionally hidden.

Options context can inform Dashboard, Decision Queue, Trade Readiness, Research, AI evidence and Signal Validation, but it never silently changes the deterministic desk Score or Action.

### Trade Readiness
Trade Readiness is a separate context layer: `READY`, `CONDITIONAL`, `CAUTION`, or `INCOMPLETE`. It considers horizon-appropriate data quality, event risk, market/global alignment, options context and contradictions. A BUY score can remain unchanged while Readiness becomes CAUTION because, for example, CPI is imminent or options/global context conflicts.

### Signal Validation
DE.PULSE can save de-duplicated decision snapshots and later compare them with real historical outcomes (1D/3D/5D/10D, MFE and MAE). This validates whether context adds value. Signal Validation never reweights or changes formulas automatically.

### Where configuration and diagnostics live
Settings contains provider configuration and connection tests. Maintenance contains operational Data Capabilities, current source/mode/freshness, Signal Validation, provider health and QA history. See **Capabilities & Limitations** for what is active, degraded, proxy-backed, unavailable, or ready for a premium upgrade.

The sidebar Data Engine is a detailed live-health surface inside the scrollable sidebar. It keeps provider/runtime/data-health rows and status dots visible directly at normal desktop widths; deeper troubleshooting remains in Maintenance. Expandable panels elsewhere preserve their open/closed state across live data rerenders. Maintenance also verifies that the running backend version/build matches the packaged renderer so stale build labels cannot silently ship.

## Getting Started

Start the runtime from the header. Use Demo Mode to learn the interface without provider credentials. Use Live Mode after configuring Finnhub and Alpaca in Settings. DE.PULSE restores your watchlists, selected ticker, settings, and persisted market/research cache between launches.

A practical first workflow is:

- Open Dashboard and read Market Context plus the Decision Queue.
- Open Discovery when you want new candidates.
- Open the recommended Day, Swing, or Long-Term desk for the actual setup.
- Use Research only when you want deeper company, catalyst, filing, or AI evidence.
- Return to the same desk and symbol after Research.

Research is optional. A stock does not need AI approval before the desk can analyze it.

## Understanding Market Sessions

DE.PULSE uses Eastern Time for US market-session logic and also displays Pacific Time where operationally useful.

- Overnight: approximately 8:00 PM–4:00 AM ET when supported by the configured Alpaca overnight mode.
- Pre-market: 4:00–9:30 AM ET.
- Regular market: 9:30 AM–4:00 PM ET.
- After-hours: 4:00–8:00 PM ET.
- Weekend/closed periods are shown explicitly.

Day Trade is session-aware. Overnight or pre-market movement is useful context before the opening bell; once regular trading begins, DE.PULSE emphasizes the regular-session open gap, VWAP, participation, liquidity, and live momentum rather than pretending overnight and regular trading are the same session.

## Data-State Labels

Always read the data-state label with the number.

- LIVE means a current live market stream is being used for that session.
- INDICATIVE means the provider has supplied an overnight/reference value that should not be treated as an exchange-grade regular-session trade.
- SNAPSHOT means the latest value came from a provider snapshot rather than an active stream.
- CURRENT means a non-stream dataset such as news/fundamentals is within its intended freshness window.
- CACHED means persisted data is being reused while a newer provider refresh may occur.
- STALE means the last successful update is older than the expected freshness policy.
- RECONNECTING or DEGRADED means the runtime is operating with reduced provider quality or attempting to restore a preferred feed.

DE.PULSE preserves provenance so a disk-restored value is not mislabeled LIVE.

## Header and Market Instruments

The header shows the market session, data health, runtime controls, and common market instruments. SPY, QQQ, DIA, IWM, GLD, SLV, TLT, USO, and XLK use the same canonical Master Symbol Store as watchlist symbols. VIX remains on a dedicated special-index path because it is not treated as an ordinary equity ticker.

Market Instruments are context, not a separate portfolio. If SPY appears in the header and several desks, DE.PULSE reuses one canonical SPY record and subscription lifecycle.

## Dashboard

Dashboard answers one question: **What deserves attention now?**

The priority order is:

1. **Market Regime + Horizon Lens** — contextual overall regime plus Day/Swing/Long lenses.
2. **Decision Queue** — the centerpiece: Action/Score, priority, Readiness, Research state, Key Driver and Risk/Conflict.
3. **Catalyst / Event Risk** — fresh News (≤72h), Filings (≤30d) and Earnings (±14d), with View All.
4. **Horizon Overview** — one best setup per horizon plus candidate count and engine state.
5. **Options Intelligence** — shown when material.
6. **More Market Context** — always visible lower on the Dashboard with global, macro and reaction detail; it is not an expand/collapse disclosure.

### Decision Queue

The Decision Queue is intentionally not a portfolio or trade blotter. It combines useful candidates from Day, Swing, and Long-Term while avoiding meaningless duplicate rows. Each row tells you which desk is relevant, the existing deterministic setup score/state, event risk, and data freshness.

Use **Go to Desk** when you already understand the company. Use **Research** when you need deeper evidence first.

## Discovery

Discovery answers: **What new symbols should I investigate?**

Choose Day Trade, Swing, or Long-Term mode. The candidate finder ranks symbols for that horizon and labels the result **Discovery Rank**. Discovery Rank is not the same as the desk Action/Score.

The funnel is:

- Scanned: universe considered.
- Ranked: candidates with a computed Discovery priority.
- Qualified: candidates passing your current filters.
- Staged: symbols you saved for further investigation.

From a candidate you can:

- Stage it for later.
- Add it directly to the current horizon desk.
- Open Research.

A scanner-only symbol may not yet have all desk-grade history/fundamentals. Staging or adding it hydrates the canonical symbol state without starting a duplicate quote subscription.

## Day Trade Desk

Day Trade answers: **Is there a valid setup for the current trading session?**

The existing Action, Score, Entry, Target, Invalidation, Quality, and plan calculations remain the deterministic decision model.

### Must-have Day context

The first visible context layer is deliberately short:

- Session and data provenance.
- Overnight/Pre-Market Move before the open, or Open Gap once regular trading is underway.
- VWAP relationship.
- RVOL.
- Bid/ask spread.
- Dollar volume.
- ATR %.

These inputs describe participation, liquidity, volatility, and where price sits in the current session. They do not modify the deterministic Score in this test build.

### More Day context

Open **More context** only when needed:

- RSI.
- EMA 9 / EMA 21.
- Support / Resistance.

The optional section is hidden by default to reduce indicator overload.

### Session carryover

Before 9:30 AM ET, overnight/pre-market movement helps explain whether a symbol is arriving at the open with unusual momentum. At the regular open, the context changes: the opening gap and regular-session VWAP become more important. DE.PULSE does not reset your watchlist or historical evidence, but the visible session metric is reweighted for the current session.

## Swing Desk

Swing answers: **Is there a technically credible 5–20 session setup?**

### Must-have Swing context

- 20 / 50 / 200 trend structure.
- Relative Strength vs SPY.
- ATR %.
- RVOL.
- Support / Resistance.
- Earnings proximity.

These cover trend, market-relative performance, volatility, participation, key levels, and event risk without importing Day-only spread/VWAP noise.

### More Swing context

- Daily RSI.

RSI is intentionally secondary. It can help interpret extension or pullback state, but it should not dominate a trend/relative-strength swing decision.

## Long-Term Desk

Long-Term answers: **Does this symbol combine a credible long-horizon setup with acceptable business/fundamental evidence?**

The deterministic Long Action/Score remains unchanged in this test build. The new context makes the supporting evidence clearer.

### Must-have Long-Term context

- Revenue Growth.
- EPS Growth.
- Operating Margin.
- Free Cash Flow.
- Debt / Equity.
- P/E / Forward P/E.
- SMA 50 / SMA 200.
- Relative Strength vs SPY.

### More Long-Term context

- ROE.
- 52-week position.
- Current Ratio.
- Dividend Yield.

These remain accessible without forcing every long-term row to become a wall of ratios. Intraday spread, VWAP, gap, and short-term EMA metrics are intentionally not part of Long-Term context.

## Watchlists and Adding Symbols

Day, Swing, and Long-Term each have a permanent watchlist. The compact list is for fast scanning; the detailed setup table is for Action/Score and planning zones. A dedicated empty Add Symbol row remains at the bottom of the detailed table.

Clicking a symbol in the compact or detailed list selects it and moves to its detailed section. Switching between desk tabs starts at the top of the newly selected desk; selecting a ticker is the deliberate exception that can scroll to that symbol.

The same symbol may exist in multiple desk watchlists without creating multiple provider subscriptions.

## Engine On / Off

Each desk engine can be enabled or disabled independently.

When a desk engine is OFF:

- Its watchlist is preserved.
- Shared market data can continue because other pages/desks may need it.
- That horizon's analysis is paused.
- No duplicate subscription is created or removed merely because a second desk uses the symbol.

When turned ON:

- Missing required history is requested immediately.
- Relevant derived cache entries become eligible for recomputation.
- Analysis resumes using the shared canonical data.

## Action, Score, Quality, and Planning Zones

The deterministic trading desk layer includes:

- Action / state.
- Score.
- Entry range.
- Target range.
- Invalidation.
- Risk/reward context.
- Setup Quality.
- Earnings/event risk.

In v13.3.0 Stable, the additional Day/Swing/Long context and AI Research Confirmation remain separate. A favorable AI opinion does not silently turn a desk score from 70 into 85.

## Research Workspace

Research answers: **What deeper evidence should I know about this symbol before I trust or reject the setup?**

The workspace keeps four focused primary tabs:

- Overview.
- Fundamentals.
- Catalysts.
- Filings.

Technical setup indicators are not repeated as a separate Research tab. The appropriate Day/Swing/Long desk remains the canonical technical setup view.

### Research Overview

Overview shows:

- Symbol/price/freshness snapshot.
- Next earnings, mapped catalyst count, SEC count, and Discovery Rank where available.
- Day, Swing, and Long horizon cards with Trade Readiness.
- Sourced Earnings/Guidance, SEC/Filing Risk, Global/Macro and material Options evidence before optional AI confirmation.
- AI Research Confirmation when requested, as a second opinion after sourced evidence.

Use each horizon card's **Open Desk** action to open the relevant desk while preserving the selected symbol.

### Fundamentals

Fundamentals are arranged vertically as **Growth & Business Quality → Cash Flow & Balance Sheet → Valuation → Supporting Fundamentals**. This is especially useful when evaluating Long-Term candidates but remains available for any symbol.

### Catalysts

Catalysts combines selected-symbol news and earnings context. Only a short list is shown by default and headline summaries are collapsed. Use **All News** or **All Earnings** when you need the wider secondary view.

### Filings

Filings shows selected-symbol SEC filings and the derived SEC intelligence summary. Use **All Filings** for the wider tracked-symbol view.

## Earnings and SEC Intelligence by Horizon

Earnings and filings are interpreted according to the trading horizon instead of being copied into every desk.

### Day Trade

Day Trade shows earnings or SEC context only when it is a current-session catalyst. Before earnings it can show the expectation and event timing. After results are available it can show EPS/revenue expected vs actual, surprise %, guidance state when explicitly supported by evidence, and the overnight/pre-market reaction. A fresh material filing can appear as a catalyst; old routine filings do not occupy Day screen space.

### Swing

Swing uses compact earnings/guidance and filing-risk context: recent beat/miss, guidance Raised/Maintained/Lowered/Withdrawn/Not Provided, whether an earnings gap is holding/fading when data supports it, next-earnings proximity, and material filing/event risk. Full filing history stays in Research.

### Long-Term

Long-Term keeps the full SEC Intelligence Summary and deeper earnings trend because multi-quarter business evidence matters at this horizon. Friendly filing names are primary, for example **Quarterly Report · 10-Q**, **Annual Report · 10-K**, **Material Company Update · 8-K**, and **Insider Transaction · Form 4**. Raw form codes are secondary reference.

Guidance is never invented. If the available evidence does not explicitly support a guidance change, DE.PULSE shows **NOT PROVIDED** rather than inferring a value.

## AI Copilot and Research Confirmation

Ask AI is an optional second opinion. Its visible response is intentionally structured rather than a long essay.

The Research Target is explicit. Ask AI from Dashboard, Discovery, Day, Swing, Long-Term, or Research passes the ticker and the relevant horizon context automatically. You can also select a known symbol or enter any ticker directly. You do not have to restate the setup in a blank chat.

The result contains:

- Verdict: FAVORABLE, MIXED, CAUTION, or INFORMATIONAL.
- Confidence.
- Up to three key reasons.
- Up to three risks.
- One catalyst.
- Best-fit horizon.
- One user-controlled next action.
- A one-line summary.
- Expandable Evidence & details.

### What happens if AI is favorable?

Nothing automatic happens.

You may explicitly choose:

- Open Research.
- Open the recommended Day/Swing/Long desk.
- Add the symbol to the recommended desk if it is not already there.

AI cannot place an order, change your deterministic desk Score, or add a symbol without your click. Its role is **Research Confirmation**, not autonomous decision execution.

## Secondary News, Earnings, and SEC Views

News, Earnings, and SEC remain available, but they are no longer primary sidebar destinations. This reduces navigation duplication.

Open them from:

- Dashboard Catalyst Pulse.
- Research Catalysts.
- Research Filings.

The symbol-focused Research workspace should be the normal deep-dive path. The secondary views are useful when you want the full tracked-universe list.

## Data Engine

The Data Engine sidebar is the detailed live health surface. It reports central service state rather than running separate diagnostic subscriptions.

Status text is never hidden to force a fixed row height. Frequently changing rows reserve a comfortable minimum height to reduce visual jitter; simpler rows stay compact naturally, and longer real operational messages expand to remain fully readable.

Important services include:

- Finnhub WebSocket.
- Finnhub REST.
- Alpaca IEX / overnight stream.
- Alpaca snapshot.
- Historical bars.
- Fundamentals.
- News.
- Earnings.
- SEC Filings.
- Cache Refresh.
- Scanner.
- VIX.
- AI.

Long values wrap rather than hiding important information.

## Settings

Settings is for configuration, not diagnostics.

### Configuration order

Settings starts with a **Configuration Summary**, then follows the setup path: **Core Market Data → Market Data Modes → Global/Macro/Options → SEC/Filing Configuration → Trading Engines → Signal Engine → AI Integration → Application**.

### Data Connections

Configure Finnhub and Alpaca credentials first. Save/Test confirms credentials without changing your watchlists. Global/macro/options provider configuration and SEC contact identity remain with the data-source settings rather than being separated into unrelated lower sections.

### Market Data

Choose Data Mode and Overnight Data Mode. Dropdowns always show a visible chevron so they cannot be confused with plain text fields.

### Trading Engines

Enable/disable Day, Swing, and Long analysis.

### Signal Engine

Signal Profile, Market Context Strength, and Earnings Risk Penalty remain configuration inputs for the existing deterministic model.

### AI Integration

Choose Groq, OpenRouter, or Gemini and configure the relevant key/model. Provider selection affects AI Research Confirmation only; it does not change deterministic desk calculations.

### Application

Auto-Start Runtime is a compact right-aligned control in the Settings header; **QUIT DE.PULSE** remains in the Application lifecycle section. The persistent Save Settings bar remains accessible while scrolling and must not cover page content.

## Cache Freshness and Performance

DE.PULSE caches data that does not need to be refetched continuously. The goal is not "cache forever"; it is "refresh according to how quickly the information can change."

The v13.3.0 Stable policy is:

- Live quotes/trades: event-driven streaming with controlled fallback.
- Provider REST snapshots: normal provider fallback cadence.
- Historical bars: 15-minute normal refresh plus immediate hydration when required.
- News: approximately every 10 minutes during active market sessions and approximately every 30 minutes overnight/closed.
- Earnings calendar: approximately every 2 hours, with recent reported quarters retained so expected-vs-actual can be compared.
- SEC filings: approximately every 30 minutes.
- Fundamentals: approximately every 24 hours, with refresh also becoming eligible when new financial evidence makes the cached fundamentals stale.
- Disk cache: dirty-state/fingerprint-aware persistence; active changes are persisted promptly, while unchanged cache content is not rewritten just because a timer fired.
- Pre-Market Preparation: once per trading day during approximately 3:15–3:50 AM ET.
- Weekly Integrity: Saturday, non-destructive.

### Pre-Market Preparation does not stop overnight quotes

US equities can still have overnight data before 4:00 AM ET. DE.PULSE therefore does not stop or restart quote streams. Pre-Market Preparation is a low-priority targeted refresh of stale/missing historical bars, fundamentals, earnings, news, SEC, and derived prerequisites. Existing good cache entries remain available while refreshes run.

### Stale-while-refresh behavior

If a stable dataset is available from cache, DE.PULSE can keep displaying it with its CACHED/STALE provenance while a newer refresh is attempted. A provider failure does not require deleting a known-good cache.

## Maintenance

Maintenance is for inspect/diagnose/clean/verify/repair tasks.

### Weekly Maintenance

Weekly Maintenance is diagnostic and non-destructive. It does not clear the market cache or change trading subscriptions.

### Freshness & Cache Policy

Maintenance shows Data Freshness, the latest Pre-Market Preparation status, Weekly Integrity status, and the expected cadence for each data class.

Use **Refresh Stale Data Now** to request a low-priority stable-data refresh. Use **Run Pre-Market Prep Now** to run the same targeted preparation outside its scheduled window. Use **Run Integrity Check** for the non-destructive integrity pass. All three keep live quote streams active.

Use **Clear Market Cache** only when you intentionally want destructive market-cache removal. The app preserves settings/watchlists and immediately requests required historical bars again when the Live runtime is active.

## Common Workflows

### New stock found in Discovery

- Open Discovery.
- Choose a horizon and run the scanner.
- Read Discovery Rank and the concise reason tags.
- If you understand the company, add/open the recommended desk first.
- If the name is unfamiliar or catalyst-heavy, open Research first.
- Review Fundamentals/Catalysts/Filings as needed.
- Optionally Ask AI for a concise Research Confirmation.
- Explicitly open the recommended desk.
- Make your own decision from the desk Action/Score plus supporting evidence.

### Stock you already know

- Open Dashboard Decision Queue or the relevant desk.
- Select the ticker.
- Read the deterministic setup and must-have horizon context.
- Open **More context** only if required.
- Open Research only when you need deeper evidence.
- Return to the same desk and symbol.

## Common Troubleshooting

### Price is visible but chart is empty

Historical bars may still be hydrating. Adding a symbol, re-enabling a desk engine, starting Live runtime with missing bars, or clearing the market cache requests required bars immediately rather than waiting for the normal 15-minute history cycle.

### Quote is not marked LIVE

Read the provider/source and session. Overnight may be INDICATIVE depending on entitlement and selected overnight mode. A disk-restored quote is marked CACHED until a provider replaces it.

### Day context shows Overnight Move instead of Open Gap

That is intentional before the regular open. During the regular or after-hours session, the context switches to the regular-session opening gap.

### AI gave a favorable result but nothing changed

That is intentional. AI is non-binding. Use the explicit **Open Desk** or **Add to Desk** action if you agree with the research.

### Research has fewer technical indicators than before

That is intentional. Technical setup evidence belongs to Day/Swing/Long. Research is for company, catalyst, filing, and optional AI evidence.

### Data looks old

Check the Data Engine freshness timestamp and Maintenance > Freshness & Cache Policy. You can use **Refresh Stale Data Now** without clearing the whole cache.

## Data and Privacy Notes

API keys are stored locally by the app's existing credential mechanism and are intentionally excluded from exported profiles. Research/AI prompts use the currently selected DE.PULSE context; do not treat model output as guaranteed fact. Verify material filings, earnings dates, and provider freshness before acting on them.

> **v14.0.3 Stable:** Global/macro/options context, Trade Readiness, Signal Validation, Decision Queue/Research/AI correlation, and the v13.3 horizon architecture are part of the Stable product. The validated deterministic Action/Score model remains unchanged.


### Options IV change
When two consecutive real option snapshots are available, DE.PULSE shows **ΔIV** as the percentage-point change in average chain IV. It is context only; it does not change deterministic Action/Score.

## v14.0.3 UI consistency patch

This patch does not change deterministic trading Scores/Actions. It improves readability and stability in shared decision surfaces: Global/Macro Driver cards no longer collide, Market Instruments no longer pause/jiggle on hover, Event & Data Context uses a stable responsive key/value layout, SEC Intelligence values no longer overpower labels, and very old cached quote context is visually identified as stale evidence rather than current data.

## v14.3.0 Improvement Build

### Preparation and event-driven reaction flow
- **Pre-Market Prep** prepares due/stale supporting data before premarket without indiscriminately clearing good cache.
- **Earnings & Material Catalyst Reaction Watch** is event-driven. It activates only for a relevant scheduled earnings event in its applicable session/date or for a material unexpected News/SEC trigger. Routine headlines do not create a dedicated watcher run.
- **Market Open Prep** runs shortly before the U.S. regular open on trading days and reconciles the latest overnight/premarket evidence into trading-readiness flags. It is a checkpoint, not a frozen snapshot.
- Side-panel Data Engine and Maintenance expose structured state for these jobs; trading tabs receive the resulting context rather than duplicated operational telemetry.

### New contextual intelligence
DE.PULSE can now derive concise **Rates, Credit, Financial Conditions, Dollar, Inflation, Labor, Energy**, and ticker-level **Liquidity Health** states from configured real providers. These states are contextual evidence. They do not silently change deterministic Setup Scores.

### Save behavior and notifications
Save/action feedback appears in the global header notification slot. The header remains the same height; long feedback truncates safely and saving from a long Settings page preserves the current scroll/focus context instead of returning the user to the top.

### Dashboard More Market Context
More Market Context remains always visible. The evidence is grouped into Market Direction & Risk, Breadth & Leadership, Rates & Credit, Global Markets, and Commodities & Inflation. Spacing is balanced to remove hollow hero cards and cramped subsection transitions without compacting the overall interface.


## v14.3.6 Stable usability patch

- **SEC insider filings:** Form 4 entries now say Insider BUY, Insider SELL, Insider OTHER, or Mixed when the filing contains multiple transaction classes. BUY/SELL is reserved for genuine P/S transaction codes; grants, exercises, tax withholding, gifts, conversions and similar rows remain OTHER.
- **Desk membership:** Day, Swing and Long watchlist membership is shown with small shared DAY / SWING / LONG pills in desk watchlists and Discovery. A ticker can belong to more than one desk at the same time.
- **Market Regime at a glance:** SPY and QQQ show current/last real price plus session change. A compact directional state-only pill appears between QQQ and Data Confidence. It is derived from the existing horizon Market Regime evidence and does not alter setup scores.
- **Global Market Drivers:** the right-side state word now includes a small status light and summarizes the combined driver-family evidence. The evidence line shows supportive/headwind counts and confidence without repeating the state.
- **Catalyst Watch:** RUNNING means the short manual Evaluate action only. Normal event-driven states are READY, ARMED, TRIGGERED and REACTION.
- **Macro Rates:** core official rates can remain healthy when usable Treasury data exists even if optional FRED enrichment is temporarily unavailable; provider-specific health remains explicit in Maintenance.
- **Discovery and placement:** candidate rows use explicit columns so liquidity, spread, membership and action controls do not overlap. At narrower supported widths, the candidate table scrolls inside its panel instead of causing page-level overflow.

## v18.3 STABLE — hosted backup/recovery behavior

v18.3 STABLE includes the certified private, integrity-checked persistence backup/migration format for hosted-state operation. It preserves account/session hashes, per-user workspaces and canonical market/intelligence history, so archive files must be protected like other sensitive application data. API/provider secrets are not included in this database archive.

Hosted readiness now reflects the canonical database rather than only process startup. A database interruption can leave `/api/health` live while `/api/ready` reports unavailable; DE.PULSE keeps bounded/coalesced persistence work and recovers only after database health is stable again. This is intentional truthful degradation, not a silent fallback to a different local state store.
