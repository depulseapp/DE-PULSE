#!/usr/bin/env python3
from __future__ import annotations

import argparse
import json
import os
from pathlib import Path
import shutil
import subprocess

DEPRECATED_TEMP_MARKERS = (
    "-release-certification",
    "-stable-promotion",
    "-certification-trigger",
    "-cert-trigger",
    "-promotion-trigger",
    "-dispatch",
    "-retry",
    "-fallback",
    "-trigger",
)


def run(*args: str, check: bool = True) -> subprocess.CompletedProcess[str]:
    return subprocess.run(args, check=check, text=True, capture_output=True)


def is_deprecated_release_temp(name: str) -> bool:
    return name.startswith("v") and any(marker in name for marker in DEPRECATED_TEMP_MARKERS)


def open_pr_heads() -> set[str] | None:
    repo = os.environ.get("GITHUB_REPOSITORY", "").strip()
    if not repo or not shutil.which("gh"):
        return None
    result = run(
        "gh", "pr", "list", "--repo", repo, "--state", "open", "--limit", "200",
        "--json", "headRefName", "--jq", ".[ ].headRefName".replace("[ ]", "[]"),
        check=False,
    )
    if result.returncode != 0:
        return None
    return {line.strip() for line in result.stdout.splitlines() if line.strip()}


def main() -> int:
    p = argparse.ArgumentParser(description="Audit/delete merged and deprecated DE.PULSE release-temp branches")
    p.add_argument("--apply", action="store_true")
    p.add_argument("--json-out")
    p.add_argument("--remote", default="origin")
    p.add_argument("--base", default="main")
    args = p.parse_args()

    run("git", "fetch", "--prune", args.remote, f"+refs/heads/*:refs/remotes/{args.remote}/*")
    base_ref = f"refs/remotes/{args.remote}/{args.base}"
    base_sha = run("git", "rev-parse", base_ref).stdout.strip()
    raw = run(
        "git", "for-each-ref", f"refs/remotes/{args.remote}",
        "--format=%(refname:short)\t%(objectname)"
    ).stdout

    protected = {args.base, "HEAD"}
    open_heads = open_pr_heads()
    merged: list[dict[str, str]] = []
    deprecated_temp: list[dict[str, str]] = []
    retained: list[dict[str, str]] = []

    for line in raw.splitlines():
        if not line.strip():
            continue
        ref, sha = line.split("\t", 1)
        prefix = f"{args.remote}/"
        name = ref[len(prefix):] if ref.startswith(prefix) else ref
        if name in protected or name.startswith("HEAD ->"):
            continue

        item = {"branch": name, "sha": sha}
        if open_heads is not None and name in open_heads:
            item["reason"] = "OPEN_PR"
            retained.append(item)
            continue

        ancestor = run("git", "merge-base", "--is-ancestor", sha, base_sha, check=False)
        if ancestor.returncode == 0:
            item["reason"] = "FULLY_MERGED"
            merged.append(item)
        elif is_deprecated_release_temp(name) and open_heads is not None:
            item["reason"] = "DEPRECATED_RELEASE_TEMP_NO_OPEN_PR"
            deprecated_temp.append(item)
        else:
            item["reason"] = "UNIQUE_OR_ACTIVE"
            retained.append(item)

    deleted_merged: list[dict[str, str]] = []
    deleted_temp: list[dict[str, str]] = []
    if args.apply:
        for bucket, deleted in ((merged, deleted_merged), (deprecated_temp, deleted_temp)):
            for item in bucket:
                name = item["branch"]
                result = run("git", "push", args.remote, "--delete", name, check=False)
                if result.returncode != 0:
                    raise SystemExit(f"failed deleting branch {name}: {result.stderr.strip()}")
                deleted.append(item)

    report = {
        "schema": "DE.PULSE-BRANCH-HYGIENE-2",
        "base": args.base,
        "baseSha": base_sha,
        "mode": "APPLY" if args.apply else "DRY_RUN",
        "openPrHeadsResolved": open_heads is not None,
        "mergedCandidates": merged,
        "deprecatedReleaseTempCandidates": deprecated_temp,
        "deletedMerged": deleted_merged,
        "deletedDeprecatedReleaseTemp": deleted_temp,
        "retained": retained,
        "policy": (
            "Never delete a branch with an open PR. Delete branches already contained in main. "
            "After open-PR resolution succeeds, also delete closed/orphaned versioned trigger/retry/dispatch/release-certification/stable-promotion branches."
        ),
    }
    text = json.dumps(report, indent=2, sort_keys=True) + "\n"
    print(text, end="")
    if args.json_out:
        path = Path(args.json_out)
        path.parent.mkdir(parents=True, exist_ok=True)
        path.write_text(text)
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
