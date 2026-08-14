#!/usr/bin/env python3
"""Run the full responsive/layout integrity matrix in one bounded browser process.

Repeated Playwright subprocess launches can leave browser/process handles pending on
some certification hosts, making a third shard appear hung after the assertions in
prior shards passed. One full-matrix invocation is both stronger evidence and more
deterministic: all 15 viewports and all surfaces are exercised together under a
single explicit timeout.
"""
from __future__ import annotations
import os, subprocess, sys, time
from pathlib import Path
R=Path(__file__).resolve().parent

env=os.environ.copy()
env.pop('DEPULSE_VIEWPORT_SLICE', None)
started=time.time()
try:
    p=subprocess.run(
        [sys.executable,'responsive_ui_test.py'],
        cwd=R, env=env, text=True, capture_output=True, timeout=240,
    )
    out=(p.stdout or '')+(p.stderr or '')
    print(out.strip(), flush=True)
    if p.returncode != 0:
        raise SystemExit(p.returncode)
except subprocess.TimeoutExpired as e:
    out=(e.stdout or '')+(e.stderr or '')
    print(out.strip(), flush=True)
    print('Responsive/Layout Integrity: INFRA FAIL · full matrix timeout', flush=True)
    raise SystemExit(124)
print(f'Responsive/Layout Integrity: PASS · 15/15 viewports in one bounded full-matrix run · {time.time()-started:.2f}s')
