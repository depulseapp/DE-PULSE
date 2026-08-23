# DE.PULSE — Current Handoff

**SUPERSEDES ALL PRIOR CHAT HANDOFFS**

**GitHub source of truth:** `depulseapp/DE-PULSE`  
**Canonical machine state:** `governance/current-state.json`  
**Certified Stable:** `v18.9.1-stable`  
**Certified Stable candidate:** `e55d8d25b15cec2ffb0f5411bc358bc40b359cf9`  
**Certified Stable qualified source:** `d7276c3421dd2b4529ac2a987466be3cffa05678`  
**Certified Stable fingerprint:** `0062f46dea5690d0b3fcd8a9ed3b1f71ebe1522c7dee2cb218e9d36b9e0076ff`  
**Certified Stable build ID:** `v18.9.1-stable-20260821`  
**Immutable predecessor resume checkpoint release:** `v18.9.0` / `v18.9.0-stable`  
**Active work slice:** #73 / `ADAPT-ROOT-CONVERGENCE-001`  
**Active branch:** `adapt-root-convergence-001`  
**Closure ledger:** `governance/work-slices/ADAPT-ROOT-CONVERGENCE-001/closure.json`  
**Public product version consumed:** none  
**Product behavior change:** none intended.

## Resume Rule

1. Fetch live `main` and `adapt-root-convergence-001` first; another session may have advanced either.
2. Read `governance/current-state.json`, issue #73, this handoff, and the #73 work-slice/closure/root-disposition contracts.
3. Continue actual root convergence; do not restart #70 and do not begin TradeInsight Settings/API-key UX until #73 closes.
4. Preserve `v18.9.1-stable`, all published Stable assets/tags/evidence, the #70 GitHub-plan waiver, G0–G16, No Execution and canonical product/provider owners.

## Current root-convergence truth

#70 is complete and remains immutable history. It materially reduced root clutter and prevented new arbitrary root recurrence, but intentionally retained inherited v17/v18 compatibility classes with migrate-when-touched dispositions. #73 now owns the final repository-structure convergence requested after reviewing the live root.

Every tracked root file must receive exactly one `KEEP`, `MOVE`, `CONSOLIDATE` or `DELETE` disposition through `tools/ci/root_convergence_gate.py`. Historical/version-scoped non-Go evidence moves to governed release/history ownership. Obsolete copies may be deleted only after active-reference and assertion/evidence proof. Package-main Go source is not moved cosmetically; cohesive extraction must preserve tests, dependencies and private package access without artificial test-only production exports.

`WAIVER-GITHUB-MAIN-PROTECTION-001` remains retained under `governance/work-slices/ADAPT-CI-CONVERGENCE-001/main-protection-waiver.json`; actual `main` protection remains unavailable on the current GitHub plan. PR-first and exact-head Fast/Qualified compensating controls remain mandatory.

## Exactly one next action

Execute #73 Wave 1: relocate the historical v17/v18 non-Go root evidence to governed release/history owners, atomically rebind every active consumer, then run Fast on the exact resulting head before proceeding to root tooling and versioned-Go convergence.

Permanent product boundaries remain unchanged: US Equities Processing; No Execution; Smart Provider Router v2 sole routing owner; direct SEC/EDGAR authority for Form 4; GLD/SLV/USO actionable exceptions.
