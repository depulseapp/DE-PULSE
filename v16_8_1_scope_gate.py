#!/usr/bin/env python3
import json, subprocess, sys
from pathlib import Path
R=Path(__file__).resolve().parent
errs=[]
def run(cmd):
    p=subprocess.run(cmd,cwd=R,text=True,capture_output=True)
    if p.returncode: errs.append('command failed: '+' '.join(cmd)+'\n'+(p.stdout+p.stderr)[-3000:])
run(['go','test','-count=1','-run','^TestV1681','./...'])
run(['node','v16_8_1_renderer_test.js'])
run([sys.executable,'data_utility_gate.py'])
run([sys.executable,'data_health_policy_gate.py'])
run([sys.executable,'release_identity.py','--verify'])
scope=json.loads((R/'renderer/qa/v16.8.1-master-scope.json').read_text())
if len(scope.get('scope_lock',[])) != 6: errs.append('immutable v16.8.1 scope must contain six approved hardening/process items')
trace=(R/'renderer/qa/v16.8.1-traceability.md').read_text()
if '27 FULL / 3 PARTIAL / 0 MISSING' not in trace: errs.append('roadmap status truth missing')
if errs:
    print('v16.8.1 scope gate: FAIL'); [print(' -',e) for e in errs]; raise SystemExit(2)
print('v16.8.1 scope gate: PASS · 6/6 hardening/process clauses · roadmap unchanged 27 FULL / 3 PARTIAL / 0 MISSING')
