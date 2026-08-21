# DE.PULSE — Current Handoff

**SUPERSEDES ALL PRIOR CHAT HANDOFFS**

**GitHub source of truth:** `depulseapp/DE-PULSE`  
**Certified Stable:** `v18.8.2-stable`  
**Certified Stable candidate:** `e51831b8269c3ae673edc93eb0ec88a0a954344f`  
**Stable continuity / PR base:** `78378889e52c2ed3e0c3458aea6fbf36efe97ab3`  
**Engineering branch:** `v18.9.0-development`  
**PR:** #62 — Draft during release-candidate identity promotion  
**Durable scope:** issue #61 / `ADAPT-TRADEINSIGHT-001`  
**Candidate package identity:** `18.9.0` / `v18.9.0-stable-20260821`  
**Pre-RC product head:** `6cdf250b725815277e007b329a8459c1cbf78a5c` — Fast #478 PASS / Qualified #152 PASS through G10 before release-identity promotion.

## Current truth

v18.9.0 implementation is complete for every capability whose executable production contract is proven. The release-candidate promotion changes package/release identity, so the earlier G10 evidence remains implementation evidence but is not exact-head merge authority. The new RC head must earn fresh Fast and full Qualified before merge.

The immutable v18.8.2 Stable checkpoints intentionally remain anchored to v18.8.2 while v18.9.0 is an in-flight candidate. They must not be rewritten merely to make candidate version strings match.

## Implemented v18.9.0 scope

- Smart Provider Router v2 remains the sole executable provider-routing authority.
- TradeInsight configuration is runtime-secret-only through `TIDATA_API_KEY` with legacy `TRADEINSIGHT_API_KEY` compatibility.
- TradeInsight daily adjusted OHLCV participates only in the canonical `Historical Bars` route as fallback/backfill; it never claims intraday support.
- Weekly bars are derived from canonical daily bars; dividend/split evidence merges into the existing corporate-action ledger.
- `bulk-history` is SHADOW-admitted as bounded client-side fan-out over the verified per-symbol `/ohlc` endpoint. Canonical symbols are normalized/deduplicated, VIX is excluded, requests remain sequential, and the hard ceiling is 50 symbols.
- `bulk-history` admission is executable: if gated, TradeInsight collapses to one eligible symbol instead of silently performing multi-symbol fan-out.
- Congressional Trading Intelligence uses the validated `/trading-data/v1/congress/v1/trades` contract and is wired SHADOW-only into the canonical Research alternative-evidence seam after direct SEC refresh.
- Congress evidence has no user-facing/deterministic Day/Swing/Long impact and an optional TradeInsight failure cannot downgrade healthy Research readiness.
- Shared ProviderTelemetry, provider capability lifecycle, freshness, cache, persistence and canonical state owners are reused. No TradeInsight-specific router/cache/store/scanner/scheduler/Market Mode/SEC owner exists.

## Contract-gated capabilities

These remain intentionally non-executable until an exact production REST endpoint and output schema are independently verified:

- TradeInsight SEC Form 4 enrichment — Direct SEC/EDGAR remains authoritative.
- Top movers — MCP `get_top_movers` discovery does not authorize a guessed REST endpoint.
- Ticker/company search — MCP `search_ticker` discovery does not authorize a guessed REST endpoint.

These are explicit issue-#61 dispositions, not forgotten scope. Future admission must reuse the existing canonical SEC, Opportunity Radar and symbol-validation owners.

## Release-candidate promotion

The release-capable v18.9.0 head aligns:
- `release_identity.json` and `VERSION.txt` to v18.9.0 Stable package identity;
- `app_bootstrap.go` runtime version/build identity;
- renderer cache/title identity plus last-loaded `release-identity-v18.9.0.js`;
- `release/v18.9.0/release_contract.json`;
- `release/v18.9.0/run_full_certification.sh` for exact-source G12;
- this authoritative cross-assistant handoff.

The canonical three-workflow CI surface remains unchanged. `release_identity.json` makes PR #62 eligible for the existing merge-triggered `.github/workflows/release.yml` G11–G16 flow after fresh exact-head Fast + Qualified.

## Protected boundaries

US Equities Processing; GLD/SLV/USO actionable exceptions; deterministic Day/Swing/Long truth; Smart Provider Router v2 sole routing authority; direct SEC/EDGAR Form 4 authority; canonical freshness/cache/persistence/state owners; Opportunity Radar sole ranking owner; canonical symbol validation final authority; No Execution; no Python or MCP production dependency.

## Exactly one next action

Inspect the automatic **CI Fast** result on the atomic v18.9.0 release-candidate head. If and only if that exact head passes Fast, mark the same PR #62 Ready for Review once so the existing full Qualified workflow re-certifies the exact RC head before any merge.

## Resume rule

Any ChatGPT account, Codex session, Claude or human maintainer must read `AGENTS.md` / `CLAUDE.md`, `governance/AI-ASSISTANT-PORTABILITY-CONTRACT.md`, this file, issue #61/comments, PR #62, `release_identity.json`, `release/v18.9.0/release_contract.json`, both `.depulse-certification/resume/` checkpoints and live GitHub branch/check state. GitHub objects and executable evidence outrank chat memory. No upload of an old chat handoff is required.
