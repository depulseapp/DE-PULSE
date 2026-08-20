#!/usr/bin/env python3
from __future__ import annotations

import json
from pathlib import Path
import subprocess
import sys

ROOT = Path(__file__).resolve().parents[2]
MANIFEST = ROOT / 'adaptive-governance' / 'V18.8.1-ZERO-MISS-RECONCILIATION.json'
EXPECTED_PACKETS = [
    'ADAPT-CI-001','ADAPT-CI-002','ADAPT-CI-003','ADAPT-CI-004',
    'ADAPT-DATA-001','ADAPT-DATA-002','ADAPT-ARCH-001','ADAPT-UI-001',
    'ADAPT-QA-001','ADAPT-GOV-001','ADAPT-COST-001','ADAPT-RECON-001',
    'ADAPT-UX-RESEARCH-001','ADAPT-SYMBOL-001','ADAPT-READINESS-001',
    'ADAPT-FRESHNESS-001','ADAPT-RESEARCH-002',
]
EXPECTED_RELEASES = ['v18.5.2','v18.6.0','v18.6.1','v18.7.0','v18.8.0']


def fail(errors: list[str]) -> int:
    print('DE.PULSE v18.8.1 zero-miss reconciliation: FAIL', file=sys.stderr)
    for error in errors:
        print(f' - {error}', file=sys.stderr)
    return 1


def main() -> int:
    errors: list[str] = []
    if not MANIFEST.is_file():
        return fail(['zero-miss reconciliation manifest missing'])
    data = json.loads(MANIFEST.read_text(encoding='utf-8'))
    if data.get('schema') != 'DE.PULSE-V18.8.1-ZERO-MISS-RECONCILIATION-1':
        errors.append('manifest schema drift')
    if data.get('releaseLine') != 'v18.8.1-development':
        errors.append('release line drift')

    historical = data.get('historicalAuthority', {})
    if historical.get('expectedRows') != 296:
        errors.append('historical reconciliation expected-row identity drift')
    for key in ('ledger','identity','gate'):
        path = historical.get(key)
        if not path or not (ROOT / path).exists():
            errors.append(f'historical authority missing: {key}={path}')
    if historical.get('disposition') != 'DELEGATED_IMMUTABLE':
        errors.append('historical authority must remain delegated/immutable')

    releases = data.get('postLedgerReleaseAuthorities', [])
    release_ids = [row.get('release') for row in releases]
    if release_ids != EXPECTED_RELEASES:
        errors.append(f'post-ledger release authority coverage drift: {release_ids!r}')
    for row in releases:
        evidence = row.get('evidence')
        if not evidence or not (ROOT / evidence).exists():
            errors.append(f"post-ledger release evidence missing: {row.get('release')} -> {evidence}")
        if not row.get('disposition'):
            errors.append(f"post-ledger release disposition missing: {row.get('release')}")

    rows = data.get('currentPackets', [])
    ids = [row.get('id') for row in rows]
    if ids != EXPECTED_PACKETS:
        errors.append('v18.8.1 packet ordering/coverage drift; every frozen packet must appear exactly once')
    if len(ids) != len(set(ids)):
        errors.append('duplicate v18.8.1 packet identity')

    pending: list[str] = []
    for row in rows:
        packet = row.get('id') or '<unknown>'
        if not row.get('owner'):
            errors.append(f'{packet}: current owner missing')
        elif row.get('state') == 'CLOSED' and not (ROOT / row['owner']).exists():
            errors.append(f"{packet}: closed owner path missing: {row['owner']}")
        evidence = row.get('evidence') or []
        if not evidence:
            errors.append(f'{packet}: evidence binding missing')
        if row.get('state') == 'CLOSED':
            for path in evidence:
                if not (ROOT / path).exists():
                    errors.append(f'{packet}: closed evidence path missing: {path}')
        else:
            pending.append(packet)
        if not row.get('disposition'):
            errors.append(f'{packet}: disposition missing')

    if data.get('requireAllClosed'):
        if pending:
            errors.append('qualification closure requested with pending packets: ' + ', '.join(pending))
    elif not pending:
        errors.append('requireAllClosed must be true once every packet is CLOSED')

    old_gate = ROOT / str(historical.get('gate', ''))
    if old_gate.is_file():
        result = subprocess.run([sys.executable, str(old_gate)], cwd=ROOT, check=False)
        if result.returncode != 0:
            errors.append('immutable historical reconciliation authority no longer passes')

    if errors:
        return fail(errors)
    print('DE.PULSE v18.8.1 zero-miss reconciliation: PASS')
    print('immutable historical 296-row authority delegated, not copied: PASS')
    print('post-ledger certified release authority chain v18.5.2 -> v18.8.0: PASS')
    print(f'v18.8.1 packet identities covered exactly once: {len(rows)}/{len(EXPECTED_PACKETS)}')
    print('all packets closed: ' + ('YES' if data.get('requireAllClosed') else 'NO - coherent development still in progress'))
    return 0


if __name__ == '__main__':
    raise SystemExit(main())
