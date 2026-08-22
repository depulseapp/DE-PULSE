#!/usr/bin/env python3
"""Fail-closed active work-slice closure-ledger validation.

The normal Fast invocation validates that the active work slice has a complete,
canonical machine-readable gap ledger. It does not pretend open implementation
work is complete. Final closure/G16 invokes --require-closed, which fails until
every blocking gap is VERIFIED with evidence.
"""
from __future__ import annotations

import argparse
import json
from collections import Counter
from pathlib import Path
import sys
from typing import Any

ROOT = Path(__file__).resolve().parents[2]
STATE_PATH = ROOT / "governance" / "current-state.json"
ALLOWED_STATUSES = {"OPEN", "IMPLEMENTED_UNVERIFIED", "BLOCKED_EXTERNAL", "VERIFIED"}
CLOSED_WORK_SLICE_STATES = {"READY_FOR_CLOSURE", "COMPLETE", "COMPLETED", "CLOSED", "DELIVERED"}
REQUIRED_70_GAPS = {
    "FAST-640-QUALIFIED-PATH",
    "PLANNER-EVIDENCE-OWNER-ROUTING",
    "RETIRED-TEST-EQUIVALENCE",
    "SOURCE-HEALTH-DEBT",
    "ACTIVE-VERSIONED-TEST-MIGRATION",
    "PACKAGE-DECOMPOSITION",
    "PERMANENT-ROOT-ALLOWLIST",
    "ASSET-REGISTRY-OWNERSHIP",
    "RELEASE-IDENTITY-FANOUT",
    "SEMVER-RELEASE-CUTOVER",
    "MAIN-PROTECTION-RULESET",
    "ARTIFACT-ATTESTATION-SBOM",
    "CURRENT-STATE-OVERLAY-CONVERGENCE",
    "G16-ROOT-CI-EFFICIENCY",
    "FINAL-QUALIFIED",
}


def load(path: Path) -> dict[str, Any]:
    return json.loads(path.read_text(encoding="utf-8"))


def nonempty_strings(value: object) -> bool:
    return isinstance(value, list) and bool(value) and all(isinstance(item, str) and item.strip() for item in value)


def validate(require_closed: bool) -> tuple[list[str], Counter[str], str]:
    errors: list[str] = []
    if not STATE_PATH.is_file():
        return ["missing governance/current-state.json"], Counter(), ""

    state = load(STATE_PATH)
    active = state.get("activeWorkSlice", {})
    work_slice_id = str(active.get("workSliceId", "")).strip()
    if not work_slice_id:
        return ["current-state activeWorkSlice.workSliceId missing"], Counter(), ""

    work_dir = ROOT / "governance" / "work-slices" / work_slice_id
    work_path = work_dir / "work-slice.json"
    closure_path = work_dir / "closure.json"
    if not work_path.is_file():
        errors.append(f"missing canonical work-slice contract: {work_path.relative_to(ROOT)}")
        return errors, Counter(), work_slice_id
    if not closure_path.is_file():
        errors.append(f"missing executable closure ledger: {closure_path.relative_to(ROOT)}")
        return errors, Counter(), work_slice_id

    work = load(work_path)
    closure = load(closure_path)
    if closure.get("schema") != "DE.PULSE-WORK-SLICE-CLOSURE-1":
        errors.append("unsupported work-slice closure schema")
    for field, expected in (
        ("workSliceId", work_slice_id),
        ("issue", active.get("issue")),
        ("branch", active.get("branch")),
    ):
        if closure.get(field) != expected:
            errors.append(f"closure/current-state mismatch for {field}: {closure.get(field)!r} != {expected!r}")
    if work.get("workSliceId") != work_slice_id:
        errors.append("work-slice/current-state workSliceId mismatch")
    if work.get("issue") != active.get("issue"):
        errors.append("work-slice/current-state issue mismatch")
    if work.get("branch") != active.get("branch"):
        errors.append("work-slice/current-state branch mismatch")
    expected_ledger = f"governance/work-slices/{work_slice_id}/closure.json"
    if work.get("closureLedger") != expected_ledger:
        errors.append(f"work-slice must name canonical closure ledger {expected_ledger}")
    if closure.get("allGapsMustBeVerified") is not True:
        errors.append("closure ledger must require all gaps VERIFIED")
    if not str(closure.get("closurePolicy", "")).strip():
        errors.append("closure policy text missing")

    gaps = closure.get("gaps")
    if not isinstance(gaps, list) or not gaps:
        errors.append("closure ledger gaps must be a non-empty array")
        return errors, Counter(), work_slice_id

    seen: set[str] = set()
    statuses: Counter[str] = Counter()
    for index, gap in enumerate(gaps):
        if not isinstance(gap, dict):
            errors.append(f"gap[{index}] must be an object")
            continue
        gid = str(gap.get("id", "")).strip()
        if not gid:
            errors.append(f"gap[{index}] missing id")
            continue
        if gid in seen:
            errors.append(f"duplicate closure gap id: {gid}")
        seen.add(gid)
        status = str(gap.get("status", "")).strip()
        statuses[status] += 1
        if status not in ALLOWED_STATUSES:
            errors.append(f"{gid}: unsupported status {status!r}")
        if gap.get("blocksIssueClosure") is not True:
            errors.append(f"{gid}: blocksIssueClosure must be true for #70")
        if not str(gap.get("owner", "")).strip():
            errors.append(f"{gid}: owner missing")
        if not nonempty_strings(gap.get("implementationPaths")):
            errors.append(f"{gid}: implementationPaths must contain real repository/settings owners")
        if not nonempty_strings(gap.get("evidenceRequired")):
            errors.append(f"{gid}: evidenceRequired must be non-empty")
        if not str(gap.get("closureCondition", "")).strip():
            errors.append(f"{gid}: closureCondition missing")
        evidence = gap.get("evidence", [])
        if evidence is not None and not isinstance(evidence, list):
            errors.append(f"{gid}: evidence must be an array when present")
        if status == "VERIFIED" and not nonempty_strings(evidence):
            errors.append(f"{gid}: VERIFIED requires non-empty executable evidence references")

    if work_slice_id == "ADAPT-CI-CONVERGENCE-001":
        missing = sorted(REQUIRED_70_GAPS - seen)
        extra = sorted(seen - REQUIRED_70_GAPS)
        if missing:
            errors.append("#70 closure ledger missing mandatory gap ids: " + ", ".join(missing))
        if extra:
            # Additional discovered gaps are allowed only when explicitly registered;
            # keep them visible rather than failing solely for being stricter.
            print("additional registered #70 closure gaps: " + ", ".join(extra))

    unresolved = [
        str(gap.get("id"))
        for gap in gaps
        if isinstance(gap, dict) and gap.get("blocksIssueClosure") is True and gap.get("status") != "VERIFIED"
    ]
    work_status = str(work.get("status", active.get("status", ""))).strip().upper()
    final_required = require_closed or work_status in CLOSED_WORK_SLICE_STATES
    if final_required and unresolved:
        errors.append("work-slice closure blocked by unresolved gaps: " + ", ".join(sorted(unresolved)))

    return errors, statuses, work_slice_id


def main() -> int:
    parser = argparse.ArgumentParser(description="Validate DE.PULSE active work-slice executable closure ledger")
    parser.add_argument("--require-closed", action="store_true", help="Fail unless every blocking gap is VERIFIED")
    args = parser.parse_args()

    errors, statuses, work_slice_id = validate(args.require_closed)
    print(f"DE.PULSE active work-slice closure ledger: {work_slice_id or 'UNKNOWN'}")
    for status in sorted(ALLOWED_STATUSES):
        print(f"{status}: {statuses.get(status, 0)}")
    unresolved = sum(count for status, count in statuses.items() if status != "VERIFIED")
    print(f"unresolved blocking gaps: {unresolved}")
    print("documentation-only closure: PROHIBITED")
    if errors:
        print("DE.PULSE work-slice closure ledger: FAIL", file=sys.stderr)
        for error in errors:
            print(" - " + error, file=sys.stderr)
        return 1
    if args.require_closed:
        print("all blocking gaps executable/evidence verified: PASS")
    else:
        print("ledger completeness/enforcement contract: PASS")
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
