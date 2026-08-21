# DE.PULSE — Current Handoff

**SUPERSEDES ALL PRIOR CHAT HANDOFFS**

**GitHub source of truth:** `depulseapp/DE-PULSE`  
**Certified Stable:** `v18.8.2-stable`  
**Certified candidate:** `e51831b8269c3ae673edc93eb0ec88a0a954344f`  
**Certified fingerprint:** `a3b8851f32ef251054ac92ffdd0a9f2ed24e34b44bc45f2fa47cd97da5792247`  
**Build ID:** `v18.8.2-stable-20260820`  
**Qualified source head:** `66e59e4e5f803ca53520797e5eb6e9d3fe72e84c`  
**Final Fast:** #437 / `32435845178` — PASS  
**Final Qualified:** #151 / `32435920048` — PASS  
**Final Release:** #31 / `32436189650` — G11–G16 PASS  
**Next release line:** `v18.9.0 — TradeInsight Full Capability SHADOW Integration`

## v18.8.2 — COMPLETE / STABLE

Issue #57 / `ADAPT-FRESHNESS-001 REOPENED` is resolved by v18.8.2 Stable. The bounded repair makes the existing canonical Market Intelligence breadth universe participate in quote freshness/recovery accountability and prevents degraded/unavailable Market Tradeability evidence from being presented as an observed numeric zero.

Protected architecture remains unchanged: Smart Provider Router v2 is sole routing authority; canonical freshness/recovery + routed refresh are sole recovery owners; existing multi-feed allocation remains sole subscription owner; BroadSnapshotBroker remains the canonical reuse owner; deterministic Day/Swing/Long, U.S. Equities Processing, GLD/SLV/USO actionable exceptions and No Execution are preserved.

## Qualification and release history

Product PR #59 carried the bounded runtime/test repair. Its release-capable exact head `186dd18bcd33a2d891b3df738478ba88cf7b98b6` passed Fast #435 / `32434635563` and full Qualified #150 / `32434742951`, then squash-merged as `d855607426bc56372656d3b0baad67611aae7a96`.

Automatic Release #30 / `32435511692` passed G11 and stopped at G12 on a stale literal README presentation-heading assertion in `version_consistency_test.py`. This was classified as **CI/release-harness failure**, not product/runtime failure. It did not reach native packaging or publication.

Release-harness recovery PR #60 followed the repository's existing post-merge harness-recovery contract: the same canonical `v18.8.2-development` line was recreated from the failed merged candidate, only harness/governance files changed, and recovery head `66e59e4e5f803ca53520797e5eb6e9d3fe72e84c` earned fresh Fast #437 and full Qualified #151 before merge.

Release #31 / `32436189650` then passed G11 exact-source provenance, G12 full certification, macOS Apple Silicon and Windows x64 G13/G14 actual packaged-runtime audits, G15 Release Assurance, same-run exact-artifact publication with **no rebuild**, and G16 durable handoff evidence.

The immutable tag `v18.8.2-stable` points exactly to certified candidate `e51831b8269c3ae673edc93eb0ec88a0a954344f`. Branch hygiene removed `v18.8.2-development` after merge. Durable evidence is in `release/v18.8.2/stable-evidence-manifest.json` and both `.depulse-certification/resume/` checkpoints.

## v18.9.0 mandatory entry

`ADAPT-TRADEINSIGHT-001` is next. G0 must first enumerate the complete capability surface available to the configured TradeInsight beta account/API and classify every useful capability by purpose. Congressional Trading, SEC Form 4 enrichment secondary to direct SEC, and historical OHLCV fallback/backfill are mandatory minimum roles, not a cap.

Disposition each capability as `USE`, `CORROBORATE`, `FALLBACK`, `STORE_FOR_HISTORY`, `FUTURE`, or `NOT_USEFUL`. Any executable use must flow through Smart Provider Router v2 with entitlement/rights, freshness, budget, cache/retention, disagreement, Market Mode, SHADOW telemetry, promotion and graceful-degradation rules. Full capability discovery never means blindly calling every endpoint.

## Exactly one next action

Perform **v18.9.0 G0 exact-baseline / TradeInsight full-capability discovery** from current GitHub truth and report the findings plus a bounded proposed G1 scope before changing product source. Do not begin v18.9.0 implementation before G0 diagnosis and G1 freeze.

## Resume rule

Any ChatGPT account, Claude or other assistant must read `AGENTS.md` / `CLAUDE.md`, `governance/AI-ASSISTANT-PORTABILITY-CONTRACT.md`, this file, the four CURRENT Adaptive overlays, `release_identity.json`, both `.depulse-certification/resume/` checkpoints, `release/v18.8.2/stable-evidence-manifest.json`, the v18.8.2 release contract, issue #57 history, PRs #59/#60 and live GitHub branch/workflow/tag state before changing source. GitHub is the durable source of truth; never resume from model memory alone.
