# DE.PULSE Repository Instructions

These instructions apply to the entire repository and to every ChatGPT/Codex account or compatible coding agent.

## Start here — never from model memory

Before planning, editing or claiming status:

1. Read `governance/AI-ASSISTANT-PORTABILITY-CONTRACT.md`.
2. Read `governance/CI-EFFICIENCY-CONTRACT.md`; its branch/PR/Actions rules are mandatory and exist to prevent CI event amplification and unnecessary runner cost.
3. Read `governance/README.md` and its source-of-truth hierarchy.
4. Read `handoff/CURRENT.md`.
5. Inspect the actual GitHub default branch, active branch, open PR, current HEAD, latest Stable tag/release, checks and artifacts.
6. Read `release_identity.json`, `.depulse-certification/resume/build-checkpoint.json` and `.depulse-certification/resume/release-evidence-checkpoint.json`.
7. Read `governance/ROADMAP.md`, `adaptive-governance/ADAPTIVE_BUILD_PLAN.md`, `ADAPTIVE_BUILD_PROCESS.md` and `ADAPTIVE_DELIVERY_PROCESS.md`.
8. Run or inspect `python3 adaptive_resume_gate.py` and `python3 tools/ci/workflow_policy.py`.
9. Reconcile disagreements and resume from the last trustworthy PASS / earliest required G0–G16 gate.

Actual GitHub source, immutable tags/releases, PR/check/artifact state and package evidence outrank chat summaries. Do not ask the user to restate context already committed in GitHub.

## Permanent constraints

- GitHub is the durable source of truth; temporary workspaces and conversation memory are disposable.
- Use G0–G16 only. Do not create G17+.
- Follow LOOKUP → COMPARE → CLASSIFY → DECIDE → UPDATE.
- Follow REUSE → CONSOLIDATE → REFACTOR → DELETE/REPLACE → ADD.
- Preserve the U.S. Equities Processing Boundary, No Execution Boundary, deterministic Market Mode ownership and governed SHADOW → VALIDATED → APPROVED → PRODUCTION lifecycle.
- Never commit secrets or copy user API keys into continuity files, tests, logs or prompts.
- Keep requirements open until code, behavior and required package evidence prove closure.
- Normal releases use one `v<version>-development` branch and one PR to `main`; do not create release-certification, Stable-promotion, trigger, dispatch, retry or fallback branches/PRs.
- CI retries must rerun failed jobs or update the same branch/PR. Creating GitHub events is never a reason to create a commit, branch or PR.
- CI Fast must not dispatch Release. Qualified is candidate-only. Release is the single merged-release-PR G11–G16 certify-and-publish path.

## Before ending meaningful work

Commit durable changes to the active branch; update `handoff/CURRENT.md`; update the machine checkpoint when a release candidate exists; keep exactly one next action; and ensure another ChatGPT/Codex/Claude account can resume without this conversation.
