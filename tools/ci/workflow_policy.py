#!/usr/bin/env python3
from __future__ import annotations

import os
from pathlib import Path
import subprocess
import sys

ALLOWED = {"ci-fast.yml", "ci-qualified.yml", "release.yml"}
FORBIDDEN_WORKFLOW_NAME_FRAGMENTS = (
    "-retry",
    "-monitor",
    "-probe",
    "-recovery",
    "-certification",
    "-publish",
)
DEPRECATED_RELEASE_BRANCH_MARKERS = (
    "-release-certification",
    "-stable-promotion",
    "-certification-trigger",
    "-cert-trigger",
    "-promotion-trigger",
    "-dispatch",
    "-retry",
    "-fallback",
)


def fail(label: str, items: list[str]) -> int:
    print("DE.PULSE workflow policy: FAIL", file=sys.stderr)
    print(f"{label}: " + ", ".join(items), file=sys.stderr)
    return 1


def run_gate(root: Path, filename: str, label: str) -> int:
    gate = root / filename
    if not gate.is_file():
        return fail("missing gate", [filename])
    result = subprocess.run([sys.executable, str(gate)], cwd=root, check=False)
    if result.returncode != 0:
        print("DE.PULSE workflow policy: FAIL", file=sys.stderr)
        print(f"{label} failed", file=sys.stderr)
        return result.returncode
    return 0


def branch_name_contract() -> int:
    branch = (os.environ.get("GITHUB_HEAD_REF") or os.environ.get("GITHUB_REF_NAME") or "").strip()
    if not branch:
        return 0
    if branch.startswith("v") and any(marker in branch for marker in DEPRECATED_RELEASE_BRANCH_MARKERS):
        return fail("deprecated release-temp branch naming is prohibited", [branch])
    return 0


def canonical_workflow_contract(workflows: Path) -> int:
    ci_fast = (workflows / "ci-fast.yml").read_text(encoding="utf-8")
    qualified = (workflows / "ci-qualified.yml").read_text(encoding="utf-8")
    release = (workflows / "release.yml").read_text(encoding="utf-8")

    fast_required = (
        "types: [opened, synchronize, reopened]",
        "- main",
        "workflow_dispatch:",
        "cancel-in-progress: true",
        "branch-hygiene:",
        "tools/ci/branch_hygiene.py --apply",
        "release dispatch: NOT PERMITTED from CI Fast",
    )
    fast_forbidden = (
        "types: [opened, synchronize, reopened, closed]",
        "'v*-development'",
        "'v*-release-certification'",
        "'v*-stable-promotion'",
        "release-dispatch:",
        "gh workflow run release.yml",
        "actions: write",
        "DE-PULSE-ci-impact-${{ github.sha }}",
    )
    missing = [x for x in fast_required if x not in ci_fast]
    forbidden = [x for x in fast_forbidden if x in ci_fast]
    if missing:
        return fail("CI Fast efficiency contract missing", missing)
    if forbidden:
        return fail("CI Fast duplicate-trigger/dispatcher contract violated", forbidden)

    qualified_required = (
        "types: [ready_for_review]",
        "workflow_dispatch:",
        "workflow_call:",
        "PR_HEAD_SHA: ${{ github.event.pull_request.head.sha }}",
        "needs.context.outputs.lane == 'ci-harness'",
        "ci-harness) require \"$PORTABILITY\"",
        "full) require \"$BACKEND\"; require \"$RENDERER\"; require \"$BROWSER\"",
    )
    qualified_forbidden = (
        "paths:",
        "types: [opened",
        "types: [synchronize",
        "types: [closed",
    )
    missing = [x for x in qualified_required if x not in qualified]
    forbidden = [x for x in qualified_forbidden if x in qualified]
    if missing:
        return fail("CI Qualified candidate-only contract missing", missing)
    if forbidden:
        return fail("CI Qualified must not run on routine development updates", forbidden)

    release_required = (
        "types: [closed]",
        "- 'release_identity.json'",
        "github.event.pull_request.merged == true",
        "github.event.pull_request.merge_commit_sha",
        "G12 Full certification",
        "G13/G14 macOS Apple Silicon",
        "G13/G14 Windows x64",
        "G15 Release Assurance",
        "Publish exact same-run certified artifacts",
        "tools/release/verify_promotion_evidence.py",
        "gh release create",
        "gh release upload",
        '"mode": "SINGLE_RUN_CERTIFY_AND_PUBLISH"',
        '"noRebuildPublication": true',
    )
    release_forbidden = (
        "workflow_dispatch:",
        "certification_run_id",
        "promotion_sha",
        "stable-promotion",
        "release-certification",
        "gh workflow run",
        "run-id:",
        "publish:\n",
    )
    missing = [x for x in release_required if x not in release]
    forbidden = [x for x in release_forbidden if x in release]
    if missing:
        return fail("single-run Release G11-G16 contract missing", missing)
    if forbidden:
        return fail("legacy release dispatch/promotion contract still present", forbidden)

    print("CI Fast single-event development contract: PASS")
    print("CI Qualified ready-candidate contract: PASS")
    print("Release single merged-PR certify-and-publish contract: PASS")
    print("premium runner separation: PASS")
    return 0


def main() -> int:
    root = Path(__file__).resolve().parents[2]
    workflows = root / ".github" / "workflows"
    present = sorted(
        p.name
        for p in workflows.glob("*")
        if p.is_file() and p.suffix.lower() in {".yml", ".yaml"}
    )
    unexpected = [name for name in present if name not in ALLOWED]
    missing = sorted(ALLOWED - set(present))
    forbidden = [name for name in present if any(x in name for x in FORBIDDEN_WORKFLOW_NAME_FRAGMENTS)]

    if missing:
        return fail("missing canonical workflows", missing)
    if unexpected:
        return fail("unexpected active workflows", unexpected)
    if forbidden:
        return fail("forbidden one-off workflow naming", forbidden)
    if branch_name_contract() != 0:
        return 1
    if canonical_workflow_contract(workflows) != 0:
        return 1
    if run_gate(root, "dependency_readiness_gate.py", "dependency/provider readiness contract") != 0:
        return 1
    if run_gate(root, "ai_continuous_eval_gate.py", "AI continuous eval/rights contract") != 0:
        return 1

    print("DE.PULSE workflow policy: PASS")
    print("active workflows: " + ", ".join(present))
    print("branch/retry event-amplification prevention: PASS")
    print("dependency/provider readiness: PASS")
    print("AI continuous eval/rights: PASS")
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
