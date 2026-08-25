# CURRENT Adaptive Build Process

**Canonical machine state:** `governance/current-state.json`  
**Certified Stable:** `v18.9.1-stable`  
**Completed provider onboarding:** #95 / PR #101 / merge `2eab9bd38b0a75a116de46e531015ed699ed7308`  
**Active process work:** #102 / `ADAPT-POST-STABLE-CONTINUITY-001` / `adapt-post-stable-continuity-001`  
**Canonical closure ledger:** `governance/work-slices/ADAPT-POST-STABLE-CONTINUITY-001/closure.json`  
**Separate product residual:** #94 / `ADAPT-PROVIDER-TELEMETRY-001` — not started  
**Parent program:** #65 / `ADAPT-PROVIDER-INTELLIGENCE-010`  
**Future hosted program:** #66 remains blocked/not started.

For #102, follow LOOKUP -> COMPARE -> CLASSIFY -> DECIDE -> UPDATE and remain evidence-first:

- fetch live `main` and the live #102 branch before each mutation;
- verify immutable `v18.9.1-stable`, PR #69, Fast #581, Qualified #163, Release #33 and retained artifacts before writing checkpoint evidence;
- never infer artifact IDs/digests, merge SHAs or gate conclusions from filenames or prose;
- reconcile stale lower-authority checkpoints to GitHub truth; never rewrite GitHub/release history to make a stale checkpoint pass;
- reuse the existing `DE.PULSE-STABLE-EVIDENCE-1` retrospective manifest pattern and explicitly preserve Stable immutability/no-rebuild;
- update #95 machine closure from real executable evidence only; do not reopen completed provider implementation;
- project #102 consistently through current-state, handoff, the four CURRENT Adaptive layers, CI convergence and gap closure;
- keep `blocksNextProductCapability=true` while #102 is active; #94 remains separate and #66 remains blocked;
- do not modify product/provider behavior, release identity, Router v2, Data Health, lifecycle, freshness, persistence, rights or security semantics;
- use one coherent process branch/PR; do not make commits merely to trigger CI;
- run canonical Fast on the exact final candidate, then deliberate impact-selected Qualified on the identical head;
- any source/governance mutation after qualification creates a new candidate and invalidates earlier #102 evidence;
- merge only after re-fetching live `main` and using the expected-head guard;
- no direct-main patch, no gate waiver, no Release workflow and no Stable/public SemVer publication for #102.

## Retained Data Health invariants

The completed **#80 -> #81/#82/#83/#78 -> #84 -> #95** provider architecture remains inherited executable authority. Registration-aware recurrence, canonical freshness, scoped degradation/recovery, Smart Provider Router v2, lifecycle/readiness, fault recovery, cache/persistence reuse, direct SEC/EDGAR Form 4 authority, U.S. equities, GLD/SLV/USO and No Execution remain unchanged.

## Acceptance discipline

Source/governance edits are not closure evidence by themselves. #102 remains fail-closed until its closure ledger is satisfied by executable continuity gates, exact-head Fast, impact-selected Qualified and guarded merge evidence.
