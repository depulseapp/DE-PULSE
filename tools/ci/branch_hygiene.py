#!/usr/bin/env python3
from __future__ import annotations

import argparse
import json
from pathlib import Path
import subprocess


def run(*args: str, check: bool = True) -> subprocess.CompletedProcess[str]:
    return subprocess.run(args, check=check, text=True, capture_output=True)


def main() -> int:
    p = argparse.ArgumentParser(description='Audit/delete only branches fully merged into main')
    p.add_argument('--apply', action='store_true')
    p.add_argument('--json-out')
    p.add_argument('--remote', default='origin')
    p.add_argument('--base', default='main')
    args = p.parse_args()

    run('git', 'fetch', '--prune', args.remote, '+refs/heads/*:refs/remotes/%s/*' % args.remote)
    base_ref = f'refs/remotes/{args.remote}/{args.base}'
    base_sha = run('git', 'rev-parse', base_ref).stdout.strip()
    raw = run(
        'git', 'for-each-ref', f'refs/remotes/{args.remote}',
        '--format=%(refname:short)\t%(objectname)'
    ).stdout

    protected = {args.base, 'HEAD'}
    merged: list[dict[str, str]] = []
    retained: list[dict[str, str]] = []
    for line in raw.splitlines():
        if not line.strip():
            continue
        ref, sha = line.split('\t', 1)
        prefix = f'{args.remote}/'
        name = ref[len(prefix):] if ref.startswith(prefix) else ref
        if name in protected or name.startswith('HEAD ->'):
            continue
        result = run('git', 'merge-base', '--is-ancestor', sha, base_sha, check=False)
        item = {'branch': name, 'sha': sha}
        if result.returncode == 0:
            merged.append(item)
        else:
            retained.append(item)

    deleted: list[dict[str, str]] = []
    if args.apply:
        for item in merged:
            name = item['branch']
            result = run('git', 'push', args.remote, '--delete', name, check=False)
            if result.returncode != 0:
                raise SystemExit(f"failed deleting merged branch {name}: {result.stderr.strip()}")
            deleted.append(item)

    report = {
        'schema': 'DE.PULSE-BRANCH-HYGIENE-1',
        'base': args.base,
        'baseSha': base_sha,
        'mode': 'APPLY' if args.apply else 'DRY_RUN',
        'mergedCandidates': merged,
        'deleted': deleted,
        'retainedUniqueOrDiverged': retained,
        'policy': 'Only branch tips already contained in main are deletion candidates. Unique/diverged tips are never deleted automatically.',
    }
    text = json.dumps(report, indent=2, sort_keys=True) + '\n'
    print(text, end='')
    if args.json_out:
        path = Path(args.json_out)
        path.parent.mkdir(parents=True, exist_ok=True)
        path.write_text(text)
    return 0


if __name__ == '__main__':
    raise SystemExit(main())
