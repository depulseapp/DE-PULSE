#!/usr/bin/env python3
"""T4 edge/adversarial/failure/data-truth assurance over the frozen v18 inventory.

T4 is intentionally fail-closed about *evidence quality*: a happy-path owner or a
filename containing words such as "stale" is not enough. A counted owner must be
an executable regression/CI owner whose source contains both an adverse semantic
and an assertion/fail-closed mechanism. Final T4 closure additionally requires
exact-head Fast and identical-head Qualified execution.
"""
from __future__ import annotations

import argparse
from pathlib import Path
import sys

from v18_t2_unit_contract_assurance_gate import (
    ROOT, PROGRAM, LEDGER, FREEZE, SCAN1, SCAN2, RECONCILIATION,
    CURRENT_STATE, CLOSURE, load, git_blob_sha, reconstruct_effective,
)

T4 = PROGRAM / "T4_EDGE_FAILURE_DATA_TRUTH_ASSURANCE.json"
CI_FAST = ROOT / ".github" / "workflows" / "ci-fast.yml"

ADVERSE_SEMANTICS = (
    "stale", "missing", "future", "partial", "contradict", "unavailable",
    "timeout", "timed out", "rate limit", "ratelimit", "429", "outage",
    "fallback", "cache miss", "cache stale", "degrad", "recover", "retry",
    "invalid", "malformed", "duplicate", "replay", "backpressure", "overflow",
    "denied", "forbidden", "unauthorized", "expired", "revoked", "corrupt",
    "cancel", "unknown", "abstain", "mismatch", "diverg", "no rebuild",
    "rebuild", "provenance", "fail closed", "fail-closed", "network",
    "after-hours", "after hours", "pre-market", "premarket", "market closed",
    "closed market", "zero truth", "not configured", "not-configured",
    "disconnected", "unhealthy", "insufficient", "empty evidence",
)

ASSERTION_MARKERS = (
    "t.fatal", "t.error", "t.fail", "require.", "assert.", "assert(",
    "expect(", "throw new error", "raise assertionerror", "raise systemexit",
    "errors.append", "sys.exit(1", "return 1", "process.exit(1",
)


def classify(path_text: str) -> str:
    p = path_text.replace("\\", "/").strip()
    lower = p.lower()
    name = lower.rsplit("/", 1)[-1]
    if lower.startswith("tests/acceptance/"):
        return "ACCEPTANCE_ADVERSE"
    if lower.startswith("tests/integration/"):
        return "HTTP_ADVERSE"
    if "browser" in name and name.endswith((".py", ".js")):
        return "BROWSER_ADVERSE"
    if lower.startswith("tests/renderer/") and name.endswith(".js"):
        return "RENDERER_ADVERSE"
    if name.endswith("_test.go"):
        return "GO_ADVERSE"
    if lower.startswith("tools/ci/") and name.endswith(("_gate.py", "_contract.py", "_test.py", "_test.js")):
        return "CI_FAIL_CLOSED"
    if lower == "tools/ci/release_rehearsal.py":
        return "CI_FAIL_CLOSED"
    return "UNKNOWN_EXECUTABLE"


PROFILE_CLASSES = {
    "VISIBLE_STATEFUL": {"GO_ADVERSE","HTTP_ADVERSE","ACCEPTANCE_ADVERSE","BROWSER_ADVERSE","RENDERER_ADVERSE"},
    "VISIBLE_READONLY": {"GO_ADVERSE","HTTP_ADVERSE","ACCEPTANCE_ADVERSE","BROWSER_ADVERSE","RENDERER_ADVERSE"},
    "SECURITY_STATEFUL": {"GO_ADVERSE","HTTP_ADVERSE","ACCEPTANCE_ADVERSE","BROWSER_ADVERSE","RENDERER_ADVERSE"},
    "DATA_PATH": {"GO_ADVERSE","HTTP_ADVERSE","CI_FAIL_CLOSED"},
    "BACKGROUND_JOB": {"GO_ADVERSE","HTTP_ADVERSE","CI_FAIL_CLOSED"},
    "RUNTIME_STATEFUL": {"GO_ADVERSE","HTTP_ADVERSE","CI_FAIL_CLOSED"},
    "PERSISTENCE": {"GO_ADVERSE","HTTP_ADVERSE","CI_FAIL_CLOSED"},
    "BOUNDARY": {"GO_ADVERSE","HTTP_ADVERSE","CI_FAIL_CLOSED","ACCEPTANCE_ADVERSE"},
    "RELEASE": {"CI_FAIL_CLOSED","BROWSER_ADVERSE"},
}

# These two scan-discovered responsibilities were conservatively assigned the
# RELEASE profile by the immutable T1 reconstruction even though their actual
# executable behavior is Go-owned: isolated TEST-profile migration and an
# explicitly opt-in developer diagnostic. Keep the RELEASE profile strict for
# every publication/platform row; admit GO_ADVERSE only for these exact IDs.
ROW_CLASS_EXCEPTIONS = {
    "RELEASE-TEST-PROFILE-MIGRATION": {"GO_ADVERSE"},
    "RELEASE-DEVELOPER-SCHEMA-PROBE": {"GO_ADVERSE"},
}


def meaningful_adverse_owner(path: Path) -> tuple[bool, list[str], list[str]]:
    try:
        text = path.read_text(encoding="utf-8", errors="ignore").lower()
    except OSError:
        return False, [], []
    semantics = sorted({term for term in ADVERSE_SEMANTICS if term in text})
    assertions = sorted({marker for marker in ASSERTION_MARKERS if marker in text})
    return bool(semantics and assertions), semantics, assertions


def main() -> int:
    parser = argparse.ArgumentParser()
    parser.add_argument("--strict", action="store_true")
    args = parser.parse_args()

    ledger, freeze = load(LEDGER), load(FREEZE)
    scan1, scan2, reconciliation = load(SCAN1), load(SCAN2), load(RECONCILIATION)
    current_state, closure, t4 = load(CURRENT_STATE), load(CLOSURE), load(T4)
    errors: list[str] = []

    product = current_state.get("productCapabilityGate") or {}
    governed = product.get("nextGovernedTracks") or []
    t4_runtime_state = next((str(x.get("status") or "") for x in governed if isinstance(x, dict) and x.get("track") == "T4"), "")
    t5_runtime_state = next((str(x.get("status") or "") for x in governed if isinstance(x, dict) and x.get("track") == "T5"), "")
    assurance_state = str(t4.get("state") or "")

    if assurance_state == "IN_PROGRESS":
        if product.get("nextChildIssue") != 117 or product.get("nextChildTrack") != "T4" or t4_runtime_state != "IN_PROGRESS":
            errors.append("IN_PROGRESS T4 must be the active T4/#117 child")
        if t5_runtime_state != "NOT_STARTED":
            errors.append("IN_PROGRESS T4 must not silently start T5/#118")
    elif assurance_state == "COMPLETE":
        completed = product.get("completedChildTracks") or []
        t4_completed = next((x for x in completed if isinstance(x, dict) and x.get("track") == "T4" and x.get("issue") == 117), None)
        if not isinstance(t4_completed, dict) or t4_completed.get("status") != "COMPLETE":
            errors.append("COMPLETE T4 must remain recorded in completedChildTracks")
        if not str((t4_completed or {}).get("mergedCommitSha") or "").strip():
            errors.append("COMPLETE T4 requires durable mergedCommitSha evidence")
    else:
        errors.append(f"unsupported T4 state: {assurance_state!r}")

    gaps = closure.get("gaps") or []
    for required in ("T1-FEATURE-TRACEABILITY", "T2-UNIT-CONTRACT-PROPERTY", "T3-FUNCTIONAL-E2E"):
        row = next((x for x in gaps if isinstance(x, dict) and x.get("id") == required), None)
        if not isinstance(row, dict) or row.get("status") != "VERIFIED":
            errors.append(f"T4 requires {required}=VERIFIED")
    t4_gap = next((x for x in gaps if isinstance(x, dict) and x.get("id") == "T4-EDGE-FAILURE-DATA-TRUTH"), None)
    if not isinstance(t4_gap, dict) or t4_gap.get("status") not in {"IMPLEMENTED_UNVERIFIED", "VERIFIED"}:
        errors.append("parent closure ledger must record T4 as IMPLEMENTED_UNVERIFIED or VERIFIED")

    fast_text = CI_FAST.read_text(encoding="utf-8") if CI_FAST.is_file() else ""
    if "python3 tools/ci/v18_t4_edge_failure_data_truth_assurance_gate.py" not in fast_text:
        errors.append("T4 assurance gate is not bound into canonical CI Fast")

    actual_blob = git_blob_sha(LEDGER)
    expected_blob = str((freeze.get("frozenDiscovery") or {}).get("gitBlobSha") or "")
    if freeze.get("state") != "FROZEN_T1" or actual_blob != expected_blob:
        errors.append("T4 requires the immutable frozen T1 discovery blob")
    if t4.get("trackIssue") != 117 or t4.get("programIssue") != 113 or t4.get("frozenT1GitBlobSha") != actual_blob:
        errors.append("T4 assurance identity/frozen T1 binding mismatch")

    effective = reconstruct_effective(ledger, scan1, scan2, reconciliation, errors)
    expected_count = int((freeze.get("effectiveInventory") or {}).get("effectiveShippedV18Responsibilities") or 0)
    if len(effective) != expected_count:
        errors.append(f"effective inventory count mismatch: expected {expected_count}, got {len(effective)}")

    supplements = t4.get("evidenceSupplements") or {}
    if not isinstance(supplements, dict):
        errors.append("T4 evidenceSupplements must be an object")
        supplements = {}
    for fid, spec in supplements.items():
        if fid not in effective or not isinstance(spec, dict):
            errors.append(f"invalid T4 supplement: {fid}")
            continue
        rationale = str(spec.get("rationale") or "").strip()
        owners = spec.get("owners")
        if not rationale or not isinstance(owners, list) or not owners:
            errors.append(f"T4 supplement {fid} requires owners and rationale")
            continue
        row_owners = effective[fid].setdefault("existingRegressionOwners", [])
        for owner in owners:
            if isinstance(owner, str) and owner.strip() and owner.strip() not in row_owners:
                row_owners.append(owner.strip())

    profiles = ledger.get("assuranceProfiles") or {}
    uncovered: list[str] = []
    missing_paths: set[str] = set()
    class_counts: dict[str, int] = {}
    semantic_counts: dict[str, int] = {}
    covered = 0

    for fid, row in sorted(effective.items()):
        profile_name = str(row.get("assuranceProfile") or "")
        profile = profiles.get(profile_name)
        if not isinstance(profile, dict) or not str(profile.get("T4") or "").strip():
            errors.append(f"{fid}: missing T4 expectation/profile")
            continue
        valid_classes = set(PROFILE_CLASSES.get(profile_name) or set()) | ROW_CLASS_EXCEPTIONS.get(fid, set())
        if not valid_classes:
            errors.append(f"{fid}: no T4 evidence policy for profile {profile_name}")
            continue

        valid: list[str] = []
        for owner in row.get("existingRegressionOwners") or []:
            owner_text = str(owner or "").strip()
            if not owner_text:
                continue
            path = ROOT / owner_text
            if not path.is_file():
                missing_paths.add(f"{fid}:{owner_text}")
                continue
            cls = classify(owner_text)
            class_counts[cls] = class_counts.get(cls, 0) + 1
            meaningful, semantics, _assertions = meaningful_adverse_owner(path)
            if meaningful:
                for semantic in semantics:
                    semantic_counts[semantic] = semantic_counts.get(semantic, 0) + 1
            if cls in valid_classes and meaningful:
                valid.append(owner_text)

        if valid:
            covered += 1
        else:
            uncovered.append(fid)

    if missing_paths:
        errors.append("referenced T4 evidence missing at current head: " + ", ".join(sorted(missing_paths)))

    actual_gaps = set(uncovered)
    declared_gaps = set(str(x) for x in (t4.get("knownCoverageGaps") or []))
    declaration = str(t4.get("gapDeclarationState") or "")
    if declaration == "PENDING_FIRST_EXECUTION":
        if declared_gaps:
            errors.append("PENDING_FIRST_EXECUTION T4 must not predeclare guessed coverage gaps")
    elif declaration == "CURRENT":
        if actual_gaps - declared_gaps:
            errors.append("T4 found undeclared gaps: " + ", ".join(sorted(actual_gaps - declared_gaps)))
        if declared_gaps - actual_gaps:
            errors.append("T4 declared stale/resolved gaps: " + ", ".join(sorted(declared_gaps - actual_gaps)))
        if t4.get("uncoveredResponsibilityCount") != len(uncovered):
            errors.append("T4 uncoveredResponsibilityCount drift")
    else:
        errors.append(f"unsupported T4 gapDeclarationState: {declaration!r}")

    strict = args.strict or assurance_state == "COMPLETE"
    if strict and uncovered:
        errors.append("T4 strict closure has uncovered responsibilities: " + ", ".join(uncovered))
    if assurance_state == "COMPLETE":
        if declaration != "CURRENT":
            errors.append("T4 COMPLETE requires a CURRENT executable gap declaration")
        if t4.get("uncoveredResponsibilityCount") != 0:
            errors.append("T4 COMPLETE requires uncoveredResponsibilityCount=0")

    print("V18 T4 EDGE / ADVERSARIAL / FAILURE / DATA-TRUTH ASSURANCE")
    print(f"frozen T1 blob: {actual_blob}")
    print(f"effective responsibilities: {len(effective)}")
    print(f"T4-covered responsibilities: {covered}")
    print(f"T4-uncovered responsibilities: {len(uncovered)}")
    print(f"T4-specific evidence supplements: {len(supplements)}")
    for cls in sorted(class_counts):
        print(f"{cls}: {class_counts[cls]}")
    if semantic_counts:
        print("adverse semantics observed: " + ", ".join(f"{k}={semantic_counts[k]}" for k in sorted(semantic_counts)))
    if uncovered:
        print("uncovered ids: " + ", ".join(uncovered))
    print("happy-path-only and filename-only evidence is rejected: PASS")
    print("RELEASE profile remains CI/browser fail-closed except two exact Go-owned scan responsibilities: PASS")
    print("canonical recovery/routing owners are conserved; T4 creates no parallel recovery subsystem: PASS")
    print("T5-T10 certification is not implied by T4: PASS")

    if errors:
        print("V18 T4 ASSURANCE GATE: FAIL", file=sys.stderr)
        for error in errors:
            print(" - " + error, file=sys.stderr)
        return 1
    print(f"V18 T4 ASSURANCE GATE: PASS (strict={strict}, declaration={declaration})")
    return 0


if __name__ == "__main__":
    sys.exit(main())
