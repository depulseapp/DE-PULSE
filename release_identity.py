#!/usr/bin/env python3
from __future__ import annotations
import argparse, json, re
from pathlib import Path
ROOT=Path(__file__).resolve().parent
IDENTITY=ROOT/'release_identity.json'
RELEASE_COUPLED_ASSETS=(
    'renderer.js',
    'live-dom-reconcile.js',
    'watchlist-v18.5.1.js',
    'watchlist-v18.5.1.css',
    'header-v18.5.1.js',
    'ui-v18.5.1.css',
    'surface-consolidation-v18.6.js',
    'surface-consolidation-v18.6.css',
    'documentation-access-v18.6.js',
)

def load():
    x=json.loads(IDENTITY.read_text())
    for k in ('version','display_version','build_id','channel','stable_baseline','previous_stable','scope','bundle_version','runtime_config','application_bundle'):
        if not str(x.get(k,'')).strip(): raise SystemExit(f'release identity missing {k}')
    return x

def sync(x):
    (ROOT/'VERSION.txt').write_text(
        f"{x['display_version']}\n"
        f"Build: {x['build_id']}\n"
        f"Channel: {x['channel']}\n"
        f"Stable baseline: {x['stable_baseline']}\n"
        f"Previous Stable: {x['previous_stable']}\n"
        f"Scope: {x['scope']}\n"
        f"Application bundle: {x['application_bundle']}\n"
        "Deterministic Day/Swing/Long formulas: unchanged / v14.3.7-compatible\n"
        f"Runtime/config continuity: {x['runtime_config']} (preserves compatible prior Stable settings and API keys)\n"
        "Release learning: adaptive G0-G16 Release Learning Registry active · Build Process v2 pre-freeze qualification + canonical identity enabled\n"
        "Adaptive contracts: US Equities Processing Boundary · Data Utility/Evidence Value · Performance/Scalability · Data Health/Freshness/Cache · Testing · Intelligence/Learning preserved\n"
    )
    p=ROOT/'app_bootstrap.go'; s=p.read_text()
    s=re.sub(r'const appVersion = "[^"]+"', f'const appVersion = "{x["version"]}"', s)
    s=re.sub(r'const buildID = "[^"]+"', f'const buildID = "{x["build_id"]}"', s)
    s=re.sub(r'const releaseChannel = "[^"]+"', f'const releaseChannel = "{x["channel"]}"', s)
    p.write_text(s)
    p=ROOT/'renderer/renderer.js'; s=p.read_text()
    s=re.sub(r"const EXPECTED_RELEASE_VERSION='[^']+';", f"const EXPECTED_RELEASE_VERSION='{x['version']}';", s)
    s=re.sub(r"const EXPECTED_BUILD_ID='[^']+';", f"const EXPECTED_BUILD_ID='{x['build_id']}';", s)
    p.write_text(s)
    p=ROOT/'renderer/index.html'; s=p.read_text()
    s=re.sub(r'<title>DE\.PULSE v[^<]+</title>', f"<title>DE.PULSE v{x['version']}</title>", s)
    for asset in RELEASE_COUPLED_ASSETS:
        s=re.sub(rf'{re.escape(asset)}\?v=[0-9.]+', f"{asset}?v={x['version']}", s)
    p.write_text(s)
    for name in ('certification_plan.json','ci_pipeline_plan.json'):
        p=ROOT/name; d=json.loads(p.read_text()); d['version']=x['version']
        if name=='ci_pipeline_plan.json':
            d.setdefault('policy',{})['baseline']=x['previous_stable']+' Stable'
            d['policy']['release_channel']=x['channel']
            d['policy']['canonical_release_identity']='release_identity.json'
            d['policy']['pre_freeze_qualification']=True
            d['policy']['unique_test_evidence']=True
        p.write_text(json.dumps(d,indent=2)+"\n")

def verify(x):
    errs=[]
    version=(ROOT/'VERSION.txt').read_text()
    boot=(ROOT/'app_bootstrap.go').read_text()
    index=(ROOT/'renderer/index.html').read_text()
    cert=json.loads((ROOT/'certification_plan.json').read_text())
    ci=json.loads((ROOT/'ci_pipeline_plan.json').read_text())
    checks=[
      (x['display_version'] in version,'VERSION display'),(x['build_id'] in version,'VERSION build'),
      (f'Previous Stable: {x["previous_stable"]}' in version,'VERSION predecessor'),
      (f'const appVersion = "{x["version"]}"' in boot,'appVersion'),(f'const buildID = "{x["build_id"]}"' in boot,'buildID'),
      (f'const releaseChannel = "{x["channel"]}"' in boot,'release channel'),
      (f"const EXPECTED_RELEASE_VERSION='{x['version']}';" in (ROOT/'renderer/renderer.js').read_text(),'renderer version'),
      (f"const EXPECTED_BUILD_ID='{x['build_id']}';" in (ROOT/'renderer/renderer.js').read_text(),'renderer build'),
      (f"<title>DE.PULSE v{x['version']}</title>" in index,'HTML title'),
      (cert.get('version')==x['version'],'certification plan version'),(ci.get('version')==x['version'],'CI plan version'),
      (ci.get('policy',{}).get('baseline')==x['previous_stable']+' Stable','CI baseline'),
    ]
    checks.extend((f"{asset}?v={x['version']}" in index,f'{asset} cache-bust version') for asset in RELEASE_COUPLED_ASSETS)
    errs.extend(label for ok,label in checks if not ok)
    if errs: raise SystemExit('Release identity: FAIL · '+', '.join(errs))
    print(f"Release identity: PASS · {x['version']} · {x['build_id']}")

def main():
    ap=argparse.ArgumentParser(); ap.add_argument('--sync',action='store_true'); ap.add_argument('--verify',action='store_true'); a=ap.parse_args(); x=load()
    if a.sync: sync(x)
    if a.verify or not a.sync: verify(x)
if __name__=='__main__': main()
