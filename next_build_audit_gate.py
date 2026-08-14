#!/usr/bin/env python3
"""G16 transition gate: authorize a new DE.PULSE minor only from a closed Stable handoff."""
from pathlib import Path
import hashlib,json,os,subprocess,sys
ROOT=Path(__file__).resolve().parent
closure_path=Path(os.environ.get('DEPULSE_PREVIOUS_CLOSURE', ROOT/'renderer/qa/v16.0.4-working-candidate-closure.json'))
errors=[]
if not closure_path.exists():
    errors.append(f'previous-build closure missing: {closure_path}')
else:
    c=json.loads(closure_path.read_text())
    status=str(c.get('status','')).upper()
    if 'CLOSED' not in status or 'STABLE' not in status or status.startswith('NOT '):
        errors.append(f'previous build is not a CLOSED STABLE release: status={c.get("status")}')
    auth=str(c.get('authorization','')).upper()
    if 'V16.4.0' not in auth or 'AUTHORIZED' not in auth or 'NOT AUTHORIZED' in auth:
        errors.append(f'v16.4.0 next minor not explicitly authorized: {c.get("authorization")}')
    if c.get('masterScope',{}).get('knownP0P1FeatureTruthDefects',1)!=0:
        errors.append('previous closure reports P0/P1 feature-truth defects')
    if not c.get('protectedStable',{}).get('untouched'):
        errors.append('closure does not protect Stable')
    stable_path=os.environ.get('DEPULSE_STABLE_SOURCE',''); expected=c.get('protectedStable',{}).get('sourceSha256','')
    if stable_path and expected:
        got=hashlib.sha256(Path(stable_path).read_bytes()).hexdigest()
        if got!=expected: errors.append(f'protected Stable hash changed: {got}')

# Only after closure is valid do we perform a small fresh transition audit. It
# reuses checkpointed checks and never recursively re-runs the heavy G12 matrix.
if not errors:
    ids=['g2_source_health','g10_professional_expert','g11_fresh_adversarial','g11_approved_scope']
    cmd=[sys.executable,'certification_runner.py']
    for cid in ids: cmd += ['--check',cid]
    p=subprocess.run(cmd,cwd=ROOT,text=True)
    if p.returncode:
        errors.append('fresh checkpointed transition audit failed or is blocked')
if errors:
    print('G16 Next-Build Authorization Gate: FAIL')
    print('\n'.join('- '+e for e in errors));sys.exit(1)
print('G16 Next-Build Authorization Gate: PASS · previous Stable closed and next minor explicitly authorized')
