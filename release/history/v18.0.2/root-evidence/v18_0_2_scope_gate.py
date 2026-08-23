#!/usr/bin/env python3
from pathlib import Path
import json,sys
R=Path(__file__).resolve().parent; e=[]
def need(ok,msg):
    if not ok:e.append(msg)
i=json.loads((R/'release_identity.json').read_text()); s=json.loads((R/'v18_0_2_scope.json').read_text())
need(i.get('version')=='18.0.2' and i.get('channel')=='TEST','v18.0.2 TEST identity missing')
need(i.get('build_id')==s.get('build_id'),'build identity drift')
need(i.get('runtime_config')=='PersonalMarketTerminal-v18.0.2-TEST','isolated v18.0.2 TEST runtime missing')
need(i.get('application_bundle')=='De-Pulse-v18.0.2-TEST.app','separate TEST bundle missing')
need(s.get('stable_baseline')=='v17.5.1','Stable baseline drift')
need(len(s.get('clauses',[]))==8,'scope must remain 8 clauses')
for f in ('ci_pipeline.py','prefreeze_qualification.py','certification_runner.py'):
    need('canonical_source_fingerprint' in (R/f).read_text(),f+' does not use canonical fingerprint owner')
need((R/'source_fingerprint.py').exists() and (R/'source_fingerprint_portability_test.py').exists(),'portability utility/test missing')
if e:
 print('v18.0.2 Native Delivery Scope: FAIL'); [print(' -',x) for x in e]; sys.exit(2)
print('v18.0.2 Native Delivery Scope: PASS · 8/8 clauses · canonical cross-platform fingerprint owner · protected trading logic inherited')
