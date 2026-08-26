#!/usr/bin/env python3
"""Fail-closed T10 future-proof regression/conservation/portability assurance."""
from __future__ import annotations

import argparse
import json
from pathlib import Path
import subprocess
import sys

ROOT = Path(__file__).resolve().parents[2]
PROGRAM = ROOT / "governance" / "programs" / "ADAPT-V18-FINAL-CLOSURE-10-10-001"
T10 = PROGRAM / "T10_FUTURE_PROOF_ZERO_GAP_CERTIFICATION.json"
LEDGER = PROGRAM / "feature-assurance-ledger.json"
FREEZE = PROGRAM / "feature-assurance-ledger-freeze.json"
RECON = PROGRAM / "T1_FINAL_RECONCILIATION.json"
SCAN1 = PROGRAM / "T1_INDEPENDENT_OMISSION_SCAN.json"
SCAN2 = PROGRAM / "T1_INDEPENDENT_OMISSION_SCAN_2.json"
CURRENT = ROOT / "governance" / "current-state.json"
WORK_SLICE = ROOT / "governance" / "work-slices" / "ADAPT-V18-FINAL-CLOSURE-10-10-001" / "work-slice.json"
CLOSURE = ROOT / "governance" / "work-slices" / "ADAPT-V18-FINAL-CLOSURE-10-10-001" / "closure.json"
V19_GATE = ROOT / "tools" / "ci" / "requirement_conservation_gate.py"
V19_LEDGER = ROOT / "governance" / "programs" / "ADAPT-HOSTED-SYNC-001" / "requirement-conservation.json"
HANDOFF = ROOT / "handoff" / "CURRENT.md"
AGENTS = ROOT / "AGENTS.md"
CLAUDE = ROOT / "CLAUDE.md"
PORTABILITY = ROOT / "governance" / "AI-ASSISTANT-PORTABILITY-CONTRACT.md"
ADAPTIVE_RESUME = ROOT / Path("tools/ci/adaptive_resume_gate.py")
RETIRED_EQUIVALENCE = ROOT / "tools" / "ci" / "retired_test_equivalence_gate.py"
WORKFLOW_POLICY = ROOT / "tools" / "ci" / "workflow_policy.py"
IMPACT = ROOT / "tools" / "ci" / "impact_plan.py"
IMPACT_TEST = ROOT / "tools" / "ci" / "impact_plan_self_test.py"
CI_FAST = ROOT / ".github" / "workflows" / "ci-fast.yml"
CI_QUALIFIED = ROOT / ".github" / "workflows" / "ci-qualified.yml"
RELEASE = ROOT / ".github" / "workflows" / "release.yml"
MANIFEST = ROOT / "release" / "v18.10.0" / "certification-manifest.json"
VERSIONING = ROOT / "governance" / "versioning-policy.json"

PRIOR_TRACK_ROWS = (
    "T1-FEATURE-TRACEABILITY",
    "T2-UNIT-CONTRACT-PROPERTY",
    "T3-FUNCTIONAL-E2E",
    "T4-EDGE-FAILURE-DATA-TRUTH",
    "T5-PERSISTENCE-LIFECYCLE-RECOVERY",
    "T6-SECURITY-ROLES-RIGHTS",
    "T7-UI-UX-IA-CONTENT",
    "T8-PERFORMANCE-CONCURRENCY-SOAK",
    "T9-PACKAGED-CROSS-PLATFORM-RELEASE",
)
EXPECTED_OPEN_GAPS = {
    "T10-EXACT-HEAD-FAST",
    "T10-IDENTICAL-HEAD-QUALIFIED",
    "T10-G11-G16-STABLE-PUBLICATION",
}
EXECUTABLE_SUFFIXES = (".go", ".js", ".py", ".sh", ".ps1")


def load(path: Path) -> dict:
    return json.loads(path.read_text(encoding="utf-8-sig"))


def text(path: Path) -> str:
    return path.read_text(encoding="utf-8", errors="replace")


def parse_core(value: str) -> tuple[int, int, int]:
    raw = str(value).strip().removeprefix("v").removesuffix("-stable")
    parts = raw.split(".")
    if len(parts) != 3 or not all(p.isdigit() for p in parts):
        raise ValueError(value)
    return tuple(int(p) for p in parts)


def stable_tag(version: str) -> str:
    policy = load(VERSIONING)
    v = str(version).strip().removeprefix("v").removesuffix("-stable")
    cutover = str(policy.get("effectiveAfterProductVersion", "18.9.1"))
    pattern = policy.get("legacyStableTagPattern") if parse_core(v) <= parse_core(cutover) else policy.get("futureStableTagPattern")
    return str(pattern).format(productVersion=v)


def closure_row(closure: dict, row_id: str) -> dict | None:
    return next((row for row in closure.get("gaps", []) if isinstance(row, dict) and row.get("id") == row_id), None)


def governed_status(product: dict, track: str) -> str:
    return next((str(row.get("status") or "") for row in product.get("nextGovernedTracks", [])
                 if isinstance(row, dict) and row.get("track") == track), "")


def owner_exists(owner: str) -> bool:
    value = str(owner or "").strip()
    if not value or value.startswith("#"):
        return False
    path = ROOT / value
    return path.is_file() and path.suffix.lower() in EXECUTABLE_SUFFIXES


def regression_ownership_errors() -> tuple[list[str], int]:
    errors: list[str] = []
    ledger = load(LEDGER)
    scan1 = load(SCAN1)
    scan2 = load(SCAN2)
    recon = load(RECON)
    freeze = load(FREEZE)

    features = ledger.get("features") or []
    omissions1 = scan1.get("omissionsFound") or []
    omissions2 = scan2.get("omissionsFound") or []
    exclusions = recon.get("excludedFutureSourceCarryForward") or []
    excluded_ids = {str(row.get("id") or "") for row in exclusions if isinstance(row, dict)}
    physical = {str(row.get("id") or ""): row for row in features if isinstance(row, dict) and row.get("id")}
    parent_overrides = recon.get("scanParentOverrides") or {}
    parent_by_category = recon.get("scanParentByCategory") or {}

    effective_ids: list[str] = []
    for row in features:
        if not isinstance(row, dict):
            errors.append("feature-assurance ledger contains a non-object row")
            continue
        row_id = str(row.get("id") or "")
        if row_id in excluded_ids:
            continue
        effective_ids.append(row_id)
        owners = [str(x) for x in row.get("durableRegressionOwner") or []]
        if "#123" not in owners:
            errors.append(f"{row_id}: T10/#123 durable regression binding missing")
        if not any(owner_exists(owner) for owner in owners):
            errors.append(f"{row_id}: no existing executable durable regression owner")

    for source_name, rows in (("scan1", omissions1), ("scan2", omissions2)):
        for row in rows:
            if not isinstance(row, dict):
                errors.append(f"{source_name}: omission scan contains a non-object row")
                continue
            row_id = str(row.get("id") or "")
            category = str(row.get("category") or "")
            effective_ids.append(row_id)
            parent_id = str(parent_overrides.get(row_id) or parent_by_category.get(category) or "")
            parent = physical.get(parent_id)
            if not parent_id or not isinstance(parent, dict):
                errors.append(f"{row_id}: omission-scan responsibility has no valid canonical parent")
                continue
            tests = [str(x) for x in row.get("tests") or []]
            if not tests:
                tests = [str(x) for x in parent.get("existingRegressionOwners") or []]
            if not tests:
                errors.append(f"{row_id}: omission-scan responsibility and canonical parent have no regression tests")
            elif not any(owner_exists(test) for test in tests):
                errors.append(f"{row_id}: omission-scan regression owner no longer exists")
            durable = [str(x) for x in parent.get("durableRegressionOwner") or []]
            if "#123" not in durable:
                errors.append(f"{row_id}: canonical parent {parent_id} lost T10/#123 durable binding")
            if not any(owner_exists(owner) for owner in durable):
                errors.append(f"{row_id}: canonical parent {parent_id} has no executable durable regression owner")

    if len(effective_ids) != 180:
        errors.append(f"effective T1 regression responsibility count drifted: {len(effective_ids)} != 180")
    if len(set(effective_ids)) != len(effective_ids):
        errors.append("effective T1 regression responsibility IDs are not unique")
    if (freeze.get("effectiveInventory") or {}).get("effectiveShippedV18Responsibilities") != 180:
        errors.append("immutable T1 freeze no longer records 180 effective shipped-v18 responsibilities")
    if recon.get("effectiveShippedV18ResponsibilityCount") != 180 or recon.get("unexplainedGapCount") != 0:
        errors.append("T1 final reconciliation must remain 180 responsibilities / zero unexplained gaps")
    return errors, len(effective_ids)


def main() -> int:
    parser = argparse.ArgumentParser()
    parser.add_argument("--strict", action="store_true")
    args = parser.parse_args()
    errors: list[str] = []

    required = (
        T10, LEDGER, FREEZE, RECON, SCAN1, SCAN2, CURRENT, WORK_SLICE, CLOSURE,
        V19_GATE, V19_LEDGER, HANDOFF, AGENTS, CLAUDE, PORTABILITY, ADAPTIVE_RESUME,
        RETIRED_EQUIVALENCE, WORKFLOW_POLICY, IMPACT, IMPACT_TEST, CI_FAST,
        CI_QUALIFIED, RELEASE, MANIFEST, VERSIONING,
    )
    for path in required:
        if not path.is_file():
            errors.append(f"required T10 owner missing: {path.relative_to(ROOT)}")
    if errors:
        print("V18 T10 ASSURANCE GATE: FAIL", file=sys.stderr)
        for error in errors:
            print(" - " + error, file=sys.stderr)
        return 1

    t10 = load(T10)
    current = load(CURRENT)
    work_slice = load(WORK_SLICE)
    closure = load(CLOSURE)
    product = current.get("productCapabilityGate") or {}
    state = str(t10.get("state") or "")

    if t10.get("schema") != "DE.PULSE-V18-T10-FUTURE-PROOF-ZERO-GAP-CERTIFICATION-1":
        errors.append("T10 assurance schema mismatch")
    if t10.get("programIssue") != 113 or t10.get("trackIssue") != 123 or t10.get("targetVersion") != "v18.10.0":
        errors.append("T10 assurance identity must remain #113/#123/v18.10.0")
    if state not in {"IN_PROGRESS", "RELEASE_AUTHORIZED", "COMPLETE"}:
        errors.append(f"unsupported T10 state: {state!r}")

    for row_id in PRIOR_TRACK_ROWS:
        row = closure_row(closure, row_id)
        if not isinstance(row, dict) or row.get("status") != "VERIFIED":
            errors.append(f"T10 requires immutable VERIFIED prior closure: {row_id}")

    completed = {str(row.get("track")): row for row in product.get("completedChildTracks") or [] if isinstance(row, dict)}
    for number in range(1, 10):
        track = f"T{number}"
        row = completed.get(track)
        if not row or row.get("status") != "COMPLETE" or not row.get("fastRunId") or not row.get("qualifiedRunId"):
            errors.append(f"T10 requires durable exact-head Fast/Qualified completion for {track}")

    ownership_errors, ownership_count = regression_ownership_errors()
    errors.extend(ownership_errors)

    conservation = subprocess.run([sys.executable, str(V19_GATE)], cwd=ROOT, check=False, text=True, capture_output=True)
    if conservation.returncode != 0:
        errors.append("permanent v19/#66 requirement-conservation gate failed: " + conservation.stdout.strip() + conservation.stderr.strip())
    v19 = load(V19_LEDGER)
    if len(v19.get("requirements") or []) != 72:
        errors.append("T10 requires exactly 72 conserved v19/#66 rows")

    ci_fast = text(CI_FAST)
    for marker in (
        "python3 tools/ci/requirement_conservation_gate.py",
        "python3 tools/ci/v18_t10_future_proof_zero_gap_certification_gate.py",
        "python3 tools/ci/adaptive_resume_gate.py",
        "python3 tools/ci/retired_test_equivalence_gate.py",
    ):
        if marker not in ci_fast:
            errors.append(f"CI Fast lost T10 fail-closed owner: {marker}")

    manifest = load(MANIFEST)
    python_gates = {tuple(row) for row in manifest.get("pythonGates", []) if isinstance(row, list)}
    for gate in (
        ("python3", "tools/ci/requirement_conservation_gate.py"),
        ("python3", "tools/ci/v18_t10_future_proof_zero_gap_certification_gate.py"),
        ("python3", "tools/ci/adaptive_resume_gate.py"),
        ("python3", "tools/ci/retired_test_equivalence_gate.py"),
    ):
        if gate not in python_gates:
            errors.append(f"v18.10 certification manifest lost required T10 gate: {gate[1]}")

    impact = text(IMPACT)
    impact_test = text(IMPACT_TEST)
    for marker in ("T9_CROSS_LAYER_ASSURANCE_FILES", "release/v18.10.0/certification-manifest.json", "native_macos_required", "native_windows_required"):
        if marker not in impact:
            errors.append(f"Planner v3 lost final-v18 full-evidence routing marker: {marker}")
    for marker in ("T9 manifest must select both native packages", "unknown path must select full evidence graph"):
        if marker not in impact_test:
            errors.append(f"Planner v3 self-test lost final-v18 fail-closed invariant: {marker}")

    for adapter_path in (AGENTS, CLAUDE):
        body = text(adapter_path)
        if "governance/AI-ASSISTANT-PORTABILITY-CONTRACT.md" not in body or "handoff/CURRENT.md" not in body:
            errors.append(f"{adapter_path.name} no longer points to GitHub-only portability owners")
    portability = text(PORTABILITY)
    for marker in ("GitHub source-of-truth hierarchy", "Mandatory fresh-session algorithm", "No upload of an old chat handoff is required"):
        if marker not in portability:
            errors.append(f"portability contract lost required T10 marker: {marker}")
    handoff = text(HANDOFF)
    if "SUPERSEDES ALL PRIOR CHAT HANDOFFS" not in handoff or "Exactly one next action" not in handoff:
        errors.append("current handoff is not a single portable GitHub resume authority")

    t10_row = closure_row(closure, "T10-FUTURE-PROOF-ZERO-GAP-CERTIFICATION")
    if not isinstance(t10_row, dict):
        errors.append("T10 closure row missing")
        t10_row = {}
    gap_rows = [row for row in t10.get("knownCoverageGaps", []) if isinstance(row, dict)]
    open_gap_ids = {str(row.get("id") or "") for row in gap_rows if row.get("status") == "OPEN"}

    if state == "IN_PROGRESS":
        if work_slice.get("nextTrack") != {"track": "T10", "issue": 123}:
            errors.append("IN_PROGRESS T10 requires work-slice nextTrack T10/#123")
        if product.get("nextChildIssue") != 123 or product.get("nextChildTrack") != "T10":
            errors.append("IN_PROGRESS T10 requires current-state active child T10/#123")
        if governed_status(product, "T10") != "IN_PROGRESS":
            errors.append("IN_PROGRESS T10 requires current-state nextGovernedTracks T10 IN_PROGRESS")
        if (product.get("downstreamAssuranceState") or {}).get("T10Started") is not True:
            errors.append("IN_PROGRESS T10 requires downstreamAssuranceState.T10Started=true")
        if t10_row.get("status") not in {"OPEN", "IMPLEMENTED_UNVERIFIED"}:
            errors.append("IN_PROGRESS T10 closure row must remain OPEN or IMPLEMENTED_UNVERIFIED")
        if open_gap_ids != EXPECTED_OPEN_GAPS:
            errors.append(f"IN_PROGRESS T10 open gaps must be exactly {sorted(EXPECTED_OPEN_GAPS)}")
        if t10.get("releaseAuthorized") is not False or t10.get("stablePublished") is not False:
            errors.append("IN_PROGRESS T10 may not authorize release or claim Stable publication")
    elif state == "RELEASE_AUTHORIZED":
        if open_gap_ids != {"T10-G11-G16-STABLE-PUBLICATION"}:
            errors.append("RELEASE_AUTHORIZED T10 must retain only final G11-G16 Stable-publication gap")
        if t10.get("releaseAuthorized") is not True or t10.get("stablePublished") is not False:
            errors.append("RELEASE_AUTHORIZED requires releaseAuthorized=true and stablePublished=false")
        if t10_row.get("status") != "IMPLEMENTED_UNVERIFIED":
            errors.append("RELEASE_AUTHORIZED T10 remains IMPLEMENTED_UNVERIFIED until publication/G16")
        if not t10.get("finalSourceSha") or not t10.get("fastRunId") or not t10.get("qualifiedRunId"):
            errors.append("RELEASE_AUTHORIZED T10 requires exact-head Fast/Qualified evidence")
    elif state == "COMPLETE":
        if gap_rows:
            errors.append("COMPLETE T10 cannot retain coverage gaps")
        if t10.get("releaseAuthorized") is not True or t10.get("stablePublished") is not True:
            errors.append("COMPLETE T10 requires release authorization and Stable publication")
        if t10_row.get("status") != "VERIFIED":
            errors.append("COMPLETE T10 closure row must be VERIFIED")
        stable = current.get("stable") or {}
        if stable.get("productVersion") != "18.10.0" or stable.get("tag") != stable_tag("18.10.0"):
            errors.append("COMPLETE T10 requires current-state Stable v18.10.0 publication under canonical tag policy")

        reservation_complete = product.get("reservationStatus") in {"COMPLETE", "CLOSED"}
        post_audit_pass = (
            product.get("postClosureSourceOverlapAuditRequired") is True
            and str(product.get("postClosureSourceOverlapAuditStatus") or "").upper() == "PASS"
            and bool(product.get("postClosureSourceOverlapAuditMergeSha"))
        )
        later_reservation_active = (
            product.get("reservedWorkSliceId") != "ADAPT-V18-FINAL-CLOSURE-10-10-001"
            and product.get("reservationStatus") in {"ACTIVE", "IN_PROGRESS", "G1_RESERVED"}
            and product.get("reservedIssue") not in {None, 113, 123}
            and bool(product.get("reservedBranch"))
        )
        if not reservation_complete and not (post_audit_pass and later_reservation_active):
            errors.append(
                "COMPLETE T10 requires either the completed v18.10 closure reservation or a later governed reservation after PASS post-v18 source-overlap audit"
            )
        if later_reservation_active:
            t10_completed = completed.get("T10")
            if not t10_completed or t10_completed.get("status") != "COMPLETE" or not t10_completed.get("releaseRunId"):
                errors.append("post-v18 reservation requires durable completed T10 child evidence")
            downstream = product.get("downstreamAssuranceState") or {}
            if downstream.get("T10State") != "COMPLETE" or downstream.get("T10ReleaseRunId") != t10.get("releaseRunId"):
                errors.append("post-v18 reservation requires current-state T10 COMPLETE / release-run binding")
            if product.get("nextChildIssue") is not None or product.get("nextChildTrack") is not None:
                errors.append("post-v18 reservation may not reopen completed v18 T10 child progression")

        if governed_status(product, "T10") not in {"COMPLETE", "VERIFIED"}:
            errors.append("COMPLETE T10 requires current-state T10 final status")
        if not t10.get("releaseRunId") or not t10.get("g16Evidence"):
            errors.append("COMPLETE T10 requires durable Release/G16 evidence")
        if product.get("postClosureSourceOverlapAuditRequired") is not True:
            errors.append("COMPLETE T10 must preserve mandatory post-v18.10 source-overlap audit before v19 G1")

    if t10.get("effectiveRegressionResponsibilityCount") != 180:
        errors.append("T10 artifact must record 180 durable regression responsibilities")
    if t10.get("v19ConservedRequirementCount") != 72:
        errors.append("T10 artifact must record 72 v19 conserved requirements")
    ordering = t10.get("publicationOrdering") or {}
    if ordering.get("prePublicationState") != "RELEASE_AUTHORIZED" or ordering.get("postPublicationState") != "COMPLETE":
        errors.append("T10 publication ordering must avoid circular completion while preserving post-publication verification")

    strict = args.strict or state == "COMPLETE"
    if strict and gap_rows:
        errors.append("strict T10 closure cannot retain coverage gaps")

    print("V18 T10 FUTURE-PROOF / ZERO-GAP / PORTABILITY ASSURANCE")
    print(f"state: {state}")
    print(f"durable v18 regression responsibilities: {ownership_count}/180")
    print("v19/#66 conserved rows: 72")
    print("portable resume: GitHub-only contract + current handoff + canonical resume gate")
    print("retired-test equivalence: permanent fail-closed control retained")
    print(f"open gaps: {len(open_gap_ids)}")
    if errors:
        print("V18 T10 ASSURANCE GATE: FAIL", file=sys.stderr)
        for error in errors:
            print(" - " + error, file=sys.stderr)
        return 1
    print(f"V18 T10 ASSURANCE GATE: PASS (strict={strict})")
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
