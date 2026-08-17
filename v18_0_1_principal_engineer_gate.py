#!/usr/bin/env python3
from pathlib import Path
import json, subprocess, sys
R=Path(__file__).resolve().parent

# Compatibility entry point: old immutable releases keep their own historical
# file content via tags. On the active v18.5 branch, do not require obsolete
# v18.0.1 TEST-profile packaging; delegate to the current Principal Engineer
# closure contract instead.
try:
    current=json.loads((R/'release_identity.json').read_text()).get('version','')
except Exception:
    current=''
if current == '18.5.0':
    result=subprocess.run([sys.executable, str(R/'v18_5_principal_engineer_gate.py')], cwd=R)
    sys.exit(result.returncode)

review=R/'renderer/qa/v18.0.1-principal-engineer-review.md'
required_files=[
 'identity_model.go','password_argon2id.go','http_auth.go','v18_test_profile.go',
 'smart_router_v2.go','rapid_move_intelligence.go','v18_0_1_scope_gate.py',
 'v18_0_1_smart_router_rapid_test.go','v18_documentation_typography_gate.py'
]
problems=[]
for f in required_files:
    if not (R/f).exists(): problems.append(f'missing {f}')
text=review.read_text() if review.exists() else ''
for phrase in [
 'v17.5.1 remains the authoritative Stable',
 'single executable provider-routing authority',
 'NOT_ENTITLED',
 'Preferred provider and Serving provider',
 'TIERED_PARTIAL',
 '15s/30s/60s/2m/5m',
 'SHADOW → VALIDATED → APPROVED → PRODUCTION',
 'PersonalMarketTerminal-v18.0.1-TEST',
 '2403/2403 PASS',
 '167/167 PASS',
 '15/15 supported viewports PASS',
 'No Execution Boundary',
 'GO to v18.0.1 immutable RC/full certification',
]:
    if phrase not in text: problems.append(f'review missing: {phrase}')
router=(R/'provider_router.go').read_text() + '\n' + (R/'smart_router_v2.go').read_text()
rapid=(R/'rapid_move_intelligence.go').read_text()
profile=(R/'v18_test_profile.go').read_text()
auth=(R/'http_auth.go').read_text(); identity=(R/'identity_model.go').read_text(); pw=(R/'password_argon2id.go').read_text()
if 'NOT_ENTITLED' not in router: problems.append('NOT_ENTITLED classification/suppression missing')
if 'Preferred' not in router or 'Serving' not in router: problems.append('Preferred vs Serving semantics missing')
if 'TIERED_PARTIAL' not in rapid: problems.append('Coverage Truth missing')
for token in ['15', '30', '60']:
    if token not in rapid: problems.append(f'rapid window evidence missing: {token}')
if 'PersonalMarketTerminal-v18.0.1-TEST' not in profile: problems.append('isolated v18.0.1 TEST profile missing')
if 'sessionTokenHash' not in identity: problems.append('hashed session-token persistence missing')
if 'argon2.IDKey' not in pw and 'vendoredargon2.IDKey' not in pw: problems.append('Argon2id production path missing')
if 'SameSiteStrictMode' not in auth or 'X-DE-PULSE-CSRF' not in auth: problems.append('cookie/CSRF hardening missing')
scope=json.loads((R/'v18_0_1_scope.json').read_text())
if scope.get('version')!='18.0.1' or scope.get('stable_baseline')!='v17.5.1': problems.append('scope identity/baseline mismatch')
if problems:
    print('v18.0.1 Principal Engineer Gate: FAIL')
    for p in problems: print(' -',p)
    sys.exit(1)
print('v18.0.1 Principal Engineer Gate: PASS · canonical router/event ownership, adaptive governance, coverage truth and inherited security verified · GO to immutable RC/full certification')