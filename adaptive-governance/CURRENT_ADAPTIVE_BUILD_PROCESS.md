# CURRENT Adaptive Build Process

**Canonical machine state:** `governance/current-state.json`  
**Certified Stable:** `v18.9.1-stable`  
**Retained completed process authority:** #73 / `ADAPT-ROOT-CONVERGENCE-001` / `adapt-root-convergence-001`  
**Completed foundation:** #80 / `ADAPT-DATAHEALTH-BASELINE-001`  
**Completed Router adoption:** #81 / `ADAPT-PROVIDER-ROUTER-PRODUCTION-001` / PR #87 / merge `1870dd3881dbe7f6463f242e35fdc19e70d9ae15`  
**Active product work:** #82 / `ADAPT-DATAHEALTH-RUNTIME-001` / `adapt-datahealth-runtime-001`  
**Parent program:** #79 / `ADAPT-PROVIDER-PRODUCTION-001`

For #82 and all remaining Data Health slices, implementation remains executable-first and canonical-owner-first:
- reuse Smart Provider Router v2 as the sole general market-data routing/admission owner; #82 must not re-route around it;
- reuse canonical freshness, cache, persistence, telemetry, state, provider workload and validation owners before adding anything;
- REUSE/CONSOLIDATE/REFACTOR before ADD, and fail closed on a new parallel provider-specific health subsystem;
- derive degradation from affected required evidence and scope it by capability/consumer/symbol/session before wider escalation;
- reuse valid warm state only when canonical freshness policy allows; stale/unknown evidence must never be relabeled healthy;
- make eligible fallback/revalidation/recovery automatic and hysteresis-protected;
- protect critical decision-support evidence before optional/background work; shed lower-value work before core truth;
- bound scanner/prep/event/research/background fan-out and distinguish local overload from provider/capability failure;
- preserve explicit direct-authority/public paths rather than score-overriding deterministic truth boundaries;
- use focused fault/regression evidence before canonical Fast exact-head PASS, then deliberate impact-selected Qualified exact-head PASS;
- do not add temporary/permanent workflow families or weaken G0–G16/source-health/architecture gates.

The governed program sequence token remains `#81/#82/#83/#78/#84`; #81 is complete, so current execution is #82 → #83 + #78 → #84. Documentation alone never verifies closure.
