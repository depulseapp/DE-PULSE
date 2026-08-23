#!/usr/bin/env python3
"""Permanent #70 repository-root ownership guard."""
from __future__ import annotations

from collections import Counter
import json
from pathlib import Path
import re
import subprocess
import sys

ROOT = Path(__file__).resolve().parents[2]
POLICY = ROOT / "governance" / "root-layout-policy.json"
MIGRATIONS = ROOT / "governance" / "repository-migrations.json"


def git(*args: str) -> str:
    return subprocess.run(("git", *args), cwd=ROOT, check=True, text=True, capture_output=True).stdout


def paths_at(commit: str | None = None) -> set[str]:
    raw = git("ls-tree", "-r", "--name-only", commit) if commit else git("ls-files")
    return {line.strip() for line in raw.splitlines() if line.strip()}


def root_paths(paths: set[str]) -> set[str]:
    return {path for path in paths if "/" not in path}


def main() -> int:
    errors: list[str] = []
    if not POLICY.is_file() or not MIGRATIONS.is_file():
        print("DE.PULSE permanent root ownership: FAIL", file=sys.stderr)
        print(" - missing root policy or repository migration registry", file=sys.stderr)
        return 1

    policy = json.loads(POLICY.read_text(encoding="utf-8"))
    migrations = json.loads(MIGRATIONS.read_text(encoding="utf-8"))
    baseline = str(policy.get("baselineCommit", "")).strip()
    if not baseline:
        errors.append("baselineCommit missing")
        baseline = "HEAD"
    if policy.get("grandfatherExistingBaselineRootFiles") is not False:
        errors.append("final root ownership must disable grandfatherExistingBaselineRootFiles")
    if policy.get("migrationRegisteredRootTargetsAreOwned") is not True:
        errors.append("migrationRegisteredRootTargetsAreOwned must be true")

    canonical = {str(x) for x in policy.get("canonicalRootFiles", []) if str(x)}
    final_evidence = policy.get("finalRootEvidenceFiles", {})
    if not isinstance(final_evidence, dict):
        errors.append("finalRootEvidenceFiles must be an object")
        final_evidence = {}
    for path, meta in final_evidence.items():
        if "/" in path or not isinstance(meta, dict) or not all(str(meta.get(key, "")).strip() for key in ("owner", "reason")):
            errors.append(f"invalid explicit final root evidence owner: {path}")

    classes = policy.get("baselineRootOwnershipClasses", [])
    compiled: list[tuple[str, re.Pattern[str]]] = []
    if not isinstance(classes, list) or not classes:
        errors.append("baselineRootOwnershipClasses must be a small non-empty list")
        classes = []
    for row in classes:
        if not isinstance(row, dict):
            errors.append("baseline ownership class must be an object")
            continue
        cid = str(row.get("id", "")).strip()
        pattern = str(row.get("pattern", "")).strip()
        owner = str(row.get("owner", "")).strip()
        scope = str(row.get("scope", "")).strip()
        disposition = str(row.get("disposition", "")).strip()
        if not all((cid, pattern, owner, scope, disposition)):
            errors.append(f"incomplete baseline ownership class: {cid or '<unnamed>'}")
            continue
        if scope != "EXISTING_BASELINE_ONLY":
            errors.append(f"baseline ownership class {cid} must be EXISTING_BASELINE_ONLY")
        try:
            compiled.append((cid, re.compile(pattern)))
        except re.error as exc:
            errors.append(f"invalid ownership regex {cid}: {exc}")

    migration_targets: dict[str, dict[str, object]] = {}
    rows = migrations.get("moves", [])
    for row in rows if isinstance(rows, list) else []:
        if not isinstance(row, dict):
            continue
        new = str(row.get("newPath", "")).strip()
        if not new or "/" in new:
            continue
        if not all(str(row.get(key, "")).strip() for key in ("owner", "reason", "removalCondition")):
            errors.append(f"root migration target lacks durable ownership metadata: {new}")
            continue
        migration_targets[new] = row

    current_all = paths_at()
    baseline_all = paths_at(baseline)
    current = root_paths(current_all)
    baseline_root = root_paths(baseline_all)
    new_root = current - baseline_root
    ownership = Counter()

    for path in sorted(current):
        if path in canonical:
            ownership["CANONICAL"] += 1
            continue
        if path in final_evidence:
            ownership["EXPLICIT_FINAL_EVIDENCE"] += 1
            continue
        if path in migration_targets:
            ownership["REGISTERED_MIGRATION_TARGET"] += 1
            continue
        if path in baseline_root:
            matches = [cid for cid, pattern in compiled if pattern.fullmatch(path)]
            if not matches:
                errors.append(f"unowned inherited root path: {path}")
            elif len(matches) > 1:
                errors.append(f"ambiguous inherited root ownership {path}: {', '.join(matches)}")
            else:
                ownership[matches[0]] += 1
            continue
        errors.append(f"unowned new root path: {path}")

    for path in sorted(new_root):
        if path not in canonical and path not in final_evidence and path not in migration_targets:
            errors.append(f"new root recurrence is not canonical/final-evidence/registered: {path}")

    aliases = migrations.get("temporaryAliases", [])
    if aliases:
        errors.append("temporaryAliases must be empty at #70 final root ownership")

    for line in git("ls-files", "-s").splitlines():
        meta, path = line.split("\t", 1)
        mode = meta.split()[0]
        if "/" not in path and mode == "120000":
            errors.append(f"root compatibility symlink is prohibited at closure: {path}")

    print("DE.PULSE permanent root ownership")
    print(f"baseline root files: {len(baseline_root)}")
    print(f"current root files: {len(current)}")
    print(f"root reduction: {len(baseline_root) - len(current)}")
    print(f"new root paths: {len(new_root)}")
    print(f"explicit final root evidence owners: {len(final_evidence)}")
    for key in sorted(ownership):
        print(f"owner[{key}]: {ownership[key]}")
    print("blanket baseline grandfathering: DISABLED")
    print("new root recurrence: CANONICAL_OR_EXPLICIT_FINAL_EVIDENCE_OR_REGISTERED_MIGRATION_ONLY")
    print("temporary root aliases/symlinks: PROHIBITED")

    if errors:
        print("DE.PULSE permanent root ownership: FAIL", file=sys.stderr)
        for error in errors:
            print(" - " + error, file=sys.stderr)
        return 1
    print("DE.PULSE permanent root ownership: PASS")
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
