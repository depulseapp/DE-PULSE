#!/usr/bin/env python3
"""DE.PULSE G2 — whole-source health, architecture fit and reuse audit."""
from __future__ import annotations
from collections import Counter, defaultdict
from pathlib import Path
import hashlib,re,sys
ROOT=Path(__file__).resolve().parent
import json
HEALTH_POLICY=json.loads((ROOT/'source_health_baseline.json').read_text())
FILE_MAX=int(HEALTH_POLICY.get('production_file_max_lines',1200))
PROD=sorted(p for p in ROOT.glob('*.go') if not p.name.endswith('_test.go'))
errors=[]
def fail(x): errors.append(x)
prod_text='\n'.join(p.read_text(errors='ignore') for p in PROD)
# Count identifiers once. The previous gate rescanned the entire production corpus
# for every declared function, which became unnecessarily expensive as DE.PULSE grew.
# This preserves the same whole-corpus name-reference semantics while keeping G2 CI bounded.
prod_identifier_counts=Counter(re.findall(r'\b[A-Za-z_]\w*\b', prod_text))
# Maintenance debt and module budget.
for p in PROD:
    lines=p.read_text(errors='ignore').splitlines()
    if len(lines)>FILE_MAX: fail(f'{p.name}: {len(lines)} lines exceeds {FILE_MAX:,}-line responsibility budget')
    for i,line in enumerate(lines,1):
        if re.search(r'\b(TODO|FIXME|HACK)\b',line): fail(f'{p.name}:{i}: unresolved maintenance marker')
# Unreferenced production functions/methods.
func_decl=re.compile(r'(?m)^func\s+(?:\([^\n]*?\)\s*)?([A-Za-z_]\w*)\s*\(')
for p in PROD:
    for m in func_decl.finditer(p.read_text(errors='ignore')):
        name=m.group(1)
        if name in {'main','init'}: continue
        if prod_identifier_counts.get(name, 0)==1:
            fail(f'{p.name}: production helper {name} has no production reference')
# Exact duplicate sizeable Go function bodies.
def bodies(path):
    s=path.read_text(errors='ignore'); head=re.compile(r'(?m)^func\s+(?:\([^\n]*?\)\s*)?([A-Za-z_]\w*)\s*\([^\n]*?\)[^{]*\{')
    for m in head.finditer(s):
        i=m.end()-1;d=0;j=i;quote=None;esc=False;lc=False;bc=False
        while j<len(s):
            c=s[j];n=s[j+1] if j+1<len(s) else ''
            if lc:
                if c=='\n': lc=False
            elif bc:
                if c=='*' and n=='/':bc=False;j+=1
            elif quote:
                if quote=='`':
                    if c=='`':quote=None
                elif esc:esc=False
                elif c=='\\':esc=True
                elif c==quote:quote=None
            else:
                if c=='/' and n=='/':lc=True;j+=1
                elif c=='/' and n=='*':bc=True;j+=1
                elif c in {'"',"'",'`'}:quote=c
                elif c=='{':d+=1
                elif c=='}':
                    d-=1
                    if d==0:
                        b=s[i+1:j];b=re.sub(r'//.*','',b);b=re.sub(r'/\*.*?\*/','',b,flags=re.S);b=re.sub(r'\s+',' ',b).strip()
                        if len(b)>=180:yield m.group(1),b
                        break
            j+=1
dup=defaultdict(list)
for p in PROD:
    for name,b in bodies(p):dup[hashlib.sha256(b.encode()).hexdigest()].append((p.name,name))
for rows in dup.values():
    if len(rows)>1:
        rowset=set(rows)
        platform_parity={('persistence_backend_sqlite.go','Stats'),('persistence_backend_windows.go','Stats')}
        if rowset == platform_parity:
            continue
        fail('exact duplicate production function bodies: '+', '.join(f'{p}:{n}' for p,n in rows))

# Adaptive Build Process v2 / release metadata truth.
for required in ('release_identity.json','data_utility_registry.json','data_health_policy.json','prefreeze_qualification.py'):
    if not (ROOT/required).exists(): fail('Adaptive Build Process v2 artifact missing: '+required)
try:
    import json
    ident=json.loads((ROOT/'release_identity.json').read_text())
    ci=json.loads((ROOT/'ci_pipeline_plan.json').read_text())
    if ci.get('version') != ident.get('version'): fail('CI pipeline release version drift')
    if ci.get('policy',{}).get('baseline') != ident.get('previous_stable','')+' Stable': fail('CI pipeline baseline drift')
    if not ci.get('policy',{}).get('pre_freeze_qualification'): fail('pre-freeze qualification policy missing')
except Exception as exc:
    fail('release metadata validation failed: '+str(exc))

# Renderer orphan/duplicate checks.
js=(ROOT/'renderer/renderer.js').read_text(errors='ignore')
js_identifier_counts=Counter(re.findall(r'(?<![A-Za-z0-9_$])[A-Za-z_$][\w$]*(?![A-Za-z0-9_$])', js))
js_decl=re.compile(r'(?m)^function\s+([A-Za-z_$][\w$]*)\s*\(')
for m in js_decl.finditer(js):
    name=m.group(1)
    if js_identifier_counts.get(name, 0)==1: fail(f'renderer.js: named renderer helper {name} has no production reference')
# Canonical ownership / retired paths.
for token,why in {
    'refreshTwelveFundamentalFallback':'retired Twelve fundamentals fallback',
    '"twelve-fundamentals"':'retired Twelve fundamentals health key',
    '/institutional_holders?symbol=':'retired Twelve institutional-holder fallback',
}.items():
    if token in prod_text: fail(f'{why} returned: {token}')
if not re.search(r'"Fundamentals"\s*:\s*\{"Finnhub",\s*"SEC",\s*"yfinance"\}',prod_text): fail('Fundamentals route must remain Finnhub -> SEC -> yfinance')
if len(re.findall(r'func\s+\(e \*Engine\)\s+executeProviderRoute\s*\(',prod_text))!=1: fail('Provider Router must have exactly one executeProviderRoute authority')
if len(re.findall(r'var\s+canonicalSymbolClassifications\s*=',prod_text))!=1: fail('Market Intelligence classification must have exactly one backend canonical owner')
if 'const sectorETF=' in js: fail('renderer-owned ticker/sector map returned; classification must remain backend-owned')
print('DE.PULSE G2 — Source Health + Architecture Fit + Reuse Audit')
print(f'Production Go files: {len(PROD)}')
print(f'Production Go lines: {sum(len(p.read_text(errors="ignore").splitlines()) for p in PROD)}')
print('Orphan production Go helpers: 0' if not any('production helper' in e for e in errors) else 'Orphan production Go helpers: FAIL')
print('Orphan named renderer helpers: 0' if not any('renderer helper' in e for e in errors) else 'Orphan named renderer helpers: FAIL')
print('Exact duplicate sizeable Go bodies: 0' if not any('duplicate production' in e for e in errors) else 'Exact duplicate sizeable Go bodies: FAIL')
if errors:
    for e in errors: print('FAIL:',e)
    sys.exit(1)
print('G2 PASS — minimum necessary code · one canonical owner per responsibility · REUSE/CONSOLIDATE/REFACTOR/DELETE before ADD')
