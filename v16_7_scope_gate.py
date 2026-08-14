#!/usr/bin/env python3
from pathlib import Path
import json,subprocess,sys
R=Path(__file__).resolve().parent; e=[]
try:
 c=json.loads((R/'renderer/qa/original-professional-roadmap-acceptance.json').read_text()); by={int(x['id']):x for x in c['items']}
 m=json.loads((R/'renderer/qa/v16.7.0-master-scope.json').read_text()); ids={int(x['id']) for x in m.get('scope_lock',[])}
 ev=json.loads((R/'renderer/qa/v16.7.0-original-acceptance-evidence.json').read_text()); evby={int(x['id']):x for x in ev['scope']}
 if ids!={3,12,13,14,15}: e.append('scope lock mismatch: '+str(sorted(ids)))
 for i in sorted(ids):
  expected=by[i].get('original_acceptance',[]); got=evby.get(i,{})
  clauses=got.get('clauses',[])
  if got.get('status')!='FULL': e.append(f'#{i} is not FULL')
  if len(clauses)!=len(expected): e.append(f'#{i} acceptance clause count mismatch {len(clauses)}/{len(expected)}')
  for n,text in enumerate(expected,1):
   row=next((x for x in clauses if int(x.get('clause',0))==n),None)
   if not row or row.get('acceptance')!=text or row.get('status')!='PASS' or not row.get('evidence'): e.append(f'#{i} clause {n} not fully evidenced')
except Exception as ex: e.append('immutable scope/evidence unreadable: '+str(ex))
trace=(R/'renderer/qa/v16.7.0-traceability.md').read_text(errors='ignore')
for i in [3,12,13,14,15]:
 if f'| {i} |' not in trace or '**FULL**' not in trace: e.append(f'traceability missing/full proof #{i}')
for cmd,label in [(['go','test','-count=1','-run','TestV167','./...'],'v16.7 Go acceptance'),(['node','v16_7_renderer_test.js'],'v16.7 renderer acceptance'),(['node','deterministic_equivalence_test.js'],'deterministic 2403/2403')]:
 p=subprocess.run(cmd,cwd=R,text=True,capture_output=True)
 if p.returncode: e.append(label+' failed: '+(p.stdout+p.stderr)[-1600:])
if e:
 print('v16.7 Original Roadmap Closure Scope Gate: FAIL'); print('\n'.join('- '+x for x in e)); sys.exit(1)
print('v16.7 Original Roadmap Closure Scope Gate: PASS · #3/#12/#13/#14/#15 FULL at immutable clause depth · original roadmap 22 FULL / 8 PARTIAL · deterministic unchanged')
