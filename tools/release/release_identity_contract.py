#!/usr/bin/env python3
from __future__ import annotations

import argparse
import json
from pathlib import Path
import re
import sys
from typing import Any

ROOT = Path(__file__).resolve().parents[2]
POLICY_PATH = ROOT / "governance" / "versioning-policy.json"
IDENTITY_PATH = ROOT / "release_identity.json"
CURRENT_STATE_PATH = ROOT / "governance" / "current-state.json"
SEMVER_RE = re.compile(r"^(0|[1-9]\d*)\.(0|[1-9]\d*)\.(0|[1-9]\d*)(?:-([0-9A-Za-z.-]+))?(?:\+([0-9A-Za-z.-]+))?$")


def load(path: Path) -> dict[str, Any]:
    return json.loads(path.read_text(encoding="utf-8"))


def parse_semver(value: str) -> tuple[int, int, int, str | None]:
    raw = value.strip().removeprefix("v")
    match = SEMVER_RE.fullmatch(raw)
    if not match:
        raise ValueError(f"not a valid SemVer product version: {value!r}")
    return int(match.group(1)), int(match.group(2)), int(match.group(3)), match.group(4)


def core(value: str) -> tuple[int, int, int]:
    major, minor, patch, _ = parse_semver(value)
    return major, minor, patch


def stable_tag_for_version(policy: dict[str, Any], version: str) -> str:
    parsed = parse_semver(version)
    if parsed[3] is not None:
        raise ValueError("Stable product version cannot contain a SemVer prerelease component")
    cutover = core(str(policy["effectiveAfterProductVersion"]))
    pattern = policy["legacyStableTagPattern"] if core(version) <= cutover else policy["futureStableTagPattern"]
    return str(pattern).format(productVersion=version)


def normalize_product_version(value: str) -> str:
    raw = value.strip()
    if raw.startswith("v"):
        raw = raw[1:]
    if raw.endswith("-stable"):
        raw = raw[: -len("-stable")]
    parse_semver(raw)
    return raw


def verify() -> list[str]:
    errors: list[str] = []
    policy = load(POLICY_PATH)
    identity = load(IDENTITY_PATH)
    current = load(CURRENT_STATE_PATH)
    stable = current.get("stable", {})
    work = current.get("activeWorkSlice", {})

    if policy.get("schema") != "DE.PULSE-VERSIONING-POLICY-1":
        errors.append("VERSIONING_POLICY_SCHEMA")
    if policy.get("decision") != "PROSPECTIVE_SEMVER":
        errors.append("VERSIONING_DECISION_NOT_PROSPECTIVE_SEMVER")

    version = str(identity.get("version", "")).strip()
    try:
        parse_semver(version)
    except ValueError as exc:
        errors.append(f"PRODUCT_VERSION_INVALID:{exc}")
        return errors

    if str(identity.get("display_version")) != f"DE.PULSE v{version}":
        errors.append("DISPLAY_VERSION_MISMATCH")
    if str(identity.get("channel")) != "STABLE":
        errors.append("CHANNEL_NOT_STABLE")

    build_id = str(identity.get("build_id", "")).strip()
    if not build_id or build_id in {version, f"v{version}"}:
        errors.append("BUILD_ID_NOT_INDEPENDENT")

    bundle = str(identity.get("bundle_version", "")).strip()
    if not bundle.isdigit() or int(bundle) <= 0:
        errors.append("PLATFORM_BUILD_NUMBER_INVALID")

    stable_version = str(stable.get("productVersion", "")).strip()
    stable_build = str(stable.get("platformBuildNumber", "")).strip()
    if not stable_version or not stable_build.isdigit():
        errors.append("CURRENT_STABLE_IDENTITY_INCOMPLETE")
    else:
        try:
            candidate_core = core(version)
            stable_core = core(stable_version)
            if candidate_core < stable_core:
                errors.append("PRODUCT_VERSION_REGRESSION")
            elif candidate_core == stable_core:
                if bundle != stable_build:
                    errors.append("CURRENT_STABLE_PLATFORM_BUILD_MISMATCH")
                expected_tag = stable_tag_for_version(policy, version)
                if str(stable.get("tag")) != expected_tag:
                    errors.append("CURRENT_STABLE_TAG_POLICY_MISMATCH")
                if build_id != str(stable.get("buildId")):
                    errors.append("CURRENT_STABLE_BUILD_ID_MISMATCH")
            elif int(bundle) <= int(stable_build):
                errors.append("PLATFORM_BUILD_NUMBER_NOT_MONOTONIC")
        except ValueError as exc:
            errors.append(f"CURRENT_STABLE_VERSION_INVALID:{exc}")

    previous_raw = str(identity.get("previous_stable", "")).strip()
    try:
        previous_version = normalize_product_version(previous_raw)
        if core(previous_version) >= core(version):
            errors.append("PREVIOUS_STABLE_NOT_OLDER")
    except ValueError as exc:
        errors.append(f"PREVIOUS_STABLE_INVALID:{exc}")

    if not str(work.get("workSliceId", "")).strip():
        errors.append("WORK_SLICE_ID_MISSING")
    if work.get("publicProductVersion") is not None and str(work.get("publicProductVersion")) == str(work.get("workSliceId")):
        errors.append("WORK_SLICE_CONFLATED_WITH_PRODUCT_VERSION")
    if work.get("type") == "PROCESS_RELEASE_ENGINEERING" and work.get("publicProductVersion") is not None:
        errors.append("PROCESS_WORK_SLICE_MUST_NOT_CONSUME_PRODUCT_VERSION")

    for field in ("candidateSha", "sourceFingerprint", "buildId", "evidenceSchemaVersion"):
        if not stable.get(field):
            errors.append(f"CURRENT_STABLE_{field.upper()}_MISSING")
    if stable.get("candidateSha") == stable.get("sourceFingerprint"):
        errors.append("CANDIDATE_SHA_CONFLATED_WITH_SOURCE_FINGERPRINT")

    manifest_path = ROOT / "release" / f"v{version}" / "certification-manifest.json"
    if manifest_path.is_file():
        manifest = load(manifest_path)
        if manifest.get("productVersion") != version:
            errors.append("G12_MANIFEST_PRODUCT_VERSION_MISMATCH")
        if not manifest.get("workSliceId"):
            errors.append("G12_MANIFEST_WORK_SLICE_MISSING")
        if not isinstance(manifest.get("evidenceSchemaVersion"), int):
            errors.append("G12_MANIFEST_EVIDENCE_SCHEMA_INVALID")

    return errors


def main() -> int:
    parser = argparse.ArgumentParser(description="DE.PULSE product/work-slice/build/evidence identity contract")
    parser.add_argument("--verify", action="store_true")
    parser.add_argument("--stable-tag", default="")
    parser.add_argument("--previous-stable-tag", action="store_true")
    args = parser.parse_args()

    policy = load(POLICY_PATH)
    identity = load(IDENTITY_PATH)

    if args.stable_tag:
        print(stable_tag_for_version(policy, normalize_product_version(args.stable_tag)))
        return 0
    if args.previous_stable_tag:
        print(stable_tag_for_version(policy, normalize_product_version(str(identity["previous_stable"]))))
        return 0

    errors = verify()
    if errors:
        print("DE.PULSE release identity contract: FAIL", file=sys.stderr)
        for error in errors:
            print(f" - {error}", file=sys.stderr)
        return 1

    version = str(identity["version"])
    print("DE.PULSE release identity contract: PASS")
    print(f"productVersion={version}")
    print(f"stableTag={stable_tag_for_version(policy, version)}")
    print(f"workSliceId={load(CURRENT_STATE_PATH)['activeWorkSlice']['workSliceId']}")
    print(f"buildId={identity['build_id']}")
    print(f"platformBuildNumber={identity['bundle_version']}")
    print("product/work-slice/source/build/evidence identities: DISTINCT")
    print("prospective public versioning: Semantic Versioning 2.0.0")
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
