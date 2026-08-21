# DE.PULSE — Current Handoff

**SUPERSEDES ALL PRIOR CHAT HANDOFFS**

**GitHub source of truth:** `depulseapp/DE-PULSE`  
**Certified Stable:** `v18.8.2-stable`  
**Certified candidate:** `e51831b8269c3ae673edc93eb0ec88a0a954344f`  
**Certified fingerprint:** `a3b8851f32ef251054ac92ffdd0a9f2ed24e34b44bc45f2fa47cd97da5792247`  
**Build ID:** `v18.8.2-stable-20260820`  
**Active release line:** `v18.9.0 — TradeInsight Full Capability SHADOW Integration`  
**Active branch:** `v18.9.0-development`  
**Durable scope:** issue #61 / `ADAPT-TRADEINSIGHT-001`  
**Latest product/test source at this checkpoint:** `0a0a04c861bd9cecd68b8120a3298495e143211b`  
**Latest Fast proof for that source:** #446 / `32491146838` — PASS

## v18.8.2 — COMPLETE / STABLE

`v18.8.2-stable` remains the immutable certified baseline. Its release evidence and G11–G16 package proof remain authoritative and are not superseded by the v18.9 development line.

## v18.9.0 — current governed state

Issue #61 is the durable G0/G1 capability-discovery and bounded-scope contract. Do **not** restart v18.9 from chat memory and do not broaden beyond that issue without updating the durable contract first.

The historical-route G2 blocker discovered during v18.9 mapping is resolved in source:
- scheduled intraday and daily/weekly history refreshes route through the single canonical dataset identifier `Historical Bars`;
- mode remains request context rather than a second provider-route name;
- TradeInsight joins the existing canonical history route after Alpaca and before Twelve Data/yfinance;
- TradeInsight is admitted only to the documented daily/weekly mode and is never presented as an intraday provider.

The stale Extreme-30 compatibility assertion was reconciled in commit `88b2dbc504ece2b9990f1698cabd7c8f234d37ed`. Fast #444 passed that exact repaired head.

TradeInsight history/corporate-action work now implemented on the same branch:
- native Go Bearer-token adapter using runtime env configuration (`TIDATA_API_KEY`, legacy alias `TRADEINSIGHT_API_KEY`); no secret is committed;
- bounded paginated REST fetch with response-size and page safety limits, Retry-After reporting and key redaction;
- daily adjusted OHLCV fallback/backfill through the canonical history owner, with weekly bars derived from canonical daily bars;
- no TradeInsight intraday assumption;
- dividend and split fields from admitted daily-history responses normalize into the existing canonical `CorporateAction` ledger;
- supplemental semantic duplicates preserve the existing canonical action rather than overwriting it;
- focused v18.9 regression tests cover auth/pagination/errors, daily-only admission, canonical route membership, adjusted history, weekly derivation and corporate-action normalization/precedence.

Product commits after the router repair:
- `2b4a5381c09a09fa12c0957be88b9ed97476ffa5` — normalize TradeInsight dividend/split evidence into the canonical corporate-action ledger;
- `0a0a04c861bd9cecd68b8120a3298495e143211b` — focused corporate-action regression coverage.

Fast #446 passed exact product/test source `0a0a04c861bd9cecd68b8120a3298495e143211b`.

## Architecture preserved

- Smart Provider Router v2 remains the sole executable routing authority.
- Canonical freshness/recovery and routed refresh remain the sole freshness/recovery owners.
- Existing multi-feed allocation remains the sole live-subscription owner.
- Existing corporate-action ledger/truth builder remains the sole corporate-action owner.
- Direct SEC/EDGAR remains authoritative for Form 4; TradeInsight may only corroborate/enrich when that capability is implemented.
- Opportunity Radar remains the sole mover/candidate-ranking owner.
- Deterministic Day/Swing/Long truth, U.S. Equities Processing, GLD/SLV/USO actionable exceptions and No Execution remain protected.
- No Python or MCP production dependency was added.

## Still open — do not claim v18.9 complete

G2/G3 remains open for capability-specific entitlement/schema readiness beyond publicly evidenced daily history. Required remaining scope from issue #61 includes:
- configured-key entitlement/schema probe by capability without leaking secrets;
- Congressional Trading Intelligence SHADOW ingestion/normalization;
- SEC Form 4 corroboration/enrichment with direct-SEC precedence and source-family de-duplication;
- selective bounded bulk-history behavior where the configured entitlement proves it;
- top-mover SHADOW evidence into existing Opportunity Radar only, but no REST endpoint may be invented from MCP-only evidence;
- ticker/company lookup fallback only after an admitted runtime endpoint/schema is proven;
- truthful provider capability/diagnostic visibility and adaptive usefulness/promotion evidence;
- missing-key/optional-provider behavior must never create a DATA DEGRADED cascade;
- full G5–G16 qualification and release proof are not yet earned.

Public vendor evidence currently proves the beta REST base, Bearer auth, daily OHLCV and `/congress/v1/trades`. The public beta page does not by itself prove production REST paths for MCP `get_top_movers` or `search_ticker`; do not guess them.

## Exactly one next action

Complete the remaining **G2/G3 TradeInsight dependency/readiness map**: locate the exact existing Congressional/SEC/Radar/symbol/telemetry owner seams and perform a configured-key, redacted entitlement/schema probe for the next admitted capability. Then implement the smallest SHADOW slice through those existing owners, beginning with congressional trades because `/congress/v1/trades` is publicly evidenced. Do not add unverified endpoint paths or a TradeInsight-specific subsystem.

## Resume rule

Any ChatGPT account, Claude or other assistant must read `AGENTS.md` / `CLAUDE.md`, `governance/AI-ASSISTANT-PORTABILITY-CONTRACT.md`, issue #61, this file, the four CURRENT Adaptive overlays, `release_identity.json`, both `.depulse-certification/resume/` checkpoints and live GitHub branch/check state before changing source. Reconcile this handoff against the actual branch head first; GitHub source and executable checks outrank prose.
