#!/usr/bin/env python3
from __future__ import annotations

from pathlib import Path
import subprocess
import sys

ALLOWED = {"ci-fast.yml", "ci-qualified.yml", "release.yml"}
FORBIDDEN_FRAGMENTS = (
    "-retry",
    "-monitor",
    "-probe",
    "-recovery",
    "-certification",
    "-publish",
)


def run_gate(root: Path, filename: str, label: str) -> int:
    gate = root / filename
    if not gate.is_file():
        print("DE.PULSE workflow policy: FAIL", file=sys.stderr)
        print(f"missing {filename}", file=sys.stderr)
        return 1
    result = subprocess.run([sys.executable, str(gate)], cwd=root, check=False)
    if result.returncode != 0:
        print("DE.PULSE workflow policy: FAIL", file=sys.stderr)
        print(f"{label} failed", file=sys.stderr)
        return result.returncode
    return 0


def release_dispatch_contract(workflows: Path) -> int:
    ci_fast = (workflows / "ci-fast.yml").read_text(encoding="utf-8")
    release = (workflows / "release.yml").read_text(encoding="utf-8")
    root = workflows.parents[1]
    verifier_path = root / "tools" / "release" / "verify_promotion_evidence.py"
    verifier = verifier_path.read_text(encoding="utf-8") if verifier_path.is_file() else ""

    ci_required = (
        "types: [opened, synchronize, reopened, closed]",
        "endsWith(github.ref, '-release-certification')",
        "endsWith(github.ref, '-stable-promotion')",
        "github.event.sender.login == github.repository_owner",
        "github.event_name == 'pull_request'",
        "github.event.action != 'closed'",
        "endsWith(github.event.pull_request.base.ref, '-release-certification')",
        "github.event.action == 'closed'",
        "github.event.pull_request.merged == true",
        "endsWith(github.event.pull_request.base.ref, '-stable-promotion')",
        "github.event.pull_request.merged_by.login == github.repository_owner",
        "github.actor == github.repository_owner",
        'release_ref="${PR_BASE_REF:-}"',
        'event_sha="${PR_MERGE_SHA:-}"',
        '"v${release_line}-release-certification")',
        '"v${release_line}-stable-promotion")',
        '.depulse-certification/resume/release-evidence-checkpoint.json',
        "cur['releaseCandidateCommit']",
        "cur['sourceFingerprint']",
        "cur['canonicalReleaseRun']",
        "cur['promotionState']=='READY_NOT_PROMOTED'",
        "ev['G16']['status']=='PASS_CLOSED'",
        "certification_run_id",
        "promotion_sha",
        'repos/${GITHUB_REPOSITORY}/issues/${TRACKING_PR}/comments',
    )
    ci_missing = [fragment for fragment in ci_required if fragment not in ci_fast]
    ci_forbidden = (
        "github.event.head_commit.author.username",
        "gh pr comment",
        "endsWith(github.event.pull_request.base.ref, '-stable-promotion') && github.event.action != 'closed'",
    )
    ci_present_forbidden = [fragment for fragment in ci_forbidden if fragment in ci_fast]

    release_required = (
        "certification_run_id:",
        "Publish already-certified assets without rebuilding them",
        "if: ${{ !inputs.publish }}",
        "G15 Promotion / exact no-rebuild publication",
        "needs.g15.result == 'skipped'",
        "github-token: ${{ github.token }}",
        "repository: ${{ github.repository }}",
        "run-id: ${{ inputs.certification_run_id }}",
        "tools/release/verify_promotion_evidence.py",
        '--certified-run-head "$CANDIDATE_SHA"',
        '--source-fingerprint "$EXPECTED_FP"',
        "--out promotion-verification.json",
        "cat > release-notes.md <<'EOF'",
        "cur['promotionState']=='READY_NOT_PROMOTED'",
        "releases/tags/$tag",
        "git/refs/tags/$tag",
        "gh api --method DELETE",
        "gh release create",
        "gh release upload",
        'release-assets/macos/De-Pulse-v${{ inputs.version }}-Stable-macOS-Apple-Silicon.zip',
        'release-assets/windows/De-Pulse-v${{ inputs.version }}-Stable-Windows-x64.zip',
        "release-assets/macos/G13-G14-macOS-Apple-Silicon.json",
        "release-assets/windows/G13-G14-Windows-x64.json",
        "release-assets/g15/G15-Release-Assurance.json",
        "promotion-verification.json",
        "G12/G13/G14/G15 are not rerun in promotion mode",
        '"noRebuildPublication": true',
    )
    release_missing = [fragment for fragment in release_required if fragment not in release]
    release_forbidden = (
        "git merge-base --is-ancestor '${{ inputs.candidate_sha }}'",
        'git merge-base --is-ancestor "$CANDIDATE_SHA"',
        "cat > release-notes.md <<EOF\n          # DE.PULSE",
        "release-assets/macos/*",
        "release-assets/windows/*",
        "release-assets/g15/*",
    )
    release_present_forbidden = [fragment for fragment in release_forbidden if fragment in release]

    verifier_required = (
        "DE.PULSE-STABLE-PROMOTION-VERIFY-1",
        "DE.PULSE-G15-ASSURANCE-2",
        "DE.PULSE-G13-G14-NATIVE-2",
        "certifiedSourceSha",
        "sourceFingerprint",
        "artifactSha256",
        "noExecutionBoundary",
        "promotionAuthorized",
        "all(v == 'PASS' for v in evidence_checks.values())",
        "noRebuild",
    )
    verifier_missing = [fragment for fragment in verifier_required if fragment not in verifier]
    if not verifier_path.is_file():
        verifier_missing.insert(0, "tools/release/verify_promotion_evidence.py")

    for job_marker in ("g12:\n", "macos:\n", "windows:\n", "g15:\n"):
        pos = release.find(job_marker)
        if pos < 0:
            release_missing.append(f"job marker {job_marker.strip()}")
            continue
        window = release[pos:pos + 500]
        if "!inputs.publish" not in window:
            release_missing.append(f"{job_marker.strip()} publish=false guard")

    if ci_missing or ci_present_forbidden or release_missing or release_present_forbidden or verifier_missing:
        print("DE.PULSE workflow policy: FAIL", file=sys.stderr)
        if ci_missing:
            print("release dispatcher contract missing: " + ", ".join(ci_missing), file=sys.stderr)
        if ci_present_forbidden:
            print("release dispatcher forbidden contract fragments: " + ", ".join(ci_present_forbidden), file=sys.stderr)
        if release_missing:
            print("no-rebuild release contract missing: " + ", ".join(release_missing), file=sys.stderr)
        if release_present_forbidden:
            print("no-rebuild release contract forbidden fragments: " + ", ".join(release_present_forbidden), file=sys.stderr)
        if verifier_missing:
            print("promotion evidence verifier contract missing: " + ", ".join(verifier_missing), file=sys.stderr)
        return 1

    print("release dispatcher certification + merged-PR Stable promotion authorization: PASS")
    print("cross-run exact-artifact no-rebuild Stable promotion contract: PASS")
    print("quoted release notes + fail-closed stale partial tag recovery: PASS")
    print("explicit unique Stable release asset allowlist: PASS")
    return 0


def g12_browser_contract(root: Path) -> int:
    full_cert_path = root / "release" / "v18.6.0" / "run_full_certification.sh"
    hierarchy_path = root / "release" / "v18.6.0" / "browser_ui_hierarchy_test.py"
    if not full_cert_path.is_file() or not hierarchy_path.is_file():
        print("DE.PULSE workflow policy: FAIL", file=sys.stderr)
        print("v18.6 G12 browser contract files are missing", file=sys.stderr)
        return 1

    full_cert = full_cert_path.read_text(encoding="utf-8")
    hierarchy = hierarchy_path.read_text(encoding="utf-8")
    required = (
        "release/v18.6.0/browser_ui_hierarchy_test.py",
        "release/v18.6.0/browser_watchlist_membership_test.py",
    )
    forbidden = (
        "release/v18.5.1/browser_ui_hierarchy_test.py",
        "release/v18.5.1/browser_watchlist_membership_test.py",
    )
    missing = [fragment for fragment in required if fragment not in full_cert]
    stale = [fragment for fragment in forbidden if fragment in full_cert]
    hierarchy_required = (
        'RELEASE_VERSION = json.loads((ROOT / "release_identity.json")',
        'f"ui-v18.5.1.css?v={RELEASE_VERSION}"',
        'f"header-v18.5.1.js?v={RELEASE_VERSION}"',
    )
    hierarchy_missing = [fragment for fragment in hierarchy_required if fragment not in hierarchy]
    if "?v=18.5.2" in hierarchy:
        hierarchy_missing.append("no hard-coded ?v=18.5.2 cache-buster")

    if missing or stale or hierarchy_missing:
        print("DE.PULSE workflow policy: FAIL", file=sys.stderr)
        if missing:
            print("G12 browser contract missing current proofs: " + ", ".join(missing), file=sys.stderr)
        if stale:
            print("G12 browser contract still invokes superseded proofs: " + ", ".join(stale), file=sys.stderr)
        if hierarchy_missing:
            print("v18.6 hierarchy proof contract mismatch: " + ", ".join(hierarchy_missing), file=sys.stderr)
        return 1

    print("v18.6 G12 browser proof selection and release cache-buster contract: PASS")
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
    forbidden = [name for name in present if any(x in name for x in FORBIDDEN_FRAGMENTS)]

    if missing or unexpected or forbidden:
        print("DE.PULSE workflow policy: FAIL", file=sys.stderr)
        if missing:
            print("missing canonical workflows: " + ", ".join(missing), file=sys.stderr)
        if unexpected:
            print("unexpected active workflows: " + ", ".join(unexpected), file=sys.stderr)
        if forbidden:
            print("forbidden one-off workflow naming: " + ", ".join(forbidden), file=sys.stderr)
        return 1

    if release_dispatch_contract(workflows) != 0:
        return 1
    if g12_browser_contract(root) != 0:
        return 1
    if run_gate(root, "dependency_readiness_gate.py", "dependency/provider readiness contract") != 0:
        return 1
    if run_gate(root, "ai_continuous_eval_gate.py", "AI continuous eval/rights contract") != 0:
        return 1

    print("DE.PULSE workflow policy: PASS")
    print("active workflows: " + ", ".join(present))
    print("G12 browser proof selection: PASS")
    print("dependency/provider readiness: PASS")
    print("AI continuous eval/rights: PASS")
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
