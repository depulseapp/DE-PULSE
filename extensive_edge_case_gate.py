#!/usr/bin/env python3
"""Permanent extensive edge-case and failure-mode gate for DE.PULSE."""
from pathlib import Path
import subprocess,sys
ROOT=Path(__file__).resolve().parent
errors=[]
commands=[
 ['go','test','-count=1','-run','TestV160(2|3)Edge','./...'],
 ['node','tests/acceptance/professional_expert_runtime_test.js'],
 [sys.executable,'edge_case_hardening_test.py'],
]
for cmd in commands:
 p=subprocess.run(cmd,cwd=ROOT,text=True,capture_output=True)
 if p.returncode: errors.append('failed: '+' '.join(cmd)+'\n'+p.stdout[-6000:]+'\n'+p.stderr[-6000:])
if errors:
 print('Extensive Edge-Case / Failure-Mode Gate: FAIL');print('\n'.join(errors));sys.exit(1)
print('Extensive Edge-Case / Failure-Mode Gate: PASS · timestamp, stale/missing/conflicting/partial/pagination/restart/runtime cases verified')
