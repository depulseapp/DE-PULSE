#!/usr/bin/env python3
from pathlib import Path
import json,sys,subprocess
R=Path(__file__).resolve().parent; go='\n'.join(p.read_text(errors='ignore') for p in R.glob('*.go') if not p.name.endswith('_test.go')); tests='\n'.join(p.read_text(errors='ignore') for p in R.glob('*_test.go')); idx=(R/'renderer/index.html').read_text(); m=json.loads((R/'renderer/qa/v16.0.5-master-scope.json').read_text()); e=[]
req={'V16.0.5-01','V16.0.5-02','V16.0.5-03','V16.0.5-04'}
if m.get('count')!=4 or {x['id'] for x in m.get('requirements',[])}!=req:e.append('scope identity/count/IDs mismatch')
if 'valid contemporaneous provider observations' not in go:e.append('Provider Reconciliation contemporaneity invariant missing')
for tok in ['TestV1600ProviderReconciliationDoesNotUseNonContemporaneousObservationAsConflict','TestV1600ProviderReconciliationAllowsContemporaneousPriorCloseEvidence']:
    if tok not in tests:e.append('Provider Reconciliation closure regression missing: '+tok)
if any(x in idx for x in ['WC ·','TEST ·','test-build-badge']):e.append('release status leaked into primary header')
if '<small>Personal Market Intelligence</small>' not in idx:e.append('Stable header subtitle missing')
if 'filepath.Join(base, "PersonalMarketTerminal")' not in go or 'PersonalMarketTerminal-v16-TEST' in go:e.append('Stable runtime/config continuity missing')
for gate in ['v16_0_4_scope_gate.py','approved_scope_gate.py']:
    q=subprocess.run([sys.executable,str(R/gate)],cwd=R,capture_output=True,text=True)
    if q.returncode:e.append(gate+' failed: '+(q.stdout+q.stderr)[-500:])
if e: print('v16.0.5 Inherited Closure Scope Gate: FAIL'); print('\n'.join('- '+x for x in e)); sys.exit(1)
print('v16.0.5 Inherited Closure Scope Gate: PASS · 4/4 predecessor requirements preserved')
