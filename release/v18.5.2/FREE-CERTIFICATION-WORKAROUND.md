# DE.PULSE v18.5.2 — Exact-Source Certification Paths

No release gate is waived because of CI budget, runner, or tooling availability. The canonical certification logic remains `release/v18.5.2/run_free_certification.sh`.

## Preferred checkpoint-triggered GitHub lane

`.github/workflows/v18.5.2-hotfix-certification.yml` is checkpoint-triggered rather than ordinary-push triggered. A deliberate fingerprint-excluded request at:

`.depulse-certification/resume/qualification-request.json`

starts the GitHub lane on `v18.5.2-development`. Ordinary source pushes remain quiet, preserving the permanent Adaptive CI cost/qualification contract.

The GitHub lane checks out the exact triggering SHA, installs only the lightweight Python Playwright driver, uses the runner's installed Google Chrome, executes the canonical script with `DEPULSE_EXPECTED_SHA` bound to `github.sha`, and retains the resulting `.depulse-certification/manual-v18.5.2/<source-sha>/` evidence as a workflow artifact.

If Actions billing, runner availability, or permissions prevent the checkpoint lane from starting, classify that as CI/infra state rather than an application failure and use the local lane below. Required evidence is never waived.

## Local fallback lane

From a clean checkout of the pull-request head:

```bash
git switch v18.5.2-development
git pull --ff-only
DEPULSE_EXPECTED_SHA="$(git rev-parse HEAD)" \
  bash release/v18.5.2/run_free_certification.sh
```

The runner performs:

- canonical identity and version consistency;
- GitHub-backed adaptive resume/portability enforcement;
- adaptive functionality utility and provider-to-Market-Mode integration enforcement;
- `gofmt` drift detection and `go vet`;
- full, race-enabled, and randomized Go suites;
- renderer syntax and deterministic-equivalence contracts;
- all v18.5.1 recovery browser tests;
- v18.5.2 tracked-symbol, configurable-identity, and Settings save-bar Chromium tests.

Evidence is written under `.depulse-certification/manual-v18.5.2/<source-sha>/`.

## Release path

The G0–G16 model remains authoritative:

1. G0/G1 — verify canonical release identity and scope/baseline truth.
2. G2–G10 — complete exact-source qualification and coverage reconciliation.
3. G11/G12 — bind/freeze the qualifying candidate and complete authoritative certification according to the active release checkpoint.
4. G13/G14 — package and runtime-audit native macOS Apple Silicon and Windows x64 artifacts from the certified candidate.
5. G15 — promote only after source, browser, package, provenance, runtime and artifact-hash evidence pass.
6. G16 — publish the final retrospective/handoff and durable recovery state.

A G0–G12 source/browser PASS does not authorize Stable promotion without required G13/G14 native package/runtime evidence and G15 assurance.
