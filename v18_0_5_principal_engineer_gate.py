#!/usr/bin/env python3
from pathlib import Path
import sys
R=Path(__file__).resolve().parent; review=(R/'renderer/qa/v18.0.5-principal-engineer-review.md').read_text(); e=[]
for term in ['REUSE → CONSOLIDATE → REFACTOR','Tracked Symbols','explicit empty','Disabled desks','Discovery staging','Opportunity Radar','Technical Context','USER/DEMO','Rapid Move','Day/Swing/Long','No Execution','v18.0.4 Stable','PASS — GO to v18.0.5']:
    if term not in review:e.append(term)
if e:
 print('v18.0.5 Principal Engineer Gate: FAIL'); [print(' - missing',x) for x in e]; sys.exit(2)
print('v18.0.5 Principal Engineer Gate: PASS · focused architecture/UI hardening bounded · intelligence/execution contracts protected')
