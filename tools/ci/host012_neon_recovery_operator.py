#!/usr/bin/env python3
"""Run the real HOST-012 recovery drill against a controlled Neon project.

This file is provider-specific operator glue only. Canonical DE.PULSE lifecycle,
anti-resurrection replay and final evidence validation remain owned by:
- host012_managed_recovery_fixture.py
- host012_managed_recovery_drill.py
- host012_managed_recovery_evidence_gate.py

The Neon API key is accepted only through NEON_API_KEY. Connection URIs returned
by Neon remain in process memory/environment and are never written to evidence or
printed. This operator refuses production.
"""
from __future__ import annotations

import argparse
from datetime import datetime, timezone
import hashlib
import json
import os
from pathlib import Path
import subprocess
import sys
import time
from typing import Any
from urllib.error import HTTPError, URLError
from urllib.parse import parse_qsl, urlencode, urlparse, urlunparse
from urllib.request import Request, urlopen

import host012_managed_recovery_drill as canonical_drill

ROOT = Path(__file__).resolve().parents[2]
API_BASE = "https://console.neon.tech/api/v2"
PROVIDER = "Neon"
DATABASE_SERVICE = "Neon Postgres"
EVIDENCE_SCHEMA = "DE.PULSE-HOST012-MANAGED-RECOVERY-EVIDENCE-1"
API_KEY_ENV = "NEON_API_KEY"
SOURCE_URL_ENV = "DEPULSE_MANAGED_RECOVERY_SOURCE_URL"
RESTORED_URL_ENV = "DEPULSE_MANAGED_RECOVERY_RESTORED_URL"
USER_ENV = "DEPULSE_MANAGED_RECOVERY_USER_ID"
ENV_CLASS_ENV = "DEPULSE_MANAGED_RECOVERY_ENV_CLASS"
ENCRYPTION_BASIS = "https://neon.com/security"


class OperatorError(RuntimeError):
    pass


def utc_now() -> datetime:
    return datetime.now(timezone.utc)


def iso_utc(value: datetime) -> str:
    return value.astimezone(timezone.utc).isoformat(timespec="milliseconds").replace("+00:00", "Z")


def parse_iso(value: str) -> datetime:
    text = str(value).strip()
    if text.endswith("Z"):
        text = text[:-1] + "+00:00"
    parsed = datetime.fromisoformat(text)
    if parsed.tzinfo is None:
        raise OperatorError("timestamp must contain timezone information")
    return parsed.astimezone(timezone.utc)


def load_json(path: Path) -> dict[str, Any]:
    value = json.loads(path.read_text(encoding="utf-8"))
    if not isinstance(value, dict):
        raise OperatorError(f"{path} must contain a JSON object")
    return value


def write_json(path: Path, value: dict[str, Any]) -> None:
    path.parent.mkdir(parents=True, exist_ok=True)
    path.write_text(json.dumps(value, indent=2, sort_keys=True) + "\n", encoding="utf-8")
    try:
        path.chmod(0o600)
    except OSError:
        pass


def file_sha256(path: Path) -> str:
    digest = hashlib.sha256()
    with path.open("rb") as handle:
        for chunk in iter(lambda: handle.read(1024 * 1024), b""):
            digest.update(chunk)
    return digest.hexdigest()


def exact_head_sha() -> str:
    result = subprocess.run(
        ["git", "rev-parse", "HEAD"], cwd=ROOT, check=True, capture_output=True, text=True
    )
    sha = result.stdout.strip().lower()
    if len(sha) != 40 or any(ch not in "0123456789abcdef" for ch in sha):
        raise OperatorError("git rev-parse HEAD did not return a full SHA")
    return sha


def force_verify_full(raw: str) -> str:
    parsed = urlparse(raw.strip())
    if parsed.scheme not in {"postgres", "postgresql"} or not parsed.hostname:
        raise OperatorError("Neon returned an invalid PostgreSQL connection URI")
    query = dict(parse_qsl(parsed.query, keep_blank_values=True))
    query["sslmode"] = "verify-full"
    return urlunparse((parsed.scheme, parsed.netloc, parsed.path, parsed.params, urlencode(query), parsed.fragment))


def sanitized_http_error(exc: HTTPError) -> str:
    detail = ""
    try:
        body = exc.read().decode("utf-8", errors="replace")
        parsed = json.loads(body)
        if isinstance(parsed, dict):
            detail = str(parsed.get("message") or parsed.get("error") or "")
    except Exception:
        detail = ""
    suffix = f": {detail[:300]}" if detail else ""
    return f"Neon API HTTP {exc.code}{suffix}"


def api_request(
    token: str,
    method: str,
    path: str,
    *,
    query: dict[str, str] | None = None,
    body: dict[str, Any] | None = None,
    timeout: float = 30.0,
) -> dict[str, Any]:
    url = API_BASE + path
    if query:
        url += "?" + urlencode(query)
    payload = None if body is None else json.dumps(body).encode("utf-8")
    request = Request(
        url,
        data=payload,
        method=method,
        headers={
            "Authorization": f"Bearer {token}",
            "Accept": "application/json",
            "Content-Type": "application/json",
            "User-Agent": "DE.PULSE-HOST012-Recovery-Operator/1",
        },
    )
    try:
        with urlopen(request, timeout=timeout) as response:
            raw = response.read()
    except HTTPError as exc:
        raise OperatorError(sanitized_http_error(exc)) from exc
    except URLError as exc:
        raise OperatorError(f"Neon API network failure: {exc.reason}") from exc
    if not raw:
        return {}
    decoded = json.loads(raw.decode("utf-8"))
    if not isinstance(decoded, dict):
        raise OperatorError("Neon API returned a non-object JSON response")
    return decoded


def select_database(token: str, project_id: str, branch_id: str, requested: str | None) -> str:
    if requested:
        return requested
    response = api_request(token, "GET", f"/projects/{project_id}/branches/{branch_id}/databases")
    databases = response.get("databases")
    if not isinstance(databases, list) or not databases:
        raise OperatorError("Neon source branch has no database")
    names = [str(item.get("name") or "") for item in databases if isinstance(item, dict)]
    if "neondb" in names:
        return "neondb"
    for name in names:
        if name:
            return name
    raise OperatorError("Neon database list contained no usable name")


def select_role(token: str, project_id: str, branch_id: str, requested: str | None) -> str:
    if requested:
        return requested
    response = api_request(token, "GET", f"/projects/{project_id}/branches/{branch_id}/roles")
    roles = response.get("roles")
    if not isinstance(roles, list) or not roles:
        raise OperatorError("Neon source branch has no role")
    names = [str(item.get("name") or "") for item in roles if isinstance(item, dict)]
    if "neondb_owner" in names:
        return "neondb_owner"
    for name in names:
        if name:
            return name
    raise OperatorError("Neon role list contained no usable name")


def connection_uri(token: str, project_id: str, branch_id: str, database: str, role: str) -> str:
    response = api_request(
        token,
        "GET",
        f"/projects/{project_id}/connection_uri",
        query={"branch_id": branch_id, "database_name": database, "role_name": role, "pooled": "false"},
    )
    uri = response.get("uri")
    if not isinstance(uri, str) or not uri.strip():
        raise OperatorError("Neon connection URI response did not contain uri")
    return force_verify_full(uri)


def run_checked(command: list[str], env: dict[str, str]) -> None:
    completed = subprocess.run(command, cwd=ROOT, env=env, check=False)
    if completed.returncode != 0:
        raise OperatorError(f"operator subprocess failed: {Path(command[1]).name if len(command) > 1 else command[0]}")


def fixture_env(source_uri: str, environment_class: str) -> dict[str, str]:
    env = os.environ.copy()
    env[SOURCE_URL_ENV] = source_uri
    env[ENV_CLASS_ENV] = environment_class
    return env


def run_seed(source_uri: str, environment_class: str, artifact: Path) -> dict[str, Any]:
    run_checked(
        [
            sys.executable,
            "tools/ci/host012_managed_recovery_fixture.py",
            "seed",
            "--artifact",
            str(artifact),
            "--confirm-mutate-source",
        ],
        fixture_env(source_uri, environment_class),
    )
    return load_json(artifact)


def run_delete(source_uri: str, environment_class: str, seed: Path, artifact: Path) -> dict[str, Any]:
    run_checked(
        [
            sys.executable,
            "tools/ci/host012_managed_recovery_fixture.py",
            "delete",
            "--artifact",
            str(artifact),
            "--seed-artifact",
            str(seed),
            "--confirm-mutate-source",
        ],
        fixture_env(source_uri, environment_class),
    )
    return load_json(artifact)


def create_pitr_branch(token: str, project_id: str, source_branch_id: str, restore_point: str, name: str) -> dict[str, Any]:
    response = api_request(
        token,
        "POST",
        f"/projects/{project_id}/branches",
        body={"branch": {"name": name, "parent_id": source_branch_id, "parent_timestamp": restore_point}},
    )
    branch = response.get("branch")
    if not isinstance(branch, dict) or not branch.get("id"):
        raise OperatorError("Neon PITR branch creation returned no branch id")
    return branch


def is_provider_conflict(exc: OperatorError) -> bool:
    return str(exc).startswith("Neon API HTTP 423")


def wait_for_branch_ready(
    token: str, project_id: str, branch_id: str, *, timeout_seconds: int = 180
) -> dict[str, Any]:
    deadline = time.monotonic() + timeout_seconds
    last_state = "unknown"
    while time.monotonic() < deadline:
        try:
            response = api_request(token, "GET", f"/projects/{project_id}/branches/{branch_id}")
            branch = response.get("branch")
            if isinstance(branch, dict):
                last_state = str(branch.get("current_state") or "unknown")
                if last_state == "ready":
                    return branch
        except OperatorError as exc:
            if not is_provider_conflict(exc):
                raise
            last_state = str(exc)
        time.sleep(3)
    raise OperatorError(f"Neon PITR branch did not become ready: {last_state}")


def create_compute(
    token: str, project_id: str, branch_id: str, *, timeout_seconds: int = 180
) -> dict[str, Any]:
    deadline = time.monotonic() + timeout_seconds
    last_error = ""
    while time.monotonic() < deadline:
        try:
            response = api_request(
                token,
                "POST",
                f"/projects/{project_id}/endpoints",
                body={"endpoint": {"branch_id": branch_id, "type": "read_write"}},
            )
            endpoint = response.get("endpoint")
            if not isinstance(endpoint, dict) or not endpoint.get("id"):
                raise OperatorError("Neon restored branch compute creation returned no endpoint id")
            return endpoint
        except OperatorError as exc:
            if not is_provider_conflict(exc):
                raise
            last_error = str(exc)
            time.sleep(3)
    raise OperatorError(f"Neon restored branch compute did not become schedulable: {last_error}")


def wait_for_connection_uri(
    token: str, project_id: str, branch_id: str, database: str, role: str, *, timeout_seconds: int = 180
) -> str:
    deadline = time.monotonic() + timeout_seconds
    last_error = ""
    while time.monotonic() < deadline:
        try:
            return connection_uri(token, project_id, branch_id, database, role)
        except OperatorError as exc:
            last_error = str(exc)
            time.sleep(3)
    raise OperatorError(f"restored Neon connection URI did not become ready: {last_error}")


def run_privacy_drill(
    source_uri: str,
    restored_uri: str,
    environment_class: str,
    user_id: str,
    delete_artifact: Path,
    privacy_artifact: Path,
) -> dict[str, Any]:
    env = os.environ.copy()
    env[SOURCE_URL_ENV] = source_uri
    env[RESTORED_URL_ENV] = restored_uri
    env[USER_ENV] = user_id
    env[ENV_CLASS_ENV] = environment_class
    run_checked(
        [
            sys.executable,
            "tools/ci/host012_managed_recovery_drill.py",
            "--artifact",
            str(privacy_artifact),
            "--fixture-evidence",
            str(delete_artifact),
            "--confirm-replace-restored-target",
        ],
        env,
    )
    return load_json(privacy_artifact)


def build_backup_metadata(
    *, project: dict[str, Any], project_id: str, source_branch_id: str, restore_point: str, candidate_sha: str
) -> dict[str, Any]:
    retention_seconds = float(project.get("history_retention_seconds") or 0)
    if retention_seconds <= 0:
        raise OperatorError("Neon project does not expose a positive history retention window")
    return {
        "schema": "DE.PULSE-HOST012-NEON-BACKUP-METADATA-1",
        "provider": PROVIDER,
        "projectId": project_id,
        "sourceBranchId": source_branch_id,
        "candidateSha": candidate_sha,
        "requestedRestorePoint": restore_point,
        "historyRetentionSeconds": retention_seconds,
        "recoveryWindowHours": retention_seconds / 3600.0,
        "retentionDays": retention_seconds / 86400.0,
        "encrypted": True,
        "pitrEnabled": True,
        "encryptionBasis": ENCRYPTION_BASIS,
        "backupId": f"neon-history:{project_id}:{source_branch_id}:{restore_point}",
        "capturedAt": iso_utc(utc_now()),
        "containsSecrets": False,
    }


def build_restore_metadata(
    *,
    project_id: str,
    source_branch_id: str,
    restored_branch: dict[str, Any],
    endpoint: dict[str, Any],
    restore_point: str,
    restore_started: datetime,
    restore_completed: datetime,
    candidate_sha: str,
) -> dict[str, Any]:
    return {
        "schema": "DE.PULSE-HOST012-NEON-RESTORE-METADATA-1",
        "provider": PROVIDER,
        "projectId": project_id,
        "sourceBranchId": source_branch_id,
        "restoredBranchId": str(restored_branch.get("id") or ""),
        "restoredBranchName": str(restored_branch.get("name") or ""),
        "endpointId": str(endpoint.get("id") or ""),
        "requestedRestorePoint": restore_point,
        "restoreStartedAt": iso_utc(restore_started),
        "restoreCompletedAt": iso_utc(restore_completed),
        "candidateSha": candidate_sha,
        "pointInTimeBranching": True,
        "containsSecrets": False,
    }


def build_final_evidence(
    *,
    candidate_sha: str,
    environment_class: str,
    project: dict[str, Any],
    project_id: str,
    source_branch_id: str,
    restored_branch_id: str,
    restore_point: str,
    restore_started: datetime,
    restore_completed: datetime,
    backup_metadata: Path,
    restore_metadata: Path,
    privacy_artifact: Path,
    privacy: dict[str, Any],
    delete_fixture: dict[str, Any],
    rpo_target_seconds: float,
    rto_target_seconds: float,
) -> dict[str, Any]:
    deletion_at = datetime.fromtimestamp(float(delete_fixture["tombstoneDeletedAt"]) / 1000.0, tz=timezone.utc)
    restore_point_dt = parse_iso(restore_point)
    measured_rpo = max(0.0, (deletion_at - restore_point_dt).total_seconds())
    measured_rto = max(0.0, (restore_completed - restore_started).total_seconds())
    retention_seconds = float(project.get("history_retention_seconds") or 0)
    privacy_assertions = privacy.get("privacyAssertions")
    if not isinstance(privacy_assertions, dict):
        raise OperatorError("privacy artifact lacks privacyAssertions")
    return {
        "schema": EVIDENCE_SCHEMA,
        "evidenceMode": "MANAGED_HOSTED_OPERATOR_DRILL",
        "requirements": ["HOST-012", "HOST-016"],
        "candidateSha": candidate_sha,
        "environment": {
            "managed": True,
            "provider": PROVIDER,
            "databaseService": DATABASE_SERVICE,
            "region": str(project.get("region_id") or "unknown"),
            "class": environment_class,
            "sourceDatabaseId": source_branch_id,
            "restoredDatabaseId": restored_branch_id,
        },
        "backup": {
            "encrypted": True,
            "pitrEnabled": True,
            "backupId": f"neon-history:{project_id}:{source_branch_id}:{restore_point}",
            "restoreId": restored_branch_id,
            "recoveryWindowHours": retention_seconds / 3600.0,
            "retentionDays": retention_seconds / 86400.0,
            "requestedRestorePoint": restore_point,
            "restoreCompletedAt": iso_utc(restore_completed),
        },
        "operator": {
            "drillId": f"host012-neon-{candidate_sha[:12]}-{int(restore_started.timestamp())}",
            "role": "DE.PULSE development managed-recovery operator",
            "changeReference": f"git:{candidate_sha}",
            "executedAt": iso_utc(restore_started),
        },
        "recovery": {
            "serviceBlockedDuringRecovery": True,
            "restoredFromPointBeforeDeletion": bool((privacy.get("beforeReplay") or {}).get("liveUserPresent")),
            "authoritativeTombstoneSource": f"current-source:{source_branch_id}:Application.deleteAccountData",
            "tombstoneReplayCompleted": bool((privacy.get("canonicalReplay") or {}).get("restartVerified")),
            "serviceEnabledAfterVerification": True,
            "serviceEnablementScope": "isolated-development-recovery-target; no end-user routing changed",
            "rpoTargetSeconds": rpo_target_seconds,
            "measuredRPOSeconds": measured_rpo,
            "rtoTargetSeconds": rto_target_seconds,
            "measuredRTOSeconds": measured_rto,
        },
        "privacyAssertions": {
            "deletedUsersResurrected": privacy_assertions.get("deletedUsersResurrected"),
            "liveDeletedProfiles": privacy_assertions.get("liveDeletedProfiles"),
            "personalWorkspaceRowsForDeletedUsers": privacy_assertions.get("personalWorkspaceRowsForDeletedUsers"),
            "activeSessionsForDeletedUsers": privacy_assertions.get("activeSessionsForDeletedUsers"),
            "activeDevicesForDeletedUsers": privacy_assertions.get("activeDevicesForDeletedUsers"),
            "tombstonesPresentAfterReplay": privacy_assertions.get("tombstonesPresentAfterReplay"),
            "canonicalHealthCheckPassed": privacy_assertions.get("canonicalHealthCheckPassed"),
            "accountDataProjectionChecked": privacy_assertions.get("accountDataProjectionChecked"),
        },
        "artifacts": [
            {"kind": "backup-metadata", "id": backup_metadata.name, "sha256": file_sha256(backup_metadata)},
            {"kind": "restore-metadata", "id": restore_metadata.name, "sha256": file_sha256(restore_metadata)},
            {"kind": "post-restore-privacy-verification", "id": privacy_artifact.name, "sha256": file_sha256(privacy_artifact)},
        ],
    }


def delete_branch(token: str, project_id: str, branch_id: str) -> None:
    api_request(token, "DELETE", f"/projects/{project_id}/branches/{branch_id}")


def main(argv: list[str] | None = None) -> int:
    parser = argparse.ArgumentParser(description="Run the real HOST-012 recovery drill against Neon")
    parser.add_argument("--project-id", required=True)
    parser.add_argument("--source-branch-id", required=True)
    parser.add_argument("--database-name")
    parser.add_argument("--role-name")
    parser.add_argument("--environment-class", default="stage", choices=("development", "test", "stage"))
    parser.add_argument("--output-dir", default=".depulse-host012-neon")
    parser.add_argument("--rpo-target-seconds", type=float, default=300.0)
    parser.add_argument("--rto-target-seconds", type=float, default=900.0)
    parser.add_argument("--confirm-source-mutation", action="store_true")
    parser.add_argument("--confirm-pitr-restore", action="store_true")
    parser.add_argument("--delete-restored-branch-after-validation", action="store_true")
    args = parser.parse_args(argv)

    if not args.confirm_source_mutation or not args.confirm_pitr_restore:
        print("HOST-012 Neon recovery operator: FAIL", file=sys.stderr)
        print("both --confirm-source-mutation and --confirm-pitr-restore are required", file=sys.stderr)
        return 1
    if args.rpo_target_seconds <= 0 or args.rto_target_seconds <= 0:
        print("HOST-012 Neon recovery operator: FAIL", file=sys.stderr)
        print("RPO/RTO targets must be positive", file=sys.stderr)
        return 1
    token = os.environ.get(API_KEY_ENV, "").strip()
    if not token:
        print("HOST-012 Neon recovery operator: FAIL", file=sys.stderr)
        print(f"{API_KEY_ENV} is required and must be supplied as a secret environment variable", file=sys.stderr)
        return 1

    out = Path(args.output_dir).expanduser().resolve()
    out.mkdir(parents=True, exist_ok=True)
    seed_path = out / "seed-fixture.json"
    delete_path = out / "delete-fixture.json"
    backup_path = out / "backup-metadata.json"
    restore_path = out / "restore-metadata.json"
    privacy_path = out / "post-restore-privacy-verification.json"
    evidence_path = out / "managed-recovery-evidence.json"
    validation_path = out / "managed-recovery-validation.json"

    restored_branch_id = ""
    try:
        candidate_sha = exact_head_sha()
        project_response = api_request(token, "GET", f"/projects/{args.project_id}")
        project = project_response.get("project")
        if not isinstance(project, dict):
            raise OperatorError("Neon project lookup returned no project")
        source_response = api_request(token, "GET", f"/projects/{args.project_id}/branches/{args.source_branch_id}")
        source_branch = source_response.get("branch")
        if not isinstance(source_branch, dict) or str(source_branch.get("id") or "") != args.source_branch_id:
            raise OperatorError("Neon source branch lookup did not match requested branch")
        if str(source_branch.get("project_id") or args.project_id) != args.project_id:
            raise OperatorError("Neon source branch project mismatch")

        database = select_database(token, args.project_id, args.source_branch_id, args.database_name)
        role = select_role(token, args.project_id, args.source_branch_id, args.role_name)
        source_uri = connection_uri(token, args.project_id, args.source_branch_id, database, role)

        seed = run_seed(source_uri, args.environment_class, seed_path)
        user_id = str(seed.get("userId") or "").strip()
        if not user_id:
            raise OperatorError("seed fixture did not return a recovery user id")

        time.sleep(2)
        restore_point_dt = utc_now()
        restore_point = iso_utc(restore_point_dt)
        write_json(
            backup_path,
            build_backup_metadata(
                project=project,
                project_id=args.project_id,
                source_branch_id=args.source_branch_id,
                restore_point=restore_point,
                candidate_sha=candidate_sha,
            ),
        )

        time.sleep(2)
        deleted = run_delete(source_uri, args.environment_class, seed_path, delete_path)
        if str(deleted.get("userId") or "") != user_id:
            raise OperatorError("delete fixture user does not match seed fixture")
        if parse_iso(restore_point) >= datetime.fromtimestamp(float(deleted["tombstoneDeletedAt"]) / 1000.0, tz=timezone.utc):
            raise OperatorError("provider restore point is not before the canonical account deletion")

        restore_started = utc_now()
        branch_name = f"host012-pitr-{candidate_sha[:8]}-{int(restore_started.timestamp())}"
        restored_branch = create_pitr_branch(
            token, args.project_id, args.source_branch_id, restore_point, branch_name
        )
        restored_branch_id = str(restored_branch["id"])
        if restored_branch_id == args.source_branch_id:
            raise OperatorError("Neon PITR target must be a distinct branch")
        restored_branch = wait_for_branch_ready(token, args.project_id, restored_branch_id)
        endpoint = create_compute(token, args.project_id, restored_branch_id)
        restored_uri = wait_for_connection_uri(token, args.project_id, restored_branch_id, database, role)

        privacy = run_privacy_drill(
            source_uri, restored_uri, args.environment_class, user_id, delete_path, privacy_path
        )
        restore_completed = utc_now()
        write_json(
            restore_path,
            build_restore_metadata(
                project_id=args.project_id,
                source_branch_id=args.source_branch_id,
                restored_branch=restored_branch,
                endpoint=endpoint,
                restore_point=restore_point,
                restore_started=restore_started,
                restore_completed=restore_completed,
                candidate_sha=candidate_sha,
            ),
        )
        evidence = build_final_evidence(
            candidate_sha=candidate_sha,
            environment_class=args.environment_class,
            project=project,
            project_id=args.project_id,
            source_branch_id=args.source_branch_id,
            restored_branch_id=restored_branch_id,
            restore_point=restore_point,
            restore_started=restore_started,
            restore_completed=restore_completed,
            backup_metadata=backup_path,
            restore_metadata=restore_path,
            privacy_artifact=privacy_path,
            privacy=privacy,
            delete_fixture=deleted,
            rpo_target_seconds=args.rpo_target_seconds,
            rto_target_seconds=args.rto_target_seconds,
        )
        write_json(evidence_path, evidence)

        status = canonical_drill.run_final_validation(
            evidence_path,
            backup_path,
            restore_path,
            privacy_path,
            candidate_sha,
            validation_path,
        )
        if status != 0:
            raise OperatorError("final managed recovery evidence gate failed")

        result = {
            "status": "PASS",
            "candidateSha": candidate_sha,
            "provider": PROVIDER,
            "projectId": args.project_id,
            "sourceBranchId": args.source_branch_id,
            "restoredBranchId": restored_branch_id,
            "userIdSha256": str(seed.get("userIdSha256") or ""),
            "restorePoint": restore_point,
            "evidenceSha256": file_sha256(evidence_path),
            "validationSha256": file_sha256(validation_path),
            "containsConnectionUris": False,
            "restoredBranchDeletedAfterValidation": bool(args.delete_restored_branch_after_validation),
        }
        if args.delete_restored_branch_after_validation:
            delete_branch(token, args.project_id, restored_branch_id)
            restored_branch_id = ""
        print(json.dumps(result, indent=2, sort_keys=True))
        return 0
    except (OperatorError, OSError, ValueError, KeyError, subprocess.CalledProcessError, json.JSONDecodeError) as exc:
        print("HOST-012 Neon recovery operator: FAIL", file=sys.stderr)
        print(str(exc), file=sys.stderr)
        if restored_branch_id:
            print(f"recovery target retained for inspection: {restored_branch_id}", file=sys.stderr)
        return 1


if __name__ == "__main__":
    raise SystemExit(main())
