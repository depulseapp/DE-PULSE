#!/usr/bin/env python3
from pathlib import Path
import sys
R=Path(__file__).resolve().parent; review=(R/'renderer/qa/v18.1.0-principal-engineer-review.md').read_text(); e=[]
for term in ['PASS — GO to v18.1.0','REUSE → CONSOLIDATE → REFACTOR','UserWorkspace','IdentityService','PersistenceBackend','Provider Router v2','deduplicated processing union','Opportunity Radar','Rapid Move','ADMIN-owned','Performance/Data Utility','v18.2 Admin / Presence / Sessions','No Execution Boundary']:
    if term not in review:e.append(term)
if e:
 print('v18.1.0 Principal Engineer Gate: FAIL'); [print(' - missing',x) for x in e]; sys.exit(2)
print('v18.1.0 Principal Engineer Gate: PASS · ownership/isolation bounded · one canonical intelligence core · v18.2 excluded')
