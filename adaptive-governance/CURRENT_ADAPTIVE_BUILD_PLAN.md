# CURRENT Adaptive Build Plan

**Canonical machine state:** `governance/current-state.json`  
**Certified Stable:** `v18.9.1-stable`  
**v18 state:** CLOSED by executable evidence.  
**Active process-only planning slice:** #110 / `ADAPT-V19-ZERO-MISS-PLAN-001`  
**Future hosted umbrella:** #66 — no product implementation slice reserved.  
**Detailed version plan:** `governance/V19_V20_ZERO_MISS_PLAN.md`  
**Machine requirement ledger:** `governance/v19-v20-requirement-conservation.json`

## Zero-miss build rule

Every future work packet must carry parent conservation requirement IDs and one primary responsibility. Before G1 implementation, current source is classified `INHERITED`, `EXTEND_EXISTING_OWNER`, `REPLACE_OR_CONSOLIDATE`, or `NEW_RESIDUAL`. REUSE → CONSOLIDATE → REFACTOR → DELETE/REPLACE → ADD remains mandatory.

A packet cannot close because its headline works while a sub-contract is missing. Its `Not complete until` criteria in `governance/V19_V20_ZERO_MISS_PLAN.md`, all applicable conservation rows, required Mac/Windows/Web adapters and exact evidence must reconcile together. `UNASSIGNED`, `OPEN_WITHOUT_TARGET`, duplicate primary ownership or unexplained carry-forward blocks the band closure.

Each dependency band has a no-feature zero-gap checkpoint. That checkpoint verifies completeness; it does not become a dumping ground for forgotten implementation.

## Retained Adaptive Data Health build owners

The canonical inputs remain `governance/data-health/provider-capability-matrix.json`, `governance/data-health/data-health-slo.json`, and `governance/data-health/provider-fetch-paths.json`. They continue binding the **Adaptive Roadmap**, **Build Plan**, **Build Process** and **Delivery Process** to provider capability classification, canonical freshness/SLO, degradation/recovery and fetch-path ownership.

Completed #79/#84, #92, #95, #94 and #107 remain inherited executable foundations. Smart Provider Router v2, provider lifecycle/readiness, canonical freshness, subscription, persistence/cache, transport telemetry, semantic reconciliation and degradation owners cannot be replaced. Provider usefulness remains `ADVISORY_ONLY` unless a separately governed future adaptive slice is explicitly promoted.

## Future sequence control

The version sequence is authoritative in `governance/ROADMAP.md`; detailed responsibilities/dependencies/closure criteria are in `governance/V19_V20_ZERO_MISS_PLAN.md`. The machine ledger maps issue #66 body + all seven addenda, issue #65 strategic inheritance, issue #57 v18 closure and issue #110 planning requirements to those versions.

#110 is planning/governance only and temporarily blocks the next product capability. After #110 closes, re-fetch live `main`, perform source overlap for the earliest planned version, and reserve exactly one real residual. If `v19.0.0` is already fully satisfied by current `provider_data_rights.go` / `ProviderRegistration` evidence, close or skip it by evidence rather than duplicating equivalent code.
