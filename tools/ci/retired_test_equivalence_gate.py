#!/usr/bin/env python3
"""Fail-closed equivalence contract for #70 policy-retired pre-v17 tests/gates."""
from __future__ import annotations

import json
from pathlib import Path
import re
import subprocess
import sys

ROOT = Path(__file__).resolve().parents[2]
MIGRATIONS = ROOT / "governance" / "repository-migrations.json"
POLICY = ROOT / "governance" / "root-layout-policy.json"
GO_TEST_RE = re.compile(r"(?m)^func\s+(Test[A-Za-z0-9_]+|Benchmark[A-Za-z0-9_]+|Fuzz[A-Za-z0-9_]+)\s*\(")
VERSIONED_ROOT_RE = re.compile(r"^v(?P<major>[0-9]+)(?:[_\.-]|$)", re.IGNORECASE)
ALLOWED = {
    "SUPERSEDED_BY_CURRENT_CAPABILITY_EVIDENCE",
    "MIGRATED_TO_CAPABILITY_OWNER",
    "JUSTIFIED_HISTORICAL_SCOPE_RETIREMENT",
}


def git(*args: str, check: bool = True) -> subprocess.CompletedProcess[str]:
    return subprocess.run(("git", *args), cwd=ROOT, text=True, capture_output=True, check=check)


def tracked_paths(commit: str | None = None) -> list[str]:
    result = git("ls-tree", "-r", "--name-only", commit) if commit else git("ls-files")
    return sorted(line.strip() for line in result.stdout.splitlines() if line.strip())


def read_at(commit: str, path: str) -> str:
    result = git("show", f"{commit}:{path}", check=False)
    return result.stdout if result.returncode == 0 else ""


def root_version_major(path: str) -> int | None:
    if "/" in path:
        return None
    match = VERSIONED_ROOT_RE.match(path)
    return int(match.group("major")) if match else None


def retired_root(path: str, retained: set[int]) -> bool:
    major = root_version_major(path)
    return major is not None and major not in retained


def go_ids_by_path(paths: list[str], commit: str | None = None) -> dict[str, list[str]]:
    out: dict[str, list[str]] = {}
    for path in paths:
        if not path.endswith("_test.go"):
            continue
        text = read_at(commit, path) if commit else (ROOT / path).read_text(encoding="utf-8", errors="replace")
        package_dir = Path(path).parent.as_posix()
        ids = sorted({f"{package_dir}:{m.group(1)}" for m in GO_TEST_RE.finditer(text)})
        if ids:
            out[path] = ids
    return out


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


def flatten_go_groups(groups: object, errors: list[str]) -> dict[str, dict[str, object]]:
    out: dict[str, dict[str, object]] = {}
    if not isinstance(groups, list):
        errors.append("goGroups must be a list")
        return out
    for group in groups:
        if not isinstance(group, dict):
            errors.append("goGroups entry must be an object")
            continue
        for key in ("id", "capability", "disposition", "reason"):
            if not str(group.get(key, "")).strip():
                errors.append(f"go group missing {key}")
        disposition = str(group.get("disposition", ""))
        if disposition not in ALLOWED or disposition == "JUSTIFIED_HISTORICAL_SCOPE_RETIREMENT":
            errors.append(f"invalid Go-family disposition: {group.get('id')}={disposition}")
        bindings = group.get("bindings", [])
        if not isinstance(bindings, list) or not bindings:
            errors.append(f"go group missing bindings: {group.get('id')}")
            bindings = []
        sources = group.get("sources", [])
        if not isinstance(sources, list) or not sources:
            errors.append(f"go group missing sources: {group.get('id')}")
            continue
        for source in sources:
            if not isinstance(source, dict):
                errors.append(f"go source row invalid in {group.get('id')}")
                continue
            path = str(source.get("path", "")).strip()
            if not path or path in out:
                errors.append(f"duplicate/empty retired Go source mapping: {path}")
                continue
            out[path] = {"identities": source.get("identities"), "bindings": bindings, "disposition": disposition, "group": group.get("id")}
    return out


def flatten_exec_groups(groups: object, errors: list[str]) -> dict[str, dict[str, object]]:
    out: dict[str, dict[str, object]] = {}
    if not isinstance(groups, list):
        errors.append("executableGroups must be a list")
        return out
    for group in groups:
        if not isinstance(group, dict):
            errors.append("executableGroups entry must be an object")
            continue
        for key in ("id", "capability", "disposition", "reason"):
            if not str(group.get(key, "")).strip():
                errors.append(f"executable group missing {key}")
        disposition = str(group.get("disposition", ""))
        if disposition not in ALLOWED:
            errors.append(f"invalid executable disposition: {group.get('id')}={disposition}")
        bindings = group.get("bindings", [])
        if not isinstance(bindings, list) or not bindings:
            errors.append(f"executable group missing bindings: {group.get('id')}")
            bindings = []
        paths = group.get("sourcePaths", [])
        if not isinstance(paths, list) or not paths:
            errors.append(f"executable group missing sourcePaths: {group.get('id')}")
            continue
        for raw in paths:
            path = str(raw).strip()
            if not path or path in out:
                errors.append(f"duplicate/empty retired executable mapping: {path}")
                continue
            if disposition == "JUSTIFIED_HISTORICAL_SCOPE_RETIREMENT" and "scope" not in path.lower():
                errors.append(f"historical scope-only retirement used for non-scope path: {path}")
            out[path] = {"bindings": bindings, "disposition": disposition, "group": group.get("id")}
    return out


def main() -> int:
    errors: list[str] = []
    migrations = json.loads(MIGRATIONS.read_text(encoding="utf-8"))
    policy = json.loads(POLICY.read_text(encoding="utf-8"))
    registry_path = ROOT / "governance" / "retired-test-equivalence.json"
    if not registry_path.is_file():
        errors.append("retired-test equivalence registry missing")
        registry: dict[str, object] = {}
    else:
        registry = json.loads(registry_path.read_text(encoding="utf-8"))
    binding_rel = str(registry.get("bindingCatalog", "")).strip()
    binding_path = ROOT / binding_rel if binding_rel else ROOT / "governance" / "retired-test-evidence-bindings.json"
    if not binding_rel or not binding_path.is_file():
        errors.append("retired-test equivalence registry bindingCatalog missing/unresolvable")
        bindings: dict[str, object] = {}
    else:
        binding_doc = json.loads(binding_path.read_text(encoding="utf-8"))
        if binding_doc.get("schema") != "DE.PULSE-RETIRED-TEST-EVIDENCE-BINDINGS-1":
            errors.append("unsupported retired-test evidence-binding schema")
        bindings = binding_doc.get("bindings", {}) if isinstance(binding_doc.get("bindings"), dict) else {}

    baseline = str(policy.get("baselineCommit", "")).strip()
    if not baseline or baseline != migrations.get("baselineCommit") or baseline != registry.get("baselineCommit"):
        errors.append("retired-test baseline must exactly match root/migration/registry baseline")
    retained = {int(v) for v in policy.get("retainedRootVersionMajors", [])}
    if retained != {17, 18}:
        errors.append(f"retained majors must remain exactly v17/v18, got {sorted(retained)}")

    baseline_paths = tracked_paths(baseline) if baseline else []
    current_paths = tracked_paths()
    current_set = set(current_paths)
    retired_go_paths = sorted(p for p in baseline_paths if p.endswith("_test.go") and retired_root(p, retained))
    retired_go = go_ids_by_path(retired_go_paths, baseline)
    baseline_exec_all = {p for p in baseline_paths if executable_candidate(p)}
    retired_exec = sorted(p for p in baseline_exec_all if retired_root(p, retained))
    current_go = go_ids_by_path(current_paths)

    if registry.get("schema") != "DE.PULSE-RETIRED-TEST-EQUIVALENCE-1":
        errors.append("unsupported retired-test equivalence schema")
    expected = registry.get("expected", {})
    actual_counts = {"retiredGoSourceFiles": len(retired_go), "retiredGoIdentities": sum(len(ids) for ids in retired_go.values()), "retiredExecutablePaths": len(retired_exec)}
    if expected != actual_counts:
        errors.append(f"registry expected counts drift: expected={expected} actual={actual_counts}")

    mapped_go = flatten_go_groups(registry.get("goGroups"), errors)
    mapped_exec = flatten_exec_groups(registry.get("executableGroups"), errors)
    if set(mapped_go) != set(retired_go):
        errors.append("retired Go family mapping mismatch missing=" + repr(sorted(set(retired_go) - set(mapped_go))) + " extra=" + repr(sorted(set(mapped_go) - set(retired_go))))
    if set(mapped_exec) != set(retired_exec):
        errors.append("retired executable mapping mismatch missing=" + repr(sorted(set(retired_exec) - set(mapped_exec))) + " extra=" + repr(sorted(set(mapped_exec) - set(retired_exec))))

    referenced_bindings: set[str] = set()
    for path, row in mapped_go.items():
        actual = len(retired_go.get(path, []))
        if row.get("identities") != actual:
            errors.append(f"retired Go identity count mismatch {path}: registry={row.get('identities')} baseline={actual}")
        if path in current_set:
            errors.append(f"retired Go source unexpectedly returned to current tree: {path}")
        referenced_bindings.update(str(x) for x in row.get("bindings", []))
    for path, row in mapped_exec.items():
        if path in current_set:
            errors.append(f"retired executable unexpectedly returned to current tree: {path}")
        referenced_bindings.update(str(x) for x in row.get("bindings", []))

    current_go_count = sum(len(ids) for ids in current_go.values())
    for name in sorted(referenced_bindings):
        binding = bindings.get(name)
        if not isinstance(binding, dict):
            errors.append(f"unknown retired-test evidence binding: {name}")
            continue
        for rel in binding.get("evidencePaths", []):
            if not (ROOT / str(rel)).is_file():
                errors.append(f"retired-test evidence path missing for {name}: {rel}")
        for check in binding.get("executionChecks", []):
            if not isinstance(check, dict):
                errors.append(f"invalid execution check for {name}")
                continue
            rel = str(check.get("path", "")).strip()
            token = str(check.get("contains", ""))
            path = ROOT / rel
            if not rel or not token or not path.is_file():
                errors.append(f"retired-test execution anchor missing/invalid for {name}: {rel}")
                continue
            if token not in path.read_text(encoding="utf-8", errors="replace"):
                errors.append(f"retired-test execution anchor drift for {name}: {rel} lacks {token!r}")
        minimum = int(binding.get("minimumCurrentGoIdentities", 0))
        if minimum and current_go_count < minimum:
            errors.append(f"current Go test identity floor failed for {name}: {current_go_count} < {minimum}")

    for required_source in ("v16_8_performance_test.go", "v16_10_performance_test.go"):
        row = mapped_go.get(required_source, {})
        if row.get("disposition") != "MIGRATED_TO_CAPABILITY_OWNER" or "PERFORMANCE_500_SCALE_CURRENT" not in row.get("bindings", []):
            errors.append(f"performance test family lacks explicit migrated current owner: {required_source}")
    for required_source in ("v16_8_performance_gate.py", "v16_10_performance_gate.py"):
        row = mapped_exec.get(required_source, {})
        if row.get("disposition") != "MIGRATED_TO_CAPABILITY_OWNER" or "PERFORMANCE_500_SCALE_CURRENT" not in row.get("bindings", []):
            errors.append(f"performance gate lacks explicit migrated current owner: {required_source}")

    if "age/version alone never justifies retirement" not in str(registry.get("qualityPolicy", "")).lower():
        errors.append("equivalence registry must explicitly prohibit age/version-only retirement")

    print("DE.PULSE retired-test equivalence")
    print(f"baseline: {baseline}")
    print(f"retired Go source families mapped: {len(mapped_go)}/{len(retired_go)}")
    print(f"retired Go identities covered: {sum(len(ids) for ids in retired_go.values())}")
    print(f"retired executable paths mapped: {len(mapped_exec)}/{len(retired_exec)}")
    print(f"current Go test identities: {current_go_count}")
    print(f"current evidence bindings used: {len(referenced_bindings)}")
    print("performance ceilings: MIGRATED_TO_CAPABILITY_OWNER")
    print("age/version-only retirement: PROHIBITED")
    if errors:
        print("DE.PULSE retired-test equivalence: FAIL", file=sys.stderr)
        for error in errors:
            print(" - " + error, file=sys.stderr)
        return 1
    print("retired source-family coverage: 100%")
    print("current evidence path/execution-anchor validation: PASS")
    print("historical scope-only retirement bounded to scope-named gates: PASS")
    print("DE.PULSE retired-test equivalence: PASS")
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
