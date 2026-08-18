# DE.PULSE v18.5.2 — No-Cost Certification Workaround

GitHub Actions automatic execution is paused because the repository Actions budget is exhausted. This does not waive any release gate and is not recorded as an application failure.

## Release path

The existing G0–G16 model remains authoritative:

1. G0/G1 — verify canonical release identity and frozen scope.
2. G10 — commit the implementation to `v18.5.2-development`.
3. G11 — bind certification to the exact clean commit SHA.
4. G12 — run the full local source and Chromium certification lane.
5. G13/G14 — package and runtime-audit the native macOS and Windows applications from that same SHA.
6. G15 — promote only after source, browser, package, provenance, and artifact hashes pass.
7. G16 — publish the final handoff and recovery index.

## One-command local lane

From a clean checkout of the pull-request head:

```bash
git switch v18.5.2-development
git pull --ff-only
DEPULSE_EXPECTED_SHA="$(git rev-parse HEAD)" \
  bash release/v18.5.2/run_free_certification.sh
```

The runner performs:

- canonical identity and version consistency;
- `gofmt` drift detection and `go vet`;
- full, race-enabled, and randomized Go suites;
- renderer syntax and deterministic-equivalence contracts;
- all v18.5.1 recovery browser tests;
- v18.5.2 tracked-symbol, configurable-identity, and Settings save-bar Chromium tests.

Evidence is written outside product source under
`.depulse-certification/manual-v18.5.2/<source-sha>/`.

## Required local tools

- Git
- the repository's supported Go toolchain
- Node.js
- Python 3 with Playwright
- installed Google Chrome or Chromium

If Python Playwright is missing, install only the lightweight driver package:

```bash
python3 -m pip install --user playwright
```

The runner uses the already installed Chrome/Chromium executable and does not download a browser.

## Actions policy while budget is unavailable

`.github/workflows/v18.5.2-hotfix-certification.yml` is retained but has only
`workflow_dispatch`. Push and pull-request triggers are disabled so ordinary
development does not consume Actions minutes. Do not manually dispatch it until
budget is restored.

A local PASS is source/browser evidence only. It does not authorize Stable promotion without the G13/G14 native package and runtime audit.
