# DE.PULSE — Current Handoff

**SUPERSEDES ALL PRIOR CHAT HANDOFFS**

**GitHub source of truth:** `depulseapp/DE-PULSE`  
**Certified Stable:** `v18.8.2-stable`  
**Stable baseline / PR base:** `78378889e52c2ed3e0c3458aea6fbf36efe97ab3`  
**Certified candidate:** `e51831b8269c3ae673edc93eb0ec88a0a954344f`  
**Certified fingerprint:** `a3b8851f32ef251054ac92ffdd0a9f2ed24e34b44bc45f2fa47cd97da5792247`  
**Build ID:** `v18.8.2-stable-20260820`  
**Active release line:** `v18.9.0 — TradeInsight Full Capability SHADOW Integration`  
**Active branch:** `v18.9.0-development`  
**Draft PR:** #62  
**Durable scope:** issue #61 / `ADAPT-TRADEINSIGHT-001`  
**Latest product/admission commit immediately before this handoff update:** `4424bf2af54d4c5b5894773a2e58a4aa1474b37d`  
**Latest known Fast proof before that admission-only commit:** Fast #470 at `c9402b431489c6935437fb4f66801372ab213905` — PASS  
**Latest-head qualification:** must be read live from GitHub; do not infer a PASS from an earlier SHA.

## v18.8.2 — COMPLETE / STABLE

`v18.8.2-stable` remains the immutable certified baseline. Its prior G11–G16 release evidence remains authoritative and is not superseded by the v18.9 development line.

## v18.9.0 — implemented executable scope

Issue #61 is the durable G0/G1 contract. Do **not** restart v18.9 from chat memory and do not broaden beyond that issue without updating the durable contract first.

### Canonical history / provider integration

Implemented on the existing owners only:
- Smart Provider Router v2 remains the sole executable routing authority.
- Scheduled/history requests use the one canonical dataset `Historical Bars`; mode is request context, not a second provider route.
- TradeInsight is a daily-history fallback member of that canonical route and is never presented as an intraday provider.
- Native Go Bearer-token adapter uses runtime-only `TIDATA_API_KEY` (legacy `TRADEINSIGHT_API_KEY` accepted); no secret is committed.
- Pagination is bounded; errors are classifiable/redacted; `Retry-After` is retained; runtime calls feed shared `ProviderTelemetry`.
- Adjusted daily OHLCV writes into canonical history state; weekly bars are derived from canonical daily bars.
- Dividend/split evidence merges into the existing canonical `CorporateAction` ledger and semantic duplicates preserve existing canonical truth.
- Multi-symbol history/backfill is client-side fan-out over the verified per-ticker `/ohlc` endpoint only. `historyRouteSymbols()` supplies the canonical deduplicated symbol set, the TradeInsight loop is sequential, and the provider attempt is capped at 50 symbols. No server-side bulk REST endpoint is assumed or invented.
- Commit `4424bf2af54d4c5b5894773a2e58a4aa1474b37d` therefore admits `bulk-history` as **SHADOW**, with regression coverage that locks in canonical ownership, sequential/deduplicated behavior, the 50-symbol cap, and explicit-promotion requirement.

### Congressional Trading Intelligence

Implemented and SHADOW-only:
- configured-key manual live validation on 2026-08-21 proved `GET /trading-data/v1/congress/v1/trades?limit=1` returns HTTP 200 and a `{data:[...]}` envelope;
- captured normalized schema includes amount/asset/house/member/member owner/ticker/trader/transaction date/type and stable transaction identifiers; `tx_hash` is the preferred stable ID;
- redacted/local schema fingerprint: `ee80a28688a81ac4dca5a8a47c46ea4ca6966034785813725a9c9ab7a09c9426`;
- `ResearchEngine.refreshCongressShadow()` runs only at the canonical Research alternative-evidence seam after direct SEC refresh;
- Congress evidence has no user-facing output, does not alter deterministic Day/Swing/Long truth, and an optional TradeInsight failure cannot downgrade Research readiness;
- shared provider telemetry/lifecycle truth are reused; no TradeInsight-specific persistence/router/freshness owner exists.

## Capabilities intentionally contract-gated

These are **not omitted scope**; they are explicit issue-#61 dispositions and must remain non-executable until the exact production contract is proven:

- **SEC Form 4 enrichment:** vendor availability is advertised, but the current official public `TradeInsight-Info/tidata` tree does not expose an exact production REST endpoint/response schema. Do not guess. Direct SEC/EDGAR remains authoritative even after future admission.
- **Top movers:** current official `tidata/mcp/README.md` documents MCP `get_top_movers`, but no production REST path/output schema is verified. Keep GATED; when admitted, normalize only into the existing Market Activity / Opportunity Radar candidate path.
- **Ticker/company search:** current official MCP README documents `search_ticker`, but no production REST path/output schema is verified. Keep GATED; when admitted, use only as fallback/corroboration behind canonical U.S.-equity symbol validation.
- **Generic duplicate market-price surfaces / vendor-derived scores:** FUTURE unless a configured entitlement exposes an independently useful, contract-verified capability.
- **MCP/Python SDK:** reference/developer semantics only; they are not production runtime dependencies.

Official SDK recheck on 2026-08-21 confirms `tidata/tifinance/multi.py` implements multi-ticker download by iterating per-symbol `Ticker.history()`, while `Ticker.history()` uses `/ohlc` and currently supports daily (`1d`) only. This corroborates the bounded client-side fan-out design and does not authorize a server bulk endpoint.

## Architecture and product boundaries preserved

- Smart Provider Router v2 = sole routing authority.
- Canonical freshness/recovery/cache/persistence owners remain unchanged.
- Existing multi-feed allocator remains the sole live-subscription owner.
- Existing corporate-action ledger remains the sole corporate-action owner.
- Direct SEC/EDGAR remains authoritative for Form 4.
- Opportunity Radar remains the sole mover/candidate-ranking owner.
- Canonical symbol validation remains final; U.S. Equities boundary and GLD/SLV/USO actionable exceptions remain protected.
- Adaptive Intelligence may learn provider usefulness from shared evidence/telemetry but cannot auto-promote SHADOW capability to live authority.
- Deterministic Day/Swing/Long truth is not overridden by TradeInsight evidence.
- No Execution.
- No Python or MCP production dependency.
- No new router, cache, scanner, scheduler, state store, Market Mode engine, SEC truth owner, or freshness owner was introduced.

## Qualification truth

- Fast #470 passed the exact prior head `c9402b431489c6935437fb4f66801372ab213905`.
- `4424bf2` is an admission-contract/test change after that PASS. At the time this handoff update was prepared, GitHub had not yet surfaced a latest-head commit status/workflow result for `4424bf2`; therefore **do not claim latest-head Fast/G5 PASS unless GitHub now shows it**.
- G5–G16 qualification/release proof must be earned from the current immutable candidate SHA in the normal workflow. Do not merge or release from this handoff alone.

## Exactly next action

1. Re-fetch the actual head of `v18.9.0-development` and PR #62 before changing anything.
2. Read the latest issue #61 comments and live checks.
3. If latest-head Fast is not green, run/diagnose the existing normal qualification path; do not create a duplicate workflow/branch.
4. When Fast is green, continue the existing G5–G16 qualification sequence as far as permitted **without merging or releasing**.
5. Keep Form 4, movers and ticker search gated unless an exact configured-key production REST endpoint + schema is captured/redacted and the issue admission contract is updated first.

## Resume rule

Any ChatGPT account, Claude or other assistant must read `AGENTS.md` / `CLAUDE.md`, `governance/AI-ASSISTANT-PORTABILITY-CONTRACT.md`, issue #61 and comments, this file, the four CURRENT Adaptive overlays, `release_identity.json`, both `.depulse-certification/resume/` checkpoints, PR #62, and live GitHub branch/check state before changing source. GitHub source and executable checks outrank prose. Never create a replacement v18.9 branch merely because a chat handoff is incomplete.