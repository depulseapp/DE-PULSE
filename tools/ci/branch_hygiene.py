#!/usr/bin/env python3
from __future__ import annotations

import argparse
import json
import os
from pathlib import Path
import re
import shutil
import subprocess

DEPRECATED_TEMP_MARKERS = (
    "-release-certification",
    "-stable-certification",
    "-stable-promotion",
    "-stable-tag-tooling",
    "-native-certification",
    "-certification-trigger",
    "-cert-trigger",
    "-promotion-trigger",
    "-promotion-hardening",
    "-g10-recovery",
    "-dispatch",
    "-retry",
    "-fallback",
    "-trigger",
)
VERSION_PREFIX = re.compile(r"^(v\d+\.\d+(?:\.\d+)?)")


def run(*args: str, check: bool = True) -> subprocess.CompletedProcess[str]:
    return subprocess.run(args, check=check, text=True, capture_output=True)


def is_deprecated_release_temp(name: str) -> bool:
    return name.startswith("v") and any(marker in name for marker in DEPRECATED_TEMP_MARKERS)


def pr_heads(state: str) -> set[str] | None:
    repo = os.environ.get("GITHUB_REPOSITORY", "").strip()
    if not repo or not shutil.which("gh"):
        return None
    result = run(
        "gh", "pr", "list", "--repo", repo, "--state", state, "--limit", "200",
        "--json", "headRefName", "--jq", ".[].headRefName",
        check=False,
    )
    if result.returncode != 0:
        return None
    return {line.strip() for line in result.stdout.splitlines() if line.strip()}


def stable_line_closed(name: str) -> bool:
    match = VERSION_PREFIX.match(name)
    if not match:
        return False
    prefix = match.group(1)
    if prefix.count(".") >= 2:
        pattern = f"{prefix}-stable"
    else:
        pattern = f"{prefix}.*-stable"
    return bool(run("git", "tag", "--list", pattern).stdout.strip())


def main() -> int:
    p = argparse.ArgumentParser(description="Audit/delete merged and obsolete DE.PULSE branches")
    p.add_argument("--apply", action="store_true")
    p.add_argument("--json-out")
    p.add_argument("--remote", default="origin")
    p.add_argument("--base", default="main")
    args = p.parse_args()

    run("git", "fetch", "--prune", "--tags", args.remote, f"+refs/heads/*:refs/remotes/{args.remote}/*")
    base_ref = f"refs/remotes/{args.remote}/{args.base}"
    base_sha = run("git", "rev-parse", base_ref).stdout.strip()
    raw = run(
        "git", "for-each-ref", f"refs/remotes/{args.remote}",
        "--format=%(refname:short)\t%(objectname)"
    ).stdout

    protected = {args.base, "HEAD"}
    open_heads = pr_heads("open")
    merged_heads = pr_heads("merged")
    merged_ancestor: list[dict[str, str]] = []
    merged_pr: list[dict[str, str]] = []
    closed_stable_line: list[dict[str, str]] = []
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

        if merged_heads is not None and name in merged_heads:
            item["reason"] = "MERGED_PR_HEAD"
            merged_pr.append(item)
            continue

        ancestor = run("git", "merge-base", "--is-ancestor", sha, base_sha, check=False)
        if ancestor.returncode == 0:
            item["reason"] = "FULLY_MERGED"
            merged_ancestor.append(item)
        elif open_heads is not None and stable_line_closed(name):
            item["reason"] = "STABLE_LINE_ALREADY_PUBLISHED"
            closed_stable_line.append(item)
        elif is_deprecated_release_temp(name) and open_heads is not None:
            item["reason"] = "DEPRECATED_RELEASE_TEMP_NO_OPEN_PR"
            deprecated_temp.append(item)
        else:
            item["reason"] = "UNIQUE_OR_ACTIVE"
            retained.append(item)

    deleted: list[dict[str, str]] = []
    candidates = merged_pr + merged_ancestor + closed_stable_line + deprecated_temp
    if args.apply:
        seen: set[str] = set()
        for item in candidates:
            name = item["branch"]
            if name in seen:
                continue
            seen.add(name)
            result = run("git", "push", args.remote, "--delete", name, check=False)
            if result.returncode != 0:
                raise SystemExit(f"failed deleting branch {name}: {result.stderr.strip()}")
            deleted.append(item)

    report = {
        "schema": "DE.PULSE-BRANCH-HYGIENE-3",
        "base": args.base,
        "baseSha": base_sha,
        "mode": "APPLY" if args.apply else "DRY_RUN",
        "openPrHeadsResolved": open_heads is not None,
        "mergedPrHeadsResolved": merged_heads is not None,
        "mergedPrHeadCandidates": merged_pr,
        "mergedAncestorCandidates": merged_ancestor,
        "closedStableLineCandidates": closed_stable_line,
        "deprecatedReleaseTempCandidates": deprecated_temp,
        "deleted": deleted,
        "retained": retained,
        "policy": (
            "Never delete a branch with an open PR. Delete merged PR heads (including squash-merged heads), "
            "branches already contained in main, versioned branches whose Stable line is already published, "
            "and closed/orphaned release-temp branches. Fail conservative when PR state cannot be resolved."
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
