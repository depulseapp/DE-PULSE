#!/usr/bin/env python3
from __future__ import annotations

import argparse
import json
import os
from pathlib import Path
import re
import shutil
import socket
import subprocess
import sys
import time


EVIDENCE_SCHEMA = "DE.PULSE-HOST016-POSTGRES-HA-EVIDENCE-1"
REPLICATION_ROLE = "depulse_host016_replica"
REPLICATION_PASSWORD = "depulse-host016-ephemeral-replication-only"


class DrillFailure(RuntimeError):
    pass


def run(
    *args: str,
    env: dict[str, str] | None = None,
    check: bool = True,
) -> subprocess.CompletedProcess[str]:
    result = subprocess.run(
        args,
        text=True,
        stdout=subprocess.PIPE,
        stderr=subprocess.PIPE,
        env=env,
        check=False,
    )
    if check and result.returncode != 0:
        detail = (result.stderr or result.stdout or "command failed").strip()
        raise DrillFailure(f"{Path(args[0]).name} failed: {detail[:1000]}")
    return result


def docker(*args: str, check: bool = True) -> subprocess.CompletedProcess[str]:
    return run("docker", *args, check=check)


def discover_primary_container() -> str:
    result = docker("ps", "--format", "{{.ID}}\t{{.Ports}}")
    candidates = []
    for line in result.stdout.splitlines():
        container_id, _, ports = line.partition("\t")
        if re.search(r"(?:0\.0\.0\.0|\[::\]):5432->5432/tcp", ports):
            candidates.append(container_id)
    if len(candidates) != 1:
        raise DrillFailure(f"expected one PostgreSQL service on host port 5432, found {len(candidates)}")
    return candidates[0]


def container_network_and_ip(container_id: str) -> tuple[str, str]:
    payload = json.loads(docker("inspect", container_id).stdout)[0]
    networks = payload.get("NetworkSettings", {}).get("Networks", {})
    if len(networks) != 1:
        raise DrillFailure(f"expected one Actions service network, found {len(networks)}")
    network, config = next(iter(networks.items()))
    ip_address = str(config.get("IPAddress", "")).strip()
    if not network or not ip_address:
        raise DrillFailure("PostgreSQL service network identity is incomplete")
    return network, ip_address


def wait_container_ready(container: str, *, attempts: int = 60) -> None:
    for _ in range(attempts):
        result = docker(
            "exec",
            container,
            "pg_isready",
            "-U",
            "depulse_ci",
            "-d",
            "depulse_test",
            check=False,
        )
        if result.returncode == 0:
            return
        time.sleep(1)
    raise DrillFailure(f"PostgreSQL container {container} did not become ready")


def psql(container: str, sql: str) -> str:
    result = docker(
        "exec",
        container,
        "psql",
        "-v",
        "ON_ERROR_STOP=1",
        "-At",
        "-U",
        "depulse_ci",
        "-d",
        "depulse_test",
        "-c",
        sql,
    )
    return result.stdout.strip()


def wait_replication(primary: str, standby: str, *, attempts: int = 90) -> tuple[str, str]:
    primary_lsn = psql(primary, "SELECT pg_current_wal_lsn()")
    if not re.fullmatch(r"[0-9A-F]+/[0-9A-F]+", primary_lsn):
        raise DrillFailure("primary did not return a valid WAL LSN")
    for _ in range(attempts):
        replay_lsn = psql(standby, "SELECT COALESCE(pg_last_wal_replay_lsn()::text,'')")
        caught_up = psql(
            standby,
            f"SELECT COALESCE(pg_last_wal_replay_lsn() >= '{primary_lsn}'::pg_lsn, false)",
        )
        if caught_up == "t":
            return primary_lsn, replay_lsn
        time.sleep(1)
    raise DrillFailure("standby did not replay through the primary fixture LSN")


def wait_promoted(container: str, *, attempts: int = 60) -> None:
    for _ in range(attempts):
        if psql(container, "SELECT NOT pg_is_in_recovery()") == "t":
            return
        time.sleep(1)
    raise DrillFailure("standby did not become the read-write primary")


def require_outage_observed() -> None:
    try:
        connection = socket.create_connection(("127.0.0.1", 5432), timeout=1)
    except OSError:
        return
    connection.close()
    raise DrillFailure("primary fault injection did not interrupt the application endpoint")


def run_fixture(phase: str, database_url: str, candidate_sha: str) -> None:
    fixture_env = os.environ.copy()
    fixture_env.update(
        {
            "DEPULSE_TEST_POSTGRES_URL": database_url,
            "DEPULSE_HOST016_FAILOVER_PHASE": phase,
            "DEPULSE_HOST016_CANDIDATE_SHA": candidate_sha,
        }
    )
    result = run(
        "go",
        "test",
        "-tags",
        "postgres",
        "-count=1",
        "-v",
        ".",
        "-run",
        "^TestHOST016PostgresPhysicalFailoverFixture$",
        env=fixture_env,
    )
    required = f"--- PASS: TestHOST016PostgresPhysicalFailoverFixture/{phase}"
    if required not in result.stdout:
        raise DrillFailure(f"HOST-016 {phase} fixture did not execute to PASS")


def parse_args() -> argparse.Namespace:
    parser = argparse.ArgumentParser(description="Run the governed HOST-016 PostgreSQL physical failover drill")
    parser.add_argument("--candidate-sha", required=True)
    parser.add_argument("--database-url", required=True)
    parser.add_argument("--postgres-image", required=True)
    parser.add_argument("--output-dir", required=True, type=Path)
    return parser.parse_args()


def main() -> int:
    args = parse_args()
    candidate_sha = args.candidate_sha.strip().lower()
    if not re.fullmatch(r"[0-9a-f]{40}", candidate_sha):
        raise SystemExit("--candidate-sha must be a full hexadecimal Git SHA")
    if "@sha256:" not in args.postgres_image:
        raise SystemExit("--postgres-image must be digest pinned")
    if shutil.which("docker") is None or shutil.which("go") is None:
        raise SystemExit("docker and go are required")

    args.output_dir.mkdir(parents=True, exist_ok=True)
    evidence_path = args.output_dir / "host016-postgres-ha-evidence.json"
    evidence: dict[str, object] = {
        "schema": EVIDENCE_SCHEMA,
        "candidateSha": candidate_sha,
        "environmentClass": "github-actions-ephemeral",
        "postgresImage": args.postgres_image,
        "status": "FAIL",
        "topology": "one read-write primary plus one physical streaming standby",
        "fault": "stop and remove the primary after standby WAL replay convergence",
        "applicationEndpoint": "same host/port/database identity before and after standby promotion",
        "providerAutomaticFailoverClaimed": False,
        "startedAt": time.strftime("%Y-%m-%dT%H:%M:%SZ", time.gmtime()),
    }
    suffix = candidate_sha[:12]
    volume = f"depulse-host016-{suffix}"
    standby = f"depulse-host016-standby-{suffix}"
    promoted = f"depulse-host016-promoted-{suffix}"
    primary = ""
    primary_stopped_at = 0.0
    try:
        primary = discover_primary_container()
        network, primary_ip = container_network_and_ip(primary)
        expected_image_id = docker("image", "inspect", args.postgres_image, "--format", "{{.Id}}").stdout.strip()
        primary_image_id = json.loads(docker("inspect", primary).stdout)[0].get("Image", "")
        if not expected_image_id or primary_image_id != expected_image_id:
            raise DrillFailure("primary service does not use the governed digest-pinned PostgreSQL image")

        psql(
            primary,
            f"DROP ROLE IF EXISTS {REPLICATION_ROLE}; "
            f"CREATE ROLE {REPLICATION_ROLE} WITH REPLICATION LOGIN PASSWORD '{REPLICATION_PASSWORD}'",
        )
        docker(
            "exec",
            "--user",
            "postgres",
            primary,
            "sh",
            "-ec",
            f"printf '%s\\n' 'host replication {REPLICATION_ROLE} all scram-sha-256' >> \"$PGDATA/pg_hba.conf\"; pg_ctl reload -D \"$PGDATA\"",
        )

        docker("volume", "create", volume)
        docker(
            "run",
            "--rm",
            "-v",
            f"{volume}:/var/lib/postgresql/data",
            args.postgres_image,
            "sh",
            "-ec",
            "chown -R postgres:postgres /var/lib/postgresql/data",
        )
        backup_env = os.environ.copy()
        backup_env["PGPASSWORD"] = REPLICATION_PASSWORD
        run(
            "docker",
            "run",
            "--rm",
            "--user",
            "postgres",
            "--network",
            network,
            "-e",
            "PGPASSWORD",
            "-v",
            f"{volume}:/var/lib/postgresql/data",
            args.postgres_image,
            "pg_basebackup",
            "-h",
            primary_ip,
            "-p",
            "5432",
            "-U",
            REPLICATION_ROLE,
            "-D",
            "/var/lib/postgresql/data",
            "-Fp",
            "-Xs",
            "-R",
            "-P",
            env=backup_env,
        )
        docker(
            "run",
            "-d",
            "--name",
            standby,
            "--network",
            network,
            "-v",
            f"{volume}:/var/lib/postgresql/data",
            args.postgres_image,
            "postgres",
        )
        wait_container_ready(standby)
        if psql(standby, "SELECT pg_is_in_recovery()") != "t":
            raise DrillFailure("replica started without physical recovery mode")

        run_fixture("seed", args.database_url, candidate_sha)
        primary_lsn, standby_lsn = wait_replication(primary, standby)

        primary_stopped_at = time.monotonic()
        docker("stop", "--time", "2", primary)
        require_outage_observed()
        docker("stop", "--time", "2", standby)
        docker("rm", primary)
        primary = ""
        docker("rm", standby)

        docker(
            "run",
            "-d",
            "--name",
            promoted,
            "--network",
            network,
            "-p",
            "5432:5432",
            "-v",
            f"{volume}:/var/lib/postgresql/data",
            args.postgres_image,
            "postgres",
        )
        wait_container_ready(promoted)
        docker(
            "exec",
            "--user",
            "postgres",
            promoted,
            "pg_ctl",
            "promote",
            "-D",
            "/var/lib/postgresql/data",
            "-w",
            "-t",
            "30",
        )
        wait_promoted(promoted)
        rto_seconds = round(time.monotonic() - primary_stopped_at, 3)
        run_fixture("verify", args.database_url, candidate_sha)
        promoted_lsn = psql(promoted, "SELECT pg_current_wal_lsn()")

        evidence.update(
            {
                "status": "PASS",
                "primaryFixtureLsn": primary_lsn,
                "standbyReplayLsnBeforeFault": standby_lsn,
                "promotedPrimaryLsn": promoted_lsn,
                "fixtureRpo": "ZERO_CONFIRMED_THROUGH_REPLAYED_FIXTURE_LSN",
                "measuredApplicationRtoSeconds": rto_seconds,
                "assertions": [
                    "physical standby replayed through the committed fixture LSN",
                    "the original application endpoint was unavailable after primary loss",
                    "standby promotion restored a read-write endpoint at the same DSN",
                    "the canonical hosted PostgreSQL backend reconnected and reopened",
                    "two tenant identities remained independently owned",
                    "workspace ownership remained bound to canonical tenant identity",
                    "the deletion tombstone and privacy-blank workspace survived",
                    "tampered cross-tenant restore failed closed after failover",
                    "legacy aggregate authority remained retired",
                ],
            }
        )
    except Exception as exc:  # evidence must survive every failure path
        evidence["failure"] = str(exc)[:1000]
    finally:
        evidence["finishedAt"] = time.strftime("%Y-%m-%dT%H:%M:%SZ", time.gmtime())
        evidence_path.write_text(json.dumps(evidence, indent=2, sort_keys=True) + "\n", encoding="utf-8")
        for container in (promoted, standby):
            docker("rm", "-f", container, check=False)
        if primary:
            docker("start", primary, check=False)
        docker("volume", "rm", "-f", volume, check=False)

    if evidence["status"] != "PASS":
        print(f"HOST-016 PostgreSQL HA drill failed: {evidence.get('failure', 'unknown failure')}", file=sys.stderr)
        return 1
    print(f"PASS: HOST-016 physical PostgreSQL failover evidence written to {evidence_path}")
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
