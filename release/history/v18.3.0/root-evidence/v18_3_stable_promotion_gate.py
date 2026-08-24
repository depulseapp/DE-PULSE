#!/usr/bin/env python3
from pathlib import Path
import json,sys
R=Path(__file__).resolve().parent;e=[]
def need(ok,msg):
 if not ok:e.append(msg)
i=json.loads((R/'release_identity.json').read_text());boot=(R/'app_bootstrap.go').read_text();readme=(R/'README.md').read_text();user=(R/'renderer/docs/user.md').read_text();dev=(R/'renderer/docs/developer.md').read_text();lim=(R/'renderer/docs/limitations.md').read_text();man=json.loads((R/'renderer/qa/manifest.json').read_text())
need(i.get('version')=='18.3.0','version drift');need(i.get('channel')=='STABLE','channel is not STABLE');need(i.get('build_id')=='v18.3.0-stable-postgresql-hosted-shared-state-20260815','Stable build identity drift');need(i.get('runtime_config')=='PersonalMarketTerminal','Stable desktop runtime must use canonical profile');need(i.get('application_bundle')=='De-Pulse.app','Stable application bundle drift');need(i.get('previous_stable')=='v18.2.0' and i.get('patch_predecessor')=='v18.2.0','v18.2 predecessor drift')
need('const releaseChannel = \"STABLE\"' in boot,'runtime release channel not STABLE');need('const buildID = \"v18.3.0-stable-postgresql-hosted-shared-state-20260815\"' in boot,'runtime Stable build ID drift');need(readme.startswith('# DE.PULSE v18.3.0 STABLE — PostgreSQL / Hosted Shared State'),'README Stable title drift');need('**Runtime/config:** `PersonalMarketTerminal`' in '\n'.join(readme.splitlines()[:12]),'README Stable runtime drift')
need('## v18.3.0 STABLE — PostgreSQL / hosted shared state' in '\n'.join(user.splitlines()[:15]),'user docs Stable section drift');need('## v18.3.0 STABLE — PostgreSQL / hosted shared-state architecture' in '\n'.join(dev.splitlines()[:15]),'developer docs Stable section drift');need('## v18.3.0 STABLE — hosted persistence boundaries' in '\n'.join(lim.splitlines()[:15]),'limitations Stable section drift')
first=man.get('releases',[{}])[0];need(first.get('version')=='18.3.0' and first.get('status')=='STABLE' and first.get('buildId')==i.get('build_id'),'QA manifest Stable identity drift')
for f in ['README.md','renderer/docs/user.md','renderer/docs/developer.md','renderer/docs/limitations.md']:need('No Execution' in (R/f).read_text(),f+' lost No Execution boundary')
need('PersonalMarketTerminal-v18.3.0-TEST' in (R/'v18_test_profile.go').read_text(),'historical TEST isolation contract lost');need('if releaseChannel == \"STABLE\"' in (R/'v18_test_profile.go').read_text(),'canonical Stable runtime routing lost')
if e:print('v18.3 Stable Promotion Gate: FAIL');[print(' -',x) for x in e];sys.exit(2)
print('v18.3 Stable Promotion Gate: PASS · identity-only promotion · canonical desktop Stable runtime · hosted PostgreSQL boundary preserved')
