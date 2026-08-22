#!/usr/bin/env python3
"""Digest-preserving Stable release asset publication for DE.PULSE.

Existing Stable asset names are never overwritten. If an asset already exists,
its downloaded bytes must hash exactly to the local certified artifact; exact
matches are idempotent reuse and any mismatch is a release-integrity failure.
"""
from __future__ import annotations

import argparse
import hashlib
import json
import os
from pathlib import Path
import subprocess
import sys
import tempfile


def sha256(path: Path) -> str:
    h = hashlib.sha256()
    with path.open("rb") as fh:
        for block in iter(lambda: fh.read(1024 * 1024), b""):
            h.update(block)
    return h.hexdigest()


def decide_asset(existing_digest: str | None, local_digest: str) -> str:
    if not local_digest:
        raise ValueError("local asset digest is required")
    if existing_digest is None:
        return "UPLOAD"
    if existing_digest == local_digest:
        return "REUSE"
    raise ValueError(
        f"Stable asset digest mismatch: existing={existing_digest} local={local_digest}; overwrite is forbidden"
    )


def run(argv: list[str], *, check: bool = True) -> subprocess.CompletedProcess[str]:
    return subprocess.run(argv, text=True, capture_output=True, check=check)


def release_assets(repo: str, tag: str) -> dict[str, dict[str, object]]:
    result = run(
        ["gh", "release", "view", tag, "--repo", repo, "--json", "assets"],
        check=False,
    )
    if result.returncode != 0:
        raise RuntimeError(f"Stable release {tag} does not exist: {result.stderr.strip()}")
    payload = json.loads(result.stdout)
    rows: dict[str, dict[str, object]] = {}
    for item in payload.get("assets", []) or []:
        name = str(item.get("name", "")).strip()
        if not name:
            continue
        if name in rows:
            raise RuntimeError(f"duplicate existing Stable asset name: {name}")
        rows[name] = item
    return rows


def download_existing(repo: str, tag: str, name: str, directory: Path) -> Path:
    result = run(
        [
            "gh",
            "release",
            "download",
            tag,
            "--repo",
            repo,
            "--pattern",
            name,
            "--dir",
            str(directory),
        ],
        check=False,
    )
    if result.returncode != 0:
        raise RuntimeError(f"failed to download existing Stable asset {name}: {result.stderr.strip()}")
    path = directory / name
    if not path.is_file():
        raise RuntimeError(f"downloaded Stable asset missing after gh release download: {name}")
    return path


def publish(repo: str, tag: str, assets: list[Path]) -> dict[str, object]:
    if not assets:
        raise ValueError("at least one asset is required")
    existing = release_assets(repo, tag)
    seen: set[str] = set()
    results: list[dict[str, str]] = []

    with tempfile.TemporaryDirectory(prefix="depulse-stable-assets-") as tmp:
        tmpdir = Path(tmp)
        for asset in assets:
            path = asset.resolve()
            if not path.is_file():
                raise ValueError(f"local certified asset missing: {asset}")
            name = path.name
            if name in seen:
                raise ValueError(f"duplicate local Stable asset name: {name}")
            seen.add(name)
            local_digest = sha256(path)

            if name in existing:
                existing_path = download_existing(repo, tag, name, tmpdir)
                existing_digest = sha256(existing_path)
                action = decide_asset(existing_digest, local_digest)
                print(f"REUSE: {name} sha256={local_digest}")
            else:
                action = decide_asset(None, local_digest)
                result = run(["gh", "release", "upload", tag, "--repo", repo, str(path)], check=False)
                if result.returncode != 0:
                    raise RuntimeError(f"failed to upload Stable asset {name}: {result.stderr.strip()}")
                print(f"UPLOAD: {name} sha256={local_digest}")

            results.append({"name": name, "sha256": local_digest, "action": action})

    return {
        "schema": "DE.PULSE-STABLE-ASSET-IMMUTABILITY-1",
        "tag": tag,
        "assets": results,
        "differingExistingBytesOverwrite": "FORBIDDEN",
    }


def self_test() -> int:
    digest_a = hashlib.sha256(b"same").hexdigest()
    digest_b = hashlib.sha256(b"different").hexdigest()
    assert decide_asset(None, digest_a) == "UPLOAD"
    assert decide_asset(digest_a, digest_a) == "REUSE"
    try:
        decide_asset(digest_b, digest_a)
    except ValueError:
        pass
    else:
        raise AssertionError("differing Stable bytes must fail closed")
    print("DE.PULSE Stable asset immutability self-test: PASS")
    print("absent asset -> upload: PASS")
    print("same digest -> idempotent reuse/no-op: PASS")
    print("different digest -> fail, never overwrite: PASS")
    return 0


def main() -> int:
    parser = argparse.ArgumentParser(description="Publish immutable DE.PULSE Stable release assets")
    parser.add_argument("--repo", default=os.environ.get("GITHUB_REPOSITORY", ""))
    parser.add_argument("--tag", default="")
    parser.add_argument("--asset", action="append", default=[])
    parser.add_argument("--json-out", default="")
    parser.add_argument("--self-test", action="store_true")
    args = parser.parse_args()

    if args.self_test:
        return self_test()
    if not args.repo or not args.tag:
        raise SystemExit("--repo/GITHUB_REPOSITORY and --tag are required")

    payload = publish(args.repo, args.tag, [Path(item) for item in args.asset])
    text = json.dumps(payload, indent=2, sort_keys=True) + "\n"
    if args.json_out:
        Path(args.json_out).write_text(text, encoding="utf-8")
    print(text, end="")
    return 0


if __name__ == "__main__":
    try:
        raise SystemExit(main())
    except (ValueError, RuntimeError) as exc:
        print(f"DE.PULSE Stable asset immutability: FAIL - {exc}", file=sys.stderr)
        raise SystemExit(2)
