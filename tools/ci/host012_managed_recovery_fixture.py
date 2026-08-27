#!/usr/bin/env python3
"""Prepare the canonical HOST-012 managed PostgreSQL recovery lifecycle fixture.

The fixture deliberately performs only application-owned lifecycle operations on a
controlled non-production managed PostgreSQL source. It never provisions, creates,
simulates, or claims a provider backup/PITR restore.

Recommended operator sequence:
1. run ``seed`` to create an active account with session, trusted device and
   personal workspace state through canonical DE.PULSE lifecycle owners;
2. ensure the managed provider has a recoverable backup/PITR point containing the
   seeded state;
3. run ``delete`` with the seed artifact to exercise canonical durable account
   deletion on the current source;
4. restore a point from before deletion to a distinct target through the provider;
5. run ``host012_managed_recovery_drill.py`` against current source + restored
   target, optionally supplying the delete fixture artifact for cross-checking.

Database URLs are accepted only through environment variables and are never
printed. Generated passwords and device fingerprints stay in-process and are not
written to evidence.
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
SCHEMA = "DE.PULSE-HOST012-MANAGED-RECOVERY-LIFECYCLE-FIXTURE-1"
ACK = "HOST012_MANAGED_RECOVERY_LIFECYCLE_FIXTURE"
SOURCE_ENV = "DEPULSE_MANAGED_RECOVERY_SOURCE_URL"
ENV_CLASS_ENV = "DEPULSE_MANAGED_RECOVERY_ENV_CLASS"
USER_ENV = "DEPULSE_MANAGED_RECOVERY_USER_ID"


def fail(message: str) -> int:
    print("HOST-012 managed recovery lifecycle fixture: FAIL", file=sys.stderr)
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


def user_sha256(user_id: str) -> str:
    return hashlib.sha256(user_id.encode("utf-8")).hexdigest()


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


def validate_artifact(
    data: dict[str, Any],
    *,
    candidate_sha: str,
    environment_class: str,
    phase: str,
    expected_user_id: str | None = None,
    expected_database_identity: str | None = None,
) -> list[str]:
    errors: list[str] = []
    if data.get("schema") != SCHEMA:
        errors.append(f"fixture artifact schema must be {SCHEMA}")
    if data.get("candidateSha") != candidate_sha:
        errors.append("fixture artifact candidateSha does not match the exact checked-out head")
    if data.get("environmentClass") != environment_class:
        errors.append("fixture artifact environmentClass does not match the requested environment")
    if data.get("phase") != phase:
        errors.append(f"fixture artifact phase must be {phase}")
    if data.get("providerBackupPitrClaimed") is not False:
        errors.append("fixture artifact must not claim provider backup/PITR truth")

    database_identity = data.get("databaseIdentity")
    if not isinstance(database_identity, str) or not database_identity.strip():
        errors.append("fixture artifact databaseIdentity is required")
    elif expected_database_identity is not None and database_identity != expected_database_identity:
        errors.append("fixture artifact databaseIdentity does not match the seed source")

    user_id = data.get("userId")
    user_hash = data.get("userIdSha256")
    if not isinstance(user_id, str) or not user_id.strip():
        errors.append("fixture artifact userId is required")
    else:
        if expected_user_id is not None and user_id != expected_user_id:
            errors.append("fixture artifact userId does not match the seeded account")
        if user_hash != user_sha256(user_id):
            errors.append("fixture artifact userIdSha256 does not match userId")

    if not isinstance(data.get("tenantId"), str) or not data.get("tenantId", "").strip():
        errors.append("fixture artifact tenantId is required")
    if not isinstance(data.get("canonicalLifecycleOwner"), str) or not data.get("canonicalLifecycleOwner", "").strip():
        errors.append("fixture artifact canonicalLifecycleOwner is required")
    if not isinstance(data.get("createdAt"), str) or not data.get("createdAt", "").strip():
        errors.append("fixture artifact createdAt is required")

    if phase == "seed":
        if data.get("accountStatus") != "ACTIVE":
            errors.append("seed fixture accountStatus must be ACTIVE")
        if not isinstance(data.get("deviceCount"), int) or data.get("deviceCount", 0) < 1:
            errors.append("seed fixture must contain at least one trusted device")
        if not isinstance(data.get("sessionCount"), int) or data.get("sessionCount", 0) < 1:
            errors.append("seed fixture must contain at least one session")
        if data.get("personalWorkspacePresent") is not True:
            errors.append("seed fixture must contain personal workspace state")
        if data.get("tombstoneReason") not in (None, ""):
            errors.append("seed fixture must not contain a tombstone reason")
    elif phase == "delete":
        if data.get("accountStatus") != "DISABLED":
            errors.append("delete fixture accountStatus must be DISABLED")
        if data.get("tombstoneReason") != "USER_REQUESTED":
            errors.append("delete fixture must contain the canonical USER_REQUESTED tombstone reason")
        if not isinstance(data.get("tombstoneDeletedAt"), int) or data.get("tombstoneDeletedAt", 0) <= 0:
            errors.append("delete fixture must contain a durable tombstone timestamp")
        if data.get("deviceCount") != 0:
            errors.append("delete fixture deviceCount must be 0")
        if data.get("sessionCount") != 0:
            errors.append("delete fixture sessionCount must be 0")
        if data.get("personalWorkspacePresent") is not False:
            errors.append("delete fixture must prove the workspace is privacy blank")
    else:
        errors.append(f"unsupported fixture phase {phase}")
    return errors


def main(argv: list[str] | None = None) -> int:
    parser = argparse.ArgumentParser(description="Prepare the HOST-012 managed recovery lifecycle fixture")
    parser.add_argument("phase", choices=("seed", "delete"))
    parser.add_argument("--artifact", required=True, help="Secret-free fixture evidence JSON output path")
    parser.add_argument(
        "--seed-artifact",
        help="Seed artifact from the same exact candidate/source; required for delete unless DEPULSE_MANAGED_RECOVERY_USER_ID is set",
    )
    parser.add_argument(
        "--confirm-mutate-source",
        action="store_true",
        help="Required acknowledgement that the controlled non-production managed source will be mutated",
    )
    args = parser.parse_args(argv)

    if not args.confirm_mutate_source:
        return fail("--confirm-mutate-source is required")
    missing = [name for name in (SOURCE_ENV, ENV_CLASS_ENV) if not os.environ.get(name, "").strip()]
    if missing:
        return fail("required environment variables are missing: " + ", ".join(missing))
    environment_class = os.environ[ENV_CLASS_ENV].strip().lower()
    if environment_class == "production":
        return fail("production sources are prohibited for this development fixture")
    if environment_class not in {"development", "test", "stage"}:
        return fail(f"{ENV_CLASS_ENV} must be development, test, or stage")

    try:
        candidate_sha = exact_head_sha()
    except (OSError, subprocess.CalledProcessError, RuntimeError) as exc:
        return fail(f"cannot resolve exact Git head: {exc}")

    seed: dict[str, Any] | None = None
    user_id = os.environ.get(USER_ENV, "").strip()
    if args.phase == "delete":
        if args.seed_artifact:
            seed_path = Path(args.seed_artifact).expanduser().resolve()
            try:
                seed = load_json(seed_path)
            except (OSError, json.JSONDecodeError, ValueError) as exc:
                return fail(f"cannot read seed fixture artifact: {exc}")
            seed_errors = validate_artifact(
                seed,
                candidate_sha=candidate_sha,
                environment_class=environment_class,
                phase="seed",
            )
            if seed_errors:
                return fail("invalid seed fixture artifact: " + "; ".join(seed_errors))
            seed_user = str(seed.get("userId") or "").strip()
            if user_id and user_id != seed_user:
                return fail(f"{USER_ENV} does not match --seed-artifact")
            user_id = seed_user
        if not user_id:
            return fail(f"delete requires --seed-artifact or {USER_ENV}")

    artifact_path = Path(args.artifact).expanduser().resolve()
    fixture_env = os.environ.copy()
    fixture_env["DEPULSE_MANAGED_RECOVERY_FIXTURE_ACK"] = ACK
    fixture_env["DEPULSE_MANAGED_RECOVERY_FIXTURE_PHASE"] = args.phase
    fixture_env["DEPULSE_MANAGED_RECOVERY_CANDIDATE_SHA"] = candidate_sha
    fixture_env["DEPULSE_MANAGED_RECOVERY_FIXTURE_ARTIFACT_PATH"] = str(artifact_path)
    if user_id:
        fixture_env[USER_ENV] = user_id

    command = [
        "go",
        "test",
        "-tags",
        "postgres",
        "-count=1",
        "-run",
        "^TestHOST012ManagedRecoveryLifecycleFixture$",
        ".",
    ]
    completed = subprocess.run(command, cwd=ROOT, env=fixture_env, check=False)
    if completed.returncode != 0:
        return fail(f"canonical managed recovery lifecycle {args.phase} fixture failed")

    try:
        artifact = load_json(artifact_path)
    except (OSError, json.JSONDecodeError, ValueError) as exc:
        return fail(f"cannot read lifecycle fixture artifact: {exc}")
    artifact_errors = validate_artifact(
        artifact,
        candidate_sha=candidate_sha,
        environment_class=environment_class,
        phase=args.phase,
        expected_user_id=user_id or None,
        expected_database_identity=str(seed.get("databaseIdentity")) if seed is not None else None,
    )
    if artifact_errors:
        return fail("; ".join(artifact_errors))

    result: dict[str, Any] = {
        "status": "PASS",
        "phase": args.phase,
        "candidateSha": candidate_sha,
        "artifactSha256": file_sha256(artifact_path),
        "userId": artifact["userId"],
        "userIdSha256": artifact["userIdSha256"],
        "managedBackupPitrClaimed": False,
    }
    if args.phase == "seed":
        result["nextAction"] = "Ensure provider backup/PITR can recover this live seeded state before running the delete phase."
    else:
        result["nextAction"] = "Restore a provider PITR point from before deletion to a distinct target, then run the HOST-012 privacy replay drill."
    print(json.dumps(result, indent=2, sort_keys=True))
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
