# DE.PULSE v18.5 — Active Workflow Manifest

Status: CURRENT v18.5 RELEASE TOOLING

This manifest exists to make hosted-runner ownership explicit and to prevent historical or automatic workflows from silently consuming GitHub Actions minutes.

## Active release workflow

| Purpose | Workflow | Trigger file | Runner policy |
|---|---|---|---|
| Complete v18.5 Stable closure from the authoritative TEST G12 evidence through G11–G16 | `.github/workflows/v18.5-stable-closure.yml` on `v18.5-stable-promotion-prep` | `.depulse-certification/triggers/v18.5-stable-closure.json` | **One intentional Actions run:** small Ubuntu G11 identity/source freeze → one authoritative Stable G12 → macOS Apple Silicon + Windows x64 in parallel → G15 exact-artifact assurance → no-rebuild tag/Release/G16 governance sync |

## Trigger invariant
The workflow is dormant unless `.depulse-certification/triggers/v18.5-stable-closure.json` is the changed path. The trigger must contain `run: true`, an incremented `nonce`, and the exact authoritative G12/RC/expected-Stable-fingerprint bindings. Ordinary source, documentation, governance, native-audit-tooling or workflow-file changes cannot intentionally start release compute.

## Single-run retry behavior
The G11 job is idempotent:
- if `v18.5.0-stable-candidate` still equals the frozen RC, it verifies and applies the reviewed 12-file Stable identity patch and creates exactly one identity-only commit;
- if that exact one-commit Stable candidate already exists from an earlier partial run, it verifies its parent, 12-file diff and exact fingerprint, then reuses it rather than creating another commit.

GitHub's failed-job rerun should be preferred for downstream transient failures so already-successful expensive jobs are not repeated unnecessarily.

## Retired v18.5 workflows
- `v18.5-development/.github/workflows/v18.5-dev-ci.yml` — retired after exact G0–G5 evidence passed.
- `v18.5-development/.github/workflows/v18.5-g10-prefreeze.yml` — retired after exact G6–G10 evidence passed.
- `v18.5-release-certification/.github/workflows/v18.5-release-certification.yml` — retired after authoritative TEST G11/G12 run `32009262146`; exact evidence/artifacts remain the audit trail.
- `v18.5-native-retry/.github/workflows/v18.5-native-retry.yml` — retired because a TEST native pass followed by a fresh Stable native pass duplicated expensive platform evidence.
- `.github/workflows/v18.5-stable-candidate.yml` — superseded by single-run Stable closure.
- `.github/workflows/v18.5-stable-certification.yml` — superseded by single-run Stable closure.
- `.github/workflows/v18.5-stable-publish.yml` — superseded by single-run Stable closure.

Retirement means the workflow file was removed after its useful evidence or design purpose was superseded. Git history and immutable source/evidence references remain available.

## Why there is only one native pass
The frozen TEST product semantics already passed authoritative G12. The Stable transition changes only the reviewed/allowlisted release identity/generated surfaces and is followed immediately inside the same closure run by a fresh exact Stable G11–G15 certification. Therefore the macOS Apple Silicon and Windows x64 G13/G14 audits run once—on the actual final Stable binaries that will be published. Stable G15 verifies those exact artifacts and hashes before publication.

This removes duplicate native compute without removing, waiving or weakening G13, G14 or G15.

## Canonical closure tooling
- G11 deterministic identity/source freeze: `certification/v18_5_stable_g11_prepare.sh`
- reviewed identity patch: `certification/v18_5_stable_identity.patch`
- identity patch provenance: `certification/v18_5_stable_identity_patch.json`
- G12 authoritative certification: `certification/v18_5_stable_g12.sh`
- macOS Apple Silicon G13/G14: `certification/v18_5_g14_macos_arm64_generic.sh`
- Windows x64 G13/G14: `certification/v18_5_g14_windows_x64_generic.ps1`
- G15 bundler: `certification/v18_5_stable_g15.py`
- G16 publication metadata verification: `certification/v18_5_stable_publish_prepare.py`
- no-rebuild publication/G16 sync: `certification/v18_5_stable_publish.sh`

## Cost / evidence rule
No release gate may be waived to save cost. Reuse exact unaffected evidence, run cheap prerequisites before expensive jobs, run macOS and Windows only once on the exact final closure candidate, and publish already-certified artifacts without rebuilding.

This manifest is governed by `governance/GITHUB_ACTIONS_EFFICIENCY_CONTRACT.md` and does not add a new G0–G16 gate.
