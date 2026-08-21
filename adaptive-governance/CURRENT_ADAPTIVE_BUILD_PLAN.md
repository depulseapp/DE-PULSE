# DE.PULSE — Current Adaptive Build Plan

**Immutable Stable:** `v18.8.2-stable`  
**Stable candidate:** `e51831b8269c3ae673edc93eb0ec88a0a954344f`  
**Stable fingerprint:** `a3b8851f32ef251054ac92ffdd0a9f2ed24e34b44bc45f2fa47cd97da5792247`  
**Active development branch:** none  
**Active PR:** none  
**Current line:** `v18.9.0 — G0 TradeInsight capability discovery / G1 planning`.

## v18.8.2 closure — G0–G16 PASS

Issue #57 was the sole bounded product scope. The implementation reused the existing Market Intelligence breadth universe, canonical quote freshness/recovery, existing multi-feed allocation, Smart Provider Router v2 and VIX special-index path; it added no second owner. Unavailable/degraded Market Tradeability evidence is now distinct from a genuine numeric zero.

Final release evidence:
- Fast #437 / `32435845178`: PASS on exact recovery/release head `66e59e4e5f803ca53520797e5eb6e9d3fe72e84c`;
- Qualified #151 / `32435920048`: PASS across backend/full/race/randomized, renderer, Chrome, WebKit and CI/harness;
- Release #31 / `32436189650`: G11–G16 PASS;
- certified candidate `e51831b8269c3ae673edc93eb0ec88a0a954344f`;
- `v18.8.2-stable` exact tag match and no-rebuild publication PASS.

Release #30 / `32435511692` is retained as classified historical evidence: G11 PASS, G12 harness failure caused by a stale README presentation-heading assertion. Recovery PR #60 changed release-harness/governance only, requalified fully, then Release #31 closed the release. No product/runtime behavior or package identity changed in that recovery.

## v18.9.0 G0/G1 plan

Before any product-source implementation:
1. re-read current GitHub truth and provider/routing governance;
2. enumerate the complete configured TradeInsight beta capability surface and entitlement/rights;
3. map each capability to a real DE.PULSE consumer/purpose;
4. classify it `USE`, `CORROBORATE`, `FALLBACK`, `STORE_FOR_HISTORY`, `FUTURE`, or `NOT_USEFUL`;
5. identify source-independence/double-counting risks, freshness, cache/retention, rate-limit/budget and Market Mode implications;
6. preserve Smart Provider Router v2 as sole executable routing authority;
7. propose one bounded G1 scope with explicit exclusions and tests before changing product source.

Mandatory minimum candidate roles are Congressional Trading, SEC Form 4 enrichment secondary to direct SEC, and historical OHLCV fallback/backfill. They are not a cap on useful discovery.

## Protected boundaries

Smart Provider Router v2 sole routing authority; canonical freshness/recovery sole freshness owner; existing multi-feed allocator sole subscription owner; BroadSnapshotBroker canonical reuse owner; deterministic Day/Swing/Long; U.S. Equities Processing; GLD/SLV/USO actionable exceptions; No Execution.

## Exactly one next action

Perform the v18.9.0 G0 exact-baseline / TradeInsight full-capability discovery and report findings plus bounded proposed G1 scope before product-source implementation.
