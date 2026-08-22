#!/usr/bin/env python3
"""#70 repository migration-safety guardrail.

Protects root-layout, stale active references, executable modes, case-sensitive
paths, Go test discovery identities and source-tree cleanliness while files move.
The migration registry is the explicit exception mechanism; historical evidence
is immutable and is not treated as a current executable path consumer.
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
VERSIONED_ROOT_RE = re.compile(r"^v(?P<major>[0-9]+)(?:[_\.-]|$)", re.IGNORECASE)


def git(*args: str, check: bool = True) -> subprocess.CompletedProcess[str]:
    return subprocess.run(("git", *args), cwd=ROOT, check=check, text=True, capture_output=True)


def tracked_paths(commit: str | None = None) -> list[str]:
    result = git("ls-tree", "-r", "--name-only", commit) if commit else git("ls-files")
    return sorted(line.strip() for line in result.stdout.splitlines() if line.strip())


def tracked_modes(commit: str | None = None) -> dict[str, str]:
    result = git("ls-tree", "-r", commit) if commit else git("ls-files", "-s")
    modes: dict[str, str] = {}
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
    out: list[str] = []
    for path in paths:
        p = Path(path)
        if p.suffix.lower() in TEXT_SUFFIXES or p.name in {"Dockerfile", "Makefile"}:
            out.append(path)
    return out


def root_version_major(path: str) -> int | None:
    if "/" in path:
        return None
    match = VERSIONED_ROOT_RE.match(path)
    return int(match.group("major")) if match else None


def policy_retired_root_path(path: str, retained_majors: set[int]) -> bool:
    major = root_version_major(path)
    return major is not None and major not in retained_majors


def current_active_reference_texts(current_paths: list[str]) -> dict[str, str]:
    """Return current executable/control and runtime text, excluding history/prose."""
    texts: dict[str, str] = {}
    try:
        from tools.ci.legacy_test_gate_inventory import active_executable_text
        texts.update(active_executable_text())
    except Exception as exc:
        raise RuntimeError(f"cannot resolve canonical executable closure: {exc}") from exc

    for path in current_paths:
        p = ROOT / path
        if path.endswith(".go") and not path.endswith("_test.go"):
            texts[path] = p.read_text(encoding="utf-8", errors="replace")
            continue
        if path.startswith("renderer/") and not path.startswith(("renderer/qa/", "renderer/docs/")):
            if p.suffix.lower() in {".js", ".html", ".css"}:
                texts[path] = p.read_text(encoding="utf-8", errors="replace")
    return texts


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

    retained_majors_raw = policy.get("retainedRootVersionMajors", [])
    try:
        retained_majors = {int(value) for value in retained_majors_raw}
    except Exception:
        retained_majors = set()
        errors.append("retainedRootVersionMajors must contain integers")
    if retained_majors != {17, 18}:
        errors.append(f"#70 current-root retention must be exactly v17/v18, got {sorted(retained_majors)}")

    retention = migrations.get("rootVersionRetention", {})
    if set(retention.get("retainedMajors", [])) != retained_majors:
        errors.append("root-layout and migration retention-major policy drift")
    if str(retention.get("scope", "")).strip() != "CURRENT_REPOSITORY_ROOT_ONLY":
        errors.append("root-version retirement scope must be CURRENT_REPOSITORY_ROOT_ONLY")
    if str(retention.get("historicalGitAndReleaseEvidence", "")).strip() != "IMMUTABLE":
        errors.append("historical Git/release evidence immutability contract missing")

    current_paths = tracked_paths()
    baseline_paths = tracked_paths(baseline)
    current_set = set(current_paths)

    lowered: dict[str, list[str]] = {}
    for path in current_paths:
        lowered.setdefault(path.casefold(), []).append(path)
    for variants in lowered.values():
        if len(variants) > 1:
            errors.append("case-colliding tracked paths: " + ", ".join(sorted(variants)))

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

    retired_root = sorted(path for path in current_root if policy_retired_root_path(path, retained_majors))
    if retired_root:
        errors.append("pre-v17 version-scoped files remain in current root: " + ", ".join(retired_root[:30]))

    for directory in policy.get("requiredTopLevelDirectories", []):
        if not any(path == directory or path.startswith(directory + "/") for path in current_paths):
            errors.append(f"required top-level owner missing: {directory}")

    try:
        active_refs = current_active_reference_texts(current_paths)
    except Exception as exc:
        errors.append(str(exc))
        active_refs = {}

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
        scope = str(row.get("referenceScope", "ALL_TEXT")).strip().upper()
        if scope == "ACTIVE_RUNTIME_AND_CONTROL":
            scan = active_refs
        elif scope == "ALL_TEXT":
            scan = {
                path: (ROOT / path).read_text(encoding="utf-8", errors="replace")
                for path in textual_tracked_files(current_paths)
            }
        else:
            errors.append(f"migration move {old} has unsupported referenceScope {scope}")
            scan = {}
        for path, text in scan.items():
            if path in allowed_refs:
                continue
            if old in text:
                errors.append(f"stale active reference after move {old} -> {new}: {path}")

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

    baseline_tests_all = go_test_identities(baseline_paths, baseline)
    baseline_retained_paths = [p for p in baseline_paths if not policy_retired_root_path(p, retained_majors)]
    baseline_tests = go_test_identities(baseline_retained_paths, baseline)
    policy_retired_test_count = len(baseline_tests_all - baseline_tests)
    current_tests = go_test_identities(current_paths)
    rename_map = migrations.get("testIdentityRenames", {})
    retired_tests = set(migrations.get("removedGoTestIdentities", []))
    expected_tests = {str(rename_map.get(identity, identity)) for identity in baseline_tests if identity not in retired_tests}
    missing_tests = sorted(expected_tests - current_tests)
    if missing_tests:
        errors.append("Go test identities disappeared without explicit evidence mapping: " + ", ".join(missing_tests[:20]))

    baseline_exec_all = {path for path in baseline_paths if executable_candidate(path)}
    baseline_exec = {path for path in baseline_exec_all if not policy_retired_root_path(path, retained_majors)}
    policy_retired_exec_count = len(baseline_exec_all - baseline_exec)
    current_exec = {path for path in current_paths if executable_candidate(path)}
    retired_exec = set(migrations.get("retiredExecutablePaths", []))
    expected_exec = {canonicalize(path, move_map) for path in baseline_exec if path not in retired_exec}
    missing_exec = sorted(expected_exec - current_exec)
    if missing_exec:
        errors.append("executable test/gate paths disappeared without registered move/retirement: " + ", ".join(missing_exec[:20]))

    aliases = migrations.get("temporaryAliases", [])
    for alias in aliases if isinstance(aliases, list) else []:
        if not isinstance(alias, dict) or not all(str(alias.get(key, "")).strip() for key in ("path", "owner", "reason", "expiry", "removalCondition")):
            errors.append("temporary alias missing path/owner/reason/expiry/removalCondition")

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

    status = git("status", "--porcelain", "--untracked-files=all").stdout.strip()
    if status:
        errors.append("migration guard left/observed an unclean governed checkout: " + status.replace("\n", " | ")[:1000])

    print("DE.PULSE repository migration safety")
    print(f"baseline tracked paths: {len(baseline_paths)}")
    print(f"current tracked paths: {len(current_paths)}")
    print(f"baseline root files: {len(baseline_root)}")
    print(f"current root files: {len(current_root)}")
    print("retained root version majors: " + ", ".join(f"v{x}" for x in sorted(retained_majors)))
    print(f"policy-retired pre-v17 Go test identities: {policy_retired_test_count}")
    print(f"policy-retired pre-v17 executable paths: {policy_retired_exec_count}")
    print(f"registered moves: {len(move_map)}")
    print(f"baseline retained Go test identities: {len(baseline_tests)}")
    print(f"current Go test identities: {len(current_tests)}")
    print(f"baseline retained executable gate/test paths: {len(baseline_exec)}")
    print(f"current executable gate/test paths: {len(current_exec)}")
    if errors:
        print("DE.PULSE repository migration safety: FAIL", file=sys.stderr)
        for error in errors:
            print(" - " + error, file=sys.stderr)
        return 1
    print("case-sensitive path safety: PASS")
    print("root recurrence + v17/v18 retention guard: PASS")
    print("stale active-reference registry: PASS")
    print("executable mode conservation: PASS")
    print("non-retired Go test identity conservation: PASS")
    print("non-retired executable path identity conservation: PASS")
    print("historical Git/release evidence immutability boundary: PASS")
    print("governed checkout cleanliness: PASS")
    print("DE.PULSE repository migration safety: PASS")
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
