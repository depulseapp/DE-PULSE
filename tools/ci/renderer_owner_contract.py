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
HEADER=ROOT/'renderer'/'market-header-ui.js'
LEGACY_HEADER=ROOT/'renderer'/'header-v18.5.1.js'
QUALIFIED=ROOT/'.github'/'workflows'/'ci-qualified.yml'
NODE_TEST=ROOT/'tests/renderer/documentation_ui_owner_test.js'
HEADER_NODE_TEST=ROOT/'tests'/'renderer'/'market_header_owner_test.js'
BROWSER_TEST=ROOT/'tools'/'ci'/'documentation_owner_browser_test.py'
HIERARCHY_TEST=ROOT/'release'/'v18.5.1'/'browser_ui_hierarchy_test.py'


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
    header=HEADER.read_text(encoding='utf-8') if HEADER.is_file() else ''
    qualified=QUALIFIED.read_text(encoding='utf-8')
    node_test=NODE_TEST.read_text(encoding='utf-8') if NODE_TEST.is_file() else ''
    header_node_test=HEADER_NODE_TEST.read_text(encoding='utf-8') if HEADER_NODE_TEST.is_file() else ''
    browser_test=BROWSER_TEST.read_text(encoding='utf-8') if BROWSER_TEST.is_file() else ''
    hierarchy_test=HIERARCHY_TEST.read_text(encoding='utf-8') if HIERARCHY_TEST.is_file() else ''

    scripts=re.findall(r'<script\s+src="([^"]+)"',index)
    names=[x.split('?',1)[0] for x in scripts]
    required=['renderer.js','documentation-ui.js','market-header-ui.js','documentation-access-v18.6.js']
    for name in required:
        if names.count(name)!=1: errors.append(f'{name} must load exactly once; got {names.count(name)}')
    if 'header-v18.5.1.js' in names:
        errors.append('version-stacked header-v18.5.1.js must not remain an active runtime owner')
    if all(name in names for name in ('renderer.js','documentation-ui.js','documentation-access-v18.6.js')):
        if not (names.index('renderer.js')<names.index('documentation-ui.js')<names.index('documentation-access-v18.6.js')):
            errors.append('load order must be renderer.js -> documentation-ui.js -> documentation-access-v18.6.js')
    if all(name in names for name in ('renderer.js','market-header-ui.js')):
        if not names.index('renderer.js')<names.index('market-header-ui.js'):
            errors.append('market-header-ui.js must load after renderer.js so it decorates the canonical chrome owner')

    if 'documentation-ui-v' in index or 'documentation-ui-v' in owner:
        errors.append('new long-lived Documentation capability owner must not use version-stacked naming')
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

    if not HEADER.is_file():
        errors.append('stable Market Header capability owner is missing: renderer/market-header-ui.js')
    if 'market-header-ui-v' in index or 'market-header-ui-v' in header:
        errors.append('Market Header capability owner must remain release-neutral, not version-stacked')
    required_header=(
        "const OWNER='renderer/market-header-ui.js'",
        "state:'ACTIVE_OWNER_WITH_COMPAT_ALIAS'",
        'ensureSecondaryMarketStatus',
        'updateChrome = function updateChromeMarketHeader',
        'globalThis.__DE_PULSE_MARKET_HEADER_UI__',
        'globalThis.__DE_PULSE_RENDERER_OWNERS__',
        'registry.marketHeader',
        "compatibilityAliases:['__v1851HeaderContracts']",
        'globalThis.__v1851HeaderContracts=api',
        "'market-pulse-ribbon'",
        "'market-clocks'",
        "'data-runtime-control'",
    )
    for token in required_header:
        if token not in header: errors.append(f'Market Header owner contract missing: {token}')
    if LEGACY_HEADER.is_file() and 'legacyCompatibilityFile' not in header:
        errors.append('legacy header compatibility file must be explicitly classified by the active owner')
    if 'header-v18.5.1.js' not in hierarchy_test or 'ensureSecondaryMarketStatus' not in hierarchy_test:
        errors.append('historical header hierarchy regression fixture unexpectedly lost')

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

    evidence_tokens=(
        'run: node tests/renderer/documentation_ui_owner_test.js',
        'run: python3 tools/ci/documentation_owner_browser_test.py --engine chrome',
        'run: python3 tools/ci/documentation_owner_browser_test.py --engine webkit',
    )
    for token in evidence_tokens:
        if token not in qualified: errors.append(f'primary owner evidence missing from Qualified: {token}')
    if 'Documentation capability owner regression PASS' not in node_test:
        errors.append('Node Documentation owner regression proof missing')
    if "require('./market_header_owner_test.js')" not in node_test:
        errors.append('Qualified renderer owner regression must transitively execute canonical Market Header owner test')
    if not HEADER_NODE_TEST.is_file():
        errors.append('canonical Market Header behavior test is missing: tests/renderer/market_header_owner_test.js')
    header_test_tokens=(
        'Market Header capability owner regression PASS',
        "owner.owner,'renderer/market-header-ui.js'",
        "__v1851HeaderContracts===__DE_PULSE_MARKET_HEADER_UI__",
        "updateChrome('SPY')",
        'wrapper must not multiply base update calls',
        'ensure must be idempotent',
        'header-v18.5.1.js',
    )
    for token in header_test_tokens:
        if token not in header_node_test: errors.append(f'Market Header behavior regression proof missing: {token}')
    for token in ("choices=('chrome','webkit')", "data-render-owner=\"documentation-ui\"", "__DE_PULSE_DOCUMENTATION_UI__"):
        if token not in browser_test: errors.append(f'primary-engine browser owner proof missing: {token}')

    if errors: return fail(errors)
    print('DE.PULSE renderer owner contract: PASS')
    print('Documentation capability-oriented runtime owner: renderer/documentation-ui.js')
    print('Market Header capability-oriented runtime owner: renderer/market-header-ui.js')
    print('Market Header canonical behavior regression is transitively bound to Qualified renderer owner evidence: PASS')
    print('Market Header legacy v18.5.1 file retained only as historical compatibility evidence, not loaded by runtime: PASS')
    print('Market Header compatibility alias __v1851HeaderContracts preserved: PASS')
    print('load order owner before role-access decorator: PASS')
    print('Node Documentation + Market Header owner evidence binding: PASS')
    print('Chrome + WebKit Documentation owner evidence binding: PASS')
    print('legacy monolith Documentation fallback retained and explicitly gated for later physical deletion: PASS')
    print('version-stacked new-owner naming prevention: PASS')
    return 0


if __name__=='__main__':
    raise SystemExit(main())
