#!/usr/bin/env python3
from __future__ import annotations

import json
from pathlib import Path
import re
import sys

ROOT = Path(__file__).resolve().parents[2]
MANIFEST = ROOT / "governance" / "toolchain-manifest.json"
WORKFLOWS = (
    ROOT / ".github" / "workflows" / "ci-fast.yml",
    ROOT / ".github" / "workflows" / "ci-qualified.yml",
    ROOT / ".github" / "workflows" / "release.yml",
)


def fail(errors: list[str]) -> int:
    print("DE.PULSE toolchain contract: FAIL", file=sys.stderr)
    for error in errors:
        print(f" - {error}", file=sys.stderr)
    return 1


def main() -> int:
    errors: list[str] = []
    manifest = json.loads(MANIFEST.read_text(encoding="utf-8"))
    if manifest.get("schema") != "DE.PULSE-TOOLCHAIN-MANIFEST-1":
        errors.append("toolchain manifest schema mismatch")

    go = str(manifest.get("go", {}).get("version", ""))
    node = str(manifest.get("node", {}).get("version", ""))
    python = str(manifest.get("python", {}).get("version", ""))
    playwright = str(manifest.get("playwright", {}).get("version", ""))
    for label, value, pattern in (
        ("Go", go, r"\d+\.\d+\.\d+"),
        ("Node", node, r"\d+\.\d+\.\d+"),
        ("Python", python, r"\d+\.\d+\.\d+"),
        ("Playwright", playwright, r"\d+\.\d+\.\d+"),
    ):
        if not re.fullmatch(pattern, value):
            errors.append(f"{label} must be exact patch version, got {value!r}")

    lock = ROOT / str(manifest.get("playwright", {}).get("lockFile", ""))
    if not lock.is_file() or f"playwright=={playwright}" not in lock.read_text(encoding="utf-8"):
        errors.append("Playwright lock does not match canonical manifest")

    runners = manifest.get("runnerImages", {})
    required_runners = {
        str(runners.get("linux", "")),
        str(runners.get("macos", "")),
        str(runners.get("windows", "")),
    }
    for workflow in WORKFLOWS:
        text = workflow.read_text(encoding="utf-8")
        for token in (
            f"GO_VERSION: '{go}'",
            f"NODE_VERSION: '{node}'",
            f"PYTHON_VERSION: '{python}'",
        ):
            if token not in text:
                errors.append(f"{workflow.name} does not use canonical exact toolchain token {token}")
        for runner in required_runners:
            if runner and runner not in text and workflow.name != "ci-fast.yml":
                errors.append(f"{workflow.name} missing governed runner image {runner}")

    executor = (ROOT / "tools" / "release" / "run_full_certification.py").read_text(encoding="utf-8")
    for token in (
        "resolvedToolchain",
        "go version",
        "node --version",
        "python3 --version",
        "importlib.metadata",
        "RUNNER_OS",
        "ImageOS",
    ):
        if token not in executor:
            errors.append(f"canonical G12 executor does not record resolved toolchain identity: {token}")

    if errors:
        return fail(errors)
    print("DE.PULSE toolchain contract: PASS")
    print(f"Go={go} Node={node} Python={python} Playwright={playwright}")
    print("exact patch selectors in Fast/Qualified/Release: PASS")
    print("canonical runner image contract: PASS")
    print("resolved release provenance recording: PASS")
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
