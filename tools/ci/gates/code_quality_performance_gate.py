#!/usr/bin/env python3
"""DE.PULSE G7 — Code Quality + Maintainability + Performance Health.

G2 architecture/source health is intentionally a separate earlier gate. G7 verifies
that the implemented candidate is formatted, statically sound, bounded in common
resource patterns, documented where invariants are non-obvious, and still passes
focused performance/concurrency regressions.
"""
from __future__ import annotations
from pathlib import Path
import argparse
import re
import subprocess
import sys

ROOT=Path(__file__).resolve().parents[3]
ap=argparse.ArgumentParser()
ap.add_argument('--static-only',action='store_true',help='run formatting/static/comment/resource hygiene only; performance checks are checkpointed separately by certification_runner.py')
args=ap.parse_args()
PROD=sorted(p for p in ROOT.glob('*.go') if not p.name.endswith('_test.go'))
errors=[]

def run(cmd, timeout, label):
    try:
        p=subprocess.run(cmd,cwd=ROOT,text=True,capture_output=True,timeout=timeout)
        if p.returncode:
            errors.append(f"{label} failed:\n{p.stdout[-5000:]}\n{p.stderr[-5000:]}")
        return p.stdout
    except subprocess.TimeoutExpired as e:
        errors.append(f"{label} timed out:\n{str(e.stdout or '')[-2500:]}\n{str(e.stderr or '')[-2500:]}")
        return ''

# Formatting/static correctness.
out=run(['gofmt','-l',*[p.name for p in ROOT.glob('*.go')]],60,'gofmt')
if out.strip():
    errors.append('gofmt required for: '+', '.join(out.split()))
run(['go','vet','./...'],90,'go vet')

# Resource/lifecycle hygiene that is cheap enough to block every build.
for path in PROD:
    lines=path.read_text(errors='ignore').splitlines()
    for i,line in enumerate(lines,1):
        if 'time.Tick(' in line:
            errors.append(f'{path.name}:{i} uses time.Tick; use NewTicker + Stop so lifecycle ownership is explicit')
        if 'http.Client{' in line:
            window=' '.join(lines[i-1:min(len(lines),i+5)])
            if 'Timeout:' not in window:
                errors.append(f'{path.name}:{i} constructs http.Client without an explicit Timeout')

# Developer comments are required around canonical/trading/concurrency invariants.
critical=['catalystQuoteReactionActive','executeProviderRoute','buildProviderReconciliation','selectedTickerEarningsEvidence','catalystMaterialEventComponent','requiredMarketContextComponent','buildResearchPackageTruth','buildCorporateActionTruth','buildEvidenceSnapshot','refreshAlpacaRawHistoryForCorporateActions','marketIntelligenceBarEvidence','marketTradeability']
all_lines=[]
for path in PROD:
    for i,line in enumerate(path.read_text(errors='ignore').splitlines()): all_lines.append((path.name,i,line))
for fn in critical:
    found=None
    rx=re.compile(r'^func\s+(?:\([^)]*\)\s*)?'+re.escape(fn)+r'\s*\(')
    for path in PROD:
        lines=path.read_text(errors='ignore').splitlines()
        for idx,line in enumerate(lines):
            if rx.search(line): found=(path,idx,lines); break
        if found: break
    if not found:
        errors.append(f'critical invariant function missing: {fn}'); continue
    path,idx,lines=found; prior='\n'.join(lines[max(0,idx-6):idx]).strip()
    if '//' not in prior: errors.append(f'{path.name}:{fn} needs a WHY/invariant developer comment immediately above it')

# Standalone mode preserves the complete G7 behavior, but runs each expensive
# performance/concurrency case independently. The permanent certification runner
# checkpoints these as separate checks so an interruption never loses prior PASSes.
if not args.static_only:
    targeted=[
        'TestV1330CacheFingerprintSkipsUnchangedPhysicalRewrite',
        'TestExtreme30_19ConcurrencyAndRaceStress',
        'TestExtreme30_20PerformanceLongRunNoUnboundedGoroutineGrowth',
        'TestV1604StopTimeoutPersistsCurrentCacheBeforeReturning',
        'TestV1604QuoteLevelCatalystTrackingIsEventDriven',
    ]
    for name in targeted:
        run(['go','test','-count=1','-run',f'^{name}$','.'],90,f'G7 targeted regression {name}')

print('DE.PULSE G7 — Code Quality + Maintainability + Performance Health')
if errors:
    for e in errors: print('FAIL:',e)
    sys.exit(1)
if args.static_only:
    print('G7 STATIC PASS — gofmt clean · go vet PASS · bounded HTTP/ticker lifecycle · critical WHY comments present')
else:
    print('G7 PASS — static quality + independently bounded performance/concurrency regressions PASS')
