# CURRENT Adaptive CI Convergence

**Certified Stable:** `v18.10.0` — immutable  
**Canonical rebaseline:** `governance/V19_V20_REBASELINE.md`  
**Machine current-state authority:** `governance/current-state.json`  
**Active version:** `v19.0.0`  
**Active work slice:** `ADAPT-HOSTED-TRUST-FOUNDATION-001`  
**Active issue/PR/branch:** #148 / PR #149 / `adapt-hosted-trust-foundation-001`  
**Canonical closure ledger:** `governance/work-slices/ADAPT-HOSTED-TRUST-FOUNDATION-001/closure.json`

Exactly three routine workflow families remain canonical: `ci-fast.yml`, `ci-qualified.yml`, and `release.yml`. G0-G16 remains the only gate model.

## Version-first CI rule

Requirements/backlog rows are evidence units, not CI/release units. Current CI strategy is:

`coherent changed candidate -> exact-head Fast -> impact/risk-selected Qualified -> G10 full-version reconciliation -> Release G11-G16 only for an immutable release candidate`

Do not create requirement-sized workflow runs, branches, certification branches, releases or public versions merely because a requirement row advanced. Do not merge unrelated product scope just to reduce Actions minutes.

Frozen v18 T1-T10 remains conserved baseline evidence. Unchanged historical assurance may be reused only through the established equivalence/fingerprint rules; changed/dependent behavior receives fresh evidence.

## Required planning inputs

Impact Planner and existing gates must consume/reflect, where applicable:
- `governance/V19_V20_REBASELINE.md`;
- `governance/programs/V19-V20-REBASELINE/backlog-version-matrix.json`;
- `governance/programs/V19-V20-REBASELINE/host-requirement-version-map.json`;
- #170 cross-integration/Market-Regime obligations;
- #171 UI/data-density/intelligence-maturity obligations.

These do not add G17. They strengthen requirement/owner/test selection and G10 reconciliation inside existing gates.

## Current executable state

The original hosted-helper reachability defect is resolved. On exact-head Fast #1165 for `f36417fda84e063d5a9cafcc31c464b051f5b3af`, recursive/package-aware source health reported zero unregistered orphan production Go helpers. Exact-head Fast #1166 for `08f0ffc64be8f05ad1dd6bb5155114d9cd60d3be` then passed Canonical workflow policy and Recursive source-health, including the restored Adaptive Data Health contract. Its current blocker is repository/current-state projection convergence, not product helper reachability.

GitHub-hosted artifact attestation remains mandatory where supported. Where repository/account constraints prevent hosted attestation enforcement, existing exact-hash provenance, promotion verification and SBOM remain compensating evidence.

## Exactly one next action

Restore the seven CURRENT/handoff projections to the canonical machine state and closure ledger, then obtain a fresh exact-head Fast on PR #149. Do not start v19.1.0 or spend Qualified/Release budget until v19.0.0 Fast evidence is green and current-version closure criteria are reconciled.
