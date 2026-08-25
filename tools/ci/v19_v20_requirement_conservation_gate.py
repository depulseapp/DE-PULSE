#!/usr/bin/env python3
"""Fail closed on v19/v20 requirement loss or unowned future work.

This is an existing G2/G10 governance check, not a new top-level gate. It makes
issue-to-version conservation executable without deciding product behavior.
"""
from __future__ import annotations

import json
import re
from pathlib import Path

ROOT = Path(__file__).resolve().parents[2]
LEDGER = ROOT / "governance" / "v19-v20-requirement-conservation.json"
ROADMAP = ROOT / "governance" / "ROADMAP.md"
PLAN = ROOT / "governance" / "V19_V20_ZERO_MISS_PLAN.md"

REQUIRED_SOURCE_REFS = {
    "issue:57",
    "issue:65",
    "issue:66:body",
    "issue:66:comment:5376362643",
    "issue:66:comment:5376382961",
    "issue:66:comment:5376545349",
    "issue:66:comment:5376571835",
    "issue:66:comment:5376601126",
    "issue:66:comment:5376640936",
    "issue:66:comment:5376734976",
    "issue:110",
}
ALLOWED_DISPOSITIONS = {"INHERITED", "IMPLEMENT_IN", "FUTURE_BLOCKED"}
FORBIDDEN_STATES = {
    "UNASSIGNED",
    "OPEN_WITHOUT_TARGET",
    "DUPLICATE_PRIMARY_OWNER",
    "UNEXPLAINED_CARRY_FORWARD",
    "EVIDENCE_MISSING_AT_BAND_CLOSURE",
}
VERSION_RE = re.compile(r"^v(?:19|20)\.\d+\.\d+$")


def load_json(path: Path) -> dict:
    return json.loads(path.read_text(encoding="utf-8"))


def main() -> int:
    errors: list[str] = []
    for path in (LEDGER, ROADMAP, PLAN):
        if not path.is_file():
            errors.append(f"missing required zero-miss authority: {path.relative_to(ROOT)}")
    if errors:
        print("\n".join(errors))
        return 2

    ledger = load_json(LEDGER)
    roadmap = ROADMAP.read_text(encoding="utf-8")
    plan = PLAN.read_text(encoding="utf-8")

    if ledger.get("schema") != "DE.PULSE-V19-V20-REQUIREMENT-CONSERVATION-1":
        errors.append("unsupported v19/v20 conservation ledger schema")
    if ledger.get("planningIssue") != 110 or ledger.get("hostedProgramIssue") != 66:
        errors.append("planning/hosted issue identity mismatch")

    coverage = ledger.get("sourceCoverage", [])
    covered = {
        str(item.get("sourceRef", "")).strip()
        for item in coverage
        if isinstance(item, dict) and str(item.get("status", "")).strip() == "MAPPED"
    }
    missing_sources = sorted(REQUIRED_SOURCE_REFS - covered)
    if missing_sources:
        errors.append("unmapped GitHub requirement sources: " + ", ".join(missing_sources))

    v18 = ledger.get("v18Closure", {})
    if not isinstance(v18, dict) or v18.get("state") != "CLOSED_BY_EXECUTABLE_EVIDENCE":
        errors.append("v18 closure state must remain CLOSED_BY_EXECUTABLE_EVIDENCE")
    if isinstance(v18, dict) and v18.get("openCorrectiveCount") != 0:
        errors.append("v18 openCorrectiveCount is non-zero; v19 must remain blocked until reconciled")

    closure_versions = ledger.get("bandClosures", [])
    if not isinstance(closure_versions, list) or not closure_versions:
        errors.append("bandClosures must be a non-empty array")
    else:
        for version in closure_versions:
            if not isinstance(version, str) or not VERSION_RE.match(version):
                errors.append(f"invalid band closure version: {version!r}")
                continue
            if version not in roadmap or version not in plan:
                errors.append(f"band closure {version} missing from roadmap or zero-miss plan")

    requirements = ledger.get("requirements", [])
    if not isinstance(requirements, list) or not requirements:
        errors.append("requirements must be a non-empty array")
        requirements = []

    seen_ids: set[str] = set()
    for index, row in enumerate(requirements):
        if not isinstance(row, dict):
            errors.append(f"requirement[{index}] must be an object")
            continue
        rid = str(row.get("id", "")).strip()
        if not rid:
            errors.append(f"requirement[{index}] missing id")
            continue
        if rid in seen_ids:
            errors.append(f"duplicate requirement id: {rid}")
        seen_ids.add(rid)

        source_refs = row.get("sourceRefs")
        if not isinstance(source_refs, list) or not source_refs or not all(isinstance(x, str) and x.strip() for x in source_refs):
            errors.append(f"{rid}: sourceRefs must be non-empty strings")
        elif any(ref not in covered for ref in source_refs):
            errors.append(f"{rid}: references a source not declared MAPPED in sourceCoverage")

        if not str(row.get("requirement", "")).strip():
            errors.append(f"{rid}: requirement text missing")
        if not str(row.get("canonicalOwner", "")).strip():
            errors.append(f"{rid}: canonicalOwner missing")

        disposition = str(row.get("disposition", "")).strip()
        state = str(row.get("state", "")).strip()
        target = str(row.get("targetVersion", "")).strip()
        if disposition not in ALLOWED_DISPOSITIONS:
            errors.append(f"{rid}: unsupported disposition {disposition!r}")
        if state in FORBIDDEN_STATES:
            errors.append(f"{rid}: blocking conservation state {state}")
        if not target:
            errors.append(f"{rid}: targetVersion missing")

        if disposition == "INHERITED" and not str(row.get("evidence", "")).strip():
            errors.append(f"{rid}: INHERITED requires evidence")
        if disposition == "FUTURE_BLOCKED" and not str(row.get("blocker", "")).strip():
            errors.append(f"{rid}: FUTURE_BLOCKED requires blocker")
        if disposition == "IMPLEMENT_IN" and VERSION_RE.match(target):
            if target not in roadmap:
                errors.append(f"{rid}: target {target} missing from governance/ROADMAP.md")
            if target not in plan:
                errors.append(f"{rid}: target {target} missing from governance/V19_V20_ZERO_MISS_PLAN.md")

    if len(seen_ids) < 50:
        errors.append(f"requirement ledger unexpectedly small ({len(seen_ids)} rows); investigate requirement loss")

    roadmap_versions = set(re.findall(r"v(?:19|20)\.\d+\.\d+", roadmap))
    plan_versions = set(re.findall(r"v(?:19|20)\.\d+\.\d+", plan))
    if roadmap_versions != plan_versions:
        missing_from_roadmap = sorted(plan_versions - roadmap_versions)
        missing_from_plan = sorted(roadmap_versions - plan_versions)
        if missing_from_roadmap:
            errors.append("versions missing from roadmap: " + ", ".join(missing_from_roadmap))
        if missing_from_plan:
            errors.append("versions missing from zero-miss plan: " + ", ".join(missing_from_plan))

    if errors:
        print("DE.PULSE v19/v20 requirement conservation: FAIL")
        for error in errors:
            print(f"- {error}")
        return 2

    print("DE.PULSE v19/v20 requirement conservation: PASS")
    print(f"- mapped GitHub sources: {len(covered)}")
    print(f"- conserved requirements: {len(seen_ids)}")
    print(f"- planned v19/v20 versions: {len(roadmap_versions)}")
    print(f"- zero-gap/major closure checkpoints: {len(closure_versions)}")
    print("- v18 closure: CLOSED_BY_EXECUTABLE_EVIDENCE / 0 open corrective")
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
