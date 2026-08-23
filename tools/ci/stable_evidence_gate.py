#!/usr/bin/env python3
from __future__ import annotations

import json
import os
from pathlib import Path
import re
import subprocess
import sys

ROOT = Path(__file__).resolve().parents[2]
CHECKPOINT = ROOT / ".depulse-certification" / "resume" / "release-evidence-checkpoint.json"


def fail(errors: list[str]) -> int:
    print("DE.PULSE Stable evidence gate: FAIL", file=sys.stderr)
    for error in errors:
        print(f" - {error}", file=sys.stderr)
    return 1


def load_json(path: Path) -> dict:
    try:
        value = json.loads(path.read_text(encoding="utf-8"))
        return value if isinstance(value, dict) else {}
    except Exception:
        return {}


def branch_name() -> str:
    env = (
        os.environ.get("GITHUB_HEAD_REF")
        or os.environ.get("GITHUB_REF_NAME")
        or ""
    ).strip()
    if env:
        return env.removeprefix("refs/heads/")
    proc = subprocess.run(
        ["git", "branch", "--show-current"],
        cwd=ROOT,
        text=True,
        capture_output=True,
        check=False,
    )
    return proc.stdout.strip()


def run_local_gate(filename: str, label: str, extra_args: list[str] | None = None) -> list[str]:
    gate = ROOT / "tools" / "ci" / filename
    if not gate.is_file():
        return [f"{label} missing: {gate.relative_to(ROOT)}"]
    command = [sys.executable, str(gate), *(extra_args or [])]
    result = subprocess.run(command, cwd=ROOT, check=False)
    if result.returncode != 0:
        return [f"{label} failed with exit code {result.returncode}"]
    return []


def registered_process_target() -> tuple[str, list[str]]:
    """Return the already-published current Stable candidate for a governed process branch.

    This is deliberately narrower than a generic existing-tag exception. It lets
    the canonical release-state coherence owner exercise its existing exact-target
    REUSE semantics only when current-state + work-slice metadata prove that the
    branch is process/release engineering on top of an already-published Stable.
    """
    state_path = ROOT / "governance" / "current-state.json"
    if not state_path.is_file():
        return "", []

    state = load_json(state_path)
    stable = state.get("stable", {}) if isinstance(state.get("stable"), dict) else {}
    active = (
        state.get("activeWorkSlice", {})
        if isinstance(state.get("activeWorkSlice"), dict)
        else {}
    )
    branch = branch_name()
    registered_branch = str(active.get("branch", "")).strip()
    if (
        active.get("type") != "PROCESS_RELEASE_ENGINEERING"
        or not registered_branch
        or branch != registered_branch
    ):
        return "", []

    errors: list[str] = []
    identity = load_json(ROOT / "release_identity.json")
    version = str(identity.get("version", "")).strip().lstrip("v")
    build_id = str(identity.get("build_id", "")).strip()
    expected_tag = f"v{version}-stable" if version else ""
    candidate = str(stable.get("candidateSha", "")).strip()
    fingerprint = str(stable.get("sourceFingerprint", "")).strip()
    work_slice_id = str(active.get("workSliceId", "")).strip()
    work_slice_path = (
        ROOT / "governance" / "work-slices" / work_slice_id / "work-slice.json"
        if work_slice_id
        else ROOT / "governance" / "work-slices" / "<missing>" / "work-slice.json"
    )
    work_slice = load_json(work_slice_path)

    if identity.get("channel") != "STABLE":
        errors.append("registered process branch requires unchanged STABLE release identity")
    if str(stable.get("productVersion", "")).strip().lstrip("v") != version:
        errors.append("registered process Stable productVersion / release identity drift")
    if str(stable.get("buildId", "")).strip() != build_id:
        errors.append("registered process Stable buildId / release identity drift")
    if str(stable.get("platformBuildNumber", "")).strip() != str(identity.get("bundle_version", "")).strip():
        errors.append("registered process platform build number / release identity drift")
    if str(stable.get("tag", "")).strip() != expected_tag:
        errors.append("registered process Stable tag / release identity drift")
    if stable.get("publication") != "PASS_NO_REBUILD":
        errors.append("registered process Stable publication must remain PASS_NO_REBUILD")
    if not re.fullmatch(r"[0-9a-f]{40}", candidate):
        errors.append(f"registered process Stable candidate invalid: {candidate!r}")
    if not re.fullmatch(r"[0-9a-f]{64}", fingerprint):
        errors.append("registered process Stable source fingerprint invalid")
    if active.get("publicProductVersion") is not None:
        errors.append("registered process work slice must not consume a public product version")
    if active.get("productBehaviorChange") is not False:
        errors.append("registered process work slice must declare productBehaviorChange=false")

    if not work_slice:
        errors.append(f"registered process work-slice metadata missing: {work_slice_path.relative_to(ROOT)}")
    else:
        if str(work_slice.get("workSliceId", "")) != work_slice_id:
            errors.append("registered process current-state/work-slice id drift")
        if work_slice.get("type") != "PROCESS_RELEASE_ENGINEERING":
            errors.append("registered process work-slice type drift")
        if work_slice.get("publicProductVersion") is not None:
            errors.append("process work-slice metadata must not consume a public product version")
        if work_slice.get("productBehaviorChange") is not False:
            errors.append("process work-slice metadata must declare productBehaviorChange=false")
        if str(work_slice.get("branch", "")) != branch:
            errors.append("registered process work-slice branch drift")
        if str(work_slice.get("baselineCandidateSha", "")) != candidate:
            errors.append("registered process baseline/current Stable candidate drift")
        if str(work_slice.get("baselineSourceFingerprint", "")) != fingerprint:
            errors.append("registered process baseline/current Stable fingerprint drift")
        if str(work_slice.get("baselineBuildId", "")) != build_id:
            errors.append("registered process baseline/current Stable build drift")
        if str(work_slice.get("stableProductVersionAtStart", "")).lstrip("v") != version:
            errors.append("registered process Stable-at-start / release identity drift")
        if work_slice.get("blocksNextProductCapability") is not True:
            errors.append("registered process work slice must keep next product capability blocked")

    capability_gate = (
        state.get("productCapabilityGate", {})
        if isinstance(state.get("productCapabilityGate"), dict)
        else {}
    )
    if capability_gate.get("blocked") is not True:
        errors.append("next product capability must remain blocked during process work")
    if capability_gate.get("blockedByIssue") != active.get("issue"):
        errors.append("product capability blocker / active process issue drift")

    return candidate if not errors else "", errors


def main() -> int:
    errors: list[str] = []
    process_candidate, process_errors = registered_process_target()
    errors.extend(process_errors)
    errors.extend(run_local_gate("ci_hardening_gate.py", "CI hardening contract"))
    errors.extend(
        run_local_gate(
            "release_state_coherence_self_test.py",
            "Release State Coherence self-test",
        )
    )
    coherence_args = (
        ["--g11-candidate-sha", process_candidate]
        if process_candidate
        else []
    )
    errors.extend(
        run_local_gate(
            "release_state_coherence.py",
            "Release State Coherence",
            coherence_args,
        )
    )
    if errors:
        return fail(errors)

    try:
        checkpoint = json.loads(CHECKPOINT.read_text(encoding="utf-8"))
    except Exception as exc:
        return fail([f"authoritative checkpoint unreadable: {exc}"])

    release = str(checkpoint.get("release", ""))
    manifest_path = ROOT / "release" / release / "stable-evidence-manifest.json"
    if not manifest_path.is_file():
        return fail(
            [
                f"durable Stable manifest missing for {release}: "
                f"{manifest_path.relative_to(ROOT)}"
            ]
        )

    try:
        manifest = json.loads(manifest_path.read_text(encoding="utf-8"))
    except Exception as exc:
        return fail([f"Stable manifest unreadable: {exc}"])

    stable = checkpoint.get("stable", {})
    candidate = checkpoint.get("currentCandidate", {})
    evidence = checkpoint.get("evidence", {})

    expected = {
        "schema": "DE.PULSE-STABLE-EVIDENCE-1",
        "release": release,
        "status": "STABLE_PUBLISHED",
        "buildId": stable.get("buildId"),
        "stableTag": stable.get("tag"),
        "certifiedCandidate": stable.get("promotionCommit"),
        "qualifiedSourceHead": candidate.get("qualificationHead"),
        "sourceFingerprint": stable.get("sourceFingerprint"),
    }
    for key, value in expected.items():
        if manifest.get(key) != value:
            errors.append(
                f"{key} mismatch: manifest={manifest.get(key)!r} checkpoint={value!r}"
            )

    runs = manifest.get("runs", {})
    g10 = evidence.get("G10", {})
    expected_runs = {
        "fast": g10.get("fastRun"),
        "qualified": g10.get("qualifiedRun"),
        "productFullQualified": g10.get("productFullQualifiedRun"),
        "canonicalReleaseG11G16": stable.get("certificationRun"),
    }
    for key, value in expected_runs.items():
        if runs.get(key) != value:
            errors.append(
                f"run {key} mismatch: manifest={runs.get(key)!r} checkpoint={value!r}"
            )

    artifact_map = {
        "macOSAppleSilicon": "G13-G14-macOS-Apple-Silicon",
        "windowsX64": "G13-G14-Windows-x64",
        "g15ReleaseAssurance": "G15",
        "g16Closure": "G16",
    }
    artifacts = manifest.get("artifacts", {})
    for manifest_key, evidence_key in artifact_map.items():
        source = evidence.get(evidence_key, {})
        row = artifacts.get(manifest_key, {})
        if row.get("artifactId") != source.get("artifactId"):
            errors.append(f"artifact id mismatch for {manifest_key}")
        if row.get("digest") != source.get("artifactDigest"):
            errors.append(f"artifact digest mismatch for {manifest_key}")
        digest = str(row.get("digest", ""))
        if not re.fullmatch(r"sha256:[0-9a-f]{64}", digest):
            errors.append(f"invalid sha256 digest for {manifest_key}: {digest!r}")

    provenance = manifest.get("provenance", {})
    if (
        provenance.get("authoritativeCheckpoint")
        != ".depulse-certification/resume/release-evidence-checkpoint.json"
    ):
        errors.append("manifest must point to the canonical release evidence checkpoint")
    if provenance.get("immutableStableTagRemainsAuthority") is not True:
        errors.append("manifest must preserve immutable Stable tag authority")
    if provenance.get("manifestCommittedAfterStable") is not True:
        errors.append("retrospective manifest must state it was committed after Stable")
    if provenance.get("manifestDoesNotRedefineStableArtifact") is not True:
        errors.append("retrospective manifest must explicitly not redefine Stable")

    if errors:
        return fail(errors)

    print("DE.PULSE Stable evidence gate: PASS")
    print("CI hardening contract: PASS")
    print("Release State Coherence aggregate preflight: PASS")
    if process_candidate:
        print(
            "registered current Stable exact-tag reuse: PASS · "
            f"candidate={process_candidate}"
        )
    print(
        "durable prior-Stable manifest: "
        f"release/{release}/stable-evidence-manifest.json"
    )
    print(
        "immutable Stable authority preserved: "
        f"{manifest['stableTag']} -> {manifest['certifiedCandidate']}"
    )
    print("Fast/Qualified/Release run binding: PASS")
    print("native/G15/G16 artifact digest binding: PASS")
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
