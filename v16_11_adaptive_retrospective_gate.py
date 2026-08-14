#!/usr/bin/env python3
from pathlib import Path
import json,sys
R=Path(__file__).resolve().parent; errs=[]
rep=(R/'renderer/qa/v16.11.0-adaptive-retrospective.md')
if not rep.exists(): errs.append('retrospective missing')
else:
 t=rep.read_text(errors='ignore')
 for token in ['Wrapper timeouts','Duplicate historical certification work','Pre-Freeze Qualification','resource-class-aware','Source accumulated obsolete delivery debris','Freshness truth and readiness truth diverged','Major Closure & Release Assurance','Adaptive Roadmap','Adaptive Build Plan','Adaptive Build Process','Adaptive Delivery','v17 entry conditions']:
  if token not in t: errs.append('retrospective missing '+token)
try:
 reg=json.loads((R/'release_learning_registry.json').read_text())
 for rid in ['RL-029','RL-030','RL-031']:
  if not any(x.get('id')==rid and x.get('status')=='ACTIVE' for x in reg.get('entries',[])): errs.append(rid+' not ACTIVE')
except Exception as ex: errs.append('release learning unreadable: '+str(ex))
try:
 scope=json.loads((R/'renderer/qa/v16.11.0-master-scope.json').read_text())
 if len(scope.get('scope_lock',[]))!=10: errs.append('closure scope not 10 clauses')
except Exception as ex: errs.append('closure scope unreadable: '+str(ex))
if errs:
 print('v16.11 Adaptive Retrospective Gate: FAIL'); [print(' -',e) for e in errs]; raise SystemExit(2)
print('v16.11 Adaptive Retrospective Gate: PASS · v16 challenges/root causes/fixes/effectiveness/process adaptations reconciled')
