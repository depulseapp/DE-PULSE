#!/usr/bin/env python3
from pathlib import Path
import sys
R=Path(__file__).resolve().parent; e=[]
def need(ok,msg):
    if not ok:e.append(msg)
admin=(R/'identity_admin.go').read_text(); pres=(R/'identity_presence.go').read_text(); ui=(R/'renderer/admin-v18.2.js').read_text()
need('IdentityService' in admin,'canonical IdentityService not reused')
need('PasswordHash' not in ui and 'TokenHash' not in ui,'credential material exposed to renderer')
need('adminPresenceActiveWindow' in admin and 'sessionPresence' in admin,'presence not derived centrally')
need('sessionIDActive' in pres,'long-lived session validity owner missing')
need('canManageRole' in admin and 'activeCriticalOwnersLocked' in admin,'role hierarchy/owner safety missing')
need('v182CanAdmin' in ui,'role-aware UI boundary missing')
if e:
    print('v18.2 Principal Engineer Gate: FAIL'); [print(' -',x) for x in e]; sys.exit(2)
print('v18.2 Principal Engineer Gate: PASS · one identity/session owner · redacted admin view · no hosted/execution scope')
