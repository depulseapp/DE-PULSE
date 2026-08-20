# DE.PULSE — Active Workflow Manifest

**Status:** PERMANENT / ACTIVE CI AUTHORITY  
**Scope:** operational workflow surface only; G0–G16 remains the permanent release model.

This manifest defines the only normal active GitHub Actions entry points. Historical/version-specific workflow files belong in Git history and immutable Stable evidence, not as active orchestration.

## Canonical active workflows

| Workflow | Responsibility | Normal trigger / invocation |
|---|---|---|
| `.github/workflows/ci-fast.yml` | Cheap development preflight, governance/provenance/release-state policy, affected Go/renderer checks, exact-head Fast status; on `main`, conservative branch hygiene plus a cheap post-Stable continuity sentinel only. | PR opened/synchronize/reopened; `main` push for hygiene + continuity sentinel; manual dispatch. |
| `.github/workflows/ci-qualified.yml` | Parameterized affected-area qualification, backend/race/shuffle, renderer/browser contracts and CI-harness portability. | PR `ready_for_review`; manual/reusable invocation with exact SHA + lane. |
| `.github/workflows/release.yml` | Single merged-release-PR G11–G16 path: immutable candidate, full certification, independent macOS/Windows package/runtime audit, G15 assurance, no-rebuild Stable publication and G16 workflow evidence. | A merged PR to `main` whose head is `v*-development` and whose changed paths include `release_identity.json` or `release.yml`. |

No other active `.github/workflows/*.yml` or `.yaml` file is allowed. `tools/ci/workflow_policy.py` enforces the allowlist.

## Version-neutral execution

Version, candidate SHA, source fingerprint, build ID, qualification lane and package target are configuration/evidence—not reasons to create workflow families such as `vX.Y-*.yml`, retry, monitor, recovery, certification or publish workflows.

## Exact-head development invariant

Fast is one PR event stream. Qualified starts only when the same PR is deliberately marked Ready. Release starts only after exact-head Fast + Qualified evidence exists and the release PR is merged. Failed evidence is classified before rerun; create no new branch/PR just to manufacture another event.

## Post-Stable continuity invariant

Release G16 produces workflow/release evidence but cannot silently update GitHub repository metadata after publication. Before the next product line begins, the actual Stable tag/Release must be reconciled into:
- both `.depulse-certification/resume/` checkpoints;
- `release/<version>/stable-evidence-manifest.json`;
- `handoff/CURRENT.md`;
- all four CURRENT Adaptive overlays.

`tools/ci/post_stable_continuity_gate.py` is the canonical cheap repository-level sentinel. A `main` checkout carrying a later STABLE identity than the durable Stable checkpoint is a continuity failure, not an in-flight candidate. The main-push sentinel detects this without rerunning expensive product qualification.

## Provenance invariant

Cross-platform source identity uses canonical Git object bytes via `source_fingerprint.py --mode git`. Publication reuses exact certified native artifacts and never rebuilds them.

## Branch/repository hygiene

Permanent branch model:
- `main`;
- one active `v<version>-development` branch during product work;
- genuinely short-lived bounded fix/continuity branches when needed.

RC is an immutable SHA/checkpoint. Stable is an immutable tag + GitHub Release. Historical retry/certification/promotion branches are not archival evidence.

## Continuous-improvement invariant

G16 learning is incomplete when a recurring release/process issue is only described. Material lessons must harden canonical workflow/tool/test behavior or receive an evidence-backed `NO_IMPLEMENTATION_CHANGE_REQUIRED` disposition.

Governed by the Adaptive Roadmap/Build Plan/Build Process/Delivery Process, `governance/GITHUB_ACTIONS_EFFICIENCY_CONTRACT.md`, `governance/CI-EFFICIENCY-CONTRACT.md` and `governance/REPOSITORY_STRUCTURE_CONTRACT.md`.
