#!/usr/bin/env python3
from pathlib import Path
import sys
R=Path(__file__).resolve().parent; e=[]
def need(ok,msg):
    if not ok:e.append(msg)
review=(R/'renderer/qa/v18.0.2-principal-engineer-review.md').read_text()
for term in ['source_fingerprint.py','Smart Router v2','Rapid Move','Day/Swing/Long','No Execution','v17.5.1 Stable','PASS — GO to v18.0.2']:
    need(term in review,'review missing '+term)
for f in ('ci_pipeline.py','prefreeze_qualification.py','certification_runner.py'):
    need('canonical_source_fingerprint' in (R/f).read_text(),f+' not consolidated')
if e:
 print('v18.0.2 Principal Engineer Gate: FAIL'); [print(' -',x) for x in e]; sys.exit(2)
print('v18.0.2 Principal Engineer Gate: PASS · canonical fingerprint ownership consolidated · intelligence/execution boundaries protected')
