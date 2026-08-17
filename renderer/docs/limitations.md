# DE.PULSE — Capabilities & Limitations

## v18.4.0 STABLE — commercial/data-rights boundary

v18.4 does **not** assert that any provider is commercially approved, redistributable or approved for AI/LLM use merely because an endpoint is public, configured or operational. Current default provider rights are `UNREVIEWED` / `NOT_ASSERTED`, and commercial readiness is fail-closed until provider-specific evidence is bound. Legal/licensing interpretation remains an explicit review responsibility rather than an inferred runtime fact.

Hosted request quotas are abuse/capacity safeguards, not provider-entitlement limits and not a guarantee against upstream rate limits. They do not change desktop behavior or provider routing. v18.4 remains research/intelligence/decision support only; the **No Execution** boundary is unchanged.

## v18.3.0 STABLE — hosted persistence boundaries

v18.3 STABLE has passed PostgreSQL 17 repository parity, migration/export, backup/restore, bounded contention/recovery, hosted liveness/readiness, and required macOS Apple Silicon / Windows x64 release certification. This does not by itself represent the broader commercial/data-rights/security hardening reserved for v18.4.

PostgreSQL is not treated as a cure-all for `DATA DEGRADED`: provider pressure, freshness, queues and DB pressure remain separate capability-level concerns. Hosted mode must fail readiness truthfully when its canonical database is unavailable and must not silently fall back to per-machine local state. Personal workspace isolation does not imply separate market engines; canonical market intelligence is intentionally shared and deduplicated. Broader commercial/data-rights/security hardening remains v18.4. The **No Execution** boundary is unchanged.

## v18.2.0 STABLE — administration / presence boundaries

The promoted Stable build uses the canonical `PersonalMarketTerminal` profile; the separate v18.2 TEST profile remains historical and is not the Stable runtime target.

v18.2 adds local/shared-process administrative user and session operations but is **not yet the true hosted multi-user deployment**. Presence reflects authenticated session activity as ACTIVE / IDLE / OFFLINE; it is operational presence, not device identity, geolocation or guaranteed human attention. Session history is bounded operational context.

Administrators can manage only roles below their authority; own-account destructive changes and elimination of the final active OWNER/SUPER_OWNER are blocked. Temporary-password creation/reset requires the user to replace the password on sign-in and revokes prior sessions. Hosted PostgreSQL, browser-scale concurrency, backup/restore and shared-server recovery remain v18.3. Broader commercial/data-rights/security hardening remains v18.4. The **No Execution** boundary is unchanged.

## v18.1.0 STABLE — multi-user boundaries

v18.1 isolates personal market state, but it is not yet the full hosted multi-user product. User-administration, presence/session operations and lifecycle controls remain v18.2; true shared hosted PostgreSQL/browser deployment remains v18.3. Provider entitlements, canonical market evidence and global operational policy are intentionally shared rather than duplicated per account. Broad-market Opportunity Radar and Rapid Move / Market Shock intelligence may therefore be visible to multiple users even when it was not triggered by their personal watchlist.

The Stable build uses the canonical `PersonalMarketTerminal` profile and preserves compatible v18.0.6 state in place; TEST-profile clone/isolation behavior is not used by the Stable channel. DE.PULSE remains decision support, not a profit predictor, and the **No Execution** boundary remains unchanged.

## v18.0.6 STABLE — current capability boundary

Rapid Move / Market Shock coverage remains truthful **TIERED_PARTIAL**, not full-market surveillance. Adaptive outcome learning remains SHADOW-only and cannot auto-promote policy changes. DE.PULSE remains research/intelligence/decision-support only; no execution, paper trading, portfolio/P&L or order-management scope is introduced.

## v18.0.5 STABLE — Current truth boundaries

v18.0.5 improves symbol-management correctness and presentation; it does not expand market-data entitlement or guarantee opportunity detection. Remove All clears user-managed Day/Swing/Long tracked membership but does not remove protected market-context processing. USER/DEMO diagnostic suppression is presentation/privacy hardening, not evidence deletion. Opportunity Radar remains bounded investigation priority and Research remains decision support, not a profit predictor. v18.1 multi-user architecture is not included. Protected deterministic formulas and the **No Execution** boundary are unchanged.

## v18.0.4 STABLE — Delivery status

v18.0.4 is Stable after native Apple Silicon and Windows x64 G14 audits and G15 Release Assurance passed. It does not expand market coverage, provider entitlements, adaptive authority, or the No Execution boundary. Stable runtime uses the canonical `PersonalMarketTerminal` profile and preserves compatible prior Stable state.

## v18.0.3 TEST — Delivery status

v18.0.3 remains TEST until native Apple Silicon and Windows x64 G14 audits pass and G15 returns READY_FOR_PROMOTION. This patch hardens native runtime/resource portability and test/delivery correctness; it does not expand market coverage, provider entitlements, adaptive-authority state, or the No Execution boundary.

## v18.0.2 TEST — Delivery status

v18.0.2 remains TEST until physical/native Apple Silicon and Windows x64 G14 audits pass and G15 returns READY_FOR_PROMOTION. This patch corrects release fingerprint portability only; it does not expand provider entitlements, market coverage, or the No Execution boundary.


## v18.0.1 TEST — Smart Router / Rapid Move truth boundaries

Rapid Move detection is **tiered partial coverage**, not a claim that every U.S. stock is streamed at 15s/30s/60s resolution. Exact short-window detection applies to symbols receiving current canonical observations; broad-market seeds may promote additional symbols only when provider capabilities/entitlements support them. Provider outages, entitlement limits, halts, delayed headlines and unavailable corroborating sources can still delay or limit conclusions.

A +5%/60s move is a strong baseline condition, not an automatic trading signal. DE.PULSE validates materiality using context and suppresses known mechanical/noisy cases such as current split/corporate-action discontinuities, severe source conflict and low-priced/liquidity-risk bursts. `UNEXPLAINED` means price evidence is material before a catalyst is confirmed; it does not imply hidden information or predict continuation. Adaptive learning remains non-authoritative until separately validated/approved.

Smart Router scoring does not manufacture provider entitlements: capability truth is learned from configured/observed provider behavior and revalidated after bounded cooldowns. The permanent **No Execution Boundary** remains unchanged.

## v18.0.0 TEST — Current identity/session boundaries

v18.0 is a **TEST-only security foundation**, not a multi-user commercial release. It currently bootstraps one local OWNER account and provides password-backed login/logout, opaque server-side sessions, role enforcement foundations and isolated TEST-profile migration. Per-user symbol ownership, user-administration UI, presence/session operations, configurable inactivity/weekend policies and PostgreSQL hosted shared state remain later v18 slices.

The local TEST server uses secure cookie attributes appropriate to its transport: `HttpOnly` + `SameSite=Strict`, with `Secure` enabled when TLS is present. State-changing authenticated requests require a CSRF token. Full brute-force/rate-limit policy, broader adversarial security qualification and commercial/data-rights hardening remain v18.4 scope.

The permanent **No Execution Boundary** remains unchanged. The v18 TEST profile is intentionally separate from v17.5.1 Stable. First-run migration requires Stable to be closed briefly so its persisted profile can be copied consistently; Stable may be reopened afterward. The TEST build must not be treated as the authoritative Stable release until its own qualification and later promotion process are complete.


## v17.5.1 STABLE — Current limitations / truth boundaries

DE.PULSE remains **decision support, not a profit predictor**. Completion of v17 does not guarantee profitable decisions or eliminate provider/network outages. Provider entitlements, publication delays, exchange/source corrections and live-market conditions can still limit evidence. Local SQLite persistence, workload budgets and freshness/readiness safeguards improve resilience but do not create missing market data. Actionable instruments remain **U.S.-listed only**. v18 authentication/multi-user/session/quota capabilities are intentionally absent. There is **no execution**, paper trading, Portfolio/P&L, Journal, OMS/blotter or broker routing. Native macOS SQLite acceptance passed; Windows x64 SQLite packaging/compile contract is certified, while physical Windows-host acceptance remains a documented non-blocking residual.

## v17 release evolution — all checkpoints incorporated into Stable

### v17.5.0 — Major Closure
The full v17 scope was freshly reconstructed/certified and native macOS SQLite runtime acceptance completed. The historical RC wording shown before v17.5.1 was a release-documentation defect, not an unfinished v17 architecture state.

### v17.4.0 — UX / operational hardening
Grouped preparation exceptions reduce repetitive presentation without hiding genuine risks. Freshness failures remain conservative and can block readiness; genuine current liquidity/event/extended risks remain drillable per symbol.

### v17.3.0 — Performance / reliability hardening
Runtime SLOs are defensive acceptance budgets, not provider entitlement guarantees. CPU/storage/write thresholds and queue limits protect local responsiveness; external provider availability still constrains evidence.

### v17.2.0 — Canonical intelligence / data efficiency
Material-change propagation reduces unnecessary downstream work, but it intentionally does not suppress freshness threshold crossings, selected/promoted-symbol work or material events. Historical evidence quality still bounds reproducibility.

### v17.1.0 — Runtime load / backpressure
Tier 0/1 work is protected before broad/background work, but sustained provider outages/rate limits can still degrade critical evidence and must be reported truthfully. Lower-priority work may be delayed or shed under pressure by design.

### v17.0.0 — Persistence foundation
Persisted state accelerates warm start but is never automatically LIVE/CURRENT: timestamps/session/provider freshness policy are re-evaluated on read. SQLite retention is bounded; the database is not an unbounded raw-data dump.

## v16.11.0 Stable — Major Closure truth boundaries

Major Closure materially increases assurance but cannot guarantee profitable decisions, eliminate provider outages, or substitute for target-host/provider-entitlement testing. DE.PULSE remains **decision support, not a profit predictor**. Physical native macOS/Windows launch/sleep-wake and authenticated provider entitlement acceptance remain target-host responsibilities when unavailable on the certification host.

Freshness truth is now coupled to readiness safety: stale/cached/history-only current-price evidence cannot show a clean Day `READY` state. This lowers confidence/readiness only and does not rewrite deterministic scores. Community/social evidence remains untrusted and source-policy controlled. Opportunity Radar remains bounded U.S.-equity investigation priority, not an execution or guaranteed-opportunity engine. Shadow remains read-only until separately validated/approved.

The original professional roadmap remains **30 FULL / 0 PARTIAL / 0 MISSING**, v16.10 remains **10/10**, and actionable instruments remain **U.S.-listed only**. A future major-family transition requires its own fresh Major Closure build rather than relying only on accumulated historical audit reports.

## v16.10.0 Stable — Opportunity Radar limitations

Opportunity Radar is **decision support, not a profit predictor**. It can miss opportunities or surface false positives because free/provider entitlements, snapshot cadence, delayed data, liquidity changes and event timing constrain evidence. Broad scanning is deliberately bounded rather than a permanent live stream of the full U.S. market. A Radar promotion means “investigate now,” not BUY/SELL and not guaranteed volatility continuation.

Session-normalized RVOL and range expansion depend on available current/prior bars and truthful provider timestamps. Illiquid, low-price, wide-spread and low-dollar-volume spikes are gated out where evidence supports that judgment. Adaptive freshness/cache changes remain bounded; Shadow experiments are read-only and cannot mutate production. Physical native macOS/Windows launch/sleep-wake and authenticated provider-entitlement acceptance remain target-host actions. Original professional roadmap remains **30 FULL / 0 PARTIAL / 0 MISSING** and all actionable symbols remain **U.S.-listed only**.

## v16.9.0 Stable — Current Limitations / Truth Boundaries

The original professional roadmap is **30 FULL / 0 PARTIAL / 0 MISSING**, but completion does not remove external data-rights, entitlement or evidence limitations. Community sources are permitted/user-authorized only; scraping/bypass is rejected. Social/community evidence is untrusted and source-policy controlled. Telegram, Reddit, X, Discord and WhatsApp content may be context-only for AI/retention unless explicit rights permit more. Commercial/hosted use still requires a source-by-source rights audit.

WTI/Brent context uses the best canonical evidence actually available. `USO` is not WTI spot, and DE.PULSE does not claim CL/BZ continuous-contract or roll-adjusted futures semantics without a source that truly provides them. Historical scenario replay requires retained cutoff-safe evidence; missing scenarios remain `UNAVAILABLE` and no lookahead, mock-live mode or paper trading is allowed. Actionable symbols remain **U.S.-listed only**; international evidence remains selective **Global Context**. Physical native macOS/Windows launch/UAT and real credential/entitlement behavior remain target-host responsibilities.

## v16.8.1 Stable — Current Limitations / Truth Boundaries

DE.PULSE actionable symbol processing is intentionally **U.S.-listed only**. International equities are not supported as watchlist/desk/Research/Discovery instruments. Global macro/index/FX/commodity evidence remains selective context for the U.S. market and does not automatically trigger full U.S. Event Mode. U.S.-listing verification depends on the available recovery/provider metadata in live mode. Global Context may omit routine foreign releases by design when their U.S.-market utility is low. Physical native macOS/Windows launch/UAT remains host-dependent. Original-roadmap status remains **27 FULL / 3 PARTIAL / 0 MISSING** until v16.9.


## v16.8.0 Stable — Current Limitations / Truth Boundaries

Heat broad-market mode is coverage-truthful and does **not** claim full S&P 500 coverage when canonical constituents are unavailable. GEX is a structural gamma×OI context proxy and never claims measured dealer positioning. Seasonality is descriptive, not predictive. Liquidity/slippage cannot claim order-book depth when depth evidence is unavailable. Smart Notifications are material state changes only; no standalone Alerts workspace exists. Physical native macOS/Windows launch/UAT remains host-dependent. Remaining original-roadmap PARTIAL items are #10, #11 and #20.


> **v16.7.0 Stable acceptance limitations:** the five closure items are complete at their original approved depth using the currently available canonical evidence, but tracked-universe Breadth is still explicitly **not exchange-wide breadth**. Provider entitlements determine whether all source fields are available at runtime; unavailable evidence must degrade honestly. Physical native macOS/Windows launch/sleep-wake and authenticated live-provider entitlement acceptance still require the target machines/credentials and are not claimed by the Linux certification host.

> The remaining original-roadmap PARTIAL items after v16.7 are **#6 Heat Map, #8 Seasonality, #9 GEX, #10 Community, #11 Oil/Energy, #20 Replay, #21 Smart Notifications and #27 Liquidity/Slippage**. Their current placement in v16.8/v16.9 is adaptive and may MOVE/MERGE/SPLIT after the v16.7 Roadmap Adaptation Review, but original acceptance clauses may not silently disappear.

## Core Capabilities & Limitations Reference

> **v16.6.0 Stable acceptance limitations:** the release reconciles all 30 professional capabilities and closes the Master Symbol Store empty-state defect, but physical native macOS/Windows launch and authenticated provider-entitlement checks still require the target machines/credentials. Linux-host certification does not claim those checks were physically executed.

> **Master Symbol truth:** clearing **User Master Symbol Store** removes user-managed DAY/SWING/LONG membership and can legitimately reach zero. Protected **System Market Context** (for example SPY/QQQ and required market instruments) and the VIX special-index path remain available for market context; they are separate from the user-removable count.

> **v16.5.0 Stable context limitations:** Sentiment is a transparent composite of available canonical components, not a prediction. Missing inputs lower confidence. Heat-map coverage is limited to the explicit tracked sector benchmark universe and only current members vote. GEX is withheld when gamma/OI coverage or recency is insufficient and must not be interpreted as measured dealer positioning.

> **Community limitations:** Community Intelligence is **manual/user-authorized only**, always **UNTRUSTED COMMUNITY INTELLIGENCE**, may be incomplete/biased, and is never a verified source of truth or deterministic trading input. No automatic social/community scraping is performed.

> **Oil/Energy limitations:** Official EIA data and tradable market instruments are kept distinct. `USO` is a futures-based ETF, not WTI spot; basis, futures-curve and roll effects can diverge. Provider publication schedules, revisions and entitlements remain external constraints. Standard `De-Pulse.app` identity and compatible prior Stable API keys/settings are preserved.

> **v16.4.0 Stable AI limitations:** AI is an evidence-grounded second opinion, not a deterministic decision owner. External/news/SEC/provider/web text is untrusted. Missing evidence remains missing; an unavailable provider route remains unavailable. AI confidence is not win probability, and AI cannot change Score/Action, Market Tradeability, Trade Readiness, membership or execution. Material-evidence caching improves cost/latency but does not make stale evidence fresh.

> **Routing limitations:** Efficient/Balanced/Deep use only providers for which you configured credentials; provider/model availability, entitlement, quota and latency can still degrade. Manual mode never silently switches providers.

> **Stable identity/continuity:** `De-Pulse.app` is the standard macOS application identity; compatible prior Stable secrets/settings use `PersonalMarketTerminal`.

> **v16.3.0 Stable — Capabilities & limitations:** Seasonality is descriptive historical context only; current scope prioritizes SPY/QQQ and requires explicit sample counts/date ranges. Calibration groups frozen deterministic outcomes but does **not** convert Setup Score into win probability and cannot auto-learn/rewrite production weights. Correlation/Concentration uses aligned daily returns and candidate/watchlist context; it has no Portfolio, positions/P&L or brokerage knowledge and never automatically rejects a trade.

> **Outcome/replay limitations:** OHLC bars cannot always establish intrabar ordering. If Entry/Target/Invalidation share one unresolved bar, DE.PULSE reports `AMBIGUOUS` rather than guessing. Missing post-snapshot history becomes `PENDING`, `PARTIAL` or `UNAVAILABLE`; it is never zero/success. Exact replay requires the v16.3 frozen feature/evidence/settings package. Older snapshots may only support a `SCENARIO` replay, and Long-Term legacy replay is withheld when historically knowable fundamentals cannot be proven. Corporate-action-adjusted history is handled through canonical corporate-action truth; replay is refused when adjustment look-ahead cannot be removed safely.

> **External constraints:** Historical depth, provider corrections, exchange calendars and entitlements still bound what can be proven. Physical native macOS/Windows acceptance remains a target-host release step. These limitations do not relax the permanent no-lookahead, no-false-Fresh/Ready/Live, or deterministic Score/Action protections inherited from v16.0–v16.2.

> **v16.2.0 TEST — Capabilities & limitations:** Event Intelligence is only as complete as the canonical source evidence. News materiality is a deterministic keyword/context classification, not AI certainty. Headline clustering groups defensible near-identical normalized stories but cannot guarantee every syndicated variation is merged. Economic Surprise is unavailable unless both Actual and Forecast are sourced. Fed phases/timeline are shown only when the official/canonical calendar provides distinct evidence. Reaction Intelligence reports observed moves around captured event times; it does not prove causation. Smart Notifications are in-app attention items and do not provide OS push delivery. Event risk can degrade Trade Readiness/Queue priority but never changes deterministic Score/Action.

> **Data limitations:** Live provider credentials/entitlements determine available News, calendar enrichment and reaction data. Missing/stale/future-skew evidence remains `UNAVAILABLE`, `DATA DEGRADED`, or omitted as appropriate; missing values are never coerced to zero/neutral. Physical native macOS/Windows acceptance remains a target-host release step.


> **v16.0.6:** Final v16.0 consolidation/delivery closure. Carries all v16.0.5 professional-truth fixes, keeps the Stable-style primary header, restores the normal native app launch path with binaries prepackaged, and makes the complete release/handoff artifact set mandatory. No v16.1 Market Intelligence/Tradeability feature is included. Deterministic Day/Swing/Long Score/Action formulas remain frozen; v15.1.2 remains the protected Stable baseline.

> **v15.1.2 Stable:** Freshness now distinguishes source observation age from successful reconciliation age and uses dataset/session/event-aware cadence. Provider entitlements, exchange/source delays, native target-OS launch acceptance, and authenticated live-provider acceptance remain external constraints. Packages remain unsigned/not notarized unless separately signed.


> **v15.0.1 Stable:** yfinance is recovery-only and CBOE is VIX-specific; provider freshness/source labels must remain truthful. Native macOS/Windows launch and authenticated entitlement acceptance still require the installed target environment. Packages remain unsigned/not notarized unless separately signed.

> **v14.3.7 Stable:** Font-fit visual patch. Existing intentional truncation/scroll behavior is preserved where content is designed to be secondary or horizontally scrollable; no information model or trading behavior changed.


> **v14.3.6 Stable:** Treasury official yields now keep core Macro Rates usable when optional FRED enrichment is unavailable; FRED-specific credit/conditions/USD capability is still reported separately and truthfully. The patch also includes SEC transaction clarity and the approved Discovery/regime/watchlist UI corrections. Native target-OS launch and authenticated provider entitlement acceptance still require the installed environment.

> **v14.3.4 Stable:** Native target-OS launch and authenticated live-provider entitlement acceptance still require the installed environment. Header notifications use end ellipsis only after fitting to the minimum readable size.


> **v14.3.3 Stable:** Live coverage is entitlement-aware and prioritizes trading-relevant symbols; Alpaca Basic IEX is not consolidated SIP. Passive symbols may use snapshots by design. Header notification text may ellipsize only after dynamic font fitting reaches the minimum readable size. Native target-OS launch acceptance remains required.

> **v14.3.2 Completion Final:** Header notifications are constrained to the existing header lane and truncate when space is limited; Data Engine manual actions require an active runtime and still obey provider/entitlement/materiality rules. Native desktop-shell acceptance remains required on target OS.

> **v14.3.1 Completion Patch:** Provider availability now requires verified health; scheduled earnings remain ARMED until actual release evidence appears; Twelve Data fallback is entitlement-aware; Market Open context is broader but still contextual. Native desktop-shell acceptance remains required on target OS.

> **v14.3.0 Improvement Build:** Built from exact v14.2.2 Stable. Provider capability/entitlement truthfulness, event-driven catalyst reactions, Market Open Prep, and new derived contextual states are added without changing deterministic Setup Scores or fabricating unavailable data.

> **v14.2 Stable Hotfix:** v14.0.3 Stable remains the functional/architecture baseline; v14.2 changes only transient notification layout/cleanup, forced empty-space sizing, the sidebar Data Engine Details row, and release identity.


This tab is operationally honest by design. A fixable bug is not a limitation: it must be fixed before Stable. A limitation is an external, entitlement, platform, latency, coverage, or source-quality constraint that DE.PULSE cannot truthfully remove.

## Production data invariant
Outside Demo Mode, DE.PULSE displays only sourced real data. The permitted quality path is:

`LIVE/DIRECT → OFFICIAL/CURRENT → DELAYED → REAL PROXY → REAL LAST-KNOWN-GOOD CACHE → UNAVAILABLE`

The application never generates plausible-looking production market values when a provider fails.

## Data Capability Matrix

| Capability | Current source / mode | Fallback | Current limitation | Trading impact | Upgrade path / status |
|---|---|---|---|---|---|
| U.S. equities | Alpaca IEX/snapshots primary | Finnhub secondary; Twelve Data tertiary recovery | Provider entitlement/rate limits can degrade coverage | Live Day/Swing context can become delayed/degraded | Implemented |
| Overnight U.S. | Alpaca Auto/Indicative/Live | Last-known-good real cache | Overnight entitlement/coverage differs from regular session | Treat indicative values as context, not exchange-grade regular trades | Implemented |
| VIX | Dedicated true-index/provider path | Real cache / unavailable | Not treated as a normal equity symbol | Volatility context may degrade independently | Implemented |
| Global markets | Optional Twelve Data direct adapter + official/public completed-session data where robust + real U.S. ETF proxies | Official/public → real proxy → cache | TWSE TAIEX official close is implemented; other underlyings may require direct/licensed feed and otherwise remain clearly labeled proxy/unavailable | Useful for contextual read-through; direct/official/proxy provenance is preserved separately | Implemented · Free/Proxy mode · Direct/Premium upgrade available |
| U.S. futures | Implemented direct ES/NQ/RTY provider contract + optional Twelve Data adapter | SPY/QQQ/IWM/overnight real proxy/cache | True licensed/streaming ES/NQ/RTY depends on provider entitlement | Overnight macro reaction is contextual when direct futures are unavailable | Implemented interface · Premium/direct entitlement optional |
| Rates/credit | U.S. Treasury core curve + optional FRED enrichment | Cached official rates / unavailable | Treasury provides core yields; FRED expands credit/conditions/USD and is not tick-level | Suitable for macro regime, not institutional microsecond reaction | Implemented |
| Macro calendars | Fed/BLS/BEA/ECB/BOJ/China official/public pages | Cached real events | Page formats can change; some events are date-only | Unknown release times are never invented; Event Mode requires known time | Implemented |
| Fast macro consensus + actual | `FastMacroProvider` contract; official/free actuals remain current operating mode | Official release / unavailable | No premium low-latency provider is configured by default; free/official stack does not guarantee first-seconds consensus + actual | Do not use as institutional first-second CPI/NFP feed | Provider interface ready · Premium upgrade available |
| Options Intelligence | Alpaca options snapshots; OPRA when entitled, indicative fallback in AUTO | Real indicative data / unavailable | Current snapshot path may not expose open interest; entitlement controls OPRA | Options remains context; unsupported OI is never guessed | Implemented · Premium quality depends on entitlement |
| Options expected move | Nearest-expiry near-the-money sourced IV | Unavailable | Requires underlying price + suitable IV contracts | No expected move is shown when inputs are incomplete | Implemented |
| Earnings guidance | Structured sourced values when present; qualitative evidence from mapped company/SEC context | Not Provided | Numeric old/new guidance cannot be inferred safely from absent structured evidence | Missing guidance remains missing | Implemented evidence gate · richer transcript/provider upgrade optional |
| SEC Intelligence | SEC/company filing services | Real cache | Filing delivery/metadata timing can vary | Day only surfaces fresh material filing catalysts; Long gets full context | Implemented |
| AI Copilot | OpenRouter / Groq / Gemini | Provider-specific error/fallback behavior | Provider quotas/token limits and model quality vary | Advisory only; never changes Score/Action | Implemented structured evidence + retry |
| Signal Validation | Canonical decision snapshots + later real historical bars | Pending until real future bars exist | No outcome can be known before the future bar exists | Evidence accumulates over time; formulas are never auto-reweighted | Implemented |

## Horizon-specific Options policy
- **Dashboard:** only material watchlist/market options context and contradictions.
- **Day:** immediate IV/volume/event confirmation or conflict.
- **Swing:** persistent positioning, IV trend/expected move and catalyst context.
- **Long-Term:** only material longer-horizon positioning/risk; short-lived flow is intentionally hidden.
- **Discovery:** options is intentionally hidden during broad candidate generation and becomes relevant after a candidate is staged/selected, avoiding high-cost chain polling for the scan universe.

## Cross-module integration matrix

| Capability | Dashboard | Day | Swing | Long | Queue | Discovery | Research | AI | Settings | Maintenance | Docs |
|---|---|---|---|---|---|---|---|---|---|---|---|
| Global/Macro | Implemented | Implemented | Implemented | Implemented | Implemented | Context used after qualification | Implemented | Implemented | Config | Capabilities | Implemented |
| Options | Material pulse | Immediate context | Persistent context | Selective context | Confirmation/conflict | Intentionally hidden until staged | Deep context | Relevant evidence only | Config | Capabilities | Implemented |
| Trade Readiness | Queue/overview | Implemented | Implemented | Implemented | Implemented | Not a scanner rank | Implemented | Evidence only | Not Relevant | Validation | Implemented |
| Signal Validation | Not primary | Snapshot source | Snapshot source | Snapshot source | Context captured | Not Relevant | Diagnostic evidence | Not formula authority | Not Relevant | Primary diagnostics | Implemented |

## Current UI / integration guarantees
- Expandable panels preserve user open/closed state across live rerenders.
- Sidebar Data Engine shows detailed live provider/runtime/data-health rows directly at normal desktop widths; Maintenance retains deeper troubleshooting diagnostics.
- Options Intelligence uses compact DE.PULSE terminal styling rather than large flat gray rows.
- Maintenance verifies backend + renderer build identity; stale v13.x build labels are treated as a mismatch, not silently displayed as current.
- Global, macro, options, earnings and SEC evidence is evaluated for cross-module impact. “Intentionally Hidden” is used only when the metric would be horizon-inappropriate noise.

## Platform/package limitations
Native macOS Gatekeeper/WKWebView/Finder behavior and Windows shell/Explorer presentation require final acceptance on the target OS. Authenticated private-provider behavior requires the user's real credentials/entitlements. Automated tests can validate request structure, fallbacks and error handling without exposing secrets, but cannot manufacture an entitlement.

## Explicit non-capabilities
DE.PULSE does not provide order entry, brokerage execution, positions/P&L, Portfolio, Journal, generic Alerts, Workspace Visibility, autonomous trading, or an Options Trading Desk.


- Options IV change requires at least two consecutive real option snapshots; the first valid snapshot has no ΔIV. Open interest/unusual-flow classification remains unavailable unless a future source provides those fields truthfully.

## v14.0.3 stale-cache presentation

A cached quote can remain available as truthful last-known-good evidence when live providers are unavailable. When desk context shows a cache older than 30 minutes, DE.PULSE now emphasizes that age visually. This is a presentation safeguard, not a fabricated refresh; the application still follows real data → real cache → unavailable semantics.

## v14.3.0 Additional Capability / Limitation Notes

- Provider Capability Registry distinguishes **AVAILABLE**, **PLAN LIMITED**, **NOT ENTITLED**, and **TEMPORARILY UNAVAILABLE**. A configured API key is not treated as proof of entitlement.
- Finnhub recommendation/price-target/insider/institutional datasets are contextual and may be plan-limited. DE.PULSE degrades truthfully rather than depending on them.
- Alpaca Basic/IEX and other entitlements can limit full-market activity or quote quality. Liquidity Health reflects the data actually received and does not imply full-SIP depth.
- Official FRED/BLS/EIA data has release-specific cadence and is unsuitable for pretending to be tick-level data. Derived states preserve source freshness.
- BLS actuals do not imply Wall Street consensus. DE.PULSE does not invent surprise values when a sourced estimate/consensus is unavailable.
- EIA electricity/STEO thematic data is best-effort sector/research context and is not a Day Setup Score input.
- Direct Twelve Data global/FX availability remains entitlement-aware; the fallback/provenance hierarchy is preserved.
- Premarket/after-hours earnings or catalyst reactions are **not regular-session confirmation**. Trade Readiness/Market Open Prep revalidate price extension, event risk, and liquidity.
- Historical profitability/edge testing remains deferred and is not claimed by this release.

## v15.0.1 Provider consolidation policy
- **Alpha Vantage** is demoted to emergency / low-frequency fallback only and remains a removal candidate after dependency validation. New active routing must not depend on it.
- **Massive** and **FMP** are intentionally **not added** in this release. DE.PULSE will not add another general-purpose market-data provider unless a verified capability gap remains after the approved Alpaca / Finnhub / Twelve Data / Marketaux / SEC / FRED core plus yfinance/CBOE recovery stack is evaluated.
- yfinance remains automatic recovery-only; CBOE remains VIX-specific validation/delayed fallback. Neither is presented as a configurable credentialed provider.

## v18.3 STABLE — archive and outage limits

Persistence archives contain credential/session hashes and must be stored privately. Restore is empty-target-only by default; replacing existing canonical state requires the explicit `replace` mode. The archive does not carry provider/API secrets from `secrets.json`.

During a prolonged canonical-database outage, DE.PULSE bounds queued persistence memory and may shed reproducible derived-feature records before immutable evidence/decision/outcome lineage. Readiness remains unavailable until recovery hysteresis succeeds. v18.5 still requires realistic supported-load closure for prolonged outage, queue saturation, restart/warm-start and recovery-flap behavior before the v18 major line is considered fully closed.
