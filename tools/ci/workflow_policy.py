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


def require_tokens(label: str, text: str, required: tuple[str, ...], forbidden: tuple[str, ...] = ()) -> int:
    missing = [token for token in required if token not in text]
    bad = [token for token in forbidden if token in text]
    if missing:
        return fail(f"{label} missing", missing)
    if bad:
        return fail(f"{label} forbidden contract", bad)
    return 0


def canonical_workflow_contract(workflows: Path) -> int:
    ci_fast = (workflows / "ci-fast.yml").read_text(encoding="utf-8")
    qualified = (workflows / "ci-qualified.yml").read_text(encoding="utf-8")
    release = (workflows / "release.yml").read_text(encoding="utf-8")
    root = workflows.parents[1]
    branch_hygiene = (root / "tools" / "ci" / "branch_hygiene.py").read_text(encoding="utf-8")

    if require_tokens(
        "CI Fast efficiency/exact-head/capability-test/post-Stable-continuity contract",
        ci_fast,
        (
            "types: [opened, synchronize, reopened]",
            "- main",
            "workflow_dispatch:",
            "cancel-in-progress: true",
            "if: ${{ github.event_name == 'pull_request' || github.event_name == 'workflow_dispatch' }}",
            "ref: ${{ github.event.pull_request.head.sha || github.sha }}",
            "python3 tools/ci/impact_plan.py",
            "python3 tools/ci/post_stable_continuity_gate.py",
            "python3 tools/ci/stable_evidence_gate.py",
            "tools/ci/branch_hygiene.py --apply",
            "node tests/renderer/surface_consolidation_test.js",
            "node tests/renderer/documentation_access_test.js",
            "DE.PULSE/fast-head",
            "release dispatch: NOT PERMITTED from CI Fast",
        ),
        (
            "types: [opened, synchronize, reopened, closed]",
            "gh workflow run release.yml",
            "actions: write",
            "node v18_6_surface_consolidation_test.js",
            "node v18_6_documentation_access_test.js",
        ),
    ) != 0:
        return 1

    if require_tokens(
        "CI Qualified Planner v3/exact-head/dependency/native evidence contract",
        qualified,
        (
            "types: [ready_for_review]",
            "workflow_dispatch:",
            "workflow_call:",
            "base_sha:",
            "target_ref:",
            "--target-ref \"$target\"",
            "--requested-lane \"$requested\"",
            "portability_required: ${{ steps.resolve.outputs.portability_required }}",
            "backend_required: ${{ steps.resolve.outputs.backend_required }}",
            "renderer_required: ${{ steps.resolve.outputs.renderer_required }}",
            "chrome_required: ${{ steps.resolve.outputs.chrome_required }}",
            "webkit_required: ${{ steps.resolve.outputs.webkit_required }}",
            "security_rights_required: ${{ steps.resolve.outputs.security_rights_required }}",
            "db_integration_required: ${{ steps.resolve.outputs.db_integration_required }}",
            "native_macos_required: ${{ steps.resolve.outputs.native_macos_required }}",
            "native_windows_required: ${{ steps.resolve.outputs.native_windows_required }}",
            "if: ${{ needs.context.outputs.backend_required == 'true' }}",
            "if: ${{ needs.context.outputs.renderer_required == 'true' }}",
            "if: ${{ needs.context.outputs.chrome_required == 'true' }}",
            "if: ${{ needs.context.outputs.webkit_required == 'true' }}",
            "name: Primary Chrome behavior",
            "name: Primary WebKit browser compatibility",
            "run: python3 tools/ci/webkit_browser_test.py",
            "name: Qualified macOS native lifecycle rehearsal",
            "bash tools/release/native_macos.sh",
            "name: Qualified Windows native runtime rehearsal",
            "tools/release/native_windows.ps1",
            "name: Qualified persistence / DB integration",
            "name: Qualified security / data-rights contracts",
            "Require Planner v3 selected jobs to pass",
            "needs.context.outputs.selected_jobs",
            "actions: read",
            "Collect CI telemetry",
            "tools/ci/ci_telemetry.py",
            "DE-PULSE-ci-telemetry-${{ github.run_id }}-${{ needs.context.outputs.sha }}",
            "retention-days: 30",
            "DE.PULSE/qualified-head",
        ),
        (
            "paths:",
            "types: [opened",
            "types: [synchronize",
            "types: [closed",
            "browser: [chromium, webkit]",
            "playwright install --with-deps firefox",
            "run: python3 tools/ci/webkit_targeted_test.py",
        ),
    ) != 0:
        return 1

    if require_tokens(
        "single-run Release G11-G16 contract",
        release,
        (
            "types: [closed]",
            "- 'release_identity.json'",
            "- '.github/workflows/release.yml'",
            "statuses: read",
            "github.event.pull_request.merged == true",
            "github.event.pull_request.base.ref == 'main'",
            "startsWith(github.event.pull_request.head.ref, 'v')",
            "endsWith(github.event.pull_request.head.ref, '-development')",
            "github.event.pull_request.merge_commit_sha",
            "github.event.pull_request.head.sha",
            "DE.PULSE/fast-head",
            "DE.PULSE/qualified-head",
            "require_status",
            'test "$source_fp" = "$candidate_fp"',
            "G12 Full certification",
            "G13/G14 macOS Apple Silicon",
            "G13/G14 Windows x64",
            "G15 Release Assurance",
            "Publish exact same-run certified artifacts",
            "tools/release/verify_promotion_evidence.py",
            "git ls-remote --refs origin",
            "gh release create",
            "gh release upload",
            '"mode": "SINGLE_RUN_CERTIFY_AND_PUBLISH"',
            '"noRebuildPublication": true',
        ),
        (
            "workflow_dispatch:",
            "certification_run_id:",
            "promotion_sha:",
            "stable-promotion",
            "release-certification",
            "gh workflow run",
            "gh run list",
            "READY_NOT_PROMOTED",
            "git/ref/tags/$tag",
        ),
    ) != 0:
        return 1

    hygiene_required = (
        'pr_heads("open")',
        'pr_heads("merged")',
        "MERGED_PR_HEAD",
        "STABLE_LINE_ALREADY_PUBLISHED",
        "stable_line_closed",
        "DE.PULSE-BRANCH-HYGIENE-3",
    )
    missing = [token for token in hygiene_required if token not in branch_hygiene]
    if missing:
        return fail("branch hygiene squash/stable-line contract missing", missing)

    print("CI Fast single-event exact-head development contract: PASS")
    print("CI Qualified Planner v3 deterministic job selection: PASS")
    print("Qualified trustworthy merge-base/manual target binding: PASS")
    print("Chrome + WebKit browser evidence ownership: PASS")
    print("macOS + Windows native rehearsal ownership separation: PASS")
    print("Qualified DB + security/data-rights dependency evidence: PASS")
    print("Qualified telemetry/evidence retention contract: PASS")
    print("Release exact G10-head status / merged-candidate evidence binding: PASS")
    print("Release single merged-PR certify-and-publish contract: PASS")
    print("squash-merged/stable-line branch hygiene: PASS")
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

    gates = (
        ("tools/ci/workflow_structural_lint.py", "zero-network workflow structural lint"),
        ("tools/ci/impact_plan_self_test.py", "CI impact planner v3 contract"),
        ("tools/ci/legacy_test_gate_inventory.py", "legacy test/gate inventory contract"),
        ("tools/ci/reproducibility_gate.py", "CI reproducibility/dependency/permission contract"),
        ("tools/ci/browser_risk_routing_gate.py", "Chrome/WebKit primary browser routing contract"),
        ("tools/ci/renderer_owner_contract.py", "capability-oriented renderer ownership contract"),
        ("tools/ci/ci_telemetry_self_test.py", "CI telemetry/amplification contract"),
        ("tools/ci/post_stable_continuity_gate.py", "post-Stable repository continuity contract"),
        ("tools/ci/stable_evidence_gate.py", "durable Stable evidence contract"),
        ("tools/ci/release_rehearsal.py", "pre-merge release rehearsal contract"),
        ("dependency_readiness_gate.py", "dependency/provider readiness contract"),
        ("ai_continuous_eval_gate.py", "AI continuous eval/rights contract"),
    )
    for filename, label in gates:
        if run_gate(root, filename, label) != 0:
            return 1

    print("DE.PULSE workflow policy: PASS")
    print("active workflows: " + ", ".join(present))
    print("branch/retry event-amplification prevention: PASS")
    print("zero-network workflow structural lint: PASS")
    print("CI impact planner v3 self-test: PASS")
    print("legacy test/gate inventory contract: PASS")
    print("CI reproducibility/dependency/permission contract: PASS")
    print("Chrome/WebKit primary browser routing contract: PASS")
    print("capability-oriented renderer ownership contract: PASS")
    print("CI telemetry/amplification contract: PASS")
    print("post-Stable repository continuity contract: PASS")
    print("durable Stable evidence contract: PASS")
    print("pre-merge release rehearsal: PASS")
    print("dependency/provider readiness: PASS")
    print("AI continuous eval/rights: PASS")
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
