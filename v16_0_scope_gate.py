#!/usr/bin/env python3
"""v16.0.3 committed-scope proof. Heavy professional/edge authorization is a separate gate."""
from pathlib import Path
import json,subprocess,sys
ROOT=Path(__file__).resolve().parent
m=json.loads((ROOT/'renderer/qa/v16.0.3-approved-scope.json').read_text())
errors=[]
if m.get('version')!='16.0.3 TEST' or m.get('count')!=4 or len(m.get('items',[]))!=4:
 errors.append('v16.0.3 scope manifest identity/count mismatch')
for item in m.get('items',[]):
 for e in item.get('evidence',[]):
  p=ROOT/e['file']
  if not p.exists() or e['contains'] not in p.read_text(errors='ignore'):
   errors.append(f"missing evidence #{item.get('id')}: {e['file']}::{e['contains']}")
commands=[
 ['go','test','-count=1','-run','TestV1603(Professional|Edge|Authorization|FaultInjection)','./...'],
 ['go','test','-count=1','./...'],
]
for cmd in commands:
 p=subprocess.run(cmd,cwd=ROOT,text=True,capture_output=True)
 if p.returncode: errors.append('failed: '+' '.join(cmd)+'\n'+p.stdout[-6000:]+'\n'+p.stderr[-6000:])
if errors:
 print('v16.0.3 Professional Authorization Hardening Scope Gate: FAIL');print('\n'.join(errors));sys.exit(1)
print('v16.0.3 Scope Gate: 4/4 CLOSED · feature + professional + edge + independent authorization fixtures PASS')
