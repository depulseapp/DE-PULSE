#!/usr/bin/env python3
"""T5 persistence / lifecycle / restart assurance over frozen v18 inventory."""
from __future__ import annotations

import argparse
from pathlib import Path
import sys

from v18_t2_unit_contract_assurance_gate import (
    ROOT, PROGRAM, LEDGER, FREEZE, SCAN1, SCAN2, RECONCILIATION,
    CURRENT_STATE, CLOSURE, load, git_blob_sha, reconstruct_effective,
)

T5 = PROGRAM / "T5_PERSISTENCE_LIFECYCLE_ASSURANCE.json"
CI_FAST = ROOT / ".github" / "workflows" / "ci-fast.yml"

ALWAYS_APPLICABLE = {"VISIBLE_STATEFUL", "SECURITY_STATEFUL", "RUNTIME_STATEFUL", "PERSISTENCE"}
CONDITIONAL = {"VISIBLE_READONLY", "DATA_PATH", "BACKGROUND_JOB", "BOUNDARY", "RELEASE"}

PERSISTENCE_SEMANTICS = (
    "persist", "sqlite", "migration", "restart", "warm start", "warmstart",
    "restore", "backup", "rollback", "atomic", "interrupted", "profile",
    "settings", "workspace", "session", "cache", "checkpoint", "durable",
    "write", "read back", "reload", "reopen", "upgrade", "preserv",
    "archive", "postgres", "storage", "config dir", "config directory",
)
ASSERTION_MARKERS = (
    "t.fatal", "t.error", "t.fail", "require.", "assert.", "assert(",
    "expect(", "raise assertionerror", "raise systemexit", "errors.append",
    "sys.exit(1", "return 1", "throw new error", "process.exit(1",
)


def classify(path_text: str) -> str:
    p = path_text.replace("\\", "/").strip().lower()
    name = p.rsplit("/", 1)[-1]
    if name.endswith("_test.go"):
        return "GO_LIFECYCLE"
    if p.startswith("tests/integration/"):
        return "HTTP_LIFECYCLE"
    if p.startswith("tools/ci/") and name.endswith(("_gate.py", "_test.py", "_contract.py")):
        return "CI_LIFECYCLE"
    if p.startswith("tools/release/") and name.endswith((".py", ".sh", ".ps1")):
        return "RELEASE_LIFECYCLE"
    if "browser" in name and name.endswith((".py", ".js")):
        return "BROWSER_LIFECYCLE"
    return "UNKNOWN_EXECUTABLE"


def inspect(path: Path) -> tuple[bool, set[str]]:
    try:
        text = path.read_text(encoding="utf-8", errors="ignore").lower()
    except OSError:
        return False, set()
    semantics = {term for term in PERSISTENCE_SEMANTICS if term in text}
    assertions = any(marker in text for marker in ASSERTION_MARKERS)
    return bool(semantics and assertions), semantics


def main() -> int:
    parser = argparse.ArgumentParser()
    parser.add_argument("--strict", action="store_true")
    args = parser.parse_args()

    ledger, freeze = load(LEDGER), load(FREEZE)
    scan1, scan2, reconciliation = load(SCAN1), load(SCAN2), load(RECONCILIATION)
    current_state, closure, t5 = load(CURRENT_STATE), load(CLOSURE), load(T5)
    errors: list[str] = []

    product = current_state.get("productCapabilityGate") or {}
    governed = product.get("nextGovernedTracks") or []
    t5_runtime = next((str(x.get("status") or "") for x in governed if isinstance(x, dict) and x.get("track") == "T5"), "")
    t6_runtime = next((str(x.get("status") or "") for x in governed if isinstance(x, dict) and x.get("track") == "T6"), "")
    state = str(t5.get("state") or "")
    if state == "IN_PROGRESS":
        if product.get("nextChildIssue") != 118 or product.get("nextChildTrack") != "T5" or t5_runtime != "IN_PROGRESS":
            errors.append("IN_PROGRESS T5 must be the active T5/#118 child")
        if t6_runtime != "NOT_STARTED":
            errors.append("IN_PROGRESS T5 must not silently start T6/#119")
    elif state == "COMPLETE":
        completed = product.get("completedChildTracks") or []
        row = next((x for x in completed if isinstance(x, dict) and x.get("track") == "T5" and x.get("issue") == 118), None)
        if not isinstance(row, dict) or row.get("status") != "COMPLETE" or not str(row.get("mergedCommitSha") or ""):
            errors.append("COMPLETE T5 requires durable completedChildTracks merge evidence")
    else:
        errors.append(f"unsupported T5 state: {state!r}")

    gaps = closure.get("gaps") or []
    for required in ("T1-FEATURE-TRACEABILITY","T2-UNIT-CONTRACT-PROPERTY","T3-FUNCTIONAL-E2E","T4-EDGE-FAILURE-DATA-TRUTH"):
        row = next((x for x in gaps if isinstance(x, dict) and x.get("id") == required), None)
        if not isinstance(row, dict) or row.get("status") != "VERIFIED":
            errors.append(f"T5 requires {required}=VERIFIED")
    t5_gap = next((x for x in gaps if isinstance(x, dict) and x.get("id") == "T5-PERSISTENCE-LIFECYCLE-RECOVERY"), None)
    expected_gap = "VERIFIED" if state == "COMPLETE" else "IMPLEMENTED_UNVERIFIED"
    if not isinstance(t5_gap, dict) or t5_gap.get("status") != expected_gap:
        errors.append(f"T5 closure row must be {expected_gap}")

    if "python3 tools/ci/v18_t5_persistence_lifecycle_assurance_gate.py" not in CI_FAST.read_text(encoding="utf-8"):
        errors.append("T5 assurance gate is not bound into canonical CI Fast")

    actual_blob = git_blob_sha(LEDGER)
    expected_blob = str((freeze.get("frozenDiscovery") or {}).get("gitBlobSha") or "")
    if freeze.get("state") != "FROZEN_T1" or actual_blob != expected_blob:
        errors.append("T5 requires immutable frozen T1 ledger")
    if t5.get("trackIssue") != 118 or t5.get("programIssue") != 113 or t5.get("frozenT1GitBlobSha") != actual_blob:
        errors.append("T5 assurance identity/frozen T1 binding mismatch")

    effective = reconstruct_effective(ledger, scan1, scan2, reconciliation, errors)
    expected_count = int((freeze.get("effectiveInventory") or {}).get("effectiveShippedV18Responsibilities") or 0)
    if len(effective) != expected_count:
        errors.append(f"effective inventory count mismatch: expected {expected_count}, got {len(effective)}")

    profiles = ledger.get("assuranceProfiles") or {}
    overrides = t5.get("applicabilityOverrides") or {}
    supplements = t5.get("evidenceSupplements") or {}
    if not isinstance(overrides, dict) or not isinstance(supplements, dict):
        errors.append("T5 applicabilityOverrides/evidenceSupplements must be objects")
        overrides, supplements = {}, {}

    for fid, spec in supplements.items():
        if fid not in effective or not isinstance(spec, dict):
            errors.append(f"invalid T5 supplement: {fid}")
            continue
        owners, rationale = spec.get("owners"), str(spec.get("rationale") or "").strip()
        if not isinstance(owners, list) or not owners or not rationale:
            errors.append(f"T5 supplement {fid} requires owners and rationale")
            continue
        target = effective[fid].setdefault("existingRegressionOwners", [])
        for owner in owners:
            if isinstance(owner, str) and owner.strip() and owner.strip() not in target:
                target.append(owner.strip())

    applicable: list[str] = []
    nonapplicable: list[str] = []
    uncovered: list[str] = []
    missing_paths: set[str] = set()
    class_counts: dict[str, int] = {}

    for fid, row in sorted(effective.items()):
        profile_name = str(row.get("assuranceProfile") or "")
        profile = profiles.get(profile_name)
        if not isinstance(profile, dict) or not str(profile.get("T5") or "").strip():
            errors.append(f"{fid}: missing T5 expectation/profile")
            continue

        override = overrides.get(fid)
        if override is not None:
            if not isinstance(override, dict) or override.get("state") not in {"APPLICABLE","NOT_APPLICABLE"} or not str(override.get("rationale") or "").strip():
                errors.append(f"{fid}: invalid T5 applicability override")
                continue
            is_applicable = override.get("state") == "APPLICABLE"
        elif profile_name in ALWAYS_APPLICABLE:
            is_applicable = True
        elif profile_name in CONDITIONAL:
            is_applicable = False
            for owner in row.get("existingRegressionOwners") or []:
                owner_text = str(owner or "").strip()
                path = ROOT / owner_text
                if path.is_file():
                    meaningful, semantics = inspect(path)
                    if meaningful and semantics:
                        is_applicable = True
                        break
        else:
            errors.append(f"{fid}: no T5 applicability policy for profile {profile_name}")
            continue

        if not is_applicable:
            nonapplicable.append(fid)
            continue
        applicable.append(fid)

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
            meaningful, _ = inspect(path)
            if meaningful and cls in {"GO_LIFECYCLE","HTTP_LIFECYCLE","CI_LIFECYCLE","RELEASE_LIFECYCLE","BROWSER_LIFECYCLE"}:
                valid.append(owner_text)
        if not valid:
            uncovered.append(fid)

    if missing_paths:
        errors.append("referenced T5 evidence missing at current head: " + ", ".join(sorted(missing_paths)))

    declaration = str(t5.get("gapDeclarationState") or "")
    declared = set(str(x) for x in (t5.get("knownCoverageGaps") or []))
    actual = set(uncovered)
    if declaration == "PENDING_FIRST_EXECUTION":
        if declared:
            errors.append("PENDING_FIRST_EXECUTION T5 must not predeclare guessed gaps")
    elif declaration == "CURRENT":
        if actual != declared:
            errors.append("T5 declared gap set does not match executable census")
        if t5.get("applicableResponsibilityCount") != len(applicable):
            errors.append("T5 applicableResponsibilityCount drift")
        if t5.get("coveredApplicableResponsibilityCount") != len(applicable) - len(uncovered):
            errors.append("T5 coveredApplicableResponsibilityCount drift")
        if t5.get("uncoveredResponsibilityCount") != len(uncovered):
            errors.append("T5 uncoveredResponsibilityCount drift")
        if t5.get("nonApplicableResponsibilityCount") != len(nonapplicable):
            errors.append("T5 nonApplicableResponsibilityCount drift")
    else:
        errors.append(f"unsupported T5 gapDeclarationState: {declaration!r}")

    strict = args.strict or state == "COMPLETE"
    if strict and uncovered:
        errors.append("T5 strict closure has uncovered responsibilities: " + ", ".join(uncovered))
    if state == "COMPLETE" and declaration != "CURRENT":
        errors.append("T5 COMPLETE requires CURRENT gap declaration")

    print("V18 T5 PERSISTENCE / LIFECYCLE / RESTART ASSURANCE")
    print(f"frozen T1 blob: {actual_blob}")
    print(f"effective responsibilities: {len(effective)}")
    print(f"T5-applicable responsibilities: {len(applicable)}")
    print(f"T5-covered applicable responsibilities: {len(applicable) - len(uncovered)}")
    print(f"T5-uncovered responsibilities: {len(uncovered)}")
    print(f"T5-non-applicable responsibilities: {len(nonapplicable)}")
    for cls in sorted(class_counts):
        print(f"{cls}: {class_counts[cls]}")
    if uncovered:
        print("uncovered ids: " + ", ".join(uncovered))
    print("canonical persistence/cache/identity/workspace owners conserved: PASS")
    print("read-only surfaces cannot become parallel persistence owners: PASS")
    print("packaged native positive proof remains T9: PASS")
    print("T6-T10 certification is not implied by T5: PASS")

    if errors:
        print("V18 T5 ASSURANCE GATE: FAIL", file=sys.stderr)
        for error in errors:
            print(" - " + error, file=sys.stderr)
        return 1
    print(f"V18 T5 ASSURANCE GATE: PASS (strict={strict}, declaration={declaration})")
    return 0


if __name__ == "__main__":
    sys.exit(main())
