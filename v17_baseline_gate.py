#!/usr/bin/env python3
import json, sys
from pathlib import Path
R=Path(__file__).resolve().parent
x=json.loads((R/'v17_baseline_contract.json').read_text())
i=json.loads((R/'release_identity.json').read_text())
errors=[]
expected={
 'stable':'v16.11.0',
 'source_zip_sha256':'d7da932d380dbeece82753b45fa9165e58d22cde55bab040625ac68f81f476f7',
 'clean_source_fingerprint':'f4f7fc930615a9276ec885147a3198ee622c613771d456165df2524410f98275',
 'source_zip_entries':384,
 'v16_to_v17':'GO',
 'release_learning_through':'RL-031'
}
for k,v in expected.items():
    if x.get(k)!=v: errors.append(f'{k} drift: {x.get(k)!r}')
if i.get('stable_baseline')!='v16.11.0' or i.get('previous_stable')!='v16.11.0':
    errors.append('current v17 identity must remain based on v16.11.0 Stable')
if not str(i.get('version','')).startswith('17.') or i.get('channel') not in {'TEST','RC','STABLE'}:
    errors.append('current v17 identity must preserve the certified v16.11 baseline across TEST/RC/STABLE')
if errors:
    print('v17 baseline gate: FAIL')
    for e in errors: print(' -',e)
    sys.exit(1)
print('v17 baseline gate: PASS · certified v16.11.0 provenance anchored · v16→v17 GO retained')
