#!/usr/bin/env python3
from pathlib import Path
import json,sys
R=Path(__file__).resolve().parent; e=[]
a=(R/'v16_10_information_architecture_data_efficiency_audit.md').read_text(errors='ignore')
for t in ['Discovery/Scanner remains the only scanner owner','REUSE MORE','Rapid Price Dislocation Watch','Short Trade Plan','SPY / QQQ drawdown from ATH','Provider Router','no separate tab/service']:
    if t not in a: e.append('IA audit missing '+t)
try:
    du=json.loads((R/'data_utility_registry.json').read_text()); names={x.get('dataset') for x in du.get('datasets',[])}
    if not {'Opportunity Radar','Market Activity'}.issubset(names): e.append('utility mapping incomplete')
    scope=json.loads((R/'renderer/qa/v16.10.0-master-scope.json').read_text())
    if len(scope.get('scope_lock',[]))!=10: e.append('scope lock incomplete')
except Exception as ex: e.append('structured evidence unreadable: '+str(ex))
if e:
    print('v16.10 IA/Data Efficiency Gate: FAIL'); [print(' -',x) for x in e]; sys.exit(1)
print('v16.10 IA/Data Efficiency Gate: PASS · canonical placement/data reuse classified before feature promotion')
