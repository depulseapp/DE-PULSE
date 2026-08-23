#!/usr/bin/env python3
from __future__ import annotations

import importlib.util
import json
from pathlib import Path
import tempfile
import sys

spec = importlib.util.spec_from_file_location('coherence', Path(__file__).with_name('release_state_coherence.py'))
mod = importlib.util.module_from_spec(spec)
assert spec and spec.loader
sys.modules['coherence'] = mod
spec.loader.exec_module(mod)

COUPLED = mod.RELEASE_COUPLED_ASSETS
assert 'market-header-ui.js' in COUPLED, 'active Market Header owner must be release-coupled'
assert 'documentation-ui.js' in COUPLED, 'active Documentation owner must be release-coupled'
assert 'header-v18.5.1.js' not in COUPLED, 'inactive legacy header must not be release-coupled'


def write(root: Path, rel: str, text: str) -> None:
    p = root / rel
    p.parent.mkdir(parents=True, exist_ok=True)
    p.write_text(text, encoding='utf-8')


def write_json(root: Path, rel: str, value: dict) -> None:
    write(root, rel, json.dumps(value, indent=2) + '\n')


def write_renderer_assets(root: Path, version: str, build_id: str) -> None:
    for asset in COUPLED:
        body = f'fixture asset: {asset}\n'
        if asset == 'renderer.js':
            body = (
                f"const EXPECTED_RELEASE_VERSION='{version}';\n"
                f"const EXPECTED_BUILD_ID='{build_id}';\n"
            )
        write(root, f'renderer/{asset}', body)


def write_index(root: Path, version: str) -> None:
    entries = []
    for asset in COUPLED:
        token = mod.git_blob_token(root / 'renderer' / asset)
        entries.append(f'{asset}?v={token}')
    write(root, 'renderer/index.html', f'<title>DE.PULSE v{version}</title>\n' + '\n'.join(entries))


def fixture(root: Path) -> tuple[str, str, str]:
    candidate = 'a' * 40
    fp = 'b' * 64
    build_id = 'v1.2.3-stable-test'
    release = 'v1.2.3'
    tag = 'v1.2.3-stable'
    write_json(root, '.depulse-certification/resume/release-evidence-checkpoint.json', {
        'release': release,
        'channel': 'STABLE',
        'stable': {'tag': tag, 'promotionCommit': candidate, 'sourceFingerprint': fp, 'buildId': build_id, 'promotionRun': 3, 'certificationRun': 3},
        'currentCandidate': {'qualificationHead': 'c' * 40},
        'evidence': {'G10': {'fastRun': 1, 'qualifiedRun': 2}},
    })
    write_json(root, '.depulse-certification/resume/build-checkpoint.json', {
        'release': release, 'channel': 'STABLE', 'buildId': build_id, 'sourceFingerprint': fp,
        'certifiedStable': {'tag': tag, 'promotionCommit': candidate, 'sourceFingerprint': fp},
    })
    write_json(root, f'release/{release}/stable-evidence-manifest.json', {
        'schema': 'DE.PULSE-STABLE-EVIDENCE-1', 'release': release, 'status': 'STABLE_PUBLISHED',
        'buildId': build_id, 'stableTag': tag, 'certifiedCandidate': candidate, 'sourceFingerprint': fp,
    })
    write_json(root, 'governance/current-state.json', {
        'schema': 'DE.PULSE-CURRENT-STATE-1',
        'stable': {
            'productVersion': '1.2.3',
            'tag': tag,
            'candidateSha': candidate,
            'sourceFingerprint': fp,
            'buildId': build_id,
        },
    })
    write_json(root, 'release_identity.json', {
        'version': '1.2.3', 'display_version': 'DE.PULSE v1.2.3', 'build_id': build_id, 'channel': 'STABLE',
        'stable_baseline': 'v1.2.2', 'previous_stable': 'v1.2.2',
    })
    write(root, 'VERSION.txt', f'DE.PULSE v1.2.3\nBuild: {build_id}\nPrevious Stable: v1.2.2\n')
    write(root, 'app_bootstrap.go', f'const appVersion = "1.2.3"\nconst buildID = "{build_id}"\nconst releaseChannel = "STABLE"\n')
    write_renderer_assets(root, '1.2.3', build_id)
    write_index(root, '1.2.3')
    write(root, 'handoff/CURRENT.md', f'{tag}\n{candidate}\n{fp}\n{build_id}\n')
    write_json(root, 'release/v1.2.3/certification-manifest.json', {
        'schema': 'DE.PULSE-G12-EVIDENCE-MANIFEST-1', 'productVersion': '1.2.3',
        'workSliceId': 'TEST-SLICE', 'evidenceSchemaVersion': 1,
        'releaseContract': 'release/v1.2.3/release_contract.json',
    })
    write(root, 'tools/release/run_full_certification.py', '# canonical G12 executor\n')
    return candidate, fp, tag


def main() -> int:
    with tempfile.TemporaryDirectory() as td:
        root = Path(td)
        candidate, fp, tag = fixture(root)
        tags = {tag: candidate}
        lookup = lambda _root, name: tags.get(name, '')
        ok = mod.collect_errors(root, tag_lookup=lookup)
        assert not ok.errors, ok.errors
        assert ok.target_tag_action == 'REUSE', ok.target_tag_action

        build = json.loads((root / '.depulse-certification/resume/build-checkpoint.json').read_text())
        build['release'] = 'v1.2.2'
        write_json(root, '.depulse-certification/resume/build-checkpoint.json', build)
        manifest = json.loads((root / 'release/v1.2.3/stable-evidence-manifest.json').read_text())
        manifest['sourceFingerprint'] = 'd' * 64
        write_json(root, 'release/v1.2.3/stable-evidence-manifest.json', manifest)
        write(root, 'handoff/CURRENT.md', 'stale handoff\n')
        index = (root / 'renderer/index.html').read_text()
        renderer_token = mod.git_blob_token(root / 'renderer/renderer.js')
        write(root, 'renderer/index.html', index.replace(f'renderer.js?v={renderer_token}', 'renderer.js?v=0'))
        tags[tag] = 'e' * 40
        bad = mod.collect_errors(root, tag_lookup=lookup)
        codes = '\n'.join(bad.errors)
        for required in (
            'CHECKPOINT_RELEASE_MISMATCH',
            'MANIFEST_FINGERPRINT_MISMATCH',
            'HANDOFF_STABLE_TAG_MISMATCH',
            'CACHE_IDENTITY_MISMATCH',
            'IMMUTABLE_STABLE_TAG_CONFLICT',
        ):
            assert required in codes, (required, bad.errors)
        assert len(bad.errors) >= 5

        missingroot = root / 'missing-g12'
        missing_candidate, _, missing_tag = fixture(missingroot)
        (missingroot / 'release/v1.2.3/certification-manifest.json').unlink()
        missing = mod.collect_errors(
            missingroot,
            tag_lookup=lambda _r, name: missing_candidate if name == missing_tag else '',
        )
        assert any('G12_MANIFEST_MISSING' in error for error in missing.errors), missing.errors

    print('DE.PULSE Release State Coherence self-test: PASS')
    print('active release-coupled owner set: PASS')
    print('inactive legacy header exclusion: PASS')
    print('coherent current Stable fixture: PASS')
    print('predecessor resume checkpoint separation: PASS')
    print('content-derived cache identity fixture: PASS')
    print('current Stable tag exact-reuse/conflict preflight: PASS')
    print('multi-mismatch aggregation: PASS')
    print('immutable Stable tag conflict detection: PASS')
    print('canonical version-neutral G12 manifest/executor requirement: PASS')
    return 0


if __name__ == '__main__':
    raise SystemExit(main())
