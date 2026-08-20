#!/usr/bin/env python3
from __future__ import annotations

import json
import os
from pathlib import Path
import subprocess
import sys

ROOT = Path(__file__).resolve().parents[2]


def load_json(path: Path) -> dict:
    try:
        value = json.loads(path.read_text(encoding="utf-8"))
        return value if isinstance(value, dict) else {}
    except Exception:
        return {}


def branch_name() -> str:
    env = (os.environ.get("GITHUB_HEAD_REF") or os.environ.get("GITHUB_REF_NAME") or "").strip()
    if env:
        return env.removeprefix("refs/heads/")
    proc = subprocess.run(["git", "branch", "--show-current"], cwd=ROOT, text=True, capture_output=True, check=False)
    return proc.stdout.strip()


def fail(items: list[str]) -> int:
    print("DE.PULSE post-Stable continuity: FAIL", file=sys.stderr)
    for item in items:
        print(f" - {item}", file=sys.stderr)
    return 1


def main() -> int:
    errors: list[str] = []
    identity = load_json(ROOT / "release_identity.json")
    build = load_json(ROOT / ".depulse-certification/resume/build-checkpoint.json")
    evidence = load_json(ROOT / ".depulse-certification/resume/release-evidence-checkpoint.json")

    identity_version = str(identity.get("version", "")).lstrip("v")
    checkpoint_version = str(build.get("release", "")).lstrip("v")
    evidence_version = str(evidence.get("release", "")).lstrip("v")
    branch = branch_name()

    if not identity_version or not checkpoint_version or not evidence_version:
        errors.append("release identity/build checkpoint/release evidence version missing")
        return fail(errors)

    if checkpoint_version != evidence_version:
        errors.append(f"build/evidence Stable mismatch: {checkpoint_version} != {evidence_version}")

    if identity_version != checkpoint_version:
        expected = f"v{identity_version}-development"
        if branch in {"main", "master"}:
            errors.append(
                f"default branch carries STABLE identity v{identity_version} while durable Stable checkpoint is v{checkpoint_version}; post-Stable continuity reconciliation is required"
            )
        elif branch != expected:
            errors.append(
                f"identity/checkpoint differ outside the expected in-flight branch {expected}: current branch={branch or '<detached>'}"
            )
    else:
        stable_tag = f"v{identity_version}-stable"
        manifest = ROOT / "release" / f"v{identity_version}" / "stable-evidence-manifest.json"
        if not manifest.is_file():
            errors.append(f"Stable evidence manifest missing: {manifest.relative_to(ROOT)}")

        handoff = (ROOT / "handoff/CURRENT.md").read_text(encoding="utf-8", errors="ignore")
        if f"**Certified Stable:** `{stable_tag}`" not in handoff:
            errors.append("handoff/CURRENT.md does not identify the aligned Stable tag")

        current_docs = [
            ROOT / "adaptive-governance/CURRENT_ADAPTIVE_ROADMAP.md",
            ROOT / "adaptive-governance/CURRENT_ADAPTIVE_BUILD_PLAN.md",
            ROOT / "adaptive-governance/CURRENT_ADAPTIVE_BUILD_PROCESS.md",
            ROOT / "adaptive-governance/CURRENT_ADAPTIVE_DELIVERY_PROCESS.md",
        ]
        for path in current_docs:
            text = path.read_text(encoding="utf-8", errors="ignore") if path.is_file() else ""
            if stable_tag not in text:
                errors.append(f"{path.relative_to(ROOT)} is not aligned to {stable_tag}")

    if errors:
        return fail(errors)

    mode = "STABLE_ALIGNED" if identity_version == checkpoint_version else "IN_FLIGHT_CANDIDATE"
    print(
        "DE.PULSE post-Stable continuity: PASS · "
        f"mode={mode} · branch={branch or '<detached>'} · identity=v{identity_version} · durable Stable=v{checkpoint_version}"
    )
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
