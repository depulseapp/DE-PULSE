#!/usr/bin/env python3
"""#70 repository migration-safety guardrail.

Protects root-layout, stale references, executable modes, case-sensitive paths,
Go test discovery identities and source-tree cleanliness while files are moved.
The migration registry is the explicit exception mechanism; no implicit wrappers.
"""
from __future__ import annotations

import json
from pathlib import Path
import re
import subprocess
import sys

ROOT = Path(__file__).resolve().parents[2]
POLICY_PATH = ROOT / "governance" / "root-layout-policy.json"
MIGRATIONS_PATH = ROOT / "governance" / "repository-migrations.json"
STATE_PATH = ROOT / "governance" / "current-state.json"
TEXT_SUFFIXES = {
    ".go", ".py", ".js", ".mjs", ".cjs", ".sh", ".ps1", ".json", ".yml", ".yaml",
    ".md", ".html", ".css", ".txt", ".toml", ".xml",
}
GO_TEST_RE = re.compile(r"(?m)^func\s+(Test[A-Za-z0-9_]+|Benchmark[A-Za-z0-9_]+|Fuzz[A-Za-z0-9_]+)\s*\(")


def git(*args: str, check: bool = True) -> subprocess.CompletedProcess[str]:
    return subprocess.run(("git", *args), cwd=ROOT, check=check, text=True, capture_output=True)


def tracked_paths(commit: str | None = None) -> list[str]:
    if commit:
        result = git("ls-tree", "-r", "--name-only", commit)
    else:
        result = git("ls-files")
    return sorted(line.strip() for line in result.stdout.splitlines() if line.strip())


def tracked_modes(commit: str | None = None) -> dict[str, str]:
    if commit:
        result = git("ls-tree", "-r", commit)
        modes: dict[str, str] = {}
        for line in result.stdout.splitlines():
            # <mode> <type> <sha>\t<path>
            meta, path = line.split("\t", 1)
            modes[path] = meta.split()[0]
        return modes
    result = git("ls-files", "-s")
    modes = {}
    for line in result.stdout.splitlines():
        meta, path = line.split("\t", 1)
        modes[path] = meta.split()[0]
    return modes


def read_at(commit: str, path: str) -> str:
    result = git("show", f"{commit}:{path}", check=False)
    return result.stdout if result.returncode == 0 else ""


def go_test_identities(paths: list[str], commit: str | None = None) -> set[str]:
    identities: set[str] = set()
    for path in paths:
        if not path.endswith("_test.go"):
            continue
        text = read_at(commit, path) if commit else (ROOT / path).read_text(encoding="utf-8", errors="replace")
        package_dir = str(Path(path).parent.as_posix())
        for match in GO_TEST_RE.finditer(text):
            identities.add(f"{package_dir}:{match.group(1)}")
    return identities


def executable_candidate(path: str) -> bool:
    p = Path(path)
    name = p.name.lower()
    if path.startswith(("tools/ci/", "tools/release/", "tools/dev/")):
        return p.suffix.lower() in {".py", ".js", ".sh", ".ps1"}
    if path.startswith("release/"):
        return name.endswith(("_test.py", "_test.js", "_gate.py", ".sh", ".ps1"))
    if "/" not in path:
        return name.endswith(("_gate.py", "_test.py", "_test.js"))
    return False


def canonicalize(path: str, moves: dict[str, str]) -> str:
    return moves.get(path, path)


def textual_tracked_files(paths: list[str]) -> list[str]:
    out = []
    for path in paths:
        p = Path(path)
        if p.suffix.lower() in TEXT_SUFFIXES or p.name in {"Dockerfile", "Makefile"}:
            out.append(path)
    return out


def main() -> int:
    errors: list[str] = []
    for required in (POLICY_PATH, MIGRATIONS_PATH, STATE_PATH):
        if not required.is_file():
            errors.append(f"missing migration-safety authority: {required.relative_to(ROOT)}")
    if errors:
        for error in errors:
            print("FAIL:", error, file=sys.stderr)
        return 1

    policy = json.loads(POLICY_PATH.read_text())
    migrations = json.loads(MIGRATIONS_PATH.read_text())
    state = json.loads(STATE_PATH.read_text())
    baseline = str(policy.get("baselineCommit") or migrations.get("baselineCommit") or "").strip()
    if not baseline or git("cat-file", "-e", f"{baseline}^{{commit}}", check=False).returncode != 0:
        errors.append("root/migration baseline commit is missing or not resolvable")
        baseline = "HEAD"
    if migrations.get("baselineCommit") != policy.get("baselineCommit"):
        errors.append("root-layout and migration registries disagree on baseline commit")

    current_paths = tracked_paths()
    baseline_paths = tracked_paths(baseline)
    current_set = set(current_paths)
    baseline_set = set(baseline_paths)

    # Case-sensitive/cross-platform path safety.
    lowered: dict[str, list[str]] = {}
    for path in current_paths:
        lowered.setdefault(path.casefold(), []).append(path)
    for variants in lowered.values():
        if len(variants) > 1:
            errors.append("case-colliding tracked paths: " + ", ".join(sorted(variants)))

    # Root recurrence guard: baseline debt may remain while migrating, but no new
    # arbitrary root file can appear without explicit transitional registration.
    canonical_root = set(policy.get("canonicalRootFiles", []))
    transitional = policy.get("transitionalRootFiles", {})
    if not isinstance(transitional, dict):
        errors.append("transitionalRootFiles must be an object")
        transitional = {}
    current_root = {path for path in current_paths if "/" not in path}
    baseline_root = {path for path in baseline_paths if "/" not in path}
    new_root = current_root - baseline_root
    allowed_new_root = canonical_root | set(transitional)
    for path in sorted(new_root - allowed_new_root):
        errors.append(f"unregistered new root file: {path}")
    for path, meta in transitional.items():
        if path not in current_root:
            continue
        if not isinstance(meta, dict) or not all(str(meta.get(key, "")).strip() for key in ("owner", "reason", "expiry", "removalCondition")):
            errors.append(f"transitional root file lacks owner/reason/expiry/removalCondition: {path}")

    # Required top-level durable owners.
    for directory in policy.get("requiredTopLevelDirectories", []):
        if not any(path == directory or path.startswith(directory + "/") for path in current_paths):
            errors.append(f"required top-level owner missing: {directory}")

    # Explicit move registry + stale active reference scan.
    move_rows = migrations.get("moves", [])
    move_map: dict[str, str] = {}
    for row in move_rows if isinstance(move_rows, list) else []:
        if not isinstance(row, dict):
            errors.append("migration move entry must be an object")
            continue
        old = str(row.get("oldPath", "")).strip()
        new = str(row.get("newPath", "")).strip()
        if not old or not new:
            errors.append("migration move requires oldPath/newPath")
            continue
        move_map[old] = new
        if old in current_set:
            errors.append(f"registered old path still exists: {old}")
        if new not in current_set:
            errors.append(f"registered new path missing: {new}")
        for key in ("owner", "reason", "removalCondition"):
            if not str(row.get(key, "")).strip():
                errors.append(f"migration move {old} missing {key}")
        allowed_refs = set(row.get("allowedReferenceFiles", []))
        allowed_refs.add(MIGRATIONS_PATH.relative_to(ROOT).as_posix())
        for path in textual_tracked_files(current_paths):
            if path in allowed_refs:
                continue
            try:
                text = (ROOT / path).read_text(encoding="utf-8", errors="replace")
            except OSError:
                continue
            if old in text:
                errors.append(f"stale reference after move {old} -> {new}: {path}")

    # Preserve executable file mode for unchanged files and registered moves.
    baseline_modes = tracked_modes(baseline)
    current_modes = tracked_modes()
    for path, mode in baseline_modes.items():
        if mode != "100755":
            continue
        if path in current_set and current_modes.get(path) != mode:
            errors.append(f"executable mode lost: {path} baseline={mode} current={current_modes.get(path)}")
        if path in move_map:
            new = move_map[path]
            if current_modes.get(new) != mode:
                errors.append(f"executable mode lost across move: {path} -> {new}")

    # Go test identity conservation protects package-local discovery even when
    # test filenames are capability-renamed. Explicit identity remapping is needed
    # before package decomposition changes package paths.
    baseline_tests = go_test_identities(baseline_paths, baseline)
    current_tests = go_test_identities(current_paths)
    rename_map = migrations.get("testIdentityRenames", {})
    retired_tests = set(migrations.get("removedGoTestIdentities", []))
    expected_tests = {str(rename_map.get(identity, identity)) for identity in baseline_tests if identity not in retired_tests}
    missing_tests = sorted(expected_tests - current_tests)
    if missing_tests:
        errors.append("Go test identities disappeared without explicit evidence mapping: " + ", ".join(missing_tests[:20]))

    # Conserved executable path identity for script/gate/test owners. A move must
    # be registered; retirement must be explicit.
    baseline_exec = {path for path in baseline_paths if executable_candidate(path)}
    current_exec = {path for path in current_paths if executable_candidate(path)}
    retired_exec = set(migrations.get("retiredExecutablePaths", []))
    expected_exec = {canonicalize(path, move_map) for path in baseline_exec if path not in retired_exec}
    missing_exec = sorted(expected_exec - current_exec)
    if missing_exec:
        errors.append("executable test/gate paths disappeared without registered move/retirement: " + ", ".join(missing_exec[:20]))

    # Temporary aliases are exceptional and must expire.
    aliases = migrations.get("temporaryAliases", [])
    for alias in aliases if isinstance(aliases, list) else []:
        if not isinstance(alias, dict) or not all(str(alias.get(key, "")).strip() for key in ("path", "owner", "reason", "expiry", "removalCondition")):
            errors.append("temporary alias missing path/owner/reason/expiry/removalCondition")

    # Version-independent active work-slice truth.
    active = state.get("activeWorkSlice", {})
    work_slice_id = str(active.get("workSliceId", "")).strip()
    if not work_slice_id:
        errors.append("current-state activeWorkSlice.workSliceId missing")
    else:
        work_slice_path = ROOT / "governance" / "work-slices" / work_slice_id / "work-slice.json"
        if not work_slice_path.is_file():
            errors.append(f"canonical work-slice metadata missing: {work_slice_path.relative_to(ROOT)}")
        else:
            work_slice = json.loads(work_slice_path.read_text())
            if work_slice.get("workSliceId") != work_slice_id:
                errors.append("current-state/work-slice ID mismatch")
            if work_slice.get("publicProductVersion") is not None:
                errors.append("process-hardening work slice consumed a public product version")

    # Normal guards should leave a governed checkout clean. Ignored CI scratch
    # output is intentionally absent from porcelain status due to root .gitignore.
    status = git("status", "--porcelain", "--untracked-files=all").stdout.strip()
    if status:
        errors.append("migration guard left/observed an unclean governed checkout: " + status.replace("\n", " | ")[:1000])

    print("DE.PULSE repository migration safety")
    print(f"baseline tracked paths: {len(baseline_paths)}")
    print(f"current tracked paths: {len(current_paths)}")
    print(f"baseline root files: {len(baseline_root)}")
    print(f"current root files: {len(current_root)}")
    print(f"registered moves: {len(move_map)}")
    print(f"baseline Go test identities: {len(baseline_tests)}")
    print(f"current Go test identities: {len(current_tests)}")
    print(f"baseline executable gate/test paths: {len(baseline_exec)}")
    print(f"current executable gate/test paths: {len(current_exec)}")
    if errors:
        print("DE.PULSE repository migration safety: FAIL", file=sys.stderr)
        for error in errors:
            print(" - " + error, file=sys.stderr)
        return 1
    print("case-sensitive path safety: PASS")
    print("root recurrence guard: PASS")
    print("stale-reference registry: PASS")
    print("executable mode conservation: PASS")
    print("Go test identity conservation: PASS")
    print("executable path identity conservation: PASS")
    print("governed checkout cleanliness: PASS")
    print("DE.PULSE repository migration safety: PASS")
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
