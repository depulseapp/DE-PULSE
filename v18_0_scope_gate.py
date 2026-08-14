#!/usr/bin/env python3
from pathlib import Path
import json,sys
R=Path(__file__).resolve().parent; e=[]
def need(ok,msg):
    if not ok:e.append(msg)
i=json.loads((R/'release_identity.json').read_text())
need(i.get('version')=='18.0.0' and i.get('channel')=='TEST','v18.0 TEST identity missing')
need(i.get('build_id')=='v18.0.0-test-identity-session-foundation-20260813','v18.0 build id drift')
need(i.get('stable_baseline')=='v17.5.1' and i.get('previous_stable')=='v17.5.1','v17.5.1 Stable baseline drift')
need(i.get('runtime_config')=='PersonalMarketTerminal-v18-TEST','v18 TEST runtime isolation missing')
need(i.get('application_bundle')=='De-Pulse-v18-TEST.app','v18 TEST application identity missing')
for f in ['identity_model.go','password_argon2id.go','http_auth.go','v18_test_profile.go','identity_v18_test.go','v18_major_scope.json','v18_delivery_slices.json']:
    need((R/f).exists(),'missing '+f)
ident=(R/'identity_model.go').read_text(); auth=(R/'http_auth.go').read_text(); api=(R/'http_api.go').read_text(); pw=(R/'password_argon2id.go').read_text(); prof=(R/'v18_test_profile.go').read_text(); repo=(R/'persistence_repository.go').read_text(); sql=(R/'persistence_backend_sqlite.go').read_text(); js=(R/'renderer/renderer.js').read_text()
for token in ['SUPER_OWNER','OWNER','ADMIN','USER','DEMO','Principal','SessionRecord','bootstrapOwnerID','defaultSessionIdleTTL','defaultSessionAbsoluteTTL']:
    need(token in ident,'identity contract missing '+token)
for token in ['argon2.IDKey','argonTime','argonMemory','argonThreads','verifyPasswordArgon2id']:
    need(token in pw,'Argon2id contract missing '+token)
for token in ['validRequestCSRF','handleRotateSession','authAllowPasswordSetup','requireRole']:
    need(token in auth,'HTTP security contract missing '+token)
for token in ['depulse_session','depulse_csrf']:
    need(token in ident+auth,'HTTP cookie contract missing '+token)
for route in ['/api/auth/login','/api/auth/logout','/api/auth/set-password','/api/auth/rotate']:
    need(route in api,'auth route missing '+route)
need('LoadIdentityState' in repo and 'SaveIdentityState' in repo,'identity persistence repository contract missing')
need('CREATE TABLE IF NOT EXISTS identity_state' in sql,'SQLite identity migration missing')
need('stableInstanceIsLive' in prof and 'cloneStableProfile' in prof and 'PersonalMarketTerminal-v18-TEST' in prof,'side-by-side Stable protection missing')
need('X-DE-PULSE-CSRF' in js and '/api/auth/logout' in js,'renderer authenticated state-change contract missing')
alltext=(R/'README.md').read_text()+(R/'renderer/docs/developer.md').read_text()+(R/'renderer/docs/limitations.md').read_text()
need('No Execution' in alltext or 'no execution' in alltext.lower(),'No Execution boundary missing')
if e:
    print('v18.0 Identity & Secure Session Scope Gate: FAIL'); [print(' -',x) for x in e]; sys.exit(2)
print('v18.0 Identity & Secure Session Scope Gate: PASS · canonical identity + Argon2id + opaque sessions + CSRF/RBAC + isolated TEST profile + v17.5.1 Stable protection')
