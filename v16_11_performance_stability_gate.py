#!/usr/bin/env python3
from pathlib import Path
import re,subprocess,sys
R=Path(__file__).resolve().parent
errs=[]
def run(cmd,timeout,label):
    try:
        p=subprocess.run(cmd,cwd=R,text=True,capture_output=True,timeout=timeout)
    except subprocess.TimeoutExpired as ex:
        errs.append(label+' timeout'); return ''
    out=p.stdout+p.stderr
    if p.returncode: errs.append(label+' failed:\n'+out[-3500:])
    return out
run(['go','test','-count=1','-run','^TestV1611MajorClosure','./...'],90,'v16.11 long-run/capacity regressions')
out=run(['go','test','-run','^$','-bench','^BenchmarkV1611MajorClosurePromotionSelection500$','-benchmem','-benchtime=100x','.'],90,'v16.11 500-candidate closure benchmark')
m=re.search(r'BenchmarkV1611MajorClosurePromotionSelection500-\d+\s+\d+\s+([0-9.]+) ns/op\s+([0-9.]+) B/op\s+(\d+) allocs/op',out)
if not m: errs.append('v16.11 benchmark output unreadable')
else:
    ns=float(m.group(1)); bts=float(m.group(2)); alloc=int(m.group(3))
    # Catastrophic-regression ceilings, intentionally much wider than observed host values.
    if ns>10_000_000: errs.append(f'promotion selection latency regression {ns} ns/op')
    if bts>5_000_000: errs.append(f'promotion selection allocation regression {bts} B/op')
    if alloc>1000: errs.append(f'promotion selection alloc-count regression {alloc}')
if errs:
    print('v16.11 Major Performance / Capacity / Long-Run Gate: FAIL')
    for e in errs: print(' -',e)
    raise SystemExit(2)
print('v16.11 Major Performance / Capacity / Long-Run Gate: PASS')
if m: print(f' · Radar promotion selection 500 candidates: {m.group(1)} ns/op · {m.group(2)} B/op · {m.group(3)} allocs/op')
print(' · 10,000-cycle bounded promotion state + 2,000-symbol rotation/cadence stress PASS')
