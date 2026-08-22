#!/usr/bin/env python3
from __future__ import annotations

from pathlib import Path
import sys

ROOT = Path(__file__).resolve().parents[2]


def fail(errors: list[str]) -> int:
    print('DE.PULSE CI hardening gate: FAIL', file=sys.stderr)
    for error in errors:
        print(f' - {error}', file=sys.stderr)
    return 1


def require_order(text: str, labels: list[tuple[str, str]], errors: list[str], scope: str) -> None:
    positions: list[tuple[str, int]] = []
    for label, token in labels:
        pos = text.find(token)
        if pos < 0:
            errors.append(f'{scope}: missing {label}: {token}')
        positions.append((label, pos))
    valid = [item for item in positions if item[1] >= 0]
    for (left, lpos), (right, rpos) in zip(valid, valid[1:]):
        if lpos >= rpos:
            errors.append(f'{scope}: expected {left} before {right}')


def main() -> int:
    errors: list[str] = []
    fast = (ROOT / '.github/workflows/ci-fast.yml').read_text(encoding='utf-8')
    qualified = (ROOT / '.github/workflows/ci-qualified.yml').read_text(encoding='utf-8')
    release = (ROOT / '.github/workflows/release.yml').read_text(encoding='utf-8')
    stable_gate = (ROOT / 'tools/ci/stable_evidence_gate.py').read_text(encoding='utf-8')
    planner = (ROOT / 'tools/ci/impact_plan.py').read_text(encoding='utf-8')

    require_order(fast, [
        ('Python setup', 'actions/setup-python@'),
        ('impact plan', 'name: Deterministic impact plan'),
        ('workflow/coherence policy', 'name: Canonical workflow policy'),
        ('release identity', 'name: Release identity'),
        ('Python syntax', 'name: Python syntax'),
        ('Go setup', 'actions/setup-go@'),
        ('Node setup', 'actions/setup-node@'),
        ('Go tests', 'name: Go full unit/package suite'),
    ], errors, 'CI Fast cheap-first')

    for token in (
        'default: adaptive',
        '- adaptive',
        'requested="${INPUT_LANE:-adaptive}"',
        '--requested-lane "$requested"',
        '--target-ref "$target"',
        'base_sha:',
        'target_ref:',
        'selected_jobs: ${{ steps.resolve.outputs.selected_jobs }}',
    ):
        if token not in qualified:
            errors.append(f'CI Qualified safe Planner v3 manual default missing: {token}')
    if 'default: full' in qualified:
        errors.append('CI Qualified manual/workflow-call default must not be full')

    for token in (
        'def resolve_base(base: str, head: str, target_ref: str)',
        'trustworthy CI base unavailable',
        'HEAD^ fallback is prohibited',
        'apply_lane_override',
        'if requested == "adaptive"',
        'unknownPaths',
        'failClosedFull',
    ):
        if token not in planner:
            errors.append(f'Planner v3 trustworthy-base/fail-closed contract missing: {token}')
    if 'HEAD^' in planner and 'HEAD^ fallback is prohibited' not in planner:
        errors.append('Planner v3 must not retain an executable HEAD^ fallback')

    g11_start = release.find('\n  g11:\n')
    g12_start = release.find('\n  g12:\n')
    policy = release.find('python3 tools/ci/workflow_policy.py', g11_start if g11_start >= 0 else 0)
    if min(g11_start, g12_start, policy) < 0 or not (g11_start < policy < g12_start):
        errors.append('Release G11 must execute canonical workflow/coherence policy before G12')
    if 'release_state_coherence.py' not in stable_gate:
        errors.append('Stable evidence gate must invoke Release State Coherence')
    if 'release_state_coherence_self_test.py' not in stable_gate:
        errors.append('Stable evidence gate must invoke Release State Coherence self-test')

    for token in (
        'git ls-remote --refs origin "refs/tags/$tag"',
        'existing_ref" != "$CANDIDATE_SHA"',
    ):
        if token not in release:
            errors.append(f'publication collision guard missing: {token}')

    if errors:
        return fail(errors)

    print('DE.PULSE CI hardening gate: PASS')
    print('Fast cheap-first setup ordering: PASS')
    print('manual/reusable Qualified adaptive Planner v3 default: PASS')
    print('trustworthy full-delta base / no HEAD^ fallback: PASS')
    print('unknown-path fail-closed contract: PASS')
    print('G11 coherence preflight before G12: PASS')
    print('publication collision defense-in-depth: PASS')
    return 0


if __name__ == '__main__':
    raise SystemExit(main())
