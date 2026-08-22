#!/usr/bin/env python3
from __future__ import annotations

import argparse
import json
from pathlib import Path
import re
import subprocess
import sys

ROOT = Path(__file__).resolve().parents[2]
SELF = Path(__file__).resolve()
ROOT_POLICY = ROOT / "governance" / "root-layout-policy.json"

FIRST_WAVE_RENAMES = {
    "v18_6_ai_hardening_test.go": "ai_hardening_test.go",
    "v18_6_broad_snapshot_broker_test.go": "broad_snapshot_broker_test.go",
    "v18_6_documentation_access_test.go": "documentation_access_test.go",
    "v18_6_session_intelligence_coordinator_test.go": "session_intelligence_coordinator_test.go",
    "v18_6_surface_consolidation_test.js": "tests/renderer/surface_consolidation_test.js",
    "v18_6_documentation_access_test.js": "tests/renderer/documentation_access_test.js",
}

LEGACY_REFERENCE_ALLOWLIST = {
    "tools/ci/legacy_test_gate_inventory.py": set(FIRST_WAVE_RENAMES),
    "tools/ci/workflow_policy.py": {
        "v18_6_surface_consolidation_test.js",
        "v18_6_documentation_access_test.js",
    },
}

VERSIONED_PREFIX = re.compile(r"^v\d+(?:[_\.]\d+)*[_-]", re.IGNORECASE)
EXECUTABLE_REFERENCE_RE = re.compile(
    r"(?P<path>[A-Za-z0-9_.\-/]+\.(?:py|js|sh|ps1))",
    re.IGNORECASE,
)


def git(*args: str, check: bool = True) -> subprocess.CompletedProcess[str]:
    return subprocess.run(
        ("git", *args),
        cwd=ROOT,
        check=check,
        text=True,
        capture_output=True,
    )


def relative(path: Path) -> str:
    return path.resolve().relative_to(ROOT.resolve()).as_posix()


def load_json(path: Path) -> dict[str, object]:
    return json.loads(path.read_text(encoding="utf-8"))


def canonical_g12_executor() -> Path:
    path = ROOT / "tools" / "release" / "run_full_certification.py"
    if not path.is_file():
        raise RuntimeError(f"canonical version-neutral G12 executor missing: {relative(path)}")
    return path


def current_release_manifest() -> Path:
    identity = load_json(ROOT / "release_identity.json")
    version = str(identity.get("version", "")).strip()
    if not version:
        raise RuntimeError("release_identity.json version missing")
    path = ROOT / "release" / f"v{version}" / "certification-manifest.json"
    if not path.is_file():
        raise RuntimeError(f"current release G12 manifest missing: {relative(path)}")
    return path


def active_control_files() -> tuple[Path, ...]:
    # Historical certification plans and per-version G12 shell orchestrators are
    # deliberately excluded. Current G12 ownership is the version-neutral
    # tools/release executor plus the current declarative release manifest.
    return (
        ROOT / ".github" / "workflows" / "ci-fast.yml",
        ROOT / ".github" / "workflows" / "ci-qualified.yml",
        ROOT / ".github" / "workflows" / "release.yml",
        ROOT / "tools" / "ci" / "workflow_policy.py",
        canonical_g12_executor(),
        current_release_manifest(),
    )


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
    for path in active_control_files():
        if not path.is_file():
            raise RuntimeError(f"active control missing: {relative(path)}")
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
    if rel.startswith(("tools/ci/", "tools/release/", "tools/dev/")):
        return path.suffix.lower() in {".py", ".js", ".sh", ".ps1"}
    if path.parent == ROOT:
        return (
            name.endswith("_gate.py")
            or name.endswith("_test.py")
            or name.endswith("_test.js")
        )
    if rel.startswith("release/"):
        # Every version-specific full-certification shell is historical after
        # canonical G12 convergence. Current release behavior is driven by the
        # declarative manifest already included in active_control_files().
        if name == "run_full_certification.sh":
            return False
        # Canonical G12 may deliberately reuse older capability/browser tests.
        # Those tests remain active evidence; historical G12 orchestrators do not.
        return (
            name.endswith("_test.py")
            or name.endswith("_test.js")
            or name.endswith("_gate.py")
        )
    return False


def active_executable_text() -> dict[str, str]:
    """Resolve executable test/gate closure from current canonical controls."""
    texts = active_control_text()
    pending = list(texts.values())
    visited_paths = {(ROOT / rel).resolve() for rel in texts}

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
    consumers = [
        name for name, text in controls.items() if rel in text or path.name in text
    ]

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


def tracked_root_files(commit: str | None = None) -> set[str]:
    if commit:
        result = git("ls-tree", "-r", "--name-only", commit)
    else:
        result = git("ls-files")
    return {
        line.strip()
        for line in result.stdout.splitlines()
        if line.strip() and "/" not in line.strip()
    }


def retained_assets() -> dict[str, dict[str, object]]:
    registry = load_json(ROOT / "source_retained_assets.json")
    rows = registry.get("assets", [])
    if not isinstance(rows, list):
        raise RuntimeError("source_retained_assets.json assets must be a list")
    out: dict[str, dict[str, object]] = {}
    for row in rows:
        if isinstance(row, dict):
            path = str(row.get("path", "")).strip()
            if path:
                out[path] = row
    return out


def classify_root_file(
    name: str,
    controls: dict[str, str],
    policy: dict[str, object],
    baseline_root: set[str],
    retained: dict[str, dict[str, object]],
) -> dict[str, object]:
    path = ROOT / name
    canonical = set(policy.get("canonicalRootFiles", []))
    transitional = (
        policy.get("transitionalRootFiles", {})
        if isinstance(policy.get("transitionalRootFiles", {}), dict)
        else {}
    )
    consumers = [rel for rel, text in controls.items() if name in text]

    if name in canonical:
        classification = "CANONICAL_ROOT"
        reason = "Explicit steady-state root allowlist owner."
    elif name in transitional:
        classification = "TRANSITIONAL_ROOT"
        reason = "Explicit temporary root exception with owner/expiry/removal condition."
    elif name in retained:
        classification = "RETAINED_ASSET"
        reason = str(
            retained[name].get(
                "purpose",
                "Explicit retained source asset; move requires atomic consumer updates.",
            )
        )
    elif name.endswith("_test.go"):
        classification = "ACTIVE_TEST"
        reason = "Package-main Go regression; migration must conserve test identity and package-local access."
    elif name.endswith(".go"):
        classification = "ACTIVE_RUNTIME"
        reason = "Current package-main production source; package decomposition follows recursive source-health guardrails."
    elif consumers:
        classification = "ACTIVE_TOOL"
        reason = "Referenced by current canonical executable control closure; move consumers atomically before removal."
    elif VERSIONED_PREFIX.match(name):
        classification = "MIGRATION_CANDIDATE"
        reason = "Grandfathered version-scoped root material; assertion/evidence mapping required before move/delete."
    elif name in {
        "certification_plan.json",
        "certification_runner.py",
        "ci_pipeline.py",
        "ci_pipeline_plan.json",
    }:
        classification = "MIGRATION_CANDIDATE"
        reason = "Known stale/competing legacy orchestration debt; consolidate/retire in #70 Wave 2 after equivalence proof."
    elif name in baseline_root:
        classification = "MIGRATION_CANDIDATE"
        reason = "Grandfathered Stable-baseline root file; explicit KEEP/MOVE/CONSOLIDATE/DELETE disposition required by #70."
    else:
        classification = "UNCLASSIFIED_NEW"
        reason = "New root file is neither canonical nor registered transitional debt."

    return {
        "path": name,
        "classification": classification,
        "consumers": sorted(consumers),
        "sizeBytes": path.stat().st_size if path.is_file() else 0,
        "reason": reason,
    }


def root_layout_inventory(controls: dict[str, str]) -> dict[str, object]:
    if not ROOT_POLICY.is_file():
        return {
            "error": "governance/root-layout-policy.json missing",
            "rows": [],
            "counts": {},
        }
    policy = load_json(ROOT_POLICY)
    baseline = str(policy.get("baselineCommit", "")).strip()
    if not baseline:
        return {"error": "root-layout baselineCommit missing", "rows": [], "counts": {}}

    baseline_root = tracked_root_files(baseline)
    current_root = tracked_root_files()
    retained = retained_assets()
    rows = [
        classify_root_file(name, controls, policy, baseline_root, retained)
        for name in sorted(current_root)
    ]
    counts: dict[str, int] = {}
    for row in rows:
        key = str(row["classification"])
        counts[key] = counts.get(key, 0) + 1

    return {
        "schema": "DE.PULSE-ROOT-LAYOUT-INVENTORY-1",
        "baselineCommit": baseline,
        "baselineRootFileCount": len(baseline_root),
        "currentRootFileCount": len(rows),
        "canonicalRootFiles": sorted(policy.get("canonicalRootFiles", [])),
        "newTrackedRootFiles": sorted(current_root - baseline_root),
        "rows": rows,
        "counts": counts,
        "automaticDeletionInference": "PROHIBITED",
    }


def inventory() -> dict[str, object]:
    controls = active_executable_text()
    rows = []
    for path in sorted(
        (p for p in ROOT.iterdir() if p.is_file() and is_target(p)),
        key=lambda p: p.name,
    ):
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
        "schema": "DE.PULSE-LEGACY-TEST-GATE-INVENTORY-2",
        "scope": "root version-stacked executable tests/gates plus complete tracked root-layout classification",
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
            "staleCertificationPlanIsCurrentControl": False,
            "historicalReleaseG12IsCurrentControl": False,
            "canonicalVersionNeutralG12IsCurrentControl": True,
        },
        "activeControlFiles": sorted(relative(path) for path in active_control_files()),
        "activeExecutableConsumerCount": len(controls),
        "legacyPathConsumers": legacy_path_consumers(controls),
        "counts": counts,
        "rows": rows,
        "firstWave": first_wave,
        "rootLayout": root_layout_inventory(controls),
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
            errors.append(
                f"automatic inventory may not infer SAFE_TO_REMOVE: {row.get('path')}"
            )

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

    active_controls = report.get("activeControlFiles", [])
    if "certification_plan.json" in active_controls:
        errors.append(
            "legacy certification_plan.json may not be a current executable control owner"
        )
    canonical_g12 = "tools/release/run_full_certification.py"
    if canonical_g12 not in active_controls:
        errors.append(
            f"canonical version-neutral G12 missing from executable control closure: {canonical_g12}"
        )
    identity = load_json(ROOT / "release_identity.json")
    current_manifest = f"release/v{identity['version']}/certification-manifest.json"
    if current_manifest not in active_controls:
        errors.append(
            f"current release G12 manifest missing from control closure: {current_manifest}"
        )
    version_shells = [
        rel
        for rel in active_controls
        if rel.startswith("release/") and rel.endswith("/run_full_certification.sh")
    ]
    if version_shells:
        errors.append(
            "version-specific release G12 shell incorrectly treated as current control: "
            + ", ".join(version_shells)
        )

    root_layout = report.get("rootLayout", {})
    if not isinstance(root_layout, dict) or root_layout.get("error"):
        errors.append(str(root_layout.get("error", "root-layout inventory missing")))
    else:
        root_rows = root_layout.get("rows", [])
        if not isinstance(root_rows, list) or not root_rows:
            errors.append("complete tracked root-layout inventory unexpectedly empty")
        for row in root_rows if isinstance(root_rows, list) else []:
            if row.get("classification") == "UNCLASSIFIED_NEW":
                errors.append(
                    f"new root file lacks canonical/transitional disposition: {row.get('path')}"
                )
        if root_layout.get("automaticDeletionInference") != "PROHIBITED":
            errors.append("automatic root deletion inference must remain PROHIBITED")
        if not any(
            isinstance(row, dict) and row.get("classification") == "RETAINED_ASSET"
            for row in root_rows if isinstance(root_rows, list)
        ):
            errors.append("retained branding/source asset is not explicitly classified")

    fast = (ROOT / ".github" / "workflows" / "ci-fast.yml").read_text(
        encoding="utf-8"
    )
    for required in (
        "node tests/renderer/surface_consolidation_test.js",
        "node tests/renderer/documentation_access_test.js",
        "python3 v18_5_1_v17_v18_reconciliation_gate.py",
    ):
        if required not in fast:
            errors.append(
                f"Fast consumer missing governed cleanup/reconciliation proof: {required}"
            )
    for forbidden in (
        "node v18_6_surface_consolidation_test.js",
        "node v18_6_documentation_access_test.js",
    ):
        if forbidden in fast:
            errors.append(f"Fast still consumes legacy renderer-test path: {forbidden}")

    return errors


def main() -> int:
    parser = argparse.ArgumentParser(
        description=(
            "Inventory DE.PULSE legacy executables and the complete tracked root "
            "layout without guessing deletion safety"
        )
    )
    parser.add_argument("--json-out")
    args = parser.parse_args()

    try:
        report = inventory()
        errors = validate(report)
    except Exception as exc:
        print("DE.PULSE legacy/root inventory: FAIL", file=sys.stderr)
        print(f" - {exc}", file=sys.stderr)
        return 1

    text = json.dumps(report, indent=2, sort_keys=True) + "\n"
    if args.json_out:
        out = Path(args.json_out)
        out.parent.mkdir(parents=True, exist_ok=True)
        out.write_text(text, encoding="utf-8")

    if errors:
        print("DE.PULSE legacy/root inventory: FAIL", file=sys.stderr)
        for error in errors:
            print(f" - {error}", file=sys.stderr)
        return 1

    counts = report["counts"]
    root_layout = report["rootLayout"]
    print("DE.PULSE legacy/root inventory: PASS")
    print(
        "root version-stacked executables classified: "
        + str(sum(counts.values()))
    )
    print(
        "classifications: "
        + ", ".join(f"{key}={value}" for key, value in sorted(counts.items()))
    )
    print(
        "active executable consumer closure: "
        + str(report["activeExecutableConsumerCount"])
    )
    print("active control files: " + ", ".join(report["activeControlFiles"]))
    print("stale certification_plan.json current-control ownership: REMOVED")
    print("canonical version-neutral G12 current-control ownership: tools/release/run_full_certification.py")
    print("historical release G12 current-control ownership: EXCLUDED")
    print(
        "complete tracked root files classified: "
        + str(root_layout["currentRootFileCount"])
    )
    print(
        "root classifications: "
        + ", ".join(
            f"{key}={value}" for key, value in sorted(root_layout["counts"].items())
        )
    )
    print("new unclassified root files: 0")
    print("retained branding/assets classification: PASS")
    print("first-wave capability-oriented renames/moves: PASS")
    print("first-wave legacy executable references: NONE")
    print("conserved v17/v18 reconciliation inventory remains Fast-bound: PASS")
    print("automatic deletion inference: PROHIBITED")
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
