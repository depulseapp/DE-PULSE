# CURRENT Adaptive Build Plan

**Canonical machine state:** `governance/current-state.json`  
**Certified Stable:** `v18.9.1-stable`  
**Retained completed process authority:** #73 / `ADAPT-ROOT-CONVERGENCE-001` / `adapt-root-convergence-001`  
**Active product work:** #80 / `ADAPT-DATAHEALTH-BASELINE-001` / `adapt-datahealth-baseline-001`  
**Parent program:** #79 / `ADAPT-PROVIDER-PRODUCTION-001`

#80 is the current build foundation. Required executable outputs are:
1. `governance/data-health/provider-capability-matrix.json` — exhaustive provider/authoritative-source and capability/dataset ownership.
2. `governance/data-health/data-health-slo.json` — healthy coverage, freshness, fallback, degradation, recovery, quota and consumer-impact contracts.
3. `governance/data-health/provider-fetch-paths.json` — every production external fetch/runtime path classified `MIGRATE`, `DIRECT_AUTHORITY` or justified `N/A`.
4. A canonical CI/source-health recurrence check rejecting new unclassified providers/capabilities/fetch paths.
5. Explicit propagation into Adaptive Roadmap, Build Plan, Build Process and Delivery Process.

Inventory must cover all currently active providers and authoritative/public sources, including the U.S.-equities stack, macro/reference sources, AI providers where they consume external provider capability, direct SEC/EDGAR authority, CBOE VIX reference/history, U.S. Treasury/BEA/BLS/EIA official macro sources, TWSE official close, and any additional runtime source found by executable source audit.

#80 records remediation disposition; #81 performs general Router v2 adoption fixes and #82 common health/recovery/load-shedding behavior. Do not prematurely implement a parallel router or provider-specific health stack in #80.
