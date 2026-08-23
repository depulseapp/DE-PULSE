# DE.PULSE — Current Handoff

**SUPERSEDES ALL PRIOR CHAT HANDOFFS**

**GitHub source of truth:** `depulseapp/DE-PULSE`  
**Canonical machine state:** `governance/current-state.json`  
**Certified Stable:** `v18.9.1-stable`  
**Certified Stable candidate:** `e55d8d25b15cec2ffb0f5411bc358bc40b359cf9`  
**Certified Stable qualified source:** `d7276c3421dd2b4529ac2a987466be3cffa05678`  
**Certified Stable fingerprint:** `0062f46dea5690d0b3fcd8a9ed3b1f71ebe1522c7dee2cb218e9d36b9e0076ff`  
**Certified Stable build ID:** `v18.9.1-stable-20260821`  
**Certified release run:** `32546555659`  
**Immutable predecessor resume checkpoint release:** `v18.9.0` / `v18.9.0-stable`  
**Active work:** #70 / `ADAPT-CI-CONVERGENCE-001`  
**Active branch / PR:** `adapt-ci-convergence-001` / Draft PR #71  
**Work-slice state:** FINAL_QUALIFICATION  
**Public product version consumed by #70:** none.

## Resume Rule

1. Fetch the CURRENT live head of `adapt-ci-convergence-001`; another session/process may have advanced it.
2. Read `governance/current-state.json`, issue #70 and current comments, `governance/work-slices/ADAPT-CI-CONVERGENCE-001/closure.json`, and the active external-control waiver below.
3. Inspect commits since the v18.9.1 Stable baseline so implemented work is never duplicated.
4. Continue actual executable closure from the exact current state; GitHub objects and CI evidence outrank chat memory.
5. Do not begin the next product capability while #70 remains open.

## GitHub main-protection truth and approved waiver

GitHub still truthfully reports `main` as unprotected. The repository owner configured the ruleset named `DE.PULSE main protection`, but GitHub states that rulesets are not enforced for this private organization repository on the current plan unless the organization upgrades to GitHub Team. The owner explicitly declined that upgrade on 2026-08-23.

This limitation is not relabeled as technical enforcement. It is governed by `WAIVER-GITHUB-MAIN-PROTECTION-001` at `governance/work-slices/ADAPT-CI-CONVERGENCE-001/main-protection-waiver.json` for the single closure gap `MAIN-PROTECTION-RULESET`.

The executable `tools/ci/work_slice_closure_gate.py` validates that waiver fail-closed and requires the compensating controls to remain mandatory: PR-first development, exact-head `DE.PULSE/fast-head`, exact-head `DE.PULSE/qualified-head`, no direct main push, no force push/deletion, canonical G11-G16 release, and exact-SHA/fingerprint provenance. If GitHub enforcement becomes available or repository ownership/visibility/maintainer conditions change, the waiver must be revalidated and actual protection becomes required again.

## Exactly one next action

Obtain final exact-head Fast + full Qualified evidence on the waiver/governance head. If both pass and the head does not move, PR #71 may proceed to merge/#70 closure under the approved external-control waiver. Do not begin TradeInsight Settings/API-key UX before #70 is truthfully closed.

## Current implementation checkpoint

All substantive #70 implementation gaps are already evidenced on the prior exact head. The only factual external control remains GitHub main protection, now handled through the bounded owner-approved waiver above. This source update itself changes the exact head, so final Fast + Qualified must be rerun before closure; previous qualification remains historical evidence and is not reused as exact-head certification for the newer governance head.

Permanent product boundaries remain unchanged: US Equities Processing; No Execution; Smart Provider Router v2 as sole routing owner; direct SEC/EDGAR authority for Form 4; GLD/SLV/USO actionable exceptions.
