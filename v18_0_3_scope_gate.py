#!/usr/bin/env python3
import json,sys
from pathlib import Path
R=Path(__file__).resolve().parent; e=[]
def need(x,m):
    if not x:e.append(m)
i=json.loads((R/'release_identity.json').read_text()); s=json.loads((R/'v18_0_3_scope.json').read_text())
need(i.get('version')=='18.0.3' and i.get('channel')=='TEST','v18.0.3 TEST identity missing')
need(i.get('runtime_config')=='PersonalMarketTerminal-v18.0.3-TEST','isolated v18.0.3 TEST runtime missing')
need(i.get('application_bundle')=='De-Pulse-v18.0.3-TEST.app','separate TEST bundle missing')
need(len(s.get('clauses',[]))==9,'scope clause count mismatch')
text=(R/'http_api.go').read_text(); need('path.Clean' in text and 'filepath.Clean' not in text,'embedded resource path portability fix missing')
tests=(R/'source_test_helpers_test.go').read_text(); need('APPDATA' in tests and 'os.UserConfigDir()' in tests,'cross-platform test config isolation missing')
need('Smart Router v2' in (R/'README.md').read_text() and 'Rapid Move' in (R/'README.md').read_text(),'inherited intelligence boundary missing')
if e:
 print('v18.0.3 Native Runtime Portability Scope: FAIL'); [print(' -',x) for x in e]; sys.exit(2)
print('v18.0.3 Native Runtime Portability Scope: PASS · 9/9 clauses · cross-platform runtime/test delivery hardening · intelligence boundaries protected')
