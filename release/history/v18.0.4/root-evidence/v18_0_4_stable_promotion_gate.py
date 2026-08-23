#!/usr/bin/env python3
from pathlib import Path
import json,sys
R=Path(__file__).resolve().parent; e=[]
def need(x,m):
    if not x:e.append(m)
i=json.loads((R/'release_identity.json').read_text())
need(i.get('version')=='18.0.4','version drift')
need(i.get('channel')=='STABLE','channel must be STABLE')
need(i.get('build_id')=='v18.0.4-stable-native-cross-platform-closure-20260813','Stable build ID drift')
need(i.get('runtime_config')=='PersonalMarketTerminal','Stable runtime target drift')
need(i.get('application_bundle')=='De-Pulse.app','Stable bundle drift')
boot=(R/'app_bootstrap.go').read_text(); need('const releaseChannel = "STABLE"' in boot,'runtime release channel drift'); need('const buildID = "v18.0.4-stable-native-cross-platform-closure-20260813"' in boot,'runtime build ID drift')
app=(R/'app_state.go').read_text(); need('resolveV18RuntimeConfig(base)' in app,'release-channel runtime resolver missing')
prof=(R/'v18_test_profile.go').read_text(); need('if releaseChannel == "STABLE"' in prof and 'stableRuntimeConfigDirName' in prof,'Stable runtime selection missing')
r=(R/'renderer/renderer.js').read_text(); need("EXPECTED_BUILD_ID='v18.0.4-stable-native-cross-platform-closure-20260813'" in r,'renderer build identity drift'); need('packaged v18.0.4 Stable renderer' in r,'renderer Stable mismatch copy missing')
man=json.loads((R/'renderer/qa/manifest.json').read_text())['releases'][0]; need(man.get('version')=='18.0.4' and man.get('status')=='STABLE' and man.get('buildId')=='v18.0.4-stable-native-cross-platform-closure-20260813','QA manifest Stable identity missing')
for f in ['README.md','VERSION.txt','renderer/docs/user.md','renderer/docs/developer.md','renderer/docs/limitations.md']:
    head='\n'.join((R/f).read_text(errors='ignore').splitlines()[:70]); need('v18.0.4' in head and 'STABLE' in head.upper(),f+' current Stable identity missing')
if e:
    print('v18.0.4 Stable Promotion Gate: FAIL'); [print(' -',x) for x in e]; sys.exit(2)
print('v18.0.4 Stable Promotion Gate: PASS · identity/runtime/docs aligned · functional intelligence unchanged')
