#!/usr/bin/env python3
from pathlib import Path
import subprocess,json,sys
R=Path(__file__).resolve().parent; errs=[]
def run(cmd,label):
 p=subprocess.run(cmd,cwd=R,text=True,capture_output=True)
 if p.returncode: errs.append(label+' failed:\n'+(p.stdout+p.stderr)[-2500:])
run([sys.executable,'data_utility_gate.py'],'Data Utility')
run([sys.executable,'v16_11_source_hygiene_gate.py'],'Source Hygiene')
try:
 pol=json.loads((R/'source_health_baseline.json').read_text())
 if pol.get('schema')!='DE.PULSE-SOURCE-HEALTH-POLICY-2' or not pol.get('major_closure_requires_source_review'): errs.append('active source-health policy incomplete')
 reg=json.loads((R/'data_utility_registry.json').read_text())
 if len(reg.get('datasets',[]))<14: errs.append('dataset utility registry unexpectedly incomplete')
except Exception as ex: errs.append('closure utility/source policy unreadable: '+str(ex))
rep=(R/'renderer/qa/v16.11.0-data-utility-source-hygiene-audit.md')
if not rep.exists(): errs.append('closure utility/source audit report missing')
else:
 t=rep.read_text(errors='ignore')
 for token in ['KEEP / ACTIVE','REUSE MORE','STOP / REMOVE','obsolete v16.1 TEST handoff DOCX','tracked decomposition watchlist','PASS']:
  if token not in t: errs.append('closure utility/source report missing '+token)
if errs:
 print('v16.11 Data Utility / Source Hygiene Closure: FAIL'); [print(' -',e) for e in errs]; raise SystemExit(2)
print('v16.11 Data Utility / Source Hygiene Closure: PASS · active data has consumers/retention justification · source debris removed · maintainability policy active')
