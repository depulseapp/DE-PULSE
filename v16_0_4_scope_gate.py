#!/usr/bin/env python3
from pathlib import Path
import json,sys
R=Path(__file__).resolve().parent
go='\n'.join(p.read_text(errors='ignore') for p in R.glob('*.go') if not p.name.endswith('_test.go'))
tests='\n'.join(p.read_text(errors='ignore') for p in R.glob('*_test.go'))
render=(R/'renderer/renderer.js').read_text(); trace=(R/'renderer/qa/v16.0.4-traceability.md').read_text(); m=json.loads((R/'renderer/qa/v16.0.4-master-scope.json').read_text())
req={'V16.0.4-17-01','V16.0.4-17-02','V16.0.4-17-03','V16.0.4-25-01','V16.0.4-25-02','V16.0.4-25-03','V16.0.4-ES-01','V16.0.4-30-01','V16.0.4-30-02','V16.0.4-30-03','V16.0.4-UX-01'}; e=[]
if m.get('count')!=11 or {x['id'] for x in m.get('requirements',[])}!=req:e.append('scope identity/count/IDs mismatch')
for tok in ['SymbolLifecycleState','Catalyst & Material Event Context','Required Market Context','buildEvidenceSnapshot','selectedTickerEarningsEvidence','catalyst-watch']:
    if tok not in go:e.append('implementation token missing: '+tok)
for tok in ['TestV1604LifecycleCanonicalStatesAndProvenance','TestV1604NoActiveCatalystCannotBeFreshWithStaleNewsCheck','TestV1604MarketContextWorstRequiredDependencyAndWeekendSemantics','TestV1604EvidenceSnapshotChangesOnlyForMaterialEvidence']:
    if tok not in tests:e.append('verification token missing: '+tok)
for tok in ['lifecycle','Corporate Actions & Symbol Lifecycle','Evidence Snapshot']:
    if tok not in render:e.append('surface token missing: '+tok)
for rid in req:
    if rid not in trace:e.append('traceability missing '+rid)
if e: print('v16.0.4 Immutable Master Scope Gate: FAIL'); print('\n'.join('- '+x for x in e)); sys.exit(1)
print('v16.0.4 Immutable Master Scope Gate: PASS · 11/11 requirements represented')
