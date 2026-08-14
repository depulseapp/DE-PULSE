#!/usr/bin/env python3
from pathlib import Path
import json,sys
R=Path(__file__).resolve().parent; errors=[]
def need(ok,msg):
    if not ok: errors.append(msg)
try:
    frozen=json.loads((R/'v17_major_scope.json').read_text()); matrix=json.loads((R/'v17_5_major_closure_scope_matrix.json').read_text()); slices=json.loads((R/'v17_delivery_slices.json').read_text()); ident=json.loads((R/'release_identity.json').read_text())
    native=ident.get('version') in {'17.5.0','17.5.1'} and ident.get('channel') in {'TEST','RC','STABLE'}
    inherited=str(ident.get('version','')).startswith('18.') and ident.get('channel') in {'TEST','RC','STABLE'} and ident.get('stable_baseline')=='v17.5.1'
    need(native or inherited,'v17.5 closure/inherited release identity missing')
    if native: need(ident.get('stable_baseline')=='v16.11.0','v16.11 Stable baseline drifted')
    rows=matrix.get('items',[])
    need(len(rows)==20 and [x.get('item') for x in rows]==frozen.get('items'),'closure matrix must exactly reconstruct frozen 20-item v17 scope')
    for row in rows:
        need(bool(row.get('canonical_owner')),f"#{row.get('id')} canonical owner missing")
        for f in row.get('code_owners',[]): need((R/f).exists(),f"#{row.get('id')} code owner missing: {f}")
        ev=row.get('fresh_executable_evidence',[]); need(bool(ev),f"#{row.get('id')} fresh evidence owner missing")
        for f in ev: need((R/f).exists(),f"#{row.get('id')} evidence missing: {f}")
    by={x.get('id'):x for x in slices.get('slices',[])}
    for sid in ['v17.0','v17.1','v17.2','v17.3','v17.4']:
        need(str(by.get(sid,{}).get('status','')).startswith('SOURCE TEST CHECKPOINT') or by.get(sid,{}).get('status') in {'TEST QUALIFIED','SOURCE QUALIFIED','COMPLETE'},sid+' is not a completed source checkpoint')
    need(by.get('v17-major-closure',{}).get('status') in {'IN DEVELOPMENT','PREFREEZE QUALIFIED','RC','STABLE'},'v17 Major Closure slice not active')
    for row in matrix.get('slice_milestones',[]):
        need((R/row['scope_gate']).exists(),row['id']+' scope gate missing')
    need(len(matrix.get('permanent_contracts',[]))>=6,'permanent contract reconstruction incomplete')
except Exception as ex:
    errors.append('matrix unreadable: '+str(ex))
if errors:
    print('v17.5 Major Closure Scope Matrix: FAIL'); [print(' -',e) for e in errors]; sys.exit(2)
print('v17.5 Major Closure Scope Matrix: PASS · frozen 20/20 v17 scope + v17.0-v17.4 slices + permanent boundaries mapped to current owners/evidence')
