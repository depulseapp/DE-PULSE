#!/usr/bin/env python3
"""Run one deterministic bounded Go race-detector shard.

Each shard is intentionally a separate process so G12 can checkpoint PASS results
and resume after runner interruption without rerunning already-clean shards.
"""
from __future__ import annotations
import argparse, re, subprocess, sys, time
from pathlib import Path

ROOT = Path(__file__).resolve().parents[3]
ap = argparse.ArgumentParser()
ap.add_argument("--shard", type=int, required=True, help="1-based shard number")
ap.add_argument("--shards", type=int, default=10)
args = ap.parse_args()
if args.shards < 1 or args.shard < 1 or args.shard > args.shards:
    ap.error("--shard must be between 1 and --shards")

q = None
for attempt in (1, 2):
    try:
        q = subprocess.run(["go", "test", "-list", "^Test", "."], cwd=ROOT, text=True, capture_output=True, timeout=45)
        break
    except subprocess.TimeoutExpired as exc:
        print(f"Race shard {args.shard}/{args.shards}: INFRA test-enumeration timeout attempt {attempt}/2", file=sys.stderr)
        print((exc.stdout or "")[-3000:], file=sys.stderr)
        print((exc.stderr or "")[-3000:], file=sys.stderr)
        if attempt == 1:
            time.sleep(1)
if q is None:
    print(f"Race shard {args.shard}/{args.shards}: INFRASTRUCTURE FAIL · repeated test-enumeration timeout; no race test was started", file=sys.stderr)
    sys.exit(2)
if q.returncode:
    print(q.stdout[-6000:]); print(q.stderr[-6000:]); sys.exit(q.returncode)
names = [x.strip() for x in q.stdout.splitlines() if x.startswith("Test")]
shard = names[args.shard - 1 :: args.shards]
if not shard:
    print(f"Race shard {args.shard}/{args.shards}: PASS (empty shard)")
    sys.exit(0)
pattern = "^(" + "|".join(re.escape(x) for x in shard) + ")$"
print(f"Race shard {args.shard}/{args.shards}: tests={len(shard)} discovered_total={len(names)}", flush=True)
p = subprocess.run(
    ["go", "test", "-race", "-count=1", "-timeout=100s", "-run", pattern, "."],
    cwd=ROOT, text=True, capture_output=True, timeout=120,
)
print(p.stdout, end="")
print(p.stderr, end="", file=sys.stderr)
if p.returncode:
    print(f"Race shard {args.shard}/{args.shards}: FAIL", file=sys.stderr)
    sys.exit(p.returncode)
print(f"Race shard {args.shard}/{args.shards}: PASS · {len(shard)} tests")
