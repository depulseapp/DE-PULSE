# CURRENT Adaptive Delivery Process

**Canonical machine state:** `governance/current-state.json`  
**Certified Stable:** `v18.9.1-stable`  
**Completed provider onboarding:** #95 / PR #101  
**Completed continuity process:** #102 / `ADAPT-POST-STABLE-CONTINUITY-001` / PR #103 + closure PR #104  
**Active product work:** #94 / `ADAPT-PROVIDER-TELEMETRY-001` / `adapt-provider-telemetry-001`  
**Canonical closure ledger:** `governance/work-slices/ADAPT-PROVIDER-TELEMETRY-001/closure.json`  
**Parent program:** #65 / `ADAPT-PROVIDER-INTELLIGENCE-010`  
**Future hosted program:** #66 remains blocked/not started.

## #94 delivery order

1. Complete the narrow provider telemetry implementation and executable regressions on `adapt-provider-telemetry-001` without changing Router ordering/admission/lifecycle.
2. Keep all closure-ledger behavioral gaps `IMPLEMENTED_UNVERIFIED` until the final exact candidate is tested.
3. Open one Draft PR only when the source/governance packet is coherent; canonical Fast must pass on that exact head.
4. Fix only evidence-backed failures. Once Fast is green, make the same PR Ready to trigger Planner v3 impact-selected Qualified on the identical head.
5. Do not mutate source after Qualified PASS. Record immutable candidate/Fast/Qualified evidence, re-fetch live `main`, and merge only with expected-head protection.
6. Confirm main branch hygiene and continuity remain green and no Release workflow was triggered.
7. Update final closure evidence truthfully; no Stable/public SemVer release belongs to #94 alone.

## Acceptance boundary

Privileged Maintenance must show transport completion/success/median/P95 truth separately from semantic usefulness. Usefulness must remain evidence-floor guarded, restart-safe, bounded and `ADVISORY_ONLY`; normal USER/DEMO roles must not receive the new internal projection. Smart Provider Router v2 remains sole routing authority.

Completed #79/#84, #92, #95 and #102 remain inherited foundations. Direct SEC/EDGAR, canonical freshness/degradation/cache/persistence/telemetry/reconciliation/lifecycle owners, U.S. equities, GLD/SLV/USO and No Execution remain unchanged. #66 stays blocked.
