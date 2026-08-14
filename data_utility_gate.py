#!/usr/bin/env python3
import json, re
from pathlib import Path
R=Path(__file__).resolve().parent
reg=json.loads((R/'data_utility_registry.json').read_text())
errs=[]; seen=set()
for d in reg.get('datasets',[]):
    name=str(d.get('dataset','')).strip()
    if not name or name in seen: errs.append('duplicate/missing dataset '+name)
    seen.add(name)
    if not str(d.get('owner','')).strip(): errs.append(name+': owner missing')
    if not str(d.get('purpose','')).strip(): errs.append(name+': purpose missing')
    if not d.get('consumers') and not str(d.get('retention_reason','')).strip(): errs.append(name+': no consumer/retention reason')
    if str(d.get('utility','')).upper() not in {'ACTIVE','STRATEGIC'}: errs.append(name+': invalid utility class')
source=(R/'data_freshness.go').read_text()
for expected in ['Quotes','VIX','Intraday Bars','Daily / Weekly History','News','Earnings','SEC Filings','Fundamentals','Global','Macro','Options']:
    if expected not in seen: errs.append('freshness dataset not registered: '+expected)
    if f'{{"{expected}"' not in source: errs.append('freshness registry/source drift: '+expected)
if errs:
    print('Data Utility: FAIL')
    for e in errs: print(' -',e)
    raise SystemExit(2)
print(f'Data Utility: PASS · {len(seen)} datasets have owners, purposes and consumers/retention justification')
