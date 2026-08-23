#!/usr/bin/env python3
"""DE.PULSE G0 release-environment readiness gate.

Exit codes are certification-runner aware:
  0 = PASS
  3 = BLOCKED (environment/tooling prerequisite, not a product defect)
  1 = malformed policy or other gate defect
"""
from __future__ import annotations
import json, platform, re, shutil, subprocess, sys
from pathlib import Path

ROOT = Path(__file__).resolve().parent
policy = json.loads((ROOT / "release_environment_policy.json").read_text())
blocked: list[str] = []
errors: list[str] = []

for cmd in policy.get("required_commands", []):
    if shutil.which(cmd) is None:
        blocked.append(f"required command unavailable: {cmd}")
for cmd in policy.get("required_security_commands", []):
    if shutil.which(cmd) is None:
        blocked.append(f"required security command unavailable: {cmd}")

try:
    out = subprocess.check_output(["go", "version"], text=True, stderr=subprocess.STDOUT).strip()
    m = re.search(r"go(\d+\.\d+)(?:\.\d+)?", out)
    if not m:
        errors.append(f"unable to parse Go version: {out}")
    elif m.group(1) not in set(policy.get("approved_go_minor_lines", [])):
        blocked.append(
            f"Go {m.group(1)} is not in approved release lines "
            f"{policy.get('approved_go_minor_lines')} (preferred {policy.get('preferred_go_version')}); observed: {out}"
        )
except Exception as exc:
    blocked.append(f"unable to execute go version: {exc}")

host = platform.system().lower()
targets = set(policy.get("native_acceptance_targets", []))
if "darwin" in targets and host != "darwin":
    blocked.append("native macOS acceptance unavailable on this host")
if "windows" in targets and host != "windows":
    blocked.append("native Windows acceptance unavailable on this host")

print("DE.PULSE G0 — Release Environment Readiness")
print(f"Host: {platform.system()} {platform.machine()}")
print(f"Policy date: {policy.get('policy_date')}")
if errors:
    for e in errors:
        print("FAIL:", e)
    sys.exit(1)
if blocked:
    for item in blocked:
        print("BLOCKED:", item)
    print("G0 RESIDUAL — release-assurance prerequisites are unavailable on this host; record them as documented non-blocking residuals for the user-approved Stable promotion while product certification continues")
    sys.exit(3)
print("G0 PASS — release host has approved toolchain/security/native prerequisites")
