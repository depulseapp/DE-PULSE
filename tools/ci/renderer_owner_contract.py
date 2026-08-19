#!/usr/bin/env python3
from __future__ import annotations

from pathlib import Path
import re
import sys

ROOT=Path(__file__).resolve().parents[2]
INDEX=ROOT/'renderer'/'index.html'
MONOLITH=ROOT/'renderer'/'renderer.js'
OWNER=ROOT/'renderer'/'documentation-ui.js'
ACCESS=ROOT/'renderer'/'documentation-access-v18.6.js'


def fail(errors:list[str])->int:
    print('DE.PULSE renderer owner contract: FAIL',file=sys.stderr)
    for error in errors: print(f' - {error}',file=sys.stderr)
    return 1


def main()->int:
    errors=[]
    index=INDEX.read_text(encoding='utf-8')
    monolith=MONOLITH.read_text(encoding='utf-8')
    owner=OWNER.read_text(encoding='utf-8') if OWNER.is_file() else ''
    access=ACCESS.read_text(encoding='utf-8')

    scripts=re.findall(r'<script\s+src="([^"]+)"',index)
    names=[x.split('?',1)[0] for x in scripts]
    required=['renderer.js','documentation-ui.js','documentation-access-v18.6.js']
    for name in required:
        if names.count(name)!=1: errors.append(f'{name} must load exactly once; got {names.count(name)}')
    if all(name in names for name in required):
        if not (names.index('renderer.js')<names.index('documentation-ui.js')<names.index('documentation-access-v18.6.js')):
            errors.append('load order must be renderer.js -> documentation-ui.js -> documentation-access-v18.6.js')

    if 'documentation-ui-v' in index or 'documentation-ui-v' in owner:
        errors.append('new long-lived capability owner must not use version-stacked naming')
    required_owner=(
        "const OWNER='renderer/documentation-ui.js'",
        "state:'ACTIVE_OWNER_WITH_LEGACY_FALLBACK'",
        'renderMarkdown=documentationMarkdown',
        'hydrateDocumentation=documentationHydrate',
        'renderDocumentation=documentationRender',
        'globalThis.__DE_PULSE_DOCUMENTATION_UI__',
        'globalThis.__DE_PULSE_RENDERER_OWNERS__',
        "data-render-owner=\"documentation-ui\"",
    )
    for token in required_owner:
        if token not in owner: errors.append(f'Documentation owner contract missing: {token}')

    legacy=('function renderMarkdown(', 'async function hydrateDocumentation(', 'function renderDocumentation(')
    for token in legacy:
        if token not in monolith: errors.append(f'expected transitional monolith fallback missing unexpectedly: {token}')
    if 'legacy-architecture-diagram' not in owner:
        errors.append('transitional dependency on monolith architectureDiagram must remain explicit until extracted')
    if 'deletionGate' not in owner:
        errors.append('owner must declare the legacy deletion gate')

    if 'const v186BaseHydrateDocumentation=hydrateDocumentation;' not in access:
        errors.append('role-access decorator must wrap the active hydration owner')
    if 'const v186BaseRenderDocumentation=renderDocumentation;' not in access:
        errors.append('role-access decorator must wrap the active render owner')

    if errors: return fail(errors)
    print('DE.PULSE renderer owner contract: PASS')
    print('Documentation capability-oriented runtime owner: renderer/documentation-ui.js')
    print('load order owner before role-access decorator: PASS')
    print('legacy monolith fallback retained and explicitly gated for later physical deletion: PASS')
    print('version-stacked new-owner naming prevention: PASS')
    return 0


if __name__=='__main__':
    raise SystemExit(main())
