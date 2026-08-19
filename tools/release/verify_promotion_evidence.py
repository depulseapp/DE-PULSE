#!/usr/bin/env python3
from __future__ import annotations

import argparse
import hashlib
import json
from pathlib import Path


def sha256(path: Path) -> str:
    h = hashlib.sha256()
    with path.open('rb') as f:
        for chunk in iter(lambda: f.read(1024 * 1024), b''):
            h.update(chunk)
    return h.hexdigest()


def load(path: Path) -> dict:
    data = json.loads(path.read_text(encoding='utf-8'))
    if not isinstance(data, dict):
        raise AssertionError(f'{path}: expected JSON object')
    return data


def main() -> int:
    p = argparse.ArgumentParser(description='Verify exact DE.PULSE certified artifacts before no-rebuild Stable promotion.')
    p.add_argument('--g15', required=True)
    p.add_argument('--mac-evidence', required=True)
    p.add_argument('--windows-evidence', required=True)
    p.add_argument('--mac-artifact', required=True)
    p.add_argument('--windows-artifact', required=True)
    p.add_argument('--version', required=True)
    p.add_argument('--build-id', required=True)
    p.add_argument('--source-fingerprint', required=True)
    p.add_argument('--certified-run-head', required=True)
    p.add_argument('--out', required=True)
    a = p.parse_args()

    g15_path = Path(a.g15)
    mac_ev_path = Path(a.mac_evidence)
    win_ev_path = Path(a.windows_evidence)
    mac_art_path = Path(a.mac_artifact)
    win_art_path = Path(a.windows_artifact)

    g15 = load(g15_path)
    mac = load(mac_ev_path)
    win = load(win_ev_path)

    release = a.version if a.version.startswith('v') else f'v{a.version}'
    expected = {
        'release': release,
        'buildId': a.build_id,
        'sourceFingerprint': a.source_fingerprint,
    }

    assert g15.get('schema') == 'DE.PULSE-G15-ASSURANCE-2', g15
    assert g15.get('status') == 'PASS', g15
    assert g15.get('promotionAuthorized') is True, g15
    assert g15.get('noExecutionBoundary') == 'PRESERVED', g15
    assert g15.get('certifiedSourceSha') == a.certified_run_head, g15
    for k, v in expected.items():
        assert g15.get(k) == v, (k, g15.get(k), v)

    native_by_platform = {x.get('platform'): x for x in g15.get('nativePlatforms', []) if isinstance(x, dict)}
    assert set(native_by_platform) == {'macOS Apple Silicon', 'Windows x64'}, native_by_platform

    checks = [
        (mac, 'macOS Apple Silicon', mac_art_path),
        (win, 'Windows x64', win_art_path),
    ]
    verified = []
    for evidence, platform, artifact_path in checks:
        assert evidence.get('schema') == 'DE.PULSE-G13-G14-NATIVE-2', evidence
        assert evidence.get('status') == 'PASS', evidence
        assert evidence.get('platform') == platform, evidence
        assert evidence.get('certifiedSourceSha') == a.certified_run_head, evidence
        for k, v in expected.items():
            assert evidence.get(k) == v, (platform, k, evidence.get(k), v)
        assert artifact_path.name == evidence.get('artifact'), (artifact_path, evidence.get('artifact'))
        digest = sha256(artifact_path)
        assert digest == evidence.get('artifactSha256'), (platform, digest, evidence.get('artifactSha256'))
        g15_platform = native_by_platform[platform]
        assert g15_platform.get('artifact') == artifact_path.name, g15_platform
        assert g15_platform.get('sha256') == digest, g15_platform
        evidence_checks = evidence.get('checks', {})
        assert evidence_checks and all(v == 'PASS' for v in evidence_checks.values()), evidence_checks
        verified.append({'platform': platform, 'artifact': artifact_path.name, 'sha256': digest})

    out = {
        'schema': 'DE.PULSE-STABLE-PROMOTION-VERIFY-1',
        'status': 'PASS',
        'release': release,
        'buildId': a.build_id,
        'certifiedSourceSha': a.certified_run_head,
        'sourceFingerprint': a.source_fingerprint,
        'noExecutionBoundary': 'PRESERVED',
        'noRebuild': True,
        'artifacts': verified,
    }
    out_path = Path(a.out)
    out_path.parent.mkdir(parents=True, exist_ok=True)
    out_path.write_text(json.dumps(out, indent=2, sort_keys=True) + '\n', encoding='utf-8')
    print('PASS: exact certified macOS/Windows artifacts + G15 evidence verified for no-rebuild Stable promotion')
    return 0


if __name__ == '__main__':
    raise SystemExit(main())
