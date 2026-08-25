#!/usr/bin/env python3
from __future__ import annotations

import argparse
import hashlib
import json
import os
import re
import subprocess
import sys
from dataclasses import dataclass
from pathlib import Path
from typing import Callable

ROOT = Path(__file__).resolve().parents[2] if Path(__file__).resolve().parts[-3:-1] == ('tools','ci') else Path.cwd()
RELEASE_COUPLED_ASSETS = (
    'renderer.js', 'documentation-ui.js', 'live-dom-reconcile.js',
    'watchlist-ui.js', 'watchlist-desk.css', 'market-header-ui.js',
    'ui-layout-contracts.css', 'surface-consolidation.js',
    'surface-consolidation.css', 'documentation-access.js',
)

@dataclass
class Result:
    errors: list[str]
    stable_release: str = ''
    stable_tag: str = ''
    stable_candidate: str = ''
    current_version: str = ''
    target_tag: str = ''
    target_tag_action: str = 'NOT_EVALUATED'


def read_json(path: Path, errors: list[str], code: str) -> dict:
    try:
        value = json.loads(path.read_text(encoding='utf-8'))
        if not isinstance(value, dict):
            raise ValueError('root must be an object')
        return value
    except Exception as exc:
        errors.append(f'{code}: {path.as_posix()}: {exc}')
        return {}


def read_text(path: Path, errors: list[str], code: str) -> str:
    try:
        return path.read_text(encoding='utf-8')
    except Exception as exc:
        errors.append(f'{code}: {path.as_posix()}: {exc}')
        return ''


def parse_semver(value: str) -> tuple[int, int, int] | None:
    match = re.fullmatch(r'(?:v)?(\d+)\.(\d+)\.(\d+)', value.strip())
    return tuple(int(x) for x in match.groups()) if match else None


def normalize_release(value: str) -> str:
    value = value.strip()
    return value if value.startswith('v') else f'v{value}'


def git_tag_lookup(root: Path, tag: str) -> str:
    proc = subprocess.run(['git', 'rev-parse', '-q', '--verify', f'refs/tags/{tag}^{{commit}}'], cwd=root, text=True, capture_output=True, check=False)
    return proc.stdout.strip() if proc.returncode == 0 else ''


def add_mismatch(errors: list[str], code: str, actual, expected) -> None:
    if actual != expected:
        errors.append(f'{code}: actual={actual!r} expected={expected!r}')


def git_blob_token(path: Path) -> str:
    data = path.read_bytes()
    return hashlib.sha1(f'blob {len(data)}\0'.encode() + data).hexdigest()[:16]


def cache_identity_ok(root: Path, index: str, name: str) -> bool:
    path = root / 'renderer' / name
    return path.is_file() and f'{name}?v={git_blob_token(path)}' in index


def current_stable(state: dict, errors: list[str]) -> tuple[str, str, str, str, str]:
    stable = state.get('stable', {}) if isinstance(state.get('stable'), dict) else {}
    version = str(stable.get('productVersion', '')).strip().lstrip('v')
    release = f'v{version}' if version else ''
    tag = str(stable.get('tag', '')).strip()
    candidate = str(stable.get('candidateSha', '')).strip()
    fingerprint = str(stable.get('sourceFingerprint', '')).strip()
    build_id = str(stable.get('buildId', '')).strip()
    if not parse_semver(version): errors.append(f'CURRENT_STATE_STABLE_VERSION_INVALID: {version!r}')
    if not re.fullmatch(r'[0-9a-f]{40}', candidate): errors.append(f'CURRENT_STATE_STABLE_CANDIDATE_INVALID: {candidate!r}')
    if not re.fullmatch(r'[0-9a-f]{64}', fingerprint): errors.append(f'CURRENT_STATE_STABLE_FINGERPRINT_INVALID: {fingerprint!r}')
    if version and not tag: errors.append('CURRENT_STATE_STABLE_TAG_MISSING')
    if version and not build_id: errors.append('CURRENT_STATE_STABLE_BUILD_ID_MISSING')
    return release, tag, candidate, fingerprint, build_id


def active_release_closure(root: Path, state: dict, identity: dict, errors: list[str]) -> tuple[bool, dict]:
    capability = state.get('productCapabilityGate', {}) if isinstance(state.get('productCapabilityGate'), dict) else {}
    status = str(capability.get('reservationStatus', '')).strip().upper()
    work_slice_rel = str(capability.get('workSlicePath', '')).strip()
    work_slice = read_json(root / work_slice_rel, errors, 'RELEASE_CLOSURE_WORK_SLICE_UNREADABLE') if work_slice_rel else {}
    if not work_slice or work_slice.get('type') != 'PRODUCT_RELEASE_CLOSURE' or status not in {'ACTIVE', 'IN_PROGRESS'}:
        return False, work_slice

    current_version = str(identity.get('version', '')).strip().lstrip('v')
    previous_stable = str(identity.get('previous_stable', '')).strip().lstrip('v')
    add_mismatch(errors, 'RELEASE_CLOSURE_RESERVED_ID_MISMATCH', str(capability.get('reservedWorkSliceId', '')).strip(), str(work_slice.get('workSliceId', '')).strip())
    add_mismatch(errors, 'RELEASE_CLOSURE_RESERVED_ISSUE_MISMATCH', capability.get('reservedIssue'), work_slice.get('issue'))
    add_mismatch(errors, 'RELEASE_CLOSURE_RESERVED_BRANCH_MISMATCH', str(capability.get('reservedBranch', '')).strip(), str(work_slice.get('branch', '')).strip())
    add_mismatch(errors, 'RELEASE_CLOSURE_PUBLIC_VERSION_MISMATCH', str(work_slice.get('publicProductVersion', '')).strip().lstrip('v'), current_version)
    add_mismatch(errors, 'RELEASE_CLOSURE_BASELINE_VERSION_MISMATCH', str(work_slice.get('stableProductVersionAtStart', '')).strip().lstrip('v'), previous_stable)
    add_mismatch(errors, 'RELEASE_CLOSURE_TARGET_STABLE_MISMATCH', str(work_slice.get('targetStable', '')).strip(), f'v{current_version}-stable')
    if work_slice.get('productBehaviorChange') is not True:
        errors.append('RELEASE_CLOSURE_PRODUCT_BEHAVIOR_CHANGE_INVALID: must be true')
    if work_slice.get('blocksNextProductCapability') is not True:
        errors.append('RELEASE_CLOSURE_NEXT_CAPABILITY_BLOCK_INVALID: must remain true')
    current_semver = parse_semver(current_version)
    previous_semver = parse_semver(previous_stable)
    if not current_semver or not previous_semver or current_semver <= previous_semver:
        errors.append(f'RELEASE_CLOSURE_VERSION_NOT_MONOTONIC: current={current_version!r} previous={previous_stable!r}')
    return True, work_slice


def collect_errors(root: Path, *, g11_candidate_sha: str = '', tag_lookup: Callable[[Path, str], str] | None = None) -> Result:
    errors: list[str] = []
    tag_lookup = tag_lookup or git_tag_lookup
    build = read_json(root / '.depulse-certification/resume/build-checkpoint.json', errors, 'BUILD_CHECKPOINT_UNREADABLE')
    evidence = read_json(root / '.depulse-certification/resume/release-evidence-checkpoint.json', errors, 'RELEASE_EVIDENCE_UNREADABLE')
    identity = read_json(root / 'release_identity.json', errors, 'RELEASE_IDENTITY_UNREADABLE')
    state = read_json(root / 'governance/current-state.json', errors, 'CURRENT_STATE_UNREADABLE')

    checkpoint_release = normalize_release(str(evidence.get('release', build.get('release', '')))) if (evidence.get('release') or build.get('release')) else ''
    checkpoint_stable = evidence.get('stable', {}) if isinstance(evidence.get('stable'), dict) else {}
    checkpoint_tag = str(checkpoint_stable.get('tag', '')).strip()
    checkpoint_candidate = str(checkpoint_stable.get('promotionCommit', '')).strip()
    checkpoint_fp = str(checkpoint_stable.get('sourceFingerprint', '')).strip()
    checkpoint_build = str(checkpoint_stable.get('buildId', '')).strip()
    if checkpoint_release:
        if build:
            add_mismatch(errors, 'CHECKPOINT_RELEASE_MISMATCH', str(build.get('release', '')), checkpoint_release)
            add_mismatch(errors, 'CHECKPOINT_CHANNEL_MISMATCH', str(build.get('channel', '')), 'STABLE')
            add_mismatch(errors, 'CHECKPOINT_BUILD_ID_MISMATCH', str(build.get('buildId', '')), checkpoint_build)
            add_mismatch(errors, 'CHECKPOINT_FINGERPRINT_MISMATCH', str(build.get('sourceFingerprint', '')), checkpoint_fp)
            certified = build.get('certifiedStable', {}) if isinstance(build.get('certifiedStable'), dict) else {}
            add_mismatch(errors, 'CHECKPOINT_TAG_MISMATCH', str(certified.get('tag', '')), checkpoint_tag)
            add_mismatch(errors, 'CHECKPOINT_CANDIDATE_MISMATCH', str(certified.get('promotionCommit', '')), checkpoint_candidate)
        manifest = read_json(root / 'release' / checkpoint_release / 'stable-evidence-manifest.json', errors, 'STABLE_MANIFEST_UNREADABLE')
        if manifest:
            add_mismatch(errors, 'MANIFEST_RELEASE_MISMATCH', str(manifest.get('release', '')), checkpoint_release)
            add_mismatch(errors, 'MANIFEST_STATUS_MISMATCH', str(manifest.get('status', '')), 'STABLE_PUBLISHED')
            add_mismatch(errors, 'MANIFEST_TAG_MISMATCH', str(manifest.get('stableTag', '')), checkpoint_tag)
            add_mismatch(errors, 'MANIFEST_CANDIDATE_MISMATCH', str(manifest.get('certifiedCandidate', '')), checkpoint_candidate)
            add_mismatch(errors, 'MANIFEST_FINGERPRINT_MISMATCH', str(manifest.get('sourceFingerprint', '')), checkpoint_fp)
            add_mismatch(errors, 'MANIFEST_BUILD_ID_MISMATCH', str(manifest.get('buildId', '')), checkpoint_build)

    stable_release, stable_tag, stable_candidate, stable_fingerprint, stable_build_id = current_stable(state, errors)
    current_version = str(identity.get('version', '')).strip()
    current_build = str(identity.get('build_id', '')).strip()
    current_channel = str(identity.get('channel', '')).strip()
    display_version = str(identity.get('display_version', '')).strip()
    previous_stable = str(identity.get('previous_stable', '')).strip().lstrip('v')
    release_closure, release_work_slice = active_release_closure(root, state, identity, errors)

    if release_closure:
        expected_published_release = f'v{previous_stable}' if previous_stable else ''
        add_mismatch(errors, 'RELEASE_CLOSURE_STABLE_VERSION_MISMATCH', stable_release, expected_published_release)
        add_mismatch(errors, 'RELEASE_CLOSURE_CHECKPOINT_VERSION_MISMATCH', checkpoint_release, expected_published_release)
        add_mismatch(errors, 'RELEASE_CLOSURE_STABLE_TAG_MISMATCH', stable_tag, checkpoint_tag)
        add_mismatch(errors, 'RELEASE_CLOSURE_STABLE_CANDIDATE_MISMATCH', stable_candidate, checkpoint_candidate)
        add_mismatch(errors, 'RELEASE_CLOSURE_STABLE_FINGERPRINT_MISMATCH', stable_fingerprint, checkpoint_fp)
        add_mismatch(errors, 'RELEASE_CLOSURE_STABLE_BUILD_MISMATCH', stable_build_id, checkpoint_build)
        add_mismatch(errors, 'RELEASE_CLOSURE_BASELINE_CANDIDATE_MISMATCH', str(release_work_slice.get('baselineCandidateSha', '')).strip(), stable_candidate)
        add_mismatch(errors, 'RELEASE_CLOSURE_BASELINE_FINGERPRINT_MISMATCH', str(release_work_slice.get('baselineSourceFingerprint', '')).strip(), stable_fingerprint)
        add_mismatch(errors, 'RELEASE_CLOSURE_BASELINE_BUILD_MISMATCH', str(release_work_slice.get('baselineBuildId', '')).strip(), stable_build_id)
        stable_platform = str((state.get('stable') or {}).get('platformBuildNumber', '')).strip()
        try:
            if int(current_build and identity.get('bundle_version', 0)) <= int(stable_platform):
                errors.append(f'RELEASE_CLOSURE_PLATFORM_BUILD_NOT_MONOTONIC: candidate={identity.get("bundle_version")!r} stable={stable_platform!r}')
        except Exception:
            errors.append('RELEASE_CLOSURE_PLATFORM_BUILD_INVALID')
    else:
        add_mismatch(errors, 'CURRENT_STATE_IDENTITY_VERSION_MISMATCH', current_version, stable_release.removeprefix('v'))
        add_mismatch(errors, 'CURRENT_STATE_IDENTITY_BUILD_MISMATCH', current_build, stable_build_id)
    if current_channel != 'STABLE': errors.append(f'IDENTITY_CHANNEL_INVALID: {current_channel!r}')

    handoff = read_text(root / 'handoff/CURRENT.md', errors, 'HANDOFF_UNREADABLE')
    for code, token in (('HANDOFF_STABLE_TAG_MISMATCH', stable_tag), ('HANDOFF_CANDIDATE_MISMATCH', stable_candidate), ('HANDOFF_FINGERPRINT_MISMATCH', stable_fingerprint), ('HANDOFF_BUILD_ID_MISMATCH', stable_build_id)):
        if handoff and token and token not in handoff: errors.append(f'{code}: {token!r} not present in handoff/CURRENT.md')

    version_text = read_text(root / 'VERSION.txt', errors, 'VERSION_UNREADABLE')
    for code, token in (('VERSION_DISPLAY_MISMATCH', display_version), ('VERSION_BUILD_MISMATCH', current_build), ('VERSION_PREDECESSOR_MISMATCH', f'Previous Stable: v{previous_stable}')):
        if version_text and token and token not in version_text: errors.append(f'{code}: {token!r} not present')

    boot = read_text(root / 'app_bootstrap.go', errors, 'APP_BOOTSTRAP_UNREADABLE')
    if boot:
        if f'const appVersion = "{current_version}"' not in boot: errors.append('APP_VERSION_MISMATCH: app_bootstrap.go does not match release_identity.version')
        if f'const buildID = "{current_build}"' not in boot: errors.append('APP_BUILD_ID_MISMATCH: app_bootstrap.go does not match release_identity.build_id')
        if f'const releaseChannel = "{current_channel}"' not in boot: errors.append('APP_CHANNEL_MISMATCH: app_bootstrap.go does not match release_identity.channel')

    index = read_text(root / 'renderer/index.html', errors, 'RENDERER_INDEX_UNREADABLE')
    if index and current_version:
        if f'<title>DE.PULSE v{current_version}</title>' not in index: errors.append('RENDERER_TITLE_MISMATCH: index.html title does not match release identity')
        for asset in RELEASE_COUPLED_ASSETS:
            if not cache_identity_ok(root, index, asset): errors.append(f'CACHE_IDENTITY_MISMATCH: {asset} is not bound to its Git-content identity')

    contract: dict = {}
    contract_path = root / 'release' / f'v{current_version}' / 'release_contract.json'
    if current_version and contract_path.exists(): contract = read_json(contract_path, errors, 'RELEASE_CONTRACT_UNREADABLE')
    historical_identity_asset = str(contract.get('identity_asset', '')).strip() if contract else ''
    runtime_identity_asset = str(identity.get('renderer_identity_asset', '')).strip()
    identity_asset = runtime_identity_asset or historical_identity_asset
    if runtime_identity_asset and re.search(r'(?:^|[-_])v\d+(?:[._-]\d+)+', runtime_identity_asset, re.I):
        errors.append(f'CURRENT_RENDERER_IDENTITY_ASSET_VERSIONED: {runtime_identity_asset!r}')
    if identity_asset:
        if '/' in identity_asset or '\\' in identity_asset or identity_asset.startswith('.'):
            errors.append(f'IDENTITY_ASSET_INVALID: {identity_asset!r}')
        else:
            overlay_path = root / 'renderer' / identity_asset
            overlay = read_text(overlay_path, errors, 'IDENTITY_OVERLAY_UNREADABLE')
            if overlay and f"DEPULSE_RELEASE_VERSION = '{current_version}'" not in overlay: errors.append('RENDERER_VERSION_MISMATCH: identity overlay version mismatch')
            if overlay and f"DEPULSE_RELEASE_BUILD_ID = '{current_build}'" not in overlay: errors.append('RENDERER_BUILD_ID_MISMATCH: identity overlay build mismatch')
            if index and not cache_identity_ok(root, index, identity_asset): errors.append('IDENTITY_OVERLAY_CACHE_IDENTITY_MISMATCH: identity overlay cache key is not content-derived')
            if runtime_identity_asset and historical_identity_asset and runtime_identity_asset != historical_identity_asset:
                if '/' in historical_identity_asset or '\\' in historical_identity_asset or historical_identity_asset.startswith('.'):
                    errors.append(f'HISTORICAL_IDENTITY_ASSET_INVALID: {historical_identity_asset!r}')
                else:
                    archived = root / 'release' / 'history' / f'v{current_version}' / 'renderer' / historical_identity_asset
                    if not archived.is_file():
                        errors.append(f'HISTORICAL_IDENTITY_ARCHIVE_MISSING: {archived.relative_to(root).as_posix()}')
                    elif overlay_path.is_file() and archived.read_bytes() != overlay_path.read_bytes():
                        errors.append('HISTORICAL_IDENTITY_ARCHIVE_MISMATCH: archived certified overlay differs from current neutral runtime owner')
    else:
        renderer = read_text(root / 'renderer' / 'renderer.js', errors, 'RENDERER_UNREADABLE')
        if renderer and f"const EXPECTED_RELEASE_VERSION='{current_version}';" not in renderer: errors.append('RENDERER_VERSION_MISMATCH: renderer.js release version mismatch')
        if renderer and f"const EXPECTED_BUILD_ID='{current_build}';" not in renderer: errors.append('RENDERER_BUILD_ID_MISMATCH: renderer.js build id mismatch')

    existing = tag_lookup(root, stable_tag) if stable_tag else ''
    if stable_tag and not existing: errors.append(f'IMMUTABLE_STABLE_TAG_MISSING: {stable_tag}')
    elif stable_tag and existing != stable_candidate: errors.append(f'IMMUTABLE_STABLE_TAG_CONFLICT: {stable_tag} -> {existing}, expected {stable_candidate}')

    g12_path = root / 'release' / f'v{current_version}' / 'certification-manifest.json'
    if current_version and not g12_path.is_file(): errors.append(f'G12_MANIFEST_MISSING: {g12_path.relative_to(root).as_posix()}')
    elif current_version:
        g12 = read_json(g12_path, errors, 'G12_MANIFEST_UNREADABLE')
        if g12:
            add_mismatch(errors, 'G12_MANIFEST_SCHEMA_MISMATCH', str(g12.get('schema', '')), 'DE.PULSE-G12-EVIDENCE-MANIFEST-1')
            add_mismatch(errors, 'G12_MANIFEST_VERSION_MISMATCH', str(g12.get('productVersion', '')), current_version)
            if not str(g12.get('workSliceId', '')).strip(): errors.append('G12_MANIFEST_WORK_SLICE_MISSING: workSliceId is required')
            if not isinstance(g12.get('evidenceSchemaVersion'), int): errors.append('G12_MANIFEST_EVIDENCE_SCHEMA_INVALID: evidenceSchemaVersion must be an integer')
    if not (root / 'tools/release/run_full_certification.py').is_file(): errors.append('G12_CANONICAL_EXECUTOR_MISSING: tools/release/run_full_certification.py')

    target_tag = stable_tag if current_version == stable_release.removeprefix('v') else (f'v{current_version}-stable' if release_closure else f'v{current_version}')
    target_action = 'NOT_EVALUATED'
    if target_tag:
        existing_target = tag_lookup(root, target_tag)
        if g11_candidate_sha:
            if not re.fullmatch(r'[0-9a-f]{40}', g11_candidate_sha): errors.append(f'G11_CANDIDATE_SHA_INVALID: {g11_candidate_sha!r}')
            elif not existing_target: target_action = 'CREATE'
            elif existing_target == g11_candidate_sha: target_action = 'REUSE'
            else:
                target_action = 'CONFLICT'
                errors.append(f'TARGET_TAG_CONFLICT: {target_tag} -> {existing_target}, candidate={g11_candidate_sha}')
        else:
            target_action = 'REUSE' if existing_target == stable_candidate else ('CONFLICT' if existing_target else 'MISSING')
    return Result(errors, stable_release, stable_tag, stable_candidate, current_version, target_tag, target_action)


def main() -> int:
    parser = argparse.ArgumentParser()
    parser.add_argument('--root', type=Path, default=ROOT)
    parser.add_argument('--g11-candidate-sha', default='')
    parser.add_argument('--json-out', type=Path)
    args = parser.parse_args()
    candidate = args.g11_candidate_sha.strip()
    if not candidate and os.environ.get('GITHUB_WORKFLOW') == 'DE.PULSE | Release G11-G16':
        proc = subprocess.run(['git', 'rev-parse', 'HEAD'], cwd=args.root.resolve(), text=True, capture_output=True, check=False)
        if proc.returncode == 0: candidate = proc.stdout.strip()
    result = collect_errors(args.root.resolve(), g11_candidate_sha=candidate)
    payload = {'schema':'DE.PULSE-RELEASE-STATE-COHERENCE-3','status':'PASS' if not result.errors else 'FAIL','stableRelease':result.stable_release,'stableTag':result.stable_tag,'stableCandidate':result.stable_candidate,'currentVersion':result.current_version,'targetTag':result.target_tag,'targetTagAction':result.target_tag_action,'errors':result.errors}
    if args.json_out:
        args.json_out.parent.mkdir(parents=True, exist_ok=True)
        args.json_out.write_text(json.dumps(payload, indent=2) + '\n', encoding='utf-8')
    if result.errors:
        print(f'DE.PULSE Release State Coherence: FAIL ({len(result.errors)} issue(s))', file=sys.stderr)
        for error in result.errors: print(f' - {error}', file=sys.stderr)
        return 1
    print('DE.PULSE Release State Coherence: PASS')
    print(f'current Stable machine state: {result.stable_tag} -> {result.stable_candidate}')
    print('predecessor resume checkpoint integrity: PRESERVED')
    print('release-coupled renderer cache identity: CONTENT_DERIVED')
    print(f'target tag action: {result.target_tag_action}')
    return 0

if __name__ == '__main__':
    raise SystemExit(main())
