#!/usr/bin/env python3
from __future__ import annotations

import json
from pathlib import Path
import re
import sys

ROOT = Path(__file__).resolve().parents[2]
LOCK_PATH = ROOT / "tools" / "ci" / "ci_dependency_lock.json"
USE_RE = re.compile(r"^\s*(?:-\s+)?uses:\s+([^@\s]+)@([0-9a-f]{40})(?:\s+#\s*(\S+))?\s*$")
USE_DIRECTIVE_RE = re.compile(r"^\s*(?:-\s+)?uses:\s+")
WRITE_RE = re.compile(r"^\s+([a-zA-Z0-9_-]+):\s*write\s*$")


def fail(errors: list[str]) -> int:
    print("DE.PULSE CI reproducibility gate: FAIL", file=sys.stderr)
    for error in errors:
        print(f" - {error}", file=sys.stderr)
    return 1


def main() -> int:
    errors: list[str] = []
    try:
        lock = json.loads(LOCK_PATH.read_text(encoding="utf-8"))
    except Exception as exc:
        return fail([f"dependency lock unreadable: {exc}"])

    if lock.get("schema") != "DE.PULSE-CI-DEPENDENCY-LOCK-1":
        errors.append("unexpected dependency lock schema")

    scoped = lock.get("scope", {}).get("immutable_action_pin_workflows", [])
    if scoped != [".github/workflows/ci-fast.yml", ".github/workflows/ci-qualified.yml"]:
        errors.append("immutable Action pin scope drift")

    actions = lock.get("actions", {})
    expected_actions = {
        "actions/checkout",
        "actions/setup-python",
        "actions/setup-go",
        "actions/setup-node",
        "actions/upload-artifact",
    }
    if set(actions) != expected_actions:
        errors.append(f"dependency lock Action set drift: {sorted(actions)}")

    for name, meta in actions.items():
        sha = str(meta.get("sha", ""))
        label = str(meta.get("label", ""))
        if not re.fullmatch(r"[0-9a-f]{40}", sha):
            errors.append(f"{name} is not pinned to a 40-hex commit SHA")
        if not label:
            errors.append(f"{name} version label missing")

    allowed_write = {
        ".github/workflows/ci-fast.yml": {"statuses": 1, "contents": 1},
        ".github/workflows/ci-qualified.yml": {"statuses": 1},
    }
    forbidden_permissions = set(lock.get("permission_policy", {}).get("forbidden_in_scoped_workflows", []))

    for relative in scoped:
        path = ROOT / relative
        if not path.is_file():
            errors.append(f"missing scoped workflow: {relative}")
            continue
        text = path.read_text(encoding="utf-8")
        lines = text.splitlines()
        if "permissions:\n  contents: read" not in text:
            errors.append(f"{relative} must keep top-level contents: read")
        if "\t" in text:
            errors.append(f"{relative} contains tab indentation")
        if any(line.rstrip() != line for line in lines):
            errors.append(f"{relative} contains trailing whitespace")

        seen_actions: set[str] = set()
        for line_no, line in enumerate(lines, 1):
            if not USE_DIRECTIVE_RE.match(line):
                continue
            target = line.split("uses:", 1)[1].strip()
            if target.startswith("./"):
                continue
            match = USE_RE.match(line)
            if not match:
                errors.append(f"{relative}:{line_no} external Action is not immutable-SHA pinned: {target}")
                continue
            action, sha, label = match.groups()
            seen_actions.add(action)
            meta = actions.get(action)
            if not meta:
                errors.append(f"{relative}:{line_no} Action missing from dependency lock: {action}")
                continue
            if sha != meta.get("sha"):
                errors.append(f"{relative}:{line_no} {action} SHA differs from dependency lock")
            if label != meta.get("label"):
                errors.append(f"{relative}:{line_no} {action} version comment must be # {meta.get('label')}")

        if "actions/checkout" not in seen_actions or "actions/setup-python" not in seen_actions:
            errors.append(f"{relative} missing required pinned checkout/setup-python Action")

        writes: dict[str, int] = {}
        for line in lines:
            match = WRITE_RE.match(line)
            if match:
                key = match.group(1)
                writes[key] = writes.get(key, 0) + 1
        if writes != allowed_write[relative]:
            errors.append(f"{relative} write-permission set drift: {writes}; expected {allowed_write[relative]}")
        for forbidden in forbidden_permissions:
            if forbidden in text:
                errors.append(f"{relative} contains forbidden permission: {forbidden}")

    browser = lock.get("browser", {})
    requirements_path = ROOT / str(browser.get("requirements_file", ""))
    if not requirements_path.is_file():
        errors.append("browser requirements file missing")
    else:
        requirements = [
            line.strip()
            for line in requirements_path.read_text(encoding="utf-8").splitlines()
            if line.strip() and not line.lstrip().startswith("#")
        ]
        expected = [f"playwright=={browser.get('playwright')}"]
        if requirements != expected:
            errors.append(f"browser requirements must be exactly {expected}; got {requirements}")

    qualified = (ROOT / ".github" / "workflows" / "ci-qualified.yml").read_text(encoding="utf-8")
    required_browser_contract = (
        "cache: 'pip'",
        "cache-dependency-path: tools/ci/browser-requirements.txt",
        "python -m pip install --disable-pip-version-check -r tools/ci/browser-requirements.txt",
    )
    for token in required_browser_contract:
        if token not in qualified:
            errors.append(f"Qualified browser reproducibility/cache contract missing: {token}")
    if re.search(r"pip\s+install[^\n]*\bplaywright(?:\s|$)", qualified):
        errors.append("Qualified workflow must install Playwright only through the pinned requirements file")

    deferred = lock.get("scope", {}).get("deferred_release_workflow", {})
    if deferred.get("path") != ".github/workflows/release.yml" or not deferred.get("reason"):
        errors.append("release-workflow pin deferral must remain explicit and reasoned")

    if errors:
        return fail(errors)

    print("DE.PULSE CI reproducibility gate: PASS")
    print("Fast/Qualified third-party Actions immutable-SHA pinned: PASS")
    print("Playwright exact-version requirements + safe pip cache: PASS")
    print("scoped least-privilege permission contract: PASS")
    print("release.yml pinning explicitly deferred to next release-capable product slice: PASS")
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
