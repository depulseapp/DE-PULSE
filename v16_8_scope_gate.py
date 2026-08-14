#!/usr/bin/env python3
from pathlib import Path
import json,subprocess,sys
R=Path(__file__).resolve().parent; e=[]
try:
 c=json.loads((R/'renderer/qa/original-professional-roadmap-acceptance.json').read_text()); by={int(x['id']):x for x in c['items']}
 m=json.loads((R/'renderer/qa/v16.8.0-master-scope.json').read_text()); ids={int(x['id']) for x in m.get('scope_lock',[])}
 ev=json.loads((R/'renderer/qa/v16.8.0-original-acceptance-evidence.json').read_text()); evby={int(x['id']):x for x in ev['scope']}
 if ids!={6,8,9,21,27}: e.append('scope lock mismatch: '+str(sorted(ids)))
 for i in sorted(ids):
  expected=by[i].get('original_acceptance',[]); got=evby.get(i,{})
  clauses=got.get('clauses',[])
  if got.get('status')!='FULL': e.append(f'#{i} is not FULL')
  if len(clauses)!=len(expected): e.append(f'#{i} acceptance clause count mismatch {len(clauses)}/{len(expected)}')
  for n,text in enumerate(expected,1):
   row=next((x for x in clauses if int(x.get('clause',0))==n),None)
   if not row or row.get('acceptance')!=text or row.get('status')!='PASS' or not row.get('evidence'): e.append(f'#{i} clause {n} not fully evidenced')
 status=json.loads((R/'renderer/qa/v16.8.0-original-roadmap-status.json').read_text())
 counts={k:sum(1 for x in status['items'] if x['status']==k) for k in ['FULL','PARTIAL','MISSING']}
 if counts!={'FULL':27,'PARTIAL':3,'MISSING':0}: e.append('roadmap status mismatch '+str(counts))
except Exception as ex: e.append('immutable scope/evidence unreadable: '+str(ex))
trace=(R/'renderer/qa/v16.8.0-traceability.md').read_text(errors='ignore')
for i in [6,8,9,21,27]:
 if f'| {i} |' not in trace: e.append(f'traceability missing #{i}')
for cmd,label in [(['go','test','-count=1','-run','TestV168','./...'],'v16.8 Go acceptance'),(['node','v16_8_renderer_test.js'],'v16.8 renderer acceptance'),(['node','deterministic_equivalence_test.js'],'deterministic 2403/2403')]:
 p=subprocess.run(cmd,cwd=R,text=True,capture_output=True)
 if p.returncode: e.append(label+' failed: '+(p.stdout+p.stderr)[-2000:])
if e:
 print('v16.8 Original Roadmap Closure Scope Gate: FAIL'); print('\n'.join('- '+x for x in e)); sys.exit(1)
print('v16.8 Original Roadmap Closure Scope Gate: PASS · #6/#8/#9/#21/#27 FULL · roadmap 27 FULL / 3 PARTIAL / 0 MISSING · deterministic unchanged')
