#!/usr/bin/env python3
"""Fail-closed T8 performance/load/soak/concurrency/resource assurance."""
from __future__ import annotations

import argparse
import json
from pathlib import Path
import sys

ROOT = Path(__file__).resolve().parents[2]
PROGRAM = ROOT / "governance" / "programs" / "ADAPT-V18-FINAL-CLOSURE-10-10-001"
T8 = PROGRAM / "T8_PERFORMANCE_LOAD_SOAK_CONCURRENCY_RESOURCE_ASSURANCE.json"
CURRENT = ROOT / "governance" / "current-state.json"
CLOSURE = ROOT / "governance" / "work-slices" / "ADAPT-V18-FINAL-CLOSURE-10-10-001" / "closure.json"
WORKLOAD = ROOT / "workload_backpressure_regression_test.go"
ACTIVE = ROOT / "active_market_reliability_regression_test.go"
SLO = ROOT / "runtime_slo.go"
CI_FAST = ROOT / ".github" / "workflows" / "ci-fast.yml"
CI_QUALIFIED = ROOT / ".github" / "workflows" / "ci-qualified.yml"
T5 = PROGRAM / "T5_PERSISTENCE_LIFECYCLE_ASSURANCE.json"

REQUIRED_CONTROL_AREAS = {
    "provider pressure and protected capacity",
    "scanner and scheduled/background job pressure",
    "persistence pressure and lifecycle recovery",
    "bounded queue and permit convergence",
    "duplicate upstream work coalescing",
    "goroutine and concurrency safety",
    "CPU, memory and GC resource safety",
    "lock/contention risk",
    "UI/API latency budgets",
    "race detector",
    "randomized execution order",
    "bounded soak/repeated-run recovery",
}
REQUIRED_RULES = {
    "reuseCanonicalRuntimeOwners",
    "noParallelPerformanceOrBackpressureSubsystem",
    "deterministicBoundedCiEvidencePreferredOverFragileWallClockBenchmarks",
    "raceAndRandomizedExecutionRemainRequired",
    "providerPressureMustPreserveProtectedCriticalCapacity",
    "queuesAndPermitsMustConvergeAfterCancellationAndRecovery",
    "duplicateProviderWorkMustRemainCoalesced",
    "persistencePressureMustPreserveIntegrityAndRecovery",
    "resourceEvidenceMustCoverGoroutineCpuMemoryGcAndLockRiskWhereMaterial",
    "latencyEvidenceMustUseRepositoryOwnedMeaningfulBudgets",
    "t9OrT10CertificationNotImplied",
}


def load(path: Path) -> dict:
    with path.open("r", encoding="utf-8") as handle:
        return json.load(handle)


def contains(path: Path, *needles: str) -> bool:
    text = path.read_text(encoding="utf-8", errors="ignore")
    return all(needle in text for needle in needles)


def closure_row(closure: dict, row_id: str) -> dict | None:
    return next((row for row in closure.get("gaps", []) if isinstance(row, dict) and row.get("id") == row_id), None)


def governed_status(product: dict, track: str) -> str:
    return next((str(row.get("status") or "") for row in product.get("nextGovernedTracks", []) if isinstance(row, dict) and row.get("track") == track), "")


def main() -> int:
    parser = argparse.ArgumentParser()
    parser.add_argument("--strict", action="store_true")
    args = parser.parse_args()
    errors: list[str] = []

    for path in (T8, CURRENT, CLOSURE, WORKLOAD, ACTIVE, SLO, CI_FAST, CI_QUALIFIED, T5):
        if not path.is_file():
            errors.append(f"required T8 owner missing: {path.relative_to(ROOT)}")
    if errors:
        for error in errors:
            print(" - " + error, file=sys.stderr)
        return 1

    t8 = load(T8)
    current = load(CURRENT)
    closure = load(CLOSURE)
    t5 = load(T5)
    product = current.get("productCapabilityGate") or {}
    state = str(t8.get("state") or "")

    if t8.get("schema") != "DE.PULSE-V18-T8-PERFORMANCE-LOAD-SOAK-CONCURRENCY-RESOURCE-ASSURANCE-1" or t8.get("programIssue") != 113 or t8.get("trackIssue") != 121:
        errors.append("T8 assurance identity/schema mismatch")
    if set(t8.get("requiredControlAreas") or []) != REQUIRED_CONTROL_AREAS:
        errors.append("T8 required control-area matrix drifted")
    rules = t8.get("rules") or {}
    for key in sorted(REQUIRED_RULES):
        if rules.get(key) is not True:
            errors.append(f"T8 rule must remain true: {key}")

    for prior in range(1, 8):
        row = closure_row(closure, f"T{prior}-" + {
            1: "FEATURE-TRACEABILITY",
            2: "UNIT-CONTRACT-PROPERTY",
            3: "FUNCTIONAL-E2E",
            4: "EDGE-FAILURE-DATA-TRUTH",
            5: "PERSISTENCE-LIFECYCLE-RECOVERY",
            6: "SECURITY-ROLES-RIGHTS",
            7: "UI-UX-IA-CONTENT",
        }[prior])
        if not isinstance(row, dict) or row.get("status") != "VERIFIED":
            errors.append(f"T8 requires T{prior} durable VERIFIED closure")

    if t5.get("state") != "COMPLETE" or t5.get("uncoveredResponsibilityCount") != 0:
        errors.append("T8 persistence pressure must inherit a zero-gap COMPLETE T5 lifecycle owner")

    workload_text = WORKLOAD.read_text(encoding="utf-8", errors="ignore")
    for marker in (
        "TestT8RepeatedWorkloadPressureConvergesWithoutLeakedPermitsOrQueue",
        "baseline.MaxQueue", "baseline.Capacity", "expectedCanceled", "expectedCompleted",
        "did not recover admission after pressure",
    ):
        if marker not in workload_text:
            errors.append(f"canonical workload owner lost T8 pressure evidence: {marker}")

    active_text = ACTIVE.read_text(encoding="utf-8", errors="ignore")
    for marker in (
        "TestV1870ActiveMarketDuplicateSnapshotBurstCoalesces",
        "calls.Load() != 1",
        "TestV1870ActiveMarketProviderPressureIsBoundedAndTruthful",
        "PressureState != \"PROTECTED\"",
        "TestT8CanonicalRuntimeSLOBudgetsFailClosedAtOwnedThresholds",
        '"Interactive API p95"', '"Provider queue"', '"Persistence queue age"',
        '"DB write rate"', '"CPU utilization"', '"Local process budget"',
        '"Startup/warm-start time"', '"Storage growth"',
    ):
        if marker not in active_text:
            errors.append(f"active-market/T8 canonical evidence missing marker: {marker}")

    slo_text = SLO.read_text(encoding="utf-8", errors="ignore")
    for budget in (
        'add("Interactive API p95"', 'p95 > 250', 'p95 > 1000',
        'add("Provider queue"', 'providerOldest > 2000', 'providerOldest > 5000',
        'add("Persistence queue age"', 'persistAge > 2000', 'persistAge > 10000',
        'add("DB write rate"', 'RowsWrittenLastMin > 600', 'RowsWrittenLastMin > 3000',
        'add("CPU utilization"', 'CPUUtilizationPct > 70', 'CPUUtilizationPct > 90',
        'add("Local process budget"', 'Goroutines > 600', 'HeapAllocBytes > 768*1024*1024',
        'add("Startup/warm-start time"', 'BootstrapDurationMs > 1000', 'BootstrapDurationMs > 5000',
        'add("Provider request budgets"', 'maxBudget >= 80', 'maxBudget >= 100',
        'add("Storage growth"', 'StorageGrowthBytes > 250*1024*1024', 'StorageGrowthBytes > 1024*1024*1024',
    ):
        if budget not in slo_text:
            errors.append(f"canonical runtime SLO budget missing/drifted: {budget}")

    qualified = CI_QUALIFIED.read_text(encoding="utf-8", errors="ignore")
    for marker in (
        "go test -race -count=1 ./...",
        "go test -shuffle=on -count=1 ./...",
        "Persistence-focused regression",
        "Test.*(SQLite|Persistence|Migration|Canonical|Cache)",
    ):
        if marker not in qualified:
            errors.append(f"Qualified CI lost T8 evidence lane: {marker}")

    fast = CI_FAST.read_text(encoding="utf-8", errors="ignore")
    if "python3 tools/ci/v18_t8_performance_load_soak_concurrency_resource_assurance_gate.py" not in fast:
        errors.append("T8 assurance gate is not bound into canonical CI Fast")

    resolved = {str(row.get("id")): row for row in t8.get("resolvedCoverageFindings", []) if isinstance(row, dict)}
    for finding in ("T8-PERSISTENCE-PRESSURE", "T8-PROCESS-RESOURCE-SAFETY", "T8-LATENCY-BUDGETS"):
        row = resolved.get(finding)
        if not isinstance(row, dict) or row.get("status") != "RESOLVED" or not row.get("evidence"):
            errors.append(f"T8 material finding must be RESOLVED with evidence: {finding}")

    platform = t8.get("platformScope") or {}
    if platform.get("linux") != "CI_TEST_ONLY" or platform.get("requiredPackagedReleaseTargets") != ["macOS Apple Silicon", "Windows x64"]:
        errors.append("T8 must not promote Linux into a packaged release target or weaken T9 platform scope")

    gaps = [row for row in t8.get("knownCoverageGaps", []) if isinstance(row, dict)]
    gap_ids = {str(row.get("id")) for row in gaps if row.get("status") == "OPEN"}
    t8_closure = closure_row(closure, "T8-PERFORMANCE-CONCURRENCY-SOAK")
    if state == "IN_PROGRESS":
        if product.get("nextChildIssue") != 121 or product.get("nextChildTrack") != "T8" or governed_status(product, "T8") != "IN_PROGRESS":
            errors.append("IN_PROGRESS T8 must be the active T8/#121 child in current-state")
        if governed_status(product, "T9") != "NOT_STARTED":
            errors.append("IN_PROGRESS T8 must not silently start T9")
        if not isinstance(t8_closure, dict) or t8_closure.get("status") != "IMPLEMENTED_UNVERIFIED":
            errors.append("IN_PROGRESS T8 closure row must be IMPLEMENTED_UNVERIFIED")
        if gap_ids != {"T8-EXACT-HEAD-QUALIFICATION"}:
            errors.append("after executable implementation, only exact-head qualification may remain open in T8")
    elif state == "COMPLETE":
        completed = product.get("completedChildTracks") or []
        row = next((x for x in completed if isinstance(x, dict) and x.get("track") == "T8" and x.get("issue") == 121), None)
        if not isinstance(row, dict) or row.get("status") != "COMPLETE" or not str(row.get("frozenHeadSha") or ""):
            errors.append("COMPLETE T8 requires durable completedChildTracks exact-head evidence")
        if not isinstance(t8_closure, dict) or t8_closure.get("status") != "VERIFIED":
            errors.append("COMPLETE T8 closure row must be VERIFIED")
        if gaps:
            errors.append("COMPLETE T8 cannot retain known coverage gaps")
    else:
        errors.append(f"unsupported T8 state: {state!r}")

    strict = args.strict or state == "COMPLETE"
    if strict and gaps:
        errors.append("strict T8 closure cannot retain coverage gaps")

    print("V18 T8 PERFORMANCE / LOAD / SOAK / CONCURRENCY / RESOURCE ASSURANCE")
    print(f"state: {state}")
    print(f"required control areas: {len(REQUIRED_CONTROL_AREAS)}")
    print(f"resolved material findings: {len(resolved)}")
    print(f"open gaps: {len(gap_ids)}")
    print("canonical SLO budgets: API/provider/persistence/CPU/goroutine/heap/startup/provider-rate/storage bound")
    print("race detector + randomized order + persistence integration: bound to Qualified")
    print("Linux: CI_TEST_ONLY; packaged release scope remains macOS Apple Silicon + Windows x64")
    print("T9/T10 certification is not implied by T8")

    if errors:
        print("V18 T8 ASSURANCE GATE: FAIL", file=sys.stderr)
        for error in errors:
            print(" - " + error, file=sys.stderr)
        return 1
    print(f"V18 T8 ASSURANCE GATE: PASS (strict={strict})")
    return 0


if __name__ == "__main__":
    sys.exit(main())
