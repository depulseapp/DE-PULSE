#!/usr/bin/env python3
"""Block v17 desktop packaging if any supported packaged desktop target can silently lose SQLite."""
from pathlib import Path
import subprocess, sys, os, tempfile
R=Path(__file__).resolve().parent
errors=[]
fallback=(R/'persistence_backend_fallback.go').read_text(errors='ignore') if (R/'persistence_backend_fallback.go').exists() else ''
windows=(R/'persistence_backend_windows.go').read_text(errors='ignore') if (R/'persistence_backend_windows.go').exists() else ''
native=(R/'persistence_backend_sqlite.go').read_text(errors='ignore') if (R/'persistence_backend_sqlite.go').exists() else ''
if '//go:build !cgo && !windows' not in fallback:
    errors.append('fallback build tag must exclude Windows so Windows cannot silently select JSON persistence')
if 'winsqlite3.dll' not in windows or 'func newPersistenceBackend' not in windows or 'return "sqlite"' not in windows:
    errors.append('Windows x64 SQLite backend via winsqlite3.dll missing/incomplete')
for symbol in ['sqlite3_open_v2','sqlite3_prepare_v2','sqlite3_step','sqlite3_column_blob','sqlite3_close_v2']:
    if symbol not in windows:
        errors.append(f'Windows SQLite API contract missing {symbol}')
if '//go:build cgo && !windows' not in native or 'sqlite3_open_v2' not in native:
    errors.append('native macOS/Linux SQLite backend missing/incomplete')
if 'file-fallback' not in fallback:
    errors.append('explicit unsupported-host fallback disappeared; source-test fallback should remain attributable')
# Prove the packaged Windows x64 target compiles with CGO disabled. This is the
# exact cross-compile mode used by the existing release workflow.
if not errors:
    out=Path(tempfile.gettempdir())/'depulse-v17-windows-persistence-test.exe'
    env=os.environ.copy(); env.update({'GOOS':'windows','GOARCH':'amd64','CGO_ENABLED':'0'})
    p=subprocess.run(['go','test','-c','-o',str(out)],cwd=R,env=env,text=True,capture_output=True)
    if p.returncode != 0:
        errors.append('Windows x64 CGO-disabled compile failed: '+(p.stdout+p.stderr)[-1200:].replace('\n',' | '))
    try: out.unlink()
    except FileNotFoundError: pass
if errors:
    print('v17.0 Cross-Platform Persistence Packaging Gate: BLOCK')
    for e in errors: print(' -',e)
    sys.exit(2)
print('v17.0 Cross-Platform Persistence Packaging Gate: PASS · macOS/Linux native SQLite + Windows x64 system SQLite · no Windows JSON fallback')
