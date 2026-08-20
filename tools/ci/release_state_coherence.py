#!/usr/bin/env python3
from __future__ import annotations

import argparse
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
    'renderer.js',
    'documentation-ui.js',
    'live-dom-reconcile.js',
    'watchlist-v18.5.1.js',
    'watchlist-v18.5.1.css',
    'market-header-ui.js',
    'ui-v18.5.1.css',
    'surface-consolidation-v18.6.js',
    'surface-consolidation-v18.6.css',
    'documentation-access-v18.6.js',
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
    m = re.fullmatch(r'(?:v)?(\d+)\.(\d+)\.(\d+)', value.strip())
    if not m:
        return None
    return tuple(int(x) for x in m.groups())  # type: ignore[return-value]


def normalize_release(value: str) -> str:
    value = value.strip()
    return value if value.startswith('v') else f'v{value}'


def git_tag_lookup(root: Path, tag: str) -> str:
    proc = subprocess.run(
        ['git', 'rev-parse', '-q', '--verify', f'refs/tags/{tag}^{{commit}}'],
        cwd=root,
        text=True,
        capture_output=True,
        check=False,
    )
    return proc.stdout.strip() if proc.returncode == 0 else ''


def add_mismatch(errors: list[str], code: str, actual, expected) -> None:
    if actual != expected:
        errors.append(f'{code}: actual={actual!r} expected={expected!r}')


def collect_errors(
    root: Path,
    *,
    g11_candidate_sha: str = '',
    tag_lookup: Callable[[Path, str], str] | None = None,
) -> Result:
    errors: list[str] = []
    tag_lookup = tag_lookup or git_tag_lookup

    build = read_json(root / '.depulse-certification/resume/build-checkpoint.json', errors, 'BUILD_CHECKPOINT_UNREADABLE')
    evidence = read_json(root / '.depulse-certification/resume/release-evidence-checkpoint.json', errors, 'RELEASE_EVIDENCE_UNREADABLE')
    identity = read_json(root / 'release_identity.json', errors, 'RELEASE_IDENTITY_UNREADABLE')

    stable_release = normalize_release(str(evidence.get('release', build.get('release', '')))) if (evidence.get('release') or build.get('release')) else ''
    stable = evidence.get('stable', {}) if isinstance(evidence.get('stable'), dict) else {}
    stable_tag = str(stable.get('tag', '')).strip()
    stable_candidate = str(stable.get('promotionCommit', '')).strip()
    stable_fingerprint = str(stable.get('sourceFingerprint', '')).strip()
    stable_build_id = str(stable.get('buildId', '')).strip()

    if not stable_release:
        errors.append('STABLE_RELEASE_MISSING: release not found in durable checkpoints')
    if not stable_tag:
        errors.append('STABLE_TAG_MISSING: stable.tag missing from release evidence checkpoint')
    if not stable_candidate:
        errors.append('STABLE_CANDIDATE_MISSING: stable.promotionCommit missing from release evidence checkpoint')
    if not re.fullmatch(r'[0-9a-f]{64}', stable_fingerprint):
        errors.append(f'STABLE_FINGERPRINT_INVALID: {stable_fingerprint!r}')

    if build:
        add_mismatch(errors, 'CHECKPOINT_RELEASE_MISMATCH', str(build.get('release', '')), stable_release)
        add_mismatch(errors, 'CHECKPOINT_CHANNEL_MISMATCH', str(build.get('channel', '')), 'STABLE')
        add_mismatch(errors, 'CHECKPOINT_BUILD_ID_MISMATCH', str(build.get('buildId', '')), stable_build_id)
        add_mismatch(errors, 'CHECKPOINT_FINGERPRINT_MISMATCH', str(build.get('sourceFingerprint', '')), stable_fingerprint)
        certified = build.get('certifiedStable', {}) if isinstance(build.get('certifiedStable'), dict) else {}
        add_mismatch(errors, 'CHECKPOINT_TAG_MISMATCH', str(certified.get('tag', '')), stable_tag)
        add_mismatch(errors, 'CHECKPOINT_CANDIDATE_MISMATCH', str(certified.get('promotionCommit', '')), stable_candidate)
        add_mismatch(errors, 'CHECKPOINT_CERTIFIED_FINGERPRINT_MISMATCH', str(certified.get('sourceFingerprint', '')), stable_fingerprint)

    if stable_release:
        manifest_path = root / 'release' / stable_release / 'stable-evidence-manifest.json'
        manifest = read_json(manifest_path, errors, 'STABLE_MANIFEST_UNREADABLE')
        if manifest:
            add_mismatch(errors, 'MANIFEST_RELEASE_MISMATCH', str(manifest.get('release', '')), stable_release)
            add_mismatch(errors, 'MANIFEST_STATUS_MISMATCH', str(manifest.get('status', '')), 'STABLE_PUBLISHED')
            add_mismatch(errors, 'MANIFEST_TAG_MISMATCH', str(manifest.get('stableTag', '')), stable_tag)
            add_mismatch(errors, 'MANIFEST_CANDIDATE_MISMATCH', str(manifest.get('certifiedCandidate', '')), stable_candidate)
            add_mismatch(errors, 'MANIFEST_FINGERPRINT_MISMATCH', str(manifest.get('sourceFingerprint', '')), stable_fingerprint)
            add_mismatch(errors, 'MANIFEST_BUILD_ID_MISMATCH', str(manifest.get('buildId', '')), stable_build_id)

    handoff = read_text(root / 'handoff/CURRENT.md', errors, 'HANDOFF_UNREADABLE')
    if handoff:
        expected_tokens = {
            'HANDOFF_STABLE_TAG_MISMATCH': stable_tag,
            'HANDOFF_CANDIDATE_MISMATCH': stable_candidate,
            'HANDOFF_FINGERPRINT_MISMATCH': stable_fingerprint,
            'HANDOFF_BUILD_ID_MISMATCH': stable_build_id,
        }
        for code, token in expected_tokens.items():
            if token and token not in handoff:
                errors.append(f'{code}: {token!r} not present in handoff/CURRENT.md')

    current_version = str(identity.get('version', '')).strip()
    display_version = str(identity.get('display_version', '')).strip()
    current_build = str(identity.get('build_id', '')).strip()
    current_channel = str(identity.get('channel', '')).strip()
    previous_stable = str(identity.get('previous_stable', '')).strip()
    stable_baseline = str(identity.get('stable_baseline', '')).strip()

    if not parse_semver(current_version):
        errors.append(f'IDENTITY_VERSION_INVALID: {current_version!r}')
    if current_channel != 'STABLE':
        errors.append(f'IDENTITY_CHANNEL_INVALID: {current_channel!r}')
    if current_version and not current_build.startswith(f'v{current_version}-'):
        errors.append(f'IDENTITY_BUILD_ID_VERSION_MISMATCH: version={current_version!r} build={current_build!r}')

    version_text = read_text(root / 'VERSION.txt', errors, 'VERSION_UNREADABLE')
    if version_text:
        for code, token in (
            ('VERSION_DISPLAY_MISMATCH', display_version),
            ('VERSION_BUILD_MISMATCH', current_build),
            ('VERSION_PREDECESSOR_MISMATCH', f'Previous Stable: {previous_stable}'),
        ):
            if token and token not in version_text:
                errors.append(f'{code}: {toke!r} not present')

    boot = read_text(root / 'app_bootstrap.go', errors, 'APP_BOOTSTRAP_UNREADABLE')
    if boot:
        if current_version and f'const appVersion = "{current_version}"' not in boot:
            errors.append('APP_VERSION_MISMATCH: app_bootstrap.go does not match release_identity.version')
        if current_build and f'const buildID = "{current_build}"' not in boot:
            errors.append('APP_BUILD_ID_MISMATCH: app_bootstrap.go does not match release_identity.build_id')
        if current_channel and f'const releaseChannel = "{current_channel}"' not in boot:
            errors.append('APP_CHANNEL_MISMATCH: app_bootstrap.go does not match release_identity.channel')

    index = read_text(root / 'renderer/index.html', errors, 'RENDERER_INDEX_UNREADABLE')
    if index and current_version:
        if f'<title>DE.PULSE v{current_version}</title>' not in index:
            errors.append('RENDERER_TITLE_MISMATCH: index.html title does not match release identity')
        for asset in RELEASE_COUPLED_ASSETS:
            if f'{asset}?v={current_version}' not in index:
                errors.append(f'CACHE_BUST_MISMATCH: {asset} is not cache-busted with v{current_version}')

    contract = {}
    if current_version:
        contract_path = root / 'release' / f'v{current_version}' / 'release_contract.json'
        if contract_path.exists():
            contract = read_json(contract_path, errors, 'RELEASE_CONTRACT_UNREADABLE')
    identity_asset = str(contract.get('identity_asset', '')).strip() if contract else ''
    if identity_asset:
        if '/' in identity_asset or '\\' in identity_asset or identity_asset.startswith('.'):
            errors.append(f'IDENTITY_ASSET_INVALID: {identity_asset!r}')
        else:
            overlay = read_text(root / 'renderer' / identity_asset, errors, 'IDENTITY_OVERLAY_UNREADABLE')
            if overlay:
                if f"DEPULSE_RELEASE_VERSION = '{current_version}'" not in overlay:
                    errors.append('RENDERER_VERSION_MISMATCH: identity overlay version mismatch')
                if f"DEPULSE_RELEASE_BUILD_ID = '{current_build}'" not in overlay:
                    errors.append('RENDERER_BUILD_ID_MISMATCH: identity overlay build mismatch')
            if index and f'{identity_asset}?v={current_version}' not in index:
                errors.append('IDENTITY_OVERLAY_CACHE_BUST_MISMATCH: identity overlay cache key mismatch')
    else:
        renderer = read_text(root / 'renderer/renderer.js', errors, 'RENDERER_UNREADABLE')
        if renderer:
            if f"const EXPECTED_RELEASE_VERSION='{current_version}';" not in renderer:
                errors.append('RENDERER_VERSION_MISMATCH: renderer.js release version mismatch')
            if f"const EXPECTED_BUILD_ID='{current_build}';" not in renderer:
                errors.append('RENDERER_BUILD_ID_MISMATCH: renderer.js build id mismatch')

    stable_semver = parse_semver(stable_release)
    current_semver = parse_semver(current_version)
    if stable_semver and current_semver:
        if current_semver < stable_semver:
            errors.append(f'IDENTITY_BEHIND_STABLE: current={current_version} stable={stable_release}')
        elif current_semver == stable_semver:
            add_mismatch(errors, 'CURRENT_STABLE_BUILD_MISMATCH', current_build, stable_build_id)
        else:
            if previous_stable != stable_release:
                errors.append(f'PREVIOUS_STABLE_MISMATCH: candidate={previous_stable!r} certified={stable_release!r}')
            if stable_baseline != stable_release:
                errors.append(f'STABLE_BASELINE_MISMATCH: candidate={stable_baseline!r} certified={stable_release!r}')

    if stable_tag:
        existing_stable = tag_lookup(root, stable_tag)
        if not existing_stable:
            errors.append(f'IMMUTABLE_STABLE_TAG_MISSING: {stable_tag}')
        elif stable_candidate and existing_stable != stable_candidate:
            errors.append(f'IMMUTABLE_STABLE_TAG_CONFLICT: {stable_tag} -> {existing_stable}, expected {stable_candidate}')

    target_tag = f'v{current_version}-stable' if current_version else ''
    target_action = 'NOT_EVALUATED'
    if target_tag:
        existing_target = tag_lookup(root, target_tag)
        if g11_candidate_sha:
            if not re.fullmatch(r'[0-9a-f]{40}', g11_candidate_sha):
                errors.append(f'G11_CANDIDATE_SHA_INVALID: {g11_candidate_sha!r}')
            elif not existing_target:
                target_action = 'CREATE'
            elif existing_target == g11_candidate_sha:
                target_action = 'REUSE'
            else:
                target_action = 'CONFLICT'
                errors.append(f'TARGET_TAG_CONFLICT: {target_tag} -> {existing_target}, candidate={g11_candidate_sha}')
            scaffold = root / 'release' / f'v{current_version}' / 'run_full_certification.sh'
            if not scaffold.is_file():
                errors.append(f'RELEASE_SCAFFOLD_MISSING: {scaffold.relative_to(root).as_posix()}')
        elif current_semver and stable_semver and current_semver > stable_semver and existing_target:
            target_action = 'CONFLICT'
            errors.append(f'TARGET_TAG_ALREADY_EXISTS: next-release target {target_tag} -> {existing_target}')
        elif current_semver and stable_semver and current_semver == stable_semver:
            target_action = 'REUSE' if existing_target == stable_candidate else ('CONFLICT' if existing_target else 'MISSINGˆ[ÙN‚ˆ\™Ù]ØXİ[ÛˆH	ĞP”ÑS•	ÈYˆ›İ^\İ[™×İ\™Ù][ÙH	ÑVTÕÉÂ‚ˆ™]\›ˆ™\İ[
ˆ\œ›ÜœÏY\œ›ÜœËˆİX›WÜ™[X\ÙO\İX›WÜ™[X\ÙKˆİX›WİYÏ\İX›WİYËˆİX›WØØ[™Y]O\İX›WØØ[™Y]Kˆİ\œ™[İ™\œÚ[ÛXİ\œ™[İ™\œÚ[Û‹ˆ\™Ù]İYÏ]\™Ù]İYËˆ\™Ù]İY×ØXİ[Û]\™Ù]ØXİ[Û‹ˆ
B‚‚™YˆXZ[Š
HOˆ[‚ˆ\œÙ\ˆH\™Ü\œÙK\™İ[Y[\œÙ\Š
Bˆ\œÙ\‹˜YØ\™İ[Y[
	ËK\›Ûİ	Ë\OT]Y˜][T“ÓÕ
Bˆ\œÙ\‹˜YØ\™İ[Y[
	ËKYÌLKXØ[™Y]K\ÚIËY˜][IÉÊBˆ\œÙ\‹˜YØ\™İ[Y[
	ËKZœÛÛ‹[İ]	Ë\OT]
Bˆ\™ÜÈH\œÙ\‹œ\œÙWØ\™ÜÊ
B‚ˆÌLWØØ[™Y]WÜÚHH\™ÜË™ÌLWØØ[™Y]WÜÚKœİš\

BˆYˆ›İÌLWØØ[™Y]WÜÚH[™ÜË™[š\›Û‹™Ù]
	ÑÒUP—ÕÓÔ’Ñ“ÕÉÊHOH	ÑK”SÑH™[X\ÙHÌLKQÌM‰Î‚ˆ›ØÈHİXœ›ØÙ\ÜËœ[ŠÉÙÚ]	Ë	Ü™]‹\\œÙIË	ÒPQ	×KİÙX\™ÜËœ›Ûİœ™\ÛÛ™J
K^UYKØ\\™WÛİ]]UYKÚXÚÏQ˜[ÙJBˆYˆ›ØËœ™]\›˜ÛÙHOH‚ˆÌLWØØ[™Y]WÜÚHH›ØËœİİ]œİš\

Bˆ™\İ[HÛÛXİÙ\œ›ÜœÊ\™ÜËœ›Ûİœ™\ÛÛ™J
KÌLWØØ[™Y]WÜÚOYÌLWØØ[™Y]WÜÚJBˆ^[ØYHÂˆ	ÜØÚ[XIÎˆ	ÑK”SÑKT‘SPTÑKTÕUKPÓÒT‘SÑKLIËˆ	Üİ]\ÉÎˆ	ÔTÔÉÈYˆ›İ™\İ[™\œ›ÜœÈ[ÙH	ÑRS	Ëˆ	ÜİX›T™[X\ÙIÎˆ™\İ[œİX›WÜ™[X\ÙKˆ	ÜİX›UYÉÎˆ™\İ[œİX›WİYËˆ	ÜİX›PØ[™Y]IÎˆ™\İ[œİX›WØØ[™Y]Kˆ	Øİ\œ™[™\œÚ[Û‰Îˆ™\İ[˜İ\œ™[İ™\œÚ[Û‹ˆ	İ\™Ù]YÉÎˆ™\İ[\™Ù]İYËˆ	İ\™Ù]YĞXİ[Û‰Îˆ™\İ[\™Ù]İY×ØXİ[Û‹ˆ	Ù\œ›ÜœÉÎˆ™\İ[™\œ›ÜœËˆBˆYˆ\™ÜËšœÛÛ—Ûİ]‚ˆ\™ÜËšœÛÛ—Ûİ]Üš]Wİ^
œÛÛ‹™[\Ê^[ØY[™[LŠH
È	×‰Ë[˜ÛÙ[™ÏIİ]‹N	ÊB‚ˆYˆ™\İ[™\œ›ÜœÎ‚ˆš[
ˆ‘K”SÑH™[X\ÙHİ]HÛÚ\™[˜ÙNˆRS
Û[Š™\İ[™\œ›ÜœÊ_H\ÜİYJÊJH‹š[O\Ş\Ëœİ\œŠBˆ›Üˆ\œ›Üˆ[ˆ™\İ[™\œ›ÜœÎ‚ˆš[
‰ÈHÙ\œ›ÜŸIËš[O\Ş\Ëœİ\œŠBˆ™]\›ˆB‚ˆš[
	ÑK”SÑH™[X\ÙHİ]HÛÚ\™[˜ÙNˆTÔÉÊBˆš[
‰ØÙ\YšYYİX›NˆÜ™\İ[œİX›WİYßHOˆÜ™\İ[œİX›WØØ[™Y]_IÊBˆš[
‰Øİ\œ™[™[X\ÙHY[]NˆÜ™\İ[˜İ\œ™[İ™\œÚ[ÛŸIÊBˆš[
‰İ\™Ù]YÈİ]NˆÜ™\İ[\™Ù]İY×ØXİ[ÛŸH
Ü™\İ[\™Ù]İYßJIÊBˆ™]\›ˆ‚‚šYˆ×Û˜[YW×ÈOH	××ÛXZ[—×ÉÎ‚ˆ˜Z\ÙHŞ\İ[Q^]
XZ[Š
JB