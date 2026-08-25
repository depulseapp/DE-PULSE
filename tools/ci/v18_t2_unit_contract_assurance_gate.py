#!/usr/bin/env python3
"""Deterministic T2 unit/contract/static/property assurance audit.

T2 does not certify functional/E2E behavior (T3), adverse-state/data-truth
behavior (T4), persistence lifecycle (T5), security (T6), UI quality (T7),
performance (T8), packaged platforms (T9), or final closure (T10).

The audit reconstructs the frozen T1 effective inventory and requires every
shipped-v18 responsibility to have a non-empty T2 expectation and at least one
current executable unit/contract/static evidence owner. Browser/acceptance/E2E
evidence is intentionally not accepted as a T2 substitute.
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

T2_CLASSES = {
    "GO_UNIT_PACKAGE",
    "RENDERER_NODE_CONTRACT",
    "CI_STATIC_CONTRACT",
}
NON_T2_CLASSES = {
    "ACCEPTANCE_E2E",
    "BROWSER_E2E",
    "PLATFORM_RELEASE_EVIDENCE",
    "UNKNOWN_EXECUTABLE",
}


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
    header = f"blob {len(data)}\0".encode("utf-8")
    return hashlib.sha1(header + data).hexdigest()


def classify_evidence(path_text: str) -> str:
    p = path_text.replace("\\", "/").strip()
    name = p.rsplit("/", 1)[-1].lower()
    lower = p.lower()

    if lower.startswith("tests/acceptance/"):
        return "ACCEPTANCE_E2E"
    if lower.startswith("packaging/") or lower.startswith("tools/release/"):
        return "PLATFORM_RELEASE_EVIDENCE"
    # Static gates/contracts are T2 evidence even when their subject is browser
    # routing. Executable browser harnesses remain browser evidence below.
    if lower.startswith("tools/ci/") and (
        name.endswith("_gate.py")
        or name.endswith("_contract.py")
        or name.endswith("_contract.js")
    ):
        return "CI_STATIC_CONTRACT"
    if "browser" in name and (name.endswith(".py") or name.endswith(".js")):
        return "BROWSER_E2E"
    if name.endswith("_test.go"):
        return "GO_UNIT_PACKAGE"
    if lower.startswith("tests/renderer/") and name.endswith(".js"):
        return "RENDERER_NODE_CONTRACT"
    if lower.startswith("tools/ci/") and (
        name.endswith("_test.py") or name.endswith("_test.js")
    ):
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
    parser.add_argument(
        "--strict",
        action="store_true",
        help="Fail if any effective responsibility lacks acceptable T2 evidence.",
    )
    args = parser.parse_args()

    ledger = load(LEDGER)
    freeze = load(FREEZE)
    scan1 = load(SCAN1)
    scan2 = load(SCAN2)
    reconciliation = load(RECONCILIATION)
    t2 = load(T2)
    current_state = load(CURRENT_STATE)
    closure = load(CLOSURE)
    errors: list[str] = []

    product = current_state.get("productCapabilityGate") or {}
    governed = product.get("nextGovernedTracks") or []
    t2_state = next((str(item.get("status") or "") for item in governed if isinstance(item, dict) and item.get("track") == "T2"), "")
    t3_state = next((str(item.get("status") or "") for item in governed if isinstance(item, dict) and item.get("track") == "T3"), "")
    if product.get("nextChildIssue") != 115 or product.get("nextChildTrack") != "T2" or t2_state != "IN_PROGRESS":
        errors.append("current-state must identify T2/#115 as the active IN_PROGRESS child")
    if t3_state != "NOT_STARTED":
        errors.append("T2 audit must not silently start T3/#116")

    gaps = closure.get("gaps") or []
    t1_gap = next((item for item in gaps if isinstance(item, dict) and item.get("id") == "T1-FEATURE-TRACEABILITY"), None)
    t2_gap = next((item for item in gaps if isinstance(item, dict) and item.get("id") == "T2-UNIT-CONTRACT-PROPERTY"), None)
    if not isinstance(t1_gap, dict) or t1_gap.get("status") != "VERIFIED":
        errors.append("T2 requires parent closure ledger T1-FEATURE-TRACEABILITY=VERIFIED")
    if not isinstance(t2_gap, dict) or t2_gap.get("status") not in {"IMPLEMENTED_UNVERIFIED", "VERIFIED"}:
        errors.append("parent closure ledger must record active T2 as IMPLEMENTED_UNVERIFIED or VERIFIED")

    ci_fast = CI_FAST.read_text(encoding="utf-8")
    ci_token = "python3 tools/ci/v18_t2_unit_contract_assurance_gate.py"
    if ci_token not in ci_fast:
        errors.append("T2 assurance gate is not bound into canonical CI Fast")

    if freeze.get("state") != "FROZEN_T1":
        errors.append("T2 requires frozen T1 manifest state FROZEN_T1")
    frozen_discovery = freeze.get("frozenDiscovery") or {}
    expected_blob = str(frozen_discovery.get("gitBlobSha") or "").strip()
    actual_blob = git_blob_sha(LEDGER)
    if expected_blob != actual_blob:
        errors.append(f"frozen T1 ledger blob mismatch: expected {expected_blob}, got {actual_blob}")

    if t2.get("trackIssue") != 115 or t2.get("programIssue") != 113:
        errors.append("T2 assurance identity must bind program #113 / track #115")
    if t2.get("frozenT1GitBlobSha") != actual_blob:
        errors.append("T2 assurance contract is not bound to the frozen T1 discovery blob")
    if t2.get("state") not in {"IN_PROGRESS", "COMPLETE"}:
        errors.append(f"unsupported T2 state: {t2.get('state')!r}")

    effective = reconstruct_effective(ledger, scan1, scan2, reconciliation, errors)
    expected_count = int((freeze.get("effectiveInventory") or {}).get("effectiveShippedV18Responsibilities") or 0)
    if len(effective) != expected_count:
        errors.append(f"effective inventory count mismatch: expected {expected_count}, got {len(effective)}")

    supplements = t2.get("evidenceSupplements") or {}
    if not isinstance(supplements, dict):
        errors.append("T2 evidenceSupplements must be an object")
        supplements = {}
    for fid, spec in supplements.items():
        if fid not in effective:
            errors.append(f"T2 evidence supplement references unknown effective responsibility: {fid}")
            continue
        if not isinstance(spec, dict):
            errors.append(f"T2 evidence supplement {fid} must be an object")
            continue
        rationale = str(spec.get("rationale") or "").strip()
        owners = spec.get("owners")
        if not rationale:
            errors.append(f"T2 evidence supplement {fid} requires a material rationale")
        if not isinstance(owners, list) or not owners or not all(isinstance(x, str) and x.strip() for x in owners):
            errors.append(f"T2 evidence supplement {fid} requires non-empty owners")
            continue
        row_owners = effective[fid].setdefault("existingRegressionOwners", [])
        for owner in owners:
            owner_text = owner.strip()
            if owner_text not in row_owners:
                row_owners.append(owner_text)

    profiles = ledger.get("assuranceProfiles") or {}
    uncovered: list[str] = []
    unknown_paths: list[str] = []
    class_counts: dict[str, int] = {key: 0 for key in sorted(T2_CLASSES | NON_T2_CLASSES)}
    covered = 0

    for fid, row in sorted(effective.items()):
        profile_name = str(row.get("assuranceProfile") or "").strip()
        profile = profiles.get(profile_name)
        if not isinstance(profile, dict):
            errors.append(f"{fid}: unknown assurance profile {profile_name!r}")
            continue
        expectation = str(profile.get("T2") or "").strip()
        if not expectation:
            errors.append(f"{fid}: missing T2 expectation")
            continue

        owners = row.get("existingRegressionOwners")
        if not isinstance(owners, list) or not owners:
            uncovered.append(fid)
            continue

        valid_for_t2: list[str] = []
        for owner in owners:
            owner_text = str(owner or "").strip()
            if not owner_text:
                continue
            path = ROOT / owner_text
            if not path.is_file():
                unknown_paths.append(f"{fid}:{owner_text}")
                continue
            cls = classify_evidence(owner_text)
            class_counts[cls] = class_counts.get(cls, 0) + 1
            if cls in T2_CLASSES:
                valid_for_t2.append(owner_text)

        if valid_for_t2:
            covered += 1
        else:
            uncovered.append(fid)

    if unknown_paths:
        errors.append("referenced regression evidence missing at current head: " + ", ".join(sorted(unknown_paths)))

    strict = args.strict or t2.get("state") == "COMPLETE"
    declared_gaps = set(str(x) for x in (t2.get("knownCoverageGaps") or []))
    actual_gaps = set(uncovered)
    if t2.get("state") == "IN_PROGRESS":
        declaration_state = str(t2.get("gapDeclarationState") or "").strip()
        if declaration_state not in {"PENDING_FIRST_EXECUTION", "CURRENT"}:
            errors.append(f"unsupported T2 gapDeclarationState: {declaration_state!r}")
        if declaration_state == "CURRENT":
            undeclared = sorted(actual_gaps - declared_gaps)
            stale = sorted(declared_gaps - actual_gaps)
            if undeclared:
                errors.append("T2 audit found undeclared coverage gaps: " + ", ".join(undeclared))
            if stale:
                errors.append("T2 knownCoverageGaps contains stale/resolved ids: " + ", ".join(stale))
            if t2.get("uncoveredResponsibilityCount") != len(uncovered):
                errors.append(
                    "T2 uncoveredResponsibilityCount drift: "
                    f"declared={t2.get('uncoveredResponsibilityCount')!r} actual={len(uncovered)}"
                )
    if strict and uncovered:
        errors.append("T2 strict closure has uncovered responsibilities: " + ", ".join(uncovered))
    if t2.get("state") == "COMPLETE":
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
