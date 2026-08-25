#!/usr/bin/env python3
"""Fail-closed structural gate for the v18 final feature-assurance ledger.

T1 uses this gate before inventory freeze. It intentionally does not prove T2-T9
behavior; it proves that the machine ledger is structurally complete enough for
those tracks to consume without silently dropping source-discovered rows.
"""

from __future__ import annotations

import json
from pathlib import Path
import sys

ROOT = Path(__file__).resolve().parents[2]
PROGRAM = ROOT / "governance" / "programs" / "ADAPT-V18-FINAL-CLOSURE-10-10-001"
LEDGER = PROGRAM / "feature-assurance-ledger.json"
SCAN = PROGRAM / "T1_INDEPENDENT_OMISSION_SCAN.json"

REQUIRED_ROW_KEYS = {
    "id",
    "name",
    "category",
    "requirementProvenance",
    "canonicalSourceOwners",
    "consumers",
    "existingRegressionOwners",
    "positiveFunctionalEvidenceExpectation",
    "assuranceProfile",
    "durableRegressionOwner",
    "currentAssuranceState",
    "blockingStates",
    "uiDisposition",
}


def load_json(path: Path) -> dict:
    try:
        value = json.loads(path.read_text(encoding="utf-8"))
    except Exception as exc:  # pragma: no cover - command-line diagnostics
        raise SystemExit(f"FAIL: cannot read {path.relative_to(ROOT)}: {exc}") from exc
    if not isinstance(value, dict):
        raise SystemExit(f"FAIL: {path.relative_to(ROOT)} must contain a JSON object")
    return value


def nonempty_list(row: dict, key: str) -> bool:
    value = row.get(key)
    return isinstance(value, list) and bool(value) and all(isinstance(x, str) and x.strip() for x in value)


def main() -> int:
    ledger = load_json(LEDGER)
    scan = load_json(SCAN)

    rules = ledger.get("rules") or {}
    allowed_states = set(rules.get("allowedFeatureStates") or [])
    allowed_blocking = set(rules.get("blockingStates") or [])
    allowed_dispositions = set(ledger.get("uiDispositionValues") or [])
    profiles = ledger.get("assuranceProfiles") or {}
    rows = ledger.get("features")

    errors: list[str] = []
    if not isinstance(rows, list) or not rows:
        errors.append("features must be a non-empty array")
        rows = []

    ids: set[str] = set()
    for index, row in enumerate(rows):
        label = f"features[{index}]"
        if not isinstance(row, dict):
            errors.append(f"{label} is not an object")
            continue
        missing = sorted(REQUIRED_ROW_KEYS - set(row))
        if missing:
            errors.append(f"{label} missing keys: {', '.join(missing)}")
        feature_id = str(row.get("id") or "").strip()
        if not feature_id:
            errors.append(f"{label} has empty id")
        elif feature_id in ids:
            errors.append(f"duplicate feature id: {feature_id}")
        else:
            ids.add(feature_id)
        for key in ("requirementProvenance", "canonicalSourceOwners", "consumers", "existingRegressionOwners", "durableRegressionOwner"):
            if not nonempty_list(row, key):
                errors.append(f"{feature_id or label} has empty/invalid {key}")
        expectation = str(row.get("positiveFunctionalEvidenceExpectation") or "").strip()
        if not expectation:
            errors.append(f"{feature_id or label} has no positive functional evidence expectation")
        profile = str(row.get("assuranceProfile") or "").strip()
        if profile not in profiles:
            errors.append(f"{feature_id or label} references unknown assuranceProfile {profile!r}")
        else:
            profile_keys = set(profiles[profile])
            expected_tracks = {f"T{i}" for i in range(2, 10)}
            missing_tracks = sorted(expected_tracks - profile_keys)
            if missing_tracks:
                errors.append(f"assuranceProfile {profile} missing downstream expectations: {', '.join(missing_tracks)}")
        state = str(row.get("currentAssuranceState") or "").strip()
        if state not in allowed_states:
            errors.append(f"{feature_id or label} has invalid assurance state {state!r}")
        disposition = str(row.get("uiDisposition") or "").strip()
        if disposition not in allowed_dispositions:
            errors.append(f"{feature_id or label} has invalid uiDisposition {disposition!r}")
        blocking = row.get("blockingStates")
        if not isinstance(blocking, list):
            errors.append(f"{feature_id or label} blockingStates must be an array")
        else:
            unknown = sorted({str(x) for x in blocking} - allowed_blocking)
            if unknown:
                errors.append(f"{feature_id or label} has unknown blocking states: {', '.join(unknown)}")

    omissions = scan.get("omissionsFound") or []
    omission_ids = {
        str(item.get("id") or "").strip()
        for item in omissions
        if isinstance(item, dict) and str(item.get("id") or "").strip()
    }
    missing_scan_rows = sorted(omission_ids - ids)
    if missing_scan_rows:
        errors.append(
            "independent omission scan rows not yet represented in ledger: "
            + ", ".join(missing_scan_rows)
        )

    implementation_misses = scan.get("implementationMisses") or []
    if implementation_misses:
        corrective_rows = {
            str(row.get("id") or "")
            for row in rows
            if isinstance(row, dict) and row.get("currentAssuranceState") == "CORRECTIVE_REQUIRED"
        }
        for miss in implementation_misses:
            miss_id = str((miss or {}).get("featureId") or (miss or {}).get("id") or "").strip() if isinstance(miss, dict) else ""
            if miss_id and miss_id not in corrective_rows:
                errors.append(f"implementation miss {miss_id} lacks CORRECTIVE_REQUIRED ledger row")

    discovery_complete = ledger.get("discoveryComplete") is True
    if discovery_complete:
        if ledger.get("unexplainedGapCount") != 0:
            errors.append("discoveryComplete requires unexplainedGapCount == 0")
        for row in rows:
            if isinstance(row, dict) and row.get("blockingStates"):
                errors.append(f"discoveryComplete but {row.get('id')} still has blockingStates")
        if scan.get("state") != "COMPLETE":
            errors.append("discoveryComplete requires independent omission scan state COMPLETE")

    if errors:
        print("V18 FEATURE ASSURANCE LEDGER GATE: FAIL")
        for error in errors:
            print(f"- {error}")
        return 1

    print(
        "V18 FEATURE ASSURANCE LEDGER GATE: PASS "
        f"({len(rows)} unique rows; {len(omission_ids)} independent-scan rows covered; discoveryComplete={discovery_complete})"
    )
    return 0


if __name__ == "__main__":
    sys.exit(main())
