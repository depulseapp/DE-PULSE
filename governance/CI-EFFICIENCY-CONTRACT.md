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

### Coherent-change batching / event minimization

For automation/connector-driven maintenance and for one coherent correction packet, prepare all related file mutations first and advance the development branch **once** with one Git tree/commit whenever practical. Do not use per-file GitHub writes that each move the PR head and therefore generate a separate `pull_request.synchronize` Fast run.

A fresh Fast run is expected for every genuinely new candidate SHA; this is a quality requirement and must not be disabled. The efficiency rule is therefore to minimize unnecessary candidate SHAs, not to suppress validation of a changed candidate.

Multiple commits while a PR is open are acceptable only when they represent genuinely independent review units or when a real defect discovered by CI requires a corrected candidate. Metadata-only or handoff-only commits must not be created merely to manufacture a CI event.

Normal target for a coherent packet: `one development branch → one Draft PR → one batched candidate update → one Fast → Ready → one Qualified → merge`; additional Fast/Qualified runs are justified only by an actual candidate change or an explicitly classified same-SHA infrastructure retry.

## CI Fast

The CI Fast workflow listens to PR `opened`, `synchronize` and `reopened`, plus `main` pushes so branch hygiene can run. The **Fast validation job is skipped on `main` pushes**; a merge therefore does not repeat the Fast tests already proven on the PR head. It does not run on PR close and does not run independently on development-branch pushes.

For PR events, Fast explicitly checks out and validates `pull_request.head.sha`, not GitHub's synthetic PR merge SHA. A successful run records the durable commit status `DE.PULSE/fast-head` on that exact source commit.

Fast owns cheap, impact-aware validation: workflow policy, portability/resume contracts, source provenance, release identity, formatting/vet/unit tests when Go changes, renderer syntax and focused regressions when renderer code changes.

CI Fast must never dispatch Release and must not upload an impact-plan artifact on every successful run. Its concurrency group cancels obsolete runs for the same PR/ref.

On a `main` push only the branch-hygiene job is intended to execute; that job is independent of Fast-test success and removes obsolete branches conservatively.

## CI Qualified / G10

A release PR remains Draft during active development. Marking the candidate `Ready for review` is the normal automatic G10 trigger.

Qualified explicitly checks out the exact PR head. Product candidates run full backend, renderer and browser qualification on Linux. macOS/Windows portability runners are reserved for `ci-harness` changes because final native macOS/Windows runtime truth is already mandatory in G13/G14.

A successful Qualified evidence job records the durable commit status `DE.PULSE/qualified-head` on the exact candidate SHA. These exact-SHA Fast and Qualified statuses are the release handoff evidence; they avoid relying on Actions run metadata that may refer to GitHub's synthetic PR merge commit.

If the candidate changes after a product/test failure, return the PR to Draft, apply fixes on the same branch, then mark it Ready again. If the SHA is unchanged and the failure is infrastructure-only, rerun only the failed job(s). Do not create a retry branch.

## Release / G11–G16

Release has one routine trigger: a release PR is merged, that PR changes `release_identity.json`, the base branch is `main`, and the source branch follows `v*-development`.

Before G12 can begin, G11 must fail closed unless:
- the exact release PR head carries successful `DE.PULSE/fast-head` and `DE.PULSE/qualified-head` commit statuses;
- the PR head SHA still matches GitHub's durable pull-request source ref;
- the canonical source fingerprint of that qualified PR head equals the merged candidate fingerprint;
- release identity and the version-specific certification script resolve correctly.

The merged commit is then the immutable candidate. One Release workflow performs G11 provenance, G12 full certification, G13/G14 macOS Apple Silicon package/runtime audit, G13/G14 Windows x64 package/runtime audit, G15 release assurance, publication of those exact same-run certified artifacts without rebuild, and G16 closure evidence.

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

On a `main` push, run branch hygiene without rerunning Fast tests. Never delete a branch with an open PR. Delete:
- merged PR head branches, including squash-merged heads that are not Git ancestors of `main`;
- branches fully contained in `main`;
- versioned branches whose corresponding Stable line is already published;
- closed/orphaned trigger, retry, dispatch, certification, promotion and similar release-temp branches.

If GitHub PR state cannot be resolved, fail conservative and retain uncertain unique branches.

## Non-negotiable quality boundaries

This efficiency contract does not weaken G0–G16, source fingerprinting, deterministic Day/Swing/Long behavior, Smart Provider Router ownership, U.S. Equities Processing Boundary, permanent No Execution Boundary, native macOS Apple Silicon delivery, Windows x64 delivery, evidence integrity or Stable immutability.
