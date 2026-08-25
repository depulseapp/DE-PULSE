# CURRENT Adaptive Roadmap

**Canonical machine state:** `governance/current-state.json`  
**Certified Stable:** `v18.9.1-stable`  
**Completed provider/data-health program:** #79 with final #84 / PR #91  
**Completed canonical identity:** #92 / PR #93  
**Completed provider onboarding:** #95 / `ADAPT-PROVIDER-ONBOARDING-001` / PR #101  
**Completed continuity process:** #102 / `ADAPT-POST-STABLE-CONTINUITY-001` / PR #103 + closure PR #104  
**Active product work:** #94 / `ADAPT-PROVIDER-TELEMETRY-001` / `adapt-provider-telemetry-001`  
**Closure ledger:** `governance/work-slices/ADAPT-PROVIDER-TELEMETRY-001/closure.json`  
**Parent program:** #65 / `ADAPT-PROVIDER-INTELLIGENCE-010`  
**Future hosted program:** #66 remains blocked/not started.

## Current roadmap decision

The fresh post-#102 #65 semantic-overlap audit re-proved #94 as the next real dependency-ordered residual. Existing transport telemetry already computes completion, success percentage and P50/P95; #94 projects those diagnostics into privileged Maintenance without a duplicate calculation owner. Existing provider reconciliation already owns AGREED/CONFLICT/SINGLE SOURCE/STALE truth; #94 derives bounded semantic-usefulness aggregates only from those decisions and persists only compact aggregate/dedup state through the canonical persistence manager.

Semantic usefulness remains `ADVISORY_ONLY`. Sparse or single-source evidence stays `INSUFFICIENT`; stale, invalid and non-contemporaneous evidence cannot manufacture a provider score. Smart Provider Router v2 ordering/admission/lifecycle remains unchanged.

## Retained foundation and boundaries

The inherited sequence **#80 baseline -> #81 Router v2 adoption + #82 common health/recovery -> #83 lifecycle + #78 TradeInsight admission -> #84 zero-gap closure -> #92 canonical identity -> #95 onboarding/adaptation** remains executable authority. Smart Provider Router v2 remains the sole general routing/admission owner; direct SEC/EDGAR remains Form 4 authority; canonical freshness/degradation/cache/persistence/telemetry/reconciliation/lifecycle owners remain unchanged; U.S. equities only, GLD/SLV/USO actionable exceptions and No Execution remain permanent.

## Delivery boundary

#94 must earn exact-head Fast and impact-selected Qualified on one unchanged candidate before expected-head merge. No Stable/public SemVer release belongs to #94 alone. #66 remains future-blocked.
