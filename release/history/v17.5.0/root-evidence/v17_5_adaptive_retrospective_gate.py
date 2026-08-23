#!/usr/bin/env python3
from pathlib import Path
import json
R=Path(__file__).resolve().parent
errs=[]
d=json.loads((R/'release_learning_registry.json').read_text())
entries={e.get('id'):e for e in d.get('entries',[])}
for rid in ['RL-032','RL-033','RL-034','RL-035','RL-036','RL-037']:
    e=entries.get(rid)
    if not e: errs.append(rid+' missing')
    else:
        if e.get('status')!='ACTIVE': errs.append(rid+' not ACTIVE')
        if not e.get('owningGates') or any(not str(g).startswith('G') for g in e.get('owningGates',[])): errs.append(rid+' missing existing gate ownership')
        if any(int(str(g)[1:])>16 for g in e.get('owningGates',[]) if str(g)[1:].isdigit()): errs.append(rid+' introduced gate beyond G16')
        if not e.get('permanentProtection'): errs.append(rid+' missing permanent protection')
report=R/'renderer/qa/v17.5.0-adaptive-retrospective.md'
if not report.exists(): errs.append('adaptive retrospective report missing')
else:
    txt=report.read_text(errors='ignore')
    for token in ['RL-032','RL-033','RL-034','RL-035','RL-036','RL-037','G0–G16','Retrospective verdict']:
        if token not in txt: errs.append('retrospective missing '+token)
if errs:
    print('v17.5 Adaptive Retrospective Gate: FAIL')
    for e in errs: print(' -',e)
    raise SystemExit(2)
print('v17.5 Adaptive Retrospective Gate: PASS · Release Learning active through RL-037 · G0-G16 ownership preserved')
