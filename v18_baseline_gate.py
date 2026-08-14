#!/usr/bin/env python3
from pathlib import Path
import json,sys
R=Path(__file__).resolve().parent
x=json.loads((R/'v18_baseline_contract.json').read_text()); i=json.loads((R/'release_identity.json').read_text()); e=[]
def need(ok,msg):
    if not ok:e.append(msg)
expected={
 'stable':'v17.5.1',
 'stable_build_id':'v17.5.1-stable-release-identity-documentation-hardening-20260813',
 'source_zip_sha256':'0cdf81b9438914393a381ec8f64d8cad26b36db8759cb3c21bceb90a34ac5b6c',
 'clean_source_fingerprint':'4ac43087344bd46ac92081be5b8f9a919cc4eb59a7d4cbb138995b40b58371e8',
 'macos_zip_sha256':'f11307d9e8d048ae46f5acebd20e74f265b6ba40e00e5413ebaa73b520bc23f5',
 'windows_zip_sha256':'263e0abbbca3a03d16869a11551fd37b41de32bbb9c05cdeab1899904429ba02',
 'final_manifest_sha256':'86f860b723ac98f1b99ef76884e39f5ad427efed0188516b6483d7250955a8f5',
 'v17_to_v18':'GO', 'no_stable_overwrite':True}
for k,v in expected.items(): need(x.get(k)==v,f'{k} drift: {x.get(k)!r}')
need(str(i.get('version','')).startswith('18.') and i.get('channel') in {'TEST','RC','STABLE'},'current v18 identity missing')
need(i.get('stable_baseline')=='v17.5.1','v18 major lineage must remain anchored to v17.5.1 Stable')
previous=str(i.get('previous_stable',''))
need(previous=='v17.5.1' or (previous.startswith('v18.') and previous!=f"v{i.get('version','')}"),'previous Stable must be the certified predecessor while v17.5.1 remains the major anchor')

if i.get('channel')=='STABLE':
    need(i.get('runtime_config')=='PersonalMarketTerminal','v18 Stable must use canonical Stable runtime profile')
else:
    need(i.get('runtime_config')!='PersonalMarketTerminal','v18 TEST/RC may not write into the Stable runtime profile')
if e:
 print('v18 baseline gate: FAIL'); [print(' -',x) for x in e]; sys.exit(1)
print('v18 baseline gate: PASS · certified v17.5.1 provenance anchored · v17→v18 GO retained · release-channel runtime target protected')
