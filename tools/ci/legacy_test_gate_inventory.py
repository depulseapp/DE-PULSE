#!/usr/bin/env python3
from __future__ import annotations

import argparse
import json
from pathlib import Path
import re
import sys

ROOT = Path(__file__).resolve().parents[2]

ACTIVE_CONTROL_FILES = (
    ROOT / ".github" / "workflows" / "ci-fast.yml",
    ROOT / ".github" / "workflows" / "ci-qualified.yml",
    ROOT / ".github" / "workflows" / "release.yml",
    ROOT / "certification_plan.json",
    ROOT / "tools" / "ci" / "workflow_policy.py",
)

FIRST_WAVE_RENAMES = {
    "v18_6_ai_hardening_test.go": "ai_hardening_test.go",
    "v18_6_broad_snapshot_broker_test.go": "broad_snapshot_broker_test.go",
    "v18_6_documentation_access_test.go": "documentation_access_test.go",
    "v18_6_session_intelligence_coordinator_test.go": "session_intelligence_coordinator_test.go",
    "v18_6_surface_consolidation_test.js": "tests/renderer/surface_consolidation_test.js",
    "v18_6_documentation_access_test.js": "tests/renderer/documentation_access_test.js",
}

VERSIONED_PREFIX = re.compile(r"^v\d+(?:[_\.]\d+)*[_-]", re.IGNORECASE)


def is_target(path: Path) -> bool:
    if path.parent != ROOT or not VERSIONED_PREFIX.match(path.name):
        return False
    name = path.name.lower()
    return (
        name.endswith("_test.go")
        or name.endswith("_test.js")
        or name.endswith("_test.py")
        or name.endswith("_gate.py")
    )


def active_control_text() -> dict[str, str]:
    out: dict[str, str] = {}
    for path in ACTIVE_CONTROL_FILES:
        if path.is_file():
            out[path.relative_to(ROOT).as_posix()] = path.read_text(encoding="utf-8", errors="replace")
    return out


def classify(path: Path, controls: dict[str, str]) -> tuple[str, list[str], str]:
    rel = path.relative_to(ROOT).as_posix()
    consumers = [name for name, text in controls.items() if rel in text or path.name in text]

    if path.name.endswith("_test.go"):
        return (
            "ACTIVE_REQUIRED",
            ["go test ./... (package discovery)", *consumers],
            "Go package regression remains active; rename/consolidate only with equivalent package-local coverage.",
        )
    if consumers:
        return (
            "ACTIVE_REQUIRED",
            consumers,
            "Explicit current CI/certification consumer exists; migrate consumer and assertions atomically before old path removal.",
        )
    return (
        "UNREFERENCED_USEFUL",
        [],
        "No direct current Fast/Qualified/Release/certification-plan consumer; assertion/evidence mapping is still required before deletion.",
    )


def inventory() -> dict[str, object]:
    controls = active_control_text()
    rows = []
    for path in sorted((p for p in ROOT.iterdir() if p.is_file() and is_target(p)), key=lambda p: p.name):
        classification, consumers, condition = classify(path, controls)
        rows.append(
            {
                "path": path.name,
                "classification": classification,
                "consumers": consumers,
                "sizeBytes": path.stat().st_size,
                "deletionOrMigrationCondition": condition,
            }
        )

    counts: dict[str, int] = {}
    for row in rows:
        key = str(row["classification"])
        counts[key] = counts.get(key, 0) + 1

    first_wave = []
    for old, new in FIRST_WAVE_RENAMES.items():
        old_path = ROOT / old
        new_path = ROOT / new
        first_wave.append(
            {
                "oldPath": old,
                "newPath": new,
                "oldPathAbsent": not old_path.exists(),
                "newPathPresent": new_path.is_file(),
            }
        )

    return {
        "schema": "DE.PULSE-LEGACY-TEST-GATE-INVENTORY-1",
        "scope": "root version-stacked executable tests/gates",
        "classificationVocabulary": [
            "ACTIVE_REQUIRED",
            "ACTIVE_DUPLICATE",
            "UNREFERENCED_USEFUL",
            "HISTORICAL_EVIDENCE",
            "SAFE_TO_REMOVE",
        ],
        "classificationPolicy": {
            "goTestDiscovery": "ACTIVE_REQUIRED",
            "directCurrentControlReference": "ACTIVE_REQUIRED",
            "unreferencedExecutableDefault": "UNREFERENCED_USEFUL",
            "safeRemoval": "never inferred automatically; requires explicit assertion/evidence review",
        },
        "counts": counts,
        "rows": rows,
        "firstWave": first_wave,
    }


def validate(report: dict[str, object]) -> list[str]:
    errors: list[str] = []
    rows = report["rows"]
    if not isinstance(rows, list) or not rows:
        errors.append("version-stacked root inventory unexpectedly empty")

    for row in rows if isinstance(rows, list) else []:
        if row.get("classification") not in {
            "ACTIVE_REQUIRED",
            "ACTIVE_DUPLICATE",
            "UNREFERENCED_USEFUL",
            "HISTORICAL_EVIDENCE",
            "SAFE_TO_REMOVE",
        }:
            errors.append(f"invalid classification for {row.get('path')}")
        if row.get("classification") == "SAFE_TO_REMOVE":
            errors.append(f"automatic inventory may not infer SAFE_TO_REMOVE: {row.get('path')}")

    for item in report["firstWave"] if isinstance(report["firstWave"], list) else []:
        if not item.get("oldPathAbsent"):
            errors.append(f"first-wave old path still present: {item.get('oldPath')}")
        if not item.get("newPathPresent"):
            errors.append(f"first-wave new path missing: {item.get('newPath')}")

    fast = (ROOT / ".github" / "workflows" / "ci-fast.yml").read_text(encoding="utf-8")
    for required in (
        "node tests/renderer/surface_consolidation_test.js",
        "node tests/renderer/documentation_access_test.js",
        "python3 v18_5_1_v17_v18_reconciliation_gate.py",
    ):
        if required not in fast:
            errors.append(f"Fast consumer missing governed cleanup/reconciliation proof: {required}")
    for forbidden in (
        "node v18_6_surface_consolidation_test.js",
        "node v18_6_documentation_access_test.js",
    ):
        if forbidden in fast:
            errors.append(f"Fast still consumes legacy renderer-test path: {forbidden}")

    return errors


def main() -> int:
    parser = argparse.ArgumentParser(description="Inventory version-stacked DE.PULSE root tests/gates without guessing deletion safety")
    parser.add_argument("--json-out")
    args = parser.parse_args()

    report = inventory()
    errors = validate(report)
    text = json.dumps(report, indent=2, sort_keys=True) + "\n"
    if args.json_out:
        out = Path(args.json_out)
        out.parent.mkdir(parents=True, exist_ok=True)
        out.write_text(text, encoding="utf-8")

    if errors:
        print("DE.PULSE legacy test/gate inventory: FAIL", file=sys.stderr)
        for error in errors:
            print(f" - {error}", file=sys.stderr)
        return 1

    counts = report["counts"]
    print("DE.PULSE legacy test/gate inventory: PASS")
    print("root version-stacked executables classified: " + str(sum(counts.values())))
    print("classifications: " + ", ".join(f"{k}={v}" for k, v in sorted(counts.items())))
    print("first-wave capability-oriented renames/moves: PASS")
    print("conserved v17/v18 reconciliation inventory remains Fast-bound: PASS")
    print("automatic deletion inference: PROHIBITED")
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
