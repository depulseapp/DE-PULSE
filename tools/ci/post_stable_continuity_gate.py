#!/usr/bin/env python3
"""Fail-closed post-Stable continuity and branch/state alignment gate."""
from __future__ import annotations

import json
import os
from pathlib import Path
import re
import subprocess
import sys

ROOT = Path(__file__).resolve().parents[2]
CLOSED = {"COMPLETE", "COMPLETED", "CLOSED", "DELIVERED", "VERIFIED"}
ACTIVE = {"ACTIVE", "IN_PROGRESS"}


def load_json(path: Path) -> dict:
    try:
        value = json.loads(path.read_text(encoding="utf-8"))
        return value if isinstance(value, dict) else {}
    except Exception:
        return {}


def git(*args: str) -> subprocess.CompletedProcess[str]:
    return subprocess.run(["git", *args], cwd=ROOT, text=True, capture_output=True, check=False)


def branch_name() -> str:
    env = (os.environ.get("GITHUB_HEAD_REF") or os.environ.get("GITHUB_REF_NAME") or "").strip()
    return env.removeprefix("refs/heads/") if env else git("branch", "--show-current").stdout.strip()


def clean_version(value: object) -> str:
    raw = str(value or "").strip()
    if raw.startswith("v"):
        raw = raw[1:]
    if raw.endswith("-stable"):
        raw = raw[:-7]
    return raw


def version_tuple(value: object) -> tuple[int, int, int] | None:
    match = re.fullmatch(r"(\d+)\.(\d+)\.(\d+)", clean_version(value))
    return tuple(int(x) for x in match.groups()) if match else None


def stable_tag_for(version: object) -> str:
    policy = load_json(ROOT / "governance" / "versioning-policy.json")
    current = version_tuple(version)
    cutover = version_tuple(policy.get("effectiveAfterProductVersion"))
    if not current or not cutover:
        return f"v{clean_version(version)}"
    pattern = policy.get("legacyStableTagPattern") if current <= cutover else policy.get("futureStableTagPattern")
    return str(pattern).format(productVersion=clean_version(version))


def fail(items: list[str]) -> int:
    print("DE.PULSE post-Stable continuity: FAIL", file=sys.stderr)
    for item in items:
        print(f" - {item}", file=sys.stderr)
    return 1


def stable_documentation_errors(stable_tag: str) -> list[str]:
    errors: list[str] = []
    handoff_path = ROOT / "handoff" / "CURRENT.md"
    handoff = handoff_path.read_text(encoding="utf-8", errors="ignore") if handoff_path.is_file() else ""
    if f"**Certified Stable:** `{stable_tag}`" not in handoff:
        errors.append("handoff/CURRENT.md does not identify the aligned Stable tag")
    for rel in (
        "adaptive-governance/CURRENT_ADAPTIVE_ROADMAP.md",
        "adaptive-governance/CURRENT_ADAPTIVE_BUILD_PLAN.md",
        "adaptive-governance/CURRENT_ADAPTIVE_BUILD_PROCESS.md",
        "adaptive-governance/CURRENT_ADAPTIVE_DELIVERY_PROCESS.md",
    ):
        path = ROOT / rel
        text = path.read_text(encoding="utf-8", errors="ignore") if path.is_file() else ""
        if stable_tag not in text:
            errors.append(f"{rel} is not aligned to {stable_tag}")
    return errors


def validate_stable_state(identity: dict, build: dict, evidence: dict) -> list[str]:
    errors: list[str] = []
    state = load_json(ROOT / "governance" / "current-state.json")
    stable = state.get("stable", {}) if isinstance(state.get("stable"), dict) else {}
    version = clean_version(identity.get("version"))
    expected_tag = stable_tag_for(version)
    checkpoint_version = clean_version(build.get("release"))
    evidence_version = clean_version(evidence.get("release"))

    if clean_version(stable.get("productVersion")) != version:
        errors.append("current-state Stable productVersion / release identity drift")
    if str(stable.get("tag", "")) != expected_tag:
        errors.append(f"current-state Stable tag mismatch: {stable.get('tag')!r} != {expected_tag!r}")
    if str(stable.get("buildId", "")) != str(identity.get("build_id", "")):
        errors.append("current-state Stable buildId / release identity drift")
    if str(stable.get("platformBuildNumber", "")) != str(identity.get("bundle_version", "")):
        errors.append("current-state platform build number / release identity drift")
    if stable.get("publication") != "PASS_NO_REBUILD":
        errors.append("current-state Stable publication must remain PASS_NO_REBUILD")
    if checkpoint_version != version or evidence_version != version:
        errors.append("Stable identity/build checkpoint/release evidence must align after publication")

    certified = build.get("certifiedStable", {}) if isinstance(build.get("certifiedStable"), dict) else {}
    ev_stable = evidence.get("stable", {}) if isinstance(evidence.get("stable"), dict) else {}
    if certified.get("tag") != expected_tag or ev_stable.get("tag") != expected_tag:
        errors.append("Stable checkpoint tag / versioning-policy drift")
    if certified.get("promotionCommit") != stable.get("candidateSha"):
        errors.append("build checkpoint promotion candidate / current Stable drift")
    if ev_stable.get("promotionCommit") != stable.get("candidateSha"):
        errors.append("release evidence promotion candidate / current Stable drift")
    if certified.get("sourceFingerprint") != stable.get("sourceFingerprint"):
        errors.append("build checkpoint fingerprint / current Stable drift")
    if ev_stable.get("sourceFingerprint") != stable.get("sourceFingerprint"):
        errors.append("release evidence fingerprint / current Stable drift")

    manifest = ROOT / "release" / f"v{version}" / "stable-evidence-manifest.json"
    if not manifest.is_file():
        errors.append(f"Stable evidence manifest missing: {manifest.relative_to(ROOT)}")
    else:
        m = load_json(manifest)
        for field, expected in (
            ("release", f"v{version}"),
            ("status", "STABLE_PUBLISHED"),
            ("stableTag", expected_tag),
            ("certifiedCandidate", stable.get("candidateSha")),
            ("qualifiedSourceHead", stable.get("qualifiedSourceSha")),
            ("sourceFingerprint", stable.get("sourceFingerprint")),
            ("buildId", stable.get("buildId")),
        ):
            if m.get(field) != expected:
                errors.append(f"Stable evidence manifest {field} drift")
    errors.extend(stable_documentation_errors(expected_tag))
    return errors


def validate_inflight(identity: dict, build: dict, branch: str) -> tuple[str, list[str]]:
    errors: list[str] = []
    state = load_json(ROOT / "governance" / "current-state.json")
    stable = state.get("stable", {}) if isinstance(state.get("stable"), dict) else {}
    active = state.get("activeWorkSlice", {}) if isinstance(state.get("activeWorkSlice"), dict) else {}
    gate = state.get("productCapabilityGate", {}) if isinstance(state.get("productCapabilityGate"), dict) else {}
    identity_version = clean_version(identity.get("version"))
    checkpoint_version = clean_version(build.get("release"))
    expected_candidate_branch = f"v{identity_version}-development"

    work_rel = str(gate.get("workSlicePath", "")).strip()
    product = load_json(ROOT / work_rel) if work_rel else {}
    product_branch = str(gate.get("reservedBranch", "")).strip()
    process_branch = str(active.get("branch", "")).strip()
    closure_branch = str(active.get("closureBranch", "")).strip()

    if branch == expected_candidate_branch:
        return "IN_FLIGHT_CANDIDATE", errors

    if branch == product_branch and product and str(gate.get("reservationStatus", "")).upper() in ACTIVE:
        if gate.get("blocked") is not False:
            errors.append("registered product work requires unblocked product capability gate")
        if product.get("type") == "PRODUCT_RELEASE_CLOSURE":
            previous = clean_version(identity.get("previous_stable"))
            if checkpoint_version != previous:
                errors.append("release closure checkpoint must remain immediate predecessor Stable")
            if clean_version(stable.get("productVersion")) != previous:
                errors.append("release closure must preserve predecessor Stable until publication")
            if str(stable.get("tag", "")) != stable_tag_for(previous):
                errors.append("release closure predecessor Stable tag drift")
            if product.get("targetStable") != stable_tag_for(identity_version):
                errors.append("release closure targetStable / versioning-policy drift")
            return "PRODUCT_RELEASE_CLOSURE", errors
        return "PRODUCT_WORK_SLICE", errors

    if branch in {process_branch, closure_branch} and active.get("type") == "PROCESS_RELEASE_ENGINEERING":
        if active.get("publicProductVersion") is not None or active.get("productBehaviorChange") is not False:
            errors.append("process work slice must remain product-version neutral")
        return "PROCESS_WORK_SLICE_CLOSURE" if branch == closure_branch else "PROCESS_WORK_SLICE", errors

    errors.append(
        "identity/checkpoint differ outside an allowed candidate/product/process branch: "
        f"candidate={expected_candidate_branch}, product={product_branch or '<none>'}, "
        f"process={process_branch or '<none>'}, closure={closure_branch or '<none>'}, current={branch or '<detached>'}"
    )
    return "INVALID", errors


def main() -> int:
    identity = load_json(ROOT / "release_identity.json")
    build = load_json(ROOT / ".depulse-certification" / "resume" / "build-checkpoint.json")
    evidence = load_json(ROOT / ".depulse-certification" / "resume" / "release-evidence-checkpoint.json")
    errors: list[str] = []
    identity_version = clean_version(identity.get("version"))
    checkpoint_version = clean_version(build.get("release"))
    evidence_version = clean_version(evidence.get("release"))
    branch = branch_name()

    if not identity_version or not checkpoint_version or not evidence_version:
        return fail(["release identity/build checkpoint/release evidence version missing"])
    if checkpoint_version != evidence_version:
        return fail([f"build/evidence Stable mismatch: {checkpoint_version} != {evidence_version}"])

    if identity_version == checkpoint_version:
        mode = "STABLE_ALIGNED"
        errors.extend(validate_stable_state(identity, build, evidence))
    else:
        mode, inflight_errors = validate_inflight(identity, build, branch)
        errors.extend(inflight_errors)

    if errors:
        return fail(errors)

    print(
        "DE.PULSE post-Stable continuity: PASS · "
        f"mode={mode} · branch={branch or '<detached>'} · identity=v{identity_version} · durable checkpoint=v{checkpoint_version}"
    )
    if mode == "STABLE_ALIGNED":
        print("current Stable / checkpoint / manifest / documentation / SemVer tag policy: PASS")
    if mode == "PRODUCT_RELEASE_CLOSURE":
        print("active release-closure candidate / predecessor Stable separation: PASS")
        print("publication remains deferred until canonical release gates: PASS")
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
