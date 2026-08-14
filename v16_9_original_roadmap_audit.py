#!/usr/bin/env python3
"""Independent original professional-roadmap audit for v16.9.
Validates clause evidence and executes each unique regression family once per source.
Designed to run unchanged from the Exact Source ZIP after packaging.
"""
from pathlib import Path
import concurrent.futures,json,subprocess,sys
R=Path(__file__).resolve().parent; errors=[]

def validate_evidence(version, expected_ids, contract_by):
    try:
        ev=json.loads((R/f'renderer/qa/{version}-original-acceptance-evidence.json').read_text())
        evby={int(x['id']):x for x in ev.get('scope',[])}
        if set(evby)!=set(expected_ids): errors.append(f'{version} evidence IDs mismatch')
        for i in expected_ids:
            exp=contract_by[i].get('original_acceptance',[]); got=evby.get(i,{}); clauses=got.get('clauses',[])
            if got.get('status')!='FULL' or len(clauses)!=len(exp): errors.append(f'{version} #{i} evidence incomplete'); continue
            for n,text in enumerate(exp,1):
                row=next((x for x in clauses if int(x.get('clause',0))==n),None)
                if not row or row.get('acceptance')!=text or row.get('status')!='PASS' or not row.get('evidence'):
                    errors.append(f'{version} #{i} clause {n} not fully evidenced')
    except Exception as ex: errors.append(f'{version} evidence unreadable: {ex}')

def run(item):
    cmd,label=item
    p=subprocess.run(cmd,cwd=R,text=True,capture_output=True)
    return label,p.returncode,(p.stdout+p.stderr)[-3500:]

try:
    contract=json.loads((R/'renderer/qa/original-professional-roadmap-acceptance.json').read_text())
    if len(contract.get('items',[]))!=30: errors.append('immutable contract is not 30 items')
    byc={int(x['id']):x for x in contract.get('items',[])}
    status=json.loads((R/'renderer/qa/v16.9.0-original-roadmap-status.json').read_text())
    by={int(x['id']):x for x in status.get('items',[])}
    if set(by)!=set(range(1,31)): errors.append('roadmap status IDs are not exactly 1..30')
    for i in range(1,31):
        if by.get(i,{}).get('status')!='FULL': errors.append(f'#{i} not FULL')
    if status.get('summary')!='30 FULL / 0 PARTIAL / 0 MISSING': errors.append('summary mismatch')
    m6=json.loads((R/'renderer/qa/v16.6.0-master-scope.json').read_text())
    if {int(x['id']) for x in m6.get('scope_lock',[])}!=set(range(1,31)): errors.append('v16.6 integrated 30-item baseline incomplete')
    validate_evidence('v16.7.0',[3,12,13,14,15],byc)
    validate_evidence('v16.8.0',[6,8,9,21,27],byc)
    validate_evidence('v16.9.0',[10,11,20],byc)
except Exception as ex: errors.append('contract/status unreadable: '+str(ex))

jobs=[
 (['go','test','-count=1','-run','TestV166','./...'],'v16.6 integrated professional baseline'),
 (['go','test','-count=1','-run','TestV167','./...'],'v16.7 clause regression'),
 (['go','test','-count=1','-run','TestV168','./...'],'v16.8 clause regression'),
 (['go','test','-count=1','-run','^TestV169','./...'],'v16.9 final clause regression'),
 (['bash','-lc','node v16_6_renderer_test.js && node v16_7_renderer_test.js && node v16_8_renderer_test.js && node v16_9_renderer_test.js'],'renderer integration chain'),
 (['node','v16_3_renderer_test.js'],'replay no-lookahead/state isolation'),
 (['node','deterministic_equivalence_test.js'],'deterministic equivalence 2403/2403'),
]
with concurrent.futures.ThreadPoolExecutor(max_workers=3) as ex:
    for label,rc,out in [f.result() for f in [ex.submit(run,j) for j in jobs]]:
        if rc: errors.append(label+' failed:\n'+out)
if errors:
    print('INDEPENDENT ORIGINAL ROADMAP AUDIT: FAIL'); [print(' -',x) for x in errors]; raise SystemExit(2)
print('INDEPENDENT ORIGINAL ROADMAP AUDIT: PASS · 30 FULL / 0 PARTIAL / 0 MISSING · unique regression evidence executed once')
