#!/usr/bin/env python3
from pathlib import Path
import json,sys
R=Path(__file__).resolve().parent; go='\n'.join(p.read_text(errors='ignore') for p in R.glob('*.go') if not p.name.endswith('_test.go')); main=(R/'main.go').read_text(); html=(R/'renderer/index.html').read_text(); readme=(R/'README.md').read_text(); m=json.loads((R/'renderer/qa/v16.0.6-master-scope.json').read_text()); e=[]
req={f'V16.0.6-{i:02d}' for i in range(1,6)}
if m.get('count')!=5 or {x.get('id') for x in m.get('requirements',[])}!=req:e.append('scope identity/count/IDs mismatch')
if 'filepath.Join(base, "PersonalMarketTerminal")' not in go or 'PersonalMarketTerminal-v16-TEST' in go:e.append('Stable runtime/config continuity missing')
if any(x in html for x in ['WC ·','TEST ·','test-build-badge']):e.append('primary header release-status badge regression')
if 'NewApplication()' not in main or 'openAppWindow' not in main or 'certification_runner' in main:e.append('normal-launch application path missing or contaminated')
if 'deterministic Day/Swing/Long' not in readme or not any(x in readme for x in ['cannot silently mutate','never silently rewrite']):e.append('deterministic contract missing')
if e: print('v16.0.6 Inherited Consolidation Scope Gate: FAIL'); print('\n'.join('- '+x for x in e)); sys.exit(1)
print('v16.0.6 Inherited Consolidation Scope Gate: PASS · 5/5 requirements preserved')
