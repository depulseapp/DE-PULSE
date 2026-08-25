# CURRENT Adaptive CI Convergence

**Canonical machine state:** `governance/current-state.json`  
**Certified Stable:** `v18.9.1-stable`  
**Completed provider onboarding:** #95 / PR #101 / merge `2eab9bd38b0a75a116de46e531015ed699ed7308`  
**Active process work:** #102 / `ADAPT-POST-STABLE-CONTINUITY-001` / `adapt-post-stable-continuity-001`  
**Canonical active closure ledger:** `governance/work-slices/ADAPT-POST-STABLE-CONTINUITY-001/closure.json`  
**Separate product residual:** #94 / `ADAPT-PROVIDER-TELEMETRY-001` — not started  
**Parent program:** #65 / `ADAPT-PROVIDER-INTELLIGENCE-010`  
**Future hosted program:** #66 remains blocked/not started.

#70/#73 remain the completed CI/repository control-plane foundation. #102 uses only the three canonical workflows: `.github/workflows/ci-fast.yml`, `.github/workflows/ci-qualified.yml`, and `.github/workflows/release.yml`; however the Release workflow is **not triggered** for #102 because release identity does not change and this process slice creates no Stable/public SemVer release.

The merge-push Fast #966 failure that opened #102 was isolated to the unchanged `post_stable_continuity_gate.py`: default-branch Stable identity was v18.9.1 while durable resume checkpoints still described v18.9.0. Branch hygiene passed and product Fast validation was skipped on `main` exactly as the CI-efficiency contract requires.

#102 repairs lower-authority continuity evidence rather than weakening the sentinel. The process candidate must preserve exact v18.9.1 PR #69 / Fast #581 / Qualified #163 / Release #33 evidence and retained G12/macOS/Windows/G15/G16 artifact IDs/digests.

The completed #80-#84 provider/Data Health architecture and #95 registration-aware recurrence remain inherited executable authority. Canonical `ProviderRegistration` identities and Router v2 members must still classify against the provider-capability matrix, including historical Router `SEC` -> durable `SEC EDGAR` identity. Smart Provider Router v2 remains the sole general routing/admission authority; direct SEC/EDGAR remains Form 4 authority.

#102 qualification is exact-head: canonical **CI Fast** must PASS first, then Planner v3 impact-selected **CI Qualified** must PASS on the identical SHA. No lane is assumed before the planner evaluates the complete merge-base delta. Any candidate mutation invalidates prior evidence.

After qualification, merge requires a fresh live-main expected-head check. The next `main` push must show branch hygiene and a green post-Stable continuity sentinel. No gate waiver, temporary workflow, retry branch, direct-main patch or Stable release is permitted.
