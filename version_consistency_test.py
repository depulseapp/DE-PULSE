#!/usr/bin/env python3
import json,re
from pathlib import Path
ROOT=Path(__file__).resolve().parent
I=json.loads((ROOT/'release_identity.json').read_text())
VERSION=I['version']; BUILD=I['build_id']; PREV=I['previous_stable']; STABLE=I['stable_baseline']; CHANNEL=I['channel']; DISPLAY=I['display_version']
errs=[]
def need(ok,msg):
    if not ok: errs.append(msg)
boot=(ROOT/'app_bootstrap.go').read_text(); renderer=(ROOT/'renderer/renderer.js').read_text(); readme=(ROOT/'README.md').read_text(); version=(ROOT/'VERSION.txt').read_text()
patch_path=ROOT/'renderer'/f'watchlist-desk-contract-v{VERSION}.js'
patch=patch_path.read_text() if patch_path.exists() else ''
index=(ROOT/'renderer/index.html').read_text()
need(f'const appVersion = "{VERSION}"' in boot,'appVersion mismatch')
need(f'const buildID = "{BUILD}"' in boot,'buildID mismatch')
renderer_identity=(f"const EXPECTED_RELEASE_VERSION='{VERSION}';" in renderer and f"const EXPECTED_BUILD_ID='{BUILD}';" in renderer)
patch_identity=(f"DEPULSE_PATCH_VERSION = '{VERSION}'" in patch and f"DEPULSE_PATCH_BUILD_ID = '{BUILD}'" in patch and f'watchlist-desk-contract-v{VERSION}.js?v={VERSION}' in index)
need(renderer_identity or patch_identity,'renderer/patch version mismatch')
need(readme.startswith(f'# DE.PULSE v{VERSION}'),'README title mismatch')
need(BUILD in readme and f'Current Stable baseline:** {PREV}' in readme,'README immediate Stable predecessor mismatch')
need(f'**Channel:** {CHANNEL}' in '\n'.join(readme.splitlines()[:8]),'README current channel mismatch')
need(DISPLAY in version and BUILD in version and f'Stable baseline: {STABLE}' in version and f'Previous Stable: {PREV}' in version,'VERSION.txt mismatch')
for doc in ('user.md','developer.md','limitations.md'):
    text=(ROOT/'renderer/docs'/doc).read_text()
    expected={'user.md':'# DE.PULSE — User documentation','developer.md':'# DE.PULSE — Developer documentation','limitations.md':'# DE.PULSE — Capabilities & Limitations'}[doc]
    need(text.splitlines()[0]==expected,f'{doc} canonical first heading mismatch')
    need(f'v{VERSION} {CHANNEL}' in '\n'.join(text.splitlines()[:30]),f'{doc} current build section missing')
ci=json.loads((ROOT/'ci_pipeline_plan.json').read_text()); cert=json.loads((ROOT/'certification_plan.json').read_text())
need(ci.get('version')==VERSION,'CI plan current-release identity mismatch')
need(ci.get('policy',{}).get('baseline')==PREV+' Stable','CI plan predecessor mismatch')
need(cert.get('version')==VERSION,'certification plan current-release identity mismatch')
if errs:
    print('Version consistency: FAIL')
    for e in errs: print(' -',e)
    raise SystemExit(2)
print(f'Version consistency: PASS · v{VERSION} canonical release identity aligned across runtime, renderer patch, docs and CI metadata')
