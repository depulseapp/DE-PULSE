#!/usr/bin/env python3
"""Historical identity contract for the conserved v17/v18 reconciliation ledger."""

from __future__ import annotations

from copy import deepcopy

EXPECTED_SCHEMA = "DE.PULSE-HISTORICAL-RECONCILIATION-IDENTITY-1"
EXPECTED_LEDGER_SCHEMA = "DE.PULSE-V18.5.1-V17-V18-IMPLEMENTATION-RECONCILIATION-1"
EXPECTED_RELEASE = "v18.5.1"
EXPECTED_STABLE_TAG = "v18.5.0-stable"
EXPECTED_STABLE_COMMIT = "0d37ca35f5fc3ad89cebed506cc5a4c2d6a7a680"
EXPECTED_BRANCH = "v18.5.1-development"
EXPECTED_LEDGER_PATH = "release/v18.5.1/V17-V18-IMPLEMENTATION-RECONCILIATION.json"
CURRENT_RELEASE_IDENTITY_PATH = "release_identity.json"


def historical_identity_errors(ledger: dict, identity: dict) -> list[str]:
    errors: list[str] = []
    historical = identity.get("historicalReconciliation", {})
    separation = identity.get("separationContract", {})
    baseline = ledger.get("baseline", {})

    if identity.get("schema") != EXPECTED_SCHEMA:
        errors.append("historical reconciliation identity schema drift")
    if identity.get("identityRole") != "IMMUTABLE_HISTORICAL_RECONCILIATION_BASELINE":
        errors.append("historical reconciliation identity role drift")
    if historical.get("ledgerPath") != EXPECTED_LEDGER_PATH:
        errors.append("historical reconciliation ledger path drift")
    if historical.get("ledgerSchema") != EXPECTED_LEDGER_SCHEMA:
        errors.append("historical reconciliation ledger schema binding drift")
    if historical.get("reconciliationRelease") != EXPECTED_RELEASE:
        errors.append("historical reconciliation release identity drift")
    if historical.get("incomingStableTag") != EXPECTED_STABLE_TAG:
        errors.append("historical incoming Stable tag drift")
    if historical.get("incomingStableCommit") != EXPECTED_STABLE_COMMIT:
        errors.append("historical incoming Stable commit drift")
    if historical.get("reconciliationBranch") != EXPECTED_BRANCH:
        errors.append("historical reconciliation branch drift")

    if separation.get("currentReleaseIdentityPath") != CURRENT_RELEASE_IDENTITY_PATH:
        errors.append("current release identity separation path drift")
    if separation.get("historicalBaselineSource") != "HISTORICAL_IDENTITY_MANIFEST_ONLY":
        errors.append("historical baseline source policy drift")
    if separation.get("currentReleaseIdentityRole") != "CURRENT_RELEASE_ONLY":
        errors.append("current release identity role policy drift")
    if separation.get("deriveHistoricalBaselineFromCurrentReleaseIdentity") is not False:
        errors.append("historical baseline must never derive from current release identity")

    if ledger.get("schema") != historical.get("ledgerSchema"):
        errors.append("ledger schema differs from immutable historical identity")
    if ledger.get("release") != historical.get("reconciliationRelease"):
        errors.append("ledger release differs from immutable historical identity")
    if baseline.get("currentStableTag") != historical.get("incomingStableTag"):
        errors.append("ledger historical Stable tag differs from immutable identity")
    if baseline.get("currentStableCommit") != historical.get("incomingStableCommit"):
        errors.append("ledger historical Stable commit differs from immutable identity")
    if baseline.get("reconciliationBranch") != historical.get("reconciliationBranch"):
        errors.append("ledger historical reconciliation branch differs from immutable identity")
    return errors


def historical_identity_self_test_errors(identity: dict) -> list[str]:
    """Cheap mutation tests proving historical identity does not follow current release state."""
    errors: list[str] = []
    canonical = {
        "schema": EXPECTED_LEDGER_SCHEMA,
        "release": EXPECTED_RELEASE,
        "baseline": {
            "currentStableTag": EXPECTED_STABLE_TAG,
            "currentStableCommit": EXPECTED_STABLE_COMMIT,
            "reconciliationBranch": EXPECTED_BRANCH,
        },
    }
    if historical_identity_errors(canonical, identity):
        errors.append("historical identity canonical self-test unexpectedly failed")
        return errors

    mutations = (
        ("release", "v99.0.0"),
        ("baseline.currentStableTag", "v99.0.0-stable"),
        ("baseline.currentStableCommit", "f" * 40),
        ("baseline.reconciliationBranch", "v99.0.0-development"),
    )
    for field, value in mutations:
        candidate = deepcopy(canonical)
        if field == "release":
            candidate["release"] = value
        else:
            candidate["baseline"][field.split(".", 1)[1]] = value
        if not historical_identity_errors(candidate, identity):
            errors.append(f"historical identity self-test failed to reject mutation: {field}")

    # Current release identity is deliberately not an input to historical_identity_errors.
    # This locks separation without coupling the historical baseline to whichever release
    # is current when the reconciliation gate runs.
    separation = identity.get("separationContract", {})
    if separation.get("deriveHistoricalBaselineFromCurrentReleaseIdentity") is not False:
        errors.append("historical/current release separation self-test failed")
    return errors
