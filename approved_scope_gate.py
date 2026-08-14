#!/usr/bin/env python3
"""Blocking G11 gate: immutable Approved Scope traceability completeness.

Behavioral suites are deliberately not nested here. The checkpointed G12 runner
owns full Go, renderer, content, version, Extreme, deterministic, HTTP and
responsive proof exactly once on the frozen source fingerprint.
"""
from pathlib import Path
import json, sys
ROOT=Path(__file__).resolve().parent
manifest=json.loads((ROOT/'renderer/qa/v15.1.2-approved-scope.json').read_text())
errors=[]
if manifest.get('count')!=48 or len(manifest.get('items',[]))!=48:
    errors.append('scope manifest must contain exactly 48 approved items')
ids=[x.get('id') for x in manifest.get('items',[])]
if ids!=list(range(1,49)): errors.append(f'item ids are not exactly 1..48: {ids}')
prod='\n'.join(p.read_text(errors='ignore') for p in ROOT.glob('*.go') if not p.name.endswith('_test.go'))+'\n'+(ROOT/'renderer/renderer.js').read_text(errors='ignore')+'\n'+(ROOT/'renderer/index.html').read_text(errors='ignore')
for item in manifest.get('items',[]):
    evidence=item.get('evidence') or []
    if not evidence: errors.append(f"#{item.get('id')} has no evidence")
    for ev in evidence:
        path=ROOT/ev['file']
        # Historical manifests recorded the then-owning filename. After the v16.1 G2 decomposition, ownership may move while behavior remains canonical.
        if path.exists() and ev['contains'] in path.read_text(errors='ignore'):
            continue
        if ev['contains'] not in prod:
            errors.append(f"#{item['id']} evidence marker missing from production source set: {ev['contains']}")
if errors:
    print('Approved Scope Traceability & Implementation Completeness: FAIL')
    print('\n'.join(errors[:80])); sys.exit(1)
print('Approved Scope Traceability & Implementation Completeness: 48/48 PASS · behavioral proof owned by checkpointed G12 checks')
