#!/usr/bin/env python3
from __future__ import annotations

import argparse
import json
from pathlib import Path
import subprocess

PROCESS_ONLY_PREFIXES = (
    ".github/workflows/",
    "tools/ci/",
    "tools/release/",
    "adaptive-governance/",
    "governance/",
    "handoff/",
    ".depulse-certification/resume/",
)
PROCESS_ONLY_EXACT = {
    "source_fingerprint.py",
    "README.md",
    "AGENTS.md",
    "CLAUDE.md",
}


def git(*args: str, check: bool = True) -> subprocess.CompletedProcess[str]:
    return subprocess.run(("git", *args), check=check, text=True, capture_output=True)


def resolve_base(base: str, head: str) -> str:
    candidate = base.strip()
    if candidate and set(candidate) != {"0"} and git("cat-file", "-e", f"{candidate}^{{commit}}", check=False).returncode == 0:
        return candidate
    parent = git("rev-parse", f"{head}^", check=False)
    if parent.returncode == 0:
        return parent.stdout.strip()
    return head


def is_process_only(path: str) -> bool:
    return path in PROCESS_ONLY_EXACT or path.startswith(PROCESS_ONLY_PREFIXES)


def main() -> int:
    parser = argparse.ArgumentParser(description="Plan DE.PULSE CI lanes from a deterministic Git diff")
    parser.add_argument("--base", default="")
    parser.add_argument("--head", default="HEAD")
    parser.add_argument("--github-output")
    parser.add_argument("--json-out")
    args = parser.parse_args()

    head = git("rev-parse", args.head).stdout.strip()
    base = resolve_base(args.base, head)
    raw = git("diff", "--name-only", base, head).stdout
    changed = sorted({line.strip() for line in raw.splitlines() if line.strip()})

    process_only = bool(changed) and all(is_process_only(path) for path in changed)
    go_required = any(path.endswith(".go") or path in {"go.mod", "go.sum"} for path in changed)
    node_required = any(path.endswith((".js", ".mjs", ".cjs")) for path in changed)
    qualified_lane = "ci-harness" if process_only else "full"

    plan = {
        "schema": "DE.PULSE-CI-IMPACT-PLAN-1",
        "baseSha": base,
        "headSha": head,
        "changedPaths": changed,
        "processOnly": process_only,
        "goRequired": go_required,
        "nodeRequired": node_required,
        "qualifiedLane": qualified_lane,
        "reason": (
            "Only canonical CI/release/governance/handoff tooling changed; run harness + portability only."
            if process_only
            else "Product, test, dependency, release identity, or other non-process content changed; run full qualified coverage."
        ),
    }
    text = json.dumps(plan, indent=2, sort_keys=True) + "\n"
    print(text, end="")

    if args.json_out:
        out = Path(args.json_out)
        out.parent.mkdir(parents=True, exist_ok=True)
        out.write_text(text)
    if args.github_output:
        out = Path(args.github_output)
        with out.open("a", encoding="utf-8") as f:
            f.write(f"qualified_lane={qualified_lane}\n")
            f.write(f"go_required={str(go_required).lower()}\n")
            f.write(f"node_required={str(node_required).lower()}\n")
            f.write(f"process_only={str(process_only).lower()}\n")
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
