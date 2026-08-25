#!/usr/bin/env python3
"""Deterministic T2 unit/contract/static/property assurance audit.

The gate is lifecycle-aware: while T2 is IN_PROGRESS it must be the sole active
child; after T2 is COMPLETE it remains a durable regression gate while later
closure tracks advance. T2 never certifies T3-T10.
"""
from __future__ import annotations

import argparse
import hashlib
import json
from pathlib import Path
import sys
from typing import Any

ROOT = Path(__file__).resolve().parents[2]
PROGRAM = ROOT / "governance" / "programs" / "ADAPT-V18-FINAL-CLOSURE-10-10-001"
LEDGER = PROGRAM / "feature-assurance-ledger.json"
FREEZE = PROGRAM / "feature-assurance-ledger-freeze.json"
SCAN1 = PROGRAM / "T1_INDEPENDENT_OMISSION_SCAN.json"
SCAN2 = PROGRAM / "T1_INDEPENDENT_OMISSION_SCAN_2.json"
RECONCILIATION = PROGRAM / "T1_FINAL_RECONCILIATION.json"
T2 = PROGRAM / "T2_UNIT_CONTRACT_ASSURANCE.json"
CURRENT_STATE = ROOT / "governance" / "current-state.json"
CLOSURE = ROOT / "governance" / "work-slices" / "ADAPT-V18-FINAL-CLOSURE-10-10-001" / "closure.json"
CI_FAST = ROOT / ".github" / "workflows" / "ci-fast.yml"

T2_CLASSES = {"GO_UNIT_PACKAGE", "RENDERER_NODE_CONTRACT", "CI_STATIC_CONTRACT"}
NON_T2_CLASSES = {"ACCEPTANCE_E2E", "BROWSER_E2E", "PLATFORM_RELEASE_EVIDENCE", "UNKNOWN_EXECUTABLE"}


def load(path: Path) -> dict[str, Any]:
    try:
        value = json.loads(path.read_text(encoding="utf-8"))
    except Exception as exc:
        raise SystemExit(f"FAIL: cannot read {path.relative_to(ROOT)}: {exc}") from exc
    if not isinstance(value, dict):
        raise SystemExit(f"FAIL: {path.relative_to(ROOT)} must contain a JSON object")
    return value


def git_blob_sha(path: Path) -> str:
    data = path.read_bytes()
    return hashlib.sha1(f"blob {len(data)}\0".encode("utf-8") + data).hexdigest()


def classify_evidence(path_text: str) -> str:
    p = path_text.replace("\\", "/").strip()
    lower = p.lower()
    name = lower.rsplit("/", 1)[-1]
    if lower.startswith("tests/acceptance/"):
        return "ACCEPTANCE_E2E"
    if lower.startswith("packaging/") or lower.startswith("tools/release/"):
        return "PLATFORM_RELEASE_EVIDENCE"
    if lower.startswith("tools/ci/") and name.endswith(("_gate.py", "_contract.py", "_contract.js")):
        return "CI_STATIC_CONTRACT"
    if "browser" in name and name.endswith((".py", ".js")):
        return "BROWSER_E2E"
    if name.endswith("_test.go"):
        return "GO_UNIT_PACKAGE"
    if lower.startswith("tests/renderer/") and name.endswith(".js"):
        return "RENDERER_NODE_CONTRACT"
    if lower.startswith("tools/ci/") and name.endswith(("_test.py", "_test.js")):
        return "CI_STATIC_CONTRACT"
    return "UNKNOWN_EXECUTABLE"


def parent_for(item: dict[str, Any], reconciliation: dict[str, Any]) -> str:
    item_id = str(item.get("id") or "").strip()
    category = str(item.get("category") or "").strip()
    overrides = reconciliation.get("scanParentOverrides") or {}
    defaults = reconciliation.get("scanParentByCategory") or {}
    return str(overrides.get(item_id) or defaults.get(category) or "").strip()


def reconstruct_effective(
    ledger: dict[str, Any],
    scan1: dict[str, Any],
    scan2: dict[str, Any],
    reconciliation: dict[str, Any],
    errors: list[str],
) -> dict[str, dict[str, Any]]:
    features = ledger.get("features")
    if not isinstance(features, list) or not features:
        errors.append("frozen ledger features must be a non-empty array")
        return {}

    physical: dict[str, dict[str, Any]] = {}
    for row in features:
        if not isinstance(row, dict):
            errors.append("frozen ledger contains a non-object feature")
            continue
        fid = str(row.get("id") or "").strip()
        if not fid:
            errors.append("frozen ledger contains feature with empty id")
            continue
        if fid in physical:
            errors.append(f"duplicate frozen feature id: {fid}")
            continue
        physical[fid] = row

    excluded = {
        str(item.get("id") or "").strip()
        for item in (reconciliation.get("excludedFutureSourceCarryForward") or [])
        if isinstance(item, dict)
    }
    effective = {fid: dict(row) for fid, row in physical.items() if fid not in excluded}

    for scan_name, scan in ((SCAN1.name, scan1), (SCAN2.name, scan2)):
        omissions = scan.get("omissionsFound")
        if not isinstance(omissions, list):
            errors.append(f"{scan_name} omissionsFound must be an array")
            continue
        for item in omissions:
            if not isinstance(item, dict):
                errors.append(f"{scan_name} contains a non-object responsibility")
                continue
            item_id = str(item.get("id") or "").strip()
            if not item_id:
                errors.append(f"{scan_name} contains a responsibility with empty id")
                continue
            if item_id in effective:
                errors.append(f"effective responsibility duplicated: {item_id}")
                continue
            parent_id = parent_for(item, reconciliation)
            parent = physical.get(parent_id)
            if parent is None:
                errors.append(f"{scan_name}:{item_id} has no valid canonical parent ({parent_id!r})")
                continue
            tests = item.get("tests")
            if not isinstance(tests, list) or not tests:
                tests = parent.get("existingRegressionOwners") or []
            effective[item_id] = {
                "id": item_id,
                "name": item_id,
                "category": item.get("category"),
                "assuranceProfile": parent.get("assuranceProfile"),
                "existingRegressionOwners": list(tests),
                "_parentFeature": parent_id,
                "_scan": scan_name,
            }
    return effective


def main() -> int:
    parser = argparse.ArgumentParser()
    parser.add_argument("--strict", action="store_true")
    args = parser.parse_args()

    ledger, freeze = load(LEDGER), load(FREEZE)
    scan1, scan2, reconciliation = load(SCAN1), load(SCAN2), load(RECONCILIATION)
    t2, current_state, closure = load(T2), load(CURRENT_STATE), load(CLOSURE)
    errors: list[str] = []

    product = current_state.get("productCapabilityGate") or {}
    governed = product.get("nextGovernedTracks") or []
    t2_runtime_state = next((str(x.get("status") or "") for x in governed if isinstance(x, dict) and x.get("track") == "T2"), "")
    t3_runtime_state = next((str(x.get("status") or "") for x in governed if isinstance(x, dict) and x.get("track") == "T3"), "")
    assurance_state = str(t2.get("state") or "")

    if assurance_state == "IN_PROGRESS":
        if product.get("nextChildIssue") != 115 or product.get("nextChildTrack") != "T2" or t2_runtime_state != "IN_PROGRESS":
            errors.append("IN_PROGRESS T2 must be the active T2/#115 child")
        if t3_runtime_state != "NOT_STARTED":
            errors.append("IN_PROGRESS T2 must not silently start T3/#116")
    elif assurance_state == "COMPLETE":
        completed = product.get("completedChildTracks") or []
        t2_completed = next((x for x in completed if isinstance(x, dict) and x.get("track") == "T2" and x.get("issue") == 115), None)
        if not isinstance(t2_completed, dict) or t2_completed.get("status") != "COMPLETE":
            errors.append("COMPLETE T2 must remain recorded in completedChildTracks")
        if not str((t2_completed or {}).get("mergedCommitSha") or "").strip():
            errors.append("COMPLETE T2 requires durable mergedCommitSha evidence")
    else:
        errors.append(f"unsupported T2 state: {assurance_state!r}")

    gaps = closure.get("gaps") or []
    t1_gap = next((x for x in gaps if isinstance(x, dict) and x.get("id") == "T1-FEATURE-TRACEABILITY"), None)
    t2_gap = next((x for x in gaps if isinstance(x, dict) and x.get("id") == "T2-UNIT-CONTRACT-PROPERTY"), None)
    if not isinstance(t1_gap, dict) or t1_gap.get("status") != "VERIFIED":
        errors.append("T2 requires T1-FEATURE-TRACEABILITY=VERIFIED")
    expected_t2_gap_state = "VERIFIED" if assurance_state == "COMPLETE" else "IMPLEMENTED_UNVERIFIED"
    if not isinstance(t2_gap, dict) or t2_gap.get("status") != expected_t2_gap_state:
        errors.append(f"T2 closure row must be {expected_t2_gap_state} while assurance state is {assurance_state}")

    if "python3 tools/ci/v18_t2_unit_contract_assurance_gate.py" not in CI_FAST.read_text(encoding="utf-8"):
        errors.append("T2 assurance gate is not bound into canonical CI Fast")

    actual_blob = git_blob_sha(LEDGER)
    expected_blob = str((freeze.get("frozenDiscovery") or {}).get("gitBlobSha") or "").strip()
    if freeze.get("state") != "FROZEN_T1" or expected_blob != actual_blob:
        errors.append(f"frozen T1 ledger blob mismatch: expected {expected_blob}, got {actual_blob}")
    if t2.get("trackIssue") != 115 or t2.get("programIssue") != 113 or t2.get("frozenT1GitBlobSha") != actual_blob:
        errors.append("T2 assurance identity/frozen T1 binding mismatch")

    effective = reconstruct_effective(ledger, scan1, scan2, reconciliation, errors)
    expected_count = int((freeze.get("effectiveInventory") or {}).get("effectiveShippedV18Responsibilities") or 0)
    if len(effective) != expected_count:
        errors.append(f"effective inventory count mismatch: expected {expected_count}, got {len(effective)}")

    supplements = t2.get("evidenceSupplements") or {}
    if not isinstance(supplements, dict):
        errors.append("T2 evidenceSupplements must be an object")
        supplements = {}
    for fid, spec in supplements.items():
        if fid not in effective or not isinstance(spec, dict):
            errors.append(f"invalid T2 evidence supplement: {fid}")
            continue
        owners = spec.get("owners")
        rationale = str(spec.get("rationale") or "").strip()
        if not isinstance(owners, list) or not owners or not rationale:
            errors.append(f"T2 evidence supplement {fid} requires owners and rationale")
            continue
        row_owners = effective[fid].setdefault("existingRegressionOwners", [])
        for owner in owners:
            if isinstance(owner, str) and owner.strip() and owner.strip() not in row_owners:
                row_owners.append(owner.strip())

    profiles = ledger.get("assuranceProfiles") or {}
    uncovered: list[str] = []
    missing_paths: list[str] = []
    class_counts: dict[str, int] = {key: 0 for key in sorted(T2_CLASSES | NON_T2_CLASSES)}
    covered = 0
    for fid, row in sorted(effective.items()):
        profile_name = str(row.get("assuranceProfile") or "").strip()
        profile = profiles.get(profile_name)
        if not isinstance(profile, dict) or not str(profile.get("T2") or "").strip():
            errors.append(f"{fid}: missing T2 expectation/profile")
            continue
        valid = []
        for owner in row.get("existingRegressionOwners") or []:
            owner_text = str(owner or "").strip()
            if not owner_text:
                continue
            path = ROOT / owner_text
            if not path.is_file():
                missing_paths.append(f"{fid}:{owner_text}")
                continue
            cls = classify_evidence(owner_text)
            class_counts[cls] = class_counts.get(cls, 0) + 1
            if cls in T2_CLASSES:
                valid.append(owner_text)
        if valid:
            covered += 1
        else:
            uncovered.append(fid)

    if missing_paths:
        errors.append("referenced regression evidence missing at current head: " + ", ".join(sorted(missing_paths)))

    strict = args.strict or assurance_state == "COMPLETE"
    actual_gaps = set(uncovered)
    declared_gaps = set(str(x) for x in (t2.get("knownCoverageGaps") or []))
    if assurance_state == "IN_PROGRESS":
        declaration = str(t2.get("gapDeclarationState") or "")
        if declaration not in {"PENDING_FIRST_EXECUTION", "CURRENT"}:
            errors.append(f"unsupported T2 gapDeclarationState: {declaration!r}")
        if declaration == "CURRENT":
            if actual_gaps - declared_gaps:
                errors.append("T2 found undeclared gaps: " + ", ".join(sorted(actual_gaps - declared_gaps)))
            if declared_gaps - actual_gaps:
                errors.append("T2 declared stale/resolved gaps: " + ", ".join(sorted(declared_gaps - actual_gaps)))
            if t2.get("uncoveredResponsibilityCount") != len(uncovered):
                errors.append("T2 uncoveredResponsibilityCount drift")
    if strict and uncovered:
        errors.append("T2 strict closure has uncovered responsibilities: " + ", ".join(uncovered))
    if assurance_state == "COMPLETE":
        if t2.get("uncoveredResponsibilityCount") != 0:
            errors.append("T2 COMPLETE requires uncoveredResponsibilityCount=0")
        if t2.get("exactHeadQualificationRequired") is not True:
            errors.append("T2 COMPLETE must require exact-head qualification")

    print("V18 T2 UNIT / CONTRACT / STATIC / PROPERTY ASSURANCE")
    print(f"frozen T1 blob: {actual_blob}")
    print(f"effective responsibilities: {len(effective)}")
    print(f"T2-covered responsibilities: {covered}")
    print(f"T2-uncovered responsibilities: {len(uncovered)}")
    print(f"T2-specific evidence supplements: {len(supplements)}")
    for cls in sorted(class_counts):
        print(f"{cls}: {class_counts[cls]}")
    if uncovered:
        print("uncovered ids: " + ", ".join(uncovered))
    print("acceptance/browser/platform-only evidence cannot satisfy T2: PASS")
    print("completed T2 remains enforceable after governed handoff: PASS")
    print("T3-T10 certification is not implied by T2: PASS")

    if errors:
        print("V18 T2 ASSURANCE GATE: FAIL", file=sys.stderr)
        for error in errors:
            print(" - " + error, file=sys.stderr)
        return 1
    print(f"V18 T2 ASSURANCE GATE: PASS (strict={strict})")
    return 0


if __name__ == "__main__":
    sys.exit(main())
