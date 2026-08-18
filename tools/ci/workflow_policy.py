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
    required = (
        "endsWith(github.ref, '-release-certification')",
        "endsWith(github.ref, '-stable-promotion')",
        "github.event.sender.login == github.repository_owner",
        '"v${release_line}-release-certification") publish=false',
        '"v${release_line}-stable-promotion") publish=true',
    )
    missing = [fragment for fragment in required if fragment not in ci_fast]
    forbidden = "github.event.head_commit.author.username"
    if missing or forbidden in ci_fast:
        print("DE.PULSE workflow policy: FAIL", file=sys.stderr)
        if missing:
            print("release dispatcher contract missing: " + ", ".join(missing), file=sys.stderr)
        if forbidden in ci_fast:
            print("release dispatcher uses unreliable head_commit author authorization", file=sys.stderr)
        return 1
    print("release dispatcher authorization/publish boundary: PASS")
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
    if run_gate(root, "dependency_readiness_gate.py", "dependency/provider readiness contract") != 0:
        return 1
    if run_gate(root, "ai_continuous_eval_gate.py", "AI continuous eval/rights contract") != 0:
        return 1

    print("DE.PULSE workflow policy: PASS")
    print("active workflows: " + ", ".join(present))
    print("dependency/provider readiness: PASS")
    print("AI continuous eval/rights: PASS")
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
