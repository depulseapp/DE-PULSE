#!/usr/bin/env python3
from pathlib import Path
import json,sys,subprocess
R=Path(__file__).resolve().parent; go='\n'.join(p.read_text(errors='ignore') for p in R.glob('*.go') if not p.name.endswith('_test.go')); tests=(R/'v16_1_1_test.go').read_text(); m=json.loads((R/'renderer/qa/v16.1.1-master-scope.json').read_text()); e=[]
req={'V16.1.1-MI-HORIZON-01','V16.1.1-MI-FRESH-01','V16.1.1-MI-LIQ-01'}
if m.get('count')!=3 or {x['id'] for x in m.get('requirements',[])}!=req:e.append('patch scope identity/count mismatch')
for tok in ['marketIntelligenceBarEvidence','case "long":','key, lookback = "weekly", 26','hasValidSpread']:
    if tok not in go:e.append('implementation invariant missing '+tok)
for tok in ['TestV1611LongStructureRequiresWeeklyEvidence','TestV1611StructureRejectsStaleHistoricalBars','TestV1611RelativeStrengthRejectsStaleBarsAndUsesEvidenceTimestamp','TestV1611LiquidityWithoutValidBidAskIsUnknown']:
    if tok not in tests:e.append('patch regression missing '+tok)
q=subprocess.run(['go','test','-count=1','-run','TestV1611','./...'],cwd=R,capture_output=True,text=True)
if q.returncode:e.append('v16.1.1 regressions failed: '+(q.stdout+q.stderr)[-700:])
if e: print('v16.1.1 Truth-Hardening Scope Gate: FAIL'); print('\n'.join('- '+x for x in e)); sys.exit(1)
print('v16.1.1 Truth-Hardening Scope Gate: PASS · 3/3')
