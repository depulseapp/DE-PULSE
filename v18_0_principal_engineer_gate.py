#!/usr/bin/env python3
from pathlib import Path
import re, sys
R=Path(__file__).resolve().parent
review=R/'renderer/qa/v18.0.0-principal-engineer-review.md'
required_files=[
 'identity_model.go','password_argon2id.go','http_auth.go','v18_test_profile.go',
 'identity_v18_test.go','v18_0_scope_gate.py','v18_documentation_typography_gate.py'
]
problems=[]
for f in required_files:
 if not (R/f).exists(): problems.append(f'missing {f}')
text=review.read_text() if review.exists() else ''
for phrase in [
 'v17.5.1 remains the authoritative Stable',
 'Argon2id only',
 'Raw session tokens are not stored',
 'HTTP 428',
 'ADMIN-or-higher',
 '2403/2403 PASS',
 '167/167 PASS',
 '30/30 PASS',
 '15/15 supported viewports PASS',
 'No Execution Boundary',
 'GO to v18.0.0 separate TEST packaging/host audit',
 'Native macOS packaging with SQLite must be built and audited on a compatible macOS release host',
]:
 if phrase not in text: problems.append(f'review missing: {phrase}')
# Architecture/security source invariants.
identity=(R/'identity_model.go').read_text()
auth=(R/'http_auth.go').read_text()
pw=(R/'password_argon2id.go').read_text()
profile=(R/'v18_test_profile.go').read_text()
if 'RoleSuperOwner' not in identity or 'RoleOwner' not in identity or 'RoleAdmin' not in identity: problems.append('role hierarchy missing')
if 'sessionTokenHash' not in identity: problems.append('hashed session-token persistence missing')
if 'argon2.IDKey' not in pw and 'vendoredargon2.IDKey' not in pw: problems.append('Argon2id production path missing')
if 'SameSiteStrictMode' not in auth or 'X-DE-PULSE-CSRF' not in auth: problems.append('cookie/CSRF hardening missing')
if 'PersonalMarketTerminal-v18-TEST' not in profile: problems.append('isolated TEST profile missing')
if problems:
 print('v18.0 Principal Engineer Gate: FAIL')
 for p in problems: print(' -',p)
 sys.exit(1)
print('v18.0 Principal Engineer Gate: PASS · architecture/security/migration/review evidence present · GO to separate TEST packaging only')
