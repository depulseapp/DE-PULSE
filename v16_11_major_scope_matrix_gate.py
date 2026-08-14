#!/usr/bin/env python3
"""v16.11 Major Closure matrix/code-truth gate.
Static/structural by design: executable behavior is owned by independent fresh
checkpoints so this gate never recursively nests expensive historical audits.
"""
from pathlib import Path
import json,sys
R=Path(__file__).resolve().parent
errs=[]
try:
    m=json.loads((R/'v16_11_major_closure_scope_matrix.json').read_text())
    c=json.loads((R/'renderer/qa/original-professional-roadmap-acceptance.json').read_text())
    rows=m.get('original_professional_roadmap',[])
    if len(rows)!=30 or {int(x['id']) for x in rows}!=set(range(1,31)):
        errs.append('closure matrix must contain exactly original roadmap IDs 1..30')
    cb={int(x['id']):x for x in c.get('items',[])}
    for row in rows:
        i=int(row['id']); exp=cb.get(i,{})
        if row.get('name')!=exp.get('name'): errs.append(f'#{i} name drift')
        if row.get('acceptance')!=exp.get('original_acceptance'): errs.append(f'#{i} immutable acceptance drift')
        if not row.get('canonical_owner'): errs.append(f'#{i} canonical owner missing')
        owners=row.get('code_owners',[])
        if not owners: errs.append(f'#{i} code owners missing')
        for f in owners:
            if not (R/f).exists(): errs.append(f'#{i} current code owner missing: {f}')
        ev=row.get('fresh_executable_evidence',[])
        if not ev: errs.append(f'#{i} executable evidence owner missing')
        for f in ev:
            if not (R/f).exists(): errs.append(f'#{i} evidence file missing: {f}')
    ms=m.get('major_family_milestones',[])
    expected={'16.1.0':7,'16.1.1':3,'16.2.0':7,'16.3.0':5,'16.4.0':3,'16.5.0':5,'16.6.0':30,'16.7.0':5,'16.8.0':5,'16.8.1':6,'16.9.0':3,'16.10.0':10}
    if {x.get('version') for x in ms}!=set(expected): errs.append('major-family milestone coverage mismatch')
    for x in ms:
        v=x['version']
        if int(x.get('scope_count',-1))!=expected[v]: errs.append(v+' scope count drift')
        for f in (x.get('scope_contract'),x.get('fresh_gate')):
            if not f or not (R/f).exists(): errs.append(v+' evidence owner missing '+str(f))
    if len(m.get('permanent_contracts',[]))<10: errs.append('permanent contract coverage incomplete')
except Exception as ex:
    errs.append('matrix unreadable: '+str(ex))
if errs:
    print('v16.11 Major Closure Scope Matrix: FAIL')
    for e in errs: print(' -',e)
    raise SystemExit(2)
print('v16.11 Major Closure Scope Matrix: PASS · 30 original items + 12 v16 milestones + permanent contracts mapped to current code/evidence')
