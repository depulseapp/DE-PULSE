# DE.PULSE v15.0.0 — 43-Item Requirement Traceability

Baseline: frozen DE.PULSE v14.3.7 Stable source. A requirement is PASS only when implemented and exercised by its cited validation; final Stable status is assigned only after the frozen-source retest and package integrity gates pass.

| # | Requirement | Implementation | Validation | Status |
|---:|---|---|---|---|
| 1 | Protect v14.3.7 Stable baseline | v15 build identity and exact-source freeze workflow | version + exact-source validation | PASS |
| 2 | Fix VIX stale classification | Twelve → yfinance → CBOE route; dataset-aware VIX freshness | Go VIX freshness/router tests + responsive diagnostics | PASS |
| 3 | Rebuild Data Freshness | FreshnessDiagnostic/Summary + Maintenance rendering/actions | Go + renderer + responsive tests | PASS |
| 4 | Build Provider Router | ProviderRouterSnapshot and centralized route registry | Go route exactness + renderer tests | PASS |
| 5 | Approved routing chains | US equities/VIX/history/news/earnings/fundamentals/SEC/macro routes | Go route exactness tests | PASS |
| 6 | Feed Capability Registry | v14.3 registry extended for Twelve/yfinance/CBOE/Marketaux | Go/renderer registry validation | PASS |
| 7 | Prevent duplicate API calls | shared engine/master-symbol fan-out retained and router centralized | cross-module/reuse tests | PASS |
| 8 | Intelligent provider failover | provider health/fallback route state | Go provider-router tests | PASS |
| 9 | Provider circuit breakers | 3-failure open / timed recovery | Go circuit breaker tests | PASS |
| 10 | Add yfinance fallback only | VIX/history/earnings/fundamentals recovery adapters | Go routing + source review | PASS |
| 11 | Add CBOE for VIX | official VIX close/history fallback/validation | Go VIX route tests | PASS |
| 12 | Strengthen Twelve Data | VIX primary + history recovery + global role | Go routing/capability tests | PASS |
| 13 | Keep core providers | Alpaca/Finnhub/Twelve/Marketaux/SEC/FRED registry | registry/UI validation | PASS |
| 14 | Demote Alpha Vantage | not used in active v15 primary routes; emergency legacy only | route audit | PASS |
| 15 | Do not add Massive/FMP | no adapters/routes added | source audit | PASS |
| 16 | Provider Router in Data Engine | Maintenance Provider Router panel with health/circuit/fallback | responsive renderer test | PASS |
| 17 | SPY/QQQ ticker font = price | v15 CSS regime quote typography | responsive layout audit | PASS |
| 18 | Right-align regime/header panels | v15 regime-grid/pill alignment CSS | responsive layout audit | PASS |
| 19 | Center Action / Score | shared header/cell centering CSS | focused responsive test | PASS |
| 20 | Fix table heading placement | shared table grid/column alignment CSS | responsive layout audit | PASS |
| 21 | Global Market Driver Last Updated | actual runtime refresh-based label | renderer + responsive test | PASS |
| 22 | Desk clicks preserve scroll | in-place membership mutation with scroll capture/restore | responsive interaction test | PASS |
| 23 | Visible desk membership state | aria-pressed + active/inactive button styling | renderer/responsive tests | PASS |
| 24 | Cross-desk add controls | POST /api/desk/membership toggle API | Go + interaction tests | PASS |
| 25 | Canonical desk membership | shared dedicated watchlist state helpers | Go membership tests | PASS |
| 26 | Cross-module desk sync | state broadcast/bootstrap redraw; no duplicate membership | Go + HTTP + renderer tests | PASS |
| 27 | Master Market Symbols panel | compact full-width dashboard panel | responsive renderer test | PASS |
| 28 | Priority placement | panel after Decision Queue and before lower-priority catalyst detail | placement test | PASS |
| 29 | Truthful Master status dots | uses canonical symbolDataRecord/quote state | renderer logic test | PASS |
| 30 | Global ticker removal | /api/master-symbol/remove removes all three desks | Go + HTTP tests | PASS |
| 31 | Undo global removal | /api/master-symbol/restore + toast Undo | Go + renderer tests | PASS |
| 32 | Preserve context on removal | scroll capture/restore around mutation | interaction test | PASS |
| 33 | Responsive Master panel | wrap/overflow CSS | 15-viewport matrix | PASS |
| 34 | Data Freshness release gate | v15 freshness tests included in blocking gate | Go/renderer/responsive gate | PASS |
| 35 | Provider Router release gate | route/failure/circuit/recovery validation | Go + HTTP gate | PASS |
| 36 | Cross-module integration gate | membership/master/provider/SEC correlation | Go + HTTP + renderer gate | PASS |
| 37 | Layout Integrity gate | 18 surfaces × 15 viewports + focused checks | responsive matrix | PASS |
| 38 | No overlap/clipping | blocking overflow/control geometry checks | responsive matrix | PASS |
| 39 | Difficult UI states | fixtures include stale/delayed/empty/long/expanded states | renderer/responsive tests | PASS |
| 40 | Supported layout validation | MacBook/desktop/tablet + Windows scaling/Retina DPR | responsive matrix | PASS |
| 41 | Fix row hover jitter | stable row geometry hover CSS and non-nested controls | focused hover geometry test | PASS |
| 42 | Desk toggle/removal rule | last-desk protection; multi-desk remove; add inactive | Go + HTTP + renderer tests | PASS |
| 43 | SEC Intelligence regression | Form 4 BUY/SELL/OTHER + rich recent transaction rendering | Go classification + Research/Filings/Long visibility tests | PASS |

## Blocking release rule

All 43 items, unit/race/vet, deterministic trading equivalence, professional trader acceptance, HTTP/cross-module workflows, content/version consistency, responsive/layout integrity, exact-source freeze/retest, and platform package verification must pass before `STABLE · PASS` is published.
