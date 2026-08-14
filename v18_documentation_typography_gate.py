#!/usr/bin/env python3
from pathlib import Path
import re,sys
R=Path(__file__).resolve().parent
errors=[]
def need(ok,msg):
    if not ok: errors.append(msg)

docs={
 'renderer/docs/user.md':'DE.PULSE — User documentation',
 'renderer/docs/developer.md':'DE.PULSE — Developer documentation',
 'renderer/docs/limitations.md':'DE.PULSE — Capabilities & Limitations',
}
for name,title in docs.items():
    lines=(R/name).read_text().splitlines()
    h1=[x for x in lines if re.match(r'^#\s+',x)]
    need(len(h1)==1,f'{name}: expected exactly one H1, found {len(h1)}')
    need(bool(lines) and lines[0]==f'# {title}',f'{name}: page H1/title drift')
    need(any(re.match(r'^##\s+v18\.0\.0 TEST',x) for x in lines),f'{name}: v18 TEST section must be H2')
    need(not any(re.match(r'^#\s+v\d',x) for x in lines),f'{name}: release/version heading incorrectly uses H1')

css=(R/'renderer/styles.css').read_text()
need('.doc-content h1{font-size:21px!important' in css,'desktop documentation H1 target missing')
need('.doc-content h2{font-size:15px!important' in css,'desktop documentation H2 target missing')
need('.doc-content h3{font-size:12.5px!important' in css,'desktop documentation H3 target missing')
need('@media(max-width:760px){.doc-content h1{font-size:18px!important}.doc-content h2{font-size:14px!important}.doc-content h3{font-size:12px!important}}' in css,'mobile documentation hierarchy missing')
renderer=(R/'renderer/renderer.js').read_text()
need("out.push(`<h1>" in renderer and "out.push(`<h2>" in renderer and "out.push(`<h3>" in renderer,'renderer Markdown heading semantics missing')
if errors:
    print('v18 Documentation Typography Hierarchy: FAIL')
    for e in errors: print(' -',e)
    sys.exit(2)
print('v18 Documentation Typography Hierarchy: PASS · exactly one page H1 · release H2 · nested H3 · responsive scale')
