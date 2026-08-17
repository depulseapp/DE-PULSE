#!/usr/bin/env python3
"""DE.PULSE v18.5.1 G1 scope-conservation diagnostic/freeze gate.

This is an owning check inside existing G1. It does not create a new top-level gate.
The gate deliberately fails while the reconciliation ledger is still unplaced or
unowned, so expensive qualification cannot start from an implicitly moving scope.
"""
from __future__ import annotations

import argparse
import json
import sys
from collections import Counter
from pathlib import Path

ROOT = Path(__file__).resolve().parents[2]
LEDGER = ROOT / "release" / "v18.5.1" / "V17-V18-IMPLEMENTATION-RECONCILIATION.json"
REPORT = ROOT / "evidence" / "v18.5.1-G1-scope-diagnostics.json"

EXPECTED_RELEASE = "v18.5.1"
EXPECTED_BRANCH = "v18.5.1-development"
EXPECTED_BASELINE_TAG = "v18.5.0-stable"
EXPECTED_BASELINE_COMMIT = "0d37ca35f5fc3ad89cebed506cc5a4c2d6a7a680"
EXPECTED_ROWS = 295
EXPECTED_CURRENT_IDS = [
    "COPY-18.5.1-001",
    "SYMBOL-18.5.1-001",
    "SYMBOL-18.5.1-002",
    "NAV-18.5.1-001",
    "RESEARCH-v15.1.0-17-19-REOPENED",
    "VERSION-18.5.1-002",
    "HOVER-18.5.1-001",
    "AUDIT-18-UI-001",
    "AUDIT-18-CI-001",
    "AUDIT-18-QA-001",
]
FINAL_STATUSES = {"FRESH_PASS", "INTENTIONALLY_SUPERSEDED", "NOT_APPLICABLE"}
FUTURE_STATUSES = {"ROADMAP_PLACED_FUTURE"}
OPEN_STATUSES = {
    "REVALIDATION_REQUIRED",
    "REVALIDATION_REQUIRED_HIGH",
    "REOPENED",
    "NOT_IMPLEMENTED_CONFIRMED",
    "PLACED_NEXT_V18",
    "ROADMAP_BOUND_FOUNDATION_OPEN",
}


def collect_rows(node, path="$"):
    out = []
    if isinstance(node, dict):
        if isinstance(node.get("id"), str) and isinstance(node.get("status"), str):
            out.append((path, node))
        for key, value in node.items():
            out.extend(collect_rows(value, f"{path}.{key}"))
    elif isinstance(node, list):
        for idx, value in enumerate(node):
            out.extend(collect_rows(value, f"{path}[{idx}]"))
    return out


def has_future_placement(row):
    if row.get("status") in FUTURE_STATUSES:
        return bool(row.get("placement") or row.get("candidateRelease"))
    classification = row.get("classification")
    candidate = row.get("candidateRelease") or row.get("completionLane") or row.get("placement")
    return classification in {"PLACED_NEXT_V18", "FINAL_CLOSURE_BLOCKER"} and bool(candidate)


def source_owner(row):
    for key in ("owner", "sourceOwner", "codeOwner", "currentSourceOwner"):
        if row.get(key):
            return row.get(key)
    return None


def main():
    parser = argparse.ArgumentParser()
    parser.add_argument("--diagnose", action="store_true", help="write diagnostics and exit 0")
    args = parser.parse_args()

    data = json.loads(LEDGER.read_text())
    baseline = data.get("baseline", {})
    train = data.get("adaptiveReleaseTrain", {})
    current = train.get("currentSlice", {})
    rows = collect_rows(data)
    row_ids = [row["id"] for _, row in rows]
    counts = Counter(row_ids)
    duplicates = sorted(k for k, v in counts.items() if v > 1)
    unique = {row["id"]: row for _, row in rows}

    errors = []
    warnings = []
    if data.get("release") != EXPECTED_RELEASE:
        errors.append(f"ledger release is {data.get('release')!r}, expected {EXPECTED_RELEASE}")
    if baseline.get("currentStableTag") != EXPECTED_BASELINE_TAG:
        errors.append("incoming Stable tag drift")
    if baseline.get("currentStableCommit") != EXPECTED_BASELINE_COMMIT:
        errors.append("incoming Stable commit drift")
    if baseline.get("reconciliationBranch") != EXPECTED_BRANCH:
        errors.append("reconciliation branch drift")
    if baseline.get("totalTrackedRows") != EXPECTED_ROWS:
        errors.append(f"declared tracked-row count is {baseline.get('totalTrackedRows')}, expected {EXPECTED_ROWS}")
    if duplicates:
        errors.append(f"duplicate tracked IDs: {duplicates[:20]}")
    missing_current = [rid for rid in EXPECTED_CURRENT_IDS if rid not in unique]
    if missing_current:
        errors.append(f"current-slice IDs missing from ledger: {missing_current}")
    assigned = current.get("assignedIds", [])
    if assigned != EXPECTED_CURRENT_IDS:
        errors.append("adaptiveReleaseTrain.currentSlice.assignedIds differs from the audited ten-ID recovery scope")

    unplaced = []
    missing_owner = []
    unknown_status = []
    allowed = set(data.get("allowedStatuses", [])) | {"ROADMAP_BOUND_FOUNDATION_OPEN"}
    for path, row in rows:
        rid = row["id"]
        status = row["status"]
        if status not in allowed:
            unknown_status.append({"id": rid, "status": status, "path": path})
        if rid in EXPECTED_CURRENT_IDS:
            if source_owner(row) is None:
                missing_owner.append({"id": rid, "status": status, "path": path})
            continue
        if status in FINAL_STATUSES:
            continue
        if status in FUTURE_STATUSES and has_future_placement(row):
            continue
        if status in OPEN_STATUSES:
            if not has_future_placement(row):
                unplaced.append({"id": rid, "status": status, "path": path})
            if source_owner(row) is None:
                missing_owner.append({"id": rid, "status": status, "path": path})

    declared_count = baseline.get("totalTrackedRows")
    if len(unique) != declared_count:
        warnings.append(
            f"recursive id/status discovery found {len(unique)} unique rows vs declared {declared_count}; "
            "review section accounting before freeze"
        )

    frozen = current.get("planningState") == "G1_SCOPE_FROZEN" and current.get("scopeFreezeState") == "FROZEN"
    report = {
        "schema": "DE.PULSE-v18.5.1-G1-SCOPE-DIAGNOSTICS-1",
        "release": EXPECTED_RELEASE,
        "declaredTrackedRows": declared_count,
        "discoveredUniqueIdStatusRows": len(unique),
        "assignedIds": assigned,
        "expectedAssignedIds": EXPECTED_CURRENT_IDS,
        "planningState": current.get("planningState"),
        "scopeFreezeState": current.get("scopeFreezeState"),
        "frozenStateRecorded": frozen,
        "duplicateIds": duplicates,
        "missingCurrentIds": missing_current,
        "unplacedOpenRows": unplaced,
        "missingCurrentOwnerRows": missing_owner,
        "unknownStatusRows": unknown_status,
        "errors": errors,
        "warnings": warnings,
        "decision": "G1_FREEZE_ELIGIBLE" if not errors and not unplaced and not missing_owner and frozen else "G1_BLOCKED",
    }
    REPORT.parent.mkdir(parents=True, exist_ok=True)
    REPORT.write_text(json.dumps(report, indent=2) + "\n")

    print(f"G1 diagnostics: uniqueRows={len(unique)} declared={declared_count} duplicates={len(duplicates)} unplaced={len(unplaced)} missingOwner={len(missing_owner)} frozen={frozen}")
    for item in errors[:20]:
        print("ERROR:", item)
    for item in warnings[:20]:
        print("WARN:", item)
    if unplaced:
        print("First unplaced IDs:", ", ".join(x["id"] for x in unplaced[:30]))
    if missing_owner:
        print("First ownership gaps:", ", ".join(x["id"] for x in missing_owner[:30]))

    if args.diagnose:
        return 0
    if report["decision"] != "G1_FREEZE_ELIGIBLE":
        print("G1 Scope Conservation: FAIL")
        return 1
    print("G1 Scope Conservation: PASS")
    return 0


if __name__ == "__main__":
    sys.exit(main())
