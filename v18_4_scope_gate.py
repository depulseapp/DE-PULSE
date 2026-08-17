#!/usr/bin/env python3
from pathlib import Path
import json, subprocess
R=Path(__file__).resolve().parent
errors=[]
def need(ok,msg):
    if not ok: errors.append(msg)
scope=json.loads((R/'v18_4_scope.json').read_text())
contract=json.loads((R/'v18_4_g0_g3_contract.json').read_text())
need(scope.get('version')=='18.4.0','v18.4 version drift')
need(scope.get('incomingStable')=='v18.3.0','incoming Stable drift')
need(scope.get('incomingStableCommit')=='35ef2a161e6a01c4f22fd0c9aecbe8920a7dace2','G0 Stable commit drift')
need(scope.get('incomingStableFingerprint')=='d33c591287aae2c6b07b2be082494aaa76277a327f7d6ea01e87cae1dcbcd272','G0 Stable fingerprint drift')
expected=[
 'secrets_security','auth_session_csrf_cookie_hardening','adversarial_authorization',
 'provider_entitlement_data_rights_metadata','hosted_commercial_readiness',
 'quota_abuse_safeguards','observability','licensing_redistribution_ai_use_suitability']
need(scope.get('clauses')==expected,'G1 clause set/order drift')
need(contract.get('G0',{}).get('commit')==scope.get('incomingStableCommit'),'G0 contract/scope mismatch')
for f in ('release/v18.4.0/G0-EXACT-BASELINE.md','release/v18.4.0/G1-IMMUTABLE-SCOPE.md','release/v18.4.0/G2-ARCHITECTURE-DATA-UTILITY.md','release/v18.4.0/G3-DESIGN-DEPENDENCY-READINESS.md'):
    need((R/f).exists(),f+' missing')
protected=['scanner.go','preparation_types_liquidity.go','validation_learning.go','signal_validation.go','persistence_backend_postgres.go','persistence_repository.go']
if (R/'.git').exists():
    try:
        changed=subprocess.check_output(['git','diff','--name-only',scope['incomingStableCommit']+'...HEAD','--',*protected],cwd=R,text=True).splitlines()
        need(not changed,'protected v18.3/market owner drift: '+', '.join(changed))
    except Exception as exc:
        errors.append('protected-source comparison failed: '+str(exc))
if errors:
    print('v18.4 scope gate: FAIL')
    for e in errors: print(' -',e)
    raise SystemExit(2)
print('v18.4 scope gate: PASS · G0-G3 frozen · 8/8 clauses · protected v18.3/market owners unchanged')
