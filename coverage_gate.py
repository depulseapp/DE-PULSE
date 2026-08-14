#!/usr/bin/env python3
"""G12 statement-coverage gate with an explicit regression floor."""
from __future__ import annotations
import argparse, os, re, subprocess, sys
from pathlib import Path

ROOT = Path(__file__).resolve().parent
ap = argparse.ArgumentParser()
default_cert_dir = Path(os.environ.get("DEPULSE_CERT_DIR", str(ROOT.parent / f".depulse-certification-{ROOT.name}"))).resolve()
ap.add_argument("--output", default=str(default_cert_dir / "artifacts" / "coverage.out"))
ap.add_argument("--minimum", type=float, default=50.0)
args = ap.parse_args()
out = Path(args.output)
if not out.is_absolute():
    out = ROOT / out
out.parent.mkdir(parents=True, exist_ok=True)
p = subprocess.run(["go", "test", "-count=1", f"-coverprofile={out}", "./..."], cwd=ROOT, text=True, capture_output=True, timeout=180)
print(p.stdout, end=""); print(p.stderr, end="", file=sys.stderr)
if p.returncode:
    sys.exit(p.returncode)
q = subprocess.run(["go", "tool", "cover", "-func", str(out)], cwd=ROOT, text=True, capture_output=True, timeout=45)
print(q.stdout, end=""); print(q.stderr, end="", file=sys.stderr)
if q.returncode:
    sys.exit(q.returncode)
m = re.search(r"(?m)^total:\s+\(statements\)\s+([0-9.]+)%", q.stdout)
if not m:
    print("Coverage Gate: FAIL unable to parse total statement coverage", file=sys.stderr); sys.exit(1)
value = float(m.group(1))
if value < args.minimum:
    print(f"Coverage Gate: FAIL {value:.1f}% < required {args.minimum:.1f}%", file=sys.stderr); sys.exit(1)
print(f"Coverage Gate: PASS · {value:.1f}% statements · floor {args.minimum:.1f}%")
