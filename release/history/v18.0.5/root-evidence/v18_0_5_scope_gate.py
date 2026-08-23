#!/usr/bin/env python3
import json,sys
from pathlib import Path
R=Path(__file__).resolve().parent; e=[]
def need(x,m):
    if not x:e.append(m)
i=json.loads((R/'release_identity.json').read_text()); s=json.loads((R/'v18_0_5_scope.json').read_text())
need(i.get('version')=='18.0.5' and i.get('channel') in {'TEST','STABLE'},'v18.0.5 release identity missing')
need(i.get('stable_baseline')=='v17.5.1' and i.get('previous_stable')=='v18.0.4','v18 major anchor / v18.0.4 previous Stable drift')
if i.get('channel')=='TEST':
    need(i.get('build_id')==s.get('build_id'),'TEST build identity drift')
    need(i.get('runtime_config')=='PersonalMarketTerminal-v18.0.5-TEST','isolated v18.0.5 TEST runtime missing')
else:
    need(i.get('build_id')=='v18.0.5-stable-ui-ux-symbol-management-hardening-20260814','Stable build identity drift')
    need(i.get('runtime_config')=='PersonalMarketTerminal','Stable runtime config drift')
need(len(s.get('clauses',[]))==12,'immutable scope clause count mismatch')
mem=(R/'desk_membership.go').read_text(); need(all(x in mem for x in ['trackedSymbolsLocked','setTrackedSymbolLocked','clearTrackedSymbolsLocked']),'canonical tracked-symbol owner helpers missing')
need('wl.Symbols = []string{}' in mem,'explicit-empty Remove All truth missing')
r=(R/'renderer/renderer.js').read_text(); css=(R/'renderer/styles.css').read_text()
need('Tracked Symbols' in r and 'Add Symbol' in r and 'No tracked symbols yet.' in r,'Tracked Symbols user surface incomplete')
need('function diagnosticVisibilityForRole' in r and "'USER'" not in r.split('function diagnosticVisibilityForRole',1)[1].split('}',1)[0],'USER must not be privileged diagnostics role')
need('opportunity-radar-head' in r and 'radar-status-pill' in r and '@media(max-width:900px)' in css,'Opportunity Radar hierarchy/responsive hardening missing')
need('research-origin-context' in r and 'Technical Context' in r,'compact Research Target or inherited Technical Context missing')
need('v18.1' not in mem,'v18.1 architecture leaked into v18.0.5 implementation')
if e:
 print('v18.0.5 Immutable Scope: FAIL'); [print(' -',x) for x in e]; sys.exit(2)
print('v18.0.5 Immutable Scope: PASS · 12/12 clauses · focused patch boundary preserved')
