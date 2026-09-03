#!/usr/bin/env python3
"""Self-tests for the HOST-012 managed recovery evidence validator."""

from __future__ import annotations

import copy
import sys

from host012_managed_recovery_evidence_gate import validate_evidence

TEST_SHA = "a" * 40
DIGEST_A = "1" * 64
DIGEST_B = "2" * 64
DIGEST_C = "3" * 64


def valid_evidence() -> dict:
    return {
        "schema": "DE.PULSE-HOST012-MANAGED-RECOVERY-EVIDENCE-1",
        "requirements": ["HOST-012", "HOST-016"],
        "candidateSha": TEST_SHA,
        "evidenceMode": "MANAGED_HOSTED_OPERATOR_DRILL",
        "environment": {
            "managed": True,
            "provider": "example-managed-postgres",
            "databaseService": "managed-postgresql",
            "class": "stage",
            "region": "ca-west",
            "sourceDatabaseId": "stage-primary-001",
            "restoredDatabaseId": "stage-pitr-restore-001",
        },
        "backup": {
            "encrypted": True,
            "pitrEnabled": True,
            "backupId": "backup-001",
            "restoreId": "restore-001",
            "recoveryWindowHours": 168,
            "retentionDays": 14,
            "requestedRestorePoint": "2026-08-27T09:00:00Z",
            "restoreCompletedAt": "2026-08-27T10:30:00Z",
        },
        "operator": {
            "drillId": "HOST012-RECOVERY-001",
            "role": "recovery-operator",
            "changeReference": "operator-drill-001",
            "executedAt": "2026-08-27T10:00:00Z",
        },
        "recovery": {
            "serviceBlockedDuringRecovery": True,
            "restoredFromPointBeforeDeletion": True,
            "authoritativeTombstoneSource": "canonical-post-delete-state-001",
            "tombstoneReplayCompleted": True,
            "serviceEnabledAfterVerification": True,
            "rpoTargetSeconds": 3600,
            "measuredRPOSeconds": 120,
            "rtoTargetSeconds": 3600,
            "measuredRTOSeconds": 900,
        },
        "privacyAssertions": {
            "deletedUsersResurrected": 0,
            "liveDeletedProfiles": 0,
            "personalWorkspaceRowsForDeletedUsers": 0,
            "activeSessionsForDeletedUsers": 0,
            "tombstonesPresentAfterReplay": 1,
            "canonicalHealthCheckPassed": True,
            "accountDataProjectionChecked": True,
        },
        "artifacts": [
            {"kind": "backup-metadata", "id": "provider-backup-001", "sha256": DIGEST_A},
            {"kind": "restore-metadata", "id": "provider-restore-001", "sha256": DIGEST_B},
            {"kind": "post-restore-privacy-verification", "id": "privacy-check-001", "sha256": DIGEST_C},
        ],
    }


def require_failure(name: str, mutate) -> None:
    evidence = valid_evidence()
    mutate(evidence)
    errors = validate_evidence(evidence, expected_sha=TEST_SHA)
    if not errors:
        raise AssertionError(f"{name}: invalid evidence unexpectedly passed")


def main() -> int:
    baseline_errors = validate_evidence(valid_evidence(), expected_sha=TEST_SHA)
    if baseline_errors:
        raise AssertionError(f"valid managed recovery evidence failed: {baseline_errors}")

    require_failure("simulated-mode", lambda e: e.__setitem__("evidenceMode", "SIMULATED"))
    require_failure("local-provider", lambda e: e["environment"].__setitem__("provider", "docker"))
    require_failure("same-restore-target", lambda e: e["environment"].__setitem__("restoredDatabaseId", e["environment"]["sourceDatabaseId"]))
    require_failure("missing-pitr", lambda e: e["backup"].__setitem__("pitrEnabled", False))
    require_failure("tombstone-replay-missing", lambda e: e["recovery"].__setitem__("tombstoneReplayCompleted", False))
    require_failure("resurrected-user", lambda e: e["privacyAssertions"].__setitem__("deletedUsersResurrected", 1))
    require_failure("rpo-breach", lambda e: e["recovery"].__setitem__("measuredRPOSeconds", 7200))
    require_failure("missing-provider-artifact", lambda e: e.__setitem__("artifacts", e["artifacts"][1:]))
    require_failure("secret-bearing-field", lambda e: e["environment"].__setitem__("databaseUrl", "postgres://forbidden"))

    mismatched = copy.deepcopy(valid_evidence())
    mismatch_errors = validate_evidence(mismatched, expected_sha="b" * 40)
    if not any("expected source SHA" in error for error in mismatch_errors):
        raise AssertionError(f"source-SHA mismatch was not rejected: {mismatch_errors}")

    print("HOST-012 managed recovery evidence gate self-test: PASS")
    return 0


if __name__ == "__main__":
    sys.exit(main())
