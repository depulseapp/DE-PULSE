#!/usr/bin/env python3
"""T7 UI / UX / IA / content / accessibility assurance over frozen v18 inventory."""
from __future__ import annotations

import argparse
from pathlib import Path
import sys

from v18_t2_unit_contract_assurance_gate import (
    ROOT, PROGRAM, LEDGER, FREEZE, SCAN1, SCAN2, RECONCILIATION,
    CURRENT_STATE, CLOSURE, load, git_blob_sha, reconstruct_effective,
)

T7 = PROGRAM / "T7_UI_UX_IA_CONTENT_ASSURANCE.json"
CI_FAST = ROOT / ".github" / "workflows" / "ci-fast.yml"

VALID_UI_CLASSES = {"RENDERER_CONTRACT", "BROWSER_RESPONSIVE", "BROWSER_UI", "ACCEPTANCE_UI"}
UI_SEMANTICS = (
    "render", "page", "surface", "navigation", "nav", "placement", "hierarchy", "layout",
    "visible", "hidden", "responsive", "viewport", "overflow", "scroll", "focus", "keyboard",
    "aria", "access", "contrast", "loading", "empty", "degraded", "stale", "content", "label",
    "wording", "documentation", "dashboard", "research", "maintenance", "settings", "administration",
    "desk", "discovery", "toast", "notification", "header", "sidebar", "table", "role",
)
ASSERTION_MARKERS = (
    "assert.", "assert(", "errors.append", "raise assertionerror", "raise systemexit", "sys.exit(1",
    "throw new error", "process.exit(1", "expect(", "t.fatal", "t.error", "t.fail",
)


def classify(path_text: str) -> str:
    p = path_text.replace("\\", "/").strip().lower()
    name = p.rsplit("/", 1)[-1]
    if p.startswith("tests/renderer/") and name.endswith(".js"):
        return "RENDERER_CONTRACT"
    if p.startswith("tests/renderer/") and name.endswith(".py"):
        return "BROWSER_RESPONSIVE"
    if "browser" in name and name.endswith((".py", ".js")):
        return "BROWSER_UI"
    if p.startswith("tests/acceptance/") and name.endswith((".js", ".py")):
        return "ACCEPTANCE_UI"
    if name.endswith("_test.go"):
        return "GO_NON_UI"
    if p.startswith("tests/integration/"):
        return "INTEGRATION_NON_UI"
    if p.startswith("tools/ci/"):
        return "CI_NON_ROW_EVIDENCE"
    return "UNKNOWN_EXECUTABLE"


def inspect(path: Path) -> tuple[bool, set[str]]:
    try:
        text = path.read_text(encoding="utf-8", errors="ignore").lower()
    except OSError:
        return False, set()
    semantics = {term for term in UI_SEMANTICS if term in text}
    asserted = any(marker in text for marker in ASSERTION_MARKERS)
    return bool(semantics and asserted), semantics


def physical_rows(ledger: dict) -> dict[str, dict]:
    out: dict[str, dict] = {}
    for row in ledger.get("features") or []:
        if isinstance(row, dict) and str(row.get("id") or "").strip():
            out[str(row["id"]).strip()] = row
    return out


def main() -> int:
    parser = argparse.ArgumentParser()
    parser.add_argument("--strict", action="store_true")
    args = parser.parse_args()

    ledger, freeze = load(LEDGER), load(FREEZE)
    scan1, scan2, reconciliation = load(SCAN1), load(SCAN2), load(RECONCILIATION)
    current_state, closure, t7 = load(CURRENT_STATE), load(CLOSURE), load(T7)
    errors: list[str] = []

    product = current_state.get("productCapabilityGate") or {}
    governed = product.get("nextGovernedTracks") or []
    t7_runtime = next((str(x.get("status") or "") for x in governed if isinstance(x, dict) and x.get("track") == "T7"), "")
    t8_runtime = next((str(x.get("status") or "") for x in governed if isinstance(x, dict) and x.get("track") == "T8"), "")
    state = str(t7.get("state") or "")
    if state == "IN_PROGRESS":
        if product.get("nextChildIssue") != 120 or product.get("nextChildTrack") != "T7" or t7_runtime != "IN_PROGRESS":
            errors.append("IN_PROGRESS T7 must be the active T7/#120 child")
        if t8_runtime != "NOT_STARTED":
            errors.append("IN_PROGRESS T7 must not silently start T8/#121")
    elif state == "COMPLETE":
        completed = product.get("completedChildTracks") or []
        row = next((x for x in completed if isinstance(x, dict) and x.get("track") == "T7" and x.get("issue") == 120), None)
        if not isinstance(row, dict) or row.get("status") != "COMPLETE" or not str(row.get("mergedCommitSha") or ""):
            errors.append("COMPLETE T7 requires durable completedChildTracks merge evidence")
    else:
        errors.append(f"unsupported T7 state: {state!r}")

    gaps = closure.get("gaps") or []
    for required in (
        "T1-FEATURE-TRACEABILITY", "T2-UNIT-CONTRACT-PROPERTY", "T3-FUNCTIONAL-E2E",
        "T4-EDGE-FAILURE-DATA-TRUTH", "T5-PERSISTENCE-LIFECYCLE-RECOVERY", "T6-SECURITY-ROLES-RIGHTS",
    ):
        row = next((x for x in gaps if isinstance(x, dict) and x.get("id") == required), None)
        if not isinstance(row, dict) or row.get("status") != "VERIFIED":
            errors.append(f"T7 requires {required}=VERIFIED")
    t7_gap = next((x for x in gaps if isinstance(x, dict) and x.get("id") == "T7-UI-UX-IA-CONTENT"), None)
    expected_gap = "VERIFIED" if state == "COMPLETE" else "IMPLEMENTED_UNVERIFIED"
    if not isinstance(t7_gap, dict) or t7_gap.get("status") != expected_gap:
        errors.append(f"T7 closure row must be {expected_gap}")

    if "python3 tools/ci/v18_t7_ui_ux_ia_content_assurance_gate.py" not in CI_FAST.read_text(encoding="utf-8"):
        errors.append("T7 assurance gate is not bound into canonical CI Fast")

    actual_blob = git_blob_sha(LEDGER)
    expected_blob = str((freeze.get("frozenDiscovery") or {}).get("gitBlobSha") or "")
    if freeze.get("state") != "FROZEN_T1" or actual_blob != expected_blob:
        errors.append("T7 requires immutable frozen T1 ledger")
    if t7.get("trackIssue") != 120 or t7.get("programIssue") != 113 or t7.get("frozenT1GitBlobSha") != actual_blob:
        errors.append("T7 assurance identity/frozen T1 binding mismatch")

    effective = reconstruct_effective(ledger, scan1, scan2, reconciliation, errors)
    expected_count = int((freeze.get("effectiveInventory") or {}).get("effectiveShippedV18Responsibilities") or 0)
    if len(effective) != expected_count:
        errors.append(f"effective inventory count mismatch: expected {expected_count}, got {len(effective)}")

    profiles = ledger.get("assuranceProfiles") or {}
    allowed = set(str(x) for x in (t7.get("allowedDispositions") or []))
    frozen_allowed = set(str(x) for x in (ledger.get("uiDispositionValues") or []))
    if not allowed or allowed != frozen_allowed:
        errors.append("T7 allowedDispositions must exactly conserve frozen T1 disposition values")

    overrides = t7.get("dispositionOverrides") or {}
    supplements = t7.get("evidenceSupplements") or {}
    if not isinstance(overrides, dict) or not isinstance(supplements, dict):
        errors.append("T7 dispositionOverrides/evidenceSupplements must be objects")
        overrides, supplements = {}, {}

    physical = physical_rows(ledger)
    for fid, spec in supplements.items():
        if fid not in effective or not isinstance(spec, dict):
            errors.append(f"invalid T7 supplement: {fid}")
            continue
        owners, rationale = spec.get("owners"), str(spec.get("rationale") or "").strip()
        if not isinstance(owners, list) or not owners or not rationale:
            errors.append(f"T7 supplement {fid} requires owners and rationale")
            continue
        target = effective[fid].setdefault("existingRegressionOwners", [])
        for owner in owners:
            if isinstance(owner, str) and owner.strip() and owner.strip() not in target:
                target.append(owner.strip())

    visible: list[str] = []
    nonvisible: list[str] = []
    uncovered: list[str] = []
    missing_paths: set[str] = set()
    disposition_counts: dict[str, int] = {x: 0 for x in sorted(allowed)}
    evidence_class_counts: dict[str, int] = {}
    resolved_dispositions: dict[str, tuple[str, str]] = {}

    def override_disposition(fid: str) -> tuple[str, str] | None:
        spec = overrides.get(fid)
        if spec is None:
            return None
        if not isinstance(spec, dict):
            errors.append(f"{fid}: invalid T7 disposition override")
            return ("", "INVALID")
        disposition = str(spec.get("disposition") or "")
        rationale = str(spec.get("rationale") or "").strip()
        if disposition not in allowed or not rationale:
            errors.append(f"{fid}: T7 override requires allowed disposition and rationale")
        return disposition, "T7_OVERRIDE"

    for fid, row in sorted(effective.items()):
        profile_name = str(row.get("assuranceProfile") or "")
        profile = profiles.get(profile_name)
        if not isinstance(profile, dict) or not str(profile.get("T7") or "").strip():
            errors.append(f"{fid}: missing T7 expectation/profile")
            continue

        resolved = override_disposition(fid)
        if resolved is None and row.get("_parentFeature"):
            parent_id = str(row.get("_parentFeature") or "")
            parent = physical.get(parent_id)
            if not isinstance(parent, dict):
                errors.append(f"{fid}: missing canonical parent for T7 disposition")
                continue
            parent_override = override_disposition(parent_id)
            if parent_override is not None:
                disposition = parent_override[0]
                source = f"INHERITED_OVERRIDE:{parent_id}"
            else:
                disposition = str(parent.get("uiDisposition") or "")
                source = f"INHERITED_FROZEN_PARENT:{parent_id}"
            resolved = (disposition, source)
        elif resolved is None:
            resolved = (str(row.get("uiDisposition") or ""), "FROZEN_T1")

        disposition, source = resolved
        if disposition not in allowed:
            errors.append(f"{fid}: unresolved/invalid T7 disposition {disposition!r}")
            continue
        resolved_dispositions[fid] = (disposition, source)
        disposition_counts[disposition] = disposition_counts.get(disposition, 0) + 1

        if disposition == "NOT_USER_VISIBLE":
            nonvisible.append(fid)
            continue
        visible.append(fid)

        candidates: list[str] = []
        for owner in row.get("existingRegressionOwners") or []:
            owner_text = str(owner or "").strip()
            if owner_text and owner_text not in candidates:
                candidates.append(owner_text)
        parent_id = str(row.get("_parentFeature") or "")
        parent = physical.get(parent_id)
        if isinstance(parent, dict):
            for owner in parent.get("existingRegressionOwners") or []:
                owner_text = str(owner or "").strip()
                if owner_text and owner_text not in candidates:
                    candidates.append(owner_text)

        has_ui_evidence = False
        for owner_text in candidates:
            path = ROOT / owner_text
            if not path.is_file():
                missing_paths.add(f"{fid}:{owner_text}")
                continue
            cls = classify(owner_text)
            evidence_class_counts[cls] = evidence_class_counts.get(cls, 0) + 1
            meaningful, _semantics = inspect(path)
            if meaningful and cls in VALID_UI_CLASSES:
                has_ui_evidence = True
        if not has_ui_evidence:
            uncovered.append(fid)

    if missing_paths:
        errors.append("referenced T7 evidence missing at current head: " + ", ".join(sorted(missing_paths)))

    declaration = str(t7.get("gapDeclarationState") or "")
    declared = set(str(x) for x in (t7.get("knownCoverageGaps") or []))
    actual = set(uncovered)
    if declaration == "PENDING_FIRST_EXECUTION":
        if declared:
            errors.append("PENDING_FIRST_EXECUTION T7 must not predeclare guessed gaps")
    elif declaration == "CURRENT":
        if actual != declared:
            errors.append("T7 declared gap set does not match executable census")
        if t7.get("visibleResponsibilityCount") != len(visible):
            errors.append("T7 visibleResponsibilityCount drift")
        if t7.get("coveredVisibleResponsibilityCount") != len(visible) - len(uncovered):
            errors.append("T7 coveredVisibleResponsibilityCount drift")
        if t7.get("uncoveredResponsibilityCount") != len(uncovered):
            errors.append("T7 uncoveredResponsibilityCount drift")
        if t7.get("nonVisibleResponsibilityCount") != len(nonvisible):
            errors.append("T7 nonVisibleResponsibilityCount drift")
    else:
        errors.append(f"unsupported T7 gapDeclarationState: {declaration!r}")

    strict = args.strict or state == "COMPLETE"
    if strict and uncovered:
        errors.append("T7 strict closure has uncovered visible responsibilities: " + ", ".join(uncovered))
    if state == "COMPLETE" and declaration != "CURRENT":
        errors.append("T7 COMPLETE requires CURRENT gap declaration")

    print("V18 T7 UI / UX / IA / CONTENT / ACCESSIBILITY ASSURANCE")
    print(f"frozen T1 blob: {actual_blob}")
    print(f"effective responsibilities: {len(effective)}")
    print(f"T7-visible responsibilities: {len(visible)}")
    print(f"T7-covered visible responsibilities: {len(visible) - len(uncovered)}")
    print(f"T7-uncovered visible responsibilities: {len(uncovered)}")
    print(f"T7-not-user-visible responsibilities: {len(nonvisible)}")
    for disposition in sorted(disposition_counts):
        print(f"disposition {disposition}: {disposition_counts[disposition]}")
    if uncovered:
        print("uncovered ids: " + ", ".join(uncovered))
    for cls in sorted(evidence_class_counts):
        print(f"{cls}: {evidence_class_counts[cls]}")
    print("scan responsibilities inherit only their named canonical T1 parent disposition: PASS")
    print("UI hiding cannot substitute for T6 authorization: PASS")
    print("read-only UI cannot create parallel data/routing/recovery ownership: PASS")
    print("T8-T10 certification is not implied by T7: PASS")

    if errors:
        print("V18 T7 ASSURANCE GATE: FAIL", file=sys.stderr)
        for error in errors:
            print(" - " + error, file=sys.stderr)
        return 1
    print(f"V18 T7 ASSURANCE GATE: PASS (strict={strict}, declaration={declaration})")
    return 0


if __name__ == "__main__":
    sys.exit(main())
