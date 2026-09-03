#!/usr/bin/env python3
"""Zero-network contract tests for the HOST-012 Neon recovery operator."""
from __future__ import annotations

from datetime import datetime, timedelta, timezone
import json
from pathlib import Path
import tempfile

import host012_managed_recovery_evidence_gate as evidence_gate
import host012_neon_recovery_operator as operator


def require(condition: bool, message: str) -> None:
    if not condition:
        raise AssertionError(message)


def test_force_verify_full() -> None:
    raw = "postgresql://user:password@example.neon.tech/neondb?sslmode=require&channel_binding=require"
    updated = operator.force_verify_full(raw)
    require("sslmode=verify-full" in updated, "sslmode must be verify-full")
    require("channel_binding=require" in updated, "channel_binding must be conserved")
    require("user:password@example.neon.tech" in updated, "connection authority must be conserved in memory")
    try:
        operator.force_verify_full("sqlite:///tmp/example.db")
    except operator.OperatorError:
        pass
    else:
        raise AssertionError("non-PostgreSQL URI must be rejected")


def test_neon_async_branch_and_compute_retry() -> None:
    original_api_request = operator.api_request
    original_sleep = operator.time.sleep
    calls = {"branch": 0, "compute": 0}

    def fake_api_request(token: str, method: str, path: str, **kwargs: object) -> dict[str, object]:
        del token, kwargs
        if method == "GET" and path.endswith("/branches/br-restored"):
            calls["branch"] += 1
            state = "init" if calls["branch"] == 1 else "ready"
            return {"branch": {"id": "br-restored", "name": "host012-pitr", "current_state": state}}
        if method == "POST" and path.endswith("/endpoints"):
            calls["compute"] += 1
            if calls["compute"] == 1:
                raise operator.OperatorError(
                    "Neon API HTTP 423: project already has running conflicting operations, scheduling of new ones is prohibited"
                )
            return {"endpoint": {"id": "ep-restored"}}
        raise AssertionError(f"unexpected zero-network API request: {method} {path}")

    try:
        operator.api_request = fake_api_request
        operator.time.sleep = lambda _: None
        branch = operator.wait_for_branch_ready("secret-not-used", "project-example", "br-restored", timeout_seconds=2)
        endpoint = operator.create_compute("secret-not-used", "project-example", "br-restored", timeout_seconds=2)
    finally:
        operator.api_request = original_api_request
        operator.time.sleep = original_sleep

    require(branch["current_state"] == "ready", "PITR branch must be observed READY before compute creation")
    require(endpoint["id"] == "ep-restored", "compute creation must succeed after transient provider conflict")
    require(calls["branch"] == 2, "branch readiness polling must retry non-ready state")
    require(calls["compute"] == 2, "HTTP 423 conflict must be retried exactly until successful in this fixture")


def test_secret_free_metadata_and_final_evidence() -> None:
    candidate = "a" * 40
    project = {
        "region_id": "aws-us-east-2",
        "history_retention_seconds": 21600,
    }
    restore_point_dt = datetime(2026, 8, 27, 22, 0, 0, tzinfo=timezone.utc)
    restore_started = restore_point_dt + timedelta(seconds=10)
    restore_completed = restore_started + timedelta(seconds=20)
    restore_point = operator.iso_utc(restore_point_dt)
    deleted = {
        "tombstoneDeletedAt": int((restore_point_dt + timedelta(seconds=5)).timestamp() * 1000),
    }
    privacy = {
        "beforeReplay": {"liveUserPresent": True, "tombstonePresent": False},
        "canonicalReplay": {"restartVerified": True},
        "privacyAssertions": {
            "deletedUsersResurrected": 0,
            "liveDeletedProfiles": 0,
            "personalWorkspaceRowsForDeletedUsers": 0,
            "activeSessionsForDeletedUsers": 0,
            "activeDevicesForDeletedUsers": 0,
            "tombstonesPresentAfterReplay": 1,
            "canonicalHealthCheckPassed": True,
            "accountDataProjectionChecked": True,
        },
    }

    with tempfile.TemporaryDirectory() as tmp:
        root = Path(tmp)
        backup = root / "backup-metadata.json"
        restore = root / "restore-metadata.json"
        privacy_path = root / "privacy.json"
        operator.write_json(
            backup,
            operator.build_backup_metadata(
                project=project,
                project_id="project-example",
                source_branch_id="br-source",
                restore_point=restore_point,
                candidate_sha=candidate,
            ),
        )
        operator.write_json(
            restore,
            operator.build_restore_metadata(
                project_id="project-example",
                source_branch_id="br-source",
                restored_branch={"id": "br-restored", "name": "host012-pitr"},
                endpoint={"id": "ep-restored"},
                restore_point=restore_point,
                restore_started=restore_started,
                restore_completed=restore_completed,
                candidate_sha=candidate,
            ),
        )
        operator.write_json(privacy_path, privacy)
        final = operator.build_final_evidence(
            candidate_sha=candidate,
            environment_class="stage",
            project=project,
            project_id="project-example",
            source_branch_id="br-source",
            restored_branch_id="br-restored",
            restore_point=restore_point,
            restore_started=restore_started,
            restore_completed=restore_completed,
            backup_metadata=backup,
            restore_metadata=restore,
            privacy_artifact=privacy_path,
            privacy=privacy,
            delete_fixture=deleted,
            rpo_target_seconds=300,
            rto_target_seconds=900,
        )
        errors = evidence_gate.validate_evidence(final, expected_sha=candidate)
        require(not errors, "final evidence builder must satisfy canonical gate: " + "; ".join(errors))

        rendered = json.dumps({"backup": json.loads(backup.read_text()), "restore": json.loads(restore.read_text()), "final": final})
        normalized = rendered.lower()
        for forbidden in (
            "postgresql://",
            "postgres://",
            "connectionstring",
            "connection_uri",
            '"password"',
            '"apikey"',
            '"api_key"',
            '"token"',
        ):
            require(forbidden not in normalized, f"sanitized evidence leaked forbidden material: {forbidden}")

        kinds = {item["kind"] for item in final["artifacts"]}
        require(
            kinds == {"backup-metadata", "restore-metadata", "post-restore-privacy-verification"},
            "final evidence must bind exactly the three canonical recovery artifact kinds",
        )
        require(final["environment"]["sourceDatabaseId"] != final["environment"]["restoredDatabaseId"], "restore target must be distinct")
        require(final["recovery"]["measuredRPOSeconds"] == 5.0, "RPO must be measured from restore point to deletion")
        require(final["recovery"]["measuredRTOSeconds"] == 20.0, "RTO must be measured across provider restore plus verification window")


def main() -> int:
    test_force_verify_full()
    test_neon_async_branch_and_compute_retry()
    test_secret_free_metadata_and_final_evidence()
    print("HOST-012 Neon managed-recovery operator self-test: PASS")
    print("provider API network calls: NONE")
    print("transient Neon HTTP 423 async conflicts: RETRIED")
    print("connection URI evidence persistence: PROHIBITED")
    print("canonical final evidence gate composition: PASS")
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
