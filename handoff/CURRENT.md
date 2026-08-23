# DE.PULSE — Current Handoff

**SUPERSEDES ALL PRIOR CHAT HANDOFFS**

**GitHub source of truth:** `depulseapp/DE-PULSE`  
**Canonical machine state:** `governance/current-state.json`  
**Canonical #70 closure ledger:** `governance/work-slices/ADAPT-CI-CONVERGENCE-001/closure.json`  
**Certified Stable:** `v18.9.1-stable`  
**Certified Stable candidate:** `e55d8d25b15cec2ffb0f5411bc358bc40b359cf9`  
**Certified Stable qualified source:** `d7276c3421dd2b4529ac2a987466be3cffa05678`  
**Certified Stable fingerprint:** `0062f46dea5690d0b3fcd8a9ed3b1f71ebe1522c7dee2cb218e9d36b9e0076ff`  
**Certified Stable build ID:** `v18.9.1-stable-20260821`  
**Certified release run:** `32546555659`  
**Immutable predecessor resume checkpoint release:** `v18.9.0` / `v18.9.0-stable`  
**Completed work:** #70 / `ADAPT-CI-CONVERGENCE-001`  
**Implementation branch:** `adapt-ci-convergence-001`  
**Closure branch:** `adapt-ci-convergence-001-closure`  
**Merged implementation PR / commit:** #71 / `6ba8c0c2b486bdbbebac4611f440741d0588c65f`  
**Work-slice state:** COMPLETE  
**Next reserved capability:** TradeInsight Settings/API-key UX.

## Resume Rule

1. Fetch live `main` first; concurrent sessions may advance it.
2. Read `governance/current-state.json` and this handoff before starting new work.
3. #70 is complete only with its machine closure ledger, final qualification binding and bounded external-control waiver retained.
4. Start the next product capability as a **new work slice/branch/PR** from current `main`; do not reopen or extend #70 for product scope.

## #70 closure truth

The exact implementation candidate `ffeb1640174a744cf85578e26bbed7abd828cee1` earned Fast #742 / `32656818767` and Qualified #175 / `32656912135`, both PASS. Qualified covered CI/harness, backend full Go, race detector, randomized package order, persistence/DB, security/data-rights, renderer, Chrome, WebKit, actual packaged macOS lifecycle rehearsal and actual packaged Windows runtime rehearsal. PR #71 merged with expected-head protection as `6ba8c0c2b486bdbbebac4611f440741d0588c65f`.

Post-run immutable evidence is bound at `governance/work-slices/ADAPT-CI-CONVERGENCE-001/final-qualification-evidence.json`. This avoids falsely editing the pre-run `FINAL-QUALIFIED` ledger entry after qualification and thereby creating an endless requalification loop.

## GitHub main-protection truth and approved waiver

`main` remains factually unprotected on the current GitHub plan. The configured `DE.PULSE main protection` ruleset is not technically enforced for this private organization repository without GitHub Team, and the owner declined that upgrade. This remains `MAIN-PROTECTION-RULESET = BLOCKED_EXTERNAL`, not a technical PASS.

The only accepted exception is `WAIVER-GITHUB-MAIN-PROTECTION-001` at `governance/work-slices/ADAPT-CI-CONVERGENCE-001/main-protection-waiver.json`. `tools/ci/work_slice_closure_gate.py` validates it fail-closed and preserves PR-first development, exact-head Fast/Qualified, no direct-main/force-push/deletion policy, canonical G11-G16 and exact-SHA/fingerprint provenance. The waiver must be revalidated if plan, repository visibility/ownership, platform capability, or maintainer population changes.

## Exactly one next action

Create the next governed product work slice for **TradeInsight Settings/API-key UX** from current `main`, re-checking the master provider program and dependencies before implementation. Do not assign a public product version merely because a new work slice starts; #70 established work-slice identity separately from public SemVer release grouping.

Permanent product boundaries remain unchanged: US Equities Processing; No Execution; Smart Provider Router v2 sole routing owner; direct SEC/EDGAR authority for Form 4; GLD/SLV/USO actionable exceptions.
