#!/usr/bin/env python3
"""Permanent professional expert trader/investor acceptance gate for DE.PULSE."""
from pathlib import Path
import subprocess,sys
ROOT=Path(__file__).resolve().parents[3]
errors=[]
commands=[
 ['go','test','-count=1','-run','TestV160(2|3)Professional','./...'],
 ['node','tests/acceptance/professional_expert_runtime_test.js'],
 ['node','tests/acceptance/trader_acceptance_test.js'],
 ['node','tools/ci/deterministic_equivalence_test.js'],
]
for cmd in commands:
 p=subprocess.run(cmd,cwd=ROOT,text=True,capture_output=True)
 if p.returncode: errors.append('failed: '+' '.join(cmd)+'\n'+p.stdout[-6000:]+'\n'+p.stderr[-6000:])
if errors:
 print('Professional Expert Trader / Investor Acceptance Gate: FAIL');print('\n'.join(errors));sys.exit(1)
print('Professional Expert Trader / Investor Acceptance Gate: PASS · decision truth, explainability, active runtime and deterministic guardrails verified')
