#!/usr/bin/env python3
"""DE.PULSE CI/CD execution model metadata and isolated-job helper.

The final RC heavy suite remains owned by certification_runner.py. This helper
records the dependency/resource graph and can run a single command from an
isolated source copy, writing logs outside that copy. It intentionally does not
launch all HEAVY jobs concurrently.
"""
from pathlib import Path
import argparse, json, os, shutil, subprocess, tempfile, time
from source_fingerprint import canonical_source_fingerprint
ROOT=Path(__file__).resolve().parent
PLAN=ROOT/'ci_pipeline_plan.json'
EXCLUDE={'.depulse-certification','__pycache__','.git'}
def fingerprint(root=ROOT):
 return canonical_source_fingerprint(root, excluded_dirs=EXCLUDE, excluded_suffixes={'.log','.tmp','.out','.exe'})
def main():
 ap=argparse.ArgumentParser(); ap.add_argument('--print-plan',action='store_true'); ap.add_argument('--fingerprint',action='store_true'); ap.add_argument('--run',nargs=argparse.REMAINDER); a=ap.parse_args()
 if a.print_plan: print(PLAN.read_text()); return
 if a.fingerprint: print(fingerprint()); return
 if a.run:
  fp=fingerprint(); ev=ROOT.parent/f'.depulse-ci-evidence-{fp[:12]}'; ev.mkdir(exist_ok=True)
  with tempfile.TemporaryDirectory(prefix='depulse-ci-') as td:
   work=Path(td)/'src'; shutil.copytree(ROOT,work,ignore=shutil.ignore_patterns(*EXCLUDE))
   started=time.time(); p=subprocess.run(a.run,cwd=work,text=True,capture_output=True)
   rec={'source_fingerprint':fp,'command':a.run,'returncode':p.returncode,'duration_seconds':round(time.time()-started,3)}
   (ev/'single-job.json').write_text(json.dumps(rec,indent=2)+'\n'); (ev/'single-job.log').write_text(p.stdout+p.stderr)
   print(json.dumps(rec)); raise SystemExit(p.returncode)
 ap.print_help()
if __name__=='__main__': main()
