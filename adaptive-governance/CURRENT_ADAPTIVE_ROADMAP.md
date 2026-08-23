# DE.PULSE — Current Adaptive Roadmap

**Canonical machine state:** `governance/current-state.json`  
**Certified Stable:** `v18.9.1-stable`  
**Stable candidate:** `e55d8d25b15cec2ffb0f5411bc358bc40b359cf9`  
**Stable qualified source:** `d7276c3421dd2b4529ac2a987466be3cffa05678`  
**Stable release run:** `32546555659`  
**Active work slice:** #70 / `ADAPT-CI-CONVERGENCE-001`  
**Active branch / PR:** `adapt-ci-convergence-001` / Draft PR #71  
**Product behavior change:** none; process/release-engineering/repository convergence only.

## Current ordering

1. Finish #70 executable convergence and exact-head evidence.
2. Do not start/merge a new product-version capability while `productCapabilityGate.blocked=true` in `governance/current-state.json`.
3. After #70 closes, re-baseline the next product capability from GitHub durable scope; prospective public versioning follows `governance/versioning-policy.json`.

## Permanent product boundaries

US Equities Processing, No Execution, Smart Provider Router v2 sole routing ownership, direct SEC/EDGAR Form 4 authority, and GLD/SLV/USO actionable exceptions remain unchanged.

Evergreen roadmap detail remains in `adaptive-governance/ADAPTIVE_ROADMAP.md`; this CURRENT file is intentionally only the live projection and must not duplicate stale release identity.
