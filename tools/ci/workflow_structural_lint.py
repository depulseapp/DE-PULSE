#!/usr/bin/env python3
from __future__ import annotations

from pathlib import Path
import re
import sys

ROOT = Path(__file__).resolve().parents[2]
WORKFLOWS = ROOT / ".github" / "workflows"
ALLOWED = ("ci-fast.yml", "ci-qualified.yml", "release.yml")
TOP_LEVEL = {"name", "on", "permissions", "concurrency", "env", "jobs"}
MAPPING_RE = re.compile(r"^(\s*)(?:-\s+)?([A-Za-z0-9_.${}\[\]'\" -]+):(?:\s|$)")


def fail(errors: list[str]) -> int:
    print("DE.PULSE workflow structural lint: FAIL", file=sys.stderr)
    for error in errors:
        print(f" - {error}", file=sys.stderr)
    return 1


def lint_file(path: Path) -> list[str]:
    errors: list[str] = []
    text = path.read_text(encoding="utf-8")
    lines = text.splitlines()
    relative = path.relative_to(ROOT).as_posix()

    if not text.endswith("\n"):
        errors.append(f"{relative}: missing final newline")
    if "\t" in text:
        errors.append(f"{relative}: tab indentation is prohibited")
    if any(line.rstrip() != line for line in lines):
        errors.append(f"{relative}: trailing whitespace")
    if text.count("${{") != text.count("}}"):
        errors.append(f"{relative}: unbalanced GitHub expression delimiters")

    top_seen: set[str] = set()
    jobs_line = None
    for number, line in enumerate(lines, 1):
        stripped = line.strip()
        if not stripped or stripped.startswith("#") or stripped in {"|", ">"}:
            continue
        indent = len(line) - len(line.lstrip(" "))
        if indent % 2 != 0:
            errors.append(f"{relative}:{number}: indentation must use two-space levels")
        if indent == 0:
            match = re.match(r"^([A-Za-z0-9_-]+):(?:\s|$)", line)
            if not match:
                errors.append(f"{relative}:{number}: invalid top-level mapping line")
                continue
            key = match.group(1)
            if key not in TOP_LEVEL:
                errors.append(f"{relative}:{number}: unexpected top-level key {key!r}")
            if key in top_seen:
                errors.append(f"{relative}:{number}: duplicate top-level key {key!r}")
            top_seen.add(key)
            if key == "jobs":
                jobs_line = number

    for required in ("name", "on", "permissions", "jobs"):
        if required not in top_seen:
            errors.append(f"{relative}: missing top-level {required!r}")
    if jobs_line is None:
        return errors

    job_keys: set[str] = set()
    in_jobs = False
    for number, line in enumerate(lines, 1):
        if line == "jobs:":
            in_jobs = True
            continue
        if in_jobs and line and not line.startswith(" "):
            in_jobs = False
        if not in_jobs or not line.strip() or line.lstrip().startswith("#"):
            continue
        indent = len(line) - len(line.lstrip(" "))
        if indent == 2:
            match = re.match(r"^  ([A-Za-z0-9_-]+):(?:\s|$)", line)
            if not match:
                errors.append(f"{relative}:{number}: malformed job mapping")
                continue
            key = match.group(1)
            if key in job_keys:
                errors.append(f"{relative}:{number}: duplicate job id {key!r}")
            job_keys.add(key)

    if not job_keys:
        errors.append(f"{relative}: no jobs declared")
    return errors


def main() -> int:
    errors: list[str] = []
    present = sorted(p.name for p in WORKFLOWS.glob("*.y*ml") if p.is_file())
    if present != sorted(ALLOWED):
        errors.append(f"workflow set mismatch: {present}")
    for name in ALLOWED:
        path = WORKFLOWS / name
        if path.is_file():
            errors.extend(lint_file(path))
        else:
            errors.append(f"missing workflow: {name}")
    if errors:
        return fail(errors)
    print("DE.PULSE workflow structural lint: PASS")
    print("workflow set / top-level mappings / duplicate job ids: PASS")
    print("indentation / whitespace / GitHub-expression balance: PASS")
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
