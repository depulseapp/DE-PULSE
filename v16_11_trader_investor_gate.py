#!/usr/bin/env python3
"""v16.11 Professional Trader / Investor real-money acceptance owner.

RL-027/028: this gate owns only closure-specific real-money scenarios/report.
Professional expert, deterministic, replay and recent-Go regressions remain
independent blocking certification checkpoints and are not recursively rerun here.
"""
from pathlib import Path
import subprocess,sys
R=Path(__file__).resolve().parent
errs=[]
try:
    p=subprocess.run(['node','v16_11_real_money_acceptance_test.js'],cwd=R,text=True,capture_output=True,timeout=90)
    if p.returncode: errs.append('fresh real-money renderer scenarios failed:\n'+(p.stdout+p.stderr)[-5000:])
except subprocess.TimeoutExpired:
    errs.append('fresh real-money renderer scenarios timeout')
report=R/'renderer/qa/v16.11.0-trader-investor-review.md'
if not report.exists(): errs.append('trader/investor review report missing')
else:
    txt=report.read_text(errors='ignore')
    for token in ['Real-Money Review','stale live-provider quote','DATA DEGRADED','Wide/risky liquidity','Community/Telegram-style evidence','Closure defect found and fixed','Verdict']:
        if token not in txt: errs.append('review report missing '+token)
# Ensure the closure-specific regression and independent evidence owners remain present.
for f in ['v16_11_real_money_acceptance_test.js','professional_expert_gate.py','deterministic_equivalence_test.js','v16_3_renderer_test.js']:
    if not (R/f).exists(): errs.append('required independent real-money/protection evidence owner missing: '+f)
if errs:
    print('v16.11 Professional Trader / Investor Real-Money Gate: FAIL')
    for e in errs: print(' -',e)
    raise SystemExit(2)
print('v16.11 Professional Trader / Investor Real-Money Gate: PASS · closure-specific freshness/liquidity/event/community/Radar truth; deterministic/professional/replay evidence remains independently blocking')
