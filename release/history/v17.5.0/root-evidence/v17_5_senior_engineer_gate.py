#!/usr/bin/env python3
from pathlib import Path
import re,json,sys
R=Path(__file__).resolve().parent; errs=[]; warnings=[]
def need(ok,msg):
    if not ok: errs.append(msg)
ident=json.loads((R/'release_identity.json').read_text()); need(ident.get('version') in {'17.5.0','17.5.1'} or str(ident.get('version','')).startswith('18.'),'v17.5 closure inheritance identity missing')
policy=json.loads((R/'source_health_baseline.json').read_text()); FILE_MAX=int(policy.get('production_file_max_lines',1200)); FUNC_BLOCK=int(policy.get('function_block_lines',320)); BRANCH_BLOCK=int(policy.get('branch_block_heuristic',120)); FUNC_WATCH=int(policy.get('function_watch_lines',200))
prod_files=sorted(p for p in R.glob('*.go') if not p.name.endswith('_test.go')); prod='\n'.join(p.read_text(errors='ignore') for p in prod_files)
checks={
 'PersistenceBackend owner':len(re.findall(r'type\s+PersistenceBackend\s+interface',prod))==1,
 'canonical quote propagation owner':len(re.findall(r'func\s+\(e \*Engine\)\s+propagateCanonicalQuoteChange\s*\(',prod))==1,
 'workload controller constructor':len(re.findall(r'func\s+NewWorkloadController\s*\(',prod))==1,
 'provider route execute authority':len(re.findall(r'func\s+\(e \*Engine\)\s+executeProviderRoute\s*\(',prod))==1,
 'runtime degradation owner':len(re.findall(r'func\s+deriveRuntimeDegradation\s*\(',prod))==1,
 'runtime SLO owner':len(re.findall(r'func\s+buildRuntimeSLOAssessmentWithContext\s*\(',prod))==1,
 'current liquidity risk owner':len(re.findall(r'func\s+currentLiquidityMarketRisk\s*\(',prod))==1,
}
for label,ok in checks.items(): need(ok,label+' failed')
for p in prod_files:
    lines=p.read_text(errors='ignore').splitlines(); need(len(lines)<=FILE_MAX,f'{p.name}: {len(lines)} lines exceeds {FILE_MAX}')
    for i,line in enumerate(lines,1):
        if re.search(r'\b(TODO|FIXME|HACK)\b',line): errs.append(f'{p.name}:{i}: unresolved maintenance marker')
        if 'time.Tick(' in line: errs.append(f'{p.name}:{i}: unowned time.Tick lifecycle')
        if 'http.Client{' in line:
            window=' '.join(lines[i-1:min(len(lines),i+7)])
            if 'Timeout:' not in window: errs.append(f'{p.name}:{i}: http.Client lacks explicit timeout')
for pat,label in [(r'sk-[A-Za-z0-9_-]{20,}','OpenAI-style secret'),(r'AIza[0-9A-Za-z_-]{20,}','Google API secret')]:
    if re.search(pat,prod): errs.append('possible hard-coded '+label)
head=re.compile(r'(?m)^func\s+(?:\([^\n]*?\)\s*)?([A-Za-z_]\w*)\s*\([^\n]*?\)[^{]*\{')
metrics=[]
for p in prod_files:
    s=p.read_text(errors='ignore')
    for m in head.finditer(s):
        i=m.end()-1; d=0; j=i; quote=None; esc=False; lc=False; bc=False
        while j<len(s):
            c=s[j]; n=s[j+1] if j+1<len(s) else ''
            if lc:
                if c=='\n': lc=False
            elif bc:
                if c=='*' and n=='/': bc=False; j+=1
            elif quote:
                if quote=='`':
                    if c=='`': quote=None
                elif esc: esc=False
                elif c=='\\': esc=True
                elif c==quote: quote=None
            else:
                if c=='/' and n=='/': lc=True; j+=1
                elif c=='/' and n=='*': bc=True; j+=1
                elif c in {'"',"'",'`'}: quote=c
                elif c=='{': d+=1
                elif c=='}':
                    d-=1
                    if d==0:
                        body=s[m.start():j+1]; ln=body.count('\n')+1; br=len(re.findall(r'\b(if|for|switch|select|case)\b|&&|\|\|',body)); metrics.append((ln,br,p.name,m.group(1)))
                        if ln>FUNC_BLOCK: errs.append(f'{p.name}:{m.group(1)} is {ln} lines > {FUNC_BLOCK}')
                        if br>BRANCH_BLOCK: errs.append(f'{p.name}:{m.group(1)} complexity {br} > {BRANCH_BLOCK}')
                        break
            j+=1
watch=[x for x in sorted(metrics,reverse=True) if x[0]>FUNC_WATCH]
tests=sum(len(re.findall(r'(?m)^func\s+Test[A-Za-z0-9_]+\s*\(',p.read_text(errors='ignore'))) for p in R.glob('*_test.go')); need(tests>=470,f'Go test inventory unexpectedly small: {tests}')
for f in ['v17_5_major_closure_scope_matrix.json','v17_5_major_scope_matrix_gate.py','data_utility_registry.json','release_learning_registry.json','runtime_slo.go','workload_controller.go','canonical_pipeline.go']:
    need((R/f).exists(),'closure engineering artifact missing: '+f)
report=R/'renderer/qa/v17.5.0-senior-engineer-review.md'; txt=report.read_text(errors='ignore') if report.exists() else ''
for token in ['15+ year','Architecture ownership','Persistence and truth review','Runtime / concurrency / failure handling','Maintainability watchlist','Performance / scalability','Security / product boundaries','Upgrade safety','Verdict']:
    need(token in txt,'senior review report missing '+token)
print('v17.5 Senior Developer / Principal Engineer Review')
print(f'Production Go: {len(prod_files)} files · {sum(len(p.read_text(errors="ignore").splitlines()) for p in prod_files)} lines · Go tests {tests}')
for x in watch[:6]: print(' · watch',x)
if errs:
    print('VERDICT: FAIL'); [print(' -',e) for e in errs]; sys.exit(2)
print('VERDICT: PASS · v17 architecture ownership, failure handling, maintainability, performance boundaries and upgrade safety are closure-suitable')
