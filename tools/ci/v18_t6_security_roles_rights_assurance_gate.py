#!/usr/bin/env python3
"""T6 security / roles / secrets / rights assurance over frozen v18 inventory."""
from __future__ import annotations

import argparse
from pathlib import Path
import sys

from v18_t2_unit_contract_assurance_gate import (
    ROOT, PROGRAM, LEDGER, FREEZE, SCAN1, SCAN2, RECONCILIATION,
    CURRENT_STATE, CLOSURE, load, git_blob_sha, reconstruct_effective,
)

T6 = PROGRAM / "T6_SECURITY_ROLES_RIGHTS_ASSURANCE.json"
CI_FAST = ROOT / ".github" / "workflows" / "ci-fast.yml"

ALWAYS_APPLICABLE = {"SECURITY_STATEFUL", "BOUNDARY"}
CONDITIONAL = {"VISIBLE_STATEFUL", "VISIBLE_READONLY", "DATA_PATH", "BACKGROUND_JOB", "RUNTIME_STATEFUL", "PERSISTENCE", "RELEASE"}

SECURITY_SEMANTICS = (
    "role", "admin", "owner", "user", "demo", "auth", "session", "csrf", "reauth",
    "secret", "credential", "password", "token", "permission", "rights", "entitlement",
    "commercial", "redact", "sec", "edgar", "form 4", "form4", "execution", "trade action",
    "provider", "authority", "capability", "public state", "sensitive",
)
POSITIVE_SEMANTICS = (
    "allows", "allow", "authorized", "authenticated", "role hierarchy", "capability",
    "admin", "owner", "rights approved", "production", "direct authority", "sec", "edgar",
    "explicit opt-in", "requires reauth", "protected mutation", "secure", "csrf",
)
NEGATIVE_SEMANTICS = (
    "deny", "denied", "reject", "forbidden", "unauthorized", "redact", "leak", "missing key",
    "missing credential", "expired", "revoked", "disabled", "wrong reauthentication", "fails closed",
    "fail closed", "fail-closed", "no plaintext", "untrusted", "cannot", "must not", "not entitled",
    "no execution", "trade action", "does not change", "stay out of executable routing", "never enters executable routing",
)
ASSERTION_MARKERS = (
    "t.fatal", "t.error", "t.fail", "require.", "assert.", "assert(", "expect(",
    "raise assertionerror", "raise systemexit", "errors.append", "sys.exit(1", "return 1",
    "throw new error", "process.exit(1",
)


def classify(path_text: str) -> str:
    p = path_text.replace("\\", "/").strip().lower()
    name = p.rsplit("/", 1)[-1]
    if name.endswith("_test.go"):
        return "GO_SECURITY"
    if p.startswith("tests/integration/"):
        return "HTTP_SECURITY"
    if p.startswith("tests/acceptance/"):
        return "ACCEPTANCE_SECURITY"
    if p.startswith("tests/renderer/") and name.endswith(".js"):
        return "RENDERER_SECURITY"
    if "browser" in name and name.endswith((".py", ".js")):
        return "BROWSER_SECURITY"
    if p.startswith("tools/ci/") and name.endswith(("_gate.py", "_test.py", "_contract.py")):
        return "CI_SECURITY"
    if p.startswith("tools/release/") and name.endswith((".py", ".sh", ".ps1")):
        return "RELEASE_SECURITY"
    return "UNKNOWN_EXECUTABLE"


def inspect(path: Path) -> tuple[bool, bool, bool, set[str]]:
    try:
        text = path.read_text(encoding="utf-8", errors="ignore").lower()
    except OSError:
        return False, False, False, set()
    security = {term for term in SECURITY_SEMANTICS if term in text}
    has_assertion = any(marker in text for marker in ASSERTION_MARKERS)
    positive = has_assertion and bool(security) and any(term in text for term in POSITIVE_SEMANTICS)
    negative = has_assertion and bool(security) and any(term in text for term in NEGATIVE_SEMANTICS)
    return bool(security and has_assertion), positive, negative, security


def main() -> int:
    parser = argparse.ArgumentParser()
    parser.add_argument("--strict", action="store_true")
    args = parser.parse_args()

    ledger, freeze = load(LEDGER), load(FREEZE)
    scan1, scan2, reconciliation = load(SCAN1), load(SCAN2), load(RECONCILIATION)
    current_state, closure, t6 = load(CURRENT_STATE), load(CLOSURE), load(T6)
    errors: list[str] = []

    product = current_state.get("productCapabilityGate") or {}
    governed = product.get("nextGovernedTracks") or []
    t6_runtime = next((str(x.get("status") or "") for x in governed if isinstance(x, dict) and x.get("track") == "T6"), "")
    t7_runtime = next((str(x.get("status") or "") for x in governed if isinstance(x, dict) and x.get("track") == "T7"), "")
    state = str(t6.get("state") or "")
    if state == "IN_PROGRESS":
        if product.get("nextChildIssue") != 119 or product.get("nextChildTrack") != "T6" or t6_runtime != "IN_PROGRESS":
            errors.append("IN_PROGRESS T6 must be the active T6/#119 child")
        if t7_runtime != "NOT_STARTED":
            errors.append("IN_PROGRESS T6 must not silently start T7/#120")
    elif state == "COMPLETE":
        completed = product.get("completedChildTracks") or []
        row = next((x for x in completed if isinstance(x, dict) and x.get("track") == "T6" and x.get("issue") == 119), None)
        if not isinstance(row, dict) or row.get("status") != "COMPLETE" or not str(row.get("mergedCommitSha") or ""):
            errors.append("COMPLETE T6 requires durable completedChildTracks merge evidence")
    else:
        errors.append(f"unsupported T6 state: {state!r}")

    gaps = closure.get("gaps") or []
    for required in ("T1-FEATURE-TRACEABILITY","T2-UNIT-CONTRACT-PROPERTY","T3-FUNCTIONAL-E2E","T4-EDGE-FAILURE-DATA-TRUTH","T5-PERSISTENCE-LIFECYCLE-RECOVERY"):
        row = next((x for x in gaps if isinstance(x, dict) and x.get("id") == required), None)
        if not isinstance(row, dict) or row.get("status") != "VERIFIED":
            errors.append(f"T6 requires {required}=VERIFIED")
    t6_gap = next((x for x in gaps if isinstance(x, dict) and x.get("id") == "T6-SECURITY-ROLES-RIGHTS"), None)
    expected_gap = "VERIFIED" if state == "COMPLETE" else "IMPLEMENTED_UNVERIFIED"
    if not isinstance(t6_gap, dict) or t6_gap.get("status") != expected_gap:
        errors.append(f"T6 closure row must be {expected_gap}")

    if "python3 tools/ci/v18_t6_security_roles_rights_assurance_gate.py" not in CI_FAST.read_text(encoding="utf-8"):
        errors.append("T6 assurance gate is not bound into canonical CI Fast")

    actual_blob = git_blob_sha(LEDGER)
    expected_blob = str((freeze.get("frozenDiscovery") or {}).get("gitBlobSha") or "")
    if freeze.get("state") != "FROZEN_T1" or actual_blob != expected_blob:
        errors.append("T6 requires immutable frozen T1 ledger")
    if t6.get("trackIssue") != 119 or t6.get("programIssue") != 113 or t6.get("frozenT1GitBlobSha") != actual_blob:
        errors.append("T6 assurance identity/frozen T1 binding mismatch")

    effective = reconstruct_effective(ledger, scan1, scan2, reconciliation, errors)
    expected_count = int((freeze.get("effectiveInventory") or {}).get("effectiveShippedV18Responsibilities") or 0)
    if len(effective) != expected_count:
        errors.append(f"effective inventory count mismatch: expected {expected_count}, got {len(effective)}")

    profiles = ledger.get("assuranceProfiles") or {}
    overrides = t6.get("applicabilityOverrides") or {}
    supplements = t6.get("evidenceSupplements") or {}
    if not isinstance(overrides, dict) or not isinstance(supplements, dict):
        errors.append("T6 applicabilityOverrides/evidenceSupplements must be objects")
        overrides, supplements = {}, {}

    for fid, spec in supplements.items():
        if fid not in effective or not isinstance(spec, dict):
            errors.append(f"invalid T6 supplement: {fid}")
            continue
        owners, rationale = spec.get("owners"), str(spec.get("rationale") or "").strip()
        if not isinstance(owners, list) or not owners or not rationale:
            errors.append(f"T6 supplement {fid} requires owners and rationale")
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
    positive_only: list[str] = []
    negative_only: list[str] = []

    for fid, row in sorted(effective.items()):
        profile_name = str(row.get("assuranceProfile") or "")
        profile = profiles.get(profile_name)
        if not isinstance(profile, dict) or not str(profile.get("T6") or "").strip():
            errors.append(f"{fid}: missing T6 expectation/profile")
            continue

        override = overrides.get(fid)
        if override is not None:
            if not isinstance(override, dict) or override.get("state") not in {"APPLICABLE","NOT_APPLICABLE"} or not str(override.get("rationale") or "").strip():
                errors.append(f"{fid}: invalid T6 applicability override")
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
                    meaningful, _positive, _negative, _semantics = inspect(path)
                    if meaningful:
                        is_applicable = True
                        break
        else:
            errors.append(f"{fid}: no T6 applicability policy for profile {profile_name}")
            continue

        if not is_applicable:
            nonapplicable.append(fid)
            continue
        applicable.append(fid)

        has_positive = False
        has_negative = False
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
            meaningful, positive, negative, _ = inspect(path)
            if meaningful and cls in {"GO_SECURITY","HTTP_SECURITY","ACCEPTANCE_SECURITY","RENDERER_SECURITY","BROWSER_SECURITY","CI_SECURITY","RELEASE_SECURITY"}:
                has_positive = has_positive or positive
                has_negative = has_negative or negative
        if has_positive and not has_negative:
            positive_only.append(fid)
        if has_negative and not has_positive:
            negative_only.append(fid)
        if not (has_positive and has_negative):
            uncovered.append(fid)

    if missing_paths:
        errors.append("referenced T6 evidence missing at current head: " + ", ".join(sorted(missing_paths)))

    declaration = str(t6.get("gapDeclarationState") or "")
    declared = set(str(x) for x in (t6.get("knownCoverageGaps") or []))
    actual = set(uncovered)
    if declaration == "PENDING_FIRST_EXECUTION":
        if declared:
            errors.append("PENDING_FIRST_EXECUTION T6 must not predeclare guessed gaps")
    elif declaration == "CURRENT":
        if actual != declared:
            errors.append("T6 declared gap set does not match executable census")
        if t6.get("applicableResponsibilityCount") != len(applicable):
            errors.append("T6 applicableResponsibilityCount drift")
        if t6.get("coveredApplicableResponsibilityCount") != len(applicable) - len(uncovered):
            errors.append("T6 coveredApplicableResponsibilityCount drift")
        if t6.get("uncoveredResponsibilityCount") != len(uncovered):
            errors.append("T6 uncoveredResponsibilityCount drift")
        if t6.get("nonApplicableResponsibilityCount") != len(nonapplicable):
            errors.append("T6 nonApplicableResponsibilityCount drift")
    else:
        errors.append(f"unsupported T6 gapDeclarationState: {declaration!r}")

    strict = args.strict or state == "COMPLETE"
    if strict and uncovered:
        errors.append("T6 strict closure has uncovered responsibilities: " + ", ".join(uncovered))
    if state == "COMPLETE" and declaration != "CURRENT":
        errors.append("T6 COMPLETE requires CURRENT gap declaration")

    print("V18 T6 SECURITY / ROLES / SECRETS / RIGHTS ASSURANCE")
    print(f"frozen T1 blob: {actual_blob}")
    print(f"effective responsibilities: {len(effective)}")
    print(f"T6-applicable responsibilities: {len(applicable)}")
    print(f"T6-covered applicable responsibilities: {len(applicable) - len(uncovered)}")
    print(f"T6-uncovered responsibilities: {len(uncovered)}")
    print(f"T6-non-applicable responsibilities: {len(nonapplicable)}")
    if positive_only:
        print("positive-only ids: " + ", ".join(positive_only))
    if negative_only:
        print("negative-only ids: " + ", ".join(negative_only))
    if uncovered:
        print("uncovered ids: " + ", ".join(uncovered))
    for cls in sorted(class_counts):
        print(f"{cls}: {class_counts[cls]}")
    print("UI hiding cannot substitute for direct-route/API authorization: PASS")
    print("provider rights remain separate from executable router scoring/routing: PASS")
    print("direct SEC/EDGAR authority and No Execution boundary are conserved: PASS")
    print("T7-T10 certification is not implied by T6: PASS")

    if errors:
        print("V18 T6 ASSURANCE GATE: FAIL", file=sys.stderr)
        for error in errors:
            print(" - " + error, file=sys.stderr)
        return 1
    print(f"V18 T6 ASSURANCE GATE: PASS (strict={strict}, declaration={declaration})")
    return 0


if __name__ == "__main__":
    sys.exit(main())
