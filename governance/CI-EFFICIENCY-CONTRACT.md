# DE.PULSE CI Efficiency & Release Flow Contract

## Purpose

Keep DE.PULSE release quality at G0–G16 while preventing branch, pull-request, workflow-run, runner-minute and artifact amplification. GitHub remains the durable source of truth; CI exists to verify product quality, not to manufacture events.

## Canonical workflow set

Exactly three active workflow files are permitted:

1. `.github/workflows/ci-fast.yml`
2. `.github/workflows/ci-qualified.yml`
3. `.github/workflows/release.yml`

Do not create version-specific, retry, monitor, probe, recovery, certification or publication workflow files.

## Branch and PR model

For a normal release line use one long-lived development branch, for example `v18.6.1-development`, and one PR from that branch to `main`.

Short-lived `feature/*`, `fix/*` or `agent/*` branches are allowed only for genuinely independent code/process changes. They are not allowed merely to trigger Actions or retry CI.

The following versioned branch purposes are prohibited going forward: `*-release-certification`, `*-stable-promotion`, `*-certification-trigger`, `*-cert-trigger`, `*-promotion-trigger`, `*-dispatch*`, `*-retry*`, and `*-fallback*`.

A CI retry never creates a new branch or PR.

## CI Fast

CI Fast runs on PR `opened`, `synchronize` and `reopened`, plus `main` pushes for post-merge sanity and branch hygiene. It does not run on PR close and does not run independently on development-branch pushes, preventing the same commit from receiving both push and PR runs.

Fast owns cheap, impact-aware validation: workflow policy, portability/resume contracts, source provenance, release identity, formatting/vet/unit tests when Go changes, renderer syntax and focused regressions when renderer code changes.

CI Fast must never dispatch Release and must not upload an impact-plan artifact on every successful run. Its concurrency group cancels obsolete runs for the same PR/ref.

## CI Qualified / G10

A release PR remains Draft during active development. Marking the candidate `Ready for review` is the normal automatic G10 trigger.

Qualified runs once for the candidate SHA. Product candidates run full backend, renderer and browser qualification on Linux. macOS/Windows portability runners are reserved for `ci-harness` changes because final native macOS/Windows runtime truth is already mandatory in G13/G14.

If the candidate changes after a product/test failure, return the PR to Draft, apply fixes on the same branch, then mark it Ready again. If the SHA is unchanged and the failure is infrastructure-only, rerun only the failed job(s). Do not create a retry branch.

## Release / G11–G16

Release has one routine trigger: a release PR is merged and that PR changes `release_identity.json`.

The merged commit is the immutable candidate. One Release workflow performs G11 provenance, G12 full certification, G13/G14 macOS Apple Silicon package/runtime audit, G13/G14 Windows x64 package/runtime audit, G15 release assurance, publication of those exact same-run certified artifacts without rebuild, and G16 closure evidence.

There is no release-certification branch, no Stable-promotion branch, no trigger PR and no second promotion workflow run.

If Release fails, Stable remains the previous immutable version. Fix source in the normal development branch/PR flow and produce a new candidate; infrastructure failures on the same SHA use GitHub's failed-job rerun capability.

## Retry classification

- `INFRA_FAIL`: rerun only failed job(s) on the same SHA.
- `CI_HARNESS_FAIL`: fix the canonical harness on the same branch; requalify the affected candidate.
- `GATE_TEST_FAIL`: fix the canonical test or product defect on the same branch; rerun from the earliest invalidated gate.
- `PRODUCT_FAIL`: fix product code on the same development branch and same release PR.
- Never use a new branch, PR or metadata-only commit merely to create an Actions event.

## Runner and artifact budget

Use Linux for routine Fast and product Qualified lanes. Use macOS/Windows in Qualified only for CI/harness portability changes and in Release for authoritative native G13/G14 package/runtime audits.

Keep diagnostic artifacts only where they materially help investigation. Final release/native/G15/G16 evidence may use longer retention; routine impact plans and successful Fast diagnostics should remain in the job summary rather than becoming persistent artifacts.

## Branch hygiene

After a successful `main` Fast run, delete branches fully contained in `main`. Also remove deprecated versioned trigger/retry/dispatch/release-certification/stable-promotion branches when they have no open PR. Never delete a branch that still has an open PR.

## Non-negotiable quality boundaries

This efficiency contract does not weaken G0–G16, source fingerprinting, deterministic Day/Swing/Long behavior, Smart Provider Router ownership, U.S. Equities Processing Boundary, permanent No Execution Boundary, native macOS Apple Silicon delivery, Windows x64 delivery, evidence integrity or Stable immutability.
