#!/usr/bin/env python3
import argparse,hashlib,json,pathlib,datetime
SHA='857c7999891225f5ac8dc26eaf1f2018dfd72ccbf64ef23b11a103e169679199'; FP='5f5285071f4a46ed6ab83026c158d58936eecdea38961940622faa8cfe4d7f90'; BUILD='v18.0.6-stable-smart-provider-router-rapid-move-market-shock-hardening-20260814'
REQ=['sourceSha','sourceFingerprint','nativeCompile','nativeSQLite','packageLauncherExecution','releaseIdentity','stableRuntimeContinuity','stableSecretsPreserved','bootstrapOwner','csrfPasswordSetup','restartPersistence','credentialLoginAfterRestart','smartRouterV2Runtime','rapidMoveRuntime','coverageTruth']
def h(p):
 x=hashlib.sha256(); f=open(p,'rb')
 for b in iter(lambda:f.read(1048576),b''): x.update(b)
 f.close(); return x.hexdigest()
def load(root,plat):
 r=pathlib.Path(root); e=json.loads((r/'g14-evidence.json').read_text(encoding='utf-8-sig'))
 assert e['gate']=='G14' and e['status']=='PASS' and e['release']=='v18.0.6 STABLE' and e['platform']==plat and e['sourceSha256']==SHA and e['sourceFingerprint']==FP and e['buildId']==BUILD
 miss=[k for k in REQ if e.get('checks',{}).get(k) is not True]
 if miss: raise SystemExit(f'{plat}: missing {miss}')
 a=r/e['artifact']['name']; actual=h(a)
 if actual!=e['artifact']['sha256']: raise SystemExit(f'{plat}: artifact SHA mismatch')
 return a,actual
def main():
 p=argparse.ArgumentParser(); p.add_argument('--mac',required=True); p.add_argument('--windows',required=True); p.add_argument('--out',required=True); a=p.parse_args(); m,ms=load(a.mac,'macos-arm64'); w,ws=load(a.windows,'windows-x64'); o=pathlib.Path(a.out); o.mkdir(parents=True,exist_ok=True)
 j={'schemaVersion':1,'gate':'G15','status':'PASS','decision':'STABLE_CERTIFIED','release':'v18.0.6 STABLE','buildId':BUILD,'sourceSha256':SHA,'sourceFingerprint':FP,'nativeTargets':{'macos-arm64':{'status':'PASS','artifact':m.name,'sha256':ms},'windows-x64':{'status':'PASS','artifact':w.name,'sha256':ws}},'generatedAt':datetime.datetime.now(datetime.timezone.utc).isoformat(),'note':'User-authorized v18.0.6 Stable promotion. Both required native targets passed G14 against the exact Stable source fingerprint.'}
 (o/'G15-STABLE-RELEASE-ASSURANCE.json').write_text(json.dumps(j,indent=2)+'\n'); (o/'G15-STABLE-RELEASE-ASSURANCE.txt').write_text('PASS: DE.PULSE v18.0.6 STABLE native release assurance. STABLE_CERTIFIED.\n'); print('PASS: G15 v18.0.6 STABLE release assurance — STABLE_CERTIFIED')
if __name__=='__main__': main()
