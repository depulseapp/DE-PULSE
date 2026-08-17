#!/usr/bin/env python3
"""DE.PULSE v18.5 Major Closure matrix gate.

Fail closed unless the machine matrix exactly covers the frozen v18.5 scope,
all 15 ADR-GDI scenarios, and the mandatory native GitHub Release delivery set.
"""
from pathlib import Path
import json
import sys

ROOT = Path(__file__).resolve().parent
SCOPE_PATH = ROOT / "v18_5_scope.json"
MATRIX_PATH = ROOT / "v18_5_major_closure.json"
errors = []


def need(condition, message):
    if not condition:
        errors.append(message)


scope = json.loads(SCOPE_PATH.read_text())
matrix = json.loads(MATRIX_PATH.read_text())

need(matrix.get("schema") == "DE.PULSE-v18.5-MAJOR-CLOSURE-1", "major-closure schema drift")
need(matrix.get("version") == "18.5.0", "major-closure version drift")
need(matrix.get("scopeFile") == "v18_5_scope.json", "scope-file binding drift")

dimensions = matrix.get("closureDimensions", [])
expected_dimensions = scope.get("closureDimensions", [])
actual_dimensions = [row.get("id") for row in dimensions]
need(actual_dimensions == expected_dimensions, f"closure dimensions/order drift: {actual_dimensions}")
need(len(dimensions) == 10, f"expected 10 closure dimensions, found {len(dimensions)}")
for row in dimensions:
    ident = row.get("id", "<missing>")
    need(row.get("ownerGate") in {"G6", "G7", "G8", "G9", "G10"}, f"{ident}: invalid ownerGate")
    evidence = row.get("evidence")
    need(isinstance(evidence, list) and evidence and all(isinstance(x, str) and x.strip() for x in evidence), f"{ident}: closure evidence missing")

scenarios = matrix.get("adrGdiScenarios", [])
expected_scenarios = scope.get("adrGdiRequiredScenarios", [])
actual_scenarios = [row.get("id") for row in scenarios]
need(actual_scenarios == expected_scenarios, f"ADR-GDI scenarios/order drift: {actual_scenarios}")
need(len(scenarios) == 15, f"expected 15 ADR-GDI scenarios, found {len(scenarios)}")
for row in scenarios:
    ident = row.get("id", "<missing>")
    need(row.get("ownerGate") in {"G6", "G7", "G8", "G9", "G10"}, f"{ident}: invalid ownerGate")
    for field in ("selectorOrGate", "expectedTruth", "evidencePath"):
        need(isinstance(row.get(field), str) and row[field].strip(), f"{ident}: blank {field}")
    need(isinstance(row.get("requiresRealService"), bool), f"{ident}: requiresRealService must be boolean")
    need(row.get("status") == "required", f"{ident}: status must be required")
    serialized = json.dumps(row).lower()
    need("implied" not in serialized and "unverified" not in serialized, f"{ident}: implied/unverified evidence is forbidden")

by_id = {row.get("id"): row for row in scenarios}
for required_real in (
    "postgres_pressure_slow_unavailable",
    "restart_warm_start_recovery",
    "packaged_runtime_degradation_ux_diagnostics",
):
    need(by_id.get(required_real, {}).get("requiresRealService") is True, f"{required_real}: real-service execution is mandatory")

native = matrix.get("nativeDeliveryRequirements", {})
for field in (
    "macosAppleSiliconRunnableZip",
    "windowsX64RunnableZip",
    "exactSourceZip",
    "sha256Manifest",
    "certificationEvidence",
    "githubReleaseListsRunnableAssets",
):
    need(native.get(field) == "required", f"native delivery requirement missing/not fail-closed: {field}")

envelope = matrix.get("certificationOperatingEnvelope", {})
need(envelope.get("meaning") == "certification envelope only; not an advertised production capacity guarantee", "operating-envelope meaning drift")
need(envelope.get("radarCandidateBenchmark") == 500, "500-candidate benchmark envelope drift")
need(envelope.get("promotionStateStressCycles") == 10000, "10,000-cycle stress envelope drift")
need(envelope.get("symbolRotationStressCount") == 2000, "2,000-symbol stress envelope drift")
need(envelope.get("postgresConfiguredMaxConnectionsCeiling") == 128, "PostgreSQL pool ceiling envelope drift")
need(isinstance(envelope.get("evidence"), list) and len(envelope["evidence"]) >= 3, "operating-envelope evidence missing")

# The archive contract makes actual Release assets mandatory; a tag is not enough.
archive = (ROOT / "governance/RELEASE_ARTIFACT_ARCHIVE.md").read_text(errors="ignore")
for phrase in (
    "A tag alone is not sufficient",
    "macOS Apple Silicon runnable ZIP",
    "Windows x64 runnable ZIP",
    "GitHub Release",
):
    need(phrase in archive, f"release archive contract missing native delivery phrase: {phrase}")

if errors:
    print("v18.5 Major Closure Matrix Gate: FAIL")
    for error in errors:
        print(" -", error)
    sys.exit(2)

print("v18.5 Major Closure Matrix Gate: PASS · 10/10 closure dimensions · 15/15 ADR-GDI scenarios · native delivery fail-closed")
