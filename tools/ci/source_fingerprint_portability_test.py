#!/usr/bin/env python3
from pathlib import Path, PurePosixPath, PureWindowsPath
import sys
import tempfile

ROOT = Path(__file__).resolve().parents[2]
if str(ROOT) not in sys.path:
    sys.path.insert(0, str(ROOT))
from source_fingerprint import canonical_rel_key, canonical_source_fingerprint

def need(ok, msg):
    if not ok:
        raise SystemExit('Source fingerprint portability: FAIL · '+msg)

need(canonical_rel_key(PurePosixPath('renderer/docs/user.md')) == 'renderer/docs/user.md', 'POSIX canonical key')
need(canonical_rel_key(PureWindowsPath(r'renderer\docs\user.md')) == 'renderer/docs/user.md', 'Windows canonical key')
with tempfile.TemporaryDirectory() as td:
    r=Path(td)
    (r/'A').mkdir(); (r/'renderer'/'docs').mkdir(parents=True)
    (r/'A'/'one.txt').write_bytes(b'one')
    (r/'renderer'/'docs'/'user.md').write_bytes(b'user')
    a=canonical_source_fingerprint(r)
    b=canonical_source_fingerprint(r)
    need(a==b and len(a)==64, 'deterministic canonical digest')
print('Source fingerprint portability: PASS · canonical POSIX path key is identical for POSIX/Windows path models')