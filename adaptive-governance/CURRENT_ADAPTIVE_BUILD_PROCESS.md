# CURRENT Adaptive Build Process

**Canonical machine state:** `governance/current-state.json`  
**Certified Stable:** `v18.9.1-stable`  
**Retained completed process authority:** #73 / `ADAPT-ROOT-CONVERGENCE-001` / `adapt-root-convergence-001`  
**Active product work:** #80 / `ADAPT-DATAHEALTH-BASELINE-001` / `adapt-datahealth-baseline-001`  
**Parent program:** #79 / `ADAPT-PROVIDER-PRODUCTION-001`

For #80 and all child Data Health slices, implementation must be executable-first and canonical-owner-first:
- source-audit every provider/capability/external fetch path before claiming coverage;
- reuse Smart Provider Router v2 for general market-data routing/admission;
- classify legitimate direct-authority/public paths explicitly rather than forcing rank-based substitution;
- reuse canonical freshness, cache, persistence, telemetry, state and validation owners;
- fail closed on unclassified new providers/capabilities/runtime external-fetch paths;
- keep degradation scoped to affected capability/consumer and never suppress a genuine evidence gap merely to appear healthy;
- use local/static/focused evidence before canonical Fast, then deliberate impact-selected Qualified;
- do not add temporary/permanent workflow families or weaken G0–G16/source-health/architecture gates.

Every known gap discovered during #80 must be either resolved in #80 acceptance or assigned an executable disposition to #81/#82/#83/#78/#84. Documentation alone never verifies a closure gap.
