# DE.PULSE — Current Adaptive CI / Versioning Convergence

**Program:** #70 / `ADAPT-CI-CONVERGENCE-001`  
**Status:** FINAL QUALIFICATION PENDING  
**Certified Stable:** `v18.9.1-stable`  
**Active branch / PR:** `adapt-ci-convergence-001` / Draft PR #71  
**Canonical machine state:** `governance/current-state.json`  
**Closure ledger:** `governance/work-slices/ADAPT-CI-CONVERGENCE-001/closure.json`

#70 implementation is materially complete. Repository convergence includes the three-workflow control plane, Planner v3, trustworthy base binding, targeted native rehearsals, canonical version-neutral G12, release identity/toolchain contracts, immutable publication checks, legacy CI retirement/history ownership, recursive source health, migration/equivalence gates, active version-named test-owner migration, cohesive production package decomposition, stable retained-asset ownership, prospective SemVer, G16 CI-efficiency evidence, and final-binary SPDX SBOM/provenance integration.

## GitHub main-protection external-control waiver

GitHub `main` remains factually unprotected because the configured private-repository ruleset cannot be enforced on the current organization plan without upgrading to GitHub Team. The owner declined that upgrade. This is not treated as a technical PASS.

The single gap `MAIN-PROTECTION-RULESET` is governed by `WAIVER-GITHUB-MAIN-PROTECTION-001` at `governance/work-slices/ADAPT-CI-CONVERGENCE-001/main-protection-waiver.json`. `tools/ci/work_slice_closure_gate.py` validates the waiver fail-closed, restricts it to this exact #70 control, requires explicit residual-risk acceptance, and requires the PR-first/exact-head Fast/exact-head Qualified/no-direct-main/no-force-push/G11-G16/provenance compensating controls. The waiver must be retired when technical enforcement becomes available or its revalidation triggers fire.

## Closure rule

Documentation alone still cannot close #70. Every non-waived blocking gap must be `VERIFIED`; the factual `BLOCKED_EXTERNAL` protection gap is closure-satisfying only while the machine-readable waiver above validates. Final exact-head `DE.PULSE/fast-head` and full `DE.PULSE/qualified-head` are still required on the current waiver/governance head before merge/closure.

The durable full acceptance contract remains issue #70 plus `adaptive-governance/ADAPTIVE_CI_OPERATING_CONTRACT.md`; this CURRENT file is the live projection from `governance/current-state.json`.
