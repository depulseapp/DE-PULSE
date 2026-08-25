# CURRENT Adaptive Build Process

**Canonical machine state:** `governance/current-state.json`  
**Certified Stable:** `v18.9.1-stable`  
**v18 state:** CLOSED by executable evidence.  
**Active process-only planning slice:** #110 / `ADAPT-V19-ZERO-MISS-PLAN-001`  
**Branch:** `adapt-v19-zero-miss-plan-001`  
**Future hosted umbrella:** #66 — not broadly started and no product implementation slice reserved.  
**Detailed plan:** `governance/V19_V20_ZERO_MISS_PLAN.md`  
**Requirement ledger:** `governance/v19-v20-requirement-conservation.json`

## Zero-miss execution process

Future implementation follows:

**LOOKUP → COMPARE → CLASSIFY → DECIDE → UPDATE → exact-head evidence → band reconciliation.**

At each G0/G1:
1. fetch live `main`, issue/PR state and current executable evidence;
2. locate the canonical owner for every requirement row targeted by the candidate slice;
3. classify overlap as `INHERITED`, `EXTEND_EXISTING_OWNER`, `REPLACE_OR_CONSOLIDATE`, or `NEW_RESIDUAL`;
4. freeze exactly one primary responsibility with explicit dependencies and `Not complete until` criteria;
5. attach the relevant conservation IDs to the work-slice contract and closure ledger;
6. implement all REQUIRED Mac/Windows/Web adapters for a shared capability in that same responsibility;
7. run exact-head Fast and impact-selected Qualified;
8. at band closure, reconcile every applicable row before advancing.

A later version cannot absorb an unfinished prerequisite simply to keep schedule/version numbering. A no-feature closure version cannot be used to hide feature work that should have been implemented in an owning slice.

`tools/ci/v19_v20_requirement_conservation_gate.py` is the machine enforcement inside existing G2/G10. It fails for missing mapped GitHub source coverage, illegal/unassigned states, missing targets, roadmap/plan version drift, missing closure checkpoints or a newly declared v18 corrective count.

## Retained Adaptive Data Health process

The inherited **#81/#82/#83/#78/#84** process remains in force: unclassified provider/fetch paths **fail closed**, **canonical freshness** remains authoritative, and **Smart Provider Router v2** remains the only general routing/admission owner. #95 registration recurrence and #94 observational usefulness remain extensions of those owners rather than parallel provider, health, cache, persistence, telemetry, reconciliation or lifecycle systems.

Direct SEC/EDGAR remains Form 4 authority. U.S. equities processing, GLD/SLV/USO actionable exceptions, No Execution, the canonical multi-feed subscription owner, existing identity/session/calendar owners and G0-G16 remain protected.

## Cross-platform execution

Shared capability is complete only when all G1 `REQUIRED` clients pass equivalent domain/API/state behavior. Mac-only, Windows-only or Web-only success is diagnostic and does not permit a later shared domain to begin while material parity debt remains.

## Current process action

Finish #110 projection/checkpoint convergence and exact-head CI. Once #110 closes, perform a fresh source-overlap audit for `v19.0.0`; if current rights/registration code already satisfies all mapped requirements, record inheritance rather than rewriting it. Reserve only the first genuine residual.
