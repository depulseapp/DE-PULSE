#!/usr/bin/env python3
from __future__ import annotations

import json
from pathlib import Path
import re
import sys

ROOT = Path(__file__).resolve().parents[2]
CHECKPOINT = ROOT / ".depulse-certification" / "resume" / "release-evidence-checkpoint.json"


def fail(errors: list[str]) -> int:
    print("DE.PULSE Stable evidence gate: FAIL", file=sys.stderr)
    for error in errors:
        print(f" - {error}", file=sys.stderr)
    return 1


def main() -> int:
    errors: list[str] = []
    try:
        checkpoint = json.loads(CHECKPOINT.read_text(encoding="utf-8"))
    except Exception as exc:
        return fail([f"authoritative checkpoint unreadable: {exc}"])

    release = str(checkpoint.get("release", ""))
    version = release.removeprefix("v")
    manifest_path = ROOT / "release" / version / "stable-evidence-manifest.json"
    if not manifest_path.is_file():
        return fail([f"durable Stable manifest missing for {release}: {manifest_path.relative_to(ROOT)}"])

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
            errors.append(f"{key} mismatch: manifest={manifest.get(key)!r} checkpoint={value!r}")

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
            errors.append(f"run {key} mismatch: manifest={runs.get(key)!r} checkpoint={value!r}")

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
    if provenance.get("authoritativeCheckpoint") != ".depulse-certification/resume/release-evidence-checkpoint.json":
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
    print(f"durable Stable manifest: release/{version}/stable-evidence-manifest.json")
    print(f"immutable Stable authority preserved: {manifest['stableTag']} -> {manifest['certifiedCandidate']}")
    print("Fast/Qualified/Release run binding: PASS")
    print("native/G15/G16 artifact digest binding: PASS")
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
