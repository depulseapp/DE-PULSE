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

For automation/connector-driven maintenance and one coherent correction packet, prepare related file mutations on the development branch **before opening the PR whenever practical**. If a PR is already open, advance its head once per coherent correction rather than performing a sequence of incomplete per-file commits that each generate `pull_request.synchronize` Fast runs.

A fresh Fast run is expected for every genuinely new candidate SHA; this is a quality requirement and must not be disabled. Efficiency means minimizing unnecessary candidate SHAs, not suppressing validation of a changed candidate.

Multiple commits while a PR is open are acceptable only when they are genuinely independent review units or when CI discovers a real defect requiring a corrected candidate. Metadata-only or handoff-only commits must not be created merely to manufacture a CI event.

Normal target: `one development branch → batch coherent work → one Draft PR → one Fast → Ready → one Qualified → merge`; additional Fast/Qualified runs require an actual candidate change or an explicitly classified same-SHA infrastructure retry.

## CI Fast

CI Fast listens to PR `opened`, `synchronize` and `reopened`, plus `main` pushes so branch hygiene can run. The **Fast validation job is skipped on `main` pushes**; a merge therefore does not repeat the Fast tests already proven on the PR head. It does not run on PR close and does not run independently on development-branch pushes.

For PR events, Fast validates `pull_request.head.sha`, not GitHub's synthetic PR merge SHA. A successful run records `DE.PULSE/fast-head` on that exact source commit.

Fast owns cheap, impact-aware validation: workflow policy, portability/resume contracts, source provenance, release identity, formatting/vet/unit tests when Go changes, renderer syntax and focused regressions when renderer code changes.

CI Fast never dispatches Release and does not upload an impact-plan artifact on every successful run. Its concurrency group cancels obsolete runs for the same PR/ref.

On a `main` push only branch hygiene is intended to execute.

## CI Qualified / G10

A release PR remains Draft during active development. Marking the exact candidate `Ready for review` is the normal automatic G10 trigger.

Qualified checks out the exact PR head. Process-only CI/harness candidates use Ubuntu/macOS/Windows portability; product candidates run affected/full backend, renderer and browser qualification according to Impact Planner.

**Chrome and WebKit are the two primary browser engines.** Chrome owns broad behavior coverage. WebKit owns co-primary compatibility evidence for core renderer/UI interactions. Full/browser candidates require both; renderer/UI or WebKit-harness risk requires WebKit; unaffected backend/provider/process-only work does not pay browser runtime.

A successful Qualified evidence job records `DE.PULSE/qualified-head` on the exact candidate SHA.

If the candidate changes after a product/test failure, return the PR to Draft, fix the same branch, obtain new Fast, then mark the same PR Ready again. If the SHA is unchanged and the failure is infrastructure-only, rerun only failed work when there is a reasonable recovery signal. Do not create a retry branch and do not loop same-SHA retries indefinitely.

## Qualified telemetry contract

Each Qualified run retains a compact JSON telemetry artifact for 30 days and renders the same core information in the job summary. Telemetry includes:

- exact candidate SHA and selected lane;
- per-job queue seconds and execution seconds;
- Linux/macOS/Windows/unknown runner consumption;
- Chrome/WebKit dependency setup duration and `setup-python` pip cache-hit state when those lanes run;
- current PR-branch Fast/Qualified/Release run counts;
- conservative workflow-amplification warnings.

Telemetry is a cost/efficiency **proxy**, not a substitute for GitHub billing. It must not invent currency rates; `actualCurrencyCost` remains null in CI telemetry and financial billing remains external source-of-truth data.

Amplification thresholds are diagnostic warnings, not quality bypasses. Legitimate defect corrections may exceed the normal 1 Fast / 1 Qualified target.

## Durable Stable evidence

Transient CI logs and artifacts may expire, so each published Stable should have a compact repo-durable evidence manifest bound to the authoritative release-evidence checkpoint. It indexes:

- Stable tag, certified candidate, qualification source SHA and fingerprint;
- canonical Fast, Qualified, full-product Qualified and G11–G16 run IDs;
- required gate states;
- macOS, Windows, G15 and G16 artifact IDs/digests.

A retrospective manifest committed after a Stable publication must explicitly state that it does **not** redefine the immutable Stable tag or binaries. Impact Planner may treat only `release/v*/stable-evidence-manifest.json` as process-only retrospective evidence; executable release scripts remain full-qualification scope.

## Workflow linting and reproducibility

DE.PULSE uses two complementary controls:

1. zero-network structural lint for canonical workflow set, top-level/job mappings, indentation/whitespace and GitHub-expression balance;
2. DE.PULSE semantic workflow policy for triggers, exact-head behavior, browser policy, permissions, telemetry, release topology and prohibited amplification patterns.

Fast/Qualified external Actions and browser dependencies remain immutable/pinned. Release workflow Action pins are updated with the next genuine release-capable product slice so their G11–G16/native proof is useful rather than a process-only paid rerun.

## Release / G11–G16

Release has one routine trigger: a merged release PR changes `release_identity.json` or the canonical Release workflow, targets `main`, and originates from `v*-development`.

Before G12, G11 fails closed unless:
- the exact release PR head has successful `DE.PULSE/fast-head` and `DE.PULSE/qualified-head` statuses;
- the durable PR source ref still matches that source head;
- the canonical source fingerprint of the qualified head equals the merged candidate fingerprint;
- release identity and version-specific certification script resolve correctly.

One Release workflow performs G11 provenance, G12 full certification, G13/G14 macOS Apple Silicon and Windows x64 package/runtime audit, G15 assurance, publication of exact same-run certified artifacts without rebuild, and G16 closure.

There is no release-certification branch, Stable-promotion branch, trigger PR or second promotion workflow run.

If Release fails, Stable remains the previous immutable version. Source/harness defects use normal same-branch/PR correction; same-SHA infrastructure failures use bounded failed-job reruns.

## Retry classification

- `INFRA_FAIL`: rerun only failed unchanged-SHA work when a recovery signal exists.
- `CI_HARNESS_FAIL`: fix canonical harness on same branch; requalify affected candidate.
- `GATE_TEST_FAIL`: fix canonical test or product defect on same branch; rerun from earliest invalidated gate.
- `PRODUCT_FAIL`: fix product code on the same development branch and PR.
- Never use a new branch, PR or metadata-only commit merely to create an Actions event.

## Runner and artifact budget

Use Linux for routine Fast and broad deterministic/product qualification. Use macOS/Windows in Qualified only when the impact plan requires portability/WebKit, and in Release for authoritative native G13/G14 proof.

Keep diagnostic artifacts only where they materially help. Qualified telemetry is intentionally tiny and 30-day. Stable release/native/G15/G16 truth is indexed durably; successful Fast impact plans remain in summaries rather than persistent artifacts.

## Branch hygiene

On a `main` push, run branch hygiene without rerunning Fast tests. Never delete a branch with an open PR. Delete:
- merged PR head branches, including squash-merged heads that are not Git ancestors of `main`;
- branches fully contained in `main`;
- versioned branches whose corresponding Stable line is already published;
- closed/orphaned trigger, retry, dispatch, certification, promotion and similar release-temp branches.

If GitHub PR state cannot be resolved, fail conservative and retain uncertain unique branches.

## Non-negotiable quality boundaries

This efficiency contract does not weaken G0–G16, source fingerprinting, deterministic Day/Swing/Long behavior, Chrome + WebKit primary evidence when applicable, Smart Provider Router ownership, U.S. Equities Processing Boundary, permanent No Execution Boundary, native macOS Apple Silicon delivery, Windows x64 delivery, evidence integrity or Stable immutability.
