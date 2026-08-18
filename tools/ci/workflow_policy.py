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

    dependency_gate = root / "dependency_readiness_gate.py"
    if not dependency_gate.is_file():
        print("DE.PULSE workflow policy: FAIL", file=sys.stderr)
        print("missing dependency_readiness_gate.py", file=sys.stderr)
        return 1
    dependency_result = subprocess.run([sys.executable, str(dependency_gate)], cwd=root, check=False)
    if dependency_result.returncode != 0:
        print("DE.PULSE workflow policy: FAIL", file=sys.stderr)
        print("dependency/provider readiness contract failed", file=sys.stderr)
        return dependency_result.returncode

    print("DE.PULSE workflow policy: PASS")
    print("active workflows: " + ", ".join(present))
    print("dependency/provider readiness: PASS")
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
