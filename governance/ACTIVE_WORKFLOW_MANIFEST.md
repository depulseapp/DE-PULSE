# DE.PULSE — Active Workflow Manifest

**Status:** PERMANENT / ACTIVE CI AUTHORITY  
**Scope:** operational workflow surface only; G0–G16 remains the permanent release model.

This manifest defines the only normal active GitHub Actions entry points for DE.PULSE. Historical/version-specific workflow files are retained by Git history and immutable Stable tags/releases, not by leaving obsolete orchestration active.

## Canonical active workflows

| Workflow | Responsibility | Normal trigger / invocation |
|---|---|---|
| `.github/workflows/ci-fast.yml` | Cheap development preflight, workflow/provenance policy, formatting/static/unit checks, renderer syntax, and conservative merged-branch hygiene after successful `main` pushes. | PRs; pushes to `main` and active `v*-development`; manual dispatch. |
| `.github/workflows/ci-qualified.yml` | Parameterized affected-area qualification, full backend/race/shuffle, renderer/browser contracts, and Linux/macOS/Windows CI-harness portability. | CI/release-tooling PR changes; manual/reusable invocation with exact SHA + lane. |
| `.github/workflows/release.yml` | Parameterized G11–G16: exact candidate, G12 full certification, independent macOS/Windows G13/G14, G15 assurance, optional no-rebuild Stable publication, G16 evidence. | Deliberate `workflow_dispatch` only. |

No other `.github/workflows/*.yml` or `.yaml` file is allowed by normal policy. `tools/ci/workflow_policy.py` enforces the allowlist.

## Version-neutral execution

Release version, candidate SHA, source fingerprint, build ID, certification script, gate/lane, target platform and publication intent are inputs/configuration. They must not create workflow families such as:

- `vX.Y-*.yml`;
- `*-retry.yml`;
- `*-monitor.yml`;
- `*-probe.yml`;
- `*-recovery.yml`;
- `*-certification.yml`;
- `*-publish.yml`.

A deliberate migration of the permanent CI architecture may temporarily require migration tooling, but its approval, scope, rollback and retirement criteria must be explicit before merge.

## Failure and retry invariant

Every failure is classified before rerun:

- `INFRA_FAIL` → rerun failed job(s) in the same workflow/SHA where possible;
- `CI_HARNESS_FAIL` → fix the canonical shared harness/tool, add prevention coverage, rerun the same affected workflow/lane;
- `GATE_TEST_FAIL` → fix the canonical test contract, rerun that gate and invalidated dependents;
- `PRODUCT_FAIL` → correct source, run the same workflow on the new SHA from the earliest invalidated gate;
- `EXPECTED_NOOP` → record and continue;
- `SUPERSEDED` → preserve history without recreating obsolete evidence.

Independent PASS evidence is retained when exact source/package/test/dependency fingerprints remain equivalent. A Windows failure does not automatically rerun macOS, and vice versa.

## Provenance invariant

Cross-platform source identity is computed from canonical raw Git object bytes using `source_fingerprint.py --mode git`. Filesystem-materialized fingerprints may be retained as diagnostics but are not authoritative release provenance because checkout behavior can differ by OS.

`ci-qualified.yml` proves the same Git-object fingerprint on Ubuntu, macOS Apple Silicon and Windows x64 and parses the reusable native harness on the corresponding native shell.

## Release invariant

`release.yml` follows:

**G11 immutable candidate → G12 full certification → G13/G14 macOS + Windows in parallel → G15 assurance → optional no-rebuild publication → G16 evidence.**

Publication downloads and uploads the exact already-certified artifacts. It does not rebuild native packages.

## Branch/repository hygiene

After successful `main` Fast CI, `tools/ci/branch_hygiene.py` may delete only remote branch tips already proven ancestors of `main`. Unique/diverged branches are never deleted automatically; they remain for explicit reconciliation.

Permanent branch model:

- `main`;
- one active release development branch when product work is active;
- genuinely short-lived feature/fix branches.

RC is an immutable SHA/checkpoint. Stable is an immutable tag + GitHub Release. Historical RC/retry/certification/promotion branches are not the archival model.

## Continuous-improvement invariant

G16 learning is incomplete when a recurring CI/process problem is merely documented. A material lesson must either:

1. harden the canonical workflow/tool/test with regression evidence; or
2. receive an explicit evidence-backed `NO_IMPLEMENTATION_CHANGE_REQUIRED` disposition.

The objective is a smaller, more reliable and more reusable CI system after each release—not an accumulating collection of special-case workflows.

Governed by:
- `adaptive-governance/ADAPTIVE_ROADMAP.md`;
- `adaptive-governance/ADAPTIVE_BUILD_PLAN.md`;
- `adaptive-governance/ADAPTIVE_BUILD_PROCESS.md`;
- `adaptive-governance/ADAPTIVE_DELIVERY_PROCESS.md`;
- `governance/GITHUB_ACTIONS_EFFICIENCY_CONTRACT.md`;
- `governance/REPOSITORY_STRUCTURE_CONTRACT.md`.
