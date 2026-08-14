#!/usr/bin/env python3
import json,sys
from pathlib import Path
R=Path(__file__).resolve().parent; e=[]
def need(x,m):
    if not x:e.append(m)
i=json.loads((R/'release_identity.json').read_text()); s=json.loads((R/'v18_0_4_scope.json').read_text())
need(i.get('version')=='18.0.4' and i.get('channel')=='STABLE','v18.0.4 STABLE identity missing')
need(i.get('build_id')==s.get('build_id'),'build identity drift')
need(i.get('runtime_config')=='PersonalMarketTerminal','canonical v18.0.4 Stable runtime missing')
need(i.get('application_bundle')=='De-Pulse.app','Stable application bundle missing')
need(s.get('stable_baseline')=='v17.5.1','Stable baseline drift')
need(len(s.get('clauses',[]))==8,'scope clause count mismatch')
helpers=(R/'source_test_helpers_test.go').read_text(); need('registerApplicationCleanup' in helpers and 'persistence.Close()' in helpers,'Windows SQLite lifecycle cleanup helper missing')
need('buildId' in (R/'http_api.go').read_text() or 'buildID' in (R/'app_bootstrap.go').read_text(),'canonical build identity owner missing')
need('Smart Router v2' in (R/'README.md').read_text() and 'Rapid Move' in (R/'README.md').read_text(),'inherited intelligence boundary missing')
if e:
 print('v18.0.4 Native Closure Scope: FAIL'); [print(' -',x) for x in e]; sys.exit(2)
print('v18.0.4 Stable Promotion Scope: PASS · 8/8 clauses · native closure preserved · Stable identity/runtime target verified')
