#!/usr/bin/env python3
from __future__ import annotations

import json
from pathlib import Path
import subprocess
import tempfile

ROOT = Path(__file__).resolve().parents[2]
G16 = ROOT / "tools" / "release" / "g16_delivery_evidence.py"
ROOT_GATE = ROOT / "tools" / "ci" / "root_ownership_gate.py"


def main() -> int:
    with tempfile.TemporaryDirectory() as d:
        tmp = Path(d)
        telemetry = tmp / "telemetry.json"
        ownership = tmp / "root-ownership.json"
        out = tmp / "g16.json"
        telemetry.write_text(json.dumps({
            "runnerMinutesByPlatform": {"linux": 1.0, "macos": 2.0, "windows": 3.0, "unknown": 0.0},
            "trustworthyEvidenceAccounting": {
                "successfulEvidenceUnits": 4,
                "avoidedWork": {
                    "impactRoutedJobRunsSkipped": 2,
                    "skippedJobs": ["A", "B"],
                    "browserSetupOperationsSkipped": 1,
                    "runnerMinutesAvoided": None,
                    "runnerMinutesAvoidedReason": "counterfactual not observed",
                },
            },
        }), encoding="utf-8")
        subprocess.run([
            "python3", str(ROOT_GATE), "--json-out", str(ownership)
        ], cwd=ROOT, check=True, capture_output=True, text=True)
        subprocess.run([
            "python3", str(G16),
            "--telemetry-json", str(telemetry),
            "--root-ownership-json", str(ownership),
            "--version", "18.9.2",
            "--stable-tag", "v18.9.2",
            "--source-pr", "999",
            "--qualified-source-head", "a" * 40,
            "--candidate-sha", "b" * 40,
            "--source-fingerprint", "c" * 64,
            "--build-id", "test-build",
            "--release-run-id", "123",
            "--run-attempt", "2",
            "--publication-state", "EXACT_ALREADY_PUBLISHED",
            "--output", str(out),
        ], cwd=ROOT, check=True, capture_output=True, text=True)
        data = json.loads(out.read_text(encoding="utf-8"))
        assert data["schema"] == "DE.PULSE-G16-CI-DELIVERY-6"
        assert data["repositoryClosure"]["baselineRootFiles"] >= data["repositoryClosure"]["currentRootFiles"]
        assert data["repositoryClosure"]["registeredMoveRenameConsolidationCount"] > 0
        assert data["repositoryClosure"]["remainingIntentionalRootOwners"]
        assert data["ciEfficiency"]["runnerMinutesConsumed"] == 6.0
        assert data["ciEfficiency"]["workflowRunAttempt"] == 2
        assert data["ciEfficiency"]["rerunsObserved"] == 1
        assert data["ciEfficiency"]["rerunsAvoided"] is None
        assert data["evidenceReuse"] is True
        assert data["qualityRequirementsRemoved"] is False

    print("DE.PULSE G16 delivery evidence self-test: PASS")
    print("root before/after + migration path accounting: PASS")
    print("remaining intentional root-owner accounting: PASS")
    print("measured runner-minute accounting: PASS")
    print("observed rerun + no fabricated avoided-rerun contract: PASS")
    print("evidence reuse + no-quality-removal contract: PASS")
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
