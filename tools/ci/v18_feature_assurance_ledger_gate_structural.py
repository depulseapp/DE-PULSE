#!/usr/bin/env python3
"""Fail-closed T1 feature reality/quality audit gate for final v18 closure.

The on-disk representation is normalized on purpose:
- feature-assurance-ledger.json owns canonical high-level rows and T2-T9 profiles;
- two independent omission scans retain source/runtime-discovered responsibilities;
- T1_FINAL_RECONCILIATION.json maps every scan responsibility to exactly one
  canonical parent and explicitly excludes evidenced future-program source;
- T1_QUALITY_AUDIT.json supplies reusable profile assessments plus narrow
  responsibility overrides.

This gate resolves that storage form into an effective per-responsibility audit.
It never treats a closed issue, a test name, or inherited profile text as proof
that later T2-T9 behavioral assurance has passed.
"""

from __future__ import annotations

import argparse
import copy
import json
from pathlib import Path
import sys

ROOT = Path(__file__).resolve().parents[2]
PROGRAM = ROOT / "governance" / "programs" / "ADAPT-V18-FINAL-CLOSURE-10-10-001"
LEDGER = PROGRAM / "feature-assurance-ledger.json"
SCAN_1 = PROGRAM / "T1_INDEPENDENT_OMISSION_SCAN.json"
SCAN_2 = PROGRAM / "T1_INDEPENDENT_OMISSION_SCAN_2.json"
RECONCILIATION = PROGRAM / "T1_FINAL_RECONCILIATION.json"
QUALITY = PROGRAM / "T1_QUALITY_AUDIT.json"
CONTRACT = ROOT / "governance" / "feature-closure-audit-contract.json"

BASE_ROW_KEYS = {
    "id",
    "name",
    "category",
    "requirementProvenance",
    "canonicalSourceOwners",
    "consumers",
    "existingRegressionOwners",
    "positiveFunctionalEvidenceExpectation",
    "assuranceProfile",
    "durableRegressionOwner",
    "currentAssuranceState",
    "blockingStates",
    "uiDisposition",
}

QUALITY_KEYS = {
    "productUtilityAssessment",
    "correctnessAssessment",
    "architectureAssessment",
    "intelligenceAssessment",
    "dataTruthAssessment",
    "maintainabilityAssessment",
    "staleCodeAssessment",
    "sourceDisposition",
    "securityRightsAssessment",
    "persistenceLifecycleAssessment",
    "performanceEfficiencyAssessment",
    "observabilityAssessment",
    "improvementOpportunities",
    "closureDecision",
}

BLOCKING_DECISIONS = {
    "IMPROVE_BEFORE_CLOSURE",
    "REFACTOR_BEFORE_CLOSURE",
    "REWRITE_BEFORE_CLOSURE",
    "MERGE_OR_REMOVE_BEFORE_CLOSURE",
    "EXTERNAL_BLOCKED_WITH_EVIDENCE",
}


def load_json(path: Path) -> dict:
    try:
        value = json.loads(path.read_text(encoding="utf-8"))
    except Exception as exc:  # pragma: no cover - command-line diagnostics
        raise SystemExit(f"FAIL: cannot read {path.relative_to(ROOT)}: {exc}") from exc
    if not isinstance(value, dict):
        raise SystemExit(f"FAIL: {path.relative_to(ROOT)} must contain a JSON object")
    return value


def nonempty_string(value: object) -> bool:
    return isinstance(value, str) and bool(value.strip())


def nonempty_string_list(value: object) -> bool:
    return (
        isinstance(value, list)
        and bool(value)
        and all(nonempty_string(item) for item in value)
    )


def merge_dict(base: dict, override: dict) -> dict:
    """Deep-merge only dictionaries; lists/scalars are intentional replacements."""
    out = copy.deepcopy(base)
    for key, value in override.items():
        if isinstance(value, dict) and isinstance(out.get(key), dict):
            out[key] = merge_dict(out[key], value)
        else:
            out[key] = copy.deepcopy(value)
    return out


def check_repo_path(path_text: str) -> bool:
    if not nonempty_string(path_text):
        return False
    if path_text.startswith("#") or path_text.startswith("issue:") or path_text.startswith("PR#"):
        return True
    return (ROOT / path_text).exists()


def validate_base_row(row: dict, label: str, profiles: dict, allowed_states: set[str], allowed_dispositions: set[str], errors: list[str]) -> None:
    missing = sorted(BASE_ROW_KEYS - set(row))
    if missing:
        errors.append(f"{label} missing base keys: {', '.join(missing)}")
    for key in ("requirementProvenance", "canonicalSourceOwners", "consumers", "existingRegressionOwners", "durableRegressionOwner"):
        if not nonempty_string_list(row.get(key)):
            errors.append(f"{label} has empty/invalid {key}")
    if not nonempty_string(row.get("positiveFunctionalEvidenceExpectation")):
        errors.append(f"{label} has no positive functional evidence expectation")
    profile = str(row.get("assuranceProfile") or "").strip()
    if profile not in profiles:
        errors.append(f"{label} references unknown assuranceProfile {profile!r}")
    else:
        expected_tracks = {f"T{i}" for i in range(2, 10)}
        missing_tracks = sorted(expected_tracks - set(profiles[profile]))
        if missing_tracks:
            errors.append(f"assuranceProfile {profile} missing downstream expectations: {', '.join(missing_tracks)}")
    state = str(row.get("currentAssuranceState") or "").strip()
    if state not in allowed_states:
        errors.append(f"{label} has invalid assurance state {state!r}")
    disposition = str(row.get("uiDisposition") or "").strip()
    if disposition not in allowed_dispositions:
        errors.append(f"{label} has invalid uiDisposition {disposition!r}")
    if not isinstance(row.get("blockingStates"), list):
        errors.append(f"{label} blockingStates must be an array")

    for owner in row.get("canonicalSourceOwners") or []:
        if not check_repo_path(owner):
            errors.append(f"{label} canonical source owner does not exist at current head: {owner}")
    for test in row.get("existingRegressionOwners") or []:
        if not check_repo_path(test):
            errors.append(f"{label} regression owner does not exist at current head: {test}")


def resolve_scan_row(item: dict, scan_name: str, physical: dict[str, dict], reconciliation: dict, errors: list[str]) -> dict | None:
    item_id = str(item.get("id") or "").strip()
    category = str(item.get("category") or "").strip()
    if not item_id or not category:
        errors.append(f"{scan_name} contains omission without id/category")
        return None
    overrides = reconciliation.get("scanParentOverrides") or {}
    defaults = reconciliation.get("scanParentByCategory") or {}
    parent_id = str(overrides.get(item_id) or defaults.get(category) or "").strip()
    if not parent_id:
        errors.append(f"{scan_name}:{item_id} has no deterministic canonical parent for category {category}")
        return None
    parent = physical.get(parent_id)
    if parent is None:
        errors.append(f"{scan_name}:{item_id} parent {parent_id} is not a canonical high-level ledger row")
        return None

    owners = item.get("owners")
    if not nonempty_string_list(owners):
        errors.append(f"{scan_name}:{item_id} has no direct source-owner evidence")
        return None
    tests = item.get("tests")
    if not nonempty_string_list(tests):
        tests = parent.get("existingRegressionOwners")
    reason = str(item.get("reason") or "").strip()
    expectation = reason or (
        f"Responsibility {item_id} must preserve the parent contract of {parent_id} "
        "through the source owners discovered by the independent T1 scan."
    )
    provenance = list(parent.get("requirementProvenance") or []) + [f"t1-scan:{scan_name}"]

    return {
        "id": item_id,
        "name": item_id,
        "category": category,
        "requirementProvenance": provenance,
        "canonicalSourceOwners": list(owners),
        "consumers": list(parent.get("consumers") or []),
        "existingRegressionOwners": list(tests or []),
        "positiveFunctionalEvidenceExpectation": expectation,
        "assuranceProfile": parent.get("assuranceProfile"),
        "durableRegressionOwner": list(parent.get("durableRegressionOwner") or []),
        "currentAssuranceState": parent.get("currentAssuranceState"),
        "blockingStates": [],
        "uiDisposition": parent.get("uiDisposition"),
        "_parentFeature": parent_id,
        "_sourceScan": scan_name,
    }


def validate_quality(effective: dict, quality: dict, contract: dict, errors: list[str]) -> dict | None:
    feature_id = str(effective.get("id") or "")
    profile = str(effective.get("assuranceProfile") or "")
    defaults = quality.get("profileDefaults") or {}
    default = defaults.get(profile)
    if not isinstance(default, dict):
        errors.append(f"{feature_id} has no quality profile default for {profile!r}")
        return None
    overrides = quality.get("overrides") or {}
    override = overrides.get(feature_id) or {}
    if not isinstance(override, dict):
        errors.append(f"{feature_id} quality override must be an object")
        return None
    resolved = merge_dict(default, override)

    missing = sorted(QUALITY_KEYS - set(resolved))
    if missing:
        errors.append(f"{feature_id} effective quality audit missing: {', '.join(missing)}")
        return resolved
    for key in (
        "productUtilityAssessment",
        "correctnessAssessment",
        "architectureAssessment",
        "dataTruthAssessment",
        "maintainabilityAssessment",
        "staleCodeAssessment",
        "sourceDisposition",
        "securityRightsAssessment",
        "persistenceLifecycleAssessment",
        "performanceEfficiencyAssessment",
        "observabilityAssessment",
        "closureDecision",
    ):
        if not nonempty_string(resolved.get(key)):
            errors.append(f"{feature_id} quality field {key} is empty")
    if not isinstance(resolved.get("improvementOpportunities"), list):
        errors.append(f"{feature_id} improvementOpportunities must be an array")

    intelligence = resolved.get("intelligenceAssessment")
    if not isinstance(intelligence, dict):
        errors.append(f"{feature_id} intelligenceAssessment must be an object")
    else:
        expected_fields = {"currentLevel", "expectedLevel", "fit", "gaps", "rationale"}
        missing_intel = sorted(expected_fields - set(intelligence))
        if missing_intel:
            errors.append(f"{feature_id} intelligenceAssessment missing: {', '.join(missing_intel)}")
        allowed_levels = set(contract.get("intelligenceExpectationValues") or [])
        allowed_fits = set(contract.get("intelligenceFitValues") or [])
        for key in ("currentLevel", "expectedLevel"):
            if intelligence.get(key) not in allowed_levels:
                errors.append(f"{feature_id} has invalid intelligence {key}: {intelligence.get(key)!r}")
        if intelligence.get("fit") not in allowed_fits:
            errors.append(f"{feature_id} has invalid intelligence fit: {intelligence.get('fit')!r}")
        if not isinstance(intelligence.get("gaps"), list):
            errors.append(f"{feature_id} intelligence gaps must be an array")
        if not nonempty_string(intelligence.get("rationale")):
            errors.append(f"{feature_id} intelligence rationale is empty")

    if resolved.get("sourceDisposition") not in set(contract.get("sourceDispositionValues") or []):
        errors.append(f"{feature_id} has invalid sourceDisposition {resolved.get('sourceDisposition')!r}")
    if resolved.get("closureDecision") not in set(contract.get("closureDecisionValues") or []):
        errors.append(f"{feature_id} has invalid closureDecision {resolved.get('closureDecision')!r}")
    if resolved.get("closureDecision") in BLOCKING_DECISIONS and not resolved.get("correctiveIssue"):
        errors.append(f"{feature_id} blocking closure decision lacks correctiveIssue")
    if resolved.get("closureDecision") == "DEFER_NON_BLOCKING_IMPROVEMENT" and not nonempty_string(resolved.get("deferredImprovementOwner")):
        errors.append(f"{feature_id} deferred improvement lacks deferredImprovementOwner")
    if "uiDisposition" in override:
        effective["uiDisposition"] = override["uiDisposition"]
    return resolved


def main() -> int:
    parser = argparse.ArgumentParser()
    parser.add_argument(
        "--freeze",
        action="store_true",
        help="Require final T1 frozen state in addition to structural/effective audit completeness.",
    )
    args = parser.parse_args()

    ledger = load_json(LEDGER)
    scan1 = load_json(SCAN_1)
    scan2 = load_json(SCAN_2)
    reconciliation = load_json(RECONCILIATION)
    quality = load_json(QUALITY)
    contract = load_json(CONTRACT)

    errors: list[str] = []
    rules = ledger.get("rules") or {}
    profiles = ledger.get("assuranceProfiles") or {}
    allowed_states = set(rules.get("allowedFeatureStates") or [])
    allowed_dispositions = set(ledger.get("uiDispositionValues") or [])

    rows = ledger.get("features")
    if not isinstance(rows, list) or not rows:
        errors.append("features must be a non-empty array")
        rows = []

    physical: dict[str, dict] = {}
    for index, row in enumerate(rows):
        if not isinstance(row, dict):
            errors.append(f"features[{index}] is not an object")
            continue
        feature_id = str(row.get("id") or "").strip()
        if not feature_id:
            errors.append(f"features[{index}] has empty id")
            continue
        if feature_id in physical:
            errors.append(f"duplicate high-level feature id: {feature_id}")
            continue
        physical[feature_id] = row
        validate_base_row(row, feature_id, profiles, allowed_states, allowed_dispositions, errors)

    exclusions = reconciliation.get("excludedFutureSourceCarryForward") or []
    excluded_ids: set[str] = set()
    for item in exclusions:
        if not isinstance(item, dict):
            errors.append("excludedFutureSourceCarryForward item is not an object")
            continue
        item_id = str(item.get("id") or "").strip()
        if not item_id:
            errors.append("excludedFutureSourceCarryForward item has empty id")
            continue
        if item_id in excluded_ids:
            errors.append(f"duplicate future-source exclusion: {item_id}")
        excluded_ids.add(item_id)
        if item_id not in physical:
            errors.append(f"future-source exclusion {item_id} does not exist in discovery ledger")
        if not nonempty_string(item.get("reason")) or not nonempty_string_list(item.get("provenance")):
            errors.append(f"future-source exclusion {item_id} lacks reason/provenance")
        for owner in item.get("canonicalSourceOwners") or []:
            if not check_repo_path(owner):
                errors.append(f"future-source exclusion {item_id} owner missing at current head: {owner}")

    effective: dict[str, dict] = {
        feature_id: copy.deepcopy(row)
        for feature_id, row in physical.items()
        if feature_id not in excluded_ids
    }

    scan_counts: dict[str, int] = {}
    for scan_name, scan in ((SCAN_1.name, scan1), (SCAN_2.name, scan2)):
        omissions = scan.get("omissionsFound")
        if not isinstance(omissions, list):
            errors.append(f"{scan_name} omissionsFound must be an array")
            continue
        scan_counts[scan_name] = len(omissions)
        local_ids: set[str] = set()
        for item in omissions:
            if not isinstance(item, dict):
                errors.append(f"{scan_name} contains non-object omission")
                continue
            item_id = str(item.get("id") or "").strip()
            if item_id in local_ids:
                errors.append(f"{scan_name} duplicate omission id: {item_id}")
                continue
            local_ids.add(item_id)
            if item_id in excluded_ids:
                errors.append(f"{scan_name}:{item_id} is both shipped responsibility and future-source exclusion")
                continue
            if item_id in effective:
                errors.append(f"effective responsibility id appears more than once: {item_id}")
                continue
            resolved = resolve_scan_row(item, scan_name, physical, reconciliation, errors)
            if resolved is None:
                continue
            effective[item_id] = resolved
            validate_base_row(resolved, item_id, profiles, allowed_states, allowed_dispositions, errors)

    # Every reconciliation override must point at a real scan responsibility so stale
    # exception text cannot silently accumulate.
    scan_ids = {
        str(item.get("id") or "").strip()
        for scan in (scan1, scan2)
        for item in (scan.get("omissionsFound") or [])
        if isinstance(item, dict)
    }
    for override_id in (reconciliation.get("scanParentOverrides") or {}):
        if override_id not in scan_ids:
            errors.append(f"stale scanParentOverride does not match a scanned responsibility: {override_id}")

    quality_resolved: dict[str, dict] = {}
    for feature_id, row in effective.items():
        quality_row = validate_quality(row, quality, contract, errors)
        if quality_row is not None:
            quality_resolved[feature_id] = quality_row

    unknown_quality_overrides = sorted(set((quality.get("overrides") or {})) - set(effective))
    if unknown_quality_overrides:
        errors.append("quality overrides reference non-effective responsibilities: " + ", ".join(unknown_quality_overrides))

    for finding in quality.get("repositoryWideFindings") or []:
        if not isinstance(finding, dict):
            errors.append("repositoryWideFindings item is not an object")
            continue
        if finding.get("closureDecision") == "DEFER_NON_BLOCKING_IMPROVEMENT" and not nonempty_string(finding.get("deferredImprovementOwner")):
            errors.append(f"repository-wide deferred finding {finding.get('id')} lacks owner")

    freeze_requested = args.freeze or reconciliation.get("state") == "COMPLETE" or quality.get("state") == "COMPLETE" or ledger.get("discoveryComplete") is True
    if freeze_requested:
        if reconciliation.get("state") != "COMPLETE":
            errors.append("freeze requires T1_FINAL_RECONCILIATION state COMPLETE")
        if quality.get("state") != "COMPLETE":
            errors.append("freeze requires T1_QUALITY_AUDIT state COMPLETE")
        if ledger.get("discoveryComplete") is not True:
            errors.append("freeze requires feature-assurance-ledger discoveryComplete=true")
        if ledger.get("unexplainedGapCount") != 0:
            errors.append("freeze requires feature-assurance-ledger unexplainedGapCount == 0")
        if reconciliation.get("unexplainedGapCount") != 0:
            errors.append("freeze requires final reconciliation unexplainedGapCount == 0")
        unresolved = quality.get("unresolvedBlockingDecisions") or []
        if unresolved:
            errors.append("freeze has unresolved blocking decisions: " + ", ".join(str(x) for x in unresolved))
        blocking_features = [
            feature_id
            for feature_id, resolved in quality_resolved.items()
            if resolved.get("closureDecision") in BLOCKING_DECISIONS
        ]
        if blocking_features:
            errors.append("freeze has blocking closure decisions: " + ", ".join(sorted(blocking_features)))
        for feature_id, row in effective.items():
            if row.get("blockingStates"):
                errors.append(f"freeze but {feature_id} still has blockingStates")
        if any(item.get("state") != "QUALIFIED" for item in reconciliation.get("correctiveIssues") or [] if isinstance(item, dict)):
            errors.append("freeze requires every T1 correctiveIssues entry state QUALIFIED")

    if errors:
        print("V18 FEATURE REALITY & QUALITY AUDIT GATE: FAIL")
        for error in errors:
            print(f"- {error}")
        return 1

    print(
        "V18 FEATURE REALITY & QUALITY AUDIT GATE: PASS "
        f"({len(physical)} discovery rows; {len(excluded_ids)} explicit future-source exclusion; "
        f"{scan_counts.get(SCAN_1.name, 0)} scan-1 responsibilities; "
        f"{scan_counts.get(SCAN_2.name, 0)} scan-2 responsibilities; "
        f"{len(effective)} effective shipped-v18 responsibilities; freeze={freeze_requested})"
    )
    if quality.get("unresolvedBlockingDecisions"):
        print("T1 structural audit is complete enough for corrective qualification; freeze remains blocked by: " + ", ".join(quality["unresolvedBlockingDecisions"]))
    return 0


if __name__ == "__main__":
    sys.exit(main())
