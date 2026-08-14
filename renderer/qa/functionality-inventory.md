# DE.PULSE v14.0.3 — Permanent Functionality Inventory

This inventory is a blocking Gate 3 artifact. Every release must exercise the functions below, including normal, loading, empty, degraded and unavailable states where applicable. “Cross-module” means the capability is evaluated across relevant surfaces even when intentionally hidden from a horizon.

## Shell / navigation / lifecycle
- [x] DE.PULSE brand and approved icon assets.
- [x] Market session + next transition, ET/PT clocks, data status.
- [x] START DATA / STOP DATA runtime lifecycle.
- [x] Market Instruments strip and VIX dedicated-index behavior.
- [x] Below-header success/warning/error notification area; hidden state consumes zero layout height.
- [x] Sidebar order: Day → Swing → Long; Discovery → Research → AI; Maintenance → Settings → Documentation.
- [x] Bundled SVG navigation icons, keyboard/accessibility labels and collapsed-sidebar tooltips.
- [x] Compact anchored Data Engine footer; full diagnostics reside in Maintenance.
- [x] Persistent expandable panels across live/runtime rerenders.
- [x] Clean Quit lifecycle and single-instance/local-session security.

## Dashboard
- [x] Market Regime with confidence and compact Day/Swing/Long Horizon Lens.
- [x] Global Market Driver contribution to contextual overall regime without changing deterministic desk formulas.
- [x] Broad-market options contribution only when SPY/QQQ options context is real/reliable; ticker options do not move the whole-market regime.
- [x] Decision Queue is the primary action surface after Market Regime.
- [x] Queue fields: ticker, horizon, Action, Score, Priority, Trade Readiness, Research state, Key Driver, Risk/Conflict.
- [x] Queue Why? correlation: global/macro, options, earnings and event reaction.
- [x] One best setup per Day/Swing/Long horizon + candidate count/engine state.
- [x] Catalyst Pulse with News ≤72h, Filings ≤30d, Earnings ±14d and View All.
- [x] Compact Options Intelligence pulse with P/C, IV, ΔIV, expected move and provenance.
- [x] More Market Context progressive disclosure with direct/official/proxy global evidence, macro horizons and event reaction.

## Day Trade Desk
- [x] Existing deterministic Action / Score / Entry / Target / Invalidation / Quality unchanged.
- [x] Session, overnight/pre-market/open-gap, VWAP, RVOL, spread, dollar volume, ATR%.
- [x] Current catalyst and fresh material SEC event only.
- [x] VIX/rates/FX/global context and High-Impact event countdown/reaction.
- [x] Day-specific options context (P/C, IV/ΔIV, immediate confirmation/conflict).
- [x] Trade Readiness consumes Day data quality, session, VWAP/RVOL/liquidity, catalyst, event/global/options contradictions.
- [x] Optional RSI, EMA 9/21 and support/resistance behind progressive disclosure.

## Swing Desk
- [x] Existing deterministic Action / Score / zones unchanged.
- [x] 20/50/200 trend, RS vs SPY, ATR%, RVOL, support/resistance, earnings proximity.
- [x] Genuine 5D/20D rates, USD, credit, semiconductor/global context.
- [x] Swing options positioning/IV/expected-move context.
- [x] Compact filing/catalyst risk.
- [x] Trade Readiness consumes trend/RS/volume/levels/event/macro/global/options/data-quality context.

## Long-Term Desk
- [x] Existing deterministic Action / Score / zones unchanged.
- [x] Revenue/EPS growth, operating margin, FCF, leverage, valuation, SMA 50/200, RS vs SPY.
- [x] Earnings/guidance trend and full SEC Intelligence.
- [x] Genuine 1M/3M rates, real-yields, credit, USD, China/global-growth context where real history exists.
- [x] Options intentionally hidden unless materially long-dated.
- [x] Trade Readiness consumes fundamentals, guidance, valuation/leverage, technical trend, persistent macro regimes and completeness.

## Global Market Driver Engine
- [x] Shared fetch/calculate-once canonical global context.
- [x] U.S. risk / VIX / breadth / sectors.
- [x] Asia, Europe, rates, FX, credit, commodities and semiconductor read-through.
- [x] Optional Twelve Data direct-global adapter and ES/NQ/RTY direct futures interface.
- [x] Official/public international close path (TWSE TAIEX where robustly available).
- [x] Real ETF proxy fallback: EWY/EWT/EWJ/EWH/MCHI/FXI/VGK/FEZ etc.
- [x] Preserve direct/official-close/proxy evidence separately; never label proxy as underlying.
- [x] AUTO / Direct Only / Free First / Proxy Only routing.
- [x] Real cache / UNAVAILABLE final states; no production fabrication.
- [x] Logical-family correlation avoids treating direct + close + proxy as independent markets.

## Macro / High-Impact Event Engine
- [x] Fed/FOMC, CPI/PCE, employment, GDP, ISM, retail sales, claims.
- [x] ECB/Euro-area, BOJ/Japan and China NBS/PBOC/Customs public event sources where robustly parseable.
- [x] Lifecycle: UPCOMING → RELEASED → MARKET REACTION → RESOLVED.
- [x] HIGH / MEDIUM / LOW impact.
- [x] Date-only official events remain date-only; no invented release time.
- [x] T−15 affected-symbol/sector map and market-critical prewarm.
- [x] FRED/Treasury/BLS/BEA/EIA/direct-global refresh during preparation.
- [x] Nonessential fundamentals/filing background work deferred while event mode is active.
- [x] Decision Queue context is pre-prepared from shared canonical state.
- [x] Reaction snapshots at +5s/+30s/+1m/+5m/+15m/+1h when runtime cadence/data permit.
- [x] Reactions use frozen pre-event price baseline and record app-internal capture latency.
- [x] AI stays asynchronous/off the critical path.

## Macro / official data sources
- [x] FRED: rates, real yields and credit history when configured.
- [x] U.S. Treasury official nominal/real yield XML feed, managed internally.
- [x] BLS official CPI/core CPI/payroll/unemployment actuals, public mode; optional registration key stored safely.
- [x] BEA public GDP/PCE releases.
- [x] EIA WTI/Brent actuals when free EIA key is configured.
- [x] Fed/ECB/BOJ/China official/public calendars/decisions where robustly parseable.
- [x] Fast Macro Consensus + Actual provider contract retained for future premium source; free stack never claims first-second consensus latency.

## Options Intelligence
- [x] Provider abstraction with Alpaca implementation; future provider can plug in without UI redesign.
- [x] OPRA when entitled; AUTO may use real indicative fallback.
- [x] Calls/puts volume, put/call volume ratio, IV, ΔIV after consecutive real snapshots.
- [x] Nearest-expiry near-the-money expected move only when sufficient real IV exists.
- [x] Open interest/unusual-flow unavailable unless a future source supplies it truthfully.
- [x] Dashboard + Day + Swing + selective Long + Queue + Research + AI + Readiness + Signal Validation correlation.
- [x] No Options Trading Desk / no execution / no options P&L.

## Discovery / Research / Decision Queue workflow
- [x] Discovery Day/Swing/Long candidate generation and Discovery Rank separate from trading Score.
- [x] Scanned → Ranked → Qualified → Staged workflow.
- [x] Research tabs: Overview / Fundamentals / Catalysts / Filings; no duplicate Technicals tab.
- [x] Context-aware navigation preserves ticker, origin, horizon and relevant evidence.
- [x] Queue consumes shared state without duplicate subscriptions/API pulls.

## Earnings / SEC
- [x] Expectations, actuals, surprise, Beat/In-line/Miss, guidance state, trend and reaction where sourced.
- [x] Conservative numeric old-vs-new revenue/EPS guidance extraction only from explicit guidance/outlook evidence.
- [x] Prior structured range retained for comparison; missing evidence remains NOT PROVIDED.
- [x] SEC friendly names globally, raw form secondary.
- [x] Day fresh material filing only; Swing compact risk; Long/Research full intelligence.

## AI Copilot
- [x] OpenRouter / Groq / Gemini remain optional providers.
- [x] Research Confirmation / Risk Review / Catalyst Review / Custom Question use compact review-specific evidence.
- [x] Arbitrary ticker; no hard-coded symbol.
- [x] Loading/selected/error states and one smaller-context retry on token/context-limit error.
- [x] Cache key includes ticker + horizon + review type + evidence fingerprint.
- [x] AI cannot trade, alter watchlists, or modify deterministic Score/Action formulas.

## Signal Validation
- [x] De-duplicated ticker/horizon decision snapshots.
- [x] Score/Action, Readiness, Market Regime, Global/Macro, Options, Event Risk.
- [x] Research state, Queue priority, Key Driver and contradictions.
- [x] Later real 1D/3D/5D/10D outcomes + MFE/MAE in Live mode only.
- [x] No automatic formula reweighting.

## Settings
- [x] Finnhub / Alpaca existing providers and AI providers.
- [x] Global mode / Options mode / Macro Event Mode.
- [x] FRED / BLS optional registration / EIA / Twelve Data provider cards.
- [x] Provider Test verifies endpoint/auth/capability rather than non-empty field only.
- [x] Official no-credential sources managed internally rather than fake credential fields.
- [x] Secrets masked after save, excluded from logs/AI/QA/profile export and preserved by migration.
- [x] Standardized dropdown/control styling and keyboard states.

## Maintenance
- [x] Build Identity Verification against backend + packaged renderer version/build.
- [x] System/Data Engine detailed diagnostics.
- [x] Data Capabilities: U.S. equities, Global Indices, U.S. Futures, Macro Events/Actuals, Fast Macro, Options, Signal Validation, Rates/Credit, direct provider.
- [x] Current source/mode/status/freshness/provenance and premium-upgrade readiness.
- [x] Cache/storage controls, Data Freshness, Pre-Market Prep, Weekly Integrity.
- [x] Signal Validation diagnostic view.
- [x] Latest five QA reports.

## Documentation / release
- [x] User Documentation.
- [x] Developer Documentation.
- [x] Capabilities & Limitations.
- [x] Data Capability / fallback / premium-upgrade documentation.
- [x] Permanent Functionality Inventory (this file).
- [x] Requirement-level Scope & Traceability matrix.
- [x] Twelve blocking gates, including Cross-Module Integration/Correlation/Reuse and the separate Tab-by-Tab UI Placement & Information Hierarchy audit.
- [x] Exact Source ZIP re-extraction/retest before final build.
- [x] macOS ARM64 + Intel and Windows x64 package verification, asset/version/hash integrity.

## v14.0.3 UI consistency patch
- [x] Global / Macro Driver cards use collision-safe responsive card layout with separate title/status/detail rows.
- [x] Market Instruments hover/focus keeps identical geometry and continuous ticker motion; no hover pause/jiggle.
- [x] Event & Data Context uses stable header actions plus responsive key/value rows; very old cache is visually warned, never presented as current.
- [x] SEC Intelligence Summary uses consistent key/value typography so values do not overpower labels.
- [x] Focused responsive regressions measure card overlap, ticker hover geometry/play state, Event/Data row overlap, and SEC font hierarchy.

## v14.2 notification-layout hotfix
- [x] v14.0.3 Stable remains the functional/trading/data baseline.
- [x] Transient below-header notifications are centered while visible.
- [x] Notification host preserves its below-header class through show/hide cycles.
- [x] Hidden/dismissed notifications consume zero layout space above Market Instruments.
- [x] v14.1 UI redesign is not part of v14.2.

## v14.2 conservative UI follow-up
- [x] Sidebar Data Engine Details row removed; essential health/runtime remains directly visible.
- [x] Deeper Data Engine diagnostics remain in Maintenance.
- [x] Summary, candidate, news/feed, event and generic empty-state panels no longer keep unnecessary forced minimum height.
- [x] v14.0.3 page hierarchy/interaction is preserved; no v14.1 redesign is included.

## v14.3.0 Improvement Build
- [x] True app-shell toast is independent of page layout and reuses one viewport host.
- [x] Save/action context preserves scroll/focus on long pages.
- [x] Market Open Prep is distinct from Pre-Market Prep and has automatic/manual paths.
- [x] Side Data Engine directly shows Pre-Market Prep, Market Open Prep and event-driven catalyst watcher structured state.
- [x] Earnings & Material Catalyst Reaction Watch activates selectively for scheduled applicable earnings or material unexpected News/SEC events; it is not a daily always-on watcher.
- [x] Liquidity Health uses current bid/ask/spread/size/freshness evidence without fabricating missing depth.
- [x] Provider Capability Registry is entitlement-aware.
- [x] FRED/BLS/EIA/Twelve/Finnhub/Alpaca P1/P2 context is cadence-aware and canonically reused.
- [x] New context is available to relevant Readiness/Queue/Research/AI/horizon surfaces without changing frozen deterministic Setup Scores.
- [x] Dashboard More Market Context uses coherent grouping and corrected spacing rather than appending more provider panels.
- [x] 13 blocking release gates are used, including separate Gate 3 correlation/reuse and Gate 4 tab placement/hierarchy audits.
- [x] Historical Edge Testing remains excluded/deferred.



## v14.3.1 Completion Patch
- [x] Sidebar fixed scroll-fade pseudo-element removed; no black horizontal Data Engine overlay/jitter on scroll or Settings rerender.
- [x] Twelve Data statistics/earnings/institutional data is entitlement-aware fallback/validation only and cannot overwrite populated primary canonical values.
- [x] Market Open checkpoint explicitly reconciles VIX, global/macro, premarket gap/range/volume, options and meaningful horizon changes.
- [x] Day OUTSIDE ENTRY ZONE uses the actual deterministic trade-plan entry range after a current Market Open checkpoint.
- [x] Earnings & Material Catalyst Watch distinguishes ARMED, TRIGGERED and REACTION.
- [x] Provider Capability Registry requires verified provider health before AVAILABLE.
- [x] Provider adapter fixtures cover Finnhub intelligence, Alpaca calendar/activity/corporate actions, Twelve FX/fallback and entitlement truthfulness.
- [x] Slow official macro cadence is startup/prep/event-driven with conservative fallback refresh rather than frequent polling.
- [x] Gate 3 cross-module correlation/reuse and Gate 4 tab-by-tab placement were re-audited for the completion scope.
- [x] Deterministic Day/Swing/Long Setup Score formulas remain frozen.


## v14.3.2 Completion Final
- Global header notification/alert slot with stable header height, no overlap/reflow, truncation for long messages, and accessibility live-region behavior.
- Detailed Data Engine manual actions: Pre-Market Prep Run Now, Market Open Prep Run Now, Catalyst Evaluate, Refresh Due Data, Global/FX Refresh, Provider Capability Recheck, conditional VIX Refresh.
- Manual actions reuse production pipelines and preserve canonical/provider/materiality/deterministic-score protections.


## v14.3.3 Stable

- Dynamic Multi-Feed Subscription Manager with reserved capacity and prioritized promotion/demotion.
- GLD / SLV / USO pinned tradable live coverage.
- Canonical Live Coverage State and Live Market Coverage diagnostics.
- Persistent manual-action lifecycle status with start/completion notifications.
- Side Data Engine generic Manual Actions removed; readiness controls retained.
- Maintenance Data Engine Detail contextual action buttons.
- Header notification lane dynamic fitting and stable header height.
- Adaptive Data Engine geometry: frequently changing values reserve a three-line minimum, simple rows remain natural height, and all real status text can expand without clipping or ellipsis; status typography remains secondary.
- Content/copy consistency and canonical build-identity release checks.

## v14.3.4 Stable

- [x] Side readiness-state typography is secondary to the preparation-job title.
- [x] Maintenance Data Engine action/status/timestamp controls use one consistent centered action cluster.
- [x] QA & Release History renders only the latest five releases with complete status/date/summary/report metadata.
- [x] Decision Queue Why? evidence opens by default.
- [x] Freshness & Cache Policy cards use content-driven heights rather than forced tall cards.
- [x] Header long-message fitting preserves the beginning and uses end ellipsis only after minimum readable font size.
- [x] Settings Auto-Start Runtime is a compact right-aligned header control and preserves saved behavior.
- [x] Maintenance VIX exposes Refresh/Retry through the canonical VIX pipeline.
- [x] Deterministic Day/Swing/Long Setup Scores remain unchanged.


## v14.3.5 Stable

- Maintenance action control visuals refined into integrated DE.PULSE mini-panels.
- Side readiness dots smaller and inline with title.
- Settings Auto-Start premium status capsule retained to the right of Settings but positioned closer to the heading.
- No trading/data-engine functional changes.
