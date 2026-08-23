#!/usr/bin/env python3
"""Fail-closed active work-slice closure-ledger validation.

Normal Fast validates a complete canonical machine-readable gap ledger. Final
closure/G16 invokes --require-closed. Blocking gaps normally require VERIFIED
status. A narrowly scoped external platform control may satisfy closure only
through an explicit machine-readable owner-approved waiver validated here; the
factual blocked state is never relabeled as technically enforced.
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
WAIVABLE_EXTERNAL_GAPS = {("ADAPT-CI-CONVERGENCE-001", "MAIN-PROTECTION-RULESET")}
REQUIRED_MAIN_PROTECTION_COMPENSATING_CONTROLS = {
    "PR_FIRST_DEVELOPMENT",
    "EXACT_HEAD_FAST_STATUS",
    "EXACT_HEAD_QUALIFIED_STATUS",
    "NO_DIRECT_MAIN_PUSH_POLICY",
    "NO_FORCE_PUSH_POLICY",
    "CANONICAL_RELEASE_G11_G16",
    "EXACT_SHA_FINGERPRINT_PROVENANCE",
}
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


def validate_external_waiver(
    work_slice_id: str,
    issue: object,
    work: dict[str, Any],
    gap: dict[str, Any],
    errors: list[str],
) -> bool:
    """Validate the one bounded external-control waiver allowed by #70."""
    gid = str(gap.get("id", "")).strip()
    if (work_slice_id, gid) not in WAIVABLE_EXTERNAL_GAPS:
        errors.append(f"{gid}: external waiver is not permitted for this work slice/gap")
        return False

    mapping = work.get("externalControlWaivers", {})
    if not isinstance(mapping, dict):
        errors.append("work-slice externalControlWaivers must be an object")
        return False
    waiver_rel = str(mapping.get(gid, "")).strip()
    if not waiver_rel:
        errors.append(f"{gid}: BLOCKED_EXTERNAL requires a registered waiver path")
        return False
    waiver_path = ROOT / waiver_rel
    if not waiver_path.is_file():
        errors.append(f"{gid}: registered waiver file missing: {waiver_rel}")
        return False

    start_errors = len(errors)
    waiver = load(waiver_path)
    if waiver.get("schema") != "DE.PULSE-EXTERNAL-CONTROL-WAIVER-1":
        errors.append(f"{gid}: unsupported external waiver schema")
    if waiver.get("status") != "APPROVED":
        errors.append(f"{gid}: external waiver status must be APPROVED")
    if waiver.get("workSliceId") != work_slice_id:
        errors.append(f"{gid}: waiver workSliceId mismatch")
    if waiver.get("issue") != issue:
        errors.append(f"{gid}: waiver issue mismatch")
    if waiver.get("gapId") != gid:
        errors.append(f"{gid}: waiver gapId mismatch")
    if waiver.get("scope") != "GITHUB_MAIN_PROTECTION_ONLY":
        errors.append(f"{gid}: waiver scope must remain GITHUB_MAIN_PROTECTION_ONLY")
    if waiver.get("noProductBehaviorChange") is not True:
        errors.append(f"{gid}: waiver must assert noProductBehaviorChange=true")
    if waiver.get("noReleaseEvidenceInvalidation") is not True:
        errors.append(f"{gid}: waiver must preserve existing release evidence")

    actual = waiver.get("actualState", {})
    if not isinstance(actual, dict):
        errors.append(f"{gid}: waiver actualState missing")
    else:
        if actual.get("repository") != "depulseapp/DE-PULSE":
            errors.append(f"{gid}: waiver repository mismatch")
        if actual.get("mainProtected") is not False:
            errors.append(f"{gid}: waiver must truthfully retain mainProtected=false")
        if actual.get("rulesetConfigured") is not True:
            errors.append(f"{gid}: configured ruleset evidence missing")
        if actual.get("rulesetEnforced") is not False:
            errors.append(f"{gid}: waiver must truthfully retain rulesetEnforced=false")
        if actual.get("enforcementAvailability") != "UNAVAILABLE_CURRENT_PLAN":
            errors.append(f"{gid}: unexpected enforcement availability classification")

    limitation = waiver.get("limitation", {})
    if not isinstance(limitation, dict):
        errors.append(f"{gid}: waiver limitation missing")
    else:
        if limitation.get("provider") != "GitHub":
            errors.append(f"{gid}: waiver provider must be GitHub")
        if limitation.get("category") != "PLAN_ENFORCEMENT_LIMITATION":
            errors.append(f"{gid}: waiver category mismatch")
        if limitation.get("enforcementAvailable") is not False:
            errors.append(f"{gid}: waiver may only apply while technical enforcement is unavailable")
        if not str(limitation.get("detail", "")).strip():
            errors.append(f"{gid}: waiver limitation detail missing")

    decision = waiver.get("ownerDecision", {})
    if not isinstance(decision, dict):
        errors.append(f"{gid}: ownerDecision missing")
    else:
        if decision.get("approved") is not True:
            errors.append(f"{gid}: owner approval missing")
        if decision.get("upgradeDecision") != "DECLINED":
            errors.append(f"{gid}: owner upgrade decision must be DECLINED for this waiver")
        if decision.get("scopeLimited") is not True:
            errors.append(f"{gid}: waiver must be explicitly scope-limited")
        if not str(decision.get("decisionRecordedAt", "")).strip():
            errors.append(f"{gid}: owner decision date missing")

    risk = waiver.get("riskAcceptance", {})
    if not isinstance(risk, dict) or risk.get("accepted") is not True:
        errors.append(f"{gid}: residual risk must be explicitly accepted")
    elif not nonempty_strings(risk.get("residualRisks")):
        errors.append(f"{gid}: residual risks must be enumerated")

    controls = waiver.get("compensatingControls", [])
    control_ids: set[str] = set()
    if not isinstance(controls, list) or not controls:
        errors.append(f"{gid}: compensatingControls must be non-empty")
    else:
        for item in controls:
            if not isinstance(item, dict):
                errors.append(f"{gid}: compensating control must be an object")
                continue
            cid = str(item.get("id", "")).strip()
            if not cid:
                errors.append(f"{gid}: compensating control id missing")
                continue
            control_ids.add(cid)
            if item.get("mandatory") is not True:
                errors.append(f"{gid}: compensating control {cid} must remain mandatory")
            if not str(item.get("detail", "")).strip():
                errors.append(f"{gid}: compensating control {cid} detail missing")
    missing_controls = sorted(REQUIRED_MAIN_PROTECTION_COMPENSATING_CONTROLS - control_ids)
    if missing_controls:
        errors.append(f"{gid}: missing required compensating controls: {', '.join(missing_controls)}")

    if not nonempty_strings(waiver.get("revalidationTriggers")):
        errors.append(f"{gid}: revalidationTriggers must be non-empty")
    if not str(waiver.get("retirementCondition", "")).strip():
        errors.append(f"{gid}: retirementCondition missing")

    return len(errors) == start_errors


def validate(require_closed: bool) -> tuple[list[str], Counter[str], str, set[str]]:
    errors: list[str] = []
    if not STATE_PATH.is_file():
        return ["missing governance/current-state.json"], Counter(), "", set()

    state = load(STATE_PATH)
    active = state.get("activeWorkSlice", {})
    work_slice_id = str(active.get("workSliceId", "")).strip()
    if not work_slice_id:
        return ["current-state activeWorkSlice.workSliceId missing"], Counter(), "", set()

    work_dir = ROOT / "governance" / "work-slices" / work_slice_id
    work_path = work_dir / "work-slice.json"
    closure_path = work_dir / "closure.json"
    if not work_path.is_file():
        errors.append(f"missing canonical work-slice contract: {work_path.relative_to(ROOT)}")
        return errors, Counter(), work_slice_id, set()
    if not closure_path.is_file():
        errors.append(f"missing executable closure ledger: {closure_path.relative_to(ROOT)}")
        return errors, Counter(), work_slice_id, set()

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
        errors.append("closure ledger must retain allGapsMustBeVerified=true; only validated external-control waivers may satisfy a factual BLOCKED_EXTERNAL gap")
    if not str(closure.get("closurePolicy", "")).strip():
        errors.append("closure policy text missing")

    gaps = closure.get("gaps")
    if not isinstance(gaps, list) or not gaps:
        errors.append("closure ledger gaps must be a non-empty array")
        return errors, Counter(), work_slice_id, set()

    seen: set[str] = set()
    statuses: Counter[str] = Counter()
    waived_external: set[str] = set()
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
        if status == "BLOCKED_EXTERNAL":
            if validate_external_waiver(work_slice_id, active.get("issue"), work, gap, errors):
                waived_external.add(gid)

    if work_slice_id == "ADAPT-CI-CONVERGENCE-001":
        missing = sorted(REQUIRED_70_GAPS - seen)
        extra = sorted(seen - REQUIRED_70_GAPS)
        if missing:
            errors.append("#70 closure ledger missing mandatory gap ids: " + ", ".join(missing))
        if extra:
            print("additional registered #70 closure gaps: " + ", ".join(extra))

    unresolved = [
        str(gap.get("id"))
        for gap in gaps
        if isinstance(gap, dict)
        and gap.get("blocksIssueClosure") is True
        and gap.get("status") != "VERIFIED"
        and str(gap.get("id")) not in waived_external
    ]
    work_status = str(work.get("status", active.get("status", ""))).strip().upper()
    final_required = require_closed or work_status in CLOSED_WORK_SLICE_STATES
    if final_required and unresolved:
        errors.append("work-slice closure blocked by unresolved gaps: " + ", ".join(sorted(unresolved)))

    return errors, statuses, work_slice_id, waived_external


def main() -> int:
    parser = argparse.ArgumentParser(description="Validate DE.PULSE active work-slice executable closure ledger")
    parser.add_argument(
        "--require-closed",
        action="store_true",
        help="Fail unless every blocking gap is VERIFIED or is the narrowly validated external-control waiver",
    )
    args = parser.parse_args()

    errors, statuses, work_slice_id, waived_external = validate(args.require_closed)
    print(f"DE.PULSE active work-slice closure ledger: {work_slice_id or 'UNKNOWN'}")
    for status in sorted(ALLOWED_STATUSES):
        print(f"{status}: {statuses.get(status, 0)}")
    print(f"validated external-control waivers: {len(waived_external)}")
    for gid in sorted(waived_external):
        print(f"waiver-satisfied factual BLOCKED_EXTERNAL: {gid}")
    unresolved = sum(count for status, count in statuses.items() if status != "VERIFIED") - len(waived_external)
    print(f"unresolved blocking gaps: {max(unresolved, 0)}")
    print("documentation-only closure: PROHIBITED")
    print("external waiver does not claim technical enforcement: PASS")
    if errors:
        print("DE.PULSE work-slice closure ledger: FAIL", file=sys.stderr)
        for error in errors:
            print(" - " + error, file=sys.stderr)
        return 1
    if args.require_closed:
        print("all blocking gaps verified or explicitly external-control-waived: PASS")
    else:
        print("ledger completeness/enforcement contract: PASS")
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
