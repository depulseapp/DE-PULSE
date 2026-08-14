#!/usr/bin/env python3
from pathlib import Path
import re,subprocess,sys
R=Path(__file__).resolve().parent
p=subprocess.run(['go','test','-run','^$','-bench','BenchmarkV1610','-benchmem','-benchtime=50x','.'],cwd=R,text=True,capture_output=True,timeout=120)
out=p.stdout+p.stderr
if p.returncode:
 print('v16.10 Performance Gate: FAIL\n'+out[-4000:]); sys.exit(1)
rows=re.findall(r'^(BenchmarkV1610\S+)-\d+\s+\d+\s+([0-9.]+) ns/op\s+([0-9.]+) B/op\s+(\d+) allocs/op',out,re.M)
if len(rows)<2:
 print('v16.10 Performance Gate: FAIL · benchmark output unreadable\n'+out[-4000:]); sys.exit(1)
# Broad ceilings intentionally guard catastrophic regressions without encoding host-specific micro-optimizations.
limits={
 'BenchmarkV1610OpportunityRadar500Candidates':(50_000_000,20_000_000),
 'BenchmarkV1610OpportunityMetricEnrichment500Candidates':(150_000_000,40_000_000),
}
errs=[]
for name,ns,bts,alloc in rows:
 ns=float(ns); bts=float(bts); key=name.split('-')[0]; lim=limits.get(key)
 if lim and (ns>lim[0] or bts>lim[1]): errs.append(f'{key} regression ns/op={ns} B/op={bts}')
if errs:
 print('v16.10 Performance Gate: FAIL'); print('\n'.join('- '+x for x in errs)); sys.exit(1)
print('v16.10 Performance Gate: PASS · bounded 500-candidate Radar selection/enrichment')
for r in rows: print(' ·',r[0],r[1]+' ns/op',r[2]+' B/op',r[3]+' allocs/op')
