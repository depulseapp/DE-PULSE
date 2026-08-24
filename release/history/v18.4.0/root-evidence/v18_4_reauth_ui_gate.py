#!/usr/bin/env python3
from pathlib import Path

text = Path('renderer/renderer.js').read_text()
required = {
    'reauth endpoint': "'/api/auth/reauth'",
    'machine reauth code': "data.code==='REAUTH_REQUIRED'",
    'masked password input': 'type="password"',
    'password autocomplete': 'autocomplete="current-password"',
    'single retry after reauth': 'if(await requestRecentAuthentication())',
    'wrong-password inline handling': 'r.status===403||r.status===429',
}
errors = [name for name, needle in required.items() if needle not in text]
if "if(r.status===401||r.status===428){location.replace('/')" not in text:
    errors.append('generic auth redirect remains after dedicated reauth branch')
if errors:
    print('v18.4 reauth UI gate: FAIL')
    for item in errors:
        print(' -', item)
    raise SystemExit(2)
print('v18.4 reauth UI gate: PASS · masked password · machine-coded 428 · bounded single retry')
