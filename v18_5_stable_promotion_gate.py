#!/usr/bin/env python3
from pathlib import Path
import json,sys
R=Path(__file__).resolve().parent;e=[]
def need(ok,msg):
 if not ok:e.append(msg)
i=json.loads((R/'release_identity.json').read_text()); b=(R/'app_bootstrap.go').read_text(); r=(R/'README.md').read_text(); m=json.loads((R/'renderer/qa/manifest.json').read_text()); lim=(R/'renderer/docs/limitations.md').read_text(); qa=(R/'renderer/qa/v18.5.0.txt').read_text()
need(i.get('version')=='18.5.0','version drift'); need(i.get('channel')=='STABLE','channel not STABLE'); need(i.get('build_id')=='v18.5.0-stable-major-closure-release-assurance-20260817','Stable build ID drift')
need(i.get('stable_baseline')=='v17.5.1','v18 provenance anchor drift'); need(i.get('previous_stable')=='v18.4.0' and i.get('patch_predecessor')=='v18.4.0','predecessor drift')
need(i.get('runtime_config')=='PersonalMarketTerminal','Stable runtime profile drift'); need(i.get('application_bundle')=='De-Pulse.app','Stable bundle drift')
need('const releaseChannel = "STABLE"' in b,'runtime channel drift'); need('const buildID = "v18.5.0-stable-major-closure-release-assurance-20260817"' in b,'runtime build drift')
need(r.startswith('# DE.PULSE v18.5.0 STABLE — Major Closure & Release Assurance'),'README Stable identity drift')
q=m.get('releases',[{}])[0]; need(q.get('version')=='18.5.0' and q.get('status')=='STABLE' and q.get('buildId')=='v18.5.0-stable-major-closure-release-assurance-20260817','QA manifest Stable identity drift')
need('Final v18.5 STABLE publication requires G0-G15' in lim,'candidate limitations must keep certification fail-closed')
need('final publication requires exact Stable G11-G15 certification' in qa,'QA candidate must not claim unearned G15')
for f in ['README.md','renderer/docs/user.md','renderer/docs/limitations.md']: need('No Execution' in (R/f).read_text(),f+' lost No Execution boundary')
if e: print('v18.5 Stable Promotion Gate: FAIL'); [print(' -',x) for x in e]; sys.exit(2)
print('v18.5 Stable Promotion Gate: PASS · identity-only promotion · provenance preserved · certification claims fail-closed')
