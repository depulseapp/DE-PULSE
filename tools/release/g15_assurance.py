#!/usr/bin/env python3
from __future__ import annotations

import argparse
import datetime as dt
import hashlib
import json
from pathlib import Path


def sha256(path: Path) -> str:
    h = hashlib.sha256()
    with path.open('rb') as f:
        for chunk in iter(lambda: f.read(1024 * 1024), b''):
            h.update(chunk)
    return h.hexdigest()


def load_evidence(path: Path, source_sha: str, source_fp: str, build_id: str) -> dict:
    data = json.loads(path.read_text(encoding='utf-8-sig'))
    assert data['status'] == 'PASS', data
    assert data['certifiedSourceSha'] == source_sha, data
    assert data['sourceFingerprint'] == source_fp, data
    assert data['buildId'] == build_id, data
    assert all(value == 'PASS' for value in data['checks'].values()), data['checks']
    artifact = path.parent / data['artifact']
    assert artifact.is_file(), artifact
    actual = sha256(artifact)
    assert actual == data['artifactSha256'], (actual, data['artifactSha256'])
    return data


def main() -> int:
    p = argparse.ArgumentParser()
    p.add_argument('--mac-evidence', required=True)
    p.add_argument('--windows-evidence', required=True)
    p.add_argument('--source-sha', required=True)
    p.add_argument('--source-fingerprint', required=True)
    p.add_argument('--build-id', required=True)
    p.add_argument('--version', required=True)
    p.add_argument('--out-dir', required=True)
    args = p.parse_args()

    mac = load_evidence(Path(args.mac_evidence), args.source_sha, args.source_fingerprint, args.build_id)
    win = load_evidence(Path(args.windows_evidence), args.source_sha, args.source_fingerprint, args.build_id)
    assert mac['platform'] == 'macOS Apple Silicon'
    assert win['platform'] == 'Windows x64'

    out = Path(args.out_dir)
    out.mkdir(parents=True, exist_ok=True)
    manifest = {
        'schema': 'DE.PULSE-G15-ASSURANCE-2',
        'release': f'v{args.version}',
        'status': 'PASS',
        'promotionAuthorized': True,
        'certifiedSourceSha': args.source_sha,
        'sourceFingerprint': args.source_fingerprint,
        'buildId': args.build_id,
        'nativePlatforms': [
            {'platform': mac['platform'], 'artifact': mac['artifact'], 'sha256': mac['artifactSha256']},
            {'platform': win['platform'], 'artifact': win['artifact'], 'sha256': win['artifactSha256']},
        ],
        'noExecutionBoundary': 'PRESERVED',
        'generatedAt': dt.datetime.now(dt.timezone.utc).isoformat(),
    }
    (out / 'G15-Release-Assurance.json').write_text(json.dumps(manifest, indent=2, sort_keys=True) + '\n')
    lines = [f'# DE.PULSE v{args.version} G15 NATIVE ARTIFACT SHA-256']
    for item in manifest['nativePlatforms']:
        lines.append(f"{item['sha256']}  {item['artifact']}")
    (out / f'De-Pulse-v{args.version}-G15-SHA256.txt').write_text('\n'.join(lines) + '\n')
    print(f'PASS: v{args.version} G15 release assurance; both required native lanes PASS.')
    return 0


if __name__ == '__main__':
    raise SystemExit(main())
