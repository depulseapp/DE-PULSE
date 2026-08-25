# CURRENT Adaptive Build Plan

**Canonical machine state:** `governance/current-state.json`  
**Certified Stable:** `v18.9.1-stable`  
**Completed provider/data-health program:** #79 with final #84 / PR #91  
**Completed canonical identity:** #92 / PR #93  
**Completed provider onboarding:** #95 / `ADAPT-PROVIDER-ONBOARDING-001` / PR #101  
**Completed continuity process:** #102 / `ADAPT-POST-STABLE-CONTINUITY-001` / PR #103 + closure PR #104  
**Active product work:** #94 / `ADAPT-PROVIDER-TELEMETRY-001` / `adapt-provider-telemetry-001`  
**Work-slice:** `governance/work-slices/ADAPT-PROVIDER-TELEMETRY-001/work-slice.json`  
**Closure ledger:** `governance/work-slices/ADAPT-PROVIDER-TELEMETRY-001/closure.json`  
**Parent program:** #65 / `ADAPT-PROVIDER-INTELLIGENCE-010`  
**Future hosted program:** #66 remains blocked/not started.

## Active build packet

1. Reuse `ProviderRequestDiagnostics` as the transport reliability authority; expose completed requests, success/failure ratio, success percentage, P50/median, P95 and average latency only in privileged Maintenance.
2. Reuse `buildProviderReconciliation` as semantic evidence truth. Aggregate only canonical reconciliation observations; do not create another reconciliation algorithm.
3. Keep sparse/single-source evidence `INSUFFICIENT`; use an evidence floor before emitting an agreement percentage. Stale, invalid and non-contemporaneous evidence is excluded and provider attribution is never guessed.
4. Persist one bounded aggregate/dedup derived feature through the existing `PersistenceManager`; no schema or parallel persistence/cache owner.
5. Keep all usefulness diagnostics `ADVISORY_ONLY`; do not feed Smart Provider Router v2 ranking, admission, lifecycle or promotion.
6. Prove transport projection, semantic aggregation, evidence eligibility, bounded restart safety, privileged role projection and routing invariance with executable regressions.
7. Run canonical Fast on the final exact head, then impact-selected Qualified on that identical head, then expected-head merge after fresh main re-fetch. Do not trigger Release.

## Retained owners

Completed #79/#84, #92, #95 and #102 remain inherited. Smart Provider Router v2, Data Health, freshness/degradation, persistence, transport telemetry, reconciliation, direct SEC/EDGAR, U.S. equities, GLD/SLV/USO and No Execution boundaries remain unchanged.
