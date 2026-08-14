#!/usr/bin/env python3
from pathlib import Path
import json,re,sys
R=Path(__file__).resolve().parent
e=[]
def need(ok,msg):
    if not ok:e.append(msg)
canonical=['dashboard','market-intelligence','day','swing','long','discovery','research','ai','maintenance','settings','documentation']
expected_audit=['dashboard','market-intelligence','day','swing','long','discovery','research','ai','administration','maintenance','settings','documentation']
html=(R/'renderer/index.html').read_text()
nav=re.findall(r'<button class="nav(?: active)?" data-page="([^"]+)"',html)
need(nav==canonical,f'canonical nav drift: {nav}')
audit=json.loads((R/'renderer/qa/role-aware-tab-audit-v18.2.json').read_text())
tabs=audit.get('tabs',[])
need([x.get('id') for x in tabs]==expected_audit,'role audit does not cover canonical tabs plus conditional Administration in intended order')
need(audit.get('personas')==['OWNER','ADMIN','USER','DEMO'],'role persona matrix drift')
layout=audit.get('layoutRules',{})
for key in ['removeUnauthorizedFromLayout','noReservedGeometry','removeEmptySectionShells','naturalGridReflow','noOrphanHeadersDividers','primaryHierarchyIndependentOfRole','privilegedContentNotAutoPromoted','noPageOverflowOverlapClipping','keyboardFocusOrderMustRemainContinuous']:
    need(layout.get(key) is True,key+' missing')
by={x.get('id'):x for x in tabs}
need(by['administration'].get('kind')=='conditional-privileged','Administration must remain conditional/privileged, not a universal tab')
need(by['administration'].get('OWNER')=='SHOW_FULL' and by['administration'].get('ADMIN')=='SHOW_DELEGATED','Administration role boundary drift')
need(by['administration'].get('USER')=='HIDE_AND_REJECT' and by['administration'].get('DEMO')=='HIDE_AND_REJECT','Administration must be absent/rejected for USER/DEMO')
need(by['maintenance'].get('USER')=='HIDE_AND_REJECT' and by['maintenance'].get('DEMO')=='HIDE_AND_REJECT','Maintenance must be hidden/rejected for USER/DEMO')
need(by['maintenance'].get('ADMIN')=='HIDE_UNLESS_DELEGATED','ADMIN Maintenance must not inherit OWNER diagnostics by role name alone')
need(by['settings'].get('USER')=='COMPOSE_PERSONAL' and by['settings'].get('ADMIN')=='COMPOSE_PERSONAL_AND_EXPLICIT_DELEGATIONS','Settings role composition drift')
shell=audit.get('globalShell',{})
need(shell.get('runtimeStartStop',{}).get('ADMIN')=='HIDE_UNLESS_DELEGATED','ADMIN runtime control must be explicit-delegation only')
need(shell.get('runtimeStartStop',{}).get('USER')=='HIDE' and shell.get('runtimeStartStop',{}).get('DEMO')=='HIDE','shared runtime control must not be user/demo visible')
need(shell.get('dataEngineSidebar',{}).get('ADMIN')=='HIDE_UNLESS_DELEGATED','ADMIN Data Engine must be explicit-delegation only')
need(shell.get('dataEngineSidebar',{}).get('USER')=='HIDE' and shell.get('dataEngineSidebar',{}).get('DEMO')=='HIDE','Data Engine must stay hidden from USER/DEMO')
g9=audit.get('g9',{})
need(g9.get('required') is True and g9.get('viewportFamilies')==['desktop','tablet','narrow-browser'],'G9 role-layout matrix missing')
contract=(R/'adaptive-governance/ROLE_AWARE_UI_COMPOSITION_CONTRACT.md').read_text()
for phrase in ['Role changes composition, not information hierarchy','Seamless composition and no-gap rule','G9 role-aware layout audit','Do not design the OWNER page and subtract cards']:
    need(phrase in contract,'contract phrase missing: '+phrase)
if e:
    print('Role-Aware UI Composition Gate: FAIL')
    [print(' -',x) for x in e]
    sys.exit(2)
print('Role-Aware UI Composition Gate: PASS · 11 canonical + conditional Administration · OWNER/ADMIN/USER/DEMO · hierarchy/no-gap contract frozen')
