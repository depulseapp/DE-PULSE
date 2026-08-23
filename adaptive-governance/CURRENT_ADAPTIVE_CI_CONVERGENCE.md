# DE.PULSE — Current Adaptive CI / Repository Convergence

**Canonical machine state:** `governance/current-state.json`  
**Certified Stable:** `v18.9.1-stable`  
**Active work slice:** #73 / `ADAPT-ROOT-CONVERGENCE-001`  
**Active branch:** `adapt-root-convergence-001`  
**Closure ledger:** `governance/work-slices/ADAPT-ROOT-CONVERGENCE-001/closure.json`

#70 CI convergence is complete and remains the control-plane foundation. #73 extends repository hygiene through `tools/ci/root_convergence_gate.py`, which assigns every tracked root file exactly one `KEEP`, `MOVE`, `CONSOLIDATE` or `DELETE` disposition and progressively tightens enforcement by implementation phase. No parallel workflow family is introduced.

`WAIVER-GITHUB-MAIN-PROTECTION-001` remains the bounded current-plan exception under the completed #70 work slice; it does not claim technical branch protection. Exact-head Fast/Qualified remains mandatory for #73 closure.
