#!/usr/bin/env python3
from pathlib import Path
import json,sys
R=Path(__file__).resolve().parent;e=[]
def need(ok,msg):
    if not ok:e.append(msg)
i=json.loads((R/'release_identity.json').read_text());boot=(R/'app_bootstrap.go').read_text();r=(R/'renderer/renderer.js').read_text();readme=(R/'README.md').read_text();user=(R/'renderer/docs/user.md').read_text();dev=(R/'renderer/docs/developer.md').read_text();lim=(R/'renderer/docs/limitations.md').read_text();man=json.loads((R/'renderer/qa/manifest.json').read_text());qa=(R/'renderer/qa/v18.4.0.txt').read_text();cert=json.loads((R/'certification_plan.json').read_text())
build='v18.4.0-stable-security-commercial-readiness-20260816'
need(i.get('version')=='18.4.0','version drift');need(i.get('channel')=='STABLE','channel is not STABLE');need(i.get('build_id')==build,'Stable build identity drift');need(i.get('runtime_config')=='PersonalMarketTerminal','Stable desktop runtime must use canonical profile');need(i.get('application_bundle')=='De-Pulse.app','Stable application bundle drift');need(i.get('previous_stable')=='v18.3.0' and i.get('patch_predecessor')=='v18.3.0','v18.3 predecessor drift')
need('const releaseChannel = "STABLE"' in boot,'runtime release channel not STABLE');need(f'const buildID = "{build}"' in boot,'runtime Stable build ID drift');need(f"EXPECTED_BUILD_ID='{build}'" in r,'renderer Stable build identity drift')
need(readme.startswith('# DE.PULSE v18.4.0 STABLE — Security / Commercial Readiness Hardening'),'README Stable title drift');need('**Runtime/config:** `PersonalMarketTerminal`' in '\n'.join(readme.splitlines()[:12]),'README Stable runtime drift')
need('## v18.4.0 STABLE — Security / commercial readiness hardening' in '\n'.join(user.splitlines()[:15]),'user docs Stable section drift');need('## v18.4.0 STABLE — Security / commercial-readiness architecture' in '\n'.join(dev.splitlines()[:15]),'developer docs Stable section drift');need('## v18.4.0 STABLE — commercial/data-rights boundary' in '\n'.join(lim.splitlines()[:15]),'limitations Stable section drift')
first=man.get('releases',[{}])[0];need(first.get('version')=='18.4.0' and first.get('status')=='STABLE' and first.get('buildId')==build,'QA manifest Stable identity drift');need(build in qa and 'No Execution Boundary preserved' in qa,'v18.4 Stable QA record drift')
labels={x.get('id'):x.get('label') for x in cert.get('checks',[])};need(labels.get('g0_release_identity')=='Canonical v18.4.0 STABLE release identity','certification-plan Stable identity label drift')
for f in ['README.md','renderer/docs/user.md','renderer/docs/developer.md','renderer/docs/limitations.md']:need('No Execution' in (R/f).read_text(),f+' lost No Execution boundary')
prof=(R/'v18_test_profile.go').read_text();need('PersonalMarketTerminal-v18.4.0-TEST' in prof,'historical TEST isolation contract lost');need('if releaseChannel == "STABLE"' in prof,'canonical Stable runtime routing lost')
if e:print('v18.4 Stable Promotion Gate: FAIL');[print(' -',x) for x in e];sys.exit(2)
print('v18.4 Stable Promotion Gate: PASS · identity-only promotion · canonical desktop Stable runtime · security/commercial boundaries preserved')
