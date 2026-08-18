# DE.PULSE — Current Authoritative Handoff

**SUPERSEDES ALL PRIOR CHAT HANDOFFS.** GitHub is authoritative; chat history and temporary workspaces are advisory only.

**Release:** `v18.5.2`  
**Active branch:** `v18.6-development`  
**Stable status:** v18.5.2 STABLE / G0–G16 CLOSED  
**Repository:** `depulseapp/DE-PULSE`  
**Stable tag:** `v18.5.2-stable`  
**Stable promotion commit:** `d30e54db4908ca57c52ae298cc4ada3416fab46b`  
**Certified source checkout:** `e9ca615cf7ab97ac476128c81ee9ae2d7340c0d9`  
**Certified source fingerprint:** `807de082d43e83d1d3548bca9350d13b72ef4dc71a848940b73b63e4b4d215b0`  
**Runtime build ID:** `v18.5.2-stable-hotfix-20260817`  
**Last updated:** 2026-08-18 America/Vancouver

## Resume rule

Read `AGENTS.md` or `CLAUDE.md`, then reconcile this handoff, `release_identity.json`, both `.depulse-certification/resume/` checkpoints, actual branch/PR/check state and the immutable Stable tag/release. `v18.5.2-stable` remains the certified product authority; later process/governance/handoff commits do not redefine that release.

## v18.5.2 certified release

v18.5.2 remains fully closed and unchanged. Deterministic Day/Swing/Long formulas, the U.S. Equities Processing Boundary, the permanent No Execution Boundary and deterministic/statistical Market Mode ownership remain protected. TradeInsight remains `NOT_IMPLEMENTED` with `NONE` production influence until its governed implementation and validation lifecycle is complete.

Certified native packages remain:

- macOS Apple Silicon SHA-256: `91f7d64a433474c4efbed0bd5c7d065508b2ce6eeae0544b1c217849fe62a4ae`
- Windows x64 SHA-256: `6431562a67bcd55db6ebab2e6a09724006119910ea8adbc358afb8644a752326`

## CI & Repository Hygiene Consolidation — completed

The behavior-neutral consolidation required before normal v18.6 product work was merged in PR #14 at `c6cf481a32a8151306e930574d2b8465cb1fd094`.

Normal active workflows are exactly:

1. `.github/workflows/ci-fast.yml`
2. `.github/workflows/ci-qualified.yml`
3. `.github/workflows/release.yml`

The consolidation includes impact-selected lanes, workflow allowlisting, canonical Git-object cross-platform provenance, reusable native release harnesses, preserved independent platform PASS, same-workflow retry/resume, no-rebuild publication and conservative branch hygiene.

PR #14 final head `1ee74e1902860000e165d51bf349e7918dba851f` passed:

- CI Qualified run `32106754567` — SUCCESS
- CI Fast run `32106754569` — SUCCESS

The current remote inventory contains 16 branches including `main`. Automatic hygiene removes only tips already contained in `main`; unique/diverged tips remain for explicit reconciliation. `v18.6-development` was identical to post-consolidation `main` before the continuity repair below.

## Resume-contract compatibility repair — active

Fresh-session reconciliation found a portability mismatch between `adaptive_resume_gate.py` and the certified-Stable checkpoint schema. The gate expected development-style top-level fields while the Stable checkpoint keeps certified identity inside `certifiedStable` and Stable/post-release identity rules.

Current repair commits on `v18.6-development`:

- `4748a4e2a4f90bc08d468bf5cfdeca8257c5dc3b` — support both active-candidate and certified-Stable checkpoint schemas while preserving identity validation;
- `d585f6be590ca16ddf48dd6fe656a8047ea96139` — run the adaptive resume portability gate inside canonical CI Fast.

This is behavior-neutral process/continuity hardening and does not modify v18.5.2 Stable product behavior.

## Residuals

- Role-specific tab/navigation composition remains future work.
- TradeInsight implementation and governed promotion remain future work with `NONE` production influence.
- Fifteen non-main branches remain for explicit semantic reconciliation.
- Normal v18.6 product scope is not frozen yet; G1 must be evidence-selected from the conserved recovery ledger.

## Exactly one next action

**Finish and merge the resume-contract reconciliation on `v18.6-development`: update the machine build checkpoint to this active branch/process state, obtain canonical CI Fast PASS on the repaired portability gate, merge the repair to `main`, and only then begin v18.6 G0–G3 evidence-selected scope freeze.**

## Provider-neutral continuation instruction

> Connect to `depulseapp/DE-PULSE`, read `AGENTS.md` or `CLAUDE.md`, verify `v18.5.2-stable` at promotion commit `d30e54db4908ca57c52ae298cc4ada3416fab46b`, inspect `v18.6-development`, actual PR/check state, `handoff/CURRENT.md` and both resume checkpoints, then resume the single next action above. Preserve G0–G16, GitHub account/assistant independence, the U.S. Equities Processing Boundary, No Execution Boundary, deterministic Market Mode ownership, TradeInsight `NOT_IMPLEMENTED/NONE` truth and the certified v18.5.2 package hashes.
