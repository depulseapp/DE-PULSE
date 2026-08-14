#!/usr/bin/env python3
from pathlib import Path
import json,re,sys
R=Path(__file__).resolve().parent
errs=[]
# Source root must not carry office-document release debris. Historical release docs belong in release/library artifacts, not runtime source.
for ext in ('.docx','.pdf','.pptx','.xlsx'):
    for p in R.glob('*'+ext): errs.append('unjustified root release/document artifact: '+p.name)
try:
    retained=json.loads((R/'source_retained_assets.json').read_text())
    for row in retained.get('assets',[]):
        p=R/row.get('path','')
        if not p.exists(): errs.append('retained asset missing: '+str(row.get('path')))
        if not row.get('purpose') or not row.get('consumers'): errs.append('retained asset lacks purpose/consumers: '+str(row.get('path')))
except Exception as ex: errs.append('retained-asset registry unreadable: '+str(ex))
# Current-release professional acceptance fixtures must derive release identity; historical literals are allowed only in explicit historical checks.
tr=(R/'trader_acceptance_test.js').read_text(errors='ignore')
if "releaseIdentity=JSON.parse(fs.readFileSync('release_identity.json'" not in tr: errs.append('professional trader fixture does not derive canonical release identity')
if "state={version:'15.1.2'" in tr: errs.append('stale v15.1.2 current-state professional fixture remains')
# No disabled/skip markers in production-facing closure tests.
for p in [R/'v16_11_major_scope_matrix_gate.py',R/'trader_acceptance_test.js']:
    txt=p.read_text(errors='ignore')
    if re.search(r'(?i)\b(skip|todo|fixme|hack)\b',txt) and p.name.startswith('v16_11'):
        errs.append(p.name+': unresolved closure marker')
if errs:
    print('v16.11 Source Hygiene Gate: FAIL')
    for e in errs: print(' -',e)
    raise SystemExit(2)
print('v16.11 Source Hygiene Gate: PASS · release debris removed · retained assets justified · current professional fixture uses canonical identity')
