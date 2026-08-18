#!/usr/bin/env python3
"""Inventory and release gate for complete v17/v18 implementation reconciliation."""

from __future__ import annotations

import argparse
import json
import re
import sys
from pathlib import Path

ROOT = Path(__file__).resolve().parent
LEDGER = ROOT / "release" / "v18.5.1" / "V17-V18-IMPLEMENTATION-RECONCILIATION.json"

BLOCKING = {
    "REVALIDATION_REQUIRED",
    "REVALIDATION_REQUIRED_HIGH",
    "REOPENED",
    "NOT_IMPLEMENTED_CONFIRMED",
    "PLACED_NEXT_V18",
}
FINAL = {
    "FRESH_PASS",
    "INTENTIONALLY_SUPERSEDED",
    "NOT_APPLICABLE",
    "ROADMAP_PLACED_FUTURE",
}
REQUIRED_BLOCKERS = {
    "IMPL-18-TRADEINSIGHT-001",
    "CONVO-V18-003",
    "COPY-18.5.1-001",
    "SYMBOL-18.5.1-001",
    "SYMBOL-18.5.1-002",
    "NAV-18.5.1-001",
    "HOVER-18.5.1-001",
    "HEADER-18.5.1-001",
    "CI-ADAPTIVE-18.5.1-001",
    "RESEARCH-v15.1.0-17-19-REOPENED",
    "VERSION-18.5.1-002",
    "IMPL-18-UTILITY-001",
    "IMPL-18-UTILITY-002",
    "IMPL-18-UTILITY-003",
    "IMPL-18-UTILITY-004",
    "IMPL-18-DOC-001",
    "IMPL-17-DEPS-001",
}
RELEASE_SCOPE_FILES = [
    "v18_0_scope.json",
    "v18_0_1_scope.json",
    "v18_0_2_scope.json",
    "v18_0_3_scope.json",
    "v18_0_4_scope.json",
    "v18_0_5_scope.json",
    "v18_0_6_scope.json",
    "v18_1_scope.json",
    "v18_2_scope.json",
    "v18_3_scope.json",
    "v18_4_scope.json",
    "v18_5_scope.json",
]


def load(path: Path):
    return json.loads(path.read_text(encoding="utf-8"))


def slug(value: object) -> str:
    result = re.sub(r"[^A-Z0-9]+", "-", str(value).upper()).strip("-")
    return result[:80]


def scope_entries(path: str, data: dict) -> list[dict]:
    rows: list[dict] = []

    key_prefix = {
        "items": "items",
        "scope_lock": "scope-lock",
        "clauses": "clauses",
        "firstSlice": "first-slice",
        "closureDimensions": "closure",
        "adrGdiRequiredScenarios": "adr-gdi",
        "releaseBlockers": "release-blocker",
        "protectedContracts": "protected",
    }

    def append(section: str, values: object) -> None:
        if not isinstance(values, list):
            return
        for index, item in enumerate(values, 1):
            if isinstance(item, dict):
                key = item.get("id") or f"{path}#{key_prefix[section]}-{index}"
                title = item.get("name") or item.get("requirement") or str(item)
            else:
                key = f"{path}#{key_prefix[section]}-{index}"
                title = str(item)
            rows.append(
                {
                    "id": f"{slug(path.replace('.json', ''))}-{slug(key)}",
                    "title": title,
                    "section": section,
                }
            )

    append("items", data.get("items"))
    append("scope_lock", data.get("scope_lock"))
    append("clauses", data.get("clauses"))
    append("firstSlice", data.get("firstSlice"))
    append("closureDimensions", data.get("closureDimensions"))
    append("adrGdiRequiredScenarios", data.get("adrGdiRequiredScenarios"))
    append("releaseBlockers", data.get("releaseBlockers"))
    append("protectedContracts", data.get("protectedContracts"))
    return rows


def collect_status_rows(value: object, path: str = "ledger") -> list[tuple[str, str, dict]]:
    rows: list[tuple[str, str, dict]] = []
    if isinstance(value, dict):
        status = value.get("status")
        if isinstance(status, str):
            rows.append((path, status, value))
        for key, child in value.items():
            rows.extend(collect_status_rows(child, f"{path}.{key}"))
    elif isinstance(value, list):
        for index, child in enumerate(value):
            rows.extend(collect_status_rows(child, f"{path}[{index}]"))
    return rows


def main() -> int:
    parser = argparse.ArgumentParser()
    mode_group = parser.add_mutually_exclusive_group()
    mode_group.add_argument(
        "--slice-release",
        action="store_true",
        help="enforce the current evidence-selected v18.x slice; remaining work must have complete next-slice placement",
    )
    mode_group.add_argument(
        "--major-closure",
        action="store_true",
        help="enforce zero-gap final v18 major closure",
    )
    mode_group.add_argument(
        "--release",
        action="store_true",
        help="backward-compatible alias for --major-closure",
    )
    args = parser.parse_args()
    slice_mode = args.slice_release
    major_mode = args.major_closure or args.release

    errors: list[str] = []
    ledger = load(LEDGER)

    if ledger.get("schema") != "DE.PULSE-V18.5.1-V17-V18-IMPLEMENTATION-RECONCILIATION-1":
        errors.append("ledger schema drift")
    if ledger.get("release") != "v18.5.1":
        errors.append("ledger release drift")

    train = ledger.get("adaptiveReleaseTrain", {})
    current_slice = train.get("currentSlice", {})
    if train.get("model") != "EVIDENCE_SELECTED_V18_X":
        errors.append("adaptive release-train model missing or drifted")
    if current_slice.get("release") != ledger.get("release"):
        errors.append("current adaptive slice does not match ledger release")
    if current_slice.get("finalMajorClosure") is not False:
        errors.append("v18.5.1 must not be predeclared final major closure")
    if train.get("finalClosureDesignation") != "EVIDENCE_SELECTED_ONLY_AFTER_ZERO_GAP_READINESS":
        errors.append("final closure designation policy drift")


    inherited = load(ROOT / "renderer" / "qa" / "v15.1.2-approved-scope.json")
    inherited_rows = ledger.get("inheritedApprovedScope", {}).get("items", [])
    if inherited.get("count") != 48 or len(inherited.get("items", [])) != 48:
        errors.append("canonical inherited approved scope is not 48")
    if len(inherited_rows) != 48:
        errors.append("ledger does not inventory all 48 inherited requirements")
    inherited_origin = {row.get("originId") for row in inherited_rows}
    if inherited_origin != {item.get("id") for item in inherited.get("items", [])}:
        errors.append("inherited requirement IDs differ from canonical scope")

    v17_major = load(ROOT / "v17_major_scope.json")
    v17_matrix = load(ROOT / "v17_5_major_closure_scope_matrix.json")
    v17_rows = ledger.get("v17", {}).get("items", [])
    if len(v17_major.get("items", [])) != 20 or v17_matrix.get("scope_count") != 20:
        errors.append("canonical v17 scope is not 20")
    if [row.get("title") for row in v17_rows] != v17_major.get("items", []):
        errors.append("ledger v17 items differ from canonical v17 major scope")
    matrix_titles = [row.get("item") for row in v17_matrix.get("items", [])]
    if matrix_titles != v17_major.get("items", []):
        errors.append("v17 closure matrix differs from frozen v17 major scope")
    for row in v17_matrix.get("items", []):
        for ref in row.get("code_owners", []) + row.get("fresh_executable_evidence", []):
            if not (ROOT / ref).exists():
                errors.append(f"v17 evidence reference missing: {ref}")

    v18_major = load(ROOT / "v18_major_scope.json")
    delivery = load(ROOT / "v18_delivery_slices.json")
    workstreams = ledger.get("v18", {}).get("workstreams", [])
    if [row.get("title") for row in workstreams] != v18_major.get("workstreams", []):
        errors.append("ledger v18 workstreams differ from canonical v18 major scope")
    if len(workstreams) != len(delivery.get("slices", [])):
        errors.append("v18 major workstream/delivery-slice cardinality mismatch")
    for row in workstreams:
        if not row.get("deliverySlice"):
            errors.append(f"v18 workstream lacks delivery placement: {row.get('title')}")

    ledger_scopes = {
        scope.get("source"): scope for scope in ledger.get("v18", {}).get("releaseScopes", [])
    }
    for path in RELEASE_SCOPE_FILES:
        if path not in ledger_scopes:
            errors.append(f"release scope missing from ledger: {path}")
            continue
        canonical = scope_entries(path, load(ROOT / path))
        recorded = [
            {key: row.get(key) for key in ("id", "title", "section")}
            for row in ledger_scopes[path].get("entries", [])
        ]
        if recorded != canonical:
            errors.append(f"release-scope entry drift: {path}")

    remediation = load(ROOT / "functionality_utility_remediation.json")
    carry = ledger.get("functionalityUtilityCarryForward", {})
    carry_rows = carry.get("items", [])
    expected_carry = remediation.get("items", [])
    if remediation.get("targetRelease") != "v18.3.0":
        errors.append("functionality utility remediation target drift")
    if len(expected_carry) != 13 or len(carry_rows) != 13:
        errors.append("functionality utility carry-forward is not 13")
    expected_disposition = [
        (row.get("name"), row.get("action"), row.get("priority"))
        for row in expected_carry
    ]
    recorded_disposition = [
        (row.get("name"), row.get("action"), row.get("priority"))
        for row in carry_rows
    ]
    if recorded_disposition != expected_disposition:
        errors.append("functionality utility carry-forward differs from canonical remediation")
    v18_3_scope = load(ROOT / "v18_3_scope.json")
    v18_3_text = json.dumps(v18_3_scope, sort_keys=True)
    silently_sliced = [
        row.get("name") for row in expected_carry
        if str(row.get("name", "")).lower() in v18_3_text.lower()
    ]
    if silently_sliced:
        errors.append("audit assumption drift: v18.3 scope now contains remediation names; reclassify ledger")


    release_entries = [
        row
        for scope in ledger.get("v18", {}).get("releaseScopes", [])
        for row in scope.get("entries", [])
    ]
    tracked_rows = (
        inherited_rows
        + v17_rows
        + workstreams
        + release_entries
        + carry_rows
        + ledger.get("confirmedImplementationMisses", [])
        + ledger.get("escapedDefects", [])
        + ledger.get("conversationalCommitments", [])
        + ledger.get("roadmapPlacedNotCurrentMisses", [])
        + ledger.get("independentAudit20260817", {}).get("findings", [])
    )
    tracked_ids = [row.get("id") for row in tracked_rows]
    declared_total = ledger.get("baseline", {}).get("totalTrackedRows")
    if declared_total != len(tracked_rows):
        errors.append(
            f"tracked-row conservation mismatch: baseline={declared_total}, actual={len(tracked_rows)}"
        )
    if any(not requirement_id for requirement_id in tracked_ids):
        errors.append("tracked row lacks immutable ID")
    duplicate_ids = sorted(
        requirement_id
        for requirement_id in set(tracked_ids)
        if requirement_id and tracked_ids.count(requirement_id) > 1
    )
    if duplicate_ids:
        errors.append("duplicate tracked requirement IDs: " + ", ".join(duplicate_ids))

    present_blockers = {
        row.get("id") for row in ledger.get("confirmedImplementationMisses", [])
    } | {
        row.get("id") for row in ledger.get("escapedDefects", [])
    } | {
        row.get("id") for row in ledger.get("conversationalCommitments", [])
    }
    missing_blockers = REQUIRED_BLOCKERS - present_blockers
    if missing_blockers:
        errors.append("known blocker missing from ledger: " + ", ".join(sorted(missing_blockers)))

    allowed = set(ledger.get("allowedStatuses", []))
    status_rows = collect_status_rows(ledger)
    for path, status, _ in status_rows:
        if status not in allowed and path != "ledger":
            errors.append(f"unsupported status {status!r} at {path}")

    blockers = [(path, status, row) for path, status, row in status_rows if status in BLOCKING]
    id_rows = {
        row.get("id"): (path, status, row)
        for path, status, row in status_rows
        if row.get("id")
    }

    def validate_fresh_pass(path: str, row: dict) -> None:
        evidence = row.get("closureEvidence")
        if not isinstance(evidence, dict):
            errors.append(f"FRESH_PASS lacks closureEvidence at {path}")
            return
        for key in ("sourceCommit", "regression", "macOS", "windows"):
            if not evidence.get(key):
                errors.append(f"FRESH_PASS lacks {key} evidence at {path}")

    if slice_mode:
        if current_slice.get("planningState") != "FROZEN":
            errors.append("current slice G1 is not FROZEN")
        assigned = current_slice.get("assignedIds", [])
        if not assigned:
            errors.append("current slice has no assigned requirement IDs")
        if len(assigned) != len(set(assigned)):
            errors.append("current slice contains duplicate requirement IDs")
        for requirement_id in assigned:
            hit = id_rows.get(requirement_id)
            if not hit:
                errors.append(f"current-slice ID missing from ledger: {requirement_id}")
                continue
            row_path, row_status, row = hit
            if row_status not in FINAL:
                errors.append(f"current-slice item is not final: {requirement_id}={row_status}")
            if row_status == "FRESH_PASS":
                validate_fresh_pass(row_path, row)

        assigned_set = set(assigned)
        for row_path, row_status, row in status_rows:
            requirement_id = row.get("id")
            if row_path == "ledger" or not requirement_id or requirement_id in assigned_set:
                continue
            if row_status == "PLACED_NEXT_V18":
                for key in ("targetRelease", "owner", "reason", "userImpact"):
                    if not row.get(key):
                        errors.append(f"next-slice placement lacks {key}: {requirement_id}")
            elif row_status not in FINAL:
                errors.append(
                    f"unassigned applicable row lacks final/next-slice disposition: "
                    f"{requirement_id}={row_status}"
                )
            elif row_status == "FRESH_PASS":
                validate_fresh_pass(row_path, row)

    if major_mode:
        if ledger.get("status") != "CLOSED":
            errors.append("major-closure ledger is not CLOSED")
        if blockers:
            errors.append(f"{len(blockers)} blocking/revalidation/next-slice rows remain")
        for row_path, row_status, row in status_rows:
            if row_status == "FRESH_PASS":
                validate_fresh_pass(row_path, row)
            elif row_status not in FINAL and row_path != "ledger":
                errors.append(f"non-final major-closure disposition at {row_path}: {row_status}")

    if errors:
        print("v17/v18 implementation reconciliation gate: FAIL")
        for error in errors:
            print(f" - {error}")
        return 2

    mode = "slice-release" if slice_mode else "major-closure" if major_mode else "inventory"
    print(
        "v17/v18 implementation reconciliation gate: PASS"
        f" · mode={mode}"
        f" · inherited={len(inherited_rows)}"
        f" · v17={len(v17_rows)}"
        f" · v18-workstreams={len(workstreams)}"
        f" · v18-release-entries={sum(len(x.get('entries', [])) for x in ledger_scopes.values())}"
        f" · functionality-remediation={len(carry_rows)}"
        f" · current-slice={current_slice.get('release', 'unknown')}"
        f" · slice-planning={current_slice.get('planningState', 'unknown')}"
        f" · open-blocking/revalidation/next={len(blockers)}"
    )
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
