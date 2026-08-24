# CURRENT Adaptive Build Plan

**Canonical machine state:** `governance/current-state.json`  
**Certified Stable:** `v18.9.1-stable`  
**Completed foundation:** #80 / `ADAPT-DATAHEALTH-BASELINE-001`  
**Completed Router adoption:** #81 / PR #87 / merge `1870dd3881dbe7f6463f242e35fdc19e70d9ae15`  
**Completed runtime Data Health:** #82 / Fast #894 / Qualified #187 / PR #88 / merge `4882b6d53c0c34463239faae752b86de377fb19a`  
**Active product work:** #83 / `ADAPT-PROVIDER-LIFECYCLE-001` / `adapt-provider-lifecycle-001`  
**Parent program:** #79 / `ADAPT-PROVIDER-PRODUCTION-001`

#80 produced executable `provider-capability-matrix.json`, `data-health-slo.json` and `provider-fetch-paths.json`; #81 routed every general MIGRATE path through Smart Provider Router v2; #82 completed canonical scoped degradation, warm/fallback truth, recovery hysteresis and bounded load shedding.

#83 required build outputs are:
1. Reuse Router v2 capability state, ProviderTelemetry, canonical freshness and validation evidence in one common capability-readiness evaluator.
2. Govern the lifecycle vocabulary as `SHADOW → VALIDATED → APPROVED → PRODUCTION` with explicit promotion records only.
3. Compute deterministic `READY_FOR_PROMOTION` from observation depth, success/error/auth stability, freshness, schema/semantic integrity, p95 latency, 401/403/429/5xx, quota/headroom, fallback correctness, corroboration/disagreement, provenance independence, consumer utility/outcomes and truth-boundary compliance where measurable.
4. Keep runtime suppression/cooldown/probe/recovery automatic through existing Router v2/Data Health owners without changing governed lifecycle.
5. Classify all 26 #80 provider/capability rows with lifecycle, authority, evidence status and readiness/N/A.
6. Keep direct SEC/EDGAR and other first-party authorities outside rank-promotion semantics.
7. Make TradeInsight consume this framework at SHADOW; #78 remains the separate validation/promotion adoption packet.
8. Keep lifecycle/readiness diagnostics non-secret and reason-bearing.
9. Maintain Adaptive Roadmap, Build Plan, Build Process and Delivery Process alignment.

The prior sequence `#81/#82/#83/#78/#84` is advanced through completed #81 and #82; remaining execution is **#83 → #78 → #84**. No parallel lifecycle/health subsystem and no routing-authority change is permitted.
