#!/usr/bin/env python3
"""Early G11 Stable-publication feasibility for DE.PULSE.

This runs before expensive G12/G13/G14 work. It rejects impossible branch,
version, predecessor, tag, or release states and classifies exact already-
published candidates as idempotent reuse so they are not rebuilt or republished.
"""
from __future__ import annotations

import argparse
import json
import os
from pathlib import Path
import re
import subprocess
import sys
from typing import Any

ROOT = Path(__file__).resolve().parents[2]
STABLE_TAG_RE = re.compile(r"^v(\d+)\.(\d+)\.(\d+)-stable$")
SHA_RE = re.compile(r"^[0-9a-f]{40}$")


def run(argv: list[str], *, check: bool = True) -> subprocess.CompletedProcess[str]:
    return subprocess.run(argv, cwd=ROOT, text=True, capture_output=True, check=check)


def load_identity() -> dict[str, Any]:
    return json.loads((ROOT / "release_identity.json").read_text(encoding="utf-8"))


def previous_stable_tag(identity: dict[str, Any]) -> str:
    raw = str(identity.get("previous_stable", "")).strip()
    if not raw:
        raise ValueError("release_identity.previous_stable is required")
    if raw.endswith("-stable"):
        return raw
    if not raw.startswith("v"):
        raw = "v" + raw
    return raw + "-stable"


def tag_key(tag: str) -> tuple[int, int, int]:
    match = STABLE_TAG_RE.fullmatch(tag)
    if not match:
        raise ValueError(f"not a canonical historical Stable tag: {tag}")
    return tuple(int(x) for x in match.groups())


def fetch_tags() -> None:
    result = run(["git", "fetch", "--tags", "--force", "origin"], check=False)
    if result.returncode != 0:
        raise RuntimeError("unable to refresh remote Stable tags: " + result.stderr.strip())


def resolve_tag_commit(tag: str) -> str:
    result = run(["git", "rev-list", "-n", "1", f"refs/tags/{tag}"], check=False)
    return result.stdout.strip() if result.returncode == 0 else ""


def latest_stable_tag() -> str:
    result = run(["git", "tag", "--list", "v*-stable"], check=True)
    tags = [line.strip() for line in result.stdout.splitlines() if STABLE_TAG_RE.fullmatch(line.strip())]
    return max(tags, key=tag_key) if tags else ""


def release_metadata(repo: str, tag: str) -> dict[str, Any] | None:
    result = run(
        [
            "gh",
            "release",
            "view",
            tag,
            "--repo",
            repo,
            "--json",
            "tagName,targetCommitish,isDraft,isPrerelease,assets",
        ],
        check=False,
    )
    if result.returncode != 0:
        return None
    return json.loads(result.stdout)


def decide(
    *,
    identity: dict[str, Any],
    candidate_sha: str,
    source_head_ref: str,
    existing_tag_sha: str,
    latest_tag: str,
    release: dict[str, Any] | None,
) -> dict[str, Any]:
    version = str(identity.get("version", "")).strip()
    channel = str(identity.get("channel", "")).strip()
    if channel != "STABLE":
        raise ValueError(f"release channel must be STABLE, got {channel!r}")
    if not version:
        raise ValueError("release identity version is required")
    if not SHA_RE.fullmatch(candidate_sha):
        raise ValueError("candidate SHA must be a 40-hex Git commit")

    expected_ref = f"v{version}-development"
    if source_head_ref != expected_ref:
        raise ValueError(f"release PR head must be {expected_ref}, got {source_head_ref!r}")

    baseline = str(identity.get("stable_baseline", "")).strip()
    previous = str(identity.get("previous_stable", "")).strip()
    if not baseline or not previous or baseline != previous:
        raise ValueError(
            f"stable_baseline and previous_stable must agree before release: {baseline!r} vs {previous!r}"
        )

    tag = f"v{version}-stable"
    previous_tag = previous_stable_tag(identity)
    if existing_tag_sha and existing_tag_sha != candidate_sha:
        raise ValueError(
            f"immutable Stable tag collision: {tag} points to {existing_tag_sha}, candidate is {candidate_sha}"
        )

    if not existing_tag_sha:
        if latest_tag != previous_tag:
            raise ValueError(
                f"latest Stable predecessor mismatch: expected {previous_tag}, observed {latest_tag or '<none>'}"
            )
        if release is not None:
            raise ValueError(f"Stable release {tag} exists without a resolvable matching tag")
        mode = "BUILD_AND_PUBLISH"
    else:
        if latest_tag not in {tag, previous_tag}:
            raise ValueError(
                f"latest Stable line is incompatible with target {tag}: observed {latest_tag or '<none>'}"
            )
        if release is None:
            mode = "BUILD_AND_PUBLISH"
        else:
            if bool(release.get("isDraft")) or bool(release.get("isPrerelease")):
                raise ValueError(f"existing Stable release {tag} must not be draft/prerelease")
            if str(release.get("tagName", tag)) != tag:
                raise ValueError(f"release tag identity mismatch for {tag}")
            mode = "REUSE_PUBLISHED"

    manifest = f"release/v{version}/certification-manifest.json"
    return {
        "schema": "DE.PULSE-G11-PUBLICATION-PREFLIGHT-1",
        "productVersion": version,
        "tag": tag,
        "previousStableTag": previous_tag,
        "candidateSha": candidate_sha,
        "sourceHeadRef": source_head_ref,
        "existingTagSha": existing_tag_sha or None,
        "latestStableTag": latest_tag or None,
        "releaseExists": release is not None,
        "publicationMode": mode,
        "certificationManifest": manifest,
    }


def self_test() -> int:
    identity = {
        "version": "18.9.1",
        "channel": "STABLE",
        "stable_baseline": "v18.9.0",
        "previous_stable": "v18.9.0",
    }
    candidate = "a" * 40
    created = decide(
        identity=identity,
        candidate_sha=candidate,
        source_head_ref="v18.9.1-development",
        existing_tag_sha="",
        latest_tag="v18.9.0-stable",
        release=None,
    )
    assert created["publicationMode"] == "BUILD_AND_PUBLISH"
    reused = decide(
        identity=identity,
        candidate_sha=candidate,
        source_head_ref="v18.9.1-development",
        existing_tag_sha=candidate,
        latest_tag="v18.9.1-stable",
        release={"tagName": "v18.9.1-stable", "isDraft": False, "isPrerelease": False},
    )
    assert reused["publicationMode"] == "REUSE_PUBLISHED"
    try:
        decide(
            identity=identity,
            candidate_sha=candidate,
            source_head_ref="v18.9.1-development",
            existing_tag_sha="b" * 40,
            latest_tag="v18.9.1-stable",
            release=None,
        )
    except ValueError:
        pass
    else:
        raise AssertionError("mismatched Stable tag must fail closed")
    try:
        decide(
            identity=identity,
            candidate_sha=candidate,
            source_head_ref="v18.9.1-development",
            existing_tag_sha="",
            latest_tag="v18.8.2-stable",
            release=None,
        )
    except ValueError:
        pass
    else:
        raise AssertionError("wrong Stable predecessor must fail closed")
    print("DE.PULSE G11 publication preflight self-test: PASS")
    print("new publication feasibility: PASS")
    print("exact already-published idempotent reuse: PASS")
    print("tag collision / predecessor mismatch fail-closed: PASS")
    return 0


def main() -> int:
    parser = argparse.ArgumentParser(description="DE.PULSE early G11 publication feasibility")
    parser.add_argument("--repo", default=os.environ.get("GITHUB_REPOSITORY", ""))
    parser.add_argument("--candidate-sha", default="")
    parser.add_argument("--source-head-ref", default="")
    parser.add_argument("--json-out", default="")
    parser.add_argument("--self-test", action="store_true")
    args = parser.parse_args()

    if args.self_test:
        return self_test()
    if not args.repo:
        raise SystemExit("--repo or GITHUB_REPOSITORY is required")
    if not args.candidate_sha or not args.source_head_ref:
        raise SystemExit("--candidate-sha and --source-head-ref are required")

    identity = load_identity()
    version = str(identity.get("version", "")).strip()
    manifest = ROOT / "release" / f"v{version}" / "certification-manifest.json"
    if not manifest.is_file():
        raise SystemExit(f"canonical G12 manifest missing: {manifest.relative_to(ROOT)}")

    fetch_tags()
    tag = f"v{version}-stable"
    payload = decide(
        identity=identity,
        candidate_sha=args.candidate_sha.strip(),
        source_head_ref=args.source_head_ref.strip(),
        existing_tag_sha=resolve_tag_commit(tag),
        latest_tag=latest_stable_tag(),
        release=release_metadata(args.repo, tag),
    )
    text = json.dumps(payload, indent=2, sort_keys=True) + "\n"
    print(text, end="")
    if args.json_out:
        Path(args.json_out).write_text(text, encoding="utf-8")
    return 0


if __name__ == "__main__":
    try:
        raise SystemExit(main())
    except (ValueError, RuntimeError) as exc:
        print(f"DE.PULSE G11 publication preflight: FAIL - {exc}", file=sys.stderr)
        raise SystemExit(2)
