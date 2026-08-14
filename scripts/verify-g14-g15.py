#!/usr/bin/env python3
import argparse, hashlib, json, pathlib, sys, datetime
EXPECTED_SHA='69f87c3fd2f94e9adc8fe5e3fb273843bd5b07cc019f62eeaf1e4ad688adb66f'
EXPECTED_FP='a1ff742baf04176d6122da338fc2360c0d21c7220977c6901299b755ce2cfc5b'
EXPECTED_BUILD='v18.0.4-test-native-windows-lifecycle-g14-harness-hardening-20260813'
REQUIRED=['sourceSha','sourceFingerprint','nativeCompile','nativeSQLite','packageLauncherExecution','releaseIdentity','stableProfileMigration','stableSecretsMigration','stableIsolation','bootstrapOwner','csrfPasswordSetup','restartPersistence','credentialLoginAfterRestart','smartRouterV2Runtime','rapidMoveRuntime','coverageTruth']

def sha256(p):
    h=hashlib.sha256()
    with open(p,'rb') as f:
        for b in iter(lambda:f.read(1024*1024),b''): h.update(b)
    return h.hexdigest()

def load(root, platform):
    root=pathlib.Path(root); ep=root/'g14-evidence.json'
    if not ep.is_file(): raise SystemExit(f'MISSING: {ep}')
    e=json.loads(ep.read_text(encoding='utf-8-sig'))
    assert e['gate']=='G14' and e['status']=='PASS'
    assert e['platform']==platform
    assert e['sourceSha256']==EXPECTED_SHA and e['sourceFingerprint']==EXPECTED_FP and e['buildId']==EXPECTED_BUILD
    missing=[k for k in REQUIRED if e.get('checks',{}).get(k) is not True]
    if missing: raise SystemExit(f'{platform}: failed/missing checks: {missing}')
    art=root/e['artifact']['name']
    if not art.is_file(): raise SystemExit(f'{platform}: artifact missing: {art}')
    actual=sha256(art)
    if actual!=e['artifact']['sha256']: raise SystemExit(f'{platform}: artifact SHA mismatch {actual}')
    return e,art,actual

def main():
    ap=argparse.ArgumentParser(); ap.add_argument('--mac',required=True); ap.add_argument('--windows',required=True); ap.add_argument('--out',required=True); a=ap.parse_args()
    m,mp,ms=load(a.mac,'macos-arm64'); w,wp,ws=load(a.windows,'windows-x64')
    out=pathlib.Path(a.out); out.mkdir(parents=True,exist_ok=True)
    result={'schemaVersion':1,'gate':'G15','status':'PASS','decision':'READY_FOR_PROMOTION','release':'v18.0.4 TEST','buildId':EXPECTED_BUILD,'sourceSha256':EXPECTED_SHA,'sourceFingerprint':EXPECTED_FP,'nativeTargets':{'macos-arm64':{'status':'PASS','artifact':mp.name,'sha256':ms},'windows-x64':{'status':'PASS','artifact':wp.name,'sha256':ws}},'generatedAt':datetime.datetime.now(datetime.timezone.utc).isoformat(),'note':'G15 evidence says both required native targets passed G14. Stable promotion remains an explicit release action.'}
    (out/'G15-RELEASE-ASSURANCE.json').write_text(json.dumps(result,indent=2)+'\n')
    (out/'G15-RELEASE-ASSURANCE.txt').write_text('PASS: v18.0.4 required native targets passed G14. READY_FOR_PROMOTION.\n')
    print('PASS: G15 native release assurance — READY_FOR_PROMOTION')
if __name__=='__main__': main()
