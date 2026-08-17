#!/usr/bin/env python3
"""Blocking approved-scope traceability completeness gate.

The immutable v15.1.2 48-item manifest records the original owning source
markers. Later behavior-preserving decomposition/UX naming may move ownership.
Current releases may provide an explicit, reviewed current-owner override
registry without mutating the historical manifest.

Behavioral suites are deliberately not nested here. The checkpointed G12 runner
owns full Go, renderer, content, version, Extreme, deterministic, HTTP and
responsive proof exactly once on the frozen source fingerprint.
"""
from pathlib import Path
import json, sys

ROOT=Path(__file__).resolve().parent
manifest=json.loads((ROOT/'renderer/qa/v15.1.2-approved-scope.json').read_text())
override_path=ROOT/'renderer/qa/v18.5-approved-scope-current-owners.json'
overrides={}
if override_path.exists():
    registry=json.loads(override_path.read_text())
    if registry.get('schemaVersion') != 1 or registry.get('release') != 'v18.5.0':
        print('Approved Scope Traceability & Implementation Completeness: FAIL')
        print('current-owner override registry schema/release mismatch')
        sys.exit(1)
    overrides=registry.get('overrides') or {}

errors=[]
if manifest.get('count')!=48 or len(manifest.get('items',[]))!=48:
    errors.append('scope manifest must contain exactly 48 approved items')
ids=[x.get('id') for x in manifest.get('items',[])]
if ids!=list(range(1,49)): errors.append(f'item ids are not exactly 1..48: {ids}')

# Overrides must only rebind an existing item and must carry a human-readable
# reason plus non-empty concrete source markers. They cannot add/remove scope.
valid_ids={str(i) for i in range(1,49)}
for key,row in overrides.items():
    if key not in valid_ids:
        errors.append(f'current-owner override refers to unknown scope item {key}')
        continue
    if not str(row.get('reason','')).strip():
        errors.append(f'#{key} current-owner override missing reason')
    ev=row.get('evidence') or []
    if not ev:
        errors.append(f'#{key} current-owner override has no evidence')
    for marker in ev:
        if not str(marker.get('file','')).strip() or not str(marker.get('contains','')).strip():
            errors.append(f'#{key} current-owner override contains blank file/marker')

# Production source fallback retains the historic decomposition behavior.
prod='\n'.join(p.read_text(errors='ignore') for p in ROOT.glob('*.go') if not p.name.endswith('_test.go'))+'\n'+(ROOT/'renderer/renderer.js').read_text(errors='ignore')+'\n'+(ROOT/'renderer/index.html').read_text(errors='ignore')

for item in manifest.get('items',[]):
    override=overrides.get(str(item.get('id')))
    evidence=(override or item).get('evidence') or []
    if not evidence: errors.append(f"#{item.get('id')} has no evidence")
    for ev in evidence:
        path=ROOT/ev['file']
        marker=ev['contains']
        # Prefer the exact declared owner. Historical manifests may then fall
        # back to the production corpus when decomposition moved a function.
        if path.exists() and marker in path.read_text(errors='ignore'):
            continue
        if override:
            errors.append(f"#{item['id']} current-owner evidence missing: {ev['file']} :: {marker}")
        elif marker not in prod:
            errors.append(f"#{item['id']} evidence marker missing from production source set: {marker}")

if errors:
    print('Approved Scope Traceability & Implementation Completeness: FAIL')
    print('\n'.join(errors[:80])); sys.exit(1)
print(f'Approved Scope Traceability & Implementation Completeness: 48/48 PASS · current-owner overrides={len(overrides)} · historical manifest preserved · behavioral proof owned by checkpointed G12 checks')
