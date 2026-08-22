#!/usr/bin/env python3
"""Adaptive Build Process v2 pre-freeze qualification.

Runs cheap/medium read-only release-truth checks before an immutable RC is
created. Evidence is source-fingerprint scoped and written outside the source
tree. Permanent Build Resume Protocol metadata under `.depulse-certification`
is deliberately excluded from the product fingerprint so recording progress
does not invalidate the candidate it describes.
"""
from __future__ import annotations
import concurrent.futures, datetime as dt, json, os, signal, subprocess, sys, time
from pathlib import Path

ROOT = Path(__file__).resolve().parents[2]
if str(ROOT) not in sys.path:
    sys.path.insert(0, str(ROOT))
from source_fingerprint import canonical_source_fingerprint

FINGERPRINT_EXCLUDED_DIRS = {'.git', '__pycache__', '.depulse-certification'}
FINGERPRINT_EXCLUDED_SUFFIXES = {'.log', '.tmp', '.out', '.exe'}


def fingerprint():
    return canonical_source_fingerprint(ROOT, excluded_dirs=FINGERPRINT_EXCLUDED_DIRS, excluded_suffixes=FINGERPRINT_EXCLUDED_SUFFIXES)


def run_job(job, outdir):
    jid, cmd, timeout = job
    result_path = outdir / f'{jid}.json'
    if result_path.exists():
        try:
            prior = json.loads(result_path.read_text())
            if prior.get('status') == 'PASS' and prior.get('command') == cmd:
                prior = dict(prior); prior['reused'] = True; return prior
        except Exception: pass
    started=time.time(); log=outdir/f'{jid}.log'; rc=1; status='INFRA FAIL'
    try:
        with log.open('w',encoding='utf-8',errors='replace') as lf:
            p=subprocess.Popen(cmd,cwd=ROOT,text=True,stdout=lf,stderr=subprocess.STDOUT,start_new_session=True)
            try:
                rc=p.wait(timeout=timeout); status='PASS' if rc==0 else 'FAIL'
            except subprocess.TimeoutExpired:
                try: os.killpg(p.pid,signal.SIGTERM); p.wait(timeout=5)
                except Exception:
                    try: os.killpg(p.pid,signal.SIGKILL)
                    except Exception: pass
                lf.write('\nTIMEOUT\n'); rc=124; status='INFRA FAIL'
    except FileNotFoundError as exc:
        with log.open('a',encoding='utf-8',errors='replace') as lf: lf.write(f'\nMISSING EXECUTABLE: {exc}\n')
        rc=127; status='INFRA FAIL'
    except Exception as exc:
        with log.open('a',encoding='utf-8',errors='replace') as lf: lf.write(f'\nINFRA EXCEPTION: {type(exc).__name__}: {exc}\n')
        rc=125; status='INFRA FAIL'
    result={'id':jid,'status':status,'returncode':rc,'duration_seconds':round(time.time()-started,3),'command':cmd,'log':str(log),'reused':False}
    result_path.write_text(json.dumps(result,indent=2)+'\n'); return result


def wave(name,jobs,max_workers,outdir):
    print(f'[{name}] {len(jobs)} jobs · max concurrency {max_workers}',flush=True); results=[]
    if max_workers==1:
        for job in jobs:
            result=run_job(job,outdir); results.append(result); state='RESUME' if result.get('reused') else result['status']; print(f"  [{state}] {result['id']} · {result['duration_seconds']}s",flush=True)
        return sorted(results,key=lambda x:x['id'])
    with concurrent.futures.ThreadPoolExecutor(max_workers=max_workers) as executor:
        futures={executor.submit(run_job,job,outdir):job[0] for job in jobs}
        for future in concurrent.futures.as_completed(futures):
            result=future.result(); results.append(result); state='RESUME' if result.get('reused') else result['status']; print(f"  [{state}] {result['id']} · {result['duration_seconds']}s",flush=True)
    return sorted(results,key=lambda x:x['id'])


def write_summary(outdir,fp,rows,status,started):
    total=sum(r['duration_seconds'] for r in rows); fresh=sum(r['duration_seconds'] for r in rows if not r.get('reused')); wall=time.time()-started; ratio=round(fresh/max(wall,0.001),2) if fresh else 0.0
    summary={'schema':'DE.PULSE-PREFREEZE-QUALIFICATION-2','generated_at':dt.datetime.now(dt.timezone.utc).isoformat(),'source_fingerprint':fp,'status':status,'wall_seconds':round(wall,3),'evidence_summed_job_seconds':round(total,3),'fresh_job_seconds':round(fresh,3),'reused_count':sum(1 for r in rows if r.get('reused')),'current_run_parallel_efficiency_ratio':ratio,'resume_protocol':'adaptive-governance/BUILD_RESUME_PROTOCOL.md','checkpoint_metadata_fingerprint_excluded':True,'results':rows}
    (outdir/'summary.json').write_text(json.dumps(summary,indent=2)+'\n'); return summary


def main():
    fp=fingerprint(); outdir=ROOT.parent/f'.depulse-prefreeze-{fp[:12]}'; outdir.mkdir(exist_ok=True); py=sys.executable
    light=[
        ('v18_baseline',[py,'v18_baseline_gate.py'],30),('release_identity',[py,'release_identity.py','--verify'],30),('adaptive_resume',[py,'tools/ci/adaptive_resume_gate.py'],30),('v183_scope',[py,'v18_3_scope_gate.py'],30),('v183_principal',[py,'v18_3_principal_engineer_gate.py'],30),('v182_principal_inherited',[py,'v18_2_principal_engineer_gate.py'],30),('v181_principal_inherited',[py,'v18_1_principal_engineer_gate.py'],30),('v1805_ui',['node','v18_0_5_renderer_test.js'],60),('fingerprint_portability',[py,'tools/ci/source_fingerprint_portability_test.py'],30),('v18_typography',[py,'v18_documentation_typography_gate.py'],30),('version',[py,'version_consistency_test.py'],30),('documentation',[py,'tools/ci/documentation_governance_gate.py'],60),('data_utility',[py,'tools/ci/data_utility_gate.py'],45),('functionality_utility',[py,'tools/ci/functionality_utility_checkpoint_gate.py'],45),('data_health',[py,'tools/ci/data_health_policy_gate.py'],30),('source_health',[py,'tools/ci/source_health_architecture_gate.py'],120),('content',[py,'tools/ci/content_copy_audit_test.py'],60),('renderer',['node','renderer_logic_test.js'],90),
    ]
    medium_go=[
        ('go_vet',['go','vet','./...'],150),('v183_postgres_hosted_foundation',['go','test','-count=1','-run','^TestV183','./...'],180),('v182_admin_presence_sessions',['go','test','-count=1','-run','^TestV182','./...'],180),('v181_multi_user',['go','test','-count=1','-run','^TestV181','./...'],180),('v1806_router_shock',['go','test','-count=1','-run','^TestV1806','./...'],120),('v1801_smart_router_rapid',['go','test','-count=1','-run','^TestV1801','./...'],180),('v1805_symbols',['go','test','-count=1','-run','^TestV1805','./...'],120),('v180_security',['go','test','-count=1','-run','^TestV180','./...'],180),('full_go',['go','test','-count=1','-json','-timeout=120s','./...'],180),('cgo_fallback',['bash','-lc','CGO_ENABLED=0 go test -count=1 -json -timeout=120s ./...'],180),('windows_compile',['bash','-lc',"tmp=$(mktemp -d); trap 'rm -rf \\\"$tmp\\\"' EXIT; CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go test -c -o \\\"$tmp/depulse-v18.test.exe\\\" . && test -s \\\"$tmp/depulse-v18.test.exe\\\""],180),
    ]
    medium_other=[('http_workflow',[py,'http_workflow_test.py'],300),('deterministic',['node','tools/ci/deterministic_equivalence_test.js'],150),('professional_trader',['node','trader_acceptance_test.js'],150),('responsive',[py,'tools/ci/responsive_ui_sharded_gate.py'],320)]
    started=time.time(); rows=[]; rows+=wave('LIGHT',light,4,outdir)
    if any(r['status']!='PASS' for r in rows): write_summary(outdir,fp,rows,'FAIL',started); return 2
    rows+=wave('MEDIUM-GO',medium_go,1,outdir)
    if any(r['status']!='PASS' for r in rows): write_summary(outdir,fp,rows,'FAIL',started); return 2
    rows+=wave('MEDIUM-OTHER',medium_other,2,outdir); status='PASS' if all(r['status']=='PASS' for r in rows) else 'FAIL'; summary=write_summary(outdir,fp,rows,status,started)
    print(f"Pre-freeze qualification: {status} · wall {summary['wall_seconds']:.1f}s · evidence jobs {summary['evidence_summed_job_seconds']:.1f}s · fresh jobs {summary['fresh_job_seconds']:.1f}s · reused {summary['reused_count']} · current-run concurrency ratio {summary['current_run_parallel_efficiency_ratio']}x")
    return 0 if status=='PASS' else 2

if __name__=='__main__': raise SystemExit(main())
