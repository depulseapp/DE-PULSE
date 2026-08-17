#!/usr/bin/env python3
"""Validate the DE.PULSE v18.5.1 G1 immutable recovery-scope overlay."""
from __future__ import annotations

import argparse
import hashlib
import json
import sys
from collections import Counter
from pathlib import Path

ROOT = Path(__file__).resolve().parents[2]
LEDGER = ROOT / "release" / "v18.5.1" / "V17-V18-IMPLEMENTATION-RECONCILIATION.json"
MAP = ROOT / "release" / "v18.5.1" / "V18.5.1-G1-PLACEMENT-MAP.json"
FREEZE = ROOT / "release" / "v18.5.1" / "G1-SCOPE-FREEZE.json"
REPORT = ROOT / "evidence" / "v18.5.1-G1-scope-diagnostics.json"
EXPECTED_ROWS = 295
EXPECTED_IDS = [
    "COPY-18.5.1-001", "SYMBOL-18.5.1-001", "SYMBOL-18.5.1-002", "NAV-18.5.1-001",
    "RESEARCH-v15.1.0-17-19-REOPENED", "VERSION-18.5.1-002", "HOVER-18.5.1-001",
    "AUDIT-18-UI-001", "AUDIT-18-CI-001", "AUDIT-18-QA-001",
]


def collect(node, path="$"):
    rows = []
    if isinstance(node, dict):
        if isinstance(node.get("id"), str) and isinstance(node.get("status"), str):
            rows.append((path, node))
        for key, value in node.items():
            rows.extend(collect(value, f"{path}.{key}"))
    elif isinstance(node, list):
        for idx, value in enumerate(node):
            rows.extend(collect(value, f"{path}[{idx}]"))
    return rows


def sha256(path):
    return hashlib.sha256(path.read_bytes()).hexdigest()


def main():
    parser = argparse.ArgumentParser()
    parser.add_argument("--diagnose", action="store_true")
    args = parser.parse_args()
    errors = []
    ledger = json.loads(LEDGER.read_text())
    baseline = ledger.get("baseline", {})
    current = ledger.get("adaptiveReleaseTrain", {}).get("currentSlice", {})
    rows = collect(ledger)
    ids = [row["id"] for _, row in rows]
    counts = Counter(ids)
    duplicates = sorted(k for k, v in counts.items() if v > 1)
    if len(ids) != EXPECTED_ROWS or len(set(ids)) != EXPECTED_ROWS:
        errors.append(f"ledger conservation mismatch rows={len(ids)} unique={len(set(ids))}")
    if duplicates:
        errors.append(f"duplicate IDs: {duplicates[:20]}")
    if baseline.get("totalTrackedRows") != EXPECTED_ROWS:
        errors.append("declared tracked-row count drift")
    if baseline.get("currentStableTag") != "v18.5.0-stable" or baseline.get("currentStableCommit") != "0d37ca35f5fc3ad89cebed506cc5a4c2d6a7a680":
        errors.append("incoming Stable baseline drift")
    if baseline.get("reconciliationBranch") != "v18.5.1-development":
        errors.append("reconciliation branch drift")
    if current.get("assignedIds") != EXPECTED_IDS:
        errors.append("audited ten-ID current slice drift")

    placement = None
    freeze = None
    if not MAP.exists():
        errors.append("missing G1 placement map")
    else:
        placement = json.loads(MAP.read_text())
    if not FREEZE.exists():
        errors.append("missing G1 freeze manifest")
    else:
        freeze = json.loads(FREEZE.read_text())

    if placement:
        placed = placement.get("placements", [])
        placed_ids = [item.get("id") for item in placed]
        if placement.get("trackedRows") != EXPECTED_ROWS or len(placed) != EXPECTED_ROWS or set(placed_ids) != set(ids):
            errors.append("placement map does not conserve all 295 ledger IDs")
        if len(set(placed_ids)) != len(placed_ids):
            errors.append("placement map contains duplicate IDs")
        current_map = [item.get("id") for item in placed if item.get("disposition") == "CURRENT_V18.5.1"]
        if set(current_map) != set(EXPECTED_IDS) or len(current_map) != len(EXPECTED_IDS):
            errors.append("placement map current-slice disposition is not exactly the ten frozen IDs")
        stray_current = [item.get("id") for item in placed if item.get("release") == "v18.5.1" and item.get("id") not in EXPECTED_IDS]
        if stray_current:
            errors.append(f"non-scope IDs assigned to v18.5.1: {stray_current[:20]}")
        for item in placed:
            for field in ("id", "path", "status", "disposition", "release", "owner", "lane", "regression", "impact", "reason"):
                if not item.get(field):
                    errors.append(f"placement {item.get('id')} missing {field}")
                    break
        if placement.get("sourceLedgerSha256") != sha256(LEDGER):
            errors.append("placement map source-ledger fingerprint mismatch")

    if freeze:
        if freeze.get("gate") != "G1" or freeze.get("freezeState") != "FROZEN" or freeze.get("decision") != "G1_SCOPE_FROZEN":
            errors.append("freeze manifest is not an explicit G1 FROZEN decision")
        if freeze.get("currentSliceIds") != EXPECTED_IDS:
            errors.append("freeze manifest ten-ID scope mismatch")
        src = freeze.get("sourceLedger", {})
        pm = freeze.get("placementMap", {})
        if src.get("sha256") != sha256(LEDGER) or src.get("trackedRows") != EXPECTED_ROWS:
            errors.append("freeze manifest ledger identity mismatch")
        if MAP.exists() and (pm.get("sha256") != sha256(MAP) or pm.get("rows") != EXPECTED_ROWS):
            errors.append("freeze manifest placement-map identity mismatch")
        if freeze.get("auditIntegrationCommit") != "664418e143a969b63cd2169616278dd54e501d6b":
            errors.append("audit integration baseline drift")

    report = {
        "schema": "DE.PULSE-v18.5.1-G1-SCOPE-DIAGNOSTICS-2",
        "release": "v18.5.1",
        "ledgerRows": len(ids),
        "uniqueLedgerRows": len(set(ids)),
        "duplicateIds": duplicates,
        "ledgerPlanningState": current.get("planningState"),
        "ledgerScopeFreezeState": current.get("scopeFreezeState"),
        "freezeOverlayPresent": FREEZE.exists(),
        "placementMapPresent": MAP.exists(),
        "errors": errors,
        "decision": "G1_SCOPE_FROZEN" if not errors else "G1_BLOCKED",
        "note": "The original audit ledger remains the immutable input snapshot; G1 freeze is recorded by the hash-bound overlay manifest and complete placement map.",
    }
    REPORT.parent.mkdir(parents=True, exist_ok=True)
    REPORT.write_text(json.dumps(report, indent=2) + "\n")
    print(f"G1 scope gate: rows={len(ids)} unique={len(set(ids))} overlay={FREEZE.exists()} map={MAP.exists()} errors={len(errors)}")
    for error in errors[:40]:
        print("ERROR:", error)
    if errors and not args.diagnose:
        print("G1 Scope Conservation: FAIL")
        return 1
    if errors:
        print("G1 Scope Conservation: DIAGNOSTIC BLOCKED")
    else:
        print("G1 Scope Conservation: PASS · 295/295 conserved · exact ten-ID v18.5.1 scope · placement/ownership map hash-bound")
    return 0


if __name__ == "__main__":
    sys.exit(main())
