# DE.PULSE — Current Adaptive Delivery Process

**Canonical machine state:** `governance/current-state.json`  
**Certified Stable:** `v18.9.1-stable`  
**Active work slice:** #70 / `ADAPT-CI-CONVERGENCE-001`  
**Active branch / PR:** `adapt-ci-convergence-001` / Draft PR #71.

## Current delivery truth

#70 is a process/release-engineering packet and consumes no public product version. Delivery remains exact-head Fast → exact-head Qualified → immutable release candidate → one canonical G11–G16 release lifecycle. Stable publication must never rebuild certified artifacts or overwrite differing Stable bytes. Prospective tag/version semantics come from `governance/versioning-policy.json`; historical Stable tags remain immutable.

Do not promote/merge this packet until the closure ledger is satisfied with executable evidence, including actual GitHub `main` protection/ruleset evidence.

Evergreen delivery rules remain in `adaptive-governance/ADAPTIVE_DELIVERY_PROCESS.md`.
