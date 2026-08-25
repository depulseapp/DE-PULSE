# CURRENT Adaptive CI Convergence

**Canonical machine state:** `governance/current-state.json`  
**Certified Stable:** `v18.9.1-stable`  
**Active process-only planning slice:** #110 / `ADAPT-V19-ZERO-MISS-PLAN-001`  
**Branch:** `adapt-v19-zero-miss-plan-001`  
**Closure ledger:** `governance/work-slices/ADAPT-V19-ZERO-MISS-PLAN-001/closure.json`  
**Future hosted umbrella:** #66 — no product slice reserved.

#70/#73 remain the CI/repository control-plane foundation. #110 adds no workflow family and no new top-level G17+ gate. It extends canonical CI Fast with one cheap deterministic G2/G10 governance step:

`python3 tools/ci/v19_v20_requirement_conservation_gate.py`

The check validates the machine requirement ledger, mandatory GitHub source coverage, v18 closure guard, legal dispositions/targets, matching Roadmap/zero-miss-plan version sets and required zero-gap closure versions. Missing/unassigned/unexplained future work therefore fails early before expensive product lanes.

Exactly three canonical workflows remain authoritative: CI Fast, CI Qualified and Release. No retry/certification/promotion workflow or branch is added. #110 remains Draft during active planning, must earn canonical exact-head Fast, then become Ready and earn impact-selected Qualified on the identical SHA before expected-head merge. No Release/Stable publication belongs to #110.

The `$5` budget remains a soft cost alert, never a quality cap. Planning/metadata work should stay process-only where impact planning permits, while CI/harness changes receive the portability/backend evidence selected by Planner v3.

After #110 closure, its branch should be removed by normal branch hygiene and the next product slice must be selected only after fresh source-overlap G0/G1 from live `main`.
