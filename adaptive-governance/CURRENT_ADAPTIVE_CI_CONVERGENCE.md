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

The original hosted-helper reachability defect remains resolved. Exact-head Fast #1167 / run `32980268194` on `170477f222226d3ae7f5da41cb0822a655ab5e54` passed Canonical workflow policy, Recursive source-health, Adaptive Data Health conservation, repository migration/current-state projection convergence, active work-slice closure-ledger validation, resume/post-Stable continuity, and v18 T1-T9 assurance. It then failed only at V19 requirement conservation because `governance/ROADMAP.md` had lost the durable `Zero-Miss Future-Version Conservation` section identity while the underlying 72 HOST rows and conservation rules remained intact.

Current branch head `716bfa1495de18061484ad26f81f1e830c520378` restores that canonical zero-miss section together with source-overlap-before-G1 and zero-gap-before-next-band semantics while retaining coherent version/build planning. No gate, Stable artifact or packet-sized release model was reintroduced.

GitHub-hosted artifact attestation remains mandatory where supported. Where repository/account constraints prevent hosted attestation enforcement, existing exact-hash provenance, promotion verification and SBOM remain compensating evidence.

## Exactly one next action

Obtain fresh exact-head Fast on the current PR #149 head, verify V19 requirement conservation and the remaining formatting/vet/unit/static gates, then continue fixing only the first truthful executable failure. Do not start v19.1.0 or spend Qualified/Release budget until v19.0.0 Fast evidence is green and current-version closure criteria are reconciled.
