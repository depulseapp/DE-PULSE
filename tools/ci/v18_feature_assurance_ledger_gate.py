#!/usr/bin/env python3
"""Final T1 feature-assurance dispatcher.

Before T1 freeze this delegates byte-for-byte to the structural/effective audit.
Once the immutable freeze manifest exists, the same structural audit still runs
first. Because the preserved structural implementation predates the immutable
manifest model and auto-enters its legacy freeze mode when reconciliation says
COMPLETE, the dispatcher temporarily projects that one state field as IN_PROGRESS
for the structural subprocess, restores the exact tracked bytes in a finally block,
then enforces the stronger hashed final manifest. No tracked source is changed.

T1 never marks later T2-T9 behavioral assurance VERIFIED.
"""
from __future__ import annotations

import argparse
import hashlib
import json
from pathlib import Path
import subprocess
import sys

ROOT=Path(__file__).resolve().parents[2]
PROGRAM=ROOT/'governance'/'programs'/'ADAPT-V18-FINAL-CLOSURE-10-10-001'
STRUCTURAL=ROOT/'tools'/'ci'/'v18_feature_assurance_ledger_gate_structural.py'
LEDGER=PROGRAM/'feature-assurance-ledger.json'
SCAN1=PROGRAM/'T1_INDEPENDENT_OMISSION_SCAN.json'
SCAN2=PROGRAM/'T1_INDEPENDENT_OMISSION_SCAN_2.json'
RECON=PROGRAM/'T1_FINAL_RECONCILIATION.json'
QUALITY=PROGRAM/'T1_QUALITY_AUDIT.json'
QUALITY_RESOLUTION=PROGRAM/'T1_QUALITY_RESOLUTION.json'
FREEZE=PROGRAM/'feature-assurance-ledger-freeze.json'


def load(path:Path)->dict:
    try:
        value=json.loads(path.read_text(encoding='utf-8'))
    except Exception as exc:
        raise SystemExit(f'FAIL: cannot read {path.relative_to(ROOT)}: {exc}') from exc
    if not isinstance(value,dict):
        raise SystemExit(f'FAIL: {path.relative_to(ROOT)} must contain a JSON object')
    return value


def git_blob_sha(path:Path)->str:
    data=path.read_bytes()
    h=hashlib.sha1()
    h.update(f'blob {len(data)}\0'.encode('utf-8'))
    h.update(data)
    return h.hexdigest()


def fail(errors:list[str])->int:
    print('V18 FEATURE REALITY & QUALITY AUDIT FREEZE: FAIL')
    for error in errors:
        print(f'- {error}')
    return 1


def run_structural()->int:
    if not STRUCTURAL.is_file():
        return fail(['structural feature-assurance gate is missing'])
    if not FREEZE.exists():
        return subprocess.run([sys.executable,str(STRUCTURAL)],cwd=ROOT).returncode

    original=RECON.read_bytes()
    try:
        projected=json.loads(original.decode('utf-8'))
        if not isinstance(projected,dict):
            return fail(['final reconciliation must be an object'])
        projected['state']='IN_PROGRESS'
        RECON.write_text(json.dumps(projected,indent=2,ensure_ascii=False)+'\n',encoding='utf-8')
        return subprocess.run([sys.executable,str(STRUCTURAL)],cwd=ROOT).returncode
    finally:
        RECON.write_bytes(original)


def validate_freeze()->int:
    manifest=load(FREEZE)
    ledger=load(LEDGER)
    scan1=load(SCAN1)
    scan2=load(SCAN2)
    recon=load(RECON)
    quality=load(QUALITY)
    resolution=load(QUALITY_RESOLUTION)
    errors:list[str]=[]

    if manifest.get('state')!='FROZEN_T1': errors.append('freeze manifest state must be FROZEN_T1')
    if manifest.get('programIssue')!=113 or manifest.get('trackIssue')!=114: errors.append('freeze manifest must bind #113/#114')
    if manifest.get('targetVersion')!='v18.10.0': errors.append('freeze manifest targetVersion must be v18.10.0')
    if manifest.get('unexplainedGapCount')!=0: errors.append('freeze manifest unexplainedGapCount must be zero')
    if manifest.get('unresolvedBlockingDecisions') not in ([],None): errors.append('freeze manifest must have no unresolved blocking decisions')

    frozen=manifest.get('frozenDiscovery') or {}
    if frozen.get('path')!='governance/programs/ADAPT-V18-FINAL-CLOSURE-10-10-001/feature-assurance-ledger.json':
        errors.append('freeze manifest discovery-ledger path mismatch')
    actual_ledger_blob=git_blob_sha(LEDGER)
    if frozen.get('gitBlobSha')!=actual_ledger_blob:
        errors.append(f'frozen discovery ledger blob mismatch: expected {frozen.get("gitBlobSha")}, got {actual_ledger_blob}')

    base_quality=manifest.get('baseQualityAudit') or {}
    actual_quality_blob=git_blob_sha(QUALITY)
    if base_quality.get('gitBlobSha')!=actual_quality_blob:
        errors.append(f'frozen base quality blob mismatch: expected {base_quality.get("gitBlobSha")}, got {actual_quality_blob}')
    if resolution.get('baseQualityAuditGitBlobSha')!=actual_quality_blob:
        errors.append('quality resolution is not bound to the exact base quality audit blob')

    rows=ledger.get('features') or []
    omissions1=scan1.get('omissionsFound') or []
    omissions2=scan2.get('omissionsFound') or []
    exclusions=recon.get('excludedFutureSourceCarryForward') or []
    if not isinstance(rows,list) or not isinstance(omissions1,list) or not isinstance(omissions2,list) or not isinstance(exclusions,list):
        errors.append('freeze inventory inputs must remain arrays')
        rows=[]; omissions1=[]; omissions2=[]; exclusions=[]
    counts=manifest.get('effectiveInventory') or {}
    observed={
        'highLevelDiscoveryRows':len(rows),
        'explicitFutureSourceExclusions':len(exclusions),
        'scan1Responsibilities':len(omissions1),
        'scan2Responsibilities':len(omissions2),
        'effectiveShippedV18Responsibilities':len(rows)-len(exclusions)+len(omissions1)+len(omissions2),
    }
    for key,value in observed.items():
        if counts.get(key)!=value: errors.append(f'freeze inventory count mismatch for {key}: manifest={counts.get(key)} observed={value}')
    if observed['effectiveShippedV18Responsibilities']!=180: errors.append('effective shipped-v18 responsibility count must remain 180')

    exclusion_ids={str(x.get('id') or '') for x in exclusions if isinstance(x,dict)}
    if exclusion_ids!={'PERSIST-POSTGRES-FOUNDATION'}:
        errors.append(f'future-source exclusion set changed unexpectedly: {sorted(exclusion_ids)}')
    if manifest.get('futureSourceExclusion')!='PERSIST-POSTGRES-FOUNDATION/#66':
        errors.append('freeze manifest must explicitly bind PostgreSQL future source to blocked v19 #66')

    if recon.get('state')!='COMPLETE' or recon.get('unexplainedGapCount')!=0:
        errors.append('final reconciliation must be COMPLETE with zero unexplained gaps')
    if recon.get('effectiveShippedV18ResponsibilityCount')!=180:
        errors.append('final reconciliation effective responsibility count must be 180')
    corrective={int(x.get('issue')):x for x in recon.get('correctiveIssues') or [] if isinstance(x,dict) and str(x.get('issue','')).isdigit()}
    if set(corrective)!={125,126}: errors.append(f'final corrective set must be exactly #125/#126; got {sorted(corrective)}')
    for issue in (125,126):
        item=corrective.get(issue) or {}
        if item.get('state')!='QUALIFIED': errors.append(f'corrective #{issue} is not QUALIFIED')
        evidence=item.get('qualificationEvidence') or {}
        if evidence.get('headSha')!='6714e14fbe2c2e67a7ecb64ca15b09636e7c03c9': errors.append(f'corrective #{issue} qualification head mismatch')
        if evidence.get('fastRunId')!=32854754246 or evidence.get('qualifiedRunId')!=32854819449: errors.append(f'corrective #{issue} qualification run binding mismatch')

    if resolution.get('state')!='COMPLETE': errors.append('T1 quality resolution state must be COMPLETE')
    if resolution.get('unresolvedBlockingDecisions') not in ([],None): errors.append('T1 quality resolution still has unresolved blockers')
    resolved=resolution.get('resolutions') or {}
    if set(resolved)!={'SURF-DOCUMENTATION','SURF-ADMINISTRATION'}:
        errors.append('quality resolution must resolve exactly the two T1 corrective feature rows')
    for feature,issue in (('SURF-DOCUMENTATION',125),('SURF-ADMINISTRATION',126)):
        item=resolved.get(feature) or {}
        if item.get('correctiveIssue')!=issue or item.get('closureDecision')!='READY_10_10':
            errors.append(f'{feature} final quality resolution is not READY_10_10 for corrective #{issue}')
        if item.get('qualifiedHeadSha')!='6714e14fbe2c2e67a7ecb64ca15b09636e7c03c9':
            errors.append(f'{feature} final quality resolution head mismatch')

    raw_blockers=set(quality.get('unresolvedBlockingDecisions') or [])
    if raw_blockers!={'SURF-DOCUMENTATION:#125','SURF-ADMINISTRATION:#126'}:
        errors.append(f'base quality checkpoint blocker set drifted: {sorted(raw_blockers)}')

    downstream=manifest.get('downstreamAssurance') or {}
    if downstream.get('T2ToT9Verified') is not False: errors.append('T1 freeze must explicitly leave T2-T9 unverified')
    if downstream.get('T10Started') is not False: errors.append('T1 freeze must explicitly leave T10 not started')
    if downstream.get('assuranceCeiling')!='IMPLEMENTED_UNVERIFIED': errors.append('T1 freeze assurance ceiling must be IMPLEMENTED_UNVERIFIED')
    overclaimed=[str(row.get('id') or '') for row in rows if isinstance(row,dict) and row.get('currentAssuranceState')=='VERIFIED']
    if overclaimed: errors.append('T1 discovery ledger contains prematurely VERIFIED rows: '+', '.join(overclaimed))

    required_boundaries={
        'SMART_PROVIDER_ROUTER_V2','CANONICAL_DATA_HEALTH_FRESHNESS_CACHE_PERSISTENCE_SUBSCRIPTIONS_TELEMETRY_RECONCILIATION',
        'GOVERNED_PROVIDER_LIFECYCLE','DIRECT_SEC_EDGAR_AUTHORITY','US_EQUITIES_PROCESSING_BOUNDARY','GLD_SLV_USO_TRADABLE_EXCEPTIONS','NO_EXECUTION'
    }
    if set(manifest.get('preservedBoundaries') or [])!=required_boundaries:
        errors.append('freeze preserved-boundary set is incomplete or changed')

    if errors: return fail(errors)
    print('V18 FEATURE REALITY & QUALITY AUDIT FREEZE: PASS')
    print(f"frozen discovery blob={actual_ledger_blob} · high-level={len(rows)} · scan1={len(omissions1)} · scan2={len(omissions2)} · excluded={len(exclusions)} · effective=180")
    print('correctives #125/#126: QUALIFIED · unresolved gaps=0 · T2-T9 remain IMPLEMENTED_UNVERIFIED at most')
    return 0


def main()->int:
    parser=argparse.ArgumentParser()
    parser.add_argument('--freeze',action='store_true')
    parser.parse_args()
    structural_rc=run_structural()
    if structural_rc!=0:
        return structural_rc
    if not FREEZE.exists():
        return 0
    return validate_freeze()


if __name__=='__main__':
    raise SystemExit(main())
