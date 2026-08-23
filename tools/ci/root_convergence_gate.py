#!/usr/bin/env python3
"""Executable full-root disposition and #73 repository architecture guard."""
from __future__ import annotations

import argparse
from collections import Counter
import json
from pathlib import Path
import re
import subprocess
import sys

ROOT = Path(__file__).resolve().parents[2]
POLICY = ROOT / "governance" / "work-slices" / "ADAPT-ROOT-CONVERGENCE-001" / "root-convergence-policy.json"
ROOT_POLICY = ROOT / "governance" / "root-layout-policy.json"
STATE = ROOT / "governance" / "current-state.json"
SUPPORTED_SCHEMAS = {"DE.PULSE-ROOT-CONVERGENCE-POLICY-1", "DE.PULSE-ROOT-CONVERGENCE-POLICY-2"}
VERSIONED = re.compile(r"^v(?:17|18)(?:[_\.-]|$)", re.I)
VERSIONED_GO = re.compile(r"(?:^|_)v(?:17|18)(?:[_\.]\d+)*(?:_|\.|$)", re.I)
VERSIONED_RENDERER = re.compile(r"(?:^|[-_])v(?:17|18)(?:[._-]\d+)+(?:[._-]|$)", re.I)
CLOSED = {"READY_FOR_CLOSURE", "COMPLETE", "COMPLETED", "CLOSED", "DELIVERED"}


def tracked_root_files() -> list[str]:
    result = subprocess.run(("git", "ls-files"), cwd=ROOT, check=True, text=True, capture_output=True)
    return sorted(x.strip() for x in result.stdout.splitlines() if x.strip() and "/" not in x.strip())


def classify(name: str, canonical: set[str], remaining: set[str]) -> tuple[str, str]:
    if name in canonical:
        return "KEEP", "Canonical steady-state project root owner."
    if name.endswith(".go"):
        if name in remaining or VERSIONED_GO.search(name):
            return "CONSOLIDATE", "Active version-named package-main Go owner; rename or cohesively extract."
        return "KEEP", "Current package-main Go source/test owner; relocation is architectural, not cosmetic."
    if VERSIONED.match(name):
        return "MOVE", "Historical/version-scoped non-Go evidence belongs under governed release/history ownership."
    if Path(name).suffix.lower() in {".py", ".js", ".sh", ".ps1"}:
        return "MOVE", "Reusable root tooling/evidence should converge to a semantic tools or tests owner."
    if Path(name).suffix.lower() in {".json", ".md", ".txt"}:
        return "MOVE", "Non-canonical root metadata should converge to governance/release/test ownership."
    return "MOVE", "Non-canonical root file requires an explicit canonical owner outside repository root."


def versioned_renderer_runtime_candidates() -> list[str]:
    renderer = ROOT / "renderer"
    if not renderer.is_dir():
        return []
    return sorted(
        path.relative_to(ROOT).as_posix()
        for path in renderer.iterdir()
        if path.is_file() and VERSIONED_RENDERER.search(path.name)
    )


def main() -> int:
    parser = argparse.ArgumentParser()
    parser.add_argument("--json-out")
    args = parser.parse_args()
    errors: list[str] = []

    if not POLICY.is_file() or not ROOT_POLICY.is_file() or not STATE.is_file():
        print("DE.PULSE root convergence: FAIL", file=sys.stderr)
        return 1

    policy = json.loads(POLICY.read_text(encoding="utf-8"))
    root_policy = json.loads(ROOT_POLICY.read_text(encoding="utf-8"))
    state = json.loads(STATE.read_text(encoding="utf-8"))
    schema = str(policy.get("schema", ""))
    if schema not in SUPPORTED_SCHEMAS:
        errors.append("unsupported root-convergence policy schema")
    if policy.get("workSliceId") != "ADAPT-ROOT-CONVERGENCE-001" or policy.get("issue") != 73:
        errors.append("root-convergence policy identity drift")

    if schema == "DE.PULSE-ROOT-CONVERGENCE-POLICY-2":
        target_rel = str(policy.get("architectureTarget", "")).strip()
        if not target_rel or not (ROOT / target_rel).is_file():
            errors.append("repository architecture target missing")
        else:
            try:
                target = json.loads((ROOT / target_rel).read_text(encoding="utf-8"))
                if target.get("schema") != "DE.PULSE-REPOSITORY-ARCHITECTURE-TARGET-1":
                    errors.append("unsupported repository architecture target schema")
                if target.get("workSliceId") != "ADAPT-ROOT-CONVERGENCE-001" or target.get("issue") != 73:
                    errors.append("repository architecture target identity drift")
            except Exception as exc:
                errors.append(f"repository architecture target invalid: {exc}")

        targets = policy.get("versionedGoTargets", {})
        if not isinstance(targets, dict):
            errors.append("versionedGoTargets must be an object")
            targets = {}
        remaining_declared = {str(x) for x in policy.get("remainingVersionedGo", [])}
        if set(targets) != remaining_declared:
            errors.append("versionedGoTargets keys must exactly match remainingVersionedGo")
        for old, new in targets.items():
            if not str(new).endswith(".go") or VERSIONED_GO.search(str(new)):
                errors.append(f"versioned Go target must be a version-neutral Go filename: {old} -> {new}")

    canonical = {str(x) for x in root_policy.get("canonicalRootFiles", [])}
    remaining = {str(x) for x in policy.get("remainingVersionedGo", [])}
    files = tracked_root_files()
    rows: list[dict[str, str]] = []
    counts: Counter[str] = Counter()
    for name in files:
        disposition, reason = classify(name, canonical, remaining)
        rows.append({"path": name, "disposition": disposition, "reason": reason})
        counts[disposition] += 1

    actual = {n for n in files if n.endswith(".go") and VERSIONED_GO.search(n)}
    unregistered = sorted(actual - remaining)
    stale = sorted(remaining - actual)
    if unregistered:
        errors.append("unregistered version-named Go root owners: " + ", ".join(unregistered))

    phase = str(policy.get("phase", "")).upper()
    historical = sorted(n for n in files if VERSIONED.match(n) and not n.endswith(".go"))
    tools = sorted(
        n for n in files
        if Path(n).suffix.lower() in {".py", ".js", ".sh", ".ps1"}
        and n not in canonical
        and not VERSIONED.match(n)
    )
    moves = [row["path"] for row in rows if row["disposition"] == "MOVE"]
    renderer_versioned = versioned_renderer_runtime_candidates()

    historical_clean_phases = {"ARCHITECTURE_REAUDIT", "HISTORICAL_NON_GO_CLEAN", "TOOLING_CLEAN", "GO_CONVERGED", "FINAL"}
    tooling_clean_phases = {"TOOLING_CLEAN", "GO_CONVERGED", "FINAL"}
    go_clean_phases = {"GO_CONVERGED", "FINAL"}
    if phase in historical_clean_phases and historical:
        errors.append("historical v17/v18 non-Go root files remain: " + ", ".join(historical[:20]))
    if phase in tooling_clean_phases and tools:
        errors.append("non-canonical root tooling remains: " + ", ".join(tools[:20]))
    if phase in go_clean_phases and actual:
        errors.append("version-named Go root owners remain: " + ", ".join(sorted(actual)))

    active = state.get("activeWorkSlice", {})
    status = str(active.get("status", "")).upper() if isinstance(active, dict) else ""
    if phase == "FINAL" or status in CLOSED:
        if moves:
            errors.append("closure has unresolved MOVE dispositions: " + ", ".join(moves[:20]))
        if stale:
            errors.append("closure policy still registers converged versioned Go paths: " + ", ".join(stale))
        if renderer_versioned:
            errors.append("closure has version-named renderer owners requiring archive/rename proof: " + ", ".join(renderer_versioned[:20]))

    report = {
        "schema": "DE.PULSE-ROOT-DISPOSITION-INVENTORY-2",
        "workSliceId": "ADAPT-ROOT-CONVERGENCE-001",
        "issue": 73,
        "phase": phase,
        "status": "FAIL" if errors else "PASS",
        "rootFileCount": len(files),
        "dispositionCounts": dict(sorted(counts.items())),
        "historicalVersionedNonGoRootCount": len(historical),
        "versionedGoRootCount": len(actual),
        "versionedRendererCandidateCount": len(renderer_versioned),
        "nonCanonicalRootToolCount": len(tools),
        "unclassifiedCount": 0,
        "rows": rows,
        "errors": errors,
    }
    if args.json_out:
        output = Path(args.json_out)
        output.parent.mkdir(parents=True, exist_ok=True)
        output.write_text(json.dumps(report, indent=2, sort_keys=True) + "\n", encoding="utf-8")

    print("DE.PULSE #73 repository architecture convergence")
    print(f"root files: {len(files)}")
    print("dispositions: " + ", ".join(f"{key}={counts[key]}" for key in sorted(counts)))
    print(f"historical versioned non-Go root: {len(historical)}")
    print(f"versioned Go root: {len(actual)}")
    print(f"version-named renderer candidates: {len(renderer_versioned)}")
    print(f"non-canonical root tooling: {len(tools)}")
    print("unclassified root files: 0")
    print("MOVE root files: " + (", ".join(moves) if moves else "NONE"))
    print("CONSOLIDATE versioned Go: " + (", ".join(sorted(actual)) if actual else "NONE"))
    print("REVIEW versioned renderer: " + (", ".join(renderer_versioned) if renderer_versioned else "NONE"))
    if errors:
        print("DE.PULSE root convergence: FAIL", file=sys.stderr)
        for error in errors:
            print(" - " + error, file=sys.stderr)
        return 1
    print("DE.PULSE root convergence: PASS")
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
