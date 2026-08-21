# DE.PULSE — Current Adaptive Roadmap

**Certified Stable:** `v18.9.0-stable`  
**Certified candidate:** `9ea81cddae4875ae15d3719ca028519a36c597b6`  
**Certified fingerprint:** `a8719090c341c874dbd1279cc31ad98e84075d5701c46a800bf951340780ecb9`  
**Build ID:** `v18.9.0-stable-20260821`  
**Active corrective program:** issue #65 / `ADAPT-PROVIDER-INTELLIGENCE-010`  
**Immediate blocker / next patch:** issue #64 / `ADAPT-RUNTIME-CRASH-001` -> `v18.9.1`.

## v18.9.0 — COMPLETE / IMMUTABLE STABLE

Issue #61 / `ADAPT-TRADEINSIGHT-001` is closed completed. Exact source head `9e86b5e731f7a585cc77c1521f3639fc7a208efc` passed Fast #481 and Qualified #153. Merged candidate `9ea81cddae4875ae15d3719ca028519a36c597b6` passed canonical Release #32 through G11–G16. Durable release evidence is `release/v18.9.0/stable-evidence-manifest.json`.

The post-release audit found that the architecture is sound but v18.9.0 did not fully realize the intended adaptive multi-provider/product UX. Those findings are not retroactively added to the immutable Stable artifact; they are governed by #65 as small v18.9.x patches.

## Permanent Small-Patch Operating Rule

DE.PULSE prefers **many small, dependency-ordered, complete patches** over heavy multi-domain builds.

- One primary responsibility per patch; only tightly coupled support work may accompany it.
- No stability + routing + provider-expansion + UX bundles.
- Each patch starts from an exact G0 baseline and immutable G1 scope.
- Before the next patch starts, the current patch gets implementation-miss review, focused regression proof, runtime/browser proof where applicable, open-issue reconciliation and durable handoff.
- A known implementation miss must be fixed in-scope or explicitly registered against a later patch before closure; it may not disappear into chat memory.
- One development branch + one PR per patch; no retry/certification branch families and no duplicate CI runs.
- G0–G16 remains the only release model.

## v18.9.x ordered patch roadmap

1. **v18.9.1 — Runtime crash corrective ONLY** — #64. Diagnose/fix the real macOS Apple Silicon SIGABRT from evidence/reproduction; preserve user state/API keys; add lifecycle regression and actual packaged macOS proof.
2. **v18.9.2 — TradeInsight Settings/API-key UX ONLY.** Existing Data Provider Settings/secret owner; masked Save/Test/Clear; truthful status; environment override only as developer/runtime fallback.
3. **v18.9.3 — Coverage-aware Smart Provider Router core ONLY.** Upgrade first-success behavior to requirement/coverage-aware fulfillment; cache/gap first; merge/provenance/re-evaluate; validation lifecycle separated from serving role.
4. **v18.9.4 — Canonical company identity + all-desk presentation ONLY.** Shared symbol/company identity; `APP - AppLovin : In Entry Zone` with symbol-only fallback; reused by desks/Research/Discovery/Add Symbol.
5. **v18.9.5 — Market Data Modes + capability diagnostics ONLY.** Behavior-oriented Adaptive modes rather than provider-brand modes; capability-level source/freshness/coverage diagnostics; no separate TradeInsight mode.
6. **v18.9.6 — TradeInsight SEC Form 4 enrichment ONLY.** Contract-validated SHADOW-first enrichment/corroboration; direct SEC/EDGAR authoritative; source-family de-duplication.
7. **v18.9.7 — TradeInsight ticker/company search ONLY.** Contract-validated fallback/corroboration through canonical symbol validation/company identity; U.S.-equity boundary final.
8. **v18.9.8 — TradeInsight movers/ranking evidence ONLY.** Contract-validated candidate evidence into Opportunity Radar; existing scanner/ranker remains canonical; SHADOW-first usefulness proof.
9. **v18.9.9 — Remaining useful TradeInsight capability sweep ONLY.** Every useful entitlement gets explicit disposition and consumer; retest Congress/history/corporate actions under coverage-aware routing; no invented endpoints or Python/MCP production dependency.
10. **v18.9.10 — Provider efficiency + Adaptive Intelligence telemetry ONLY.** Coverage filled, gaps, calls avoided, cache hits, provider usefulness, latency/rate-limit/freshness/conflict telemetry, bounded fan-out and runtime-load proof.
11. **v18.9.11 — Whole v18.9.x professional closure audit ONLY.** End-to-end implementation-miss audit, #57/#64 regression, deterministic Day/Swing/Long equivalence, actual macOS/Windows packaged proof, Adaptive Intelligence Scorecard and zero unexplained carry-forward/orphan useful capability/duplicate owner.

## Permanent adaptive provider architecture

DE.PULSE operates as:

`consumer requirement -> current cache/state -> exact missing coverage -> eligible-provider ranking -> targeted acquisition -> canonical merge/provenance -> coverage re-evaluation -> next provider only if still needed -> synthesized consumer state`

A successful provider response is not enough to stop if required coverage/freshness/fields/quality remain incomplete. No fixed global chain such as `Alpaca -> TradeInsight -> Twelve Data -> yfinance` is the decision model; static ordering is at most a prior/tiebreaker. Smart Provider Router v2 remains sole executable routing authority.

## v19+ inheritance

v19 Professional Data Infrastructure and v20 Adaptive Intelligence maturation inherit: small-patch discipline, coverage-aware fulfillment, validation-lifecycle/serving-role separation, canonical company identity, behavior-oriented Market Data Modes, provider-usefulness telemetry and full useful-capability disposition. Future providers plug into these contracts rather than creating fixed fallback chains or parallel subsystems.

Permanent constraints: U.S. Equities Processing, No Execution, Smart Provider Router v2 sole routing owner, canonical freshness/recovery sole freshness owner, existing multi-feed allocator sole subscription owner, BroadSnapshotBroker canonical reuse owner, direct SEC/EDGAR authoritative, GLD/SLV/USO actionable tradable exceptions and deterministic Day/Swing/Long truth protected.

## Exactly one next action

Perform issue #64 / v18.9.1 G0 crash diagnosis from concrete macOS evidence or deterministic reproduction. Do not start v18.9.2 until v18.9.1 is closed with truthful evidence or the crash is proven external/non-product.
