# CURRENT Adaptive CI Convergence

**Canonical machine state:** `governance/current-state.json`  
**Certified Stable:** `v18.9.1-stable`  
**Completed provider onboarding:** #95 / PR #101  
**Completed continuity process:** #102 / `ADAPT-POST-STABLE-CONTINUITY-001` / PR #103 + closure PR #104  
**Active product work:** #94 / `ADAPT-PROVIDER-TELEMETRY-001` / `adapt-provider-telemetry-001`  
**Canonical closure ledger:** `governance/work-slices/ADAPT-PROVIDER-TELEMETRY-001/closure.json`  
**Parent program:** #65 / `ADAPT-PROVIDER-INTELLIGENCE-010`  
**Future hosted program:** #66 remains blocked/not started.

#70/#73 remain the completed CI/repository control-plane foundation. #94 uses only the canonical CI Fast and CI Qualified workflows; release identity remains `v18.9.1-stable`, so no Release workflow belongs to this slice.

## Required #94 evidence

- exact candidate identity and source fingerprint;
- canonical Fast exact-head PASS, including current-state/work-slice/source-health, Go formatting/vet/tests and impacted renderer checks selected by the workflow;
- Planner v3 impact selection on the same candidate;
- impact-selected Qualified PASS on the identical head, including every selected evidence owner;
- source proof that semantic usefulness remains `ADVISORY_ONLY` and Smart Provider Router v2 has no usefulness dependency;
- expected-head protected merge after a fresh `main` fetch;
- post-merge branch hygiene/continuity remains green and no Release workflow is triggered.

The active closure ledger remains fail-closed until these immutable GitHub objects exist. No gate waiver, retry branch, direct-main patch or temporary workflow is permitted.

Completed #80-#84, #92, #95 and #102 remain inherited executable authority. Smart Provider Router v2 remains sole general routing/admission authority; direct SEC/EDGAR, canonical freshness/degradation/cache/persistence/telemetry/reconciliation/lifecycle, U.S. equities, GLD/SLV/USO and No Execution remain unchanged.
