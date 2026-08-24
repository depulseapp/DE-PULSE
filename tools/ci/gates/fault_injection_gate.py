#!/usr/bin/env python3
"""Professional truth fault-injection gate."""
from pathlib import Path
import subprocess,sys
ROOT=Path(__file__).resolve().parents[3]
p=subprocess.run(['go','test','-count=1','-run','TestV1603FaultInjection','./...'],cwd=ROOT,text=True,capture_output=True)
if p.returncode:
    print('Professional Truth Fault-Injection Gate: FAIL');print(p.stdout[-6000:]);print(p.stderr[-6000:]);sys.exit(1)
print('Professional Truth Fault-Injection Gate: PASS · false-Fresh, future-time, partial-backfill, and State/Why defects rejected')
