#!/usr/bin/env python3
"""G12 randomized-order regression gate with reproducible plan-supplied seeds.

The certification plan owns the exact release seed set. This runner deliberately
does not hard-code a release/version whitelist: any positive bounded integer seed
can be checkpointed and reproduced, while omitting ``--seed`` runs the current
default x10 set. Heavy suites remain separately blocking.
"""
from __future__ import annotations

import argparse
import subprocess
import sys
from pathlib import Path

ROOT = Path(__file__).resolve().parent
SEEDS = [1670701, 1670702, 1670703, 1670704, 1670705,
         1670706, 1670707, 1670708, 1670709, 1670710]
HEAVY_SKIP = (
    r"^(TestExtreme|TestV14StressSixtyMovingSymbolsNoDroppedCanonicalState|"
    r"TestV1512StopTimeoutRemainsStoppingUntilWorkersExit|"
    r"TestV1604StopTimeoutPersistsCurrentCacheBeforeReturning)$"
)

ap = argparse.ArgumentParser()
ap.add_argument(
    "--seed",
    type=int,
    help="run one reproducible checkpointable positive seed; omit to run the current default x10 set",
)
args = ap.parse_args()
if args.seed is not None and not (1 <= args.seed <= 2_147_483_647):
    ap.error("--seed must be a positive 32-bit integer")
run_seeds = [args.seed] if args.seed is not None else SEEDS

for i, seed in enumerate(run_seeds, 1):
    label = f"seed={seed}" if args.seed is not None else f"run {i}/{len(run_seeds)} seed={seed}"
    print(f"Randomized Test Order Gate: {label}", flush=True)
    try:
        p = subprocess.run(
            ["go", "test", "-count=1", f"-shuffle={seed}", "-timeout=60s", "-skip", HEAVY_SKIP, "."],
            cwd=ROOT,
            text=True,
            capture_output=True,
            timeout=75,
        )
    except subprocess.TimeoutExpired as exc:
        print(f"Randomized Test Order Gate: FAIL timeout seed={seed}")
        print((exc.stdout or "")[-5000:])
        print((exc.stderr or "")[-5000:])
        sys.exit(1)
    if p.returncode:
        print(f"Randomized Test Order Gate: FAIL seed={seed}")
        print(p.stdout[-5000:])
        print(p.stderr[-5000:])
        sys.exit(p.returncode)

if args.seed is not None:
    print(
        f"Randomized Test Order Gate: PASS seed={args.seed} · "
        "fast/order-sensitive corpus · dedicated heavy suites remain blocking"
    )
else:
    print(
        "Randomized Test Order Gate: PASS x10 · fast/order-sensitive corpus · "
        "dedicated heavy suites remain blocking · seeds=" + ",".join(map(str, SEEDS))
    )
