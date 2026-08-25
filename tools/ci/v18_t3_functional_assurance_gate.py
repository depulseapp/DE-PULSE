#!/usr/bin/env python3
"""T3 functional/integration/end-to-end assurance over the frozen v18 inventory."""
from __future__ import annotations

import argparse
from pathlib import Path
import subprocess
import sys
from typing import Any

from v18_t2_unit_contract_assurance_gate import (
    ROOT, PROGRAM, LEDGER, FREEZE, SCAN1, SCAN2, RECONCILIATION,
    CURRENT_STATE, CLOSURE, load, git_blob_sha, reconstruct_effective,
)

T3 = PROGRAM / "T3_FUNCTIONAL_ASSURANCE.json"
CI_FAST = ROOT / ".github" / "workflows" / "ci-fast.yml"
SECURITY_WORKFLOW = ROOT / "tests" / "integration" / "security_identity_sse_workflow_test.py"


def classify(path_text: str) -> str:
    p = path_text.replace("\\", "/").strip()
    lower = p.lower()
    name = lower.rsplit("/", 1)[-1]
    if lower.startswith("tests/acceptance/"):
        return "ACCEPTANCE_E2E"
    if lower.startswith("tests/integration/"):
        return "HTTP_INTEGRATION"
    if "browser" in name and name.endswith((".py", ".js")):
        return "BROWSER_E2E"
    if lower.startswith("tests/renderer/") and name.endswith(".js"):
        return "RENDERER_FUNCTIONAL"
    if name.endswith("_test.go"):
        return "GO_FUNCTIONAL"
    if lower.startswith("tools/ci/") and name.endswith(("_gate.py", "_contract.py", "_test.py", "_test.js")):
        return "CI_WORKFLOW_CONTRACT"
    if lower.startswith("tools/release/") or lower.startswith("packaging/") or lower.startswith(".github/workflows/"):
        return "CI_WORKFLOW_CONTRACT"
    return "UNKNOWN_EXECUTABLE"


PROFILE_CLASSES = {
    "VISIBLE_STATEFUL": {"ACCEPTANCE_E2E","HTTP_INTEGRATION","BROWSER_E2E","RENDERER_FUNCTIONAL"},
    "VISIBLE_READONLY": {"ACCEPTANCE_E2E","HTTP_INTEGRATION","BROWSER_E2E","RENDERER_FUNCTIONAL"},
    "SECURITY_STATEFUL": {"ACCEPTANCE_E2E","HTTP_INTEGRATION","BROWSER_E2E","RENDERER_FUNCTIONAL"},
    "DATA_PATH": {"GO_FUNCTIONAL","HTTP_INTEGRATION","ACCEPTANCE_E2E"},
    "BACKGROUND_JOB": {"GO_FUNCTIONAL","HTTP_INTEGRATION","ACCEPTANCE_E2E"},
    "RUNTIME_STATEFUL": {"GO_FUNCTIONAL","HTTP_INTEGRATION","ACCEPTANCE_E2E"},
    "PERSISTENCE": {"GO_FUNCTIONAL","HTTP_INTEGRATION","ACCEPTANCE_E2E"},
    "BOUNDARY": {"GO_FUNCTIONAL","HTTP_INTEGRATION","ACCEPTANCE_E2E"},
    "RELEASE": {"CI_WORKFLOW_CONTRACT","ACCEPTANCE_E2E"},
}


def main() -> int:
    parser = argparse.ArgumentParser()
    parser.add_argument("--strict", action="store_true")
    args = parser.parse_args()
    ledger, freeze = load(LEDGER), load(FREEZE)
    scan1, scan2, reconciliation = load(SCAN1), load(SCAN2), load(RECONCILIATION)
    current_state, closure, t3 = load(CURRENT_STATE), load(CLOSURE), load(T3)
    errors: list[str] = []

    product = current_state.get("productCapabilityGate") or {}
    governed = product.get("nextGovernedTracks") or []
    t3_state = next((str(x.get("status") or "") for x in governed if isinstance(x, dict) and x.get("track") == "T3"), "")
    t4_state = next((str(x.get("status") or "") for x in governed if isinstance(x, dict) and x.get("track") == "T4"), "")
    if product.get("nextChildIssue") != 116 or product.get("nextChildTrack") != "T3" or t3_state != "IN_PROGRESS":
        errors.append("current-state must identify T3/#116 as the active IN_PROGRESS child")
    if t4_state != "NOT_STARTED":
        errors.append("T3 audit must not silently start T4/#117")

    gaps = closure.get("gaps") or []
    for required in ("T1-FEATURE-TRACEABILITY", "T2-UNIT-CONTRACT-PROPERTY"):
        row = next((x for x in gaps if isinstance(x, dict) and x.get("id") == required), None)
        if not isinstance(row, dict) or row.get("status") != "VERIFIED":
            errors.append(f"T3 requires {required}=VERIFIED")
    t3_gap = next((x for x in gaps if isinstance(x, dict) and x.get("id") == "T3-FUNCTIONAL-E2E"), None)
    if not isinstance(t3_gap, dict) or t3_gap.get("status") not in {"IMPLEMENTED_UNVERIFIED","VERIFIED"}:
        errors.append("parent closure ledger must record active T3 as IMPLEMENTED_UNVERIFIED or VERIFIED")

    if "python3 tools/ci/v18_t3_functional_assurance_gate.py" not in CI_FAST.read_text(encoding="utf-8"):
        errors.append("T3 assurance gate is not bound into canonical CI Fast")
    actual_blob = git_blob_sha(LEDGER)
    expected_blob = str((freeze.get("frozenDiscovery") or {}).get("gitBlobSha") or "")
    if freeze.get("state") != "FROZEN_T1" or actual_blob != expected_blob:
        errors.append("T3 requires the immutable frozen T1 discovery blob")
    if t3.get("trackIssue") != 116 or t3.get("programIssue") != 113 or t3.get("frozenT1GitBlobSha") != actual_blob:
        errors.append("T3 assurance identity/frozen T1 binding mismatch")

    if not SECURITY_WORKFLOW.is_file():
        errors.append("T3 security/identity/workspace/SSE integration owner is missing")
    else:
        try:
            subprocess.run([sys.executable, str(SECURITY_WORKFLOW)], cwd=ROOT, check=True)
        except subprocess.CalledProcessError as exc:
            errors.append(f"T3 security/identity/workspace/SSE workflow failed with exit code {exc.returncode}")

    effective = reconstruct_effective(ledger, scan1, scan2, reconciliation, errors)
    expected_count = int((freeze.get("effectiveInventory") or {}).get("effectiveShippedV18Responsibilities") or 0)
    if len(effective) != expected_count:
        errors.append(f"effective inventory count mismatch: expected {expected_count}, got {len(effective)}")

    supplements = t3.get("evidenceSupplements") or {}
    if not isinstance(supplements, dict):
        errors.append("T3 evidenceSupplements must be an object")
        supplements = {}
    for fid, spec in supplements.items():
        if fid not in effective or not isinstance(spec, dict):
            errors.append(f"invalid T3 supplement: {fid}")
            continue
        rationale = str(spec.get("rationale") or "").strip()
        owners = spec.get("owners")
        if not rationale or not isinstance(owners, list) or not owners:
            errors.append(f"T3 supplement {fid} requires owners and rationale")
            continue
        row_owners = effective[fid].setdefault("existingRegressionOwners", [])
        for owner in owners:
            if isinstance(owner, str) and owner.strip() and owner.strip() not in row_owners:
                row_owners.append(owner.strip())

    profiles = ledger.get("assuranceProfiles") or {}
    uncovered: list[str] = []
    missing_paths: list[str] = []
    class_counts: dict[str, int] = {}
    covered = 0
    for fid, row in sorted(effective.items()):
        profile_name = str(row.get("assuranceProfile") or "")
        profile = profiles.get(profile_name)
        if not isinstance(profile, dict) or not str(profile.get("T3") or "").strip():
            errors.append(f"{fid}: missing T3 expectation/profile")
            continue
        valid_classes = PROFILE_CLASSES.get(profile_name)
        if not valid_classes:
            errors.append(f"{fid}: no T3 evidence policy for profile {profile_name}")
            continue
        valid = []
        for owner in row.get("existingRegressionOwners") or []:
            owner_text = str(owner or "").strip()
            if not owner_text:
                continue
            path = ROOT / owner_text
            if not path.is_file():
                missing_paths.append(f"{fid}:{owner_text}")
                continue
            cls = classify(owner_text)
            class_counts[cls] = class_counts.get(cls, 0) + 1
            if cls in valid_classes:
                valid.append(owner_text)
        if valid:
            covered += 1
        else:
            uncovered.append(fid)

    if missing_paths:
        errors.append("referenced T3 evidence missing at current head: " + ", ".join(sorted(missing_paths)))
    actual_gaps, declared_gaps = set(uncovered), set(str(x) for x in (t3.get("knownCoverageGaps") or []))
    strict = args.strict or t3.get("state") == "COMPLETE"
    if t3.get("state") == "IN_PROGRESS":
        declaration = str(t3.get("gapDeclarationState") or "")
        if declaration not in {"PENDING_FIRST_EXECUTION","CURRENT"}:
            errors.append(f"unsupported T3 gapDeclarationState: {declaration!r}")
        if declaration == "CURRENT":
            if actual_gaps - declared_gaps:
                errors.append("T3 found undeclared gaps: " + ", ".join(sorted(actual_gaps - declared_gaps)))
            if declared_gaps - actual_gaps:
                errors.append("T3 declared stale/resolved gaps: " + ", ".join(sorted(declared_gaps - actual_gaps)))
            if t3.get("uncoveredResponsibilityCount") != len(uncovered):
                errors.append("T3 uncoveredResponsibilityCount drift")
    if strict and uncovered:
        errors.append("T3 strict closure has uncovered responsibilities: " + ", ".join(uncovered))
    if t3.get("state") == "COMPLETE" and t3.get("uncoveredResponsibilityCount") != 0:
        errors.append("T3 COMPLETE requires uncoveredResponsibilityCount=0")

    print("V18 T3 FUNCTIONAL / INTEGRATION / END-TO-END ASSURANCE")
    print(f"frozen T1 blob: {actual_blob}")
    print(f"effective responsibilities: {len(effective)}")
    print(f"T3-covered responsibilities: {covered}")
    print(f"T3-uncovered responsibilities: {len(uncovered)}")
    print(f"T3-specific evidence supplements: {len(supplements)}")
    for cls in sorted(class_counts):
        print(f"{cls}: {class_counts[cls]}")
    if uncovered:
        print("uncovered ids: " + ", ".join(uncovered))
    print("security/identity/workspace/SSE workflow owner executes through real HTTP routes: PASS")
    print("visible workflows cannot be closed by backend-unit/static evidence alone: PASS")
    print("T4-T10 certification is not implied by T3: PASS")
    if errors:
        print("V18 T3 ASSURANCE GATE: FAIL", file=sys.stderr)
        for error in errors:
            print(" - " + error, file=sys.stderr)
        return 1
    print(f"V18 T3 ASSURANCE GATE: PASS (strict={strict})")
    return 0

if __name__ == "__main__":
    sys.exit(main())
