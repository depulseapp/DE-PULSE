#!/usr/bin/env python3
from __future__ import annotations

import json
from pathlib import Path
import subprocess
import tempfile

ROOT = Path(__file__).resolve().parents[2]
TOOL = ROOT / "tools" / "ci" / "ci_telemetry.py"

JOBS = {
    "jobs": [
        {
            "name": "CI/harness contract",
            "status": "completed",
            "conclusion": "success",
            "created_at": "2026-08-19T20:00:00Z",
            "started_at": "2026-08-19T20:00:05Z",
            "completed_at": "2026-08-19T20:01:05Z",
            "labels": ["ubuntu-24.04"]
        },
        {
            "name": "Primary WebKit compatibility",
            "status": "completed",
            "conclusion": "success",
            "created_at": "2026-08-19T20:00:00Z",
            "started_at": "2026-08-19T20:00:10Z",
            "completed_at": "2026-08-19T20:02:10Z",
            "labels": ["macos-15"]
        },
        {
            "name": "Qualified renderer contracts",
            "status": "completed",
            "conclusion": "skipped",
            "created_at": "2026-08-19T20:00:00Z",
            "started_at": "2026-08-19T20:00:00Z",
            "completed_at": "2026-08-19T20:00:00Z",
            "labels": ["ubuntu-24.04"]
        },
        {
            "name": "Qualified evidence summary",
            "status": "in_progress",
            "conclusion": None,
            "created_at": "2026-08-19T20:02:11Z",
            "started_at": "2026-08-19T20:02:12Z",
            "completed_at": None,
            "labels": ["ubuntu-24.04"]
        }
    ]
}

RUNS_OK = {
    "workflow_runs": [
        {"name": "DE.PULSE | CI Fast", "created_at": "2026-08-19T19:59:00Z"},
        {"name": "DE.PULSE | CI Qualified", "created_at": "2026-08-19T20:00:00Z"},
        {"name": "DE.PULSE | CI Fast", "created_at": "2026-08-18T19:59:00Z"}
    ]
}

RUNS_WARN = {
    "workflow_runs": [
        *[{"name": "DE.PULSE | CI Fast", "created_at": "2026-08-19T20:00:00Z"} for _ in range(4)],
        *[{"name": "DE.PULSE | CI Qualified", "created_at": "2026-08-19T20:00:00Z"} for _ in range(3)]
    ]
}


def run_case(tmp: Path, runs: dict, name: str) -> dict:
    jobs = tmp / f"jobs-{name}.json"
    runs_path = tmp / f"runs-{name}.json"
    out = tmp / f"out-{name}.json"
    jobs.write_text(json.dumps(JOBS), encoding="utf-8")
    runs_path.write_text(json.dumps(runs), encoding="utf-8")
    subprocess.run([
        "python3", str(TOOL),
        "--jobs-json", str(jobs),
        "--runs-json", str(runs_path),
        "--workflow", "DE.PULSE | CI Qualified",
        "--run-id", "123",
        "--candidate-sha", "a" * 40,
        "--lane", "ci-harness",
        "--webkit-required", "true",
        "--pr-created-at", "2026-08-19T19:00:00Z",
        "--chrome-cache-hit", "",
        "--chrome-setup-seconds", "",
        "--webkit-cache-hit", "true",
        "--webkit-setup-seconds", "12",
        "--output", str(out),
    ], cwd=ROOT, check=True, capture_output=True, text=True)
    return json.loads(out.read_text(encoding="utf-8"))


def main() -> int:
    with tempfile.TemporaryDirectory() as d:
        tmp = Path(d)
        ok = run_case(tmp, RUNS_OK, "ok")
        assert ok["schema"] == "DE.PULSE-CI-TELEMETRY-2"
        assert ok["runnerSecondsByPlatform"] == {"linux": 60, "macos": 120, "windows": 0, "unknown": 0}, ok
        assert ok["runnerMinutesByPlatform"]["macos"] == 2.0
        assert ok["browserDependencySetup"]["webkit"] == {"pipCacheHit": True, "setupSeconds": 12}
        assert ok["browserDependencySetup"]["chrome"] == {"pipCacheHit": None, "setupSeconds": None}
        assert ok["workflowAmplification"]["counts"] == {"fast": 1, "qualified": 1, "release": 0}
        assert ok["workflowAmplification"]["status"] == "OK"

        evidence = ok["trustworthyEvidenceAccounting"]
        assert evidence["successfulEvidenceUnits"] == 2, evidence
        assert evidence["runnerMinutesConsumed"] == 3.0, evidence
        assert evidence["runnerMinutesPerTrustworthyEvidence"] == 1.5, evidence
        assert evidence["avoidedWork"]["impactRoutedJobRunsSkipped"] == 1, evidence
        assert evidence["avoidedWork"]["browserSetupOperationsSkipped"] == 0, evidence
        assert evidence["avoidedWork"]["observedBrowserSetupSeconds"] == 12, evidence
        assert evidence["avoidedWork"]["browserSetupCacheHits"] == 1, evidence
        assert evidence["avoidedWork"]["runnerMinutesAvoided"] is None
        assert evidence["avoidedWork"]["setupSecondsAvoided"] is None

        assert ok["costAccounting"]["actualCurrencyCost"] is None
        assert ok["costAccounting"]["currencyCostPerTrustworthyEvidence"] is None
        assert ok["costAccounting"]["runnerMinutesPerTrustworthyEvidence"] == 1.5

        warn = run_case(tmp, RUNS_WARN, "warn")
        assert warn["workflowAmplification"]["status"] == "WARN"
        assert len(warn["workflowAmplification"]["warnings"]) == 2

    print("DE.PULSE CI telemetry self-test: PASS")
    print("queue/runtime/platform accounting: PASS")
    print("browser setup/cache signals: PASS")
    print("workflow amplification warning thresholds: PASS")
    print("trustworthy evidence normalization: PASS")
    print("impact-routed avoided-job accounting: PASS")
    print("no fabricated avoided-minute/currency-cost contract: PASS")
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
