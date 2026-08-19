#!/usr/bin/env python3
from __future__ import annotations

from pathlib import Path
import re
import sys

ROOT = Path(__file__).resolve().parents[2]
WORKFLOWS = ROOT / ".github" / "workflows"
ALLOWED = ("ci-fast.yml", "ci-qualified.yml", "release.yml")
TOP_LEVEL = {"name", "on", "permissions", "concurrency", "env", "jobs"}


def fail(errors: list[str]) -> int:
    print("DE.PULSE workflow structural lint: FAIL", file=sys.stderr)
    for error in errors:
        print(f" - {error}", file=sys.stderr)
    return 1


def github_expressions_balanced(text: str) -> bool:
    """Validate only GitHub expression openers.

    Literal `}}` is valid inside shell/Python/JSON content, so comparing global
    opener/closer counts creates false positives. Every `${{` must instead find
    its own subsequent `}}`; unrelated closing braces are intentionally ignored.
    """
    position = 0
    while True:
        start = text.find("${{", position)
        if start < 0:
            return True
        end = text.find("}}", start + 3)
        if end < 0:
            return False
        nested = text.find("${{", start + 3, end)
        if nested >= 0:
            return False
        position = end + 2


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
    if not github_expressions_balanced(text):
        errors.append(f"{relative}: unbalanced or nested GitHub expression delimiters")

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


def self_test() -> list[str]:
    errors: list[str] = []
    valid = "run: echo '${{ github.sha }}' && python -c 'print({\"x\": {}})'\n"
    missing_close = "run: echo '${{ github.sha'\n"
    nested = "run: echo '${{ github.event && ${{ github.sha }} }}'\n"
    extra_literal_close = "run: python -c 'print({{\"x\":1}})'\n"
    if not github_expressions_balanced(valid):
        errors.append("self-test rejected a valid GitHub expression with literal shell/JSON braces")
    if github_expressions_balanced(missing_close):
        errors.append("self-test accepted an unclosed GitHub expression")
    if github_expressions_balanced(nested):
        errors.append("self-test accepted nested GitHub expression openers")
    if not github_expressions_balanced(extra_literal_close):
        errors.append("self-test rejected unrelated literal closing braces")
    return errors


def main() -> int:
    errors: list[str] = self_test()
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
    print("GitHub-expression parser self-test: PASS")
    print("workflow set / top-level mappings / duplicate job ids: PASS")
    print("indentation / whitespace / GitHub-expression balance: PASS")
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
