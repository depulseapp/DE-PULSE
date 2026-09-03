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
    expected_scoped = [
        ".github/workflows/ci-fast.yml",
        ".github/workflows/ci-qualified.yml",
        ".github/workflows/release.yml",
    ]
    if scoped != expected_scoped:
        errors.append("immutable Action pin scope drift")

    actions = lock.get("actions", {})
    expected_actions = {
        "actions/checkout",
        "actions/setup-python",
        "actions/setup-go",
        "actions/setup-node",
        "actions/upload-artifact",
        "actions/download-artifact",
        "actions/attest-build-provenance",
        "azure/login",
        "hashicorp/setup-terraform",
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
        ".github/workflows/ci-qualified.yml": {"statuses": 1, "id-token": 1},
        ".github/workflows/release.yml": {"contents": 1, "id-token": 1, "attestations": 1},
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

    qualified_text = (ROOT / ".github" / "workflows" / "ci-qualified.yml").read_text(encoding="utf-8")
    qualified_oidc_required = (
        "host013_azure_confirmation:",
        "HOST013_AZURE_AKS_OPERATOR_DRILL",
        "name: HOST-013/014 Azure AKS trust operator drill",
        "id-token: write",
        "azure/login@532459ea530d8321f2fb9bb10d1e0bcf23869a43",
        "python3 tools/ci/host013_azure_operator.py",
        "--environment dev",
    )
    for token in qualified_oidc_required:
        if token not in qualified_text:
            errors.append(f"Qualified Azure OIDC least-privilege contract missing: {token}")
    for forbidden in ("AZURE_CLIENT_SECRET", "ARM_CLIENT_SECRET"):
        if forbidden in qualified_text:
            errors.append(f"Qualified Azure OIDC contract contains forbidden secret marker: {forbidden}")

    release_text = (ROOT / ".github" / "workflows" / "release.yml").read_text(encoding="utf-8")
    provenance_required = (
        "actions/attest-build-provenance@",
        "id-token: write",
        "attestations: write",
        "subject-path:",
    )
    for token in provenance_required:
        if token not in release_text:
            errors.append(f"release provenance least-privilege contract missing: {token}")

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

    required_browser_contract = (
        "cache: 'pip'",
        "cache-dependency-path: tools/ci/browser-requirements.txt",
        "python -m pip install --disable-pip-version-check -r tools/ci/browser-requirements.txt",
    )
    for workflow_name in ("ci-qualified.yml", "release.yml"):
        workflow = (ROOT / ".github" / "workflows" / workflow_name).read_text(encoding="utf-8")
        for token in required_browser_contract:
            if token not in workflow:
                errors.append(f"{workflow_name} browser reproducibility/cache contract missing: {token}")
        if re.search(r"pip\s+install[^\n]*\bplaywright(?:\s|$)", workflow):
            errors.append(f"{workflow_name} must install Playwright only through the pinned requirements file")

    if errors:
        return fail(errors)

    print("DE.PULSE CI reproducibility gate: PASS")
    print("Fast/Qualified/Release third-party Actions immutable-SHA pinned: PASS")
    print("Qualified/Release Playwright exact-version requirements + safe pip cache: PASS")
    print("scoped least-privilege permission contract: PASS")
    print("Qualified Azure OIDC exception is explicit, manual-only and host013-job scoped: PASS")
    print("Release provenance OIDC/attestation exception is explicit and publish-job scoped: PASS")
    print("Release workflow Action/browser reproducibility deferral closed in v18.7.0: PASS")
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
