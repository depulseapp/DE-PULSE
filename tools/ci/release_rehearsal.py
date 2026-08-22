#!/usr/bin/env python3
from __future__ import annotations

from pathlib import Path
import re
import subprocess
import sys


def fail(message: str) -> int:
    print(f"DE.PULSE release rehearsal: FAIL - {message}", file=sys.stderr)
    return 1


def decide_tag_action(existing_sha: str, candidate_sha: str) -> str:
    existing = existing_sha.strip()
    candidate = candidate_sha.strip()
    if not candidate:
        raise ValueError("candidate SHA is required")
    if not existing:
        return "CREATE"
    if existing == candidate:
        return "REUSE"
    raise ValueError(f"Stable tag mismatch: existing={existing} candidate={candidate}")


def extract_publish_block(text: str) -> str:
    marker = "\n  publish:\n"
    start = text.find(marker)
    if start < 0:
        raise ValueError("publish job missing")
    tail = text[start + 1 :]
    g16 = tail.find("\n  g16:\n")
    return tail if g16 < 0 else tail[:g16]


def main() -> int:
    root = Path(__file__).resolve().parents[2]
    release_path = root / ".github" / "workflows" / "release.yml"
    if not release_path.is_file():
        return fail("canonical release workflow missing")

    text = release_path.read_text(encoding="utf-8")
    required = (
        "types: [closed]",
        "- 'release_identity.json'",
        "- '.github/workflows/release.yml'",
        "group: depulse-stable-release",
        "cancel-in-progress: false",
        "github.event.pull_request.merged == true",
        "github.event.pull_request.base.ref == 'main'",
        "github.event.pull_request.merge_commit_sha",
        "github.event.pull_request.head.sha",
        "DE.PULSE/fast-head",
        "DE.PULSE/qualified-head",
        'test "$source_fp" = "$candidate_fp"',
        "certification-manifest.json",
        "tools/release/run_full_certification.py",
        "G12 Full certification",
        "G13/G14 macOS Apple Silicon",
        "G13/G14 Windows x64",
        "G15 Release Assurance",
        "Publish exact same-run certified artifacts",
        "tools/release/verify_promotion_evidence.py",
        'git ls-remote --refs origin "refs/tags/$tag"',
        'existing_ref" != "$CANDIDATE_SHA"',
        'gh release create "$tag"',
        '--target "$CANDIDATE_SHA"',
        'gh release upload "$tag"',
        "immutable Stable asset conflict",
        "existing Stable asset byte-identical; reuse",
        "G16 Adaptive release handoff evidence",
    )
    missing = [token for token in required if token not in text]
    if missing:
        return fail("missing release invariant(s): " + ", ".join(missing))
    if "--clobber" in text:
        return fail("Stable release publication must never use --clobber")

    try:
        publish = extract_publish_block(text)
    except ValueError as exc:
        return fail(str(exc))

    forbidden_build_patterns = (
        r"\bgo\s+build\b",
        r"\bgo\s+test\b",
        r"native_macos\.sh",
        r"native_windows\.ps1",
    )
    for pattern in forbidden_build_patterns:
        if re.search(pattern, publish):
            return fail(f"publish job rebuild risk matched {pattern!r}")

    candidate = "abc123"
    if decide_tag_action("", candidate) != "CREATE":
        return fail("absent tag model must CREATE")
    if decide_tag_action(candidate, candidate) != "REUSE":
        return fail("same-SHA tag model must REUSE")
    try:
        decide_tag_action("different", candidate)
    except ValueError:
        pass
    else:
        return fail("mismatched Stable tag model must fail closed")

    hardening = subprocess.run(
        [sys.executable, str(root / "tools" / "ci" / "ci_hardening_gate.py")],
        cwd=root,
        check=False,
    )
    if hardening.returncode != 0:
        return fail("CI/release identity, versioning, toolchain hardening contract failed")

    print("DE.PULSE pre-merge release rehearsal: PASS")
    print("G11 exact-head status + fingerprint contract: PASS")
    print("G12/G13/G14/G15/G16 topology contract: PASS")
    print("same-run no-rebuild publication contract: PASS")
    print("Stable tag absent/same/mismatch model: PASS")
    print("identity/versioning/toolchain hardening contracts: PASS")
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
