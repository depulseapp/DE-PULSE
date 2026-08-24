#!/usr/bin/env python3
from pathlib import Path
import json,sys
R=Path(__file__).resolve().parent;e=[]
def need(ok,msg):
 if not ok:e.append(msg)
i=json.loads((R/'release_identity.json').read_text());boot=(R/'app_bootstrap.go').read_text();readme=(R/'README.md').read_text();docs=(R/'renderer/docs/user.md').read_text();man=json.loads((R/'renderer/qa/manifest.json').read_text())
need(i.get('version')=='18.2.0','version drift');need(i.get('channel')=='STABLE','channel is not STABLE');need(i.get('build_id')=='v18.2.0-stable-admin-presence-sessions-20260814','Stable build identity drift');need(i.get('runtime_config')=='PersonalMarketTerminal','Stable runtime must use canonical profile');need(i.get('application_bundle')=='De-Pulse.app','Stable application bundle drift');need(i.get('previous_stable')=='v18.1.0' and i.get('patch_predecessor')=='v18.1.0','v18.1 predecessor drift')
need('const releaseChannel = \"STABLE\"' in boot,'runtime release channel not STABLE');need('const buildID = \"v18.2.0-stable-admin-presence-sessions-20260814\"' in boot,'runtime Stable build ID drift');need(readme.startswith('# DE.PULSE v18.2.0 STABLE — Admin / Presence / Sessions'),'README Stable title drift');need('**Runtime/config:** `PersonalMarketTerminal`' in '\n'.join(readme.splitlines()[:12]),'README Stable runtime drift');need('## v18.2.0 STABLE — Administration, Presence & Sessions' in '\n'.join(docs.splitlines()[:20]),'user docs Stable section drift');first=man.get('releases',[{}])[0];need(first.get('version')=='18.2.0' and first.get('status')=='STABLE' and first.get('buildId')==i.get('build_id'),'QA manifest Stable identity drift')
for f in ['README.md','renderer/docs/user.md','renderer/docs/limitations.md']:need('No Execution' in (R/f).read_text(),f+' lost No Execution boundary')
if e:print('v18.2 Stable Promotion Gate: FAIL');[print(' -',x) for x in e];sys.exit(2)
print('v18.2 Stable Promotion Gate: PASS · identity-only promotion boundary + canonical Stable runtime')
