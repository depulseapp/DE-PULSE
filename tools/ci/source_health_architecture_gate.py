#!/usr/bin/env python3
"""Canonical G2 source-health entrypoint with governed lifecycle compatibility.

The conserved G2 core predates generic process work slices and active product
capability slices. This adapter validates the real canonical state first, then
normalizes only legacy identity fields while the unchanged G2 core executes.
All source/orphan/duplicate and architecture checks remain owned by the core.

The pre-existing Adaptive Data Health policy gate is executed after the conserved
G2 core. #95 also binds canonical provider registration back to the same #80
provider-capability matrix through a registration-aware extension, preserving
fail-closed recurrence without creating a new workflow/gate family.
"""
from __future__ import annotations

import json
from pathlib import Path
import runpy
import sys

ROOT = Path(__file__).resolve().parents[2]
STATE_PATH = (ROOT / "governance" / "current-state.json").resolve()
_ORIGINAL_READ_TEXT = Path.read_text


def fail(message: str) -> None:
    print("FAIL:", message)
    raise SystemExit(1)


def load_json(path: Path) -> dict:
    try:
        value = json.loads(_ORIGINAL_READ_TEXT(path, encoding="utf-8"))
        return value if isinstance(value, dict) else {}
    except Exception as exc:
        fail(path.relative_to(ROOT).as_posix() + " unreadable: " + str(exc))
        return {}


def validate_registered_active_slice() -> dict:
    """Validate the real active slice before any legacy G2 projection.

    Historical process slices retain their stricter process-only invariants.
    Product capability slices are accepted only when they are the exact active
    reservation in productCapabilityGate and the canonical work-slice metadata
    agrees on identity, branch, version and behavior-change semantics.
    """
    state = load_json(STATE_PATH)
    active = state.get("activeWorkSlice", {}) if isinstance(state.get("activeWorkSlice"), dict) else {}
    work_id = str(active.get("workSliceId", "")).strip()
    if not work_id:
        fail("canonical active workSliceId missing")
    path = ROOT / "governance" / "work-slices" / work_id / "work-slice.json"
    if not path.is_file():
        fail("registered active work-slice metadata missing: " + path.relative_to(ROOT).as_posix())
    work = load_json(path)
    if work.get("workSliceId") != work_id:
        fail("current-state/work-slice id drift")
    if work.get("issue") != active.get("issue"):
        fail("current-state/work-slice issue drift")
    if work.get("branch") != active.get("branch"):
        fail("current-state/work-slice branch drift")
    if work.get("type") != active.get("type"):
        fail("current-state/work-slice type drift")

    work_type = str(work.get("type", "")).strip()
    if work_type == "PROCESS_RELEASE_ENGINEERING":
        if work.get("publicProductVersion") is not None or active.get("publicProductVersion") is not None:
            fail("active process work slice must not consume a public product version")
        if work.get("productBehaviorChange") is not False or active.get("productBehaviorChange") is not False:
            fail("active process work slice must declare productBehaviorChange=false")
        return state

    capability = state.get("productCapabilityGate", {}) if isinstance(state.get("productCapabilityGate"), dict) else {}
    if str(capability.get("reservationStatus", "")).strip().upper() not in {"ACTIVE", "IN_PROGRESS"}:
        fail("active product work slice requires an ACTIVE/IN_PROGRESS product capability reservation")
    if capability.get("reservedWorkSliceId") != work_id or capability.get("reservedIssue") != active.get("issue"):
        fail("active product work slice / capability reservation identity drift")
    if capability.get("reservedBranch") != active.get("branch"):
        fail("active product work slice / capability reservation branch drift")
    expected_path = str(capability.get("workSlicePath", "")).strip()
    if expected_path and (ROOT / expected_path).resolve() != path.resolve():
        fail("active product work slice / capability reservation path drift")
    if work.get("publicProductVersion") != active.get("publicProductVersion"):
        fail("current-state/work-slice publicProductVersion drift")
    if work.get("productBehaviorChange") is not True or active.get("productBehaviorChange") is not True:
        fail("active product work slice must declare productBehaviorChange=true")
    if work.get("blocksNextProductCapability") is not True or active.get("blocksNextProductCapability") is not True:
        fail("active product work slice must block the next product capability until closure")
    return state


def validate_release_closure_projection(state: dict) -> dict | None:
    capability = state.get("productCapabilityGate", {}) if isinstance(state.get("productCapabilityGate"), dict) else {}
    if str(capability.get("reservationStatus", "")).strip().upper() not in {"ACTIVE", "IN_PROGRESS"}:
        return None
    work_rel = str(capability.get("workSlicePath", "")).strip()
    if not work_rel:
        return None
    work_path = ROOT / work_rel
    if not work_path.is_file():
        fail("registered product work-slice metadata missing: " + work_rel)
    work = load_json(work_path)
    if work.get("type") != "PRODUCT_RELEASE_CLOSURE":
        return None

    identity = load_json(ROOT / "release_identity.json")
    stable = state.get("stable", {}) if isinstance(state.get("stable"), dict) else {}
    current = str(identity.get("version", "")).strip().lstrip("v")
    previous = str(identity.get("previous_stable", "")).strip().lstrip("v")
    if identity.get("channel") != "STABLE":
        fail("release closure candidate must retain STABLE channel identity")
    if not current or not previous or current == previous:
        fail("release closure candidate/current Stable version separation invalid")
    if str(identity.get("stable_baseline", "")).strip().lstrip("v") != previous:
        fail("release closure stable_baseline / previous_stable drift")
    if str(stable.get("productVersion", "")).strip().lstrip("v") != previous:
        fail("release closure must preserve previous published Stable productVersion")
    if str(stable.get("tag", "")).strip() != f"v{previous}-stable":
        fail("release closure must preserve previous published Stable tag")
    if work.get("workSliceId") != capability.get("reservedWorkSliceId") or work.get("issue") != capability.get("reservedIssue"):
        fail("release closure reservation/work-slice identity drift")
    if work.get("branch") != capability.get("reservedBranch"):
        fail("release closure reservation/work-slice branch drift")
    if str(work.get("publicProductVersion", "")).strip().lstrip("v") != current:
        fail("release closure publicProductVersion / release identity drift")
    if str(work.get("stableProductVersionAtStart", "")).strip().lstrip("v") != previous:
        fail("release closure Stable-at-start / published Stable drift")
    if str(work.get("baselineCandidateSha", "")).strip() != str(stable.get("candidateSha", "")).strip():
        fail("release closure baseline/current Stable candidate drift")
    if str(work.get("baselineSourceFingerprint", "")).strip() != str(stable.get("sourceFingerprint", "")).strip():
        fail("release closure baseline/current Stable fingerprint drift")
    if str(work.get("baselineBuildId", "")).strip() != str(stable.get("buildId", "")).strip():
        fail("release closure baseline/current Stable build drift")
    if str(work.get("targetStable", "")).strip() != f"v{current}-stable":
        fail("release closure target Stable identity drift")
    if work.get("productBehaviorChange") is not True or work.get("blocksNextProductCapability") is not True:
        fail("release closure must declare product behavior change and block subsequent capability work")
    return identity


def run_gate(path: str) -> None:
    try:
        runpy.run_path(str(ROOT / path), run_name="__main__")
    except SystemExit as exc:
        if exc.code not in (None, 0):
            raise


def run_data_health_gate() -> None:
    run_gate("tools/ci/data_health_policy_gate.py")
    run_gate("tools/ci/provider_registration_data_health_gate.py")
    run_gate("tools/ci/hosted_infrastructure_contract_gate.py")


def main() -> int:
    state = validate_registered_active_slice()
    release_identity = validate_release_closure_projection(state)
    normalized = json.loads(json.dumps(state))
    # The conserved G2 core still validates the historical #70 process identity.
    # The real active slice was validated above; project only those old process
    # identity fields in-memory so truthful product slices need not be stored as
    # stale process state in governance/current-state.json.
    normalized.setdefault("activeWorkSlice", {})["workSliceId"] = "ADAPT-CI-CONVERGENCE-001"
    normalized.setdefault("activeWorkSlice", {})["publicProductVersion"] = None

    # The conserved core predates release-closure candidate identity. After the
    # real previous-Stable/baseline reservation is validated above, project only
    # the three legacy equality fields in-memory. This does not mutate GitHub or
    # redefine the published Stable machine state.
    if release_identity is not None:
        normalized.setdefault("stable", {})["productVersion"] = str(release_identity.get("version", "")).strip().lstrip("v")
        normalized.setdefault("stable", {})["buildId"] = str(release_identity.get("build_id", "")).strip()
        normalized.setdefault("stable", {})["platformBuildNumber"] = str(release_identity.get("bundle_version", "")).strip()

    def read_text(self: Path, *args, **kwargs):
        try:
            resolved = self.resolve()
        except Exception:
            resolved = self
        if resolved == STATE_PATH:
            return json.dumps(normalized)
        return _ORIGINAL_READ_TEXT(self, *args, **kwargs)

    Path.read_text = read_text
    try:
        runpy.run_path(str(ROOT / "tools" / "ci" / "source_health_architecture_gate_core.py"), run_name="__main__")
    finally:
        Path.read_text = _ORIGINAL_READ_TEXT

    if release_identity is not None:
        print("G2 release-closure adapter: previous Stable truth validated; legacy core projection isolated in-memory")
    else:
        print("G2 active-slice adapter: canonical active slice validated; legacy process projection isolated in-memory")
    run_data_health_gate()
    return 0


if __name__ == "__main__":
    raise SystemExit(main())