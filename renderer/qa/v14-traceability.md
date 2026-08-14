# DE.PULSE v14.3.7 — Requirement-Level Scope & Traceability

Release traceability artifact. Status is not inferred from the presence of a tab: every requirement must have implementation, engineering/functional/trader/cross-module validation and documentation. The exact v14.3.4 Stable Source ZIP is the baseline for v14.3.6; the final v14.3.6 Source ZIP must be freshly extracted and fully retested before platform builds. Package integrity remains a blocking release condition.

| # | Requirement | Implementation | Engineering / Functional Test | Trader / Cross-Module Test | Documentation | Status |
|---|---|---|---|---|---|---|
| 1 | Preserve v13.3 deterministic Day/Swing/Long Score/Action | Existing `planFor`/desk math unchanged | `deterministic_equivalence_test.js` 2,403 cases | Day/Swing/Long acceptance scenarios | Developer + Capabilities | PASS |
| 2 | Real-data-only outside Demo | provider adapters/fallbacks return real/cache/unavailable only | Go provider/degradation tests + HTTP Demo isolation | Provider-disconnect scenario | User/Developer/Limitations | PASS |
| 3 | Finnhub primary / Alpaca controlled fallback / VIX special | existing canonical provider paths retained | Go + renderer + HTTP invariants | fast-market/provider-fallback scenario | Developer | PASS |
| 4 | Canonical Master Symbol Store / no duplicate subscriptions | existing shared state/subscription manager retained | unit/race/stress + renderer store tests | all tabs consume shared state | Developer | PASS |
| 5 | Global Driver Engine shared/canonical | `deriveGlobalMarketContext`, canonical runtime snapshot | Go driver tests | Dashboard/Day/Swing/Long/Queue/Research/AI correlation | User/Developer/Limitations | PASS |
| 6 | Direct global provider runtime path | `DirectGlobalProvider`, Twelve Data adapter | mocked search/quote provider tests | direct vs fallback provenance | Developer/Limitations | PASS |
| 7 | Direct ES/NQ/RTY interface | `DirectFuturesProvider`, Twelve Data futures adapter | mocked futures provider test | Dashboard/global/event context when available | Developer/Limitations | PASS |
| 8 | Official/public international-close path + session awareness where robustly available | TWSE official TAIEX close adapter; other markets remain direct-provider/proxy/unavailable rather than fabricated | TWSE mocked official response + global routing tests | official close retained alongside current proxy; completed local session stays current context | User/Developer/Limitations | PASS |
| 9 | Direct → official/public → proxy → cache → unavailable | provider-mode routing + persisted `globalDirect` + ETF proxy layer | Go routing/provenance tests | Market Context progressive disclosure | Developer/Limitations | PASS |
| 10 | AUTO / Direct / Free First / Proxy modes | Settings + `deriveGlobalMarketContext` routing | renderer settings + Go routing tests | no horizon formula impact | User/Developer | PASS |
| 11 | Global drivers correlate to overall Market Regime | contextual `marketRegime()` global component, grouped logical families | renderer logic + 2,403 equivalence | Dashboard regime/desk divergence scenarios | User/Developer | PASS |
| 12 | Broad-market options may influence regime; ticker options may not | SPY/QQQ-only `broadOptionsRegimeComponent` | renderer logic | trader contradiction tests | User/Limitations | PASS |
| 13 | True/proxy breadth not watchlist breadth | broad ETF/sector proxy driver, explicit REAL PROXY label | renderer/provider tests | Dashboard context review | User/Limitations | PASS |
| 14 | FRED rates/credit | `refreshFRED` with multi-horizon history | local mocked FRED + Go tests | Swing/Long macro context | Developer/Limitations | PASS |
| 15 | Treasury official nominal/real yields | `refreshTreasury` | official XML mocked adapter test | Market regime + horizons | User/Developer | PASS |
| 16 | BLS actuals | `refreshBLSActuals` | public API mocked adapter test | macro event/Research evidence | Developer/Limitations | PASS |
| 17 | BEA actuals | `refreshBEAActuals` | public release mocked adapter test | Long/global context | Developer/Limitations | PASS |
| 18 | EIA actuals | `refreshEIAActuals` | mocked EIA adapter test | Energy/global context | Developer/Limitations | PASS |
| 19 | Macro event universe US/EU/JP/CN | expanded official/public sources | parser/degraded tests | horizon-specific event visibility | User/Developer | PASS |
| 20 | Event lifecycle | `eventLifecycle` incl MARKET REACTION | Go lifecycle test | Dashboard/Day/Queue event interpretation | User/Developer | PASS |
| 21 | High-Impact T−15 prep | `eventModeLoop`, affected map, history/provider prewarm, deferred background refresh | Go/event-mode + no-repeat prep tests | FOMC/CPI scenario | User/Developer | PASS |
| 22 | Reaction snapshots vs event-time baseline | frozen pre-event baseline + 5/30/60/300/900/3600 sec records | Go event reaction tests | Day/Queue/Global panel | User/Developer | PASS |
| 23 | AI off macro critical path | event processing never calls AI; AI remains user/asynchronous | static/Go assertions | macro event scenario | Developer | PASS |
| 24 | Horizon Lens + Dashboard order | compact regime/Horizon Lens → Queue → Overview → Catalyst → Context | 284 responsive/DOM/cross-device checks | professional information-hierarchy review | User | PASS |
| 25 | Decision Queue enriched correlation | Action/Score/Priority/Readiness/Research/Key Driver/Risk + Why | renderer logic/UI | cross-module queue workflow | User/Developer | PASS |
| 26 | Catalyst freshness + View All | 72h news / 30d filing / ±14d earnings filters | renderer logic/UI | dashboard clutter review | User | PASS |
| 27 | Swing genuine 5D/20D context | FRED/Treasury deltas + UUP/SMH series returns | Go/renderer logic | Swing scenario | User/Developer | PASS |
| 28 | Long genuine 1M/3M regimes | rates/real yield/credit + USD/China trend history | Go/renderer logic | Long investor scenario | User/Developer | PASS |
| 29 | Trade Readiness deep horizon inputs | `tradeReadiness` Day/Swing/Long orchestration | renderer logic | trader scenarios; deterministic scores unchanged | User/Developer | PASS |
| 30 | Options provider abstraction | `OptionsIntelligenceProvider`, Alpaca implementation | Go options tests | all relevant surfaces reviewed | Developer/Limitations | PASS |
| 31 | Options IV/P-C/ΔIV/expected move | real option snapshot aggregation + corrected ATM expected move | Go options tests | Day/Swing/Long policy | User/Limitations | PASS |
| 32 | Options compact styling | `.options-intel-row` terminal styling | responsive visual DOM check | Dashboard clutter review | User | PASS |
| 33 | Signal Validation Research + Queue context | extended `SignalSnapshot` and renderer record payload | HTTP/Go/renderer tests | validation diagnostic review | User/Developer | PASS |
| 34 | Numeric old-vs-new guidance only from evidence | conservative guidance parser + news/8-K evidence enrichment | Go explicit/negative parsing tests | earnings mixed-result scenario | User/Limitations | PASS |
| 35 | AI compact-context hotfix | review-specific evidence builder + one reduced retry + cache fingerprint | Go/renderer AI tests | AI disagreement remains non-binding | User/Developer | PASS |
| 36 | SEC horizon behavior / friendly names | existing SEC Intelligence retained | renderer/HTTP tests | Day/Swing/Long review | User/Developer | PASS |
| 37 | Persistent expanders | shared expansion-state manager across rerenders | explicit Playwright rerender test | all progressive-disclosure surfaces | User/Developer | PASS |
| 38 | Rich anchored sidebar Data Engine | fixed readable status zone outside scrollable nav; three trading-readiness controls remain in side panel while deeper operations/diagnostics live in Maintenance | explicit Playwright anchor/control-placement tests | shell all surfaces | User | PASS |
| 39 | Consistent SVG navigation | inline bundled SVG icon system | static/responsive tests | shell visual review | User | PASS |
| 40 | Build identity always current | backend constants + renderer expected ID + Maintenance verification + package metadata | release-identity Go test + UI DOM check | Maintenance operational truth | Developer/QA | PASS |
| 41 | Maintenance Data Capabilities complete | `buildCapabilities` + build identity + diagnostics | renderer/HTTP/Go | user can distinguish active/degraded/premium-ready | User/Limitations | PASS |
| 42 | Three documentation tabs | user/developer/limitations files and renderer tabs | renderer tests | documentation matches product | all three docs | PASS |
| 43 | Performance / duplicate-code cleanup | duplicate renderer functions consolidated; TD symbol-search cache; event prewarm de-dupe | race/vet/stress/static duplicate checks | no freshness reduction | Developer/QA | PASS |
| 44 | Permanent Functionality Inventory | `renderer/qa/functionality-inventory.md` | Gate 3 inventory audit | cross-module inventory | QA | PASS |
| 45 | Gate 11 Cross-Module Integration | capability matrix + canonical correlation architecture | cross-module renderer/HTTP tests | Dashboard/Desks/Queue/Research/AI review | Developer/Limitations | PASS |
| 46 | Migration/backward compatibility | schema-tolerant load/profile migration retains old state/secrets behavior | migration tests | user workflow preserved | Developer/QA | PASS |
| 47 | Exact package integrity | Source ZIP → fresh extract → retest → platform builds → hashes | exact-source test/build/package verification | n/a | external QA + SHA | PASS |

## Cross-Module Integration Matrix

Status vocabulary: **Integrated**, **Intentionally Hidden**, **Not Relevant**.

| Capability | Dashboard / Regime | Day | Swing | Long | Queue | Discovery | Research | AI | Readiness | Validation | Settings | Maintenance | Docs |
|---|---|---|---|---|---|---|---|---|---|---|---|---|---|
| Global / Macro | Integrated | Integrated immediate | Integrated 5D/20D | Integrated 1M/3M | Integrated | Context after qualification | Integrated | Relevant evidence | Integrated | Captured | Config | Capabilities | Integrated |
| High-Impact Event Mode | Countdown + reaction | Immediate | Proximity/persistence | Persistent regime only | Integrated | Risk after qualification | Integrated | Off critical path | Integrated | Captured | Enable | Operational truth | Integrated |
| Options | Broad SPY/QQQ regime + material pulse | Immediate | Positioning/IV | Selective long-dated | Integrated | Intentionally Hidden until staged | Integrated | Relevant evidence | Integrated | Captured | Config | Capabilities | Integrated |
| Earnings | Material pulse | Immediate catalyst | Event proximity | Guidance/trend | Integrated | Qualification input | Full evidence | Relevant evidence | Integrated | Captured | Not Relevant | Capability | Integrated |
| SEC | Material pulse | Fresh material only | Compact risk | Full | Integrated | Qualification input | Full evidence | Relevant evidence | Integrated | Captured | Config | Capability | Integrated |
| Trade Readiness | Queue/overview | Integrated | Integrated | Integrated | Integrated | Not a scanner rank | Integrated | Evidence only | Primary | Captured | Not Relevant | Diagnostic | Integrated |
| Signal Validation | Diagnostic summary | Snapshot source | Snapshot source | Snapshot source | Context captured | Qualification evidence later | Diagnostic | Not formula authority | Context captured | Primary | Not Relevant | Primary | Integrated |

## Gate 11 rule
A capability fails the release if it works locally but relevant canonical evidence is stranded in one tab, duplicated by parallel fetching/state, or fails to surface useful confirmation/contradiction elsewhere. A metric can be **Intentionally Hidden** when exposing it would be horizon-inappropriate noise, but the decision must be explicit and tested.

## v14.0.3 patch traceability

| # | Requirement | Implementation | Engineering / Functional Test | Trader / Cross-Module Test | Documentation | Status |
|---|---|---|---|---|---|---|
| 48 | Global / Macro cards never overlap | collision-safe two-row responsive driver card CSS | focused Playwright geometry/overflow check + full matrix | Dashboard/Global context remains fully readable | User/QA | PASS |
| 49 | Market Instruments hover never jiggles | stable box geometry + continuous ticker animation on hover/focus | focused Playwright before/after geometry + animation-state check | common header behavior across every module | User/QA | PASS |
| 50 | Event & Data Context responsive hierarchy | stable action header + two-column responsive key/value grid | focused Playwright overlap/overflow check | Day/Swing/Long detail context | User/QA | PASS |
| 51 | SEC Intelligence key/value typography consistency | quieter values, stronger but balanced labels, controlled wrapping | focused Playwright font-size/overlap check | Long-Term/Research SEC readability | User/QA | PASS |
| 52 | Very old cached quote presentation is explicit | `very-stale` context state + warning treatment after 30 minutes | renderer/UI regression | prevents stale evidence looking current in desk context | User/Limitations | PASS |

## v14.2 notification-layout hotfix traceability

Baseline: exact v14.0.3 Stable implementation. v14.1 redesign code is intentionally excluded.

| Requirement | Implementation | Regression | Status |
|---|---|---|---|
| Center transient below-header notification | `below-header-notification[aria-hidden="false"]` uses bounded width, auto horizontal margins and centered content | Responsive UI visible-notification geometry check | PASS when release suite passes |
| Dismissed alert leaves no Market Instruments gap | `toast()` preserves `below-header-notification`; `hideToast()` restores hidden host; hidden CSS is `display:none` with zero dimensions/margins | Responsive UI post-dismiss height/gap check | PASS when release suite passes |
| Repeated Save Settings cycles do not accumulate spacer | Single host reused; base class preserved on every show/hide cycle | Renderer static lifecycle assertion + responsive geometry | PASS when release suite passes |
| Preserve v14.0.3 trading/data baseline | No formula/provider/freshness logic changed | Deterministic equivalence, trader acceptance, Go/HTTP regressions | PASS when release suite passes |

| Remove Data Engine Details row | Sidebar markup no longer renders the Details action; Maintenance remains diagnostic surface | Responsive/static DOM assertion | Cleaner earlier-style Data Engine | PASS when release suite passes |
| Reduce forced empty panel space | Natural-height overrides for summary/candidate/feed/event/empty states with v14.0.3 structure unchanged | Responsive min-height/padding assertions | Less wasted panel space without redesign | PASS when release suite passes |

## v14.3.0 Improvement Build traceability

| Requirement | Implementation / Evidence | Blocking gate |
|---|---|---|
| True app-shell toast (historical v14.3.0 approach; superseded by v14.3.2+ header notification lane) | Historical implementation only; current release uses the fixed-height global header center lane | 2, 5 |
| Save preserves location | Captures/restores window/main/Settings scroll plus focused control with `preventScroll` | 2, 4, 5 |
| Market Open Prep | 9:20–9:25 AM ET trading-day window + manual endpoint; readiness flags; live processing continues | 1, 2, 3 |
| Detailed prep state | Side Data Engine + Maintenance expose Pre-Market, Market Open and event-driven catalyst watcher status | 3, 4, 5 |
| P1/P2 provider expansion | Finnhub/Alpaca/FRED/BLS/EIA/Twelve adapters + derived canonical context + entitlement-aware capability registry | 1, 2, 3, 6 |
| Event-driven Earnings & Material Catalyst Watch | Scheduled earnings applicable date/session or material News/SEC trigger only; normal days remain dormant | 1, 2, 3, 6 |
| More Market Context correction | Five coherent evidence groups; corrected hero dead space, VIX row allocation and section/grid spacing | 4, 5 |
| Frozen deterministic scores | Context/Readiness can change; exact v14.2.2 Setup Score equivalence remains required | 1, 8, 9, 12 |
| Gate 3 audit | `v14.3.0-integration-placement-audit.md` maps every capability to relevant surfaces | 3 |
| Gate 4 audit | Same artifact audits every tab/surface for hierarchy, relevance, grouping and non-compaction | 4 |

### Permanent v14.3.0 blocking gate order
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



## v14.3.1 Completion Patch traceability

| Requirement | Implementation / evidence | Gate(s) |
|---|---|---|
| Remove sidebar black overlay/jitter | fixed pseudo-element scroll fades removed; responsive regression inspects CSS rules and sidebar state across rerender | 4, 5 |
| Complete Twelve fallback | Historical v14 Twelve fundamentals fallback (retired from the active runtime in v15); current regression protects Twelve FX plus the canonical `Finnhub → SEC → yfinance` Fundamentals route | 1, 2, 3, 6 |
| Complete Market Open reconciliation | `MarketOpenCheckpoint`, premarket snapshots, VIX/global/macro/options/horizon context | 1, 2, 3, 8, 9 |
| Real entry-zone flag | renderer Trade Readiness compares current price with deterministic Day plan entry low/high | 2, 3, 9 |
| ARMED/TRIGGERED/REACTION | catalyst watcher requires sourced release confirmation before reaction state | 1, 2, 3, 9 |
| Capability truthfulness | `capabilityStatusFromHealth` verified-positive requirement + unit test | 1, 2, 10 |
| Provider adapter fixtures | Finnhub, Alpaca and Twelve mocked endpoint tests | 1, 7 |
| Efficient slow-data cadence | FRED 12h fallback; official macro 24h fallback; startup/prep/high-impact event refresh retained | 6, 8, 10 |
| Frozen score formulas | deterministic equivalence suite remains release-blocking | 1, 8, 9, 12 |

See `v14.3.1-integration-placement-audit.md` for Gate 3/4 surface decisions.


## v14.3.2 Completion Final traceability
- Header notification feedback → global header center lane; fixed header height; overflow/truncation regression; Save-scroll preservation.
- Data Engine manual actions → same production prep/refresh/evaluate pipelines; runtime/entitlement/materiality guards; header feedback.
- Gate 3/4 decisions → `v14.3.2-integration-placement-audit.md`.
- Deterministic Day/Swing/Long Setup Score formulas → unchanged / equivalence-gated.


## v14.3.3 Stable traceability

- Dynamic live allocation → `multiFeedAllocationWithHints` / `multiFeedAllocation`, `liveSymbols`, `alpacaIEXSymbols`, `LiveCoverageState`.
- Tradable exceptions → GLD / SLV / USO pinned at highest live priority.
- Provider capacity truth → Finnhub 50 max with five-slot app reserve; Alpaca 30 max with five-slot app reserve; passive context uses snapshot pool.
- Canonical reuse → `RuntimeSnapshot.liveCoverage` is consumed by shared data-quality/context paths rather than page-specific provider logic.
- Side Data Engine → only Pre-Market Prep, Market Open Prep, and Earnings & Material Catalyst Watch retain direct controls.
- Maintenance → contextual provider/data-engine Run/Refresh/Recheck/Reconnect actions with persistent manual action lifecycle state.
- Header notification → fixed header height, full center lane, short-copy centering, long-copy dynamic fit, ellipsis only after minimum font size.
- Adaptive Data Engine geometry → frequently changing values reserve a three-line minimum, simple values use natural height, and status text is never clipped/ellipsized; status typography uses secondary hierarchy.
- Gate 3 → every new data/capability is marked Integrated / Not Relevant / Intentionally Hidden in `v14.3.3-integration-placement-audit.md`.
- Gate 5 → Content & Copy Consistency Audit covers visible UI, notifications, tooltips, errors, docs, and QA artifacts.
- Gate 10/11 → canonical build identity must match across backend/public state, renderer, native title, package metadata, final artifacts, and QA.
- Permanent exclusions and deterministic score protections remain unchanged. Historical Edge Testing remains deferred.


## v14.3.4 Stable traceability

- Gate 4/5: readiness typography, Maintenance action geometry, content-driven cache cards, Settings-header Auto-Start placement, and Decision Queue Why? default expansion are blocking UI checks.
- Gate 5/10: QA & Release History is limited to the latest five and requires status/date/summary/file metadata.
- Gate 2/3/8: Maintenance VIX Refresh/Retry reuses the canonical production VIX pipeline.
- Gate 5: long header notifications preserve the beginning of the message; font fitting precedes end ellipsis.
- Gates 10/11/12: build identity, platform metadata, exact-source freeze, and final retest are blocking.


## v14.3.7 font-fit audit

- Baseline: exact v14.3.6 Stable source.
- Scope: rendered audit of Dashboard, Day, Swing, Long, Discovery (Day/Swing/Long), Research (Overview/Fundamentals/Catalysts/Filings), AI Copilot, News, Earnings, SEC Filings, Maintenance, Documentation, and Settings across the existing responsive matrix.
- Rule: no content removal, no placement/order change, no workflow/data/scoring/provider change. Only typography fit is permitted for verified compact-panel overflow; a compact cell may lose internal padding only when required to keep the reduced text inside its existing bounds.
- Corrected: compact Decision Queue headline labels/values, Day/Swing Event & Data Context heading, Discovery signed numeric cells.
- Intentional behavior retained: Discovery long secondary description ellipsis and horizontally scrollable narrow tables are not treated as font overflow.

## v14.3.6 Stable traceability

- Baseline: exact v14.3.5 Stable source.
- SEC Form 4: `form4TransactionMeaning` + `enrichForm4` classify P as BUY, S as SELL, and compensation/exercise/tax/gift/conversion codes as OTHER; richer transaction facts flow through the canonical filing object into SEC Filings/Research/desk context without changing deterministic scores.
- Side Data Engine: the existing readiness status dot is moved after Run Now / Evaluate; no new preparation state or endpoint is introduced.
- Options Settings: provider capability card is content-driven; layout-only change, no options feed or entitlement semantics change.
- Canonical watchlist membership: `deskMembershipStrip` reads the permanent Day/Swing/Long watchlists and is reused by all desk watchlist tables and Discovery; no duplicate membership state.
- Market Regime: SPY/QQQ pills reuse canonical quote state on the second state row beside Data Confidence; no additional API calls.
- Discovery Gate 4/5 audit: explicit columns and responsive rules prevent text/text, text/button, and liquidity/spread collisions while preserving normal typography and spacing.
- Macro Rates: `macro-rates` is aggregate operational health from official Treasury core rates plus optional FRED enrichment; `fred-rates` remains source-specific provider health. Treasury can keep core rates healthy when FRED is temporarily unavailable, while FRED is never falsely advertised AVAILABLE from aggregate Treasury health. Rates Intelligence falls back from DGS10/DGS2 to UST10Y/UST2Y.
- Capability truthfulness: negative health terms are evaluated before positive substring matching so `unavailable` can never accidentally match `available`.
- Gate 3: all new/changed facts reuse canonical filing, quote, watchlist, macro metric, and health state; no duplicate fetch/subscription/scoring path.
- Deterministic Day/Swing/Long Setup Score formulas remain frozen and equivalence-gated.


### v14.3.6 final placement additions

- Market Regime state row → existing horizon regime score + canonical SPY/QQQ quote facts → state-only STRONG BULLISH/BULLISH/LEAN BULLISH/NEUTRAL/LEAN BEARISH/BEARISH/STRONG BEARISH presentation → Day/Swing/Long and common Market Regime. No scoring mutation.
- Global Market Drivers indicator → existing `GlobalMarketContext.Tone` aggregate across logical driver families → one status light + state word in the section header. Evidence counts/confidence remain the explanatory layer.
- Decision Queue four headline tiles → presentation-only font-fit correction for Score/Priority/Readiness/Research; no queue-priority or readiness logic change.
- Settings Auto-Start → presentation-only header-position correction; persisted setting and runtime startup behavior unchanged.
- Catalyst Watch persistent RUNNING regression → event state machine corrected to READY/ARMED/TRIGGERED/REACTION; RUNNING remains a short-lived manual Evaluate lifecycle state only.
- Overall Gate 4/5 fit requirement → changed components are blocking failures if text/buttons collide, values wrap outside their cards, large unexplained dead space is introduced, or responsive layout causes document-level horizontal overflow.
