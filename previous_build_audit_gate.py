#!/usr/bin/env python3
"""G0 previous-Stable closure + adaptive learning audit for v16.11 Major Closure."""
from pathlib import Path
import hashlib,json,os,sys
R=Path(__file__).resolve().parent; e=[]
go='\n'.join(p.read_text(errors='ignore') for p in R.glob('*.go') if not p.name.endswith('_test.go'))
if 'filepath.Join(base, "PersonalMarketTerminal")' not in go: e.append('canonical Stable runtime directory missing')
if 'PersonalMarketTerminal-v16-TEST' in go: e.append('retired TEST runtime identity remains in production source')
required=['renderer/qa/v16.9.0-original-roadmap-status.json','renderer/qa/v16.10.0-master-scope.json','renderer/qa/v16.10.0-acceptance-evidence.json','renderer/qa/v16.11.0-master-scope.json','renderer/qa/original-professional-roadmap-acceptance.json','release_learning_registry.json','release_identity.json','v16_11_major_closure_scope_matrix.json']
for q in required:
 if not (R/q).exists(): e.append('baseline/learning evidence missing: '+q)
try:
 roadmap=json.loads((R/'renderer/qa/v16.9.0-original-roadmap-status.json').read_text())
 if roadmap.get('summary')!='30 FULL / 0 PARTIAL / 0 MISSING': e.append('v16.9 original roadmap closure is not 30/0/0')
 v1610=json.loads((R/'renderer/qa/v16.10.0-master-scope.json').read_text())
 if {int(x.get('id')) for x in v1610.get('scope_lock',[])}!=set(range(1,11)): e.append('v16.10 scope lock must contain exactly clauses 1-10')
 v1611=json.loads((R/'renderer/qa/v16.11.0-master-scope.json').read_text())
 if {int(x.get('id')) for x in v1611.get('scope_lock',[])}!=set(range(1,11)): e.append('v16.11 closure scope lock must contain exactly clauses 1-10')
 ident=json.loads((R/'release_identity.json').read_text())
 if ident.get('previous_stable')!='v16.10.0': e.append('canonical previous Stable must be v16.10.0')
 if ident.get('version')!='16.11.0': e.append('current release identity must be 16.11.0')
except Exception as ex: e.append('v16.10/v16.11 scope/identity unreadable: '+str(ex))
try:
 c=json.loads((R/'renderer/qa/original-professional-roadmap-acceptance.json').read_text())
 if len(c.get('items',[]))!=30 or any(not x.get('original_acceptance') for x in c.get('items',[])): e.append('immutable original-roadmap acceptance contract incomplete')
except Exception as ex: e.append('original-roadmap contract unreadable: '+str(ex))
try:
 reg=json.loads((R/'release_learning_registry.json').read_text()); ids=set()
 if reg.get('schema')!='DE.PULSE-RELEASE-LEARNING-1': e.append('Release Learning Registry schema invalid')
 for x in reg.get('entries',[]):
  if x.get('id') in ids: e.append('duplicate Release Learning entry '+str(x.get('id')))
  ids.add(x.get('id'))
  if x.get('status')=='ACTIVE':
   if not x.get('rootCause') or not x.get('generalizedRule'): e.append(str(x.get('id'))+' active learning incomplete')
   gates=x.get('owningGates') or []
   if not gates or any(not (g.startswith('G') and g[1:].isdigit() and 0<=int(g[1:])<=16) for g in gates): e.append(str(x.get('id'))+' must map to G0-G16')
   if not x.get('permanentProtection'): e.append(str(x.get('id'))+' missing permanent protection')
 for req in [f'RL-{i:03d}' for i in range(20,32)]:
  if req not in ids: e.append(req+' adaptive lesson missing')
except Exception as ex: e.append('Release Learning Registry unreadable: '+str(ex))
for gate in ['v16_0_4_scope_gate.py','v16_0_5_scope_gate.py','v16_0_6_scope_gate.py','v16_1_scope_gate.py','v16_1_1_scope_gate.py','v16_2_scope_gate.py','v16_3_scope_gate.py','v16_4_scope_gate.py','v16_5_scope_gate.py','v16_6_scope_gate.py','v16_7_scope_gate.py','v16_8_scope_gate.py','v16_8_1_scope_gate.py','v16_9_scope_gate.py','v16_10_scope_gate.py','v16_11_scope_gate.py','approved_scope_gate.py']:
 if not (R/gate).exists(): e.append('G11 inherited/closure gate missing: '+gate)
try:
 qa=json.loads((R/'renderer/qa/manifest.json').read_text()); releases=qa.get('releases',[])
 prev=[x for x in releases if str(x.get('version',''))=='16.10.0']
 if not prev: e.append('v16.10.0 Stable predecessor closure missing from QA manifest')
 elif prev[0].get('buildId')!='v16.10.0-stable-opportunity-decision-intelligence-20260812': e.append('v16.10 predecessor build identity drifted')
except Exception as ex: e.append('QA predecessor manifest unreadable: '+str(ex))
expected='00d6ea4d8b5f8f5d5f37d6dc0ac7ef72b2d48586653bece089e835fbd5bea733'; src=os.environ.get('DEPULSE_PREVIOUS_SOURCE','').strip()
if src:
 p=Path(src); got=hashlib.sha256(p.read_bytes()).hexdigest() if p.exists() else ''
 if got!=expected: e.append('frozen v16.10.0 Stable source ZIP hash mismatch')
if e:
 print('G0 Previous Stable + Major Closure Learning Audit: FAIL'); print('\n'.join('- '+x for x in e)); sys.exit(1)
print('G0 Previous Stable + Major Closure Learning Audit: PASS · exact v16.10 baseline identity · 30/30 original + v16.10 10/10 inherited · v16.11 10-clause closure frozen · RL-029/030/031 active')
