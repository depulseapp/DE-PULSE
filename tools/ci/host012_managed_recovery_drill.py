#!/usr/bin/env python3
"""Run the provider-neutral HOST-012 managed PostgreSQL privacy replay drill.

The managed provider/operator must create the PITR-restored database before this
runner is invoked. Database URLs are accepted only through environment variables
so credentials are not placed in command history or emitted into evidence.

This runner:
1. binds the drill to the exact checked-out Git SHA;
2. executes the canonical tagged PostgreSQL privacy replay regression against a
   current post-delete source and a distinct pre-delete PITR restore target;
3. verifies the secret-free privacy artifact produced by that canonical test;
4. optionally binds real backup/restore artifacts to the final managed recovery
   evidence packet and runs the fail-closed evidence gate.

It never provisions, simulates, or claims a managed backup/PITR operation.
"""

from __future__ import annotations

import argparse
import hashlib
import json
import os
from pathlib import Path
import subprocess
import sys
from typing import Any

ROOT = Path(__file__).resolve().parents[2]
PRIVACY_SCHEMA = "DE.PULSE-HOST012-MANAGED-RECOVERY-PRIVACY-VERIFICATION-1"
FINAL_SCHEMA = "DE.PULSE-HOST012-MANAGED-RECOVERY-EVIDENCE-1"
ACK = "HOST012_MANAGED_PITR_OPERATOR_DRILL"
REQUIRED_DATABASE_ENV = (
    "DEPULSE_MANAGED_RECOVERY_SOURCE_URL",
    "DEPULSE_MANAGED_RECOVERY_RESTORED_URL",
    "DEPULSE_MANAGED_RECOVERY_USER_ID",
    "DEPULSE_MANAGED_RECOVERY_ENV_CLASS",
)


def fail(message: str) -> int:
    print("HOST-012 managed recovery drill: FAIL", file=sys.stderr)
    print(message, file=sys.stderr)
    return 1


def load_json(path: Path) -> dict[str, Any]:
    value = json.loads(path.read_text(encoding="utf-8"))
    if not isinstance(value, dict):
        raise ValueError(f"{path} must contain a JSON object")
    return value


def file_sha256(path: Path) -> str:
    digest = hashlib.sha256()
    with path.open("rb") as handle:
        for chunk in iter(lambda: handle.read(1024 * 1024), b""):
            digest.update(chunk)
    return digest.hexdigest()


def exact_head_sha() -> str:
    result = subprocess.run(
        ["git", "rev-parse", "HEAD"],
        cwd=ROOT,
        check=True,
        capture_output=True,
        text=True,
    )
    sha = result.stdout.strip().lower()
    if len(sha) != 40 or any(ch not in "0123456789abcdef" for ch in sha):
        raise RuntimeError("git rev-parse HEAD did not return a full SHA")
    return sha


def validate_privacy_artifact(data: dict[str, Any], candidate_sha: str) -> list[str]:
    errors: list[str] = []
    if data.get("schema") != PRIVACY_SCHEMA:
        errors.append(f"privacy artifact schema must be {PRIVACY_SCHEMA}")
    if data.get("candidateSha") != candidate_sha:
        errors.append("privacy artifact candidateSha does not match the exact checked-out head")
    if data.get("verificationMode") != "CANONICAL_ARCHIVE_TOMBSTONE_REPLAY":
        errors.append("privacy artifact must come from canonical archive tombstone replay")
    source = data.get("sourceTombstone")
    if not isinstance(source, dict) or source.get("present") is not True or not isinstance(source.get("updatedAt"), int) or source.get("updatedAt", 0) <= 0:
        errors.append("privacy artifact lacks an authoritative source tombstone")
    before = data.get("beforeReplay")
    if not isinstance(before, dict) or before.get("liveUserPresent") is not True or before.get("tombstonePresent") is not False:
        errors.append("privacy artifact does not prove the PITR target was from before deletion")
    replay = data.get("canonicalReplay")
    if not isinstance(replay, dict):
        errors.append("privacy artifact canonicalReplay is missing")
    else:
        if replay.get("owner") != "enforceArchiveAccountDeletionPrivacy":
            errors.append("privacy replay did not use the canonical anti-resurrection owner")
        if replay.get("restoreMode") != "replace":
            errors.append("privacy replay must use canonical replace-restore semantics")
        if replay.get("restartVerified") is not True:
            errors.append("privacy replay restart verification is required")
        if replay.get("sourceNotRestoreTarget") is not True:
            errors.append("privacy drill must keep the authoritative source distinct from the replay/restore target")
    privacy = data.get("privacyAssertions")
    if not isinstance(privacy, dict):
        errors.append("privacyAssertions is missing")
    else:
        for key in (
            "deletedUsersResurrected",
            "liveDeletedProfiles",
            "personalWorkspaceRowsForDeletedUsers",
            "activeSessionsForDeletedUsers",
            "activeDevicesForDeletedUsers",
        ):
            if privacy.get(key) != 0:
                errors.append(f"privacyAssertions.{key} must be 0")
        if not isinstance(privacy.get("tombstonesPresentAfterReplay"), int) or privacy.get("tombstonesPresentAfterReplay", 0) < 1:
            errors.append("privacyAssertions.tombstonesPresentAfterReplay must be >= 1")
        if privacy.get("canonicalHealthCheckPassed") is not True:
            errors.append("privacyAssertions.canonicalHealthCheckPassed must be true")
        if privacy.get("accountDataProjectionChecked") is not True:
            errors.append("privacyAssertions.accountDataProjectionChecked must be true")
    user_hash = data.get("userIdSha256")
    if not isinstance(user_hash, str) or len(user_hash) != 64 or any(ch not in "0123456789abcdef" for ch in user_hash.lower()):
        errors.append("privacy artifact must contain only a SHA-256 user identifier")
    return errors


def require_artifact_digest(evidence: dict[str, Any], kind: str, expected: str) -> str | None:
    artifacts = evidence.get("artifacts")
    if not isinstance(artifacts, list):
        return "final evidence artifacts must be an array"
    matches = [item for item in artifacts if isinstance(item, dict) and item.get("kind") == kind]
    if len(matches) != 1:
        return f"final evidence must contain exactly one {kind} artifact"
    digest = str(matches[0].get("sha256") or "").lower()
    if digest != expected.lower():
        return f"final evidence {kind} sha256 does not match the supplied artifact file"
    return None


def run_final_validation(
    evidence_path: Path,
    backup_artifact: Path,
    restore_artifact: Path,
    privacy_artifact: Path,
    candidate_sha: str,
    validation_out: Path | None,
) -> int:
    try:
        evidence = load_json(evidence_path)
    except (OSError, json.JSONDecodeError, ValueError) as exc:
        return fail(f"cannot read final evidence: {exc}")
    if evidence.get("schema") != FINAL_SCHEMA:
        return fail(f"final evidence schema must be {FINAL_SCHEMA}")

    artifact_checks = (
        ("backup-metadata", backup_artifact),
        ("restore-metadata", restore_artifact),
        ("post-restore-privacy-verification", privacy_artifact),
    )
    for kind, path in artifact_checks:
        if not path.is_file():
            return fail(f"{kind} artifact is missing: {path}")
        message = require_artifact_digest(evidence, kind, file_sha256(path))
        if message:
            return fail(message)

    gate = ROOT / "tools" / "ci" / "host012_managed_recovery_evidence_gate.py"
    command = [
        sys.executable,
        str(gate),
        "--evidence",
        str(evidence_path),
        "--expected-sha",
        candidate_sha,
    ]
    if validation_out is not None:
        command.extend(["--json-out", str(validation_out)])
    return subprocess.run(command, cwd=ROOT, check=False).returncode


def main(argv: list[str] | None = None) -> int:
    parser = argparse.ArgumentParser(description="Run the HOST-012 managed PostgreSQL privacy replay drill")
    parser.add_argument("--artifact", required=True, help="Secret-free privacy verification JSON output path")
    parser.add_argument(
        "--confirm-replace-restored-target",
        action="store_true",
        help="Required acknowledgement that the distinct PITR-restored target may be replaced by the canonical replay",
    )
    parser.add_argument("--final-evidence", help="Optional operator/provider final evidence JSON to validate after the drill")
    parser.add_argument("--backup-artifact", help="Provider backup metadata artifact file; required with --final-evidence")
    parser.add_argument("--restore-artifact", help="Provider restore metadata artifact file; required with --final-evidence")
    parser.add_argument("--validation-out", help="Optional machine-readable final evidence validation result")
    args = parser.parse_args(argv)

    if not args.confirm_replace_restored_target:
        return fail("--confirm-replace-restored-target is required; the PITR-restored target is destructively reconciled")
    missing = [name for name in REQUIRED_DATABASE_ENV if not os.environ.get(name, "").strip()]
    if missing:
        return fail("required environment variables are missing: " + ", ".join(missing))
    environment_class = os.environ["DEPULSE_MANAGED_RECOVERY_ENV_CLASS"].strip().lower()
    if environment_class == "production":
        return fail("production recovery targets are prohibited for this development drill")
    if environment_class not in {"development", "test", "stage"}:
        return fail("DEPULSE_MANAGED_RECOVERY_ENV_CLASS must be development, test, or stage")

    try:
        candidate_sha = exact_head_sha()
    except (OSError, subprocess.CalledProcessError, RuntimeError) as exc:
        return fail(f"cannot resolve exact Git head: {exc}")

    artifact_path = Path(args.artifact).expanduser().resolve()
    drill_env = os.environ.copy()
    drill_env["DEPULSE_MANAGED_RECOVERY_ACK"] = ACK
    drill_env["DEPULSE_MANAGED_RECOVERY_CANDIDATE_SHA"] = candidate_sha
    drill_env["DEPULSE_MANAGED_RECOVERY_PRIVACY_ARTIFACT_PATH"] = str(artifact_path)

    command = [
        "go",
        "test",
        "-tags",
        "postgres",
        "-count=1",
        "-run",
        "^TestHOST012ManagedRecoveryPrivacyReplayDrill$",
        ".",
    ]
    completed = subprocess.run(command, cwd=ROOT, env=drill_env, check=False)
    if completed.returncode != 0:
        return fail("canonical managed recovery privacy replay regression failed")

    try:
        privacy = load_json(artifact_path)
    except (OSError, json.JSONDecodeError, ValueError) as exc:
        return fail(f"cannot read privacy verification artifact: {exc}")
    privacy_errors = validate_privacy_artifact(privacy, candidate_sha)
    if privacy_errors:
        return fail("; ".join(privacy_errors))

    result: dict[str, Any] = {
        "status": "PASS",
        "candidateSha": candidate_sha,
        "privacyArtifactSha256": file_sha256(artifact_path),
        "privacyReplayVerified": True,
        "managedBackupPitrClaimed": False,
        "note": "The canonical privacy replay passed. Managed backup/PITR truth still comes from provider/operator artifacts and the final evidence gate.",
    }

    if args.final_evidence:
        if not args.backup_artifact or not args.restore_artifact:
            return fail("--backup-artifact and --restore-artifact are required with --final-evidence")
        validation_out = Path(args.validation_out).expanduser().resolve() if args.validation_out else None
        final_status = run_final_validation(
            Path(args.final_evidence).expanduser().resolve(),
            Path(args.backup_artifact).expanduser().resolve(),
            Path(args.restore_artifact).expanduser().resolve(),
            artifact_path,
            candidate_sha,
            validation_out,
        )
        if final_status != 0:
            return fail("provider/operator managed recovery evidence did not pass the final evidence gate")
        result["managedRecoveryEvidenceValidated"] = True
    else:
        result["managedRecoveryEvidenceValidated"] = False

    print(json.dumps(result, indent=2, sort_keys=True))
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
