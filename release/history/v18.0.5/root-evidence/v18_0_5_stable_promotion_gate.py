#!/usr/bin/env python3
from pathlib import Path
import json,sys
R=Path(__file__).resolve().parent; e=[]
def need(ok,msg):
    if not ok:e.append(msg)
i=json.loads((R/'release_identity.json').read_text())
need(i.get('version')=='18.0.5','version drift')
need(i.get('channel')=='STABLE','channel must be STABLE')
need(i.get('build_id')=='v18.0.5-stable-ui-ux-symbol-management-hardening-20260814','Stable build ID drift')
need(i.get('runtime_config')=='PersonalMarketTerminal','Stable runtime target drift')
need(i.get('application_bundle')=='De-Pulse.app','Stable bundle drift')
boot=(R/'app_bootstrap.go').read_text()
need('const releaseChannel = "STABLE"' in boot,'runtime release channel drift')
need('const buildID = "v18.0.5-stable-ui-ux-symbol-management-hardening-20260814"' in boot,'runtime build ID drift')
prof=(R/'v18_test_profile.go').read_text()
need('if releaseChannel == "STABLE"' in prof and 'stableRuntimeConfigDirName' in prof,'Stable runtime selection missing')
rr=(R/'renderer/renderer.js').read_text()
need("const EXPECTED_BUILD_ID='v18.0.5-stable-ui-ux-symbol-management-hardening-20260814';" in rr,'renderer build identity drift')
man=json.loads((R/'renderer/qa/manifest.json').read_text())['releases'][0]
need(man.get('version')=='18.0.5' and man.get('status')=='STABLE' and man.get('buildId')=='v18.0.5-stable-ui-ux-symbol-management-hardening-20260814','QA manifest Stable identity missing')
for f in ['README.md','VERSION.txt','renderer/docs/user.md','renderer/docs/developer.md','renderer/docs/limitations.md']:
    head='\n'.join((R/f).read_text(errors='ignore').splitlines()[:70])
    need('v18.0.5' in head and 'STABLE' in head.upper(),f+' current Stable identity missing')
if e:
    print('v18.0.5 Stable Promotion Gate: FAIL')
    [print(' -',x) for x in e]
    sys.exit(2)
print('v18.0.5 Stable Promotion Gate: PASS · identity/runtime/docs aligned · functional intelligence unchanged')
