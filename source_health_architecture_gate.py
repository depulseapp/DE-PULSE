#!/usr/bin/env python3
"""DE.PULSE G2 — recursive whole-source health, architecture fit and reuse audit."""
from __future__ import annotations

from collections import Counter, defaultdict
from pathlib import Path
import hashlib
import json
import re
import sys

ROOT = Path(__file__).resolve().parent
HEALTH_POLICY = json.loads((ROOT / "source_health_baseline.json").read_text())
FILE_MAX = int(HEALTH_POLICY.get("production_file_max_lines", 1200))
EXCLUDED_DIRS = {".git", "vendor", "node_modules", "testdata", ".depulse-certification"}


def rel(path: Path) -> str:
    return path.relative_to(ROOT).as_posix()


def production_go(path: Path) -> bool:
    if not path.is_file() or path.suffix != ".go" or path.name.endswith("_test.go"):
        return False
    relative = path.relative_to(ROOT)
    return not any(part in EXCLUDED_DIRS for part in relative.parts)


PROD = sorted((p for p in ROOT.rglob("*.go") if production_go(p)), key=rel)
errors: list[str] = []


def fail(message: str) -> None:
    errors.append(message)


prod_text = "\n".join(p.read_text(errors="ignore") for p in PROD)
# Preserve the historical whole-corpus reference semantics while making discovery
# recursive. This prevents package moves from silently disappearing from G2.
prod_identifier_counts = Counter(re.findall(r"\b[A-Za-z_]\w*\b", prod_text))

# Maintenance debt and module budget across every production Go package.
for p in PROD:
    lines = p.read_text(errors="ignore").splitlines()
    if len(lines) > FILE_MAX:
        fail(f"{rel(p)}: {len(lines)} lines exceeds {FILE_MAX:,}-line responsibility budget")
    for i, line in enumerate(lines, 1):
        if re.search(r"\b(TODO|FIXME|HACK)\b", line):
            fail(f"{rel(p)}:{i}: unresolved maintenance marker")

# Unreferenced production functions/methods using the conserved whole-corpus
# name-reference rule. Exported cross-package users remain visible in the corpus.
func_decl = re.compile(r"(?m)^func\s+(?:\([^\n]*?\)\s*)?([A-Za-z_]\w*)\s*\(")
for p in PROD:
    for match in func_decl.finditer(p.read_text(errors="ignore")):
        name = match.group(1)
        if name in {"main", "init"}:
            continue
        if prod_identifier_counts.get(name, 0) == 1:
            fail(f"{rel(p)}: production helper {name} has no production reference")

# Exact duplicate sizeable Go function bodies across the recursive production tree.
def bodies(path: Path):
    source = path.read_text(errors="ignore")
    head = re.compile(r"(?m)^func\s+(?:\([^\n]*?\)\s*)?([A-Za-z_]\w*)\s*\([^\n]*?\)[^{]*\{")
    for match in head.finditer(source):
        i = match.end() - 1
        depth = 0
        j = i
        quote = None
        escaped = False
        line_comment = False
        block_comment = False
        while j < len(source):
            char = source[j]
            nxt = source[j + 1] if j + 1 < len(source) else ""
            if line_comment:
                if char == "\n":
                    line_comment = False
            elif block_comment:
                if char == "*" and nxt == "/":
                    block_comment = False
                    j += 1
            elif quote:
                if quote == "`":
                    if char == "`":
                        quote = None
                elif escaped:
                    escaped = False
                elif char == "\\":
                    escaped = True
                elif char == quote:
                    quote = None
            else:
                if char == "/" and nxt == "/":
                    line_comment = True
                    j += 1
                elif char == "/" and nxt == "*":
                    block_comment = True
                    j += 1
                elif char in {'"', "'", "`"}:
                    quote = char
                elif char == "{":
                    depth += 1
                elif char == "}":
                    depth -= 1
                    if depth == 0:
                        body = source[i + 1 : j]
                        body = re.sub(r"//.*", "", body)
                        body = re.sub(r"/\*.*?\*/", "", body, flags=re.S)
                        body = re.sub(r"\s+", " ", body).strip()
                        if len(body) >= 180:
                            yield match.group(1), body
                        break
            j += 1


duplicates = defaultdict(list)
for p in PROD:
    for name, body in bodies(p):
        duplicates[hashlib.sha256(body.encode()).hexdigest()].append((rel(p), name))
for rows in duplicates.values():
    if len(rows) > 1:
        rowset = set(rows)
        platform_parity = {
            ("persistence_backend_sqlite.go", "Stats"),
            ("persistence_backend_windows.go", "Stats"),
        }
        if rowset == platform_parity:
            continue
        fail("exact duplicate production function bodies: " + ", ".join(f"{p}:{n}" for p, n in rows))

# Canonical release/current-state truth. The stale v18.6 ci_pipeline_plan.json is
# deliberately no longer a G2 current-authority dependency; #70 will retire that
# competing orchestration only after consumer/evidence equivalence is proven.
for required in (
    "release_identity.json",
    "data_utility_registry.json",
    "data_health_policy.json",
    "prefreeze_qualification.py",
    "governance/current-state.json",
    "governance/work-slices/ADAPT-CI-CONVERGENCE-001/work-slice.json",
):
    if not (ROOT / required).exists():
        fail("Adaptive Build Process artifact missing: " + required)
try:
    ident = json.loads((ROOT / "release_identity.json").read_text())
    state = json.loads((ROOT / "governance" / "current-state.json").read_text())
    stable = state.get("stable", {})
    active = state.get("activeWorkSlice", {})
    if stable.get("productVersion") != ident.get("version"):
        fail("current-state Stable productVersion / release identity drift")
    if stable.get("buildId") != ident.get("build_id"):
        fail("current-state Stable buildId / release identity drift")
    if stable.get("platformBuildNumber") != ident.get("bundle_version"):
        fail("current-state platform build number / release identity drift")
    if active.get("workSliceId") != "ADAPT-CI-CONVERGENCE-001":
        fail("current-state active work slice is not ADAPT-CI-CONVERGENCE-001")
    if active.get("publicProductVersion") is not None:
        fail("#70 process work slice must not consume a public product version")
except Exception as exc:
    fail("canonical release/current-state validation failed: " + str(exc))

# Renderer orphan/duplicate checks.
js = (ROOT / "renderer" / "renderer.js").read_text(errors="ignore")
js_identifier_counts = Counter(re.findall(r"(?<![A-Za-z0-9_$])[A-Za-z_$][\w$]*(?![A-Za-z0-9_$])", js))
js_decl = re.compile(r"(?m)^function\s+([A-Za-z_$][\w$]*)\s*\(")
for match in js_decl.finditer(js):
    name = match.group(1)
    if js_identifier_counts.get(name, 0) == 1:
        fail(f"renderer/renderer.js: named renderer helper {name} has no production reference")

# Canonical ownership / retired paths.
for token, why in {
    "refreshTwelveFundamentalFallback": "retired Twelve fundamentals fallback",
    '"twelve-fundamentals"': "retired Twelve fundamentals health key",
    "/institutional_holders?symbol=": "retired Twelve institutional-holder fallback",
}.items():
    if token in prod_text:
        fail(f"{why} returned: {token}")
if not re.search(r'"Fundamentals"\s*:\s*\{"Finnhub",\s*"SEC",\s*"yfinance"\}', prod_text):
    fail("Fundamentals route must remain Finnhub -> SEC -> yfinance")
if len(re.findall(r"func\s+\(e \*Engine\)\s+executeProviderRoute\s*\(", prod_text)) != 1:
    fail("Provider Router must have exactly one executeProviderRoute authority")
if len(re.findall(r"var\s+canonicalSymbolClassifications\s*=", prod_text)) != 1:
    fail("Market Intelligence classification must have exactly one backend canonical owner")
if "const sectorETF=" in js:
    fail("renderer-owned ticker/sector map returned; classification must remain backend-owned")

package_dirs = sorted({rel(p.parent) if p.parent != ROOT else "." for p in PROD})
print("DE.PULSE G2 — Recursive Source Health + Architecture Fit + Reuse Audit")
print(f"Production Go files: {len(PROD)}")
print(f"Production Go packages/directories: {len(package_dirs)}")
print("Production package directories: " + ", ".join(package_dirs))
print(f"Production Go lines: {sum(len(p.read_text(errors='ignore').splitlines()) for p in PROD)}")
print("Orphan production Go helpers: 0" if not any("production helper" in e for e in errors) else "Orphan production Go helpers: FAIL")
print("Orphan named renderer helpers: 0" if not any("renderer helper" in e for e in errors) else "Orphan named renderer helpers: FAIL")
print("Exact duplicate sizeable Go bodies: 0" if not any("duplicate production" in e for e in errors) else "Exact duplicate sizeable Go bodies: FAIL")
print("Canonical current-state owner: governance/current-state.json")
if errors:
    for error in errors:
        print("FAIL:", error)
    sys.exit(1)
print("G2 PASS — recursive production discovery · one canonical owner per responsibility · REUSE/CONSOLIDATE/REFACTOR/DELETE before ADD")
