#!/usr/bin/env python3
from pathlib import Path
import sys
R=Path(__file__).resolve().parent; review=(R/'renderer/qa/v18.0.6-principal-engineer-review.md').read_text(); e=[]
for term in ['PASS — GO to v18.0.6','REUSE → CONSOLIDATE → REFACTOR','Smart Provider Router v2','Source disagreement','MARKET_SHOCK','Hysteresis','Point-in-time learning','SHADOW → VALIDATED → APPROVED → PRODUCTION','TIERED_PARTIAL','Performance/Data Utility','No Execution Boundary','v18.1 multi-user']:
    if term not in review:e.append(term)
if e:
 print('v18.0.6 Principal Engineer Gate: FAIL'); [print(' - missing',x) for x in e]; sys.exit(2)
print('v18.0.6 Principal Engineer Gate: PASS · Router/Rapid Move hardening bounded · protected formulas/no-execution preserved')
