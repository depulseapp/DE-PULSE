# CURRENT Adaptive CI Convergence

**Certified Stable:** `v18.10.0` — immutable  
**Canonical rebaseline:** `governance/V19_V20_REBASELINE.md`  
**Active version:** `v19.0.0`  
**Active issue/PR/branch:** #148 / PR #149 / `adapt-hosted-trust-foundation-001`

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

The pre-rebaseline product head `c5d0713d16f95522fd013123a78bc7cc58dc2422` failed Fast #1141 / `32929281393` at recursive source-health because hosted identity/session helpers were not production-referenced. Rebaseline-only governance commits do not make that product failure pass. A fresh exact-head Fast is required after the product reachability gap is corrected.

GitHub-hosted artifact attestation remains mandatory where supported. Where repository/account constraints prevent hosted attestation enforcement, existing exact-hash provenance, promotion verification and SBOM remain compensating evidence.

## Exactly one next action

Fix the v19.0.0 identity/session production-reachability/source-health gap on PR #149, then run exact-head Fast on the coherent candidate. Do not start v19.1.0 or spend Qualified/Release budget on an unready version head.
