#!/usr/bin/env python3
"""v17 Major Closure trader/investor structural evidence owner.
Behavioral tests are independent certification jobs per RL-027.
"""
from pathlib import Path
R=Path(__file__).resolve().parent
errs=[]
report=R/'renderer/qa/v17.5.0-trader-investor-review.md'
if not report.exists(): errs.append('v17.5 trader/investor review report missing')
else:
    txt=report.read_text(errors='ignore')
    for token in ['Real-Money Review','Freshness and readiness truth','OPEN-002 closure','Critical Decision Data','No Execution Boundary','Closure defect found and fixed','Verdict']:
        if token not in txt: errs.append('review report missing '+token)
for f in ['v17_5_real_money_acceptance_test.js','v16_11_real_money_acceptance_test.js','deterministic_equivalence_test.js','professional_expert_gate.py','trader_acceptance_test.js','http_workflow_test.py']:
    if not (R/f).exists(): errs.append('independent real-money evidence owner missing: '+f)
if errs:
    print('v17.5 Professional Trader / Investor Gate: FAIL')
    for e in errs: print(' -',e)
    raise SystemExit(2)
print('v17.5 Professional Trader / Investor Gate: PASS · v17-specific real-money review mapped; behavioral evidence remains independently blocking')
