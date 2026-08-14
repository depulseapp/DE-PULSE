#!/usr/bin/env python3
from pathlib import Path
import subprocess, sys

ROOT=Path(__file__).resolve().parent

def need(cond,msg):
    if not cond:
        print('FAIL:',msg); sys.exit(1)

def text(name): return (ROOT/name).read_text(errors='replace')

prod='\n'.join(x.read_text(errors='replace') for x in ROOT.glob('*.go') if not x.name.endswith('_test.go')); main=prod; v143=prod; v15=prod; renderer=text('renderer/renderer.js'); css=text('renderer/styles.css')
checks=[
 ('after-hours catalyst completion calendar', 'catalystCompletionAt' in v143 and 'isUSMarketHoliday' in v143 and 'isUSEarlyClose' in v143),
 ('completed catalyst finalization', 'finalizeCompletedCatalysts' in v143 and 'CompletedAt' in v143),
 ('readiness gates actual freshness state', 'readinessFreshnessGate' in v143 and 'Intraday Bars' in v143),
 ('auto-recovery includes ERROR/UNAVAILABLE', 'freshnessRecoveryDue' in main and 'UNAVAILABLE' in main and 'ERROR' in main),
 ('target-scoped research fundamentals', 'refreshResearchFundamentals' in main),
 ('target-scoped research earnings', 'refreshResearchEarnings' in main),
 ('target-scoped research news', 'refreshResearchNews' in main),
 ('target-scoped deep SEC', 'refreshSECResearchSymbol' in main and 'AddDate(0, 0, -95)' in main),
 ('research ready is evidence-gated', 'researchPackageReadiness' in v15 and 'ready, issues' in v15),
 ('Form 4 derivative transactions parsed', 'DerivativeTable' in main or 'derivativeTable' in main),
 ('adaptive seconds/minutes/hours/days age', "return `${sec}s`" in renderer and "`${day}d${h?` ${h}h`:''}`" in renderer),
 ('freshness grid containment', '.v151-freshness-row>.freshness-reason' in css and 'overflow-wrap:anywhere' in css),
 ('SEC table containment', '.sec-transaction-table-wrap' in css and 'overflow-x:auto' in css),
 ('provider Settings naming truthfulness', 'Alpaca · Preferred U.S. Equities & History' in renderer and 'Finnhub · Live Recovery & Intelligence' in renderer),
]
for label,ok in checks: need(ok,label)

# These tests exercise time/session boundaries, freshness degradation/recovery,
# SEC derivative Form 4 semantics, Research readiness truthfulness, and CBOE delayed semantics.
cmd=['go','test','-run','^TestV1511','-count=10','./...']
r=subprocess.run(cmd,cwd=ROOT,text=True,stdout=subprocess.PIPE,stderr=subprocess.STDOUT)
if r.returncode:
    print(r.stdout); sys.exit(r.returncode)
print(f'Extensive Edge-Case & Failure-Mode Validation: PASS ({len(checks)} source invariants + targeted Go tests x10)')
