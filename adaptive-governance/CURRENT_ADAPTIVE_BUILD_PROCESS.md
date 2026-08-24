# CURRENT Adaptive Build Process

**Canonical machine state:** `governance/current-state.json`  
**Certified Stable:** `v18.9.1-stable`  
**Retained completed process authority:** #73 / `ADAPT-ROOT-CONVERGENCE-001` / `adapt-root-convergence-001`  
**Completed foundation:** #80 / `ADAPT-DATAHEALTH-BASELINE-001`  
**Active product work:** #81 / `ADAPT-PROVIDER-ROUTER-PRODUCTION-001` / `adapt-provider-router-production-001`  
**Parent program:** #79 / `ADAPT-PROVIDER-PRODUCTION-001`

For #81 and all child Data Health slices, implementation remains executable-first and canonical-owner-first:
- use the #80 provider/capability/fetch-path inventory as the migration ledger; do not re-invent provider scope;
- reuse Smart Provider Router v2 as the sole general market-data routing/admission owner;
- keep provider-specific HTTP/normalization inside existing loaders while Router v2 selects/adopts attempts;
- preserve explicit direct-authority/public paths rather than score-overriding deterministic truth boundaries;
- reuse canonical freshness, cache, persistence, telemetry, state and validation owners;
- preserve existing cache/single-flight/coalescing before adding fetches; no blind provider fan-out;
- fail closed on unclassified new providers/capabilities/runtime external-fetch paths and on new general bypasses;
- keep degradation/recovery/load shedding scoped to #82 rather than building a parallel #81 health stack;
- use local/static/focused evidence before canonical Fast exact-head PASS, then deliberate impact-selected Qualified exact-head PASS;
- do not add temporary/permanent workflow families or weaken G0–G16/source-health/architecture gates.

Every #80 `MIGRATE` row must either be executable through Router v2 in #81 or receive a new evidence-backed authority classification. The remaining sequence is #81/#82/#83/#78/#84; documentation alone never verifies closure.
