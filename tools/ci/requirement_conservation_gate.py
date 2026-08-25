#!/usr/bin/env python3
"""Permanent fail-closed conservation gate for the reserved v19/#66 program."""
from __future__ import annotations

import json
from pathlib import Path
import re
import sys

ROOT = Path(__file__).resolve().parents[2]
PROGRAM = ROOT / "governance" / "programs" / "ADAPT-HOSTED-SYNC-001"
LEDGER = PROGRAM / "requirement-conservation.json"
PLAN = PROGRAM / "V19_ZERO_MISS_PLAN.md"
CURRENT = ROOT / "governance" / "current-state.json"
ROADMAP = ROOT / "governance" / "ROADMAP.md"

EXPECTED_IDS = [f"HOST-{n:03d}" for n in range(1, 73)]
REQUIRED_RULES = (
    "noUnassignedApplicableRequirement",
    "sourceOverlapBeforeG1",
    "onePrimaryResponsibilityPerVersion",
    "bandClosureRequiredBeforeNextBand",
)
REQUIRED_BOUNDARIES = {
    "US_EQUITIES_PROCESSING",
    "NO_EXECUTION",
    "SMART_PROVIDER_ROUTER_V2_SOLE_ROUTING_OWNER",
    "DIRECT_SEC_EDGAR_FORM4_AUTHORITY",
    "GLD_SLV_USO_ACTIONABLE_EXCEPTIONS",
    "NO_PARALLEL_CANONICAL_OWNERS",
}


def load(path: Path) -> dict:
    return json.loads(path.read_text(encoding="utf-8-sig"))


def main() -> int:
    errors: list[str] = []
    for path in (LEDGER, PLAN, CURRENT, ROADMAP):
        if not path.is_file():
            errors.append(f"required conservation owner missing: {path.relative_to(ROOT)}")
    if errors:
        for error in errors:
            print(" - " + error, file=sys.stderr)
        return 1

    ledger = load(LEDGER)
    current = load(CURRENT)
    rules = ledger.get("rules") or {}
    rows = ledger.get("requirements") or []
    boundaries = {str(x) for x in rules.get("permanentBoundaries") or []}

    if ledger.get("schema") != "DE.PULSE-V19-REQUIREMENT-CONSERVATION-1":
        errors.append("v19 conservation schema drifted")
    program = ledger.get("program") or {}
    planning = ledger.get("planning") or {}
    if program.get("issue") != 66 or program.get("scopeId") != "ADAPT-HOSTED-SYNC-001":
        errors.append("v19 conservation program identity must remain #66 / ADAPT-HOSTED-SYNC-001")
    if planning.get("issue") != 110 or planning.get("scopeId") != "ADAPT-V19-ZERO-MISS-PLAN-001":
        errors.append("v19 conservation planning identity must remain #110 / ADAPT-V19-ZERO-MISS-PLAN-001")
    if not isinstance(rows, list):
        errors.append("v19 conservation requirements must be an array")
        rows = []

    ids = [str(row.get("id") or "") for row in rows if isinstance(row, dict)]
    if len(rows) != 72 or ids != EXPECTED_IDS:
        errors.append("v19 conservation ledger must retain exactly ordered HOST-001..HOST-072")
    if len(set(ids)) != len(ids):
        errors.append("v19 conservation ledger contains duplicate requirement IDs")

    for key in REQUIRED_RULES:
        if rules.get(key) is not True:
            errors.append(f"v19 conservation rule must remain true: {key}")

    missing_boundaries = sorted(REQUIRED_BOUNDARIES - boundaries)
    if missing_boundaries:
        errors.append("v19 permanent boundaries missing: " + ", ".join(missing_boundaries))
    lockstep = str(rules.get("crossPlatformLockstep") or "")
    if not all(token in lockstep for token in ("Mac", "Windows", "Web", "REQUIRED")):
        errors.append("v19 crossPlatformLockstep must keep Mac + Windows + Web required-or-justified-N/A semantics")

    for index, row in enumerate(rows, start=1):
        if not isinstance(row, dict):
            errors.append(f"HOST-{index:03d} row is not an object")
            continue
        row_id = str(row.get("id") or "")
        for field in ("requirement", "origin", "planningDisposition", "canonicalOwner", "plannedVersion", "evidenceState"):
            value = row.get(field)
            if not isinstance(value, str) or not value.strip():
                errors.append(f"{row_id or index}: required field missing: {field}")
        disposition = str(row.get("planningDisposition") or "").upper()
        if any(token in disposition for token in ("UNASSIGNED", "UNKNOWN", "TBD")):
            errors.append(f"{row_id}: planningDisposition may not be unassigned/unknown")
        planned_version = str(row.get("plannedVersion") or "")
        if not re.fullmatch(r"v19\.\d+\.\d+", planned_version):
            errors.append(f"{row_id}: plannedVersion must remain an explicit v19 patch")
        dependencies = row.get("dependencies")
        if dependencies is not None and not isinstance(dependencies, list):
            errors.append(f"{row_id}: dependencies must be an array when present")

    gate = current.get("productCapabilityGate") or {}
    blocked_issue = gate.get("futureBlockedIssue")
    if blocked_issue == 66:
        if rules.get("productProgramReserved") is not False:
            errors.append("#66 is blocked in current-state, so productProgramReserved must remain false")
        if gate.get("futureProgramBlockerIssue") != 113:
            errors.append("v19/#66 must remain blocked by v18.10 parent #113 before publication")
    elif rules.get("productProgramReserved") not in (False, True):
        errors.append("productProgramReserved must be boolean")

    plan_text = PLAN.read_text(encoding="utf-8", errors="replace")
    roadmap_text = ROADMAP.read_text(encoding="utf-8", errors="replace")
    for marker in ("HOST-001", "HOST-072", "source-overlap", "band", "zero-gap"):
        if marker.lower() not in plan_text.lower():
            errors.append(f"v19 zero-miss plan lost conservation marker: {marker}")
    for marker in ("Zero-Miss Future-Version Conservation", "requirement-conservation.json", "v19"):
        if marker not in roadmap_text:
            errors.append(f"canonical roadmap lost v19 conservation marker: {marker}")

    print("V19 REQUIREMENT CONSERVATION")
    print(f"rows: {len(rows)} / expected 72")
    print(f"first/last: {ids[0] if ids else '-'} / {ids[-1] if ids else '-'}")
    print(f"program reserved: {rules.get('productProgramReserved')}")
    print("source-overlap before G1: REQUIRED")
    print("band zero-gap closure before next band: REQUIRED")
    if errors:
        print("V19 REQUIREMENT CONSERVATION GATE: FAIL", file=sys.stderr)
        for error in errors:
            print(" - " + error, file=sys.stderr)
        return 1
    print("V19 REQUIREMENT CONSERVATION GATE: PASS")
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
