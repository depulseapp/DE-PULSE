#!/usr/bin/env python3
"""Permanent DE.PULSE Extreme Production release gate.

Ensures all 30 named extreme-test categories remain present, then executes the
Go test matrix. Native target-OS GUI launch and authenticated live-provider
entitlement acceptance remain environment-dependent and are intentionally
reported separately in the release QA report.
"""
from pathlib import Path
import re
import subprocess
import sys

ROOT = Path(__file__).resolve().parent
TEST_FILE = ROOT / "extreme_30_matrix_test.go"
text = TEST_FILE.read_text(encoding="utf-8")
found = {int(x) for x in re.findall(r"func\s+TestExtreme30_(\d{2})", text)}
expected = set(range(1, 31))
missing = sorted(expected - found)
extra = sorted(found - expected)
if missing or extra:
    print(f"Extreme Production gate: FAIL missing={missing} extra={extra}")
    sys.exit(1)

proc = subprocess.run(
    ["go", "test", "-count=1", "-run", r"^TestExtreme30_", "."],
    cwd=ROOT,
    text=True,
    capture_output=True,
)
if proc.returncode != 0:
    sys.stdout.write(proc.stdout)
    sys.stderr.write(proc.stderr)
    print("Extreme Production gate: FAIL")
    sys.exit(proc.returncode)
print("Extreme Production / Professional Trader matrix: 30/30 categories present · PASS")
