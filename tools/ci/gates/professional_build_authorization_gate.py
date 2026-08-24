#!/usr/bin/env python3
"""Checkpoint-aware professional authorization aggregate.

The constituent G10/G11 checks are independent certification-runner processes.
If the wrapper is interrupted, successful checks remain checkpointed and resume
instead of being repeated from the beginning.
"""
from pathlib import Path
import subprocess,sys
ROOT=Path(__file__).resolve().parents[3]
ids=['g10_professional_expert','g10_trader_acceptance','g10_v162_professional','g10_v163_professional_go','g10_v163_replay_renderer','g11_edge_cases','g11_prior_authorization_regressions','g11_fault_injection','g11_fresh_adversarial','g11_v162_scope','g11_v163_scope','g10_v168_professional_go','g10_v168_renderer','g11_v168_scope']
cmd=[sys.executable,'certification_runner.py']
for cid in ids: cmd += ['--check',cid]
p=subprocess.run(cmd,cwd=ROOT,text=True)
if p.returncode:
    print('Professional Build Authorization Gate: FAIL/BLOCKED · see checkpointed certification summary')
    sys.exit(p.returncode)
print('Professional Build Authorization Gate: PASS · professional/trader/current v16.8 + inherited validation/replay/edge/fault/adversarial/scope checks')
