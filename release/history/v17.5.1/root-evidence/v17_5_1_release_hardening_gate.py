#!/usr/bin/env python3
from pathlib import Path
import json,sys,re
R=Path(__file__).resolve().parent; e=[]
def need(ok,msg):
    if not ok:e.append(msg)
i=json.loads((R/'release_identity.json').read_text())
need(i.get('version')=='17.5.1','version must be 17.5.1')
need(i.get('channel')=='STABLE','channel must be STABLE')
need(i.get('build_id')=='v17.5.1-stable-release-identity-documentation-hardening-20260813','Stable build ID drift')
need(i.get('patch_predecessor')=='v17.5.0','patch predecessor missing')
for doc in ('user.md','developer.md','limitations.md'):
    text=(R/'renderer/docs'/doc).read_text(errors='ignore'); head='\n'.join(text.splitlines()[:55])
    need('v17.5.1 STABLE' in head,doc+' current Stable identity missing')
    for v in ('v17.0.0','v17.1.0','v17.2.0','v17.3.0','v17.4.0','v17.5.0'):
        need(v in head,doc+' missing '+v+' delivery history')
    for stale in ('This is still a TEST candidate','not a Stable promotion','Final v17 Major Closure must still complete'):
        need(stale not in head,doc+' stale promotion wording: '+stale)
r=(R/'renderer/renderer.js').read_text(errors='ignore')
need("EXPECTED_RELEASE_VERSION='17.5.1'" in r,'renderer version identity drift')
need("EXPECTED_BUILD_ID='v17.5.1-stable-release-identity-documentation-hardening-20260813'" in r,'renderer build identity drift')
need('packaged v17.5.1 Stable renderer' in r,'stale mismatch message')
need('packaged v16.8.1 Stable renderer' not in r,'legacy v16.8.1 mismatch message remains')
boot=(R/'app_bootstrap.go').read_text(errors='ignore')
need('const appVersion = "17.5.1"' in boot,'runtime appVersion drift')
need('const buildID = "v17.5.1-stable-release-identity-documentation-hardening-20260813"' in boot,'runtime buildID drift')
man=json.loads((R/'renderer/qa/manifest.json').read_text()); first=man.get('releases',[{}])[0]
need(first.get('version')=='17.5.1' and first.get('status')=='STABLE','QA manifest current Stable identity missing')
reg=json.loads((R/'release_learning_registry.json').read_text())
need(any(x.get('id')=='RL-038' and x.get('status')=='ACTIVE' for x in reg.get('entries',[])),'RL-038 missing')
if e:
    print('v17.5.1 Release Identity & Documentation Hardening: FAIL'); [print(' -',x) for x in e]; sys.exit(2)
print('v17.5.1 Release Identity & Documentation Hardening: PASS · Stable identity + complete v17.0-v17.5 in-app history + stale-copy protections')
