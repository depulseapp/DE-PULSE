#!/usr/bin/env python3
from pathlib import Path
import re,subprocess,sys
R=Path(__file__).resolve().parent
p=subprocess.run(['go','test','-run','^$','-bench','BenchmarkV168','-benchmem','-benchtime=100x','.'],cwd=R,text=True,capture_output=True)
out=p.stdout+p.stderr
if p.returncode:
 print('v16.8 Performance Gate: FAIL\n'+out[-4000:]); sys.exit(1)
rows=re.findall(r'^(BenchmarkV168\S+)-\d+\s+\d+\s+([0-9.]+) ns/op\s+([0-9.]+) B/op\s+(\d+) allocs/op',out,re.M)
if len(rows)<2:
 print('v16.8 Performance Gate: FAIL · benchmark output unreadable\n'+out[-4000:]); sys.exit(1)
limits={'BenchmarkV168HeatMap500Symbols':(150_000_000,12_000_000),'BenchmarkV168Liquidity500Symbols':(250_000_000,30_000_000)}
errs=[]
for name,ns,bts,alloc in rows:
 ns=float(ns); bts=float(bts); key=name.split('-')[0]
 lim=limits.get(key)
 if lim and (ns>lim[0] or bts>lim[1]): errs.append(f'{key} regression ns/op={ns} B/op={bts}')
if errs:
 print('v16.8 Performance Gate: FAIL'); print('\n'.join('- '+x for x in errs)); sys.exit(1)
print('v16.8 Performance Gate: PASS · 500-symbol heat/liquidity benchmark within bounded release ceilings')
for r in rows: print(' ·',r[0],r[1]+' ns/op',r[2]+' B/op',r[3]+' allocs/op')
