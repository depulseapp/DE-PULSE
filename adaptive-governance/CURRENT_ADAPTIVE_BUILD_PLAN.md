# CURRENT Adaptive Build Plan

**Canonical machine state:** `governance/current-state.json`  
**Certified Stable:** `v18.9.1-stable`  
**Retained completed process authority:** #73 / `ADAPT-ROOT-CONVERGENCE-001` / `adapt-root-convergence-001`  
**Completed foundation:** #80 / `ADAPT-DATAHEALTH-BASELINE-001`  
**Active product work:** #81 / `ADAPT-PROVIDER-ROUTER-PRODUCTION-001` / `adapt-provider-router-production-001`  
**Parent program:** #79 / `ADAPT-PROVIDER-PRODUCTION-001`

#80 produced the executable foundation: `provider-capability-matrix.json`, `data-health-slo.json`, `provider-fetch-paths.json`, and recurrence protection. Those artifacts remain authoritative inputs to #81/#82/#83/#78/#84.

#81 required build outputs are:
1. Migrate every `MIGRATE` general provider path from `provider-fetch-paths.json` into Smart Provider Router v2.
2. Preserve provider-specific loaders and canonical cache/persistence/coalescing owners; Router v2 owns only admission/selection/fallback ordering.
3. Preserve explicit direct authorities, including SEC/EDGAR Form 4 and official public macro/reference sources.
4. Keep provider selection capability-aware and bounded using entitlement/configuration, health/circuits, freshness/usefulness evidence, latency, quota/headroom, backpressure, work tier and cost/materiality where applicable.
5. Extend existing Smart Router/Data Health regression owners for bypass recurrence, eligibility/authority ordering, rate-limit/quota pressure, bounded fallback, anti-flapping and non-secret route reasons.
6. Maintain CURRENT Adaptive Roadmap, Build Plan, Build Process and Delivery Process alignment.

Known #80 migration debt begins with Alpaca U.S. market calendar, canonical U.S. asset universe/corporate actions and Twelve Data direct global/futures acquisition. Same-provider endpoint fallback may remain inside a provider attempt when bounded; general provider choice must not bypass Router v2.

#82 follows for common scoped degradation/recovery/load shedding. Do not fold #82 into a parallel #81 health subsystem.
