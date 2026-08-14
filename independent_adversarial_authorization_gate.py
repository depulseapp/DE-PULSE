#!/usr/bin/env python3
from pathlib import Path
import subprocess,sys
ROOT=Path(__file__).resolve().parent
checks=[
 (['go','test','-count=1','-run','TestV160(3Authorization|4Authorization)','./...'],'inherited v16 authorization'),
 (['go','test','-count=1','-run','TestV1611','./...'],'v16.1.1 escaped Market Intelligence truth'),
]
for cmd,label in checks:
 p=subprocess.run(cmd,cwd=ROOT,text=True,capture_output=True)
 if p.returncode:
  print('Independent Adversarial Authorization Gate: FAIL · '+label); print((p.stdout+p.stderr)[-7000:]); sys.exit(1)
print('Independent Adversarial Authorization Gate: PASS · fresh v16.0.x plus v16.1.1 escaped-truth combinations verified')
