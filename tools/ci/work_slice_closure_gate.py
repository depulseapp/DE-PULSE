#!/usr/bin/env python3
"""Fail-closed active work-slice closure-ledger validation.

Blocking gaps normally require VERIFIED status. Two narrowly governed special
bindings exist for #70: the factual GitHub-plan external-control waiver and the
post-run immutable final-qualification evidence binding. Neither mechanism may
be generalized to hide ordinary implementation gaps.
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
FINAL_EVIDENCE_BINDABLE_GAPS = {("ADAPT-CI-CONVERGENCE-001", "FINAL-QUALIFIED")}
REQUIRED_MAIN_PROTECTION_COMPENSATING_CONTROLS = {
    "PR_FIRST_DEVELOPMENT",
    "EXACT_HEAD_FAST_STATUS",
    "EXACT_HEAD_QUALIFIED_STATUS",
    "NO_DIRECT_MAIN_PUSH_POLICY",
    "NO_FORCE_PUSH_POLICY",
    "CANONICAL_RELEASE_G11_G16",
    "EXACT_SHA_FINGERPRINT_PROVENANCE",
}
REQUIRED_FINAL_EVIDENCE_OWNERS = {
    "CI_HARNESS",
    "BACKEND_FULL_GO",
    "RACE_DETECTOR",
    "RANDOMIZED_PACKAGE_ORDER",
    "PERSISTENCE_DB",
    "SECURITY_DATA_RIGHTS",
    "RENDERER",
    "CHROME",
    "WEBKIT",
    "NATIVE_MACOS_PACKAGED_LIFECYCLE",
    "NATIVE_WINDOWS_PACKAGED_RUNTIME",
}
REQUIRED_70_GAPS = {
    "FAST-640-QUALIFIED-PATH", "PLANNER-EVIDENCE-OWNER-ROUTING",
    "RETIRED-TEST-EQUIVALENCE", "SOURCE-HEALTH-DEBT",
    "ACTIVE-VERSIONED-TEST-MIGRATION", "PACKAGE-DECOMPOSITION",
    "PERMANENT-ROOT-ALLOWLIST", "ASSET-REGISTRY-OWNERSHIP",
    "RELEASE-IDENTITY-FANOUT", "SEMVER-RELEASE-CUTOVER",
    "MAIN-PROTECTION-RULESET", "ARTIFACT-ATTESTATION-SBOM",
    "CURRENT-STATE-OVERLAY-CONVERGENCE", "G16-ROOT-CI-EFFICIENCY",
    "FINAL-QUALIFIED",
}


def load(path: Path) -> dict[str, Any]:
    return json.loads(path.read_text(encoding="utf-8"))


def nonempty_strings(value: object) -> bool:
    return isinstance(value, list) and bool(value) and all(isinstance(item, str) and item.strip() for item in value)


def projected_active_work(state: dict[str, Any]) -> tuple[dict[str, Any], str]:
    """Use an in-progress product reservation as the current closure authority.

    `activeWorkSlice` may intentionally retain a completed process/repository
    slice such as #73. Once that process slice unblocks a reserved product
    capability, the active product reservation is what current Fast/Qualified
    closure enforcement must inspect.
    """
    product = state.get("productCapabilityGate", {})
    if isinstance(product, dict) and str(product.get("reservationStatus", "")).strip().upper() == "IN_PROGRESS":
        return {
            "workSliceId": product.get("reservedWorkSliceId"),
            "issue": product.get("reservedIssue"),
            "branch": product.get("reservedBranch"),
            "closureLedger": product.get("closureLedger"),
            "status": product.get("reservationStatus"),
        }, "PRODUCT_CAPABILITY"
    active = state.get("activeWorkSlice", {})
    return (active if isinstance(active, dict) else {}), "WORK_SLICE"


def validate_external_waiver(work_slice_id: str, issue: object, work: dict[str, Any], gap: dict[str, Any], errors: list[str]) -> bool:
    gid = str(gap.get("id", "")).strip()
    if (work_slice_id, gid) not in WAIVABLE_EXTERNAL_GAPS:
        errors.append(f"{gid}: external waiver is not permitted for this work slice/gap")
        return False
    mapping = work.get("externalControlWaivers", {})
    waiver_rel = str(mapping.get(gid, "")).strip() if isinstance(mapping, dict) else ""
    if not waiver_rel:
        errors.append(f"{gid}: BLOCKED_EXTERNAL requires a registered waiver path")
        return False
    waiver_path = ROOT / waiver_rel
    if not waiver_path.is_file():
        errors.append(f"{gid}: registered waiver file missing: {waiver_rel}")
        return False
    start = len(errors)
    waiver = load(waiver_path)
    checks = (
        (waiver.get("schema") == "DE.PULSE-EXTERNAL-CONTROL-WAIVER-1", "unsupported external waiver schema"),
        (waiver.get("status") == "APPROVED", "external waiver status must be APPROVED"),
        (waiver.get("workSliceId") == work_slice_id, "waiver workSliceId mismatch"),
        (waiver.get("issue") == issue, "waiver issue mismatch"),
        (waiver.get("gapId") == gid, "waiver gapId mismatch"),
        (waiver.get("scope") == "GITHUB_MAIN_PROTECTION_ONLY", "waiver scope mismatch"),
        (waiver.get("noProductBehaviorChange") is True, "waiver must assert noProductBehaviorChange=true"),
        (waiver.get("noReleaseEvidenceInvalidation") is True, "waiver must preserve release evidence"),
    )
    for ok, message in checks:
        if not ok:
            errors.append(f"{gid}: {message}")
    actual = waiver.get("actualState", {})
    if not isinstance(actual, dict) or actual.get("repository") != "depulseapp/DE-PULSE" or actual.get("mainProtected") is not False or actual.get("rulesetConfigured") is not True or actual.get("rulesetEnforced") is not False or actual.get("enforcementAvailability") != "UNAVAILABLE_CURRENT_PLAN":
        errors.append(f"{gid}: factual unenforced main/ruleset state changed or is incomplete")
    limitation = waiver.get("limitation", {})
    if not isinstance(limitation, dict) or limitation.get("provider") != "GitHub" or limitation.get("category") != "PLAN_ENFORCEMENT_LIMITATION" or limitation.get("enforcementAvailable") is not False or not str(limitation.get("detail", "")).strip():
        errors.append(f"{gid}: GitHub plan limitation contract invalid")
    decision = waiver.get("ownerDecision", {})
    if not isinstance(decision, dict) or decision.get("approved") is not True or decision.get("upgradeDecision") != "DECLINED" or decision.get("scopeLimited") is not True or not str(decision.get("decisionRecordedAt", "")).strip():
        errors.append(f"{gid}: explicit owner decision contract invalid")
    risk = waiver.get("riskAcceptance", {})
    if not isinstance(risk, dict) or risk.get("accepted") is not True or not nonempty_strings(risk.get("residualRisks")):
        errors.append(f"{gid}: residual risk acceptance invalid")
    controls = waiver.get("compensatingControls", [])
    ids: set[str] = set()
    if isinstance(controls, list):
        for item in controls:
            if isinstance(item, dict) and str(item.get("id", "")).strip():
                cid = str(item["id"]).strip()
                ids.add(cid)
                if item.get("mandatory") is not True or not str(item.get("detail", "")).strip():
                    errors.append(f"{gid}: compensating control {cid} must remain mandatory and documented")
    missing = sorted(REQUIRED_MAIN_PROTECTION_COMPENSATING_CONTROLS - ids)
    if missing:
        errors.append(f"{gid}: missing required compensating controls: {', '.join(missing)}")
    if not nonempty_strings(waiver.get("revalidationTriggers")) or not str(waiver.get("retirementCondition", "")).strip():
        errors.append(f"{gid}: waiver revalidation/retirement contract missing")
    return len(errors) == start


def validate_final_qualification_binding(work_slice_id: str, issue: object, work: dict[str, Any], gap: dict[str, Any], errors: list[str]) -> bool:
    gid = str(gap.get("id", "")).strip()
    if (work_slice_id, gid) not in FINAL_EVIDENCE_BINDABLE_GAPS:
        return False
    rel = str(work.get("finalQualificationEvidence", "")).strip()
    if not rel:
        errors.append(f"{gid}: finalQualificationEvidence path missing")
        return False
    path = ROOT / rel
    if not path.is_file():
        errors.append(f"{gid}: final qualification evidence file missing: {rel}")
        return False
    start = len(errors)
    evidence = load(path)
    if evidence.get("schema") != "DE.PULSE-WORK-SLICE-FINAL-QUALIFICATION-1" or evidence.get("status") != "VERIFIED":
        errors.append(f"{gid}: final qualification evidence schema/status invalid")
    if evidence.get("workSliceId") != work_slice_id or evidence.get("issue") != issue or evidence.get("gapId") != gid:
        errors.append(f"{gid}: final qualification identity mismatch")
    candidate = str(evidence.get("candidateSha", "")).strip()
    merge_sha = str(evidence.get("mergeCommitSha", "")).strip()
    if len(candidate) != 40 or len(merge_sha) != 40:
        errors.append(f"{gid}: candidate/merge SHA must be full immutable SHAs")
    fast = evidence.get("fast", {})
    qualified = evidence.get("qualified", {})
    if not isinstance(fast, dict) or fast.get("context") != "DE.PULSE/fast-head" or fast.get("conclusion") != "success" or fast.get("candidateSha") != candidate or not isinstance(fast.get("runId"), int):
        errors.append(f"{gid}: exact-head Fast binding invalid")
    if not isinstance(qualified, dict) or qualified.get("context") != "DE.PULSE/qualified-head" or qualified.get("conclusion") != "success" or qualified.get("candidateSha") != candidate or not isinstance(qualified.get("runId"), int):
        errors.append(f"{gid}: exact-head Qualified binding invalid")
    owners = set(qualified.get("evidenceOwners", [])) if isinstance(qualified, dict) and isinstance(qualified.get("evidenceOwners"), list) else set()
    missing = sorted(REQUIRED_FINAL_EVIDENCE_OWNERS - owners)
    if missing:
        errors.append(f"{gid}: final Qualified evidence owners missing: {', '.join(missing)}")
    merge = evidence.get("merge", {})
    if not isinstance(merge, dict) or merge.get("merged") is not True or merge.get("expectedHeadSha") != candidate or merge.get("mergeCommitSha") != merge_sha or merge.get("pullRequest") != 71 or merge.get("base") != "main":
        errors.append(f"{gid}: expected-head protected merge binding invalid")
    if work.get("mergedPullRequest") != 71 or work.get("mergedCommitSha") != merge_sha:
        errors.append(f"{gid}: work-slice merge binding mismatch")
    if evidence.get("noProductBehaviorChange") is not True or evidence.get("stableReleaseEvidenceUnchanged") is not True:
        errors.append(f"{gid}: final binding must preserve product/release boundaries")
    return len(errors) == start


def validate(require_closed: bool) -> tuple[list[str], Counter[str], str, set[str], set[str]]:
    errors: list[str] = []
    if not STATE_PATH.is_file():
        return ["missing governance/current-state.json"], Counter(), "", set(), set()
    state = load(STATE_PATH)
    active, _authority = projected_active_work(state)
    work_slice_id = str(active.get("workSliceId", "")).strip()
    if not work_slice_id:
        return ["canonical current active work workSliceId missing"], Counter(), "", set(), set()
    work_dir = ROOT / "governance" / "work-slices" / work_slice_id
    work_path, closure_path = work_dir / "work-slice.json", work_dir / "closure.json"
    if not work_path.is_file() or not closure_path.is_file():
        return ["missing canonical work-slice contract or closure ledger"], Counter(), work_slice_id, set(), set()
    work, closure = load(work_path), load(closure_path)
    if closure.get("schema") != "DE.PULSE-WORK-SLICE-CLOSURE-1":
        errors.append("unsupported work-slice closure schema")
    for field, expected in (("workSliceId", work_slice_id), ("issue", active.get("issue")), ("branch", active.get("branch"))):
        if closure.get(field) != expected:
            errors.append(f"closure/current-state mismatch for {field}: {closure.get(field)!r} != {expected!r}")
    if work.get("workSliceId") != work_slice_id or work.get("issue") != active.get("issue") or work.get("branch") != active.get("branch"):
        errors.append("work-slice/current-state identity mismatch")
    expected_ledger = f"governance/work-slices/{work_slice_id}/closure.json"
    if work.get("closureLedger") != expected_ledger:
        errors.append(f"work-slice must name canonical closure ledger {expected_ledger}")
    if active.get("closureLedger") and active.get("closureLedger") != expected_ledger:
        errors.append(f"current-state must name canonical closure ledger {expected_ledger}")
    if closure.get("allGapsMustBeVerified") is not True:
        errors.append("closure ledger must retain allGapsMustBeVerified=true")
    if not str(closure.get("closurePolicy", "")).strip():
        errors.append("closure policy text missing")
    gaps = closure.get("gaps")
    if not isinstance(gaps, list) or not gaps:
        return errors + ["closure ledger gaps must be a non-empty array"], Counter(), work_slice_id, set(), set()
    seen: set[str] = set()
    statuses: Counter[str] = Counter()
    waived: set[str] = set()
    bound: set[str] = set()
    for index, gap in enumerate(gaps):
        if not isinstance(gap, dict):
            errors.append(f"gap[{index}] must be an object")
            continue
        gid = str(gap.get("id", "")).strip()
        status = str(gap.get("status", "")).strip()
        if not gid:
            errors.append(f"gap[{index}] missing id")
            continue
        if gid in seen:
            errors.append(f"duplicate closure gap id: {gid}")
        seen.add(gid)
        statuses[status] += 1
        if status not in ALLOWED_STATUSES:
            errors.append(f"{gid}: unsupported status {status!r}")
        if gap.get("blocksIssueClosure") is not True:
            errors.append(f"{gid}: blocksIssueClosure must be true for the active work slice")
        if not str(gap.get("owner", "")).strip() or not nonempty_strings(gap.get("implementationPaths")) or not nonempty_strings(gap.get("evidenceRequired")) or not str(gap.get("closureCondition", "")).strip():
            errors.append(f"{gid}: incomplete owner/path/evidence/closure contract")
        evidence = gap.get("evidence", [])
        if evidence is not None and not isinstance(evidence, list):
            errors.append(f"{gid}: evidence must be an array")
        if status == "VERIFIED" and not nonempty_strings(evidence):
            errors.append(f"{gid}: VERIFIED requires evidence")
        if status == "BLOCKED_EXTERNAL" and validate_external_waiver(work_slice_id, active.get("issue"), work, gap, errors):
            waived.add(gid)
        if gid == "FINAL-QUALIFIED" and status != "VERIFIED" and validate_final_qualification_binding(work_slice_id, active.get("issue"), work, gap, errors):
            bound.add(gid)
    if work_slice_id == "ADAPT-CI-CONVERGENCE-001":
        missing = sorted(REQUIRED_70_GAPS - seen)
        if missing:
            errors.append("#70 closure ledger missing mandatory gap ids: " + ", ".join(missing))
    unresolved = [
        str(g.get("id"))
        for g in gaps
        if isinstance(g, dict)
        and g.get("blocksIssueClosure") is True
        and g.get("status") != "VERIFIED"
        and str(g.get("id")) not in waived
        and str(g.get("id")) not in bound
    ]
    work_status = str(work.get("status", active.get("status", ""))).strip().upper()
    if (require_closed or work_status in CLOSED_WORK_SLICE_STATES) and unresolved:
        errors.append("work-slice closure blocked by unresolved gaps: " + ", ".join(sorted(unresolved)))
    if work_status in CLOSED_WORK_SLICE_STATES and work.get("blocksNextProductCapability") is not False:
        errors.append("closed work slice must set blocksNextProductCapability=false")
    return errors, statuses, work_slice_id, waived, bound


def main() -> int:
    parser = argparse.ArgumentParser(description="Validate DE.PULSE active work-slice executable closure ledger")
    parser.add_argument("--require-closed", action="store_true")
    args = parser.parse_args()
    errors, statuses, work_slice_id, waived, bound = validate(args.require_closed)
    print(f"DE.PULSE active work-slice closure ledger: {work_slice_id or 'UNKNOWN'}")
    for status in sorted(ALLOWED_STATUSES):
        print(f"{status}: {statuses.get(status, 0)}")
    print(f"validated external-control waivers: {len(waived)}")
    print(f"validated post-run evidence bindings: {len(bound)}")
    for gid in sorted(waived):
        print(f"waiver-satisfied factual BLOCKED_EXTERNAL: {gid}")
    for gid in sorted(bound):
        print(f"runtime-evidence-satisfied static ledger state: {gid}")
    unresolved_count = sum(count for status, count in statuses.items() if status != "VERIFIED") - len(waived) - len(bound)
    print(f"unresolved blocking gaps: {max(unresolved_count, 0)}")
    print("documentation-only closure: PROHIBITED")
    print("special bindings are scope-limited and fail-closed: PASS")
    if errors:
        print("DE.PULSE work-slice closure ledger: FAIL", file=sys.stderr)
        for error in errors:
            print(" - " + error, file=sys.stderr)
        return 1
    if args.require_closed:
        print("all blocking gaps verified or narrowly evidence-satisfied: PASS")
    else:
        print("ledger completeness/enforcement contract: PASS")
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
