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
    write_json(root, 'release_identity.json', {
        'version': '1.2.3', 'display_version': 'DE.PULSE v1.2.3', 'build_id': build_id, 'channel': 'STABLE',
        'stable_baseline': 'v1.2.2', 'previous_stable': 'v1.2.2',
    })
    write(root, 'VERSION.txt', f'DE.PULSE v1.2.3\nBuild: {build_id}\nPrevious Stable: v1.2.2\n')
    write(root, 'app_bootstrap.go', f'const appVersion = "1.2.3"\nconst buildID = "{build_id}"\nconst releaseChannel = "STABLE"\n')
    write(root, 'renderer/renderer.js', f"const EXPECTED_RELEASE_VERSION='1.2.3';\nconst EXPECTED_BUILD_ID='{build_id}';\n")
    index = '<title>DE.PULSE v1.2.3</title>\n' + '\n'.join(f'{a}?v=1.2.3' for a in COUPLED)
    write(root, 'renderer/index.html', index)
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

        build = json.loads((root / '.depulse-certification/resume/build-checkpoint.json').read_text())
        build['release'] = 'v1.2.2'
        write_json(root, '.depulse-certification/resume/build-checkpoint.json', build)
        manifest = json.loads((root / 'release/v1.2.3/stable-evidence-manifest.json').read_text())
        manifest['sourceFingerprint'] = 'd' * 64
        write_json(root, 'release/v1.2.3/stable-evidence-manifest.json', manifest)
        write(root, 'handoff/CURRENT.md', 'stale handoff\n')
        index = (root / 'renderer/index.html').read_text().replace('renderer.js?v=1.2.3', 'renderer.js?v=0.0.0')
        write(root, 'renderer/index.html', index)
        tags[tag] = 'e' * 40
        bad = mod.collect_errors(root, tag_lookup=lookup)
        codes = '\n'.join(bad.errors)
        for required in ('CHECKPOINT_RELEASE_MISMATCH', 'MANIFEST_FINGERPRINT_MISMATCH', 'HANDOFF_STABLE_TAG_MISMATCH', 'CACHE_BUST_MISMATCH', 'IMMUTABLE_STABLE_TAG_CONFLICT'):
            assert required in codes, (required, bad.errors)
        assert len(bad.errors) >= 5

        nextroot = root / 'g11-next'
        next_candidate_stable, _, next_stable_tag = fixture(nextroot)
        next_tags = {next_stable_tag: next_candidate_stable}
        next_lookup = lambda _root, name: next_tags.get(name, '')
        next_version = '1.2.4'
        next_build = 'v1.2.4-stable-test'
        write_json(nextroot, 'release_identity.json', {
            'version': next_version, 'display_version': 'DE.PULSE v1.2.4', 'build_id': next_build, 'channel': 'STABLE',
            'stable_baseline': 'v1.2.3', 'previous_stable': 'v1.2.3',
        })
        write(nextroot, 'VERSION.txt', f'DE.PULSE v1.2.4\nBuild: {next_build}\nPrevious Stable: v1.2.3\n')
        write(nextroot, 'app_bootstrap.go', f'const appVersion = "1.2.4"\nconst buildID = "{next_build}"\nconst releaseChannel = "STABLE"\n')
        write(nextroot, 'renderer/renderer.js', f"const EXPECTED_RELEASE_VERSION='1.2.4';\nconst EXPECTED_BUILD_ID='{next_build}';\n")
        next_index = '<title>DE.PULSE v1.2.4</title>\n' + '\n'.join(f'{a}?v=1.2.4' for a in COUPLED)
        write(nextroot, 'renderer/index.html', next_index)
        write_json(nextroot, 'release/v1.2.4/certification-manifest.json', {
            'schema': 'DE.PULSE-G12-EVIDENCE-MANIFEST-1', 'productVersion': next_version,
            'workSliceId': 'NEXT-SLICE', 'evidenceSchemaVersion': 1,
            'releaseContract': 'release/v1.2.4/release_contract.json',
        })
        g11_candidate = 'c' * 40
        create = mod.collect_errors(nextroot, g11_candidate_sha=g11_candidate, tag_lookup=next_lookup)
        assert not create.errors, create.errors
        assert create.target_tag_action == 'CREATE', create.target_tag_action
        next_tags['v1.2.4-stable'] = g11_candidate
        reuse = mod.collect_errors(nextroot, g11_candidate_sha=g11_candidate, tag_lookup=next_lookup)
        assert not reuse.errors, reuse.errors
        assert reuse.target_tag_action == 'REUSE', reuse.target_tag_action
        next_tags['v1.2.4-stable'] = 'd' * 40
        conflict = mod.collect_errors(nextroot, g11_candidate_sha=g11_candidate, tag_lookup=next_lookup)
        assert any('TARGET_TAG_CONFLICT' in error for error in conflict.errors), conflict.errors

        missing_manifest = nextroot / 'release/v1.2.4/certification-manifest.json'
        missing_manifest.unlink()
        missing = mod.collect_errors(nextroot, g11_candidate_sha=g11_candidate, tag_lookup=lambda _r, n: next_candidate_stable if n == next_stable_tag else '')
        assert any('G12_MANIFEST_MISSING' in error for error in missing.errors), missing.errors

    print('DE.PULSE Release State Coherence self-test: PASS')
    print('active release-coupled owner set: PASS')
    print('inactive legacy header exclusion: PASS')
    print('coherent Stable fixture: PASS')
    print('multi-mismatch aggregation: PASS')
    print('immutable Stable tag conflict detection: PASS')
    print('G11 target tag absent/exact/conflict preflight: PASS')
    print('canonical version-neutral G12 manifest/executor requirement: PASS')
    return 0


if __name__ == '__main__':
    raise SystemExit(main())
