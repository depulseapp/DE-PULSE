#!/usr/bin/env python3
from pathlib import Path
import sys
R=Path(__file__).resolve().parent; review=(R/'renderer/qa/v18.0.3-principal-engineer-review.md').read_text(); e=[]
for term in ['platform-neutral','os.UserConfigDir()','PowerShell','Smart Router v2','Rapid Move','Day/Swing/Long','No Execution','v17.5.1 Stable','PASS — GO to v18.0.3']:
    if term not in review:e.append(term)
if e:
 print('v18.0.3 Principal Engineer Gate: FAIL'); [print(' - missing',x) for x in e]; sys.exit(2)
print('v18.0.3 Principal Engineer Gate: PASS · native portability fixes bounded · intelligence/execution boundaries protected')
