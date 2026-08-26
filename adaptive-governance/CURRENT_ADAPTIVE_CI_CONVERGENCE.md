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

The original hosted-helper reachability defect remains resolved. Fast #1165 / run `32978721012` proved recursive/package-aware source health with zero unregistered orphan production Go helpers. Fast #1167 / run `32980268194` on `170477f222226d3ae7f5da41cb0822a655ab5e54` passed Canonical workflow policy, Recursive Source Health, Adaptive Data Health conservation, repository migration/current-state projection convergence, active work-slice closure validation, resume/post-Stable continuity, and v18 T1-T9 assurance including functional/integration and security/rights. It then failed only at V19 requirement conservation because `governance/ROADMAP.md` had lost the durable `Zero-Miss Future-Version Conservation` section identity while the underlying 72 HOST rows and conservation rules remained intact.

That roadmap conservation contract is restored: source-overlap classification remains required before G1 implementation and zero-gap closure remains required before advancing the owning dependency band, while coherent version/build planning remains authoritative.

The HOST-004..007 tenant/identity/device/session stage is now `IMPLEMENTED_UNVERIFIED` in the canonical closure ledger. HOST-008..023 and final exact-head Fast/Qualified remain OPEN.

Fast #1168 / run `32984624055` ended `startup_failure` before creating any job. GitHub reported zero jobs: no runner checkout and no DE.PULSE workflow step or gate executed. Therefore #1168 is external CI-startup evidence; it is neither a product PASS nor a DE.PULSE product/gate failure and creates no waiver.

No executable exact-head Fast has yet run on the branch after the HOST-004..007 closure-evidence update. Do not consume old-head or zero-job evidence as current qualification.

GitHub-hosted artifact attestation remains mandatory where supported. Where repository/account constraints prevent hosted attestation enforcement, existing exact-hash provenance, promotion verification and SBOM remain compensating evidence; this does not relax the exact-head Fast/Qualified requirement.

## Exactly one next action

Obtain executable exact-head Fast on PR #149 only when GitHub Actions successfully starts jobs. If it executes, continue from the first truthful failure or PASS. If GitHub returns another zero-job `startup_failure`, retain the external CI blocker without repeated retry churn, waiver, product-scope advance, Qualified/Release spend, or v19.1.0 start.
