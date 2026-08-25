# CURRENT Adaptive Delivery Process

**Canonical machine state:** `governance/current-state.json`  
**Certified Stable:** `v18.9.1-stable`  
**Completed provider onboarding:** #95 / PR #101 / merge `2eab9bd38b0a75a116de46e531015ed699ed7308`  
**Active process work:** #102 / `ADAPT-POST-STABLE-CONTINUITY-001` / `adapt-post-stable-continuity-001`  
**Canonical closure ledger:** `governance/work-slices/ADAPT-POST-STABLE-CONTINUITY-001/closure.json`  
**Separate product residual:** #94 / `ADAPT-PROVIDER-TELEMETRY-001` — not started  
**Parent program:** #65 / `ADAPT-PROVIDER-INTELLIGENCE-010`  
**Future hosted program:** #66 remains blocked/not started.

Delivery is PR-first, exact-head governed and evidence-bound. #102 is a process/release-continuity correction with no product behavior change and no public release.

## Before qualification

The #102 branch must prove:
- `.depulse-certification/resume/build-checkpoint.json` and `release-evidence-checkpoint.json` truthfully describe immutable v18.9.1 Stable rather than v18.9.0;
- `release/v18.9.1/stable-evidence-manifest.json` binds actual GitHub run/artifact IDs and digests and explicitly cannot redefine/rebuild/republish Stable;
- #95 is projected complete using its real Fast #965 / Qualified #193 / PR #101 merge evidence;
- `governance/current-state.json`, `handoff/CURRENT.md`, Roadmap, Build Plan, Build Process, Delivery Process, CI Convergence and Gap Closure agree on #102 as the active process work;
- #94 remains separate/not started and #66 remains blocked/not started;
- the unchanged post-Stable continuity, portability, source-health/current-state and closure gates pass;
- no product source/provider/router/Data Health/lifecycle/freshness/cache/persistence behavior changed.

## Retained Data Health delivery contract

The completed #80-#84 provider/Data Health foundation remains delivery authority under completed #95 and active #102. The exact inherited recurrence contract remains: **canonical Fast exact-head PASS** must precede **Qualified exact-head PASS** on the identical candidate. Smart Provider Router v2 remains the sole general routing/admission authority, direct **SEC/EDGAR** remains Form 4 authority, canonical freshness/degradation semantics remain unchanged, and **No Execution** remains permanent.

## Canonical qualification order

1. Re-fetch live `main` and live `adapt-post-stable-continuity-001`; inspect the complete delta.
2. Commit one coherent continuity candidate before opening the PR whenever practical.
3. Open one Draft PR for #102.
4. Run canonical Fast exact-head PASS on the PR head.
5. Mark the same exact candidate Ready so Planner v3 runs impact-selected Qualified exact-head PASS.
6. Any candidate mutation after Fast/Qualified invalidates the corresponding evidence and must be requalified.
7. Re-fetch live `main`; merge only with expected-head protection and only if the qualified PR head is unchanged.
8. Inspect the resulting main-push branch-hygiene/post-Stable continuity run. The continuity sentinel must be green without weakening it.
9. Record immutable run/merge evidence on issue #102. No Release workflow or Stable/public SemVer release is run for this process slice.

## Closure boundary

#102 closes only when the durable checkpoints, retrospective Stable manifest, #95 post-merge projections, all CURRENT surfaces, exact-head Fast/Qualified and guarded merge/main-sentinel evidence reconcile. Product work does not start until a fresh #65 overlap audit selects the next genuine residual.
