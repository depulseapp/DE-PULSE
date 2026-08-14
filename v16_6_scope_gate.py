#!/usr/bin/env python3
from pathlib import Path
import json,subprocess,sys
R=Path(__file__).resolve().parent; e=[]
try:
 m=json.loads((R/'renderer/qa/v16.6.0-master-scope.json').read_text()); ids={int(x['id']) for x in m.get('scope_lock',[])}
 if ids!=set(range(1,31)): e.append(f'professional reconciliation IDs mismatch: {sorted(ids)}')
 if not any(x.get('id')=='V16.6-MS-01' for x in m.get('defect_fixes',[])): e.append('Master Symbol defect scope missing')
except Exception as ex: e.append('master scope unreadable: '+str(ex))
trace=(R/'renderer/qa/v16.6.0-traceability.md').read_text(errors='ignore')
for i in range(1,31):
 if f'| {i} |' not in trace: e.append(f'traceability missing professional item {i}')
for cmd,label in [(['go','test','-count=1','-run','TestV166','./...'],'v16.6 integration/defect regressions'),(['node','v16_6_renderer_test.js'],'v16.6 renderer integration'),(['node','deterministic_equivalence_test.js'],'deterministic 2403/2403')]:
 p=subprocess.run(cmd,cwd=R,text=True,capture_output=True)
 if p.returncode: e.append(label+' failed: '+(p.stdout+p.stderr)[-1200:])
if e:
 print('v16.6 Full Professional Integration Scope Gate: FAIL'); print('\n'.join('- '+x for x in e)); sys.exit(1)
print('v16.6 Full Professional Integration Scope Gate: PASS · 30/30 reconciled · V16.6-MS-01 closed · deterministic unchanged')
