#!/usr/bin/env python3
from __future__ import annotations

import argparse
import json
from pathlib import Path
import re
import sys

ROOT = Path(__file__).resolve().parents[2]
SELF = Path(__file__).resolve()

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

# These files intentionally carry old-path strings as negative-policy/migration
# markers. They are not executable consumers of the legacy paths.
LEGACY_REFERENCE_ALLOWLIST = {
    "tools/ci/legacy_test_gate_inventory.py": set(FIRST_WAVE_RENAMES),
    "tools/ci/workflow_policy.py": {
        "v18_6_surface_consolidation_test.js",
        "v18_6_documentation_access_test.js",
    },
}

VERSIONED_PREFIX = re.compile(r"^v\d+(?:[_\.]\d+)*[_-]", re.IGNORECASE)
EXECUTABLE_REFERENCE_RE = re.compile(r"(?P<path>[A-Za-z0-9_.\-/]+\.(?:py|js))")


def relative(path: Path) -> str:
    return path.resolve().relative_to(ROOT.resolve()).as_posix()


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
            out[relative(path)] = path.read_text(encoding="utf-8", errors="replace")
    return out


def executable_candidate(path: Path) -> bool:
    if not path.is_file() or path.resolve() == SELF:
        return False
    try:
        rel = relative(path)
    except ValueError:
        return False
    name = path.name.lower()
    if rel.startswith("tools/ci/") or rel.startswith("tools/release/"):
        return path.suffix.lower() in {".py", ".js"}
    if path.parent == ROOT:
        return name.endswith("_gate.py") or name.endswith("_test.py") or name.endswith("_test.js")
    if rel.startswith("release/"):
        return name.endswith("_test.py") or name.endswith("_test.js") or name.endswith("_gate.py")
    return False


def active_executable_text() -> dict[str, str]:
    """Resolve the current executable test/gate closure from canonical controls.

    This catches indirect consumers such as workflow_policy.py ->
    ai_continuous_eval_gate.py -> a renamed Go test file without declaring every
    historical root script active merely because it still exists.
    """
    texts = active_control_text()
    pending = list(texts.values())
    visited_paths = {ROOT / rel for rel in texts}

    while pending:
        text = pending.pop()
        for match in EXECUTABLE_REFERENCE_RE.finditer(text):
            token = match.group("path").lstrip("./")
            path = (ROOT / token).resolve()
            if path in visited_paths or not executable_candidate(path):
                continue
            try:
                path.relative_to(ROOT.resolve())
            except ValueError:
                continue
            body = path.read_text(encoding="utf-8", errors="replace")
            texts[relative(path)] = body
            visited_paths.add(path)
            pending.append(body)
    return texts


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
            "Current executable control closure references this path; migrate consumer and assertions atomically before old path removal.",
        )
    return (
        "UNREFERENCED_USEFUL",
        [],
        "No current executable control-plane consumer; assertion/evidence mapping is still required before deletion.",
    )


def legacy_path_consumers(texts: dict[str, str]) -> dict[str, list[str]]:
    consumers: dict[str, list[str]] = {}
    for old in FIRST_WAVE_RENAMES:
        hits: list[str] = []
        for rel, text in texts.items():
            if old not in text:
                continue
            allowed = LEGACY_REFERENCE_ALLOWLIST.get(rel, set())
            if old in allowed:
                continue
            hits.append(rel)
        consumers[old] = sorted(hits)
    return consumers


def inventory() -> dict[str, object]:
    controls = active_executable_text()
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
            "currentExecutableControlClosureReference": "ACTIVE_REQUIRED",
            "unreferencedExecutableDefault": "UNREFERENCED_USEFUL",
            "safeRemoval": "never inferred automatically; requires explicit assertion/evidence review",
        },
        "activeExecutableConsumerCount": len(controls),
        "legacyPathConsumers": legacy_path_consumers(controls),
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

    legacy_consumers = report.get("legacyPathConsumers", {})
    if isinstance(legacy_consumers, dict):
        for old, consumers in legacy_consumers.items():
            if consumers:
                errors.append(
                    "first-wave old path still referenced by current executable consumer: "
                    + f"{old} -> {', '.join(consumers)}"
                )

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
    print("active executable consumer closure: " + str(report["activeExecutableConsumerCount"]))
    print("first-wave capability-oriented renames/moves: PASS")
    print("first-wave legacy executable references: NONE")
    print("conserved v17/v18 reconciliation inventory remains Fast-bound: PASS")
    print("automatic deletion inference: PROHIBITED")
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
