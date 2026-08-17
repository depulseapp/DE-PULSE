# DE.PULSE v18.5.0 — Post-Tag G16 Governance Sync Plan

Status: PREPARED / EXECUTE ONLY AFTER `v18.5.0-stable` IS PUBLISHED AND VERIFIED

## Purpose
The immutable Stable tag must remain the exact product source that passed final Stable G11–G15. Process improvements developed while certification was runner-constrained must also survive into the next development baseline. These goals are compatible by applying G16 governance metadata to `main` **after** the immutable Stable tag is created, without moving or modifying that tag.

## Invariants
1. `v18.5.0-stable` remains permanently bound to the exact G15-certified Stable source commit.
2. GitHub Release assets remain bound to that same tagged source and G15 evidence.
3. Post-tag G16 sync MUST NOT change product/runtime/provider/persistence/intelligence semantics represented by the Stable tag.
4. G16 sync is runner-free repository governance work and does not require re-certifying the immutable Stable tag.
5. v18.5.1 branches from the resulting G16 `main` head, not from an older pre-governance Stable branch.

## G16 sync content
At minimum synchronize these permanent process artifacts from the v18.5 tooling line onto post-tag `main`:
- `governance/GITHUB_ACTIONS_EFFICIENCY_CONTRACT.md`
- final v18.5 G16 handoff / release provenance references under `release/v18.5.0/`
- the current workflow ownership/retirement manifest needed to explain historical v18.5 CI
- the v18.5.1 repository archetype build plan / repository structure contract already approved for the next build

Do **not** carry dormant v18.5 release trigger files forward as active v18.5.1 release automation. v18.5.1 creates its own version-scoped intent triggers only when its exact closure workflow is defined.

## Workflow cleanup after Stable
- v18.5 development auto workflows: retired.
- v18.5 legacy release-certification workflow: retired.
- v18.5 TEST-native duplicate workflow: retired.
- v18.5 Stable candidate/certification/publication workflows: preserve as historical release tooling on the tooling branch; do not make them auto-active on v18.5.1.

## v18.5.1 inheritance
v18.5.1 must inherit these rules from G16 `main`:
- GitHub is durable source/release archive, not default compute for every edit.
- expensive release workflows are intent-triggered only.
- reuse exact unaffected evidence when valid.
- cheap/static/focused checks precede expensive jobs.
- one fresh authoritative final-candidate G12.
- one normal macOS Apple Silicon + Windows x64 native pass on the final closure candidate, in parallel after shared prerequisites.
- no publication rebuild.
- no certification waiver for budget/cost.

## Execution record
When final Stable publication succeeds, replace this prepared plan with/alongside an executed G16 record containing:
- immutable Stable tag and source commit;
- Stable fingerprint/source ZIP SHA;
- Stable certification run ID;
- GitHub Release verification result;
- post-tag G16 `main` governance commit SHA;
- v18.5.1 starting branch/commit.

This plan does not add a top-level gate; it is execution detail inside permanent G16 Adaptive Retrospective / Handoff.
