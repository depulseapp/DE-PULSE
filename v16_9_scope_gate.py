#!/usr/bin/env python3
from pathlib import Path
import json,subprocess,sys
R=Path(__file__).resolve().parent; errs=[]
def run(cmd,label):
 p=subprocess.run(cmd,cwd=R,text=True,capture_output=True)
 if p.returncode: errs.append(label+' failed:\n'+(p.stdout+p.stderr)[-3500:])
try:
 c=json.loads((R/'renderer/qa/original-professional-roadmap-acceptance.json').read_text()); by={int(x['id']):x for x in c['items']}
 m=json.loads((R/'renderer/qa/v16.9.0-master-scope.json').read_text()); ids={int(x['id']) for x in m.get('scope_lock',[])}
 ev=json.loads((R/'renderer/qa/v16.9.0-original-acceptance-evidence.json').read_text()); evby={int(x['id']):x for x in ev['scope']}
 if ids!={10,11,20}: errs.append('scope lock mismatch '+str(sorted(ids)))
 for i in sorted(ids):
  expected=by[i].get('original_acceptance',[]); got=evby.get(i,{}); clauses=got.get('clauses',[])
  if got.get('status')!='FULL': errs.append(f'#{i} is not FULL')
  if len(clauses)!=len(expected): errs.append(f'#{i} acceptance clause count mismatch {len(clauses)}/{len(expected)}')
  for n,text in enumerate(expected,1):
   row=next((x for x in clauses if int(x.get('clause',0))==n),None)
   if not row or row.get('acceptance')!=text or row.get('status')!='PASS' or not row.get('evidence'): errs.append(f'#{i} clause {n} not fully evidenced')
 status=json.loads((R/'renderer/qa/v16.9.0-original-roadmap-status.json').read_text())
 counts={k:sum(1 for x in status['items'] if x['status']==k) for k in ['FULL','PARTIAL','MISSING']}
 if counts!={'FULL':30,'PARTIAL':0,'MISSING':0}: errs.append('roadmap status mismatch '+str(counts))
 policy=json.loads((R/'community_source_policy.json').read_text())
 for p in ['X','REDDIT','TELEGRAM','DISCORD','WHATSAPP','MANUAL']:
  if p not in policy.get('platforms',{}): errs.append('community source policy missing '+p)
except Exception as ex: errs.append('scope/evidence unreadable: '+str(ex))
trace=(R/'renderer/qa/v16.9.0-traceability.md').read_text(errors='ignore')
for i in [10,11,20]:
 if f'| {i} |' not in trace: errs.append(f'traceability missing #{i}')
run(['go','test','-count=1','-run','^TestV169','./...'],'v16.9 Go acceptance')
run(['node','v16_9_renderer_test.js'],'v16.9 renderer acceptance')
run(['node','v16_3_renderer_test.js'],'replay no-lookahead/live-state isolation')
run(['node','deterministic_equivalence_test.js'],'deterministic 2403/2403')
if errs:
 print('v16.9 Original Roadmap Closure Scope Gate: FAIL'); [print(' -',e) for e in errs]; raise SystemExit(2)
print('v16.9 Original Roadmap Closure Scope Gate: PASS · #10/#11/#20 FULL · roadmap 30 FULL / 0 PARTIAL / 0 MISSING · deterministic unchanged')
