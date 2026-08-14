#!/usr/bin/env python3
from pathlib import Path
import json,sys
R=Path(__file__).resolve().parent; e=[]
def need(ok,msg):
    if not ok:e.append(msg)
i=json.loads((R/'release_identity.json').read_text())
B='v18.0.6-stable-smart-provider-router-rapid-move-market-shock-hardening-20260814'
need(i.get('version')=='18.0.6','version drift')
need(i.get('channel')=='STABLE','channel must be STABLE')
need(i.get('build_id')==B,'Stable build ID drift')
need(i.get('runtime_config')=='PersonalMarketTerminal','Stable runtime target drift')
need(i.get('application_bundle')=='De-Pulse.app','Stable bundle drift')
boot=(R/'app_bootstrap.go').read_text(); need('const releaseChannel = "STABLE"' in boot,'runtime release channel drift'); need(f'const buildID = "{B}"' in boot,'runtime build ID drift')
prof=(R/'v18_test_profile.go').read_text(); need('if releaseChannel == "STABLE"' in prof and 'stableRuntimeConfigDirName' in prof,'Stable runtime selection missing')
rr=(R/'renderer/renderer.js').read_text(); need(f"const EXPECTED_BUILD_ID='{B}';" in rr,'renderer build identity drift')
man=json.loads((R/'renderer/qa/manifest.json').read_text())['releases'][0]
need(man.get('version')=='18.0.6' and man.get('status')=='STABLE' and man.get('buildId')==B,'QA manifest Stable identity missing')
rapid=(R/'rapid_move_intelligence.go').read_text(); need('rapid-move-v1.1.0' in rapid and 'rapid-move-learning-v1.0.0' in rapid,'Rapid Move hardening missing')
need(all(x in rapid for x in ['MARKET_SHOCK','rapidMoveApplyHysteresis','SHADOW','VALIDATED','APPROVED','PRODUCTION']),'Rapid Move governance/hysteresis drift')
for f in ['README.md','VERSION.txt','renderer/docs/user.md','renderer/docs/developer.md','renderer/docs/limitations.md']:
    head='\n'.join((R/f).read_text(errors='ignore').splitlines()[:90]); need('v18.0.6' in head and 'STABLE' in head.upper(),f+' current Stable identity missing')
if e:
    print('v18.0.6 Stable Promotion Gate: FAIL'); [print(' -',x) for x in e]; sys.exit(2)
print('v18.0.6 Stable Promotion Gate: PASS · identity/runtime/docs aligned · Router/Rapid Move hardening preserved')
