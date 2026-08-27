#!/usr/bin/env python3
"""Validate real managed-hosted recovery evidence for the HOST-012 privacy dependency.

This gate deliberately does not perform or simulate a backup/PITR restore. It validates the
machine-readable evidence packet produced by an operator-run managed PostgreSQL recovery
exercise. Passing this gate is necessary evidence hygiene, but it is not a substitute for the
provider/operator drill that produced the packet.
"""

from __future__ import annotations

import argparse
import hashlib
import json
import re
import sys
from datetime import datetime
from pathlib import Path
from typing import Any

SCHEMA = "DE.PULSE-HOST012-MANAGED-RECOVERY-EVIDENCE-1"
REQUIRED_REQUIREMENTS = {"HOST-012", "HOST-016"}
EVIDENCE_MODE = "MANAGED_HOSTED_OPERATOR_DRILL"
SHA_RE = re.compile(r"^[0-9a-f]{40}$")
SHA256_RE = re.compile(r"^[0-9a-f]{64}$")
DISALLOWED_PROVIDER_VALUES = {
    "ci",
    "docker",
    "github-actions",
    "local",
    "localhost",
    "mock",
    "simulated",
    "testcontainer",
}
SECRET_KEY_NAMES = {
    "apikey",
    "connectionstring",
    "databaseurl",
    "password",
    "privatekey",
    "rawsessiontoken",
    "secret",
    "sessiontoken",
    "token",
}
REQUIRED_ARTIFACT_KINDS = {
    "backup-metadata",
    "restore-metadata",
    "post-restore-privacy-verification",
}


def _norm_key(value: str) -> str:
    return re.sub(r"[^a-z0-9]", "", value.lower())


def _nonempty(value: Any) -> bool:
    return isinstance(value, str) and bool(value.strip())


def _positive_number(value: Any) -> bool:
    return isinstance(value, (int, float)) and not isinstance(value, bool) and value > 0


def _nonnegative_number(value: Any) -> bool:
    return isinstance(value, (int, float)) and not isinstance(value, bool) and value >= 0


def _parse_timestamp(value: Any, field: str, errors: list[str]) -> datetime | None:
    if not _nonempty(value):
        errors.append(f"{field} must be a non-empty ISO-8601 timestamp")
        return None
    text = value.strip()
    try:
        if text.endswith("Z"):
            text = text[:-1] + "+00:00"
        parsed = datetime.fromisoformat(text)
    except ValueError:
        errors.append(f"{field} must be a valid ISO-8601 timestamp")
        return None
    if parsed.tzinfo is None:
        errors.append(f"{field} must include a timezone")
        return None
    return parsed


def _walk_for_secret_keys(value: Any, path: str, errors: list[str]) -> None:
    if isinstance(value, dict):
        for key, child in value.items():
            normalized = _norm_key(str(key))
            if normalized in SECRET_KEY_NAMES:
                errors.append(f"secret-bearing field is forbidden in recovery evidence: {path}{key}")
            _walk_for_secret_keys(child, f"{path}{key}.", errors)
    elif isinstance(value, list):
        for index, child in enumerate(value):
            _walk_for_secret_keys(child, f"{path}[{index}].", errors)


def _require_dict(parent: dict[str, Any], key: str, errors: list[str]) -> dict[str, Any]:
    value = parent.get(key)
    if not isinstance(value, dict):
        errors.append(f"{key} must be an object")
        return {}
    return value


def validate_evidence(data: Any, expected_sha: str | None = None) -> list[str]:
    errors: list[str] = []
    if not isinstance(data, dict):
        return ["evidence root must be an object"]

    _walk_for_secret_keys(data, "", errors)

    if data.get("schema") != SCHEMA:
        errors.append(f"schema must be {SCHEMA}")
    if data.get("evidenceMode") != EVIDENCE_MODE:
        errors.append(f"evidenceMode must be {EVIDENCE_MODE}; local/simulated proof is not accepted")

    requirements = data.get("requirements")
    if not isinstance(requirements, list) or not REQUIRED_REQUIREMENTS.issubset(set(requirements)):
        errors.append("requirements must include HOST-012 and HOST-016")

    candidate_sha = data.get("candidateSha")
    if not _nonempty(candidate_sha) or not SHA_RE.fullmatch(candidate_sha.strip().lower()):
        errors.append("candidateSha must be a 40-character lowercase Git SHA")
    elif expected_sha and candidate_sha.strip().lower() != expected_sha.strip().lower():
        errors.append("candidateSha does not match the expected source SHA")

    environment = _require_dict(data, "environment", errors)
    if environment.get("managed") is not True:
        errors.append("environment.managed must be true")
    provider = environment.get("provider")
    if not _nonempty(provider):
        errors.append("environment.provider is required")
    elif provider.strip().lower() in DISALLOWED_PROVIDER_VALUES:
        errors.append("environment.provider must identify a real managed database provider, not local/CI/simulated infrastructure")
    if not _nonempty(environment.get("databaseService")):
        errors.append("environment.databaseService is required")
    if not _nonempty(environment.get("region")):
        errors.append("environment.region is required to make residency explicit")
    if environment.get("class") not in {"development", "test", "stage", "production"}:
        errors.append("environment.class must be development, test, stage, or production")
    source_database_id = environment.get("sourceDatabaseId")
    restored_database_id = environment.get("restoredDatabaseId")
    if not _nonempty(source_database_id) or not _nonempty(restored_database_id):
        errors.append("environment sourceDatabaseId and restoredDatabaseId are required")
    elif source_database_id.strip() == restored_database_id.strip():
        errors.append("restoredDatabaseId must identify a distinct recovery target")

    backup = _require_dict(data, "backup", errors)
    if backup.get("encrypted") is not True:
        errors.append("backup.encrypted must be true")
    if backup.get("pitrEnabled") is not True:
        errors.append("backup.pitrEnabled must be true")
    for key in ("backupId", "restoreId"):
        if not _nonempty(backup.get(key)):
            errors.append(f"backup.{key} is required")
    if not _positive_number(backup.get("recoveryWindowHours")):
        errors.append("backup.recoveryWindowHours must be > 0")
    if not _positive_number(backup.get("retentionDays")):
        errors.append("backup.retentionDays must be > 0")
    restore_point = _parse_timestamp(backup.get("requestedRestorePoint"), "backup.requestedRestorePoint", errors)
    restore_completed = _parse_timestamp(backup.get("restoreCompletedAt"), "backup.restoreCompletedAt", errors)
    if restore_point and restore_completed and restore_point >= restore_completed:
        errors.append("backup.requestedRestorePoint must precede backup.restoreCompletedAt")

    operator = _require_dict(data, "operator", errors)
    for key in ("drillId", "role", "changeReference"):
        if not _nonempty(operator.get(key)):
            errors.append(f"operator.{key} is required")
    executed_at = _parse_timestamp(operator.get("executedAt"), "operator.executedAt", errors)
    if executed_at and restore_completed and executed_at > restore_completed:
        errors.append("operator.executedAt must not be after backup.restoreCompletedAt")

    recovery = _require_dict(data, "recovery", errors)
    if recovery.get("serviceBlockedDuringRecovery") is not True:
        errors.append("recovery.serviceBlockedDuringRecovery must be true")
    if recovery.get("restoredFromPointBeforeDeletion") is not True:
        errors.append("recovery.restoredFromPointBeforeDeletion must be true for anti-resurrection proof")
    if not _nonempty(recovery.get("authoritativeTombstoneSource")):
        errors.append("recovery.authoritativeTombstoneSource is required")
    if recovery.get("tombstoneReplayCompleted") is not True:
        errors.append("recovery.tombstoneReplayCompleted must be true before service enablement")
    if recovery.get("serviceEnabledAfterVerification") is not True:
        errors.append("recovery.serviceEnabledAfterVerification must be true")

    for metric in ("rpo", "rto"):
        target_key = f"{metric}TargetSeconds"
        measured_key = f"measured{metric.upper()}Seconds"
        target = recovery.get(target_key)
        measured = recovery.get(measured_key)
        if not _positive_number(target):
            errors.append(f"recovery.{target_key} must be > 0")
        if not _nonnegative_number(measured):
            errors.append(f"recovery.{measured_key} must be >= 0")
        if _positive_number(target) and _nonnegative_number(measured) and measured > target:
            errors.append(f"recovery.{measured_key} exceeds {target_key}")

    privacy = _require_dict(data, "privacyAssertions", errors)
    for key in (
        "deletedUsersResurrected",
        "liveDeletedProfiles",
        "personalWorkspaceRowsForDeletedUsers",
        "activeSessionsForDeletedUsers",
    ):
        if privacy.get(key) != 0:
            errors.append(f"privacyAssertions.{key} must be 0")
    if not isinstance(privacy.get("tombstonesPresentAfterReplay"), int) or privacy.get("tombstonesPresentAfterReplay", 0) < 1:
        errors.append("privacyAssertions.tombstonesPresentAfterReplay must be >= 1")
    if privacy.get("canonicalHealthCheckPassed") is not True:
        errors.append("privacyAssertions.canonicalHealthCheckPassed must be true")
    if privacy.get("accountDataProjectionChecked") is not True:
        errors.append("privacyAssertions.accountDataProjectionChecked must be true")

    artifacts = data.get("artifacts")
    seen_kinds: set[str] = set()
    if not isinstance(artifacts, list) or not artifacts:
        errors.append("artifacts must contain provider/operator evidence references")
    else:
        for index, artifact in enumerate(artifacts):
            if not isinstance(artifact, dict):
                errors.append(f"artifacts[{index}] must be an object")
                continue
            kind = artifact.get("kind")
            if _nonempty(kind):
                seen_kinds.add(kind.strip())
            else:
                errors.append(f"artifacts[{index}].kind is required")
            if not _nonempty(artifact.get("id")):
                errors.append(f"artifacts[{index}].id is required")
            digest = artifact.get("sha256")
            if not _nonempty(digest) or not SHA256_RE.fullmatch(digest.strip().lower()):
                errors.append(f"artifacts[{index}].sha256 must be a 64-character digest")
        missing_kinds = REQUIRED_ARTIFACT_KINDS - seen_kinds
        if missing_kinds:
            errors.append("artifacts missing required kinds: " + ", ".join(sorted(missing_kinds)))

    return errors


def evidence_digest(data: Any) -> str:
    canonical = json.dumps(data, sort_keys=True, separators=(",", ":"), ensure_ascii=False).encode("utf-8")
    return hashlib.sha256(canonical).hexdigest()


def main(argv: list[str] | None = None) -> int:
    parser = argparse.ArgumentParser(description="Validate HOST-012 managed PostgreSQL recovery evidence")
    parser.add_argument("--evidence", required=True, help="Path to operator-produced recovery evidence JSON")
    parser.add_argument("--expected-sha", help="Require evidence to bind to this exact source SHA")
    parser.add_argument("--json-out", help="Optional path for a machine-readable validation result")
    args = parser.parse_args(argv)

    evidence_path = Path(args.evidence)
    try:
        data = json.loads(evidence_path.read_text(encoding="utf-8"))
    except (OSError, json.JSONDecodeError) as exc:
        result = {"status": "FAIL", "errors": [f"cannot read evidence: {exc}"]}
    else:
        errors = validate_evidence(data, expected_sha=args.expected_sha)
        result = {
            "status": "PASS" if not errors else "FAIL",
            "schema": SCHEMA,
            "candidateSha": data.get("candidateSha") if isinstance(data, dict) else None,
            "evidenceSha256": evidence_digest(data),
            "errors": errors,
            "note": "PASS validates evidence completeness; the managed provider/operator drill remains the source of infrastructure truth.",
        }

    rendered = json.dumps(result, indent=2, sort_keys=True)
    if args.json_out:
        Path(args.json_out).write_text(rendered + "\n", encoding="utf-8")
    print(rendered)
    return 0 if result["status"] == "PASS" else 1


if __name__ == "__main__":
    sys.exit(main())
