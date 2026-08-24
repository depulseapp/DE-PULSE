#!/usr/bin/env python3
from pathlib import Path
import json,sys
R=Path(__file__).resolve().parent; e=[]
def need(ok,msg):
    if not ok:e.append(msg)
i=json.loads((R/'release_identity.json').read_text()); s=json.loads((R/'v18_2_scope.json').read_text())
need(i.get('version')=='18.2.0' and i.get('channel') in {'TEST','STABLE'},'v18.2 identity missing')
need(i.get('previous_stable')=='v18.1.0' and i.get('patch_predecessor')=='v18.1.0','v18.1.0 predecessor drift')
if i.get('channel')=='TEST':
    need(i.get('build_id')=='v18.2.0-test-admin-presence-sessions-20260814','TEST build identity drift')
    need(i.get('runtime_config')=='PersonalMarketTerminal-v18.2.0-TEST','isolated TEST runtime missing')
    need(i.get('application_bundle')=='De-Pulse-v18.2.0-TEST.app','TEST bundle missing')
need(s.get('incomingStableCommit')=='4b48c61ee103897e834f4ff7d24af07d4275a62c','G0 commit drift')
need(s.get('incomingStableSourceSha256')=='84a8f2ce4ca7a174b1b3a18cf17904bbbbb4528139c14644b9e07a287e2f0419','G0 source SHA drift')
need(len(s.get('clauses',[]))==12,'scope clause count drift')
a=(R/'identity_admin.go').read_text(); h=(R/'http_admin_identity.go').read_text(); api=(R/'http_api.go').read_text(); ui=(R/'renderer/admin-v18.2.js').read_text()
for token in ['adminSnapshot','adminCreateUser','adminSetUserRole','adminSetUserStatus','adminResetPassword','adminRevokeSession']:
    need(token in a,token+' missing')
for route in ['/api/admin/identity','/api/admin/users/create','/api/admin/users/role','/api/admin/users/status','/api/admin/users/reset-password','/api/admin/sessions/revoke']:
    need(route in api or route in ui,route+' missing')
need('sessionIDActive(p.SessionID)' in api and 'touchSessionID(p.SessionID)' in api,'SSE lifecycle integration missing')
need("['SUPER_OWNER','OWNER','ADMIN']" in ui and 'SESSION TRUTH' in ui,'role/presence UI boundary missing')
need('postgres' not in a.lower() and 'postgres' not in h.lower(),'v18.3 hosted scope leaked into v18.2')
if e:
    print('v18.2 Scope Gate: FAIL'); [print(' -',x) for x in e]; sys.exit(2)
print('v18.2 Scope Gate: PASS · 12/12 clauses · canonical identity/session ownership preserved')
